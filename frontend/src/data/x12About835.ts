// What each element of the Medicare 835 is, what it is for, the values it
// takes, and how Medicare fills it: the entries behind the 835 viewer's
// lookup panel (frontend/src/docs/viewer/Lookup.tsx, rows from
// lookupRows835.ts) and the validator app's hover text (codeHints.ts), which
// reads the same entries and shows their sentences without the citations.
//
// One entry per use of an element. The key is `${scope}:${ref}`, with
// `@${qualifier}` when the segment's own qualifier picks the use
// ('claim:DTM02@050', 'claim:NM109@TT', 'payer:PER04@CX'); the scopes are
// envelope, header, payer, payee, summary, claim, line, provider, and
// trailer, and lookupRows835.ts builds the same keys from its rows. A
// segment-level entry ('claim:CLP') says what the segment as a whole is.
//
// Every sentence here stands on a cited public document, named in `sources`
// the way the rest of this folder cites them; where no document on file
// says a thing, the entry says so instead of guessing. The X12 TR3 is not
// in the repository and nothing here is taken from it. Page numbers are PDF
// viewer pages, counted from 1; scripts/ssr-check-lookup.mjs checks that
// every citation names a document under docs/reference/ (or a listed public
// web page) and a page inside it.

export interface About {
  // What it is: the element's own name and what the value stands for.
  is: string
  // What it is for and why the transaction carries it.
  purpose?: string
  // The values it takes: format and precision, or the code list with the
  // meanings the documents give.
  values?: string
  // How Medicare fills it: which lines of business send it and when.
  medicare?: string
  // Where the sentences above come from.
  sources: string[]
}

// ---------------------------------------------------------------- the public documents

export const CMS_FIELDS = 'https://www.cms.gov/medicare/billing/electronicbillingeditrans/downloads/835-flatfile.pdf'
export const CMS_OVERVIEW = 'https://www.cms.gov/medicare/coding-billing/electronic-billing/health-care-payment-remittance-advice'
export const CMS_MANUAL = 'https://www.cms.gov/Regulations-and-Guidance/Guidance/Manuals/downloads/clm104c22.pdf'
export const CMS_MANUAL_2003 = 'https://www.cms.gov/regulations-and-guidance/guidance/manuals/downloads/clm104c22_remit_phase3_apr24-03_r2.pdf'
export const CMS_TI_GUIDE = 'https://www.cms.gov/Medicare/Billing/ElectronicBillingEDITrans/Downloads/5010A1835CG.pdf'
export const CMS_MLN_FAQ = 'https://www.cms.gov/Outreach-and-Education/Medicare-Learning-Network-MLN/MLNProducts/Downloads/ICN905367TextOnly.pdf'
export const CMS_MREP_GUIDE = 'https://www.cms.gov/Research-Statistics-Data-and-Systems/CMS-Information-Technology/AccesstoDataApplication/downloads/EasyPrintUserGuide.pdf'
export const CEDI_835_GUIDE = 'https://www.ngscedi.com/technical-specifications'
export const X12_EXAMPLES = 'https://x12.org/examples/005010x221'
export const X12_GROUP_CODES = 'https://x12.org/codes/claim-adjustment-group-codes'
export const X12_CARC = 'https://x12.org/codes/claim-adjustment-reason-codes'
export const X12_RARC = 'https://x12.org/codes/remittance-advice-remark-codes'
export const X12_PLB_REASONS = 'https://x12.org/codes/provider-adjustment-reason-codes'

// Citation forms. Each names one document under docs/reference/ (or one
// public web page) and a PDF viewer page; the check script knows the list.
const FF = (p: string) => `CMS 835 Flat File, PDF p.${p}`
const CEDI = (p: number, printed: number) => `CEDI 835 CG, PDF p.${p} (printed p.${printed})`
const TI = (p: number) => `CMS 835 TI CG (2012), PDF p.${p}`
const IOM = (p: number) => `IOM 100-04 ch.22 (Rev. 11427), PDF p.${p}`
const IOM03 = (p: number) => `IOM 100-04 ch.22 (2003 ed.), PDF p.${p}`
const MLN = (p: number) => `MLN Remittance Advice Resources and FAQs (ICN 905367), PDF p.${p}`
const MREP = (p: number) => `MREP User Guide v4.4 (2016), PDF p.${p} (printed p.${p - 2}), ${CMS_MREP_GUIDE}`
const T907 = (p: number) => `Transmittal 907 (2011), PDF p.${p}`
const T579 = (p: number) => `Transmittal 579 (2009), PDF p.${p}`
const CEDI837P = (p: number, printed: number) => `CEDI 837P CG, PDF p.${p} (printed p.${printed})`
const FORMATS_README = 'docs/reference/file-formats/README.md'
const API_DOC = 'docs/API.md'
const X12_GROUPS = `X12 Claim Adjustment Group Codes, ${X12_GROUP_CODES}`
const X12_CARCS = `X12 Claim Adjustment Reason Codes, ${X12_CARC}`
const X12_RARCS = `X12 Remittance Advice Remark Codes, ${X12_RARC}`
const X12_PLB = `X12 Provider Adjustment Reason Codes, ${X12_PLB_REASONS}`

// Sentences reused across entries.
const AMOUNT_FORMAT =
  'An X12 decimal (type R): digits with an optional decimal point and a leading minus sign, no plus sign, no trailing zeros after the point and no thousands separators. The CMS flat file holds it as a signed number with two decimal places (S9(8)V99).'
const DATE_FORMAT = 'A date, CCYYMMDD (type DT, 8 digits).'
const COUNT_FORMAT = 'A count, an X12 decimal (type R) held in the flat file as 15 characters.'
const SILENT = 'The documents on file list the values but do not spell out what each one means.'
const NOT_USED_MEDICARE = 'The CMS 835 flat file marks it not used by Medicare (N/U), so Medicare never sends it.'

const ENVELOPE_SOURCES = [FF('1'), CEDI(20, 14), CEDI(21, 15), FORMATS_README]

const entries: Record<string, About> = {}

function put(keys: string | string[], about: About) {
  const list = typeof keys === 'string' ? [keys] : keys
  for (const key of list) entries[key] = about
}

// ---------------------------------------------------------------- envelope

put('envelope:ISA', {
  is: 'Interchange Control Header: the first segment of an X12 file. It opens the interchange, the envelope that holds every functional group and transaction in the file, and its sixteen elements have fixed widths.',
  purpose:
    'It says who sent the file and who receives it, when it was built, which control number pairs it with the IEA trailer, and which characters the file uses as its repetition and component separators, so a receiver can read everything that follows.',
  medicare: 'The CMS flat file calls every ISA element translator generated: the translator that turns the flat file into X12 fills them, not the claims system. The DME gateway lists what it sends in each one.',
  sources: ENVELOPE_SOURCES,
})
put('envelope:ISA01', {
  is: 'Authorization Information Qualifier: says what kind of authorization information, if any, ISA02 carries.',
  values: '00. Medicare sends 00, with ten blank spaces in ISA02.',
  sources: [FF('1'), CEDI(20, 14)],
})
put('envelope:ISA02', {
  is: 'Authorization Information: a ten-character field for authorization information, blank when ISA01 is 00.',
  values: 'Ten spaces. Medicare sends ten blank spaces.',
  sources: [FF('1'), CEDI(20, 14)],
})
put('envelope:ISA03', {
  is: 'Security Information Qualifier: says what kind of security information, if any, ISA04 carries.',
  values: '00. Medicare sends 00, with spaces in ISA04.',
  sources: [FF('1'), CEDI(20, 14)],
})
put('envelope:ISA04', {
  is: 'Security Information: a ten-character field for security information, blank when ISA03 is 00.',
  values: 'Ten spaces. Medicare sends spaces.',
  sources: [FF('1'), CEDI(20, 14)],
})
put('envelope:ISA05', {
  is: 'Interchange ID Qualifier for the sender: says what kind of identifier ISA06 is.',
  values:
    'The flat file lists 01, 14, 20, 27, 28, 29, 30, 33, and ZZ. The DME gateway accepts 27, 28, and ZZ and sends 27. ZZ means a mutually defined identifier; the 2003 companion document had intermediaries use ZZ because trading partners always agree the sender identifier between themselves. The other codes are listed without their meanings.',
  sources: [FF('1'), CEDI(20, 14), IOM03(10)],
})
put('envelope:ISA06', {
  is: 'Interchange Sender ID: the identifier of the party that sent the interchange, fifteen characters, padded with spaces.',
  medicare:
    'For DME remittances it is the DME MAC contractor ID: 16013 for jurisdiction A, 17013 for B, 18003 for C, and 19003 for D. The same number appears in GS02 and in the payer loop REF*2U.',
  sources: [FF('1'), CEDI(20, 14), CEDI(21, 15), CEDI(24, 18)],
})
put('envelope:ISA07', {
  is: 'Interchange ID Qualifier for the receiver: says what kind of identifier ISA08 is.',
  values:
    'The flat file lists the same codes as ISA05. Medicare sends 29, the qualifier the 2003 companion document tied to the NPI, with ZZ (mutually defined) as the alternative when a locally defined qualifier is used.',
  sources: [FF("1"), CEDI(20, 14), IOM03(10)],
})
put('envelope:ISA08', {
  is: 'Interchange Receiver ID: the identifier of the party that receives the interchange, fifteen characters, padded with spaces.',
  medicare: 'For DME remittances it is the submitter number the CEDI gateway assigned to the trading partner. The header REF*EV, when present, must match it.',
  sources: [FF('1'), CEDI(20, 14), IOM03(14)],
})
put('envelope:ISA09', {
  is: 'Interchange Date: the date the interchange was built, YYMMDD.',
  values: 'Six digits, YYMMDD. The flat file stores the same date with the century added and calls it translator generated.',
  sources: [FF('1')],
})
put('envelope:ISA10', {
  is: 'Interchange Time: the time the interchange was built, HHMM.',
  values: 'Four digits, HHMM, on a 24-hour clock. Translator generated.',
  sources: [FF('1')],
})
put('envelope:ISA11', {
  is: 'Repetition Separator: the character the file uses to separate repeats of an element. In version 5010 this position is a delimiter, not a data element.',
  values: 'One character. The DME gateway sends ^. The 2011 flat-file change log records the change from a value to a delimiter in 5010.',
  sources: [FF('1'), CEDI(21, 15), T907(281)],
})
put('envelope:ISA12', {
  is: 'Interchange Control Version Number: the version of the interchange control segments.',
  values: '00501 for version 5010.',
  sources: [FF('1')],
})
put('envelope:ISA13', {
  is: 'Interchange Control Number: a nine-digit number that identifies this interchange. The IEA trailer repeats it, so a receiver can pair the two.',
  values:
    'Nine digits. The flat file builds it as YJJJ00000: Y is the last digit of the cycle year, JJJ the Julian day of the cycle, and the last five a sequence starting at 00001. IEA02 carries the same value.',
  sources: [FF('1'), FF('15')],
})
put('envelope:ISA14', {
  is: 'Acknowledgement Requested.',
  values: '0. Medicare sends 0.',
  sources: [FF('1'), CEDI(21, 15)],
})
put('envelope:ISA15', {
  is: 'Usage Indicator.',
  values: 'P or T. ' + SILENT,
  sources: [FF('1')],
})
put('envelope:ISA16', {
  is: 'Component Element Separator: the character the file uses between the components of a composite element, such as the parts of SVC01.',
  values: 'One character, translator generated. The DME gateway publishes the delimiters it uses on request and promises not to use them inside data.',
  sources: [FF('1'), CEDI(21, 15)],
})

put('envelope:GS', {
  is: 'Functional Group Header: opens a functional group, the envelope inside the interchange that holds one or more transactions of the same kind.',
  purpose: 'It names the kind of transaction (HP, a payment advice), the application sender and receiver, the date and time, a control number the GE trailer repeats, and the version of the implementation guide the transactions follow.',
  sources: [FF('1'), FF('2'), CEDI(21, 15), FORMATS_README],
})
put('envelope:GS01', {
  is: 'Functional Identifier Code: the kind of transaction in the group.',
  values: 'HP, Health Care Claim Payment/Advice.',
  sources: [FF('1'), CEDI(7, 1)],
})
put('envelope:GS02', {
  is: "Application Sender's Code: the sending application's identifier.",
  medicare: 'For DME remittances, the DME MAC contractor ID (16013, 17013, 18003, 19003 for jurisdictions A to D).',
  sources: [FF('1'), CEDI(21, 15)],
})
put('envelope:GS03', {
  is: "Application Receiver's Code: the receiving application's identifier.",
  medicare: 'For DME remittances, the submitter number the CEDI gateway assigned. The 2003 companion document let it be the provider number, chain ID, network vendor ID, or EDI submitter number, as the trading partner agreement sets.',
  sources: [FF('1'), CEDI(21, 15), IOM03(12)],
})
put('envelope:GS04', {
  is: 'Date: the date the functional group was built.',
  values: 'CCYYMMDD, eight digits; the flat file says see ISA09.',
  sources: [FF('2')],
})
put('envelope:GS05', {
  is: 'Time: the time the functional group was built.',
  values: 'HHMM, optionally with seconds; the flat file says see ISA10.',
  sources: [FF('2')],
})
put('envelope:GS06', {
  is: 'Group Control Number: identifies this functional group; the GE trailer repeats it.',
  values: 'Up to nine digits. The flat file increments it by one, beginning at 001; GE02 carries the same value.',
  sources: [FF('2'), FF('14')],
})
put('envelope:GS07', {
  is: 'Responsible Agency Code: the agency responsible for the standard.',
  values: 'X. The CEDI acronym list names ASC the Accredited Standards Committee and X12 the standards body.',
  sources: [FF('2'), CEDI(32, 26)],
})
put('envelope:GS08', {
  is: 'Version, Release, and Industry Identifier Code: the implementation guide the transactions follow.',
  values: '005010X221A1; the 2011 flat-file change log records the A1 suffix as the errata release.',
  sources: [FF('2'), CEDI(21, 15)],
})

