// Names and coded-value text for the Medicare 837P (Health Care Claim: Professional,
// 005010X222A1) claim, built only from public sources, in the same shape and
// citation style as x12Names835.ts (read that file first for the convention this
// one follows).
//
// Citation form throughout this file:
//   "Transmittal 907, PDF p.N" - N is the PDF viewer page of the 283-page
//   docs/reference/file-formats/cms-837-835-flatfiles-r907otn-2011.pdf, counted
//   from 1. Its 837P 5010A1 flat-file attachment runs PDF p.133-241 (CLM begins
//   PDF p.154) and prints its own footer page numbers 1-109 starting at PDF
//   p.133 (PDF page = footer page + 132). This is NOT the attachment's own
//   "TR3 page" or "Page #" column, which indexes the paid X12 implementation
//   guide, a different document.
//
//   "CEDI 837P CG, PDF p.N (printed p.M)" - N is the PDF viewer page of the
//   56-page docs/reference/file-formats/cedi-837p-companion-guide.pdf; M is
//   that page's own printed footer number (N = M + 8 for every content page,
//   M = 1 starting at PDF p.9; the first 8 PDF pages are a cover and front
//   matter the guide numbers with roman numerals). The guide's own tables cite
//   a "Page #" column that indexes the TR3, not this PDF or the flat file.
//
//   "x12.org/codes/..." - a current, externally maintained X12 code list page,
//   used only for the handful of code sets X12 publishes outside the paid TR3
//   (e.g. claim adjustment group codes - the same list the 835 side of this
//   project already cites).
//
// Scope notes:
//   - The 837P 5010A1 flat-file attachment (Transmittal 907, PDF p.133-241) is
//     the primary source for every segment box
//     title and every element's positional name below. The CEDI 837P
//     Companion Guide (56 pages) is the
//     source for every Medicare-restricted coded value and every segment/loop
//     Medicare does not accept (its own Table 5). The historical 2009
//     docs/reference/file-formats/cms-837-flatfiles-r579otn-2009.pdf fills a
//     gap Transmittal 907 leaves open: Transmittal 907 names each qualifier
//     element's box but, unlike some rows of the 835 flat file, does not
//     restate the 2-3 character code value itself (HL03, NM101, REF01,
//     DTP01, AMT01, QTY01, CRC01, HI0n-1, ...) in its own Values column for
//     most boxes. Transmittal 579 is the predecessor 837P flat file for the
//     same 005010 version and, for exactly the qualifier positions where 907
//     is silent, states that Values column; a handful of the Go enum types
//     below cite it directly (e.g. NM101_NAME, HL03_NAME, REF01_NAME,
//     DTP01_NAME, AMT01_NAME, QTY01_NAME, CRC01_NAME, HI0x_1_NAME). Its
//     citation form is "Transmittal 579, PDF p.N" - N is the PDF viewer page
//     of the 251-page
//     docs/reference/file-formats/cms-837-flatfiles-r579otn-2009.pdf, counted
//     from 1 (confirmed against its own printed footer page numbers, which
//     run a constant 5 pages behind the PDF viewer page). Where neither 907
//     nor 579 states a qualifier's code value, this file still names the
//     segment by its box title only and does not invent the code
//     character(s); the comment says so at each such entry.
//   - Transmittal 907's 837P attachment names every loop's box (e.g. "BILLING
//     PROVIDER NAME") but, unlike some rows of the 835 flat file, does not
//     restate the qualifier's own 2-3 character code value (NM101, REF01,
//     DTP01, AMT01, CAS01...) anywhere in the rows this file draws on. Where a
//     qualifier's code value is given below, it is because that literal code
//     appears in the CEDI guide's own worked 837P transaction (section 11.2,
//     Figures 3-6) or in one of its Medicare-restriction tables (section 10) -
//     each such value is cited individually at the point it is used. Where a
//     qualifier's code value is not confirmed in either source (most REF01,
//     DTP01, NM101 and AMT01 instances), this file names the segment by its
//     box title only and does not invent the code character(s); the comment
//     says so at each such entry.
//   - CPT procedure code descriptors are AMA-licensed and are never written
//     here; SV101-2/SV501-2/SVD03-2 (Procedure Code) keep only their
//     positional name, never a value table.
//   - Every value below is a synthetic/structural code name, never PHI.
//   - Names written elsewhere in this repository are not consulted and not
//     reused here.
//
// Go enum types this file's values are meant to back (a type is
// listed here even where this file gives no value table for it - per the note
// above, an element with no confirmed public value table keeps its own
// reference, e.g. "CLM05-1", and the type is still named so the Go package
// knows the seam exists):
//   EntityIdentifierNM101, EntityTypeQualifierNM102,
//   IdentificationCodeQualifierNM108, ReferenceIdentificationQualifierREF01,
//   DateTimeQualifierDTP01, AmountQualifierAMT01, QuantityQualifierQTY01,
//   PayerResponsibilitySBR01, RelationshipSBR02, RelationshipPAT01,
//   ClaimFilingIndicatorSBR09, ClaimFrequencyCLM05_3, PlaceOfServiceCLM05_1,
//   FacilityCodeQualifierCLM05_2, ProviderOrSupplierSignatureCLM06,
//   AssignmentOrPlanParticipationCLM07, BenefitsAssignmentCertificationCLM08,
//   ReleaseOfInformationCLM09, PatientSignatureSourceCLM10,
//   RelatedCausesCodeCLM11_1, DelayReasonCLM20, ProviderTypePRV01,
//   UnitOrBasisForMeasurementCode, ClaimAdjustmentGroupCodeCAS01,
//   AttachmentReportTypeCodePWK01, AttachmentTransmissionCodePWK02,
//   DiagnosisTypeCodeHI0x_1, CodeCategoryCRC01, ConditionIndicatorCRC0x,
//   YesNoConditionResponseCode, ContactFunctionCodePER01,
//   CommunicationNumberQualifierPER03, PricingMethodologyHCP01.

// ===========================================================================
// Envelope: ISA / GS / ST / BHT / SE / GE / IEA
// ===========================================================================

export const SEGMENT_NAME: Partial<Record<string, string>> = {
  ISA: 'Interchange Control Header', // Transmittal 907, PDF p.133
  GS: 'Functional Group Header', // Transmittal 907, PDF p.133-134
  ST: 'Transaction Set Header', // Transmittal 907, PDF p.134
  BHT: 'Beginning of Hierarchical Transaction', // Transmittal 907, PDF p.134-135
  SE: 'Transaction Set Trailer', // Transmittal 907, PDF p.238 (footer 106)
  GE: 'Functional Group Trailer', // Transmittal 907, PDF p.238 (footer 106)
  IEA: 'Interchange Control Trailer', // Transmittal 907, PDF p.239 (footer 107)
}

export const ISA_EL_NAME: string[] = [
  'Authorization Information Qualifier', 'Authorization Information',
  'Security Information Qualifier', 'Security Information',
  'Interchange ID Qualifier', 'Interchange Sender ID',
  'Interchange ID Qualifier', 'Interchange Receiver ID',
  'Interchange Date', 'Interchange Time', 'Repetition Separator',
  'Interchange Control Version Number', 'Interchange Control Number',
  'Acknowledgement Requested', 'Usage Indicator', 'Component Element Separator',
] // Transmittal 907, PDF p.133 (footer 1)

export const GS_EL_NAME: string[] = [
  'Functional Identifier Code', 'Application Sender Code', 'Application Receiver Code', // PDF p.133 (footer 1)
  'Date', 'Time', 'Group Control Number', 'Responsible Agency Code', 'Version/Release/Industry Identifier Code', // PDF p.134 (footer 2)
] // Transmittal 907

export const ST_EL_NAME: string[] = ['Transaction Set Identifier Code', 'Transaction Set Control Number', 'Implementation Guide Version Name'] // Transmittal 907, PDF p.134 (footer 2)

export const BHT_EL_NAME: string[] = [
  'Hierarchical Structure Code', 'Transaction Set Purpose Code', 'Originator Application Transaction ID',
  'Transaction Set Creation Date', 'Transaction Set Creation Time', 'Claim or Encounter ID',
] // Transmittal 907, PDF p.134-135 (footer 2-3)
// CEDI restricts Medicare submissions to BHT02 = 00 (Original) and BHT06 = CH
// (Chargeable). CEDI 837P CG, PDF p.30 (printed p.22; Table 7)
export const BHT02_NAME: Record<string, string> = { '00': 'Original' } // CEDI 837P CG, PDF p.30 (printed p.22; Table 7)
export const BHT06_NAME: Record<string, string> = { CH: 'Chargeable' } // CEDI 837P CG, PDF p.30 (printed p.22; Table 7)

export const SE_EL_NAME: string[] = ['Transaction Segment Count', 'Transaction Set Control Number'] // Transmittal 907, PDF p.238 (footer 106)
export const GE_EL_NAME: string[] = ['Number of Transaction Sets Included', 'Group Control Number'] // Transmittal 907, PDF p.238 (footer 106)
export const IEA_EL_NAME: string[] = ['Number of Included Functional Groups', 'Interchange Control Number'] // Transmittal 907, PDF p.239 (footer 107)

// GS01/ST01 transaction-type qualifier, same combined name either way, as in x12Names835.ts TRANSACTION_TYPE_NAME.
export const TRANSACTION_TYPE_NAME: Record<string, string> = {
  HC: 'Health Care Claim', // Transmittal 907 (GS01 value shown in the CEDI worked example, GS*HC*...); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3)
  '837': 'Health Care Claim', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, ST*837*...)
}

// ===========================================================================
// 1000A/1000B: Submitter and Receiver
// ===========================================================================

export const LOOP_NM1_NAME: Partial<Record<string, string>> = {
  '1000A': 'Submitter Name', // Transmittal 907, PDF p.135 (footer 3)
  '1000B': 'Receiver Name', // Transmittal 907, PDF p.136 (footer 4)
  '2010AA': 'Billing Provider Name', // Transmittal 907, PDF p.138 (footer 6)
  '2010AB': 'Pay-To Address Name', // Transmittal 907, PDF p.140 (footer 8)
  '2010AC': 'Pay-To Plan Name', // Transmittal 907, PDF p.141 (footer 9) - not used by Medicare, see NOT_USED_BY_MEDICARE
  '2010BA': 'Subscriber Name', // Transmittal 907, PDF p.144 (footer 12)
  '2010BB': 'Payer Name', // Transmittal 907, PDF p.147 (footer 15)
  '2010CA': 'Patient Name', // Transmittal 907, PDF p.150 (footer 18) - loop not used by Medicare, see NOT_USED_BY_MEDICARE
  '2310A': 'Referring Provider Name', // Transmittal 907, PDF p.184 (footer 52)
  '2310B': 'Rendering Provider Name', // Transmittal 907, PDF p.185 (footer 53)
  '2310C': 'Service Facility Location Name', // Transmittal 907, PDF p.186 (footer 54)
  '2310D': 'Supervising Provider Name', // Transmittal 907, PDF p.188 (footer 56)
  '2310E': 'Ambulance Pick Up Location', // Transmittal 907, PDF p.189 (footer 57)
  '2310F': 'Ambulance Drop Off Location', // Transmittal 907, PDF p.189-190 (footer 57-58)
  '2330A': 'Other Subscriber Name', // Transmittal 907, PDF p.195 (footer 63)
  '2330B': 'Other Payer Name', // Transmittal 907, PDF p.197 (footer 65)
  '2330C': 'Other Payer Referring Provider', // Transmittal 907, PDF p.200 (footer 68) - not used by Medicare
  '2330D': 'Other Payer Rendering Provider', // Transmittal 907, PDF p.201 (footer 69) - not used by Medicare
  '2330E': 'Other Payer Service Facility Location', // Transmittal 907, PDF p.202 (footer 70) - not used by Medicare
  '2330F': 'Other Payer Supervising Provider', // Transmittal 907, PDF p.203 (footer 71) - not used by Medicare
  '2330G': 'Other Payer Billing Provider', // Transmittal 907, PDF p.204 (footer 72) - not used by Medicare
  '2420A': 'Rendering Provider Name', // Transmittal 907, PDF p.225 (footer 93)
  '2420B': 'Purchased Service Provider Name', // Transmittal 907, PDF p.226 (footer 94)
  '2420C': 'Service Facility Location Name', // Transmittal 907, PDF p.227 (footer 95)
  '2420D': 'Supervising Provider Name', // Transmittal 907, PDF p.229 (footer 97)
  '2420E': 'Ordering Provider Name', // Transmittal 907, PDF p.230 (footer 98)
  '2420F': 'Referring Provider Name', // Transmittal 907, PDF p.231 (footer 99)
  '2420G': 'Ambulance Pick Up Location', // Transmittal 907, PDF p.232 (footer 100)
  '2420H': 'Ambulance Drop Off Location', // Transmittal 907, PDF p.233 (footer 101)
}

