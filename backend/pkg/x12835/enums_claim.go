package x12835

// ClaimStatusCLP02 is CLP02, the Claim Status Code.
//
// None of CLP02's nine values is named in text by either reviewed document
// (docs/design/835-examples.md's Unnamed codes list "CLP02's 9 values"
// together); docs/design/835-examples.md's examples establish each value's
// real-world role (e.g. 22 is a reversal, 1 the following correction;
// 2/3/20/21/23 are a coordination-of-benefits chain) without ever citing a
// published short name for the code itself, so every value here is named by
// its code.
type ClaimStatusCLP02 string

const (
	ClaimStatusCLP02NotSent ClaimStatusCLP02 = ""
	ClaimStatusCLP021       ClaimStatusCLP02 = "1"
	ClaimStatusCLP022       ClaimStatusCLP02 = "2"
	ClaimStatusCLP023       ClaimStatusCLP02 = "3"
	ClaimStatusCLP024       ClaimStatusCLP02 = "4"
	ClaimStatusCLP0219      ClaimStatusCLP02 = "19"
	ClaimStatusCLP0220      ClaimStatusCLP02 = "20"
	ClaimStatusCLP0221      ClaimStatusCLP02 = "21"
	ClaimStatusCLP0222      ClaimStatusCLP02 = "22"
	ClaimStatusCLP0223      ClaimStatusCLP02 = "23"
	ClaimStatusCLP0225      ClaimStatusCLP02 = "25"
)

// claimStatusCLP02Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var claimStatusCLP02Codes = []codeEntry{
	{"1", "Processed as Primary", "Processed as Primary — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text (Open item)"},
	{"2", "Processed as Secondary", "Processed as Secondary — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"3", "Processed as Tertiary", "Processed as Tertiary — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"4", "Denied", "Denied — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"19", "Processed as Primary, Forwarded to Additional Payer(s)", "Processed as Primary, Forwarded to Additional Payer(s) — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"20", "Processed as Secondary, Forwarded to Additional Payer(s)", "Processed as Secondary, Forwarded to Additional Payer(s) — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"21", "Processed as Tertiary, Forwarded to Additional Payer(s)", "Processed as Tertiary, Forwarded to Additional Payer(s) — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"22", "Reversal of Previous Payment", "Reversal of Previous Payment — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"23", "Not Our Claim, Forwarded to Additional Payer(s)", "Not Our Claim, Forwarded to Additional Payer(s) — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.5; CEDI 835 CG, PDF p.25 (printed p.19); not named in text"},
	{"25", "Predetermination Pricing Only - No Payment", "Predetermination Pricing Only - No Payment — X12 005010X221A1 TR3, CLP02 (transcribed from the standard; verify against the printed guide); not in the CMS 835 Flat File's list (PDF p.5), so Medicare is not documented as sending it"},
}