put('envelope:ST', {
  is: 'Transaction Set Header: opens one 835 transaction, one remittance with its one BPR payment. The SE trailer closes it.',
  sources: [FF('1'), CEDI(23, 17), FORMATS_README],
})
put('envelope:ST01', {
  is: 'Transaction Set Identifier Code: the kind of transaction.',
  values: '835, Health Care Claim Payment/Advice.',
  sources: [FF('1'), CEDI(7, 1)],
})
put('envelope:ST02', {
  is: 'Transaction Set Control Number: identifies this transaction inside its group; SE02 repeats it.',
  values: 'Four to nine characters; the DME gateway counts up from 0001. SE02 carries the same value.',
  sources: [FF('1'), CEDI(23, 17), FF('14')],
})

// ---------------------------------------------------------------- header

put('header:BPR', {
  is: 'Financial Information: the one segment that describes the payment itself: how it is made, how much, from which bank account to which, and on what date.',
  purpose:
    'The remittance is the companion to a payment, and this segment is the payment. The check or EFT amount here is the sum of every claim-level payment in the transaction less the provider-level adjustments, which is how a transaction balances.',
  medicare: 'Medicare issues one check or EFT when payment is due, covering all benefits due for the claims itemized in the remittance. MACs may release the remittance before the payment date but never after it.',
  sources: [FF('1'), IOM(2), IOM(3), IOM(5), MLN(2), MLN(8)],
})
put('header:BPR01', {
  is: 'Transaction Handling Code: how the payment and the remittance data travel: together, separately, or as information only.',
  values:
    'The flat file lists C, D, H, I, P, U, and X. The 2003 companion documents say U and X do not apply to Medicare. The documents on file do not spell the codes out; X12\'s public example 1, "Dollars and Data Sent Separately", pairs I with an ACH payment.',
  sources: [FF('1'), IOM03(12), IOM03(34), X12_EXAMPLES],
})
put('header:BPR02', {
  is: 'Total Actual Provider Payment Amount: the amount of the check or EFT this remittance accompanies.',
  purpose:
    'The check amount is the sum of all claim-level payments (CLP04), including claim- and service-level adjustments, less any provider-level adjustments (PLB). That equation is remittance balancing at the transaction level, and every Medicare 835 must balance.',
  values: AMOUNT_FORMAT + ' Zero on a notice of no payment.',
  medicare: 'The validator keeps it on the payment row for a one-side balance check of each transaction and never compares it between a MAC file and a vendor file.',
  sources: [FF('1'), IOM(3), IOM(4), MLN(8), MLN(11), API_DOC],
})
put('header:BPR03', {
  is: 'Credit or Debit Flag Code: whether the payment credits or debits the receiver.',
  values: 'C, Credit. The DME gateway and the 2012 companion guide say D (debit) does not apply to Medicare.',
  sources: [FF('1'), CEDI(24, 18), TI(6)],
})
put('header:BPR04', {
  is: 'Payment Method Code: how the money moves.',
  values:
    'ACH, CHK, or NON. ACH is payment by electronic funds transfer through the automated clearing house; when BPR04 is ACH the bank fields BPR05 to BPR15 are required. CHK and NON are listed without text; the Easy Print guide reads the trace number as the check or EFT number for any of the three. BOP and FWT do not apply to Medicare.',
  sources: [FF('1'), CEDI(24, 18), TI(6), IOM(4), IOM03(34), MREP(214)],
})
put('header:BPR05', {
  is: 'Payment Format Code: the format of the electronic payment.',
  values: 'CCP or CTX. Required when BPR04 is ACH. ' + SILENT,
  sources: [FF('1'), IOM03(34)],
})
put('header:BPR06', {
  is: "Depository Financial Institution (DFI) Identification Number Qualifier for the payer's bank: what kind of number BPR07 is.",
  values: '01. Required when BPR04 is ACH; 04 does not apply to Medicare.',
  sources: [FF('1'), CEDI(24, 18), TI(6), IOM03(34)],
})
put('header:BPR07', {
  is: "Sender DFI Identifier: the payer's bank, identified as BPR06 says.",
  values: 'Three to twelve characters. Required when BPR04 is ACH.',
  sources: [FF('1'), IOM03(34)],
})
put('header:BPR08', {
  is: "Account Number Qualifier for the payer's account: what kind of account BPR09 names.",
  values: 'DA. Required when BPR04 is ACH. ' + SILENT,
  sources: [FF('1'), IOM03(34)],
})
put('header:BPR09', {
  is: "Sender Bank Account Number: the payer's account the money comes from.",
  values: 'Up to 35 characters. Required when BPR04 is ACH.',
  sources: [FF('1'), IOM03(34)],
})
put('header:BPR10', {
  is: 'Payer Identifier: the ten-character identifier of the payer as originator of the payment.',
  values: 'Ten characters. TRN03 must carry the same value.',
  sources: [FF('1'), IOM03(13)],
})
put('header:BPR11', {
  is: 'Originating Company Supplemental Code: a further code for the originating company.',
  values: 'Nine characters, equal to TRN04 when sent.',
  medicare: 'Reported by Part A (institutional) only; the flat file says not used by Part B, and the 2003 carrier document says it should not be reported.',
  sources: [FF('1'), TI(6), IOM03(34)],
})
put('header:BPR12', {
  is: "DFI Identification Number Qualifier for the receiver's bank: what kind of number BPR13 is.",
  values: '01. Required when BPR04 is ACH; 04 does not apply to Medicare.',
  sources: [FF('2'), CEDI(24, 18), TI(6), IOM03(34)],
})
put('header:BPR13', {
  is: "Receiver or Provider Bank ID Number: the payee's bank, identified as BPR12 says.",
  values: 'Three to twelve characters. Required when BPR04 is ACH.',
  sources: [FF('2'), IOM03(34)],
})
put('header:BPR14', {
  is: "Account Number Qualifier for the payee's account: what kind of account BPR15 names.",
  values: 'DA or SG. Required when BPR04 is ACH. ' + SILENT,
  sources: [FF('2'), IOM03(34)],
})
put('header:BPR15', {
  is: "Receiver or Provider Account Number: the payee's account the money goes to.",
  values: 'Up to 35 characters. Required when BPR04 is ACH.',
  sources: [FF('2'), IOM03(34)],
})
put('header:BPR16', {
  is: 'Check Issue or EFT Effective Date: the date of the payment.',
  purpose: 'The paper remittance prints it as the check or EFT issue date. MACs may release the remittance before this date but never after it.',
  values: DATE_FORMAT,
  medicare: 'The validator keeps it on the payment row as the payment date.',
  sources: [FF('2'), IOM(5), IOM03(66), MREP(214), API_DOC],
})

put('header:TRN@1', {
  is: 'Reassociation Trace Number: the check or EFT number, with the payer that originated it, so the money and the advice can be matched when they travel separately.',
  medicare: 'When no payment is issued, intermediaries put the remittance advice number here.',
  sources: [FF('2'), IOM03(13), MLN(2)],
})
put('header:TRN01@1', {
  is: 'Trace Type Code: the kind of trace.',
  values: '1, the only value Medicare sends; the segment it opens is the Reassociation Trace Number.',
  sources: [FF('2')],
})
put('header:TRN02@1', {
  is: 'Check or EFT Trace Number: the number of the check or the electronic transfer this remittance accompanies.',
  values: 'Up to 50 characters. The paper remittance prints it as the check or EFT trace number; when no payment is issued, the remittance advice number stands in.',
  sources: [FF('2'), IOM03(13), IOM03(66), MREP(214)],
})
put('header:TRN03@1', {
  is: 'Payer Identifier: the payer that originated the trace number.',
  values: 'Ten characters; must equal BPR10.',
  sources: [FF('2'), IOM03(13)],
})
put('header:TRN04@1', {
  is: 'Originating Company Supplemental Code: a further code for the originating company.',
  values: 'Up to 50 characters, equal to BPR11.',
  medicare: 'Reported by Part A only; it does not apply to carriers.',
  sources: [FF('2'), TI(6), IOM03(34)],
})

put(['header:CUR', 'header:CUR01', 'header:CUR02'], {
  is: 'Foreign Currency Information: the currency of the payment when it is not the default.',
  medicare: 'Medicare does not use it: the flat file and the 2012 companion guide both say the segment does not apply.',
  sources: [FF('2'), TI(6)],
})

put('header:REF@EV', {
  is: 'Receiver Identification: identifies the receiver of the transaction when it is someone other than the payee, such as a clearinghouse or billing service.',
  medicare: 'Required when the 835 goes to any entity other than the payee; the value must correspond to ISA08.',
  sources: [FF('2'), IOM03(13), IOM03(34)],
})
put('header:REF01@EV', {
  is: 'Reference Identification Qualifier: the kind of reference.',
  values: 'EV; the flat file names the segment Receiver Identification.',
  sources: [FF('2')],
})
put('header:REF02@EV', {
  is: "Receiver Identifier: the receiver's identifier, matching ISA08.",
  values: 'Up to 50 characters.',
  sources: [FF('2'), IOM03(14)],
})
put('header:REF@F2', {
  is: "Version Identification: a local version number for the payer's implementation of the 835.",
  purpose: 'The 2003 carrier document explains it: a local version number marks a programming change made after implementation, such as a fix, so a provider inquiry about the change can be answered.',
  medicare: 'Carriers send it; it does not apply to intermediaries.',
  sources: [FF('2'), IOM03(14), IOM03(35)],
})
put('header:REF01@F2', {
  is: 'Reference Identification Qualifier: the kind of reference.',
  values: 'F2; the flat file names the segment Version Identification.',
  sources: [FF('2')],
})
put('header:REF02@F2', {
  is: 'Version ID Code: the locally assigned version number.',
  values: 'Up to 50 characters.',
  sources: [FF('2'), IOM03(35)],
})

put('header:DTM@405', {
  is: 'Production Date: the date the remittance data was produced.',
  medicare: 'Sent when the date of the 835 differs from the cutoff date of the adjudication run that produced it.',
  sources: [FF('2'), IOM03(14), IOM03(35)],
})
put('header:DTM01@405', {
  is: 'Date/Time Qualifier: the kind of date.',
  values: '405, Production.',
  sources: [FF('2')],
})
put('header:DTM02@405', {
  is: 'Production Date.',
  values: DATE_FORMAT,
  medicare: 'The validator keeps it on the payment row as the production date.',
  sources: [FF('2'), API_DOC],
})

// ---------------------------------------------------------------- payer loop 1000A

put('payer:N1@PR', {
  is: 'Payer Identification: the loop that names the payer, the Medicare contractor that adjudicated the claims and makes the payment, with its address and contacts.',
  sources: [FF('2'), IOM03(14), IOM03(35)],
})
put('payer:N101@PR', {
  is: 'Entity Identifier Code: the role of the party named.',
  values: 'PR, Payer.',
  sources: [FF('2')],
})
put('payer:N102@PR', {
  is: 'Payer Name: the name of the paying contractor.',
  values: 'Up to 60 characters; the flat file widened every name to 60 for HIGLAS, the accounting system.',
  sources: [FF('2')],
})
put('payer:N103@PR', {
  is: 'Identification Code Qualifier: the kind of identifier N104 is.',
  values: 'XV, the CMS PlanID (the 837 names it the Health Care Financing Administration National PlanID). The 2003 documents sent XV once the PlanID took effect.',
  sources: [FF('2'), T907(9), IOM03(14), IOM03(35)],
})
put('payer:N104@PR', {
  is: "Payer Identifier: the payer's identifier of the kind N103 says.",
  values: 'Two to 80 characters.',
  medicare:
    'The validator reads the trimmed value as the MAC workload of the match key: claims pair only within one workload, and the Workloads page groups by it.',
  sources: [FF('2'), API_DOC],
})
put('payer:N3', {
  is: "Payer Address: the payer's street address, one or two lines.",
  sources: [FF('2'), IOM03(14)],
})
put('payer:N301', {
  is: 'Payer Address Line: the first line of the street address.',
  values: 'Up to 55 characters.',
  sources: [FF('2')],
})
put('payer:N302', {
  is: 'Payer Address Line: the second line, when the address has one (a suite number, for example).',
  values: 'Up to 55 characters.',
  sources: [FF('2'), IOM03(14)],
})
put('payer:N4', {
  is: "Payer City, State, Zip: the rest of the payer's address.",
  sources: [FF('2')],
})
put('payer:N401', { is: 'Payer City Name.', values: 'Two to 30 characters.', sources: [FF('2')] })
put('payer:N402', { is: 'Payer State Code: the two-letter state or province code.', values: 'Two characters.', sources: [FF('2')] })
put('payer:N403', { is: 'Payer Postal Zone or ZIP Code.', values: 'Three to 15 characters.', sources: [FF('2')] })

