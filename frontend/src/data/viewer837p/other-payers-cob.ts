// Examples in the sidebar group "Other payers (COB)" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import medicareSecondary from './files/medicare-secondary.837?raw'
import lineLevelOtherPayer from './files/line-level-other-payer.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const OTHER_PAYERS_COB: ViewerExample[] = [
  {
    // Opens to: Other Subscriber Information (2320).
    id: 'medicare-secondary',
    group: 'Other payers (COB)',
    title: 'Medicare secondary',
    file: medicareSecondary,
  },
  {
    // Opens to: Line Adjudication Information (2430).
    id: 'line-level-other-payer',
    group: 'Other payers (COB)',
    title: 'Line-level other payer',
    file: lineLevelOtherPayer,
  },
]
