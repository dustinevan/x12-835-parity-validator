package x12837i

// ClaimCLM holds one CLM segment (loop 2300, "Claim Information") and
// everything the institutional 837 nests under it, down through its
// service lines. Transmittal 907, PDF p.23 onward; frontend/src/data/
// x12Names837i.ts's "Claim (2300)" section names every element below.
type ClaimCLM struct {
	Index int `json:"index,omitempty"` // the CLM's position in the file

	PatientControlCLM01 string `json:"patientControlCLM01,omitempty"` // Patient Control Number; CLM01; Transmittal 907, PDF p.23
	TotalChargeCLM02    Amount `json:"totalChargeCLM02,omitempty"`    // Total Claim Charge Amount; CLM02; PDF p.23

	FilingIndicatorCLM03 string `json:"filingIndicatorCLM03,omitempty"` // Claim Filing Indicator Code (N/U); CLM03; PDF p.23
	ClaimTypeCLM04       string `json:"claimTypeCLM04,omitempty"`       // Non-Institutional Claim Type Code (N/U); CLM04; PDF p.23

	// CLM05: the Health Care Service Location Information composite.
	// CLM05-1 and CLM05-3 are NUBC Type of Bill digits (facility type +
	// bill classification, and frequency); Transmittal 907's own "Values"
	// column is blank for both (PDF p.23), and the NUBC's Official UB-04
	// Data Specifications Manual that would name them is a paid
	// publication, not a public source available to this project — see
	// the NUBC note at the top of frontend/src/data/x12Names837i.ts. The
	// two are therefore plain strings with no list: a value is kept as
	// sent and never judged.
	FacilityTypeCLM05_1          string                       `json:"facilityTypeCLM05_1,omitempty"`          // Place of Service Code; CLM05-1; PDF p.23 (external code list, NUBC Type of Bill digits 1-2; no public values)
	FacilityCodeQualifierCLM05_2 FacilityCodeQualifierCLM05_2 `json:"facilityCodeQualifierCLM05_2,omitempty"` // Facility Code Qualifier; CLM05-2; PDF p.23
	ClaimFrequencyCLM05_3        string                       `json:"claimFrequencyCLM05_3,omitempty"`        // Claim Frequency Code; CLM05-3; PDF p.23 (external code list, NUBC Type of Bill digit 4; no public values)

	SignatureIndicatorCLM06 string              `json:"signatureIndicatorCLM06,omitempty"` // Provider or Supplier Signature Indicator (N/U); CLM06; PDF p.23
	AssignmentCLM07         AssignmentCLM07     `json:"assignmentCLM07,omitempty"`         // Assignment or Plan Participation Code; CLM07; PDF p.23
	BenefitsAssignCLM08     BenefitsAssignCLM08 `json:"benefitsAssignCLM08,omitempty"`     // Benefits Assignment Certification Indicator; CLM08; PDF p.23
	ReleaseOfInfoCLM09      ReleaseOfInfoCLM09  `json:"releaseOfInfoCLM09,omitempty"`      // Release of Information Code; CLM09; PDF p.23

	SignatureSourceCLM10        string           `json:"signatureSourceCLM10,omitempty"`        // Patient Signature Source Code (N/U); CLM10; PDF p.23
	RelatedCausesCLM11          string           `json:"relatedCausesCLM11,omitempty"`          // Related Causes Information composite (N/U), kept as one text with its component separators; CLM11; PDF p.23
	SpecialProgramCLM12         string           `json:"specialProgramCLM12,omitempty"`         // Special Program Indicator (N/U); CLM12; PDF p.23
	YesNoCLM13                  string           `json:"yesNoCLM13,omitempty"`                  // Yes/No Condition or Response Code (N/U); CLM13; PDF p.23
	LevelOfServiceCLM14         string           `json:"levelOfServiceCLM14,omitempty"`         // Level of Service Code (N/U); CLM14; PDF p.23
	YesNoCLM15                  string           `json:"yesNoCLM15,omitempty"`                  // Yes/No Condition or Response Code (N/U); CLM15; PDF p.23
	ParticipationAgreementCLM16 string           `json:"participationAgreementCLM16,omitempty"` // Participation Agreement (N/U); CLM16; PDF p.23
	StatusCLM17                 string           `json:"statusCLM17,omitempty"`                 // Claim Status Code (N/U); CLM17; PDF p.24
	YesNoCLM18                  string           `json:"yesNoCLM18,omitempty"`                  // Yes/No Condition or Response Code (N/U); CLM18; PDF p.24
	SubmissionReasonCLM19       string           `json:"submissionReasonCLM19,omitempty"`       // Claim Submission Reason Code (N/U); CLM19; PDF p.24
	DelayReasonCLM20            DelayReasonCLM20 `json:"delayReasonCLM20,omitempty"`            // Delay Reason Code; CLM20; PDF p.24

	// Dates (DTP, loop 2300, qualifier-specific). DTP_CLAIM_NAME in
	// frontend/src/data/x12Names837i.ts. Each is one DateDTP read in the
	// form its qualifier allows (Transmittal 907, PDF p.24): 096 TM, an
	// hour alone; 434 RD8 or D8, a date range or one date; 435 DT or D8, a
	// date with an hour or one date; 050 D8. A DTP in another form, or a
	// repeat, is kept raw.
	DischargeHourDTP096    DateDTP `json:"dischargeHourDTP096,omitempty"`    // Discharge Hour; DTP*096; PDF p.24
	StatementDatesDTP434   DateDTP `json:"statementDatesDTP434,omitempty"`   // Statement From-Through Dates; DTP*434; PDF p.24
	AdmissionDTP435        DateDTP `json:"admissionDTP435,omitempty"`        // Admission Date and Hour; DTP*435; PDF p.24
	RepricerReceivedDTP050 DateDTP `json:"repricerReceivedDTP050,omitempty"` // Repricer Received Date; DTP*050; PDF p.24

	InstitutionalCL1 InstitutionalCL1 `json:"institutionalCL1,omitempty"` // CL1, Institutional Claim Code; PDF p.24-25

	SupplementalInfo []SupplementalPWK `json:"supplementalInfo,omitempty"` // PWK, Claim Supplemental Information; PDF p.25
	Contracts        []ContractCN1     `json:"contracts,omitempty"`        // CN1, Contract Information; PDF p.25-26

	// AMT (loop 2300). AMT01_NAME in frontend/src/data/x12Names837i.ts
	// documents only F3 at this loop; the table's other codes (D, EAF,
	// A8, GT, N8) belong to loop 2320 and loop 2400, not here.
	PatientEstimatedAmountAMTF3 AmountAMT   `json:"patientEstimatedAmountAMTF3,omitempty"` // Patient Estimated Amount Due; AMT*F3; PDF p.26
	Amounts                     []AmountAMT `json:"amounts,omitempty"`                     // any other claim-level AMT, or a repeat of F3

	// REF (loop 2300), named by qualifier per REF_CLAIM_NAME in
	// frontend/src/data/x12Names837i.ts (source: Transmittal 907,
	// PDF p.26-28). Each holds the first REF with its qualifier that
	// carries a value; the qualifier is carried by the field name and by
	// the segment itself, the same convention x12835's ServiceLineSVC uses
	// for its own REF fields.
	ServiceAuthExceptionREF4N     ReferenceREF `json:"serviceAuthExceptionREF4N,omitempty"`     // Service Authorization Exception Code; REF*4N; PDF p.26 (REF02 is itself coded 1-7, REF_SERVICE_AUTH_VALUES, no per-code text given)
	ReferralNumberREF9F           ReferenceREF `json:"referralNumberREF9F,omitempty"`           // Referral Number; REF*9F; PDF p.26
	PriorAuthREFG1                ReferenceREF `json:"priorAuthREFG1,omitempty"`                // Prior Authorization; REF*G1; PDF p.26
	PayerClaimControlREFF8        ReferenceREF `json:"payerClaimControlREFF8,omitempty"`        // Payer Claim Control Number; REF*F8; PDF p.27 (used on a replacement/void claim to name the original claim's ICN)
	RepricedClaimREF9A            ReferenceREF `json:"repricedClaimREF9A,omitempty"`            // Repriced Claim Number; REF*9A; PDF p.27
	AdjustedRepricedClaimREF9C    ReferenceREF `json:"adjustedRepricedClaimREF9C,omitempty"`    // Adjusted Repriced Claim Number; REF*9C; PDF p.27
	InvestigationalDeviceREFLX    ReferenceREF `json:"investigationalDeviceREFLX,omitempty"`    // Investigational Device Exemption Number; REF*LX; PDF p.27
	TransmissionIntermediaryREFD9 ReferenceREF `json:"transmissionIntermediaryREFD9,omitempty"` // Claim Identifier for Transmission Intermediaries; REF*D9; PDF p.27-28 (REF02 labeled "Value Added Network Trace Number")
	AutoAccidentStateREFLU        ReferenceREF `json:"autoAccidentStateREFLU,omitempty"`        // Auto Accident State or Province; REF*LU; PDF p.28
	MedicalRecordREFEA            ReferenceREF `json:"medicalRecordREFEA,omitempty"`            // Medical Record Number; REF*EA; PDF p.28
	DemonstrationProjectREFP4     ReferenceREF `json:"demonstrationProjectREFP4,omitempty"`     // Demonstration Project Identifier; REF*P4; PDF p.28
	PROApprovalREFG4              ReferenceREF `json:"proApprovalREFG4,omitempty"`              // Peer Review Organization (PRO) Approval Number; REF*G4; PDF p.28

	// References holds any claim-level REF not one of the qualified
	// fields above (including the unnamed REF*CN noted in
	// frontend/src/data/x12Names837i.ts just above K3), or a repeat of one.
	References []ReferenceREF `json:"references,omitempty"`

	Files []FileInfoK3 `json:"files,omitempty"` // K3, File Information; PDF p.29

	Notes []NoteNTE `json:"notes,omitempty"` // NTE, Claim Note; PDF p.29

	EPSDTReferralsCRCZZ []ReferralCRC `json:"epsdtReferralsCRCZZ,omitempty"` // CRC, EPSDT Referral; PDF p.29

	// HI groups (loop 2300). Each role is one qualifier of the shared HI0n
	// composite (HealthInfoHI); HI_TYPE_NAME in
	// frontend/src/data/x12Names837i.ts names each qualifier, pairing the
	// ICD-9 and ICD-10 codes for the same concept under one role here,
	// since a given claim sends only one of the pair. The singular fields
	// hold at most one composite; the plural fields can repeat across
	// several HI segments (e.g. multiple Other Diagnosis codes). Each
	// composite remembers its segment and position, so the HI segments
	// are written back as they came.
	PrincipalDiagnosisHIABK HealthInfoHI   `json:"principalDiagnosisHIABK,omitempty"` // ABK/BK; PDF p.30
	AdmittingDiagnosisHIABJ HealthInfoHI   `json:"admittingDiagnosisHIABJ,omitempty"` // ABJ/BJ; PDF p.34
	PatientReasonsHIAPR     []HealthInfoHI `json:"patientReasonsHIAPR,omitempty"`     // APR/PR; PDF p.39 (HI01-HI03)
	ExternalCausesHIABN     []HealthInfoHI `json:"externalCausesHIABN,omitempty"`     // ABN/BN; PDF p.44 (HI01-HI12)
	OtherDiagnosesHIABF     []HealthInfoHI `json:"otherDiagnosesHIABF,omitempty"`     // ABF/BF; PDF p.51
	PrincipalProcedureHIBBR HealthInfoHI   `json:"principalProcedureHIBBR,omitempty"` // BBR/BR; PDF p.58
	OtherProceduresHIBBQ    []HealthInfoHI `json:"otherProceduresHIBBQ,omitempty"`    // BBQ/BQ; PDF p.62
	OccurrenceSpansHIBI     []HealthInfoHI `json:"occurrenceSpansHIBI,omitempty"`     // BI; PDF p.67 (NUBC code value; see NUBC note)
	OccurrenceCodesHIBH     []HealthInfoHI `json:"occurrenceCodesHIBH,omitempty"`     // BH; PDF p.72 (NUBC code value; see NUBC note)
	ValueCodesHIBE          []HealthInfoHI `json:"valueCodesHIBE,omitempty"`          // BE; PDF p.77 (NUBC code value; see NUBC note)
	ConditionCodesHIBG      []HealthInfoHI `json:"conditionCodesHIBG,omitempty"`      // BG; PDF p.81 (NUBC code value; see NUBC note)
	TreatmentCodesHITC      []HealthInfoHI `json:"treatmentCodesHITC,omitempty"`      // TC; PDF p.85
	DRGHIDR                 HealthInfoHI   `json:"drgHIDR,omitempty"`                 // DR; PDF p.48

	Pricing []PricingHCP `json:"pricing,omitempty"` // HCP, Claim Pricing/Repricing Information; PDF p.91-92

	// Claim-level providers (2310A-F). NM1_ROLE_NAME in
	// frontend/src/data/x12Names837i.ts names each role by its NM101 code.
	AttendingNM171       AttendingNM171       `json:"attendingNM171,omitempty"`       // 2310A, role 71; PDF p.92
	OperatingNM172       OperatingNM172       `json:"operatingNM172,omitempty"`       // 2310B, role 72; PDF p.93
	OtherOperatingNM1ZZ  OtherOperatingNM1ZZ  `json:"otherOperatingNM1ZZ,omitempty"`  // 2310C, role ZZ; PDF p.94
	RenderingNM182       RenderingNM182       `json:"renderingNM182,omitempty"`       // 2310D, role 82; PDF p.94
	ServiceFacilityNM177 ServiceFacilityNM177 `json:"serviceFacilityNM177,omitempty"` // 2310E, role 77; PDF p.95
	ReferringNM1DN       ReferringNM1DN       `json:"referringNM1DN,omitempty"`       // 2310F, role DN; PDF p.96-97

	// OtherSubscribers holds the 2320/2330A-I coordination-of-benefits
	// loops; see other_subscriber.go.
	OtherSubscribers []OtherSubscriberSBR `json:"otherSubscribers,omitempty"`

	ServiceLines []ServiceLineLX `json:"serviceLines,omitempty"` // LX, loop 2400, one per service line

	Raw []RawSegment `json:"raw,omitempty"` // segments within this claim the model does not place
}

