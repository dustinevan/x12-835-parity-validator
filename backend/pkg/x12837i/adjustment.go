package x12837i

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"

// AdjustmentCAS holds one CAS segment: a group code followed by up to six
// reason/amount/quantity triples (CAS02-04, CAS05-07, ... CAS17-19). The
// segment is named "Claim Level Adjustments" at loop 2320 and "Line
// Adjustment" at loop 2430 (frontend/src/data/x12Names837i.ts
// CAS_CLAIM_NAME / CAS_LINE_NAME); the field holding this struct carries
// that distinction, not a field here. Copied in shape from
// backend/pkg/x12835/adjustment.go, reusing the same CARC code package: the
// Adjustment Reason Code list is shared across transaction types, not
// specific to the 835.
type AdjustmentCAS struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	GroupCAS01 CASGroupCode01 `json:"groupCAS01,omitempty"` // Claim Adjustment Group Code — Transmittal 907, PDF p.98 (2320), p.122 (2430)
	Reasons    []CASReason    `json:"reasons,omitempty"`    // CAS02-19, up to six triples, in position: an empty triple between two sent ones stays empty
}

// CASReason is one reason/amount/quantity triple within a CAS segment.
type CASReason struct {
	Reason   carc.Code `json:"reason,omitempty"`   // Adjustment Reason Code — Transmittal 907, PDF p.98 (2320), p.122 (2430)
	Amount   Amount    `json:"amount,omitempty"`   // Adjustment Amount — same pages
	Quantity Amount    `json:"quantity,omitempty"` // Adjustment Quantity — same pages
}
