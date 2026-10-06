// Examples in the sidebar group "Other payers (COB)" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import medicareSecondaryOutpatient from './files/medicare-secondary-outpatient.837?raw'
import medicareSecondaryInpatient from './files/medicare-secondary-inpatient.837?raw'
import lineLevelOtherPayerPayment from './files/line-level-other-payer-payment.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment. Every example below has exactly one billing provider, one
// subscriber, and one claim, so every `open` path shares this prefix.
const CLAIM0 = 'ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0.CLAIM#0'

export const OTHER_PAYERS_COB: ViewerExample[] = [
  {
    // Opens to: Other Subscriber Information (2320).
    id: 'medicare-secondary-outpatient',
    group: 'Other payers (COB)',
    title: 'Medicare secondary (outpatient)',
    file: medicareSecondaryOutpatient,
    open: [`${CLAIM0}.OTHER-SUBSCRIBER#0`],
  },
  {
    // Opens to: Other Subscriber Information (2320).
    id: 'medicare-secondary-inpatient',
    group: 'Other payers (COB)',
    title: 'Medicare secondary (inpatient)',
    file: medicareSecondaryInpatient,
    open: [`${CLAIM0}.OTHER-SUBSCRIBER#0`],
  },
  {
    // Opens to: Institutional Service Line (2400), the 2430 line
    // adjudication sub-loop.
    id: 'line-level-other-payer-payment',
    group: 'Other payers (COB)',
    title: 'Line-level other payer payment',
    file: lineLevelOtherPayerPayment,
    open: [`${CLAIM0}.LINE#0.LINE-ADJUDICATION`],
  },
]
