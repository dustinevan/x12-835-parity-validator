import { useEffect, useState } from 'react'

// Typed client for docs/API.md. Every endpoint is GET under /api/ and returns
// JSON. Amounts, rates, and dates stay strings exactly as the server sends
// them; format.ts turns them into display text.

// ---------------------------------------------------------------- shared shapes

export type ClaimType = 'institutional' | 'professional' | 'dme'
export type Outcome = 'equal' | 'different' | 'missing' | 'duplicated' | 'extra' | 'set_aside'
export type Unmatched = 'missing' | 'extra' | 'duplicated' | 'set_aside'
export type Sort = 'delta' | '-delta' | 'paid' | 'charge'

export interface Key {
  mbi: string
  date_of_service: string
  total_charge: string
  tax_id: string
  npi: string
  workload: string
}

// A claim as stored: its bytes as text and its file's separators, which the
// claim screen reads in the browser (docs/compare835/file.ts parseClaim).
export interface Claim {
  segments: string // the claim's bytes as text, CLP through its last line
  delimiters: Delimiters // the four separators its file used
}

// The one-character separators of a claim's file, so its segments text can
// be read back (the docs site's 835 reader, docs/compare835/file.ts).
export interface Delimiters {
  element: string
  component: string
  repetition: string
  terminator: string
}

export interface Difference {
  element: string
  occurrence: number
  line_control_number: string | null
  mac_value: string | null
  vendor_value: string | null
  amount_difference: string | null
}

export interface Submission {
  prefix: string
  side: 'mac' | 'vendor'
  label: string
  submitted_on: string | null
  status: string
  claim_count: number
}

export interface Measure {
  claim_type: ClaimType | 'all'
  workload: string
  mac_claims: number
  equal: number
  different: number
  missing: number
  extra: number
  set_aside: number
  duplicated: number
  outcome_match_rate: string | null
  paid_claims: number
  median_deviation: string | null
  p95_deviation: string | null
}

export interface Tag {
  tag: string
  label: string
}

// ---------------------------------------------------------------- responses

export interface SubmissionsResponse {
  submissions: Submission[]
}

// previous and previous_measures, like DifferenceRow.previous_claims, give the
// same vendor's previous submission and its figures; no page shows them.
export interface MeasuresResponse {
  submission: Submission
  previous: Submission | null
  measures: Measure[]
  previous_measures: Measure[]
}

export interface DifferenceRow extends Tag {
  claims: number
  dollars: string
  previous_claims: number | null
}

export interface DifferencesResponse {
  different_claims: number
  rows: DifferenceRow[]
}

// /api/agreement: MAC claims counted by their outcome in two vendor
// submissions (equal, different, missing, or duplicated in each); pairs
// with no claims are left out. same_tags counts, among claims different in
// both, those whose differing elements are the same in both. The request
// and the response both name the two submissions v1 and v2.
export type Compared = 'equal' | 'different' | 'missing' | 'duplicated'

export interface AgreementCell {
  v1: Compared
  v2: Compared
  claims: number
  same_tags: number
}

export interface AgreementResponse {
  v1: Submission
  v2: Submission
  cells: AgreementCell[]
}

export interface Bucket {
  from: string
  to: string
  count: number
}

export interface ClaimRow {
  match_key: string
  key: Key
  claim_type: ClaimType
  mac_paid: string | null // null when that side is absent
  vendor_paid: string | null
  delta: string | null
}

// /api/claims: one page of the claims the combination filter keeps, with
// their count, the paid difference added up, and its histogram.
export interface ClaimsResponse {
  claims: number
  dollars: string
  histogram: Bucket[]
  rows: ClaimRow[]
  page: number
  pages: number
}

// The experts' combination filter: claims carrying every tag in `all`, at
// least one tag in `any`, and no tag in `none`. Empty lists are left out.
export interface TagFilter {
  all: string[]
  any: string[]
  none: string[]
}

export interface ClaimResponse {
  key: Key
  outcome: Outcome
  mac: Claim | null
  vendor: Claim | null // the first of vendor_claims
  vendor_claims: Claim[] // every vendor claim with the key, in file order (several when duplicated)
  differences: Difference[]
}

export interface UnmatchedRow {
  match_key: string
  key: Key
  claim_type: ClaimType
  patient_control_number: string
  paid: string
}

export interface UnmatchedResponse {
  outcome: Unmatched
  count: number
  share: string
  by_claim_type: { claim_type: ClaimType; count: number }[]
  by_workload: { workload: string; count: number }[]
  rows: UnmatchedRow[]
}

