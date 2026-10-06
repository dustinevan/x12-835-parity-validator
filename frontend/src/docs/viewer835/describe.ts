// Turns a parsed segment or box opening/closing into the row of display
// cells the page renders: one cell per file element, its element reference
// and its value on the first line, its public name beneath (Cell.name;
// undefined means the cell has no name line; viewer/Parts.tsx draws it).
// Every name comes from frontend/src/data/x12Names835.ts, which is built only
// from the spec's tables; a coded value with no entry there (the spec's
// "Open" list) renders as its raw code alone, with no name on hover — see the
// `label()` helper below.
//
// Two kinds of name table are in play, and they answer different
// questions. A "qualifier code" table (e.g. N1_NAME, DTM_NAME, HEADER_REF_NAME)
// names what a CODE means — N1_NAME.PR is "Payer Identification", the
// meaning of the code "PR". An "element position" table (e.g. N1_PR_EL_NAME,
// DTM01_EL_NAME) names the ELEMENT ITSELF — N1_PR_EL_NAME[0] is "Entity
// Identifier Code", N101's own name regardless of what code it holds. A
// cell's name line always comes from the second kind; its value, when the
// element is itself a code, is resolved through the first kind.
//
// A segment whose row-level identity depends on a qualifier value (N1, NM1,
// PER, REF outside the payer/payee loops, TRN, DTM, LQ) leads its row with
// its bare segment ID only — the resolved meaning lives entirely in the
// qualifier's own cell, never in both places (docs/design/835-viewer.md,
// "no duplicate meaning between a row's leading name and a cell"). A segment
// whose name is fixed regardless of any qualifier (ISA, GS, ST, BPR, CLP,
// SVC, MOA, MIA, PLB, SE, GE, IEA, RDM, CUR) or that varies only by which
// loop holds it, not by any value inside it (N3, N4, REF inside the
// payer/payee loops, CAS, AMT, QTY) leads with that true name instead.
import {
  AMT01_EL_NAME,
  BPR03_NAME,
  BPR_EL_NAME,
  CARC,
  CAS01_EL_NAME,
  CAS_AMOUNT_EL_NAME,
  CAS_GROUP_NAME,
  CAS_QTY_EL_NAME,
  CAS_REASON_EL_NAME,
  CLAIM_AMT02_EL_NAME,
  CLAIM_AMT_NAME,
  CLAIM_CAS_NAME,
  CLAIM_PER_NAME,
  CLAIM_QTY_NAME,
  CLAIM_REF_NAME,
  CLP_EL_NAME,
  DTM01_EL_NAME,
  DTM02_EL_NAME_BY_QUALIFIER,
  DTM_NAME,
  GE_EL_NAME,
  GS_EL_NAME,
  HEADER_REF_NAME,
  ID_QUALIFIER_NAME,
  IEA_EL_NAME,
  ISA_EL_NAME,
  LQ_EL_NAME,
  LQ_NAME,
  LX01_EL_NAME,
  MIA_EL_NAME,
  MOA_EL_NAME,
  N1_NAME,
  N1_PE_EL_NAME,
  N1_PR_EL_NAME,
  N3_PE_EL_NAME,
  N3_PR_EL_NAME,
  N4_PE_EL_NAME,
  N4_PR_EL_NAME,
  NM101_EL_NAME,
  NM102_EL_NAME,
  NM1_74_EL_NAME,
  NM1_82_EL_NAME,
  NM1_NAME,
  NM1_PR_EL_NAME,
  NM1_QC_EL_NAME,
  NM1_TT_EL_NAME,
  PAYEE_LOOP_NAME,
  PAYER_LOOP_NAME,
  PAYER_PER_NAME,
  PER01_EL_NAME,
  PER_CLAIM_EL_NAME,
  PER_PAYER_EL_NAME,
  PLB01_EL_NAME,
  PLB02_EL_NAME,
  PLB_AMOUNT_EL_NAME,
  PLB_COMPOSITE_EL_NAME,
  PLB_REASON_NAME,
  QTY01_EL_NAME,
  RARC,
  RDM_EL_NAME,
  REF_CLAIM_EL_NAME,
  REF_HEADER_EL_NAME,
  REF_LINE_EL_NAME,
  REF_PAYEE_EL_NAME,
  REF_PAYER_EL_NAME,
  SE_EL_NAME,
  SEGMENT_NAME,
  SERVICE_AMT02_EL_NAME,
  SERVICE_AMT_NAME,
  SERVICE_CAS_NAME,
  SERVICE_REF_NAME,
  ST_EL_NAME,
  SVC_EL_NAME,
  SVC_QUALIFIER_NAME,
  TRANSACTION_TYPE_NAME,
  TRN_EL_NAME,
  TRN_NAME,
  TS2_EL_NAME,
  TS3_EL_NAME,
} from '../../data/x12Names835'
import { cell, label, lcell, pad, ref, value, vcell, type Cell, type LabelPart, type Part, type SegmentDisplay } from '../viewer/display'
import { CAS_GROUP_DOT } from './colors'
import type { Box, Segment } from './parse'

