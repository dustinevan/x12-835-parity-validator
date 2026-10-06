import { FRONT_FILTERS, type FrontFilter, type FrontValue } from '../api'
import { CLAIM_TYPE_LABEL, OUTCOME_LABEL, statusLabel } from '../format'

// The overview's selection: which values of each dimension the numbers
// are over. It lives in the address under the names /api/front takes
// (docs/API.md), one parameter per dimension, a comma-separated list of
// values; a dimension absent from the address means every value. A list
// with nothing in it (every value left out) is written as "~", a value no
// row carries, so the server answers with nothing rather than with everything.
// The parameters are codes and bands, never identifiers, so an address
// carries no PHI.

export interface Dimension {
  name: FrontFilter
  label: string
}

// Every dimension, in the order the response lists them and the breakdown
// tables show them.
export const DIMENSIONS: Dimension[] = [
  { name: 'jurisdiction', label: 'Jurisdiction' },
  { name: 'part', label: 'Part' },
  { name: 'workload', label: 'Workload' },
  { name: 'claim_type', label: 'Claim type' },
  { name: 'adjudication', label: 'Inpatient or outpatient' },
  { name: 'mac_status', label: 'MAC claim status' },
  { name: 'vendor_status', label: 'Vendor claim status' },
  { name: 'frequency', label: 'Claim frequency' },
  { name: 'paid_band', label: 'MAC paid band' },
  { name: 'service_month', label: 'Service month' },
  { name: 'payment_band', label: 'Payment difference band' },
  { name: 'outcome', label: 'Outcome' },
]

export const dimensionLabel = (name: FrontFilter): string => DIMENSIONS.find((d) => d.name === name)?.label ?? name

// The dimensions offered as boxes in the selection header (SelectionHeader.tsx),
// in the order the header shows them: eight boxes of tabs, then the
// workload box, which lists the workloads under their jurisdictions. Part,
// MAC paid band, and payment difference band have no box: they stay filters
// of /api/front, so a link can still narrow them, and the page shows them as
// content instead (the breakdown tables and the bands chart).
export const BOX_DIMENSIONS: FrontFilter[] = [
  'jurisdiction', 'claim_type', 'adjudication', 'mac_status', 'vendor_status', 'frequency', 'service_month', 'outcome', 'workload',
]

// The dimensions only the address can narrow.
export const ADDRESS_ONLY: FrontFilter[] = FRONT_FILTERS.filter((name) => !BOX_DIMENSIONS.includes(name))

// The sixteen MAC jurisdictions in CMS's order, the order /api/macs lists
// them in (backend/pkg/macs): the twelve A/B jurisdictions, then the four
// DME ones, so the jurisdiction box's rows of four end with the DME row.
// scripts/ssr-check-app.mjs checks this list against the /api/macs fixture.
export const JURISDICTIONS = ['J5', 'J6', 'J8', 'J15', 'JE', 'JF', 'JH', 'JJ', 'JK', 'JL', 'JM', 'JN', 'JA', 'JB', 'JC', 'JD']

export type OverviewSelection = Partial<Record<FrontFilter, string[]>>

export const NONE = '~'

export function readSelection(params: URLSearchParams): OverviewSelection {
  const out: OverviewSelection = {}
  for (const name of FRONT_FILTERS) {
    const text = params.get(name)
    if (text === null) continue
    out[name] = text === NONE ? [] : text.split(',')
  }
  return out
}

export function writeSelection(params: URLSearchParams, selection: OverviewSelection): void {
  for (const name of FRONT_FILTERS) {
    const list = selection[name]
    if (list === undefined) params.delete(name)
    else params.set(name, list.length === 0 ? NONE : list.join(','))
  }
}

// The request's filter parameters: the address's own spelling.
export function filtersOf(selection: OverviewSelection): Partial<Record<FrontFilter, string>> {
  const out: Partial<Record<FrontFilter, string>> = {}
  for (const name of FRONT_FILTERS) {
    const list = selection[name]
    if (list !== undefined) out[name] = list.length === 0 ? NONE : list.join(',')
  }
  return out
}

// One text for a selection, to key a request on.
export const selectionKey = (selection: OverviewSelection): string => new URLSearchParams(filtersOf(selection)).toString()

export const isNarrowed = (selection: OverviewSelection): boolean => FRONT_FILTERS.some((name) => selection[name] !== undefined)

