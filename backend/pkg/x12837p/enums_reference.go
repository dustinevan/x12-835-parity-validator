package x12837p

// Generic cross-loop qualifiers for the shared segment types every loop
// reuses: REF (Reference), DTP (Date), AMT (Amount), and QTY (Quantity).
// Each qualifier type is shared by every instance of its segment in the
// transaction, not owned by one loop.

// ReferenceIdentificationQualifierREF01 is REF01, the Reference Identification
// Qualifier, shared by every REF segment in the transaction. PRV02 shares
// this type (ProviderPRV's QualifierPRV02); PXC is its one value.
//
// The list is the union of the REF01 codes the 837P implementation guide
// allows across its REF boxes. Transmittal 907 names every REF box's title
// but, like NM101, rarely restates REF01's own code character
// (frontend/src/data/x12Names837p.ts's header note). The first three values
// below are confirmed in the CEDI 837P Companion Guide's own worked
// transaction or a Medicare-restriction table. Most of the rest come from
// Transmittal 579, the original 005010 837P flat file, whose Values column
// prints REF01's code at every REF box; the parser places each REF in its
// box's field by this code. A code that is the only value of its box(es)
// takes the box title as its short name; a code listed beside others in a
// secondary-identification box (SY, 0B, 1G, 2U, FY, NF, Y4, G2, LU, VY) takes
// the guide's own name for the code.
type ReferenceIdentificationQualifierREF01 string

const (
	ReferenceIdentificationQualifierREF01NotSent ReferenceIdentificationQualifierREF01 = ""
	ReferenceIdentificationQualifierREF01EI      ReferenceIdentificationQualifierREF01 = "EI"
	ReferenceIdentificationQualifierREF016R      ReferenceIdentificationQualifierREF01 = "6R"
	ReferenceIdentificationQualifierREF01XZ      ReferenceIdentificationQualifierREF01 = "XZ"
	ReferenceIdentificationQualifierREF01SY      ReferenceIdentificationQualifierREF01 = "SY"
	ReferenceIdentificationQualifierREF010B      ReferenceIdentificationQualifierREF01 = "0B"
	ReferenceIdentificationQualifierREF011G      ReferenceIdentificationQualifierREF01 = "1G"
	ReferenceIdentificationQualifierREF012U      ReferenceIdentificationQualifierREF01 = "2U"
	ReferenceIdentificationQualifierREF01FY      ReferenceIdentificationQualifierREF01 = "FY"
	ReferenceIdentificationQualifierREF01NF      ReferenceIdentificationQualifierREF01 = "NF"
	ReferenceIdentificationQualifierREF01Y4      ReferenceIdentificationQualifierREF01 = "Y4"
	ReferenceIdentificationQualifierREF011W      ReferenceIdentificationQualifierREF01 = "1W"
	ReferenceIdentificationQualifierREF01G2      ReferenceIdentificationQualifierREF01 = "G2"
	ReferenceIdentificationQualifierREF01LU      ReferenceIdentificationQualifierREF01 = "LU"
	ReferenceIdentificationQualifierREF014N      ReferenceIdentificationQualifierREF01 = "4N"
	ReferenceIdentificationQualifierREF01F5      ReferenceIdentificationQualifierREF01 = "F5"
	ReferenceIdentificationQualifierREF01EW      ReferenceIdentificationQualifierREF01 = "EW"
	ReferenceIdentificationQualifierREF019F      ReferenceIdentificationQualifierREF01 = "9F"
	ReferenceIdentificationQualifierREF01G1      ReferenceIdentificationQualifierREF01 = "G1"
	ReferenceIdentificationQualifierREF01F8      ReferenceIdentificationQualifierREF01 = "F8"
	ReferenceIdentificationQualifierREF01X4      ReferenceIdentificationQualifierREF01 = "X4"
	ReferenceIdentificationQualifierREF019A      ReferenceIdentificationQualifierREF01 = "9A"
	ReferenceIdentificationQualifierREF019C      ReferenceIdentificationQualifierREF01 = "9C"
	ReferenceIdentificationQualifierREF01LX      ReferenceIdentificationQualifierREF01 = "LX"
	ReferenceIdentificationQualifierREF01D9      ReferenceIdentificationQualifierREF01 = "D9"
	ReferenceIdentificationQualifierREF01EA      ReferenceIdentificationQualifierREF01 = "EA"
	ReferenceIdentificationQualifierREF01P4      ReferenceIdentificationQualifierREF01 = "P4"
	ReferenceIdentificationQualifierREF011J      ReferenceIdentificationQualifierREF01 = "1J"
	ReferenceIdentificationQualifierREF01T4      ReferenceIdentificationQualifierREF01 = "T4"
	ReferenceIdentificationQualifierREF019B      ReferenceIdentificationQualifierREF01 = "9B"
	ReferenceIdentificationQualifierREF019D      ReferenceIdentificationQualifierREF01 = "9D"
	ReferenceIdentificationQualifierREF01F4      ReferenceIdentificationQualifierREF01 = "F4"
	ReferenceIdentificationQualifierREF01BT      ReferenceIdentificationQualifierREF01 = "BT"
	ReferenceIdentificationQualifierREF01VY      ReferenceIdentificationQualifierREF01 = "VY"
	ReferenceIdentificationQualifierREF01PXC     ReferenceIdentificationQualifierREF01 = "PXC"
)

