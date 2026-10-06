package x12835

// ProcedureSVC01 is the composite medical procedure identifier carried by
// SVC01 and, in the same shape, by SVC06 (the original procedure, before
// adjudication). CMS 835 Flat File, PDF p.11 names the composite
// "Composite Medical Procedure Identifier" but does not give its
// sub-components individual names; the shape below (qualifier, procedure
// code, up to four modifiers, description, and a second code) is the
// composite's standard layout (C003), not an individually cited name. The
// 835 guide marks the eighth component not used; a file that sends it
// keeps it.
type ProcedureSVC01 struct {
	Qualifier   ProductQualifierSVC01 `json:"qualifier,omitempty"`   // code-system qualifier, e.g. HC = Healthcare Common Procedure Coding System (SVC_QUALIFIER_NAME in frontend/src/data/x12Names835.ts)
	Code        string                `json:"code,omitempty"`        // the procedure code; CPT codes are kept as their code (AMA-licensed, see docs/design/835-viewer.md)
	Modifiers   [4]string             `json:"modifiers,omitempty"`   // up to four procedure modifiers
	Description string                `json:"description,omitempty"` // free-text description, sent in place of a code at the payer's discretion
	Code2       string                `json:"code2,omitempty"`       // C003-08 Product/Service ID, the second code; not used by the 835 guide, kept as sent
}

// ServiceLineSVC holds one SVC segment (loop 2110) and everything the 835
// nests under it.
type ServiceLineSVC struct {
	Index int `json:"index,omitempty"` // the SVC segment's position in the file

	ProcedureSVC01   ProcedureSVC01 `json:"procedureSVC01,omitempty"`   // SVC01; CMS 835 Flat File, PDF p.11
	ChargeSVC02      Amount         `json:"chargeSVC02,omitempty"`      // Line Item Charge Amount; PDF p.11
	PaymentSVC03     Amount         `json:"paymentSVC03,omitempty"`     // Line Item Provider Payment; PDF p.11
	RevenueCodeSVC04 string         `json:"revenueCodeSVC04,omitempty"` // NUBC Revenue Code; PDF p.11 (external code list)
	UnitsSVC05       Amount         `json:"unitsSVC05,omitempty"`       // Units of Service Paid Count; PDF p.12

	// SVC06 is the same composite shape as SVC01, holding the original
	// (pre-adjudication) procedure; SVC07 is its matching original unit count.
	OriginalProcedureSVC06 ProcedureSVC01 `json:"originalProcedureSVC06,omitempty"` // PDF p.12
	OriginalUnitsSVC07     Amount         `json:"originalUnitsSVC07,omitempty"`     // Original Units of Service Count; PDF p.12

	// CMS 835 Flat File, PDF p.12 names the loop 2110 DTM segment and its
	// DTM02 "Service Date" and lists 150, 151, 472 as the DTM01 values. The
	// page names no qualifier individually, so all three fields carry the
	// segment's name and the field name carries the qualifier. Each is
	// present when its DateQualifierDTM01 is not empty.
	ServiceDateDTM472 DateTimeDTM `json:"serviceDateDTM472,omitempty"` // DTM01 472; CMS 835 Flat File, PDF p.12
	ServiceDateDTM150 DateTimeDTM `json:"serviceDateDTM150,omitempty"` // DTM01 150; CMS 835 Flat File, PDF p.12
	ServiceDateDTM151 DateTimeDTM `json:"serviceDateDTM151,omitempty"` // DTM01 151; CMS 835 Flat File, PDF p.12

	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"` // CAS, "Service Adjustment" at this level (SERVICE_CAS_NAME)

	// REF qualifiers SERVICE_REF_NAME names individually at this level, each
	// the first REF with that qualifier and a value; the field name carries
	// the qualifier. A field is present when its Value is not empty.
	LineControlREF6R ReferenceREF   `json:"lineControlREF6R,omitempty"` // Line Item Control Number; CMS 835 Flat File, PDF p.13
	ServiceIDREFLU   ReferenceREF   `json:"serviceIDREFLU,omitempty"`   // Service Identification; PDF p.12
	RenderingREFHPI  ReferenceREF   `json:"renderingREFHPI,omitempty"`  // Rendering Provider Information; PDF p.13
	PolicyREF0K      ReferenceREF   `json:"policyREF0K,omitempty"`      // Health Care Policy Identification; PDF p.13
	References       []ReferenceREF `json:"references,omitempty"`       // any loop 2110 REF not one of the qualified fields above, or a repeat of one

	Amounts    []AmountAMT   `json:"amounts,omitempty"`    // AMT, "Service Supplemental Amount" (SERVICE_AMT_NAME); PDF p.13
	Quantities []QuantityQTY `json:"quantities,omitempty"` // QTY; not documented at this level in the reviewed sources (CLAIM_QTY_NAME only covers loop 2100) — kept for symmetry and any file that sends one

	Remarks []RemarkLQ `json:"remarks,omitempty"` // LQ, qualifier HE, "Health Care Remarks Codes"; PDF p.13

	Raw []RawSegment `json:"raw,omitempty"` // segments within this service line the model does not place
}
