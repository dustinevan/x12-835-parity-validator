package x12837p

// Provider- and contact-level coded elements: PRV (provider specialty) and
// PER (contact) segments, shared across the billing, rendering, ordering,
// and service-facility provider loops.

// ProviderTypePRV01 is PRV01, the Provider Code, on a PRV segment (2000A
// Billing Provider, 2310A/2420F Referring Provider, 2310B/2420A Rendering
// Provider Specialty Information).
//
// The list is the three codes the 837P implementation guide allows at PRV01:
// BI at the billing provider, PE at a rendering provider, and RF at a
// referring provider. Neither Transmittal 907 nor the CEDI 837P Companion
// Guide names PRV01's individual codes in text; frontend/src/data/x12Names837p.ts
// PRV_EL_NAME gives only the element's own position name, "Provider Code," at
// every PRV instance. Transmittal 579's Values column lists BI at 2000A and
// PE at 2310B and 2420A (PDF p.9, 58, 98); the names are the guide's.
type ProviderTypePRV01 string

const (
	ProviderTypePRV01NotSent ProviderTypePRV01 = ""
	ProviderTypePRV01BI      ProviderTypePRV01 = "BI"
	ProviderTypePRV01PE      ProviderTypePRV01 = "PE"
	ProviderTypePRV01RF      ProviderTypePRV01 = "RF"
)

var providerTypePRV01List = codeList{
	{"BI", "Billing", "Billing — the 2000A Billing Provider Specialty Information box's only PRV01 value, Transmittal 579, PDF p.9 (Values column); the name is the guide's: " + tr3},
	{"PE", "Performing", "Performing — the only PRV01 value of the 2310B and 2420A Rendering Provider Specialty Information boxes, Transmittal 579, PDF p.58, 98 (Values column); the name is the guide's: " + tr3},
	{"RF", "Referring", guide("Referring", "PRV01 at 2310A and 2420F")},
}

func (v ProviderTypePRV01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ProviderTypePRV01) Known() bool { return providerTypePRV01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ProviderTypePRV01) ShortDescription() string { return providerTypePRV01List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ProviderTypePRV01) Description() string { return providerTypePRV01List.description(string(v)) }

// ParseProviderTypePRV01 returns text as a ProviderTypePRV01, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseProviderTypePRV01(text string) (ProviderTypePRV01, error) {
	return ProviderTypePRV01(text), providerTypePRV01List.check("PRV01", text)
}

// ContactFunctionCodePER01 is PER01, the Contact Function Code, shared by every
// PER segment in the transaction (1000A Submitter, 2010AA Billing Provider,
// 2310C Service Facility Location, 2420E Ordering Provider contact boxes).
//
// The 837P implementation guide allows the one code, IC, at every PER.
// Transmittal 579's Values column lists it at every PER box (PDF p.8, 13, 20,
// 26, 60, 104), and the CEDI worked transaction sends it ("PER*IC*..."); the
// name is the guide's.
type ContactFunctionCodePER01 string

const (
	ContactFunctionCodePER01NotSent ContactFunctionCodePER01 = ""
	ContactFunctionCodePER01IC      ContactFunctionCodePER01 = "IC"
)

var contactFunctionCodePER01List = codeList{
	{"IC", "Information Contact", "Information Contact — every PER box's only PER01 value, Transmittal 579, PDF p.8, 13, 20, 26, 60, 104 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"PER*IC*CONTACT NAME*TE*...\"); the name is the guide's: " + tr3},
}

func (v ContactFunctionCodePER01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ContactFunctionCodePER01) Known() bool { return contactFunctionCodePER01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ContactFunctionCodePER01) ShortDescription() string {
	return contactFunctionCodePER01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ContactFunctionCodePER01) Description() string {
	return contactFunctionCodePER01List.description(string(v))
}

// ParseContactFunctionCodePER01 returns text as a ContactFunctionCodePER01,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseContactFunctionCodePER01(text string) (ContactFunctionCodePER01, error) {
	return ContactFunctionCodePER01(text), contactFunctionCodePER01List.check("PER01", text)
}

// CommunicationNumberQualifierPER03 is PER03 (the same position repeated at
// PER05 and PER07), the Communication Number Qualifier, shared by every PER
// segment in the transaction.
//
// The list is the four codes the 837P implementation guide allows across
// PER03, PER05, and PER07 (EX only at PER05 and PER07). Transmittal 579's
// Values column lists "EM, EX, FX, TE" for the three positions (PDF p.8, 13,
// 104; fewer at some boxes), and the CEDI worked transaction sends TE; the
// names are the guide's.
type CommunicationNumberQualifierPER03 string

const (
	CommunicationNumberQualifierPER03NotSent CommunicationNumberQualifierPER03 = ""
	CommunicationNumberQualifierPER03EM      CommunicationNumberQualifierPER03 = "EM"
	CommunicationNumberQualifierPER03EX      CommunicationNumberQualifierPER03 = "EX"
	CommunicationNumberQualifierPER03FX      CommunicationNumberQualifierPER03 = "FX"
	CommunicationNumberQualifierPER03TE      CommunicationNumberQualifierPER03 = "TE"
)

var communicationNumberQualifierPER03List = codeList{
	{"EM", "Electronic Mail", "Electronic Mail — Transmittal 579, PDF p.8, 13, 104 (Values column, \"EM, EX, FX, TE\"); the name is the guide's: " + tr3},
	{"EX", "Telephone Extension", "Telephone Extension — Transmittal 579, PDF p.8, 13, 104 (Values column, \"EM, EX, FX, TE\"); the name is the guide's: " + tr3},
	{"FX", "Facsimile", "Facsimile — Transmittal 579, PDF p.8, 13, 104 (Values column, \"EM, EX, FX, TE\"); the name is the guide's: " + tr3},
	{"TE", "Telephone", "Telephone — Transmittal 579, PDF p.8, 13, 104 (Values column, \"EM, EX, FX, TE\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"PER*IC*CONTACT NAME*TE*...\"); the name is the guide's: " + tr3},
}

func (v CommunicationNumberQualifierPER03) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v CommunicationNumberQualifierPER03) Known() bool {
	return communicationNumberQualifierPER03List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v CommunicationNumberQualifierPER03) ShortDescription() string {
	return communicationNumberQualifierPER03List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v CommunicationNumberQualifierPER03) Description() string {
	return communicationNumberQualifierPER03List.description(string(v))
}

// ParseCommunicationNumberQualifierPER03 returns text as a
// CommunicationNumberQualifierPER03, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseCommunicationNumberQualifierPER03(text string) (CommunicationNumberQualifierPER03, error) {
	return CommunicationNumberQualifierPER03(text), communicationNumberQualifierPER03List.check("PER03", text)
}
