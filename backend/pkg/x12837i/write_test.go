package x12837i

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRoundTrip parses every 837I on file and writes it back, and requires
// the bytes to come out exactly as they went in: the three fixtures under
// testdata/x12837i, every taxonomy fixture's source.837 that is an 837I,
// and every example file the docs site's 837I viewer draws. A difference
// is reported at its first byte, with the segment around it.
func TestRoundTrip(t *testing.T) {
	var files []string
	for _, pattern := range []string{"testdata/x12837i/*.837", "frontend/src/data/viewer837i/files/*.837"} {
		matches, err := filepath.Glob(filepath.Join(repoRoot, pattern))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, matches...)
	}
	sources, err := filepath.Glob(filepath.Join(repoRoot, "testdata/taxonomy/*/source.837"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range sources {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte("*005010X223A2~")) {
			files = append(files, path)
		}
	}
	if len(files) < 3+29+42 {
		t.Fatalf("found %d files, want the 3 fixtures, 29 viewer examples, and 42 taxonomy 837Is", len(files))
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
			var out bytes.Buffer
			if err := Write(&out, isa); err != nil {
				t.Fatalf("Write = %v, want nil", err)
			}
			if !bytes.Equal(out.Bytes(), data) {
				t.Error(firstDifference(data, out.Bytes()))
			}
		})
	}
}

// firstDifference describes where two files first differ: the byte offset
// and the segment around it on each side.
func firstDifference(want, got []byte) string {
	i := 0
	for i < len(want) && i < len(got) && want[i] == got[i] {
		i++
	}
	around := func(data []byte) string {
		start := bytes.LastIndexByte(data[:min(i, len(data))], '~') + 1
		end := bytes.IndexByte(data[min(i, len(data)):], '~')
		if end < 0 {
			end = len(data)
		} else {
			end += min(i, len(data)) + 1
		}
		return string(data[start:end])
	}
	return fmt.Sprintf("written bytes differ from the input at byte %d: input %q, written %q (%d and %d bytes)", i, around(want), around(got), len(want), len(got))
}

// TestWritePositions checks that Write refuses a tree whose segment
// positions have a gap or a repeat, since either would mean a segment was
// lost or written twice.
func TestWritePositions(t *testing.T) {
	base := parseText(t, readFile(t, inpatientFile))
	tests := []struct {
		name  string
		index int
		want  string
	}{
		{"gap", 99, "x12837i: writing: no segment at position 3"},
		{"repeat", 2, "x12837i: writing: two segments at position 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isa := base
			isa.FunctionalGroups = append([]FunctionalGroupGS(nil), base.FunctionalGroups...)
			isa.FunctionalGroups[0].Transactions = append([]TransactionST(nil), base.FunctionalGroups[0].Transactions...)
			isa.FunctionalGroups[0].Transactions[0].Header.Index = tt.index // the BHT sits at position 3
			err := Write(&bytes.Buffer{}, isa)
			if err == nil || err.Error() != tt.want {
				t.Errorf("Write error = %v, want %q", err, tt.want)
			}
		})
	}
}