// NM101/NM102 own element-position names are the same at every NM1 instance.
export const NM101_EL_NAME = 'Entity Identifier Code' // Transmittal 907, every NM1 instance
export const NM102_EL_NAME = 'Entity Type Qualifier' // Transmittal 907, every NM1 instance
// NM103-NM109 by loop. NM106 (Name Prefix) is Not Used at every 837P NM1 and
// is omitted; array order is [Last/Org Name, First Name, Middle Name/Initial,
// Suffix, ID Code Qualifier, Primary Identifier].
export const NM1_ROLE_EL_NAME: Partial<Record<string, string[]>> = {
  '1000A': ['Submitter Last or Organization Name', 'Submitter First Name', 'Submitter Middle Name or Initial', 'Name Suffix (Not Used)', 'Identification Code Qualifier', 'Submitter Identifier'], // Transmittal 907, PDF p.135
  '1000B': ['Receiver Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Receiver Primary Identifier'], // Transmittal 907, PDF p.136
  '2010AA': ['Billing Provider Last or Organizational Name', 'Billing Provider First Name', 'Billing Provider Middle Name or Initial', 'Billing Provider Name Suffix', 'Identification Code Qualifier', 'Billing Provider Identifier'], // Transmittal 907, PDF p.138
  '2010AB': ['Pay-to Provider Last or Organization Name', 'Pay-to Provider First Name', 'Pay-to Provider Middle Name', 'Pay-to Provider Name Suffix', '(Not Used)', '(Not Used)'], // Transmittal 907, PDF p.140
  '2010AC': ['Pay to Plan Organizational Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Pay to Plan Primary ID'], // Transmittal 907, PDF p.141-142
  '2010BA': ['Subscriber Last Name', 'Subscriber First Name', 'Subscriber Middle Name or Initial', 'Subscriber Name Suffix', 'Identification Code Qualifier', 'Subscriber Primary Identifier'], // Transmittal 907, PDF p.144-145
  '2010BB': ['Payer Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Payer Identifier'], // Transmittal 907, PDF p.147
  '2010CA': ['Patient Last Name', 'Patient First Name', 'Patient Middle Name or Initial', 'Patient Name Suffix', 'Identification Code Qualifier', 'Patient Primary Identifier'], // Transmittal 907, PDF p.150
  '2310A': ['Referring Provider Last Name', 'Referring Provider First Name', 'Referring Provider Middle Name or Initial', 'Referring Provider Name Suffix', 'Identification Code Qualifier', 'Referring Provider Identifier'], // Transmittal 907, PDF p.184
  '2310B': ['Rendering Provider Last or Organization Name', 'Rendering Provider First Name', 'Rendering Provider Middle Name or Initial', 'Rendering Provider Name Suffix', 'Identification Code Qualifier', 'Rendering Provider Identifier'], // Transmittal 907, PDF p.185
  '2310C': ['Laboratory or Facility Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Laboratory or Facility Primary Identifier'], // Transmittal 907, PDF p.186
  '2310D': ['Supervising Provider Last Name', 'Supervising Provider First Name', 'Supervising Provider Middle Name or Initial', 'Supervising Provider Name Suffix', 'Identification Code Qualifier', 'Supervising Provider Identifier'], // Transmittal 907, PDF p.188
  '2310E': ['Name Last or Organization Name', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)'], // Transmittal 907, PDF p.189
  '2310F': ['Ambulance Drop Off Location', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)'], // Transmittal 907, PDF p.189-190
  '2330A': ['Other Insured Last Name', 'Other Insured First Name', 'Other Insured Middle Name', 'Other Insured Name Suffix', 'Identification Code Qualifier', 'Other Insured Identifier'], // Transmittal 907, PDF p.195
  '2330B': ['Other Payer Organization Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Other Payer Primary Identifier'], // Transmittal 907, PDF p.197
  '2420A': ['Rendering Provider Last or Organization Name', 'Rendering Provider First Name', 'Rendering Provider Middle Name or Initial', 'Rendering Provider Name Suffix', 'Identification Code Qualifier', 'Rendering Provider Identifier'], // Transmittal 907, PDF p.225
  '2420B': ['Name Last or Organization Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Purchase Service Primary Identifier'], // Transmittal 907, PDF p.226
  '2420C': ['Laboratory or Facility Name', '(Not Used)', '(Not Used)', '(Not Used)', 'Identification Code Qualifier', 'Laboratory or Facility Primary Identifier'], // Transmittal 907, PDF p.227
  '2420D': ['Supervising Provider Last Name', 'Supervising Provider First Name', 'Supervising Provider Middle Name or Initial', 'Supervising Provider Name Suffix', 'Identification Code Qualifier', 'Supervising Provider Identifier'], // Transmittal 907, PDF p.229
  '2420E': ['Ordering Provider Last Name', 'Ordering Provider First Name', 'Ordering Provider Middle Name or Initial', 'Ordering Provider Name Suffix', 'Identification Code Qualifier', 'Ordering Provider Identifier'], // Transmittal 907, PDF p.230
  '2420F': ['Name Last or Organization Name', 'Referring Provider First Name', 'Referring Provider Middle Name or Initial', 'Referring Provider Name Suffix', 'Identification Code Qualifier', 'Referring Provider Identifier'], // Transmittal 907, PDF p.231
  '2420G': ['Name Last or Organization Name', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)'], // Transmittal 907, PDF p.232
  '2420H': ['Ambulance Drop Off Location', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)', '(Not Used)'], // Transmittal 907, PDF p.233
}

// NM101 code values. Transmittal 907 names every NM1 instance's loop/box
// (NM1_ROLE_EL_NAME above) but does not restate these qualifier characters
// in its own Values column for most of them. Five are confirmed instead in
// the CEDI guide's own worked 837P transaction (section 11.2, Figure 3,
// cited individually below); the rest are confirmed in Transmittal 579's
// Values column, the predecessor flat file for the same 005010 version (see
// the file header). DN and P3 are an ambiguity in the sources:
// both are Transmittal 579's listed values for the same three Referring
// Provider boxes (2310A, 2330C, 2420F) and neither source says which code is
// which, so both keep the same name. Every other loop in LOOP_NM1_NAME keeps
// its box title as the only public name; its NM101 code value is not
// independently confirmed in either source and is not guessed here.
export const NM101_CONFIRMED_NAME: Partial<Record<string, string>> = {
  41: 'Submitter', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "NM1*41*2*EDIT TEST SUBMITTER"); Transmittal 579, PDF p.8 (1000A)
  40: 'Receiver', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "NM1*40*2*DME MAC JURISDICTION B"); Transmittal 579, PDF p.9 (1000B)
  85: 'Billing Provider', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "NM1*85*2*BILLING PROVIDER"); Transmittal 579, PDF p.11 (2010AA) and p.77 (2330G, Other Payer Billing Provider)
  87: 'Pay-to Provider', // Transmittal 579, PDF p.13 (2010AB, Values column) - the box's only NM101 value
  PE: 'Pay-to Plan', // Transmittal 579, PDF p.14 (2010AC, Values column) - the box's only NM101 value; loop not used by Medicare
  PR: 'Payer', // Transmittal 579, PDF p.20 (2010BB) and p.69 (2330B, Values column) - each box's only NM101 value
  IL: 'Insured or Subscriber', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "NM1*IL*1*LASTNAME*FIRSTNAME"); Transmittal 579, PDF p.17 (2010BA) and p.68 (2330A, Other Subscriber); same code x12Names835.ts does not name (835 uses N1/NM1 PR/PE/QC instead)
  QC: 'Patient', // Transmittal 579, PDF p.23 (2010CA, Values column) - the box's only NM101 value; loop not used by Medicare
  DN: 'Referring Provider (DN or P3, ambiguous)', // Transmittal 579, PDF p.57, 73, 104 (2310A/2330C/2420F, Values column "DN, P3"); the source does not say which code is which
  P3: 'Referring Provider (DN or P3, ambiguous)', // Transmittal 579, PDF p.57, 73, 104 (2310A/2330C/2420F, Values column "DN, P3"); the source does not say which code is which
  82: 'Rendering Provider', // Transmittal 579, PDF p.58, 97, 74 (2310B/2420A Rendering Provider boxes and 2330D Other Payer Rendering Provider, Values column) - each box's only NM101 value
  77: 'Service Facility Location', // Transmittal 579, PDF p.59, 100, 75 (2310C/2420C Service Facility Location boxes and 2330E Other Payer Service Facility Location, Values column) - each box's only NM101 value
  DQ: 'Supervising Provider', // Transmittal 579, PDF p.61, 101, 76 (2310D/2420D Supervising Provider boxes and 2330F Other Payer Supervising Provider, Values column) - each box's only NM101 value
  PW: 'Ambulance Pick Up Location', // Transmittal 579, PDF p.62, 105 (2310E/2420G, Values column) - each box's only NM101 value
  45: 'Ambulance Drop Off Location', // Transmittal 579, PDF p.63, 106 (2310F/2420H, Values column) - each box's only NM101 value
  QB: 'Purchased Service Provider', // Transmittal 579, PDF p.99 (2420B, Values column) - the box's only NM101 value
  DK: 'Ordering Provider', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "NM1*DK*1*PROVIDER*ORDERING"); Transmittal 579, PDF p.102 (2420E)
}
export const NM108_CONFIRMED_NAME: Partial<Record<string, string>> = {
  46: 'Electronic Transmitter Identification Number', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, submitter/receiver "*****46*B08XXXXXX")
  XX: 'National Provider Identifier', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, billing/ordering provider "*****XX*1234567893"); same meaning as x12Names835.ts ID_QUALIFIER_NAME.XX
  MI: 'Member Identification Number', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, subscriber "****MI*##########"); matches CEDI 837P CG, PDF p.33 (printed p.25; Table 15) "NM108 Subscriber Identification Code Qualifier MI"
  PI: 'Payer Identification', // CEDI 837P CG, PDF p.35 (printed p.27; Table 17) "NM108 Payer Identification Code Qualifier PI"
}

