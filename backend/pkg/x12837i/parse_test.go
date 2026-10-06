package x12837i

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"
)

const repoRoot = "../../.."

const (
	inpatientFile  = "testdata/x12837i/inpatient-drg.837"
	outpatientFile = "testdata/x12837i/outpatient-revenue-lines.837"
	secondaryFile  = "testdata/x12837i/medicare-secondary.837"
)

// TestParseFiles parses every 837I under testdata/x12837i, and every 837I
// (GS08 005010X223A2) among the taxonomy fixtures' source.837 files, and
// checks that each segment of the input lands in the tree, as a typed
// segment or a RawSegment. An HI segment counts once per composite it
// sends on both sides, since its composites go to separate fields.
func TestParseFiles(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(repoRoot, "testdata/x12837i/*.837"))
	if err != nil {
		t.Fatal(err)
	}
	sources, err := filepath.Glob(filepath.Join(repoRoot, "testdata/taxonomy/*/source.837"))
	if err != nil {
		t.Fatal(err)
	}
	institutional := 0
	for _, path := range sources {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte("*005010X223A2~")) {
			files = append(files, path)
			institutional++
		}
	}
	if len(files) == institutional || institutional == 0 {
		t.Fatalf("found %d files, %d of them taxonomy 837Is, want some of each", len(files), institutional)
	}
	for _, path := range files {
		t.Run(strings.TrimPrefix(path, repoRoot+"/"), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			isa, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse = %v, want nil", err)
			}
			if got, want := countTree(isa), countInput(data); got != want {
				t.Errorf("countTree = %d, want %d (segments in the input)", got, want)
			}
		})
	}
}

// countInput counts the segments in data without the parser: it splits on
// the byte after ISA16, which is byte 105 of the fixed-width ISA, and counts
// the pieces that are not blank. An HI segment counts its non-empty
// elements instead.
func countInput(data []byte) int {
	start := bytes.Index(data, []byte("ISA"))
	element := data[start+3]
	terminator := data[start+105]
	n := 0
	for _, piece := range bytes.Split(data[start:], []byte{terminator}) {
		piece = bytes.TrimSpace(piece)
		if len(piece) == 0 {
			continue
		}
		fields := strings.Split(string(piece), string(element))
		if fields[0] == "HI" {
			n += nonEmpty(fields[1:])
			continue
		}
		n++
	}
	return n
}

func nonEmpty(fields []string) int {
	n := 0
	for _, f := range fields {
		if f != "" {
			n++
		}
	}
	return n
}

// countRaw counts raw segments the way countInput counts input segments.
func countRaw(raw []RawSegment) int {
	n := 0
	for _, r := range raw {
		if r.ID == "HI" {
			n += nonEmpty(r.Elements)
			continue
		}
		n++
	}
	return n
}

func present(b bool) int {
	if b {
		return 1
	}
	return 0
}

// refsFilled counts the REF fields that were placed: a placed REF carries
// its qualifier.
func refsFilled(refs ...ReferenceREF) int {
	n := 0
	for _, r := range refs {
		n += present(r.QualifierREF01 != "")
	}
	return n
}

// amtsFilled counts the AMT fields that were placed: a placed AMT carries
// its qualifier.
func amtsFilled(amts ...AmountAMT) int {
	n := 0
	for _, a := range amts {
		n += present(a.QualifierAMT01 != "")
	}
	return n
}

func named(n NameNM1) int {
	return present(placed(n))
}

// countTree counts the segments the tree holds: each typed segment that was
// present, plus each RawSegment; an HI segment counts its composites.
func countTree(isa InterchangeISA) int {
	n := 2 + countRaw(isa.Raw) // ISA, IEA
	for _, g := range isa.FunctionalGroups {
		n += 2 + countRaw(g.Raw) // GS, GE
		for _, tr := range g.Transactions {
			n += countTransaction(tr)
		}
	}
	return n
}

func countTransaction(tr TransactionST) int {
	n := 2 + countRaw(tr.Raw) + present(tr.Header.Set) // ST, SE
	n += named(tr.SubmitterNM141.Name) + len(tr.SubmitterNM141.Contacts) + countRaw(tr.SubmitterNM141.Raw)
	n += named(tr.ReceiverNM140.Name) + countRaw(tr.ReceiverNM140.Raw)
	for _, b := range tr.BillingProviders {
		n += 1 + countRaw(b.Raw) + present(b.Provider.Set) + present(b.Currency.Set)
		n += named(b.BillingProviderNM185) + present(b.Address.Set) + present(b.Location.Set)
		n += refsFilled(b.TaxIDREFEI) + len(b.References) + len(b.Contacts)
		n += named(b.PayToAddressNM187) + present(b.PayToAddress.Set) + present(b.PayToLocation.Set)
		for _, s := range b.Subscribers {
			n += countSubscriber(s)
		}
	}
	return n
}

func countSubscriber(s SubscriberHL) int {
	n := 1 + countRaw(s.Raw) + present(s.Subscriber.Set)
	n += named(s.SubscriberNM1IL) + present(s.Address.Set) + present(s.Location.Set) + present(s.Demographics.Set)
	n += refsFilled(s.SecondaryIDREFSY, s.ClaimNumberREFY4) + len(s.References)
	n += named(s.PayerNM1PR) + present(s.PayerAddress.Set) + present(s.PayerLocation.Set) + len(s.PayerReferences)
	for _, p := range s.Patients {
		n += 1 + countRaw(p.Raw) + present(p.Patient.Set)
		n += named(p.PatientNM1QC) + present(p.Address.Set) + present(p.Location.Set) + present(p.Demographics.Set) + len(p.References)
		for _, c := range p.Claims {
			n += countClaim(c)
		}
	}
	for _, c := range s.Claims {
		n += countClaim(c)
	}
	return n
}

// countProvider counts an NM1-led provider loop: the NM1 when it was
// placed, its REFs, and its raw segments.
func countProvider(name NameNM1, references []ReferenceREF, raw []RawSegment) int {
	return named(name) + len(references) + countRaw(raw)
}

