package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func smallConfig(out string) config {
	return config{
		Out:         out,
		Claims:      150,
		Seed:        7,
		Submissions: []string{"2026-10-02", "2026-10-16", "2026-10-30"},
		Vendors:     []string{"vendor-1", "vendor-2"},
	}
}

// coverageConfig is large enough for every outcome, band, status, and
// pattern to appear at the profiles' rates.
func coverageConfig(out string) config {
	cfg := smallConfig(out)
	cfg.Claims = 6000
	return cfg
}

func readTree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		files[rel] = data
		return nil
	})
	if err != nil {
		t.Fatalf("reading %s: %v", root, err)
	}
	return files
}

func readTruth(t *testing.T, data []byte) []truthRecord {
	t.Helper()
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	var out []truthRecord
	for sc.Scan() {
		var r truthRecord
		if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
			t.Fatalf("truth line: %v", err)
		}
		out = append(out, r)
	}
	return out
}

func TestSameSeedSameBytes(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	if _, err := run(smallConfig(a)); err != nil {
		t.Fatal(err)
	}
	if _, err := run(smallConfig(b)); err != nil {
		t.Fatal(err)
	}
	fa, fb := readTree(t, a), readTree(t, b)
	if len(fa) != len(fb) {
		t.Fatalf("file counts differ: %d and %d", len(fa), len(fb))
	}
	for name, data := range fa {
		if !bytes.Equal(data, fb[name]) {
			t.Errorf("%s differs between two runs with the same seed", name)
		}
	}

	c := t.TempDir()
	other := smallConfig(c)
	other.Seed = 8
	if _, err := run(other); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(fa["truth.jsonl"], readTree(t, c)["truth.jsonl"]) {
		t.Error("a different seed produced the same truth file")
	}
}

// cents parses an X12 decimal into cents.
func cents(t *testing.T, s string) int64 {
	t.Helper()
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	whole, frac, _ := strings.Cut(s, ".")
	frac = (frac + "00")[:2]
	n, err := strconv.ParseInt(whole+frac, 10, 64)
	if err != nil {
		t.Fatalf("parsing amount %q: %v", s, err)
	}
	if neg {
		return -n
	}
	return n
}

// casTotal sums the amounts in a CAS segment and the part owed by the patient.
func casTotal(t *testing.T, el []string) (total, patient int64) {
	for i := 3; i < len(el); i += 3 {
		if el[i] == "" {
			continue
		}
		a := cents(t, el[i])
		total += a
		if el[1] == "PR" {
			patient += a
		}
	}
	return total, patient
}

func TestEveryRemittanceBalances(t *testing.T) {
	out := t.TempDir()
	if _, err := run(smallConfig(out)); err != nil {
		t.Fatal(err)
	}
	files := readTree(t, out)
	checked := 0
	for name, data := range files {
		if !strings.HasSuffix(name, ".835") {
			continue
		}
		checkFile(t, name, data)
		checked++
	}
	if checked == 0 {
		t.Fatal("no 835 files written")
	}
}

