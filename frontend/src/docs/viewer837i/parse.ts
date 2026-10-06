// The 837I's own loop grammar — which segment opens which box, and how deep
// each one nests. Institutional claim shape (005010X223A2): ISA ⊃ GS ⊃ ST ⊃
// (BHT, 1000A Submitter, 1000B Receiver, 2000A Billing Provider HL ⊃
// (2010AA Billing Provider, 2010AB Pay-To Address, 2010AC Pay-To Plan) ⊃
// 2000B Subscriber HL ⊃ (2010BA Subscriber, 2010BB Payer) ⊃
// [2000C Patient HL ⊃ 2010CA Patient] ⊃ 2300 Claim ⊃ (CL1, HI groups,
// 2310A-F claim-level provider loops, 2320 Other Subscriber Information ⊃
// 2330A-I, 2400 Service Line LX ⊃ (2410 Drug, 2420A-D line providers, 2430
// Line Adjudication))); SE, GE, and IEA are the closing edges of their
// envelopes. Segments and elements stay in file order; an element that is
// present but empty stays in the array as '' (it renders as a faint ·) —
// only a trailing, never-sent element is actually absent from the array.
// The generic tree shapes and the walks over an already-built tree (envelope
// detection, path seeding, box lookup) are transaction-agnostic and live in
// frontend/src/docs/viewer/tree.ts; only the grammar below — what each box
// type is, and which segment sequence builds one — is this transaction's
// own, hand-verified against every file under
// frontend/src/data/viewer837i/files/.
//
// Box-path convention (mirrors frontend/src/docs/viewer835/parse.ts's):
// '#n' indexes the nth sibling of a box type that can repeat at its level —
// BILLING-HL, SUBSCRIBER-HL, PATIENT-HL, CLAIM, OTHER-SUBSCRIBER, LINE.
// Every other box type (the 1000A/B, 2010, 2310, 2330, 2410, 2420, 2430
// loops) appears at most once under its parent, so it carries no index —
// e.g. 'ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0.CLAIM#1.LINE#0.DRUG' opens
// the second claim's first line's drug sub-loop.
import * as tree from '../viewer/tree'
import { splitX12 } from '../viewer/x12split'

export type { Segment } from '../viewer/tree'
export { allBoxPaths, ancestorPaths, childBoxCount, envelopePaths, findBox, seedOpenPaths } from '../viewer/tree'

// The file's segments, split on the delimiters its own ISA declares (element
// separator, component separator ISA16, segment terminator), so a file need
// not use `* : ~`, break its lines, or avoid CRLF — frontend/src/docs/viewer/
// x12split.ts. Throws when the text has no readable ISA header.
export function splitSegments(raw: string): tree.Segment[] {
  return splitX12(raw).segments
}

export type BoxType =
  | 'ISA'
  | 'GS'
  | 'ST'
  | 'SUBMITTER'
  | 'RECEIVER'
  | 'BILLING-HL'
  | 'BILLING-PROVIDER'
  | 'PAY-TO-ADDRESS'
  | 'PAY-TO-PLAN'
  | 'SUBSCRIBER-HL'
  | 'SUBSCRIBER'
  | 'PAYER'
  | 'PATIENT-HL'
  | 'PATIENT'
  | 'CLAIM'
  | 'ATTENDING'
  | 'OPERATING'
  | 'OTHER-OPERATING'
  | 'RENDERING'
  | 'SERVICE-FACILITY'
  | 'REFERRING'
  | 'OTHER-SUBSCRIBER'
  | 'OTHER-SUBSCRIBER-NAME'
  | 'OTHER-PAYER'
  | 'OTHER-PAYER-ATTENDING'
  | 'OTHER-PAYER-OPERATING'
  | 'OTHER-PAYER-OTHER-OPERATING'
  | 'OTHER-PAYER-SERVICE-FACILITY'
  | 'OTHER-PAYER-RENDERING'
  | 'OTHER-PAYER-REFERRING'
  | 'OTHER-PAYER-BILLING-PROVIDER'
  | 'LINE'
  | 'DRUG'
  | 'LINE-OPERATING'
  | 'LINE-OTHER-OPERATING'
  | 'LINE-RENDERING'
  | 'LINE-REFERRING'
  | 'LINE-ADJUDICATION'

export type Box = tree.Box<BoxType>
export type Node = tree.Node<BoxType>

// Segment IDs that end every loop (the envelope's own closing edges) —
// never a plain row of any box below ST.
const ENVELOPE_CLOSERS = new Set(['SE', 'GE', 'IEA'])

