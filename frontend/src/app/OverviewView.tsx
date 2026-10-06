import { Fragment, useMemo, useState } from 'react'
import { Link } from 'react-router'
import type { AgreementCell, Compared, FrontResponse, FrontSums, FrontTag, FrontValue, MacsResponse, Outcome } from '../api'
import { Pending } from '../components/ViewState'
import { count, money, signedMoney, statusLabel, sum } from '../format'
import { SlotChip } from './VendorPicker'
import { CodeHint } from './CodeHint'
import { HELP, type Term } from './help'
import { Info } from './Info'
import { slotNumber, useSelection, type Side } from './selection'
import { Breakdowns } from './Breakdowns'
import { ColumnHeads, GRID, Gap, LegacyCell, LegacyChip, MATRIX_GRID, SideCell, Slots, Tile, legacyOf, slotCells, type OverviewSide } from './overviewColumns'
import { SelectionHeader } from './SelectionHeader'
import { OUTCOMES, PAYMENT_BANDS, hasVendorAmount, isSelected, valueLabel, type OverviewSelection } from './overviewFilters'

// The overview's view, drawn from one /api/front response per picked
// submission and, with both slots picked, the /api/agreement response for
// the pair. The page (pages/Overview.tsx) fetches and owns the address;
// this file only draws and reports clicks. Top to bottom: the three
// columns' names, the selection header (SelectionHeader.tsx, which stays at
// the top of the window as the rest scrolls), every claim by outcome, how
// the two agree, payment and status agreement among pairs, what differs,
// and the breakdowns (Breakdowns.tsx). Every table sits on one grid, the
// label then Legacy, 1, and 2 (overviewColumns.tsx). Legacy has figures
// only where the legacy side has its own, read from the response's legacy
// block: the outcome table's total, the MAC claims beside each MAC status,
// and the breakdowns by what the MAC claim carries; everywhere else its
// column stays empty. Every number is a count or a sum; where a figure is
// one count out of another, both are printed.

export type { OverviewSide }

const help = HELP.overview

// A link's side parameter: side=2 for slot 2; slot 1 is the default.
const onSide = (side: Side) => ({ side: side === 'v2' ? '2' : undefined })

// Which Claims-page filter an overview selection can carry: the claim type
// when exactly one is selected. The Claims page has no other dimension.
function claimsFilter(selection: OverviewSelection): Record<string, string | undefined> {
  const types = selection.claim_type
  return { claim_type: types !== undefined && types.length === 1 ? types[0] : undefined }
}

// ---------------------------------------------------------------- every claim, by outcome

interface OutcomeRow {
  outcome: Outcome
  label: string
  fill: string // the bar segment
  ink: string // the count on the segment
  swatch: string // the legend's square
  term: Term
  view: '/claims' | '/unmatched' // where a count opens (Selection.link)
}

// In the order of overviewFilters.ts OUTCOMES: the outcomes with a MAC claim,
// then extra, so a vendor's extra claims come last on its bar.
const OUTCOME_ROWS: Record<Outcome, OutcomeRow> = {
  equal: { outcome: 'equal', label: 'Equal', fill: 'bg-ok', ink: 'text-white', swatch: 'bg-ok', term: help.terms.equal, view: '/claims' },
  different: { outcome: 'different', label: 'Different', fill: 'bg-warn', ink: 'text-black', swatch: 'bg-warn', term: help.terms.different, view: '/claims' },
  missing: {
    outcome: 'missing', label: 'Missing from vendor', fill: 'border-2 border-dashed border-ink-2 bg-bg', ink: 'text-ink',
    swatch: 'border border-dashed border-ink-2', term: help.terms.missing, view: '/claims',
  },
  duplicated: {
    outcome: 'duplicated', label: 'Duplicated by vendor', fill: 'border-2 border-dotted border-ink-2 bg-surface', ink: 'text-ink',
    swatch: 'border border-dotted border-ink-2 bg-surface', term: help.terms.duplicated, view: '/claims',
  },
  set_aside: { outcome: 'set_aside', label: 'Set aside', fill: 'bg-mark-off', ink: 'text-white', swatch: 'bg-mark-off', term: help.terms.setAside, view: '/unmatched' },
  extra: { outcome: 'extra', label: 'Extra on vendor', fill: 'bg-ink-2', ink: 'text-white', swatch: 'bg-ink-2', term: help.terms.extra, view: '/unmatched' },
}
const OUTCOME_ORDER = OUTCOMES.map((o) => OUTCOME_ROWS[o as Outcome])

