import { useState } from 'react'
import type { FrontFilter, MacNumber, MacsResponse } from '../api'
import { HELP } from './help'
import { Info } from './Info'
import { slotNumber } from './selection'
import type { OverviewSide } from './overviewColumns'
import {
  ADDRESS_ONLY, BOX_DIMENSIONS, boxValues, dimensionLabel, isNarrowed, selectAll, selectNone, tabLabel, toggleValue, valueLabel, type BoxValue, type OverviewSelection,
} from './overviewFilters'
import { statesText, workloadGroups } from './macLines'
import { Folder } from './MacTree'

// The overview's selection header: what the numbers are over, one box per
// dimension (overviewFilters.ts BOX_DIMENSIONS). It sticks to the top of the
// window while the page scrolls under it, so a change shows in the tables
// below without scrolling back. It stays in the page's flow, so the content
// starts where it ends, and it is never taller than 40% of the window: past
// that it scrolls its own boxes.
//
// A box is the dimension's name with All and None at its top and its values
// as a tight grid of small tabs, as many to a row as COLUMNS says: the
// jurisdictions four to a row in CMS's order, the others in a row or a short
// grid. The boxes sit side by side and wrap. A click on a value picks it and
// a second click unpicks it (overviewFilters.ts toggleValue); with nothing
// picked a box counts every value, which All restores, and None leaves every
// value out. A box lists only the values with claims under what the other
// boxes pick (overviewFilters.ts boxValues), so picking J5 leaves only J5's
// workloads; a picked value with none stays, dimmed, so it can be unpicked.
// A tab carries no count.
//
// The ninth box, Workload, starts closed: one line with how many workloads
// it lists, or how many are picked. Open, it lists the workloads under their
// jurisdiction codes as chips smaller than the other boxes' tabs, each
// roll-up a tinted group headed by its folder and number as the MAC map draws
// it (macLines.ts workloadGroups, from /api/macs). The status line beside
// the header's title says when new numbers are on their way, or a request
// failed.

const help = HELP.overview

// How many tabs a box sets in a row; a box with fewer values is as wide as they are.
const COLUMNS: Partial<Record<FrontFilter, number>> = {
  jurisdiction: 4,
  claim_type: 3,
  adjudication: 3,
  mac_status: 8,
  vendor_status: 8,
  frequency: 4,
  service_month: 12,
  outcome: 6,
}

// The dimensions whose values are codes, set in monospace.
const CODED = new Set<FrontFilter>(['jurisdiction', 'mac_status', 'vendor_status', 'frequency', 'service_month', 'workload'])

const TAB = 'border px-1.5 py-0.5 text-[11px] leading-4 whitespace-nowrap'
const CHIP = 'border px-1 text-[10px] leading-4 whitespace-nowrap' // the workload box's, a size smaller

// A picked value is dark; picked with no claims under the rest of the
// selection, dark and dimmed; not picked, plain.
const valueTone = (v: BoxValue) =>
  v.picked
    ? v.claims > 0 ? 'border-ink bg-ink font-semibold text-bg' : 'border-ink/30 bg-ink/30 font-semibold text-bg'
    : 'border-line bg-bg text-ink-2 hover:border-ink-2 hover:text-ink'

// All and None are outlined while they hold.
const endTone = (on: boolean) => (on ? 'border-ink font-semibold text-ink' : 'border-transparent text-ink-2 hover:text-ink')

interface BoxProps {
  name: FrontFilter
  values: BoxValue[]
  selection: OverviewSelection
  onSelection: (next: OverviewSelection) => void
}

function Ends({ name, selection, onSelection }: Omit<BoxProps, 'values'>) {
  const list = selection[name]
  const none = list !== undefined && list.length === 0
  return (
    <span className="ml-auto flex shrink-0 gap-0.5">
      <button type="button" data-tab="all" aria-pressed={list === undefined} onClick={() => onSelection(selectAll(selection, name))}
        className={`${TAB} ${endTone(list === undefined)}`}>
        All
      </button>
      <button type="button" data-tab="none" aria-pressed={none} onClick={() => onSelection(selectNone(selection, name))}
        className={`${TAB} ${endTone(none)}`}>
        None
      </button>
    </span>
  )
}

