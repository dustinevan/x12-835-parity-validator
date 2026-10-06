package tagtotal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for tag totals.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const deleteBySubmissionSQL = `DELETE FROM public.tag_total WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("tagtotal: deleting submission: %w", err)
	}
	return nil
}

const listSQL = `
SELECT
    submission_prefix
    , claim_type
    , mac_workload_id
    , tag
    , claims
    , dollars::text
FROM public.tag_total
WHERE submission_prefix = @submissionPrefix
    AND claim_type = @claimType
    AND mac_workload_id = @workload
ORDER BY claims DESC, tag`

// List returns one claim type and workload's tag totals (measure.All for
// the total over either), most claims first, ties by tag.
func (t TableAccess) List(ctx context.Context, submissionPrefix, claimType, workload string) ([]TagTotal, error) {
	rows, err := t.db.Query(ctx, listSQL, pgx.NamedArgs{
		"submissionPrefix": submissionPrefix,
		"claimType":        claimType,
		"workload":         workload,
	})
	if err != nil {
		return nil, fmt.Errorf("tagtotal: listing: %w", err)
	}
	defer rows.Close()
	var out []TagTotal
	for rows.Next() {
		var tt TagTotal
		err := rows.Scan(
			&tt.SubmissionPrefix,
			&tt.ClaimType,
			&tt.MACWorkloadID,
			&tt.Tag,
			&tt.Claims,
			&tt.Dollars,
		)
		if err != nil {
			return nil, fmt.Errorf("tagtotal: scanning: %w", err)
		}
		out = append(out, tt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tagtotal: listing: %w", err)
	}
	return out, nil
}
