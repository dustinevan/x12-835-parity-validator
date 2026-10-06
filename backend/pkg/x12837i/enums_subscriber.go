package x12837i

import "fmt"

// PayerResponsibilitySBR01 is SBR01, the Payer Responsibility Sequence
// Number Code, shared by the 2000B Subscriber Information SBR and the
// 2320 Other Subscriber Information SBR.
//
// Transmittal 907 lists the full value set A-H, P, S, T, U (PDF p.16) but
// only names P/S/T in text, as X12's own Primary/Secondary/Tertiary; A is
// additionally named here because X12's external code list itself names
// A "Payer Responsibility Four" (the list continues B..H as "...Five"
// through "...Eleven", a naming convention of the X12 standard code set,
// not of this project). frontend/src/data/x12Names837i.ts's SBR01_NAME
// table gives that one example (A) and stops there; B-H and U are
// code-named rather than guessing the remaining "...Five" through
// "...Eleven"/"Unknown" labels, since the source table does not state
// them.
type PayerResponsibilitySBR01 string

const (
	PayerResponsibilitySBR01A PayerResponsibilitySBR01 = "A"
	PayerResponsibilitySBR01B PayerResponsibilitySBR01 = "B"
	PayerResponsibilitySBR01C PayerResponsibilitySBR01 = "C"
	PayerResponsibilitySBR01D PayerResponsibilitySBR01 = "D"
	PayerResponsibilitySBR01E PayerResponsibilitySBR01 = "E"
	PayerResponsibilitySBR01F PayerResponsibilitySBR01 = "F"
	PayerResponsibilitySBR01G PayerResponsibilitySBR01 = "G"
	PayerResponsibilitySBR01H PayerResponsibilitySBR01 = "H"
	PayerResponsibilitySBR01P PayerResponsibilitySBR01 = "P"
	PayerResponsibilitySBR01S PayerResponsibilitySBR01 = "S"
	PayerResponsibilitySBR01T PayerResponsibilitySBR01 = "T"
	PayerResponsibilitySBR01U PayerResponsibilitySBR01 = "U"
)

// The codes the guide allows for SBR01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var payerResponsibilitySBR01Codes = []codeEntry{
	{"A", "Payer Responsibility Four", "Payer Responsibility Four — Transmittal 907, PDF p.16; X12 Payer Responsibility Sequence Number Code external list"},
	{"B", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"C", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"D", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"E", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"F", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"G", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"H", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
	{"P", "Primary", "Primary — common X12 Payer Responsibility Sequence Number code"},
	{"S", "Secondary", "Secondary — common X12 Payer Responsibility Sequence Number code"},
	{"T", "Tertiary", "Tertiary — common X12 Payer Responsibility Sequence Number code"},
	{"U", "", "Transmittal 907, PDF p.16 (value set A-H, P, S, T, U); not named in text or in frontend/src/data/x12Names837i.ts"},
}

