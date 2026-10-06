package x12837i

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
// struct, and records the segment's position in the file in the struct's
// Index. An array indexed by element number collects what each element's
// parse reports: a malformed amount, date, time, or count. Any one fails
// the segment, and errors.Join reports them all.
//
// A coded element is converted to its string type directly, whatever its
// text: a code off the guide's list is kept as sent, and Known reports it.
// Nothing about a code can fail a segment.
//
// checkLast's limit is the number of elements X12 defines for the segment.
// Where the 837I guide uses fewer, the positions past the model's named
// fields go to the struct's ExtraElements as sent, so nothing a file sends
// in an X12-defined position is dropped; text past X12's own last position
// is an error.
//
// So each element's line reads: field, its slot = parse(element n).

func parseISA(s inputSegment) (InterchangeISA, error) {
	if err := s.checkLast(16); err != nil {
		return InterchangeISA{}, err
	}
	isa := InterchangeISA{Index: s.index}
	var errs [17]error
	isa.AuthQualifierISA01 = s.el(1)
	isa.AuthInfoISA02 = s.el(2)
	isa.SecQualifierISA03 = s.el(3)
	isa.SecInfoISA04 = s.el(4)
	isa.SenderIDQualifierISA05 = s.el(5)
	isa.SenderIDISA06 = s.el(6)
	isa.ReceiverIDQualifierISA07 = s.el(7)
	isa.ReceiverIDISA08 = s.el(8)
	isa.DateISA09, errs[9] = parseInterchangeDate(s, 9)
	isa.TimeISA10, errs[10] = parseTime(s, 10)
	isa.RepetitionSepISA11 = s.el(11)
	isa.VersionISA12 = s.el(12)
	isa.ControlNumberISA13 = s.el(13)
	isa.AckRequestedISA14 = s.el(14)
	isa.UsageISA15 = s.el(15)
	isa.ComponentSepISA16 = s.el(16)
	if err := errors.Join(errs[:]...); err != nil {
		return InterchangeISA{}, err
	}
	return isa, nil
}

func parseGS(s inputSegment) (FunctionalGroupGS, error) {
	if err := s.checkLast(8); err != nil {
		return FunctionalGroupGS{}, err
	}
	g := FunctionalGroupGS{Index: s.index}
	var errs [9]error
	g.FunctionalIDGS01 = TransactionTypeGS01_ST01(s.el(1))
	g.SenderCodeGS02 = s.el(2)
	g.ReceiverCodeGS03 = s.el(3)
	g.DateGS04, errs[4] = parseDate(s, 4)
	g.TimeGS05, errs[5] = parseTime(s, 5)
	g.ControlNumberGS06 = s.el(6)
	g.AgencyGS07 = s.el(7)
	g.VersionGS08 = s.el(8)
	if err := errors.Join(errs[:]...); err != nil {
		return FunctionalGroupGS{}, err
	}
	return g, nil
}

func parseST(s inputSegment) (TransactionST, error) {
	if err := s.checkLast(3); err != nil {
		return TransactionST{}, err
	}
	t := TransactionST{Index: s.index}
	t.TransactionSetIDST01 = TransactionTypeGS01_ST01(s.el(1))
	t.ControlNumberST02 = s.el(2)
	t.VersionST03 = s.el(3)
	return t, nil
}

// parseTrailer reads an SE, GE, or IEA: a count of digits in element 01,
// returned as the digits the file sent, and a control number in element
// 02.
func parseTrailer(s inputSegment) (string, string, error) {
	if err := s.checkLast(2); err != nil {
		return "", "", err
	}
	if _, err := parseWholeNumber(s, 1); err != nil {
		return "", "", err
	}
	return s.el(1), s.el(2), nil
}

func parseBHT(s inputSegment) (HeaderBHT, error) {
	if err := s.checkLast(6); err != nil {
		return HeaderBHT{}, err
	}
	h := HeaderBHT{Set: true, Index: s.index}
	var errs [7]error
	h.StructureCodeBHT01 = s.el(1)
	h.PurposeCodeBHT02 = s.el(2)
	h.OriginatorTransactionIDBHT03 = s.el(3)
	h.CreationDateBHT04, errs[4] = parseDate(s, 4)
	h.CreationTimeBHT05, errs[5] = parseTime(s, 5)
	h.ClaimOrEncounterBHT06 = ClaimOrEncounterBHT06(s.el(6))
	if err := errors.Join(errs[:]...); err != nil {
		return HeaderBHT{}, err
	}
	return h, nil
}