var referenceIdentificationQualifierREF01List = codeList{
	{"EI", "Billing Provider Tax Identification", "Billing Provider Tax Identification — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"REF*EI*123456789\"); frontend/src/data/x12Names837p.ts BILLING_PROVIDER_REF_EL_NAME; also listed for the 2010AC, 2010BB, and 2330B identification boxes, Transmittal 579, PDF p.16, 21, 71; the guide's own name for the code is Employer's Identification Number"},
	{"6R", "Line Item Control Number", "Line Item Control Number — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"REF*6R*001H0000000140309201\"); Transmittal 579, PDF p.89; frontend/src/data/x12Names837p.ts SERVICE_REF_NAME.lineItemControl; the guide's own name for the code is Provider Control Number"},
	{"XZ", "Pharmacy Prescription Number", "Pharmacy Prescription Number — required with REF01 = XZ whenever the 2400 SV1 segment carries modifier J1; CEDI 837P CG, PDF p.47 (printed p.39; Table 41); frontend/src/data/x12Names837p.ts DRUG_ASSOCIATION_REF01_NAME"},
	{"SY", "Social Security Number", "Social Security Number — Transmittal 579, PDF p.12 (2010AA tax identification box, \"EI, SY\"), p.19 (2010BA), p.69 (2330A); the name is the guide's: " + tr3},
	{"0B", "State License Number", "State License Number — Transmittal 579, PDF p.12 (2010AA, \"0B, 1G\") and the 2310A-D, 2330C-F, and 2420A/D/E/F secondary identification boxes; the name is the guide's: " + tr3},
	{"1G", "Provider UPIN Number", "Provider UPIN Number — Transmittal 579, PDF p.12 (2010AA, \"0B, 1G\") and the 2310A/B/D, 2330C/D/F, and 2420A/D/E/F secondary identification boxes; the name is the guide's: " + tr3},
	{"2U", "Payer Identification Number", "Payer Identification Number — Transmittal 579, PDF p.15 (2010AC, \"2U, FY, NF\"), p.21 (2010BB), p.71 (2330B); the name is the guide's: " + tr3},
	{"FY", "Claim Office Number", "Claim Office Number — Transmittal 579, PDF p.15 (2010AC, \"2U, FY, NF\"), p.21 (2010BB), p.71 (2330B); the name is the guide's: " + tr3},
	{"NF", "National Association of Insurance Commissioners (NAIC) Code", "National Association of Insurance Commissioners (NAIC) Code — Transmittal 579, PDF p.15 (2010AC, \"2U, FY, NF\"), p.21 (2010BB), p.71 (2330B); the name is the guide's: " + tr3},
	{"Y4", "Agency Claim Number", "Agency Claim Number — Transmittal 579, PDF p.19 (2010BA Property and Casualty Claim Number box) and p.25 (2010CA); the name is the guide's: " + tr3},
	{"1W", "Member Identification Number", guide("Member Identification Number", "REF01 at the 2010CA Property and Casualty Patient Identifier box, beside SY")},
	{"G2", "Provider Commercial Number", "Provider Commercial Number — Transmittal 579, PDF p.22 (2010BB, \"G2, LU\") and the 2310A-D, 2330C-G, 2420A, and 2420C-F secondary identification boxes; the name is the guide's: " + tr3},
	{"LU", "Location Number", "Location Number — Transmittal 579, PDF p.22 (2010BB, \"G2, LU\") and the 2310B-D, 2330D-G, 2420A, 2420C, and 2420D secondary identification boxes; the name is the guide's: " + tr3},
	{"4N", "Service Authorization Exception Code", "Service Authorization Exception Code — the 2300 box's only REF01 value, Transmittal 579, PDF p.32; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.serviceAuthException; the guide's own name for the code is Special Payment Reference Number"},
	{"F5", "Mandatory Medicare (Section 4081) Crossover Indicator", "Mandatory Medicare (Section 4081) Crossover Indicator — the 2300 box's only REF01 value, Transmittal 579, PDF p.32; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.medicareCrossoverIndicator; the guide's own name for the code is Medicare Version Code"},
	{"EW", "Mammography Certification Number", "Mammography Certification Number — the only REF01 value of the 2300 and 2400 boxes, Transmittal 579, PDF p.33, 90; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.mammographyCertification"},
	{"9F", "Referral Number", "Referral Number — the only REF01 value of the 2300, 2330B, and 2400 referral boxes, Transmittal 579, PDF p.33, 72, 92; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.referral"},
	{"G1", "Prior Authorization", "Prior Authorization — the only REF01 value of the 2300, 2330B, and 2400 prior authorization boxes, Transmittal 579, PDF p.34, 71, 89; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.priorAuthorization; the guide's own name for the code is Prior Authorization Number"},
	{"F8", "Payer Claim Control Number", "Payer Claim Control Number — the only REF01 value of the 2300 box and the 2330B Other Payer Claim Control Number box, Transmittal 579, PDF p.34, 73; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.payerClaimControlNumber; the guide's own name for the code is Original Reference Number"},
	{"X4", "Clinical Laboratory Improvement Amendment (CLIA) Number", "Clinical Laboratory Improvement Amendment (CLIA) Number — the only REF01 value of the 2300 and 2400 CLIA boxes, Transmittal 579, PDF p.35, 90; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.clia"},
	{"9A", "Repriced Claim Number", "Repriced Claim Number — the 2300 box's only REF01 value, Transmittal 579, PDF p.35; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.repricedClaim; the guide's own name for the code is Repriced Claim Reference Number"},
	{"9C", "Adjusted Repriced Claim Number", "Adjusted Repriced Claim Number — the 2300 box's only REF01 value, Transmittal 579, PDF p.36; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.adjustedRepricedClaim; the guide's own name for the code is Adjusted Repriced Claim Reference Number"},
	{"LX", "Investigational Device Exemption Number", "Investigational Device Exemption Number — the 2300 box's only REF01 value, Transmittal 579, PDF p.36; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.investigationalDeviceExemption; the guide's own name for the code is Qualified Products List"},
	{"D9", "Claim Identifier for Transmission Intermediaries", "Claim Identifier for Transmission Intermediaries — the 2300 box's only REF01 value, Transmittal 579, PDF p.37; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.claimIdentifierForTransmissionIntermediaries; the guide's own name for the code is Claim Number"},
	{"EA", "Medical Record Number", "Medical Record Number — the 2300 box's only REF01 value, Transmittal 579, PDF p.37; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.medicalRecordNumber; the guide's own name for the code is Medical Record Identification Number"},
	{"P4", "Demonstration Project Identifier", "Demonstration Project Identifier — the 2300 box's only REF01 value, Transmittal 579, PDF p.38; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.demonstrationProject; the guide's own name for the code is Project Code"},
	{"1J", "Care Plan Oversight", "Care Plan Oversight — the 2300 box's only REF01 value, Transmittal 579, PDF p.38; frontend/src/data/x12Names837p.ts CLAIM_REF_NAME.carePlanOversight; the guide's own name for the code is Facility ID Number"},
	{"T4", "Other Payer Claim Adjustment Indicator", "Other Payer Claim Adjustment Indicator — the 2330B box's only REF01 value, Transmittal 579, PDF p.72; frontend/src/data/x12Names837p.ts OTHER_PAYER_REF_NAME; the guide's own name for the code is Signal Code"},
	{"9B", "Repriced Line Item Reference Number", "Repriced Line Item Reference Number — the 2400 box's only REF01 value, Transmittal 579, PDF p.88; frontend/src/data/x12Names837p.ts SERVICE_REF_NAME.repricedLineItem"},
	{"9D", "Adjusted Repriced Line Item Reference Number", "Adjusted Repriced Line Item Reference Number — the 2400 box's only REF01 value, Transmittal 579, PDF p.88; frontend/src/data/x12Names837p.ts SERVICE_REF_NAME.adjustedRepricedLineItem"},
	{"F4", "Referring Clinical Laboratory Improvement Amendment (CLIA) Facility Identification", "Referring Clinical Laboratory Improvement Amendment (CLIA) Facility Identification — the 2400 box's only REF01 value, Transmittal 579, PDF p.91; frontend/src/data/x12Names837p.ts SERVICE_REF_NAME.referringClia; the guide's own name for the code is Facility Certification Number"},
	{"BT", "Immunization Batch Number", "Immunization Batch Number — the 2400 box's only REF01 value, Transmittal 579, PDF p.91; frontend/src/data/x12Names837p.ts SERVICE_REF_NAME.immunizationBatch; the guide's own name for the code is Batch Number"},
	{"VY", "Link Sequence Number", "Link Sequence Number — Transmittal 579, PDF p.97 (2410 Prescription or Compound Drug Association Number box, \"VY, XZ\"); the name is the guide's: " + tr3},
	{"PXC", "Health Care Provider Taxonomy Code", "Health Care Provider Taxonomy Code — PRV02's only value at 2000A, 2310B, and 2420A, Transmittal 579, PDF p.9, 58, 98 (Values column); the name is the guide's: " + tr3},
}

