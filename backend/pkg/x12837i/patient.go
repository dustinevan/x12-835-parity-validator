package x12837i

// ---------------------------------------------------------------------------
// Patient: 2000C (HL, PAT), 2010CA (patient NM1, N3, N4, DMG, REF).
//
// Transmittal 907 marks every PAT element except PAT01 Not Used, and marks
// HL03 value 23 ("Patient") itself as not sent for Medicare, where the
// patient is always the subscriber and loop 2000C does not appear
// (HL03_NAME's own note, PDF p.20). This file models the loop fully anyway,
// for other payers' 837Is and for completeness of the type, same reasoning
// as x12835 modeling MIA/MOA/TS2 segments that only some claims carry.
// ---------------------------------------------------------------------------

// PatientRelationshipPAT holds one PAT segment. Every element except PAT01
// (Patient Location Code, Employment Status Code, Student Status Code, the
// DTP-shaped death-date pair, Unit or Basis for Measurement Code, Patient
// Weight, Pregnancy Indicator) is Not Used in this flat file; they are kept
// as sent in ExtraElements. Set is false when the segment did not appear
// at all. Transmittal 907, PDF p.20 (2000C).
type PatientRelationshipPAT struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	RelationshipPAT01 RelationshipSBR02 `json:"relationshipPAT01,omitempty"` // Individual Relationship Code; the shared type also used for SBR02 (see subscriber.go) — PAT01_VALUES excludes "18" (Self), which cannot apply here

	// ExtraElements holds PAT02 through PAT09 as sent. X12 defines them;
	// the 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// PatientHL is the 2000C loop: the HL and PAT segments, and the patient's
// own NM1/N3/N4/DMG/REF (2010CA).
type PatientHL struct {
	Hierarchy HierarchyHL            `json:"hierarchy,omitempty"` // HL; Transmittal 907, PDF p.20
	Patient   PatientRelationshipPAT `json:"patient,omitempty"`   // PAT; PDF p.20

	// 2010CA. Transmittal 907, PDF p.20-23.
	PatientNM1QC NameNM1        `json:"patientNM1QC,omitempty"` // NM1*QC, Patient; PDF p.20-21 (NM1_PATIENT_EL_NAME — the ID Code Qualifier and Patient Primary Identifier positions are Not Used in this flat file)
	Address      AddressN3      `json:"address,omitempty"`      // N3
	Location     LocationN4     `json:"location,omitempty"`     // N4
	Demographics DemographicDMG `json:"demographics,omitempty"` // DMG; PDF p.22, identical shape to 2010BA (see subscriber.go's DemographicDMG)

	// References (2010CA REF): 1W and SY both carry the combined name
	// "Property and Casualty Patient Identifier" in REF_PATIENT_NAME, while
	// Y4 carries its own distinct name, "Property and Casualty Claim
	// Number"; all three are kept generically here rather than split into
	// individually named fields, since this loop is not used for Medicare.
	// PDF p.22-23.
	References []ReferenceREF `json:"references,omitempty"`

	Claims []ClaimCLM `json:"claims,omitempty"` // the 2300 loops that follow this loop

	Raw []RawSegment `json:"raw,omitempty"` // segments within 2000C/2010CA this model does not place
}
