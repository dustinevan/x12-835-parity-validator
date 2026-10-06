package x12837p

// ProcedureSV101 is the Composite Medical Procedure Identifier: the shared
// shape of SV101 (loop 2400's Professional Service), SV501 (loop 2400's
// Durable Medical Equipment Service), and SVD03 (loop 2430's Line
// Adjudication Information). SV1_EL_NAME, SV5_EL_NAME, and SVD_EL_NAME in
// frontend/src/data/x12Names837p.ts all give this same qualifier/code/four
// modifiers/description shape; it is the identical composite the 835 side
// models as ProcedureSVC01 (backend/pkg/x12835/service_line.go), reused
// there for SVC01 and SVC06 the same way. Medicare requires the qualifier HC
// (Healthcare Common Procedure Coding System) at every one of these
// positions — SV101_1_NAME, SVD03_1_NAME; CEDI 837P CG, PDF p.42 (printed
// p.34; Table 33), p.49 (printed p.41; Table 46). The eighth component is
// Not Used, and a file that sends it keeps its text in Extra.
type ProcedureSV101 struct {
	Qualifier   string    `json:"qualifier,omitempty"`   // Product or Service ID Qualifier; Medicare requires HC. No enum type covers this code set at its several 837P positions (SV101-1, SV501-1, SVD03-1, LIN02, HCP09), so it is kept as text.
	Code        string    `json:"code,omitempty"`        // Procedure Code; CPT codes are kept as their code (AMA-licensed, see docs/design/835-viewer.md's equivalent rule on the 835 side)
	Modifiers   [4]string `json:"modifiers,omitempty"`   // up to four procedure modifiers
	Description string    `json:"description,omitempty"` // free-text description, sent in place of a code at the submitter's discretion

	Extra []string `json:"extra,omitempty"` // components after the seventh, kept as sent
}

// ProfessionalServiceSV1 holds the SV1 segment (loop 2400, Professional
// Service). Set is false when the line carried no SV1. SV106, SV110, SV113,
// and SV114 are Not Used and kept as plain text (SV108, "Monetary Amount,"
// is numeric, so it is kept as an Amount). SV105 is the Place of Service
// Code, an external CMS list the repository does not hold, so it stays
// plain text. Transmittal 907, PDF p.205-206 (footer 73-74).
type ProfessionalServiceSV1 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ProcedureSV101      ProcedureSV101                `json:"procedureSV101,omitempty"`      // SV101 composite; Medicare requires SV101-1 = HC — SV101_1_NAME; PDF p.205 (footer 73)
	ChargeSV102         Amount                        `json:"chargeSV102,omitempty"`         // Line Item Charge Amount; PDF p.205-206 (footer 73-74)
	UnitSV103           UnitOrBasisForMeasurementCode `json:"unitSV103,omitempty"`           // Unit or Basis for Measurement Code; PDF p.206 (footer 74)
	UnitCountSV104      Amount                        `json:"unitCountSV104,omitempty"`      // Service Unit Count; PDF p.206 (footer 74)
	PlaceOfServiceSV105 string                        `json:"placeOfServiceSV105,omitempty"` // Place of Service Code; the same external list as CLM05-1; PDF p.206 (footer 74)

	ServiceTypeSV106 string `json:"serviceTypeSV106,omitempty"` // Service Type Code (N/U) — Transmittal 907, PDF p.206 (footer 74)

	DiagnosisPointersSV107 [4]int   `json:"diagnosisPointersSV107,omitempty"` // SV107 composite, up to four 1-based pointers into the claim's HI diagnosis list (SV107-1 required, SV107-2/3/4 situational); 0 means that position was not sent, since a real pointer is never 0; PDF p.206 (footer 74)
	PointersExtraSV107     []string `json:"pointersExtraSV107,omitempty"`     // components after SV107-4, kept as sent

	MonetaryAmountSV108 Amount `json:"monetaryAmountSV108,omitempty"` // Monetary Amount (N/U) — Transmittal 907, PDF p.206 (footer 74)

	EmergencySV109 YesNoConditionResponseCode `json:"emergencySV109,omitempty"` // Emergency Indicator; PDF p.206 (footer 74)

	MultipleProcedureSV110 string `json:"multipleProcedureSV110,omitempty"` // Multiple Procedure Code (N/U) — Transmittal 907, PDF p.206 (footer 74)

	EPSDTSV111          YesNoConditionResponseCode `json:"epsdtSV111,omitempty"`          // EPSDT Indicator; PDF p.206 (footer 74)
	FamilyPlanningSV112 YesNoConditionResponseCode `json:"familyPlanningSV112,omitempty"` // Family Planning Indicator; PDF p.206 (footer 74)

	ReviewCodeSV113  string `json:"reviewCodeSV113,omitempty"`  // Review Code (N/U) — Transmittal 907, PDF p.206 (footer 74)
	ReviewValueSV114 string `json:"reviewValueSV114,omitempty"` // National or Local Assigned Review Value (N/U) — Transmittal 907, PDF p.206 (footer 74)

	CoPayStatusSV115 string `json:"coPayStatusSV115,omitempty"` // Co-Pay Status Code; external code list, not confirmed in reviewed sources; PDF p.206 (footer 74)

	Extra []string `json:"extra,omitempty"` // elements after SV115, kept as sent
}

