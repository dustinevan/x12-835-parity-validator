// Package submission reads and writes the submission table: one row per set of
// 835 files that arrived together, the MAC's or a vendor's. Load writes
// the row; compare and summarize advance its status.
package submission

import "time"

const (
	SideMAC    = "mac"
	SideVendor = "vendor"

	// A submission's status moves loaded -> compared -> summarized, or
	// failed. The MAC submission stays loaded. The views show a vendor
	// submission only when it is summarized.
	StatusLoaded     = "loaded"
	StatusCompared   = "compared"
	StatusSummarized = "summarized"
	StatusFailed     = "failed" // compare sets it when a workload fails; summarize refuses it
)

// Submission is one row. A zero SubmittedOn or ComparedAt, or an empty
// ComparedSet, means NULL. ComparedSet names the compared set (which
// elements the comparison left out) of the compare run that wrote the
// submission's rows.
type Submission struct {
	Prefix      string
	Side        string
	Label       string
	SubmittedOn time.Time
	Status      string
	ClaimCount  int64
	LoadedAt    time.Time
	ComparedAt  time.Time
	ComparedSet string
}

// NewSubmission is what the loader knows before the claims are written.
// A zero SubmittedOn is stored as NULL (the MAC submission has no date).
type NewSubmission struct {
	Prefix      string
	Side        string
	Label       string
	SubmittedOn time.Time
	ClaimCount  int64
}
