package x12837p

// Subscriber- and patient-level coded elements: SBR (2000B Subscriber
// Information, reused by 2320 Other Subscriber Information) and PAT (2000C
// Patient Information, not used by Medicare).

// PayerResponsibilitySBR01 is SBR01, the Payer Responsibility Sequence Number
// Code, on the 2000B Subscriber Information SBR segment and the 2320 Other
// Subscriber Information SBR segment.
//
// The list is every code the 837P implementation guide allows at SBR01: the
// primary, secondary, and tertiary payers, the fourth through eleventh, and
// unknown. Medicare restricts 2000B SBR01 to P and S; CEDI 837P CG, PDF p.33
// (printed p.25; Table 13), confirmed in the worked example "SBR*P*18*...MB"
// (CEDI 837P CG, PDF p.50, printed p.42, Figure 3). The other names are the
// guide's.
type PayerResponsibilitySBR01 string

const (
	PayerResponsibilitySBR01NotSent PayerResponsibilitySBR01 = ""
	PayerResponsibilitySBR01P       PayerResponsibilitySBR01 = "P"
	PayerResponsibilitySBR01S       PayerResponsibilitySBR01 = "S"
	PayerResponsibilitySBR01T       PayerResponsibilitySBR01 = "T"
	PayerResponsibilitySBR01A       PayerResponsibilitySBR01 = "A"
	PayerResponsibilitySBR01B       PayerResponsibilitySBR01 = "B"
	PayerResponsibilitySBR01C       PayerResponsibilitySBR01 = "C"
	PayerResponsibilitySBR01D       PayerResponsibilitySBR01 = "D"
	PayerResponsibilitySBR01E       PayerResponsibilitySBR01 = "E"
	PayerResponsibilitySBR01F       PayerResponsibilitySBR01 = "F"
	PayerResponsibilitySBR01G       PayerResponsibilitySBR01 = "G"
	PayerResponsibilitySBR01H       PayerResponsibilitySBR01 = "H"
	PayerResponsibilitySBR01U       PayerResponsibilitySBR01 = "U"
)

var payerResponsibilitySBR01List = codeList{
	{"P", "Primary", "Primary — CEDI 837P CG, PDF p.33 (printed p.25; Table 13); frontend/src/data/x12Names837p.ts SBR01_NAME"},
	{"S", "Secondary", "Secondary — CEDI 837P CG, PDF p.33 (printed p.25; Table 13); frontend/src/data/x12Names837p.ts SBR01_NAME"},
	{"T", "Tertiary", guide("Tertiary", "SBR01")},
	{"A", "Payer Responsibility Four", guide("Payer Responsibility Four", "SBR01")},
	{"B", "Payer Responsibility Five", guide("Payer Responsibility Five", "SBR01")},
	{"C", "Payer Responsibility Six", guide("Payer Responsibility Six", "SBR01")},
	{"D", "Payer Responsibility Seven", guide("Payer Responsibility Seven", "SBR01")},
	{"E", "Payer Responsibility Eight", guide("Payer Responsibility Eight", "SBR01")},
	{"F", "Payer Responsibility Nine", guide("Payer Responsibility Nine", "SBR01")},
	{"G", "Payer Responsibility Ten", guide("Payer Responsibility Ten", "SBR01")},
	{"H", "Payer Responsibility Eleven", guide("Payer Responsibility Eleven", "SBR01")},
	{"U", "Unknown", guide("Unknown", "SBR01")},
}

func (v PayerResponsibilitySBR01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v PayerResponsibilitySBR01) Known() bool { return payerResponsibilitySBR01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v PayerResponsibilitySBR01) ShortDescription() string {
	return payerResponsibilitySBR01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v PayerResponsibilitySBR01) Description() string {
	return payerResponsibilitySBR01List.description(string(v))
}

// ParsePayerResponsibilitySBR01 returns text as a PayerResponsibilitySBR01,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParsePayerResponsibilitySBR01(text string) (PayerResponsibilitySBR01, error) {
	return PayerResponsibilitySBR01(text), payerResponsibilitySBR01List.check("SBR01", text)
}

