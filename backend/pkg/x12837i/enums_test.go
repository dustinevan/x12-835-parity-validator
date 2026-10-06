package x12837i

import (
	"encoding/json"
	"testing"
)

// TestCodeLists checks, for every coded type, that its list has no empty
// code, no repeated code, and a description for every entry. A short name
// is legitimately "" for a code-named constant.
func TestCodeLists(t *testing.T) {
	tests := []struct {
		name string
		list []codeEntry
	}{
		{"TransactionTypeGS01_ST01", transactionTypeGS01_ST01Codes},
		{"ClaimOrEncounterBHT06", claimOrEncounterBHT06Codes},
		{"HierarchicalLevelHL03", hierarchicalLevelHL03Codes},
		{"PartyRoleNM101", partyRoleNM101Codes},
		{"EntityTypeNM102", entityTypeNM102Codes},
		{"IDQualifierNM108", idQualifierNM108Codes},
		{"DateQualifierDTP01", dateQualifierDTP01Codes},
		{"ProviderRoleCodePRV01", providerRoleCodePRV01Codes},
		{"PayerResponsibilitySBR01", payerResponsibilitySBR01Codes},
		{"RelationshipSBR02", relationshipSBR02Codes},
		{"ClaimFilingIndicatorSBR09", claimFilingIndicatorSBR09Codes},
		{"SexDMG03", sexDMG03Codes},
		{"FacilityCodeQualifierCLM05_2", facilityCodeQualifierCLM05_2Codes},
		{"AssignmentCLM07", assignmentCLM07Codes},
		{"BenefitsAssignCLM08", benefitsAssignCLM08Codes},
		{"ReleaseOfInfoCLM09", releaseOfInfoCLM09Codes},
		{"DelayReasonCLM20", delayReasonCLM20Codes},
		{"PWKReportTypeCode01", pwkReportTypeCode01Codes},
		{"PWKTransmissionCode02", pwkTransmissionCode02Codes},
		{"DiagnosisQualifierHI01_1", diagnosisQualifierHI01_1Codes},
		{"PresentOnAdmissionHI0n_9", presentOnAdmissionHI0n_9Codes},
		{"CRCCategoryCRC01", crcCategoryCRC01Codes},
		{"CRCCertificationCRC02", crcCertificationCRC02Codes},
		{"CRCConditionIndicator", crcConditionIndicatorCodes},
		{"HCPPricingMethodologyHCP01", hcpPricingMethodologyHCP01Codes},
		{"HCPRejectReasonHCP13", hcpRejectReasonHCP13Codes},
		{"HCPPolicyComplianceHCP14", hcpPolicyComplianceHCP14Codes},
		{"HCPExceptionCodeHCP15", hcpExceptionCodeHCP15Codes},
		{"NTENoteCode01", nteNoteCode01Codes},
		{"ProductServiceQualifier", productServiceQualifierCodes},
		{"UnitOrBasisOfMeasurement", unitOrBasisOfMeasurementCodes},
		{"NDCQualifierLIN02", ndcQualifierLIN02Codes},
		{"CASGroupCode01", casGroupCode01Codes},
		{"OtherInsuranceIndicator", otherInsuranceIndicatorCodes},
	}
	for _, tt := range tests {
		if len(tt.list) == 0 {
			t.Errorf("%s: the list is empty", tt.name)
			continue
		}
		seen := map[string]bool{}
		for i, e := range tt.list {
			if e.code == "" || seen[e.code] {
				t.Errorf("%s[%d]: code %q is empty or repeated", tt.name, i, e.code)
			}
			seen[e.code] = true
			if e.description == "" {
				t.Errorf("%s[%d] (%s): description is empty", tt.name, i, e.code)
			}
		}
	}
}