func countClaim(c ClaimCLM) int {
	n := 1 + countRaw(c.Raw)
	n += present(c.DischargeHourDTP096.Set) + present(c.StatementDatesDTP434.Set) + present(c.AdmissionDTP435.Set) + present(c.RepricerReceivedDTP050.Set)
	n += present(c.InstitutionalCL1.Set) + len(c.SupplementalInfo) + len(c.Contracts)
	n += amtsFilled(c.PatientEstimatedAmountAMTF3) + len(c.Amounts)
	n += refsFilled(c.ServiceAuthExceptionREF4N, c.ReferralNumberREF9F, c.PriorAuthREFG1, c.PayerClaimControlREFF8,
		c.RepricedClaimREF9A, c.AdjustedRepricedClaimREF9C, c.InvestigationalDeviceREFLX, c.TransmissionIntermediaryREFD9,
		c.AutoAccidentStateREFLU, c.MedicalRecordREFEA, c.DemonstrationProjectREFP4, c.PROApprovalREFG4)
	n += len(c.References) + len(c.Files) + len(c.Notes) + len(c.EPSDTReferralsCRCZZ)
	for _, hi := range []HealthInfoHI{c.PrincipalDiagnosisHIABK, c.AdmittingDiagnosisHIABJ, c.PrincipalProcedureHIBBR, c.DRGHIDR} {
		n += present(hi.Set)
	}
	for _, his := range [][]HealthInfoHI{c.PatientReasonsHIAPR, c.ExternalCausesHIABN, c.OtherDiagnosesHIABF, c.OtherProceduresHIBBQ,
		c.OccurrenceSpansHIBI, c.OccurrenceCodesHIBH, c.ValueCodesHIBE, c.ConditionCodesHIBG, c.TreatmentCodesHITC} {
		n += len(his)
	}
	n += len(c.Pricing)
	n += countProvider(c.AttendingNM171.Name, c.AttendingNM171.References, c.AttendingNM171.Raw) + present(c.AttendingNM171.Specialty.Set)
	n += countProvider(c.OperatingNM172.Name, c.OperatingNM172.References, c.OperatingNM172.Raw)
	n += countProvider(c.OtherOperatingNM1ZZ.Name, c.OtherOperatingNM1ZZ.References, c.OtherOperatingNM1ZZ.Raw)
	n += countProvider(c.RenderingNM182.Name, c.RenderingNM182.References, c.RenderingNM182.Raw)
	n += countProvider(c.ServiceFacilityNM177.Name, c.ServiceFacilityNM177.References, c.ServiceFacilityNM177.Raw)
	n += countProvider(c.ReferringNM1DN.Name, c.ReferringNM1DN.References, c.ReferringNM1DN.Raw)
	for _, o := range c.OtherSubscribers {
		n += countOtherSubscriber(o)
	}
	for _, l := range c.ServiceLines {
		n += countServiceLine(l)
	}
	return n
}

func countOtherSubscriber(o OtherSubscriberSBR) int {
	n := 1 + countRaw(o.Raw) + len(o.Adjustments)
	n += present(o.MOA.Set) + present(o.MIA.Set) + present(o.OtherInsurance.Set)
	n += amtsFilled(o.OtherPayerPaidAMTD, o.RemainingLiabilityAMTEAF, o.NonCoveredAMTA8) + len(o.Amounts)
	n += countProvider(o.OtherInsuredNM1IL.Name, o.OtherInsuredNM1IL.References, o.OtherInsuredNM1IL.Raw)
	p := o.OtherPayerNM1PR
	n += countProvider(p.Name, p.References, p.Raw) + present(p.CheckOrRemittanceDTP573.Set)
	n += refsFilled(p.SecondaryIDREF2U, p.SecondaryIDREFEI, p.SecondaryIDREFFY, p.SecondaryIDREFNF,
		p.PriorAuthREFG1, p.ReferralNumberREF9F, p.AdjustmentIndicatorREFT4, p.ClaimControlREFF8)
	n += countProvider(o.OtherPayerAttendingNM171.Name, o.OtherPayerAttendingNM171.References, o.OtherPayerAttendingNM171.Raw)
	n += countProvider(o.OtherPayerOperatingNM172.Name, o.OtherPayerOperatingNM172.References, o.OtherPayerOperatingNM172.Raw)
	n += countProvider(o.OtherPayerOtherOperatingNM1ZZ.Name, o.OtherPayerOtherOperatingNM1ZZ.References, o.OtherPayerOtherOperatingNM1ZZ.Raw)
	n += countProvider(o.OtherPayerServiceFacilityNM177.Name, o.OtherPayerServiceFacilityNM177.References, o.OtherPayerServiceFacilityNM177.Raw)
	n += countProvider(o.OtherPayerRenderingNM182.Name, o.OtherPayerRenderingNM182.References, o.OtherPayerRenderingNM182.Raw)
	n += countProvider(o.OtherPayerReferringNM1DN.Name, o.OtherPayerReferringNM1DN.References, o.OtherPayerReferringNM1DN.Raw)
	n += countProvider(o.OtherPayerBillingNM185.Name, o.OtherPayerBillingNM185.References, o.OtherPayerBillingNM185.Raw)
	return n
}

func countServiceLine(l ServiceLineLX) int {
	n := 1 + countRaw(l.Raw) + present(l.Service.Set) + len(l.SupplementalInfo)
	n += present(l.ServiceDateDTP472.Set) + len(l.Dates)
	n += refsFilled(l.LineControlREF6R, l.RepricedLineREF9B, l.AdjustedRepricedLineREF9D) + len(l.References)
	n += amtsFilled(l.ServiceTaxAMTGT, l.FacilityTaxAMTN8) + len(l.Amounts) + present(l.Pricing.Set)
	for _, d := range l.Drugs {
		n += 1 + countRaw(d.Raw) + present(d.Quantity.Set) + refsFilled(d.LinkSequenceREFVY, d.PrescriptionREFXZ) + len(d.References)
	}
	for _, p := range []LineProviderNM1{l.OperatingNM172, l.OtherOperatingNM1ZZ, l.RenderingNM182, l.ReferringNM1DN} {
		n += countProvider(p.Name, p.References, p.Raw)
	}
	for _, a := range l.Adjudications {
		n += 1 + countRaw(a.Raw) + len(a.Adjustments) + present(a.CheckOrRemittanceDateDTP573.Set) + len(a.Dates)
		n += amtsFilled(a.RemainingLiabilityAMTEAF) + len(a.Amounts)
	}
	return n
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot, path))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func parseText(t *testing.T, text string) InterchangeISA {
	t.Helper()
	isa, err := Parse([]byte(text))
	if err != nil {
		t.Fatalf("Parse = %v, want nil", err)
	}
	return isa
}