export const LOOP_N3_NAME: Partial<Record<string, string>> = {
  '2010AA': 'Billing Provider Address', // Transmittal 907, PDF p.138 (footer 6)
  '2010AB': 'Pay-to Provider Address', // Transmittal 907, PDF p.140-141 (footer 8-9)
  '2010AC': 'Pay-to Plan Address', // Transmittal 907, PDF p.142 (footer 10)
  '2010BA': 'Subscriber Address', // Transmittal 907, PDF p.145 (footer 13)
  '2010BB': 'Payer Address', // Transmittal 907, PDF p.148 (footer 16)
  '2010CA': 'Patient Address', // Transmittal 907, PDF p.150 (footer 18)
  '2310C': 'Service Facility Location Address', // Transmittal 907, PDF p.187 (footer 55)
  '2310E': 'Ambulance Pick Up Location Address', // Transmittal 907, PDF p.189 (footer 57)
  '2310F': 'Ambulance Drop Off Location Address', // Transmittal 907, PDF p.190 (footer 58)
  '2330A': 'Other Subscriber Address', // Transmittal 907, PDF p.196 (footer 64)
  '2330B': 'Other Payer Address', // Transmittal 907, PDF p.198 (footer 66)
  '2420C': 'Service Facility Location Address', // Transmittal 907, PDF p.228 (footer 96)
  '2420E': 'Ordering Provider Address', // Transmittal 907, PDF p.230 (footer 98)
  '2420G': 'Ambulance Pick Up Location Address', // Transmittal 907, PDF p.233 (footer 101)
  '2420H': 'Ambulance Drop Off Location Address', // Transmittal 907, PDF p.234 (footer 102)
}
export const LOOP_N4_NAME: Partial<Record<string, string>> = {
  '2010AA': 'Billing Provider City/State/ZIP Code', // Transmittal 907, PDF p.138-139 (footer 6-7)
  '2010AB': 'Pay-to Provider City/State/ZIP Code', // Transmittal 907, PDF p.141 (footer 9)
  '2010AC': 'Pay-to Plan City/State/ZIP Code', // Transmittal 907, PDF p.142 (footer 10)
  '2010BA': 'Subscriber City/State/ZIP Code', // Transmittal 907, PDF p.145 (footer 13)
  '2010BB': 'Payer City/State/ZIP Code', // Transmittal 907, PDF p.148 (footer 16)
  '2010CA': 'Patient City/State/ZIP Code', // Transmittal 907, PDF p.151 (footer 19-20)
  '2310C': 'Service Facility Location City/State/ZIP Code', // Transmittal 907, PDF p.187 (footer 55)
  '2310E': 'Ambulance Pick Up Location City/State/ZIP', // Transmittal 907, PDF p.189 (footer 57)
  '2310F': 'Ambulance Drop Off Location City/State/ZIP', // Transmittal 907, PDF p.190 (footer 58)
  '2330A': 'Other Subscriber City/State/ZIP Code', // Transmittal 907, PDF p.196 (footer 64)
  '2330B': 'Other Payer City/State/ZIP Code', // Transmittal 907, PDF p.198 (footer 66)
  '2420C': 'Service Facility Location City/State/ZIP', // Transmittal 907, PDF p.228 (footer 96)
  '2420E': 'Ordering Provider City/State/ZIP Code', // Transmittal 907, PDF p.230-231 (footer 98-99)
  '2420G': 'Ambulance Pick Up Location City/State/ZIP', // Transmittal 907, PDF p.233-234 (footer 101-102)
  '2420H': 'Ambulance Drop Off Location City/State/ZIP', // Transmittal 907, PDF p.234 (footer 102)
}
// N3/N4 element-position names are the same shape at every instance.
export const N3_EL_NAME: string[] = ['Address Line', 'Address Line'] // Transmittal 907, every N3 instance
export const N4_EL_NAME: string[] = ['City Name', 'State or Province Code', 'Postal Zone or ZIP Code', 'Country Code'] // N404-407 are usually Not Used (Location Qualifier/Identifier, Country Subdivision Code) except N404/N407 situational; Transmittal 907, every N4 instance

export const SUBMITTER_PER_EL_NAME: string[] = [
  'Contact Function Code', 'Submitter Contact Name',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number', 'Contact Inquiry Reference',
] // Transmittal 907, PDF p.135-136 (footer 3-4), segment box "Submitter EDI Contact Information"

// ===========================================================================
// 2000A/2010AA/AB: Billing Provider
// ===========================================================================

export const HL_EL_NAME: string[] = ['Hierarchical ID Number', 'Hierarchical Parent ID Number', 'Hierarchical Level Code', 'Hierarchical Child Code'] // Transmittal 907, every HL instance (2000A PDF p.136, 2000B PDF p.143, 2000C PDF p.149)
export const LOOP_HL_NAME: Partial<Record<string, string>> = {
  '2000A': 'Billing Provider Hierarchical Level', // Transmittal 907, PDF p.136 (footer 4)
  '2000B': 'Subscriber Hierarchical Level', // Transmittal 907, PDF p.143 (footer 11)
  '2000C': 'Patient Hierarchical Level', // Transmittal 907, PDF p.149 (footer 17) - not used by Medicare, see NOT_USED_BY_MEDICARE
}
// HL03 code values. Transmittal 907 names each loop (LOOP_HL_NAME above) but
// does not restate HL03's own code characters; Transmittal 579's Values
// column does, at each loop's own HL.
export const HL03_NAME: Record<string, string> = {
  20: 'Billing Provider', // Transmittal 579, PDF p.9 (2000A, Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*1**20*1")
  22: 'Subscriber', // Transmittal 579, PDF p.16 (2000B, Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*2*1*22*0")
  23: 'Patient', // Transmittal 579, PDF p.22 (2000C, Values column); loop not used by Medicare, see NOT_USED_BY_MEDICARE
}

export const PRV_EL_NAME: string[] = ['Provider Code', 'Reference Identification Qualifier', 'Provider Taxonomy Code'] // PRV04-06 Not Used; Transmittal 907, every PRV instance
export const LOOP_PRV_NAME: Partial<Record<string, string>> = {
  '2000A': 'Billing Provider Specialty Information', // Transmittal 907, PDF p.137 (footer 5)
  '2310B': 'Rendering Provider Specialty Information', // Transmittal 907, PDF p.185-186 (footer 53-54)
  '2420A': 'Rendering Provider Specialty Information', // Transmittal 907, PDF p.225 (footer 93)
}

export const CUR_EL_NAME: string[] = ['Entity Identifier Code', 'Currency Code'] // CUR03-21 Not Used; Transmittal 907, PDF p.137 (footer 5)
export const CUR_NAME = 'Foreign Currency Information' // Transmittal 907, PDF p.137 (footer 5); not supported by Medicare, see NOT_USED_BY_MEDICARE

export const BILLING_PROVIDER_REF_NAME = 'Billing Provider Tax Identification' // Transmittal 907, PDF p.139 (footer 7)
export const BILLING_PROVIDER_REF_EL_NAME: [string, string] = ['Reference Identification Qualifier', 'Billing Provider Tax Identification Number'] // Transmittal 907, PDF p.139 (footer 7); REF01 = EI confirmed, CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "REF*EI*123456789")
export const BILLING_PROVIDER_UPIN_REF_NAME = 'Billing Provider UPIN/License Information' // Transmittal 907, PDF p.139 (footer 7); not used by Medicare, see NOT_USED_BY_MEDICARE

export const BILLING_PROVIDER_PER_NAME = 'Billing Provider Contact Information' // Transmittal 907, PDF p.139-140 (footer 7-8)
export const BILLING_PROVIDER_PER_EL_NAME: string[] = [
  'Contact Function Code', 'Billing Provider Contact Name',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
] // PER09 Not Used; Transmittal 907, PDF p.139-140 (footer 7-8)

export const PAY_TO_PLAN_SECONDARY_REF_NAME = 'Pay-To Plan Secondary Identification' // Transmittal 907, PDF p.142 (footer 10) - loop not used by Medicare
export const PAY_TO_PLAN_TAX_REF_NAME = 'Pay-To Plan Tax Identification' // Transmittal 907, PDF p.143 (footer 11) - loop not used by Medicare
export const REF_SECONDARY_ID_EL_NAME: [string, string] = ['Reference Identification Qualifier', 'Secondary Identifier'] // generic shape shared by every "Secondary Identification"/"Secondary Identifier" REF box in this file; Transmittal 907

// ===========================================================================
// 2000B/C, 2010BA/BB/CA: Subscriber and Patient
// ===========================================================================

export const SBR_EL_NAME: string[] = [
  'Payer Responsibility Sequence Number Code', 'Individual Relationship Code', 'Subscriber Group or Policy Number',
  'Subscriber Group Name', 'Insurance Type Code', 'Claim Filing Indicator Code',
] // SBR06-08 Not Used; Transmittal 907, PDF p.143-144 (footer 11-12)
export const OTHER_SUBSCRIBER_SBR_EL_NAME: string[] = [
  'Payer Responsibility Sequence Number Code', 'Individual Relationship Code', 'Insured Group or Policy Number',
  'Other Insured Group Name', 'Insurance Type Code', 'Claim Filing Indicator Code',
] // Transmittal 907, PDF p.191-192 (footer 59-60), loop 2320 "Other Subscriber Information"

// Medicare restricts 2000B SBR01 to P or S, and because the subscriber is
// always the same person as the patient under Medicare, SBR02 is always 18
// (Self) and SBR09 is always MB (Medicare Part B). CEDI 837P CG, PDF p.33
// (printed p.25; Table 13), confirmed in the worked example "SBR*P*18*...MB"
// (CEDI 837P CG, PDF p.50, printed p.42, Figure 3).
export const SBR01_NAME: Record<string, string> = {
  P: 'Primary', // CEDI 837P CG, PDF p.33 (printed p.25; Table 13)
  S: 'Secondary', // CEDI 837P CG, PDF p.33 (printed p.25; Table 13)
}
export const SBR02_NAME: Record<string, string> = {
  '18': 'Self', // CEDI 837P CG, PDF p.33 (printed p.25; Table 13); the only value Medicare's own subscriber loop sends
}
export const SBR09_NAME: Record<string, string> = {
  MB: 'Medicare Part B', // CEDI 837P CG, PDF p.33 (printed p.25; Table 13); loop 2320's SBR09 cannot equal MA or MB (PDF p.40, printed p.32, Table 27)
}

export const PAT_EL_NAME: string[] = [
  'Individual Relationship Code', 'Patient Location Code', 'Employment Status Code', 'Student Status Code',
  'Date Time Period Format Qualifier', 'Patient Death Date', 'Unit or Basis for Measurement Code', 'Patient Weight', 'Pregnancy Indicator',
] // PAT01-04 Not Used in 2000B (Medicare: subscriber is always the patient); Transmittal 907, PDF p.144 (footer 12); 2000C PAT not used by Medicare at all, see NOT_USED_BY_MEDICARE
export const PAT_NAME = 'Patient Information' // Transmittal 907, PDF p.144 (footer 12)

export const DMG_EL_NAME: string[] = ['Date Time Period Format Qualifier', 'Birth Date', 'Gender Code'] // DMG04-11 Not Used (Marital Status, Race/Ethnicity, Citizenship Status, Country Code, Basis of Verification, Quantity, Code List Qualifier Code, Industry Code); Transmittal 907, PDF p.145-146 (footer 13-14)
export const DMG_NAME = 'Subscriber Demographic Information' // Transmittal 907, PDF p.145 (footer 13)

export const SUBSCRIBER_PC_REF_NAME = 'Property and Casualty Claim Number' // Transmittal 907, PDF p.146 (footer 14)
export const SUBSCRIBER_PC_PER_NAME = 'Property and Casualty Subscriber Contact Information' // Transmittal 907, PDF p.146-147 (footer 14-15)
export const SUBSCRIBER_SECONDARY_REF_NAME = 'Subscriber Secondary Identification' // Transmittal 907, PDF p.146 (footer 14); not used by Medicare, see NOT_USED_BY_MEDICARE

export const PAYER_SECONDARY_REF_NAME = 'Payer Secondary Identification' // Transmittal 907, PDF p.148 (footer 16); not used by Medicare, see NOT_USED_BY_MEDICARE
export const PAYER_LOOP_BILLING_PROVIDER_REF_NAME = 'Billing Provider Secondary Identification' // Transmittal 907, PDF p.148-149 (footer 16-17), appears inside loop 2010BB; not used by Medicare, see NOT_USED_BY_MEDICARE

export const PATIENT_PC_REF_NAME = 'Property and Casualty Patient Identifier' // Transmittal 907, PDF p.152-153 (footer 20-21); loop 2010CA is itself not used by Medicare
export const PATIENT_PC_PER_NAME = 'Property and Casualty Patient Contact Information' // Transmittal 907, PDF p.153 (footer 21)

// Every loop this file marks "not used by Medicare" is drawn from CEDI 837P
// CG, PDF p.26-27 (printed p.18-19; Table 5 "Segment/Elements Not Accepted by
// Medicare"). A claim built with any of these present is rejected, not
// merely ignored.
export const NOT_USED_BY_MEDICARE: string[] = [
  '2000A CUR (Foreign Currency Information)',
  '2010AA REF (Billing Provider UPIN/License Information)',
  '2010AC (Pay-to Plan Loop, in its entirety)',
  '2010BA REF (Subscriber Secondary Identification, REF01=SY)',
  '2010BB REF (Payer Secondary Identification)',
  '2010BB REF (Billing Provider Secondary Identification)',
  '2000C (Patient Hierarchical Level, in its entirety - the subscriber is always the patient)',
  '2000C PAT (Patient Information, in its entirety)',
  '2010CA (Patient Name Loop, in its entirety)',
  '2300 CN1 (Contract Information)',
  '2300 REF (Mandatory Medicare [Section 4081] Crossover Indicator)',
  '2300 REF (Payer Claim Control Number)',
  '2330C (Other Payer Referring Provider Loop, in its entirety)',
  '2330D (Other Payer Rendering Provider Loop, in its entirety)',
  '2330E (Other Payer Service Facility Location Loop, in its entirety)',
  '2330F (Other Payer Supervising Provider Loop, in its entirety)',
  '2330G (Other Payer Billing Provider Loop, in its entirety)',
  '2400 CN1 (Contract Information)',
  '2400/2300 HCP (Line/Claim Pricing/Repricing Information)',
] // CEDI 837P CG, PDF p.26-27 (printed p.18-19; Table 5)

