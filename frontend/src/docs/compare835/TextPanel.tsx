import { useEffect, useRef, useState } from 'react'
import { COLS, COL_LABEL, type Col } from './rows'

// The raw text of each column's source, beside the table: an example's file
// as stored while the column shows an example, or whatever the reader has
// pasted or typed. Editing a box makes its text the column's source (the page
// decides what that does to the select and the address; this file only draws
// the boxes). A text that does not read shows its error under its box, and
// the box keeps the text as typed.
//
// Beside the table once the page is 93.5rem wide (the tree, the table, this
// panel, and their gaps and padding side by side); below the table under
// that. The width is the page's own, not the window's, so the panel moves
// below when the lookup panel opens beside the page and leaves too little
// room (viewer/LookupFrame.tsx). Compare835.tsx widens its page cap by the
// width used here (w-80) and its gap.

function TextBox({ col, text, error, onText }: { col: Col; text: string; error?: string; onText: (text: string) => void }) {
  const id = `compare835-text-${col}`
  const [note, setNote] = useState<'copied' | 'failed' | null>(null)
  const timer = useRef<number | undefined>(undefined)
  useEffect(() => () => window.clearTimeout(timer.current), [])

  const copy = async () => {
    let result: 'copied' | 'failed' = 'copied'
    try {
      await navigator.clipboard.writeText(text)
    } catch {
      result = 'failed' // no clipboard here, or the browser refused
    }
    setNote(result)
    window.clearTimeout(timer.current)
    timer.current = window.setTimeout(() => setNote(null), 1500)
  }

  const control = 'text-link hover:underline disabled:text-ink-2 disabled:no-underline disabled:opacity-50'
  return (
    <div>
      <div className="flex items-baseline gap-3">
        <label htmlFor={id} className="font-semibold text-ink">{COL_LABEL[col]}</label>
        <span className="ml-auto flex gap-3" aria-live="polite">
          <button type="button" onClick={copy} disabled={text === ''} className={control}>
            {note === 'copied' ? 'Copied' : note === 'failed' ? 'Not copied' : 'Copy'}
          </button>
          <button type="button" onClick={() => onText('')} disabled={text === ''} className={control}>Clear</button>
        </span>
      </div>
      <textarea
        id={id}
        value={text}
        onChange={(e) => onText(e.target.value)}
        wrap="off"
        spellCheck={false}
        autoComplete="off"
        autoCorrect="off"
        autoCapitalize="off"
        placeholder="Paste an 835 here"
        aria-invalid={error !== undefined}
        aria-describedby={error !== undefined ? `${id}-error` : undefined}
        className="mt-1 block h-40 w-full resize-none overflow-auto whitespace-pre border border-line bg-bg px-2 py-1.5 font-mono text-[11px] leading-4 text-ink placeholder:text-ink-2/60"
      />
      {error !== undefined && (
        <p id={`${id}-error`} role="alert" className="mt-1 text-[11px] leading-snug text-danger [overflow-wrap:anywhere]">{error}</p>
      )}
    </div>
  )
}

export function TextPanel({ texts, errors, onText }: {
  texts: Record<Col, string>
  errors: Partial<Record<Col, string>>
  onText: (col: Col, text: string) => void
}) {
  return (
    <div className="flex w-full max-w-[52rem] flex-col gap-2 @min-[93.5rem]:sticky @min-[93.5rem]:top-6 @min-[93.5rem]:max-h-[calc(100svh-3rem)] @min-[93.5rem]:w-80 @min-[93.5rem]:shrink-0 @min-[93.5rem]:self-start @min-[93.5rem]:overflow-y-auto">
      <aside aria-label="Raw text of each column" className="border border-line text-xs">
        <div className="border-b border-line px-3 py-2 font-semibold text-ink">Raw text</div>
        <div className="flex flex-col gap-4 p-3">
          {COLS.map((col) => (
            <TextBox key={col} col={col} text={texts[col]} error={errors[col]} onText={(text) => onText(col, text)} />
          ))}
        </div>
      </aside>
      <p className="text-[11px] leading-snug text-ink-2">Pasted text stays on this page: it is not put in the address, not saved in this browser, and not sent anywhere.</p>
    </div>
  )
}
