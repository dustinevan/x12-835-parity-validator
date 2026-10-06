package x12837p

// Service-line-level (2400) coded elements shared across the segments that
// carry them.

// UnitOrBasisForMeasurementCode is the Unit or Basis for Measurement Code, the
// external X12 355 list shared by every element that carries it (PAT07,
// CR101/CR105, CR302, SV103, SV502, CTP05-1, HCP11).
//
// The list is the union of the codes the 837P implementation guide allows at
// those positions: 01 at PAT07; LB at CR101; DH at CR105; MO at CR302; MJ and
// UN at SV103 and SV502; F2, GR, ME, ML, and UN at CTP05-1; MJ and UN at
// HCP11. Transmittal 579's Values column lists the codes for each position:
// MJ, UN for SV103 (PDF p.79); LB for CR101 and DH for CR105 (p.39, 82); MO
// for CR302 (p.83); F2, GR, ME, ML, UN for CTP05-1 (p.96). CTP05-1 must equal
// UN when the loop 2410 LIN03 NDC is associated with an Oral Cancer Drug HCPCS
// code; CEDI 837P CG, PDF p.47 (printed p.39; Table 41). The names are the
// guide's.
type UnitOrBasisForMeasurementCode string

const (
	UnitOrBasisForMeasurementCodeNotSent UnitOrBasisForMeasurementCode = ""
	UnitOrBasisForMeasurementCodeUN      UnitOrBasisForMeasurementCode = "UN"
	UnitOrBasisForMeasurementCodeMJ      UnitOrBasisForMeasurementCode = "MJ"
	UnitOrBasisForMeasurementCodeLB      UnitOrBasisForMeasurementCode = "LB"
	UnitOrBasisForMeasurementCodeDH      UnitOrBasisForMeasurementCode = "DH"
	UnitOrBasisForMeasurementCodeMO      UnitOrBasisForMeasurementCode = "MO"
	UnitOrBasisForMeasurementCodeF2      UnitOrBasisForMeasurementCode = "F2"
	UnitOrBasisForMeasurementCodeGR      UnitOrBasisForMeasurementCode = "GR"
	UnitOrBasisForMeasurementCodeME      UnitOrBasisForMeasurementCode = "ME"
	UnitOrBasisForMeasurementCodeML      UnitOrBasisForMeasurementCode = "ML"
	UnitOrBasisForMeasurementCode01      UnitOrBasisForMeasurementCode = "01"
)

var unitOrBasisForMeasurementCodeList = codeList{
	{"UN", "Unit", "Unit — CTP05-1 must equal UN when the loop 2410 LIN03 NDC is associated with an Oral Cancer Drug HCPCS code; CEDI 837P CG, PDF p.47 (printed p.39; Table 41); frontend/src/data/x12Names837p.ts CTP05_1_NAME; also listed for SV103, Transmittal 579, PDF p.79"},
	{"MJ", "Minutes", "Minutes — SV103, Transmittal 579, PDF p.79 (Values column, \"MJ, UN\"); the name is the guide's: " + tr3},
	{"LB", "Pound", "Pound — CR101, Transmittal 579, PDF p.39, 82 (Values column); the name is the guide's: " + tr3},
	{"DH", "Miles", "Miles — CR105, Transmittal 579, PDF p.39, 82 (Values column); the name is the guide's: " + tr3},
	{"MO", "Months", "Months — CR302, Transmittal 579, PDF p.83 (Values column); the name is the guide's: " + tr3},
	{"F2", "International Unit", "International Unit — CTP05-1, Transmittal 579, PDF p.96 (Values column, \"F2, GR, ME, ML, UN\"); the name is the guide's: " + tr3},
	{"GR", "Gram", "Gram — CTP05-1, Transmittal 579, PDF p.96 (Values column, \"F2, GR, ME, ML, UN\"); the name is the guide's: " + tr3},
	{"ME", "Milligram", "Milligram — CTP05-1, Transmittal 579, PDF p.96 (Values column, \"F2, GR, ME, ML, UN\"); the name is the guide's: " + tr3},
	{"ML", "Milliliter", "Milliliter — CTP05-1, Transmittal 579, PDF p.96 (Values column, \"F2, GR, ME, ML, UN\"); the name is the guide's: " + tr3},
	{"01", "Actual Pounds", guide("Actual Pounds", "PAT07 at 2000B and 2000C")},
}

func (v UnitOrBasisForMeasurementCode) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v UnitOrBasisForMeasurementCode) Known() bool {
	return unitOrBasisForMeasurementCodeList.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v UnitOrBasisForMeasurementCode) ShortDescription() string {
	return unitOrBasisForMeasurementCodeList.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v UnitOrBasisForMeasurementCode) Description() string {
	return unitOrBasisForMeasurementCodeList.description(string(v))
}

// ParseUnitOrBasisForMeasurementCode returns text as a
// UnitOrBasisForMeasurementCode, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseUnitOrBasisForMeasurementCode(text string) (UnitOrBasisForMeasurementCode, error) {
	return UnitOrBasisForMeasurementCode(text), unitOrBasisForMeasurementCodeList.check("Unit or Basis for Measurement Code", text)
}