put('payer:REF@2U', {
  is: 'Additional Payer Identification: a further identifier for the payer.',
  medicare:
    'For DME remittances it carries the DME MAC contractor ID (16013, 17013, 18003, or 19003). The 2003 documents required it before the PlanID took effect and allowed it alongside the PlanID afterwards. EO, HI, and NF are not used.',
  sources: [FF('3'), CEDI(24, 18), TI(6), IOM03(15), IOM03(36)],
})
put('payer:REF01@2U', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values: '2U; the flat file names the segment Additional Payer Identification.',
  sources: [FF('3'), CEDI(24, 18)],
})
put('payer:REF02@2U', {
  is: 'Additional Payer ID: the identifier itself.',
  values: 'Up to 50 characters. For DME, the contractor ID.',
  medicare: 'The project\'s MAC workload lookup notes that the DME guide prints the contractor numbers here, in ISA06, and in GS02.',
  sources: [FF('3'), CEDI(24, 18), 'backend/pkg/macs/macs.go'],
})

const PER_PAYER_SOURCES = [FF('3'), IOM03(36), MREP(214)]
put('payer:PER@CX', {
  is: 'Payer Business Contact Information: whom to contact at the payer about the remittance.',
  medicare: 'The 2003 carrier document recommends it but leaves reporting contact information to each contractor; contact data comes from the contractor, not the claims system.',
  sources: PER_PAYER_SOURCES,
})
put('payer:PER@BL', {
  is: 'Payer Technical Contact Information: whom to contact at the payer about the transmission itself.',
  medicare: 'The flat file marks this one required. Contact data comes from the contractor, not the claims system.',
  sources: PER_PAYER_SOURCES,
})
put('payer:PER@1C', {
  is: "Payer Web Site: the payer's web address.",
  values:
    'The flat file lists the qualifier as IC, and the 2011 flat-file change log records the correction from 1C to IC "to reflect the value correctly per TR3"; the Easy Print guide reads the web site from PER04 when PER01 is IC. The worked examples on this site spell it 1C.',
  sources: [FF("3"), T907(281), MREP(214)],
})
const perContact = (scope: string, qualifier: string, who: string) => {
  put(`${scope}:PER01@${qualifier}`, {
    is: 'Contact Function Code: the kind of contact this segment names.',
    values: qualifier === '1C' ? 'IC in the flat file (see the segment\'s entry about 1C).' : `${qualifier}.`,
    sources: [FF(scope === 'payer' ? '3' : '11')],
  })
  put(`${scope}:PER02@${qualifier}`, {
    is: `${who} Contact Name: the name of the person or office to contact.`,
    values: qualifier === '1C' ? 'Not used in the web-site PER.' : 'Up to 60 characters.',
    sources: [FF(scope === 'payer' ? '3' : '11')],
  })
  const pairs: [number, string][] = qualifier === '1C' ? [[3, 'first']] : [[3, 'first'], [5, 'second'], [7, 'third']]
  for (const [n, kind] of pairs) {
    put(`${scope}:PER${String(n).padStart(2, '0')}@${qualifier}`, {
      is: `Communication Number Qualifier: what kind of number the ${kind} communication number (PER${String(n + 1).padStart(2, '0')}) is.`,
      values:
        'The flat file lists EM, EX, FX, TE, and UR by position. The Easy Print guide reads EM as an email address, TE as a telephone number, EX as its extension, and the web-site PER\'s UR as the web address. FX is listed without text.',
      sources: [FF(scope === 'payer' ? '3' : '11'), MREP(214)],
    })
    put(`${scope}:PER${String(n + 1).padStart(2, '0')}@${qualifier}`, {
      is: `${who} Contact Communication Number: the ${kind} number or address, of the kind PER${String(n).padStart(2, '0')} says.`,
      values: 'Up to 256 characters.',
      sources: [FF(scope === 'payer' ? '3' : '11')],
    })
  }
}
perContact('payer', 'CX', 'Payer')
perContact('payer', 'BL', 'Payer')
perContact('payer', '1C', 'Payer')
put('payer:PER03@1C', {
  is: 'Communication Number Qualifier: what kind of address PER04 is.',
  values: 'UR, a web address.',
  sources: [FF('3'), MREP(214)],
})
put('payer:PER04@1C', {
  is: "Payer Contact Communication Number: the payer's web address.",
  values: 'Up to 256 characters.',
  sources: [FF('3'), MREP(214)],
})

// ---------------------------------------------------------------- payee loop 1000B

put('payee:N1@PE', {
  is: 'Payee Identification: the loop that names the provider or supplier being paid, with its address and tax identifier.',
  purpose: 'The payee is identified by its NPI here; the manual says the provider-level identification on an 835 is one of the taxpayer identification number, the CMS PlanID, or the NPI.',
  sources: [FF('3'), IOM(2), IOM03(15), IOM03(36)],
})
put('payee:N101@PE', {
  is: 'Entity Identifier Code: the role of the party named.',
  values: 'PE, Payee.',
  sources: [FF('3')],
})
put('payee:N102@PE', {
  is: 'Payee Name: the name of the provider or supplier being paid.',
  values: 'Up to 60 characters (widened to 60 for HIGLAS).',
  sources: [FF('3')],
})
put('payee:N103@PE', {
  is: 'Identification Code Qualifier: the kind of identifier N104 is.',
  values:
    'XX, FI, or XV. XX is the National Provider Identifier; XV is the CMS PlanID. FI is listed without text; the 2003 carrier document had carriers send FI until the NPI took effect, and XX after.',
  sources: [FF('3'), CEDI(32, 26), T907(9), IOM03(36)],
})
put('payee:N104@PE', {
  is: "Payee ID Code: the payee's identifier of the kind N103 says, normally the NPI.",
  values: 'Two to 80 characters.',
  medicare: 'The validator takes the payee NPI of the match key from here.',
  sources: [FF('3'), MREP(214), API_DOC],
})
put('payee:N3', { is: "Payee Address: the payee's street address, one or two lines.", sources: [FF('3'), IOM03(36)] })
put('payee:N301', { is: 'Payee Address Line: the first line of the street address.', values: 'Up to 55 characters.', sources: [FF('3')] })
put('payee:N302', { is: 'Payee Address Line: the second line, when there is one.', values: 'Up to 55 characters.', sources: [FF('3'), IOM03(37)] })
put('payee:N4', { is: "Payee City, State, Zip: the rest of the payee's address.", sources: [FF('3'), FF('4')] })
put('payee:N401', { is: 'Payee City Name.', values: 'Two to 30 characters.', sources: [FF('3')] })
put('payee:N402', { is: 'Payee State Code.', values: 'Two characters.', sources: [FF('4')] })
put('payee:N403', { is: 'Payee Postal Zone or ZIP Code.', values: 'Three to 15 characters.', sources: [FF('4')] })

put('payee:REF@TJ', {
  is: "Payee Additional Identification: the payee's taxpayer identification number.",
  medicare: 'Required for Medicare to report the TIN once the NPI took effect; before that, carriers used 1C (the Medicare provider number) or 1G (the UPIN) here.',
  sources: [FF('4'), TI(6), IOM03(16), IOM03(37)],
})
put('payee:REF01@TJ', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values: 'TJ, the taxpayer identification number.',
  sources: [FF('4'), TI(6)],
})
put('payee:REF02@TJ', {
  is: 'Additional Payee ID Number: the TIN.',
  values: 'Up to 50 characters.',
  medicare: 'The validator takes the tax ID of the match key from here.',
  sources: [FF('4'), TI(6), API_DOC],
})

put('payee:RDM', {
  is: 'Remittance Delivery Method: how a remittance is delivered when the payment and the remittance travel separately.',
  sources: [FF('4')],
})
put('payee:RDM01', {
  is: 'Report Transmission Code: the delivery method.',
  values: 'BM, EM, FT, or OL. ' + SILENT,
  sources: [FF('4')],
})
put('payee:RDM02', { is: 'Name: the name the remittance is delivered to.', values: 'Up to 60 characters.', sources: [FF('4')] })
put('payee:RDM03', { is: 'Communication Number: the address or number it is delivered to.', values: 'Up to 256 characters.', sources: [FF('4')] })

// ---------------------------------------------------------------- loop 2000, the header number and provider summaries

put('summary:LX', {
  is: 'Header Number: opens a group of claims inside the transaction. The claim loops that follow belong to it until the next LX.',
  purpose: 'Medicare uses the number to separate assigned from non-assigned claims; the paper remittance prints the assignment flag from it.',
  sources: [FF("4"), CEDI(25, 19), TI(7), IOM03(67), MREP(214)],
})
put('summary:LX01', {
  is: 'Assigned Number: the number of this group of claims.',
  values:
    'Part B and DME send 1 for assigned claims and 0 for non-assigned claims (the shared systems write 000001 and 000000). Part A (FISS) sends TTYYMM: the facility type code, year, and month.',
  sources: [FF('4'), CEDI(25, 19), TI(7), MREP(214)],
})

put('summary:TS3', {
  is: 'Provider Summary Information: totals for one provider and type of bill over a fiscal period, at the head of a group of claims.',
  medicare: 'Part A only: intermediaries send it when applicable, carriers never. Of its 24 positions Medicare fills eleven; the flat file marks the rest not used.',
  sources: [FF('4'), FF('5'), IOM03(16), IOM03(37)],
})
const ts3 = (ref: string, is: string, medicare: string, page: string, values?: string) =>
  put(`summary:${ref}`, { is, values, medicare, sources: [FF(page), IOM03(ref <= "TS310" ? 16 : ref <= "TS320" ? 17 : 18), IOM03(37)] })
ts3('TS301', 'Provider Identifier: the provider the totals belong to.', 'The NPI; the flat file keeps the legacy number for VA pricing.', '4', 'Up to 50 characters.')
ts3('TS302', 'Facility Code Value: the type of bill the totals belong to.', 'The type of bill; the flat file says Part B will use either 11 or 99.', '4', 'Two characters.')
ts3('TS303', 'Fiscal Period Date: the end of the fiscal period the totals cover.', 'Sent by intermediaries.', '4', DATE_FORMAT)
ts3('TS304', 'Total Claim Count: the number of claims in the group.', 'Sent by intermediaries.', '4', COUNT_FORMAT)
ts3('TS305', 'Total Claim Charge Amount: the submitted charges of those claims added together.', 'Sent by intermediaries.', '4', AMOUNT_FORMAT)
for (const [ref, is] of [
  ['TS306', 'Total Covered Charge Amount.'],
  ['TS307', 'Total Noncovered Charge Amount.'],
  ['TS308', 'Total Denied Charge Amount.'],
  ['TS309', 'Total Provider Amount: the payments to the provider for the period, including interest.'],
  ['TS310', 'Total Interest Amount.'],
  ['TS311', 'Total Contractual Adjustment Amount.'],
  ['TS312', 'Total Gramm-Rudman Reduction Amount.'],
  ['TS314', 'Total Blood Deductible Amount.'],
  ['TS316', 'Total Coinsurance Amount.'],
  ['TS319', 'Total Deductible Amount.'],
] as const) {
  ts3(ref, is, 'The CMS 835 flat file marks it not used (N/U); the 2003 companion document had intermediaries send it when the amount applied.', ref <= 'TS316' ? '4' : '4')
}
ts3('TS313', 'Total MSP Payer Amount: payments by payers primary to Medicare for the claims in the period, including coinsurance and deductible amounts another payer paid.', 'Part A only.', '4', AMOUNT_FORMAT)
ts3('TS315', 'Total Non-Lab Charge Amount: the covered charges less the charges for lines paid on the clinical laboratory or orthotics and prosthetics fee schedules.', 'Part A only.', '4', AMOUNT_FORMAT)
ts3('TS317', 'Total HCPCS Reported Charge Amount: the reported charges for the lines coded with HCPCS.', 'Part A only.', '4', AMOUNT_FORMAT)
ts3('TS318', 'Total HCPCS Payable Amount: the benefits allowed for HCPCS lines paid on the clinical laboratory or orthotics and prosthetics fee schedules.', 'Part A only.', '4', AMOUNT_FORMAT)
ts3('TS320', 'Total Professional Component Amount: professional component amounts paid to the provider.', 'Part A only.', '4', AMOUNT_FORMAT)
ts3('TS321', 'Total MSP Patient Liability Met Amount: patient liability amounts (PR adjustments) that other payers satisfied.', 'Part A only.', '5', AMOUNT_FORMAT)
ts3('TS322', 'Total Patient Reimbursement Amount: refunds Medicare made to patients on behalf of the provider.', 'Part A only.', '5', AMOUNT_FORMAT)
ts3('TS323', 'Total PIP Claim Count: the claims paid under periodic interim payment.', 'Part A only; sent when the provider was on PIP for any part of the period.', '5', COUNT_FORMAT)
ts3('TS324', 'Total PIP Adjustment Amount: the periodic interim payment adjustment.', 'Part A only; sent when the provider was on PIP for any part of the period.', '5', AMOUNT_FORMAT)

