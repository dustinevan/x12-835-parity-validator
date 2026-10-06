// Package cell reads and writes the cell table: the summary cube one
// summarize run writes for a vendor submission, one row per distinct
// combination of its twelve dimensions with counts and sums. The Overview
// reads it with a Filter: WHERE on the dimensions and SUM over the rows,
// never a group of the comparison rows.
package cell

import "time"

// The dimensions, by the names the API uses for their filters, in the order
// the API lists them.
const (
	DimensionJurisdiction = "jurisdiction"
	DimensionPart         = "part"
	DimensionWorkload     = "workload"
	DimensionClaimType    = "claim_type"
	DimensionAdjudication = "adjudication"
	DimensionMACStatus    = "mac_status"
	DimensionVendorStatus = "vendor_status"
	DimensionFrequency    = "frequency"
	DimensionPaidBand     = "paid_band"
	DimensionServiceMonth = "service_month"
	DimensionPaymentBand  = "payment_band"
	DimensionOutcome      = "outcome"
)

// Dimensions lists every dimension in the API's order.
var Dimensions = []string{
	DimensionJurisdiction,
	DimensionPart,
	DimensionWorkload,
	DimensionClaimType,
	DimensionAdjudication,
	DimensionMACStatus,
	DimensionVendorStatus,
	DimensionFrequency,
	DimensionPaidBand,
	DimensionServiceMonth,
	DimensionPaymentBand,
	DimensionOutcome,
}

// The paid bands, in order: the MAC paid amount (the vendor's when extra).
var PaidBands = []string{"zero", "under 100", "under 1000", "under 10000", "10000 and more"}

// The payment bands, in order: |vendor paid - MAC paid| on a pair, "match"
// being within five cents, the comparison's five-cent rule. A row that is
// not a pair has "".
var PaymentBands = []string{"match", "up to 1.00", "up to 1%", "up to 5%", "over 5%"}

// Filter selects cells of one submission. Each list keeps the rows whose
// dimension has one of its values; an empty list keeps every row. A value
// may be "", the not-applicable bucket. ServiceMonth holds first-of-month
// dates.
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

// Sums are the cube's counts and sums added over a set of cells. Amounts
// are exact decimal strings, trailing fractional zeros removed. Claims
// counts every comparison row; MACClaims counts the rows that have a MAC
// claim, which is every outcome but extra.
type Sums struct {
	Claims          int64
	MACClaims       int64
	MACPaid         string
	VendorPaid      string
	Billed          string
	NetDifference   string
	GrossDifference string
	SameStatus      int64
}

// OutcomeSums is the sums of one outcome.
type OutcomeSums struct {
	Outcome string
	Sums
}

// ValueSums is the sums of one value of a dimension.
type ValueSums struct {
	Value string
	Sums
}

// StatusCell counts the pairs with one MAC status and one vendor status.
type StatusCell struct {
	MACStatus    string
	VendorStatus string
	Claims       int64
	MACPaid      string
	VendorPaid   string
}
