package x12837p

import (
	"encoding/json"
	"strings"
	"testing"
)

// codeLists names every coded element's list, so one test can check them
// all the same way.
var codeLists = map[string]codeList{
	"ClaimAdjustmentGroupCodeCAS01":         claimAdjustmentGroupCodeCAS01List,
	"EntityIdentifierNM101":                 entityIdentifierNM101List,
	"EntityTypeQualifierNM102":              entityTypeQualifierNM102List,
	"IdentificationCodeQualifierNM108":      identificationCodeQualifierNM108List,
	"ClaimFrequencyCLM05_3":                 claimFrequencyCLM05_3List,
	"FacilityCodeQualifierCLM05_2":          facilityCodeQualifierCLM05_2List,
	"ProviderOrSupplierSignatureCLM06":      providerOrSupplierSignatureCLM06List,
	"AssignmentOrPlanParticipationCLM07":    assignmentOrPlanParticipationCLM07List,
	"BenefitsAssignmentCertificationCLM08":  benefitsAssignmentCertificationCLM08List,
	"ReleaseOfInformationCLM09":             releaseOfInformationCLM09List,
	"PatientSignatureSourceCLM10":           patientSignatureSourceCLM10List,
	"RelatedCausesCodeCLM11_1":              relatedCausesCodeCLM11_1List,
	"DelayReasonCLM20":                      delayReasonCLM20List,
	"AttachmentReportTypeCodePWK01":         attachmentReportTypeCodePWK01List,
	"AttachmentTransmissionCodePWK02":       attachmentTransmissionCodePWK02List,
	"DiagnosisTypeCodeHI0x_1":               diagnosisTypeCodeHI0x_1List,
	"CodeCategoryCRC01":                     codeCategoryCRC01List,
	"ConditionIndicatorCRC0x":               conditionIndicatorCRC0xList,
	"YesNoConditionResponseCode":            yesNoConditionResponseCodeList,
	"PricingMethodologyHCP01":               pricingMethodologyHCP01List,
	"ProviderTypePRV01":                     providerTypePRV01List,
	"ContactFunctionCodePER01":              contactFunctionCodePER01List,
	"CommunicationNumberQualifierPER03":     communicationNumberQualifierPER03List,
	"UnitOrBasisForMeasurementCode":         unitOrBasisForMeasurementCodeList,
	"PayerResponsibilitySBR01":              payerResponsibilitySBR01List,
	"RelationshipSBR02":                     relationshipSBR02List,
	"RelationshipPAT01":                     relationshipPAT01List,
	"ClaimFilingIndicatorSBR09":             claimFilingIndicatorSBR09List,
	"InsuranceTypeSBR05":                    insuranceTypeSBR05List,
	"ReferenceIdentificationQualifierREF01": referenceIdentificationQualifierREF01List,
	"DateTimeQualifierDTP01":                dateTimeQualifierDTP01List,
	"AmountQualifierAMT01":                  amountQualifierAMT01List,
	"QuantityQualifierQTY01":                quantityQualifierQTY01List,
	"HierarchicalLevelHL03":                 hierarchicalLevelHL03List,
	"HierarchicalChildCodeHL04":             hierarchicalChildCodeHL04List,
	"FunctionalIDGS01":                      functionalIDGS01List,
	"AgencyGS07":                            agencyGS07List,
	"HierarchicalStructureBHT01":            hierarchicalStructureBHT01List,
	"TransactionSetPurposeBHT02":            transactionSetPurposeBHT02List,
	"TransactionTypeBHT06":                  transactionTypeBHT06List,
	"AuthQualifierISA01":                    authQualifierISA01List,
	"SecQualifierISA03":                     secQualifierISA03List,
	"InterchangeIDQualifierISA05":           interchangeIDQualifierISA05List,
	"AckRequestedISA14":                     ackRequestedISA14List,
	"UsageISA15":                            usageISA15List,
	"TransactionSetIDST01":                  transactionSetIDST01List,
}

// TestCodeLists checks every list: it has at least one code, no code is
// empty or repeated, every entry has a description, and a description
// starts with the entry's short name when it has one.
func TestCodeLists(t *testing.T) {
	if len(codeLists) != 46 {
		t.Errorf("codeLists has %d lists, want 46", len(codeLists))
	}
	for name, list := range codeLists {
		if len(list) == 0 {
			t.Errorf("%s: the list is empty", name)
			continue
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
			if e.short != "" && !strings.HasPrefix(e.description, e.short) {
				t.Errorf("%s[%d] (%s): description %q does not start with the short name %q", name, i, e.code, e.description, e.short)
			}
		}
	}
}

// TestEnumMethods checks the methods on one string enum for the three
// kinds of value: not sent, on the list, and off the list.
func TestEnumMethods(t *testing.T) {
	tests := []struct {
		value     EntityIdentifierNM101
		known     bool
		wantShort string
		wantDesc  bool // whether Description returns text
	}{
		{EntityIdentifierNM101NotSent, false, "not sent", true},
		{EntityIdentifierNM10185, true, "Billing Provider", true},
		{EntityIdentifierNM101("ZZ"), false, "", false},
	}
	for _, tt := range tests {
		if got := tt.value.String(); got != string(tt.value) {
			t.Errorf("EntityIdentifierNM101(%q).String() = %q, want the value itself", tt.value, got)
		}
		if got := tt.value.Known(); got != tt.known {
			t.Errorf("EntityIdentifierNM101(%q).Known() = %v, want %v", tt.value, got, tt.known)
		}
		if got := tt.value.ShortDescription(); got != tt.wantShort {
			t.Errorf("EntityIdentifierNM101(%q).ShortDescription() = %q, want %q", tt.value, got, tt.wantShort)
		}
		if got := tt.value.Description(); (got != "") != tt.wantDesc {
			t.Errorf("EntityIdentifierNM101(%q).Description() = %q, want text: %v", tt.value, got, tt.wantDesc)
		}
	}
}

// TestEnumParse checks Parse and the JSON round trip on one enum: "" and a
// listed code parse without error, an off-list code parses with an error
// but keeps its text, and the JSON form is the code text in every case.
func TestEnumParse(t *testing.T) {
	tests := []struct {
		text    string
		wantErr bool
	}{
		{"", false},
		{"22", false},
		{"19", true},
	}
	for _, tt := range tests {
		got, err := ParseHierarchicalLevelHL03(tt.text)
		if string(got) != tt.text || (err != nil) != tt.wantErr {
			t.Errorf("ParseHierarchicalLevelHL03(%q) = %q, %v, want %q, error %v", tt.text, got, err, tt.text, tt.wantErr)
		}
		data, err := json.Marshal(got)
		if err != nil || string(data) != `"`+tt.text+`"` {
			t.Errorf("json.Marshal(HierarchicalLevelHL03(%q)) = %s, %v, want the quoted text", tt.text, data, err)
		}
		var back HierarchicalLevelHL03
		if err := json.Unmarshal(data, &back); err != nil || back != got {
			t.Errorf("json.Unmarshal(%s) = %q, %v, want %q, nil", data, back, err, got)
		}
	}
}
