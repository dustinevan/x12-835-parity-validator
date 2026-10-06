// Every element and coded value of the Medicare 837I (institutional claim,
// 005010X223A2), built only from x12Names837i.ts (names and citations read
// from Transmittal 907, the CMS Claims Processing Manual Ch. 25, and X12's
// public pages; read that file's own header first for the exact citation
// convention and its NUBC-code-list note). Read by the 837I viewer's lookup
// panel (frontend/src/docs/viewer/Lookup.tsx, opened from
// docs/pages/Viewer837I.tsx). Nothing here is invented: an element the
// public documents do not name has no `name`, and the panel shows its
// reference instead; a NUBC-maintained code list (Type of Bill, Admission
// Type/Source, Patient Status, Condition/Occurrence/Value Code) is not
// public, so its codes keep their bare reference too.
//
// Element numbering follows x12Names837i.ts's own array order (position 1,
// 2, 3…), matching frontend/src/data/lookupRows835.ts's convention, except
// where the file's own comment gives an explicit composite breakdown
// ("CLM05 composite", HI0n-1/HI0n-2…, the NM1 role arrays, which start at
// position 3). A few arrays (CLM, SV2, SVD, HCP, MIA) mix a composite's own
// descriptive label in with its numbered sub-parts in a way Transmittal 907
// does not resolve on the page cited; those are numbered by plain counting
// and may not match the X12 standard's own element number exactly, while
// the name and citation text stay the public document's own.
import * as I from './x12Names837i'
import type { LookupCodeRow, LookupElementRow } from '../docs/viewer/Lookup'

export const TRANSMITTAL_907 = 'https://www.cms.gov/files/document/cms-837-835-flatfiles-r907otn-2011.pdf'
export const X12_837I_EXAMPLES = 'https://x12.org/examples/005010x223'
export const CMS_CLAIMS_MANUAL_CH25 = 'https://www.cms.gov/regulations-and-guidance/guidance/manuals/downloads/clm104c25.pdf'

const pad = (n: number) => String(n).padStart(2, '0')
const seq = (n: number): string[] => Array.from({ length: n }, (_, i) => pad(i + 1))

interface RowContext {
  segmentName?: string
}

function elementRowsWithRefs(segment: string, where: string, refs: string[], names: (string | undefined)[], source: string, context: RowContext = {}): LookupElementRow[] {
  return names.map((name, i) => {
    const row: LookupElementRow = { ref: `${segment}${refs[i]}`, segment, where, source }
    if (name) row.name = name
    if (context.segmentName) row.segmentName = context.segmentName
    return row
  })
}

function elementRows(segment: string, where: string, names: (string | undefined)[], source: string, context: RowContext = {}): LookupElementRow[] {
  return elementRowsWithRefs(segment, where, seq(names.length), names, source, context)
}

// `segmentName`, when given, must equal the one element row's own
// `segmentName` this code belongs under — frontend/src/docs/viewer/
// Lookup.tsx nests a code only under the element row(s) sharing it, so two
// boxes that reuse the same bare ref (REF01 at a dozen loops, NM101 at
// every NM1) don't bleed into each other. Left out only for a code
// genuinely shared by more than one box (NM101 role codes, CAS group codes
// at both the claim and the line level).
function codeRow(segment: string, element: string, value: string, name: string | undefined, where: string, source: string, full?: string, segmentName?: string): LookupCodeRow {
  const row: LookupCodeRow = { segment, element, code: value, where, source }
  if (name) row.name = name
  if (full) row.full = full
  if (segmentName) row.segmentName = segmentName
  return row
}

function codesFromRecord(segment: string, element: string, names: Partial<Record<string, string>> | undefined, where: string, source: string, segmentName?: string): LookupCodeRow[] {
  if (!names) return []
  return Object.entries(names)
    .filter((e): e is [string, string] => e[1] !== undefined)
    .map(([value, name]) => codeRow(segment, element, value, name, where, source, undefined, segmentName))
}

function valuesOnly(segment: string, element: string, values: string[] | undefined, where: string, source: string, segmentName?: string): LookupCodeRow[] {
  return (values ?? []).map((value) => codeRow(segment, element, value, undefined, where, source, undefined, segmentName))
}

const t907 = (pages: string) => `Transmittal 907, PDF p.${pages}`
const ch25 = (section: string, fl: string) => `CMS Claims Processing Manual Ch. 25, §${section} (${fl})`

