// Examples in the sidebar group "Envelope and submitter" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import officeVisit from './files/office-visit.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const ENVELOPE_AND_SUBMITTER: ViewerExample[] = [
  {
    // Opens to: Claim Information (2300).
    id: 'office-visit',
    group: 'Envelope and submitter',
    title: 'Office visit',
    file: officeVisit,
  },
]
