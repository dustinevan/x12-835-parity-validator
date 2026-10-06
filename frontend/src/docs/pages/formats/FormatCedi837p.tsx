import type { ReactNode } from 'react'
import { DocSection } from '../../../components/DocSection'
import { Figure } from '../../../components/Figure'
import { PageTitle } from '../../../components/PageTitle'

// A visual digest of the CEDI 837P Standard Companion Guide, version 5.6
// (January 2025), based on 005010X222A1. Every fact cites a PDF viewer page.


const INK = 'var(--color-ink)'
const INK2 = 'var(--color-ink-2)'
const LINE = 'var(--color-line)'
const ACCENT = 'var(--color-accent)'
const SURFACE = 'var(--color-surface)'

// ---------------------------------------------------------------------------
// Figure 1: production and test transaction flows (PDF pp. 16-17)

function FlowBox({ title, children, strong }: { title: string; children: ReactNode; strong?: boolean }) {
  return (
    <div className={`border bg-surface p-3 text-[12px] leading-5 text-ink-2 ${strong ? 'border-accent' : 'border-line'}`}>
      <div className="mb-1 font-semibold text-ink">{title}</div>
      {children}
    </div>
  )
}

function Two({ label }: { label: string }) {
  return (
    <div className="flex flex-col items-center justify-center gap-1 px-1 text-[11px] text-ink-2">
      <svg width="44" height="16" viewBox="0 0 44 16" aria-hidden="true">
        <path d="M6 8 H38 M6 8 L12 3 M6 8 L12 13 M38 8 L32 3 M38 8 L32 13" stroke={INK2} strokeWidth="1.5" fill="none" />
      </svg>
      <span className="font-mono">{label}</span>
    </div>
  )
}

function ProductionFlow() {
  return (
    <div className="min-w-[760px] space-y-6">
      <div className="grid grid-cols-[1.2fr_auto_0.5fr_auto_1.6fr_auto_1fr] items-stretch gap-2">
        <div className="text-center text-[11px] font-semibold uppercase tracking-wide text-ink">Trading partners</div>
        <div />
        <div className="text-center text-[11px] font-semibold uppercase tracking-wide text-ink">Mailbox</div>
        <div />
        <div className="text-center text-[11px] font-semibold uppercase tracking-wide text-ink">CEDI processing</div>
        <div />
        <div className="text-center text-[11px] font-semibold uppercase tracking-wide text-ink">DME MACs</div>

        {/* Partner column: three ways in */}
        <div className="row-span-2 flex flex-col gap-2">
          <FlowBox title="sFTP Gateway">
            In: 837 claims, NCPDP claims, 276 status requests, X.509 certificate
            <br />
            Out: acknowledgments, <span className="text-accent">835 ERA</span>, 277 responses
          </FlowBox>
          <FlowBox title="Claims Portal">
            In: 837 claims
            <br />
            Out: acknowledgments, <span className="text-accent">835 ERA</span>
          </FlowBox>
          <FlowBox title="Internet Gateway">
            In: 276 status requests
            <br />
            Out: acknowledgments, <span className="text-accent">835 ERA</span>, 277 responses
          </FlowBox>
        </div>
        <div className="row-span-2 flex flex-col justify-around">
          <Two label="NSV" />
          <Two label="Internet" />
          <Two label="Internet" />
        </div>
        <div className="row-span-2 flex items-center justify-center border border-line bg-surface p-2 text-center text-[12px] font-semibold text-ink">
          One trading-partner mailbox for all inbound and outbound files
        </div>
        <div className="row-span-2 flex flex-col justify-around">
          <Two label="" />
          <Two label="" />
        </div>
        <FlowBox title="Inbound files">
          <ul className="list-disc pl-4">
            <li>Level I and II front-end edits</li>
            <li>837: return TA1, TRN, 999, 277CA</li>
            <li>Assign a CCN to accepted claims</li>
            <li>Split files by DME MAC, using the beneficiary state on the claim</li>
          </ul>
        </FlowBox>
        <div className="flex items-center"><Arrow dir="right" /></div>
        <FlowBox title="Receive and process">
          Receive claims and status requests. Process claims. <span className="text-accent">Create the 835 remittance advice flat file</span> and the 277 flat file.
        </FlowBox>

        <FlowBox title="Outbound files" strong>
          <span className="text-accent">Translate the remittance advice flat file to an X12 835 ERA</span> and deliver it to the mailbox. Translate the 277 flat file to X12 277.
        </FlowBox>
        <div className="flex items-center"><Arrow dir="left" /></div>
        <FlowBox title="Send to CEDI" strong>
          835 remittance advice flat file; 277 claim status response flat file.
        </FlowBox>
      </div>

      <div className="border-t border-line pt-4">
        <div className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-ink">Test flow, PDF p. 17: what differs</div>
        <div className="grid grid-cols-3 gap-2 text-[12px] leading-5 text-ink-2">
          <div className="border border-line p-3">Inbound test files receive the same front-end edits and acknowledgments. The figure shows no forwarding of test claims to a DME MAC.</div>
          <div className="border border-line p-3">DME MAC column: <span className="text-ink">create a TEST remittance advice flat file</span> and send it to CEDI.</div>
          <div className="border border-accent p-3">CEDI translates the test flat file to an X12 835 ERA and delivers it to the mailbox, the same final step as production.</div>
        </div>
      </div>
    </div>
  )
}

