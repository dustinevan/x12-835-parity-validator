// The rows of the 835 compare page: three whole example files, a reference
// and two compared (A and B), as one row per element in file order within each
// segment group, one column per file. Nothing here reads a database: a
// compared column's cell is marked where its value differs from the
// reference's value for the same element, the one way for every element
// (isMarked below), with the comparison's five-cent rule for dollar amounts.
//
// The file is shown in the viewer's order: the payment header (BPR, TRN, REF,
// DTM, CUR), the payer and payee N1 loops, the provider summary (LX, TS3,
// TS2), then each claim followed by its service lines, then the provider
// adjustments (PLB). The envelope is not listed (file.ts). Claims are placed
// side by side when they share a match key (file.ts); service lines are
// paired by REF*6R, then by position (../../linePairing.ts).
//
// Rows are built from the segments each file's parse holds (file.ts), named by
// the viewer's own describe.ts, so every element a file sends has a row, with
// the raw value the file carries. Adjustments (CAS) are one row per group and
// reason code, and remark codes (MOA, MIA, LQ) one row per code.
import {
  CARC, CAS_AMOUNT_EL_NAME, CAS_GROUP_NAME, CAS_QTY_EL_NAME, CLAIM_CAS_NAME, CLP_EL_NAME, LQ_EL_NAME, MIA_EL_NAME, MOA_EL_NAME, RARC, SEGMENT_NAME,
  SERVICE_CAS_NAME, SVC_EL_NAME,
} from '../../data/x12Names835'
import { AMOUNT_MATCH_CENTS } from '../../data/comparedFields'
import { compare, ordinal, withinCents } from '../../format'
import { lineTitle, pairedLines, type Slot } from '../../linePairing'
import { label, type Part, type SegmentDisplay } from '../viewer/display'
import { describeOpening, describeRow } from '../viewer835/describe'
import type { BoxType, Segment } from '../viewer835/parse'
import type { Loop, ParsedClaim, ParsedFile, ParsedLine } from './file'

// ---------------------------------------------------------------- columns

// The page's three columns: the reference, then the two compared files, A and
// B. Their service lines pair through the slots of ../../linePairing.ts
// (addClaimGroup below maps the columns onto them).
export type Col = 'ref' | 'a' | 'b'
export const COLS: Col[] = ['ref', 'a', 'b']
export const COL_LABEL: Record<Col, string> = { ref: 'Reference', a: 'A', b: 'B' }
export type Sources<T> = Partial<Record<Col, T>>

function mapCols<T, U>(src: Sources<T>, f: (t: T) => U): Sources<U> {
  const out: Sources<U> = {}
  for (const col of COLS) {
    const t = src[col]
    if (t !== undefined) out[col] = f(t)
  }
  return out
}

// ---------------------------------------------------------------- the rows

// The parts of the list, in the order the list and the tree show them.
export type Scope = 'payment' | 'payer' | 'payee' | 'summary' | 'claim' | 'line' | 'adjustment'
export const SCOPES: Scope[] = ['payment', 'payer', 'payee', 'summary', 'claim', 'line', 'adjustment']
export const SCOPE_TITLE: Record<Scope, string> = {
  payment: 'Payment',
  payer: 'Payer',
  payee: 'Payee',
  summary: 'Provider summary',
  claim: 'Claim',
  line: 'Service line',
  adjustment: 'Provider adjustment',
}

// The loop each part sits in, which names the segments of a loop that vary
// with it (describe.ts).
const HOST: Record<Scope, BoxType> = { payment: 'ST', payer: 'N1-PR', payee: 'N1-PE', summary: 'LX', claim: 'CLP', line: 'SVC', adjustment: 'ST' }

// A row's key, or the note beside it: the code as shown, its public name (the
// code again where it has none), the text a hover shows, and the code the
// lookup is searched for when the code is clicked. A label with no `lookup`
// is not an element's code (a CAS row's key, CO-45, is a group code and a
// reason code, both values) and is not a link.
export interface Label {
  codes: string
  names: string
  hint?: string
  lookup?: string
}

const elementLabel = (ref: string, name?: string): Label => ({ codes: ref, names: name ?? ref, lookup: ref })

