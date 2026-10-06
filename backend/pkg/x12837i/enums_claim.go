package x12837i

import "fmt"

// AssignmentCLM07 is CLM07, the Assignment or Plan Participation Code.
type AssignmentCLM07 string

const (
	AssignmentCLM07A AssignmentCLM07 = "A"
	AssignmentCLM07B AssignmentCLM07 = "B"
	AssignmentCLM07C AssignmentCLM07 = "C"
)

// The codes the guide allows for CLM07, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var assignmentCLM07Codes = []codeEntry{
	{"A", "Assigned", "Assigned — Transmittal 907, PDF p.23"},
	{"B", "Assignment Accepted on Clinical Lab Services Only", "Assignment Accepted on Clinical Lab Services Only — Transmittal 907, PDF p.23"},
	{"C", "Not Assigned", "Not Assigned — Transmittal 907, PDF p.23"},
}

// String returns the code as the file sent it.
func (v AssignmentCLM07) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v AssignmentCLM07) Known() bool { return knownCode(assignmentCLM07Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v AssignmentCLM07) ShortDescription() string { return shortOf(assignmentCLM07Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v AssignmentCLM07) Description() string { return descriptionOf(assignmentCLM07Codes, string(v)) }

// ParseAssignmentCLM07 reads one CLM07 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseAssignmentCLM07(code string) (AssignmentCLM07, error) {
	v := AssignmentCLM07(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CLM07: a value not on the list")
}

// BenefitsAssignCLM08 is CLM08, the Benefits Assignment Certification
// Indicator. The same code set is reused at OI03 (loop 2320), see
// OtherInsuranceIndicator.
type BenefitsAssignCLM08 string

const (
	BenefitsAssignCLM08N BenefitsAssignCLM08 = "N"
	BenefitsAssignCLM08W BenefitsAssignCLM08 = "W"
	BenefitsAssignCLM08Y BenefitsAssignCLM08 = "Y"
)

// The codes the guide allows for CLM08, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var benefitsAssignCLM08Codes = []codeEntry{
	{"N", "No", "No — Transmittal 907, PDF p.23"},
	{"W", "Not Applicable", "Not Applicable — Transmittal 907, PDF p.23"},
	{"Y", "Yes", "Yes — Transmittal 907, PDF p.23"},
}

// String returns the code as the file sent it.
func (v BenefitsAssignCLM08) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v BenefitsAssignCLM08) Known() bool { return knownCode(benefitsAssignCLM08Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v BenefitsAssignCLM08) ShortDescription() string {
	return shortOf(benefitsAssignCLM08Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v BenefitsAssignCLM08) Description() string {
	return descriptionOf(benefitsAssignCLM08Codes, string(v))
}

// ParseBenefitsAssignCLM08 reads one CLM08 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseBenefitsAssignCLM08(code string) (BenefitsAssignCLM08, error) {
	v := BenefitsAssignCLM08(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CLM08: a value not on the list")
}

// ReleaseOfInfoCLM09 is CLM09, the Release of Information Code. The same
// code set is reused at OI06 (loop 2320), see OtherInsuranceIndicator.
type ReleaseOfInfoCLM09 string

const (
	ReleaseOfInfoCLM09I ReleaseOfInfoCLM09 = "I"
	ReleaseOfInfoCLM09Y ReleaseOfInfoCLM09 = "Y"
)

// The codes the guide allows for CLM09, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var releaseOfInfoCLM09Codes = []codeEntry{
	{"I", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes — Transmittal 907, PDF p.23"},
	{"Y", "Yes, Provider has a Signed Statement Permitting Release of Medical Billing Data Related to a Claim", "Yes, Provider has a Signed Statement Permitting Release of Medical Billing Data Related to a Claim — Transmittal 907, PDF p.23"},
}

// String returns the code as the file sent it.
func (v ReleaseOfInfoCLM09) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v ReleaseOfInfoCLM09) Known() bool { return knownCode(releaseOfInfoCLM09Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v ReleaseOfInfoCLM09) ShortDescription() string {
	return shortOf(releaseOfInfoCLM09Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v ReleaseOfInfoCLM09) Description() string {
	return descriptionOf(releaseOfInfoCLM09Codes, string(v))
}

// ParseReleaseOfInfoCLM09 reads one CLM09 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseReleaseOfInfoCLM09(code string) (ReleaseOfInfoCLM09, error) {
	v := ReleaseOfInfoCLM09(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CLM09: a value not on the list")
}

// DelayReasonCLM20 is CLM20, the Delay Reason Code. Transmittal 907 gives
// only the code list, no per-code text.
type DelayReasonCLM20 string

const (
	DelayReasonCLM201  DelayReasonCLM20 = "1"
	DelayReasonCLM202  DelayReasonCLM20 = "2"
	DelayReasonCLM203  DelayReasonCLM20 = "3"
	DelayReasonCLM204  DelayReasonCLM20 = "4"
	DelayReasonCLM205  DelayReasonCLM20 = "5"
	DelayReasonCLM206  DelayReasonCLM20 = "6"
	DelayReasonCLM207  DelayReasonCLM20 = "7"
	DelayReasonCLM208  DelayReasonCLM20 = "8"
	DelayReasonCLM209  DelayReasonCLM20 = "9"
	DelayReasonCLM2010 DelayReasonCLM20 = "10"
	DelayReasonCLM2011 DelayReasonCLM20 = "11"
	DelayReasonCLM2015 DelayReasonCLM20 = "15"
)

// The codes the guide allows for CLM20, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var delayReasonCLM20Codes = []codeEntry{
	{"1", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"2", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"3", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"4", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"5", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"6", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"7", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"8", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"9", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"10", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"11", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
	{"15", "", "Transmittal 907, PDF p.24; code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v DelayReasonCLM20) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v DelayReasonCLM20) Known() bool { return knownCode(delayReasonCLM20Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v DelayReasonCLM20) ShortDescription() string { return shortOf(delayReasonCLM20Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v DelayReasonCLM20) Description() string {
	return descriptionOf(delayReasonCLM20Codes, string(v))
}

// ParseDelayReasonCLM20 reads one CLM20 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseDelayReasonCLM20(code string) (DelayReasonCLM20, error) {
	v := DelayReasonCLM20(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CLM20: a value not on the list")
}

// PWKReportTypeCode01 is PWK01, the Report Type Code, shared by the
// claim-level PWK (2300) and the line-level PWK (2400). The line-level set
// is a narrower list (Transmittal 907, PDF p.112) almost entirely
// contained in the claim-level set (PDF p.25); A4 is the one code the
// line-level list adds. Both lists give only the code, no per-code text.
type PWKReportTypeCode01 string

const (
	PWKReportTypeCode0103 PWKReportTypeCode01 = "03"
	PWKReportTypeCode0104 PWKReportTypeCode01 = "04"
	PWKReportTypeCode0105 PWKReportTypeCode01 = "05"
	PWKReportTypeCode0106 PWKReportTypeCode01 = "06"
	PWKReportTypeCode0107 PWKReportTypeCode01 = "07"
	PWKReportTypeCode0108 PWKReportTypeCode01 = "08"
	PWKReportTypeCode0109 PWKReportTypeCode01 = "09"
	PWKReportTypeCode0110 PWKReportTypeCode01 = "10"
	PWKReportTypeCode0111 PWKReportTypeCode01 = "11"
	PWKReportTypeCode0113 PWKReportTypeCode01 = "13"
	PWKReportTypeCode0115 PWKReportTypeCode01 = "15"
	PWKReportTypeCode0121 PWKReportTypeCode01 = "21"
	PWKReportTypeCode01A3 PWKReportTypeCode01 = "A3"
	PWKReportTypeCode01AS PWKReportTypeCode01 = "AS"
	PWKReportTypeCode01AM PWKReportTypeCode01 = "AM"
	PWKReportTypeCode01B2 PWKReportTypeCode01 = "B2"
	PWKReportTypeCode01B3 PWKReportTypeCode01 = "B3"
	PWKReportTypeCode01B4 PWKReportTypeCode01 = "B4"
	PWKReportTypeCode01BR PWKReportTypeCode01 = "BR"
	PWKReportTypeCode01BS PWKReportTypeCode01 = "BS"
	PWKReportTypeCode01BT PWKReportTypeCode01 = "BT"
	PWKReportTypeCode01CB PWKReportTypeCode01 = "CB"
	PWKReportTypeCode01CK PWKReportTypeCode01 = "CK"
	PWKReportTypeCode01CT PWKReportTypeCode01 = "CT"
	PWKReportTypeCode01D2 PWKReportTypeCode01 = "D2"
	PWKReportTypeCode01DA PWKReportTypeCode01 = "DA"
	PWKReportTypeCode01DB PWKReportTypeCode01 = "DB"
	PWKReportTypeCode01DG PWKReportTypeCode01 = "DG"
	PWKReportTypeCode01DJ PWKReportTypeCode01 = "DJ"
	PWKReportTypeCode01DS PWKReportTypeCode01 = "DS"
	PWKReportTypeCode01EB PWKReportTypeCode01 = "EB"
	PWKReportTypeCode01HC PWKReportTypeCode01 = "HC"
	PWKReportTypeCode01HR PWKReportTypeCode01 = "HR"
	PWKReportTypeCode01I5 PWKReportTypeCode01 = "I5"
	PWKReportTypeCode01IR PWKReportTypeCode01 = "IR"
	PWKReportTypeCode01LA PWKReportTypeCode01 = "LA"
	PWKReportTypeCode01M1 PWKReportTypeCode01 = "M1"
	PWKReportTypeCode01MT PWKReportTypeCode01 = "MT"
	PWKReportTypeCode01NN PWKReportTypeCode01 = "NN"
	PWKReportTypeCode01OB PWKReportTypeCode01 = "OB"
	PWKReportTypeCode01OC PWKReportTypeCode01 = "OC"
	PWKReportTypeCode01OD PWKReportTypeCode01 = "OD"
	PWKReportTypeCode01OE PWKReportTypeCode01 = "OE"
	PWKReportTypeCode01OX PWKReportTypeCode01 = "OX"
	PWKReportTypeCode01OZ PWKReportTypeCode01 = "OZ"
	PWKReportTypeCode01P4 PWKReportTypeCode01 = "P4"
	PWKReportTypeCode01P5 PWKReportTypeCode01 = "P5"
	PWKReportTypeCode01PE PWKReportTypeCode01 = "PE"
	PWKReportTypeCode01PN PWKReportTypeCode01 = "PN"
	PWKReportTypeCode01PO PWKReportTypeCode01 = "PO"
	PWKReportTypeCode01PQ PWKReportTypeCode01 = "PQ"
	PWKReportTypeCode01PY PWKReportTypeCode01 = "PY"
	PWKReportTypeCode01PZ PWKReportTypeCode01 = "PZ"
	PWKReportTypeCode01RB PWKReportTypeCode01 = "RB"
	PWKReportTypeCode01RR PWKReportTypeCode01 = "RR"
	PWKReportTypeCode01RT PWKReportTypeCode01 = "RT"
	PWKReportTypeCode01RX PWKReportTypeCode01 = "RX"
	PWKReportTypeCode01SG PWKReportTypeCode01 = "SG"
	PWKReportTypeCode01V5 PWKReportTypeCode01 = "V5"
	PWKReportTypeCode01XP PWKReportTypeCode01 = "XP"
	PWKReportTypeCode01A4 PWKReportTypeCode01 = "A4"
)

// The codes the guide allows for PWK01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var pwkReportTypeCode01Codes = []codeEntry{
	{"03", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"04", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"05", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"06", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"07", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"08", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"09", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"10", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"11", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"13", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"15", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"21", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"A3", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"AS", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"AM", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"B2", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"B3", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"B4", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"BR", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"BS", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"BT", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"CB", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"CK", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"CT", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"D2", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"DA", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"DB", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"DG", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"DJ", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"DS", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"EB", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"HC", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"HR", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"I5", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"IR", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"LA", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"M1", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"MT", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"NN", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OB", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OC", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OD", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OE", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OX", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"OZ", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"P4", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"P5", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PE", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PN", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PO", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PQ", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PY", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"PZ", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"RB", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"RR", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"RT", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"RX", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"SG", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"V5", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"XP", "", "Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755); code list only, no per-code text given in the flat file"},
	{"A4", "", "Transmittal 907, PDF p.112 (line-level PWK01); code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v PWKReportTypeCode01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v PWKReportTypeCode01) Known() bool { return knownCode(pwkReportTypeCode01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v PWKReportTypeCode01) ShortDescription() string {
	return shortOf(pwkReportTypeCode01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v PWKReportTypeCode01) Description() string {
	return descriptionOf(pwkReportTypeCode01Codes, string(v))
}

// ParsePWKReportTypeCode01 reads one PWK01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParsePWKReportTypeCode01(code string) (PWKReportTypeCode01, error) {
	v := PWKReportTypeCode01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing PWK01: a value not on the list")
}

// PWKTransmissionCode02 is PWK02, the Report Transmission Code, shared by
// the claim-level PWK (2300) and the line-level PWK (2400). The line-level
// set (SG, V5, XP) is disjoint from the claim-level set; both are merged
// here into one code list.
type PWKTransmissionCode02 string

const (
	PWKTransmissionCode02AA PWKTransmissionCode02 = "AA"
	PWKTransmissionCode02BM PWKTransmissionCode02 = "BM"
	PWKTransmissionCode02EL PWKTransmissionCode02 = "EL"
	PWKTransmissionCode02EM PWKTransmissionCode02 = "EM"
	PWKTransmissionCode02FT PWKTransmissionCode02 = "FT"
	PWKTransmissionCode02FX PWKTransmissionCode02 = "FX"
	PWKTransmissionCode02SG PWKTransmissionCode02 = "SG"
	PWKTransmissionCode02V5 PWKTransmissionCode02 = "V5"
	PWKTransmissionCode02XP PWKTransmissionCode02 = "XP"
)

// The codes the guide allows for PWK02, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var pwkTransmissionCode02Codes = []codeEntry{
	{"AA", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"BM", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"EL", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"EM", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"FT", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"FX", "", "Transmittal 907, PDF p.25 (claim-level PWK02); code list only, no per-code text given in the flat file"},
	{"SG", "", "Transmittal 907, PDF p.112 (line-level PWK02); code list only, no per-code text given in the flat file"},
	{"V5", "", "Transmittal 907, PDF p.112 (line-level PWK02); code list only, no per-code text given in the flat file"},
	{"XP", "", "Transmittal 907, PDF p.112 (line-level PWK02); code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v PWKTransmissionCode02) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v PWKTransmissionCode02) Known() bool { return knownCode(pwkTransmissionCode02Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v PWKTransmissionCode02) ShortDescription() string {
	return shortOf(pwkTransmissionCode02Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v PWKTransmissionCode02) Description() string {
	return descriptionOf(pwkTransmissionCode02Codes, string(v))
}

// ParsePWKTransmissionCode02 reads one PWK02 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParsePWKTransmissionCode02(code string) (PWKTransmissionCode02, error) {
	v := PWKTransmissionCode02(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing PWK02: a value not on the list")
}

// DiagnosisQualifierHI01_1 is HI01-1 (and the equivalent first component of
// every repeated HI0n composite in the transaction): which kind of code
// the composite carries. Reused for every HI segment instance —
// frontend/src/data/x12Names837i.ts's HI_TYPE_NAME table, named once.
// ICD-10 ("AB*"-prefixed) and ICD-9 (unprefixed) qualifiers for the same
// concept are paired codes, not alternate names for one code.
type DiagnosisQualifierHI01_1 string

const (
	DiagnosisQualifierHI01_1ABK DiagnosisQualifierHI01_1 = "ABK"
	DiagnosisQualifierHI01_1BK  DiagnosisQualifierHI01_1 = "BK"
	DiagnosisQualifierHI01_1ABJ DiagnosisQualifierHI01_1 = "ABJ"
	DiagnosisQualifierHI01_1BJ  DiagnosisQualifierHI01_1 = "BJ"
	DiagnosisQualifierHI01_1APR DiagnosisQualifierHI01_1 = "APR"
	DiagnosisQualifierHI01_1PR  DiagnosisQualifierHI01_1 = "PR"
	DiagnosisQualifierHI01_1ABN DiagnosisQualifierHI01_1 = "ABN"
	DiagnosisQualifierHI01_1BN  DiagnosisQualifierHI01_1 = "BN"
	DiagnosisQualifierHI01_1DR  DiagnosisQualifierHI01_1 = "DR"
	DiagnosisQualifierHI01_1ABF DiagnosisQualifierHI01_1 = "ABF"
	DiagnosisQualifierHI01_1BF  DiagnosisQualifierHI01_1 = "BF"
	DiagnosisQualifierHI01_1BBR DiagnosisQualifierHI01_1 = "BBR"
	DiagnosisQualifierHI01_1BR  DiagnosisQualifierHI01_1 = "BR"
	DiagnosisQualifierHI01_1BBQ DiagnosisQualifierHI01_1 = "BBQ"
	DiagnosisQualifierHI01_1BQ  DiagnosisQualifierHI01_1 = "BQ"
	DiagnosisQualifierHI01_1BI  DiagnosisQualifierHI01_1 = "BI"
	DiagnosisQualifierHI01_1BH  DiagnosisQualifierHI01_1 = "BH"
	DiagnosisQualifierHI01_1BE  DiagnosisQualifierHI01_1 = "BE"
	DiagnosisQualifierHI01_1BG  DiagnosisQualifierHI01_1 = "BG"
	DiagnosisQualifierHI01_1TC  DiagnosisQualifierHI01_1 = "TC"
)

// The codes the guide allows for HI0n-1, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var diagnosisQualifierHI01_1Codes = []codeEntry{
	{"ABK", "Principal Diagnosis (ICD-10-CM)", "Principal Diagnosis (ICD-10-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BK", "Principal Diagnosis (ICD-9-CM)", "Principal Diagnosis (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"ABJ", "Admitting Diagnosis (ICD-10-CM)", "Admitting Diagnosis (ICD-10-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BJ", "Admitting Diagnosis (ICD-9-CM)", "Admitting Diagnosis (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"APR", "Patient's Reason for Visit (ICD-10-CM)", "Patient's Reason for Visit (ICD-10-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"PR", "Patient's Reason for Visit (ICD-9-CM)", "Patient's Reason for Visit (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"ABN", "External Cause of Injury (ICD-10-CM)", "External Cause of Injury (ICD-10-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BN", "External Cause of Injury (ICD-9-CM)", "External Cause of Injury (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"DR", "Diagnosis Related Group (DRG) Code", "Diagnosis Related Group (DRG) Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"ABF", "Other Diagnosis (ICD-10-CM)", "Other Diagnosis (ICD-10-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BF", "Other Diagnosis (ICD-9-CM)", "Other Diagnosis (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BBR", "Principal Procedure (ICD-10-PCS)", "Principal Procedure (ICD-10-PCS) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BR", "Principal Procedure (ICD-9-CM)", "Principal Procedure (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BBQ", "Other Procedure (ICD-10-PCS)", "Other Procedure (ICD-10-PCS) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BQ", "Other Procedure (ICD-9-CM)", "Other Procedure (ICD-9-CM) — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BI", "Occurrence Span Code", "Occurrence Span Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BH", "Occurrence Code", "Occurrence Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BE", "Value Code", "Value Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"BG", "Condition Code", "Condition Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
	{"TC", "Treatment Code", "Treatment Code — Transmittal 907, PDF p.30 (ABK/BK); p.34 (ABJ/BJ); p.39 (APR/PR); p.44 (ABN/BN); p.48 (DR); p.51 (ABF/BF); p.58 (BBR/BR); p.62 (BBQ/BQ); p.67 (BI); p.72 (BH); p.77 (BE); p.81 (BG); p.85 (TC)"},
}

// String returns the code as the file sent it.
func (v DiagnosisQualifierHI01_1) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v DiagnosisQualifierHI01_1) Known() bool {
	return knownCode(diagnosisQualifierHI01_1Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v DiagnosisQualifierHI01_1) ShortDescription() string {
	return shortOf(diagnosisQualifierHI01_1Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v DiagnosisQualifierHI01_1) Description() string {
	return descriptionOf(diagnosisQualifierHI01_1Codes, string(v))
}

// ParseDiagnosisQualifierHI01_1 reads one HI0n-1 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseDiagnosisQualifierHI01_1(code string) (DiagnosisQualifierHI01_1, error) {
	v := DiagnosisQualifierHI01_1(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HI0n-1: a value not on the list")
}

// PresentOnAdmissionHI0n_9 is HI0n-9, the Present on Admission Indicator,
// present only on the diagnosis-bearing HI composites (Principal/Other/
// Admitting Diagnosis, External Cause of Injury).
type PresentOnAdmissionHI0n_9 string

const (
	PresentOnAdmissionHI0n_9Y PresentOnAdmissionHI0n_9 = "Y"
	PresentOnAdmissionHI0n_9N PresentOnAdmissionHI0n_9 = "N"
	PresentOnAdmissionHI0n_9U PresentOnAdmissionHI0n_9 = "U"
	PresentOnAdmissionHI0n_9W PresentOnAdmissionHI0n_9 = "W"
)

// The codes the guide allows for HI0n-9, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var presentOnAdmissionHI0n_9Codes = []codeEntry{
	{"Y", "Yes", "Yes — Transmittal 907, PDF p.30"},
	{"N", "No", "No — Transmittal 907, PDF p.30"},
	{"U", "Unknown", "Unknown — Transmittal 907, PDF p.30"},
	{"W", "Clinically Undetermined", "Clinically Undetermined — Transmittal 907, PDF p.30"},
}

// String returns the code as the file sent it.
func (v PresentOnAdmissionHI0n_9) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v PresentOnAdmissionHI0n_9) Known() bool {
	return knownCode(presentOnAdmissionHI0n_9Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v PresentOnAdmissionHI0n_9) ShortDescription() string {
	return shortOf(presentOnAdmissionHI0n_9Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v PresentOnAdmissionHI0n_9) Description() string {
	return descriptionOf(presentOnAdmissionHI0n_9Codes, string(v))
}

// ParsePresentOnAdmissionHI0n_9 reads one HI0n-9 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParsePresentOnAdmissionHI0n_9(code string) (PresentOnAdmissionHI0n_9, error) {
	v := PresentOnAdmissionHI0n_9(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HI0n-9: a value not on the list")
}

// CRCCategoryCRC01 is CRC01, the Code Category on the EPSDT Referral CRC
// segment. Transmittal 907 lists only one value.
type CRCCategoryCRC01 string

const (
	CRCCategoryCRC01ZZ CRCCategoryCRC01 = "ZZ"
)

// The codes the guide allows for CRC01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var crcCategoryCRC01Codes = []codeEntry{
	{"ZZ", "Mutually Defined", "Mutually Defined — Transmittal 907, PDF p.29 (EPSDT Referral code category)"},
}

// String returns the code as the file sent it.
func (v CRCCategoryCRC01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v CRCCategoryCRC01) Known() bool { return knownCode(crcCategoryCRC01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v CRCCategoryCRC01) ShortDescription() string { return shortOf(crcCategoryCRC01Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v CRCCategoryCRC01) Description() string {
	return descriptionOf(crcCategoryCRC01Codes, string(v))
}

// ParseCRCCategoryCRC01 reads one CRC01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseCRCCategoryCRC01(code string) (CRCCategoryCRC01, error) {
	v := CRCCategoryCRC01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CRC01: a value not on the list")
}

// CRCCertificationCRC02 is CRC02, the Certification Condition Code Applies
// Indicator on the EPSDT Referral CRC segment.
type CRCCertificationCRC02 string

const (
	CRCCertificationCRC02N CRCCertificationCRC02 = "N"
	CRCCertificationCRC02Y CRCCertificationCRC02 = "Y"
)

// The codes the guide allows for CRC02, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var crcCertificationCRC02Codes = []codeEntry{
	{"N", "No", "No — Transmittal 907, PDF p.29"},
	{"Y", "Yes", "Yes — Transmittal 907, PDF p.29"},
}

// String returns the code as the file sent it.
func (v CRCCertificationCRC02) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v CRCCertificationCRC02) Known() bool { return knownCode(crcCertificationCRC02Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v CRCCertificationCRC02) ShortDescription() string {
	return shortOf(crcCertificationCRC02Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v CRCCertificationCRC02) Description() string {
	return descriptionOf(crcCertificationCRC02Codes, string(v))
}

// ParseCRCCertificationCRC02 reads one CRC02 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseCRCCertificationCRC02(code string) (CRCCertificationCRC02, error) {
	v := CRCCertificationCRC02(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CRC02: a value not on the list")
}

// CRCConditionIndicator is the Condition Indicator carried at CRC03, CRC04,
// and CRC05 (CRC06-07 are Not Used) — one shared code set reused across
// all three positions, same shape as the HI0n-1 reuse pattern.
type CRCConditionIndicator string

const (
	CRCConditionIndicatorAV CRCConditionIndicator = "AV"
	CRCConditionIndicatorNU CRCConditionIndicator = "NU"
	CRCConditionIndicatorS2 CRCConditionIndicator = "S2"
	CRCConditionIndicatorST CRCConditionIndicator = "ST"
)

// The codes the guide allows for CRC03-05, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var crcConditionIndicatorCodes = []codeEntry{
	{"AV", "", "Transmittal 907, PDF p.29; code list only, no per-code text given in the flat file"},
	{"NU", "", "Transmittal 907, PDF p.29; code list only, no per-code text given in the flat file"},
	{"S2", "", "Transmittal 907, PDF p.29; code list only, no per-code text given in the flat file"},
	{"ST", "", "Transmittal 907, PDF p.29; code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v CRCConditionIndicator) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v CRCConditionIndicator) Known() bool { return knownCode(crcConditionIndicatorCodes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v CRCConditionIndicator) ShortDescription() string {
	return shortOf(crcConditionIndicatorCodes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v CRCConditionIndicator) Description() string {
	return descriptionOf(crcConditionIndicatorCodes, string(v))
}

// ParseCRCConditionIndicator reads one CRC03-05 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseCRCConditionIndicator(code string) (CRCConditionIndicator, error) {
	v := CRCConditionIndicator(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing CRC03-05: a value not on the list")
}

// HCPPricingMethodologyHCP01 is HCP01, the Pricing Methodology, shared by
// the claim-level HCP (2300) and the line-level HCP (2400). Transmittal
// 907 lists the full 00-14 code set but names only 00 in text (as the
// common X12 Pricing Methodology code); the rest are code-named.
type HCPPricingMethodologyHCP01 string

const (
	HCPPricingMethodologyHCP0100 HCPPricingMethodologyHCP01 = "00"
	HCPPricingMethodologyHCP0101 HCPPricingMethodologyHCP01 = "01"
	HCPPricingMethodologyHCP0102 HCPPricingMethodologyHCP01 = "02"
	HCPPricingMethodologyHCP0103 HCPPricingMethodologyHCP01 = "03"
	HCPPricingMethodologyHCP0104 HCPPricingMethodologyHCP01 = "04"
	HCPPricingMethodologyHCP0105 HCPPricingMethodologyHCP01 = "05"
	HCPPricingMethodologyHCP0106 HCPPricingMethodologyHCP01 = "06"
	HCPPricingMethodologyHCP0107 HCPPricingMethodologyHCP01 = "07"
	HCPPricingMethodologyHCP0108 HCPPricingMethodologyHCP01 = "08"
	HCPPricingMethodologyHCP0109 HCPPricingMethodologyHCP01 = "09"
	HCPPricingMethodologyHCP0110 HCPPricingMethodologyHCP01 = "10"
	HCPPricingMethodologyHCP0111 HCPPricingMethodologyHCP01 = "11"
	HCPPricingMethodologyHCP0112 HCPPricingMethodologyHCP01 = "12"
	HCPPricingMethodologyHCP0113 HCPPricingMethodologyHCP01 = "13"
	HCPPricingMethodologyHCP0114 HCPPricingMethodologyHCP01 = "14"
)

// The codes the guide allows for HCP01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var hcpPricingMethodologyHCP01Codes = []codeEntry{
	{"00", "Zero Pricing (Not Covered Under Contract)", "Zero Pricing (Not Covered Under Contract) — common X12 Pricing Methodology code; Transmittal 907, PDF p.91"},
	{"01", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"02", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"03", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"04", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"05", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"06", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"07", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"08", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"09", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"10", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"11", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"12", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"13", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
	{"14", "", "Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)"},
}

// String returns the code as the file sent it.
func (v HCPPricingMethodologyHCP01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v HCPPricingMethodologyHCP01) Known() bool {
	return knownCode(hcpPricingMethodologyHCP01Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v HCPPricingMethodologyHCP01) ShortDescription() string {
	return shortOf(hcpPricingMethodologyHCP01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v HCPPricingMethodologyHCP01) Description() string {
	return descriptionOf(hcpPricingMethodologyHCP01Codes, string(v))
}

// ParseHCPPricingMethodologyHCP01 reads one HCP01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseHCPPricingMethodologyHCP01(code string) (HCPPricingMethodologyHCP01, error) {
	v := HCPPricingMethodologyHCP01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HCP01: a value not on the list")
}

// HCPRejectReasonHCP13 is HCP13, the Reject Reason Code. Transmittal 907
// gives only the code list, no per-code text.
type HCPRejectReasonHCP13 string

const (
	HCPRejectReasonHCP13T1 HCPRejectReasonHCP13 = "T1"
	HCPRejectReasonHCP13T2 HCPRejectReasonHCP13 = "T2"
	HCPRejectReasonHCP13T3 HCPRejectReasonHCP13 = "T3"
	HCPRejectReasonHCP13T4 HCPRejectReasonHCP13 = "T4"
	HCPRejectReasonHCP13T5 HCPRejectReasonHCP13 = "T5"
	HCPRejectReasonHCP13T6 HCPRejectReasonHCP13 = "T6"
)

// The codes the guide allows for HCP13, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var hcpRejectReasonHCP13Codes = []codeEntry{
	{"T1", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"T2", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"T3", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"T4", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"T5", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"T6", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v HCPRejectReasonHCP13) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v HCPRejectReasonHCP13) Known() bool { return knownCode(hcpRejectReasonHCP13Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v HCPRejectReasonHCP13) ShortDescription() string {
	return shortOf(hcpRejectReasonHCP13Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v HCPRejectReasonHCP13) Description() string {
	return descriptionOf(hcpRejectReasonHCP13Codes, string(v))
}

// ParseHCPRejectReasonHCP13 reads one HCP13 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseHCPRejectReasonHCP13(code string) (HCPRejectReasonHCP13, error) {
	v := HCPRejectReasonHCP13(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HCP13: a value not on the list")
}

// HCPPolicyComplianceHCP14 is HCP14, the Policy Compliance Code.
// Transmittal 907 gives only the code list, no per-code text.
type HCPPolicyComplianceHCP14 string

const (
	HCPPolicyComplianceHCP141 HCPPolicyComplianceHCP14 = "1"
	HCPPolicyComplianceHCP142 HCPPolicyComplianceHCP14 = "2"
	HCPPolicyComplianceHCP143 HCPPolicyComplianceHCP14 = "3"
	HCPPolicyComplianceHCP144 HCPPolicyComplianceHCP14 = "4"
	HCPPolicyComplianceHCP145 HCPPolicyComplianceHCP14 = "5"
)

// The codes the guide allows for HCP14, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var hcpPolicyComplianceHCP14Codes = []codeEntry{
	{"1", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"2", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"3", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"4", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"5", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v HCPPolicyComplianceHCP14) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v HCPPolicyComplianceHCP14) Known() bool {
	return knownCode(hcpPolicyComplianceHCP14Codes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v HCPPolicyComplianceHCP14) ShortDescription() string {
	return shortOf(hcpPolicyComplianceHCP14Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v HCPPolicyComplianceHCP14) Description() string {
	return descriptionOf(hcpPolicyComplianceHCP14Codes, string(v))
}

// ParseHCPPolicyComplianceHCP14 reads one HCP14 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseHCPPolicyComplianceHCP14(code string) (HCPPolicyComplianceHCP14, error) {
	v := HCPPolicyComplianceHCP14(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HCP14: a value not on the list")
}

// HCPExceptionCodeHCP15 is HCP15, the Exception Code. Transmittal 907
// gives only the code list, no per-code text.
type HCPExceptionCodeHCP15 string

const (
	HCPExceptionCodeHCP151 HCPExceptionCodeHCP15 = "1"
	HCPExceptionCodeHCP152 HCPExceptionCodeHCP15 = "2"
	HCPExceptionCodeHCP153 HCPExceptionCodeHCP15 = "3"
	HCPExceptionCodeHCP154 HCPExceptionCodeHCP15 = "4"
	HCPExceptionCodeHCP155 HCPExceptionCodeHCP15 = "5"
	HCPExceptionCodeHCP156 HCPExceptionCodeHCP15 = "6"
)

// The codes the guide allows for HCP15, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var hcpExceptionCodeHCP15Codes = []codeEntry{
	{"1", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"2", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"3", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"4", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"5", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
	{"6", "", "Transmittal 907, PDF p.92; code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v HCPExceptionCodeHCP15) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v HCPExceptionCodeHCP15) Known() bool { return knownCode(hcpExceptionCodeHCP15Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v HCPExceptionCodeHCP15) ShortDescription() string {
	return shortOf(hcpExceptionCodeHCP15Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v HCPExceptionCodeHCP15) Description() string {
	return descriptionOf(hcpExceptionCodeHCP15Codes, string(v))
}

// ParseHCPExceptionCodeHCP15 reads one HCP15 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseHCPExceptionCodeHCP15(code string) (HCPExceptionCodeHCP15, error) {
	v := HCPExceptionCodeHCP15(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing HCP15: a value not on the list")
}

// NTENoteCode01 is NTE01, the Note Reference Code, shared by every NTE
// instance in the transaction: Claim Note and Billing Note (both 2300,
// Transmittal 907, PDF p.29) and Third Party Organization Notes (2400).
// The three value sets
// do not overlap; merged here into one code list. Transmittal 907 gives
// only the codes, no per-code text.
type NTENoteCode01 string

const (
	NTENoteCode01ALG NTENoteCode01 = "ALG"
	NTENoteCode01DCP NTENoteCode01 = "DCP"
	NTENoteCode01DGN NTENoteCode01 = "DGN"
	NTENoteCode01DME NTENoteCode01 = "DME"
	NTENoteCode01MED NTENoteCode01 = "MED"
	NTENoteCode01NTR NTENoteCode01 = "NTR"
	NTENoteCode01ODT NTENoteCode01 = "ODT"
	NTENoteCode01RHB NTENoteCode01 = "RHB"
	NTENoteCode01RLH NTENoteCode01 = "RLH"
	NTENoteCode01RNH NTENoteCode01 = "RNH"
	NTENoteCode01SET NTENoteCode01 = "SET"
	NTENoteCode01SFM NTENoteCode01 = "SFM"
	NTENoteCode01SPT NTENoteCode01 = "SPT"
	NTENoteCode01UPI NTENoteCode01 = "UPI"
	NTENoteCode01ADD NTENoteCode01 = "ADD"
	NTENoteCode01TPO NTENoteCode01 = "TPO"
)

// The codes the guide allows for NTE01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var nteNoteCode01Codes = []codeEntry{
	{"ALG", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"DCP", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"DGN", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"DME", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"MED", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"NTR", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"ODT", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"RHB", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"RLH", "Reasons Patient Leaves Home", "Reasons Patient Leaves Home — X12 005010X223A2 TR3, NTE NTE01, loop 2300 Claim Note (transcribed from the standard; verify against the printed guide); not in Transmittal 907's list"},
	{"RNH", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"SET", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"SFM", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"SPT", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"UPI", "", "Transmittal 907, PDF p.29 (Claim Note); code list only, no per-code text given in the flat file"},
	{"ADD", "", "Transmittal 907, PDF p.29 (Billing Note); code list only, no per-code text given in the flat file"},
	{"TPO", "", "Transmittal 907, PDF p.114 (Third Party Organization Notes); code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v NTENoteCode01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v NTENoteCode01) Known() bool { return knownCode(nteNoteCode01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v NTENoteCode01) ShortDescription() string { return shortOf(nteNoteCode01Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v NTENoteCode01) Description() string { return descriptionOf(nteNoteCode01Codes, string(v)) }

// ParseNTENoteCode01 reads one NTE01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseNTENoteCode01(code string) (NTENoteCode01, error) {
	v := NTENoteCode01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing NTE01: a value not on the list")
}
