package x12837i

// FunctionalGroupGS is one functional group: the GS header, every
// transaction set it carries, and the matching GE trailer. Element names
// and their citations are frontend/src/data/x12Names837i.ts's GS_EL_NAME and
// GE_EL_NAME arrays (source: Transmittal 907, PDF p.7-8, 124).
type FunctionalGroupGS struct {
	Index int `json:"index,omitempty"` // the GS's position in the file

	FunctionalIDGS01 TransactionTypeGS01_ST01 `json:"functionalIDGS01,omitempty"` // Functional Identifier Code — Transmittal 907, PDF p.7
	SenderCodeGS02   string                   `json:"senderCodeGS02,omitempty"`   // Application Sender's Code — Transmittal 907, PDF p.7
	ReceiverCodeGS03 string                   `json:"receiverCodeGS03,omitempty"` // Application Receiver's Code — Transmittal 907, PDF p.7

	DateGS04 Date `json:"dateGS04,omitempty"` // Date — Transmittal 907, PDF p.8
	TimeGS05 Time `json:"timeGS05,omitempty"` // Time — Transmittal 907, PDF p.8

	ControlNumberGS06 string `json:"controlNumberGS06,omitempty"` // Group Control Number — Transmittal 907, PDF p.8

	// AgencyGS07 (Responsible Agency Code) is kept as a plain string, not a
	// coded type - see InterchangeISA's doc comment for why.
	AgencyGS07  string `json:"agencyGS07,omitempty"`  // Responsible Agency Code — Transmittal 907, PDF p.8
	VersionGS08 string `json:"versionGS08,omitempty"` // Version/Release/Industry Identifier Code — Transmittal 907, PDF p.8

	Transactions []TransactionST `json:"transactions,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between GS and GE, outside any ST, that the model does not place

	// Trailer: GE. Envelope counts, not business data; the count is kept as
	// the digits the file sent.
	IndexGE              int    `json:"indexGE,omitempty"`              // the GE's position in the file
	TransactionCountGE01 string `json:"transactionCountGE01,omitempty"` // Number of Transaction Sets Included — Transmittal 907, PDF p.124
	ControlNumberGE02    string `json:"controlNumberGE02,omitempty"`    // Group Control Number — Transmittal 907, PDF p.124
}
