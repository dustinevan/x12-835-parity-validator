package x12837p

// ReceiverNM140 is the 1000B Receiver Name loop: the NM1*40 segment alone.
// Source: Transmittal 907, PDF p.136 (footer 4); LOOP_NM1_NAME["1000B"] and
// NM1_ROLE_EL_NAME["1000B"] in frontend/src/data/x12Names837p.ts, which
// gives no PER table for this loop. NM101 is EntityIdentifierNM10140 on a
// compliant file.
type ReceiverNM140 struct {
	Name NameNM1 `json:"name,omitempty"` // NM1*40 — Transmittal 907, PDF p.136 (footer 4)

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop that the model does not place
}
