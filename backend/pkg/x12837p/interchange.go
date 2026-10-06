package x12837p

// Delimiters are the two bytes the ISA fixes by position rather than by
// element, plus the file's line breaks: Element is the element separator
// (the byte after "ISA"), Segment the segment terminator (the byte after
// ISA16), and AfterSegment the whitespace the file puts after every segment
// terminator ("\n", "\r\n", or "" when the segments run together). The
// component separator is ISA16 and the repetition separator ISA11. Parse
// fills all three from the file; Write uses them, so a file writes back
// with its own layout. AfterSegment is "" when the file's gaps are not all
// the same, and the writer then puts nothing between segments.
type Delimiters struct {
	Element      string `json:"element,omitempty"`
	Segment      string `json:"segment,omitempty"`
	AfterSegment string `json:"afterSegment,omitempty"`
}

// The ISA01/03/05/07/14/15 qualifiers below are generic X12 005010
// interchange-envelope elements, not specific to the 837P transaction they
// carry. Each list is the one the 837P implementation guide allows.
// Transmittal 907's 837P attachment does not restate their code characters;
// Transmittal 579's 837P attachment lists some in its Values column (PDF
// p.6); the names are the guide's.

// AuthQualifierISA01 is ISA01, the Authorization Information Qualifier.
type AuthQualifierISA01 string

const (
	AuthQualifierISA01NotSent AuthQualifierISA01 = ""
	AuthQualifierISA0100      AuthQualifierISA01 = "00"
	AuthQualifierISA0103      AuthQualifierISA01 = "03"
)

var authQualifierISA01List = codeList{
	{"00", "No Authorization Information Present", "No Authorization Information Present (No Meaningful Information in I02) — Transmittal 579, PDF p.6 (Values column, \"00, 03\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
	{"03", "Additional Data Identification", "Additional Data Identification — Transmittal 579, PDF p.6 (Values column, \"00, 03\"); the name is the guide's: " + tr3},
}

func (v AuthQualifierISA01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AuthQualifierISA01) Known() bool { return authQualifierISA01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AuthQualifierISA01) ShortDescription() string { return authQualifierISA01List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AuthQualifierISA01) Description() string {
	return authQualifierISA01List.description(string(v))
}

// ParseAuthQualifierISA01 returns text as an AuthQualifierISA01, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseAuthQualifierISA01(text string) (AuthQualifierISA01, error) {
	return AuthQualifierISA01(text), authQualifierISA01List.check("ISA01", text)
}

// SecQualifierISA03 is ISA03, the Security Information Qualifier.
type SecQualifierISA03 string

const (
	SecQualifierISA03NotSent SecQualifierISA03 = ""
	SecQualifierISA0300      SecQualifierISA03 = "00"
	SecQualifierISA0301      SecQualifierISA03 = "01"
)

var secQualifierISA03List = codeList{
	{"00", "No Security Information Present", "No Security Information Present (No Meaningful Information in I04) — Transmittal 579, PDF p.6 (Values column, \"00, 01\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
	{"01", "Password", "Password — Transmittal 579, PDF p.6 (Values column, \"00, 01\"); the name is the guide's: " + tr3},
}

func (v SecQualifierISA03) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v SecQualifierISA03) Known() bool { return secQualifierISA03List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v SecQualifierISA03) ShortDescription() string { return secQualifierISA03List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v SecQualifierISA03) Description() string { return secQualifierISA03List.description(string(v)) }

// ParseSecQualifierISA03 returns text as a SecQualifierISA03, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseSecQualifierISA03(text string) (SecQualifierISA03, error) {
	return SecQualifierISA03(text), secQualifierISA03List.check("ISA03", text)
}

// InterchangeIDQualifierISA05 is the Interchange ID Qualifier, used for
// both ISA05 (sender) and ISA07 (receiver).
//
// The list is every code the 837P implementation guide allows at ISA05 and
// ISA07. Transmittal 579's Values column lists ZZ alone (PDF p.6), and the
// CEDI worked transaction sends it; the other names are the guide's.
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

