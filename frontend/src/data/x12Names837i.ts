// Names and coded-value text for the Medicare 837I (institutional claim, HIPAA
// version 005010X223A2) — same role and citation conventions as
// frontend/src/data/x12Names835.ts (read that file first; this one copies its
// shape). Built only from public sources:
//
//   - "Transmittal 907, PDF p.N": docs/reference/file-formats/
//     cms-837-835-flatfiles-r907otn-2011.pdf (CR 7409, "AS OF May 25, 2011"),
//     a 283-page PDF. N is the PDF VIEWER page, counted from 1 — NOT the
//     sheet's own printed footer number, which
//     restarts at 1 on the institutional layout's first page; footer page F
//     is PDF page F+6 throughout the 837I section (PDF pp. 7-119). This is
//     also not the flat file's own "Loop ID"/"Seg. ID" columns, which are
//     X12 loop and segment identifiers, not page references.
//   - "Transmittal 579, PDF p.N": docs/reference/file-formats/
//     cms-837-flatfiles-r579otn-2009.pdf, the original 2009 CR (6575) for the
//     same 005010 837I flat file, used only where 907 is silent. 907 is
//     never silent on anything cited below from it alone, so 579 is
//     not cited by number here; it is the same kind of internal
//     flat-file layout as 907, not a source of NUBC code values either (see
//     the NUBC note below).
//   - "CMS Claims Processing Manual Ch. 25, §N (FL M)": Medicare Claims
//     Processing Manual, Pub. 100-04, Chapter 25 ("Completing and Processing
//     the Form CMS-1450 Data Set"), retrieved from cms.gov (not stored in
//     this repo) on 2026-09-30, Rev. 12423 (12-20-23). Cited by its own
//     section number and UB-04 Form Locator (FL), which are stable across
//     revisions; it has no useful PDF page numbers of its own for this
//     purpose. Only explains institutional concepts and provider-identifier
//     qualifiers (0B/1G/G2); it does not list the NUBC code sets (see below).
//   - X12's public pages (x12.org/examples/005010x223, x12.org/codes): used
//     only where cited inline; x12.org/codes does not publish a Reference
//     Identification Qualifier (element 128) list, so REF01 qualifier names
//     below come only from the flat file's own segment captions.
//
// No CEDI 837I companion guide exists in this repo (only cedi-837p-companion-
// guide.pdf, for the professional claim) and none is cited here.
//
// NUBC code values: Form CMS-1450 (UB-04) Type of Bill, Admission Type,
// Admission Source (Point of Origin), Patient (Discharge) Status, Condition
// Code, Occurrence Code, Occurrence Span Code, and Value Code are all
// maintained by the National Uniform Billing Committee (NUBC), not X12 or
// CMS. Transmittal 907 leaves every one of these elements' "Values" column
// blank, and CMS Claims Processing Manual Ch. 25 explicitly defers each of
// them to "Medicare contractors" and "the NUBC (www.nubc.org) via the NUBC's
// Official UB-04 Data Specifications Manual" (§§75.1-75.6) — a paid
// publication, not a public source available to this project. These
// elements below therefore keep only their X12 position and
// reference; no code list is given, and no name is invented.
//
// CPT codes are AMA-licensed and are never given descriptive text here (the
// 837I carries them only as HCPCS/CPT codes in SV202-2/SVD03-2, never as
// prose). Revenue codes (SV201, SVD04) and HCPCS Level II codes are NUBC/CMS
// codes, not AMA text, and are fine as bare codes; their full code lists are
// likewise NUBC-maintained and not reproduced here.
//
// Element-position arrays are 0-indexed from element 01 (index 0 = element
// 01), matching x12Names835.ts's ISA_EL_NAME convention. Segments whose shape
// repeats identically across every instance in this transaction (HL, N3, N4,
// generic NM1 positions 1-2, PER, REF01/02) are named once and reused, with a
// note; segments that Medicare's institutional claim does not use (N/U in
// the flat file) are included for completeness, with that noted.
//
// Go enum types this file's data implies (named <Concept><Segment><Element>,
// numbered sub-positions with underscore, e.g. CLM05_1):
//   TransactionTypeGS01_ST01, ClaimOrEncounterBHT06, HierarchicalLevelHL03,
//   PartyRoleNM101, EntityTypeNM102, IDQualifierNM108, ProviderRoleCodePRV01,
//   PayerResponsibilitySBR01, RelationshipSBR02 (shared with PAT01),
//   ClaimFilingIndicatorSBR09, SexDMG03, FacilityTypeCLM05_1 (unresolved — see
//   NUBC note), FacilityCodeQualifierCLM05_2, ClaimFrequencyCLM05_3
//   (unresolved), AssignmentCLM07, BenefitsAssignCLM08, ReleaseOfInfoCLM09,
//   DelayReasonCLM20, AdmissionTypeCL101 (unresolved), AdmissionSourceCL102
//   (unresolved), PatientStatusCL103 (unresolved), PWKReportTypeCode01,
//   PWKTransmissionCode02, DiagnosisQualifierHI01_1, PresentOnAdmissionHI0n_9,
//   CRCCategoryCRC01, CRCCertificationCRC02, CRCConditionIndicator,
//   HCPPricingMethodologyHCP01, UnitOrBasisOfMeasurement, HCPRejectReasonHCP13,
//   HCPPolicyComplianceHCP14, HCPExceptionCodeHCP15, ProductServiceQualifier
//   (SV202_1/SVD03_1), NDCQualifierLIN02, CASGroupCode01, OtherInsuranceIndicator
//   (OI03/OI06), NTENoteCode01, DateQualifierDTP01 (per-loop, see DTP tables).

// ---------------------------------------------------------------------------
// Bare segment names (loop-independent).
// ---------------------------------------------------------------------------
export const SEGMENT_NAME: Partial<Record<string, string>> = {
  ISA: 'Interchange Control Header', // Transmittal 907, PDF p.7
  GS: 'Functional Group Header', // Transmittal 907, PDF p.7
  ST: 'Transaction Set Header', // Transmittal 907, PDF p.8
  BHT: 'Beginning of Hierarchical Transaction', // Transmittal 907, PDF p.8
  HL: 'Hierarchical Level', // Transmittal 907, PDF p.10 (first instance, 2000A); same shape at every HL
  SBR: 'Subscriber Information', // Transmittal 907, PDF p.16
  PAT: 'Patient Information', // Transmittal 907, PDF p.20
  DMG: 'Demographic Information', // Transmittal 907, PDF p.17 (subscriber); p.22 (patient)
  CLM: 'Claim Information', // Transmittal 907, PDF p.23
  DTP: 'Date/Time Period', // Transmittal 907, PDF p.24 (claim-level; also used at several other loops, see DTP tables)
  CL1: 'Institutional Claim Code', // Transmittal 907, PDF p.24
  PWK: 'Claim Supplemental Information', // Transmittal 907, PDF p.25 (claim-level); 'Line Supplemental Information' at 2400, PDF p.112
  CN1: 'Contract Information', // Transmittal 907, PDF p.25
  AMT: 'Amount', // name varies by loop/qualifier, see AMT tables
  REF: 'Reference Identification', // name varies by loop/qualifier, see REF tables
  K3: 'File Information', // Transmittal 907, PDF p.29
  NTE: 'Note', // name varies by loop ('Claim Note' / 'Billing Note' / 'Third Party Organization Notes'), see NTE table
  CRC: 'EPSDT Referral', // Transmittal 907, PDF p.29
  HI: 'Health Care Information Codes', // name varies by diagnosis/procedure/code type, see HI_TYPE_NAME
  HCP: 'Claim Pricing/Repricing Information', // Transmittal 907, PDF p.91 (claim-level); 'Line Pricing/Repricing Information' at 2400, PDF p.114
  PRV: 'Provider Information', // Transmittal 907, PDF p.10 ('Billing Provider Specialty Information'); 'Attending Provider Specialty Information' at 2310A, p.92
  NM1: 'Individual or Organizational Name', // role varies by loop, see NM1_ROLE_NAME
  N3: 'Address', // role varies by loop (e.g. 'Billing Provider Address'); same positional shape everywhere
  N4: 'City, State, ZIP Code', // role varies by loop; same positional shape everywhere
  PER: 'Administrative Communications Contact', // same shape at every instance, see PER table
  MIA: 'Inpatient Adjudication Information', // Transmittal 907, PDF p.100 (loop 2320, coordination of benefits)
  MOA: 'Outpatient Adjudication Information', // Transmittal 907, PDF p.102 (loop 2320, coordination of benefits)
  OI: 'Other Health Insurance Information', // Transmittal 907, PDF p.100
  LX: 'Service Line Number', // Transmittal 907, PDF p.111
  SV2: 'Institutional Service Line', // Transmittal 907, PDF p.111
  LIN: 'Drug Identification', // Transmittal 907, PDF p.115
  CTP: 'Drug Quantity', // Transmittal 907, PDF p.116
  SVD: 'Line Adjudication Information', // Transmittal 907, PDF p.121
  CAS: 'Claim/Line Adjustment', // name varies by loop ('Claim Level Adjustments' at 2320; 'Line Adjustment' at 2430)
  SE: 'Transaction Set Trailer', // Transmittal 907, PDF p.124
  GE: 'Functional Group Trailer', // Transmittal 907, PDF p.124
  IEA: 'Interchange Control Trailer', // Transmittal 907, PDF p.125
}

