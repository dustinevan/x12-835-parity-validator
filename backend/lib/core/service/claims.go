package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/payment"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/macs"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/x12835"
)

// Claim types.
const (
	Institutional = "institutional"
	Professional  = "professional"
	DME           = "dme"
)

// Delimiters are a file's four separators, read from its ISA segment and
// stored with each claim so the claim's bytes can be parsed on their own.
type Delimiters struct {
	Element    byte
	Component  byte
	Repetition byte
	Segment    byte
}

// String is the stored form: element, component, repetition, segment
// terminator, one byte each.
func (d Delimiters) String() string {
	return string([]byte{d.Element, d.Component, d.Repetition, d.Segment})
}

func ParseDelimiters(s string) (Delimiters, error) {
	if len(s) != 4 {
		return Delimiters{}, fmt.Errorf("service: parsing delimiters %q: want 4 bytes", s)
	}
	return Delimiters{Element: s[0], Component: s[1], Repetition: s[2], Segment: s[3]}, nil
}

// delimitersOf reads the four separators from a file's ISA: the element
// separator is the byte after "ISA", the repetition separator is ISA11, the
// component separator is ISA16, and the segment terminator is the byte
// after ISA16. x12835.Parse reads the same bytes the same way; this is the
// file's one ISA, since Parse accepts one interchange.
func delimitersOf(data []byte) (Delimiters, error) {
	start := 0
	for start < len(data) && isSpace(data[start]) {
		start++
	}
	if !bytes.HasPrefix(data[start:], []byte("ISA")) || len(data) < start+4 {
		return Delimiters{}, fmt.Errorf("service: reading delimiters: expected ISA")
	}
	element := data[start+3]
	separators := 0
	isa16 := -1
	for i := start + 3; i < len(data); i++ {
		if data[i] != element {
			continue
		}
		separators++
		if separators == 16 {
			isa16 = i + 1
			break
		}
	}
	if isa16 < 0 || isa16+1 >= len(data) {
		return Delimiters{}, fmt.Errorf("service: reading delimiters: incomplete ISA")
	}
	fields := bytes.Split(data[start:isa16], []byte{element})
	if len(fields[11]) != 1 {
		return Delimiters{}, fmt.Errorf("service: reading delimiters: ISA11 is not one byte")
	}
	return Delimiters{Element: element, Component: data[isa16], Repetition: fields[11][0], Segment: data[isa16+1]}, nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}

// Key is the six fields a MAC claim and a vendor claim are matched on.
type Key struct {
	MBI           string // 2100 NM1*QC NM109
	DateOfService string // YYYYMMDD: DTM*232, else the earliest 2110 DTM*472
	TotalCharge   string // CLP03
	TaxID         string // 1000B REF*TJ
	NPI           string // 1000B N1*PE N104
	Workload      string // 1000A N1*PR N104
}

// normalizeKey puts each field in one spelling, so that two systems that
// agree on the facts produce the same bytes:
//
//	MBI            trimmed, upper case
//	DateOfService  must be a valid CCYYMMDD date; kept as CCYYMMDD
//	TotalCharge    exact decimal with trailing fractional zeros removed:
//	               "447.00" -> "447", "1469.50" -> "1469.5", "0.00" -> "0"
//	TaxID          trimmed, hyphens removed
//	NPI, Workload  trimmed
func normalizeKey(k Key) (Key, error) {
	var n Key
	n.MBI = strings.ToUpper(strings.TrimSpace(k.MBI))
	n.DateOfService = strings.TrimSpace(k.DateOfService)
	if _, err := time.Parse("20060102", n.DateOfService); err != nil {
		return Key{}, fmt.Errorf("service: normalizing key: date of service %q: %w", k.DateOfService, err)
	}
	charge, err := x12835.ParseDecimal([]byte(strings.TrimSpace(k.TotalCharge)))
	if err != nil {
		return Key{}, fmt.Errorf("service: normalizing key: total charge %q: %w", k.TotalCharge, err)
	}
	n.TotalCharge = charge.String()
	n.TaxID = strings.ReplaceAll(strings.TrimSpace(k.TaxID), "-", "")
	n.NPI = strings.TrimSpace(k.NPI)
	n.Workload = strings.TrimSpace(k.Workload)
	return n, nil
}

// MatchKey is the SHA-256 of the six normalized key fields. The hashed
// bytes are the fields in this order, UTF-8, separated by one 0x1F byte
// (ASCII unit separator, which X12 data cannot contain), with no
// separator before the first or after the last:
//
//	MBI 0x1F DateOfService 0x1F TotalCharge 0x1F TaxID 0x1F NPI 0x1F Workload
//
// In SQL: sha256(convert_to(concat_ws(chr(31), mbi, dos,
// trim_scale(charge)::text, tax_id, npi, workload), 'UTF8')).
func MatchKey(k Key) ([]byte, error) {
	n, err := normalizeKey(k)
	if err != nil {
		return nil, err
	}
	text := strings.Join([]string{n.MBI, n.DateOfService, n.TotalCharge, n.TaxID, n.NPI, n.Workload}, "\x1f")
	sum := sha256.Sum256([]byte(text))
	return sum[:], nil
}

