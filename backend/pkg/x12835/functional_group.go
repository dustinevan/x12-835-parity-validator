package x12835

// FunctionalGroupGS is one functional group: the GS header, every
// transaction set it carries, and the matching GE trailer. Element names
// and their citations are frontend/src/data/x12Names835.ts's GS_EL_NAME and
// GE_EL_NAME arrays (source: CMS 835 Flat File, PDF p.1-2, 14).
type FunctionalGroupGS struct {
	Index int `json:"index,omitempty"` // the GS segment's position in the file

	FunctionalIDGS01 FunctionalIDGS01 `json:"functionalIDGS01,omitempty"` // Functional Identifier Code — CMS 835 Flat File, PDF p.1
	SenderCodeGS02   string           `json:"senderCodeGS02,omitempty"`   // Application Sender's Code — CMS 835 Flat File, PDF p.1
	ReceiverCodeGS03 string           `json:"receiverCodeGS03,omitempty"` // Application Receiver's Code — CMS 835 Flat File, PDF p.1

	DateGS04 Date `json:"dateGS04,omitempty"` // Date — CMS 835 Flat File, PDF p.2
	TimeGS05 Time `json:"timeGS05,omitempty"` // Time — CMS 835 Flat File, PDF p.2

	ControlNumberGS06 string     `json:"controlNumberGS06,omitempty"` // Group Control Number — CMS 835 Flat File, PDF p.2
	AgencyGS07        AgencyGS07 `json:"agencyGS07,omitempty"`        // Responsible Agency Code — CMS 835 Flat File, PDF p.2
	VersionGS08       string     `json:"versionGS08,omitempty"`       // Version/Release/Industry ID code — CMS 835 Flat File, PDF p.2

	Transactions []TransactionST `json:"transactions,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between GS and GE, outside any ST, that the model does not place

	// Trailer: GE. Envelope counts, not business data; the parser checks
	// TransactionCountGE01 against the transactions actually present and
	// ControlNumberGE02 against ControlNumberGS06.
	TransactionCountGE01 int    `json:"transactionCountGE01,omitempty"` // # Transaction Sets Included — CMS 835 Flat File, PDF p.14
	ControlNumberGE02    string `json:"controlNumberGE02,omitempty"`    // Group Control # — CMS 835 Flat File, PDF p.14
}
