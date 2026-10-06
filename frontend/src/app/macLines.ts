import type { MacJurisdiction, MacNumber, MacPart, MacRollUp, MacState, MacsResponse } from '../api'

// The MAC map's lines (MacTree.tsx), built from /api/macs with no React, so
// scripts/ssr-check-mactree.mjs can check them: what each line says, and
// which lines a search keeps. The tree is the backend's (backend/pkg/macs
// Tree); nothing here places a number, it only words the lines. The
// overview's workload box groups its workloads from the same tree
// (workloadGroups, at the end).

export type MacNodeKind = 'jurisdiction' | 'part' | 'rollup' | 'workload' | 'hhh' | 'other'

export interface MacNode {
  key: string
  kind: MacNodeKind
  id?: string // the five-digit number, on a roll-up's or a number's line
  title: string // 'JM', 'Part A', '11001', '11201'
  detail: string // the contractor, the Part's meaning, a roll-up's states in quotes, a number's states
  states: MacState[] // what a search finds the line by, and its hover text
  codes?: string // the USPS codes on a line of their own, for an HHH number's area
  note?: string // the backend's note, on a line of its own
  aside?: string // 'no roll-up named', for a workload directly under a Part that has roll-ups
  terms: string[] // the words a search matches the line by, lowercase
  children: MacNode[]
}

// What each Part's claims are, on its heading: the claims and the shared
// system that processes them.
export const PART_MEANING: Record<string, string> = {
  A: 'institutional claims, FISS',
  B: 'professional claims, MCS',
  DME: 'durable medical equipment',
}

const partTitle = (part: string) => (part === 'DME' ? 'DME' : part ? `Part ${part}` : 'Part not given')

// Joins words as English does: 'SC', 'VA and WV', 'AL, GA, and TN'.
export function listed(words: string[]): string {
  if (words.length <= 2) return words.join(' and ')
  return `${words.slice(0, -1).join(', ')}, and ${words[words.length - 1]}`
}

const MOST_NAMES = 4 // a number with more states than this lists their codes; the names show on hover

// A number's states as its line shows them: the names for a few, the codes
// for many.
export function statesText(states: MacState[]): string {
  if (states.length === 0) return 'no states named'
  if (states.length > MOST_NAMES) return states.map((s) => s.code).join(', ')
  return states.map((s) => s.name).join(', ')
}

// The words of a text, lowercase, for the search.
function wordsOf(...texts: (string | null | undefined)[]): string[] {
  const out: string[] = []
  for (const text of texts) {
    if (!text) continue
    for (const word of text.toLowerCase().split(/[^a-z0-9*-]+/)) if (word) out.push(word)
  }
  return out
}

function stateTerms(states: MacState[]): string[] {
  const out: string[] = []
  for (const s of states) out.push(s.code.toLowerCase(), ...wordsOf(s.name))
  return out
}

function numberNode(key: string, kind: MacNodeKind, n: MacNumber, aside?: string): MacNode {
  const node: MacNode = {
    key,
    kind,
    id: n.id,
    title: n.id,
    detail: statesText(n.states),
    states: n.states,
    note: n.note ?? undefined,
    aside,
    terms: [...stateTerms(n.states), ...wordsOf(n.note)],
    children: [],
  }
  if (kind === 'hhh') {
    // An HHH number's note is its HHH region: 'home health and hospice, Region C'.
    node.detail = n.note ? `home health and hospice, ${n.note}` : 'home health and hospice'
    node.codes = n.states.map((s) => s.code).join(', ')
    node.note = undefined
    node.terms.push('hhh', 'home', 'health', 'hospice')
  }
  return node
}

function partNode(j: MacJurisdiction, p: MacPart): MacNode {
  const key = `${j.jurisdiction}:${p.part}`
  const children: MacNode[] = []
  for (const r of p.roll_ups) {
    const scope = r.states.map((s) => s.code)
    if (r.carries_hhh) scope.push('the HHH claims')
    children.push({
      key: `${key}:${r.id}`,
      kind: 'rollup',
      id: r.id,
      title: r.id,
      detail: `“${listed(scope)}”`,
      // A roll-up is found by its number and its note; its states belong to
      // the workloads under it, so a search for a state shows those.
      states: r.states,
      note: r.note ?? undefined,
      terms: ['roll-up', ...wordsOf(r.note)],
      children: r.workloads.map((w) => numberNode(`${key}:${r.id}:${w.id}`, 'workload', w)),
    })
  }
  const aside = p.roll_ups.length > 0 ? 'no roll-up named' : undefined
  for (const w of p.workloads) children.push(numberNode(`${key}:${w.id}`, 'workload', w, aside))
  for (const w of p.hhh) children.push(numberNode(`${key}:${w.id}`, 'hhh', w))
  for (const w of p.other) children.push(numberNode(`${key}:${w.id}`, 'other', w))
  return {
    key,
    kind: 'part',
    title: partTitle(p.part),
    detail: PART_MEANING[p.part] ?? '',
    states: [],
    terms: p.part === 'DME' ? ['dme'] : [],
    children,
  }
}

