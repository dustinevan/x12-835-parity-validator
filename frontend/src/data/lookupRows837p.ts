// Every element and coded value of the Medicare 837P (professional claim,
// 005010X222A1), built only from x12Names837p.ts (names and citations read
// from Transmittal 907 and the CEDI 837P companion guide; read that file's
// own header first for the exact citation convention). Read by the 837P
// viewer's lookup panel (frontend/src/docs/viewer/Lookup.tsx, opened from
// docs/pages/Viewer837P.tsx). Nothing here is invented: an element the
// public documents do not name has no `name`, and the panel shows its
// reference instead.
//
// Element numbering: most segments number their positions in the order
// x12Names837p.ts's own array lists them (position 1, 2, 3…), matching the
// convention frontend/src/data/lookupRows835.ts already uses. Two kinds of
// segment need an explicit position list instead of plain counting, and get
// one below: a composite element whose sub-parts are cited as "CLM05-1",
// "SV101-2", and so on, and the NM1 role arrays, where x12Names837p.ts's own
// comment says position 6 (Name Prefix) is omitted from every array rather
// than held as a placeholder ("NM106…is omitted"). Positions this file
// cannot confirm from an explicit citation (CR1, CR2, CRC, HCP, MIA/MOA-style
// repeat slots elsewhere) are numbered by plain counting and may not match
// the X12 standard's own element number exactly; the name and citation text
// are still the public document's own.
import { CARC } from './x12Names835'
import * as P from './x12Names837p'
import type { LookupCodeRow, LookupElementRow } from '../docs/viewer/Lookup'

export const TRANSMITTAL_907 = 'https://www.cms.gov/files/document/cms-837-835-flatfiles-r907otn-2011.pdf'
export const CEDI_837P_GUIDE = 'https://www.ngscedi.com/technical-specifications'
export const X12_837P_EXAMPLES = 'https://x12.org/examples/005010x222'

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
// genuinely shared by more than one box (the CARC list, reused at both
// 2320 and 2430's CAS02).
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
const cedi = (pdf: string, printed: string) => `CEDI 837P CG, PDF p.${pdf} (printed p.${printed})`

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
const PATIENT = 'Loop 2010CA · patient'
const CLAIM = 'Loop 2300 · claim'
const REFERRING = 'Loop 2310A · referring provider'
const RENDERING = 'Loop 2310B · rendering provider'
const SERVICE_FACILITY = 'Loop 2310C · service facility location'
const SUPERVISING = 'Loop 2310D · supervising provider'
const AMB_PICKUP = 'Loop 2310E · ambulance pick-up location'
const AMB_DROPOFF = 'Loop 2310F · ambulance drop-off location'
const OTHER_SUBSCRIBER = 'Loop 2320/2330A · other subscriber (coordination of benefits)'
const OTHER_PAYER = 'Loop 2330B · other payer (coordination of benefits)'
const SERVICE_LINE = 'Loop 2400 · service line'
const DRUG = 'Loop 2410 · drug identification'
const LINE_RENDERING = 'Loop 2420A · rendering provider'
const LINE_PURCHASED = 'Loop 2420B · purchased service provider'
const LINE_SERVICE_FACILITY = 'Loop 2420C · service facility location'
const LINE_SUPERVISING = 'Loop 2420D · supervising provider'
const LINE_ORDERING = 'Loop 2420E · ordering provider'
const LINE_REFERRING = 'Loop 2420F · referring provider'
const LINE_AMB_PICKUP = 'Loop 2420G · ambulance pick-up location'
const LINE_AMB_DROPOFF = 'Loop 2420H · ambulance drop-off location'
const LINE_ADJUDICATION = 'Loop 2430 · line adjudication (coordination of benefits)'
const SUPPORTING_DOC = 'Loop 2440 · supporting documentation'
const TRAILER = 'Trailer'

// NM1 loop metadata: this loop's own box title (P.LOOP_NM1_NAME), the
// `where` scope above, and the citation for its NM1 segment, transcribed
// from x12Names837p.ts's own per-loop comment.
const NM1_LOOP: Partial<Record<string, { where: string; source: string }>> = {
  '1000A': { where: SUBMITTER, source: t907('135 (footer 3)') },
  '1000B': { where: RECEIVER, source: t907('136 (footer 4)') },
  '2010AA': { where: BILLING_PROVIDER, source: t907('138 (footer 6)') },
  '2010AB': { where: PAY_TO_ADDRESS, source: t907('140 (footer 8)') },
  '2010AC': { where: PAY_TO_PLAN, source: t907('141 (footer 9)') },
  '2010BA': { where: SUBSCRIBER, source: t907('144 (footer 12)') },
  '2010BB': { where: PAYER, source: t907('147 (footer 15)') },
  '2010CA': { where: PATIENT, source: t907('150 (footer 18)') },
  '2310A': { where: REFERRING, source: t907('184 (footer 52)') },
  '2310B': { where: RENDERING, source: t907('185 (footer 53)') },
  '2310C': { where: SERVICE_FACILITY, source: t907('186 (footer 54)') },
  '2310D': { where: SUPERVISING, source: t907('188 (footer 56)') },
  '2310E': { where: AMB_PICKUP, source: t907('189 (footer 57)') },
  '2310F': { where: AMB_DROPOFF, source: t907('189-190 (footer 57-58)') },
  '2330A': { where: OTHER_SUBSCRIBER, source: t907('195 (footer 63)') },
  '2330B': { where: OTHER_PAYER, source: t907('197 (footer 65)') },
  '2420A': { where: LINE_RENDERING, source: t907('225 (footer 93)') },
  '2420B': { where: LINE_PURCHASED, source: t907('226 (footer 94)') },
  '2420C': { where: LINE_SERVICE_FACILITY, source: t907('227 (footer 95)') },
  '2420D': { where: LINE_SUPERVISING, source: t907('229 (footer 97)') },
  '2420E': { where: LINE_ORDERING, source: t907('230 (footer 98)') },
  '2420F': { where: LINE_REFERRING, source: t907('231 (footer 99)') },
  '2420G': { where: LINE_AMB_PICKUP, source: t907('232 (footer 100)') },
  '2420H': { where: LINE_AMB_DROPOFF, source: t907('233 (footer 101)') },
}

