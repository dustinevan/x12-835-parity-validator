package x12835

// EntityIDN101 is N101, the payer/payee role carried by an N1 segment.
type EntityIDN101 string

const (
	EntityIDN101NotSent EntityIDN101 = ""
	EntityIDN101PR      EntityIDN101 = "PR"
	EntityIDN101PE      EntityIDN101 = "PE"
)

// entityIDN101Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var entityIDN101Codes = []codeEntry{
	{"PR", "Payer Identification", "Payer Identification — CMS 835 Flat File, PDF p.2; frontend/src/data/x12Names835.ts N1_NAME"},
	{"PE", "Payee Identification", "Payee Identification — CMS 835 Flat File, PDF p.3; frontend/src/data/x12Names835.ts N1_NAME"},
}

func (v EntityIDN101) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v EntityIDN101) Known() bool { return known(entityIDN101Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v EntityIDN101) ShortDescription() string { return shortOf(entityIDN101Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v EntityIDN101) Description() string { return descriptionOf(entityIDN101Codes, string(v)) }

// ParseEntityIDN101 reads one N101 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseEntityIDN101(code string) (EntityIDN101, error) {
	return EntityIDN101(code), checkListed("N101", entityIDN101Codes, code)
}

// IDQualifierN103 is N103, the ID Code Qualifier on the payee's N1 segment
// (loop 1000B).
type IDQualifierN103 string

const (
	IDQualifierN103NotSent IDQualifierN103 = ""
	IDQualifierN103XX      IDQualifierN103 = "XX"
	IDQualifierN103FI      IDQualifierN103 = "FI"
	IDQualifierN103XV      IDQualifierN103 = "XV"
)

// idQualifierN103Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var idQualifierN103Codes = []codeEntry{
	{"XX", "National Provider Identifier", "National Provider Identifier — CMS 835 Flat File, PDF p.3, p.7 (NM109 shows the same XX to NPI mapping); CEDI 835 CG, PDF p.32 (printed p.26; Acronym List); frontend/src/data/x12Names835.ts ID_QUALIFIER_NAME"},
	{"FI", "Federal Taxpayer's Identification Number", "Federal Taxpayer's Identification Number — X12 005010X221A1 TR3, N103 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.3; not named in text (docs/design/835-examples.md \"Items not built\": a pre-NPI legacy qualifier the flat file still lists as valid)"},
	{"XV", "Centers for Medicare and Medicaid Services PlanID", "Centers for Medicare and Medicaid Services PlanID — X12 005010X221A1 TR3, N103 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.3; not named in text (same note as FI)"},
}

func (v IDQualifierN103) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v IDQualifierN103) Known() bool { return known(idQualifierN103Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v IDQualifierN103) ShortDescription() string { return shortOf(idQualifierN103Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v IDQualifierN103) Description() string { return descriptionOf(idQualifierN103Codes, string(v)) }

// ParseIDQualifierN103 reads one N103 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseIDQualifierN103(code string) (IDQualifierN103, error) {
	return IDQualifierN103(code), checkListed("N103", idQualifierN103Codes, code)
}

// ContactFunctionPER01 is PER01, the Contact Function Code.
//
// CX has two names in the two sourced tables:
// frontend/src/data/x12Names835.ts's PAYER_PER_NAME.CX names it
// "Payer Business Contact Information" (payer loop, 1000A; CMS 835 Flat
// File, PDF p.3) while the same file's CLAIM_PER_NAME.CX names it "Claim
// Contact Information" (claim loop, 2100; CMS 835 Flat File, PDF p.11) - the
// same one-letter-pair code, two different sourced meanings depending on
// which loop the PER segment sits in. A flat code-to-name list can only
// carry one label per code, so ContactFunctionPER01CX is named by its code
// rather than picking one of the two sourced meanings over the other; a
// struct that embeds this value alongside its loop (as the other 835 files
// in this package do) carries the actual meaning. BL and 1C have no such
// conflict: both sources show them only in the payer loop.
type ContactFunctionPER01 string

const (
	ContactFunctionPER01NotSent ContactFunctionPER01 = ""
	ContactFunctionPER01CX      ContactFunctionPER01 = "CX"
	ContactFunctionPER01BL      ContactFunctionPER01 = "BL"
	ContactFunctionPER011C      ContactFunctionPER01 = "1C"
	ContactFunctionPER01IC      ContactFunctionPER01 = "IC"
)

// contactFunctionPER01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var contactFunctionPER01Codes = []codeEntry{
	{"CX", "Payers Claim Office", "Payers Claim Office — X12 005010X221A1 TR3, PER01 (transcribed from the standard; verify against the printed guide); ambiguous, see the type doc comment; CMS 835 Flat File, PDF p.3, p.11"},
	{"BL", "Payer Technical Contact Information", "Payer Technical Contact Information — CMS 835 Flat File, PDF p.3; frontend/src/data/x12Names835.ts PAYER_PER_NAME (payer loop, 1000A, only)"},
	{"1C", "Payer Web Site", "Payer Web Site — CMS 835 Flat File, PDF p.3; frontend/src/data/x12Names835.ts PAYER_PER_NAME (payer loop, 1000A, only)"},
	{"IC", "Information Contact", "Information Contact — X12 005010X221A1 TR3, PER01 (transcribed from the standard; verify against the printed guide); the payer's web-site contact in loop 1000A; the 2023 CMS 835 Flat File (PDF p.3) and Transmittal 907 (PDF p.281) print IC, where the files on file that send 1C are kept under that constant"},
}

