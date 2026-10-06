import type { ReactNode } from 'react'
import { DocSection } from '../../../components/DocSection'
import { Figure } from '../../../components/Figure'
import { PageTitle } from '../../../components/PageTitle'

// Boils down the CEDI 835 companion guide
// (CEDI 835 Standard Companion Guide, version 5.6, January 2025, 005010X221A1).
// Every fact cites a PDF viewer page number, counted from 1.

const MONO = 'IBM Plex Mono, monospace'

// ---------------------------------------------------------------------------
// Figure 1: the production flow, drawn from the guide's Figure 1 (PDF p. 13).

function Box({ x, y, w, h, title, lines, strong }: { x: number; y: number; w: number; h: number; title: string; lines: string[]; strong?: boolean }) {
  return (
    <g>
      <rect x={x} y={y} width={w} height={h} fill="var(--color-surface)" stroke={strong ? 'var(--color-accent)' : 'var(--color-line)'} strokeWidth={strong ? 1.8 : 1} />
      <text x={x + 12} y={y + 22} fontFamily={MONO} fontSize="12.5" fontWeight={600} fill={strong ? 'var(--color-accent)' : 'var(--color-ink)'}>{title}</text>
      {lines.map((line, i) => (
        <text key={line} x={x + 12} y={y + 42 + i * 16} fontFamily={MONO} fontSize="11" fill="var(--color-ink-2)">{line}</text>
      ))}
    </g>
  )
}

function ProductionFlow() {
  return (
    <svg viewBox="0 0 980 450" width="100%" role="img" style={{ minWidth: 760 }}>
      <title>
        Claims flow from trading partners through the CEDI mailbox and front-end edits to the DME MACs. Remittance flows back: the DME MAC creates an 835 remittance advice flat file, CEDI translates it to an X12 835 and places it in the trading partner mailbox, and the partner retrieves it through one of three gateways.
      </title>
      <defs>
        <marker id="cedi835-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0,0 L10,5 L0,10 z" fill="var(--color-ink-2)" />
        </marker>
        <marker id="cedi835-arrow-accent" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0,0 L10,5 L0,10 z" fill="var(--color-accent)" />
        </marker>
      </defs>

      {/* Lane headers */}
      {[
        { x: 20, w: 220, label: 'Trading partners' },
        { x: 260, w: 480, label: 'CEDI processing' },
        { x: 760, w: 200, label: 'DME MACs' },
      ].map((lane) => (
        <g key={lane.label}>
          <rect x={lane.x} y={14} width={lane.w} height={28} fill="none" stroke="var(--color-line)" />
          <text x={lane.x + lane.w / 2} y={33} textAnchor="middle" fontFamily={MONO} fontSize="12.5" fontWeight={600} fill="var(--color-ink)">{lane.label}</text>
        </g>
      ))}

      {/* Row labels */}
      <text x={20} y={66} fontFamily={MONO} fontSize="11" fill="var(--color-ink-2)">inbound: claims →</text>
      <text x={20} y={252} fontFamily={MONO} fontSize="11" fill="var(--color-accent)">← outbound: remittance (the 835 path)</text>

      {/* Inbound row */}
      <Box x={20} y={76} w={220} h={120} title="Submit claims" lines={['X12 837 claims', 'NCPDP claims', 'X12 276 status requests', '(sFTP gateway or', ' claims portal)']} />
      <Box x={375} y={76} w={365} h={120} title="Inbound files" lines={['Level I & II front-end edits', 'return TA1, TRN, 999, 277CA', 'assign CCN to accepted claims', 'route to a DME MAC by the', "beneficiary's state code"]} />
      <Box x={760} y={76} w={200} h={120} title="Process claims" lines={['adjudicate claims', 'process claim status', 'requests']} />

      {/* Mailbox column */}
      <rect x={265} y={76} width={90} height={310} fill="var(--color-bg)" stroke="var(--color-ink-2)" strokeDasharray="4 3" />
      <text x={310} y={216} textAnchor="middle" fontFamily={MONO} fontSize="12" fontWeight={600} fill="var(--color-ink)">Trading</text>
      <text x={310} y={232} textAnchor="middle" fontFamily={MONO} fontSize="12" fontWeight={600} fill="var(--color-ink)">Partner</text>
      <text x={310} y={248} textAnchor="middle" fontFamily={MONO} fontSize="12" fontWeight={600} fill="var(--color-ink)">Mailbox</text>
      <text x={310} y={272} textAnchor="middle" fontFamily={MONO} fontSize="10" fill="var(--color-ink-2)">835 kept</text>
      <text x={310} y={286} textAnchor="middle" fontFamily={MONO} fontSize="10" fill="var(--color-ink-2)">45 days</text>

      <line x1={240} y1={120} x2={263} y2={120} stroke="var(--color-ink-2)" markerEnd="url(#cedi835-arrow)" />
      <line x1={355} y1={120} x2={373} y2={120} stroke="var(--color-ink-2)" markerEnd="url(#cedi835-arrow)" />
      <line x1={740} y1={120} x2={758} y2={120} stroke="var(--color-ink-2)" markerEnd="url(#cedi835-arrow)" />

      {/* Outbound row */}
      <Box x={760} y={266} w={200} h={120} title="① Create flat file" lines={['835 Remittance', 'Advice Flat File', '(MAC internal layout,', ' not defined here)', 'sent to CEDI']} />
      <Box x={375} y={266} w={365} h={120} strong title="② Translate flat → X12 835" lines={['Remittance Advice Flat File', '→ X12 835 ERA', 'deliver to the trading', 'partner mailbox']} />
      <Box x={20} y={266} w={220} h={120} title="③ Retrieve the 835" lines={['sFTP gateway  (via NSV)', 'claims portal (Internet)', 'Internet gateway', ' (CAQH CORE, X.509 cert)']} />

      <line x1={760} y1={326} x2={742} y2={326} stroke="var(--color-accent)" strokeWidth={1.8} markerEnd="url(#cedi835-arrow-accent)" />
      <line x1={375} y1={326} x2={357} y2={326} stroke="var(--color-accent)" strokeWidth={1.8} markerEnd="url(#cedi835-arrow-accent)" />
      <line x1={265} y1={326} x2={242} y2={326} stroke="var(--color-accent)" strokeWidth={1.8} markerEnd="url(#cedi835-arrow-accent)" />

      <text x={20} y={418} fontFamily={MONO} fontSize="11" fill="var(--color-ink-2)">Test flow (PDF p. 14): the DME MAC creates a TEST remittance advice flat file; CEDI translates and delivers it the same way.</text>
      <text x={20} y={436} fontFamily={MONO} fontSize="11" fill="var(--color-ink-2)">Any submission path (837P, NCPDP D.0, or paper) gets its remittance back as an X12 835 (PDF p. 22).</text>
    </svg>
  )
}

