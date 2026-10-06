package x12835

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// WriteClaim writes one claim, its CLP through the last segment of its
// last service line, with the delimiters given: the inverse of ParseClaim.
func WriteClaim(w io.Writer, c ClaimCLP, d Delimiters) error {
	if d.Element == 0 || d.Component == 0 || d.Terminator == 0 {
		return fmt.Errorf("x12835: writing claim: the delimiters are not set")
	}
	x := writer{d: d}
	x.claim(c)
	_, err := w.Write(x.buf.Bytes())
	if err != nil {
		return fmt.Errorf("x12835: writing claim: %w", err)
	}
	return nil
}

// writer builds the output. Each method writes one segment or one loop; a
// loop's segments are collected as items and written in Index order, so
// typed segments and raw ones interleave as the file had them.
type writer struct {
	d   Delimiters
	buf bytes.Buffer
}

// item is one segment or loop waiting to be written at its position.
type item struct {
	index int
	write func()
}

// writeInOrder writes items by Index, keeping the order given for equal
// indexes (the guide's order, for a struct built by hand).
func (x *writer) writeInOrder(items []item) {
	sort.SliceStable(items, func(i, j int) bool { return items[i].index < items[j].index })
	for _, it := range items {
		it.write()
	}
}

// segment writes one segment: its ID, its elements with trailing empty
// ones left off, the terminator, and the suffix.
func (x *writer) segment(id string, elements ...string) {
	n := len(elements)
	for n > 0 && elements[n-1] == "" {
		n--
	}
	x.buf.WriteString(id)
	for _, e := range elements[:n] {
		x.buf.WriteByte(x.d.Element)
		x.buf.WriteString(e)
	}
	x.buf.WriteByte(x.d.Terminator)
	x.buf.WriteString(x.d.Suffix)
}

// composite joins components with the component separator, trailing empty
// ones left off.
func (x *writer) composite(parts ...string) string {
	n := len(parts)
	for n > 0 && parts[n-1] == "" {
		n--
	}
	return strings.Join(parts[:n], string(x.d.Component))
}

func (x *writer) raw(r RawSegment) {
	x.segment(r.ID, r.Elements...)
}

func (x *writer) raws(raws []RawSegment) []item {
	var items []item
	for _, r := range raws {
		r := r
		items = append(items, item{r.Index, func() { x.raw(r) }})
	}
	return items
}

// ---------------------------------------------------------------- values

// amountText is an amount as the file carried it, or, for one built by
// hand, its canonical spelling.
func amountText(a Amount) string {
	if !a.Set {
		return ""
	}
	if a.Text != "" {
		return a.Text
	}
	return a.Value.String()
}

// dateText is CCYYMMDD.
func dateText(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}

// interchangeDateText is ISA09's YYMMDD.
func interchangeDateText(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}

