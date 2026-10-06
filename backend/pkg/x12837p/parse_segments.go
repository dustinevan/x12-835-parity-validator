package x12837p

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/carc"
	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/rarc"
)

// Each parse function below reads one segment's elements into its typed
// struct. A coded element is read as its text (the enum types are string
// types; see enums.go), so a code off the list is kept exactly as sent and
// never fails the segment. Elements after the last position the model
// names go to the struct's Extra, and components after the last a
// composite names go to that composite's own extra field, so nothing is
// dropped. An array indexed by element number collects what each numeric
// element's parse reports: a malformed amount, date, time, count, or
// composite fails the segment, and errors.Join reports them all. Each
// struct records the segment's position in the file in Index.

func parseISA(s inputSegment) (InterchangeISA, error) {
	var isa InterchangeISA
	var errs [17]error
	isa.AuthQualifierISA01 = AuthQualifierISA01(s.el(1))
	isa.AuthInfoISA02 = s.el(2)
	isa.SecQualifierISA03 = SecQualifierISA03(s.el(3))
	isa.SecInfoISA04 = s.el(4)
	isa.SenderIDQualifierISA05 = InterchangeIDQualifierISA05(s.el(5))
	isa.SenderIDISA06 = s.el(6)
	isa.ReceiverIDQualifierISA07 = InterchangeIDQualifierISA05(s.el(7))
	isa.ReceiverIDISA08 = s.el(8)
	isa.DateISA09, errs[9] = parseInterchangeDate(s, 9)
	isa.TimeISA10, errs[10] = parseTime(s, 10)
	isa.RepetitionSepISA11 = s.el(11)
	isa.VersionISA12 = s.el(12)
	isa.ControlNumberISA13 = s.el(13)
	isa.AckRequestedISA14 = AckRequestedISA14(s.el(14))
	isa.UsageISA15 = UsageISA15(s.el(15))
	isa.ComponentSepISA16 = s.el(16)
	if err := errors.Join(errs[:]...); err != nil {
		return InterchangeISA{}, err
	}
	return isa, nil
}

func parseGS(s inputSegment) (FunctionalGroupGS, error) {
	g := FunctionalGroupGS{Index: s.index}
	var errs [9]error
	g.FunctionalIDGS01 = FunctionalIDGS01(s.el(1))
	g.SenderCodeGS02 = s.el(2)
	g.ReceiverCodeGS03 = s.el(3)
	g.DateGS04, errs[4] = parseDate(s, 4)
	g.TimeGS05, errs[5] = parseTime(s, 5)
	g.ControlNumberGS06 = s.el(6)
	g.AgencyGS07 = AgencyGS07(s.el(7))
	g.VersionGS08 = s.el(8)
	g.Extra = s.extra(8)
	if err := errors.Join(errs[:]...); err != nil {
		return FunctionalGroupGS{}, err
	}
	if !g.FunctionalIDGS01.Known() {
		return FunctionalGroupGS{}, fmt.Errorf("x12837p: parsing GS01 in segment %d: %q is not a health care claim functional group", s.index, g.FunctionalIDGS01)
	}
	return g, nil
}

func parseST(s inputSegment) (TransactionST, error) {
	t := TransactionST{Index: s.index}
	t.TransactionSetIDST01 = TransactionSetIDST01(s.el(1))
	t.ControlNumberST02 = s.el(2)
	t.ImplementationGuideVersionST03 = s.el(3)
	t.Extra = s.extra(3)
	if !t.TransactionSetIDST01.Known() {
		return TransactionST{}, fmt.Errorf("x12837p: parsing ST01 in segment %d: %q is not an 837 transaction set", s.index, t.TransactionSetIDST01)
	}
	return t, nil
}

// parseTrailer reads an SE, GE, or IEA: a count in element 01 and a control
// number in element 02.
func parseTrailer(s inputSegment) (int, string, error) {
	count, err := parseWholeNumber(s, 1)
	if err != nil {
		return 0, "", err
	}
	return count, s.el(2), nil
}

func parseBHT(s inputSegment) (HeaderBHT, error) {
	h := HeaderBHT{Set: true, Index: s.index}
	var errs [7]error
	h.StructureBHT01 = HierarchicalStructureBHT01(s.el(1))
	h.PurposeBHT02 = TransactionSetPurposeBHT02(s.el(2))
	h.OriginatorTransactionIDBHT03 = s.el(3)
	h.DateBHT04, errs[4] = parseDate(s, 4)
	h.TimeBHT05, errs[5] = parseTime(s, 5)
	h.TransactionTypeBHT06 = TransactionTypeBHT06(s.el(6))
	h.Extra = s.extra(6)
	if err := errors.Join(errs[:]...); err != nil {
		return HeaderBHT{}, err
	}
	return h, nil
}