// Whether the numbers count a value's claims: every value when the
// dimension has no list, else the values in it.
export const isSelected = (selection: OverviewSelection, name: FrontFilter, value: string): boolean => {
  const list = selection[name]
  return list === undefined || list.includes(value)
}

// Whether a value is picked: in the dimension's list. With no list (All)
// nothing is picked, and the numbers count every value.
export const isPicked = (selection: OverviewSelection, name: FrontFilter, value: string): boolean => selection[name]?.includes(value) ?? false

// The selection after a click on one value's tab: a value not picked joins
// the dimension's list, a picked one leaves it. From All or None the click
// picks that value alone. Unpicking the last one removes the list, so the
// dimension counts every value again.
export function toggleValue(selection: OverviewSelection, name: FrontFilter, value: string): OverviewSelection {
  const list = selection[name] ?? []
  const next = list.includes(value) ? list.filter((v) => v !== value) : [...list, value]
  return withList(selection, name, next.length === 0 ? undefined : next)
}

export const selectAll = (selection: OverviewSelection, name: FrontFilter): OverviewSelection => withList(selection, name, undefined)
export const selectNone = (selection: OverviewSelection, name: FrontFilter): OverviewSelection => withList(selection, name, [])

function withList(selection: OverviewSelection, name: FrontFilter, list: string[] | undefined): OverviewSelection {
  const next = { ...selection }
  if (list === undefined) delete next[name]
  else next[name] = list
  return next
}

// ---------------------------------------------------------------- the values a box shows

export interface BoxValue {
  value: string
  claims: number // under what the other boxes pick, the picked submissions added together
  picked: boolean
}

// The values a box shows: those with claims under the selection on the
// other dimensions, in either picked submission (each response's facet for
// the dimension, which /api/front adds under the selection on every other
// dimension and lists only with claims), and every value the box has
// picked, even one with no claims there, so it can be unpicked. A value
// with no claims that is not picked is not shown: picking J5 leaves only
// J5's workloads in the workload box. `facets` holds each picked
// submission's values for the dimension.
export function boxValues(facets: FrontValue[][], name: FrontFilter, selection: OverviewSelection): BoxValue[] {
  const claims = new Map<string, number>()
  for (const values of facets) for (const v of values) claims.set(v.value, (claims.get(v.value) ?? 0) + v.claims)
  const picked = selection[name] ?? []
  const shown = new Set<string>(picked)
  for (const [value, n] of claims) if (n > 0) shown.add(value)
  return [...shown]
    .map((value) => ({ value, claims: claims.get(value) ?? 0, picked: picked.includes(value) }))
    .sort((x, y) => boxOrder(name, x.value, y.value))
}

// A box's order: the jurisdictions in CMS's order, anything else after
// them; every other dimension in its own order (compareValues).
function boxOrder(name: FrontFilter, x: string, y: string): number {
  if (name !== 'jurisdiction') return compareValues(name, x, y)
  const rank = (v: string) => {
    const at = JURISDICTIONS.indexOf(v)
    return at >= 0 ? at : JURISDICTIONS.length + (v === '' ? 1 : 0)
  }
  return rank(x) - rank(y) || x.localeCompare(y)
}

// ---------------------------------------------------------------- labels and order

export const PAID_BANDS = ['zero', 'under 100', 'under 1000', 'under 10000', '10000 and more']
// The payment bands, how far apart a pair's two payments (CLP04) are: match
// is within five cents, which the comparison counts as the same amount;
// the empty value is a row that is not a pair.
export const PAYMENT_BANDS = ['match', 'up to 1.00', 'up to 1%', 'up to 5%', 'over 5%', '']
// The outcomes with a MAC claim first, then extra, as the overview's bars
// draw them.
export const OUTCOMES = ['equal', 'different', 'missing', 'duplicated', 'set_aside', 'extra']

const PAID_BAND_LABEL: Record<string, string> = {
  zero: '$0',
  'under 100': 'under $100',
  'under 1000': '$100 to under $1,000',
  'under 10000': '$1,000 to under $10,000',
  '10000 and more': '$10,000 and more',
}

