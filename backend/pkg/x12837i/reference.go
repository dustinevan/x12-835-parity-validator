package x12837i

// ReferenceREF holds one REF segment. The loops that give a qualifier its
// own field (the billing provider's REF*EI, the claim's REF*G1, and the
// rest) hold one of these there too, so the segment's position and its
// qualifier travel with the value; every other REF goes in the loop's
// References, where the qualifier says what the value means.
// frontend/src/data/x12Names837i.ts names what each qualifier means at each
// level (REF_CLAIM_NAME, PROVIDER_REF_NAME, REF_SUBSCRIBER_NAME,
// REF_PAYER_NAME, REF_LINE_NAME, and the rest of that file's per-loop REF
// tables).
//
// Unlike backend/pkg/x12835's ReferenceREF, QualifierREF01 is kept as a
// plain string, not a coded type: this transaction's REF qualifiers
// (0B/1G/G2, 2U/EI/FY/NF, 6R/9B/9D, Y4, and the rest) vary by loop with no
// single combined code list the way x12835's one merged RefQualifierREF01
// covers every 835 REF, and no REF type is in
// frontend/src/data/x12Names837i.ts's list of implied coded types.
type ReferenceREF struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	QualifierREF01 string `json:"qualifierREF01,omitempty"` // Reference Identification Qualifier
	Value          string `json:"value,omitempty"`          // REF02; kept as text even when all digits

	// ExtraElements holds REF03 (Description) and REF04 (the Reference
	// Identifier composite, as one text) as sent. X12 defines them; the
	// 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// DateDTP holds one DTP segment: DTP01 (the qualifier), DTP02 (the Date
// Time Period Format Qualifier), and DTP03 read in the form DTP02 names.
// The 837I sends four forms, and one struct holds them all so a DTP is one
// thing wherever it sits:
//
//   - D8, one CCYYMMDD date, in Date.
//   - RD8, a CCYYMMDD-CCYYMMDD range, in Date and DateThrough (for example
//     qualifier 434, Statement From-Through Dates).
//   - DT, CCYYMMDDHHMM, a date and an hour, in Date and Time (qualifier
//     435, Admission Date and Hour).
//   - TM, HHMM, an hour alone, in Time (qualifier 096, Discharge Hour).
//
// The fields the form does not use keep Set false. The loops that give a
// qualifier its own field hold one of these there (ClaimCLM's four DTPs,
// the service line's DTP*472, the two DTP*573s); every other DTP goes in a
// Dates list or is kept raw. Set is false when the segment did not appear
// at all. DTP01 is always "Date/Time Qualifier" (Transmittal 907, PDF
// p.24), consistent across every DTP instance in the transaction.
type DateDTP struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	QualifierDTP01 DateQualifierDTP01 `json:"qualifierDTP01,omitempty"` // Date/Time Qualifier — Transmittal 907, PDF p.24
	FormatDTP02    string             `json:"formatDTP02,omitempty"`    // Date Time Period Format Qualifier: D8, RD8, DT, or TM — Transmittal 907, PDF p.24
	Date           Date               `json:"date,omitempty"`           // DTP03's date: the D8 date, the RD8 range's first date, or the DT date
	DateThrough    Date               `json:"dateThrough,omitempty"`    // the RD8 range's second date; Set is false otherwise
	Time           Time               `json:"time,omitempty"`           // the DT hour and minute, or the TM time; Set is false otherwise
}

// AmountAMT holds one AMT segment. frontend/src/data/x12Names837i.ts's
// AMT01_NAME shows the qualifier's meaning varies by loop and qualifier
// (claim-level F3/GT/N8, COB-level D/EAF/A8); no AMT01 type is in that
// file's list of implied coded types, so QualifierAMT01 is kept as a plain
// string, unlike backend/pkg/x12835's AmountAMT. The loops that give a
// qualifier its own field hold one of these there; every other AMT goes in
// the loop's Amounts.
type AmountAMT struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	QualifierAMT01 string `json:"qualifierAMT01,omitempty"` // Amount Qualifier Code
	Amount         Amount `json:"amount,omitempty"`

	// ExtraElements holds AMT03 (Credit/Debit Flag Code) as sent. X12
	// defines it; the 837I guide marks it Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}
