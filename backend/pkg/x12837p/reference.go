package x12837p

// ReferenceREF holds one REF segment, whether a loop placed it in a
// qualifier-specific field or in a catch-all list; the qualifier travels
// with the value either way, so no meaning is lost. CLAIM_REF_NAME,
// SERVICE_REF_NAME, LOOP_REF_NAME, and the other *_REF_NAME tables in
// frontend/src/data/x12Names837p.ts name what each qualifier means at its
// loop. REF03 and REF04 are Not Used at every 837P REF, and a file that
// sends them keeps their text in Extra.
//
// ReferenceIdentificationQualifierREF01 is defined in enums_reference.go.
type ReferenceREF struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	QualifierREF01 ReferenceIdentificationQualifierREF01 `json:"qualifierREF01,omitempty"` // Reference Identification Qualifier
	Value          string                                `json:"value,omitempty"`          // REF02; kept as text even when all digits

	Extra []string `json:"extra,omitempty"` // elements after REF02, kept as sent
}

// DateDTP holds one DTP segment whose qualifier is not already fixed by the
// field name holding it. DTP01 is always "Date/Time Qualifier" and DTP02 is
// always "Date Time Period Format Qualifier" (DTP01_EL_NAME/DTP02_EL_NAME
// in frontend/src/data/x12Names837p.ts, every DTP instance); what DTP03
// means for a given qualifier is CLAIM_DTP_NAME/CLAIM_DTP03_EL_NAME,
// SERVICE_DTP_NAME/SERVICE_DTP03_EL_NAME, and the other *_DTP_NAME tables
// in that file.
//
// FormatDTP02 is kept as received ("D8" for a single CCYYMMDD date, "RD8"
// for a CCYYMMDD-CCYYMMDD range; a handful of boxes may send other
// qualifiers) because the 837P reuses this one segment shape for both
// single dates and ranges. Date is the parsed value only when FormatDTP02
// is "D8"; it is zero (Set false) for a range or any other format, whose
// text is read from RawDTP03 instead. RawDTP03 always holds DTP03 exactly
// as received, independent of whether Date could be parsed, so no exact
// evidence depends on this struct's own interpretation of the format.
//
// DateTimeQualifierDTP01 is defined in enums_reference.go.
type DateDTP struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	QualifierDTP01 DateTimeQualifierDTP01 `json:"qualifierDTP01,omitempty"` // Date/Time Qualifier
	FormatDTP02    string                 `json:"formatDTP02,omitempty"`    // Date Time Period Format Qualifier, e.g. "D8" or "RD8"
	Date           Date                   `json:"date,omitempty"`           // parsed only when FormatDTP02 is "D8"
	RawDTP03       string                 `json:"rawDTP03,omitempty"`       // DTP03 exactly as received

	Extra []string `json:"extra,omitempty"` // elements after DTP03, kept as sent
}

// AmountAMT holds one AMT segment. AMT01_EL_NAME in
// frontend/src/data/x12Names837p.ts gives the generic "Amount Qualifier
// Code" position name, and AMT02_EL_NAME_BY_BOX/SERVICE_AMT02_EL_NAME_BY_BOX
// name what AMT02 means at a given loop and qualifier. AMT03 is Not Used at
// every 837P AMT, and a file that sends it keeps its text in Extra.
//
// AmountQualifierAMT01 is defined in enums_reference.go.
type AmountAMT struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	QualifierAMT01 AmountQualifierAMT01 `json:"qualifierAMT01,omitempty"` // Amount Qualifier Code
	Amount         Amount               `json:"amount,omitempty"`         // AMT02

	Extra []string `json:"extra,omitempty"` // elements after AMT02, kept as sent
}

// QuantityQTY holds one QTY segment (AMBULANCE_PATIENT_COUNT_QTY_NAME,
// OBSTETRIC_ANESTHESIA_QTY_NAME in frontend/src/data/x12Names837p.ts).
// QTY01_EL_NAME gives the generic "Quantity Qualifier" position name, and
// QTY02_EL_NAME_BY_BOX names what QTY02 means at a given qualifier. QTY03
// and QTY04 are Not Used, and a file that sends them keeps their text in
// Extra.
//
// QuantityQualifierQTY01 is defined in enums_reference.go.
type QuantityQTY struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	QualifierQTY01 QuantityQualifierQTY01 `json:"qualifierQTY01,omitempty"` // Quantity Qualifier
	Quantity       Amount                 `json:"quantity,omitempty"`       // QTY02

	Extra []string `json:"extra,omitempty"` // elements after QTY02, kept as sent
}

// FileInfoK3 holds one K3 segment, "File Information" (2300 and 2400, up to
// 10 each): one fixed-format text record in K301. K302 and K303 are Not
// Used, and a file that sends them keeps their text in Extra.
// CLAIM_K3_NAME, SERVICE_K3_NAME, K3_EL_NAME; Transmittal 907, PDF p.166,
// 220.
type FileInfoK3 struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	TextK301 string `json:"textK301,omitempty"` // Fixed Format Information

	Extra []string `json:"extra,omitempty"` // elements after K301, kept as sent
}