func (v ReferenceIdentificationQualifierREF01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ReferenceIdentificationQualifierREF01) Known() bool {
	return referenceIdentificationQualifierREF01List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ReferenceIdentificationQualifierREF01) ShortDescription() string {
	return referenceIdentificationQualifierREF01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ReferenceIdentificationQualifierREF01) Description() string {
	return referenceIdentificationQualifierREF01List.description(string(v))
}

// ParseReferenceIdentificationQualifierREF01 returns text as a
// ReferenceIdentificationQualifierREF01, with an error when the text is
// neither empty nor on the list. The value is kept either way.
func ParseReferenceIdentificationQualifierREF01(text string) (ReferenceIdentificationQualifierREF01, error) {
	return ReferenceIdentificationQualifierREF01(text), referenceIdentificationQualifierREF01List.check("REF01", text)
}

// DateTimeQualifierDTP01 is DTP01, the Date/Time Qualifier, shared by every DTP
// segment in the transaction.
//
// The list is the union of the DTP01 codes the 837P implementation guide
// allows across its DTP boxes. Transmittal 907 names every DTP box's own
// title (CLAIM_DTP_NAME, SERVICE_DTP_NAME) but, for all but one of them, does
// not restate DTP01's own numeric qualifier (frontend/src/data/x12Names837p.ts's
// header note). The Admission Date box's 435 is confirmed by a CEDI Medicare
// presence requirement. The rest come from Transmittal 579, the original
// 005010 837P flat file, whose Values column prints DTP01's code at every DTP
// box; the parser places each DTP in its box's field by this code. A code
// that is the only value of its box takes the box title as its short name;
// the Disability Dates (314, 360, 361), Assumed and Relinquished Care (090,
// 091), and Test Date (738, 739) boxes list several codes without saying
// which is which, so those take the guide's own names.
type DateTimeQualifierDTP01 string

