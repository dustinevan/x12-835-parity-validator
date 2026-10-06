package x12837p

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
)

const repoRoot = "../../.."

const (
	officeVisit       = "testdata/x12837p/office-visit.837"
	dependentPatient  = "testdata/x12837p/dependent-patient.837"
	medicareSecondary = "testdata/x12837p/medicare-secondary.837"
)

// professionalFiles returns every *.837 under testdata/ and
// frontend/src/data whose GS08 names the 837P implementation guide
// (005010X222). The 837I files under testdata/taxonomy share the extension
// and are skipped.
func professionalFiles(t *testing.T) []string {
	t.Helper()
	var paths []string
	for _, dir := range []string{"testdata", "frontend/src/data"} {
		err := filepath.WalkDir(filepath.Join(repoRoot, dir), func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".837" {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if bytes.Contains(data, []byte("005010X222")) {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return paths
}

// TestParseFiles parses every 837P file under testdata/ and
// frontend/src/data, and checks that each segment of the input lands in the
// tree, as a typed segment or a RawSegment.
func TestParseFiles(t *testing.T) {
	paths := professionalFiles(t)
	if len(paths) < 4 {
		t.Fatalf("found %d 837P files, want the three in testdata/x12837p and the taxonomy's", len(paths))
	}
	for _, path := range paths {
		name := strings.TrimPrefix(path, repoRoot+"/")
		t.Run(name, func(t *testing.T) {
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
// the pieces that are not blank.
func countInput(data []byte) int {
	start := bytes.Index(data, []byte("ISA"))
	terminator := data[start+105]
	n := 0
	for _, piece := range bytes.Split(data[start:], []byte{terminator}) {
		if len(bytes.TrimSpace(piece)) > 0 {
			n++
		}
	}
	return n
}

// countTree counts the segments the tree holds: each typed segment that was
// present, plus each RawSegment.
func countTree(isa InterchangeISA) int {
	n := 2 + len(isa.Raw) // ISA, IEA
	for _, g := range isa.FunctionalGroups {
		n += 2 + len(g.Raw) // GS, GE
		for _, tr := range g.Transactions {
			n += countTransaction(tr)
		}
	}
	return n
}

func countTransaction(tr TransactionST) int {
	n := 2 + len(tr.Raw) + present(tr.Header.Set) // ST, SE
	n += loop(tr.SubmitterNM141.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, len(tr.SubmitterNM141.Contacts), len(tr.SubmitterNM141.Raw))
	n += loop(tr.ReceiverNM140.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, len(tr.ReceiverNM140.Raw))
	for _, b := range tr.BillingProviders {
		n += 1 + len(b.Raw) + present(b.Specialty.CodePRV01 != ProviderTypePRV01NotSent) + present(b.Currency.Set)
		p := b.ProviderNM185
		n += loop(p.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(p.Address.Set), present(p.Location.Set),
			reference(p.TaxIDREFEI), len(p.References), contact(p.Contact), len(p.Raw))
		n += loop(b.PayToProviderNM187.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(b.PayToProviderNM187.Address.Set),
			present(b.PayToProviderNM187.Location.Set), len(b.PayToProviderNM187.Raw))
		n += loop(b.PayToPlanNM1PE.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(b.PayToPlanNM1PE.Address.Set),
			present(b.PayToPlanNM1PE.Location.Set), len(b.PayToPlanNM1PE.References), len(b.PayToPlanNM1PE.Raw))
		for _, s := range b.Subscribers {
			n += 1 + len(s.Raw) + present(s.Info.Set) + present(s.PatientInfo.Set)
			n += loop(s.SubscriberNM1IL.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(s.SubscriberNM1IL.Address.Set),
				present(s.SubscriberNM1IL.Location.Set), present(s.SubscriberNM1IL.Demographics.Set), len(s.SubscriberNM1IL.References),
				contact(s.SubscriberNM1IL.Contact), len(s.SubscriberNM1IL.Raw))
			n += loop(s.PayerNM1PR.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(s.PayerNM1PR.Address.Set),
				present(s.PayerNM1PR.Location.Set), len(s.PayerNM1PR.References), len(s.PayerNM1PR.Raw))
			for _, c := range s.Claims {
				n += countClaim(c)
			}
			for _, pat := range s.Patients {
				n += 1 + len(pat.Raw) + present(pat.PatientInfo.Set)
				n += loop(pat.PatientNM1QC.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent, present(pat.PatientNM1QC.Address.Set),
					present(pat.PatientNM1QC.Location.Set), present(pat.PatientNM1QC.Demographics.Set), len(pat.PatientNM1QC.References),
					contact(pat.PatientNM1QC.Contact), len(pat.PatientNM1QC.Raw))
				for _, c := range pat.Claims {
					n += countClaim(c)
				}
			}
		}
	}
	return n
}

func countClaim(c ClaimCLM) int {
	n := 1 + len(c.Raw) + len(c.SupplementalInfoPWK) + len(c.FileInfoK3) + len(c.References)
	for _, d := range []DateDTP{
		c.OnsetDTP431, c.InitialTreatmentDTP454, c.LastSeenDTP304, c.AcuteManifestationDTP453, c.AccidentDTP439,
		c.LastMenstrualPeriodDTP484, c.LastXRayDTP455, c.HearingVisionPrescriptionDTP471, c.DisabilityDTP314, c.LastWorkedDTP297,
		c.AuthorizedReturnToWorkDTP296, c.AdmissionDTP435, c.DischargeDTP096, c.AssumedOrRelinquishedCareDTP090[0],
		c.AssumedOrRelinquishedCareDTP090[1], c.PropertyCasualtyFirstContactDTP444, c.RepricerReceivedDTP050,
	} {
		n += present(d.QualifierDTP01 != DateTimeQualifierDTP01NotSent)
	}
	for _, ref := range []ReferenceREF{
		c.ServiceAuthExceptionREF4N, c.MedicareCrossoverIndicatorREFF5, c.MammographyCertificationREFEW, c.ReferralREF9F,
		c.PriorAuthorizationREFG1, c.PayerClaimControlNumberREFF8, c.CLIAREFX4, c.RepricedClaimREF9A, c.AdjustedRepricedClaimREF9C,
		c.InvestigationalDeviceExemptionREFLX, c.TransmissionIntermediaryREFD9, c.MedicalRecordNumberREFEA,
		c.DemonstrationProjectREFP4, c.CarePlanOversightREF1J,
	} {
		n += reference(ref)
	}
	n += present(c.ContractInfoCN1.Set) + present(c.PatientAmountPaidAMTF5.QualifierAMT01 != AmountQualifierAMT01NotSent) + present(c.NoteNTE.Set)
	n += present(c.AmbulanceTransportCR1.Set) + present(c.SpinalManipulationCR2.Set) + present(c.PricingHCP.Set)
	n += len(c.AmbulanceCertificationCRC07) + len(c.VisionConditionCRCE1) + present(c.HomeboundCRC75.Set) + present(c.EPSDTReferralCRCZZ.Set)
	n += present(c.DiagnosisCodesHIABK.Set) + present(c.AnesthesiaProceduresHIBP.Set) + len(c.ConditionCodesHIBG)

	n += loop(c.ReferringProvider.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, len(c.ReferringProvider.References), len(c.ReferringProvider.Raw))
	n += loop(c.RenderingProviderNM182.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, specialty(c.RenderingProviderNM182.Specialty),
		len(c.RenderingProviderNM182.References), len(c.RenderingProviderNM182.Raw))
	f := c.ServiceFacilityNM177
	n += loop(f.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, present(f.Address.Set), present(f.Location.Set),
		len(f.References), contact(f.Contact), len(f.Raw))
	n += loop(c.SupervisingProviderNM1DQ.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, len(c.SupervisingProviderNM1DQ.References), len(c.SupervisingProviderNM1DQ.Raw))
	n += loop(c.AmbulancePickupNM1PW.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, present(c.AmbulancePickupNM1PW.Address.Set),
		present(c.AmbulancePickupNM1PW.Location.Set), len(c.AmbulancePickupNM1PW.Raw))
	n += loop(c.AmbulanceDropoffNM145.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, present(c.AmbulanceDropoffNM145.Address.Set),
		present(c.AmbulanceDropoffNM145.Location.Set), len(c.AmbulanceDropoffNM145.Raw))
	for _, o := range c.OtherSubscribers {
		n += countOtherSubscriber(o)
	}
	for _, l := range c.ServiceLines {
		n += countServiceLine(l)
	}
	return n
}

func countOtherSubscriber(o OtherSubscriberSBR) int {
	n := 1 + len(o.Adjustments) + len(o.Raw) + present(o.OtherInsurance.Set) + present(o.OutpatientAdjudication.Set)
	n += amountSlot(o.PayerPaidAMTD) + amountSlot(o.NonCoveredAMTA8) + amountSlot(o.RemainingLiabilityAMTEAF)
	s := o.OtherSubscriberNM1IL
	n += loop(s.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, present(s.Address.Set), present(s.Location.Set),
		reference(s.SecondaryID), len(s.Raw))
	p := o.OtherPayerNM1PR
	n += loop(p.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, present(p.Address.Set), present(p.Location.Set),
		present(p.RemittanceDateDTP573.QualifierDTP01 != DateTimeQualifierDTP01NotSent), len(p.SecondaryID),
		reference(p.PriorAuthorizationREFG1), reference(p.ReferralREF9F), reference(p.ClaimAdjustmentIndicatorREFT4),
		reference(p.ClaimControlNumberREFF8), len(p.Raw))
	for _, r := range []OtherPayerRoleNM1{o.OtherPayerReferringProvider, o.OtherPayerRenderingProviderNM182,
		o.OtherPayerServiceFacilityNM177, o.OtherPayerSupervisingProviderNM1DQ, o.OtherPayerBillingProviderNM185} {
		n += loop(r.Name.EntityIDNM101 != EntityIdentifierNM101NotSent, len(r.References), len(r.Raw))
	}
	return n
}

func countServiceLine(l ServiceLineLX) int {
	n := 1 + len(l.Raw) + present(l.Service.Set) + present(l.DME.Set) + len(l.Attachments)
	n += present(l.AmbulanceTransport.Set) + present(l.DMECertification.Set)
	n += len(l.AmbulanceCertificationCRC07) + present(l.HospiceEmployeeCRC70.Set) + present(l.DMEConditionCRC09.Set)
	for _, d := range []DateDTP{l.ServiceDateDTP472, l.PrescriptionDTP471, l.CertificationRevisionDTP607, l.BeginTherapyDTP463,
		l.LastCertificationDTP461, l.LastTreatmentOrTherapyDTP304, l.ShippedDTP011, l.LastXRayDTP455, l.InitialTreatmentDTP454} {
		n += present(d.QualifierDTP01 != DateTimeQualifierDTP01NotSent)
	}
	n += len(l.TestPerformedDTP738)
	n += present(l.AmbulancePatientCountQTYPT.QualifierQTY01 != QuantityQualifierQTY01NotSent)
	n += present(l.ObstetricAdditionalUnitsQTYFL.QualifierQTY01 != QuantityQualifierQTY01NotSent)
	n += len(l.TestResults) + present(l.Contract.Set)
	for _, ref := range []ReferenceREF{l.RepricedLineItemREF9B, l.AdjustedRepricedLineItemREF9D, l.LineControlREF6R, l.MammographyCertificationREFEW,
		l.CLIAREFX4, l.ReferringCLIAREFF4, l.ImmunizationBatchREFBT} {
		n += reference(ref)
	}
	n += len(l.PriorAuthorizationREFG1) + len(l.ReferralREF9F)
	n += amountSlot(l.SalesTaxAMTT) + amountSlot(l.PostageClaimedAMTF4)
	n += len(l.FileInfoK3) + present(l.LineNote.Set) + present(l.ThirdPartyNote.Set)
	n += present(l.PurchasedService.Set) + present(l.Pricing.Set)
	n += loop(l.Drug.Set, present(l.Drug.Pricing.Set), reference(l.Drug.AssociationREFXZ), len(l.Drug.Raw))
	for _, p := range []LineProviderNM1{l.RenderingNM182, l.PurchasedServiceNM1QB, l.ServiceFacilityNM177,
		l.SupervisingNM1DQ, l.OrderingNM1DK, l.ReferringNM12420F, l.AmbulancePickupNM1PW, l.AmbulanceDropoffNM145} {
		n += loop(p.Set, specialty(p.Specialty), present(p.Address.Set), present(p.Location.Set), contact(p.Contact),
			len(p.References), len(p.Raw))
	}
	for _, a := range l.Adjudications {
		n += 1 + len(a.Adjustments) + len(a.Raw)
		n += present(a.AdjudicationDateDTP573.QualifierDTP01 != DateTimeQualifierDTP01NotSent)
		n += amountSlot(a.RemainingLiabilityAMTEAF)
	}
	for _, f := range l.Forms {
		n += 1 + len(f.Answers) + len(f.Raw)
	}
	return n
}

// loop counts a loop's segments: its first segment, when the loop was
// placed, plus the counts of what it holds.
func loop(placed bool, counts ...int) int {
	if !placed {
		return 0
	}
	n := 1
	for _, c := range counts {
		n += c
	}
	return n
}

func contact(c ContactPER) int {
	return present(c.FunctionPER01 != ContactFunctionCodePER01NotSent)
}

func specialty(p ProviderPRV) int {
	return present(p.CodePRV01 != ProviderTypePRV01NotSent)
}

// reference counts a REF held in a qualifier-specific field.
func reference(r ReferenceREF) int {
	return present(r.QualifierREF01 != ReferenceIdentificationQualifierREF01NotSent)
}

// amountSlot counts an AMT held in a qualifier-specific field.
func amountSlot(a AmountAMT) int {
	return present(a.QualifierAMT01 != AmountQualifierAMT01NotSent)
}

func present(b bool) int {
	if b {
		return 1
	}
	return 0
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

// amount returns an Amount as the parser builds it from text: the parsed
// value and the text itself.
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

// cents returns an amount in hundredths. Every amount in the synthetic files
// has at most two decimal places.
func cents(a Amount) int64 {
	c := a.Value.Coefficient
	for s := a.Value.Scale; s < 2; s++ {
		c *= 10
	}
	return c
}

func firstTransaction(isa InterchangeISA) TransactionST {
	return isa.FunctionalGroups[0].Transactions[0]
}

// TestParseOfficeVisit spot-checks typed values in office-visit.837: the
// envelope, the submitter, the billing provider, the subscriber, and a
// Medicare primary claim with two lines. Each segment's Index is its
// 0-based position in the file.
func TestParseOfficeVisit(t *testing.T) {
	isa := parseText(t, readFile(t, officeVisit))
	tr := firstTransaction(isa)
	if len(tr.BillingProviders) != 1 || len(tr.BillingProviders[0].Subscribers) != 1 {
		t.Fatalf("billing providers, subscribers = %d, want 1 and 1", len(tr.BillingProviders))
	}
	billing := tr.BillingProviders[0]
	subscriber := billing.Subscribers[0]
	if len(subscriber.Claims) != 1 || len(subscriber.Claims[0].ServiceLines) != 2 {
		t.Fatalf("claims = %d, want 1 with 2 lines", len(subscriber.Claims))
	}
	claim := subscriber.Claims[0]
	first := claim.ServiceLines[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"delimiters", isa.Delimiters, Delimiters{Element: "*", Segment: "~", AfterSegment: "\n"}},
		{"ISA09", isa.DateISA09, date(2026, 9, 15)},
		{"ISA14", isa.AckRequestedISA14, AckRequestedISA141},
		{"ISA15", isa.UsageISA15, UsageISA15T},
		{"GS index", tr.Index, 2},
		{"BHT", tr.Header, HeaderBHT{Set: true, Index: 3, StructureBHT01: HierarchicalStructureBHT010019, PurposeBHT02: TransactionSetPurposeBHT0200,
			OriginatorTransactionIDBHT03: "SYNBATCH101", DateBHT04: date(2026, 9, 15), TimeBHT05: Time{Hour: 10, Minute: 30, Set: true},
			TransactionTypeBHT06: TransactionTypeBHT06CH}},
		{"submitter PER", tr.SubmitterNM141.Contacts, []ContactPER{{Index: 5, FunctionPER01: ContactFunctionCodePER01IC, NamePER02: "SYNTHETIC CONTACT",
			QualifierPER03: CommunicationNumberQualifierPER03TE, CommPER04: "8005550100",
			QualifierPER05: CommunicationNumberQualifierPER03EM, CommPER06: "CONTACT@EXAMPLE.TEST"}}},
		{"receiver NM109", tr.ReceiverNM140.Name.IDNM109, "SYNMAC001"},
		{"2000A HL", billing.HL, HierarchyHL{Index: 7, IDNumberHL01: "1", LevelHL03: HierarchicalLevelHL0320, ChildCodeHL04: HierarchicalChildCodeHL041}},
		{"2000A PRV", billing.Specialty, ProviderPRV{Index: 8, CodePRV01: ProviderTypePRV01BI, QualifierPRV02: ReferenceIdentificationQualifierREF01PXC, TaxonomyPRV03: "207Q00000X"}},
		{"2010AA NM103", billing.ProviderNM185.NM1.LastNM103, "SYNTHETIC FAMILY CLINIC"},
		{"2010AA REF*EI", billing.ProviderNM185.TaxIDREFEI, ReferenceREF{Index: 12, QualifierREF01: ReferenceIdentificationQualifierREF01EI, Value: "990000001"}},
		{"2010AA N4", billing.ProviderNM185.Location, LocationN4{Set: true, Index: 11, CityN401: "SYNTHETIC CITY", StateN402: "CO", ZipN403: "802021234"}},
		{"2000B HL", subscriber.HL, HierarchyHL{Index: 13, IDNumberHL01: "2", ParentIDNumberHL02: "1", LevelHL03: HierarchicalLevelHL0322, ChildCodeHL04: HierarchicalChildCodeHL040}},
		{"SBR", subscriber.Info, SubscriberSBR{Set: true, Index: 14, PayerResponsibilitySBR01: PayerResponsibilitySBR01P, RelationshipSBR02: RelationshipSBR0218,
			ClaimFilingIndicatorSBR09: ClaimFilingIndicatorSBR09MB}},
		{"2010BA NM1", subscriber.SubscriberNM1IL.NM1, NameNM1{Index: 15, EntityIDNM101: EntityIdentifierNM101IL, EntityTypeNM102: EntityTypeQualifierNM1021,
			LastNM103: "SYNTHETIC", FirstNM104: "ALEX", MiddleNM105: "Q", IDQualifierNM108: IdentificationCodeQualifierNM108MI, IDNM109: "1SYN0000AA01"}},
		{"2010BA DMG", subscriber.SubscriberNM1IL.Demographics, DemographicDMG{Set: true, Index: 18, FormatDMG01: "D8", BirthDateDMG02: date(1950, 1, 1), GenderDMG03: "F"}},
		{"2010BB NM1", subscriber.PayerNM1PR.NM1.EntityIDNM101, EntityIdentifierNM101PR},
		{"CLM index", claim.Index, 20},
		{"CLM02", claim.TotalChargeCLM02, amount(t, "175.50")},
		{"CLM05-1", claim.PlaceOfServiceCLM05_1, "11"},
		{"CLM05-2", claim.FacilityQualifierCLM05_2, FacilityCodeQualifierCLM05_2B},
		{"CLM05-3", claim.FrequencyCLM05_3, ClaimFrequencyCLM05_31},
		{"CLM06-09", []any{claim.SignatureCLM06, claim.AssignmentCLM07, claim.BenefitsAssignmentCLM08, claim.ReleaseOfInfoCLM09},
			[]any{ProviderOrSupplierSignatureCLM06Y, AssignmentOrPlanParticipationCLM07A, BenefitsAssignmentCertificationCLM08Y, ReleaseOfInformationCLM09Y}},
		{"DTP*431", claim.OnsetDTP431, DateDTP{Index: 21, QualifierDTP01: DateTimeQualifierDTP01431, FormatDTP02: "D8", Date: date(2026, 9, 1), RawDTP03: "20260901"}},
		{"REF*G1", claim.PriorAuthorizationREFG1, ReferenceREF{Index: 22, QualifierREF01: ReferenceIdentificationQualifierREF01G1, Value: "SYNAUTH01"}},
		{"HI", claim.DiagnosisCodesHIABK, HealthCareInfoHI{Set: true, Index: 23, Codes: []HealthCareCodeHI{{TypeHI0x1: DiagnosisTypeCodeHI0x_1ABK, CodeHI0x2: "J069"}, {TypeHI0x1: DiagnosisTypeCodeHI0x_1ABF, CodeHI0x2: "R05"}}}},
		{"2310B NM109", claim.RenderingProviderNM182.Name.IDNM109, "1999999992"},
		{"2310B PRV01", claim.RenderingProviderNM182.Specialty.CodePRV01, ProviderTypePRV01PE},
		{"LX01", first.AssignedNumberLX01, "1"},
		{"LX index", first.Index, 26},
		{"SV101", first.Service.ProcedureSV101, ProcedureSV101{Qualifier: "HC", Code: "99213", Modifiers: [4]string{"25"}}},
		{"SV1 index", first.Service.Index, 27},
		{"SV102", first.Service.ChargeSV102, amount(t, "125.50")},
		{"SV103", first.Service.UnitSV103, UnitOrBasisForMeasurementCodeUN},
		{"SV104", first.Service.UnitCountSV104, amount(t, "1")},
		{"SV107", first.Service.DiagnosisPointersSV107, [4]int{1, 2}},
		{"DTP*472", first.ServiceDateDTP472.Date, date(2026, 9, 10)},
		{"REF*6R", first.LineControlREF6R.Value, "SYNLINE0001"},
		{"line raw", len(first.Raw), 0},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
	var lines int64
	for _, l := range claim.ServiceLines {
		lines += cents(l.Service.ChargeSV102)
	}
	if lines != cents(claim.TotalChargeCLM02) {
		t.Errorf("SV102 total = %d cents, CLM02 = %d cents, want equal", lines, cents(claim.TotalChargeCLM02))
	}
}

// TestParseDependentPatient checks that a 2000C patient loop attaches to the
// 2000B its HL02 names and that the claim after it belongs to the patient.
func TestParseDependentPatient(t *testing.T) {
	tr := firstTransaction(parseText(t, readFile(t, dependentPatient)))
	subscriber := tr.BillingProviders[0].Subscribers[0]
	if len(subscriber.Patients) != 1 {
		t.Fatalf("patients = %d, want 1", len(subscriber.Patients))
	}
	patient := subscriber.Patients[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"BHT05 with seconds", tr.Header.TimeBHT05, Time{Hour: 8, Minute: 15, HasSeconds: true, Set: true}},
		{"2000B HL04", subscriber.HL.ChildCodeHL04, HierarchicalChildCodeHL041},
		{"SBR02 not sent", subscriber.Info.RelationshipSBR02, RelationshipSBR02NotSent},
		{"SBR09 commercial", subscriber.Info.ClaimFilingIndicatorSBR09, ClaimFilingIndicatorSBR09CI},
		{"SBR09 on the list", subscriber.Info.ClaimFilingIndicatorSBR09.Known(), true},
		{"subscriber claims", len(subscriber.Claims), 0},
		{"2000C HL", patient.HL, HierarchyHL{Index: 19, IDNumberHL01: "3", ParentIDNumberHL02: "2", LevelHL03: HierarchicalLevelHL0323, ChildCodeHL04: HierarchicalChildCodeHL040}},
		{"PAT01", patient.PatientInfo, PatientInfoPAT{Set: true, Index: 20, RelationshipPAT01: RelationshipPAT0119}},
		{"2010CA NM1", patient.PatientNM1QC.NM1, NameNM1{Index: 21, EntityIDNM101: EntityIdentifierNM101QC, EntityTypeNM102: EntityTypeQualifierNM1021, LastNM103: "SYNTHETIC", FirstNM104: "RILEY"}},
		{"2010CA DMG", patient.PatientNM1QC.Demographics.BirthDateDMG02, date(2015, 3, 3)},
		{"patient claims", len(patient.Claims), 1},
		{"patient claim CLM01", patient.Claims[0].PatientControlCLM01, "SYNPCN0002"},
		{"patient claim SV102", patient.Claims[0].ServiceLines[0].Service.ChargeSV102, amount(t, "90.00")},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
}

// TestParseMedicareSecondary checks the 2320/2330 coordination-of-benefits
// loops and the 2430 line adjudication in medicare-secondary.837, and that
// the file's amounts balance: the lines add up to the claim, each line's
// other-payer payment plus its adjustments equals its charge, and the line
// payments add up to the 2320 AMT*D.
func TestParseMedicareSecondary(t *testing.T) {
	tr := firstTransaction(parseText(t, readFile(t, medicareSecondary)))
	subscriber := tr.BillingProviders[0].Subscribers[0]
	claim := subscriber.Claims[0]
	if len(claim.OtherSubscribers) != 1 || len(claim.ServiceLines) != 2 {
		t.Fatalf("other subscribers, lines = %d, %d, want 1, 2", len(claim.OtherSubscribers), len(claim.ServiceLines))
	}
	other := claim.OtherSubscribers[0]
	adjudication := claim.ServiceLines[0].Adjudications[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"2000B SBR01", subscriber.Info.PayerResponsibilitySBR01, PayerResponsibilitySBR01S},
		{"2320 SBR", other.Subscriber, SubscriberSBR{Set: true, Index: 21, PayerResponsibilitySBR01: PayerResponsibilitySBR01P, RelationshipSBR02: RelationshipSBR0218,
			GroupOrPolicyNumberSBR03: "SYNEMPGRP03", ClaimFilingIndicatorSBR09: ClaimFilingIndicatorSBR09CI}},
		{"AMT*D", other.PayerPaidAMTD, AmountAMT{Index: 22, QualifierAMT01: AmountQualifierAMT01D, Amount: amount(t, "120.00")}},
		{"AMT*A8 not sent", other.NonCoveredAMTA8, AmountAMT{}},
		{"AMT*EAF", other.RemainingLiabilityAMTEAF, AmountAMT{Index: 23, QualifierAMT01: AmountQualifierAMT01EAF, Amount: amount(t, "30.00")}},
		{"OI", other.OtherInsurance, OtherInsuranceOI{Set: true, Index: 24, BenefitsAssignmentOI03: BenefitsAssignmentCertificationCLM08Y,
			SignatureSourceOI04: PatientSignatureSourceCLM10P, ReleaseOfInfoOI06: ReleaseOfInformationCLM09Y}},
		{"2330A NM109", other.OtherSubscriberNM1IL.Name.IDNM109, "SYNEMP0003"},
		{"2330A N4", other.OtherSubscriberNM1IL.Location.ZipN403, "850021234"},
		{"2330B NM109", other.OtherPayerNM1PR.Name.IDNM109, "SYNEGHP03"},
		{"2330B DTP*573", other.OtherPayerNM1PR.RemittanceDateDTP573.Date, date(2026, 9, 15)},
		{"2330B REF*F8", other.OtherPayerNM1PR.ClaimControlNumberREFF8, ReferenceREF{Index: 30, QualifierREF01: ReferenceIdentificationQualifierREF01F8, Value: "SYNPRIMARY0003"}},
		{"SVD", adjudication.ProcedureSVD03, ProcedureSV101{Qualifier: "HC", Code: "20610"}},
		{"SVD05", adjudication.PaidUnitsSVD05, amount(t, "1")},
		{"2430 CAS", adjudication.Adjustments, []AdjustmentCAS{
			{Index: 36, GroupCAS01: ClaimAdjustmentGroupCodeCAS01CO, Reasons: []CASReason{{Reason: carc.Code("45"), Amount: amount(t, "40.00")}}},
			{Index: 37, GroupCAS01: ClaimAdjustmentGroupCodeCAS01PR, Reasons: []CASReason{{Reason: carc.Code("1"), Amount: amount(t, "20.00")}}},
		}},
		{"2430 DTP*573", adjudication.AdjudicationDateDTP573.QualifierDTP01, DateTimeQualifierDTP01573},
		{"claim raw", len(claim.Raw), 0},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
	var charges, paid int64
	for i, l := range claim.ServiceLines {
		charges += cents(l.Service.ChargeSV102)
		a := l.Adjudications[0]
		paid += cents(a.PaidAmountSVD02)
		accounted := cents(a.PaidAmountSVD02)
		for _, cas := range a.Adjustments {
			for _, r := range cas.Reasons {
				accounted += cents(r.Amount)
			}
		}
		if accounted != cents(l.Service.ChargeSV102) {
			t.Errorf("line %d: SVD02 plus CAS = %d cents, SV102 = %d cents, want equal", i+1, accounted, cents(l.Service.ChargeSV102))
		}
	}
	if charges != cents(claim.TotalChargeCLM02) {
		t.Errorf("SV102 total = %d cents, CLM02 = %d cents, want equal", charges, cents(claim.TotalChargeCLM02))
	}
	if paid != cents(other.PayerPaidAMTD.Amount) {
		t.Errorf("SVD02 total = %d cents, AMT*D = %d cents, want equal", paid, cents(other.PayerPaidAMTD.Amount))
	}
}

// TestParseHierarchyParent appends a second 2000A and then a 2000B whose
// HL02 names the first 2000A, and checks that the subscriber attaches to
// the HL it names, not to the most recent one.
func TestParseHierarchyParent(t *testing.T) {
	text := strings.Replace(readFile(t, officeVisit), "SE*33*0001~",
		"HL*3**20*1~NM1*85*2*SECOND CLINIC*****XX*1999999950~HL*4*1*22*0~SBR*P*18*******MB~SE*37*0001~", 1)
	tr := firstTransaction(parseText(t, text))
	if len(tr.BillingProviders) != 2 {
		t.Fatalf("billing providers = %d, want 2", len(tr.BillingProviders))
	}
	if got := len(tr.BillingProviders[0].Subscribers); got != 2 {
		t.Errorf("first billing provider's subscribers = %d, want 2", got)
	}
	if got := len(tr.BillingProviders[1].Subscribers); got != 0 {
		t.Errorf("second billing provider's subscribers = %d, want 0", got)
	}
	if got := tr.BillingProviders[0].Subscribers[1].HL.IDNumberHL01; got != "4" {
		t.Errorf("second subscriber HL01 = %q, want \"4\"", got)
	}
}

// coded is what every coded element's type offers: its text and whether
// the text is on its list.
type coded interface {
	String() string
	Known() bool
}

// TestParseOffList changes one coded value in office-visit.837 to a code
// off its list and checks that the field holds the text as sent and reports
// Known false.
func TestParseOffList(t *testing.T) {
	claimOf := func(isa InterchangeISA) ClaimCLM {
		return firstTransaction(isa).BillingProviders[0].Subscribers[0].Claims[0]
	}
	tests := []struct {
		name     string
		old, new string
		get      func(InterchangeISA) coded
		want     string
	}{
		{"HL04", "HL*2*1*22*0~", "HL*2*1*22*7~", func(isa InterchangeISA) coded {
			return firstTransaction(isa).BillingProviders[0].Subscribers[0].HL.ChildCodeHL04
		}, "7"},
		{"CLM05-2 component", "*11:B:1*", "*11:X:1*", func(isa InterchangeISA) coded { return claimOf(isa).FacilityQualifierCLM05_2 }, "X"},
		{"NM102", "NM1*82*1*", "NM1*82*3*", func(isa InterchangeISA) coded { return claimOf(isa).RenderingProviderNM182.Name.EntityTypeNM102 }, "3"},
		{"HI02-1 component", "*ABF:R05~", "*ABX:R05~", func(isa InterchangeISA) coded { return claimOf(isa).DiagnosisCodesHIABK.Codes[1].TypeHI0x1 }, "ABX"},
		{"HI02-9 component", "*ABF:R05~", "*ABF:R05:::::::U~", func(isa InterchangeISA) coded { return claimOf(isa).DiagnosisCodesHIABK.Codes[1].YesNoHI0x9 }, "U"},
		{"SV103 on a line", "*50.00*UN*", "*50.00*ZZ*", func(isa InterchangeISA) coded { return claimOf(isa).ServiceLines[1].Service.UnitSV103 }, "ZZ"},
		{"SBR09", "SBR*P*18*******MB~", "SBR*P*18*******QQ~", func(isa InterchangeISA) coded {
			return firstTransaction(isa).BillingProviders[0].Subscribers[0].Info.ClaimFilingIndicatorSBR09
		}, "QQ"},
		{"ISA15", "*1*T*:~", "*1*X*:~", func(isa InterchangeISA) coded { return isa.UsageISA15 }, "X"},
	}
	base := readFile(t, officeVisit)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("office-visit.837 does not contain %q", tt.old)
			}
			text := strings.Replace(base, tt.old, tt.new, 1)
			field := tt.get(parseText(t, text))
			if field.String() != tt.want || field.Known() {
				t.Errorf("field = %q, Known %v, want %q, Known false", field.String(), field.Known(), tt.want)
			}
			got, err := Bytes(parseText(t, text))
			if err != nil || string(got) != text {
				t.Errorf("the off-list code did not round-trip: %v", err)
			}
		})
	}
}

// TestParseRaw inserts or changes segments in office-visit.837 so that the
// model has no slot for them, and checks that each is kept as a RawSegment
// on the struct that encloses it, with its file position, and that nothing
// is dropped.
func TestParseRaw(t *testing.T) {
	claimOf := func(isa InterchangeISA) ClaimCLM {
		return firstTransaction(isa).BillingProviders[0].Subscribers[0].Claims[0]
	}
	tests := []struct {
		name     string
		old, new string
		se       string // the SE that replaces SE*33* when the change adds segments
		get      func(InterchangeISA) []RawSegment
		want     []RawSegment
	}{
		{
			name: "unmodeled segment in a claim", old: "REF*G1*SYNAUTH01~", new: "REF*G1*SYNAUTH01~\nXYZ*1~", se: "SE*34*",
			get:  func(isa InterchangeISA) []RawSegment { return claimOf(isa).Raw },
			want: []RawSegment{{Index: 23, ID: "XYZ", Elements: []string{"1"}}},
		},
		{
			name: "second DTP*431 in a claim", old: "DTP*431*D8*20260901~", new: "DTP*431*D8*20260901~\nDTP*431*D8*20260902~", se: "SE*34*",
			get:  func(isa InterchangeISA) []RawSegment { return claimOf(isa).Raw },
			want: []RawSegment{{Index: 22, ID: "DTP", Elements: []string{"431", "D8", "20260902"}}},
		},
		{
			name: "NM1 role with no claim slot", old: "LX*1~", new: "NM1*ZZ*1*OTHER~\nREF*G2*X1~\nLX*1~", se: "SE*35*",
			get: func(isa InterchangeISA) []RawSegment { return claimOf(isa).Raw },
			want: []RawSegment{
				{Index: 26, ID: "NM1", Elements: []string{"ZZ", "1", "OTHER"}},
				{Index: 27, ID: "REF", Elements: []string{"G2", "X1"}},
			},
		},
		{
			name: "second REF*6R on a line", old: "REF*6R*SYNLINE0002~", new: "REF*6R*SYNLINE0002~\nREF*6R*SYNLINE0003~", se: "SE*34*",
			get:  func(isa InterchangeISA) []RawSegment { return claimOf(isa).ServiceLines[1].Raw },
			want: []RawSegment{{Index: 34, ID: "REF", Elements: []string{"6R", "SYNLINE0003"}}},
		},
		{
			name: "HL whose parent is not found", old: "HL*2*1*22*0~", new: "HL*2*9*22*0~",
			get: func(isa InterchangeISA) []RawSegment {
				raw := firstTransaction(isa).Raw
				return raw[:1]
			},
			want: []RawSegment{{Index: 13, ID: "HL", Elements: []string{"2", "9", "22", "0"}}},
		},
		{
			name: "unmodeled segment in a group", old: "ST*837*0001", new: "XYZ*1~\nST*837*0001",
			get:  func(isa InterchangeISA) []RawSegment { return isa.FunctionalGroups[0].Raw },
			want: []RawSegment{{Index: 2, ID: "XYZ", Elements: []string{"1"}}},
		},
	}
	base := readFile(t, officeVisit)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("office-visit.837 does not contain %q", tt.old)
			}
			text := strings.Replace(base, tt.old, tt.new, 1)
			if tt.se != "" {
				text = strings.Replace(text, "SE*33*", tt.se, 1)
			}
			isa := parseText(t, text)
			if got := tt.get(isa); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Raw = %+v, want %+v", got, tt.want)
			}
			if got, want := countTree(isa), countInput([]byte(text)); got != want {
				t.Errorf("countTree = %d, want %d", got, want)
			}
			got, err := Bytes(isa)
			if err != nil || string(got) != text {
				t.Errorf("the raw segments did not round-trip in place: %v\n%s", err, firstDifference(got, []byte(text)))
			}
		})
	}
}

