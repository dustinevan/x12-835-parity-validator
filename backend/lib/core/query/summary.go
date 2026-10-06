// Package query holds SQL that reads more than one table. Each method has
// its own complete statement directly above it and runs it itself; nothing
// SQL is shared between methods.
package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/core/service"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// SummaryQuery is the summarize stage's SQL: it writes measure, tag_total,
// cell, and tag_cell from the compare stage's comparison and difference
// rows (and, for measure, the MAC submission's claim rows). Each method is
// one INSERT ... SELECT. The caller (lib/core/control/submission) deletes the
// rows being replaced first, in the same short transaction.
type SummaryQuery struct {
	db db.Console
}

func NewSummaryQuery(db db.Console) SummaryQuery {
	return SummaryQuery{db: db}
}

// queryCreateMeasures writes one row per grouping: all/all, per claim
// type, per workload, and per claim type and workload.
//
// A set-aside comparison row stands for every MAC claim sharing its key,
// so it counts those claims (from the MAC submission's claim rows). A
// duplicated row stands for its one MAC claim, like a missing row: the
// vendor's copies are not MAC claims. The outcome match rate divides by
// MAC claims not set aside, so a duplicated key stays in the denominator
// and, being unmatched, counts against the vendor.
// Deviation is |vendor paid - MAC paid| / MAC paid over matched pairs
// where the MAC paid more than zero. Its median and 95th percentile use
// percentile_disc on NUMERIC: an observed value, computed exactly, with no
// floating-point interpolation (percentile_cont works in double precision).
// trim_scale gives equal ratios one spelling, so ties pick the same text.
const queryCreateMeasures = `
WITH mac_claims_per_key AS (
    SELECT match_key, count(*) AS claims
    FROM public.claim
    WHERE submission_prefix = @macPrefix
    GROUP BY match_key
),
result AS (
    SELECT
        c.claim_type
        , c.mac_workload_id
        , c.outcome
        , CASE c.outcome
            WHEN 'extra' THEN 0
            WHEN 'set_aside' THEN COALESCE(k.claims, 1)
            ELSE 1
          END AS mac_claims
        , c.outcome IN ('equal', 'different') AS matched
        , c.mac_claim_status_code = c.vendor_claim_status_code AS same_status
        , CASE WHEN c.outcome IN ('equal', 'different') AND c.mac_claim_payment_amount > 0
            THEN trim_scale(abs(c.vendor_claim_payment_amount - c.mac_claim_payment_amount) / c.mac_claim_payment_amount)
          END AS deviation
    FROM public.comparison c
    LEFT JOIN mac_claims_per_key k
        ON c.outcome = 'set_aside' AND k.match_key = c.match_key
    WHERE c.submission_prefix = @submissionPrefix
),
totals AS (
    SELECT
        CASE WHEN GROUPING(claim_type) = 1 THEN 'all' ELSE claim_type END AS claim_type
        , CASE WHEN GROUPING(mac_workload_id) = 1 THEN 'all' ELSE mac_workload_id END AS mac_workload_id
        , sum(mac_claims)::bigint AS mac_claims
        , count(*) FILTER (WHERE outcome = 'equal') AS equal_count
        , count(*) FILTER (WHERE outcome = 'different') AS different_count
        , count(*) FILTER (WHERE outcome = 'missing') AS missing_count
        , count(*) FILTER (WHERE outcome = 'extra') AS extra_count
        , COALESCE(sum(mac_claims) FILTER (WHERE outcome = 'set_aside'), 0)::bigint AS set_aside_count
        , count(*) FILTER (WHERE outcome = 'duplicated') AS duplicated_count
        , count(*) FILTER (WHERE matched AND same_status) AS same_status_count
        , count(deviation) AS paid_claims
        , percentile_disc(0.5) WITHIN GROUP (ORDER BY deviation) AS median_deviation
        , percentile_disc(0.95) WITHIN GROUP (ORDER BY deviation) AS p95_deviation
    FROM result
    GROUP BY GROUPING SETS ((), (claim_type), (mac_workload_id), (claim_type, mac_workload_id))
)
INSERT INTO public.measure (
    submission_prefix
    , claim_type
    , mac_workload_id
    , mac_claims
    , equal_count
    , different_count
    , missing_count
    , extra_count
    , set_aside_count
    , duplicated_count
    , outcome_match_rate
    , paid_claims
    , median_deviation
    , p95_deviation
)
SELECT
    @submissionPrefix
    , claim_type
    , mac_workload_id
    , mac_claims
    , equal_count
    , different_count
    , missing_count
    , extra_count
    , set_aside_count
    , duplicated_count
    , same_status_count::numeric / NULLIF(mac_claims - set_aside_count, 0)
    , paid_claims
    , median_deviation
    , p95_deviation
FROM totals`

