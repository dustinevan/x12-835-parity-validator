package submission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for submissions.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const createSQL = `
INSERT INTO public.submission (
    prefix
    , side
    , label
    , submitted_on
    , status
    , claim_count
    , loaded_at
) VALUES (
    @prefix, @side, @label, @submittedOn, @status, @claimCount, @loadedAt
)`

// Create inserts a submission with status 'loaded'.
func (t TableAccess) Create(ctx context.Context, ns NewSubmission) error {
	var submittedOn pgtype.Date
	if !ns.SubmittedOn.IsZero() {
		submittedOn = pgtype.Date{Time: ns.SubmittedOn, Valid: true}
	}
	_, err := t.db.Exec(ctx, createSQL, pgx.NamedArgs{
		"prefix":      ns.Prefix,
		"side":        ns.Side,
		"label":       ns.Label,
		"submittedOn": submittedOn,
		"status":      StatusLoaded,
		"claimCount":  ns.ClaimCount,
		"loadedAt":    time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("submission: creating: %w", err)
	}
	return nil
}

const deleteSQL = `DELETE FROM public.submission WHERE prefix = @prefix`

func (t TableAccess) Delete(ctx context.Context, prefix string) error {
	if _, err := t.db.Exec(ctx, deleteSQL, pgx.NamedArgs{"prefix": prefix}); err != nil {
		return fmt.Errorf("submission: deleting: %w", err)
	}
	return nil
}

const setComparedSQL = `
UPDATE public.submission
SET
    status = @status
    , compared_at = @comparedAt
    , compared_set = @comparedSet
WHERE
    prefix = @prefix`

// SetCompared marks a vendor submission compared, now, under the compared
// set named comparedSet (service.CompareSet.Name).
func (t TableAccess) SetCompared(ctx context.Context, prefix, comparedSet string) error {
	_, err := t.db.Exec(ctx, setComparedSQL, pgx.NamedArgs{
		"prefix":      prefix,
		"status":      StatusCompared,
		"comparedAt":  time.Now().UTC(),
		"comparedSet": comparedSet,
	})
	if err != nil {
		return fmt.Errorf("submission: setting compared: %w", err)
	}
	return nil
}

const clearSummarizedSQL = `
UPDATE public.submission
SET
    status = @compared
WHERE
    prefix = @prefix
    AND status = @summarized`

// ClearSummarized sets a summarized submission back to compared, which hides
// it from the views while its summaries are rebuilt. Any other status is
// left as it is.
func (t TableAccess) ClearSummarized(ctx context.Context, prefix string) error {
	_, err := t.db.Exec(ctx, clearSummarizedSQL, pgx.NamedArgs{
		"prefix":     prefix,
		"compared":   StatusCompared,
		"summarized": StatusSummarized,
	})
	if err != nil {
		return fmt.Errorf("submission: clearing summarized: %w", err)
	}
	return nil
}

const clearComparedSQL = `
UPDATE public.submission
SET
    status = @status
    , compared_at = NULL
    , compared_set = NULL
WHERE
    prefix = @prefix`

// ClearCompared sets a vendor submission back to loaded, with no compared
// set, before its comparison rows are rewritten. The views hide it and
// Summarize refuses it until Compare finishes; a run that dies leaves it
// there.
func (t TableAccess) ClearCompared(ctx context.Context, prefix string) error {
	_, err := t.db.Exec(ctx, clearComparedSQL, pgx.NamedArgs{
		"prefix": prefix,
		"status": StatusLoaded,
	})
	if err != nil {
		return fmt.Errorf("submission: clearing compared: %w", err)
	}
	return nil
}

const setSummarizedSQL = `
UPDATE public.submission
SET
    status = @status
WHERE
    prefix = @prefix`

// SetSummarized marks a vendor submission summarized, which shows it in the
// views.
func (t TableAccess) SetSummarized(ctx context.Context, prefix string) error {
	_, err := t.db.Exec(ctx, setSummarizedSQL, pgx.NamedArgs{
		"prefix": prefix,
		"status": StatusSummarized,
	})
	if err != nil {
		return fmt.Errorf("submission: setting summarized: %w", err)
	}
	return nil
}

const setFailedSQL = `
UPDATE public.submission
SET
    status = @status
WHERE
    prefix = @prefix`

// SetFailed marks a submission failed: a stage could not finish and left
// its rows partly written. The views hide it, Summarize refuses it, and
// rerunning the stage clears it.
func (t TableAccess) SetFailed(ctx context.Context, prefix string) error {
	_, err := t.db.Exec(ctx, setFailedSQL, pgx.NamedArgs{
		"prefix": prefix,
		"status": StatusFailed,
	})
	if err != nil {
		return fmt.Errorf("submission: setting failed: %w", err)
	}
	return nil
}

const getSQL = `
SELECT
    prefix
    , side
    , label
    , submitted_on
    , status
    , claim_count
    , loaded_at
    , compared_at
    , COALESCE(compared_set, '')
FROM public.submission
WHERE prefix = @prefix`

// Get returns db.ErrNotFound when no submission has the prefix.
func (t TableAccess) Get(ctx context.Context, prefix string) (Submission, error) {
	var d Submission
	var submittedOn pgtype.Date
	var comparedAt pgtype.Timestamptz
	err := t.db.QueryRow(ctx, getSQL, pgx.NamedArgs{"prefix": prefix}).Scan(
		&d.Prefix,
		&d.Side,
		&d.Label,
		&submittedOn,
		&d.Status,
		&d.ClaimCount,
		&d.LoadedAt,
		&comparedAt,
		&d.ComparedSet,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Submission{}, fmt.Errorf("submission: getting %q: %w", prefix, db.ErrNotFound)
	}
	if err != nil {
		return Submission{}, fmt.Errorf("submission: getting: %w", err)
	}
	if submittedOn.Valid {
		d.SubmittedOn = submittedOn.Time
	}
	if comparedAt.Valid {
		d.ComparedAt = comparedAt.Time.UTC()
	}
	d.LoadedAt = d.LoadedAt.UTC()
	return d, nil
}

const macSQL = `
SELECT
    prefix
    , side
    , label
    , submitted_on
    , status
    , claim_count
    , loaded_at
    , compared_at
    , COALESCE(compared_set, '')
FROM public.submission
WHERE side = 'mac'`

// MAC returns the one MAC submission. None is db.ErrNotFound; more than one
// is an error, because comparisons would not know which to use.
func (t TableAccess) MAC(ctx context.Context) (Submission, error) {
	rows, err := t.db.Query(ctx, macSQL)
	if err != nil {
		return Submission{}, fmt.Errorf("submission: getting the MAC submission: %w", err)
	}
	defer rows.Close()
	var found []Submission
	for rows.Next() {
		var d Submission
		var submittedOn pgtype.Date
		var comparedAt pgtype.Timestamptz
		err := rows.Scan(
			&d.Prefix,
			&d.Side,
			&d.Label,
			&submittedOn,
			&d.Status,
			&d.ClaimCount,
			&d.LoadedAt,
			&comparedAt,
			&d.ComparedSet,
		)
		if err != nil {
			return Submission{}, fmt.Errorf("submission: getting the MAC submission: %w", err)
		}
		if submittedOn.Valid {
			d.SubmittedOn = submittedOn.Time
		}
		if comparedAt.Valid {
			d.ComparedAt = comparedAt.Time.UTC()
		}
		d.LoadedAt = d.LoadedAt.UTC()
		found = append(found, d)
	}
	if err := rows.Err(); err != nil {
		return Submission{}, fmt.Errorf("submission: getting the MAC submission: %w", err)
	}
	switch len(found) {
	case 0:
		return Submission{}, fmt.Errorf("submission: getting the MAC submission: %w", db.ErrNotFound)
	case 1:
		return found[0], nil
	}
	return Submission{}, fmt.Errorf("submission: getting the MAC submission: %d submissions have side 'mac'", len(found))
}

const listSQL = `
SELECT
    prefix
    , side
    , label
    , submitted_on
    , status
    , claim_count
    , loaded_at
    , compared_at
    , COALESCE(compared_set, '')
FROM public.submission
ORDER BY
    side = 'mac' DESC
    , label
    , submitted_on
    , prefix`

// List returns every submission, whatever its status: the MAC first, then
// vendors by label and date.
func (t TableAccess) List(ctx context.Context) ([]Submission, error) {
	rows, err := t.db.Query(ctx, listSQL)
	if err != nil {
		return nil, fmt.Errorf("submission: listing: %w", err)
	}
	defer rows.Close()
	var out []Submission
	for rows.Next() {
		var d Submission
		var submittedOn pgtype.Date
		var comparedAt pgtype.Timestamptz
		err := rows.Scan(
			&d.Prefix,
			&d.Side,
			&d.Label,
			&submittedOn,
			&d.Status,
			&d.ClaimCount,
			&d.LoadedAt,
			&comparedAt,
			&d.ComparedSet,
		)
		if err != nil {
			return nil, fmt.Errorf("submission: listing: %w", err)
		}
		if submittedOn.Valid {
			d.SubmittedOn = submittedOn.Time
		}
		if comparedAt.Valid {
			d.ComparedAt = comparedAt.Time.UTC()
		}
		d.LoadedAt = d.LoadedAt.UTC()
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("submission: listing: %w", err)
	}
	return out, nil
}
