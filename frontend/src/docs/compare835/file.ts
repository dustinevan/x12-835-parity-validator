// One 835 file read for the list page: the viewer's loop grammar
// (../viewer835/parse.ts) sorted into the parts the list shows, and each
// claim's match key worked out from the file the way the validator does it.
//
// The delimiters come from the file's own ISA (viewer835/parse.ts
// splitSegments), so an example and a pasted file are read alike; a
// composite reads with ':' whatever separator the file declares. The
// envelope (ISA, GS, ST, SE, GE, IEA) is read but not listed: its control
// numbers, dates, and counts differ between any two files by construction.
// parseFile throws an Error with an "x12: ..." message when the text is not
// one whole interchange: no readable ISA, or no IEA, or segments after it.
import { parseLoopTree, splitSegments, type Box, type Segment } from '../viewer835/parse'
import { splitWith } from '../viewer/x12split'

// An N1 loop (the head is its N1) or an LX group's own rows (the head is its
// LX; the segments are TS3 and TS2).
export interface Loop {
  head: Segment
  segs: Segment[]
}

export interface ParsedLine {
  head: Segment // SVC
  segs: Segment[]
  line_control_number: string // REF*6R; '' when the line has none
}

// The six fields a claim is matched on (backend/lib/core/service/claims.go,
// Key), each in the one spelling the validator hashes.
export interface MatchKey {
  mbi: string
  dateOfService: string
  charge: string
  taxId: string
  npi: string
  workload: string
  id: string // the six fields joined, the same facts giving the same id
}

export interface ParsedClaim {
  head: Segment // CLP
  segs: Segment[] // the claim-level segments after CLP, before the first SVC
  lines: ParsedLine[]
  key?: MatchKey
  keyError?: string // why no key could be worked out
}

export interface ParsedFile {
  header: Segment[] // BPR, TRN, REF, DTM, CUR … directly under ST
  payers: Loop[] // N1*PR loops
  payees: Loop[] // N1*PE loops
  summaries: Loop[] // LX with its TS3 and TS2
  claims: ParsedClaim[] // every CLP in file order
  adjustments: Segment[] // PLB
  envelope: Segment[] // ISA, GS, ST, SE, GE, IEA: read, not listed
}

// ---------------------------------------------------------------- the loop tree

export function parseFile(raw: string): ParsedFile {
  const out: ParsedFile = { header: [], payers: [], payees: [], summaries: [], claims: [], adjustments: [], envelope: [] }

  const claimOf = (box: Box): ParsedClaim => {
    const segs: Segment[] = []
    const lines: ParsedLine[] = []
    for (const node of box.children) {
      if (node.kind === 'segment') {
        segs.push(node.segment)
        continue
      }
      const inner = node.box.children.flatMap((c) => (c.kind === 'segment' ? [c.segment] : []))
      const control = inner.find((s) => s.id === 'REF' && s.elements[0] === '6R')?.elements[1] ?? ''
      lines.push({ head: node.box.opening, segs: inner, line_control_number: control })
    }
    return { head: box.opening, segs, lines }
  }

  const visit = (box: Box) => {
    switch (box.boxType) {
      case 'ISA':
      case 'GS':
      case 'ST':
        out.envelope.push(box.opening)
        for (const node of box.children) {
          if (node.kind === 'box') visit(node.box)
          else if (box.boxType === 'ST') (node.segment.id === 'PLB' ? out.adjustments : out.header).push(node.segment)
          else out.envelope.push(node.segment)
        }
        if (box.closing) out.envelope.push(box.closing)
        break
      case 'N1-PR':
      case 'N1-PE': {
        const loop: Loop = { head: box.opening, segs: box.children.flatMap((c) => (c.kind === 'segment' ? [c.segment] : [])) }
        ;(box.boxType === 'N1-PR' ? out.payers : out.payees).push(loop)
        break
      }
      case 'LX': {
        const loop: Loop = { head: box.opening, segs: [] }
        for (const node of box.children) {
          if (node.kind === 'segment') loop.segs.push(node.segment)
          else out.claims.push(claimOf(node.box))
        }
        out.summaries.push(loop)
        break
      }
    }
  }
  const root = parseLoopTree(raw)
  if (root.closing?.id !== 'IEA') throw new Error('x12: parsing interchange: the text ends before an IEA segment')
  visit(root)

  // The keys come from the flat segment list, which is where the validator
  // reads them; the nth CLP there is the nth claim box in the tree.
  const flat = splitSegments(raw)
  if (segmentsOf(out).length !== flat.length) throw new Error('x12: parsing interchange: segments follow the IEA')
  const keys = claimKeys(flat)
  if (keys.length === out.claims.length) {
    out.claims.forEach((c, i) => {
      if ('id' in keys[i]) c.key = keys[i]
      else c.keyError = keys[i].error
    })
  } else {
    for (const c of out.claims) c.keyError = 'the file holds a CLP outside an LX group'
  }
  return out
}

