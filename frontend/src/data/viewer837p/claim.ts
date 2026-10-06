// Examples in the sidebar group "Claim" (docs/design/837p-examples.md, Part 2).
import type { ViewerExample } from '../viewer835Examples'
import correctedClaim from './files/corrected-claim.837?raw'
import priorAuthorization from './files/prior-authorization.837?raw'
import specialistReferral from './files/specialist-referral.837?raw'
import attachment from './files/attachment.837?raw'
import multipleDiagnoses from './files/multiple-diagnoses.837?raw'
import repricedClaim from './files/repriced-claim.837?raw'
import specialProgramIdentifiers from './files/special-program-identifiers.837?raw'
import homeboundAndEpsdtReferral from './files/homebound-and-epsdt-referral.837?raw'
import vision from './files/vision.837?raw'
import chiropractic from './files/chiropractic.837?raw'

// The 837P viewer (frontend/src/docs/viewer837p/parse.ts) defines the
// box-path convention `open` uses. Each example's intended loop is named in
// the comment above it, matching the design doc's "Opens to" line.

export const CLAIM: ViewerExample[] = [
  {
    // Opens to: Claim Information (2300).
    id: 'corrected-claim',
    group: 'Claim',
    title: 'Corrected claim',
    file: correctedClaim,
  },
  {
    // Opens to: Prior Authorization Number (2300).
    id: 'prior-authorization',
    group: 'Claim',
    title: 'Prior authorization',
    file: priorAuthorization,
  },
  {
    // Opens to: Referring Provider Name (2310A).
    id: 'specialist-referral',
    group: 'Claim',
    title: 'Specialist referral',
    file: specialistReferral,
  },
  {
    // Opens to: Claim Supplemental Information (2300).
    id: 'attachment',
    group: 'Claim',
    title: 'Attachment',
    file: attachment,
  },
  {
    // Opens to: Health Care Diagnosis Code (2300).
    id: 'multiple-diagnoses',
    group: 'Claim',
    title: 'Multiple diagnoses',
    file: multipleDiagnoses,
  },
  {
    // Opens to: Claim Pricing/Repricing Information (2300).
    id: 'repriced-claim',
    group: 'Claim',
    title: 'Repriced claim',
    file: repricedClaim,
  },
  {
    // Opens to: Claim Information (2300).
    id: 'special-program-identifiers',
    group: 'Claim',
    title: 'Special program identifiers',
    file: specialProgramIdentifiers,
  },
  {
    // Opens to: Homebound Indicator (2300).
    id: 'homebound-and-epsdt-referral',
    group: 'Claim',
    title: 'Homebound and EPSDT referral',
    file: homeboundAndEpsdtReferral,
  },
  {
    // Opens to: Patient Condition Information: Vision (2300).
    id: 'vision',
    group: 'Claim',
    title: 'Vision',
    file: vision,
  },
  {
    // Opens to: Spinal Manipulation Service Information (2300).
    id: 'chiropractic',
    group: 'Claim',
    title: 'Chiropractic',
    file: chiropractic,
  },
]
