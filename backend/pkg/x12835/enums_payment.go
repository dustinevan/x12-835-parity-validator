package x12835

// HandlingBPR01 is BPR01, the Transaction Handling Code.
//
// The flat file lists all seven of these as valid BPR01 values; neither
// reviewed document names what any of them means in text. A real Medicare
// remittance that carries claim payments sends C or I; D/H/P/U/X describe
// non-monetary handling (docs/design/835-examples.md's "Items not built"
// note) and are included here because the flat file lists them as valid, not
// because an example produces them.
type HandlingBPR01 string

const (
	HandlingBPR01NotSent HandlingBPR01 = ""
	HandlingBPR01C       HandlingBPR01 = "C"
	HandlingBPR01D       HandlingBPR01 = "D"
	HandlingBPR01H       HandlingBPR01 = "H"
	HandlingBPR01I       HandlingBPR01 = "I"
	HandlingBPR01P       HandlingBPR01 = "P"
	HandlingBPR01U       HandlingBPR01 = "U"
	HandlingBPR01X       HandlingBPR01 = "X"
)

// handlingBPR01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var handlingBPR01Codes = []codeEntry{
	{"C", "Payment Accompanies Remittance Advice", "Payment Accompanies Remittance Advice — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text (docs/design/835-examples.md Open items)"},
	{"D", "Make Payment Only", "Make Payment Only — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
	{"H", "Notification Only", "Notification Only — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
	{"I", "Remittance Information Only", "Remittance Information Only — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text; docs/design/835-examples.md cross-references X12 public Example 01 (\"Dollars and Data Sent Separately\"), which pairs this code with an ACH payment"},
	{"P", "Prenotification of Future Transfers", "Prenotification of Future Transfers — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
	{"U", "Split Payment and Remittance", "Split Payment and Remittance — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
	{"X", "Handling Party's Option to Split Payment and Remittance", "Handling Party's Option to Split Payment and Remittance — X12 005010X221A1 TR3, BPR01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
}

func (v HandlingBPR01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v HandlingBPR01) Known() bool { return known(handlingBPR01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v HandlingBPR01) ShortDescription() string { return shortOf(handlingBPR01Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v HandlingBPR01) Description() string { return descriptionOf(handlingBPR01Codes, string(v)) }

// ParseHandlingBPR01 reads one BPR01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseHandlingBPR01(code string) (HandlingBPR01, error) {
	return HandlingBPR01(code), checkListed("BPR01", handlingBPR01Codes, code)
}

// CreditDebitBPR03 is BPR03, the Credit or Debit Flag Code.
//
// CEDI 835 CG, PDF p.24 (printed p.18) states only C applies to Medicare
// ("D" is a valid HIPAA code but Medicare never sends it); D is listed
// because the standard allows it.
type CreditDebitBPR03 string

const (
	CreditDebitBPR03NotSent CreditDebitBPR03 = ""
	CreditDebitBPR03C       CreditDebitBPR03 = "C"
	CreditDebitBPR03D       CreditDebitBPR03 = "D"
)

// creditDebitBPR03Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var creditDebitBPR03Codes = []codeEntry{
	{"C", "Credit", "Credit — CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.24 (printed p.18; Table 6)"},
	{"D", "Debit", "Debit — X12 005010X221A1 TR3, BPR03 (transcribed from the standard; verify against the printed guide)"},
}

func (v CreditDebitBPR03) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v CreditDebitBPR03) Known() bool { return known(creditDebitBPR03Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v CreditDebitBPR03) ShortDescription() string { return shortOf(creditDebitBPR03Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v CreditDebitBPR03) Description() string {
	return descriptionOf(creditDebitBPR03Codes, string(v))
}

// ParseCreditDebitBPR03 reads one BPR03 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseCreditDebitBPR03(code string) (CreditDebitBPR03, error) {
	return CreditDebitBPR03(code), checkListed("BPR03", creditDebitBPR03Codes, code)
}

// PayMethodBPR04 is BPR04, the Payment Method Code.
type PayMethodBPR04 string

const (
	PayMethodBPR04NotSent PayMethodBPR04 = ""
	PayMethodBPR04ACH     PayMethodBPR04 = "ACH"
	PayMethodBPR04CHK     PayMethodBPR04 = "CHK"
	PayMethodBPR04NON     PayMethodBPR04 = "NON"
	PayMethodBPR04BOP     PayMethodBPR04 = "BOP"
	PayMethodBPR04FWT     PayMethodBPR04 = "FWT"
)

// payMethodBPR04Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var payMethodBPR04Codes = []codeEntry{
	{"ACH", "Automated Clearing House (ACH)", "Automated Clearing House (ACH) — X12 005010X221A1 TR3, BPR04 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.24 (printed p.18); not named in text in either document, but the abbreviation (Automated Clearing House) is unambiguous"},
	{"CHK", "Check", "Check — X12 005010X221A1 TR3, BPR04 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.24 (printed p.18); not named in text; abbreviation (check) unambiguous"},
	{"NON", "Non-Payment Data", "Non-Payment Data — X12 005010X221A1 TR3, BPR04 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.24 (printed p.18); not named in text (X12's own abbreviation is \"Non-Payment Data,\" used for $0 remittances)"},
	{"BOP", "Financial Institution Option", "Financial Institution Option — X12 005010X221A1 TR3, BPR04 (transcribed from the standard; verify against the printed guide)"},
	{"FWT", "Federal Reserve Funds/Wire Transfer - Nonrepetitive", "Federal Reserve Funds/Wire Transfer - Nonrepetitive — X12 005010X221A1 TR3, BPR04 (transcribed from the standard; verify against the printed guide)"},
}

func (v PayMethodBPR04) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v PayMethodBPR04) Known() bool { return known(payMethodBPR04Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v PayMethodBPR04) ShortDescription() string { return shortOf(payMethodBPR04Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v PayMethodBPR04) Description() string { return descriptionOf(payMethodBPR04Codes, string(v)) }

// ParsePayMethodBPR04 reads one BPR04 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParsePayMethodBPR04(code string) (PayMethodBPR04, error) {
	return PayMethodBPR04(code), checkListed("BPR04", payMethodBPR04Codes, code)
}

// PayFormatBPR05 is BPR05, the Payment Format Code.
//
// Populated only when BPR04 is ACH (CMS 835 Flat File, PDF p.1).
type PayFormatBPR05 string

const (
	PayFormatBPR05NotSent PayFormatBPR05 = ""
	PayFormatBPR05CCP     PayFormatBPR05 = "CCP"
	PayFormatBPR05CTX     PayFormatBPR05 = "CTX"
)

// payFormatBPR05Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var payFormatBPR05Codes = []codeEntry{
	{"CCP", "Cash Concentration/Disbursement plus Addenda (CCD+) (ACH)", "Cash Concentration/Disbursement plus Addenda (CCD+) (ACH) — X12 005010X221A1 TR3, BPR05 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
	{"CTX", "Corporate Trade Exchange (CTX) (ACH)", "Corporate Trade Exchange (CTX) (ACH) — X12 005010X221A1 TR3, BPR05 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text"},
}

func (v PayFormatBPR05) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v PayFormatBPR05) Known() bool { return known(payFormatBPR05Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v PayFormatBPR05) ShortDescription() string { return shortOf(payFormatBPR05Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v PayFormatBPR05) Description() string { return descriptionOf(payFormatBPR05Codes, string(v)) }

// ParsePayFormatBPR05 reads one BPR05 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParsePayFormatBPR05(code string) (PayFormatBPR05, error) {
	return PayFormatBPR05(code), checkListed("BPR05", payFormatBPR05Codes, code)
}

// DFIQualifierBPR06 is the DFI ID Number Qualifier, used for both BPR06
// (sender DFI, populated when BPR04 is ACH) and BPR12 (receiver/provider
// DFI).
//
// Neither reviewed document names this element's individual codes.
type DFIQualifierBPR06 string

const (
	DFIQualifierBPR06NotSent DFIQualifierBPR06 = ""
	DFIQualifierBPR0601      DFIQualifierBPR06 = "01"
	DFIQualifierBPR0604      DFIQualifierBPR06 = "04"
)

// dfiQualifierBPR06Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var dfiQualifierBPR06Codes = []codeEntry{
	{"01", "ABA Transit Routing Number Including Check Digits (9 digits)", "ABA Transit Routing Number Including Check Digits (9 digits) — X12 005010X221A1 TR3, BPR06 and BPR12 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1-2; not named in text; the only value seen in this repo's example 835 files"},
	{"04", "Canadian Bank Branch and Institution Number", "Canadian Bank Branch and Institution Number — X12 005010X221A1 TR3, BPR06 and BPR12 (transcribed from the standard; verify against the printed guide)"},
}

func (v DFIQualifierBPR06) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v DFIQualifierBPR06) Known() bool { return known(dfiQualifierBPR06Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v DFIQualifierBPR06) ShortDescription() string {
	return shortOf(dfiQualifierBPR06Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v DFIQualifierBPR06) Description() string {
	return descriptionOf(dfiQualifierBPR06Codes, string(v))
}

// ParseDFIQualifierBPR06 reads one BPR06 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseDFIQualifierBPR06(code string) (DFIQualifierBPR06, error) {
	return DFIQualifierBPR06(code), checkListed("BPR06", dfiQualifierBPR06Codes, code)
}

// AccountQualifierBPR08 is the Account Number Qualifier, used for both BPR08
// (sender account) and BPR14 (receiver/provider account).
//
// Neither reviewed document names this element's individual codes.
type AccountQualifierBPR08 string

const (
	AccountQualifierBPR08NotSent AccountQualifierBPR08 = ""
	AccountQualifierBPR08DA      AccountQualifierBPR08 = "DA"
	AccountQualifierBPR08SG      AccountQualifierBPR08 = "SG"
)

// accountQualifierBPR08Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var accountQualifierBPR08Codes = []codeEntry{
	{"DA", "Demand Deposit", "Demand Deposit — X12 005010X221A1 TR3, BPR08 and BPR14 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1-2; not named in text; the only value seen in this repo's example 835 files"},
	{"SG", "Savings", "Savings — X12 005010X221A1 TR3, BPR08 and BPR14 (transcribed from the standard; verify against the printed guide); the guide allows SG in BPR14 (the receiver's account) only; BPR08 takes DA alone"},
}

func (v AccountQualifierBPR08) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AccountQualifierBPR08) Known() bool { return known(accountQualifierBPR08Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AccountQualifierBPR08) ShortDescription() string {
	return shortOf(accountQualifierBPR08Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AccountQualifierBPR08) Description() string {
	return descriptionOf(accountQualifierBPR08Codes, string(v))
}

// ParseAccountQualifierBPR08 reads one BPR08 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAccountQualifierBPR08(code string) (AccountQualifierBPR08, error) {
	return AccountQualifierBPR08(code), checkListed("BPR08", accountQualifierBPR08Codes, code)
}

// TraceTypeTRN01 is TRN01, the Trace Type Code.
type TraceTypeTRN01 string

const (
	TraceTypeTRN01NotSent TraceTypeTRN01 = ""
	TraceTypeTRN011       TraceTypeTRN01 = "1"
)

// traceTypeTRN01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var traceTypeTRN01Codes = []codeEntry{
	{"1", "Reassociation Trace Number", "Reassociation Trace Number — CMS 835 Flat File, PDF p.2; frontend/src/data/x12Names835.ts TRN_NAME"},
}

func (v TraceTypeTRN01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v TraceTypeTRN01) Known() bool { return known(traceTypeTRN01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v TraceTypeTRN01) ShortDescription() string { return shortOf(traceTypeTRN01Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v TraceTypeTRN01) Description() string { return descriptionOf(traceTypeTRN01Codes, string(v)) }

// ParseTraceTypeTRN01 reads one TRN01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseTraceTypeTRN01(code string) (TraceTypeTRN01, error) {
	return TraceTypeTRN01(code), checkListed("TRN01", traceTypeTRN01Codes, code)
}