// InstitutionalCL1 holds the CL1 segment (loop 2300, Institutional Claim
// Code). Set is false when the segment did not appear at all. CL101-103
// are NUBC-maintained code values (UB-04 FL14, FL15, FL17); Transmittal
// 907's "Values" column is blank for all three (PDF p.24-25), and CMS
// Claims Processing Manual Ch. 25 §§75.1-75.2 confirm they come from
// Medicare contractors and the NUBC's non-public Official UB-04 Data
// Specifications Manual, not CMS itself — see the NUBC note at the top of
// frontend/src/data/x12Names837i.ts. The three are therefore plain strings
// with no list, the same situation as CLM05-1 and CLM05-3 above.
type InstitutionalCL1 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	AdmissionTypeCL101   string `json:"admissionTypeCL101,omitempty"`   // Admission Type Code; CL101; PDF p.24 (external code list, NUBC FL14; no public values)
	AdmissionSourceCL102 string `json:"admissionSourceCL102,omitempty"` // Admission Source Code; CL102; PDF p.24 (external code list, NUBC FL15; no public values)
	PatientStatusCL103   string `json:"patientStatusCL103,omitempty"`   // Patient Status Code; CL103; PDF p.25 (external code list, NUBC FL17; no public values)
	NursingHomeCL104     string `json:"nursingHomeCL104,omitempty"`     // Nursing Home Residential Status Code (N/U); CL104; PDF p.25
}