const outcomeOf = (data: FrontResponse | undefined, outcome: Outcome): FrontSums | undefined => data?.outcomes.find((o) => o.outcome === outcome)

function useOutcomeLink(selection: OverviewSelection) {
  const { link } = useSelection()
  return (side: Side, row: OutcomeRow) => {
    if (row.view === '/unmatched') return link(row.view, { outcome: row.outcome, ...onSide(side) })
    return link(row.view, { [side === 'v2' ? 'v2_outcome' : 'v1_outcome']: row.outcome, ...onSide(side), ...claimsFilter(selection) })
  }
}

// One bar per column on one scale, each with its mark: legacy's is the MAC
// claims under the selection on what the MAC claim carries; a vendor's is
// its results by outcome under the whole selection, its extra claims last.
function OutcomeBars({ sides, selection }: { sides: OverviewSide[]; selection: OverviewSelection }) {
  const outcomeLink = useOutcomeLink(selection)
  const legacy = legacyOf(sides)?.totals
  const top = Math.max(1, legacy?.claims ?? 0, ...sides.map((s) => s.data?.totals.claims ?? 0))
  const width = (n: number) => `${(n / top) * 100}%`
  return (
    <div className="space-y-2">
      {legacy && (
        <div className="flex items-center gap-3" data-bar="legacy">
          <LegacyChip />
          <div className="flex h-9 flex-1">
            {legacy.claims > 0 && (
              <div title={`Legacy: ${count(legacy.claims)} MAC claims`} style={{ width: width(legacy.claims), minWidth: 5 }}
                className="flex items-center justify-center overflow-hidden border-2 border-ink bg-surface font-mono text-xs font-semibold whitespace-nowrap text-ink">
                {count(legacy.claims)} MAC claims
              </div>
            )}
          </div>
        </div>
      )}
      {sides.map((s) => (
        <div key={s.side} className="flex items-center gap-3" data-bar={s.side}>
          <SlotChip side={s.side} />
          <div className="flex h-9 flex-1">
            {OUTCOME_ORDER.map((o) => {
              const n = outcomeOf(s.data, o.outcome)?.claims ?? 0
              if (n === 0) return null
              return (
                <Link key={o.outcome} to={outcomeLink(s.side, o)} title={`${o.label}: ${count(n)}`} style={{ width: width(n), minWidth: 5 }}
                  className={`${o.fill} flex items-center justify-center overflow-hidden border-r-2 border-bg font-mono text-xs font-semibold whitespace-nowrap ${o.ink} hover:opacity-85`}>
                  {n / top > 0.08 && count(n)}
                </Link>
              )
            })}
          </div>
        </div>
      ))}
    </div>
  )
}

