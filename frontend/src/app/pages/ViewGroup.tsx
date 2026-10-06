import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, nameOf, sharesVendor, useApi, type Bucket, type ClaimType, type Compared, type Difference, type DifferenceRow, type ClaimRow, type Sort, type TagFilter } from '../../api'
import { CLAIM_TYPE_LABEL, compare, count, isZero, money, signedMoney } from '../../format'
import { claimLink } from '../links'
import { HELP, type Term } from '../help'
import { Info } from '../Info'
import { CodeHint } from '../CodeHint'
import { SideToggle, SlotChip } from '../VendorPicker'
import { useSelection, type Side } from '../selection'

// Claims in one vendor submission chosen by tag: Must have / Any of / Must not
// have, kept in the URL (all=, any=, none=). One Must have tag shows that kind
// on a real example; every selection gets the paid-difference histogram and a
// paged table sorted by the server. An Overview agreement cell adds each
// slot's result (v1_outcome=, v2_outcome=). Reads /api/claims, /api/claim,
// /api/differences.

const help = HELP.claims

const CLAIM_TYPES: ClaimType[] = ['professional', 'institutional', 'dme']
const SORTS: Sort[] = ['delta', '-delta', 'paid', 'charge']
const COMPARED: Compared[] = ['equal', 'different', 'missing', 'duplicated']
const SIDES: Side[] = ['v1', 'v2']

// ---------------------------------------------------------------- histogram

function bucketLabel(b: Bucket) {
  if (compare(b.from, b.to) === 0) return signedMoney(b.from)
  return `${signedMoney(b.from)} – ${signedMoney(b.to)}`
}

const zeroBucket = (b: Bucket) => isZero(b.from) && isZero(b.to)

