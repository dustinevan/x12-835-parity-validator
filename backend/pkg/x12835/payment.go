package x12835

// PaymentBPR is the BPR segment: the payment header carried by every 835
// transaction. Element names and citations are
// frontend/src/data/x12Names835.ts's BPR_EL_NAME array (source: CMS 835
// Flat File, PDF p.1-2). Set is false when the segment did not appear at
// all. BPR17 to BPR21 are the positions the standard defines after the
// effective date; the 835 guide marks them not used, and a file that sends
// them keeps them.
type PaymentBPR struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	HandlingBPR01    HandlingBPR01    `json:"handlingBPR01,omitempty"`    // Transaction Handling Code — CMS 835 Flat File, PDF p.1
	TotalPaidBPR02   Amount           `json:"totalPaidBPR02,omitempty"`   // Total Actual Provider Payment Amt — CMS 835 Flat File, PDF p.1
	CreditDebitBPR03 CreditDebitBPR03 `json:"creditDebitBPR03,omitempty"` // Credit or Debit Flag Code — CMS 835 Flat File, PDF p.1
	PayMethodBPR04   PayMethodBPR04   `json:"payMethodBPR04,omitempty"`   // Payment Method Code — CMS 835 Flat File, PDF p.1
	PayFormatBPR05   PayFormatBPR05   `json:"payFormatBPR05,omitempty"`   // Payment Format Code — CMS 835 Flat File, PDF p.1

	SenderDFIQualifierBPR06  DFIQualifierBPR06     `json:"senderDFIQualifierBPR06,omitempty"`  // DFI ID # Qualifier — CMS 835 Flat File, PDF p.1
	SenderDFIBPR07           string                `json:"senderDFIBPR07,omitempty"`           // Sender DFI Identifier — CMS 835 Flat File, PDF p.1
	SenderAcctQualifierBPR08 AccountQualifierBPR08 `json:"senderAcctQualifierBPR08,omitempty"` // Acct # Qualifier — CMS 835 Flat File, PDF p.1
	SenderAcctBPR09          string                `json:"senderAcctBPR09,omitempty"`          // Sender Bank Acct # — CMS 835 Flat File, PDF p.1

	PayerIDBPR10       string `json:"payerIDBPR10,omitempty"`       // Payer Identifier — CMS 835 Flat File, PDF p.1
	OrigCoSupCodeBPR11 string `json:"origCoSupCodeBPR11,omitempty"` // Originating Co Supplemental Code — CMS 835 Flat File, PDF p.1

	ReceiverDFIQualifierBPR12  DFIQualifierBPR06     `json:"receiverDFIQualifierBPR12,omitempty"`  // DFI ID # Qualifier — CMS 835 Flat File, PDF p.2
	ReceiverDFIBPR13           string                `json:"receiverDFIBPR13,omitempty"`           // Receiver or Provider Bank ID # — CMS 835 Flat File, PDF p.2
	ReceiverAcctQualifierBPR14 AccountQualifierBPR08 `json:"receiverAcctQualifierBPR14,omitempty"` // Acct # Qualifier — CMS 835 Flat File, PDF p.2
	ReceiverAcctBPR15          string                `json:"receiverAcctBPR15,omitempty"`          // Receiver or Provider Acct # — CMS 835 Flat File, PDF p.2

	EffectiveDateBPR16 Date `json:"effectiveDateBPR16,omitempty"` // Check Issue or EFT Effective Date — CMS 835 Flat File, PDF p.2

	BusinessFunctionBPR17   string `json:"businessFunctionBPR17,omitempty"`   // BPR17 Business Function Code; not used by the 835 guide, kept as sent
	ThirdDFIQualifierBPR18  string `json:"thirdDFIQualifierBPR18,omitempty"`  // BPR18 (DFI) ID Number Qualifier; not used by the 835 guide, kept as sent
	ThirdDFIBPR19           string `json:"thirdDFIBPR19,omitempty"`           // BPR19 (DFI) Identification Number; not used by the 835 guide, kept as sent
	ThirdAcctQualifierBPR20 string `json:"thirdAcctQualifierBPR20,omitempty"` // BPR20 Account Number Qualifier; not used by the 835 guide, kept as sent
	ThirdAcctBPR21          string `json:"thirdAcctBPR21,omitempty"`          // BPR21 Account Number; not used by the 835 guide, kept as sent
}