const (
	DateTimeQualifierDTP01NotSent DateTimeQualifierDTP01 = ""
	DateTimeQualifierDTP01435     DateTimeQualifierDTP01 = "435"
	DateTimeQualifierDTP01431     DateTimeQualifierDTP01 = "431"
	DateTimeQualifierDTP01454     DateTimeQualifierDTP01 = "454"
	DateTimeQualifierDTP01304     DateTimeQualifierDTP01 = "304"
	DateTimeQualifierDTP01453     DateTimeQualifierDTP01 = "453"
	DateTimeQualifierDTP01439     DateTimeQualifierDTP01 = "439"
	DateTimeQualifierDTP01484     DateTimeQualifierDTP01 = "484"
	DateTimeQualifierDTP01455     DateTimeQualifierDTP01 = "455"
	DateTimeQualifierDTP01471     DateTimeQualifierDTP01 = "471"
	DateTimeQualifierDTP01314     DateTimeQualifierDTP01 = "314"
	DateTimeQualifierDTP01360     DateTimeQualifierDTP01 = "360"
	DateTimeQualifierDTP01361     DateTimeQualifierDTP01 = "361"
	DateTimeQualifierDTP01297     DateTimeQualifierDTP01 = "297"
	DateTimeQualifierDTP01296     DateTimeQualifierDTP01 = "296"
	DateTimeQualifierDTP01096     DateTimeQualifierDTP01 = "096"
	DateTimeQualifierDTP01090     DateTimeQualifierDTP01 = "090"
	DateTimeQualifierDTP01091     DateTimeQualifierDTP01 = "091"
	DateTimeQualifierDTP01444     DateTimeQualifierDTP01 = "444"
	DateTimeQualifierDTP01050     DateTimeQualifierDTP01 = "050"
	DateTimeQualifierDTP01573     DateTimeQualifierDTP01 = "573"
	DateTimeQualifierDTP01472     DateTimeQualifierDTP01 = "472"
	DateTimeQualifierDTP01607     DateTimeQualifierDTP01 = "607"
	DateTimeQualifierDTP01463     DateTimeQualifierDTP01 = "463"
	DateTimeQualifierDTP01461     DateTimeQualifierDTP01 = "461"
	DateTimeQualifierDTP01738     DateTimeQualifierDTP01 = "738"
	DateTimeQualifierDTP01739     DateTimeQualifierDTP01 = "739"
	DateTimeQualifierDTP01011     DateTimeQualifierDTP01 = "011"
)