put('summary:TS2', {
  is: 'Provider Supplemental Summary Information: inpatient prospective payment totals for the provider and type of bill over the period.',
  medicare: 'Part A only, sent when there were inpatient PPS payments; the flat file marks it not used for Part B.',
  sources: [FF('5'), IOM03(18), IOM03(37)],
})
const ts2 = (ref: string, is: string, when: string, values = AMOUNT_FORMAT) =>
  put(`summary:${ref}`, { is, values, medicare: `Part A only; ${when}`, sources: [FF("5"), IOM03(ref <= "TS212" ? 18 : 19)] })
ts2('TS201', 'Total DRG Amount: the diagnosis-related group payments added together.', 'the one required position when the segment is sent.')
ts2('TS202', 'Total Federal Specific Amount: the federal-specific operating DRG amounts.', 'sent when any were paid.')
ts2('TS203', 'Total Hospital Specific Amount: the hospital-specific operating DRG amounts.', 'sent when any were paid.')
ts2('TS204', 'Total Disproportionate Amount: disproportionate share payments.', 'sent when any were paid.')
ts2('TS205', 'Total Capital Amount: capital payments other than capital outliers.', 'sent when any were paid.')
ts2('TS206', 'Total Indirect Medical Education Amount.', 'sent when any were paid.')
ts2('TS207', 'Total Outlier Day Count.', 'sent when day outlier payments were made.', COUNT_FORMAT)
ts2('TS208', 'Total Day Outlier Amount.', 'sent when day outlier payments were made.')
ts2('TS209', 'Total Cost Outlier Amount.', 'sent when cost outlier payments were made.')
ts2('TS210', 'Average DRG Length of Stay: the arithmetic average length of stay for the DRGs in the transmission.', 'sent when DRG payments were made.', COUNT_FORMAT)
ts2('TS211', 'Total Discharge Count.', 'sent when there were discharges.', COUNT_FORMAT)
ts2('TS212', 'Total Cost Report Day Count.', 'sent when there were cost report days.', COUNT_FORMAT)
ts2('TS213', 'Total Covered Day Count.', 'sent when there were covered days.', COUNT_FORMAT)
ts2('TS214', 'Total Noncovered Day Count.', 'sent when there were non-covered days.', COUNT_FORMAT)
ts2('TS215', 'Total MSP Pass-Through Amount.', 'sent when MSP pass-through amounts applied.')
ts2('TS216', 'Average DRG Weight.', 'sent when DRG payments were made; the flat file holds it as S9(3)V9999.', COUNT_FORMAT)
ts2('TS217', 'Total PPS Capital FSP DRG Amount: the federal-specific portion of capital DRG payments.', 'sent when any were paid.')
ts2('TS218', 'Total PPS Capital HSP DRG Amount: the hospital-specific portion of capital DRG payments.', 'sent when any were paid.')
ts2('TS219', 'Total PPS DSH DRG Amount: the disproportionate share portion of DRG payments.', 'sent when any were paid.')

// ---------------------------------------------------------------- loop 2100, the claim

put('claim:CLP', {
  is: 'Claim Payment Information: the segment that opens one claim\'s loop. It carries the claim\'s identifiers, its status, the charge, the payment, and the patient\'s share, and everything that follows until the next CLP or the PLB belongs to this claim.',
  purpose:
    'The remittance itemizes each claim and service line so the provider can match the adjudication decisions to the claims it submitted. The validator pairs MAC and vendor claims by a match key built from six fields: one of them here (the total charge, CLP03), and the rest in the claim\'s NM1*QC and DTM*232 (or the earliest line date), the payee loop (NPI and tax ID), and the payer loop (the MAC workload).',
  sources: [FF('5'), FF('6'), CEDI(25, 19), MLN(2), API_DOC],
})
put('claim:CLP01', {
  is: 'Patient Control Number: the claim identifier the provider assigned when it billed, returned unchanged (the paper remittance prints it as the account number, ACNT).',
  purpose: 'It lets the provider find the claim in its own accounts receivable.',
  values: 'Up to 38 characters. Carriers print a single 0 when none was received on the 837 (CLM01).',
  medicare: 'The validator compares it as echoed claim data; it is not part of the match key.',
  sources: [FF('5'), IOM03(67), MREP(215), 'frontend/src/data/comparedFields.ts'],
})
put('claim:CLP02', {
  is: 'Claim Status Code: how the claim came out of adjudication: paid as primary, secondary, or tertiary, denied, forwarded on, or reversed.',
  purpose: 'It is the outcome state of the claim. The validator compares it, and its summaries count the pairs whose status codes agree.',
  values:
    'Medicare sends 1, 2, 3, 4, 19, 20, 21, 22, or 23. The Easy Print guide names 1 (Processed as Primary), 2 (Processed as Secondary), 19 (Processed as Primary, Forwarded to Additional Payer(s)), 20 (Processed as Secondary, Forwarded to Additional Payer(s)), and 21 (Processed as Tertiary, Forwarded to Additional Payer(s)); the MLN booklet calls 22 a reversal of previous payment. Codes 3, 4, and 23 are listed by the flat file without text in the documents on file. 25 (Predetermination Pricing Only, No Payment) and 27 do not apply to Medicare.',
  sources: [FF('5'), CEDI(25, 19), TI(7), IOM03(37), MREP(101), MREP(103), MREP(104), MLN(10), API_DOC],
})
put('claim:CLP03', {
  is: 'Total Claim Charge Amount: what the provider billed for the whole claim.',
  purpose: 'Balancing starts here: the total paid must equal the submitted charges plus or minus the adjustments reported, at the claim level and across the transaction.',
  values: AMOUNT_FORMAT,
  medicare: 'A match-key field: the validator pairs claims on it and leaves it out of the comparison, since a pair is equal on it by construction.',
  sources: [FF('5'), IOM(4), IOM03(68), MREP(216), API_DOC],
})
put('claim:CLP04', {
  is: 'Claim Payment Amount: what Medicare pays the provider for this claim (the paper remittance\'s PROV PD, the calculated payment to the provider).',
  purpose:
    'The sum of every CLP04 in the transaction, less the provider-level adjustments, is the check amount. On an adjusted claim the remittance shows the amount paid for the voided claim as previously paid; no remittance is issued for the voided claim itself.',
  values: AMOUNT_FORMAT + ' Negative on a reversal.',
  medicare: 'The validator compares it, and the dollar numbers in its summaries are built from it.',
  sources: [FF('6'), IOM(3), IOM03(69), MREP(217), MLN(8), API_DOC],
})
put('claim:CLP05', {
  is: 'Patient Responsibility Amount: the part of the claim the patient owes (PT RESP on the paper remittance).',
  purpose: 'Providers may bill the beneficiary only for amounts adjusted under group code PR; this element totals that share at the claim level.',
  values: AMOUNT_FORMAT,
  medicare: 'Carriers send it when there is any patient financial responsibility for amounts Medicare did not pay; the 2003 intermediary document says it does not apply to intermediaries.',
  sources: [FF('6'), MLN(2), IOM03(19), IOM03(37), MREP(216)],
})
put('claim:CLP06', {
  is: 'Claim Filing Indicator Code: the kind of coverage the claim was processed under.',
  values: 'MA, Medicare Part A, required for Part A; MB, Medicare Part B, required for Part B and sent for DME.',
  sources: [FF('6'), CEDI(25, 19), TI(7), IOM03(19), IOM03(38)],
})
put('claim:CLP07', {
  is: "Payer Claim Control Number: the payer's own identifier for the claim, the internal control number (ICN) on a Medicare remittance.",
  purpose: 'It is the number to quote back to the contractor about this claim.',
  values: 'Up to 50 characters.',
  medicare: 'Required for Medicare. The validator leaves it out of the comparison: each system assigns its own.',
  sources: [FF('6'), IOM03(19), IOM03(67), MREP(215), 'frontend/src/data/comparedFields.ts'],
})
put('claim:CLP08', {
  is: 'Facility Code Value: the facility type of the claim (the paper remittance prints it as the type of bill).',
  values: 'Two characters; the documents on file list no values for it.',
  medicare: 'Required for Medicare; the 2011 flat-file change log widened its description to cover Part B.',
  sources: [FF("6"), IOM03(19), IOM03(38), IOM03(59), T907(281)],
})
put('claim:CLP09', {
  is: 'Claim Frequency Code: the third position of the type of bill. The 837I carries the same code in CLM05-3; its values are NUBC code values the documents on file do not list.',
  values: 'One character; the documents on file list no values for it.',
  medicare: 'Intermediaries send it; it does not apply to carriers. The flat file marks it required when the information was received on the claim.',
  sources: [FF('6'), IOM03(19), IOM03(38), T907(23)],
})
put('claim:CLP10', {
  is: 'Patient Status Code.',
  medicare: NOT_USED_MEDICARE + ' The validator leaves it out of the comparison for that reason.',
  sources: [FF('6'), 'frontend/src/data/comparedFields.ts'],
})
put('claim:CLP11', {
  is: 'Diagnosis Related Group (DRG) Code: the group the inpatient stay was paid under.',
  values: 'Up to four characters.',
  medicare: 'Part A only, sent when a DRG payment was made.',
  sources: [FF('6'), IOM03(19)],
})
put('claim:CLP12', {
  is: 'DRG Weight: the relative weight of that DRG.',
  values: 'A decimal; the flat file holds it as S9(3)V9999.',
  medicare: 'Part A only, sent when a DRG payment was made.',
  sources: [FF("6"), IOM03(20), T907(259)],
})
put('claim:CLP13', {
  is: 'Discharge Fraction.',
  values: 'A decimal; the flat file holds it as S9(4)V999.',
  medicare: 'Sent when a discharge fraction was a factor in the payment to an institution.',
  sources: [FF("6"), IOM03(20), T907(259)],
})
put('claim:CLP14', {
  is: 'Yes/No Condition or Response Code.',
  medicare: NOT_USED_MEDICARE + ' The validator leaves it out of the comparison for that reason.',
  sources: [FF('6'), 'frontend/src/data/comparedFields.ts'],
})

const casAbout = (scope: 'claim' | 'line') => {
  const level = scope === 'claim' ? 'claim' : 'service line'
  const page = scope === 'claim' ? '6' : '12'
  const sources = scope === 'claim' ? [FF('6'), IOM(3), IOM(11), IOM03(38), MLN(8)] : [FF('12'), FF('13'), IOM(3), IOM(11), IOM03(44), MLN(8)]
  put(`${scope}:CAS`, {
    is: `${scope === 'claim' ? 'Claim' : 'Service'} Adjustment: the adjustments that explain why the ${level}'s payment differs from its charge. Each is a group code, a reason code, and an amount, with an optional quantity, and one segment holds up to six of them under one group code.`,
    purpose:
      'The difference between the submitted charge and the payment must be accounted for, or the 835 does not balance. A positive amount reduces the payment and a negative amount increases it. The same adjustment never appears at both the claim and the line level.',
    medicare:
      scope === 'claim'
        ? 'Inpatient remittances report adjustments at the claim level only. Carriers report every adjustment at the line level, so a Part B claim has no claim-level CAS.'
        : 'Carriers report every adjustment at the line level and send a line CAS whenever the amount paid for a service differs from the amount billed. Inpatient remittances carry no line-level adjustments.',
    sources,
  })
  put(`${scope}:CAS01`, {
    is: 'Claim Adjustment Group Code: the general category of the adjustment, which also says who is financially responsible for the amount.',
    values:
      'CO, Contractual Obligations: a contractual agreement or a regulatory requirement produced the adjustment; generally a write-off for the provider, not billed to the patient. OA, Other Adjustments: used when no other group applies. PR, Patient Responsibility: an amount that may be billed to the patient, typically deductible and copayment. Medicare contractors use only these three; PI is not used, and the flat file notes that CR has been deleted.',
    medicare: 'The validator labels each adjustment with its group and reason (CAS*CO*45), and the X12 list behind the codes is public.',
    sources: [FF(page), IOM(12), MLN(2), CEDI(scope === 'claim' ? 25 : 27, scope === 'claim' ? 19 : 21), TI(scope === 'claim' ? 7 : 9), X12_GROUPS],
  })
  const reason = (which: string) => ({
    is: `Adjustment Reason Code: the Claim Adjustment Reason Code (CARC) that explains the ${which} adjustment in this segment.`,
    purpose:
      'A reason code explains a financial adjustment: a denial, a reduction, or an increase in payment. The 835 allows one reason code per adjustment amount, and some codes are so general that a remark code must accompany them.',
    values: 'A code from the X12 CARC list, maintained by the Claim Adjustment Status and Reason Code Maintenance Committee and updated three times a year; the validator app shows each code\'s text on hover.',
    medicare: 'MACs must use the latest valid codes, and a code deactivated by the maintainers may still appear in a reversal that mirrors an older remittance.',
    sources: [FF(page), IOM(12), IOM(13), X12_CARCS],
  })
  const amount = (which: string) => ({
    is: `Adjustment Amount: the money adjusted for the ${which} reason code in this segment.`,
    values: AMOUNT_FORMAT + ' Positive reduces the payment; negative increases it.',
    medicare: 'The validator compares it under the group and reason beside it (for example CAS*PR*1) and sums the dollars on it in the summaries.',
    sources: [FF(page), IOM(3), MLN(8), API_DOC],
  })
  const quantity = (which: string) => ({
    is: `Adjustment Quantity: the units of service the ${which} adjustment applies to, when a quantity was adjusted.`,
    values: 'A decimal; the 2011 flat-file change log widened it to S9(7)V999 to allow fractional units.',
    medicare:
      'Intermediaries send it; the 2003 document had the visits of a home health episode of four or fewer reported here. Carriers do not use it. The validator compares it under the adjustment\'s own label with QTY added (CAS*CO*45*QTY).',
    sources: [FF(page), IOM03(scope === "claim" ? 20 : 29), IOM03(27), IOM03(44), T907(282), API_DOC],
  })
  const ordinals = ['first', 'second', 'third', 'fourth', 'fifth', 'sixth']
  for (let i = 0; i < 6; i++) {
    const base = 2 + i * 3
    const pad = (n: number) => String(n).padStart(2, '0')
    put(`${scope}:CAS${pad(base)}`, reason(ordinals[i]))
    put(`${scope}:CAS${pad(base + 1)}`, amount(ordinals[i]))
    put(`${scope}:CAS${pad(base + 2)}`, quantity(ordinals[i]))
  }
}
casAbout('claim')
casAbout('line')