function Arrow({ dir }: { dir: 'left' | 'right' }) {
  const d = dir === 'right' ? 'M2 8 H34 M34 8 L28 3 M34 8 L28 13' : 'M34 8 H2 M2 8 L8 3 M2 8 L8 13'
  return (
    <svg width="36" height="16" viewBox="0 0 36 16" aria-hidden="true">
      <path d={d} stroke={ACCENT} strokeWidth="2" fill="none" />
    </svg>
  )
}

// ---------------------------------------------------------------------------
// Figure 2: the acknowledgment and remittance sequence

type Msg = { from: number; to: number; label: string; note?: string; accent?: boolean; self?: boolean }

const LANES = [
  { x: 110, name: 'Trading partner', sub: 'supplier, biller, clearinghouse' },
  { x: 390, name: 'CEDI', sub: 'front end and translator' },
  { x: 670, name: 'DME MAC', sub: 'jurisdictions A–D' },
]

const MSGS: Msg[] = [
  { from: 0, to: 1, label: '837P file (ISA … IEA)' },
  { from: 1, to: 0, label: 'TRN report', note: 'File recognized as X12; partner authorized' },
  { from: 1, to: 0, label: 'TA1', note: 'Only when the ISA/IEA envelope has an error' },
  { from: 1, to: 0, label: '999', note: 'X12 5010 syntax: accept or reject' },
  { from: 1, to: 0, label: '277CA', note: 'Per claim: accept, or reject with a status' },
  { from: 1, to: 2, label: 'Accepted claims', note: 'Routed by beneficiary state (2010BA N402)' },
  { from: 2, to: 2, label: 'Adjudicate', note: 'Create 835 remittance advice flat file', self: true },
  { from: 2, to: 1, label: '835 flat file', accent: true },
  { from: 1, to: 1, label: 'Translate', note: 'Flat file → X12 835 ERA', self: true, accent: true },
  { from: 1, to: 0, label: '835 ERA in mailbox', accent: true },
]

function AckSequence() {
  const top = 70
  const step = 38
  const height = top + MSGS.length * step + 20
  return (
    <svg viewBox={`0 0 780 ${height}`} className="w-full min-w-[680px]" role="img" aria-label="Sequence of files between trading partner, CEDI, and DME MAC">
      <defs>
        <marker id="cedi837p-seq" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill={INK2} />
        </marker>
        <marker id="cedi837p-seq-accent" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0 0 L10 5 L0 10 z" fill={ACCENT} />
        </marker>
      </defs>
      {LANES.map((lane) => (
        <g key={lane.name}>
          <rect x={lane.x - 85} y={8} width={170} height={42} fill={SURFACE} stroke={LINE} />
          <text x={lane.x} y={26} textAnchor="middle" fontSize="13" fontWeight="600" fill={INK}>{lane.name}</text>
          <text x={lane.x} y={42} textAnchor="middle" fontSize="10.5" fill={INK2}>{lane.sub}</text>
          <line x1={lane.x} y1={50} x2={lane.x} y2={height - 6} stroke={LINE} strokeDasharray="4 4" />
        </g>
      ))}
      {MSGS.map((m, i) => {
        const y = top + i * step
        const color = m.accent ? ACCENT : INK2
        const marker = m.accent ? 'url(#cedi837p-seq-accent)' : 'url(#cedi837p-seq)'
        if (m.self) {
          const x = LANES[m.from].x
          const left = m.from === 2
          const bx = left ? x - 150 : x + 8
          return (
            <g key={i}>
              <rect x={bx} y={y - 14} width={142} height={30} fill={SURFACE} stroke={color} />
              <text x={bx + 71} y={y - 1} textAnchor="middle" fontSize="11" fontWeight="600" fill={m.accent ? ACCENT : INK}>{m.label}</text>
              <text x={bx + 71} y={y + 11} textAnchor="middle" fontSize="9.5" fill={INK2}>{m.note}</text>
            </g>
          )
        }
        const x1 = LANES[m.from].x
        const x2 = LANES[m.to].x
        const mid = (x1 + x2) / 2
        return (
          <g key={i}>
            <line x1={x1} y1={y} x2={x2 + (x2 > x1 ? -3 : 3)} y2={y} stroke={color} strokeWidth={m.accent ? 2 : 1.3} markerEnd={marker} />
            <text x={mid} y={y - 5} textAnchor="middle" fontSize="11.5" fontWeight="600" fill={m.accent ? ACCENT : INK} fontFamily="var(--font-mono)">{m.label}</text>
            {m.note && <text x={mid} y={y + 13} textAnchor="middle" fontSize="10" fill={INK2}>{m.note}</text>}
          </g>
        )
      })}
    </svg>
  )
}

