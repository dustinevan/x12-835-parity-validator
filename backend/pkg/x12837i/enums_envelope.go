package x12837i

import "fmt"

// TransactionTypeGS01_ST01 is the transaction-type code, carried on both
// GS01 (Functional Identifier Code) and ST01 (Transaction Set Identifier
// Code). Transmittal 907 gives both positions the same combined name
// (frontend/src/data/x12Names837i.ts TRANSACTION_TYPE_NAME), the same
// pattern as x12835's FunctionalIDGS01/TransactionSetIDST01 pair, modeled
// here as one shared type since the two codes mean the same thing.
type TransactionTypeGS01_ST01 string

const (
	TransactionTypeGS01_ST01HC  TransactionTypeGS01_ST01 = "HC"
	TransactionTypeGS01_ST01837 TransactionTypeGS01_ST01 = "837"
)

// The codes the guide allows for GS01/ST01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var transactionTypeGS01_ST01Codes = []codeEntry{
	{"HC", "Health Care Claim", "Health Care Claim — Transmittal 907, PDF p.7 (GS01), PDF p.8 (ST01)"},
	{"837", "Health Care Claim", "Health Care Claim — Transmittal 907, PDF p.7 (GS01), PDF p.8 (ST01)"},
}

// String returns the code as the file sent it.
func (v TransactionTypeGS01_ST01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v TransactionTypeGS01_ST01) Known() bool {
	return knownCode(transactionTypeGS01_ST01Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v TransactionTypeGS01_ST01) ShortDescription() string {
	return shortOf(transactionTypeGS01_ST01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v TransactionTypeGS01_ST01) Description() string {
	return descriptionOf(transactionTypeGS01_ST01Codes, string(v))
}

// ParseTransactionTypeGS01_ST01 reads one GS01/ST01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseTransactionTypeGS01_ST01(code string) (TransactionTypeGS01_ST01, error) {
	v := TransactionTypeGS01_ST01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing GS01/ST01: a value not on the list")
}

// ClaimOrEncounterBHT06 is BHT06, the Claim or Encounter Identifier.
//
// Transmittal 907 lists the value set "31, CH, RP" without naming any of
// the three in text; CH (Chargeable) and RP (Reporting) are X12's
// documented claim/encounter codes, but neither the flat file nor this
// package invents English names beyond what the source gives, so all
// three are code-named.
type ClaimOrEncounterBHT06 string

const (
	ClaimOrEncounterBHT0631 ClaimOrEncounterBHT06 = "31"
	ClaimOrEncounterBHT06CH ClaimOrEncounterBHT06 = "CH"
	ClaimOrEncounterBHT06RP ClaimOrEncounterBHT06 = "RP"
)

// The codes the guide allows for BHT06, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var claimOrEncounterBHT06Codes = []codeEntry{
	{"31", "", "Transmittal 907, PDF p.9; not named in text"},
	{"CH", "", "Transmittal 907, PDF p.9; not named in text"},
	{"RP", "", "Transmittal 907, PDF p.9; not named in text"},
}

// String returns the code as the file sent it.
func (v ClaimOrEncounterBHT06) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v ClaimOrEncounterBHT06) Known() bool { return knownCode(claimOrEncounterBHT06Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v ClaimOrEncounterBHT06) ShortDescription() string {
	return shortOf(claimOrEncounterBHT06Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v ClaimOrEncounterBHT06) Description() string {
	return descriptionOf(claimOrEncounterBHT06Codes, string(v))
}

// ParseClaimOrEncounterBHT06 reads one BHT06 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseClaimOrEncounterBHT06(code string) (ClaimOrEncounterBHT06, error) {
	v := ClaimOrEncounterBHT06(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing BHT06: a value not on the list")
}

// HierarchicalLevelHL03 is HL03, the Hierarchical Level Code, shared by
// every HL segment in the transaction (2000A billing provider, 2000B
// subscriber, 2000C patient).
type HierarchicalLevelHL03 string

const (
	HierarchicalLevelHL0320 HierarchicalLevelHL03 = "20"
	HierarchicalLevelHL0322 HierarchicalLevelHL03 = "22"
	HierarchicalLevelHL0323 HierarchicalLevelHL03 = "23"
)

// The codes the guide allows for HL03, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var hierarchicalLevelHL03Codes = []codeEntry{
	{"20", "Billing Provider", "Billing Provider — Transmittal 907, PDF p.10 (2000A); PDF p.16 (2000B); PDF p.20 (2000C, not used for Medicare, where the patient is always the subscriber)"},
	{"22", "Subscriber", "Subscriber — Transmittal 907, PDF p.10 (2000A); PDF p.16 (2000B); PDF p.20 (2000C, not used for Medicare, where the patient is always the subscriber)"},
	{"23", "Patient", "Patient — Transmittal 907, PDF p.10 (2000A); PDF p.16 (2000B); PDF p.20 (2000C, not used for Medicare, where the patient is always the subscriber)"},
}

// String returns the code as the file sent it.
func (v HierarchicalLevelHL03) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v HierarchicalLevelHL03) Known() bool { return knownCode(hierarchicalLevelHL03Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v HierarchicalLevelHL03) ShortDescription() string {
	return shortOf(hierarchicalLevelHL03Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v HierarchicalLevelHL03) Description() string {
	return descriptionOf(hierarchicalLevelHL03Codes, string(v))
}

// ParseHierarchicalLevelHL03 reads one HL03 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseHierarchicalLevelHL03(code string) (HierarchicalLevelHL03, error) {
	v := HierarchicalLevelHL03(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HL03: a value not on the list")
}
