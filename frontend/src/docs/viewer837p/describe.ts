// Turns a parsed 837P segment or box opening/closing into the row of
// display cells the page renders — the 837P's own version of
// frontend/src/docs/viewer835/describe.ts; read that file's header first for
// the shared conventions (a cell's code, value, and name; qualifier-driven
// vs fixed row names; the cell-building helpers now shared in
// frontend/src/docs/viewer/display.ts). Every public name comes from
// frontend/src/data/x12Names837p.ts, built only from Transmittal 907 and the
// CEDI 837P Companion Guide; a coded value with no entry there renders as
// its raw code alone, never invented.
//
// A positional note this file leans on throughout: most of x12Names837p.ts's
// element-name arrays name every position directly (array[i] is element
// i+1, matching the raw file 1:1, including a "Not Used" position that
// still gets a name and simply arrives empty — confirmed against the 24
// example files for CLM03/04, CR1, CRC, PAT, SBR01-05, DMG, PRV, MEA, PS1,
// CTP, LIN). A few arrays instead *compact away* a Not Used position's name
// entirely, so the array is shorter than the real element count and a
// later array index lands on a later real element than its own position
// would suggest — confirmed against the files for CLM (the CLM05 composite
// adds 2 extra name slots for 1 raw element, then CLM06+ runs name-index =
// raw-index + 2; the CLM11 composite adds 3 more, so CLM12+ runs
// name-index = raw-index + 5), SBR (SBR06-08 compacted away, so SBR09's
// name sits at array index 5, not 8), SVD (SVD04 compacted away, so SVD05/
// SVD06's names sit at array index 9/10, not 4/5), OI (OI01/02/05 compacted
// away, so OI03/04/06's names sit at array index 0/1/2, not 2/3/5), and CR2
// (CR201-07/209/212 compacted away entirely — only CR208, CR210, CR211 are
// named, at array index 0/1/2). Each case below is commented with the real
// element index it reads.
import { CARC, RARC } from '../../data/x12Names835'
import {
  AMBULANCE_DROPOFF_NM1_EL_NAME,
  AMBULANCE_PATIENT_COUNT_QTY_NAME,
  AMBULANCE_PICKUP_NM1_EL_NAME,
  AMT01_COB_NAME,
  AMT01_EL_NAME,
  AMT02_EL_NAME_BY_BOX,
  BHT02_NAME,
  BHT06_NAME,
  BHT_EL_NAME,
  BILLING_PROVIDER_PER_EL_NAME,
  BILLING_PROVIDER_PER_NAME,
  BILLING_PROVIDER_REF_EL_NAME,
  BILLING_PROVIDER_REF_NAME,
  CAS01_EL_NAME,
  CAS_AMOUNT_EL_NAME,
  CAS_GROUP_NAME,
  CAS_QTY_EL_NAME,
  CAS_REASON_EL_NAME,
  CLAIM_AMT_EL_NAME,
  CLAIM_AMT_NAME,
  CLAIM_CRC_NAME,
  CLAIM_CR1_NAME,
  CLAIM_DTP03_EL_NAME,
  CLAIM_DTP_NAME,
  CLAIM_K3_NAME,
  CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME,
  CLAIM_NTE_NAME,
  CLAIM_PWK_NAME,
  CLAIM_REF02_EL_NAME,
  CLAIM_REF_NAME,
  CLM05_3_NAME,
  CLM_EL_NAME,
  CLM_NAME,
  CN1_EL_NAME,
  CN1_NAME,
  CR1_EL_NAME,
  CR2_EL_NAME,
  CR2_NAME,
  CR3_EL_NAME,
  CR3_NAME,
  CRC_EL_NAME,
  CTP05_1_NAME,
  CTP_EL_NAME,
  CTP_NAME,
  CUR_EL_NAME,
  CUR_NAME,
  DMG_EL_NAME,
  DMG_NAME,
  DRUG_ASSOCIATION_REF01_NAME,
  DRUG_ASSOCIATION_REF02_EL_NAME,
  DRUG_ASSOCIATION_REF_NAME,
  DTP01_EL_NAME,
  DTP02_EL_NAME,
  FRM_EL_NAME,
  FRM_NAME,
  GE_EL_NAME,
  GS_EL_NAME,
  HCP_EL_NAME,
  HCP_NAME,
  HI_CONDITION_NAME,
  HI_DIAGNOSIS_NAME,
  HI_NOT_USED_EL_NAME,
  HL_EL_NAME,
  IEA_EL_NAME,
  ISA_EL_NAME,
  LIN02_NAME,
  LIN_EL_NAME,
  LIN_NAME,
  LINE_ADJUDICATION_DTP03_EL_NAME,
  LINE_ADJUDICATION_DTP_NAME,
  LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME,
  LINE_HCP_NAME,
  LINE_NTE_NAME,
  LOOP_N3_NAME,
  LOOP_N4_NAME,
  LOOP_NM1_NAME,
  LOOP_REF02_EL_NAME,
  LOOP_REF_NAME,
  LQ_EL_NAME,
  LQ_NAME,
  LX01_EL_NAME,
  LX_NAME,
  MEA_EL_NAME,
  MEA_NAME,
  MOA_EL_NAME,
  MOA_NAME,
  N3_EL_NAME,
  N4_EL_NAME,
  NM101_CONFIRMED_NAME,
  NM101_EL_NAME,
  NM102_EL_NAME,
  NM108_CONFIRMED_NAME,
  NM1_ROLE_EL_NAME,
  OBSTETRIC_ANESTHESIA_QTY_NAME,
  OI_EL_NAME,
  OI_NAME,
  ORDERING_PROVIDER_PER_EL_NAME,
  ORDERING_PROVIDER_PER_NAME,
  OTHER_PAYER_DTP03_EL_NAME,
  OTHER_PAYER_DTP_NAME,
  OTHER_PAYER_REF02_EL_NAME,
  OTHER_PAYER_REF_NAME,
  OTHER_SUBSCRIBER_SBR_EL_NAME,
  PAT_EL_NAME,
  PAT_NAME,
  POSTAGE_CLAIMED_AMT_NAME,
  PRV_EL_NAME,
  PS1_EL_NAME,
  PS1_NAME,
  PWK_EL_NAME,
  QTY01_EL_NAME,
  QTY02_EL_NAME_BY_BOX,
  SALES_TAX_AMT_NAME,
  SBR01_NAME,
  SBR02_NAME,
  SBR09_NAME,
  SBR_EL_NAME,
  SE_EL_NAME,
  SEGMENT_NAME,
  SERVICE_AMT02_EL_NAME_BY_BOX,
  SERVICE_CR1_NAME,
  SERVICE_CRC_NAME,
  SERVICE_DTP03_EL_NAME,
  SERVICE_DTP_NAME,
  SERVICE_K3_NAME,
  SERVICE_PWK_NAME,
  SERVICE_REF02_EL_NAME,
  SERVICE_REF_NAME,
  ST_EL_NAME,
  SUBMITTER_PER_EL_NAME,
  SUBSCRIBER_PC_PER_NAME,
  SUBSCRIBER_PC_REF_NAME,
  SV101_1_NAME,
  SV1_EL_NAME,
  SV1_NAME,
  SV5_EL_NAME,
  SV5_NAME,
  SVD03_1_NAME,
  SVD_EL_NAME,
  SVD_NAME,
  THIRD_PARTY_NTE_NAME,
  TRANSACTION_TYPE_NAME,
} from '../../data/x12Names837p'
import { cell, label, lcell, pad, ref, value, vcell, type Cell, type Part, type SegmentDisplay } from '../viewer/display'
import type { Segment } from '../viewer/tree'
import type { Box } from './parse'
import type { BoxType } from './rules'

