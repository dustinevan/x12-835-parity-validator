package x12835

import "github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"

// InpatientMIA holds one MIA segment (loop 2100), every element. CMS 835
// Flat File, PDF p.8-9; the four remark positions resolve as RARC, same as
// MOA's (docs/design/835-viewer.md, MIA05/20-23 row). Set is false when
// the segment did not appear at all.
type InpatientMIA struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	CoveredDaysMIA01      Amount    `json:"coveredDaysMIA01,omitempty"`      // Covered Days or Visits Count; PDF p.8
	OperatingOutlierMIA02 Amount    `json:"operatingOutlierMIA02,omitempty"` // PPS Operating Outlier Amount; PDF p.8
	PsychDaysMIA03        Amount    `json:"psychDaysMIA03,omitempty"`        // Lifetime Psychiatric Days Count; PDF p.9
	DRGAmountMIA04        Amount    `json:"drgAmountMIA04,omitempty"`        // Claim DRG Amount; PDF p.9
	RemarkMIA05           rarc.Code `json:"remarkMIA05,omitempty"`           // Claim Payment Remark Code; PDF p.9
	DSHAmountMIA06        Amount    `json:"dshAmountMIA06,omitempty"`        // Claim DSH Amount; PDF p.9
	MSPPassThruMIA07      Amount    `json:"mspPassThruMIA07,omitempty"`      // Claim MSP Pass Thru Amt; PDF p.9
	PPSCapitalMIA08       Amount    `json:"ppsCapitalMIA08,omitempty"`       // Claim PPS Capital Amount; PDF p.9
	CapitalFSPMIA09       Amount    `json:"capitalFSPMIA09,omitempty"`       // PPS Capital FSP DRG Amt; PDF p.9
	CapitalHSPMIA10       Amount    `json:"capitalHSPMIA10,omitempty"`       // PPS Capital HSP DRG Amt; PDF p.9
	CapitalDSHMIA11       Amount    `json:"capitalDSHMIA11,omitempty"`       // PPS Capital DSH DRG Amt; PDF p.9
	OldCapitalMIA12       Amount    `json:"oldCapitalMIA12,omitempty"`       // Old Capital Amount; PDF p.9
	CapitalIMEMIA13       Amount    `json:"capitalIMEMIA13,omitempty"`       // PPS Capital IME Amount; PDF p.9
	OperHSPSpecMIA14      Amount    `json:"operHSPSpecMIA14,omitempty"`      // PPS Oper HSP Spec DRG Amt; PDF p.9
	CostReportDaysMIA15   Amount    `json:"costReportDaysMIA15,omitempty"`   // Cost Report Day Count; PDF p.9
	OperFSPSpecMIA16      Amount    `json:"operFSPSpecMIA16,omitempty"`      // PPS Oper FSP Spec DRG Amt; PDF p.9
	PPSOutlierMIA17       Amount    `json:"ppsOutlierMIA17,omitempty"`       // Claim PPS Outlier Amount; PDF p.9
	IndirectTeachingMIA18 Amount    `json:"indirectTeachingMIA18,omitempty"` // Claim Indirect Teaching; PDF p.9
	NonPayProfCompMIA19   Amount    `json:"nonPayProfCompMIA19,omitempty"`   // Non Pay Prof Comp Amt; PDF p.9
	RemarkMIA20           rarc.Code `json:"remarkMIA20,omitempty"`           // Claim Payment Remark Code; PDF p.9
	RemarkMIA21           rarc.Code `json:"remarkMIA21,omitempty"`           // Claim Payment Remark Code; PDF p.9
	RemarkMIA22           rarc.Code `json:"remarkMIA22,omitempty"`           // Claim Payment Remark Code; PDF p.9
	RemarkMIA23           rarc.Code `json:"remarkMIA23,omitempty"`           // Claim Payment Remark Code; PDF p.9
	CapitalExceptionMIA24 Amount    `json:"capitalExceptionMIA24,omitempty"` // PPS Capital Exception Amt; PDF p.9
}

// OutpatientMOA holds one MOA segment (loop 2100), every element. CMS 835
// Flat File, PDF p.10; the five remark positions resolve as RARC. Set is
// false when the segment did not appear at all.
type OutpatientMOA struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	ReimbursementRateMOA01  Amount    `json:"reimbursementRateMOA01,omitempty"`  // Reimbursement Rate; PDF p.10
	HCPCSPayableMOA02       Amount    `json:"hcpcsPayableMOA02,omitempty"`       // Claim HCPCS Payable Amount; PDF p.10
	RemarkMOA03             rarc.Code `json:"remarkMOA03,omitempty"`             // Claim Payment Remark Code; PDF p.10
	RemarkMOA04             rarc.Code `json:"remarkMOA04,omitempty"`             // Claim Payment Remark Code; PDF p.10
	RemarkMOA05             rarc.Code `json:"remarkMOA05,omitempty"`             // Claim Payment Remark Code; PDF p.10
	RemarkMOA06             rarc.Code `json:"remarkMOA06,omitempty"`             // Claim Payment Remark Code; PDF p.10
	RemarkMOA07             rarc.Code `json:"remarkMOA07,omitempty"`             // Claim Payment Remark Code; PDF p.10
	ESRDPaymentMOA08        Amount    `json:"esrdPaymentMOA08,omitempty"`        // Claim ESRD Payment Amount; PDF p.10
	NonpayableProfCompMOA09 Amount    `json:"nonpayableProfCompMOA09,omitempty"` // Nonpayable Professional Comp Amt; PDF p.10
}