func parseNM1(s inputSegment) (NameNM1, error) {
	if err := s.checkLast(12); err != nil {
		return NameNM1{}, err
	}
	n := NameNM1{Index: s.index}
	n.EntityIDNM101 = PartyRoleNM101(s.el(1))
	n.EntityTypeNM102 = EntityTypeNM102(s.el(2))
	n.LastNM103 = s.el(3)
	n.FirstNM104 = s.el(4)
	n.MiddleNM105 = s.el(5)
	n.PrefixNM106 = s.el(6)
	n.SuffixNM107 = s.el(7)
	n.IDQualifierNM108 = IDQualifierNM108(s.el(8))
	n.IDNM109 = s.el(9)
	n.ExtraElements = extraElements(s, 9)
	return n, nil
}

func parsePER(s inputSegment) (ContactPER, error) {
	if err := s.checkLast(9); err != nil {
		return ContactPER{}, err
	}
	return ContactPER{
		Index:          s.index,
		FunctionPER01:  s.el(1),
		NamePER02:      s.el(2),
		QualifierPER03: s.el(3),
		CommPER04:      s.el(4),
		QualifierPER05: s.el(5),
		CommPER06:      s.el(6),
		QualifierPER07: s.el(7),
		CommPER08:      s.el(8),
		ExtraElements:  extraElements(s, 8),
	}, nil
}

func parseN3(s inputSegment) (AddressN3, error) {
	if err := s.checkLast(2); err != nil {
		return AddressN3{}, err
	}
	return AddressN3{Set: true, Index: s.index, Line1N301: s.el(1), Line2N302: s.el(2)}, nil
}

func parseN4(s inputSegment) (LocationN4, error) {
	if err := s.checkLast(7); err != nil {
		return LocationN4{}, err
	}
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
	}, nil
}

func parseREF(s inputSegment) (ReferenceREF, error) {
	if err := s.checkLast(4); err != nil {
		return ReferenceREF{}, err
	}
	return ReferenceREF{Index: s.index, QualifierREF01: s.el(1), Value: s.el(2), ExtraElements: extraElements(s, 2)}, nil
}

func parseHL(s inputSegment) (HierarchyHL, error) {
	if err := s.checkLast(4); err != nil {
		return HierarchyHL{}, err
	}
	h := HierarchyHL{Index: s.index}
	h.IDNumberHL01 = s.el(1)
	h.ParentIDNumberHL02 = s.el(2)
	h.LevelCodeHL03 = HierarchicalLevelHL03(s.el(3))
	h.ChildCodeHL04 = s.el(4)
	return h, nil
}

func parsePRV(s inputSegment) (ProviderPRV, error) {
	if err := s.checkLast(6); err != nil {
		return ProviderPRV{}, err
	}
	p := ProviderPRV{Set: true, Index: s.index}
	p.ProviderCodePRV01 = ProviderRoleCodePRV01(s.el(1))
	p.RefQualifierPRV02 = s.el(2)
	p.TaxonomyCodePRV03 = s.el(3)
	p.ExtraElements = extraElements(s, 3)
	return p, nil
}

func parseCUR(s inputSegment) (CurrencyCUR, error) {
	if err := s.checkLast(21); err != nil {
		return CurrencyCUR{}, err
	}
	return CurrencyCUR{Set: true, Index: s.index, EntityIDCUR01: s.el(1), CurrencyCodeCUR02: s.el(2), ExtraElements: extraElements(s, 2)}, nil
}

func parseSBR(s inputSegment) (SubscriberSBR, error) {
	if err := s.checkLast(9); err != nil {
		return SubscriberSBR{}, err
	}
	b := SubscriberSBR{Set: true, Index: s.index}
	b.PayerResponsibilitySBR01 = PayerResponsibilitySBR01(s.el(1))
	b.RelationshipSBR02 = RelationshipSBR02(s.el(2))
	b.GroupNumberSBR03 = s.el(3)
	b.GroupNameSBR04 = s.el(4)
	b.InsuranceTypeSBR05 = s.el(5)
	b.CoordinationSBR06 = s.el(6)
	b.ResponseSBR07 = s.el(7)
	b.EmploymentSBR08 = s.el(8)
	b.ClaimFilingIndicatorSBR09 = ClaimFilingIndicatorSBR09(s.el(9))
	return b, nil
}

