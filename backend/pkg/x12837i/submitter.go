package x12837i

// SubmitterNM141 is the 1000A loop: the submitter's NM1 segment (NM101
// fixed at "41") and its PER contacts. The NM1 uses the shared NameNM1
// shape, as every other NM1 loop in this package does. Its positions'
// public names at this loop are frontend/src/data/x12Names837i.ts's
// NM1_SUBMITTER_EL_NAME array (source: Transmittal 907, PDF p.9-10):
// NM103 Submitter Last or Organization Name, NM104 Submitter First Name,
// NM105 Submitter Middle Name or Initial, NM106 Name Prefix (N/U), NM107
// Name Suffix (N/U), NM108 Identification Code Qualifier, NM109 Submitter
// Identifier.
type SubmitterNM141 struct {
	Name NameNM1 `json:"name,omitempty"` // NM1*41 — Transmittal 907, PDF p.9

	Contacts []ContactPER `json:"contacts,omitempty"` // PER — Transmittal 907, PDF p.9-10

	Raw []RawSegment `json:"raw,omitempty"` // segments within the loop, after the NM1, that the model does not place
}
