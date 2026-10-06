import { useEffect, useId, useRef, type RefObject } from 'react'
import type { PageHelp } from './help'

// "How to read this page": a panel that slides over the right edge of the
// window, so the page layout does not move. Opening moves focus to its
// close button; Escape or the close button closes it and returns focus to
// the ? button. A click outside closes it and leaves focus where the click
// put it. The ? button itself toggles it, so it is not an outside click.

export function HelpDrawer({ help, open, onClose, opener }: {
  help: PageHelp | undefined
  open: boolean
  onClose: (restoreFocus: boolean) => void
  opener: RefObject<HTMLButtonElement | null>
}) {
  const panel = useRef<HTMLDivElement>(null)
  const closeButton = useRef<HTMLButtonElement>(null)
  const titleId = useId()

  useEffect(() => {
    if (!open) return
    closeButton.current?.focus()
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose(true)
    }
    const onDown = (e: PointerEvent) => {
      const target = e.target as Node
      if (panel.current?.contains(target) || opener.current?.contains(target)) return
      onClose(false)
    }
    document.addEventListener('keydown', onKey)
    document.addEventListener('pointerdown', onDown)
    return () => {
      document.removeEventListener('keydown', onKey)
      document.removeEventListener('pointerdown', onDown)
    }
  }, [open, onClose, opener])

  return (
    <div
      ref={panel}
      role="dialog"
      aria-modal="false"
      aria-labelledby={titleId}
      inert={!open}
      className={`fixed top-0 right-0 bottom-0 z-50 flex w-[26rem] max-w-full flex-col border-l border-line bg-bg shadow-xl transition-[translate,visibility] duration-200 ease-out motion-reduce:transition-none ${open ? 'visible translate-x-0' : 'invisible translate-x-full'}`}
    >
      <div className="flex items-center justify-between border-b border-line px-5 py-3">
        <h2 id={titleId} className="text-sm font-semibold text-ink">How to read this page</h2>
        <button ref={closeButton} type="button" onClick={() => onClose(true)} aria-label="Close"
          className="border border-line px-2 py-0.5 font-mono text-[11px] text-ink-2 hover:text-ink">
          close
        </button>
      </div>
      {help && (
        <div className="flex-1 overflow-y-auto px-5 py-5">
          <div className="font-serif text-xl text-ink">{help.title}</div>
          <p className="mt-1 text-sm leading-6 text-ink-2">{help.purpose}</p>
          {help.guide.map((section) => (
            <section key={section.heading} className="mt-5">
              <h3 className="mb-1 text-xs font-semibold text-ink">{section.heading}</h3>
              <p className="text-sm leading-6 text-ink">{section.text}</p>
            </section>
          ))}
        </div>
      )}
    </div>
  )
}