// checkFile checks one file's arithmetic and conventions: every line and
// claim balances, BPR02 is the claims less the PLBs, SE01 counts the
// segments, CLP09 is sent on institutional claims only, and a reversal
// names the claim it undoes in REF*F8.
func checkFile(t *testing.T, name string, data []byte) {
	t.Helper()
	var segs [][]string
	for _, s := range strings.Split(string(data), "~") {
		s = strings.TrimSpace(s)
		if s != "" {
			segs = append(segs, strings.Split(s, "*"))
		}
	}

	var stStart int
	var paid, plb, bpr int64
	var claimCharge, claimPaid, claimPatient, claimCAS, claimPR int64
	var lineCharge, linePaid, lineCAS int64
	inClaim, inLine, reversal, sawF8 := false, false, false, false
	closeLine := func() {
		if inLine && lineCharge-linePaid != lineCAS {
			t.Errorf("%s: line does not balance: charge %d paid %d adjustments %d", name, lineCharge, linePaid, lineCAS)
		}
		inLine = false
	}
	closeClaim := func() {
		closeLine()
		if !inClaim {
			return
		}
		if claimCharge-claimPaid != claimCAS {
			t.Errorf("%s: claim does not balance: charge %d paid %d adjustments %d", name, claimCharge, claimPaid, claimCAS)
		}
		if claimPatient != claimPR {
			t.Errorf("%s: CLP05 %d is not the PR total %d", name, claimPatient, claimPR)
		}
		if reversal && !sawF8 {
			t.Errorf("%s: a reversal (CLP02 22) carries no REF*F8", name)
		}
		inClaim = false
	}
	for i, el := range segs {
		switch el[0] {
		case "ST":
			stStart = i
			paid, plb = 0, 0
		case "BPR":
			bpr = cents(t, el[2])
		case "CLP":
			closeClaim()
			inClaim = true
			claimCharge, claimPaid, claimPatient = cents(t, el[3]), cents(t, el[4]), cents(t, el[5])
			claimCAS, claimPR = 0, 0
			paid += claimPaid
			reversal, sawF8 = el[2] == "22", false
			institutional := el[6] == "MA"
			hasFrequency := len(el) > 9 && el[9] != ""
			if institutional && !hasFrequency {
				t.Errorf("%s: an institutional claim has no CLP09", name)
			}
			if !institutional && hasFrequency {
				t.Errorf("%s: a non-institutional claim carries CLP09 %q", name, el[9])
			}
		case "REF":
			if inClaim && !inLine && el[1] == "F8" {
				sawF8 = true
			}
		case "SVC":
			closeLine()
			inLine = true
			lineCharge, linePaid, lineCAS = cents(t, el[2]), cents(t, el[3]), 0
		case "CAS":
			total, patient := casTotal(t, el)
			claimCAS += total
			claimPR += patient
			if inLine {
				lineCAS += total
			}
		case "PLB":
			closeClaim()
			for j := 4; j < len(el); j += 2 {
				plb += cents(t, el[j])
			}
		case "SE":
			closeClaim()
			if bpr != paid-plb {
				t.Errorf("%s: BPR02 %d is not claim payments %d less PLB %d", name, bpr, paid, plb)
			}
			if n, _ := strconv.Atoi(el[1]); n != i-stStart+1 {
				t.Errorf("%s: SE01 %d, counted %d segments", name, n, i-stStart+1)
			}
		}
	}
}

// paymentBand is the Overview's payment-difference band, as the
// validator's summary cube defines it: "match" is within matchCents, the
// five-cent rule, so a pair the truth records no CLP04 difference for is a
// match.
func paymentBand(diff, macPaid int64) string {
	if diff < 0 {
		diff = -diff
	}
	switch {
	case diff <= matchCents:
		return "match"
	case diff <= 100:
		return "up to 1.00"
	case macPaid > 0 && diff*100 <= macPaid:
		return "up to 1%"
	case macPaid > 0 && diff*100 <= 5*macPaid:
		return "up to 5%"
	}
	return "over 5%"
}

// isMoney reports whether a truth difference's element is a dollar amount,
// which the comparison compares under the five-cent rule.
func isMoney(element string) bool {
	switch element {
	case "CLP03", "CLP04", "CLP05", "MIA02", "MIA04", "SVC02", "SVC03":
		return true
	}
	return strings.HasPrefix(element, "CAS*") || strings.HasPrefix(element, "AMT*")
}