// SupplementalPWK holds one PWK segment: Claim Supplemental Information at
// loop 2300 (ClaimCLM.SupplementalInfo) and Line Supplemental Information
// at loop 2400 (ServiceLineLX.SupplementalInfo), the identical shape at
// both (Transmittal 907, PDF p.25 claim-level, PDF p.112-113 line-level).
// PWKReportTypeCode01 and PWKTransmissionCode02 (enums_claim.go) hold the
// union of the two levels' code lists.
type SupplementalPWK struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	ReportTypePWK01        PWKReportTypeCode01   `json:"reportTypePWK01,omitempty"`        // Report Type Code; PWK01
	TransmissionPWK02      PWKTransmissionCode02 `json:"transmissionPWK02,omitempty"`      // Report Transmission Code; PWK02
	CopiesNeededPWK03      string                `json:"copiesNeededPWK03,omitempty"`      // Report Copies Needed (N/U); PWK03
	EntityIDPWK04          string                `json:"entityIDPWK04,omitempty"`          // Entity Identifier Code (N/U); PWK04
	IDQualifierPWK05       string                `json:"idQualifierPWK05,omitempty"`       // Identification Code Qualifier; PWK05 (no code list given in the reviewed sources)
	AttachmentControlPWK06 string                `json:"attachmentControlPWK06,omitempty"` // Attachment Control Number; PWK06

	// ExtraElements holds PWK07 through PWK09 (Description, the Actions
	// Indicated composite as one text, Request Category Code) as sent. X12
	// defines them; the 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// ContractCN1 holds one CN1 segment (loop 2300, Contract Information).
