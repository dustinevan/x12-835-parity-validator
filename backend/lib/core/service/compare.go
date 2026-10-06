package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/x12835"
)

// AmountMatchCents is the five-cent rule: two dollar amounts, both sent,
// whose exact difference is this many cents or less, either way, are the
// same amount. The comparison records no difference for them, so they
// leave no tag, and a pair whose only differences were such amounts is
// equal. The rule is for rounding: a percentage (the 2 percent
// sequestration reduction, 20 percent coinsurance) rounded per line
// instead of per claim, or truncated instead of rounded, moves an amount
// by a cent or a few (taxonomy T2.09, T3.16, T4.19). The rule covers money
// only: counts, units, days, rates, dates, and codes compare exactly, and
// an amount sent on one side only is still a difference. The summary
// cube's "match" payment band uses the same constant
// (lib/core/query/summary.go).
const AmountMatchCents = 5

// amountMatchLimit is AmountMatchCents as an exact decimal: 0.05.
var amountMatchLimit = x12835.Decimal{Coefficient: AmountMatchCents, Scale: 2}

// Difference is one element that differs between a MAC claim and the
// vendor claim matched to it. Values are the typed values as text: an
// amount as its exact decimal without trailing fractional zeros (80.00
// reads "80"), a date as CCYYMMDD, a code as sent. An empty value means the
// element is absent on that side. AmountDifference is vendor minus MAC,
// exact, for a dollar amount, with an absent side counted as zero (a
// missing CAS*CO*253 of 2.30 is -2.30), and always more than
// AmountMatchCents cents either way when both sides sent the amount; empty
// for anything that is not money (a code, a date, a count, an adjustment
// quantity).
type Difference struct {
	Element           string // the tag: 'CLP04', 'SVC03', 'CAS*CO*253', 'CAS*CO*253*QTY', 'MOA', 'RAW*K3'...
	Occurrence        int    // 1st, 2nd... of this element in the claim
	LineControlNumber string // REF*6R of the service line, for line-level elements
	MACValue          string
	VendorValue       string
	AmountDifference  string
}

// DiffClaims compares a MAC claim with its matched vendor claim under a
// compared set and returns every differing element. Each side's bytes are
// parsed with that side's delimiters (x12835.ParseClaim). Claim-level
// elements pair by label; service lines pair by REF*6R, and lines left over
// pair by position; a value that repeats under one label (the 2nd CAS*CO*45
// in a line) pairs with the same repeat on the other side. Dollar amounts
// compare under the five-cent rule (AmountMatchCents); everything else
// compares exactly.
func DiffClaims(set CompareSet, mac, vendor claim.Claim) ([]Difference, error) {
	m, err := parseRow(mac)
	if err != nil {
		return nil, err
	}
	v, err := parseRow(vendor)
	if err != nil {
		return nil, err
	}
	return diffClaims(set, m, v)
}

// parseRow parses a stored claim's bytes with its delimiters.
func parseRow(c claim.Claim) (x12835.ClaimCLP, error) {
	d, err := ParseDelimiters(c.Delimiters)
	if err != nil {
		return x12835.ClaimCLP{}, err
	}
	parsed, err := x12835.ParseClaim(c.Segments, d.Element, d.Component, d.Repetition, d.Segment)
	if err != nil {
		return x12835.ClaimCLP{}, fmt.Errorf("service: parsing %s payment %d claim %d: %w", c.FileName, c.PaymentNumber, c.ClaimNumber, err)
	}
	return parsed, nil
}