var dateTimeQualifierDTP01List = codeList{
	{"435", "Admission", "Admission — the Date - Admission box (2300 DTP); CEDI 837P CG, PDF p.37 (printed p.29, below Table 18): \"DTP01 = '435' must be present\" when 2400 SV105 is 21, 51, or 61; Transmittal 579, PDF p.30; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.admission"},
	{"431", "Onset of Current Illness/Symptom", "Onset of Current Illness/Symptom — the 2300 box's only DTP01 value, Transmittal 579, PDF p.27; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.onsetOfIllness"},
	{"454", "Initial Treatment", "Initial Treatment — the only DTP01 value of the 2300 and 2400 Initial Treatment boxes, Transmittal 579, PDF p.28, 86; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.initialTreatment"},
	{"304", "Last Seen Date", "Last Seen Date — the only DTP01 value of the 2300 Last Seen Date and 2400 Date Last Seen Date boxes, Transmittal 579, PDF p.28, 85; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.lastSeen; the guide's own name for the code is Latest Visit or Consultation"},
	{"453", "Acute Manifestation", "Acute Manifestation — the 2300 box's only DTP01 value, Transmittal 579, PDF p.28; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.acuteManifestation; the guide's own name for the code is Acute Manifestation of a Chronic Condition"},
	{"439", "Accident", "Accident — the 2300 box's only DTP01 value, Transmittal 579, PDF p.28; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.accident"},
	{"484", "Last Menstrual Period", "Last Menstrual Period — the 2300 box's only DTP01 value, Transmittal 579, PDF p.28; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.lastMenstrualPeriod"},
	{"455", "Last X-Ray Date", "Last X-Ray Date — the only DTP01 value of the 2300 and 2400 Last X-Ray Date boxes, Transmittal 579, PDF p.29, 86; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.lastXRay"},
	{"471", "Prescription Date", "Prescription Date — the only DTP01 value of the 2300 Hearing and Vision Prescription Date and 2400 Prescription Date boxes, Transmittal 579, PDF p.29, 84; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.hearingVisionPrescription"},
	{"314", "Disability", "Disability — one of three DTP01 values (314, 360, 361) of the 2300 Date - Disability Dates box, Transmittal 579, PDF p.29; the name is the guide's: " + tr3},
	{"360", "Initial Disability Period Start", "Initial Disability Period Start — one of three DTP01 values (314, 360, 361) of the 2300 Date - Disability Dates box, Transmittal 579, PDF p.29; the name is the guide's: " + tr3},
	{"361", "Initial Disability Period End", "Initial Disability Period End — one of three DTP01 values (314, 360, 361) of the 2300 Date - Disability Dates box, Transmittal 579, PDF p.29; the name is the guide's: " + tr3},
	{"297", "Last Worked", "Last Worked — the 2300 box's only DTP01 value, Transmittal 579, PDF p.29; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.lastWorked; the guide's own name for the code is Initial Disability Period Last Day Worked"},
	{"296", "Authorized Return to Work", "Authorized Return to Work — the 2300 box's only DTP01 value, Transmittal 579, PDF p.29; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.authorizedReturnToWork; the guide's own name for the code is Initial Disability Period Return To Work"},
	{"096", "Discharge", "Discharge — the 2300 box's only DTP01 value, Transmittal 579, PDF p.30; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.discharge"},
	{"090", "Report Start", "Report Start — one of two DTP01 values (090, 091) of the 2300 Date - Assumed and Relinquished Care Dates box, Transmittal 579, PDF p.30; the assumed care date; the name is the guide's: " + tr3},
	{"091", "Report End", "Report End — one of two DTP01 values (090, 091) of the 2300 Date - Assumed and Relinquished Care Dates box, Transmittal 579, PDF p.30; the relinquished care date; the name is the guide's: " + tr3},
	{"444", "Property and Casualty Date of First Contact", "Property and Casualty Date of First Contact — the 2300 box's only DTP01 value, Transmittal 579, PDF p.30; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.propertyCasualtyFirstContact; the guide's own name for the code is First Visit or Consultation"},
	{"050", "Repricer Received Date", "Repricer Received Date — the 2300 box's only DTP01 value, Transmittal 579, PDF p.30; frontend/src/data/x12Names837p.ts CLAIM_DTP_NAME.repricerReceived; the guide's own name for the code is Received"},
	{"573", "Check or Remittance Date", "Check or Remittance Date — the only DTP01 value of the 2330B Claim Check or Remittance Date and 2430 Line Check or Remittance Date boxes, Transmittal 579, PDF p.70, 110; the guide's own name for the code is Date Claim Paid"},
	{"472", "Service Date", "Service Date — the 2400 box's only DTP01 value, Transmittal 579, PDF p.84; frontend/src/data/x12Names837p.ts SERVICE_DTP_NAME.serviceDate"},
	{"607", "Certification Revision/Recertification Date", "Certification Revision/Recertification Date — the 2400 box's only DTP01 value, Transmittal 579, PDF p.84; frontend/src/data/x12Names837p.ts SERVICE_DTP_NAME.certificationRevision"},
	{"463", "Begin Therapy Date", "Begin Therapy Date — the 2400 box's only DTP01 value, Transmittal 579, PDF p.85; frontend/src/data/x12Names837p.ts SERVICE_DTP_NAME.beginTherapy"},
	{"461", "Last Certification Date", "Last Certification Date — the 2400 box's only DTP01 value, Transmittal 579, PDF p.85; frontend/src/data/x12Names837p.ts SERVICE_DTP_NAME.lastCertification"},
	{"738", "Most Recent Hemoglobin or Hematocrit or Both", "Most Recent Hemoglobin or Hematocrit or Both — one of two DTP01 values (738, 739) of the 2400 Date - Test Date box, Transmittal 579, PDF p.85; the name is the guide's: " + tr3},
	{"739", "Most Recent Serum Creatine", "Most Recent Serum Creatine — one of two DTP01 values (738, 739) of the 2400 Date - Test Date box, Transmittal 579, PDF p.85; the name is the guide's: " + tr3},
	{"011", "Shipped Date", "Shipped Date — the 2400 box's only DTP01 value, Transmittal 579, PDF p.85; frontend/src/data/x12Names837p.ts SERVICE_DTP_NAME.shipped"},
}