// ---------------------------------------------------------------------------
// Qualifier → slug maps, derived from the 24 example files under
// frontend/src/data/viewer837p/files/ (the file named at each entry below;
// x12Names837p.ts documents the slug's public name but, per its own stated
// scope, not which literal code selects it — see that file's header).
// ---------------------------------------------------------------------------

const CLAIM_DTP_SLUG: Record<string, keyof typeof CLAIM_DTP_NAME> = {
  '431': 'onsetOfIllness', // multiple-diagnoses.837
  '454': 'initialTreatment', // multiple-diagnoses.837
  '304': 'lastSeen', // multiple-diagnoses.837
  '453': 'acuteManifestation', // multiple-diagnoses.837
  '439': 'accident', // property-and-casualty-claim.837
  '484': 'lastMenstrualPeriod', // anesthesia-minutes.837
  '455': 'lastXRay', // chiropractic.837
  '471': 'hearingVisionPrescription', // vision.837
  '314': 'disability', // property-and-casualty-claim.837
  '297': 'lastWorked', // property-and-casualty-claim.837
  '296': 'authorizedReturnToWork', // property-and-casualty-claim.837
  '435': 'admission', // group-practice-multiple-renderers.837 — confirmed, x12Names837p.ts CLAIM_DTP_NAME.admission's own comment
  '096': 'discharge', // group-practice-multiple-renderers.837
  '090': 'assumedOrRelinquishedCare', // group-practice-multiple-renderers.837 (first of 2 occurrences)
  '091': 'assumedOrRelinquishedCare', // group-practice-multiple-renderers.837 (second of 2 occurrences)
  '444': 'propertyCasualtyFirstContact', // property-and-casualty-claim.837
  '050': 'repricerReceived', // repriced-claim.837
}

const SERVICE_DTP_SLUG: Record<string, keyof typeof SERVICE_DTP_NAME> = {
  '472': 'serviceDate', // every service-line example
  '471': 'prescription', // durable-medical-equipment.837
  '607': 'certificationRevision', // durable-medical-equipment.837
  '463': 'beginTherapy', // durable-medical-equipment.837
  '461': 'lastCertification', // durable-medical-equipment.837
  '304': 'lastTreatmentOrTherapy', // incident-to-service.837
  '738': 'testPerformed', // laboratory-test.837 (first of 2 occurrences)
  '739': 'testPerformed', // laboratory-test.837 (second of 2 occurrences)
  '011': 'shipped', // durable-medical-equipment.837
  '455': 'lastXRay', // chiropractic.837 (line level)
  '454': 'initialTreatment', // multiple-diagnoses.837 (line level)
}

const CLAIM_REF_SLUG: Record<string, keyof typeof CLAIM_REF_NAME> = {
  '4N': 'serviceAuthException', // office-visit.837
  F8: 'payerClaimControlNumber', // corrected-claim.837
  EW: 'mammographyCertification', // multiple-diagnoses.837
  EA: 'medicalRecordNumber', // multiple-diagnoses.837
  '9F': 'referral', // specialist-referral.837
  G1: 'priorAuthorization', // prior-authorization.837
  X4: 'clia', // laboratory-test.837
  '9A': 'repricedClaim', // repriced-claim.837
  '9C': 'adjustedRepricedClaim', // repriced-claim.837
  LX: 'investigationalDeviceExemption', // special-program-identifiers.837 — a REF01 value, unrelated to this viewer's LX (service line) box type
  P4: 'demonstrationProject', // special-program-identifiers.837
  '1J': 'carePlanOversight', // special-program-identifiers.837
  D9: 'claimIdentifierForTransmissionIntermediaries', // special-program-identifiers.837, REF01=D9 per Transmittal 907's own errata example (PDF p.244)
}

// Loop 2400's "Prior Authorization" REF box covers both prior-authorization
// and referral numbers under one qualifier range — SERVICE_REF02_EL_NAME.
// priorAuthorization is itself named "Prior Authorization or Referral
// Number", unlike the claim level's two separate boxes.
const SERVICE_REF_SLUG: Record<string, keyof typeof SERVICE_REF_NAME> = {
  '6R': 'lineItemControl', // every service-line example — REF01=6R confirmed, CEDI 837P CG Figure 3
  BT: 'immunizationBatch', // office-visit.837
  X4: 'clia', // laboratory-test.837 (line level)
  F4: 'referringClia', // laboratory-test.837
  '9B': 'repricedLineItem', // repriced-claim.837
  '9D': 'adjustedRepricedLineItem', // repriced-claim.837
  G1: 'priorAuthorization', // prior-authorization.837 (line level)
  '9F': 'priorAuthorization', // specialist-referral.837 (line level, reused — see docs/design/837p-examples.md "S-REF-PRIORAUTH (reused)")
}

const OTHER_PAYER_REF_SLUG: Record<string, keyof typeof OTHER_PAYER_REF_NAME> = {
  EI: 'secondaryIdentification', // medicare-secondary.837
  G1: 'priorAuthorization', // medicare-secondary.837
  '9F': 'referral', // medicare-secondary.837
  T4: 'claimAdjustmentIndicator', // medicare-secondary.837
}

const CLAIM_AMT_SLUG: Record<string, keyof typeof AMT02_EL_NAME_BY_BOX> = {
  D: 'cobPayerPaid', // medicare-secondary.837, line-level-other-payer.837 — AMT01=D confirmed, AMT01_COB_NAME
  A8: 'cobNonCovered', // medicare-secondary.837
  EAF: 'remainingPatientLiability', // medicare-secondary.837 (2320), line-level-other-payer.837 (2430, same qualifier)
}

const SERVICE_AMT_SLUG: Record<string, keyof typeof SERVICE_AMT02_EL_NAME_BY_BOX> = {
  T: 'salesTax', // durable-medical-equipment.837
  F4: 'postageClaimed', // durable-medical-equipment.837
}

const SERVICE_QTY_SLUG: Record<string, keyof typeof QTY02_EL_NAME_BY_BOX> = {
  PT: 'ambulancePatientCount', // ambulance.837
  FL: 'obstetricAdditionalUnits', // anesthesia-minutes.837
}