// claimType classifies a claim from its typed fields and its payer:
//
//	institutional  the claim has a statement date (DTM*232 or DTM*233), an
//	               MIA segment, or a service line with a revenue code (SVC04)
//	dme            otherwise, when the payer name (1000A N1*PR N102) contains
//	               "DME", as DME MAC payer names do ("SYNTHETIC DME MAC JC")
//	professional   otherwise
func claimType(c x12835.ClaimCLP, payerName string) string {
	if c.StatementFromDTM232.Date.Set || c.StatementToDTM233.Date.Set || c.InpatientMIA.Set {
		return Institutional
	}
	for i := range c.ServiceLines {
		if c.ServiceLines[i].RevenueCodeSVC04 != "" {
			return Institutional
		}
	}
	if strings.Contains(strings.ToUpper(payerName), "DME") {
		return DME
	}
	return Professional
}

// dateOfService is DTM*232 (statement from), else the earliest service
// line date (DTM*472, or DTM*150 when a line gives a period).
func dateOfService(c x12835.ClaimCLP) x12835.Date {
	if c.StatementFromDTM232.Date.Set {
		return c.StatementFromDTM232.Date
	}
	var earliest x12835.Date
	for i := range c.ServiceLines {
		for _, d := range [2]x12835.Date{c.ServiceLines[i].ServiceDateDTM472.Date, c.ServiceLines[i].ServiceDateDTM150.Date} {
			if d.Set && (!earliest.Set || dateNumber(d) < dateNumber(earliest)) {
				earliest = d
			}
		}
	}
	return earliest
}

// dateNumber is a date as CCYYMMDD in one integer, so dates order as their
// CCYYMMDD text does.
func dateNumber(d x12835.Date) int {
	return d.Year*10000 + d.Month*100 + d.Day
}

// payee is what a payment's header says about everyone in it.
type payee struct {
	workload  string // 1000A N1*PR N104
	payerName string // 1000A N1*PR N102
	npi       string // 1000B N1*PE N104
	taxID     string // the last REF*TJ among the header's, the payer's, and the payee's references
}

func payeeOf(t x12835.TransactionST) payee {
	p := payee{workload: t.PayerN1.IdentifierN104, payerName: t.PayerN1.NameN102, npi: t.PayeeN1.IdentifierN104}
	for _, refs := range [][]x12835.ReferenceREF{t.HeaderRefs, t.PayerN1.References, t.PayeeN1.References} {
		for _, r := range refs {
			if r.RefQualifierREF01 == x12835.RefQualifierREF01TJ {
				p.taxID = r.Value
			}
		}
	}
	return p
}

// claimSpan locates one claim's bytes in a file: from the first byte of its
// CLP through the terminator of its last segment, and where the claim sits:
// the payment (ST) it is in, counting from 1, and its number within that
// payment.
type claimSpan struct {
	start, end     int
	payment, claim int
}

// endsClaim reports whether a segment ends the claim before it. A claim
// runs from its CLP through the segment before the next CLP, LX, PLB, or
// envelope segment: the boundaries x12835.Parse gives a claim.
func endsClaim(tag string) bool {
	switch tag {
	case "CLP", "LX", "PLB", "ISA", "IEA", "GS", "GE", "ST", "SE":
		return true
	}
	return false
}

// claimSpans finds every claim's bytes in a file by scanning its segments
// with the file's delimiters. Whitespace between segments belongs to no
// claim, so a span starts at the CLP's first byte and ends right after its
// last segment's terminator.
func claimSpans(data []byte, d Delimiters) []claimSpan {
	var spans []claimSpan
	open := false
	pos := 0
	last := 0 // the byte after the last segment's terminator
	payment, claim := 0, 0
	for pos < len(data) {
		for pos < len(data) && isSpace(data[pos]) {
			pos++
		}
		if pos == len(data) {
			break
		}
		n := bytes.IndexByte(data[pos:], d.Segment)
		if n < 0 {
			n = len(data) - pos
		} else {
			n++
		}
		tag := data[pos : pos+n]
		if i := bytes.IndexByte(tag, d.Element); i >= 0 {
			tag = tag[:i]
		}
		tag = bytes.TrimRight(tag, " \t\r\n"+string(d.Segment))
		if open && endsClaim(string(tag)) {
			spans[len(spans)-1].end = last
			open = false
		}
		switch string(tag) {
		case "ST":
			payment++
			claim = 0
		case "CLP":
			claim++
			spans = append(spans, claimSpan{start: pos, payment: payment, claim: claim})
			open = true
		}
		pos += n
		last = pos
	}
	if open {
		spans[len(spans)-1].end = last
	}
	return spans
}

