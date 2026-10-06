// Splits a whole X12 file into segments using the delimiters the file itself
// declares in its ISA header, the way backend/pkg/x12837p/parse.go does: the
// element separator is the byte after "ISA", the component separator is
// ISA16 (the byte after the 16th element separator), and the segment
// terminator is the byte after ISA16. A pasted file can therefore use
// something other than `* : ~`, run on one long line, or end its lines with
// CRLF.
//
// Each Segment keeps its text exactly as the file sent it in `raw` (what the
// raw-file panel shows and what hover sync matches on) and its elements
// split on the file's own element separator. The viewers' describe code reads
// a composite element by splitting on ':', so a component separator other
// than ':' is rewritten to ':' inside every element except the ISA's (whose
// ISA16 is the separator itself). A literal ':' inside a data element of such
// a file therefore reads as a component boundary; only composite elements are
// ever split, and those do not carry free text.
import type { Segment } from './tree'

export interface Delimiters {
  element: string
  component: string
  terminator: string
}

export interface SplitFile {
  segments: Segment[]
  delimiters: Delimiters
}

function isSpace(c: string): boolean {
  return /\s/.test(c)
}

// Reads the three separators from the ISA at the start of the text (leading
// whitespace, including a byte-order mark, is skipped). `start` is the index
// of the "I" in "ISA".
export function readDelimiters(text: string): { delimiters: Delimiters; start: number } {
  let start = 0
  while (start < text.length && isSpace(text[start])) start++
  if (!text.startsWith('ISA', start) || text.length < start + 4) {
    throw new Error('x12: reading the delimiters: the text does not start with ISA')
  }
  const element = text[start + 3]

  // ISA16 is the one byte after the 16th element separator, and the segment
  // terminator is the byte after ISA16.
  let isa16 = -1
  let separators = 0
  for (let i = start + 3; i < text.length; i++) {
    if (text[i] !== element) continue
    separators++
    if (separators === 16) {
      isa16 = i + 1
      break
    }
  }
  if (isa16 < 0 || isa16 + 1 >= text.length) {
    throw new Error('x12: reading the ISA: the header is incomplete (it needs 16 element separators, then the component separator and the segment terminator)')
  }
  const component = text[isa16]
  const terminator = text[isa16 + 1]
  if (element === component || element === terminator || component === terminator) {
    throw new Error('x12: reading the ISA: the element separator, component separator, and segment terminator must differ')
  }
  return { delimiters: { element, component, terminator }, start }
}

// What follows a segment's text on its own line in the raw-file panel: the
// file's terminator, or nothing when the terminator is itself a line break.
export function terminatorSuffix(d: Delimiters): string {
  return isSpace(d.terminator) ? '' : d.terminator
}

// Splits text that starts at a segment (an interchange from its ISA, or one
// claim from its CLP) on the delimiters given, the way splitX12 does once it
// has read them from the ISA.
export function splitWith(text: string, d: Delimiters): Segment[] {
  const { element, component, terminator } = d
  const segments: Segment[] = []
  for (const piece of text.split(terminator)) {
    const text = piece.trim() // also drops the CR of a CRLF line end
    if (text === '') continue
    const parts = text.split(element)
    const id = parts[0]
    if (!/^[A-Z0-9]{2,3}$/.test(id)) {
      const shown = id.length > 12 ? `${id.slice(0, 12)}…` : id
      throw new Error(`x12: splitting segment ${segments.length + 1}: the segment ID "${shown}" is not 2 or 3 capital letters or digits`)
    }
    const elements = parts.slice(1)
    segments.push({
      id,
      elements: component === ':' || id === 'ISA' ? elements : elements.map((e) => e.split(component).join(':')),
      raw: text,
    })
  }
  return segments
}

export function splitX12(raw: string): SplitFile {
  const { delimiters, start } = readDelimiters(raw)
  return { segments: splitWith(raw.slice(start), delimiters), delimiters }
}
