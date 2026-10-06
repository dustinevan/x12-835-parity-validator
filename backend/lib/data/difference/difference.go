// Package difference reads and writes the difference table: one row per
// differing element in a compared claim pair.
package difference

// Difference is one row. An empty value is NULL: the element is absent on
// that side. An empty AmountDifference is NULL: the element is not an amount.
type Difference struct {
	SubmissionPrefix  string
	MatchKey          []byte
	Element           string
	Occurrence        int
	LineControlNumber string
	MACValue          string
	VendorValue       string
	AmountDifference  string
}
