// This file holds loop 2400 (service line, LX/SV2) and its sub-loops: 2410
// (drug identification, LIN/CTP), 2420A-D (line-level providers, by role),
// and 2430 (line adjudication information, SVD). The line-level PWK and HCP
// segments use claim.go's SupplementalPWK and PricingHCP, the same shape
// at both levels. Values over pointers; parse.go fills these types.

package x12837i

// ProcedureSV202 is the Composite Medical Procedure Identifier (C003)
// carried by SV202 and, in the same shape, by SVD03 (loop 2430's
// adjudicated procedure). Transmittal 907 names the composite but not its
// sub-components individually (PDF p.111-112 for SV202, p.122 for SVD03);
// the shape below (qualifier, procedure code, up to four modifiers,
// description, and the composite's eighth component) is the composite's
// standard layout, following backend/pkg/x12835's ProcedureSVC01, the
// 835's equivalent composite.
type ProcedureSV202 struct {
	Qualifier   ProductServiceQualifier `json:"qualifier,omitempty"`   // SV202-1 / SVD03-1, Product or Service ID Qualifier; same set reused at HCP09
	Code        string                  `json:"code,omitempty"`        // SV202-2 / SVD03-2, Procedure Code; CPT codes are kept as their code only (AMA-licensed)
	Modifiers   [4]string               `json:"modifiers,omitempty"`   // SV202-3..6 / SVD03-3..6, up to four procedure modifiers
	Description string                  `json:"description,omitempty"` // SV202-7 / SVD03-7, free-text description, sent in place of a code

	// ProductServiceID is the composite's eighth component (C003-08,
	// Product/Service ID), kept as sent. X12 defines it; the 837I guide
	// marks it Not Used.
	ProductServiceID string `json:"productServiceID,omitempty"`
}

// InstitutionalServiceSV2 holds the SV2 segment, the Institutional Service
// Line (loop 2400). Transmittal 907, PDF p.111-112. Set is false when the
// segment did not appear at all.
type InstitutionalServiceSV2 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	RevenueCodeSV201       string                   `json:"revenueCodeSV201,omitempty"`       // Revenue Code; NUBC code, external list, not reproduced here
	ProcedureSV202         ProcedureSV202           `json:"procedureSV202,omitempty"`         // Composite Medical Procedure Identifier
	ChargeSV203            Amount                   `json:"chargeSV203,omitempty"`            // Line Item Charge Amount
	UnitSV204              UnitOrBasisOfMeasurement `json:"unitSV204,omitempty"`              // Unit or Basis for Measurement Code; Transmittal 907 gives SV204 no value list of its own — this reuses UnitOrBasisOfMeasurement (HCP11/CTP05-1's set) as the same abstract X12 element 355, not independently confirmed for SV204
	UnitsSV205             Amount                   `json:"unitsSV205,omitempty"`             // Service Units/Days Count
	UnitRateSV206          Amount                   `json:"unitRateSV206,omitempty"`          // Unit Rate (N/U)
	DeniedChargeSV207      Amount                   `json:"deniedChargeSV207,omitempty"`      // Line Item Denied or Non-Covered Charge Amount
	ConditionSV208         string                   `json:"conditionSV208,omitempty"`         // Yes/No Condition or Response Code (N/U)
	ResidentialStatusSV209 string                   `json:"residentialStatusSV209,omitempty"` // Nursing Home Residential Status Code (N/U)
	LevelOfCareSV210       string                   `json:"levelOfCareSV210,omitempty"`       // Level of Care Code (N/U)
}

// DrugQuantityCTP holds the CTP segment of a 2410 loop, the Drug Quantity.
// Transmittal 907, PDF p.115-116. CTP05 is itself a composite, "Composite
// Unit of Measure" (C001); only its first component (the code qualifier)
// is documented, and the rest of the composite and CTP06-11 are Not Used,
// so they are kept as sent. Set is false when the segment did not appear
// at all.
type DrugQuantityCTP struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	ClassOfTradeCTP01    string                   `json:"classOfTradeCTP01,omitempty"`    // Class of Trade Code (N/U); PDF p.115-116
	PriceIdentifierCTP02 string                   `json:"priceIdentifierCTP02,omitempty"` // Price Identifier Code (N/U); PDF p.115-116
	UnitPriceCTP03       Amount                   `json:"unitPriceCTP03,omitempty"`       // Unit Price (N/U); PDF p.115-116
	UnitCountCTP04       Amount                   `json:"unitCountCTP04,omitempty"`       // National Drug Unit Count; PDF p.116
	UnitCTP05_1          UnitOrBasisOfMeasurement `json:"unitCTP05_1,omitempty"`          // Composite Unit of Measure, first component; PDF p.116

	// UnitCTP05Rest holds components 2 through 15 of the CTP05 composite
	// (exponents and multipliers) as sent, trailing empties dropped. X12
	// defines them; the 837I guide marks them Not Used.
	UnitCTP05Rest []string `json:"unitCTP05Rest,omitempty"`

	// ExtraElements holds CTP06 through CTP11 as sent. X12 defines them;
	// the 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// DrugLIN holds one 2410 loop: the LIN drug identification segment, its CTP