// The outcome table: one row per outcome, each vendor's results and what it
// paid; Legacy has no figure per outcome, since an outcome is a vendor's
// result, so its column holds only the total: the MAC claims and what the
// MAC paid on them.
function EveryClaim({ sides, selection }: { sides: OverviewSide[]; selection: OverviewSelection }) {
  const outcomeLink = useOutcomeLink(selection)
  const legacy = legacyOf(sides)
  const rule = 'border-t border-line pt-2'
  return (
    <Tile title="Every claim, by outcome" info={help.terms.outcomes} right="claims · paid (CLP04)">
      <OutcomeBars sides={sides} selection={selection} />
      <div className="mt-4 flex flex-wrap items-center gap-x-5 gap-y-1.5">
        <span className="flex items-center gap-1.5 text-[11px] text-ink-2"><span className="h-2.5 w-2.5 shrink-0 border border-ink bg-surface" />Legacy MAC claims</span>
        {OUTCOME_ORDER.map((o) => (
          <span key={o.outcome} className="flex items-center gap-1.5 text-[11px] text-ink-2">
            <span className={`h-2.5 w-2.5 shrink-0 ${o.swatch}`} />
            {o.label}
            <Info term={o.term} />
          </span>
        ))}
      </div>
      <div className="mt-5 overflow-x-auto">
        <div className="grid min-w-max items-center gap-x-6 gap-y-2" style={GRID} data-table="outcomes" data-legacy="1">
          <span className="flex items-center gap-1 text-[11px] text-ink-2">Outcome<Info term={help.terms.columns} /></span>
          <ColumnHeads sides={sides} legacy={['MAC claims', 'MAC paid']} vendor={['claims', 'paid']} />
          {OUTCOME_ORDER.map((o) => (
            <Fragment key={o.outcome}>
              <span className="flex items-center gap-1.5 text-xs text-ink"><span className={`h-2.5 w-2.5 shrink-0 ${o.swatch}`} />{o.label}</span>
              <Gap column="legacy" />
              {slotCells(sides, (s) => (
                <SideCell sums={outcomeOf(s.data, o.outcome)} amount={hasVendorAmount('outcome', o.outcome)} to={outcomeLink(s.side, o)} />
              ))}
            </Fragment>
          ))}
          <span className={`${rule} text-xs font-semibold text-ink`}>Total</span>
          <span className={rule} data-total="legacy"><LegacyCell sums={legacy?.totals} /></span>
          {slotCells(sides, (s) => <span className={rule}><SideCell sums={s.data?.totals} /></span>)}
        </div>
      </div>
    </Tile>
  )
}

// ---------------------------------------------------------------- how the two agree

const COMPARED: Compared[] = ['equal', 'different', 'missing', 'duplicated']
const COMPARED_LABEL: Record<Compared, string> = { equal: 'Equal', different: 'Different', missing: 'Missing', duplicated: 'Duplicated' }

// Both equal is green; one different is light amber; both different, the
// cell to read first, is darker amber; a missing or duplicated side is plain.
function cellTint(outcome1: Compared, outcome2: Compared) {
  if (outcome1 === 'missing' || outcome2 === 'missing' || outcome1 === 'duplicated' || outcome2 === 'duplicated') return 'bg-surface'
  if (outcome1 === 'equal' && outcome2 === 'equal') return 'bg-ok/15'
  if (outcome1 === 'different' && outcome2 === 'different') return 'bg-warn/40'
  return 'bg-warn/15'
}

// A share printed beside a count: that count ÷ the grid's total, both shown.
const shareText = (share: number) => (share > 0 && share < 0.01 ? `${(share * 1000).toFixed(1)} per 1,000` : `${(share * 100).toFixed(1)}%`)

// The four columns of slot 2's outcomes, across the page grid's three figure
// columns, so the grid's edges meet the page's Legacy and slot 2 lines.
const ACROSS = { gridColumn: '2 / span 3', gridTemplateColumns: 'repeat(4, minmax(0, 1fr))' }

