package x12837p

// Claim-level (2300) coded elements: the CLM composite, PWK, HI, CRC, HCP,
// and the generic Yes/No code set as the claim loop uses it.

// ClaimFrequencyCLM05_3 is CLM05-3, the Claim Frequency Code, the third
// component of the CLM05 composite.
//
// The 837P implementation guide defers this element to an external code
// source, the National Uniform Billing Committee's claim frequency type
// codes, which the repository does not hold; so the list here is Medicare's
// one value, not the guide's. Medicare restricts CLM05-3 to 1 (Original);
// CEDI 837P CG, PDF p.36 (printed p.28; Table 18). Any other frequency code a
// file sends is kept as its text and reports Known false.
type ClaimFrequencyCLM05_3 string

const (
	ClaimFrequencyCLM05_3NotSent ClaimFrequencyCLM05_3 = ""
	ClaimFrequencyCLM05_31       ClaimFrequencyCLM05_3 = "1"
)

var claimFrequencyCLM05_3List = codeList{
	{"1", "Original", "Original — CEDI 837P CG, PDF p.36 (printed p.28; Table 18); frontend/src/data/x12Names837p.ts CLM05_3_NAME"},
}

func (v ClaimFrequencyCLM05_3) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ClaimFrequencyCLM05_3) Known() bool { return claimFrequencyCLM05_3List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ClaimFrequencyCLM05_3) ShortDescription() string {
	return claimFrequencyCLM05_3List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ClaimFrequencyCLM05_3) Description() string {
	return claimFrequencyCLM05_3List.description(string(v))
}

// ParseClaimFrequencyCLM05_3 returns text as a ClaimFrequencyCLM05_3, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseClaimFrequencyCLM05_3(text string) (ClaimFrequencyCLM05_3, error) {
	return ClaimFrequencyCLM05_3(text), claimFrequencyCLM05_3List.check("CLM05-3", text)
}

// FacilityCodeQualifierCLM05_2 is CLM05-2, the Facility Code Qualifier, the
// second component of the CLM05 composite.
//
// The 837P implementation guide allows the one code, B. Transmittal 579's
// Values column lists it (PDF p.27); the name is the guide's.
type FacilityCodeQualifierCLM05_2 string

const (
	FacilityCodeQualifierCLM05_2NotSent FacilityCodeQualifierCLM05_2 = ""
	FacilityCodeQualifierCLM05_2B       FacilityCodeQualifierCLM05_2 = "B"
)

var facilityCodeQualifierCLM05_2List = codeList{
	{"B", "Place of Service Codes for Professional or Dental Services", "Place of Service Codes for Professional or Dental Services — Transmittal 579, PDF p.27 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"12>B>1\"); the name is the guide's: " + tr3},
}

func (v FacilityCodeQualifierCLM05_2) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v FacilityCodeQualifierCLM05_2) Known() bool {
	return facilityCodeQualifierCLM05_2List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v FacilityCodeQualifierCLM05_2) ShortDescription() string {
	return facilityCodeQualifierCLM05_2List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v FacilityCodeQualifierCLM05_2) Description() string {
	return facilityCodeQualifierCLM05_2List.description(string(v))
}

// ParseFacilityCodeQualifierCLM05_2 returns text as a
// FacilityCodeQualifierCLM05_2, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseFacilityCodeQualifierCLM05_2(text string) (FacilityCodeQualifierCLM05_2, error) {
	return FacilityCodeQualifierCLM05_2(text), facilityCodeQualifierCLM05_2List.check("CLM05-2", text)
}

// ProviderOrSupplierSignatureCLM06 is CLM06, the Provider or Supplier Signature
// Indicator.
//
// The 837P implementation guide allows N and Y. Transmittal 579's Values
// column lists "N, Y" (PDF p.27); the names are the guide's.
type ProviderOrSupplierSignatureCLM06 string

const (
	ProviderOrSupplierSignatureCLM06NotSent ProviderOrSupplierSignatureCLM06 = ""
	ProviderOrSupplierSignatureCLM06N       ProviderOrSupplierSignatureCLM06 = "N"
	ProviderOrSupplierSignatureCLM06Y       ProviderOrSupplierSignatureCLM06 = "Y"
)

var providerOrSupplierSignatureCLM06List = codeList{
	{"N", "No", "No — Transmittal 579, PDF p.27 (Values column, \"N, Y\"); the name is the guide's: " + tr3},
	{"Y", "Yes", "Yes — Transmittal 579, PDF p.27 (Values column, \"N, Y\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
}

func (v ProviderOrSupplierSignatureCLM06) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ProviderOrSupplierSignatureCLM06) Known() bool {
	return providerOrSupplierSignatureCLM06List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ProviderOrSupplierSignatureCLM06) ShortDescription() string {
	return providerOrSupplierSignatureCLM06List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ProviderOrSupplierSignatureCLM06) Description() string {
	return providerOrSupplierSignatureCLM06List.description(string(v))
}

// ParseProviderOrSupplierSignatureCLM06 returns text as a
// ProviderOrSupplierSignatureCLM06, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseProviderOrSupplierSignatureCLM06(text string) (ProviderOrSupplierSignatureCLM06, error) {
	return ProviderOrSupplierSignatureCLM06(text), providerOrSupplierSignatureCLM06List.check("CLM06", text)
}

// AssignmentOrPlanParticipationCLM07 is CLM07, the Assignment or Plan
// Participation Code.
//
// The 837P implementation guide allows A, B, and C. Transmittal 579's Values
// column lists "A, B, C" (PDF p.27); the names are the guide's.
type AssignmentOrPlanParticipationCLM07 string

const (
	AssignmentOrPlanParticipationCLM07NotSent AssignmentOrPlanParticipationCLM07 = ""
	AssignmentOrPlanParticipationCLM07A       AssignmentOrPlanParticipationCLM07 = "A"
	AssignmentOrPlanParticipationCLM07B       AssignmentOrPlanParticipationCLM07 = "B"
	AssignmentOrPlanParticipationCLM07C       AssignmentOrPlanParticipationCLM07 = "C"
)

var assignmentOrPlanParticipationCLM07List = codeList{
	{"A", "Assigned", "Assigned — Transmittal 579, PDF p.27 (Values column, \"A, B, C\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
	{"B", "Assignment Accepted on Clinical Lab Services Only", "Assignment Accepted on Clinical Lab Services Only — Transmittal 579, PDF p.27 (Values column, \"A, B, C\"); the name is the guide's: " + tr3},
	{"C", "Not Assigned", "Not Assigned — Transmittal 579, PDF p.27 (Values column, \"A, B, C\"); the name is the guide's: " + tr3},
}

func (v AssignmentOrPlanParticipationCLM07) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AssignmentOrPlanParticipationCLM07) Known() bool {
	return assignmentOrPlanParticipationCLM07List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AssignmentOrPlanParticipationCLM07) ShortDescription() string {
	return assignmentOrPlanParticipationCLM07List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AssignmentOrPlanParticipationCLM07) Description() string {
	return assignmentOrPlanParticipationCLM07List.description(string(v))
}

// ParseAssignmentOrPlanParticipationCLM07 returns text as an
// AssignmentOrPlanParticipationCLM07, with an error when the text is
// neither empty nor on the list. The value is kept either way.
func ParseAssignmentOrPlanParticipationCLM07(text string) (AssignmentOrPlanParticipationCLM07, error) {
	return AssignmentOrPlanParticipationCLM07(text), assignmentOrPlanParticipationCLM07List.check("CLM07", text)
}