export interface ListRow {
  scope: Scope
  group: string // the tree group: the segment, with NM1 split by role ('CLP', 'NM1*QC', 'DTM')
  groupName?: string // that segment's public name, or its qualifier's
  leaf: string // the tree leaf: the element reference ('CLP03'), or the first of a list
  leafCodes: string // how the leaf reads ('MOA03–MOA07')
  name?: string // the element's public name
  key: Label // the row's key
  context?: Label // a repeated segment with its qualifier or place: 'DTM*232', 'NM1*QC', 'PLB (2nd)'
  amount: boolean // compared as an exact decimal
  money: boolean // a dollar amount: within five cents of the reference is the same amount (sameValue)
  values: Sources<string> // absent where the file does not carry the element, or carries it empty
  parts: Sources<Part> // how the viewer describes the value, for its hover
  from: Sources<Segment[]> // the segments the values were read from
}

export interface ListSection {
  id: string
  title: string
  scope: Scope
  held: Record<Col, boolean> // which files carry this part at all
  note?: string
  rows: ListRow[]
}

// ---------------------------------------------------------------- what an element is

const pad2 = (n: number) => String(n).padStart(2, '0')

// MOA and MIA elements that are dollar amounts; the other segments' amounts
// are the cells the viewer marks as amounts (describe.ts).
const AMOUNT_POSITIONS: Record<string, Set<number>> = {
  MOA: new Set([2, 8, 9]),
  MIA: new Set([2, 4, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17, 18, 19, 24]),
}

// The numbers the viewer marks as amounts that are not dollars: the DRG
// weight and the discharge fraction. They compare exactly, as in the
// comparison (backend/lib/core/service/compare.go counts them as numbers,
// not money); every other amount row is a dollar amount.
const NOT_MONEY = new Set(['CLP12', 'CLP13'])

// The positions of a segment that carry a remark code, shown as one list.
const REMARK_POSITIONS: Record<string, number[]> = { MOA: [3, 4, 5, 6, 7], MIA: [5, 20, 21, 22, 23], LQ: [2] }

// The segments whose first element qualifies them (a role, a date kind, a
// reference kind), so two of one kind are different things.
const QUALIFIED = new Set(['NM1', 'DTM', 'REF', 'AMT', 'QTY', 'LQ', 'PER'])

interface Fields {
  scope: Scope
  group: string
  groupName?: string
  context?: Label
}

// What the viewer calls a segment and its elements, across the columns: the
// first lead name, the first qualifier name, and each element's name.
interface Names {
  lead?: string
  qualifier?: string
  elements: Map<string, string>
}

function describeSeg(seg: Segment, scope: Scope, head: boolean): SegmentDisplay {
  const host = HOST[scope]
  return head ? describeOpening({ path: '', boxType: host, opening: seg, children: [] }) : describeRow(seg, host)
}

function namesOf(shown: Sources<SegmentDisplay>, qualified: boolean): Names {
  const names: Names = { elements: new Map() }
  for (const col of COLS) {
    const display = shown[col]
    if (!display) continue
    names.lead ??= display.lead.name
    const first = display.cells[0]?.part
    if (qualified && first?.kind === 'label') names.qualifier ??= first.name
    for (const cell of display.cells) if (cell.name && !names.elements.has(cell.ref)) names.elements.set(cell.ref, cell.name)
  }
  return names
}

// One element of one segment, across the columns.
function elementRow(f: Fields, cls: string, p: number, segs: Sources<Segment>, shown: Sources<SegmentDisplay>, names: Names): ListRow {
  const ref = `${cls}${pad2(p)}`
  const values: Sources<string> = {}
  const parts: Sources<Part> = {}
  for (const col of COLS) {
    const value = segs[col]?.elements[p - 1]
    if (value !== undefined && value !== '') values[col] = value
    const part = shown[col]?.cells.find((c) => c.ref === ref)?.part
    if (part) parts[col] = part
  }
  const name = names.elements.get(ref)
  const amount = AMOUNT_POSITIONS[cls]?.has(p) || COLS.some((c) => {
    const part = parts[c]
    return part?.kind === 'value' && part.amount === true
  })
  return {
    ...f,
    leaf: ref,
    leafCodes: ref,
    name,
    key: elementLabel(ref, name),
    amount,
    money: amount && !NOT_MONEY.has(ref),
    values,
    parts,
    from: mapCols(segs, (s) => [s]),
  }
}

