// A dev-only correctness check for each example's raw file: SE01 must equal
// the true ST..SE segment count, and BPR02 must equal the sum of every
// CLP04 minus every PLB adjustment amount (PLB04, 06, 08, 10, 12, 14) — the
// same arithmetic the spec's closing paragraph states, extended to cover
// provider-level adjustments. Run once per example at module load (see
// Viewer835.tsx); throws on a mismatch so a bad example is caught
// immediately instead of silently mis-rendering.
import { splitSegments } from './parse'

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

  const bpr = segments.find((s) => s.id === 'BPR')
  const bpr02 = bpr ? Number(bpr.elements[1]) : NaN
  const clpSum = segments.filter((s) => s.id === 'CLP').reduce((sum, s) => sum + Number(s.elements[3]), 0)
  // PLB holds up to six reason/identifier/amount triples after PLB01/02; the
  // amount is the second element of each triple — elements[3,5,7,9,11,13].
  const plbSum = segments
    .filter((s) => s.id === 'PLB')
    .reduce((sum, s) => {
      let slotSum = 0
      for (let i = 3; i < s.elements.length; i += 2) {
        if (s.elements[i] !== undefined && s.elements[i] !== '') slotSum += Number(s.elements[i])
      }
      return sum + slotSum
    }, 0)
  const expected = Math.round((clpSum - plbSum) * 100) / 100
  if (bpr02 !== expected) {
    ok = false
    details.push(`BPR02=${bpr02} but sum(CLP04)=${clpSum} minus sum(PLB amounts)=${plbSum} = ${expected}`)
  } else {
    details.push(`BPR02=${bpr02} matches sum(CLP04)=${clpSum} minus sum(PLB amounts)=${plbSum}`)
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
    throw new Error(`835 viewer: example data failed its balance check — ${failures.join(' | ')}`)
  }
}