// BenefitsAssignmentCertificationCLM08 is CLM08, the Benefits Assignment
// Certification Indicator; OI03 at 2320 shares the type.
//
// The 837P implementation guide allows N, W, and Y at both positions.
// Transmittal 579's Values column lists "N, W, Y" for CLM08 and OI03 (PDF
// p.27, 67); the names are the guide's.
type BenefitsAssignmentCertificationCLM08 string

const (
	BenefitsAssignmentCertificationCLM08NotSent BenefitsAssignmentCertificationCLM08 = ""
	BenefitsAssignmentCertificationCLM08N       BenefitsAssignmentCertificationCLM08 = "N"
	BenefitsAssignmentCertificationCLM08W       BenefitsAssignmentCertificationCLM08 = "W"
	BenefitsAssignmentCertificationCLM08Y       BenefitsAssignmentCertificationCLM08 = "Y"
)

var benefitsAssignmentCertificationCLM08List = codeList{
	{"N", "No", "No — Transmittal 579, PDF p.27 and p.67 (Values column, \"N, W, Y\"); the name is the guide's: " + tr3},
	{"W", "Not Applicable", "Not Applicable — Transmittal 579, PDF p.27 and p.67 (Values column, \"N, W, Y\"); the name is the guide's: " + tr3},
	{"Y", "Yes", "Yes — Transmittal 579, PDF p.27 and p.67 (Values column, \"N, W, Y\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
}

func (v BenefitsAssignmentCertificationCLM08) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v BenefitsAssignmentCertificationCLM08) Known() bool {
	return benefitsAssignmentCertificationCLM08List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v BenefitsAssignmentCertificationCLM08) ShortDescription() string {
	return benefitsAssignmentCertificationCLM08List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v BenefitsAssignmentCertificationCLM08) Description() string {
	return benefitsAssignmentCertificationCLM08List.description(string(v))
}

// ParseBenefitsAssignmentCertificationCLM08 returns text as a
// BenefitsAssignmentCertificationCLM08, with an error when the text is
// neither empty nor on the list. The value is kept either way.
func ParseBenefitsAssignmentCertificationCLM08(text string) (BenefitsAssignmentCertificationCLM08, error) {
	return BenefitsAssignmentCertificationCLM08(text), benefitsAssignmentCertificationCLM08List.check("CLM08", text)
}

// ReleaseOfInformationCLM09 is CLM09, the Release of Information Code; OI06 at
// 2320 shares the type.
//
// The 837P implementation guide allows I and Y at both positions.
// Transmittal 579's Values column lists "I, Y" for CLM09 and OI06 (PDF p.27,
// 67); the names are the guide's.
type ReleaseOfInformationCLM09 string

const (
	ReleaseOfInformationCLM09NotSent ReleaseOfInformationCLM09 = ""
	ReleaseOfInformationCLM09I       ReleaseOfInformationCLM09 = "I"
	ReleaseOfInformationCLM09Y       ReleaseOfInformationCLM09 = "Y"
)

var releaseOfInformationCLM09List = codeList{
	{"I", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes", "Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes — Transmittal 579, PDF p.27 and p.67 (Values column, \"I, Y\"); the name is the guide's: " + tr3},
	{"Y", "Yes, Provider has a Signed Statement Permitting Release of Medical Billing Data Related to a Claim", "Yes, Provider has a Signed Statement Permitting Release of Medical Billing Data Related to a Claim — Transmittal 579, PDF p.27 and p.67 (Values column, \"I, Y\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3); the name is the guide's: " + tr3},
}

func (v ReleaseOfInformationCLM09) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ReleaseOfInformationCLM09) Known() bool {
	return releaseOfInformationCLM09List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ReleaseOfInformationCLM09) ShortDescription() string {
	return releaseOfInformationCLM09List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ReleaseOfInformationCLM09) Description() string {
	return releaseOfInformationCLM09List.description(string(v))
}

// ParseReleaseOfInformationCLM09 returns text as a ReleaseOfInformationCLM09,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseReleaseOfInformationCLM09(text string) (ReleaseOfInformationCLM09, error) {
	return ReleaseOfInformationCLM09(text), releaseOfInformationCLM09List.check("CLM09", text)
}

// PatientSignatureSourceCLM10 is CLM10, the Patient Signature Source Code;
// OI04 at 2320 shares the type.
//
// The 837P implementation guide allows the one code, P, at both positions.
// Transmittal 579's Values column lists it for CLM10 and OI04 (PDF p.27, 67);
// the name is the guide's.
type PatientSignatureSourceCLM10 string

const (
	PatientSignatureSourceCLM10NotSent PatientSignatureSourceCLM10 = ""
	PatientSignatureSourceCLM10P       PatientSignatureSourceCLM10 = "P"
)

var patientSignatureSourceCLM10List = codeList{
	{"P", "Signature generated by provider because the patient was not physically present for services", "Signature generated by provider because the patient was not physically present for services — Transmittal 579, PDF p.27 and p.67 (Values column); the name is the guide's: " + tr3},
}

func (v PatientSignatureSourceCLM10) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v PatientSignatureSourceCLM10) Known() bool {
	return patientSignatureSourceCLM10List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v PatientSignatureSourceCLM10) ShortDescription() string {
	return patientSignatureSourceCLM10List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v PatientSignatureSourceCLM10) Description() string {
	return patientSignatureSourceCLM10List.description(string(v))
}

// ParsePatientSignatureSourceCLM10 returns text as a
// PatientSignatureSourceCLM10, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParsePatientSignatureSourceCLM10(text string) (PatientSignatureSourceCLM10, error) {
	return PatientSignatureSourceCLM10(text), patientSignatureSourceCLM10List.check("CLM10", text)
}

// RelatedCausesCodeCLM11_1 is CLM11-1, the Related Causes Code, the first
// component of the CLM11 composite (CLM11-2 repeats the same code set in its
// second component).
//
// The 837P implementation guide allows AA, EM, and OA. Transmittal 579's
// Values column lists "AA, EM, OA" (PDF p.27); the names are the guide's.
type RelatedCausesCodeCLM11_1 string

const (
	RelatedCausesCodeCLM11_1NotSent RelatedCausesCodeCLM11_1 = ""
	RelatedCausesCodeCLM11_1AA      RelatedCausesCodeCLM11_1 = "AA"
	RelatedCausesCodeCLM11_1EM      RelatedCausesCodeCLM11_1 = "EM"
	RelatedCausesCodeCLM11_1OA      RelatedCausesCodeCLM11_1 = "OA"
)

var relatedCausesCodeCLM11_1List = codeList{
	{"AA", "Auto Accident", "Auto Accident — Transmittal 579, PDF p.27 (Values column, \"AA, EM, OA\"); the name is the guide's: " + tr3},
	{"EM", "Employment", "Employment — Transmittal 579, PDF p.27 (Values column, \"AA, EM, OA\"); the name is the guide's: " + tr3},
	{"OA", "Other Accident", "Other Accident — Transmittal 579, PDF p.27 (Values column, \"AA, EM, OA\"); the name is the guide's: " + tr3},
}

