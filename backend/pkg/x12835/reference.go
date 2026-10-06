package x12835

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"

// Index, on every segment struct in this package, is the segment's position
// in the file, 0-based, the ISA being 0 (ParseClaim counts from the CLP,
// which is 0), the same count RawSegment.Index uses. Write puts a loop's
// segments back in Index order, so a file reads in and writes out in its
// own order; a struct built by hand with every Index zero writes in the
// implementation guide's order.

// ReferenceREF holds one REF segment: the qualifier, the value, and the two
// positions the standard defines after them, which the 835 guide marks not
// used but a file may send. HEADER_REF_NAME, CLAIM_REF_NAME, and
// SERVICE_REF_NAME in frontend/src/data/x12Names835.ts name what each
// qualifier means at its level. The service line keeps four qualifiers in
// fields of their own (LineControlREF6R and the others in service_line.go);
// everything else goes in a References slice, qualifier and all.
type ReferenceREF struct {
	Index             int               `json:"index,omitempty"`
	RefQualifierREF01 RefQualifierREF01 `json:"refQualifierREF01,omitempty"` // Reference Identification Qualifier
	Value             string            `json:"value,omitempty"`             // REF02; kept as text even when all digits
	DescriptionREF03  string            `json:"descriptionREF03,omitempty"`  // REF03 Description; not used by the 835 guide, kept as sent
	SubReferenceREF04 string            `json:"subReferenceREF04,omitempty"` // REF04 Reference Identifier (composite C040); not used by the 835 guide, kept as sent with its component separators
}

// DateTimeDTM holds one DTM segment: the qualifier, the date, and the four
// positions the standard defines after them (a time, a time code, and a
// period), which the 835 guide marks not used but a file may send. The
// header's DTM*405 and the claim's and service line's dated fields all use
// it, so a DTM writes back out whole. DTM01 is always "Date/Time Qualifier"
// (CMS 835 Flat File, PDF p.2, 10-12); what DTM02 means for a given
// qualifier is DTM_NAME / DTM02_EL_NAME_BY_QUALIFIER in
// frontend/src/data/x12Names835.ts. A field holding one is present when
// DateQualifierDTM01 is not empty.
type DateTimeDTM struct {
	Index              int                `json:"index,omitempty"`
	DateQualifierDTM01 DateQualifierDTM01 `json:"dateQualifierDTM01,omitempty"`
	Date               Date               `json:"date,omitempty"`
	TimeDTM03          string             `json:"timeDTM03,omitempty"`         // DTM03 Time; kept as sent
	TimeCodeDTM04      string             `json:"timeCodeDTM04,omitempty"`     // DTM04 Time Code; kept as sent
	PeriodFormatDTM05  string             `json:"periodFormatDTM05,omitempty"` // DTM05 Date Time Period Format Qualifier; kept as sent
	PeriodDTM06        string             `json:"periodDTM06,omitempty"`       // DTM06 Date Time Period; kept as sent
}

// AmountAMT holds one AMT segment whose qualifier is not resolved to an
// individual meaning at its level: CLAIM_AMT_NAME (loop 2100) and
// SERVICE_AMT_NAME (loop 2110) in frontend/src/data/x12Names835.ts show the
// segment carries one combined name regardless of AMT01's value.
type AmountAMT struct {
	Index                int                  `json:"index,omitempty"`
	AmountQualifierAMT01 AmountQualifierAMT01 `json:"amountQualifierAMT01,omitempty"`
	Amount               Amount               `json:"amount,omitempty"`
	CreditDebitAMT03     string               `json:"creditDebitAMT03,omitempty"` // AMT03 Credit/Debit Flag Code; not used by the 835 guide, kept as sent
}

// QuantityQTY holds one QTY segment (claim level, loop 2100;
// CLAIM_QTY_NAME in frontend/src/data/x12Names835.ts — the segment carries
// one combined name, "Claim Supplemental Information Quantity", regardless
// of QTY01's value).
type QuantityQTY struct {
	Index                  int                    `json:"index,omitempty"`
	QuantityQualifierQTY01 QuantityQualifierQTY01 `json:"quantityQualifierQTY01,omitempty"`
	Amount                 Amount                 `json:"amount,omitempty"`
	UnitOfMeasureQTY03     string                 `json:"unitOfMeasureQTY03,omitempty"` // QTY03 Composite Unit of Measure (C001); not used by the 835 guide, kept as sent with its component separators
	FreeFormQTY04          string                 `json:"freeFormQTY04,omitempty"`      // QTY04 Free-form Information; not used by the 835 guide, kept as sent
}

// RemarkLQ holds one LQ segment (service line health care remark codes;
// LQ_NAME in frontend/src/data/x12Names835.ts — qualifier HE, "Health Care
// Remarks Codes", is the one Medicare documents). LQ02 is read as a RARC,
// per docs/design/835-viewer.md's LQ instruction; under another qualifier
// it is kept as the code sent.
type RemarkLQ struct {
	Index               int                 `json:"index,omitempty"`
	RemarkQualifierLQ01 RemarkQualifierLQ01 `json:"remarkQualifierLQ01,omitempty"`
	RARC                rarc.Code           `json:"rarc,omitempty"`
}