function ValueTab({ name, v, selection, onSelection, chip = false, title }: Omit<BoxProps, 'values'> & { v: BoxValue; chip?: boolean; title?: string }) {
  const label = title ?? valueLabel(name, v.value)
  const empty = v.picked && v.claims === 0
  return (
    <button type="button" data-tab={`${name}:${v.value}`} aria-pressed={v.picked} data-zero={empty ? '' : undefined}
      title={empty ? `${label}: no claims under the rest of the selection` : label}
      onClick={() => onSelection(toggleValue(selection, name, v.value))}
      className={[chip ? CHIP : TAB, CODED.has(name) ? 'font-mono' : '', valueTone(v)].filter(Boolean).join(' ')}>
      {tabLabel(name, v.value)}
    </button>
  )
}

const Nothing = () => <p className="px-1.5 pb-1.5 text-[11px] text-ink-2">No values with claims</p>

function Box({ name, values, selection, onSelection }: BoxProps) {
  const label = dimensionLabel(name)
  const columns = Math.max(1, Math.min(COLUMNS[name] ?? 4, values.length))
  return (
    <div role="group" aria-label={label} data-tab-group={name} className="border border-line">
      <div className="flex items-center gap-3 px-1.5 pt-1">
        <span className="text-[11px] font-semibold whitespace-nowrap text-ink">{label}</span>
        <Ends name={name} selection={selection} onSelection={onSelection} />
      </div>
      {values.length === 0 ? <Nothing /> : (
        <div className="grid gap-0.5 p-1.5" style={{ gridTemplateColumns: `repeat(${columns}, ${name === 'jurisdiction' ? '1fr' : 'auto'})` }} data-columns={columns}>
          {values.map((v) => <ValueTab key={v.value} name={name} v={v} selection={selection} onSelection={onSelection} />)}
        </div>
      )}
    </div>
  )
}

// A workload number's hover: the number and its states.
const numberTitle = (n: MacNumber) => `${n.id}: ${statesText(n.states)}`

