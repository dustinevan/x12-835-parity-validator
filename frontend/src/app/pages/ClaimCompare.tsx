import { useMemo, type ReactNode } from 'react'
import { Link, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, isNotFound, nameOf, sharesVendor, useApi, type Claim, type ClaimResponse, type Key, type Loaded, type Submission } from '../../api'
import { parseClaim, type ParsedClaim } from '../../docs/compare835/file'
import { count, money } from '../../format'
import type { Col, Sources } from '../../linePairing'
import { ClaimCompareView } from '../ClaimCompareView'
import { HELP } from '../help'
import { Info } from '../Info'
import { useSelection, type Selection, type Side } from '../selection'

// The claim screen: the MAC-processed claim beside the claim with the same
// match key in each of the two picked submissions, as the 835 compare page
// draws three files, one row per element. Two /api/claim calls (one per
// slot) each return the same MAC claim; every claim comes with its segments
// text and its file's delimiters, read here into the compare page's shape
// (docs/compare835/file.ts parseClaim). The page opens from a claim link on
// another page (match_key in the address); picking other submissions keeps
// the claim.

const help = HELP.claim

// ---------------------------------------------------------------- one vendor

// One vendor column: its submission, its /api/claim result, and what the
// column can show. A 404 means the submission has no result for this key.
interface Vendor {
  slot: Side
  submission: Submission
  data?: ClaimResponse
  error?: Error // a failure other than 404
  absent: boolean
  loading: boolean
  claim: Claim | null // a compared claim, or an extra or duplicated one shown for reference
}

function vendorOf(slot: Side, submission: Submission, loaded: Loaded<ClaimResponse>): Vendor {
  const data = loaded.data
  const shown = data?.outcome === 'equal' || data?.outcome === 'different' || data?.outcome === 'extra' || data?.outcome === 'duplicated'
  return {
    slot,
    submission,
    data,
    error: isNotFound(loaded.error) ? undefined : loaded.error,
    absent: isNotFound(loaded.error),
    loading: loaded.loading,
    claim: data && shown ? data.vendor : null,
  }
}

// The claim with this match key in each picked submission: one /api/claim
// call per slot (both return the same MAC claim), none until the address
// carries a match key. A slot with no submission has no entry.
function useVendors(sel: Selection, matchKey: string | undefined): Vendor[] {
  const v1 = sel.v1?.prefix
  const v2 = sel.v2?.prefix
  const claim1 = useApi(v1 && matchKey ? `claim:${v1}|${matchKey}` : null, (signal) => api.claim({ submission: v1!, match_key: matchKey! }, signal))
  const claim2 = useApi(v2 && matchKey ? `claim:${v2}|${matchKey}` : null, (signal) => api.claim({ submission: v2!, match_key: matchKey! }, signal))
  const vendors: Vendor[] = []
  if (sel.v1) vendors.push(vendorOf('v1', sel.v1, claim1))
  if (sel.v2) vendors.push(vendorOf('v2', sel.v2, claim2))
  return vendors
}

// What a column's head says under the submission: the outcome decides it,
// and only a compared, equal pair is a Match.
function headlineOf(v: Vendor): string {
  if (v.loading) return 'Loading'
  if (v.absent) return 'No claim in this submission'
  if (!v.data) return 'Could not load'
  const n = v.data.differences.length
  switch (v.data.outcome) {
    case 'equal':
      return 'Match'
    case 'different':
      return `${count(n)} ${n === 1 ? 'difference' : 'differences'}`
    case 'missing':
      return 'No claim in this submission'
    case 'duplicated':
      return `first of ${count(v.data.vendor_claims.length)} copies, for reference`
    case 'extra':
      return 'No MAC claim to compare'
    case 'set_aside':
      return 'Set aside, not compared'
  }
}

// ---------------------------------------------------------------- the page

const KEY_FIELDS: [keyof Key, string][] = [
  ['mbi', 'Patient'],
  ['date_of_service', 'Service date'],
  ['total_charge', 'Billed'],
  ['npi', 'Payee NPI'],
  ['tax_id', 'Tax ID'],
  ['workload', 'Workload'],
]

