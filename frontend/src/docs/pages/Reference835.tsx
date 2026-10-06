import { Link } from 'react-router'
import { DocSection } from '../../components/DocSection'
import { PageTitle } from '../../components/PageTitle'
import { CEDI_835_GUIDE, CMS_FIELDS, CMS_MANUAL, CMS_MANUAL_2003, CMS_MLN_FAQ, CMS_MREP_GUIDE, CMS_OVERVIEW, CMS_TI_GUIDE, X12_CARC, X12_EXAMPLES, X12_GROUP_CODES, X12_RARC } from '../../data/lookupRows835'
import {
  BPR_EL_NAME, CAS01_EL_NAME, CAS_AMOUNT_EL_NAME, CAS_GROUP_NAME, CAS_REASON_EL_NAME, CLP_EL_NAME, DTM01_EL_NAME, DTM_NAME, GE_EL_NAME,
  GS_EL_NAME, ID_QUALIFIER_NAME, IEA_EL_NAME, ISA_EL_NAME, LQ_EL_NAME, MOA_EL_NAME, N1_NAME, N1_PE_EL_NAME, NM1_QC_EL_NAME, PLB_AMOUNT_EL_NAME,
  SE_EL_NAME, SERVICE_REF_NAME, ST_EL_NAME, SVC_EL_NAME, SVC_QUALIFIER_NAME, TRN_EL_NAME, TRN_NAME,
} from '../../data/x12Names835'
import { APP_URL } from '../appUrl'

// Every name on this page is a public document's, cited where it appears:
// "CMS 835 Flat File, PDF p.N" counts PDF viewer pages of
// docs/reference/file-formats/cms-835-flatfile.pdf; "CEDI 835 CG, PDF p.N
// (printed p.M)" is the CEDI 835 companion guide. The names themselves are
// read from frontend/src/data/x12Names835.ts, where each carries its page.

const LINK = 'text-link underline decoration-line underline-offset-4 hover:decoration-current'

function Source({ href, children }: { href: string; children: string }) {
  return <a href={href} target="_blank" rel="noreferrer" className={LINK}>{children}</a>
}

function Example({ id, children }: { id: string; children: string }) {
  return <Link to={`/835-viewer?example=${id}`} className={LINK}>{children}</Link>
}

const cms = (pages: string) => `CMS 835 Flat File, PDF p.${pages}`

// The amounts an 835 carries, by element, with the flat file's name for each.
const AMOUNTS: [string, string, string][] = [
  ['BPR02', BPR_EL_NAME[1], cms('1')],
  ['CLP03', CLP_EL_NAME[2], cms('5')],
  ['CLP04', CLP_EL_NAME[3], cms('6')],
  ['CLP05', CLP_EL_NAME[4], cms('6')],
  ['CAS03, CAS06, …', CAS_AMOUNT_EL_NAME, cms('6 (claim), 12–13 (service line)')],
  ['SVC02', SVC_EL_NAME[1], cms('11')],
  ['SVC03', SVC_EL_NAME[2], cms('11')],
  ['PLB04, PLB06, …', PLB_AMOUNT_EL_NAME, cms('13–14')],
]

const IDENTIFIERS: [string, string, string][] = [
  ['ISA13 ↔ IEA02', `${ISA_EL_NAME[12]} ↔ ${IEA_EL_NAME[1]}`, cms('1, 15')],
  ['GS06 ↔ GE02', `${GS_EL_NAME[5]} ↔ ${GE_EL_NAME[1]}`, cms('2, 14')],
  ['ST02 ↔ SE02', `${ST_EL_NAME[1]} ↔ ${SE_EL_NAME[1]}`, cms('1, 14')],
  ['TRN02, TRN03', `${TRN_EL_NAME[1]}; ${TRN_EL_NAME[2]} (segment: ${TRN_NAME['1']})`, cms('2')],
  ['N104, with N103 = XX', `${N1_PE_EL_NAME[3]}; XX is the ${ID_QUALIFIER_NAME.XX}`, `${cms('3, 7')}; CEDI 835 CG, PDF p.32 (printed p.26)`],
  ['CLP01', CLP_EL_NAME[0], cms('5')],
  ['CLP07', CLP_EL_NAME[6], cms('6')],
  ['NM109, with NM101 = QC', NM1_QC_EL_NAME[6], cms('6–7')],
  ['REF02, with REF01 = 6R', SERVICE_REF_NAME['6R'], cms('13')],
]

