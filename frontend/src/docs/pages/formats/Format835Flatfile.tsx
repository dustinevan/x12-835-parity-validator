import type { ReactNode } from 'react'
import { DocSection } from '../../../components/DocSection'
import { Figure } from '../../../components/Figure'
import { PageTitle } from '../../../components/PageTitle'

// A one-page digest of the CMS 835 flat-file layout (5010A1, updated
// January 10, 2023).
// Every position, width, and repeat count below is copied from that PDF;
// page numbers are PDF viewer pages, counted from 1.

type Tone = 'key' | 'money' | 'id' | 'code' | 'date' | 'text' | 'filler'

const TONE: Record<Tone, string> = {
  key: 'bg-accent/15',
  money: 'bg-ok/15',
  id: 'bg-series-2/20',
  code: 'bg-series-3/15',
  date: 'bg-accent/5',
  text: 'bg-surface',
  filler: 'bg-bg',
}

const LEGEND: { tone: Tone; label: string }[] = [
  { tone: 'key', label: 'Record key' },
  { tone: 'money', label: 'Amount or quantity' },
  { tone: 'id', label: 'Identifier' },
  { tone: 'code', label: 'Code' },
  { tone: 'text', label: 'Other' },
]

interface Span {
  id: string
  start: number
  len: number
  tone: Tone
  name?: string
}

// A byte ruler: one cell per field, width proportional to its length.
function ByteRuler({ spans, minCell = '2.75rem' }: { spans: Span[]; minCell?: string }) {
  return (
    <div className="flex min-w-[46rem] border border-line">
      {spans.map((s) => (
        <div
          key={`${s.id}-${s.start}`}
          title={s.name ? `${s.id} ${s.name}` : s.id}
          className={`${TONE[s.tone]} border-r border-line px-1 py-2 last:border-r-0`}
          style={{ flexGrow: s.len, flexBasis: 0, minWidth: minCell }}
        >
          <div className="truncate font-mono text-[11px] font-semibold text-ink">{s.id}</div>
          <div className="truncate font-mono text-[10px] text-ink-2">
            {s.start}
            {s.len > 1 ? `–${s.start + s.len - 1}` : ''}
          </div>
        </div>
      ))}
    </div>
  )
}

function Legend() {
  return (
    <div className="mt-3 flex flex-wrap gap-4 text-[12px] text-ink-2">
      {LEGEND.map((l) => (
        <span key={l.tone} className="flex items-center gap-1.5">
          <span className={`inline-block h-3 w-3 border border-line ${TONE[l.tone]}`} />
          {l.label}
        </span>
      ))}
    </div>
  )
}

// Figure 2: the record key.
const KEY: Span[] = [
  { id: 'Loop ID', start: 1, len: 6, tone: 'key' },
  { id: 'Loop Seq', start: 7, len: 4, tone: 'key' },
  { id: 'Seg ID', start: 11, len: 4, tone: 'key' },
  { id: 'Seg Seq', start: 15, len: 4, tone: 'key' },
]

// Figure 4: the CLP (claim) record, PDF pp. 5–6.
const CLP: Span[] = [
  { id: 'Key', start: 1, len: 18, tone: 'key', name: '2100 / CLP record key' },
  { id: 'CLP01', start: 19, len: 38, tone: 'id', name: 'Patient control number' },
  { id: 'CLP02', start: 57, len: 2, tone: 'code', name: 'Claim status code' },
  { id: 'CLP03', start: 59, len: 18, tone: 'money', name: 'Total claim charge amount' },
  { id: 'CLP04', start: 77, len: 18, tone: 'money', name: 'Claim payment amount' },
  { id: 'CLP05', start: 95, len: 18, tone: 'money', name: 'Patient responsibility amount' },
  { id: 'CLP06', start: 113, len: 2, tone: 'code', name: 'Claim filing indicator (MA/MB)' },
  { id: 'CLP07', start: 115, len: 50, tone: 'id', name: 'Payer claim control number' },
  { id: 'CLP08', start: 165, len: 2, tone: 'code', name: 'Facility code value' },
  { id: 'CLP09', start: 167, len: 1, tone: 'code', name: 'Claim frequency code' },
  { id: 'CLP11', start: 168, len: 4, tone: 'code', name: 'DRG code (Part A only)' },
  { id: 'CLP12', start: 172, len: 15, tone: 'money', name: 'DRG weight (Part A only)' },
  { id: 'CLP13', start: 187, len: 10, tone: 'money', name: 'Discharge fraction' },
]

