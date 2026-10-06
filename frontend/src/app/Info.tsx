import { useEffect, useId, useLayoutEffect, useRef, useState, type MouseEvent } from 'react'
import { createPortal } from 'react-dom'
import type { Term } from './help'

// A small ⓘ beside a number, measure, chart title, or column header. Hovering
// shows the definition from help.ts; a click keeps it open until a second
// click, Escape, or a click elsewhere. The popover is placed in document.body
// with fixed coordinates, so a scrolling table or a clipped chart never cuts
// it off. The button stops its click there, so an ⓘ inside a clickable row
// or header opens the definition without sorting or navigating.

const GUTTER = 16 // px kept clear of the viewport edge
const GAP = 6 // px between the button and the popover
const DELAY = 140 // ms before hover opens or closes

type Mode = 'closed' | 'hover' | 'pinned'

export function Info({ term }: { term: Term }) {
  const id = useId()
  const button = useRef<HTMLButtonElement>(null)
  const popover = useRef<HTMLDivElement>(null)
  const timer = useRef<number | undefined>(undefined)
  const [mode, setMode] = useState<Mode>('closed')
  const [place, setPlace] = useState<{ top: number; left: number } | null>(null)
  const open = mode !== 'closed'

  const cancel = () => window.clearTimeout(timer.current)
  const later = (next: (m: Mode) => Mode) => {
    cancel()
    timer.current = window.setTimeout(() => setMode(next), DELAY)
  }
  const hoverIn = () => later((m) => (m === 'closed' ? 'hover' : m))
  const hoverOut = () => later((m) => (m === 'hover' ? 'closed' : m))

  const toggle = (e: MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    cancel()
    setMode((m) => (m === 'pinned' ? 'closed' : 'pinned'))
  }

  // Below the button, or above it when the viewport has no room below;
  // horizontally centered on it and kept inside the viewport.
  useLayoutEffect(() => {
    if (!open) return
    const update = () => {
      const b = button.current?.getBoundingClientRect()
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
      if (e.key !== 'Escape') return
      setMode('closed')
      if (popover.current?.contains(document.activeElement)) button.current?.focus()
    }
    const onDown = (e: PointerEvent) => {
      const target = e.target as Node
      if (button.current?.contains(target) || popover.current?.contains(target)) return
      setMode('closed')
    }
    document.addEventListener('keydown', onKey)
    document.addEventListener('pointerdown', onDown)
    return () => {
      document.removeEventListener('keydown', onKey)
      document.removeEventListener('pointerdown', onDown)
    }
  }, [open])

  useEffect(() => () => window.clearTimeout(timer.current), [])

  return (
    <>
      <button
        ref={button}
        type="button"
        onClick={toggle}
        onMouseEnter={hoverIn}
        onMouseLeave={hoverOut}
        aria-label={`About ${term.label}`}
        aria-expanded={open}
        aria-controls={open ? id : undefined}
        aria-describedby={open ? id : undefined}
        className={`inline-flex h-4 w-4 shrink-0 items-center justify-center align-middle hover:text-ink ${open ? 'text-ink' : 'text-ink-2'}`}
      >
        <svg viewBox="0 0 16 16" className="h-3 w-3" aria-hidden>
          <circle cx="8" cy="8" r="6.75" fill="none" stroke="currentColor" strokeWidth="1.3" />
          <circle cx="8" cy="4.9" r="1" fill="currentColor" />
          <path d="M8 7.3v4.4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
        </svg>
      </button>
      {open && createPortal(
        <div
          ref={popover}
          id={id}
          role="tooltip"
          onMouseEnter={cancel}
          onMouseLeave={hoverOut}
          style={{ position: 'fixed', top: place?.top ?? 0, left: place?.left ?? 0, visibility: place ? 'visible' : 'hidden' }}
          className="z-[60] w-max max-w-[min(20rem,calc(100vw-2rem))] border border-line bg-bg px-3 py-2 text-left font-sans text-xs leading-5 font-normal tracking-normal normal-case text-ink shadow-lg"
        >
          <div className="font-semibold">{term.label}</div>
          <div className="mt-0.5">{term.text}</div>
        </div>,
        document.body,
      )}
    </>
  )
}