// ClaimCLP holds one CLP segment (loop 2100) and everything the 835 nests
// under it, down through its service lines.
type ClaimCLP struct {
	Index int `json:"index,omitempty"` // the CLP segment's position in the file; 0 for a claim read by ParseClaim

	PatientControlCLP01    string               `json:"patientControlCLP01,omitempty"`    // Patient Control #; CMS 835 Flat File, PDF p.5
	StatusCLP02            ClaimStatusCLP02     `json:"statusCLP02,omitempty"`            // Claim Status Code; PDF p.5
	ChargeCLP03            Amount               `json:"chargeCLP03,omitempty"`            // Total Claim Charge Amount; PDF p.5
	PaymentCLP04           Amount               `json:"paymentCLP04,omitempty"`           // Claim Payment Amount; PDF p.6
	PatientRespCLP05       Amount               `json:"patientRespCLP05,omitempty"`       // Patient Responsibility Amount; PDF p.6
	FilingIndicatorCLP06   FilingIndicatorCLP06 `json:"filingIndicatorCLP06,omitempty"`   // Claim Filing Indicator Code; PDF p.6
	PayerClaimControlCLP07 string               `json:"payerClaimControlCLP07,omitempty"` // Payer Claim Control #; PDF p.6
	FacilityCLP08          string               `json:"facilityCLP08,omitempty"`          // Facility Code Value; PDF p.6 (external code list)
	FrequencyCLP09         string               `json:"frequencyCLP09,omitempty"`         // Claim Frequency Code; PDF p.6 (external code list)
	PatientStatusCLP10     string               `json:"patientStatusCLP10,omitempty"`     // Patient Status Code; PDF p.6 (external code list, not covered by an enum here)
	DRGCLP11               string               `json:"drgCLP11,omitempty"`               // DRG Code; PDF p.6
	DRGWeightCLP12         Amount               `json:"drgWeightCLP12,omitempty"`         // DRG Weight; PDF p.6
	DischargeFractionCLP13 Amount               `json:"dischargeFractionCLP13,omitempty"` // Discharge Fraction; PDF p.6
	YesNoCLP14             string               `json:"yesNoCLP14,omitempty"`             // Yes/No Condition or Response Code; PDF p.6 (not covered by an enum here)

	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"` // CAS, "Claim Adjustment" at this level (CLAIM_CAS_NAME); PDF p.6

	// NM1s (loop 2100). NM1_NAME in frontend/src/data/x12Names835.ts names
	// each role; the field name carries the qualifier.
	PatientNM1QC          NameNM1 `json:"patientNM1QC,omitempty"`          // Patient Name; PDF p.6
	CorrectedPatientNM174 NameNM1 `json:"correctedPatientNM174,omitempty"` // Corrected Patient/Insured Name; PDF p.7
	RenderingNM182        NameNM1 `json:"renderingNM182,omitempty"`        // Service Provider Name; PDF p.7
	CrossoverNM1TT        NameNM1 `json:"crossoverNM1TT,omitempty"`        // Crossover Carrier Name; PDF p.7
	CorrectedPayerNM1PR   NameNM1 `json:"correctedPayerNM1PR,omitempty"`   // Corrected Priority Payer Name; PDF p.7-8

	InpatientMIA  InpatientMIA  `json:"inpatientMIA,omitempty"`  // MIA, present only for an inpatient claim
	OutpatientMOA OutpatientMOA `json:"outpatientMOA,omitempty"` // MOA, present only for an outpatient claim

	// REF (loop 2100). CLAIM_REF_NAME gives qualifiers EA and F8 the same
	// combined name, so both stay in one generic slice rather than separate
	// fields.
	References []ReferenceREF `json:"references,omitempty"`

	// Dates (loop 2100). DTM_NAME in frontend/src/data/x12Names835.ts. Each
	// field holds the first DTM with its qualifier that carries a date, and is
	// present when its DateQualifierDTM01 is not empty.
	StatementFromDTM232      DateTimeDTM `json:"statementFromDTM232,omitempty"`      // Statement From-To Date (from); PDF p.10
	StatementToDTM233        DateTimeDTM `json:"statementToDTM233,omitempty"`        // Statement From-To Date (to); PDF p.10
	ReceivedDTM050           DateTimeDTM `json:"receivedDTM050,omitempty"`           // Claim Received Date; PDF p.11
	CoverageExpirationDTM036 DateTimeDTM `json:"coverageExpirationDTM036,omitempty"` // Coverage Expiration Date; PDF p.10

	ContactPER ContactPER `json:"contactPER,omitempty"` // PER, qualifier CX, claim-level contact, present when its FunctionPER01 is not empty; PDF p.11

	Amounts    []AmountAMT   `json:"amounts,omitempty"`    // AMT, "Claim Supplemental Information" (CLAIM_AMT_NAME); PDF p.11
	Quantities []QuantityQTY `json:"quantities,omitempty"` // QTY, "Claim Supplemental Information Quantity" (CLAIM_QTY_NAME); PDF p.11

	ServiceLines []ServiceLineSVC `json:"serviceLines,omitempty"` // SVC, loop 2110, one per service line
	Raw          []RawSegment     `json:"raw,omitempty"`          // segments within this claim the model does not place
}
