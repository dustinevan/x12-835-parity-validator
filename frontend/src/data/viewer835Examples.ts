// The "835 viewer" page is a small library of worked examples, each a full
// raw 835 showing one thing. Examples live in one file per sidebar group
// under ./viewer835/ (their raw 835s in ./viewer835/files/<id>.835); add an
// example by adding one entry to its group's file. The sidebar lists the
// groups in the order below. docs/design/835-examples.md lists which
// examples exist and which parts of the 835 each one covers.
import { PAYMENT } from './viewer835/payment'
import { PAYER_AND_PAYEE } from './viewer835/payer-and-payee'
import { PROVIDER_SUMMARY } from './viewer835/provider-summary'
import { CLAIM } from './viewer835/claim'
import { SERVICE_LINE } from './viewer835/service-line'
import { PROVIDER_ADJUSTMENT } from './viewer835/provider-adjustment'

export interface ViewerExample {
  id: string
  group: string
  title: string
  // Raw X12 text: ISA through IEA.
  file: string
  // Box path IDs (see frontend/src/docs/viewer835/parse.ts) that stay open
  // when "Collapse all" is pressed — the part this example shows — beyond
  // the envelope (ISA/GS/ST). Every box starts open. A path is
  // the dotted chain of box types from the root, with '#n' for the nth
  // sibling of a repeating box (LX, CLP, SVC) and '-PR'/'-PE' for the N1
  // loops — e.g. 'ISA.GS.ST.LX#0.CLP#1' opens the second claim, and
  // 'ISA.GS.ST.LX#0.CLP#1.SVC#2' also opens that claim's third line.
  // Opening a box also opens its ancestors, so listing the deepest box you
  // want visible is enough.
  open?: string[]
}

export const VIEWER_835_EXAMPLES: ViewerExample[] = [
  ...PAYMENT,
  ...PAYER_AND_PAYEE,
  ...PROVIDER_SUMMARY,
  ...CLAIM,
  ...SERVICE_LINE,
  ...PROVIDER_ADJUSTMENT,
]
