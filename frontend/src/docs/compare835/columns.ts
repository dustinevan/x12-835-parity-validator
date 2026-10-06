// Which example fills each of the list's three columns: the reference and the
// two compared (A and B), kept in the address as ref=, a=, and b= (example
// ids, the ones the 835 viewer's sidebar lists). A column may instead hold
// text the reader pasted: that is page state, never in the address (the
// helpers at the end of this file read it).
import { VIEWER_835_EXAMPLES } from '../../data/viewer835Examples'
import { parseFile, type ParsedFile } from './file'
import { COLS, listOf, type Col, type ListSection, type Sources } from './rows'

// Three examples that hold one claim: the same patient (MBI), date of
// service, billed amount, payee NPI, and payer workload, so the claim's rows
// sit side by side and what differs is the rest of each file.
export const DEFAULT_COLUMNS: Record<Col, string> = {
  ref: 'paid-by-ach',
  a: 'fee-schedule-reduction',
  b: 'overpayment-recovered',
}

export const exampleIds = VIEWER_835_EXAMPLES.map((e) => e.id)

// Each example parsed once.
const parsedCache = new Map<string, ParsedFile>()
export function parsedOf(id: string): ParsedFile {
  let parsed = parsedCache.get(id)
  if (!parsed) {
    const example = VIEWER_835_EXAMPLES.find((e) => e.id === id) ?? VIEWER_835_EXAMPLES[0]
    parsed = parseFile(example.file)
    parsedCache.set(id, parsed)
  }
  return parsed
}

// The match keys of every claim in an example.
export const keysOf = (id: string): Set<string> => new Set(parsedOf(id).claims.flatMap((c) => (c.key ? [c.key.id] : [])))

// The other examples that hold a claim with a key the example holds, in the
// order the viewer's sidebar lists them.
export function partnersOf(id: string): string[] {
  const keys = keysOf(id)
  return exampleIds.filter((other) => other !== id && [...keysOf(other)].some((k) => keys.has(k)))
}

// The three examples for an address. With none named, the defaults. With
// some named, a column left out takes the first example that shares a claim
// with the reference (so a link from the viewer, which names only the
// reference, opens on rows that line up), then the defaults, then any other
// example, none twice.
export function resolveColumns(given: Record<Col, string | null>): Record<Col, string> {
  const known = (id: string | null) => (id !== null && exampleIds.includes(id) ? id : null)
  const named = { ref: known(given.ref), a: known(given.a), b: known(given.b) }
  if (named.ref === null && named.a === null && named.b === null) return { ...DEFAULT_COLUMNS }

  const ref = named.ref ?? DEFAULT_COLUMNS.ref
  const used = new Set([ref, named.a, named.b])
  const pool = [...partnersOf(ref), DEFAULT_COLUMNS.a, DEFAULT_COLUMNS.b, ...exampleIds]
  const next = () => {
    const id = pool.find((p) => !used.has(p)) ?? ref
    used.add(id)
    return id
  }
  const a = named.a ?? next()
  const b = named.b ?? next()
  return { ref, a, b }
}

// ---------------------------------------------------------------- pasted text

// The file an example holds, one segment per line, as stored: what a column's
// box shows while the column's source is that example.
export const exampleTextOf = (id: string): string => (VIEWER_835_EXAMPLES.find((e) => e.id === id) ?? VIEWER_835_EXAMPLES[0]).file

// One column read: its file, or why there is none. Blank text is no file and
// no error (the column has no rows). Text that does not read is an error and
// no file, so a column never keeps the last good parse. Nothing here throws.
export interface Reading {
  file?: ParsedFile
  error?: string
}

const messageOf = (err: unknown) => (err instanceof Error ? err.message : String(err))

export function readingOfText(text: string): Reading {
  if (text.trim() === '') return {}
  let file: ParsedFile
  try {
    file = parseFile(text)
  } catch (err) {
    return { error: messageOf(err) }
  }
  // The rows are built for this file alone, so a shape the rows cannot hold
  // is reported under its box and never reaches the page.
  try {
    listOf({ ref: file })
  } catch (err) {
    return { error: `x12: building the rows: ${messageOf(err)}` }
  }
  return { file }
}

// The list for the columns that have a file. A column with none (blank text,
// or text that did not read) is left out, and its cells hold no value. If the
// rows cannot be built, the list is empty and the reason comes back for the
// page to show.
export function listOfReadings(readings: Record<Col, Reading>): { sections: ListSection[]; error?: string } {
  const files: Sources<ParsedFile> = {}
  for (const col of COLS) {
    const file = readings[col].file
    if (file) files[col] = file
  }
  try {
    return { sections: listOf(files) }
  } catch (err) {
    return { sections: [], error: `x12: building the rows: ${messageOf(err)}` }
  }
}
