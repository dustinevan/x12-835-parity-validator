// Examples in the sidebar group "Payment" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import paidByAch from './files/paid-by-ach.835?raw'
import paidByCheck from './files/paid-by-check.835?raw'
import noPaymentZeroTotal from './files/no-payment-zero-total.835?raw'
import multipleClaimsOnePayment from './files/multiple-claims-one-payment.835?raw'

export const PAYMENT: ViewerExample[] = [
  {
    id: 'paid-by-ach',
    group: 'Payment',
    title: 'Paid by ACH',
    file: paidByAch,
  },
  {
    id: 'paid-by-check',
    group: 'Payment',
    title: 'Paid by check',
    file: paidByCheck,
  },
  {
    id: 'no-payment-zero-total',
    group: 'Payment',
    title: 'No payment, zero total',
    file: noPaymentZeroTotal,
  },
  {
    id: 'multiple-claims-one-payment',
    group: 'Payment',
    title: 'Multiple claims, one payment',
    file: multipleClaimsOnePayment,
  },
]
