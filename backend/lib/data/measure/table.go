package measure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for measures.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const deleteBySubmissionSQL = `DELETE FROM public.measure WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("measure: deleting submission: %w", err)
	}
	return nil
}

const getSQL = `
SELECT
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
    , COALESCE(outcome_match_rate::text, '')
    , paid_claims
    , COALESCE(median_deviation::text, '')
    , COALESCE(p95_deviation::text, '')
FROM public.measure
WHERE submission_prefix = @submissionPrefix
    AND claim_type = @claimType
    AND mac_workload_id = @workload`

// Get returns one row: All for claimType or workload is the total over
// it. It returns db.ErrNotFound when the submission has no claims of that
// claim type and workload.
func (t TableAccess) Get(ctx context.Context, submissionPrefix, claimType, workload string) (Measure, error) {
	var m Measure
	err := t.db.QueryRow(ctx, getSQL, pgx.NamedArgs{
		"submissionPrefix": submissionPrefix,
		"claimType":        claimType,
		"workload":         workload,
	}).Scan(
		&m.SubmissionPrefix,
		&m.ClaimType,
		&m.MACWorkloadID,
		&m.MACClaims,
		&m.EqualCount,
		&m.DifferentCount,
		&m.MissingCount,
		&m.ExtraCount,
		&m.SetAsideCount,
		&m.DuplicatedCount,
		&m.OutcomeMatchRate,
		&m.PaidClaims,
		&m.MedianDeviation,
		&m.P95Deviation,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Measure{}, fmt.Errorf("measure: getting: %w", db.ErrNotFound)
	}
	if err != nil {
		return Measure{}, fmt.Errorf("measure: getting: %w", err)
	}
	return m, nil
}

const listBySubmissionSQL = `
SELECT
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
    , COALESCE(outcome_match_rate::text, '')
    , paid_claims
    , COALESCE(median_deviation::text, '')
    , COALESCE(p95_deviation::text, '')
FROM public.measure
WHERE submission_prefix = @submissionPrefix
ORDER BY claim_type <> 'all', claim_type, mac_workload_id <> 'all', mac_workload_id`

// ListBySubmission returns every measure row for a submission, totals first.
func (t TableAccess) ListBySubmission(ctx context.Context, submissionPrefix string) ([]Measure, error) {
	rows, err := t.db.Query(ctx, listBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix})
	if err != nil {
		return nil, fmt.Errorf("measure: listing: %w", err)
	}
	defer rows.Close()
	var out []Measure
	for rows.Next() {
		var m Measure
		err := rows.Scan(
			&m.SubmissionPrefix,
			&m.ClaimType,
			&m.MACWorkloadID,
			&m.MACClaims,
			&m.EqualCount,
			&m.DifferentCount,
			&m.MissingCount,
			&m.ExtraCount,
			&m.SetAsideCount,
			&m.DuplicatedCount,
			&m.OutcomeMatchRate,
			&m.PaidClaims,
			&m.MedianDeviation,
			&m.P95Deviation,
		)
		if err != nil {
			return nil, fmt.Errorf("measure: scanning: %w", err)
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("measure: listing: %w", err)
	}
	return out, nil
}