// Every NM1 loop: NM101/NM102 (the same two positions everywhere) plus this
// loop's own NM103-109, skipping NM106 (Name Prefix) per x12Names837p.ts's
// own note that every 837P NM1 array omits it rather than placeholding it.
const NM1_REFS = ['01', '02', '03', '04', '05', '07', '08', '09']
function nm1LoopRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [loopId, names] of Object.entries(P.NM1_ROLE_EL_NAME)) {
    const meta = NM1_LOOP[loopId]
    if (!meta || !names) continue
    out.push(
      ...elementRowsWithRefs('NM1', meta.where, NM1_REFS, [P.NM101_EL_NAME, P.NM102_EL_NAME, ...names], meta.source, {
        segmentName: P.LOOP_NM1_NAME[loopId],
      }),
    )
  }
  return out
}

// N3/N4 (address) loops: same two/four positions everywhere; only the box
// title, scope, and citation change per loop.
const N3_LOOP: Partial<Record<string, { where: string; source: string }>> = {
  '2010AA': { where: BILLING_PROVIDER, source: t907('138 (footer 6)') },
  '2010AB': { where: PAY_TO_ADDRESS, source: t907('140-141 (footer 8-9)') },
  '2010AC': { where: PAY_TO_PLAN, source: t907('142 (footer 10)') },
  '2010BA': { where: SUBSCRIBER, source: t907('145 (footer 13)') },
  '2010BB': { where: PAYER, source: t907('148 (footer 16)') },
  '2010CA': { where: PATIENT, source: t907('150 (footer 18)') },
  '2310C': { where: SERVICE_FACILITY, source: t907('187 (footer 55)') },
  '2310E': { where: AMB_PICKUP, source: t907('189 (footer 57)') },
  '2310F': { where: AMB_DROPOFF, source: t907('190 (footer 58)') },
  '2330A': { where: OTHER_SUBSCRIBER, source: t907('196 (footer 64)') },
  '2330B': { where: OTHER_PAYER, source: t907('198 (footer 66)') },
  '2420C': { where: LINE_SERVICE_FACILITY, source: t907('228 (footer 96)') },
  '2420E': { where: LINE_ORDERING, source: t907('230 (footer 98)') },
  '2420G': { where: LINE_AMB_PICKUP, source: t907('233 (footer 101)') },
  '2420H': { where: LINE_AMB_DROPOFF, source: t907('234 (footer 102)') },
}
const N4_LOOP: Partial<Record<string, { where: string; source: string }>> = {
  '2010AA': { where: BILLING_PROVIDER, source: t907('138-139 (footer 6-7)') },
  '2010AB': { where: PAY_TO_ADDRESS, source: t907('141 (footer 9)') },
  '2010AC': { where: PAY_TO_PLAN, source: t907('142 (footer 10)') },
  '2010BA': { where: SUBSCRIBER, source: t907('145 (footer 13)') },
  '2010BB': { where: PAYER, source: t907('148 (footer 16)') },
  '2010CA': { where: PATIENT, source: t907('151 (footer 19-20)') },
  '2310C': { where: SERVICE_FACILITY, source: t907('187 (footer 55)') },
  '2310E': { where: AMB_PICKUP, source: t907('189 (footer 57)') },
  '2310F': { where: AMB_DROPOFF, source: t907('190 (footer 58)') },
  '2330A': { where: OTHER_SUBSCRIBER, source: t907('196 (footer 64)') },
  '2330B': { where: OTHER_PAYER, source: t907('198 (footer 66)') },
  '2420C': { where: LINE_SERVICE_FACILITY, source: t907('228 (footer 96)') },
  '2420E': { where: LINE_ORDERING, source: t907('230-231 (footer 98-99)') },
  '2420G': { where: LINE_AMB_PICKUP, source: t907('233-234 (footer 101-102)') },
  '2420H': { where: LINE_AMB_DROPOFF, source: t907('234 (footer 102)') },
}
function addressLoopRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [loopId, meta] of Object.entries(N3_LOOP)) {
    if (!meta) continue
    out.push(...elementRows('N3', meta.where, P.N3_EL_NAME, meta.source, { segmentName: P.LOOP_N3_NAME[loopId] }))
  }
  for (const [loopId, meta] of Object.entries(N4_LOOP)) {
    if (!meta) continue
    out.push(...elementRows('N4', meta.where, P.N4_EL_NAME, meta.source, { segmentName: P.LOOP_N4_NAME[loopId] }))
  }
  return out
}

