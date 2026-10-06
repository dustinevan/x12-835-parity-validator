// Every element and coded value of the 835, as the CMS 835 flat file lays
// it out for Medicare, with the name a public document gives it and the page
// that names it. Built only from x12Names835.ts (names and citations, each
// cited to its page), x12Echo835.ts (elements a public document says carry the
// value submitted on the 837), and x12About835.ts (what each element is,
// what it is for, its values, and how Medicare fills it, each sentence
// cited). Read by the 835 viewer's lookup panel
// (frontend/src/docs/viewer/Lookup.tsx, opened from
// docs/pages/Viewer835.tsx), the 835 Reference, and the validator app's
// hover explanations (codeHints.ts). Nothing here is invented: an element the
// public documents do not name has no `name`, and the pages show its
// reference instead.
//
// Each element row also says what the validator does with the element,
// read from the compared set the Compared fields page draws
// (comparedFields.ts stateOf), so the lookup and that page cannot disagree.
//
// Citation form: "CMS 835 Flat File, PDF p.N" counts PDF viewer pages of the
// 15-page docs/reference/file-formats/cms-835-flatfile.pdf, from 1; "CEDI 835
// CG, PDF p.N (printed p.M)" is the 33-page CEDI companion guide (N the PDF
// page, M its printed footer); x12.org addresses are X12's public code lists.
import {
  AMT01_EL_NAME, BPR03_NAME, BPR_EL_NAME, CARC, CAS01_EL_NAME, CAS_AMOUNT_EL_NAME, CAS_GROUP_NAME, CAS_QTY_EL_NAME,
  CAS_REASON_EL_NAME, CLAIM_AMT02_EL_NAME, CLAIM_AMT_NAME, CLAIM_CAS_NAME, CLAIM_PER_NAME, CLAIM_QTY_NAME, CLAIM_REF_NAME,
  CLP_EL_NAME, DTM01_EL_NAME, DTM02_EL_NAME_BY_QUALIFIER, DTM_NAME, GE_EL_NAME, GS_EL_NAME, HEADER_REF_NAME,
  ID_QUALIFIER_NAME, IEA_EL_NAME, ISA_EL_NAME, LQ_EL_NAME, LQ_NAME, LX01_EL_NAME, MIA_EL_NAME, MOA_EL_NAME, N1_NAME,
  N1_PE_EL_NAME, N1_PR_EL_NAME, N3_PE_EL_NAME, N3_PR_EL_NAME, N4_PE_EL_NAME, N4_PR_EL_NAME, NM101_EL_NAME, NM102_EL_NAME,
  NM1_74_EL_NAME, NM1_82_EL_NAME, NM1_NAME, NM1_PR_EL_NAME, NM1_QC_EL_NAME, NM1_TT_EL_NAME, PAYEE_LOOP_NAME,
  PAYER_LOOP_NAME, PAYER_PER_NAME, PER01_EL_NAME, PER_CLAIM_EL_NAME, PER_PAYER_EL_NAME, PLB01_EL_NAME, PLB02_EL_NAME,
  PLB_AMOUNT_EL_NAME, PLB_COMPOSITE_EL_NAME, PLB_REASON_NAME, QTY01_EL_NAME, RARC, RDM_EL_NAME, REF_CLAIM_EL_NAME,
  REF_HEADER_EL_NAME, REF_LINE_EL_NAME, REF_PAYEE_EL_NAME, REF_PAYER_EL_NAME, SEGMENT_NAME, SERVICE_AMT02_EL_NAME,
  SERVICE_AMT_NAME, SERVICE_CAS_NAME, SERVICE_REF_NAME, SE_EL_NAME, ST_EL_NAME, SVC_EL_NAME, SVC_QUALIFIER_NAME,
  TRANSACTION_TYPE_NAME, TRN_EL_NAME, TRN_NAME, TS2_EL_NAME, TS3_EL_NAME,
  CLP02_NAME, CLP06_NAME, LX01_NAME, NM102_NAME, ID_QUALIFIER_NAME_MORE, CLAIM_AMT_QUALIFIER_NAME, SERVICE_AMT_QUALIFIER_NAME, QTY_QUALIFIER_NAME, LINE_DTM_NAME, REF_QUALIFIER_NAME,
} from './x12Names835'
import { X12_ECHO_837 } from './x12Echo835'
import { CODE_ABOUT_835, X12_CARC, X12_GROUP_CODES, X12_PLB_REASONS, X12_RARC, aboutOf, type About } from './x12About835'
import { stateOf, type LeafState } from './comparedFields'

