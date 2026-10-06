package claim

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for claims.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

// CreateMany bulk-loads claims with COPY.
func (t TableAccess) CreateMany(ctx context.Context, claims []Claim) (int64, error) {
	columns := []string{
		"submission_prefix",
		"match_key",
		"file_name",
		"payment_number",
		"claim_number",
		"medicare_beneficiary_identifier",
		"date_of_service",
		"total_claim_charge_amount",
		"federal_taxpayer_identification_number",
		"payee_identification_code",
		"mac_workload_id",
		"patient_control_number",
		"claim_type",
		"claim_status_code",
		"claim_payment_amount",
		"payer_claim_control_number",
		"claim_frequency_code",
		"adjudication",
		"jurisdiction",
		"part",
		"delimiters",
		"segments",
	}
	rows := pgx.CopyFromSlice(len(claims), func(i int) ([]any, error) {
		c := claims[i]
		charge, err := numeric(c.TotalClaimChargeAmount)
		if err != nil {
			return nil, fmt.Errorf("CLP03: %w", err)
		}
		paid, err := numeric(c.ClaimPaymentAmount)
		if err != nil {
			return nil, fmt.Errorf("CLP04: %w", err)
		}
		return []any{
			c.SubmissionPrefix,
			c.MatchKey,
			c.FileName,
			c.PaymentNumber,
			c.ClaimNumber,
			c.MedicareBeneficiaryIdentifier,
			pgtype.Date{Time: c.DateOfService, Valid: true},
			charge,
			c.FederalTaxpayerIdentificationNumber,
			c.PayeeIdentificationCode,
			c.MACWorkloadID,
			c.PatientControlNumber,
			c.ClaimType,
			c.ClaimStatusCode,
			paid,
			c.PayerClaimControlNumber,
			c.ClaimFrequencyCode,
			c.Adjudication,
			c.Jurisdiction,
			c.Part,
			c.Delimiters,
			c.Segments,
		}, nil
	})
	n, err := t.db.CopyFrom(ctx, pgx.Identifier{"public", "claim"}, columns, rows)
	if err != nil {
		return 0, fmt.Errorf("claim: copying: %w", err)
	}
	return n, nil
}

const deleteBySubmissionSQL = `DELETE FROM public.claim WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("claim: deleting submission: %w", err)
	}
	return nil
}

const listBySubmissionSQL = `
SELECT
    submission_prefix
    , match_key
    , file_name
    , payment_number
    , claim_number
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount::text
    , federal_taxpayer_identification_number
    , payee_identification_code
    , mac_workload_id
    , patient_control_number
    , claim_type
    , claim_status_code
    , claim_payment_amount::text
    , payer_claim_control_number
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , delimiters
    , segments
FROM public.claim
WHERE submission_prefix = @submissionPrefix
ORDER BY file_name, payment_number, claim_number`

// ListBySubmission returns a submission's claims in file order.
func (t TableAccess) ListBySubmission(ctx context.Context, submissionPrefix string) ([]Claim, error) {
	rows, err := t.db.Query(ctx, listBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix})
	if err != nil {
		return nil, fmt.Errorf("claim: listing submission: %w", err)
	}
	defer rows.Close()
	var out []Claim
	for rows.Next() {
		var c Claim
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
			&c.ClaimFrequencyCode,
			&c.Adjudication,
			&c.Jurisdiction,
			&c.Part,
			&c.Delimiters,
			&c.Segments,
		)
		if err != nil {
			return nil, fmt.Errorf("claim: scanning: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("claim: listing submission: %w", err)
	}
	return out, nil
}

const listByMatchKeySQL = `
SELECT
    submission_prefix
    , match_key
    , file_name
    , payment_number
    , claim_number
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount::text
    , federal_taxpayer_identification_number
    , payee_identification_code
    , mac_workload_id
    , patient_control_number
    , claim_type
    , claim_status_code
    , claim_payment_amount::text
    , payer_claim_control_number
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , delimiters
    , segments
FROM public.claim
WHERE submission_prefix = @submissionPrefix
    AND match_key = @matchKey
ORDER BY file_name, payment_number, claim_number`

// ListByMatchKey returns the claims in one submission that share a match key:
// usually one; several in the MAC submission when the key is set aside, or
// in a vendor submission when it is duplicated; none when missing.
func (t TableAccess) ListByMatchKey(ctx context.Context, submissionPrefix string, matchKey []byte) ([]Claim, error) {
	rows, err := t.db.Query(ctx, listByMatchKeySQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix, "matchKey": matchKey})
	if err != nil {
		return nil, fmt.Errorf("claim: listing by match key: %w", err)
	}
	defer rows.Close()
	var out []Claim
	for rows.Next() {
		var c Claim
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
			&c.ClaimFrequencyCode,
			&c.Adjudication,
			&c.Jurisdiction,
			&c.Part,
			&c.Delimiters,
			&c.Segments,
		)
		if err != nil {
			return nil, fmt.Errorf("claim: scanning: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("claim: listing by match key: %w", err)
	}
	return out, nil
}

// numeric parses an exact decimal string for COPY, which sends NUMERIC in
// binary and so cannot take a Go string directly.
func numeric(s string) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("parsing decimal %q: %w", s, err)
	}
	return n, nil
}
