import { Fragment, type ReactNode } from 'react'
import { Link } from 'react-router'
import type { FrontFilter, FrontResponse, FrontSums, FrontValue, LegacySums, LegacyValue, Submission } from '../api'
import { compare, count, money } from '../format'
import { SlotChip } from './VendorPicker'
import type { Term } from './help'
import { Info } from './Info'
import type { Side } from './selection'
import { LEGACY_DIMENSIONS, compareValues } from './overviewFilters'

// The parts every table on the overview shares. Every table sits on one
// fixed grid of four columns of one width: the label, then Legacy (the
// MAC-processed side), then slot 1, then slot 2, each figure column holding
// a count and an amount. A column keeps its place in every table: where a
// table has no legacy figure of its own, or a slot has no submission, its
// column stays empty (Gap) and the others do not move, so the eye finds
// Legacy, 1, and 2 at the same place all the way down the page.

export interface OverviewSide {
  side: Side
  submission: Submission
  data: FrontResponse | undefined
  loading: boolean
  error: Error | undefined
}

// The two slots, in the grid's order.
export const SLOTS: Side[] = ['v1', 'v2']

// The one grid: the label, Legacy, slot 1, slot 2, each 15rem.
export const GRID = { gridTemplateColumns: 'repeat(4, 15rem)' }

// The grid of the status table, whose slot columns hold a matrix each: the
// same columns while a matrix fits in 15rem, wider where one does not.
export const MATRIX_GRID = { gridTemplateColumns: '15rem 15rem repeat(2, minmax(15rem, max-content))' }

// A sum of paid amounts with its sign: a reversal pays a negative amount.
export const paid = (amount: string) => (compare(amount, '0') < 0 ? `−${money(amount)}` : money(amount))

// The legacy column's mark, beside the slots' numbers.
export function LegacyChip() {
  return (
    <span className="inline-flex h-5 w-5 shrink-0 items-center justify-center border border-ink bg-ink font-mono text-[11px] font-semibold text-bg">
      L
    </span>
  )
}

// An empty column: Legacy in a table where the legacy side has no figure of
// its own, or a slot with no submission picked. It holds the column's place
// and shows nothing.
export const Gap = ({ column }: { column: 'legacy' | Side }) => <span aria-hidden data-gap={column} />

// One cell per slot, in the grid's order: `cell` for a picked slot, a Gap
// for a slot with no submission.
export function slotCells(sides: OverviewSide[], cell: (s: OverviewSide) => ReactNode): ReactNode[] {
  return SLOTS.map((slot) => {
    const s = sides.find((x) => x.side === slot)
    return <Fragment key={slot}>{s ? cell(s) : <Gap column={slot} />}</Fragment>
  })
}

// The legacy figures every response carries: the same in each, since both
// slots are compared with one MAC submission and the legacy figures are cut
// only by what the MAC claim carries.
export const legacyOf = (sides: OverviewSide[]) => sides.find((s) => s.data)?.data?.legacy

// A cell's two figures in two fixed widths, right-aligned, so counts line up
// over counts and amounts over amounts in every table: a count on the left
// (a link when `to` is given), an amount or a note on the right.
export function Slots({ left, right, to, dim, strong }: {
  left: ReactNode
  right: ReactNode
  to?: string
  dim?: boolean // a value left out of the selection
  strong?: boolean // the right figure is the line's main figure
}) {
  const leftClass = `w-16 shrink-0 text-right font-mono text-sm font-semibold tabular-nums ${dim ? 'text-ink-2' : 'text-ink'}`
  return (
    <span className="flex items-baseline justify-end gap-3">
      {to ? <Link to={to} className={`${leftClass} hover:text-link`}>{left}</Link> : <span className={leftClass}>{left}</span>}
      <span className={`w-28 shrink-0 text-right font-mono tabular-nums ${strong ? 'text-sm text-ink' : 'text-xs text-ink-2'}`}>{right}</span>
    </span>
  )
}

