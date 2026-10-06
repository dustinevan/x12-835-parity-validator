import { useEffect, useRef } from 'react'
import { RawFile } from './RawFile'
import type { EditableRaw } from './useEditableRaw'

const LINK_BUTTON = 'text-ink-2 underline decoration-dotted decoration-line underline-offset-2 hover:text-ink'

// The raw-file panel under the boxes, shared by the three viewers: the
// read view (RawFile, with hover sync) or, after Edit, a textarea holding the
// same text. Every change re-parses and the boxes above follow; Done returns
// to the read view showing the edited text; Reset (shown whenever the text
// differs from the example's) restores the example's own file. When the text
// does not parse, the boxes keep the last text that did and the parser's error
// reads under the panel (see useEditableRaw.ts).
export function RawPanel<B extends string>({ raw, hoveredRaw }: { raw: EditableRaw<B>; hoveredRaw: string | null }) {
  const area = useRef<HTMLTextAreaElement>(null)
  useEffect(() => {
    if (raw.editing) area.current?.focus()
  }, [raw.editing])

  return (
    <div className="mt-6">
      <div className="mb-1.5 flex items-center justify-between gap-3">
        <div className="font-mono text-[11px] uppercase tracking-wide text-ink-2">Raw file</div>
        <div className="flex shrink-0 gap-3 text-xs">
          {raw.edited && (
            <button type="button" onClick={raw.reset} className={LINK_BUTTON}>
              Reset
            </button>
          )}
          <button type="button" onClick={() => raw.setEditing(!raw.editing)} aria-pressed={raw.editing} className={LINK_BUTTON}>
            {raw.editing ? 'Done' : 'Edit'}
          </button>
        </div>
      </div>

      {raw.editing ? (
        <textarea
          ref={area}
          value={raw.text}
          onChange={(e) => raw.setText(e.target.value)}
          aria-label="Raw file text"
          spellCheck={false}
          autoComplete="off"
          autoCapitalize="off"
          wrap="off"
          className="block h-72 w-full resize-none overflow-auto whitespace-pre border border-line bg-surface p-4 font-mono text-[12px] text-ink focus:outline-2 focus:outline-accent"
        />
      ) : (
        <RawFile file={raw.text} hoveredRaw={hoveredRaw} />
      )}

      {raw.error && (
        <p className="mt-1.5 break-words text-[11px] text-danger">
          {raw.error}. The boxes above still show the last text that parsed.
        </p>
      )}
      <p className="mt-1.5 text-[11px] text-ink-2">Pasted or edited text stays in this page: it is not saved, not in the address, and not sent anywhere.</p>
    </div>
  )
}