var interchangeIDQualifierISA05List = codeList{
	{"ZZ", "Mutually Defined", "Mutually Defined — ISA05's and ISA07's only value, Transmittal 579, PDF p.6 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
	{"01", "Duns (Dun & Bradstreet)", guide("Duns (Dun & Bradstreet)", "ISA05 and ISA07")},
	{"14", "Duns Plus Suffix", guide("Duns Plus Suffix", "ISA05 and ISA07")},
	{"20", "Health Industry Number (HIN)", guide("Health Industry Number (HIN)", "ISA05 and ISA07")},
	{"27", "Carrier Identification Number as assigned by Health Care Financing Administration (HCFA)", guide("Carrier Identification Number as assigned by Health Care Financing Administration (HCFA)", "ISA05 and ISA07")},
	{"28", "Fiscal Intermediary Identification Number as assigned by Health Care Financing Administration (HCFA)", guide("Fiscal Intermediary Identification Number as assigned by Health Care Financing Administration (HCFA)", "ISA05 and ISA07")},
	{"29", "Medicare Provider and Supplier Identification Number as assigned by Health Care Financing Administration (HCFA)", guide("Medicare Provider and Supplier Identification Number as assigned by Health Care Financing Administration (HCFA)", "ISA05 and ISA07")},
	{"30", "U.S. Federal Tax Identification Number", guide("U.S. Federal Tax Identification Number", "ISA05 and ISA07")},
	{"33", "National Association of Insurance Commissioners Company Code (NAIC)", guide("National Association of Insurance Commissioners Company Code (NAIC)", "ISA05 and ISA07")},
}

func (v InterchangeIDQualifierISA05) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v InterchangeIDQualifierISA05) Known() bool {
	return interchangeIDQualifierISA05List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v InterchangeIDQualifierISA05) ShortDescription() string {
	return interchangeIDQualifierISA05List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v InterchangeIDQualifierISA05) Description() string {
	return interchangeIDQualifierISA05List.description(string(v))
}

// ParseInterchangeIDQualifierISA05 returns text as an
// InterchangeIDQualifierISA05, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseInterchangeIDQualifierISA05(text string) (InterchangeIDQualifierISA05, error) {
	return InterchangeIDQualifierISA05(text), interchangeIDQualifierISA05List.check("ISA05", text)
}

// AckRequestedISA14 is ISA14, the Acknowledgement Requested code.
type AckRequestedISA14 string

const (
	AckRequestedISA14NotSent AckRequestedISA14 = ""
	AckRequestedISA140       AckRequestedISA14 = "0"
	AckRequestedISA141       AckRequestedISA14 = "1"
)

var ackRequestedISA14List = codeList{
	{"0", "No Interchange Acknowledgment Requested", "No Interchange Acknowledgment Requested — Transmittal 579, PDF p.6 (Values column, \"0, 1\"); the name is the guide's: " + tr3},
	{"1", "Interchange Acknowledgment Requested (TA1)", "Interchange Acknowledgment Requested (TA1) — Transmittal 579, PDF p.6 (Values column, \"0, 1\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
}

func (v AckRequestedISA14) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AckRequestedISA14) Known() bool { return ackRequestedISA14List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AckRequestedISA14) ShortDescription() string { return ackRequestedISA14List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AckRequestedISA14) Description() string { return ackRequestedISA14List.description(string(v)) }

// ParseAckRequestedISA14 returns text as an AckRequestedISA14, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseAckRequestedISA14(text string) (AckRequestedISA14, error) {
	return AckRequestedISA14(text), ackRequestedISA14List.check("ISA14", text)
}

// UsageISA15 is ISA15, the Usage Indicator.
type UsageISA15 string

const (
	UsageISA15NotSent UsageISA15 = ""
	UsageISA15P       UsageISA15 = "P"
	UsageISA15T       UsageISA15 = "T"
)

var usageISA15List = codeList{
	{"P", "Production Data", "Production Data — Transmittal 579, PDF p.6 (837P Values column, \"P, T\"); CMS 835 Flat File, PDF p.1 (Values \"P,T\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
	{"T", "Test Data", "Test Data — Transmittal 579, PDF p.6 (837P Values column, \"P, T\"); CMS 835 Flat File, PDF p.1 (Values \"P,T\"); the name is the guide's: " + tr3},
}

func (v UsageISA15) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v UsageISA15) Known() bool { return usageISA15List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v UsageISA15) ShortDescription() string { return usageISA15List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v UsageISA15) Description() string { return usageISA15List.description(string(v)) }

// ParseUsageISA15 returns text as a UsageISA15, with an error when the text
// is neither empty nor on the list. The value is kept either way.
func ParseUsageISA15(text string) (UsageISA15, error) {
	return UsageISA15(text), usageISA15List.check("ISA15", text)
}