export type { Cell, JoinedPart, LabelPart, Part, SegmentDisplay, ValuePart } from '../viewer/display'
export { describeCompact } from '../viewer/display'

function carcLabel(code: string): LabelPart {
  const entry = CARC[code]
  return label(code, entry?.short, { full: entry?.full })
}
function rarcLabel(code: string): LabelPart {
  const entry = RARC[code]
  return label(code, entry?.short, { full: entry?.full })
}

// SVC01 and SVC06 are both the same composite shape: a code-system
// qualifier, then the procedure code (and possibly modifiers, not present in
// this page's examples) — SVC01 is what was adjudicated, SVC06 is the
// original, as-billed code. Rendered as one cell (the composite's two
// components joined tight), not split across two cells.
function svc01Part(raw: string): Part {
  const [qualifier, code] = raw.split(':')
  const qualifierPart = label(qualifier, SVC_QUALIFIER_NAME[qualifier])
  if (!code) return qualifierPart
  return { kind: 'joined', parts: [qualifierPart, value(code)] }
}

// CAS holds up to six reason/amount/quantity triplets after its group code;
// each triplet gets three cells (the quantity cell only when the file sent one).
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

// MIA holds five remark-code positions among otherwise plain counts/amounts:
// MIA05 (index 4) and MIA20-MIA23 (indices 19-22). Every other position is a
// plain value, per the page's rule that only coded values get a public name.
const MIA_REMARK_INDEXES = new Set([4, 19, 20, 21, 22])
function miaCells(e: string[]): Cell[] {
  return e.map((v, i) => {
    const cellRef = ref('MIA', i + 1)
    const name = MIA_EL_NAME[i]
    return MIA_REMARK_INDEXES.has(i) && v ? cell(cellRef, name, rarcLabel(v)) : vcell(cellRef, name, v)
  })
}

// PLB repeats a reason:identifier composite + amount up to six times after
// PLB01 (provider identifier) and PLB02 (fiscal period date), both plain
// values. The reason code is named from the public Provider Adjustment
// Reason Code list; the identifier is free text, joined tight to its reason
// the same way a composite element's sub-parts are elsewhere on this page.
function plbAdjustmentCells(e: string[]): Cell[] {
  const cells: Cell[] = []
  for (let i = 2; i + 1 < e.length; i += 2) {
    const composite = e[i]
    const amount = e[i + 1]
    if (composite === undefined) break
    const [reason, identifier] = composite.split(':')
    const reasonPart = label(reason, PLB_REASON_NAME[reason])
    const part: Part = identifier !== undefined ? { kind: 'joined', parts: [reasonPart, value(identifier)] } : reasonPart
    cells.push(cell(ref('PLB', i + 1), PLB_COMPOSITE_EL_NAME, part))
    cells.push(vcell(ref('PLB', i + 2), PLB_AMOUNT_EL_NAME, amount, true))
  }
  return cells
}

