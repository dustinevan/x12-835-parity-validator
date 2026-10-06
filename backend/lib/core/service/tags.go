package service

import (
	"regexp"
	"strings"
)

// Tag is a differing element as the views show it: the element label the
// comparison writes (the tag) and a plain label.
type Tag struct {
	Tag   string `json:"tag"`
	Label string `json:"label"`
}

// tagLabels gives the common tags a plain label.
var tagLabels = map[string]string{
	"CLP02":      "Claim status",
	"CLP03":      "Billed charge",
	"CLP04":      "Claim payment",
	"CLP05":      "Patient responsibility",
	"CLP06":      "Claim filing indicator",
	"CLP08":      "Facility type",
	"CLP09":      "Claim frequency",
	"CLP11":      "Diagnosis-related group",
	"CLP12":      "DRG weight",
	"CLP13":      "Discharge fraction",
	"MIA":        "Inpatient remark codes",
	"MIA04":      "Inpatient DRG amount",
	"MOA":        "Outpatient remark codes",
	"LQ":         "Line remark codes",
	"AMT*AU":     "Claim coverage amount",
	"AMT*I":      "Interest",
	"AMT*B6":     "Line allowed amount",
	"DTM*050":    "Claim received date",
	"DTM*232":    "Statement start date",
	"DTM*233":    "Statement end date",
	"DTM*472":    "Line service date",
	"NM1*QC":     "Patient identifier",
	"NM1*74":     "Corrected insured",
	"NM1*82":     "Rendering provider",
	"NM1*TT":     "Crossover carrier",
	"NM1*PR":     "Corrected priority payer",
	"REF*6R":     "Line control number",
	"SVC01":      "Procedure code",
	"SVC02":      "Line charge",
	"SVC03":      "Line payment",
	"SVC04":      "Revenue code",
	"SVC05":      "Paid units",
	"CAS*CO*45":  "Charge above the fee schedule",
	"CAS*CO*59":  "Multiple-procedure reduction",
	"CAS*CO*97":  "Bundled into another service",
	"CAS*CO*119": "Benefit maximum reached",
	"CAS*CO*253": "Sequestration adjustment",
	"CAS*PR*1":   "Deductible",
	"CAS*PR*2":   "Coinsurance",
	"CAS*PR*3":   "Copayment",
}

// TagOf describes a tag. An unplaced segment (RAW*K3) is described by its
// segment ID, and a CAS quantity (CAS*CO*45*QTY) by its adjustment. A tag
// outside the table is its own label.
func TagOf(tag string) Tag {
	if label, ok := tagLabels[tag]; ok {
		return Tag{Tag: tag, Label: label}
	}
	if id, ok := strings.CutPrefix(tag, "RAW*"); ok {
		return Tag{Tag: tag, Label: "Unplaced segment " + id}
	}
	if adjustment, ok := strings.CutSuffix(tag, "*QTY"); ok && strings.HasPrefix(adjustment, "CAS*") {
		return Tag{Tag: tag, Label: TagOf(adjustment).Label + ", quantity"}
	}
	return Tag{Tag: tag, Label: tag}
}

// tagPattern is the shape of every element label the comparison writes: an
// unplaced segment by its ID (RAW*K3), a CAS adjustment's quantity by its
// group and reason (CAS*CO*253*QTY), and otherwise a segment tag with an
// element number (CLP04), a segment with qualifiers (CAS*CO*253, NM1*TT),
// or a bare segment (MOA).
var tagPattern = regexp.MustCompile(`^(RAW\*[A-Z0-9]{2,3}|CAS\*[A-Z0-9]{1,5}\*[A-Z0-9]{1,5}\*QTY|[A-Z][A-Z0-9]{1,4}(\*[A-Z0-9]{1,5}){0,2})$`)

// ValidTag reports whether s has the shape of a tag. It does not check
// that any claim carries it.
func ValidTag(s string) bool {
	return tagPattern.MatchString(s)
}
