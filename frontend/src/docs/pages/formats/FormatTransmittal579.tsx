import type { ReactNode } from 'react'
import { DocSection } from '../../../components/DocSection'
import { Figure } from '../../../components/Figure'
import { PageTitle } from '../../../components/PageTitle'

// CMS Transmittal 579 (Change Request 6575, October 19, 2009), boiled down.
// Every fact below comes from CMS Transmittal 579 (2009).
// Page numbers are PDF viewer pages, counted from 1.


// ---------- small building blocks ----------

function Box({ title, sub, tone = 'plain' }: { title: string; sub?: string; tone?: 'plain' | 'accent' | 'flag' }) {
  const border = tone === 'accent' ? 'border-accent' : tone === 'flag' ? 'border-flag' : 'border-line'
  return (
    <div className={`min-w-[9rem] flex-1 border-2 ${border} bg-surface px-3 py-2`}>
      <div className="text-[13px] font-semibold text-ink">{title}</div>
      {sub && <div className="mt-0.5 text-[12px] leading-5 text-ink-2">{sub}</div>}
    </div>
  )
}

function Arrow({ label, down = false }: { label?: string; down?: boolean }) {
  return (
    <div className={`flex shrink-0 items-center justify-center text-ink-2 ${down ? 'flex-col py-1' : 'flex-col px-1'}`}>
      {label && <span className="font-mono text-[10px] leading-4">{label}</span>}
      <span className="font-mono text-lg leading-5">{down ? '↓' : '→'}</span>
    </div>
  )
}

// One cell of a byte ruler: width proportional to byte length.
interface Span {
  start: number
  len: number
  label: string
  value?: string
  tone?: 'key' | 'data' | 'id' | 'count' | 'blank'
}

const TONE: Record<NonNullable<Span['tone']>, string> = {
  key: 'bg-surface border-ink-2',
  data: 'bg-bg border-line',
  id: 'bg-bg border-accent',
  count: 'bg-bg border-ok',
  blank: 'bg-bg border-line border-dashed',
}

function Ruler({ spans, total, unit = 14 }: { spans: Span[]; total: number; unit?: number }) {
  return (
    <div style={{ width: total * unit }} className="min-w-full">
      <div className="flex">
        {spans.map((s) => (
          <div
            key={s.start}
            style={{ width: s.len * unit }}
            className={`shrink-0 border-2 ${TONE[s.tone ?? 'data']} -ml-[2px] first:ml-0 px-1 py-1`}
          >
            {s.value !== undefined && (
              <div className="truncate font-mono text-[12px] text-ink" title={s.value}>
                {s.value === '' ? ' ' : s.value}
              </div>
            )}
            <div className="truncate text-[11px] leading-4 text-ink-2" title={s.label}>{s.label}</div>
          </div>
        ))}
      </div>
      <div className="flex">
        {spans.map((s) => (
          <div key={s.start} style={{ width: s.len * unit }} className="shrink-0 pt-1 font-mono text-[10px] text-ink-2">
            {s.start}
            {s.len > 1 ? `–${s.start + s.len - 1}` : ''}
          </div>
        ))}
      </div>
    </div>
  )
}

function Legend({ items }: { items: { tone: NonNullable<Span['tone']>; label: string }[] }) {
  return (
    <div className="mt-4 flex flex-wrap gap-4 text-[12px] text-ink-2">
      {items.map((i) => (
        <span key={i.label} className="flex items-center gap-1.5">
          <span className={`inline-block h-3 w-5 border-2 ${TONE[i.tone]}`} />
          {i.label}
        </span>
      ))}
    </div>
  )
}

function Seg({ children, tone = 'plain', note }: { children: ReactNode; tone?: 'plain' | 'internal' | 'claim'; note?: string }) {
  const cls =
    tone === 'internal'
      ? 'border-flag text-ink'
      : tone === 'claim'
        ? 'border-accent text-ink'
        : 'border-line text-ink-2'
  return (
    <div className="flex items-baseline gap-3">
      <div className={`border-l-4 ${cls} bg-surface px-2 py-0.5 font-mono text-[12px]`}>{children}</div>
      {note && <div className="text-[12px] text-ink-2">{note}</div>}
    </div>
  )
}

