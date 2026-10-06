package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// CompareQuery is the compare stage's SQL: the MAC workloads a vendor
// submission is compared in, the stream of claim pairs Go compares, and the
// four outcomes the database decides on its own (set_aside, duplicated,
// missing, extra), each one INSERT ... SELECT. The caller
// (lib/core/control/submission) runs them in that order, one workload at a
// time. The match key contains the workload, so every key, and so every
// outcome, falls in exactly one workload.
//
// Which keys repeat was decided at load: claim_repeat holds one row per
// (submission, key) with more than one claim. Every statement here reads
// that table instead of counting claims per key again. A key is compared
// when neither side has a row; set aside when the MAC side has one, whatever
// the vendor sent; duplicated when only the vendor side has one and the MAC
// has the key; missing when the MAC has it once and the vendor not at all;
// extra when only the vendor has it, once or more.
type CompareQuery struct {
	db db.Console
}

func NewCompareQuery(db db.Console) CompareQuery {
	return CompareQuery{db: db}
}

const queryListWorkloads = `
SELECT DISTINCT mac_workload_id
FROM public.claim
WHERE submission_prefix IN (@macPrefix, @vendorPrefix)
ORDER BY mac_workload_id`

// ListWorkloads returns the distinct MAC workloads of the MAC submission's
// claims and the vendor submission's claims together, in workload order.
func (q CompareQuery) ListWorkloads(ctx context.Context, macPrefix, vendorPrefix string) ([]string, error) {
	rows, err := q.db.Query(ctx, queryListWorkloads, pgx.NamedArgs{"macPrefix": macPrefix, "vendorPrefix": vendorPrefix})
	if err != nil {
		return nil, fmt.Errorf("query: listing workloads: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var workload string
		if err := rows.Scan(&workload); err != nil {
			return nil, fmt.Errorf("query: scanning workload: %w", err)
		}
		out = append(out, workload)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query: listing workloads: %w", err)
	}
	return out, nil
}

// queryPairs joins the MAC claims and the vendor claims of one workload on
// match key, leaving out the keys either side repeats (a claim_repeat row
// on that side), and returns both claims' columns. The two NOT EXISTS
// tests are anti-joins against the small claim_repeat table on its primary
// key, so the rows stream straight out of the join, with no grouping of
// the workload's claims before the first row. It has no ORDER BY: a sort
// would make the database gather the whole workload before returning the
// first row, and the tables are keyed, so row order carries no meaning.
const queryPairs = `
SELECT
    m.match_key
    , m.mac_workload_id
    , m.claim_type
    , m.medicare_beneficiary_identifier
    , m.date_of_service
    , m.total_claim_charge_amount::text
    , m.claim_status_code
    , m.claim_payment_amount::text
    , m.claim_frequency_code
    , m.adjudication
    , m.jurisdiction
    , m.part
    , m.delimiters
    , m.segments
    , v.claim_status_code
    , v.claim_payment_amount::text
    , v.delimiters
    , v.segments
FROM public.claim m
JOIN public.claim v
    ON v.submission_prefix = @vendorPrefix
    AND v.mac_workload_id = @workload
    AND v.match_key = m.match_key
WHERE m.submission_prefix = @macPrefix
    AND m.mac_workload_id = @workload
    AND NOT EXISTS (
        SELECT 1
        FROM public.claim_repeat r
        WHERE r.submission_prefix = @macPrefix
            AND r.match_key = m.match_key
    )
    AND NOT EXISTS (
        SELECT 1
        FROM public.claim_repeat r
        WHERE r.submission_prefix = @vendorPrefix
            AND r.match_key = m.match_key
    )`

// Pairs streams every key of one workload that is on exactly one claim in
// the MAC submission and exactly one in the vendor submission, with both
// claims. The caller reads one pair at a time with Next and Pair, so no
// workload is ever held in memory. The statement runs for the whole
// workload, while the caller compares each pair as it arrives, so the
// connection's statement timeout must be zero or longer than that (vdb's
// APP_SERVER_DB_STATEMENT_TIMEOUT defaults to 0s, no limit).
func (q CompareQuery) Pairs(ctx context.Context, macPrefix, vendorPrefix, workload string) (PairStream, error) {
	rows, err := q.db.Query(ctx, queryPairs, pgx.NamedArgs{"macPrefix": macPrefix, "vendorPrefix": vendorPrefix, "workload": workload})
	if err != nil {
		return PairStream{}, fmt.Errorf("query: streaming pairs: %w", err)
	}
	return PairStream{rows: rows}, nil
}

// PairStream reads the pairs of one workload one at a time as the database
// returns them. Use it as a bufio.Scanner: Next, then Pair, until Next is
// false; then Err. Close releases the connection; the caller defers it.
type PairStream struct {
	rows   pgx.Rows
	mac    claim.Claim
	vendor claim.Claim
	err    error
}

// Next reads the next pair. It returns false at the end of the stream or
// after an error, which Err reports.
func (s *PairStream) Next() bool {
	if s.err != nil || !s.rows.Next() {
		return false
	}
	var mac, vendor claim.Claim
	err := s.rows.Scan(
		&mac.MatchKey,
		&mac.MACWorkloadID,
		&mac.ClaimType,
		&mac.MedicareBeneficiaryIdentifier,
		&mac.DateOfService,
		&mac.TotalClaimChargeAmount,
		&mac.ClaimStatusCode,
		&mac.ClaimPaymentAmount,
		&mac.ClaimFrequencyCode,
		&mac.Adjudication,
		&mac.Jurisdiction,
		&mac.Part,
		&mac.Delimiters,
		&mac.Segments,
		&vendor.ClaimStatusCode,
		&vendor.ClaimPaymentAmount,
		&vendor.Delimiters,
		&vendor.Segments,
	)
	if err != nil {
		s.err = fmt.Errorf("query: scanning pair: %w", err)
		s.rows.Close()
		return false
	}
	s.mac, s.vendor = mac, vendor
	return true
}

// Pair returns the pair Next read: the MAC claim with its match key,
// workload, claim type, beneficiary, date of service, charge, status,
// payment, frequency code, adjudication, jurisdiction, Part, delimiters,
// and segments filled, and the vendor claim with its
// status, payment, delimiters, and segments. Every other field is zero.
func (s *PairStream) Pair() (claim.Claim, claim.Claim) {
	return s.mac, s.vendor
}

// Err returns the first error the stream met, if any. Call it once Next has
// returned false.
func (s *PairStream) Err() error {
	if s.err != nil {
		return s.err
	}
	if err := s.rows.Err(); err != nil {
		return fmt.Errorf("query: streaming pairs: %w", err)
	}
	return nil
}

// Close releases the stream's connection. Calling it more than once is safe.
func (s *PairStream) Close() {
	s.rows.Close()
}

// queryCreateSetAside writes one set_aside row per key of the workload that
// the MAC submission repeats (a MAC claim_repeat row), whatever the vendor
// sent for it. The row's columns come from the first MAC claim of the key
// in file order; the vendor columns are NULL.
const queryCreateSetAside = `
WITH picked AS (
    SELECT DISTINCT ON (m.match_key)
        m.match_key
        , m.mac_workload_id
        , m.claim_type
        , m.medicare_beneficiary_identifier
        , m.date_of_service
        , m.total_claim_charge_amount
        , m.claim_status_code
        , m.claim_payment_amount
        , m.claim_frequency_code
        , m.adjudication
        , m.jurisdiction
        , m.part
    FROM public.claim m
    JOIN public.claim_repeat r
        ON r.submission_prefix = @macPrefix
        AND r.match_key = m.match_key
    WHERE m.submission_prefix = @macPrefix
        AND m.mac_workload_id = @workload
    ORDER BY m.match_key, m.file_name, m.payment_number, m.claim_number
)
INSERT INTO public.comparison (
    submission_prefix
    , match_key
    , outcome
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , mac_claim_status_code
    , vendor_claim_status_code
    , mac_claim_payment_amount
    , vendor_claim_payment_amount
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , difference_count
    , tags
)
SELECT
    @vendorPrefix::text
    , match_key
    , 'set_aside'
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , claim_status_code
    , NULL
    , claim_payment_amount
    , NULL
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , 0
    , '{}'
FROM picked`

// CreateSetAside writes a vendor submission's set_aside rows for one MAC
// workload and returns how many it wrote.
func (q CompareQuery) CreateSetAside(ctx context.Context, vendorPrefix, macPrefix, workload string) (int64, error) {
	tag, err := q.db.Exec(ctx, queryCreateSetAside, pgx.NamedArgs{"vendorPrefix": vendorPrefix, "macPrefix": macPrefix, "workload": workload})
	if err != nil {
		return 0, fmt.Errorf("query: creating set-aside rows: %w", err)
	}
	return tag.RowsAffected(), nil
}

// queryCreateDuplicated writes one duplicated row per key of the workload
// that the MAC has on exactly one claim (no MAC claim_repeat row) and the
// vendor sent more than once (a vendor claim_repeat row, which load writes
// only for a key on two or more of the vendor's claims). The row's columns
// come from that MAC claim; the vendor columns are NULL, since no vendor
// copy was chosen.
const queryCreateDuplicated = `
INSERT INTO public.comparison (
    submission_prefix
    , match_key
    , outcome
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , mac_claim_status_code
    , vendor_claim_status_code
    , mac_claim_payment_amount
    , vendor_claim_payment_amount
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , difference_count
    , tags
)
SELECT
    @vendorPrefix::text
    , m.match_key
    , 'duplicated'
    , m.mac_workload_id
    , m.claim_type
    , m.medicare_beneficiary_identifier
    , m.date_of_service
    , m.total_claim_charge_amount
    , m.claim_status_code
    , NULL
    , m.claim_payment_amount
    , NULL
    , m.claim_frequency_code
    , m.adjudication
    , m.jurisdiction
    , m.part
    , 0
    , '{}'
FROM public.claim m
JOIN public.claim_repeat r
    ON r.submission_prefix = @vendorPrefix
    AND r.match_key = m.match_key
WHERE m.submission_prefix = @macPrefix
    AND m.mac_workload_id = @workload
    AND NOT EXISTS (
        SELECT 1
        FROM public.claim_repeat mr
        WHERE mr.submission_prefix = @macPrefix
            AND mr.match_key = m.match_key
    )`

// CreateDuplicated writes a vendor submission's duplicated rows for one MAC
// workload and returns how many it wrote.
func (q CompareQuery) CreateDuplicated(ctx context.Context, vendorPrefix, macPrefix, workload string) (int64, error) {
	tag, err := q.db.Exec(ctx, queryCreateDuplicated, pgx.NamedArgs{"vendorPrefix": vendorPrefix, "macPrefix": macPrefix, "workload": workload})
	if err != nil {
		return 0, fmt.Errorf("query: creating duplicated rows: %w", err)
	}
	return tag.RowsAffected(), nil
}

// queryCreateMissing writes one missing row per key of the workload that is
// on exactly one MAC claim (no MAC claim_repeat row) and on no vendor claim.
// The row's columns come from that MAC claim; the vendor columns are NULL.
const queryCreateMissing = `
INSERT INTO public.comparison (
    submission_prefix
    , match_key
    , outcome
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , mac_claim_status_code
    , vendor_claim_status_code
    , mac_claim_payment_amount
    , vendor_claim_payment_amount
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , difference_count
    , tags
)
SELECT
    @vendorPrefix::text
    , m.match_key
    , 'missing'
    , m.mac_workload_id
    , m.claim_type
    , m.medicare_beneficiary_identifier
    , m.date_of_service
    , m.total_claim_charge_amount
    , m.claim_status_code
    , NULL
    , m.claim_payment_amount
    , NULL
    , m.claim_frequency_code
    , m.adjudication
    , m.jurisdiction
    , m.part
    , 0
    , '{}'
FROM public.claim m
WHERE m.submission_prefix = @macPrefix
    AND m.mac_workload_id = @workload
    AND NOT EXISTS (
        SELECT 1
        FROM public.claim_repeat r
        WHERE r.submission_prefix = @macPrefix
            AND r.match_key = m.match_key
    )
    AND NOT EXISTS (
        SELECT 1
        FROM public.claim v
        WHERE v.submission_prefix = @vendorPrefix
            AND v.mac_workload_id = @workload
            AND v.match_key = m.match_key
    )`

// CreateMissing writes a vendor submission's missing rows for one MAC
// workload and returns how many it wrote.
func (q CompareQuery) CreateMissing(ctx context.Context, vendorPrefix, macPrefix, workload string) (int64, error) {
	tag, err := q.db.Exec(ctx, queryCreateMissing, pgx.NamedArgs{"vendorPrefix": vendorPrefix, "macPrefix": macPrefix, "workload": workload})
	if err != nil {
		return 0, fmt.Errorf("query: creating missing rows: %w", err)
	}
	return tag.RowsAffected(), nil
}

// queryCreateExtra writes one extra row per key of the workload that is on
// at least one vendor claim and on no MAC claim; a key the vendor repeated
// is still extra, not duplicated. The row's columns come from the first
// vendor claim of the key in file order; the MAC status and payment are
// NULL.
const queryCreateExtra = `
WITH picked AS (
    SELECT DISTINCT ON (v.match_key)
        v.match_key
        , v.mac_workload_id
        , v.claim_type
        , v.medicare_beneficiary_identifier
        , v.date_of_service
        , v.total_claim_charge_amount
        , v.claim_status_code
        , v.claim_payment_amount
        , v.claim_frequency_code
        , v.adjudication
        , v.jurisdiction
        , v.part
    FROM public.claim v
    WHERE v.submission_prefix = @vendorPrefix
        AND v.mac_workload_id = @workload
        AND NOT EXISTS (
            SELECT 1
            FROM public.claim m
            WHERE m.submission_prefix = @macPrefix
                AND m.mac_workload_id = @workload
                AND m.match_key = v.match_key
        )
    ORDER BY v.match_key, v.file_name, v.payment_number, v.claim_number
)
INSERT INTO public.comparison (
    submission_prefix
    , match_key
    , outcome
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , mac_claim_status_code
    , vendor_claim_status_code
    , mac_claim_payment_amount
    , vendor_claim_payment_amount
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , difference_count
    , tags
)
SELECT
    @vendorPrefix::text
    , match_key
    , 'extra'
    , mac_workload_id
    , claim_type
    , medicare_beneficiary_identifier
    , date_of_service
    , total_claim_charge_amount
    , NULL
    , claim_status_code
    , NULL
    , claim_payment_amount
    , claim_frequency_code
    , adjudication
    , jurisdiction
    , part
    , 0
    , '{}'
FROM picked`

// CreateExtra writes a vendor submission's extra rows for one MAC workload
// and returns how many it wrote.
func (q CompareQuery) CreateExtra(ctx context.Context, vendorPrefix, macPrefix, workload string) (int64, error) {
	tag, err := q.db.Exec(ctx, queryCreateExtra, pgx.NamedArgs{"vendorPrefix": vendorPrefix, "macPrefix": macPrefix, "workload": workload})
	if err != nil {
		return 0, fmt.Errorf("query: creating extra rows: %w", err)
	}
	return tag.RowsAffected(), nil
}
