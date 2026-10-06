package x12837i

import "fmt"

// FacilityCodeQualifierCLM05_2 is CLM05-2, the Facility Code Qualifier,
// the second component of the CLM05 composite. Transmittal 907 lists only
// one value.
type FacilityCodeQualifierCLM05_2 string

const (
	FacilityCodeQualifierCLM05_2A FacilityCodeQualifierCLM05_2 = "A"
)

// The codes the guide allows for CLM05-2, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var facilityCodeQualifierCLM05_2Codes = []codeEntry{
	{"A", "Uniform Billing Claim Form Bill Type", "Uniform Billing Claim Form Bill Type — Transmittal 907, PDF p.23 (CLM05-2's only listed value)"},
}

// String returns the code as the file sent it.
func (v FacilityCodeQualifierCLM05_2) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v FacilityCodeQualifierCLM05_2) Known() bool {
	return knownCode(facilityCodeQualifierCLM05_2Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v FacilityCodeQualifierCLM05_2) ShortDescription() string {
	return shortOf(facilityCodeQualifierCLM05_2Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v FacilityCodeQualifierCLM05_2) Description() string {
	return descriptionOf(facilityCodeQualifierCLM05_2Codes, string(v))
}

// ParseFacilityCodeQualifierCLM05_2 reads one CLM05-2 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseFacilityCodeQualifierCLM05_2(code string) (FacilityCodeQualifierCLM05_2, error) {
	v := FacilityCodeQualifierCLM05_2(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CLM05-2: a value not on the list")
}