const ENVELOPE = 'Envelope'
const SUBMITTER = 'Loop 1000A · submitter'
const RECEIVER = 'Loop 1000B · receiver'
const BILLING_HL = 'Loop 2000A · billing provider'
const BILLING_PROVIDER = 'Loop 2010AA · billing provider'
const PAY_TO_ADDRESS = 'Loop 2010AB · pay-to address'
const PAY_TO_PLAN = 'Loop 2010AC · pay-to plan'
const SUBSCRIBER_HL = 'Loop 2000B · subscriber'
const SUBSCRIBER = 'Loop 2010BA · subscriber'
const PAYER = 'Loop 2010BB · payer'
const PATIENT_HL = 'Loop 2000C · patient (not used by Medicare)'
const PATIENT = 'Loop 2010CA · patient (not used by Medicare)'
const CLAIM = 'Loop 2300 · claim'
const ATTENDING = 'Loop 2310A · attending provider'
const OPERATING = 'Loop 2310B · operating physician'
const OTHER_OPERATING = 'Loop 2310C · other operating physician'
const RENDERING = 'Loop 2310D · rendering provider'
const SERVICE_FACILITY = 'Loop 2310E · service facility location'
const REFERRING = 'Loop 2310F · referring provider'
const OTHER_SUBSCRIBER = 'Loop 2320/2330A · other subscriber (coordination of benefits)'
const OTHER_PAYER = 'Loop 2330B · other payer (coordination of benefits)'
const OTHER_PAYER_PROVIDER = 'Loop 2330C-I · other payer provider loops (coordination of benefits)'
const SERVICE_LINE = 'Loop 2400 · service line'
const LINE_OPERATING = 'Loop 2420A · operating physician'
const LINE_OTHER_OPERATING = 'Loop 2420B · other operating physician'
const LINE_RENDERING = 'Loop 2420C · rendering provider'
const LINE_REFERRING = 'Loop 2420D · referring provider'
const LINE_ADJUDICATION = 'Loop 2430 · line adjudication (coordination of benefits)'
const PROVIDER_LOOPS = 'Loop 2310A-F/2420A-D · provider secondary identification'
const TRAILER = 'Trailer'

// NM1: the generic NM101/NM102 positions, plus this loop's own NM103-109
// (every 837I array is a full 7-item, gap-free list, unlike 837P's — no
// position list needed beyond starting at 3).
const NM1_REFS = ['03', '04', '05', '06', '07', '08', '09']
function nm1Row(where: string, names: string[], source: string, segmentName: string): LookupElementRow[] {
  return elementRowsWithRefs('NM1', where, ['01', '02', ...NM1_REFS], [I.NM101_EL_NAME, I.NM102_EL_NAME, ...names], source, { segmentName })
}

// 2300 HI: one 5-position composite (type qualifier, code, date format,
// date, amount) repeated up to 12 times; the composite's 9th position
// (present-on-admission) applies only to the diagnosis-bearing instances.
function hiRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (let n = 1; n <= 12; n++) {
    out.push(...elementRowsWithRefs('HI', CLAIM, [`${pad(n)}-1`, `${pad(n)}-2`, `${pad(n)}-3`, `${pad(n)}-4`, `${pad(n)}-5`], I.HI_EL_NAME, t907('24'), { segmentName: `Health Care Information Codes, repeat ${n} of 12` }))
    out.push(...elementRowsWithRefs('HI', CLAIM, [`${pad(n)}-9`], [I.HI_POA_EL_NAME], t907('30'), { segmentName: `Health Care Information Codes, repeat ${n} of 12 (diagnosis composites only)` }))
  }
  return out
}

// CAS (2320 claim-level, 2430 line-level): a group code, then up to six
// reason/amount/quantity triplets — the same shape
// frontend/src/data/lookupRows835.ts and lookupRows837p.ts already give the
// 835's and 837P's own CAS; x12Names837i.ts does not state a repeat limit,
// so this follows the standard X12 maximum (six) rather than guess lower.
function casRows(where: string, source: string, segmentName: string): LookupElementRow[] {
  const refs = ['01']
  const names: (string | undefined)[] = [I.CAS01_EL_NAME]
  for (let i = 0; i < 6; i++) {
    const n = 2 + i * 3
    refs.push(pad(n), pad(n + 1), pad(n + 2))
    names.push(I.CAS_REASON_EL_NAME, I.CAS_AMOUNT_EL_NAME, I.CAS_QTY_EL_NAME)
  }
  return elementRowsWithRefs('CAS', where, refs, names, source, { segmentName })
}