// Figure 6: the CAS (adjustment) record, PDF pp. 6 and 12.
const CAS: Span[] = [
  { id: 'Key', start: 1, len: 18, tone: 'key' },
  { id: 'CAS01', start: 19, len: 2, tone: 'code', name: 'Group code' },
  ...[21, 59, 97, 135, 173, 211].flatMap((reasonStart, i): Span[] => [
    { id: `CAS${String(2 + i * 3).padStart(2, '0')}`, start: reasonStart, len: 5, tone: 'code', name: 'Reason code' },
    { id: `CAS${String(3 + i * 3).padStart(2, '0')}`, start: reasonStart + 5, len: 18, tone: 'money', name: 'Amount' },
    { id: `CAS${String(4 + i * 3).padStart(2, '0')}`, start: reasonStart + 23, len: 15, tone: 'money', name: 'Quantity' },
  ]),
]

// Figure 7: where each record's last defined field ends.
const EXTENTS: { seg: string; loop: string; end: number; page: string }[] = [
  { seg: 'ISA', loop: 'envelope', end: 104, page: '1' },
  { seg: 'BPR', loop: 'header', end: 183, page: '1–2' },
  { seg: 'CLP', loop: '2100', end: 196, page: '5–6' },
  { seg: 'NM1', loop: '2100 patient', end: 234, page: '6–8' },
  { seg: 'CAS', loop: '2100 / 2110', end: 248, page: '6, 12' },
  { seg: 'SVC', loop: '2110', end: 328, page: '11–12' },
  { seg: 'MOA', loop: '2100', end: 332, page: '9–10' },
  { seg: 'PLB', loop: 'trailer', end: 496, page: '13–14' },
  { seg: 'MIA', loop: '2100', end: 601, page: '8–9' },
  { seg: 'PER', loop: '1000A / 2100', end: 854, page: '3, 11' },
]

interface Node {
  seg: string
  name: string
  repeat?: string
  note?: string
}

interface Level {
  loop: string
  title: string
  repeat?: string
  page: string
  nodes: Node[]
  children?: Level[]
}

// Figure 3: record order and repeats, copied from the layout's Loop,
// Loop Repeat, and Record Repeat columns.
const TREE: Level = {
  loop: 'ISA · GS · ST',
  title: 'Envelope and transaction header',
  page: '1–2',
  nodes: [
    { seg: 'BPR', name: 'Payment amount, method, dates' },
    { seg: 'TRN', name: 'Check or EFT trace number' },
    { seg: 'REF', name: 'Receiver ID (EV), version (F2)' },
    { seg: 'DTM', name: 'Production date (405)' },
  ],
  children: [
    {
      loop: '1000A',
      title: 'Payer',
      page: '2–3',
      nodes: [
        { seg: 'N1 · N3 · N4', name: 'Name and address' },
        { seg: 'REF', name: 'Additional payer ID', repeat: '4' },
        { seg: 'PER', name: 'Business, technical, web contacts' },
      ],
    },
    {
      loop: '1000B',
      title: 'Payee (provider)',
      page: '3–4',
      nodes: [
        { seg: 'N1 · N3 · N4', name: 'Name, ID, address' },
        { seg: 'REF', name: 'Additional payee ID', repeat: '>1' },
        { seg: 'RDM', name: 'Remittance delivery method' },
      ],
    },
    {
      loop: '2000',
      title: 'Claim group header',
      repeat: '>1',
      page: '4–5',
      nodes: [
        { seg: 'LX', name: 'Header number' },
        { seg: 'TS3', name: 'Provider summary (many fields Part A only)' },
        { seg: 'TS2', name: 'Supplemental summary', note: 'not used for Part B' },
      ],
      children: [
        {
          loop: '2100',
          title: 'One claim',
          repeat: '>1',
          page: '5–11',
          nodes: [
            { seg: 'CLP', name: 'Claim IDs, status, charge, payment' },
            { seg: 'CAS', name: 'Claim adjustment', repeat: '99' },
            { seg: 'NM1', name: 'Patient, corrected name, rendering provider, crossover and priority payer' },
            { seg: 'MIA or MOA', name: 'Inpatient or outpatient adjudication', note: 'one or the other, not both' },
            { seg: 'REF', name: 'Other claim ID', repeat: '5' },
            { seg: 'DTM', name: 'Statement dates (2), expiration, received' },
            { seg: 'PER', name: 'Claim contact', repeat: '2' },
            { seg: 'AMT', name: 'Supplemental amount', repeat: '13' },
            { seg: 'QTY', name: 'Supplemental quantity', repeat: '14' },
          ],
          children: [
            {
              loop: '2110',
              title: 'One service line',
              repeat: '999',
              page: '11–13',
              nodes: [
                { seg: 'SVC', name: 'Procedure, charge, payment, units' },
                { seg: 'DTM', name: 'Service date', repeat: '2' },
                { seg: 'CAS', name: 'Service adjustment', repeat: '99' },
                { seg: 'REF', name: 'Service ID, line control number, rendering provider, policy' },
                { seg: 'AMT', name: 'Service supplemental amount' },
                { seg: 'LQ', name: 'Remark code', repeat: '99' },
              ],
            },
          ],
        },
      ],
    },
    {
      loop: 'PLB · SE · GE · IEA',
      title: 'Provider adjustments and trailers',
      page: '13–15',
      nodes: [
        { seg: 'PLB', name: 'Provider-level adjustment, up to 6 per record' },
        { seg: 'SE', name: 'Segment count; SE02 = ST02' },
        { seg: 'GE', name: 'Count of transactions; GE02 = GS06' },
        { seg: 'IEA', name: 'Count of groups; IEA02 = ISA13' },
      ],
    },
  ],
}