func parseHL(s inputSegment) HierarchyHL {
	return HierarchyHL{
		Index:              s.index,
		IDNumberHL01:       s.el(1),
		ParentIDNumberHL02: s.el(2),
		LevelHL03:          HierarchicalLevelHL03(s.el(3)),
		ChildCodeHL04:      HierarchicalChildCodeHL04(s.el(4)),
		Extra:              s.extra(4),
	}
}

func parseNM1(s inputSegment) NameNM1 {
	return NameNM1{
		Index:            s.index,
		EntityIDNM101:    EntityIdentifierNM101(s.el(1)),
		EntityTypeNM102:  EntityTypeQualifierNM102(s.el(2)),
		LastNM103:        s.el(3),
		FirstNM104:       s.el(4),
		MiddleNM105:      s.el(5),
		PrefixNM106:      s.el(6),
		SuffixNM107:      s.el(7),
		IDQualifierNM108: IdentificationCodeQualifierNM108(s.el(8)),
		IDNM109:          s.el(9),
		Extra:            s.extra(9),
	}
}

func parseN3(s inputSegment) AddressN3 {
	return AddressN3{Set: true, Index: s.index, Line1N301: s.el(1), Line2N302: s.el(2), Extra: s.extra(2)}
}

func parseN4(s inputSegment) LocationN4 {
	return LocationN4{
		Set:                    true,
		Index:                  s.index,
		CityN401:               s.el(1),
		StateN402:              s.el(2),
		ZipN403:                s.el(3),
		CountryN404:            s.el(4),
		LocationQualifierN405:  s.el(5),
		LocationIdentifierN406: s.el(6),
		SubdivisionN407:        s.el(7),
		Extra:                  s.extra(7),
	}
}

func parsePER(s inputSegment) ContactPER {
	return ContactPER{
		Index:                 s.index,
		FunctionPER01:         ContactFunctionCodePER01(s.el(1)),
		NamePER02:             s.el(2),
		QualifierPER03:        CommunicationNumberQualifierPER03(s.el(3)),
		CommPER04:             s.el(4),
		QualifierPER05:        CommunicationNumberQualifierPER03(s.el(5)),
		CommPER06:             s.el(6),
		QualifierPER07:        CommunicationNumberQualifierPER03(s.el(7)),
		CommPER08:             s.el(8),
		InquiryReferencePER09: s.el(9),
		Extra:                 s.extra(9),
	}
}

func parseREF(s inputSegment) ReferenceREF {
	return ReferenceREF{Index: s.index, QualifierREF01: ReferenceIdentificationQualifierREF01(s.el(1)), Value: s.el(2), Extra: s.extra(2)}
}

func parsePRV(s inputSegment) ProviderPRV {
	return ProviderPRV{
		Index:          s.index,
		CodePRV01:      ProviderTypePRV01(s.el(1)),
		QualifierPRV02: ReferenceIdentificationQualifierREF01(s.el(2)),
		TaxonomyPRV03:  s.el(3),
		Extra:          s.extra(3),
	}
}

func parseCUR(s inputSegment) CurrencyCUR {
	return CurrencyCUR{Set: true, Index: s.index, EntityIDCUR01: s.el(1), CurrencyCodeCUR02: s.el(2), Extra: s.extra(2)}
}

func parseSBR(s inputSegment) SubscriberSBR {
	return SubscriberSBR{
		Set:                       true,
		Index:                     s.index,
		PayerResponsibilitySBR01:  PayerResponsibilitySBR01(s.el(1)),
		RelationshipSBR02:         RelationshipSBR02(s.el(2)),
		GroupOrPolicyNumberSBR03:  s.el(3),
		GroupNameSBR04:            s.el(4),
		InsuranceTypeSBR05:        InsuranceTypeSBR05(s.el(5)),
		CoordinationSBR06:         s.el(6),
		ResponseSBR07:             s.el(7),
		EmploymentSBR08:           s.el(8),
		ClaimFilingIndicatorSBR09: ClaimFilingIndicatorSBR09(s.el(9)),
		Extra:                     s.extra(9),
	}
}

func parsePAT(s inputSegment) (PatientInfoPAT, error) {
	p := PatientInfoPAT{Set: true, Index: s.index}
	var errs [10]error
	p.RelationshipPAT01 = RelationshipPAT01(s.el(1))
	p.LocationPAT02 = s.el(2)
	p.EmploymentStatusPAT03 = s.el(3)
	p.StudentStatusPAT04 = s.el(4)
	p.FormatPAT05 = s.el(5)
	p.DeathDatePAT06, errs[6] = parseDate(s, 6)
	p.WeightUOMPAT07 = UnitOrBasisForMeasurementCode(s.el(7))
	p.WeightPAT08, errs[8] = parseAmount(s, 8)
	p.PregnancyPAT09 = YesNoConditionResponseCode(s.el(9))
	p.Extra = s.extra(9)
	if err := errors.Join(errs[:]...); err != nil {
		return PatientInfoPAT{}, err
	}
	return p, nil
}