// String returns the code as the file sent it.
func (v PayerResponsibilitySBR01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v PayerResponsibilitySBR01) Known() bool {
	return knownCode(payerResponsibilitySBR01Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v PayerResponsibilitySBR01) ShortDescription() string {
	return shortOf(payerResponsibilitySBR01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v PayerResponsibilitySBR01) Description() string {
	return descriptionOf(payerResponsibilitySBR01Codes, string(v))
}

// ParsePayerResponsibilitySBR01 reads one SBR01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParsePayerResponsibilitySBR01(code string) (PayerResponsibilitySBR01, error) {
	v := PayerResponsibilitySBR01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing SBR01: a value not on the list")
}

// RelationshipSBR02 is the Individual Relationship Code, shared by SBR02
// (2000B subscriber coverage and 2320 other-subscriber coverage) and PAT01
// (2000C patient loop) — frontend/src/data/x12Names837i.ts's
// RELATIONSHIP_NAME table, named once and reused.
type RelationshipSBR02 string

const (
	RelationshipSBR0201 RelationshipSBR02 = "01"
	RelationshipSBR0218 RelationshipSBR02 = "18"
	RelationshipSBR0219 RelationshipSBR02 = "19"
	RelationshipSBR0220 RelationshipSBR02 = "20"
	RelationshipSBR0221 RelationshipSBR02 = "21"
	RelationshipSBR0239 RelationshipSBR02 = "39"
	RelationshipSBR0240 RelationshipSBR02 = "40"
	RelationshipSBR0253 RelationshipSBR02 = "53"
	RelationshipSBR02G8 RelationshipSBR02 = "G8"
)

// The codes the guide allows for SBR02, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var relationshipSBR02Codes = []codeEntry{
	{"01", "Spouse", "Spouse — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"18", "Self", "Self — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"19", "Child", "Child — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"20", "Employee", "Employee — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"21", "Unknown", "Unknown — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"39", "Organ Donor", "Organ Donor — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"40", "Cadaver Donor", "Cadaver Donor — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"53", "Life Partner", "Life Partner — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
	{"G8", "Other Relationship", "Other Relationship — common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (SBR02, 2000B, listing only 18); PDF p.20 (PAT01, 2000C); PDF p.97 (SBR02, 2320 other subscriber)"},
}

// String returns the code as the file sent it.
func (v RelationshipSBR02) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v RelationshipSBR02) Known() bool { return knownCode(relationshipSBR02Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v RelationshipSBR02) ShortDescription() string {
	return shortOf(relationshipSBR02Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v RelationshipSBR02) Description() string {
	return descriptionOf(relationshipSBR02Codes, string(v))
}

// ParseRelationshipSBR02 reads one SBR02 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseRelationshipSBR02(code string) (RelationshipSBR02, error) {
	v := RelationshipSBR02(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing SBR02: a value not on the list")
}

// ClaimFilingIndicatorSBR09 is SBR09, the Claim Filing Indicator Code,
// shared by the 2000B Subscriber Information SBR and the 2320 Other
// Subscriber Information SBR (identical value set at both).
type ClaimFilingIndicatorSBR09 string

const (
	ClaimFilingIndicatorSBR0911 ClaimFilingIndicatorSBR09 = "11"
	ClaimFilingIndicatorSBR0912 ClaimFilingIndicatorSBR09 = "12"
	ClaimFilingIndicatorSBR0913 ClaimFilingIndicatorSBR09 = "13"
	ClaimFilingIndicatorSBR0914 ClaimFilingIndicatorSBR09 = "14"
	ClaimFilingIndicatorSBR0915 ClaimFilingIndicatorSBR09 = "15"
	ClaimFilingIndicatorSBR0916 ClaimFilingIndicatorSBR09 = "16"
	ClaimFilingIndicatorSBR0917 ClaimFilingIndicatorSBR09 = "17"
	ClaimFilingIndicatorSBR09AM ClaimFilingIndicatorSBR09 = "AM"
	ClaimFilingIndicatorSBR09BL ClaimFilingIndicatorSBR09 = "BL"
	ClaimFilingIndicatorSBR09CH ClaimFilingIndicatorSBR09 = "CH"
	ClaimFilingIndicatorSBR09CI ClaimFilingIndicatorSBR09 = "CI"
	ClaimFilingIndicatorSBR09DS ClaimFilingIndicatorSBR09 = "DS"
	ClaimFilingIndicatorSBR09FI ClaimFilingIndicatorSBR09 = "FI"
	ClaimFilingIndicatorSBR09HM ClaimFilingIndicatorSBR09 = "HM"
	ClaimFilingIndicatorSBR09LM ClaimFilingIndicatorSBR09 = "LM"
	ClaimFilingIndicatorSBR09MA ClaimFilingIndicatorSBR09 = "MA"
	ClaimFilingIndicatorSBR09MB ClaimFilingIndicatorSBR09 = "MB"
	ClaimFilingIndicatorSBR09MC ClaimFilingIndicatorSBR09 = "MC"
	ClaimFilingIndicatorSBR09OF ClaimFilingIndicatorSBR09 = "OF"
	ClaimFilingIndicatorSBR09TV ClaimFilingIndicatorSBR09 = "TV"
	ClaimFilingIndicatorSBR09VA ClaimFilingIndicatorSBR09 = "VA"
	ClaimFilingIndicatorSBR09WC ClaimFilingIndicatorSBR09 = "WC"
	ClaimFilingIndicatorSBR09ZZ ClaimFilingIndicatorSBR09 = "ZZ"
)

// The codes the guide allows for SBR09, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var claimFilingIndicatorSBR09Codes = []codeEntry{
	{"11", "Other Non-Federal Programs", "Other Non-Federal Programs — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"12", "Preferred Provider Organization (PPO)", "Preferred Provider Organization (PPO) — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"13", "Point of Service (POS)", "Point of Service (POS) — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"14", "Exclusive Provider Organization (EPO)", "Exclusive Provider Organization (EPO) — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"15", "Indemnity Insurance", "Indemnity Insurance — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"16", "Health Maintenance Organization (HMO) Medicare Risk", "Health Maintenance Organization (HMO) Medicare Risk — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"17", "Dental Maintenance Organization", "Dental Maintenance Organization — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"AM", "Automobile Medical", "Automobile Medical — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"BL", "Blue Cross/Blue Shield", "Blue Cross/Blue Shield — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"CH", "Champus", "Champus — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"CI", "Commercial Insurance Co.", "Commercial Insurance Co. — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"DS", "Disability", "Disability — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"FI", "Federal Employees Program", "Federal Employees Program — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"HM", "Health Maintenance Organization", "Health Maintenance Organization — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"LM", "Liability Medical", "Liability Medical — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"MA", "Medicare Part A", "Medicare Part A — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"MB", "Medicare Part B", "Medicare Part B — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"MC", "Medicaid", "Medicaid — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"OF", "Other Federal Program", "Other Federal Program — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"TV", "Title V", "Title V — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"VA", "Veterans Affairs Plan", "Veterans Affairs Plan — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"WC", "Workers' Compensation Health Claim", "Workers' Compensation Health Claim — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
	{"ZZ", "Mutually Defined", "Mutually Defined — common X12 Claim Filing Indicator Code; Transmittal 907, PDF p.16 (2000B SBR09); identical set at PDF p.97 (2320)"},
}

// String returns the code as the file sent it.
func (v ClaimFilingIndicatorSBR09) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v ClaimFilingIndicatorSBR09) Known() bool {
	return knownCode(claimFilingIndicatorSBR09Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v ClaimFilingIndicatorSBR09) ShortDescription() string {
	return shortOf(claimFilingIndicatorSBR09Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v ClaimFilingIndicatorSBR09) Description() string {
	return descriptionOf(claimFilingIndicatorSBR09Codes, string(v))
}

// ParseClaimFilingIndicatorSBR09 reads one SBR09 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseClaimFilingIndicatorSBR09(code string) (ClaimFilingIndicatorSBR09, error) {
	v := ClaimFilingIndicatorSBR09(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing SBR09: a value not on the list")
}

// SexDMG03 is DMG03, the Gender Code on a DMG segment, shared by the
// 2010BA subscriber DMG and the 2010CA patient DMG (identical shape).
type SexDMG03 string

const (
	SexDMG03F SexDMG03 = "F"
	SexDMG03M SexDMG03 = "M"
	SexDMG03U SexDMG03 = "U"
)

// The codes the guide allows for DMG03, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var sexDMG03Codes = []codeEntry{
	{"F", "Female", "Female — Transmittal 907, PDF p.17"},
	{"M", "Male", "Male — Transmittal 907, PDF p.17"},
	{"U", "Unknown", "Unknown — Transmittal 907, PDF p.17"},
}

// String returns the code as the file sent it.
func (v SexDMG03) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v SexDMG03) Known() bool { return knownCode(sexDMG03Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v SexDMG03) ShortDescription() string { return shortOf(sexDMG03Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v SexDMG03) Description() string { return descriptionOf(sexDMG03Codes, string(v)) }

// ParseSexDMG03 reads one DMG03 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseSexDMG03(code string) (SexDMG03, error) {
	v := SexDMG03(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing DMG03: a value not on the list")
}