// ---------------------------------------------------------------- lists and adjustments

// The tree leaf of each kind of remark list, and the element's public name.
const REMARK_LEAF: Record<string, { leaf: string; leafCodes: string; name: string }> = {
  MOA: { leaf: 'MOA03', leafCodes: 'MOA03–MOA07', name: MOA_EL_NAME[2] },
  MIA: { leaf: 'MIA05', leafCodes: 'MIA05, MIA20–MIA23', name: MIA_EL_NAME[4] },
  LQ: { leaf: 'LQ02', leafCodes: 'LQ02', name: LQ_EL_NAME[1] },
}

// One row per remark code: MA01 in the reference and not in A is a row where
// only the reference has a code. `lists` are the segments each column's codes
// were read from.
function remarkRows(f: Fields, cls: string, codes: Sources<string[]>, lists: Sources<Segment[]>): ListRow[] {
  const leaf = REMARK_LEAF[cls]
  const distinct: string[] = []
  for (const col of COLS) for (const code of codes[col] ?? []) if (!distinct.includes(code)) distinct.push(code)
  return distinct.map((code) => {
    const values: Sources<string> = {}
    const parts: Sources<Part> = {}
    for (const col of COLS) {
      if (!codes[col]?.includes(code)) continue
      values[col] = code
      parts[col] = label(code, RARC[code]?.short, { full: RARC[code]?.full })
    }
    return {
      ...f,
      ...leaf,
      key: { codes: leaf.leafCodes, names: leaf.name, lookup: leaf.leaf },
      amount: false,
      money: false,
      values,
      parts,
      from: lists,
    }
  })
}

function remarkCodes(list: Segment[] | undefined, cls: string, q: string): string[] | undefined {
  if (!list) return undefined
  const out: string[] = []
  for (const s of list) {
    if (s.id !== cls || (cls === 'LQ' && (s.elements[0] ?? '') !== q)) continue
    for (const p of REMARK_POSITIONS[cls]) if (s.elements[p - 1]) out.push(s.elements[p - 1])
  }
  return out
}

interface Adjusted {
  id: string // 'CO-45#1': the group and reason, and which of that pair this is in the scope
  group: string
  reason: string
  n: number
}

// Every CAS in the scope as one row per group and reason code (the second
// CO-45 in a scope pairs with the second), then a row per adjustment
// quantity.
function casRows(scope: Scope, segs: Sources<Segment[]>): ListRow[] {
  const amounts = new Map<Col, Map<string, string>>()
  const quantities = new Map<Col, Map<string, string>>()
  const order: Adjusted[] = []
  for (const col of COLS) {
    const list = segs[col]
    if (!list) continue
    const amount = new Map<string, string>()
    const quantity = new Map<string, string>()
    const seen = new Map<string, number>()
    for (const s of list) {
      if (s.id !== 'CAS') continue
      const e = s.elements
      // CAS02 to CAS19: up to six reason, amount, quantity triples.
      for (let i = 1; i < e.length && i <= 16; i += 3) {
        const reason = e[i] ?? ''
        if (reason === '') continue
        const group = e[0] ?? ''
        const n = (seen.get(`${group}-${reason}`) ?? 0) + 1
        seen.set(`${group}-${reason}`, n)
        const id = `${group}-${reason}#${n}`
        if (e[i + 1]) amount.set(id, e[i + 1])
        if (e[i + 2]) quantity.set(id, e[i + 2])
        if (!order.some((o) => o.id === id)) order.push({ id, group, reason, n })
      }
    }
    amounts.set(col, amount)
    quantities.set(col, quantity)
  }

  const f: Fields = { scope, group: 'CAS', groupName: scope === 'claim' ? CLAIM_CAS_NAME : SERVICE_CAS_NAME }
  const from: Sources<Segment[]> = {}
  for (const col of COLS) if (segs[col]) from[col] = segs[col]!.filter((s) => s.id === 'CAS')

  const rows: ListRow[] = []
  for (const [found, leaf, leafCodes, name] of [
    [amounts, 'CAS03', 'CAS03, 06, …', CAS_AMOUNT_EL_NAME],
    [quantities, 'CAS04', 'CAS04, 07, …', CAS_QTY_EL_NAME],
  ] as const) {
    for (const { id, group, reason, n } of order) {
      const values: Sources<string> = {}
      for (const col of COLS) {
        const v = found.get(col)?.get(id)
        if (v !== undefined) values[col] = v
      }
      if (Object.keys(values).length === 0) continue
      const code = `${group}-${reason}`
      const meaning = `${CAS_GROUP_NAME[group] ?? group} · ${CARC[reason]?.short ?? reason}`
      rows.push({
        ...f,
        leaf,
        leafCodes,
        name,
        key: { codes: `${code}${ordinal(n)}`, names: `${meaning}${ordinal(n)}` },
        context: { codes: leafCodes, names: name, lookup: leaf },
        amount: found === amounts,
        money: found === amounts,
        values,
        parts: {},
        from,
      })
    }
  }
  return rows
}

