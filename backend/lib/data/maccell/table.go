package maccell

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for the legacy cube.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const deleteBySubmissionSQL = `DELETE FROM public.mac_cell WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("maccell: deleting submission: %w", err)
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
		"frequency":        list(f.Frequency),
		"paidBand":         list(f.PaidBand),
		"serviceMonth":     dates(f.ServiceMonth),
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
const totalsSQL = `
SELECT
    COALESCE(sum(claims), 0)::bigint
    , trim_scale(COALESCE(sum(paid), 0))::text
    , trim_scale(COALESCE(sum(billed), 0))::text
FROM public.mac_cell
WHERE submission_prefix = @submissionPrefix
    AND (cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
    AND (cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
    AND (cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
    AND (cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
    AND (cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
    AND (cardinality(@macStatus::text[]) = 0 OR claim_status_code = ANY (@macStatus::text[]))
    AND (cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
    AND (cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
    AND (cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))`

// Totals adds every cell the filter keeps.
func (t TableAccess) Totals(ctx context.Context, f Filter) (Sums, error) {
	var s Sums
	if err := t.db.QueryRow(ctx, totalsSQL, args(f, "")).Scan(&s.Claims, &s.Paid, &s.Billed); err != nil {
		return Sums{}, fmt.Errorf("maccell: adding totals: %w", err)
	}
	return s, nil
}

// facetSQL adds the cells by the values of one dimension, named by
// @dimension, under the filter on every other dimension: the dimension's
// own filter is left out, so a reader sees what each of its values holds
// under the rest of the selection. The value is text; a service month
// reads YYYY-MM. Rows come in value order.
const facetSQL = `
WITH picked AS (
    SELECT
        CASE @dimension::text
            WHEN 'jurisdiction' THEN jurisdiction
            WHEN 'part' THEN part
            WHEN 'workload' THEN mac_workload_id
            WHEN 'claim_type' THEN claim_type
            WHEN 'adjudication' THEN adjudication
            WHEN 'mac_status' THEN claim_status_code
            WHEN 'frequency' THEN claim_frequency_code
            WHEN 'paid_band' THEN paid_band
            WHEN 'service_month' THEN to_char(service_month, 'YYYY-MM')
        END AS value
        , claims
        , paid
        , billed
    FROM public.mac_cell
    WHERE submission_prefix = @submissionPrefix
        AND (@dimension::text = 'jurisdiction' OR cardinality(@jurisdiction::text[]) = 0 OR jurisdiction = ANY (@jurisdiction::text[]))
        AND (@dimension::text = 'part' OR cardinality(@part::text[]) = 0 OR part = ANY (@part::text[]))
        AND (@dimension::text = 'workload' OR cardinality(@workload::text[]) = 0 OR mac_workload_id = ANY (@workload::text[]))
        AND (@dimension::text = 'claim_type' OR cardinality(@claimType::text[]) = 0 OR claim_type = ANY (@claimType::text[]))
        AND (@dimension::text = 'adjudication' OR cardinality(@adjudication::text[]) = 0 OR adjudication = ANY (@adjudication::text[]))
        AND (@dimension::text = 'mac_status' OR cardinality(@macStatus::text[]) = 0 OR claim_status_code = ANY (@macStatus::text[]))
        AND (@dimension::text = 'frequency' OR cardinality(@frequency::text[]) = 0 OR claim_frequency_code = ANY (@frequency::text[]))
        AND (@dimension::text = 'paid_band' OR cardinality(@paidBand::text[]) = 0 OR paid_band = ANY (@paidBand::text[]))
        AND (@dimension::text = 'service_month' OR cardinality(@serviceMonth::date[]) = 0 OR service_month = ANY (@serviceMonth::date[]))
)
SELECT
    value
    , sum(claims)::bigint
    , trim_scale(sum(paid))::text
    , trim_scale(sum(billed))::text
FROM picked
GROUP BY value
ORDER BY value`

// Facet adds the cells by the values of one dimension (a Dimension
// constant) under the filter on the other dimensions, in value order. An
// unknown dimension is an error.
func (t TableAccess) Facet(ctx context.Context, f Filter, dimension string) ([]ValueSums, error) {
	known := false
	for _, d := range Dimensions {
		if d == dimension {
			known = true
		}
	}
	if !known {
		return nil, fmt.Errorf("maccell: adding by dimension: unknown dimension %q", dimension)
	}
	rows, err := t.db.Query(ctx, facetSQL, args(f, dimension))
	if err != nil {
		return nil, fmt.Errorf("maccell: adding by %s: %w", dimension, err)
	}
	defer rows.Close()
	var out []ValueSums
	for rows.Next() {
		var v ValueSums
		if err := rows.Scan(&v.Value, &v.Claims, &v.Paid, &v.Billed); err != nil {
			return nil, fmt.Errorf("maccell: scanning %s: %w", dimension, err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("maccell: adding by %s: %w", dimension, err)
	}
	return out, nil
}
