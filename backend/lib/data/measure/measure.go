// Package measure reads and deletes the measure table: one vendor
// submission's headline measures, written by the summarize stage and read by
// the headline views. The rows are built from other tables, so the
// INSERT ... SELECT is in lib/core/query (SummaryQuery).
package measure

// All is the claim_type and mac_workload_id value of a total row.
const All = "all"

// Measure is one row. Every count except ExtraCount counts MAC claims: a
// set-aside key counts each MAC claim that shares it, and a duplicated key
// (the MAC has it once, the vendor sent it more than once) counts its one
// MAC claim. Rates and deviations are decimal strings; empty is NULL (no
// claims to measure).
type Measure struct {
	SubmissionPrefix string
	ClaimType        string
	MACWorkloadID    string
	MACClaims        int64
	EqualCount       int64
	DifferentCount   int64
	MissingCount     int64
	ExtraCount       int64
	SetAsideCount    int64
	DuplicatedCount  int64
	OutcomeMatchRate string
	PaidClaims       int64
	MedianDeviation  string
	P95Deviation     string
}