func parseDMG(s inputSegment) (DemographicDMG, error) {
	birth, err := parseDate(s, 2)
	if err != nil {
		return DemographicDMG{}, err
	}
	return DemographicDMG{Set: true, Index: s.index, FormatDMG01: s.el(1), BirthDateDMG02: birth, GenderDMG03: s.el(3), Extra: s.extra(3)}, nil
}

// parseDTP reads a DTP. DTP03 is always kept as received in RawDTP03; when
// DTP02 is D8 it is also read as a CCYYMMDD date.
func parseDTP(s inputSegment) (DateDTP, error) {
	d := DateDTP{Index: s.index}
	d.QualifierDTP01 = DateTimeQualifierDTP01(s.el(1))
	d.FormatDTP02 = s.el(2)
	d.RawDTP03 = s.el(3)
	d.Extra = s.extra(3)
	if d.FormatDTP02 == "D8" {
		date, err := parseDate(s, 3)
		if err != nil {
			return DateDTP{}, err
		}
		d.Date = date
	}
	return d, nil
}

func parseCLM(s inputSegment, component byte) (ClaimCLM, error) {
	c := ClaimCLM{Index: s.index}
	var errs [21]error
	c.PatientControlCLM01 = s.el(1)
	c.TotalChargeCLM02, errs[2] = parseAmount(s, 2)
	c.ClaimFilingIndicatorCLM03 = s.el(3)
	c.ClaimTypeCLM04 = s.el(4)
	location, locationExtra := s.components(5, 3, component)
	c.PlaceOfServiceCLM05_1 = nth(location, 1)
	c.FacilityQualifierCLM05_2 = FacilityCodeQualifierCLM05_2(nth(location, 2))
	c.FrequencyCLM05_3 = ClaimFrequencyCLM05_3(nth(location, 3))
	c.LocationExtraCLM05 = locationExtra
	c.SignatureCLM06 = ProviderOrSupplierSignatureCLM06(s.el(6))
	c.AssignmentCLM07 = AssignmentOrPlanParticipationCLM07(s.el(7))
	c.BenefitsAssignmentCLM08 = BenefitsAssignmentCertificationCLM08(s.el(8))
	c.ReleaseOfInfoCLM09 = ReleaseOfInformationCLM09(s.el(9))
	c.SignatureSourceCLM10 = PatientSignatureSourceCLM10(s.el(10))
	causes, causesExtra := s.components(11, 5, component)
	c.RelatedCausesCLM11[0] = RelatedCausesCodeCLM11_1(nth(causes, 1))
	c.RelatedCausesCLM11[1] = RelatedCausesCodeCLM11_1(nth(causes, 2))
	c.RelatedCausesCLM11_3 = nth(causes, 3)
	c.AutoAccidentStateCLM11_4 = nth(causes, 4)
	c.CountryCLM11_5 = nth(causes, 5)
	c.RelatedCausesExtraCLM11 = causesExtra
	c.SpecialProgramCLM12 = s.el(12)
	c.YesNoCLM13 = s.el(13)
	c.LevelOfServiceCLM14 = s.el(14)
	c.YesNoCLM15 = s.el(15)
	c.ParticipationCLM16 = s.el(16)
	c.ClaimStatusCLM17 = s.el(17)
	c.YesNoCLM18 = s.el(18)
	c.SubmissionReasonCLM19 = s.el(19)
	c.DelayReasonCLM20 = DelayReasonCLM20(s.el(20))
	c.Extra = s.extra(20)
	if err := errors.Join(errs[:]...); err != nil {
		return ClaimCLM{}, err
	}
	return c, nil
}

func parsePWK(s inputSegment) (AttachmentPWK, error) {
	a := AttachmentPWK{Index: s.index}
	var err error
	a.ReportTypePWK01 = AttachmentReportTypeCodePWK01(s.el(1))
	a.TransmissionPWK02 = AttachmentTransmissionCodePWK02(s.el(2))
	a.ReportCopiesPWK03, err = parseAmount(s, 3)
	a.EntityIDPWK04 = s.el(4)
	a.IDQualifierPWK05 = s.el(5)
	a.ControlNumberPWK06 = s.el(6)
	a.Extra = s.extra(6)
	if err != nil {
		return AttachmentPWK{}, err
	}
	return a, nil
}

