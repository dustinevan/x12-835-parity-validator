// Turns a parsed 837I segment or box opening/closing into the row of
// display cells the page renders — same role and conventions as
// frontend/src/docs/viewer835/describe.ts (read that file first; this one
// copies its shape). Every name comes from
// frontend/src/data/x12Names837i.ts, which is built only from the spec's
// tables; a coded value with no entry there (the file's own "Values" list
// with no per-code text, or an element Transmittal 907 marks Not Used)
// renders as its raw code alone, with no name on hover.
//
// A segment whose row-level identity depends on a qualifier value (NM1, REF
// outside a fixed single-qualifier loop, DTP, NTE, AMT) leads its row with
// its bare segment ID only — the resolved meaning lives entirely in the
// qualifier's own cell. A segment whose name is fixed regardless of any
// qualifier (ISA, GS, ST, HL, CLM, SV2, LIN, SVD, CRC, …) leads with that
// true name instead.
//
// `codeCell` is this file's one addition over the 835 pattern: a coded
// element that is present but empty (CL1's three NUBC positions are often
// sent blank — "CL1***01" — a UB-04 field the claim did not need) renders as
// the faint value dot every other empty element gets, not as a blank label
// with no marker; a non-empty code still gets the label treatment (the raw
// code, with the resolved name on hover where one exists).
import { CARC, RARC } from '../../data/x12Names835'
import {
  AMT01_NAME,
  BHT_EL_NAME,
  CAS01_EL_NAME,
  CAS_AMOUNT_EL_NAME,
  CAS_CLAIM_NAME,
  CAS_GROUP_NAME,
  CAS_LINE_NAME,
  CAS_QTY_EL_NAME,
  CAS_REASON_EL_NAME,
  CL1_EL_NAME,
  CLM05_2_NAME,
  CLM07_NAME,
  CLM08_NAME,
  CLM09_NAME,
  CLM_EL_NAME,
  CN1_EL_NAME,
  CRC01_NAME,
  CRC02_NAME,
  CRC_EL_NAME,
  CTP_EL_NAME,
  CUR01_NAME,
  CUR_EL_NAME,
  DMG01_NAME,
  DMG03_NAME,
  DMG_EL_NAME,
  DTP01_EL_NAME,
  DTP472_NAME,
  DTP_CLAIM_NAME,
  DTP_OTHER_PAYER_NAME,
  GE_EL_NAME,
  GS_EL_NAME,
  HCP01_NAME,
  HCP_EL_NAME,
  HI_POA_NAME,
  HI_TYPE_NAME,
  HL03_NAME,
  HL_EL_NAME,
  IEA_EL_NAME,
  ISA_EL_NAME,
  K3_EL_NAME,
  LIN02_NAME,
  LIN_EL_NAME,
  LX01_EL_NAME,
  MIA_EL_NAME,
  MOA_EL_NAME,
  NM101_EL_NAME,
  NM102_EL_NAME,
  NM102_NAME,
  NM108_NAME,
  NM1_ATTENDING_EL_NAME,
  NM1_BILLING_EL_NAME,
  NM1_OPERATING_EL_NAME,
  NM1_OTHER_OPERATING_EL_NAME,
  NM1_OTHER_PAYER_EL_NAME,
  NM1_OTHER_SUBSCRIBER_EL_NAME,
  NM1_PATIENT_EL_NAME,
  NM1_PAYER_EL_NAME,
  NM1_PAYTOPLAN_EL_NAME,
  NM1_PAYTO_EL_NAME,
  NM1_RECEIVER_EL_NAME,
  NM1_REFERRING_EL_NAME,
  NM1_RENDERING_EL_NAME,
  NM1_ROLE_NAME,
  NM1_SERVICE_FACILITY_EL_NAME,
  NM1_SUBMITTER_EL_NAME,
  NM1_SUBSCRIBER_EL_NAME,
  NTE_BILLING_NAME,
  NTE_CLAIM_NAME,
  NTE_EL_NAME,
  NTE_LINE_NAME,
  OI03_NAME,
  OI06_NAME,
  OI_EL_NAME,
  PAT_EL_NAME,
  PER01_NAME,
  PER_EL_NAME,
  PWK_EL_NAME,
  PRV01_NAME,
  PRV02_NAME,
  PRV_EL_NAME,
  PROVIDER_REF_NAME,
  REF_BILLING_TAX_NAME,
  REF_CLAIM_NAME,
  REF_LINE_NAME,
  REF_OTHER_PAYER_NAME,
  REF_OTHER_SUBSCRIBER_NAME,
  REF_PATIENT_NAME,
  REF_PAYER_NAME,
  REF_PAYTOPLAN_NAME,
  REF_PRESCRIPTION_NAME,
  REF_SUBSCRIBER_NAME,
  RELATIONSHIP_NAME,
  SBR01_NAME,
  SBR09_NAME,
  SBR_EL_NAME,
  SEGMENT_NAME,
  SE_EL_NAME,
  ST_EL_NAME,
  SVD_EL_NAME,
  SV2_EL_NAME,
  TRANSACTION_TYPE_NAME,
} from '../../data/x12Names837i'
import { cell, label, pad, ref, value, vcell, type Cell, type LabelPart, type Part, type SegmentDisplay, type ValuePart } from '../viewer/display'
import { CAS_GROUP_DOT } from './colors'
import type { Box, BoxType, Segment } from './parse'

