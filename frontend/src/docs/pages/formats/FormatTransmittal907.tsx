import { DocSection } from '../../../components/DocSection'
import { Figure } from '../../../components/Figure'
import { PageTitle } from '../../../components/PageTitle'

// CMS Transmittal 907 (Change Request 7409, June 14, 2011), boiled down to figures.
// Every fact on this page comes from the PDF; captions cite PDF viewer page numbers.

const MONO = 'IBM Plex Mono, monospace'

type Tone = 'loop' | 'lseq' | 'seg' | 'sseq' | 'data' | 'fill'

const TONE: Record<Tone, { cell: string; label: string }> = {
  loop: { cell: 'bg-series-1/15 border-series-1/60 text-ink', label: 'text-series-1' },
  lseq: { cell: 'bg-series-2/15 border-series-2/60 text-ink', label: 'text-series-2' },
  seg: { cell: 'bg-series-3/15 border-series-3/60 text-ink', label: 'text-series-3' },
  sseq: { cell: 'bg-ok/15 border-ok/60 text-ink', label: 'text-ok' },
  data: { cell: 'bg-surface border-line text-ink', label: 'text-ink-2' },
  fill: { cell: 'bg-bg border-line border-dashed text-ink-2', label: 'text-ink-2' },
}

interface Part {
  text: string
  label: string
  tone: Tone
}

// A run of fixed-width bytes, one cell per byte. Spaces show as a middle dot
// so space fill is visible. Each part is labeled underneath with its byte range.
function ByteRun({ parts, start = 1 }: { parts: Part[]; start?: number }) {
  let pos = start
  return (
    <div className="flex flex-wrap items-start gap-y-3">
      {parts.map((part, i) => {
        const from = pos
        pos += part.text.length
        const to = pos - 1
        return (
          <div key={i} className="flex flex-col">
            <div className="flex">
              {part.text.split('').map((ch, j) => (
                <span
                  key={j}
                  className={`-ml-px flex h-7 w-[1.15rem] items-center justify-center border font-mono text-[12px] first:ml-0 ${TONE[part.tone].cell}`}
                >
                  {ch === ' ' ? <span className="text-ink-2/60">·</span> : ch}
                </span>
              ))}
            </div>
            <div className={`mt-1 pr-2 font-mono text-[10.5px] leading-4 ${TONE[part.tone].label}`}>
              <div>{from === to ? from : `${from}–${to}`}</div>
              <div className="max-w-[12rem]">{part.label}</div>
            </div>
          </div>
        )
      })}
    </div>
  )
}

// Figure 1: the 283 pages as one proportional bar.
const DOC_MAP: { pages: string; count: number; title: string; tone: string }[] = [
  { pages: '1–6', count: 6, title: 'Transmittal notice: summary, background, requirements', tone: 'bg-accent' },
  { pages: '7–125', count: 119, title: 'Attachment 1 · 837I (institutional claim) layout', tone: 'bg-series-1/50' },
  { pages: '126–132', count: 7, title: '837I change log, column legend, error-record example', tone: 'bg-series-1/25' },
  { pages: '133–239', count: 107, title: 'Attachment 2 · 837P (professional claim) layout', tone: 'bg-series-3/50' },
  { pages: '240–245', count: 6, title: '837P legend, change log, error-record and response examples', tone: 'bg-series-3/25' },
  { pages: '246–280', count: 35, title: 'Attachment 3 · 835 (payment and remittance advice) layout', tone: 'bg-series-2/60' },
  { pages: '281–283', count: 3, title: '835 change log and column legend', tone: 'bg-series-2/30' },
]

