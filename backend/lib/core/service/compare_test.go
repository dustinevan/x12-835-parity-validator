package service

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/payment"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/x12835"
)

func TestMatchKey(t *testing.T) {
	base := Key{MBI: "5EN4VN7KY98", DateOfService: "20260105", TotalCharge: "447.00", TaxID: "999615452", NPI: "1999002557", Workload: "91101"}
	want := sha256.Sum256([]byte("5EN4VN7KY98\x1f20260105\x1f447\x1f999615452\x1f1999002557\x1f91101"))

	tests := []struct {
		name    string
		key     Key
		same    bool // hashes to want
		wantErr bool
	}{
		{"as written", base, true, false},
		{"charge without cents", Key{"5EN4VN7KY98", "20260105", "447", "999615452", "1999002557", "91101"}, true, false},
		{"charge with one decimal", Key{"5EN4VN7KY98", "20260105", "447.0", "999615452", "1999002557", "91101"}, true, false},
		{"MBI lower case and padded", Key{" 5en4vn7ky98 ", "20260105", "447.00", "999615452", "1999002557", "91101"}, true, false},
		{"tax ID with hyphen", Key{"5EN4VN7KY98", "20260105", "447.00", "99-9615452", "1999002557", "91101"}, true, false},
		{"different charge", Key{"5EN4VN7KY98", "20260105", "447.01", "999615452", "1999002557", "91101"}, false, false},
		{"different workload", Key{"5EN4VN7KY98", "20260105", "447.00", "999615452", "1999002557", "91102"}, false, false},
		{"invalid date", Key{"5EN4VN7KY98", "20261305", "447.00", "999615452", "1999002557", "91101"}, false, true},
		{"charge not a number", Key{"5EN4VN7KY98", "20260105", "4X7", "999615452", "1999002557", "91101"}, false, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MatchKey(test.key)
			if test.wantErr {
				if err == nil {
					t.Fatalf("MatchKey: got no error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Equal(got, want[:]) != test.same {
				t.Fatalf("MatchKey(%+v) = %x; same as base: got %v, want %v", test.key, got, !test.same, test.same)
			}
		})
	}
}

// fixtureClaim splits one taxonomy fixture file that holds one claim.
func fixtureClaim(t *testing.T, dir, name string) claim.Claim {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testdata", "taxonomy", dir, name))
	if err != nil {
		t.Fatal(err)
	}
	claims, _, err := Split(dir, name, data)
	if err != nil {
		t.Fatal(err)
	}
	if len(claims) != 1 {
		t.Fatalf("%s/%s: got %d claims, want 1", dir, name, len(claims))
	}
	if !bytes.HasPrefix(claims[0].Segments, []byte("CLP*")) || !bytes.Contains(data, claims[0].Segments) {
		t.Fatalf("%s/%s: claim bytes are not an exact slice of the file starting at CLP", dir, name)
	}
	return claims[0]
}

// TestSplitFault checks that a malformed amount fails the file and that the
// error names the file, the payment, and the claim it sits in.
func TestSplitFault(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testdata", "taxonomy", "T4.02", "mac.835"))
	if err != nil {
		t.Fatal(err)
	}
	broken := bytes.Replace(data, []byte("CLP*"), []byte("CLP*X*1*1X0.00*"), 1)
	if bytes.Equal(broken, data) {
		t.Fatal("fixture has no CLP to break")
	}
	_, _, err = Split("mac", "mac.835", broken)
	if err == nil {
		t.Fatal("Split accepted a malformed CLP03")
	}
	for _, want := range []string{"splitting mac.835", "payment 1 claim 1", "CLP03"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Split error %q does not name %q", err, want)
		}
	}
}

