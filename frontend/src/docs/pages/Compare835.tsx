import { useMemo, useState, type ReactNode } from 'react'
import { Link, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty } from '../../components/ViewState'
import { AMOUNT_MATCH_CENTS } from '../../data/comparedFields'
import { VIEWER_835_EXAMPLES } from '../../data/viewer835Examples'
import { exampleTextOf, listOfReadings, parsedOf, readingOfText, resolveColumns, type Reading } from '../compare835/columns'
import { ListGrid } from '../compare835/ListGrid'
import { ListTree } from '../compare835/ListTree'
import { PresetBar, PresetText } from '../compare835/PresetControls'
import { COLS, COL_LABEL, rowDiffers, type Col, type ListRow, type ListSection } from '../compare835/rows'
import { TextPanel } from '../compare835/TextPanel'
import { rowPath, treeOf } from '../compare835/tree'
import { usePresets } from '../compare835/usePresets'
import { LookupButton, LookupFrame } from '../viewer/LookupFrame'
import { useLookup } from '../viewer/useLookup'
import { LOOKUP_835 } from '../viewer835/lookup'

// Three 835s compared: one row per element in file order, a reference on the
// left and two files compared against it, beside a tree that narrows the
// rows to the elements wanted, and a panel holding each column's raw text.
// The rows and the marks are built in compare835/rows.ts. A column is an
// example (ref=, a=, b= in the address) or the text pasted into its box.
// Pasted text is page state only: never in the address or in browser storage
// (a real 835 is PHI), and a column's parameter is left out of the address
// while its text is the source. The tree's hidden nodes (hide=) and the
// only-differences switch (differs=1) are kept in the address, so a link
// reproduces the view of its examples. Which tree groups are open is local
// to the page. A row's key and the tree's groups and elements show the code
// with the public name beneath it, and every code is a link into the 835
// lookup panel beside the page (lookup=1 and q= in the address;
// viewer/LookupFrame.tsx).
//
// Presets (compare835/presets.ts) are named hidden sets, chosen from a select
// at the top of the tree and read or written as text in a box under the table
// and the Raw text panel. The built-in ones are worked out from the tree's own
// keys; the reader's own are kept in this browser (docs.835-compare.presets).
// A preset's choice goes through setHidden below, so it writes hide= and the
// saved hidden set as a click on the tree does; nothing else about presets is
// in the address, and the preset box's text is page state only, like the Raw
// text.

// ---------------------------------------------------------------- small parts

function Switch({ checked, onChange, label }: { checked: boolean; onChange: (on: boolean) => void; label: string }) {
  return (
    <button type="button" role="switch" aria-checked={checked} onClick={() => onChange(!checked)} className="flex items-center gap-2 text-sm text-ink">
      <span className={`flex h-4 w-8 shrink-0 items-center border px-0.5 ${checked ? 'justify-end border-ink bg-ink' : 'justify-start border-line bg-bg'}`}>
        <span className={`h-2.5 w-2.5 ${checked ? 'bg-bg' : 'bg-ink-2'}`} />
      </span>
      {label}
    </button>
  )
}

// The select's entry for a column that holds pasted text (never an example id).
const PASTED = '#pasted'

// The example each column shows, picked from the viewer's own list, or
// "Pasted text" while the column holds text from its box.
function ExampleSelect({ col, value, pasted, onChange }: { col: Col; value: string; pasted: boolean; onChange: (id: string) => void }) {
  const groups = useMemo(() => {
    const map = new Map<string, typeof VIEWER_835_EXAMPLES>()
    for (const ex of VIEWER_835_EXAMPLES) map.set(ex.group, [...(map.get(ex.group) ?? []), ex])
    return [...map.entries()]
  }, [])
  return (
    <>
      <span className="block font-semibold">{COL_LABEL[col]}</span>
      <select aria-label={`${COL_LABEL[col]} example`} value={pasted ? PASTED : value} onChange={(e) => { if (e.target.value !== PASTED) onChange(e.target.value) }}
        className="mt-1 w-full border border-line bg-bg px-2 py-1 text-xs font-normal text-ink">
        {pasted && <option value={PASTED}>Pasted text</option>}
        {groups.map(([group, examples]) => (
          <optgroup key={group} label={group}>
            {examples.map((ex) => <option key={ex.id} value={ex.id}>{ex.title}</option>)}
          </optgroup>
        ))}
      </select>
    </>
  )
}

