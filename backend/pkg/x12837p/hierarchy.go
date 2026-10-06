package x12837p

// HierarchyHL holds one HL segment's own four elements. HL_EL_NAME in
// frontend/src/data/x12Names837p.ts gives the generic position names,
// shared by every HL instance (2000A, 2000B, 2000C; source: Transmittal
// 907, PDF p.136, 143, 149). The parser places each HL loop by its LevelHL03
// code and attaches it to its parent by ParentIDNumberHL02 (see Parse), so
// the loop struct embedding this value (BillingProviderHL, SubscriberHL,
// PatientHL) always agrees with LevelHL03.
//
// HL01 and HL02 are AN 1-12 in both flat files (Transmittal 907, PDF p.136;
// Transmittal 579, PDF p.9), so they are kept as text, even when all digits.
type HierarchyHL struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	IDNumberHL01       string                    `json:"idNumberHL01,omitempty"`       // Hierarchical ID Number — unique within the transaction
	ParentIDNumberHL02 string                    `json:"parentIDNumberHL02,omitempty"` // Hierarchical Parent ID Number — Not Used, so empty, at the top level (2000A)
	LevelHL03          HierarchicalLevelHL03     `json:"levelHL03,omitempty"`          // Hierarchical Level Code
	ChildCodeHL04      HierarchicalChildCodeHL04 `json:"childCodeHL04,omitempty"`      // Hierarchical Child Code

	Extra []string `json:"extra,omitempty"` // elements after HL04, kept as sent
}

// HierarchicalLevelHL03 is HL03, the Hierarchical Level Code.
//
// The 837P implementation guide allows 20, 22, and 23. Transmittal 907 names
// each HL loop's own box (LOOP_HL_NAME in frontend/src/data/x12Names837p.ts:
// "Billing Provider Hierarchical Level", "Subscriber Hierarchical Level",
// "Patient Hierarchical Level") but does not restate HL03's code characters.
// Transmittal 579, the original 005010 837P flat file, prints them in its
// Values column: 20 at 2000A, 22 at 2000B, and 23 at 2000C. The CEDI
// companion guide's worked transaction shows 20 and 22 ("HL*1**20*1",
// "HL*2*1*22*0"). The names are the guide's.
type HierarchicalLevelHL03 string

const (
	HierarchicalLevelHL03NotSent HierarchicalLevelHL03 = ""
	HierarchicalLevelHL0320      HierarchicalLevelHL03 = "20"
	HierarchicalLevelHL0322      HierarchicalLevelHL03 = "22"
	HierarchicalLevelHL0323      HierarchicalLevelHL03 = "23"
)

var hierarchicalLevelHL03List = codeList{
	{"20", "Information Source", `Information Source — the 2000A Billing Provider Hierarchical Level's HL03 value, Transmittal 579, PDF p.9 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*1**20*1"); the name is the guide's: ` + tr3},
	{"22", "Subscriber", `Subscriber — the 2000B Subscriber Hierarchical Level's HL03 value, Transmittal 579, PDF p.16 (Values column); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*2*1*22*0"); the name is the guide's: ` + tr3},
	{"23", "Dependent", "Dependent — the 2000C Patient Hierarchical Level's HL03 value, Transmittal 579, PDF p.22 (Values column); the loop is not used by Medicare; the name is the guide's: " + tr3},
}

func (v HierarchicalLevelHL03) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v HierarchicalLevelHL03) Known() bool { return hierarchicalLevelHL03List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v HierarchicalLevelHL03) ShortDescription() string {
	return hierarchicalLevelHL03List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v HierarchicalLevelHL03) Description() string {
	return hierarchicalLevelHL03List.description(string(v))
}

// ParseHierarchicalLevelHL03 returns text as a HierarchicalLevelHL03, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseHierarchicalLevelHL03(text string) (HierarchicalLevelHL03, error) {
	return HierarchicalLevelHL03(text), hierarchicalLevelHL03List.check("HL03", text)
}

// HierarchicalChildCodeHL04 is HL04, the Hierarchical Child Code (whether
// this HL has a subordinate HL beneath it).
//
// The 837P implementation guide allows 0 and 1. Transmittal 579's Values
// column lists 1 at 2000A, "0, 1" at 2000B, and 0 at 2000C (PDF p.9, 16, 22);
// the CEDI worked transaction shows both. The names are the guide's.
type HierarchicalChildCodeHL04 string

const (
	HierarchicalChildCodeHL04NotSent HierarchicalChildCodeHL04 = ""
	HierarchicalChildCodeHL040       HierarchicalChildCodeHL04 = "0"
	HierarchicalChildCodeHL041       HierarchicalChildCodeHL04 = "1"
)

var hierarchicalChildCodeHL04List = codeList{
	{"0", "No Subordinate HL Segment in This Hierarchical Structure", `No Subordinate HL Segment in This Hierarchical Structure — Transmittal 579, PDF p.16 and p.22 (Values column, 2000B and 2000C); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*2*1*22*0"); the name is the guide's: ` + tr3},
	{"1", "Additional Subordinate HL Data Segment in This Hierarchical Structure", `Additional Subordinate HL Data Segment in This Hierarchical Structure — Transmittal 579, PDF p.9 and p.16 (Values column, 2000A and 2000B); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, "HL*1**20*1"); the name is the guide's: ` + tr3},
}

func (v HierarchicalChildCodeHL04) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v HierarchicalChildCodeHL04) Known() bool {
	return hierarchicalChildCodeHL04List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v HierarchicalChildCodeHL04) ShortDescription() string {
	return hierarchicalChildCodeHL04List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v HierarchicalChildCodeHL04) Description() string {
	return hierarchicalChildCodeHL04List.description(string(v))
}

// ParseHierarchicalChildCodeHL04 returns text as a
// HierarchicalChildCodeHL04, with an error when the text is neither empty
// nor on the list. The value is kept either way.
func ParseHierarchicalChildCodeHL04(text string) (HierarchicalChildCodeHL04, error) {
	return HierarchicalChildCodeHL04(text), hierarchicalChildCodeHL04List.check("HL04", text)
}