// amount builds the Amount the parser builds from text: the decimal and
// the text itself.
func amount(t *testing.T, text string) Amount {
	t.Helper()
	value, err := ParseDecimal([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	return Amount{Value: value, Text: text, Set: true}
}

func date(year, month, day int) Date {
	return Date{Year: year, Month: month, Day: day, Set: true}
}

// cents returns an Amount in hundredths, for the balance checks; every
// amount in the fixtures has at most two decimal places.
func cents(t *testing.T, a Amount) int64 {
	t.Helper()
	if a.Value.Scale > 2 {
		t.Fatalf("amount %+v has more than two decimal places", a)
	}
	c := a.Value.Coefficient
	for s := a.Value.Scale; s < 2; s++ {
		c *= 10
	}
	return c
}

// firstClaim returns the first claim under the first subscriber.
func firstClaim(isa InterchangeISA) ClaimCLM {
	return isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Claims[0]
}

// TestParseInpatient spot-checks typed values in the inpatient stay: the
// delimiters, the envelope, the hierarchy, the claim dates, CL1, the HI
// composites (including the DRG and an occurrence span's date range) with
// their positions, the attending provider's PRV, and the service lines, and
// checks that the line charges add up to CLM02. Each placed segment's
// Index is checked against its position in the file.
func TestParseInpatient(t *testing.T) {
	text := readFile(t, inpatientFile)
	at := func(prefix string) int { return segmentIndex(t, text, prefix) }
	isa := parseText(t, text)
	tr := isa.FunctionalGroups[0].Transactions[0]
	provider := tr.BillingProviders[0]
	subscriber := provider.Subscribers[0]
	claim := subscriber.Claims[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"delimiters and layout", []any{isa.ElementSeparator, isa.SegmentTerminator, isa.SegmentBreak, isa.ComponentSepISA16}, []any{byte('*'), byte('~'), "\n", ":"}},
		{"ISA09", isa.DateISA09, date(2026, 10, 1)},
		{"ISA10", isa.TimeISA10, Time{Hour: 10, Minute: 30, Set: true}},
		{"IEA", []any{isa.IndexIEA, isa.GroupCountIEA01}, []any{46, "1"}},
		{"GS01", isa.FunctionalGroups[0].FunctionalIDGS01, TransactionTypeGS01_ST01HC},
		{"ST01", tr.TransactionSetIDST01, TransactionTypeGS01_ST01837},
		{"BHT06", tr.Header.ClaimOrEncounterBHT06, ClaimOrEncounterBHT06CH},
		{"submitter NM108", tr.SubmitterNM141.Name.IDQualifierNM108, IDQualifierNM10846},
		{"submitter PER04", tr.SubmitterNM141.Contacts[0].CommPER04, "5555550100"},
		{"receiver NM109", tr.ReceiverNM140.Name.IDNM109, "SYNMAC0101"},
		{"HL03", provider.Hierarchy.LevelCodeHL03, HierarchicalLevelHL0320},
		{"billing PRV", provider.Provider, ProviderPRV{Set: true, Index: at("PRV*BI*"), ProviderCodePRV01: ProviderRoleCodePRV01BI, RefQualifierPRV02: "PXC", TaxonomyCodePRV03: "282N00000X"}},
		{"billing NM109", provider.BillingProviderNM185.IDNM109, "1234567893"},
		{"REF*EI", provider.TaxIDREFEI, ReferenceREF{Index: at("REF*EI*"), QualifierREF01: "EI", Value: "990000001"}},
		{"subscriber HL02", subscriber.Hierarchy.ParentIDNumberHL02, "1"},
		{"SBR01", subscriber.Subscriber.PayerResponsibilitySBR01, PayerResponsibilitySBR01P},
		{"SBR09", subscriber.Subscriber.ClaimFilingIndicatorSBR09, ClaimFilingIndicatorSBR09MA},
		{"DMG", subscriber.Demographics, DemographicDMG{Set: true, Index: at("DMG*D8*19450615"), DateFormatDMG01: "D8", BirthDateDMG02: date(1945, 6, 15), SexDMG03: SexDMG03F}},
		{"payer NM103", subscriber.PayerNM1PR.LastNM103, "SYNTHETIC MAC"},
		{"CLM index", claim.Index, at("CLM*SYNINPT-0001")},
		{"CLM02", claim.TotalChargeCLM02, amount(t, "12500.00")},
		{"CLM05-1, a NUBC code kept as text", claim.FacilityTypeCLM05_1, "11"},
		{"CLM05-2", claim.FacilityCodeQualifierCLM05_2, FacilityCodeQualifierCLM05_2A},
		{"CLM05-3, a NUBC code kept as text", claim.ClaimFrequencyCLM05_3, "1"},
		{"DTP*096", claim.DischargeHourDTP096, DateDTP{Set: true, Index: at("DTP*096*"), QualifierDTP01: DateQualifierDTP01096, FormatDTP02: "TM", Time: Time{Hour: 11, Minute: 30, Set: true}}},
		{"DTP*434", claim.StatementDatesDTP434, DateDTP{Set: true, Index: at("DTP*434*"), QualifierDTP01: DateQualifierDTP01434, FormatDTP02: "RD8", Date: date(2026, 9, 20), DateThrough: date(2026, 9, 24)}},
		{"DTP*435", claim.AdmissionDTP435, DateDTP{Set: true, Index: at("DTP*435*"), QualifierDTP01: DateQualifierDTP01435, FormatDTP02: "DT", Date: date(2026, 9, 20), Time: Time{Hour: 8, Minute: 15, Set: true}}},
		{"CL1, NUBC codes kept as text", claim.InstitutionalCL1, InstitutionalCL1{Set: true, Index: at("CL1*1*1*01"), AdmissionTypeCL101: "1", AdmissionSourceCL102: "1", PatientStatusCL103: "01"}},
		{"REF*EA", claim.MedicalRecordREFEA.Value, "SYNMRN0001"},
		{"principal diagnosis", claim.PrincipalDiagnosisHIABK, HealthInfoHI{Set: true, Index: at("HI*ABK:I214"), Position: 1, QualifierHI01_1: DiagnosisQualifierHI01_1ABK, CodeHI01_2: "I214", PresentOnAdmissionHI0n_9: PresentOnAdmissionHI0n_9Y}},
		{"admitting diagnosis", claim.AdmittingDiagnosisHIABJ.CodeHI01_2, "R079"},
		{"other diagnoses POA", []PresentOnAdmissionHI0n_9{claim.OtherDiagnosesHIABF[0].PresentOnAdmissionHI0n_9, claim.OtherDiagnosesHIABF[1].PresentOnAdmissionHI0n_9}, []PresentOnAdmissionHI0n_9{PresentOnAdmissionHI0n_9Y, PresentOnAdmissionHI0n_9N}},
		{"other diagnoses positions in one HI", []int{claim.OtherDiagnosesHIABF[0].Index, claim.OtherDiagnosesHIABF[0].Position, claim.OtherDiagnosesHIABF[1].Index, claim.OtherDiagnosesHIABF[1].Position}, []int{at("HI*ABF:E119"), 1, at("HI*ABF:E119"), 2}},
		{"principal procedure", claim.PrincipalProcedureHIBBR, HealthInfoHI{Set: true, Index: at("HI*BBR:"), Position: 1, QualifierHI01_1: DiagnosisQualifierHI01_1BBR, CodeHI01_2: "02703DZ", DateQualifierHI01_3: "D8", DateHI01_4: date(2026, 9, 21)}},
		{"DRG", claim.DRGHIDR.CodeHI01_2, "280"},
		{"occurrence span", claim.OccurrenceSpansHIBI, []HealthInfoHI{{Set: true, Index: at("HI*BI:74"), Position: 1, QualifierHI01_1: DiagnosisQualifierHI01_1BI, CodeHI01_2: "74", DateQualifierHI01_3: "RD8", DateHI01_4: date(2026, 9, 21), DateThroughHI01_4: date(2026, 9, 22)}}},
		{"value code amount", claim.ValueCodesHIBE[0].AmountHI01_5, amount(t, "4")},
		{"condition code", claim.ConditionCodesHIBG[0].CodeHI01_2, "39"},
		{"attending PRV01", claim.AttendingNM171.Specialty.ProviderCodePRV01, ProviderRoleCodePRV01AT},
		{"attending REF", claim.AttendingNM171.References, []ReferenceREF{{Index: at("REF*G2*SYNATT01"), QualifierREF01: "G2", Value: "SYNATT01"}}},
		{"operating NM103", claim.OperatingNM172.Name.LastNM103, "SYNTHETIC"},
		{"service lines", len(claim.ServiceLines), 3},
		{"first line LX01", claim.ServiceLines[0].NumberLX01, "1"},
		{"first line SV204, SV205", []any{claim.ServiceLines[0].Service.UnitSV204, claim.ServiceLines[0].Service.UnitsSV205}, []any{UnitOrBasisOfMeasurementDA, amount(t, "4")}},
		{"third line SV202", claim.ServiceLines[2].Service.ProcedureSV202, ProcedureSV202{Qualifier: ProductServiceQualifierHC, Code: "93458"}},
		{"SE as sent", []any{tr.IndexSE, tr.SegmentCountSE01, tr.ControlNumberSE02}, []any{at("SE*43*0001"), "43", "0001"}},
		{"transaction Raw", len(tr.Raw), 0},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
	var lines int64
	for _, l := range claim.ServiceLines {
		lines += cents(t, l.Service.ChargeSV203)
	}
	if want := cents(t, claim.TotalChargeCLM02); lines != want {
		t.Errorf("sum of SV203 = %d cents, want CLM02 = %d cents", lines, want)
	}
}

// TestParseOutpatient spot-checks the outpatient visit: the pay-to address
// loop, a PER with an extension, the patient's reasons for visit (an HI
// that repeats its composite), an occurrence code with a date, and the
// revenue lines' procedure, service date, and line control number, and
// checks that the line charges add up to CLM02.
func TestParseOutpatient(t *testing.T) {
	text := readFile(t, outpatientFile)
	at := func(prefix string) int { return segmentIndex(t, text, prefix) }
	isa := parseText(t, text)
	provider := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0]
	claim := firstClaim(isa)
	first := claim.ServiceLines[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"billing N302", provider.Address.Line2N302, "SUITE 4"},
		{"billing PER05, PER06", []string{provider.Contacts[0].QualifierPER05, provider.Contacts[0].CommPER06}, []string{"EX", "12"}},
		{"pay-to NM102", provider.PayToAddressNM187.EntityTypeNM102, EntityTypeNM1022},
		{"pay-to N301", provider.PayToAddress.Line1N301, "PO BOX 900"},
		{"pay-to N403", provider.PayToLocation.ZipN403, "220310900"},
		{"CLM02", claim.TotalChargeCLM02, amount(t, "412.50")},
		{"REF*G1", claim.PriorAuthREFG1.Value, "SYNPA0002"},
		{"CL1 with only CL103", claim.InstitutionalCL1, InstitutionalCL1{Set: true, Index: at("CL1***01"), PatientStatusCL103: "01"}},
		{"reasons for visit", []string{claim.PatientReasonsHIAPR[0].CodeHI01_2, claim.PatientReasonsHIAPR[1].CodeHI01_2}, []string{"M545", "R262"}},
		{"reasons for visit positions", []int{claim.PatientReasonsHIAPR[0].Position, claim.PatientReasonsHIAPR[1].Position}, []int{1, 2}},
		{"occurrence code date", claim.OccurrenceCodesHIBH[0].DateHI01_4, date(2026, 9, 14)},
		{"rendering NM103", claim.RenderingNM182.Name.LastNM103, "SYNTHETIC"},
		{"first line SV202", first.Service.ProcedureSV202, ProcedureSV202{Qualifier: ProductServiceQualifierHC, Code: "99284", Modifiers: [4]string{"25"}}},
		{"first line DTP*472", first.ServiceDateDTP472, DateDTP{Set: true, Index: at("DTP*472*D8*20260915"), QualifierDTP01: DateQualifierDTP01472, FormatDTP02: "D8", Date: date(2026, 9, 15)}},
		{"first line REF*6R", first.LineControlREF6R.Value, "SYNLINE0001"},
		{"second line SV203", claim.ServiceLines[1].Service.ChargeSV203, amount(t, "87.50")},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
	var lines int64
	for _, l := range claim.ServiceLines {
		lines += cents(t, l.Service.ChargeSV203)
	}
	if want := cents(t, claim.TotalChargeCLM02); lines != want {
		t.Errorf("sum of SV203 = %d cents, want CLM02 = %d cents", lines, want)
	}
}

// TestParseSecondary spot-checks the Medicare-secondary claim: Medicare's
// SBR, the other payer's 2320/2330 loops, and each line's 2430
// adjudication. It checks the coordination-of-benefits balance: each
// line's charge equals the other payer's line payment plus its line
// adjustments, the line payments add up to AMT*D, and CLM02 less every
// adjustment equals AMT*D.
func TestParseSecondary(t *testing.T) {
	text := readFile(t, secondaryFile)
	at := func(prefix string) int { return segmentIndex(t, text, prefix) }
	isa := parseText(t, text)
	subscriber := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0]
	claim := subscriber.Claims[0]
	other := claim.OtherSubscribers[0]
	adjudication := claim.ServiceLines[0].Adjudications[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"SBR01", subscriber.Subscriber.PayerResponsibilitySBR01, PayerResponsibilitySBR01S},
		{"SBR05", subscriber.Subscriber.InsuranceTypeSBR05, "12"},
		{"external cause", claim.ExternalCausesHIABN[0].CodeHI01_2, "W010XXA"},
		{"other subscribers", len(claim.OtherSubscribers), 1},
		{"other SBR", other.SBR, SubscriberSBR{Set: true, Index: at("SBR*P*18*SYNGRP77"), PayerResponsibilitySBR01: PayerResponsibilitySBR01P, RelationshipSBR02: RelationshipSBR0218, GroupNumberSBR03: "SYNGRP77", GroupNameSBR04: "SYNTHETIC EMPLOYER PLAN", ClaimFilingIndicatorSBR09: ClaimFilingIndicatorSBR09CI}},
		{"AMT*D", other.OtherPayerPaidAMTD, AmountAMT{Index: at("AMT*D*450.00"), QualifierAMT01: "D", Amount: amount(t, "450.00")}},
		{"AMT*EAF", other.RemainingLiabilityAMTEAF.Amount, amount(t, "60.00")},
		{"OI03, OI06", []OtherInsuranceIndicator{other.OtherInsurance.BenefitsAssignOI03, other.OtherInsurance.ReleaseOfInfoOI06}, []OtherInsuranceIndicator{OtherInsuranceIndicatorY, OtherInsuranceIndicatorY}},
		{"MOA03", other.MOA.RemarkMOA03, rarc.Code("MA01")},
		{"2330A NM109", other.OtherInsuredNM1IL.Name.IDNM109, "SYNEMP7700"},
		{"2330B NM109", other.OtherPayerNM1PR.Name.IDNM109, "SYNPAYER77"},
		{"2330B DTP*573", other.OtherPayerNM1PR.CheckOrRemittanceDTP573, DateDTP{Set: true, Index: at("DTP*573*D8*20260920"), QualifierDTP01: DateQualifierDTP01573, FormatDTP02: "D8", Date: date(2026, 9, 20)}},
		{"2330B REF*F8", other.OtherPayerNM1PR.ClaimControlREFF8.Value, "SYNCLAIM7700"},
		{"SVD03", adjudication.ProcedureSVD03, ProcedureSV202{Qualifier: ProductServiceQualifierHC, Code: "99284"}},
		{"SVD05", adjudication.PaidUnitsSVD05, amount(t, "1")},
		{"line CAS", adjudication.Adjustments, []AdjustmentCAS{
			{Index: at("CAS*CO*45*60.00"), GroupCAS01: CASGroupCode01CO, Reasons: []CASReason{{Reason: carc.Code("45"), Amount: amount(t, "60.00")}}},
			{Index: at("CAS*PR*2*40.00"), GroupCAS01: CASGroupCode01PR, Reasons: []CASReason{{Reason: carc.Code("2"), Amount: amount(t, "40.00")}}},
		}},
		{"2430 DTP*573", adjudication.CheckOrRemittanceDateDTP573.Date, date(2026, 9, 20)},
		{"2430 AMT*EAF", adjudication.RemainingLiabilityAMTEAF.Amount, amount(t, "40.00")},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
	var paid, adjusted int64
	for _, l := range claim.ServiceLines {
		a := l.Adjudications[0]
		var lineAdjusted int64
		for _, cas := range a.Adjustments {
			for _, r := range cas.Reasons {
				lineAdjusted += cents(t, r.Amount)
			}
		}
		if got, want := cents(t, a.PaymentSVD02)+lineAdjusted, cents(t, l.Service.ChargeSV203); got != want {
			t.Errorf("line %s: SVD02 + CAS = %d cents, want SV203 = %d cents", l.NumberLX01, got, want)
		}
		paid += cents(t, a.PaymentSVD02)
		adjusted += lineAdjusted
	}
	if want := cents(t, other.OtherPayerPaidAMTD.Amount); paid != want {
		t.Errorf("sum of SVD02 = %d cents, want AMT*D = %d cents", paid, want)
	}
	if got, want := cents(t, claim.TotalChargeCLM02)-adjusted, cents(t, other.OtherPayerPaidAMTD.Amount); got != want {
		t.Errorf("CLM02 - CAS = %d cents, want AMT*D = %d cents", got, want)
	}
}