function DocumentMap() {
  return (
    <div className="min-w-[40rem]">
      <div className="flex h-10 w-full overflow-hidden border border-line">
        {DOC_MAP.map((s) => (
          <div
            key={s.pages}
            className={`${s.tone} flex min-w-[6px] items-center justify-center border-r border-bg font-mono text-[11px] text-ink last:border-r-0`}
            style={{ flexGrow: s.count, flexBasis: 0 }}
            title={`PDF pp. ${s.pages}: ${s.title}`}
          >
            {s.count > 30 ? `pp. ${s.pages}` : ''}
          </div>
        ))}
      </div>
      <ul className="mt-4 grid gap-x-6 gap-y-1.5 text-[13px] md:grid-cols-2">
        {DOC_MAP.map((s) => (
          <li key={s.pages} className="flex items-baseline gap-2">
            <span className={`${s.tone} inline-block h-3 w-3 shrink-0 border border-line`} />
            <span className="w-16 shrink-0 font-mono text-[11px] text-ink-2">pp. {s.pages}</span>
            <span className="text-ink">{s.title}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}

// Figure 2: why the transmittal exists.
const HISTORY = [
  { when: 'Nov 10, 2010', what: 'CR 7202 / Transmittal 2090', note: 'Publishes 5010 flat files for 837I, 837P, and 835.' },
  { when: 'CR 7065', what: 'Fractional units', note: 'Contractors must accept and report fractional units for ambulance services.' },
  { when: 'The gap', what: 'CAS cannot hold fractions', note: 'The CR 7202 flat files do not support fractional units in the claim adjustment (CAS) fields.', gap: true },
  { when: 'Apr 15, 2011', what: 'Transmittal 873', note: 'Earlier issue of this change request; rescinded and replaced.' },
  { when: 'Jun 14, 2011', what: 'Transmittal 907 / CR 7409', note: 'Replaces 873; corrects typos and omissions. Effective July 1, implemented July 5, 2011.', current: true },
]

function WhyItExists() {
  return (
    <div className="min-w-[44rem] space-y-6">
      <ol className="grid grid-cols-5 gap-px border border-line bg-line">
        {HISTORY.map((h, i) => (
          <li key={i} className={`relative p-3 ${h.current ? 'bg-accent/10' : 'bg-surface'}`}>
            <div className={`font-mono text-[11px] ${h.gap ? 'text-flag' : 'text-ink-2'}`}>{h.when}</div>
            <div className={`mt-1 text-[13px] font-semibold ${h.gap ? 'text-flag' : 'text-ink'}`}>{h.what}</div>
            <div className="mt-1 text-[12px] leading-5 text-ink-2">{h.note}</div>
            {i < HISTORY.length - 1 && (
              <span className="absolute -right-2 top-3 z-10 font-mono text-[13px] text-ink-2">→</span>
            )}
          </li>
        ))}
      </ol>
      <div className="grid grid-cols-[auto_1fr] items-center gap-x-6 gap-y-3 text-[13px]">
        <span className="font-mono text-[11px] text-ink-2">the fix</span>
        <div className="flex flex-wrap items-center gap-3">
          <span className="text-ink">CAS adjustment quantity (CAS04, 07, 10, 13, 16, 19) becomes</span>
          <ByteRun
            parts={[
              { text: '9999999', label: '7 whole digits', tone: 'data' },
              { text: '999', label: '3 implied decimals', tone: 'lseq' },
            ]}
          />
          <code className="font-mono text-[12px] text-ink">PIC S9(7)V999</code>
        </div>
        <span className="font-mono text-[11px] text-ink-2">not changed</span>
        <span className="text-ink-2">
          The older 4010A1 flat files keep whole units. Remittances to providers not yet on 5010 show whole units even when
          payment used fractional mileage.
        </span>
        <span className="font-mono text-[11px] text-ink-2">837I rule</span>
        <span className="text-ink-2">
          Requirement 7409.7: the institutional shared system stores service units (SV205) as a decimal{' '}
          <code className="font-mono text-[12px]">9(6)V9</code> when the revenue code is 0540, and as a whole number{' '}
          <code className="font-mono text-[12px]">9(7)</code> otherwise.
        </span>
      </div>
    </div>
  )
}

// Figure 3: the 18-byte record key.
function RecordKey() {
  return (
    <div className="min-w-[40rem] space-y-6">
      <ByteRun
        parts={[
          { text: '2300  ', label: 'Loop ID · left-justified, space-filled', tone: 'loop' },
          { text: '0001', label: 'Loop sequence · right-justified, zero-filled', tone: 'lseq' },
          { text: 'CLM ', label: 'Segment ID · left-justified, space-filled', tone: 'seg' },
          { text: '0120', label: 'Segment sequence · right-justified, zero-filled', tone: 'sseq' },
          { text: 'PARTB               ', label: 'Data from byte 19: CLM01, 20 bytes (837P)', tone: 'data' },
          { text: '…', label: 'more fields at their Start positions', tone: 'fill' },
        ]}
      />
      <table className="w-full max-w-[44rem] text-left text-[13px]">
        <thead>
          <tr className="border-b border-line text-[11px] text-ink-2">
            <th className="py-2 pr-4 font-mono font-normal">Bytes</th>
            <th className="py-2 pr-4 font-mono font-normal">Part</th>
            <th className="py-2 pr-4 font-mono font-normal">Width</th>
            <th className="py-2 font-mono font-normal">Legend's own example</th>
          </tr>
        </thead>
        <tbody className="text-ink">
          {[
            ['1–6', 'Loop ID', '6', '"2010AA"'],
            ['7–10', 'Loop sequence', '4', '"0001"'],
            ['11–14', 'Segment ID', '4', '"REF "'],
            ['15–18', 'Segment sequence', '4', '"0001"'],
          ].map(([b, p, w, e], i) => (
            <tr key={b} className="border-b border-line">
              <td className="py-2 pr-4 font-mono text-[12px]">{b}</td>
              <td className={`py-2 pr-4 ${[TONE.loop, TONE.lseq, TONE.seg, TONE.sseq][i].label}`}>{p}</td>
              <td className="py-2 pr-4 font-mono text-[12px]">{w}</td>
              <td className="py-2 font-mono text-[12px]">{e}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

// Figure 4: the 835 as a tree of loops, with where each piece sits in the PDF.
interface Node {
  key: string
  name: string
  records: string
  repeat?: string
  pages: string
  children?: Node[]
  highlight?: boolean
}

const TREE_835: Node = {
  key: 'ISA',
  name: 'Interchange',
  records: 'ISA … IEA',
  pages: '246, 280',
  children: [
    {
      key: 'GS',
      name: 'Functional group',
      records: 'GS … GE',
      pages: '246–247, 279–280',
      children: [
        {
          key: 'ST',
          name: 'Transaction (one 835)',
          records: 'ST · BPR · TRN · REF · DTM … SE',
          pages: '247–249, 279',
          children: [
            { key: '1000A', name: 'Payer', records: 'N1 · N3 · N4 · REF ×4 · PER', pages: '249–251' },
            { key: '1000B', name: 'Payee', records: 'N1 · N3 · N4 · REF', pages: '251–252' },
            {
              key: '2000',
              name: 'Header number',
              records: 'LX · TS3 · TS2',
              repeat: 'loop >1',
              pages: '252–258',
              children: [
                {
                  key: '2100',
                  name: 'Claim',
                  records: 'CLP · CAS ×99 · NM1 · MIA or MOA · REF · DTM · PER · AMT ×13 · QTY ×14',
                  repeat: 'loop >1',
                  pages: '258–271',
                  highlight: true,
                  children: [
                    {
                      key: '2110',
                      name: 'Service line',
                      records: 'SVC · DTM ×2 · CAS ×99 · REF · AMT · LQ ×99',
                      repeat: 'loop 999',
                      pages: '271–276',
                    },
                  ],
                },
              ],
            },
            { key: 'PLB', name: 'Provider-level adjustment', records: 'PLB', repeat: '>1', pages: '277–279' },
          ],
        },
      ],
    },
  ],
}

function TreeNode({ node }: { node: Node }) {
  return (
    <div className={`border-l-2 ${node.highlight ? 'border-accent' : 'border-line'} pl-4`}>
      <div className={`flex flex-wrap items-baseline gap-x-3 py-1.5 ${node.highlight ? 'bg-accent/10 -ml-4 pl-4' : ''}`}>
        <span className="w-14 shrink-0 font-mono text-[13px] font-semibold text-series-1">{node.key}</span>
        <span className="text-[13px] font-semibold text-ink">{node.name}</span>
        {node.repeat && <span className="border border-line px-1.5 font-mono text-[10.5px] text-ink-2">{node.repeat}</span>}
        <span className="font-mono text-[11.5px] text-ink-2">{node.records}</span>
        <span className="ml-auto pr-2 font-mono text-[10.5px] text-ink-2">p. {node.pages}</span>
      </div>
      {node.children && (
        <div className="space-y-1 pb-1">
          {node.children.map((c) => (
            <TreeNode key={c.key} node={c} />
          ))}
        </div>
      )}
    </div>
  )
}

// Figure 5: one claim is a run of records, each keyed by the same 18-byte scheme.
const CLAIM_RECORDS: { loop: string; seg: string; what: string; repeat?: string; service?: boolean }[] = [
  { loop: '2100', seg: 'CLP', what: 'Claim: patient control #, status, charge, payment, patient responsibility, payer claim control #' },
  { loop: '2100', seg: 'CAS', what: 'Claim adjustments: group, reason, amount, quantity', repeat: 'up to 99 records' },
  { loop: '2100', seg: 'NM1', what: 'Patient, corrected patient, service provider, crossover carrier, corrected priority payer names' },
  { loop: '2100', seg: 'MIA', what: 'Inpatient or outpatient (MOA) adjudication information' },
  { loop: '2100', seg: 'DTM', what: 'Statement dates, claim received date' },
  { loop: '2100', seg: 'AMT', what: 'Claim supplemental amounts', repeat: 'up to 13' },
  { loop: '2110', seg: 'SVC', what: 'Service line: procedure code, charge, payment, units', service: true },
  { loop: '2110', seg: 'DTM', what: 'Service dates', repeat: 'up to 2', service: true },
  { loop: '2110', seg: 'CAS', what: 'Service adjustments', repeat: 'up to 99', service: true },
  { loop: '2110', seg: 'LQ', what: 'Remark codes', repeat: 'up to 99', service: true },
  { loop: '2110', seg: 'SVC', what: 'Next service line … (loop repeats up to 999)', service: true },
  { loop: '2100', seg: 'CLP', what: 'Next claim starts a new run' },
]

function ClaimAsRecords() {
  return (
    <div className="min-w-[44rem]">
      <div className="mb-2 grid grid-cols-[13rem_1fr] font-mono text-[10.5px] text-ink-2">
        <span>record key, bytes 1–18</span>
        <span>data from byte 19</span>
      </div>
      <div className="space-y-1">
        {CLAIM_RECORDS.map((r, i) => {
          const last = i === CLAIM_RECORDS.length - 1
          return (
            <div
              key={i}
              className={`grid grid-cols-[13rem_1fr] items-center ${r.service ? 'pl-5' : ''} ${last ? 'opacity-50' : ''}`}
            >
              <div className="flex font-mono text-[12px]">
                <span className={`border px-1 ${TONE.loop.cell}`}>{r.loop.padEnd(6, '·')}</span>
                <span className={`-ml-px border px-1 ${TONE.lseq.cell}`}>nnnn</span>
                <span className={`-ml-px border px-1 ${TONE.seg.cell}`}>{r.seg.padEnd(4, '·')}</span>
                <span className={`-ml-px border px-1 ${TONE.sseq.cell}`}>nnnn</span>
              </div>
              <div className={`flex items-baseline gap-3 border-l-2 pl-3 text-[12.5px] ${r.service ? 'border-series-3' : 'border-accent'}`}>
                <span className="text-ink">{r.what}</span>
                {r.repeat && <span className="font-mono text-[10.5px] text-ink-2">{r.repeat}</span>}
              </div>
            </div>
          )
        })}
      </div>
      <p className="mt-3 font-mono text-[10.5px] text-ink-2">
        · = space fill. nnnn = four-digit sequence, zero-filled. Service-line records indented.
      </p>
    </div>
  )
}

// Figure 6: the 835 CLP record drawn to scale, plus the anatomy of an amount field.
const CLP_FIELDS: { id: string; name: string; start: number; len: number; amount?: boolean }[] = [
  { id: 'key', name: 'record key', start: 1, len: 18 },
  { id: 'CLP01', name: 'patient control #', start: 19, len: 38 },
  { id: 'CLP02', name: 'claim status', start: 57, len: 2 },
  { id: 'CLP03', name: 'total charge', start: 59, len: 18, amount: true },
  { id: 'CLP04', name: 'payment', start: 77, len: 18, amount: true },
  { id: 'CLP05', name: 'patient resp.', start: 95, len: 18, amount: true },
  { id: 'CLP06', name: 'filing indicator', start: 113, len: 2 },
  { id: 'CLP07', name: 'payer claim control #', start: 115, len: 50 },
  { id: 'CLP08', name: 'facility code', start: 165, len: 2 },
  { id: 'CLP09', name: 'frequency code', start: 167, len: 1 },
  { id: 'CLP11', name: 'DRG code', start: 168, len: 4 },
  { id: 'CLP12', name: 'DRG weight', start: 172, len: 15 },
  { id: 'CLP13', name: 'discharge fraction', start: 187, len: 10 },
]

// Label rows and x positions for fields too narrow to label inside the bar,
// chosen so leader lines do not cross other labels.
const NARROW_LABELS: Record<string, { row: number; lx: number }> = {
  CLP02: { row: 0, lx: 290 },
  CLP06: { row: 0, lx: 570 },
  CLP11: { row: 0, lx: 866 },
  CLP08: { row: 1, lx: 690 },
  CLP09: { row: 2, lx: 690 },
}

function ClpRecord() {
  const px = 5
  const x0 = 10
  const barY = 58
  const barH = 40
  const narrow = CLP_FIELDS.filter((f) => f.len < 10)
  return (
    <svg viewBox="0 0 1000 250" width="100%" role="img" style={{ minWidth: 760 }}>
      <title>The 835 CLP claim record, 196 bytes, drawn to scale by byte position.</title>
      {/* byte ruler */}
      {[1, 19, 57, 77, 95, 115, 165, 172, 187, 196].map((b) => (
        <g key={b}>
          <line x1={x0 + (b - 1) * px} y1={40} x2={x0 + (b - 1) * px} y2={barY} stroke="var(--color-line)" />
          <text x={x0 + (b - 1) * px + 2} y={36} fontFamily={MONO} fontSize="10" fill="var(--color-ink-2)">
            {b}
          </text>
        </g>
      ))}
      <text x={x0} y={16} fontFamily={MONO} fontSize="11" fill="var(--color-ink-2)">
        byte position →
      </text>
      {CLP_FIELDS.map((f) => {
        const x = x0 + (f.start - 1) * px
        const w = f.len * px
        const fill = f.id === 'key' ? 'var(--color-series-1)' : f.amount ? 'var(--color-series-2)' : 'var(--color-surface)'
        return (
          <g key={f.id}>
            <rect
              x={x}
              y={barY}
              width={w}
              height={barH}
              fill={fill}
              fillOpacity={f.id === 'key' || f.amount ? 0.22 : 1}
              stroke="var(--color-ink-2)"
              strokeWidth="0.8"
            />
            {f.len >= 10 && (
              <>
                <text x={x + 4} y={barY + 16} fontFamily={MONO} fontSize="11" fill="var(--color-ink)">
                  {f.id}
                </text>
                <text x={x + 4} y={barY + 31} fontFamily={MONO} fontSize="9.5" fill="var(--color-ink-2)">
                  {f.len * px > 70 ? f.name : `${f.len}B`}
                </text>
              </>
            )}
          </g>
        )
      })}
      {/* narrow fields: labels below with leader lines */}
      {narrow.map((f) => {
        const cx = x0 + (f.start - 1) * px + (f.len * px) / 2
        const place = NARROW_LABELS[f.id]
        const ly = barY + barH + 20 + place.row * 28
        const lx = place.lx
        return (
          <g key={f.id}>
            <path
              d={`M ${cx} ${barY + barH} V ${ly - 13} H ${lx}`}
              fill="none"
              stroke="var(--color-ink-2)"
              strokeWidth="0.8"
            />
            <text x={lx} y={ly} fontFamily={MONO} fontSize="10.5" fill="var(--color-ink)">
              {f.id}
            </text>
            <text x={lx} y={ly + 12} fontFamily={MONO} fontSize="9.5" fill="var(--color-ink-2)">
              {f.name} · {f.len}B
            </text>
          </g>
        )
      })}
      {/* legend */}
      <rect x={x0} y={214} width={12} height={12} fill="var(--color-series-2)" fillOpacity={0.22} stroke="var(--color-ink-2)" strokeWidth="0.8" />
      <text x={x0 + 18} y={224} fontFamily={MONO} fontSize="10.5" fill="var(--color-ink-2)">
        amount field, 18 bytes (anatomy below)
      </text>
      <rect x={x0 + 300} y={214} width={12} height={12} fill="var(--color-series-1)" fillOpacity={0.22} stroke="var(--color-ink-2)" strokeWidth="0.8" />
      <text x={x0 + 318} y={224} fontFamily={MONO} fontSize="10.5" fill="var(--color-ink-2)">
        18-byte record key
      </text>
      <text x={x0 + 520} y={224} fontFamily={MONO} fontSize="10.5" fill="var(--color-ink-2)">
        CLP10 and CLP14 are not used and take no bytes
      </text>
    </svg>
  )
}

function AmountAnatomy() {
  return (
    <div className="mt-6 space-y-5 border-t border-line pt-5">
      <div className="font-mono text-[11px] text-ink-2">
        every 835 amount field · PIC X(18) redefined as S9(8)V99 + X(8) filler
      </div>
      <ByteRun
        parts={[
          { text: '99999999', label: '8 whole digits, right-justified, zero-filled', tone: 'data' },
          { text: '99', label: '2 cents digits; decimal point implied, not stored', tone: 'lseq' },
          { text: '        ', label: 'filler: spaces', tone: 'fill' },
        ]}
      />
      <div className="grid gap-x-8 gap-y-3 text-[12.5px] md:grid-cols-2">
        <div>
          <div className="mb-1.5 font-mono text-[11px] text-ink-2">no amount to report</div>
          <ByteRun parts={[{ text: '                  ', label: 'all 18 bytes spaces', tone: 'fill' }]} />
        </div>
        <div>
          <div className="mb-1.5 font-mono text-[11px] text-ink-2">allowed FISS exception</div>
          <ByteRun
            parts={[
              { text: '0000000000', label: 'zeros', tone: 'data' },
              { text: '        ', label: 'spaces', tone: 'fill' },
            ]}
          />
        </div>
      </div>
    </div>
  )
}

// Figure 7: records the flat file carries that X12 does not.
function ErrorRecords() {
  return (
    <div className="min-w-[48rem] space-y-6">
      <div className="grid grid-cols-[1fr_auto_1fr] items-center gap-4 text-[12.5px]">
        <div className="border border-line bg-surface p-3">
          <div className="font-mono text-[10.5px] text-ink-2">provider's 837P, as received</div>
          <pre className="mt-1 font-mono text-[12px] text-ink">CLM{'>'}PARTB{'>'}<span className="text-flag">92.511</span>{'>'}{'>'}{'>'}11+B+1{'>'}Y{'>'}A{'>'}Y{'>'}Y{'>'}P</pre>
          <div className="mt-1 text-ink-2">Error: CLM02 has too many decimal positions (code 697).</div>
        </div>
        <span className="font-mono text-ink-2">→</span>
        <div className="border border-line bg-surface p-3">
          <div className="font-mono text-[10.5px] text-ink-2">translator writes flat records for the edit module</div>
          <div className="mt-1 font-mono text-[12px] text-ink">
            2300 CLM record, then a new <span className="text-series-3">STC</span> record right after it
          </div>
          <div className="mt-1 text-ink-2">The STC record carries the error; the X12 claim had no STC segment.</div>
        </div>
      </div>
      <div>
        <div className="mb-2 font-mono text-[11px] text-ink-2">the inserted STC error record, laid out by the STC definition</div>
        <ByteRun
          parts={[
            { text: '2300  ', label: 'loop ID', tone: 'loop' },
            { text: '0001', label: 'loop seq', tone: 'lseq' },
            { text: 'STC ', label: 'seg ID', tone: 'seg' },
            { text: '0120', label: 'seg seq: same as the CLM it follows', tone: 'sseq' },
            { text: 'A7   ', label: 'STC01-1 category', tone: 'data' },
            { text: '697  ', label: 'STC01-2 status code', tone: 'data' },
            { text: '   ', label: 'STC01-3 entity', tone: 'fill' },
            { text: '20090819', label: 'STC02 date 9(8)', tone: 'data' },
            { text: 'U ', label: 'STC03 action', tone: 'data' },
            { text: '000000925A', label: 'STC04 charge S9(8)V99', tone: 'lseq' },
            { text: 'A7   ', label: 'STC10-1', tone: 'data' },
            { text: '178  ', label: 'STC10-2', tone: 'data' },
          ]}
        />
      </div>
      <ul className="list-disc space-y-1 pl-5 text-[12.5px] text-ink-2">
        <li>
          The charge 92.511 arrives as <code className="font-mono text-ink">000000925A</code>: ten bytes, no decimal point,
          last digit replaced by a letter. The document does not define that letter. In COBOL signed fields a trailing{' '}
          <code className="font-mono text-ink">A</code> commonly means +1, which reads as 92.51. <em>Verify</em> the sign
          encoding with the producer before decoding.
        </li>
        <li>
          The segment sequence 0120 matches the 999 acknowledgment, which reports the CLM error at segment position 120
          (p. 245).
        </li>
        <li>
          The 837I layout has similar internal-only records: a <code className="font-mono text-ink">DTP</code> with value
          +RC and a <code className="font-mono text-ink">REF</code> with value +PR whose data is P (paper), K (keyshop), or O
          (optical character recognition) (p. 8).
        </li>
      </ul>
    </div>
  )
}

// Same business field, different flat-file storage.
const SAME_FIELD = [
  ['837I CLM01 patient control #', 'AN 1–38', 'X(20)', '19', '20', '23'],
  ['837P CLM01 patient account #', 'AN 1–38', 'X(20)', '19', '20', '154'],
  ['835 CLP01 patient control #', 'AN 1–38', 'X(38)', '19', '38', '258'],
  ['837I CLM02 total charge', 'R 1–18', 'S9(8)V99', '39', '10', '23'],
  ['837P CLM02 total charge', 'R 1–18', 'S9(5)V99 + X(11)', '39', '18', '154'],
  ['835 CLP03 total charge', 'R 1–18', 'S9(8)V99 + X(8)', '59', '18', '258'],
]

const WHERE = [
  ['Summary of changes; effective and implementation dates', '1'],
  ['Background: CR 7202, CR 7065, fractional units; 4010A1 not updated', '3'],
  ['Business requirements 7409.1–7409.8', '4'],
  ['837I layout begins (ISA record, byte positions)', '7'],
  ['837I internal DTP +RC and REF +PR records', '8'],
  ['837I CLM claim record', '23'],
  ['837I change log (including "all amounts … 18 to 10")', '126–128'],
  ['837I column legend and 18-byte key definition', '129–130'],
  ['837I STC error-record example', '131–132'],
  ['837P layout begins', '133'],
  ['837P CLM claim record', '154'],
  ['837P column legend and 18-byte key definition', '240–241'],
  ['837P change log', '242–243'],
  ['837P STC error record, 999 and 277CA examples', '244–245'],
  ['835 layout begins; amount-field rule stated in each page header', '246'],
  ['835 BPR payment and TRN trace records', '247–248'],
  ['835 1000A payer and 1000B payee loops', '249–252'],
  ['835 2000 header-number loop (LX, TS3, TS2)', '252–258'],
  ['835 2100 CLP claim record; CAS claim adjustments', '258–259'],
  ['835 2110 SVC service record and its DTM, CAS, REF, AMT, LQ', '271–276'],
  ['835 PLB provider-level adjustment', '277–279'],
  ['835 SE, GE, IEA trailers', '279–280'],
  ['835 change log', '281–282'],
  ['835 column legend and 18-byte key definition', '283'],
]

export function FormatTransmittal907() {
  return (
    <div className="mx-auto max-w-[90rem] px-4 py-12 md:px-10">
      <PageTitle
        lede={
          <>
            CMS Transmittal 907, Change Request (CR) 7409, June 14, 2011. Pub. 100-20 One-Time Notification, 283 pages.
          </>
        }
      >
        Transmittal 907: 837 and 835 flat files
      </PageTitle>

      <DocSection id="gist" title="The gist">
        <p>
          This is a Centers for Medicare &amp; Medicaid Services (CMS) instruction to Medicare Administrative Contractors
          (MACs) and the shared claims-processing systems. Its three attachments are the fixed-position “flat file”
          layouts that Medicare systems use internally for the 837I institutional claim, the 837P professional claim, and
          the 835 payment and remittance advice. Every X12 segment becomes one or more fixed-width records that start with
          an 18-byte key, so one claim spans many records.
        </p>
        <p>
          For a validator that compares MAC-processed 835s against vendor-processed 835s, Attachment 3 shows one form a
          supplied MAC file could take: the same 835 fields at fixed byte positions, with implied decimal points and no
          delimiters.
        </p>
        <Figure number={1} caption="What the 283 pages contain, drawn to scale. Over 90% of the document is the three attachment layouts. (PDF pp. 1–283)">
          <DocumentMap />
        </Figure>
      </DocSection>

      <DocSection id="why" title="Why this transmittal exists">
        <p>The narrative is short: the earlier flat files could not hold fractional units in the claim adjustment segment (CAS).</p>
        <Figure number={2} caption="From CR 7202 to Transmittal 907, and the field-level fix. Dates and requirement 7409.7 from pp. 1, 3–4; CAS change from change logs pp. 242, 282.">
          <WhyItExists />
        </Figure>
        <ul className="list-disc space-y-1 pl-5 text-[14px]">
          <li>7409.1–7409.3: contractors implement the 837I, 837P, and 835 flat files (Attachments 1–3).</li>
          <li>7409.4–7409.5: the remittance print tools (Medicare Remit Easy Print and PC Print) follow Attachment 3.</li>
          <li>7409.7–7409.8: the Fiscal Intermediary Shared System (FISS) moves ambulance units as decimals in and out of the claim record.</li>
        </ul>
      </DocSection>

      <DocSection id="key" title="The 18-byte record key">
        <p>
          Every record begins with the same key. The data fields that follow sit at the Start and Length positions in each
          layout row. These are byte positions in the record, not X12 delimiter offsets.
        </p>
        <Figure number={3} caption="The key, using the document's own sample record (837P CLM, p. 244). Justification and fill rules from the column legend, pp. 129–130, 241, 283.">
          <RecordKey />
        </Figure>
      </DocSection>

      <DocSection id="tree" title="The 835 as loops and records">
        <Figure number={4} caption="The 835 flat file's loop tree. The loop ID column gives the key's first 6 bytes; ×N is the Record Repeat limit. Claim loop highlighted. (pp. 246–280)">
          <TreeNode node={TREE_835} />
        </Figure>
        <Figure number={5} caption="One claim in the 835 flat file is a run of records: a 2100 CLP record, its claim-level records, then each 2110 service line's records. Selected records shown; repeat limits from pp. 258–276.">
          <ClaimAsRecords />
        </Figure>
      </DocSection>

      <DocSection id="clp" title="The claim record, byte by byte">
        <p>
          Amount fields use one rule across the 835 layout, printed at the top of every page: right-justified and
          zero-filled in the first 10 bytes, space-filled in bytes 11–18.
        </p>
        <Figure number={6} caption="835 CLP record to scale (pp. 258–259); DRG is diagnosis-related group. Below it, the amount-field rule from the page header (p. 246 onward).">
          <ClpRecord />
          <AmountAnatomy />
        </Figure>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[40rem] text-left text-[13px]">
            <caption className="mb-2 text-left font-mono text-[11px] text-ink-2">
              Same business field, different storage. X12 allows up to 38 characters; the 837 flat files keep 20.
            </caption>
            <thead>
              <tr className="border-b border-line font-mono text-[11px] text-ink-2">
                <th className="py-2 pr-4 font-normal">Field</th>
                <th className="py-2 pr-4 font-normal">X12 size</th>
                <th className="py-2 pr-4 font-normal">COBOL picture</th>
                <th className="py-2 pr-4 font-normal">Start</th>
                <th className="py-2 pr-4 font-normal">Length</th>
                <th className="py-2 font-normal">PDF p.</th>
              </tr>
            </thead>
            <tbody className="text-ink">
              {SAME_FIELD.map((r) => (
                <tr key={r[0]} className="border-b border-line">
                  <td className="py-2 pr-4">{r[0]}</td>
                  <td className="py-2 pr-4 font-mono text-[12px]">{r[1]}</td>
                  <td className="py-2 pr-4 font-mono text-[12px]">{r[2]}</td>
                  <td className="py-2 pr-4 font-mono text-[12px]">{r[3]}</td>
                  <td className="py-2 pr-4 font-mono text-[12px]">{r[4]}</td>
                  <td className="py-2 font-mono text-[12px]">{r[5]}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DocSection>

      <DocSection id="extra" title="Records that X12 does not have">
        <Figure number={7} caption="The commercial off-the-shelf (COTS) translator inserts STC error records after a segment with errors and passes them to the Common Edit Module (CEM). Sample and STC layout from p. 244; 999 example p. 245.">
          <ErrorRecords />
        </Figure>
      </DocSection>

      <DocSection id="validator" title="What this means for the validator">
        <ul className="list-disc space-y-2 pl-5">
          <li>
            A claim is a run of records, not one line. In the 835 layout it starts at a <code>2100</code> <code>CLP</code>{' '}
            record and continues through its 2100 and 2110 records until the next <code>CLP</code>, a <code>PLB</code>, or{' '}
            <code>SE</code>.
          </li>
          <li>
            The 18-byte key places a record inside one file. It does not identify a claim across systems; cross-system
            matching uses the claim's own content: the beneficiary identifier, the date of service, the total
            charge, the payee's tax ID and NPI, and the MAC workload.
          </li>
          <li>
            Amounts have an implied decimal point and fixed width. Decode them to exact decimals and keep the raw bytes.
            An empty amount can be 18 spaces or, for FISS, zeros then spaces; exact comparison keeps those spellings
            distinct, and any equivalence is a later, labeled policy.
          </li>
          <li>
            Flat widths can be narrower than X12 allows (<code>CLM01</code>: 20 of 38). A difference can come from the
            layout rather than from adjudication.
          </li>
          <li>
            This is a 2011 edition, and the change logs show widths and positions moving between revisions. Identify a
            supplied file's layout revision before applying any offsets from this page.
          </li>
        </ul>
      </DocSection>

      <DocSection id="where" title="Where to look in the PDF">
        <p className="text-[13px]">PDF viewer page numbers, counted from 1. Each attachment restarts its printed page numbers.</p>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[32rem] text-left text-[13px]">
            <thead>
              <tr className="border-b border-line font-mono text-[11px] text-ink-2">
                <th className="py-2 pr-4 font-normal">Topic</th>
                <th className="py-2 font-normal">PDF pages</th>
              </tr>
            </thead>
            <tbody className="text-ink">
              {WHERE.map(([topic, pages]) => (
                <tr key={topic} className="border-b border-line">
                  <td className="py-2 pr-4">{topic}</td>
                  <td className="py-2 font-mono text-[12px]">{pages}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </DocSection>
    </div>
  )
}