const CLAIM_CRC_SLUG: Record<string, keyof typeof CLAIM_CRC_NAME> = {
  '07': 'ambulanceCertification', // ambulance.837
  E1: 'visionCondition', // vision.837
  '75': 'homebound', // homebound-and-epsdt-referral.837
  ZZ: 'epsdtReferral', // homebound-and-epsdt-referral.837
}

const SERVICE_CRC_SLUG: Record<string, keyof typeof SERVICE_CRC_NAME> = {
  '07': 'ambulanceCertification', // ambulance.837 (line level)
  '09': 'dmeConditionIndicator', // durable-medical-equipment.837
}

// HI01-1/HI02-1/... Diagnosis Type Code qualifiers seen in the examples:
// ABK/BK/ABF select the Diagnosis box; BG selects Condition Information.
// Neither source names these letters in text (docs/design/837p-examples.md,
// "Open items"), so only which BOX a composite belongs in is resolved here
// — the qualifier itself stays an unresolved code, with no name on hover.
const HI_CONDITION_QUALIFIERS = new Set(['BG'])

// ---------------------------------------------------------------------------
// Small helpers
// ---------------------------------------------------------------------------

function carcLabel(code: string) {
  const entry = CARC[code]
  return label(code, entry?.short, { full: entry?.full })
}
function rarcLabel(code: string) {
  const entry = RARC[code]
  return label(code, entry?.short, { full: entry?.full })
}

// A composite element joined tight with ':' between parts, every part
// unresolved unless the caller overwrites one afterward — the same shape
// SVC01 uses on the 835 side (frontend/src/docs/viewer835/describe.ts,
// svc01Part): SV101/SV501/SVD03's product-or-service qualifier plus
// procedure code plus modifiers, and CLM05's place-of-service/
// facility-qualifier/frequency triplet.
function joinedParts(parts: string[]): Part {
  const labels = parts.map((p) => label(p, undefined))
  return labels.length === 1 ? labels[0] : { kind: 'joined', parts: labels }
}

// CAS holds up to six adjustment reason/amount/quantity triplets after its
// group code — identical shape to the 835 side.
function casTripletCells(e: string[]): Cell[] {
  const cells: Cell[] = []
  for (let i = 1; i + 1 < e.length; i += 3) {
    const reason = e[i]
    const amount = e[i + 1]
    const quantity = e[i + 2]
    if (reason === undefined) break
    cells.push(cell(ref('CAS', i + 1), CAS_REASON_EL_NAME, carcLabel(reason)))
    cells.push(vcell(ref('CAS', i + 2), CAS_AMOUNT_EL_NAME, amount, true))
    if (quantity !== undefined) cells.push(vcell(ref('CAS', i + 3), CAS_QTY_EL_NAME, quantity))
  }
  return cells
}

// HI holds up to 12 composites (HI01-HI12), each a qualifier:code pair
// (HI0n-1, HI0n-2) followed by HI0n-3 through HI0n-9, which the 837P marks
// Not Used but the standard defines, so a file can carry text there. One
// cell per composite present, drawn as it always was: no per-cell caption
// (every composite would otherwise share the same pair of names, which would
// mislabel whichever part a reader is looking at) — the lead already names
// the whole row. A component after the code that carries text gets its own
// cell right after its composite's cell, with its element reference (HI01-4)
// and its sourced name beneath (HI_NOT_USED_EL_NAME); an empty one gets no
// cell, since the reference already says which position it is. A composite
// that stops at the code draws exactly as before.
function hiCells(e: string[]): Cell[] {
  const cells: Cell[] = []
  e.forEach((composite, i) => {
    const parts = composite.split(':')
    const [qualifier, code] = parts
    const qualifierPart = label(qualifier, undefined)
    const part: Part = code !== undefined ? { kind: 'joined', parts: [qualifierPart, value(code)] } : qualifierPart
    cells.push(cell(ref('HI', i + 1), undefined, part))
    for (let k = 2; k < parts.length; k++) {
      if (parts[k] === '') continue
      // HI0n-5 (Monetary Amount) and HI0n-6 (Quantity) are numbers, flagged
      // as the viewer's other numbers are (display.ts ValuePart.amount).
      cells.push(vcell(`${ref('HI', i + 1)}-${k + 1}`, HI_NOT_USED_EL_NAME[k - 2], parts[k], k === 4 || k === 5))
    }
  })
  return cells
}

// PWK02's four Medicare-accepted values (BM/FX/EL/FT) have no decoded
// English text in either source (x12Names837p.ts, PWK02_MEDICARE_CODES'
// own comment) — rendered as an unresolved code either way.
function pwkCells(e: string[]): Cell[] {
  const n = PWK_EL_NAME
  const cells: Cell[] = [lcell(ref('PWK', 1), n[0], e[0], undefined)]
  if (e[1] !== undefined) cells.push(lcell(ref('PWK', 2), n[1], e[1], undefined))
  for (let i = 2; i < e.length; i++) cells.push(vcell(ref('PWK', i + 1), n[i], e[i]))
  return cells
}

// A provider/payer/subscriber NM1 loop's own element names, by its public
// X12 loop id (LOOP_NM1_NAME/NM1_ROLE_EL_NAME keys) — shared by every box
// type that is just "a name, maybe an address, maybe a secondary ID or
// contact" with no further nesting (frontend/src/docs/viewer837p/parse.ts's
// parseSimpleLoop covers the same set of box types this maps).
const LOOP_ID: Record<string, string> = {
  SUBMITTER: '1000A',
  RECEIVER: '1000B',
  BILLING_PROVIDER: '2010AA',
  PAYTO: '2010AB',
  SUBSCRIBER: '2010BA',
  PAYER: '2010BB',
  REFERRING: '2310A',
  RENDERING: '2310B',
  SERVICE_FACILITY: '2310C',
  SUPERVISING: '2310D',
  AMB_PICKUP: '2310E',
  AMB_DROPOFF: '2310F',
  OTHER_SUBSCRIBER_NAME: '2330A',
  OTHER_PAYER: '2330B',
  RENDERING_LINE: '2420A',
  PURCHASED_SERVICE: '2420B',
  SERVICE_FACILITY_LINE: '2420C',
  SUPERVISING_LINE: '2420D',
  ORDERING: '2420E',
  REFERRING_LINE: '2420F',
  AMB_PICKUP_LINE: '2420G',
  AMB_DROPOFF_LINE: '2420H',
}

function nm1OpeningDisplay(e: string[], loopId: string): SegmentDisplay {
  const roleNames = NM1_ROLE_EL_NAME[loopId]
  const cells: Cell[] = [
    lcell(ref('NM1', 1), NM101_EL_NAME, e[0], NM101_CONFIRMED_NAME[e[0]]),
    lcell(ref('NM1', 2), NM102_EL_NAME, e[1], undefined), // NM102 entity type — not independently confirmed
  ]
  for (let i = 2; i < e.length; i++) {
    const cellRef = ref('NM1', i + 1)
    const name = roleNames?.[i - 2]
    if (i === 7) cells.push(lcell(cellRef, name, e[7], NM108_CONFIRMED_NAME[e[7]]))
    else cells.push(vcell(cellRef, name, e[i]))
  }
  return { lead: label('NM1', LOOP_NM1_NAME[loopId]), cells }
}