// Split reads one 835 file into the typed model and returns a claim row for
// every claim in it, in file order: the key fields from the typed claim and
// its payment's payer and payee, and the claim's exact bytes (from its CLP
// through its last segment's terminator); and a payment row for every
// transaction (ST..SE) in it, in file order: the header the transaction's
// claims sit under (paymentRow). A malformed amount, date, or count
// anywhere in the file fails the file; the error names the file and, when
// the fault is inside a claim, the payment and claim.
func Split(submissionPrefix, fileName string, data []byte) ([]claim.Claim, []payment.Payment, error) {
	d, err := delimitersOf(data)
	if err != nil {
		return nil, nil, fmt.Errorf("service: splitting %s: %w", fileName, err)
	}
	spans := claimSpans(data, d)
	isa, err := x12835.Parse(data)
	if err != nil {
		return nil, nil, fmt.Errorf("service: splitting %s: %w", fileName, locateFault(err, data, d, spans))
	}
	typed := 0
	transactions := 0
	for _, group := range isa.FunctionalGroups {
		for _, transaction := range group.Transactions {
			transactions++
			for _, lx := range transaction.HeaderNumbers {
				typed += len(lx.Claims)
			}
		}
	}
	if typed != len(spans) {
		return nil, nil, fmt.Errorf("service: splitting %s: the typed model holds %d claims, the file has %d CLP segments", fileName, typed, len(spans))
	}

	out := make([]claim.Claim, 0, typed)
	payments := make([]payment.Payment, 0, transactions)
	paymentNumber := 0
	for _, group := range isa.FunctionalGroups {
		for _, transaction := range group.Transactions {
			paymentNumber++
			p := payeeOf(transaction)
			header, err := paymentRow(submissionPrefix, fileName, paymentNumber, transaction, p)
			if err != nil {
				return nil, nil, fmt.Errorf("service: splitting %s: payment %d: %w", fileName, paymentNumber, err)
			}
			payments = append(payments, header)
			claimNumber := 0
			for _, lx := range transaction.HeaderNumbers {
				for i := range lx.Claims {
					claimNumber++
					span := spans[len(out)]
					row, err := claimRow(submissionPrefix, fileName, paymentNumber, claimNumber, lx.Claims[i], p, d, data[span.start:span.end])
					if err != nil {
						return nil, nil, fmt.Errorf("service: splitting %s: payment %d claim %d: %w", fileName, paymentNumber, claimNumber, err)
					}
					out = append(out, row)
				}
			}
		}
	}
	return out, payments, nil
}

// paymentRow builds one payment row from a transaction's header: the BPR,
// TRN, and DTM*405 segments, the payer and payee loops (the tax ID as the
// match key reads it, from payeeOf), and the counts over its claims and
// PLB segments. The PLB total is the sum of every adjustment amount sent
// (PLB04, PLB06, ... PLB14), 0 with no PLB. A code is written as the text
// the file sent, as the comparison writes one.
func paymentRow(submissionPrefix, fileName string, paymentNumber int, t x12835.TransactionST, p payee) (payment.Payment, error) {
	claims := 0
	for _, lx := range t.HeaderNumbers {
		claims += len(lx.Claims)
	}
	var plbTotal x12835.Decimal
	for _, plb := range t.ProviderAdjustments {
		for _, a := range plb.Adjustments {
			if !a.Amount.Set {
				continue
			}
			sum, err := plbTotal.Add(a.Amount.Value)
			if err != nil {
				return payment.Payment{}, fmt.Errorf("PLB: %w", err)
			}
			plbTotal = sum
		}
	}
	return payment.Payment{
		SubmissionPrefix: submissionPrefix,
		FileName:         fileName,
		PaymentNumber:    paymentNumber,
		PayerIdentifier:  t.PayerN1.IdentifierN104,
		PayerName:        t.PayerN1.NameN102,
		PayeeIdentifier:  t.PayeeN1.IdentifierN104,
		PayeeName:        t.PayeeN1.NameN102,
		PayeeTaxID:       p.taxID,
		PayeeState:       t.PayeeN1.Location.StateN402,
		PaymentDate:      dateValue(t.Payment.EffectiveDateBPR16),
		PaymentMethod:    t.Payment.PayMethodBPR04.String(),
		HandlingCode:     t.Payment.HandlingBPR01.String(),
		TraceNumber:      t.Trace.TraceNumberTRN02,
		ProductionDate:   dateValue(t.ProductionDateDTM405.Date),
		TotalPaid:        amountText(t.Payment.TotalPaidBPR02),
		ClaimCount:       claims,
		PLBCount:         len(t.ProviderAdjustments),
		PLBTotal:         plbTotal.String(),
	}, nil
}