func parseCN1(s inputSegment) (ContractCN1, error) {
	c := ContractCN1{Set: true, Index: s.index}
	var errs [7]error
	c.TypeCN101 = s.el(1)
	c.AmountCN102, errs[2] = parseAmount(s, 2)
	c.PercentageCN103, errs[3] = parseAmount(s, 3)
	c.CodeCN104 = s.el(4)
	c.TermsDiscountPercentCN105, errs[5] = parseAmount(s, 5)
	c.VersionCN106 = s.el(6)
	c.Extra = s.extra(6)
	if err := errors.Join(errs[:]...); err != nil {
		return ContractCN1{}, err
	}
	return c, nil
}

func parseAMT(s inputSegment) (AmountAMT, error) {
	amount, err := parseAmount(s, 2)
	if err != nil {
		return AmountAMT{}, err
	}
	return AmountAMT{Index: s.index, QualifierAMT01: AmountQualifierAMT01(s.el(1)), Amount: amount, Extra: s.extra(2)}, nil
}

func parseK3(s inputSegment) FileInfoK3 {
	return FileInfoK3{Index: s.index, TextK301: s.el(1), Extra: s.extra(1)}
}

func parseNTE(s inputSegment) NoteNTE {
	return NoteNTE{Set: true, Index: s.index, ReferenceCodeNTE01: s.el(1), TextNTE02: s.el(2), Extra: s.extra(2)}
}

func parseCR1(s inputSegment) (AmbulanceTransportCR1, error) {
	c := AmbulanceTransportCR1{Set: true, Index: s.index}
	var errs [11]error
	c.WeightUnitCR101 = UnitOrBasisForMeasurementCode(s.el(1))
	c.WeightCR102, errs[2] = parseAmount(s, 2)
	c.TransportCodeCR103 = s.el(3)
	c.ReasonCR104 = s.el(4)
	c.DistanceUnitCR105 = UnitOrBasisForMeasurementCode(s.el(5))
	c.DistanceCR106, errs[6] = parseAmount(s, 6)
	c.AddressInfoCR107 = s.el(7)
	c.AddressInfoCR108 = s.el(8)
	c.RoundTripPurposeCR109 = s.el(9)
	c.StretcherPurposeCR110 = s.el(10)
	c.Extra = s.extra(10)
	if err := errors.Join(errs[:]...); err != nil {
		return AmbulanceTransportCR1{}, err
	}
	return c, nil
}

func parseCR2(s inputSegment) (ClaimSpinalManipulationCR2, error) {
	c := ClaimSpinalManipulationCR2{Set: true, Index: s.index}
	var errs [12]error
	c.TreatmentSeriesCR201, errs[1] = parseAmount(s, 1)
	c.TreatmentCountCR202, errs[2] = parseAmount(s, 2)
	c.SubluxationLevelCR203 = s.el(3)
	c.SubluxationLevelCR204 = s.el(4)
	c.UnitCR205 = s.el(5)
	c.TreatmentPeriodCountCR206, errs[6] = parseAmount(s, 6)
	c.MonthlyTreatmentCountCR207, errs[7] = parseAmount(s, 7)
	c.ConditionCodeCR208 = s.el(8)
	c.ComplicationIndicatorCR209 = s.el(9)
	c.ConditionDescriptionCR210 = s.el(10)
	c.ConditionDescriptionCR211 = s.el(11)
	c.Extra = s.extra(11)
	if err := errors.Join(errs[:]...); err != nil {
		return ClaimSpinalManipulationCR2{}, err
	}
	return c, nil
}

func parseCRC(s inputSegment) ConditionCRC {
	c := ConditionCRC{Set: true, Index: s.index}
	c.CodeCategoryCRC01 = CodeCategoryCRC01(s.el(1))
	c.CertConditionCRC02 = YesNoConditionResponseCode(s.el(2))
	for n := 3; n <= 7; n++ {
		c.ConditionsCRC0x[n-3] = ConditionIndicatorCRC0x(s.el(n))
	}
	c.Extra = s.extra(7)
	return c
}

// parseHI reads every composite sent through HI12, an empty one between two
// sent ones included, so each code keeps its position; elements after HI12
// go to Extra. Each composite holds the nine components the X12 standard
// defines: qualifier, code, date format qualifier, date (D8, one date;
// RD8, a range), amount, quantity, version identifier, industry code, and
// a yes/no code. Components 3 through 9 are Not Used in the 837P flat file
// and are kept all the same, and a component after the ninth goes to the
// composite's Extra.
func parseHI(s inputSegment, component byte) (HealthCareInfoHI, error) {
	h := HealthCareInfoHI{Set: true, Index: s.index, Extra: s.extra(12)}
	var errs [13]error
	for n := 1; n <= len(s.elements) && n <= 12; n++ {
		parts, extra := s.components(n, 9, component)
		var code HealthCareCodeHI
		var dateErr, amountErr, quantityErr error
		code.TypeHI0x1 = DiagnosisTypeCodeHI0x_1(nth(parts, 1))
		code.CodeHI0x2 = nth(parts, 2)
		code.DateQualifierHI0x3 = nth(parts, 3)
		code.DateHI0x4, code.DateThroughHI0x4, dateErr = periodText(s, componentRef(s, n, 4), code.DateQualifierHI0x3, nth(parts, 4))
		code.AmountHI0x5, amountErr = amountText(s, componentRef(s, n, 5), nth(parts, 5))
		code.QuantityHI0x6, quantityErr = amountText(s, componentRef(s, n, 6), nth(parts, 6))
		code.VersionIdentifierHI0x7 = nth(parts, 7)
		code.IndustryCodeHI0x8 = nth(parts, 8)
		code.YesNoHI0x9 = YesNoConditionResponseCode(nth(parts, 9))
		code.Extra = extra
		errs[n] = errors.Join(dateErr, amountErr, quantityErr)
		h.Codes = append(h.Codes, code)
	}
	if err := errors.Join(errs[:]...); err != nil {
		return HealthCareInfoHI{}, err
	}
	return h, nil
}

