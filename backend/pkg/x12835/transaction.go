package x12835

// TransactionST is one 835 transaction set: the ST header, the payment
// header, the payer and payee loops, every header-number loop with its
// claims, any provider-level adjustments, and the matching SE trailer.
// Element names and citations for ST and SE are
// frontend/src/data/x12Names835.ts's ST_EL_NAME and SE_EL_NAME arrays
// (source: CMS 835 Flat File, PDF p.1, 14).
type TransactionST struct {
	Index int `json:"index,omitempty"` // the ST segment's position in the file

	TransactionSetIDST01 TransactionSetIDST01 `json:"transactionSetIDST01,omitempty"` // Transaction Set Identifier Code — CMS 835 Flat File, PDF p.1
	ControlNumberST02    string               `json:"controlNumberST02,omitempty"`    // Transaction Set Control Number — CMS 835 Flat File, PDF p.1

	// ImplementationConventionST03 is marked N/U on the flat file page, but
	// files send it (005010X221A1), so the model keeps it.
	ImplementationConventionST03 string `json:"implementationConventionST03,omitempty"` // Implementation Convention Reference — CMS 835 Flat File, PDF p.1

	Payment  PaymentBPR  `json:"payment,omitempty"`  // BPR — CMS 835 Flat File, PDF p.1-2
	Trace    TraceTRN    `json:"trace,omitempty"`    // TRN — CMS 835 Flat File, PDF p.2
	Currency CurrencyCUR `json:"currency,omitempty"` // CUR, rarely used — CMS 835 Flat File, PDF p.2

	// HeaderRefs are the REF segments that sit directly under ST, before any
	// loop (qualifiers EV and F2; HEADER_REF_NAME in
	// frontend/src/data/x12Names835.ts).
	HeaderRefs []ReferenceREF `json:"headerRefs,omitempty"`

	// ProductionDateDTM405 is the one DTM the header carries (qualifier 405;
	// DTM_NAME in frontend/src/data/x12Names835.ts). It is present when its
	// DateQualifierDTM01 is not empty.
	ProductionDateDTM405 DateTimeDTM `json:"productionDateDTM405,omitempty"`

	PayerN1 PartyN1 `json:"payerN1,omitempty"` // N1*PR loop (1000A) — CMS 835 Flat File, PDF p.2-3
	PayeeN1 PartyN1 `json:"payeeN1,omitempty"` // N1*PE loop (1000B) — CMS 835 Flat File, PDF p.3-4

	HeaderNumbers []HeaderNumberLX `json:"headerNumbers,omitempty"` // LX loops — CMS 835 Flat File, PDF p.4 onward

	ProviderAdjustments []ProviderAdjustmentPLB `json:"providerAdjustments,omitempty"` // PLB, after every LX loop — CMS 835 Flat File, PDF p.13-14

	Raw []RawSegment `json:"raw,omitempty"` // segments directly under ST that the model does not place, including every segment of an N1 loop that is neither the first PR nor the first PE loop

	SegmentCountSE01  int    `json:"segmentCountSE01,omitempty"`  // Transaction Segment Count — CMS 835 Flat File, PDF p.14
	ControlNumberSE02 string `json:"controlNumberSE02,omitempty"` // Transaction Set Control # — CMS 835 Flat File, PDF p.14
}
