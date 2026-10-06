// Examples in the sidebar group "Subscriber and patient" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import subscriberIsPatient from './files/subscriber-is-patient.837?raw'
import dependentPatient from './files/dependent-patient.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment.

export const SUBSCRIBER_AND_PATIENT: ViewerExample[] = [
  {
    // Opens to: Subscriber Name (2010BA).
    id: 'subscriber-is-patient',
    group: 'Subscriber and patient',
    title: 'Subscriber is patient',
    file: subscriberIsPatient,
    open: ['ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0.SUBSCRIBER'],
  },
  {
    // Opens to: Patient Name (2010CA).
    id: 'dependent-patient',
    group: 'Subscriber and patient',
    title: 'Dependent patient',
    file: dependentPatient,
    open: ['ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0.PATIENT-HL#0.PATIENT'],
  },
]
