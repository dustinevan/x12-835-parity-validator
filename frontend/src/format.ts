// Display helpers for the API's strings. Amounts arrive as exact decimal
// strings; arithmetic on them (a difference, a sum) uses scaled integers so
// nothing is rounded. Numbers are used only for bar widths and percentages.

// ---------------------------------------------------------------- exact decimals

interface Scaled {
  units: bigint
  scale: number
}

function parse(text: string): Scaled {
  const s = text.trim()
  const negative = s.startsWith('-')
  const body = s.replace(/^[+-]/, '')
  const [whole, fraction = ''] = body.split('.')
  const digits = `${whole || '0'}${fraction}`.replace(/[^0-9]/g, '') || '0'
  const units = BigInt(digits)
  return { units: negative ? -units : units, scale: fraction.length }
}

function rescale(v: Scaled, scale: number): bigint {
  return v.units * 10n ** BigInt(scale - v.scale)
}

function render(units: bigint, scale: number): string {
  const negative = units < 0n
  const digits = (negative ? -units : units).toString().padStart(scale + 1, '0')
  const whole = digits.slice(0, digits.length - scale)
  const fraction = scale > 0 ? `.${digits.slice(digits.length - scale)}` : ''
  return `${negative ? '-' : ''}${whole}${fraction}`
}

// a − b, exactly.
export function subtract(a: string, b: string): string {
  const x = parse(a)
  const y = parse(b)
  const scale = Math.max(x.scale, y.scale)
  return render(rescale(x, scale) - rescale(y, scale), scale)
}

// The sum of the amounts, exactly.
export function sum(amounts: string[]): string {
  const parsed = amounts.map(parse)
  const scale = Math.max(2, ...parsed.map((p) => p.scale))
  let total = 0n
  for (const p of parsed) total += rescale(p, scale)
  return render(total, scale)
}

// −1, 0, or 1; '10.0' and '10.00' compare equal.
export function compare(a: string, b: string): number {
  const x = parse(a)
  const y = parse(b)
  const scale = Math.max(x.scale, y.scale)
  const d = rescale(x, scale) - rescale(y, scale)
  return d < 0n ? -1 : d > 0n ? 1 : 0
}

export const isZero = (a: string) => compare(a, '0') === 0

// Whether a and b are `cents` cents apart or less, either way, edges
// included, exactly: 80.00 and 79.95 are within 5 cents, 80.00 and 79.94 are
// not.
export function withinCents(a: string, b: string, cents: number): boolean {
  const x = parse(a)
  const y = parse(b)
  const scale = Math.max(x.scale, y.scale, 2)
  const d = rescale(x, scale) - rescale(y, scale)
  const limit = BigInt(cents) * 10n ** BigInt(scale - 2)
  return d <= limit && d >= -limit
}

// A decimal-looking value: amounts compare as numbers, codes as text.
export const isAmount = (s: string) => /^[+-]?\d+(\.\d+)?$/.test(s.trim()) && s.includes('.')

// ---------------------------------------------------------------- display

export const count = (n: number) => n.toLocaleString('en-US')

function group(whole: string) {
  return whole.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// '$1,234.50' from '1234.5'; exact, at least two decimals, sign dropped.
export function money(amount: string): string {
  const s = amount.trim().replace(/^[+-]/, '')
  const [whole, fraction = ''] = s.split('.')
  return `$${group(whole.replace(/^0+(?=\d)/, '') || '0')}.${fraction.padEnd(2, '0')}`
}

// '+$1.26', '−$1.26', or '$0.00'.
export function signedMoney(amount: string): string {
  const c = compare(amount, '0')
  if (c === 0) return '$0.00'
  return `${c > 0 ? '+' : '−'}${money(amount)}`
}

// A fraction as a percent: '0.941' → '94.1%'.
export function rate(fraction: string | number | null | undefined, digits = 1): string {
  if (fraction === null || fraction === undefined) return '—'
  return `${(Number(fraction) * 100).toFixed(digits)}%`
}

export const CLAIM_TYPE_LABEL: Record<string, string> = {
  institutional: 'Institutional',
  professional: 'Professional',
  dme: 'DME',
}

export const PART_LABEL: Record<string, string> = {
  institutional: 'Part A',
  professional: 'Part B',
  dme: 'DME',
}

export const OUTCOME_LABEL: Record<string, string> = {
  equal: 'Equal',
  different: 'Different',
  missing: 'Missing from vendor',
  duplicated: 'Duplicated by vendor',
  extra: 'Extra on vendor',
  set_aside: 'Set aside',
}

// CLP02 claim status codes (X12 005010X221A1).
export const CLAIM_STATUS: Record<string, string> = {
  '1': 'Processed as primary',
  '2': 'Processed as secondary',
  '3': 'Processed as tertiary',
  '4': 'Denied',
  '19': 'Primary, forwarded',
  '20': 'Secondary, forwarded',
  '21': 'Tertiary, forwarded',
  '22': 'Reversal of previous payment',
  '23': 'Not our claim, forwarded',
}

export const statusLabel = (code: string) => CLAIM_STATUS[code] ?? `Status ${code}`

// '' is kept (an explicit empty slot, selection.tsx), unlike undefined/null.
export function withParams(path: string, params: Record<string, string | undefined | null>) {
  const search = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) if (v !== undefined && v !== null) search.set(k, v)
  const text = search.toString()
  return text ? `${path}?${text}` : path
}

// The share of MAC claims whose vendor claim is exactly equal, as a fraction.
// The denominator is every MAC claim except those set aside for sharing a
// key: a missing vendor claim counts against the rate, and an extra vendor
// claim is not a MAC claim. Every view shows "exactly equal" this one way.
export function equalRate(m: { mac_claims: number; set_aside: number; equal: number }): number | undefined {
  const eligible = m.mac_claims - m.set_aside
  return eligible <= 0 ? undefined : m.equal / eligible
}

// " (2nd)" for the second occurrence of something, "" for the first.
export const ordinal = (n: number) => (n === 1 ? '' : ` (${n}${n === 2 ? 'nd' : n === 3 ? 'rd' : 'th'})`)