const PAYMENT_BAND_LABEL: Record<string, string> = {
  match: 'match',
  'up to 1.00': 'up to $1.00 apart',
  'up to 1%': 'up to 1% of MAC paid',
  'up to 5%': 'up to 5% of MAC paid',
  'over 5%': 'over 5% of MAC paid',
  '': 'not a pair',
}

// The claim frequency codes the project's documents name; others show as sent.
const FREQUENCY_LABEL: Record<string, string> = { '1': 'original', '7': 'replacement', '8': 'void' }

export function valueLabel(name: FrontFilter, value: string): string {
  if (value === '') {
    switch (name) {
      case 'payment_band': return PAYMENT_BAND_LABEL['']
      case 'adjudication': return 'neither (no MIA or MOA)'
      case 'mac_status': return 'no MAC claim'
      case 'vendor_status': return 'no vendor claim'
      default: return 'unknown'
    }
  }
  switch (name) {
    case 'claim_type': return CLAIM_TYPE_LABEL[value] ?? value
    case 'outcome': return OUTCOME_LABEL[value] ?? value
    case 'mac_status':
    case 'vendor_status': return `${value} · ${statusLabel(value)}`
    case 'part': return value === 'DME' ? 'DME' : `Part ${value}`
    case 'adjudication': return value === 'inpatient' ? 'Inpatient (MIA)' : value === 'outpatient' ? 'Outpatient (MOA)' : value
    case 'frequency': return FREQUENCY_LABEL[value] ? `${value} · ${FREQUENCY_LABEL[value]}` : value
    case 'paid_band': return PAID_BAND_LABEL[value] ?? value
    case 'payment_band': return PAYMENT_BAND_LABEL[value] ?? value
    default: return value
  }
}

const OUTCOME_TAB_LABEL: Record<string, string> = {
  equal: 'Equal',
  different: 'Different',
  missing: 'Missing',
  duplicated: 'Duplicated',
  extra: 'Extra',
  set_aside: 'Set aside',
}

// A tab's text, shorter than valueLabel so a box stays a tight grid: a
// claim status or frequency tab shows its code, an inpatient or outpatient
// tab and an outcome tab one word. The tab's title carries valueLabel.
export function tabLabel(name: FrontFilter, value: string): string {
  switch (name) {
    case 'mac_status':
    case 'vendor_status': return value === '' ? 'no claim' : value
    case 'adjudication': return value === '' ? 'Neither' : value === 'inpatient' ? 'Inpatient' : value === 'outpatient' ? 'Outpatient' : value
    case 'frequency': return value === '' ? 'unknown' : value
    case 'outcome': return OUTCOME_TAB_LABEL[value] ?? valueLabel(name, value)
    default: return valueLabel(name, value)
  }
}

// Bands and outcomes in their own order; everything else as the server
// sorted it (by value), with the empty value last.
export function compareValues(name: FrontFilter, x: string, y: string): number {
  const order = name === 'paid_band' ? PAID_BANDS : name === 'payment_band' ? PAYMENT_BANDS : name === 'outcome' ? OUTCOMES : undefined
  if (order) return order.indexOf(x) - order.indexOf(y)
  if (x === '' || y === '') return x === '' ? 1 : -1
  const nx = Number(x)
  const ny = Number(y)
  if (!Number.isNaN(nx) && !Number.isNaN(ny) && x !== '' && y !== '') return nx - ny
  return x.localeCompare(y)
}

// ---------------------------------------------------------------- the legacy column

// The dimensions the MAC claim itself carries, in the response's order: the
// legacy figures (/api/front's legacy block) are cut by these and by
// nothing a vendor decides, so a breakdown by one of them has a Legacy
// column. The other three, the vendor's claim status, the payment band, and
// the outcome, are a vendor's results, and their breakdowns have none.
export const LEGACY_DIMENSIONS: FrontFilter[] = [
  'jurisdiction', 'part', 'workload', 'claim_type', 'adjudication', 'mac_status', 'frequency', 'paid_band', 'service_month',
]

// False where the comparison keeps no amount from the vendor: missing,
// duplicated, and set-aside rows carry no vendor claim, and neither do the
// rows with no vendor status.
export function hasVendorAmount(dimension?: FrontFilter, value?: string): boolean {
  if (dimension === 'outcome') return value !== 'missing' && value !== 'duplicated' && value !== 'set_aside'
  if (dimension === 'vendor_status') return value !== ''
  return true
}
