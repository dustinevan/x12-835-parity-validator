package x12837p

// ClaimCLM holds one CLM segment (loop 2300) and everything the 837P nests
// under it: dates, supplemental information, claim-level providers
// (2310A-F), coordination-of-benefits subscribers (2320/2330A-G), and
// service lines (2400). Citations below are to
// docs/reference/file-formats/cms-837-835-flatfiles-r907otn-2011.pdf
// ("Transmittal 907, PDF p.N") and
// docs/reference/file-formats/cedi-837p-companion-guide.pdf ("CEDI 837P CG,
// PDF p.N"), per the citation form documented at the top of
// frontend/src/data/x12Names837p.ts. Most REF01, DTP01, NM101, and AMT01
// qualifier code characters are not restated in either source for this
// loop; where that is true, the field below is named from its box title
// only and the comment says so, matching x12Names837p.ts's own citation
// discipline.
type ClaimCLM struct {
	Index int `json:"index,omitempty"` // the CLM's position in the file, 0-based; the ISA is 0

	// CLM01-CLM20. CLM03, CLM04, and CLM13-CLM19 are Not Used in 005010A1
	// (CLM03's Claim Filing Indicator Code lives in SBR09 instead); they are
	// still given fields here, kept as plain text, so an exact difference in
	// a position the standard forbids is never silently dropped (every
	// difference is retained).
	PatientControlCLM01 string `json:"patientControlCLM01,omitempty"` // Patient Control Number (PDF literal box text: "Patient Account Number"); CLM_NAME; Transmittal 907, PDF p.155
	TotalChargeCLM02    Amount `json:"totalChargeCLM02,omitempty"`    // Total Claim Charge Amount; PDF p.155

	ClaimFilingIndicatorCLM03 string `json:"claimFilingIndicatorCLM03,omitempty"` // Claim Filing Indicator Code; Not Used in 005010A1 (SBR09 carries this instead); PDF p.155
	ClaimTypeCLM04            string `json:"claimTypeCLM04,omitempty"`            // Non-Institutional Claim Type Code; Not Used in 005010A1; PDF p.155

	// CLM05, the "Health Care Service Location Information" composite.
	// CLM05-1 is the Place of Service Code, an external CMS list the
	// repository does not hold, so it stays plain text.
	PlaceOfServiceCLM05_1    string                       `json:"placeOfServiceCLM05_1,omitempty"`    // Place of Service Code; PDF p.155
	FacilityQualifierCLM05_2 FacilityCodeQualifierCLM05_2 `json:"facilityQualifierCLM05_2,omitempty"` // Facility Code Qualifier; PDF p.155
	FrequencyCLM05_3         ClaimFrequencyCLM05_3        `json:"frequencyCLM05_3,omitempty"`         // Claim Frequency Code; PDF p.156; Medicare restricts this to 1 (Original), CEDI 837P CG, PDF p.36 (printed p.28; Table 18)
	LocationExtraCLM05       []string                     `json:"locationExtraCLM05,omitempty"`       // components after CLM05-3, kept as sent

	SignatureCLM06          ProviderOrSupplierSignatureCLM06     `json:"signatureCLM06,omitempty"`          // Provider or Supplier Signature Indicator; PDF p.156
	AssignmentCLM07         AssignmentOrPlanParticipationCLM07   `json:"assignmentCLM07,omitempty"`         // Assignment or Plan Participation Code; PDF p.156
	BenefitsAssignmentCLM08 BenefitsAssignmentCertificationCLM08 `json:"benefitsAssignmentCLM08,omitempty"` // Benefits Assignment Certification Indicator; PDF p.156
	ReleaseOfInfoCLM09      ReleaseOfInformationCLM09            `json:"releaseOfInfoCLM09,omitempty"`      // Release of Information Code; PDF p.156
	SignatureSourceCLM10    PatientSignatureSourceCLM10          `json:"signatureSourceCLM10,omitempty"`    // Patient Signature Source Code; PDF p.156

	// CLM11, the "Related Causes Information" composite. CLM11-1 (required
	// whenever CLM11 is sent) and CLM11-2 (situational, a second cause) share
	// one code list, so both positions reuse the single RelatedCausesCodeCLM11_1
	// type the names file lists; CLM11-3 is Not Used and kept as plain text.
	RelatedCausesCLM11       [2]RelatedCausesCodeCLM11_1 `json:"relatedCausesCLM11,omitempty"`       // CLM11-1, CLM11-2; PDF p.156
	RelatedCausesCLM11_3     string                      `json:"relatedCausesCLM11_3,omitempty"`     // Related Causes Code (N/U) — PDF p.156
	AutoAccidentStateCLM11_4 string                      `json:"autoAccidentStateCLM11_4,omitempty"` // Auto Accident State or Province Code; PDF p.156
	CountryCLM11_5           string                      `json:"countryCLM11_5,omitempty"`           // Country Code; PDF p.156
	RelatedCausesExtraCLM11  []string                    `json:"relatedCausesExtraCLM11,omitempty"`  // components after CLM11-5, kept as sent

	SpecialProgramCLM12 string `json:"specialProgramCLM12,omitempty"` // Special Program Indicator; no confirmed value table in the reviewed sources; PDF p.156

	YesNoCLM13            string `json:"yesNoCLM13,omitempty"`            // Yes/No Condition or Response Code; Not Used in 005010A1; PDF p.157
	LevelOfServiceCLM14   string `json:"levelOfServiceCLM14,omitempty"`   // Level of Service Code; Not Used in 005010A1; PDF p.157
	YesNoCLM15            string `json:"yesNoCLM15,omitempty"`            // Yes/No Condition or Response Code; Not Used in 005010A1; PDF p.157
	ParticipationCLM16    string `json:"participationCLM16,omitempty"`    // Participation Agreement; Not Used in 005010A1; PDF p.157
	ClaimStatusCLM17      string `json:"claimStatusCLM17,omitempty"`      // Claim Status Code; Not Used in 005010A1; PDF p.157
	YesNoCLM18            string `json:"yesNoCLM18,omitempty"`            // Yes/No Condition or Response Code; Not Used in 005010A1; PDF p.157
	SubmissionReasonCLM19 string `json:"submissionReasonCLM19,omitempty"` // Claim Submission Reason Code; Not Used in 005010A1; PDF p.157

	DelayReasonCLM20 DelayReasonCLM20 `json:"delayReasonCLM20,omitempty"` // Delay Reason Code; PDF p.157; Medicare ignores any value submitted here, CEDI 837P CG, PDF p.36 (printed p.28; Table 18)

	Extra []string `json:"extra,omitempty"` // elements after CLM20, kept as sent

	// Dates (loop 2300, DTP). CLAIM_DTP_NAME / CLAIM_DTP03_EL_NAME in
	// frontend/src/data/x12Names837p.ts name each box. Each box's DTP01 code
	// is the one Transmittal 579's Values column gives it (see
	// DateTimeQualifierDTP01), and the parser places a DTP in the field for
	// its code; Admission's 435 is also confirmed by CEDI 837P CG, PDF p.37
	// (printed p.29). DateDTP keeps DTP01, DTP02, and DTP03 as received, as
	// the service line's dates do. A DTP whose box is already filled, or
	// whose DTP01 matches no box, is kept in Raw.
	OnsetDTP431                        DateDTP    `json:"onsetDTP431,omitempty"`                        // Date - Onset of Current Illness/Symptom; DTP01 431; PDF p.155
	InitialTreatmentDTP454             DateDTP    `json:"initialTreatmentDTP454,omitempty"`             // Date - Initial Treatment; DTP01 454; PDF p.155
	LastSeenDTP304                     DateDTP    `json:"lastSeenDTP304,omitempty"`                     // Date - Last Seen Date; DTP01 304; PDF p.155
	AcuteManifestationDTP453           DateDTP    `json:"acuteManifestationDTP453,omitempty"`           // Date - Acute Manifestation; DTP01 453; PDF p.155-156
	AccidentDTP439                     DateDTP    `json:"accidentDTP439,omitempty"`                     // Date - Accident; DTP01 439; PDF p.156
	LastMenstrualPeriodDTP484          DateDTP    `json:"lastMenstrualPeriodDTP484,omitempty"`          // Date - Last Menstrual Period; DTP01 484; PDF p.156
	LastXRayDTP455                     DateDTP    `json:"lastXRayDTP455,omitempty"`                     // Date - Last X-Ray Date; DTP01 455; PDF p.156
	HearingVisionPrescriptionDTP471    DateDTP    `json:"hearingVisionPrescriptionDTP471,omitempty"`    // Date - Hearing and Vision Prescription Date; DTP01 471; PDF p.156
	DisabilityDTP314                   DateDTP    `json:"disabilityDTP314,omitempty"`                   // Date - Disability Dates; DTP01 314, 360, or 361, D8 or RD8; one occurrence; PDF p.156-157
	LastWorkedDTP297                   DateDTP    `json:"lastWorkedDTP297,omitempty"`                   // Date - Last Worked; DTP01 297; PDF p.157
	AuthorizedReturnToWorkDTP296       DateDTP    `json:"authorizedReturnToWorkDTP296,omitempty"`       // Date - Authorized Return to Work; DTP01 296; PDF p.157
	AdmissionDTP435                    DateDTP    `json:"admissionDTP435,omitempty"`                    // Date - Admission; DTP01 435; PDF p.157
	DischargeDTP096                    DateDTP    `json:"dischargeDTP096,omitempty"`                    // Date - Discharge; DTP01 096; PDF p.157
	AssumedOrRelinquishedCareDTP090    [2]DateDTP `json:"assumedOrRelinquishedCareDTP090,omitempty"`    // Date - Assumed and Relinquished Care Dates; DTP01 090 or 091; repeats twice, filled in file order; PDF p.157-158
	PropertyCasualtyFirstContactDTP444 DateDTP    `json:"propertyCasualtyFirstContactDTP444,omitempty"` // Date - Property and Casualty Date of First Contact; DTP01 444; PDF p.158
	RepricerReceivedDTP050             DateDTP    `json:"repricerReceivedDTP050,omitempty"`             // Date - Repricer Received Date; DTP01 050; PDF p.158

	// PWK, "Claim Supplemental Information"; up to 10 occurrences.
	// CLAIM_PWK_NAME; PDF p.158. Medicare restricts PWK02 to BM/FX/EL/FT and
	// only the first occurrence (claim or line) is used in adjudication,
	// CEDI 837P CG, PDF p.38 (printed p.30; Table 21).
	SupplementalInfoPWK []AttachmentPWK `json:"supplementalInfoPWK,omitempty"`

	// CN1, "Contract Information"; not used by Medicare at either 2300 or
	// 2400 (NOT_USED_BY_MEDICARE). CN1_NAME; PDF p.159.
	ContractInfoCN1 ContractCN1 `json:"contractInfoCN1,omitempty"`

	// AMT, "Patient Amount Paid"; one fixed box, AMT01 F5 (Transmittal 579,
	// PDF p.32). The parser fills it from the first AMT*F5 that carries an
	// amount; any other AMT is kept in Raw. CLAIM_AMT_NAME; PDF p.159.
	PatientAmountPaidAMTF5 AmountAMT `json:"patientAmountPaidAMTF5,omitempty"`

	// REF (loop 2300), one field per qualifier box; CLAIM_REF_NAME /
	// CLAIM_REF02_EL_NAME in frontend/src/data/x12Names837p.ts, every box
	// Max=1 in the flat file. The REF01 code after each comment is the box's
	// only value in Transmittal 579's Values column; the parser fills a
	// field from the first REF with that code and a non-empty REF02.
	ServiceAuthExceptionREF4N           ReferenceREF `json:"serviceAuthExceptionREF4N,omitempty"`           // Service Authorization Exception Code; REF01 4N; PDF p.159
	MedicareCrossoverIndicatorREFF5     ReferenceREF `json:"medicareCrossoverIndicatorREFF5,omitempty"`     // Mandatory Medicare (Section 4081) Crossover Indicator; REF01 F5; not used by Medicare; PDF p.160
	MammographyCertificationREFEW       ReferenceREF `json:"mammographyCertificationREFEW,omitempty"`       // Mammography Certification Number; REF01 EW; PDF p.160
	ReferralREF9F                       ReferenceREF `json:"referralREF9F,omitempty"`                       // Referral Number; REF01 9F; PDF p.160-161
	PriorAuthorizationREFG1             ReferenceREF `json:"priorAuthorizationREFG1,omitempty"`             // Prior Authorization; REF01 G1; PDF p.161
	PayerClaimControlNumberREFF8        ReferenceREF `json:"payerClaimControlNumberREFF8,omitempty"`        // Payer Claim Control Number; REF01 F8; not used by Medicare; PDF p.161-162
	CLIAREFX4                           ReferenceREF `json:"cliaREFX4,omitempty"`                           // Clinical Laboratory Improvement Amendment (CLIA) Number; REF01 X4; PDF p.162
	RepricedClaimREF9A                  ReferenceREF `json:"repricedClaimREF9A,omitempty"`                  // Repriced Claim Number; REF01 9A; PDF p.162-163
	AdjustedRepricedClaimREF9C          ReferenceREF `json:"adjustedRepricedClaimREF9C,omitempty"`          // Adjusted Repriced Claim Number; REF01 9C; PDF p.163
	InvestigationalDeviceExemptionREFLX ReferenceREF `json:"investigationalDeviceExemptionREFLX,omitempty"` // Investigational Device Exemption Number; REF01 LX; PDF p.163-164
	TransmissionIntermediaryREFD9       ReferenceREF `json:"transmissionIntermediaryREFD9,omitempty"`       // Claim Identifier for Transmission Intermediaries; REF01 D9; PDF p.164
	MedicalRecordNumberREFEA            ReferenceREF `json:"medicalRecordNumberREFEA,omitempty"`            // Medical Record Number; REF01 EA; PDF p.164-165
	DemonstrationProjectREFP4           ReferenceREF `json:"demonstrationProjectREFP4,omitempty"`           // Demonstration Project Identifier; REF01 P4; PDF p.165
	CarePlanOversightREF1J              ReferenceREF `json:"carePlanOversightREF1J,omitempty"`              // Care Plan Oversight; REF01 1J; PDF p.165

	// References catches any loop 2300 REF not one of the qualified fields
	// above, a repeat of one, or one whose REF02 is empty.
	References []ReferenceREF `json:"references,omitempty"`

	// K3, "File Information"; each occurrence, up to 10. CLAIM_K3_NAME,
	// K3_EL_NAME; PDF p.166.
	FileInfoK3 []FileInfoK3 `json:"fileInfoK3,omitempty"`

	// NTE, "Claim Note". CLAIM_NTE_NAME; PDF p.166.
	NoteNTE NoteNTE `json:"noteNTE,omitempty"`

	// CR1, "Ambulance Transport Information". CLAIM_CR1_NAME; PDF p.166-167.
	AmbulanceTransportCR1 AmbulanceTransportCR1 `json:"ambulanceTransportCR1,omitempty"`

	// CR2, "Spinal Manipulation Service Information". CR2_NAME; PDF p.167.
	SpinalManipulationCR2 ClaimSpinalManipulationCR2 `json:"spinalManipulationCR2,omitempty"`

	// CRC (loop 2300), one field per box, placed by CRC01 (Transmittal 579's
	// codes: 07, E1-E3, 75, ZZ); CLAIM_CRC_NAME; PDF p.168-169. Ambulance
	// Certification and Patient Condition Information: Vision each repeat up
	// to 3 times; Homebound Indicator and EPSDT Referral occur at most once.
	AmbulanceCertificationCRC07 []ConditionCRC `json:"ambulanceCertificationCRC07,omitempty"` // CRC01 07; PDF p.168
	VisionConditionCRCE1        []ConditionCRC `json:"visionConditionCRCE1,omitempty"`        // CRC01 E1, E2, or E3; PDF p.168
	HomeboundCRC75              ConditionCRC   `json:"homeboundCRC75,omitempty"`              // CRC01 75; PDF p.168
	EPSDTReferralCRCZZ          ConditionCRC   `json:"epsdtReferralCRCZZ,omitempty"`          // CRC01 ZZ; PDF p.169

	// HI (loop 2300). Health Care Diagnosis Code, Anesthesia Related
	// Procedure, and Condition Information each repeat the same nine-position
	// composite (HealthCareCodeHI) up to 12 times (HI01-HI12); Medicare
	// populates only the first two positions of each composite (HI0x-3
	// through HI0x-9 Not Used), and the parser keeps all nine.
	// HI_DIAGNOSIS_NAME / HI_ANESTHESIA_NAME / HI_CONDITION_NAME; PDF
	// p.169-183. The parser places an HI segment in its box by HI01-1 (ABK
	// or BK, BP, BG; Transmittal 579, PDF p.42, 46, 51), keeping each
	// composite in position, an empty one included. The diagnosis and
	// anesthesia boxes occur once; the condition box occurs up to twice, so
	// it holds one segment per occurrence. Diagnosis codes (ICD) are kept
	// as plain text, same as every other code value in this package; this
	// is distinct from CPT/HCPCS procedure codes, which carry an
	// AMA-licensing restriction noted where those appear (service line, not
	// here).
	DiagnosisCodesHIABK      HealthCareInfoHI   `json:"diagnosisCodesHIABK,omitempty"`      // PDF p.169-174
	AnesthesiaProceduresHIBP HealthCareInfoHI   `json:"anesthesiaProceduresHIBP,omitempty"` // PDF p.174-178
	ConditionCodesHIBG       []HealthCareInfoHI `json:"conditionCodesHIBG,omitempty"`       // up to 2 segments; PDF p.178-183

	// HCP, "Claim Pricing/Repricing Information"; not used by Medicare
	// (NOT_USED_BY_MEDICARE). HCP_NAME; PDF p.181-182.
	PricingHCP PricingHCP `json:"pricingHCP,omitempty"`

	// Claim-level providers, loops 2310A-F. Each is its own X12 loop (its
	// own Loop ID in the flat file), so each carries its own Raw rather
	// than bubbling into this claim's Raw below. See LOOP_NM1_NAME,
	// NM1_ROLE_EL_NAME, LOOP_N3_NAME, LOOP_N4_NAME, LOOP_PRV_NAME,
	// LOOP_REF_NAME/LOOP_REF02_EL_NAME in frontend/src/data/x12Names837p.ts.
	ReferringProvider        ClaimReferringProviderNM1   `json:"referringProvider,omitempty"`        // 2310A; Transmittal 907, PDF p.184
	RenderingProviderNM182   ClaimRenderingProviderNM1   `json:"renderingProviderNM182,omitempty"`   // 2310B; PDF p.185-186
	ServiceFacilityNM177     ClaimServiceFacilityNM1     `json:"serviceFacilityNM177,omitempty"`     // 2310C; PDF p.186-188
	SupervisingProviderNM1DQ ClaimSupervisingProviderNM1 `json:"supervisingProviderNM1DQ,omitempty"` // 2310D; PDF p.188-189
	AmbulancePickupNM1PW     ClaimAmbulancePickupNM1     `json:"ambulancePickupNM1PW,omitempty"`     // 2310E; PDF p.189
	AmbulanceDropoffNM145    ClaimAmbulanceDropoffNM1    `json:"ambulanceDropoffNM145,omitempty"`    // 2310F; PDF p.189-190

	// OtherSubscribers, loop 2320 (coordination of benefits), one per other
	// payer; each carries its own 2330A-G NM1s. See other_subscriber.go.
	OtherSubscribers []OtherSubscriberSBR `json:"otherSubscribers,omitempty"`

	// ServiceLines, loop 2400, one per line of service.
	ServiceLines []ServiceLineLX `json:"serviceLines,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments directly within loop 2300 the model does not place
}

// The segment types below appear at both the claim (2300) and the service
// line (2400) with the same elements, so one type serves both loops; the
// field holding a value says which loop and box it came from.

// AttachmentPWK holds one PWK segment: "Claim Supplemental Information"
// (2300, up to 10), "Line Supplemental Information" (2400, up to 10), and
// the 2400 Durable Medical Equipment Certificate of Medical Necessity
// Indicator box. PWK03/04 are Not Used and kept as plain text; PWK07-09 are
// Not Used too and come after the last position the model names, so a file
// that sends them keeps their text in Extra. PWK05 has no confirmed code
// table in frontend/src/data/x12Names837p.ts and is kept as text.
// Transmittal 907, PDF p.158, 209; CLAIM_PWK_NAME, SERVICE_PWK_NAME,
// PWK_EL_NAME.
type AttachmentPWK struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ReportTypePWK01   AttachmentReportTypeCodePWK01   `json:"reportTypePWK01,omitempty"`   // Attachment Report Type Code
	TransmissionPWK02 AttachmentTransmissionCodePWK02 `json:"transmissionPWK02,omitempty"` // Attachment Transmission Code

	ReportCopiesPWK03 Amount `json:"reportCopiesPWK03,omitempty"` // Report Copies Needed (N/U) — Transmittal 907, PDF p.158
	EntityIDPWK04     string `json:"entityIDPWK04,omitempty"`     // Entity Identifier Code (N/U) — Transmittal 907, PDF p.158

	IDQualifierPWK05   string `json:"idQualifierPWK05,omitempty"`   // Identification Code Qualifier
	ControlNumberPWK06 string `json:"controlNumberPWK06,omitempty"` // Attachment Control Number

	Extra []string `json:"extra,omitempty"` // elements after PWK06, kept as sent
}

// ContractCN1 holds one CN1 segment, "Contract Information" (2300 and
// 2400). Not used by Medicare. Set is false when the segment did not appear
// at all. Transmittal 907, PDF p.159, 215; CN1_NAME / CN1_EL_NAME.
type ContractCN1 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	TypeCN101                 string `json:"typeCN101,omitempty"`                 // Contract Type Code; no confirmed value table in the reviewed sources
	AmountCN102               Amount `json:"amountCN102,omitempty"`               // Contract Amount
	PercentageCN103           Amount `json:"percentageCN103,omitempty"`           // Contract Percentage
	CodeCN104                 string `json:"codeCN104,omitempty"`                 // Contract Code
	TermsDiscountPercentCN105 Amount `json:"termsDiscountPercentCN105,omitempty"` // Terms Discount Percent
	VersionCN106              string `json:"versionCN106,omitempty"`              // Contract Version Identifier

	Extra []string `json:"extra,omitempty"` // elements after CN106, kept as sent
}

// NoteNTE holds one NTE segment: "Claim Note" (2300), "Line Note" and
// "Third Party Organization Notes" (2400). Set is false when the segment
// did not appear at all. Transmittal 907, PDF p.166, 220; CLAIM_NTE_NAME,
// LINE_NTE_NAME, THIRD_PARTY_NTE_NAME, NTE_EL_NAME.
type NoteNTE struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ReferenceCodeNTE01 string `json:"referenceCodeNTE01,omitempty"` // Note Reference Code; no confirmed value table in frontend/src/data/x12Names837p.ts
	TextNTE02          string `json:"textNTE02,omitempty"`          // Note Text

	Extra []string `json:"extra,omitempty"` // elements after NTE02, kept as sent
}

// AmbulanceTransportCR1 holds one CR1 segment, "Ambulance Transport
// Information" (2300 and 2400). CR103, CR107, and CR108 are Not Used and
// kept as plain text. Set is false when the segment did not appear at all.
// Transmittal 907, PDF p.166-167, 210; CLAIM_CR1_NAME, SERVICE_CR1_NAME,
// CR1_EL_NAME.
type AmbulanceTransportCR1 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	WeightUnitCR101    UnitOrBasisForMeasurementCode `json:"weightUnitCR101,omitempty"`    // Unit or Basis for Measurement Code
	WeightCR102        Amount                        `json:"weightCR102,omitempty"`        // Patient Weight
	TransportCodeCR103 string                        `json:"transportCodeCR103,omitempty"` // Ambulance Transport Code (N/U) — Transmittal 907, PDF p.166
	ReasonCR104        string                        `json:"reasonCR104,omitempty"`        // Ambulance Transport Reason Code; no confirmed value table in the reviewed sources
	DistanceUnitCR105  UnitOrBasisForMeasurementCode `json:"distanceUnitCR105,omitempty"`  // Unit or Basis for Measurement Code
	DistanceCR106      Amount                        `json:"distanceCR106,omitempty"`      // Transport Distance

	AddressInfoCR107 string `json:"addressInfoCR107,omitempty"` // Address Information (N/U) — Transmittal 907, PDF p.166
	AddressInfoCR108 string `json:"addressInfoCR108,omitempty"` // Address Information (N/U) — Transmittal 907, PDF p.166

	RoundTripPurposeCR109 string `json:"roundTripPurposeCR109,omitempty"` // Round Trip Purpose Description
	StretcherPurposeCR110 string `json:"stretcherPurposeCR110,omitempty"` // Stretcher Purpose Description

	Extra []string `json:"extra,omitempty"` // elements after CR110, kept as sent
}

// ClaimSpinalManipulationCR2 holds the CR2 segment (loop 2300, "Spinal
// Manipulation Service Information"). CR201-07 and CR209 are Not Used and
// kept as plain text; CR212 is after the last position the model names
// (CR2 has 11 elements), so a file that sends it keeps its text in Extra.
// Set is false when the segment did not appear at all. Transmittal 907,
// PDF p.167; CR2_NAME / CR2_EL_NAME.
type ClaimSpinalManipulationCR2 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	TreatmentSeriesCR201       Amount `json:"treatmentSeriesCR201,omitempty"`       // Treatment Series Number (N/U) — Transmittal 907, PDF p.167
	TreatmentCountCR202        Amount `json:"treatmentCountCR202,omitempty"`        // Treatment Count (N/U) — Transmittal 907, PDF p.167
	SubluxationLevelCR203      string `json:"subluxationLevelCR203,omitempty"`      // Subluxation Level Code (N/U) — Transmittal 907, PDF p.167
	SubluxationLevelCR204      string `json:"subluxationLevelCR204,omitempty"`      // Subluxation Level Code (N/U) — Transmittal 907, PDF p.167
	UnitCR205                  string `json:"unitCR205,omitempty"`                  // Unit or Basis for Measurement Code (N/U) — Transmittal 907, PDF p.167
	TreatmentPeriodCountCR206  Amount `json:"treatmentPeriodCountCR206,omitempty"`  // Treatment Period Count (N/U) — Transmittal 907, PDF p.167
	MonthlyTreatmentCountCR207 Amount `json:"monthlyTreatmentCountCR207,omitempty"` // Monthly Treatment Count (N/U) — Transmittal 907, PDF p.167

	ConditionCodeCR208 string `json:"conditionCodeCR208,omitempty"` // Patient Condition Code; no confirmed value table in the reviewed sources

	ComplicationIndicatorCR209 string `json:"complicationIndicatorCR209,omitempty"` // Complication Indicator (N/U) — Transmittal 907, PDF p.167

	ConditionDescriptionCR210 string `json:"conditionDescriptionCR210,omitempty"` // Patient Condition Description
	ConditionDescriptionCR211 string `json:"conditionDescriptionCR211,omitempty"` // Patient Condition Description

	Extra []string `json:"extra,omitempty"` // elements after CR211, kept as sent
}

// ConditionCRC holds one CRC segment, at 2300 (CLAIM_CRC_NAME: Ambulance
// Certification, Patient Condition Information: Vision, Homebound
// Indicator, EPSDT Referral) or 2400 (SERVICE_CRC_NAME: Ambulance
// Certification, Hospice Employee Indicator, Condition Indicator Durable
// Medical Equipment); the box is carried by the field holding this struct,
// not by a field here. CRC02 is the X12 1073 Yes/No code at every box,
// though some boxes title it differently (Certification Condition
// Indicator, Hospice Employed Provider Indicator). The number of condition
// positions a box populates varies (Transmittal 907 gives one shared
// element array, CRC_EL_NAME); unused slots stay empty. Set is false when
// the occurrence was not sent. Transmittal 907, PDF p.168-169, 210-211.
type ConditionCRC struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	CodeCategoryCRC01  CodeCategoryCRC01          `json:"codeCategoryCRC01,omitempty"`  // Code Category
	CertConditionCRC02 YesNoConditionResponseCode `json:"certConditionCRC02,omitempty"` // Certification Condition Indicator
	ConditionsCRC0x    [5]ConditionIndicatorCRC0x `json:"conditionsCRC0x,omitempty"`    // CRC03-CRC07, Condition Indicator/Condition Code

	Extra []string `json:"extra,omitempty"` // elements after CRC07, kept as sent
}

// HealthCareInfoHI holds one HI segment: up to twelve HealthCareCodeHI
// composites (HI01-HI12), each kept in its position, an empty one included.
// Set is false when the segment did not appear at all. A file that sends
// elements after HI12 keeps their text in Extra.
type HealthCareInfoHI struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	Codes []HealthCareCodeHI `json:"codes,omitempty"` // HI01 through HI12

	Extra []string `json:"extra,omitempty"` // elements after HI12, kept as sent
}

// HealthCareCodeHI is the nine-position composite HI (loop 2300) repeats up
// to 12 times (HI01-HI12) for each of its three box names - Health Care
// Diagnosis Code, Anesthesia Related Procedure, and Condition Information -
// which field this struct is stored in on ClaimCLM carries that
// distinction. Medicare populates only the first two positions (HI0x-3
// through HI0x-9 are Not Used in Transmittal 907's 837P), but the X12
// standard defines all nine, so each has a field and the parser keeps the
// text it finds there; a Not Used position is never an error (see Parse's
// doc comment), and a component after the ninth keeps its text in Extra.
// Transmittal 907 names all nine at every composite of all three boxes, PDF
// p.169-183; HI_DIAGNOSIS_EL_NAME / HI_ANESTHESIA_EL_NAME /
// HI_CONDITION_EL_NAME.
//
// The nine components and their types follow the 837I's HealthInfoHI, with
// two differences. This type has no Set field: HealthCareInfoHI keeps every
// composite in its position, an empty one included, so an unsent composite
// is the zero value. And HI0x-9 is the Yes/No Condition or Response Code
// the 837P flat file names, where the 837I's HI0n-9 is its Present on
// Admission Indicator.
//
// HI0x-4 is read as the format in HI0x-3 says: D8 is one date (DateHI0x4),
// RD8 a range (DateHI0x4 and DateThroughHI0x4). A date under any other
// format is an error, since no field holds it, as in the 837I. A malformed
// date, amount, or quantity is an error naming its component, such as
// "HI01-4".
type HealthCareCodeHI struct {
	TypeHI0x1 DiagnosisTypeCodeHI0x_1 `json:"typeHI0x1,omitempty"` // HI0x-1
	CodeHI0x2 string                  `json:"codeHI0x2,omitempty"` // HI0x-2

	DateQualifierHI0x3 string `json:"dateQualifierHI0x3,omitempty"` // Date Time Period Format Qualifier; Not Used in 005010A1; PDF p.169-183; D8 or RD8
	DateHI0x4          Date   `json:"dateHI0x4,omitempty"`          // Date Time Period; Not Used in 005010A1; PDF p.169-183; the D8 date, or an RD8 range's first date
	DateThroughHI0x4   Date   `json:"dateThroughHI0x4,omitempty"`   // an RD8 range's second date; Set is false for a single D8 date, as in the 837I's HealthInfoHI
	AmountHI0x5        Amount `json:"amountHI0x5,omitempty"`        // Monetary Amount; Not Used in 005010A1; PDF p.169-183
	QuantityHI0x6      Amount `json:"quantityHI0x6,omitempty"`      // Quantity; Not Used in 005010A1; PDF p.169-183

	VersionIdentifierHI0x7 string                     `json:"versionIdentifierHI0x7,omitempty"` // Version Identifier; Not Used in 005010A1; PDF p.169-183
	IndustryCodeHI0x8      string                     `json:"industryCodeHI0x8,omitempty"`      // Industry Code; Not Used in 005010A1; PDF p.169-183
	YesNoHI0x9             YesNoConditionResponseCode `json:"yesNoHI0x9,omitempty"`             // Yes/No Condition or Response Code; Not Used in 005010A1; PDF p.169-183

	Extra []string `json:"extra,omitempty"` // components after HI0x-9, kept as sent
}

// PricingHCP holds one HCP segment, "Claim Pricing/Repricing Information"
// (2300) or "Line Pricing/Repricing Information" (2400); not used by
// Medicare. HCP08 is Not Used and kept as plain text. Set is false when the
// segment did not appear at all. Transmittal 907, PDF p.184, 216-217;
// HCP_NAME / LINE_HCP_NAME / HCP_EL_NAME.
type PricingHCP struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	MethodologyHCP01       PricingMethodologyHCP01 `json:"methodologyHCP01,omitempty"`       // Pricing Methodology
	AllowedAmountHCP02     Amount                  `json:"allowedAmountHCP02,omitempty"`     // Repriced Allowed Amount
	SavingAmountHCP03      Amount                  `json:"savingAmountHCP03,omitempty"`      // Repriced Saving Amount
	OrgIDHCP04             string                  `json:"orgIDHCP04,omitempty"`             // Repricing Organization Identifier
	PerDiemOrFlatRateHCP05 Amount                  `json:"perDiemOrFlatRateHCP05,omitempty"` // Repricing Per Diem or Flat Rate Amount
	APGCodeHCP06           string                  `json:"apgCodeHCP06,omitempty"`           // Repriced Approved Ambulatory Patient Group Code
	APGAmountHCP07         Amount                  `json:"apgAmountHCP07,omitempty"`         // Repriced Approved Ambulatory Patient Group Amount

	ProductServiceIDHCP08 string `json:"productServiceIDHCP08,omitempty"` // Product/Service ID (N/U) — Transmittal 907, PDF p.184

	IDQualifierHCP09      string                        `json:"idQualifierHCP09,omitempty"`      // Product or Service ID Qualifier; no confirmed value table in the reviewed sources
	ApprovedHCPCSHCP10    string                        `json:"approvedHCPCSHCP10,omitempty"`    // Repriced Approved HCPCS Code
	UnitHCP11             UnitOrBasisForMeasurementCode `json:"unitHCP11,omitempty"`             // Unit or Basis for Measurement Code
	ApprovedUnitsHCP12    Amount                        `json:"approvedUnitsHCP12,omitempty"`    // Repriced Approved Service Unit Count
	RejectReasonHCP13     string                        `json:"rejectReasonHCP13,omitempty"`     // Reject Reason Code
	PolicyComplianceHCP14 string                        `json:"policyComplianceHCP14,omitempty"` // Policy Compliance Code
	ExceptionHCP15        string                        `json:"exceptionHCP15,omitempty"`        // Exception Code

	Extra []string `json:"extra,omitempty"` // elements after HCP15, kept as sent
}

// ClaimReferringProviderNM1 holds loop 2310A, "Referring Provider Name":
// NM1 plus up to 3 REF occurrences. Transmittal 907, PDF p.184-185.
type ClaimReferringProviderNM1 struct {
	Name       NameNM1        `json:"name,omitempty"`       // NM1; NM1_ROLE_EL_NAME["2310A"]
	References []ReferenceREF `json:"references,omitempty"` // REF, "Referring Provider Secondary Identification"; up to 3

	Raw []RawSegment `json:"raw,omitempty"`
}

// ClaimRenderingProviderNM1 holds loop 2310B, "Rendering Provider Name":
// NM1, PRV, and up to 4 REF occurrences. This is the claim-level rendering
// provider; the line-level rendering provider (loop 2420A, same box title)
// is a distinct struct owned by service_line.go. Transmittal 907, PDF
// p.185-186.
type ClaimRenderingProviderNM1 struct {
	Name       NameNM1        `json:"name,omitempty"`       // NM1; NM1_ROLE_EL_NAME["2310B"]
	Specialty  ProviderPRV    `json:"specialty,omitempty"`  // PRV, "Rendering Provider Specialty Information"
	References []ReferenceREF `json:"references,omitempty"` // REF, "Rendering Provider Secondary Identification"; up to 4

	Raw []RawSegment `json:"raw,omitempty"`
}

// ClaimServiceFacilityNM1 holds loop 2310C, "Service Facility Location
// Name": NM1, N3, N4, up to 3 REF occurrences, and PER. Transmittal 907,
// PDF p.186-188.
type ClaimServiceFacilityNM1 struct {
	Name       NameNM1        `json:"name,omitempty"`       // NM1; NM1_ROLE_EL_NAME["2310C"]
	Address    AddressN3      `json:"address,omitempty"`    // N3
	Location   LocationN4     `json:"location,omitempty"`   // N4
	References []ReferenceREF `json:"references,omitempty"` // REF, "Service Facility Location Secondary Identification"; up to 3
	Contact    ContactPER     `json:"contact,omitempty"`    // PER, SERVICE_FACILITY_PER_NAME, "Service Facility Contact Information"

	Raw []RawSegment `json:"raw,omitempty"`
}

// ClaimSupervisingProviderNM1 holds loop 2310D, "Supervising Provider
// Name": NM1 and up to 4 REF occurrences (no PRV at this loop).
// Transmittal 907, PDF p.188-189.
type ClaimSupervisingProviderNM1 struct {
	Name       NameNM1        `json:"name,omitempty"`       // NM1; NM1_ROLE_EL_NAME["2310D"]
	References []ReferenceREF `json:"references,omitempty"` // REF, "Supervising Provider Secondary Identifier"; up to 4

	Raw []RawSegment `json:"raw,omitempty"`
}

// ClaimAmbulancePickupNM1 holds loop 2310E, "Ambulance Pick Up Location":
// NM1, N3, N4 (no REF). Only NM101 (Entity Identifier Code) and NM103
// (AMBULANCE_PICKUP_NM1_EL_NAME, "Name Last or Organization Name") are
// populated; NM102 and NM104-NM112 are Not Used. Transmittal 907, PDF
// p.189.
type ClaimAmbulancePickupNM1 struct {
	Name     NameNM1    `json:"name,omitempty"`     // NM1
	Address  AddressN3  `json:"address,omitempty"`  // N3
	Location LocationN4 `json:"location,omitempty"` // N4

	Raw []RawSegment `json:"raw,omitempty"`
}

// ClaimAmbulanceDropoffNM1 holds loop 2310F, "Ambulance Drop Off
// Location": NM1, N3, N4 (no REF), same populated-position shape as
// ClaimAmbulancePickupNM1 (AMBULANCE_DROPOFF_NM1_EL_NAME). Transmittal
// 907, PDF p.189-190.
type ClaimAmbulanceDropoffNM1 struct {
	Name     NameNM1    `json:"name,omitempty"`     // NM1
	Address  AddressN3  `json:"address,omitempty"`  // N3
	Location LocationN4 `json:"location,omitempty"` // N4

	Raw []RawSegment `json:"raw,omitempty"`
}
