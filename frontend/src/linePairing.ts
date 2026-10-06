// How the service lines of one claim in up to three columns are paired, and
// what a column of such a grid is called. Used by the docs site's 835
// compare page (docs/compare835/rows.ts), whose rows the validator app's
// claim screen also draws (app/pages/ClaimCompare.tsx); it knows nothing of
// either page's rows. The first column is the one the others are paired
// against: the MAC claim in the app, the reference file in the docs. The
// other two are the app's two slots, v1 and v2; the docs page puts its two
// compared files in them.

export type Slot = 'v1' | 'v2'
export type Col = 'mac' | Slot
export const SIDES: Slot[] = ['v1', 'v2']

export type Sources<T> = Partial<Record<Col, T>>

// The comparison's own pairing (compare.go pairLines): lines whose REF*6R is
// present and unique on both sides pair by it; the rest pair in order.
// partner[i] is the vendor line paired with MAC line i, or -1.
export function pairLines<T extends { line_control_number: string }>(mac: T[], vendor: T[]) {
  const tally = (lines: T[]) => {
    const n = new Map<string, number>()
    for (const l of lines) if (l.line_control_number) n.set(l.line_control_number, (n.get(l.line_control_number) ?? 0) + 1)
    return n
  }
  const macCount = tally(mac)
  const vendorCount = tally(vendor)
  const vendorIndex = new Map(vendor.map((l, j) => [l.line_control_number, j]))
  const partner = mac.map(() => -1)
  const taken = vendor.map(() => false)
  mac.forEach((l, i) => {
    const c = l.line_control_number
    if (!c || macCount.get(c) !== 1 || vendorCount.get(c) !== 1) return
    const j = vendorIndex.get(c)!
    partner[i] = j
    taken[j] = true
  })
  let next = 0
  mac.forEach((_, i) => {
    if (partner[i] >= 0) return
    while (next < vendor.length && taken[next]) next++
    if (next < vendor.length) {
      partner[i] = next
      taken[next] = true
    }
  })
  return { partner, unpaired: vendor.flatMap((_, j) => (taken[j] ? [] : [j])) }
}

// One line of the grid: the MAC line with each vendor's paired line, or a
// vendor line with no MAC partner.
export interface PairedLine<T> {
  src: Sources<T>
  // The REF*6R the comparison records on this line's differences, per
  // slot: the MAC line's when there is one, else that slot's own line's.
  control: Partial<Record<Slot, string>>
  macIndex: number | null // the line's place among the MAC lines; null when it is not on the MAC claim
}

// One entry per MAC line, with each vendor's paired line beside it, then
// the vendor lines with no MAC partner. Vendor lines with the same REF*6R
// share an entry.
export function pairedLines<T extends { line_control_number: string }>(mac: T[] | undefined, vendors: Partial<Record<Slot, T[]>>): PairedLine<T>[] {
  const macLines = mac ?? []
  const pairing: Partial<Record<Slot, ReturnType<typeof pairLines<T>>>> = {}
  for (const slot of SIDES) {
    const lines = vendors[slot]
    if (lines) pairing[slot] = pairLines(macLines, lines)
  }
  const out: PairedLine<T>[] = macLines.map((m, i) => {
    const src: Sources<T> = { mac: m }
    for (const slot of SIDES) {
      const j = pairing[slot]?.partner[i] ?? -1
      if (j >= 0) src[slot] = vendors[slot]![j]
    }
    const c = m.line_control_number
    return { src, control: { v1: c, v2: c }, macIndex: i }
  })

  const unpaired: { src: Sources<T>; control: Partial<Record<Slot, string>> }[] = []
  for (const slot of SIDES) {
    for (const j of pairing[slot]?.unpaired ?? []) {
      const line = vendors[slot]![j]
      const c = line.line_control_number
      const same = c ? unpaired.find((u) => u.src[slot] === undefined && Object.values(u.control).includes(c)) : undefined
      if (same) {
        same.src[slot] = line
        same.control[slot] = c
      } else {
        unpaired.push({ src: { [slot]: line }, control: { [slot]: c } })
      }
    }
  }
  for (const u of unpaired) out.push({ src: u.src, control: u.control, macIndex: null })
  return out
}

// "Line 2 · 71046-26 · L001202", or "Line not on the MAC claim · …" for a
// vendor line with no MAC partner ("Line · …" when there is no MAC claim).
// `notOnFirst` is the words for a line the first column does not hold.
export function lineTitle<T extends { line_control_number: string }>(line: PairedLine<T>, hasMac: boolean, procedure: (t: T) => string, notOnFirst = 'Line not on the MAC claim'): string {
  const first = (line.src.mac ?? line.src.v1 ?? line.src.v2)!
  const c = first.line_control_number
  const tail = `${procedure(first)}${c ? ` · ${c}` : ''}`
  return line.macIndex !== null ? `Line ${line.macIndex + 1} · ${tail}` : `${hasMac ? notOnFirst : 'Line'} · ${tail}`
}