// parseDMG reads a DMG whose DMG01 the caller has checked is D8 or not
// sent, so DMG02 is a CCYYMMDD date.
func parseDMG(s inputSegment) (DemographicDMG, error) {
	if err := s.checkLast(11); err != nil {
		return DemographicDMG{}, err
	}
	d := DemographicDMG{Set: true, Index: s.index, DateFormatDMG01: s.el(1)}
	birth, err := parseDate(s, 2)
	if err != nil {
		return DemographicDMG{}, err
	}
	d.BirthDateDMG02 = birth
	d.SexDMG03 = SexDMG03(s.el(3))
	d.ExtraElements = extraElements(s, 3)
	return d, nil
}

func parsePAT(s inputSegment) (PatientRelationshipPAT, error) {
	if err := s.checkLast(9); err != nil {
		return PatientRelationshipPAT{}, err
	}
	p := PatientRelationshipPAT{Set: true, Index: s.index}
	p.RelationshipPAT01 = RelationshipSBR02(s.el(1))
	p.ExtraElements = extraElements(s, 1)
	return p, nil
}

// parseCLM reads a CLM, splitting the CLM05 composite (facility type,
// facility code qualifier, claim frequency). CLM11, the Related Causes
// composite, is Not Used and kept as its received text.
func parseCLM(s inputSegment, component byte) (ClaimCLM, error) {
	if err := s.checkLast(20); err != nil {
		return ClaimCLM{}, err
	}
	c := ClaimCLM{Index: s.index}
	var errs [21]error
	c.PatientControlCLM01 = s.el(1)
	c.TotalChargeCLM02, errs[2] = parseAmount(s, 2)
	c.FilingIndicatorCLM03 = s.el(3)
	c.ClaimTypeCLM04 = s.el(4)
	var location []string
	location, errs[5] = components(s, 5, component, 3)
	c.FacilityTypeCLM05_1 = nth(location, 1)
	c.FacilityCodeQualifierCLM05_2 = FacilityCodeQualifierCLM05_2(nth(location, 2))
	c.ClaimFrequencyCLM05_3 = nth(location, 3)
	c.SignatureIndicatorCLM06 = s.el(6)
	c.AssignmentCLM07 = AssignmentCLM07(s.el(7))
	c.BenefitsAssignCLM08 = BenefitsAssignCLM08(s.el(8))
	c.ReleaseOfInfoCLM09 = ReleaseOfInfoCLM09(s.el(9))
	c.SignatureSourceCLM10 = s.el(10)
	c.RelatedCausesCLM11 = s.el(11)
	c.SpecialProgramCLM12 = s.el(12)
	c.YesNoCLM13 = s.el(13)
	c.LevelOfServiceCLM14 = s.el(14)
	c.YesNoCLM15 = s.el(15)
	c.ParticipationAgreementCLM16 = s.el(16)
	c.StatusCLM17 = s.el(17)
	c.YesNoCLM18 = s.el(18)
	c.SubmissionReasonCLM19 = s.el(19)
	c.DelayReasonCLM20 = DelayReasonCLM20(s.el(20))
	if err := errors.Join(errs[:]...); err != nil {
		return ClaimCLM{}, err
	}
	return c, nil
}

// parseDTP reads a DTP whose DTP02 the caller has checked is D8 or RD8,
// into a DateDTP with its date or its range.
func parseDTP(s inputSegment) (DateDTP, error) {
	d := DateDTP{Set: true, Index: s.index}
	d.QualifierDTP01 = DateQualifierDTP01(s.el(1))
	d.FormatDTP02 = s.el(2)
	from, through, err := parseDTPPeriod(s)
	if err != nil {
		return DateDTP{}, err
	}
	d.Date, d.DateThrough = from, through
	return d, nil
}

// parseDTPPeriod reads DTP03 in the format DTP02 names: D8, one CCYYMMDD
// date, returned first; or RD8, a CCYYMMDD-CCYYMMDD range, returned as its
// first and second dates.
func parseDTPPeriod(s inputSegment) (Date, Date, error) {
	if err := s.checkLast(3); err != nil {
		return Date{}, Date{}, err
	}
	return periodText(s, s.name(3), s.el(2), s.el(3))
}

