package x12837p

// SubmitterNM141 is the 1000A Submitter Name loop: the NM1*41 segment and
// its Submitter EDI Contact Information PER segment(s). Source: Transmittal
// 907, PDF p.135-136 (footer 3-4); LOOP_NM1_NAME["1000A"],
// NM1_ROLE_EL_NAME["1000A"], and SUBMITTER_PER_EL_NAME in
// frontend/src/data/x12Names837p.ts. NM101 is EntityIdentifierNM10141 on a
// compliant file.
type SubmitterNM141 struct {
	Name NameNM1 `json:"name,omitempty"` // NM1*41 — Transmittal 907, PDF p.135 (footer 3)

	// Contacts is the Submitter EDI Contact Information PER segment(s).
	// SUBMITTER_PER_EL_NAME's own nine positions include PER09 (Contact
	// Inquiry Reference), unlike every other PER instance in this package -
	// see ContactPER's InquiryReferencePER09 field.
	Contacts []ContactPER `json:"contacts,omitempty"` // Transmittal 907, PDF p.135-136 (footer 3-4)

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop that the model does not place
}
