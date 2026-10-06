package x12837p

// 2320/2330 Other Subscriber and Payer (Coordination of Benefits) coded
// elements.

// ClaimAdjustmentGroupCodeCAS01 is CAS01, the Claim Adjustment Group Code,
// shared by every CAS segment in the transaction (2320 Other Subscriber
// claim-level adjustments and 2430 Line Adjudication line-level adjustments).
//
// The list is the five codes the 837P implementation guide allows at CAS01,
// the same externally maintained X12 list the 835 side of this project
// already uses (x12.org/codes/claim-adjustment-group-codes);
// frontend/src/data/x12Names837p.ts CAS_GROUP_NAME. A CAS segment must not be
// present in loop 2320 when 2000B SBR01 = P (CEDI 837P CG, PDF p.40, printed
// p.32; Table 28).
type ClaimAdjustmentGroupCodeCAS01 string

const (
	ClaimAdjustmentGroupCodeCAS01NotSent ClaimAdjustmentGroupCodeCAS01 = ""
	ClaimAdjustmentGroupCodeCAS01CO      ClaimAdjustmentGroupCodeCAS01 = "CO"
	ClaimAdjustmentGroupCodeCAS01OA      ClaimAdjustmentGroupCodeCAS01 = "OA"
	ClaimAdjustmentGroupCodeCAS01PR      ClaimAdjustmentGroupCodeCAS01 = "PR"
	ClaimAdjustmentGroupCodeCAS01PI      ClaimAdjustmentGroupCodeCAS01 = "PI"
	ClaimAdjustmentGroupCodeCAS01CR      ClaimAdjustmentGroupCodeCAS01 = "CR"
)

var claimAdjustmentGroupCodeCAS01List = codeList{
	{"CO", "Contractual Obligation", "Contractual Obligation — x12.org/codes/claim-adjustment-group-codes; frontend/src/data/x12Names837p.ts CAS_GROUP_NAME"},
	{"OA", "Other Adjustment", "Other Adjustment — x12.org/codes/claim-adjustment-group-codes; frontend/src/data/x12Names837p.ts CAS_GROUP_NAME"},
	{"PR", "Patient Responsibility", "Patient Responsibility — x12.org/codes/claim-adjustment-group-codes; frontend/src/data/x12Names837p.ts CAS_GROUP_NAME"},
	{"PI", "Payer Initiated Reductions", "Payer Initiated Reductions — x12.org/codes/claim-adjustment-group-codes; frontend/src/data/x12Names837p.ts CAS_GROUP_NAME"},
	{"CR", "Correction and Reversals", "Correction and Reversals — x12.org/codes/claim-adjustment-group-codes; frontend/src/data/x12Names837p.ts CAS_GROUP_NAME"},
}

func (v ClaimAdjustmentGroupCodeCAS01) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v ClaimAdjustmentGroupCodeCAS01) Known() bool {
	return claimAdjustmentGroupCodeCAS01List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v ClaimAdjustmentGroupCodeCAS01) ShortDescription() string {
	return claimAdjustmentGroupCodeCAS01List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v ClaimAdjustmentGroupCodeCAS01) Description() string {
	return claimAdjustmentGroupCodeCAS01List.description(string(v))
}

// ParseClaimAdjustmentGroupCodeCAS01 returns text as a
// ClaimAdjustmentGroupCodeCAS01, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseClaimAdjustmentGroupCodeCAS01(text string) (ClaimAdjustmentGroupCodeCAS01, error) {
	return ClaimAdjustmentGroupCodeCAS01(text), claimAdjustmentGroupCodeCAS01List.check("CAS01", text)
}
