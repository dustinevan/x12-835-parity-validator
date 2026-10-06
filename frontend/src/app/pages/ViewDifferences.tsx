import { useState, type ReactNode } from 'react'
import { Link, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, useApi, type ClaimType, type DifferenceRow, type Measure } from '../../api'
import { CLAIM_TYPE_LABEL, PART_LABEL, compare, count, equalRate, isZero, money, sum } from '../../format'
import { HELP, type Term } from '../help'
import { Info } from '../Info'
import { CodeHint } from '../CodeHint'
import { SideToggle } from '../VendorPicker'
import { useSelection } from '../selection'

// Differences between MAC-processed and vendor-processed 835s for one vendor
// submission, grouped by the element that differs (the tag). Reads
// /api/differences and /api/measures.

type Effect = 'codes' | 'patient' | 'paid' | 'outcome'

const EFFECTS: Record<Effect, { label: string; fill: string }> = {
  paid: { label: 'Amounts', fill: 'bg-series-1' },
  patient: { label: 'Patient owes', fill: 'bg-[#009d9a]' },
  outcome: { label: 'Claim outcome', fill: 'bg-ink' },
  codes: { label: 'Codes only', fill: 'bg-mark-off' },
}

// What a tag touches, read from the element itself.
function effectOf(row: DifferenceRow): Effect {
  if (row.tag === 'CLP02') return 'outcome'
  if (row.tag === 'CLP05' || row.tag.startsWith('CAS*PR*')) return 'patient'
  if (!isZero(row.dollars)) return 'paid'
  return 'codes'
}

const help = HELP.differences

const TOP = 10
const CLAIM_TYPES: ClaimType[] = ['institutional', 'professional', 'dme']

// ---------------------------------------------------------------- numbers

function percent(part: number, whole: number) {
  if (whole === 0) return '0%'
  const p = (part / whole) * 100
  return `${p < 1 ? p.toFixed(1) : p.toFixed(0)}%`
}

// ---------------------------------------------------------------- controls

