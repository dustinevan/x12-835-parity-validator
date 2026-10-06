package x12837i

// TransactionST is one 837I transaction set: the ST header, the BHT
// beginning-of-transaction segment, the submitter and receiver loops, every
// billing provider hierarchical loop with its subscribers and claims, and
// the matching SE trailer. Element names and citations for ST and SE are
// frontend/src/data/x12Names837i.ts's ST_EL_NAME and SE_EL_NAME arrays
// (source: Transmittal 907, PDF p.8, 124).
//
// BillingProviders holds the repeating 2000A loop (BillingProviderHL), each
// nesting its own 2000B subscribers and, under them, 2000C patients and
// 2300 claims, placed by each HL's parent link (HL02).
type TransactionST struct {
	Index int `json:"index,omitempty"` // the ST's position in the file

	TransactionSetIDST01 TransactionTypeGS01_ST01 `json:"transactionSetIDST01,omitempty"` // Transaction Set Identifier Code — Transmittal 907, PDF p.8
	ControlNumberST02    string                   `json:"controlNumberST02,omitempty"`    // Transaction Set Control Number — Transmittal 907, PDF p.8
	VersionST03          string                   `json:"versionST03,omitempty"`          // Implementation Guide Version Name (005010X223A2) — Transmittal 907, PDF p.8

	Header HeaderBHT `json:"header,omitempty"` // BHT — Transmittal 907, PDF p.8-9

	SubmitterNM141 SubmitterNM141 `json:"submitterNM141,omitempty"` // 1000A — Transmittal 907, PDF p.9-10
	ReceiverNM140  ReceiverNM140  `json:"receiverNM140,omitempty"`  // 1000B — Transmittal 907, PDF p.10

	BillingProviders []BillingProviderHL `json:"billingProviders,omitempty"` // 2000A loops

	Raw []RawSegment `json:"raw,omitempty"` // segments directly under ST that the model does not place

	// Trailer: SE. The count is kept as the digits the file sent; the
	// parser checks it against the segments from ST through SE.
	IndexSE           int    `json:"indexSE,omitempty"`           // the SE's position in the file
	SegmentCountSE01  string `json:"segmentCountSE01,omitempty"`  // Transaction Segment Count — Transmittal 907, PDF p.124
	ControlNumberSE02 string `json:"controlNumberSE02,omitempty"` // Transaction Set Control Number — Transmittal 907, PDF p.124
}
