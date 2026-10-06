package carc

import "testing"

// TestEntries checks the table: every code is non-empty and listed once,
// the index finds every row, and every status is one of the two that Entry
// and /api/codes document, so a refreshed list with a new status fails here
// instead of reaching the API unannounced.
func TestEntries(t *testing.T) {
	seen := map[string]bool{}
	for i, e := range entries {
		if e.Code == "" || seen[e.Code] {
			t.Errorf("entries[%d]: code %q is empty or repeated", i, e.Code)
		}
		seen[e.Code] = true
		if index[Code(e.Code)] != i {
			t.Errorf("index[%q] = %d, want %d", e.Code, index[Code(e.Code)], i)
		}
		if e.Status != "current" && e.Status != "deactivated" {
			t.Errorf("entries[%d] (%s): status %q, want current or deactivated", i, e.Code, e.Status)
		}
	}
	if len(index) != len(entries) {
		t.Errorf("index has %d codes, entries %d", len(index), len(entries))
	}
}

// TestSpotCheck checks a few codes against the published list, so an
// update that drops or garbles rows fails loudly.
func TestSpotCheck(t *testing.T) {
	tests := []struct {
		code            Code
		wantShort       string
		wantDescription string
		wantStatus      string
	}{
		{"45", "Charge exceeds fee schedule/maximum allowable or contracted/legislated fee arrangement", "Charge exceeds fee schedule/maximum allowable or contracted/legislated fee arrangement. Usage: This adjustment amount cannot equal the total service or claim charge amount; and must not duplicate provider adjustment amounts (payments and contractual reductions) that have resulted from prior payer(s) adjudication. (Use only with Group Codes PR or CO depending upon liability)", "current"},
		{"253", "Sequestration - reduction in federal payment", "Sequestration - reduction in federal payment", "current"},
	}
	for _, tt := range tests {
		if got := tt.code.ShortDescription(); got != tt.wantShort {
			t.Errorf("%s.ShortDescription() = %q, want %q", tt.code, got, tt.wantShort)
		}
		if got := tt.code.Description(); got != tt.wantDescription {
			t.Errorf("%s.Description() = %q, want %q", tt.code, got, tt.wantDescription)
		}
		if got := tt.code.Status(); got != tt.wantStatus {
			t.Errorf("%s.Status() = %q, want %q", tt.code, got, tt.wantStatus)
		}
	}
}

// TestParse checks that every listed code parses to itself with no error,
// that an empty value is CodeNotSent, and that a value off the list comes
// back as itself with an error, so nothing received is ever dropped.
func TestParse(t *testing.T) {
	for _, c := range Codes() {
		parsed, err := Parse(c.String())
		if err != nil || parsed != c || !parsed.Known() {
			t.Errorf("Parse(%q) = %v, %v, want %v known, nil", c.String(), parsed, err, c)
		}
	}
	if parsed, err := Parse(""); err != nil || parsed != CodeNotSent || parsed.Known() {
		t.Errorf("Parse(\"\") = %v, %v, want CodeNotSent, nil", parsed, err)
	}
	parsed, err := Parse("not-a-real-code")
	if err == nil || parsed != Code("not-a-real-code") || parsed.Known() || parsed.Description() != "" {
		t.Errorf("Parse(unknown) = %v, %v, want the value kept, an error", parsed, err)
	}
}
