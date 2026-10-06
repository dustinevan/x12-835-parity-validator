// What the docs site's Compared fields page draws: the leave-out list of
// compare set 1, and a state for every element of the file tree. The page shows
// the set the code applies, set "1" in backend/lib/core/service/compareset.go:
// everything inside a claim is compared except the leave-out entries. The
// record of how the set was decided is docs/design/compared-fields.md.
//
// LEFT_OUT is the code's list. scripts/ssr-check-compared-fields.mjs reads
// compareset.go and fails when the two differ, so change them together.

// The five-cent rule: two dollar amounts within this many cents of each other
// are the same amount to the comparison, and every other value compares
// exactly. It is AmountMatchCents in backend/lib/core/service/compare.go;
// scripts/ssr-check-compared-fields.mjs reads that constant and fails when the
// two differ.
export const AMOUNT_MATCH_CENTS = 5

// ---------------------------------------------------------------- the leave-out list

// The two levels the comparison finds a difference at, as compareset.go names
// them.
export type Level = 'claim' | 'line'

// One entry of the leave-out list, in the order compareset.go lists them.
export interface LeaveOut {
  level: Level
  // The tag as the comparison labels it. A tag ending in * covers every
  // qualifier: "QTY*" is every QTY, whatever its QTY01.
  tag: string
  publicName: string
  // One line: why the element is left out.
  reason: string
  // The keys of the file tree the entry covers (see "the file in order").
  keys: string[]
}

export const LEFT_OUT: LeaveOut[] = [
  {
    level: 'claim',
    tag: 'CLP03',
    publicName: 'Total Claim Charge Amount',
    reason: 'A match-key field, equal as an exact decimal within a pair.',
    keys: ['claim.CLP.CLP03'],
  },
  {
    level: 'claim',
    tag: 'NM1*QC',
    publicName: 'Patient Identifier',
    reason: 'A match-key field, equal within a pair apart from the case and spacing the key ignores.',
    keys: ['claim.NM1*QC'],
  },
  {
    level: 'claim',
    tag: 'DTM*232',
    publicName: 'Statement From-To Date',
    reason: "The key's date of service, equal whenever both sides send it.",
    keys: ['claim.DTM.DTM*232'],
  },
  {
    level: 'claim',
    tag: 'CLP07',
    publicName: 'Payer Claim Control #',
    reason: 'Each system assigns its own.',
    keys: ['claim.CLP.CLP07'],
  },
  {
    level: 'claim',
    tag: 'REF*F8',
    publicName: 'Other Claim Related Identification',
    reason: "Carries the original claim's control number.",
    keys: ['claim.REF.REF*F8'],
  },
  {
    level: 'claim',
    tag: 'CLP10',
    publicName: 'Patient Status Code',
    reason: 'Marked not used by Medicare (CMS 835 Flat File p.6).',
    keys: ['claim.CLP.CLP10'],
  },
  {
    level: 'claim',
    tag: 'CLP14',
    publicName: 'Yes/No Condition or Response Code',
    reason: 'Marked not used by Medicare (CMS 835 Flat File p.6).',
    keys: ['claim.CLP.CLP14'],
  },
  {
    level: 'claim',
    tag: 'PER',
    publicName: 'Claim Contact Information',
    reason: "The payer's contact for the claim (CMS 835 Flat File p.11).",
    keys: ['claim.PER'],
  },
  {
    level: 'line',
    tag: 'REF*6R',
    publicName: 'Line Item Control Number',
    reason: 'Pairs the lines, and every line difference already records it.',
    keys: ['line.REF.REF*6R'],
  },
  {
    level: 'line',
    tag: 'QTY*',
    publicName: 'Service Supplemental Quantity',
    reason: 'Marked not used by Medicare (CMS 835 Flat File p.13). Claim-level QTY stays compared.',
    keys: ['line.QTY'],
  },
]

// ---------------------------------------------------------------- the file in order: a state on every element

// The page's first section draws the elements the 28 worked examples send, in
// file order, each with a state:
//   compared           inside a claim, and read by the comparison;
//   compared, watched  compared, and one of the four WATCHED elements;
//   left out           inside a claim, and on the LEFT_OUT list or one of the
//                      NAME_ELEMENTS;
//   outside the claim  not inside a claim or a service line, so never compared;
//   unlisted           inside a claim, and covered by none of the rules below.
//                      It should stay empty: the check fails when it is not.

export type LeafState = 'compared' | 'compared, watched' | 'left out' | 'outside the claim' | 'unlisted'

// The states in the order a tally lists them.
export const LEAF_STATES: LeafState[] = ['compared', 'compared, watched', 'left out', 'outside the claim', 'unlisted']

// A state and, for the hover on its tag, the plain sentence behind it.
export interface LeafMark {
  state: LeafState
  why?: string
}

// A name for a set of tree keys. The label is for readers of this file.
export interface KeyGroup {
  label: string
  keys: string[]
}

