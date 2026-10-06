// The 837P's own loop grammar, built from the loop-rule table in
// frontend/src/docs/viewer837p/rules.ts (which segment, and which qualifier
// where the segment ID is reused, opens which box) plus this file's
// knowledge of each box's own plain-segment content and nesting order,
// grounded in the 24 worked examples under frontend/src/data/viewer837p/files/
// and docs/design/837p-examples.md. The generic tree shapes and the walks
// over an already-built tree are transaction-agnostic and live in
// frontend/src/docs/viewer/tree.ts.
//
// The whole transaction is one descending chain, not a flat list of
// envelope-level loops the way the 835 is (ISA ⊃ GS ⊃ ST ⊃ several
// loops side by side): ISA ⊃ GS ⊃ ST ⊃ BHT, submitter, receiver, then
// billing provider HL ⊃ subscriber HL ⊃ claim CLM ⊃ (claim-level provider
// loops, other subscribers, service line LX ⊃ (drug, line providers, line
// adjudication, forms)) — each level's own parse function consumes
// everything nested inside it, all the way to SE.
import * as tree from '../viewer/tree'
import { splitX12 } from '../viewer/x12split'

export type { Segment } from '../viewer/tree'
export { allBoxPaths, ancestorPaths, childBoxCount, envelopePaths, findBox, seedOpenPaths } from '../viewer/tree'

import { CHILDREN, LOOP_RULES, type BoxType } from './rules'

export type { BoxType } from './rules'
export type Box = tree.Box<BoxType>
export type Node = tree.Node<BoxType>

// The file's segments, split on the delimiters its own ISA declares (element
// separator, component separator ISA16, segment terminator), so a file need
// not use `* : ~`, break its lines, or avoid CRLF — frontend/src/docs/viewer/
// x12split.ts. Throws when the text has no readable ISA header.
export function splitSegments(raw: string): tree.Segment[] {
  return splitX12(raw).segments
}

const ENVELOPE_CLOSERS = new Set(['SE', 'GE', 'IEA'])

// Builds a qualifier-code → box-type dispatch map straight from the
// loop-rule table, scoped to one parent's valid children — e.g. CLM's own
// NM1 dispatch (DN→REFERRING, 82→RENDERING, 77→SERVICE_FACILITY, …) versus
// LX's (the same four-plus codes, to the line-level box types instead).
// This is what makes the parser table-driven rather than a hardcoded
// per-level switch: add a row to rules.ts and both the matching parent's
// dispatch and the tint/describe lookups pick it up.
function nm1Dispatch(parentType: BoxType): Record<string, BoxType> {
  const dispatch: Record<string, BoxType> = {}
  for (const childType of CHILDREN[parentType] ?? []) {
    const rule = LOOP_RULES.find((r) => r.boxType === childType && r.segmentId === 'NM1')
    if (rule) for (const code of rule.qualifierCodes ?? []) dispatch[code] = childType
  }
  return dispatch
}

const CLM_NM1 = nm1Dispatch('CLM')
const LX_NM1 = nm1Dispatch('LX')

// Ambulance pick-up/drop-off (2310E/F, 2420G/H) carry only a bare name plus
// address — no REF or PER, per Transmittal 907's own "(Not Used)" marking
// of every other NM1 sub-element in those two loops (x12Names837p.ts,
// AMBULANCE_PICKUP_NM1_EL_NAME/AMBULANCE_DROPOFF_NM1_EL_NAME). Every other
// provider loop this file draws carries address plus REF/PER.
const AMBULANCE_TYPES = new Set<BoxType>(['AMB_PICKUP', 'AMB_DROPOFF', 'AMB_PICKUP_LINE', 'AMB_DROPOFF_LINE'])
function providerPlainIds(boxType: BoxType): string[] {
  return AMBULANCE_TYPES.has(boxType) ? ['N3', 'N4'] : ['N3', 'N4', 'DMG', 'REF', 'PER']
}

