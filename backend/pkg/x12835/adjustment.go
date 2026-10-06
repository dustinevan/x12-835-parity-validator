package x12835

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"

// AdjustmentCAS holds one CAS segment: a group code followed by up to six
// reason/amount/quantity triples (CAS02-04, CAS05-07, ... CAS17-19). The
// segment is named "Claim Adjustment" at the claim level (loop 2100) and
// "Service Adjustment" at the service line level (loop 2110) — see
// CLAIM_CAS_NAME and SERVICE_CAS_NAME in frontend/src/data/x12Names835.ts;
// the field holding this struct (claim-level or service-line Adjustments)
// carries that distinction, not a field here.
type AdjustmentCAS struct {
	Index      int           `json:"index,omitempty"`
	GroupCAS01 AdjGroupCAS01 `json:"groupCAS01,omitempty"` // Claim Adjustment Group Code; CMS 835 Flat File, PDF p.6 (claim loop), p.12 (service loop)
	Reasons    []CASReason   `json:"reasons,omitempty"`    // CAS02-19, up to six triples, in position: an empty triple between two sent ones is kept
}

// CASReason is one reason/amount/quantity triple within a CAS segment.
type CASReason struct {
	Reason   carc.Code `json:"reason,omitempty"`   // Adjustment Reason Code; CMS 835 Flat File, PDF p.6 (claim loop), p.12-13 (service loop)
	Amount   Amount    `json:"amount,omitempty"`   // Adjustment Amount; same pages
	Quantity Amount    `json:"quantity,omitempty"` // Adjustment Quantity; same pages
}