func (v ContactFunctionPER01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v ContactFunctionPER01) Known() bool { return known(contactFunctionPER01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v ContactFunctionPER01) ShortDescription() string {
	return shortOf(contactFunctionPER01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v ContactFunctionPER01) Description() string {
	return descriptionOf(contactFunctionPER01Codes, string(v))
}

// ParseContactFunctionPER01 reads one PER01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseContactFunctionPER01(code string) (ContactFunctionPER01, error) {
	return ContactFunctionPER01(code), checkListed("PER01", contactFunctionPER01Codes, code)
}

// CommQualifierPER03 is the Communication Number Qualifier, used for PER03,
// PER05, and PER07 alike.
type CommQualifierPER03 string

const (
	CommQualifierPER03NotSent CommQualifierPER03 = ""
	CommQualifierPER03TE      CommQualifierPER03 = "TE"
	CommQualifierPER03UR      CommQualifierPER03 = "UR"
	CommQualifierPER03EM      CommQualifierPER03 = "EM"
	CommQualifierPER03EX      CommQualifierPER03 = "EX"
	CommQualifierPER03FX      CommQualifierPER03 = "FX"
)

// commQualifierPER03Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var commQualifierPER03Codes = []codeEntry{
	{"TE", "Telephone", "Telephone — X12 005010X221A1 TR3, PER03, PER05, and PER07 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.3, p.11; not named in text (frontend/src/data/x12Names835.ts only names the PER03/05/07 positions themselves, e.g. \"Communication # Qualifier\", not their codes)"},
	{"UR", "Uniform Resource Locator (URL)", "Uniform Resource Locator (URL) — X12 005010X221A1 TR3, PER03, PER05, and PER07 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.3; not named in text (same note as TE)"},
	{"EM", "Electronic Mail", "Electronic Mail — X12 005010X221A1 TR3, PER03, PER05, and PER07 (transcribed from the standard; verify against the printed guide)"},
	{"EX", "Telephone Extension", "Telephone Extension — X12 005010X221A1 TR3, PER03, PER05, and PER07 (transcribed from the standard; verify against the printed guide)"},
	{"FX", "Facsimile", "Facsimile — X12 005010X221A1 TR3, PER03, PER05, and PER07 (transcribed from the standard; verify against the printed guide)"},
}

func (v CommQualifierPER03) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v CommQualifierPER03) Known() bool { return known(commQualifierPER03Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v CommQualifierPER03) ShortDescription() string {
	return shortOf(commQualifierPER03Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v CommQualifierPER03) Description() string {
	return descriptionOf(commQualifierPER03Codes, string(v))
}

// ParseCommQualifierPER03 reads one PER03 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseCommQualifierPER03(code string) (CommQualifierPER03, error) {
	return CommQualifierPER03(code), checkListed("PER03", commQualifierPER03Codes, code)
}

// DeliveryMethodRDM01 is RDM01, the Report Transmission Code on the
// Remittance Delivery Method segment (loop 1000B).
//
// docs/design/835-examples.md's payee-by-npi example lists these four as
// RDM01's valid set; neither reviewed document names any of them in text.
type DeliveryMethodRDM01 string

const (
	DeliveryMethodRDM01NotSent DeliveryMethodRDM01 = ""
	DeliveryMethodRDM01BM      DeliveryMethodRDM01 = "BM"
	DeliveryMethodRDM01EM      DeliveryMethodRDM01 = "EM"
	DeliveryMethodRDM01FT      DeliveryMethodRDM01 = "FT"
	DeliveryMethodRDM01OL      DeliveryMethodRDM01 = "OL"
)

// deliveryMethodRDM01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var deliveryMethodRDM01Codes = []codeEntry{
	{"BM", "By Mail", "By Mail — X12 005010X221A1 TR3, RDM01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.4; not named in text; the value used in this repo's payee-by-npi.835 example"},
	{"EM", "E-Mail", "E-Mail — X12 005010X221A1 TR3, RDM01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.4; not named in text"},
	{"FT", "File Transfer", "File Transfer — X12 005010X221A1 TR3, RDM01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.4; not named in text"},
	{"OL", "On-Line", "On-Line — X12 005010X221A1 TR3, RDM01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.4; not named in text"},
}

func (v DeliveryMethodRDM01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v DeliveryMethodRDM01) Known() bool { return known(deliveryMethodRDM01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v DeliveryMethodRDM01) ShortDescription() string {
	return shortOf(deliveryMethodRDM01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v DeliveryMethodRDM01) Description() string {
	return descriptionOf(deliveryMethodRDM01Codes, string(v))
}

// ParseDeliveryMethodRDM01 reads one RDM01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseDeliveryMethodRDM01(code string) (DeliveryMethodRDM01, error) {
	return DeliveryMethodRDM01(code), checkListed("RDM01", deliveryMethodRDM01Codes, code)
}
