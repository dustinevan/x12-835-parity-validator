// Names and coded-value text for the "835 viewer" docs page
// (frontend/src/docs/pages/Viewer835.tsx), its lookup panel and the 835
// Reference (through lookupRows835.ts), and the validator app's hover
// explanations (codeHints.ts), built only from the table in
// docs/design/835-viewer.md ("Codes shown, by their public name").
// Each entry's comment repeats that table's Source column. A code the
// spec's "Open" list leaves unresolved has no entry here and the pages
// render it as its raw code. Every name here comes from that table, not
// from any other naming in this repository.

// Bare segment names (no qualifier involved). Citation form throughout this
// file: "CMS 835 Flat File, PDF p.N" (N counts PDF viewer pages of the
// 15-page docs/reference/file-formats/cms-835-flatfile.pdf, from 1) and
// "CEDI 835 CG, PDF p.N (printed p.M)" (N is the PDF viewer page of the
// 33-page docs/reference/file-formats/cedi-835-companion-guide.pdf; M is
// that same page's own printed footer number, always N-6 in this guide —
// the two differ because its first 6 PDF pages are a cover and front
// matter the guide does not number). A bare "p.N" with neither label is
// always a continuation of the citation just before it on the same line.
export const SEGMENT_NAME: Partial<Record<string, string>> = {
  ISA: 'Interchange Control Header', // CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.20-21 (printed p.14-15; Table 3)
  GS: 'Functional Group Header', // CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.21 (printed p.15; Table 4)
  ST: 'Transaction Set Header', // CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.23 (printed p.17; Table 5)
  BPR: 'Financial Information', // CMS 835 Flat File, PDF p.1; CEDI 835 CG, PDF p.24 (printed p.18; Table 6)
  CUR: 'Foreign Currency Information', // CMS 835 Flat File, PDF p.2 (not used by Medicare, per the flat file's own comment); docs/design/835-examples.md H-CUR
  LX: 'Header Number', // CMS 835 Flat File, PDF p.4
  TS3: 'Provider Summary Information', // CMS 835 Flat File, PDF p.4-5; docs/design/835-examples.md PS-TS3
  TS2: 'Provider Supplemental Summary Information', // CMS 835 Flat File, PDF p.5; docs/design/835-examples.md PS-TS2
  CLP: 'Claim Payment Information', // CEDI 835 CG, PDF p.25 (printed p.19; Table 9); CMS 835 Flat File, PDF p.5 names it "Claim Level Data"; the two disagree, and this name is the CEDI/TR3-matching one
  MIA: 'Inpatient Adjudication Information', // CMS 835 Flat File, PDF p.8-9; docs/design/835-examples.md C-MIA
  MOA: 'Outpatient Adjudication Information', // CMS 835 Flat File, PDF p.9
  SVC: 'Service Payment Information', // CMS 835 Flat File, PDF p.11
  PLB: 'Provider Level Adjustment', // CMS 835 Flat File, PDF p.13; docs/design/835-examples.md "Provider Level Adjustment (PLB)"
  SE: 'Transaction Set Trailer', // CMS 835 Flat File, PDF p.14; CEDI 835 CG, PDF p.32 (printed p.26; Acronym List, row "ST/SE")
  GE: 'Functional Group Trailer', // CMS 835 Flat File, PDF p.14; CEDI 835 CG, PDF p.31 (printed p.25; Acronym List, row "GS/GE")
  IEA: 'Interchange Control Trailer', // CMS 835 Flat File, PDF p.14; CEDI 835 CG, PDF p.32 (printed p.26; Acronym List, row "ISA/IEA")
  RDM: 'Remittance Delivery Method', // CMS 835 Flat File, PDF p.4; docs/design/835-examples.md PE-RDM
}

// GS01 / ST01: the transaction-type qualifier, same combined name either way.
export const TRANSACTION_TYPE_NAME: Record<string, string> = {
  HP: 'Health Care Claim Payment/Advice', // CEDI 835 CG, PDF p.7 (printed p.1)
  '835': 'Health Care Claim Payment/Advice', // CEDI 835 CG, PDF p.7 (printed p.1)
}

// TRN01 trace-type qualifier.
export const TRN_NAME: Record<string, string> = {
  '1': 'Reassociation Trace Number', // CMS 835 Flat File, PDF p.2
}

// DTM01 date-qualifier, wherever it appears (header, claim, or service level).
export const DTM_NAME: Record<string, string> = {
  '405': 'Production Date', // CMS 835 Flat File, PDF p.2
  '050': 'Claim Received Date', // CMS 835 Flat File, PDF p.11
  '232': 'Statement From-To Date', // CMS 835 Flat File, PDF p.10; docs/design/835-examples.md C-DTM-STMT
  '233': 'Statement From-To Date', // CMS 835 Flat File, PDF p.10; docs/design/835-examples.md C-DTM-STMT
  '036': 'Coverage Expiration Date', // CMS 835 Flat File, PDF p.10; docs/design/835-examples.md C-DTM-COVEXP
  '472': 'Service Date', // CMS 835 Flat File, PDF p.12
}

// Header-level REF (directly under ST, before any loop).
export const HEADER_REF_NAME: Record<string, string> = {
  EV: 'Receiver Identification', // CMS 835 Flat File, PDF p.2; docs/design/835-examples.md H-REF-EV
  F2: 'Version Identification', // CMS 835 Flat File, PDF p.2; docs/design/835-examples.md H-REF-F2
}

