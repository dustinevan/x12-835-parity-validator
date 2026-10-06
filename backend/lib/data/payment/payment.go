// Package payment reads and writes the payment table: one row per 835
// transaction (ST..SE) per submission, the payment header the transaction's
// claims sit under, so results can be cut by the MAC's payment context.
package payment

import "time"

// Payment is one row: the header of one 835 transaction. Amounts are
// decimal strings exactly as received. An empty string, or a zero time, is
// a column the file did not send, written as NULL.
type Payment struct {
	SubmissionPrefix string
	FileName         string    // relative to the submission prefix
	PaymentNumber    int       // 1st, 2nd... ST..SE in the file
	PayerIdentifier  string    // 1000A N1*PR N104
	PayerName        string    // 1000A N1*PR N102
	PayeeIdentifier  string    // 1000B N1*PE N104
	PayeeName        string    // 1000B N1*PE N102
	PayeeTaxID       string    // REF*TJ
	PayeeState       string    // 1000B N4 N402
	PaymentDate      time.Time // BPR16
	PaymentMethod    string    // BPR04
	HandlingCode     string    // BPR01
	TraceNumber      string    // TRN02
	ProductionDate   time.Time // DTM*405
	TotalPaid        string    // BPR02; for the transaction's own balance check, never compared between sides
	ClaimCount       int       // CLP segments in the transaction
	PLBCount         int       // PLB segments in the transaction
	PLBTotal         string    // the sum of the PLB adjustment amounts sent; "0" with none
}
