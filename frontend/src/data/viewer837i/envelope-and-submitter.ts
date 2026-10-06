// Examples in the sidebar group "Envelope and submitter" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import submitterAndReceiver from './files/submitter-and-receiver.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment.

export const ENVELOPE_AND_SUBMITTER: ViewerExample[] = [
  {
    // Opens to: Submitter Name (1000A).
    id: 'submitter-and-receiver',
    group: 'Envelope and submitter',
    title: 'Submitter and receiver',
    file: submitterAndReceiver,
    open: ['ISA.GS.ST.SUBMITTER'],
  },
]
