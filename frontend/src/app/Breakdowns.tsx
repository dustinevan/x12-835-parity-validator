import { useState } from 'react'
import type { FrontFilter } from '../api'
import { statusLabel } from '../format'
import { CodeHint } from './CodeHint'
import { HELP } from './help'
import { Info } from './Info'
import { ColumnHeads, GRID, Gap, LegacyCell, SideCell, slotCells, valuesOf, type OverviewSide } from './overviewColumns'
import { DIMENSIONS, LEGACY_DIMENSIONS, hasVendorAmount, isSelected, valueLabel, type OverviewSelection } from './overviewFilters'

// The overview's breakdowns: one small table per dimension the endpoint
// returns, every one of the twelve, rows the dimension's values and the
// page's columns, Legacy, 1, and 2. Each table is the response's facet for
// its dimension, added under the selection on every other dimension and not
// its own, so a value left out of the selection is still listed, struck
// through, with what it would add back; the rows not struck through add up
// to the totals. A table by something the MAC claim carries has a Legacy
// column, from the response's legacy block; a table by something a vendor
// decides (its claim status, the payment band, the outcome) has none, and
// its Legacy column stays empty. A long table shows its first rows and
// opens on a click; the rows it hides stay in the page, hidden, so the page
// check can add them up.

const help = HELP.overview

const FIRST = 6 // rows a long table shows before it is opened
const LONG = 8 // a table with more rows than this starts closed

// "By jurisdiction", "By MAC claim status": the label with its first letter
// lowered unless the label starts with an abbreviation.
const titleOf = (label: string) => `By ${/^[A-Z][a-z]/.test(label) ? label[0].toLowerCase() + label.slice(1) : label}`

// A row's label: a claim status shows its code with the code's hover, then
// the status's name.
function RowLabel({ name, value }: { name: FrontFilter; value: string }) {
  if ((name === 'mac_status' || name === 'vendor_status') && value !== '') {
    return (
      <>
        <span className="font-mono"><CodeHint kind="status" code={value} label={statusLabel(value)}>{value}</CodeHint></span>
        <span className="ml-2">{statusLabel(value)}</span>
      </>
    )
  }
  return <>{valueLabel(name, value)}</>
}

function Breakdown({ name, label, sides, selection }: { name: FrontFilter; label: string; sides: OverviewSide[]; selection: OverviewSelection }) {
  const [open, setOpen] = useState(false)
  const rows = valuesOf(sides, name)
  const legacy = LEGACY_DIMENSIONS.includes(name)
  const long = rows.length > LONG
  const leftOut = rows.filter((r) => !isSelected(selection, name, r.value)).length
  return (
    <section className="border border-line" data-breakdown={name} data-legacy={legacy ? '1' : '0'}>
      <div className="flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1 border-b border-line px-4 py-2">
        <h3 className="text-xs font-semibold text-ink">{titleOf(label)}</h3>
        <span className="font-mono text-[11px] text-ink-2">
          {rows.length} {rows.length === 1 ? 'value' : 'values'}{leftOut > 0 && ` · ${leftOut} left out`}
        </span>
      </div>
      <div className="overflow-x-auto px-4 py-3">
        <div className="grid min-w-max items-center gap-x-6 gap-y-1.5" style={GRID}>
          <span className="text-[11px] text-ink-2">{label}</span>
          <ColumnHeads sides={sides} legacy={legacy ? ['MAC claims', 'MAC paid'] : undefined} vendor={['claims', 'paid']} />
          {rows.length === 0 && <span className="col-span-full font-mono text-[11px] text-ink-2">No claims under the rest of the selection.</span>}
          {rows.map((r, i) => {
            const on = isSelected(selection, name, r.value)
            return (
              <div key={r.value} className={!long || open || i < FIRST ? 'contents' : 'hidden'} data-row={r.value}>
                <span className={`min-w-0 truncate text-xs ${on ? 'text-ink' : 'text-ink-2 line-through'}`} title={valueLabel(name, r.value)}
                  data-left-out={on ? undefined : name}>
                  <RowLabel name={name} value={r.value} />
                </span>
                {legacy ? <LegacyCell sums={r.legacy} dim={!on} selected={on} /> : <Gap column="legacy" />}
                {slotCells(sides, (s) => <SideCell sums={r.sums[s.side]} amount={hasVendorAmount(name, r.value)} dim={!on} />)}
              </div>
            )
          })}
        </div>
        {long && (
          <button type="button" onClick={() => setOpen((o) => !o)} aria-expanded={open} className="mt-2 text-[11px] text-link hover:underline">
            {open ? `Show the first ${FIRST}` : `Show all ${rows.length}`}
          </button>
        )}
      </div>
    </section>
  )
}

export function Breakdowns({ sides, selection }: { sides: OverviewSide[]; selection: OverviewSelection }) {
  return (
    <section data-section="breakdowns">
      <h2 className="mb-3 flex items-center gap-1 text-xs font-semibold text-ink">Breakdowns<Info term={help.terms.breakdowns} /></h2>
      <div className="space-y-4">
        {DIMENSIONS.map((d) => <Breakdown key={d.name} name={d.name} label={d.label} sides={sides} selection={selection} />)}
      </div>
    </section>
  )
}
