package macs

import (
	"slices"
	"testing"
)

// jurisdictions is the list of MAC jurisdictions with their standard
// geographic assignments as USPS codes, checked against CMS's MAC
// jurisdiction pages (cms.gov, "Who are the MACs") on 2026-09-22. The table may use no other
// jurisdiction, and every jurisdiction here needs rows in the table.
var jurisdictions = map[string][]string{
	"J5":  {"IA", "KS", "MO", "NE"},
	"J6":  {"IL", "MN", "WI"},
	"J8":  {"IN", "MI"},
	"J15": {"KY", "OH"},
	"JE":  {"CA", "HI", "NV", "AS", "GU", "MP"},
	"JF":  {"AK", "AZ", "ID", "MT", "ND", "OR", "SD", "UT", "WA", "WY"},
	"JH":  {"AR", "CO", "LA", "MS", "NM", "OK", "TX"},
	"JJ":  {"AL", "GA", "TN"},
	"JK":  {"CT", "ME", "MA", "NH", "NY", "RI", "VT"},
	"JL":  {"DE", "DC", "MD", "NJ", "PA"},
	"JM":  {"NC", "SC", "VA", "WV"},
	"JN":  {"FL", "PR", "VI"},
	"JA":  {"CT", "DC", "DE", "MA", "MD", "ME", "NH", "NJ", "NY", "PA", "RI", "VT"},
	"JB":  {"IL", "IN", "KY", "MI", "MN", "OH", "WI"},
	"JC":  {"AL", "AR", "CO", "FL", "GA", "LA", "MS", "NC", "NM", "OK", "PR", "SC", "TN", "TX", "VA", "VI", "WV"},
	"JD":  {"AK", "AZ", "CA", "HI", "IA", "ID", "KS", "MO", "MT", "ND", "NE", "NV", "OR", "SD", "UT", "WA", "WY", "AS", "GU", "MP"},
}

// dme lists the DME jurisdictions among the jurisdiction list's keys.
var dme = []string{"JA", "JB", "JC", "JD"}

// jurisdictionOf maps the first two digits of a number to the jurisdiction
// that now holds it. The two digits are the jurisdiction indicator in CMS's
// numbering (CR 5651, 2007): A/B jurisdictions 1 to 15 as numbered before
// the mergers and renamings that produced today's letters, DME 16 to 19.
var jurisdictionOf = map[string]string{
	"01": "JE", "02": "JF", "03": "JF", "04": "JH", "05": "J5",
	"06": "J6", "07": "JH", "08": "J8", "09": "JN", "10": "JJ",
	"11": "JM", "12": "JL", "13": "JK", "14": "JK", "15": "J15",
	"16": "JA", "17": "JB", "18": "JC", "19": "JD",
}

// TestIDs checks that every ID is five digits and that none repeats.
func TestIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, w := range Workloads {
		if len(w.ID) != 5 {
			t.Errorf("ID %q has %d characters, want 5", w.ID, len(w.ID))
		}
		for _, c := range w.ID {
			if c < '0' || c > '9' {
				t.Errorf("ID %q has a character that is not a digit", w.ID)
				break
			}
		}
		if seen[w.ID] {
			t.Errorf("ID %q appears more than once", w.ID)
		}
		seen[w.ID] = true
	}
}

// TestJurisdictions checks every row's jurisdiction against the jurisdiction list's
// list and against the first two digits of its ID, keeps DME numbers with
// the DME jurisdictions, and checks that every jurisdiction in the jurisdiction list's
// list has rows: both Parts for an A/B jurisdiction, DME for the others.
func TestJurisdictions(t *testing.T) {
	parts := map[string]map[string]bool{}
	for _, w := range Workloads {
		if _, ok := jurisdictions[w.Jurisdiction]; !ok {
			t.Errorf("%s: jurisdiction %q is not in the jurisdiction list", w.ID, w.Jurisdiction)
		}
		if len(w.ID) == 5 && jurisdictionOf[w.ID[:2]] != w.Jurisdiction {
			t.Errorf("%s: jurisdiction %q, but the first two digits belong to %q", w.ID, w.Jurisdiction, jurisdictionOf[w.ID[:2]])
		}
		if isDME := slices.Contains(dme, w.Jurisdiction); isDME != (w.Part == "DME") {
			t.Errorf("%s: jurisdiction %q with Part %q, want Part DME only for the DME jurisdictions", w.ID, w.Jurisdiction, w.Part)
		}
		if w.Contractor == "" {
			t.Errorf("%s: contractor is empty", w.ID)
		}
		if parts[w.Jurisdiction] == nil {
			parts[w.Jurisdiction] = map[string]bool{}
		}
		parts[w.Jurisdiction][w.Part] = true
	}
	for j := range jurisdictions {
		want := []string{"A", "B"}
		if slices.Contains(dme, j) {
			want = []string{"DME"}
		}
		for _, p := range want {
			if !parts[j][p] {
				t.Errorf("jurisdiction %s has no Part %s row", j, p)
			}
		}
	}
}

