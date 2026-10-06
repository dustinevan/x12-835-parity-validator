// The CAS adjustment group (CO/OA/PR) marker, shared between a CAS row's own
// marker dot and that group's slice in the maximized line's bar (Bar.tsx),
// so the two visibly match — docs/design/835-viewer.md "Color". A
// fixed, non-judgment assignment: CO/OA/PR are neutral X12 category codes,
// not a quality verdict. Monochrome, no hue: green is reserved elsewhere in
// this project for an 835 element that carries the value submitted on the
// 837, and amber already means "difference" elsewhere in this project —
// reusing either here would misread as one of those meanings. Each group is
// the page's own ink (`bg-ink`) at a different weight instead: CO at 55%
// opacity, PR at 35%, and OA a diagonal-stripe fill (full-strength ink, no
// opacity step, so three groups stay distinguishable at a glance without a
// fourth opacity level). Kept separate from the three structural hues
// (series-1/2/3 — see BoxView.tsx) so a CAS marker never reads as "part of"
// its enclosing claim or service line box's own tint.
// CO and PR are literal Tailwind opacity-modifier classes (not built from
// interpolated strings) so the build's static class scan always picks them
// up. OA is a small fixed CSS class (`.cas-pattern-oa`, frontend/src/index.css)
// carrying a repeating-linear-gradient in the page's ink color, since
// Tailwind has no utility for a diagonal-stripe fill; it is defined once in
// both light and dark mode so it stays subtle in each.
export const CAS_GROUP_DOT: Record<string, string> = {
  CO: 'bg-ink/55',
  OA: 'cas-pattern-oa',
  PR: 'bg-ink/35',
}
export const CAS_GROUP_BAR: Record<string, string> = {
  CO: 'bg-ink/55',
  OA: 'cas-pattern-oa',
  PR: 'bg-ink/35',
}
