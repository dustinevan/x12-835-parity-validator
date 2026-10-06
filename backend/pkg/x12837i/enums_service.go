package x12837i

import "fmt"

// ProductServiceQualifier is the Product or Service ID Qualifier carried by
// SV202-1, SVD03-1, and HCP09 — the same code set reused at all three
// positions (frontend/src/data/x12Names837i.ts's SV202_1_VALUES,
// SVD03_1_VALUES, HCP09_VALUES). Transmittal 907 gives only this bare value
// list (PDF p.112 SV202-1, p.122 SVD03-1, p.91/114-115 HCP09); no individual
// code's meaning is given in the reviewed sources, so every value here is
// named by its code, following the coded-element scheme in enums.go.
type ProductServiceQualifier string

const (
	ProductServiceQualifierER ProductServiceQualifier = "ER"
	ProductServiceQualifierHC ProductServiceQualifier = "HC"
	ProductServiceQualifierHP ProductServiceQualifier = "HP"
	ProductServiceQualifierIV ProductServiceQualifier = "IV"
	ProductServiceQualifierWK ProductServiceQualifier = "WK"
)

// The codes the guide allows for product/service ID qualifier, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var productServiceQualifierCodes = []codeEntry{
	{"ER", "", "Transmittal 907, PDF p.112 (SV202-1), p.122 (SVD03-1), p.91/114-115 (HCP09); code list only, no per-code text given in the flat file"},
	{"HC", "", "Transmittal 907, PDF p.112 (SV202-1), p.122 (SVD03-1), p.91/114-115 (HCP09); code list only, no per-code text given in the flat file"},
	{"HP", "", "Transmittal 907, PDF p.112 (SV202-1), p.122 (SVD03-1), p.91/114-115 (HCP09); code list only, no per-code text given in the flat file"},
	{"IV", "", "Transmittal 907, PDF p.112 (SV202-1), p.122 (SVD03-1), p.91/114-115 (HCP09); code list only, no per-code text given in the flat file"},
	{"WK", "", "Transmittal 907, PDF p.112 (SV202-1), p.122 (SVD03-1), p.91/114-115 (HCP09); code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v ProductServiceQualifier) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v ProductServiceQualifier) Known() bool {
	return knownCode(productServiceQualifierCodes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v ProductServiceQualifier) ShortDescription() string {
	return shortOf(productServiceQualifierCodes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v ProductServiceQualifier) Description() string {
	return descriptionOf(productServiceQualifierCodes, string(v))
}

// ParseProductServiceQualifier reads one product/service ID qualifier value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseProductServiceQualifier(code string) (ProductServiceQualifier, error) {
	v := ProductServiceQualifier(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing product/service ID qualifier: a value not on the list")
}

// UnitOrBasisOfMeasurement is the X12 Unit or Basis for Measurement Code
// (element 355) wherever it occurs with a documented value list in this
// model: HCP11 (frontend/src/data/x12Names837i.ts's HCP11_VALUES, "Days/
// Units", Transmittal 907 PDF p.91, identical shape at line-level p.114-115)
// and CTP05-1 (CTP05_1_VALUES, PDF p.116). It is also used for SV2's own
// SV204, which Transmittal 907 gives no separate value list for; that reuse
// is this model's inference, not an independently cited SV204 list.
type UnitOrBasisOfMeasurement string

const (
	UnitOrBasisOfMeasurementDA UnitOrBasisOfMeasurement = "DA"
	UnitOrBasisOfMeasurementUN UnitOrBasisOfMeasurement = "UN"
	UnitOrBasisOfMeasurementF2 UnitOrBasisOfMeasurement = "F2"
	UnitOrBasisOfMeasurementGR UnitOrBasisOfMeasurement = "GR"
	UnitOrBasisOfMeasurementME UnitOrBasisOfMeasurement = "ME"
	UnitOrBasisOfMeasurementML UnitOrBasisOfMeasurement = "ML"
)

// The codes the guide allows for unit or basis for measurement code, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var unitOrBasisOfMeasurementCodes = []codeEntry{
	{"DA", "", "Transmittal 907, PDF p.91 (HCP11, \"Days/Units\"), identical shape at line-level p.114-115; code list only, no per-code text given in the flat file"},
	{"UN", "", "Transmittal 907, PDF p.91 (HCP11); also PDF p.116 (CTP05-1); code list only, no per-code text given in the flat file"},
	{"F2", "", "Transmittal 907, PDF p.116 (CTP05-1, Composite Unit of Measure); code list only, no per-code text given in the flat file"},
	{"GR", "", "Transmittal 907, PDF p.116 (CTP05-1); code list only, no per-code text given in the flat file"},
	{"ME", "", "Transmittal 907, PDF p.116 (CTP05-1); code list only, no per-code text given in the flat file"},
	{"ML", "", "Transmittal 907, PDF p.116 (CTP05-1); code list only, no per-code text given in the flat file"},
}

// String returns the code as the file sent it.
func (v UnitOrBasisOfMeasurement) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v UnitOrBasisOfMeasurement) Known() bool {
	return knownCode(unitOrBasisOfMeasurementCodes, string(v))
}

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v UnitOrBasisOfMeasurement) ShortDescription() string {
	return shortOf(unitOrBasisOfMeasurementCodes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v UnitOrBasisOfMeasurement) Description() string {
	return descriptionOf(unitOrBasisOfMeasurementCodes, string(v))
}

// ParseUnitOrBasisOfMeasurement reads one unit or basis for measurement code value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseUnitOrBasisOfMeasurement(code string) (UnitOrBasisOfMeasurement, error) {
	v := UnitOrBasisOfMeasurement(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing unit or basis for measurement code: a value not on the list")
}

// NDCQualifierLIN02 is LIN02, the Product or Service ID Qualifier on the
// drug identification segment (loop 2410).
type NDCQualifierLIN02 string

const (
	NDCQualifierLIN02N4 NDCQualifierLIN02 = "N4"
)

// The codes the guide allows for LIN02, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var ndcQualifierLIN02Codes = []codeEntry{
	{"N4", "National Drug Code in 5-4-2 Format", "National Drug Code in 5-4-2 Format — Transmittal 907, PDF p.115 (LIN02_NAME)"},
}

// String returns the code as the file sent it.
func (v NDCQualifierLIN02) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v NDCQualifierLIN02) Known() bool { return knownCode(ndcQualifierLIN02Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v NDCQualifierLIN02) ShortDescription() string {
	return shortOf(ndcQualifierLIN02Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v NDCQualifierLIN02) Description() string {
	return descriptionOf(ndcQualifierLIN02Codes, string(v))
}

// ParseNDCQualifierLIN02 reads one LIN02 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseNDCQualifierLIN02(code string) (NDCQualifierLIN02, error) {
	v := NDCQualifierLIN02(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing LIN02: a value not on the list")
}