// A PER segment's qualifier+value pairs beyond PER01/02 (PER03/04,
// PER05/06, PER07/08): a communication-number qualifier (Open list, shown
// unresolved) followed by the communication number itself.
function perCommunicationCells(e: string[], names: string[]): Cell[] {
  const cells: Cell[] = []
  for (let i = 2; i < e.length; i += 2) {
    const nameIdx = i - 1 // names[] starts at PER02
    cells.push(lcell(ref('PER', i + 1), names[nameIdx], e[i], undefined))
    if (e[i + 1] !== undefined) cells.push(vcell(ref('PER', i + 2), names[nameIdx + 1], e[i + 1]))
  }
  return cells
}

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
        cells: [lcell(ref('ST', 1), ST_EL_NAME[0], e[0], TRANSACTION_TYPE_NAME[e[0]]), vcell(ref('ST', 2), ST_EL_NAME[1], e[1])],
      }
    case 'N1-PR': {
      const n = N1_PR_EL_NAME
      const cells: Cell[] = [lcell(ref('N1', 1), n[0], e[0], N1_NAME[e[0]]), vcell(ref('N1', 2), n[1], e[1])]
      if (e[2] !== undefined) cells.push(lcell(ref('N1', 3), n[2], e[2], ID_QUALIFIER_NAME[e[2]]))
      if (e[3] !== undefined) cells.push(vcell(ref('N1', 4), n[3], e[3]))
      return { lead: label('N1', undefined), cells }
    }
    case 'N1-PE': {
      const n = N1_PE_EL_NAME
      const cells: Cell[] = [lcell(ref('N1', 1), n[0], e[0], N1_NAME[e[0]]), vcell(ref('N1', 2), n[1], e[1])]
      if (e[2] !== undefined) cells.push(lcell(ref('N1', 3), n[2], e[2], ID_QUALIFIER_NAME[e[2]]))
      if (e[3] !== undefined) cells.push(vcell(ref('N1', 4), n[3], e[3]))
      return { lead: label('N1', undefined), cells }
    }
    case 'LX':
      return { lead: label('LX', SEGMENT_NAME.LX), cells: [vcell(ref('LX', 1), LX01_EL_NAME, e[0])] }
    case 'CLP': {
      // CLP01-07 are always present; CLP08-14 (facility code, frequency,
      // CLP10 "not used", DRG code, DRG weight, discharge fraction, CLP14
      // "not used") are rendered whenever the file actually sends them —
      // a trailing element the file never sends is absent from the array.
      const n = CLP_EL_NAME
      const cells: Cell[] = [
        vcell(ref('CLP', 1), n[0], e[0]),
        lcell(ref('CLP', 2), n[1], e[1], undefined), // CLP02 status — Open list
        vcell(ref('CLP', 3), n[2], e[2], true),
        vcell(ref('CLP', 4), n[3], e[3], true),
        vcell(ref('CLP', 5), n[4], e[4], true),
        lcell(ref('CLP', 6), n[5], e[5], undefined), // CLP06 filing indicator — Open list
        vcell(ref('CLP', 7), n[6], e[6]),
      ]
      for (let i = 7; i < e.length && i < 14; i++) {
        const number = i === 11 || i === 12 // CLP12 DRG weight, CLP13 discharge fraction
        cells.push(vcell(ref('CLP', i + 1), n[i], e[i], number))
      }
      return { lead: label('CLP', SEGMENT_NAME.CLP), cells }
    }
    case 'SVC': {
      const n = SVC_EL_NAME
      const cells: Cell[] = [
        cell(ref('SVC', 1), n[0], svc01Part(e[0])),
        vcell(ref('SVC', 2), n[1], e[1], true),
        vcell(ref('SVC', 3), n[2], e[2], true),
      ]
      if (e[3] !== undefined) cells.push(vcell(ref('SVC', 4), n[3], e[3]))
      if (e[4] !== undefined) cells.push(vcell(ref('SVC', 5), n[4], e[4]))
      if (e[5]) cells.push(cell(ref('SVC', 6), n[5], svc01Part(e[5]))) // SVC06: original (as-billed) procedure composite
      if (e[6] !== undefined) cells.push(vcell(ref('SVC', 7), n[6], e[6])) // SVC07 original units of service
      return { lead: label('SVC', SEGMENT_NAME.SVC), cells }
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

// Plain rows inside a box: `host` is that box's type, which disambiguates a
// segment ID whose name depends on which loop it sits in (N3/N4/REF/PER/
// CAS/AMT/QTY all carry different public names at different levels).
export function describeRow(segment: Segment, host: Box['boxType']): SegmentDisplay {
  const e = segment.elements
  const id = segment.id

  if (id === 'N3' || id === 'N4' || id === 'REF') {
    const loopTable = host === 'N1-PR' ? PAYER_LOOP_NAME : host === 'N1-PE' ? PAYEE_LOOP_NAME : undefined
    if (loopTable) {
      if (id === 'REF') {
        const n = host === 'N1-PR' ? REF_PAYER_EL_NAME : REF_PAYEE_EL_NAME
        return {
          lead: label('REF', loopTable.REF),
          cells: [lcell(ref('REF', 1), n[0], e[0], undefined), vcell(ref('REF', 2), n[1], e[1])],
        }
      }
      const n = host === 'N1-PR' ? (id === 'N3' ? N3_PR_EL_NAME : N4_PR_EL_NAME) : id === 'N3' ? N3_PE_EL_NAME : N4_PE_EL_NAME
      return { lead: label(id, loopTable[id]), cells: e.map((v, i) => vcell(`${id}${pad(i + 1)}`, n[i], v)) }
    }
    if (id === 'REF') {
      // Header REF (EV, F2) sits directly under ST; claim REF (EA, F8, ...)
      // and line REF (6R, LU, HPI, 0K, ...) sit under CLP and SVC. Each is a
      // qualifier-driven name: it leads bare, and the qualifier's own cell
      // carries the resolved meaning.
      const resolvedName = host === 'ST' ? HEADER_REF_NAME[e[0]] : host === 'CLP' ? CLAIM_REF_NAME[e[0]] : host === 'SVC' ? SERVICE_REF_NAME[e[0]] : undefined
      const n = host === 'ST' ? REF_HEADER_EL_NAME[e[0]] : host === 'CLP' ? REF_CLAIM_EL_NAME : host === 'SVC' ? REF_LINE_EL_NAME[e[0]] : undefined
      return {
        lead: label('REF', undefined),
        cells: [lcell(ref('REF', 1), n?.[0], e[0], resolvedName), vcell(ref('REF', 2), n?.[1], e[1])],
      }
    }
  }

  switch (id) {
    case 'BPR': {
      const n = BPR_EL_NAME
      return {
        lead: label('BPR', SEGMENT_NAME.BPR),
        cells: [
          lcell(ref('BPR', 1), n[0], e[0], undefined), // BPR01 — Open list
          vcell(ref('BPR', 2), n[1], e[1], true),
          lcell(ref('BPR', 3), n[2], e[2], BPR03_NAME[e[2]]),
          lcell(ref('BPR', 4), n[3], e[3], undefined), // BPR04 — Open list
          lcell(ref('BPR', 5), n[4], e[4], undefined), // BPR05 — Open list
          ...e.slice(5, 15).map((v, i) => vcell(ref('BPR', i + 6), n[i + 5], v)),
          ...(e[15] !== undefined ? [vcell(ref('BPR', 16), n[15], e[15])] : []),
        ],
      }
    }
    case 'TRN': {
      const n = TRN_EL_NAME
      const cells = [lcell(ref('TRN', 1), n[0], e[0], TRN_NAME[e[0]]), vcell(ref('TRN', 2), n[1], e[1])]
      if (e[2] !== undefined) cells.push(vcell(ref('TRN', 3), n[2], e[2]))
      if (e[3] !== undefined) cells.push(vcell(ref('TRN', 4), n[3], e[3]))
      return { lead: label('TRN', undefined), cells }
    }
    case 'DTM':
      return {
        lead: label('DTM', undefined),
        cells: [lcell(ref('DTM', 1), DTM01_EL_NAME, e[0], DTM_NAME[e[0]]), vcell(ref('DTM', 2), DTM02_EL_NAME_BY_QUALIFIER[e[0]], e[1])],
      }
    case 'PER': {
      // Payer-loop PER (CX/BL/1C) and claim-loop PER (CX) carry different names.
      const resolveTable = host === 'N1-PR' ? PAYER_PER_NAME : host === 'CLP' ? CLAIM_PER_NAME : undefined
      const nameTable = host === 'N1-PR' ? PER_PAYER_EL_NAME[e[0]] : host === 'CLP' ? PER_CLAIM_EL_NAME : undefined
      const cells: Cell[] = [lcell(ref('PER', 1), PER01_EL_NAME, e[0], resolveTable?.[e[0]])]
      if (e[1] !== undefined) cells.push(vcell(ref('PER', 2), nameTable?.[0], e[1]))
      cells.push(...perCommunicationCells(e, nameTable ?? []))
      return { lead: label('PER', undefined), cells }
    }
    case 'NM1': {
      const roleArr = e[0] === 'QC' ? NM1_QC_EL_NAME : e[0] === '74' ? NM1_74_EL_NAME : e[0] === '82' ? NM1_82_EL_NAME : e[0] === 'TT' ? NM1_TT_EL_NAME : e[0] === 'PR' ? NM1_PR_EL_NAME : undefined
      const cells: Cell[] = [
        lcell(ref('NM1', 1), NM101_EL_NAME, e[0], NM1_NAME[e[0]]),
        lcell(ref('NM1', 2), NM102_EL_NAME, e[1], undefined), // NM102 entity type — Open list
      ]
      for (let i = 2; i < e.length; i++) {
        const cellRef = ref('NM1', i + 1)
        const name = roleArr?.[i - 2]
        if (i === 7) cells.push(lcell(cellRef, name, e[7], ID_QUALIFIER_NAME[e[7]]))
        else cells.push(vcell(cellRef, name, e[i]))
      }
      return { lead: label('NM1', undefined), cells }
    }
    case 'MOA':
      return {
        lead: label('MOA', SEGMENT_NAME.MOA),
        cells: [
          vcell(ref('MOA', 1), MOA_EL_NAME[0], e[0]),
          vcell(ref('MOA', 2), MOA_EL_NAME[1], e[1], true),
          ...e.slice(2).map((code, i) => cell(ref('MOA', i + 3), MOA_EL_NAME[i + 2], rarcLabel(code))),
        ],
      }
    case 'MIA':
      return { lead: label('MIA', SEGMENT_NAME.MIA), cells: miaCells(e) }
    case 'CAS': {
      const name = host === 'CLP' ? CLAIM_CAS_NAME : SERVICE_CAS_NAME
      const group = lcell(ref('CAS', 1), CAS01_EL_NAME, e[0], CAS_GROUP_NAME[e[0]])
      group.marker = CAS_GROUP_DOT[e[0]]
      return { lead: label('CAS', name), cells: [group, ...casTripletCells(e)] }
    }
    case 'AMT': {
      const name = host === 'CLP' ? CLAIM_AMT_NAME : host === 'SVC' ? SERVICE_AMT_NAME : undefined
      const el02Name = host === 'CLP' ? CLAIM_AMT02_EL_NAME : host === 'SVC' ? SERVICE_AMT02_EL_NAME : undefined
      return {
        lead: label('AMT', name),
        cells: [lcell(ref('AMT', 1), AMT01_EL_NAME, e[0], undefined), vcell(ref('AMT', 2), el02Name, e[1], true)],
      }
    }
    case 'QTY': {
      const name = host === 'CLP' ? CLAIM_QTY_NAME : undefined
      return {
        lead: label('QTY', name),
        cells: [lcell(ref('QTY', 1), QTY01_EL_NAME, e[0], undefined), vcell(ref('QTY', 2), name, e[1])],
      }
    }
    case 'LQ':
      return {
        lead: label('LQ', undefined),
        cells: [lcell(ref('LQ', 1), LQ_EL_NAME[0], e[0], LQ_NAME[e[0]]), cell(ref('LQ', 2), LQ_EL_NAME[1], rarcLabel(e[1]))],
      }
    case 'PLB':
      return {
        lead: label('PLB', SEGMENT_NAME.PLB),
        cells: [vcell(ref('PLB', 1), PLB01_EL_NAME, e[0]), vcell(ref('PLB', 2), PLB02_EL_NAME, e[1]), ...plbAdjustmentCells(e)],
      }
    case 'RDM': {
      const n = RDM_EL_NAME
      const cells = [lcell(ref('RDM', 1), n[0], e[0], undefined)] // RDM01 — Open list
      if (e[1] !== undefined) cells.push(vcell(ref('RDM', 2), n[1], e[1]))
      if (e[2] !== undefined) cells.push(vcell(ref('RDM', 3), n[2], e[2]))
      return { lead: label('RDM', SEGMENT_NAME.RDM), cells }
    }
    case 'TS3':
      return { lead: label('TS3', SEGMENT_NAME.TS3), cells: e.map((v, i) => vcell(ref('TS3', i + 1), TS3_EL_NAME[i], v)) }
    case 'TS2':
      return { lead: label('TS2', SEGMENT_NAME.TS2), cells: e.map((v, i) => vcell(ref('TS2', i + 1), TS2_EL_NAME[i], v)) }
    default:
      return { lead: label(id, SEGMENT_NAME[id]), cells: e.map((v, i) => vcell(`${id}${pad(i + 1)}`, undefined, v)) }
  }
}