// The tree's lines, in the backend's order.
export function macNodes(data: MacsResponse): MacNode[] {
  return data.jurisdictions.map((j) => ({
    key: j.jurisdiction,
    kind: 'jurisdiction' as const,
    title: j.jurisdiction,
    detail: j.contractor,
    states: [],
    terms: [j.jurisdiction.toLowerCase(), ...wordsOf(j.contractor)],
    children: j.parts.map((p) => partNode(j, p)),
  }))
}

// Whether one search word answers a line: digits match the start of the
// line's number; a word of one or two letters matches one of the line's
// words whole (a USPS code, a jurisdiction, the C of Region C), so "va" is
// Virginia and not Nevada; a longer word matches the start of one (a state's
// name, a contractor, a note).
function answers(node: MacNode, word: string): boolean {
  if (/^\d+$/.test(word)) return node.id !== undefined && node.id.startsWith(word)
  if (word.length <= 2) return node.terms.includes(word)
  return node.terms.some((t) => t.startsWith(word))
}

// The lines a search keeps. Every word must be answered by a line or a line
// above it: "palmetto virginia" keeps Virginia's numbers under Palmetto's
// jurisdictions. A line that answers the last of the words keeps every line
// under it; a line above a kept line stays to hold it.
export function searchTree(nodes: MacNode[], query: string): MacNode[] {
  const words = query.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (words.length === 0) return nodes
  const keep = (node: MacNode, left: string[]): MacNode | null => {
    const rest = left.filter((w) => !answers(node, w))
    if (rest.length === 0) return node
    const children: MacNode[] = []
    for (const child of node.children) {
      const kept = keep(child, rest)
      if (kept) children.push(kept)
    }
    return children.length > 0 ? { ...node, children } : null
  }
  const out: MacNode[] = []
  for (const node of nodes) {
    const kept = keep(node, words)
    if (kept) out.push(kept)
  }
  return out
}

// How many numbers (roll-ups included) the lines hold.
export function numbersIn(nodes: MacNode[]): number {
  let n = 0
  for (const node of nodes) n += (node.id !== undefined ? 1 : 0) + numbersIn(node.children)
  return n
}

// One line as the drawing writes it: the title and the detail, and the
// aside in parentheses. The SSR check compares these with the drawing.
export function lineText(node: MacNode): string {
  const title = node.kind === 'rollup' ? `roll-up ${node.title}` : node.kind === 'part' ? `${node.title}  (${node.detail})` : node.title
  const parts = node.kind === 'part' ? [title] : [title, node.detail]
  if (node.aside) parts.push(`(${node.aside})`)
  return parts.filter(Boolean).join('  ')
}

// ---------------------------------------------------------------- the overview's workload box

// One item of a jurisdiction's group in the workload box: a roll-up, drawn
// as a heading over the workloads under it that the box lists (`own` when
// the roll-up's number is itself one it lists), or a number under no
// roll-up.
export type WorkloadItem =
  | { kind: 'rollup'; id: string; own: boolean; rollUp: MacRollUp; workloads: MacNumber[] }
  | { kind: 'number'; id: string; number: MacNumber }

export interface WorkloadGroup {
  jurisdiction: string
  contractor: string
  items: WorkloadItem[]
}

// The workload box's two tiers (SelectionHeader.tsx): the numbers it lists,
// under their jurisdictions in the tree's order, Part A before Part B as the
// tree has them, a roll-up heading the workloads under it as the MAC map
// draws it. A number the tree does not hold is unplaced, listed last in the
// order given.
export function workloadGroups(data: MacsResponse, shown: string[]): { groups: WorkloadGroup[]; unplaced: string[] } {
  const want = new Set(shown)
  const placed = new Set<string>()
  // Places a listed number the first time the tree has it, once.
  const take = (id: string) => {
    if (!want.has(id) || placed.has(id)) return false
    placed.add(id)
    return true
  }
  const groups: WorkloadGroup[] = []
  for (const j of data.jurisdictions) {
    const items: WorkloadItem[] = []
    for (const p of j.parts) {
      for (const r of p.roll_ups) {
        const own = take(r.id)
        const workloads = r.workloads.filter((w) => take(w.id))
        if (own || workloads.length > 0) items.push({ kind: 'rollup', id: r.id, own, rollUp: r, workloads })
      }
      for (const n of [...p.workloads, ...p.hhh, ...p.other]) if (take(n.id)) items.push({ kind: 'number', id: n.id, number: n })
    }
    if (items.length > 0) groups.push({ jurisdiction: j.jurisdiction, contractor: j.contractor, items })
  }
  return { groups, unplaced: shown.filter((id) => !placed.has(id)) }
}