// Ambulance pick-up/drop-off (2310E/F, 2420G/H): a bare name only (NM103),
// per Transmittal 907's own "(Not Used)" marking of every other NM1
// sub-element in those two loops.
function ambulanceLocationOpeningDisplay(e: string[], loopId: string, names: readonly [string]): SegmentDisplay {
  return {
    lead: label('NM1', LOOP_NM1_NAME[loopId]),
    cells: [lcell(ref('NM1', 1), NM101_EL_NAME, e[0], NM101_CONFIRMED_NAME[e[0]]), vcell(ref('NM1', 3), names[0], e[2])],
  }
}

// N3/N4/REF rows inside a simple NM1 loop, named generically by the loop's
// own public id — the same approach the 835 side uses for its payer/payee
// loop (frontend/src/docs/viewer835/describe.ts, PAYER_LOOP_NAME/PAYEE_LOOP_NAME).
function loopAddressRow(id: 'N3' | 'N4', e: string[], loopId: string): SegmentDisplay {
  const title = id === 'N3' ? LOOP_N3_NAME[loopId] : LOOP_N4_NAME[loopId]
  const names = id === 'N3' ? N3_EL_NAME : N4_EL_NAME
  return { lead: label(id, title), cells: e.map((v, i) => vcell(`${id}${pad(i + 1)}`, names[i], v)) }
}

function loopRefRow(e: string[], loopId: string): SegmentDisplay {
  const title = LOOP_REF_NAME[loopId]
  const el2 = LOOP_REF02_EL_NAME[loopId]
  return {
    lead: label('REF', title),
    cells: [lcell(ref('REF', 1), 'Reference Identification Qualifier', e[0], undefined), vcell(ref('REF', 2), el2, e[1])],
  }
}

// ---------------------------------------------------------------------------
// Opening rows
// ---------------------------------------------------------------------------

