package x12835

// AuthQualifierISA01 is ISA01, the Authorization Information Qualifier.
//
// Neither CMS 835 Flat File nor the CEDI companion guide names ISA01's
// individual codes in text (frontend/src/data/x12Names835.ts's ISA_EL_NAME
// only gives the element's own position name, "Authorization Information
// Qualifier"); the names below are transcribed from the X12 standard.
type AuthQualifierISA01 string

const (
	AuthQualifierISA01NotSent AuthQualifierISA01 = ""
	AuthQualifierISA0100      AuthQualifierISA01 = "00"
	AuthQualifierISA0103      AuthQualifierISA01 = "03"
)

// authQualifierISA01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var authQualifierISA01Codes = []codeEntry{
	{"00", "No Authorization Information Present (No Meaningful Information in I02)", "No Authorization Information Present (No Meaningful Information in I02) — X12 005010X221A1 TR3, ISA01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text in either reviewed document; the only value seen in this repo's example 835 files"},
	{"03", "Additional Data Identification", "Additional Data Identification — X12 005010X221A1 TR3, ISA01 (transcribed from the standard; verify against the printed guide)"},
}

func (v AuthQualifierISA01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AuthQualifierISA01) Known() bool { return known(authQualifierISA01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AuthQualifierISA01) ShortDescription() string {
	return shortOf(authQualifierISA01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AuthQualifierISA01) Description() string {
	return descriptionOf(authQualifierISA01Codes, string(v))
}

// ParseAuthQualifierISA01 reads one ISA01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAuthQualifierISA01(code string) (AuthQualifierISA01, error) {
	return AuthQualifierISA01(code), checkListed("ISA01", authQualifierISA01Codes, code)
}

// SecQualifierISA03 is ISA03, the Security Information Qualifier.
//
// Same sourcing note as AuthQualifierISA01: the position is named, not its
// individual codes.
type SecQualifierISA03 string

const (
	SecQualifierISA03NotSent SecQualifierISA03 = ""
	SecQualifierISA0300      SecQualifierISA03 = "00"
	SecQualifierISA0301      SecQualifierISA03 = "01"
)

// secQualifierISA03Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var secQualifierISA03Codes = []codeEntry{
	{"00", "No Security Information Present (No Meaningful Information in I04)", "No Security Information Present (No Meaningful Information in I04) — X12 005010X221A1 TR3, ISA03 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text; the only value seen in this repo's example 835 files"},
	{"01", "Password", "Password — X12 005010X221A1 TR3, ISA03 (transcribed from the standard; verify against the printed guide)"},
}