// Secondary-identification REF loops: REF01/REF02 everywhere, the loop's
// own REF02 label overriding the generic one.
const REF_LOOP: Partial<Record<string, { where: string; source: string }>> = {
  '2310A': { where: REFERRING, source: t907('184-185 (footer 52-53)') },
  '2310B': { where: RENDERING, source: t907('186 (footer 54)') },
  '2310C': { where: SERVICE_FACILITY, source: t907('187-188 (footer 55-56)') },
  '2310D': { where: SUPERVISING, source: t907('188-189 (footer 56-57)') },
  '2330A': { where: OTHER_SUBSCRIBER, source: t907('196-197 (footer 64-65)') },
  '2420A': { where: LINE_RENDERING, source: t907('225-226 (footer 93-94)') },
  '2420B': { where: LINE_PURCHASED, source: t907('226-227 (footer 94-95)') },
  '2420C': { where: LINE_SERVICE_FACILITY, source: t907('228 (footer 96)') },
  '2420D': { where: LINE_SUPERVISING, source: t907('229-230 (footer 97-98)') },
  '2420E': { where: LINE_ORDERING, source: t907('231 (footer 99)') },
  '2420F': { where: LINE_REFERRING, source: t907('231-232 (footer 99-100)') },
}
function secondaryIdLoopRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [loopId, meta] of Object.entries(REF_LOOP)) {
    if (!meta) continue
    const names: [string, string] = [P.REF_SECONDARY_ID_EL_NAME[0], P.LOOP_REF02_EL_NAME[loopId] ?? P.REF_SECONDARY_ID_EL_NAME[1]]
    out.push(...elementRows('REF', meta.where, names, meta.source, { segmentName: P.LOOP_REF_NAME[loopId] }))
  }
  return out
}

// Claim-level DTP boxes: one DTP01/DTP02/DTP03 row per slug in
// P.CLAIM_DTP_NAME, each with its own citation.
const CLAIM_DTP_SOURCE: Record<string, string> = {
  onsetOfIllness: t907('155 (footer 23)'),
  initialTreatment: t907('155 (footer 23)'),
  lastSeen: t907('155 (footer 23)'),
  acuteManifestation: t907('155-156 (footer 23-24)'),
  accident: t907('156 (footer 24)'),
  lastMenstrualPeriod: t907('156 (footer 24)'),
  lastXRay: t907('156 (footer 24)'),
  hearingVisionPrescription: t907('156 (footer 24)'),
  disability: t907('156-157 (footer 24-25)'),
  lastWorked: t907('157 (footer 25)'),
  authorizedReturnToWork: t907('157 (footer 25)'),
  admission: t907('157 (footer 25)'),
  discharge: t907('157 (footer 25)'),
  assumedOrRelinquishedCare: t907('157-158 (footer 25-26)'),
  propertyCasualtyFirstContact: t907('158 (footer 26)'),
  repricerReceived: t907('158 (footer 26)'),
}
function claimDtpRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [slug, segmentName] of Object.entries(P.CLAIM_DTP_NAME)) {
    out.push(...elementRows('DTP', CLAIM, [P.DTP01_EL_NAME, P.DTP02_EL_NAME, P.CLAIM_DTP03_EL_NAME[slug]], CLAIM_DTP_SOURCE[slug], { segmentName }))
  }
  return out
}

// Claim-level and line-level REF boxes keyed by slug, in P.CLAIM_REF_NAME /
// P.SERVICE_REF_NAME, each a single REF01/REF02 pair.
const CLAIM_REF_SOURCE: Record<string, string> = {
  serviceAuthException: t907('159 (footer 27)'),
  medicareCrossoverIndicator: t907('160 (footer 28)'),
  mammographyCertification: t907('160 (footer 28)'),
  referral: t907('160-161 (footer 28-29)'),
  priorAuthorization: t907('161 (footer 29)'),
  payerClaimControlNumber: t907('161-162 (footer 29-30)'),
  clia: t907('162 (footer 30)'),
  repricedClaim: t907('162-163 (footer 30-31)'),
  adjustedRepricedClaim: t907('163 (footer 31)'),
  investigationalDeviceExemption: t907('163-164 (footer 31-32)'),
  claimIdentifierForTransmissionIntermediaries: t907('164 (footer 32)'),
  medicalRecordNumber: t907('164-165 (footer 32-33)'),
  demonstrationProject: t907('165 (footer 33)'),
  carePlanOversight: t907('165 (footer 33)'),
}
function claimRefRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [slug, segmentName] of Object.entries(P.CLAIM_REF_NAME)) {
    out.push(...elementRows('REF', CLAIM, [undefined, P.CLAIM_REF02_EL_NAME[slug]], CLAIM_REF_SOURCE[slug], { segmentName }))
  }
  return out
}
const SERVICE_REF_SOURCE: Record<string, string> = {
  repricedLineItem: t907('215-216 (footer 83-84)'),
  adjustedRepricedLineItem: t907('216 (footer 84)'),
  priorAuthorization: t907('216 (footer 84)'),
  lineItemControl: t907('217 (footer 85)'),
  mammographyCertification: t907('217 (footer 85)'),
  clia: t907('218 (footer 86)'),
  referringClia: t907('218 (footer 86)'),
  immunizationBatch: t907('219 (footer 87)'),
  referral: t907('219 (footer 87)'),
}
function serviceRefRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [slug, segmentName] of Object.entries(P.SERVICE_REF_NAME)) {
    out.push(...elementRows('REF', SERVICE_LINE, [undefined, P.SERVICE_REF02_EL_NAME[slug]], SERVICE_REF_SOURCE[slug], { segmentName }))
  }
  return out
}