function Badge({ children, tone = 'ink' }: { children: ReactNode; tone?: 'ink' | 'accent' }) {
  const cls = tone === 'accent' ? 'border-accent text-accent' : 'border-line text-ink-2'
  return <span className={`ml-2 inline-block border px-1.5 font-mono text-[10px] leading-4 ${cls}`}>{children}</span>
}

function LoopBox({ level, depth = 0 }: { level: Level; depth?: number }) {
  return (
    <div className={`border ${depth >= 2 ? 'border-accent' : 'border-line'} bg-bg p-3`}>
      <div className="mb-2 flex flex-wrap items-baseline gap-x-2">
        <span className="font-mono text-[12px] font-semibold text-accent">{level.loop}</span>
        <span className="text-[13px] font-semibold text-ink">{level.title}</span>
        {level.repeat && <Badge tone="accent">loop repeats {level.repeat}</Badge>}
        <span className="ml-auto font-mono text-[10px] text-ink-2">p. {level.page}</span>
      </div>
      <ul className="space-y-1">
        {level.nodes.map((n) => (
          <li key={`${n.seg}-${n.name}`} className="flex items-baseline gap-2 text-[12px] leading-5">
            <span className="w-24 shrink-0 font-mono text-ink">{n.seg}</span>
            <span className="text-ink-2">
              {n.name}
              {n.repeat && <Badge>×{n.repeat}</Badge>}
              {n.note && <span className="ml-2 italic">({n.note})</span>}
            </span>
          </li>
        ))}
      </ul>
      {level.children && (
        <div className="mt-3 space-y-3 border-l-2 border-line pl-3">
          {level.children.map((c) => (
            <LoopBox key={c.loop} level={c} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  )
}

// One 18-byte amount field, drawn cell by cell.
function AmountCells({ cells, label }: { cells: string; label: string }) {
  return (
    <div className="flex items-center gap-4">
      <div className="w-44 shrink-0 text-[12px] text-ink-2">{label}</div>
      <div className="flex">
        {cells.split('').map((ch, i) => (
          <div
            key={i}
            className={`flex h-8 w-7 items-center justify-center border border-line font-mono text-[13px] text-ink ${
              i < 10 ? 'bg-ok/15' : 'bg-bg'
            } ${i === 8 ? 'border-l-2 border-l-accent' : ''}`}
          >
            {ch === ' ' ? <span className="text-ink-2">·</span> : ch}
          </div>
        ))}
      </div>
    </div>
  )
}

function Th({ children }: { children: ReactNode }) {
  return <th className="py-2 pr-4 text-left text-xs font-semibold text-ink">{children}</th>
}

function Td({ children, mono }: { children: ReactNode; mono?: boolean }) {
  return <td className={`py-2 pr-4 align-top ${mono ? 'font-mono text-xs text-ink' : 'text-sm'}`}>{children}</td>
}

export function Format835Flatfile() {
  const maxEnd = Math.max(...EXTENTS.map((e) => e.end))
  return (
    <div className="mx-auto max-w-6xl px-10 py-12">
      <PageTitle
        lede={
          <>
            <strong className="text-ink">835 Flat File – 5010A1</strong>, Centers for Medicare &amp; Medicaid
            Services (CMS), updated January 10, 2023 (15 pages).
          </>
        }
      >
        CMS 835 flat-file layout
      </PageTitle>

      <DocSection id="gist" title="The gist">
        <p>
          This is a fixed-position record layout: the same data as an X12 835 remittance advice (the
          payer&rsquo;s explanation of claim payments), but with <strong>one record per segment</strong>,
          each record starting with an <strong>18-byte key</strong> and every field at a fixed start
          position and length. CMS publishes it as the internal Medicare representation of the 835; its
          comments cite conventions of FISS (Fiscal Intermediary Shared System) and MCS (Multi-Carrier
          System) and requests from HIGLAS (Healthcare Integrated General Ledger Accounting System).
        </p>
        <p>
          For the validator it matters because a Medicare administrative contractor (MAC) output may
          arrive in this shape rather than as delimited X12. Every field can be located by byte
          position and mapped back to the X12 element it came from, so the two sides can be compared
          element by element.
        </p>
      </DocSection>

      <DocSection id="row" title="How to read one row of the layout">
        <Figure
          number={1}
          caption="Each layout row has two halves: X12 attributes on the left, flat-file position on the right. Row CLP04 shown (PDF p. 6; column headers p. 1)."
        >
          <div className="grid min-w-[46rem] grid-cols-[3fr_2fr] gap-3">
            <div className="border border-line">
              <div className="border-b border-line bg-surface px-3 py-1.5 text-[11px] font-semibold text-ink">
                X12 element attributes
              </div>
              <div className="grid grid-cols-6 text-[12px]">
                {[
                  ['Element', 'CLP04'],
                  ['Name', 'Claim Payment Amount'],
                  ['COBOL PIC', 'X(18) group'],
                  ['ID (type)', 'R = decimal'],
                  ['Min/Max', '1–18'],
                  ['Usage', 'R = required'],
                ].map(([h, v]) => (
                  <div key={h} className="border-r border-line p-2 last:border-r-0">
                    <div className="text-[10px] text-ink-2">{h}</div>
                    <div className="font-mono text-ink">{v}</div>
                  </div>
                ))}
              </div>
            </div>
            <div className="border-2 border-accent">
              <div className="border-b border-line bg-accent/15 px-3 py-1.5 text-[11px] font-semibold text-ink">
                X12 flat file
              </div>
              <div className="grid grid-cols-4 text-[12px]">
                {[
                  ['Loop ID', '2100'],
                  ['Seg ID', 'CLP'],
                  ['Start', '77'],
                  ['Length', '18'],
                ].map(([h, v]) => (
                  <div key={h} className="border-r border-line p-2 last:border-r-0">
                    <div className="text-[10px] text-ink-2">{h}</div>
                    <div className="font-mono text-ink">{v}</div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <ul className="mt-4 grid min-w-[46rem] grid-cols-2 gap-x-8 gap-y-1 text-[12px] text-ink-2">
            <li><span className="font-mono text-ink">Usage</span>: R required, S situational, N/U not used.</li>
            <li><span className="font-mono text-ink">Page #</span>: page in the X12 implementation guide (TR3, Technical Report Type 3).</li>
            <li><span className="font-mono text-ink">Loop Repeat / Record Repeat</span>: how many times the loop or record can occur.</li>
            <li><span className="font-mono text-ink">Comment</span>: CMS notes, such as &ldquo;Only Part A&rdquo; or &ldquo;Not used by Medicare&rdquo;.</li>
            <li>Rows highlighted yellow in the PDF changed in this edition.</li>
            <li>Most N/U elements have no Start or Length: they take no space in the flat record.</li>
          </ul>
        </Figure>
      </DocSection>

      <DocSection id="key" title="Every record starts with an 18-byte key">
        <Figure
          number={2}
          caption="The record key occupies positions 1–18 (widths 6, 4, 4, 4 in the column header, PDF p. 1). Every segment header row lists Start 1, Length 18; the first data field starts at 19."
        >
          <div className="min-w-[46rem]">
            <div className="flex">
              <div style={{ flexGrow: 18, flexBasis: 0 }}>
                <ByteRuler spans={KEY} minCell="4rem" />
              </div>
              <div
                style={{ flexGrow: 22, flexBasis: 0 }}
                className="flex items-center border border-l-0 border-dashed border-line px-3 font-mono text-[11px] text-ink-2"
              >
                19 → data fields at fixed Start and Length
              </div>
            </div>
            <div className="mt-4 grid grid-cols-3 gap-3 text-[12px]">
              {[
                { rec: 'Payer name', key: '1000A · N1', from: 'N101 at 19, length 3' },
                { rec: 'Claim', key: '2100 · CLP', from: 'CLP01 at 19, length 38' },
                { rec: 'Service line', key: '2110 · SVC', from: 'SVC01-1 at 19, length 2' },
              ].map((r) => (
                <div key={r.key} className="border border-line bg-surface p-2">
                  <div className="text-ink-2">{r.rec}</div>
                  <div className="font-mono text-ink">{r.key}</div>
                  <div className="font-mono text-[11px] text-ink-2">{r.from}</div>
                </div>
              ))}
            </div>
            <p className="mt-3 text-[12px] text-ink-2">
              Records outside the numbered loops (ISA, GS, ST, BPR, TRN, REF, DTM, PLB, SE, GE, IEA) have a blank Loop ID.
              The layout leaves the Loop Seq and Seg Seq columns empty; it does not show example sequence values.
            </p>
          </div>
        </Figure>
      </DocSection>

      <DocSection id="order" title="Record order and repeats">
        <p>Nesting follows the X12 835 loops. Each box is a loop; badges show how many times a record can repeat.</p>
        <Figure
          number={3}
          caption="Loops and records in layout order, with Loop Repeat and Record Repeat values (PDF pp. 1–15). Loop 2100 is one claim; loop 2110 is one service line inside it."
        >
          <div className="min-w-[40rem]">
            <LoopBox level={TREE} />
          </div>
        </Figure>
      </DocSection>

      <DocSection id="claim" title="The claim record, byte by byte">
        <Figure
          number={4}
          caption="CLP record: positions 1–196. CLP10 (patient status) and CLP14 are not used and take no bytes, so CLP11 follows CLP09 directly (PDF pp. 5–6). Hover a cell for the field name."
        >
          <ByteRuler spans={CLP} />
          <Legend />
          <div className="mt-4 grid min-w-[46rem] grid-cols-3 gap-x-6 gap-y-1 text-[12px] text-ink-2">
            {CLP.slice(1).map((s) => (
              <div key={s.id}>
                <span className="font-mono text-ink">{s.id}</span> {s.name}
              </div>
            ))}
          </div>
        </Figure>
        <p className="text-sm">
          Two claim identifiers sit in this record: <code>CLP01</code>, the provider&rsquo;s patient control number
          (38 bytes), and <code>CLP07</code>, the payer claim control number (50 bytes). <code>CLP02</code> allows status
          codes 1, 2, 3, 4, 19, 20, 21, 22, 23.
        </p>
      </DocSection>

      <DocSection id="amounts" title="How an amount is stored">
        <Figure
          number={5}
          caption="Amount fields are an 18-byte group redefined as a signed number with two implied decimals, S9(8)V99, plus an 8-byte filler. Rule from the footnote on every page: zero-filled and right-justified in positions 1–10, spaces in 11–18; no data means 18 spaces, except FISS may zero-fill 1–10 (PDF p. 1). Example values are illustrative."
        >
          <div className="min-w-[46rem] space-y-2">
            <div className="flex items-center gap-4">
              <div className="w-44 shrink-0" />
              <div className="flex font-mono text-[10px] text-ink-2">
                {Array.from({ length: 18 }, (_, i) => (
                  <div key={i} className="w-7 text-center">{i + 1}</div>
                ))}
              </div>
            </div>
            <AmountCells label="80.00" cells="0000008000        " />
            <AmountCells label="1,234.56" cells="0000123456        " />
            <AmountCells label="No data (standard)" cells="                  " />
            <AmountCells label="No data (FISS may)" cells="0000000000        " />
            <div className="flex items-center gap-4 pt-2 text-[12px] text-ink-2">
              <div className="w-44 shrink-0" />
              <div className="flex gap-4">
                <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-3 border border-line bg-ok/15" />S9(8)V99 digits</span>
                <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-0.5 bg-accent" />implied decimal point</span>
                <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-3 border border-line bg-bg" />filler X(8)</span>
              </div>
            </div>
          </div>
        </Figure>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[36rem]">
            <thead className="border-b border-line">
              <tr><Th>Kind</Th><Th>Width</Th><Th>Numeric part</Th><Th>Example</Th></tr>
            </thead>
            <tbody className="[&_tr]:border-b [&_tr]:border-line">
              <tr><Td>Amount</Td><Td mono>X(18)</Td><Td mono>S9(8)V99 + X(8)</Td><Td>CLP04, SVC03, CAS03 (p. 6, 11)</Td></tr>
              <tr><Td>Count</Td><Td mono>X(15)</Td><Td mono>S9(3) + X(12)</Td><Td>MIA03 lifetime psychiatric days (p. 8)</Td></tr>
              <tr><Td>Rate</Td><Td mono>X(10)</Td><Td mono>S9(4)V9999 + X(2)</Td><Td>MOA01 reimbursement rate (p. 10)</Td></tr>
              <tr><Td>Date</Td><Td mono>X(8)</Td><Td mono>9(8) CCYYMMDD</Td><Td>DTM02, BPR16 (pp. 2, 10)</Td></tr>
            </tbody>
          </table>
        </div>
        <p className="text-sm">
          The layout does not say how the sign of a negative amount is written. Confirm that with real sample files before decoding.
        </p>
      </DocSection>

      <DocSection id="cas" title="Repeated elements become fixed slots">
        <Figure
          number={6}
          caption="CAS record: one group code, then six reason / amount / quantity slots at fixed positions, ending at 248. Up to 99 CAS records per claim (2100) and per service line (2110) (PDF pp. 6, 12)."
        >
          <ByteRuler spans={CAS} minCell="2.5rem" />
          <Legend />
          <div className="mt-5 grid min-w-[46rem] grid-cols-[auto_1fr] items-center gap-x-4 gap-y-2 text-[12px]">
            <span className="text-ink-2">X12 segment</span>
            <code className="font-mono text-ink">CAS*PR*1*20.00~</code>
            <span className="text-ink-2">Flat record</span>
            <span className="font-mono text-ink">
              [key 1–18] <span className="bg-series-3/15 px-1">PR</span> at 19 ·{' '}
              <span className="bg-series-3/15 px-1">1&nbsp;&nbsp;&nbsp;&nbsp;</span> at 21 ·{' '}
              <span className="bg-ok/15 px-1">0000002000········</span> at 26 · the rest spaces
            </span>
          </div>
          <p className="mt-3 min-w-[46rem] text-[12px] text-ink-2">
            Illustrative mapping. The same pattern holds for PLB (up to six adjustment slots, p. 13–14), MIA/MOA remark codes,
            and the three PER contact numbers (256 bytes each).
          </p>
        </Figure>
      </DocSection>

      <DocSection id="lengths" title="Records are not one fixed length">
        <Figure
          number={7}
          caption="Last byte of the last defined field, per record type (Start + Length − 1 from the layout). The document does not state a total record length or padding rule."
        >
          <svg viewBox="0 0 720 300" className="w-full min-w-[40rem]" role="img" aria-label="Bar chart of record extents">
            {EXTENTS.map((e, i) => {
              const y = 10 + i * 28
              const w = (e.end / maxEnd) * 480
              return (
                <g key={e.seg}>
                  <text x={0} y={y + 14} className="fill-current text-ink font-mono" fontSize={12}>{e.seg}</text>
                  <text x={40} y={y + 14} className="fill-current text-ink-2" fontSize={10}>{e.loop}</text>
                  <rect x={130} y={y} width={(18 / maxEnd) * 480} height={18} style={{ fill: 'var(--accent)' }} opacity={0.45} />
                  <rect x={130 + (18 / maxEnd) * 480} y={y} width={w - (18 / maxEnd) * 480} height={18} style={{ fill: 'var(--series-1)' }} opacity={0.2} />
                  <text x={136 + w} y={y + 14} className="fill-current text-ink font-mono" fontSize={11}>{e.end}</text>
                  <text x={720} y={y + 14} textAnchor="end" className="fill-current text-ink-2" fontSize={10}>p. {e.page}</text>
                </g>
              )
            })}
          </svg>
        </Figure>
      </DocSection>

      <DocSection id="validator" title="What this means for the validator">
        <ul className="list-disc space-y-2 pl-5">
          <li>
            <strong>Decode per record type.</strong> Read the 18-byte key, pick the layout row set for that Loop ID and
            Seg ID, then slice fields by Start and Length. Records differ in length (Figure 7).
          </li>
          <li>
            <strong>Normalize before comparing with X12.</strong> Trim fixed-width padding, turn the implied-decimal
            amount into a decimal value, and treat an all-space amount as absent. A zero-filled blank from FISS is
            ambiguous with a real 0.00; keep the original bytes as evidence.
          </li>
          <li>
            <strong>Widths differ from X12.</strong> Flat widths are storage widths: names are widened to 60 bytes and TRN04
            from 30 to 50. Where a row&rsquo;s PIC and Length disagree (NM105 shows X(2) with Length 25, p. 7), slice by Start and Length.
          </li>
          <li>
            <strong>The key is not a claim identity.</strong> Loop and segment sequence numbers are positions within a file.
            Match claims on the six-field key: the beneficiary identifier, the date of service, the total charge, the payee's tax ID and NPI, and the MAC workload; CLP01 and CLP07 are not in it, since each system assigns its own.
          </li>
          <li>
            <strong>Part A and Part B fill different records.</strong> MIA or MOA (never both), TS2 not used for Part B,
            many TS3 fields Part A only. Missing records can be expected, not errors.
          </li>
          <li>
            <strong>Pin the edition.</strong> This is the January 10, 2023 layout. Confirm it and the sign encoding against
            supplied files before fixing offsets in code.
          </li>
        </ul>
      </DocSection>

      <DocSection id="where" title="Where to look in the PDF">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[36rem]">
            <thead className="border-b border-line">
              <tr><Th>Topic</Th><Th>Records</Th><Th>PDF page</Th></tr>
            </thead>
            <tbody className="[&_tr]:border-b [&_tr]:border-line">
              <tr><Td>Column headers, record key widths, amount-field rule</Td><Td mono>header, footnote</Td><Td>1</Td></tr>
              <tr><Td>Interchange, group, transaction headers; payment</Td><Td mono>ISA GS ST BPR</Td><Td>1–2</Td></tr>
              <tr><Td>Trace, receiver, version, production date</Td><Td mono>TRN REF DTM</Td><Td>2</Td></tr>
              <tr><Td>Payer (1000A)</Td><Td mono>N1 N3 N4 REF PER</Td><Td>2–3</Td></tr>
              <tr><Td>Payee (1000B)</Td><Td mono>N1 N3 N4 REF RDM</Td><Td>3–4</Td></tr>
              <tr><Td>Claim group header (2000)</Td><Td mono>LX TS3 TS2</Td><Td>4–5</Td></tr>
              <tr><Td>Claim (2100)</Td><Td mono>CLP CAS NM1</Td><Td>5–8</Td></tr>
              <tr><Td>Inpatient and outpatient adjudication</Td><Td mono>MIA MOA</Td><Td>8–10</Td></tr>
              <tr><Td>Claim references, dates, contact, supplemental</Td><Td mono>REF DTM PER AMT QTY</Td><Td>10–11</Td></tr>
              <tr><Td>Service line (2110)</Td><Td mono>SVC DTM CAS REF AMT LQ</Td><Td>11–13</Td></tr>
              <tr><Td>Provider-level adjustments</Td><Td mono>PLB</Td><Td>13–14</Td></tr>
              <tr><Td>Trailers and control numbers</Td><Td mono>SE GE IEA</Td><Td>14–15</Td></tr>
            </tbody>
          </table>
        </div>
      </DocSection>
    </div>
  )
}