// periodText reads text as the date format qualifier format names: D8 or
// RD8. ref names the element or component, e.g. "DTP03" or "HI01-4", for
// errors. Empty text is two Dates with Set false.
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
			return Date{}, Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: date range is not CCYYMMDD-CCYYMMDD", ref, s.index)
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
	return Date{}, Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: its date format qualifier is not D8 or RD8", ref, s.index)
}

// parseDateHour reads element n in the format format names: DT,
// CCYYMMDDHHMM, a date and an hour; or D8, a date alone.
func parseDateHour(s inputSegment, n int, format string) (Date, Time, error) {
	text := s.el(n)
	if format == "D8" {
		date, err := parseDate(s, n)
		return date, Time{}, err
	}
	if len(text) != 12 || !allDigits(text) {
		return Date{}, Time{}, fmt.Errorf("x12837i: parsing %s in segment %d: date and time is not CCYYMMDDHHMM", s.name(n), s.index)
	}
	date, err := dateText(s, s.name(n), text[:8])
	if err != nil {
		return Date{}, Time{}, err
	}
	hour, err := timeText(s, s.name(n), text[8:])
	if err != nil {
		return Date{}, Time{}, err
	}
	return date, hour, nil
}

func parseCL1(s inputSegment) (InstitutionalCL1, error) {
	if err := s.checkLast(4); err != nil {
		return InstitutionalCL1{}, err
	}
	c := InstitutionalCL1{Set: true, Index: s.index}
	c.AdmissionTypeCL101 = s.el(1)
	c.AdmissionSourceCL102 = s.el(2)
	c.PatientStatusCL103 = s.el(3)
	c.NursingHomeCL104 = s.el(4)
	return c, nil
}

func parsePWK(s inputSegment) (SupplementalPWK, error) {
	if err := s.checkLast(9); err != nil {
		return SupplementalPWK{}, err
	}
	p := SupplementalPWK{Index: s.index}
	p.ReportTypePWK01 = PWKReportTypeCode01(s.el(1))
	p.TransmissionPWK02 = PWKTransmissionCode02(s.el(2))
	p.CopiesNeededPWK03 = s.el(3)
	p.EntityIDPWK04 = s.el(4)
	p.IDQualifierPWK05 = s.el(5)
	p.AttachmentControlPWK06 = s.el(6)
	p.ExtraElements = extraElements(s, 6)
	return p, nil
}

func parseCN1(s inputSegment) (ContractCN1, error) {
	if err := s.checkLast(6); err != nil {
		return ContractCN1{}, err
	}
	c := ContractCN1{Index: s.index}
	var errs [7]error
	c.ContractTypeCN101 = s.el(1)
	c.ContractAmountCN102, errs[2] = parseAmount(s, 2)
	c.ContractPercentCN103, errs[3] = parseAmount(s, 3)
	c.ContractCodeCN104 = s.el(4)
	c.TermsDiscountCN105, errs[5] = parseAmount(s, 5)
	c.ContractVersionCN106 = s.el(6)
	if err := errors.Join(errs[:]...); err != nil {
		return ContractCN1{}, err
	}
	return c, nil
}

func parseAMT(s inputSegment) (AmountAMT, error) {
	if err := s.checkLast(3); err != nil {
		return AmountAMT{}, err
	}
	amount, err := parseAmount(s, 2)
	if err != nil {
		return AmountAMT{}, err
	}
	return AmountAMT{Index: s.index, QualifierAMT01: s.el(1), Amount: amount, ExtraElements: extraElements(s, 2)}, nil
}

func parseK3(s inputSegment) (FileInfoK3, error) {
	if err := s.checkLast(3); err != nil {
		return FileInfoK3{}, err
	}
	return FileInfoK3{Index: s.index, FixedFormatK301: s.el(1), ExtraElements: extraElements(s, 1)}, nil
}

func parseNTE(s inputSegment) (NoteNTE, error) {
	if err := s.checkLast(2); err != nil {
		return NoteNTE{}, err
	}
	n := NoteNTE{Index: s.index}
	n.NoteCodeNTE01 = NTENoteCode01(s.el(1))
	n.TextNTE02 = s.el(2)
	return n, nil
}