// RelationshipSBR02 is SBR02, the Individual Relationship Code, on the 2000B
// Subscriber Information SBR segment and the 2320 Other Subscriber
// Information SBR segment.
//
// The list is every code the 837P implementation guide allows at the two
// SBR02 positions: 18 (Self) alone at 2000B, since a patient who is not the
// subscriber is sent in the 2000C loop instead, and the nine codes at 2320.
// Under Medicare the subscriber is always the patient, so 2000B SBR02 is
// always 18; CEDI 837P CG, PDF p.33 (printed p.25; Table 13), confirmed in
// the worked example "SBR*P*18*...". The other names are the guide's.
type RelationshipSBR02 string

const (
	RelationshipSBR02NotSent RelationshipSBR02 = ""
	RelationshipSBR0218      RelationshipSBR02 = "18"
	RelationshipSBR0201      RelationshipSBR02 = "01"
	RelationshipSBR0219      RelationshipSBR02 = "19"
	RelationshipSBR0220      RelationshipSBR02 = "20"
	RelationshipSBR0221      RelationshipSBR02 = "21"
	RelationshipSBR0239      RelationshipSBR02 = "39"
	RelationshipSBR0240      RelationshipSBR02 = "40"
	RelationshipSBR0253      RelationshipSBR02 = "53"
	RelationshipSBR02G8      RelationshipSBR02 = "G8"
)

var relationshipSBR02List = codeList{
	{"18", "Self", "Self — the only value Medicare's own subscriber loop sends, and the only value the guide allows at 2000B; CEDI 837P CG, PDF p.33 (printed p.25; Table 13); frontend/src/data/x12Names837p.ts SBR02_NAME"},
	{"01", "Spouse", guide("Spouse", "SBR02 at 2320")},
	{"19", "Child", guide("Child", "SBR02 at 2320")},
	{"20", "Employee", guide("Employee", "SBR02 at 2320")},
	{"21", "Unknown", guide("Unknown", "SBR02 at 2320")},
	{"39", "Organ Donor", guide("Organ Donor", "SBR02 at 2320")},
	{"40", "Cadaver Donor", guide("Cadaver Donor", "SBR02 at 2320")},
	{"53", "Life Partner", guide("Life Partner", "SBR02 at 2320")},
	{"G8", "Other Relationship", guide("Other Relationship", "SBR02 at 2320")},
}

func (v RelationshipSBR02) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v RelationshipSBR02) Known() bool { return relationshipSBR02List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v RelationshipSBR02) ShortDescription() string { return relationshipSBR02List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v RelationshipSBR02) Description() string { return relationshipSBR02List.description(string(v)) }

// ParseRelationshipSBR02 returns text as a RelationshipSBR02, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseRelationshipSBR02(text string) (RelationshipSBR02, error) {
	return RelationshipSBR02(text), relationshipSBR02List.check("SBR02", text)
}

// RelationshipPAT01 is PAT01, the Individual Relationship Code, on the 2000C
// Patient Information PAT segment.
//
// The list is every code the 837P implementation guide allows at PAT01. The
// 2000C Patient Hierarchical Level and its PAT segment are not used by
// Medicare at all (frontend/src/data/x12Names837p.ts NOT_USED_BY_MEDICARE:
// "2000C (Patient Hierarchical Level, in its entirety..."), and PAT01-04 are
// Not Used at 2000B as well, since the subscriber is always the patient
// (frontend/src/data/x12Names837p.ts PAT_EL_NAME). Transmittal 579's Values
// column lists the 2000C PAT01 codes (PDF p.22, "01, 19, 20, 21, 39, 40, 53,
// G8"); the names are the guide's.
type RelationshipPAT01 string

