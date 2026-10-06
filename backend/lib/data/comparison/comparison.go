// Package comparison reads and writes the comparison table: one row per
// (vendor submission, match key), the result most views read.
package comparison

import "time"

const (
	OutcomeEqual      = "equal"
	OutcomeDifferent  = "different"
	OutcomeMissing    = "missing"
	OutcomeExtra      = "extra"
	OutcomeSetAside   = "set_aside"
	OutcomeDuplicated = "duplicated"
)

// Comparison is one row. An empty status code or payment amount is NULL:
// the MAC side is NULL when extra, the vendor side when missing, duplicated,
// or set aside.
type Comparison struct {
	SubmissionPrefix              string
	MatchKey                      []byte
	Outcome                       string
	MACWorkloadID                 string
	ClaimType                     string
	MedicareBeneficiaryIdentifier string
	DateOfService                 time.Time
	TotalClaimChargeAmount        string
	MACClaimStatusCode            string
	VendorClaimStatusCode         string
	MACClaimPaymentAmount         string
	VendorClaimPaymentAmount      string
	ClaimFrequencyCode            string // CLP09 of the MAC claim (the vendor's when extra); "" when not sent
	Adjudication                  string // "inpatient", "outpatient", or "", from the same claim
	Jurisdiction                  string // the workload's jurisdiction from the same claim; "" when unknown
	Part                          string // "A", "B", "DME", or ""
	DifferenceCount               int
	Tags                          []string
}

// Filter selects comparison rows for the views. An empty ClaimType,
// Workload, or Outcome selects all. All, Any, and None are tag lists:
// every tag in All, at least one in Any, none of None. A nil or empty
// list selects everything. A non-empty OtherOutcome keeps only the keys
// whose outcome in the submission OtherSubmissionPrefix is OtherOutcome.
type Filter struct {
	SubmissionPrefix      string
	ClaimType             string
	Workload              string
	Outcome               string
	All                   []string
	Any                   []string
	None                  []string
	OtherSubmissionPrefix string
	OtherOutcome          string
}

// Totals is what the rows a Filter selects add up to: how many there are,
// and the sum of |vendor paid - MAC paid| over them, where a row with
// either side absent adds nothing.
type Totals struct {
	Claims    int64
	PaidDelta string
}

// AgreementCell counts the MAC claims with one outcome in the first
// submission (slot v1) and one in the second (slot v2). SameTags counts
// those different in both whose differing elements are the same in both.
type AgreementCell struct {
	V1Outcome string
	V2Outcome string
	Claims    int64
	SameTags  int64
}

// Bucket is one histogram bar: the smallest and largest paid delta in it,
// and how many claims it holds.
type Bucket struct {
	From  string
	To    string
	Count int64
}
