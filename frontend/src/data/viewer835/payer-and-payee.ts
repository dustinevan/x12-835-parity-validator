// Examples in the sidebar group "Payer and payee" (docs/design/835-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import payeeByNpi from './files/payee-by-npi.835?raw'
import payerContacts from './files/payer-contacts.835?raw'

export const PAYER_AND_PAYEE: ViewerExample[] = [
  {
    id: 'payee-by-npi',
    group: 'Payer and payee',
    title: 'Payee by NPI',
    file: payeeByNpi,
    open: ['ISA.GS.ST.N1-PE'],
  },
  {
    id: 'payer-contacts',
    group: 'Payer and payee',
    title: 'Payer contacts',
    file: payerContacts,
    open: ['ISA.GS.ST.N1-PR'],
  },
]