// The agreement grid: each MAC claim once, by its outcome in
// slot 1 (rows) and in slot 2 (columns), from /api/agreement, whose cells
// name the two outcomes v1 and v2. It covers both whole submissions: the
// selection does not narrow it. Slot 1's outcomes sit in the page grid's
// label column and slot 2's across its figure columns. Each cell opens
// Claims for those claims, showing the slot that differs (slot 2 when only
// slot 2 differs, else slot 1); the toggle there shows the other.
function Agreement({ cells }: { cells: AgreementCell[] }) {
  const { link } = useSelection()
  const total = cells.reduce((n, c) => n + c.claims, 0)
  const find = (outcome1: Compared, outcome2: Compared) => cells.find((c) => c.v1 === outcome1 && c.v2 === outcome2)
  return (
    <section className="border border-line px-4 py-5" data-section="agreement">
      <div className="mb-4 flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1">
        <h2 className="flex items-center gap-1 text-xs font-semibold text-ink">How the two agree, per MAC claim<Info term={help.terms.agreement} /></h2>
        <span className="font-mono text-[11px] text-ink-2">
          {count(total)} MAC claims, set-aside claims left out · whole submissions, not narrowed by the selection · a share is the cell ÷ {count(total)}
        </span>
      </div>
      <div className="overflow-x-auto">
        <div className="grid min-w-max items-stretch gap-x-6 gap-y-1" style={GRID} data-grid="agreement">
          <span />
          <div className="grid gap-1" style={ACROSS}>
            {COMPARED.map((outcome2) => (
              <span key={outcome2} className="flex items-center gap-2 px-3 pb-1 text-xs font-semibold text-ink"><SlotChip side="v2" />{COMPARED_LABEL[outcome2]}</span>
            ))}
          </div>
          {COMPARED.map((outcome1) => (
            <Fragment key={outcome1}>
              <span className="flex items-center gap-2 text-xs font-semibold text-ink"><SlotChip side="v1" />{COMPARED_LABEL[outcome1]}</span>
              <div className="grid gap-1" style={ACROSS}>
                {COMPARED.map((outcome2) => {
                  const cell = find(outcome1, outcome2)
                  if (!cell) return <span key={outcome2} className="flex items-center px-3 py-3 font-mono text-sm text-ink-2">—</span>
                  const side: Side = outcome1 === 'equal' && outcome2 !== 'equal' ? 'v2' : 'v1'
                  return (
                    <Link key={outcome2} to={link('/claims', { v1_outcome: outcome1, v2_outcome: outcome2, ...onSide(side) })} data-agreement={`${outcome1}:${outcome2}`}
                      className={`block px-3 py-3 hover:outline hover:outline-1 hover:outline-ink ${cellTint(outcome1, outcome2)}`}>
                      <div className="flex items-baseline gap-2">
                        <span className="font-mono text-lg font-semibold text-ink">{count(cell.claims)}</span>
                        <span className="font-mono text-[11px] text-ink-2">{shareText(total === 0 ? 0 : cell.claims / total)}</span>
                      </div>
                      {cell.same_tags > 0 && (
                        <div className="mt-1 flex items-center gap-1 text-[11px] text-ink">
                          {count(cell.same_tags)} differ in the same elements<Info term={help.terms.sameElements} />
                        </div>
                      )}
                    </Link>
                  )
                })}
              </div>
            </Fragment>
          ))}
        </div>
      </div>
    </section>
  )
}

// ---------------------------------------------------------------- payment agreement

// The pairs of one side, equal plus different claims, with what each side
// paid on them.
function pairsOf(data: FrontResponse | undefined) {
  const equal = outcomeOf(data, 'equal')
  const different = outcomeOf(data, 'different')
  const both = [equal, different].filter((x): x is FrontSums => x !== undefined)
  return {
    claims: both.reduce((n, x) => n + x.claims, 0),
    mac_claims: both.reduce((n, x) => n + x.mac_claims, 0),
    mac_paid: sum(both.map((x) => x.mac_paid)),
    vendor_paid: sum(both.map((x) => x.vendor_paid)),
  }
}

// The payment bands of one side, every band the facet lists (it is added
// under every other dimension's filter), and those of them the selection
// keeps.
function bandsOf(data: FrontResponse | undefined, selection: OverviewSelection) {
  const all = new Map<string, FrontValue>()
  for (const v of data?.dimensions.find((d) => d.name === 'payment_band')?.values ?? []) all.set(v.value, v)
  const kept = (band: string) => (isSelected(selection, 'payment_band', band) ? all.get(band)?.claims ?? 0 : 0)
  return { all, kept }
}

interface Line {
  label: string
  info?: Term
  read: (s: OverviewSide) => { left: string; right: string }
}

