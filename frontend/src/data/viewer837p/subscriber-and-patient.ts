// Examples in the sidebar group "Subscriber and patient" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import subscriberIsPatient from './files/subscriber-is-patient.837?raw'
import propertyAndCasualtyClaim from './files/property-and-casualty-claim.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const SUBSCRIBER_AND_PATIENT: ViewerExample[] = [
  {
    // Opens to: Subscriber Information (2000B).
    id: 'subscriber-is-patient',
    group: 'Subscriber and patient',
    title: 'Subscriber is the patient',
    file: subscriberIsPatient,
  },
  {
    // Opens to: Property and Casualty Claim Number (2010BA).
    id: 'property-and-casualty-claim',
    group: 'Subscriber and patient',
    title: 'Workers\' comp or auto accident claim',
    file: propertyAndCasualtyClaim,
  },
]