// timeText is HHMM, or HHMMSS when the element carried seconds.
func timeText(t Time) string {
	if !t.Set {
		return ""
	}
	if t.HasSeconds {
		return fmt.Sprintf("%02d%02d%02d", t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d%02d", t.Hour, t.Minute)
}

func countText(n int) string {
	return strconv.Itoa(n)
}

// ---------------------------------------------------------------- envelope

func (x *writer) interchange(isa InterchangeISA) {
	// The ISA is fixed width: all sixteen elements are written, empty or not.
	x.buf.WriteString("ISA")
	for _, e := range []string{
		isa.AuthQualifierISA01.String(), isa.AuthInfoISA02,
		isa.SecQualifierISA03.String(), isa.SecInfoISA04,
		isa.SenderIDQualifierISA05.String(), isa.SenderIDISA06,
		isa.ReceiverIDQualifierISA07.String(), isa.ReceiverIDISA08,
		interchangeDateText(isa.DateISA09), timeText(isa.TimeISA10),
		isa.RepetitionSepISA11, isa.VersionISA12, isa.ControlNumberISA13,
		isa.AckRequestedISA14.String(), isa.UsageISA15.String(), isa.ComponentSepISA16,
	} {
		x.buf.WriteByte(x.d.Element)
		x.buf.WriteString(e)
	}
	x.buf.WriteByte(x.d.Terminator)
	x.buf.WriteString(x.d.Suffix)

	var items []item
	for i := range isa.FunctionalGroups {
		g := isa.FunctionalGroups[i]
		items = append(items, item{g.Index, func() { x.group(g) }})
	}
	items = append(items, x.raws(isa.Raw)...)
	x.writeInOrder(items)
	x.segment("IEA", countText(isa.GroupCountIEA01), isa.ControlNumberIEA02)
}

func (x *writer) group(g FunctionalGroupGS) {
	x.segment("GS", g.FunctionalIDGS01.String(), g.SenderCodeGS02, g.ReceiverCodeGS03, dateText(g.DateGS04), timeText(g.TimeGS05),
		g.ControlNumberGS06, g.AgencyGS07.String(), g.VersionGS08)
	var items []item
	for i := range g.Transactions {
		t := g.Transactions[i]
		items = append(items, item{t.Index, func() { x.transaction(t) }})
	}
	items = append(items, x.raws(g.Raw)...)
	x.writeInOrder(items)
	x.segment("GE", countText(g.TransactionCountGE01), g.ControlNumberGE02)
}

func (x *writer) transaction(t TransactionST) {
	x.segment("ST", t.TransactionSetIDST01.String(), t.ControlNumberST02, t.ImplementationConventionST03)
	var items []item
	if t.Payment.Set {
		items = append(items, item{t.Payment.Index, func() { x.payment(t.Payment) }})
	}
	if t.Trace.Set {
		items = append(items, item{t.Trace.Index, func() { x.trace(t.Trace) }})
	}
	if t.Currency.Set {
		items = append(items, item{t.Currency.Index, func() { x.currency(t.Currency) }})
	}
	for i := range t.HeaderRefs {
		r := t.HeaderRefs[i]
		items = append(items, item{r.Index, func() { x.reference(r) }})
	}
	if t.ProductionDateDTM405.DateQualifierDTM01 != DateQualifierDTM01NotSent {
		items = append(items, item{t.ProductionDateDTM405.Index, func() { x.dateTime(t.ProductionDateDTM405) }})
	}
	if t.PayerN1.EntityIDN101 != EntityIDN101NotSent {
		items = append(items, item{t.PayerN1.Index, func() { x.party(t.PayerN1) }})
	}
	if t.PayeeN1.EntityIDN101 != EntityIDN101NotSent {
		items = append(items, item{t.PayeeN1.Index, func() { x.party(t.PayeeN1) }})
	}
	for i := range t.HeaderNumbers {
		lx := t.HeaderNumbers[i]
		items = append(items, item{lx.Index, func() { x.headerNumber(lx) }})
	}
	for i := range t.ProviderAdjustments {
		plb := t.ProviderAdjustments[i]
		items = append(items, item{plb.Index, func() { x.providerAdjustment(plb) }})
	}
	items = append(items, x.raws(t.Raw)...)
	x.writeInOrder(items)
	x.segment("SE", countText(t.SegmentCountSE01), t.ControlNumberSE02)
}

func (x *writer) payment(b PaymentBPR) {
	x.segment("BPR", b.HandlingBPR01.String(), amountText(b.TotalPaidBPR02), b.CreditDebitBPR03.String(), b.PayMethodBPR04.String(), b.PayFormatBPR05.String(),
		b.SenderDFIQualifierBPR06.String(), b.SenderDFIBPR07, b.SenderAcctQualifierBPR08.String(), b.SenderAcctBPR09,
		b.PayerIDBPR10, b.OrigCoSupCodeBPR11,
		b.ReceiverDFIQualifierBPR12.String(), b.ReceiverDFIBPR13, b.ReceiverAcctQualifierBPR14.String(), b.ReceiverAcctBPR15,
		dateText(b.EffectiveDateBPR16),
		b.BusinessFunctionBPR17, b.ThirdDFIQualifierBPR18, b.ThirdDFIBPR19, b.ThirdAcctQualifierBPR20, b.ThirdAcctBPR21)
}

func (x *writer) trace(t TraceTRN) {
	x.segment("TRN", t.TraceTypeTRN01.String(), t.TraceNumberTRN02, t.PayerIDTRN03, t.OrigCoSupCodeTRN04)
}

func (x *writer) currency(c CurrencyCUR) {
	x.segment("CUR", c.EntityIDCUR01, c.CurrencyCodeCUR02, amountText(c.ExchangeRateCUR03),
		c.EntityIDCUR04, c.CurrencyCodeCUR05, c.MarketCUR06,
		c.DateQualifierCUR07, c.DateCUR08, c.TimeCUR09,
		c.DateQualifierCUR10, c.DateCUR11, c.TimeCUR12,
		c.DateQualifierCUR13, c.DateCUR14, c.TimeCUR15,
		c.DateQualifierCUR16, c.DateCUR17, c.TimeCUR18,
		c.DateQualifierCUR19, c.DateCUR20, c.TimeCUR21)
}

func (x *writer) reference(r ReferenceREF) {
	x.segment("REF", r.RefQualifierREF01.String(), r.Value, r.DescriptionREF03, r.SubReferenceREF04)
}

func (x *writer) dateTime(d DateTimeDTM) {
	x.segment("DTM", d.DateQualifierDTM01.String(), dateText(d.Date), d.TimeDTM03, d.TimeCodeDTM04, d.PeriodFormatDTM05, d.PeriodDTM06)
}

// ---------------------------------------------------------------- parties

func (x *writer) party(n PartyN1) {
	x.segment("N1", n.EntityIDN101.String(), n.NameN102, n.IDQualifierN103.String(), n.IdentifierN104, n.RelationshipN105, n.EntityIDN106)
	var items []item
	if n.Address.Set {
		items = append(items, item{n.Address.Index, func() { x.address(n.Address) }})
	}
	if n.Location.Set {
		items = append(items, item{n.Location.Index, func() { x.location(n.Location) }})
	}
	for i := range n.References {
		r := n.References[i]
		items = append(items, item{r.Index, func() { x.reference(r) }})
	}
	for i := range n.Contacts {
		c := n.Contacts[i]
		items = append(items, item{c.Index, func() { x.contact(c) }})
	}
	if n.DeliveryRDM.Set {
		items = append(items, item{n.DeliveryRDM.Index, func() { x.delivery(n.DeliveryRDM) }})
	}
	items = append(items, x.raws(n.Raw)...)
	x.writeInOrder(items)
}

func (x *writer) address(a AddressN3) {
	x.segment("N3", a.Line1N301, a.Line2N302)
}

func (x *writer) location(l LocationN4) {
	x.segment("N4", l.CityN401, l.StateN402, l.ZipN403, l.CountryN404, l.LocationQualifierN405, l.LocationN406, l.CountrySubdivisionN407)
}

func (x *writer) contact(c ContactPER) {
	x.segment("PER", c.FunctionPER01.String(), c.NamePER02, c.QualifierPER03.String(), c.CommPER04, c.QualifierPER05.String(), c.CommPER06,
		c.QualifierPER07.String(), c.CommPER08, c.InquiryReferencePER09)
}

func (x *writer) delivery(d DeliveryRDM) {
	x.segment("RDM", d.ReportTransmissionRDM01.String(), d.NameRDM02, d.CommNumberRDM03, d.ReferenceRDM04, d.ReferenceRDM05)
}

// ---------------------------------------------------------------- header numbers

func (x *writer) headerNumber(lx HeaderNumberLX) {
	x.segment("LX", countText(lx.NumberLX01))
	var items []item
	if lx.ProviderSummary.Set {
		items = append(items, item{lx.ProviderSummary.Index, func() { x.providerSummary(lx.ProviderSummary) }})
	}
	if lx.InstitutionalSummary.Set {
		items = append(items, item{lx.InstitutionalSummary.Index, func() { x.institutionalSummary(lx.InstitutionalSummary) }})
	}
	for i := range lx.Claims {
		c := lx.Claims[i]
		items = append(items, item{c.Index, func() { x.claim(c) }})
	}
	items = append(items, x.raws(lx.Raw)...)
	x.writeInOrder(items)
}

func (x *writer) providerSummary(t ProviderSummaryTS3) {
	x.segment("TS3", t.ProviderIDTS301, t.FacilityCodeTS302, dateText(t.FiscalPeriodDateTS303),
		amountText(t.ClaimCountTS304), amountText(t.ClaimChangeAmountTS305), amountText(t.CoveredChargeAmountTS306),
		amountText(t.NoncoveredChargeAmountTS307), amountText(t.DeniedChargeAmountTS308), amountText(t.ProviderAmountTS309),
		amountText(t.InterestAmountTS310), amountText(t.ContractualAdjustmentAmountTS311), amountText(t.GrammRudmanReductionAmountTS312),
		amountText(t.MSPPayerAmountTS313), amountText(t.BloodDeductibleAmountTS314), amountText(t.NonLabChargeAmountTS315),
		amountText(t.CoinsuranceAmountTS316), amountText(t.HCPCSReportedChargeAmountTS317), amountText(t.HCPCSPayableAmountTS318),
		amountText(t.DeductibleAmountTS319), amountText(t.ProfessionalComponentAmountTS320), amountText(t.MSPPatientLiabilityMetAmountTS321),
		amountText(t.PatientReimbursementAmountTS322), amountText(t.PIPClaimCountTS323), amountText(t.PIPAdjustmentAmountTS324))
}

func (x *writer) institutionalSummary(t InstitutionalSummaryTS2) {
	x.segment("TS2", amountText(t.DRGAmountTS201), amountText(t.FederalSpecificAmountTS202), amountText(t.HospitalSpecificAmountTS203),
		amountText(t.DisproportionateAmountTS204), amountText(t.CapitalAmountTS205), amountText(t.IndirectMedicalEducationAmountTS206),
		amountText(t.OutlierDayCountTS207), amountText(t.DayOutlierAmountTS208), amountText(t.CostOutlierAmountTS209),
		amountText(t.AverageDRGLengthOfStayTS210), amountText(t.DischargeCountTS211), amountText(t.CostReportDayCountTS212),
		amountText(t.CoveredDayCountTS213), amountText(t.NoncoveredDayCountTS214), amountText(t.MSPPassThroughAmountTS215),
		amountText(t.AverageDRGWeightTS216), amountText(t.PPSCapitalFSPDRGAmountTS217), amountText(t.PSPCapitalHSPDRGAmountTS218),
		amountText(t.PPSDSHDRGAmountTS219))
}

// ---------------------------------------------------------------- claims

func (x *writer) claim(c ClaimCLP) {
	x.segment("CLP", c.PatientControlCLP01, c.StatusCLP02.String(), amountText(c.ChargeCLP03), amountText(c.PaymentCLP04), amountText(c.PatientRespCLP05),
		c.FilingIndicatorCLP06.String(), c.PayerClaimControlCLP07, c.FacilityCLP08, c.FrequencyCLP09, c.PatientStatusCLP10, c.DRGCLP11,
		amountText(c.DRGWeightCLP12), amountText(c.DischargeFractionCLP13), c.YesNoCLP14)
	var items []item
	for i := range c.Adjustments {
		a := c.Adjustments[i]
		items = append(items, item{a.Index, func() { x.adjustment(a) }})
	}
	for _, n := range []NameNM1{c.PatientNM1QC, c.CorrectedPatientNM174, c.RenderingNM182, c.CrossoverNM1TT, c.CorrectedPayerNM1PR} {
		if n.EntityIDNM101 != EntityIDNM101NotSent {
			n := n
			items = append(items, item{n.Index, func() { x.name(n) }})
		}
	}
	if c.InpatientMIA.Set {
		items = append(items, item{c.InpatientMIA.Index, func() { x.inpatient(c.InpatientMIA) }})
	}
	if c.OutpatientMOA.Set {
		items = append(items, item{c.OutpatientMOA.Index, func() { x.outpatient(c.OutpatientMOA) }})
	}
	for i := range c.References {
		r := c.References[i]
		items = append(items, item{r.Index, func() { x.reference(r) }})
	}
	for _, d := range []DateTimeDTM{c.StatementFromDTM232, c.StatementToDTM233, c.ReceivedDTM050, c.CoverageExpirationDTM036} {
		if d.DateQualifierDTM01 != DateQualifierDTM01NotSent {
			d := d
			items = append(items, item{d.Index, func() { x.dateTime(d) }})
		}
	}
	if c.ContactPER.FunctionPER01 != ContactFunctionPER01NotSent {
		items = append(items, item{c.ContactPER.Index, func() { x.contact(c.ContactPER) }})
	}
	for i := range c.Amounts {
		a := c.Amounts[i]
		items = append(items, item{a.Index, func() { x.amountSegment(a) }})
	}
	for i := range c.Quantities {
		q := c.Quantities[i]
		items = append(items, item{q.Index, func() { x.quantity(q) }})
	}
	for i := range c.ServiceLines {
		l := c.ServiceLines[i]
		items = append(items, item{l.Index, func() { x.serviceLine(l) }})
	}
	items = append(items, x.raws(c.Raw)...)
	x.writeInOrder(items)
}

func (x *writer) adjustment(a AdjustmentCAS) {
	elements := []string{a.GroupCAS01.String()}
	for _, r := range a.Reasons {
		elements = append(elements, r.Reason.String(), amountText(r.Amount), amountText(r.Quantity))
	}
	x.segment("CAS", elements...)
}

func (x *writer) name(n NameNM1) {
	x.segment("NM1", n.EntityIDNM101.String(), n.EntityTypeNM102.String(), n.LastNM103, n.FirstNM104, n.MiddleNM105, n.PrefixNM106, n.SuffixNM107,
		n.IDQualifierNM108.String(), n.IDNM109, n.RelationshipNM110, n.EntityIDNM111, n.LastNM112)
}

func (x *writer) inpatient(m InpatientMIA) {
	x.segment("MIA", amountText(m.CoveredDaysMIA01), amountText(m.OperatingOutlierMIA02), amountText(m.PsychDaysMIA03), amountText(m.DRGAmountMIA04),
		m.RemarkMIA05.String(), amountText(m.DSHAmountMIA06), amountText(m.MSPPassThruMIA07), amountText(m.PPSCapitalMIA08),
		amountText(m.CapitalFSPMIA09), amountText(m.CapitalHSPMIA10), amountText(m.CapitalDSHMIA11), amountText(m.OldCapitalMIA12),
		amountText(m.CapitalIMEMIA13), amountText(m.OperHSPSpecMIA14), amountText(m.CostReportDaysMIA15), amountText(m.OperFSPSpecMIA16),
		amountText(m.PPSOutlierMIA17), amountText(m.IndirectTeachingMIA18), amountText(m.NonPayProfCompMIA19),
		m.RemarkMIA20.String(), m.RemarkMIA21.String(), m.RemarkMIA22.String(), m.RemarkMIA23.String(), amountText(m.CapitalExceptionMIA24))
}

func (x *writer) outpatient(m OutpatientMOA) {
	x.segment("MOA", amountText(m.ReimbursementRateMOA01), amountText(m.HCPCSPayableMOA02),
		m.RemarkMOA03.String(), m.RemarkMOA04.String(), m.RemarkMOA05.String(), m.RemarkMOA06.String(), m.RemarkMOA07.String(),
		amountText(m.ESRDPaymentMOA08), amountText(m.NonpayableProfCompMOA09))
}

func (x *writer) amountSegment(a AmountAMT) {
	x.segment("AMT", a.AmountQualifierAMT01.String(), amountText(a.Amount), a.CreditDebitAMT03)
}

func (x *writer) quantity(q QuantityQTY) {
	x.segment("QTY", q.QuantityQualifierQTY01.String(), amountText(q.Amount), q.UnitOfMeasureQTY03, q.FreeFormQTY04)
}

func (x *writer) serviceLine(l ServiceLineSVC) {
	x.segment("SVC", x.procedure(l.ProcedureSVC01), amountText(l.ChargeSVC02), amountText(l.PaymentSVC03), l.RevenueCodeSVC04, amountText(l.UnitsSVC05),
		x.procedure(l.OriginalProcedureSVC06), amountText(l.OriginalUnitsSVC07))
	var items []item
	for _, d := range []DateTimeDTM{l.ServiceDateDTM472, l.ServiceDateDTM150, l.ServiceDateDTM151} {
		if d.DateQualifierDTM01 != DateQualifierDTM01NotSent {
			d := d
			items = append(items, item{d.Index, func() { x.dateTime(d) }})
		}
	}
	for i := range l.Adjustments {
		a := l.Adjustments[i]
		items = append(items, item{a.Index, func() { x.adjustment(a) }})
	}
	for _, r := range []ReferenceREF{l.LineControlREF6R, l.ServiceIDREFLU, l.RenderingREFHPI, l.PolicyREF0K} {
		if r.Value != "" {
			r := r
			items = append(items, item{r.Index, func() { x.reference(r) }})
		}
	}
	for i := range l.References {
		r := l.References[i]
		items = append(items, item{r.Index, func() { x.reference(r) }})
	}
	for i := range l.Amounts {
		a := l.Amounts[i]
		items = append(items, item{a.Index, func() { x.amountSegment(a) }})
	}
	for i := range l.Quantities {
		q := l.Quantities[i]
		items = append(items, item{q.Index, func() { x.quantity(q) }})
	}
	for i := range l.Remarks {
		r := l.Remarks[i]
		items = append(items, item{r.Index, func() { x.remark(r) }})
	}
	items = append(items, x.raws(l.Raw)...)
	x.writeInOrder(items)
}

// procedure writes the SVC01/SVC06 composite.
func (x *writer) procedure(p ProcedureSVC01) string {
	return x.composite(p.Qualifier.String(), p.Code, p.Modifiers[0], p.Modifiers[1], p.Modifiers[2], p.Modifiers[3], p.Description, p.Code2)
}

func (x *writer) remark(r RemarkLQ) {
	x.segment("LQ", r.RemarkQualifierLQ01.String(), r.RARC.String())
}

func (x *writer) providerAdjustment(p ProviderAdjustmentPLB) {
	elements := []string{p.ProviderPLB01, dateText(p.FiscalPeriodPLB02)}
	for _, a := range p.Adjustments {
		elements = append(elements, x.composite(a.Reason.String(), a.Identifier), amountText(a.Amount))
	}
	x.segment("PLB", elements...)
}
