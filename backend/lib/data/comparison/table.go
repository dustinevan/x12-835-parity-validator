package comparison

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for comparisons.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

// CreateMany bulk-loads comparisons with COPY.
func (t TableAccess) CreateMany(ctx context.Context, comparisons []Comparison) (int64, error) {
	columns := []string{
		"submission_prefix",
		"match_key",
		"outcome",
		"mac_workload_id",
		"claim_type",
		"medicare_beneficiary_identifier",
		"date_of_service",
		"total_claim_charge_amount",
		"mac_claim_status_code",
		"vendor_claim_status_code",
		"mac_claim_payment_amount",
		"vendor_claim_payment_amount",
		"claim_frequency_code",
		"adjudication",
		"jurisdiction",
		"part",
		"difference_count",
		"tags",
	}
	rows := pgx.CopyFromSlice(len(comparisons), func(i int) ([]any, error) {
		c := comparisons[i]
		charge, err := numeric(c.TotalClaimChargeAmount)
		if err != nil {
			return nil, err
		}
		macPaid, err := numeric(c.MACClaimPaymentAmount)
		if err != nil {
			return nil, err
		}
		vendorPaid, err := numeric(c.VendorClaimPaymentAmount)
		if err != nil {
			return nil, err
		}
		tags := c.Tags
		if tags == nil {
			tags = []string{}
		}
		return []any{
			c.SubmissionPrefix,
			c.MatchKey,
			c.Outcome,
			c.MACWorkloadID,
			c.ClaimType,
			c.MedicareBeneficiaryIdentifier,
			pgtype.Date{Time: c.DateOfService, Valid: true},
			charge,
			text(c.MACClaimStatusCode),
			text(c.VendorClaimStatusCode),
			macPaid,
			vendorPaid,
			c.ClaimFrequencyCode,
			c.Adjudication,
			c.Jurisdiction,
			c.Part,
			c.DifferenceCount,
			tags,
		}, nil
	})
	n, err := t.db.CopyFrom(ctx, pgx.Identifier{"public", "comparison"}, columns, rows)
	if err != nil {
		return 0, fmt.Errorf("comparison: copying: %w", err)
	}
	return n, nil
}

const deleteBySubmissionSQL = `DELETE FROM public.comparison WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("comparison: deleting submission: %w", err)
	}
	return nil
}

const getSQL = `
SELECT
    submission_prefix
    , match_key
    , outcome
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount::text
    , COALESCE(mac_claim_status_code, '')
    , COALESCE(vendor_claim_status_code, '')
    , COALESCE(mac_claim_payment_amount::text, '')
    , COALESCE(vendor_claim_payment_amount::text, '')
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , difference_count
    , tags
FROM public.comparison
WHERE submission_prefix = @submissionPrefix
    AND match_key = @matchKey`

// Get returns db.ErrNotFound when the submission has no row for the key.
func (t TableAccess) Get(ctx context.Context, submissionPrefix string, matchKey []byte) (Comparison, error) {
	var c Comparison
	err := t.db.QueryRow(ctx, getSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "matchKey": matchKey}).Scan(
		&c.SubmissionPrefix,
		&c.MatchKey,
		&c.Outcome,
		&c.MACWorkloadID,
		&c.ClaimType,
		&c.MedicareBeneficiaryIdentifier,
		&c.DateOfService,
		&c.TotalClaimChargeAmount,
		&c.MACClaimStatusCode,
		&c.VendorClaimStatusCode,
		&c.MACClaimPaymentAmount,
		&c.VendorClaimPaymentAmount,
		&c.ClaimFrequencyCode,
		&c.Adjudication,
		&c.Jurisdiction,
		&c.Part,
		&c.DifferenceCount,
		&c.Tags,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Comparison{}, fmt.Errorf("comparison: getting: %w", db.ErrNotFound)
	}
	if err != nil {
		return Comparison{}, fmt.Errorf("comparison: getting: %w", err)
	}
	return c, nil
}

// numeric parses an exact decimal string for COPY; empty is NULL.
func numeric(s string) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	if s == "" {
		return n, nil
	}
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("comparison: parsing decimal %q: %w", s, err)
	}
	return n, nil
}

// text is NULL for an empty string.
func text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// histogramSQL puts each paid delta on its side of zero (below, zero,
// above) and splits each side into @buckets equal-width buckets between
// that side's smallest and largest delta. A bucket reports the smallest
// and largest delta it actually holds, so its bounds are observed values,
// exact, with no rounding. Empty buckets are left out.
const histogramSQL = `
WITH picked AS (
    SELECT c.vendor_claim_payment_amount - c.mac_claim_payment_amount AS delta
    FROM public.comparison c
    WHERE c.submission_prefix = @submissionPrefix
        AND (@claimType::text = '' OR c.claim_type = @claimType::text)
        AND (@workload::text = '' OR c.mac_workload_id = @workload::text)
        AND (@outcome::text = '' OR c.outcome = @outcome::text)
        AND (@all::text[] IS NULL OR c.tags @> @all::text[])
        AND (@any::text[] IS NULL OR cardinality(@any::text[]) = 0 OR c.tags && @any::text[])
        AND (@none::text[] IS NULL OR NOT c.tags && @none::text[])
        AND (@otherOutcome::text = '' OR EXISTS (
            SELECT 1
            FROM public.comparison o
            WHERE o.submission_prefix = @otherPrefix::text
                AND o.match_key = c.match_key
                AND o.outcome = @otherOutcome::text
        ))
        AND c.mac_claim_payment_amount IS NOT NULL
        AND c.vendor_claim_payment_amount IS NOT NULL
),
sides AS (
    SELECT sign(delta) AS side, min(delta) AS low, max(delta) AS high
    FROM picked
    GROUP BY sign(delta)
),
bucketed AS (
    SELECT
        p.delta
        , s.side
        , CASE WHEN s.low = s.high THEN 1
            ELSE LEAST(width_bucket(p.delta, s.low, s.high, @buckets::int), @buckets::int)
          END AS bucket
    FROM picked p
    JOIN sides s ON s.side = sign(p.delta)
)
SELECT min(delta)::text, max(delta)::text, count(*)
FROM bucketed
GROUP BY side, bucket
ORDER BY min(delta)`

// Histogram buckets the paid deltas of the rows a Filter selects.
func (t TableAccess) Histogram(ctx context.Context, f Filter, bucketsPerSide int) ([]Bucket, error) {
	rows, err := t.db.Query(ctx, histogramSQL, pgx.NamedArgs{
		"submissionPrefix": f.SubmissionPrefix,
		"claimType":        f.ClaimType,
		"workload":         f.Workload,
		"outcome":          f.Outcome,
		"all":              f.All,
		"any":              f.Any,
		"none":             f.None,
		"otherPrefix":      f.OtherSubmissionPrefix,
		"otherOutcome":     f.OtherOutcome,
		"buckets":          bucketsPerSide,
	})
	if err != nil {
		return nil, fmt.Errorf("comparison: bucketing: %w", err)
	}
	defer rows.Close()
	var out []Bucket
	for rows.Next() {
		var b Bucket
		if err := rows.Scan(&b.From, &b.To, &b.Count); err != nil {
			return nil, fmt.Errorf("comparison: scanning bucket: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("comparison: bucketing: %w", err)
	}
	return out, nil
}

// totalsSQL counts the rows a Filter selects and adds up |vendor paid - MAC
// paid| over them; a row with either side absent adds nothing.
const totalsSQL = `
SELECT
    count(*)
    , COALESCE(sum(abs(c.vendor_claim_payment_amount - c.mac_claim_payment_amount)), 0)::text
