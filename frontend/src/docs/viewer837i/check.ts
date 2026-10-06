// A dev-only correctness check for each example's raw file: SE01 must equal
// the true ST..SE segment count, and every claim's CLM02 must equal the sum
// of that same claim's service lines' SV203 — the examples' balance rule
// (docs/design/837i-examples.md, top-of-document "Rules the examples
// follow"). A line whose revenue code (SV201) is '0001' ("Total Charges",
// UB-04 Form Locator 47 / CMS Claims Processing Manual Ch. 25) restates the
// claim's other lines rather than adding to them, so it is excluded from
// the sum — confirmed against every multi-line example below, including
// inpatient-stay (two real revenue lines plus a restating 0001 line).
// A file can carry more than one claim (interim-bill's two CLM occurrences);
// each is checked against only its own service lines. Run once per example
// at module load (see Viewer837I.tsx); throws on a mismatch so a bad
// example is caught immediately instead of silently mis-rendering.
import { splitSegments, type Segment } from './parse'

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

  let claimNumber = 0
  let current: { clm: Segment; sum: number } | null = null
  const finishClaim = () => {
    if (!current) return
    claimNumber += 1
    const clm02 = Number(current.clm.elements[1])
    const sum = Math.round(current.sum * 100) / 100
    if (clm02 !== sum) {
      ok = false
      details.push(`claim ${claimNumber} (${current.clm.elements[0]}): CLM02=${clm02} but sum(SV203, excluding revenue 0001)=${sum}`)
    } else {
      details.push(`claim ${claimNumber} (${current.clm.elements[0]}): CLM02=${clm02} matches sum(SV203)`)
    }
  }
  for (const seg of segments) {
    if (seg.id === 'CLM') {
      finishClaim()
      current = { clm: seg, sum: 0 }
    } else if (seg.id === 'SV2' && current && seg.elements[0] !== '0001') {
      current.sum += Number(seg.elements[2] ?? '0')
    }
  }
  finishClaim()

  return { ok, details }
}

export function verifyExamplesOrThrow(examples: { id: string; file: string }[]): void {
  const failures: string[] = []
  for (const example of examples) {
    const result = verifySample(example.file)
    if (!result.ok) failures.push(`${example.id}: ${result.details.join('; ')}`)
  }
  if (failures.length > 0) {
    throw new Error(`837I viewer: example data failed its balance check — ${failures.join(' | ')}`)
  }
}