// ---------------------------------------------------------------- the match key

// A segment that ends the claim before it (claims.go endsClaim): a claim runs
// from its CLP through the segment before the next of these.
const ENDS_CLAIM = new Set(['CLP', 'LX', 'TS3', 'TS2', 'PLB', 'SE', 'GE', 'IEA', 'ST'])

interface Payee {
  workload: string // N1*PR N104
  npi: string // N1*PE N104
  taxId: string // REF*TJ read outside a claim
}

// The key of every claim in the file, in file order, worked out as
// claims.go Split and claimRow do: the payee fields as the header says them
// at the claim's CLP (the last N1*PR, N1*PE, and REF*TJ read outside a claim,
// reset at each ST), and the claim's own fields from its segments.
function claimKeys(segments: Segment[]): (MatchKey | { error: string })[] {
  const keys: (MatchKey | { error: string })[] = []
  let payee: Payee = { workload: '', npi: '', taxId: '' }
  let open: { payee: Payee; segs: Segment[] } | undefined
  const close = () => {
    if (open) keys.push(keyOf(open.segs, open.payee))
    open = undefined
  }
  for (const s of segments) {
    if (open && ENDS_CLAIM.has(s.id)) close()
    switch (s.id) {
      case 'ST':
        payee = { workload: '', npi: '', taxId: '' }
        break
      case 'N1':
        if (s.elements[0] === 'PR') payee = { ...payee, workload: s.elements[3] ?? '' }
        else if (s.elements[0] === 'PE') payee = { ...payee, npi: s.elements[3] ?? '' }
        break
      case 'REF':
        if (!open && s.elements[0] === 'TJ') payee = { ...payee, taxId: s.elements[1] ?? '' }
        break
      case 'CLP':
        open = { payee, segs: [] }
        break
    }
    if (open) open.segs.push(s)
  }
  close()
  return keys
}

function keyOf(segs: Segment[], payee: Payee): MatchKey | { error: string } {
  const mbi = (segs.find((s) => s.id === 'NM1' && s.elements[0] === 'QC')?.elements[8] ?? '').trim().toUpperCase()
  const dateOfService = serviceDate(segs).trim()
  if (!validDate(dateOfService)) return { error: `date of service “${dateOfService}” is not a date` }
  const charge = normalizedCharge((segs[0].elements[2] ?? '').trim())
  if (charge === null) return { error: `billed amount “${segs[0].elements[2] ?? ''}” is not a number` }
  const taxId = payee.taxId.trim().replaceAll('-', '')
  const npi = payee.npi.trim()
  const workload = payee.workload.trim()
  return { mbi, dateOfService, charge, taxId, npi, workload, id: [mbi, dateOfService, charge, taxId, npi, workload].join('\x1f') }
}