// edit replaces old with new once in text, adds added to the SE01 count
// (the segments the edit adds inside the transaction), and fails when text
// does not contain old.
func edit(t *testing.T, text, old, new string, added int) string {
	t.Helper()
	if !strings.Contains(text, old) {
		t.Fatalf("text does not contain %q", old)
	}
	text = strings.Replace(text, old, new, 1)
	if added == 0 {
		return text
	}
	start := strings.Index(text, "\nSE*") + len("\nSE*")
	end := start + strings.IndexByte(text[start:], '*')
	count, err := strconv.Atoi(text[start:end])
	if err != nil {
		t.Fatal(err)
	}
	return text[:start] + strconv.Itoa(count+added) + text[end:]
}

// segmentIndex returns the 0-based file position of the segment that
// starts with prefix: the number of segment terminators before it.
func segmentIndex(t *testing.T, text, prefix string) int {
	t.Helper()
	i := strings.Index(text, "\n"+prefix)
	if i < 0 {
		t.Fatalf("text has no segment starting %q", prefix)
	}
	return strings.Count(text[:i], "~")
}

// roundTrips writes a parsed tree back and fails unless the bytes are the
// text it was parsed from.
func roundTrips(t *testing.T, text string, isa InterchangeISA) {
	t.Helper()
	var out bytes.Buffer
	if err := Write(&out, isa); err != nil {
		t.Fatalf("Write = %v, want nil", err)
	}
	if out.String() != text {
		t.Error(firstDifference([]byte(text), out.Bytes()))
	}
}

