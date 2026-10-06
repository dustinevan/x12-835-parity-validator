package views

import (
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/submission"
)

func TestRatio(t *testing.T) {
	tests := []struct {
		n, d int64
		want string
	}{
		{8, 1000, "0.008"},
		{1, 3, "0.333333"},
		{2, 3, "0.666667"},
		{1000, 1000, "1"},
		{0, 1000, "0"},
		{0, 0, "0"},
	}
	for _, tt := range tests {
		if got := ratio(tt.n, tt.d); got != tt.want {
			t.Errorf("ratio(%d, %d) = %q, want %q", tt.n, tt.d, got, tt.want)
		}
	}
}

// TestShown checks the status gate: the MAC submission once loaded, a
// vendor submission only once summarized.
func TestShown(t *testing.T) {
	tests := []struct {
		side, status string
		want         bool
	}{
		{submission.SideMAC, submission.StatusLoaded, true},
		{submission.SideVendor, submission.StatusLoaded, false},
		{submission.SideVendor, submission.StatusCompared, false},
		{submission.SideVendor, submission.StatusSummarized, true},
		{submission.SideVendor, submission.StatusFailed, false},
	}
	for _, tt := range tests {
		d := submission.Submission{Side: tt.side, Status: tt.status}
		if got := shown(d); got != tt.want {
			t.Errorf("shown(%s, %s) = %v, want %v", tt.side, tt.status, got, tt.want)
		}
	}
}