// A key names an element of the tree the 835 compare page draws, as that page
// names them (compare835/tree.ts): the part, the group, and the element,
// joined by dots, such as 'claim.CLP.CLP04' or 'line.SVC.SVC02'. NM1 is one
// group per role ('claim.NM1*82'). A segment that carries a qualifier (DTM,
// REF, AMT, QTY, LQ, PER) has one element per qualified tag, written the way
// the comparison labels it ('claim.DTM.DTM*050'); compare835/unionTree.ts says
// why. A key also stands for everything under it ('claim.NM1*82' is every
// element of that role), and a part of a key that ends in * stands for any
// qualifier ('claim.NM1*.NM103' is NM103 of every role).
//
// Qualifier elements (NM101, REF01, AMT01, CAS01 alone, and the like) have no
// keys of their own; they form part of the tag. NM102 and NM108 are read as
// such: they take the state of their role, so 'claim.NM1*82' below covers
// NM101, NM102, NM108, and NM109 of that role.

// The keys of a segment's numbered elements: numbered('MIA', [1, 4], [24, 24])
// is MIA01 to MIA04 and MIA24.
const numbered = (segment: string, ...ranges: [number, number][]): string[] =>
  ranges.flatMap(([from, to]) => Array.from({ length: to - from + 1 }, (_, i) => `${segment}${String(from + i).padStart(2, '0')}`))

// The keys of qualified tags under a group: tagged('claim.AMT', 'AU', 'DY') is
// 'claim.AMT.AMT*AU' and 'claim.AMT.AMT*DY'.
const tagged = (group: string, ...kinds: string[]): string[] => kinds.map((kind) => `${group}.${group.split('.').pop()}*${kind}`)

// The qualifier ranges ZK to ZO and ZL to ZO, letter by letter.
const ZK_TO_ZO = ['ZK', 'ZL', 'ZM', 'ZN', 'ZO']
const ZL_TO_ZO = ['ZL', 'ZM', 'ZN', 'ZO']

// The name elements of an NM1, NM103 to NM107, in every role. The tree tags
// them left out, and LEFT_OUT has no entry for them: the comparison reads an
// NM1's identifier, NM109, under its role, and not its name elements
// (compareset.go says so beside the set).
export const NAME_ELEMENTS: KeyGroup = {
  label: 'NM103 to NM107 of every NM1',
  keys: ['claim.NM1*.NM103', 'claim.NM1*.NM104', 'claim.NM1*.NM105', 'claim.NM1*.NM106', 'claim.NM1*.NM107'],
}

// What the hover on a name element says.
const NAME_ELEMENTS_WHY = "Not read: the comparison reads an NM1's identifier, NM109, and not its name elements."

// The four watched elements: compared and recorded like every other element.
// Whether an element belongs in the compared set is settled from the real
// files, and the set is named on every run, so a change is visible.
export const WATCHED: KeyGroup[] = [
  { label: 'DTM*050', keys: ['claim.DTM.DTM*050'] },
  { label: 'AMT*I', keys: ['claim.AMT.AMT*I'] },
  { label: 'AMT*NL', keys: ['claim.AMT.AMT*NL'] },
  { label: 'NM1*TT', keys: ['claim.NM1*TT'] },
]

const WATCHED_WHY =
  'Compared and recorded like every other element. Whether an element belongs in the compared set is settled from the real files, and the set is named on every run, so a change is visible.'

