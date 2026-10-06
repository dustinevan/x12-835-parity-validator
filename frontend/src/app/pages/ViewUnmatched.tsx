import { Link, useNavigate, useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, useApi, type Key, type Unmatched, type UnmatchedResponse, type UnmatchedRow } from '../../api'
import { CLAIM_TYPE_LABEL, count, money } from '../../format'
import { claimLink } from '../links'
import { HELP } from '../help'
import { Info } from '../Info'
import { SideToggle } from '../VendorPicker'
import { useSelection } from '../selection'

// Claims on one side only (missing, extra), claims the vendor sent more than
// once (duplicated), and claims set aside because several MAC claims share
// one match key, for one vendor submission. Reads /api/unmatched and
// /api/measures.

const help = HELP.unmatched

const TABS: { id: Unmatched; label: string; side: string }[] = [
  { id: 'missing', label: 'Missing', side: 'MAC claims with no vendor claim' },
  { id: 'extra', label: 'Extra', side: 'Vendor claims with no MAC claim' },
  { id: 'duplicated', label: 'Duplicated', side: 'MAC claims the vendor sent more than once' },
  { id: 'set_aside', label: 'Set aside', side: 'MAC claims sharing a key' },
]

// ---------------------------------------------------------------- format

const pct = (fraction: number) => `${(fraction * 100).toFixed(2)}%`
const perThousand = (fraction: number) => `${(fraction * 1000).toFixed(1)} per 1,000`

// ---------------------------------------------------------------- pieces

function Absent() {
  return <span className="inline-block border border-dashed border-line px-2 font-mono text-xs text-ink-2">absent</span>
}

// The vendor side of a duplicated key: it has the claim, more than once.
function Repeated() {
  return <span className="inline-block border border-dashed border-warn px-2 font-mono text-xs text-ink">sent more than once</span>
}

function Headline({ tab, data, of }: { tab: (typeof TABS)[number]; data: UnmatchedResponse; of?: number }) {
  const share = Number(data.share)
  const keys = tab.id === 'set_aside' ? new Set(data.rows.map((r) => r.match_key)).size : undefined
  return (
    <div className="flex flex-wrap items-center gap-x-8 gap-y-3 bg-warn/15 px-5 py-5">
      <span className="flex h-10 min-w-10 items-center justify-center rounded-full bg-warn px-2 text-lg font-bold text-black">{count(data.count)}</span>
      <div>
        <div className="flex items-center gap-2 font-serif text-2xl text-ink">{count(data.count)} {tab.id === 'set_aside' ? 'claims set aside' : `${tab.label.toLowerCase()} claims`}<Info term={help.terms[tab.id]} /></div>
        <div className="text-xs text-ink-2">{tab.side}{keys !== undefined && ` · ${count(keys)} keys shown`}</div>
      </div>
      <div className="ml-auto flex gap-8 text-right">
        <div>
          <div className="flex items-center justify-end gap-1 text-[11px] text-ink-2">Share<Info term={help.terms.share} /></div>
          <div className="font-mono text-xl font-semibold text-ink">{pct(share)}</div>
        </div>
        <div>
          <div className="text-[11px] text-ink-2">Rate</div>
          <div className="font-mono text-xl font-semibold text-ink">{perThousand(share)}</div>
        </div>
        {of !== undefined && (
          <div>
            <div className="text-[11px] text-ink-2">MAC claims</div>
            <div className="font-mono text-xl text-ink-2">{count(of)}</div>
          </div>
        )}
      </div>
    </div>
  )
}

