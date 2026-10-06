import { submissionsOf, type Submission } from '../api'
import { SELECTION } from './help'
import { Info } from './Info'
import { slotNumber, useSelection, type Side } from './selection'

// The two controls over selection.tsx. VendorStrip sits in the navigation on
// every page: one submission selector per slot, each listing every vendor
// submission by vendor, dates oldest first, plus a "No submission" choice
// (the choice reads "No submission yet" for a slot whose default vendor
// has not submitted). Any two submissions can be picked, two dates of one
// vendor included. SideToggle sits under the title of each
// single-submission page and chooses which of the two picked submissions the
// page shows, offering only slots that have one.

const SELECT = 'border border-line bg-bg px-2 py-1 font-mono text-xs text-ink hover:border-ink-2'
const NO_SUBMISSION = '' // the submission <select>'s value for "No submission"

function Slot({ slot }: { slot: Side }) {
  const sel = useSelection()
  const current = slot === 'v1' ? sel.v1 : sel.v2
  // A slot never chosen and with no default vendor has nothing to show yet.
  const none = !current && sel.carry[slot] === undefined && sel.vendors[slot === 'v1' ? 0 : 1] === undefined ? 'No submission yet' : 'No submission'
  return (
    <div className="flex items-center gap-2">
      <SlotChip side={slot} />
      <select aria-label={`Submission ${slotNumber(slot)}`} value={current?.prefix ?? NO_SUBMISSION} onChange={(e) => sel.setSlot(slot, e.target.value)} className={SELECT}>
        <option value={NO_SUBMISSION}>{none}</option>
        {sel.vendors.map((label) => (
          <optgroup key={label} label={label}>
            {submissionsOf(sel.all, label).map((d) => <option key={d.prefix} value={d.prefix}>{label} · {d.submitted_on ?? d.prefix}</option>)}
          </optgroup>
        ))}
      </select>
    </div>
  )
}


export function VendorStrip() {
  const sel = useSelection()
  if (sel.error) return <span className="font-mono text-[11px] text-ink-2">Submissions could not be loaded: {sel.error.message}</span>
  if (sel.vendors.length === 0) {
    return <span className="font-mono text-[11px] text-ink-2">{sel.loading ? 'Loading submissions…' : 'No vendor submission has been compared yet'}</span>
  }
  return (
    <div className="flex flex-wrap items-center gap-x-8 gap-y-1">
      <Slot slot="v1" />
      <Slot slot="v2" />
      <Info term={SELECTION.picker} />
    </div>
  )
}

// A slot's number, used wherever a row, mark, or filter belongs to one of the
// two picked submissions.
export function SlotChip({ side }: { side: Side }) {
  return (
    <span className="inline-flex h-5 w-5 shrink-0 items-center justify-center border border-ink font-mono text-[11px] font-semibold text-ink">
      {slotNumber(side)}
    </span>
  )
}

// Selected is light and bold; unselected is grayed, as the app's other tabs.
export function SideToggle() {
  const sel = useSelection()
  const options = ([['v1', sel.v1], ['v2', sel.v2]] as const).filter((o): o is readonly [Side, Submission] => o[1] !== undefined)
  if (options.length === 0) return null
  const sameLabel = options.length === 2 && options[0][1].label === options[1][1].label
  return (
    <div className="flex items-center gap-2">
      <span className="text-xs text-ink-2">Showing</span>
      <div className="inline-flex border border-line" role="tablist" aria-label="Vendor shown">
        {options.map(([side, d]) => {
          const selected = sel.side === side
          return (
            <button key={side} type="button" role="tab" aria-selected={selected} onClick={() => sel.setSide(side)}
              className={`flex items-baseline gap-2 border-r border-line px-4 py-2 text-sm last:border-r-0 ${selected ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
              {d.label}
              <span className={`font-mono text-[11px] font-normal ${sameLabel ? 'text-ink' : 'text-ink-2'}`}>{d.submitted_on}</span>
            </button>
          )
        })}
      </div>
      <Info term={SELECTION.side} />
    </div>
  )
}
