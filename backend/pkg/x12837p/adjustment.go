package x12837p

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"

// AdjustmentCAS holds one CAS segment: a group code followed by up to six
// reason/amount/quantity triples (CAS02-04, CAS05-07, ... CAS17-19).
// CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME (loop 2320, Other Subscriber claim-level
// adjustments) and LINE_ADJUSTMENT_CAS_NAME (loop 2430, Line Adjudication
// line-level adjustments) in frontend/src/data/x12Names837p.ts name the
// segment differently by loop; the field holding this struct carries that
// distinction, not a field here. A CAS segment must not be present in loop
// 2320 when 2000B SBR01 = P (CEDI 837P CG, PDF p.40, printed p.32; Table 28).
//
// ClaimAdjustmentGroupCodeCAS01 is defined in enums_cob.go. Reason reuses
// carc.Code (backend/pkg/carc), the same CARC list the 835 side of this
// project already embeds — frontend/src/data/x12Names837p.ts's own
// "CARC/RARC reuse note" says the 2320/2430 CAS segments use that same
// externally maintained reason-code set, not a separately sourced one.
type AdjustmentCAS struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	GroupCAS01 ClaimAdjustmentGroupCodeCAS01 `json:"groupCAS01,omitempty"` // Claim Adjustment Group Code — Transmittal 907, every CAS instance; CAS01_EL_NAME
	Reasons    []CASReason                   `json:"reasons,omitempty"`    // CAS02-19, up to six triples; an empty triple between two sent ones keeps its place

	Extra []string `json:"extra,omitempty"` // elements after CAS19, kept as sent
}

// CASReason is one reason/amount/quantity triple within a CAS segment.
type CASReason struct {
	Reason   carc.Code `json:"reason,omitempty"`   // Adjustment Reason Code — Transmittal 907, every CAS instance; CAS_REASON_EL_NAME
	Amount   Amount    `json:"amount,omitempty"`   // Adjustment Amount — CAS_AMOUNT_EL_NAME
	Quantity Amount    `json:"quantity,omitempty"` // Adjustment Quantity — CAS_QTY_EL_NAME
}
