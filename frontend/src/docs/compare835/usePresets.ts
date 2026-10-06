import { useEffect, useMemo, useRef, useState } from 'react'
import {
  builtInPresets, inTreeOrder, loadSavedPresets, matchingOption, nameProblem, normalizeName, storeSavedPresets, viewOf, withoutPreset, withPreset, type Preset,
} from './presets'
import { idsOf, type TreeNode } from './tree'

// The state behind the page's presets (presets.ts holds the rules): the
// built-in presets for the tree, the reader's saved ones, and which of them
// the current hidden set matches. The hidden set itself stays where the page
// keeps it (the address and storage, Compare835.tsx); choosing or saving a
// preset hands its ids to `setHidden`, so a preset writes hide= and storage
// the way a click on the tree does. Nothing else about a preset reaches the
// address, and the text of the preset box is never passed through here.

export interface PresetOption extends Preset {
  value: string // the select's value: 'built:Claim level', 'saved:My view'
  saved: boolean
}

// A hidden set a choice has just written, before the address shows it: the
// router applies a change of address a moment after the click, as a transition,
// so for a moment the address still holds the old set. Until it follows, the
// select and the text box describe this set; they would otherwise show the old
// one for that moment and then jump.
interface Pending {
  from: string // the address's view when the choice was made
  ids: string[] // the hidden set the choice wrote
}

export function usePresets(tree: TreeNode[], hidden: Set<string>, setHidden: (ids: Set<string>) => void) {
  const [saved, setSaved] = useState<Preset[]>(loadSavedPresets)
  // The option the reader last chose or saved. It decides between two options
  // that show the same rows, and nothing else.
  const [chosen, setChosen] = useState('')

  const options = useMemo<PresetOption[]>(
    () => [
      ...builtInPresets(tree).map((p): PresetOption => ({ ...p, value: `built:${p.name}`, saved: false })),
      ...saved.map((p): PresetOption => ({ ...p, value: `saved:${p.name}`, saved: true })),
    ],
    [tree, saved],
  )

  // What the hidden set shows, and the option that shows the same (null:
  // "Custom"). The hidden set is the address's, or the one a choice has just
  // written while the address has not yet followed (see Pending).
  const [pending, setPending] = useState<Pending | null>(null)
  const addressView = useMemo(() => viewOf(tree, hidden), [tree, hidden])
  const waiting = pending !== null && addressView === pending.from
  useEffect(() => {
    if (pending !== null && addressView !== pending.from) setPending(null) // the address has followed
  }, [pending, addressView])
  const ids = waiting ? pending.ids : hidden
  const view = waiting ? viewOf(tree, pending.ids) : addressView
  const current = useMemo(() => matchingOption(tree, options, view, chosen), [tree, options, view, chosen])

  const hide = useMemo(() => inTreeOrder(tree, ids), [tree, ids]) // the hidden set, in tree order
  const known = useMemo(() => new Set(tree.flatMap(idsOf)), [tree]) // every key this tree holds

  // Writes `next` as the hidden set, through the page.
  const write = (next: string[]) => {
    setPending({ from: addressView, ids: next })
    setHidden(new Set(next))
  }

  const choose = (value: string) => {
    const option = options.find((o) => o.value === value)
    if (option === undefined) return
    setChosen(value)
    write(option.hide)
  }

  // Makes `next` the hidden set (the text box's Apply).
  const apply = (next: string[]) => {
    setChosen('')
    write(next)
  }

  // Keeps `keys` under `name`, in place of a saved preset of that name, and
  // makes them the hidden set and the preset the chosen option. Says why not
  // when the name cannot be used, or whether the browser kept the list.
  const save = (name: string, keys: string[]): { error: string } | { stored: boolean } => {
    const error = nameProblem(name)
    if (error !== undefined) return { error }
    const title = normalizeName(name)
    const next = withPreset(saved, { name: title, hide: keys })
    setSaved(next)
    setChosen(`saved:${title}`)
    write(keys)
    return { stored: storeSavedPresets(next) }
  }

  // Deletes the selected preset when it is a saved one.
  const remove = () => {
    if (current === null || !current.saved) return
    const next = withoutPreset(saved, current.name)
    setSaved(next)
    setChosen('')
    storeSavedPresets(next)
  }

  return {
    options,
    current,
    view,
    hide,
    known,
    viewAfter: (keys: string[]) => viewOf(tree, keys), // the view `keys` would show here
    choose,
    apply,
    save,
    remove,
  }
}

export type PresetState = ReturnType<typeof usePresets>

// ---------------------------------------------------------------- one line of feedback

export interface Message {
  kind: 'error' | 'note'
  text: string
}

// A line under a control. An error stays until the reader acts again; a note
// goes after `ms` (0 keeps it too).
export function useMessage() {
  const [message, setMessage] = useState<Message | null>(null)
  const timer = useRef<number | undefined>(undefined)
  useEffect(() => () => window.clearTimeout(timer.current), [])

  const clear = () => {
    window.clearTimeout(timer.current)
    setMessage(null)
  }
  const show = (kind: Message['kind'], text: string, ms: number = kind === 'note' ? 1500 : 0) => {
    window.clearTimeout(timer.current)
    setMessage({ kind, text })
    if (ms > 0) timer.current = window.setTimeout(() => setMessage(null), ms)
  }
  return { message, show, clear }
}
