package x12837i

// ReceiverNM140 is the 1000B loop: the receiver's NM1 segment (NM101 fixed
// at "40"). Transmittal 907 carries no PER or other segment at this loop.
// The NM1 uses the shared NameNM1 shape. Its positions' public names at
// this loop are frontend/src/data/x12Names837i.ts's NM1_RECEIVER_EL_NAME
// array (source: Transmittal 907, PDF p.10): NM103 Receiver Name, NM104-107
// Not Used, NM108 Identification Code Qualifier, NM109 Receiver Primary
// Identifier.
type ReceiverNM140 struct {
	Name NameNM1 `json:"name,omitempty"` // NM1*40 — Transmittal 907, PDF p.10

	Raw []RawSegment `json:"raw,omitempty"` // segments within the loop, after the NM1, that the model does not place
}