// LineDMEServiceSV5 holds the SV5 segment (loop 2400, Durable Medical
// Equipment Service): an addition to SV1 carried only on a DME line. Set is
// false when the segment did not appear at all. Transmittal 907, PDF
// p.208-209 (footer 76-77); SV5_NAME, SV5_EL_NAME.
type LineDMEServiceSV5 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ProcedureSV501                ProcedureSV101                `json:"procedureSV501,omitempty"`                // SV501 composite, same shape as SV101
	UnitSV502                     UnitOrBasisForMeasurementCode `json:"unitSV502,omitempty"`                     // Unit or Basis for Measurement Code
	LengthOfNecessitySV503        Amount                        `json:"lengthOfNecessitySV503,omitempty"`        // Length of Medical Necessity
	RentalPriceSV504              Amount                        `json:"rentalPriceSV504,omitempty"`              // DME Rental Price
	PurchasePriceSV505            Amount                        `json:"purchasePriceSV505,omitempty"`            // DME Purchase Price
	RentalUnitPriceIndicatorSV506 string                        `json:"rentalUnitPriceIndicatorSV506,omitempty"` // Rental Unit Price Indicator; external code list, not confirmed in reviewed sources

	Extra []string `json:"extra,omitempty"` // elements after SV506, kept as sent
}

// LineDMECertificationCR3 holds the CR3 segment (loop 2400, Durable Medical
// Equipment Certification). The CEDI 837P Companion Guide's own
// Medicare-restricted tables dropped this box as of version 5.6 (January
// 2025; PDF p.55, printed p.48; Table 50 change log), but it remains part of
// the flat file layout — CR3_NAME, CR3_EL_NAME; Transmittal 907, PDF p.210
// (footer 78). Set is false when the segment did not appear at all.
type LineDMECertificationCR3 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	CertificationTypeCR301 string                        `json:"certificationTypeCR301,omitempty"` // Certification Type Code; external code list, not confirmed in reviewed sources
	UnitCR302              UnitOrBasisForMeasurementCode `json:"unitCR302,omitempty"`              // Unit or Basis for Measurement Code
	DurationCR303          Amount                        `json:"durationCR303,omitempty"`          // Durable Medical Equipment Duration

	Extra []string `json:"extra,omitempty"` // elements after CR303, kept as sent
}

