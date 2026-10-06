import { useMemo, useState, type ReactNode } from 'react'
import { useSearchParams } from 'react-router'
import { Empty } from '../components/ViewState'
import { stateOf, type LeafState } from '../data/comparedFields'
import { ListGrid, type RowDecoration } from '../docs/compare835/ListGrid'
import { ListTree } from '../docs/compare835/ListTree'
import { PresetBar, PresetText } from '../docs/compare835/PresetControls'
import type { ParsedClaim } from '../docs/compare835/file'
import { COLS, listOfClaims, rowDiffers, type Col, type ListRow, type ListSection } from '../docs/compare835/rows'
import { rowPath, treeOf } from '../docs/compare835/tree'
import { usePresets } from '../docs/compare835/usePresets'
import type { Slot } from '../linePairing'

// The claim screen's view: the docs site's 835 compare, with one claim per
// column from the database instead of a whole example file. The MAC claim is
// the reference; the claims of slots 1 and 2 are compared against it, one
// row per element in file order, marked where a value differs from the MAC's
// (the same browser rule as the docs page, rows.ts isMarked, which counts two
// dollar amounts within five cents as the same amount, as the comparison
// does). Beside each row that the compared set does not compare, a tag says
// so (data/comparedFields.ts).
// The tree, the Only differences switch, and the presets are the docs page's
// own code, with the same address state (hide=, differs=1) and the same
// browser storage, so a preset made on either page works on both.

// The browser may refuse storage (a private window, cleared site data), so
// every access is guarded and the page works without it. The key is the
// docs compare page's, on purpose.
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

// ---------------------------------------------------------------- the compared set's state of a row

// The compared-fields key of a row: its part, its group, and its element, or
// the qualified tag for a segment whose first element qualifies it (the key
// the Compared fields page uses, so the two pages give one answer).
const TAGGED = new Set(['DTM', 'REF', 'AMT', 'QTY', 'LQ', 'PER'])

function stateKeyOf(row: ListRow): string {
  if (!TAGGED.has(row.group)) return `${row.scope}.${row.group}.${row.leaf}`
  const segment = COLS.map((c) => row.from[c]?.[0]).find((s) => s !== undefined)
  const kind = segment?.elements[0] ?? ''
  return `${row.scope}.${row.group}.${kind ? `${row.group}*${kind}` : row.group}`
}

const STATE_TEXT: Partial<Record<LeafState, string>> = { 'compared, watched': 'watched', 'left out': 'left out', unlisted: 'unlisted', 'outside the claim': 'outside the claim' }

export function decorateRow(_section: ListSection, row: ListRow): RowDecoration | undefined {
  const { state, why } = stateOf(stateKeyOf(row))
  const text = STATE_TEXT[state]
  if (text === undefined) return undefined
  return {
    className: state === 'left out' ? 'opacity-60' : undefined,
    note: <span data-state={state} title={why} className="mt-0.5 inline-block border border-line px-1 font-mono text-[10px] font-normal text-ink-2">{text}</span>,
  }
}

// ---------------------------------------------------------------- the view

// The docs page's columns, as this screen fills and names them: the MAC
// claim in its reference column, slot 1's claim in its first compared column
// and slot 2's in the second, those two named by the slots' numbers.
const LABELS: Record<Col, string> = { ref: 'Reference', a: '1', b: '2' }

export function ClaimCompareView({ claims, head, error }: {
  claims: Partial<Record<'mac' | Slot, ParsedClaim>> // the MAC claim, and each slot's
  head: Record<'mac' | Slot, ReactNode> // each column's header: the submission and what its cell holds
  error?: string // why a column could not be read, shown under the grid
}) {
  const [params, setParams] = useSearchParams()
  const sections = useMemo(() => listOfClaims({ ref: claims.mac, a: claims.v1, b: claims.v2 }, LABELS), [claims])
  const columnHead: Record<Col, ReactNode> = { ref: head.mac, a: head.v1, b: head.v2 }
  const tree = useMemo(() => treeOf(sections), [sections])

  // The hidden set and the switch live in the address, as on the docs page;
  // the hidden set is also kept in this browser so it outlives the page.
  const hideParam = params.get('hide')
  const hideText = hideParam ?? loadHidden()
  const hidden = useMemo(() => new Set(hideText.split(',').filter(Boolean)), [hideText])
  const onlyDifferences = params.get('differs') === '1'
  const change = (apply: (next: URLSearchParams) => void) =>
    setParams((prev) => {
      const next = new URLSearchParams(prev)
      apply(next)
      return next
    }, { replace: true })
  const setHidden = (ids: Set<string>) => {
    const text = [...ids].join(',')
    saveHidden(text)
    change((next) => next.set('hide', text))
  }
  const setOnlyDifferences = (on: boolean) => change((next) => (on ? next.set('differs', '1') : next.delete('differs')))
  const presets = usePresets(tree, hidden, setHidden)
  const [allOpen] = useState(true)

  const visible = (row: ListRow) => !rowPath(row).some((id) => hidden.has(id)) && (!onlyDifferences || rowDiffers(row))
  const shown: ListSection[] = sections.map((s) => ({ ...s, rows: s.rows.filter(visible) })).filter((s) => s.rows.length > 0)
  const total = sections.reduce((n, s) => n + s.rows.length, 0)
  const shownCount = shown.reduce((n, s) => n + s.rows.length, 0)

  // The tree sizes itself with container queries (ListTree.tsx: full width,
  // then a fixed column from 72rem of container), so the section is the
  // container.
  return (
    <section className="@container border border-line">
      <div className="flex flex-col gap-6 p-4 @min-[72rem]:flex-row @min-[72rem]:items-start">
        <ListTree tree={tree} hidden={hidden} onChange={setHidden} top={<PresetBar presets={presets} />} />
        <div className="min-w-0 flex-1">
          <div className="mb-1 flex max-w-[52rem] flex-wrap items-center gap-x-6 gap-y-2">
            <Switch checked={onlyDifferences} onChange={setOnlyDifferences} label="Only differences" />
            <span className="ml-auto font-mono text-[11px] text-ink-2">{shownCount} of {total} rows</span>
          </div>
          {allOpen && <ListGrid sections={shown} head={columnHead} decorate={decorateRow} />}
          {shownCount === 0 && (
            <Empty>{error ?? (total === 0 ? 'No claim could be read for this key.' : 'No row to show. Check more in the tree, or turn off Only differences.')}</Empty>
          )}
          {error && shownCount > 0 && <p className="mt-2 font-mono text-[11px] text-danger" role="alert">{error}</p>}
        </div>
      </div>
      <div className="border-t border-line p-4">
        <PresetText presets={presets} />
      </div>
    </section>
  )
}
