// Examples in the sidebar group "Billing provider" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import payToAddress from './files/pay-to-address.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const BILLING_PROVIDER: ViewerExample[] = [
  {
    // Opens to: Pay-to Address Name (2010AB).
    id: 'pay-to-address',
    group: 'Billing provider',
    title: 'Pay-to address',
    file: payToAddress,
  },
]