const (
	RelationshipPAT01NotSent RelationshipPAT01 = ""
	RelationshipPAT0101      RelationshipPAT01 = "01"
	RelationshipPAT0119      RelationshipPAT01 = "19"
	RelationshipPAT0120      RelationshipPAT01 = "20"
	RelationshipPAT0121      RelationshipPAT01 = "21"
	RelationshipPAT0139      RelationshipPAT01 = "39"
	RelationshipPAT0140      RelationshipPAT01 = "40"
	RelationshipPAT0153      RelationshipPAT01 = "53"
	RelationshipPAT01G8      RelationshipPAT01 = "G8"
)

var relationshipPAT01List = codeList{
	{"01", "Spouse", "Spouse — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"19", "Child", "Child — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"20", "Employee", "Employee — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"21", "Unknown", "Unknown — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"39", "Organ Donor", "Organ Donor — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"40", "Cadaver Donor", "Cadaver Donor — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"53", "Life Partner", "Life Partner — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
	{"G8", "Other Relationship", "Other Relationship — Transmittal 579, PDF p.22 (2000C PAT01 Values column); the name is the guide's: " + tr3},
}

func (v RelationshipPAT01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v RelationshipPAT01) Known() bool { return relationshipPAT01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v RelationshipPAT01) ShortDescription() string { return relationshipPAT01List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v RelationshipPAT01) Description() string { return relationshipPAT01List.description(string(v)) }

// ParseRelationshipPAT01 returns text as a RelationshipPAT01, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseRelationshipPAT01(text string) (RelationshipPAT01, error) {
	return RelationshipPAT01(text), relationshipPAT01List.check("PAT01", text)
}

// ClaimFilingIndicatorSBR09 is SBR09, the Claim Filing Indicator Code, on the
// 2000B Subscriber Information SBR segment and the 2320 Other Subscriber
// Information SBR segment.
//
// The list is every code the 837P implementation guide allows at SBR09.
// Medicare restricts 2000B SBR09 to MB; CEDI 837P CG, PDF p.33 (printed p.25;
// Table 13), confirmed in the worked example "SBR*P*18*...MB" (CEDI 837P CG,
// PDF p.50, printed p.42, Figure 3). Loop 2320's own SBR09 cannot equal MA or
// MB (CEDI 837P CG, PDF p.40, printed p.32, Table 27). The other names are the
// guide's.
type ClaimFilingIndicatorSBR09 string

const (
	ClaimFilingIndicatorSBR09NotSent ClaimFilingIndicatorSBR09 = ""
	ClaimFilingIndicatorSBR09MB      ClaimFilingIndicatorSBR09 = "MB"
	ClaimFilingIndicatorSBR09MA      ClaimFilingIndicatorSBR09 = "MA"
	ClaimFilingIndicatorSBR0911      ClaimFilingIndicatorSBR09 = "11"
	ClaimFilingIndicatorSBR0912      ClaimFilingIndicatorSBR09 = "12"
	ClaimFilingIndicatorSBR0913      ClaimFilingIndicatorSBR09 = "13"
	ClaimFilingIndicatorSBR0914      ClaimFilingIndicatorSBR09 = "14"
	ClaimFilingIndicatorSBR0915      ClaimFilingIndicatorSBR09 = "15"
	ClaimFilingIndicatorSBR0916      ClaimFilingIndicatorSBR09 = "16"
	ClaimFilingIndicatorSBR0917      ClaimFilingIndicatorSBR09 = "17"
	ClaimFilingIndicatorSBR09AM      ClaimFilingIndicatorSBR09 = "AM"
	ClaimFilingIndicatorSBR09BL      ClaimFilingIndicatorSBR09 = "BL"
	ClaimFilingIndicatorSBR09CH      ClaimFilingIndicatorSBR09 = "CH"
	ClaimFilingIndicatorSBR09CI      ClaimFilingIndicatorSBR09 = "CI"
	ClaimFilingIndicatorSBR09DS      ClaimFilingIndicatorSBR09 = "DS"
	ClaimFilingIndicatorSBR09FI      ClaimFilingIndicatorSBR09 = "FI"
	ClaimFilingIndicatorSBR09HM      ClaimFilingIndicatorSBR09 = "HM"
	ClaimFilingIndicatorSBR09LM      ClaimFilingIndicatorSBR09 = "LM"
	ClaimFilingIndicatorSBR09MC      ClaimFilingIndicatorSBR09 = "MC"
	ClaimFilingIndicatorSBR09OF      ClaimFilingIndicatorSBR09 = "OF"
	ClaimFilingIndicatorSBR09TV      ClaimFilingIndicatorSBR09 = "TV"
	ClaimFilingIndicatorSBR09VA      ClaimFilingIndicatorSBR09 = "VA"
	ClaimFilingIndicatorSBR09WC      ClaimFilingIndicatorSBR09 = "WC"
	ClaimFilingIndicatorSBR09ZZ      ClaimFilingIndicatorSBR09 = "ZZ"
)