// ---------- figure data ----------

const JURISDICTIONS = Array.from({ length: 15 }, (_, i) => i + 1)
const IN_SCOPE = new Set([1, 3, 4, 5, 10, 13, 14])

const REQUIREMENTS = [
  { id: '6575.1', who: 'A/B MACs, CEDI', text: 'Translate each syntactically correct 837P into a Medicare Part B claim flat file.' },
  { id: '6575.2', who: 'A/B MACs', text: 'Translate each syntactically correct 837I into a Medicare Part A claim flat file.' },
  { id: '6575.3', who: 'A/B MACs, CEDI', text: 'Run “alpha” testing with contractor or translator-vendor test data.' },
  { id: '6575.4', who: 'A/B MACs, CEDI', text: 'Produce test 837 flat files in the attached formats.' },
  { id: '6575.5', who: 'FISS, MCS, VMS', text: 'Shared systems use the attached spreadsheets to implement.' },
  { id: '6575.6', who: 'A/B MACs, DME MACs, CEDI', text: 'Review the draft receipt and control record; final layout comes in later change requests.' },
]

const KEY_SPANS: Span[] = [
  { start: 1, len: 6, label: 'Loop ID', value: '2300', tone: 'key' },
  { start: 7, len: 4, label: 'Loop seq', value: '0001', tone: 'key' },
  { start: 11, len: 4, label: 'Seg ID', value: 'STC', tone: 'key' },
  { start: 15, len: 4, label: 'Seg seq', value: '0120', tone: 'key' },
]

// STC error record from the sample on PDF p. 116, positioned by the STC layout on pp. 116-117.
const STC_SPANS: Span[] = [
  ...KEY_SPANS,
  { start: 19, len: 5, label: 'STC01-1 category', value: 'A7', tone: 'data' },
  { start: 24, len: 5, label: 'STC01-2 status', value: '697', tone: 'data' },
  { start: 29, len: 3, label: 'STC01-3 entity', value: '', tone: 'blank' },
  { start: 32, len: 8, label: 'STC02 date', value: '20090819', tone: 'data' },
  { start: 40, len: 2, label: 'STC03', value: 'U', tone: 'data' },
  { start: 42, len: 10, label: 'STC04 charges S9(8)V99', value: '000000925A', tone: 'count' },
  { start: 52, len: 5, label: 'STC10-1', value: 'A7', tone: 'data' },
  { start: 57, len: 5, label: 'STC10-2', value: '178', tone: 'data' },
  { start: 62, len: 3, label: 'STC10-3', value: '', tone: 'blank' },
]

// Draft receipt and control record, 350 bytes, PDF pp. 246-251. Grouped by purpose.
const CTR_SPANS: Span[] = [
  { start: 1, len: 4, label: 'CTR01 record ID', tone: 'key' },
  { start: 5, len: 7, label: 'CTR02–04 type, direction, line', tone: 'data' },
  { start: 12, len: 57, label: 'CTR05–11 copied from ISA: sender, receiver, date, time, ISA13 control number, version, test/prod', tone: 'id' },
  { start: 69, len: 58, label: 'CTR12–16 data centers, contractor ID, cycle date, time', tone: 'data' },
  { start: 127, len: 30, label: 'CTR17 local data center tracking', tone: 'data' },
  { start: 157, len: 46, label: 'CTR18–25 SENT totals + error switches', tone: 'count' },
  { start: 203, len: 46, label: 'CTR26–33 RECEIVED totals + error switches', tone: 'count' },
  { start: 249, len: 47, label: 'CTR34 claim control number exchange area', tone: 'id' },
  { start: 296, len: 55, label: 'CTR35 filler (spaces)', tone: 'blank' },
]

