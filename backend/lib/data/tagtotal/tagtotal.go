// Package tagtotal reads and deletes the tag_total table: for each tag of
// a vendor submission's different claims, how many carry it and the dollars
// of that element, by claim type and workload. The summarize stage writes
// it (lib/core/query SummaryQuery); the differences view reads it.
package tagtotal

// TagTotal is one row. ClaimType and MACWorkloadID are "all" on total
// rows, as in measure. Dollars is a decimal string.
type TagTotal struct {
	SubmissionPrefix string
	ClaimType        string
	MACWorkloadID    string
	Tag              string
	Claims           int64
	Dollars          string
}