// NM1 role qualifiers (NM101), by context, to the box type they open. The
// same six roles open a claim-level provider loop (2310A-F), an other-payer
// provider loop (2330C-I, plus 2330I's own Billing Provider role 85), or a
// line-level provider loop (2420A-D, a narrower four of the six) — see
// frontend/src/data/x12Names837i.ts's NM1_ROLE_NAME and
// OTHER_PAYER_PROVIDER_LOOP_NAME/LINE_PROVIDER_LOOP_NAME comments.
const CLAIM_PROVIDER_BOX: Partial<Record<string, BoxType>> = {
  '71': 'ATTENDING',
  '72': 'OPERATING',
  ZZ: 'OTHER-OPERATING',
  '82': 'RENDERING',
  '77': 'SERVICE-FACILITY',
  DN: 'REFERRING',
}
const OTHER_PAYER_PROVIDER_BOX: Partial<Record<string, BoxType>> = {
  '71': 'OTHER-PAYER-ATTENDING',
  '72': 'OTHER-PAYER-OPERATING',
  ZZ: 'OTHER-PAYER-OTHER-OPERATING',
  '77': 'OTHER-PAYER-SERVICE-FACILITY',
  '82': 'OTHER-PAYER-RENDERING',
  DN: 'OTHER-PAYER-REFERRING',
  '85': 'OTHER-PAYER-BILLING-PROVIDER',
}
const LINE_PROVIDER_BOX: Partial<Record<string, BoxType>> = {
  '72': 'LINE-OPERATING',
  ZZ: 'LINE-OTHER-OPERATING',
  '82': 'LINE-RENDERING',
  DN: 'LINE-REFERRING',
}

// Trailing rows a "party" box (any loop that opens with a plain NM1 and
// holds only address/identifier detail below it — 1000A/B, every 2010, 2310,
// 2330, 2420 loop) can carry, in whatever order the file sends them: address
// (N3/N4), demographics (DMG, subscriber/patient only), a date (DTP,
// 2330B's remittance date only), secondary identifiers (REF), specialty
// (PRV, 2310A attending only), and a contact (PER, 1000A/2010AA only). A
// segment never in this set ends the party box — always either a sibling
// loop's own opening segment or the enclosing loop's closing edge.
const PARTY_ROW_IDS = new Set(['N3', 'N4', 'DMG', 'DTP', 'REF', 'PRV', 'PER'])