export type { Cell, JoinedPart, LabelPart, Part, SegmentDisplay, ValuePart } from '../viewer/display'
export { describeCompact } from '../viewer/display'

// PER01 ('Contact Function Code') is the same at every PER instance, but
// x12Names837i.ts names only its code→meaning table (PER01_NAME), not the
// element position itself — see that file's PER comment.
const PER01_EL_NAME = 'Contact Function Code' // Transmittal 907, PDF p.9

function carcLabel(code: string): LabelPart {
  const entry = CARC[code]
  return label(code, entry?.short, { full: entry?.full })
}
function rarcLabel(code: string): LabelPart {
  const entry = RARC[code]
  return label(code, entry?.short, { full: entry?.full })
}

// A coded element: present and non-empty gets the label treatment (the raw
// code, with the resolved name on hover when `table` has one);
// present-but-empty falls back to a plain value so it renders the same
// faint "·" every other empty element gets, instead of a blank label.
function codeCell(cellRef: string, name: string | undefined, code: string | undefined, table?: Record<string, string | undefined>): Cell {
  return code ? cell(cellRef, name, label(code, table?.[code])) : vcell(cellRef, name, code)
}

// Same as `codeCell`, for the (REF, DTP) cases whose resolved name is
// already computed by a per-host lookup function rather than a plain table.
function codeCellNamed(cellRef: string, name: string | undefined, code: string | undefined, resolvedName: string | undefined): Cell {
  return code ? cell(cellRef, name, label(code, resolvedName)) : vcell(cellRef, name, code)
}

// SV202 (service procedure), SVD03 (other payer's procedure) are both the
// same composite shape as the 835's SVC01/SVC06: a code-system qualifier,
// then the procedure code. SV202-1/SVD03-1's qualifier values (ER/HC/HP/
// IV/WK) have no public name beyond the code itself (Transmittal 907 gives
// only the value list), so the qualifier always renders unresolved.
function qualifierCodePart(raw: string | undefined): Part {
  if (!raw) return value('')
  const [qualifier, code] = raw.split(':')
  const qualifierPart = label(qualifier, undefined)
  if (!code) return qualifierPart
  return { kind: 'joined', parts: [qualifierPart, value(code)] }
}

