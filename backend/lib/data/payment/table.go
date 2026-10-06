package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys/db"
)

// TableAccess provides database operations for payments. Nothing reads the
// table yet, so it has no list.
type TableAccess struct {
	db db.Console
}

func NewTableAccess(db db.Console) TableAccess {
	return TableAccess{db: db}
}

// CreateMany bulk-loads payments with COPY.
func (t TableAccess) CreateMany(ctx context.Context, payments []Payment) (int64, error) {
	columns := []string{
		"submission_prefix",
		"file_name",
		"payment_number",
		"payer_identifier",
		"payer_name",
		"payee_identifier",
		"payee_name",
		"payee_tax_id",
		"payee_state",
		"payment_date",
		"payment_method",
		"handling_code",
		"trace_number",
		"production_date",
		"total_paid",
		"claim_count",
		"plb_count",
		"plb_total",
	}
	rows := pgx.CopyFromSlice(len(payments), func(i int) ([]any, error) {
		p := payments[i]
		totalPaid, err := numeric(p.TotalPaid)
		if err != nil {
			return nil, fmt.Errorf("BPR02: %w", err)
		}
		plbTotal, err := numeric(p.PLBTotal)
		if err != nil {
			return nil, fmt.Errorf("PLB total: %w", err)
		}
		return []any{
			p.SubmissionPrefix,
			p.FileName,
			p.PaymentNumber,
			p.PayerIdentifier,
			p.PayerName,
			p.PayeeIdentifier,
			p.PayeeName,
			text(p.PayeeTaxID),
			text(p.PayeeState),
			date(p.PaymentDate),
			text(p.PaymentMethod),
			text(p.HandlingCode),
			text(p.TraceNumber),
			date(p.ProductionDate),
			totalPaid,
			p.ClaimCount,
			p.PLBCount,
			plbTotal,
		}, nil
	})
	n, err := t.db.CopyFrom(ctx, pgx.Identifier{"public", "payment"}, columns, rows)
	if err != nil {
		return 0, fmt.Errorf("payment: copying: %w", err)
	}
	return n, nil
}

const deleteBySubmissionSQL = `DELETE FROM public.payment WHERE submission_prefix = @submissionPrefix`

func (t TableAccess) DeleteBySubmission(ctx context.Context, submissionPrefix string) error {
	if _, err := t.db.Exec(ctx, deleteBySubmissionSQL, pgx.NamedArgs{"submissionPrefix": submissionPrefix}); err != nil {
		return fmt.Errorf("payment: deleting submission: %w", err)
	}
	return nil
}

// text writes an empty string as NULL.
func text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// date writes a zero time as NULL.
func date(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: !t.IsZero()}
}

// numeric parses an exact decimal string for COPY, which sends NUMERIC in
// binary and so cannot take a Go string directly. An empty string is NULL.
func numeric(s string) (pgtype.Numeric, error) {
	if s == "" {
		return pgtype.Numeric{}, nil
	}
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("parsing decimal %q: %w", s, err)
	}
	return n, nil
}
