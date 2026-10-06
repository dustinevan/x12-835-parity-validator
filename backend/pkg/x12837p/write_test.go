package x12837p

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestRoundTrip parses every 837P file under testdata/ and
// frontend/src/data and writes it back, and requires the bytes to match
// the input exactly: every element, every component, every raw segment,
// the ISA's padding, and the file's own line breaks.
func TestRoundTrip(t *testing.T) {
	paths := professionalFiles(t)
	if len(paths) < 4 {
		t.Fatalf("found %d 837P files, want the three in testdata/x12837p, the taxonomy's, and the docs site's", len(paths))
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
			got, err := Bytes(isa)
			if err != nil {
				t.Fatalf("Bytes = %v, want nil", err)
			}
			if !bytes.Equal(got, data) {
				t.Errorf("round trip differs:\n%s", firstDifference(got, data))
			}
		})
	}
}

// firstDifference describes where two byte strings first differ: the
// offset, and the segment of each around it.
func firstDifference(got, want []byte) string {
	n := len(got)
	if len(want) < n {
		n = len(want)
	}
	at := n
	for i := 0; i < n; i++ {
		if got[i] != want[i] {
			at = i
			break
		}
	}
	window := func(b []byte) string {
		start := at - 60
		if start < 0 {
			start = 0
		}
		end := at + 60
		if end > len(b) {
			end = len(b)
		}
		return strings.ReplaceAll(string(b[start:end]), "\n", "\\n")
	}
	return "at byte " + itoa(at) + " (got " + itoa(len(got)) + " bytes, want " + itoa(len(want)) + ")\n got: " + window(got) + "\nwant: " + window(want)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

// TestWriteExtraKept changes office-visit.837 so that segments carry text
// past the model's last position, and checks that the text lands in the
// typed fields and comes back out on the round trip.
func TestWriteExtraKept(t *testing.T) {
	base := readFile(t, officeVisit)
	tests := []struct {
		name     string
		old, new string
		get      func(ClaimCLM) []string
	}{
		{"element after LX01", "LX*1~", "LX*1*2~", func(c ClaimCLM) []string { return c.ServiceLines[0].Extra }},
		{"component after SV101-7", "SV1*HC:87880*", "SV1*HC:87880:A:B:C:D:TEXT:X*", func(c ClaimCLM) []string { return c.ServiceLines[1].Service.ProcedureSV101.Extra }},
		{"component after HI01-9", "HI*ABK:J069*", "HI*ABK:J069::::::::X*", func(c ClaimCLM) []string { return c.DiagnosisCodesHIABK.Codes[0].Extra }},
		{"element after NM109", "NM1*82*1*SYNTHETIC*JORDAN****XX*1999999992~", "NM1*82*1*SYNTHETIC*JORDAN****XX*1999999992*34*Q~", func(c ClaimCLM) []string { return c.RenderingProviderNM182.Name.Extra }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("office-visit.837 does not contain %q", tt.old)
			}
			text := strings.Replace(base, tt.old, tt.new, 1)
			isa := parseText(t, text)
			claim := firstTransaction(isa).BillingProviders[0].Subscribers[0].Claims[0]
			if got := tt.get(claim); len(got) == 0 {
				t.Errorf("extra text was not kept")
			}
			got, err := Bytes(isa)
			if err != nil {
				t.Fatalf("Bytes = %v, want nil", err)
			}
			if string(got) != text {
				t.Errorf("round trip differs:\n%s", firstDifference(got, []byte(text)))
			}
		})
	}
}