// TestResultMEA holds one MEA segment (loop 2400, Test Results). Repeats up
// to 5 times. MEA01 and MEA02 are coded elements with no publicly confirmed
// value table in the reviewed sources and no assigned enum type, so both are
// kept as text — MEA_NAME, MEA_EL_NAME; Transmittal 907, PDF p.214-215
// (footer 82-83).
type TestResultMEA struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ReferenceIDMEA01 string `json:"referenceIDMEA01,omitempty"` // Measurement Reference Identification Code; external code list, not confirmed in reviewed sources
	QualifierMEA02   string `json:"qualifierMEA02,omitempty"`   // Measurement Qualifier; external code list, not confirmed in reviewed sources
	ResultMEA03      Amount `json:"resultMEA03,omitempty"`      // Test Result

	Extra []string `json:"extra,omitempty"` // elements after MEA03, kept as sent
}

// PurchasedServicePS1 holds the PS1 segment (loop 2400, Purchased Service
// Information): who performed a purchased service and what it cost.
// Transmittal 907, PDF p.220-221 (footer 88-89); PS1_NAME, PS1_EL_NAME. Set
// is false when the segment did not appear at all.
type PurchasedServicePS1 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ProviderIDPS101   string `json:"providerIDPS101,omitempty"`   // Purchased Service Provider Identifier
	ChargeAmountPS102 Amount `json:"chargeAmountPS102,omitempty"` // Purchased Service Charge Amount

	Extra []string `json:"extra,omitempty"` // elements after PS102, kept as sent
}

// DrugPricingCTP holds the CTP segment of loop 2410 (Drug Pricing). Set is
// false when the loop carried no CTP. CTP01-03 are Not Used and kept as
// plain text (CTP03, "Unit Price," is numeric, so it is kept as an Amount).
// CTP05 is a composite whose first component is the unit; a file that
// sends more components keeps their text in UnitCTP05Extra. Transmittal
// 907, PDF p.223-224 (footer 91-92).
type DrugPricingCTP struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ClassOfTradeCTP01    string `json:"classOfTradeCTP01,omitempty"`    // Class of Trade Code (N/U) — Transmittal 907, PDF p.223 (footer 91)
	PriceIdentifierCTP02 string `json:"priceIdentifierCTP02,omitempty"` // Price Identifier Code (N/U) — Transmittal 907, PDF p.223 (footer 91)
	UnitPriceCTP03       Amount `json:"unitPriceCTP03,omitempty"`       // Unit Price (N/U) — Transmittal 907, PDF p.223 (footer 91)

	UnitCountCTP04 Amount                        `json:"unitCountCTP04,omitempty"` // National Drug Unit Count; PDF p.223 (footer 91)
	UnitCTP05_1    UnitOrBasisForMeasurementCode `json:"unitCTP05_1,omitempty"`    // Unit or Basis for Measurement Code (CTP05, first component); required UN for an Oral Cancer Drug NDC — CTP05_1_NAME; PDF p.223-224 (footer 91-92)
	UnitCTP05Extra []string                      `json:"unitCTP05Extra,omitempty"` // components after CTP05-1, kept as sent

	Extra []string `json:"extra,omitempty"` // elements after CTP05, kept as sent
}

// DrugLIN holds loop 2410 (Drug Identification): LIN, CTP, and REF. Present
// when the service line reports an NDC; required for a DME claim whenever
// loop 2400 SV101-1 is a default HCPCS code (CEDI 837P CG, PDF p.46, printed
// p.38, Table 39). Transmittal 907, PDF p.222-224 (footer 90-92). Set is
// false when the loop did not appear at all. LIN01 is Not Used and kept as
// plain text; LIN04 onward keep their text in Extra.
type DrugLIN struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the LIN's position in the file, 0-based; the ISA is 0

	AssignedIDLIN01 string `json:"assignedIDLIN01,omitempty"` // Assigned Identification (N/U) — Transmittal 907, PDF p.222 (footer 90)
	QualifierLIN02  string `json:"qualifierLIN02,omitempty"`  // Product/Service ID Qualifier; Medicare requires N4, National Drug Code — LIN02_NAME; PDF p.222 (footer 90)
	NDCLIN03        string `json:"ndcLIN03,omitempty"`        // National Drug Code; PDF p.222 (footer 90)

	Extra []string `json:"extra,omitempty"` // elements after LIN03, kept as sent

	Pricing DrugPricingCTP `json:"pricing,omitempty"` // CTP

	AssociationREFXZ ReferenceREF `json:"associationREFXZ,omitempty"` // REF, Prescription or Compound Drug Association Number; REF01 = XZ required with SV1 modifier J1 — DRUG_ASSOCIATION_REF01_NAME, DRUG_ASSOCIATION_REF_NAME; PDF p.224 (footer 92)

	Raw []RawSegment `json:"raw,omitempty"` // segments within this drug loop the model does not place
}