// TestParseHierarchy checks HL placement by HL03 and HL02: a patient loop
// (2000C) under its subscriber, a subscriber placed under an earlier billing
// provider than the latest one, and an HL with no place kept raw with its
// loop.
func TestParseHierarchy(t *testing.T) {
	base := readFile(t, outpatientFile)

	t.Run("patient loop", func(t *testing.T) {
		text := edit(t, base, "HL*2*1*22*0~", "HL*2*1*22*1~", 0)
		text = edit(t, text, "CLM*SYNOUTP-0002*", "HL*3*2*23*0~\nPAT*01~\nNM1*QC*1*SYNTHETIC*DEPENDENT~\nDMG*D8*19800101*U~\nCLM*SYNOUTP-0002*", 4)
		isa := parseText(t, text)
		subscriber := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0]
		if len(subscriber.Claims) != 0 || len(subscriber.Patients) != 1 {
			t.Fatalf("subscriber claims, patients = %d, %d, want 0, 1", len(subscriber.Claims), len(subscriber.Patients))
		}
		patient := subscriber.Patients[0]
		tests := []struct {
			name string
			got  any
			want any
		}{
			{"HL02", patient.Hierarchy.ParentIDNumberHL02, "2"},
			{"PAT01", patient.Patient, PatientRelationshipPAT{Set: true, Index: segmentIndex(t, text, "PAT*01"), RelationshipPAT01: RelationshipSBR0201}},
			{"NM104", patient.PatientNM1QC.FirstNM104, "DEPENDENT"},
			{"DMG03", patient.Demographics.SexDMG03, SexDMG03U},
			{"claims", len(patient.Claims), 1},
		}
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
			}
		}
		if got, want := countTree(isa), countInput([]byte(text)); got != want {
			t.Errorf("countTree = %d, want %d", got, want)
		}
		roundTrips(t, text, isa)
	})

	t.Run("subscriber under an earlier billing provider", func(t *testing.T) {
		text := edit(t, base, "\nSE*", "\nHL*3**20*1~\nNM1*85*2*SYNTHETIC SECOND HOSPITAL*****XX*1234567893~\nHL*4*1*22*0~\nSBR*P*18*******MA~\nSE*", 4)
		isa := parseText(t, text)
		providers := isa.FunctionalGroups[0].Transactions[0].BillingProviders
		if len(providers) != 2 || len(providers[0].Subscribers) != 2 || len(providers[1].Subscribers) != 0 {
			t.Fatalf("providers = %d with %d and %d subscribers, want 2 with 2 and 0", len(providers), len(providers[0].Subscribers), len(providers[1].Subscribers))
		}
		if got := providers[0].Subscribers[1].Hierarchy.IDNumberHL01; got != "4" {
			t.Errorf("second subscriber HL01 = %q, want \"4\"", got)
		}
		if got, want := countTree(isa), countInput([]byte(text)); got != want {
			t.Errorf("countTree = %d, want %d", got, want)
		}
		roundTrips(t, text, isa)
	})

	sbr := RawSegment{ID: "SBR", Elements: []string{"P", "18", "", "", "", "", "", "", "MA"}}
	tests := []struct {
		name   string
		insert string
		want   []RawSegment // Index is set from the file
	}{
		{"parent not found", "HL*3*9*22*0~\nSBR*P*18*******MA~\n", []RawSegment{{ID: "HL", Elements: []string{"3", "9", "22", "0"}}, sbr}},
		{"level off the list", "HL*3*1*99*0~\nSBR*P*18*******MA~\n", []RawSegment{{ID: "HL", Elements: []string{"3", "1", "99", "0"}}, sbr}},
		{"billing provider with a parent", "HL*3*1*20*0~\nPRV*BI*PXC*282N00000X~\n", []RawSegment{{ID: "HL", Elements: []string{"3", "1", "20", "0"}}, {ID: "PRV", Elements: []string{"BI", "PXC", "282N00000X"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := edit(t, base, "\nSE*", "\n"+tt.insert+"SE*", 2)
			isa := parseText(t, text)
			tr := isa.FunctionalGroups[0].Transactions[0]
			first := segmentIndex(t, text, "HL*3*")
			want := make([]RawSegment, len(tt.want))
			for i, r := range tt.want {
				r.Index = first + i
				want[i] = r
			}
			if !reflect.DeepEqual(tr.Raw, want) {
				t.Errorf("transaction Raw = %+v, want %+v", tr.Raw, want)
			}
			if got, want := countTree(isa), countInput([]byte(text)); got != want {
				t.Errorf("countTree = %d, want %d", got, want)
			}
			roundTrips(t, text, isa)
		})
	}
}

// TestParseOffList changes one coded value in a fixture to a code off its
// list and checks that the field holds the text as sent, that Known reports
// false, and that the file writes back unchanged.
func TestParseOffList(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		old, new string
		get      func(InterchangeISA) (string, bool)
		want     string
	}{
		{
			name: "BHT06", path: inpatientFile, old: "*1030*CH~", new: "*1030*ZZ~",
			get: func(isa InterchangeISA) (string, bool) {
				v := isa.FunctionalGroups[0].Transactions[0].Header.ClaimOrEncounterBHT06
				return v.String(), v.Known()
			},
			want: "ZZ",
		},
		{
			name: "PRV01", path: inpatientFile, old: "PRV*BI*", new: "PRV*XX*",
			get: func(isa InterchangeISA) (string, bool) {
				v := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Provider.ProviderCodePRV01
				return v.String(), v.Known()
			},
			want: "XX",
		},
		{
			name: "SBR09", path: inpatientFile, old: "*******MA~", new: "*******XX~",
			get: func(isa InterchangeISA) (string, bool) {
				v := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Subscriber.ClaimFilingIndicatorSBR09
				return v.String(), v.Known()
			},
			want: "XX",
		},
		{
			name: "HI01-9", path: inpatientFile, old: "HI*ABK:I214:::::::Y~", new: "HI*ABK:I214:::::::X~",
			get: func(isa InterchangeISA) (string, bool) {
				v := firstClaim(isa).PrincipalDiagnosisHIABK.PresentOnAdmissionHI0n_9
				return v.String(), v.Known()
			},
			want: "X",
		},
		{
			name: "NM108 in the attending provider loop", path: inpatientFile, old: "ATTENDING****XX*", new: "ATTENDING****ZZ*",
			get: func(isa InterchangeISA) (string, bool) {
				v := firstClaim(isa).AttendingNM171.Name.IDQualifierNM108
				return v.String(), v.Known()
			},
			want: "ZZ",
		},
		{
			name: "SV202-1", path: inpatientFile, old: "SV2*0481*HC:", new: "SV2*0481*ZZ:",
			get: func(isa InterchangeISA) (string, bool) {
				p := firstClaim(isa).ServiceLines[2].Service.ProcedureSV202
				return p.Qualifier.String() + ":" + p.Code, p.Qualifier.Known()
			},
			want: "ZZ:93458",
		},
		{
			name: "OI03", path: secondaryFile, old: "OI***Y***Y~", new: "OI***Q***Y~",
			get: func(isa InterchangeISA) (string, bool) {
				v := firstClaim(isa).OtherSubscribers[0].OtherInsurance.BenefitsAssignOI03
				return v.String(), v.Known()
			},
			want: "Q",
		},
		{
			name: "CAS01", path: secondaryFile, old: "CAS*CO*45*60.00~", new: "CAS*XX*45*60.00~",
			get: func(isa InterchangeISA) (string, bool) {
				v := firstClaim(isa).ServiceLines[0].Adjudications[0].Adjustments[0].GroupCAS01
				return v.String(), v.Known()
			},
			want: "XX",
		},
		{
			name: "DTP01 on a service line goes to Dates", path: outpatientFile, old: "REF*6R*SYNLINE0001~", new: "REF*6R*SYNLINE0001~\nDTP*999*RD8*20260901-20260902~",
			get: func(isa InterchangeISA) (string, bool) {
				d := firstClaim(isa).ServiceLines[0].Dates[0]
				return d.QualifierDTP01.String() + " " + dateOut(d.Date) + "-" + dateOut(d.DateThrough), d.QualifierDTP01.Known()
			},
			want: "999 20260901-20260902",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := strings.Count(tt.new, "~") - strings.Count(tt.old, "~")
			text := edit(t, readFile(t, tt.path), tt.old, tt.new, delta)
			isa := parseText(t, text)
			got, known := tt.get(isa)
			if got != tt.want || known {
				t.Errorf("value, known = %q, %v, want %q, false", got, known, tt.want)
			}
			roundTrips(t, text, isa)
		})
	}
}