// diffClaims compares two parsed claims: the claim level, then each service
// line pair in pairLines order. The differences are then numbered by
// element within the claim.
func diffClaims(set CompareSet, mac, vendor x12835.ClaimCLP) ([]Difference, error) {
	d := differ{set: set, level: LevelClaim}
	d.claim(mac, vendor)

	d.level = LevelLine
	var none x12835.ServiceLineSVC
	for _, p := range pairLines(mac.ServiceLines, vendor.ServiceLines) {
		m, v := none, none
		d.line = ""
		if p.vendor >= 0 {
			v = vendor.ServiceLines[p.vendor]
			d.line = v.LineControlREF6R.Value
		}
		if p.mac >= 0 {
			m = mac.ServiceLines[p.mac]
			d.line = m.LineControlREF6R.Value
		}
		d.serviceLine(m, v)
	}
	if d.err != nil {
		return nil, d.err
	}

	seen := map[string]int{}
	for i := range d.out {
		seen[d.out[i].Element]++
		d.out[i].Occurrence = seen[d.out[i].Element]
	}
	return d.out, nil
}

// differ collects the differences of one claim pair, leaving out what the
// set leaves out at the level being compared.
type differ struct {
	set   CompareSet
	level string // LevelClaim or LevelLine
	line  string // REF*6R of the service line being compared; empty at the claim level
	out   []Difference
	err   error
}

func (d *differ) add(element, macValue, vendorValue, amountDifference string) {
	if d.set.leavesOut(d.level, element) {
		return
	}
	d.out = append(d.out, Difference{Element: element, LineControlNumber: d.line, MACValue: macValue, VendorValue: vendorValue, AmountDifference: amountDifference})
}

// text compares a text element.
func (d *differ) text(element, mac, vendor string) {
	if mac != vendor {
		d.add(element, mac, vendor, "")
	}
}

// amount compares a dollar amount under the five-cent rule (sameMoney).
// The difference carries vendor minus MAC.
func (d *differ) amount(element string, mac, vendor x12835.Amount) {
	if sameMoney(mac, vendor) {
		return
	}
	d.add(element, amountText(mac), amountText(vendor), d.amountDifference(element, mac, vendor))
}

// amountDifference is vendor minus MAC, exact. A side that did not send
// the element counts as zero, so an adjustment or an allowed amount on one
// side only shows its whole amount. A subtraction that does not fit a
// Decimal is recorded as the differ's error.
func (d *differ) amountDifference(element string, mac, vendor x12835.Amount) string {
	var m, v x12835.Decimal
	if mac.Set {
		m = mac.Value
	}
	if vendor.Set {
		v = vendor.Value
	}
	difference, err := v.Sub(m)
	if err != nil {
		d.err = errors.Join(d.err, fmt.Errorf("service: comparing %s: %w", element, err))
		return ""
	}
	return difference.String()
}

// count compares a number that is not money (a day count, units, a DRG
// weight, an adjustment quantity): no amount difference.
func (d *differ) count(element string, mac, vendor x12835.Amount) {
	if !sameAmount(mac, vendor) {
		d.add(element, amountText(mac), amountText(vendor), "")
	}
}

func (d *differ) date(element string, mac, vendor x12835.Date) {
	if mac != vendor {
		d.add(element, dateText(mac), dateText(vendor), "")
	}
}

// sameAmount is exact decimal equality, for a number that is not money:
// not sent equals not sent, and two sent values are equal when their
// Decimals are (ParseDecimal keeps no trailing fractional zeros, so 10.0
// and 10.00 are one Decimal).
func sameAmount(a, b x12835.Amount) bool {
	return a.Set == b.Set && (!a.Set || a.Value == b.Value)
}

// sameMoney is how two dollar amounts compare: not sent equals not sent,
// an amount sent on one side only differs, and two sent amounts are the
// same when they are equal or within AmountMatchCents cents of each other.
func sameMoney(mac, vendor x12835.Amount) bool {
	if mac.Set != vendor.Set {
		return false
	}
	return !mac.Set || mac.Value == vendor.Value || withinMatch(mac.Value, vendor.Value)
}

// withinMatch reports whether vendor minus MAC is between -0.05 and 0.05,
// edges included, computed exactly: 0.05 minus the difference and the
// difference plus 0.05 are both zero or more. A subtraction or addition
// that does not fit a Decimal means the amounts are far apart, not within.
func withinMatch(mac, vendor x12835.Decimal) bool {
	difference, err := vendor.Sub(mac)
	if err != nil {
		return false
	}
	above, err := amountMatchLimit.Sub(difference)
	if err != nil || above.Coefficient < 0 {
		return false
	}
	below, err := difference.Add(amountMatchLimit)
	return err == nil && below.Coefficient >= 0
}

