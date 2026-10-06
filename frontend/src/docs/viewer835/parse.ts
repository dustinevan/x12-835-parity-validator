// The 835's own loop grammar — which segment opens which box, and how deep
// each one nests: ISA ⊃ GS ⊃ ST ⊃ (header segments, the N1 PR loop, the N1 PE
// loop, LX ⊃ CLP ⊃ SVC); SE, GE, and IEA are the closing edges of their
// envelopes. Segments and elements stay in file order; an element that is
// present but empty stays in the array as '' (it renders as a faint ·) —
// only a trailing, never-sent element is actually absent from the array.
// The generic tree shapes and the walks over an already-built tree
// (envelope detection, path seeding, box lookup) are transaction-agnostic
// and live in frontend/src/docs/viewer/tree.ts; only the grammar below —
// what each box type is, and which segment sequence builds one — is this
// transaction's own.
import * as tree from '../viewer/tree'
import { splitX12 } from '../viewer/x12split'

export type { Segment } from '../viewer/tree'
export { allBoxPaths, ancestorPaths, childBoxCount, envelopePaths, findBox, seedOpenPaths } from '../viewer/tree'

export type BoxType = 'ISA' | 'GS' | 'ST' | 'N1-PR' | 'N1-PE' | 'LX' | 'CLP' | 'SVC'
export type Box = tree.Box<BoxType>
export type Node = tree.Node<BoxType>

// The file's segments, split on the delimiters its own ISA declares (element
// separator, component separator ISA16, segment terminator), so an example
// and a pasted file are read alike: a file need not use `* : ~`, break its
// lines, or avoid CRLF (frontend/src/docs/viewer/x12split.ts, after
// backend/pkg/x12835/parse.go). A component separator other than ':' reads as
// ':' in every element but the ISA's, so composites compare and describe the
// same across files. Throws an Error with an "x12: ..." message when the text
// has no readable ISA or a segment ID that is not 2 or 3 capital letters or
// digits.
export function splitSegments(raw: string): tree.Segment[] {
  return splitX12(raw).segments
}

// Segment IDs that end the enclosing loop one level up, used as lookahead
// stop conditions below (a segment never named here is a plain row of
// whatever loop is currently open).
const ENVELOPE_CLOSERS = new Set(['SE', 'GE', 'IEA'])

export function parseLoopTree(raw: string): Box {
  const segments = splitSegments(raw)
  let i = 0
  const peek = () => segments[i]
  const take = () => segments[i++]

  function parseN1(parentPath: string): Box {
    const opening = take()
    const boxType: BoxType = opening.elements[0] === 'PE' ? 'N1-PE' : 'N1-PR'
    const path = `${parentPath}.${boxType}`
    const children: Node[] = []
    while (peek() && ['N3', 'N4', 'REF', 'PER', 'RDM'].includes(peek().id)) {
      children.push({ kind: 'segment', segment: take() })
    }
    return { path, boxType, opening, children }
  }

  function parseSVC(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.SVC#${index}`
    const children: Node[] = []
    while (peek() && !['SVC', 'CLP', 'LX', 'PLB', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      children.push({ kind: 'segment', segment: take() })
    }
    return { path, boxType: 'SVC', opening, children }
  }

  function parseCLP(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.CLP#${index}`
    const children: Node[] = []
    let svcIndex = 0
    while (peek() && !['CLP', 'LX', 'PLB', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      if (peek().id === 'SVC') children.push({ kind: 'box', box: parseSVC(path, svcIndex++) })
      else children.push({ kind: 'segment', segment: take() })
    }
    return { path, boxType: 'CLP', opening, children }
  }

  function parseLX(parentPath: string, index: number): Box {
    const opening = take()
    const path = `${parentPath}.LX#${index}`
    const children: Node[] = []
    let clpIndex = 0
    while (peek() && !['LX', 'PLB', ...ENVELOPE_CLOSERS].includes(peek().id)) {
      if (peek().id === 'CLP') children.push({ kind: 'box', box: parseCLP(path, clpIndex++) })
      else children.push({ kind: 'segment', segment: take() }) // TS3, TS2
    }
    return { path, boxType: 'LX', opening, children }
  }

  function parseST(parentPath: string): Box {
    const opening = take()
    const path = `${parentPath}.ST`
    const children: Node[] = []
    let lxIndex = 0
    while (peek() && peek().id !== 'SE') {
      const id = peek().id
      if (id === 'N1') children.push({ kind: 'box', box: parseN1(path) })
      else if (id === 'LX') children.push({ kind: 'box', box: parseLX(path, lxIndex++) })
      else children.push({ kind: 'segment', segment: take() }) // BPR, TRN, DTM, CUR, PLB, ...
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

  return parseISA()
}