export function describeOpening(box: Box): SegmentDisplay {
  const e = box.opening.elements
  switch (box.boxType) {
    case 'ISA':
      return { lead: label('ISA', SEGMENT_NAME.ISA), cells: e.map((v, i) => vcell(ref('ISA', i + 1), ISA_EL_NAME[i], v)) }
    case 'GS':
      return {
        lead: label('GS', SEGMENT_NAME.GS),
        cells: [
          lcell(ref('GS', 1), GS_EL_NAME[0], e[0], TRANSACTION_TYPE_NAME[e[0]]),
          ...e.slice(1).map((v, i) => vcell(ref('GS', i + 2), GS_EL_NAME[i + 1], v)),
        ],
      }
    case 'ST':
      return {
        lead: label('ST', SEGMENT_NAME.ST),
        cells: [lcell(ref('ST', 1), ST_EL_NAME[0], e[0], TRANSACTION_TYPE_NAME[e[0]]), vcell(ref('ST', 2), ST_EL_NAME[1], e[1]), vcell(ref('ST', 3), ST_EL_NAME[2], e[2])],
      }
    case 'BILLING_HL':
    case 'SUBSCRIBER_HL': {
      const title = box.boxType === 'BILLING_HL' ? 'Billing Provider Hierarchical Level' : 'Subscriber Hierarchical Level'
      return { lead: label('HL', title), cells: e.map((v, i) => vcell(ref('HL', i + 1), HL_EL_NAME[i], v)) }
    }
    case 'CLM': {
      // Direct positional through CLM04; CLM05 is a 3-part composite (one
      // raw element, three names); CLM06-10 direct again; CLM11 is a
      // 4-part composite the examples never split (still one raw
      // element); CLM12 (Special Program Indicator) is the last position
      // any example sends — see this file's header note.
      const cells: Cell[] = [vcell(ref('CLM', 1), CLM_EL_NAME[0], e[0]), vcell(ref('CLM', 2), CLM_EL_NAME[1], e[1], true)]
      if (e[2] !== undefined) cells.push(vcell(ref('CLM', 3), CLM_EL_NAME[2], e[2]))
      if (e[3] !== undefined) cells.push(vcell(ref('CLM', 4), CLM_EL_NAME[3], e[3]))
      if (e[4] !== undefined) {
        const parts = e[4].split(':')
        const composite = joinedParts(parts)
        if (composite.kind === 'joined' && parts[2] !== undefined) composite.parts[2] = label(parts[2], CLM05_3_NAME[parts[2]])
        cells.push(cell(ref('CLM', 5), CLM_EL_NAME[4], composite))
      }
      if (e[5] !== undefined) cells.push(lcell(ref('CLM', 6), CLM_EL_NAME[7], e[5], undefined))
      if (e[6] !== undefined) cells.push(lcell(ref('CLM', 7), CLM_EL_NAME[8], e[6], undefined))
      if (e[7] !== undefined) cells.push(lcell(ref('CLM', 8), CLM_EL_NAME[9], e[7], undefined))
      if (e[8] !== undefined) cells.push(lcell(ref('CLM', 9), CLM_EL_NAME[10], e[8], undefined))
      if (e[9] !== undefined) cells.push(lcell(ref('CLM', 10), CLM_EL_NAME[11], e[9], undefined))
      if (e[10] !== undefined) cells.push(vcell(ref('CLM', 11), CLM_EL_NAME[12], e[10]))
      if (e[11] !== undefined) cells.push(vcell(ref('CLM', 12), CLM_EL_NAME[16], e[11]))
      return { lead: label('CLM', CLM_NAME), cells }
    }
    case 'SUBMITTER':
      return nm1OpeningDisplay(e, LOOP_ID.SUBMITTER)
    case 'RECEIVER':
      return nm1OpeningDisplay(e, LOOP_ID.RECEIVER)
    case 'BILLING_PROVIDER':
      return nm1OpeningDisplay(e, LOOP_ID.BILLING_PROVIDER)
    case 'PAYTO':
      return nm1OpeningDisplay(e, LOOP_ID.PAYTO)
    case 'SUBSCRIBER':
      return nm1OpeningDisplay(e, LOOP_ID.SUBSCRIBER)
    case 'PAYER':
      return nm1OpeningDisplay(e, LOOP_ID.PAYER)
    case 'REFERRING':
      return nm1OpeningDisplay(e, LOOP_ID.REFERRING)
    case 'RENDERING':
      return nm1OpeningDisplay(e, LOOP_ID.RENDERING)
    case 'SERVICE_FACILITY':
      return nm1OpeningDisplay(e, LOOP_ID.SERVICE_FACILITY)
    case 'SUPERVISING':
      return nm1OpeningDisplay(e, LOOP_ID.SUPERVISING)
    case 'OTHER_SUBSCRIBER_NAME':
      return nm1OpeningDisplay(e, LOOP_ID.OTHER_SUBSCRIBER_NAME)
    case 'OTHER_PAYER':
      return nm1OpeningDisplay(e, LOOP_ID.OTHER_PAYER)
    case 'RENDERING_LINE':
      return nm1OpeningDisplay(e, LOOP_ID.RENDERING_LINE)
    case 'PURCHASED_SERVICE':
      return nm1OpeningDisplay(e, LOOP_ID.PURCHASED_SERVICE)
    case 'SERVICE_FACILITY_LINE':
      return nm1OpeningDisplay(e, LOOP_ID.SERVICE_FACILITY_LINE)
    case 'SUPERVISING_LINE':
      return nm1OpeningDisplay(e, LOOP_ID.SUPERVISING_LINE)
    case 'ORDERING':
      return nm1OpeningDisplay(e, LOOP_ID.ORDERING)
    case 'REFERRING_LINE':
      return nm1OpeningDisplay(e, LOOP_ID.REFERRING_LINE)
    case 'AMB_PICKUP':
      return ambulanceLocationOpeningDisplay(e, LOOP_ID.AMB_PICKUP, AMBULANCE_PICKUP_NM1_EL_NAME)
    case 'AMB_DROPOFF':
      return ambulanceLocationOpeningDisplay(e, LOOP_ID.AMB_DROPOFF, AMBULANCE_DROPOFF_NM1_EL_NAME)
    case 'AMB_PICKUP_LINE':
      return ambulanceLocationOpeningDisplay(e, LOOP_ID.AMB_PICKUP, AMBULANCE_PICKUP_NM1_EL_NAME)
    case 'AMB_DROPOFF_LINE':
      return ambulanceLocationOpeningDisplay(e, LOOP_ID.AMB_DROPOFF, AMBULANCE_DROPOFF_NM1_EL_NAME)
    case 'OTHER_SUBSCRIBER': {
      // SBR01-05 direct positional; SBR06-08 compacted away (array has no
      // name for them); SBR09 (Claim Filing Indicator) is real raw index 8,
      // reading OTHER_SUBSCRIBER_SBR_EL_NAME[5] — see this file's header note.
      const n = OTHER_SUBSCRIBER_SBR_EL_NAME
      const cells: Cell[] = [lcell(ref('SBR', 1), n[0], e[0], SBR01_NAME[e[0]])]
      if (e[1] !== undefined) cells.push(vcell(ref('SBR', 2), n[1], e[1]))
      if (e[2] !== undefined) cells.push(vcell(ref('SBR', 3), n[2], e[2]))
      if (e[3] !== undefined) cells.push(vcell(ref('SBR', 4), n[3], e[3]))
      if (e[4] !== undefined) cells.push(vcell(ref('SBR', 5), n[4], e[4]))
      if (e[8] !== undefined) cells.push(vcell(ref('SBR', 9), n[5], e[8]))
      return { lead: label('SBR', 'Other Subscriber Information'), cells }
    }
    case 'LX':
      return { lead: label('LX', LX_NAME), cells: [vcell(ref('LX', 1), LX01_EL_NAME, e[0])] }
    case 'DRUG':
      // LIN01 (a sequence number) is compacted away; LIN02/LIN03 read raw
      // index 1/2.
      return {
        lead: label('LIN', LIN_NAME),
        cells: [lcell(ref('LIN', 2), LIN_EL_NAME[0], e[1], LIN02_NAME[e[1]]), vcell(ref('LIN', 3), LIN_EL_NAME[1], e[2])],
      }
    case 'LINE_ADJUDICATION': {
      // SVD01/02 direct; SVD03 is a composite (product qualifier + code);
      // SVD04 is compacted away, so SVD05/SVD06 read SVD_EL_NAME[9]/[10],
      // not [4]/[5] — see this file's header note.
      const n = SVD_EL_NAME
      const cells: Cell[] = [vcell(ref('SVD', 1), n[0], e[0]), vcell(ref('SVD', 2), n[1], e[1], true)]
      if (e[2] !== undefined) {
        const parts = e[2].split(':')
        const composite = joinedParts(parts)
        if (composite.kind === 'joined') composite.parts[0] = label(parts[0], SVD03_1_NAME[parts[0]])
        cells.push(cell(ref('SVD', 3), n[2], composite))
      }
      if (e[4] !== undefined) cells.push(vcell(ref('SVD', 5), n[9], e[4]))
      if (e[5] !== undefined) cells.push(vcell(ref('SVD', 6), n[10], e[5]))
      return { lead: label('SVD', SVD_NAME), cells }
    }
    case 'FORMS':
      return {
        lead: label('LQ', LQ_NAME),
        cells: [lcell(ref('LQ', 1), LQ_EL_NAME[0], e[0], undefined), vcell(ref('LQ', 2), LQ_EL_NAME[1], e[1])],
      }
  }
}

export function describeClosing(box: Box): SegmentDisplay {
  const seg = box.closing
  if (!seg) return { lead: label('', undefined), cells: [] }
  const e = seg.elements
  const names = seg.id === 'SE' ? SE_EL_NAME : seg.id === 'GE' ? GE_EL_NAME : seg.id === 'IEA' ? IEA_EL_NAME : undefined
  return {
    lead: label(seg.id, SEGMENT_NAME[seg.id]),
    cells: e.map((v, i) => vcell(`${seg.id}${pad(i + 1)}`, names?.[i], v)),
  }
}

// ---------------------------------------------------------------------------
// Plain rows
// ---------------------------------------------------------------------------

const SIMPLE_LOOP_HOSTS = new Set<BoxType>([
  'SUBMITTER',
  'RECEIVER',
  'BILLING_PROVIDER',
  'PAYTO',
  'SUBSCRIBER',
  'PAYER',
  'REFERRING',
  'RENDERING',
  'SERVICE_FACILITY',
  'SUPERVISING',
  'AMB_PICKUP',
  'AMB_DROPOFF',
  'OTHER_SUBSCRIBER_NAME',
  'OTHER_PAYER',
  'RENDERING_LINE',
  'PURCHASED_SERVICE',
  'SERVICE_FACILITY_LINE',
  'SUPERVISING_LINE',
  'ORDERING',
  'REFERRING_LINE',
  'AMB_PICKUP_LINE',
  'AMB_DROPOFF_LINE',
])