// TestParseExtraElements sends text in positions X12 defines but the model
// does not name, and checks that each lands in the struct's ExtraElements
// (or the composite's own field) and writes back unchanged.
func TestParseExtraElements(t *testing.T) {
	tests := []struct {
		name     string
		old, new string
		added    int
		get      func(InterchangeISA) any
		want     any
	}{
		{
			name: "NM110 and NM111 on the billing provider",
			old:  "NM1*85*2*SYNTHETIC GENERAL HOSPITAL*****XX*1234567893~", new: "NM1*85*2*SYNTHETIC GENERAL HOSPITAL*****XX*1234567893*A*B~",
			get: func(isa InterchangeISA) any {
				return isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].BillingProviderNM185.ExtraElements
			},
			want: []string{"A", "B"},
		},
		{
			name: "REF03 on a claim REF",
			old:  "REF*EA*SYNMRN0001~", new: "REF*EA*SYNMRN0001*DESC~",
			get:  func(isa InterchangeISA) any { return firstClaim(isa).MedicalRecordREFEA.ExtraElements },
			want: []string{"DESC"},
		},
		{
			name: "PER09 after two empty positions",
			old:  "PER*IC*SYNTHETIC BILLING OFFICE*TE*5555550100*EM*BILLING@EXAMPLE.TEST~", new: "PER*IC*SYNTHETIC BILLING OFFICE*TE*5555550100*EM*BILLING@EXAMPLE.TEST***REF9~",
			get: func(isa InterchangeISA) any {
				p := isa.FunctionalGroups[0].Transactions[0].SubmitterNM141.Contacts[0]
				return []any{p.QualifierPER07, p.CommPER08, p.ExtraElements}
			},
			want: []any{"", "", []string{"REF9"}},
		},
		{
			name: "the eighth component of SV202",
			old:  "SV2*0481*HC:93458*", new: "SV2*0481*HC:93458::::::ID8*",
			get: func(isa InterchangeISA) any {
				return firstClaim(isa).ServiceLines[2].Service.ProcedureSV202.ProductServiceID
			},
			want: "ID8",
		},
		{
			name: "CTP05's other components and CTP06",
			old:  "SV2*0481*HC:93458*3000.00*UN*1~", new: "SV2*0481*HC:93458*3000.00*UN*1~\nLIN**N4*63323026201~\nCTP****1*UN:2:3*X6~",
			added: 2,
			get: func(isa InterchangeISA) any {
				q := firstClaim(isa).ServiceLines[2].Drugs[0].Quantity
				return []any{q.UnitCTP05_1, q.UnitCTP05Rest, q.ExtraElements}
			},
			want: []any{UnitOrBasisOfMeasurementUN, []string{"2", "3"}, []string{"X6"}},
		},
		{
			name: "DMG04",
			old:  "DMG*D8*19450615*F~", new: "DMG*D8*19450615*F*M~",
			get: func(isa InterchangeISA) any {
				return isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Demographics.ExtraElements
			},
			want: []string{"M"},
		},
	}
	base := readFile(t, inpatientFile)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := edit(t, base, tt.old, tt.new, tt.added)
			isa := parseText(t, text)
			if got := tt.get(isa); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
			if got, want := countTree(isa), countInput([]byte(text)); got != want {
				t.Errorf("countTree = %d, want %d", got, want)
			}
			roundTrips(t, text, isa)
		})
	}
}