// TestParseErrors checks that Parse rejects malformed input and that each
// error names the element and the segment.
func TestParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		old, new string
		want     string
	}{
		{"SE01", "SE*33*0001", "SE*32*0001", `x12837p: checking SE01 of transaction "0001": SE01 is 32, ST through SE is 33 segments`},
		{"SE02", "SE*33*0001", "SE*33*0002", `x12837p: checking SE02: ST02 is "0001", SE02 is "0002"`},
		{"GE01", "GE*1*101~", "GE*2*101~", `x12837p: checking GE01 of group "101": GE01 is 2, the group holds 1 transactions`},
		{"GE02", "GE*1*101~", "GE*1*107~", `x12837p: checking GE02: GS06 is "101", GE02 is "107"`},
		{"IEA01", "IEA*1*", "IEA*2*", `x12837p: checking IEA01: IEA01 is 2, the interchange holds 1 functional groups`},
		{"IEA02", "IEA*1*000000101", "IEA*1*000000106", `x12837p: checking IEA02: ISA13 is "000000101", IEA02 is "000000106"`},
		{"GS01 not a claim", "GS*HC*", "GS*HP*", `x12837p: parsing GS01 in segment 1: "HP" is not a health care claim functional group`},
		{"ST01 not an 837", "ST*837*0001", "ST*835*0001", `x12837p: parsing ST01 in segment 2: "835" is not an 837 transaction set`},
		{"malformed amount", "CLM*SYNPCN0001*175.50", "CLM*SYNPCN0001*17S.50", "x12837p: parsing CLM02 in segment 20: "},
		{"not a calendar date", "DMG*D8*19500101", "DMG*D8*19501301", "x12837p: parsing DMG02 in segment 18: not a calendar date"},
		{"D8 date not CCYYMMDD", "DTP*431*D8*20260901", "DTP*431*D8*2026091", "x12837p: parsing DTP03 in segment 21: date is not CCYYMMDD"},
		{"diagnosis pointer 0", "***1:2~", "***0:2~", "x12837p: parsing SV107-1 in segment 27: diagnosis pointer is not a whole number from 1"},
		{"HI date not a calendar date", "HI*ABK:J069*", "HI*ABK:J069:D8:20260931*", "x12837p: parsing HI01-4 in segment 23: not a calendar date"},
		{"HI date range without a dash", "HI*ABK:J069*", "HI*ABK:J069:RD8:20260901*", "x12837p: parsing HI01-4 in segment 23: date range is not CCYYMMDD-CCYYMMDD"},
		{"HI date format not D8 or RD8", "HI*ABK:J069*", "HI*ABK:J069:DT:202609010800*", "x12837p: parsing HI01-4 in segment 23: its date format qualifier is not D8 or RD8"},
		{"HI amount malformed", "HI*ABK:J069*", "HI*ABK:J069:::12,50*", "x12837p: parsing HI01-5 in segment 23: "},
		{"HI quantity malformed in a later composite", "*ABF:R05~", "*ABF:R05::::2X~", "x12837p: parsing HI02-6 in segment 23: "},
		{"missing SE", "SE*33*0001~", "", "x12837p: parsing segment 34: GE inside the transaction at segment 2 (missing SE)"},
		{"input after IEA", "IEA*1*000000101~", "IEA*1*000000101~TA1*1~", "x12837p: parsing segment 37: input continues after the IEA"},
		{"missing terminator", "IEA*1*000000101~", "IEA*1*000000101", "x12837p: splitting segment 36: missing segment terminator"},
	}
	base := readFile(t, officeVisit)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("office-visit.837 does not contain %q", tt.old)
			}
			_, err := Parse([]byte(strings.Replace(base, tt.old, tt.new, 1)))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Parse error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

