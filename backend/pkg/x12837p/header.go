package x12837p

// HeaderBHT is the BHT segment: the beginning of the hierarchical
// transaction, carried directly under ST before any loop. Element names and
// citations are frontend/src/data/x12Names837p.ts's BHT_EL_NAME array
// (source: Transmittal 907, PDF p.134-135, footer 2-3). Set is false when
// the segment did not appear at all.
type HeaderBHT struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	StructureBHT01 HierarchicalStructureBHT01 `json:"structureBHT01,omitempty"` // Hierarchical Structure Code
	PurposeBHT02   TransactionSetPurposeBHT02 `json:"purposeBHT02,omitempty"`   // Transaction Set Purpose Code

	OriginatorTransactionIDBHT03 string `json:"originatorTransactionIDBHT03,omitempty"` // Originator Application Transaction ID

	DateBHT04 Date `json:"dateBHT04,omitempty"` // Transaction Set Creation Date
	TimeBHT05 Time `json:"timeBHT05,omitempty"` // Transaction Set Creation Time

	TransactionTypeBHT06 TransactionTypeBHT06 `json:"transactionTypeBHT06,omitempty"` // Claim or Encounter ID (X12's own BHT06 Transaction Type Code)

	Extra []string `json:"extra,omitempty"` // elements after BHT06, kept as sent
}

// HierarchicalStructureBHT01 is BHT01, the Hierarchical Structure Code.
//
// The 837P implementation guide allows the one code, 0019. Transmittal 579's
// Values column lists it (PDF p.7), also shown in the CEDI worked
// transaction; the name is the guide's.
type HierarchicalStructureBHT01 string

const (
	HierarchicalStructureBHT01NotSent HierarchicalStructureBHT01 = ""
	HierarchicalStructureBHT010019    HierarchicalStructureBHT01 = "0019"
)

var hierarchicalStructureBHT01List = codeList{
	{"0019", "Information Source, Subscriber, Dependent", `Information Source, Subscriber, Dependent — Transmittal 579, PDF p.7 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "BHT*0019*00*..."); the name is the guide's: ` + tr3},
}

func (v HierarchicalStructureBHT01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v HierarchicalStructureBHT01) Known() bool {
	return hierarchicalStructureBHT01List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v HierarchicalStructureBHT01) ShortDescription() string {
	return hierarchicalStructureBHT01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v HierarchicalStructureBHT01) Description() string {
	return hierarchicalStructureBHT01List.description(string(v))
}

// ParseHierarchicalStructureBHT01 returns text as a
// HierarchicalStructureBHT01, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseHierarchicalStructureBHT01(text string) (HierarchicalStructureBHT01, error) {
	return HierarchicalStructureBHT01(text), hierarchicalStructureBHT01List.check("BHT01", text)
}

// TransactionSetPurposeBHT02 is BHT02, the Transaction Set Purpose Code.
//
// The 837P implementation guide allows 00 and 18. Medicare restricts
// submissions to BHT02 = 00 (Original). CEDI 837P CG, PDF p.30 (printed p.22;
// Table 7); BHT02_NAME in frontend/src/data/x12Names837p.ts.
type TransactionSetPurposeBHT02 string

const (
	TransactionSetPurposeBHT02NotSent TransactionSetPurposeBHT02 = ""
	TransactionSetPurposeBHT0200      TransactionSetPurposeBHT02 = "00"
	TransactionSetPurposeBHT0218      TransactionSetPurposeBHT02 = "18"
)

var transactionSetPurposeBHT02List = codeList{
	{"00", "Original", "Original — CEDI 837P CG, PDF p.30 (printed p.22; Table 7); the only value Medicare submissions send"},
	{"18", "Reissue", guide("Reissue", "BHT02") + "; not a value Medicare submissions send (CEDI 837P CG, PDF p.30)"},
}

func (v TransactionSetPurposeBHT02) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v TransactionSetPurposeBHT02) Known() bool {
	return transactionSetPurposeBHT02List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v TransactionSetPurposeBHT02) ShortDescription() string {
	return transactionSetPurposeBHT02List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v TransactionSetPurposeBHT02) Description() string {
	return transactionSetPurposeBHT02List.description(string(v))
}

// ParseTransactionSetPurposeBHT02 returns text as a
// TransactionSetPurposeBHT02, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseTransactionSetPurposeBHT02(text string) (TransactionSetPurposeBHT02, error) {
	return TransactionSetPurposeBHT02(text), transactionSetPurposeBHT02List.check("BHT02", text)
}

// TransactionTypeBHT06 is BHT06. Transmittal 907 labels the position
// "Claim or Encounter ID" (BHT_EL_NAME); it is X12's own BHT06 Transaction
// Type Code element underneath that label. The 837P implementation guide
// allows 31, CH, and RP. Medicare restricts submissions to BHT06 = CH
// (Chargeable). CEDI 837P CG, PDF p.30 (printed p.22; Table 7); BHT06_NAME in
// frontend/src/data/x12Names837p.ts.
type TransactionTypeBHT06 string

const (
	TransactionTypeBHT06NotSent TransactionTypeBHT06 = ""
	TransactionTypeBHT06CH      TransactionTypeBHT06 = "CH"
	TransactionTypeBHT06RP      TransactionTypeBHT06 = "RP"
	TransactionTypeBHT0631      TransactionTypeBHT06 = "31"
)

var transactionTypeBHT06List = codeList{
	{"CH", "Chargeable", "Chargeable — CEDI 837P CG, PDF p.30 (printed p.22; Table 7); the only value Medicare submissions send"},
	{"RP", "Reporting", guide("Reporting", "BHT06") + "; not a value Medicare submissions send (CEDI 837P CG, PDF p.30)"},
	{"31", "Subrogation Demand", guide("Subrogation Demand", "BHT06") + "; not a value Medicare submissions send (CEDI 837P CG, PDF p.30)"},
}

func (v TransactionTypeBHT06) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v TransactionTypeBHT06) Known() bool { return transactionTypeBHT06List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v TransactionTypeBHT06) ShortDescription() string {
	return transactionTypeBHT06List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v TransactionTypeBHT06) Description() string {
	return transactionTypeBHT06List.description(string(v))
}

// ParseTransactionTypeBHT06 returns text as a TransactionTypeBHT06, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseTransactionTypeBHT06(text string) (TransactionTypeBHT06, error) {
	return TransactionTypeBHT06(text), transactionTypeBHT06List.check("BHT06", text)
}