func parseHCP(s inputSegment) (PricingHCP, error) {
	h := PricingHCP{Set: true, Index: s.index}
	var errs [16]error
	h.MethodologyHCP01 = PricingMethodologyHCP01(s.el(1))
	h.AllowedAmountHCP02, errs[2] = parseAmount(s, 2)
	h.SavingAmountHCP03, errs[3] = parseAmount(s, 3)
	h.OrgIDHCP04 = s.el(4)
	h.PerDiemOrFlatRateHCP05, errs[5] = parseAmount(s, 5)
	h.APGCodeHCP06 = s.el(6)
	h.APGAmountHCP07, errs[7] = parseAmount(s, 7)
	h.ProductServiceIDHCP08 = s.el(8)
	h.IDQualifierHCP09 = s.el(9)
	h.ApprovedHCPCSHCP10 = s.el(10)
	h.UnitHCP11 = UnitOrBasisForMeasurementCode(s.el(11))
	h.ApprovedUnitsHCP12, errs[12] = parseAmount(s, 12)
	h.RejectReasonHCP13 = s.el(13)
	h.PolicyComplianceHCP14 = s.el(14)
	h.ExceptionHCP15 = s.el(15)
	h.Extra = s.extra(15)
	if err := errors.Join(errs[:]...); err != nil {
		return PricingHCP{}, err
	}
	return h, nil
}

// parseCAS reads the group code and every reason/amount/quantity triple
// sent through CAS19, an empty triple between two sent ones included, so
// each triple keeps its position; elements after CAS19 go to Extra.
func parseCAS(s inputSegment) (AdjustmentCAS, error) {
	a := AdjustmentCAS{Index: s.index, Extra: s.extra(19)}
	var errs [20]error
	a.GroupCAS01 = ClaimAdjustmentGroupCodeCAS01(s.el(1))
	for n := 2; n <= len(s.elements) && n <= 17; n += 3 {
		var r CASReason
		r.Reason = carc.Code(s.el(n))
		r.Amount, errs[n+1] = parseAmount(s, n+1)
		r.Quantity, errs[n+2] = parseAmount(s, n+2)
		a.Reasons = append(a.Reasons, r)
	}
	if err := errors.Join(errs[:]...); err != nil {
		return AdjustmentCAS{}, err
	}
	return a, nil
}

func parseOI(s inputSegment) OtherInsuranceOI {
	return OtherInsuranceOI{
		Set:                      true,
		Index:                    s.index,
		ClaimFilingIndicatorOI01: s.el(1),
		ClaimSubmitReasonOI02:    s.el(2),
		BenefitsAssignmentOI03:   BenefitsAssignmentCertificationCLM08(s.el(3)),
		SignatureSourceOI04:      PatientSignatureSourceCLM10(s.el(4)),
		ProviderAgreementOI05:    s.el(5),
		ReleaseOfInfoOI06:        ReleaseOfInformationCLM09(s.el(6)),
		Extra:                    s.extra(6),
	}
}

func parseMOA(s inputSegment) (OutpatientAdjudicationMOA, error) {
	m := OutpatientAdjudicationMOA{Set: true, Index: s.index}
	var errs [10]error
	m.ReimbursementRateMOA01, errs[1] = parseAmount(s, 1)
	m.HCPCSPayableAmountMOA02, errs[2] = parseAmount(s, 2)
	m.RemarkMOA03 = rarc.Code(s.el(3))
	m.RemarkMOA04 = rarc.Code(s.el(4))
	m.RemarkMOA05 = rarc.Code(s.el(5))
	m.RemarkMOA06 = rarc.Code(s.el(6))
	m.RemarkMOA07 = rarc.Code(s.el(7))
	m.ESRDPaymentMOA08, errs[8] = parseAmount(s, 8)
	m.NonPayableProfCompMOA09, errs[9] = parseAmount(s, 9)
	m.Extra = s.extra(9)
	if err := errors.Join(errs[:]...); err != nil {
		return OutpatientAdjudicationMOA{}, err
	}
	return m, nil
}