// dateValue is a sent date as a UTC midnight time. A date not sent is the
// zero time, which the payment table writes as NULL.
func dateValue(d x12835.Date) time.Time {
	if !d.Set {
		return time.Time{}
	}
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

// locateFault names the claim a whole-file parse error sits in, by parsing
// each claim's bytes on their own (x12835.ParseClaim reads a claim the way
// Parse does): the first claim that fails on its own is the one, and the
// file's error is returned with its payment and claim number in front. An
// error outside every claim (the envelope, a payment header) is returned
// as it is.
func locateFault(err error, data []byte, d Delimiters, spans []claimSpan) error {
	for _, span := range spans {
		if _, claimErr := x12835.ParseClaim(data[span.start:span.end], d.Element, d.Component, d.Repetition, d.Segment); claimErr != nil {
			return fmt.Errorf("payment %d claim %d: %w", span.payment, span.claim, err)
		}
	}
	return err
}

// claimRow builds one claim row from a typed claim, its payment's payee,
// and its bytes.
func claimRow(submissionPrefix, fileName string, paymentNumber, claimNumber int, c x12835.ClaimCLP, p payee, d Delimiters, data []byte) (claim.Claim, error) {
	date := dateOfService(c)
	if !date.Set {
		return claim.Claim{}, fmt.Errorf("no date of service (DTM*232, or a service line DTM*472 or DTM*150)")
	}
	if !c.ChargeCLP03.Set {
		return claim.Claim{}, fmt.Errorf("CLP03 not sent")
	}
	if !c.PaymentCLP04.Set {
		return claim.Claim{}, fmt.Errorf("CLP04 not sent")
	}
	key := Key{
		MBI:           c.PatientNM1QC.IDNM109,
		DateOfService: dateText(date),
		TotalCharge:   c.ChargeCLP03.Value.String(),
		TaxID:         p.taxID,
		NPI:           p.npi,
		Workload:      p.workload,
	}
	matchKey, err := MatchKey(key)
	if err != nil {
		return claim.Claim{}, err
	}

	segments := make([]byte, len(data))
	copy(segments, data)
	// The jurisdiction and Part are a lookup on the workload; a number the
	// table does not hold leaves both empty.
	workload, _ := macs.Lookup(strings.TrimSpace(key.Workload))
	return claim.Claim{
		SubmissionPrefix:                    submissionPrefix,
		MatchKey:                            matchKey,
		FileName:                            fileName,
		PaymentNumber:                       paymentNumber,
		ClaimNumber:                         claimNumber,
		MedicareBeneficiaryIdentifier:       strings.TrimSpace(key.MBI),
		DateOfService:                       time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC),
		TotalClaimChargeAmount:              key.TotalCharge,
		FederalTaxpayerIdentificationNumber: key.TaxID,
		PayeeIdentificationCode:             key.NPI,
		MACWorkloadID:                       strings.TrimSpace(key.Workload), // the key hashes the trimmed value; compare partitions by this column
		PatientControlNumber:                c.PatientControlCLP01,
		ClaimType:                           claimType(c, p.payerName),
		ClaimStatusCode:                     statusText(c),
		ClaimPaymentAmount:                  c.PaymentCLP04.Value.String(),
		PayerClaimControlNumber:             c.PayerClaimControlCLP07,
		ClaimFrequencyCode:                  c.FrequencyCLP09,
		Adjudication:                        adjudication(c),
		Jurisdiction:                        workload.Jurisdiction,
		Part:                                workload.Part,
		Delimiters:                          d.String(),
		Segments:                            segments,
	}, nil
}

// adjudication says which Medicare adjudication segment the claim carries:
// inpatient for MIA, outpatient for MOA, and nothing for neither. A claim
// with both (which the guide does not allow) reads as inpatient.
func adjudication(c x12835.ClaimCLP) string {
	switch {
	case c.InpatientMIA.Set:
		return "inpatient"
	case c.OutpatientMOA.Set:
		return "outpatient"
	}
	return ""
}

// ParseMatchKey reads a match key as the URL carries it: 64 lowercase hex
// characters.
func ParseMatchKey(s string) ([]byte, error) {
	if len(s) != 64 || strings.ToLower(s) != s {
		return nil, fmt.Errorf("service: parsing match key: want 64 lowercase hex characters: %w", ErrInput)
	}
	k, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("service: parsing match key: want 64 lowercase hex characters: %w", ErrInput)
	}
	return k, nil
}
