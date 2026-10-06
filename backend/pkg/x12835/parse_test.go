package x12835

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"
)

const repoRoot = "../../.."

// TestParseFiles parses every example 835 the 835 viewer draws and every
// 835 under testdata/taxonomy, and checks that each segment of the input
// lands in the tree, as a typed segment or a RawSegment.
func TestParseFiles(t *testing.T) {
	examples, err := filepath.Glob(filepath.Join(repoRoot, "frontend/src/data/viewer835/files/*.835"))
	if err != nil {
		t.Fatal(err)
	}
	taxonomy, err := filepath.Glob(filepath.Join(repoRoot, "testdata/taxonomy/*/*.835"))
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) == 0 || len(taxonomy) == 0 {
		t.Fatalf("found %d viewer example files and %d taxonomy files, want some of each", len(examples), len(taxonomy))
	}
	for _, path := range append(examples, taxonomy...) {
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
	n := 2 + len(tr.Raw) // ST, SE
	n += present(tr.Payment.Set) + present(tr.Trace.Set) + present(tr.Currency.Set)
	n += len(tr.HeaderRefs)
	n += present(tr.ProductionDateDTM405.DateQualifierDTM01 != DateQualifierDTM01NotSent)
	for _, party := range []PartyN1{tr.PayerN1, tr.PayeeN1} {
		if party.EntityIDN101 == EntityIDN101NotSent {
			continue
		}
		n += 1 + len(party.References) + len(party.Contacts) + len(party.Raw)
		n += present(party.Address.Set) + present(party.Location.Set) + present(party.DeliveryRDM.Set)
	}
	for _, lx := range tr.HeaderNumbers {
		n += 1 + len(lx.Raw)
		n += present(lx.ProviderSummary.Set) + present(lx.InstitutionalSummary.Set)
		for _, c := range lx.Claims {
			n += countClaim(c)
		}
	}
	n += len(tr.ProviderAdjustments)
	return n
}