function KeyStrip({ claimKey }: { claimKey: Key }) {
  return (
    <div className="mb-4 flex flex-wrap items-stretch border border-line">
      {KEY_FIELDS.map(([field, label]) => (
        <div key={field} className="border-r border-line px-4 py-2.5 last:border-r-0">
          <div className="text-[11px] text-ink-2">{label}</div>
          <div className="font-mono text-sm text-ink">{field === 'total_charge' ? money(claimKey[field]) : claimKey[field]}</div>
        </div>
      ))}
      <div className="flex items-center px-3">
        <Info term={help.terms.matchKey} />
      </div>
    </div>
  )
}

// A claim read into the compare page's shape, or why it could not be.
function readClaim(claim: Claim | null | undefined): { parsed?: ParsedClaim; error?: string } {
  if (!claim) return {}
  try {
    return { parsed: parseClaim(claim.segments, claim.delimiters) }
  } catch (err) {
    return { error: err instanceof Error ? err.message : String(err) }
  }
}

function headOf(v: Vendor, qualified: boolean): ReactNode {
  return (
    <>
      <span className="block font-semibold">{nameOf(v.submission, qualified)}</span>
      <span className="block font-mono text-[11px] font-normal text-ink-2">{headlineOf(v)}</span>
    </>
  )
}

export function ClaimCompare() {
  const sel = useSelection()
  const [params] = useSearchParams()
  const matchKey = params.get('match_key') ?? undefined
  const vendors = useVendors(sel, matchKey)

  const settled = vendors.length > 0 && vendors.every((v) => !v.loading)
  const ready = settled && vendors.some((v) => v.data)
  const results = vendors.flatMap((v) => (v.data ? [v.data] : []))
  const mac = results.find((r) => r.mac)?.mac ?? null
  const qualified = sharesVendor(vendors[0]?.submission, vendors[1]?.submission)

  const read = useMemo(() => {
    const claims: Sources<ParsedClaim> = {}
    const errors: string[] = []
    const take = (col: Col, claim: Claim | null | undefined, name: string) => {
      const r = readClaim(claim)
      if (r.parsed) claims[col] = r.parsed
      if (r.error) errors.push(`${name}: ${r.error}`)
    }
    take('mac', mac, 'MAC')
    for (const v of vendors) take(v.slot, v.claim, v.submission.label)
    return { claims, error: errors.length ? errors.join(' · ') : undefined }
    // The claims change with the responses; the vendors' loaded data is what varies.
  }, [mac, vendors.map((v) => `${v.slot}:${v.claim?.segments ?? ''}`).join('|')])

  const head: Record<Col, ReactNode> = {
    mac: (
      <>
        <span className="block font-semibold">MAC-processed</span>
        {!mac && <span className="block font-mono text-[11px] font-normal text-ink-2">no MAC claim</span>}
      </>
    ),
    v1: vendors.find((v) => v.slot === 'v1') ? headOf(vendors.find((v) => v.slot === 'v1')!, qualified) : <span className="text-ink-2">—</span>,
    v2: vendors.find((v) => v.slot === 'v2') ? headOf(vendors.find((v) => v.slot === 'v2')!, qualified) : <span className="text-ink-2">—</span>,
  }

  return (
    <div className="mx-auto max-w-[96rem] px-10 py-12">
      <PageTitle purpose={help.purpose}>Compare claim</PageTitle>
      {sel.unselected ? (
        <Empty>Select a submission in the navigation.</Empty>
      ) : !matchKey ? (
        <p className="text-sm text-ink-2">
          Open a claim from <Link to={sel.link('/claims')} className="text-link hover:underline">Claims</Link> or{' '}
          <Link to={sel.link('/unmatched')} className="text-link hover:underline">Unmatched</Link>.
        </p>
      ) : ready ? (
        <>
          <KeyStrip claimKey={results[0].key} />
          <ClaimCompareView claims={read.claims} head={head} error={read.error} />
        </>
      ) : (
        <Pending errors={[sel.error, ...vendors.map((v) => v.error)]} loading={sel.loading || vendors.some((v) => v.loading)}
          empty="Neither picked submission has a claim with this match key" />
      )}
    </div>
  )
}
