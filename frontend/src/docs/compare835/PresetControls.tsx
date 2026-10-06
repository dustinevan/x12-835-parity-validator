import { useState } from 'react'
import { MAX_NAME, formatPreset, parsePreset, type Preset } from './presets'
import { useMessage, type Message, type PresetState } from './usePresets'

// The two places a preset is chosen, saved, and read as text, both drawn from
// the state in usePresets.ts:
//   - PresetBar, at the top of the selector (ListTree.tsx's `top`): a select of
//     the built-in presets, the reader's saved ones, and "Custom" when the
//     hidden set matches none; Delete for a saved one; Save as… for the
//     current hidden set.
//   - PresetText, under the table and the Raw text panel: the current preset as
//     plain text in a box, with Copy, Apply, and Save as….
// The text in the box is page state only. It is not put in the address, not
// kept in browser storage, and not sent anywhere; Save as… keeps the keys it
// holds, under a name, in this browser (presets.ts).

const CONTROL = 'text-link hover:underline disabled:text-ink-2 disabled:no-underline disabled:opacity-50'
const FIELD = 'min-w-0 border border-line bg-bg px-2 py-1 text-xs text-ink'
const NAME_FIELD = `flex-1 ${FIELD} placeholder:text-ink-2/60`
const CUSTOM = '#custom' // the select's value for "Custom", which is never a preset's value

function MessageLine({ message }: { message: Message | null }) {
  if (message === null) return null
  const error = message.kind === 'error'
  return (
    <p role={error ? 'alert' : 'status'} className={`mt-1 text-[11px] leading-snug [overflow-wrap:anywhere] ${error ? 'text-danger' : 'text-ink-2'}`}>
      {message.text}
    </p>
  )
}

// ---------------------------------------------------------------- at the top of the selector

export function PresetBar({ presets }: { presets: PresetState }) {
  const [name, setName] = useState('')
  const { message, show, clear } = useMessage()
  const built = presets.options.filter((o) => !o.saved)
  const mine = presets.options.filter((o) => o.saved)

  const save = () => {
    const result = presets.save(name, presets.hide)
    if ('error' in result) show('error', result.error)
    else if (result.stored) show('note', 'Saved')
    else show('note', 'Saved for this visit only: the browser refused storage.', 6000)
  }

  return (
    <div className="flex flex-col gap-3 border-b border-line px-3 py-2">
      <div>
        <div className="flex items-baseline gap-3">
          <label htmlFor="compare835-preset" className="font-semibold text-ink">Preset</label>
          <button type="button" onClick={() => { clear(); presets.remove() }} disabled={presets.current === null || !presets.current.saved} className={`ml-auto ${CONTROL}`}>
            Delete
          </button>
        </div>
        <select id="compare835-preset" value={presets.current?.value ?? CUSTOM} onChange={(e) => { clear(); presets.choose(e.target.value) }} className={`mt-1 w-full ${FIELD}`}>
          <optgroup label="Built in">
            {built.map((o) => <option key={o.value} value={o.value}>{o.name}</option>)}
          </optgroup>
          {mine.length > 0 && (
            <optgroup label="Saved">
              {mine.map((o) => <option key={o.value} value={o.value}>{o.name}</option>)}
            </optgroup>
          )}
          {presets.current === null && <option value={CUSTOM}>Custom</option>}
        </select>
      </div>
      <form onSubmit={(e) => { e.preventDefault(); save() }}>
        <label htmlFor="compare835-preset-name" className="font-semibold text-ink">Save as…</label>
        <div className="mt-1 flex items-center gap-2">
          <input id="compare835-preset-name" value={name} onChange={(e) => { setName(e.target.value); clear() }} maxLength={MAX_NAME} placeholder="Name"
            autoComplete="off" spellCheck={false} className={NAME_FIELD} />
          <button type="submit" className={CONTROL}>Save</button>
        </div>
        <MessageLine message={message} />
      </form>
    </div>
  )
}

// ---------------------------------------------------------------- under the table