// ---------------------------------------------------------------------------
// Figure 2: the translation boundary (PDF pp. 13, 30).

function TranslationBoundary() {
  const x12 = ['ISA*00*…*27*17013*29*B08XXXXXX…~', 'GS*HP*17013*B08XXXXXX*…*005010X221A1~', 'ST*835*0001~', 'BPR*I*72.59*C*CHK*…~', 'CLP*TC15*1*100*73.42*18.35*MB*…~', 'SVC*HC:E0650:RR*100*73.42**1~', '…', 'SE*42*0001~']
  return (
    <div className="grid items-stretch gap-4 md:grid-cols-[1fr_auto_1.3fr]">
      <div className="border border-line bg-surface p-4">
        <div className="mb-2 font-mono text-[12px] font-semibold text-ink">DME MAC processing system</div>
        <div className="mb-3 font-mono text-[11px] text-ink-2">835 Remittance Advice Flat File</div>
        <div className="space-y-1">
          {['record', 'record', 'record', 'record', 'record'].map((label, i) => (
            <div key={i} className="flex h-5 border border-line bg-bg">
              <div className="w-1/4 border-r border-line bg-line/60" />
              <div className="w-1/6 border-r border-line" />
              <div className="w-1/3 border-r border-line" />
              <div className="flex-1" />
              <span className="sr-only">{label}</span>
            </div>
          ))}
        </div>
        <p className="mt-3 text-[12px] leading-5 text-ink-2">Fixed-position records. This guide names the file but does <strong className="text-ink">not</strong> define its layout.</p>
      </div>
      <div className="flex flex-col items-center justify-center px-2 text-center">
        <div className="font-mono text-[12px] font-semibold text-accent">CEDI translates</div>
        <svg viewBox="0 0 80 20" width="80" height="20" aria-hidden="true"><path d="M0,10 H70 M62,3 L72,10 L62,17" fill="none" stroke="var(--color-accent)" strokeWidth="2" /></svg>
        <div className="font-mono text-[11px] text-ink-2">into standard X12</div>
      </div>
      <div className="border border-accent bg-surface p-4">
        <div className="mb-2 font-mono text-[12px] font-semibold text-ink">X12 835 ERA in the mailbox</div>
        <pre className="overflow-x-auto font-mono text-[11px] leading-5 text-ink">{x12.join('\n')}</pre>
        <p className="mt-3 text-[12px] leading-5 text-ink-2">Delimited segments, from the guide's example file (PDF p. 30).</p>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 3: envelope nesting of the guide's example file (PDF p. 30).

function Seg({ children }: { children: ReactNode }) {
  return <code className="whitespace-nowrap bg-surface px-1.5 py-0.5 font-mono text-[11.5px] text-ink">{children}</code>
}

function EnvelopeNesting() {
  return (
    <div className="min-w-[640px] font-mono text-[12px]">
      <div className="border-2 border-ink-2 p-3">
        <div className="mb-2 flex flex-wrap items-baseline justify-between gap-2">
          <span className="font-semibold text-ink">Interchange (ISA … IEA)</span>
          <span className="text-ink-2">sender 17013 (DME MAC JB) → receiver B08XXXXXX (CEDI submitter number)</span>
        </div>
        <div className="border-2 border-accent p-3">
          <div className="mb-2 flex flex-wrap items-baseline justify-between gap-2">
            <span className="font-semibold text-accent">Functional group (GS … GE)</span>
            <span className="text-ink-2">GS01 HP · GS08 005010X221A1</span>
          </div>
          <div className="grid gap-3 md:grid-cols-2">
            {[
              {
                st: '0001',
                bpr: 'BPR*I*72.59*C*CHK',
                claims: ['CLP*TC15 · 1 service line', 'CLP*TC9 · 1 service line'],
                plb: 'PLB · 2 provider adjustments',
                se: 'SE*42*0001',
              },
              {
                st: '0002',
                bpr: 'BPR*H*0*C*NON',
                claims: ['CLP*TC35 · 2 service lines'],
                plb: 'no PLB',
                se: 'SE*31*0002',
              },
            ].map((t) => (
              <div key={t.st} className="border border-line bg-bg p-3">
                <div className="mb-2 font-semibold text-ink">Transaction ST*835*{t.st}</div>
                <div className="space-y-1.5">
                  <div><Seg>{t.bpr}</Seg> <span className="text-ink-2">one payment</span></div>
                  <div className="text-ink-2">payer N1*PR, <Seg>REF*2U*17013</Seg>, payee N1*PE</div>
                  <div className="border-l-2 border-line pl-2">
                    <div className="text-ink-2">LX*1 (Loop 2000)</div>
                    {t.claims.map((c) => <div key={c} className="pl-3 text-ink">{c}</div>)}
                  </div>
                  <div className="text-ink-2">{t.plb}</div>
                  <div><Seg>{t.se}</Seg> <span className="text-ink-2">segment count, control number</span></div>
                </div>
              </div>
            ))}
          </div>
          <div className="mt-2 text-ink-2"><Seg>GE*2*1</Seg> two transactions in this group</div>
        </div>
        <div className="mt-2 text-ink-2"><Seg>IEA*1*000000001</Seg> one group; matches ISA13</div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 4: ISA and GS element by element (PDF pp. 20-21, example p. 30).

type Kind = 'fixed' | 'route' | 'open'
const KIND_STYLE: Record<Kind, string> = {
  fixed: 'border-accent text-accent',
  route: 'border-ink-2 text-ink',
  open: 'border-line text-ink-2 border-dashed',
}

const ISA: { el: string; name: string; value: string; kind: Kind; rule?: string }[] = [
  { el: '01', name: 'auth qualifier', value: '00', kind: 'fixed', rule: 'sends 00' },
  { el: '02', name: 'auth info', value: '␣×10', kind: 'fixed', rule: '10 spaces' },
  { el: '03', name: 'security qualifier', value: '00', kind: 'fixed', rule: 'sends 00' },
  { el: '04', name: 'security info', value: '␣×10', kind: 'fixed', rule: 'spaces' },
  { el: '05', name: 'sender qualifier', value: '27', kind: 'fixed', rule: 'sends 27' },
  { el: '06', name: 'sender ID', value: '17013', kind: 'route', rule: 'DME MAC ID' },
  { el: '07', name: 'receiver qualifier', value: '29', kind: 'fixed', rule: 'sends 29' },
  { el: '08', name: 'receiver ID', value: 'B08XXXXXX', kind: 'route', rule: 'CEDI submitter no.' },
  { el: '09', name: 'date', value: '190124', kind: 'open' },
  { el: '10', name: 'time', value: '1022', kind: 'open' },
  { el: '11', name: 'repetition sep.', value: '^', kind: 'fixed', rule: 'CEDI puts ^' },
  { el: '12', name: 'version', value: '00501', kind: 'open' },
  { el: '13', name: 'control no.', value: '000000001', kind: 'open' },
  { el: '14', name: 'ack requested', value: '0', kind: 'fixed', rule: 'sends 0' },
  { el: '15', name: 'usage', value: 'P', kind: 'open' },
  { el: '16', name: 'component sep.', value: ':', kind: 'open' },
]

const GS: { el: string; name: string; value: string; kind: Kind; rule?: string }[] = [
  { el: '01', name: 'functional ID', value: 'HP', kind: 'open' },
  { el: '02', name: 'app sender', value: '17013', kind: 'route', rule: 'DME MAC ID' },
  { el: '03', name: 'app receiver', value: 'B08XXXXXX', kind: 'route', rule: 'CEDI submitter no.' },
  { el: '04', name: 'date', value: '20190124', kind: 'open' },
  { el: '05', name: 'time', value: '10222270', kind: 'open' },
  { el: '06', name: 'group control', value: '1', kind: 'open' },
  { el: '07', name: 'agency', value: 'X', kind: 'open' },
  { el: '08', name: 'version', value: '005010X221A1', kind: 'fixed', rule: 'required' },
]

function ElementStrip({ seg, cells }: { seg: string; cells: typeof ISA }) {
  return (
    <div className="mb-5">
      <div className="mb-1.5 font-mono text-[12px] font-semibold text-ink">{seg}</div>
      <div className="flex flex-wrap gap-1">
        {cells.map((c) => (
          <div key={c.el} className={`min-w-[92px] border bg-surface px-2 py-1.5 font-mono ${KIND_STYLE[c.kind]}`}>
            <div className="text-[10px] text-ink-2">{seg}{c.el}</div>
            <div className="text-[13px] font-semibold">{c.value}</div>
            <div className="text-[10px] leading-4 text-ink-2">{c.name}</div>
            {c.rule && <div className="text-[10px] leading-4">{c.rule}</div>}
          </div>
        ))}
      </div>
    </div>
  )
}

function EnvelopeElements() {
  return (
    <div>
      <ElementStrip seg="ISA" cells={ISA} />
      <ElementStrip seg="GS" cells={GS} />
      <div className="flex flex-wrap gap-4 font-mono text-[11px] text-ink-2">
        <span><span className="mr-1 inline-block h-3 w-3 border border-accent align-middle" /> value set by the guide</span>
        <span><span className="mr-1 inline-block h-3 w-3 border border-ink-2 align-middle" /> routing identifier set by the guide</span>
        <span><span className="mr-1 inline-block h-3 w-3 border border-dashed border-line align-middle" /> guide says only: follow the TR3; value from the example</span>
      </div>
      <div className="mt-5 grid max-w-xl grid-cols-4 gap-px border border-line bg-line font-mono text-[12px]">
        {[['JA', '16013'], ['JB', '17013'], ['JC', '18003'], ['JD', '19003']].map(([j, id]) => (
          <div key={j} className="bg-surface p-2 text-center">
            <div className="text-ink-2">DME MAC {j}</div>
            <div className="text-[14px] font-semibold text-ink">{id}</div>
          </div>
        ))}
      </div>
      <div className="mt-1 font-mono text-[11px] text-ink-2">The same contractor ID appears in ISA06, GS02, and Loop 1000A REF02 (qualifier 2U).</div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 5: Medicare usage rules laid on the 835 structure (PDF pp. 23-29).

interface Rule {
  where: string
  element: string
  uses?: string[]
  notUsed?: string[]
  note?: string
  page: number
  segmentNotUsed?: boolean
}

const USAGE: { group: string; rows: Rule[] }[] = [
  {
    group: 'Header',
    rows: [
      { where: 'ST', element: 'ST02 control number', note: 'one-by-one counter starting at 0001', page: 23 },
      { where: 'BPR', element: 'BPR03 credit/debit flag', uses: ['C'], notUsed: ['D'], page: 24 },
      { where: 'BPR', element: 'BPR04 payment method', uses: ['ACH', 'CHK', 'NON'], notUsed: ['BOP', 'FWT'], page: 24 },
      { where: 'BPR', element: 'BPR06 / BPR12 bank ID qualifier', uses: ['01'], notUsed: ['04'], page: 24 },
      { where: '1000A REF', element: 'REF01 / REF02 payer ID', uses: ['2U', '16013', '17013', '18003', '19003'], page: 24 },
    ],
  },
  {
    group: 'Loop 2000',
    rows: [{ where: 'LX', element: 'LX01 assigned number', uses: ['1', '0'], note: '1 = assigned, 0 = non-assigned', page: 25 }],
  },
  {
    group: 'Loop 2100 claim',
    rows: [
      { where: 'CLP', element: 'CLP02 claim status', uses: ['1', '2', '3', '4', '19', '20', '21', '22', '23'], notUsed: ['25'], page: 25 },
      { where: 'CLP', element: 'CLP06 claim filing indicator', uses: ['MB'], note: 'MB for DME (MA also listed)', page: 25 },
      { where: 'CAS', element: 'CAS01 adjustment group', uses: ['CO', 'OA', 'PR'], notUsed: ['PI'], page: 25 },
      { where: 'NM1 patient', element: 'NM108 ID qualifier', uses: ['MI'], page: 25 },
      { where: 'NM1 insured', element: 'whole segment', segmentNotUsed: true, page: 26 },
      { where: 'NM1 crossover', element: 'NM108 ID qualifier', uses: ['PI', 'XV'], notUsed: ['AD', 'FI', 'NI', 'PP'], note: 'more than one secondary payer → claim remark N89', page: 26 },
      { where: 'REF other claim ID', element: 'REF01 qualifier', uses: ['28', '6P', 'EA', 'F8'], notUsed: ['1L', '1W', '9A', '9C', 'BB', 'CE', 'G1', 'G3', 'IG'], page: 26 },
      { where: 'REF rendering prov.', element: 'whole segment', segmentNotUsed: true, page: 26 },
      { where: 'AMT', element: 'AMT01 qualifier', uses: ['AU', 'DY', 'F5', 'I', 'NL', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'], notUsed: ['D8', 'T', 'T2'], page: 26 },
      { where: 'QTY', element: 'QTY01 qualifier', uses: ['CA', 'CD', 'LA', 'OU', 'ZK', 'ZL', 'ZM', 'ZN', 'ZO'], notUsed: ['LE', 'NE', 'NR', 'PS', 'VS'], page: 27 },
    ],
  },
  {
    group: 'Loop 2110 service line',
    rows: [
      { where: 'SVC', element: 'SVC01-1 / SVC06-1 code qualifier', uses: ['HC', 'NU', 'N4', 'HP'], note: 'only these four', page: 27 },
      { where: 'CAS', element: 'CAS01 adjustment group', uses: ['CO', 'OA', 'PR'], notUsed: ['PI'], page: 27 },
      { where: 'REF service ID', element: 'REF01 qualifier', uses: ['LU', '1S', 'APC', 'RB'], notUsed: ['BB', 'E9', 'G1', 'G3'], page: 27 },
      { where: 'REF rendering prov.', element: 'REF01 qualifier', uses: ['HPI', 'SY', 'TJ', '1C'], notUsed: ['0B', '1A', '1B', '1D', '1H', '1J', 'D3', 'G2'], page: 28 },
      { where: 'REF policy', element: 'REF01 qualifier', uses: ['0K'], note: 'REF02 carries the LCD/NCD code', page: 28 },
      { where: 'AMT', element: 'AMT01 qualifier', uses: ['B6', 'KH', '2K', 'ZL', 'ZM', 'ZN', 'ZO'], notUsed: ['T', 'T2'], page: 28 },
      { where: 'LQ', element: 'LQ01 remark list', uses: ['HE'], note: 'only HE', page: 28 },
    ],
  },
  {
    group: 'Summary',
    rows: [
      {
        where: 'PLB',
        element: 'PLB03-1 adjustment reason',
        uses: ['50', '51', '72', '90', 'AP', 'B2', 'B3', 'BD', 'BN', 'C5', 'CS', 'CV', 'DM', 'E3', 'FB', 'GO', 'HM', 'IP', 'IS', 'IR', 'J1', 'L3', 'L6', 'LE', 'LS', 'OA', 'OB', 'PI', 'PL', 'RA', 'RE', 'SL', 'TL', 'WO', 'WU'],
        notUsed: ['AH', 'AM', 'CR', 'CT', 'CW', 'FC'],
        page: 29,
      },
    ],
  },
]

function Chip({ code, off }: { code: string; off?: boolean }) {
  return (
    <span className={`inline-block border px-1.5 py-px font-mono text-[11px] leading-5 ${off ? 'border-line text-ink-2 line-through decoration-danger' : 'border-accent/60 bg-accent/10 text-ink'}`}>
      {code}
    </span>
  )
}

function UsageMap() {
  return (
    <div className="min-w-[760px]">
      <div className="mb-3 flex flex-wrap gap-4 font-mono text-[11px] text-ink-2">
        <span><Chip code="CO" /> Medicare uses</span>
        <span><Chip code="PI" off /> does not apply to Medicare</span>
        <span><span className="border border-danger px-1.5 py-px text-danger">not used</span> whole segment not used</span>
      </div>
      {USAGE.map((g, gi) => (
        <div key={g.group} className="flex">
          <div className="flex w-5 flex-col items-center">
            <div className={`w-px flex-1 ${gi === 0 ? 'bg-transparent' : 'bg-line'}`} />
            <div className="h-2.5 w-2.5 bg-accent" />
            <div className={`w-px flex-1 ${gi === USAGE.length - 1 ? 'bg-transparent' : 'bg-line'}`} />
          </div>
          <div className="flex-1 py-2 pl-3" style={{ paddingLeft: gi >= 2 && gi <= 3 ? 12 + (gi - 1) * 16 : 12 }}>
            <div className="mb-1 font-mono text-[12px] font-semibold text-ink">{g.group}</div>
            <div className="divide-y divide-line border border-line">
              {g.rows.map((r) => (
                <div key={`${r.where}-${r.element}`} className="grid grid-cols-[150px_200px_1fr_52px] items-start gap-2 bg-bg px-2 py-1.5 text-[12px]">
                  <div className="font-mono text-ink">{r.where}</div>
                  <div className="text-ink-2">{r.element}</div>
                  <div className="flex flex-wrap items-center gap-1">
                    {r.segmentNotUsed && <span className="border border-danger px-1.5 py-px font-mono text-[11px] text-danger">not used by Medicare</span>}
                    {r.uses?.map((c) => <Chip key={`u${c}`} code={c} />)}
                    {r.notUsed?.map((c) => <Chip key={`n${c}`} code={c} off />)}
                    {r.note && <span className="text-[11px] text-ink-2">{r.note}</span>}
                  </div>
                  <div className="text-right font-mono text-[11px] text-ink-2">p. {r.page}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

// ---------------------------------------------------------------------------
// Figure 6: the money in the guide's example balances (PDF p. 30).

function Bar({ billed, parts }: { billed: number; parts: { label: string; amount: number; tone: string }[] }) {
  return (
    <div className="flex h-7 w-full border border-line">
      {parts.map((p) => (
        <div key={p.label} className={`flex items-center justify-center overflow-hidden whitespace-nowrap font-mono text-[10.5px] ${p.tone}`} style={{ width: `${(p.amount / billed) * 100}%` }}>
          {p.label}
        </div>
      ))}
    </div>
  )
}

const PAID = 'bg-accent/25 text-ink'
const PR = 'bg-series-2/25 text-ink'
const CO = 'bg-line text-ink-2'

function MoneyBalance() {
  const claims = [
    { id: 'TC15', billed: 100, paid: 73.42, pr: 18.35, co: 8.23 },
    { id: 'TC9', billed: 450, paid: 150.88, pr: 37.72, co: 261.4 },
  ]
  return (
    <div className="grid gap-6 lg:grid-cols-2">
      <div>
        <div className="mb-2 font-mono text-[12px] font-semibold text-ink">Each claim: billed = paid + adjustments</div>
        {claims.map((c) => (
          <div key={c.id} className="mb-3">
            <div className="mb-1 flex justify-between font-mono text-[11px] text-ink-2">
              <span>CLP*{c.id} · billed {c.billed.toFixed(2)}</span>
              <span>{c.paid.toFixed(2)} + {c.pr.toFixed(2)} + {c.co.toFixed(2)} = {c.billed.toFixed(2)}</span>
            </div>
            <Bar
              billed={c.billed}
              parts={[
                { label: `paid ${c.paid}`, amount: c.paid, tone: PAID },
                { label: `PR ${c.pr}`, amount: c.pr, tone: PR },
                { label: `CO ${c.co}`, amount: c.co, tone: CO },
              ]}
            />
          </div>
        ))}
        <div className="font-mono text-[11px] leading-5 text-ink-2">
          paid = CLP04 · PR = patient responsibility, CLP05 and CAS*PR*2 · CO = contractual, CAS*CO*45 (service-level CAS)
        </div>
      </div>
      <div>
        <div className="mb-2 font-mono text-[12px] font-semibold text-ink">Transaction 0001: claims paid − provider adjustments = BPR02</div>
        <div className="grid grid-cols-[1fr_auto] gap-x-6 gap-y-1 border border-line bg-surface p-3 font-mono text-[12.5px]">
          <span className="text-ink-2">CLP04 claim TC15</span><span className="text-right text-ink">73.42</span>
          <span className="text-ink-2">CLP04 claim TC9</span><span className="text-right text-ink">+ 150.88</span>
          <span className="text-ink-2">PLB reason FB</span><span className="text-right text-ink">− 56.39</span>
          <span className="text-ink-2">PLB reason 90</span><span className="text-right text-ink">− 95.32</span>
          <span className="border-t border-line pt-1 font-semibold text-accent">BPR02 payment (CHK)</span>
          <span className="border-t border-line pt-1 text-right font-semibold text-accent">72.59</span>
        </div>
        <div className="mt-4 mb-2 font-mono text-[12px] font-semibold text-ink">Transaction 0002: nothing paid</div>
        <div className="border border-line bg-surface p-3 font-mono text-[12px] leading-6 text-ink-2">
          <div><span className="text-ink">BPR*H*0*C*NON</span> · payment 0.00, method NON</div>
          <div><span className="text-ink">CLP*TC35*4*300*0</span> · billed 300.00, paid 0.00</div>
          <div className="pl-3">SVC E1390 250 → CAS*CO*176*250</div>
          <div className="pl-3">SVC E0431 50 → CAS*CO*176*50</div>
        </div>
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------

const WHERE: [string, string][] = [
  ['Scope and purpose: supplements the TR3, does not replace it', '7'],
  ['Trading partner types and ERA enrollment (NPI and PTAN)', '10'],
  ['835 testing not required; recommended', '12'],
  ['Figure 1, production flow (flat file → X12 835)', '13'],
  ['Figure 2, test flow; 45-day mailbox retention; gateways', '14–15'],
  ['Who answers remittance-content questions (the DME MAC)', '18'],
  ['Section 6: ISA values (Table 3)', '20–21'],
  ['Section 6: GS values (Table 4); outbound delimiters', '21'],
  ['GS/GE and ST/SE follow the TR3; business rules 7.1–7.3', '22'],
  ['ST02 counter (Table 5)', '23'],
  ['BPR and Loop 1000A payer ID (Tables 6–7)', '24'],
  ['Loops 2000 and 2100 rules (Tables 8–17)', '25–27'],
  ['Loop 2110 rules (Tables 18–24)', '27–28'],
  ['PLB provider adjustment codes (Table 25)', '29'],
  ['Figure 3, complete example 835 file', '30'],
  ['Acronym list (Table 26)', '31–32'],
  ['Version history (Table 27)', '33'],
]

export function FormatCedi835() {
  return (
    <div className="mx-auto max-w-[90rem] px-10 py-12">
      <PageTitle
        lede={
          <>
            <strong className="text-ink">CEDI 835 Standard Companion Guide</strong>, version 5.6, January 2025, based on
            X12 005010X221A1.
          </>
        }
      >
        CEDI 835 Companion Guide
      </PageTitle>

      <DocSection id="gist" title="The gist">
        <p>
          This guide tells <strong>trading partners</strong> (suppliers, billing services, clearinghouses) how Medicare durable
          medical equipment (DME) remittance reaches them. The four DME Medicare Administrative Contractors (MACs) produce
          remittance data as an internal flat file; Common Electronic Data Interchange (CEDI), the DME front-end gateway,
          translates it into a standard X12 835 Electronic Remittance Advice (ERA) and places it in the partner's mailbox
          (PDF pp. 13, 30).
        </p>
        <p>
          For the validator, it defines what the legacy DME 835s look like: fixed envelope values, the contractor IDs in them, and
          which codes Medicare does and does not use. It covers DME only; it is not the flat-file layout and it
          does not repeat the full X12 Technical Report Type 3 (TR3) implementation guide (PDF p. 7).
        </p>
        <ul className="list-disc space-y-1 pl-5 text-[14px]">
          <li>Acronyms used below: NSV = Network Service Vendor; CCN = Claim Control Number; NCPDP = National Council for Prescription Drug Programs; CAQH CORE = Council for Affordable Quality Healthcare, Committee on Operating Rules for Information Exchange; TA1, 999, 277CA = acknowledgment transactions; TRN = CEDI's Transaction Acknowledgement report (PDF pp. 31–32).</li>
          <li>The guide's own tables have a "Page #" column. Those numbers point into the TR3, not into this PDF. Page numbers on this page are PDF viewer pages.</li>
        </ul>
      </DocSection>

      <DocSection id="flow" title="How DME remittance travels">
        <Figure number={1} caption="Production transaction flow, redrawn from the guide's Figure 1 (PDF p. 13); test flow and submission-method rule from PDF pp. 14 and 22. The 835 path is the lower row, right to left.">
          <ProductionFlow />
        </Figure>
        <Figure number={2} caption="The translation boundary. The DME MAC processing system generates the remittance content; CEDI translates it into the X12 835 format (PDF pp. 13, 30). Segment text is abbreviated from the guide's example (PDF p. 30).">
          <TranslationBoundary />
        </Figure>
      </DocSection>

      <DocSection id="envelopes" title="Envelopes (section 6)">
        <Figure number={3} caption="Nesting of the guide's complete example file: one interchange, one functional group, two 835 transactions, each with its own BPR payment line (PDF p. 30). GS/GE and ST/SE follow the TR3 (PDF p. 22).">
          <EnvelopeNesting />
        </Figure>
        <Figure number={4} caption="ISA and GS elements. Rules from Tables 3 and 4 (PDF pp. 20–21); values from the example file (PDF p. 30). Delimiters are not fixed by the guide: partners ask CEDI which to expect (PDF p. 21). The example uses * between elements, : inside composites, ~ after segments.">
          <EnvelopeElements />
        </Figure>
      </DocSection>

      <DocSection id="usage" title="What Medicare uses inside the 835">
        <p>Section 10 lists only the places where Medicare narrows the TR3. Everything not listed follows the TR3 as written.</p>
        <Figure number={5} caption="Medicare-specific 835 usage by loop, from Tables 5–25 (PDF pp. 23–29). Two segments are not used at all: 2100 NM1 Insured Name and 2100 REF Rendering Provider. LCD/NCD = Local or National Coverage Determination.">
          <UsageMap />
        </Figure>
      </DocSection>

      <DocSection id="example" title="The example file, balanced">
        <Figure number={6} caption="Amounts from the guide's example 835 (PDF p. 30). Each claim's paid amount plus its adjustments equals the billed amount; the transaction's payment equals the claims' paid amounts minus provider-level adjustments. Arithmetic checked by hand.">
          <MoneyBalance />
        </Figure>
      </DocSection>

      <DocSection id="validator" title="What this means for the validator">
        <ul className="list-disc space-y-2 pl-5">
          <li>Legacy DME 835 bytes are CEDI's translation of a DME MAC flat file (PDF pp. 13, 30). Differences can arise at adjudication or at translation, and the file alone does not say which.</li>
          <li>Read delimiters from each file's ISA segment. CEDI does not publish a fixed set; the example uses <code>*</code>, <code>:</code>, <code>~</code>, and <code>^</code> (PDF pp. 21, 30).</li>
          <li>One file can carry several 835 transactions, one BPR payment each (PDF p. 30). Compare claims within their transaction, and keep envelope and control numbers (ISA13, GS06, ST02) apart from claim content.</li>
          <li>The contractor ID (16013, 17013, 18003, 19003) appears in ISA06, GS02, and 1000A REF02 and identifies the DME jurisdiction of each legacy file (PDF pp. 20–21, 24).</li>
          <li>The code restrictions in Figure 5 can back a separate conformance check. For example, a CAS group code other than CO, OA, or PR would not come from Medicare. That check stays outside the exact comparison evidence.</li>
          <li>The payment and claim balance in Figure 6 gives a small synthetic fixture with known totals for testing.</li>
        </ul>
      </DocSection>

      <DocSection id="where" title="Where to look in the PDF">
        <div className="overflow-x-auto">
          <table className="w-full max-w-3xl border border-line text-[13px]">
            <thead>
              <tr className="bg-surface text-left text-ink">
                <th className="border-b border-line px-3 py-2 font-semibold">Topic</th>
                <th className="border-b border-line px-3 py-2 font-semibold">PDF page</th>
              </tr>
            </thead>
            <tbody>
              {WHERE.map(([topic, page]) => (
                <tr key={topic} className="border-b border-line last:border-b-0">
                  <td className="px-3 py-1.5">{topic}</td>
                  <td className="px-3 py-1.5 font-mono text-ink">{page}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DocSection>
    </div>
  )
}
