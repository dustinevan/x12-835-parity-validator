import { Link, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, useApi, type ClaimType, type Measure, type Tag } from '../../api'
import { CLAIM_TYPE_LABEL, PART_LABEL, count, equalRate } from '../../format'
import { HELP } from '../help'
import { Info } from '../Info'
import { CodeHint } from '../CodeHint'
import { SideToggle } from '../VendorPicker'
import { useSelection } from '../selection'

// By-workload summary for one vendor submission: outcome-state match per MAC
// workload as tiles sized by claim count, then ranked by difference rate.
// Reads /api/measures and /api/differences per workload.

const help = HELP.workloads

type TypeFilter = 'all' | ClaimType

interface Result {
  workload: string
  type?: ClaimType
  claims: number
  outcome?: number // fraction
  exact?: number // fraction
}

const rateOf = (m: Measure) => (m.outcome_match_rate === null ? undefined : Number(m.outcome_match_rate))

// One row per workload, with its claim type when the rollup has
// type-by-workload rows.
function results(measures: Measure[]): Result[] {
  const typed = measures.filter((m) => m.workload !== 'all' && m.claim_type !== 'all')
  const rows = typed.length ? typed : measures.filter((m) => m.workload !== 'all' && m.claim_type === 'all')
  return rows.map((m) => ({
    workload: m.workload,
    type: m.claim_type === 'all' ? undefined : m.claim_type,
    claims: m.mac_claims,
    outcome: rateOf(m),
    exact: equalRate(m),
  }))
}

// ---------------------------------------------------------------- formatting

const TARGET = 0.95
const pct = (r?: number) => (r === undefined ? '—' : `${(r * 100).toFixed(1)}%`)

// Four bins against the 95% target: two passing, two short.
const BINS = [
  { label: '≥ 98%', test: (r: number) => r >= 0.98, fill: 'bg-ok', ink: 'text-white' },
  { label: '95–98%', test: (r: number) => r >= TARGET, fill: 'bg-ok/45', ink: 'text-ink' },
  { label: '90–95%', test: (r: number) => r >= 0.9, fill: 'bg-warn/55', ink: 'text-ink' },
  { label: '< 90%', test: () => true, fill: 'bg-warn', ink: 'text-black' },
]
const bin = (r?: number) => (r === undefined ? { label: '', fill: 'bg-surface', ink: 'text-ink-2' } : BINS.find((b) => b.test(r))!)
const meets = (r: Result) => r.outcome !== undefined && r.outcome >= TARGET

function name(r: Result) {
  return { id: r.workload, part: r.type ? PART_LABEL[r.type] : '' }
}

// ---------------------------------------------------------------- tabs