export function describeRow(segment: Segment, host: BoxType): SegmentDisplay {
  const e = segment.elements
  const id = segment.id

  // N3/N4/REF inside a simple provider/payer/subscriber NM1 loop are named
  // generically by that loop's own public id, except the handful of loops
  // with a single fixed-purpose REF box named individually below.
  if (SIMPLE_LOOP_HOSTS.has(host) && (id === 'N3' || id === 'N4' || id === 'REF')) {
    const loopId = LOOP_ID[host]
    if (id === 'REF') {
      if (host === 'BILLING_PROVIDER')
        return {
          lead: label('REF', BILLING_PROVIDER_REF_NAME),
          cells: [lcell(ref('REF', 1), BILLING_PROVIDER_REF_EL_NAME[0], e[0], undefined), vcell(ref('REF', 2), BILLING_PROVIDER_REF_EL_NAME[1], e[1])],
        }
      if (host === 'SUBSCRIBER' && e[0] === 'Y4')
        return {
          lead: label('REF', SUBSCRIBER_PC_REF_NAME),
          cells: [lcell(ref('REF', 1), 'Reference Identification Qualifier', e[0], undefined), vcell(ref('REF', 2), SUBSCRIBER_PC_REF_NAME, e[1])],
        }
      if (host === 'OTHER_PAYER') {
        const slug = OTHER_PAYER_REF_SLUG[e[0]]
        const title = slug ? OTHER_PAYER_REF_NAME[slug] : undefined
        const el2 = slug ? OTHER_PAYER_REF02_EL_NAME[slug] : undefined
        return { lead: label('REF', undefined), cells: [lcell(ref('REF', 1), undefined, e[0], title), vcell(ref('REF', 2), el2, e[1])] }
      }
      return loopRefRow(e, loopId)
    }
    return loopAddressRow(id, e, loopId)
  }
  if (host === 'SUBSCRIBER' && id === 'PER')
    return {
      lead: label('PER', SUBSCRIBER_PC_PER_NAME),
      cells: [lcell(ref('PER', 1), 'Contact Function Code', e[0], undefined), ...e.slice(1).map((v, i) => vcell(ref('PER', i + 2), undefined, v))],
    }
  if (host === 'BILLING_PROVIDER' && id === 'PER') {
    const n = BILLING_PROVIDER_PER_EL_NAME
    return {
      lead: label('PER', BILLING_PROVIDER_PER_NAME),
      cells: [lcell(ref('PER', 1), n[0], e[0], undefined), ...e.slice(1).map((v, i) => vcell(ref('PER', i + 2), n[i + 1], v))],
    }
  }
  if (host === 'ORDERING' && id === 'PER') {
    const n = ORDERING_PROVIDER_PER_EL_NAME
    return {
      lead: label('PER', ORDERING_PROVIDER_PER_NAME),
      cells: [lcell(ref('PER', 1), n[0], e[0], undefined), ...e.slice(1).map((v, i) => vcell(ref('PER', i + 2), n[i + 1], v))],
    }
  }
  if (host === 'SUBMITTER' && id === 'PER') {
    const n = SUBMITTER_PER_EL_NAME
    return {
      lead: label('PER', 'Submitter EDI Contact Information'),
      cells: [lcell(ref('PER', 1), n[0], e[0], undefined), ...e.slice(1).map((v, i) => vcell(ref('PER', i + 2), n[i + 1], v))],
    }
  }
  if (host === 'DRUG') {
    if (id === 'CTP') {
      // CTP01-03 compacted away; CTP04/CTP05-1 read raw index 3/4.
      const n = CTP_EL_NAME
      const cells: Cell[] = []
      if (e[3] !== undefined) cells.push(vcell(ref('CTP', 4), n[0], e[3]))
      if (e[4] !== undefined) cells.push(lcell(ref('CTP', 5), n[1], e[4], CTP05_1_NAME[e[4]]))
      return { lead: label('CTP', CTP_NAME), cells }
    }
    if (id === 'REF')
      return {
        lead: label('REF', DRUG_ASSOCIATION_REF_NAME),
        cells: [lcell(ref('REF', 1), 'Reference Identification Qualifier', e[0], DRUG_ASSOCIATION_REF01_NAME[e[0]]), vcell(ref('REF', 2), DRUG_ASSOCIATION_REF02_EL_NAME, e[1])],
      }
  }
  if (host === 'LINE_ADJUDICATION') {
    if (id === 'CAS') {
      const group = lcell(ref('CAS', 1), CAS01_EL_NAME, e[0], CAS_GROUP_NAME[e[0]])
      return { lead: label('CAS', 'Line Adjustment'), cells: [group, ...casTripletCells(e)] }
    }
    if (id === 'DTP')
      return {
        lead: label('DTP', LINE_ADJUDICATION_DTP_NAME),
        cells: [lcell(ref('DTP', 1), DTP01_EL_NAME, e[0], undefined), vcell(ref('DTP', 2), LINE_ADJUDICATION_DTP03_EL_NAME, e[1])],
      }
    if (id === 'AMT')
      return {
        lead: label('AMT', LINE_ADJUDICATION_REMAINING_LIABILITY_AMT_NAME),
        cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], undefined), vcell(ref('AMT', 2), AMT02_EL_NAME_BY_BOX.remainingPatientLiability, e[1], true)],
      }
  }
  if (host === 'OTHER_SUBSCRIBER') {
    if (id === 'CAS') {
      const group = lcell(ref('CAS', 1), CAS01_EL_NAME, e[0], CAS_GROUP_NAME[e[0]])
      return { lead: label('CAS', CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME), cells: [group, ...casTripletCells(e)] }
    }
    if (id === 'AMT') {
      const slug = CLAIM_AMT_SLUG[e[0]]
      const title = slug ? AMT02_EL_NAME_BY_BOX[slug] : undefined
      return {
        lead: label('AMT', title),
        cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], slug ? AMT01_COB_NAME[e[0]] : undefined), vcell(ref('AMT', 2), title, e[1], true)],
      }
    }
    if (id === 'OI') {
      // OI01/02/05 compacted away; OI03/04/06 read raw index 2/3/5.
      const cells: Cell[] = []
      if (e[0] !== undefined) cells.push(vcell(ref('OI', 1), undefined, e[0]))
      if (e[1] !== undefined) cells.push(vcell(ref('OI', 2), undefined, e[1]))
      if (e[2] !== undefined) cells.push(lcell(ref('OI', 3), OI_EL_NAME[0], e[2], undefined))
      if (e[3] !== undefined) cells.push(lcell(ref('OI', 4), OI_EL_NAME[1], e[3], undefined))
      if (e[4] !== undefined) cells.push(vcell(ref('OI', 5), undefined, e[4]))
      if (e[5] !== undefined) cells.push(lcell(ref('OI', 6), OI_EL_NAME[2], e[5], undefined))
      return { lead: label('OI', OI_NAME), cells }
    }
    if (id === 'MOA') {
      const n = MOA_EL_NAME
      return {
        lead: label('MOA', MOA_NAME),
        cells: [
          vcell(ref('MOA', 1), n[0], e[0]),
          vcell(ref('MOA', 2), n[1], e[1], true),
          ...e.slice(2).map((code, i) => (code ? cell(ref('MOA', i + 3), n[i + 2], rarcLabel(code)) : vcell(ref('MOA', i + 3), n[i + 2], code))),
        ],
      }
    }
  }
  if (host === 'OTHER_PAYER' && id === 'DTP')
    return {
      lead: label('DTP', OTHER_PAYER_DTP_NAME),
      cells: [lcell(ref('DTP', 1), DTP01_EL_NAME, e[0], undefined), vcell(ref('DTP', 2), OTHER_PAYER_DTP03_EL_NAME, e[1])],
    }

  switch (id) {
    case 'BHT': {
      const n = BHT_EL_NAME
      return {
        lead: label('BHT', SEGMENT_NAME.BHT),
        cells: [
          vcell(ref('BHT', 1), n[0], e[0]),
          lcell(ref('BHT', 2), n[1], e[1], BHT02_NAME[e[1]]),
          vcell(ref('BHT', 3), n[2], e[2]),
          vcell(ref('BHT', 4), n[3], e[3]),
          vcell(ref('BHT', 5), n[4], e[4]),
          lcell(ref('BHT', 6), n[5], e[5], BHT06_NAME[e[5]]),
        ],
      }
    }
    case 'PRV':
      return {
        lead: label('PRV', 'Billing Provider Specialty Information'),
        cells: [lcell(ref('PRV', 1), PRV_EL_NAME[0], e[0], undefined), lcell(ref('PRV', 2), PRV_EL_NAME[1], e[1], undefined), vcell(ref('PRV', 3), PRV_EL_NAME[2], e[2])],
      }
    case 'CUR':
      return { lead: label('CUR', CUR_NAME), cells: e.map((v, i) => vcell(ref('CUR', i + 1), CUR_EL_NAME[i], v)) }
    case 'SBR': {
      // SBR01-05 direct positional; SBR06-08 compacted away; SBR09 reads
      // SBR_EL_NAME[5] from raw index 8 — see this file's header note.
      const n = SBR_EL_NAME
      const cells: Cell[] = [lcell(ref('SBR', 1), n[0], e[0], SBR01_NAME[e[0]])]
      if (e[1] !== undefined) cells.push(lcell(ref('SBR', 2), n[1], e[1], SBR02_NAME[e[1]]))
      if (e[2] !== undefined) cells.push(vcell(ref('SBR', 3), n[2], e[2]))
      if (e[3] !== undefined) cells.push(vcell(ref('SBR', 4), n[3], e[3]))
      if (e[4] !== undefined) cells.push(vcell(ref('SBR', 5), n[4], e[4]))
      if (e[8] !== undefined) cells.push(lcell(ref('SBR', 9), n[5], e[8], SBR09_NAME[e[8]]))
      return { lead: label('SBR', 'Subscriber Information'), cells }
    }
    case 'PAT':
      return { lead: label('PAT', PAT_NAME), cells: e.map((v, i) => vcell(ref('PAT', i + 1), PAT_EL_NAME[i], v)) }
    case 'DMG':
      return {
        lead: label('DMG', DMG_NAME),
        cells: [lcell(ref('DMG', 1), DMG_EL_NAME[0], e[0], undefined), vcell(ref('DMG', 2), DMG_EL_NAME[1], e[1]), lcell(ref('DMG', 3), DMG_EL_NAME[2], e[2], undefined)],
      }
    case 'DTP': {
      if (host === 'CLM') {
        const slug = CLAIM_DTP_SLUG[e[0]]
        const title = slug ? CLAIM_DTP_NAME[slug] : undefined
        const el3 = slug ? CLAIM_DTP03_EL_NAME[slug] : undefined
        return {
          lead: label('DTP', title),
          cells: [lcell(ref('DTP', 1), DTP01_EL_NAME, e[0], undefined), lcell(ref('DTP', 2), DTP02_EL_NAME, e[1], undefined), vcell(ref('DTP', 3), el3, e[2])],
        }
      }
      if (host === 'LX') {
        const slug = SERVICE_DTP_SLUG[e[0]]
        const title = slug ? SERVICE_DTP_NAME[slug] : undefined
        const el3 = slug ? SERVICE_DTP03_EL_NAME[slug] : undefined
        return {
          lead: label('DTP', title),
          cells: [lcell(ref('DTP', 1), DTP01_EL_NAME, e[0], undefined), lcell(ref('DTP', 2), DTP02_EL_NAME, e[1], undefined), vcell(ref('DTP', 3), el3, e[2])],
        }
      }
      return { lead: label('DTP', undefined), cells: [lcell(ref('DTP', 1), DTP01_EL_NAME, e[0], undefined), vcell(ref('DTP', 2), DTP02_EL_NAME, e[1]), vcell(ref('DTP', 3), undefined, e[2])] }
    }
    case 'PWK':
      return { lead: label('PWK', host === 'LX' ? SERVICE_PWK_NAME : CLAIM_PWK_NAME), cells: pwkCells(e) }
    case 'CN1':
      return { lead: label('CN1', CN1_NAME), cells: e.map((v, i) => vcell(ref('CN1', i + 1), CN1_EL_NAME[i], v)) }
    case 'AMT': {
      if (host === 'CLM')
        return {
          lead: label('AMT', CLAIM_AMT_NAME),
          cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], undefined), vcell(ref('AMT', 2), CLAIM_AMT_EL_NAME[1], e[1], true)],
        }
      if (host === 'LX') {
        const slug = SERVICE_AMT_SLUG[e[0]]
        const title = slug ? SERVICE_AMT02_EL_NAME_BY_BOX[slug] : undefined
        return {
          lead: label('AMT', slug === 'salesTax' ? SALES_TAX_AMT_NAME : slug === 'postageClaimed' ? POSTAGE_CLAIMED_AMT_NAME : title),
          cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], undefined), vcell(ref('AMT', 2), title, e[1], true)],
        }
      }
      return { lead: label('AMT', undefined), cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], undefined), vcell(ref('AMT', 2), undefined, e[1], true)] }
    }
    case 'REF': {
      if (host === 'CLM') {
        const slug = CLAIM_REF_SLUG[e[0]]
        const title = slug ? CLAIM_REF_NAME[slug] : undefined
        const el2 = slug ? CLAIM_REF02_EL_NAME[slug] : undefined
        return { lead: label('REF', undefined), cells: [lcell(ref('REF', 1), undefined, e[0], title), vcell(ref('REF', 2), el2, e[1])] }
      }
      if (host === 'LX') {
        const slug = SERVICE_REF_SLUG[e[0]]
        const title = slug ? SERVICE_REF_NAME[slug] : undefined
        const el2 = slug ? SERVICE_REF02_EL_NAME[slug] : undefined
        return { lead: label('REF', undefined), cells: [lcell(ref('REF', 1), undefined, e[0], title), vcell(ref('REF', 2), el2, e[1])] }
      }
      return { lead: label('REF', undefined), cells: [lcell(ref('REF', 1), undefined, e[0], undefined), vcell(ref('REF', 2), undefined, e[1])] }
    }
    case 'K3':
      return { lead: label('K3', host === 'LX' ? SERVICE_K3_NAME : CLAIM_K3_NAME), cells: [vcell(ref('K3', 1), 'Fixed Format Information', e[0])] }
    case 'NTE': {
      const title = host === 'LX' ? (e[0] === 'TPO' ? THIRD_PARTY_NTE_NAME : LINE_NTE_NAME) : CLAIM_NTE_NAME
      return { lead: label('NTE', title), cells: [lcell(ref('NTE', 1), 'Note Reference Code', e[0], undefined), vcell(ref('NTE', 2), 'Note Text', e[1])] }
    }
    case 'CR1':
      return { lead: label('CR1', host === 'LX' ? SERVICE_CR1_NAME : CLAIM_CR1_NAME), cells: e.map((v, i) => vcell(ref('CR1', i + 1), CR1_EL_NAME[i], v)) }
    case 'CR2':
      // CR201-07/209/212 compacted away; only CR208/CR210/CR211 are named,
      // reading raw index 7/9/10 — see this file's header note.
      return {
        lead: label('CR2', CR2_NAME),
        cells: [
          ...(e[7] !== undefined ? [lcell(ref('CR2', 8), CR2_EL_NAME[0], e[7], undefined)] : []),
          ...(e[9] !== undefined ? [vcell(ref('CR2', 10), CR2_EL_NAME[1], e[9])] : []),
          ...(e[10] !== undefined ? [vcell(ref('CR2', 11), CR2_EL_NAME[2], e[10])] : []),
        ],
      }
    case 'CR3':
      return { lead: label('CR3', CR3_NAME), cells: e.map((v, i) => vcell(ref('CR3', i + 1), CR3_EL_NAME[i], v)) }
    case 'CRC': {
      const slug = host === 'LX' ? SERVICE_CRC_SLUG[e[0]] : CLAIM_CRC_SLUG[e[0]]
      const title = slug ? (host === 'LX' ? SERVICE_CRC_NAME[slug] : CLAIM_CRC_NAME[slug]) : undefined
      const n = CRC_EL_NAME
      return {
        lead: label('CRC', title),
        cells: [
          lcell(ref('CRC', 1), n[0], e[0], undefined),
          lcell(ref('CRC', 2), n[1], e[1], undefined),
          ...e.slice(2).map((v, i) => (v ? lcell(ref('CRC', i + 3), n[i + 2], v, undefined) : vcell(ref('CRC', i + 3), n[i + 2], v))),
        ],
      }
    }
    case 'HI': {
      const isCondition = HI_CONDITION_QUALIFIERS.has(e[0]?.split(':')[0] ?? '')
      return { lead: label('HI', isCondition ? HI_CONDITION_NAME : HI_DIAGNOSIS_NAME), cells: hiCells(e) }
    }
    case 'HCP': {
      // HCP01-04 direct positional; HCP08's own compaction is never
      // exercised by these examples, which never send past HCP04.
      const n = HCP_EL_NAME
      const cells: Cell[] = [lcell(ref('HCP', 1), n[0], e[0], undefined), vcell(ref('HCP', 2), n[1], e[1], true)]
      if (e[2] !== undefined) cells.push(vcell(ref('HCP', 3), n[2], e[2], true))
      if (e[3] !== undefined) cells.push(vcell(ref('HCP', 4), n[3], e[3]))
      return { lead: label('HCP', host === 'LX' ? LINE_HCP_NAME : HCP_NAME), cells }
    }
    case 'QTY': {
      const slug = SERVICE_QTY_SLUG[e[0]]
      const title = slug === 'ambulancePatientCount' ? AMBULANCE_PATIENT_COUNT_QTY_NAME : slug === 'obstetricAdditionalUnits' ? OBSTETRIC_ANESTHESIA_QTY_NAME : undefined
      const el2 = slug ? QTY02_EL_NAME_BY_BOX[slug] : undefined
      return {
        lead: label('QTY', title),
        cells: [lcell(ref('QTY', 1), QTY01_EL_NAME, e[0], undefined), vcell(ref('QTY', 2), el2, e[1])],
      }
    }
    case 'MEA':
      return {
        lead: label('MEA', MEA_NAME),
        cells: [lcell(ref('MEA', 1), MEA_EL_NAME[0], e[0], undefined), lcell(ref('MEA', 2), MEA_EL_NAME[1], e[1], undefined), vcell(ref('MEA', 3), MEA_EL_NAME[2], e[2])],
      }
    case 'PS1':
      return { lead: label('PS1', PS1_NAME), cells: [vcell(ref('PS1', 1), PS1_EL_NAME[0], e[0]), vcell(ref('PS1', 2), PS1_EL_NAME[1], e[1], true)] }
    case 'FRM':
      return { lead: label('FRM', FRM_NAME), cells: e.map((v, i) => vcell(ref('FRM', i + 1), FRM_EL_NAME[i], v)) }
    case 'SV1': {
      // SV101 composite (7 sub-parts, one raw element); SV102-105 direct
      // positional from there (raw index 1-4); SV106 never sent; SV107 (a
      // diagnosis-pointer composite, itself up to 4 digits joined with
      // ':') shown as one cell holding exactly what the file sent.
      const n = SV1_EL_NAME
      const parts1 = e[0].split(':')
      const composite = joinedParts(parts1)
      if (composite.kind === 'joined') composite.parts[0] = label(parts1[0], SV101_1_NAME[parts1[0]])
      const cells: Cell[] = [cell(ref('SV1', 1), n[0], composite)]
      if (e[1] !== undefined) cells.push(vcell(ref('SV1', 2), n[7], e[1], true))
      if (e[2] !== undefined) cells.push(lcell(ref('SV1', 3), n[8], e[2], undefined))
      if (e[3] !== undefined) cells.push(vcell(ref('SV1', 4), n[9], e[3]))
      if (e[4]) cells.push(lcell(ref('SV1', 5), n[10], e[4], undefined))
      if (e[6] !== undefined) cells.push(vcell(ref('SV1', 7), n[11], e[6]))
      return { lead: label('SV1', SV1_NAME), cells }
    }
    case 'SV5': {
      // SV501 composite (one raw element); SV502-506 direct from raw index
      // 1, reading SV5_EL_NAME[7+(i-1)] for raw index i.
      const n = SV5_EL_NAME
      const composite = joinedParts(e[0].split(':'))
      const cells: Cell[] = [cell(ref('SV5', 1), n[0], composite)]
      if (e[1] !== undefined) cells.push(lcell(ref('SV5', 2), n[7], e[1], undefined))
      for (let i = 2; i < e.length; i++) cells.push(vcell(ref('SV5', i + 1), n[7 + (i - 1)], e[i]))
      return { lead: label('SV5', SV5_NAME), cells }
    }
    default:
      return { lead: label(id, SEGMENT_NAME[id]), cells: e.map((v, i) => vcell(`${id}${pad(i + 1)}`, undefined, v)) }
  }
}