func parseLX(s inputSegment) ServiceLineLX {
	return ServiceLineLX{Index: s.index, AssignedNumberLX01: s.el(1), Extra: s.extra(1)}
}

func parseSV1(s inputSegment, component byte) (ProfessionalServiceSV1, error) {
	v := ProfessionalServiceSV1{Set: true, Index: s.index}
	var errs [16]error
	v.ProcedureSV101, errs[1] = parseProcedure(s, 1, component)
	v.ChargeSV102, errs[2] = parseAmount(s, 2)
	v.UnitSV103 = UnitOrBasisForMeasurementCode(s.el(3))
	v.UnitCountSV104, errs[4] = parseAmount(s, 4)
	v.PlaceOfServiceSV105 = s.el(5)
	v.ServiceTypeSV106 = s.el(6)
	v.DiagnosisPointersSV107, v.PointersExtraSV107, errs[7] = parsePointers(s, 7, component)
	v.MonetaryAmountSV108, errs[8] = parseAmount(s, 8)
	v.EmergencySV109 = YesNoConditionResponseCode(s.el(9))
	v.MultipleProcedureSV110 = s.el(10)
	v.EPSDTSV111 = YesNoConditionResponseCode(s.el(11))
	v.FamilyPlanningSV112 = YesNoConditionResponseCode(s.el(12))
	v.ReviewCodeSV113 = s.el(13)
	v.ReviewValueSV114 = s.el(14)
	v.CoPayStatusSV115 = s.el(15)
	v.Extra = s.extra(15)
	if err := errors.Join(errs[:]...); err != nil {
		return ProfessionalServiceSV1{}, err
	}
	return v, nil
}

func parseSV5(s inputSegment, component byte) (LineDMEServiceSV5, error) {
	d := LineDMEServiceSV5{Set: true, Index: s.index}
	var errs [7]error
	d.ProcedureSV501, errs[1] = parseProcedure(s, 1, component)
	d.UnitSV502 = UnitOrBasisForMeasurementCode(s.el(2))
	d.LengthOfNecessitySV503, errs[3] = parseAmount(s, 3)
	d.RentalPriceSV504, errs[4] = parseAmount(s, 4)
	d.PurchasePriceSV505, errs[5] = parseAmount(s, 5)
	d.RentalUnitPriceIndicatorSV506 = s.el(6)
	d.Extra = s.extra(6)
	if err := errors.Join(errs[:]...); err != nil {
		return LineDMEServiceSV5{}, err
	}
	return d, nil
}

func parseCR3(s inputSegment) (LineDMECertificationCR3, error) {
	duration, err := parseAmount(s, 3)
	if err != nil {
		return LineDMECertificationCR3{}, err
	}
	return LineDMECertificationCR3{
		Set:                    true,
		Index:                  s.index,
		CertificationTypeCR301: s.el(1),
		UnitCR302:              UnitOrBasisForMeasurementCode(s.el(2)),
		DurationCR303:          duration,
		Extra:                  s.extra(3),
	}, nil
}

func parseQTY(s inputSegment) (QuantityQTY, error) {
	quantity, err := parseAmount(s, 2)
	if err != nil {
		return QuantityQTY{}, err
	}
	return QuantityQTY{Index: s.index, QualifierQTY01: QuantityQualifierQTY01(s.el(1)), Quantity: quantity, Extra: s.extra(2)}, nil
}

func parseMEA(s inputSegment) (TestResultMEA, error) {
	result, err := parseAmount(s, 3)
	if err != nil {
		return TestResultMEA{}, err
	}
	return TestResultMEA{Index: s.index, ReferenceIDMEA01: s.el(1), QualifierMEA02: s.el(2), ResultMEA03: result, Extra: s.extra(3)}, nil
}

func parsePS1(s inputSegment) (PurchasedServicePS1, error) {
	charge, err := parseAmount(s, 2)
	if err != nil {
		return PurchasedServicePS1{}, err
	}
	return PurchasedServicePS1{Set: true, Index: s.index, ProviderIDPS101: s.el(1), ChargeAmountPS102: charge, Extra: s.extra(2)}, nil
}

func parseLIN(s inputSegment) DrugLIN {
	return DrugLIN{Set: true, Index: s.index, AssignedIDLIN01: s.el(1), QualifierLIN02: s.el(2), NDCLIN03: s.el(3), Extra: s.extra(3)}
}