const TOTALS = [
  { what: 'Records (ISA–IEA, not counting CTR)', sent: 'CTR18 · 157–167 · 9(11)', recv: 'CTR26 · 203–213' },
  { what: 'Units: CLM for 837, CLP for 835, TRN for 276/277', sent: 'CTR20 · 170–178 · 9(9)', recv: 'CTR28 · 216–224' },
  { what: 'Submitted-charge dollars (837 and pharmacy only)', sent: 'CTR22 · 181–193 · 9(11)v99', recv: 'CTR30 · 227–239' },
  { what: 'ST–SE transaction pairs', sent: 'CTR24 · 196–200 · 9(5)', recv: 'CTR32 · 242–246' },
]

const CCN_SPLITS = [
  { sys: 'FISS', parts: [{ len: 23, label: 'CCN 23' }, { len: 23, label: 'CCN 23' }, { len: 1, label: '' }] },
  { sys: 'MCS', parts: [{ len: 13, label: '13' }, { len: 10, label: 'sp' }, { len: 13, label: '13' }, { len: 10, label: 'sp' }, { len: 1, label: '' }] },
  { sys: 'VMS', parts: [{ len: 15, label: '15' }, { len: 9, label: 'sp' }, { len: 15, label: '15' }, { len: 9, label: 'sp' }] },
]

const PDF_MAP = [
  { what: 'Summary, dates, jurisdictions, background on the ICN REF and STC tab', pages: '1–3' },
  { what: 'Business requirements 6575.1–6575.6', pages: '3–4' },
  { what: '837P (professional) flat-file layout', pages: '6–112' },
  { what: '837P internal DTP +RC and REF +PR after GS', pages: '7' },
  { what: '837P claim loop 2300 CLM; REF F8 payer claim control number', pages: '26, 34' },
  { what: '837P internal REF +CN in loop 2300', pages: '39' },
  { what: '837P legend (18-byte key), change log', pages: '113–115' },
  { what: '837P STC error tab: sample records, STC layout, 999 and 277CA examples', pages: '116–119' },
  { what: '837I (institutional) flat-file layout', pages: '122–235' },
  { what: '837I internal DTP +RC / REF +PR; CLM; REF +CN', pages: '123, 137, 143' },
  { what: '837I change log, legend, STC error tab', pages: '236–243' },
  { what: 'Draft receipt and control record: 276/277, 837 and pharmacy, 835 variants', pages: '246–251' },
]

// ---------- page ----------