func (v RelatedCausesCodeCLM11_1) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v RelatedCausesCodeCLM11_1) Known() bool { return relatedCausesCodeCLM11_1List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v RelatedCausesCodeCLM11_1) ShortDescription() string {
	return relatedCausesCodeCLM11_1List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v RelatedCausesCodeCLM11_1) Description() string {
	return relatedCausesCodeCLM11_1List.description(string(v))
}

// ParseRelatedCausesCodeCLM11_1 returns text as a RelatedCausesCodeCLM11_1,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseRelatedCausesCodeCLM11_1(text string) (RelatedCausesCodeCLM11_1, error) {
	return RelatedCausesCodeCLM11_1(text), relatedCausesCodeCLM11_1List.check("CLM11-1", text)
}

// DelayReasonCLM20 is CLM20, the Delay Reason Code.
//
// The list is every code the 837P implementation guide allows at CLM20.
// Medicare ignores any value submitted here (CEDI 837P CG, PDF p.36, printed
// p.28, Table 18); the names are the guide's.
type DelayReasonCLM20 string

const (
	DelayReasonCLM20NotSent DelayReasonCLM20 = ""
	DelayReasonCLM201       DelayReasonCLM20 = "1"
	DelayReasonCLM202       DelayReasonCLM20 = "2"
	DelayReasonCLM203       DelayReasonCLM20 = "3"
	DelayReasonCLM204       DelayReasonCLM20 = "4"
	DelayReasonCLM205       DelayReasonCLM20 = "5"
	DelayReasonCLM206       DelayReasonCLM20 = "6"
	DelayReasonCLM207       DelayReasonCLM20 = "7"
	DelayReasonCLM208       DelayReasonCLM20 = "8"
	DelayReasonCLM209       DelayReasonCLM20 = "9"
	DelayReasonCLM2010      DelayReasonCLM20 = "10"
	DelayReasonCLM2011      DelayReasonCLM20 = "11"
	DelayReasonCLM2015      DelayReasonCLM20 = "15"
)

var delayReasonCLM20List = codeList{
	{"1", "Proof of Eligibility Unknown or Unavailable", guide("Proof of Eligibility Unknown or Unavailable", "CLM20")},
	{"2", "Litigation", guide("Litigation", "CLM20")},
	{"3", "Authorization Delays", guide("Authorization Delays", "CLM20")},
	{"4", "Delay in Certifying Provider", guide("Delay in Certifying Provider", "CLM20")},
	{"5", "Delay in Supplying Billing Forms", guide("Delay in Supplying Billing Forms", "CLM20")},
	{"6", "Delay in Delivery of Custom-made Appliances", guide("Delay in Delivery of Custom-made Appliances", "CLM20")},
	{"7", "Third Party Processing Delay", guide("Third Party Processing Delay", "CLM20")},
	{"8", "Delay in Eligibility Determination", guide("Delay in Eligibility Determination", "CLM20")},
	{"9", "Original Claim Rejected or Denied Due to a Reason Unrelated to the Billing Limitation Rules", guide("Original Claim Rejected or Denied Due to a Reason Unrelated to the Billing Limitation Rules", "CLM20")},
	{"10", "Administration Delay in the Prior Approval Process", guide("Administration Delay in the Prior Approval Process", "CLM20")},
	{"11", "Other", guide("Other", "CLM20")},
	{"15", "Natural Disaster", guide("Natural Disaster", "CLM20")},
}

func (v DelayReasonCLM20) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v DelayReasonCLM20) Known() bool { return delayReasonCLM20List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v DelayReasonCLM20) ShortDescription() string { return delayReasonCLM20List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v DelayReasonCLM20) Description() string { return delayReasonCLM20List.description(string(v)) }

// ParseDelayReasonCLM20 returns text as a DelayReasonCLM20, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseDelayReasonCLM20(text string) (DelayReasonCLM20, error) {
	return DelayReasonCLM20(text), delayReasonCLM20List.check("CLM20", text)
}

// AttachmentReportTypeCodePWK01 is PWK01, the Attachment Report Type Code, on a
// PWK segment (2300 Claim Supplemental Information, 2400 Line Supplemental
// Information, and the 2400 DME Certificate of Medical Necessity box).
//
// The list is every code the 837P implementation guide allows at PWK01, with
// the guide's names. Neither Transmittal 907 nor the CEDI 837P Companion
// Guide names PWK01's individual codes in text; frontend/src/data/x12Names837p.ts
// PWK_EL_NAME gives only the element's own position name. Transmittal 579
// lists CT at the DME Certificate of Medical Necessity box (PDF p.82).
type AttachmentReportTypeCodePWK01 string