// Claim-level REF (loop 2100), qualifiers EA/6P/28/F8/CE share one combined name.
export const CLAIM_REF_NAME: Record<string, string> = {
  EA: 'Other Claim Related Identification', // CMS 835 Flat File, PDF p.10; CEDI 835 CG, PDF p.26 (printed p.20; Table 14); docs/design/835-examples.md C-REF-OTHERCLAIM
  F8: 'Other Claim Related Identification', // CMS 835 Flat File, PDF p.10; CEDI 835 CG, PDF p.26 (printed p.20; Table 14); docs/design/835-examples.md C-REF-OTHERCLAIM (F8 = Original Reference Number, used to link a correction to the claim it corrects)
}

// Payer-loop PER (1000A), by qualifier.
export const PAYER_PER_NAME: Record<string, string> = {
  CX: 'Payer Business Contact Information', // CMS 835 Flat File, PDF p.3; docs/design/835-examples.md PR-PER-CX
  BL: 'Payer Technical Contact Information', // CMS 835 Flat File, PDF p.3; docs/design/835-examples.md PR-PER-BL
  '1C': 'Payer Web Site', // CMS 835 Flat File, PDF p.3; docs/design/835-examples.md PR-PER-1C
}

// Claim-level PER (loop 2100), by qualifier.
export const CLAIM_PER_NAME: Record<string, string> = {
  CX: 'Claim Contact Information', // CMS 835 Flat File, PDF p.11; docs/design/835-examples.md C-PER
}

// LQ01 qualifier (service line health care remark codes segment).
export const LQ_NAME: Record<string, string> = {
  HE: 'Health Care Remarks Codes', // CMS 835 Flat File, PDF p.13; CEDI 835 CG, PDF p.28 (printed p.22; Table 24); docs/design/835-examples.md S-LQ
}

// N101 payer/payee role.
export const N1_NAME: Record<string, string> = {
  PR: 'Payer Identification', // CMS 835 Flat File, PDF p.2
  PE: 'Payee Identification', // CMS 835 Flat File, PDF p.3
}

// N3/N4/REF inside the payer loop (1000A) carry different names than the
// same segment IDs inside the payee loop (1000B).
export const PAYER_LOOP_NAME: Partial<Record<string, string>> = {
  N3: 'Payer Address', // CMS 835 Flat File, PDF p.2
  N4: 'Payer City, State, Zip', // CMS 835 Flat File, PDF p.2
  REF: 'Additional Payer Identification', // CMS 835 Flat File, PDF p.3
}
export const PAYEE_LOOP_NAME: Partial<Record<string, string>> = {
  N3: 'Payee Address', // CMS 835 Flat File, PDF p.3
  N4: 'Payee City, State, Zip', // CMS 835 Flat File, PDF p.3-4 (Payee City Name on p.3; Payee State Code and Postal Zone/ZIP on p.4)
  REF: 'Payee Additional Identification', // CMS 835 Flat File, PDF p.4; docs/design/835-examples.md PE-REF (qualifier TJ)
}

// ID-type qualifier (N103 / NM108) shared across loops.
export const ID_QUALIFIER_NAME: Record<string, string> = {
  XX: 'National Provider Identifier', // CMS 835 Flat File, PDF p.7 (NM109 content shown as "NPI"); CEDI 835 CG, PDF p.32 (printed p.26; Acronym List)
}

// NM101 role qualifier in the claim loop.
export const NM1_NAME: Record<string, string> = {
  QC: 'Patient Name', // CMS 835 Flat File, PDF p.6
  '74': 'Corrected Patient/Insured Name', // CMS 835 Flat File, PDF p.7; docs/design/835-examples.md C-NM1-74
  '82': 'Service Provider Name', // CMS 835 Flat File, PDF p.7
  TT: 'Crossover Carrier Name', // CMS 835 Flat File, PDF p.7; docs/design/835-examples.md C-NM1-TT
  PR: 'Corrected Priority Payer Name', // CMS 835 Flat File, PDF p.7-8; docs/design/835-examples.md C-NM1-PR
}

// BPR03 credit/debit flag.
export const BPR03_NAME: Record<string, string> = {
  C: 'Credit', // CMS 835 Flat File, PDF p.1 ("Credit or Debit Flag Code"); CEDI 835 CG, PDF p.24 (printed p.18; only C applies to Medicare)
}

// SVC01 composite's first component (code-system qualifier).
export const SVC_QUALIFIER_NAME: Record<string, string> = {
  HC: 'Healthcare Common Procedure Coding System', // CEDI 835 CG, PDF p.31 (printed p.25; Acronym List, "HCPCS"); CMS 835 Flat File, PDF p.11 lists HC as a Product/Service ID Qualifier value
}

// CAS: the same segment, named differently by the loop it sits in.
export const CLAIM_CAS_NAME = 'Claim Adjustment' // CMS 835 Flat File, PDF p.6; CEDI 835 CG, PDF p.25 (printed p.19; Table 10); docs/design/835-examples.md C-CAS
export const SERVICE_CAS_NAME = 'Service Adjustment' // CMS 835 Flat File, PDF p.12
export const CAS_GROUP_NAME: Record<string, string> = {
  CO: 'Contractual Obligation', // x12.org/codes/claim-adjustment-group-codes
  OA: 'Other Adjustment', // x12.org/codes/claim-adjustment-group-codes
  PR: 'Patient Responsibility', // x12.org/codes/claim-adjustment-group-codes
}

