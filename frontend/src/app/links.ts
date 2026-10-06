import { withParams } from '../format'

// Where one claim opens: the claim page, with the hex match key and the
// selection's v1, v2, and side (Selection.carry), so the claim page shows
// the MAC claim beside both picked submissions.
export function claimLink(selection: Record<string, string | undefined>, matchKey: string) {
  return withParams('/claim', { ...selection, match_key: matchKey })
}