func amountText(a x12835.Amount) string {
	if !a.Set {
		return ""
	}
	return a.Value.String()
}

// dateText writes a date as CCYYMMDD, the form the file carries.
func dateText(d x12835.Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}

// claim compares everything above the service lines.
func (d *differ) claim(m, v x12835.ClaimCLP) {
	d.text("CLP01", m.PatientControlCLP01, v.PatientControlCLP01)
	d.text("CLP02", statusText(m), statusText(v))
	d.amount("CLP03", m.ChargeCLP03, v.ChargeCLP03)
	d.amount("CLP04", m.PaymentCLP04, v.PaymentCLP04)
	d.amount("CLP05", m.PatientRespCLP05, v.PatientRespCLP05)
	d.text("CLP06", filingText(m), filingText(v))
	d.text("CLP07", m.PayerClaimControlCLP07, v.PayerClaimControlCLP07)
	d.text("CLP08", m.FacilityCLP08, v.FacilityCLP08)
	d.text("CLP09", m.FrequencyCLP09, v.FrequencyCLP09)
	d.text("CLP10", m.PatientStatusCLP10, v.PatientStatusCLP10)
	d.text("CLP11", m.DRGCLP11, v.DRGCLP11)
	d.count("CLP12", m.DRGWeightCLP12, v.DRGWeightCLP12)
	d.count("CLP13", m.DischargeFractionCLP13, v.DischargeFractionCLP13)
	d.text("CLP14", m.YesNoCLP14, v.YesNoCLP14)

	var macBuf, vendorBuf [16]keyedValue
	d.keyed("CAS*", appendAdjustments(macBuf[:0], m.Adjustments), appendAdjustments(vendorBuf[:0], v.Adjustments))

	d.text("NM1*QC", m.PatientNM1QC.IDNM109, v.PatientNM1QC.IDNM109)
	d.text("NM1*74", m.CorrectedPatientNM174.IDNM109, v.CorrectedPatientNM174.IDNM109)
	d.text("NM1*82", m.RenderingNM182.IDNM109, v.RenderingNM182.IDNM109)
	d.text("NM1*TT", m.CrossoverNM1TT.IDNM109, v.CrossoverNM1TT.IDNM109)
	d.text("NM1*PR", m.CorrectedPayerNM1PR.IDNM109, v.CorrectedPayerNM1PR.IDNM109)

	d.outpatient(m, v)
	d.inpatient(m, v)

	d.keyed("REF*", appendReferences(macBuf[:0], m.References), appendReferences(vendorBuf[:0], v.References))

	d.date("DTM*232", m.StatementFromDTM232.Date, v.StatementFromDTM232.Date)
	d.date("DTM*233", m.StatementToDTM233.Date, v.StatementToDTM233.Date)
	d.date("DTM*050", m.ReceivedDTM050.Date, v.ReceivedDTM050.Date)
	d.date("DTM*036", m.CoverageExpirationDTM036.Date, v.CoverageExpirationDTM036.Date)

	d.text("PER", contactText(m.ContactPER), contactText(v.ContactPER))

	d.keyed("AMT*", appendAmounts(macBuf[:0], m.Amounts), appendAmounts(vendorBuf[:0], v.Amounts))
	d.keyed("QTY*", appendQuantities(macBuf[:0], m.Quantities), appendQuantities(vendorBuf[:0], v.Quantities))
	d.keyed("RAW*", appendRaw(macBuf[:0], m.Raw), appendRaw(vendorBuf[:0], v.Raw))
}

// statusText is CLP02 as sent.
func statusText(c x12835.ClaimCLP) string {
	return c.StatusCLP02.String()
}

// filingText is CLP06 as sent.
func filingText(c x12835.ClaimCLP) string {
	return c.FilingIndicatorCLP06.String()
}

