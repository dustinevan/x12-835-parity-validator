// Package tagcell reads and writes the tag_cell table: the summary cube
// for what differs, one row per cell of the dimensions and tag, with the
// different pairs carrying the tag and the dollars of that element in them.
package tagcell

import "time"

// Filter selects tag cells of one submission the way cell.Filter selects
// cells: each list keeps the rows whose dimension has one of its values; an
// empty list keeps every row; a value may be "", the not-applicable
// bucket. ServiceMonth holds first-of-month dates.
type Filter struct {
	SubmissionPrefix string
	Jurisdiction     []string
	Part             []string
	Workload         []string
	ClaimType        []string
	Adjudication     []string
	MACStatus        []string
	VendorStatus     []string
	Frequency        []string
	PaidBand         []string
	ServiceMonth     []time.Time
	PaymentBand      []string
	Outcome          []string
}

// TagSums is one tag's totals over a set of cells: the different pairs
// carrying it and the dollars of the element in them, an exact decimal
// string.
type TagSums struct {
	Tag     string
	Claims  int64
	Dollars string
}
