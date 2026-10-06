package x12837i

import "fmt"

// ProviderRoleCodePRV01 is PRV01, the Provider Code identifying which
// provider a PRV segment describes: billing (2000A) or attending
// (2310A), the only two PRV instances in the institutional claim.
type ProviderRoleCodePRV01 string

const (
	ProviderRoleCodePRV01BI ProviderRoleCodePRV01 = "BI"
	ProviderRoleCodePRV01AT ProviderRoleCodePRV01 = "AT"
)

// The codes the guide allows for PRV01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var providerRoleCodePRV01Codes = []codeEntry{
	{"BI", "Billing", "Billing — Transmittal 907, PDF p.10 (2000A Billing Provider Specialty Information); PDF p.92 (2310A Attending Provider Specialty Information)"},
	{"AT", "Attending", "Attending — Transmittal 907, PDF p.10 (2000A Billing Provider Specialty Information); PDF p.92 (2310A Attending Provider Specialty Information)"},
}

// String returns the code as the file sent it.
func (v ProviderRoleCodePRV01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v ProviderRoleCodePRV01) Known() bool { return knownCode(providerRoleCodePRV01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v ProviderRoleCodePRV01) ShortDescription() string {
	return shortOf(providerRoleCodePRV01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v ProviderRoleCodePRV01) Description() string {
	return descriptionOf(providerRoleCodePRV01Codes, string(v))
}

// ParseProviderRoleCodePRV01 reads one PRV01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseProviderRoleCodePRV01(code string) (ProviderRoleCodePRV01, error) {
	v := ProviderRoleCodePRV01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing PRV01: a value not on the list")
}