// The public documents behind every name here (defined beside the entries
// that cite them, in x12About835.ts).
export { CMS_FIELDS, CMS_OVERVIEW, CMS_MANUAL, CMS_MANUAL_2003, CMS_TI_GUIDE, CMS_MLN_FAQ, CMS_MREP_GUIDE, CEDI_835_GUIDE, X12_EXAMPLES, X12_GROUP_CODES, X12_CARC, X12_RARC, X12_PLB_REASONS } from './x12About835'

// The parts of the transaction a row sits in, in file order.
export type Scope = 'envelope' | 'header' | 'payer' | 'payee' | 'summary' | 'claim' | 'line' | 'provider' | 'trailer'

// What the validator does with an element: the state the Compared fields
// page gives it, or 'part of the tag' for a qualifier element, whose value
// names the element beside it (DTM*050, CAS*CO*45) and is never compared
// on its own. `why` is the sentence behind the state; `detail` lists the
// states of each qualifier when one row stands for several kinds (AMT02 at
// the claim level is ten kinds of amount).
export interface ValidatorNote {
  state: LeafState | 'part of the tag'
  why: string
  detail?: string
}

export interface Element835 {
  ref: string // the element reference: ISA08, CLP04, PLB03
  name?: string // its public name; absent when no public document names it
  segment: string
  segmentName?: string // the segment's public name at this use
  tag?: string // the segment with its qualifier at this use: 'DTM*050', 'NM1*QC' (the form the comparison's labels and the compare page use)
  where: string // the part of the 835 it sits in
  scope: Scope
  use: string // the key of this use in x12About835.ts: 'claim:DTM02@050'
  source: string
  echo?: string // the public text saying this element carries the value submitted on the 837
  about?: About // what it is, what it is for, its values, how Medicare fills it
  segmentAbout?: About // the same for the segment at this use
  validator: ValidatorNote
}

export interface Code835 {
  segment: string
  element: string // the element the code is a value of, e.g. "REF01"
  code: string
  name?: string
  full?: string // the maintainer's complete text (CARC, RARC), when longer than `name`
  where: string
  source: string
  // The one element row this code belongs under (that row's own
  // `segmentName`), so a bare element ref repeated at several boxes in the
  // same segment (REF01 at header/claim/line, NM101 at every NM1 role)
  // doesn't show every box's codes under every other box. Left out for a
  // code genuinely shared by more than one box — CARC/RARC and the CAS
  // group codes, which apply at both the claim and the line level.
  segmentName?: string
  // What the code means beyond its name, where a document says more, with
  // the citations.
  about?: { text: string; sources: string[] }
}

const ENVELOPE = 'Envelope'
const HEADER = 'Header, before the loops'
const PAYER = 'Loop 1000A · payer'
const PAYEE = 'Loop 1000B · payee'
const SUMMARY = 'Loop 2000 · header number'
const CLAIM = 'Loop 2100 · claim'
const LINE = 'Loop 2110 · service line'
const PROVIDER = 'Provider adjustment, after the claims'
const TRAILER = 'Trailer'

const SCOPE_OF: Record<string, Scope> = {
  [ENVELOPE]: 'envelope',
  [HEADER]: 'header',
  [PAYER]: 'payer',
  [PAYEE]: 'payee',
  [SUMMARY]: 'summary',
  [CLAIM]: 'claim',
  [LINE]: 'line',
  [PROVIDER]: 'provider',
  [TRAILER]: 'trailer',
}

// The part of the 835 compare page's tree (compare835/tree.ts) a scope maps
// to, which is how comparedFields.ts keys an element. The envelope and the
// trailer are not in that tree.
const PART_OF: Partial<Record<Scope, string>> = {
  header: 'payment',
  payer: 'payer',
  payee: 'payee',
  summary: 'summary',
  claim: 'claim',
  line: 'line',
  provider: 'adjustment',
}

const cms = (pages: string) => `CMS 835 Flat File, PDF p.${pages}`
const cedi = (pdf: string, printed: string) => `CEDI 835 CG, PDF p.${pdf} (printed p.${printed})`
const pad = (n: number) => String(n).padStart(2, '0')

interface Context {
  segmentName?: string
  // The segment's own qualifier value at this use (REF01, DTM01, NM101,
  // PER01, LQ01), which decides the element names, the 837 echo key, and
  // the entry key.
  qualifier?: string
}

// ---------------------------------------------------------------- what the validator does with an element

const OUTSIDE = 'Not inside a claim: the comparison reads claims and their service lines only.'
const ENVELOPE_NOTE = 'The envelope is read for its delimiters and its transactions and compared nowhere; its control numbers differ between any two files by construction.'
const TAG_NOTE = 'Not compared as a value: it names the kind of the element beside it, and the comparison labels that element with it (DTM*050, CAS*CO*45, NM1*TT).'

