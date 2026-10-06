// Package claim reads and writes the claim table: one row per claim (CLP
// with its service lines) per submission, holding the claim's exact bytes.
package claim

import "time"

// Claim is one row. Amounts are decimal strings exactly as received.
type Claim struct {
	SubmissionPrefix                    string
	MatchKey                            []byte // SHA-256, 32 bytes
	FileName                            string // relative to the submission prefix
	PaymentNumber                       int    // 1st, 2nd... ST..SE in the file
	ClaimNumber                         int    // 1st, 2nd... claim in the payment
	MedicareBeneficiaryIdentifier       string
	DateOfService                       time.Time
	TotalClaimChargeAmount              string
	FederalTaxpayerIdentificationNumber string
	PayeeIdentificationCode             string
	MACWorkloadID                       string
	PatientControlNumber                string
	ClaimType                           string
	ClaimStatusCode                     string
	ClaimPaymentAmount                  string
	PayerClaimControlNumber             string
	ClaimFrequencyCode                  string // CLP09 as sent; "" when not sent
	Adjudication                        string // "inpatient" (MIA), "outpatient" (MOA), or ""
	Jurisdiction                        string // the workload's MAC jurisdiction (backend/pkg/macs); "" when unknown
	Part                                string // "A", "B", "DME", or ""
	Delimiters                          string // element, component, repetition, segment terminator
	Segments                            []byte
}