// CLM05 (the UB-04 Type of Bill, split into place-of-service, facility code
// qualifier, and claim frequency) — three NUBC/X12 sub-elements joined
// tight, matching the file's own composite shape.
function clm05Part(raw: string | undefined): Part {
  const [pos, fac, freq] = (raw ?? '').split(':')
  const parts: LabelPart[] = []
  if (pos !== undefined) parts.push(label(pos, undefined))
  if (fac !== undefined) parts.push(label(fac, CLM05_2_NAME[fac]))
  if (freq !== undefined) parts.push(label(freq, undefined))
  if (parts.length === 0) return value('')
  if (parts.length === 1) return parts[0]
  return { kind: 'joined', parts }
}

// CAS holds up to six reason/amount/quantity triplets after its group code;
// each triplet gets three cells (the quantity cell only when the file sent
// one) — identical shape to the 835's own CAS (frontend/src/docs/viewer835/describe.ts).
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

// MIA holds five remark-code positions among otherwise plain counts/amounts
// — the same segment shape and relative positions as the 835's MIA
// (x12Names837i.ts's own MIA_EL_NAME comment): MIA05 (index 4) and
// MIA20-23 (indices 19-22).
const MIA_REMARK_INDEXES = new Set([4, 19, 20, 21, 22])
function miaCells(e: string[]): Cell[] {
  return e.map((v, i) => {
    const cellRef = ref('MIA', i + 1)
    const name = MIA_EL_NAME[i]
    return MIA_REMARK_INDEXES.has(i) && v ? cell(cellRef, name, rarcLabel(v)) : vcell(cellRef, name, v)
  })
}

// SBR (2000B's own, and the identical shape at 2320's Other Subscriber
// Information box opening).
function sbrCells(e: string[]): Cell[] {
  const n = SBR_EL_NAME
  return e.map((v, i) => {
    const cellRef = ref('SBR', i + 1)
    if (i === 0) return codeCell(cellRef, n[0], v, SBR01_NAME)
    if (i === 1) return codeCell(cellRef, n[1], v, RELATIONSHIP_NAME)
    if (i === 8) return codeCell(cellRef, n[8], v, SBR09_NAME)
    return vcell(cellRef, n[i], v)
  })
}

// A PER segment's qualifier+value pairs beyond PER01/02 (PER03/04, PER05/06,
// PER07/08): a communication-number qualifier (no public name beyond the
// code) followed by the communication number itself.
function perCommunicationCells(e: string[]): Cell[] {
  const cells: Cell[] = []
  for (let i = 2; i < e.length; i += 2) {
    const nameIdx = i - 1 // PER_EL_NAME starts at PER02
    cells.push(codeCell(ref('PER', i + 1), PER_EL_NAME[nameIdx], e[i]))
    if (e[i + 1] !== undefined) cells.push(vcell(ref('PER', i + 2), PER_EL_NAME[nameIdx + 1], e[i + 1]))
  }
  return cells
}

// HI carries up to twelve "Healthcare Code Information" composites in one
// segment (HI01 through HI12); value-codes and condition-codes each pack
// two same-type composites into one HI segment, so every element of the
// split array is its own cell, not just the first. Sub-parts beyond the
// type/code pair only ever appear together with a meaning this file states:
// a date (occurrence/procedure), an amount (value code), or a Present on
// Admission flag (diagnosis composites' 9th sub-part) — never more than one
// of the three in the examples below, so they render in that order without
// needing to know which HI qualifier sent them.
function hiCompositeCells(e: string[]): Cell[] {
  return e.map((raw, i) => {
    const sub = raw.split(':')
    const type = sub[0]
    const code = sub[1]
    const date = sub[3]
    const amount = sub[4]
    const poa = sub[8]
    const parts: (LabelPart | ValuePart)[] = [label(type, HI_TYPE_NAME[type])]
    if (code !== undefined) parts.push(value(code))
    if (date) parts.push(value(date))
    if (amount) parts.push(value(amount, true))
    if (poa) parts.push(label(poa, HI_POA_NAME[poa]))
    const part: Part = parts.length === 1 ? parts[0] : { kind: 'joined', parts }
    return cell(ref('HI', i + 1), undefined, part)
  })
}