// ===========================================================================
// 2300: Claim
// ===========================================================================

export const CLM_NAME = 'Claim Information' // Transmittal 907, PDF p.155 (footer 23)
export const CLM_EL_NAME: string[] = [
  'Patient Control Number', 'Total Claim Charge Amount', 'Claim Filing Indicator Code', 'Non-Institutional Claim Type Code',
  'Place of Service Code', 'Facility Code Qualifier', 'Claim Frequency Code', // CLM05-1/05-2/05-3
  'Provider or Supplier Signature Indicator', 'Assignment or Plan Participation Code', 'Benefits Assignment Certification Indicator',
  'Release of Information Code', 'Patient Signature Source Code',
  'Related Causes Code', 'Related Causes Code', 'Auto Accident State or Province Code', 'Country Code', // CLM11-1/11-2/11-4/11-5 (CLM11-3 Not Used)
  'Special Program Indicator', 'Yes/No Condition or Response Code', 'Level of Service Code', 'Yes/No Condition or Response Code',
  'Participation Agreement', 'Claim Status Code', 'Yes/No Condition or Response Code', 'Claim Submission Reason Code', 'Delay Reason Code',
] // CLM03/04/13-19 Not Used (CLM03 Claim Filing Indicator Code lives in SBR09 instead); Transmittal 907, PDF p.155-157 (footer 23-25)
// Medicare restricts CLM05-3 to 1 (Original) and ignores any value submitted
// in CLM20. CEDI 837P CG, PDF p.36 (printed p.28; Table 18).
export const CLM05_3_NAME: Record<string, string> = { '1': 'Original' } // CEDI 837P CG, PDF p.36 (printed p.28; Table 18)

// 2300 DTP boxes. DTP01/DTP02 are always "Date/Time Qualifier" and "Date Time
// Period Format Qualifier"; DTP03's name is the box's own title, keyed here
// by a short slug since neither source restates DTP01's numeric code for
// most of these boxes (the one exception, Admission Date = qualifier 435, is
// called out below with its own citation).
export const DTP01_EL_NAME = 'Date/Time Qualifier' // Transmittal 907, every DTP instance
export const DTP02_EL_NAME = 'Date Time Period Format Qualifier' // Transmittal 907, every DTP instance
export const CLAIM_DTP_NAME: Record<string, string> = {
  onsetOfIllness: 'Date - Onset of Current Illness/Symptom', // Transmittal 907, PDF p.155 (footer 23)
  initialTreatment: 'Date - Initial Treatment', // Transmittal 907, PDF p.155 (footer 23)
  lastSeen: 'Date - Last Seen Date', // Transmittal 907, PDF p.155 (footer 23)
  acuteManifestation: 'Date - Acute Manifestation', // Transmittal 907, PDF p.155-156 (footer 23-24)
  accident: 'Date - Accident', // Transmittal 907, PDF p.156 (footer 24)
  lastMenstrualPeriod: 'Date - Last Menstrual Period', // Transmittal 907, PDF p.156 (footer 24)
  lastXRay: 'Date - Last X-Ray Date', // Transmittal 907, PDF p.156 (footer 24)
  hearingVisionPrescription: 'Date - Hearing and Vision Prescription Date', // Transmittal 907, PDF p.156 (footer 24)
  disability: 'Date - Disability Dates', // Transmittal 907, PDF p.156-157 (footer 24-25)
  lastWorked: 'Date - Last Worked', // Transmittal 907, PDF p.157 (footer 25)
  authorizedReturnToWork: 'Date - Authorized Return to Work', // Transmittal 907, PDF p.157 (footer 25)
  admission: 'Date - Admission', // Transmittal 907, PDF p.157 (footer 25); DTP01 = 435, CEDI 837P CG, PDF p.37 (printed p.29, below Table 18): "DTP01 = '435' must be present" when 2400 SV105 is 21, 51 or 61
  discharge: 'Date - Discharge', // Transmittal 907, PDF p.157 (footer 25)
  assumedOrRelinquishedCare: 'Date - Assumed and Relinquished Care Dates', // Transmittal 907, PDF p.157-158 (footer 25-26); repeats twice, one occurrence per date
  propertyCasualtyFirstContact: 'Date - Property and Casualty Date of First Contact', // Transmittal 907, PDF p.158 (footer 26)
  repricerReceived: 'Date - Repricer Received Date', // Transmittal 907, PDF p.158 (footer 26)
}
export const CLAIM_DTP03_EL_NAME: Record<string, string> = {
  onsetOfIllness: 'Onset of Current Illness or Injury Date',
  initialTreatment: 'Initial Treatment Date',
  lastSeen: 'Last Seen Date',
  acuteManifestation: 'Acute Manifestation Date',
  accident: 'Accident Date',
  lastMenstrualPeriod: 'Last Menstrual Period Date',
  lastXRay: 'Last X-Ray Date',
  hearingVisionPrescription: 'Prescription Date',
  disability: 'Disability From Date',
  lastWorked: 'Last Worked Date',
  authorizedReturnToWork: 'Work Return Date',
  admission: 'Related Hospitalization Admission Date',
  discharge: 'Related Hospitalization Discharge Date',
  assumedOrRelinquishedCare: 'Assumed or Relinquished Care Date',
  propertyCasualtyFirstContact: 'Property and Casualty First Contact Date',
  repricerReceived: 'Repricer Received Date',
} // Transmittal 907, same pages as CLAIM_DTP_NAME above (one row per box)

export const CLAIM_PWK_NAME = 'Claim Supplemental Information' // Transmittal 907, PDF p.158 (footer 26); repeats up to 10 times
export const SERVICE_PWK_NAME = 'Line Supplemental Information' // Transmittal 907, PDF p.209 (footer 77); repeats up to 10 times
export const PWK_EL_NAME: string[] = ['Attachment Report Type Code', 'Attachment Transmission Code', 'Identification Code Qualifier', 'Attachment Control Number'] // PWK03/04/07-09 Not Used; Transmittal 907, PDF p.158 (footer 26), p.209 (footer 77)
// Medicare restricts PWK02 to these four values. CEDI 837P CG, PDF p.38
// (printed p.30; Table 21). Only the first PWK iteration at either the claim
// or line level is used in adjudication.
export const PWK02_MEDICARE_CODES: string[] = ['BM', 'FX', 'EL', 'FT'] // CEDI 837P CG, PDF p.38 (printed p.30; Table 21) - the guide gives these four characters without expanding them to English, so no name table is given here per this file's citation rule

export const CN1_NAME = 'Contract Information' // Transmittal 907, PDF p.159 (footer 27); not used by Medicare at either 2300 or 2400, see NOT_USED_BY_MEDICARE
export const CN1_EL_NAME: string[] = ['Contract Type Code', 'Contract Amount', 'Contract Percentage', 'Contract Code', 'Terms Discount Percent', 'Contract Version Identifier'] // Transmittal 907, PDF p.159 (footer 27)

export const CLAIM_AMT_NAME = 'Patient Amount Paid' // Transmittal 907, PDF p.159 (footer 27)
export const CLAIM_AMT_EL_NAME: [string, string] = ['Amount Qualifier Code', 'Patient Amount Paid'] // AMT03 Not Used; Transmittal 907, PDF p.159 (footer 27)

// 2300 REF boxes, same slug convention as CLAIM_DTP_NAME.
export const CLAIM_REF_NAME: Record<string, string> = {
  serviceAuthException: 'Service Authorization Exception Code', // Transmittal 907, PDF p.159 (footer 27)
  medicareCrossoverIndicator: 'Mandatory Medicare (Section 4081) Crossover Indicator', // Transmittal 907, PDF p.160 (footer 28); not used by Medicare, see NOT_USED_BY_MEDICARE
  mammographyCertification: 'Mammography Certification Number', // Transmittal 907, PDF p.160 (footer 28)
  referral: 'Referral Number', // Transmittal 907, PDF p.160-161 (footer 28-29)
  priorAuthorization: 'Prior Authorization', // Transmittal 907, PDF p.161 (footer 29)
  payerClaimControlNumber: 'Payer Claim Control Number', // Transmittal 907, PDF p.161-162 (footer 29-30); not used by Medicare, see NOT_USED_BY_MEDICARE
  clia: 'Clinical Laboratory Improvement Amendment (CLIA) Number', // Transmittal 907, PDF p.162 (footer 30)
  repricedClaim: 'Repriced Claim Number', // Transmittal 907, PDF p.162-163 (footer 30-31)
  adjustedRepricedClaim: 'Adjusted Repriced Claim Number', // Transmittal 907, PDF p.163 (footer 31)
  investigationalDeviceExemption: 'Investigational Device Exemption Number', // Transmittal 907, PDF p.163-164 (footer 31-32)
  claimIdentifierForTransmissionIntermediaries: 'Claim Identifier for Transmission Intermediaries', // Transmittal 907, PDF p.164 (footer 32)
  medicalRecordNumber: 'Medical Record Number', // Transmittal 907, PDF p.164-165 (footer 32-33)
  demonstrationProject: 'Demonstration Project Identifier', // Transmittal 907, PDF p.165 (footer 33)
  carePlanOversight: 'Care Plan Oversight', // Transmittal 907, PDF p.165 (footer 33)
}
export const CLAIM_REF02_EL_NAME: Record<string, string> = {
  serviceAuthException: 'Service Authorization Exception Code',
  medicareCrossoverIndicator: 'Medicare Section 4081 Indicator',
  mammographyCertification: 'Mammography Certification Number',
  referral: 'Referral Number',
  priorAuthorization: 'Prior Authorization or Referral Number',
  payerClaimControlNumber: 'Payer Claim Control Number',
  clia: 'Clinical Laboratory Improvement Amendment Number',
  repricedClaim: 'Repriced Claim Reference Number',
  adjustedRepricedClaim: 'Adjusted Repriced Claim Reference Number',
  investigationalDeviceExemption: 'Investigational Device Exemption Number',
  claimIdentifierForTransmissionIntermediaries: 'Value Added Network Trace Number',
  medicalRecordNumber: 'Medical Record Number',
  demonstrationProject: 'Demonstration Project Identifier',
  carePlanOversight: 'Care Plan Oversight Number',
} // Transmittal 907, same pages as CLAIM_REF_NAME above

export const K3_EL_NAME = 'Fixed Format Information' // K302/K303 Not Used; Transmittal 907, PDF p.166 (footer 34), p.221 (footer 88)
export const CLAIM_K3_NAME = 'File Information' // Transmittal 907, PDF p.166 (footer 34); repeats up to 10 times
export const SERVICE_K3_NAME = 'File Information' // Transmittal 907, PDF p.221 (footer 88); repeats up to 10 times

export const CLAIM_NTE_NAME = 'Claim Note' // Transmittal 907, PDF p.166 (footer 34)
export const NTE_EL_NAME: [string, string] = ['Note Reference Code', 'Note Text'] // Transmittal 907, every NTE instance

export const CLAIM_CR1_NAME = 'Ambulance Transport Information' // Transmittal 907, PDF p.166 (footer 34)
export const SERVICE_CR1_NAME = 'Ambulance Transport Information' // Transmittal 907, PDF p.210 (footer 78)
// Position labels confirmed directly against the flat-file rows (Transmittal
// 907, PDF p.166-167): CR101, CR102, CR104 (Ambulance Transport Reason
// Code), CR105, CR106 (Transport Distance), CR109, CR110; CR103 and CR107-08
// are Not Used and are skipped, matching docs/design/837p-examples.md's
// "CR104 (transport reason), CR106 (distance)" and AmbulanceTransportCR1's
// field names (claim.go).
export const CR1_EL_NAME: string[] = [
  'Unit or Basis for Measurement Code', 'Patient Weight', 'Ambulance Transport Reason Code', // CR101, CR102, CR104
  'Unit or Basis for Measurement Code', 'Transport Distance', 'Round Trip Purpose Description', 'Stretcher Purpose Description', // CR105, CR106, CR109, CR110
] // CR103, CR107-08 Not Used; Transmittal 907, PDF p.166-167 (footer 34-35), p.210 (footer 78)

