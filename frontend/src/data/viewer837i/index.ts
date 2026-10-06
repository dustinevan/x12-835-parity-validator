// The "837I viewer" page is a small library of worked examples, each a full
// raw 837I showing one thing. Examples live in one file per sidebar group
// under ./viewer837i/ (their raw 837Is in ./viewer837i/files/<id>.837); add
// an example by adding one entry to its group's file. The sidebar lists the
// groups in the order below, matching docs/design/837i-examples.md Part 2's
// own group order (which lists the examples and the parts of the 837I
// each one covers). `ViewerExample` (id/group/title/file/open?) is the
// same shape the 835 viewer's examples use — see
// frontend/src/data/viewer835Examples.ts — reused as-is since nothing about
// it is transaction-specific.
import type { ViewerExample } from '../viewer835Examples'
import { ENVELOPE_AND_SUBMITTER } from './envelope-and-submitter'
import { BILLING_PROVIDER } from './billing-provider'
import { SUBSCRIBER_AND_PATIENT } from './subscriber-and-patient'
import { CLAIM } from './claim'
import { SERVICE_LINE } from './service-line'
import { OTHER_PAYERS_COB } from './other-payers-cob'

export type { ViewerExample }

export const VIEWER_837I_EXAMPLES: ViewerExample[] = [
  ...ENVELOPE_AND_SUBMITTER,
  ...BILLING_PROVIDER,
  ...SUBSCRIBER_AND_PATIENT,
  ...CLAIM,
  ...SERVICE_LINE,
  ...OTHER_PAYERS_COB,
]