// TestPartAndClaimType checks Part against the last digit of a state-level
// workload number, which CMS defines as the claim type: 1 is Part A and 2 is
// Part B. Roll-up, HHH, and DME numbers have 0 as the third digit and are
// left to TestJurisdictions.
func TestPartAndClaimType(t *testing.T) {
	want := map[byte]string{'1': "A", '2': "B"}
	for _, w := range Workloads {
		if len(w.ID) != 5 || w.ID[2] == '0' {
			continue
		}
		if got := want[w.ID[4]]; got != w.Part {
			t.Errorf("%s: Part %q, but the last digit means Part %q", w.ID, w.Part, got)
		}
	}
}

// TestStates checks that every state is a two-letter code inside the
// jurisdiction's own states. A number ending in 4 is an HHH number; its
// states are the HHH area, which differs from the jurisdiction's, so it is
// left out.
func TestStates(t *testing.T) {
	for _, w := range Workloads {
		if len(w.ID) == 5 && w.ID[4] == '4' {
			continue
		}
		for _, s := range w.States {
			if len(s) != 2 {
				t.Errorf("%s: state %q is not a two-letter code", w.ID, s)
			}
			if !slices.Contains(jurisdictions[w.Jurisdiction], s) {
				t.Errorf("%s: state %s is outside jurisdiction %s", w.ID, s, w.Jurisdiction)
			}
		}
	}
}

// TestLookup checks a few numbers against the CMS documents they came from,
// and that numbers the table does not hold, including the synthetic IDs
// tools/gen writes, are not found.
func TestLookup(t *testing.T) {
	tests := []struct {
		id           string
		wantOK       bool
		jurisdiction string
		part         string
		contractor   string
		states       []string
	}{
		{"05101", true, "J5", "A", wps, []string{"IA"}},
		{"04112", true, "JH", "B", novitas, []string{"CO"}},
		{"13292", true, "JK", "B", ngs, []string{"NY"}},
		{"01182", true, "JE", "B", noridian, []string{"CA"}},
		{"15001", true, "J15", "A", cgs, []string{"KY", "OH"}},
		{"16013", true, "JA", "DME", noridian, []string{"CT", "DC", "DE", "MA", "MD", "ME", "NH", "NJ", "NY", "PA", "RI", "VT"}},
		{"05901", true, "J5", "A", wps, nil},
		{"91101", false, "", "", "", nil},
		{"16003", false, "", "", "", nil},
		{"", false, "", "", "", nil},
		{"5101", false, "", "", "", nil},
		{" 05101", false, "", "", "", nil},
	}
	for _, tt := range tests {
		got, ok := Lookup(tt.id)
		if ok != tt.wantOK {
			t.Errorf("Lookup(%q) found = %v, want %v", tt.id, ok, tt.wantOK)
			continue
		}
		if got.Jurisdiction != tt.jurisdiction || got.Part != tt.part || got.Contractor != tt.contractor || !slices.Equal(got.States, tt.states) {
			t.Errorf("Lookup(%q) = %+v, want jurisdiction %q, part %q, contractor %q, states %v",
				tt.id, got, tt.jurisdiction, tt.part, tt.contractor, tt.states)
		}
	}
}

// TestLookupEveryRow checks that Lookup returns each row of the table for
// its own ID.
func TestLookupEveryRow(t *testing.T) {
	for _, w := range Workloads {
		got, ok := Lookup(w.ID)
		if !ok {
			t.Errorf("Lookup(%q) found = false, want true", w.ID)
			continue
		}
		if got.ID != w.ID || got.Jurisdiction != w.Jurisdiction || got.Part != w.Part || got.Contractor != w.Contractor || !slices.Equal(got.States, w.States) {
			t.Errorf("Lookup(%q) = %+v, want %+v", w.ID, got, w)
		}
	}
}