// What the box holds when the reader has typed or pasted: the text, and the
// view of the hidden set it goes with. The box shows it while that is the
// current view; any other change of the hidden set (a click in the tree,
// another preset) shows the current preset's own text again.
interface Draft {
  text: string
  view: string
}

const some = (ids: string[]) => (ids.length <= 3 ? ids.join(', ') : `${ids.slice(0, 3).join(', ')} and ${ids.length - 3} more`)

export function PresetText({ presets }: { presets: PresetState }) {
  const [name, setName] = useState('')
  const [draft, setDraft] = useState<Draft | null>(null)
  const { message, show, clear } = useMessage()

  const text = draft !== null && draft.view === presets.view ? draft.text : formatPreset(presets.current?.name ?? '', presets.hide)

  // The box's text as a preset, or the error under the box and null.
  const read = (): Preset | null => {
    const parsed = parsePreset(text)
    if ('error' in parsed) {
      show('error', parsed.error)
      return null
    }
    return parsed
  }

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(text)
      show('note', 'Copied')
    } catch {
      show('error', 'Not copied: the browser refused clipboard access.') // no clipboard here, or it refused
    }
  }

  // Makes the box's keys the hidden set, and its name (when it has one) the
  // name in the Save as… field. The text stays as it was written.
  const apply = () => {
    const parsed = read()
    if (parsed === null) return
    const missing = parsed.hide.filter((id) => !presets.known.has(id))
    presets.apply(parsed.hide)
    if (parsed.name !== '') setName(parsed.name)
    setDraft({ text, view: presets.viewAfter(parsed.hide) })
    if (missing.length === 0) show('note', 'Applied')
    else show('note', `Applied. Not in the files shown: ${some(missing)}.`, 0)
  }

  // Keeps the box's keys under the name in the field (else the name in the text),
  // and makes them the hidden set.
  const save = () => {
    const parsed = read()
    if (parsed === null) return
    const result = presets.save(name.trim() !== '' ? name : parsed.name, parsed.hide)
    if ('error' in result) {
      show('error', result.error)
      return
    }
    setDraft(null)
    if (result.stored) show('note', 'Saved')
    else show('note', 'Saved for this visit only: the browser refused storage.', 6000)
  }

  return (
    <>
      <section aria-label="Preset as text" className="mt-6 border border-line text-xs">
        <div className="border-b border-line px-3 py-2 font-semibold text-ink">Preset text</div>
        <div className="flex flex-col gap-4 p-3 sm:flex-row sm:items-start">
          <div className="w-full min-w-0 flex-1">
            <textarea
              id="compare835-preset-text"
              aria-label="Preset text"
              value={text}
              onChange={(e) => {
                clear()
                setDraft({ text: e.target.value, view: presets.view })
              }}
              wrap="off"
              spellCheck={false}
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              className="block h-48 w-full resize-none overflow-auto whitespace-pre border border-line bg-bg px-2 py-1.5 font-mono text-[11px] leading-4 text-ink"
            />
            <MessageLine message={message} />
          </div>
          <div className="flex w-full shrink-0 flex-col gap-3 sm:w-56">
            <div className="flex gap-3">
              <button type="button" onClick={copy} className={CONTROL}>Copy</button>
              <button type="button" onClick={apply} className={CONTROL}>Apply</button>
            </div>
            <form onSubmit={(e) => { e.preventDefault(); save() }}>
              <label htmlFor="compare835-preset-text-name" className="font-semibold text-ink">Save as…</label>
              <div className="mt-1 flex items-center gap-2">
                <input id="compare835-preset-text-name" value={name} onChange={(e) => { setName(e.target.value); clear() }} maxLength={MAX_NAME} placeholder="Name"
                  autoComplete="off" spellCheck={false} className={NAME_FIELD} />
                <button type="submit" className={CONTROL}>Save</button>
              </div>
            </form>
          </div>
        </div>
      </section>
      <p className="mt-2 max-w-[52rem] text-[11px] leading-snug text-ink-2">
        The text stays on this page: it is not put in the address, not saved in this browser, and not sent anywhere. Save as… keeps only the element keys, under the name, in this browser.
      </p>
    </>
  )
}