export const CR2_NAME = 'Spinal Manipulation Service Information' // Transmittal 907, PDF p.167 (footer 35)
export const CR2_EL_NAME: string[] = ['Patient Condition Code', 'Patient Condition Description', 'Patient Condition Description'] // CR201-07, CR209, CR212 Not Used; Transmittal 907, PDF p.167 (footer 35)

export const CLAIM_CRC_NAME: Record<string, string> = {
  ambulanceCertification: 'Ambulance Certification', // Transmittal 907, PDF p.168 (footer 36); repeats up to 3 times
  visionCondition: 'Patient Condition Information: Vision', // Transmittal 907, PDF p.168 (footer 36); repeats up to 3 times
  homebound: 'Homebound Indicator', // Transmittal 907, PDF p.168 (footer 36)
  epsdtReferral: 'EPSDT Referral', // Transmittal 907, PDF p.169 (footer 37)
}
export const CRC_EL_NAME: string[] = ['Code Category', 'Certification Condition Indicator', 'Condition Indicator', 'Condition Indicator', 'Condition Indicator', 'Condition Code', 'Condition Code'] // field count/usage varies slightly by box; Transmittal 907, PDF p.168-169 (footer 36-37)

// HI - Health Care Diagnosis Code, Anesthesia Related Procedure, and
// Condition Information each repeat the same 9-part composite (HI0n-1
// through HI0n-9) up to 12 times (HI01-HI12). Only the first two positions
// are ever populated by Medicare; HI0n-3 through HI0n-9 are Not Used. The
// flat file still names each of the nine at every composite of all three
// boxes, so each has a name below (HI_NOT_USED_EL_NAME for HI0n-3 through
// HI0n-9, the same at every box), and backend/pkg/x12837p's
// HealthCareCodeHI keeps a field for each.
export const HI_DIAGNOSIS_NAME = 'Health Care Diagnosis Code' // Transmittal 907, PDF p.169 (footer 37); loop 2300, up to 12 composites
export const HI_ANESTHESIA_NAME = 'Anesthesia Related Procedure' // Transmittal 907, PDF p.174 (footer 42); loop 2300, up to 12 composites
export const HI_CONDITION_NAME = 'Condition Information' // Transmittal 907, PDF p.178 (footer 46); loop 2300, up to 12 composites
// HI0n-3 through HI0n-9, in position order; all Not Used (N/U). The PDF prints
// two of these with one lowercase word ("Industry code", "Yes/No Condition or
// response Code"); they are capitalized here as this file's other uses of the
// same two elements are (DMG11, CLM13).
export const HI_NOT_USED_EL_NAME: string[] = [
  'Date Time Period Format Qualifier', 'Date Time Period', 'Monetary Amount', 'Quantity', 'Version Identifier', 'Industry Code', 'Yes/No Condition or Response Code',
] // Transmittal 907, PDF p.169-183 (footer 37-51): the same seven names at HI01 through HI12 of all three boxes
export const HI_DIAGNOSIS_EL_NAME: string[] = ['Diagnosis Type Code', 'Diagnosis Code', ...HI_NOT_USED_EL_NAME] // HI0n-1, HI0n-2, then HI0n-3 through HI0n-9; Transmittal 907, PDF p.169-174 (footer 37-42)
export const HI_ANESTHESIA_EL_NAME: string[] = ['Code List Qualifier Code', 'Anesthesia Related Surgical Procedure', ...HI_NOT_USED_EL_NAME] // HI0n-1, HI0n-2, then HI0n-3 through HI0n-9; Transmittal 907, PDF p.174-178 (footer 42-46)
export const HI_CONDITION_EL_NAME: string[] = ['Code List Qualifier Code', 'Condition Code', ...HI_NOT_USED_EL_NAME] // HI0n-1, HI0n-2, then HI0n-3 through HI0n-9; Transmittal 907, PDF p.178-183 (footer 46-51)
// HI01-1/HI02-1/... Diagnosis Type Code qualifier values. Transmittal 907
// never restates the ICD-9/ICD-10 qualifier letters; these are the X12
// 005010 element 1270 values, not independently re-confirmed in either
// source reviewed here, so no value table is given - the element keeps its
// own reference (HI0x-1) per this file's citation rule.

export const HCP_NAME = 'Claim Pricing/Repricing Information' // Transmittal 907, PDF p.181-182 (footer 49-50); loop 2300; not used by Medicare, see NOT_USED_BY_MEDICARE
export const LINE_HCP_NAME = 'Line Pricing/Repricing Information' // Transmittal 907, PDF p.217 (footer 85); loop 2400; not used by Medicare, see NOT_USED_BY_MEDICARE
export const HCP_EL_NAME: string[] = [
  'Pricing Methodology', 'Repriced Allowed Amount', 'Repriced Saving Amount', 'Repricing Organization Identifier',
  'Repricing Per Diem or Flat Rate Amount', 'Repriced Approved Ambulatory Patient Group Code', 'Repriced Approved Ambulatory Patient Group Amount',
  'Product or Service ID Qualifier', 'Repriced Approved HCPCS Code', 'Unit or Basis for Measurement Code', 'Repriced Approved Service Unit Count',
  'Reject Reason Code', 'Policy Compliance Code', 'Exception Code',
] // HCP08 Not Used; Transmittal 907, PDF p.181-182 (footer 49-50), p.216-217 (footer 84-85)

// ===========================================================================
// 2310A-F: Claim-Level Providers
// ===========================================================================

export const LOOP_REF_NAME: Partial<Record<string, string>> = {
  '2310A': 'Referring Provider Secondary Identification', // Transmittal 907, PDF p.184-185 (footer 52-53); repeats up to 3 times
  '2310B': 'Rendering Provider Secondary Identification', // Transmittal 907, PDF p.186 (footer 54); repeats up to 4 times
  '2310C': 'Service Facility Location Secondary Identification', // Transmittal 907, PDF p.187-188 (footer 55-56); repeats up to 3 times
  '2310D': 'Supervising Provider Secondary Identifier', // Transmittal 907, PDF p.188-189 (footer 56-57); repeats up to 4 times
  '2330A': 'Other Subscriber Secondary Identification', // Transmittal 907, PDF p.196-197 (footer 64-65)
  '2420A': 'Rendering Provider Secondary Identification', // Transmittal 907, PDF p.225-226 (footer 93-94); repeats up to 20 times
  '2420B': 'Purchased Service Provider Secondary Identification', // Transmittal 907, PDF p.226-227 (footer 94-95); repeats up to 20 times
  '2420C': 'Service Facility Location Secondary Identification', // Transmittal 907, PDF p.228 (footer 96); repeats up to 3 times
  '2420D': 'Supervising Provider Secondary Identification', // Transmittal 907, PDF p.229-230 (footer 97-98); repeats up to 20 times
  '2420E': 'Ordering Provider Secondary Identification', // Transmittal 907, PDF p.231 (footer 99); repeats up to 20 times
  '2420F': 'Referring Provider Secondary Identification', // Transmittal 907, PDF p.231-232 (footer 99-100); repeats up to 20 times
}
export const LOOP_REF02_EL_NAME: Partial<Record<string, string>> = {
  '2310A': 'Referring Provider Secondary Identifier',
  '2310B': 'Rendering Provider Secondary Identifier',
  '2310C': 'Laboratory or Facility Secondary Identifier',
  '2310D': 'Supervising Provider Secondary Identifier',
  '2330A': 'Other Insured Additional Identifier', // restricted to a 9-digit numeric value with position-level rules, CEDI 837P CG, PDF p.41 (printed p.33; Table 31)
  '2420A': 'Rendering Provider Secondary Identifier',
  '2420B': 'Purchased Service Provider Secondary Identifier',
  '2420C': 'Laboratory or Facility Secondary Identifier',
  '2420D': 'Supervising Provider Secondary Identifier',
  '2420E': 'Ordering Provider Secondary Identifier',
  '2420F': 'Referring Provider Secondary Identifier',
} // Transmittal 907, same pages as LOOP_REF_NAME above

export const AMBULANCE_PICKUP_NM1_EL_NAME: [string] = ['Name Last or Organization Name'] // NM102-NM112 Not Used; Transmittal 907, PDF p.189 (footer 57), loops 2310E/2420G
export const AMBULANCE_DROPOFF_NM1_EL_NAME: [string] = ['Ambulance Drop Off Location'] // NM102-NM112 Not Used; Transmittal 907, PDF p.189-190, 233-234 (footer 57-58, 101-102), loops 2310F/2420H

// ===========================================================================
// 2320/2330A-G: Other Subscriber and Payer (COB)
// ===========================================================================

export const CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME = 'Claim Level Adjustments' // Transmittal 907, PDF p.192 (footer 60); loop 2320; repeats up to 5 times
export const LINE_ADJUSTMENT_CAS_NAME = 'Line Adjustment' // Transmittal 907, PDF p.235 (footer 103); loop 2430; repeats up to 5 times
// CAS carries up to six adjustment reason/amount/quantity triplets per
// occurrence (CAS02-04, 05-07, 08-10, 11-13, 14-16, 17-19); every triplet
// shares the same three position names.
export const CAS01_EL_NAME = 'Claim Adjustment Group Code' // Transmittal 907, every CAS instance
export const CAS_REASON_EL_NAME = 'Adjustment Reason Code' // Transmittal 907, every CAS instance
export const CAS_AMOUNT_EL_NAME = 'Adjustment Amount' // Transmittal 907, every CAS instance
export const CAS_QTY_EL_NAME = 'Adjustment Quantity' // Transmittal 907, every CAS instance
// Claim Adjustment Group Code (CAS01) is the same externally maintained X12
// list the 835 side of this project already uses.
export const CAS_GROUP_NAME: Record<string, string> = {
  CO: 'Contractual Obligation', // x12.org/codes/claim-adjustment-group-codes
  OA: 'Other Adjustment', // x12.org/codes/claim-adjustment-group-codes
  PR: 'Patient Responsibility', // x12.org/codes/claim-adjustment-group-codes
  PI: 'Payer Initiated Reductions', // x12.org/codes/claim-adjustment-group-codes
  CR: 'Correction and Reversals', // x12.org/codes/claim-adjustment-group-codes
}
// CAS segment must not be present in loop 2320 when 2000B SBR01 = P. CEDI
// 837P CG, PDF p.40 (printed p.32; Table 28).

export const COB_PAYER_PAID_AMT_NAME = 'COB Payer Paid Amount' // Transmittal 907, PDF p.194 (footer 62)
export const COB_NON_COVERED_AMT_NAME = 'COB Total Non-Covered Amount' // Transmittal 907, PDF p.194 (footer 62)
export const OTHER_SUBSCRIBER_REMAINING_LIABILITY_AMT_NAME = 'Remaining Patient Liability' // Transmittal 907, PDF p.194 (footer 62); loop 2320
export const LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME = 'Remaining Patient Liability' // Transmittal 907, PDF p.237 (footer 105); loop 2430
export const AMT01_EL_NAME = 'Amount Qualifier Code' // Transmittal 907, every AMT instance
export const AMT02_EL_NAME_BY_BOX: Record<string, string> = {
  cobPayerPaid: 'Payer Paid Amount',
  cobNonCovered: 'Non-Covered Charge Amount',
  remainingPatientLiability: 'Remaining Patient Liability',
} // AMT03 Not Used at every instance; Transmittal 907, PDF p.194 (footer 62), p.237 (footer 105)
// Medicare requires one occurrence of the 2320 AMT with AMT01 = D (Payer
// Paid Amount coordination of benefits) whenever 2000B SBR01 = S. CEDI 837P
// CG, PDF p.40 (printed p.32; Table 29).
export const AMT01_COB_NAME: Record<string, string> = { D: 'Payer Amount Paid' } // CEDI 837P CG, PDF p.40 (printed p.32; Table 29)

export const OI_NAME = 'Other Insurance Coverage Information' // Transmittal 907, PDF p.194-195 (footer 62-63)
export const OI_EL_NAME: string[] = ['Benefits Assignment Certification Indicator', 'Patient Signature Source Code', 'Release of Information Code'] // OI01/02/05 Not Used; Transmittal 907, PDF p.194-195 (footer 62-63)

