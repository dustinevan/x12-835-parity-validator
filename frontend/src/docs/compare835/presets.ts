// Presets for the 835 compare page's selector (ListTree.tsx).
//
// A preset is a name and a hidden set, spelled the way the address's hide=
// spells it: the node ids of tree.ts ('line', 'claim.CLP', 'claim.CLP.CLP03',
// 'line.REF.REF02'). There are two kinds.
//   - Built in, defined below. Each is worked out from the tree the page has
//     built from the three files, so every id in one is a node that tree
//     holds, and a group the files do not have is never named.
//   - The reader's own, kept in this browser's local storage under
//     docs.835-compare.presets as [{ name, hide }]. Nothing about a preset is
//     in the address beyond hide=.
// A preset also has a plain-text form to copy and paste (formatPreset and
// parsePreset), one key per line, the keys being the ids above:
//
//   preset: Claim level
//   hide:
//     line
//
// Nothing in this file uses React. A preset holds keys and never a value from
// a file. Text pasted into the page's preset box is read here and never kept:
// the page holds it as page state, not in storage and not in the address.
import { SCOPES } from './rows'
import { checkOf, hiddenKeeping, idsOf, leavesOf, type TreeNode } from './tree'

export interface Preset {
  name: string // '' when the text gave none
  hide: string[] // node ids, each once
}

export type ParsedPreset = Preset | { error: string }

// ---------------------------------------------------------------- the keys

export const MAX_NAME = 60

// A key is a node id: a part, then optionally a group and an element, joined
// by dots ('line', 'claim.NM1*QC', 'line.REF.REF02'). Which groups and elements
// exist depends on the files, so only the shape is checked: a known part
// first, then at most two names with no spaces, dots, or commas (a comma
// separates the ids in hide=).
const KEY = new RegExp(`^(?:${SCOPES.join('|')})(?:\\.[^\\s.,]+){0,2}$`)
export const isKey = (text: string): boolean => text.length <= 80 && KEY.test(text)

// A name as it is kept: runs of white space become one space, and the ends are trimmed.
export const normalizeName = (name: string): string => name.replace(/\s+/g, ' ').trim()

// ---------------------------------------------------------------- built in

// A built-in preset names groups of two parts, the claim and the service line,
// by their segment ID ('CLP', 'CAS'; 'NM1' stands for every 'NM1*…' group).
// `keeps` says whether the preset shows the rows of one such group.
// `keepsOthers` says whether it shows the rest of the file: the payment, payer,
// payee, provider summary, and provider adjustment parts.
export interface BuiltIn {
  name: string
  keeps: (part: 'claim' | 'line', segment: string) => boolean
  keepsOthers: boolean
}

export const BUILT_IN: BuiltIn[] = [
  { name: 'Everything', keeps: () => true, keepsOthers: true },
  // Every service-line box hidden.
  { name: 'Claim level', keeps: (part) => part !== 'line', keepsOthers: true },
  // The claim-level boxes hidden except the CLP line that identifies the claim.
  { name: 'Service lines', keeps: (part, segment) => part !== 'claim' || segment === 'CLP', keepsOthers: true },
  // CLP, SVC, CAS and AMT at both levels shown; the rest hidden.
  { name: 'Amounts and adjustments', keeps: (_, segment) => ['CLP', 'CAS', 'SVC', 'AMT'].includes(segment), keepsOthers: false },
  // CLP, MOA, MIA, LQ and DTM at both levels shown; the rest hidden.
  { name: 'Codes and dates', keeps: (_, segment) => ['CLP', 'MOA', 'MIA', 'LQ', 'DTM'].includes(segment), keepsOthers: false },
]

export const BUILT_IN_NAMES: string[] = BUILT_IN.map((b) => b.name)

// The built-in presets for this tree, in the order the select lists them.
export function builtInPresets(tree: TreeNode[]): Preset[] {
  return BUILT_IN.map(({ name, keeps, keepsOthers }) => ({
    name,
    hide: hiddenKeeping(tree, (leaf) => {
      const [part, group = ''] = leaf.id.split('.')
      if (part !== 'claim' && part !== 'line') return keepsOthers
      return keeps(part, group.split('*')[0])
    }),
  }))
}

// ---------------------------------------------------------------- what a hidden set shows

// The leaves a hidden set hides in this tree, as one text. Two hidden sets
// with the same text show the same rows: 'line' and 'line,line.SVC' are one
// view, and an id the tree has no node for hides nothing. The page compares
// views to tell which preset the current hidden set matches.
export function viewOf(tree: TreeNode[], hide: Iterable<string>): string {
  const hidden = new Set(hide)
  return tree
    .flatMap(leavesOf)
    .filter((leaf) => checkOf(leaf, hidden) === 'off')
    .map((leaf) => leaf.id)
    .join(',')
}

// Of `options`, the one that shows what `view` is the view of: the one whose
// value is `chosen` when several show the same rows, else the first; null when
// none does (the page's "Custom").
export function matchingOption<T extends { hide: string[]; value: string }>(tree: TreeNode[], options: T[], view: string, chosen: string): T | null {
  const matches = options.filter((o) => viewOf(tree, o.hide) === view)
  return matches.find((o) => o.value === chosen) ?? matches[0] ?? null
}

// The ids in the order the tree lists them; ids the tree has no node for
// follow, in the order given. Each id once.
export function inTreeOrder(tree: TreeNode[], ids: Iterable<string>): string[] {
  const place = new Map(tree.flatMap(idsOf).map((id, at): [string, number] => [id, at]))
  return [...new Set(ids)]
    .map((id, i) => ({ id, at: place.get(id) ?? place.size + i }))
    .sort((x, y) => x.at - y.at)
    .map((x) => x.id)
}