// The CEDI 835 companion guide's acronym list, Table 26 (PDF p.31–32,
// printed p.25–26), for the acronyms the 835 pages use; CARC and RARC are
// the titles of X12's public code lists.
const CEDI_ACRONYMS: [string, string][] = [
  ['835', '835 Electronic Remittance Advice transaction'],
  ['ERA', 'Electronic Remittance Advice'],
  ['EDI', 'Electronic Data Interchange'],
  ['ISA/IEA', 'ISA – Interchange Control Header / IEA – Interchange Control Trailer'],
  ['GS/GE', 'GS – Functional Group Header / GE – Functional Group Trailer'],
  ['ST/SE', 'ST – Transaction Set Header / SE – Transaction Set Trailer'],
  ['TR3', 'Technical Report Type 3'],
  ['X12', 'A standards development organization that develops EDI standards and related documents for national and global markets'],
  ['HCPCS', 'Healthcare Common Procedure Coding System'],
  ['NPI', 'National Provider Identifier'],
  ['PTAN', 'Provider Transaction Access Number'],
  ['CCN', 'Claim Control Number'],
  ['COB', 'Coordination of Benefits'],
  ['DME', 'Durable Medical Equipment'],
  ['FFS', 'Medicare Fee-For-Service'],
  ['MAC', 'Medicare Administrative Contractor'],
  ['FISS', 'Fiscal Intermediary Standard System'],
  ['HIPAA', 'Health Insurance Portability and Accountability Act of 1996'],
  ['PHI', 'Protected Health Information'],
  ['CAQH CORE', 'Council for Affordable Quality Healthcare - Committee on Operating Rules for Information Exchange'],
]

const CELL = 'py-3 pr-5 align-top'

