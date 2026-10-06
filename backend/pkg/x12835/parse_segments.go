package x12835

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
// struct, every position the standard defines for the segment. A coded
// element is kept as the text it carried, whether or not the code is on
// its list (Known says which); an amount, date, time, or count that does
// not parse fails the segment, and errors.Join reports every such element
// at once. So each element's line reads: field = parse(element n).

func parseISA(s inputSegment) (InterchangeISA, error) {
	if err := s.checkLast(16); err != nil {
		return InterchangeISA{}, err
	}
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
	if err := s.checkLast(8); err != nil {
		return FunctionalGroupGS{}, err
	}
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
	t.TransactionSetIDST01 = TransactionSetIDST01(s.el(1))
	t.ControlNumberST02 = s.el(2)
	t.ImplementationConventionST03 = s.el(3)
	return t, nil
}

// parseTrailer reads an SE, GE, or IEA: a count in element 01 and a control
// number in element 02.
func parseTrailer(s inputSegment) (int, string, error) {
	if err := s.checkLast(2); err != nil {
		return 0, "", err
	}
	count, err := parseWholeNumber(s, 1)
	if err != nil {
		return 0, "", err
	}
	return count, s.el(2), nil
}

func parseBPR(s inputSegment) (PaymentBPR, error) {
	if err := s.checkLast(21); err != nil {
		return PaymentBPR{}, err
	}
	b := PaymentBPR{Index: s.index, Set: true}
	var errs [22]error
	b.HandlingBPR01 = HandlingBPR01(s.el(1))
	b.TotalPaidBPR02, errs[2] = parseAmount(s, 2)
	b.CreditDebitBPR03 = CreditDebitBPR03(s.el(3))
	b.PayMethodBPR04 = PayMethodBPR04(s.el(4))
	b.PayFormatBPR05 = PayFormatBPR05(s.el(5))
	b.SenderDFIQualifierBPR06 = DFIQualifierBPR06(s.el(6))
	b.SenderDFIBPR07 = s.el(7)
	b.SenderAcctQualifierBPR08 = AccountQualifierBPR08(s.el(8))
	b.SenderAcctBPR09 = s.el(9)
	b.PayerIDBPR10 = s.el(10)
	b.OrigCoSupCodeBPR11 = s.el(11)
	b.ReceiverDFIQualifierBPR12 = DFIQualifierBPR06(s.el(12))
	b.ReceiverDFIBPR13 = s.el(13)
	b.ReceiverAcctQualifierBPR14 = AccountQualifierBPR08(s.el(14))
	b.ReceiverAcctBPR15 = s.el(15)
	b.EffectiveDateBPR16, errs[16] = parseDate(s, 16)
	b.BusinessFunctionBPR17 = s.el(17)
	b.ThirdDFIQualifierBPR18 = s.el(18)
	b.ThirdDFIBPR19 = s.el(19)
	b.ThirdAcctQualifierBPR20 = s.el(20)
	b.ThirdAcctBPR21 = s.el(21)
	if err := errors.Join(errs[:]...); err != nil {
		return PaymentBPR{}, err
	}
	return b, nil
}

func parseTRN(s inputSegment) (TraceTRN, error) {
	if err := s.checkLast(4); err != nil {
		return TraceTRN{}, err
	}
	t := TraceTRN{Index: s.index, Set: true}
	t.TraceTypeTRN01 = TraceTypeTRN01(s.el(1))
	t.TraceNumberTRN02 = s.el(2)
	t.PayerIDTRN03 = s.el(3)
	t.OrigCoSupCodeTRN04 = s.el(4)
	return t, nil
}

func parseCUR(s inputSegment) (CurrencyCUR, error) {
	if err := s.checkLast(21); err != nil {
		return CurrencyCUR{}, err
	}
	rate, err := parseAmount(s, 3)
	if err != nil {
		return CurrencyCUR{}, err
	}
	return CurrencyCUR{
		Index: s.index, Set: true,
		EntityIDCUR01: s.el(1), CurrencyCodeCUR02: s.el(2), ExchangeRateCUR03: rate,
		EntityIDCUR04: s.el(4), CurrencyCodeCUR05: s.el(5), MarketCUR06: s.el(6),
		DateQualifierCUR07: s.el(7), DateCUR08: s.el(8), TimeCUR09: s.el(9),
		DateQualifierCUR10: s.el(10), DateCUR11: s.el(11), TimeCUR12: s.el(12),
		DateQualifierCUR13: s.el(13), DateCUR14: s.el(14), TimeCUR15: s.el(15),
		DateQualifierCUR16: s.el(16), DateCUR17: s.el(17), TimeCUR18: s.el(18),
		DateQualifierCUR19: s.el(19), DateCUR20: s.el(20), TimeCUR21: s.el(21),
	}, nil
}