func parseCRC(s inputSegment) (ReferralCRC, error) {
	if err := s.checkLast(7); err != nil {
		return ReferralCRC{}, err
	}
	r := ReferralCRC{Index: s.index}
	r.CategoryCRC01 = CRCCategoryCRC01(s.el(1))
	r.CertificationCRC02 = CRCCertificationCRC02(s.el(2))
	r.Condition1CRC03 = CRCConditionIndicator(s.el(3))
	r.Condition2CRC04 = CRCConditionIndicator(s.el(4))
	r.Condition3CRC05 = CRCConditionIndicator(s.el(5))
	r.ExtraElements = extraElements(s, 5)
	return r, nil
}

// hiQualifiers returns the first component (HI0n-1) of each composite an HI
// segment sends, skipping an empty composite.
func hiQualifiers(s inputSegment, component byte) []DiagnosisQualifierHI01_1 {
	var qualifiers []DiagnosisQualifierHI01_1
	for n := 1; n <= len(s.elements); n++ {
		if s.el(n) == "" {
			continue
		}
		first, _, _ := strings.Cut(s.el(n), string(component))
		qualifiers = append(qualifiers, DiagnosisQualifierHI01_1(first))
	}
	return qualifiers
}

// parseHI reads every composite an HI segment sends, skipping an empty one:
// qualifier, code, date format qualifier and date (D8, one date; RD8, a
// range), amount, quantity, version identifier, industry code, and
// present-on-admission indicator. Each composite records the segment's
// position and its own element number. Components 7 and 8 (Version
// Identifier and Industry Code) are Not Used and kept as plain text.
func parseHI(s inputSegment, component byte) ([]HealthInfoHI, error) {
	if err := s.checkLast(12); err != nil {
		return nil, err
	}
	var infos []HealthInfoHI
	var errs [13]error
	for n := 1; n <= len(s.elements); n++ {
		if s.el(n) == "" {
			continue
		}
		parts, err := components(s, n, component, 9)
		if err != nil {
			errs[n] = err
			continue
		}
		hi := HealthInfoHI{Set: true, Index: s.index, Position: n}
		hi.QualifierHI01_1 = DiagnosisQualifierHI01_1(nth(parts, 1))
		hi.CodeHI01_2 = nth(parts, 2)
		hi.DateQualifierHI01_3 = nth(parts, 3)
		var dateErr, amountErr, quantityErr error
		hi.DateHI01_4, hi.DateThroughHI01_4, dateErr = periodText(s, componentRef(s, n, 4), hi.DateQualifierHI01_3, nth(parts, 4))
		hi.AmountHI01_5, amountErr = amountText(s, componentRef(s, n, 5), nth(parts, 5))
		hi.QuantityHI01_6, quantityErr = amountText(s, componentRef(s, n, 6), nth(parts, 6))
		hi.VersionIdentifierHI0n_7 = nth(parts, 7)
		hi.IndustryCodeHI0n_8 = nth(parts, 8)
		hi.PresentOnAdmissionHI0n_9 = PresentOnAdmissionHI0n_9(nth(parts, 9))
		errs[n] = errors.Join(dateErr, amountErr, quantityErr)
		infos = append(infos, hi)
	}
	if err := errors.Join(errs[:]...); err != nil {
		return nil, err
	}
	return infos, nil
}

func parseHCP(s inputSegment) (PricingHCP, error) {
	if err := s.checkLast(15); err != nil {
		return PricingHCP{}, err
	}
	h := PricingHCP{Set: true, Index: s.index}
	var errs [16]error
	h.PricingMethodologyHCP01 = HCPPricingMethodologyHCP01(s.el(1))
	h.AllowedAmountHCP02, errs[2] = parseAmount(s, 2)
	h.SavingsAmountHCP03, errs[3] = parseAmount(s, 3)
	h.OrganizationHCP04 = s.el(4)
	h.PerDiemAmountHCP05, errs[5] = parseAmount(s, 5)
	h.ApprovedDRGHCP06 = s.el(6)
	h.ApprovedAmountHCP07, errs[7] = parseAmount(s, 7)
	h.ProductServiceHCP08 = s.el(8)
	h.ProductQualifierHCP09 = ProductServiceQualifier(s.el(9))
	h.ApprovedHCPCSHCP10 = s.el(10)
	h.UnitOrBasisHCP11 = UnitOrBasisOfMeasurement(s.el(11))
	h.QuantityHCP12, errs[12] = parseAmount(s, 12)
	h.RejectReasonHCP13 = HCPRejectReasonHCP13(s.el(13))
	h.PolicyComplianceHCP14 = HCPPolicyComplianceHCP14(s.el(14))
	h.ExceptionHCP15 = HCPExceptionCodeHCP15(s.el(15))
	if err := errors.Join(errs[:]...); err != nil {
		return PricingHCP{}, err
	}
	return h, nil
}