// REF at the service level (2110), by qualifier.
export const SERVICE_REF_NAME: Record<string, string> = {
  '6R': 'Line Item Control Number', // CMS 835 Flat File, PDF p.13
  LU: 'Service Identification', // CMS 835 Flat File, PDF p.12; CEDI 835 CG, PDF p.27 (printed p.21; Table 20); docs/design/835-examples.md S-REF-SVCID
  HPI: 'Rendering Provider Information', // CMS 835 Flat File, PDF p.13; CEDI 835 CG, PDF p.28 (printed p.22; Table 21); docs/design/835-examples.md S-REF-RENDERING
  '0K': 'Health Care Policy Identification', // CMS 835 Flat File, PDF p.13; CEDI 835 CG, PDF p.28 (printed p.22; Table 22); docs/design/835-examples.md S-REF-HCPOLICY
}

// AMT: the same segment, named differently by the loop it sits in. The name
// applies regardless of AMT01 (qualifiers B6, I, ZK are all in the spec's
// Open list).
export const CLAIM_AMT_NAME = 'Claim Supplemental Information' // CMS 835 Flat File, PDF p.11; docs/design/835-examples.md C-AMT
export const SERVICE_AMT_NAME = 'Service Supplemental Amount' // CMS 835 Flat File, PDF p.13

// QTY at the claim level (2100): CMS 835 Flat File, PDF p.11. The name applies
// regardless of QTY01 (qualifier CA is in the spec's Open list — no
// individual QTY qualifier is named in either cited document).
export const CLAIM_QTY_NAME = 'Claim Supplemental Information Quantity' // CMS 835 Flat File, PDF p.11; docs/design/835-examples.md C-QTY

// PLB03-1/05-1/07-1/09-1/11-1/13-1: the Provider Adjustment Reason Code,
// X12's public list (x12.org/codes/provider-adjustment-reason-codes,
// element 426). Seven codes Medicare's MACs use are not on the list's
// current public edition and are left unresolved (shown as their code):
// C5, CV, DM, GO, LS, OA (distinct from the CAS group code OA above), RE.
export const PLB_REASON_NAME: Record<string, string> = {
  '50': 'The amount of the late charge, late claim filing penalty, or Medicare late cost report penalty',
  '51': 'Late filing interest penalty assessment amount',
  '72': 'Provider refund amount',
  '90': 'Early payment allowance amount',
  AP: 'Advanced or accelerated payment amount',
  B2: 'Rebate amount',
  B3: 'Recovery amount (from the provider, for an overpayment based on payments from other payers)',
  BD: 'Bad debt amount',
  BN: 'Bonus amount',
  CS: 'Adjustment amount, detailed information provided separately',
  E3: 'Withholding amount',
  FB: 'Non-claim related balance forward amount',
  HM: 'Hemophilia clotting factor add-on payment amount',
  IP: 'Incentive payment amount',
  IR: 'Internal Revenue Service 1099 withholding amount',
  IS: 'Lump sum based on an interim rate',
  J1: 'Amount not reimbursed based on a demonstration program or other limitation',
  L3: 'Penalty amount',
  L6: 'Interest amount',
  LE: 'Internal Revenue Service non-1099 withholding amount',
  OB: 'Affiliated provider(s) offset amount',
  PI: 'Periodic Interim Payment (PIP) lump sum amount',
  PL: 'Final payment or settlement amount',
  RA: 'Retroactive adjustment amount',
  SL: 'Student loan garnishment amount',
  TL: 'Third party liability determination amount',
  WO: 'Overpayment recovery amount',
  WU: 'Non-Internal Revenue Service withholding amount related to a federal payment levy program',
}

// Element-position names: the name a cell shows beneath its own code (ISA01,
// CLP04, ...) — a different thing from the tables above, which name a
// CODE's resolved value (e.g. N1_NAME.PR names the code "PR", not the N101
// element itself). Every array below is 0-indexed from element 01 (index 0
// = element 01). Source throughout: CMS 835 Flat File, PDF p.N — PDF viewer
// pages of the 15-page flat file, counted from 1. This is NOT the flat
// file's own "Page #" column: that column
// is a reference into the paid X12 TR3 implementation guide, a different,
// much longer document, and its numbers (visible in the flat file itself)
// run past 200 — nowhere close to a page of this 15-page PDF. Where that TR3
// page is worth keeping as a cross-reference it is noted in parentheses as
// "TR3 p.N per the flat file's Page # column"; a later unlabeled "p.N" on
// the same line continues the citation just before it.

export const ISA_EL_NAME: string[] = [
  'Authorization Information Qualifier', 'Authorization Information',
  'Security Information Qualifier', 'Security Information',
  'Interchange ID Qualifier', 'Interchange Sender ID',
  'Interchange ID Qualifier', 'Interchange Receiver ID',
  'Interchange Date', 'Interchange Time', 'Repetition Separator',
  'Interchange Control Version Number', 'Interchange Control Number',
  'Acknowledgement Requested', 'Usage Indicator', 'Component Element Separator',
] // CMS 835 Flat File, PDF p.1

export const GS_EL_NAME: string[] = [
  'Functional Identifier Code', "Application Sender's Code", "Application Receiver's Code", // CMS 835 Flat File, PDF p.1
  'Date', 'Time', 'Group Control Number', 'Responsible Agency Code', 'Version/Release/Industry ID code', // p.2
]