const (
	AttachmentReportTypeCodePWK01NotSent AttachmentReportTypeCodePWK01 = ""
	AttachmentReportTypeCodePWK0103      AttachmentReportTypeCodePWK01 = "03"
	AttachmentReportTypeCodePWK0104      AttachmentReportTypeCodePWK01 = "04"
	AttachmentReportTypeCodePWK0105      AttachmentReportTypeCodePWK01 = "05"
	AttachmentReportTypeCodePWK0106      AttachmentReportTypeCodePWK01 = "06"
	AttachmentReportTypeCodePWK0107      AttachmentReportTypeCodePWK01 = "07"
	AttachmentReportTypeCodePWK0108      AttachmentReportTypeCodePWK01 = "08"
	AttachmentReportTypeCodePWK0109      AttachmentReportTypeCodePWK01 = "09"
	AttachmentReportTypeCodePWK0110      AttachmentReportTypeCodePWK01 = "10"
	AttachmentReportTypeCodePWK0111      AttachmentReportTypeCodePWK01 = "11"
	AttachmentReportTypeCodePWK0113      AttachmentReportTypeCodePWK01 = "13"
	AttachmentReportTypeCodePWK0115      AttachmentReportTypeCodePWK01 = "15"
	AttachmentReportTypeCodePWK0121      AttachmentReportTypeCodePWK01 = "21"
	AttachmentReportTypeCodePWK01A3      AttachmentReportTypeCodePWK01 = "A3"
	AttachmentReportTypeCodePWK01A4      AttachmentReportTypeCodePWK01 = "A4"
	AttachmentReportTypeCodePWK01AM      AttachmentReportTypeCodePWK01 = "AM"
	AttachmentReportTypeCodePWK01AS      AttachmentReportTypeCodePWK01 = "AS"
	AttachmentReportTypeCodePWK01B2      AttachmentReportTypeCodePWK01 = "B2"
	AttachmentReportTypeCodePWK01B3      AttachmentReportTypeCodePWK01 = "B3"
	AttachmentReportTypeCodePWK01B4      AttachmentReportTypeCodePWK01 = "B4"
	AttachmentReportTypeCodePWK01BR      AttachmentReportTypeCodePWK01 = "BR"
	AttachmentReportTypeCodePWK01BS      AttachmentReportTypeCodePWK01 = "BS"
	AttachmentReportTypeCodePWK01BT      AttachmentReportTypeCodePWK01 = "BT"
	AttachmentReportTypeCodePWK01CB      AttachmentReportTypeCodePWK01 = "CB"
	AttachmentReportTypeCodePWK01CK      AttachmentReportTypeCodePWK01 = "CK"
	AttachmentReportTypeCodePWK01CT      AttachmentReportTypeCodePWK01 = "CT"
	AttachmentReportTypeCodePWK01D2      AttachmentReportTypeCodePWK01 = "D2"
	AttachmentReportTypeCodePWK01DA      AttachmentReportTypeCodePWK01 = "DA"
	AttachmentReportTypeCodePWK01DB      AttachmentReportTypeCodePWK01 = "DB"
	AttachmentReportTypeCodePWK01DG      AttachmentReportTypeCodePWK01 = "DG"
	AttachmentReportTypeCodePWK01DJ      AttachmentReportTypeCodePWK01 = "DJ"
	AttachmentReportTypeCodePWK01DS      AttachmentReportTypeCodePWK01 = "DS"
	AttachmentReportTypeCodePWK01EB      AttachmentReportTypeCodePWK01 = "EB"
	AttachmentReportTypeCodePWK01HC      AttachmentReportTypeCodePWK01 = "HC"
	AttachmentReportTypeCodePWK01HR      AttachmentReportTypeCodePWK01 = "HR"
	AttachmentReportTypeCodePWK01I5      AttachmentReportTypeCodePWK01 = "I5"
	AttachmentReportTypeCodePWK01IR      AttachmentReportTypeCodePWK01 = "IR"
	AttachmentReportTypeCodePWK01LA      AttachmentReportTypeCodePWK01 = "LA"
	AttachmentReportTypeCodePWK01M1      AttachmentReportTypeCodePWK01 = "M1"
	AttachmentReportTypeCodePWK01MT      AttachmentReportTypeCodePWK01 = "MT"
	AttachmentReportTypeCodePWK01NN      AttachmentReportTypeCodePWK01 = "NN"
	AttachmentReportTypeCodePWK01OB      AttachmentReportTypeCodePWK01 = "OB"
	AttachmentReportTypeCodePWK01OC      AttachmentReportTypeCodePWK01 = "OC"
	AttachmentReportTypeCodePWK01OD      AttachmentReportTypeCodePWK01 = "OD"
	AttachmentReportTypeCodePWK01OE      AttachmentReportTypeCodePWK01 = "OE"
	AttachmentReportTypeCodePWK01OX      AttachmentReportTypeCodePWK01 = "OX"
	AttachmentReportTypeCodePWK01OZ      AttachmentReportTypeCodePWK01 = "OZ"
	AttachmentReportTypeCodePWK01P4      AttachmentReportTypeCodePWK01 = "P4"
	AttachmentReportTypeCodePWK01P5      AttachmentReportTypeCodePWK01 = "P5"
	AttachmentReportTypeCodePWK01PE      AttachmentReportTypeCodePWK01 = "PE"
	AttachmentReportTypeCodePWK01PN      AttachmentReportTypeCodePWK01 = "PN"
	AttachmentReportTypeCodePWK01PO      AttachmentReportTypeCodePWK01 = "PO"
	AttachmentReportTypeCodePWK01PQ      AttachmentReportTypeCodePWK01 = "PQ"
	AttachmentReportTypeCodePWK01PY      AttachmentReportTypeCodePWK01 = "PY"
	AttachmentReportTypeCodePWK01PZ      AttachmentReportTypeCodePWK01 = "PZ"
	AttachmentReportTypeCodePWK01RB      AttachmentReportTypeCodePWK01 = "RB"
	AttachmentReportTypeCodePWK01RR      AttachmentReportTypeCodePWK01 = "RR"
	AttachmentReportTypeCodePWK01RT      AttachmentReportTypeCodePWK01 = "RT"
	AttachmentReportTypeCodePWK01RX      AttachmentReportTypeCodePWK01 = "RX"
	AttachmentReportTypeCodePWK01SG      AttachmentReportTypeCodePWK01 = "SG"
	AttachmentReportTypeCodePWK01V5      AttachmentReportTypeCodePWK01 = "V5"
	AttachmentReportTypeCodePWK01XP      AttachmentReportTypeCodePWK01 = "XP"
)

var attachmentReportTypeCodePWK01List = codeList{
	{"03", "Report Justifying Treatment Beyond Utilization Guidelines", guide("Report Justifying Treatment Beyond Utilization Guidelines", "PWK01")},
	{"04", "Drugs Administered", guide("Drugs Administered", "PWK01")},
	{"05", "Treatment Diagnosis", guide("Treatment Diagnosis", "PWK01")},
	{"06", "Initial Assessment", guide("Initial Assessment", "PWK01")},
	{"07", "Functional Goals", guide("Functional Goals", "PWK01")},
	{"08", "Plan of Treatment", guide("Plan of Treatment", "PWK01")},
	{"09", "Progress Report", guide("Progress Report", "PWK01")},
	{"10", "Continued Treatment", guide("Continued Treatment", "PWK01")},
	{"11", "Chemical Analysis", guide("Chemical Analysis", "PWK01")},
	{"13", "Certified Test Report", guide("Certified Test Report", "PWK01")},
	{"15", "Justification for Admission", guide("Justification for Admission", "PWK01")},
	{"21", "Recovery Plan", guide("Recovery Plan", "PWK01")},
	{"A3", "Allergies/Sensitivities Document", guide("Allergies/Sensitivities Document", "PWK01")},
	{"A4", "Autopsy Report", guide("Autopsy Report", "PWK01")},
	{"AM", "Ambulance Certification", guide("Ambulance Certification", "PWK01")},
	{"AS", "Admission Summary", guide("Admission Summary", "PWK01")},
	{"B2", "Prescription", guide("Prescription", "PWK01")},
	{"B3", "Physician Order", guide("Physician Order", "PWK01")},
	{"B4", "Referral Form", guide("Referral Form", "PWK01")},
	{"BR", "Benchmark Testing Results", guide("Benchmark Testing Results", "PWK01")},
	{"BS", "Baseline", guide("Baseline", "PWK01")},
	{"BT", "Blanket Test Results", guide("Blanket Test Results", "PWK01")},
	{"CB", "Chiropractic Justification", guide("Chiropractic Justification", "PWK01")},
	{"CK", "Consent Form(s)", guide("Consent Form(s)", "PWK01")},
	{"CT", "Certification", "Certification — the 2400 Durable Medical Equipment Certificate of Medical Necessity Indicator box's PWK01 value, Transmittal 579, PDF p.82; the name is the guide's: " + tr3},
	{"D2", "Drug Profile Document", guide("Drug Profile Document", "PWK01")},
	{"DA", "Dental Models", guide("Dental Models", "PWK01")},
	{"DB", "Durable Medical Equipment Prescription", guide("Durable Medical Equipment Prescription", "PWK01")},
	{"DG", "Diagnostic Report", guide("Diagnostic Report", "PWK01")},
	{"DJ", "Discharge Monitoring Report", guide("Discharge Monitoring Report", "PWK01")},
	{"DS", "Discharge Summary", guide("Discharge Summary", "PWK01")},
	{"EB", "Explanation of Benefits (Coordination of Benefits or Medicare Secondary Payor)", guide("Explanation of Benefits (Coordination of Benefits or Medicare Secondary Payor)", "PWK01")},
	{"HC", "Health Certificate", guide("Health Certificate", "PWK01")},
	{"HR", "Health Clinic Records", guide("Health Clinic Records", "PWK01")},
	{"I5", "Immunization Record", guide("Immunization Record", "PWK01")},
	{"IR", "State School Immunization Records", guide("State School Immunization Records", "PWK01")},
	{"LA", "Laboratory Results", guide("Laboratory Results", "PWK01")},
	{"M1", "Medical Record Attachment", guide("Medical Record Attachment", "PWK01")},
	{"MT", "Models", guide("Models", "PWK01")},
	{"NN", "Nursing Notes", guide("Nursing Notes", "PWK01")},
	{"OB", "Operative Note", guide("Operative Note", "PWK01")},
	{"OC", "Oxygen Content Averaging Report", guide("Oxygen Content Averaging Report", "PWK01")},
	{"OD", "Orders and Treatments Document", guide("Orders and Treatments Document", "PWK01")},
	{"OE", "Objective Physical Examination (including vital signs) Document", guide("Objective Physical Examination (including vital signs) Document", "PWK01")},
	{"OX", "Oxygen Therapy Certification", guide("Oxygen Therapy Certification", "PWK01")},
	{"OZ", "Support Data for Claim", guide("Support Data for Claim", "PWK01")},
	{"P4", "Pathology Report", guide("Pathology Report", "PWK01")},
	{"P5", "Patient Medical History Document", guide("Patient Medical History Document", "PWK01")},
	{"PE", "Parenteral or Enteral Certification", guide("Parenteral or Enteral Certification", "PWK01")},
	{"PN", "Physical Therapy Notes", guide("Physical Therapy Notes", "PWK01")},
	{"PO", "Prosthetics or Orthotic Certification", guide("Prosthetics or Orthotic Certification", "PWK01")},
	{"PQ", "Paramedical Results", guide("Paramedical Results", "PWK01")},
	{"PY", "Physician's Report", guide("Physician's Report", "PWK01")},
	{"PZ", "Physical Therapy Certification", guide("Physical Therapy Certification", "PWK01")},
	{"RB", "Radiology Films", guide("Radiology Films", "PWK01")},
	{"RR", "Radiology Reports", guide("Radiology Reports", "PWK01")},
	{"RT", "Report of Tests and Analysis Report", guide("Report of Tests and Analysis Report", "PWK01")},
	{"RX", "Renewable Oxygen Content Averaging Report", guide("Renewable Oxygen Content Averaging Report", "PWK01")},
	{"SG", "Symptoms Document", guide("Symptoms Document", "PWK01")},
	{"V5", "Death Notification", guide("Death Notification", "PWK01")},
	{"XP", "Photographs", guide("Photographs", "PWK01")},
}

