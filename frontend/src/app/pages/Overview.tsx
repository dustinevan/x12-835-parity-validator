import { useEffect, useMemo, useRef, useState } from 'react'
import { useSearchParams } from 'react-router'
import { PageTitle } from '../../components/PageTitle'
import { Empty, Pending } from '../../components/ViewState'
import { api, useApi, type FrontFilter, type FrontResponse } from '../../api'
import { HELP } from '../help'
import { useSelection, type Side } from '../selection'
import { readSelection, selectionKey, writeSelection, type OverviewSelection } from '../overviewFilters'
import { OverviewView, type OverviewSide } from '../OverviewView'

// The overview: one /api/front request per picked submission, again once the
// selection (the address's filter parameters, which the selection header's
// boxes write) has held still for a moment, so a run of quick clicks sends
// one request per slot. While a newer response is in flight the last one
// stays on the screen with a quiet mark, so the numbers never blank out.
// With both slots picked, the agreement grid reads /api/agreement once per
// pair; the selection does not narrow it. The MAC map's tree (/api/macs,
// static and cacheable) is read once, to group the workload box's workloads
// under their jurisdictions.

const help = HELP.overview

const SETTLE = 200 // ms the selection holds still before its requests go out

// The last response for a submission, kept while the next one loads; dropped
// when the submission changes.
function useLast(prefix: string | undefined, data: FrontResponse | undefined): FrontResponse | undefined {
  const last = useRef<{ prefix: string; data: FrontResponse } | undefined>(undefined)
  if (prefix !== undefined && data !== undefined) last.current = { prefix, data }
  return data ?? (prefix !== undefined && last.current?.prefix === prefix ? last.current.data : undefined)
}

// `key` once it has held still for `ms`; the first key at once.
function useSettled(key: string, ms: number): string {
  const [settled, setSettled] = useState(key)
  useEffect(() => {
    if (key === settled) return
    const timer = window.setTimeout(() => setSettled(key), ms)
    return () => window.clearTimeout(timer)
  }, [key, settled, ms])
  return settled
}

export function Overview() {
  const sel = useSelection()
  const [params, setParams] = useSearchParams()
  const selection = useMemo(() => readSelection(params), [params])
  const wanted = selectionKey(selection)
  const key = useSettled(wanted, SETTLE)
  // The request's filters, read back from the settled key (selectionKey is
  // the filters written as a query string).
  const filters = Object.fromEntries(new URLSearchParams(key)) as Partial<Record<FrontFilter, string>>

  const { v1, v2 } = sel
  const load1 = useApi(v1 ? `front:${v1.prefix}|${key}` : null, (signal) => api.front({ submission: v1!.prefix, ...filters }, signal))
  const load2 = useApi(v2 ? `front:${v2.prefix}|${key}` : null, (signal) => api.front({ submission: v2!.prefix, ...filters }, signal))
  const agreement = useApi(v1 && v2 ? `agreement:${v1.prefix}|${v2.prefix}` : null, (signal) => api.agreement({ v1: v1!.prefix, v2: v2!.prefix }, signal))
  const macs = useApi('macs', (signal) => api.macs(signal))
  const data1 = useLast(v1?.prefix, load1.data)
  const data2 = useLast(v2?.prefix, load2.data)

  const onSelection = (next: OverviewSelection) =>
    setParams((prev) => {
      const out = new URLSearchParams(prev)
      writeSelection(out, next)
      return out
    }, { replace: true })

  const sides: OverviewSide[] = [
    { side: 'v1' as Side, submission: v1, data: data1, loading: load1.loading, error: load1.error },
    { side: 'v2' as Side, submission: v2, data: data2, loading: load2.loading, error: load2.error },
  ].flatMap((s) => (s.submission ? [{ ...s, submission: s.submission }] : []))
  const ready = sides.length > 0 && sides.some((s) => s.data !== undefined)

  return (
    <div className="mx-auto max-w-[72rem] px-10 py-12">
      <PageTitle purpose={help.purpose}>Overview</PageTitle>
      {sel.unselected ? (
        <Empty>Select a submission in the navigation.</Empty>
      ) : ready ? (
        <OverviewView sides={sides} selection={selection} onSelection={onSelection} pending={wanted !== key} macs={macs.data}
          agreement={{ cells: agreement.data?.cells, error: agreement.error, loading: agreement.loading }} />
      ) : (
        <Pending errors={[sel.error, ...sides.map((s) => s.error)]} loading={sel.loading || sides.some((s) => s.loading)} empty="No vendor submission has been compared yet" />
      )}
    </div>
  )
}