// Payment agreement is a comparison of each vendor's pairs with the MAC
// claims in them: the legacy side has no figure of its own here, so its
// column stays empty in both the table and the bands chart.
function Payment({ sides, selection }: { sides: OverviewSide[]; selection: OverviewSelection }) {
  const bands = PAYMENT_BANDS.filter((b) => b !== '')
  const perSide = new Map(sides.map((s) => [s.side, bandsOf(s.data, selection)]))
  const kept = (s: OverviewSide, band: string) => perSide.get(s.side)?.kept(band) ?? 0
  const top = Math.max(1, ...sides.flatMap((s) => bands.map((b) => perSide.get(s.side)?.all.get(b)?.claims ?? 0)))
  const lines: Line[] = [
    { label: 'Same payment (CLP04)', info: help.terms.samePayment, read: (s) => ({ left: count(kept(s, 'match')), right: 'within $0.05' }) },
    { label: 'Different payment', read: (s) => ({ left: count(bands.filter((b) => b !== 'match').reduce((n, b) => n + kept(s, b), 0)), right: '' }) },
    { label: 'Net difference', info: help.terms.netGross, read: (s) => ({ left: '', right: s.data ? signedMoney(s.data.totals.net_difference) : '—' }) },
    { label: 'Gross difference', read: (s) => ({ left: '', right: s.data ? money(s.data.totals.gross_difference) : '—' }) },
    {
      label: 'Same status (CLP02)', info: help.terms.sameStatus,
      read: (s) => (s.data ? { left: count(s.data.totals.same_status), right: `of ${count(pairsOf(s.data).claims)} pairs` } : { left: '—', right: '' }),
    },
  ]
  return (
    <Tile title="Payment agreement among pairs" info={help.terms.payment} right="a pair is one MAC claim and one vendor claim with the same key">
      <div className="overflow-x-auto">
        <div className="grid min-w-max items-center gap-x-6 gap-y-2" style={GRID} data-table="payment" data-legacy="0">
          <span />
          <ColumnHeads sides={sides} vendor={['pairs', 'paid']} />
          <span className="flex items-center gap-1 text-xs text-ink">Pairs (equal and different)<Info term={help.terms.pairsPaid} /></span>
          <Gap column="legacy" />
          {slotCells(sides, (s) => <SideCell sums={s.data ? pairsOf(s.data) : undefined} />)}
          {lines.map((line) => (
            <Fragment key={line.label}>
              <span className="flex items-center gap-1 text-xs text-ink-2">{line.label}{line.info && <Info term={line.info} />}</span>
              <Gap column="legacy" />
              {slotCells(sides, (s) => {
                const { left, right } = line.read(s)
                return <Slots left={left} right={right} strong={left === ''} />
              })}
            </Fragment>
          ))}
        </div>
      </div>
      <div className="mt-6">
        <div className="mb-2 flex items-center gap-1 text-[11px] text-ink-2">Pairs by how far apart the two payments are<Info term={help.terms.bands} /></div>
        <div className="overflow-x-auto">
          <div className="grid min-w-max items-center gap-x-6 gap-y-2" style={GRID} data-chart="payment-bands" data-legacy="0">
            <span />
            <ColumnHeads sides={sides} vendor="pairs" start />
            {bands.map((band) => {
              const on = isSelected(selection, 'payment_band', band)
              return (
                <Fragment key={band}>
                  <span className={`truncate text-xs ${on ? 'text-ink' : 'text-ink-2 line-through'}`} data-left-out={on ? undefined : 'payment_band'}
                    title={on ? undefined : 'Left out of the selection by the address'}>
                    {valueLabel('payment_band', band)}
                  </span>
                  <Gap column="legacy" />
                  {slotCells(sides, (s) => {
                    const n = perSide.get(s.side)?.all.get(band)?.claims ?? 0
                    return (
                      <span className="flex items-center gap-2" data-band={`${s.side}:${band}`}>
                        <span className={`h-2.5 shrink-0 ${on ? 'bg-series-1' : 'bg-mark-off'}`} style={{ width: `${(n / top) * 75}%`, minWidth: n > 0 ? 2 : 0 }} />
                        <span className={`font-mono text-[11px] tabular-nums ${on ? 'text-ink' : 'text-ink-2'}`}>{count(n)}</span>
                      </span>
                    )
                  })}
                </Fragment>
              )
            })}
          </div>
        </div>
      </div>
    </Tile>
  )
}