func parseREF(s inputSegment) (ReferenceREF, error) {
	if err := s.checkLast(4); err != nil {
		return ReferenceREF{}, err
	}
	return ReferenceREF{
		Index:             s.index,
		RefQualifierREF01: RefQualifierREF01(s.el(1)),
		Value:             s.el(2),
		DescriptionREF03:  s.el(3),
		SubReferenceREF04: s.el(4),
	}, nil
}

// parseDTM reads a DTM whole: the qualifier, the date, and the four
// positions after them.
func parseDTM(s inputSegment) (DateTimeDTM, error) {
	if err := s.checkLast(6); err != nil {
		return DateTimeDTM{}, err
	}
	date, err := parseDate(s, 2)
	if err != nil {
		return DateTimeDTM{}, err
	}
	return DateTimeDTM{
		Index:              s.index,
		DateQualifierDTM01: DateQualifierDTM01(s.el(1)),
		Date:               date,
		TimeDTM03:          s.el(3),
		TimeCodeDTM04:      s.el(4),
		PeriodFormatDTM05:  s.el(5),
		PeriodDTM06:        s.el(6),
	}, nil
}

func parseN1(s inputSegment) (PartyN1, error) {
	if err := s.checkLast(6); err != nil {
		return PartyN1{}, err
	}
	return PartyN1{
		Index:            s.index,
		EntityIDN101:     EntityIDN101(s.el(1)),
		NameN102:         s.el(2),
		IDQualifierN103:  IDQualifierN103(s.el(3)),
		IdentifierN104:   s.el(4),
		RelationshipN105: s.el(5),
		EntityIDN106:     s.el(6),
	}, nil
}

func parseN3(s inputSegment) (AddressN3, error) {
	if err := s.checkLast(2); err != nil {
		return AddressN3{}, err
	}
	return AddressN3{Index: s.index, Set: true, Line1N301: s.el(1), Line2N302: s.el(2)}, nil
}

func parseN4(s inputSegment) (LocationN4, error) {
	if err := s.checkLast(7); err != nil {
		return LocationN4{}, err
	}
	return LocationN4{
		Index: s.index, Set: true,
		CityN401: s.el(1), StateN402: s.el(2), ZipN403: s.el(3), CountryN404: s.el(4),
		LocationQualifierN405: s.el(5), LocationN406: s.el(6), CountrySubdivisionN407: s.el(7),
	}, nil
}

func parsePER(s inputSegment) (ContactPER, error) {
	if err := s.checkLast(9); err != nil {
		return ContactPER{}, err
	}
	return ContactPER{
		Index:                 s.index,
		FunctionPER01:         ContactFunctionPER01(s.el(1)),
		NamePER02:             s.el(2),
		QualifierPER03:        CommQualifierPER03(s.el(3)),
		CommPER04:             s.el(4),
		QualifierPER05:        CommQualifierPER03(s.el(5)),
		CommPER06:             s.el(6),
		QualifierPER07:        CommQualifierPER03(s.el(7)),
		CommPER08:             s.el(8),
		InquiryReferencePER09: s.el(9),
	}, nil
}

func parseRDM(s inputSegment) (DeliveryRDM, error) {
	if err := s.checkLast(5); err != nil {
		return DeliveryRDM{}, err
	}
	return DeliveryRDM{
		Index: s.index, Set: true,
		ReportTransmissionRDM01: DeliveryMethodRDM01(s.el(1)),
		NameRDM02:               s.el(2),
		CommNumberRDM03:         s.el(3),
		ReferenceRDM04:          s.el(4),
		ReferenceRDM05:          s.el(5),
	}, nil
}

