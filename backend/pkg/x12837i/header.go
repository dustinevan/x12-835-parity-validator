package x12837i

// HeaderBHT is the BHT segment: the beginning of the hierarchical
// transaction, carried once per transaction set, directly under ST.
// Element names and citations are frontend/src/data/x12Names837i.ts's
// BHT_EL_NAME array (source: Transmittal 907, PDF p.8-9).
//
// StructureCodeBHT01 and PurposeCodeBHT02 are kept as plain strings, not
// coded types: frontend/src/data/x12Names837i.ts's own list of the coded
// types this file's data implies does not include them (BHT01 sends the
// single observed value "0019"; BHT02's values, "00"/"18", are recorded in
// BHT02_NAME but not modeled as a Go type here). Set is false when the
// segment did not appear at all.
type HeaderBHT struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	StructureCodeBHT01 string `json:"structureCodeBHT01,omitempty"` // Hierarchical Structure Code — Transmittal 907, PDF p.8 (value "0019")
	PurposeCodeBHT02   string `json:"purposeCodeBHT02,omitempty"`   // Transaction Set Purpose Code — Transmittal 907, PDF p.8 ("00" Original, "18" Reissue)

	OriginatorTransactionIDBHT03 string `json:"originatorTransactionIDBHT03,omitempty"` // Originator Application Transaction ID — Transmittal 907, PDF p.8

	CreationDateBHT04 Date `json:"creationDateBHT04,omitempty"` // Transaction Set Creation Date — Transmittal 907, PDF p.9
	CreationTimeBHT05 Time `json:"creationTimeBHT05,omitempty"` // Transaction Set Creation Time — Transmittal 907, PDF p.9

	ClaimOrEncounterBHT06 ClaimOrEncounterBHT06 `json:"claimOrEncounterBHT06,omitempty"` // Claim or Encounter ID — Transmittal 907, PDF p.9
}
