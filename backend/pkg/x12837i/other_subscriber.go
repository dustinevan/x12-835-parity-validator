package x12837i

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"

// OtherSubscriberSBR holds the 2320 Other Subscriber Information loop and
// everything nested under it (2330A-I): one other payer's coordination-of-
// benefits information for this claim. SubscriberSBR is the shared SBR
// struct, defined for the primary subscriber loop (2000B) and reused here
// unchanged — SBR_OTHER_EL_NAME in frontend/src/data/x12Names837i.ts notes
// the identical shape (Transmittal 907, PDF p.97).
type OtherSubscriberSBR struct {
	SBR SubscriberSBR `json:"sbr,omitempty"`

	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"` // CAS, "Claim Level Adjustments" at 2320; PDF p.98 (CAS_CLAIM_NAME)

	MOA OutpatientMOA `json:"moa,omitempty"` // MOA, present only for an outpatient claim; PDF p.102
	MIA InpatientMIA  `json:"mia,omitempty"` // MIA, present only for an inpatient claim; PDF p.100-101

	OtherInsurance OtherHealthInsuranceOI `json:"otherInsurance,omitempty"` // OI; PDF p.100

	// AMT (2320). AMT01_NAME in frontend/src/data/x12Names837i.ts.
	OtherPayerPaidAMTD       AmountAMT   `json:"otherPayerPaidAMTD,omitempty"`       // Other Payer Paid Amount, COB Payer Paid Amount; AMT*D; PDF p.99
	RemainingLiabilityAMTEAF AmountAMT   `json:"remainingLiabilityAMTEAF,omitempty"` // Remaining Patient Liability; AMT*EAF; PDF p.99
	NonCoveredAMTA8          AmountAMT   `json:"nonCoveredAMTA8,omitempty"`          // COB Total Non-Covered Amount; AMT*A8; PDF p.100
	Amounts                  []AmountAMT `json:"amounts,omitempty"`                  // any other 2320 AMT, or a repeat of the above

	OtherInsuredNM1IL OtherInsuredNM1IL `json:"otherInsuredNM1IL,omitempty"` // 2330A, Other Subscriber Name; PDF p.102-103
	OtherPayerNM1PR   OtherPayerNM1PR   `json:"otherPayerNM1PR,omitempty"`   // 2330B, Other Payer Name; PDF p.104-106

	// Other payer provider loops (2330C-I). OTHER_PAYER_PROVIDER_LOOP_NAME
	// in frontend/src/data/x12Names837i.ts. 2330D-H reuse the same types
	// claim.go defines for the claim-level provider loops (2310B-F), since
	// Transmittal 907 gives them the identical NM1+REF shape at both
	// loops. 2330C (Attending) and 2330I (Billing) get their own types
	// below: 2330C omits the PRV that 2310A's AttendingNM171 carries, and
	// 2330I has no 2310 counterpart at all (the claim's own billing
	// provider is loop 2010AA, BillingProviderHL, with a different shape).
	OtherPayerAttendingNM171       OtherPayerAttendingNM171 `json:"otherPayerAttendingNM171,omitempty"`       // 2330C, role 71; PDF p.106
	OtherPayerOperatingNM172       OperatingNM172           `json:"otherPayerOperatingNM172,omitempty"`       // 2330D, role 72; PDF p.107
	OtherPayerOtherOperatingNM1ZZ  OtherOperatingNM1ZZ      `json:"otherPayerOtherOperatingNM1ZZ,omitempty"`  // 2330E, role ZZ; PDF p.108
	OtherPayerServiceFacilityNM177 ServiceFacilityNM177     `json:"otherPayerServiceFacilityNM177,omitempty"` // 2330F, role 77; PDF p.108-109
	OtherPayerRenderingNM182       RenderingNM182           `json:"otherPayerRenderingNM182,omitempty"`       // 2330G, role 82; PDF p.109
	OtherPayerReferringNM1DN       ReferringNM1DN           `json:"otherPayerReferringNM1DN,omitempty"`       // 2330H, role DN; PDF p.110
	OtherPayerBillingNM185         OtherPayerBillingNM185   `json:"otherPayerBillingNM185,omitempty"`         // 2330I, role 85; PDF p.111

	Raw []RawSegment `json:"raw,omitempty"` // segments within this other-subscriber loop the model does not place
}