// TestParseExtraKept checks that text after the last position the model
// names, in an element or a component, never fails Parse and lands in the
// struct's extra field, in order, with interior empties kept.
func TestParseExtraKept(t *testing.T) {
	base := readFile(t, officeVisit)
	tests := []struct {
		name     string
		old, new string
		get      func(ClaimCLM) []string
		want     []string
	}{
		{"element after LX01", "LX*1~", "LX*1*2~", func(c ClaimCLM) []string { return c.ServiceLines[0].Extra }, []string{"2"}},
		{"elements after LX01 with a gap", "LX*1~", "LX*1**3~", func(c ClaimCLM) []string { return c.ServiceLines[0].Extra }, []string{"", "3"}},
		{"component after SV101-7", "SV1*HC:87880*", "SV1*HC:87880:A:B:C:D:TEXT:X*", func(c ClaimCLM) []string { return c.ServiceLines[1].Service.ProcedureSV101.Extra }, []string{"X"}},
		{"component after HI01-9", "HI*ABK:J069*", "HI*ABK:J069::::::::X*", func(c ClaimCLM) []string { return c.DiagnosisCodesHIABK.Codes[0].Extra }, []string{"X"}},
		{"component after SV107-4", "***1:2~", "***1:2:3:4:5~", func(c ClaimCLM) []string { return c.ServiceLines[0].Service.PointersExtraSV107 }, []string{"5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("office-visit.837 does not contain %q", tt.old)
			}
			isa := parseText(t, strings.Replace(base, tt.old, tt.new, 1))
			claim := firstTransaction(isa).BillingProviders[0].Subscribers[0].Claims[0]
			if got := tt.get(claim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extra = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestParseNotUsedPositionsKept checks that text in a Medicare "Not Used"
// position never fails Parse and lands in that position's typed field: the
// parser models every X12-defined position, and Medicare's usage rules are
// a caller concern (see Parse's doc comment).
func TestParseNotUsedPositionsKept(t *testing.T) {
	base := readFile(t, officeVisit)

	t.Run("SBR06-08 element", func(t *testing.T) {
		text := strings.Replace(base, "SBR*P*18*******MB", "SBR*P*18****1*2*3*MB", 1)
		isa := parseText(t, text)
		sbr := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Info
		if sbr.CoordinationSBR06 != "1" || sbr.ResponseSBR07 != "2" || sbr.EmploymentSBR08 != "3" {
			t.Errorf("SBR06/07/08 = %q/%q/%q, want \"1\"/\"2\"/\"3\"", sbr.CoordinationSBR06, sbr.ResponseSBR07, sbr.EmploymentSBR08)
		}
	})

	t.Run("N405-06 element", func(t *testing.T) {
		text := strings.Replace(base, "N4*SYNTHETIC CITY*CO*802021234~", "N4*SYNTHETIC CITY*CO*802021234**IN*SYNLOC001~", 1)
		isa := parseText(t, text)
		n4 := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].ProviderNM185.Location
		if n4.LocationQualifierN405 != "IN" || n4.LocationIdentifierN406 != "SYNLOC001" {
			t.Errorf("N405/N406 = %q/%q, want \"IN\"/\"SYNLOC001\"", n4.LocationQualifierN405, n4.LocationIdentifierN406)
		}
	})

	t.Run("CLM11-3 component", func(t *testing.T) {
		text := strings.Replace(base, "11:B:1*Y*A*Y*Y~", "11:B:1*Y*A*Y*Y**AA::X~", 1)
		isa := parseText(t, text)
		claim := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Claims[0]
		if claim.RelatedCausesCLM11_3 != "X" {
			t.Errorf("RelatedCausesCLM11_3 = %q, want \"X\"", claim.RelatedCausesCLM11_3)
		}
	})

	// The first composite carries every component, in the shape the 837I
	// fixtures use: a D8 date in 3 and 4, an amount in 5, a quantity in 6,
	// a version in 7, a code in 8, and a yes/no in 9. The second is empty,
	// and the third carries an RD8 range.
	t.Run("HI0x-3 through HI0x-9 components", func(t *testing.T) {
		text := strings.Replace(base, "HI*ABK:J069*ABF:R05~",
			"HI*ABK:J069:D8:20260901:12.50:2.5:V1:IND:Y**ABF:R05:RD8:20260901-20260905~", 1)
		isa := parseText(t, text)
		claim := isa.FunctionalGroups[0].Transactions[0].BillingProviders[0].Subscribers[0].Claims[0]
		codes := claim.DiagnosisCodesHIABK.Codes
		if len(codes) != 3 {
			t.Fatalf("DiagnosisCodesHIABK holds %d composites, want 3 (the empty one keeps its position)", len(codes))
		}
		full, empty, span := codes[0], codes[1], codes[2]
		tests := []struct {
			name string
			got  any
			want any
		}{
			{"HI01-1", full.TypeHI0x1, DiagnosisTypeCodeHI0x_1ABK},
			{"HI01-2", full.CodeHI0x2, "J069"},
			{"HI01-3", full.DateQualifierHI0x3, "D8"},
			{"HI01-4", full.DateHI0x4, date(2026, 9, 1)},
			{"HI01-4 second date", full.DateThroughHI0x4, Date{}},
			{"HI01-5", full.AmountHI0x5, amount(t, "12.50")},
			{"HI01-6", full.QuantityHI0x6, amount(t, "2.5")},
			{"HI01-7", full.VersionIdentifierHI0x7, "V1"},
			{"HI01-8", full.IndustryCodeHI0x8, "IND"},
			{"HI01-9", full.YesNoHI0x9, YesNoConditionResponseCodeY},
			{"HI02 empty", empty, HealthCareCodeHI{}},
			{"HI03 RD8", span, HealthCareCodeHI{TypeHI0x1: DiagnosisTypeCodeHI0x_1ABF, CodeHI0x2: "R05", DateQualifierHI0x3: "RD8",
				DateHI0x4: date(2026, 9, 1), DateThroughHI0x4: date(2026, 9, 5)}},
		}
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
			}
		}
		got, err := Bytes(isa)
		if err != nil || string(got) != text {
			t.Errorf("the full composites did not round-trip: %v\n%s", err, firstDifference(got, []byte(text)))
		}
	})
}

// TestParseLineEndings checks that whitespace between segments, and a
// carriage return before a newline terminator, do not change the tree
// beyond the recorded line break, and that the line break is recorded.
func TestParseLineEndings(t *testing.T) {
	text := readFile(t, officeVisit)
	want := parseText(t, text)
	got := parseText(t, strings.ReplaceAll(text, "~\n", "~\r\n"))
	if got.Delimiters.AfterSegment != "\r\n" || want.Delimiters.AfterSegment != "\n" {
		t.Errorf("AfterSegment = %q and %q, want \"\\r\\n\" and \"\\n\"", got.Delimiters.AfterSegment, want.Delimiters.AfterSegment)
	}
	got.Delimiters = want.Delimiters
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tree after CRLF line endings differs from the original's")
	}
	// Uneven gaps leave AfterSegment empty, since no one text follows every
	// terminator.
	uneven := parseText(t, strings.Replace(text, "BHT*", "\n\nBHT*", 1))
	if uneven.Delimiters.AfterSegment != "" {
		t.Errorf("AfterSegment with uneven gaps = %q, want \"\"", uneven.Delimiters.AfterSegment)
	}
}