// NM1 roles inside the claim. The name positions NM103 to NM107 are the same
// shape in every role; what differs is who is named and which identifier
// NM109 carries.
const nm1Role = (role: string, segment: About, names: { last: string; first: string; middle: string; suffix: string; page: string }, id: About, qualifier: About, entity: About) => {
  put(`claim:NM1@${role}`, segment)
  put(`claim:NM101@${role}`, { is: 'Entity Identifier Code: the role of the party this NM1 names.', values: `${role}.`, sources: [FF(names.page)] })
  put(`claim:NM102@${role}`, entity)
  put(`claim:NM103@${role}`, { is: names.last, values: 'Up to 60 characters.', medicare: 'Names are not compared by the validator; only the identifier in NM109 is.', sources: [FF(names.page), 'frontend/src/data/comparedFields.ts'] })
  put(`claim:NM104@${role}`, { is: names.first, values: 'Up to 35 characters.', medicare: 'Not compared.', sources: [FF(names.page), 'frontend/src/data/comparedFields.ts'] })
  put(`claim:NM105@${role}`, { is: names.middle, values: 'Up to 25 characters.', medicare: 'Not compared.', sources: [FF(names.page), 'frontend/src/data/comparedFields.ts'] })
  put(`claim:NM106@${role}`, { is: 'Name Prefix.', medicare: NOT_USED_MEDICARE, sources: [FF(names.page)] })
  put(`claim:NM107@${role}`, { is: names.suffix, values: 'Up to 10 characters.', medicare: 'Not compared.', sources: [FF(names.page), 'frontend/src/data/comparedFields.ts'] })
  put(`claim:NM108@${role}`, qualifier)
  put(`claim:NM109@${role}`, id)
}
const PERSON_OR_ENTITY = (values: string, page: string, extra: string[] = []) => ({
  is: 'Entity Type Qualifier: whether the party named is a person or an organization.',
  values: `${values} The 837 names the codes 1, Person, and 2, Non-Person Entity.`,
  sources: [FF(page), T907(9), ...extra],
})
nm1Role(
  'QC',
  {
    is: 'Patient Name: the beneficiary the claim is for, with the Medicare identifier the claim was processed under.',
    purpose: 'The remittance names the patient so the provider can tie the claim back to the person treated.',
    medicare: 'Protected health information: the validator never writes the identifier or the name to logs or addresses, and the beneficiary is never a grouping axis.',
    sources: [FF('6'), FF('7'), CEDI(25, 19), API_DOC],
  },
  { last: 'Patient Last Name.', first: 'Patient First Name.', middle: 'Patient Middle Name.', suffix: 'Patient Name Suffix.', page: '6' },
  {
    is: 'Patient Identifier: the Medicare identifier the claim was processed under.',
    values: 'The flat file says HIC number or MBI (the Medicare Beneficiary Identifier); up to 80 characters.',
    medicare:
      'A match-key field: the validator pairs claims on it and leaves it out of the comparison. Where a claim was submitted with a valid HICN, the MBI is returned in the Corrected Patient/Insured Name NM1*74 instead.',
    sources: [FF('7'), MLN(13), API_DOC, 'frontend/src/data/comparedFields.ts'],
  },
  {
    is: 'Identification Code Qualifier: the kind of identifier NM109 is.',
    values: 'MI, Member Identification Number: the DME guide says use MI, and the 837 gives the code the same name. The 2012 companion guide listed HN, the earlier qualifier.',
    sources: [FF('7'), CEDI(25, 19), CEDI837P(33, 25), TI(7)],
  },
  PERSON_OR_ENTITY('1, a person.', '6'),
)
nm1Role(
  '74',
  {
    is: 'Corrected Patient/Insured Name: the patient\'s name or identifier as corrected during processing, when the claim came in with a different one.',
    purpose: 'It tells the provider what to bill under next time. For Medicare the insured is the patient.',
    medicare:
      'Carriers send it when the name as received was corrected; intermediaries when the Medicare number was corrected. During the MBI transition the MBI is returned here when the claim carried a valid HICN.',
    sources: [FF('7'), IOM03(39), IOM03(22), MLN(13)],
  },
  { last: 'Corrected Patient/Insured Last Name.', first: 'Corrected Patient/Insured First Name.', middle: 'Corrected Patient/Insured Middle Name.', suffix: 'Corrected Patient Name Suffix.', page: '7' },
  {
    is: 'Corrected Insured Identification Indicator: the corrected identifier (an MBI, when one is returned).',
    values: 'Up to 80 characters. The flat file points to the MBI instructions of CR 9700 for when to fill it.',
    medicare: 'The validator compares it under NM1*74.',
    sources: [FF('7'), MLN(13), 'frontend/src/data/comparedFields.ts'],
  },
  {
    is: 'Identification Code Qualifier: the kind of identifier NM109 is.',
    values: 'C, the code the flat file lists for a corrected insured identifier. ' + SILENT,
    sources: [FF('7')],
  },
  PERSON_OR_ENTITY('1, a person; the 2012 companion guide says 2 does not apply to Medicare here.', '7', [TI(7)]),
)
nm1Role(
  '82',
  {
    is: 'Service Provider Name: the rendering provider, when it is not the payee.',
    purpose: 'It says who performed the services when that differs from who is paid. A rendering provider that differs for one service is named at that line instead.',
    medicare:
      'Medicare reports the number and not the name (the 2003 carrier document says do not report NM103 to NM105). When more than one rendering provider other than the payee performed services, the claim level names the one who performed more of them, or the first listed.',
    sources: [FF('7'), IOM03(39), TI(7)],
  },
  { last: 'Rendering Provider Last or Organization Name.', first: 'Rendering Provider First Name.', middle: 'Rendering Provider Middle Name.', suffix: 'Rendering Provider Name Suffix.', page: '7' },
  {
    is: "Rendering Provider Identifier: the rendering provider's NPI.",
    values: 'The NPI; up to 80 characters.',
    medicare: 'The validator compares it under NM1*82 as echoed claim data.',
    sources: [FF('7'), TI(7), 'frontend/src/data/comparedFields.ts'],
  },
  {
    is: 'Identification Code Qualifier: the kind of identifier NM109 is.',
    values: 'XX, the National Provider Identifier.',
    sources: [FF('7'), TI(7), CEDI(32, 26)],
  },
  PERSON_OR_ENTITY('1 or 2.', '7'),
)
nm1Role(
  'TT',
  {
    is: 'Crossover Carrier Name: the supplemental payer Medicare forwarded the claim to under a coordination of benefits agreement.',
    purpose: 'It tells the provider where the claim went next. The 835 can name only one crossover payer; when the claim went to more than one, remark code N89 says so in a claim-level remark position.',
    medicare: 'Sent when claim data is transferred to another payer under a COB agreement.',
    sources: [FF('7'), IOM(3), TI(8), CEDI(26, 20), IOM03(40)],
  },
  { last: 'COB Carrier Name: the name of the payer the claim was forwarded to.', first: 'First Name: not used at this role (the flat file marks NM104 to NM107 N/U).', middle: 'Middle Name: not used at this role.', suffix: 'Name Suffix: not used at this role.', page: '7' },
  {
    is: "COB Carrier Identifier: the crossover payer's identifier.",
    values: 'Up to 80 characters: the PlanID when effective, otherwise another number with PI.',
    medicare: 'The validator compares it under NM1*TT and records it like every other element; it is one of the four watched elements. Whether an element belongs in the compared set is settled from the real files, and the set is named on every run, so a change is visible.',
    sources: [FF('7'), IOM03(40), 'frontend/src/data/comparedFields.ts'],
  },
  {
    is: 'Identification Code Qualifier: the kind of identifier NM109 is.',
    values: 'PI, Payer Identification, or XV, the CMS PlanID. AD, FI, NI, and PP do not apply to Medicare.',
    sources: [FF('7'), TI(8), CEDI(26, 20), CEDI837P(35, 27), T907(9)],
  },
  PERSON_OR_ENTITY('2, an organization.', '7'),
)
nm1Role(
  'PR',
  {
    is: 'Corrected Priority Payer Name: the payer that should have been billed first, named when Medicare denied or rejected the claim because another payer is primary.',
    purpose: 'This segment notifies the provider whom to bill first. It is not used when the crossover carrier segment applies.',
    sources: [FF('7'), FF('8'), IOM03(23), IOM03(40)],
  },
  { last: 'Corrected Priority Payer Name.', first: 'First Name: not used at this role (the flat file marks NM104 to NM107 N/U).', middle: 'Middle Name: not used at this role.', suffix: 'Name Suffix: not used at this role.', page: '8' },
  {
    is: "Corrected Priority Payer ID: the primary payer's identifier.",
    values: 'Up to 80 characters: the PlanID when effective; before that, carriers zero-filled it.',
    medicare: 'The validator compares it under NM1*PR.',
    sources: [FF('8'), IOM03(41), 'frontend/src/data/comparedFields.ts'],
  },
  {
    is: 'Identification Code Qualifier: the kind of identifier NM109 is.',
    values: 'PI, Payer Identification, or XV, the CMS PlanID. AD, FI, NI, and PP do not apply to Medicare.',
    sources: [FF('8'), TI(8), CEDI837P(35, 27), T907(9)],
  },
  PERSON_OR_ENTITY('2, an organization.', '7'),
)

put('claim:MIA', {
  is: 'Inpatient Adjudication Information: the inpatient claim\'s payment components and its claim-level remark codes. A claim carries MIA or MOA, never both.',
  purpose: 'Remark codes that apply to the whole claim go here on an inpatient claim, and the prospective payment amounts behind the payment are itemized.',
  medicare: 'Part A only, sent when there has been inpatient care; carriers never send it.',
  sources: [FF('8'), FF('9'), IOM(15), IOM03(24), IOM03(41)],
})
const mia = (ref: string, is: string, when: string, page: string, values = AMOUNT_FORMAT) =>
  put(`claim:${ref}`, { is, values, medicare: `Part A only; ${when}`, sources: [FF(page), IOM03(ref <= "MIA16" ? 24 : 25)] })