// ---------------------------------------------------------------- status agreement

const byCode = (x: string, y: string) => Number(x) - Number(y) || x.localeCompare(y)

// One slot's matrix: its own claim status codes across, in a column grid
// that fills the slot's column.
const codesAcross = (n: number) => ({ gridTemplateColumns: `repeat(${Math.max(1, n)}, minmax(2.5rem, 1fr))` })

// The MAC claim status down the label column, the legacy MAC claims with
// each status in the Legacy column, then each vendor's matrix of its own
// claim status across, slot 1's in its column and slot 2's in its, in the
// agreement grid's cell style. A matrix too wide for its column widens it.
function StatusMatrix({ sides }: { sides: OverviewSide[] }) {
  const groups = new Map(sides.map((s) => {
    const cells = s.data?.status_matrix ?? []
    return [s.side, { cells, vendorCodes: [...new Set(cells.map((c) => c.vendor_status))].sort(byCode), pairs: cells.reduce((n, c) => n + c.claims, 0) }]
  }))
  const legacy = new Map((legacyOf(sides)?.dimensions.find((d) => d.name === 'mac_status')?.values ?? []).map((v) => [v.value, v]))
  const macCodes = [...new Set([...groups.values()].flatMap((g) => g.cells.map((c) => c.mac_status)))].sort(byCode)
  const right = sides.map((s) => `${slotNumber(s.side)}: ${count(groups.get(s.side)?.pairs ?? 0)} pairs`).join(' · ')
  return (
    <Tile title="Status agreement among pairs" info={help.terms.matrix} right={right}>
      {macCodes.length === 0 ? (
        <p className="font-mono text-[11px] text-ink-2">No pairs under this selection.</p>
      ) : (
        <div className="overflow-x-auto">
          <div className="grid min-w-max items-center gap-x-6 gap-y-1" style={MATRIX_GRID} data-table="status" data-legacy="1">
            <span className="self-end text-[11px] text-ink-2">MAC claim status ↓</span>
            <span data-column="legacy" className="flex flex-col items-end gap-1 self-end">
              <span className="flex items-center gap-2 text-xs font-semibold text-ink"><LegacyChip />Legacy</span>
              <span className="flex justify-end gap-3 text-[11px] text-ink-2"><span className="w-16 shrink-0 text-right">MAC claims</span><span className="w-28 shrink-0" /></span>
            </span>
            {slotCells(sides, (s) => {
              const g = groups.get(s.side)
              return (
                <span data-column={s.side} className="flex flex-col gap-1 self-end">
                  <span className="flex items-center gap-2 text-xs font-semibold text-ink"><SlotChip side={s.side} />{s.submission.label}<span className="font-normal text-[11px] text-ink-2">its status →</span></span>
                  <span className="grid gap-1" style={codesAcross(g?.vendorCodes.length ?? 0)}>
                    {(g?.vendorCodes ?? []).map((v) => (
                      <span key={v} className="px-1.5 text-right font-mono text-xs text-ink"><CodeHint kind="status" code={v} label={statusLabel(v)}>{v}</CodeHint></span>
                    ))}
                  </span>
                </span>
              )
            })}
            {macCodes.map((m) => (
              <Fragment key={m}>
                <span className="font-mono text-xs whitespace-nowrap text-ink">
                  <CodeHint kind="status" code={m} label={statusLabel(m)}>{m}</CodeHint>
                  <span className="ml-2 font-sans text-[11px] text-ink-2">{statusLabel(m)}</span>
                </span>
                <span data-legacy-status={m}><LegacyCell sums={legacy.get(m)} amount={false} /></span>
                {slotCells(sides, (s) => {
                  const g = groups.get(s.side)
                  return (
                    <span className="grid gap-1" style={codesAcross(g?.vendorCodes.length ?? 0)}>
                      {g?.vendorCodes.length === 0 && <span className="px-1.5 py-1.5 text-right font-mono text-ink-2/50">·</span>}
                      {(g?.vendorCodes ?? []).map((v) => {
                        const cell = g?.cells.find((c) => c.mac_status === m && c.vendor_status === v)
                        const same = m === v
                        return (
                          <span key={v} data-status-cell={cell ? (same ? 'same' : 'other') : undefined}
                            className={`px-1.5 py-1.5 text-right font-mono text-xs tabular-nums ${cell ? `font-semibold ${same ? 'bg-ok/15 text-ink' : 'bg-warn/20 text-ink'}` : 'text-ink-2/50'}`}>
                            {cell ? count(cell.claims) : '·'}
                          </span>
                        )
                      })}
                    </span>
                  )
                })}
              </Fragment>
            ))}
          </div>
        </div>
      )}
    </Tile>
  )
}