// Which element-name array an NM1-opened box uses — every NM1-opened box
// type except the ones with their own SBR/HL/CLM/LX opening (see parse.ts's
// BoxType). 2330C-I and 2420A-D reuse their 2310 counterparts' shape, per
// x12Names837i.ts's own OTHER_PAYER_PROVIDER_LOOP_NAME/LINE_PROVIDER_LOOP_NAME
// comments.
const NM1_EL_NAME_BY_BOX: Partial<Record<BoxType, string[]>> = {
  SUBMITTER: NM1_SUBMITTER_EL_NAME,
  RECEIVER: NM1_RECEIVER_EL_NAME,
  'BILLING-PROVIDER': NM1_BILLING_EL_NAME,
  'PAY-TO-ADDRESS': NM1_PAYTO_EL_NAME,
  'PAY-TO-PLAN': NM1_PAYTOPLAN_EL_NAME,
  SUBSCRIBER: NM1_SUBSCRIBER_EL_NAME,
  PAYER: NM1_PAYER_EL_NAME,
  PATIENT: NM1_PATIENT_EL_NAME,
  ATTENDING: NM1_ATTENDING_EL_NAME,
  OPERATING: NM1_OPERATING_EL_NAME,
  'OTHER-OPERATING': NM1_OTHER_OPERATING_EL_NAME,
  RENDERING: NM1_RENDERING_EL_NAME,
  'SERVICE-FACILITY': NM1_SERVICE_FACILITY_EL_NAME,
  REFERRING: NM1_REFERRING_EL_NAME,
  'OTHER-SUBSCRIBER-NAME': NM1_OTHER_SUBSCRIBER_EL_NAME,
  'OTHER-PAYER': NM1_OTHER_PAYER_EL_NAME,
  'OTHER-PAYER-ATTENDING': NM1_ATTENDING_EL_NAME,
  'OTHER-PAYER-OPERATING': NM1_OPERATING_EL_NAME,
  'OTHER-PAYER-OTHER-OPERATING': NM1_OTHER_OPERATING_EL_NAME,
  'OTHER-PAYER-SERVICE-FACILITY': NM1_SERVICE_FACILITY_EL_NAME,
  'OTHER-PAYER-RENDERING': NM1_RENDERING_EL_NAME,
  'OTHER-PAYER-REFERRING': NM1_REFERRING_EL_NAME,
  'OTHER-PAYER-BILLING-PROVIDER': NM1_BILLING_EL_NAME,
  'LINE-OPERATING': NM1_OPERATING_EL_NAME,
  'LINE-OTHER-OPERATING': NM1_OTHER_OPERATING_EL_NAME,
  'LINE-RENDERING': NM1_RENDERING_EL_NAME,
  'LINE-REFERRING': NM1_REFERRING_EL_NAME,
}

