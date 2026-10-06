import { useCallback, useEffect, useId, useMemo, useRef, useState, type RefObject } from 'react'
import { useSearchParams } from 'react-router'
import { api, useApi, type MacsResponse } from '../api'
import { Failed, Loading } from '../components/ViewState'
import { MAC_MAP } from './help'
import { Info } from './Info'
import { lineText, macNodes, numbersIn, searchTree, type MacNode, type MacNodeKind } from './macLines'

// The MAC map: a panel at the right of the app's pages that draws every MAC
// workload number as a tree, from /api/macs (backend/pkg/macs Tree). It sits
// beside the page, never over it: a column as tall as the window that stays
// in view while the page scrolls and scrolls its own list, so opening it
// narrows the page's column by the panel's width (the docs site's lookup
// panel, viewer/LookupFrame.tsx, is the same arrangement). The top padding
// leaves the corner controls clear.
//
// Whether it is open is the address's macs=1, which every link carries
// (selection.tsx), so it stays open from page to page and a copied link opens
// it. The button in the corner controls (Shell.tsx) opens and closes it;
// Escape closes it, after an open help drawer or definition has taken its
// own Escape. Which jurisdictions are closed and the search are local to the
// panel.

export const MAC_MAP_ID = 'mac-map'

export interface MacMapState {
  open: boolean
  setOpen: (open: boolean) => void
  button: RefObject<HTMLButtonElement | null>
}

export function useMacMap(): MacMapState {
  const [params, setParams] = useSearchParams()
  const open = params.get('macs') === '1'
  const button = useRef<HTMLButtonElement>(null)

  // Closing from inside the panel puts focus back on the button, since the
  // panel leaves the page.
  const setOpen = useCallback((on: boolean) => {
    const inside = !on && (document.getElementById(MAC_MAP_ID)?.contains(document.activeElement) ?? false)
    setParams((prev) => {
      const next = new URLSearchParams(prev)
      if (on) next.set('macs', '1')
      else next.delete('macs')
      return next
    })
    if (inside) button.current?.focus()
  }, [setParams])

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== 'Escape' || e.defaultPrevented) return
      if (document.querySelector('[role="dialog"]:not([inert]), [role="tooltip"]')) return
      setOpen(false)
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open, setOpen])

  return { open, setOpen, button }
}

// The button in the corner controls.
export function MacMapButton({ map }: { map: MacMapState }) {
  return (
    <button ref={map.button} type="button" onClick={() => map.setOpen(!map.open)} aria-pressed={map.open}
      aria-controls={map.open ? MAC_MAP_ID : undefined}
      className={`shrink-0 border px-2 py-1 font-mono text-[11px] ${map.open ? 'border-ink bg-ink text-bg' : 'border-line bg-bg text-ink-2 hover:text-ink'}`}>
      MAC map
    </button>
  )
}

// The panel's column, beside the page's (Shell.tsx). The tree is fetched
// when the panel opens; the response is cacheable, so opening it again reads
// the browser's copy.
export function MacMapPanel({ onClose }: { onClose: () => void }) {
  const macs = useApi('macs', (signal) => api.macs(signal))
  return (
    <aside id={MAC_MAP_ID} aria-label="MAC map" className="sticky top-0 flex h-svh w-[22rem] shrink-0 flex-col border-l border-line bg-bg pt-14">
      <MacTreeView data={macs.data} error={macs.error} onClose={onClose} />
    </aside>
  )
}

// A folder, for a roll-up. In the tree it sits at the end of the roll-up's
// branch, in the space a workload's branch crosses, so the roll-up's number
// starts where the numbers beside it start; the overview's workload box
// (SelectionHeader.tsx) heads its roll-ups with it too.
export function Folder({ className }: { className: string }) {
  return (
    <svg viewBox="0 0 16 16" className={className} aria-hidden>
      <path d="M1.75 3.25h4.5l1.5 1.75h6.5v7.75H1.75z" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinejoin="round" />
    </svg>
  )
}