func parseTS3(s inputSegment) (ProviderSummaryTS3, error) {
	if err := s.checkLast(24); err != nil {
		return ProviderSummaryTS3{}, err
	}
	t := ProviderSummaryTS3{Index: s.index, Set: true}
	var errs [25]error
	t.ProviderIDTS301 = s.el(1)
	t.FacilityCodeTS302 = s.el(2)
	t.FiscalPeriodDateTS303, errs[3] = parseDate(s, 3)
	t.ClaimCountTS304, errs[4] = parseAmount(s, 4)
	t.ClaimChangeAmountTS305, errs[5] = parseAmount(s, 5)
	t.CoveredChargeAmountTS306, errs[6] = parseAmount(s, 6)
	t.NoncoveredChargeAmountTS307, errs[7] = parseAmount(s, 7)
	t.DeniedChargeAmountTS308, errs[8] = parseAmount(s, 8)
	t.ProviderAmountTS309, errs[9] = parseAmount(s, 9)
	t.InterestAmountTS310, errs[10] = parseAmount(s, 10)
	t.ContractualAdjustmentAmountTS311, errs[11] = parseAmount(s, 11)
	t.GrammRudmanReductionAmountTS312, errs[12] = parseAmount(s, 12)
	t.MSPPayerAmountTS313, errs[13] = parseAmount(s, 13)
	t.BloodDeductibleAmountTS314, errs[14] = parseAmount(s, 14)
	t.NonLabChargeAmountTS315, errs[15] = parseAmount(s, 15)
	t.CoinsuranceAmountTS316, errs[16] = parseAmount(s, 16)
	t.HCPCSReportedChargeAmountTS317, errs[17] = parseAmount(s, 17)
	t.HCPCSPayableAmountTS318, errs[18] = parseAmount(s, 18)
	t.DeductibleAmountTS319, errs[19] = parseAmount(s, 19)
	t.ProfessionalComponentAmountTS320, errs[20] = parseAmount(s, 20)
	t.MSPPatientLiabilityMetAmountTS321, errs[21] = parseAmount(s, 21)
	t.PatientReimbursementAmountTS322, errs[22] = parseAmount(s, 22)
	t.PIPClaimCountTS323, errs[23] = parseAmount(s, 23)
	t.PIPAdjustmentAmountTS324, errs[24] = parseAmount(s, 24)
	if err := errors.Join(errs[:]...); err != nil {
		return ProviderSummaryTS3{}, err
	}
	return t, nil
}

func parseTS2(s inputSegment) (InstitutionalSummaryTS2, error) {
	if err := s.checkLast(19); err != nil {
		return InstitutionalSummaryTS2{}, err
	}
	t := InstitutionalSummaryTS2{Index: s.index, Set: true}
	var errs [20]error
	t.DRGAmountTS201, errs[1] = parseAmount(s, 1)
	t.FederalSpecificAmountTS202, errs[2] = parseAmount(s, 2)
	t.HospitalSpecificAmountTS203, errs[3] = parseAmount(s, 3)
	t.DisproportionateAmountTS204, errs[4] = parseAmount(s, 4)
	t.CapitalAmountTS205, errs[5] = parseAmount(s, 5)
	t.IndirectMedicalEducationAmountTS206, errs[6] = parseAmount(s, 6)
	t.OutlierDayCountTS207, errs[7] = parseAmount(s, 7)
	t.DayOutlierAmountTS208, errs[8] = parseAmount(s, 8)
	t.CostOutlierAmountTS209, errs[9] = parseAmount(s, 9)
	t.AverageDRGLengthOfStayTS210, errs[10] = parseAmount(s, 10)
	t.DischargeCountTS211, errs[11] = parseAmount(s, 11)
	t.CostReportDayCountTS212, errs[12] = parseAmount(s, 12)
	t.CoveredDayCountTS213, errs[13] = parseAmount(s, 13)
	t.NoncoveredDayCountTS214, errs[14] = parseAmount(s, 14)
	t.MSPPassThroughAmountTS215, errs[15] = parseAmount(s, 15)
	t.AverageDRGWeightTS216, errs[16] = parseAmount(s, 16)
	t.PPSCapitalFSPDRGAmountTS217, errs[17] = parseAmount(s, 17)
	t.PSPCapitalHSPDRGAmountTS218, errs[18] = parseAmount(s, 18)
	t.PPSDSHDRGAmountTS219, errs[19] = parseAmount(s, 19)
	if err := errors.Join(errs[:]...); err != nil {
		return InstitutionalSummaryTS2{}, err
	}
	return t, nil
}

