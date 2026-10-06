package query

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/comparison"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// ClaimListQuery reads the claim list for the views: the comparison rows a
// comparison.Filter selects, joined to the claims they stand for. The
// list's totals and histogram read the comparison table alone
// (comparison.TableAccess Totals and Histogram).
//
// In every statement the tag lists may be NULL (a nil slice) or empty;
// either selects everything for that list.
type ClaimListQuery struct {
	db db.Console
}

func NewClaimListQuery(db db.Console) ClaimListQuery {
	return ClaimListQuery{db: db}
}

// Sort orders for List.
const (
	SortDelta           = "delta"  // paid delta, smallest first
	SortDeltaDescending = "-delta" // paid delta, largest first
	SortPaid            = "paid"   // MAC paid, largest first
	SortCharge          = "charge" // billed charge, largest first
)

// Row is a comparison with the key fields and CLP01 of the claim it stands
// for: the MAC claim, or the vendor claim when the outcome is extra.
type Row struct {
	MatchKey                            []byte
	Outcome                             string
	ClaimType                           string
	MedicareBeneficiaryIdentifier       string
	DateOfService                       time.Time
	TotalClaimChargeAmount              string
	FederalTaxpayerIdentificationNumber string
	PayeeIdentificationCode             string
	MACWorkloadID                       string
	PatientControlNumber                string
	MACClaimPaymentAmount               string // empty is NULL
	VendorClaimPaymentAmount            string // empty is NULL
	Delta                               string // vendor minus MAC paid; empty when either is NULL
	Tags                                []string
}

// queryListClaims reads each selected comparison with the key fields of
// the claim it stands for: the first MAC claim with the key in file order,
// or the vendor's claim when the key is extra. The tags filters use the
// GIN index on comparison.tags.
const queryListClaims = `
SELECT
    c.match_key
    , c.outcome
    , c.claim_type
    , c.medicare_beneficiary_identifier
    , c.date_of_service
    , c.total_claim_charge_amount::text
    , k.federal_taxpayer_identification_number
    , k.payee_identification_code
    , c.mac_workload_id
    , k.patient_control_number
    , COALESCE(c.mac_claim_payment_amount::text, '')
    , COALESCE(c.vendor_claim_payment_amount::text, '')
    , COALESCE((c.vendor_claim_payment_amount - c.mac_claim_payment_amount)::text, '')
    , c.tags
FROM public.comparison c
JOIN LATERAL (
    SELECT
        cl.federal_taxpayer_identification_number
        , cl.payee_identification_code
        , cl.patient_control_number
    FROM public.claim cl
    WHERE cl.submission_prefix = CASE WHEN c.outcome = 'extra' THEN c.submission_prefix ELSE @macPrefix::text END
        AND cl.match_key = c.match_key
    ORDER BY cl.file_name, cl.payment_number, cl.claim_number
    LIMIT 1
) k ON true
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
ORDER BY
    CASE WHEN @sort::text = 'delta' THEN c.vendor_claim_payment_amount - c.mac_claim_payment_amount END ASC NULLS LAST
    , CASE WHEN @sort::text = '-delta' THEN c.vendor_claim_payment_amount - c.mac_claim_payment_amount END DESC NULLS LAST
    , CASE WHEN @sort::text = 'paid' THEN c.mac_claim_payment_amount END DESC NULLS LAST
    , CASE WHEN @sort::text = 'charge' THEN c.total_claim_charge_amount END DESC
    , c.match_key
LIMIT @limit
OFFSET @offset`

// List returns one page of the rows a Filter selects. Ties break on
// match key, so pages are stable.
func (q ClaimListQuery) List(ctx context.Context, f comparison.Filter, macPrefix, sort string, limit, offset int) ([]Row, error) {
	rows, err := q.db.Query(ctx, queryListClaims, pgx.NamedArgs{
		"submissionPrefix": f.SubmissionPrefix,
		"claimType":        f.ClaimType,
		"workload":         f.Workload,
		"outcome":          f.Outcome,
		"all":              f.All,
		"any":              f.Any,
		"none":             f.None,
		"otherPrefix":      f.OtherSubmissionPrefix,
		"otherOutcome":     f.OtherOutcome,
		"macPrefix":        macPrefix,
		"sort":             sort,
		"limit":            limit,
		"offset":           offset,
	})
	if err != nil {
		return nil, fmt.Errorf("query: listing claims: %w", err)
	}
	defer rows.Close()
	var out []Row
	for rows.Next() {
		var r Row
		err := rows.Scan(
			&r.MatchKey,
			&r.Outcome,
			&r.ClaimType,
			&r.MedicareBeneficiaryIdentifier,
			&r.DateOfService,
			&r.TotalClaimChargeAmount,
			&r.FederalTaxpayerIdentificationNumber,
			&r.PayeeIdentificationCode,
			&r.MACWorkloadID,
			&r.PatientControlNumber,
			&r.MACClaimPaymentAmount,
			&r.VendorClaimPaymentAmount,
			&r.Delta,
			&r.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("query: scanning claim row: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query: listing claims: %w", err)
	}
	return out, nil
}