// One line of the tree, for anything under a jurisdiction. data-node holds
// the line as the drawing writes it (macLines.ts lineText), for the SSR check.
// Every node is one line and never wraps: the tree's shape carries the
// meaning, and the longer text (state names, a note, the codes) is the hover.
//
// How a line is set, so the number is what the eye finds first:
//   - The number leads, in monospace one size above the text. Every number
//     at one depth starts at the same place and is five digits wide, so the
//     numbers stand in a column and the text after them starts at a second.
//   - A Part's heading is the small capitals label the app's sections use,
//     quiet beside the numbers.
//   - A roll-up is a tinted band holding its folder, number, and states, the
//     way a table sets a group's heading row; its workloads hang from its
//     number. A workload is plain, its branch running up to its number.
//   - The aside sits at the end of the line, joined by a dotted leader.
//   - Every line is 22px tall, as the text sets it; the number's line box is
//     the smaller one, so the larger number adds no height.
function Line({ node }: { node: MacNode }) {
  if (node.kind === 'part') {
    return (
      <div className="flex items-baseline gap-3 whitespace-nowrap py-0.5 leading-[18px]" title={node.detail}>
        <span className="shrink-0 font-mono text-[11px] font-semibold uppercase tracking-wider text-ink-2">{node.title}</span>
        <span className="min-w-0 truncate text-[11px] text-ink-2">{node.detail}</span>
      </div>
    )
  }
  const names = node.states.map((s) => `${s.code} ${s.name}`).join('\n')
  // An HHH number's region is in its line's text and nowhere else, so its
  // hover starts with that text, for when the end of the line is cut off.
  const hover = [node.kind === 'hhh' ? node.detail : '', names, node.codes, node.note].filter(Boolean).join('\n') || undefined
  const rollup = node.kind === 'rollup'
  return (
    <div title={hover}
      className={`flex items-baseline gap-2 whitespace-nowrap py-0.5 leading-[18px] ${rollup ? 'relative -ml-4 bg-surface pl-4 pr-2' : ''}`}>
      {rollup && <Folder className="absolute left-0.5 top-1/2 h-3 w-3 -translate-y-1/2 text-ink-2" />}
      <span className="shrink-0 font-mono text-[13px] leading-4 text-ink">{node.title}</span>
      {rollup && <span className="-mr-0.5 shrink-0 text-[11px] text-ink-2">roll-up</span>}
      {node.kind === 'hhh' && (
        <span className="-mr-0.5 shrink-0 self-center bg-accent/15 px-1 font-mono text-[10px] leading-4 text-link">HHH</span>
      )}
      <span className={`min-w-0 truncate text-[12px] ${rollup ? 'text-ink-2' : 'text-ink'}`}>{node.detail}</span>
      {node.aside && (
        <>
          <span aria-hidden className="mb-1 min-w-3 flex-1 self-end border-b border-dotted border-line" />
          <span className="shrink-0 text-[11px] text-ink-2">{node.aside}</span>
        </>
      )}
    </div>
  )
}

// Where each kind of line starts, and the length of its tick. The text sits
// on a 20px grid from the jurisdiction's code: a Part's heading starts under
// the code, and a number starts 20px in from the line it hangs from. A
// workload's tick stops 4px short of its number; a roll-up's stops at its
// band.
const PLACE: Record<MacNodeKind, string> = {
  jurisdiction: '',
  part: 'pl-3.5 after:w-2.5',
  rollup: 'pl-5 after:w-1',
  workload: 'pl-5 after:w-4',
  hhh: 'pl-5 after:w-4',
  other: 'pl-5 after:w-4',
}

// The lines under a jurisdiction, with the drawing's branches: each line
// carries its own piece of the trunk and a tick to it, level with the
// number's middle, and the last line's piece stops at its tick.
function Branches({ nodes, depth }: { nodes: MacNode[]; depth: number }) {
  return (
    <ul>
      {nodes.map((node) => (
        <li key={node.key} data-node={lineText(node)} data-depth={depth}
          className={`relative before:absolute before:left-0 before:top-0 before:h-full before:border-l before:border-line last:before:h-3 after:absolute after:left-0 after:top-[11px] after:border-t after:border-line ${PLACE[node.kind]}`}>
          <Line node={node} />
          {node.children.length > 0 && <Branches nodes={node.children} depth={depth + 1} />}
        </li>
      ))}
    </ul>
  )
}

const DIGITS: [string, string][] = [
  ['11', 'the jurisdiction, as CMS numbered them in 2007 (11 is JM)'],
  ['2', 'the state within it; 0 in a roll-up, HHH, or DME number'],
  ['0', 'changes when another contractor takes the work over; a few Part B numbers use it to split a state'],
  ['1', 'the claim type: 1 Part A, 2 Part B, 3 DME, 4 home health and hospice (HHH)'],
]