export const MOA_NAME = 'Outpatient Adjudication (MOA) Information' // Transmittal 907, PDF p.195 (footer 63)
export const MOA_EL_NAME: string[] = [
  'Reimbursement Rate', 'HCPCS Payable Amount',
  'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code',
  'End Stage Renal Disease Payment Amount', 'Non-Payable Professional Component Billed Amount',
] // Transmittal 907, PDF p.195 (footer 63)

export const OTHER_PAYER_DTP_NAME = 'Date - Claim Check or Remittance Date' // Transmittal 907, PDF p.198 (footer 66); loop 2330B
export const OTHER_PAYER_DTP03_EL_NAME = 'Adjudication or Payment Date' // Transmittal 907, PDF p.198 (footer 66)
export const LINE_ADJUDICATION_DTP_NAME = 'Line Check or Remittance Date' // Transmittal 907, PDF p.237 (footer 105); loop 2430
export const LINE_ADJUDICATION_DTP03_EL_NAME = 'Adjudication or Payment Date' // Transmittal 907, PDF p.237 (footer 105)

export const OTHER_PAYER_REF_NAME: Record<string, string> = {
  secondaryIdentification: 'Other Payer Secondary Identification', // Transmittal 907, PDF p.198 (footer 66)
  priorAuthorization: 'Other Payer Prior Authorization Number', // Transmittal 907, PDF p.198-199 (footer 66-67)
  referral: 'Other Payer Referral Number', // Transmittal 907, PDF p.199 (footer 67)
  claimAdjustmentIndicator: 'Other Payer Claim Adjustment Indicator', // Transmittal 907, PDF p.199-200 (footer 67-68)
  claimControlNumber: 'Other Payer Claim Control Number', // Transmittal 907, PDF p.200 (footer 68)
}
export const OTHER_PAYER_REF02_EL_NAME: Record<string, string> = {
  secondaryIdentification: 'Other Payer Secondary Identifier',
  priorAuthorization: 'Other Payer Prior Authorization Number',
  referral: 'Other Payer Referral Number',
  claimAdjustmentIndicator: 'Other Payer Claim Adjustment Indicator',
  claimControlNumber: 'Other Payer Claim Control Number',
} // Transmittal 907, same pages as OTHER_PAYER_REF_NAME above

export const OTHER_PAYER_ROLE_REF_NAME: Partial<Record<string, string>> = {
  '2330C': 'Other Payer Referring Provider Secondary Identification', // Transmittal 907, PDF p.200-201 (footer 68-69); not used by Medicare
  '2330D': 'Other Payer Rendering Provider Secondary Identification', // Transmittal 907, PDF p.201-202 (footer 69-70); not used by Medicare
  '2330E': 'Other Payer Service Facility Location Secondary Identification', // Transmittal 907, PDF p.202-203 (footer 70-71); not used by Medicare
  '2330F': 'Other Payer Supervising Provider Secondary Identification', // Transmittal 907, PDF p.203-204 (footer 71-72); not used by Medicare
  '2330G': 'Other Payer Billing Provider Secondary Identification', // Transmittal 907, PDF p.204-205 (footer 72-73); not used by Medicare
}

// ===========================================================================
// 2400: Service Line
// ===========================================================================

export const LX_NAME = 'Service Line' // Transmittal 907, PDF p.205 (footer 73)
export const LX01_EL_NAME = 'Assigned Number' // Transmittal 907, PDF p.205 (footer 73)

export const SV1_NAME = 'Professional Service' // Transmittal 907, PDF p.205 (footer 73)
export const SV1_EL_NAME: string[] = [
  'Product or Service ID Qualifier', 'Procedure Code', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Description', // SV101 composite
  'Line Item Charge Amount', 'Unit or Basis for Measurement Code', 'Service Unit Count', 'Place of Service Code',
  'Diagnosis Code Pointer', 'Diagnosis Code Pointer', 'Diagnosis Code Pointer', 'Diagnosis Code Pointer', // SV107 composite
  'Emergency Indicator', 'EPSDT Indicator', 'Family Planning Indicator', 'Co-Pay Status Code',
] // SV101-8, SV106, SV108, SV110, SV113-21 Not Used; Transmittal 907, PDF p.205-208 (footer 73-76)
export const SV101_1_NAME: Record<string, string> = { HC: 'Healthcare Common Procedure Coding System' } // Medicare requires SV101-1 = HC; CEDI 837P CG, PDF p.42 (printed p.34; Table 33); same code x12Names835.ts SVC_QUALIFIER_NAME.HC names on the 835 side

export const SV5_NAME = 'Durable Medical Equipment Service' // Transmittal 907, PDF p.208 (footer 76)
export const SV5_EL_NAME: string[] = [
  'Product or Service ID Qualifier', 'Procedure Code', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Description', // SV501 composite
  'Unit or Basis for Measurement Code', 'Length of Medical Necessity', 'DME Rental Price', 'DME Purchase Price', 'Rental Unit Price Indicator',
] // SV501-8, SV507 Not Used; Transmittal 907, PDF p.208-209 (footer 76-77)

export const CR3_NAME = 'Durable Medical Equipment Certification' // Transmittal 907, PDF p.210 (footer 78)
export const CR3_EL_NAME: string[] = ['Certification Type Code', 'Unit or Basis for Measurement Code', 'Durable Medical Equipment Duration'] // CR304-05 Not Used; Transmittal 907, PDF p.210 (footer 78); no longer carried in the CEDI 837P CG's own Medicare-specific tables as of version 5.6, January 2025 (PDF p.55, printed p.48; Table 50 change log, "removed Table 36 Loop 2400 CR3 Durable Medical Equipment Certification")

export const SERVICE_CRC_NAME: Record<string, string> = {
  ambulanceCertification: 'Ambulance Certification', // Transmittal 907, PDF p.210-211 (footer 78-79); repeats up to 3 times
  hospiceEmployeeIndicator: 'Hospice Employee Indicator', // Transmittal 907, PDF p.211 (footer 79)
  dmeConditionIndicator: 'Condition Indicator Durable Medical Equipment', // Transmittal 907, PDF p.211 (footer 79); no longer carried in the CEDI 837P CG's own Medicare-specific tables as of version 5.6 (PDF p.55, printed p.48; Table 50 change log, "removed Table 37 Loop 2400 CRC Condition Indicator/Durable Medical Equipment")
}

export const SERVICE_DTP_NAME: Record<string, string> = {
  serviceDate: 'Date - Service Date', // Transmittal 907, PDF p.211-212 (footer 79-80); required on every service line
  prescription: 'Date - Prescription Date', // Transmittal 907, PDF p.212 (footer 80)
  certificationRevision: 'Date - Certification Revision/Recertification Date', // Transmittal 907, PDF p.212 (footer 80); no longer carried in the CEDI 837P CG's own Medicare-specific tables as of version 5.6 (PDF p.55, printed p.48; Table 50)
  beginTherapy: 'Date - Begin Therapy Date', // Transmittal 907, PDF p.212 (footer 80); same v5.6 removal as above
  lastCertification: 'Date - Last Certification Date', // Transmittal 907, PDF p.212 (footer 80); same v5.6 removal as above
  lastTreatmentOrTherapy: 'Date - Date Last Seen Date', // Transmittal 907, PDF p.212-213 (footer 80-81)
  testPerformed: 'Date - Test Date', // Transmittal 907, PDF p.213 (footer 81); repeats up to 2 times
  shipped: 'Date - Shipped Date', // Transmittal 907, PDF p.213 (footer 81)
  lastXRay: 'Date - Last X-Ray Date', // Transmittal 907, PDF p.213 (footer 81)
  initialTreatment: 'Date - Initial Treatment Date', // Transmittal 907, PDF p.213 (footer 81)
}
export const SERVICE_DTP03_EL_NAME: Record<string, string> = {
  serviceDate: 'Service Date',
  prescription: 'Prescription Date',
  certificationRevision: 'Certification Revision or Recertification Date',
  beginTherapy: 'Begin Therapy Date',
  lastCertification: 'Last Certification Date',
  lastTreatmentOrTherapy: 'Last Treatment or Therapy Date',
  testPerformed: 'Test Performed Date',
  shipped: 'Shipped Date',
  lastXRay: 'Last X-Ray Date',
  initialTreatment: 'Initial Treatment Date',
} // Transmittal 907, same pages as SERVICE_DTP_NAME above

export const AMBULANCE_PATIENT_COUNT_QTY_NAME = 'Ambulance Patient Count' // Transmittal 907, PDF p.213-214 (footer 81-82)
export const OBSTETRIC_ANESTHESIA_QTY_NAME = 'Obstetric Anesthesia Additional Units' // Transmittal 907, PDF p.214 (footer 82)
export const QTY01_EL_NAME = 'Quantity Qualifier' // Transmittal 907, every QTY instance
export const QTY02_EL_NAME_BY_BOX: Record<string, string> = {
  ambulancePatientCount: 'Ambulance Patient Count',
  obstetricAdditionalUnits: 'Obstetric Additional Units',
} // Transmittal 907, PDF p.213-214 (footer 81-82)

export const MEA_NAME = 'Test Results' // Transmittal 907, PDF p.214-215 (footer 82-83); repeats up to 5 times
export const MEA_EL_NAME: string[] = ['Measurement Reference Identification Code', 'Measurement Qualifier', 'Test Result'] // MEA04-12 Not Used; Transmittal 907, PDF p.214-215 (footer 82-83)

export const SERVICE_REF_NAME: Record<string, string> = {
  repricedLineItem: 'Repriced Line Item Reference Number', // Transmittal 907, PDF p.215-216 (footer 83-84)
  adjustedRepricedLineItem: 'Adjusted Repriced Line Item Reference Number', // Transmittal 907, PDF p.216 (footer 84)
  priorAuthorization: 'Prior Authorization', // Transmittal 907, PDF p.216 (footer 84); repeats up to 5 times
  lineItemControl: 'Line Item Control Number', // Transmittal 907, PDF p.217 (footer 85); REF01 = 6R confirmed, CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "REF*6R*001H0000000140309201")
  mammographyCertification: 'Mammography Certification Number', // Transmittal 907, PDF p.217 (footer 85)
  clia: 'Clinical Laboratory Improvement Amendment (CLIA) Identification', // Transmittal 907, PDF p.218 (footer 86)
  referringClia: 'Referring Clinical Laboratory Improvement Amendment (CLIA) Facility Identification', // Transmittal 907, PDF p.218 (footer 86)
  immunizationBatch: 'Immunization Batch Number', // Transmittal 907, PDF p.219 (footer 87)
  referral: 'Referral Number', // Transmittal 907, PDF p.219 (footer 87); repeats up to 5 times
}
export const SERVICE_REF02_EL_NAME: Record<string, string> = {
  repricedLineItem: 'Repriced Line Item Reference Number',
  adjustedRepricedLineItem: 'Adjusted Repriced Line Item Reference Number',
  priorAuthorization: 'Prior Authorization or Referral Number',
  lineItemControl: 'Line Item Control Number',
  mammographyCertification: 'Mammography Certification Number',
  clia: 'Clinical Laboratory Improvement Amendment Number',
  referringClia: 'Referring CLIA Number',
  immunizationBatch: 'Immunization Batch Number',
  referral: 'Referral Number',
} // Transmittal 907, same pages as SERVICE_REF_NAME above

export const SALES_TAX_AMT_NAME = 'Sales Tax Amount' // Transmittal 907, PDF p.219-220 (footer 87-88)
export const POSTAGE_CLAIMED_AMT_NAME = 'Postage Claimed Amount' // Transmittal 907, PDF p.220 (footer 88)
export const SERVICE_AMT02_EL_NAME_BY_BOX: Record<string, string> = {
  salesTax: 'Sales Tax Amount',
  postageClaimed: 'Postage Claimed Amount',
} // Transmittal 907, PDF p.219-220 (footer 87-88)

export const LINE_NTE_NAME = 'Line Note' // Transmittal 907, PDF p.220 (footer 88)
export const THIRD_PARTY_NTE_NAME = 'Third Party Organization Notes' // Transmittal 907, PDF p.220 (footer 88)