FROM public.comparison c
WHERE c.submission_prefix = @submissionPrefix
    AND (@claimType::text = '' OR c.claim_type = @claimType::text)
    AND (@workload::text = '' OR c.mac_workload_id = @workload::text)
    AND (@outcome::text = '' OR c.outcome = @outcome::text)
    AND (@all::text[] IS NULL OR c.tags @> @all::text[])
    AND (@any::text[] IS NULL OR cardinality(@any::text[]) = 0 OR c.tags && @any::text[])
    AND (@none::text[] IS NULL OR NOT c.tags && @none::text[])
    AND (@otherOutcome::text = '' OR EXISTS (
        SELECT 1
        FROM public.comparison o
        WHERE o.submission_prefix = @otherPrefix::text
            AND o.match_key = c.match_key
            AND o.outcome = @otherOutcome::text
    ))`

// Totals counts the rows a Filter selects and sums their paid deltas.
func (t TableAccess) Totals(ctx context.Context, f Filter) (Totals, error) {
	var out Totals
	err := t.db.QueryRow(ctx, totalsSQL, pgx.NamedArgs{
		"submissionPrefix": f.SubmissionPrefix,
		"claimType":        f.ClaimType,
		"workload":         f.Workload,
		"outcome":          f.Outcome,
		"all":              f.All,
		"any":              f.Any,
		"none":             f.None,
		"otherPrefix":      f.OtherSubmissionPrefix,
		"otherOutcome":     f.OtherOutcome,
	}).Scan(&out.Claims, &out.PaidDelta)
	if err != nil {
		return Totals{}, fmt.Errorf("comparison: totaling: %w", err)
	}
	return out, nil
}

// agreementSQL joins two submissions' rows on match key and counts the keys
// by the pair of outcomes. Only equal, different, missing, and duplicated
// are counted: each such key stands for exactly one MAC claim, while a
// set-aside key stands for several and an extra key for none.
const agreementSQL = `
SELECT
    v1.outcome
    , v2.outcome
    , count(*)
    , count(*) FILTER (WHERE v1.outcome = 'different' AND v2.outcome = 'different' AND v1.tags = v2.tags)
FROM public.comparison v1
JOIN public.comparison v2
    ON v2.submission_prefix = @v2Prefix
    AND v2.match_key = v1.match_key
WHERE v1.submission_prefix = @v1Prefix
    AND v1.outcome IN ('equal', 'different', 'missing', 'duplicated')
    AND v2.outcome IN ('equal', 'different', 'missing', 'duplicated')
GROUP BY v1.outcome, v2.outcome
ORDER BY v1.outcome, v2.outcome`

// Agreement counts MAC claims by their outcome in the submission in slot v1
// and in the one in slot v2. Pairs with no claims are left out; the rows are
// ordered by the v1 outcome, then the v2 outcome.
func (t TableAccess) Agreement(ctx context.Context, v1Prefix, v2Prefix string) ([]AgreementCell, error) {
	rows, err := t.db.Query(ctx, agreementSQL, pgx.NamedArgs{"v1Prefix": v1Prefix, "v2Prefix": v2Prefix})
	if err != nil {
		return nil, fmt.Errorf("comparison: counting agreement: %w", err)
	}
	defer rows.Close()
	var out []AgreementCell
	for rows.Next() {
		var c AgreementCell
		if err := rows.Scan(&c.V1Outcome, &c.V2Outcome, &c.Claims, &c.SameTags); err != nil {
			return nil, fmt.Errorf("comparison: scanning agreement: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("comparison: counting agreement: %w", err)
	}
	return out, nil
}