func (v ClaimStatusCLP02) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v ClaimStatusCLP02) Known() bool { return known(claimStatusCLP02Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v ClaimStatusCLP02) ShortDescription() string { return shortOf(claimStatusCLP02Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v ClaimStatusCLP02) Description() string {
	return descriptionOf(claimStatusCLP02Codes, string(v))
}

// ParseClaimStatusCLP02 reads one CLP02 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseClaimStatusCLP02(code string) (ClaimStatusCLP02, error) {
	return ClaimStatusCLP02(code), checkListed("CLP02", claimStatusCLP02Codes, code)
}

// FilingIndicatorCLP06 is CLP06, the Claim Filing Indicator Code.
type FilingIndicatorCLP06 string

const (
	FilingIndicatorCLP06NotSent FilingIndicatorCLP06 = ""
	FilingIndicatorCLP06MA      FilingIndicatorCLP06 = "MA"
	FilingIndicatorCLP06MB      FilingIndicatorCLP06 = "MB"
	FilingIndicatorCLP0611      FilingIndicatorCLP06 = "11"
	FilingIndicatorCLP0612      FilingIndicatorCLP06 = "12"
	FilingIndicatorCLP0613      FilingIndicatorCLP06 = "13"
	FilingIndicatorCLP0614      FilingIndicatorCLP06 = "14"
	FilingIndicatorCLP0615      FilingIndicatorCLP06 = "15"
	FilingIndicatorCLP0616      FilingIndicatorCLP06 = "16"
	FilingIndicatorCLP0617      FilingIndicatorCLP06 = "17"
	FilingIndicatorCLP06AM      FilingIndicatorCLP06 = "AM"
	FilingIndicatorCLP06CH      FilingIndicatorCLP06 = "CH"
	FilingIndicatorCLP06DS      FilingIndicatorCLP06 = "DS"
	FilingIndicatorCLP06HM      FilingIndicatorCLP06 = "HM"
	FilingIndicatorCLP06LM      FilingIndicatorCLP06 = "LM"
	FilingIndicatorCLP06MC      FilingIndicatorCLP06 = "MC"
	FilingIndicatorCLP06OF      FilingIndicatorCLP06 = "OF"
	FilingIndicatorCLP06TV      FilingIndicatorCLP06 = "TV"
	FilingIndicatorCLP06VA      FilingIndicatorCLP06 = "VA"
	FilingIndicatorCLP06WC      FilingIndicatorCLP06 = "WC"
	FilingIndicatorCLP06ZZ      FilingIndicatorCLP06 = "ZZ"
)

// filingIndicatorCLP06Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var filingIndicatorCLP06Codes = []codeEntry{
	{"MA", "Medicare Part A", "Medicare Part A — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.6; CEDI 835 CG, PDF p.25 (printed p.19; Table 9); not named in text"},
	{"MB", "Medicare Part B", "Medicare Part B — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.6; CEDI 835 CG, PDF p.25 (printed p.19; Table 9); not named in text; CEDI 835 CG states Medicare sends MB for DME"},
	{"11", "Other Non-Federal Programs", "Other Non-Federal Programs — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide); (verify) whether the 835 guide's CLP06 list includes 11; the 837's SBR09 list does"},
	{"12", "Preferred Provider Organization (PPO)", "Preferred Provider Organization (PPO) — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"13", "Point of Service (POS)", "Point of Service (POS) — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"14", "Exclusive Provider Organization (EPO)", "Exclusive Provider Organization (EPO) — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"15", "Indemnity Insurance", "Indemnity Insurance — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"16", "Health Maintenance Organization (HMO) Medicare Risk", "Health Maintenance Organization (HMO) Medicare Risk — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"17", "Dental Maintenance Organization", "Dental Maintenance Organization — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"AM", "Automobile Medical", "Automobile Medical — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"CH", "Champus", "Champus — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"DS", "Disability", "Disability — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"HM", "Health Maintenance Organization", "Health Maintenance Organization — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"LM", "Liability Medical", "Liability Medical — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"MC", "Medicaid", "Medicaid — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"OF", "Other Federal Program", "Other Federal Program — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"TV", "Title V", "Title V — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"VA", "Veterans Affairs Plan", "Veterans Affairs Plan — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"WC", "Workers' Compensation Health Claim", "Workers' Compensation Health Claim — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
	{"ZZ", "Mutually Defined", "Mutually Defined — X12 005010X221A1 TR3, CLP06 (transcribed from the standard; verify against the printed guide)"},
}

func (v FilingIndicatorCLP06) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v FilingIndicatorCLP06) Known() bool { return known(filingIndicatorCLP06Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v FilingIndicatorCLP06) ShortDescription() string {
	return shortOf(filingIndicatorCLP06Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v FilingIndicatorCLP06) Description() string {
	return descriptionOf(filingIndicatorCLP06Codes, string(v))
}

// ParseFilingIndicatorCLP06 reads one CLP06 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseFilingIndicatorCLP06(code string) (FilingIndicatorCLP06, error) {
	return FilingIndicatorCLP06(code), checkListed("CLP06", filingIndicatorCLP06Codes, code)
}

// EntityIDNM101 is NM101, the Entity Identifier Code on a claim-loop (2100)
// NM1 segment.
//
// Two further claim-loop NM1 roles - Insured's Name and Other Subscriber
// Name (CMS 835 Flat File, PDF p.7-8; CEDI 835 CG Table 12) - are documented
// as not used by Medicare, and neither reviewed document gives their NM101
// code; IL and GB are listed with their names from the X12 standard.
type EntityIDNM101 string

const (
	EntityIDNM101NotSent EntityIDNM101 = ""
	EntityIDNM101QC      EntityIDNM101 = "QC"
	EntityIDNM10174      EntityIDNM101 = "74"
	EntityIDNM10182      EntityIDNM101 = "82"
	EntityIDNM101TT      EntityIDNM101 = "TT"
	EntityIDNM101PR      EntityIDNM101 = "PR"
	EntityIDNM101IL      EntityIDNM101 = "IL"
	EntityIDNM101GB      EntityIDNM101 = "GB"
)

// entityIDNM101Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var entityIDNM101Codes = []codeEntry{
	{"QC", "Patient Name", "Patient Name — CMS 835 Flat File, PDF p.6; frontend/src/data/x12Names835.ts NM1_NAME"},
	{"74", "Corrected Patient/Insured Name", "Corrected Patient/Insured Name — CMS 835 Flat File, PDF p.7; frontend/src/data/x12Names835.ts NM1_NAME"},
	{"82", "Service Provider Name", "Service Provider Name — CMS 835 Flat File, PDF p.7; frontend/src/data/x12Names835.ts NM1_NAME"},
	{"TT", "Crossover Carrier Name", "Crossover Carrier Name — CMS 835 Flat File, PDF p.7; frontend/src/data/x12Names835.ts NM1_NAME"},
	{"PR", "Corrected Priority Payer Name", "Corrected Priority Payer Name — CMS 835 Flat File, PDF p.7-8; frontend/src/data/x12Names835.ts NM1_NAME"},
	{"IL", "Insured or Subscriber", "Insured or Subscriber — X12 005010X221A1 TR3, NM101 (loop 2100) (transcribed from the standard; verify against the printed guide); the 2100 Insured Name NM1; this model has no field for it, so the segment is kept raw"},
	{"GB", "Other Insured", "Other Insured — X12 005010X221A1 TR3, NM101 (loop 2100) (transcribed from the standard; verify against the printed guide); the 2100 Other Subscriber Name NM1; kept raw like IL"},
}