export function Reference835() {
  return (
    <div className="mx-auto max-w-[90rem] px-10 py-12">
      <PageTitle lede="What an 835 is, which elements carry its money and its identifiers, and where its codes come from — with the public document behind each name.">
        835 Reference
      </PageTitle>

      <DocSection id="what-it-is" title="What an 835 is">
        <p>
          An <strong>835 is the payer’s explanation of how claims were settled</strong>, sent to the
          provider. Its formal name is Health Care Claim Payment/Advice (CEDI 835 CG, PDF p.7, printed p.1);
          the same guide’s acronym list expands “835” as the Electronic Remittance Advice transaction, ERA
          (PDF p.31, printed p.25). It tells the provider which claims were processed, what was paid, and
          why billed amounts were adjusted. One transaction can explain several claims.{' '}
          <Source href={CMS_OVERVIEW}>CMS: payment and remittance advice</Source>.
        </p>
        <p>
          The payment itself can arrive separately, through a bank transfer or a check. The <code>TRN</code>{' '}
          segment, the {TRN_NAME['1']} ({cms('2')}), carries the {TRN_EL_NAME[1]} and {TRN_EL_NAME[2]} that
          tie the advice to the money. <Source href={X12_EXAMPLES}>X12: worked 835 examples</Source>.
        </p>
        <p>
          Every part of a file, each element under its key, is on the <Link to="/835-viewer" className={LINK}>835 viewer</Link>:
          28 small complete 835s, one per thing an 835 can carry. Every element with its public name and the page
          that names it is in that viewer's <Link to="/835-viewer?lookup=1" className={LINK}>lookup panel</Link>.
        </p>
      </DocSection>

      <DocSection id="money" title="The money">
        <p>The amounts an 835 carries, by element, with the flat file’s name for each.</p>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[36rem] text-left text-sm">
            <thead><tr className="border-b border-line text-xs text-ink"><th className={CELL}>Element</th><th className={CELL}>Name</th><th className="py-3">Source</th></tr></thead>
            <tbody>{AMOUNTS.map(([element, name, source]) => (
              <tr key={element} className="border-b border-line"><td className={`${CELL} font-mono text-xs text-ink`}>{element}</td><td className={`${CELL} font-medium text-ink`}>{name}</td><td className="py-3 text-xs text-ink-2">{source}</td></tr>
            ))}</tbody>
          </table>
        </div>
        <p>
          Every example on the 835 viewer is checked to balance before it is shown (<code>frontend/src/docs/viewer835/check.ts</code>): on a service line, SVC02 − SVC03 equals that line’s adjustment amounts; on a claim, CLP03 − CLP04 equals
          its claim-level and line-level adjustment amounts, and CLP03 equals the sum of its lines’ SVC02; on the transaction, BPR02 equals
          the sum of CLP04 minus the PLB amounts. See <Example id="deductible-and-coinsurance">Deductible and coinsurance</Example>,{' '}
          <Example id="fee-schedule-reduction">Fee-schedule reduction</Example>, <Example id="overpayment-recovered">Overpayment recovered</Example>,
          and <Example id="multiple-claims-one-payment">Multiple claims, one payment</Example>.
        </p>
      </DocSection>

      <DocSection id="identifiers" title="There is more than one ID">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[40rem] text-left text-sm">
            <thead><tr className="border-b border-line text-xs text-ink"><th className={CELL}>Element</th><th className={CELL}>Name</th><th className="py-3">Source</th></tr></thead>
            <tbody>{IDENTIFIERS.map(([element, name, source]) => (
              <tr key={element} className="border-b border-line"><td className={`${CELL} font-mono text-xs text-ink`}>{element}</td><td className={`${CELL} font-medium text-ink`}>{name}</td><td className="py-3 text-xs text-ink-2">{source}</td></tr>
            ))}</tbody>
          </table>
        </div>
        <p>
          The three control numbers pair a header with its trailer; the trailers also count what they close
          ({SE_EL_NAME[0]}, {GE_EL_NAME[0]}, {IEA_EL_NAME[0]}; {cms('14–15')}). Which of these the validator uses to
          pair claims is its own rule, described in the app’s help.
        </p>
      </DocSection>

      <DocSection id="codes" title="Statuses, reasons, remarks, and qualifiers">
        <dl className="grid gap-6 md:grid-cols-2">
          <div>
            <dt className="font-semibold text-ink">CLP02 · {CLP_EL_NAME[1]}</dt>
            <dd>
              What happened to the claim ({cms('5')}). The CMS flat file and the CEDI guide list its values without naming them
              in text; see <Example id="processed-as-primary">Processed as primary</Example>, <Example id="denied">Denied</Example>, and{' '}
              <Example id="reversal-and-correction">Reversal and correction</Example>.
            </dd>
          </div>
          <div>
            <dt className="font-semibold text-ink">CAS01 · {CAS01_EL_NAME}</dt>
            <dd>
              Who the adjustment falls to ({cms('6, 12')}): CO {CAS_GROUP_NAME.CO}, PR {CAS_GROUP_NAME.PR}, OA {CAS_GROUP_NAME.OA}{' '}
              (<Source href={X12_GROUP_CODES}>X12 Claim Adjustment Group Codes</Source>).
            </dd>
          </div>
          <div>
            <dt className="font-semibold text-ink">CAS02, CAS05, … · {CAS_REASON_EL_NAME}</dt>
            <dd>
              Why the amount changed ({cms('6, 12–13')}). The values are X12’s{' '}
              <Source href={X12_CARC}>Claim Adjustment Reason Codes</Source> (CARC); the validator app shows each code’s text on hover.
            </dd>
          </div>
          <div>
            <dt className="font-semibold text-ink">Remark codes</dt>
            <dd>
              {MOA_EL_NAME[2]} in MOA03–07 ({cms('10')}) and MIA05, MIA20–23 ({cms('9')}); {LQ_EL_NAME[1]} in LQ02 on a service line ({cms('13')}).
              The values are X12’s <Source href={X12_RARC}>Remittance Advice Remark Codes</Source> (RARC).
            </dd>
          </div>
        </dl>
        <p>
          A <strong>qualifier</strong> says what the value beside it is. <code>DTM*472</code> is the {DTM01_EL_NAME} 472, {DTM_NAME['472']} ({cms('12')});{' '}
          <code>REF*6R</code> carries the {SERVICE_REF_NAME['6R']} ({cms('13')}); <code>SVC*HC:…</code> names its procedure under the{' '}
          {SVC_QUALIFIER_NAME.HC} (CEDI 835 CG, PDF p.31, printed p.25). The same two letters mean different things in different
          elements: PR in N101 is {N1_NAME.PR} ({cms('2')}); PR in CAS01 is {CAS_GROUP_NAME.PR}.
        </p>
      </DocSection>

      <DocSection id="terms" title="Acronyms">
        <p>As the CEDI 835 companion guide’s acronym list expands them (Table 26, PDF p.31–32, printed p.25–26).</p>
        <dl className="grid gap-x-8 gap-y-3 md:grid-cols-2">
          {CEDI_ACRONYMS.map(([term, meaning]) => (
            <div key={term} className="grid grid-cols-[7rem_1fr] gap-3"><dt className="font-mono text-xs font-semibold text-ink">{term}</dt><dd className="text-sm leading-6">{meaning}</dd></div>
          ))}
          <div className="grid grid-cols-[7rem_1fr] gap-3"><dt className="font-mono text-xs font-semibold text-ink">CARC</dt><dd className="text-sm leading-6">Claim Adjustment Reason Code (<Source href={X12_CARC}>X12 code list</Source>)</dd></div>
          <div className="grid grid-cols-[7rem_1fr] gap-3"><dt className="font-mono text-xs font-semibold text-ink">RARC</dt><dd className="text-sm leading-6">Remittance Advice Remark Code (<Source href={X12_RARC}>X12 code list</Source>)</dd></div>
        </dl>
      </DocSection>

      <DocSection id="sources" title="Scope and sources">
        <p>
          This page covers the 005010X221A1 835 as the CMS flat file lays it out for Medicare. It names elements and codes from
          the documents below; the implementation guide (TR3) and each contractor’s companion guide determine exact required and
          situational usage and allowed repetitions.
        </p>
        <p>The CMS flat file holds both X12 element names and an internal fixed-width layout. Its flat-file offsets and COBOL storage formats are not the offsets or widths of fields in an EDI file.</p>
        <ul className="list-disc space-y-2 pl-5 text-sm">
          <li><Source href={CMS_OVERVIEW}>CMS: Health Care Payment and Remittance Advice</Source> — purpose and the claim, line, and provider levels.</li>
          <li><Source href={CMS_FIELDS}>CMS: 835 Flat File, 5010A1, updated January 10, 2023</Source> — X12 element names, loops, and Medicare usage notes; the “PDF p.N” citations.</li>
          <li><Source href={CEDI_835_GUIDE}>CEDI: 835 Standard Companion Guide</Source> — the DME MAC gateway’s envelope values and acronym list; the “CEDI 835 CG” citations (a copy is in <code>docs/reference/file-formats/</code>).</li>
          <li><Source href={CMS_MANUAL}>CMS: Medicare Claims Processing Manual, chapter 22</Source> — remittance behavior and adjustment codes.</li>
          <li><Source href={CMS_MANUAL_2003}>CMS: the same chapter as posted in 2003</Source> — its Medicare 835 companion documents, element by element for intermediaries and carriers, and the paper-remittance crosswalks; the “IOM 100-04 ch.22 (2003 ed.)” citations in the lookup (a copy is in <code>docs/reference/manuals/</code>).</li>
          <li><Source href={CMS_TI_GUIDE}>CMS: 835 TI Companion Guide, July 2012</Source> — the elements whose values Medicare restricts, by line of business; the “CMS 835 TI CG (2012)” citations.</li>
          <li><Source href={CMS_MLN_FAQ}>CMS: Remittance Advice Resources and FAQs (MLN)</Source> — the three code sets, provider-level adjustments, and balancing in plain language.</li>
          <li><Source href={CMS_MREP_GUIDE}>CMS: Medicare Remit Easy Print user guide</Source> — the free viewer’s field crosswalk to the 835, which names the claim status values.</li>
          <li><Source href={X12_EXAMPLES}>X12: 005010X221 worked examples</Source> — complete example transactions and business scenarios.</li>
          <li><Source href={X12_GROUP_CODES}>X12: Claim Adjustment Group Codes</Source> · <Source href={X12_CARC}>Claim Adjustment Reason Codes</Source> · <Source href={X12_RARC}>Remittance Advice Remark Codes</Source>.</li>
        </ul>
        <p>To see field differences on a compared claim, open the <a href={`${APP_URL}/claim`} className={LINK}>claim page in the validator app</a>.</p>
      </DocSection>
    </div>
  )
}
