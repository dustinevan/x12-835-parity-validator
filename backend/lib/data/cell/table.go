package cell

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for the summary cube.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const deleteBySubmissionSQL = `DELETE FROM public.cell WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("cell: deleting submission: %w", err)
	}
	return nil
}

// args are a Filter's named arguments. A list that is nil is sent empty,
// since cardinality(NULL) is NULL and would keep no row; @dimension names
// the dimension whose own filter a statement leaves out ("" leaves none
// out).
func args(f Filter, dimension string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"submissionPrefix": f.SubmissionPrefix,
		"dimension":        dimension,
		"jurisdiction":     list(f.Jurisdiction),
		"part":             list(f.Part),
		"workload":         list(f.Workload),
		"claimType":        list(f.ClaimType),
		"adjudication":     list(f.Adjudication),
		"macStatus":        list(f.MACStatus),
		"vendorStatus":     list(f.VendorStatus),
		"frequency":        list(f.Frequency),
		"paidBand":         list(f.PaidBand),
		"serviceMonth":     dates(f.ServiceMonth),
		"paymentBand":      list(f.PaymentBand),
		"outcome":          list(f.Outcome),
	}
}

func list(xs []string) []string {
	if xs == nil {
		return []string{}
	}
	return xs
}

func dates(xs []time.Time) []time.Time {
	if xs == nil {
		return []time.Time{}
	}
	return xs
}

// totalsSQL adds every cell the filter keeps. Each sum is written with
// trailing fractional zeros removed, and an empty selection reads as zero.
// The second count leaves out the extra rows, the only ones with no MAC
// claim.
const totalsSQL = `
SELECT
    COALESCE(sum(claims), 0)::bigint
    , COALESCE(sum(claims) FILTER (WHERE outcome <> 'extra'), 0)::bigint
    , trim_scale(COALESCE(sum(mac_paid), 0))::text
    , trim_scale(COALESCE(sum(vendor_paid), 0))::text
    , trim_scale(COALESCE(sum(billed), 0))::text
    , trim_scale(COALESCE(sum(net_difference), 0))::text
    , trim_scale(COALESCE(sum(gross_difference), 0))::text
    , COALESCE(sum(same_status), 0)::bigint