// The other elements inside a claim the comparison reads. NM103 to NM107 of
// every role are left out before these are read.
export const COMPARED: KeyGroup[] = [
  { label: 'CLP02', keys: ['claim.CLP.CLP02'] },
  { label: 'CLP04', keys: ['claim.CLP.CLP04'] },
  { label: 'CLP05', keys: ['claim.CLP.CLP05'] },
  { label: 'CLP06', keys: ['claim.CLP.CLP06'] },
  { label: 'CLP08, CLP09', keys: ['claim.CLP.CLP08', 'claim.CLP.CLP09'] },
  { label: 'CLP11–CLP13', keys: ['claim.CLP.CLP11', 'claim.CLP.CLP12', 'claim.CLP.CLP13'] },
  { label: 'CAS (claim)', keys: ['claim.CAS.CAS03'] },
  { label: 'MIA01–04, 06–19, 24', keys: numbered('MIA', [1, 4], [6, 19], [24, 24]).map((ref) => `claim.MIA.${ref}`) },
  { label: 'MOA01, 02, 08, 09', keys: ['claim.MOA.MOA01', 'claim.MOA.MOA02', 'claim.MOA.MOA08', 'claim.MOA.MOA09'] },
  // The tree lists the remark codes of one segment as one element, named for its first position.
  { label: 'MIA05, 20–23; MOA03–07', keys: ['claim.MIA.MIA05', 'claim.MOA.MOA03'] },
  { label: 'AMT*AU, DY, ZK–ZO', keys: tagged('claim.AMT', 'AU', 'DY', ...ZK_TO_ZO) },
  { label: 'QTY*CA, CD, LA, OU, ZK–ZO', keys: tagged('claim.QTY', 'CA', 'CD', 'LA', 'OU', ...ZK_TO_ZO) },
  { label: 'DTM*036', keys: ['claim.DTM.DTM*036'] },
  { label: 'NM1*74 NM109; NM1*PR NM109', keys: ['claim.NM1*74', 'claim.NM1*PR'] },
  { label: 'SVC01, SVC06', keys: ['line.SVC.SVC01', 'line.SVC.SVC06'] },
  { label: 'SVC05, SVC07', keys: ['line.SVC.SVC05', 'line.SVC.SVC07'] },
  { label: 'SVC02', keys: ['line.SVC.SVC02'] },
  { label: 'SVC03', keys: ['line.SVC.SVC03'] },
  { label: 'CAS (line)', keys: ['line.CAS.CAS03'] },
  { label: 'LQ*HE', keys: tagged('line.LQ', 'HE') },
  { label: 'REF*0K', keys: tagged('line.REF', '0K') },
  { label: 'REF*LU, 1S, APC, RB', keys: tagged('line.REF', 'LU', '1S', 'APC', 'RB') },
  { label: 'AMT*B6, KH, 2K, ZL–ZO', keys: tagged('line.AMT', 'B6', 'KH', '2K', ...ZL_TO_ZO) },
  // Echoed claim data: values carried over from the submitted claim.
  { label: 'CLP01', keys: ['claim.CLP.CLP01'] },
  { label: 'DTM*233; DTM*472, 150, 151 (line)', keys: ['claim.DTM.DTM*233', ...tagged('line.DTM', '472', '150', '151')] },
  { label: 'REF*EA, 28, 6P', keys: tagged('claim.REF', 'EA', '28', '6P') },
  { label: 'NM1*82 NM109; REF*HPI, SY, TJ, 1C, 1G (line)', keys: ['claim.NM1*82', ...tagged('line.REF', 'HPI', 'SY', 'TJ', '1C', '1G')] },
  { label: 'SVC04', keys: ['line.SVC.SVC04'] },
  { label: 'AMT*F5', keys: ['claim.AMT.AMT*F5'] },
  // CAS adjustment quantities, compared.
  { label: 'CAS04, 07, …, 19', keys: ['claim.CAS.CAS04', 'line.CAS.CAS04'] },
]

// The segments the comparison places at each level. Any other segment in a
// claim or a service line is compared whole, as one value under its own label.
const PLACED_SEGMENTS: Record<'claim' | 'line', string[]> = {
  claim: ['CLP', 'NM1', 'MOA', 'MIA', 'DTM', 'REF', 'PER', 'AMT', 'QTY', 'CAS'],
  line: ['SVC', 'DTM', 'CAS', 'REF', 'AMT', 'QTY', 'LQ'],
}

// Whether a pattern covers a key: the pattern is the key, or its start, part
// by part; a part ending in * covers any qualifier.
function covers(pattern: string, key: string): boolean {
  const want = pattern.split('.')
  const have = key.split('.')
  return want.length <= have.length && want.every((part, i) => (part.endsWith('*') ? have[i].startsWith(part) : have[i] === part))
}

// The state of an element, from its tree key. In this order: an entry of the
// leave-out list; a name element of an NM1; anything else outside a claim; a
// watched element; any other compared element; a segment the comparison does
// not place; and what is left is inside a claim and in none of these, which is
// unlisted. Left out is only ever inside a claim: the entries' keys are at
// their level, and the name elements are claim keys.
export function stateOf(key: string): LeafMark {
  const left = LEFT_OUT.find((entry) => entry.keys.some((k) => covers(k, key)))
  if (left) return { state: 'left out', why: left.reason }
  if (NAME_ELEMENTS.keys.some((k) => covers(k, key))) return { state: 'left out', why: NAME_ELEMENTS_WHY }

  const [part = '', group = ''] = key.split('.')
  if (part !== 'claim' && part !== 'line') return { state: 'outside the claim', why: 'Not inside a claim.' }

  if (WATCHED.some((entry) => entry.keys.some((k) => covers(k, key)))) return { state: 'compared, watched', why: WATCHED_WHY }
  if (COMPARED.some((entry) => entry.keys.some((k) => covers(k, key)))) return { state: 'compared' }
  if (!PLACED_SEGMENTS[part].includes(group.split('*')[0])) return { state: 'compared' }
  return { state: 'unlisted', why: 'Neither the compared rule nor the left-out rule covers this element.' }
}

// "14 compared (2 watched) · 3 left out": how many of a group's elements are in
// each state, in the order of LEAF_STATES, leaving out a state none is in. A
// watched element is also counted as compared.
export function tallyText(states: LeafState[]): string {
  const count = (state: LeafState) => states.filter((s) => s === state).length
  const watched = count('compared, watched')
  const compared = count('compared') + watched
  return [
    compared > 0 ? `${compared} compared${watched > 0 ? ` (${watched} watched)` : ''}` : '',
    count('left out') > 0 ? `${count('left out')} left out` : '',
    count('outside the claim') > 0 ? `${count('outside the claim')} outside the claim` : '',
    count('unlisted') > 0 ? `${count('unlisted')} unlisted` : '',
  ]
    .filter((part) => part !== '')
    .join(' · ')
}
