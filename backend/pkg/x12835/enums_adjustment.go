package x12835

// AdjGroupCAS01 is CAS01, the Claim Adjustment Group Code, shared by the
// claim-loop (2100) and service-line (2110) CAS segments.
type AdjGroupCAS01 string

const (
	AdjGroupCAS01NotSent AdjGroupCAS01 = ""
	AdjGroupCAS01CO      AdjGroupCAS01 = "CO"
	AdjGroupCAS01OA      AdjGroupCAS01 = "OA"
	AdjGroupCAS01PR      AdjGroupCAS01 = "PR"
	AdjGroupCAS01PI      AdjGroupCAS01 = "PI"
)

// adjGroupCAS01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var adjGroupCAS01Codes = []codeEntry{
	{"CO", "Contractual Obligation", "Contractual Obligation — x12.org/codes/claim-adjustment-group-codes; CMS 835 Flat File, PDF p.6, 12; CEDI 835 CG, PDF p.25 (printed p.19); frontend/src/data/x12Names835.ts CAS_GROUP_NAME"},
	{"OA", "Other Adjustment", "Other Adjustment — x12.org/codes/claim-adjustment-group-codes; CMS 835 Flat File, PDF p.6, 12; frontend/src/data/x12Names835.ts CAS_GROUP_NAME"},
	{"PR", "Patient Responsibility", "Patient Responsibility — x12.org/codes/claim-adjustment-group-codes; CMS 835 Flat File, PDF p.6, 12; frontend/src/data/x12Names835.ts CAS_GROUP_NAME"},
	{"PI", "Payor Initiated Reductions", "Payor Initiated Reductions — X12 005010X221A1 TR3, CAS01 (transcribed from the standard; verify against the printed guide); x12.org/codes/claim-adjustment-group-codes; the CMS 835 Flat File (PDF p.6) lists CO, OA, and PR for Medicare"},
}

