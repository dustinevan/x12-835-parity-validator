// Examples in the sidebar group "Provider summary" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import institutionalTotals from './files/institutional-totals.835?raw'

export const PROVIDER_SUMMARY: ViewerExample[] = [
  {
    id: 'institutional-totals',
    group: 'Provider summary',
    title: 'Institutional totals',
    file: institutionalTotals,
    open: ['ISA.GS.ST.LX#0'],
  },
]