// The elements whose value is a qualifier.
const QUALIFIER_REFS = new Set(['NM101', 'REF01', 'DTM01', 'AMT01', 'QTY01', 'LQ01', 'PER01', 'CAS01'])
const CAS_REASONS = new Set(['CAS02', 'CAS05', 'CAS08', 'CAS11', 'CAS14', 'CAS17'])
const CAS_QUANTITIES = new Set(['CAS04', 'CAS07', 'CAS10', 'CAS13', 'CAS16', 'CAS19'])

// The qualifiers Medicare sends in an AMT or QTY at each level (the CMS 835
// flat file's Values column), so one row's state can list every kind.
const KINDS: Record<string, string[]> = {
  'claim:AMT': ['AU', 'DY', 'F5', 'I', 'NL', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'],
  'claim:QTY': ['CA', 'CD', 'LA', 'OU', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'],
  'line:AMT': ['B6', 'KH', '2K', 'ZL', 'ZM', 'ZN', 'ZO'],
  'line:QTY': ['ZK'],
}

const noteOf = (key: string): ValidatorNote => {
  const { state, why } = stateOf(key)
  return { state, why: why ?? 'Read by the comparison.' }
}

function validatorNote(scope: Scope, segment: string, ref: string, qualifier: string | undefined): ValidatorNote {
  if (scope === 'envelope' || scope === 'trailer') return { state: 'outside the claim', why: ENVELOPE_NOTE }
  const part = PART_OF[scope]
  if (!part || (scope !== 'claim' && scope !== 'line')) return { state: 'outside the claim', why: OUTSIDE }
  if (QUALIFIER_REFS.has(ref) || CAS_REASONS.has(ref)) return { state: 'part of the tag', why: TAG_NOTE }

  if (segment === 'NM1' && qualifier) {
    const role = `${part}.NM1*${qualifier}`
    const position = Number(ref.slice(3))
    if (position >= 3 && position <= 7) return noteOf(`${role}.${ref}`)
    return noteOf(role)
  }
  if (segment === 'CAS') return noteOf(`${part}.CAS.${CAS_QUANTITIES.has(ref) ? 'CAS04' : 'CAS03'}`)
  if ((segment === 'AMT' || segment === 'QTY') && !qualifier) {
    const kinds = KINDS[`${scope}:${segment}`] ?? []
    const byState = new Map<LeafState, string[]>()
    let why = ''
    for (const kind of kinds) {
      const mark = stateOf(`${part}.${segment}.${segment}*${kind}`)
      byState.set(mark.state, [...(byState.get(mark.state) ?? []), `${segment}*${kind}`])
      why ||= mark.why ?? ''
    }
    const [first] = byState.keys()
    const detail = [...byState.entries()].map(([state, tags]) => `${state}: ${tags.join(', ')}`).join(' · ')
    return { state: first ?? 'compared', why: why || 'Read by the comparison, each kind under its own label.', detail }
  }
  if (qualifier && ['DTM', 'REF', 'LQ', 'PER'].includes(segment)) return noteOf(`${part}.${segment}.${segment}*${qualifier}`)
  if (segment === 'MIA' && ['MIA20', 'MIA21', 'MIA22', 'MIA23'].includes(ref)) return noteOf(`${part}.MIA.MIA05`)
  if (segment === 'MOA' && ['MOA04', 'MOA05', 'MOA06', 'MOA07'].includes(ref)) return noteOf(`${part}.MOA.MOA03`)
  return noteOf(`${part}.${segment}.${ref}`)
}

// ---------------------------------------------------------------- the rows

// One row per element, numbered from 01. `page` gives the flat-file page
// for element n (1-based). A `name` of undefined is an element the public
// documents list without naming.
function rows(segment: string, where: string, names: (string | undefined)[], page: (n: number) => string, context: Context = {}): Element835[] {
  const scope = SCOPE_OF[where]
  const suffix = context.qualifier ? `@${context.qualifier}` : ''
  const segmentAbout = aboutOf(`${scope}:${segment}${suffix}`)
  return names.map((name, i) => {
    const n = i + 1
    const ref = `${segment}${pad(n)}`
    const echoKey = context.qualifier ? `${ref}@${context.qualifier}` : ref
    const echo = X12_ECHO_837[echoKey]?.source ?? (context.qualifier ? X12_ECHO_837[ref]?.source : undefined)
    const use = `${scope}:${ref}${suffix}`
    const row: Element835 = { ref, segment, where, scope, use, source: cms(page(n)), validator: validatorNote(scope, segment, ref, context.qualifier) }
    if (name) row.name = name
    if (context.segmentName) row.segmentName = context.segmentName
    if (context.qualifier) row.tag = `${segment}*${context.qualifier}`
    if (echo) row.echo = echo
    const about = aboutOf(use)
    if (about) row.about = about
    if (segmentAbout) row.segmentAbout = segmentAbout
    return row
  })
}

const dtm = (where: string, qualifier: string, page: string) =>
  rows('DTM', where, [DTM01_EL_NAME, DTM02_EL_NAME_BY_QUALIFIER[qualifier]], () => page, { segmentName: DTM_NAME[qualifier], qualifier })

// CAS01, then six reason/amount/quantity sets (CAS02–CAS19).
const casNames = [CAS01_EL_NAME, ...Array.from({ length: 6 }, () => [CAS_REASON_EL_NAME, CAS_AMOUNT_EL_NAME, CAS_QTY_EL_NAME]).flat()]

// PLB01, PLB02, then six identifier/amount pairs (PLB03–PLB14).
const plbNames = [PLB01_EL_NAME, PLB02_EL_NAME, ...Array.from({ length: 6 }, () => [PLB_COMPOSITE_EL_NAME, PLB_AMOUNT_EL_NAME]).flat()]

const nm1 = (role: string, names: string[], page: (n: number) => string) =>
  rows('NM1', CLAIM, [NM101_EL_NAME, NM102_EL_NAME, ...names], page, { segmentName: NM1_NAME[role], qualifier: role })

export const ELEMENTS_835: Element835[] = [
  ...rows('ISA', ENVELOPE, ISA_EL_NAME, () => '1', { segmentName: SEGMENT_NAME.ISA }),
  ...rows('GS', ENVELOPE, GS_EL_NAME, (n) => (n <= 3 ? '1' : '2'), { segmentName: SEGMENT_NAME.GS }),
  ...rows('ST', ENVELOPE, ST_EL_NAME, () => '1', { segmentName: SEGMENT_NAME.ST }),

  ...rows('BPR', HEADER, BPR_EL_NAME, (n) => (n <= 11 ? '1' : '2'), { segmentName: SEGMENT_NAME.BPR }),
  ...rows('TRN', HEADER, TRN_EL_NAME, () => '2', { segmentName: TRN_NAME['1'], qualifier: '1' }),
  ...rows('CUR', HEADER, [undefined, undefined], () => '2', { segmentName: SEGMENT_NAME.CUR }),
  ...rows('REF', HEADER, REF_HEADER_EL_NAME.EV, () => '2', { segmentName: HEADER_REF_NAME.EV, qualifier: 'EV' }),
  ...rows('REF', HEADER, REF_HEADER_EL_NAME.F2, () => '2', { segmentName: HEADER_REF_NAME.F2, qualifier: 'F2' }),
  ...dtm(HEADER, '405', '2'),

  ...rows('N1', PAYER, N1_PR_EL_NAME, () => '2', { segmentName: N1_NAME.PR, qualifier: 'PR' }),
  ...rows('N3', PAYER, N3_PR_EL_NAME, () => '2', { segmentName: PAYER_LOOP_NAME.N3 }),
  ...rows('N4', PAYER, N4_PR_EL_NAME, () => '2', { segmentName: PAYER_LOOP_NAME.N4 }),
  ...rows('REF', PAYER, REF_PAYER_EL_NAME, () => '3', { segmentName: PAYER_LOOP_NAME.REF, qualifier: '2U' }),
  ...rows('PER', PAYER, [PER01_EL_NAME, ...PER_PAYER_EL_NAME.CX], () => '3', { segmentName: PAYER_PER_NAME.CX, qualifier: 'CX' }),
  ...rows('PER', PAYER, [PER01_EL_NAME, ...PER_PAYER_EL_NAME.BL], () => '3', { segmentName: PAYER_PER_NAME.BL, qualifier: 'BL' }),
  ...rows('PER', PAYER, [PER01_EL_NAME, ...PER_PAYER_EL_NAME['1C']], () => '3', { segmentName: PAYER_PER_NAME['1C'], qualifier: '1C' }),

  ...rows('N1', PAYEE, N1_PE_EL_NAME, () => '3', { segmentName: N1_NAME.PE, qualifier: 'PE' }),
  ...rows('N3', PAYEE, N3_PE_EL_NAME, () => '3', { segmentName: PAYEE_LOOP_NAME.N3 }),
  ...rows('N4', PAYEE, N4_PE_EL_NAME, (n) => (n === 1 ? '3' : '4'), { segmentName: PAYEE_LOOP_NAME.N4 }),
  ...rows('REF', PAYEE, REF_PAYEE_EL_NAME, () => '4', { segmentName: PAYEE_LOOP_NAME.REF, qualifier: 'TJ' }),
  ...rows('RDM', PAYEE, RDM_EL_NAME, () => '4', { segmentName: SEGMENT_NAME.RDM }),

  ...rows('LX', SUMMARY, [LX01_EL_NAME], () => '4', { segmentName: SEGMENT_NAME.LX }),
  ...rows('TS3', SUMMARY, TS3_EL_NAME, (n) => (n <= 20 ? '4' : '5'), { segmentName: SEGMENT_NAME.TS3 }),
  ...rows('TS2', SUMMARY, TS2_EL_NAME, () => '5', { segmentName: SEGMENT_NAME.TS2 }),

  ...rows('CLP', CLAIM, CLP_EL_NAME, (n) => (n <= 3 ? '5' : '6'), { segmentName: SEGMENT_NAME.CLP }),
  ...rows('CAS', CLAIM, casNames, () => '6', { segmentName: CLAIM_CAS_NAME }),
  ...nm1('QC', NM1_QC_EL_NAME, (n) => (n <= 2 ? '6-8' : '6-7')),
  ...nm1('74', NM1_74_EL_NAME, (n) => (n <= 2 ? '6-8' : '7')),
  ...nm1('82', NM1_82_EL_NAME, (n) => (n <= 2 ? '6-8' : '7')),
  ...nm1('TT', NM1_TT_EL_NAME, (n) => (n <= 2 ? '6-8' : '7')),
  ...nm1('PR', NM1_PR_EL_NAME, (n) => (n <= 2 ? '6-8' : '7-8')),
  ...rows('MIA', CLAIM, MIA_EL_NAME, (n) => (n <= 2 ? '8' : '9'), { segmentName: SEGMENT_NAME.MIA }),
  ...rows('MOA', CLAIM, MOA_EL_NAME, () => '10', { segmentName: SEGMENT_NAME.MOA }),
  ...rows('REF', CLAIM, REF_CLAIM_EL_NAME, () => '10', { segmentName: CLAIM_REF_NAME.EA, qualifier: 'EA' }),
  ...rows('REF', CLAIM, REF_CLAIM_EL_NAME, () => '10', { segmentName: CLAIM_REF_NAME.F8, qualifier: 'F8' }),
  ...dtm(CLAIM, '232', '10'),
  ...dtm(CLAIM, '233', '10'),
  ...dtm(CLAIM, '036', '10'),
  ...dtm(CLAIM, '050', '11'),
  ...rows('PER', CLAIM, [PER01_EL_NAME, ...PER_CLAIM_EL_NAME], () => '11', { segmentName: CLAIM_PER_NAME.CX, qualifier: 'CX' }),
  ...rows('AMT', CLAIM, [AMT01_EL_NAME, CLAIM_AMT02_EL_NAME], () => '11', { segmentName: CLAIM_AMT_NAME }),
  ...rows('QTY', CLAIM, [QTY01_EL_NAME, undefined], () => '11', { segmentName: CLAIM_QTY_NAME }),

  ...rows('SVC', LINE, SVC_EL_NAME, (n) => (n <= 4 ? '11' : '12'), { segmentName: SEGMENT_NAME.SVC }),
  ...dtm(LINE, '472', '12'),
  ...rows('CAS', LINE, casNames, (n) => (n === 1 ? '12' : '12-13'), { segmentName: SERVICE_CAS_NAME }),
  ...rows('REF', LINE, REF_LINE_EL_NAME['6R'], () => '13', { segmentName: SERVICE_REF_NAME['6R'], qualifier: '6R' }),
  ...rows('REF', LINE, REF_LINE_EL_NAME.LU, () => '12', { segmentName: SERVICE_REF_NAME.LU, qualifier: 'LU' }),
  ...rows('REF', LINE, REF_LINE_EL_NAME.HPI, () => '13', { segmentName: SERVICE_REF_NAME.HPI, qualifier: 'HPI' }),
  ...rows('REF', LINE, REF_LINE_EL_NAME['0K'], () => '13', { segmentName: SERVICE_REF_NAME['0K'], qualifier: '0K' }),
  ...rows('AMT', LINE, [AMT01_EL_NAME, SERVICE_AMT02_EL_NAME], () => '13', { segmentName: SERVICE_AMT_NAME }),
  ...rows('LQ', LINE, LQ_EL_NAME, () => '13', { segmentName: LQ_NAME.HE, qualifier: 'HE' }),

  ...rows('PLB', PROVIDER, plbNames, (n) => (n <= 2 ? '13' : '13-14'), { segmentName: SEGMENT_NAME.PLB }),

  ...rows('SE', TRAILER, SE_EL_NAME, () => '14', { segmentName: SEGMENT_NAME.SE }),
  ...rows('GE', TRAILER, GE_EL_NAME, () => '14', { segmentName: SEGMENT_NAME.GE }),
  ...rows('IEA', TRAILER, IEA_EL_NAME, (n) => (n === 1 ? '14' : '15'), { segmentName: SEGMENT_NAME.IEA }),
]

const code = (segment: string, element: string, value: string, name: string | undefined, where: string, source: string, segmentName?: string): Code835 => {
  const row: Code835 = { segment, element, code: value, where, source }
  if (name) row.name = name
  if (segmentName) row.segmentName = segmentName
  const about = CODE_ABOUT_835[`${segment}:${element}:${value}`]
  if (about) row.about = about
  return row
}

const carcSource = `X12 Claim Adjustment Reason Codes, ${X12_CARC} (text embedded verbatim from backend/pkg/carc.json)`
const rarcSource = `X12 Remittance Advice Remark Codes, ${X12_RARC} (text embedded verbatim from backend/pkg/rarc.json)`

// A coded value a document lists without naming: the row says so, and the
// entry (x12About835.ts) carries whatever the documents do say about it.
const UNNAMED = undefined

// Coded values, with the name a public document gives each one. The CARC
// and RARC rows are only the codes the 835 examples on /835-viewer use; the
// full lists are X12's, and the validator app serves them at /api/codes.
export const CODES_835: Code835[] = [
  code('GS', 'GS01', 'HP', TRANSACTION_TYPE_NAME.HP, ENVELOPE, cedi('7', '1'), SEGMENT_NAME.GS),
  code('ST', 'ST01', '835', TRANSACTION_TYPE_NAME['835'], ENVELOPE, cedi('7', '1'), SEGMENT_NAME.ST),
  code('BPR', 'BPR03', 'C', BPR03_NAME.C, HEADER, `${cms('1')}; ${cedi('24', '18')}`, SEGMENT_NAME.BPR),
  code('BPR', 'BPR04', 'ACH', UNNAMED, HEADER, `${cms('1')}; ${cedi('24', '18')}`, SEGMENT_NAME.BPR),
  code('BPR', 'BPR04', 'CHK', UNNAMED, HEADER, `${cms('1')}; ${cedi('24', '18')}`, SEGMENT_NAME.BPR),
  code('BPR', 'BPR04', 'NON', UNNAMED, HEADER, `${cms('1')}; ${cedi('24', '18')}`, SEGMENT_NAME.BPR),
  code('TRN', 'TRN01', '1', TRN_NAME['1'], HEADER, cms('2'), TRN_NAME['1']),
  code('REF', 'REF01', 'EV', HEADER_REF_NAME.EV, HEADER, cms('2'), HEADER_REF_NAME.EV),
  code('REF', 'REF01', 'F2', HEADER_REF_NAME.F2, HEADER, cms('2'), HEADER_REF_NAME.F2),
  code('DTM', 'DTM01', '405', DTM_NAME['405'], HEADER, cms('2'), DTM_NAME['405']),
  code('N1', 'N101', 'PR', N1_NAME.PR, PAYER, cms('2'), N1_NAME.PR),
  code('N1', 'N103', 'XV', ID_QUALIFIER_NAME_MORE.XV, PAYER, cms('2'), N1_NAME.PR),
  code('REF', 'REF01', '2U', UNNAMED, PAYER, `${cms('3')}; ${cedi('24', '18')}`, PAYER_LOOP_NAME.REF),
  code('N1', 'N101', 'PE', N1_NAME.PE, PAYEE, cms('3'), N1_NAME.PE),
  code('N1', 'N103', 'XX', ID_QUALIFIER_NAME.XX, PAYEE, `${cms('7')}; ${cedi('32', '26')}`, N1_NAME.PE),
  code('N1', 'N103', 'FI', UNNAMED, PAYEE, cms('3'), N1_NAME.PE),
  code('REF', 'REF01', 'TJ', REF_QUALIFIER_NAME.TJ, PAYEE, `${cms('4')}; CMS 835 TI CG (2012), PDF p.6`, PAYEE_LOOP_NAME.REF),
  code('PER', 'PER01', 'CX', PAYER_PER_NAME.CX, PAYER, cms('3'), PAYER_PER_NAME.CX),
  code('PER', 'PER01', 'BL', PAYER_PER_NAME.BL, PAYER, cms('3'), PAYER_PER_NAME.BL),
  code('PER', 'PER01', '1C', PAYER_PER_NAME['1C'], PAYER, cms('3'), PAYER_PER_NAME['1C']),
  code('LX', 'LX01', '1', LX01_NAME['1'], SUMMARY, `${cedi('25', '19')}; ${cms('4')}`, SEGMENT_NAME.LX),
  code('LX', 'LX01', '0', LX01_NAME['0'], SUMMARY, `${cedi('25', '19')}; ${cms('4')}`, SEGMENT_NAME.LX),
  ...['1', '2', '3', '4', '19', '20', '21', '22', '23'].map((value) => code('CLP', 'CLP02', value, CLP02_NAME[value], CLAIM, `${cms('5')}; ${cedi('25', '19')}`, SEGMENT_NAME.CLP)),
  code('CLP', 'CLP06', 'MA', CLP06_NAME.MA, CLAIM, `${cms('6')}; CMS 835 TI CG (2012), PDF p.7`, SEGMENT_NAME.CLP),
  code('CLP', 'CLP06', 'MB', CLP06_NAME.MB, CLAIM, `${cms('6')}; ${cedi('25', '19')}`, SEGMENT_NAME.CLP),
  code('NM1', 'NM101', 'QC', NM1_NAME.QC, CLAIM, cms('6'), NM1_NAME.QC),
  code('NM1', 'NM101', '74', NM1_NAME['74'], CLAIM, cms('7'), NM1_NAME['74']),
  code('NM1', 'NM101', '82', NM1_NAME['82'], CLAIM, cms('7'), NM1_NAME['82']),
  code('NM1', 'NM101', 'TT', NM1_NAME.TT, CLAIM, cms('7'), NM1_NAME.TT),
  code('NM1', 'NM101', 'PR', NM1_NAME.PR, CLAIM, cms('7-8'), NM1_NAME.PR),
  // NM102 and NM108 values are left broad (no segmentName): each applies to
  // whichever roles carry it, and the entry says which.
  code('NM1', 'NM102', '1', NM102_NAME['1'], CLAIM, cms('6-8')),
  code('NM1', 'NM102', '2', NM102_NAME['2'], CLAIM, cms('7-8')),
  code('NM1', 'NM108', 'MI', ID_QUALIFIER_NAME_MORE.MI, CLAIM, `${cms('7')}; ${cedi('25', '19')}`, NM1_NAME.QC),
  code('NM1', 'NM108', 'C', UNNAMED, CLAIM, cms('7'), NM1_NAME['74']),
  code('NM1', 'NM108', 'XX', ID_QUALIFIER_NAME.XX, CLAIM, `${cms('7')}; ${cedi('32', '26')}`),
  code('NM1', 'NM108', 'PI', ID_QUALIFIER_NAME_MORE.PI, CLAIM, cms('7-8')),
  code('NM1', 'NM108', 'XV', ID_QUALIFIER_NAME_MORE.XV, CLAIM, cms('7-8')),
  code('REF', 'REF01', 'EA', CLAIM_REF_NAME.EA, CLAIM, `${cms('10')}; ${cedi('26', '20')}`, CLAIM_REF_NAME.EA),
  code('REF', 'REF01', 'F8', CLAIM_REF_NAME.F8, CLAIM, `${cms('10')}; ${cedi('26', '20')}`, CLAIM_REF_NAME.F8),
  code('DTM', 'DTM01', '232', DTM_NAME['232'], CLAIM, cms('10'), DTM_NAME['232']),
  code('DTM', 'DTM01', '233', DTM_NAME['233'], CLAIM, cms('10'), DTM_NAME['233']),
  code('DTM', 'DTM01', '036', DTM_NAME['036'], CLAIM, cms('10'), DTM_NAME['036']),
  code('DTM', 'DTM01', '050', DTM_NAME['050'], CLAIM, cms('11'), DTM_NAME['050']),
  code('PER', 'PER01', 'CX', CLAIM_PER_NAME.CX, CLAIM, cms('11'), CLAIM_PER_NAME.CX),
  ...['AU', 'DY', 'F5', 'I', 'NL', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'].map((value) => code('AMT', 'AMT01', value, CLAIM_AMT_QUALIFIER_NAME[value], CLAIM, `${cms('11')}; ${cedi('26', '20')}`, CLAIM_AMT_NAME)),
  ...['CA', 'CD', 'LA', 'OU', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'].map((value) => code('QTY', 'QTY01', value, QTY_QUALIFIER_NAME[value], CLAIM, `${cms('11')}; ${cedi('27', '21')}`, CLAIM_QTY_NAME)),
  code('CAS', 'CAS01', 'CO', CAS_GROUP_NAME.CO, `${CLAIM} · ${LINE}`, X12_GROUP_CODES),
  code('CAS', 'CAS01', 'OA', CAS_GROUP_NAME.OA, `${CLAIM} · ${LINE}`, X12_GROUP_CODES),
  code('CAS', 'CAS01', 'PR', CAS_GROUP_NAME.PR, `${CLAIM} · ${LINE}`, X12_GROUP_CODES),
  code('SVC', 'SVC01-1', 'HC', SVC_QUALIFIER_NAME.HC, LINE, `${cedi('31', '25')}; ${cms('11')}`, SEGMENT_NAME.SVC),
  code('SVC', 'SVC01-1', 'HP', UNNAMED, LINE, `${cms('11')}; ${cedi('27', '21')}`, SEGMENT_NAME.SVC),
  code('SVC', 'SVC01-1', 'N4', UNNAMED, LINE, `${cms('11')}; ${cedi('27', '21')}`, SEGMENT_NAME.SVC),
  code('SVC', 'SVC01-1', 'NU', UNNAMED, LINE, `${cms('11')}; ${cedi('27', '21')}`, SEGMENT_NAME.SVC),
  code('DTM', 'DTM01', '472', DTM_NAME['472'], LINE, cms('12'), DTM_NAME['472']),
  code('DTM', 'DTM01', '150', LINE_DTM_NAME['150'], LINE, cms('12'), DTM_NAME['472']),
  code('DTM', 'DTM01', '151', LINE_DTM_NAME['151'], LINE, cms('12'), DTM_NAME['472']),
  code('REF', 'REF01', '6R', SERVICE_REF_NAME['6R'], LINE, cms('13'), SERVICE_REF_NAME['6R']),
  code('REF', 'REF01', 'LU', SERVICE_REF_NAME.LU, LINE, `${cms('12')}; ${cedi('27', '21')}`, SERVICE_REF_NAME.LU),
  code('REF', 'REF01', '1S', UNNAMED, LINE, `${cms('12')}; ${cedi('27', '21')}`, SERVICE_REF_NAME.LU),
  code('REF', 'REF01', 'APC', REF_QUALIFIER_NAME.APC, LINE, `${cms('12')}; ${cedi('27', '21')}`, SERVICE_REF_NAME.LU),
  code('REF', 'REF01', 'RB', REF_QUALIFIER_NAME.RB, LINE, `${cms('12')}; ${cedi('27', '21')}`, SERVICE_REF_NAME.LU),
  code('REF', 'REF01', 'HPI', SERVICE_REF_NAME.HPI, LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME.HPI),
  code('REF', 'REF01', 'SY', UNNAMED, LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME.HPI),
  code('REF', 'REF01', 'TJ', REF_QUALIFIER_NAME.TJ, LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME.HPI),
  code('REF', 'REF01', '1C', REF_QUALIFIER_NAME['1C'], LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME.HPI),
  code('REF', 'REF01', '1G', REF_QUALIFIER_NAME['1G'], LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME.HPI),
  code('REF', 'REF01', '0K', SERVICE_REF_NAME['0K'], LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_REF_NAME['0K']),
  ...['B6', 'KH', '2K', 'ZL', 'ZM', 'ZN', 'ZO'].map((value) => code('AMT', 'AMT01', value, SERVICE_AMT_QUALIFIER_NAME[value], LINE, `${cms('13')}; ${cedi('28', '22')}`, SERVICE_AMT_NAME)),
  code('LQ', 'LQ01', 'HE', LQ_NAME.HE, LINE, `${cms('13')}; ${cedi('28', '22')}`, LQ_NAME.HE),
  ...Object.entries(PLB_REASON_NAME).map(([value, name]) => code('PLB', 'PLB03-1', value, name, PROVIDER, X12_PLB_REASONS, SEGMENT_NAME.PLB)),
  ...Object.entries(CARC).map(([value, text]) => ({
    ...code('CAS', 'CAS02', value, text.short, `${CLAIM} · ${LINE}`, carcSource),
    full: text.full,
  })),
  ...Object.entries(RARC).map(([value, text]) => ({
    ...code('MOA', 'MOA03', value, text.short, `${CLAIM} (MIA05, MIA20–23, MOA03–07) · ${LINE} (LQ02)`, rarcSource),
    full: text.full,
  })),
]