// CN101's code list (01-06, 09; Transmittal 907, PDF p.25) is a generic X12
// code, not an NUBC one, but is given without per-code text, so it stays a
// plain string rather than a coded type, the same reason CLM05-1 and
// CLM05-3 do.
type ContractCN1 struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	ContractTypeCN101    string `json:"contractTypeCN101,omitempty"`    // Contract Type Code; CN101 (code list only, values 01-06, 09; no per-code text)
	ContractAmountCN102  Amount `json:"contractAmountCN102,omitempty"`  // Contract Amount; CN102
	ContractPercentCN103 Amount `json:"contractPercentCN103,omitempty"` // Contract Percentage; CN103
	ContractCodeCN104    string `json:"contractCodeCN104,omitempty"`    // Contract Code; CN104
	TermsDiscountCN105   Amount `json:"termsDiscountCN105,omitempty"`   // Terms Discount Percent; CN105
	ContractVersionCN106 string `json:"contractVersionCN106,omitempty"` // Contract Version Identifier; CN106
}

// FileInfoK3 holds one K3 segment (loop 2300, File Information). K302
// (Record Format Code) and K303 (Composite Unit of Measure) are Not Used
// and kept as sent in ExtraElements.
type FileInfoK3 struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	FixedFormatK301 string `json:"fixedFormatK301,omitempty"` // Fixed Format Information; K301

	// ExtraElements holds K302 and K303 (the composite as one text) as
	// sent. X12 defines them; the 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// NoteNTE holds one NTE segment at the claim level (loop 2300): Claim Note