// Line-level DTP boxes, keyed by slug in P.SERVICE_DTP_NAME.
const SERVICE_DTP_SOURCE: Record<string, string> = {
  serviceDate: t907('211-212 (footer 79-80)'),
  prescription: t907('212 (footer 80)'),
  certificationRevision: t907('212 (footer 80)'),
  beginTherapy: t907('212 (footer 80)'),
  lastCertification: t907('212 (footer 80)'),
  lastTreatmentOrTherapy: t907('212-213 (footer 80-81)'),
  testPerformed: t907('213 (footer 81)'),
  shipped: t907('213 (footer 81)'),
  lastXRay: t907('213 (footer 81)'),
  initialTreatment: t907('213 (footer 81)'),
}
function serviceDtpRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [slug, segmentName] of Object.entries(P.SERVICE_DTP_NAME)) {
    out.push(...elementRows('DTP', SERVICE_LINE, [P.DTP01_EL_NAME, P.DTP02_EL_NAME, P.SERVICE_DTP03_EL_NAME[slug]], SERVICE_DTP_SOURCE[slug], { segmentName }))
  }
  return out
}

// Other-payer REF boxes, loop 2330B, keyed by slug in P.OTHER_PAYER_REF_NAME.
const OTHER_PAYER_REF_SOURCE: Record<string, string> = {
  secondaryIdentification: t907('198 (footer 66)'),
  priorAuthorization: t907('198-199 (footer 66-67)'),
  referral: t907('199 (footer 67)'),
  claimAdjustmentIndicator: t907('199-200 (footer 67-68)'),
  claimControlNumber: t907('200 (footer 68)'),
}
function otherPayerRefRows(): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (const [slug, segmentName] of Object.entries(P.OTHER_PAYER_REF_NAME)) {
    out.push(...elementRows('REF', OTHER_PAYER, [undefined, P.OTHER_PAYER_REF02_EL_NAME[slug]], OTHER_PAYER_REF_SOURCE[slug], { segmentName }))
  }
  return out
}

// 2330C-G other-payer provider loops: REF only, not used by Medicare (see
// P.NOT_USED_BY_MEDICARE), kept for completeness with that noted in `where`.
function otherPayerProviderRows(): LookupElementRow[] {
  const sources: Record<string, string> = {
    '2330C': t907('200-201 (footer 68-69)'),
    '2330D': t907('201-202 (footer 69-70)'),
    '2330E': t907('202-203 (footer 70-71)'),
    '2330F': t907('203-204 (footer 71-72)'),
    '2330G': t907('204-205 (footer 72-73)'),
  }
  const out: LookupElementRow[] = []
  for (const [loopId, segmentName] of Object.entries(P.OTHER_PAYER_ROLE_REF_NAME)) {
    out.push(...elementRows('REF', `${OTHER_PAYER} (not used by Medicare)`, P.REF_SECONDARY_ID_EL_NAME, sources[loopId], { segmentName }))
  }
  return out
}

// 2300 HI: three kinds of composite (diagnosis, anesthesia-related
// procedure, condition information), each up to 12 repeats of a 9-part
// composite (HI0n-1 through HI0n-9; Medicare uses the first two, and the flat
// file names all nine).
function hiRows(label: string, elName: string[], source: string): LookupElementRow[] {
  const out: LookupElementRow[] = []
  for (let n = 1; n <= 12; n++) {
    const refs = elName.map((_, k) => `${pad(n)}-${k + 1}`)
    out.push(...elementRowsWithRefs('HI', CLAIM, refs, elName, source, { segmentName: `${label}, repeat ${n} of 12` }))
  }
  return out
}

// 2320/2430 CAS: a group code, then up to five reason/amount/quantity
// triplets (x12Names837p.ts: "repeats up to 5 times" at both loops) —
// another payer's adjustments, relayed forward for coordination of
// benefits, not a MAC adjustment (that only appears on the 835).
function casRows(where: string, source: string, segmentName: string): LookupElementRow[] {
  const refs = ['01']
  const names: (string | undefined)[] = [P.CAS01_EL_NAME]
  for (let i = 0; i < 5; i++) {
    const n = 2 + i * 3
    refs.push(pad(n), pad(n + 1), pad(n + 2))
    names.push(P.CAS_REASON_EL_NAME, P.CAS_AMOUNT_EL_NAME, P.CAS_QTY_EL_NAME)
  }
  return elementRowsWithRefs('CAS', where, refs, names, source, { segmentName })
}