mia('MIA01', 'Covered Days or Visits Count.', 'required when the segment is sent; the 2003 document had intermediaries always enter zero.', '8', COUNT_FORMAT)
mia('MIA02', 'PPS Operating Outlier Amount.', 'sent when an operating outlier payment was made.', '8')
mia('MIA03', 'Lifetime Psychiatric Days Count.', 'sent when lifetime psychiatric days were used.', '9', COUNT_FORMAT)
mia('MIA04', 'Claim DRG Amount: the diagnosis-related group payment.', 'sent when a DRG payment was made.', '9')
put('claim:MIA05', {
  is: 'Claim Payment Remark Code: the first of up to five claim-level Remittance Advice Remark Codes (the others are MIA20 to MIA23).',
  purpose: 'Remark codes explain an adjustment further, or carry information no reason code expresses; an informational remark starting with Alert may stand without a group and reason code. Remark codes that apply to a whole inpatient claim are reported here.',
  values: 'A code from the X12 RARC list, maintained by CMS; the app shows each code\'s text on hover.',
  medicare: 'The validator compares the five positions as one ordered list under MIA.',
  sources: [FF('9'), IOM(15), X12_RARCS, API_DOC],
})
mia('MIA06', 'Claim DSH Amount: the disproportionate share payment.', 'sent when a disproportionate share amount was paid.', '9')
mia('MIA07', 'Claim MSP Pass-Through Amount.', 'sent when an MSP pass-through amount was paid.', '9')
mia('MIA08', 'Claim PPS Capital Amount.', 'sent when a PPS capital amount was paid.', '9')
mia('MIA09', 'PPS Capital FSP DRG Amount: the federal-specific portion of the capital payment.', 'sent when paid.', '9')
mia('MIA10', 'PPS Capital HSP DRG Amount: the hospital-specific portion of the capital payment.', 'sent when paid.', '9')
mia('MIA11', 'PPS Capital DSH DRG Amount: the disproportionate share portion of the capital payment.', 'sent when paid.', '9')
mia('MIA12', 'Old Capital Amount.', 'sent when an old capital amount was paid.', '9')
mia('MIA13', 'PPS Capital IME Amount: the indirect medical education portion of the capital payment.', 'sent when paid.', '9')
mia('MIA14', 'PPS Operating HSP Specific DRG Amount: the hospital-specific portion of the operating payment.', 'sent when paid.', '9')
mia('MIA15', 'Cost Report Day Count.', 'sent when cost report days apply.', '9', COUNT_FORMAT)
mia('MIA16', 'PPS Operating FSP Specific DRG Amount: the federal-specific portion of the operating payment.', 'sent when paid.', '9')
mia('MIA17', 'Claim PPS Outlier Amount.', 'sent when a PPS outlier amount was paid.', '9')
mia('MIA18', 'Claim Indirect Teaching Amount.', 'sent when an indirect teaching amount was paid.', '9')
mia('MIA19', 'Nonpayable Professional Component Amount: a professional component billed but not payable to this provider.', 'sent when one was billed.', '9')
for (const [ref, which] of [
  ['MIA20', 'second'],
  ['MIA21', 'third'],
  ['MIA22', 'fourth'],
  ['MIA23', 'fifth'],
] as const) {
  put(`claim:${ref}`, {
    is: `Claim Payment Remark Code: the ${which} claim-level remark code (MIA05 is the first).`,
    values: 'A code from the X12 RARC list.',
    medicare: `Part A only; sent when a ${which} claim-level remark code applies. Compared with MIA05 as one list.`,
    sources: [FF('9'), IOM03(25), X12_RARCS],
  })
}
mia('MIA24', 'PPS Capital Exception Amount.', 'sent when a PPS capital exception amount was paid.', '9')

put('claim:MOA', {
  is: 'Outpatient Adjudication Information: the claim-level remark codes and a few payment amounts for a claim that is not inpatient. A claim carries MOA or MIA, never both.',
  purpose: 'Remark codes that apply to a whole non-inpatient claim are reported here, such as an appeal-rights remark, or N89 when the claim was forwarded to more than one payer.',
  medicare: 'Intermediaries send it for non-inpatient care when a claim-level remark applies; carriers whenever any claim-level remark code applies. MOA01, MOA02, MOA08, and MOA09 are Part A only.',
  sources: [FF('10'), IOM(15), IOM03(25), IOM03(41)],
})
put('claim:MOA01', {
  is: 'Reimbursement Rate: the rate the claim was reimbursed at.',
  values: 'A decimal; the flat file holds it as S9(4)V9999.',
  medicare: 'Part A only, sent when reimbursement rate reporting applies; the flat file marks it not used for Part B.',
  sources: [FF('10'), IOM03(25), IOM03(41)],
})
put('claim:MOA02', {
  is: 'Claim HCPCS Payable Amount: the amount payable for the lines paid on a fee schedule.',
  values: AMOUNT_FORMAT,
  medicare: 'Part A only, sent when any line items were paid on a fee schedule basis; not used for Part B.',
  sources: [FF('10'), IOM03(25), IOM03(41)],
})
put('claim:MOA03', {
  is: 'Claim Payment Remark Code: the first of up to five claim-level Remittance Advice Remark Codes (MOA03 to MOA07).',
  purpose: 'Remark codes explain an adjustment further or carry information no reason code expresses; an informational remark starting with Alert may stand alone. Claim-level remarks for a non-inpatient claim are reported here.',
  values: 'A code from the X12 RARC list, maintained by CMS; the app shows each code\'s text on hover.',
  medicare: 'The validator compares the five positions as one ordered list under MOA.',
  sources: [FF('10'), IOM(15), X12_RARCS, MREP(215), API_DOC],
})
for (const [ref, which] of [
  ['MOA04', 'second'],
  ['MOA05', 'third'],
  ['MOA06', 'fourth'],
  ['MOA07', 'fifth'],
] as const) {
  put(`claim:${ref}`, {
    is: `Claim Payment Remark Code: the ${which} claim-level remark code (MOA03 is the first).`,
    values: 'A code from the X12 RARC list.',
    medicare: `Sent when a ${which} claim-level remark code applies. Compared with MOA03 as one list.`,
    sources: [FF('10'), IOM03(25), IOM03(41), X12_RARCS],
  })
}
put('claim:MOA08', {
  is: 'Claim ESRD Payment Amount: the end-stage renal disease payment.',
  values: AMOUNT_FORMAT,
  medicare: 'Part A only, sent when an ESRD payment was made.',
  sources: [FF('10'), IOM03(26), IOM03(41)],
})
put('claim:MOA09', {
  is: 'Nonpayable Professional Component Amount: a professional component billed but not payable to this provider.',
  values: AMOUNT_FORMAT,
  medicare: 'Part A only.',
  sources: [FF('10'), IOM03(26), IOM03(41)],
})

put('claim:REF@EA', {
  is: "Other Claim Related Identification: the provider's own medical record number for the claim, returned as sent.",
  purpose: 'Sent when the provider put a proprietary identification number on the claim, so the provider can find the record.',
  medicare: 'Intermediaries send it (the flat file marks the segment not used by Part B). The 837 names qualifier EA the Medical Record Number.',
  sources: [FF('10'), CEDI(26, 20), TI(8), IOM03(26), T579(37), T907(28)],
})
put('claim:REF01@EA', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values: 'EA, which the 837 names the Medical Record Number. Medicare also lists 28, 6P, and F8 at this position; 1L, 1W, 9A, 9C, BB, CE, G1, G3, and IG are not used. 28 and 6P are listed without text.',
  sources: [FF('10'), CEDI(26, 20), TI(8), T579(37)],
})
put('claim:REF02@EA', {
  is: 'Other Claim Related Identification: the medical record number.',
  values: 'Up to 50 characters.',
  medicare: 'The validator compares it as echoed claim data under REF*EA.',
  sources: [FF('10'), 'frontend/src/data/comparedFields.ts'],
})
put('claim:REF@F8', {
  is: "Other Claim Related Identification with qualifier F8: the original reference number, the payer claim control number of the claim this one corrects or reverses.",
  purpose: 'It links a correction to the claim it corrects. The 837 uses the same qualifier on a replacement or void claim to name the original claim\'s ICN.',
  medicare: 'The validator leaves it out of the comparison: it carries a control number each system assigns for itself.',
  sources: [FF('10'), CEDI(26, 20), T907(27), 'frontend/src/data/comparedFields.ts'],
})
put('claim:REF01@F8', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values: "F8, which the 837 names the Payer Claim Control Number: the original claim's control number.",
  sources: [FF('10'), T907(27)],
})
put('claim:REF02@F8', {
  is: "Other Claim Related Identification: the original claim's control number.",
  values: 'Up to 50 characters.',
  medicare: 'Left out of the comparison.',
  sources: [FF('10'), 'frontend/src/data/comparedFields.ts'],
})

const claimDate = (qualifier: string, segment: About, kind: string, date: About) => {
  put(`claim:DTM@${qualifier}`, segment)
  put(`claim:DTM01@${qualifier}`, { is: 'Date/Time Qualifier: the kind of date DTM02 is.', values: `${qualifier}, ${kind}.`, sources: [FF(qualifier === '050' ? '11' : '10')] })
  put(`claim:DTM02@${qualifier}`, date)
}
claimDate(
  '232',
  {
    is: 'Statement From Date: the first day of the period an institutional claim covers (the statement covers period of the bill).',
    medicare: 'Intermediaries send the statement dates at the claim level; carriers report the dates of care at the service level instead. The Easy Print guide falls back to this date when a line has no service date.',
    sources: [FF('10'), IOM03(26), IOM03(41), MREP(215)],
  },
  'Statement From Date',
  {
    is: 'Claim Date: the statement from date.',
    values: DATE_FORMAT,
    medicare: 'The date of service of the match key when present; the validator pairs claims on it and leaves it out of the comparison.',
    sources: [FF('10'), API_DOC, 'frontend/src/data/comparedFields.ts'],
  },
)
claimDate(
  '233',
  {
    is: 'Statement To Date: the last day of the period an institutional claim covers.',
    medicare: 'Intermediaries send it with the from date; the Easy Print guide falls back to it when a line has no end date.',
    sources: [FF('10'), IOM03(26), MREP(215)],
  },
  'Statement To Date',
  { is: 'Claim Date: the statement to date.', values: DATE_FORMAT, medicare: 'Compared as echoed claim data.', sources: [FF('10'), 'frontend/src/data/comparedFields.ts'] },
)
claimDate(
  '036',
  {
    is: 'Coverage Expiration Date: the date the patient\'s coverage ended, when that bears on the claim.',
    medicare: 'The flat file lists it for the claim loop; the 2003 carrier document says expiration of coverage information does not apply to carriers.',
    sources: [FF('10'), IOM03(41)],
  },
  'Expiration',
  { is: 'Coverage Expiration Date.', values: DATE_FORMAT, medicare: 'Compared.', sources: [FF('10'), 'frontend/src/data/comparedFields.ts'] },
)
claimDate(
  '050',
  {
    is: 'Claim Received Date: the date the contractor received the claim.',
    purpose: 'Receipt starts the clock for interest on a clean claim paid late, which the MLN booklet describes as interest after the 30th day from receipt.',
    medicare: 'Carriers always send it; intermediaries send it among the claim dates.',
    sources: [FF('11'), IOM03(26), IOM03(41), MLN(8)],
  },
  'Received',
  {
    is: 'Claim Received Date.',
    values: DATE_FORMAT,
    medicare: 'Compared and recorded like every other element; one of the four watched elements. Whether an element belongs in the compared set is settled from the real files, and the set is named on every run, so a change is visible.',
    sources: [FF('11'), 'frontend/src/data/comparedFields.ts'],
  },
)

put('claim:PER@CX', {
  is: 'Claim Contact Information: whom to contact about this claim.',
  medicare: 'Contractors may report contact information at the claim level or the payer level at their option, not both; the data comes from the contractor, not the claims system. The validator leaves the segment out of the comparison as payer contact data.',
  sources: [FF('11'), IOM03(42), 'frontend/src/data/comparedFields.ts'],
})
perContact('claim', 'CX', 'Claim')
put('claim:PER08@CX', {
  is: 'Communication Number Extension: the telephone extension of a number given earlier in the segment.',
  values: 'Up to 256 characters.',
  sources: [FF('11'), IOM03(42)],
})

put('claim:AMT', {
  is: 'Claim Supplemental Information: a claim-level amount of the kind the qualifier says, such as the covered charges, interest, or the amount the patient paid.',
  purpose: 'Amounts that are neither the charge, the payment, nor an adjustment, but that a provider reconciles against, such as the interest paid on the claim.',
  medicare: 'One AMT per qualifier that applies. Carriers use F5 and I; intermediaries the rest.',
  sources: [FF('11'), CEDI(26, 20), TI(9), IOM03(26), IOM03(42)],
})
put('claim:AMT01', {
  is: 'Amount Qualifier Code: what kind of amount AMT02 is.',
  values:
    'Medicare sends AU, DY, F5, I, NL, ZK, ZL, ZM, ZN, or ZO; D8, T, and T2 are not used. From the documents on file: AU is the covered charges (the paper remittance derives non-covered charges as CLP03 less this amount); DY the inpatient or partial hospitalization per diem rate; F5 the Patient Amount Paid (the 837 names it so, and carriers send it when the claim reported a patient payment); I the interest paid on the claim; ZK the hemophilia add-on; ZL the new-technology add-on and electroconvulsive therapy amount. NL, ZM, ZN, and ZO are listed without text.',
  sources: [FF("11"), CEDI(26, 20), TI(9), IOM03(26), IOM03(42), IOM03(61), IOM03(62), IOM03(69), T579(32), MREP(217)],
})
put('claim:AMT02', {
  is: 'Claim Supplemental Information Amount: the amount of the kind AMT01 says.',
  values: AMOUNT_FORMAT,
  medicare: 'The validator compares each kind under its own label (AMT*AU, AMT*F5, …). AMT*I and AMT*NL are watched, and are compared and recorded like every other kind. Whether an element belongs in the compared set is settled from the real files, and the set is named on every run, so a change is visible.',
  sources: [FF('11'), 'frontend/src/data/comparedFields.ts'],
})
put('claim:QTY', {
  is: 'Claim Supplemental Information Quantity: a claim-level count of the kind the qualifier says, such as covered days.',
  medicare: 'Intermediaries send one QTY per qualifier that applies; the segment does not apply to carriers.',
  sources: [FF('11'), CEDI(27, 21), TI(9), IOM03(27), IOM03(42)],
})
put('claim:QTY01', {
  is: 'Quantity Qualifier: what kind of count QTY02 is.',
  values:
    'Medicare sends CA, CD, LA, OU, ZK, ZL, ZM, ZN, or ZO; LE, NE, NR, PS, and VS are not used. The paper remittance prints CA as covered days or visits and the 2003 document mapped NA to non-covered days and ZK to the hemophilia covered units. The others are listed without text.',
  sources: [FF('11'), CEDI(27, 21), TI(9), IOM03(27), IOM03(59)],
})
put('claim:QTY02', {
  is: 'Claim Supplemental Information Quantity: the count of the kind QTY01 says.',
  values: COUNT_FORMAT,
  medicare: 'The validator compares each kind under its own label (QTY*CA, …).',
  sources: [FF('11'), 'frontend/src/data/comparedFields.ts'],
})