// ---------------------------------------------------------------- the text form

const HEADER = /^(preset|hide):(.*)$/i

const excerpt = (line: string) => (line.length > 40 ? `${line.slice(0, 40)}…` : line)

// A preset as text: the name line when there is a name, the "hide:" line, then
// one key per line. White space inside the name is made single spaces.
export function formatPreset(name: string, hide: readonly string[]): string {
  const title = normalizeName(name)
  return [...(title === '' ? [] : [`preset: ${title}`]), 'hide:', ...hide.map((id) => `  ${id}`)].join('\n')
}

// Reads a preset's text. The result is the name ('' when the text gives none)
// and the keys, or an error that says which line is wrong. Blank lines, the
// indentation, and CRLF line ends do not matter; a key written twice counts
// once. Nothing is thrown, and nothing is kept.
export function parsePreset(text: string): ParsedPreset {
  if (typeof text !== 'string' || text.trim() === '') return { error: 'The text is empty.' }
  const lines = text.split(/\r\n|\r|\n/)
  let name: string | undefined
  let reading = false // true once the "hide:" line is read
  const hide = new Set<string>()
  const bad: { line: number; text: string }[] = []
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i].trim()
    if (line === '') continue
    if (!reading) {
      const [, word = '', rest = ''] = HEADER.exec(line) ?? []
      if (word.toLowerCase() === 'preset' && name === undefined) {
        name = normalizeName(rest)
        continue
      }
      if (word.toLowerCase() === 'hide') {
        if (rest.trim() !== '') return { error: `Line ${i + 1}: put each key on its own line under "hide:".` }
        reading = true
        continue
      }
      return { error: `Line ${i + 1}: expected "preset: NAME" or "hide:", not "${excerpt(line)}".` }
    }
    if (isKey(line)) hide.add(line)
    else bad.push({ line: i + 1, text: line })
  }
  if (!reading) return { error: 'The text needs a "hide:" line.' }
  if (bad.length > 0) {
    const more = bad.length === 1 ? '' : bad.length === 2 ? ' 1 more line is not a key either.' : ` ${bad.length - 1} more lines are not keys either.`
    return { error: `Line ${bad[0].line} is not a key: "${excerpt(bad[0].text)}". Keys read like claim.CLP or line.REF.REF02.${more}` }
  }
  if (name !== undefined && name.length > MAX_NAME) return { error: `The name is longer than ${MAX_NAME} characters.` }
  return { name: name ?? '', hide: [...hide] }
}

// ---------------------------------------------------------------- the reader's own

export const PRESETS_KEY = 'docs.835-compare.presets'

// A name a built-in choice uses, in any case: the select's own "Custom" too.
const RESERVED = new Set([...BUILT_IN_NAMES, 'Custom'].map((n) => n.toLowerCase()))

// Why a name cannot be saved, or undefined when it can.
export function nameProblem(name: string): string | undefined {
  const title = normalizeName(name)
  if (title === '') return 'Enter a name first.'
  if (title.length > MAX_NAME) return `Use ${MAX_NAME} characters or fewer.`
  if (RESERVED.has(title.toLowerCase())) return `"${title}" is the name of a built-in choice. Use another name.`
  return undefined
}

const sameName = (x: string, y: string) => x.toLowerCase() === y.toLowerCase()

// The list with `preset` added, or in place of the saved preset of the same name.
export function withPreset(list: Preset[], preset: Preset): Preset[] {
  const at = list.findIndex((p) => sameName(p.name, preset.name))
  return at < 0 ? [...list, preset] : list.map((p, i) => (i === at ? preset : p))
}

export function withoutPreset(list: Preset[], name: string): Preset[] {
  return list.filter((p) => !sameName(p.name, name))
}

// Reads what storage held. Anything that is not a list of { name, hide }, an
// entry whose name cannot be saved or repeats an earlier one, and an entry with
// a key that is not a key, is left out; nothing is thrown. `hide` may also be
// the comma-separated text hide= holds.
export function readSavedPresets(text: string | null | undefined): Preset[] {
  if (typeof text !== 'string' || text === '') return []
  let value: unknown
  try {
    value = JSON.parse(text)
  } catch {
    return []
  }
  if (!Array.isArray(value)) return []
  let list: Preset[] = []
  for (const item of value) {
    const entry = (typeof item === 'object' && item !== null ? item : {}) as { name?: unknown; hide?: unknown }
    const hide = typeof entry.hide === 'string' ? entry.hide.split(',').filter(Boolean) : entry.hide
    if (typeof entry.name !== 'string' || !Array.isArray(hide)) continue
    if (!hide.every((id) => typeof id === 'string' && isKey(id))) continue
    const name = normalizeName(entry.name)
    if (nameProblem(name) !== undefined || list.some((p) => sameName(p.name, name))) continue
    list = [...list, { name, hide: [...new Set<string>(hide)] }]
  }
  return list
}

// The browser may refuse storage (a private window, cleared site data), so
// every access is guarded and the page works without it.
export function loadSavedPresets(): Preset[] {
  try {
    return readSavedPresets(localStorage.getItem(PRESETS_KEY))
  } catch {
    return []
  }
}

// Whether the list was kept in this browser.
export function storeSavedPresets(list: Preset[]): boolean {
  try {
    localStorage.setItem(PRESETS_KEY, JSON.stringify(list.map(({ name, hide }) => ({ name, hide }))))
    return true
  } catch {
    return false
  }
}