export function FormatTransmittal579() {
  return (
    <div className="mx-auto max-w-[90rem] px-10 py-12">
      <PageTitle
        lede={
          <>
            <strong className="text-ink">Transmittal 579</strong> is a 2009 Centers for Medicare &amp; Medicaid
            Services (CMS) instruction that told Medicare contractors how to turn each incoming X12 version 5010
            837 claim into a fixed-position “flat file” for the Medicare shared processing systems. It also added
            two Medicare-internal pieces: a front-end internal control number (ICN) on every claim and STC error
            records, plus a draft receipt-and-control record for balancing file totals. For a validator comparing
            contractor-produced 835s with vendor-produced 835s, it shows which claim identifiers exist inside
            Medicare, who assigns them, and how files are balanced by claim count and dollars.
          </>
        }
      >
        Transmittal 579: 837 flat files
      </PageTitle>

      <p className="-mt-4 mb-10 text-[13px] text-ink-2">
        Source: <em>CMS Manual System, Pub 100-20 One-Time Notification, Transmittal 579, Change Request 6575</em>,
        October 19, 2009; layouts “as of Oct 15, 2009” (837P) and “as of Oct 01, 2009” (837I). Historical:
        later transmittals changed these layouts.
      </p>

      <DocSection id="what" title="What it set in motion">
        <Figure
          number={1}
          caption="Translation path required by 6575.1–6575.2 and the STC error path described in the attachment tab. Acronyms: MAC, Medicare Administrative Contractor; CEDI, Common Electronic Data Interchange (the durable medical equipment front end); COTS, commercial off-the-shelf; CEM, Common Edits and Enhancements Module; FISS, MCS, VMS, the Medicare shared claim-processing systems. PDF pp. 3–4, 116–118."
        >
          <div className="min-w-[46rem] space-y-3">
            <div className="flex items-stretch">
              <Box title="Provider" sub="Sends X12 5010 837P (professional) or 837I (institutional)" />
              <Arrow label="837" />
              <Box title="COTS translator" sub="At the A/B MAC or CEDI front end" tone="accent" />
              <Arrow label="flat file" />
              <Box title="837P → Part B claim flat file" sub="837I → Part A claim flat file" tone="accent" />
              <Arrow />
              <Box title="Shared system" sub="FISS, MCS, VMS" />
            </div>
            <div className="flex items-stretch">
              <div className="min-w-[9rem] flex-1" />
              <div className="w-8 shrink-0" />
              <div className="flex min-w-[9rem] flex-1 flex-col items-center">
                <Arrow down label="errors as STC records" />
              </div>
              <div className="w-8 shrink-0" />
              <div className="min-w-[9rem] flex-1" />
              <div className="w-6 shrink-0" />
              <div className="min-w-[9rem] flex-1" />
            </div>
            <div className="flex items-stretch">
              <Box title="Acknowledgments" sub="999 (syntax) and 277CA (claim) examples on pp. 117–118" tone="flag" />
              <div className="flex shrink-0 items-center px-1 font-mono text-lg text-ink-2">←</div>
              <Box title="CEM" sub="Turns STC records into 277CA STC segments; adds its own edits (e.g. out of balance)" tone="flag" />
              <div className="flex-[2]" />
            </div>
          </div>
        </Figure>

        <Figure
          number={2}
          caption="Who had to act. Transmittal 579 replaced Transmittal 575 (October 9, 2009) to add FISS. Effective January 1, 2010 (implementation January 4) for A/B MACs, CEDI, FISS, and MCS; April 1, 2010 (implementation April 5) for VMS. Only seven A/B jurisdictions were instructed; the rest were left to a future change request. PDF pp. 1, 3."
        >
          <div className="min-w-[36rem]">
            <div className="grid grid-cols-[repeat(15,minmax(0,1fr))] gap-1">
              {JURISDICTIONS.map((j) => (
                <div
                  key={j}
                  className={`border-2 py-2 text-center font-mono text-[12px] ${IN_SCOPE.has(j) ? 'border-accent bg-surface font-semibold text-ink' : 'border-dashed border-line text-ink-2'}`}
                >
                  J{j}
                </div>
              ))}
            </div>
            <div className="mt-3 flex gap-6 text-[12px] text-ink-2">
              <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-5 border-2 border-accent bg-surface" /> prepare for 5010 837 now</span>
              <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-5 border-2 border-dashed border-line" /> not affected by this change request</span>
            </div>
            <div className="mt-5 grid gap-px border border-line bg-line sm:grid-cols-2 lg:grid-cols-3">
              {REQUIREMENTS.map((r) => (
                <div key={r.id} className="bg-bg p-3">
                  <div className="flex items-baseline justify-between gap-2">
                    <span className="font-mono text-[12px] font-semibold text-ink">{r.id}</span>
                    <span className="text-[11px] text-ink-2">{r.who}</span>
                  </div>
                  <p className="mt-1 text-[13px] leading-5 text-ink-2">{r.text}</p>
                </div>
              ))}
            </div>
          </div>
        </Figure>
      </DocSection>

      <DocSection id="flat-file" title="What a flat-file record looks like">
        <ul className="list-disc space-y-1 pl-5">
          <li>One X12 segment becomes one fixed-width record. Envelope segments (ISA, GS, ST, SE, GE, IEA) become records too.</li>
          <li>Every record starts with an 18-byte key; data fields follow at the layout’s <em>Start</em> and <em>Length</em> columns, from byte 19.</li>
          <li>Each field carries a COBOL picture clause, for example amounts as <code>S9(5)V99</code>: implied decimal point, no dot in the data.</li>
          <li>X12 size and flat size can differ: 837P CLM01 allows 1–38 characters in X12 but gets 20 flat bytes (p. 26).</li>
        </ul>

        <Figure
          number={3}
          caption="The 18-byte record key from the legend: loop ID left-justified and space-filled; loop sequence right-justified and zero-filled; segment ID left-justified; segment sequence zero-filled. Example values from the sample STC record. PDF p. 113 (837P) and p. 238 (837I)."
        >
          <Ruler spans={KEY_SPANS} total={18} unit={40} />
        </Figure>

        <Figure
          number={4}
          caption="Where claims live in the 837P flat file, with the Medicare-internal records marked. The internal DTP and REF follow GS; REF +PR takes P (paper), K (keyshop), or O (optical character recognition). REF +CN sits inside the claim loop; the change log says an “ICN REF” was added (23 bytes, later widened to 50) and the background says it carries the front-end ICN to the shared system. CLM loop repeats up to 100 times; LX service lines up to 50. PDF pp. 7, 26, 34, 39, 78, 114–115."
        >
          <div className="min-w-[40rem] space-y-1.5">
            <Seg>ISA · interchange header</Seg>
            <Seg>GS · functional group header</Seg>
            <Seg tone="internal" note="internal: DTP02 D8, DTP03 CCYYMMDD date">DTP +RC</Seg>
            <Seg tone="internal" note="internal: how the claim arrived: P, K, O">REF +PR</Seg>
            <Seg>ST · transaction header · BHT</Seg>
            <div className="ml-6 space-y-1.5 border-l border-line pl-4">
              <Seg>1000A submitter · 1000B receiver</Seg>
              <Seg>2000A / 2010AA billing provider</Seg>
              <Seg>2000B / 2010BA subscriber · 2010BB payer</Seg>
              <div className="ml-6 space-y-1.5 border-l border-line pl-4">
                <Seg tone="claim" note="CLM01 patient account number: X12 1–38, flat 20 bytes at 19">2300 CLM · one claim (repeat 100)</Seg>
                <Seg tone="claim" note="payer claim control number, 50 bytes">REF F8</Seg>
                <Seg tone="internal" note="internal claim-level REF; REF02 50 bytes">REF +CN</Seg>
                <Seg tone="internal" note="added after any record the translator found in error">STC …</Seg>
                <Seg>2310x / 2320 / 2330x providers, other payers</Seg>
                <div className="ml-6 space-y-1.5 border-l border-line pl-4">
                  <Seg>2400 LX · SV1 · one service line (repeat 50)</Seg>
                </div>
              </div>
            </div>
            <Seg>SE · GE · IEA · trailers</Seg>
            <div className="mt-3 flex flex-wrap gap-4 text-[12px] text-ink-2">
              <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-1 bg-flag" /> Medicare-internal, not from the provider</span>
              <span className="flex items-center gap-1.5"><span className="inline-block h-3 w-1 bg-accent" /> claim identity</span>
            </div>
          </div>
        </Figure>
      </DocSection>

      <DocSection id="stc" title="How translator errors travel: STC records">
        <p>
          When the translator finds an error, it adds an STC record after the segment in error. The CEM
          passes those records on as STC segments in the 277CA claim acknowledgment.
        </p>
        <Figure
          number={5}
          caption="One provider error traced through each stage, from the tab’s worked example. CLM02 92.511 has too many decimal places. A7 = acknowledgment / rejected for invalid information; 697 = invalid decimal precision; 178 = submitted charges. In the flat record, 92.51 is stored as 000000925A under PIC S9(8)V99; the final byte carries the last digit and the sign (signed zoned decimal; the document does not explain this encoding). The 999’s segment position 120 matches the key’s segment sequence 0120 in this sample. PDF pp. 116–118."
        >
          <div className="min-w-[48rem] space-y-4">
            <div className="grid grid-cols-[9rem_1fr] items-baseline gap-x-4 gap-y-3">
              <div className="text-[12px] font-semibold text-ink">1 · Provider 837P</div>
              <code className="font-mono text-[12px] text-ink">CLM&gt;PARTB&gt;<mark className="bg-transparent font-semibold text-flag underline">92.511</mark>&gt;&gt;&gt;11+B+1&gt;Y&gt;A&gt;Y&gt;Y&gt;P</code>
              <div className="text-[12px] font-semibold text-ink">2 · 999 response</div>
              <code className="font-mono text-[12px] text-ink">IK3&gt;CLM&gt;<span className="font-semibold text-accent">120</span>&gt;&gt;8~ IK4&gt;2&gt;782&gt;I12&gt;92.511~ … IK5&gt;E~</code>
              <div className="text-[12px] font-semibold text-ink">3 · Flat STC record</div>
              <div />
            </div>
            <Ruler spans={STC_SPANS} total={64} unit={14} />
            <Legend items={[{ tone: 'key', label: '18-byte key' }, { tone: 'data', label: 'status data' }, { tone: 'count', label: 'amount' }, { tone: 'blank', label: 'spaces' }]} />
            <div className="grid grid-cols-[9rem_1fr] items-baseline gap-x-4">
              <div className="text-[12px] font-semibold text-ink">4 · 277CA from CEM</div>
              <code className="font-mono text-[12px] text-ink">STC&gt;A7+697&gt;20090721&gt;U&gt;92.51&gt;&gt;&gt;&gt;&gt;&gt;A7+178~</code>
            </div>
          </div>
        </Figure>
      </DocSection>

      <DocSection id="control" title="The draft receipt and control record">
        <ul className="list-disc space-y-1 pl-5">
          <li>A 350-byte CTR record describes one flat-file transmission between data centers, with sent and received totals for balancing.</li>
          <li>Three variants share one layout: 276/277 status (pp. 246–247), 837 and pharmacy (pp. 248–249), 835 (pp. 250–251).</li>
          <li>CTR01 kinds: <code>CTRD</code> detail, <code>CTRE</code> last scheduled claims transmission of the day, <code>CTRR</code> resubmission to fix an out-of-balance file.</li>
          <li>Marked “Revised 4/13/2009 – Option B.” Requirement 6575.6 says the layout would be finalized in later change requests.</li>
        </ul>

        <Figure number={6} caption="The 350-byte control record, fields grouped by purpose. Positions are inclusive and start at 1. PDF pp. 246–251.">
          <Ruler spans={CTR_SPANS} total={350} unit={3.2} />
          <Legend items={[{ tone: 'key', label: 'record type' }, { tone: 'id', label: 'identifiers' }, { tone: 'count', label: 'balancing totals' }, { tone: 'blank', label: 'filler' }]} />
        </Figure>

        <Figure
          number={7}
          caption="Balancing: each total appears twice, as sent and as received, and each has a 2-byte error switch next to it (X0 processed, X1 out of balance, R0 resubmitted, R1 resubmission out of balance). For an 835 file the unit counted is the CLP claim segment. PDF pp. 246–251."
        >
          <table className="w-full min-w-[40rem] border-collapse text-[13px]">
            <thead>
              <tr className="border-b-2 border-line text-left text-ink">
                <th className="py-1.5 pr-4 font-semibold">Total</th>
                <th className="py-1.5 pr-4 font-semibold">Sent · bytes · picture</th>
                <th className="py-1.5 font-semibold">Received · bytes</th>
              </tr>
            </thead>
            <tbody className="text-ink-2">
              {TOTALS.map((t) => (
                <tr key={t.what} className="border-b border-line">
                  <td className="py-1.5 pr-4 text-ink">{t.what}</td>
                  <td className="py-1.5 pr-4 font-mono text-[12px]">{t.sent}</td>
                  <td className="py-1.5 font-mono text-[12px]">{t.recv}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </Figure>
      </DocSection>

      <DocSection id="identity" title="Claim identity: three different numbers">
        <Figure
          number={8}
          caption="CTR34 (bytes 249–295, 47 bytes) is “hold for exchange of the CCN numbers assigned to 5010 claims between MAC and SSM 4010 process.” CCN is claim control number; SSM is shared system maintainer. Each shared system uses its own width, stored as two slots. PDF pp. 247, 249, 251."
        >
          <div className="min-w-[36rem] space-y-2">
            {CCN_SPLITS.map((row) => (
              <div key={row.sys} className="flex items-center gap-3">
                <div className="w-12 font-mono text-[12px] font-semibold text-ink">{row.sys}</div>
                <div className="flex" style={{ width: 47 * 12 }}>
                  {row.parts.map((p, i) => (
                    <div
                      key={i}
                      style={{ width: p.len * 12 }}
                      className={`shrink-0 border-2 py-1 text-center font-mono text-[11px] ${p.label === 'sp' || p.label === '' ? 'border-dashed border-line text-ink-2' : 'border-accent bg-surface text-ink'} -ml-[2px] first:ml-0`}
                    >
                      {p.label}
                    </div>
                  ))}
                </div>
              </div>
            ))}
            <div className="grid gap-px border border-line bg-line pt-0 sm:grid-cols-3">
              <div className="bg-bg p-3 text-[13px]"><div className="font-semibold text-ink">CLM01 patient account number</div><div className="text-ink-2">Assigned by the provider. Flat 20 bytes (837P p. 26; 837I p. 137).</div></div>
              <div className="bg-bg p-3 text-[13px]"><div className="font-semibold text-ink">Front-end ICN (REF +CN)</div><div className="text-ink-2">Assigned by the receiving front end, passed to the shared system. REF02 50 bytes (pp. 3, 39, 115).</div></div>
              <div className="bg-bg p-3 text-[13px]"><div className="font-semibold text-ink">Shared-system CCN (CTR34)</div><div className="text-ink-2">Exchanged between contractor and shared-system processes; 13, 15, or 23 bytes (pp. 247–251).</div></div>
            </div>
          </div>
        </Figure>
      </DocSection>

      <DocSection id="validator" title="What this means for the validator">
        <ul className="list-disc space-y-2 pl-5">
          <li><strong>The 18-byte key locates a record, not a claim.</strong> It says which loop and segment occurrence a record is. It is no business identifier across systems.</li>
          <li><strong>Medicare adds its own claim numbers after receipt.</strong> The front-end ICN and shared-system CCNs are assigned inside Medicare. A system processing the same 837 elsewhere would not produce the same values, so these cannot be the cross-system match key without an agreement.</li>
          <li><strong>Balancing by claim count and dollars is an existing practice.</strong> The control record counts CLM segments in an 837 and CLP segments in an 835, and totals submitted charges. The same checks per pair can flag missing or extra claims before field comparison.</li>
          <li><strong>Flat numbers are not X12 numbers.</strong> Implied decimals and a sign carried in the last byte (000000925A for 92.51) must be decoded to exact decimals before comparing against X12 values.</li>
          <li><strong>This is historical and draft.</strong> The 837 layouts were revised later, and the control record was “draft” with finalization deferred. Confirm the actual edition of any supplied flat file before using these offsets.</li>
        </ul>
      </DocSection>

      <DocSection id="where" title="Where to look in the PDF">
        <table className="w-full border-collapse text-[13px]">
          <thead>
            <tr className="border-b-2 border-line text-left text-ink">
              <th className="py-1.5 pr-4 font-semibold">Topic</th>
              <th className="py-1.5 font-semibold">PDF pages</th>
            </tr>
          </thead>
          <tbody className="text-ink-2">
            {PDF_MAP.map((row) => (
              <tr key={row.what} className="border-b border-line">
                <td className="py-1.5 pr-4">{row.what}</td>
                <td className="py-1.5 font-mono text-[12px] text-ink">{row.pages}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </DocSection>
    </div>
  )
}
