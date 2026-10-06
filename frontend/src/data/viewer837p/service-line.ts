// Examples in the sidebar group "Service line" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import durableMedicalEquipment from './files/durable-medical-equipment.837?raw'
import anesthesiaMinutes from './files/anesthesia-minutes.837?raw'
import drugBilledByNdc from './files/drug-billed-by-ndc.837?raw'
import laboratoryTest from './files/laboratory-test.837?raw'
import purchasedDiagnosticService from './files/purchased-diagnostic-service.837?raw'
import groupPracticeMultipleRenderers from './files/group-practice-multiple-renderers.837?raw'
import incidentToService from './files/incident-to-service.837?raw'
import ambulance from './files/ambulance.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const SERVICE_LINE: ViewerExample[] = [
  {
    // Opens to: Durable Medical Equipment Service (2400).
    id: 'durable-medical-equipment',
    group: 'Service line',
    title: 'Durable medical equipment',
    file: durableMedicalEquipment,
  },
  {
    // Opens to: Professional Service (2400).
    id: 'anesthesia-minutes',
    group: 'Service line',
    title: 'Anesthesia minutes',
    file: anesthesiaMinutes,
  },
  {
    // Opens to: Drug Identification (2410).
    id: 'drug-billed-by-ndc',
    group: 'Service line',
    title: 'Drug billed by NDC',
    file: drugBilledByNdc,
  },
  {
    // Opens to: Test Results (2400).
    id: 'laboratory-test',
    group: 'Service line',
    title: 'Laboratory test',
    file: laboratoryTest,
  },
  {
    // Opens to: Purchased Service Information (2400).
    id: 'purchased-diagnostic-service',
    group: 'Service line',
    title: 'Purchased diagnostic service',
    file: purchasedDiagnosticService,
  },
  {
    // Opens to: Rendering Provider Name (2420A).
    id: 'group-practice-multiple-renderers',
    group: 'Service line',
    title: 'Group practice, multiple renderers',
    file: groupPracticeMultipleRenderers,
  },
  {
    // Opens to: Supervising Provider Name (2310D).
    id: 'incident-to-service',
    group: 'Service line',
    title: 'Incident-to service',
    file: incidentToService,
  },
  {
    // Opens to: Ambulance Transport Information (2400).
    id: 'ambulance',
    group: 'Service line',
    title: 'Ambulance',
    file: ambulance,
  },
]