// ---------------------------------------------------------------- what differs

// The segment a tag belongs to, for grouping: CLP04 is CLP, CAS*CO*45 is CAS,
// MOA, MIA, and LQ are remark lists, RAW*K3 is an unplaced segment.
const GROUP_ORDER = ['CLP', 'SVC', 'CAS', 'Remarks', 'NM1', 'REF', 'AMT', 'DTM', 'QTY', 'MIA', 'MOA', 'RAW', 'Other']

function groupOf(tag: string): string {
  if (tag.startsWith('RAW*')) return 'RAW'
  const head = tag.split('*')[0]
  if (head === 'MOA' || head === 'MIA' || head === 'LQ') return head === 'LQ' || /^(MOA|MIA)$/.test(tag) ? 'Remarks' : head
  const seg = /^[A-Z][A-Z0-9]{1,2}/.exec(head)?.[0] ?? 'Other'
  return GROUP_ORDER.includes(seg) ? seg : 'Other'
}

interface TagRow {
  tag: string
  label: string
  sums: Partial<Record<Side, FrontTag>>
}

function mergeTags(sides: OverviewSide[]): TagRow[] {
  const byTag = new Map<string, TagRow>()
  for (const s of sides) {
    for (const t of s.data?.tags ?? []) {
      const row = byTag.get(t.tag) ?? { tag: t.tag, label: t.label, sums: {} }
      row.sums[s.side] = t
      byTag.set(t.tag, row)
    }
  }
  return [...byTag.values()]
}

// Both submissions' counts added together set the order, so neither one's
// list sets it.
const claimsSum = (row: TagRow) => Object.values(row.sums).reduce((n, t) => n + (t?.claims ?? 0), 0)
const dollarsSum = (row: TagRow) => Object.values(row.sums).reduce((n, t) => n + Number(t?.dollars ?? 0), 0)