export const ST_EL_NAME: string[] = ['Transaction Set Identifier Code', 'Transaction Set Control Number'] // CMS 835 Flat File, PDF p.1

export const BPR_EL_NAME: string[] = [
  'Transaction Handling Code', 'Total Actual Provider Payment Amt', 'Credit or Debit Flag Code',
  'Payment Method Code', 'Payment Format Code',
  'DFI ID # Qualifier', 'Sender DFI Identifier', 'Acct # Qualifier', 'Sender Bank Acct #',
  'Payer Identifier', 'Originating Co Supplemental Code', // BPR01-11: CMS 835 Flat File, PDF p.1
  'DFI ID # Qualifier', 'Receiver or Provider Bank ID #', 'Acct # Qualifier', 'Receiver or Provider Acct #',
  'Check Issue or EFT Effective Date', // BPR12-16: p.2
]

export const TRN_EL_NAME: string[] = ['Trace Type Code', 'Check or EFT Trace #', 'Payer Identifier', 'Originating Co Supplemental Code'] // CMS 835 Flat File, PDF p.2

// DTM01's own position name is the same everywhere; DTM02's name is the TR3's
// per-usage date name (also used, in some cases, as the qualifier CODE's own
// name via DTM_NAME above — the two tables serve different roles, see the
// comment at the top of this section).
export const DTM01_EL_NAME = 'Date/Time Qualifier' // CMS 835 Flat File, PDF p.2, 10-12 (every instance)
export const DTM02_EL_NAME_BY_QUALIFIER: Record<string, string> = {
  '405': 'Production Date', // CMS 835 Flat File, PDF p.2
  '232': 'Claim Date', // CMS 835 Flat File, PDF p.10
  '233': 'Claim Date', // CMS 835 Flat File, PDF p.10
  '036': 'Coverage Expiration Date', // CMS 835 Flat File, PDF p.10
  '050': 'Claim Received Date', // CMS 835 Flat File, PDF p.11
  '472': 'Service Date', // CMS 835 Flat File, PDF p.12
}

export const N1_PR_EL_NAME: string[] = ['Entity Identifier Code', 'Payer Name', 'ID Code Qualifier', 'Payer Identifier'] // CMS 835 Flat File, PDF p.2
export const N1_PE_EL_NAME: string[] = ['Entity Identifier Code', 'Payee Name', 'Identification Code Qualifier', 'Payee ID Code'] // CMS 835 Flat File, PDF p.3

export const N3_PR_EL_NAME: string[] = ['Payer Address Line', 'Payer Address Line'] // CMS 835 Flat File, PDF p.2
export const N3_PE_EL_NAME: string[] = ['Payee Address Line', 'Payee Address Line'] // CMS 835 Flat File, PDF p.3
export const N4_PR_EL_NAME: string[] = ['Payer City Name', 'Payer State Code', 'Payer Postal Zone or ZIP Code'] // CMS 835 Flat File, PDF p.2
export const N4_PE_EL_NAME: string[] = ['Payee City Name', 'Payee State Code', 'Payee Postal Zone or ZIP Code'] // CMS 835 Flat File, PDF p.3-4 (Payee City Name on p.3; the other two on p.4)

export const REF_PAYER_EL_NAME: string[] = ['Reference Identification Qualifier', 'Additional Payer ID'] // CMS 835 Flat File, PDF p.3 (TR3 p.92 per the flat file's Page # column)
export const REF_PAYEE_EL_NAME: string[] = ['Reference Identification Qualifier', 'Additional Payee ID #'] // CMS 835 Flat File, PDF p.4 (TR3 p.107 per the flat file's Page # column)

// Header, claim, and line REF each carry a different element name per
// qualifier — the same tables HEADER_REF_NAME/CLAIM_REF_NAME/SERVICE_REF_NAME
// above name what the CODE means; these name the REF01/REF02 positions themselves.
export const REF_HEADER_EL_NAME: Record<string, [string, string]> = {
  EV: ['Receiver ID Qualifier', 'Receiver Identifier'], // CMS 835 Flat File, PDF p.2 (TR3 p.82 per the flat file's Page # column)
  F2: ['Receiver ID Qualifier', 'Version ID Code'], // CMS 835 Flat File, PDF p.2 (TR3 p.84 per the flat file's Page # column)
}
export const REF_CLAIM_EL_NAME: [string, string] = ['Reference ID Qualifier/Medical Record ID #', 'Other Claim Related ID/Medical Record #'] // CMS 835 Flat File, PDF p.10 (TR3 p.169 per the flat file's Page # column)
export const REF_LINE_EL_NAME: Record<string, [string, string]> = {
  '6R': ['Reference ID Qualifier', 'Line Item Control Number'], // CMS 835 Flat File, PDF p.13 (TR3 p.206 per the flat file's Page # column)
  LU: ['Reference ID Qualifier', 'Reference Identification'], // CMS 835 Flat File, PDF p.12 (TR3 p.204 per the flat file's Page # column)
  HPI: ['Reference ID Qualifier', 'Rendering Provider ID'], // CMS 835 Flat File, PDF p.13 (TR3 p.207 per the flat file's Page # column)
  '0K': ['Reference ID Qualifier', 'Reference Identification'], // CMS 835 Flat File, PDF p.13 (TR3 p.209 per the flat file's Page # column)
}