// TestWriteElementText checks the element writers on single values: the
// four DTP03 forms, a time with and without seconds, an amount with the
// file's text and one built without it, and the ISA date.
func TestWriteElementText(t *testing.T) {
	hour := Time{Hour: 8, Minute: 15, Set: true}
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"DTP TM", dtp03Out(DateDTP{Set: true, FormatDTP02: "TM", Time: Time{Hour: 11, Minute: 30, Set: true}}), "1130"},
		{"DTP DT", dtp03Out(DateDTP{Set: true, FormatDTP02: "DT", Date: date(2026, 9, 20), Time: hour}), "202609200815"},
		{"DTP RD8", dtp03Out(DateDTP{Set: true, FormatDTP02: "RD8", Date: date(2026, 9, 20), DateThrough: date(2026, 9, 24)}), "20260920-20260924"},
		{"DTP D8", dtp03Out(DateDTP{Set: true, FormatDTP02: "D8", Date: date(2026, 9, 15)}), "20260915"},
		{"time HHMM", timeOut(Time{Hour: 9, Minute: 5, Set: true}), "0905"},
		{"time HHMMSS", timeOut(Time{Hour: 9, Minute: 5, Second: 7, WithSeconds: true, Set: true}), "090507"},
		{"time not set", timeOut(Time{}), ""},
		{"amount as sent", amountOut(Amount{Value: Decimal{Coefficient: 125, Scale: 1}, Text: "12.50", Set: true}), "12.50"},
		{"amount without text", amountOut(Amount{Value: Decimal{Coefficient: 125, Scale: 1}, Set: true}), "12.5"},
		{"amount not set", amountOut(Amount{}), ""},
		{"decimal negative fraction", Decimal{Coefficient: -5, Scale: 1}.String(), "-0.5"},
		{"decimal cents", Decimal{Coefficient: 5, Scale: 2}.String(), "0.05"},
		{"decimal whole", Decimal{Coefficient: 12500}.String(), "12500"},
		{"decimal zero", Decimal{}.String(), "0"},
		{"ISA date", interchangeDateOut(date(2026, 10, 1)), "261001"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

// TestWriteDefaults checks that a tree built by hand, with no delimiters
// recorded, is written with '*', ':', and '~', its composites joined and its
// trailing empty elements left off.
func TestWriteDefaults(t *testing.T) {
	isa := InterchangeISA{
		AuthQualifierISA01: "00", AuthInfoISA02: "          ", SecQualifierISA03: "00", SecInfoISA04: "          ",
		SenderIDQualifierISA05: "ZZ", SenderIDISA06: "SENDER         ", ReceiverIDQualifierISA07: "ZZ", ReceiverIDISA08: "RECEIVER       ",
		DateISA09: date(2026, 10, 1), TimeISA10: Time{Hour: 10, Minute: 30, Set: true}, RepetitionSepISA11: "^", VersionISA12: "00501",
		ControlNumberISA13: "000000001", AckRequestedISA14: "0", UsageISA15: "T", ComponentSepISA16: ":",
		FunctionalGroups: []FunctionalGroupGS{{
			Index: 1, FunctionalIDGS01: TransactionTypeGS01_ST01HC, SenderCodeGS02: "S", ReceiverCodeGS03: "R",
			DateGS04: date(2026, 10, 1), TimeGS05: Time{Hour: 10, Minute: 30, Set: true}, ControlNumberGS06: "1", AgencyGS07: "X", VersionGS08: "005010X223A2",
			Transactions: []TransactionST{{
				Index: 2, TransactionSetIDST01: TransactionTypeGS01_ST01837, ControlNumberST02: "0001", VersionST03: "005010X223A2",
				Raw:     []RawSegment{{Index: 3, ID: "K3", Elements: []string{"TEXT", "", ""}}},
				IndexSE: 4, SegmentCountSE01: "3", ControlNumberSE02: "0001",
			}},
			IndexGE: 5, TransactionCountGE01: "1", ControlNumberGE02: "1",
		}},
		IndexIEA: 6, GroupCountIEA01: "1", ControlNumberIEA02: "000000001",
	}
	var out bytes.Buffer
	if err := Write(&out, isa); err != nil {
		t.Fatalf("Write = %v, want nil", err)
	}
	want := "ISA*00*          *00*          *ZZ*SENDER         *ZZ*RECEIVER       *261001*1030*^*00501*000000001*0*T*:~" +
		"GS*HC*S*R*20261001*1030*1*X*005010X223A2~ST*837*0001*005010X223A2~K3*TEXT~SE*3*0001~GE*1*1~IEA*1*000000001~"
	if out.String() != want {
		t.Errorf("Write =\n%s\nwant\n%s", out.String(), want)
	}
}