// TestSplitPayment checks the payment row Split builds for a transaction:
// the payer and payee, the BPR, TRN, and DTM*405 fields, and the counts and
// sums over its claims and PLB segments. T5.01's MAC side withholds a 40.00
// PLB, so its BPR02 is the claims' payments less 40; T4.02 has no PLB. For
// both, the row balances: BPR02 is the sum of the claims' CLP04 less the
// PLB total.
func TestSplitPayment(t *testing.T) {
	tests := []struct {
		dir, name string
		want      payment.Payment
	}{
		{"T5.01", "mac.835", payment.Payment{
			SubmissionPrefix: "T5.01",
			FileName:         "mac.835",
			PaymentNumber:    1,
			PayerIdentifier:  "50101",
			PayerName:        "SYNTHETIC MAC",
			PayeeIdentifier:  "1010037038",
			PayeeName:        "SYNTHETIC ORTHOPEDIC ASSOCIATES",
			PayeeTaxID:       "850110037",
			PayeeState:       "IL",
			PaymentDate:      time.Date(2026, time.March, 16, 0, 0, 0, 0, time.UTC),
			PaymentMethod:    "ACH",
			HandlingCode:     "I",
			TraceNumber:      "0001234510",
			ProductionDate:   time.Date(2026, time.March, 13, 0, 0, 0, 0, time.UTC),
			TotalPaid:        "4180.22",
			ClaimCount:       3,
			PLBCount:         1,
			PLBTotal:         "40",
		}},
		{"T4.02", "mac.835", payment.Payment{
			SubmissionPrefix: "T4.02",
			FileName:         "mac.835",
			PaymentNumber:    1,
			PayerIdentifier:  "94002",
			PayerName:        "SYNTHETIC MAC",
			PayeeIdentifier:  "1940021235",
			PayeeName:        "SYNTHETIC CLINIC 94002",
			PayeeTaxID:       "940402017",
			PayeeState:       "CO",
			PaymentDate:      time.Date(2026, time.March, 27, 0, 0, 0, 0, time.UTC),
			PaymentMethod:    "ACH",
			HandlingCode:     "I",
			TraceNumber:      "00940020001",
			ProductionDate:   time.Date(2026, time.March, 27, 0, 0, 0, 0, time.UTC),
			TotalPaid:        "99.72",
			ClaimCount:       1,
			PLBCount:         0,
			PLBTotal:         "0",
		}},
	}
	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "testdata", "taxonomy", test.dir, test.name))
			if err != nil {
				t.Fatal(err)
			}
			claims, payments, err := Split(test.dir, test.name, data)
			if err != nil {
				t.Fatal(err)
			}
			if len(payments) != 1 {
				t.Fatalf("Split: got %d payments, want 1", len(payments))
			}
			got := payments[0]
			if got != test.want {
				t.Errorf("Split payment = %+v, want %+v", got, test.want)
			}
			if len(claims) != got.ClaimCount {
				t.Errorf("Split: %d claims, payment row counts %d", len(claims), got.ClaimCount)
			}

			var paid x12835.Decimal
			for _, c := range claims {
				amount, err := x12835.ParseDecimal([]byte(c.ClaimPaymentAmount))
				if err != nil {
					t.Fatal(err)
				}
				if paid, err = paid.Add(amount); err != nil {
					t.Fatal(err)
				}
			}
			plb, err := x12835.ParseDecimal([]byte(got.PLBTotal))
			if err != nil {
				t.Fatal(err)
			}
			balance, err := paid.Sub(plb)
			if err != nil {
				t.Fatal(err)
			}
			if balance.String() != got.TotalPaid {
				t.Errorf("claims paid %s - PLB %s = %s, BPR02 is %s", paid, plb, balance, got.TotalPaid)
			}
		})
	}
}

type labeled struct {
	element, line string
}

