// Package claimrepeat reads and deletes the claim_repeat table: one row per
// match key that occurs more than once in a submission, with how many
// claims carry it. The load stage writes it from the submission's claim
// rows, so the INSERT ... SELECT is in lib/core/query (LoadQuery). The
// compare stage reads it inside its own statements (CompareQuery): a key
// with no row on either side is compared as a pair, a key with a row on the
// MAC side is set aside, and a key with a row on the vendor side only is
// duplicated.
package claimrepeat

// ClaimRepeat is one row. Claims is always more than one.
type ClaimRepeat struct {
	SubmissionPrefix string
	MatchKey         []byte
	Claims           int
}