export const ELEMENTS_837I: LookupElementRow[] = [
  // ---- Envelope ----
  ...elementRows('ISA', ENVELOPE, I.ISA_EL_NAME, t907('7'), { segmentName: I.SEGMENT_NAME.ISA }),
  ...elementRows('GS', ENVELOPE, I.GS_EL_NAME, t907('7-8'), { segmentName: I.SEGMENT_NAME.GS }),
  ...elementRows('ST', ENVELOPE, I.ST_EL_NAME, t907('8'), { segmentName: I.SEGMENT_NAME.ST }),
  ...elementRows('BHT', ENVELOPE, I.BHT_EL_NAME, t907('8-9'), { segmentName: I.SEGMENT_NAME.BHT }),
  ...elementRows('SE', TRAILER, I.SE_EL_NAME, t907('124'), { segmentName: I.SEGMENT_NAME.SE }),
  ...elementRows('GE', TRAILER, I.GE_EL_NAME, t907('124'), { segmentName: I.SEGMENT_NAME.GE }),
  ...elementRows('IEA', TRAILER, I.IEA_EL_NAME, t907('125'), { segmentName: I.SEGMENT_NAME.IEA }),

  // ---- 1000A/1000B submitter and receiver ----
  ...nm1Row(SUBMITTER, I.NM1_SUBMITTER_EL_NAME, t907('9'), I.NM1_ROLE_NAME['41']),
  ...nm1Row(RECEIVER, I.NM1_RECEIVER_EL_NAME, t907('10'), I.NM1_ROLE_NAME['40']),
  ...elementRows('PER', SUBMITTER, I.PER_EL_NAME, t907('9-10'), { segmentName: 'Submitter EDI Contact Information' }),
  ...elementRows('PER', BILLING_PROVIDER, I.PER_EL_NAME, t907('13, identical shape to the submitter PER'), { segmentName: 'Billing Provider Contact Information' }),

  // ---- 2000A billing provider, 2010AA/AB/AC ----
  ...elementRows('HL', BILLING_HL, I.HL_EL_NAME, t907('10'), { segmentName: I.HL03_NAME['20'] }),
  ...elementRows('PRV', BILLING_HL, I.PRV_EL_NAME, t907('10'), { segmentName: 'Billing Provider Specialty Information' }),
  ...elementRows('PRV', ATTENDING, I.PRV_EL_NAME, t907('92, identical shape to the billing provider PRV'), { segmentName: 'Attending Provider Specialty Information' }),
  ...elementRows('CUR', BILLING_HL, I.CUR_EL_NAME, t907('11'), { segmentName: 'Foreign Currency Information' }),
  ...nm1Row(BILLING_PROVIDER, I.NM1_BILLING_EL_NAME, t907('12'), I.NM1_ROLE_NAME['85']),
  ...elementRows('N3', BILLING_PROVIDER, I.N3_EL_NAME, t907('12'), { segmentName: 'Billing Provider Address' }),
  ...elementRows('N4', BILLING_PROVIDER, I.N4_EL_NAME, t907('12'), { segmentName: 'Billing Provider City, State, ZIP Code' }),
  ...elementRows('REF', BILLING_PROVIDER, [undefined, undefined], t907('13'), { segmentName: I.REF_BILLING_TAX_NAME }),
  ...nm1Row(PAY_TO_ADDRESS, I.NM1_PAYTO_EL_NAME, t907('13'), I.NM1_ROLE_NAME['87']),
  ...nm1Row(PAY_TO_PLAN, I.NM1_PAYTOPLAN_EL_NAME, t907('14-15'), I.NM1_ROLE_NAME.PE),
  ...elementRows('REF', PAY_TO_PLAN, [undefined, undefined], t907('15-16'), { segmentName: 'Pay-to Plan Secondary/Tax Identification (by REF01 code)' }),

  // ---- 2000B/C, 2010BA/BB/CA subscriber and patient ----
  ...elementRows('HL', SUBSCRIBER_HL, I.HL_EL_NAME, t907('16'), { segmentName: I.HL03_NAME['22'] }),
  ...elementRows('SBR', SUBSCRIBER_HL, I.SBR_EL_NAME, t907('16'), { segmentName: 'Subscriber Information' }),
  ...elementRows('SBR', OTHER_SUBSCRIBER, I.SBR_OTHER_EL_NAME, t907('95-96'), { segmentName: 'Other Subscriber Information' }),
  ...nm1Row(SUBSCRIBER, I.NM1_SUBSCRIBER_EL_NAME, t907('16-17'), I.NM1_ROLE_NAME.IL),
  ...elementRows('DMG', SUBSCRIBER, I.DMG_EL_NAME, t907('17'), { segmentName: 'Subscriber Demographic Information' }),
  ...elementRows('REF', SUBSCRIBER, [undefined, undefined], t907('18'), { segmentName: 'Subscriber Secondary Identification / Property and Casualty Claim Number (by REF01 code)' }),
  ...nm1Row(PAYER, I.NM1_PAYER_EL_NAME, t907('18-19'), I.NM1_ROLE_NAME.PR),
  ...elementRows('REF', PAYER, [undefined, undefined], t907('19-20'), { segmentName: 'Payer / Billing Provider Secondary Identification (by REF01 code)' }),
  ...elementRows('HL', PATIENT_HL, I.HL_EL_NAME, t907('20'), { segmentName: I.HL03_NAME['23'] }),
  ...elementRows('PAT', PATIENT_HL, I.PAT_EL_NAME, t907('20'), { segmentName: 'Patient Information' }),
  ...nm1Row(PATIENT, I.NM1_PATIENT_EL_NAME, t907('20-21'), I.NM1_ROLE_NAME.QC),
  ...elementRows('REF', PATIENT, [undefined, undefined], t907('22-23'), { segmentName: 'Property and Casualty Patient Identifier/Claim Number (by REF01 code)' }),

  // ---- 2300 claim ----
  ...elementRows('CLM', CLAIM, I.CLM_EL_NAME, t907('23-24'), { segmentName: 'Claim Information' }),
  ...elementRows('DTP', CLAIM, [I.DTP01_EL_NAME, undefined, undefined], t907('24'), { segmentName: I.DTP_CLAIM_NAME['096'] }),
  ...elementRows('DTP', CLAIM, [I.DTP01_EL_NAME, undefined, undefined], t907('24'), { segmentName: I.DTP_CLAIM_NAME['434'] }),
  ...elementRows('DTP', CLAIM, [I.DTP01_EL_NAME, undefined, undefined], t907('24'), { segmentName: I.DTP_CLAIM_NAME['435'] }),
  ...elementRows('DTP', CLAIM, [I.DTP01_EL_NAME, undefined, undefined], t907('24'), { segmentName: I.DTP_CLAIM_NAME['050'] }),
  ...elementRows('CL1', CLAIM, I.CL1_EL_NAME, t907('24-25'), { segmentName: I.SEGMENT_NAME.CL1 }),
  ...elementRows('PWK', CLAIM, I.PWK_EL_NAME, t907('25'), { segmentName: I.SEGMENT_NAME.PWK }),
  ...elementRows('PWK', SERVICE_LINE, I.PWK_EL_NAME, t907('106-107'), { segmentName: 'Line Supplemental Information' }),
  ...elementRows('CN1', CLAIM, I.CN1_EL_NAME, t907('25-26'), { segmentName: I.SEGMENT_NAME.CN1 }),
  ...elementRows('AMT', CLAIM, [undefined, undefined], t907('26'), { segmentName: 'Patient Estimated Amount Due / COB Amounts / Tax Amounts (by AMT01 code)' }),
  ...elementRows('REF', CLAIM, [undefined, undefined], t907('26-28'), { segmentName: 'Claim-Level Reference Identification (by REF01 code)' }),
  ...elementRows('K3', CLAIM, I.K3_EL_NAME, t907('29'), { segmentName: I.SEGMENT_NAME.K3 }),
  ...elementRows('NTE', CLAIM, I.NTE_EL_NAME, t907('29'), { segmentName: I.NTE_CLAIM_NAME }),
  ...elementRows('NTE', CLAIM, I.NTE_EL_NAME, t907('29'), { segmentName: I.NTE_BILLING_NAME }),
  ...elementRows('NTE', SERVICE_LINE, I.NTE_EL_NAME, t907('108'), { segmentName: I.NTE_LINE_NAME }),
  ...elementRows('CRC', CLAIM, I.CRC_EL_NAME, t907('29'), { segmentName: I.SEGMENT_NAME.CRC }),
  ...hiRows(),
  ...elementRows('HCP', CLAIM, I.HCP_EL_NAME, t907('91-92'), { segmentName: 'Claim Pricing/Repricing Information' }),
  ...elementRows('HCP', SERVICE_LINE, I.HCP_EL_NAME, t907('113-114'), { segmentName: 'Line Pricing/Repricing Information' }),

  // ---- 2310A-F claim-level providers ----
  ...nm1Row(ATTENDING, I.NM1_ATTENDING_EL_NAME, t907('92'), I.NM1_ROLE_NAME['71']),
  ...nm1Row(OPERATING, I.NM1_OPERATING_EL_NAME, t907('93'), I.NM1_ROLE_NAME['72']),
  ...nm1Row(OTHER_OPERATING, I.NM1_OTHER_OPERATING_EL_NAME, t907('93'), I.NM1_ROLE_NAME.ZZ),
  ...nm1Row(RENDERING, I.NM1_RENDERING_EL_NAME, t907('94'), I.NM1_ROLE_NAME['82']),
  ...nm1Row(SERVICE_FACILITY, I.NM1_SERVICE_FACILITY_EL_NAME, t907('94'), I.NM1_ROLE_NAME['77']),
  ...nm1Row(REFERRING, I.NM1_REFERRING_EL_NAME, t907('94-95'), I.NM1_ROLE_NAME.DN),
  ...elementRows('REF', ATTENDING, I.REF_PROVIDER_EL_NAME, t907('92, representative of every provider-loop REF'), { segmentName: 'Provider Secondary Identification' }),

  // ---- 2320/2330 other subscriber and payer (COB) ----
  ...elementRows('MIA', OTHER_SUBSCRIBER, I.MIA_EL_NAME, t907('100-101'), { segmentName: 'Inpatient Adjudication Information' }),
  ...elementRows('MOA', OTHER_SUBSCRIBER, I.MOA_EL_NAME, t907('102'), { segmentName: 'Outpatient Adjudication (MOA) Information' }),
  ...elementRows('OI', OTHER_SUBSCRIBER, I.OI_EL_NAME, t907('101'), { segmentName: 'Other Health Insurance Information' }),
  ...casRows(OTHER_SUBSCRIBER, t907('92'), I.CAS_CLAIM_NAME),
  ...nm1Row(OTHER_SUBSCRIBER, I.NM1_OTHER_SUBSCRIBER_EL_NAME, t907('96'), I.NM1_ROLE_NAME.IL),
  ...elementRows('REF', OTHER_SUBSCRIBER, [undefined, undefined], t907('97'), { segmentName: 'Other Subscriber Secondary Identification (by REF01 code)' }),
  ...nm1Row(OTHER_PAYER, I.NM1_OTHER_PAYER_EL_NAME, t907('97-98'), I.NM1_ROLE_NAME.PR),
  ...elementRows('DTP', OTHER_PAYER, [I.DTP01_EL_NAME, undefined, undefined], t907('99'), { segmentName: I.DTP_OTHER_PAYER_NAME['573'] }),
  ...elementRows('REF', OTHER_PAYER, [undefined, undefined], t907('99-100'), { segmentName: 'Other Payer Secondary Identification / Prior Authorization / Referral / Claim Adjustment Indicator / Claim Control Number (by REF01 code)' }),
  ...Object.entries(I.OTHER_PAYER_PROVIDER_LOOP_NAME).map(([loopId, segmentName]) =>
    nm1Row(`${OTHER_PAYER_PROVIDER} (${loopId})`, I.NM1_ATTENDING_EL_NAME, t907('100-104, reusing its 2310-series loop\'s own NM1 shape'), segmentName),
  ).flat(),

  // ---- 2400 service line ----
  ...elementRows('LX', SERVICE_LINE, [I.LX01_EL_NAME], t907('105'), { segmentName: I.SEGMENT_NAME.LX }),
  ...elementRows('SV2', SERVICE_LINE, I.SV2_EL_NAME, t907('105-106'), { segmentName: I.SEGMENT_NAME.SV2 }),
  ...elementRows('DTP', SERVICE_LINE, [I.DTP01_EL_NAME, undefined, undefined], t907('107'), { segmentName: I.DTP472_NAME }),
  ...elementRows('REF', SERVICE_LINE, [undefined, undefined], t907('107'), { segmentName: 'Line Item Control / Repriced Line Item Reference Number (by REF01 code)' }),
  ...elementRows('LIN', SERVICE_LINE, I.LIN_EL_NAME, t907('109'), { segmentName: I.SEGMENT_NAME.LIN }),
  ...elementRowsWithRefs('CTP', SERVICE_LINE, ['01', '02', '03', '04', '05-1'], I.CTP_EL_NAME, t907('109-110'), { segmentName: I.SEGMENT_NAME.CTP }),
  ...elementRows('REF', SERVICE_LINE, [undefined, undefined], t907('110-111'), { segmentName: 'Prescription or Compound Drug Association Number (by REF01 code)' }),
  ...nm1Row(LINE_OPERATING, I.NM1_OPERATING_EL_NAME, t907('112'), I.NM1_ROLE_NAME['72']),
  ...nm1Row(LINE_OTHER_OPERATING, I.NM1_OTHER_OPERATING_EL_NAME, t907('112-113'), I.NM1_ROLE_NAME.ZZ),
  ...nm1Row(LINE_RENDERING, I.NM1_RENDERING_EL_NAME, t907('114'), I.NM1_ROLE_NAME['82']),
  ...nm1Row(LINE_REFERRING, I.NM1_REFERRING_EL_NAME, t907('114-115'), I.NM1_ROLE_NAME.DN),
  ...elementRows('REF', LINE_OPERATING, I.REF_PROVIDER_EL_NAME, t907('92, representative of every provider-loop REF'), { segmentName: 'Provider Secondary Identification' }),
  ...elementRows('SVD', LINE_ADJUDICATION, I.SVD_EL_NAME, t907('115-116'), { segmentName: I.SEGMENT_NAME.SVD }),
  ...casRows(LINE_ADJUDICATION, t907('116'), I.CAS_LINE_NAME),
  ...elementRows('DTP', LINE_ADJUDICATION, [I.DTP01_EL_NAME, undefined, undefined], t907('118'), { segmentName: I.DTP573_LINE_NAME }),
  ...elementRows('AMT', LINE_ADJUDICATION, [undefined, undefined], t907('118'), { segmentName: I.AMT_LINE_LIABILITY_NAME }),
]