// TestParseRaw inserts segments the model has no slot for into the
// inpatient fixture and checks that each is kept as a RawSegment on the
// struct that encloses it, with its file position, that nothing is
// dropped, and that the file writes back unchanged.
func TestParseRaw(t *testing.T) {
	claimRaw := func(isa InterchangeISA) []RawSegment { return firstClaim(isa).Raw }
	tests := []struct {
		name     string
		old, new string
		added    int    // segments the edit adds inside the transaction, for SE01
		at       string // the start of the first raw segment, to find its file position
		get      func(InterchangeISA) []RawSegment
		want     []RawSegment // Index is set from the file
	}{
		{
			name: "unmodeled segment in a claim", old: "REF*EA*SYNMRN0001~", new: "REF*EA*SYNMRN0001~\nQTY*CA*1~", added: 1, at: "QTY*",
			get: claimRaw, want: []RawSegment{{ID: "QTY", Elements: []string{"CA", "1"}}},
		},
		{
			name: "second CL1", old: "CL1*1*1*01~", new: "CL1*1*1*01~\nCL1*2*1*01~", added: 1, at: "CL1*2*",
			get: claimRaw, want: []RawSegment{{ID: "CL1", Elements: []string{"2", "1", "01"}}},
		},
		{
			name: "claim DTP qualifier with no claim field", old: "CL1*1*1*01~", new: "DTP*472*D8*20260920~\nCL1*1*1*01~", added: 1, at: "DTP*472*",
			get: claimRaw, want: []RawSegment{{ID: "DTP", Elements: []string{"472", "D8", "20260920"}}},
		},
		{
			name: "claim DTP in a format its field cannot hold", old: "CL1*1*1*01~", new: "DTP*050*RD8*20260901-20260902~\nCL1*1*1*01~", added: 1, at: "DTP*050*",
			get: claimRaw, want: []RawSegment{{ID: "DTP", Elements: []string{"050", "RD8", "20260901-20260902"}}},
		},
		{
			name: "second principal diagnosis", old: "HI*ABJ:R079~", new: "HI*ABJ:R079~\nHI*ABK:I219~", added: 1, at: "HI*ABK:I219",
			get: claimRaw, want: []RawSegment{{ID: "HI", Elements: []string{"ABK:I219"}}},
		},
		{
			name: "HI with an off-list qualifier", old: "HI*ABJ:R079~", new: "HI*ABJ:R079~\nHI*ABF:E785*ZZZ:1~", added: 1, at: "HI*ABF:E785",
			get: claimRaw, want: []RawSegment{{ID: "HI", Elements: []string{"ABF:E785", "ZZZ:1"}}},
		},
		{
			name: "pay-to plan loop", old: "HL*2*1*22*0~", new: "NM1*PE*2*SYNTHETIC PLAN*****PI*SYNPLAN~\nN3*1 SYNTHETIC PLAN WAY~\nHL*2*1*22*0~", added: 2, at: "NM1*PE*",
			get: func(isa InterchangeISA) []RawSegment {
				return isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Raw
			},
			want: []RawSegment{
				{ID: "NM1", Elements: []string{"PE", "2", "SYNTHETIC PLAN", "", "", "", "", "PI", "SYNPLAN"}},
				{ID: "N3", Elements: []string{"1 SYNTHETIC PLAN WAY"}},
			},
		},
		{
			name: "line-level NTE", old: "SV2*0120**8000.00*DA*4~", new: "SV2*0120**8000.00*DA*4~\nNTE*TPO*SYNTHETIC NOTE~", added: 1, at: "NTE*",
			get:  func(isa InterchangeISA) []RawSegment { return firstClaim(isa).ServiceLines[0].Raw },
			want: []RawSegment{{ID: "NTE", Elements: []string{"TPO", "SYNTHETIC NOTE"}}},
		},
		{
			name: "second attending provider loop", old: "NM1*72*", new: "NM1*71*1*SYNTHETIC*SECOND~\nREF*G2*SYNATT02~\nNM1*72*", added: 2, at: "NM1*71*1*SYNTHETIC*SECOND",
			get: claimRaw,
			want: []RawSegment{
				{ID: "NM1", Elements: []string{"71", "1", "SYNTHETIC", "SECOND"}},
				{ID: "REF", Elements: []string{"G2", "SYNATT02"}},
			},
		},
		{
			name: "unmodeled segment in a group", old: "ST*837*0001", new: "XYZ*1~\nST*837*0001", at: "XYZ*",
			get:  func(isa InterchangeISA) []RawSegment { return isa.FunctionalGroups[0].Raw },
			want: []RawSegment{{ID: "XYZ", Elements: []string{"1"}}},
		},
	}
	base := readFile(t, inpatientFile)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := edit(t, base, tt.old, tt.new, tt.added)
			first := segmentIndex(t, text, tt.at)
			want := make([]RawSegment, len(tt.want))
			for i, r := range tt.want {
				r.Index = first + i
				want[i] = r
			}
			isa := parseText(t, text)
			if got := tt.get(isa); !reflect.DeepEqual(got, want) {
				t.Errorf("Raw = %+v, want %+v", got, want)
			}
			if got, want := countTree(isa), countInput([]byte(text)); got != want {
				t.Errorf("countTree = %d, want %d", got, want)
			}
			roundTrips(t, text, isa)
		})
	}
}