// ---------------------------------------------------------------- one part

// One segment of the part as it sits in each column. A repeated kind is told
// apart by its qualifier and its place among the segments of that kind.
interface Occurrence {
  id: string // 'DTM*232#1'
  cls: string // the segment ID
  q: string // the qualifier or role, '' for a segment without one
  n: number // 1 for the first of this kind in the part
  segs: Sources<Segment>
}

// The ids of each column's list merged into one: the first column's order,
// and an id only a later column has placed after the id before it there; one
// with no id before it there goes before the first later id the merged list
// has, else last.
function mergeOrder(lists: string[][]): string[] {
  const merged: string[] = []
  for (const list of lists) {
    let last = -1
    list.forEach((id, i) => {
      const at = merged.indexOf(id)
      if (at >= 0) {
        last = at
        return
      }
      if (last < 0) {
        const next = list.slice(i + 1).map((later) => merged.indexOf(later)).find((index) => index >= 0)
        last = (next ?? merged.length) - 1
      }
      merged.splice(last + 1, 0, id)
      last += 1
    })
  }
  return merged
}

function occurrencesOf(segs: Sources<Segment[]>): Occurrence[] {
  const byId = new Map<string, Occurrence>()
  const orders: string[][] = []
  for (const col of COLS) {
    const list = segs[col]
    if (!list) continue
    const seen = new Map<string, number>()
    const order: string[] = []
    for (const s of list) {
      if (s.id === 'CAS') continue
      const q = QUALIFIED.has(s.id) ? (s.elements[0] ?? '') : ''
      const key = q ? `${s.id}*${q}` : s.id
      const n = (seen.get(key) ?? 0) + 1
      seen.set(key, n)
      // Every LQ of one qualifier is one list of remark codes: one occurrence.
      if (s.id === 'LQ' && n > 1) continue
      const id = `${key}#${n}`
      let occ = byId.get(id)
      if (!occ) {
        occ = { id, cls: s.id, q, n, segs: {} }
        byId.set(id, occ)
      }
      occ.segs[col] = s
      order.push(id)
    }
    orders.push(order)
  }
  return mergeOrder(orders).map((id) => byId.get(id)!)
}

function occurrenceRows(scope: Scope, occ: Occurrence, all: Sources<Segment[]>, placed: Set<string>): ListRow[] {
  const shown = mapCols(occ.segs, (s) => describeSeg(s, scope, false))
  const names = namesOf(shown, QUALIFIED.has(occ.cls))
  const tag = occ.q ? `${occ.cls}*${occ.q}` : occ.cls
  const tagName = occ.q ? (names.qualifier ?? (names.lead ? `${names.lead} ${occ.q}` : undefined)) : undefined
  const f: Fields = {
    scope,
    group: occ.cls === 'NM1' ? `NM1*${occ.q}` : occ.cls,
    groupName: names.lead ?? names.qualifier,
    // The 2nd DTM*472 or the 2nd PLB says so; the first of a plain segment needs no note.
    context: occ.q || occ.n > 1
      ? { codes: `${tag}${ordinal(occ.n)}`, names: `${tagName ?? tag}${ordinal(occ.n)}`, hint: tagName, lookup: tag }
      : undefined,
  }
  const positions = Math.max(...COLS.map((c) => occ.segs[c]?.elements.length ?? 0))
  const remark = REMARK_POSITIONS[occ.cls]
  const rows: ListRow[] = []
  for (let p = 1; p <= positions; p++) {
    if (remark?.includes(p)) {
      // The remark codes of every MOA (or MIA) in the part are one list,
      // shown where its first position falls.
      const list = occ.cls === 'LQ' ? occ.id : occ.cls
      if (p === remark[0] && !placed.has(list)) {
        placed.add(list)
        const codes: Sources<string[]> = {}
        const lists: Sources<Segment[]> = {}
        for (const col of COLS) {
          const found = remarkCodes(all[col], occ.cls, occ.q)
          if (!found) continue
          codes[col] = found
          lists[col] = all[col]!.filter((s) => s.id === occ.cls && (occ.cls !== 'LQ' || (s.elements[0] ?? '') === occ.q))
        }
        rows.push(...remarkRows({ ...f, context: undefined }, occ.cls, codes, lists))
      }
      continue
    }
    rows.push(elementRow(f, occ.cls, p, occ.segs, shown, names))
  }
  return rows
}