// OtherInsuredNM1IL is loop 2330A, Other Subscriber Name: NM1 (role IL)
// and any REF (qualifier SY, Other Subscriber Secondary Identification).
// Transmittal 907, PDF p.102-103.
type OtherInsuredNM1IL struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// OtherPayerNM1PR is loop 2330B, Other Payer Name: NM1 (role PR), DTP
// (qualifier 573, Date Claim Check or Remittance Date), and REF
// (qualifiers 2U/EI/FY/NF/G1/9F/T4/F8). REF_OTHER_PAYER_NAME,
// DTP_OTHER_PAYER_NAME in frontend/src/data/x12Names837i.ts (source:
// Transmittal 907, PDF p.104-106).
type OtherPayerNM1PR struct {
	Name NameNM1 `json:"name,omitempty"`

	CheckOrRemittanceDTP573 DateDTP `json:"checkOrRemittanceDTP573,omitempty"` // Date Claim Check or Remittance Date; DTP*573, D8; PDF p.105

	// REF, named by qualifier per REF_OTHER_PAYER_NAME. 2U/EI/FY/NF all
	// resolve to the same public name, "Other Payer Secondary
	// Identification"; each keeps its own field since they are four
	// distinct qualifiers, the same way x12835 keeps separate fields for
	// qualifiers that share a name.
	SecondaryIDREF2U         ReferenceREF   `json:"secondaryIDREF2U,omitempty"`         // Other Payer Secondary Identification; REF*2U; PDF p.105
	SecondaryIDREFEI         ReferenceREF   `json:"secondaryIDREFEI,omitempty"`         // Other Payer Secondary Identification; REF*EI; PDF p.105
	SecondaryIDREFFY         ReferenceREF   `json:"secondaryIDREFFY,omitempty"`         // Other Payer Secondary Identification; REF*FY; PDF p.105
	SecondaryIDREFNF         ReferenceREF   `json:"secondaryIDREFNF,omitempty"`         // Other Payer Secondary Identification; REF*NF; PDF p.105
	PriorAuthREFG1           ReferenceREF   `json:"priorAuthREFG1,omitempty"`           // Other Payer Prior Authorization Number; REF*G1; PDF p.105
	ReferralNumberREF9F      ReferenceREF   `json:"referralNumberREF9F,omitempty"`      // Other Payer Referral Number; REF*9F; PDF p.105
	AdjustmentIndicatorREFT4 ReferenceREF   `json:"adjustmentIndicatorREFT4,omitempty"` // Other Payer Claim Adjustment Indicator; REF*T4; PDF p.106
	ClaimControlREFF8        ReferenceREF   `json:"claimControlREFF8,omitempty"`        // Other Payer Claim Control Number; REF*F8; PDF p.106
	References               []ReferenceREF `json:"references,omitempty"`               // any other 2330B REF, or a repeat of one above

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// OtherPayerAttendingNM171 is loop 2330C, Other Payer Attending Provider:
// NM1 (role 71) and up to four REF repeats (qualifiers 0B/1G/G2/LU, the
// same set PROVIDER_REF_NAME documents for the claim-level provider
// loops, claim.go). Unlike 2310A's AttendingNM171, this loop carries no
// PRV (OTHER_PAYER_PROVIDER_LOOP_NAME in frontend/src/data/x12Names837i.ts
// does not mention one). Transmittal 907, PDF p.106.
type OtherPayerAttendingNM171 struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// OtherPayerBillingNM185 is loop 2330I, Other Payer Billing Provider: NM1
// (role 85) and up to four REF repeats (qualifiers 0B/1G/G2/LU). This loop
// has no counterpart among claim.go's claim-level provider loops — the
// claim's own billing provider is loop 2010AA (BillingProviderHL), with a
// different (address-bearing) shape. Transmittal 907, PDF p.111.
type OtherPayerBillingNM185 struct {
	Name       NameNM1        `json:"name,omitempty"`
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// InpatientMIA holds one MIA segment (loop 2320, coordination of
// benefits), every element. X12 reuses the identical segment for
// inpatient institutional adjudication information in both the 837I and
// the 835; this type matches backend/pkg/x12835's own InpatientMIA field
// for field. Transmittal 907, PDF p.100-101; the four remark positions
// resolve as RARC the same way the 835's do (MIA_EL_NAME in
// frontend/src/data/x12Names837i.ts). Set is false when the segment did
// not appear at all.
type InpatientMIA struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	CoveredDaysMIA01      Amount    `json:"coveredDaysMIA01,omitempty"`      // Covered Days or Visits Count; PDF p.100
	OperatingOutlierMIA02 Amount    `json:"operatingOutlierMIA02,omitempty"` // PPS Operating Outlier Amount; PDF p.100
	PsychDaysMIA03        Amount    `json:"psychDaysMIA03,omitempty"`        // Lifetime Psychiatric Days Count; PDF p.100
	DRGAmountMIA04        Amount    `json:"drgAmountMIA04,omitempty"`        // Claim DRG Amount; PDF p.100
	RemarkMIA05           rarc.Code `json:"remarkMIA05,omitempty"`           // Claim Payment Remark Code; PDF p.100
	DSHAmountMIA06        Amount    `json:"dshAmountMIA06,omitempty"`        // Claim Disproportionate Share Amount; PDF p.100
	MSPPassThruMIA07      Amount    `json:"mspPassThruMIA07,omitempty"`      // Claim MSP Pass-Through Amount; PDF p.100
	PPSCapitalMIA08       Amount    `json:"ppsCapitalMIA08,omitempty"`       // Claim PPS Capital Amount; PDF p.100
	CapitalFSPMIA09       Amount    `json:"capitalFSPMIA09,omitempty"`       // PPS-Capital FSP DRG Amount; PDF p.100
	CapitalHSPMIA10       Amount    `json:"capitalHSPMIA10,omitempty"`       // PPS-Capital HSP DRG Amount; PDF p.100
	CapitalDSHMIA11       Amount    `json:"capitalDSHMIA11,omitempty"`       // PPS-Capital DSH DRG Amount; PDF p.100
	OldCapitalMIA12       Amount    `json:"oldCapitalMIA12,omitempty"`       // Old Capital Amount; PDF p.100
	CapitalIMEMIA13       Amount    `json:"capitalIMEMIA13,omitempty"`       // PPS-Capital IME Amount; PDF p.100
	OperHSPSpecMIA14      Amount    `json:"operHSPSpecMIA14,omitempty"`      // PPS-Operating Hospital Specific DRG Amount; PDF p.100
	CostReportDaysMIA15   Amount    `json:"costReportDaysMIA15,omitempty"`   // Cost Report Day Count; PDF p.100
	OperFSPSpecMIA16      Amount    `json:"operFSPSpecMIA16,omitempty"`      // PPS-Operating Federal Specific DRG Amount; PDF p.100
	PPSOutlierMIA17       Amount    `json:"ppsOutlierMIA17,omitempty"`       // Claim PPS Capital Outlier Amount; PDF p.100
	IndirectTeachingMIA18 Amount    `json:"indirectTeachingMIA18,omitempty"` // Claim Indirect Teaching Amount; PDF p.100
	NonPayProfCompMIA19   Amount    `json:"nonPayProfCompMIA19,omitempty"`   // Nonpayable Professional Component Amount; PDF p.100
	RemarkMIA20           rarc.Code `json:"remarkMIA20,omitempty"`           // Claim Payment Remark Code; PDF p.100
	RemarkMIA21           rarc.Code `json:"remarkMIA21,omitempty"`           // Claim Payment Remark Code; PDF p.100
	RemarkMIA22           rarc.Code `json:"remarkMIA22,omitempty"`           // Claim Payment Remark Code; PDF p.100
	RemarkMIA23           rarc.Code `json:"remarkMIA23,omitempty"`           // Claim Payment Remark Code; PDF p.100
	CapitalExceptionMIA24 Amount    `json:"capitalExceptionMIA24,omitempty"` // PPS-Capital Exception Amount; PDF p.100
}

// OutpatientMOA holds one MOA segment (loop 2320, coordination of
// benefits), every element, matching backend/pkg/x12835's own OutpatientMOA
// field for field (see InpatientMIA's doc comment). Transmittal 907, PDF
// p.102; the five remark positions resolve as RARC. Set is false when the
// segment did not appear at all.
type OutpatientMOA struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	ReimbursementRateMOA01  Amount    `json:"reimbursementRateMOA01,omitempty"`  // Reimbursement Rate; PDF p.102
	HCPCSPayableMOA02       Amount    `json:"hcpcsPayableMOA02,omitempty"`       // HCPCS Payable Amount; PDF p.102
	RemarkMOA03             rarc.Code `json:"remarkMOA03,omitempty"`             // Claim Payment Remark Code; PDF p.102
	RemarkMOA04             rarc.Code `json:"remarkMOA04,omitempty"`             // Claim Payment Remark Code; PDF p.102
	RemarkMOA05             rarc.Code `json:"remarkMOA05,omitempty"`             // Claim Payment Remark Code; PDF p.102
	RemarkMOA06             rarc.Code `json:"remarkMOA06,omitempty"`             // Claim Payment Remark Code; PDF p.102
	RemarkMOA07             rarc.Code `json:"remarkMOA07,omitempty"`             // Claim Payment Remark Code; PDF p.102
	ESRDPaymentMOA08        Amount    `json:"esrdPaymentMOA08,omitempty"`        // End Stage Renal Disease Payment Amount; PDF p.102
	NonpayableProfCompMOA09 Amount    `json:"nonpayableProfCompMOA09,omitempty"` // Nonpayable Professional Component Amount; PDF p.102
}