// parseCTP reads a CTP: CTP01-03, the unit count, and the first component
// of CTP05, with the rest of that composite kept in UnitCTP05Extra.
func parseCTP(s inputSegment, component byte) (DrugPricingCTP, error) {
	c := DrugPricingCTP{Set: true, Index: s.index}
	var errs [6]error
	c.ClassOfTradeCTP01 = s.el(1)
	c.PriceIdentifierCTP02 = s.el(2)
	c.UnitPriceCTP03, errs[3] = parseAmount(s, 3)
	c.UnitCountCTP04, errs[4] = parseAmount(s, 4)
	unit, unitExtra := s.components(5, 1, component)
	c.UnitCTP05_1 = UnitOrBasisForMeasurementCode(nth(unit, 1))
	c.UnitCTP05Extra = unitExtra
	c.Extra = s.extra(5)
	if err := errors.Join(errs[:]...); err != nil {
		return DrugPricingCTP{}, err
	}
	return c, nil
}

func parseSVD(s inputSegment, component byte) (LineAdjudicationSVD, error) {
	a := LineAdjudicationSVD{Index: s.index}
	var errs [7]error
	a.OtherPayerSVD01 = s.el(1)
	a.PaidAmountSVD02, errs[2] = parseAmount(s, 2)
	a.ProcedureSVD03, errs[3] = parseProcedure(s, 3, component)
	a.ProductServiceIDSVD04 = s.el(4)
	a.PaidUnitsSVD05, errs[5] = parseAmount(s, 5)
	a.BundledLineSVD06 = s.el(6)
	a.Extra = s.extra(6)
	if err := errors.Join(errs[:]...); err != nil {
		return LineAdjudicationSVD{}, err
	}
	return a, nil
}

func parseLQ(s inputSegment) FormLQ {
	return FormLQ{Index: s.index, QualifierLQ01: s.el(1), FormIDLQ02: s.el(2), Extra: s.extra(2)}
}

func parseFRM(s inputSegment) (FormAnswerFRM, error) {
	f := FormAnswerFRM{Index: s.index}
	var errs [6]error
	f.QuestionFRM01 = s.el(1)
	f.ResponseFRM02 = s.el(2)
	f.ResponseFRM03 = s.el(3)
	f.ResponseFRM04, errs[4] = parseDate(s, 4)
	f.ResponseFRM05, errs[5] = parseAmount(s, 5)
	f.Extra = s.extra(5)
	if err := errors.Join(errs[:]...); err != nil {
		return FormAnswerFRM{}, err
	}
	return f, nil
}

// parseProcedure reads the composite procedure identifier in element n
// (SV101, SV501, or SVD03): qualifier, code, four modifiers, and
// description. SV101-8 is Not Used, so the model names seven components and
// keeps any after them in Extra.
func parseProcedure(s inputSegment, n int, component byte) (ProcedureSV101, error) {
	parts, extra := s.components(n, 7, component)
	return ProcedureSV101{
		Qualifier:   nth(parts, 1),
		Code:        nth(parts, 2),
		Modifiers:   [4]string{nth(parts, 3), nth(parts, 4), nth(parts, 5), nth(parts, 6)},
		Description: nth(parts, 7),
		Extra:       extra,
	}, nil
}

// parsePointers reads the SV107 composite: up to four diagnosis code
// pointers, each a whole number from 1, with any component after the
// fourth returned as extra text. A pointer not sent stays 0, so a sent 0 is
// rejected rather than read as not sent.
func parsePointers(s inputSegment, n int, component byte) ([4]int, []string, error) {
	var pointers [4]int
	parts, extra := s.components(n, 4, component)
	for i, text := range parts {
		if text == "" {
			continue
		}
		value, err := strconv.Atoi(text)
		if err != nil || !allDigits(text) || value < 1 {
			return [4]int{}, nil, fmt.Errorf("x12837p: parsing %s-%d in segment %d: diagnosis pointer is not a whole number from 1", s.name(n), i+1, s.index)
		}
		pointers[i] = value
	}
	return pointers, extra, nil
}

// componentRef returns component k of element n's reference, e.g. "HI01-4".
func componentRef(s inputSegment, n, k int) string {
	return fmt.Sprintf("%s-%d", s.name(n), k)
}

