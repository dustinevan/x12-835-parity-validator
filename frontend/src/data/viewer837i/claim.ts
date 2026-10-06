// Examples in the sidebar group "Claim" (docs/design/837i-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import inpatientStay from './files/inpatient-stay.837?raw'
import outpatientVisit from './files/outpatient-visit.837?raw'
import emergencyRoom from './files/emergency-room.837?raw'
import skilledNursing from './files/skilled-nursing.837?raw'
import homeHealth from './files/home-health.837?raw'
import hospice from './files/hospice.837?raw'
import outpatientTherapy from './files/outpatient-therapy.837?raw'
import dialysis from './files/dialysis.837?raw'
import interimBill from './files/interim-bill.837?raw'
import replacementClaim from './files/replacement-claim.837?raw'
import drg from './files/drg.837?raw'
import occurrenceSpan from './files/occurrence-span.837?raw'
import valueCodes from './files/value-codes.837?raw'
import conditionCodes from './files/condition-codes.837?raw'
import attendingAndOperatingPhysicians from './files/attending-and-operating-physicians.837?raw'
import externalCauseOfInjury from './files/external-cause-of-injury.837?raw'
import claimAttachmentsAndAuthorization from './files/claim-attachments-and-authorization.837?raw'

// Box-path convention: frontend/src/docs/viewer837i/parse.ts's header
// comment. Every example below has exactly one billing provider and one
// subscriber, so every `open` path shares this prefix.
const HL = 'ISA.GS.ST.BILLING-HL#0.SUBSCRIBER-HL#0'

export const CLAIM: ViewerExample[] = [
  {
    // Opens to: Claim Information (2300).
    id: 'inpatient-stay',
    group: 'Claim',
    title: 'Inpatient stay',
    file: inpatientStay,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'outpatient-visit',
    group: 'Claim',
    title: 'Outpatient visit',
    file: outpatientVisit,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Institutional Service Line (2400).
    id: 'emergency-room',
    group: 'Claim',
    title: 'Emergency room',
    file: emergencyRoom,
    open: [`${HL}.CLAIM#0.LINE#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'skilled-nursing',
    group: 'Claim',
    title: 'Skilled nursing',
    file: skilledNursing,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Institutional Service Line (2400).
    id: 'home-health',
    group: 'Claim',
    title: 'Home health',
    file: homeHealth,
    open: [`${HL}.CLAIM#0.LINE#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'hospice',
    group: 'Claim',
    title: 'Hospice',
    file: hospice,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Rendering Physician Name (2310D).
    id: 'outpatient-therapy',
    group: 'Claim',
    title: 'Outpatient therapy',
    file: outpatientTherapy,
    open: [`${HL}.CLAIM#0.RENDERING`],
  },
  {
    // Opens to: Institutional Service Line (2400).
    id: 'dialysis',
    group: 'Claim',
    title: 'Dialysis',
    file: dialysis,
    open: [`${HL}.CLAIM#0.LINE#0`],
  },
  {
    // Opens to: Claim Information (2300) — both claims, the interim pair.
    id: 'interim-bill',
    group: 'Claim',
    title: 'Interim bill',
    file: interimBill,
    open: [`${HL}.CLAIM#0`, `${HL}.CLAIM#1`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'replacement-claim',
    group: 'Claim',
    title: 'Replacement claim',
    file: replacementClaim,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'drg',
    group: 'Claim',
    title: 'DRG',
    file: drg,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'occurrence-span',
    group: 'Claim',
    title: 'Occurrence span',
    file: occurrenceSpan,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'value-codes',
    group: 'Claim',
    title: 'Value codes',
    file: valueCodes,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'condition-codes',
    group: 'Claim',
    title: 'Condition codes',
    file: conditionCodes,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Attending Provider Name (2310A).
    id: 'attending-and-operating-physicians',
    group: 'Claim',
    title: 'Attending and operating physicians',
    file: attendingAndOperatingPhysicians,
    open: [`${HL}.CLAIM#0.ATTENDING`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'external-cause-of-injury',
    group: 'Claim',
    title: 'External cause of injury',
    file: externalCauseOfInjury,
    open: [`${HL}.CLAIM#0`],
  },
  {
    // Opens to: Claim Information (2300).
    id: 'claim-attachments-and-authorization',
    group: 'Claim',
    title: 'Claim attachments and authorization',
    file: claimAttachmentsAndAuthorization,
    open: [`${HL}.CLAIM#0`],
  },
]