// LineProviderNM1 holds one of the eight 2420A-H line-level provider loops:
// an NM1, with PRV, N3/N4, PER, and REF present only for the roles that
// carry them — LOOP_NM1_NAME, LOOP_PRV_NAME, LOOP_N3_NAME, LOOP_N4_NAME,
// LOOP_REF_NAME, ORDERING_PROVIDER_PER_NAME in
// frontend/src/data/x12Names837p.ts. The field holding this struct
// (RenderingNM182, PurchasedServiceNM1QB, ... AmbulanceDropoffNM145
// on ServiceLineLX below) carries which role and which of those boxes
// apply, not a field here — the same pattern x12835's PartyN1 uses for its
// payer/payee-only DeliveryRDM. Set is false when the loop did not appear at
// all.
type LineProviderNM1 struct {
	Set bool `json:"set,omitempty"`

	Name NameNM1 `json:"name,omitempty"` // NM1

	Specialty ProviderPRV `json:"specialty,omitempty"` // PRV; 2420A (Rendering) only — LOOP_PRV_NAME; PDF p.225 (footer 93)

	Address  AddressN3  `json:"address,omitempty"`  // N3; 2420C, 2420E, 2420G, 2420H only — LOOP_N3_NAME
	Location LocationN4 `json:"location,omitempty"` // N4; same four loops — LOOP_N4_NAME

	Contact ContactPER `json:"contact,omitempty"` // PER; 2420E (Ordering) only — ORDERING_PROVIDER_PER_NAME; PDF p.231 (footer 99)

	References []ReferenceREF `json:"references,omitempty"` // REF, secondary identification; 2420A/B/C/D/E/F only, not 2420G/H — LOOP_REF_NAME; repeats up to 20 times (3 for 2420C)

	Raw []RawSegment `json:"raw,omitempty"` // segments within this provider loop the model does not place
}

// LineAdjudicationSVD holds one SVD segment (loop 2430) and everything
// nested under it: how another payer adjudicated this service line, as
// reported by the submitter for coordination of benefits. One occurrence
// per other payer; CAS, DTP, and AMT belong to that same occurrence.
// Transmittal 907, PDF p.235-237 (footer 103-105); SVD_NAME, SVD_EL_NAME.
type LineAdjudicationSVD struct {
	Index int `json:"index,omitempty"` // the SVD's position in the file, 0-based; the ISA is 0

	OtherPayerSVD01 string `json:"otherPayerSVD01,omitempty"` // Other Payer Primary Identifier
	PaidAmountSVD02 Amount `json:"paidAmountSVD02,omitempty"` // Service Line Paid Amount

	ProcedureSVD03 ProcedureSV101 `json:"procedureSVD03,omitempty"` // SVD03 composite, same shape as SV101; Medicare requires SVD03-1 = HC, and rejects ER, IV, or WK — SVD03_1_NAME; CEDI 837P CG, PDF p.49 (printed p.41; Table 46)

	ProductServiceIDSVD04 string `json:"productServiceIDSVD04,omitempty"` // Product/Service ID (N/U) — Transmittal 907, PDF p.235 (footer 103)

	PaidUnitsSVD05   Amount `json:"paidUnitsSVD05,omitempty"`   // Paid Service Unit Count
	BundledLineSVD06 string `json:"bundledLineSVD06,omitempty"` // Bundled or Unbundled Line Number; kept as a string identifier even though numeric, matching this service line's own AssignedNumberLX01, which it cross-references

	Extra []string `json:"extra,omitempty"` // elements after SVD06, kept as sent

	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"` // CAS, "Line Adjustment" (LINE_ADJUSTMENT_CAS_NAME); repeats up to 5 times; PDF p.235 (footer 103)

	AdjudicationDateDTP573 DateDTP `json:"adjudicationDateDTP573,omitempty"` // DTP, "Line Check or Remittance Date"; DTP03 is the Adjudication or Payment Date — LINE_ADJUDICATION_DTP_NAME, LINE_ADJUDICATION_DTP03_EL_NAME; PDF p.237 (footer 105)

	RemainingLiabilityAMTEAF AmountAMT `json:"remainingLiabilityAMTEAF,omitempty"` // AMT, "Remaining Patient Liability" — LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME; PDF p.237 (footer 105)

	Raw []RawSegment `json:"raw,omitempty"` // segments within this adjudication occurrence the model does not place
}