// TestCodedValues checks every coded type's Parse and Known: "" is the one
// empty and parses clean, a code on the list parses clean and is known, and
// a code off the list comes back as itself, unknown, with an error.
func TestCodedValues(t *testing.T) {
	type parsed struct {
		text  string
		known bool
		err   error
	}
	tests := []struct {
		name  string
		known string
		parse func(string) parsed
	}{
		{"TransactionTypeGS01_ST01", "HC", func(s string) parsed {
			v, err := ParseTransactionTypeGS01_ST01(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"ClaimOrEncounterBHT06", "CH", func(s string) parsed {
			v, err := ParseClaimOrEncounterBHT06(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"HierarchicalLevelHL03", "20", func(s string) parsed {
			v, err := ParseHierarchicalLevelHL03(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"PartyRoleNM101", "85", func(s string) parsed { v, err := ParsePartyRoleNM101(s); return parsed{v.String(), v.Known(), err} }},
		{"EntityTypeNM102", "1", func(s string) parsed { v, err := ParseEntityTypeNM102(s); return parsed{v.String(), v.Known(), err} }},
		{"IDQualifierNM108", "XX", func(s string) parsed { v, err := ParseIDQualifierNM108(s); return parsed{v.String(), v.Known(), err} }},
		{"DateQualifierDTP01", "434", func(s string) parsed { v, err := ParseDateQualifierDTP01(s); return parsed{v.String(), v.Known(), err} }},
		{"ProviderRoleCodePRV01", "BI", func(s string) parsed {
			v, err := ParseProviderRoleCodePRV01(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"PayerResponsibilitySBR01", "P", func(s string) parsed {
			v, err := ParsePayerResponsibilitySBR01(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"RelationshipSBR02", "18", func(s string) parsed { v, err := ParseRelationshipSBR02(s); return parsed{v.String(), v.Known(), err} }},
		{"ClaimFilingIndicatorSBR09", "MA", func(s string) parsed {
			v, err := ParseClaimFilingIndicatorSBR09(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"SexDMG03", "F", func(s string) parsed { v, err := ParseSexDMG03(s); return parsed{v.String(), v.Known(), err} }},
		{"FacilityCodeQualifierCLM05_2", "A", func(s string) parsed {
			v, err := ParseFacilityCodeQualifierCLM05_2(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"AssignmentCLM07", "A", func(s string) parsed { v, err := ParseAssignmentCLM07(s); return parsed{v.String(), v.Known(), err} }},
		{"BenefitsAssignCLM08", "Y", func(s string) parsed {
			v, err := ParseBenefitsAssignCLM08(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"ReleaseOfInfoCLM09", "Y", func(s string) parsed { v, err := ParseReleaseOfInfoCLM09(s); return parsed{v.String(), v.Known(), err} }},
		{"DelayReasonCLM20", "1", func(s string) parsed { v, err := ParseDelayReasonCLM20(s); return parsed{v.String(), v.Known(), err} }},
		{"PWKReportTypeCode01", "OZ", func(s string) parsed {
			v, err := ParsePWKReportTypeCode01(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"PWKTransmissionCode02", "AA", func(s string) parsed {
			v, err := ParsePWKTransmissionCode02(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"DiagnosisQualifierHI01_1", "ABK", func(s string) parsed {
			v, err := ParseDiagnosisQualifierHI01_1(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"PresentOnAdmissionHI0n_9", "Y", func(s string) parsed {
			v, err := ParsePresentOnAdmissionHI0n_9(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"CRCCategoryCRC01", "ZZ", func(s string) parsed { v, err := ParseCRCCategoryCRC01(s); return parsed{v.String(), v.Known(), err} }},
		{"CRCCertificationCRC02", "Y", func(s string) parsed {
			v, err := ParseCRCCertificationCRC02(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"CRCConditionIndicator", "AV", func(s string) parsed {
			v, err := ParseCRCConditionIndicator(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"HCPPricingMethodologyHCP01", "00", func(s string) parsed {
			v, err := ParseHCPPricingMethodologyHCP01(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"HCPRejectReasonHCP13", "T1", func(s string) parsed {
			v, err := ParseHCPRejectReasonHCP13(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"HCPPolicyComplianceHCP14", "1", func(s string) parsed {
			v, err := ParseHCPPolicyComplianceHCP14(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"HCPExceptionCodeHCP15", "1", func(s string) parsed {
			v, err := ParseHCPExceptionCodeHCP15(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"NTENoteCode01", "ADD", func(s string) parsed { v, err := ParseNTENoteCode01(s); return parsed{v.String(), v.Known(), err} }},
		{"ProductServiceQualifier", "HC", func(s string) parsed {
			v, err := ParseProductServiceQualifier(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"UnitOrBasisOfMeasurement", "UN", func(s string) parsed {
			v, err := ParseUnitOrBasisOfMeasurement(s)
			return parsed{v.String(), v.Known(), err}
		}},
		{"NDCQualifierLIN02", "N4", func(s string) parsed { v, err := ParseNDCQualifierLIN02(s); return parsed{v.String(), v.Known(), err} }},
		{"CASGroupCode01", "CO", func(s string) parsed { v, err := ParseCASGroupCode01(s); return parsed{v.String(), v.Known(), err} }},
		{"OtherInsuranceIndicator", "Y", func(s string) parsed {
			v, err := ParseOtherInsuranceIndicator(s)
			return parsed{v.String(), v.Known(), err}
		}},
	}
	for _, tt := range tests {
		if got := tt.parse(""); got.text != "" || got.known || got.err != nil {
			t.Errorf("%s: Parse(\"\") = %+v, want \"\", not known, nil", tt.name, got)
		}
		if got := tt.parse(tt.known); got.text != tt.known || !got.known || got.err != nil {
			t.Errorf("%s: Parse(%q) = %+v, want %q, known, nil", tt.name, tt.known, got, tt.known)
		}
		if got := tt.parse("ZQ9X"); got.text != "ZQ9X" || got.known || got.err == nil {
			t.Errorf("%s: Parse(\"ZQ9X\") = %+v, want \"ZQ9X\", not known, an error", tt.name, got)
		}
	}
}

// TestCodeDescriptions checks the two description methods: a named code,
// a code-named code (no short name), and a code off the list (nothing).
func TestCodeDescriptions(t *testing.T) {
	tests := []struct {
		name      string
		short     string
		described bool
	}{
		{"named", CASGroupCode01CO.ShortDescription(), CASGroupCode01CO.Description() != ""},
		{"code-named", ClaimOrEncounterBHT06CH.ShortDescription(), ClaimOrEncounterBHT06CH.Description() != ""},
		{"off the list", CASGroupCode01("XX").ShortDescription(), CASGroupCode01("XX").Description() != ""},
	}
	wantShort := []string{"Contractual Obligation", "", ""}
	wantDescribed := []bool{true, true, false}
	for i, tt := range tests {
		if tt.short != wantShort[i] || tt.described != wantDescribed[i] {
			t.Errorf("%s: ShortDescription = %q, Description present = %v, want %q, %v", tt.name, tt.short, tt.described, wantShort[i], wantDescribed[i])
		}
	}
}

// TestCodedJSON checks that a coded value, on the list or off it, goes
// through JSON as its own text and comes back the same.
func TestCodedJSON(t *testing.T) {
	type holder struct {
		Group CASGroupCode01 `json:"group,omitempty"`
		Sex   SexDMG03       `json:"sex,omitempty"`
	}
	in := holder{Group: CASGroupCode01("XX"), Sex: SexDMG03F}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"group":"XX","sex":"F"}`; string(data) != want {
		t.Errorf("Marshal = %s, want %s", data, want)
	}
	var out holder
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("Unmarshal = %+v, want %+v", out, in)
	}
	var empty holder
	data, err = json.Marshal(empty)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{}`; string(data) != want {
		t.Errorf("Marshal of empty = %s, want %s", data, want)
	}
}