// trimTrailingEmpty removes empty strings from the end of parts. A trailing
// empty element or component carries no text.
func trimTrailingEmpty(parts []string) []string {
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// nth returns component i (1-based) of a split composite, or "" when the
// composite ends before it.
func nth(parts []string, i int) string {
	if i < 1 || i > len(parts) {
		return ""
	}
	return parts[i-1]
}

// parseAmount reads element n as an exact decimal. An empty element is an
// Amount with Set false.
func parseAmount(s inputSegment, n int) (Amount, error) {
	return amountText(s, s.name(n), s.el(n))
}

// amountText reads text as an exact decimal; ref names the element or
// component for errors. Empty text is an Amount with Set false; otherwise
// the Amount keeps the text as sent beside the parsed value.
func amountText(s inputSegment, ref, text string) (Amount, error) {
	if text == "" {
		return Amount{}, nil
	}
	value, err := ParseDecimal([]byte(text))
	if err != nil {
		return Amount{}, fmt.Errorf("x12837p: parsing %s in segment %d: %w", ref, s.index, err)
	}
	return Amount{Value: value, Text: text, Set: true}, nil
}

// parseDate reads element n as CCYYMMDD. An empty element is a Date with
// Set false.
func parseDate(s inputSegment, n int) (Date, error) {
	if s.el(n) == "" {
		return Date{}, nil
	}
	return dateText(s, s.name(n), s.el(n))
}

// dateText reads text as CCYYMMDD; ref names the element or component for
// errors. Unlike parseDate it rejects empty text, so each half of a range
// must be a date.
func dateText(s inputSegment, ref, text string) (Date, error) {
	if len(text) != 8 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: date is not CCYYMMDD", ref, s.index)
	}
	d := Date{Year: number(text[0:4]), Month: number(text[4:6]), Day: number(text[6:8]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: not a calendar date", ref, s.index)
	}
	return d, nil
}

// periodText reads text as the date format qualifier format names: D8, one
// CCYYMMDD date, returned first; or RD8, a CCYYMMDD-CCYYMMDD range, returned
// as its first and second dates. ref names the element or component, e.g.
// "HI01-4", for errors. Empty text is two Dates with Set false.
func periodText(s inputSegment, ref, format, text string) (Date, Date, error) {
	if text == "" {
		return Date{}, Date{}, nil
	}
	switch format {
	case "D8":
		date, err := dateText(s, ref, text)
		return date, Date{}, err
	case "RD8":
		from, through, ok := strings.Cut(text, "-")
		if !ok {
			return Date{}, Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: date range is not CCYYMMDD-CCYYMMDD", ref, s.index)
		}
		first, err := dateText(s, ref, from)
		if err != nil {
			return Date{}, Date{}, err
		}
		second, err := dateText(s, ref, through)
		if err != nil {
			return Date{}, Date{}, err
		}
		return first, second, nil
	}
	return Date{}, Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: its date format qualifier is not D8 or RD8", ref, s.index)
}

// parseInterchangeDate reads ISA09, which is YYMMDD, with Year taken as 2000
// plus YY, the same convention as backend/pkg/x12835.
func parseInterchangeDate(s inputSegment, n int) (Date, error) {
	text := s.el(n)
	if text == "" {
		return Date{}, nil
	}
	if len(text) != 6 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: date is not YYMMDD", s.name(n), s.index)
	}
	d := Date{Year: 2000 + number(text[0:2]), Month: number(text[2:4]), Day: number(text[4:6]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12837p: parsing %s in segment %d: not a calendar date", s.name(n), s.index)
	}
	return d, nil
}

// calendarDate reports whether d names a real day. time.Date moves an
// out-of-range month or day into a neighboring one, so a changed field
// means the date does not exist.
func calendarDate(d Date) bool {
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	return t.Year() == d.Year && int(t.Month()) == d.Month && t.Day() == d.Day
}

// parseTime reads element n as HHMM or HHMMSS. An empty element is a Time
// with Set false. Decimal seconds (HHMMSSD, HHMMSSDD) are rejected, since
// Time has no field for them.
func parseTime(s inputSegment, n int) (Time, error) {
	text := s.el(n)
	if text == "" {
		return Time{}, nil
	}
	if (len(text) != 4 && len(text) != 6) || !allDigits(text) {
		return Time{}, fmt.Errorf("x12837p: parsing %s in segment %d: time is not HHMM or HHMMSS", s.name(n), s.index)
	}
	t := Time{Hour: number(text[0:2]), Minute: number(text[2:4]), Set: true}
	if len(text) == 6 {
		t.Second = number(text[4:6])
		t.HasSeconds = true
	}
	if t.Hour > 23 || t.Minute > 59 || t.Second > 59 {
		return Time{}, fmt.Errorf("x12837p: parsing %s in segment %d: not a clock time", s.name(n), s.index)
	}
	return t, nil
}

// parseWholeNumber reads element n as a count of digits only.
func parseWholeNumber(s inputSegment, n int) (int, error) {
	text := s.el(n)
	if text == "" || !allDigits(text) {
		return 0, fmt.Errorf("x12837p: parsing %s in segment %d: not a whole number", s.name(n), s.index)
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("x12837p: parsing %s in segment %d: number too large", s.name(n), s.index)
	}
	return value, nil
}

func allDigits(text string) bool {
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return false
		}
	}
	return true
}

// number returns the value of a string of ASCII digits.
func number(digits string) int {
	n := 0
	for i := 0; i < len(digits); i++ {
		n = n*10 + int(digits[i]-'0')
	}
	return n
}
