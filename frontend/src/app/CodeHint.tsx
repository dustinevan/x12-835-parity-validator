import { useEffect, useId, useLayoutEffect, useRef, useState, type ReactNode } from 'react'
import { createPortal } from 'react-dom'
import { api, type CodesResponse } from '../api'
import { describeCarc, describeElement, describeRarc, describeStatus, type CodeExplanation } from '../data/codeHints'

// A hover/focus popover on an 835 code the app already shows: an element
// tag (CLP04, CAS*CO*253, REF*6R, MOA), a CARC or RARC code, or a CLP02
// claim status value. It reads like Info.tsx's ⓘ: a title of a few words
// saying what the thing is, then one short paragraph (data/codeHints.ts
// writes both), with no citations. It wraps the code's own text instead of
// adding an icon — the only visual change is a dotted underline — and never
// intercepts a click, so a code that is also a link or a filter chip still
// works exactly as before.

const GUTTER = 16 // px kept clear of the viewport edge
const GAP = 6 // px between the code and the popover
const DELAY = 140 // ms before hover opens or closes

// /api/codes is static reference data (docs/API.md): fetched once, lazily,
// on the first hover or focus of any code, and shared by every CodeHint.
let codesCache: CodesResponse | undefined
let codesPromise: Promise<void> | undefined
const subscribers = new Set<() => void>()

function ensureCodes() {
  if (codesCache || codesPromise) return
  codesPromise = api.codes().then(
    (data) => {
      codesCache = data
      codesPromise = undefined
      subscribers.forEach((notify) => notify())
    },
    () => {
      codesPromise = undefined
    },
  )
}

// Subscribes to the shared CARC/RARC cache only while `active` (the popover
// is open), so a page full of codes never fetches until one is actually
// shown.
function useCodesCache(active: boolean): CodesResponse | undefined {
  const [, rerender] = useState(0)
  useEffect(() => {
    if (!active) return
    ensureCodes()
    const notify = () => rerender((n) => n + 1)
    subscribers.add(notify)
    return () => {
      subscribers.delete(notify)
    }
  }, [active])
  return codesCache
}

export type CodeKind = 'element' | 'carc' | 'rarc' | 'status'

function explain(kind: CodeKind, code: string, label: string | undefined, codes: CodesResponse | undefined): CodeExplanation {
  switch (kind) {
    case 'carc':
      return describeCarc(code, codes)
    case 'rarc':
      return describeRarc(code, codes)
    case 'status':
      return describeStatus(code, label ?? code)
    case 'element':
      return describeElement(code, label, codes)
  }
}

export function CodeHint({ kind, code, label, children }: { kind: CodeKind; code: string; label?: string; children: ReactNode }) {
  const id = useId()
  const trigger = useRef<HTMLSpanElement>(null)
  const popover = useRef<HTMLDivElement>(null)
  const timer = useRef<number | undefined>(undefined)
  const [open, setOpen] = useState(false)
  const [place, setPlace] = useState<{ top: number; left: number } | null>(null)
  const codes = useCodesCache(open)

  const cancel = () => window.clearTimeout(timer.current)
  const later = (next: boolean) => {
    cancel()
    timer.current = window.setTimeout(() => setOpen(next), DELAY)
  }
  const show = () => later(true)
  const hide = () => later(false)

  useLayoutEffect(() => {
    if (!open) return
    const update = () => {
      const b = trigger.current?.getBoundingClientRect()
      const p = popover.current
      if (!b || !p) return
      const w = p.offsetWidth
      const h = p.offsetHeight
      const left = Math.max(GUTTER, Math.min(b.left + b.width / 2 - w / 2, window.innerWidth - w - GUTTER))
      const below = b.bottom + GAP
      const top = below + h > window.innerHeight - GUTTER && b.top - GAP - h > GUTTER ? b.top - GAP - h : below
      setPlace({ top, left })
    }
    update()
    window.addEventListener('scroll', update, true)
    window.addEventListener('resize', update)
    return () => {
      window.removeEventListener('scroll', update, true)
      window.removeEventListener('resize', update)
    }
  }, [open])

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open])

  useEffect(() => () => window.clearTimeout(timer.current), [])

  const explanation = open ? explain(kind, code, label, codes) : undefined

  return (
    <span
      ref={trigger}
      tabIndex={0}
      onMouseEnter={show}
      onMouseLeave={hide}
      onFocus={show}
      onBlur={hide}
      aria-describedby={open ? id : undefined}
      className="cursor-default underline decoration-dotted decoration-ink-2 underline-offset-2 outline-none"
    >
      {children}
      {open && explanation && createPortal(
        <div
          ref={popover}
          id={id}
          role="tooltip"
          onMouseEnter={cancel}
          onMouseLeave={hide}
          style={{ position: 'fixed', top: place?.top ?? 0, left: place?.left ?? 0, visibility: place ? 'visible' : 'hidden' }}
          className="z-[60] w-max max-w-[min(22rem,calc(100vw-2rem))] border border-line bg-bg px-3 py-2 text-left font-sans text-xs leading-5 font-normal tracking-normal normal-case text-ink shadow-lg"
        >
          <div className="font-semibold">
            {explanation.title}
            {explanation.deprecated && <span className="ml-1.5 border border-warn px-1 text-[10px] font-normal text-warn">not current</span>}
          </div>
          {explanation.body && <div className="mt-0.5 max-h-56 overflow-y-auto">{explanation.body}</div>}
          {explanation.loading && <div className="mt-0.5 text-ink-2">Loading the code’s text…</div>}
        </div>,
        document.body,
      )}
    </span>
  )
}