func (v SecQualifierISA03) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v SecQualifierISA03) Known() bool { return known(secQualifierISA03Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v SecQualifierISA03) ShortDescription() string {
	return shortOf(secQualifierISA03Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v SecQualifierISA03) Description() string {
	return descriptionOf(secQualifierISA03Codes, string(v))
}

// ParseSecQualifierISA03 reads one ISA03 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseSecQualifierISA03(code string) (SecQualifierISA03, error) {
	return SecQualifierISA03(code), checkListed("ISA03", secQualifierISA03Codes, code)
}

// InterchangeIDQualifierISA05 is the Interchange ID Qualifier, used for both
// ISA05 (sender) and ISA07 (receiver) - CMS 835 Flat File, PDF p.1 gives
// both positions the same element name.
//
// Neither reviewed document names ISA05/07's individual codes; the names
// below are transcribed from the X12 standard.
type InterchangeIDQualifierISA05 string

const (
	InterchangeIDQualifierISA05NotSent InterchangeIDQualifierISA05 = ""
	InterchangeIDQualifierISA05ZZ      InterchangeIDQualifierISA05 = "ZZ"
	InterchangeIDQualifierISA0501      InterchangeIDQualifierISA05 = "01"
	InterchangeIDQualifierISA0514      InterchangeIDQualifierISA05 = "14"
	InterchangeIDQualifierISA0520      InterchangeIDQualifierISA05 = "20"
	InterchangeIDQualifierISA0527      InterchangeIDQualifierISA05 = "27"
	InterchangeIDQualifierISA0528      InterchangeIDQualifierISA05 = "28"
	InterchangeIDQualifierISA0529      InterchangeIDQualifierISA05 = "29"
	InterchangeIDQualifierISA0530      InterchangeIDQualifierISA05 = "30"
	InterchangeIDQualifierISA0533      InterchangeIDQualifierISA05 = "33"
)

// interchangeIDQualifierISA05Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var interchangeIDQualifierISA05Codes = []codeEntry{
	{"ZZ", "Mutually Defined", "Mutually Defined — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text; the only value seen in this repo's example 835 files"},
	{"01", "Duns (Dun & Bradstreet)", "Duns (Dun & Bradstreet) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"14", "Duns Plus Suffix", "Duns Plus Suffix — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"20", "Health Industry Number (HIN)", "Health Industry Number (HIN) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"27", "Carrier Identification Number as assigned by Health Care Financing Administration (HCFA)", "Carrier Identification Number as assigned by Health Care Financing Administration (HCFA) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"28", "Fiscal Intermediary Identification Number as assigned by Health Care Financing Administration (HCFA)", "Fiscal Intermediary Identification Number as assigned by Health Care Financing Administration (HCFA) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"29", "Medicare Provider and Supplier Identification Number as assigned by Health Care Financing Administration (HCFA)", "Medicare Provider and Supplier Identification Number as assigned by Health Care Financing Administration (HCFA) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"30", "U.S. Federal Tax Identification Number", "U.S. Federal Tax Identification Number — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
	{"33", "National Association of Insurance Commissioners Company Code (NAIC)", "National Association of Insurance Commissioners Company Code (NAIC) — X12 005010X221A1 TR3, ISA05 and ISA07 (transcribed from the standard; verify against the printed guide)"},
}

func (v InterchangeIDQualifierISA05) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v InterchangeIDQualifierISA05) Known() bool {
	return known(interchangeIDQualifierISA05Codes, string(v))
}

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v InterchangeIDQualifierISA05) ShortDescription() string {
	return shortOf(interchangeIDQualifierISA05Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v InterchangeIDQualifierISA05) Description() string {
	return descriptionOf(interchangeIDQualifierISA05Codes, string(v))
}

// ParseInterchangeIDQualifierISA05 reads one ISA05 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseInterchangeIDQualifierISA05(code string) (InterchangeIDQualifierISA05, error) {
	return InterchangeIDQualifierISA05(code), checkListed("ISA05", interchangeIDQualifierISA05Codes, code)
}

// AckRequestedISA14 is ISA14, the Acknowledgement Requested code.
//
// Neither reviewed document names ISA14's individual codes.
type AckRequestedISA14 string

const (
	AckRequestedISA14NotSent AckRequestedISA14 = ""
	AckRequestedISA140       AckRequestedISA14 = "0"
	AckRequestedISA141       AckRequestedISA14 = "1"
)

// ackRequestedISA14Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var ackRequestedISA14Codes = []codeEntry{
	{"0", "No Interchange Acknowledgment Requested", "No Interchange Acknowledgment Requested — X12 005010X221A1 TR3, ISA14 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1; not named in text; the only value seen in this repo's example 835 files"},
	{"1", "Interchange Acknowledgment Requested (TA1)", "Interchange Acknowledgment Requested (TA1) — X12 005010X221A1 TR3, ISA14 (transcribed from the standard; verify against the printed guide)"},
}

func (v AckRequestedISA14) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AckRequestedISA14) Known() bool { return known(ackRequestedISA14Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AckRequestedISA14) ShortDescription() string {
	return shortOf(ackRequestedISA14Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AckRequestedISA14) Description() string {
	return descriptionOf(ackRequestedISA14Codes, string(v))
}

// ParseAckRequestedISA14 reads one ISA14 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAckRequestedISA14(code string) (AckRequestedISA14, error) {
	return AckRequestedISA14(code), checkListed("ISA14", ackRequestedISA14Codes, code)
}

// UsageISA15 is ISA15, the Usage Indicator.
//
// Neither reviewed document names ISA15's individual codes.
type UsageISA15 string

const (
	UsageISA15NotSent UsageISA15 = ""
	UsageISA15P       UsageISA15 = "P"
	UsageISA15T       UsageISA15 = "T"
)

// usageISA15Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var usageISA15Codes = []codeEntry{
	{"P", "Production Data", "Production Data — X12 005010X221A1 TR3, ISA15 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1 (Values \"P,T\"); not named in text"},
	{"T", "Test Data", "Test Data — X12 005010X221A1 TR3, ISA15 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.1 (Values \"P,T\"); not named in text"},
}

