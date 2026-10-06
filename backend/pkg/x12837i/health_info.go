package x12837i

// HealthInfoHI is the Healthcare Code Information composite (C022), reused
// for every HI01 through HI12 position in every HI segment instance in the
// transaction: principal/admitting/other/external-cause-of-injury
// diagnosis, patient's reason for visit, DRG code, principal/other
// procedure, occurrence code, occurrence span code, value code, condition
// code, and treatment code. Element names and citations are
// frontend/src/data/x12Names837i.ts's HI_EL_NAME array (source: Transmittal
// 907, PDF p.24, positions 1-5, representative of every HI0n composite) and
// HI_POA_EL_NAME (PDF p.30, position 9).
//
// One HI segment carries up to twelve of these, and the parser places each
// in the claim field its qualifier names, so Index (the segment's position
// in the file) and Position (the composite's element number, 1 through 12)
// together say where it came from and let Write put the segment back
// together.
//
// QuantityHI01_6 (HI0n-6), the Version Identifier (HI0n-7), and the
// Industry Code (HI0n-8) are Not Used throughout this transaction per the
// flat file and are kept as plain text (QuantityHI01_6 is numeric, so it is
// kept as an Amount), as the project's other N/U elements are.
// PresentOnAdmissionHI0n_9 applies only on the diagnosis-bearing composites
// (principal/admitting/other diagnosis, external cause of injury); it is
// empty elsewhere.
//
// DateQualifierHI01_3 (the Date Format Qualifier sub-element, "D8" where a
// date is sent) is kept as a plain string: no coded type for it is in this
// package's catalog.
type HealthInfoHI struct {
	Set      bool `json:"set,omitempty"`
	Index    int  `json:"index,omitempty"`    // the HI segment's position in the file
	Position int  `json:"position,omitempty"` // this composite's element number in that segment: HI01 is 1, HI12 is 12

	QualifierHI01_1 DiagnosisQualifierHI01_1 `json:"qualifierHI01_1,omitempty"` // Diagnosis/Procedure/Code Type Qualifier — Transmittal 907, PDF p.24
	CodeHI01_2      string                   `json:"codeHI01_2,omitempty"`      // Code — Transmittal 907, PDF p.24; NUBC-maintained code values (occurrence/occurrence span/value/condition code) are kept as bare text, not enumerated, per this package's NUBC note

	DateQualifierHI01_3 string `json:"dateQualifierHI01_3,omitempty"` // Date Format Qualifier — Transmittal 907, PDF p.24
	DateHI01_4          Date   `json:"dateHI01_4,omitempty"`          // Date — Transmittal 907, PDF p.24; for a D8 date the date, for an RD8 range its first date

	// DateThroughHI01_4 is an RD8 range's second date (an occurrence span
	// code, BI, sends From-Through dates in one HI0n-4). Set is false for a
	// single D8 date, the same convention as DateDTP's DateThrough.
	DateThroughHI01_4 Date `json:"dateThroughHI01_4,omitempty"`

	AmountHI01_5   Amount `json:"amountHI01_5,omitempty"`   // Monetary Amount — Transmittal 907, PDF p.24
	QuantityHI01_6 Amount `json:"quantityHI01_6,omitempty"` // Quantity (N/U) — Transmittal 907, PDF p.24

	VersionIdentifierHI0n_7 string `json:"versionIdentifierHI0n_7,omitempty"` // Version Identifier (N/U) — Transmittal 907, PDF p.24
	IndustryCodeHI0n_8      string `json:"industryCodeHI0n_8,omitempty"`      // Industry Code (N/U) — Transmittal 907, PDF p.24

	PresentOnAdmissionHI0n_9 PresentOnAdmissionHI0n_9 `json:"presentOnAdmissionHI0n_9,omitempty"` // Present on Admission Indicator — Transmittal 907, PDF p.30 (diagnosis composites only)
}
