package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// LoadQuery is the load stage's cross-table SQL: it writes claim_repeat
// from the submission's claim rows, and for the MAC submission mac_cell,
// the legacy cube. The caller (lib/core/control/submission) runs it in the
// load transaction after the claims are copied in, having deleted the
// submission's earlier rows first.
type LoadQuery struct {
	db db.Console
}

func NewLoadQuery(db db.Console) LoadQuery {
	return LoadQuery{db: db}
}

// queryCreateClaimRepeats writes one claim_repeat row per match key that
// occurs on more than one of the submission's claims, with the count.
const queryCreateClaimRepeats = `
INSERT INTO public.claim_repeat (
    submission_prefix
    , match_key
    , claims
)
SELECT
    submission_prefix
    , match_key
    , count(*)
FROM public.claim
WHERE submission_prefix = @submissionPrefix
GROUP BY submission_prefix, match_key
HAVING count(*) > 1`

// CreateClaimRepeats writes a submission's claim_repeat rows from its claim
// rows and returns how many keys the submission repeats.
func (q LoadQuery) CreateClaimRepeats(ctx context.Context, submissionPrefix string) (int64, error) {
	tag, err := q.db.Exec(ctx, queryCreateClaimRepeats, pgx.NamedArgs{"submissionPrefix": submissionPrefix})
	if err != nil {
		return 0, fmt.Errorf("query: creating repeated keys: %w", err)
	}
	return tag.RowsAffected(), nil
}

// queryCreateMACCells writes the legacy cube: one mac_cell row per distinct
// combination of the nine dimensions among the submission's claim rows,
// with the claims, the sum of CLP04, and the sum of CLP03. paid_band is the
// claim's CLP04 in the vendor cube's bands (queryCreateCells): zero, then
// under 100, 1000, and 10000, then 10000 and more. service_month is the
// first day of the month of the date of service.
const queryCreateMACCells = `
WITH banded AS (
    SELECT
        jurisdiction
        , part
        , mac_workload_id
        , claim_type
        , adjudication
        , claim_status_code
        , claim_frequency_code
        , CASE
            WHEN claim_payment_amount = 0 THEN 'zero'
            WHEN claim_payment_amount < 100 THEN 'under 100'
            WHEN claim_payment_amount < 1000 THEN 'under 1000'
            WHEN claim_payment_amount < 10000 THEN 'under 10000'
            ELSE '10000 and more'
          END AS paid_band
        , date_trunc('month', date_of_service)::date AS service_month
        , claim_payment_amount
        , total_claim_charge_amount
    FROM public.claim
    WHERE submission_prefix = @submissionPrefix
)
INSERT INTO public.mac_cell (
    submission_prefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , claim_status_code
    , claim_frequency_code
    , paid_band
    , service_month
    , claims
    , paid
    , billed
)
SELECT
    @submissionPrefix
    , jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , claim_status_code
    , claim_frequency_code
    , paid_band
    , service_month
    , count(*)
    , sum(claim_payment_amount)
    , sum(total_claim_charge_amount)
FROM banded
GROUP BY
    jurisdiction
    , part
    , mac_workload_id
    , claim_type
    , adjudication
    , claim_status_code
    , claim_frequency_code
    , paid_band
    , service_month`

// CreateMACCells writes the MAC submission's mac_cell rows from its claim
// rows.
func (q LoadQuery) CreateMACCells(ctx context.Context, submissionPrefix string) error {
	_, err := q.db.Exec(ctx, queryCreateMACCells, pgx.NamedArgs{"submissionPrefix": submissionPrefix})
	if err != nil {
		return fmt.Errorf("query: creating MAC cells: %w", err)
	}
	return nil
}