func (v AttachmentReportTypeCodePWK01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AttachmentReportTypeCodePWK01) Known() bool {
	return attachmentReportTypeCodePWK01List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AttachmentReportTypeCodePWK01) ShortDescription() string {
	return attachmentReportTypeCodePWK01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AttachmentReportTypeCodePWK01) Description() string {
	return attachmentReportTypeCodePWK01List.description(string(v))
}

// ParseAttachmentReportTypeCodePWK01 returns text as an
// AttachmentReportTypeCodePWK01, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseAttachmentReportTypeCodePWK01(text string) (AttachmentReportTypeCodePWK01, error) {
	return AttachmentReportTypeCodePWK01(text), attachmentReportTypeCodePWK01List.check("PWK01", text)
}

// AttachmentTransmissionCodePWK02 is PWK02, the Attachment Transmission Code,
// on a PWK segment (2300 Claim Supplemental Information, 2400 Line Supplemental
// Information).
//
// The list is every code the 837P implementation guide allows at PWK02, with
// the guide's names. Medicare restricts PWK02 to BM, FX, EL, and FT; CEDI
// 837P CG, PDF p.38 (printed p.30; Table 21); only the first PWK iteration
// at either the claim or line level is used in adjudication.
type AttachmentTransmissionCodePWK02 string

const (
	AttachmentTransmissionCodePWK02NotSent AttachmentTransmissionCodePWK02 = ""
	AttachmentTransmissionCodePWK02BM      AttachmentTransmissionCodePWK02 = "BM"
	AttachmentTransmissionCodePWK02FX      AttachmentTransmissionCodePWK02 = "FX"
	AttachmentTransmissionCodePWK02EL      AttachmentTransmissionCodePWK02 = "EL"
	AttachmentTransmissionCodePWK02FT      AttachmentTransmissionCodePWK02 = "FT"
	AttachmentTransmissionCodePWK02AA      AttachmentTransmissionCodePWK02 = "AA"
	AttachmentTransmissionCodePWK02EM      AttachmentTransmissionCodePWK02 = "EM"
)

var attachmentTransmissionCodePWK02List = codeList{
	{"BM", "By Mail", "By Mail — CEDI 837P CG, PDF p.38 (printed p.30; Table 21), one of the four codes Medicare accepts; frontend/src/data/x12Names837p.ts PWK02_MEDICARE_CODES; the name is the guide's: " + tr3},
	{"FX", "By Fax", "By Fax — CEDI 837P CG, PDF p.38 (printed p.30; Table 21), one of the four codes Medicare accepts; frontend/src/data/x12Names837p.ts PWK02_MEDICARE_CODES; the name is the guide's: " + tr3},
	{"EL", "Electronically Only", "Electronically Only — CEDI 837P CG, PDF p.38 (printed p.30; Table 21), one of the four codes Medicare accepts; frontend/src/data/x12Names837p.ts PWK02_MEDICARE_CODES; the name is the guide's: " + tr3},
	{"FT", "File Transfer", "File Transfer — CEDI 837P CG, PDF p.38 (printed p.30; Table 21), one of the four codes Medicare accepts; frontend/src/data/x12Names837p.ts PWK02_MEDICARE_CODES; the name is the guide's: " + tr3},
	{"AA", "Available on Request at Provider Site", guide("Available on Request at Provider Site", "PWK02") + "; not among the four codes Medicare accepts (CEDI 837P CG, PDF p.38)"},
	{"EM", "E-Mail", guide("E-Mail", "PWK02") + "; not among the four codes Medicare accepts (CEDI 837P CG, PDF p.38)"},
}

func (v AttachmentTransmissionCodePWK02) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AttachmentTransmissionCodePWK02) Known() bool {
	return attachmentTransmissionCodePWK02List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AttachmentTransmissionCodePWK02) ShortDescription() string {
	return attachmentTransmissionCodePWK02List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AttachmentTransmissionCodePWK02) Description() string {
	return attachmentTransmissionCodePWK02List.description(string(v))
}

// ParseAttachmentTransmissionCodePWK02 returns text as an
// AttachmentTransmissionCodePWK02, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseAttachmentTransmissionCodePWK02(text string) (AttachmentTransmissionCodePWK02, error) {
	return AttachmentTransmissionCodePWK02(text), attachmentTransmissionCodePWK02List.check("PWK02", text)
}

// DiagnosisTypeCodeHI0x_1 is HI0n-1 (HI01-1 through HI12-1), the Diagnosis Type
// Code qualifying each of the up to twelve HI composites in the 2300 Health
// Care Diagnosis Code segment.
//
// The same element (X12 1270) qualifies the composites of the 2300
// Anesthesia Related Procedure and Condition Information HI segments, where
// Transmittal 907 calls it the Code List Qualifier Code, so this one type
// serves all three HI boxes. The list is every code the 837P implementation
// guide allows across the three: ABK and BK for a principal diagnosis, ABF
// and BF for another diagnosis, BP and BO for an anesthesia-related
// procedure (HCPCS codes in the 837P, unlike the 837I's ICD procedure
// qualifiers), and BG for a condition code. Transmittal 579's Values column
// lists them per box (PDF p.42-55), and the CEDI worked transaction shows BK
// ("HI*BK>71598"). The parser places each HI segment in its box by HI01-1.
// The names are the guide's.
type DiagnosisTypeCodeHI0x_1 string