// The header cells of the three figure columns: the mark and the name, then
// what the two figures are, or one note across both (a chart's columns).
// `legacy` absent leaves the Legacy column empty: the table has no legacy
// figure. A chart's heads sit at the start of their columns, where its bars
// start.
export function ColumnHeads({ sides, legacy, vendor, start = false }: {
  sides: OverviewSide[]
  legacy?: [string, string] | string
  vendor: [string, string] | string
  start?: boolean
}) {
  const figures = (names: [string, string] | string) =>
    typeof names === 'string' ? (
      <span className={`text-[11px] text-ink-2 ${start ? 'text-left' : 'text-right'}`}>{names}</span>
    ) : (
      <span className="flex justify-end gap-3 text-[11px] text-ink-2">
        <span className="w-16 shrink-0 text-right">{names[0]}</span>
        <span className="w-28 shrink-0 text-right">{names[1]}</span>
      </span>
    )
  const head = `flex flex-col gap-1 ${start ? 'items-start' : 'items-end'}`
  return (
    <>
      {legacy === undefined ? <Gap column="legacy" /> : (
        <span data-column="legacy" className={head}>
          <span className="flex items-center gap-2 text-xs font-semibold text-ink"><LegacyChip />Legacy</span>
          {figures(legacy)}
        </span>
      )}
      {slotCells(sides, (s) => (
        <span data-column={s.side} className={head}>
          <span className="flex items-center gap-2 text-xs font-semibold text-ink"><SlotChip side={s.side} />{s.submission.label}</span>
          {figures(vendor)}
        </span>
      ))}
    </>
  )
}

// One row's legacy figures, from the response's legacy block: the MAC
// claims and what the MAC paid on them (CLP04), or the count alone where the
// table has no amounts. A dash where the legacy side has no claim with the
// row's value (the empty MAC status of claims only a vendor sent). The data
// attributes let the page check add the column up.
export function LegacyCell({ sums, amount = true, dim, selected = true }: {
  sums: Pick<LegacySums, 'claims' | 'paid'> | undefined
  amount?: boolean
  dim?: boolean
  selected?: boolean // false for a value left out of the selection
}) {
  return (
    <span className="block" data-legacy-claims={sums?.claims ?? 0} data-in={selected ? '1' : '0'}>
      {sums
        ? <Slots left={count(sums.claims)} right={amount ? paid(sums.paid) : ''} dim={dim} />
        : <Slots left="—" right={amount ? '—' : ''} dim={dim} />}
    </span>
  )
}

// One submission's figures in a row: its results there and what it paid
// (CLP04). The amount is a dash where the comparison keeps no claim of the
// vendor's (overviewFilters.ts hasVendorAmount).
export function SideCell({ sums, amount = true, to, dim }: {
  sums: Pick<FrontSums, 'claims' | 'vendor_paid'> | undefined
  amount?: boolean
  to?: string
  dim?: boolean
}) {
  if (sums === undefined) return <Slots left="—" right="—" dim={dim} />
  return <Slots left={count(sums.claims)} right={amount ? paid(sums.vendor_paid) : '—'} to={to} dim={dim} />
}

// A bordered section with a title, its ⓘ, and a note on the right.
export function Tile({ title, info, right, children }: { title: string; info: Term; right?: ReactNode; children: ReactNode }) {
  return (
    <section className="border border-line">
      <div className="flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1 border-b border-line px-4 py-2.5">
        <h2 className="flex items-center gap-1 text-xs font-semibold text-ink">{title}<Info term={info} /></h2>
        {right && <span className="font-mono text-[11px] text-ink-2">{right}</span>}
      </div>
      <div className="px-4 py-4">{children}</div>
    </section>
  )
}

export interface ValueRow {
  value: string
  sums: Partial<Record<Side, FrontValue>>
  legacy?: LegacyValue // for a dimension the MAC claim carries
}

// The values of one dimension across the picked submissions and, for a
// dimension the MAC claim carries, the legacy figures, merged by value, in
// the dimension's order.
export function valuesOf(sides: OverviewSide[], name: FrontFilter): ValueRow[] {
  const byValue = new Map<string, ValueRow>()
  const row = (value: string) => byValue.get(value) ?? { value, sums: {} }
  for (const s of sides) {
    const dimension = s.data?.dimensions.find((d) => d.name === name)
    for (const v of dimension?.values ?? []) {
      const r = row(v.value)
      r.sums[s.side] = v
      byValue.set(v.value, r)
    }
  }
  if (LEGACY_DIMENSIONS.includes(name)) {
    for (const v of legacyOf(sides)?.dimensions.find((d) => d.name === name)?.values ?? []) {
      const r = row(v.value)
      r.legacy = v
      byValue.set(v.value, r)
    }
  }
  return [...byValue.values()].sort((x, y) => compareValues(name, x.value, y.value))
}