func countClaim(c ClaimCLP) int {
	n := 1 + len(c.Adjustments) + len(c.References) + len(c.Amounts) + len(c.Quantities) + len(c.Raw)
	for _, name := range []NameNM1{c.PatientNM1QC, c.CorrectedPatientNM174, c.RenderingNM182, c.CrossoverNM1TT, c.CorrectedPayerNM1PR} {
		n += present(name.EntityIDNM101 != EntityIDNM101NotSent)
	}
	n += present(c.InpatientMIA.Set) + present(c.OutpatientMOA.Set)
	for _, d := range []DateTimeDTM{c.StatementFromDTM232, c.StatementToDTM233, c.ReceivedDTM050, c.CoverageExpirationDTM036} {
		n += present(d.DateQualifierDTM01 != DateQualifierDTM01NotSent)
	}
	n += present(c.ContactPER.FunctionPER01 != ContactFunctionPER01NotSent)
	for _, l := range c.ServiceLines {
		n += 1 + len(l.Adjustments) + len(l.References) + len(l.Amounts) + len(l.Quantities) + len(l.Remarks) + len(l.Raw)
		for _, d := range []DateTimeDTM{l.ServiceDateDTM472, l.ServiceDateDTM150, l.ServiceDateDTM151} {
			n += present(d.DateQualifierDTM01 != DateQualifierDTM01NotSent)
		}
		for _, ref := range []ReferenceREF{l.LineControlREF6R, l.ServiceIDREFLU, l.RenderingREFHPI, l.PolicyREF0K} {
			n += present(ref.Value != "")
		}
	}
	return n
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

func amount(t *testing.T, text string) Amount {
	t.Helper()
	value, err := ParseDecimal([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	return Amount{Value: value, Set: true, Text: text}
}

const twoClaims = "frontend/src/data/viewer835/files/two-claims.835"

// TestParseTwoClaims spot-checks typed values in two-claims.835.
func TestParseTwoClaims(t *testing.T) {
	isa := parseText(t, readFile(t, twoClaims))
	tr := isa.FunctionalGroups[0].Transactions[0]
	claims := tr.HeaderNumbers[0].Claims
	if len(claims) != 2 {
		t.Fatalf("claims = %d, want 2", len(claims))
	}
	line := claims[1].ServiceLines[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ISA09", isa.DateISA09, Date{Year: 2026, Month: 10, Day: 30, Set: true}},
		{"ISA10", isa.TimeISA10, Time{Hour: 9, Minute: 5, Set: true}},
		{"GS04", isa.FunctionalGroups[0].DateGS04, Date{Year: 2026, Month: 10, Day: 30, Set: true}},
		{"BPR02", tr.Payment.TotalPaidBPR02, amount(t, "229.32")},
		{"BPR04", tr.Payment.PayMethodBPR04, PayMethodBPR04ACH},
		{"TRN02", tr.Trace.TraceNumberTRN02, "0012873345"},
		{"DTM*405 qualifier", tr.ProductionDateDTM405.DateQualifierDTM01, DateQualifierDTM01405},
		{"DTM*405 date", tr.ProductionDateDTM405.Date, Date{Year: 2026, Month: 10, Day: 28, Set: true}},
		{"payer N102", tr.PayerN1.NameN102, "MEDICARE"},
		{"payer PER01s", []ContactFunctionPER01{tr.PayerN1.Contacts[0].FunctionPER01, tr.PayerN1.Contacts[1].FunctionPER01}, []ContactFunctionPER01{ContactFunctionPER01CX, ContactFunctionPER01BL}},
		{"payee REF count", len(tr.PayeeN1.References), 1},
		{"payee REF01", tr.PayeeN1.References[0].RefQualifierREF01, RefQualifierREF01TJ},
		{"payee REF02", tr.PayeeN1.References[0].Value, "987654321"},
		{"first CLP02", claims[0].StatusCLP02, ClaimStatusCLP021},
		{"first CLP04", claims[0].PaymentCLP04, amount(t, "94.08")},
		{"first MOA03", claims[0].OutpatientMOA.RemarkMOA03, rarc.Code("MA01")},
		{"first NM109", claims[0].PatientNM1QC.IDNM109, "2AG6PM4JX59"},
		{"second claim lines", len(claims[1].ServiceLines), 3},
		{"second claim, first line SVC01", line.ProcedureSVC01, ProcedureSVC01{Qualifier: ProductQualifierSVC01HC, Code: "99214"}},
		{"second claim, first line CAS01", line.Adjustments[0].GroupCAS01, AdjGroupCAS01CO},
		{"second claim, first line CO 45", line.Adjustments[0].Reasons[0], CASReason{Reason: carc.Code("45"), Amount: amount(t, "30.00")}},
		{"second claim, first line CO 253", line.Adjustments[0].Reasons[1], CASReason{Reason: carc.Code("253"), Amount: amount(t, "2.40")}},
		{"second claim, first line DTM*472", line.ServiceDateDTM472.Date, Date{Year: 2026, Month: 10, Day: 3, Set: true}},
		{"second claim, first line REF*6R", line.LineControlREF6R.Value, "L000124"},
		{"second claim, first line AMT count", len(line.Amounts), 1},
		{"second claim, first line AMT01", line.Amounts[0].AmountQualifierAMT01, AmountQualifierAMT01B6},
		{"second claim, first line AMT02", line.Amounts[0].Amount, amount(t, "150.00")},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
}

// TestParseAmountSentOrNot checks that an amount sent as 0.00 is Set and an
// amount not sent is not.
func TestParseAmountSentOrNot(t *testing.T) {
	isa := parseText(t, readFile(t, "frontend/src/data/viewer835/files/no-payment-zero-total.835"))
	tr := isa.FunctionalGroups[0].Transactions[0]
	claim := tr.HeaderNumbers[0].Claims[0]
	cas := parseText(t, readFile(t, twoClaims)).FunctionalGroups[0].Transactions[0].HeaderNumbers[0].Claims[0].ServiceLines[0].Adjustments[0]
	tests := []struct {
		name string
		got  Amount
		want Amount
	}{
		{"BPR02 sent as 0.00", tr.Payment.TotalPaidBPR02, Amount{Value: Decimal{}, Set: true, Text: "0.00"}},
		{"CLP04 sent as 0.00", claim.PaymentCLP04, Amount{Value: Decimal{}, Set: true, Text: "0.00"}},
		{"CLP12 not sent", claim.DRGWeightCLP12, Amount{}},
		{"CAS04 empty between sent elements", cas.Reasons[0].Quantity, Amount{}},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
}

// TestParsePLB checks the provider adjustments in
// cost-report-settlements-and-passthroughs.835.
func TestParsePLB(t *testing.T) {
	isa := parseText(t, readFile(t, "frontend/src/data/viewer835/files/cost-report-settlements-and-passthroughs.835"))
	tr := isa.FunctionalGroups[0].Transactions[0]
	if len(tr.ProviderAdjustments) != 3 {
		t.Fatalf("PLBs = %d, want 3", len(tr.ProviderAdjustments))
	}
	first := tr.ProviderAdjustments[0]
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"PLB01", first.ProviderPLB01, "1234567893"},
		{"PLB02", first.FiscalPeriodPLB02, Date{Year: 2026, Month: 10, Day: 1, Set: true}},
		{"first PLB pairs", len(first.Adjustments), 6},
		{"PLB03-04", first.Adjustments[0], PLBAdjustment{PLBReasonPLB03DirectMedicalEducationPassthrough, "ADJ000041", amount(t, "-50.00")}},
		{"PLB13-14", first.Adjustments[5], PLBAdjustment{PLBReasonPLB03OrganAcquisition, "ADJ000046", amount(t, "3.00")}},
		{"second PLB, PLB05-06", tr.ProviderAdjustments[1].Adjustments[1], PLBAdjustment{PLBReasonPLB03ReturnOnEquity, "ADJ000048", amount(t, "-2.00")}},
		{"third PLB", tr.ProviderAdjustments[2].Adjustments, []PLBAdjustment{{PLBReasonPLB03BadDebt, "ADJ000053", amount(t, "2.00")}}},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(tt.got, tt.want) {
			t.Errorf("%s = %+v, want %+v", tt.name, tt.got, tt.want)
		}
	}
}

// TestParseOffList changes one coded value in two-claims.835 to a code off
// its list and checks that the field holds the text as sent and reports it
// as not known: nothing received is dropped or rewritten.
func TestParseOffList(t *testing.T) {
	firstClaim := func(isa InterchangeISA) ClaimCLP {
		return isa.FunctionalGroups[0].Transactions[0].HeaderNumbers[0].Claims[0]
	}
	tests := []struct {
		name      string
		old, new  string
		get       func(InterchangeISA) (any, bool)
		wantField any
	}{
		{
			name: "CLP02", old: "CLP*PCN2291*1*", new: "CLP*PCN2291*99*",
			get: func(isa InterchangeISA) (any, bool) {
				c := firstClaim(isa)
				return c.StatusCLP02, c.StatusCLP02.Known()
			},
			wantField: ClaimStatusCLP02("99"),
		},
		{
			name: "NM108", old: "JOHN****MI*", new: "JOHN****QQ*",
			get: func(isa InterchangeISA) (any, bool) {
				c := firstClaim(isa)
				return c.PatientNM1QC.IDQualifierNM108, c.PatientNM1QC.IDQualifierNM108.Known()
			},
			wantField: IDQualifierNM108("QQ"),
		},
		{
			name: "SVC01 qualifier component", old: "SVC*HC:99213*", new: "SVC*ZZ:99213*",
			get: func(isa InterchangeISA) (any, bool) {
				l := firstClaim(isa).ServiceLines[0]
				return l.ProcedureSVC01, l.ProcedureSVC01.Qualifier.Known()
			},
			wantField: ProcedureSVC01{Qualifier: ProductQualifierSVC01("ZZ"), Code: "99213"},
		},
		{
			name: "CAS01 on a service line", old: "CAS*PR*2*24.00~", new: "CAS*XX*2*24.00~",
			get: func(isa InterchangeISA) (any, bool) {
				l := firstClaim(isa).ServiceLines[0]
				return l.Adjustments[1].GroupCAS01, l.Adjustments[1].GroupCAS01.Known()
			},
			wantField: AdjGroupCAS01("XX"),
		},
		{
			name: "N103 on the payee", old: "N1*PE*VALLEY FAMILY MED*XX*", new: "N1*PE*VALLEY FAMILY MED*ZZ*",
			get: func(isa InterchangeISA) (any, bool) {
				p := isa.FunctionalGroups[0].Transactions[0].PayeeN1
				return p.IDQualifierN103, p.IDQualifierN103.Known()
			},
			wantField: IDQualifierN103("ZZ"),
		},
		{
			name: "ISA15", old: "*0*T*:~", new: "*0*X*:~",
			get: func(isa InterchangeISA) (any, bool) {
				return isa.UsageISA15, isa.UsageISA15.Known()
			},
			wantField: UsageISA15("X"),
		},
	}
	base := readFile(t, twoClaims)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("two-claims.835 does not contain %q", tt.old)
			}
			text := strings.Replace(base, tt.old, tt.new, 1)
			isa := parseText(t, text)
			field, known := tt.get(isa)
			if !reflect.DeepEqual(field, tt.wantField) {
				t.Errorf("field = %+v, want %+v", field, tt.wantField)
			}
			if known {
				t.Errorf("Known() = true, want false")
			}
			var out strings.Builder
			if err := Write(&out, isa); err != nil {
				t.Fatal(err)
			}
			if out.String() != text {
				t.Errorf("the file with the off-list value does not write back out as it came")
			}
		})
	}
}