// TestParseErrors checks that Parse rejects malformed input and that each
// error names the element and the segment, and that each envelope check
// fails on a trailer that does not match.
func TestParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		old, new string
		want     string
	}{
		{"SE01", "SE*43*0001", "SE*42*0001", `x12837i: checking SE01 of transaction "0001": SE01 is 42, ST through SE is 43 segments`},
		{"SE02", "SE*43*0001", "SE*43*0002", `x12837i: checking SE02: ST02 is "0001", SE02 is "0002"`},
		{"GE01", "GE*1*101~", "GE*2*101~", `x12837i: checking GE01 of group "101": GE01 is 2, the group holds 1 transactions`},
		{"GE02", "GE*1*101~", "GE*1*102~", `x12837i: checking GE02: GS06 is "101", GE02 is "102"`},
		{"IEA01", "IEA*1*", "IEA*2*", `x12837i: checking IEA01: IEA01 is 2, the interchange holds 1 functional groups`},
		{"IEA02", "IEA*1*000000101", "IEA*1*000000102", `x12837i: checking IEA02: ISA13 is "000000101", IEA02 is "000000102"`},
		{"malformed amount", "*12500.00***", "*12500.0O***", "x12837i: parsing CLM02 in segment 20: "},
		{"date range without a dash", "DTP*434*RD8*20260920-20260924", "DTP*434*RD8*20260920", "x12837i: parsing DTP03 in segment 22: date range is not CCYYMMDD-CCYYMMDD"},
		{"admission hour not a clock time", "DTP*435*DT*202609200815", "DTP*435*DT*202609202475", "x12837i: parsing DTP03 in segment 23: not a clock time"},
		{"HI date not a calendar date", "HI*BBR:02703DZ:D8:20260921", "HI*BBR:02703DZ:D8:20260931", "x12837i: parsing HI01-4 in segment 29: not a calendar date"},
		{"component past the composite", "SV2*0481*HC:93458*", "SV2*0481*HC:93458:::::::X*", "x12837i: parsing SV202 in segment 43: the model ends at component 8"},
		{"element past the segment", "LX*1~", "LX*1*2~", "x12837i: parsing LX02 in segment 38: the model ends at LX01"},
		{"element past X12's own last position", "N4*TESTVILLE*MD*212440001~", "N4*TESTVILLE*MD*212440001*****X~", "x12837i: parsing N408 in segment 11: the model ends at N407"},
		{"LX01 not a number", "LX*1~", "LX*A~", "x12837i: parsing LX01 in segment 38: not a whole number"},
		{"SE01 not a number", "SE*43*0001", "SE*4A*0001", "x12837i: parsing SE01 in segment 44: not a whole number"},
		{"missing SE", "SE*43*0001~", "", "x12837i: parsing segment 44: GE inside the transaction at segment 2 (missing SE)"},
		{"input after IEA", "IEA*1*000000101~", "IEA*1*000000101~TA1*1~", "x12837i: parsing segment 47: input continues after the IEA"},
		{"missing terminator", "IEA*1*000000101~", "IEA*1*000000101", "x12837i: splitting segment 46: missing segment terminator"},
	}
	base := readFile(t, inpatientFile)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("%s does not contain %q", inpatientFile, tt.old)
			}
			_, err := Parse([]byte(strings.Replace(base, tt.old, tt.new, 1)))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Parse error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

// TestParseNotUsedPositionsKept checks that text in a Medicare "Not Used"
// position never fails Parse and lands in that position's typed field: the
// parser models every X12-defined position, and Medicare's usage rules are
// a caller concern (see Parse's doc comment).
func TestParseNotUsedPositionsKept(t *testing.T) {
	base := readFile(t, inpatientFile)
	text := strings.Replace(base, "HI*ABJ:R079~", "HI*ABJ:R079:::::V:W~", 1)
	isa := parseText(t, text)
	claim := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Claims[0]
	hi := claim.AdmittingDiagnosisHIABJ
	if hi.VersionIdentifierHI0n_7 != "V" || hi.IndustryCodeHI0n_8 != "W" {
		t.Errorf("HI0n-7/8 = %q/%q, want \"V\"/\"W\"", hi.VersionIdentifierHI0n_7, hi.IndustryCodeHI0n_8)
	}
	roundTrips(t, text, isa)
}

// TestParseLineEndings checks that whitespace between segments, and a
// carriage return before a newline terminator, do not change the tree
// beyond the layout the interchange records (SegmentBreak), and that each
// layout writes back as it came.
func TestParseLineEndings(t *testing.T) {
	tests := []struct {
		name     string
		old, new string
		break_   string
	}{
		{"tilde then CRLF", "~\n", "~\r\n", "\r\n"},
		{"no line breaks", "~\n", "~", ""},
	}
	text := readFile(t, secondaryFile)
	want := parseText(t, text)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed := strings.ReplaceAll(text, tt.old, tt.new)
			got := parseText(t, changed)
			if got.SegmentBreak != tt.break_ {
				t.Errorf("SegmentBreak = %q, want %q", got.SegmentBreak, tt.break_)
			}
			roundTrips(t, changed, got)
			got.SegmentBreak = want.SegmentBreak
			if !reflect.DeepEqual(got, want) {
				t.Errorf("tree after replacing %q with %q differs from the original's", tt.old, tt.new)
			}
		})
	}
}

// TestParseElements checks the date, date range, date-and-hour, and amount
// readers on single elements.
func TestParseElements(t *testing.T) {
	segment := func(texts ...string) inputSegment {
		return inputSegment{index: 4, id: "DTP", elements: texts}
	}
	tests := []struct {
		name    string
		parse   func() (any, error)
		want    any
		wantErr string
	}{
		{"D8", func() (any, error) { d, _, err := parseDTPPeriod(segment("472", "D8", "20240229")); return d, err }, date(2024, 2, 29), ""},
		{"RD8", func() (any, error) {
			from, to, err := parseDTPPeriod(segment("434", "RD8", "20260101-20260131"))
			return []Date{from, to}, err
		}, []Date{date(2026, 1, 1), date(2026, 1, 31)}, ""},
		{"RD8 with a bad second date", func() (any, error) {
			_, _, err := parseDTPPeriod(segment("434", "RD8", "20260101-2026013"))
			return nil, err
		}, nil, "x12837i: parsing DTP03 in segment 4: date is not CCYYMMDD"},
		{"format not D8 or RD8", func() (any, error) {
			_, _, err := parseDTPPeriod(segment("472", "DT", "202601010800"))
			return nil, err
		}, nil, "its date format qualifier is not D8 or RD8"},
		{"DT", func() (any, error) {
			d, h, err := parseDateHour(segment("435", "DT", "202602282359"), 3, "DT")
			return []any{d, h}, err
		}, []any{date(2026, 2, 28), Time{Hour: 23, Minute: 59, Set: true}}, ""},
		{"DT too short", func() (any, error) {
			_, _, err := parseDateHour(segment("435", "DT", "2026022823"), 3, "DT")
			return nil, err
		}, nil, "date and time is not CCYYMMDDHHMM"},
		{"time with seconds", func() (any, error) { return parseTime(segment("", "", "103045"), 3) }, Time{Hour: 10, Minute: 30, Second: 45, WithSeconds: true, Set: true}, ""},
		{"amount", func() (any, error) { return parseAmount(segment("", "", "-0.50"), 3) }, Amount{Value: Decimal{Coefficient: -5, Scale: 1}, Text: "-0.50", Set: true}, ""},
		{"empty amount", func() (any, error) { return parseAmount(segment(), 3) }, Amount{}, ""},
	}
	for _, tt := range tests {
		got, err := tt.parse()
		if tt.wantErr != "" {
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("%s: error = %v, want it to contain %q", tt.name, err, tt.wantErr)
			}
			continue
		}
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s = %+v, %v, want %+v, nil", tt.name, got, err, tt.want)
		}
	}
}
