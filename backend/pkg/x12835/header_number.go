package x12835

// HeaderNumberLX is the LX loop: a header number, an optional provider
// summary (TS3) and institutional supplemental summary (TS2), and the
// claims it groups. Source: CMS 835 Flat File, PDF p.4.
type HeaderNumberLX struct {
	Index      int `json:"index,omitempty"`      // the LX segment's position in the file
	NumberLX01 int `json:"numberLX01,omitempty"` // Assigned # — CMS 835 Flat File, PDF p.4

	ProviderSummary      ProviderSummaryTS3      `json:"providerSummary,omitempty"`      // TS3, optional
	InstitutionalSummary InstitutionalSummaryTS2 `json:"institutionalSummary,omitempty"` // TS2, optional

	Claims []ClaimCLP `json:"claims,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop before its first claim that the model does not place
}

// ProviderSummaryTS3 is the TS3 segment: provider-level summary totals.
// Element names and citations are frontend/src/data/x12Names835.ts's
// TS3_EL_NAME array (source: CMS 835 Flat File, PDF p.4-5). Set is false
// when the segment did not appear at all; every amount and count below is
// independently optional even when the segment is present (Amount's own
// Set flag says whether each one was actually sent).
type ProviderSummaryTS3 struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	ProviderIDTS301       string `json:"providerIDTS301,omitempty"`       // Provider Identifier — CMS 835 Flat File, PDF p.4
	FacilityCodeTS302     string `json:"facilityCodeTS302,omitempty"`     // Facility Code Value — CMS 835 Flat File, PDF p.4
	FiscalPeriodDateTS303 Date   `json:"fiscalPeriodDateTS303,omitempty"` // Fiscal Period Date — CMS 835 Flat File, PDF p.4

	ClaimCountTS304                   Amount `json:"claimCountTS304,omitempty"`                   // Total Claim Count — CMS 835 Flat File, PDF p.4
	ClaimChangeAmountTS305            Amount `json:"claimChangeAmountTS305,omitempty"`            // Total Claim Change Amount — CMS 835 Flat File, PDF p.4
	CoveredChargeAmountTS306          Amount `json:"coveredChargeAmountTS306,omitempty"`          // Total Covered Charge Amount — CMS 835 Flat File, PDF p.4
	NoncoveredChargeAmountTS307       Amount `json:"noncoveredChargeAmountTS307,omitempty"`       // Total Noncovered Charge Amount — CMS 835 Flat File, PDF p.4
	DeniedChargeAmountTS308           Amount `json:"deniedChargeAmountTS308,omitempty"`           // Total Denied Charge Amount — CMS 835 Flat File, PDF p.4
	ProviderAmountTS309               Amount `json:"providerAmountTS309,omitempty"`               // Total Provider Amount — CMS 835 Flat File, PDF p.4
	InterestAmountTS310               Amount `json:"interestAmountTS310,omitempty"`               // Total Interest Amount — CMS 835 Flat File, PDF p.4
	ContractualAdjustmentAmountTS311  Amount `json:"contractualAdjustmentAmountTS311,omitempty"`  // Total Contractual Adjustment Amount — CMS 835 Flat File, PDF p.4
	GrammRudmanReductionAmountTS312   Amount `json:"grammRudmanReductionAmountTS312,omitempty"`   // Total Gramm-Rudman Reduction Amount — CMS 835 Flat File, PDF p.4
	MSPPayerAmountTS313               Amount `json:"mspPayerAmountTS313,omitempty"`               // Total MSP Payer Amount — CMS 835 Flat File, PDF p.4
	BloodDeductibleAmountTS314        Amount `json:"bloodDeductibleAmountTS314,omitempty"`        // Total Blood Deductible Amount — CMS 835 Flat File, PDF p.4
	NonLabChargeAmountTS315           Amount `json:"nonLabChargeAmountTS315,omitempty"`           // Total Non-Lab Charge Amount — CMS 835 Flat File, PDF p.4
	CoinsuranceAmountTS316            Amount `json:"coinsuranceAmountTS316,omitempty"`            // Total Coinsurance Amount — CMS 835 Flat File, PDF p.4
	HCPCSReportedChargeAmountTS317    Amount `json:"hcpcsReportedChargeAmountTS317,omitempty"`    // Total HCPCS Reported Charge Amount — CMS 835 Flat File, PDF p.4
	HCPCSPayableAmountTS318           Amount `json:"hcpcsPayableAmountTS318,omitempty"`           // Total HCPCS Payable Amount — CMS 835 Flat File, PDF p.4
	DeductibleAmountTS319             Amount `json:"deductibleAmountTS319,omitempty"`             // Total Deductible Amount — CMS 835 Flat File, PDF p.4
	ProfessionalComponentAmountTS320  Amount `json:"professionalComponentAmountTS320,omitempty"`  // Total Professional Component Amount — CMS 835 Flat File, PDF p.4
	MSPPatientLiabilityMetAmountTS321 Amount `json:"mspPatientLiabilityMetAmountTS321,omitempty"` // Total MSP Patient Liability Met Amount — CMS 835 Flat File, PDF p.5
	PatientReimbursementAmountTS322   Amount `json:"patientReimbursementAmountTS322,omitempty"`   // Total Patient Reimbursement Amount — CMS 835 Flat File, PDF p.5
	PIPClaimCountTS323                Amount `json:"pipClaimCountTS323,omitempty"`                // Total PIP Claim Count — CMS 835 Flat File, PDF p.5
	PIPAdjustmentAmountTS324          Amount `json:"pipAdjustmentAmountTS324,omitempty"`          // Total PIP Adjustment Amount — CMS 835 Flat File, PDF p.5
}

// InstitutionalSummaryTS2 is the TS2 segment: institutional supplemental
// summary information. Element names and citations are
// frontend/src/data/x12Names835.ts's TS2_EL_NAME array (source: CMS 835
// Flat File, PDF p.5). Set is false when the segment did not appear at
// all; every field below is independently optional even when the segment
// is present.
type InstitutionalSummaryTS2 struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	DRGAmountTS201                      Amount `json:"drgAmountTS201,omitempty"`                      // Total DRG Amount — CMS 835 Flat File, PDF p.5
	FederalSpecificAmountTS202          Amount `json:"federalSpecificAmountTS202,omitempty"`          // Total Federal Specific Amount — CMS 835 Flat File, PDF p.5
	HospitalSpecificAmountTS203         Amount `json:"hospitalSpecificAmountTS203,omitempty"`         // Total Hospital Specific Amount — CMS 835 Flat File, PDF p.5
	DisproportionateAmountTS204         Amount `json:"disproportionateAmountTS204,omitempty"`         // Total Disproportionate Amount — CMS 835 Flat File, PDF p.5
	CapitalAmountTS205                  Amount `json:"capitalAmountTS205,omitempty"`                  // Total Capital Amount — CMS 835 Flat File, PDF p.5
	IndirectMedicalEducationAmountTS206 Amount `json:"indirectMedicalEducationAmountTS206,omitempty"` // Total Indirect Medical Education Amount — CMS 835 Flat File, PDF p.5
	OutlierDayCountTS207                Amount `json:"outlierDayCountTS207,omitempty"`                // Total Outlier Day Count — CMS 835 Flat File, PDF p.5
	DayOutlierAmountTS208               Amount `json:"dayOutlierAmountTS208,omitempty"`               // Total Day Outlier Amount — CMS 835 Flat File, PDF p.5
	CostOutlierAmountTS209              Amount `json:"costOutlierAmountTS209,omitempty"`              // Total Cost Outlier Amount — CMS 835 Flat File, PDF p.5
	AverageDRGLengthOfStayTS210         Amount `json:"averageDRGLengthOfStayTS210,omitempty"`         // Average DRG Length of Stay — CMS 835 Flat File, PDF p.5
	DischargeCountTS211                 Amount `json:"dischargeCountTS211,omitempty"`                 // Total Discharge Count — CMS 835 Flat File, PDF p.5
	CostReportDayCountTS212             Amount `json:"costReportDayCountTS212,omitempty"`             // Total Cost Report Day Count — CMS 835 Flat File, PDF p.5
	CoveredDayCountTS213                Amount `json:"coveredDayCountTS213,omitempty"`                // Total Covered Day Count — CMS 835 Flat File, PDF p.5
	NoncoveredDayCountTS214             Amount `json:"noncoveredDayCountTS214,omitempty"`             // Total Noncovered Day Count — CMS 835 Flat File, PDF p.5
	MSPPassThroughAmountTS215           Amount `json:"mspPassThroughAmountTS215,omitempty"`           // Total MSP Pass-Through Amount — CMS 835 Flat File, PDF p.5
	AverageDRGWeightTS216               Amount `json:"averageDRGWeightTS216,omitempty"`               // Average DRG Weight — CMS 835 Flat File, PDF p.5
	PPSCapitalFSPDRGAmountTS217         Amount `json:"ppsCapitalFSPDRGAmountTS217,omitempty"`         // Total PPS Capital FSP DRG Amount — CMS 835 Flat File, PDF p.5
	PSPCapitalHSPDRGAmountTS218         Amount `json:"pspCapitalHSPDRGAmountTS218,omitempty"`         // Total PSP Capital HSP DRG Amount — CMS 835 Flat File, PDF p.5 (spelled "PSP" in the source table itself)
	PPSDSHDRGAmountTS219                Amount `json:"ppsdshdrgAmountTS219,omitempty"`                // Total PPS DSH DRG Amount — CMS 835 Flat File, PDF p.5
}