export const CODES_837I: LookupCodeRow[] = [
  ...codesFromRecord('REF', 'REF01', I.REF_PAYTOPLAN_NAME, PAY_TO_PLAN, t907('15-16'), 'Pay-to Plan Secondary/Tax Identification (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_SUBSCRIBER_NAME, SUBSCRIBER, t907('18'), 'Subscriber Secondary Identification / Property and Casualty Claim Number (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_PAYER_NAME, PAYER, t907('19-20'), 'Payer / Billing Provider Secondary Identification (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_PATIENT_NAME, PATIENT, t907('22-23'), 'Property and Casualty Patient Identifier/Claim Number (by REF01 code)'),
  // AMT01 is left broad: F3/D/EAF/A8/GT/N8 span several AMT boxes (claim,
  // other subscriber, service line) this file gives only one generic row
  // each for.
  ...codesFromRecord('AMT', 'AMT01', I.AMT01_NAME, `${CLAIM} · ${OTHER_SUBSCRIBER} · ${SERVICE_LINE}`, t907('26, 100-101, 108')),
  ...codesFromRecord('REF', 'REF01', I.REF_CLAIM_NAME, CLAIM, t907('26-28'), 'Claim-Level Reference Identification (by REF01 code)'),
  // PROVIDER_REF_NAME is left broad: the same 0B/1G/G2 codes apply at every
  // one of the eight provider loops this file gives one representative REF
  // row for.
  ...codesFromRecord('REF', 'REF01', I.PROVIDER_REF_NAME, PROVIDER_LOOPS, `${ch25('75.6', 'FL 76/77/78-79')}`, 'Provider Secondary Identification'),
  ...codesFromRecord('OI', 'OI03', I.OI03_NAME, OTHER_SUBSCRIBER, t907('101'), 'Other Health Insurance Information'),
  ...codesFromRecord('OI', 'OI06', I.OI06_NAME, OTHER_SUBSCRIBER, t907('101'), 'Other Health Insurance Information'),
  ...codesFromRecord('REF', 'REF01', I.REF_OTHER_SUBSCRIBER_NAME, OTHER_SUBSCRIBER, t907('97'), 'Other Subscriber Secondary Identification (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_OTHER_PAYER_NAME, OTHER_PAYER, t907('99-100'), 'Other Payer Secondary Identification / Prior Authorization / Referral / Claim Adjustment Indicator / Claim Control Number (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_LINE_NAME, SERVICE_LINE, t907('107'), 'Line Item Control / Repriced Line Item Reference Number (by REF01 code)'),
  ...codesFromRecord('REF', 'REF01', I.REF_PRESCRIPTION_NAME, SERVICE_LINE, t907('110-111'), 'Prescription or Compound Drug Association Number (by REF01 code)'),
  ...codesFromRecord('BHT', 'BHT01', I.BHT01_NAME, ENVELOPE, t907('8'), I.SEGMENT_NAME.BHT),
  ...codesFromRecord('BHT', 'BHT02', I.BHT02_NAME, ENVELOPE, t907('8'), I.SEGMENT_NAME.BHT),
  ...valuesOnly('BHT', 'BHT06', I.BHT06_VALUES, ENVELOPE, t907('9'), I.SEGMENT_NAME.BHT),
  codeRow('GS', 'GS01', 'HC', I.TRANSACTION_TYPE_NAME.HC, ENVELOPE, t907('7'), undefined, I.SEGMENT_NAME.GS),
  codeRow('ST', 'ST01', '837', I.TRANSACTION_TYPE_NAME['837'], ENVELOPE, t907('8'), undefined, I.SEGMENT_NAME.ST),
  // NM101/NM102/NM108 are left broad (no segmentName): each applies to
  // every NM1 occurrence that carries that position, across every role.
  ...codesFromRecord('NM1', 'NM101', I.NM1_ROLE_NAME, 'Every NM1 instance', t907('9-104 (first use per role cited in x12Names837i.ts)')),
  ...codesFromRecord('NM1', 'NM102', I.NM102_NAME, 'Every NM1 instance', t907('9')),
  ...codesFromRecord('NM1', 'NM108', I.NM108_NAME, 'Every NM1 instance', t907('9, 12, 17')),
  codeRow('PRV', 'PRV01', 'BI', I.PRV01_NAME.BI, BILLING_HL, t907('10'), undefined, 'Billing Provider Specialty Information'),
  codeRow('PRV', 'PRV01', 'AT', I.PRV01_NAME.AT, ATTENDING, t907('92'), undefined, 'Attending Provider Specialty Information'),
  ...codesFromRecord('PRV', 'PRV02', I.PRV02_NAME, `${BILLING_HL} · ${ATTENDING}`, t907('10')),
  codeRow('CUR', 'CUR01', '85', I.CUR01_NAME['85'], BILLING_HL, t907('11'), undefined, 'Foreign Currency Information'),
  // PER01/PER03 are left broad: the same shape repeats, unnamed, at the
  // submitter and the billing provider.
  ...codesFromRecord('PER', 'PER01', I.PER01_NAME, `${SUBMITTER} · ${BILLING_PROVIDER}`, t907('9')),
  ...valuesOnly('PER', 'PER03', I.PER_QUALIFIER_VALUES, `${SUBMITTER} · ${BILLING_PROVIDER}`, t907('9-10')),
  ...codesFromRecord('HL', 'HL03', I.HL03_NAME, 'Every HL instance', t907('10, 16, 20')),
  // SBR01 is left broad: the same Payer Responsibility Sequence Number code
  // list applies at both the subscriber (2000B) and other-subscriber (2320)
  // SBR, which share this file's one SBR_EL_NAME shape.
  codeRow('SBR', 'SBR01', 'A', I.SBR01_NAME.A, `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, t907('16')),
  codeRow('SBR', 'SBR01', 'P', I.SBR01_NAME.P, `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, 'common X12 Payer Responsibility Sequence Number code'),
  codeRow('SBR', 'SBR01', 'S', I.SBR01_NAME.S, `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, 'common X12 Payer Responsibility Sequence Number code'),
  codeRow('SBR', 'SBR01', 'T', I.SBR01_NAME.T, `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, 'common X12 Payer Responsibility Sequence Number code'),
  ...valuesOnly('SBR', 'SBR01', I.SBR01_VALUES.filter((v) => !(v in I.SBR01_NAME)), `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, t907('16')),
  ...codesFromRecord('SBR', 'SBR02', I.RELATIONSHIP_NAME, `${SUBSCRIBER_HL} · ${PATIENT_HL} · ${OTHER_SUBSCRIBER}`, 'common X12 Individual Relationship Code; Transmittal 907, PDF p.16 (2000B SBR02, only value 18)'),
  ...codesFromRecord('PAT', 'PAT01', I.RELATIONSHIP_NAME, PATIENT_HL, t907('20'), 'Patient Information'),
  ...codesFromRecord('SBR', 'SBR09', I.SBR09_NAME, `${SUBSCRIBER_HL} · ${OTHER_SUBSCRIBER}`, t907('16, 96')),
  codeRow('DMG', 'DMG01', 'D8', I.DMG01_NAME.D8, `${SUBSCRIBER} · ${PATIENT}`, t907('17')),
  ...codesFromRecord('DMG', 'DMG03', I.DMG03_NAME, `${SUBSCRIBER} · ${PATIENT}`, t907('17')),
  ...codesFromRecord('CLM', 'CLM05-2', I.CLM05_2_NAME, CLAIM, t907('23'), 'Claim Information'),
  ...codesFromRecord('CLM', 'CLM07', I.CLM07_NAME, CLAIM, t907('23'), 'Claim Information'),
  ...codesFromRecord('CLM', 'CLM08', I.CLM08_NAME, CLAIM, t907('23'), 'Claim Information'),
  ...codesFromRecord('CLM', 'CLM09', I.CLM09_NAME, CLAIM, t907('23'), 'Claim Information'),
  ...valuesOnly('CLM', 'CLM20', I.CLM20_VALUES, CLAIM, t907('24'), 'Claim Information'),
  ...valuesOnly('PWK', 'PWK01', I.PWK01_VALUES_CLAIM, CLAIM, t907('25'), I.SEGMENT_NAME.PWK),
  ...valuesOnly('PWK', 'PWK02', I.PWK02_VALUES_CLAIM, CLAIM, t907('25'), I.SEGMENT_NAME.PWK),
  ...valuesOnly('PWK', 'PWK01', I.PWK01_VALUES_LINE, SERVICE_LINE, t907('106'), 'Line Supplemental Information'),
  ...valuesOnly('PWK', 'PWK02', I.PWK02_VALUES_LINE, SERVICE_LINE, t907('106'), 'Line Supplemental Information'),
  ...valuesOnly('CN1', 'CN101', I.CN101_VALUES, CLAIM, t907('25'), I.SEGMENT_NAME.CN1),
  ...valuesOnly('REF', 'REF02', I.REF_SERVICE_AUTH_VALUES, CLAIM, t907('26'), 'Claim-Level Reference Identification (by REF01 code)'),
  ...valuesOnly('NTE', 'NTE01', I.NTE01_VALUES_CLAIM, CLAIM, t907('29'), I.NTE_CLAIM_NAME),
  ...valuesOnly('NTE', 'NTE01', I.NTE01_VALUES_BILLING, CLAIM, t907('29'), I.NTE_BILLING_NAME),
  ...valuesOnly('NTE', 'NTE01', I.NTE01_VALUES_LINE, SERVICE_LINE, t907('108'), I.NTE_LINE_NAME),
  codeRow('CRC', 'CRC01', 'ZZ', I.CRC01_NAME.ZZ, CLAIM, t907('29'), undefined, I.SEGMENT_NAME.CRC),
  ...codesFromRecord('CRC', 'CRC02', I.CRC02_NAME, CLAIM, t907('29'), I.SEGMENT_NAME.CRC),
  ...valuesOnly('CRC', 'CRC03', I.CRC_CONDITION_VALUES, CLAIM, t907('29'), I.SEGMENT_NAME.CRC),
  ...codesFromRecord('HI', 'HI01-1', I.HI_TYPE_NAME, CLAIM, t907('30-85 (first use per type cited in x12Names837i.ts)')),
  ...codesFromRecord('HI', 'HI01-9', I.HI_POA_NAME, CLAIM, t907('30')),
  ...valuesOnly('HCP', 'HCP01', I.HCP01_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('91')),
  codeRow('HCP', 'HCP01', '00', I.HCP01_NAME['00'], `${CLAIM} · ${SERVICE_LINE}`, 'common X12 Pricing Methodology code'),
  ...valuesOnly('HCP', 'HCP09', I.HCP09_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('92')),
  ...valuesOnly('HCP', 'HCP11', I.HCP11_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('92')),
  ...valuesOnly('HCP', 'HCP13', I.HCP13_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('92')),
  ...valuesOnly('HCP', 'HCP14', I.HCP14_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('92')),
  ...valuesOnly('HCP', 'HCP15', I.HCP15_VALUES, `${CLAIM} · ${SERVICE_LINE}`, t907('92')),
  ...codesFromRecord('CAS', 'CAS01', I.CAS_GROUP_NAME, `${OTHER_SUBSCRIBER} · ${LINE_ADJUDICATION}`, 'https://x12.org/codes/claim-adjustment-group-codes'),
  // SV202-1 and SVD03-1 are composite sub-positions inside SV2_EL_NAME/
  // SVD_EL_NAME that this file numbers by plain counting (see the file-level
  // note at the top); they show as unattached coded values in the SV2/SVD
  // group rather than nested under a specific row.
  ...valuesOnly('SV2', 'SV202-1', I.SV202_1_VALUES, SERVICE_LINE, t907('106')),
  ...valuesOnly('CTP', 'CTP05-1', I.CTP05_1_VALUES, SERVICE_LINE, t907('110'), I.SEGMENT_NAME.CTP),
  ...codesFromRecord('LIN', 'LIN02', I.LIN02_NAME, SERVICE_LINE, t907('109'), I.SEGMENT_NAME.LIN),
  ...valuesOnly('SVD', 'SVD03-1', I.SVD03_1_VALUES, LINE_ADJUDICATION, t907('115')),
]