func TestDiffClaimsFixtures(t *testing.T) {
	tests := []struct {
		dir  string
		key  Key
		want []labeled
	}{
		{
			// Lines carry REF*6R, so they pair by it, and their order is
			// not compared: no differences.
			dir: "T1.29",
		},
		{
			// The fixture also lists BPR02, which is payment-level and not
			// part of a claim comparison.
			dir: "T3.14",
			key: Key{MBI: "3TX4E17JG32", DateOfService: "20260615", TotalCharge: "140.00", TaxID: "933144417", NPI: "1314456736", Workload: "31417"},
			want: []labeled{
				{"CLP04", ""},
				{"SVC03", "314101"},
				{"CAS*CO*253", "314101"},
			},
		},
		{
			// No REF*6R: lines pair by position, so every value on both
			// reordered lines differs.
			dir: "T4.02",
			key: Key{MBI: "4HT0U20FJ74", DateOfService: "20260305", TotalCharge: "190.00", TaxID: "940402017", NPI: "1940021235", Workload: "94002"},
			want: []labeled{
				{"SVC01", ""}, {"SVC02", ""}, {"SVC03", ""}, {"CAS*CO*45", ""}, {"CAS*CO*253", ""}, {"CAS*PR*2", ""}, {"AMT*B6", ""},
				{"SVC01", ""}, {"SVC02", ""}, {"SVC03", ""}, {"CAS*CO*45", ""}, {"CAS*CO*253", ""}, {"CAS*PR*2", ""}, {"AMT*B6", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			mac := fixtureClaim(t, test.dir, "mac.835")
			vendor := fixtureClaim(t, test.dir, "vendor.835")
			if !bytes.Equal(mac.MatchKey, vendor.MatchKey) {
				t.Fatalf("match keys differ: %x, %x", mac.MatchKey, vendor.MatchKey)
			}
			if test.key != (Key{}) {
				want, err := MatchKey(test.key)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(mac.MatchKey, want) {
					t.Fatalf("match key %x; want the fixture key's %x", mac.MatchKey, want)
				}
			}
			diffs, err := DiffClaims(compareSet1, mac, vendor)
			if err != nil {
				t.Fatal(err)
			}
			var got []labeled
			for _, d := range diffs {
				got = append(got, labeled{d.Element, d.LineControlNumber})
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("differences:\n got %v\nwant %v", got, test.want)
			}
		})
	}
}

func TestDiffClaimsValues(t *testing.T) {
	macClaim := "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~NM1*QC*1*SYNTHETIC****MI*1AA0AA0AA00~MOA***MA01*N130~" +
		"SVC*HC:99213*100.00*80.00**1~DTM*472*20260105~CAS*CO*45*0.00~REF*6R*L1~AMT*B6*100.00~LQ*HE*N19~"
	tests := []struct {
		name   string
		vendor string
		want   []Difference
	}{
		{
			// Decimal spelling, a trailing empty SVC01 component, CLP07,
			// CLP10, a REF*F8, a second REF*6R, and a line QTY are not
			// differences: the first four read the same, the rest are left
			// out by set 1.
			name: "same values under the set",
			vendor: "CLP*A1*1*100*80.0*20*MB*ICN9*11*1*30~NM1*QC*1*SYNTHETIC****MI*1AA0AA0AA00~REF*F8*ICN1~MOA***MA01*N130~" +
				"SVC*HC:99213:*100*80**1~DTM*472*20260105~CAS*CO*45*0~REF*6R*L1~REF*6R*L1X~AMT*B6*100~QTY*NE*2~LQ*HE*N19~",
		},
		{
			// CLP04 is five cents under the MAC's, the edge of the five-cent
			// rule, so it is the same amount; SVC03 is six cents under, a
			// difference. A CO-253 of 0.01 sent on one side only is a
			// difference however small.
			name: "amounts at the five-cent edge, reordered remarks, missing allowed amount",
			vendor: "CLP*A1*1*100.00*79.95*20.00*MB*ICN9*11*1~NM1*QC*1*SYNTHETIC****MI*1AA0AA0AA00~MOA***N130*MA01~" +
				"SVC*HC:99213*100.00*79.94**1~DTM*472*20260105~CAS*CO*45*0.00**253*0.01~REF*6R*L1~LQ*HE*N19~",
			want: []Difference{
				{Element: "MOA", Occurrence: 1, MACValue: "MA01 N130", VendorValue: "N130 MA01"},
				{Element: "SVC03", Occurrence: 1, LineControlNumber: "L1", MACValue: "80", VendorValue: "79.94", AmountDifference: "-0.06"},
				{Element: "CAS*CO*253", Occurrence: 1, LineControlNumber: "L1", VendorValue: "0.01", AmountDifference: "0.01"},
				{Element: "AMT*B6", Occurrence: 1, LineControlNumber: "L1", MACValue: "100", AmountDifference: "-100"},
			},
		},
		{
			// An adjustment quantity has its own tag and no amount
			// difference; a segment the model does not place is compared
			// whole under RAW*ID.
			name: "adjustment quantity and an unplaced segment",
			vendor: "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~NM1*QC*1*SYNTHETIC****MI*1AA0AA0AA00~MOA***MA01*N130~K3*EXTRA DATA~" +
				"SVC*HC:99213*100.00*80.00**1~DTM*472*20260105~CAS*CO*45*0.00**253*1.50*2~REF*6R*L1~AMT*B6*100.00~LQ*HE*N19~",
			want: []Difference{
				{Element: "RAW*K3", Occurrence: 1, VendorValue: "EXTRA DATA"},
				{Element: "CAS*CO*253", Occurrence: 1, LineControlNumber: "L1", VendorValue: "1.5", AmountDifference: "1.5"},
				{Element: "CAS*CO*253*QTY", Occurrence: 1, LineControlNumber: "L1", VendorValue: "2"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mac := claim.Claim{Delimiters: "*:^~", Segments: []byte(macClaim)}
			vendor := claim.Claim{Delimiters: "*:^~", Segments: []byte(test.vendor)}
			got, err := DiffClaims(compareSet1, mac, vendor)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("differences:\n got %+v\nwant %+v", got, test.want)
			}
		})
	}
}

// TestCompareSetLeavesOut checks the set's matching: the exact tag at its
// level, and a trailing * covering every qualifier.
func TestCompareSetLeavesOut(t *testing.T) {
	tests := []struct {
		level, tag string
		want       bool
	}{
		{LevelClaim, "CLP07", true},
		{LevelClaim, "CLP04", false},
		{LevelClaim, "QTY*CA", false},
		{LevelLine, "QTY*CA", true},
		{LevelLine, "REF*6R", true},
		{LevelClaim, "REF*6R", false},
		{LevelLine, "CLP07", false},
	}
	for _, test := range tests {
		if got := compareSet1.leavesOut(test.level, test.tag); got != test.want {
			t.Errorf("leavesOut(%s, %s) = %v, want %v", test.level, test.tag, got, test.want)
		}
	}
	if _, err := LookupCompareSet("1"); err != nil {
		t.Errorf("LookupCompareSet(1) = %v", err)
	}
	if _, err := LookupCompareSet("none"); err == nil {
		t.Error("LookupCompareSet(none) succeeded, want an error")
	}
}

// TestCompareSubmissionsOutcomes checks the matching rule on a hand-built
// pair of submissions that reaches all six outcomes: a key the MAC repeats
// is set aside whether or not the vendor sent it; a key the MAC has once is
// missing, duplicated, or compared by how many vendor claims carry it; a key
// only the vendor has is extra, once, however many times it was sent.
func TestCompareSubmissionsOutcomes(t *testing.T) {
	segments := "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~NM1*QC*1*SYNTHETIC****MI*1AA0AA0AA00~" +
		"SVC*HC:99213*100.00*80.00**1~DTM*472*20260105~CAS*CO*45*0.00~REF*6R*L1~AMT*B6*100.00~"
	paidLess := strings.Replace(segments, "*80.00*20.00*", "*79.90*20.00*", 1) // ten cents less: past the five-cent rule
	row := func(key, status, paid string, text string) claim.Claim {
		return claim.Claim{MatchKey: []byte(key), ClaimStatusCode: status, ClaimPaymentAmount: paid, Delimiters: "*:^~", Segments: []byte(text)}
	}
	macClaims := []claim.Claim{
		row("A", "1", "80", segments), row("A", "1", "80", segments), // repeated by the MAC; the vendor sent it once
		row("B", "1", "80", segments), row("B", "4", "0", segments), // repeated by the MAC; the vendor did not send it
		row("C", "1", "80", segments), // once; the vendor sent it twice
		row("D", "1", "80", segments), // once; the vendor did not send it
		row("F", "1", "80", segments), // once each, equal
		row("G", "1", "80", segments), // once each, one difference
	}
	vendorClaims := []claim.Claim{
		row("A", "1", "80", segments),
		row("C", "1", "80", segments), row("C", "1", "80", segments),
		row("E", "1", "50", segments), row("E", "1", "50", segments), // only the vendor has it, twice
		row("F", "1", "80", segments),
		row("G", "1", "79.9", paidLess),
	}
	rows, diffs, summary, err := CompareSubmissions(compareSet1, "vendor", macClaims, vendorClaims)
	if err != nil {
		t.Fatal(err)
	}

	type want struct {
		key, outcome, macStatus, vendorStatus, macPaid, vendorPaid string
		differences                                                int
	}
	wants := []want{
		{"A", comparison.OutcomeSetAside, "1", "", "80", "", 0},
		{"B", comparison.OutcomeSetAside, "1", "", "80", "", 0},
		{"C", comparison.OutcomeDuplicated, "1", "", "80", "", 0},
		{"D", comparison.OutcomeMissing, "1", "", "80", "", 0},
		{"F", comparison.OutcomeEqual, "1", "1", "80", "80", 0},
		{"G", comparison.OutcomeDifferent, "1", "1", "80", "79.9", 1},
		{"E", comparison.OutcomeExtra, "", "1", "", "50", 0},
	}
	if len(rows) != len(wants) {
		t.Fatalf("CompareSubmissions returned %d rows, want %d", len(rows), len(wants))
	}
	for i, w := range wants {
		got := want{string(rows[i].MatchKey), rows[i].Outcome, rows[i].MACClaimStatusCode, rows[i].VendorClaimStatusCode,
			rows[i].MACClaimPaymentAmount, rows[i].VendorClaimPaymentAmount, rows[i].DifferenceCount}
		if got != w {
			t.Errorf("row %d = %+v, want %+v", i, got, w)
		}
		if rows[i].SubmissionPrefix != "vendor" {
			t.Errorf("row %d submission = %q, want vendor", i, rows[i].SubmissionPrefix)
		}
	}
	wantSummary := Summary{Equal: 1, Different: 1, Missing: 1, Duplicated: 1, Extra: 1, SetAside: 2, Differences: 1}
	if summary != wantSummary {
		t.Errorf("summary = %+v, want %+v", summary, wantSummary)
	}
	if len(diffs) != 1 || string(diffs[0].MatchKey) != "G" || diffs[0].Element != "CLP04" {
		t.Errorf("differences = %+v, want one CLP04 difference for key G", diffs)
	}
}
