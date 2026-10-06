// Package maccell reads and writes the mac_cell table: the legacy cube load
// writes for the MAC submission, one row per distinct combination of the
// nine dimensions the MAC claim itself carries, with the claims and what
// was paid and billed. The Overview reads it with a Filter, WHERE on the
// dimensions and SUM over the rows, for the legacy side's own figures,
// which no vendor's results change.
package maccell

import "time"

// The dimensions, by the names the API uses for their filters, in the order
// the API lists them: the vendor cube's dimensions (lib/data/cell) less the
// three a vendor decides, its claim status, the payment band, and the
// outcome.
const (
	DimensionJurisdiction = "jurisdiction"
	DimensionPart         = "part"
	DimensionWorkload     = "workload"
	DimensionClaimType    = "claim_type"
	DimensionAdjudication = "adjudication"
	DimensionMACStatus    = "mac_status"
	DimensionFrequency    = "frequency"
	DimensionPaidBand     = "paid_band"
	DimensionServiceMonth = "service_month"
)

// Dimensions lists every dimension in the API's order.
var Dimensions = []string{
	DimensionJurisdiction,
	DimensionPart,
	DimensionWorkload,
	DimensionClaimType,
	DimensionAdjudication,
	DimensionMACStatus,
	DimensionFrequency,
	DimensionPaidBand,
	DimensionServiceMonth,
}

// Filter selects cells of the MAC submission. Each list keeps the rows
// whose dimension has one of its values; an empty list keeps every row. A
// value may be "", the not-applicable bucket. MACStatus selects by the MAC
// claim's status (claim_status_code). ServiceMonth holds first-of-month
// dates.
type Filter struct {
	SubmissionPrefix string
	Jurisdiction     []string
	Part             []string
	Workload         []string
	ClaimType        []string
	Adjudication     []string
	MACStatus        []string
	Frequency        []string
	PaidBand         []string
	ServiceMonth     []time.Time
}

// Sums are the legacy cube's counts and sums added over a set of cells:
// the MAC claims, the sum of their CLP04, and the sum of their CLP03.
// Amounts are exact decimal strings, trailing fractional zeros removed.
type Sums struct {
	Claims int64
	Paid   string
	Billed string
}

// ValueSums is the sums of one value of a dimension.
type ValueSums struct {
	Value string
	Sums
}