const (
	DiagnosisTypeCodeHI0x_1NotSent DiagnosisTypeCodeHI0x_1 = ""
	DiagnosisTypeCodeHI0x_1ABK     DiagnosisTypeCodeHI0x_1 = "ABK"
	DiagnosisTypeCodeHI0x_1BK      DiagnosisTypeCodeHI0x_1 = "BK"
	DiagnosisTypeCodeHI0x_1ABF     DiagnosisTypeCodeHI0x_1 = "ABF"
	DiagnosisTypeCodeHI0x_1BF      DiagnosisTypeCodeHI0x_1 = "BF"
	DiagnosisTypeCodeHI0x_1BP      DiagnosisTypeCodeHI0x_1 = "BP"
	DiagnosisTypeCodeHI0x_1BO      DiagnosisTypeCodeHI0x_1 = "BO"
	DiagnosisTypeCodeHI0x_1BG      DiagnosisTypeCodeHI0x_1 = "BG"
)

var diagnosisTypeCodeHI0x_1List = codeList{
	{"ABK", "International Classification of Diseases Clinical Modification (ICD-10-CM) Principal Diagnosis", "International Classification of Diseases Clinical Modification (ICD-10-CM) Principal Diagnosis — HI01-1 of the 2300 Health Care Diagnosis Code box, Transmittal 579, PDF p.42 (Values column, \"ABK, BK\"); the name is the guide's: " + tr3},
	{"BK", "International Classification of Diseases Clinical Modification (ICD-9-CM) Principal Diagnosis", "International Classification of Diseases Clinical Modification (ICD-9-CM) Principal Diagnosis — HI01-1 of the 2300 Health Care Diagnosis Code box, Transmittal 579, PDF p.42 (Values column, \"ABK, BK\"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"HI*BK>71598\"); the name is the guide's: " + tr3},
	{"ABF", "International Classification of Diseases Clinical Modification (ICD-10-CM) Diagnosis", "International Classification of Diseases Clinical Modification (ICD-10-CM) Diagnosis — HI02-1 through HI12-1 of the 2300 Health Care Diagnosis Code box, Transmittal 579, PDF p.42-46 (Values column, \"ABF, BF\"); the name is the guide's: " + tr3},
	{"BF", "International Classification of Diseases Clinical Modification (ICD-9-CM) Diagnosis", "International Classification of Diseases Clinical Modification (ICD-9-CM) Diagnosis — HI02-1 through HI12-1 of the 2300 Health Care Diagnosis Code box, Transmittal 579, PDF p.42-46 (Values column, \"ABF, BF\"); the name is the guide's: " + tr3},
	{"BP", "Health Care Financing Administration Common Procedural Coding System Principal Procedure", "Health Care Financing Administration Common Procedural Coding System Principal Procedure — HI01-1 of the 2300 Anesthesia Related Procedure box, Transmittal 579, PDF p.46 (Values column); the name is the guide's: " + tr3},
	{"BO", "Health Care Financing Administration Common Procedural Coding System", "Health Care Financing Administration Common Procedural Coding System — HI02-1 of the 2300 Anesthesia Related Procedure box, Transmittal 579, PDF p.47 (Values column); the name is the guide's: " + tr3},
	{"BG", "Condition", "Condition — HI01-1 through HI12-1 of the 2300 Condition Information box, Transmittal 579, PDF p.51-55 (Values column); the name is the guide's: " + tr3},
}

func (v DiagnosisTypeCodeHI0x_1) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v DiagnosisTypeCodeHI0x_1) Known() bool { return diagnosisTypeCodeHI0x_1List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v DiagnosisTypeCodeHI0x_1) ShortDescription() string {
	return diagnosisTypeCodeHI0x_1List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v DiagnosisTypeCodeHI0x_1) Description() string {
	return diagnosisTypeCodeHI0x_1List.description(string(v))
}

// ParseDiagnosisTypeCodeHI0x_1 returns text as a DiagnosisTypeCodeHI0x_1,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseDiagnosisTypeCodeHI0x_1(text string) (DiagnosisTypeCodeHI0x_1, error) {
	return DiagnosisTypeCodeHI0x_1(text), diagnosisTypeCodeHI0x_1List.check("HI0n-1", text)
}

// CodeCategoryCRC01 is CRC01, the Code Category, on a CRC segment (2300
// ambulanceCertification, visionCondition, homebound, and epsdtReferral boxes;
// 2400 ambulanceCertification, hospiceEmployeeIndicator, and
// dmeConditionIndicator boxes).
//
// The list is every CRC01 code the 837P implementation guide allows across
// its CRC boxes. frontend/src/data/x12Names837p.ts names each CRC box by a slug
// (CLAIM_CRC_NAME, SERVICE_CRC_NAME) rather than by CRC01's own character.
// Transmittal 579's Values column prints CRC01's code at every CRC box; the
// parser places each CRC in its box's field by it. A code that is the only
// value of its box(es) takes the box title as its short name; the Vision box
// lists three (E1, E2, E3), named here from the guide.
type CodeCategoryCRC01 string

const (
	CodeCategoryCRC01NotSent CodeCategoryCRC01 = ""
	CodeCategoryCRC0107      CodeCategoryCRC01 = "07"
	CodeCategoryCRC01E1      CodeCategoryCRC01 = "E1"
	CodeCategoryCRC01E2      CodeCategoryCRC01 = "E2"
	CodeCategoryCRC01E3      CodeCategoryCRC01 = "E3"
	CodeCategoryCRC0175      CodeCategoryCRC01 = "75"
	CodeCategoryCRC01ZZ      CodeCategoryCRC01 = "ZZ"
	CodeCategoryCRC0170      CodeCategoryCRC01 = "70"
	CodeCategoryCRC0109      CodeCategoryCRC01 = "09"
)