func (v AdjGroupCAS01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AdjGroupCAS01) Known() bool { return known(adjGroupCAS01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AdjGroupCAS01) ShortDescription() string { return shortOf(adjGroupCAS01Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AdjGroupCAS01) Description() string { return descriptionOf(adjGroupCAS01Codes, string(v)) }

// ParseAdjGroupCAS01 reads one CAS01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAdjGroupCAS01(code string) (AdjGroupCAS01, error) {
	return AdjGroupCAS01(code), checkListed("CAS01", adjGroupCAS01Codes, code)
}

// PLBReasonPLB03 is the Provider Adjustment Reason Code, the first component
// of the composite identifier in each PLB03/05/07/09/11/13 slot.
//
// CEDI 835 CG, PDF p.29 (printed p.23; Table 25) states Medicare's MACs use
// 35 of these codes, excluding AH, AM, CR, CT, CW, and FC; all 41 are
// listed here, the six Medicare does not use with that note. Of the 35, 28
// have a definition on X12's current public provider-adjustment-reason-codes
// list (frontend/src/data/x12Names835.ts PLB_REASON_NAME); the other 7 (C5,
// CV, DM, GO, LS, OA, RE) are not on that public list's current edition, so
// their name here comes from the CMS 835 Flat File's own description as
// read and reported in docs/design/835-examples.md, not from X12's public
// list. Note
// PLBReasonPLB03OA is a different code from AdjGroupCAS01OA: same two
// letters, unrelated code lists (provider adjustment reason vs. claim
// adjustment group).
type PLBReasonPLB03 string

const (
	PLBReasonPLB03NotSent                             PLBReasonPLB03 = ""
	PLBReasonPLB03LateChargeOrPenalty                 PLBReasonPLB03 = "50"
	PLBReasonPLB03LateFilingInterestPenalty           PLBReasonPLB03 = "51"
	PLBReasonPLB03ProviderRefund                      PLBReasonPLB03 = "72"
	PLBReasonPLB03EarlyPaymentAllowance               PLBReasonPLB03 = "90"
	PLBReasonPLB03AdvancedOrAcceleratedPayment        PLBReasonPLB03 = "AP"
	PLBReasonPLB03Rebate                              PLBReasonPLB03 = "B2"
	PLBReasonPLB03RecoveryFromProvider                PLBReasonPLB03 = "B3"
	PLBReasonPLB03BadDebt                             PLBReasonPLB03 = "BD"
	PLBReasonPLB03Bonus                               PLBReasonPLB03 = "BN"
	PLBReasonPLB03AdjustmentDetailedSeparately        PLBReasonPLB03 = "CS"
	PLBReasonPLB03Withholding                         PLBReasonPLB03 = "E3"
	PLBReasonPLB03BalanceForward                      PLBReasonPLB03 = "FB"
	PLBReasonPLB03HemophiliaClottingFactorAddOn       PLBReasonPLB03 = "HM"
	PLBReasonPLB03IncentivePayment                    PLBReasonPLB03 = "IP"
	PLBReasonPLB03IRS1099Withholding                  PLBReasonPLB03 = "IR"
	PLBReasonPLB03InterimRateLumpSum                  PLBReasonPLB03 = "IS"
	PLBReasonPLB03DemonstrationProgramLimitation      PLBReasonPLB03 = "J1"
	PLBReasonPLB03Penalty                             PLBReasonPLB03 = "L3"
	PLBReasonPLB03Interest                            PLBReasonPLB03 = "L6"
	PLBReasonPLB03IRSNon1099Withholding               PLBReasonPLB03 = "LE"
	PLBReasonPLB03AffiliatedProviderOffset            PLBReasonPLB03 = "OB"
	PLBReasonPLB03PeriodicInterimPayment              PLBReasonPLB03 = "PI"
	PLBReasonPLB03FinalPaymentOrSettlement            PLBReasonPLB03 = "PL"
	PLBReasonPLB03RetroactiveAdjustment               PLBReasonPLB03 = "RA"
	PLBReasonPLB03StudentLoanGarnishment              PLBReasonPLB03 = "SL"
	PLBReasonPLB03ThirdPartyLiability                 PLBReasonPLB03 = "TL"
	PLBReasonPLB03OverpaymentRecovery                 PLBReasonPLB03 = "WO"
	PLBReasonPLB03NonIRSWithholdingLevy               PLBReasonPLB03 = "WU"
	PLBReasonPLB03TemporaryAllowance                  PLBReasonPLB03 = "C5"
	PLBReasonPLB03CapitalPassthrough                  PLBReasonPLB03 = "CV"
	PLBReasonPLB03DirectMedicalEducationPassthrough   PLBReasonPLB03 = "DM"
	PLBReasonPLB03GraduateMedicalEducationPassthrough PLBReasonPLB03 = "GO"
	PLBReasonPLB03LumpSum                             PLBReasonPLB03 = "LS"
	PLBReasonPLB03OrganAcquisition                    PLBReasonPLB03 = "OA"
	PLBReasonPLB03ReturnOnEquity                      PLBReasonPLB03 = "RE"
	PLBReasonPLB03OriginationFee                      PLBReasonPLB03 = "AH"
	PLBReasonPLB03AppliedToBorrowersAccount           PLBReasonPLB03 = "AM"
	PLBReasonPLB03CapitationInterest                  PLBReasonPLB03 = "CR"
	PLBReasonPLB03CapitationPayment                   PLBReasonPLB03 = "CT"
	PLBReasonPLB03CRNAPassthru                        PLBReasonPLB03 = "CW"
	PLBReasonPLB03FundAllocation                      PLBReasonPLB03 = "FC"
)

// plbReasonPLB03Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var plbReasonPLB03Codes = []codeEntry{
	{"50", "The amount of the late charge, late claim filing penalty, or Medicare late cost report penalty", "The amount of the late charge, late claim filing penalty, or Medicare late cost report penalty — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"51", "Late filing interest penalty assessment amount", "Late filing interest penalty assessment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"72", "Provider refund amount", "Provider refund amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"90", "Early payment allowance amount", "Early payment allowance amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"AP", "Advanced or accelerated payment amount", "Advanced or accelerated payment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"B2", "Rebate amount", "Rebate amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"B3", "Recovery amount (from the provider, for an overpayment based on payments from other payers)", "Recovery amount (from the provider, for an overpayment based on payments from other payers) — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"BD", "Bad debt amount", "Bad debt amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"BN", "Bonus amount", "Bonus amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"CS", "Adjustment amount, detailed information provided separately", "Adjustment amount, detailed information provided separately — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"E3", "Withholding amount", "Withholding amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"FB", "Non-claim related balance forward amount", "Non-claim related balance forward amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"HM", "Hemophilia clotting factor add-on payment amount", "Hemophilia clotting factor add-on payment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"IP", "Incentive payment amount", "Incentive payment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"IR", "Internal Revenue Service 1099 withholding amount", "Internal Revenue Service 1099 withholding amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"IS", "Lump sum based on an interim rate", "Lump sum based on an interim rate — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"J1", "Amount not reimbursed based on a demonstration program or other limitation", "Amount not reimbursed based on a demonstration program or other limitation — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"L3", "Penalty amount", "Penalty amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"L6", "Interest amount", "Interest amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"LE", "Internal Revenue Service non-1099 withholding amount", "Internal Revenue Service non-1099 withholding amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"OB", "Affiliated provider(s) offset amount", "Affiliated provider(s) offset amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"PI", "Periodic Interim Payment (PIP) lump sum amount", "Periodic Interim Payment (PIP) lump sum amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"PL", "Final payment or settlement amount", "Final payment or settlement amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"RA", "Retroactive adjustment amount", "Retroactive adjustment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"SL", "Student loan garnishment amount", "Student loan garnishment amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"TL", "Third party liability determination amount", "Third party liability determination amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"WO", "Overpayment recovery amount", "Overpayment recovery amount — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"WU", "Non-Internal Revenue Service withholding amount related to a federal payment levy program", "Non-Internal Revenue Service withholding amount related to a federal payment levy program — x12.org/codes/provider-adjustment-reason-codes (element 426); frontend/src/data/x12Names835.ts PLB_REASON_NAME"},
	{"C5", "Temporary allowance", "Temporary allowance — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"CV", "Capital passthrough", "Capital passthrough — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"DM", "Direct medical education passthrough", "Direct medical education passthrough — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"GO", "Graduate medical education passthrough", "Graduate medical education passthrough — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"LS", "Lump sum", "Lump sum — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"OA", "Organ acquisition", "Organ acquisition — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"RE", "Return on equity", "Return on equity — CMS 835 Flat File, PDF p.13-14 description, per docs/design/835-examples.md's cost-report-settlements-and-passthroughs example; not on X12's public provider-adjustment-reason-codes list, x12.org/codes/provider-adjustment-reason-codes"},
	{"AH", "Origination Fee", "Origination Fee — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
	{"AM", "Applied to Borrower's Account", "Applied to Borrower's Account — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
	{"CR", "Capitation Interest", "Capitation Interest — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
	{"CT", "Capitation Payment", "Capitation Payment — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
	{"CW", "Certified Registered Nurse Anesthetist Passthru", "Certified Registered Nurse Anesthetist Passthru — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
	{"FC", "Fund Allocation", "Fund Allocation — X12 005010X221A1 TR3, PLB03-1 (transcribed from the standard; verify against the printed guide); CEDI 835 CG, PDF p.29 (printed p.23; Table 25): not used by Medicare's MACs"},
}

func (v PLBReasonPLB03) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v PLBReasonPLB03) Known() bool { return known(plbReasonPLB03Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v PLBReasonPLB03) ShortDescription() string { return shortOf(plbReasonPLB03Codes, string(v)) }

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v PLBReasonPLB03) Description() string { return descriptionOf(plbReasonPLB03Codes, string(v)) }

// ParsePLBReasonPLB03 reads one PLB03 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParsePLBReasonPLB03(code string) (PLBReasonPLB03, error) {
	return PLBReasonPLB03(code), checkListed("PLB03", plbReasonPLB03Codes, code)
}