// OtherHealthInsuranceOI holds one OI segment (loop 2320, Other Health
// Insurance Information). OI01 (Claim Filing Indicator Code), OI02 (Claim
// Submission Reason Code), OI04 (Patient Signature Source Code), and OI05
// (Provider Agreement Code) are Not Used.
//
// OI03 and OI06 use OtherInsuranceIndicator (enums_cob.go), the one type
// frontend/src/data/x12Names837i.ts names for both positions. Their codes
// match CLM08's and CLM09's, but OI06_NAME gives Y the short name "Yes"
// where CLM09_NAME gives the longer signed-statement text, so the CLM
// types are not reused here.
type OtherHealthInsuranceOI struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	FilingIndicatorOI01   string                  `json:"filingIndicatorOI01,omitempty"`   // Claim Filing Indicator Code (N/U); OI01; PDF p.101
	SubmissionReasonOI02  string                  `json:"submissionReasonOI02,omitempty"`  // Claim Submission Reason Code (N/U); OI02; PDF p.101
	BenefitsAssignOI03    OtherInsuranceIndicator `json:"benefitsAssignOI03,omitempty"`    // Benefits Assignment Certification Indicator; OI03; PDF p.101
	SignatureSourceOI04   string                  `json:"signatureSourceOI04,omitempty"`   // Patient Signature Source Code (N/U); OI04; PDF p.101
	ProviderAgreementOI05 string                  `json:"providerAgreementOI05,omitempty"` // Provider Agreement Code (N/U); OI05; PDF p.101
	ReleaseOfInfoOI06     OtherInsuranceIndicator `json:"releaseOfInfoOI06,omitempty"`     // Release of Information Code; OI06; PDF p.101
}