function Bars({ title, rows }: { title: string; rows: [string, number][] }) {
  const max = Math.max(...rows.map(([, n]) => n), 1)
  return (
    <div>
      <div className="mb-3 text-xs font-semibold text-ink">{title}</div>
      <div className="grid grid-cols-[6.5rem_1fr_2.5rem] items-center gap-x-3 gap-y-1.5">
        {rows.map(([label, n]) => (
          <div key={label} className="contents">
            <span className="truncate text-right font-mono text-xs text-ink-2">{label}</span>
            <div className="h-4 bg-surface">
              {n > 0 && <div className="h-full bg-warn" style={{ width: `${(n / max) * 100}%`, minWidth: 3 }} />}
            </div>
            <span className="font-mono text-xs text-ink">{count(n)}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

function BreakdownPanel({ data }: { data: UnmatchedResponse }) {
  const { link } = useSelection()
  const types = data.by_claim_type.map((t): [string, number] => [CLAIM_TYPE_LABEL[t.claim_type] ?? t.claim_type, t.count])
  const workloads = [...data.by_workload].sort((a, b) => b.count - a.count).map((w): [string, number] => [w.workload, w.count])
  return (
    <div className="grid gap-10 border-t border-line px-5 py-6 md:grid-cols-2">
      <Bars title="By claim type" rows={types} />
      <div>
        <Bars title="By workload" rows={workloads} />
        <Link to={link('/workloads')} className="mt-3 inline-block text-xs text-link hover:underline">All workloads →</Link>
      </div>
    </div>
  )
}

function KeyCells({ keyValue }: { keyValue: Key }) {
  return (
    <>
      <td className="py-2 pl-3 font-mono text-xs text-ink">{keyValue.mbi}</td>
      <td className="py-2 font-mono text-xs text-ink-2">{keyValue.date_of_service}</td>
      <td className="py-2 text-right font-mono text-xs text-ink-2">{money(keyValue.total_charge)}</td>
      <td className="py-2 pl-4 font-mono text-xs text-ink-2">{keyValue.npi}</td>
      <td className="py-2 font-mono text-xs text-ink-2">{keyValue.workload}</td>
    </>
  )
}

// The vendor column reads by outcome: absent for missing, the claim with an
// amber edge for extra, and "sent more than once" for duplicated.
function ClaimList({ rows, total, side }: { rows: UnmatchedRow[]; total: number; side: 'missing' | 'extra' | 'duplicated' }) {
  const { carry } = useSelection()
  const navigate = useNavigate()
  return (
    <div className="overflow-x-auto border-t border-line px-2 pb-4">
      <table className="w-full min-w-[56rem] text-sm">
        <thead>
          <tr className="border-b-2 border-ink text-left text-xs text-ink">
            <th className="py-2 pl-3">Patient <Info term={help.terms.matchKey} /></th>
            <th className="py-2">Service date</th>
            <th className="py-2 text-right">Billed</th>
            <th className="py-2 pl-4">Payee NPI</th>
            <th className="py-2">Workload</th>
            <th className="py-2">Type</th>
            <th className="py-2">MAC-processed <Info term={help.terms.sides} /></th>
            <th className="py-2">Vendor-processed</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((c) => {
            const present = (
              <span className="flex flex-col">
                <span className="font-mono text-xs text-ink">Paid {money(c.paid)}</span>
                <span className="font-mono text-[11px] text-ink-2">{c.patient_control_number}</span>
              </span>
            )
            return (
              <tr key={`${c.match_key}-${c.patient_control_number}`} onClick={() => navigate(claimLink(carry, c.match_key))} className="cursor-pointer border-b border-line hover:bg-surface">
                <KeyCells keyValue={c.key} />
                <td className="py-2 text-xs text-ink">{CLAIM_TYPE_LABEL[c.claim_type] ?? c.claim_type}</td>
                <td className="py-2">{side === 'extra' ? <Absent /> : present}</td>
                <td className="py-2">
                  {side === 'missing' ? <Absent /> : side === 'duplicated' ? <Repeated /> : <span className="inline-block border-l-4 border-warn bg-warn/15 py-0.5 pr-2 pl-2">{present}</span>}
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
      <div className="px-3 pt-3 font-mono text-[11px] text-ink-2">
        Showing {count(rows.length)} of {count(total)}
      </div>
    </div>
  )
}

const KEY_FIELDS: [keyof Key, string][] = [
  ['mbi', 'Patient'], ['date_of_service', 'Service date'], ['total_charge', 'Billed'], ['tax_id', 'Tax ID'], ['npi', 'Payee NPI'], ['workload', 'Workload'],
]

// Rows sharing a match key arrive adjacent; each run of one key is a group.
function collisions(rows: UnmatchedRow[]) {
  const groups: UnmatchedRow[][] = []
  for (const r of rows) {
    const last = groups[groups.length - 1]
    if (last && last[0].match_key === r.match_key) last.push(r)
    else groups.push([r])
  }
  return groups
}

function CollisionList({ rows, total }: { rows: UnmatchedRow[]; total: number }) {
  const { carry } = useSelection()
  const groups = collisions(rows)
  return (
    <div className="space-y-5 border-t border-line px-5 py-6">
      {groups.map((group) => (
        <div key={group[0].match_key} className="border border-line">
          <div className="flex flex-wrap items-stretch border-b border-line bg-surface">
            {KEY_FIELDS.map(([field, label]) => (
              <div key={field} className="border-r border-line px-4 py-2">
                <div className="text-[11px] text-ink-2">{label}</div>
                <div className="font-mono text-sm text-ink">{field === 'total_charge' ? money(group[0].key.total_charge) : group[0].key[field]}</div>
              </div>
            ))}
            <div className="ml-auto flex items-center gap-2 px-4 py-2">
              <span className="text-xs text-ink-2">{CLAIM_TYPE_LABEL[group[0].claim_type] ?? group[0].claim_type}</span>
              <span className="flex h-7 w-7 items-center justify-center rounded-full bg-warn text-sm font-bold text-black">×{group.length}</span>
            </div>
          </div>
          <div className="divide-y divide-line">
            {group.map((c, i) => (
              <Link key={`${c.patient_control_number}-${i}`} to={claimLink(carry, c.match_key)} className="grid grid-cols-[2rem_1fr_7rem] items-center gap-3 px-4 py-2 hover:bg-surface">
                <span className="font-mono text-xs text-ink-2">{i + 1}</span>
                <span className="font-mono text-xs text-ink">{c.patient_control_number}</span>
                <span className="text-right font-mono text-xs text-ink">{money(c.paid)}</span>
              </Link>
            ))}
          </div>
        </div>
      ))}
      <div className="flex items-center justify-between font-mono text-[11px] text-ink-2">
        <span className="flex items-center gap-1">Showing {count(rows.length)} of {count(total)} claims · excluded from match rates, counted here<Info term={help.terms.collisions} /></span>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------- page

export function ViewUnmatched() {
  const [params, setParams] = useSearchParams()
  const sel = useSelection()
  const submission = sel.shown
  const tabId = TABS.find((t) => t.id === params.get('outcome'))?.id ?? 'missing'
  const tab = TABS.find((t) => t.id === tabId)!

  const measures = useApi(submission ? `measures:${submission.prefix}` : null, (signal) => api.measures(submission!.prefix, signal))
  const unmatched = useApi(submission ? `unmatched:${submission.prefix}|${tabId}` : null, (signal) => api.unmatched({ submission: submission!.prefix, outcome: tabId }, signal))
  const totals = measures.data?.measures.find((m) => m.claim_type === 'all' && m.workload === 'all')
  const tabCount = (id: Unmatched) => totals?.[id]
  const go = (outcome: Unmatched) => {
    const next = new URLSearchParams()
    // '' (an explicit empty slot, sel.carry) is kept; undefined is dropped.
    for (const [k, v] of Object.entries({ ...sel.carry, outcome })) if (v !== undefined) next.set(k, v)
    setParams(next)
  }

  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12">
      <PageTitle purpose={help.purpose}>Missing, extra, duplicated, and set aside</PageTitle>

      {submission && (
        <div className="mb-6">
          <SideToggle />
          <div className="mt-3 flex flex-wrap items-center gap-x-6 gap-y-2 text-xs text-ink-2">
            {totals && <span className="font-mono">{count(totals.mac_claims)} MAC · {count(submission.claim_count)} vendor</span>}
            <span className="ml-auto flex gap-4">
              <Link to={sel.link('/')} className="text-link hover:underline">Overview</Link>
              <Link to={sel.link('/differences')} className="text-link hover:underline">Differences</Link>
            </span>
          </div>
        </div>
      )}

      <div className="mb-8 inline-flex border border-line" role="tablist">
        {TABS.map((t) => {
          const n = tabCount(t.id)
          return (
            <button key={t.id} type="button" role="tab" aria-selected={tabId === t.id} onClick={() => go(t.id)}
              className={`flex items-center gap-2 px-5 py-2.5 text-sm ${tabId === t.id ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
              {n !== undefined && <span className={`h-2.5 w-2.5 rounded-full ${n === 0 ? 'bg-ok' : 'bg-warn'}`} />}
              {t.label}
              {n !== undefined && <span className="font-mono text-xs">{count(n)}</span>}
            </button>
          )
        })}
      </div>

      {unmatched.data && submission ? (
        <section className="border border-line">
          <Headline tab={tab} data={unmatched.data} of={totals?.mac_claims} />
          <BreakdownPanel data={unmatched.data} />
          {tabId === 'set_aside'
            ? <CollisionList rows={unmatched.data.rows} total={unmatched.data.count} />
            : <ClaimList rows={unmatched.data.rows} total={unmatched.data.count} side={tabId} />}
        </section>
      ) : sel.unselected ? (
        <Empty>Select a submission in the navigation.</Empty>
      ) : (
        <Pending errors={[sel.error, unmatched.error]} loading={sel.loading || unmatched.loading} empty="No vendor submission has been compared yet" />
      )}
    </div>
  )
}