// The groups of a claim, then of a service line, in the order the list and
// the tree show them; a segment kind not named is last. The other parts keep
// the order the files send.
const ORDER: Partial<Record<Scope, string[]>> = {
  claim: ['NM1', 'MOA', 'MIA', 'DTM', 'REF', 'PER', 'AMT', 'QTY', 'CAS'],
  line: ['DTM', 'CAS', 'REF', 'AMT', 'LQ'],
}

// The head segment's own positions: CLP, SVC, N1, or LX. Every position any
// column sends gets a row.
function headRows(scope: Scope, head: Sources<Segment>): ListRow[] {
  const shown = mapCols(head, (s) => describeSeg(s, scope, true))
  const cls = (head.ref ?? head.a ?? head.b)!.id
  const names = namesOf(shown, true)
  // The viewer draws a cell for an element a file sends, and not for every
  // one (SVC06 only when it holds a value); the table names them all.
  const table = cls === 'CLP' ? CLP_EL_NAME : cls === 'SVC' ? SVC_EL_NAME : []
  table.forEach((name, i) => {
    if (!names.elements.has(`${cls}${pad2(i + 1)}`)) names.elements.set(`${cls}${pad2(i + 1)}`, name)
  })
  const f: Fields = { scope, group: cls, groupName: SEGMENT_NAME[cls] ?? names.qualifier }
  const positions = Math.max(...COLS.map((c) => head[c]?.elements.length ?? 0))
  const rows: ListRow[] = []
  for (let p = 1; p <= positions; p++) rows.push(elementRow(f, cls, p, head, shown, names))
  return rows
}

// The rows of one part (a claim, a service line, a loop): its head segment's
// positions, then the groups of ORDER, a group's segments in file order.
function scopeRows(scope: Scope, head: Sources<Segment> | undefined, segs: Sources<Segment[]>): ListRow[] {
  const placed = new Set<string>()
  const blocks = occurrencesOf(segs).map((occ) => ({ cls: occ.cls, rows: occurrenceRows(scope, occ, segs, placed) }))
  blocks.push({ cls: 'CAS', rows: casRows(scope, segs) })
  const order = ORDER[scope] ?? []
  const rank = (cls: string) => {
    const at = order.indexOf(cls)
    return at < 0 ? order.length : at
  }
  blocks.sort((x, y) => rank(x.cls) - rank(y.cls))
  return [...(head ? headRows(scope, head) : []), ...blocks.flatMap((b) => b.rows)]
}

// ---------------------------------------------------------------- the sections

const noneHeld = (): Record<Col, boolean> => ({ ref: false, a: false, b: false })

const heldOf = (files: Sources<ParsedFile>, has: (f: ParsedFile) => boolean): Record<Col, boolean> => ({
  ref: files.ref ? has(files.ref) : false,
  a: files.a ? has(files.a) : false,
  b: files.b ? has(files.b) : false,
})

// "Only in Reference and A": where, of the columns that could hold a part,
// some do not. `labels` names the columns: this page's own, or the
// validator app's claim screen's (listOfClaims).
function noteOf(held: Record<Col, boolean>, within: Col[] = COLS, labels: Record<Col, string> = COL_LABEL): string | undefined {
  if (within.every((c) => held[c])) return undefined
  const where = within.filter((c) => held[c]).map((c) => labels[c])
  return where.length ? `Only in ${where.join(' and ')}` : undefined
}

