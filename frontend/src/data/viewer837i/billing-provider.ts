// Examples in the sidebar group "Billing provider" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import billingProvider from './files/billing-provider.837?raw'
import payToAddressAndPlan from './files/pay-to-address-and-plan.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment.

export const BILLING_PROVIDER: ViewerExample[] = [
  {
    // Opens to: Billing Provider Name (2010AA).
    id: 'billing-provider',
    group: 'Billing provider',
    title: 'Billing provider',
    file: billingProvider,
    open: ['ISA.GS.ST.BILLING-HL#0.BILLING-PROVIDER'],
  },
  {
    // Opens to: Pay-to Plan Name (2010AC).
    id: 'pay-to-address-and-plan',
    group: 'Billing provider',
    title: 'Pay-to address and plan',
    file: payToAddressAndPlan,
    open: ['ISA.GS.ST.BILLING-HL#0.PAY-TO-PLAN'],
  },
]
