// A dev-only correctness check for each example's raw file — the 837P's own
// version of frontend/src/docs/viewer835/check.ts. SE01 must equal the true
// ST..SE segment count, and CLM02 (Total Claim Charge Amount) must equal the
// sum of every service line's SV102 (Line Item Charge Amount) in that
// claim, per docs/design/837p-examples.md's own balancing rule ("CLM02 =
// the sum of every service line's SV102 in that claim" — an 837P carries no
// payment of its own, so this is a charges-and-units check, not a
// charges-and-payments one). Run once per example at module load (see
// Viewer837P.tsx); throws on a mismatch so a bad example is caught
// immediately instead of silently mis-rendering.
import { splitSegments } from '../viewer/tree'

export interface CheckResult {
  ok: boolean
  details: string[]
}

export function verifySample(raw: string): CheckResult {
  const segments = splitSegments(raw)
  const details: string[] = []
  let ok = true

  const stIndex = segments.findIndex((s) => s.id === 'ST')
  const seIndex = segments.findIndex((s) => s.id === 'SE')
  const se = segments[seIndex]
  if (stIndex === -1 || seIndex === -1 || !se) {
    return { ok: false, details: ['missing ST or SE segment'] }
  }
  const trueCount = seIndex - stIndex + 1
  const se01 = Number(se.elements[0])
  if (se01 !== trueCount) {
    ok = false
    details.push(`SE01=${se01} but the true ST..SE segment count is ${trueCount}`)
  } else {
    details.push(`SE01=${se01} matches the true ST..SE segment count`)
  }

  const clm = segments.find((s) => s.id === 'CLM')
  const clm02 = clm ? Number(clm.elements[1]) : NaN
  const sv1Sum = segments.filter((s) => s.id === 'SV1').reduce((sum, s) => sum + Number(s.elements[1] ?? '0'), 0)
  const expected = Math.round(sv1Sum * 100) / 100
  if (clm02 !== expected) {
    ok = false
    details.push(`CLM02=${clm02} but sum(SV102)=${sv1Sum}`)
  } else {
    details.push(`CLM02=${clm02} matches sum(SV102)=${sv1Sum}`)
  }

  return { ok, details }
}

export function verifyExamplesOrThrow(examples: { id: string; file: string }[]): void {
  const failures: string[] = []
  for (const example of examples) {
    const result = verifySample(example.file)
    if (!result.ok) failures.push(`${example.id}: ${result.details.join('; ')}`)
  }
  if (failures.length > 0) {
    throw new Error(`837P viewer: example data failed its balance check — ${failures.join(' | ')}`)
  }
}
