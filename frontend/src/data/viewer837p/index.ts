// The "837P viewer" page is a small library of worked examples, each a full
// raw 837P showing one thing — the 837P's own version of
// frontend/src/data/viewer835Examples.ts (read that file first for the
// ViewerExample shape this reuses). Examples live in one file per sidebar
// group under ./; add an example by adding one entry to its group's file.
// The sidebar lists the groups in the order below, matching
// docs/design/837p-examples.md Part 2's own group order (which lists the
// examples and the parts of the 837P each one covers).
import type { ViewerExample } from '../viewer835Examples'
import { ENVELOPE_AND_SUBMITTER } from './envelope-and-submitter'
import { BILLING_PROVIDER } from './billing-provider'
import { SUBSCRIBER_AND_PATIENT } from './subscriber-and-patient'
import { CLAIM } from './claim'
import { SERVICE_LINE } from './service-line'
import { OTHER_PAYERS_COB } from './other-payers-cob'

export const VIEWER_837P_EXAMPLES: ViewerExample[] = [
  ...ENVELOPE_AND_SUBMITTER,
  ...BILLING_PROVIDER,
  ...SUBSCRIBER_AND_PATIENT,
  ...CLAIM,
  ...SERVICE_LINE,
  ...OTHER_PAYERS_COB,
]