// CreateMeasures writes a vendor submission's measure rows from its
// comparison rows.
func (q SummaryQuery) CreateMeasures(ctx context.Context, submissionPrefix, macPrefix string) error {
	_, err := q.db.Exec(ctx, queryCreateMeasures, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "macPrefix": macPrefix})
	if err != nil {
		return fmt.Errorf("query: creating measures: %w", err)
	}
	return nil
}

// queryCreateTagTotals counts, per tag, the different claims that carry
// it, and sums |amount_difference| of that element's difference rows in
// those claims, for the same four groupings as measure. A tag with no
// amount (a code or a date) has dollars 0. Claims and dollars are grouped
// separately and joined, so each group sums only its own rows.
const queryCreateTagTotals = `
WITH picked AS (
    SELECT c.match_key, c.claim_type, c.mac_workload_id, c.tags
    FROM public.comparison c
    WHERE c.submission_prefix = @submissionPrefix
        AND c.outcome = 'different'
),
claims AS (
    SELECT
        CASE WHEN GROUPING(picked.claim_type) = 1 THEN 'all' ELSE picked.claim_type END AS claim_type
        , CASE WHEN GROUPING(picked.mac_workload_id) = 1 THEN 'all' ELSE picked.mac_workload_id END AS mac_workload_id
        , tag
        , count(*) AS claims
    FROM picked, unnest(picked.tags) AS tag
    GROUP BY GROUPING SETS (
        (tag)
        , (picked.claim_type, tag)
        , (picked.mac_workload_id, tag)
        , (picked.claim_type, picked.mac_workload_id, tag)
    )
),
dollars AS (
    SELECT
        CASE WHEN GROUPING(p.claim_type) = 1 THEN 'all' ELSE p.claim_type END AS claim_type
        , CASE WHEN GROUPING(p.mac_workload_id) = 1 THEN 'all' ELSE p.mac_workload_id END AS mac_workload_id
        , d.element AS tag
        , sum(abs(d.amount_difference)) AS dollars
    FROM public.difference d
    JOIN picked p ON p.match_key = d.match_key
    WHERE d.submission_prefix = @submissionPrefix
    GROUP BY GROUPING SETS (
        (d.element)
        , (p.claim_type, d.element)
        , (p.mac_workload_id, d.element)
        , (p.claim_type, p.mac_workload_id, d.element)
    )
)
INSERT INTO public.tag_total (
    submission_prefix
    , claim_type
    , mac_workload_id
    , tag
    , claims
    , dollars
)
SELECT
    @submissionPrefix
    , claims.claim_type
    , claims.mac_workload_id
    , claims.tag
    , claims.claims
    , COALESCE(dollars.dollars, 0)
FROM claims
LEFT JOIN dollars
    ON dollars.claim_type = claims.claim_type
    AND dollars.mac_workload_id = claims.mac_workload_id
    AND dollars.tag = claims.tag`