function Tabs<T extends string>({ options, value, onChange }: { options: { id: T; label: string; dot?: string }[]; value: T; onChange: (id: T) => void }) {
  return (
    <div className="inline-flex flex-wrap border border-line" role="tablist">
      {options.map((o) => (
        <button key={o.id} type="button" role="tab" aria-selected={value === o.id} onClick={() => onChange(o.id)}
          className={`flex items-center gap-2 border-r border-line px-4 py-2 text-sm last:border-r-0 ${value === o.id ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
          {o.dot && <span className={`h-2 w-2 rounded-full ${o.dot}`} />}
          {o.label}
        </button>
      ))}
    </div>
  )
}

function FilterRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <>
      <span className="pt-2.5 text-right text-xs text-ink-2">{label}</span>
      <div>{children}</div>
    </>
  )
}

// ---------------------------------------------------------------- pieces

function Stat({ label, info, value, sub, dot }: { label: string; info: Term; value: string; sub?: string; dot?: string }) {
  return (
    <div className="flex-1 border-r border-line px-5 py-4 last:border-r-0">
      <div className="flex items-center gap-1.5 text-[11px] text-ink-2">
        {dot && <span className={`h-2 w-2 rounded-full ${dot}`} />}
        {label}
        <Info term={info} />
      </div>
      <div className="mt-1 font-mono text-xl font-semibold text-ink">{value}</div>
      {sub && <div className="font-mono text-[11px] text-ink-2">{sub}</div>}
    </div>
  )
}

// One bar the width of all claims with a difference; each top kind is a
// segment, the long tail is one gray segment at the end.
function Concentration({ ranked, tail }: { ranked: DifferenceRow[]; tail: DifferenceRow[] }) {
  const tailClaims = tail.reduce((s, m) => s + m.claims, 0)
  const total = ranked.reduce((s, m) => s + m.claims, 0) + tailClaims
  if (total === 0) return null
  return (
    <div>
      <div className="flex h-8 w-full overflow-hidden">
        {ranked.map((m, i) => {
          const share = m.claims / total
          return (
            <div key={m.tag} title={`${m.tag} ${m.label}: ${count(m.claims)} claims`}
              style={{ width: `${share * 100}%`, minWidth: 3, opacity: 1 - i * 0.07 }}
              className="flex items-center justify-center overflow-hidden border-r-2 border-bg bg-warn font-mono text-[11px] font-semibold whitespace-nowrap text-black">
              {share > 0.09 && <><CodeHint kind="element" code={m.tag} label={m.label}>{m.tag}</CodeHint> · {percent(m.claims, total)}</>}
            </div>
          )
        })}
        {tailClaims > 0 && (
          <div title={`Long tail: ${count(tailClaims)} claims`} style={{ width: `${(tailClaims / total) * 100}%`, minWidth: 3 }}
            className="bg-mark-off" />
        )}
      </div>
      <div className="mt-1.5 flex justify-between font-mono text-[11px] text-ink-2">
        <span>top {ranked.length} kinds · {percent(total - tailClaims, total)}</span>
        {tail.length > 0 && <span>{tail.length} more kinds · {percent(tailClaims, total)}</span>}
      </div>
    </div>
  )
}

const GRID = 'grid grid-cols-[2rem_minmax(0,1.7fr)_minmax(0,1fr)_4.5rem_minmax(0,1fr)] items-center gap-x-5'

function KindRow({ m, rank, maxClaims, maxDollars, different, groupLink }: {
  m: DifferenceRow
  rank: number
  maxClaims: number
  maxDollars: number
  different: number
  groupLink: (tag: string) => string
}) {
  const effect = EFFECTS[effectOf(m)]
  const dollars = Number(m.dollars)
  return (
    <Link to={groupLink(m.tag)}
      className={`${GRID} border-b border-line py-2.5 pr-3 hover:bg-warn/15`}>
      <span className="text-right font-mono text-xs text-ink-2">{rank}</span>

      <div className="min-w-0">
        <div className="flex items-baseline gap-2">
          <span className="shrink-0 border border-line px-1.5 font-mono text-[11px] font-semibold text-ink">
            <CodeHint kind="element" code={m.tag} label={m.label}>{m.tag}</CodeHint>
          </span>
          {m.label !== m.tag && <span className="truncate text-sm text-ink">{m.label}</span>}
        </div>
        <div className="mt-0.5 flex items-center gap-3 text-[11px] text-ink-2">
          <span className="flex items-center gap-1"><span className={`h-2 w-2 ${effect.fill}`} />{effect.label}</span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <div className="h-4 flex-1">
          <div className="h-full bg-warn" style={{ width: `${(m.claims / maxClaims) * 100}%`, minWidth: 2 }} />
        </div>
        <div className="w-[4.5rem] text-right">
          <div className="font-mono text-xs font-semibold text-ink">{count(m.claims)}</div>
        </div>
      </div>

      <span className="text-right font-mono text-xs text-ink">{percent(m.claims, different)}</span>

      <div className="flex items-center gap-2">
        {dollars > 0 ? (
          <>
            <div className="h-4 flex-1">
              <div className="h-full bg-series-1" style={{ width: `${(dollars / maxDollars) * 100}%`, minWidth: 2 }} />
            </div>
            <span className="w-[6.5rem] text-right font-mono text-xs font-semibold text-ink">{money(m.dollars)}</span>
          </>
        ) : (
          <span className="border border-dashed border-line px-2 font-mono text-[11px] text-ink-2">no amount</span>
        )}
      </div>
    </Link>
  )
}

// ---------------------------------------------------------------- page

type RankBy = 'claims' | 'dollars'

interface WorkloadTab {
  id: string
  type?: string
}

// Workloads from the measure rows, with their claim type when the rollup has
// type-by-workload rows.
function workloadsOf(measures: Measure[], claimType: string | null): WorkloadTab[] {
  const typed = measures.filter((m) => m.workload !== 'all' && m.claim_type !== 'all')
  const rows = typed.length ? typed : measures.filter((m) => m.workload !== 'all')
  const out: WorkloadTab[] = []
  for (const m of rows) {
    if (claimType && typed.length && m.claim_type !== claimType) continue
    if (!out.some((w) => w.id === m.workload)) out.push({ id: m.workload, type: m.claim_type === 'all' ? undefined : m.claim_type })
  }
  return out.sort((a, b) => a.id.localeCompare(b.id))
}

function measureFor(measures: Measure[], claimType: string | null, workload: string | null) {
  return measures.find((m) => m.claim_type === (claimType ?? 'all') && m.workload === (workload ?? 'all'))
    ?? measures.find((m) => m.claim_type === 'all' && m.workload === (workload ?? 'all'))
}

export function ViewDifferences() {
  const [params, setParams] = useSearchParams()
  const [rankBy, setRankBy] = useState<RankBy>('claims')
  const [showTail, setShowTail] = useState(false)

  const sel = useSelection()
  const submission = sel.shown
  const claimTypeParam = params.get('claim_type')
  const claimType = CLAIM_TYPES.find((t) => t === claimTypeParam)
  const workload = params.get('workload') || undefined

  const measures = useApi(submission ? `measures:${submission.prefix}` : null, (signal) => api.measures(submission!.prefix, signal))
  const differences = useApi(submission ? `differences:${submission.prefix}|${claimType ?? ''}|${workload ?? ''}` : null,
    (signal) => api.differences({ submission: submission!.prefix, claim_type: claimType, workload }, signal))

  const set = (next: Record<string, string | undefined>) => {
    const merged: Record<string, string | undefined> = { ...sel.carry, claim_type: claimType, workload, ...next }
    const search = new URLSearchParams()
    // '' (an explicit empty slot, sel.carry) is kept; undefined is dropped.
    for (const [k, v] of Object.entries(merged)) if (v !== undefined) search.set(k, v)
    setParams(search)
  }

  const measureRows = measures.data?.measures ?? []
  const workloads = workloadsOf(measureRows, claimType ?? null)
  const scope = measureFor(measureRows, claimType ?? null, workload ?? null)
  const rows = differences.data?.rows ?? []

  const sorted = [...rows].sort((x, y) => (rankBy === 'claims'
    ? y.claims - x.claims || compare(y.dollars, x.dollars)
    : compare(y.dollars, x.dollars) || y.claims - x.claims))
  const ranked = sorted.slice(0, TOP)
  const tail = sorted.slice(TOP)

  const pairs = scope ? scope.equal + scope.different : 0
  const different = differences.data?.different_claims ?? 0
  const equal = scope ? scope.equal : 0
  const equalShare = scope ? equalRate(scope) : undefined
  const paidRow = rows.find((r) => r.tag === 'CLP04')
  const maxClaims = Math.max(1, ...rows.map((m) => m.claims))
  const maxDollars = Math.max(1, ...rows.map((m) => Number(m.dollars)))
  const tailClaims = tail.reduce((s, m) => s + m.claims, 0)
  const tailDollars = sum(tail.map((m) => m.dollars))
  const groupLink = (tag: string) => sel.link('/claims', { tag, claim_type: claimType })

  const rowProps = (m: DifferenceRow) => ({ m, maxClaims, maxDollars, different, groupLink })

  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12">
      <PageTitle purpose={help.purpose}>Differences by kind</PageTitle>

      {submission && (
        <div className="mb-6">
          <SideToggle />
        </div>
      )}

      {submission && (
        <div className="mb-8 grid grid-cols-[6rem_1fr] gap-x-4 gap-y-3">
          <FilterRow label="Claim type">
            <Tabs<string> options={[{ id: 'all', label: 'All' }, ...CLAIM_TYPES.map((t) => ({ id: t, label: CLAIM_TYPE_LABEL[t] }))]} value={claimType ?? 'all'}
              onChange={(t) => set({ claim_type: t === 'all' ? undefined : t, workload: undefined })} />
          </FilterRow>
          <FilterRow label="Workload">
            <Tabs options={[{ id: 'all', label: 'All' }, ...workloads.map((w) => ({ id: w.id, label: w.type ? `${w.id} · ${PART_LABEL[w.type]}` : w.id }))]} value={workload ?? 'all'}
              onChange={(id) => set({ workload: id === 'all' ? undefined : id })} />
          </FilterRow>
        </div>
      )}

      {!differences.data ? (
        sel.unselected ? (
          <Empty>Select a submission in the navigation.</Empty>
        ) : (
          <Pending errors={[sel.error, measures.error, differences.error]} loading={sel.loading || differences.loading} empty="No vendor submission has been compared yet" />
        )
      ) : (
        <section className="border border-line">
          <div className="flex flex-wrap border-b border-line">
            <Stat label="Claim pairs compared" info={help.terms.pairsCompared} value={count(pairs)} />
            <Stat label="Exactly equal" info={help.terms.exactlyEqual} dot="bg-ok" value={equalShare === undefined ? '—' : `${(equalShare * 100).toFixed(1)}%`} sub={`${count(equal)} of ${count(scope ? scope.mac_claims - scope.set_aside : 0)} MAC claims`} />
            <Stat label="With any difference" info={help.terms.withAnyDifference} dot="bg-warn" value={count(different)} sub={percent(different, pairs)} />
            <Stat label="Kinds seen" info={help.terms.kindsSeen} value={count(rows.length)} />
            <Stat label="Medicare paid moved" info={help.terms.paidMoved} dot="bg-series-1" value={money(paidRow?.dollars ?? '0')} sub="CLP04, absolute, summed" />
          </div>

          <div className="border-b border-line px-5 py-5">
            <div className="mb-2 flex items-center gap-1 text-xs font-semibold text-ink">Claims with a difference, by kind<Info term={help.terms.concentration} /></div>
            <Concentration ranked={ranked} tail={tail} />
          </div>

          <div className="flex items-center gap-3 border-b border-line px-5 py-3">
            <span className="text-xs text-ink-2">Rank by</span>
            <Tabs<RankBy> options={[{ id: 'claims', label: 'Claims affected', dot: 'bg-warn' }, { id: 'dollars', label: 'Amount differing', dot: 'bg-series-1' }]} value={rankBy} onChange={setRankBy} />
          </div>

          <div className="overflow-x-auto px-2 pb-2">
            <div className="min-w-[58rem]">
              <div className={`${GRID} border-b-2 border-ink py-2 pr-3 text-xs text-ink`}>
                <span />
                <span className="flex items-center gap-1 font-semibold">What differs<Info term={help.terms.whatDiffers} /></span>
                <span className="flex items-center gap-1.5 font-semibold"><span className="h-2.5 w-2.5 bg-warn" />Claims affected<Info term={help.terms.claimsAffected} /></span>
                <span className="text-right font-semibold leading-tight">Share of different claims <Info term={help.terms.shareOfDifferent} /></span>
                <span className="flex items-center gap-1.5 font-semibold"><span className="h-2.5 w-2.5 bg-series-1" />Amount differing<Info term={help.terms.amountDiffering} /></span>
              </div>

              {rows.length === 0 && (
                <div className="flex items-center gap-3 bg-ok/10 px-5 py-6">
                  <span className="flex h-8 w-8 items-center justify-center rounded-full bg-ok font-bold text-white">✓</span>
                  <span className="font-serif text-xl text-ink">No differences</span>
                </div>
              )}

              {ranked.map((m, i) => <KindRow key={m.tag} rank={i + 1} {...rowProps(m)} />)}

              {tail.length > 0 && (
                <button type="button" onClick={() => setShowTail(!showTail)}
                  className={`${GRID} w-full border-b border-line bg-surface py-2.5 pr-3 text-left hover:text-ink`}>
                  <span className="text-right font-mono text-xs text-ink-2">{showTail ? '▾' : '▸'}</span>
                  <span className="text-sm text-ink">Long tail · {tail.length} more kinds</span>
                  <div className="flex items-center gap-2">
                    <div className="h-4 flex-1">
                      <div className="h-full bg-mark-off" style={{ width: `${(tailClaims / maxClaims) * 100}%`, minWidth: 2 }} />
                    </div>
                    <span className="w-[4.5rem] text-right font-mono text-xs font-semibold text-ink">{count(tailClaims)}</span>
                  </div>
                  <span className="text-right font-mono text-xs text-ink">{percent(tailClaims, different)}</span>
                  <span className="text-right font-mono text-xs font-semibold text-ink">{money(tailDollars)}</span>
                </button>
              )}

              {showTail && tail.map((m, i) => <KindRow key={m.tag} rank={TOP + i + 1} {...rowProps(m)} />)}
            </div>
          </div>

          <div className="flex flex-wrap gap-x-5 gap-y-1 border-t border-line px-5 py-3 text-xs text-ink-2">
            {(Object.keys(EFFECTS) as Effect[]).map((e) => (
              <span key={e} className="flex items-center gap-1.5"><span className={`h-2.5 w-2.5 ${EFFECTS[e].fill}`} />{EFFECTS[e].label}</span>
            ))}
            <Info term={help.terms.effects} />
            <span className="ml-auto font-mono text-[11px]">a claim with two kinds counts under each</span>
          </div>
        </section>
      )}

      <nav className="mt-8 flex flex-wrap gap-x-5 gap-y-2 text-sm">
        {[
          ['/', 'Overview'],
          ['/claims', 'One kind, its claims'],
          ['/claim', 'One claim, MAC beside both vendors'],
          ['/workloads', 'Workloads'],
          ['/unmatched', 'Unmatched claims'],
        ].map(([to, label]) => (
          <Link key={to} to={sel.link(to)} className="text-link hover:underline">{label} →</Link>
        ))}
      </nav>
    </div>
  )
}