export const ELEMENTS_837P: LookupElementRow[] = [
  // ---- Envelope ----
  ...elementRows('ISA', ENVELOPE, P.ISA_EL_NAME, t907('133 (footer 1)'), { segmentName: P.SEGMENT_NAME.ISA }),
  ...elementRows('GS', ENVELOPE, P.GS_EL_NAME, t907('133-134 (footer 1-2)'), { segmentName: P.SEGMENT_NAME.GS }),
  ...elementRows('ST', ENVELOPE, P.ST_EL_NAME, t907('134 (footer 2)'), { segmentName: P.SEGMENT_NAME.ST }),
  ...elementRows('BHT', ENVELOPE, P.BHT_EL_NAME, t907('134-135 (footer 2-3)'), { segmentName: P.SEGMENT_NAME.BHT }),

  // ---- 1000A/1000B submitter and receiver ----
  ...elementRows('PER', SUBMITTER, P.SUBMITTER_PER_EL_NAME, t907('135-136 (footer 3-4)'), { segmentName: 'Submitter EDI Contact Information' }),

  // ---- 2000A billing provider hierarchy, 2010AA/AB/AC ----
  ...elementRows('HL', BILLING_HL, P.HL_EL_NAME, t907('136 (footer 4)'), { segmentName: P.LOOP_HL_NAME['2000A'] }),
  ...elementRows('PRV', BILLING_HL, P.PRV_EL_NAME, t907('137 (footer 5)'), { segmentName: P.LOOP_PRV_NAME['2000A'] }),
  ...elementRows('CUR', BILLING_HL, P.CUR_EL_NAME, t907('137 (footer 5)'), { segmentName: P.CUR_NAME }),
  ...elementRows('REF', BILLING_PROVIDER, P.BILLING_PROVIDER_REF_EL_NAME, t907('139 (footer 7)'), { segmentName: P.BILLING_PROVIDER_REF_NAME }),
  ...elementRows('REF', BILLING_PROVIDER, P.REF_SECONDARY_ID_EL_NAME, t907('139 (footer 7)'), { segmentName: P.BILLING_PROVIDER_UPIN_REF_NAME }),
  ...elementRows('PER', BILLING_PROVIDER, P.BILLING_PROVIDER_PER_EL_NAME, t907('139-140 (footer 7-8)'), { segmentName: P.BILLING_PROVIDER_PER_NAME }),
  ...elementRows('REF', PAY_TO_PLAN, P.REF_SECONDARY_ID_EL_NAME, t907('142 (footer 10)'), { segmentName: P.PAY_TO_PLAN_SECONDARY_REF_NAME }),
  ...elementRows('REF', PAY_TO_PLAN, P.REF_SECONDARY_ID_EL_NAME, t907('143 (footer 11)'), { segmentName: P.PAY_TO_PLAN_TAX_REF_NAME }),

  // ---- NM1/N3/N4/REF loops shared by every party ----
  ...nm1LoopRows(),
  ...addressLoopRows(),
  ...secondaryIdLoopRows(),

  // ---- 2000B/C, 2010BA/BB/CA subscriber and patient ----
  ...elementRows('SBR', SUBSCRIBER_HL, P.SBR_EL_NAME, t907('143-144 (footer 11-12)'), { segmentName: 'Subscriber Information' }),
  ...elementRows('SBR', OTHER_SUBSCRIBER, P.OTHER_SUBSCRIBER_SBR_EL_NAME, t907('191-192 (footer 59-60)'), { segmentName: 'Other Subscriber Information' }),
  ...elementRows('PAT', SUBSCRIBER_HL, P.PAT_EL_NAME, t907('144 (footer 12)'), { segmentName: P.PAT_NAME }),
  ...elementRows('DMG', SUBSCRIBER, P.DMG_EL_NAME, t907('145-146 (footer 13-14)'), { segmentName: P.DMG_NAME }),
  ...elementRows('REF', SUBSCRIBER, P.REF_SECONDARY_ID_EL_NAME, t907('146 (footer 14)'), { segmentName: P.SUBSCRIBER_PC_REF_NAME }),
  ...elementRows('REF', SUBSCRIBER, P.REF_SECONDARY_ID_EL_NAME, t907('146 (footer 14)'), { segmentName: P.SUBSCRIBER_SECONDARY_REF_NAME }),
  ...elementRows('REF', PAYER, P.REF_SECONDARY_ID_EL_NAME, t907('148 (footer 16)'), { segmentName: P.PAYER_SECONDARY_REF_NAME }),
  ...elementRows('REF', PAYER, P.REF_SECONDARY_ID_EL_NAME, t907('148-149 (footer 16-17)'), { segmentName: P.PAYER_LOOP_BILLING_PROVIDER_REF_NAME }),
  ...elementRows('REF', PATIENT, P.REF_SECONDARY_ID_EL_NAME, t907('152-153 (footer 20-21)'), { segmentName: P.PATIENT_PC_REF_NAME }),
  // P.SUBSCRIBER_PC_PER_NAME and P.PATIENT_PC_PER_NAME (property-and-casualty
  // contact boxes) have no element-position array in x12Names837p.ts, so no
  // row is built for them here rather than guessing their width.

  // ---- 2300 claim ----
  ...elementRowsWithRefs(
    'CLM',
    CLAIM,
    ['01', '02', '03', '04', '05-1', '05-2', '05-3', '06', '07', '08', '09', '10', '11-1', '11-2', '11-4', '11-5', '12', '13', '14', '15', '16', '17', '18', '19', '20'],
    P.CLM_EL_NAME,
    t907('155-157 (footer 23-25)'),
    { segmentName: P.CLM_NAME },
  ),
  ...claimDtpRows(),
  ...claimRefRows(),
  ...casRows(OTHER_SUBSCRIBER, t907('192 (footer 60)'), P.CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME),
  ...elementRows('PWK', CLAIM, P.PWK_EL_NAME, t907('158 (footer 26)'), { segmentName: P.CLAIM_PWK_NAME }),
  ...elementRows('PWK', SERVICE_LINE, P.PWK_EL_NAME, t907('209 (footer 77)'), { segmentName: P.SERVICE_PWK_NAME }),
  ...elementRows('CN1', CLAIM, P.CN1_EL_NAME, t907('159 (footer 27)'), { segmentName: P.CN1_NAME }),
  ...elementRows('AMT', CLAIM, P.CLAIM_AMT_EL_NAME, t907('159 (footer 27)'), { segmentName: P.CLAIM_AMT_NAME }),
  ...elementRows('K3', CLAIM, [P.K3_EL_NAME], t907('166 (footer 34)'), { segmentName: P.CLAIM_K3_NAME }),
  ...elementRows('K3', SERVICE_LINE, [P.K3_EL_NAME], t907('221 (footer 88)'), { segmentName: P.SERVICE_K3_NAME }),
  ...elementRows('NTE', CLAIM, P.NTE_EL_NAME, t907('166 (footer 34)'), { segmentName: P.CLAIM_NTE_NAME }),
  ...elementRows('CR1', CLAIM, P.CR1_EL_NAME, t907('166-167 (footer 34-35)'), { segmentName: P.CLAIM_CR1_NAME }),
  ...elementRows('CR1', SERVICE_LINE, P.CR1_EL_NAME, t907('210 (footer 78)'), { segmentName: P.SERVICE_CR1_NAME }),
  ...elementRows('CR2', CLAIM, P.CR2_EL_NAME, t907('167 (footer 35)'), { segmentName: P.CR2_NAME }),
  ...elementRows('CRC', CLAIM, P.CRC_EL_NAME, t907('168 (footer 36)'), { segmentName: P.CLAIM_CRC_NAME.ambulanceCertification }),
  ...elementRows('CRC', CLAIM, P.CRC_EL_NAME, t907('168 (footer 36)'), { segmentName: P.CLAIM_CRC_NAME.visionCondition }),
  ...elementRows('CRC', CLAIM, P.CRC_EL_NAME, t907('168 (footer 36)'), { segmentName: P.CLAIM_CRC_NAME.homebound }),
  ...elementRows('CRC', CLAIM, P.CRC_EL_NAME, t907('169 (footer 37)'), { segmentName: P.CLAIM_CRC_NAME.epsdtReferral }),
  ...hiRows(P.HI_DIAGNOSIS_NAME, P.HI_DIAGNOSIS_EL_NAME, t907('169-174 (footer 37-42)')),
  ...hiRows(P.HI_ANESTHESIA_NAME, P.HI_ANESTHESIA_EL_NAME, t907('174-178 (footer 42-46)')),
  ...hiRows(P.HI_CONDITION_NAME, P.HI_CONDITION_EL_NAME, t907('178-183 (footer 46-51)')),
  ...elementRows('HCP', CLAIM, P.HCP_EL_NAME, t907('181-182 (footer 49-50)'), { segmentName: P.HCP_NAME }),
  ...elementRows('HCP', SERVICE_LINE, P.HCP_EL_NAME, t907('216-217 (footer 84-85)'), { segmentName: P.LINE_HCP_NAME }),

  // ---- 2320/2330 other subscriber and payer (COB) ----
  ...elementRows('AMT', OTHER_SUBSCRIBER, [P.AMT01_EL_NAME, P.AMT02_EL_NAME_BY_BOX.cobPayerPaid], t907('194 (footer 62)'), { segmentName: P.COB_PAYER_PAID_AMT_NAME }),
  ...elementRows('AMT', OTHER_SUBSCRIBER, [P.AMT01_EL_NAME, P.AMT02_EL_NAME_BY_BOX.cobNonCovered], t907('194 (footer 62)'), { segmentName: P.COB_NON_COVERED_AMT_NAME }),
  ...elementRows('AMT', OTHER_SUBSCRIBER, [P.AMT01_EL_NAME, P.AMT02_EL_NAME_BY_BOX.remainingPatientLiability], t907('194 (footer 62)'), { segmentName: P.OTHER_SUBSCRIBER_REMAINING_LIABILITY_AMT_NAME }),
  ...elementRows('AMT', LINE_ADJUDICATION, [P.AMT01_EL_NAME, P.AMT02_EL_NAME_BY_BOX.remainingPatientLiability], t907('237 (footer 105)'), { segmentName: P.LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME }),
  ...elementRows('OI', OTHER_SUBSCRIBER, P.OI_EL_NAME, t907('194-195 (footer 62-63)'), { segmentName: P.OI_NAME }),
  ...elementRows('MOA', OTHER_SUBSCRIBER, P.MOA_EL_NAME, t907('195 (footer 63)'), { segmentName: P.MOA_NAME }),
  ...elementRows('DTP', OTHER_PAYER, [P.DTP01_EL_NAME, P.DTP02_EL_NAME, P.OTHER_PAYER_DTP03_EL_NAME], t907('198 (footer 66)'), { segmentName: P.OTHER_PAYER_DTP_NAME }),
  ...elementRows('DTP', LINE_ADJUDICATION, [P.DTP01_EL_NAME, P.DTP02_EL_NAME, P.LINE_ADJUDICATION_DTP03_EL_NAME], t907('237 (footer 105)'), { segmentName: P.LINE_ADJUDICATION_DTP_NAME }),
  ...otherPayerRefRows(),
  ...otherPayerProviderRows(),

  // ---- 2400 service line ----
  ...elementRows('LX', SERVICE_LINE, [P.LX01_EL_NAME], t907('205 (footer 73)'), { segmentName: P.LX_NAME }),
  ...elementRowsWithRefs('SV1', SERVICE_LINE, ['01-1', '01-2', '01-3', '01-4', '01-5', '01-6', '01-7', '02', '03', '04', '05', '07-1', '07-2', '07-3', '07-4', '09', '11', '12', '19'], P.SV1_EL_NAME, t907('205-208 (footer 73-76)'), { segmentName: P.SV1_NAME }),
  ...elementRowsWithRefs('SV5', SERVICE_LINE, ['01-1', '01-2', '01-3', '01-4', '01-5', '01-6', '01-7', '02', '03', '04', '05', '06'], P.SV5_EL_NAME, t907('208-209 (footer 76-77)'), { segmentName: P.SV5_NAME }),
  ...elementRows('CR3', SERVICE_LINE, P.CR3_EL_NAME, t907('210 (footer 78)'), { segmentName: P.CR3_NAME }),
  ...elementRows('CRC', SERVICE_LINE, P.CRC_EL_NAME, t907('210-211 (footer 78-79)'), { segmentName: P.SERVICE_CRC_NAME.ambulanceCertification }),
  ...elementRows('CRC', SERVICE_LINE, P.CRC_EL_NAME, t907('211 (footer 79)'), { segmentName: P.SERVICE_CRC_NAME.hospiceEmployeeIndicator }),
  ...elementRows('CRC', SERVICE_LINE, P.CRC_EL_NAME, t907('211 (footer 79)'), { segmentName: P.SERVICE_CRC_NAME.dmeConditionIndicator }),
  ...serviceDtpRows(),
  ...elementRows('QTY', SERVICE_LINE, [P.QTY01_EL_NAME, P.QTY02_EL_NAME_BY_BOX.ambulancePatientCount], t907('213-214 (footer 81-82)'), { segmentName: P.AMBULANCE_PATIENT_COUNT_QTY_NAME }),
  ...elementRows('QTY', SERVICE_LINE, [P.QTY01_EL_NAME, P.QTY02_EL_NAME_BY_BOX.obstetricAdditionalUnits], t907('214 (footer 82)'), { segmentName: P.OBSTETRIC_ANESTHESIA_QTY_NAME }),
  ...elementRows('MEA', SERVICE_LINE, P.MEA_EL_NAME, t907('214-215 (footer 82-83)'), { segmentName: P.MEA_NAME }),
  ...serviceRefRows(),
  ...elementRows('AMT', SERVICE_LINE, [P.AMT01_EL_NAME, P.SERVICE_AMT02_EL_NAME_BY_BOX.salesTax], t907('219-220 (footer 87-88)'), { segmentName: P.SALES_TAX_AMT_NAME }),
  ...elementRows('AMT', SERVICE_LINE, [P.AMT01_EL_NAME, P.SERVICE_AMT02_EL_NAME_BY_BOX.postageClaimed], t907('220 (footer 88)'), { segmentName: P.POSTAGE_CLAIMED_AMT_NAME }),
  ...elementRows('NTE', SERVICE_LINE, P.NTE_EL_NAME, t907('220 (footer 88)'), { segmentName: P.LINE_NTE_NAME }),
  ...elementRows('NTE', SERVICE_LINE, P.NTE_EL_NAME, t907('220 (footer 88)'), { segmentName: P.THIRD_PARTY_NTE_NAME }),
  ...elementRows('PS1', SERVICE_LINE, P.PS1_EL_NAME, t907('220-221 (footer 88-89)'), { segmentName: P.PS1_NAME }),

  // ---- 2410 drug identification ----
  ...elementRowsWithRefs('LIN', DRUG, ['02', '03'], P.LIN_EL_NAME, t907('222 (footer 90)'), { segmentName: P.LIN_NAME }),
  ...elementRowsWithRefs('CTP', DRUG, ['04', '05-1'], P.CTP_EL_NAME, t907('223-224 (footer 91-92)'), { segmentName: P.CTP_NAME }),
  ...elementRows('REF', DRUG, [undefined, P.DRUG_ASSOCIATION_REF02_EL_NAME], t907('224 (footer 92)'), { segmentName: P.DRUG_ASSOCIATION_REF_NAME }),

  // ---- 2420E/2310C PER contacts (not loop-keyed in x12Names837p.ts) ----
  ...elementRows('PER', LINE_ORDERING, P.ORDERING_PROVIDER_PER_EL_NAME, t907('231 (footer 99)'), { segmentName: P.ORDERING_PROVIDER_PER_NAME }),
  ...elementRows('PER', SERVICE_FACILITY, P.SERVICE_FACILITY_PER_EL_NAME, t907('188 (footer 56)'), { segmentName: P.SERVICE_FACILITY_PER_NAME }),

  // ---- 2430 line adjudication ----
  ...elementRowsWithRefs('SVD', LINE_ADJUDICATION, ['01', '02', '03-1', '03-2', '03-3', '03-4', '03-5', '03-6', '03-7', '05', '06'], P.SVD_EL_NAME, t907('235 (footer 103)'), { segmentName: P.SVD_NAME }),
  ...casRows(LINE_ADJUDICATION, t907('235 (footer 103)'), P.LINE_ADJUSTMENT_CAS_NAME),

  // ---- 2440 supporting documentation ----
  ...elementRows('LQ', SUPPORTING_DOC, P.LQ_EL_NAME, t907('238 (footer 106)'), { segmentName: P.LQ_NAME }),
  ...elementRows('FRM', SUPPORTING_DOC, P.FRM_EL_NAME, t907('238 (footer 106)'), { segmentName: P.FRM_NAME }),

  // ---- Trailer ----
  ...elementRows('SE', TRAILER, P.SE_EL_NAME, t907('238 (footer 106)'), { segmentName: P.SEGMENT_NAME.SE }),
  ...elementRows('GE', TRAILER, P.GE_EL_NAME, t907('238 (footer 106)'), { segmentName: P.SEGMENT_NAME.GE }),
  ...elementRows('IEA', TRAILER, P.IEA_EL_NAME, t907('239 (footer 107)'), { segmentName: P.SEGMENT_NAME.IEA }),
]