function Histogram({ buckets }: { buckets: Bucket[] }) {
  const max = Math.max(1, ...buckets.map((b) => b.count))
  return (
    <div>
      <div className="flex h-44 items-end gap-2">
        {buckets.map((b) => (
          <div key={`${b.from}|${b.to}`} title={`${bucketLabel(b)}: ${count(b.count)} claims`}
            className="flex h-full flex-1 flex-col items-center justify-end">
            <span className="mb-1 font-mono text-[11px] text-ink-2">{count(b.count)}</span>
            <span style={{ height: `${(b.count / max) * 100}%`, minHeight: b.count > 0 ? 3 : 0 }}
              className={`w-full ${zeroBucket(b) ? 'bg-mark-off' : 'bg-warn'}`} />
          </div>
        ))}
      </div>
      <div className="mt-2 flex gap-2 border-t border-ink pt-1.5">
        {buckets.map((b) => (
          <span key={`${b.from}|${b.to}`} className="flex-1 text-center font-mono text-[11px] text-ink-2">{bucketLabel(b)}</span>
        ))}
      </div>
      <div className="mt-2 flex flex-wrap gap-x-5 gap-y-1 text-xs text-ink-2">
        <span className="flex items-center gap-1.5"><span className="h-2.5 w-2.5 bg-warn" />Paid differs</span>
        <span className="flex items-center gap-1.5"><span className="h-2.5 w-2.5 bg-mark-off" />Same paid</span>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- the kind, in one picture

function Value({ value }: { value: string | null }) {
  if (value === null) return <span className="border border-dashed border-line px-2 py-0.5 font-mono text-xs text-ink-2">absent</span>
  return <span className="bg-series-1 px-2 py-1 font-mono text-[11px] font-semibold text-white">{value}</span>
}

// The first example claim's values for this tag, MAC above vendor.
function KindPicture({ differences, tag, label }: { differences: Difference[]; tag: string; label?: string }) {
  const d = differences.find((x) => x.element === tag)
  if (!d) return null
  return (
    <div className="grid grid-cols-[4.5rem_1fr] items-center gap-x-3 gap-y-2 text-sm">
      <span className="text-right font-mono text-xs text-ink-2">MAC</span>
      <div className="flex items-center gap-2">
        <span className="font-mono text-[11px] text-ink-2"><CodeHint kind="element" code={tag} label={label}>{tag}</CodeHint></span>
        <Value value={d.mac_value} />
      </div>
      <span className="text-right font-mono text-xs text-ink-2">Vendor</span>
      <div className="flex items-center gap-2">
        <span className="font-mono text-[11px] text-ink-2"><CodeHint kind="element" code={tag} label={label}>{tag}</CodeHint></span>
        <Value value={d.vendor_value} />
        {d.amount_difference !== null && !isZero(d.amount_difference) && (
          <span className="bg-warn px-1.5 font-mono text-[11px] font-semibold text-black">{signedMoney(d.amount_difference)}</span>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- table

const COLUMNS: { label: string; info?: Term; numeric?: boolean; sort?: [Sort, Sort?] }[] = [
  { label: 'Patient ID', info: help.terms.patientId },
  { label: 'Service date', info: help.terms.serviceDate },
  { label: 'Billed', info: help.terms.billed, numeric: true, sort: ['charge'] },
  { label: 'MAC paid', numeric: true, sort: ['paid'] },
  { label: 'Vendor paid', numeric: true },
  { label: 'Difference', info: help.terms.difference, numeric: true, sort: ['-delta', 'delta'] },
  { label: 'Claim type', info: help.terms.claimType },
  { label: 'Workload', info: help.terms.workload },
]

function ClaimTable({ rows, carry, sort, onSort, page, pages, total, onPage }: {
  rows: ClaimRow[]
  carry: Record<string, string | undefined>
  sort: Sort
  onSort: (s: Sort) => void
  page: number
  pages: number
  total: number
  onPage: (p: number) => void
}) {
  const navigate = useNavigate()
  const PAGE = 25
  const first = (page - 1) * PAGE + 1
  return (
    <div>
      <div className="overflow-x-auto">
        <table className="w-full min-w-[52rem] text-sm">
          <thead>
            <tr className="border-b-2 border-ink text-left text-xs text-ink">
              {COLUMNS.map((c) => {
                const active = c.sort?.includes(sort)
                return (
                  <th key={c.label} className={`py-2 ${c.numeric ? 'pr-4 text-right' : 'pl-3'}`}>
                    {c.sort ? (
                      <button type="button" onClick={() => onSort(active && c.sort![1] && sort === c.sort![0] ? c.sort![1] : c.sort![0])}
                        className="inline-flex items-center gap-1 font-semibold hover:text-link">
                        {c.label}
                        <span className={`text-[10px] ${active ? 'text-ink' : 'text-line'}`}>{sort === 'delta' && active ? '▲' : '▼'}</span>
                      </button>
                    ) : <span className="font-semibold">{c.label}</span>}
                    {c.info && <> <Info term={c.info} /></>}
                  </th>
                )
              })}
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => {
              const differs = r.delta !== null && !isZero(r.delta)
              return (
                <tr key={r.match_key} onClick={() => navigate(claimLink(carry, r.match_key))} className="cursor-pointer border-b border-line hover:bg-surface">
                  <td className="py-2 pl-3 font-mono"><Link to={claimLink(carry, r.match_key)} className="text-link hover:underline" onClick={(e) => e.stopPropagation()}>{r.key.mbi}</Link></td>
                  <td className="py-2 pl-3 font-mono text-ink-2">{r.key.date_of_service}</td>
                  <td className="py-2 pr-4 text-right font-mono text-ink-2">{money(r.key.total_charge)}</td>
                  <td className="py-2 pr-4 text-right font-mono text-ink-2">{r.mac_paid === null ? '—' : money(r.mac_paid)}</td>
                  <td className={`py-2 pr-4 text-right font-mono ${differs ? 'font-semibold text-ink' : 'text-ink-2'}`}>{r.vendor_paid === null ? '—' : money(r.vendor_paid)}</td>
                  <td className="py-2 pr-4 text-right">
                    {differs
                      ? <span className="bg-warn px-1.5 font-mono text-[11px] font-semibold text-black">{signedMoney(r.delta ?? '0')}</span>
                      : <span className="font-mono text-xs text-ink-2">$0.00</span>}
                  </td>
                  <td className="py-2 pl-3 text-ink-2">{CLAIM_TYPE_LABEL[r.claim_type] ?? r.claim_type}</td>
                  <td className="py-2 pl-3 font-mono text-ink-2">{r.key.workload}</td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
      <div className="mt-3 flex items-center gap-3 text-xs text-ink-2">
        <span className="font-mono">{total === 0 ? '0' : `${count(first)}–${count(Math.min(total, first + rows.length - 1))}`} of {count(total)}</span>
        <div className="ml-auto inline-flex border border-line">
          <button type="button" disabled={page <= 1} onClick={() => onPage(1)} className="px-3 py-1.5 hover:text-ink disabled:opacity-40">First</button>
          <button type="button" disabled={page <= 1} onClick={() => onPage(page - 1)} className="border-l border-line px-3 py-1.5 hover:text-ink disabled:opacity-40">Previous</button>
          <span className="border-l border-line bg-bg px-3 py-1.5 font-mono font-semibold text-ink">{page} / {count(pages)}</span>
          <button type="button" disabled={page >= pages} onClick={() => onPage(page + 1)} className="border-l border-line px-3 py-1.5 hover:text-ink disabled:opacity-40">Next</button>
          <button type="button" disabled={page >= pages} onClick={() => onPage(pages)} className="border-l border-line px-3 py-1.5 hover:text-ink disabled:opacity-40">Last</button>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- combination filter

type Box = keyof TagFilter

const BOXES: { box: Box; label: string; rule: string }[] = [
  { box: 'all', label: 'Must have', rule: 'every one' },
  { box: 'any', label: 'Any of', rule: 'at least one' },
  { box: 'none', label: 'Must not have', rule: 'none of these' },
]

// How a tag looks in each box: solid amber, amber outline, struck through.
const CHIP: Record<Box, string> = {
  all: 'bg-warn text-black',
  any: 'border border-warn bg-warn/15 text-ink',
  none: 'border border-dashed border-ink-2 text-ink-2 line-through',
}

function boxOf(filter: TagFilter, tag: string): Box | undefined {
  return BOXES.find((b) => filter[b.box].includes(tag))?.box
}

// Three boxes the experts fill by clicking tags: pick a box, then click tags
// below to add them to it. Clicking a placed tag, in its box or in the list,
// takes it out.
function TagFilterPanel({ filter, tags, onChange }: {
  filter: TagFilter
  tags: DifferenceRow[]
  onChange: (f: TagFilter) => void
}) {
  const [target, setTarget] = useState<Box>('all')
  const labels = new Map(tags.map((t) => [t.tag, t.label]))
  const label = (tag: string) => labels.get(tag) ?? tag

  const without = (tag: string): TagFilter => ({
    all: filter.all.filter((t) => t !== tag),
    any: filter.any.filter((t) => t !== tag),
    none: filter.none.filter((t) => t !== tag),
  })
  const toggle = (tag: string) => {
    if (boxOf(filter, tag)) { onChange(without(tag)); return }
    const next = without(tag)
    next[target] = [...next[target], tag]
    onChange(next)
  }
  const empty = filter.all.length + filter.any.length + filter.none.length === 0

  return (
    <div className="border-b border-line px-5 py-4">
      <div className="grid gap-2 sm:grid-cols-3">
        {BOXES.map((b) => {
          const selected = target === b.box
          return (
            <div key={b.box} role="button" tabIndex={0} aria-pressed={selected}
              onClick={() => setTarget(b.box)} onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') setTarget(b.box) }}
              className={`min-h-[4.5rem] cursor-pointer border px-3 py-2 ${selected ? 'border-ink bg-bg' : 'border-line bg-surface hover:border-ink-2'}`}>
              <div className="mb-1.5 flex items-baseline gap-2">
                <span className={`text-xs font-semibold ${selected ? 'text-ink' : 'text-ink-2'}`}>{b.label}</span>
                <span className="text-[11px] text-ink-2">{b.rule}</span>
              </div>
              <div className="flex flex-wrap gap-1.5">
                {filter[b.box].map((tag) => (
                  <button key={tag} type="button" title={`${tag} — click to remove`}
                    onClick={(e) => { e.stopPropagation(); onChange(without(tag)) }}
                    className={`inline-flex items-center gap-1.5 px-2 py-0.5 text-xs ${CHIP[b.box]}`}>
                    <CodeHint kind="element" code={tag} label={label(tag)}>{label(tag)}</CodeHint><span aria-hidden className="no-underline opacity-60">×</span>
                  </button>
                ))}
              </div>
            </div>
          )
        })}
      </div>
      <div className="mt-3 flex flex-wrap items-center gap-1.5">
        <Info term={help.terms.filter} />
        {tags.map((t) => {
          const placed = boxOf(filter, t.tag)
          return (
            <button key={t.tag} type="button" title={`${t.tag} · ${count(t.claims)} claims`} onClick={() => toggle(t.tag)}
              className={`inline-flex items-center gap-1.5 px-2 py-0.5 text-xs ${placed ? CHIP[placed] : 'border border-line bg-surface text-ink-2 hover:border-ink-2 hover:text-ink'}`}>
              <CodeHint kind="element" code={t.tag} label={t.label}>{t.label}</CodeHint>
              <span className="font-mono text-[10px] opacity-70">{count(t.claims)}</span>
            </button>
          )
        })}
        {!empty && (
          <button type="button" onClick={() => onChange({ all: [], any: [], none: [] })} className="ml-2 text-xs text-link hover:underline">Clear</button>
        )}
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- page

const list = (value: string | null) => (value ? value.split(',').filter(Boolean) : [])

export function ViewGroup() {
  const [params, setParams] = useSearchParams()
  const sel = useSelection()
  const submission = sel.shown
  const claimType = CLAIM_TYPES.find((t) => t === params.get('claim_type'))
  const sort = SORTS.find((s) => s === params.get('sort')) ?? '-delta'
  const page = Math.max(1, Number(params.get('page')) || 1)

  // Each slot's result for the claims (v1_outcome=, v2_outcome=), from an
  // Overview agreement cell. The shown side's result filters this submission;
  // the other side's filters by the other picked submission, so the toggle
  // shows the same claims for either slot.
  const outcomes: Record<Side, Compared | undefined> = {
    v1: COMPARED.find((o) => o === params.get('v1_outcome')),
    v2: COMPARED.find((o) => o === params.get('v2_outcome')),
  }
  const byOutcome = outcomes.v1 !== undefined || outcomes.v2 !== undefined
  const otherSide: Side = sel.side === 'v2' ? 'v1' : 'v2'
  const other = otherSide === 'v1' ? sel.v1 : sel.v2
  const otherOutcome = other ? outcomes[otherSide] : undefined

  // Every tag in the submission, for the chip list and plain labels.
  const tags = useApi(submission ? `differences:${submission.prefix}` : null, (signal) => api.differences({ submission: submission!.prefix }, signal))

  // The filter lives in the URL. Arriving with ?tag= (from the differences
  // view) puts that tag in Must have. With no filter keys at all, open the
  // submission's most common kind; a cleared filter keeps an empty `all=` so
  // it stays cleared, and a result filter opens with no tags.
  const hasFilter = params.has('all') || params.has('any') || params.has('none')
  const tagParam = params.get('tag')
  const fallback = tagParam ?? (hasFilter || byOutcome ? undefined : tags.data?.rows[0]?.tag)
  const filter: TagFilter = {
    all: hasFilter ? list(params.get('all')) : fallback ? [fallback] : [],
    any: list(params.get('any')),
    none: list(params.get('none')),
  }
  const ready = hasFilter || byOutcome || tagParam !== null || tags.data !== undefined
  const filterKey = `${filter.all.join(',')}|${filter.any.join(',')}|${filter.none.join(',')}`
  const outcomeKey = `${outcomes[sel.side] ?? ''}|${otherOutcome ? `${other!.prefix}=${otherOutcome}` : ''}`

  const claims = useApi(submission && ready ? `claims:${submission.prefix}|${filterKey}|${outcomeKey}|${claimType ?? ''}|${sort}|${page}` : null,
    (signal) => api.claims({
      submission: submission!.prefix, ...filter, claim_type: claimType, outcome: outcomes[sel.side],
      other: otherOutcome ? other!.prefix : undefined, other_outcome: otherOutcome, sort, page,
    }, signal))

  // One example claim for the picture of the first Must have tag: every
  // claim in the table carries it.
  const pictured = filter.all[0]
  const exampleKey = pictured ? claims.data?.rows[0]?.match_key : undefined
  const example = useApi(submission && exampleKey ? `claim:${submission.prefix}|${exampleKey}` : null,
    (signal) => api.claim({ submission: submission!.prefix, match_key: exampleKey! }, signal))

  const set = (next: Record<string, string | undefined>, nextFilter: TagFilter = filter) => {
    const search = new URLSearchParams()
    // '' (an explicit empty slot, sel.carry) is kept; undefined is dropped.
    for (const [k, v] of Object.entries(sel.carry)) if (v !== undefined) search.set(k, v)
    for (const box of ['all', 'any', 'none'] as const) if (nextFilter[box].length) search.set(box, nextFilter[box].join(','))
    if (!search.has('all') && !search.has('any') && !search.has('none')) search.set('all', '')
    const plain: Record<string, string | undefined> = { claim_type: claimType, v1_outcome: outcomes.v1, v2_outcome: outcomes.v2, sort, page: String(page), ...next }
    for (const [k, v] of Object.entries(plain)) if (v) search.set(k, v)
    setParams(search)
  }

  // One Must have tag and nothing else: the page is about that one kind.
  const single = filter.all.length === 1 && filter.any.length === 0 && filter.none.length === 0
  const singleTag = single ? tags.data?.rows.find((r) => r.tag === filter.all[0]) : undefined
  const data = claims.data
  const noTags = filter.all.length + filter.any.length + filter.none.length === 0
  const picked = { v1: sel.v1, v2: sel.v2 }
  // Two submissions of one vendor are told apart by date.
  const qualified = sharesVendor(picked.v1, picked.v2)
  const outcomeTitle = SIDES.flatMap((s) => (outcomes[s] && picked[s] ? [`${nameOf(picked[s]!, qualified)} ${outcomes[s]}`] : [])).join(', ')
  const title = single ? (singleTag?.label ?? filter.all[0]) : byOutcome && noTags ? outcomeTitle : 'Claims by tag combination'
  const histogram = data?.histogram ?? []
  const same = histogram.filter(zeroBucket).reduce((s, b) => s + b.count, 0)
  const moved = histogram.reduce((s, b) => s + b.count, 0) - same

  const stats = data ? [
    { label: 'Claims', info: help.terms.claims, value: count(data.claims) },
    { label: 'Paid differs', info: help.terms.paidDiffers, value: count(moved) },
    { label: 'Same paid', info: help.terms.samePaid, value: count(same) },
    { label: 'Paid difference, total', info: help.terms.paidTotal, value: money(data.dollars) },
  ] : []

  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12">
      <nav className="mb-6 flex flex-wrap items-center gap-2 text-xs text-ink-2">
        <Link to={sel.link('/')} className="text-link hover:underline">Overview</Link>
        <span>/</span>
        <Link to={sel.link('/differences', { claim_type: claimType })} className="text-link hover:underline">Differences</Link>
        <span>/</span>
        <span className="text-ink">{title}</span>
      </nav>

      <PageTitle purpose={help.purpose}>{title}</PageTitle>

      {submission && (
        <div className="mb-8">
          <SideToggle />
        </div>
      )}

      {!submission ? (
        sel.unselected ? (
          <Empty>Select a submission in the navigation.</Empty>
        ) : (
          <Pending errors={[sel.error]} loading={sel.loading} empty="No vendor submissions" />
        )
      ) : (
        <section className="border border-line">
          <div className="flex flex-wrap items-stretch border-b border-line">
            {singleTag && (
              <div className="border-r border-line px-4 py-3">
                <div className="text-[11px] text-ink-2">Element</div>
                <div className="font-mono text-sm font-semibold text-ink"><CodeHint kind="element" code={singleTag.tag} label={singleTag.label}>{singleTag.tag}</CodeHint></div>
              </div>
            )}
            <div className="border-r border-line px-4 py-3">
              <div className="text-[11px] text-ink-2">Vendor</div>
              <div className="font-mono text-sm text-ink">{submission.label}</div>
            </div>
            <div className="border-r border-line px-4 py-3">
              <div className="text-[11px] text-ink-2">Submission</div>
              <div className="font-mono text-sm text-ink">{submission.submitted_on}</div>
            </div>
            {byOutcome && (
              <div className="flex items-center gap-4 border-r border-line px-4 py-3">
                <div>
                  <div className="flex items-center gap-1 text-[11px] text-ink-2">Result<Info term={help.terms.byOutcome} /></div>
                  <div className="flex items-center gap-3 text-sm text-ink">
                    {SIDES.map((s) => outcomes[s] && (
                      <span key={s} className="flex items-center gap-1.5"><SlotChip side={s} />{outcomes[s]}</span>
                    ))}
                  </div>
                </div>
                <button type="button" onClick={() => set({ v1_outcome: undefined, v2_outcome: undefined, page: undefined })} className="text-xs text-link hover:underline">Clear</button>
              </div>
            )}
          </div>

          <TagFilterPanel filter={filter} tags={tags.data?.rows ?? []} onChange={(f) => set({ page: undefined }, f)} />

          {!data ? (
            <div className="px-5 py-6">
              <Pending errors={[tags.error, claims.error]} loading={tags.loading || claims.loading} empty="No difference kinds in this submission" />
            </div>
          ) : (
            <>
              <div className="flex flex-wrap items-center gap-x-10 gap-y-5 bg-warn/15 px-5 py-5">
                {pictured && example.data && (
                  <KindPicture differences={example.data.differences} tag={pictured} label={tags.data?.rows.find((r) => r.tag === pictured)?.label} />
                )}
                <div className="ml-auto flex flex-wrap gap-x-8 gap-y-3">
                  {stats.map((s) => (
                    <div key={s.label} className="text-right">
                      <div className="flex items-center justify-end gap-1 text-[11px] text-ink-2">{s.label}<Info term={s.info} /></div>
                      <div className="font-mono text-xl font-semibold text-ink">{s.value}</div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="border-t border-line px-5 py-6">
                <div className="mb-4 flex items-center gap-1 text-xs font-semibold text-ink">Difference in Medicare paid, per claim<Info term={help.terms.histogram} /></div>
                <Histogram buckets={histogram} />
              </div>

              <div className="border-t border-line px-2 pt-4 pb-4">
                <div className="mb-3 flex flex-wrap items-center gap-3 px-3">
                  <div className="inline-flex border border-line" role="tablist">
                    {([undefined, ...CLAIM_TYPES] as (ClaimType | undefined)[]).map((t) => (
                      <button key={t ?? 'all'} type="button" role="tab" aria-selected={claimType === t} onClick={() => set({ claim_type: t, page: undefined })}
                        className={`flex items-center gap-2 border-r border-line px-4 py-2 text-sm last:border-r-0 ${claimType === t ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
                        {t ? CLAIM_TYPE_LABEL[t] : 'All'}
                        {claimType === t && <span className="font-mono text-[11px] font-normal text-ink-2">{count(data.claims)}</span>}
                      </button>
                    ))}
                  </div>
                  <Link to={sel.link('/workloads')} className="ml-auto text-xs text-link hover:underline">By workload →</Link>
                </div>
                <ClaimTable rows={data.rows} carry={sel.carry} sort={sort} onSort={(s) => set({ sort: s, page: undefined })}
                  page={data.page} pages={Math.max(1, data.pages)} total={data.claims} onPage={(p) => set({ page: String(p) })} />
              </div>
            </>
          )}
        </section>
      )}
    </div>
  )
}