// contactText writes a claim-level PER as one value: its nine elements
// joined by "*", without trailing empty ones, coded elements as sent.
func contactText(p x12835.ContactPER) string {
	parts := [9]string{
		p.FunctionPER01.String(),
		p.NamePER02,
		p.QualifierPER03.String(),
		p.CommPER04,
		p.QualifierPER05.String(),
		p.CommPER06,
		p.QualifierPER07.String(),
		p.CommPER08,
		p.InquiryReferencePER09,
	}
	n := len(parts)
	for n > 0 && parts[n-1] == "" {
		n--
	}
	return strings.Join(parts[:n], "*")
}

// outpatient is the MOA segment: MOA01 as a rate, MOA02, MOA08, and MOA09
// as amounts, and the five remark codes as one ordered list labelled MOA.
// Order is kept: the exact comparison records a reordering, and an
// order-insensitive view is later analysis (taxonomy T3.36).
func (d *differ) outpatient(m, v x12835.ClaimCLP) {
	if !m.OutpatientMOA.Set && !v.OutpatientMOA.Set {
		return
	}
	d.count("MOA01", m.OutpatientMOA.ReimbursementRateMOA01, v.OutpatientMOA.ReimbursementRateMOA01)
	d.amount("MOA02", m.OutpatientMOA.HCPCSPayableMOA02, v.OutpatientMOA.HCPCSPayableMOA02)
	d.amount("MOA08", m.OutpatientMOA.ESRDPaymentMOA08, v.OutpatientMOA.ESRDPaymentMOA08)
	d.amount("MOA09", m.OutpatientMOA.NonpayableProfCompMOA09, v.OutpatientMOA.NonpayableProfCompMOA09)
	mr := remarkTexts(outpatientRemarks(m.OutpatientMOA))
	vr := remarkTexts(outpatientRemarks(v.OutpatientMOA))
	if mr != vr {
		d.add("MOA", joinRemarks(mr), joinRemarks(vr), "")
	}
}

func outpatientRemarks(m x12835.OutpatientMOA) [5]rarc.Code {
	return [5]rarc.Code{m.RemarkMOA03, m.RemarkMOA04, m.RemarkMOA05, m.RemarkMOA06, m.RemarkMOA07}
}

// inpatient is the MIA segment: its amounts, its three counts (MIA01,
// MIA03, MIA15), and the five remark codes (MIA05, MIA20-23) as one ordered
// list labelled MIA.
func (d *differ) inpatient(m, v x12835.ClaimCLP) {
	if !m.InpatientMIA.Set && !v.InpatientMIA.Set {
		return
	}
	mm, vm := m.InpatientMIA, v.InpatientMIA
	d.count("MIA01", mm.CoveredDaysMIA01, vm.CoveredDaysMIA01)
	d.amount("MIA02", mm.OperatingOutlierMIA02, vm.OperatingOutlierMIA02)
	d.count("MIA03", mm.PsychDaysMIA03, vm.PsychDaysMIA03)
	d.amount("MIA04", mm.DRGAmountMIA04, vm.DRGAmountMIA04)
	d.amount("MIA06", mm.DSHAmountMIA06, vm.DSHAmountMIA06)
	d.amount("MIA07", mm.MSPPassThruMIA07, vm.MSPPassThruMIA07)
	d.amount("MIA08", mm.PPSCapitalMIA08, vm.PPSCapitalMIA08)
	d.amount("MIA09", mm.CapitalFSPMIA09, vm.CapitalFSPMIA09)
	d.amount("MIA10", mm.CapitalHSPMIA10, vm.CapitalHSPMIA10)
	d.amount("MIA11", mm.CapitalDSHMIA11, vm.CapitalDSHMIA11)
	d.amount("MIA12", mm.OldCapitalMIA12, vm.OldCapitalMIA12)
	d.amount("MIA13", mm.CapitalIMEMIA13, vm.CapitalIMEMIA13)
	d.amount("MIA14", mm.OperHSPSpecMIA14, vm.OperHSPSpecMIA14)
	d.count("MIA15", mm.CostReportDaysMIA15, vm.CostReportDaysMIA15)
	d.amount("MIA16", mm.OperFSPSpecMIA16, vm.OperFSPSpecMIA16)
	d.amount("MIA17", mm.PPSOutlierMIA17, vm.PPSOutlierMIA17)
	d.amount("MIA18", mm.IndirectTeachingMIA18, vm.IndirectTeachingMIA18)
	d.amount("MIA19", mm.NonPayProfCompMIA19, vm.NonPayProfCompMIA19)
	d.amount("MIA24", mm.CapitalExceptionMIA24, vm.CapitalExceptionMIA24)
	mr := remarkTexts(inpatientRemarks(mm))
	vr := remarkTexts(inpatientRemarks(vm))
	if mr != vr {
		d.add("MIA", joinRemarks(mr), joinRemarks(vr), "")
	}
}