// ---------------------------------------------------------------------------
// Figure 3: envelope nesting with Medicare's required values (PDF pp. 23-25, 28, 30, 50)

function Rule({ el, value, note }: { el: string; value: string; note?: string }) {
  return (
    <div className="grid grid-cols-[4.5rem_9rem_1fr] gap-2 py-0.5 text-[12px] leading-5">
      <span className="font-mono text-ink">{el}</span>
      <span className="font-mono text-accent">{value}</span>
      <span className="text-ink-2">{note}</span>
    </div>
  )
}

function EnvelopeNesting() {
  return (
    <div className="min-w-[640px] border-2 border-ink-2 p-3">
      <div className="mb-1 flex items-baseline justify-between">
        <span className="font-mono text-[13px] font-semibold text-ink">ISA … IEA · interchange</span>
        <span className="text-[11px] text-ink-2">one per file transmission</span>
      </div>
      <Rule el="ISA01/03" value="00" note="Authorization and security qualifiers; ISA02 is 10 blank spaces; ISA04 content is ignored" />
      <Rule el="ISA05/07" value="27 or ZZ" note="Sender and receiver ID qualifiers" />
      <Rule el="ISA06" value="CEDI submitter ID" note="Same value required in GS02 and 1000A NM109" />
      <Rule el="ISA08" value="16013 · 17013 · 18003 · 19003" note="DME MAC jurisdiction A · B · C · D" />
      <Rule el="ISA11" value="submitter's choice" note="Repetition separator" />
      <Rule el="ISA14" value="1" note="TA1 requested; returned only on an envelope error" />
      <div className="mt-3 border-2 border-line p-3">
        <div className="mb-1 flex items-baseline justify-between">
          <span className="font-mono text-[13px] font-semibold text-ink">GS … GE · functional group</span>
          <span className="text-[11px] text-danger">more than one GS per ISA is rejected</span>
        </div>
        <Rule el="GS02" value="CEDI submitter ID" />
        <Rule el="GS03" value="16013 · 17013 · 18003 · 19003" note="DME MAC contractor ID" />
        <Rule el="GS04" value="date" note="Must not be a future date" />
        <Rule el="GS08" value="005010X222A1" note="Version identifier Medicare expects" />
        <div className="mt-3 border border-accent p-3">
          <div className="mb-1 flex items-baseline justify-between">
            <span className="font-mono text-[13px] font-semibold text-ink">ST … SE · transaction set (837P)</span>
            <span className="text-[11px] text-danger">one line of business per ST–SE</span>
          </div>
          <Rule el="ST02" value="control number" note="Must be unique within the interchange, or CEDI rejects it" />
          <Rule el="BHT02" value="00" note="Original" />
          <Rule el="BHT06" value="CH" note="Chargeable" />
          <div className="my-2 border border-dashed border-line px-3 py-2 text-[12px] text-ink-2">
            Submitter, receiver, billing provider, subscriber, claims (CLM), service lines (SV1) — see Figure 5
          </div>
          <Rule el="SE02" value="= ST02" note="Must be greater than zero" />
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 4: delimiters, read from the ISA header (PDF p. 25; examples pp. 50-51)

function Delimiters() {
  const isa: { t: string; mark?: string }[] = [
    { t: 'ISA' }, { t: '*', mark: 'a' }, { t: '00' }, { t: '*' }, { t: '␣×10' }, { t: '*' }, { t: '00' }, { t: '*' }, { t: '␣×10' },
    { t: '*' }, { t: 'ZZ' }, { t: '*' }, { t: 'B08XXXXXX' }, { t: '*' }, { t: 'ZZ' }, { t: '*' }, { t: '17013' }, { t: '*' },
    { t: '190123' }, { t: '*' }, { t: '1100' }, { t: '*' }, { t: '^', mark: 'b' }, { t: '*' }, { t: '00501' }, { t: '*' },
    { t: '067660001' }, { t: '*' }, { t: '1' }, { t: '*' }, { t: 'P' }, { t: '*' }, { t: '>', mark: 'c' }, { t: '~', mark: 'd' },
  ]
  const legend: Record<string, string> = {
    a: 'Element separator: the character right after "ISA"',
    b: 'Repetition separator, ISA11: submitter defined',
    c: 'Component separator, ISA16: splits composites such as SV1*HC>E0627>RR',
    d: 'Segment terminator: ends every segment',
  }
  return (
    <div className="min-w-[680px] space-y-4">
      <div className="flex flex-wrap font-mono text-[13px] leading-8">
        {isa.map((p, i) => (
          <span key={i} className={p.mark ? 'relative mx-px border-2 border-accent bg-surface px-1 text-accent' : p.t.startsWith('␣') ? 'px-0.5 text-ink-2' : 'text-ink'}>
            {p.t}
            {p.mark && <sup className="ml-0.5 font-sans text-[10px]">{p.mark}</sup>}
          </span>
        ))}
      </div>
      <div className="grid grid-cols-2 gap-x-6 gap-y-1 text-[12px] text-ink-2">
        {Object.entries(legend).map(([k, v]) => (
          <div key={k}><span className="font-semibold text-accent">{k}</span> · {v}</div>
        ))}
      </div>
      <div className="grid grid-cols-3 gap-px border border-line bg-line text-[12px] leading-5">
        <div className="bg-surface p-3"><div className="mb-1 font-semibold text-ink">Inbound 837P</div>Chosen by the submitter. Any characters the TR3 allows, provided they do not occur inside data within the interchange.</div>
        <div className="bg-surface p-3"><div className="mb-1 font-semibold text-ink">Outbound TA1</div>Reuses the delimiters of the submitted file.</div>
        <div className="bg-surface p-3"><div className="mb-1 font-semibold text-ink">Outbound 999</div>Standard delimiters when nothing is rejected at the functional-group level; otherwise the submitted file's delimiters. The guide's sample 999 ends its ISA with <span className="font-mono text-ink">:~</span>, not <span className="font-mono text-ink">&gt;~</span>.</div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 5: the 837P loop tree with Medicare-specific rules (PDF pp. 26-27, 30-50)

type LoopNode = { id: string; name: string; depth: number; rules?: string[]; removed?: boolean; page: string }

const TREE: LoopNode[] = [
  { id: '1000A', name: 'Submitter', depth: 0, page: '31', rules: ['NM109 = CEDI submitter ID, same as ISA06 and GS02'] },
  { id: '1000B', name: 'Receiver', depth: 0, page: '32', rules: ['NM109 = valid DME MAC contractor code'] },
  { id: '2000A', name: 'Billing provider level', depth: 0, page: '26, 32', rules: ['CUR foreign currency not supported'] },
  { id: '2010AA', name: 'Billing provider name', depth: 1, page: '32', rules: ['Must be associated with the submitter', 'NPI required; legacy IDs rejected (p. 28)'] },
  { id: '2010AC', name: 'Pay-to plan', depth: 1, page: '26', removed: true },
  { id: '2000B', name: 'Subscriber level', depth: 1, page: '33', rules: ['HL04 = 0 (no child level)', 'SBR01 = P or S · SBR02 = 18 (subscriber is the patient) · SBR09 = MB'] },
  { id: '2010BA', name: 'Subscriber name', depth: 2, page: '34–35', rules: ['NM108 = MI · NM109 = Medicare Beneficiary Identifier', 'N402 state picks the DME MAC'] },
  { id: '2010BB', name: 'Payer name', depth: 2, page: '35', rules: ['NM108 = PI'] },
  { id: '2000C / 2010CA', name: 'Patient level', depth: 2, page: '26–27', removed: true, rules: ['Subscriber is always the patient'] },
  { id: '2300', name: 'Claim (CLM)', depth: 2, page: '36–38', rules: ['CLM01: 38 allowed, 20 stored and returned', 'CLM02 = sum of SV102 when Medicare is primary', 'CLM05-3 = 1 (original)'] },
  { id: '2320 / 2330', name: 'Other payer (COB)', depth: 3, page: '40–41', rules: ['Required when SBR01 = S: AMT01 = D', 'CAS not allowed when SBR01 = P', '2330C–2330G loops not used'] },
  { id: '2400', name: 'Service line (SV1)', depth: 3, page: '42–46', rules: ['SV101-1 = HC · SV103 = UN for DME', 'SV102 > 0 · SV104 ≤ 9999.9', 'Capped rental, 2 units: RT and LT modifiers'] },
  { id: '2410', name: 'Drug (LIN/CTP)', depth: 4, page: '46–47', rules: ['LIN02 = N4 · LIN03 = 11-position NDC'] },
  { id: '2420A/D/E/F', name: 'Line providers', depth: 4, page: '47–48', rules: ['NM105 first character A–Z'] },
  { id: '2430', name: 'Line adjudication (SVD)', depth: 4, page: '49', rules: ['SVD03-1 = HC · SVD05 ≤ 9999.9'] },
]

function LoopTree() {
  return (
    <div className="min-w-[700px] font-mono text-[12px]">
      {TREE.map((n) => (
        <div key={n.id} className="grid grid-cols-[1fr_22rem] items-start border-b border-line py-1.5">
          <div className="flex items-start" style={{ paddingLeft: `${n.depth * 1.5}rem` }}>
            {n.depth > 0 && <span className="mr-1 text-ink-2">└</span>}
            <span className={`mr-2 border px-1.5 ${n.removed ? 'border-danger text-danger line-through' : 'border-accent text-accent'}`}>{n.id}</span>
            <span className={`font-sans ${n.removed ? 'text-ink-2 line-through' : 'text-ink'}`}>{n.name}</span>
            <span className="ml-2 font-sans text-[10px] text-ink-2">p. {n.page}</span>
          </div>
          <ul className="font-sans text-[11.5px] leading-5 text-ink-2">
            {n.removed && <li className="text-danger">Must not be present</li>}
            {n.rules?.map((r) => <li key={r}>{r}</li>)}
          </ul>
        </div>
      ))}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 6: MBI byte ruler and numeric limits (PDF pp. 28, 35-49)

const MBI = ['C', 'A', 'AN', 'N', 'A', 'AN', 'N', 'A', 'A', 'N', 'N']
const MBI_KEY: Record<string, string> = {
  C: '1–9',
  A: 'A–Z except S L O I B Z',
  N: '0–9',
  AN: 'A or N',
}

const LIMITS: [string, string, string][] = [
  ['Any dollar amount', 'at most 7 characters; above 99,999.99 rejected', '28'],
  ['Any percentage', 'at most 99.99; two digits each side of the decimal', '28'],
  ['Amounts outside CAS', 'must be positive; negative rejects the claim', '28'],
  ['CLM02 total charge', '≥ 0 and ≤ 99,999.9 (as printed)', '36'],
  ['SV102 line charge', '> 0; 0, 1, or 2 decimal places', '42'],
  ['SV104 units', '> 0; ≤ 9999.9; 1 or 2 for capped rentals and pen pumps', '42'],
  ['SV503 quantity (DME)', '> 0; ≤ 999', '43'],
  ['PAT08 patient weight (DME)', '4 whole digits, 2 decimals', '34'],
  ['CTP04 drug units', '> 0; ≤ 9,999,999.999', '46'],
]

function MbiAndLimits() {
  const w = 52
  return (
    <div className="min-w-[640px] space-y-5">
      <svg viewBox={`0 0 ${MBI.length * w + 20} 96`} className="w-full max-w-[620px]" role="img" aria-label="Medicare Beneficiary Identifier: 11 positions">
        {MBI.map((c, i) => (
          <g key={i}>
            <rect x={10 + i * w} y={22} width={w} height={40} fill={c === 'C' ? SURFACE : 'none'} stroke={c === 'AN' ? ACCENT : INK2} />
            <text x={10 + i * w + w / 2} y={16} textAnchor="middle" fontSize="10" fill={INK2}>{i + 1}</text>
            <text x={10 + i * w + w / 2} y={47} textAnchor="middle" fontSize="15" fontWeight="600" fill={c === 'AN' ? ACCENT : INK} fontFamily="var(--font-mono)">{c}</text>
          </g>
        ))}
        {Object.entries(MBI_KEY).map(([k, v], i) => (
          <text key={k} x={[10, 90, 290, 390][i]} y={84} fontSize="11" fill={INK2}>
            <tspan fontWeight="600" fill={INK} fontFamily="var(--font-mono)">{k}</tspan> = {v}
          </text>
        ))}
      </svg>
      <table className="w-full text-left text-[12px]">
        <thead>
          <tr className="border-b border-line text-ink"><th className="py-1.5 pr-4">Value</th><th className="py-1.5 pr-4">Medicare limit on the 837P</th><th className="py-1.5">PDF p.</th></tr>
        </thead>
        <tbody>
          {LIMITS.map(([v, l, p]) => (
            <tr key={v} className="border-b border-line">
              <td className="py-1.5 pr-4 text-ink">{v}</td>
              <td className="py-1.5 pr-4 font-mono text-[11.5px] text-ink-2">{l}</td>
              <td className="py-1.5 text-ink-2">{p}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 7: the guide's example 837P file, one segment per row (PDF p. 50)

type Seg = { s: string; loop?: string; note?: string; hot?: boolean }

const EXAMPLE: Seg[] = [
  { s: 'ISA*00*␣×10*00*␣×10*ZZ*B08XXXXXX*ZZ*17013*190123*1100*^*00501*067660001*1*P*>', loop: 'ISA', note: 'ISA08 17013 = jurisdiction B; ISA14 = 1', hot: true },
  { s: 'GS*HC*B08XXXXXX*PCACEV2.18*20190123*1100*4001*X*005010X222A1', loop: 'GS', note: 'GS03 in the example is not one of the Table 4 codes' },
  { s: 'ST*837*000000001*005010X222A1', loop: 'ST' },
  { s: 'BHT*0019*00*000004*20190123*105958*CH', loop: 'BHT', note: '00 original · CH chargeable' },
  { s: 'NM1*41*2*EDIT TEST SUBMITTER*****46*B08XXXXXX', loop: '1000A', note: 'Submitter ID matches ISA06' },
  { s: 'PER*IC*CONTACT NAME*TE*5024441111', loop: '1000A' },
  { s: 'NM1*40*2*DME MAC JURISDICTION B*****46*17013', loop: '1000B', note: 'Receiver = DME MAC B' },
  { s: 'HL*1**20*1', loop: '2000A' },
  { s: 'NM1*85*2*BILLING PROVIDER*****XX*1234567893', loop: '2010AA', note: 'Billing provider NPI' },
  { s: 'N3*106 SOME RD*SUITE A', loop: '2010AA' },
  { s: 'N4*HENDERSONVILLE*TN*370753495', loop: '2010AA' },
  { s: 'REF*EI*123456789', loop: '2010AA' },
  { s: 'PER*IC*NAME*TE*6158228888', loop: '2010AA' },
  { s: 'HL*2*1*22*0', loop: '2000B', note: 'HL04 = 0: no patient level' },
  { s: 'SBR*P*18*******MB', loop: '2000B', note: 'P primary · 18 subscriber is the patient · MB', hot: true },
  { s: 'NM1*IL*1*LASTNAME*FIRSTNAME****MI*##########', loop: '2010BA', note: 'MI qualifier; the guide masks the MBI' },
  { s: 'N3*1234 SOME STREET', loop: '2010BA' },
  { s: 'N4*LOUISVILLE*KY*10168', loop: '2010BA', note: 'N402 = KY: routing state', hot: true },
  { s: 'DMG*D8*19270829*M', loop: '2010BA' },
  { s: 'NM1*PR*2*DME MAC JURISDICTION B*****PI*17013', loop: '2010BB', note: 'PI payer ID' },
  { s: 'N4*LOUISVILLE*KY*40223', loop: '2010BB' },
  { s: 'CLM*JB CLAIM 1*40***12>B>1*Y*A*Y*Y', loop: '2300', note: 'CLM01 patient control number · CLM02 $40 · CLM05-3 = 1', hot: true },
  { s: 'HI*BK>71598', loop: '2300', note: 'Diagnosis' },
  { s: 'LX*1', loop: '2400' },
  { s: 'SV1*HC>E0627>RR*40*UN*1***1', loop: '2400', note: 'HC · procedure E0627 · modifier RR · $40 = CLM02 · UN · 1 unit', hot: true },
  { s: 'DTP*472*D8*20181202', loop: '2400', note: 'Date of service (472)' },
  { s: 'REF*6R*001H000000001403092010', loop: '2400' },
  { s: 'NM1*DK*1*PROVIDER*ORDERING****XX*1234567893', loop: '2420E', note: 'Ordering provider' },
  { s: 'N3*2425 STREET', loop: '2420E' },
  { s: 'N4*CITY*KY*40508', loop: '2420E' },
  { s: 'SE*29*000000001', loop: 'SE', note: '29 segments ST through SE; SE02 = ST02' },
  { s: 'GE*1*4001', loop: 'GE', note: '1 transaction; GE02 = GS06' },
  { s: 'IEA*1*067660001', loop: 'IEA', note: '1 group; IEA02 = ISA13' },
]

function ExampleFile() {
  return (
    <div className="min-w-[760px]">
      {EXAMPLE.map((e, i) => (
        <div key={i} className={`grid grid-cols-[4.5rem_1fr_18rem] gap-3 border-b border-line py-0.5 text-[11.5px] leading-5 ${e.hot ? 'bg-surface' : ''}`}>
          <span className="font-mono text-ink-2">{e.loop}</span>
          <span className="break-all font-mono text-ink">{e.s}<span className="text-accent">~</span></span>
          <span className={e.hot ? 'text-accent' : 'text-ink-2'}>{e.note}</span>
        </div>
      ))}
      <div className="mt-3 grid grid-cols-3 gap-2 text-[12px] leading-5 text-ink-2">
        <div className="border border-line p-3"><span className="font-semibold text-ink">TRN report</span>: file size 994, “No input validation problems”, 1 envelope processed out of 1 (p. 51).</div>
        <div className="border border-line p-3"><span className="font-semibold text-ink">999</span>: <span className="font-mono">AK9*A*1*1*1</span>, the group was accepted (p. 51).</div>
        <div className="border border-danger p-3"><span className="font-semibold text-ink">277CA</span>: claim <span className="font-mono">JB CLAIM 1</span> rejected, <span className="font-mono">STC*A7:254</span>, for the date of service billed with this procedure code (p. 51).</div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------

const WHERE: [string, string][] = [
  ['Scope and overview', '9–10'],
  ['Trading-partner roles and enrollment', '11–14'],
  ['Testing: 100% syntax, 95% data accuracy', '15'],
  ['Process flows (Figures 1 and 2)', '16–17'],
  ['Re-transmission: all claims marked original', '17'],
  ['Gateways, mailbox, X.509 certificates', '17–20'],
  ['DME MAC jurisdictions and states', '21–22'],
  ['Section 6: ISA and GS tables', '23–25'],
  ['Section 6.1: delimiters, inbound and outbound', '25'],
  ['Segments and loops not accepted by Medicare', '26–27'],
  ['General transaction notes: amounts, NPI, one GS per ISA', '28'],
  ['TA1, 999, TRN report', '28–29'],
  ['Header, submitter, receiver', '30–32'],
  ['Subscriber, MBI format, payer', '33–35'],
  ['Claim (2300), dates, diagnoses', '36–38'],
  ['Other payer / COB (2320–2330B)', '40–41'],
  ['Service line (2400), drug (2410), line adjudication (2430)', '42–49'],
  ['SE trailer; example 837P file', '50'],
  ['Example TRN, 999, and 277CA', '51–52'],
  ['Acronyms; version history', '53–56'],
]

export function FormatCedi837p() {
  return (
    <div className="mx-auto max-w-[90rem] px-10 py-12">
      <PageTitle lede="CMS Standard Companion Guide: Health Care Claim: Professional (837P). Companion Guide version 5.6, January 2025, based on ASC X12N TR3 005010X222A1. Published by National Government Services for CEDI.">
        CEDI 837P Companion Guide
      </PageTitle>

      <DocSection id="gist" title="The gist">
        <p>
          This guide is the rulebook for sending professional claims (the X12 <strong>837P</strong> transaction)
          to Medicare's four <strong>durable medical equipment Medicare Administrative Contractors</strong> (DME
          MACs) through <strong>CEDI</strong>, the Common Electronic Data Interchange front end. Suppliers, billing
          services, clearinghouses, and software vendors use it on top of the X12 technical report (TR3): it adds
          Medicare's envelope values, delimiter handling, and value limits.
        </p>
        <p>
          For the validator it matters because it shows where a DME 835 comes from: the MAC adjudicates and writes
          a <strong>remittance advice flat file</strong>, and CEDI translates that file into the X12 835
          electronic remittance advice (ERA) the supplier receives (PDF p. 16).
        </p>
        <p className="text-sm">
          This is a companion guide for exchanging
          X12, not the internal fixed-position flat-file layout.
        </p>
      </DocSection>

      <DocSection id="flow" title="How files move">
        <Figure number={1} caption="Production transaction flows (PDF p. 16) with the test-flow difference (PDF p. 17). NSV = network service vendor; sFTP = secure file transfer protocol; NCPDP = National Council for Prescription Drug Programs; CCN = claim control number. The 835 path is in the accent color.">
          <ProductionFlow />
        </Figure>
        <Figure number={2} caption="One 837P file through CEDI, in order. TRN is a CEDI text report; TA1 is the interchange acknowledgment; 999 is the implementation acknowledgment; 277CA is the claim acknowledgment. Sources: PDF pp. 16, 28–29, 50–52.">
          <AckSequence />
        </Figure>
      </DocSection>

      <DocSection id="envelopes" title="Envelopes and delimiters (section 6)">
        <Figure number={3} caption="Envelope nesting with the values Medicare requires. ISA/IEA = interchange control header/trailer; GS/GE = functional group; ST/SE = transaction set. Sources: PDF pp. 23–25 (Tables 3–4), 28 (one GS per ISA, one line of business per ST), 30 and 50 (ST02, BHT, SE02).">
          <EnvelopeNesting />
        </Figure>
        <Figure number={4} caption="Delimiters are whatever characters the ISA header carries in its set positions. The ISA shown is the guide's example (PDF p. 50); ISA02 holds 10 blank spaces and Medicare ignores ISA04 (PDF p. 23). Outbound rules from section 6.1 (PDF p. 25); sample 999 on PDF p. 51.">
          <Delimiters />
        </Figure>
      </DocSection>

      <DocSection id="medicare-rules" title="Medicare-specific data requirements">
        <p>Section 10 lists requirements beyond the TR3. Most are fixed values, size limits, or loops Medicare forbids. Errors found before the subscriber loop (2000B) fail the whole file at that point (PDF p. 26).</p>
        <Figure number={5} caption="837P loop tree with Medicare's added rules. Struck-through loops must not be present. COB = coordination of benefits; NPI = National Provider Identifier; NDC = National Drug Code; RT/LT = right and left modifiers. Page numbers are PDF pages.">
          <LoopTree />
        </Figure>
        <Figure number={6} caption="Top: the 11-position Medicare Beneficiary Identifier (MBI) required in 2010BA NM109 (PDF p. 35). Bottom: numeric limits; values beyond them cause rejection (PDF pp. 28, 34, 36, 42–43, 46).">
          <MbiAndLimits />
        </Figure>
        <Figure number={7} caption="The guide's example 837P file (PDF p. 50), split at each ~ segment terminator and labeled by loop. Highlighted rows carry Medicare-specific values. Below: the acknowledgments the guide shows for this file (PDF p. 51).">
          <ExampleFile />
        </Figure>
      </DocSection>

      <DocSection id="validator" title="What this means for the validator">
        <ul className="list-disc space-y-2 pl-5">
          <li>
            A DME 835 is produced in two steps: the MAC writes a remittance advice flat file, and CEDI translates it
            to X12 (PDF p. 16). A difference in a compared 835 can originate in adjudication or in translation;
            record which stage each reference file represents.
          </li>
          <li>
            Claims reach a DME MAC by the beneficiary's state in 2010BA N402, not by the envelope receiver ID
            (PDF pp. 24, 32). Pair files by jurisdiction on that basis.
          </li>
          <li>
            Medicare stores and returns only 20 characters of CLM01, the patient control number (PDF p. 36).
            A returned patient control number longer than 20 characters would not match that behavior.
          </li>
          <li>
            CEDI assigns a claim control number (CCN) to accepted claims (PDF p. 16). It is a candidate
            cross-system claim key for the matching question, not a settled one.
          </li>
          <li>
            Delimiters vary by file, and CEDI's outbound files may use different ones from the inbound file
            (PDF p. 25). Read them from each ISA; never assume <code>*</code>, <code>~</code>, or <code>&gt;</code>.
          </li>
          <li>
            Test remittances come from a separate test flat file (PDF p. 17). Keep test and production outputs
            apart when building comparison pairs.
          </li>
        </ul>
      </DocSection>

      <DocSection id="where" title="Where to look in the PDF">
        <table className="w-full max-w-[48rem] text-left text-sm">
          <thead>
            <tr className="border-b border-line text-xs text-ink"><th className="py-2 pr-5">Topic</th><th className="py-2">PDF pages</th></tr>
          </thead>
          <tbody>
            {WHERE.map(([t, p]) => (
              <tr key={t} className="border-b border-line"><td className="py-2 pr-5">{t}</td><td className="py-2 font-mono text-xs">{p}</td></tr>
            ))}
          </tbody>
        </table>
        <p className="text-sm">PDF viewer pages count from the cover. Printed page numbers are 8 lower from section 1 onward.</p>
      </DocSection>
    </div>
  )
}