func (v UsageISA15) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v UsageISA15) Known() bool { return known(usageISA15Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v UsageISA15) ShortDescription() string { return shortOf(usageISA15Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v UsageISA15) Description() string { return descriptionOf(usageISA15Codes, string(v)) }

// ParseUsageISA15 reads one ISA15 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseUsageISA15(code string) (UsageISA15, error) {
	return UsageISA15(code), checkListed("ISA15", usageISA15Codes, code)
}

// FunctionalIDGS01 is GS01, the Functional Identifier Code.
type FunctionalIDGS01 string

const (
	FunctionalIDGS01NotSent FunctionalIDGS01 = ""
	FunctionalIDGS01HP      FunctionalIDGS01 = "HP"
)

// functionalIDGS01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var functionalIDGS01Codes = []codeEntry{
	{"HP", "Health Care Claim Payment/Advice", "Health Care Claim Payment/Advice — CEDI 835 CG, PDF p.7 (printed p.1); frontend/src/data/x12Names835.ts TRANSACTION_TYPE_NAME"},
}

func (v FunctionalIDGS01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v FunctionalIDGS01) Known() bool { return known(functionalIDGS01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v FunctionalIDGS01) ShortDescription() string { return shortOf(functionalIDGS01Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v FunctionalIDGS01) Description() string {
	return descriptionOf(functionalIDGS01Codes, string(v))
}

// ParseFunctionalIDGS01 reads one GS01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseFunctionalIDGS01(code string) (FunctionalIDGS01, error) {
	return FunctionalIDGS01(code), checkListed("GS01", functionalIDGS01Codes, code)
}

// AgencyGS07 is GS07, the Responsible Agency Code.
//
// Neither reviewed document names GS07's individual codes.
type AgencyGS07 string

const (
	AgencyGS07NotSent AgencyGS07 = ""
	AgencyGS07X       AgencyGS07 = "X"
)

// agencyGS07Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var agencyGS07Codes = []codeEntry{
	{"X", "Accredited Standards Committee X12", "Accredited Standards Committee X12 — X12 005010X221A1 TR3, GS07 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.2; not named in text; the only value seen in this repo's example 835 files"},
}

func (v AgencyGS07) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AgencyGS07) Known() bool { return known(agencyGS07Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AgencyGS07) ShortDescription() string { return shortOf(agencyGS07Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AgencyGS07) Description() string { return descriptionOf(agencyGS07Codes, string(v)) }

// ParseAgencyGS07 reads one GS07 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAgencyGS07(code string) (AgencyGS07, error) {
	return AgencyGS07(code), checkListed("GS07", agencyGS07Codes, code)
}

// TransactionSetIDST01 is ST01, the Transaction Set Identifier Code.
type TransactionSetIDST01 string

const (
	TransactionSetIDST01NotSent TransactionSetIDST01 = ""
	TransactionSetIDST01835     TransactionSetIDST01 = "835"
)

// transactionSetIDST01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var transactionSetIDST01Codes = []codeEntry{
	{"835", "Health Care Claim Payment/Advice", "Health Care Claim Payment/Advice — CEDI 835 CG, PDF p.7 (printed p.1); frontend/src/data/x12Names835.ts TRANSACTION_TYPE_NAME"},
}

func (v TransactionSetIDST01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v TransactionSetIDST01) Known() bool { return known(transactionSetIDST01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v TransactionSetIDST01) ShortDescription() string {
	return shortOf(transactionSetIDST01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v TransactionSetIDST01) Description() string {
	return descriptionOf(transactionSetIDST01Codes, string(v))
}

// ParseTransactionSetIDST01 reads one ST01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseTransactionSetIDST01(code string) (TransactionSetIDST01, error) {
	return TransactionSetIDST01(code), checkListed("ST01", transactionSetIDST01Codes, code)
}