func inpatientRemarks(m x12835.InpatientMIA) [5]rarc.Code {
	return [5]rarc.Code{m.RemarkMIA05, m.RemarkMIA20, m.RemarkMIA21, m.RemarkMIA22, m.RemarkMIA23}
}

// serviceLine compares one MAC service line with the vendor line paired
// with it (or with an empty line when it has no partner).
func (d *differ) serviceLine(m, v x12835.ServiceLineSVC) {
	d.procedure("SVC01", m.ProcedureSVC01, v.ProcedureSVC01)
	d.amount("SVC02", m.ChargeSVC02, v.ChargeSVC02)
	d.amount("SVC03", m.PaymentSVC03, v.PaymentSVC03)
	d.text("SVC04", m.RevenueCodeSVC04, v.RevenueCodeSVC04)
	d.count("SVC05", m.UnitsSVC05, v.UnitsSVC05)
	d.procedure("SVC06", m.OriginalProcedureSVC06, v.OriginalProcedureSVC06)
	d.count("SVC07", m.OriginalUnitsSVC07, v.OriginalUnitsSVC07)

	d.date("DTM*472", m.ServiceDateDTM472.Date, v.ServiceDateDTM472.Date)
	d.date("DTM*150", m.ServiceDateDTM150.Date, v.ServiceDateDTM150.Date)
	d.date("DTM*151", m.ServiceDateDTM151.Date, v.ServiceDateDTM151.Date)

	var macBuf, vendorBuf [16]keyedValue
	d.keyed("CAS*", appendAdjustments(macBuf[:0], m.Adjustments), appendAdjustments(vendorBuf[:0], v.Adjustments))

	macRefs, vendorRefs := macBuf[:0], vendorBuf[:0]
	macRefs = appendNamed(macRefs, "6R", m.LineControlREF6R.Value)
	macRefs = appendNamed(macRefs, "LU", m.ServiceIDREFLU.Value)
	macRefs = appendNamed(macRefs, "HPI", m.RenderingREFHPI.Value)
	macRefs = appendNamed(macRefs, "0K", m.PolicyREF0K.Value)
	macRefs = appendReferences(macRefs, m.References)
	vendorRefs = appendNamed(vendorRefs, "6R", v.LineControlREF6R.Value)
	vendorRefs = appendNamed(vendorRefs, "LU", v.ServiceIDREFLU.Value)
	vendorRefs = appendNamed(vendorRefs, "HPI", v.RenderingREFHPI.Value)
	vendorRefs = appendNamed(vendorRefs, "0K", v.PolicyREF0K.Value)
	vendorRefs = appendReferences(vendorRefs, v.References)
	d.keyed("REF*", macRefs, vendorRefs)

	d.keyed("AMT*", appendAmounts(macBuf[:0], m.Amounts), appendAmounts(vendorBuf[:0], v.Amounts))
	d.keyed("QTY*", appendQuantities(macBuf[:0], m.Quantities), appendQuantities(vendorBuf[:0], v.Quantities))
	d.lists(appendRemarks(macBuf[:0], m.Remarks), appendRemarks(vendorBuf[:0], v.Remarks))
	d.keyed("RAW*", appendRaw(macBuf[:0], m.Raw), appendRaw(vendorBuf[:0], v.Raw))
}

