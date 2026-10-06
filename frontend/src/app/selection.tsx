import { createContext, useContext, useEffect, type ReactNode } from 'react'
import { useSearchParams } from 'react-router'
import { submissionsOf, useSubmissions, vendorLabels, type Submission } from '../api'
import { withParams } from '../format'

// Which two vendor submissions the app shows. The navigation picks a
// submission for each of two slots, 1 and 2, kept in the address as v1=
// and v2= (submission prefixes). Each slot can hold any vendor submission,
// from any vendor: the two can be the latest of two vendors, two dates of
// one vendor, or the same submission twice. Slot 1 defaults to the first
// vendor's latest submission and slot 2 to the second's. A slot can also be
// explicitly empty (v1= present with no value): the user picked "No
// submission" for it, which is kept in the address and every link, unlike a
// slot absent from the address, which still means "not chosen" and defaults
// to the vendor's latest (so old links keep working). The single-submission
// pages show one slot, chosen with the toggle under their title and kept
// as side=2 (slot 1 when absent); when the chosen side's slot is empty,
// the other one is shown instead. The Overview and the claim page show both,
// an empty slot's column left empty. With both slots empty, every page asks
// the user to pick a submission instead of showing its content
// (Selection.unselected). Every link carries v1, v2, and side, so a copied
// address reproduces the view.
//
// /api/submissions is read once, here, for the whole app. Once it arrives,
// the address is rewritten in place to name both slots.

// A slot: v1 or v2 in the code and the address, 1 or 2 on the screen.
export type Side = 'v1' | 'v2'

// The number that names a slot on the screen.
export const slotNumber = (side: Side) => (side === 'v1' ? '1' : '2')

export interface Selection {
  all: Submission[] // every submission, the MAC's included
  vendors: string[] // vendor labels in the server's order
  v1: Submission | undefined
  v2: Submission | undefined
  side: Side
  shown: Submission | undefined // the slot the single-submission pages show
  // True once submissions have loaded, at least one vendor exists, and
  // neither slot is picked: the user needs to choose before a page can
  // show anything, distinct from no vendor having submitted at all.
  unselected: boolean
  loading: boolean
  error: Error | undefined
  carry: Record<string, string | undefined> // v1, v2, side, and macs, for links
  // The page at `path` (/, /claims, ...) with carry and `params` in its query.
  link: (path: string, params?: Record<string, string | undefined | null>) => string
  setSlot: (slot: Side, prefix: string) => void // any vendor submission's prefix, or '' for "No submission"
  setSide: (side: Side) => void
}

const Context = createContext<Selection | undefined>(undefined)

const vendorAt = (all: Submission[], prefix: string | null) =>
  prefix === null ? undefined : all.find((d) => d.side === 'vendor' && d.prefix === prefix)

function latest(all: Submission[], label: string | undefined) {
  if (label === undefined) return undefined
  const own = submissionsOf(all, label)
  return own[own.length - 1]
}

function resolve(all: Submission[], params: URLSearchParams) {
  const vendors = vendorLabels(all)
  const v1Raw = params.get('v1')
  const v2Raw = params.get('v2')
  let v1 = vendorAt(all, v1Raw)
  let v2 = vendorAt(all, v2Raw)
  let side: Side = params.get('side') === '2' ? 'v2' : 'v1'
  // v1= present with no value is an explicit "No submission": it is never
  // defaulted. An absent key, or one naming a submission that no longer
  // exists, still falls back to the vendor's latest.
  if (v1Raw !== '' && !v1) v1 = latest(all, vendors[0])
  if (v2Raw !== '' && !v2) v2 = latest(all, vendors[1])
  if (!v2) side = 'v1'
  // The chosen side's slot may be the empty one: show the other when it
  // has a submission.
  if (side === 'v1' && !v1 && v2) side = 'v2'
  return { vendors, v1, v2, side, v1Empty: v1Raw === '' && !v1, v2Empty: v2Raw === '' && !v2 }
}

// The address with both slots named; other parameters (match_key, filters)
// are kept in place. An explicit empty slot is kept as a bare v1= or v2=,
// not dropped, so the choice survives.
function canonical(params: URLSearchParams, s: ReturnType<typeof resolve>) {
  const next = new URLSearchParams(params)
  for (const [slot, picked, empty] of [['v1', s.v1, s.v1Empty], ['v2', s.v2, s.v2Empty]] as const) {
    if (picked) next.set(slot, picked.prefix)
    else if (empty) next.set(slot, '')
    else next.delete(slot)
  }
  if (s.side === 'v2') next.set('side', '2')
  else next.delete('side')
  return next
}

export function SelectionProvider({ children }: { children: ReactNode }) {
  const [params, setParams] = useSearchParams()
  const submissions = useSubmissions()
  const all = submissions.data?.submissions ?? []
  const resolved = resolve(all, params)
  const { v1, v2, side, vendors, v1Empty, v2Empty } = resolved

  const loaded = submissions.data !== undefined
  const text = params.toString()
  useEffect(() => {
    if (!loaded) return
    const next = canonical(params, resolve(all, params))
    if (next.toString() !== text) setParams(next, { replace: true })
    // `text` carries `params`; `all` changes only when `loaded` does.
  }, [loaded, text])

  // An explicit empty slot carries as '' (kept by format.ts withParams),
  // not omitted like an unpicked one, so it survives every link. macs=1,
  // the MAC map's open state (MacTree.tsx), is carried the same way, so the
  // map stays open from page to page.
  const carry = {
    v1: v1 ? v1.prefix : v1Empty ? '' : undefined,
    v2: v2 ? v2.prefix : v2Empty ? '' : undefined,
    side: side === 'v2' ? '2' : undefined,
    macs: params.get('macs') === '1' ? '1' : undefined,
  }
  // A new choice keeps the page's other parameters but not its page number.
  const update = (change: (next: URLSearchParams) => void) =>
    setParams((prev) => {
      const next = canonical(prev, resolve(all, prev))
      change(next)
      next.delete('page')
      return next
    })

  const value: Selection = {
    all,
    vendors,
    v1,
    v2,
    side,
    shown: side === 'v2' ? v2 : v1,
    unselected: vendors.length > 0 && !v1 && !v2,
    loading: submissions.loading,
    error: submissions.error,
    carry,
    link: (path, extra = {}) => withParams(path, { ...carry, ...extra }),
    setSlot: (slot, prefix) => update((next) => next.set(slot, prefix)),
    setSide: (s) => update((next) => (s === 'v2' ? next.set('side', '2') : next.delete('side'))),
  }
  return <Context.Provider value={value}>{children}</Context.Provider>
}

export function useSelection(): Selection {
  const selection = useContext(Context)
  if (!selection) throw new Error('useSelection needs SelectionProvider (App.tsx)')
  return selection
}