export const PS1_NAME = 'Purchased Service Information' // Transmittal 907, PDF p.220-221 (footer 88-89)
export const PS1_EL_NAME: [string, string] = ['Purchased Service Provider Identifier', 'Purchased Service Charge Amount'] // PS103 Not Used; Transmittal 907, PDF p.220-221 (footer 88-89)

// ===========================================================================
// 2410: Drug Identification
// ===========================================================================

export const LIN_NAME = 'Drug Identification' // Transmittal 907, PDF p.222 (footer 90)
export const LIN_EL_NAME: [string, string] = ['Product or Service ID Qualifier', 'National Drug Code'] // LIN01, LIN04-31 Not Used; Transmittal 907, PDF p.222 (footer 90)
// For DME claims, loop 2410 is required whenever 2400 SV101-1 is a default
// HCPCS code, and LIN02 must equal N4; LIN03 must be exactly 11 alphanumeric
// characters. CEDI 837P CG, PDF p.46 (printed p.38; Table 39).
export const LIN02_NAME: Record<string, string> = { N4: 'National Drug Code' } // CEDI 837P CG, PDF p.46 (printed p.38; Table 39)

export const CTP_NAME = 'Drug Pricing' // Transmittal 907, PDF p.223 (footer 91)
export const CTP_EL_NAME: string[] = ['National Drug Unit Count', 'Unit or Basis for Measurement Code'] // CTP01-03, CTP05-2 onward, CTP06-11 Not Used; Transmittal 907, PDF p.223-224 (footer 91-92)
// For DME claims, CTP05-1 must equal UN when the loop 2410 LIN03 NDC is
// associated with an Oral Cancer Drug HCPCS code. CEDI 837P CG, PDF p.47
// (printed p.39; Table 41).
export const CTP05_1_NAME: Record<string, string> = { UN: 'Unit' } // CEDI 837P CG, PDF p.47 (printed p.39; Table 41)

export const DRUG_ASSOCIATION_REF_NAME = 'Prescription or Compound Drug Association Number' // Transmittal 907, PDF p.224 (footer 92)
export const DRUG_ASSOCIATION_REF02_EL_NAME = 'Prescription Number' // Transmittal 907, PDF p.224 (footer 92)
// Required with REF01 = XZ whenever the 2400 SV1 segment carries modifier J1.
// CEDI 837P CG, PDF p.47 (printed p.39; Table 41).
export const DRUG_ASSOCIATION_REF01_NAME: Record<string, string> = { XZ: 'Pharmacy Prescription Number' } // CEDI 837P CG, PDF p.47 (printed p.39; Table 41)

// ===========================================================================
// 2420A-H: Line-Level Providers
// ===========================================================================
// See LOOP_NM1_NAME, NM1_ROLE_EL_NAME, LOOP_N3_NAME, LOOP_N4_NAME, and
// LOOP_REF_NAME/LOOP_REF02_EL_NAME above - 2420A-H share those same tables
// with their 2310-series counterparts (2420A Rendering, 2420B Purchased
// Service, 2420C Service Facility Location, 2420D Supervising, 2420E
// Ordering, 2420F Referring, 2420G Ambulance Pick Up, 2420H Ambulance Drop
// Off; only 2420B, 2420E and 2420F have no 2310-series analog).

export const ORDERING_PROVIDER_PER_NAME = 'Ordering Provider Contact Information' // Transmittal 907, PDF p.231 (footer 99)
export const ORDERING_PROVIDER_PER_EL_NAME: string[] = [
  'Contact Function Code', 'Ordering Provider Contact Name',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
] // PER09 Not Used; Transmittal 907, PDF p.231 (footer 99)

export const SERVICE_FACILITY_PER_NAME = 'Service Facility Contact Information' // Transmittal 907, PDF p.188 (footer 56); loop 2310C
export const SERVICE_FACILITY_PER_EL_NAME: string[] = [
  'Contact Function Code', 'Submitter Contact', // Transmittal 907's own column literally repeats the Submitter PER02 label here
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
  'Communication Number Qualifier', 'Communication Number',
] // Transmittal 907, PDF p.188 (footer 56)

// ===========================================================================
// 2430: Line Adjudication Information
// ===========================================================================

export const SVD_NAME = 'Line Adjudication Information' // Transmittal 907, PDF p.235 (footer 103)
export const SVD_EL_NAME: string[] = [
  'Other Payer Primary Identifier', 'Service Line Paid Amount',
  'Product or Service ID Qualifier', 'Procedure Code', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Modifier', 'Procedure Code Description', // SVD03 composite
  'Paid Service Unit Count', 'Bundled or Unbundled Line Number',
] // SVD03-8, SVD04 Not Used; Transmittal 907, PDF p.235 (footer 103)
export const SVD03_1_NAME: Record<string, string> = { HC: 'Healthcare Common Procedure Coding System' } // Medicare requires SVD03-1 = HC; claims with ER, IV or WK are rejected; CEDI 837P CG, PDF p.49 (printed p.41; Table 46)

// ===========================================================================
// 2440: Supporting Documentation (Form/Clinical Questionnaire)
// ===========================================================================

export const LQ_NAME = 'Form Identification Code' // Transmittal 907, PDF p.238 (footer 106)
export const LQ_EL_NAME: [string, string] = ['Code List Qualifier Code', 'Form Identifier'] // Transmittal 907, PDF p.238 (footer 106)

export const FRM_NAME = 'Supporting Documentation' // Transmittal 907, PDF p.238 (footer 106); repeats up to 99 times
export const FRM_EL_NAME: string[] = ['Question Number/Letter', 'Question Response', 'Question Response', 'Question Response', 'Question Response'] // Transmittal 907, PDF p.238 (footer 106)

// ===========================================================================
// Qualifier code values sourced from Transmittal 579
// ===========================================================================
// REF01, DTP01, AMT01, QTY01, CRC01, and HI0n-1 above each keep only their
// box title or element-position name, because Transmittal 907 does not
// restate these qualifiers' own code characters for most boxes (see the file
// header). Where a code IS independently confirmed - either in the CEDI
// guide's worked example/restriction tables (already cited inline above) or
// in Transmittal 579's own Values column - it gets an entry below, keyed by
// the qualifier code itself. A code used identically at more than one loop
// (e.g. REF01 0B/1G/G2/LU, the provider secondary-identification
// qualifiers) keeps one entry with every citing page listed; this file does
// not invent a different name for the same code at a different loop.

// REF01, Reference Identification Qualifier. Where Transmittal 579's Values
// column lists a code but neither it nor the CEDI guide states the code's
// own English meaning (it only shows which box/loop carries it), the entry
// says so rather than guessing a name - the same choice this file makes
// everywhere else (see the file header) and the same choice
// backend/pkg/x12837p's ReferenceIdentificationQualifierREF01 enum makes for
// these same codes.
export const REF01_NAME: Record<string, string> = {
  EI: 'Billing Provider Tax Identification', // CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "REF*EI*123456789"); Transmittal 579, PDF p.12 (2010AA, Values column "EI, SY"), 16 (2010AC), 21 (2010BB), 71 (2330B)
  SY: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.12 (2010AA, Values column "EI, SY"), 19 (2010BA), 69 (2330A)
  '0B': '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.12 (2010AA, Values column "0B, 1G") and the 2310A-D, 2330C-F, 2420A/D/E/F secondary identification boxes
  '1G': '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.12 (2010AA, Values column "0B, 1G") and the 2310A/B/D, 2330C/D/F, 2420A/D/E/F secondary identification boxes
  '2U': '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.15 (2010AC, Values column "2U, FY, NF"), 21 (2010BB), 71 (2330B)
  FY: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.15 (2010AC, Values column "2U, FY, NF"), 21 (2010BB), 71 (2330B)
  NF: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.15 (2010AC, Values column "2U, FY, NF"), 21 (2010BB), 71 (2330B)
  Y4: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.19 (2010BA, Property and Casualty Claim Number box) and p.25 (2010CA)
  G2: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.22 (2010BB, Values column "G2, LU") and the 2310A-D, 2330C-G, 2420A/C-F secondary identification boxes
  LU: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.22 (2010BB, Values column "G2, LU") and the 2310B-D, 2330D-G, 2420A/C/D secondary identification boxes
  '4N': 'Service Authorization Exception Code', // Transmittal 579, PDF p.32 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.serviceAuthException
  F5: 'Mandatory Medicare (Section 4081) Crossover Indicator', // Transmittal 579, PDF p.32 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.medicareCrossoverIndicator; not used by Medicare
  EW: 'Mammography Certification Number', // Transmittal 579, PDF p.33 (2300), 90 (2400, Values column) - each box's only REF01 value; matches CLAIM_REF_NAME.mammographyCertification / SERVICE_REF_NAME.mammographyCertification
  '9F': 'Referral Number', // Transmittal 579, PDF p.33 (2300), 72 (2330B), 92 (2400, Values column) - each box's only REF01 value; matches CLAIM_REF_NAME.referral / OTHER_PAYER_REF_NAME.referral / SERVICE_REF_NAME.referral
  G1: 'Prior Authorization', // Transmittal 579, PDF p.34 (2300), 71 (2330B), 89 (2400, Values column) - each box's only REF01 value; matches CLAIM_REF_NAME.priorAuthorization / OTHER_PAYER_REF_NAME.priorAuthorization / SERVICE_REF_NAME.priorAuthorization
  F8: 'Payer Claim Control Number', // Transmittal 579, PDF p.34 (2300), 73 (2330B, Values column) - each box's only REF01 value; matches CLAIM_REF_NAME.payerClaimControlNumber; not used by Medicare
  X4: 'Clinical Laboratory Improvement Amendment (CLIA) Number', // Transmittal 579, PDF p.35 (2300), 90 (2400, Values column) - each box's only REF01 value; matches CLAIM_REF_NAME.clia / SERVICE_REF_NAME.clia
  '9A': 'Repriced Claim Number', // Transmittal 579, PDF p.35 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.repricedClaim
  '9C': 'Adjusted Repriced Claim Number', // Transmittal 579, PDF p.36 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.adjustedRepricedClaim
  LX: 'Investigational Device Exemption Number', // Transmittal 579, PDF p.36 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.investigationalDeviceExemption
  D9: 'Claim Identifier for Transmission Intermediaries', // Transmittal 579, PDF p.37 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.claimIdentifierForTransmissionIntermediaries
  EA: 'Medical Record Number', // Transmittal 579, PDF p.37 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.medicalRecordNumber
  P4: 'Demonstration Project Identifier', // Transmittal 579, PDF p.38 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.demonstrationProject
  '1J': 'Care Plan Oversight', // Transmittal 579, PDF p.38 (2300, Values column) - the box's only REF01 value; matches CLAIM_REF_NAME.carePlanOversight
  T4: 'Other Payer Claim Adjustment Indicator', // Transmittal 579, PDF p.72 (2330B, Values column) - the box's only REF01 value; matches OTHER_PAYER_REF_NAME.claimAdjustmentIndicator
  '9B': 'Repriced Line Item Reference Number', // Transmittal 579, PDF p.88 (2400, Values column) - the box's only REF01 value; matches SERVICE_REF_NAME.repricedLineItem
  '9D': 'Adjusted Repriced Line Item Reference Number', // Transmittal 579, PDF p.88 (2400, Values column) - the box's only REF01 value; matches SERVICE_REF_NAME.adjustedRepricedLineItem
  F4: 'Referring Clinical Laboratory Improvement Amendment (CLIA) Facility Identification', // Transmittal 579, PDF p.91 (2400, Values column) - the box's only REF01 value; matches SERVICE_REF_NAME.referringClia
  BT: 'Immunization Batch Number', // Transmittal 579, PDF p.91 (2400, Values column) - the box's only REF01 value; matches SERVICE_REF_NAME.immunizationBatch
  VY: '(code name not given in Transmittal 579 or the CEDI guide)', // Transmittal 579, PDF p.97 (2410 Prescription or Compound Drug Association Number box, Values column "VY, XZ")
  XZ: 'Pharmacy Prescription Number', // Transmittal 579, PDF p.97 (2410 Prescription or Compound Drug Association Number box, Values column "VY, XZ"); DRUG_ASSOCIATION_REF01_NAME above confirms this same code and name independently from the CEDI guide
  PXC: 'Health Care Provider Taxonomy Code', // PRV02's only value at 2000A, 2310B, 2420A; Transmittal 579, PDF p.9, 58, 98 (Values column, code only, no name given); name confirmed from frontend/src/data/x12Names837i.ts PRV01_EL_NAME's own citation, Transmittal 907 (837I edition), PDF p.10 - a generic X12 code, not transaction-specific
}