// procedure is SVC01 or SVC06, a composite compared as one value:
// qualifier, code, modifiers, description, second code.
func (d *differ) procedure(element string, m, v x12835.ProcedureSVC01) {
	if m == v {
		return
	}
	d.text(element, procedureText(m), procedureText(v))
}

// procedureText writes a composite with ":" between components and without
// trailing empty ones.
func procedureText(p x12835.ProcedureSVC01) string {
	parts := [8]string{
		p.Qualifier.String(),
		p.Code, p.Modifiers[0], p.Modifiers[1], p.Modifiers[2], p.Modifiers[3], p.Description, p.Code2,
	}
	n := len(parts)
	for n > 0 && parts[n-1] == "" {
		n--
	}
	return strings.Join(parts[:n], ":")
}

// keyedValue is one value of a segment that repeats and is told apart by a
// qualifier: a REF, AMT, or QTY value, a CAS reason's amount or quantity,
// or an unplaced (Raw) segment. The nth value with a label on one side pairs
// with the nth on the other.
type keyedValue struct {
	key      string        // the qualifier; the group code for a CAS reason; the segment ID for Raw
	key2     string        // the reason code for a CAS reason; empty otherwise
	quantity bool          // a CAS reason's quantity, labelled *QTY after the reason
	text     string        // the value, for REF and Raw
	amount   x12835.Amount // the value, for AMT, QTY, and CAS
	money    bool          // amount is a dollar amount: a difference carries vendor minus MAC
}

func (e keyedValue) label(prefix string) string {
	label := prefix + e.key
	if e.key2 != "" {
		label += "*" + e.key2
	}
	if e.quantity {
		label += "*QTY"
	}
	return label
}

func (e keyedValue) valueText() string {
	if e.amount.Set {
		return amountText(e.amount)
	}
	return e.text
}

// absent is what an unpaired value compares with: the same label and
// kind, with nothing sent.
func (e keyedValue) absent() keyedValue {
	return keyedValue{key: e.key, key2: e.key2, quantity: e.quantity, money: e.money}
}

func sameLabel(a, b keyedValue) bool {
	return a.key == b.key && a.key2 == b.key2 && a.quantity == b.quantity
}

// keyed compares two lists of keyed values: each value with the value of
// the same label and occurrence on the other side, or with an absent value.
func (d *differ) keyed(prefix string, mac, vendor []keyedValue) {
	for i := range mac {
		j := keyedPartner(mac, i, vendor)
		if j < 0 {
			d.keyedValues(prefix, mac[i], mac[i].absent())
			continue
		}
		d.keyedValues(prefix, mac[i], vendor[j])
	}
	for j := range vendor {
		if keyedPartner(vendor, j, mac) < 0 {
			d.keyedValues(prefix, vendor[j].absent(), vendor[j])
		}
	}
}

// keyedValues compares one keyed value with its partner: a dollar amount
// (a CAS amount, an AMT) under the five-cent rule, anything else exactly.
func (d *differ) keyedValues(prefix string, mac, vendor keyedValue) {
	same := sameAmount(mac.amount, vendor.amount)
	if mac.money || vendor.money {
		same = sameMoney(mac.amount, vendor.amount)
	}
	if mac.text == vendor.text && same {
		return
	}
	element := mac.label(prefix)
	difference := ""
	if mac.money || vendor.money {
		difference = d.amountDifference(element, mac.amount, vendor.amount)
	}
	d.add(element, mac.valueText(), vendor.valueText(), difference)
}

// keyedPartner returns the index in in of the value that pairs with
// from[i]: the same label, and the same occurrence of it. It returns -1
// when there is none.
func keyedPartner(from []keyedValue, i int, in []keyedValue) int {
	nth := 0
	for k := 0; k < i; k++ {
		if sameLabel(from[k], from[i]) {
			nth++
		}
	}
	for j := range in {
		if sameLabel(in[j], from[i]) {
			if nth == 0 {
				return j
			}
			nth--
		}
	}
	return -1
}

