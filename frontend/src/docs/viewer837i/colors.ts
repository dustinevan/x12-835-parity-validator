// The CAS adjustment group (CO/OA/PR) marker, shared between a CAS row's own
// marker dot and the matching reason it annotates — same monochrome scheme
// and reasoning as frontend/src/docs/viewer835/colors.ts (read that file
// first; this one copies its shape for the 837I's own CAS segments at loop
// 2320, Other Subscriber Information, and loop 2430, Line Adjudication
// Information — coordination-of-benefits adjustments the file is reporting
// about another payer, not this transaction's own adjudication). Only
// CO appears in this page's examples; OA and PR are included for parity
// with the 835 and for when an example needs them. CR and PI (the other two
// codes in frontend/src/data/x12Names837i.ts's CAS_GROUP_NAME) are not yet
// exercised by any example — add them the same way (one more ink-weight
// step, or a second pattern class) if one needs them.
export const CAS_GROUP_DOT: Record<string, string> = {
  CO: 'bg-ink/55',
  OA: 'cas-pattern-oa',
  PR: 'bg-ink/35',
}
