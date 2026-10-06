package x12837i

// SubscriberSBR is the SBR segment: the same nine-element shape at 2000B
// (subscriber) and 2320 (other subscriber information). Element names and
// citations are frontend/src/data/x12Names837i.ts's SBR_EL_NAME array
// (source: Transmittal 907, PDF p.16, identical shape at 2320, PDF p.97
// per the ts file's SBR_OTHER_EL_NAME note). Set is false when the segment
// did not appear at all.
type SubscriberSBR struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	PayerResponsibilitySBR01 PayerResponsibilitySBR01 `json:"payerResponsibilitySBR01,omitempty"` // Payer Responsibility Sequence Number Code — Transmittal 907, PDF p.16
	RelationshipSBR02        RelationshipSBR02        `json:"relationshipSBR02,omitempty"`        // Individual Relationship Code — Transmittal 907, PDF p.16

	GroupNumberSBR03 string `json:"groupNumberSBR03,omitempty"` // Subscriber Group or Policy Number — Transmittal 907, PDF p.16
	GroupNameSBR04   string `json:"groupNameSBR04,omitempty"`   // Subscriber Group Name — Transmittal 907, PDF p.16

	InsuranceTypeSBR05 string `json:"insuranceTypeSBR05,omitempty"` // Insurance Type Code (N/U) — Transmittal 907, PDF p.16
	CoordinationSBR06  string `json:"coordinationSBR06,omitempty"`  // Coordination of Benefits Code (N/U) — Transmittal 907, PDF p.16
	ResponseSBR07      string `json:"responseSBR07,omitempty"`      // Yes/No Condition or Response Code (N/U) — Transmittal 907, PDF p.16
	EmploymentSBR08    string `json:"employmentSBR08,omitempty"`    // Employment Status Code (N/U) — Transmittal 907, PDF p.16

	ClaimFilingIndicatorSBR09 ClaimFilingIndicatorSBR09 `json:"claimFilingIndicatorSBR09,omitempty"` // Claim Filing Indicator Code — Transmittal 907, PDF p.16
}
