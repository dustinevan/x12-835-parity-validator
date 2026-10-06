package x12837p

// ===========================================================================
// 2000C / 2010CA: Patient
// ===========================================================================
//
// Loop 2000C, including its PAT segment and the entire 2010CA name loop,
// is not used by Medicare at all — the subscriber is always the patient
// under Medicare, so Medicare claims never open this loop
// (NOT_USED_BY_MEDICARE, frontend/src/data/x12Names837p.ts, citing CEDI
// 837P CG, PDF p.26-27, printed p.18-19, Table 5). It is still modeled
// here, matching how backend/pkg/x12835 keeps its own Medicare-unsupported
// segments (CUR) typed, since this package models the general 837P
// structure, not only the Medicare-restricted subset of it.

// PatientNM1 holds the 2010CA loop: the patient's NM1 (role code QC,
// confirmed by frontend/src/data/x12Names837p.ts's NM101_CONFIRMED_NAME;
// this loop is not used by Medicare regardless), its N3/N4, its
// demographic information (DMG, DemographicDMG from subscriber.go), its
// Property and Casualty Patient Identifier REF, and its contact PER.
type PatientNM1 struct {
	NM1          NameNM1        `json:"nm1,omitempty"`          // Transmittal 907, PDF p.150 (footer 18); loop not used by Medicare
	Address      AddressN3      `json:"address,omitempty"`      // Transmittal 907, PDF p.150 (footer 18)
	Location     LocationN4     `json:"location,omitempty"`     // Transmittal 907, PDF p.151 (footer 19-20)
	Demographics DemographicDMG `json:"demographics,omitempty"` // same shape as 2010BA's; Transmittal 907 names this loop's own DMG at the same positions

	// References holds the Property and Casualty Patient Identifier REF
	// (PATIENT_PC_REF_NAME, PDF p.152-153 footer 20-21); its REF01 qualifier
	// is not confirmed in either reviewed source.
	References []ReferenceREF `json:"references,omitempty"`

	Contact ContactPER `json:"contact,omitempty"` // Property and Casualty Patient Contact Information; Transmittal 907, PDF p.153 (footer 21)

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop the model does not place
}

// PatientHL holds the 2000C loop: the HL, the PAT (PatientInfoPAT from
// subscriber.go; the whole loop is not used by Medicare), and the 2010CA
// name loop it contains. There is no SBR at this level — SBR lives only at
// 2000B; a dependent patient's coverage is still the subscriber's SBR.
type PatientHL struct {
	HL HierarchyHL `json:"hl,omitempty"` // Transmittal 907, PDF p.149 (footer 17); loop not used by Medicare

	PatientInfo  PatientInfoPAT `json:"patientInfo,omitempty"`  // PAT
	PatientNM1QC PatientNM1     `json:"patientNM1QC,omitempty"` // 2010CA

	Claims []ClaimCLM `json:"claims,omitempty"` // 2300 loops that follow this HL before the next HL

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop, before 2010CA, that the model does not place
}
