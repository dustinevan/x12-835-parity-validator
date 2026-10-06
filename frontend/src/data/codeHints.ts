import type { CodeEntry, CodesResponse } from '../api'
import { CODES_835, ELEMENTS_835, type Element835, type Scope } from './lookupRows835'
import { CAS_GROUP_NAME, CAS_REASON_EL_NAME, CLP02_NAME, CLP_EL_NAME, LQ_NAME, SEGMENT_NAME } from './x12Names835'
import { CODE_ABOUT_835, aboutOf, type About } from './x12About835'

// Explains one coded value the app shows, for the CodeHint popover
// (app/CodeHint.tsx): an element tag as the comparison writes it (CLP04,
// AMT*B6, CAS*CO*253, CAS*CO*253*QTY, NM1*TT, MOA, MIA, LQ, RAW*K3), a CARC
// or RARC code, or a CLP02 claim status value. No React.
//
// A popover has two parts, and the second never repeats the first:
//   - the title, a few words saying what the thing is: the element's public
//     name (x12Names835.ts, through lookupRows835.ts) in sentence case
//     ("Claim payment amount"), or the plain label the API gives a tag where
//     the public documents give the thing no name of five words or fewer
//     ("Sequestration adjustment");
//   - the body, one short paragraph: the sentence of the element's entry
//     (x12About835.ts) that says what it is, without the name it opens with,
//     and for an adjustment or remark code the maintainer's text from
//     /api/codes (backend/pkg/carc.json, rarc.json).
// No citation is shown: the documents behind each sentence are listed in
// x12About835.ts and on the docs site's lookup.

export interface CodeExplanation {
  title: string
  body?: string
  deprecated?: boolean // a CARC or RARC its maintainers no longer list as current
  loading?: boolean // the CARC or RARC text has not arrived yet
}

const MOST_WORDS = 5 // a public name longer than this is not used as a title

const wordsIn = (text: string) => text.trim().split(/\s+/).filter(Boolean).length
const capitalized = (text: string) => text.charAt(0).toUpperCase() + text.slice(1)

// "Claim Payment Amount" reads "Claim payment amount": every word after the
// first in lower case, except a word written in capitals (DRG, HCPCS, a lone
// A) or holding a digit.
export function sentenceCase(name: string): string {
  return name
    .split(' ')
    .map((word, i) => (i === 0 || /^[A-Z]$|[A-Z].*[A-Z]|\d/.test(word) ? word : word.toLowerCase()))
    .join(' ')
}

// The label the API gave a tag, unless it is no more than the tag or a
// group and reason code pair (CO-45), which say nothing about what it is.
function plainLabel(tag: string, label: string | undefined): string | undefined {
  if (!label || label === tag || /^[A-Z]{2}-[A-Z0-9]+$/.test(label)) return undefined
  return label
}

// The title: the public name when it is short enough to be one, else the
// API's label, else whatever name there is.
function titleOf(publicName: string | undefined, label: string | undefined, fallback: string): string {
  if (publicName && wordsIn(publicName) <= MOST_WORDS) return sentenceCase(publicName)
  return label ?? (publicName ? sentenceCase(publicName) : fallback)
}

// What an entry's sentence says beyond the name it opens with: "Claim
// Payment Amount: what Medicare pays ..." gives "What Medicare pays ...".
// A sentence that is the name alone ("Discharge Fraction.") gives nothing.
function beyondName(text: string): string | undefined {
  const lead = /^([^:]{1,100}?): (.+)$/s.exec(text)
  if (lead && !/[.;]/.test(lead[1])) return capitalized(lead[2])
  return wordsIn(text) > 8 ? text : undefined
}

// The body from an entry: what the element is; where the entry names it
// and says no more, what it is for, else how Medicare sends it.
function bodyOf(about: About | undefined): string | undefined {
  if (!about) return undefined
  return beyondName(about.is) ?? about.purpose ?? about.medicare
}

// Which part of the 835 each lookup row's place is in, to read a code
// row's place the same way.
const SCOPE_OF_PLACE = new Map<string, Scope>(ELEMENTS_835.map((r) => [r.where, r.scope]))

// ---------------------------------------------------------------- element tags

// A tag for a whole remark list: the codes of MOA03 to MOA07, of MIA05 and
// MIA20 to MIA23, or of the line's LQ segments, compared as one list.
const REMARK_LISTS: Record<string, { name: string | undefined; entry: string }> = {
  MOA: { name: SEGMENT_NAME.MOA, entry: 'claim:MOA' },
  MIA: { name: SEGMENT_NAME.MIA, entry: 'claim:MIA' },
  LQ: { name: LQ_NAME.HE, entry: 'line:LQ@HE' },
}

// One element (CLP04, SVC03): its rows in the claim, then the line.
function elementRows(ref: string): Element835[] {
  const rows = ELEMENTS_835.filter((r) => r.ref === ref && (r.scope === 'claim' || r.scope === 'line'))
  return [...rows.filter((r) => r.scope === 'claim'), ...rows.filter((r) => r.scope === 'line')]
}