var codeCategoryCRC01List = codeList{
	{"07", "Ambulance Certification", "Ambulance Certification — the only CRC01 value of the 2300 and 2400 Ambulance Certification boxes, Transmittal 579, PDF p.40, 83; frontend/src/data/x12Names837p.ts CLAIM_CRC_NAME.ambulanceCertification"},
	{"E1", "Spectacle Lenses", "Spectacle Lenses — one of three CRC01 values (E1, E2, E3) of the 2300 Patient Condition Information: Vision box, Transmittal 579, PDF p.40; the name is the guide's: " + tr3},
	{"E2", "Contact Lenses", "Contact Lenses — one of three CRC01 values (E1, E2, E3) of the 2300 Patient Condition Information: Vision box, Transmittal 579, PDF p.40; the name is the guide's: " + tr3},
	{"E3", "Spectacle Frames", "Spectacle Frames — one of three CRC01 values (E1, E2, E3) of the 2300 Patient Condition Information: Vision box, Transmittal 579, PDF p.40; the name is the guide's: " + tr3},
	{"75", "Homebound Indicator", "Homebound Indicator — the 2300 box's only CRC01 value, Transmittal 579, PDF p.41; frontend/src/data/x12Names837p.ts CLAIM_CRC_NAME.homebound; the guide's own name for the code is Functional Limitations"},
	{"ZZ", "EPSDT Referral", "EPSDT Referral — the 2300 box's only CRC01 value, Transmittal 579, PDF p.41; frontend/src/data/x12Names837p.ts CLAIM_CRC_NAME.epsdtReferral; the guide's own name for the code is Mutually Defined"},
	{"70", "Hospice Employee Indicator", "Hospice Employee Indicator — the 2400 box's only CRC01 value, Transmittal 579, PDF p.83; frontend/src/data/x12Names837p.ts SERVICE_CRC_NAME.hospiceEmployeeIndicator; the guide's own name for the code is Hospice"},
	{"09", "Condition Indicator Durable Medical Equipment", "Condition Indicator Durable Medical Equipment — the 2400 box's only CRC01 value, Transmittal 579, PDF p.84; frontend/src/data/x12Names837p.ts SERVICE_CRC_NAME.dmeConditionIndicator; the guide's own name for the code is Durable Medical Equipment Certification"},
}

func (v CodeCategoryCRC01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v CodeCategoryCRC01) Known() bool { return codeCategoryCRC01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v CodeCategoryCRC01) ShortDescription() string { return codeCategoryCRC01List.short(string(v)) }

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v CodeCategoryCRC01) Description() string { return codeCategoryCRC01List.description(string(v)) }

// ParseCodeCategoryCRC01 returns text as a CodeCategoryCRC01, with an error
// when the text is neither empty nor on the list. The value is kept either
// way.
func ParseCodeCategoryCRC01(text string) (CodeCategoryCRC01, error) {
	return CodeCategoryCRC01(text), codeCategoryCRC01List.check("CRC01", text)
}

// ConditionIndicatorCRC0x is CRC03-07, the Condition Indicator/Condition Code
// positions of a CRC segment.
//
// The list is the union of the codes the 837P implementation guide allows
// across its CRC boxes: the ambulance conditions (01-12) at the 2300 and
// 2400 Ambulance Certification boxes, the vision conditions (L1-L5) at the
// Patient Condition Information: Vision box, IH at the Homebound box, AV,
// NU, S2, and ST at the EPSDT Referral box, 65 at the Hospice Employee box,
// and 38 and ZV at the DME Condition Indicator box. Neither reviewed source
// restates these positions' individual code characters
// (frontend/src/data/x12Names837p.ts CRC_EL_NAME gives only their position
// names); the codes and names are the guide's.
type ConditionIndicatorCRC0x string

const (
	ConditionIndicatorCRC0xNotSent ConditionIndicatorCRC0x = ""
	ConditionIndicatorCRC0x01      ConditionIndicatorCRC0x = "01"
	ConditionIndicatorCRC0x04      ConditionIndicatorCRC0x = "04"
	ConditionIndicatorCRC0x05      ConditionIndicatorCRC0x = "05"
	ConditionIndicatorCRC0x06      ConditionIndicatorCRC0x = "06"
	ConditionIndicatorCRC0x07      ConditionIndicatorCRC0x = "07"
	ConditionIndicatorCRC0x08      ConditionIndicatorCRC0x = "08"
	ConditionIndicatorCRC0x09      ConditionIndicatorCRC0x = "09"
	ConditionIndicatorCRC0x12      ConditionIndicatorCRC0x = "12"
	ConditionIndicatorCRC0x38      ConditionIndicatorCRC0x = "38"
	ConditionIndicatorCRC0x65      ConditionIndicatorCRC0x = "65"
	ConditionIndicatorCRC0xAV      ConditionIndicatorCRC0x = "AV"
	ConditionIndicatorCRC0xIH      ConditionIndicatorCRC0x = "IH"
	ConditionIndicatorCRC0xL1      ConditionIndicatorCRC0x = "L1"
	ConditionIndicatorCRC0xL2      ConditionIndicatorCRC0x = "L2"
	ConditionIndicatorCRC0xL3      ConditionIndicatorCRC0x = "L3"
	ConditionIndicatorCRC0xL4      ConditionIndicatorCRC0x = "L4"
	ConditionIndicatorCRC0xL5      ConditionIndicatorCRC0x = "L5"
	ConditionIndicatorCRC0xNU      ConditionIndicatorCRC0x = "NU"
	ConditionIndicatorCRC0xS2      ConditionIndicatorCRC0x = "S2"
	ConditionIndicatorCRC0xST      ConditionIndicatorCRC0x = "ST"
	ConditionIndicatorCRC0xZV      ConditionIndicatorCRC0x = "ZV"
)

var conditionIndicatorCRC0xList = codeList{
	{"01", "Patient was admitted to a hospital", guide("Patient was admitted to a hospital", "CRC03-07 at the Ambulance Certification boxes")},
	{"04", "Patient was moved by stretcher", guide("Patient was moved by stretcher", "CRC03-07 at the Ambulance Certification boxes")},
	{"05", "Patient was unconscious or in shock", guide("Patient was unconscious or in shock", "CRC03-07 at the Ambulance Certification boxes")},
	{"06", "Patient was transported in an emergency situation", guide("Patient was transported in an emergency situation", "CRC03-07 at the Ambulance Certification boxes")},
	{"07", "Patient had to be physically restrained", guide("Patient had to be physically restrained", "CRC03-07 at the Ambulance Certification boxes")},
	{"08", "Patient had visible hemorrhaging", guide("Patient had visible hemorrhaging", "CRC03-07 at the Ambulance Certification boxes")},
	{"09", "Ambulance service was medically necessary", guide("Ambulance service was medically necessary", "CRC03-07 at the Ambulance Certification boxes")},
	{"12", "Patient is confined to a bed or chair", guide("Patient is confined to a bed or chair", "CRC03-07 at the Ambulance Certification boxes")},
	{"38", "Certification signed by the physician is on file at the supplier's office", guide("Certification signed by the physician is on file at the supplier's office", "CRC03 at the 2400 Condition Indicator Durable Medical Equipment box")},
	{"65", "Open", guide("Open", "CRC03 at the 2400 Hospice Employee Indicator box")},
	{"AV", "Available - Not Used", guide("Available - Not Used", "CRC03 at the 2300 EPSDT Referral box")},
	{"IH", "Independent at Home", guide("Independent at Home", "CRC03 at the 2300 Homebound Indicator box")},
	{"L1", "General Standard of 20 Degree or .5 Diopter Sphere or Cylinder Change Met", guide("General Standard of 20 Degree or .5 Diopter Sphere or Cylinder Change Met", "CRC03-07 at the 2300 Patient Condition Information: Vision box")},
	{"L2", "Replacement Due to Loss or Theft", guide("Replacement Due to Loss or Theft", "CRC03-07 at the 2300 Patient Condition Information: Vision box")},
	{"L3", "Replacement Due to Breakage or Damage", guide("Replacement Due to Breakage or Damage", "CRC03-07 at the 2300 Patient Condition Information: Vision box")},
	{"L4", "Replacement Due to Patient Preference", guide("Replacement Due to Patient Preference", "CRC03-07 at the 2300 Patient Condition Information: Vision box")},
	{"L5", "Replacement Due to Medical Reason", guide("Replacement Due to Medical Reason", "CRC03-07 at the 2300 Patient Condition Information: Vision box")},
	{"NU", "Not Used", guide("Not Used", "CRC03 at the 2300 EPSDT Referral box")},
	{"S2", "Under Treatment", guide("Under Treatment", "CRC03 at the 2300 EPSDT Referral box")},
	{"ST", "New Services Requested", guide("New Services Requested", "CRC03 at the 2300 EPSDT Referral box")},
	{"ZV", "Replacement Item", guide("Replacement Item", "CRC03-04 at the 2400 Condition Indicator Durable Medical Equipment box")},
}