func (v DateTimeQualifierDTP01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v DateTimeQualifierDTP01) Known() bool { return dateTimeQualifierDTP01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v DateTimeQualifierDTP01) ShortDescription() string {
	return dateTimeQualifierDTP01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v DateTimeQualifierDTP01) Description() string {
	return dateTimeQualifierDTP01List.description(string(v))
}

// ParseDateTimeQualifierDTP01 returns text as a DateTimeQualifierDTP01, with
// an error when the text is neither empty nor on the list. The value is
// kept either way.
func ParseDateTimeQualifierDTP01(text string) (DateTimeQualifierDTP01, error) {
	return DateTimeQualifierDTP01(text), dateTimeQualifierDTP01List.check("DTP01", text)
}

// AmountQualifierAMT01 is AMT01, the Amount Qualifier Code, shared by every AMT
// segment in the transaction.
//
// The list is the six AMT01 codes the 837P implementation guide allows across
// its AMT boxes. Transmittal 907 names every AMT box's own title
// (CLAIM_AMT_NAME, the COB AMT boxes, SERVICE_AMT02_EL_NAME_BY_BOX) but does
// not restate AMT01's own code character for most of them. The COB Payer
// Amount Paid qualifier is confirmed by a CEDI Medicare presence requirement;
// the rest come from Transmittal 579's Values column, one code per AMT box,
// and each takes its box title as its short name. The parser places each AMT
// in its box's field by this code.
type AmountQualifierAMT01 string

