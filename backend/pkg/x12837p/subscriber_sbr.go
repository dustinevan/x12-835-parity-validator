package x12837p

// SubscriberSBR holds one SBR segment. The same nine-element shape is sent
// at both the 2000B Subscriber Information loop (SBR_EL_NAME in
// frontend/src/data/x12Names837p.ts) and the 2320 Other Subscriber
// Information loop (OTHER_SUBSCRIBER_SBR_EL_NAME), which labels SBR03/04
// slightly differently ("Insured Group or Policy Number" / "Other Insured
// Group Name" rather than "Subscriber Group or Policy Number" / "Subscriber
// Group Name") for the same two positions; the field holding this struct
// carries that distinction, not a field here. SBR06-08 are Not Used at
// every instance and kept as plain text. Set is false when the segment did
// not appear at all.
//
// PayerResponsibilitySBR01, RelationshipSBR02, and
// ClaimFilingIndicatorSBR09 are defined in enums_subscriber.go.
type SubscriberSBR struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	PayerResponsibilitySBR01 PayerResponsibilitySBR01 `json:"payerResponsibilitySBR01,omitempty"` // Payer Responsibility Sequence Number Code
	RelationshipSBR02        RelationshipSBR02        `json:"relationshipSBR02,omitempty"`        // Individual Relationship Code
	GroupOrPolicyNumberSBR03 string                   `json:"groupOrPolicyNumberSBR03,omitempty"` // Subscriber/Insured Group or Policy Number
	GroupNameSBR04           string                   `json:"groupNameSBR04,omitempty"`           // Subscriber/Other Insured Group Name
	InsuranceTypeSBR05       InsuranceTypeSBR05       `json:"insuranceTypeSBR05,omitempty"`       // Insurance Type Code

	CoordinationSBR06 string `json:"coordinationSBR06,omitempty"` // Coordination of Benefits Code (N/U) — Transmittal 907, PDF p.144
	ResponseSBR07     string `json:"responseSBR07,omitempty"`     // Yes/No Condition or Response Code (N/U) — Transmittal 907, PDF p.144
	EmploymentSBR08   string `json:"employmentSBR08,omitempty"`   // Employment Status Code (N/U) — Transmittal 907, PDF p.144

	ClaimFilingIndicatorSBR09 ClaimFilingIndicatorSBR09 `json:"claimFilingIndicatorSBR09,omitempty"` // Claim Filing Indicator Code

	Extra []string `json:"extra,omitempty"` // elements after SBR09, kept as sent
}

// InsuranceTypeSBR05 is SBR05, the Insurance Type Code, on the 2320 Other
// Subscriber Information SBR segment (Not Used at 2000B).
//
// The list is every code the 837P implementation guide allows at SBR05: the
// reasons Medicare is the secondary payer. Neither Transmittal 907 nor the
// CEDI companion guide restates SBR05's own code characters in
// frontend/src/data/x12Names837p.ts (SBR_EL_NAME and
// OTHER_SUBSCRIBER_SBR_EL_NAME give only the position's own name), so every
// code and name here is the guide's. SBR05 is not used to place a loop or
// box, and no other element shares its code list, so the type sits beside
// SubscriberSBR.
type InsuranceTypeSBR05 string

const (
	InsuranceTypeSBR05NotSent InsuranceTypeSBR05 = ""
	InsuranceTypeSBR0512      InsuranceTypeSBR05 = "12"
	InsuranceTypeSBR0513      InsuranceTypeSBR05 = "13"
	InsuranceTypeSBR0514      InsuranceTypeSBR05 = "14"
	InsuranceTypeSBR0515      InsuranceTypeSBR05 = "15"
	InsuranceTypeSBR0516      InsuranceTypeSBR05 = "16"
	InsuranceTypeSBR0541      InsuranceTypeSBR05 = "41"
	InsuranceTypeSBR0542      InsuranceTypeSBR05 = "42"
	InsuranceTypeSBR0543      InsuranceTypeSBR05 = "43"
	InsuranceTypeSBR0547      InsuranceTypeSBR05 = "47"
)

var insuranceTypeSBR05List = codeList{
	{"12", "Medicare Secondary Working Aged Beneficiary or Spouse with Employer Group Health Plan", guide("Medicare Secondary Working Aged Beneficiary or Spouse with Employer Group Health Plan", "SBR05 at 2320")},
	{"13", "Medicare Secondary End-Stage Renal Disease Beneficiary in the Mandated Coordination Period with an Employer's Group Health Plan", guide("Medicare Secondary End-Stage Renal Disease Beneficiary in the Mandated Coordination Period with an Employer's Group Health Plan", "SBR05 at 2320")},
	{"14", "Medicare Secondary, No-fault Insurance including Auto is Primary", guide("Medicare Secondary, No-fault Insurance including Auto is Primary", "SBR05 at 2320")},
	{"15", "Medicare Secondary Worker's Compensation", guide("Medicare Secondary Worker's Compensation", "SBR05 at 2320")},
	{"16", "Medicare Secondary Public Health Service (PHS) or Other Federal Agency", guide("Medicare Secondary Public Health Service (PHS) or Other Federal Agency", "SBR05 at 2320")},
	{"41", "Medicare Secondary Black Lung", guide("Medicare Secondary Black Lung", "SBR05 at 2320")},
	{"42", "Medicare Secondary Veteran's Administration", guide("Medicare Secondary Veteran's Administration", "SBR05 at 2320")},
	{"43", "Medicare Secondary Disabled Beneficiary Under Age 65 with Large Group Health Plan (LGHP)", guide("Medicare Secondary Disabled Beneficiary Under Age 65 with Large Group Health Plan (LGHP)", "SBR05 at 2320")},
	{"47", "Medicare Secondary, Other Liability Insurance is Primary", guide("Medicare Secondary, Other Liability Insurance is Primary", "SBR05 at 2320")},
}

func (v InsuranceTypeSBR05) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v InsuranceTypeSBR05) Known() bool { return insuranceTypeSBR05List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v InsuranceTypeSBR05) ShortDescription() string { return insuranceTypeSBR05List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v InsuranceTypeSBR05) Description() string {
	return insuranceTypeSBR05List.description(string(v))
}

// ParseInsuranceTypeSBR05 returns text as an InsuranceTypeSBR05, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseInsuranceTypeSBR05(text string) (InsuranceTypeSBR05, error) {
	return InsuranceTypeSBR05(text), insuranceTypeSBR05List.check("SBR05", text)
}
