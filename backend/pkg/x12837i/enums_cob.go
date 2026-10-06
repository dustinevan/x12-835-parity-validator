package x12837i

import "fmt"

// CASGroupCode01 is CAS01, the Claim Adjustment Group Code, shared by the
// 2320 Claim Level Adjustments CAS and the 2430 Line Adjustment CAS (the
// institutional claim's only two CAS instances, both in the coordination-
// of-benefits section).
type CASGroupCode01 string

const (
	CASGroupCode01CO CASGroupCode01 = "CO"
	CASGroupCode01CR CASGroupCode01 = "CR"
	CASGroupCode01OA CASGroupCode01 = "OA"
	CASGroupCode01PI CASGroupCode01 = "PI"
	CASGroupCode01PR CASGroupCode01 = "PR"
)

// The codes the guide allows for CAS01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var casGroupCode01Codes = []codeEntry{
	{"CO", "Contractual Obligation", "Contractual Obligation — x12.org/codes/claim-adjustment-group-codes"},
	{"CR", "Correction and Reversal", "Correction and Reversal — x12.org/codes/claim-adjustment-group-codes"},
	{"OA", "Other Adjustment", "Other Adjustment — x12.org/codes/claim-adjustment-group-codes"},
	{"PI", "Payer Initiated Reductions", "Payer Initiated Reductions — x12.org/codes/claim-adjustment-group-codes"},
	{"PR", "Patient Responsibility", "Patient Responsibility — x12.org/codes/claim-adjustment-group-codes"},
}

// String returns the code as the file sent it.
func (v CASGroupCode01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v CASGroupCode01) Known() bool { return knownCode(casGroupCode01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v CASGroupCode01) ShortDescription() string { return shortOf(casGroupCode01Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v CASGroupCode01) Description() string { return descriptionOf(casGroupCode01Codes, string(v)) }

// ParseCASGroupCode01 reads one CAS01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseCASGroupCode01(code string) (CASGroupCode01, error) {
	v := CASGroupCode01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CAS01: a value not on the list")
}

// OtherInsuranceIndicator merges OI03 (Benefits Assignment Certification
// Indicator) and OI06 (Release of Information Code) on the 2320 Other
// Health Insurance Information OI segment into one type: OI03 reuses the
// identical code set as CLM08 (BenefitsAssignCLM08) and OI06 reuses the
// identical code set as CLM09 (ReleaseOfInfoCLM09), so
// frontend/src/data/x12Names837i.ts's header lists them as one combined
// Go type, "OtherInsuranceIndicator (OI03/OI06)", rather than two.
type OtherInsuranceIndicator string

const (
	OtherInsuranceIndicatorN OtherInsuranceIndicator = "N"
	OtherInsuranceIndicatorW OtherInsuranceIndicator = "W"
	OtherInsuranceIndicatorY OtherInsuranceIndicator = "Y"
	OtherInsuranceIndicatorI OtherInsuranceIndicator = "I"
)

// The codes the guide allows for OI03/OI06, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var otherInsuranceIndicatorCodes = []codeEntry{
	{"N", "No", "No — Transmittal 907, PDF p.100 (OI03)"},
	{"W", "Not Applicable", "Not Applicable — Transmittal 907, PDF p.100 (OI03)"},
	{"Y", "Yes", "Yes — Transmittal 907, PDF p.100 (OI03)"},
	{"I", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes — Transmittal 907, PDF p.100 (OI06)"},
}

// String returns the code as the file sent it.
func (v OtherInsuranceIndicator) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v OtherInsuranceIndicator) Known() bool {
	return knownCode(otherInsuranceIndicatorCodes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v OtherInsuranceIndicator) ShortDescription() string {
	return shortOf(otherInsuranceIndicatorCodes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v OtherInsuranceIndicator) Description() string {
	return descriptionOf(otherInsuranceIndicatorCodes, string(v))
}

// ParseOtherInsuranceIndicator reads one OI03/OI06 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseOtherInsuranceIndicator(code string) (OtherInsuranceIndicator, error) {
	v := OtherInsuranceIndicator(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing OI03/OI06: a value not on the list")
}