// An adjustment, CAS*CO*253, or its quantity, CAS*CO*253*QTY: the label
// the API gave it, else its group's name; then the group and reason, and the
// reason's text as its maintainers write it.
function adjustment(group: string, reason: string, quantity: boolean, label: string | undefined, codes: CodesResponse | undefined): CodeExplanation {
  const groupName = CAS_GROUP_NAME[group]
  const kind = groupName ? (/adjustment$/i.test(groupName) ? sentenceCase(groupName) : `${sentenceCase(groupName)} adjustment`) : `Group ${group} adjustment`
  const title = label ?? (quantity ? `${kind}, quantity` : kind)
  const what = `${quantity ? 'The units of service adjusted' : 'The amount adjusted'} under group ${group}${groupName ? ` (${groupName})` : ''} for reason ${reason}`
  if (!codes) return { title, body: `${what}.`, loading: true }
  const entry = codes.carc[reason]
  if (!entry) return { title, body: `${what}, which is not on the adjustment reason code list.` }
  return { title, body: `${what}: ${entry.description}`, deprecated: entry.status !== 'current' }
}

// A tag with its segment's qualifier: DTM*050, NM1*TT, REF*6R, AMT*B6,
// QTY*CA. Its name is the qualifier code's; what it is comes from the entry
// for that use of the segment, or, for an AMT or QTY, the segment's entry at
// its level.
function qualified(segment: string, qualifier: string, tag: string, label: string | undefined): CodeExplanation {
  const rows = CODES_835.filter((c) => c.segment === segment && c.code === qualifier && c.element === `${segment}01`)
  const scopes = [...new Set(rows.map((c) => SCOPE_OF_PLACE.get(c.where)).filter((s): s is Scope => s === 'claim' || s === 'line'))]
  const order: Scope[] = scopes.length > 0 ? scopes : ['claim', 'line']
  const title = titleOf(rows.find((c) => c.name)?.name, label, tag)
  for (const scope of order) {
    const body = bodyOf(aboutOf(`${scope}:${segment}@${qualifier}`)) ?? bodyOf(aboutOf(`${scope}:${segment}`))
    if (body) return { title, body }
  }
  return { title }
}

// An element tag as service.TagOf writes it (docs/API.md ## Tags): CLP04,
// AMT*B6, CAS*CO*253, CAS*CO*253*QTY (the adjustment's quantity), NM1*TT,
// a remark list (MOA, MIA, LQ), or RAW*K3 (a segment the typed model has no
// slot for, compared whole). `label` is the plain label the API gave the tag.
export function describeElement(tag: string, label: string | undefined, codes: CodesResponse | undefined): CodeExplanation {
  const given = plainLabel(tag, label)
  if (tag.startsWith('RAW*')) {
    const segment = tag.slice(4)
    return {
      title: given ?? `Unplaced segment ${segment}`,
      body: `A ${segment} segment the typed 835 model has no slot for in the claim (a repeat the standard allows once, or a segment the claim loop does not define), compared whole as one value.`,
    }
  }

  const [base, ...rest] = tag.split('*')
  if (base === 'CAS' && (rest.length === 2 || (rest.length === 3 && rest[2] === 'QTY'))) {
    return adjustment(rest[0], rest[1], rest.length === 3, given, codes)
  }
  if (rest.length > 0) return qualified(base, rest.join('*'), tag, given)

  const rows = elementRows(base)
  if (rows.length > 0) {
    return { title: titleOf(rows.find((r) => r.name)?.name, given, base), body: bodyOf(rows.find((r) => r.about)?.about) }
  }
  const list = REMARK_LISTS[base]
  if (list) return { title: given ?? (list.name ? sentenceCase(list.name) : base), body: bodyOf(aboutOf(list.entry)) }
  const name = SEGMENT_NAME[base]
  return { title: titleOf(name, given, base), body: bodyOf(aboutOf(`claim:${base}`) ?? aboutOf(`line:${base}`)) }
}

// ---------------------------------------------------------------- codes

function listed(title: string, entry: CodeEntry | undefined, codes: CodesResponse | undefined, list: string): CodeExplanation {
  if (!codes) return { title, loading: true }
  if (!entry) return { title, body: `Not on the ${list} list.` }
  return { title, body: entry.description, deprecated: entry.status !== 'current' }
}

export function describeCarc(code: string, codes: CodesResponse | undefined): CodeExplanation {
  return listed(`${sentenceCase(CAS_REASON_EL_NAME)} ${code}`, codes?.carc[code], codes, 'adjustment reason code')
}

export function describeRarc(code: string, codes: CodesResponse | undefined): CodeExplanation {
  return listed(`Remark code ${code}`, codes?.rarc[code], codes, 'remark code')
}

// A CLP02 claim status value; `label` is the plain text format.ts's
// statusLabel already computed for it, and the title. The body is what the
// value means, where a document on file says, and which element it is.
export function describeStatus(code: string, label: string): CodeExplanation {
  const entry = CODE_ABOUT_835[`CLP:CLP02:${code}`]
  const lead = entry && /^([^:]{1,100}?): (.+)$/s.exec(entry.text)
  const publicName = CLP02_NAME[code]
  const meaning = lead && !/[.;]/.test(lead[1])
    ? capitalized(lead[2])
    : publicName && publicName.toLowerCase() !== label.toLowerCase() ? `${publicName}.` : undefined
  const element = `${sentenceCase(CLP_EL_NAME[1])} ${code} (CLP02).`
  return { title: label, body: meaning ? `${meaning} ${element}` : element }
}
