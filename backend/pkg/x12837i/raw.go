package x12837i

// RawSegment preserves one segment that the model does not place into a
// typed field, keeping its position and raw elements so no data is
// silently dropped. It is used both for a segment ID the model does not
// model at all, and for a segment the model has no slot for: a repeat of a
// segment it holds once, or a qualifier it does not hold at that level. A
// raw segment's elements are not checked. Copied in shape from
// backend/pkg/x12835/raw.go.
type RawSegment struct {
	Index    int      `json:"index,omitempty"`    // position of this segment in the file, 0-based; the ISA is 0.
	ID       string   `json:"id,omitempty"`       // the segment ID, e.g. "K3".
	Elements []string `json:"elements,omitempty"` // elements[0] is element 01, elements[1] is element 02, and so on.
}