// FormAnswerFRM is one FRM segment within a FormLQ: one question's number
// and its response. FRM03 carries a text response; FRM04 and FRM05 occupy
// the same response position for a date or a numeric answer respectively,
// per Transmittal 907's own composite redefinitions of Question Response.
// Transmittal 907, PDF p.238 (footer 106); FRM_EL_NAME.
type FormAnswerFRM struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	QuestionFRM01 string `json:"questionFRM01,omitempty"` // Question Number/Letter
	ResponseFRM02 string `json:"responseFRM02,omitempty"` // Question Response; external code list, not confirmed in reviewed sources
	ResponseFRM03 string `json:"responseFRM03,omitempty"` // Question Response, text
	ResponseFRM04 Date   `json:"responseFRM04,omitempty"` // Question Response, date
	ResponseFRM05 Amount `json:"responseFRM05,omitempty"` // Question Response, numeric

	Extra []string `json:"extra,omitempty"` // elements after FRM05, kept as sent
}

// FormLQ holds one LQ segment (loop 2440, Supporting Documentation: a form
// or clinical questionnaire) and its FRM answers. Transmittal 907, PDF
// p.238 (footer 106); LQ_NAME, LQ_EL_NAME.
type FormLQ struct {
	Index int `json:"index,omitempty"` // the LQ's position in the file, 0-based; the ISA is 0

	QualifierLQ01 string `json:"qualifierLQ01,omitempty"` // Code List Qualifier Code; external code list, not confirmed in reviewed sources
	FormIDLQ02    string `json:"formIDLQ02,omitempty"`    // Form Identifier

	Extra []string `json:"extra,omitempty"` // elements after LQ02, kept as sent

	Answers []FormAnswerFRM `json:"answers,omitempty"` // FRM, up to 99 per form

	Raw []RawSegment `json:"raw,omitempty"` // segments within this form the model does not place
}

