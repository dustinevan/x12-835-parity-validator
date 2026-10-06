// The transaction-agnostic display shapes and cell-building helpers every
// viewer's own `describe.ts` uses to turn each segment into its cells
// (frontend/src/docs/viewer835/describe.ts for the 835,
// frontend/src/docs/viewer837p/describe.ts for the 837P). Nothing here knows
// any segment ID, qualifier code, or public name table — those live in each
// transaction's own names file and describe module.
export interface LabelPart {
  kind: 'label'
  // The raw code, as it appears in the file (a bare segment ID, or a coded
  // value such as CLP02's 1 or a CARC).
  code: string
  // Its public name. Undefined means unresolved (no confirmed public name):
  // the code shows alone, with no hover and no name line.
  name?: string
  // CARC/RARC only: the maintainer's complete text, shown on hover.
  full?: string
}

export interface ValuePart {
  kind: 'value'
  text: string // '' renders as a faint "·"
  amount?: boolean // a number (an amount, a weight, a quantity): the compare page compares it as an exact decimal
}

// A composite element's sub-parts, rendered tight with the file's own
// component separator (':') between them and no inter-token gap, while each
// coded sub-part still shows its own name on hover.
export interface JoinedPart {
  kind: 'joined'
  parts: (LabelPart | ValuePart)[]
}

export type Part = LabelPart | ValuePart | JoinedPart

// One file element, shown as a cell: `ref` is its element reference (always
// present, needs no lookup — "ISA01", "CLP04", "REF01" …), shown as a link
// into the page's lookup with the value beside it; `name` is its public
// element name where the transaction's own names file has one, shown beneath
// (undefined: the cell has no name line); and `part` is its value's content.
export interface Cell {
  ref: string
  name?: string
  part: Part
  // A marker color (a literal Tailwind class) set on one cell so it visibly
  // matches a related mark elsewhere on the page — the 835's CAS group dot,
  // for instance (frontend/src/docs/viewer835/colors.ts).
  marker?: string
}

// A segment or box opening/closing: the row's own leading segment ID and
// name, plus one cell per element the file actually sent (an element
// present-but-empty still gets a cell; a never-sent trailing element gets
// none).
export interface SegmentDisplay {
  lead: LabelPart
  cells: Cell[]
}

export const label = (code: string, name: string | undefined, extra?: Partial<LabelPart>): LabelPart => ({
  kind: 'label',
  code,
  name,
  ...extra,
})
export const value = (text: string | undefined, amount = false): ValuePart => ({ kind: 'value', text: text ?? '', amount })
export const pad = (n: number) => String(n).padStart(2, '0')
export const ref = (segmentId: string, n: number) => `${segmentId}${pad(n)}`
export const cell = (cellRef: string, name: string | undefined, part: Part): Cell => ({ ref: cellRef, name, part })
export const vcell = (cellRef: string, name: string | undefined, v: string | undefined, amount = false): Cell => cell(cellRef, name, value(v, amount))
export const lcell = (cellRef: string, name: string | undefined, code: string, codeName: string | undefined, extra?: Partial<LabelPart>): Cell =>
  cell(cellRef, name, label(code, codeName, extra))

// A compact, single-line form of a segment display — the lead followed by
// each cell's own value content, with no per-cell code or name — used where
// space does not allow the full cell strip (the maximized box's ancestor
// trail).
export function describeCompact(display: SegmentDisplay): Part[] {
  return [display.lead, ...display.cells.map((c) => c.part)]
}
