import { useEffect, useId, useLayoutEffect, useRef, useState, type MouseEvent, type ReactNode } from 'react'
import { createPortal } from 'react-dom'

// The hover/click popover for one code or name token. Adapted from the
// validator app's Info.tsx (frontend/src/app/Info.tsx): hovering or
// focusing opens it after a short delay, a click pins it open until a
// second click, Escape, or a click elsewhere, and it is placed in
// document.body with fixed coordinates so a maximized or scrolled box never
// clips it.

const GUTTER = 16 // px kept clear of the viewport edge
const GAP = 6 // px between the trigger and the popover
const DELAY = 140 // ms before hover opens or closes

type Mode = 'closed' | 'hover' | 'pinned'

export function Term({
  trigger,
  triggerClassName,
  title,
  body,
}: {
  trigger: ReactNode
  triggerClassName: string
  title: ReactNode
  body?: ReactNode
}) {
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
        onFocus={hoverIn}
        onBlur={hoverOut}
        aria-expanded={open}
        aria-describedby={open ? id : undefined}
        className={triggerClassName}
      >
        {trigger}
      </button>
      {open &&
        createPortal(
          <div
            ref={popover}
            id={id}
            role="tooltip"
            onMouseEnter={cancel}
            onMouseLeave={hoverOut}
            style={{ position: 'fixed', top: place?.top ?? 0, left: place?.left ?? 0, visibility: place ? 'visible' : 'hidden' }}
            className="z-[60] w-max max-w-[min(22rem,calc(100vw-2rem))] border border-line bg-bg px-3 py-2 text-left font-sans text-xs leading-5 font-normal normal-case text-ink shadow-lg"
          >
            <div className="font-semibold">{title}</div>
            {body && <div className="mt-0.5">{body}</div>}
          </div>,
          document.body,
        )}
    </>
  )
}