// ---------------------------------------------------------------- saved hidden set

// The browser may refuse storage (a private window, cleared site data), so
// every access is guarded and the page works without it.
const HIDDEN_KEY = 'docs.835-compare.hide'

function loadHidden(): string {
  try {
    return localStorage.getItem(HIDDEN_KEY) ?? ''
  } catch {
    return ''
  }
}

function saveHidden(text: string) {
  try {
    localStorage.setItem(HIDDEN_KEY, text)
  } catch {
    // Nothing to do: the address still carries the set for this visit.
  }
}

// ---------------------------------------------------------------- page

// A column's reading: the example's parse (kept once per example), or the
// text in its box. `id` is null while the column holds text.
function useReading(id: string | null, text: string | null): Reading {
  return useMemo(() => (id !== null ? { file: parsedOf(id) } : readingOfText(text ?? '')), [id, text])
}

export function Compare835() {
  const [params, setParams] = useSearchParams()
  const lookup = useLookup()

  const refParam = params.get('ref')
  const aParam = params.get('a')
  const bParam = params.get('b')
  const columns = useMemo(() => resolveColumns({ ref: refParam, a: aParam, b: bParam }), [refParam, aParam, bParam])

  // The text pasted or typed into each box, by column; null while the column
  // shows an example. Page state only: it is never written to the address or
  // to browser storage, and a pick from the column's select replaces it.
  const [pasted, setPasted] = useState<Record<Col, string | null>>({ ref: null, a: null, b: null })
  const reading: Record<Col, Reading> = {
    ref: useReading(pasted.ref === null ? columns.ref : null, pasted.ref),
    a: useReading(pasted.a === null ? columns.a : null, pasted.a),
    b: useReading(pasted.b === null ? columns.b : null, pasted.b),
  }
  const { sections, error: rowsError } = useMemo(() => listOfReadings(reading), [reading.ref, reading.a, reading.b])
  const tree = useMemo(() => treeOf(sections), [sections])

  // What the address holds besides the examples: the unchecked tree nodes and
  // the switch. The hidden set is also saved in this browser, so it
  // outlives the page: an address that names hide= wins
  // (a shared link reproduces its view); otherwise the saved set applies,
  // and every change saves. An empty hide= means "all shown", kept apart from
  // an address with no hide= at all.
  const hideParam = params.get('hide')
  const hideText = hideParam ?? loadHidden()
  const hidden = useMemo(() => new Set(hideText.split(',').filter(Boolean)), [hideText])
  const onlyDifferences = params.get('differs') === '1'
  const change = (apply: (next: URLSearchParams) => void, replace: boolean) =>
    setParams((prev) => {
      const next = new URLSearchParams(prev)
      apply(next)
      return next
    }, { replace })
  const setHidden = (ids: Set<string>) => {
    const text = [...ids].join(',')
    saveHidden(text)
    change((next) => next.set('hide', text), true)
  }
  const setOnlyDifferences = (on: boolean) => change((next) => (on ? next.set('differs', '1') : next.delete('differs')), true)
  const presets = usePresets(tree, hidden, setHidden)

  // A change of source writes all three columns, so the address names what
  // the page shows: each column's example, and nothing for a column that
  // holds pasted text. The other columns keep the examples they show now.
  const exampleNow = (col: Col): string | null => (pasted[col] === null ? columns[col] : null)
  const writeColumns = (next: URLSearchParams, ids: Record<Col, string | null>) => {
    for (const c of COLS) {
      const id = ids[c]
      if (id === null) next.delete(c)
      else next.set(c, id)
    }
  }
  const pick = (col: Col, id: string) => {
    setPasted((prev) => ({ ...prev, [col]: null }))
    change((next) => writeColumns(next, { ref: exampleNow('ref'), a: exampleNow('a'), b: exampleNow('b'), [col]: id }), false)
  }
  // Editing a box (or clearing it) makes its text the column's source. The
  // address changes once, when the column leaves its example.
  const setText = (col: Col, text: string) => {
    const leaving = pasted[col] === null
    setPasted((prev) => ({ ...prev, [col]: text }))
    if (leaving) change((next) => writeColumns(next, { ref: exampleNow('ref'), a: exampleNow('a'), b: exampleNow('b'), [col]: null }), true)
  }

  const visible = (row: ListRow) => !rowPath(row).some((id) => hidden.has(id)) && (!onlyDifferences || rowDiffers(row))
  const shown: ListSection[] = sections.map((s) => ({ ...s, rows: s.rows.filter(visible) })).filter((s) => s.rows.length > 0)
  const total = sections.reduce((n, s) => n + s.rows.length, 0)
  const count = shown.reduce((n, s) => n + s.rows.length, 0)

  const select = (col: Col) => <ExampleSelect col={col} value={columns[col]} pasted={pasted[col] !== null} onChange={(id) => pick(col, id)} />
  const head: Record<Col, ReactNode> = { ref: select('ref'), a: select('a'), b: select('b') }
  const texts: Record<Col, string> = {
    ref: pasted.ref ?? exampleTextOf(columns.ref),
    a: pasted.a ?? exampleTextOf(columns.a),
    b: pasted.b ?? exampleTextOf(columns.b),
  }
  const errors = { ref: reading.ref.error, a: reading.a.error, b: reading.b.error }

  return (
    <LookupFrame lookup={lookup} data={LOOKUP_835}>
      {/* Capped at the content's width — the tree, the table, and their gaps
          and padding — so the outer box hugs what it holds. The breakpoints are
          the page's own width (its frame is a size container), not the
          window's: from 72rem the tree sits beside the table (ListTree.tsx);
          from 93.5rem the text panel sits right of the table and the cap grows
          by its 20rem and a 1.5rem gap (TextPanel.tsx); under that the panel
          sits below the table. Under 57rem the table's columns narrow
          (ListGrid.tsx). */}
      <div className="mx-auto max-w-[72rem] px-6 py-12 @min-[93.5rem]:max-w-[93.5rem]">
        <PageTitle
          lede={`Three 835s as one row per element, in file order, each an example or text pasted into its box. A compared column is marked where its value differs from the reference; two dollar amounts within $${(AMOUNT_MATCH_CENTS / 100).toFixed(2)} of each other are the same amount, as in the comparison. Claims in different files sit side by side when they share a match key: patient MBI, date of service, billed amount, tax ID, payee NPI, and payer workload. Service lines are paired by line control number (REF*6R), then by position.`}
          actions={
            <div className="flex gap-3">
              {/* The viewer opens an example; pasted reference text has none. */}
              {pasted.ref === null && (
                <Link to={`/835-viewer?example=${columns.ref}`} className="border border-line px-3 py-1.5 text-xs text-ink-2 hover:text-ink">
                  Viewer
                </Link>
              )}
              <LookupButton lookup={lookup} />
            </div>
          }
        >
          835 compare
        </PageTitle>

        <section className="border border-line">
          <div className="flex flex-col gap-6 p-4 @min-[93.5rem]:flex-row @min-[93.5rem]:items-start">
            <div className="flex min-w-0 flex-col gap-6 @min-[72rem]:flex-row @min-[72rem]:items-start @min-[93.5rem]:flex-1">
              <ListTree tree={tree} hidden={hidden} onChange={setHidden} top={<PresetBar presets={presets} />} />
              <div className="min-w-0 flex-1">
                {/* As wide as the table (its four fixed columns), so the count sits at its edge. */}
                <div className="mb-1 flex max-w-[52rem] flex-wrap items-center gap-x-6 gap-y-2">
                  <Switch checked={onlyDifferences} onChange={setOnlyDifferences} label="Only differences" />
                  <span className="ml-auto font-mono text-[11px] text-ink-2">{count} of {total} rows</span>
                </div>
                <ListGrid sections={shown} head={head} />
                {count === 0 && (
                  <Empty>
                    {rowsError ?? (total === 0 ? 'No row to show. Pick an example or paste an 835 into a box.' : 'No row to show. Check more in the tree, or turn off Only differences.')}
                  </Empty>
                )}
              </div>
            </div>
            <TextPanel texts={texts} errors={errors} onText={setText} />
          </div>
        </section>

        <PresetText presets={presets} />
      </div>
    </LookupFrame>
  )
}