// /api/front: the Overview's numbers for one submission under a selection.
// Every filter is a comma-separated list of values, absent meaning all. The
// totals, outcomes, status matrix, and tags are under the whole selection;
// each dimension's values are under the selection on every other dimension
// and not its own, so a reader sees what each value would add back. legacy
// is the MAC submission's own figures (FrontLegacy).
export const FRONT_FILTERS = [
  'jurisdiction', 'part', 'workload', 'claim_type', 'adjudication', 'mac_status', 'vendor_status', 'frequency', 'paid_band', 'service_month', 'payment_band', 'outcome',
] as const
export type FrontFilter = (typeof FRONT_FILTERS)[number]

export interface FrontSums {
  claims: number // every comparison row
  mac_claims: number // the rows with a MAC claim: every outcome but extra
  mac_paid: string
  vendor_paid: string
  billed: string
  net_difference: string // vendor minus MAC over pairs, signed
  gross_difference: string // |vendor minus MAC| over pairs
  same_status: number // pairs with the same CLP02
}

export interface FrontOutcome extends FrontSums {
  outcome: Outcome
}

export interface FrontValue extends FrontSums {
  value: string
}

export interface FrontDimension {
  name: FrontFilter
  values: FrontValue[]
}

export interface FrontCell {
  mac_status: string
  vendor_status: string
  claims: number
  mac_paid: string
  vendor_paid: string
}

export interface FrontTag extends Tag {
  claims: number
  dollars: string
}

// The legacy side's own figures, from the MAC submission's claims: the MAC
// claims (each claim once), what the MAC paid on them (CLP04), and what was
// billed (CLP03).
export interface LegacySums {
  claims: number
  paid: string
  billed: string
}

export interface LegacyValue extends LegacySums {
  value: string
}

export interface LegacyDimension {
  name: FrontFilter
  values: LegacyValue[]
}

// legacy in /api/front: the totals and the nine dimensions the MAC claim
// carries (every filter but vendor_status, payment_band, and outcome),
// under the selection on those nine only, each dimension added under the
// selection on the others. The same for every vendor submission.
export interface FrontLegacy {
  totals: LegacySums
  dimensions: LegacyDimension[]
}

export interface FrontResponse {
  submission: Submission
  totals: FrontSums
  outcomes: FrontOutcome[]
  dimensions: FrontDimension[]
  status_matrix: FrontCell[]
  tags: FrontTag[]
  legacy: FrontLegacy
}

// /api/codes: the embedded CARC and RARC reference lists, code to entry.
// Static reference data (Cache-Control: public), fetched once and cached by
// app/CodeHint.tsx.
export interface CodeEntry {
  description: string
  status: 'current' | 'deactivated'
  start?: string
  stop?: string
  last_modified?: string
  notes?: string
}

export interface CodesResponse {
  carc: Record<string, CodeEntry>
  rarc: Record<string, CodeEntry>
}

// /api/macs: every MAC workload number as a tree (backend/pkg/macs Tree),
// jurisdictions in CMS's order. Within a Part: the roll-ups with the
// workloads under each, the workloads under no roll-up (a DME
// jurisdiction's number is one), the home health and hospice numbers (the
// note is the HHH region), and the numbers the numbering does not place,
// each with a note. Static reference data (Cache-Control: public), read by
// the MAC map (app/MacTree.tsx) and the Overview's workload box.
export interface MacState {
  code: string // USPS
  name: string
}

export interface MacNumber {
  id: string // five digits
  states: MacState[]
  note: string | null
}

export interface MacRollUp extends MacNumber {
  carries_hhh: boolean // the roll-up also carries the HHH claims
  workloads: MacNumber[]
}

export interface MacPart {
  part: string // 'A', 'B', 'DME', or '' when the table gives none
  roll_ups: MacRollUp[]
  workloads: MacNumber[]
  hhh: MacNumber[]
  other: MacNumber[]
}

export interface MacJurisdiction {
  jurisdiction: string
  contractor: string
  parts: MacPart[]
}

export interface MacsResponse {
  jurisdictions: MacJurisdiction[]
}

// ---------------------------------------------------------------- requests

type Params = Record<string, string | number | undefined | null>

function query(params: Params) {
  const search = new URLSearchParams()
  for (const [name, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== '') search.set(name, String(value))
  }
  const text = search.toString()
  return text ? `?${text}` : ''
}

// An error from a handler is JSON: { requestId, error } (docs/API.md). The
// request ID matches the server's log line. A 404 or 405 for a path the
// server does not route is plain text.
function errorMessage(body: string): string {
  try {
    const parsed = JSON.parse(body) as { requestId?: unknown; error?: unknown }
    if (typeof parsed.error === 'string') {
      return typeof parsed.requestId === 'string' ? `${parsed.error} (request ${parsed.requestId})` : parsed.error
    }
  } catch {
    // Not JSON: show the text.
  }
  return body
}