// ---------------------------------------------------------------------------
// Envelope: ISA / GS / ST / BHT / SE / GE / IEA.
// ---------------------------------------------------------------------------
export const ISA_EL_NAME: string[] = [
  'Authorization Information Qualifier', 'Authorization Information',
  'Security Information Qualifier', 'Security Information',
  'Interchange ID Qualifier', 'Interchange Sender ID',
  'Interchange ID Qualifier', 'Interchange Receiver ID',
  'Interchange Date', 'Interchange Time', 'Repetition Separator',
  'Interchange Control Version Number', 'Interchange Control Number',
  'Acknowledgement Requested', 'Interchange Usage Indicator', 'Component Element Separator',
] // Transmittal 907, PDF p.7

export const GS_EL_NAME: string[] = [
  'Functional Identifier Code', "Application Sender's Code", "Application Receiver's Code",
  'Date', 'Time', 'Group Control Number', 'Responsible Agency Code', 'Version/Release/Industry Identifier Code',
] // Transmittal 907, PDF p.7-8

// GS01/ST01: the transaction-type code, same combined name either way.
export const TRANSACTION_TYPE_NAME: Record<string, string> = {
  HC: 'Health Care Claim', // Transmittal 907, PDF p.7 (GS01)
  '837': 'Health Care Claim', // Transmittal 907, PDF p.8 (ST01)
}

export const ST_EL_NAME: string[] = ['Transaction Set Identifier Code', 'Transaction Set Control Number', 'Implementation Guide Version Name'] // Transmittal 907, PDF p.8

export const BHT_EL_NAME: string[] = [
  'Hierarchical Structure Code', 'Transaction Set Purpose Code', 'Originator Application Transaction ID',
  'Transaction Set Creation Date', 'Transaction Set Creation Time', 'Claim or Encounter ID',
] // Transmittal 907, PDF p.8-9

export const BHT01_NAME: Record<string, string> = {
  '0019': 'Information Source, Information Receiver, Provider, Subscriber, Dependent', // Transmittal 907, PDF p.8 (hierarchical structure code value)
}
export const BHT02_NAME: Record<string, string> = {
  '00': 'Original', // Transmittal 907, PDF p.8
  '18': 'Reissue', // Transmittal 907, PDF p.8
}
// BHT06: Transmittal 907, PDF p.9 lists values "31, CH, RP". CH (Chargeable)
// and RP (Reporting) are X12's documented claim/encounter codes; the flat
// file gives no further description for either code or for 31, so only the
// code list is reproduced, not an invented meaning for 31.
export const BHT06_VALUES: string[] = ['31', 'CH', 'RP'] // Transmittal 907, PDF p.9

export const SE_EL_NAME: string[] = ['Transaction Segment Count', 'Transaction Set Control Number'] // Transmittal 907, PDF p.124
export const GE_EL_NAME: string[] = ['Number of Transaction Sets Included', 'Group Control Number'] // Transmittal 907, PDF p.124
export const IEA_EL_NAME: string[] = ['Number of Included Functional Groups', 'Interchange Control Number'] // Transmittal 907, PDF p.125

// ---------------------------------------------------------------------------
// NM1 role (NM101), shared across every loop that carries an NM1. The same
// code means the same role everywhere it appears (1000A/B, 2010AA/AB/AC,
// 2010BA/BB, 2010CA, 2310A-F, 2330A-I, 2420A-D).
// ---------------------------------------------------------------------------
export const NM1_ROLE_NAME: Record<string, string> = {
  '41': 'Submitter', // Transmittal 907, PDF p.9 (1000A)
  '40': 'Receiver', // Transmittal 907, PDF p.10 (1000B)
  '85': 'Billing Provider', // Transmittal 907, PDF p.12 (2010AA); also 2330I (other payer billing provider), PDF p.111
  '87': 'Pay-To Address', // Transmittal 907, PDF p.13 (2010AB)
  PE: 'Pay-To Plan', // Transmittal 907, PDF p.14 (2010AC)
  IL: 'Insured/Subscriber', // Transmittal 907, PDF p.16 (2010BA); also 2330A (other subscriber), PDF p.102
  PR: 'Payer', // Transmittal 907, PDF p.18 (2010BB); also 2330B (other payer), PDF p.104
  QC: 'Patient', // Transmittal 907, PDF p.20 (2010CA)
  '71': 'Attending Provider', // Transmittal 907, PDF p.92 (2310A); also 2330C, PDF p.106
  '72': 'Operating Physician', // Transmittal 907, PDF p.93 (2310B); also 2330D, PDF p.107, and 2420A, PDF p.118
  ZZ: 'Other Operating Physician', // Transmittal 907, PDF p.94 (2310C); also 2330E, PDF p.108, and 2420B, PDF p.118-119
  '82': 'Rendering Provider', // Transmittal 907, PDF p.94 (2310D); also 2330G, PDF p.109, and 2420C, PDF p.120
  '77': 'Service Facility Location', // Transmittal 907, PDF p.95 (2310E); also 2330F, PDF p.108-109
  DN: 'Referring Provider', // Transmittal 907, PDF p.96 (2310F); also 2330H, PDF p.110, and 2420D, PDF p.120-121
}

export const NM101_EL_NAME = 'Entity Identifier Code' // Transmittal 907, PDF p.9, consistent across every NM1 instance
export const NM102_EL_NAME = 'Entity Type Qualifier'
export const NM102_NAME: Record<string, string> = {
  '1': 'Person', // Transmittal 907, PDF p.9
  '2': 'Non-Person Entity', // Transmittal 907, PDF p.9
}
export const NM108_NAME: Record<string, string> = {
  XX: 'National Provider Identifier', // Transmittal 907, PDF p.12 (NM109 described as an NPI wherever XX is used)
  '46': 'Electronic Transmitter Identification Number', // Transmittal 907, PDF p.9 (submitter/receiver NM108)
  PI: 'Payer Identification', // Transmittal 907, PDF p.9 (pay-to plan/payer NM108)
  XV: 'Health Care Financing Administration National PlanID', // Transmittal 907, PDF p.9 (pay-to plan/payer NM108, alternate with PI)
  II: 'Standard Unique Health Identifier for each Individual in the United States', // Transmittal 907, PDF p.17 (subscriber NM108, alternate with MI)
  MI: 'Member Identification Number', // Transmittal 907, PDF p.17 (subscriber NM108) — carries the Medicare Beneficiary Identifier for Medicare claims
}