// TraceTRN is the TRN segment: the reassociation trace number. Element
// names and citations are frontend/src/data/x12Names835.ts's TRN_EL_NAME
// array (source: CMS 835 Flat File, PDF p.2). Set is false when the
// segment did not appear at all.
type TraceTRN struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	TraceTypeTRN01     TraceTypeTRN01 `json:"traceTypeTRN01,omitempty"`     // Trace Type Code — CMS 835 Flat File, PDF p.2
	TraceNumberTRN02   string         `json:"traceNumberTRN02,omitempty"`   // Check or EFT Trace # — CMS 835 Flat File, PDF p.2
	PayerIDTRN03       string         `json:"payerIDTRN03,omitempty"`       // Payer Identifier — CMS 835 Flat File, PDF p.2
	OrigCoSupCodeTRN04 string         `json:"origCoSupCodeTRN04,omitempty"` // Originating Co Supplemental Code — CMS 835 Flat File, PDF p.2
}

// CurrencyCUR is the CUR segment: Foreign Currency Information. The flat
// file's own comment says Medicare does not use it (SEGMENT_NAME in
// frontend/src/data/x12Names835.ts, citing CMS 835 Flat File, PDF p.2), and
// that file gives no per-element names for it - unlike every other segment
// here, its field comments cite the element reference only. The standard
// defines 21 positions; the 835 guide uses the first three, and the rest are
// kept as sent. Set is false when the segment did not appear at all.
type CurrencyCUR struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	EntityIDCUR01     string `json:"entityIDCUR01,omitempty"`     // CUR01, Entity Identifier Code - element reference only, no public name in x12Names835.ts.
	CurrencyCodeCUR02 string `json:"currencyCodeCUR02,omitempty"` // CUR02, Currency Code - element reference only.
	ExchangeRateCUR03 Amount `json:"exchangeRateCUR03,omitempty"` // CUR03, Exchange Rate - element reference only.

	// CUR04 to CUR21, kept as sent: a second entity and currency, a market
	// code, and six date/time triples (qualifier, date, time).
	EntityIDCUR04      string `json:"entityIDCUR04,omitempty"`
	CurrencyCodeCUR05  string `json:"currencyCodeCUR05,omitempty"`
	MarketCUR06        string `json:"marketCUR06,omitempty"`
	DateQualifierCUR07 string `json:"dateQualifierCUR07,omitempty"`
	DateCUR08          string `json:"dateCUR08,omitempty"`
	TimeCUR09          string `json:"timeCUR09,omitempty"`
	DateQualifierCUR10 string `json:"dateQualifierCUR10,omitempty"`
	DateCUR11          string `json:"dateCUR11,omitempty"`
	TimeCUR12          string `json:"timeCUR12,omitempty"`
	DateQualifierCUR13 string `json:"dateQualifierCUR13,omitempty"`
	DateCUR14          string `json:"dateCUR14,omitempty"`
	TimeCUR15          string `json:"timeCUR15,omitempty"`
	DateQualifierCUR16 string `json:"dateQualifierCUR16,omitempty"`
	DateCUR17          string `json:"dateCUR17,omitempty"`
	TimeCUR18          string `json:"timeCUR18,omitempty"`
	DateQualifierCUR19 string `json:"dateQualifierCUR19,omitempty"`
	DateCUR20          string `json:"dateCUR20,omitempty"`
	TimeCUR21          string `json:"timeCUR21,omitempty"`
}