func (v EntityIDNM101) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v EntityIDNM101) Known() bool { return known(entityIDNM101Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v EntityIDNM101) ShortDescription() string { return shortOf(entityIDNM101Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v EntityIDNM101) Description() string { return descriptionOf(entityIDNM101Codes, string(v)) }

// ParseEntityIDNM101 reads one NM101 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseEntityIDNM101(code string) (EntityIDNM101, error) {
	return EntityIDNM101(code), checkListed("NM101", entityIDNM101Codes, code)
}

// EntityTypeNM102 is NM102, the Entity Type Qualifier on an NM1 segment (1 =
// person, 2 = non-person entity, by the segment's own structure).
type EntityTypeNM102 string

const (
	EntityTypeNM102NotSent EntityTypeNM102 = ""
	EntityTypeNM1021       EntityTypeNM102 = "1"
	EntityTypeNM1022       EntityTypeNM102 = "2"
)

// entityTypeNM102Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var entityTypeNM102Codes = []codeEntry{
	{"1", "Person", "Person — X12 005010X221A1 TR3, NM102 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.6-8; not named in text (frontend/src/data/x12Names835.ts only names the NM102 position itself, \"Entity Type Qualifier,\" not its codes)"},
	{"2", "Non-Person Entity", "Non-Person Entity — X12 005010X221A1 TR3, NM102 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.6-8; not named in text"},
}

