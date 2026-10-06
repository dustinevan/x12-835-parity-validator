package x12835

// ProviderAdjustmentPLB holds one PLB segment: a provider and a fiscal
// period, followed by up to six adjustment pairs (PLB03-04 through
// PLB13-14).
type ProviderAdjustmentPLB struct {
	Index             int             `json:"index,omitempty"`
	ProviderPLB01     string          `json:"providerPLB01,omitempty"`     // Provider Identifier; CMS 835 Flat File, PDF p.13
	FiscalPeriodPLB02 Date            `json:"fiscalPeriodPLB02,omitempty"` // Fiscal Period Date; CMS 835 Flat File, PDF p.13
	Adjustments       []PLBAdjustment `json:"adjustments,omitempty"`       // PLB03-04, PLB05-06, ... PLB13-14; up to six, in position: an empty pair between two sent ones is kept
}

// PLBAdjustment is one adjustment-identifier/amount pair within a PLB
// segment (e.g. PLB03-04). The identifier is a composite: its first
// component is the Provider Adjustment Reason Code (PLB_REASON_NAME in
// frontend/src/data/x12Names835.ts, and X12's public
// provider-adjustment-reason-codes list for the codes that table leaves
// unresolved); its second component is free text.
type PLBAdjustment struct {
	Reason     PLBReasonPLB03 `json:"reason,omitempty"`     // Adjustment Identifier, first component; CMS 835 Flat File, PDF p.13-14
	Identifier string         `json:"identifier,omitempty"` // Adjustment Identifier, second component; CMS 835 Flat File, PDF p.13-14
	Amount     Amount         `json:"amount,omitempty"`     // Provider Adjustment Amount; CMS 835 Flat File, PDF p.13-14
}