// CreateTagTotals writes a vendor submission's tag_total rows from its
// comparison and difference rows.
func (q SummaryQuery) CreateTagTotals(ctx context.Context, submissionPrefix string) error {
	_, err := q.db.Exec(ctx, queryCreateTagTotals, pgx.NamedArgs{"submissionPrefix": submissionPrefix})
	if err != nil {
		return fmt.Errorf("query: creating tag totals: %w", err)
	}
	return nil
}

// queryCreateCells writes the summary cube: one row per distinct
// combination of the twelve dimensions among the submission's comparison
// rows. The bands are worked out here. paid_band is the MAC paid amount,
// or the vendor's for an extra row: zero, then under 100, 1000, and 10000,
// then 10000 and more. payment_band is |vendor paid - MAC paid| on a pair
// (equal or different): match when it is within @matchCents cents, the
// comparison's five-cent rule (service.AmountMatchCents), edge included;
// up to 1.00; up to 1% and up to 5% of the MAC paid amount when it is more
// than zero; else over 5%; and the empty value on anything but a pair. An
// absent amount adds as zero; a NULL status reads as the empty value.
// net_difference and gross_difference add pairs only, and same_status
// counts pairs whose CLP02 agree.
const queryCreateCells = `
WITH picked AS (
    SELECT
        c.jurisdiction
        , c.part
        , c.mac_workload_id
        , c.claim_type
        , c.adjudication
        , COALESCE(c.mac_claim_status_code, '') AS mac_status
        , COALESCE(c.vendor_claim_status_code, '') AS vendor_status
        , c.claim_frequency_code
        , date_trunc('month', c.date_of_service)::date AS service_month
        , c.outcome
        , c.outcome IN ('equal', 'different') AS pair
        , COALESCE(c.mac_claim_payment_amount, 0) AS mac_paid
        , COALESCE(c.vendor_claim_payment_amount, 0) AS vendor_paid
        , c.total_claim_charge_amount AS billed
        , c.mac_claim_status_code = c.vendor_claim_status_code AS same_status
        , CASE WHEN c.outcome = 'extra'
            THEN COALESCE(c.vendor_claim_payment_amount, 0)
            ELSE COALESCE(c.mac_claim_payment_amount, 0)
          END AS banded_paid
    FROM public.comparison c
    WHERE c.submission_prefix = @submissionPrefix
),
banded AS (
    SELECT
        *
        , CASE
            WHEN banded_paid = 0 THEN 'zero'
            WHEN banded_paid < 100 THEN 'under 100'
            WHEN banded_paid < 1000 THEN 'under 1000'
            WHEN banded_paid < 10000 THEN 'under 10000'
            ELSE '10000 and more'
          END AS paid_band
        , CASE
            WHEN NOT pair THEN ''
            WHEN abs(vendor_paid - mac_paid) <= @matchCents::numeric / 100 THEN 'match'
            WHEN abs(vendor_paid - mac_paid) <= 1 THEN 'up to 1.00'
            WHEN mac_paid > 0 AND abs(vendor_paid - mac_paid) <= mac_paid * 0.01 THEN 'up to 1%'
            WHEN mac_paid > 0 AND abs(vendor_paid - mac_paid) <= mac_paid * 0.05 THEN 'up to 5%'
            ELSE 'over 5%'
          END AS payment_band
    FROM picked
)
INSERT INTO public.cell (
    submission_prefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_claim_status_code
    , vendor_claim_status_code
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome
    , claims
    , mac_paid
    , vendor_paid
    , billed
    , net_difference
    , gross_difference
    , same_status
)
SELECT
    @submissionPrefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_status
    , vendor_status
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome
    , count(*)
    , sum(mac_paid)
    , sum(vendor_paid)
    , sum(billed)
    , sum(CASE WHEN pair THEN vendor_paid - mac_paid ELSE 0 END)
    , sum(CASE WHEN pair THEN abs(vendor_paid - mac_paid) ELSE 0 END)
    , count(*) FILTER (WHERE pair AND same_status)
FROM banded
GROUP BY
    jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_status
    , vendor_status
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome`