// TestWriteHandBuilt writes a small interchange built in code, with no
// indexes, and checks that it comes out in the implementation guide's
// order with the canonical spellings of its amounts, dates, and times.
func TestWriteHandBuilt(t *testing.T) {
	isa := InterchangeISA{
		Delimiters:         Delimiters{Element: "*", Segment: "~", AfterSegment: "\n"},
		AuthQualifierISA01: AuthQualifierISA0100, AuthInfoISA02: "          ",
		SecQualifierISA03: SecQualifierISA0300, SecInfoISA04: "          ",
		SenderIDQualifierISA05: InterchangeIDQualifierISA05ZZ, SenderIDISA06: "SENDER         ",
		ReceiverIDQualifierISA07: InterchangeIDQualifierISA05ZZ, ReceiverIDISA08: "RECEIVER       ",
		DateISA09: date(2026, 9, 15), TimeISA10: Time{Hour: 10, Minute: 30, Set: true},
		RepetitionSepISA11: "^", VersionISA12: "00501", ControlNumberISA13: "000000001",
		AckRequestedISA14: AckRequestedISA140, UsageISA15: UsageISA15T, ComponentSepISA16: ":",
		GroupCountIEA01: 1, ControlNumberIEA02: "000000001",
		FunctionalGroups: []FunctionalGroupGS{{
			FunctionalIDGS01: FunctionalIDGS01HC, SenderCodeGS02: "SENDER", ReceiverCodeGS03: "RECEIVER",
			DateGS04: date(2026, 9, 15), TimeGS05: Time{Hour: 10, Minute: 30, Set: true},
			ControlNumberGS06: "1", AgencyGS07: AgencyGS07X, VersionGS08: "005010X222A1",
			TransactionCountGE01: 1, ControlNumberGE02: "1",
			Transactions: []TransactionST{{
				TransactionSetIDST01: TransactionSetIDST01837, ControlNumberST02: "0001", ImplementationGuideVersionST03: "005010X222A1",
				Header: HeaderBHT{Set: true, StructureBHT01: HierarchicalStructureBHT010019, PurposeBHT02: TransactionSetPurposeBHT0200,
					OriginatorTransactionIDBHT03: "BATCH1", DateBHT04: date(2026, 9, 15), TimeBHT05: Time{Hour: 10, Minute: 30, Set: true}, TransactionTypeBHT06: TransactionTypeBHT06CH},
				BillingProviders: []BillingProviderHL{{
					HL: HierarchyHL{IDNumberHL01: "1", LevelHL03: HierarchicalLevelHL0320, ChildCodeHL04: HierarchicalChildCodeHL041},
					Subscribers: []SubscriberHL{{
						HL: HierarchyHL{IDNumberHL01: "2", ParentIDNumberHL02: "1", LevelHL03: HierarchicalLevelHL0322, ChildCodeHL04: HierarchicalChildCodeHL040},
						Claims: []ClaimCLM{{
							PatientControlCLM01: "PCN1", TotalChargeCLM02: Amount{Value: Decimal{Coefficient: 1755, Scale: 1}, Set: true},
							PlaceOfServiceCLM05_1: "11", FacilityQualifierCLM05_2: FacilityCodeQualifierCLM05_2B, FrequencyCLM05_3: ClaimFrequencyCLM05_31,
							SignatureCLM06:          ProviderOrSupplierSignatureCLM06Y,
							PriorAuthorizationREFG1: ReferenceREF{QualifierREF01: ReferenceIdentificationQualifierREF01G1, Value: "AUTH1"},
							OnsetDTP431:             DateDTP{QualifierDTP01: DateTimeQualifierDTP01431, FormatDTP02: "D8", RawDTP03: "20260901"},
							ServiceLines: []ServiceLineLX{{
								AssignedNumberLX01: "1",
								Service: ProfessionalServiceSV1{Set: true, ProcedureSV101: ProcedureSV101{Qualifier: "HC", Code: "99213"},
									ChargeSV102: Amount{Value: Decimal{Coefficient: 1755, Scale: 1}, Set: true}, UnitSV103: UnitOrBasisForMeasurementCodeUN,
									UnitCountSV104: Amount{Value: Decimal{Coefficient: 1}, Set: true}, DiagnosisPointersSV107: [4]int{1}},
							}},
						}},
					}},
				}},
				SegmentCountSE01: 10, ControlNumberSE02: "0001",
			}},
		}},
	}
	want := strings.Join([]string{
		"ISA*00*          *00*          *ZZ*SENDER         *ZZ*RECEIVER       *260915*1030*^*00501*000000001*0*T*:~",
		"GS*HC*SENDER*RECEIVER*20260915*1030*1*X*005010X222A1~",
		"ST*837*0001*005010X222A1~",
		"BHT*0019*00*BATCH1*20260915*1030*CH~",
		"HL*1**20*1~",
		"HL*2*1*22*0~",
		"CLM*PCN1*175.5***11:B:1*Y~",
		"DTP*431*D8*20260901~",
		"REF*G1*AUTH1~",
		"LX*1~",
		"SV1*HC:99213*175.5*UN*1***1~",
		"SE*10*0001~",
		"GE*1*1~",
		"IEA*1*000000001~",
		"",
	}, "\n")
	got, err := Bytes(isa)
	if err != nil {
		t.Fatalf("Bytes = %v, want nil", err)
	}
	if string(got) != want {
		t.Errorf("Bytes =\n%s\nwant\n%s", got, want)
	}
	if _, err := Parse(got); err != nil {
		t.Errorf("Parse of the written text = %v, want nil", err)
	}
}

// TestWriteDelimitersRequired checks that Write refuses an interchange
// whose delimiters are not set rather than guessing them.
func TestWriteDelimitersRequired(t *testing.T) {
	_, err := Bytes(InterchangeISA{ComponentSepISA16: ":"})
	if err == nil || !strings.Contains(err.Error(), "element separator") {
		t.Errorf("Bytes with no delimiters = %v, want an error naming the element separator", err)
	}
	_, err = Bytes(InterchangeISA{Delimiters: Delimiters{Element: "*", Segment: "~"}})
	if err == nil || !strings.Contains(err.Error(), "ISA16") {
		t.Errorf("Bytes with no ISA16 = %v, want an error naming ISA16", err)
	}
}