// "71046" or "71046-26": the SVC01 composite without its code-system qualifier.
function procedureOf(line: ParsedLine) {
  const parts = (line.head.elements[0] ?? '').split(':')
  const [code = '', ...modifiers] = parts.length > 1 ? parts.slice(1) : parts
  return `${code}${modifiers.filter(Boolean).length ? `-${modifiers.filter(Boolean).join('-')}` : ''}`
}

interface ClaimGroup {
  id: string
  claims: Sources<ParsedClaim>
  unmatched?: string // why the claim is not placed beside another file's
}

// Claims with one key, one in each file that has it, are one group. A key
// that occurs more than once in a file is not matched in any file (the
// validator sets such claims aside), and neither is a claim with no key: each
// stands alone in its own column.
function claimGroups(files: Sources<ParsedFile>): ClaimGroup[] {
  const repeated = new Set<string>()
  for (const col of COLS) {
    const seen = new Set<string>()
    for (const c of files[col]?.claims ?? []) {
      if (!c.key) continue
      if (seen.has(c.key.id)) repeated.add(c.key.id)
      seen.add(c.key.id)
    }
  }
  const byId = new Map<string, ClaimGroup>()
  const orders: string[][] = []
  for (const col of COLS) {
    const order: string[] = []
    files[col]?.claims.forEach((c, i) => {
      const matched = c.key !== undefined && !repeated.has(c.key.id)
      const id = matched ? `key:${c.key!.id}` : `${col}:${i}`
      let group = byId.get(id)
      if (!group) {
        group = { id, claims: {}, unmatched: matched ? undefined : c.key ? 'The match key repeats in a file, so no claim with it is matched' : `No match key: ${c.keyError}` }
        byId.set(id, group)
      }
      group.claims[col] = c
      order.push(id)
    })
    orders.push(order)
  }
  return mergeOrder(orders).map((id) => byId.get(id)!)
}

// The list for three files: the payment header, the payer and payee, the
// provider summary, each claim with its service lines, then the provider
// adjustments. A part no file carries has no section.
export function listOf(files: Sources<ParsedFile>): ListSection[] {
  const sections: ListSection[] = []
  const add = (section: ListSection) => {
    if (section.rows.length > 0) sections.push(section)
  }
  const count = (get: (f: ParsedFile) => unknown[]) => Math.max(0, ...COLS.map((c) => (files[c] ? get(files[c]!).length : 0)))
  const numbered = (title: string, i: number) => (i === 0 ? title : `${title} ${i + 1}`)

  const header = (f: ParsedFile) => f.header
  const paymentHeld = heldOf(files, (f) => header(f).length > 0)
  add({
    id: 'payment',
    title: SCOPE_TITLE.payment,
    scope: 'payment',
    held: paymentHeld,
    note: noteOf(paymentHeld),
    rows: scopeRows('payment', undefined, mapCols(files, header)),
  })
  for (const [scope, get] of [
    ['payer', (f: ParsedFile) => f.payers],
    ['payee', (f: ParsedFile) => f.payees],
    ['summary', (f: ParsedFile) => f.summaries],
  ] as const) {
    for (let i = 0; i < count(get); i++) {
      const loops: Sources<Loop> = {}
      for (const col of COLS) {
        const loop = files[col] ? get(files[col]!)[i] : undefined
        if (loop) loops[col] = loop
      }
      const held = heldOf(files, (f) => get(f)[i] !== undefined)
      add({
        id: `${scope}:${i}`,
        title: numbered(SCOPE_TITLE[scope], i),
        scope,
        held,
        note: noteOf(held),
        rows: scopeRows(scope, mapCols(loops, (l) => l.head), mapCols(loops, (l) => l.segs)),
      })
    }
  }

  for (const group of claimGroups(files)) addClaimGroup(add, group)

  const plb = (f: ParsedFile) => f.adjustments
  const adjustmentHeld = heldOf(files, (f) => plb(f).length > 0)
  add({
    id: 'adjustment',
    title: SCOPE_TITLE.adjustment,
    scope: 'adjustment',
    held: adjustmentHeld,
    note: noteOf(adjustmentHeld),
    rows: scopeRows('adjustment', undefined, mapCols(files, plb)),
  })
  return sections
}

