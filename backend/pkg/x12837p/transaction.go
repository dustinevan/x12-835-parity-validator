package x12837p

// TransactionST is one 837P transaction set: the ST header, the BHT header,
// the submitter and receiver loops, every billing provider hierarchy with
// its claims, and the matching SE trailer. Element names and citations for
// ST and SE are frontend/src/data/x12Names837p.ts's ST_EL_NAME and
// SE_EL_NAME arrays (source: Transmittal 907, PDF p.134 footer 2, PDF p.238
// footer 106).
//
// BillingProviders holds the 2000A Billing Provider loops
// (BillingProviderHL), each with its 2010AA-AC loops and the
// 2000B/2000C/2300 loops beneath it, in file order.
type TransactionST struct {
	Index int `json:"index,omitempty"` // the ST's position in the file, 0-based; the ISA is 0

	TransactionSetIDST01           TransactionSetIDST01 `json:"transactionSetIDST01,omitempty"`           // Transaction Set Identifier Code
	ControlNumberST02              string               `json:"controlNumberST02,omitempty"`              // Transaction Set Control Number
	ImplementationGuideVersionST03 string               `json:"implementationGuideVersionST03,omitempty"` // Implementation Guide Version Name

	Extra []string `json:"extra,omitempty"` // elements after ST03, kept as sent

	Header HeaderBHT `json:"header,omitempty"` // BHT — Transmittal 907, PDF p.134-135 (footer 2-3)

	SubmitterNM141 SubmitterNM141 `json:"submitterNM141,omitempty"` // 1000A — Transmittal 907, PDF p.135 (footer 3)
	ReceiverNM140  ReceiverNM140  `json:"receiverNM140,omitempty"`  // 1000B — Transmittal 907, PDF p.136 (footer 4)

	BillingProviders []BillingProviderHL `json:"billingProviders,omitempty"` // 2000A loops — Transmittal 907, PDF p.136 onward

	Raw []RawSegment `json:"raw,omitempty"` // segments directly under ST that the model does not place

	SegmentCountSE01  int    `json:"segmentCountSE01,omitempty"`  // Transaction Segment Count
	ControlNumberSE02 string `json:"controlNumberSE02,omitempty"` // Transaction Set Control Number
}

// TransactionSetIDST01 is ST01, the Transaction Set Identifier Code.
//
// The 837P implementation guide allows the one value, 837, visible in the
// CEDI companion guide's own worked 837P transaction (CEDI 837P CG, PDF p.50,
// printed p.42, Figure 3, "ST*837*..."; TRANSACTION_TYPE_NAME in
// frontend/src/data/x12Names837p.ts). The parser treats any other ST01 as a
// transaction that is not an 837.
type TransactionSetIDST01 string

const (
	TransactionSetIDST01NotSent TransactionSetIDST01 = ""
	TransactionSetIDST01837     TransactionSetIDST01 = "837"
)

var transactionSetIDST01List = codeList{
	{"837", "Health Care Claim", `Health Care Claim — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "ST*837*..."); frontend/src/data/x12Names837p.ts TRANSACTION_TYPE_NAME`},
}

func (v TransactionSetIDST01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v TransactionSetIDST01) Known() bool { return transactionSetIDST01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v TransactionSetIDST01) ShortDescription() string {
	return transactionSetIDST01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v TransactionSetIDST01) Description() string {
	return transactionSetIDST01List.description(string(v))
}

// ParseTransactionSetIDST01 returns text as a TransactionSetIDST01, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseTransactionSetIDST01(text string) (TransactionSetIDST01, error) {
	return TransactionSetIDST01(text), transactionSetIDST01List.check("ST01", text)
}