// lists compares remark codes: for each LQ qualifier, one ordered list of
// codes, labelled LQ for HE and LQ*qualifier otherwise.
func (d *differ) lists(mac, vendor []keyedValue) {
	for i := range mac {
		if firstOfKey(mac, i) {
			d.list(mac[i].key, mac, vendor)
		}
	}
	for j := range vendor {
		if firstOfKey(vendor, j) && !hasKey(mac, vendor[j].key) {
			d.list(vendor[j].key, mac, vendor)
		}
	}
}

func firstOfKey(list []keyedValue, i int) bool {
	for k := 0; k < i; k++ {
		if list[k].key == list[i].key {
			return false
		}
	}
	return true
}

func hasKey(list []keyedValue, key string) bool {
	for i := range list {
		if list[i].key == key {
			return true
		}
	}
	return false
}

// list compares the codes with one qualifier on each side, in order.
func (d *differ) list(key string, mac, vendor []keyedValue) {
	a, b := nextOfKey(mac, key, 0), nextOfKey(vendor, key, 0)
	same := true
	for a >= 0 || b >= 0 {
		if a < 0 || b < 0 || mac[a].text != vendor[b].text {
			same = false
			break
		}
		a, b = nextOfKey(mac, key, a+1), nextOfKey(vendor, key, b+1)
	}
	if same {
		return
	}
	label := "LQ"
	if key != "HE" {
		label = "LQ*" + key
	}
	d.add(label, joinOfKey(mac, key), joinOfKey(vendor, key), "")
}

func nextOfKey(list []keyedValue, key string, from int) int {
	for i := from; i < len(list); i++ {
		if list[i].key == key {
			return i
		}
	}
	return -1
}

func joinOfKey(list []keyedValue, key string) string {
	var codes []string
	for i := range list {
		if list[i].key == key {
			codes = append(codes, list[i].text)
		}
	}
	return strings.Join(codes, " ")
}

// adjustmentValue is one CAS reason as sent: its group and reason as text,
// its amount, and its quantity.
type adjustmentValue struct {
	group, reason    string
	amount, quantity x12835.Amount
}

// adjustmentValues lists every CAS reason sent, in file order, group and
// reason as the file spelled them.
func adjustmentValues(adjustments []x12835.AdjustmentCAS) []adjustmentValue {
	var out []adjustmentValue
	for _, a := range adjustments {
		for _, r := range a.Reasons {
			if r.Reason == carc.CodeNotSent {
				continue
			}
			out = append(out, adjustmentValue{group: a.GroupCAS01.String(), reason: r.Reason.String(), amount: r.Amount, quantity: r.Quantity})
		}
	}
	return out
}

// appendAdjustments adds two values per CAS reason sent, keyed by group and
// reason: the amount, and the quantity (labelled *QTY).
func appendAdjustments(entries []keyedValue, adjustments []x12835.AdjustmentCAS) []keyedValue {
	for _, a := range adjustmentValues(adjustments) {
		entries = append(entries,
			keyedValue{key: a.group, key2: a.reason, amount: a.amount, money: true},
			keyedValue{key: a.group, key2: a.reason, quantity: true, amount: a.quantity})
	}
	return entries
}

// appendReferences adds one value per REF, keyed by qualifier: REF02, and
// after it, joined by "*", a REF03 or REF04 when the file sent one (the
// guide marks those not used, so a file that sends them is compared on
// them too).
func appendReferences(entries []keyedValue, refs []x12835.ReferenceREF) []keyedValue {
	for _, r := range refs {
		text := r.Value
		if r.DescriptionREF03 != "" || r.SubReferenceREF04 != "" {
			text = strings.TrimRight(strings.Join([]string{r.Value, r.DescriptionREF03, r.SubReferenceREF04}, "*"), "*")
		}
		entries = append(entries, keyedValue{key: r.RefQualifierREF01.String(), text: text})
	}
	return entries
}