// The sections of one group of claims, one in each column that holds it: the
// claim, then each paired service line (the first column that holds the
// claim is the one the others' lines pair against). `labels` names the
// columns in the notes and line titles.
function addClaimGroup(add: (section: ListSection) => void, group: ClaimGroup, labels: Record<Col, string> = COL_LABEL) {
  {
    const present = COLS.filter((c) => group.claims[c])
    if (present.length === 0) return
    const held = noneHeld()
    for (const c of present) held[c] = true
    const first = group.claims[present[0]]!
    const title = group.unmatched || !first.key
      ? `Claim · ${first.head.elements[0] ?? ''}`
      : ['Claim', first.key.mbi, first.key.dateOfService, first.head.elements[2]].filter(Boolean).join(' · ')
    add({
      id: `claim:${group.id}`,
      title,
      scope: 'claim',
      held,
      note: group.unmatched ?? noteOf(held, COLS, labels),
      rows: scopeRows('claim', mapCols(group.claims, (c) => c.head), mapCols(group.claims, (c) => c.segs)),
    })

    // The first column that holds the claim is the one the others' lines
    // pair against: it takes linePairing's first place, and the rest take
    // its two slots in order.
    const base = present[0]
    const rest = present.slice(1)
    const column: Record<'mac' | Slot, Col | undefined> = { mac: base, v1: rest[0], v2: rest[1] }
    const paired = pairedLines(first.lines, {
      v1: rest[0] ? group.claims[rest[0]]!.lines : undefined,
      v2: rest[1] ? group.claims[rest[1]]!.lines : undefined,
    })
    paired.forEach((line, i) => {
      const src: Sources<ParsedLine> = {}
      for (const key of ['mac', 'v1', 'v2'] as const) {
        const col = column[key]
        const found = line.src[key]
        if (col && found) src[col] = found
      }
      const lineHeld = noneHeld()
      for (const c of COLS) lineHeld[c] = src[c] !== undefined
      add({
        id: `line:${group.id}:${i}`,
        title: lineTitle(line, true, procedureOf, `Line not in ${labels[base]}`),
        scope: 'line',
        held: lineHeld,
        note: noteOf(lineHeld, present, labels),
        rows: scopeRows('line', mapCols(src, (l) => l.head), mapCols(src, (l) => l.segs)),
      })
    })
  }
}

// The list for claims given directly, one per column, with no file around
// them: the claim section and its paired service lines, as listOf builds them
// for one group of claims that share a key. The validator app's claim screen
// reads claims from the database this way (file.ts parseClaim), with its own
// column names in `labels`.
export function listOfClaims(claims: Sources<ParsedClaim>, labels: Record<Col, string> = COL_LABEL): ListSection[] {
  const sections: ListSection[] = []
  const add = (section: ListSection) => {
    if (section.rows.length > 0) sections.push(section)
  }
  addClaimGroup(add, { id: 'claim', claims }, labels)
  return sections
}

// ---------------------------------------------------------------- marks and differences

const DECIMAL = /^[+-]?\d+(\.\d+)?$/

// Two values of one row: both absent, or the same text; a number as an exact
// decimal ('10.0' equals '10.00'); a dollar amount as the comparison counts
// it, the same when the two are within five cents of each other, edges
// included (AMOUNT_MATCH_CENTS: 80.00 and 79.95 are the same, 80.00 and 79.94
// are not); every other value as exact text, so an identifier that is all
// digits stays a string. An amount on one side only still differs.
export function sameValue(row: Pick<ListRow, 'amount' | 'money'>, x: string | undefined, y: string | undefined): boolean {
  if (x === undefined || y === undefined) return x === undefined && y === undefined
  if (row.amount && DECIMAL.test(x) && DECIMAL.test(y)) return row.money ? withinCents(x, y, AMOUNT_MATCH_CENTS) : compare(x, y) === 0
  return x === y
}

// A compared column's cell is marked where its value differs from the
// reference's, when both files carry the part the row is in; where one does
// not, the empty column says so and nothing in it is marked.
export function isMarked(section: ListSection, row: ListRow, col: Col): boolean {
  return col !== 'ref' && section.held.ref && section.held[col] && !sameValue(row, row.values[col], row.values.ref)
}

// A row differs when its value is not the same in every column.
export function rowDiffers(row: ListRow): boolean {
  return COLS.slice(1).some((c) => !sameValue(row, row.values[c], row.values.ref))
}