// quantity segment, and any REF (prescription/compound drug association
// number). There is no Set field; presence is the slice entry itself
// (ServiceLineLX.Drugs).
type DrugLIN struct {
	Index int `json:"index,omitempty"` // the LIN's position in the file

	AssignedIDLIN01 string            `json:"assignedIDLIN01,omitempty"` // Assigned Identification (N/U); Transmittal 907, PDF p.115
	QualifierLIN02  NDCQualifierLIN02 `json:"qualifierLIN02,omitempty"`  // Product or Service ID Qualifier; PDF p.115
	NDCLIN03        string            `json:"ndcLIN03,omitempty"`        // National Drug Code; PDF p.115

	// ExtraElements holds LIN04 through LIN31 (further qualifier and
	// identifier pairs) as sent. X12 defines them; the 837I guide marks
	// them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`

	Quantity DrugQuantityCTP `json:"quantity,omitempty"` // CTP; PDF p.115-116

	// REF, qualifiers VY/XZ (REF_PRESCRIPTION_NAME names each
	// individually); any other REF at this loop goes in References.
	// Transmittal 907, PDF p.116-117.
	LinkSequenceREFVY ReferenceREF   `json:"linkSequenceREFVY,omitempty"` // Link Sequence Number
	PrescriptionREFXZ ReferenceREF   `json:"prescriptionREFXZ,omitempty"` // Pharmacy Prescription Number
	References        []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// LineProviderNM1 is one 2420A-D line-level-provider loop: an NM1 plus any
// REF secondary identification. The role (Operating, Other Operating,
// Rendering, Referring) is carried by the field holding this struct on
// ServiceLineLX (OperatingNM172, OtherOperatingNM1ZZ, RenderingNM182,
// ReferringNM1DN — the same pattern backend/pkg/x12835/claim.go uses for
// its own NameNM1 fields) and by Name.EntityIDNM101. Each of these four
// loops reuses the identical NM1 shape and REF qualifier set as its 2310
// claim-level counterpart (NM1_ROLE_NAME, PROVIDER_REF_NAME), and keeps its
// REF segments the same way those loops do: in References, qualifier and
// value together.
type LineProviderNM1 struct {
	Name NameNM1 `json:"name,omitempty"`

	// REF, qualifiers 0B/1G/G2/LU. PROVIDER_REF_NAME names 0B (State
	// License Number), 1G (Provider UPIN Number), and G2 (Provider
	// Commercial Number), from CMS Claims Processing Manual Ch. 25 §75.6;
	// LU is in Transmittal 907's value set at this loop but is not named in
	// any reviewed source.
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// LineAdjudicationSVD holds one 2430 loop: the SVD line adjudication
// segment and everything nested under it (CAS line adjustments, DTP, AMT).
// One LineAdjudicationSVD exists per other payer that adjudicated this
// service line (coordination of benefits).
type LineAdjudicationSVD struct {
	Index int `json:"index,omitempty"` // the SVD's position in the file

	OtherPayerIDSVD01 string         `json:"otherPayerIDSVD01,omitempty"` // Other Payer Primary Identifier; Transmittal 907, PDF p.121
	PaymentSVD02      Amount         `json:"paymentSVD02,omitempty"`      // Service Line Paid Amount; PDF p.122
	ProcedureSVD03    ProcedureSV202 `json:"procedureSVD03,omitempty"`    // Composite Medical Procedure Identifier; same composite shape as SV202; PDF p.122
	RevenueCodeSVD04  string         `json:"revenueCodeSVD04,omitempty"`  // Revenue Code; PDF p.122
	PaidUnitsSVD05    Amount         `json:"paidUnitsSVD05,omitempty"`    // Paid Service Unit Count; PDF p.122
	BundledLineSVD06  string         `json:"bundledLineSVD06,omitempty"`  // Bundled or Unbundled Line Number; PDF p.122

	Adjustments []AdjustmentCAS `json:"adjustments,omitempty"` // CAS, "Line Adjustment" (CAS_LINE_NAME); PDF p.122

	// DTP, qualifier 573, "Line Check or Remittance Date", D8; PDF p.124.
	// Any other 2430 DTP in D8 or RD8 form goes in Dates.
	CheckOrRemittanceDateDTP573 DateDTP   `json:"checkOrRemittanceDateDTP573,omitempty"`
	Dates                       []DateDTP `json:"dates,omitempty"`

	// AMT, qualifier EAF, "Remaining Patient Liability"; PDF p.124 (the same
	// qualifier also names "Remaining Patient Liability" at claim-level loop
	// 2320, PDF p.99). Any other 2430 AMT goes in Amounts.
	RemainingLiabilityAMTEAF AmountAMT   `json:"remainingLiabilityAMTEAF,omitempty"`
	Amounts                  []AmountAMT `json:"amounts,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within this loop the model does not place
}

// ServiceLineLX holds one LX loop (2400) and everything this model nests
// under it: the SV2 institutional service line, supplemental information,
// the service date, line identifiers, line amounts, line pricing, the 2410
// drug loop, the four 2420A-D line-provider loops, and the 2430 line
// adjudication loop.
type ServiceLineLX struct {
	Index int `json:"index,omitempty"` // the LX's position in the file

	NumberLX01 string `json:"numberLX01,omitempty"` // Assigned Number, the digits as sent; Transmittal 907, PDF p.111

	Service InstitutionalServiceSV2 `json:"service,omitempty"` // SV2, the Institutional Service Line; Transmittal 907, PDF p.111-112

	// PWK, "Line Supplemental Information". Transmittal 907, PDF p.112-113.
	SupplementalInfo []SupplementalPWK `json:"supplementalInfo,omitempty"`

	// DTP, qualifier 472, "Service Date". Transmittal 907, PDF p.113
	// (printed page 107) lists DTP02 values D8 and RD8: ServiceDateDTP472
	// holds the D8 date, or the RD8 range in Date and DateThrough. Any
	// other line-level DTP in D8 or RD8 form (another qualifier, or a
	// repeat) goes in Dates.
	ServiceDateDTP472 DateDTP   `json:"serviceDateDTP472,omitempty"`
	Dates             []DateDTP `json:"dates,omitempty"`

	// REF (loop 2400). REF_LINE_NAME names these three qualifiers
	// individually; any other REF at this level goes in References.
	// Transmittal 907, PDF p.113-114.
	LineControlREF6R          ReferenceREF   `json:"lineControlREF6R,omitempty"`          // Line Item Control Number
	RepricedLineREF9B         ReferenceREF   `json:"repricedLineREF9B,omitempty"`         // Repriced Line Item Reference Number
	AdjustedRepricedLineREF9D ReferenceREF   `json:"adjustedRepricedLineREF9D,omitempty"` // Adjusted Repriced Line Item Reference Number
	References                []ReferenceREF `json:"references,omitempty"`

	// AMT (loop 2400). AMT01_NAME names these two qualifiers individually
	// at this level; any other AMT goes in Amounts. Transmittal 907, PDF
	// p.114.
	ServiceTaxAMTGT  AmountAMT   `json:"serviceTaxAMTGT,omitempty"`  // Service Tax Amount
	FacilityTaxAMTN8 AmountAMT   `json:"facilityTaxAMTN8,omitempty"` // Facility Tax Amount
	Amounts          []AmountAMT `json:"amounts,omitempty"`

	// HCP, "Line Pricing/Repricing Information". Transmittal 907, PDF
	// p.114-115 (identical shape to the claim-level HCP, PDF p.91-92).
	Pricing PricingHCP `json:"pricing,omitempty"`

	// 2410, drug identification: LIN, CTP, REF. Repeats, though Medicare
	// institutional claims typically send at most one per line.
	Drugs []DrugLIN `json:"drugs,omitempty"`

	// 2420A-D, line-level providers, by role (NM1_ROLE_NAME /
	// LINE_PROVIDER_LOOP_NAME).
	OperatingNM172      LineProviderNM1 `json:"operatingNM172,omitempty"`      // 2420A, Operating Physician
	OtherOperatingNM1ZZ LineProviderNM1 `json:"otherOperatingNM1ZZ,omitempty"` // 2420B, Other Operating Physician
	RenderingNM182      LineProviderNM1 `json:"renderingNM182,omitempty"`      // 2420C, Rendering Provider
	ReferringNM1DN      LineProviderNM1 `json:"referringNM1DN,omitempty"`      // 2420D, Referring Provider

	// 2430, line adjudication information: one per other payer that
	// adjudicated this line (coordination of benefits).
	Adjudications []LineAdjudicationSVD `json:"adjudications,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within loop 2400 the model does not place, including the line-level NTE (Third Party Organization Notes)
}