// ServiceLineLX holds one LX segment (loop 2400) and everything the 837P
// nests under it: the professional service itself, its optional DME
// addition, supplemental and certification information, line-level dates,
// quantities, references, amounts, notes, pricing, the drug loop (2410),
// the eight line-provider loops (2420A-H), other-payer line adjudication
// (2430), and supporting documentation (2440). Transmittal 907, PDF
// p.205-238 (footer 73-106).
type ServiceLineLX struct {
	Index int `json:"index,omitempty"` // the LX's position in the file, 0-based; the ISA is 0

	AssignedNumberLX01 string `json:"assignedNumberLX01,omitempty"` // Assigned Number; this service line's own 1-based position within the claim, referenced back by SVD06 and HCP; kept as a string identifier, like every identifier in this package, even though X12 types it numeric (N0) — LX_NAME, LX01_EL_NAME; PDF p.205 (footer 73)

	Extra []string `json:"extra,omitempty"` // elements after LX01, kept as sent

	Service ProfessionalServiceSV1 `json:"service,omitempty"` // SV1, Professional Service; PDF p.205-206 (footer 73-74)

	DME LineDMEServiceSV5 `json:"dme,omitempty"` // SV5, the DME alternative/addition to SV1, present only on a DME claim line — SV5_NAME; PDF p.208-209 (footer 76-77)

	Attachments []AttachmentPWK `json:"attachments,omitempty"` // PWK; "Line Supplemental Information" plus the DME Certificate of Medical Necessity Indicator box (PWK01 CT, Transmittal 579, PDF p.82), in file order — SERVICE_PWK_NAME; PDF p.209 (footer 77)

	AmbulanceTransport AmbulanceTransportCR1   `json:"ambulanceTransport,omitempty"` // CR1 — SERVICE_CR1_NAME; PDF p.210 (footer 78)
	DMECertification   LineDMECertificationCR3 `json:"dmeCertification,omitempty"`   // CR3, dropped from the CEDI companion guide's own Medicare-restricted tables as of v5.6 but still part of the flat file layout — CR3_NAME; PDF p.210 (footer 78)

	AmbulanceCertificationCRC07 []ConditionCRC `json:"ambulanceCertificationCRC07,omitempty"` // CRC01 07, repeats up to 3 — SERVICE_CRC_NAME.ambulanceCertification; PDF p.210-211 (footer 78-79)
	HospiceEmployeeCRC70        ConditionCRC   `json:"hospiceEmployeeCRC70,omitempty"`        // CRC01 70 — SERVICE_CRC_NAME.hospiceEmployeeIndicator; PDF p.211 (footer 79)
	DMEConditionCRC09           ConditionCRC   `json:"dmeConditionCRC09,omitempty"`           // CRC01 09, dropped from the CEDI companion guide's own Medicare-restricted tables as of v5.6 — SERVICE_CRC_NAME.dmeConditionIndicator; PDF p.211 (footer 79)

	// Dates (loop 2400). SERVICE_DTP_NAME in frontend/src/data/x12Names837p.ts
	// names each box by its title. The DTP01 code after each comment is the
	// box's value in Transmittal 579's Values column; the parser places a DTP
	// in the field for its code, and keeps a DTP whose box is already filled,
	// or whose DTP01 matches no box, in Raw.
	ServiceDateDTP472            DateDTP   `json:"serviceDateDTP472,omitempty"`            // DTP01 472; required on every service line — SERVICE_DTP_NAME.serviceDate; PDF p.211-212 (footer 79-80)
	PrescriptionDTP471           DateDTP   `json:"prescriptionDTP471,omitempty"`           // DTP01 471 — SERVICE_DTP_NAME.prescription; PDF p.212 (footer 80)
	CertificationRevisionDTP607  DateDTP   `json:"certificationRevisionDTP607,omitempty"`  // DTP01 607; CEDI v5.6-dropped, same as CR3 above — SERVICE_DTP_NAME.certificationRevision; PDF p.212 (footer 80)
	BeginTherapyDTP463           DateDTP   `json:"beginTherapyDTP463,omitempty"`           // DTP01 463; CEDI v5.6-dropped — SERVICE_DTP_NAME.beginTherapy; PDF p.212 (footer 80)
	LastCertificationDTP461      DateDTP   `json:"lastCertificationDTP461,omitempty"`      // DTP01 461; CEDI v5.6-dropped — SERVICE_DTP_NAME.lastCertification; PDF p.212 (footer 80)
	LastTreatmentOrTherapyDTP304 DateDTP   `json:"lastTreatmentOrTherapyDTP304,omitempty"` // DTP01 304 — SERVICE_DTP_NAME.lastTreatmentOrTherapy; PDF p.212-213 (footer 80-81)
	TestPerformedDTP738          []DateDTP `json:"testPerformedDTP738,omitempty"`          // DTP01 738 or 739; repeats up to 2 — SERVICE_DTP_NAME.testPerformed; PDF p.213 (footer 81)
	ShippedDTP011                DateDTP   `json:"shippedDTP011,omitempty"`                // DTP01 011 — SERVICE_DTP_NAME.shipped; PDF p.213 (footer 81)
	LastXRayDTP455               DateDTP   `json:"lastXRayDTP455,omitempty"`               // DTP01 455 — SERVICE_DTP_NAME.lastXRay; PDF p.213 (footer 81)
	InitialTreatmentDTP454       DateDTP   `json:"initialTreatmentDTP454,omitempty"`       // DTP01 454 — SERVICE_DTP_NAME.initialTreatment; PDF p.213 (footer 81)

	AmbulancePatientCountQTYPT    QuantityQTY `json:"ambulancePatientCountQTYPT,omitempty"`    // QTY01 PT — AMBULANCE_PATIENT_COUNT_QTY_NAME; PDF p.213-214 (footer 81-82)
	ObstetricAdditionalUnitsQTYFL QuantityQTY `json:"obstetricAdditionalUnitsQTYFL,omitempty"` // QTY01 FL — OBSTETRIC_ANESTHESIA_QTY_NAME; PDF p.214 (footer 82)

	TestResults []TestResultMEA `json:"testResults,omitempty"` // MEA, repeats up to 5 — MEA_NAME; PDF p.214-215 (footer 82-83)

	Contract ContractCN1 `json:"contract,omitempty"` // CN1; not used by Medicare — CN1_NAME; PDF p.215 (footer 83)

	// REF (loop 2400). SERVICE_REF_NAME names each qualifier individually.
	// The REF01 code after each comment is the box's only value in
	// Transmittal 579's Values column; the parser places a REF in the field
	// for its code, appends every REF*G1 and REF*9F, and keeps a repeat of a
	// single box, or a REF whose REF01 matches no box, in Raw. The parser
	// fills LineControlREF6R from the first REF*6R with a non-empty REF02.
	RepricedLineItemREF9B         ReferenceREF   `json:"repricedLineItemREF9B,omitempty"`         // REF01 9B — SERVICE_REF_NAME.repricedLineItem; PDF p.215-216 (footer 83-84)
	AdjustedRepricedLineItemREF9D ReferenceREF   `json:"adjustedRepricedLineItemREF9D,omitempty"` // REF01 9D — SERVICE_REF_NAME.adjustedRepricedLineItem; PDF p.216 (footer 84)
	PriorAuthorizationREFG1       []ReferenceREF `json:"priorAuthorizationREFG1,omitempty"`       // REF01 G1; repeats up to 5 — SERVICE_REF_NAME.priorAuthorization; PDF p.216 (footer 84)
	LineControlREF6R              ReferenceREF   `json:"lineControlREF6R,omitempty"`              // Line Item Control Number; REF01 = 6R confirmed, CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "REF*6R*001H0000000140309201") — SERVICE_REF_NAME.lineItemControl; PDF p.217 (footer 85)
	MammographyCertificationREFEW ReferenceREF   `json:"mammographyCertificationREFEW,omitempty"` // REF01 EW — SERVICE_REF_NAME.mammographyCertification; PDF p.217 (footer 85)
	CLIAREFX4                     ReferenceREF   `json:"cliaREFX4,omitempty"`                     // REF01 X4 — SERVICE_REF_NAME.clia; PDF p.218 (footer 86)
	ReferringCLIAREFF4            ReferenceREF   `json:"referringCLIAREFF4,omitempty"`            // REF01 F4 — SERVICE_REF_NAME.referringClia; PDF p.218 (footer 86)
	ImmunizationBatchREFBT        ReferenceREF   `json:"immunizationBatchREFBT,omitempty"`        // REF01 BT — SERVICE_REF_NAME.immunizationBatch; PDF p.219 (footer 87)
	ReferralREF9F                 []ReferenceREF `json:"referralREF9F,omitempty"`                 // REF01 9F; repeats up to 5 — SERVICE_REF_NAME.referral; PDF p.219 (footer 87)

	SalesTaxAMTT        AmountAMT `json:"salesTaxAMTT,omitempty"`        // AMT01 T — SALES_TAX_AMT_NAME; PDF p.219-220 (footer 87-88)
	PostageClaimedAMTF4 AmountAMT `json:"postageClaimedAMTF4,omitempty"` // AMT01 F4 — POSTAGE_CLAIMED_AMT_NAME; PDF p.220 (footer 88)

	FileInfoK3 []FileInfoK3 `json:"fileInfoK3,omitempty"` // K3, free-text fixed-format records, repeats up to 10 — SERVICE_K3_NAME, K3_EL_NAME; PDF p.220 (footer 88)

	// The two NTE boxes are told apart by NTE01: TPO is the third-party box's
	// only value, ADD and DCP the line note's (Transmittal 579, PDF p.93).
	// The parser sends an NTE with NTE01 TPO to ThirdPartyNote and any other
	// to LineNote.
	LineNote       NoteNTE `json:"lineNote,omitempty"`       // LINE_NTE_NAME; PDF p.220 (footer 88)
	ThirdPartyNote NoteNTE `json:"thirdPartyNote,omitempty"` // THIRD_PARTY_NTE_NAME; PDF p.220 (footer 88)

	PurchasedService PurchasedServicePS1 `json:"purchasedService,omitempty"` // PS1_NAME; PDF p.220-221 (footer 88-89)

	Pricing PricingHCP `json:"pricing,omitempty"` // HCP; not used by Medicare — LINE_HCP_NAME; PDF p.221 (footer 89)

	Drug DrugLIN `json:"drug,omitempty"` // loop 2410: LIN, CTP, REF; PDF p.222-224 (footer 90-92)

	// Line providers (loop 2420A-H). LOOP_NM1_NAME names each box; 2420A-H
	// "share those same tables with their 2310-series counterparts"
	// (frontend/src/data/x12Names837p.ts's own note above LOOP_REF_NAME),
	// except 2420B, 2420E, and 2420F, which have no 2310-series analog.
	RenderingNM182        LineProviderNM1 `json:"renderingNM182,omitempty"`        // PRV + REF x20 — PDF p.225-226 (footer 93-94)
	PurchasedServiceNM1QB LineProviderNM1 `json:"purchasedServiceNM1QB,omitempty"` // REF x20; no 2310-series analog — PDF p.226 (footer 94)
	ServiceFacilityNM177  LineProviderNM1 `json:"serviceFacilityNM177,omitempty"`  // N3/N4 + REF x3 — PDF p.227-228 (footer 95-96)
	SupervisingNM1DQ      LineProviderNM1 `json:"supervisingNM1DQ,omitempty"`      // REF x20 — PDF p.229-230 (footer 97-98)
	OrderingNM1DK         LineProviderNM1 `json:"orderingNM1DK,omitempty"`         // N3/N4 + REF x20 + PER; no 2310-series analog — PDF p.230-231 (footer 98-99)
	ReferringNM12420F     LineProviderNM1 `json:"referringNM12420F,omitempty"`     // REF x20; no 2310-series analog — PDF p.231-232 (footer 99-100)
	AmbulancePickupNM1PW  LineProviderNM1 `json:"ambulancePickupNM1PW,omitempty"`  // N3/N4, no REF — PDF p.232-234 (footer 100-102)
	AmbulanceDropoffNM145 LineProviderNM1 `json:"ambulanceDropoffNM145,omitempty"` // N3/N4, no REF — PDF p.233-234 (footer 101-102)

	Adjudications []LineAdjudicationSVD `json:"adjudications,omitempty"` // loop 2430: SVD, CAS, DTP, AMT — one occurrence per other payer that adjudicated this line; PDF p.235-237 (footer 103-105)

	Forms []FormLQ `json:"forms,omitempty"` // loop 2440: LQ, FRM; PDF p.238 (footer 106)

	Raw []RawSegment `json:"raw,omitempty"` // segments within this service line the model does not place
}
