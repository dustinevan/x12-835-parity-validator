package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// UnmatchedQuery reads the claims behind one outcome's comparison rows.
// claimPrefix is the MAC submission for missing, duplicated, and set aside,
// the vendor submission for extra. A set-aside key has several MAC claims.
type UnmatchedQuery struct {
	db db.Console
}

func NewUnmatchedQuery(db db.Console) UnmatchedQuery {
	return UnmatchedQuery{db: db}
}

// Count is a number of claims for one claim type or one workload.
type Count struct {
	Value string
	Count int64
}

const queryListUnmatched = `
SELECT
    cl.submission_prefix
    , cl.match_key
    , cl.file_name
    , cl.payment_number
    , cl.claim_number
    , cl.medicare_beneficiary_identifier
    , cl.date_of_service
    , cl.total_claim_charge_amount::text
    , cl.federal_taxpayer_identification_number
    , cl.payee_identification_code
    , cl.mac_workload_id
    , cl.patient_control_number
    , cl.claim_type
    , cl.claim_status_code
    , cl.claim_payment_amount::text
    , cl.payer_claim_control_number
    , cl.delimiters
    , ''::bytea
FROM public.comparison c
JOIN public.claim cl
    ON cl.submission_prefix = @claimPrefix
    AND cl.match_key = c.match_key
WHERE c.submission_prefix = @submissionPrefix
    AND c.outcome = @outcome
ORDER BY cl.match_key, cl.file_name, cl.payment_number, cl.claim_number
LIMIT @limit`

// List returns up to limit claims behind a submission's comparison rows of
// one outcome, in match key order, so claims sharing a set-aside key are
// adjacent. Segments are left empty.
func (q UnmatchedQuery) List(ctx context.Context, submissionPrefix, claimPrefix, outcome string, limit int) ([]claim.Claim, error) {
	rows, err := q.db.Query(ctx, queryListUnmatched, pgx.NamedArgs{
		"submissionPrefix": submissionPrefix,
		"claimPrefix":      claimPrefix,
		"outcome":          outcome,
		"limit":            limit,
	})
	if err != nil {
		return nil, fmt.Errorf("query: listing unmatched claims: %w", err)
	}
	defer rows.Close()
	var out []claim.Claim
	for rows.Next() {
		var c claim.Claim
		err := rows.Scan(
			&c.SubmissionPrefix,
			&c.MatchKey,
			&c.FileName,
			&c.PaymentNumber,
			&c.ClaimNumber,
			&c.MedicareBeneficiaryIdentifier,
			&c.DateOfService,
			&c.TotalClaimChargeAmount,
			&c.FederalTaxpayerIdentificationNumber,
			&c.PayeeIdentificationCode,
			&c.MACWorkloadID,
			&c.PatientControlNumber,
			&c.ClaimType,
			&c.ClaimStatusCode,
			&c.ClaimPaymentAmount,
			&c.PayerClaimControlNumber,
			&c.Delimiters,
			&c.Segments,
		)
		if err != nil {
			return nil, fmt.Errorf("query: scanning unmatched claim: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query: listing unmatched claims: %w", err)
	}
	return out, nil
}

const queryCountUnmatchedByClaimType = `
SELECT cl.claim_type, count(*)
FROM public.comparison c
JOIN public.claim cl
    ON cl.submission_prefix = @claimPrefix
    AND cl.match_key = c.match_key
WHERE c.submission_prefix = @submissionPrefix
    AND c.outcome = @outcome
GROUP BY cl.claim_type
ORDER BY count(*) DESC, cl.claim_type`

// CountByClaimType counts List's claims per claim type.
func (q UnmatchedQuery) CountByClaimType(ctx context.Context, submissionPrefix, claimPrefix, outcome string) ([]Count, error) {
	rows, err := q.db.Query(ctx, queryCountUnmatchedByClaimType, pgx.NamedArgs{
		"submissionPrefix": submissionPrefix,
		"claimPrefix":      claimPrefix,
		"outcome":          outcome,
	})
	if err != nil {
		return nil, fmt.Errorf("query: counting unmatched claims by claim type: %w", err)
	}
	defer rows.Close()
	var out []Count
	for rows.Next() {
		var c Count
		if err := rows.Scan(&c.Value, &c.Count); err != nil {
			return nil, fmt.Errorf("query: scanning claim type count: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query: counting unmatched claims by claim type: %w", err)
	}
	return out, nil
}

const queryCountUnmatchedByWorkload = `
SELECT cl.mac_workload_id, count(*)
FROM public.comparison c
JOIN public.claim cl
    ON cl.submission_prefix = @claimPrefix
    AND cl.match_key = c.match_key
WHERE c.submission_prefix = @submissionPrefix
    AND c.outcome = @outcome
GROUP BY cl.mac_workload_id
ORDER BY count(*) DESC, cl.mac_workload_id`

// CountByWorkload counts List's claims per MAC workload.
func (q UnmatchedQuery) CountByWorkload(ctx context.Context, submissionPrefix, claimPrefix, outcome string) ([]Count, error) {
	rows, err := q.db.Query(ctx, queryCountUnmatchedByWorkload, pgx.NamedArgs{
		"submissionPrefix": submissionPrefix,
		"claimPrefix":      claimPrefix,
		"outcome":          outcome,
	})
	if err != nil {
		return nil, fmt.Errorf("query: counting unmatched claims by workload: %w", err)
	}
	defer rows.Close()
	var out []Count
	for rows.Next() {
		var c Count
		if err := rows.Scan(&c.Value, &c.Count); err != nil {
			return nil, fmt.Errorf("query: scanning workload count: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query: counting unmatched claims by workload: %w", err)
	}
	return out, nil
}