export const PER01_EL_NAME = 'Contact Function Code' // CMS 835 Flat File, PDF p.3, 11, consistent across every PER instance
// PER02 onward, by qualifier; CX and BL share one shape (Payer Business /
// Technical Contact Information use identical element names), 1C (Payer Web
// Site) is shorter.
export const PER_PAYER_EL_NAME: Record<string, string[]> = {
  CX: ['Payer Contact Name', 'Communication # Qualifier', 'Payer Contact Communication #', 'Communication Number Qualifier', 'Payer Contact Communication #', 'Communication Number Qualifier', 'Payer Contact Communication #'], // CMS 835 Flat File, PDF p.3 (TR3 p.94 per the flat file's Page # column)
  BL: ['Payer Contact Name', 'Communication # Qualifier', 'Payer Contact Communication #', 'Communication Number Qualifier', 'Payer Contact Communication #', 'Communication Number Qualifier', 'Payer Contact Communication #'], // CMS 835 Flat File, PDF p.3 (TR3 p.97 per the flat file's Page # column)
  '1C': ['Name', 'Communication # Qualifier', 'Payer Contact Communication #'], // CMS 835 Flat File, PDF p.3 (TR3 p.100 per the flat file's Page # column)
}
export const PER_CLAIM_EL_NAME: string[] = [
  'Claim Contact Name', 'Communication # Qualifier', 'Claim Contact Communication #',
  'Communication # Qualifier', 'Claim Contact Communication #', 'Communication # Qualifier', 'Communication # Extension',
] // CMS 835 Flat File, PDF p.11 (TR3 p.179 per the flat file's Page # column)

export const NM101_EL_NAME = 'Entity Identifier Code' // CMS 835 Flat File, PDF p.6-8, consistent across every NM1 instance
export const NM102_EL_NAME = 'Entity Type Qualifier'
// NM103 onward, by role (NM101's value).
export const NM1_QC_EL_NAME: string[] = ['Patient Last Name', 'Patient First Name', 'Patient Middle Name', 'Name Prefix', 'Patient Name Suffix', 'ID Code Qualifier', 'Patient Identifier'] // CMS 835 Flat File, PDF p.6-7
export const NM1_74_EL_NAME: string[] = ['Corrected Patient/Ins Last Name', 'Corrected Patient/Ins First Name', 'Corrected Patient/Ins Middle Name', 'Name Prefix', 'Corrected Patient Name Suffix', 'Identification Code Qualifier', 'Corrected Ins Identification Indicator'] // CMS 835 Flat File, PDF p.7
export const NM1_82_EL_NAME: string[] = ['Rendering Provider Last/Org Name', 'Rendering Provider First Name', 'Rendering Provider Middle Name', 'Name Prefix', 'Rendering Provider Name Suffix', 'ID Code Qualifier', 'Rendering Provider Identifier'] // CMS 835 Flat File, PDF p.7
export const NM1_TT_EL_NAME: string[] = ['COB Carrier Name', 'First name', 'Middle name', 'Name Prefix', 'Name suffix', 'ID Code Qualifier', 'COB Carrier Identifier'] // CMS 835 Flat File, PDF p.7
export const NM1_PR_EL_NAME: string[] = ['Corrected Priority Payer Name', 'First name', 'middle name', 'Name Prefix', 'Name suffix', 'ID Code Qualifier', 'Corrected Priority Payer ID'] // CMS 835 Flat File, PDF p.7-8

export const LX01_EL_NAME = 'Assigned #' // CMS 835 Flat File, PDF p.4

export const TS3_EL_NAME: string[] = [
  'Provider Identifier', 'Facility Code Value', 'Fiscal Period Date',
  'Total Claim Count', 'Total Claim Change Amount', 'Total Covered Charge Amount', 'Total Noncovered Charge Amount',
  'Total Denied Charge Amount', 'Total Provider Amount', 'Total Interest Amount', 'Total Contractual Adjustment Amount',
  'Total Gramm-Rudman Reduction Amount', 'Total MSP Payer Amount', 'Total Blood Deductible Amount', 'Total Non-Lab Charge Amount', 'Total Coinsurance Amount',
  'Total HCPCS Reported Charge Amount', 'Total HCPCS Payable Amount', 'Total Deductible Amount',
  'Total Professional Component Amount', // TS301-320: CMS 835 Flat File, PDF p.4
  'Total MSP Patient Liability Met Amount', 'Total Patient Reimbursement Amount',
  'Total PIP Claim Count', 'Total PIP Adjustment Amount', // TS321-324: p.5
]
export const TS2_EL_NAME: string[] = [
  'Total DRG Amount', 'Total Federal Specific Amount', 'Total Hospital Specific Amount',
  'Total Disproportionate Amount', 'Total Capital Amount', 'Total Indirect Medical Education Amount',
  'Total Outlier Day Count', 'Total Day Outlier Amount', 'Total Cost Outlier Amount',
  'Average DRG Length of Stay', 'Total Discharge Count', 'Total Cost Report Day Count', 'Total Covered Day Count', 'Total Noncovered Day Count', 'Total MSP Pass-Through Amount',
  'Average DRG Weight', 'Total PPS Capital FSP DRG Amount', 'Total PSP Capital HSP DRG Amount', 'Total PPS DSH DRG Amount',
] // CMS 835 Flat File, PDF p.5 (all of TS201-219 fit on one page, immediately before CLP begins)

