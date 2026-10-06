package x12837p

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"

// OtherSubscriberSBR holds one occurrence of loop 2320, "Other Subscriber
// Information" - a coordination-of-benefits payer other than the one the
// claim is billed to - and everything nested under it, down through its
// 2330A-G NM1s. One ClaimCLM can carry several (ClaimCLM.OtherSubscribers).
// Citations below are to
// docs/reference/file-formats/cms-837-835-flatfiles-r907otn-2011.pdf
// ("Transmittal 907, PDF p.N") and
// docs/reference/file-formats/cedi-837p-companion-guide.pdf ("CEDI 837P CG,
// PDF p.N"), per the citation form documented at the top of
// frontend/src/data/x12Names837p.ts.
type OtherSubscriberSBR struct {
	Subscriber SubscriberSBR `json:"subscriber,omitempty"` // SBR; OTHER_SUBSCRIBER_SBR_EL_NAME; Transmittal 907, PDF p.191-192

	// CAS, "Claim Level Adjustments"; up to 5 occurrences. Must not be
	// present when loop 2000B's SBR01 = P (CEDI 837P CG, PDF p.40, printed
	// p.32; Table 28). CLAIM_LEVEL_ADJUSTMENTS_CAS_NAME; PDF p.192-194.
	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"`

	// AMT (loop 2320), three fixed boxes, not a qualifier-selected
	// catch-all. Payer Paid Amount (D) is required whenever loop 2000B's
	// SBR01 = S (CEDI 837P CG, PDF p.40, printed p.32; Table 29); each box's
	// AMT01 code is its only value in Transmittal 579's Values column (PDF
	// p.66). The parser fills a box from the first AMT with its code that
	// carries an amount, and keeps any other AMT in Raw. PDF p.194.
	PayerPaidAMTD            AmountAMT `json:"payerPaidAMTD,omitempty"`            // COB Payer Paid Amount; AMT01 D
	NonCoveredAMTA8          AmountAMT `json:"nonCoveredAMTA8,omitempty"`          // COB Total Non-Covered Amount; AMT01 A8
	RemainingLiabilityAMTEAF AmountAMT `json:"remainingLiabilityAMTEAF,omitempty"` // Remaining Patient Liability; AMT01 EAF

	OtherInsurance         OtherInsuranceOI          `json:"otherInsurance,omitempty"`         // OI, "Other Insurance Coverage Information" (required); PDF p.194-195
	OutpatientAdjudication OutpatientAdjudicationMOA `json:"outpatientAdjudication,omitempty"` // MOA, "Outpatient Adjudication (MOA) Information"; PDF p.195

	// 2330A-G NM1s. Each is its own X12 loop (its own Loop ID in the flat
	// file), so each carries its own Raw rather than bubbling into this
	// occurrence's Raw below.
	OtherSubscriberNM1IL               OtherSubscriberNM1 `json:"otherSubscriberNM1IL,omitempty"`               // 2330A; PDF p.195-197
	OtherPayerNM1PR                    OtherPayerNM1      `json:"otherPayerNM1PR,omitempty"`                    // 2330B; PDF p.197-200
	OtherPayerReferringProvider        OtherPayerRoleNM1  `json:"otherPayerReferringProvider,omitempty"`        // 2330C; not used by Medicare; PDF p.200-201
	OtherPayerRenderingProviderNM182   OtherPayerRoleNM1  `json:"otherPayerRenderingProviderNM182,omitempty"`   // 2330D; not used by Medicare; PDF p.201-202
	OtherPayerServiceFacilityNM177     OtherPayerRoleNM1  `json:"otherPayerServiceFacilityNM177,omitempty"`     // 2330E; not used by Medicare; PDF p.202-203
	OtherPayerSupervisingProviderNM1DQ OtherPayerRoleNM1  `json:"otherPayerSupervisingProviderNM1DQ,omitempty"` // 2330F; not used by Medicare; PDF p.203-204
	OtherPayerBillingProviderNM185     OtherPayerRoleNM1  `json:"otherPayerBillingProviderNM185,omitempty"`     // 2330G; not used by Medicare; PDF p.204-205

	Raw []RawSegment `json:"raw,omitempty"` // segments directly within loop 2320 the model does not place
}

// OtherInsuranceOI holds the OI segment (loop 2320, "Other Insurance
// Coverage Information", required whenever loop 2320 is sent). OI01, OI02,
// and OI05 are Not Used and kept as plain text. OI03/OI04/OI06 share the
// same three code lists CLM08/CLM10/CLM09 use at the claim level, so they
// reuse those enum types rather than inventing OI-specific ones. Set is
// false when the segment did not appear at all. Transmittal 907, PDF
// p.194-195; OI_NAME / OI_EL_NAME.
type OtherInsuranceOI struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ClaimFilingIndicatorOI01 string `json:"claimFilingIndicatorOI01,omitempty"` // Claim Filing Indicator Code (N/U) — Transmittal 907, PDF p.194
	ClaimSubmitReasonOI02    string `json:"claimSubmitReasonOI02,omitempty"`    // Claim Submission Reason Code (N/U) — Transmittal 907, PDF p.195

	BenefitsAssignmentOI03 BenefitsAssignmentCertificationCLM08 `json:"benefitsAssignmentOI03,omitempty"` // Benefits Assignment Certification Indicator
	SignatureSourceOI04    PatientSignatureSourceCLM10          `json:"signatureSourceOI04,omitempty"`    // Patient Signature Source Code

	ProviderAgreementOI05 string `json:"providerAgreementOI05,omitempty"` // Provider Agreement Code (N/U) — Transmittal 907, PDF p.195

	ReleaseOfInfoOI06 ReleaseOfInformationCLM09 `json:"releaseOfInfoOI06,omitempty"` // Release of Information Code

	Extra []string `json:"extra,omitempty"` // elements after OI06, kept as sent
}