var claimFilingIndicatorSBR09List = codeList{
	{"MB", "Medicare Part B", "Medicare Part B — CEDI 837P CG, PDF p.33 (printed p.25; Table 13); frontend/src/data/x12Names837p.ts SBR09_NAME"},
	{"MA", "Medicare Part A", guide("Medicare Part A", "SBR09")},
	{"11", "Other Non-Federal Programs", guide("Other Non-Federal Programs", "SBR09")},
	{"12", "Preferred Provider Organization (PPO)", guide("Preferred Provider Organization (PPO)", "SBR09")},
	{"13", "Point of Service (POS)", guide("Point of Service (POS)", "SBR09")},
	{"14", "Exclusive Provider Organization (EPO)", guide("Exclusive Provider Organization (EPO)", "SBR09")},
	{"15", "Indemnity Insurance", guide("Indemnity Insurance", "SBR09")},
	{"16", "Health Maintenance Organization (HMO) Medicare Risk", guide("Health Maintenance Organization (HMO) Medicare Risk", "SBR09")},
	{"17", "Dental Maintenance Organization", guide("Dental Maintenance Organization", "SBR09")},
	{"AM", "Automobile Medical", guide("Automobile Medical", "SBR09")},
	{"BL", "Blue Cross/Blue Shield", guide("Blue Cross/Blue Shield", "SBR09")},
	{"CH", "Champus", guide("Champus", "SBR09")},
	{"CI", "Commercial Insurance Co.", guide("Commercial Insurance Co.", "SBR09")},
	{"DS", "Disability", guide("Disability", "SBR09")},
	{"FI", "Federal Employees Program", guide("Federal Employees Program", "SBR09")},
	{"HM", "Health Maintenance Organization", guide("Health Maintenance Organization", "SBR09")},
	{"LM", "Liability Medical", guide("Liability Medical", "SBR09")},
	{"MC", "Medicaid", guide("Medicaid", "SBR09")},
	{"OF", "Other Federal Program", guide("Other Federal Program", "SBR09")},
	{"TV", "Title V", guide("Title V", "SBR09")},
	{"VA", "Veterans Affairs Plan", guide("Veterans Affairs Plan", "SBR09")},
	{"WC", "Workers' Compensation Health Claim", guide("Workers' Compensation Health Claim", "SBR09")},
	{"ZZ", "Mutually Defined", guide("Mutually Defined", "SBR09")},
}

func (v ClaimFilingIndicatorSBR09) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ClaimFilingIndicatorSBR09) Known() bool {
	return claimFilingIndicatorSBR09List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ClaimFilingIndicatorSBR09) ShortDescription() string {
	return claimFilingIndicatorSBR09List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ClaimFilingIndicatorSBR09) Description() string {
	return claimFilingIndicatorSBR09List.description(string(v))
}

// ParseClaimFilingIndicatorSBR09 returns text as a
// ClaimFilingIndicatorSBR09, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseClaimFilingIndicatorSBR09(text string) (ClaimFilingIndicatorSBR09, error) {
	return ClaimFilingIndicatorSBR09(text), claimFilingIndicatorSBR09List.check("SBR09", text)
}