// TestParseRaw inserts segments the model has no slot for into
// two-claims.835 and checks that each is kept as a RawSegment on the struct
// that encloses it, with its file position.
func TestParseRaw(t *testing.T) {
	transaction := func(isa InterchangeISA) TransactionST {
		return isa.FunctionalGroups[0].Transactions[0]
	}
	tests := []struct {
		name     string
		old, new string
		se       string // the SE that replaces SE*49* when the insert is inside the transaction
		get      func(InterchangeISA) []RawSegment
		want     []RawSegment
	}{
		{
			name: "unmodeled segment in a claim", old: "MOA***MA01~", new: "MOA***MA01~K3*ABC~", se: "SE*50*",
			get:  func(isa InterchangeISA) []RawSegment { return transaction(isa).HeaderNumbers[0].Claims[0].Raw },
			want: []RawSegment{{Index: 21, ID: "K3", Elements: []string{"ABC"}}},
		},
		{
			name: "second MOA in a claim", old: "MOA***MA01~", new: "MOA***MA01~MOA***MA02~", se: "SE*50*",
			get:  func(isa InterchangeISA) []RawSegment { return transaction(isa).HeaderNumbers[0].Claims[0].Raw },
			want: []RawSegment{{Index: 21, ID: "MOA", Elements: []string{"", "", "MA02"}}},
		},
		{
			name: "DTM*150 at the claim level", old: "DTM*050*20261026~", new: "DTM*050*20261026~DTM*150*20261001~", se: "SE*50*",
			get:  func(isa InterchangeISA) []RawSegment { return transaction(isa).HeaderNumbers[0].Claims[0].Raw },
			want: []RawSegment{{Index: 22, ID: "DTM", Elements: []string{"150", "20261001"}}},
		},
		{
			name: "second payee loop", old: "LX*1~", new: "N1*PE*OTHER~N3*1 MAIN~LX*1~", se: "SE*51*",
			get: func(isa InterchangeISA) []RawSegment { return transaction(isa).Raw },
			want: []RawSegment{
				{Index: 16, ID: "N1", Elements: []string{"PE", "OTHER"}},
				{Index: 17, ID: "N3", Elements: []string{"1 MAIN"}},
			},
		},
		{
			name: "unmodeled segment in a group", old: "ST*835*0001~", new: "XYZ*1~ST*835*0001~",
			get:  func(isa InterchangeISA) []RawSegment { return isa.FunctionalGroups[0].Raw },
			want: []RawSegment{{Index: 2, ID: "XYZ", Elements: []string{"1"}}},
		},
	}
	base := readFile(t, twoClaims)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("two-claims.835 does not contain %q", tt.old)
			}
			text := strings.Replace(base, tt.old, tt.new, 1)
			if tt.se != "" {
				text = strings.Replace(text, "SE*49*", tt.se, 1)
			}
			isa := parseText(t, text)
			if got := tt.get(isa); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Raw = %+v, want %+v", got, tt.want)
			}
			if got, want := countTree(isa), countInput([]byte(text)); got != want {
				t.Errorf("countTree = %d, want %d", got, want)
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
		{"SE01", "SE*49*0001", "SE*48*0001", `x12835: checking SE01 of transaction "0001": SE01 is 48, ST through SE is 49 segments`},
		{"SE02", "SE*49*0001", "SE*49*0002", `x12835: checking SE02: ST02 is "0001", SE02 is "0002"`},
		{"GE01", "GE*1*1~", "GE*2*1~", `x12835: checking GE01 of group "1": GE01 is 2, the group holds 1 transactions`},
		{"GE02", "GE*1*1~", "GE*1*7~", `x12835: checking GE02: GS06 is "1", GE02 is "7"`},
		{"IEA01", "IEA*1*", "IEA*2*", `x12835: checking IEA01: IEA01 is 2, the interchange holds 1 functional groups`},
		{"IEA02", "IEA*1*000000905", "IEA*1*000000906", `x12835: checking IEA02: ISA13 is "000000905", IEA02 is "000000906"`},
		{"malformed amount", "CLP*PCN2291*1*150.00", "CLP*PCN2291*1*15O.00", "x12835: parsing CLP03 in segment 17: "},
		{"not a calendar date", "DTM*472*20261001", "DTM*472*20261301", "x12835: parsing DTM02 in segment 23: not a calendar date"},
		{"element past the segment", "LX*1~", "LX*1*2~", "x12835: parsing LX02 in segment 16: the segment ends at LX01"},
		{"component past the composite", "SVC*HC:99213*", "SVC*HC:99213:A:B:C:D:TEXT:X:Y*", "x12835: parsing SVC01 in segment 22: the composite ends at component 8"},
		{"LX01 not a number", "LX*1~", "LX*A~", "x12835: parsing LX01 in segment 16: not a whole number"},
		{"missing SE", "SE*49*0001~", "", "x12835: parsing segment 50: GE inside the transaction at segment 2 (missing SE)"},
		{"input after IEA", "IEA*1*000000905~", "IEA*1*000000905~TA1*1~", "x12835: parsing segment 53: input continues after the IEA"},
		{"missing terminator", "IEA*1*000000905~", "IEA*1*000000905", "x12835: splitting segment 52: missing segment terminator"},
	}
	base := readFile(t, twoClaims)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(base, tt.old) {
				t.Fatalf("two-claims.835 does not contain %q", tt.old)
			}
			_, err := Parse([]byte(strings.Replace(base, tt.old, tt.new, 1)))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Parse error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