// OutpatientAdjudicationMOA holds the MOA segment (loop 2320, "Outpatient
// Adjudication (MOA) Information") - the prior payer's own remittance
// summary for this claim, reported back on the 837P. Its five remark
// positions resolve as RARC, the same scheme x12835's OutpatientMOA uses
// for the 835 side's own MOA. Set is false when the segment did not appear
// at all. Transmittal 907, PDF p.195; MOA_NAME / MOA_EL_NAME.
type OutpatientAdjudicationMOA struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	ReimbursementRateMOA01  Amount    `json:"reimbursementRateMOA01,omitempty"`  // Reimbursement Rate
	HCPCSPayableAmountMOA02 Amount    `json:"hcpcsPayableAmountMOA02,omitempty"` // HCPCS Payable Amount
	RemarkMOA03             rarc.Code `json:"remarkMOA03,omitempty"`             // Claim Payment Remark Code
	RemarkMOA04             rarc.Code `json:"remarkMOA04,omitempty"`             // Claim Payment Remark Code
	RemarkMOA05             rarc.Code `json:"remarkMOA05,omitempty"`             // Claim Payment Remark Code
	RemarkMOA06             rarc.Code `json:"remarkMOA06,omitempty"`             // Claim Payment Remark Code
	RemarkMOA07             rarc.Code `json:"remarkMOA07,omitempty"`             // Claim Payment Remark Code
	ESRDPaymentMOA08        Amount    `json:"esrdPaymentMOA08,omitempty"`        // End Stage Renal Disease Payment Amount
	NonPayableProfCompMOA09 Amount    `json:"nonPayableProfCompMOA09,omitempty"` // Non-Payable Professional Component Billed Amount

	Extra []string `json:"extra,omitempty"` // elements after MOA09, kept as sent
}

// OtherSubscriberNM1 holds loop 2330A, "Other Subscriber Name": NM1, N3,
// N4, and a single REF occurrence (REF02 "Other Insured Additional
// Identifier", restricted to a 9-digit numeric value, CEDI 837P CG, PDF
// p.41, printed p.33, Table 31). Transmittal 907, PDF p.195-197.
type OtherSubscriberNM1 struct {
	Name        NameNM1      `json:"name,omitempty"`        // NM1; NM1_ROLE_EL_NAME["2330A"]
	Address     AddressN3    `json:"address,omitempty"`     // N3
	Location    LocationN4   `json:"location,omitempty"`    // N4
	SecondaryID ReferenceREF `json:"secondaryID,omitempty"` // REF, "Other Subscriber Secondary Identification"

	Raw []RawSegment `json:"raw,omitempty"`
}

// OtherPayerNM1 holds loop 2330B, "Other Payer Name": NM1, N3, N4, DTP, and
// five REF boxes (OTHER_PAYER_REF_NAME / OTHER_PAYER_REF02_EL_NAME);
// secondary identification can repeat up to twice, the other four REF
// boxes occur at most once. The parser places each REF by REF01 (codes from
// Transmittal 579, PDF p.71-73). It keeps every secondary identification
// REF, and keeps a repeat of a single box, a REF whose REF01 matches no box,
// or one with an empty REF02 for a single box, in Raw. Transmittal 907, PDF
// p.197-200.
type OtherPayerNM1 struct {
	Name                 NameNM1        `json:"name,omitempty"`                 // NM1; NM1_ROLE_EL_NAME["2330B"]
	Address              AddressN3      `json:"address,omitempty"`              // N3
	Location             LocationN4     `json:"location,omitempty"`             // N4
	RemittanceDateDTP573 DateDTP        `json:"remittanceDateDTP573,omitempty"` // DTP, "Date - Claim Check or Remittance Date", DTP01 573 (OTHER_PAYER_DTP_NAME / OTHER_PAYER_DTP03_EL_NAME, "Adjudication or Payment Date")
	SecondaryID          []ReferenceREF `json:"secondaryID,omitempty"`          // REF, "Other Payer Secondary Identification", REF01 2U, EI, FY, or NF; up to 2

	PriorAuthorizationREFG1       ReferenceREF `json:"priorAuthorizationREFG1,omitempty"`       // REF, "Other Payer Prior Authorization Number", REF01 G1
	ReferralREF9F                 ReferenceREF `json:"referralREF9F,omitempty"`                 // REF, "Other Payer Referral Number", REF01 9F
	ClaimAdjustmentIndicatorREFT4 ReferenceREF `json:"claimAdjustmentIndicatorREFT4,omitempty"` // REF, "Other Payer Claim Adjustment Indicator", REF01 T4
	ClaimControlNumberREFF8       ReferenceREF `json:"claimControlNumberREFF8,omitempty"`       // REF, "Other Payer Claim Control Number", REF01 F8

	Raw []RawSegment `json:"raw,omitempty"`
}

// OtherPayerRoleNM1 is the loop shape shared by 2330C-G (Other Payer
// Referring/Rendering/Service Facility/Supervising/Billing Provider): NM1
// plus a repeating REF of that role's own secondary identification (up to
// 3 for 2330C-F, up to 2 for 2330G, per OTHER_PAYER_ROLE_REF_NAME). None of
// these loops is used by Medicare (NOT_USED_BY_MEDICARE's 2330C-G entries);
// the field holding this struct on OtherSubscriberSBR carries which role
// and loop it is. Transmittal 907, PDF p.200-205.
type OtherPayerRoleNM1 struct {
	Name       NameNM1        `json:"name,omitempty"`       // NM1
	References []ReferenceREF `json:"references,omitempty"` // REF, "Other Payer ... Secondary Identification"

	Raw []RawSegment `json:"raw,omitempty"`
}