FROM public.cell
WHERE submission_prefix = @submissionPrefix
    AND (cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
    AND (cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
    AND (cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
    AND (cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
    AND (cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
    AND (cardinality(@macStatus::text[]) = 0 OR mac_claim_status_code = ANY (@macStatus::text[]))
    AND (cardinality(@vendorStatus::text[]) = 0 OR vendor_claim_status_code = ANY (@vendorStatus::text[]))
    AND (cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
    AND (cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
    AND (cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))
    AND (cardinality(@paymentBand::text[]) = 0 OR payment_band = ANY (@paymentBand::text[]))
    AND (cardinality(@outcome::text[]) = 0 OR outcome = ANY (@outcome::text[]))`

// Totals adds every cell the filter keeps.
func (t TableAccess) Totals(ctx context.Context, f Filter) (Sums, error) {
	var s Sums
	err := t.db.QueryRow(ctx, totalsSQL, args(f, "")).Scan(
		&s.Claims,
		&s.MACClaims,
		&s.MACPaid,
		&s.VendorPaid,
		&s.Billed,
		&s.NetDifference,
		&s.GrossDifference,
		&s.SameStatus,
	)
	if err != nil {
		return Sums{}, fmt.Errorf("cell: adding totals: %w", err)
	}
	return s, nil
}

// outcomesSQL adds the cells the filter keeps by outcome, in the order the
// outcomes are reported: equal, different, missing, extra, set_aside,
// duplicated. An outcome with no cells has no row. The second count leaves
// out the extra rows, the only ones with no MAC claim, so it is zero on the
// extra row and the row's claims on every other.
const outcomesSQL = `
SELECT
    outcome
    , sum(claims)::bigint
    , COALESCE(sum(claims) FILTER (WHERE outcome <> 'extra'), 0)::bigint
    , trim_scale(sum(mac_paid))::text
    , trim_scale(sum(vendor_paid))::text
    , trim_scale(sum(billed))::text
    , trim_scale(sum(net_difference))::text
    , trim_scale(sum(gross_difference))::text
    , sum(same_status)::bigint
FROM public.cell
WHERE submission_prefix = @submissionPrefix
    AND (cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
    AND (cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
    AND (cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
    AND (cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
    AND (cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
    AND (cardinality(@macStatus::text[]) = 0 OR mac_claim_status_code = ANY (@macStatus::text[]))
    AND (cardinality(@vendorStatus::text[]) = 0 OR vendor_claim_status_code = ANY (@vendorStatus::text[]))
    AND (cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
    AND (cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
    AND (cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))
    AND (cardinality(@paymentBand::text[]) = 0 OR payment_band = ANY (@paymentBand::text[]))
    AND (cardinality(@outcome::text[]) = 0 OR outcome = ANY (@outcome::text[]))
GROUP BY outcome
ORDER BY array_position(ARRAY['equal', 'different', 'missing', 'extra', 'set_aside', 'duplicated'], outcome)`

// Outcomes adds the cells the filter keeps, one row per outcome present.
func (t TableAccess) Outcomes(ctx context.Context, f Filter) ([]OutcomeSums, error) {
	rows, err := t.db.Query(ctx, outcomesSQL, args(f, ""))
	if err != nil {
		return nil, fmt.Errorf("cell: adding by outcome: %w", err)
	}
	defer rows.Close()
	var out []OutcomeSums
	for rows.Next() {
		var o OutcomeSums
		err := rows.Scan(
			&o.Outcome,
			&o.Claims,
			&o.MACClaims,
			&o.MACPaid,
			&o.VendorPaid,
			&o.Billed,
			&o.NetDifference,
			&o.GrossDifference,
			&o.SameStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("cell: scanning outcome: %w", err)
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cell: adding by outcome: %w", err)
	}
	return out, nil
}

// facetSQL adds the cells by the values of one dimension, named by
// @dimension, under the filter on every other dimension: the dimension's
// own filter is left out, so a reader sees what each of its values holds
// under the rest of the selection. The value is text; a service month
// reads YYYY-MM. Rows come in value order. The second count leaves out the
// extra rows, the only ones with no MAC claim.
const facetSQL = `
WITH picked AS (
    SELECT
        CASE @dimension::text
            WHEN 'jurisdiction' THEN jurisdiction
            WHEN 'part' THEN part
            WHEN 'workload' THEN mac_workload_id
            WHEN 'claim_type' THEN claim_type
            WHEN 'adjudication' THEN adjudication
            WHEN 'mac_status' THEN mac_claim_status_code
            WHEN 'vendor_status' THEN vendor_claim_status_code
            WHEN 'frequency' THEN claim_frequency_code
            WHEN 'paid_band' THEN paid_band
            WHEN 'service_month' THEN to_char(service_month, 'YYYY-MM')
            WHEN 'payment_band' THEN payment_band
            WHEN 'outcome' THEN outcome
        END AS value
        , outcome
        , claims
        , mac_paid
        , vendor_paid
        , billed
        , net_difference
        , gross_difference
        , same_status
    FROM public.cell
    WHERE submission_prefix = @submissionPrefix
        AND (@dimension::text = 'jurisdiction' OR cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
        AND (@dimension::text = 'part' OR cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
        AND (@dimension::text = 'workload' OR cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
        AND (@dimension::text = 'claim_type' OR cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
        AND (@dimension::text = 'adjudication' OR cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
        AND (@dimension::text = 'mac_status' OR cardinality(@macStatus::text[]) = 0 OR mac_claim_status_code = ANY (@macStatus::text[]))
        AND (@dimension::text = 'vendor_status' OR cardinality(@vendorStatus::text[]) = 0 OR vendor_claim_status_code = ANY (@vendorStatus::text[]))
        AND (@dimension::text = 'frequency' OR cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
        AND (@dimension::text = 'paid_band' OR cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
        AND (@dimension::text = 'service_month' OR cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))
        AND (@dimension::text = 'payment_band' OR cardinality(@paymentBand::text[]) = 0 OR payment_band = ANY (@paymentBand::text[]))
        AND (@dimension::text = 'outcome' OR cardinality(@outcome::text[]) = 0 OR outcome = ANY (@outcome::text[]))
)
SELECT
    value
    , sum(claims)::bigint
    , COALESCE(sum(claims) FILTER (WHERE outcome <> 'extra'), 0)::bigint
    , trim_scale(sum(mac_paid))::text
    , trim_scale(sum(vendor_paid))::text
    , trim_scale(sum(billed))::text
    , trim_scale(sum(net_difference))::text
    , trim_scale(sum(gross_difference))::text
    , sum(same_status)::bigint
FROM picked
GROUP BY value
ORDER BY value`

// Facet adds the cells by the values of one dimension (a Dimension
// constant) under the filter on the other dimensions. An unknown dimension
// is an error.
func (t TableAccess) Facet(ctx context.Context, f Filter, dimension string) ([]ValueSums, error) {
	known := false
	for _, d := range Dimensions {
		if d == dimension {
			known = true
		}
	}
	if !known {
		return nil, fmt.Errorf("cell: adding by dimension: unknown dimension %q", dimension)
	}
	rows, err := t.db.Query(ctx, facetSQL, args(f, dimension))
	if err != nil {
		return nil, fmt.Errorf("cell: adding by %s: %w", dimension, err)
	}
	defer rows.Close()
	var out []ValueSums
	for rows.Next() {
		var v ValueSums
		err := rows.Scan(
			&v.Value,
			&v.Claims,
			&v.MACClaims,
			&v.MACPaid,
			&v.VendorPaid,
			&v.Billed,
			&v.NetDifference,
			&v.GrossDifference,
			&v.SameStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("cell: scanning %s: %w", dimension, err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cell: adding by %s: %w", dimension, err)
	}
	return out, nil
}

// statusMatrixSQL counts the pairs (outcomes equal and different) the
// filter keeps by MAC status and vendor status, in status order.
const statusMatrixSQL = `
SELECT
    mac_claim_status_code
    , vendor_claim_status_code
    , sum(claims)::bigint
    , trim_scale(sum(mac_paid))::text
    , trim_scale(sum(vendor_paid))::text
FROM public.cell
WHERE submission_prefix = @submissionPrefix
    AND outcome IN ('equal', 'different')
    AND (cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
    AND (cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
    AND (cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
    AND (cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
    AND (cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
    AND (cardinality(@macStatus::text[]) = 0 OR mac_claim_status_code = ANY (@macStatus::text[]))
    AND (cardinality(@vendorStatus::text[]) = 0 OR vendor_claim_status_code = ANY (@vendorStatus::text[]))
    AND (cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
    AND (cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
    AND (cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))
    AND (cardinality(@paymentBand::text[]) = 0 OR payment_band = ANY (@paymentBand::text[]))
    AND (cardinality(@outcome::text[]) = 0 OR outcome = ANY (@outcome::text[]))
GROUP BY mac_claim_status_code, vendor_claim_status_code
ORDER BY mac_claim_status_code, vendor_claim_status_code`

// StatusMatrix counts the pairs the filter keeps by MAC status and vendor
// status.
func (t TableAccess) StatusMatrix(ctx context.Context, f Filter) ([]StatusCell, error) {
	rows, err := t.db.Query(ctx, statusMatrixSQL, args(f, ""))
	if err != nil {
		return nil, fmt.Errorf("cell: counting statuses: %w", err)
	}
	defer rows.Close()
	var out []StatusCell
	for rows.Next() {
		var c StatusCell
		if err := rows.Scan(&c.MACStatus, &c.VendorStatus, &c.Claims, &c.MACPaid, &c.VendorPaid); err != nil {
			return nil, fmt.Errorf("cell: scanning status: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cell: counting statuses: %w", err)
	}
	return out, nil
}
