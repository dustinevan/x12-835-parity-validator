// Examples in the sidebar group "Service line" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import drugByNdc from './files/drug-by-ndc.837?raw'
import lineItemControlAndRepricing from './files/line-item-control-and-repricing.837?raw'
import lineLevelProviders from './files/line-level-providers.837?raw'
import lineNotesAndAttachments from './files/line-notes-and-attachments.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment. Every example below has exactly one billing provider, one
// subscriber, and one claim, so every `open` path shares this prefix.
const CLAIM0 = 'ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0.CLAIM#0'

export const SERVICE_LINE: ViewerExample[] = [
  {
    // Opens to: Drug Identification (2410).
    id: 'drug-by-ndc',
    group: 'Service line',
    title: 'Drug by NDC',
    file: drugByNdc,
    open: [`${CLAIM0}.LINE#0.DRUG`],
  },
  {
    // Opens to: Institutional Service Line (2400) — both lines, one
    // repriced and one not.
    id: 'line-item-control-and-repricing',
    group: 'Service line',
    title: 'Line item control and repricing',
    file: lineItemControlAndRepricing,
    open: [`${CLAIM0}.LINE#0`, `${CLAIM0}.LINE#1`],
  },
  {
    // Opens to: Institutional Service Line (2400).
    id: 'line-level-providers',
    group: 'Service line',
    title: 'Line-level providers',
    file: lineLevelProviders,
    open: [`${CLAIM0}.LINE#0`],
  },
  {
    // Opens to: Institutional Service Line (2400).
    id: 'line-notes-and-attachments',
    group: 'Service line',
    title: 'Line notes and attachments',
    file: lineNotesAndAttachments,
    open: [`${CLAIM0}.LINE#0`],
  },
]