// appendNamed adds a service line REF that the model holds in a field named
// for its qualifier, when the field has a value.
func appendNamed(entries []keyedValue, key, value string) []keyedValue {
	if value == "" {
		return entries
	}
	return append(entries, keyedValue{key: key, text: value})
}

func appendAmounts(entries []keyedValue, amounts []x12835.AmountAMT) []keyedValue {
	for _, a := range amounts {
		entries = append(entries, keyedValue{key: a.AmountQualifierAMT01.String(), amount: a.Amount, money: true})
	}
	return entries
}

func appendQuantities(entries []keyedValue, quantities []x12835.QuantityQTY) []keyedValue {
	for _, q := range quantities {
		entries = append(entries, keyedValue{key: q.QuantityQualifierQTY01.String(), amount: q.Amount})
	}
	return entries
}

// appendRemarks adds one value per LQ with a code, keyed by qualifier.
func appendRemarks(entries []keyedValue, remarks []x12835.RemarkLQ) []keyedValue {
	for _, r := range remarks {
		if r.RARC == rarc.CodeNotSent {
			continue
		}
		entries = append(entries, keyedValue{key: r.RemarkQualifierLQ01.String(), text: r.RARC.String()})
	}
	return entries
}

// appendRaw adds one value per segment the model did not place, keyed by
// segment ID (labelled RAW*ID, so it never collides with a modeled tag),
// with its elements joined by "*".
func appendRaw(entries []keyedValue, raws []x12835.RawSegment) []keyedValue {
	for _, r := range raws {
		entries = append(entries, keyedValue{key: r.ID, text: strings.Join(r.Elements, "*")})
	}
	return entries
}

// remarkTexts lists the remark codes a segment sent, in order, as text,
// with unsent positions removed (the rest move up).
func remarkTexts(codes [5]rarc.Code) [5]string {
	var out [5]string
	n := 0
	for _, code := range codes {
		if code == rarc.CodeNotSent {
			continue
		}
		out[n] = code.String()
		n++
	}
	return out
}

// joinRemarks writes remark codes joined by one space.
func joinRemarks(codes [5]string) string {
	n := len(codes)
	for n > 0 && codes[n-1] == "" {
		n--
	}
	return strings.Join(codes[:n], " ")
}

// linePair is one MAC line and the vendor line compared with it; -1 means
// the line has no partner on that side.
type linePair struct {
	mac, vendor int
}

// pairLines pairs lines whose REF*6R is present and unique on both sides,
// then pairs the remaining lines in order, then lists what is left.
func pairLines(mac, vendor []x12835.ServiceLineSVC) []linePair {
	macPartner := make([]int, len(mac))
	taken := make([]bool, len(vendor))
	for i := range mac {
		macPartner[i] = -1
		control := mac[i].LineControlREF6R.Value
		if control == "" || countControl(mac, control) != 1 || countControl(vendor, control) != 1 {
			continue
		}
		j := indexControl(vendor, control)
		macPartner[i] = j
		taken[j] = true
	}
	next := 0
	for i := range mac {
		if macPartner[i] >= 0 {
			continue
		}
		for next < len(vendor) && taken[next] {
			next++
		}
		if next < len(vendor) {
			macPartner[i] = next
			taken[next] = true
		}
	}
	out := make([]linePair, 0, len(mac)+len(vendor))
	for i := range mac {
		out = append(out, linePair{mac: i, vendor: macPartner[i]})
	}
	for j := range vendor {
		if !taken[j] {
			out = append(out, linePair{mac: -1, vendor: j})
		}
	}
	return out
}

func countControl(lines []x12835.ServiceLineSVC, control string) int {
	n := 0
	for i := range lines {
		if lines[i].LineControlREF6R.Value == control {
			n++
		}
	}
	return n
}

func indexControl(lines []x12835.ServiceLineSVC, control string) int {
	for i := range lines {
		if lines[i].LineControlREF6R.Value == control {
			return i
		}
	}
	return -1
}
