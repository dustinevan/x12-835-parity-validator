import { CARC, CAS_GROUP_NAME } from '../../data/x12Names835'
import { Term } from '../viewer/Term'
import { CAS_GROUP_BAR } from './colors'
import type { Box } from './parse'

interface Slice {
  amount: number
  group?: string
  reason?: string
}

// The charge split into the paid amount and every CAS amount in file order,
// across however many CAS segments the line carries. The bar sorts the
// adjustment slices by descending amount (paid always first) so the widest
// slices read left to right, matching the spec's drawing C.
function slices(svc: Box): Slice[] {
  const paid: Slice = { amount: Number(svc.opening.elements[2] ?? '0') }
  const adjustments: Slice[] = []
  for (const node of svc.children) {
    if (node.kind !== 'segment' || node.segment.id !== 'CAS') continue
    const e = node.segment.elements
    const group = e[0]
    for (let i = 1; i + 1 < e.length; i += 3) {
      const reason = e[i]
      if (reason === undefined) break
      adjustments.push({ amount: Number(e[i + 1] ?? '0'), group, reason })
    }
  }
  return [paid, ...adjustments.slice().sort((a, b) => b.amount - a.amount)]
}

// A slice's group and reason codes, each a coded value that shows its public
// name on hover (and for the reason, the maintainer's full text), as coded
// values do in the cells above.
function ReasonTerm({ group, reason }: { group: string; reason: string }) {
  const groupName = CAS_GROUP_NAME[group]
  const carc = CARC[reason]
  return (
    <>
      <Term
        trigger={group}
        triggerClassName="underline decoration-dotted decoration-line underline-offset-2 hover:text-ink"
        title={groupName ?? group}
      />{' '}
      <Term
        trigger={reason}
        triggerClassName="underline decoration-dotted decoration-line underline-offset-2 hover:text-ink"
        title={carc?.short ?? reason}
        body={carc && carc.full !== carc.short ? carc.full : undefined}
      />
    </>
  )
}

export function ServiceLineBar({ svc }: { svc: Box }) {
  const charge = Number(svc.opening.elements[1] ?? '0')
  if (!(charge > 0)) return null
  const parts = slices(svc)

  return (
    <div className="mt-3 max-w-[48rem]">
      <div className="font-mono text-sm text-ink">{svc.opening.elements[1]}</div>
      <div className="mt-1 flex h-6 w-full overflow-hidden border border-line">
        {parts.map((s, i) => (
          <div
            key={i}
            style={{ width: `${(s.amount / charge) * 100}%`, opacity: i === 0 ? 0.85 : 1 }}
            className={`h-full border-r border-bg last:border-r-0 ${i === 0 ? 'bg-ink' : (s.group && CAS_GROUP_BAR[s.group]) || 'bg-mark-off'}`}
          />
        ))}
      </div>
      <div className="mt-1 flex w-full">
        {parts.map((s, i) => (
          <div key={i} style={{ width: `${(s.amount / charge) * 100}%` }} className="shrink-0 pr-1 font-mono text-xs leading-5 text-ink">
            {s.amount.toFixed(2)}
          </div>
        ))}
      </div>
      <div className="flex w-full">
        {parts.map((s, i) => (
          <div key={i} style={{ width: `${(s.amount / charge) * 100}%` }} className="shrink-0 whitespace-normal pr-1 text-xs leading-5 text-ink-2">
            {s.group && s.reason && <ReasonTerm group={s.group} reason={s.reason} />}
          </div>
        ))}
      </div>
    </div>
  )
}
