package x12837i

// HierarchyHL is the HL segment: the same shape at every hierarchical level
// in this transaction (2000A billing provider, 2000B subscriber, 2000C
// patient). Element names and citations are
// frontend/src/data/x12Names837i.ts's HL_EL_NAME array (source: Transmittal
// 907, PDF p.10, representative of every HL instance). The loop-specific
// structural types that embed this shape (BillingProviderHL, SubscriberHL,
// PatientHL, defined elsewhere in this package) carry the loop's actual
// business content; this struct only holds the HL segment's own four
// elements.
//
// IDNumberHL01 and ParentIDNumberHL02 are kept as strings, like every other
// identifier in this package, even though they are sequential integers in
// practice. ChildCodeHL04 is kept as a string ("1" more subordinate HL
// loops follow this one, "0" none do): no coded type for it is in this
// package's catalog (frontend/src/data/x12Names837i.ts's own list of
// implied coded types has none for HL04).
type HierarchyHL struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	IDNumberHL01       string                `json:"idNumberHL01,omitempty"`       // Hierarchical ID Number — Transmittal 907, PDF p.10
	ParentIDNumberHL02 string                `json:"parentIDNumberHL02,omitempty"` // Hierarchical Parent ID Number — Transmittal 907, PDF p.10 (empty at the top level, 2000A)
	LevelCodeHL03      HierarchicalLevelHL03 `json:"levelCodeHL03,omitempty"`      // Hierarchical Level Code — Transmittal 907, PDF p.10
	ChildCodeHL04      string                `json:"childCodeHL04,omitempty"`      // Hierarchical Child Code — Transmittal 907, PDF p.10
}
