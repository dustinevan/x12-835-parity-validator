import { useDeferredValue, useEffect, useMemo, useRef, useState, type Dispatch, type SetStateAction } from 'react'
import type { SegmentDisplay } from './display'
import { allBoxPaths, findBox, type Box, type Segment } from './tree'

// The raw file as page state, for the three viewers (835, 837P, 837I): the
// selected example's own text until a reader edits or pastes over it, and the
// boxes drawn from whatever text is current.
//
// The text lives only here. The example data is never written, the text is
// not put in the address or in browser storage (a pasted 835 or 837 can be
// PHI, and nothing leaves the page), and picking another example drops it.
//
// Text goes through invalid states on every keystroke, so the text is parsed
// inside a try block and then checked (the envelope closes, every row can be
// described). A failure never reaches the page as a thrown error: the boxes
// keep the last tree that did parse and `error` carries the parser's message.
// The library checks in each viewer's check.ts are for the checked-in
// examples and do not run here.

export interface EditableRaw<B extends string> {
  root: Box<B>
  text: string
  // The text differs from the selected example's own.
  edited: boolean
  // The panel shows the textarea rather than the read view.
  editing: boolean
  // Why the current text could not be drawn; null when it can. The boxes then
  // show the last text that could.
  error: string | null
  setText: (next: string) => void
  setEditing: (on: boolean) => void
  reset: () => void
}

// A viewer's own describe functions, the same three it hands to BoxView.
export interface Describe<B extends string> {
  describeOpening: (box: Box<B>) => SegmentDisplay
  describeRow: (segment: Segment, host: B) => SegmentDisplay
  describeClosing: (box: Box<B>) => SegmentDisplay
}

type Parsed<B extends string> = { ok: true; root: Box<B> } | { ok: false; error: string }

// ISA, GS, and ST each close with their own segment in every transaction the
// viewers draw. A parser that reaches the end of the text without finding
// one hands back a box with no closing segment.
const CLOSER: Record<string, string> = { ISA: 'IEA', GS: 'GE', ST: 'SE' }

function message(err: unknown): string {
  return err instanceof Error ? err.message : String(err)
}

function shown(seg: Segment): string {
  return seg.raw.length > 48 ? `${seg.raw.slice(0, 48)}…` : seg.raw
}

// Deepest boxes first, so a missing SE is reported as the ST's rather than
// as the ISA's.
function checkEnvelope<B extends string>(box: Box<B>): void {
  for (const node of box.children) if (node.kind === 'box') checkEnvelope(node.box)
  const id = (box.opening as Segment | undefined)?.id
  const want = id === undefined ? undefined : CLOSER[id]
  if (want !== undefined && (box.closing as Segment | undefined)?.id !== want) {
    throw new Error(`x12: checking the envelope: ${id} has no closing ${want}`)
  }
}

// Describes every box and row once, the way BoxView will, so a tree that the
// parser built but the display cannot draw is rejected here and not in render.
function probeDisplay<B extends string>(box: Box<B>, d: Describe<B>): void {
  const run = (seg: Segment, fn: () => void) => {
    try {
      fn()
    } catch (err) {
      throw new Error(`x12: showing ${shown(seg)}: ${message(err)}`)
    }
  }
  run(box.opening, () => d.describeOpening(box))
  for (const node of box.children) {
    if (node.kind === 'box') probeDisplay(node.box, d)
    else run(node.segment, () => d.describeRow(node.segment, box.boxType))
  }
  if (box.closing) run(box.closing, () => d.describeClosing(box))
}

// Parses and checks one text; never throws. Exported for scripts/ssr-check-viewer-edit.mjs.
export function tryParse<B extends string>(text: string, parse: (raw: string) => Box<B>, describe: Describe<B>): Parsed<B> {
  if (text.trim() === '') return { ok: false, error: 'x12: reading the text: it is empty' }
  try {
    const root = parse(text)
    if ((root.opening as Segment | undefined)?.id !== 'ISA') {
      throw new Error('x12: reading the text: it does not start with an ISA segment')
    }
    checkEnvelope(root)
    probeDisplay(root, describe)
    return { ok: true, root }
  } catch (err) {
    return { ok: false, error: message(err) }
  }
}

export function useEditableRaw<B extends string>(example: { id: string; file: string }, parse: (raw: string) => Box<B>, describe: Describe<B>): EditableRaw<B> {
  const [edit, setEdit] = useState<{ id: string; text: string } | null>(null)
  const [editing, setEditing] = useState(false)

  // Picking another example drops the edit (React's own pattern for state
  // that depends on a prop: set it during render, before anything commits).
  if (edit !== null && edit.id !== example.id) setEdit(null)
  const text = edit !== null && edit.id === example.id ? edit.text : example.file
  const edited = text !== example.file

  // The example's own text parses as it always has (a failure there is a bug
  // in the example library and stays loud). Edited text parses a render
  // behind the keystroke, so a large paste keeps the textarea responsive.
  const exampleRoot = useMemo(() => parse(example.file), [parse, example])
  const lagging = useDeferredValue(text)
  const parseText = edited ? lagging : text
  const result = useMemo<Parsed<B>>(
    () => (parseText === example.file ? { ok: true, root: exampleRoot } : tryParse(parseText, parse, describe)),
    [parseText, example, exampleRoot, parse, describe],
  )

  const lastGood = useRef(exampleRoot)
  useEffect(() => {
    if (result.ok) lastGood.current = result.root
  }, [result])

  return {
    root: result.ok ? result.root : lastGood.current,
    text,
    edited,
    editing,
    error: result.ok ? null : result.error,
    setText: (next) => setEdit({ id: example.id, text: next }),
    setEditing,
    reset: () => setEdit(null),
  }
}

// What the viewers did when the example changed — every box open, nothing
// maximized — now also has to survive an edit. Same example: boxes that
// existed before and were collapsed stay collapsed, new boxes open, and a
// maximized box that no longer exists is let go. Other example: start over.
export function useSyncBoxState<B extends string>(
  root: Box<B>,
  exampleId: string,
  setOpenPaths: Dispatch<SetStateAction<Set<string>>>,
  setMaximized: Dispatch<SetStateAction<string | null>>,
): void {
  const seen = useRef<{ exampleId: string; paths: Set<string> } | null>(null)
  useEffect(() => {
    const paths = allBoxPaths(root)
    const before = seen.current ?? { exampleId, paths }
    seen.current = { exampleId, paths }
    if (before.exampleId !== exampleId) {
      setOpenPaths(paths)
      setMaximized(null)
      return
    }
    setOpenPaths((open) => {
      const next = new Set<string>()
      for (const path of paths) if (open.has(path) || !before.paths.has(path)) next.add(path)
      return next
    })
    setMaximized((path) => (path !== null && findBox(root, path) !== undefined ? path : null))
  }, [root, exampleId, setOpenPaths, setMaximized])
}
