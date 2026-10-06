// Examples in the sidebar group "Provider adjustment" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import overpaymentRecovered from './files/overpayment-recovered.835?raw'
import interestToTheProvider from './files/interest-to-the-provider.835?raw'
import interimAndPeriodicPayments from './files/interim-and-periodic-payments.835?raw'
import costReportSettlementsAndPassthroughs from './files/cost-report-settlements-and-passthroughs.835?raw'

export const PROVIDER_ADJUSTMENT: ViewerExample[] = [
  {
    id: 'overpayment-recovered',
    group: 'Provider adjustment',
    title: 'Overpayment recovered',
    file: overpaymentRecovered,
  },
  {
    id: 'interest-to-the-provider',
    group: 'Provider adjustment',
    title: 'Interest to the provider',
    file: interestToTheProvider,
  },
  {
    id: 'interim-and-periodic-payments',
    group: 'Provider adjustment',
    title: 'Interim and periodic payments',
    file: interimAndPeriodicPayments,
  },
  {
    id: 'cost-report-settlements-and-passthroughs',
    group: 'Provider adjustment',
    title: 'Cost-report settlements and passthroughs',
    file: costReportSettlementsAndPassthroughs,
  },
]