// TestParseLineEndings checks that whitespace between segments, and a
// carriage return before a newline terminator, do not change the tree.
func TestParseLineEndings(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		old, new string
	}{
		{"tilde then CRLF", twoClaims, "~", "~\r\n"},
		{"newline terminator then CRLF", "testdata/taxonomy/T5.43/vendor.835", "\n", "\r\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := readFile(t, tt.path)
			want := parseText(t, text)
			got := parseText(t, strings.ReplaceAll(text, tt.old, tt.new))
			// The delimiters record the whitespace, and they are not the tree.
			want.Delimiters, got.Delimiters = Delimiters{}, Delimiters{}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("tree after replacing %q with %q differs from the original's", tt.old, tt.new)
			}
		})
	}
}

// TestParseElements checks the amount, date, time, and count readers on
// single elements.
func TestParseElements(t *testing.T) {
	segment := func(text string) inputSegment {
		return inputSegment{index: 4, id: "XYZ", elements: []string{text}}
	}
	tests := []struct {
		text    string
		parse   func(inputSegment) (any, error)
		want    any
		wantErr string
	}{
		{"10.00", func(s inputSegment) (any, error) { return parseAmount(s, 1) }, Amount{Value: Decimal{Coefficient: 10}, Set: true, Text: "10.00"}, ""},
		{"-0.50", func(s inputSegment) (any, error) { return parseAmount(s, 1) }, Amount{Value: Decimal{Coefficient: -5, Scale: 1}, Set: true, Text: "-0.50"}, ""},
		{"", func(s inputSegment) (any, error) { return parseAmount(s, 1) }, Amount{}, ""},
		{"1,000", func(s inputSegment) (any, error) { return parseAmount(s, 1) }, Amount{}, "x12835: parsing XYZ01 in segment 4: "},
		{"20240229", func(s inputSegment) (any, error) { return parseDate(s, 1) }, Date{Year: 2024, Month: 2, Day: 29, Set: true}, ""},
		{"20250229", func(s inputSegment) (any, error) { return parseDate(s, 1) }, Date{}, "not a calendar date"},
		{"261030", func(s inputSegment) (any, error) { return parseDate(s, 1) }, Date{}, "date is not CCYYMMDD"},
		{"", func(s inputSegment) (any, error) { return parseDate(s, 1) }, Date{}, ""},
		{"261030", func(s inputSegment) (any, error) { return parseInterchangeDate(s, 1) }, Date{Year: 2026, Month: 10, Day: 30, Set: true}, ""},
		{"0905", func(s inputSegment) (any, error) { return parseTime(s, 1) }, Time{Hour: 9, Minute: 5, Set: true}, ""},
		{"235959", func(s inputSegment) (any, error) { return parseTime(s, 1) }, Time{Hour: 23, Minute: 59, Second: 59, HasSeconds: true, Set: true}, ""},
		{"2460", func(s inputSegment) (any, error) { return parseTime(s, 1) }, Time{}, "not a clock time"},
		{"09051234", func(s inputSegment) (any, error) { return parseTime(s, 1) }, Time{}, "time is not HHMM or HHMMSS"},
		{"049", func(s inputSegment) (any, error) { return parseWholeNumber(s, 1) }, 49, ""},
		{"-1", func(s inputSegment) (any, error) { return parseWholeNumber(s, 1) }, 0, "not a whole number"},
	}
	for _, tt := range tests {
		got, err := tt.parse(segment(tt.text))
		if tt.wantErr != "" {
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("parse(%q) error = %v, want it to contain %q", tt.text, err, tt.wantErr)
			}
			continue
		}
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parse(%q) = %+v, %v, want %+v, nil", tt.text, got, err, tt.want)
		}
	}
}