// (NTE01 ALG through UPI) or Billing Note (NTE01 ADD), both Transmittal
// 907, PDF p.29. NTENoteCode01 (enums_claim.go) also covers the line-level
// NTE (2400, Third Party Organization Notes, TPO), which ServiceLineLX does
// not model; a line-level NTE is kept as a RawSegment.
type NoteNTE struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	NoteCodeNTE01 NTENoteCode01 `json:"noteCodeNTE01,omitempty"` // Note Reference Code; NTE01
	TextNTE02     string        `json:"textNTE02,omitempty"`     // Note Text; NTE02
}

// ReferralCRC holds one CRC segment (loop 2300, EPSDT Referral). CRC06-07
// (Condition Indicator, two further repeats) are Not Used and kept as sent
// in ExtraElements.
type ReferralCRC struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	CategoryCRC01      CRCCategoryCRC01      `json:"categoryCRC01,omitempty"`      // Code Category; CRC01
	CertificationCRC02 CRCCertificationCRC02 `json:"certificationCRC02,omitempty"` // Certification Condition Code Applies Indicator; CRC02
	Condition1CRC03    CRCConditionIndicator `json:"condition1CRC03,omitempty"`    // Condition Indicator; CRC03
	Condition2CRC04    CRCConditionIndicator `json:"condition2CRC04,omitempty"`    // Condition Indicator; CRC04
	Condition3CRC05    CRCConditionIndicator `json:"condition3CRC05,omitempty"`    // Condition Indicator; CRC05

	// ExtraElements holds CRC06 and CRC07 as sent. X12 defines them; the
	// 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// PricingHCP holds one HCP segment: Claim Pricing/Repricing Information at
