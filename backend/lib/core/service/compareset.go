package service

import (
	"fmt"
	"strings"
)

// CompareSet names what the comparison leaves out of a claim. Everything
// inside a claim (CLP through its last service line) is compared except
// the entries listed here; nothing outside the claim (the envelope, the
// payment header, the N1 loops, LX, TS3, TS2, PLB) is compared under any
// set. Under every set, a compared dollar amount follows the five-cent rule
// (AmountMatchCents in compare.go) and every other value compares exactly.
// The set is applied while comparing: a left-out element is never
// written as a difference, and a pair whose only differences are left out
// is equal. Because the claim bytes are stored, rerunning compare under
// another set regenerates every left-out difference. The submission records
// the set's name (submission.compared_set), so every row can be traced to
// the set that produced it.
type CompareSet struct {
	Name     string
	LeaveOut []LeaveOut
}

// LeaveOut is one left-out element: the level it sits at, its tag as the
// comparison labels it, and why, with the source. A tag ending in * covers
// every qualifier: "QTY*" is every QTY, whatever its QTY01.
type LeaveOut struct {
	Level  string // LevelClaim or LevelLine
	Tag    string
	Reason string
}

// The two levels a difference is found at.
const (
	LevelClaim = "claim"
	LevelLine  = "line"
)

// leavesOut reports whether the set leaves out the tag at the level.
func (s CompareSet) leavesOut(level, tag string) bool {
	for _, e := range s.LeaveOut {
		if e.Level != level {
			continue
		}
		if e.Tag == tag || (strings.HasSuffix(e.Tag, "*") && strings.HasPrefix(tag, e.Tag)) {
			return true
		}
	}
	return false
}

// compareSet1 is set "1" (docs/design/compared-fields.md):
// a wide net. Every element inside a claim is compared except the match-key
// fields where they appear as values, the system-assigned claim control
// numbers, the two CLP elements Medicare does not use, the claim-level
// contact, the line control number that pairs the lines, and service-line
// quantities. CAS adjustment quantities are compared, and unplaced segments
// are compared whole. One leave-out row of compared-fields.md has no tag here
// because the comparison never reads it: an NM1 is compared by its NM109
// under its role, and the names and qualifiers (NM102 to NM108) are not
// read (compared-fields.md, identity and pairing: names differ in spelling
// and length).
var compareSet1 = CompareSet{
	Name: "1",
	LeaveOut: []LeaveOut{
		{LevelClaim, "CLP03", "Total Claim Charge Amount: a match-key field, equal as an exact decimal within a pair (compared-fields.md, identity and pairing)"},
		{LevelClaim, "NM1*QC", "Patient Identifier: a match-key field, equal within a pair apart from the case and spacing the key ignores (compared-fields.md, identity and pairing)"},
		{LevelClaim, "DTM*232", "Statement From Date: the key's date of service, equal whenever both sides send it (compared-fields.md, identity and pairing)"},
		{LevelClaim, "CLP07", "Payer Claim Control Number: each system assigns its own (compared-fields.md, system-assigned)"},
		{LevelClaim, "REF*F8", "Other Claim Related ID: carries the original claim's control number (compared-fields.md, system-assigned; taxonomy T2.40)"},
		{LevelClaim, "CLP10", "Patient Status Code: marked not used by Medicare (CMS 835 Flat File p.6)"},
		{LevelClaim, "CLP14", "Yes/No Condition or Response Code: marked not used by Medicare (CMS 835 Flat File p.6)"},
		{LevelClaim, "PER", "Claim Contact Information: the payer's contact for the claim (compared-fields.md, system-assigned; CMS 835 Flat File p.11)"},
		{LevelLine, "REF*6R", "Line Item Control Number: pairs the lines, and every line difference already records it (compared-fields.md, identity and pairing)"},
		{LevelLine, "QTY*", "Service Supplemental Quantity: marked not used by Medicare (CMS 835 Flat File p.13); claim-level QTY stays compared"},
	},
}

// compareSets is every set by name.
var compareSets = map[string]CompareSet{
	compareSet1.Name: compareSet1,
}

// LookupCompareSet returns the compared set with the name
// (APP_SERVER_COMPARE_SET). An unknown name is an error, so a misspelled
// setting fails at startup.
func LookupCompareSet(name string) (CompareSet, error) {
	s, ok := compareSets[name]
	if !ok {
		return CompareSet{}, fmt.Errorf("service: looking up compared set %q: unknown name", name)
	}
	return s, nil
}