function Tabs<T extends string>({ options, value, onChange, dot }: {
  options: { id: T; label: string }[]
  value: T
  onChange: (v: T) => void
  dot?: (v: T) => boolean
}) {
  return (
    <div className="inline-flex border border-line" role="tablist">
      {options.map((o) => (
        <button key={o.id} type="button" role="tab" aria-selected={value === o.id} onClick={() => onChange(o.id)}
          className={`flex items-center gap-2 border-r border-line px-4 py-2 text-sm last:border-r-0 ${value === o.id ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
          {dot && <span className={`h-2.5 w-2.5 rounded-full ${dot(o.id) ? 'bg-ok' : 'bg-warn'}`} />}
          {o.label}
        </button>
      ))}
    </div>
  )
}

// ---------------------------------------------------------------- tile map

function Tile({ r, max, top, link }: { r: Result; max: number; top?: Tag; link: string }) {
  const b = bin(r.outcome)
  const grow = r.claims / max
  const roomy = grow > 0.18
  const n = name(r)
  return (
    <Link to={link}
      title={`${n.id} ${n.part}: ${pct(r.outcome)} same outcome, ${count(r.claims)} claims`}
      style={{ flexGrow: r.claims, flexBasis: `${Math.max(7, grow * 22)}rem` }}
      className={`${b.fill} group relative flex h-28 min-w-[7rem] flex-col justify-between border-2 border-bg p-2.5 outline-offset-[-2px] hover:outline-2 hover:outline-ink`}>
      <div className="flex items-baseline justify-between gap-2">
        <span className={`font-mono text-sm font-semibold ${b.ink}`}>{n.id}</span>
        <span className={`text-[11px] ${b.ink} opacity-80`}>{n.part}</span>
      </div>
      <div>
        <div className={`font-mono text-xl font-semibold leading-none ${b.ink}`}>{pct(r.outcome)}</div>
        <div className={`mt-1.5 flex items-center justify-between gap-2 font-mono text-[11px] ${b.ink} opacity-85`}>
          <span>{count(r.claims)}</span>
        </div>
      </div>
      {roomy && !meets(r) && top && (
        <span className="absolute top-8 left-2.5 max-w-[90%] truncate text-[11px] text-ink">
          {top.label} <span className="font-mono text-ink-2"><CodeHint kind="element" code={top.tag} label={top.label}>{top.tag}</CodeHint></span>
        </span>
      )}
    </Link>
  )
}

function TileSection({ title, items, tops, link }: { title: string; items: Result[]; tops: Map<string, Tag>; link: (r: Result) => string }) {
  if (items.length === 0) return null
  const max = Math.max(1, ...items.map((r) => r.claims))
  const sorted = [...items].sort((a, b) => b.claims - a.claims)
  return (
    <div className="mb-6">
      <div className="mb-2 font-mono text-[11px] font-semibold uppercase tracking-wider text-ink-2">{title}</div>
      <div className="flex flex-wrap">
        {sorted.map((r) => <Tile key={`${r.workload}-${r.type}`} r={r} max={max} top={tops.get(r.workload)} link={link(r)} />)}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- ranked bars

// Ranked by difference rate, not by count: outcome-state mismatches per 1,000
// claims, with the target (50 per 1,000) marked.
function Ranked({ items, tops, link, groupLink }: { items: Result[]; tops: Map<string, Tag>; link: (r: Result) => string; groupLink: (tag: string) => string }) {
  const sorted = [...items].sort((a, b) => (a.outcome ?? 2) - (b.outcome ?? 2))
  const scale = Math.max(80, ...sorted.map((r) => (1 - (r.outcome ?? 1)) * 1000)) * 1.05
  const targetAt = (50 / scale) * 100
  return (
    <div>
      <div className="grid grid-cols-[8.5rem_1fr_6rem_12rem] items-center gap-x-3 border-b-2 border-ink pb-2 text-xs text-ink">
        <span className="flex items-center gap-1 pl-3">Workload<Info term={help.terms.workload} /></span>
        <span className="flex items-center gap-1">Different outcome, per 1,000 claims<Info term={help.terms.perThousand} /></span>
        <span className="flex items-center justify-end gap-1">Claims<Info term={help.terms.claims} /></span>
        <span className="flex items-center gap-1">Largest difference kind<Info term={help.terms.largestKind} /></span>
      </div>
      <div className="relative">
        <div className="pointer-events-none absolute inset-y-0 grid w-full grid-cols-[8.5rem_1fr_6rem_12rem] gap-x-3">
          <span />
          <div className="relative">
            <div className="absolute inset-y-0 border-l-2 border-dashed border-ink" style={{ left: `${targetAt}%` }}>
              <span className="absolute -top-0 left-1 whitespace-nowrap bg-bg px-1 font-mono text-[10px] text-ink">95% target</span>
            </div>
          </div>
        </div>
        {sorted.map((r) => {
          const per = r.outcome === undefined ? undefined : (1 - r.outcome) * 1000
          const miss = !meets(r)
          const top = tops.get(r.workload)
          const n = name(r)
          return (
            <div key={`${r.workload}-${r.type}`} className={`grid grid-cols-[8.5rem_1fr_6rem_12rem] items-center gap-x-3 border-b border-line py-1.5 text-sm ${miss ? 'bg-warn/15' : ''}`}>
              <Link to={link(r)}
                className={`border-l-4 pl-2 font-mono text-xs hover:underline ${miss ? 'border-warn font-semibold text-ink' : 'border-transparent text-ink-2'}`}>
                {n.id} <span className="font-sans text-ink-2">{n.part}</span>
              </Link>
              <div className="flex items-center gap-2">
                <div className={`h-4 ${miss ? 'bg-warn' : 'bg-ok'}`} style={{ width: `${((per ?? 0) / scale) * 100}%` }} />
                <span className="font-mono text-xs text-ink">{per === undefined ? '—' : per.toFixed(0)}</span>
              </div>
              <span className="text-right font-mono text-xs text-ink-2">{count(r.claims)}</span>
              {top ? (
                <Link to={groupLink(top.tag)} className="truncate text-xs text-ink-2 hover:text-link hover:underline">
                  {top.label} <span className="font-mono"><CodeHint kind="element" code={top.tag} label={top.label}>{top.tag}</CodeHint></span>
                </Link>
              ) : <span className="font-mono text-[11px] text-ink-2">—</span>}
            </div>
          )
        })}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- page

const TYPE_TABS: { id: TypeFilter; label: string }[] = [
  { id: 'all', label: 'All claim types' },
  { id: 'institutional', label: CLAIM_TYPE_LABEL.institutional },
  { id: 'professional', label: CLAIM_TYPE_LABEL.professional },
  { id: 'dme', label: CLAIM_TYPE_LABEL.dme },
]

const OTHER_VIEWS = [
  ['/', 'Overview'],
  ['/differences', 'Differences'],
  ['/unmatched', 'Unmatched'],
  ['/claim', 'One claim'],
]

export function ViewWorkloads() {
  const [params, setParams] = useSearchParams()
  const sel = useSelection()
  const submission = sel.shown
  const type = (TYPE_TABS.find((t) => t.id === params.get('claim_type'))?.id ?? 'all') as TypeFilter

  const measures = useApi(submission ? `measures:${submission.prefix}` : null, (signal) => api.measures(submission!.prefix, signal))
  const every = measures.data ? results(measures.data.measures) : []

  // The largest difference kind in each workload.
  const workloadIds = [...new Set(every.map((r) => r.workload))].sort()
  const tops = useApi(submission && workloadIds.length ? `tops:${submission.prefix}|${workloadIds.join(',')}` : null,
    (signal) => Promise.all(workloadIds.map((w) => api.differences({ submission: submission!.prefix, workload: w }, signal).then((d) => [w, d.rows[0]] as const))))
  const topMap = new Map<string, Tag>()
  for (const [w, row] of tops.data ?? []) if (row) topMap.set(w, row)

  const set = (next: Record<string, string | undefined>) => {
    const merged: Record<string, string | undefined> = { ...sel.carry, claim_type: type === 'all' ? undefined : type, ...next }
    const search = new URLSearchParams()
    // '' (an explicit empty slot, sel.carry) is kept; undefined is dropped.
    for (const [k, v] of Object.entries(merged)) if (v !== undefined) search.set(k, v)
    setParams(search)
  }

  const shown = every.filter((r) => type === 'all' || r.type === type)
  const claims = shown.reduce((s, r) => s + r.claims, 0)
  const weighted = (pick: (r: Result) => number | undefined) => {
    let total = 0
    let weight = 0
    for (const r of shown) {
      const v = pick(r)
      if (v === undefined) continue
      total += v * r.claims
      weight += r.claims
    }
    return weight === 0 ? undefined : total / weight
  }
  const scope = measures.data?.measures.find((m) => m.claim_type === type && m.workload === 'all')
  const outcome = scope ? rateOf(scope) : weighted((r) => r.outcome)
  const exact = scope ? equalRate(scope) : weighted((r) => r.exact)
  const meeting = shown.filter(meets).length
  const typeMeets = (t: TypeFilter) => every.filter((r) => t === 'all' || r.type === t).every(meets)

  const link = (r: Result) => sel.link('/differences', { claim_type: r.type, workload: r.workload })
  const groupLink = (tag: string) => sel.link('/claims', { tag })

  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12">
      <PageTitle purpose={help.purpose}>By workload</PageTitle>

      <nav className="-mt-6 mb-8 flex flex-wrap gap-x-4 gap-y-1 text-xs">
        {OTHER_VIEWS.map(([to, label]) => (
          <Link key={to} to={sel.link(to)} className="text-link hover:underline">{label}</Link>
        ))}
      </nav>

      {submission && (
        <>
          <div className="mb-4">
            <SideToggle />
          </div>
          <div className="mb-8">
            <Tabs options={TYPE_TABS} value={type} onChange={(t) => set({ claim_type: t === 'all' ? undefined : t })} dot={typeMeets} />
          </div>
        </>
      )}

      {!measures.data || !submission ? (
        sel.unselected ? (
          <Empty>Select a submission in the navigation.</Empty>
        ) : (
          <Pending errors={[sel.error, measures.error]} loading={sel.loading || measures.loading} empty="No vendor submission has been compared yet" />
        )
      ) : (
        <section className="border border-line">
          <div className="flex flex-wrap items-stretch border-b border-line">
            {([
              ['Workloads at or above 95%', `${meeting} of ${shown.length}`, help.terms.workloadsMeeting],
              ['Same outcome state', pct(outcome), help.terms.outcomeMatch],
              ['No difference at all', pct(exact), help.terms.exactlyEqual],
              ['MAC claims', count(claims), help.terms.macClaims],
            ] as const).map(([label, value, term]) => (
              <div key={label} className="border-r border-line px-5 py-3 last:border-r-0">
                <div className="flex items-center gap-1 text-[11px] text-ink-2">{label}<Info term={term} /></div>
                <div className="font-mono text-lg text-ink">{value}</div>
              </div>
            ))}
            <div className={`ml-auto flex items-center gap-3 px-5 ${meeting === shown.length ? 'bg-ok/10' : 'bg-warn/15'}`}>
              <span className={`flex h-9 w-9 items-center justify-center rounded-full text-base font-bold ${meeting === shown.length ? 'bg-ok text-white' : 'bg-warn text-black'}`}>
                {meeting === shown.length ? '✓' : shown.length - meeting}
              </span>
              <span className="font-serif text-xl text-ink">{meeting === shown.length ? 'All meet target' : 'below target'}</span>
            </div>
          </div>

          <div className="px-5 pt-5 pb-1">
            <div className="mb-4 flex flex-wrap items-center gap-x-5 gap-y-1 text-xs text-ink-2">
              <span className="font-semibold text-ink">Same outcome state</span>
              {BINS.map((b) => (
                <span key={b.label} className="flex items-center gap-1.5">
                  <span className={`h-3 w-3 ${b.fill}`} />
                  <span className="font-mono">{b.label}</span>
                </span>
              ))}
              <span className="ml-auto flex items-center gap-1">tile width = MAC claims<Info term={help.terms.tiles} /></span>
            </div>
            {type !== 'dme' && (
              <TileSection title="A/B workloads" items={shown.filter((r) => r.type !== 'dme')} tops={topMap} link={link} />
            )}
            {(type === 'all' || type === 'dme') && (
              <TileSection title="DME workloads" items={shown.filter((r) => r.type === 'dme')} tops={topMap} link={link} />
            )}
          </div>

          <div className="border-t border-line px-2 pt-5 pb-4">
            <div className="mb-3 px-3 text-xs font-semibold text-ink">Ranked by difference rate</div>
            <Ranked items={shown} tops={topMap} link={link} groupLink={groupLink} />
          </div>

          <div className="border-t border-line px-5 py-3 font-mono text-[11px] text-ink-2">
            {submission.label} · submission {submission.submitted_on} · click a workload for its differences
          </div>
        </section>
      )}
    </div>
  )
}