// ---------------------------------------------------------------- loop 2110, the service line

put('line:SVC', {
  is: 'Service Payment Information: opens one service line of the claim: the procedure as adjudicated, its charge, its payment, and its units, with the submitted procedure and units when they were changed.',
  purpose:
    'The remittance itemizes each service so the provider can match the line decisions to what it billed. The code reported on the claim must come back here even when it was invalid and the line is being adjusted for it, which is an exception to the rule that standard transactions carry only valid codes.',
  medicare: 'Carriers send a line for every service; intermediaries when service-level detail came on the claim, one loop per procedure. Inpatient remittances carry no line-level adjustments.',
  sources: [FF('11'), FF('12'), CEDI(27, 21), IOM(3), IOM03(27), IOM03(42), MLN(2)],
})
put('line:SVC01', {
  is: 'Composite Medical Procedure Identifier: the adjudicated procedure, as a qualifier (SVC01-1), the code (SVC01-2), and up to four modifiers (SVC01-3 to SVC01-6).',
  values:
    'SVC01-1 is HC, HP, N4, or NU: HC a HCPCS code, N4 a National Drug Code (the 2003 document says N4 waits on Medicare\'s use of NDCs for drugs), NU a revenue code, and HP the qualifier the 2012 guide lists alongside them (the 2003 document used ZZ for HIPPS codes in version 4010). Modifiers apply only to HCPCS codes. SVC01-7 (description) and SVC01-8 are not used.',
  medicare:
    'When a service is denied for an invalid HCPCS, HIPPS, NDC, or revenue code, the invalid submitted code is reported here. The paper remittance prints SVC01-2 as the paid procedure; a down-coded original appears in SVC06. The validator compares the composite.',
  sources: [FF('11'), FF('12'), CEDI(27, 21), TI(9), CEDI(31, 25), IOM03(27), IOM03(28), IOM03(43), MREP(215), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC02', {
  is: 'Line Item Charge Amount: what the provider billed for the line.',
  values: AMOUNT_FORMAT,
  medicare: 'Compared.',
  sources: [FF('11'), IOM03(67), MREP(215), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC03', {
  is: 'Line Item Provider Payment Amount: what Medicare pays for the line (the paper remittance\'s PROV PD at the line).',
  values: AMOUNT_FORMAT,
  medicare: 'Compared.',
  sources: [FF('11'), IOM03(68), MREP(216), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC04', {
  is: 'NUBC Revenue Code: the revenue code of the line, when the line carries both a HCPCS or NDC and a revenue code.',
  values: 'Up to 48 characters.',
  medicare: 'Intermediaries send it when both codes were reported for the same service; it does not apply to carriers. Compared as echoed claim data.',
  sources: [FF('11'), IOM03(28), IOM03(43), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC05', {
  is: 'Units of Service Paid Count: the units paid for the line (PD-NOS, the paid number of services, on the paper remittance).',
  values: COUNT_FORMAT,
  medicare: 'Required for Medicare; for home health visits paid per visit the covered visits are reported here. Compared.',
  sources: [FF('12'), IOM03(28), IOM03(31), MREP(215), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC06', {
  is: 'Composite Medical Procedure Identifier, submitted: the procedure as it came on the claim, sent only when adjudication changed the procedure or drug code (or a modifier).',
  purpose: 'When the paid code equals the submitted code only SVC01 appears; when they differ both do, and the paper remittance prints this one as the submitted procedure.',
  values: 'The same composite shape as SVC01. SVC06-7, the description, is populated when it was received on the claim.',
  medicare: 'Medicare would not change a revenue code or a HIPPS code during adjudication, so this appears for HCPCS and NDC changes. Compared as echoed claim data.',
  sources: [FF('12'), IOM03(28), IOM03(43), MLN(4), MREP(216), 'frontend/src/data/comparedFields.ts'],
})
put('line:SVC07', {
  is: 'Original Units of Service Count: the units as billed, sent only when the paid units differ from them.',
  values: COUNT_FORMAT,
  medicare: 'Compared as echoed claim data.',
  sources: [FF('12'), IOM03(29), IOM03(43), MREP(216), 'frontend/src/data/comparedFields.ts'],
})

put('line:DTM@472', {
  is: 'Service Date: the date, or the from and to dates, of the service on this line.',
  medicare: 'Carriers send the dates of care at the service level; for intermediaries only 472 applies. The paper remittance prints 150 as the from date and 151 as the through date.',
  sources: [FF('12'), IOM03(29), IOM03(41), IOM03(67), MREP(215)],
})
put('line:DTM01@472', {
  is: 'Date/Time Qualifier: the kind of date DTM02 is.',
  values: '472, 150, or 151: the 837 names 472 the service date, and the paper remittance reads 150 as the from date and 151 as the through date.',
  sources: [FF('12'), T579(84), IOM03(67)],
})
put('line:DTM02@472', {
  is: 'Service Date.',
  values: DATE_FORMAT,
  medicare: 'Compared as echoed claim data under its qualifier (DTM*472, DTM*150, DTM*151). When a claim has no statement from date, the match key takes its date of service from the earliest line date.',
  sources: [FF('12'), API_DOC, 'frontend/src/data/comparedFields.ts'],
})

put('line:REF@LU', {
  is: 'Service Identification: an identifier that bears on how the line was paid: the place of service, a rate code, or an ambulatory payment classification.',
  medicare: 'Required when the qualifier applies; the APC number is reported only with the first HCPCS in that APC, and a rate code when it factored in the payment.',
  sources: [FF('12'), CEDI(27, 21), TI(9), IOM03(30), IOM03(45)],
})
put('line:REF01@LU', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values:
    'LU, 1S, APC, or RB. The paper remittance reads LU as the place of service; the 2003 document reads RB as a rate code and APC as the ambulatory payment classification number; 1S is listed without text. BB, E9, G1, and G3 are not used.',
  sources: [FF('12'), CEDI(27, 21), TI(9), IOM03(30), IOM03(67), MREP(215)],
})
put('line:REF02@LU', {
  is: 'Reference Identification: the identifier of the kind REF01 says.',
  values: 'Up to 50 characters.',
  medicare: 'Compared under its qualifier (REF*LU, REF*1S, REF*APC, REF*RB).',
  sources: [FF('12'), 'frontend/src/data/comparedFields.ts'],
})
put('line:REF@6R', {
  is: "Line Item Control Number: the provider's own identifier for the line, returned as sent.",
  purpose: 'It lets the provider pair each paid line with the line it billed. Medicare does not keep it in the core system; carriers fill it from the COB data repository, and intermediaries do not use it because their claims are not split.',
  medicare: 'The validator uses it to pair service lines (by REF*6R, then by position) and does not compare its value.',
  sources: [FF('13'), IOM03(30), IOM03(45), MREP(216), API_DOC, 'frontend/src/data/comparedFields.ts'],
})
put('line:REF01@6R', { is: 'Reference Identification Qualifier: the kind of identifier.', values: '6R; the flat file names the segment Line Item Control Number.', sources: [FF('13')] })
put('line:REF02@6R', {
  is: 'Line Item Control Number: the identifier itself.',
  values: 'Up to 50 characters.',
  medicare: 'Pairs the lines; not compared as a value.',
  sources: [FF('13'), 'frontend/src/data/comparedFields.ts'],
})
put('line:REF@HPI', {
  is: 'Rendering Provider Information: the provider who performed this service, when that is someone other than the payee and other than the rendering provider named at the claim level.',
  medicare: 'The rendering provider at the service level is identified by one of the NPI, the Social Security Number, the taxpayer identification number, or the Medicare provider number.',
  sources: [FF('13'), CEDI(28, 22), TI(9), IOM(2), IOM03(45)],
})
put('line:REF01@HPI', {
  is: 'Reference Identification Qualifier: the kind of identifier.',
  values:
    'HPI, SY, TJ, 1C, or 1G. HPI carries the NPI and 1C the Medicare provider number (the 2003 document moved carriers from 1C to HPI when the NPI took effect); 1G is the UPIN; TJ the taxpayer identification number. SY is listed as well; the manual names the Social Security Number among the four identifiers without pairing it with a qualifier. 0B, 1A, 1B, 1D, 1H, 1J, D3, and G2 are not used.',
  sources: [FF('13'), CEDI(28, 22), TI(9), TI(6), IOM(2), IOM03(37), IOM03(45)],
})
put('line:REF02@HPI', {
  is: "Rendering Provider ID: the provider's identifier of the kind REF01 says.",
  values: 'Up to 50 characters.',
  medicare: 'Compared as echoed claim data under its qualifier (REF*HPI, REF*1C, …).',
  sources: [FF('13'), 'frontend/src/data/comparedFields.ts'],
})
put('line:REF@0K', {
  is: 'Health Care Policy Identification: the coverage policy the line decision rests on.',
  medicare: 'Medicare reports the LCD or NCD code (a local or national coverage determination) in REF02. Reason code 97 refers readers to this segment.',
  sources: [FF('13'), CEDI(28, 22), TI(9), X12_CARCS],
})
put('line:REF01@0K', { is: 'Reference Identification Qualifier: the kind of identifier.', values: '0K (zero K); the flat file names the segment Health Care Policy Identification.', sources: [FF('13'), MREP(216)] })
put('line:REF02@0K', {
  is: 'Reference Identification: the NCD or LCD code.',
  values: 'Up to 50 characters.',
  medicare: 'Compared under REF*0K.',
  sources: [FF('13'), 'frontend/src/data/comparedFields.ts'],
})

put('line:AMT', {
  is: 'Service Supplemental Amount: a line-level amount of the kind the qualifier says, such as the allowed amount or a late filing charge.',
  medicare: 'One AMT per qualifier that applies; carriers use B6 and KH, intermediaries DY and B6.',
  sources: [FF('13'), CEDI(28, 22), TI(9), IOM03(31), IOM03(45)],
})
put('line:AMT01', {
  is: 'Amount Qualifier Code: what kind of amount AMT02 is.',
  values:
    'Medicare sends B6, KH, 2K, ZL, ZM, ZN, or ZO; T and T2 are not used. From the documents on file: B6 is the allowed amount for the service (the paper remittance\'s ALLOWED); KH the late filing charge; DY, which the 2003 intermediary document also lists, the per diem. 2K, ZL, ZM, ZN, and ZO are listed without text.',
  sources: [FF('13'), CEDI(28, 22), TI(9), IOM03(31), IOM03(45), IOM03(67), IOM03(69), MREP(215), MREP(217)],
})
put('line:AMT02', {
  is: 'Service Supplemental Amount: the amount of the kind AMT01 says.',
  values: AMOUNT_FORMAT,
  medicare: 'Compared under its own label (AMT*B6, AMT*KH, …).',
  sources: [FF('13'), 'frontend/src/data/comparedFields.ts'],
})

put('line:LQ@HE', {
  is: 'Health Care Remark Codes: one Remittance Advice Remark Code for the line; the segment repeats for each code.',
  purpose: 'Remark codes at the service line level are reported here, explaining an adjustment further or carrying information a reason code cannot.',
  medicare: 'Sent whenever any service-level remark code applies. The guide allows up to 99 per line; the flat file records up to 19.',
  sources: [FF('13'), CEDI(28, 22), TI(10), IOM(15), IOM03(31), IOM03(45)],
})
put('line:LQ01@HE', { is: 'Code List Qualifier Code: which code list LQ02 draws on.', values: 'HE, Claim Payment Remark Codes; only HE applies to Medicare.', sources: [FF('13'), CEDI(28, 22), TI(10)] })
put('line:LQ02@HE', {
  is: 'Remark Code: the RARC.',
  values: 'A code from the X12 RARC list, maintained by CMS; the app shows each code\'s text on hover.',
  medicare: 'The validator compares a line\'s remark codes as one ordered list under LQ*HE.',
  sources: [FF('13'), X12_RARCS, MREP(216), 'frontend/src/data/comparedFields.ts'],
})

// ---------------------------------------------------------------- the provider-level adjustments

put('provider:PLB', {
  is: 'Provider Level Adjustment: adjustments to the payment that do not belong to any one claim, such as interest on late payment, the recoupment of an earlier overpayment, or an incentive payment. One segment holds up to six, each a reason code, an identifier, and an amount.',
  purpose: 'The check amount is the claim payments less these adjustments, and every provider-level adjustment must be reported here so the transaction balances. A positive amount reduces the payment and a negative one increases it.',
  medicare: 'Sent whenever there have been provider-level adjustments. The validator keeps the count and total of a transaction\'s PLB adjustments on the payment row for its balance check and compares nothing in the segment, which sits outside the claims.',
  sources: [FF('13'), FF('14'), IOM(3), IOM(11), MLN(2), MLN(8), CEDI(29, 23), API_DOC],
})
put('provider:PLB01', {
  is: 'Provider Identifier: the provider the adjustments apply to.',
  values: 'The NPI, up to 50 characters; the flat file keeps the legacy number for VA pricing.',
  sources: [FF('13')],
})
put('provider:PLB02', {
  is: 'Fiscal Period Date: the end of the provider\'s fiscal period the adjustments fall in.',
  values: DATE_FORMAT + ' Carriers take it from their provider file or default to December 31 of the current year.',
  sources: [FF('13'), IOM03(46)],
})
const plbPair = (n: number, which: string) => {
  const pad = (x: number) => String(x).padStart(2, '0')
  put(`provider:PLB${pad(n)}`, {
    is: `Adjustment Identifier: a composite of the ${which} adjustment's reason code (PLB${pad(n)}-1) and a reference identifier for it (PLB${pad(n)}-2).`,
    values:
      'PLB-1 is a Provider Adjustment Reason Code from the X12 list; Medicare uses 50, 51, 72, 90, AP, B2, B3, BD, BN, C5, CS, CV, DM, E3, FB, GO, HM, IP, IS, IR, J1, L3, L6, LE, LS, OA, OB, PI, PL, RA, RE, SL, TL, WO, and WU, and not AH, AM, CR, CT, CW, or FC. PLB-2 is the financial control number or other identifier of the adjustment: the paper remittance prints positions 3 to 19 as the financial control number or ICN, and the 2003 carrier document put a two-character Medicare adjustment code in positions 1 and 2.',
    medicare: n === 3 ? 'Required when the segment is sent.' : `Sent when there is a ${which} provider-level adjustment.`,
    sources: [FF(n <= 5 ? '13' : '14'), CEDI(29, 23), TI(10), IOM03(32), IOM03(46), IOM03(70), MREP(218), X12_PLB],
  })
  put(`provider:PLB${pad(n + 1)}`, {
    is: `Provider Adjustment Amount: the amount of the ${which} adjustment.`,
    values: AMOUNT_FORMAT + ' Positive reduces the payment; negative increases it.',
    medicare: n === 3 ? 'Required when the segment is sent.' : `Sent when there is a ${which} provider-level adjustment.`,
    sources: [FF(n <= 5 ? '13' : '14'), IOM(3), MLN(8), IOM03(46)],
  })
}
plbPair(3, 'first')
plbPair(5, 'second')
plbPair(7, 'third')
plbPair(9, 'fourth')
plbPair(11, 'fifth')
plbPair(13, 'sixth')

// ---------------------------------------------------------------- trailer

put('trailer:SE', { is: 'Transaction Set Trailer: closes the 835 transaction that the ST opened.', sources: [FF('14'), CEDI(32, 26)] })
put('trailer:SE01', {
  is: 'Transaction Segment Count: the number of segments in the transaction, ST and SE included.',
  values: 'Up to ten digits; the carrier system computes it.',
  sources: [FF('14'), IOM03(47)],
})
put('trailer:SE02', { is: 'Transaction Set Control Number: repeats ST02.', values: 'Four to nine characters, equal to ST02.', sources: [FF('14')] })
put('trailer:GE', { is: 'Functional Group Trailer: closes the functional group the GS opened.', sources: [FF('14'), CEDI(31, 25)] })
put('trailer:GE01', { is: 'Number of Transaction Sets Included: how many ST-SE pairs the group holds (transactions, not claims).', values: 'Up to six digits.', sources: [FF('14'), FORMATS_README] })
put('trailer:GE02', { is: 'Group Control Number: repeats GS06.', values: 'Up to nine digits, equal to GS06.', sources: [FF('14')] })
put('trailer:IEA', { is: 'Interchange Control Trailer: closes the interchange the ISA opened.', sources: [FF('14'), CEDI(32, 26)] })
put('trailer:IEA01', { is: 'Number of Included Functional Groups: how many GS-GE pairs the interchange holds.', values: 'Up to five digits.', sources: [FF('14')] })
put('trailer:IEA02', { is: 'Interchange Control Number: repeats ISA13.', values: 'Nine digits, equal to ISA13.', sources: [FF('15')] })

export const ABOUT_835: Record<string, About> = entries

// The entry for a use key, or undefined when none is written.
export const aboutOf = (key: string): About | undefined => entries[key]

// ---------------------------------------------------------------- coded values with text of their own

// What a coded value means, beyond its name, where a document says more.
// Keyed by `${segment}:${element}:${code}`.
export const CODE_ABOUT_835: Record<string, { text: string; sources: string[] }> = {
  'CLP:CLP02:1': { text: 'Processed as Primary: Medicare paid as the primary payer.', sources: [MREP(104)] },
  'CLP:CLP02:2': { text: 'Processed as Secondary: Medicare paid after a primary payer (a Medicare Secondary Payer claim).', sources: [MREP(103)] },
  'CLP:CLP02:3': { text: 'Listed by the flat file as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('5')] },
  'CLP:CLP02:4': { text: 'Listed by the flat file as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('5')] },
  'CLP:CLP02:19': { text: 'Processed as Primary, Forwarded to Additional Payer(s): paid as primary and crossed over.', sources: [MREP(101)] },
  'CLP:CLP02:20': { text: 'Processed as Secondary, Forwarded to Additional Payer(s).', sources: [MREP(101)] },
  'CLP:CLP02:21': { text: 'Processed as Tertiary, Forwarded to Additional Payer(s).', sources: [MREP(101)] },
  'CLP:CLP02:22': { text: 'Reversal of previous payment: the claim as previously paid is taken back, with a negative payment, before the corrected claim is reported.', sources: [MLN(10), IOM(3)] },
  'CLP:CLP02:23': { text: 'Listed by the flat file as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('5')] },
  'CLP:CLP06:MA': { text: 'Medicare Part A: required on Part A claims.', sources: [TI(7)] },
  'CLP:CLP06:MB': { text: 'Medicare Part B: required on Part B claims and sent for DME.', sources: [TI(7), CEDI(25, 19)] },
  'LX:LX01:1': { text: 'Assigned claims (Part B and DME).', sources: [CEDI(25, 19), MREP(214)] },
  'LX:LX01:0': { text: 'Non-assigned claims (Part B and DME).', sources: [CEDI(25, 19), MREP(218)] },
  'BPR:BPR04:ACH': { text: 'Payment by electronic funds transfer through the automated clearing house; the bank fields of the BPR are then required.', sources: [IOM(4), IOM03(34)] },
  'BPR:BPR04:CHK': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('1'), CEDI(24, 18)] },
  'BPR:BPR04:NON': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('1'), CEDI(24, 18)] },
  'SVC:SVC01-1:HC': { text: 'A HCPCS code; modifiers apply only to this kind.', sources: [CEDI(31, 25), IOM03(28)] },
  "SVC:SVC01-1:N4": { text: "A National Drug Code.", sources: [IOM03(27)] },
  'SVC:SVC01-1:NU': { text: 'A revenue code.', sources: [IOM03(28)] },
  'SVC:SVC01-1:HP': { text: 'Listed by the 2012 guide as applying to Medicare; the documents on file do not spell it out.', sources: [TI(9)] },
  'NM1:NM102:1': { text: 'Person.', sources: [T907(9)] },
  'NM1:NM102:2': { text: 'Non-Person Entity.', sources: [T907(9)] },
  'NM1:NM108:MI': { text: 'Member Identification Number.', sources: [CEDI837P(33, 25)] },
  'NM1:NM108:PI': { text: 'Payer Identification.', sources: [CEDI837P(35, 27)] },
  'NM1:NM108:XV': { text: 'The CMS PlanID (named by the 837 as the Health Care Financing Administration National PlanID).', sources: [T907(9)] },
  'NM1:NM108:C': { text: 'Listed by the flat file for a corrected insured identifier; the documents on file do not spell it out.', sources: [FF('7')] },
  'N1:N103:FI': { text: 'Listed as a payee identifier qualifier; before the NPI took effect carriers sent it.', sources: [FF('3'), IOM03(36)] },
  'N1:N103:XV': { text: 'The CMS PlanID.', sources: [FF('2'), T907(9)] },
  'AMT:AMT01:AU': { text: 'Covered charges; the paper remittance derives the non-covered charges as CLP03 less this amount.', sources: [IOM03(61)] },
  "AMT:AMT01:DY": { text: "Per diem rate: the inpatient or partial hospitalization per diem.", sources: [IOM03(26), IOM03(62)] },
  'AMT:AMT01:F5': { text: 'Patient Amount Paid; carriers send it when the claim reported a patient payment.', sources: [T579(32), IOM03(42)] },
  'AMT:AMT01:I': { text: 'Interest paid on the claim.', sources: [IOM03(69), MREP(217)] },
  'AMT:AMT01:NL': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'AMT:AMT01:ZM': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11'), FF('13')] },
  'AMT:AMT01:ZN': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11'), FF('13')] },
  'AMT:AMT01:ZO': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11'), FF('13')] },
  'AMT:AMT01:2K': { text: 'Listed as a value Medicare sends at the line; the documents on file do not spell it out.', sources: [CEDI(28, 22), TI(9)] },
  'QTY:QTY01:CD': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:LA': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:OU': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:ZK': { text: 'The hemophilia covered units, per the 2003 companion document.', sources: [IOM03(27)] },
  'QTY:QTY01:ZL': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:ZM': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:ZN': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'QTY:QTY01:ZO': { text: 'Listed as a value Medicare sends; the documents on file do not spell it out.', sources: [FF('11')] },
  'REF:REF01:1S': { text: 'Listed as a service identification qualifier Medicare sends; the documents on file do not spell it out.', sources: [FF('12'), CEDI(27, 21)] },
  'REF:REF01:SY': { text: 'Listed as a rendering provider qualifier; the manual names the Social Security Number among the identifiers Medicare may send here without pairing it with a qualifier.', sources: [FF('13'), IOM(2)] },
  "AMT:AMT01:ZK": { text: "The hemophilia add-on amount.", sources: [IOM03(26)] },
  "AMT:AMT01:ZL": { text: "New-technology add-on and electroconvulsive therapy amount.", sources: [IOM03(26)] },
  'AMT:AMT01:B6': { text: 'Allowed amount for the service.', sources: [IOM03(67), MREP(215)] },
  'AMT:AMT01:KH': { text: 'Late filing charge.', sources: [IOM03(69), MREP(217)] },
  'QTY:QTY01:CA': { text: 'Covered days or visits.', sources: [IOM03(59)] },
  'DTM:DTM01:150': { text: 'The from date of the service, as the paper remittance reads it.', sources: [IOM03(67), MREP(215)] },
  'DTM:DTM01:151': { text: 'The through date of the service, as the paper remittance reads it.', sources: [IOM03(67), MREP(215)] },
  'REF:REF01:LU': { text: 'The place of service, as the paper remittance reads it.', sources: [IOM03(67), MREP(215)] },
  'REF:REF01:RB': { text: 'A rate code that factored in the payment.', sources: [IOM03(30)] },
  'REF:REF01:APC': { text: 'The ambulatory payment classification number, reported with the first HCPCS in that APC.', sources: [IOM03(30)] },
  'REF:REF01:1C': { text: 'The Medicare provider number.', sources: [IOM03(37), IOM03(45)] },
  'REF:REF01:1G': { text: 'The UPIN.', sources: [IOM03(37)] },
  'REF:REF01:HPI': { text: 'The NPI, as the rendering provider identifier.', sources: [IOM03(45), IOM(2)] },
  'REF:REF01:TJ': { text: "Federal Taxpayer's Identification Number.", sources: [TI(6)] },
  'REF:REF01:EA': { text: 'Medical Record Number: the provider\'s own identifier for the record.', sources: [T579(37), IOM03(26)] },
  'REF:REF01:F8': { text: 'Original Reference Number: the control number of the claim this one corrects.', sources: [T907(27)] },
  'REF:REF01:2U': { text: 'The qualifier of the payer loop\'s additional identification; for DME the identifier is the contractor ID.', sources: [CEDI(24, 18)] },
  'REF:REF01:EV': { text: 'Receiver identification, matching ISA08.', sources: [IOM03(14)] },
  'REF:REF01:F2': { text: 'Version Code: a local version number of the payer\'s implementation.', sources: [IOM03(35)] },
  'REF:REF01:0K': { text: 'The LCD or NCD code of the coverage policy.', sources: [CEDI(28, 22)] },
  'LQ:LQ01:HE': { text: 'The health care remark codes list (RARC); only HE applies to Medicare.', sources: [TI(10)] },
  'CAS:CAS01:CO': { text: 'Contractual Obligations: a contractual agreement or regulatory requirement produced the adjustment; generally a provider write-off, not billed to the patient.', sources: [IOM(12), MLN(2)] },
  'CAS:CAS01:OA': { text: 'Other Adjustments: used when no other group code applies.', sources: [IOM(12)] },
  'CAS:CAS01:PR': { text: 'Patient Responsibility: an amount that may be billed to the patient, typically deductible and copayment.', sources: [IOM(12), MLN(2)] },
}