// TestCoverage checks that the first submission of each vendor puts
// rows in every outcome and every payment band, that every pattern in
// the catalog shows up somewhere, that no recorded difference is two
// amounts within five cents of each other, and that the MAC side carries
// the claim statuses and frequency codes the Overview breaks out by.
func TestCoverage(t *testing.T) {
	out := t.TempDir()
	cfg := coverageConfig(out)
	if _, err := run(cfg); err != nil {
		t.Fatal(err)
	}
	files := readTree(t, out)
	truth := readTruth(t, files["truth.jsonl"])

	outcomes := map[string]map[string]int{}
	bands := map[string]map[string]int{}
	patterns := map[string]int{}
	for _, r := range truth {
		if outcomes[r.Submission] == nil {
			outcomes[r.Submission] = map[string]int{}
			bands[r.Submission] = map[string]int{}
		}
		outcomes[r.Submission][r.Outcome]++
		for _, s := range r.Patterns {
			patterns[s]++
		}
		if r.Outcome != "matched" {
			continue
		}
		macPaid, diff := int64(0), int64(0)
		for _, df := range r.Differences {
			if df.Element == "CLP04" && df.Line == "" {
				macPaid = cents(t, df.MAC)
				diff = cents(t, df.Vendor) - macPaid
			}
			if isMoney(df.Element) && df.MAC != "" && df.Vendor != "" {
				if apart := cents(t, df.Vendor) - cents(t, df.MAC); apart >= -matchCents && apart <= matchCents {
					t.Errorf("%s: %s line %q is recorded, but %s and %s are within five cents", r.Submission, df.Element, df.Line, df.MAC, df.Vendor)
				}
			}
		}
		bands[r.Submission][paymentBand(diff, macPaid)]++
	}
	for _, v := range cfg.Vendors {
		first := v + "/" + cfg.Submissions[0]
		for _, o := range []string{"matched", "missing", "extra", "set_aside"} {
			if outcomes[first][o] == 0 {
				t.Errorf("%s: no %s records", first, o)
			}
		}
		for _, b := range []string{"match", "up to 1.00", "up to 1%", "up to 5%", "over 5%"} {
			if bands[first][b] == 0 {
				t.Errorf("%s: no pair in the %q payment band", first, b)
			}
		}
	}
	if n := outcomes["vendor-2/"+cfg.Submissions[0]]["duplicated"]; n == 0 {
		t.Errorf("vendor-2's first submission has no duplicated records")
	}
	for n, s := range catalog {
		if n > 0 && patterns[s.Name] == 0 {
			t.Errorf("pattern %d (%s) never shows", n, s.Name)
		}
	}

	statuses := map[string]int{}
	frequencies := map[string]int{}
	for name, data := range files {
		if !strings.HasPrefix(name, "mac/") {
			continue
		}
		for _, s := range strings.Split(string(data), "~") {
			el := strings.Split(strings.TrimSpace(s), "*")
			if el[0] != "CLP" {
				continue
			}
			statuses[el[2]]++
			if len(el) > 9 && el[9] != "" {
				frequencies[el[9]]++
			}
		}
	}
	for _, s := range []string{"1", "2", "4", "19", "22", "23"} {
		if statuses[s] == 0 {
			t.Errorf("no MAC claim with CLP02 %s", s)
		}
	}
	for _, f := range []string{"1", "7", "8"} {
		if frequencies[f] == 0 {
			t.Errorf("no institutional MAC claim with CLP09 %s", f)
		}
	}
}

func TestNPICheckDigit(t *testing.T) {
	// 1234567893 is the check-digit example from the NPI standard.
	if got := npiCheckDigit("123456789"); got != "3" {
		t.Errorf("check digit for 123456789: got %s, want 3", got)
	}
}

func TestSequestrationSpreadsClaimRounding(t *testing.T) {
	before := []Money{6275, 6275}
	claim := sequestration(before, false, false)
	line := sequestration(before, true, false)
	truncated := sequestration(before, false, true)
	if claim[0]+claim[1] != 251 || line[0]+line[1] != 252 || truncated[0]+truncated[1] != 251 {
		t.Errorf("per claim %v, per line %v, truncated %v", claim, line, truncated)
	}
	// 2 percent of 12.55 is 0.251: rounded 0.25, truncated 0.25; of 37.75
	// it is 0.755: rounded 0.76, truncated 0.75.
	if got := sequestration([]Money{3775}, false, false); got[0] != 76 {
		t.Errorf("rounded sequestration of 37.75: got %d, want 76", got[0])
	}
	if got := sequestration([]Money{3775}, false, true); got[0] != 75 {
		t.Errorf("truncated sequestration of 37.75: got %d, want 75", got[0])
	}
}