// CreateCells writes a vendor submission's cell rows from its comparison
// rows.
func (q SummaryQuery) CreateCells(ctx context.Context, submissionPrefix string) error {
	_, err := q.db.Exec(ctx, queryCreateCells, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "matchCents": service.AmountMatchCents})
	if err != nil {
		return fmt.Errorf("query: creating cells: %w", err)
	}
	return nil
}

// queryCreateTagCells writes the tag cube from the submission's different
// pairs: the same dimensions and bands as queryCreateCells, then one row
// per cell and tag, counting the pairs whose tags contain the tag and adding
// the dollars of that element in them (the sum of |amount_difference| of its
// difference rows; a tag with no amount, a code or a date, adds 0). The
// dollars are added per pair and element first, then joined to each pair's
// tags, so a pair with several difference rows of one element adds them
// once.
const queryCreateTagCells = `
WITH picked AS (
    SELECT
        c.match_key
        , c.tags
        , c.jurisdiction
        , c.part
        , c.mac_workload_id
        , c.claim_type
        , c.adjudication
        , COALESCE(c.mac_claim_status_code, '') AS mac_status
        , COALESCE(c.vendor_claim_status_code, '') AS vendor_status
        , c.claim_frequency_code
        , date_trunc('month', c.date_of_service)::date AS service_month
        , c.outcome
        , COALESCE(c.mac_claim_payment_amount, 0) AS mac_paid
        , COALESCE(c.vendor_claim_payment_amount, 0) AS vendor_paid
    FROM public.comparison c
    WHERE c.submission_prefix = @submissionPrefix
        AND c.outcome = 'different'
),
banded AS (
    SELECT
        *
        , CASE
            WHEN mac_paid = 0 THEN 'zero'
            WHEN mac_paid < 100 THEN 'under 100'
            WHEN mac_paid < 1000 THEN 'under 1000'
            WHEN mac_paid < 10000 THEN 'under 10000'
            ELSE '10000 and more'
          END AS paid_band
        , CASE
            WHEN abs(vendor_paid - mac_paid) <= @matchCents::numeric / 100 THEN 'match'
            WHEN abs(vendor_paid - mac_paid) <= 1 THEN 'up to 1.00'
            WHEN mac_paid > 0 AND abs(vendor_paid - mac_paid) <= mac_paid * 0.01 THEN 'up to 1%'
            WHEN mac_paid > 0 AND abs(vendor_paid - mac_paid) <= mac_paid * 0.05 THEN 'up to 5%'
            ELSE 'over 5%'
          END AS payment_band
    FROM picked
),
element_dollars AS (
    SELECT d.match_key, d.element, sum(abs(d.amount_difference)) AS dollars
    FROM public.difference d
    WHERE d.submission_prefix = @submissionPrefix
    GROUP BY d.match_key, d.element
),
tagged AS (
    SELECT b.*, t.tag, COALESCE(e.dollars, 0) AS dollars
    FROM banded b
    CROSS JOIN LATERAL unnest(b.tags) AS t(tag)
    LEFT JOIN element_dollars e
        ON e.match_key = b.match_key
        AND e.element = t.tag
)
INSERT INTO public.tag_cell (
    submission_prefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_claim_status_code
    , vendor_claim_status_code
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome
    , tag
    , claims
    , dollars
)
SELECT
    @submissionPrefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_status
    , vendor_status
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome
    , tag
    , count(*)
    , sum(dollars)
FROM tagged
GROUP BY
    jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , mac_status
    , vendor_status
    , claim_frequency_code
    , paid_band
    , service_month
    , payment_band
    , outcome
    , tag`

// CreateTagCells writes a vendor submission's tag_cell rows from its
// comparison and difference rows.
func (q SummaryQuery) CreateTagCells(ctx context.Context, submissionPrefix string) error {
	_, err := q.db.Exec(ctx, queryCreateTagCells, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "matchCents": service.AmountMatchCents})
	if err != nil {
		return fmt.Errorf("query: creating tag cells: %w", err)
	}
	return nil
}