func parseCLP(s inputSegment) (ClaimCLP, error) {
	if err := s.checkLast(14); err != nil {
		return ClaimCLP{}, err
	}
	c := ClaimCLP{Index: s.index}
	var errs [15]error
	c.PatientControlCLP01 = s.el(1)
	c.StatusCLP02 = ClaimStatusCLP02(s.el(2))
	c.ChargeCLP03, errs[3] = parseAmount(s, 3)
	c.PaymentCLP04, errs[4] = parseAmount(s, 4)
	c.PatientRespCLP05, errs[5] = parseAmount(s, 5)
	c.FilingIndicatorCLP06 = FilingIndicatorCLP06(s.el(6))
	c.PayerClaimControlCLP07 = s.el(7)
	c.FacilityCLP08 = s.el(8)
	c.FrequencyCLP09 = s.el(9)
	c.PatientStatusCLP10 = s.el(10)
	c.DRGCLP11 = s.el(11)
	c.DRGWeightCLP12, errs[12] = parseAmount(s, 12)
	c.DischargeFractionCLP13, errs[13] = parseAmount(s, 13)
	c.YesNoCLP14 = s.el(14)
	if err := errors.Join(errs[:]...); err != nil {
		return ClaimCLP{}, err
	}
	return c, nil
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
	a.GroupCAS01 = AdjGroupCAS01(s.el(1))
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

func parseNM1(s inputSegment) (NameNM1, error) {
	if err := s.checkLast(12); err != nil {
		return NameNM1{}, err
	}
	return NameNM1{
		Index:             s.index,
		EntityIDNM101:     EntityIDNM101(s.el(1)),
		EntityTypeNM102:   EntityTypeNM102(s.el(2)),
		LastNM103:         s.el(3),
		FirstNM104:        s.el(4),
		MiddleNM105:       s.el(5),
		PrefixNM106:       s.el(6),
		SuffixNM107:       s.el(7),
		IDQualifierNM108:  IDQualifierNM108(s.el(8)),
		IDNM109:           s.el(9),
		RelationshipNM110: s.el(10),
		EntityIDNM111:     s.el(11),
		LastNM112:         s.el(12),
	}, nil
}

func parseMIA(s inputSegment) (InpatientMIA, error) {
	if err := s.checkLast(24); err != nil {
		return InpatientMIA{}, err
	}
	m := InpatientMIA{Index: s.index, Set: true}
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
	m := OutpatientMOA{Index: s.index, Set: true}
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

func parseAMT(s inputSegment) (AmountAMT, error) {
	if err := s.checkLast(3); err != nil {
		return AmountAMT{}, err
	}
	amount, err := parseAmount(s, 2)
	if err != nil {
		return AmountAMT{}, err
	}
	return AmountAMT{
		Index:                s.index,
		AmountQualifierAMT01: AmountQualifierAMT01(s.el(1)),
		Amount:               amount,
		CreditDebitAMT03:     s.el(3),
	}, nil
}

func parseQTY(s inputSegment) (QuantityQTY, error) {
	if err := s.checkLast(4); err != nil {
		return QuantityQTY{}, err
	}
	amount, err := parseAmount(s, 2)
	if err != nil {
		return QuantityQTY{}, err
	}
	return QuantityQTY{
		Index:                  s.index,
		QuantityQualifierQTY01: QuantityQualifierQTY01(s.el(1)),
		Amount:                 amount,
		UnitOfMeasureQTY03:     s.el(3),
		FreeFormQTY04:          s.el(4),
	}, nil
}

func parseSVC(s inputSegment, component byte) (ServiceLineSVC, error) {
	if err := s.checkLast(7); err != nil {
		return ServiceLineSVC{}, err
	}
	l := ServiceLineSVC{Index: s.index}
	var errs [8]error
	l.ProcedureSVC01, errs[1] = parseProcedure(s, 1, component)
	l.ChargeSVC02, errs[2] = parseAmount(s, 2)
	l.PaymentSVC03, errs[3] = parseAmount(s, 3)
	l.RevenueCodeSVC04 = s.el(4)
	l.UnitsSVC05, errs[5] = parseAmount(s, 5)
	l.OriginalProcedureSVC06, errs[6] = parseProcedure(s, 6, component)
	l.OriginalUnitsSVC07, errs[7] = parseAmount(s, 7)
	if err := errors.Join(errs[:]...); err != nil {
		return ServiceLineSVC{}, err
	}
	return l, nil
}

// parseProcedure reads the composite procedure identifier in element n
// (SVC01 or SVC06): qualifier, code, four modifiers, description, and the
// second code, the eight components the standard defines.
func parseProcedure(s inputSegment, n int, component byte) (ProcedureSVC01, error) {
	parts := trimTrailingEmpty(strings.Split(s.el(n), string(component)))
	if len(parts) > 8 {
		return ProcedureSVC01{}, fmt.Errorf("x12835: parsing %s in segment %d: the composite ends at component 8", s.name(n), s.index)
	}
	return ProcedureSVC01{
		Qualifier:   ProductQualifierSVC01(nth(parts, 1)),
		Code:        nth(parts, 2),
		Modifiers:   [4]string{nth(parts, 3), nth(parts, 4), nth(parts, 5), nth(parts, 6)},
		Description: nth(parts, 7),
		Code2:       nth(parts, 8),
	}, nil
}

func parseLQ(s inputSegment) (RemarkLQ, error) {
	if err := s.checkLast(2); err != nil {
		return RemarkLQ{}, err
	}
	return RemarkLQ{
		Index:               s.index,
		RemarkQualifierLQ01: RemarkQualifierLQ01(s.el(1)),
		RARC:                rarc.Code(s.el(2)),
	}, nil
}

// parsePLB reads the provider, the fiscal period, and every adjustment pair
// sent, an empty pair between two sent ones included. Each pair's identifier
// is a composite: reason code, then free text.
func parsePLB(s inputSegment, component byte) (ProviderAdjustmentPLB, error) {
	if err := s.checkLast(14); err != nil {
		return ProviderAdjustmentPLB{}, err
	}
	plb := ProviderAdjustmentPLB{Index: s.index}
	var errs [15]error
	plb.ProviderPLB01 = s.el(1)
	plb.FiscalPeriodPLB02, errs[2] = parseDate(s, 2)
	for n := 3; n <= len(s.elements); n += 2 {
		parts := trimTrailingEmpty(strings.Split(s.el(n), string(component)))
		if len(parts) > 2 {
			errs[n] = fmt.Errorf("x12835: parsing %s in segment %d: the composite ends at component 2", s.name(n), s.index)
		}
		var a PLBAdjustment
		a.Reason = PLBReasonPLB03(nth(parts, 1))
		a.Identifier = nth(parts, 2)
		a.Amount, errs[n+1] = parseAmount(s, n+1)
		plb.Adjustments = append(plb.Adjustments, a)
	}
	if err := errors.Join(errs[:]...); err != nil {
		return ProviderAdjustmentPLB{}, err
	}
	return plb, nil
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

// parseAmount reads element n as an exact decimal, keeping the text it
// came as. An empty element is an Amount with Set false.
func parseAmount(s inputSegment, n int) (Amount, error) {
	text := s.el(n)
	if text == "" {
		return Amount{}, nil
	}
	value, err := ParseDecimal([]byte(text))
	if err != nil {
		return Amount{}, fmt.Errorf("x12835: parsing %s in segment %d: %w", s.name(n), s.index, err)
	}
	return Amount{Value: value, Set: true, Text: text}, nil
}

// parseDate reads element n as CCYYMMDD. An empty element is a Date with
// Set false.
func parseDate(s inputSegment, n int) (Date, error) {
	text := s.el(n)
	if text == "" {
		return Date{}, nil
	}
	if len(text) != 8 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12835: parsing %s in segment %d: date is not CCYYMMDD", s.name(n), s.index)
	}
	d := Date{Year: number(text[0:4]), Month: number(text[4:6]), Day: number(text[6:8]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12835: parsing %s in segment %d: not a calendar date", s.name(n), s.index)
	}
	return d, nil
}

// parseInterchangeDate reads ISA09, which is YYMMDD. The flat file's
// comment on ISA09 says the century 20 is appended (CMS 835 Flat File, PDF
// p.1), so Year is 2000 plus YY.
func parseInterchangeDate(s inputSegment, n int) (Date, error) {
	text := s.el(n)
	if text == "" {
		return Date{}, nil
	}
	if len(text) != 6 || !allDigits(text) {
		return Date{}, fmt.Errorf("x12835: parsing %s in segment %d: date is not YYMMDD", s.name(n), s.index)
	}
	d := Date{Year: 2000 + number(text[0:2]), Month: number(text[2:4]), Day: number(text[4:6]), Set: true}
	if !calendarDate(d) {
		return Date{}, fmt.Errorf("x12835: parsing %s in segment %d: not a calendar date", s.name(n), s.index)
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
		return Time{}, fmt.Errorf("x12835: parsing %s in segment %d: time is not HHMM or HHMMSS", s.name(n), s.index)
	}
	t := Time{Hour: number(text[0:2]), Minute: number(text[2:4]), Set: true}
	if len(text) == 6 {
		t.Second = number(text[4:6])
		t.HasSeconds = true
	}
	if t.Hour > 23 || t.Minute > 59 || t.Second > 59 {
		return Time{}, fmt.Errorf("x12835: parsing %s in segment %d: not a clock time", s.name(n), s.index)
	}
	return t, nil
}

// parseWholeNumber reads element n as a count of digits only.
func parseWholeNumber(s inputSegment, n int) (int, error) {
	text := s.el(n)
	if text == "" || !allDigits(text) {
		return 0, fmt.Errorf("x12835: parsing %s in segment %d: not a whole number", s.name(n), s.index)
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("x12835: parsing %s in segment %d: number too large", s.name(n), s.index)
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
