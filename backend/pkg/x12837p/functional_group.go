package x12837p

// FunctionalGroupGS is one functional group: the GS header, every
// transaction set it carries, and the matching GE trailer. Element names
// and citations are frontend/src/data/x12Names837p.ts's GS_EL_NAME and
// GE_EL_NAME arrays (source: Transmittal 907, PDF p.133-134 footer 1-2, PDF
// p.238 footer 106).
type FunctionalGroupGS struct {
	Index int `json:"index,omitempty"` // the GS's position in the file, 0-based; the ISA is 0

	FunctionalIDGS01 FunctionalIDGS01 `json:"functionalIDGS01,omitempty"` // Functional Identifier Code
	SenderCodeGS02   string           `json:"senderCodeGS02,omitempty"`   // Application Sender Code
	ReceiverCodeGS03 string           `json:"receiverCodeGS03,omitempty"` // Application Receiver Code

	DateGS04 Date `json:"dateGS04,omitempty"` // Date
	TimeGS05 Time `json:"timeGS05,omitempty"` // Time

	ControlNumberGS06 string     `json:"controlNumberGS06,omitempty"` // Group Control Number
	AgencyGS07        AgencyGS07 `json:"agencyGS07,omitempty"`        // Responsible Agency Code
	VersionGS08       string     `json:"versionGS08,omitempty"`       // Version/Release/Industry Identifier Code

	Extra []string `json:"extra,omitempty"` // elements after GS08, kept as sent

	Transactions []TransactionST `json:"transactions,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between GS and GE, outside any ST, that the model does not place

	// Trailer: GE. Envelope counts, not business data; the parser checks
	// TransactionCountGE01 against the transactions actually present and
	// ControlNumberGE02 against ControlNumberGS06.
	TransactionCountGE01 int    `json:"transactionCountGE01,omitempty"` // Number of Transaction Sets Included
	ControlNumberGE02    string `json:"controlNumberGE02,omitempty"`    // Group Control Number
}

// FunctionalIDGS01 is GS01, the Functional Identifier Code.
//
// The 837P implementation guide allows the one value, HC (Health Care
// Claim), visible in the CEDI companion guide's own worked 837P transaction
// (TRANSACTION_TYPE_NAME in frontend/src/data/x12Names837p.ts, citing
// Transmittal 907's own "GS*HC*..." and CEDI 837P CG, PDF p.50, printed p.42,
// Figure 3). The parser treats any other GS01 as a file that is not an 837P.
type FunctionalIDGS01 string

const (
	FunctionalIDGS01NotSent FunctionalIDGS01 = ""
	FunctionalIDGS01HC      FunctionalIDGS01 = "HC"
)

var functionalIDGS01List = codeList{
	{"HC", "Health Care Claim", `Health Care Claim — Transmittal 907 (GS01 value shown in the CEDI worked example, "GS*HC*..."); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); frontend/src/data/x12Names837p.ts TRANSACTION_TYPE_NAME`},
}

func (v FunctionalIDGS01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v FunctionalIDGS01) Known() bool { return functionalIDGS01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v FunctionalIDGS01) ShortDescription() string { return functionalIDGS01List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v FunctionalIDGS01) Description() string { return functionalIDGS01List.description(string(v)) }

// ParseFunctionalIDGS01 returns text as a FunctionalIDGS01, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseFunctionalIDGS01(text string) (FunctionalIDGS01, error) {
	return FunctionalIDGS01(text), functionalIDGS01List.check("GS01", text)
}

// AgencyGS07 is GS07, the Responsible Agency Code.
//
// The 837P implementation guide allows the one code, X. Transmittal 579's
// Values column lists it (PDF p.7), and the CEDI worked transaction sends it;
// the name is the guide's.
type AgencyGS07 string

const (
	AgencyGS07NotSent AgencyGS07 = ""
	AgencyGS07X       AgencyGS07 = "X"
)

var agencyGS07List = codeList{
	{"X", "Accredited Standards Committee X12", "Accredited Standards Committee X12 — Transmittal 579, PDF p.7 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
}

func (v AgencyGS07) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AgencyGS07) Known() bool { return agencyGS07List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AgencyGS07) ShortDescription() string { return agencyGS07List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AgencyGS07) Description() string { return agencyGS07List.description(string(v)) }

// ParseAgencyGS07 returns text as an AgencyGS07, with an error when the text
// is neither empty nor on the list. The value is kept either way.
func ParseAgencyGS07(text string) (AgencyGS07, error) {
	return AgencyGS07(text), agencyGS07List.check("GS07", text)
}
