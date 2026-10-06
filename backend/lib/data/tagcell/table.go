package tagcell

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for the tag cube.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

const deleteBySubmissionSQL = `DELETE FROM public.tag_cell WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("tagcell: deleting submission: %w", err)
	}
	return nil
}

// totalsSQL adds the tag cells the filter keeps by tag, most claims first.
// Each list argument is sent empty rather than nil, since cardinality(NULL)
// is NULL and would keep no row.
const totalsSQL = `
SELECT
    tag
    , sum(claims)::bigint
    , trim_scale(sum(dollars))::text
FROM public.tag_cell
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
GROUP BY tag
ORDER BY sum(claims) DESC, tag`

// Totals adds the tag cells the filter keeps, one row per tag, most claims
// first.
func (t TableAccess) Totals(ctx context.Context, f Filter) ([]TagSums, error) {
	rows, err := t.db.Query(ctx, totalsSQL, pgx.NamedArgs{
		"submissionPrefix": f.SubmissionPrefix,
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
	})
	if err != nil {
		return nil, fmt.Errorf("tagcell: adding by tag: %w", err)
	}
	defer rows.Close()
	var out []TagSums
	for rows.Next() {
		var s TagSums
		if err := rows.Scan(&s.Tag, &s.Claims, &s.Dollars); err != nil {
			return nil, fmt.Errorf("tagcell: scanning tag: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tagcell: adding by tag: %w", err)
	}
	return out, nil
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