// DTP01, Date/Time Qualifier. Every name below matches an existing box title
// in CLAIM_DTP_NAME, SERVICE_DTP_NAME, OTHER_PAYER_DTP_NAME, or
// LINE_ADJUDICATION_DTP_NAME above (all sourced from Transmittal 907
// directly); where Transmittal 579 lists more than one code for the same
// box without distinguishing them, every such code keeps that one shared
// box title rather than an invented distinction - the same treatment
// NM101_CONFIRMED_NAME gives DN/P3 above.
export const DTP01_NAME: Record<string, string> = {
  '435': 'Date - Admission', // CEDI 837P CG, PDF p.37 (printed p.29, below Table 18): required when 2400 SV105 is 21, 51, or 61; Transmittal 579, PDF p.30 (2300, Values column); matches CLAIM_DTP_NAME.admission
  '431': 'Date - Onset of Current Illness/Symptom', // Transmittal 579, PDF p.27 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.onsetOfIllness
  '454': 'Date - Initial Treatment', // Transmittal 579, PDF p.28 (2300), 86 (2400, Values column) - each box's only DTP01 value; matches CLAIM_DTP_NAME.initialTreatment / SERVICE_DTP_NAME.initialTreatment
  '304': 'Date - Last Seen Date', // Transmittal 579, PDF p.28 (2300 Last Seen Date), 85 (2400 Date Last Seen Date, Values column) - each box's only DTP01 value; matches CLAIM_DTP_NAME.lastSeen / SERVICE_DTP_NAME.lastTreatmentOrTherapy
  '453': 'Date - Acute Manifestation', // Transmittal 579, PDF p.28 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.acuteManifestation
  '439': 'Date - Accident', // Transmittal 579, PDF p.28 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.accident
  '484': 'Date - Last Menstrual Period', // Transmittal 579, PDF p.28 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.lastMenstrualPeriod
  '455': 'Date - Last X-Ray Date', // Transmittal 579, PDF p.29 (2300), 86 (2400, Values column) - each box's only DTP01 value; matches CLAIM_DTP_NAME.lastXRay / SERVICE_DTP_NAME.lastXRay
  '471': 'Date - Hearing and Vision Prescription Date', // Transmittal 579, PDF p.29 (2300), 84 (2400 Prescription Date, Values column) - each box's only DTP01 value; matches CLAIM_DTP_NAME.hearingVisionPrescription / SERVICE_DTP_NAME.prescription
  '314': 'Date - Disability Dates', // Transmittal 579, PDF p.29 (2300, Values column "314, 360, 361") - one of three DTP01 values for that one box; matches CLAIM_DTP_NAME.disability
  '360': 'Date - Disability Dates', // Transmittal 579, PDF p.29 (2300, Values column "314, 360, 361") - one of three DTP01 values for that one box; matches CLAIM_DTP_NAME.disability
  '361': 'Date - Disability Dates', // Transmittal 579, PDF p.29 (2300, Values column "314, 360, 361") - one of three DTP01 values for that one box; matches CLAIM_DTP_NAME.disability
  '297': 'Date - Last Worked', // Transmittal 579, PDF p.29 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.lastWorked
  '296': 'Date - Authorized Return to Work', // Transmittal 579, PDF p.29 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.authorizedReturnToWork
  '096': 'Date - Discharge', // Transmittal 579, PDF p.30 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.discharge
  '090': 'Date - Assumed and Relinquished Care Dates', // Transmittal 579, PDF p.30 (2300, Values column "090, 091") - one of two DTP01 values for that one box, repeating once per date; matches CLAIM_DTP_NAME.assumedOrRelinquishedCare
  '091': 'Date - Assumed and Relinquished Care Dates', // Transmittal 579, PDF p.30 (2300, Values column "090, 091") - one of two DTP01 values for that one box, repeating once per date; matches CLAIM_DTP_NAME.assumedOrRelinquishedCare
  '444': 'Date - Property and Casualty Date of First Contact', // Transmittal 579, PDF p.30 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.propertyCasualtyFirstContact
  '050': 'Date - Repricer Received Date', // Transmittal 579, PDF p.30 (2300, Values column) - the box's only DTP01 value; matches CLAIM_DTP_NAME.repricerReceived
  '573': 'Date - Claim Check or Remittance Date', // Transmittal 579, PDF p.70 (2330B), 110 (2430, Values column) - each box's only DTP01 value; matches OTHER_PAYER_DTP_NAME / LINE_ADJUDICATION_DTP_NAME
  '472': 'Date - Service Date', // Transmittal 579, PDF p.84 (2400, Values column) - the box's only DTP01 value; matches SERVICE_DTP_NAME.serviceDate
  '607': 'Date - Certification Revision/Recertification Date', // Transmittal 579, PDF p.84 (2400, Values column) - the box's only DTP01 value; matches SERVICE_DTP_NAME.certificationRevision
  '463': 'Date - Begin Therapy Date', // Transmittal 579, PDF p.85 (2400, Values column) - the box's only DTP01 value; matches SERVICE_DTP_NAME.beginTherapy
  '461': 'Date - Last Certification Date', // Transmittal 579, PDF p.85 (2400, Values column) - the box's only DTP01 value; matches SERVICE_DTP_NAME.lastCertification
  '738': 'Date - Test Date', // Transmittal 579, PDF p.85 (2400, Values column "738, 739") - one of two DTP01 values for that one box, repeating up to 2 times; matches SERVICE_DTP_NAME.testPerformed
  '739': 'Date - Test Date', // Transmittal 579, PDF p.85 (2400, Values column "738, 739") - one of two DTP01 values for that one box, repeating up to 2 times; matches SERVICE_DTP_NAME.testPerformed
  '011': 'Date - Shipped Date', // Transmittal 579, PDF p.85 (2400, Values column) - the box's only DTP01 value; matches SERVICE_DTP_NAME.shipped
}

// AMT01, Amount Qualifier Code.
export const AMT01_NAME: Record<string, string> = {
  D: 'Payer Amount Paid', // Medicare requires one 2320 AMT occurrence with AMT01 = D whenever 2000B SBR01 = S; CEDI 837P CG, PDF p.40 (printed p.32; Table 29); Transmittal 579, PDF p.66 (2320 COB Payer Paid Amount box); matches AMT01_COB_NAME.D
  F5: 'Patient Amount Paid', // Transmittal 579, PDF p.32 (2300, Values column) - the box's only AMT01 value; matches CLAIM_AMT_NAME
  A8: 'COB Total Non-Covered Amount', // Transmittal 579, PDF p.66 (2320, Values column) - the box's only AMT01 value; matches COB_NON_COVERED_AMT_NAME
  EAF: 'Remaining Patient Liability', // Transmittal 579, PDF p.66 (2320), 110 (2430, Values column) - each box's only AMT01 value; matches OTHER_SUBSCRIBER_REMAINING_LIABILITY_AMT_NAME / LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME
  T: 'Sales Tax Amount', // Transmittal 579, PDF p.92 (2400, Values column) - the box's only AMT01 value; matches SALES_TAX_AMT_NAME
  F4: 'Postage Claimed Amount', // Transmittal 579, PDF p.92 (2400, Values column) - the box's only AMT01 value; matches POSTAGE_CLAIMED_AMT_NAME
}

// QTY01, Quantity Qualifier.
export const QTY01_NAME: Record<string, string> = {
  PT: 'Ambulance Patient Count', // Transmittal 579, PDF p.86 (2400, Values column) - the box's only QTY01 value; matches AMBULANCE_PATIENT_COUNT_QTY_NAME
  FL: 'Obstetric Anesthesia Additional Units', // Transmittal 579, PDF p.86 (2400, Values column) - the box's only QTY01 value; matches OBSTETRIC_ANESTHESIA_QTY_NAME
}

// CRC01, Code Category. Every name below matches an existing box title in
// CLAIM_CRC_NAME or SERVICE_CRC_NAME above (sourced from Transmittal 907
// directly); where Transmittal 579 lists more than one code for the same
// box without distinguishing them (E1/E2/E3), every such code keeps that
// one shared box title rather than an invented distinction.
export const CRC01_NAME: Record<string, string> = {
  '07': 'Ambulance Certification', // Transmittal 579, PDF p.40 (2300), 83 (2400, Values column) - each box's only CRC01 value; repeats up to 3 times; matches CLAIM_CRC_NAME.ambulanceCertification / SERVICE_CRC_NAME.ambulanceCertification
  E1: 'Patient Condition Information: Vision', // Transmittal 579, PDF p.40 (2300, Values column "E1, E2, E3") - one of three CRC01 values for that one box, repeating up to 3 times; matches CLAIM_CRC_NAME.visionCondition
  E2: 'Patient Condition Information: Vision', // Transmittal 579, PDF p.40 (2300, Values column "E1, E2, E3") - one of three CRC01 values for that one box, repeating up to 3 times; matches CLAIM_CRC_NAME.visionCondition
  E3: 'Patient Condition Information: Vision', // Transmittal 579, PDF p.40 (2300, Values column "E1, E2, E3") - one of three CRC01 values for that one box, repeating up to 3 times; matches CLAIM_CRC_NAME.visionCondition
  '75': 'Homebound Indicator', // Transmittal 579, PDF p.41 (2300, Values column) - the box's only CRC01 value; matches CLAIM_CRC_NAME.homebound
  ZZ: 'EPSDT Referral', // Transmittal 579, PDF p.41 (2300, Values column) - the box's only CRC01 value; matches CLAIM_CRC_NAME.epsdtReferral
  '70': 'Hospice Employee Indicator', // Transmittal 579, PDF p.83 (2400, Values column) - the box's only CRC01 value; matches SERVICE_CRC_NAME.hospiceEmployeeIndicator
  '09': 'Condition Indicator Durable Medical Equipment', // Transmittal 579, PDF p.84 (2400, Values column) - the box's only CRC01 value; matches SERVICE_CRC_NAME.dmeConditionIndicator
}

// HI0n-1, Diagnosis/Procedure Type Code Qualifier (the 2300 Health Care
// Diagnosis Code, Anesthesia Related Procedure, and Condition Information
// composites). Transmittal 907 never restates these qualifier letters;
// Transmittal 579's Values column does, at the first and (where it differs)
// a later repetition of the same composite. No English name is given in
// either source for ABF/BF/BP/BO/BG beyond the position itself, so this
// file names them only by position and code, the same choice
// backend/pkg/x12837p's DiagnosisTypeCodeHI0x_1 makes.
export const HI0x_1_NAME: Record<string, string> = {
  ABK: 'HI01-1 of the Health Care Diagnosis Code box', // Transmittal 579, PDF p.42 (Values column "ABK, BK")
  BK: 'HI01-1 of the Health Care Diagnosis Code box', // Transmittal 579, PDF p.42 (Values column "ABK, BK"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HI*BK>71598")
  ABF: 'HI02-1 through HI12-1 of the Health Care Diagnosis Code box', // Transmittal 579, PDF p.42-46 (Values column "ABF, BF")
  BF: 'HI02-1 through HI12-1 of the Health Care Diagnosis Code box', // Transmittal 579, PDF p.42-46 (Values column "ABF, BF")
  BP: 'HI01-1 of the Anesthesia Related Procedure box', // Transmittal 579, PDF p.46 (Values column)
  BO: 'HI02-1 of the Anesthesia Related Procedure box', // Transmittal 579, PDF p.47 (Values column)
  BG: 'HI01-1 through HI12-1 of the Condition Information box', // Transmittal 579, PDF p.51-55 (Values column)
}

// ===========================================================================
// CARC/RARC reuse note
// ===========================================================================
// 2320 CAS (reporting another payer's adjustments) and 2430 CAS (reporting
// another payer's line adjustments) use the same CARC adjustment-reason code
// set the 835 side of this project already embeds as CARC in
// x12Names835.ts, sourced from backend/pkg/carc.json (X12's published CARC
// list, x12.org/codes/claim-adjustment-reason-codes). That table is not repeated
// here; import CARC from x12Names835.ts for CAS reason-code text on the
// 837P side too.