export function parseLoopTree(raw: string): Box {
  const segments = splitSegments(raw)
  let i = 0
  const peek = () => segments[i]
  const take = () => segments[i++]

  function parseParty(parentPath: string, boxType: BoxType): Box {
    const opening = take() // NM1
    const path = `${parentPath}.${boxType}`
    const children: Node[] = []
    while (peek() && PARTY_ROW_IDS.has(peek().id)) children.push({ kind: 'segment', segment: take() })
    return { path, boxType, opening, children }
  }

  function parseDrug(parentPath: string): Box {
    const opening = take() // LIN
    const path = `${parentPath}.DRUG`
    const children: Node[] = []
    while (peek() && ['CTP', 'REF'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    return { path, boxType: 'DRUG', opening, children }
  }

  function parseLineAdjudication(parentPath: string): Box {
    const opening = take() // SVD
    const path = `${parentPath}.LINE-ADJUDICATION`
    const children: Node[] = []
    while (peek() && ['CAS', 'DTP', 'AMT'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    return { path, boxType: 'LINE-ADJUDICATION', opening, children }
  }

  function parseLine(parentPath: string, index: number): Box {
    const opening = take() // LX
    const path = `${parentPath}.LINE#${index}`
    const children: Node[] = []
    while (peek() && !['LX', 'CLM', 'HL', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      const id = peek().id
      if (id === 'LIN') children.push({ kind: 'box', box: parseDrug(path) })
      else if (id === 'SVD') children.push({ kind: 'box', box: parseLineAdjudication(path) })
      else if (id === 'NM1') {
        const boxType = LINE_PROVIDER_BOX[peek().elements[0]] ?? 'LINE-OPERATING'
        children.push({ kind: 'box', box: parseParty(path, boxType) })
      } else children.push({ kind: 'segment', segment: take() }) // SV2, DTP, REF, AMT, PWK, NTE, HCP
    }
    return { path, boxType: 'LINE', opening, children }
  }

  function parseOtherSubscriber(parentPath: string, index: number): Box {
    const opening = take() // SBR
    const path = `${parentPath}.OTHER-SUBSCRIBER#${index}`
    const children: Node[] = []
    while (peek() && !['SBR', 'LX', 'CLM', 'HL', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      const id = peek().id
      if (id === 'NM1') {
        const qualifier = peek().elements[0]
        const boxType: BoxType = qualifier === 'IL' ? 'OTHER-SUBSCRIBER-NAME' : qualifier === 'PR' ? 'OTHER-PAYER' : (OTHER_PAYER_PROVIDER_BOX[qualifier] ?? 'OTHER-PAYER-ATTENDING')
        children.push({ kind: 'box', box: parseParty(path, boxType) })
      } else children.push({ kind: 'segment', segment: take() }) // CAS, AMT, OI, MIA, MOA
    }
    return { path, boxType: 'OTHER-SUBSCRIBER', opening, children }
  }

  function parseClaim(parentPath: string, index: number): Box {
    const opening = take() // CLM
    const path = `${parentPath}.CLAIM#${index}`
    const children: Node[] = []
    let otherSubIndex = 0
    let lineIndex = 0
    while (peek() && !['CLM', 'HL', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      const id = peek().id
      if (id === 'NM1') {
        const boxType = CLAIM_PROVIDER_BOX[peek().elements[0]] ?? 'ATTENDING'
        children.push({ kind: 'box', box: parseParty(path, boxType) })
      } else if (id === 'SBR') {
        children.push({ kind: 'box', box: parseOtherSubscriber(path, otherSubIndex++) })
      } else if (id === 'LX') {
        children.push({ kind: 'box', box: parseLine(path, lineIndex++) })
      } else {
        children.push({ kind: 'segment', segment: take() }) // DTP, CL1, PWK, CN1, AMT, REF, K3, NTE, CRC, HI, HCP
      }
    }
    return { path, boxType: 'CLAIM', opening, children }
  }

  function parsePatientHL(parentPath: string, index: number): Box {
    const opening = take() // HL (level code 23)
    const path = `${parentPath}.PATIENT-HL#${index}`
    const children: Node[] = []
    if (peek() && peek().id === 'PAT') children.push({ kind: 'segment', segment: take() })
    if (peek() && peek().id === 'NM1') children.push({ kind: 'box', box: parseParty(path, 'PATIENT') })
    let claimIndex = 0
    while (peek() && peek().id === 'CLM') children.push({ kind: 'box', box: parseClaim(path, claimIndex++) })
    return { path, boxType: 'PATIENT-HL', opening, children }
  }

  function parseSubscriberHL(parentPath: string, index: number): Box {
    const opening = take() // HL (level code 22)
    const path = `${parentPath}.SUBSCRIBER-HL#${index}`
    const children: Node[] = []
    if (peek() && peek().id === 'SBR') children.push({ kind: 'segment', segment: take() })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === 'IL') children.push({ kind: 'box', box: parseParty(path, 'SUBSCRIBER') })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === 'PR') children.push({ kind: 'box', box: parseParty(path, 'PAYER') })
    if (peek() && peek().id === 'HL' && peek().elements[2] === '23') {
      children.push({ kind: 'box', box: parsePatientHL(path, 0) })
    } else {
      let claimIndex = 0
      while (peek() && peek().id === 'CLM') children.push({ kind: 'box', box: parseClaim(path, claimIndex++) })
    }
    return { path, boxType: 'SUBSCRIBER-HL', opening, children }
  }

  function parseBillingHL(parentPath: string, index: number): Box {
    const opening = take() // HL (level code 20)
    const path = `${parentPath}.BILLING-HL#${index}`
    const children: Node[] = []
    while (peek() && ['PRV', 'CUR'].includes(peek().id)) children.push({ kind: 'segment', segment: take() })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === '85') children.push({ kind: 'box', box: parseParty(path, 'BILLING-PROVIDER') })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === '87') children.push({ kind: 'box', box: parseParty(path, 'PAY-TO-ADDRESS') })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === 'PE') children.push({ kind: 'box', box: parseParty(path, 'PAY-TO-PLAN') })
    let subIndex = 0
    while (peek() && peek().id === 'HL' && peek().elements[2] === '22') children.push({ kind: 'box', box: parseSubscriberHL(path, subIndex++) })
    return { path, boxType: 'BILLING-HL', opening, children }
  }

  function parseST(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.ST`
    const children: Node[] = []
    if (peek() && peek().id === 'BHT') children.push({ kind: 'segment', segment: take() })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === '41') children.push({ kind: 'box', box: parseParty(path, 'SUBMITTER') })
    if (peek() && peek().id === 'NM1' && peek().elements[0] === '40') children.push({ kind: 'box', box: parseParty(path, 'RECEIVER') })
    let billIndex = 0
    while (peek() && peek().id === 'HL' && peek().elements[2] === '20') children.push({ kind: 'box', box: parseBillingHL(path, billIndex++) })
    const closing = take() // SE
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
    throw new Error(`837I: parsing the interchange: ${segments.length - i} segment(s) follow the IEA (this viewer draws one interchange)`)
  }
  return root
}