export const CLP_EL_NAME: string[] = [
  'Patient Control #', 'Claim Status Code', 'Total Claim Charge Amount', // CLP01-03: CMS 835 Flat File, PDF p.5
  'Claim Payment Amount', 'Patient Responsibility Amount',
  'Claim Filing Indicator Code', 'Payer Claim Control #', 'Facility Code Value', 'Claim Frequency Code', 'Patient Status Code',
  'DRG Code', 'DRG Weight', 'Discharge Fraction', 'Yes/No Condition or Response Code', // CLP04-14: p.6
]

export const CAS01_EL_NAME = 'Claim Adjustment Group Code' // CMS 835 Flat File, PDF p.6 (claim loop), p.12 (service loop)
export const CAS_REASON_EL_NAME = 'Adjustment Reason Code' // CMS 835 Flat File, PDF p.6 (claim loop), p.12-13 (service loop)
export const CAS_AMOUNT_EL_NAME = 'Adjustment Amount' // CMS 835 Flat File, PDF p.6 (claim loop), p.12-13 (service loop)
export const CAS_QTY_EL_NAME = 'Adjustment Quantity' // CMS 835 Flat File, PDF p.6 (claim loop), p.12-13 (service loop)

export const MIA_EL_NAME: string[] = [
  'Covered Days or Visits Count', 'PPS Operating Outlier Amount', // MIA01-02: CMS 835 Flat File, PDF p.8
  'Lifetime Psychiatric Days Count',
  'Claim DRG Amount', 'Claim Payment Remark Code', 'Claim DSH Amount', 'Claim MSP Pass Thru Amt', 'Claim PPS Capital Amount',
  'PPS Capital FSP DRG Amt', 'PPS Capital HSP DRG Amt', 'PPS Capital DSH DRG Amt', 'Old Capital Amount', 'PPS Capital IME Amount',
  'PPS Oper HSP Spec DRG Amt', 'Cost Report Day Count', 'PPS Oper FSP Spec DRG Amt', 'Claim PPS Outlier Amount',
  'Claim Indirect Teaching', 'Non Pay Prof Comp Amt',
  'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code',
  'PPS Capital Exception Amt', // MIA03-24: p.9
]
export const MOA_EL_NAME: string[] = [
  'Reimbursement Rate', 'Claim HCPCS Payable Amount',
  'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code', 'Claim Payment Remark Code',
  'Claim ESRD Payment Amount', 'Nonpayable Professional Comp Amt',
] // CMS 835 Flat File, PDF p.10

export const SVC_EL_NAME: string[] = [
  'Composite Medical Procedure Identifier', 'Line Item Charge Amount', 'Line Item Provider Payment', 'NUBC Revenue Code', // SVC01-04: CMS 835 Flat File, PDF p.11
  'Units of Service Paid Count',
  'Composite Medical Procedure Identifier', 'Original Units of Service Count', // SVC05-07: p.12 (SVC06 is the same composite shape as SVC01)
]

export const LQ_EL_NAME: string[] = ['Code List Qualifier Code', 'Remark Code'] // CMS 835 Flat File, PDF p.13

export const PLB01_EL_NAME = 'Provider Identifier' // CMS 835 Flat File, PDF p.13
export const PLB02_EL_NAME = 'Fiscal Period Date' // CMS 835 Flat File, PDF p.13
export const PLB_COMPOSITE_EL_NAME = 'Adjustment Identifier' // CMS 835 Flat File, PDF p.13-14 ("PLB03" etc.; PLB03-1/-2 are the reason code and provider adjustment identifier)
export const PLB_AMOUNT_EL_NAME = 'Provider Adjustment Amount' // CMS 835 Flat File, PDF p.13-14

export const SE_EL_NAME: string[] = ['Transaction Segment Count', 'Transaction Set Control #'] // CMS 835 Flat File, PDF p.14
export const GE_EL_NAME: string[] = ['# Transaction Sets Included', 'Group Control #'] // CMS 835 Flat File, PDF p.14
export const IEA_EL_NAME: string[] = ['# Included Functional Groups', 'Interchange Control #'] // CMS 835 Flat File, PDF p.14-15 (IEA01 on p.14; IEA02 on p.15)

export const RDM_EL_NAME: string[] = ['Report Transmission Code', 'Name', 'Communication Number'] // CMS 835 Flat File, PDF p.4

export const AMT01_EL_NAME = 'Amount Qualifier Code' // CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service loop)
export const CLAIM_AMT02_EL_NAME = 'Claim Supplemental Information Amt' // CMS 835 Flat File, PDF p.11
export const SERVICE_AMT02_EL_NAME = 'Service Supplemental Amount' // CMS 835 Flat File, PDF p.13 (same wording as the segment name at this level — the flat file names both identically)

export const QTY01_EL_NAME = 'Quantity Qualifier' // CMS 835 Flat File, PDF p.11

// CARC/RARC text, embedded verbatim from backend/pkg/carc.json and
// backend/pkg/rarc.json, the maintainers' lists on x12.org as fetched
// 2026-09-02.
// `short` is the table's first-sentence truncation (marked with `*` there);
// `full` is the complete maintainer text, shown on hover under the short one.
export interface CodeText {
  short: string
  full: string
}

