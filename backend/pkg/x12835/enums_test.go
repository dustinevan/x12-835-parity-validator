package x12835

import "testing"

// codeLists is every coded element's list, by type name, for the checks
// below.
var codeLists = map[string][]codeEntry{
	"AuthQualifierISA01":          authQualifierISA01Codes,
	"SecQualifierISA03":           secQualifierISA03Codes,
	"InterchangeIDQualifierISA05": interchangeIDQualifierISA05Codes,
	"AckRequestedISA14":           ackRequestedISA14Codes,
	"UsageISA15":                  usageISA15Codes,
	"FunctionalIDGS01":            functionalIDGS01Codes,
	"AgencyGS07":                  agencyGS07Codes,
	"TransactionSetIDST01":        transactionSetIDST01Codes,
	"HandlingBPR01":               handlingBPR01Codes,
	"CreditDebitBPR03":            creditDebitBPR03Codes,
	"PayMethodBPR04":              payMethodBPR04Codes,
	"PayFormatBPR05":              payFormatBPR05Codes,
	"DFIQualifierBPR06":           dfiQualifierBPR06Codes,
	"AccountQualifierBPR08":       accountQualifierBPR08Codes,
	"TraceTypeTRN01":              traceTypeTRN01Codes,
	"EntityIDN101":                entityIDN101Codes,
	"IDQualifierN103":             idQualifierN103Codes,
	"ContactFunctionPER01":        contactFunctionPER01Codes,
	"CommQualifierPER03":          commQualifierPER03Codes,
	"DeliveryMethodRDM01":         deliveryMethodRDM01Codes,
	"ClaimStatusCLP02":            claimStatusCLP02Codes,
	"FilingIndicatorCLP06":        filingIndicatorCLP06Codes,
	"EntityIDNM101":               entityIDNM101Codes,
	"EntityTypeNM102":             entityTypeNM102Codes,
	"IDQualifierNM108":            idQualifierNM108Codes,
	"RefQualifierREF01":           refQualifierREF01Codes,
	"DateQualifierDTM01":          dateQualifierDTM01Codes,
	"AmountQualifierAMT01":        amountQualifierAMT01Codes,
	"QuantityQualifierQTY01":      quantityQualifierQTY01Codes,
	"AdjGroupCAS01":               adjGroupCAS01Codes,
	"PLBReasonPLB03":              plbReasonPLB03Codes,
	"ProductQualifierSVC01":       productQualifierSVC01Codes,
	"RemarkQualifierLQ01":         remarkQualifierLQ01Codes,
}

// TestCodeLists checks every list: each code is non-empty and appears
// once, and each entry has a description (a citation at the least). A
// short name may be "" when no reviewed source names the code.
func TestCodeLists(t *testing.T) {
	for name, list := range codeLists {
		if len(list) == 0 {
			t.Errorf("%s: no codes", name)
		}
		seen := map[string]bool{}
		for i, e := range list {
			if e.code == "" || seen[e.code] {
				t.Errorf("%s[%d]: code %q is empty or repeated", name, i, e.code)
			}
			seen[e.code] = true
			if e.description == "" {
				t.Errorf("%s[%d] (%s): description is empty", name, i, e.code)
			}
		}
	}
}

// TestCodeMethods checks the methods on one type: a listed code is known
// with its name and citation, an empty value is not sent and not known, and
// a value off the list is kept as it came, unknown, with no description.
func TestCodeMethods(t *testing.T) {
	tests := []struct {
		value     AdjGroupCAS01
		known     bool
		wantShort string
	}{
		{AdjGroupCAS01NotSent, false, ""},
		{AdjGroupCAS01CO, true, "Contractual Obligation"},
		{AdjGroupCAS01PI, true, "Payor Initiated Reductions"},
		{AdjGroupCAS01("ZZ"), false, ""},
	}
	for _, tt := range tests {
		if got := tt.value.Known(); got != tt.known {
			t.Errorf("AdjGroupCAS01(%q).Known() = %v, want %v", tt.value, got, tt.known)
		}
		if got := tt.value.ShortDescription(); got != tt.wantShort {
			t.Errorf("AdjGroupCAS01(%q).ShortDescription() = %q, want %q", tt.value, got, tt.wantShort)
		}
		if got := tt.value.Description(); (got != "") != tt.known {
			t.Errorf("AdjGroupCAS01(%q).Description() = %q, want text only when known", tt.value, got)
		}
		if got := tt.value.String(); got != string(tt.value) {
			t.Errorf("AdjGroupCAS01(%q).String() = %q", tt.value, got)
		}
	}
	if got := ClaimStatusCLP021.ShortDescription(); got != "Processed as Primary" {
		t.Errorf("ClaimStatusCLP021.ShortDescription() = %q, want the guide's name", got)
	}
}

// TestParseCodes checks that a Parse function keeps the value whether or
// not it is on the list, and reports an off-list value with an error that
// names the element.
func TestParseCodes(t *testing.T) {
	tests := []struct {
		code    string
		want    ClaimStatusCLP02
		wantErr bool
	}{
		{"", ClaimStatusCLP02NotSent, false},
		{"1", ClaimStatusCLP021, false},
		{"25", ClaimStatusCLP0225, false},
		{"99", ClaimStatusCLP02("99"), true},
	}
	for _, tt := range tests {
		got, err := ParseClaimStatusCLP02(tt.code)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("ParseClaimStatusCLP02(%q) = %q, %v, want %q, error %v", tt.code, got, err, tt.want, tt.wantErr)
		}
		if err != nil && err.Error() != `x12835: parsing CLP02: a value not on the list` {
			t.Errorf("ParseClaimStatusCLP02(%q) error = %q", tt.code, err)
		}
	}
}