export function describeOpening(box: Box): SegmentDisplay {
  const e = box.opening.elements
  const nm1Names = NM1_EL_NAME_BY_BOX[box.boxType]
  if (nm1Names) {
    const cells: Cell[] = [codeCell(ref('NM1', 1), NM101_EL_NAME, e[0], NM1_ROLE_NAME), codeCell(ref('NM1', 2), NM102_EL_NAME, e[1], NM102_NAME)]
    for (let i = 2; i < e.length; i++) {
      const cellRef = ref('NM1', i + 1)
      const name = nm1Names[i - 2] // each NM1_*_EL_NAME array starts at NM103
      if (i === 7) cells.push(codeCell(cellRef, name, e[7], NM108_NAME))
      else cells.push(vcell(cellRef, name, e[i]))
    }
    return { lead: label('NM1', undefined), cells }
  }

  switch (box.boxType) {
    case 'ISA':
      return { lead: label('ISA', SEGMENT_NAME.ISA), cells: e.map((v, i) => vcell(ref('ISA', i + 1), ISA_EL_NAME[i], v)) }
    case 'GS':
      return {
        lead: label('GS', SEGMENT_NAME.GS),
        cells: [codeCell(ref('GS', 1), GS_EL_NAME[0], e[0], TRANSACTION_TYPE_NAME), ...e.slice(1).map((v, i) => vcell(ref('GS', i + 2), GS_EL_NAME[i + 1], v))],
      }
    case 'ST':
      return {
        lead: label('ST', SEGMENT_NAME.ST),
        cells: [codeCell(ref('ST', 1), ST_EL_NAME[0], e[0], TRANSACTION_TYPE_NAME), vcell(ref('ST', 2), ST_EL_NAME[1], e[1])],
      }
    case 'BILLING-HL':
    case 'SUBSCRIBER-HL':
    case 'PATIENT-HL':
      return {
        lead: label('HL', SEGMENT_NAME.HL),
        cells: e.map((v, i) => (i === 2 ? codeCell(ref('HL', i + 1), HL_EL_NAME[i], v, HL03_NAME) : vcell(ref('HL', i + 1), HL_EL_NAME[i], v))),
      }
    case 'CLAIM': {
      const n = CLM_EL_NAME
      const cells: Cell[] = []
      for (let i = 0; i < e.length; i++) {
        const cellRef = ref('CLM', i + 1)
        if (i === 4) {
          cells.push(cell(cellRef, n[4], clm05Part(e[4])))
          continue
        }
        const name = i < 4 ? n[i] : n[i + 3]
        if (i === 6) cells.push(codeCell(cellRef, name, e[6], CLM07_NAME))
        else if (i === 7) cells.push(codeCell(cellRef, name, e[7], CLM08_NAME))
        else if (i === 8) cells.push(codeCell(cellRef, name, e[8], CLM09_NAME))
        else if (i === 19) cells.push(codeCell(cellRef, name, e[19])) // CLM20 Delay Reason Code — value list only
        else cells.push(vcell(cellRef, name, e[i], i === 1)) // CLM02 Total Claim Charge Amount, an amount
      }
      return { lead: label('CLM', SEGMENT_NAME.CLM), cells }
    }
    case 'OTHER-SUBSCRIBER':
      return { lead: label('SBR', SEGMENT_NAME.SBR), cells: sbrCells(e) }
    case 'LINE':
      return { lead: label('LX', SEGMENT_NAME.LX), cells: [vcell(ref('LX', 1), LX01_EL_NAME, e[0])] }
    case 'DRUG': {
      const n = LIN_EL_NAME
      const cells: Cell[] = [vcell(ref('LIN', 1), n[0], e[0])]
      if (e[1] !== undefined) cells.push(codeCell(ref('LIN', 2), n[1], e[1], LIN02_NAME))
      if (e[2] !== undefined) cells.push(vcell(ref('LIN', 3), n[2], e[2]))
      return { lead: label('LIN', SEGMENT_NAME.LIN), cells }
    }
    case 'LINE-ADJUDICATION': {
      const n = SVD_EL_NAME
      const cells: Cell[] = [vcell(ref('SVD', 1), n[0], e[0]), vcell(ref('SVD', 2), n[1], e[1], true), cell(ref('SVD', 3), n[2], qualifierCodePart(e[2]))]
      for (let i = 3; i < e.length; i++) cells.push(vcell(ref('SVD', i + 1), n[i + 8], e[i]))
      return { lead: label('SVD', SEGMENT_NAME.SVD), cells }
    }
    default:
      // Unreachable for a well-formed tree (every BoxType is handled above
      // or through NM1_EL_NAME_BY_BOX); keeps the function total.
      return { lead: label(box.opening.id, SEGMENT_NAME[box.opening.id]), cells: e.map((v, i) => vcell(`${box.opening.id}${pad(i + 1)}`, undefined, v)) }
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

// REF's resolved meaning is entirely host-dependent (the same REF segment
// shape carries a different qualifier set at every loop it sits in).
function refName(host: BoxType, qualifier: string): string | undefined {
  switch (host) {
    case 'BILLING-PROVIDER':
      return qualifier === 'EI' ? REF_BILLING_TAX_NAME : undefined
    case 'PAY-TO-PLAN':
      return REF_PAYTOPLAN_NAME[qualifier]
    case 'SUBSCRIBER':
      return REF_SUBSCRIBER_NAME[qualifier]
    case 'PAYER':
      return REF_PAYER_NAME[qualifier]
    case 'PATIENT':
      return REF_PATIENT_NAME[qualifier]
    case 'CLAIM':
      return REF_CLAIM_NAME[qualifier]
    case 'ATTENDING':
    case 'OPERATING':
    case 'OTHER-OPERATING':
    case 'RENDERING':
    case 'SERVICE-FACILITY':
    case 'REFERRING':
    case 'OTHER-PAYER-ATTENDING':
    case 'OTHER-PAYER-OPERATING':
    case 'OTHER-PAYER-OTHER-OPERATING':
    case 'OTHER-PAYER-SERVICE-FACILITY':
    case 'OTHER-PAYER-RENDERING':
    case 'OTHER-PAYER-REFERRING':
    case 'OTHER-PAYER-BILLING-PROVIDER':
    case 'LINE-OPERATING':
    case 'LINE-OTHER-OPERATING':
    case 'LINE-RENDERING':
    case 'LINE-REFERRING':
      return PROVIDER_REF_NAME[qualifier]
    case 'OTHER-SUBSCRIBER-NAME':
      return REF_OTHER_SUBSCRIBER_NAME[qualifier]
    case 'OTHER-PAYER':
      return REF_OTHER_PAYER_NAME[qualifier]
    case 'LINE':
      return REF_LINE_NAME[qualifier]
    case 'DRUG':
      return REF_PRESCRIPTION_NAME[qualifier]
    default:
      return undefined
  }
}

function dtpName(host: BoxType, qualifier: string): string | undefined {
  if (host === 'CLAIM') return DTP_CLAIM_NAME[qualifier]
  if (host === 'LINE' && qualifier === '472') return DTP472_NAME
  if ((host === 'OTHER-PAYER' || host === 'LINE-ADJUDICATION') && qualifier === '573') return DTP_OTHER_PAYER_NAME['573']
  return undefined
}

function nteName(host: BoxType, code: string): string | undefined {
  if (host === 'CLAIM') return code === 'ADD' ? NTE_BILLING_NAME : NTE_CLAIM_NAME
  if (host === 'LINE') return NTE_LINE_NAME
  return undefined
}

// Plain rows inside a box: `host` is that box's type, which disambiguates a
// segment ID whose name depends on which loop it sits in (REF, DTP, NTE,
// AMT, PRV all carry different public names at different levels).
export function describeRow(segment: Segment, host: BoxType): SegmentDisplay {
  const e = segment.elements
  const id = segment.id

  switch (id) {
    case 'BHT': {
      const n = BHT_EL_NAME
      return {
        lead: label('BHT', SEGMENT_NAME.BHT),
        cells: e.map((v, i) => vcell(ref('BHT', i + 1), n[i], v)),
      }
    }
    case 'PER': {
      const cells: Cell[] = [codeCell(ref('PER', 1), PER01_EL_NAME, e[0], PER01_NAME)]
      if (e[1] !== undefined) cells.push(vcell(ref('PER', 2), PER_EL_NAME[0], e[1]))
      cells.push(...perCommunicationCells(e))
      return { lead: label('PER', undefined), cells }
    }
    case 'PRV': {
      const n = PRV_EL_NAME
      const cells: Cell[] = [codeCell(ref('PRV', 1), n[0], e[0], PRV01_NAME), codeCell(ref('PRV', 2), n[1], e[1], PRV02_NAME)]
      if (e[2] !== undefined) cells.push(vcell(ref('PRV', 3), n[2], e[2]))
      return { lead: label('PRV', SEGMENT_NAME.PRV), cells }
    }
    case 'CUR':
      return {
        lead: label('CUR', SEGMENT_NAME.CUR),
        cells: e.map((v, i) => (i === 0 ? codeCell(ref('CUR', 1), CUR_EL_NAME[0], v, CUR01_NAME) : vcell(ref('CUR', i + 1), CUR_EL_NAME[i], v))),
      }
    case 'N3':
      return { lead: label('N3', undefined), cells: e.map((v, i) => vcell(ref('N3', i + 1), undefined, v)) }
    case 'N4':
      return { lead: label('N4', undefined), cells: e.map((v, i) => vcell(ref('N4', i + 1), undefined, v)) }
    case 'DMG':
      return {
        lead: label('DMG', SEGMENT_NAME.DMG),
        cells: [
          codeCell(ref('DMG', 1), DMG_EL_NAME[0], e[0], DMG01_NAME),
          vcell(ref('DMG', 2), DMG_EL_NAME[1], e[1]),
          codeCell(ref('DMG', 3), DMG_EL_NAME[2], e[2], DMG03_NAME),
        ],
      }
    case 'SBR':
      return { lead: label('SBR', SEGMENT_NAME.SBR), cells: sbrCells(e) }
    case 'PAT':
      return {
        lead: label('PAT', SEGMENT_NAME.PAT),
        cells: e.map((v, i) => (i === 0 ? codeCell(ref('PAT', 1), PAT_EL_NAME[0], v, RELATIONSHIP_NAME) : vcell(ref('PAT', i + 1), PAT_EL_NAME[i], v))),
      }
    case 'REF':
      return {
        lead: label('REF', undefined),
        cells: [codeCellNamed(ref('REF', 1), undefined, e[0], refName(host, e[0])), vcell(ref('REF', 2), undefined, e[1])],
      }
    case 'DTP':
      return {
        lead: label('DTP', undefined),
        cells: [codeCellNamed(ref('DTP', 1), DTP01_EL_NAME, e[0], dtpName(host, e[0])), vcell(ref('DTP', 2), undefined, e[1])],
      }
    case 'CL1':
      return { lead: label('CL1', SEGMENT_NAME.CL1), cells: e.map((v, i) => codeCell(ref('CL1', i + 1), CL1_EL_NAME[i], v)) }
    case 'PWK':
      // PWK01 (Report Type Code) and PWK02 (Report Transmission Code) each
      // carry a value list only (claim-level and line-level sets differ
      // slightly — PWK01_VALUES_CLAIM/LINE, PWK02_VALUES_CLAIM/LINE), no
      // per-code text in either reviewed source, so both render unresolved
      // regardless of host.
      return {
        lead: label('PWK', SEGMENT_NAME.PWK),
        cells: e.map((v, i) => (i === 0 || i === 1 ? codeCell(ref('PWK', i + 1), PWK_EL_NAME[i], v) : vcell(ref('PWK', i + 1), PWK_EL_NAME[i], v))),
      }
    case 'CN1':
      return {
        lead: label('CN1', SEGMENT_NAME.CN1),
        cells: e.map((v, i) => (i === 0 ? codeCell(ref('CN1', 1), CN1_EL_NAME[0], v) : vcell(ref('CN1', i + 1), CN1_EL_NAME[i], v))),
      }
    case 'AMT': {
      const name = AMT01_NAME[e[0]]
      return {
        lead: label('AMT', undefined),
        cells: [codeCell(ref('AMT', 1), undefined, e[0], AMT01_NAME), vcell(ref('AMT', 2), name, e[1], true)],
      }
    }
    case 'K3':
      return { lead: label('K3', SEGMENT_NAME.K3), cells: e.map((v, i) => vcell(ref('K3', i + 1), K3_EL_NAME[i], v)) }
    case 'NTE': {
      const name = nteName(host, e[0])
      return {
        lead: label('NTE', name),
        cells: [codeCell(ref('NTE', 1), NTE_EL_NAME[0], e[0]), vcell(ref('NTE', 2), NTE_EL_NAME[1], e[1])],
      }
    }
    case 'CRC': {
      const n = CRC_EL_NAME
      return {
        lead: label('CRC', SEGMENT_NAME.CRC),
        cells: e.map((v, i) => (i === 0 ? codeCell(ref('CRC', 1), n[0], v, CRC01_NAME) : i === 1 ? codeCell(ref('CRC', 2), n[1], v, CRC02_NAME) : codeCell(ref('CRC', i + 1), n[i], v))),
      }
    }
    case 'HI':
      return { lead: label('HI', undefined), cells: hiCompositeCells(e) }
    case 'HCP':
      return {
        lead: label('HCP', SEGMENT_NAME.HCP),
        cells: e.map((v, i) => (i === 0 ? codeCell(ref('HCP', 1), HCP_EL_NAME[0], v, HCP01_NAME) : vcell(ref('HCP', i + 1), HCP_EL_NAME[i], v))),
      }
    case 'SV2': {
      const n = SV2_EL_NAME
      const cells: Cell[] = []
      for (let i = 0; i < e.length; i++) {
        const cellRef = ref('SV2', i + 1)
        if (i === 1) {
          cells.push(cell(cellRef, n[1], qualifierCodePart(e[1])))
          continue
        }
        const name = i === 0 ? n[0] : n[i + 8]
        cells.push(vcell(cellRef, name, e[i], i === 2 || i === 6))
      }
      return { lead: label('SV2', SEGMENT_NAME.SV2), cells }
    }
    case 'CTP':
      return { lead: label('CTP', SEGMENT_NAME.CTP), cells: e.map((v, i) => vcell(ref('CTP', i + 1), CTP_EL_NAME[i], v)) }
    case 'OI':
      return {
        lead: label('OI', SEGMENT_NAME.OI),
        cells: e.map((v, i) => (i === 2 ? codeCell(ref('OI', 3), OI_EL_NAME[2], v, OI03_NAME) : i === 5 ? codeCell(ref('OI', 6), OI_EL_NAME[5], v, OI06_NAME) : vcell(ref('OI', i + 1), OI_EL_NAME[i], v))),
      }
    case 'MIA':
      return { lead: label('MIA', SEGMENT_NAME.MIA), cells: miaCells(e) }
    case 'MOA':
      return {
        lead: label('MOA', SEGMENT_NAME.MOA),
        cells: [vcell(ref('MOA', 1), MOA_EL_NAME[0], e[0]), vcell(ref('MOA', 2), MOA_EL_NAME[1], e[1], true), ...e.slice(2).map((code, i) => cell(ref('MOA', i + 3), MOA_EL_NAME[i + 2], rarcLabel(code)))],
      }
    case 'CAS': {
      const name = host === 'OTHER-SUBSCRIBER' ? CAS_CLAIM_NAME : host === 'LINE-ADJUDICATION' ? CAS_LINE_NAME : undefined
      const group = codeCell(ref('CAS', 1), CAS01_EL_NAME, e[0], CAS_GROUP_NAME)
      group.marker = CAS_GROUP_DOT[e[0]]
      return { lead: label('CAS', name), cells: [group, ...casTripletCells(e)] }
    }
    default:
      return { lead: label(id, SEGMENT_NAME[id]), cells: e.map((v, i) => vcell(`${id}${pad(i + 1)}`, undefined, v)) }
  }
}