export const CARC: Record<string, CodeText> = {
  '1': {
    short: 'Deductible Amount',
    full: 'Deductible Amount',
  },
  '2': {
    short: 'Coinsurance Amount',
    full: 'Coinsurance Amount',
  },
  '23': {
    short: 'The impact of prior payer(s) adjudication including payments and/or adjustments.',
    full:
      'The impact of prior payer(s) adjudication including payments and/or adjustments. (Use only with Group Code OA)',
  },
  '45': {
    short: 'Charge exceeds fee schedule/maximum allowable or contracted/legislated fee arrangement.',
    full:
      'Charge exceeds fee schedule/maximum allowable or contracted/legislated fee arrangement. Usage: This adjustment amount cannot equal the total service or claim charge amount; and must not duplicate provider adjustment amounts (payments and contractual reductions) that have resulted from prior payer(s) adjudication. (Use only with Group Codes PR or CO depending upon liability)',
  },
  '97': {
    short:
      'The benefit for this service is included in the payment/allowance for another service/procedure that has already been adjudicated.',
    full:
      'The benefit for this service is included in the payment/allowance for another service/procedure that has already been adjudicated. Usage: Refer to the 835 Healthcare Policy Identification Segment (loop 2110 Service Payment Information REF), if present.',
  },
  '119': {
    short: 'Benefit maximum for this time period or occurrence has been reached.',
    full: 'Benefit maximum for this time period or occurrence has been reached.',
  },
  '253': {
    short: 'Sequestration - reduction in federal payment',
    full: 'Sequestration - reduction in federal payment',
  },
  A1: {
    short: 'Claim/Service denied.',
    full:
      'Claim/Service denied. At least one Remark Code must be provided (may be comprised of either the NCPDP Reject Reason Code, or Remittance Advice Remark Code that is not an ALERT.) Usage: Use this code only when a more specific Claim Adjustment Reason Code is not available.',
  },
  B16: {
    short: "'New Patient' qualifications were not met.",
    full: "'New Patient' qualifications were not met.",
  },
}

export const RARC: Record<string, CodeText> = {
  MA01: {
    short: 'Alert: If you do not agree with what we approved for these services, you may appeal our decision.',
    full:
      'Alert: If you do not agree with what we approved for these services, you may appeal our decision. To make sure that we are fair to you, we require another individual that did not process your initial claim to conduct the appeal. However, in order to be eligible for an appeal, you must write to us within 120 days of the date you received this notice, unless you have a good reason for being late.',
  },
  M25: {
    short: 'The information furnished does not substantiate the need for this level of service.',
    full:
      'The information furnished does not substantiate the need for this level of service. If you believe the service should have been fully covered as billed, or if you did not know and could not reasonably have been expected to know that we would not pay for this level of service, or if you notified the patient in writing in advance that we would not pay for this level of service and he/she agreed in writing to pay, ask us to review your claim within 120 days of the date of this notice. If you do not request an appeal, we will, upon application from the patient, reimburse him/her for the amount you have collected from him/her in excess of any deductible and coinsurance amounts. We will recover the reimbursement from you as an overpayment.',
  },
  MA18: {
    short: "Alert: The claim information is also being forwarded to the patient's supplemental insurer.",
    full:
      "Alert: The claim information is also being forwarded to the patient's supplemental insurer. Send any questions regarding supplemental benefits to them.",
  },
  MA67: {
    short: 'Alert: Correction to a prior claim.',
    full: 'Alert: Correction to a prior claim.',
  },
  N89: {
    short:
      'Alert: Payment information for this claim has been forwarded to more than one other payer, but format limitations permit only one of the secondary payers to be identified in this remittance advice.',
    full:
      'Alert: Payment information for this claim has been forwarded to more than one other payer, but format limitations permit only one of the secondary payers to be identified in this remittance advice.',
  },
  N362: {
    short: 'The number of Days or Units of Service exceeds our acceptable maximum.',
    full: 'The number of Days or Units of Service exceeds our acceptable maximum.',
  },
}

// Names added 2026-10-04 from the public CMS documents under
// docs/reference/manuals/ and the 835 TI companion guide under
// docs/reference/file-formats/ (the table "Names added from the manuals" in
// docs/design/835-viewer.md lists them with their sources). Citation forms:
// "IOM 100-04 ch.22 (2003 ed.), PDF p.N" is the 82-page 2003 edition of the
// Claims Processing Manual's remittance chapter; "MREP User Guide v4.4
// (2016), PDF p.N" is the Medicare Remit Easy Print user guide, cited by
// its URL on cms.gov (CMS_MREP_GUIDE in x12About835.ts); "CMS 835 TI CG
// (2012), PDF p.N" the 10-page Medicare 835 transaction instructions.

// CLP02 claim status values. 3, 4, and 23 are listed by the flat file and
// named in text by none of the documents on file.
export const CLP02_NAME: Record<string, string> = {
  '1': 'Processed as Primary', // MREP User Guide v4.4 (2016), PDF p.104 (printed p.102)
  '2': 'Processed as Secondary', // MREP User Guide v4.4 (2016), PDF p.103 (printed p.101)
  '19': 'Processed as Primary, Forwarded to Additional Payer(s)', // MREP User Guide v4.4 (2016), PDF p.101 (printed p.99)
  '20': 'Processed as Secondary, Forwarded to Additional Payer(s)', // MREP User Guide v4.4 (2016), PDF p.101 (printed p.99)
  '21': 'Processed as Tertiary, Forwarded to Additional Payer(s)', // MREP User Guide v4.4 (2016), PDF p.101 (printed p.99)
  '22': 'Reversal of previous payment', // MLN Remittance Advice Resources and FAQs (ICN 905367), PDF p.10
}