export function parseLoopTree(raw: string): Box {
  const segments = splitSegments(raw)
  let i = 0
  const peek = () => segments[i]
  const take = () => segments[i++]

  // A plain NM1 loop: the name segment, then whichever of N3/N4/DMG/REF/PER
  // the file actually sends, in file order, until something that isn't one
  // of those appears. Covers every leaf provider/payer/subscriber loop this
  // transaction has (2010AA/AB/BA/BB, 2310A-F, 2330A/B, 2420A-H), and, with
  // a different plain-id set, the drug (2410), line adjudication (2430) and
  // forms (2440) loops too — none of those nest a box of their own.
  function parseSimpleLoop(parentPath: string, boxType: BoxType, plainIds: string[]): Box {
    const opening = take()
    const path = `${parentPath}.${boxType}`
    const children: Node[] = []
    while (peek() && plainIds.includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    return { path, boxType, opening, children }
  }

  function parseForms(parentPath: string): Box {
    return parseSimpleLoop(parentPath, 'FORMS', ['FRM'])
  }

  function parseLineAdjudication(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.LINE_ADJUDICATION#${index}`
    const children: Node[] = []
    while (peek() && ['CAS', 'DTP', 'AMT'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    return { path, boxType: 'LINE_ADJUDICATION', opening, children }
  }

  function parseDrug(parentPath: string): Box {
    return parseSimpleLoop(parentPath, 'DRUG', ['CTP', 'REF'])
  }

  // 2400 Service line. Plain claim-line segments (SV1, SV5, CR1, CR3, CRC,
  // PWK, DTP, QTY, MEA, REF, AMT, K3, NTE, PS1, HCP) and the five loops that
  // may nest inside one line (2410 drug, 2420A-H line providers, 2430 line
  // adjudication, 2440 forms) interleave in whatever order the file sends
  // them; ends at the next LX (a sibling line) or at SE.
  function parseLX(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.LX#${index}`
    const children: Node[] = []
    let adjIndex = 0
    while (peek() && !ENVELOPE_CLOSERS.has(peek().id) && peek().id !== 'LX') {
      const seg = peek()
      if (seg.id === 'LIN') {
        children.push({ kind: 'box', box: parseDrug(path) })
        continue
      }
      if (seg.id === 'SVD') {
        children.push({ kind: 'box', box: parseLineAdjudication(path, adjIndex++) })
        continue
      }
      if (seg.id === 'LQ') {
        children.push({ kind: 'box', box: parseForms(path) })
        continue
      }
      if (seg.id === 'NM1') {
        const boxType = LX_NM1[seg.elements[0]]
        if (boxType) {
          children.push({ kind: 'box', box: parseSimpleLoop(path, boxType, providerPlainIds(boxType)) })
          continue
        }
      }
      children.push({ kind: 'segment', segment: take() })
    }
    return { path, boxType: 'LX', opening, children }
  }

  // 2320 Other Subscriber Information, opened by the claim-level SBR (not
  // the 2000B one, which parseSubscriberHL already consumed as a plain
  // row). CAS/AMT/OI/MOA report the other payer's own adjudication; 2330A
  // names that other subscriber, 2330B that other payer.
  function parseOtherSubscriber(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.OTHER_SUBSCRIBER#${index}`
    const children: Node[] = []
    while (peek() && ['CAS', 'AMT', 'OI', 'MOA'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    if (peek()?.id === 'NM1' && peek().elements[0] === 'IL') children.push({ kind: 'box', box: parseSimpleLoop(path, 'OTHER_SUBSCRIBER_NAME', ['N3', 'N4', 'REF']) })
    if (peek()?.id === 'NM1' && peek().elements[0] === 'PR') children.push({ kind: 'box', box: parseSimpleLoop(path, 'OTHER_PAYER', ['N3', 'N4', 'DTP', 'REF']) })
    return { path, boxType: 'OTHER_SUBSCRIBER', opening, children }
  }

  // 2300 Claim. Plain claim-level segments (DTP, PWK, CN1, AMT, REF, K3,
  // NTE, CR1, CR2, CRC, HI, HCP) and the claim's own nested loops (2310A-F
  // providers, 2320 other subscribers, 2400 service lines) interleave in
  // whatever order the file sends them; ends at SE (one claim per
  // subscriber in every example this viewer has, but the loop supports a
  // second CLM — REPEATS includes 'CLM' — the same way LX repeats).
  function parseCLM(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.CLM#${index}`
    const children: Node[] = []
    let lxIndex = 0
    let otherSubIndex = 0
    while (peek() && !ENVELOPE_CLOSERS.has(peek().id) && peek().id !== 'CLM') {
      const seg = peek()
      if (seg.id === 'LX') {
        children.push({ kind: 'box', box: parseLX(path, lxIndex++) })
        continue
      }
      if (seg.id === 'SBR') {
        children.push({ kind: 'box', box: parseOtherSubscriber(path, otherSubIndex++) })
        continue
      }
      if (seg.id === 'NM1') {
        const boxType = CLM_NM1[seg.elements[0]]
        if (boxType) {
          children.push({ kind: 'box', box: parseSimpleLoop(path, boxType, providerPlainIds(boxType)) })
          continue
        }
      }
      children.push({ kind: 'segment', segment: take() }) // DTP, PWK, CN1, AMT, REF, K3, NTE, CR1, CR2, CRC, HI, HCP
    }
    return { path, boxType: 'CLM', opening, children }
  }

  // 2000B Subscriber Hierarchical Level — nested inside Billing Provider HL
  // per the real HL02 parent pointer every example carries, and itself
  // containing the claim(s) (2300 is subordinate to the subscriber, not a
  // sibling of it). The 2000C Patient Hierarchical Level never appears —
  // see rules.ts's header comment.
  function parseSubscriberHL(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.SUBSCRIBER_HL`
    const children: Node[] = []
    if (peek()?.id === 'SBR') children.push({ kind: 'segment', segment: take() })
    if (peek()?.id === 'PAT') children.push({ kind: 'segment', segment: take() })
    if (peek()?.id === 'NM1' && peek().elements[0] === 'IL') children.push({ kind: 'box', box: parseSimpleLoop(path, 'SUBSCRIBER', providerPlainIds('SUBSCRIBER')) })
    if (peek()?.id === 'NM1' && peek().elements[0] === 'PR') children.push({ kind: 'box', box: parseSimpleLoop(path, 'PAYER', ['N3', 'N4', 'REF']) })
    let claimIndex = 0
    while (peek()?.id === 'CLM') children.push({ kind: 'box', box: parseCLM(path, claimIndex++) })
    return { path, boxType: 'SUBSCRIBER_HL', opening, children }
  }

  // 2000A Billing Provider Hierarchical Level.
  function parseBillingHL(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.BILLING_HL`
    const children: Node[] = []
    while (peek() && ['PRV', 'CUR'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    if (peek()?.id === 'NM1' && peek().elements[0] === '85') children.push({ kind: 'box', box: parseSimpleLoop(path, 'BILLING_PROVIDER', providerPlainIds('BILLING_PROVIDER')) })
    if (peek()?.id === 'NM1' && peek().elements[0] === '87') children.push({ kind: 'box', box: parseSimpleLoop(path, 'PAYTO', providerPlainIds('PAYTO')) })
    if (peek()?.id === 'HL') children.push({ kind: 'box', box: parseSubscriberHL(path) })
    return { path, boxType: 'BILLING_HL', opening, children }
  }

  function parseST(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.ST`
    const children: Node[] = []
    while (peek() && peek().id !== 'SE') {
      const seg = peek()
      if (seg.id === 'NM1' && seg.elements[0] === '41') {
        children.push({ kind: 'box', box: parseSimpleLoop(path, 'SUBMITTER', ['PER']) })
        continue
      }
      if (seg.id === 'NM1' && seg.elements[0] === '40') {
        children.push({ kind: 'box', box: parseSimpleLoop(path, 'RECEIVER', ['PER']) })
        continue
      }
      if (seg.id === 'HL') {
        children.push({ kind: 'box', box: parseBillingHL(path) })
        continue
      }
      children.push({ kind: 'segment', segment: take() }) // BHT
    }
    const closing = take()
    return { path, boxType: 'ST', opening, closing, children }
  }

  function parseGS(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.GS`
    const children: Node[] = []
    while (peek() && peek().id !== 'GE') {
      if (peek().id === 'ST') children.push({ kind: 'box', box: parseST(path) })
      else children.push({ kind: 'segment', segment: take() })
    }
    const closing = take()
    return { path, boxType: 'GS', opening, closing, children }
  }

  function parseISA(): Box {
    const opening = take()
    const path = 'ISA'
    const children: Node[] = []
    while (peek() && peek().id !== 'IEA') {
      if (peek().id === 'GS') children.push({ kind: 'box', box: parseGS(path) })
      else children.push({ kind: 'segment', segment: take() })
    }
    const closing = take()
    return { path, boxType: 'ISA', opening, closing, children }
  }

  const root = parseISA()
  // One interchange per file: segments after its IEA would be dropped from
  // the tree, so say so rather than draw a file that is not the whole text.
  if (i < segments.length) {
    throw new Error(`837P: parsing the interchange: ${segments.length - i} segment(s) follow the IEA (this viewer draws one interchange)`)
  }
  return root
}