export function WorkloadBox({ values, selection, onSelection, macs, open, onOpen }: Omit<BoxProps, 'name'> & {
  macs: MacsResponse | undefined
  open: boolean
  onOpen: (open: boolean) => void
}) {
  const list = selection.workload
  const count = list === undefined
    ? values.length > 0 ? `all ${values.length}` : 'none with claims'
    : list.length === 0 ? 'none' : `${list.length} selected`
  const head = (
    <button type="button" aria-expanded={open} onClick={() => onOpen(!open)} data-workload-toggle=""
      className="flex items-center gap-1.5 text-[11px] whitespace-nowrap">
      <span aria-hidden className="w-2.5 text-ink-2">{open ? '▾' : '▸'}</span>
      <span className="font-semibold text-ink">Workload:</span>
      <span className="text-ink-2" data-workload-count="">{count}</span>
    </button>
  )
  if (!open) {
    return (
      <div role="group" aria-label="Workload" data-tab-group="workload" data-open="0" className="border border-line px-1.5 py-1">
        {head}
      </div>
    )
  }

  const byValue = new Map(values.map((v) => [v.value, v]))
  const chip = (id: string, title?: string) => {
    const v = byValue.get(id)
    return v && <ValueTab key={id} name="workload" v={v} selection={selection} onSelection={onSelection} chip title={title} />
  }
  const { groups, unplaced } = macs ? workloadGroups(macs, values.map((v) => v.value)) : { groups: [], unplaced: values.map((v) => v.value) }
  return (
    <div role="group" aria-label="Workload" data-tab-group="workload" data-open="1" className="basis-full border border-line">
      <div className="flex items-center gap-1.5 px-1.5 pt-1">
        {head}
        <Info term={help.terms.workloads} />
        <Ends name="workload" selection={selection} onSelection={onSelection} />
      </div>
      {values.length === 0 ? <Nothing /> : (
        <div className="grid gap-x-5 gap-y-1 p-1.5" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(12rem, 1fr))' }}>
          {groups.map((g) => (
            <div key={g.jurisdiction} className="flex items-start gap-1.5" data-workload-jurisdiction={g.jurisdiction}>
              <span className="w-8 shrink-0 font-mono text-[11px] leading-[18px] font-semibold text-ink" title={g.contractor}>{g.jurisdiction}</span>
              <span className="flex min-w-0 flex-wrap items-center gap-0.5">
                {g.items.map((item) => item.kind === 'rollup' ? (
                  <span key={item.id} data-rollup={item.id} title={`Roll-up ${item.id}: ${statesText(item.rollUp.states)}`}
                    className="inline-flex flex-wrap items-center gap-0.5 bg-surface py-px pr-0.5 pl-1">
                    <Folder className="h-3 w-3 shrink-0 text-ink-2" />
                    {item.own ? chip(item.id, `Roll-up ${numberTitle(item.rollUp)}`) : <span className="pr-0.5 font-mono text-[10px] text-ink-2">{item.id}</span>}
                    {item.workloads.map((w) => chip(w.id, numberTitle(w)))}
                  </span>
                ) : chip(item.id, numberTitle(item.number)))}
              </span>
            </div>
          ))}
          {unplaced.length > 0 && (
            <div className="flex items-start gap-1.5" data-workload-jurisdiction="">
              {macs && <span className="w-8 shrink-0 text-[11px] leading-[18px] text-ink-2" title="Numbers the MAC map does not list">other</span>}
              <span className="flex min-w-0 flex-wrap items-center gap-0.5">{unplaced.map((id) => chip(id))}</span>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

// What an address-only dimension is narrowed to, in a few words.
function narrowedTo(name: FrontFilter, list: string[]): string {
  if (list.length === 0) return 'none'
  if (list.length > 3) return `${list.length} values`
  return list.map((v) => valueLabel(name, v)).join(', ')
}

export function SelectionHeader({ sides, selection, onSelection, macs, updating = false }: {
  sides: OverviewSide[]
  selection: OverviewSelection
  onSelection: (next: OverviewSelection) => void
  macs?: MacsResponse // the MAC map's tree, for the workload box's groups
  updating?: boolean // newer numbers are on their way
}) {
  const [workloadsOpen, setWorkloadsOpen] = useState(false)
  const fromAddress = ADDRESS_ONLY.filter((name) => selection[name] !== undefined)
  const clearAddress = () => {
    const next = { ...selection }
    for (const name of fromAddress) delete next[name]
    onSelection(next)
  }
  const valuesOf = (name: FrontFilter) => boxValues(sides.map((s) => s.data?.dimensions.find((d) => d.name === name)?.values ?? []), name, selection)
  return (
    <section aria-label="Selection" data-tab-row="" className="sticky top-0 z-30 max-h-[40svh] overflow-y-auto border-b border-line bg-bg pt-4 pb-3">
      {/* The title line leaves its right end clear: the page's corner controls
          (Shell.tsx) sit over the header's top right while it is stuck, and
          the boxes start below them. */}
      <div className="mb-2.5 flex flex-wrap items-center gap-x-3 gap-y-1 pr-80">
        <span className="flex items-center gap-1 text-xs font-semibold text-ink">Selection<Info term={help.terms.selection} /></span>
        <button type="button" onClick={() => onSelection({})} disabled={!isNarrowed(selection)}
          className="text-[11px] text-link hover:underline disabled:text-ink-2 disabled:no-underline disabled:opacity-50">
          Reset selection
        </button>
        {updating && <span className="font-mono text-[11px] text-ink-2" role="status">updating…</span>}
        {sides.map((s) => s.error && <span key={s.side} className="font-mono text-[11px] text-danger" role="alert">Submission {slotNumber(s.side)}: {s.error.message}</span>)}
      </div>
      {/* The jurisdiction box is the tall one, so it stands at the left and the
          short boxes wrap in the space beside it, two rows deep, instead of
          leaving that space empty. */}
      <div className="flex items-start gap-2">
        <Box name="jurisdiction" values={valuesOf('jurisdiction')} selection={selection} onSelection={onSelection} />
        <div className="flex min-w-0 flex-1 flex-wrap items-start gap-2">
          {BOX_DIMENSIONS.filter((name) => name !== 'jurisdiction').map((name) => name === 'workload'
            ? <WorkloadBox key={name} values={valuesOf(name)} selection={selection} onSelection={onSelection} macs={macs} open={workloadsOpen} onOpen={setWorkloadsOpen} />
            : <Box key={name} name={name} values={valuesOf(name)} selection={selection} onSelection={onSelection} />)}
        </div>
      </div>
      {fromAddress.length > 0 && (
        <p className="mt-2.5 flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px] text-ink-2" data-address-only="">
          <span>Also narrowed by the address:</span>
          {fromAddress.map((name) => (
            <span key={name} className="text-ink">{dimensionLabel(name)}: {narrowedTo(name, selection[name] ?? [])}</span>
          ))}
          <Info term={help.terms.addressOnly} />
          <button type="button" onClick={clearAddress} className="text-link hover:underline">Clear</button>
        </p>
      )}
    </section>
  )
}