// A failed request keeps its HTTP status, so a view can tell "no such claim
// in this submission" (404) from a failure.
export class ApiError extends Error {
  readonly status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

export const isNotFound = (error: Error | undefined) => error instanceof ApiError && error.status === 404

async function get<T>(path: string, params: Params = {}, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`/api/${path}${query(params)}`, { signal, headers: { Accept: 'application/json' } })
  if (!response.ok) {
    const message = errorMessage((await response.text()).trim())
    throw new ApiError(response.status, `${response.status} ${response.statusText}${message ? `: ${message.slice(0, 200)}` : ''}`)
  }
  return (await response.json()) as T
}

export const api = {
  submissions: (signal?: AbortSignal) => get<SubmissionsResponse>('submissions', {}, signal),
  measures: (submission: string, signal?: AbortSignal) => get<MeasuresResponse>('measures', { submission }, signal),
  differences: (p: { submission: string; claim_type?: ClaimType; workload?: string }, signal?: AbortSignal) =>
    get<DifferencesResponse>('differences', p, signal),
  agreement: (p: { v1: string; v2: string }, signal?: AbortSignal) => get<AgreementResponse>('agreement', p, signal),
  // other and other_outcome go together: the claims' outcome in another
  // vendor submission.
  claims: (p: TagFilter & {
    submission: string; claim_type?: ClaimType; workload?: string; outcome?: Outcome
    other?: string; other_outcome?: Outcome; sort?: Sort; page?: number
  }, signal?: AbortSignal) =>
    get<ClaimsResponse>('claims', { ...p, all: p.all.join(','), any: p.any.join(','), none: p.none.join(',') }, signal),
  claim: (p: { submission: string; match_key: string }, signal?: AbortSignal) => get<ClaimResponse>('claim', p, signal),
  unmatched: (p: { submission: string; outcome: Unmatched }, signal?: AbortSignal) => get<UnmatchedResponse>('unmatched', p, signal),
  codes: (signal?: AbortSignal) => get<CodesResponse>('codes', {}, signal),
  macs: (signal?: AbortSignal) => get<MacsResponse>('macs', {}, signal),
  front: (p: { submission: string } & Partial<Record<FrontFilter, string>>, signal?: AbortSignal) => get<FrontResponse>('front', p, signal),
}

// ---------------------------------------------------------------- hook

export interface Loaded<T> {
  data: T | undefined
  error: Error | undefined
  loading: boolean
}

// Runs `load` whenever `key` changes; `key` is null while an input is not
// ready yet. A newer request aborts the older one, so a slow response never
// overwrites a newer one.
export function useApi<T>(key: string | null, load: (signal: AbortSignal) => Promise<T>): Loaded<T> {
  const [state, setState] = useState<{ key: string | null; data?: T; error?: Error }>({ key: null })
  useEffect(() => {
    if (key === null) return
    const controller = new AbortController()
    load(controller.signal).then(
      (data) => setState({ key, data }),
      (error: unknown) => {
        if (controller.signal.aborted) return
        setState({ key, error: error instanceof Error ? error : new Error(String(error)) })
      },
    )
    return () => controller.abort()
    // `load` is recreated each render; `key` carries its inputs.
  }, [key])
  const current = state.key === key
  return {
    data: current ? state.data : undefined,
    error: current ? state.error : undefined,
    loading: key !== null && (!current || (state.data === undefined && state.error === undefined)),
  }
}

// ---------------------------------------------------------------- submissions

export function useSubmissions() {
  return useApi('submissions', (signal) => api.submissions(signal))
}

// Vendor labels in the server's order ('Vendor 1', 'Vendor 2').
export function vendorLabels(submissions: Submission[]) {
  const out: string[] = []
  for (const d of submissions) if (d.side === 'vendor' && !out.includes(d.label)) out.push(d.label)
  return out
}

// One vendor's submissions, oldest first.
export function submissionsOf(submissions: Submission[], label: string) {
  return submissions
    .filter((d) => d.side === 'vendor' && d.label === label)
    .sort((a, b) => (a.submitted_on ?? '').localeCompare(b.submitted_on ?? ''))
}

// True when both submissions are of one vendor, so a label alone does not
// tell the two apart.
export const sharesVendor = (x: Submission | undefined, y: Submission | undefined) =>
  x !== undefined && y !== undefined && x.label === y.label

// A submission's name for a page that shows two: the vendor, and the date
// as well when `qualified` (see sharesVendor).
export const nameOf = (s: Submission, qualified: boolean) => (qualified ? `${s.label} · ${s.submitted_on ?? s.prefix}` : s.label)