export const CODES_837P: LookupCodeRow[] = [
  ...codesFromRecord('BHT', 'BHT02', P.BHT02_NAME, ENVELOPE, cedi('30', '22'), P.SEGMENT_NAME.BHT),
  ...codesFromRecord('BHT', 'BHT06', P.BHT06_NAME, ENVELOPE, cedi('30', '22'), P.SEGMENT_NAME.BHT),
  ...codesFromRecord('SBR', 'SBR01', P.SBR01_NAME, SUBSCRIBER_HL, cedi('33', '25'), 'Subscriber Information'),
  ...codesFromRecord('SBR', 'SBR02', P.SBR02_NAME, SUBSCRIBER_HL, `${cedi('33', '25')}; the only value Medicare's own subscriber loop sends`, 'Subscriber Information'),
  ...codesFromRecord('SBR', 'SBR09', P.SBR09_NAME, SUBSCRIBER_HL, `${cedi('33', '25')}; loop 2320's SBR09 cannot equal MA or MB (${cedi('40', '32')}, Table 27)`, 'Subscriber Information'),
  codeRow('GS', 'GS01', 'HC', P.TRANSACTION_TYPE_NAME.HC, ENVELOPE, `Transmittal 907 (GS01 value shown in the CEDI worked example, GS*HC*…); ${cedi('50', '42')}, Figure 3`, undefined, P.SEGMENT_NAME.GS),
  codeRow('ST', 'ST01', '837', P.TRANSACTION_TYPE_NAME['837'], ENVELOPE, `${cedi('50', '42')}, Figure 3, ST*837*…`, undefined, P.SEGMENT_NAME.ST),
  // NM101/NM108 codes below are left broad (no segmentName) where the
  // citation itself says the code applies at more than one loop; each
  // still attaches only to the NM1 occurrences that actually carry that
  // position.
  ...codesFromRecord('NM1', 'NM101', P.NM101_CONFIRMED_NAME, 'Every NM1 instance', `${cedi('50', '42')}, Figure 3`),
  codeRow('NM1', 'NM108', '46', P.NM108_CONFIRMED_NAME['46'], 'Every NM1 instance (submitter/receiver)', `${cedi('50', '42')}, Figure 3`),
  codeRow('NM1', 'NM108', 'XX', P.NM108_CONFIRMED_NAME.XX, 'Every NM1 instance (billing/ordering provider)', `${cedi('50', '42')}, Figure 3`),
  codeRow('NM1', 'NM108', 'MI', P.NM108_CONFIRMED_NAME.MI, 'Every NM1 instance (subscriber)', `${cedi('50', '42')}, Figure 3; ${cedi('33', '25')}, Table 15`, undefined, P.LOOP_NM1_NAME['2010BA']),
  codeRow('NM1', 'NM108', 'PI', P.NM108_CONFIRMED_NAME.PI, 'Every NM1 instance (payer)', `${cedi('35', '27')}, Table 17`, undefined, P.LOOP_NM1_NAME['2010BB']),
  codeRow('CLM', 'CLM05-3', '1', P.CLM05_3_NAME['1'], CLAIM, cedi('36', '28'), undefined, P.CLM_NAME),
  codeRow('SV1', 'SV101-1', 'HC', P.SV101_1_NAME.HC, SERVICE_LINE, cedi('42', '34'), undefined, P.SV1_NAME),
  codeRow('SVD', 'SVD03-1', 'HC', P.SVD03_1_NAME.HC, LINE_ADJUDICATION, cedi('49', '41'), undefined, P.SVD_NAME),
  ...valuesOnly('PWK', 'PWK02', P.PWK02_MEDICARE_CODES, `${CLAIM} · ${SERVICE_LINE}`, cedi('38', '30')),
  codeRow('AMT', 'AMT01', 'D', P.AMT01_COB_NAME.D, OTHER_SUBSCRIBER, cedi('40', '32'), undefined, P.COB_PAYER_PAID_AMT_NAME),
  ...codesFromRecord('LIN', 'LIN02', P.LIN02_NAME, DRUG, cedi('46', '38'), P.LIN_NAME),
  ...codesFromRecord('CTP', 'CTP05-1', P.CTP05_1_NAME, DRUG, cedi('47', '39'), P.CTP_NAME),
  ...codesFromRecord('REF', 'REF01', P.DRUG_ASSOCIATION_REF01_NAME, DRUG, cedi('47', '39'), P.DRUG_ASSOCIATION_REF_NAME),
  ...codesFromRecord('CAS', 'CAS01', P.CAS_GROUP_NAME, `${OTHER_SUBSCRIBER} · ${LINE_ADJUDICATION}`, 'https://x12.org/codes/claim-adjustment-group-codes'),
  ...Object.entries(CARC).map(([value, text]) => ({
    ...codeRow('CAS', 'CAS02', value, text.short, `${OTHER_SUBSCRIBER} · ${LINE_ADJUDICATION} (reusing the 835's own CARC list, per x12Names837p.ts's note)`, 'https://x12.org/codes/claim-adjustment-reason-codes (text embedded verbatim from backend/pkg/carc.json)'),
    full: text.full,
  })),
]