// The panel's content from the /api/macs response; exported for the SSR
// check, which renders it from a fixture.
export function MacTreeView({ data, error, onClose }: { data: MacsResponse | undefined; error?: Error; onClose: () => void }) {
  const [query, setQuery] = useState('')
  const [closed, setClosed] = useState<Set<string>>(() => new Set())
  const list = useRef<HTMLDivElement>(null)
  const searchId = useId()
  const all = useMemo(() => (data ? macNodes(data) : []), [data])
  const shown = useMemo(() => searchTree(all, query), [all, query])
  const total = numbersIn(all)
  const searching = query.trim() !== ''
  useEffect(() => {
    list.current?.scrollTo({ top: 0 })
  }, [query])

  const flip = (key: string) => setClosed((prev) => {
    const next = new Set(prev)
    if (!next.delete(key)) next.add(key)
    return next
  })

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="flex items-start justify-between gap-3 border-b border-line px-4 pb-3">
        <div>
          <h2 className="flex items-center gap-1 font-serif text-lg text-ink">MAC map<Info term={MAC_MAP} /></h2>
          <p className="mt-0.5 text-[11px] leading-4 text-ink-2">Every Medicare Administrative Contractor workload number, by jurisdiction and Part.</p>
        </div>
        <button type="button" onClick={onClose} aria-label="Close the MAC map" className="shrink-0 border border-line px-2 py-0.5 font-mono text-[11px] text-ink-2 hover:text-ink">
          close
        </button>
      </div>

      <div className="border-b border-line px-4 py-3">
        <label htmlFor={searchId} className="mb-1.5 block text-xs font-semibold text-ink">Find a number, state, or contractor</label>
        <input id={searchId} type="search" value={query} onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => {
            // The first Escape empties the search; the next closes the panel.
            if (e.key === 'Escape' && query !== '') {
              e.preventDefault()
              setQuery('')
            }
          }}
          placeholder="Try 11201, Virginia, or Palmetto"
          className="w-full border border-line bg-bg px-3 py-1.5 font-mono text-[13px] text-ink focus:outline-2 focus:outline-accent" />
        <div className="mt-1.5 flex items-center gap-3 text-[11px] text-ink-2">
          <span role="status">{searching ? `${numbersIn(shown)} of ${total} numbers` : `${total} numbers in ${all.length} jurisdictions`}</span>
          {!searching && data && (
            <span className="ml-auto flex gap-3">
              <button type="button" onClick={() => setClosed(new Set())} className="text-link hover:underline">Open all</button>
              <button type="button" onClick={() => setClosed(new Set(all.map((j) => j.key)))} className="text-link hover:underline">Close all</button>
            </span>
          )}
        </div>
      </div>

      <div ref={list} className="min-h-0 flex-1 overflow-y-auto px-4 py-1">
        {error ? <Failed error={error} /> : !data ? <Loading what="Loading the MAC numbers" /> : shown.length === 0 ? (
          <p className="py-6 text-[13px] text-ink-2">Nothing matches “{query}”.</p>
        ) : (
          <ul>
            {shown.map((j) => {
              // A search shows its matches in every jurisdiction, closed or not.
              const open = searching || !closed.has(j.key)
              return (
                <li key={j.key} data-node={lineText(j)} data-depth={0} className="border-b border-line py-1.5 last:border-b-0">
                  <button type="button" onClick={() => flip(j.key)} aria-expanded={open} disabled={searching}
                    className="flex w-full items-baseline gap-2 whitespace-nowrap text-left" title={j.detail}>
                    <span className="w-3 shrink-0 text-[11px] text-ink-2" aria-hidden>{open ? '▾' : '▸'}</span>
                    <span className="shrink-0 font-mono text-[13px] font-semibold text-ink">{j.title}</span>
                    <span className="min-w-0 truncate text-[11px] text-ink-2">{j.detail}</span>
                  </button>
                  {open && (
                    <div className="ml-1.5 mt-0.5">
                      <Branches nodes={j.children} depth={1} />
                    </div>
                  )}
                </li>
              )
            })}
          </ul>
        )}
      </div>

      <div className="border-t border-line px-4 py-3 text-[11px] leading-4 text-ink-2">
        <div className="mb-1.5 font-semibold text-ink">
          How a number reads, for <span className="font-mono">11201</span>
        </div>
        <dl className="grid grid-cols-[2rem_minmax(0,1fr)] gap-x-2 gap-y-0.5">
          {DIGITS.map(([digits, meaning]) => (
            <div key={meaning} className="contents">
              <dt className="font-mono text-ink">{digits}</dt>
              <dd>{meaning}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  )
}