func (v ConditionIndicatorCRC0x) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ConditionIndicatorCRC0x) Known() bool { return conditionIndicatorCRC0xList.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ConditionIndicatorCRC0x) ShortDescription() string {
	return conditionIndicatorCRC0xList.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ConditionIndicatorCRC0x) Description() string {
	return conditionIndicatorCRC0xList.description(string(v))
}

// ParseConditionIndicatorCRC0x returns text as a ConditionIndicatorCRC0x,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseConditionIndicatorCRC0x(text string) (ConditionIndicatorCRC0x, error) {
	return ConditionIndicatorCRC0x(text), conditionIndicatorCRC0xList.check("CRC0x", text)
}

// YesNoConditionResponseCode is the Yes/No Condition or Response Code, the
// external X12 1073 list as the 837P uses it at CRC02, PAT09, SV109, SV111,
// SV112, and HI0x-9.
//
// The 837P implementation guide allows N and Y at every one of those
// positions (PAT09, SV109, SV111, and SV112 allow Y alone). Transmittal 579's
// Values column lists "N, Y" for the CRC02 Certification Condition Indicator
// (PDF p.40); the names are the guide's. HI0x-9 is Not Used in the 837P
// (Transmittal 907, PDF p.169-183) and no source lists values for it, so any
// other code there, such as the 837I's U or W, is kept as text and reports
// Known false.
type YesNoConditionResponseCode string

const (
	YesNoConditionResponseCodeNotSent YesNoConditionResponseCode = ""
	YesNoConditionResponseCodeN       YesNoConditionResponseCode = "N"
	YesNoConditionResponseCodeY       YesNoConditionResponseCode = "Y"
)

var yesNoConditionResponseCodeList = codeList{
	{"N", "No", "No — Transmittal 579, PDF p.40 (2300 CRC02 Values column, \"N, Y\"); the name is the guide's: " + tr3},
	{"Y", "Yes", "Yes — Transmittal 579, PDF p.40 (2300 CRC02 Values column, \"N, Y\"); the name is the guide's: " + tr3},
}

func (v YesNoConditionResponseCode) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v YesNoConditionResponseCode) Known() bool {
	return yesNoConditionResponseCodeList.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v YesNoConditionResponseCode) ShortDescription() string {
	return yesNoConditionResponseCodeList.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v YesNoConditionResponseCode) Description() string {
	return yesNoConditionResponseCodeList.description(string(v))
}

// ParseYesNoConditionResponseCode returns text as a
// YesNoConditionResponseCode, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseYesNoConditionResponseCode(text string) (YesNoConditionResponseCode, error) {
	return YesNoConditionResponseCode(text), yesNoConditionResponseCodeList.check("Yes/No Condition or Response Code", text)
}

// PricingMethodologyHCP01 is HCP01, the Pricing Methodology, on an HCP segment
// (2300 Claim Pricing/Repricing Information, 2400 Line Pricing/Repricing
// Information).
//
// The list is every code the 837P implementation guide allows at HCP01, with
// the guide's names. HCP is not used by Medicare at either the claim or line
// level (frontend/src/data/x12Names837p.ts NOT_USED_BY_MEDICARE: "2400/2300
// HCP (Line/Claim Pricing/Repricing Information)").
type PricingMethodologyHCP01 string

const (
	PricingMethodologyHCP01NotSent PricingMethodologyHCP01 = ""
	PricingMethodologyHCP0100      PricingMethodologyHCP01 = "00"
	PricingMethodologyHCP0101      PricingMethodologyHCP01 = "01"
	PricingMethodologyHCP0102      PricingMethodologyHCP01 = "02"
	PricingMethodologyHCP0103      PricingMethodologyHCP01 = "03"
	PricingMethodologyHCP0104      PricingMethodologyHCP01 = "04"
	PricingMethodologyHCP0105      PricingMethodologyHCP01 = "05"
	PricingMethodologyHCP0106      PricingMethodologyHCP01 = "06"
	PricingMethodologyHCP0107      PricingMethodologyHCP01 = "07"
	PricingMethodologyHCP0108      PricingMethodologyHCP01 = "08"
	PricingMethodologyHCP0109      PricingMethodologyHCP01 = "09"
	PricingMethodologyHCP0110      PricingMethodologyHCP01 = "10"
	PricingMethodologyHCP0111      PricingMethodologyHCP01 = "11"
	PricingMethodologyHCP0112      PricingMethodologyHCP01 = "12"
	PricingMethodologyHCP0113      PricingMethodologyHCP01 = "13"
	PricingMethodologyHCP0114      PricingMethodologyHCP01 = "14"
)

var pricingMethodologyHCP01List = codeList{
	{"00", "Zero Pricing (Not Covered Under Contract)", guide("Zero Pricing (Not Covered Under Contract)", "HCP01")},
	{"01", "Priced as Billed at 100%", guide("Priced as Billed at 100%", "HCP01")},
	{"02", "Priced at the Standard Fee Schedule", guide("Priced at the Standard Fee Schedule", "HCP01")},
	{"03", "Priced at a Contractual Percentage", guide("Priced at a Contractual Percentage", "HCP01")},
	{"04", "Bundled Pricing", guide("Bundled Pricing", "HCP01")},
	{"05", "Bundled Pricing", guide("Bundled Pricing", "HCP01")},
	{"06", "Peer Review Pricing", guide("Peer Review Pricing", "HCP01")},
	{"07", "Flat Rate Pricing", guide("Flat Rate Pricing", "HCP01")},
	{"08", "Combination Pricing", guide("Combination Pricing", "HCP01")},
	{"09", "Maximum Allowable Pricing", guide("Maximum Allowable Pricing", "HCP01")},
	{"10", "Other Pricing", guide("Other Pricing", "HCP01")},
	{"11", "Lower of Cost", guide("Lower of Cost", "HCP01")},
	{"12", "Ratio of Cost", guide("Ratio of Cost", "HCP01")},
	{"13", "Cost Reimbursed", guide("Cost Reimbursed", "HCP01")},
	{"14", "Adjustment Pricing", guide("Adjustment Pricing", "HCP01")},
}

func (v PricingMethodologyHCP01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v PricingMethodologyHCP01) Known() bool { return pricingMethodologyHCP01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v PricingMethodologyHCP01) ShortDescription() string {
	return pricingMethodologyHCP01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v PricingMethodologyHCP01) Description() string {
	return pricingMethodologyHCP01List.description(string(v))
}

// ParsePricingMethodologyHCP01 returns text as a PricingMethodologyHCP01,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParsePricingMethodologyHCP01(text string) (PricingMethodologyHCP01, error) {
	return PricingMethodologyHCP01(text), pricingMethodologyHCP01List.check("HCP01", text)
}