// Submitter (1000A).
export const NM1_SUBMITTER_EL_NAME: string[] = ['Submitter Last or Organization Name', 'Submitter First Name', 'Submitter Middle Name or Initial', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Submitter Identifier'] // Transmittal 907, PDF p.9
// Receiver (1000B).
export const NM1_RECEIVER_EL_NAME: string[] = ['Receiver Name', 'Name First (N/U)', 'Name Middle (N/U)', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Receiver Primary Identifier'] // Transmittal 907, PDF p.10

// PER: identical shape at every instance (submitter 1000A, billing provider
// 2010AA). PER01 is always 'Contact Function Code'; PER02 onward repeats
// name/communication-number triplets.
export const PER01_NAME: Record<string, string> = {
  IC: 'Information Contact', // Transmittal 907, PDF p.9
}
export const PER_EL_NAME: string[] = [
  'Contact Name', 'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number', 'Communication Number Qualifier', 'Communication Number',
] // Transmittal 907, PDF p.9-10 (1000A); identical shape at 2010AA, PDF p.13
export const PER_QUALIFIER_VALUES: string[] = ['EM', 'EX', 'FX', 'TE'] // Transmittal 907, PDF p.9-10 (PER03/05/07; PER03's own first-use set is narrower, EM/FX/TE)

// ---------------------------------------------------------------------------
// Billing provider: 2000A (HL, PRV, CUR), 2010AA (NM1, N3, N4, REF, PER),
// 2010AB (pay-to address, NM1/N3/N4), 2010AC (pay-to plan, NM1/N3/N4/REF).
// ---------------------------------------------------------------------------
export const HL_EL_NAME: string[] = ['Hierarchical ID Number', 'Hierarchical Parent ID Number', 'Hierarchical Level Code', 'Hierarchical Child Code'] // Transmittal 907, PDF p.10 (2000A); identical shape at every HL (2000B PDF p.16, 2000C PDF p.20)
export const HL03_NAME: Record<string, string> = {
  '20': 'Billing Provider', // Transmittal 907, PDF p.10
  '22': 'Subscriber', // Transmittal 907, PDF p.16
  '23': 'Patient', // Transmittal 907, PDF p.20 (not used for Medicare, where the patient is always the subscriber)
}

export const PRV_EL_NAME: string[] = ['Provider Code', 'Reference Identification Qualifier', 'Provider Taxonomy Code'] // Transmittal 907, PDF p.10 (2000A billing); identical shape at 2310A attending, PDF p.92
export const PRV01_NAME: Record<string, string> = {
  BI: 'Billing', // Transmittal 907, PDF p.10 (2000A Billing Provider Specialty Information)
  AT: 'Attending', // Transmittal 907, PDF p.92 (2310A Attending Provider Specialty Information)
}
export const PRV02_NAME: Record<string, string> = {
  PXC: 'Health Care Provider Taxonomy Code', // Transmittal 907, PDF p.10
}

export const CUR_EL_NAME: string[] = ['Entity Identifier Code', 'Currency Code'] // Transmittal 907, PDF p.11 (2000A Foreign Currency Information; remaining CUR03-21 are Not Used)
export const CUR01_NAME: Record<string, string> = {
  '85': 'Billing Provider', // Transmittal 907, PDF p.11
}

// Address (N3) / City-State-ZIP (N4): identical positional shape at every
// loop that carries them. Only the box caption's role word changes (e.g.
// "Billing Provider Address", "Payer City, State, ZIP Code").
export const N3_EL_NAME: string[] = ['Address Line', 'Address Line'] // Transmittal 907, PDF p.12 (2010AA), representative of every N3 instance
export const N4_EL_NAME: string[] = ['City Name', 'State or Province Code', 'Postal Zone or ZIP Code', 'Country Code', 'Location Qualifier (N/U)', 'Location Identifier (N/U)', 'Country Subdivision Code'] // Transmittal 907, PDF p.12 (2010AA), representative of every N4 instance

export const NM1_BILLING_EL_NAME: string[] = ['Billing Provider Last or Organizational Name', 'Billing Provider First Name (N/U)', 'Billing Provider Middle Name (N/U)', 'Name Prefix (N/U)', 'Billing Provider Name Suffix (N/U)', 'Identification Code Qualifier', 'Billing Provider Identifier'] // Transmittal 907, PDF p.12 (2010AA)
export const REF_BILLING_TAX_NAME = 'Billing Provider Tax Identification' // Transmittal 907, PDF p.13 (2010AA REF, qualifier EI)

export const NM1_PAYTO_EL_NAME: string[] = ['Pay-to Provider Last or Organization Name (N/U)', 'Pay-to Provider First Name (N/U)', 'Pay-to Provider Middle Name (N/U)', 'Name Prefix (N/U)', 'Pay-to Provider Name Suffix (N/U)', 'Identification Code Qualifier (N/U)', 'Pay-to Provider Identifier (N/U)'] // Transmittal 907, PDF p.13-14 (2010AB) — name/ID positions are all Not Used in this flat file; only entity code + address flow
export const NM1_PAYTOPLAN_EL_NAME: string[] = ['Pay-to Plan Organizational Name', 'Name First (N/U)', 'Name Middle (N/U)', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Identification Code'] // Transmittal 907, PDF p.14-15 (2010AC)
export const REF_PAYTOPLAN_NAME: Record<string, string> = {
  '2U': 'Pay-to Plan Secondary Identification', // Transmittal 907, PDF p.15 (box caption; REF01 shared with FY, NF)
  FY: 'Pay-to Plan Secondary Identification',
  NF: 'Pay-to Plan Secondary Identification',
  EI: 'Pay-to Plan Tax Identification', // Transmittal 907, PDF p.15-16
}

// ---------------------------------------------------------------------------
// Subscriber (2000B / 2010BA) and patient (2000C / 2010CA).
// ---------------------------------------------------------------------------
export const SBR_EL_NAME: string[] = [
  'Payer Responsibility Sequence Number Code', 'Individual Relationship Code', 'Subscriber Group or Policy Number', 'Subscriber Group Name',
  'Insurance Type Code (N/U)', 'Coordination of Benefits Code (N/U)', 'Yes/No Condition or Response Code (N/U)', 'Employment Status Code (N/U)',
  'Claim Filing Indicator Code',
] // Transmittal 907, PDF p.16 (2000B); identical shape at 2320 Other Subscriber Information, PDF p.97

export const SBR01_NAME: Record<string, string> = {
  A: 'Payer Responsibility Four', // Transmittal 907, PDF p.16 (value set A-H, P, S, T, U; only P/S/T are named by the X12 standard as Primary/Secondary/Tertiary — the flat file lists the full code set without individual names for A-H/U)
  P: 'Primary', // common X12 Payer Responsibility Sequence Number code
  S: 'Secondary',
  T: 'Tertiary',
}
export const SBR01_VALUES: string[] = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'P', 'S', 'T', 'U'] // Transmittal 907, PDF p.16

// Individual Relationship Code — shared by SBR02 (subscriber/other-subscriber
// coverage) and PAT01 (patient loop).
export const RELATIONSHIP_NAME: Record<string, string> = {
  '01': 'Spouse', // common X12 Individual Relationship Code
  '18': 'Self', // Transmittal 907, PDF p.16 (SBR02, 2000B — the only value listed there; Medicare's patient is always the subscriber)
  '19': 'Child',
  '20': 'Employee',
  '21': 'Unknown',
  '39': 'Organ Donor',
  '40': 'Cadaver Donor',
  '53': 'Life Partner',
  G8: 'Other Relationship',
}
export const SBR02_VALUES_2000B: string[] = ['18'] // Transmittal 907, PDF p.16
export const PAT01_VALUES: string[] = ['01', '19', '20', '21', '39', '40', '53', 'G8'] // Transmittal 907, PDF p.20 (2000C PAT01)
export const SBR02_VALUES_2320: string[] = ['01', '18', '19', '20', '21', '39', '40', '53', 'G8'] // Transmittal 907, PDF p.97 (2320 SBR02, other subscriber)

export const SBR09_NAME: Record<string, string> = {
  '11': 'Other Non-Federal Programs', // common X12 Claim Filing Indicator Code
  '12': 'Preferred Provider Organization (PPO)',
  '13': 'Point of Service (POS)',
  '14': 'Exclusive Provider Organization (EPO)',
  '15': 'Indemnity Insurance',
  '16': 'Health Maintenance Organization (HMO) Medicare Risk',
  '17': 'Dental Maintenance Organization',
  AM: 'Automobile Medical',
  BL: 'Blue Cross/Blue Shield',
  CH: 'Champus',
  CI: 'Commercial Insurance Co.',
  DS: 'Disability',
  FI: 'Federal Employees Program',
  HM: 'Health Maintenance Organization',
  LM: 'Liability Medical',
  MA: 'Medicare Part A',
  MB: 'Medicare Part B',
  MC: 'Medicaid',
  OF: 'Other Federal Program',
  TV: 'Title V',
  VA: 'Veterans Affairs Plan',
  WC: "Workers' Compensation Health Claim",
  ZZ: 'Mutually Defined',
}
export const SBR09_VALUES: string[] = ['11', '12', '13', '14', '15', '16', '17', 'AM', 'BL', 'CH', 'CI', 'DS', 'FI', 'HM', 'LM', 'MA', 'MB', 'MC', 'OF', 'TV', 'VA', 'WC', 'ZZ'] // Transmittal 907, PDF p.16 (2000B SBR09); identical set at 2320, PDF p.97

export const NM1_SUBSCRIBER_EL_NAME: string[] = ['Subscriber Last Name', 'Subscriber First Name', 'Subscriber Middle Name or Initial', 'Name Prefix (N/U)', 'Subscriber Name Suffix', 'Identification Code Qualifier', 'Subscriber Primary Identifier'] // Transmittal 907, PDF p.16-17 (2010BA)

export const DMG_EL_NAME: string[] = ['Date Time Period Format Qualifier', 'Birth Date', 'Gender Code'] // Transmittal 907, PDF p.17 (subscriber, 2010BA); identical shape at patient 2010CA, PDF p.22
export const DMG01_NAME: Record<string, string> = { D8: 'Date Expressed in Format CCYYMMDD' } // Transmittal 907, PDF p.17
export const DMG03_NAME: Record<string, string> = {
  F: 'Female', // Transmittal 907, PDF p.17
  M: 'Male',
  U: 'Unknown',
}

export const REF_SUBSCRIBER_NAME: Record<string, string> = {
  SY: 'Subscriber Secondary Identification', // Transmittal 907, PDF p.18 (2010BA; REF02 labeled Social Security Number)
  Y4: 'Property and Casualty Claim Number', // Transmittal 907, PDF p.18 (2010BA); also 2010CA patient, PDF p.22-23
}

export const NM1_PAYER_EL_NAME: string[] = ['Payer Name', 'Name First (N/U)', 'Subscriber Name Middle or Initial (N/U)', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Payer Identifier'] // Transmittal 907, PDF p.18-19 (2010BB)
// REF (2010BB): Transmittal 907 shows two REF boxes at this loop whose box
// caption and REF02 element label appear to have been swapped with each
// other in the source spreadsheet (one captioned "Payer Secondary
// Identification" with REF02 labeled "Billing Provider Secondary
// Identifier", values 2U/EI/FY/NF; the other captioned "Billing Provider
// Secondary Identification" with REF02 labeled "Payer Additional
// Identifier", values G2/LU). Both REF segments sit in loop 2010BB (payer);
// transcribed here by loop and qualifier rather than by the possibly-swapped
// box captions.
export const REF_PAYER_NAME: Record<string, string> = {
  '2U': 'Payer Secondary Identification', // Transmittal 907, PDF p.19
  EI: 'Payer Secondary Identification',
  FY: 'Payer Secondary Identification',
  NF: 'Payer Secondary Identification',
  G2: 'Payer Secondary Identification', // Transmittal 907, PDF p.19-20
  LU: 'Payer Secondary Identification',
}

export const PAT_EL_NAME: string[] = [
  'Individual Relationship Code', 'Patient Location Code (N/U)', 'Employment Status Code (N/U)', 'Student Status Code (N/U)',
  'Date Time Period Format Qualifier (N/U)', 'Patient Death Date (N/U)', 'Unit or Basis for Measurement Code (N/U)', 'Patient Weight (N/U)', 'Pregnancy Indicator (N/U)',
] // Transmittal 907, PDF p.20 (2000C) — not used for Medicare, where the patient level is never present

export const NM1_PATIENT_EL_NAME: string[] = ['Patient Last Name', 'Patient First Name', 'Patient Middle Name or Initial', 'Name Prefix (N/U)', 'Patient Name Suffix', 'Identification Code Qualifier (N/U)', 'Patient Primary Identifier (N/U)'] // Transmittal 907, PDF p.20-21 (2010CA)
export const REF_PATIENT_NAME: Record<string, string> = {
  '1W': 'Property and Casualty Patient Identifier', // Transmittal 907, PDF p.22-23 (2010CA)
  SY: 'Property and Casualty Patient Identifier',
  Y4: 'Property and Casualty Claim Number', // Transmittal 907, PDF p.22-23 (2010CA)
}

// ---------------------------------------------------------------------------
// Claim (2300): CLM, DTP, CL1, PWK, CN1, AMT, REF, K3, NTE, CRC, HI, HCP.
// ---------------------------------------------------------------------------
export const CLM_EL_NAME: string[] = [
  'Patient Control Number', 'Total Claim Charge Amount', // CLM01-02: Transmittal 907, PDF p.23
  'Claim Filing Indicator Code (N/U)', 'Non-Institutional Claim Type Code (N/U)',
  'Health Care Service Location Information', 'Place of Service Code', 'Facility Code Qualifier', 'Claim Frequency Code', // CLM05 composite: p.23
  'Provider or Supplier Signature Indicator (N/U)', 'Assignment or Plan Participation Code', 'Benefits Assignment Certification Indicator', 'Release of Information Code', // CLM06-09: p.23
  'Patient Signature Source Code (N/U)', 'Related Causes Information (N/U)', 'Special Program Indicator (N/U)',
  'Yes/No Condition or Response Code (N/U)', 'Level of Service Code (N/U)', 'Yes/No Condition or Response Code (N/U)', 'Participation Agreement (N/U)',
  'Claim Status Code (N/U)', 'Yes/No Condition or Response Code (N/U)', 'Claim Submission Reason Code (N/U)', 'Delay Reason Code', // CLM17-20: p.24
] // Transmittal 907, PDF p.23-24

// CLM05 is a composite: CLM05-1 Place of Service Code, CLM05-2 Facility Code
// Qualifier, CLM05-3 Claim Frequency Code.
export const CLM05_2_NAME: Record<string, string> = {
  A: 'Uniform Billing Claim Form Bill Type', // Transmittal 907, PDF p.23 (CLM05-2's only listed value)
}
// CLM05-1 (facility type + bill classification digits of the UB-04 Type of
// Bill) and CLM05-3 (Claim Frequency Code, the Type of Bill's 4th digit) are
// NUBC code values; Transmittal 907's "Values" column is blank for both (PDF
// p.23). See the NUBC note at the top of this file.
export const CLM07_NAME: Record<string, string> = {
  A: 'Assigned', // Transmittal 907, PDF p.23
  B: 'Assignment Accepted on Clinical Lab Services Only',
  C: 'Not Assigned',
}
export const CLM08_NAME: Record<string, string> = {
  N: 'No', // Transmittal 907, PDF p.23
  W: 'Not Applicable',
  Y: 'Yes',
}
export const CLM09_NAME: Record<string, string> = {
  I: 'Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes', // Transmittal 907, PDF p.23
  Y: 'Yes, Provider has a Signed Statement Permitting Release of Medical Billing Data Related to a Claim',
}
export const CLM20_VALUES: string[] = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '15'] // Transmittal 907, PDF p.24 (Delay Reason Code — code list only, no per-code text given in the flat file)

// DTP (claim-level, 2300).
export const DTP01_EL_NAME = 'Date/Time Qualifier' // Transmittal 907, PDF p.24, consistent across every DTP instance
export const DTP_CLAIM_NAME: Record<string, string> = {
  '096': 'Discharge Hour', // Transmittal 907, PDF p.24
  '434': 'Statement From-Through Dates', // Transmittal 907, PDF p.24
  '435': 'Admission Date and Hour', // Transmittal 907, PDF p.24
  '050': 'Repricer Received Date', // Transmittal 907, PDF p.24
}

export const CL1_EL_NAME: string[] = ['Admission Type Code', 'Admission Source Code', 'Patient Status Code', 'Nursing Home Residential Status Code (N/U)'] // Transmittal 907, PDF p.24-25
// CL101 (Admission Type/Priority), CL102 (Admission Source/Point of Origin)
// and CL103 (Patient Status) are NUBC-maintained code values. Transmittal
// 907's "Values" column is blank for all three (PDF p.24-25). CMS Claims
// Processing Manual Ch. 25 §75.1 (FL 14, FL 15) and §75.2 (FL 17) confirm
// these codes are "available from Medicare contractors" and the NUBC's
// (non-public) Official UB-04 Data Specifications Manual, not published by
// CMS itself. No code list is given here; see the NUBC note at the top of
// this file.

export const PWK01_VALUES_CLAIM: string[] = [
  '03', '04', '05', '06', '07', '08', '09', '10', '11', '13', '15', '21', 'A3', 'AS', 'AM', 'B2', 'B3', 'B4', 'BR', 'BS', 'BT',
  'CB', 'CK', 'CT', 'D2', 'DA', 'DB', 'DG', 'DJ', 'DS', 'EB', 'HC', 'HR', 'I5', 'IR', 'LA', 'M1', 'MT', 'NN', 'OB', 'OC', 'OD',
  'OE', 'OX', 'OZ', 'P4', 'P5', 'PE', 'PN', 'PO', 'PQ', 'PY', 'PZ', 'RB', 'RR', 'RT', 'RX', 'SG', 'V5', 'XP',
] // Transmittal 907, PDF p.25 (claim-level PWK01, Report Type Code, X12 element 755 — the flat file gives only the code list, not per-code meanings)
export const PWK02_VALUES_CLAIM: string[] = ['AA', 'BM', 'EL', 'EM', 'FT', 'FX'] // Transmittal 907, PDF p.25 (claim-level PWK02, Report Transmission Code)
export const PWK_EL_NAME: string[] = ['Report Type Code', 'Report Transmission Code', 'Report Copies Needed (N/U)', 'Entity Identifier Code (N/U)', 'Identification Code Qualifier', 'Attachment Control Number'] // Transmittal 907, PDF p.25 (claim-level); identical shape at 2400 line-level, PDF p.112-113

export const CN1_EL_NAME: string[] = ['Contract Type Code', 'Contract Amount', 'Contract Percentage', 'Contract Code', 'Terms Discount Percent', 'Contract Version Identifier'] // Transmittal 907, PDF p.25-26
export const CN101_VALUES: string[] = ['01', '02', '03', '04', '05', '06', '09'] // Transmittal 907, PDF p.25 (Contract Type Code — code list only; no per-code text given in the flat file, and this is a generic X12 code, not an NUBC one)

export const AMT01_NAME: Record<string, string> = {
  F3: 'Patient Estimated Amount Due', // Transmittal 907, PDF p.26 (2300 claim-level AMT)
  D: 'Other Payer Paid Amount, COB Payer Paid Amount', // Transmittal 907, PDF p.99 (2320 COB Payer Paid Amount)
  EAF: 'Remaining Patient Liability', // Transmittal 907, PDF p.99 (2320); also 2430, PDF p.124
  A8: 'COB Total Non-Covered Amount', // Transmittal 907, PDF p.100 (2320)
  GT: 'Service Tax Amount', // Transmittal 907, PDF p.114 (2400 line-level)
  N8: 'Facility Tax Amount', // Transmittal 907, PDF p.114 (2400 line-level)
}

export const REF_CLAIM_NAME: Record<string, string> = {
  '4N': 'Service Authorization Exception Code', // Transmittal 907, PDF p.26
  '9F': 'Referral Number', // Transmittal 907, PDF p.26
  G1: 'Prior Authorization', // Transmittal 907, PDF p.26
  F8: 'Payer Claim Control Number', // Transmittal 907, PDF p.27 (used on a replacement/void claim to name the original claim's ICN)
  '9A': 'Repriced Claim Number', // Transmittal 907, PDF p.27
  '9C': 'Adjusted Repriced Claim Number', // Transmittal 907, PDF p.27
  LX: 'Investigational Device Exemption Number', // Transmittal 907, PDF p.27
  D9: 'Claim Identifier for Transmission Intermediaries', // Transmittal 907, PDF p.27-28 (REF02 labeled 'Value Added Network Trace Number')
  LU: 'Auto Accident State or Province', // Transmittal 907, PDF p.28
  EA: 'Medical Record Number', // Transmittal 907, PDF p.28
  P4: 'Demonstration Project Identifier', // Transmittal 907, PDF p.28
  G4: 'Peer Review Organization (PRO) Approval Number', // Transmittal 907, PDF p.28
}
export const REF_SERVICE_AUTH_VALUES: string[] = ['1', '2', '3', '4', '5', '6', '7'] // Transmittal 907, PDF p.26 (REF02 values for the 4N Service Authorization Exception Code)
// An additional claim-level REF (qualifier CN) appears in the source between
// the PRO Approval Number REF and K3 with no segment box caption given
// (Transmittal 907, PDF p.29); its business name is not stated there and is
// not reproduced here.

export const K3_EL_NAME: string[] = ['Fixed Format Information', 'Record Format Code (N/U)', 'Composite Unit of Measure (N/U)'] // Transmittal 907, PDF p.29

export const NTE_EL_NAME: string[] = ['Note Reference Code', 'Note Text'] // Transmittal 907, PDF p.29, same shape for every NTE instance
export const NTE_CLAIM_NAME = 'Claim Note' // Transmittal 907, PDF p.29
export const NTE_BILLING_NAME = 'Billing Note' // Transmittal 907, PDF p.29
export const NTE_LINE_NAME = 'Third Party Organization Notes' // Transmittal 907, PDF p.114 (2400 line-level)
export const NTE01_VALUES_CLAIM: string[] = ['ALG', 'DCP', 'DGN', 'DME', 'MED', 'NTR', 'ODT', 'RHB', 'RNH', 'SET', 'SFM', 'SPT', 'UPI'] // Transmittal 907, PDF p.29 (Claim Note)
export const NTE01_VALUES_BILLING: string[] = ['ADD'] // Transmittal 907, PDF p.29 (Billing Note)
export const NTE01_VALUES_LINE: string[] = ['TPO'] // Transmittal 907, PDF p.114 (Third Party Organization Notes)

export const CRC_EL_NAME: string[] = ['Code Category', 'Certification Condition Code Applies Indicator', 'Condition Indicator', 'Condition Indicator', 'Condition Indicator', 'Condition Indicator (N/U)', 'Condition Indicator (N/U)'] // Transmittal 907, PDF p.29
export const CRC01_NAME: Record<string, string> = { ZZ: 'Mutually Defined' } // Transmittal 907, PDF p.29 (EPSDT Referral code category)
export const CRC02_NAME: Record<string, string> = { N: 'No', Y: 'Yes' } // Transmittal 907, PDF p.29
export const CRC_CONDITION_VALUES: string[] = ['AV', 'NU', 'S2', 'ST'] // Transmittal 907, PDF p.29 (CRC03-05 Condition Indicator — code list only, no per-code text given in the flat file)

// HI: the "Healthcare Code Information" composite (C022), reused for every
// HI01 through HI12 position in every HI segment instance across the file.
// HI0n-6 (Quantity), HI0n-7 (Version Identifier) and HI0n-8 (Industry Code)
// are Not Used throughout. HI0n-9 (Present on Admission) is used only on the
// diagnosis-bearing HI segments (Principal/Other/Admitting Diagnosis,
// External Cause of Injury).
export const HI_EL_NAME: string[] = ['Diagnosis/Procedure/Code Type Qualifier', 'Code', 'Date Format Qualifier', 'Date', 'Monetary Amount'] // Transmittal 907, PDF p.24 (first HI instance), positions 1-5; representative of every HI0n composite
export const HI_POA_EL_NAME = 'Present on Admission Indicator' // Transmittal 907, PDF p.30 (HI0n-9, diagnosis composites only)
export const HI_POA_NAME: Record<string, string> = {
  Y: 'Yes', // Transmittal 907, PDF p.30
  N: 'No',
  U: 'Unknown',
  W: 'Clinically Undetermined',
}

// HI01-1 (and the equivalent position in every repeated composite): which
// kind of code the composite carries. The ICD-10 ("AB*"-prefixed) and
// ICD-9 (unprefixed) qualifiers are paired for the same concept.
export const HI_TYPE_NAME: Record<string, string> = {
  ABK: 'Principal Diagnosis (ICD-10-CM)', // Transmittal 907, PDF p.30
  BK: 'Principal Diagnosis (ICD-9-CM)', // Transmittal 907, PDF p.30
  ABJ: 'Admitting Diagnosis (ICD-10-CM)', // Transmittal 907, PDF p.34
  BJ: 'Admitting Diagnosis (ICD-9-CM)', // Transmittal 907, PDF p.34
  APR: "Patient's Reason for Visit (ICD-10-CM)", // Transmittal 907, PDF p.39
  PR: "Patient's Reason for Visit (ICD-9-CM)", // Transmittal 907, PDF p.39
  ABN: 'External Cause of Injury (ICD-10-CM)', // Transmittal 907, PDF p.44
  BN: 'External Cause of Injury (ICD-9-CM)', // Transmittal 907, PDF p.44
  DR: 'Diagnosis Related Group (DRG) Code', // Transmittal 907, PDF p.48 (Diagnosis Related Group Information HI segment)
  ABF: 'Other Diagnosis (ICD-10-CM)', // Transmittal 907, PDF p.51
  BF: 'Other Diagnosis (ICD-9-CM)', // Transmittal 907, PDF p.51
  BBR: 'Principal Procedure (ICD-10-PCS)', // Transmittal 907, PDF p.58
  BR: 'Principal Procedure (ICD-9-CM)', // Transmittal 907, PDF p.58
  BBQ: 'Other Procedure (ICD-10-PCS)', // Transmittal 907, PDF p.62
  BQ: 'Other Procedure (ICD-9-CM)', // Transmittal 907, PDF p.62
  BI: 'Occurrence Span Code', // Transmittal 907, PDF p.67
  BH: 'Occurrence Code', // Transmittal 907, PDF p.72
  BE: 'Value Code', // Transmittal 907, PDF p.77
  BG: 'Condition Code', // Transmittal 907, PDF p.81
  TC: 'Treatment Code', // Transmittal 907, PDF p.85
}
// Condition Code (BG), Occurrence Code (BH), Occurrence Span Code (BI) and
// Value Code (BE) are themselves NUBC-maintained code values (HI0n-2's
// content). Transmittal 907's "Values" column there is blank; see the NUBC
// note at the top of this file. Revenue codes and HCPCS Level II codes are
// NUBC/CMS codes too, used as bare codes elsewhere in this file (SV201,
// SVD04) and not enumerated for the same reason.

export const HCP_EL_NAME: string[] = [
  'Pricing Methodology', 'Repriced Allowed Amount', 'Repriced Saving Amount', 'Repricing Organization Identifier',
  'Repriced Per Diem or Flat Rate Amount', 'Repriced Approved DRG Code', 'Repriced Approved Amount', 'Product/Service ID (N/U)',
  'Product or Service ID Qualifier', 'Repriced Approved HCPCS Code', 'Unit or Basis for Measurement Code', 'Quantity',
  'Reject Reason Code', 'Policy Compliance Code', 'Exception Code',
] // Transmittal 907, PDF p.91-92 (claim-level); identical shape at 2400 line-level, PDF p.114-115
export const HCP01_NAME: Record<string, string> = {
  '00': 'Zero Pricing (Not Covered Under Contract)', // common X12 Pricing Methodology code; full 00-14 list given without per-code text below
}
export const HCP01_VALUES: string[] = ['00', '01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12', '13', '14'] // Transmittal 907, PDF p.91 (code list only, no per-code text given in the flat file beyond 00)
export const HCP09_VALUES: string[] = ['ER', 'HC', 'HP', 'IV', 'WK'] // Transmittal 907, PDF p.91 (HCP09 Product or Service ID Qualifier; same set as SV202-1/SVD03-1, see below)
export const HCP11_VALUES: string[] = ['DA', 'UN'] // Transmittal 907, PDF p.91 (Days/Units)
export const HCP13_VALUES: string[] = ['T1', 'T2', 'T3', 'T4', 'T5', 'T6'] // Transmittal 907, PDF p.92 (Reject Reason Code — code list only)
export const HCP14_VALUES: string[] = ['1', '2', '3', '4', '5'] // Transmittal 907, PDF p.92 (Policy Compliance Code — code list only)
export const HCP15_VALUES: string[] = ['1', '2', '3', '4', '5', '6'] // Transmittal 907, PDF p.92 (Exception Code — code list only)

// ---------------------------------------------------------------------------
// Claim-level providers (2310A-F): attending, operating, other operating,
// rendering, service facility, referring. Each is NM1 (+ PRV at 2310A only)
// + REF (up to 4 repeats, same qualifier set at every one of these loops).
// ---------------------------------------------------------------------------
export const NM1_ATTENDING_EL_NAME: string[] = ['Attending Provider Last Name', 'Attending Provider First Name', 'Attending Provider Middle Name', 'Name Prefix (N/U)', 'Attending Provider Name Suffix', 'Identification Code Qualifier', 'Attending Provider Identifier'] // Transmittal 907, PDF p.92
export const NM1_OPERATING_EL_NAME: string[] = ['Operating Physician Last or Organization Name', 'Operating Physician First Name', 'Rendering Provider Middle Name or Initial', 'Name Prefix (N/U)', 'Operating Physician Name Suffix', 'Identification Code Qualifier', 'Identifier'] // Transmittal 907, PDF p.93 (2310B); the flat file's own column labels NM105 "Rendering Provider Middle Name or Initial" in this operating-physician loop
export const NM1_OTHER_OPERATING_EL_NAME: string[] = ['Last or Organization Name', 'First Name', 'Middle Name', 'Name Prefix (N/U)', 'Name Suffix', 'Identification Code Qualifier', 'Identifier'] // Transmittal 907, PDF p.94 (2310C)
export const NM1_RENDERING_EL_NAME: string[] = ['Rendering Provider Last or Organization Name', 'Rendering Provider First Name', 'Rendering Provider Middle Name', 'Name Prefix (N/U)', 'Rendering Provider Name Suffix', 'Identification Code Qualifier', 'Rendering Provider Identifier'] // Transmittal 907, PDF p.94 (2310D)
export const NM1_SERVICE_FACILITY_EL_NAME: string[] = ['Laboratory or Facility Name', 'Name First (N/U)', 'Name Middle (N/U)', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Laboratory or Facility Primary Identifier'] // Transmittal 907, PDF p.95 (2310E)
export const NM1_REFERRING_EL_NAME: string[] = ['Referring Provider Last Name', 'Referring Provider First Name', 'Referring Provider Middle Name', 'Name Prefix (N/U)', 'Referring Provider Name Suffix', 'Identification Code Qualifier', 'Referring Provider Identifier'] // Transmittal 907, PDF p.96-97 (2310F)

// REF secondary identification at every one of 2310A-F (and, at identical
// shape, 2420A-D line-level providers): REF01 qualifier set is 0B/1G/G2/LU.
// CMS Claims Processing Manual Ch. 25 §75.6 (FL 76, FL 77, FL 78/79) names
// three of the four:
export const PROVIDER_REF_NAME: Record<string, string> = {
  '0B': 'State License Number', // CMS Claims Processing Manual Ch. 25, §75.6 (FL 76/77/78-79)
  '1G': 'Provider UPIN Number', // CMS Claims Processing Manual Ch. 25, §75.6 (FL 76/77/78-79)
  G2: 'Provider Commercial Number', // CMS Claims Processing Manual Ch. 25, §75.6 (FL 76/77/78-79)
  // LU: listed in Transmittal 907's REF01 value set (PDF p.92-96) but not
  // named in Chapter 25 or the flat file itself; left unresolved.
}
export const REF_PROVIDER_EL_NAME: string[] = ['Reference Identification Qualifier', 'Provider Secondary Identifier'] // Transmittal 907, PDF p.92 (2310A), representative of every provider-loop REF

// ---------------------------------------------------------------------------
// Other subscribers/payers (2320 / 2330A-I): coordination of benefits.
// ---------------------------------------------------------------------------
export const MIA_EL_NAME: string[] = [
  'Covered Days or Visits Count', 'PPS Operating Outlier Amount', 'Lifetime Psychiatric Days Count', 'Claim DRG Amount',
  'Claim Payment Remark Code', 'Claim Disproportionate Share Amount', 'Claim MSP Pass-Through Amount', 'Claim PPS Capital Amount',
  'PPS-Capital FSP DRG Amount', 'PPS-Capital HSP DRG Amount', 'PPS-Capital DSH DRG Amount', 'Old Capital Amount', 'PPS-Capital IME Amount',
  'PPS-Operating Hospital Specific DRG Amount', 'Cost Report Day Count', 'PPS-Operating Federal Specific DRG Amount', 'Claim PPS Capital Outlier Amount',
  'Claim Indirect Teaching Amount', 'Nonpayable Professional Component Amount',
  'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code',
  'PPS-Capital Exception Amount',
] // Transmittal 907, PDF p.100-101 (2320 Inpatient Adjudication Information — same segment, element shape and relative position as the 835's MIA)
export const MOA_EL_NAME: string[] = [
  'Reimbursement Rate', 'HCPCS Payable Amount', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code',
  'Claim Payment Remark Code', 'Claim Payment Remark Code', 'End Stage Renal Disease Payment Amount', 'Nonpayable Professional Component Amount',
] // Transmittal 907, PDF p.102 (2320 Outpatient Adjudication Information)
// Claim Payment Remark Code values are the CMS/CORE Remittance Advice Remark
// Code (RARC) set, already embedded in frontend/src/data/x12Names835.ts's
// RARC table (from backend/pkg/rarc.json); not duplicated here.

export const OI_EL_NAME: string[] = ['Claim Filing Indicator Code (N/U)', 'Claim Submission Reason Code (N/U)', 'Benefits Assignment Certification Indicator', 'Patient Signature Source Code (N/U)', 'Provider Agreement Code (N/U)', 'Release of Information Code'] // Transmittal 907, PDF p.100
export const OI03_NAME: Record<string, string> = { N: 'No', W: 'Not Applicable', Y: 'Yes' } // Transmittal 907, PDF p.100
export const OI06_NAME: Record<string, string> = { I: 'Informed Consent to Release Medical Information for Conditions or Diagnoses Regulated by Federal Statutes', Y: 'Yes' } // Transmittal 907, PDF p.100

export const SBR_OTHER_EL_NAME: string[] = SBR_EL_NAME // Transmittal 907, PDF p.97 (2320 Other Subscriber Information) — identical shape to 2000B's SBR

export const NM1_OTHER_SUBSCRIBER_EL_NAME: string[] = ['Other Insured Last Name', 'Other Insured First Name', 'Other Insured Middle Name or Initial', 'Name Prefix (N/U)', 'Other Insured Name Suffix', 'Identification Code Qualifier', 'Other Insured Identifier'] // Transmittal 907, PDF p.102 (2330A)
export const REF_OTHER_SUBSCRIBER_NAME: Record<string, string> = { SY: 'Other Subscriber Secondary Identification' } // Transmittal 907, PDF p.103 (2330A)

export const NM1_OTHER_PAYER_EL_NAME: string[] = ['Other Payer Organization Name', 'Name First (N/U)', 'Name Middle (N/U)', 'Name Prefix (N/U)', 'Name Suffix (N/U)', 'Identification Code Qualifier', 'Other Payer Primary Identifier'] // Transmittal 907, PDF p.104 (2330B)
export const DTP_OTHER_PAYER_NAME: Record<string, string> = { '573': 'Date Claim Check or Remittance Date' } // Transmittal 907, PDF p.105 (2330B); identical qualifier at 2430 line-level, PDF p.124
export const REF_OTHER_PAYER_NAME: Record<string, string> = {
  '2U': 'Other Payer Secondary Identification', // Transmittal 907, PDF p.105
  EI: 'Other Payer Secondary Identification',
  FY: 'Other Payer Secondary Identification',
  NF: 'Other Payer Secondary Identification',
  G1: 'Other Payer Prior Authorization Number', // Transmittal 907, PDF p.105
  '9F': 'Other Payer Referral Number', // Transmittal 907, PDF p.105
  T4: 'Other Payer Claim Adjustment Indicator', // Transmittal 907, PDF p.106
  F8: 'Other Payer Claim Control Number', // Transmittal 907, PDF p.106
}

// 2330C-I (other payer attending/operating/other operating/service
// facility/rendering/referring/billing provider) reuse the same NM1 shape as
// their 2310 counterparts above and the same REF qualifier set as
// PROVIDER_REF_NAME (0B/1G/G2/LU), per Transmittal 907, PDF p.106-111.
export const OTHER_PAYER_PROVIDER_LOOP_NAME: Record<string, string> = {
  '2330C': 'Other Payer Attending Provider', // Transmittal 907, PDF p.106, NM101 71
  '2330D': 'Other Payer Operating Physician', // Transmittal 907, PDF p.107, NM101 72
  '2330E': 'Other Payer Other Operating Physician', // Transmittal 907, PDF p.108, NM101 ZZ
  '2330F': 'Other Payer Service Facility Location', // Transmittal 907, PDF p.108-109, NM101 77
  '2330G': 'Other Payer Rendering Provider', // Transmittal 907, PDF p.109, NM101 82
  '2330H': 'Other Payer Referring Provider', // Transmittal 907, PDF p.110, NM101 DN
  '2330I': 'Other Payer Billing Provider', // Transmittal 907, PDF p.111, NM101 85
}

// ---------------------------------------------------------------------------
// Service line (2400): LX, SV2, PWK, DTP, REF, AMT, HCP; 2410 LIN/CTP/REF;
// 2420A-D line providers; 2430 SVD/CAS/DTP/AMT.
// ---------------------------------------------------------------------------
export const LX01_EL_NAME = 'Assigned Number' // Transmittal 907, PDF p.111

export const SV2_EL_NAME: string[] = [
  'Revenue Code', 'Composite Medical Procedure Identifier', 'Product or Service ID Qualifier', 'Procedure Code', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Description', 'Product/Service ID (N/U)', // SV201-202(composite): Transmittal 907, PDF p.111-112
  'Line Item Charge Amount', 'Unit or Basis for Measurement Code', 'Service Units/Days Count', 'Unit Rate (N/U)', 'Line Item Denied or Non-Covered Charge Amount', // SV203-207: p.112
  'Yes/No Condition or Response Code (N/U)', 'Nursing Home Residential Status Code (N/U)', 'Level of Care Code (N/U)', // SV208-210: p.112
]
export const SV202_1_VALUES: string[] = ['ER', 'HC', 'HP', 'IV', 'WK'] // Transmittal 907, PDF p.112 (SV202-1 Product or Service ID Qualifier; same set reused at HCP09, SVD03-1)

export const PWK01_VALUES_LINE: string[] = [
  '11', '13', '15', '21', 'A3', 'A4', 'AM', 'AS', 'B2', 'B3', 'B4', 'BR', 'BS', 'BT', 'CB', 'CT', 'D2', 'DA', 'DB', 'DG', 'DJ', 'DS',
  'EB', 'HC', 'HR', 'I5', 'IR', 'LA', 'M1', 'MT', 'NN', 'OB', 'OC', 'OD', 'OE', 'OX', 'OZ', 'P4', 'P5', 'PE', 'PN', 'PO', 'PQ', 'PY',
  'PZ', 'RB', 'RR', 'RT', 'RX', 'SG', 'V5', 'XP',
] // Transmittal 907, PDF p.112 (line-level PWK01 — narrower set than the claim-level list)
export const PWK02_VALUES_LINE: string[] = ['SG', 'V5', 'XP'] // Transmittal 907, PDF p.112 (line-level PWK02 — narrower set than the claim-level list)

export const DTP472_NAME = 'Service Date' // Transmittal 907, PDF p.113 (2400 DTP, qualifier 472)

export const REF_LINE_NAME: Record<string, string> = {
  '6R': 'Line Item Control Number', // Transmittal 907, PDF p.113
  '9B': 'Repriced Line Item Reference Number', // Transmittal 907, PDF p.113
  '9D': 'Adjusted Repriced Line Item Reference Number', // Transmittal 907, PDF p.113-114
}

export const LIN_EL_NAME: string[] = ['Assigned Identification (N/U)', 'Product or Service ID Qualifier', 'National Drug Code'] // Transmittal 907, PDF p.115; LIN04-31 (repeated Product/Service ID Qualifier/ID pairs) are all Not Used
export const LIN02_NAME: Record<string, string> = { N4: 'National Drug Code in 5-4-2 Format' } // Transmittal 907, PDF p.115

export const CTP_EL_NAME: string[] = ['Class of Trade Code (N/U)', 'Price Identifier Code (N/U)', 'Unit Price (N/U)', 'National Drug Unit Count', 'Composite Unit of Measure'] // Transmittal 907, PDF p.115-116; CTP06-11 are Not Used
export const CTP05_1_VALUES: string[] = ['F2', 'GR', 'ME', 'ML', 'UN'] // Transmittal 907, PDF p.116 (Composite Unit of Measure code qualifier)

export const REF_PRESCRIPTION_NAME: Record<string, string> = { VY: 'Link Sequence Number', XZ: 'Pharmacy Prescription Number' } // Transmittal 907, PDF p.116-117 (2410 Prescription or Compound Drug Association Number — REF01 values VY, XZ; individual meanings per common X12 usage, not stated in the flat file itself)

// 2420A-D line providers reuse the same NM1 and REF shapes as their 2310
// counterparts (operating, other operating, rendering, referring); see
// NM1_ROLE_NAME and PROVIDER_REF_NAME above.
export const LINE_PROVIDER_LOOP_NAME: Record<string, string> = {
  '2420A': 'Operating Physician', // Transmittal 907, PDF p.118, NM101 72
  '2420B': 'Other Operating Physician', // Transmittal 907, PDF p.118-119, NM101 ZZ
  '2420C': 'Rendering Provider', // Transmittal 907, PDF p.120, NM101 82
  '2420D': 'Referring Provider', // Transmittal 907, PDF p.120-121, NM101 DN
}

export const SVD_EL_NAME: string[] = [
  'Other Payer Primary Identifier', 'Service Line Paid Amount', 'Composite Medical Procedure Identifier', 'Product or Service ID Qualifier',
  'Procedure Code', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Code Description', 'Product/Service ID (N/U)',
  'Revenue Code', 'Paid Service Unit Count', 'Bundled or Unbundled Line Number',
] // Transmittal 907, PDF p.121-122
export const SVD03_1_VALUES: string[] = ['ER', 'HC', 'HP', 'IV', 'WK'] // Transmittal 907, PDF p.122 (SVD03-1; same set as SV202-1/HCP09)

export const CAS_GROUP_NAME: Record<string, string> = {
  CO: 'Contractual Obligation', // x12.org/codes/claim-adjustment-group-codes
  CR: 'Correction and Reversal', // x12.org/codes/claim-adjustment-group-codes
  OA: 'Other Adjustment', // x12.org/codes/claim-adjustment-group-codes
  PI: 'Payer Initiated Reductions', // x12.org/codes/claim-adjustment-group-codes
  PR: 'Patient Responsibility', // x12.org/codes/claim-adjustment-group-codes
}
export const CAS_CLAIM_NAME = 'Claim Level Adjustments' // Transmittal 907, PDF p.98 (2320)
export const CAS_LINE_NAME = 'Line Adjustment' // Transmittal 907, PDF p.122 (2430)
export const CAS01_EL_NAME = 'Claim Adjustment Group Code' // Transmittal 907, PDF p.98 (2320), p.122 (2430)
export const CAS_REASON_EL_NAME = 'Adjustment Reason Code' // Transmittal 907, PDF p.98, p.122 — Claim Adjustment Reason Codes (CARC), same list as embedded in x12Names835.ts's CARC table
export const CAS_AMOUNT_EL_NAME = 'Adjustment Amount'
export const CAS_QTY_EL_NAME = 'Adjustment Quantity'

export const DTP573_LINE_NAME = 'Line Check or Remittance Date' // Transmittal 907, PDF p.124 (2430 DTP, qualifier 573)
export const AMT_LINE_LIABILITY_NAME = 'Remaining Patient Liability' // Transmittal 907, PDF p.124 (2430 AMT, qualifier EAF)