// DTM*232 (statement from), else the earliest service-line date (DTM*472, or
// DTM*150 when a line gives a period).
function serviceDate(segs: Segment[]): string {
  let earliest = ''
  for (const s of segs) {
    if (s.id !== 'DTM') continue
    const date = s.elements[1] ?? ''
    switch (s.elements[0]) {
      case '232':
        return date
      case '472':
      case '150':
        if (date !== '' && (earliest === '' || date < earliest)) earliest = date
    }
  }
  return earliest
}

function validDate(text: string): boolean {
  if (!/^\d{8}$/.test(text)) return false
  const y = Number(text.slice(0, 4))
  const m = Number(text.slice(4, 6))
  const d = Number(text.slice(6, 8))
  const date = new Date(Date.UTC(y, m - 1, d))
  return date.getUTCFullYear() === y && date.getUTCMonth() === m - 1 && date.getUTCDate() === d
}

// An exact decimal with trailing fractional zeros removed ("447.00" becomes
// "447", "1469.50" becomes "1469.5", "0.00" becomes "0"); null when the text
// is not a number the validator reads.
function normalizedCharge(text: string): string | null {
  const m = /^([+-]?)(\d*)(?:\.(\d*))?$/.exec(text)
  if (!m || (m[2] === '' && (m[3] ?? '') === '')) return null
  const whole = m[2].replace(/^0+/, '') || '0'
  const fraction = (m[3] ?? '').replace(/0+$/, '')
  const zero = whole === '0' && fraction === ''
  return `${m[1] === '-' && !zero ? '-' : ''}${whole}${fraction ? `.${fraction}` : ''}`
}

// ---------------------------------------------------------------- coverage

// Every segment of the file the parts above hold, the envelope included:
// the count the file's own segments must equal.
export function segmentsOf(file: ParsedFile): Segment[] {
  return [
    ...file.envelope,
    ...file.header,
    ...file.payers.flatMap((l) => [l.head, ...l.segs]),
    ...file.payees.flatMap((l) => [l.head, ...l.segs]),
    ...file.summaries.flatMap((l) => [l.head, ...l.segs]),
    ...file.claims.flatMap((c) => [c.head, ...c.segs, ...c.lines.flatMap((l) => [l.head, ...l.segs])]),
    ...file.adjustments,
  ]
}

// ---------------------------------------------------------------- one claim

// The four separators a claim's file used, as the validator stores them
// beside the claim's bytes (docs/API.md, Claim.delimiters).
export interface ClaimDelimiters {
  element: string
  component: string
  repetition: string
  terminator: string
}

// One claim read from its own segments: the text from its CLP through the
// last segment of its last service line, as the validator keeps a claim,
// with the delimiters its file used. The shape is the one parseFile gives a
// claim, so the compare rows (rows.ts listOfClaims) read it the same way;
// no match key is worked out, since the payee fields sit outside the text.
// Throws an Error with an "x12: ..." message when the text does not start
// with CLP or runs past the claim.
export function parseClaim(text: string, d: ClaimDelimiters): ParsedClaim {
  const segments = splitWith(text.trim(), { element: d.element, component: d.component, terminator: d.terminator })
  if (segments.length === 0 || segments[0].id !== 'CLP') throw new Error('x12: parsing claim: the text does not start with CLP')
  const head = segments[0]
  const segs: Segment[] = []
  const lines: ParsedLine[] = []
  let line: ParsedLine | undefined
  for (const s of segments.slice(1)) {
    if (ENDS_CLAIM.has(s.id)) throw new Error(`x12: parsing claim: a ${s.id} segment follows the claim`)
    if (s.id === 'SVC') {
      line = { head: s, segs: [], line_control_number: '' }
      lines.push(line)
      continue
    }
    if (line === undefined) {
      segs.push(s)
      continue
    }
    line.segs.push(s)
    if (s.id === 'REF' && s.elements[0] === '6R' && line.line_control_number === '') line.line_control_number = s.elements[1] ?? ''
  }
  return { head, segs, lines }
}