// CLP06 claim filing indicator values.
export const CLP06_NAME: Record<string, string> = {
  MA: 'Medicare Part A', // CMS 835 TI CG (2012), PDF p.7 ("Required for Part A")
  MB: 'Medicare Part B', // CMS 835 TI CG (2012), PDF p.7 ("Required for Part B"); CEDI 835 CG, PDF p.25 (printed p.19) sends MB for DME
}

// LX01 header number values for Part B and DME.
export const LX01_NAME: Record<string, string> = {
  '1': 'Assigned', // CEDI 835 CG, PDF p.25 (printed p.19; Table 8); MREP User Guide v4.4 (2016), PDF p.214 (printed p.212)
  '0': 'Non-Assigned', // CEDI 835 CG, PDF p.25 (printed p.19; Table 8)
}

// NM102 entity type qualifier, the same X12 code list the 837 uses.
export const NM102_NAME: Record<string, string> = {
  '1': 'Person', // Transmittal 907, PDF p.9
  '2': 'Non-Person Entity', // Transmittal 907, PDF p.9
}

// Further identification code qualifiers (N103 / NM108), from the 837
// documents that name the same X12 code list.
export const ID_QUALIFIER_NAME_MORE: Record<string, string> = {
  MI: 'Member Identification Number', // CEDI 837P CG, PDF p.33 (printed p.25; Table 15)
  PI: 'Payer Identification', // CEDI 837P CG, PDF p.35 (printed p.27; Table 17)
  XV: 'Health Care Financing Administration National PlanID', // Transmittal 907, PDF p.9
}

// AMT01 qualifiers named by the manuals' paper-remittance crosswalks, the
// 837, and the Easy Print guide.
export const CLAIM_AMT_QUALIFIER_NAME: Record<string, string> = {
  AU: 'Covered charges', // IOM 100-04 ch.22 (2003 ed.), PDF p.61 (SPR crosswalk: COVD CHGS, AMT02 when AU)
  DY: 'Per diem rate', // IOM 100-04 ch.22 (2003 ed.), PDF p.62 (SPR crosswalk: PER DIEM RTE, AMT02 when DY)
  F5: 'Patient Amount Paid', // Transmittal 579, PDF p.32 (the 837P's AMT01 F5); IOM 100-04 ch.22 (2003 ed.), PDF p.42
  I: 'Interest', // IOM 100-04 ch.22 (2003 ed.), PDF p.69 (SPR crosswalk: INT, AMT02 when I); MREP User Guide v4.4 (2016), PDF p.217 (printed p.215)
  ZK: 'Hemophilia add-on', // IOM 100-04 ch.22 (2003 ed.), PDF p.26
  ZL: 'New technology add-on', // IOM 100-04 ch.22 (2003 ed.), PDF p.26 ("New Technology Add-on and Electroconvulsive Therapy")
}
export const SERVICE_AMT_QUALIFIER_NAME: Record<string, string> = {
  B6: 'Allowed amount', // IOM 100-04 ch.22 (2003 ed.), PDF p.67 (SPR crosswalk: ALLOWED, AMT02 when B6); MREP User Guide v4.4 (2016), PDF p.215 (printed p.213)
  KH: 'Late filing charge', // IOM 100-04 ch.22 (2003 ed.), PDF p.69 (SPR crosswalk); MREP User Guide v4.4 (2016), PDF p.217 (printed p.215)
}

// QTY01 qualifiers named by the paper-remittance crosswalk.
export const QTY_QUALIFIER_NAME: Record<string, string> = {
  CA: 'Covered days or visits', // IOM 100-04 ch.22 (2003 ed.), PDF p.59 (SPR crosswalk: COVDY, QTY02 when CA)
}

// Line-level DTM01 qualifiers 150 and 151, as the paper remittance reads them.
export const LINE_DTM_NAME: Record<string, string> = {
  '150': 'Service Date, From', // IOM 100-04 ch.22 (2003 ed.), PDF p.67 (SPR crosswalk); MREP User Guide v4.4 (2016), PDF p.215 (printed p.213)
  '151': 'Service Date, Through', // IOM 100-04 ch.22 (2003 ed.), PDF p.67 (SPR crosswalk); MREP User Guide v4.4 (2016), PDF p.215 (printed p.213)
}

// REF01 qualifiers named by the manuals and the TI companion guide.
export const REF_QUALIFIER_NAME: Record<string, string> = {
  TJ: 'Taxpayer Identification Number', // CMS 835 TI CG (2012), PDF p.6 ("Required for Medicare to report the Taxpayer Identification Number (TIN)")
  '1C': 'Medicare Provider Number', // IOM 100-04 ch.22 (2003 ed.), PDF p.37, 45
  '1G': 'UPIN', // IOM 100-04 ch.22 (2003 ed.), PDF p.37
  RB: 'Rate code', // IOM 100-04 ch.22 (2003 ed.), PDF p.30 ("when a rate code factored in the payment")
  APC: 'Ambulatory Payment Classification number', // IOM 100-04 ch.22 (2003 ed.), PDF p.30 ("The APC number will only be reported with the first HCPCS")
}