// What differs is between the legacy claim and a vendor's, so the legacy
// side has no count of its own and its column stays empty.
function Differs({ sides, selection }: { sides: OverviewSide[]; selection: OverviewSelection }) {
  const { link } = useSelection()
  const [sort, setSort] = useState<'claims' | 'dollars'>('claims')
  const [query, setQuery] = useState('')
  const rows = useMemo(() => mergeTags(sides), [sides])
  const q = query.trim().toLowerCase()
  const shown = rows.filter((r) => q === '' || r.tag.toLowerCase().includes(q) || r.label.toLowerCase().includes(q))
  const order = (x: TagRow, y: TagRow) => (sort === 'claims' ? claimsSum(y) - claimsSum(x) : dollarsSum(y) - dollarsSum(x)) || x.tag.localeCompare(y.tag)
  const groups = GROUP_ORDER.map((g) => ({ group: g, rows: shown.filter((r) => groupOf(r.tag) === g).sort(order) })).filter((g) => g.rows.length > 0)
  return (
    <Tile title="What differs, among different pairs" info={help.terms.differs} right="a claim with several differing elements counts under each">
      <div className="mb-3 flex flex-wrap items-center gap-3 text-xs">
        <input type="search" value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Find an element" aria-label="Find an element"
          className="w-48 border border-line bg-bg px-2 py-1 text-xs text-ink placeholder:text-ink-2/60" />
        <span className="text-ink-2">Sort by</span>
        <div className="inline-flex border border-line" role="tablist" aria-label="Sort">
          {(['claims', 'dollars'] as const).map((s) => (
            <button key={s} type="button" role="tab" aria-selected={sort === s} onClick={() => setSort(s)}
              className={`border-r border-line px-3 py-1 last:border-r-0 ${sort === s ? 'bg-bg font-semibold text-ink' : 'bg-surface text-ink-2 hover:text-ink'}`}>
              {s}
            </button>
          ))}
        </div>
      </div>
      {groups.length === 0 && <p className="font-mono text-[11px] text-ink-2">{rows.length === 0 ? 'No different pairs under this selection.' : 'No element matches the search.'}</p>}
      <div className="overflow-x-auto">
        {groups.map((g) => (
          <div key={g.group} className="mb-3 grid min-w-max items-center gap-x-6" style={GRID} data-table="differs" data-legacy="0">
            <span className="self-end px-2 pb-1 font-mono text-[11px] font-semibold tracking-wider text-ink-2 uppercase">{g.group}</span>
            <ColumnHeads sides={sides} vendor={['claims', 'dollars']} />
            {g.rows.map((r) => (
              <Fragment key={r.tag}>
                <span className="flex min-w-0 items-baseline gap-3 px-2 py-1">
                  <span className="w-24 shrink-0 truncate font-mono text-[11px] text-ink-2"><CodeHint kind="element" code={r.tag} label={r.label}>{r.tag}</CodeHint></span>
                  <span className="truncate text-xs text-ink">{r.label}</span>
                </span>
                <Gap column="legacy" />
                {slotCells(sides, (s) => {
                  const t = r.sums[s.side]
                  return (
                    <Slots left={t ? count(t.claims) : '—'} right={t ? money(t.dollars) : '—'}
                      to={t ? link('/claims', { tag: r.tag, ...onSide(s.side), ...claimsFilter(selection) }) : undefined} />
                  )
                })}
              </Fragment>
            ))}
          </div>
        ))}
      </div>
    </Tile>
  )
}

// ---------------------------------------------------------------- the view

export interface AgreementLoad {
  cells: AgreementCell[] | undefined
  error: Error | undefined
  loading: boolean
}

export function OverviewView({ sides, selection, onSelection, agreement, macs, pending = false }: {
  sides: OverviewSide[]
  selection: OverviewSelection
  onSelection: (next: OverviewSelection) => void
  agreement?: AgreementLoad // with both slots picked
  macs?: MacsResponse // the MAC map's tree, which groups the workload box
  pending?: boolean // the selection changed and its requests have not gone out yet
}) {
  const updating = pending || sides.some((s) => s.loading)
  return (
    <div className="space-y-10">
      <div className="flex flex-wrap items-center gap-x-8 gap-y-2">
        <span className="flex items-center gap-2">
          <LegacyChip />
          <span className="text-sm font-semibold text-ink">Legacy</span>
          <span className="font-mono text-[11px] text-ink-2">MAC-processed</span>
        </span>
        {sides.map((s) => (
          <span key={s.side} className="flex items-center gap-2">
            <SlotChip side={s.side} />
            <span className="text-sm font-semibold text-ink">{s.submission.label}</span>
            <span className="font-mono text-[11px] text-ink-2">{s.submission.submitted_on}</span>
          </span>
        ))}
        <Info term={help.terms.vendors} />
      </div>
      <SelectionHeader sides={sides} selection={selection} onSelection={onSelection} macs={macs} updating={updating} />
      <EveryClaim sides={sides} selection={selection} />
      {sides.length === 2 && agreement && (
        agreement.cells
          ? <Agreement cells={agreement.cells} />
          : <Pending errors={[agreement.error]} loading={agreement.loading} />
      )}
      <Payment sides={sides} selection={selection} />
      <StatusMatrix sides={sides} />
      <Differs sides={sides} selection={selection} />
      <Breakdowns sides={sides} selection={selection} />
    </div>
  )
}
