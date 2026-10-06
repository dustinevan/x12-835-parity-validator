package x12835

// ProductQualifierSVC01 is the first component of the SVC01 (and, for the
// as-billed composite, SVC06) Composite Medical Procedure Identifier - the
// code-system qualifier in front of the procedure code.
//
// SVC06 reuses the identical qualifier set (docs/design/835-examples.md
// S-SVC06: "same qualifier set as SVC01").
type ProductQualifierSVC01 string

const (
	ProductQualifierSVC01NotSent ProductQualifierSVC01 = ""
	ProductQualifierSVC01HC      ProductQualifierSVC01 = "HC"
	ProductQualifierSVC01HP      ProductQualifierSVC01 = "HP"
	ProductQualifierSVC01N4      ProductQualifierSVC01 = "N4"
	ProductQualifierSVC01NU      ProductQualifierSVC01 = "NU"
	ProductQualifierSVC01AD      ProductQualifierSVC01 = "AD"
	ProductQualifierSVC01ER      ProductQualifierSVC01 = "ER"
	ProductQualifierSVC01IV      ProductQualifierSVC01 = "IV"
	ProductQualifierSVC01WK      ProductQualifierSVC01 = "WK"
)

// productQualifierSVC01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var productQualifierSVC01Codes = []codeEntry{
	{"HC", "Healthcare Common Procedure Coding System", "Healthcare Common Procedure Coding System — CEDI 835 CG, PDF p.31 (printed p.25; Acronym List, \"HCPCS\"); CMS 835 Flat File, PDF p.11 lists HC as a Product/Service ID Qualifier value; frontend/src/data/x12Names835.ts SVC_QUALIFIER_NAME"},
	{"HP", "Health Insurance Prospective Payment System (HIPPS) Skilled Nursing Facility Rate Code", "Health Insurance Prospective Payment System (HIPPS) Skilled Nursing Facility Rate Code — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; CEDI 835 CG, PDF p.27 (printed p.21); not named in text; used for DME alongside HC (docs/design/835-examples.md dme-rental)"},
	{"N4", "National Drug Code in 5-4-2 Format", "National Drug Code in 5-4-2 Format — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; CEDI 835 CG, PDF p.27 (printed p.21); not named in text. docs/design/835-examples.md glosses this as a National Drug Code \"by common HIPAA code-set usage,\" but says so explicitly is \"not spelled out in either document\" - that gloss is not a published name, so it is not used as this constant's name"},
	{"NU", "National Uniform Billing Committee (NUBC) UB92 Codes", "National Uniform Billing Committee (NUBC) UB92 Codes — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; CEDI 835 CG, PDF p.27 (printed p.21); not named in text"},
	{"AD", "American Dental Association Codes", "American Dental Association Codes — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide)"},
	{"ER", "Jurisdiction Specific Procedure and Supply Codes", "Jurisdiction Specific Procedure and Supply Codes — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide)"},
	{"IV", "Home Infusion EDI Coalition (HIEC) Product/Service Code", "Home Infusion EDI Coalition (HIEC) Product/Service Code — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide)"},
	{"WK", "Advanced Billing Concepts (ABC) Codes", "Advanced Billing Concepts (ABC) Codes — X12 005010X221A1 TR3, SVC01-1 and SVC06-1 (transcribed from the standard; verify against the printed guide)"},
}

func (v ProductQualifierSVC01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v ProductQualifierSVC01) Known() bool { return known(productQualifierSVC01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v ProductQualifierSVC01) ShortDescription() string {
	return shortOf(productQualifierSVC01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v ProductQualifierSVC01) Description() string {
	return descriptionOf(productQualifierSVC01Codes, string(v))
}

// ParseProductQualifierSVC01 reads one SVC01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseProductQualifierSVC01(code string) (ProductQualifierSVC01, error) {
	return ProductQualifierSVC01(code), checkListed("SVC01", productQualifierSVC01Codes, code)
}

// RemarkQualifierLQ01 is LQ01, the Code List Qualifier Code on the service-
// line LQ segment.
type RemarkQualifierLQ01 string

const (
	RemarkQualifierLQ01NotSent RemarkQualifierLQ01 = ""
	RemarkQualifierLQ01HE      RemarkQualifierLQ01 = "HE"
	RemarkQualifierLQ01RX      RemarkQualifierLQ01 = "RX"
)

// remarkQualifierLQ01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var remarkQualifierLQ01Codes = []codeEntry{
	{"HE", "Health Care Remarks Codes", "Health Care Remarks Codes — CMS 835 Flat File, PDF p.13; CEDI 835 CG, PDF p.28 (printed p.22; Table 24); frontend/src/data/x12Names835.ts LQ_NAME"},
	{"RX", "National Council for Prescription Drug Programs Reject/Payment Codes", "National Council for Prescription Drug Programs Reject/Payment Codes — X12 005010X221A1 TR3, LQ01 (transcribed from the standard; verify against the printed guide)"},
}

func (v RemarkQualifierLQ01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v RemarkQualifierLQ01) Known() bool { return known(remarkQualifierLQ01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v RemarkQualifierLQ01) ShortDescription() string {
	return shortOf(remarkQualifierLQ01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v RemarkQualifierLQ01) Description() string {
	return descriptionOf(remarkQualifierLQ01Codes, string(v))
}

// ParseRemarkQualifierLQ01 reads one LQ01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseRemarkQualifierLQ01(code string) (RemarkQualifierLQ01, error) {
	return RemarkQualifierLQ01(code), checkListed("LQ01", remarkQualifierLQ01Codes, code)
}
