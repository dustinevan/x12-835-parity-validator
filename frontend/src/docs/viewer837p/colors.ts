// The 837P viewer's box tints, matching the 835's own rule exactly
// (frontend/src/docs/viewer835/colors.ts and the style the 835's Viewer835.tsx
// passes to the shared BoxView): boxes are marked, not filled. Every box
// keeps the page background — no tint on a box's body — and a tinted kind
// carries only a left accent in the site's blue (series-1) at one constant
// strength, plus a faint ink tint on its opening row alone. Depth (the
// indent) is what tells the nested families apart, not color; amber
// (series-2) and red/pink (series-3) are never used for a box, since both
// read as warnings on this page (the same reason the 835 side gives,
// colors.ts there). The envelope (ISA/GS/ST) and the two HL wrapper levels
// (billing provider, subscriber) stay neutral, exactly as the 835 keeps its
// envelope and its LX provider-group wrapper neutral — a structural
// wrapper carries no accent of its own, only the loops nested inside it
// do. "From the 835" green has no meaning on this page — an 837P carries
// no value an 835 later confirms — so it is never used here.
import { LOOP_RULES, UNGOVERNED_TINT, type BoxType, type Tint } from './rules'

// Every box type's tint, combining the loop-rule table (every box that
// opens on a segment+qualifier match) with the six positionally-opened
// types rules.ts names directly. Tint keeps its four marked levels
// (rules.ts) as a record of which loop family a box belongs to, even
// though every marked level renders identically below — useful to a
// reader tracing the grammar, and a seam if a future design ever wants to
// tell the families apart some other way than depth.
const TINT_BY_TYPE: Record<BoxType, Tint> = { ...UNGOVERNED_TINT } as Record<BoxType, Tint>
for (const rule of LOOP_RULES) TINT_BY_TYPE[rule.boxType] = rule.tint

export function tintOf(boxType: BoxType): Tint {
  return TINT_BY_TYPE[boxType]
}

export function tintBoxClass(tint: Tint): string {
  if (tint === 'neutral') return 'border border-line'
  return 'border border-line border-l-2 border-l-series-1'
}

// The opening row of a tinted box is the one place its tint appears; the
// box's body stays on the page background so nested boxes never stack
// tints (same rule as the 835's openingClass).
export function tintOpeningClass(tint: Tint): string {
  return tint === 'neutral' ? '' : 'bg-ink/4'
}

// The two functions frontend/src/docs/viewer/BoxView.tsx's BoxStyle wants —
// see Viewer837P.tsx, which passes these straight through.
export function boxClass(boxType: BoxType): string {
  return tintBoxClass(tintOf(boxType))
}
export function openingClass(boxType: BoxType): string {
  return tintOpeningClass(tintOf(boxType))
}