func (v EntityTypeNM102) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v EntityTypeNM102) Known() bool { return known(entityTypeNM102Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v EntityTypeNM102) ShortDescription() string { return shortOf(entityTypeNM102Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v EntityTypeNM102) Description() string { return descriptionOf(entityTypeNM102Codes, string(v)) }

// ParseEntityTypeNM102 reads one NM102 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseEntityTypeNM102(code string) (EntityTypeNM102, error) {
	return EntityTypeNM102(code), checkListed("NM102", entityTypeNM102Codes, code)
}

// IDQualifierNM108 is NM108, the ID Code Qualifier on a claim-loop NM1
// segment.
type IDQualifierNM108 string

const (
	IDQualifierNM108NotSent IDQualifierNM108 = ""
	IDQualifierNM108MI      IDQualifierNM108 = "MI"
	IDQualifierNM108XX      IDQualifierNM108 = "XX"
	IDQualifierNM10834      IDQualifierNM108 = "34"
	IDQualifierNM108BD      IDQualifierNM108 = "BD"
	IDQualifierNM108BS      IDQualifierNM108 = "BS"
	IDQualifierNM108C       IDQualifierNM108 = "C"
	IDQualifierNM108FI      IDQualifierNM108 = "FI"
	IDQualifierNM108HN      IDQualifierNM108 = "HN"
	IDQualifierNM108II      IDQualifierNM108 = "II"
	IDQualifierNM108MC      IDQualifierNM108 = "MC"
	IDQualifierNM108MR      IDQualifierNM108 = "MR"
	IDQualifierNM108PC      IDQualifierNM108 = "PC"
	IDQualifierNM108PI      IDQualifierNM108 = "PI"
	IDQualifierNM108SL      IDQualifierNM108 = "SL"
	IDQualifierNM108UP      IDQualifierNM108 = "UP"
	IDQualifierNM108XV      IDQualifierNM108 = "XV"
)

// idQualifierNM108Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var idQualifierNM108Codes = []codeEntry{
	{"MI", "Member Identification Number", "Member Identification Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.6; not named in text (docs/design/835-examples.md's Open items list \"MI\" by name as an unresolved qualifier)"},
	{"XX", "National Provider Identifier", "National Provider Identifier — CMS 835 Flat File, PDF p.7 (NM109 content shown as \"NPI\"); CEDI 835 CG, PDF p.32 (printed p.26; Acronym List); frontend/src/data/x12Names835.ts ID_QUALIFIER_NAME"},
	{"34", "Social Security Number", "Social Security Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"BD", "Blue Cross Provider Number", "Blue Cross Provider Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"BS", "Blue Shield Provider Number", "Blue Shield Provider Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"C", "Insured's Changed Unique Identification Number", "Insured's Changed Unique Identification Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"FI", "Federal Taxpayer's Identification Number", "Federal Taxpayer's Identification Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"HN", "Health Insurance Claim (HIC) Number", "Health Insurance Claim (HIC) Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"II", "Standard Unique Health Identifier for each Individual in the United States", "Standard Unique Health Identifier for each Individual in the United States — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"MC", "Medicaid Provider Number", "Medicaid Provider Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"MR", "Medicaid Recipient Identification Number", "Medicaid Recipient Identification Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"PC", "Provider Commercial Number", "Provider Commercial Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"PI", "Payor Identification", "Payor Identification — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"SL", "State License Number", "State License Number — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"UP", "Unique Physician Identification Number (UPIN)", "Unique Physician Identification Number (UPIN) — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
	{"XV", "Centers for Medicare and Medicaid Services PlanID", "Centers for Medicare and Medicaid Services PlanID — X12 005010X221A1 TR3, NM108 (loop 2100) (transcribed from the standard; verify against the printed guide)"},
}

func (v IDQualifierNM108) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v IDQualifierNM108) Known() bool { return known(idQualifierNM108Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v IDQualifierNM108) ShortDescription() string { return shortOf(idQualifierNM108Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v IDQualifierNM108) Description() string {
	return descriptionOf(idQualifierNM108Codes, string(v))
}

// ParseIDQualifierNM108 reads one NM108 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseIDQualifierNM108(code string) (IDQualifierNM108, error) {
	return IDQualifierNM108(code), checkListed("NM108", idQualifierNM108Codes, code)
}