const (
	AmountQualifierAMT01NotSent AmountQualifierAMT01 = ""
	AmountQualifierAMT01D       AmountQualifierAMT01 = "D"
	AmountQualifierAMT01F5      AmountQualifierAMT01 = "F5"
	AmountQualifierAMT01A8      AmountQualifierAMT01 = "A8"
	AmountQualifierAMT01EAF     AmountQualifierAMT01 = "EAF"
	AmountQualifierAMT01T       AmountQualifierAMT01 = "T"
	AmountQualifierAMT01F4      AmountQualifierAMT01 = "F4"
)

var amountQualifierAMT01List = codeList{
	{"D", "Payer Amount Paid", "Payer Amount Paid — Medicare requires one occurrence of the 2320 AMT with AMT01 = D whenever 2000B SBR01 = S; CEDI 837P CG, PDF p.40 (printed p.32; Table 29); Transmittal 579, PDF p.66 (2320 COB Payer Paid Amount box); frontend/src/data/x12Names837p.ts AMT01_COB_NAME; the guide's own name for the code is Payor Amount Paid"},
	{"F5", "Patient Amount Paid", "Patient Amount Paid — the 2300 box's only AMT01 value, Transmittal 579, PDF p.32; frontend/src/data/x12Names837p.ts CLAIM_AMT_NAME"},
	{"A8", "COB Total Non-Covered Amount", "COB Total Non-Covered Amount — the 2320 box's only AMT01 value, Transmittal 579, PDF p.66; the guide's own name for the code is Noncovered Charges - Actual"},
	{"EAF", "Remaining Patient Liability", "Remaining Patient Liability — the only AMT01 value of the 2320 and 2430 Remaining Patient Liability boxes, Transmittal 579, PDF p.66, 110; the guide's own name for the code is Amount Owed"},
	{"T", "Sales Tax Amount", "Sales Tax Amount — the 2400 box's only AMT01 value, Transmittal 579, PDF p.92; frontend/src/data/x12Names837p.ts SALES_TAX_AMT_NAME; the guide's own name for the code is Tax"},
	{"F4", "Postage Claimed Amount", "Postage Claimed Amount — the 2400 box's only AMT01 value, Transmittal 579, PDF p.92; frontend/src/data/x12Names837p.ts POSTAGE_CLAIMED_AMT_NAME; the guide's own name for the code is Postage Claimed"},
}