func parseOI(s inputSegment) (OtherHealthInsuranceOI, error) {
	if err := s.checkLast(6); err != nil {
		return OtherHealthInsuranceOI{}, err
	}
	o := OtherHealthInsuranceOI{Set: true, Index: s.index}
	o.FilingIndicatorOI01 = s.el(1)
	o.SubmissionReasonOI02 = s.el(2)
	o.BenefitsAssignOI03 = OtherInsuranceIndicator(s.el(3))
	o.SignatureSourceOI04 = s.el(4)
	o.ProviderAgreementOI05 = s.el(5)
	o.ReleaseOfInfoOI06 = OtherInsuranceIndicator(s.el(6))
	return o, nil
}

func parseMIA(s inputSegment) (InpatientMIA, error) {
	if err := s.checkLast(24); err != nil {
		return InpatientMIA{}, err
	}
	m := InpatientMIA{Set: true, Index: s.index}
	var errs [25]error
	m.CoveredDaysMIA01, errs[1] = parseAmount(s, 1)
	m.OperatingOutlierMIA02, errs[2] = parseAmount(s, 2)
	m.PsychDaysMIA03, errs[3] = parseAmount(s, 3)
	m.DRGAmountMIA04, errs[4] = parseAmount(s, 4)
	m.RemarkMIA05 = rarc.Code(s.el(5))
	m.DSHAmountMIA06, errs[6] = parseAmount(s, 6)
	m.MSPPassThruMIA07, errs[7] = parseAmount(s, 7)
	m.PPSCapitalMIA08, errs[8] = parseAmount(s, 8)
	m.CapitalFSPMIA09, errs[9] = parseAmount(s, 9)
	m.CapitalHSPMIA10, errs[10] = parseAmount(s, 10)
	m.CapitalDSHMIA11, errs[11] = parseAmount(s, 11)
	m.OldCapitalMIA12, errs[12] = parseAmount(s, 12)
	m.CapitalIMEMIA13, errs[13] = parseAmount(s, 13)
	m.OperHSPSpecMIA14, errs[14] = parseAmount(s, 14)
	m.CostReportDaysMIA15, errs[15] = parseAmount(s, 15)
	m.OperFSPSpecMIA16, errs[16] = parseAmount(s, 16)
	m.PPSOutlierMIA17, errs[17] = parseAmount(s, 17)
	m.IndirectTeachingMIA18, errs[18] = parseAmount(s, 18)
	m.NonPayProfCompMIA19, errs[19] = parseAmount(s, 19)
	m.RemarkMIA20 = rarc.Code(s.el(20))
	m.RemarkMIA21 = rarc.Code(s.el(21))
	m.RemarkMIA22 = rarc.Code(s.el(22))
	m.RemarkMIA23 = rarc.Code(s.el(23))
	m.CapitalExceptionMIA24, errs[24] = parseAmount(s, 24)
	if err := errors.Join(errs[:]...); err != nil {
		return InpatientMIA{}, err
	}
	return m, nil
}

func parseMOA(s inputSegment) (OutpatientMOA, error) {
	if err := s.checkLast(9); err != nil {
		return OutpatientMOA{}, err
	}
	m := OutpatientMOA{Set: true, Index: s.index}
	var errs [10]error
	m.ReimbursementRateMOA01, errs[1] = parseAmount(s, 1)
	m.HCPCSPayableMOA02, errs[2] = parseAmount(s, 2)
	m.RemarkMOA03 = rarc.Code(s.el(3))
	m.RemarkMOA04 = rarc.Code(s.el(4))
	m.RemarkMOA05 = rarc.Code(s.el(5))
	m.RemarkMOA06 = rarc.Code(s.el(6))
	m.RemarkMOA07 = rarc.Code(s.el(7))
	m.ESRDPaymentMOA08, errs[8] = parseAmount(s, 8)
	m.NonpayableProfCompMOA09, errs[9] = parseAmount(s, 9)
	if err := errors.Join(errs[:]...); err != nil {
		return OutpatientMOA{}, err
	}
	return m, nil
}

