package difference

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for differences.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

// CreateMany bulk-loads differences with COPY.
func (t TableAccess) CreateMany(ctx context.Context, differences []Difference) (int64, error) {
	columns := []string{
		"submission_prefix",
		"match_key",
		"element",
		"occurrence",
		"line_control_number",
		"mac_value",
		"vendor_value",
		"amount_difference",
	}
	rows := pgx.CopyFromSlice(len(differences), func(i int) ([]any, error) {
		d := differences[i]
		var amount pgtype.Numeric
		if d.AmountDifference != "" {
			if err := amount.Scan(d.AmountDifference); err != nil {
				return nil, fmt.Errorf("parsing amount difference %q: %w", d.AmountDifference, err)
			}
		}
		return []any{
			d.SubmissionPrefix,
			d.MatchKey,
			d.Element,
			d.Occurrence,
			text(d.LineControlNumber),
			text(d.MACValue),
			text(d.VendorValue),
			amount,
		}, nil
	})
	n, err := t.db.CopyFrom(ctx, pgx.Identifier{"public", "difference"}, columns, rows)
	if err != nil {
		return 0, fmt.Errorf("difference: copying: %w", err)
	}
	return n, nil
}

const deleteBySubmissionSQL = `DELETE FROM public.difference WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("difference: deleting submission: %w", err)
	}
	return nil
}

const listByMatchKeySQL = `
SELECT
    submission_prefix
    , match_key
    , element
    , occurrence
    , COALESCE(line_control_number, '')
    , COALESCE(mac_value, '')
    , COALESCE(vendor_value, '')
    , COALESCE(amount_difference::text, '')
FROM public.difference
WHERE submission_prefix = @submissionPrefix
    AND match_key = @matchKey
ORDER BY element, occurrence`

// ListByMatchKey returns one compared pair's differences.
func (t TableAccess) ListByMatchKey(ctx context.Context, submissionPrefix string, matchKey []byte) ([]Difference, error) {
	rows, err := t.db.Query(ctx, listByMatchKeySQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "matchKey": matchKey})
	if err != nil {
		return nil, fmt.Errorf("difference: listing: %w", err)
	}
	defer rows.Close()
	var out []Difference
	for rows.Next() {
		var d Difference
		err := rows.Scan(&d.SubmissionPrefix, &d.MatchKey, &d.Element, &d.Occurrence, &d.LineControlNumber, &d.MACValue, &d.VendorValue, &d.AmountDifference)
		if err != nil {
			return nil, fmt.Errorf("difference: scanning: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("difference: listing: %w", err)
	}
	return out, nil
}

// text is NULL for an empty string.
func text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}