func (v AmountQualifierAMT01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v AmountQualifierAMT01) Known() bool { return amountQualifierAMT01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v AmountQualifierAMT01) ShortDescription() string {
	return amountQualifierAMT01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v AmountQualifierAMT01) Description() string {
	return amountQualifierAMT01List.description(string(v))
}

// ParseAmountQualifierAMT01 returns text as an AmountQualifierAMT01, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseAmountQualifierAMT01(text string) (AmountQualifierAMT01, error) {
	return AmountQualifierAMT01(text), amountQualifierAMT01List.check("AMT01", text)
}

// QuantityQualifierQTY01 is QTY01, the Quantity Qualifier, shared by every QTY
// segment in the transaction.
//
// The list is the two QTY01 codes the 837P implementation guide allows, one
// per QTY box. Neither Transmittal 907 nor the CEDI 837P Companion Guide
// names QTY01's individual codes in text; frontend/src/data/x12Names837p.ts
// QTY01_EL_NAME gives only the element's own position name;
// QTY02_EL_NAME_BY_BOX names the boxes by slug, not by QTY01's own character.
// Transmittal 579's Values column gives one code per QTY box (PT, FL), and
// each takes its box title as its short name; the parser places each QTY in
// its box's field by it.
type QuantityQualifierQTY01 string

const (
	QuantityQualifierQTY01NotSent QuantityQualifierQTY01 = ""
	QuantityQualifierQTY01PT      QuantityQualifierQTY01 = "PT"
	QuantityQualifierQTY01FL      QuantityQualifierQTY01 = "FL"
)

var quantityQualifierQTY01List = codeList{
	{"PT", "Ambulance Patient Count", "Ambulance Patient Count — the 2400 box's only QTY01 value, Transmittal 579, PDF p.86; frontend/src/data/x12Names837p.ts AMBULANCE_PATIENT_COUNT_QTY_NAME; the guide's own name for the code is Patients"},
	{"FL", "Obstetric Anesthesia Additional Units", "Obstetric Anesthesia Additional Units — the 2400 box's only QTY01 value, Transmittal 579, PDF p.86; frontend/src/data/x12Names837p.ts OBSTETRIC_ANESTHESIA_QTY_NAME; the guide's own name for the code is Units"},
}

func (v QuantityQualifierQTY01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v QuantityQualifierQTY01) Known() bool { return quantityQualifierQTY01List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v QuantityQualifierQTY01) ShortDescription() string {
	return quantityQualifierQTY01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v QuantityQualifierQTY01) Description() string {
	return quantityQualifierQTY01List.description(string(v))
}

// ParseQuantityQualifierQTY01 returns text as a QuantityQualifierQTY01, with
// an error when the text is neither empty nor on the list. The value is
// kept either way.
func ParseQuantityQualifierQTY01(text string) (QuantityQualifierQTY01, error) {
	return QuantityQualifierQTY01(text), quantityQualifierQTY01List.check("QTY01", text)
}