// parseCAS reads the group code and every reason/amount/quantity triple
// sent, an empty triple between two sent ones included, so each triple
// keeps its position.
func parseCAS(s inputSegment) (AdjustmentCAS, error) {
	if err := s.checkLast(19); err != nil {
		return AdjustmentCAS{}, err
	}
	a := AdjustmentCAS{Index: s.index}
	var errs [20]error
	a.GroupCAS01 = CASGroupCode01(s.el(1))
	for n := 2; n <= len(s.elements); n += 3 {
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

func parseSV2(s inputSegment, component byte) (InstitutionalServiceSV2, error) {
	if err := s.checkLast(10); err != nil {
		return InstitutionalServiceSV2{}, err
	}
	v := InstitutionalServiceSV2{Set: true, Index: s.index}
	var errs [11]error
	v.RevenueCodeSV201 = s.el(1)
	v.ProcedureSV202, errs[2] = parseProcedure(s, 2, component)
	v.ChargeSV203, errs[3] = parseAmount(s, 3)
	v.UnitSV204 = UnitOrBasisOfMeasurement(s.el(4))
	v.UnitsSV205, errs[5] = parseAmount(s, 5)
	v.UnitRateSV206, errs[6] = parseAmount(s, 6)
	v.DeniedChargeSV207, errs[7] = parseAmount(s, 7)
	v.ConditionSV208 = s.el(8)
	v.ResidentialStatusSV209 = s.el(9)
	v.LevelOfCareSV210 = s.el(10)
	if err := errors.Join(errs[:]...); err != nil {
		return InstitutionalServiceSV2{}, err
	}
	return v, nil
}

// parseProcedure reads the composite procedure identifier in element n
// (SV202 or SVD03): qualifier, code, four modifiers, description, and the
// eighth component (Product/Service ID), which the guide marks Not Used.
func parseProcedure(s inputSegment, n int, component byte) (ProcedureSV202, error) {
	parts, err := components(s, n, component, 8)
	if err != nil {
		return ProcedureSV202{}, err
	}
	var p ProcedureSV202
	p.Qualifier = ProductServiceQualifier(nth(parts, 1))
	p.Code = nth(parts, 2)
	p.Modifiers = [4]string{nth(parts, 3), nth(parts, 4), nth(parts, 5), nth(parts, 6)}
	p.Description = nth(parts, 7)
	p.ProductServiceID = nth(parts, 8)
	return p, nil
}

func parseLIN(s inputSegment) (DrugLIN, error) {
	if err := s.checkLast(31); err != nil {
		return DrugLIN{}, err
	}
	d := DrugLIN{Index: s.index}
	d.AssignedIDLIN01 = s.el(1)
	d.QualifierLIN02 = NDCQualifierLIN02(s.el(2))
	d.NDCLIN03 = s.el(3)
	d.ExtraElements = extraElements(s, 3)
	return d, nil
}

// parseCTP reads a CTP. CTP05 is a composite of which the first component
// (the unit code) is named; the other fourteen are kept as sent.
func parseCTP(s inputSegment, component byte) (DrugQuantityCTP, error) {
	if err := s.checkLast(11); err != nil {
		return DrugQuantityCTP{}, err
	}
	c := DrugQuantityCTP{Set: true, Index: s.index}
	var errs [6]error
	c.ClassOfTradeCTP01 = s.el(1)
	c.PriceIdentifierCTP02 = s.el(2)
	c.UnitPriceCTP03, errs[3] = parseAmount(s, 3)
	c.UnitCountCTP04, errs[4] = parseAmount(s, 4)
	var unit []string
	unit, errs[5] = components(s, 5, component, 15)
	c.UnitCTP05_1 = UnitOrBasisOfMeasurement(nth(unit, 1))
	if len(unit) > 1 {
		c.UnitCTP05Rest = append([]string(nil), unit[1:]...)
	}
	c.ExtraElements = extraElements(s, 5)
	if err := errors.Join(errs[:]...); err != nil {
		return DrugQuantityCTP{}, err
	}
	return c, nil
}

func parseSVD(s inputSegment, component byte) (LineAdjudicationSVD, error) {
	if err := s.checkLast(6); err != nil {
		return LineAdjudicationSVD{}, err
	}
	a := LineAdjudicationSVD{Index: s.index}
	var errs [7]error
	a.OtherPayerIDSVD01 = s.el(1)
	a.PaymentSVD02, errs[2] = parseAmount(s, 2)
	a.ProcedureSVD03, errs[3] = parseProcedure(s, 3, component)
	a.RevenueCodeSVD04 = s.el(4)
	a.PaidUnitsSVD05, errs[5] = parseAmount(s, 5)
	a.BundledLineSVD06 = s.el(6)
	if err := errors.Join(errs[:]...); err != nil {
		return LineAdjudicationSVD{}, err
	}
	return a, nil
}

// extraElements returns the elements past position named, as sent (the
// splitter has already dropped trailing empties), or nil when the segment
// ends at or before it.
func extraElements(s inputSegment, named int) []string {
	if len(s.elements) <= named {
		return nil
	}
	return append([]string(nil), s.elements[named:]...)
}

// components splits element n on the component separator, removes trailing
// empty components, and returns an error when a component past last, the
// highest one X12 defines for the composite, carries text.
func components(s inputSegment, n int, component byte, last int) ([]string, error) {
	parts := trimTrailingEmpty(strings.Split(s.el(n), string(component)))
	if len(parts) > last {
		return nil, fmt.Errorf("x12837i: parsing %s in segment %d: the model ends at component %d", s.name(n), s.index, last)
	}
	return parts, nil
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

// amountText reads text as an exact decimal and keeps the text; ref names
// the element or component for errors. Empty text is an Amount with Set
// false.
func amountText(s inputSegment, ref, text string) (Amount, error) {
	if text == "" {
		return Amount{}, nil
	}
	value, err := ParseDecimal([]byte(text))
	if err != nil {
		return Amount{}, fmt.Errorf("x12837i: parsing %s in segment %d: %w", ref, s.index, err)
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
// errors.
func dateText(s inputSegment, ref, text string) (Date, error) {
	if len(text) != 8 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: date is not CCYYMMDD", ref, s.index)
	}
	d := Date{Year: number(text[0:4]), Month: number(text[4:6]), Day: number(text[6:8]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: not a calendar date", ref, s.index)
	}
	return d, nil
}

// parseInterchangeDate reads ISA09, which is YYMMDD. Year is 2000 plus YY,
// the convention backend/pkg/x12835's parser follows from the CMS 835 flat
// file's comment on ISA09.
func parseInterchangeDate(s inputSegment, n int) (Date, error) {
	text := s.el(n)
	if text == "" {
		return Date{}, nil
	}
	if len(text) != 6 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: date is not YYMMDD", s.name(n), s.index)
	}
	d := Date{Year: 2000 + number(text[0:2]), Month: number(text[2:4]), Day: number(text[4:6]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12837i: parsing %s in segment %d: not a calendar date", s.name(n), s.index)
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
	if s.el(n) == "" {
		return Time{}, nil
	}
	return timeText(s, s.name(n), s.el(n))
}

// timeText reads text as HHMM or HHMMSS; ref names the element for errors.
func timeText(s inputSegment, ref, text string) (Time, error) {
	if (len(text) != 4 && len(text) != 6) || !allDigits(text) {
		return Time{}, fmt.Errorf("x12837i: parsing %s in segment %d: time is not HHMM or HHMMSS", ref, s.index)
	}
	t := Time{Hour: number(text[0:2]), Minute: number(text[2:4]), Set: true}
	if len(text) == 6 {
		t.Second = number(text[4:6])
		t.WithSeconds = true
	}
	if t.Hour > 23 || t.Minute > 59 || t.Second > 59 {
		return Time{}, fmt.Errorf("x12837i: parsing %s in segment %d: not a clock time", ref, s.index)
	}
	return t, nil
}

// parseWholeNumber reads element n as a count of digits only.
func parseWholeNumber(s inputSegment, n int) (int, error) {
	text := s.el(n)
	if text == "" || !allDigits(text) {
		return 0, fmt.Errorf("x12837i: parsing %s in segment %d: not a whole number", s.name(n), s.index)
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("x12837i: parsing %s in segment %d: number too large", s.name(n), s.index)
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