// loop 2300 (ClaimCLM.Pricing) and Line Pricing/Repricing Information at
// loop 2400 (ServiceLineLX.Pricing), the identical shape at both
// (Transmittal 907, PDF p.91-92 claim-level, PDF p.114-115 line-level).
// The HCP01 and HCP13-15 coded types are in enums_claim.go;
// ProductServiceQualifier and UnitOrBasisOfMeasurement are in
// enums_service.go. Set is false when the segment did not appear at all.
type PricingHCP struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	PricingMethodologyHCP01 HCPPricingMethodologyHCP01 `json:"pricingMethodologyHCP01,omitempty"` // Pricing Methodology; HCP01
	AllowedAmountHCP02      Amount                     `json:"allowedAmountHCP02,omitempty"`      // Repriced Allowed Amount; HCP02
	SavingsAmountHCP03      Amount                     `json:"savingsAmountHCP03,omitempty"`      // Repriced Saving Amount; HCP03
	OrganizationHCP04       string                     `json:"organizationHCP04,omitempty"`       // Repricing Organization Identifier; HCP04
	PerDiemAmountHCP05      Amount                     `json:"perDiemAmountHCP05,omitempty"`      // Repriced Per Diem or Flat Rate Amount; HCP05
	ApprovedDRGHCP06        string                     `json:"approvedDRGHCP06,omitempty"`        // Repriced Approved DRG Code; HCP06
	ApprovedAmountHCP07     Amount                     `json:"approvedAmountHCP07,omitempty"`     // Repriced Approved Amount; HCP07
	ProductServiceHCP08     string                     `json:"productServiceHCP08,omitempty"`     // Product/Service ID (N/U); HCP08
	ProductQualifierHCP09   ProductServiceQualifier    `json:"productQualifierHCP09,omitempty"`   // Product or Service ID Qualifier; HCP09
	ApprovedHCPCSHCP10      string                     `json:"approvedHCPCSHCP10,omitempty"`      // Repriced Approved HCPCS Code; HCP10
	UnitOrBasisHCP11        UnitOrBasisOfMeasurement   `json:"unitOrBasisHCP11,omitempty"`        // Unit or Basis for Measurement Code; HCP11 (values DA, UN per HCP11_VALUES; UnitOrBasisOfMeasurement, enums_service.go, is the union with CTP05-1's own F2/GR/ME/ML/UN)
	QuantityHCP12           Amount                     `json:"quantityHCP12,omitempty"`           // Quantity; HCP12
	RejectReasonHCP13       HCPRejectReasonHCP13       `json:"rejectReasonHCP13,omitempty"`       // Reject Reason Code; HCP13
	PolicyComplianceHCP14   HCPPolicyComplianceHCP14   `json:"policyComplianceHCP14,omitempty"`   // Policy Compliance Code; HCP14
	ExceptionHCP15          HCPExceptionCodeHCP15      `json:"exceptionHCP15,omitempty"`          // Exception Code; HCP15
}

// AttendingNM171 is loop 2310A, Attending Provider: NM1 (role 71), PRV
// (Attending Provider Specialty Information — the only one of 2310A-F that
// carries a PRV), and up to four REF repeats (secondary identification,
// qualifiers 0B/1G/G2/LU). NM1_ATTENDING_EL_NAME, PROVIDER_REF_NAME,
// REF_PROVIDER_EL_NAME in frontend/src/data/x12Names837i.ts (source:
// Transmittal 907, PDF p.92).
type AttendingNM171 struct {
	Name       NameNM1        `json:"name,omitempty"`
	Specialty  ProviderPRV    `json:"specialty,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// OperatingNM172 is loop 2310B, Operating Physician: NM1 (role 72) and up
// to four REF repeats, the same qualifier set as AttendingNM171's.
// Transmittal 907, PDF p.93; the flat file's own column mislabels NM105
// there (see NM1_OPERATING_EL_NAME's note).
type OperatingNM172 struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// OtherOperatingNM1ZZ is loop 2310C, Other Operating Physician: NM1 (role
// ZZ) and up to four REF repeats. Transmittal 907, PDF p.94.
type OtherOperatingNM1ZZ struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// RenderingNM182 is loop 2310D, Rendering Provider: NM1 (role 82) and up to
// four REF repeats. Transmittal 907, PDF p.94.
type RenderingNM182 struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// ServiceFacilityNM177 is loop 2310E, Service Facility Location: NM1 (role
// 77) and up to four REF repeats. Transmittal 907, PDF p.95.
type ServiceFacilityNM177 struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// ReferringNM1DN is loop 2310F, Referring Provider: NM1 (role DN) and up to
// four REF repeats. Transmittal 907, PDF p.96-97.
type ReferringNM1DN struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}