// TestParseElements checks the element readers on single elements.
func TestParseElements(t *testing.T) {
	segment := func(elements ...string) inputSegment {
		return inputSegment{index: 4, id: "XYZ", elements: elements}
	}
	tests := []struct {
		name    string
		parse   func() (any, error)
		want    any
		wantErr string
	}{
		{"amount", func() (any, error) { return parseAmount(segment("10.00"), 1) }, Amount{Value: Decimal{Coefficient: 10}, Text: "10.00", Set: true}, ""},
		{"negative amount", func() (any, error) { return parseAmount(segment("-0.50"), 1) }, Amount{Value: Decimal{Coefficient: -5, Scale: 1}, Text: "-0.50", Set: true}, ""},
		{"amount not sent", func() (any, error) { return parseAmount(segment(""), 1) }, Amount{}, ""},
		{"amount with a comma", func() (any, error) { return parseAmount(segment("1,000"), 1) }, Amount{}, "x12837p: parsing XYZ01 in segment 4: "},
		{"leap day", func() (any, error) { return parseDate(segment("20240229"), 1) }, date(2024, 2, 29), ""},
		{"not a leap year", func() (any, error) { return parseDate(segment("20250229"), 1) }, Date{}, "not a calendar date"},
		{"interchange date", func() (any, error) { return parseInterchangeDate(segment("260915"), 1) }, date(2026, 9, 15), ""},
		{"time with seconds", func() (any, error) { return parseTime(segment("235959"), 1) }, Time{Hour: 23, Minute: 59, Second: 59, HasSeconds: true, Set: true}, ""},
		{"time out of range", func() (any, error) { return parseTime(segment("2460"), 1) }, Time{}, "not a clock time"},
		{"pointers", func() (any, error) {
			pointers, _, err := parsePointers(segment("1:3::12"), 1, ':')
			return pointers, err
		}, [4]int{1, 3, 0, 12}, ""},
		{"pointers past the fourth", func() (any, error) {
			_, extra, err := parsePointers(segment("1:2:3:4:5"), 1, ':')
			return extra, err
		}, []string{"5"}, ""},
		{"pointer not a number", func() (any, error) {
			pointers, _, err := parsePointers(segment("A"), 1, ':')
			return pointers, err
		}, [4]int{}, "diagnosis pointer is not a whole number from 1"},
		{"DTP range kept raw", func() (any, error) {
			return parseDTP(inputSegment{index: 4, id: "DTP", elements: []string{"314", "RD8", "20260101-20260131"}})
		}, DateDTP{Index: 4, QualifierDTP01: DateTimeQualifierDTP01314, FormatDTP02: "RD8", RawDTP03: "20260101-20260131"}, ""},
		{"period D8", func() (any, error) {
			first, through, err := periodText(segment(), "HI01-4", "D8", "20240229")
			return []Date{first, through}, err
		}, []Date{date(2024, 2, 29), {}}, ""},
		{"period RD8", func() (any, error) {
			first, through, err := periodText(segment(), "HI01-4", "RD8", "20260101-20260131")
			return []Date{first, through}, err
		}, []Date{date(2026, 1, 1), date(2026, 1, 31)}, ""},
		{"period not sent", func() (any, error) {
			first, through, err := periodText(segment(), "HI01-4", "", "")
			return []Date{first, through}, err
		}, []Date{{}, {}}, ""},
		{"RD8 with an empty second date", func() (any, error) {
			_, _, err := periodText(segment(), "HI01-4", "RD8", "20260101-")
			return nil, err
		}, nil, "x12837p: parsing HI01-4 in segment 4: date is not CCYYMMDD"},
		{"date with no format qualifier", func() (any, error) {
			_, _, err := periodText(segment(), "HI01-4", "", "20260101")
			return nil, err
		}, nil, "x12837p: parsing HI01-4 in segment 4: its date format qualifier is not D8 or RD8"},
		{"decimal spelling", func() (any, error) { return Decimal{Coefficient: 10463, Scale: 2}.String(), nil }, "104.63", ""},
		{"decimal below one", func() (any, error) { return Decimal{Coefficient: -5, Scale: 2}.String(), nil }, "-0.05", ""},
		{"decimal whole", func() (any, error) { return Decimal{Coefficient: 80}.String(), nil }, "80", ""},
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
