package x12837i

import (
	"fmt"
	"sort"
	"strings"
)

// outputSegment is one segment ready to write: its position, its ID, and
// its elements with trailing empties already dropped.
type outputSegment struct {
	index    int
	id       string
	elements []string
}

// collector gathers the segments of a tree, in any order, for ordered to
// sort by position. HI composites are gathered by segment and put back
// together at the end, since the claim keeps them in separate fields.
type collector struct {
	component  byte
	segments   []outputSegment
	healthInfo map[int][]HealthInfoHI
}

// add records one segment, dropping trailing empty elements.
func (c *collector) add(index int, id string, elements ...string) {
	c.segments = append(c.segments, outputSegment{index: index, id: id, elements: trimTrailingEmpty(elements)})
}

// addFixed records one segment with every element, empty or not: the ISA,
// whose sixteen elements are fixed width.
func (c *collector) addFixed(index int, id string, elements ...string) {
	c.segments = append(c.segments, outputSegment{index: index, id: id, elements: elements})
}

// raw records segments the model did not place.
func (c *collector) raw(segments []RawSegment) {
	for _, r := range segments {
		c.add(r.Index, r.ID, r.Elements...)
	}
}

// join builds one composite element from its components, dropping trailing
// empties.
func (c *collector) join(parts ...string) string {
	return strings.Join(trimTrailingEmpty(parts), string(c.component))
}

// ordered puts the HI segments together, sorts every segment by position,
// and checks that the positions run 0, 1, 2, ... with none missing and
// none repeated.
func (c *collector) ordered() ([]outputSegment, error) {
	for index, infos := range c.healthInfo {
		last := 0
		for _, hi := range infos {
			if hi.Position > last {
				last = hi.Position
			}
		}
		elements := make([]string, last)
		for _, hi := range infos {
			if hi.Position < 1 {
				return nil, fmt.Errorf("x12837i: writing: an HI composite at segment %d has no position", index)
			}
			elements[hi.Position-1] = c.healthInfoText(hi)
		}
		c.add(index, "HI", elements...)
	}
	sort.Slice(c.segments, func(i, j int) bool { return c.segments[i].index < c.segments[j].index })
	for i, s := range c.segments {
		if s.index > i {
			return nil, fmt.Errorf("x12837i: writing: no segment at position %d", i)
		}
		if s.index < i {
			return nil, fmt.Errorf("x12837i: writing: two segments at position %d", s.index)
		}
	}
	return c.segments, nil
}

// ---------------------------------------------------------------- element text

// dateText writes a Date as CCYYMMDD, or "" when it is not set.
func dateOut(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}

// interchangeDateOut writes ISA09's date as YYMMDD, the form it came in.
func interchangeDateOut(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}

// timeOut writes a Time as HHMM or HHMMSS, whichever the file sent, or ""
// when it is not set.
func timeOut(t Time) string {
	if !t.Set {
		return ""
	}
	if t.WithSeconds {
		return fmt.Sprintf("%02d%02d%02d", t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d%02d", t.Hour, t.Minute)
}

// amountOut writes an Amount as the text the file sent, or the decimal's
// own spelling for an Amount built without one, or "" when it is not set.
func amountOut(a Amount) string {
	if !a.Set {
		return ""
	}
	if a.Text != "" {
		return a.Text
	}
	return a.Value.String()
}

// periodOut writes a date or a date range in the form the format qualifier
// names: RD8 as CCYYMMDD-CCYYMMDD, anything else as the one date.
func periodOut(format string, from, through Date) string {
	if !from.Set {
		return ""
	}
	if format == "RD8" {
		return dateOut(from) + "-" + dateOut(through)
	}
	return dateOut(from)
}

// dtp03Out writes DTP03 in the form DTP02 names: TM an hour, DT a date and
// an hour, RD8 a range, and D8 a date.
func dtp03Out(d DateDTP) string {
	switch d.FormatDTP02 {
	case "TM":
		return timeOut(d.Time)
	case "DT":
		if !d.Date.Set {
			return ""
		}
		return dateOut(d.Date) + fmt.Sprintf("%02d%02d", d.Time.Hour, d.Time.Minute)
	}
	return periodOut(d.FormatDTP02, d.Date, d.DateThrough)
}

func (c *collector) procedureText(p ProcedureSV202) string {
	return c.join(string(p.Qualifier), p.Code, p.Modifiers[0], p.Modifiers[1], p.Modifiers[2], p.Modifiers[3], p.Description, p.ProductServiceID)
}

func (c *collector) healthInfoText(hi HealthInfoHI) string {
	return c.join(string(hi.QualifierHI01_1), hi.CodeHI01_2, hi.DateQualifierHI01_3,
		periodOut(hi.DateQualifierHI01_3, hi.DateHI01_4, hi.DateThroughHI01_4),
		amountOut(hi.AmountHI01_5), amountOut(hi.QuantityHI01_6),
		hi.VersionIdentifierHI0n_7, hi.IndustryCodeHI0n_8, string(hi.PresentOnAdmissionHI0n_9))
}

// ---------------------------------------------------------------- segments

func (c *collector) interchange(isa InterchangeISA) {
	c.addFixed(isa.Index, "ISA",
		isa.AuthQualifierISA01, isa.AuthInfoISA02, isa.SecQualifierISA03, isa.SecInfoISA04,
		isa.SenderIDQualifierISA05, isa.SenderIDISA06, isa.ReceiverIDQualifierISA07, isa.ReceiverIDISA08,
		interchangeDateOut(isa.DateISA09), timeOut(isa.TimeISA10), isa.RepetitionSepISA11, isa.VersionISA12,
		isa.ControlNumberISA13, isa.AckRequestedISA14, isa.UsageISA15, isa.ComponentSepISA16)
	for _, g := range isa.FunctionalGroups {
		c.group(g)
	}
	c.raw(isa.Raw)
	c.add(isa.IndexIEA, "IEA", isa.GroupCountIEA01, isa.ControlNumberIEA02)
}

func (c *collector) group(g FunctionalGroupGS) {
	c.add(g.Index, "GS", string(g.FunctionalIDGS01), g.SenderCodeGS02, g.ReceiverCodeGS03,
		dateOut(g.DateGS04), timeOut(g.TimeGS05), g.ControlNumberGS06, g.AgencyGS07, g.VersionGS08)
	for _, t := range g.Transactions {
		c.transaction(t)
	}
	c.raw(g.Raw)
	c.add(g.IndexGE, "GE", g.TransactionCountGE01, g.ControlNumberGE02)
}

func (c *collector) transaction(t TransactionST) {
	c.add(t.Index, "ST", string(t.TransactionSetIDST01), t.ControlNumberST02, t.VersionST03)
	if t.Header.Set {
		h := t.Header
		c.add(h.Index, "BHT", h.StructureCodeBHT01, h.PurposeCodeBHT02, h.OriginatorTransactionIDBHT03,
			dateOut(h.CreationDateBHT04), timeOut(h.CreationTimeBHT05), string(h.ClaimOrEncounterBHT06))
	}
	if placed(t.SubmitterNM141.Name) {
		c.nm1(t.SubmitterNM141.Name)
		for _, per := range t.SubmitterNM141.Contacts {
			c.per(per)
		}
		c.raw(t.SubmitterNM141.Raw)
	}
	if placed(t.ReceiverNM140.Name) {
		c.nm1(t.ReceiverNM140.Name)
		c.raw(t.ReceiverNM140.Raw)
	}
	for _, b := range t.BillingProviders {
		c.billingProvider(b)
	}
	c.raw(t.Raw)
	c.add(t.IndexSE, "SE", t.SegmentCountSE01, t.ControlNumberSE02)
}

func (c *collector) nm1(n NameNM1) {
	elements := []string{string(n.EntityIDNM101), string(n.EntityTypeNM102), n.LastNM103, n.FirstNM104, n.MiddleNM105,
		n.PrefixNM106, n.SuffixNM107, string(n.IDQualifierNM108), n.IDNM109}
	c.add(n.Index, "NM1", append(elements, n.ExtraElements...)...)
}

func (c *collector) per(p ContactPER) {
	elements := []string{p.FunctionPER01, p.NamePER02, p.QualifierPER03, p.CommPER04, p.QualifierPER05, p.CommPER06, p.QualifierPER07, p.CommPER08}
	c.add(p.Index, "PER", append(elements, p.ExtraElements...)...)
}

func (c *collector) n3(a AddressN3) {
	if a.Set {
		c.add(a.Index, "N3", a.Line1N301, a.Line2N302)
	}
}

func (c *collector) n4(l LocationN4) {
	if l.Set {
		c.add(l.Index, "N4", l.CityN401, l.StateN402, l.ZipN403, l.CountryN404, l.LocationQualifierN405, l.LocationIdentifierN406, l.SubdivisionN407)
	}
}

// ref records one REF when it was placed: a REF always carries its
// qualifier, so an empty one is a field that was never filled.
func (c *collector) ref(r ReferenceREF) {
	if r.QualifierREF01 != "" {
		c.add(r.Index, "REF", append([]string{r.QualifierREF01, r.Value}, r.ExtraElements...)...)
	}
}

func (c *collector) refs(rs []ReferenceREF) {
	for _, r := range rs {
		c.ref(r)
	}
}

func (c *collector) hl(h HierarchyHL) {
	c.add(h.Index, "HL", h.IDNumberHL01, h.ParentIDNumberHL02, string(h.LevelCodeHL03), h.ChildCodeHL04)
}

func (c *collector) prv(p ProviderPRV) {
	if p.Set {
		c.add(p.Index, "PRV", append([]string{string(p.ProviderCodePRV01), p.RefQualifierPRV02, p.TaxonomyCodePRV03}, p.ExtraElements...)...)
	}
}

func (c *collector) cur(u CurrencyCUR) {
	if u.Set {
		c.add(u.Index, "CUR", append([]string{u.EntityIDCUR01, u.CurrencyCodeCUR02}, u.ExtraElements...)...)
	}
}

func (c *collector) sbr(b SubscriberSBR) {
	if b.Set {
		c.add(b.Index, "SBR", string(b.PayerResponsibilitySBR01), string(b.RelationshipSBR02), b.GroupNumberSBR03, b.GroupNameSBR04,
			b.InsuranceTypeSBR05, b.CoordinationSBR06, b.ResponseSBR07, b.EmploymentSBR08, string(b.ClaimFilingIndicatorSBR09))
	}
}

func (c *collector) dmg(d DemographicDMG) {
	if d.Set {
		c.add(d.Index, "DMG", append([]string{d.DateFormatDMG01, dateOut(d.BirthDateDMG02), string(d.SexDMG03)}, d.ExtraElements...)...)
	}
}

func (c *collector) pat(p PatientRelationshipPAT) {
	if p.Set {
		c.add(p.Index, "PAT", append([]string{string(p.RelationshipPAT01)}, p.ExtraElements...)...)
	}
}

func (c *collector) dtp(d DateDTP) {
	if d.Set {
		c.add(d.Index, "DTP", string(d.QualifierDTP01), d.FormatDTP02, dtp03Out(d))
	}
}

func (c *collector) dtps(ds []DateDTP) {
	for _, d := range ds {
		c.dtp(d)
	}
}

// amt records one AMT when it was placed: an AMT always carries its
// qualifier, so an empty one is a field that was never filled.
func (c *collector) amt(a AmountAMT) {
	if a.QualifierAMT01 != "" {
		c.add(a.Index, "AMT", append([]string{a.QualifierAMT01, amountOut(a.Amount)}, a.ExtraElements...)...)
	}
}

func (c *collector) amts(as []AmountAMT) {
	for _, a := range as {
		c.amt(a)
	}
}

func (c *collector) pwks(ps []SupplementalPWK) {
	for _, p := range ps {
		c.add(p.Index, "PWK", append([]string{string(p.ReportTypePWK01), string(p.TransmissionPWK02), p.CopiesNeededPWK03,
			p.EntityIDPWK04, p.IDQualifierPWK05, p.AttachmentControlPWK06}, p.ExtraElements...)...)
	}
}

func (c *collector) hcp(h PricingHCP) {
	if h.Set {
		c.add(h.Index, "HCP", string(h.PricingMethodologyHCP01), amountOut(h.AllowedAmountHCP02), amountOut(h.SavingsAmountHCP03),
			h.OrganizationHCP04, amountOut(h.PerDiemAmountHCP05), h.ApprovedDRGHCP06, amountOut(h.ApprovedAmountHCP07),
			h.ProductServiceHCP08, string(h.ProductQualifierHCP09), h.ApprovedHCPCSHCP10, string(h.UnitOrBasisHCP11),
			amountOut(h.QuantityHCP12), string(h.RejectReasonHCP13), string(h.PolicyComplianceHCP14), string(h.ExceptionHCP15))
	}
}

func (c *collector) cas(as []AdjustmentCAS) {
	for _, a := range as {
		elements := []string{string(a.GroupCAS01)}
		for _, r := range a.Reasons {
			elements = append(elements, string(r.Reason), amountOut(r.Amount), amountOut(r.Quantity))
		}
		c.add(a.Index, "CAS", elements...)
	}
}

func (c *collector) billingProvider(b BillingProviderHL) {
	c.hl(b.Hierarchy)
	c.prv(b.Provider)
	c.cur(b.Currency)
	if placed(b.BillingProviderNM185) {
		c.nm1(b.BillingProviderNM185)
	}
	c.n3(b.Address)
	c.n4(b.Location)
	c.ref(b.TaxIDREFEI)
	c.refs(b.References)
	for _, per := range b.Contacts {
		c.per(per)
	}
	if placed(b.PayToAddressNM187) {
		c.nm1(b.PayToAddressNM187)
	}
	c.n3(b.PayToAddress)
	c.n4(b.PayToLocation)
	for _, s := range b.Subscribers {
		c.subscriber(s)
	}
	c.raw(b.Raw)
}

func (c *collector) subscriber(s SubscriberHL) {
	c.hl(s.Hierarchy)
	c.sbr(s.Subscriber)
	if placed(s.SubscriberNM1IL) {
		c.nm1(s.SubscriberNM1IL)
	}
	c.n3(s.Address)
	c.n4(s.Location)
	c.dmg(s.Demographics)
	c.ref(s.SecondaryIDREFSY)
	c.ref(s.ClaimNumberREFY4)
	c.refs(s.References)
	if placed(s.PayerNM1PR) {
		c.nm1(s.PayerNM1PR)
	}
	c.n3(s.PayerAddress)
	c.n4(s.PayerLocation)
	c.refs(s.PayerReferences)
	for _, p := range s.Patients {
		c.patient(p)
	}
	for _, cl := range s.Claims {
		c.claim(cl)
	}
	c.raw(s.Raw)
}

func (c *collector) patient(p PatientHL) {
	c.hl(p.Hierarchy)
	c.pat(p.Patient)
	if placed(p.PatientNM1QC) {
		c.nm1(p.PatientNM1QC)
	}
	c.n3(p.Address)
	c.n4(p.Location)
	c.dmg(p.Demographics)
	c.refs(p.References)
	for _, cl := range p.Claims {
		c.claim(cl)
	}
	c.raw(p.Raw)
}

// providerLoop records an NM1-led provider loop: the name, its PRV when it
// has one, its REFs, and its raw segments.
func (c *collector) providerLoop(name NameNM1, specialty ProviderPRV, refs []ReferenceREF, raw []RawSegment) {
	if !placed(name) {
		return
	}
	c.nm1(name)
	c.prv(specialty)
	c.refs(refs)
	c.raw(raw)
}

func (c *collector) healthInfos(infos ...HealthInfoHI) {
	for _, hi := range infos {
		if hi.Set {
			c.healthInfo[hi.Index] = append(c.healthInfo[hi.Index], hi)
		}
	}
}

func (c *collector) claim(cl ClaimCLM) {
	c.add(cl.Index, "CLM", cl.PatientControlCLM01, amountOut(cl.TotalChargeCLM02), cl.FilingIndicatorCLM03, cl.ClaimTypeCLM04,
		c.join(cl.FacilityTypeCLM05_1, string(cl.FacilityCodeQualifierCLM05_2), cl.ClaimFrequencyCLM05_3),
		cl.SignatureIndicatorCLM06, string(cl.AssignmentCLM07), string(cl.BenefitsAssignCLM08), string(cl.ReleaseOfInfoCLM09),
		cl.SignatureSourceCLM10, cl.RelatedCausesCLM11, cl.SpecialProgramCLM12, cl.YesNoCLM13, cl.LevelOfServiceCLM14,
		cl.YesNoCLM15, cl.ParticipationAgreementCLM16, cl.StatusCLM17, cl.YesNoCLM18, cl.SubmissionReasonCLM19,
		string(cl.DelayReasonCLM20))
	c.dtp(cl.DischargeHourDTP096)
	c.dtp(cl.StatementDatesDTP434)
	c.dtp(cl.AdmissionDTP435)
	c.dtp(cl.RepricerReceivedDTP050)
	if cl.InstitutionalCL1.Set {
		l := cl.InstitutionalCL1
		c.add(l.Index, "CL1", l.AdmissionTypeCL101, l.AdmissionSourceCL102, l.PatientStatusCL103, l.NursingHomeCL104)
	}
	c.pwks(cl.SupplementalInfo)
	for _, n := range cl.Contracts {
		c.add(n.Index, "CN1", n.ContractTypeCN101, amountOut(n.ContractAmountCN102), amountOut(n.ContractPercentCN103),
			n.ContractCodeCN104, amountOut(n.TermsDiscountCN105), n.ContractVersionCN106)
	}
	c.amt(cl.PatientEstimatedAmountAMTF3)
	c.amts(cl.Amounts)
	c.refs([]ReferenceREF{cl.ServiceAuthExceptionREF4N, cl.ReferralNumberREF9F, cl.PriorAuthREFG1, cl.PayerClaimControlREFF8,
		cl.RepricedClaimREF9A, cl.AdjustedRepricedClaimREF9C, cl.InvestigationalDeviceREFLX, cl.TransmissionIntermediaryREFD9,
		cl.AutoAccidentStateREFLU, cl.MedicalRecordREFEA, cl.DemonstrationProjectREFP4, cl.PROApprovalREFG4})
	c.refs(cl.References)
	for _, k := range cl.Files {
		c.add(k.Index, "K3", append([]string{k.FixedFormatK301}, k.ExtraElements...)...)
	}
	for _, n := range cl.Notes {
		c.add(n.Index, "NTE", string(n.NoteCodeNTE01), n.TextNTE02)
	}
	for _, r := range cl.EPSDTReferralsCRCZZ {
		c.add(r.Index, "CRC", append([]string{string(r.CategoryCRC01), string(r.CertificationCRC02), string(r.Condition1CRC03),
			string(r.Condition2CRC04), string(r.Condition3CRC05)}, r.ExtraElements...)...)
	}
	c.healthInfos(cl.PrincipalDiagnosisHIABK, cl.AdmittingDiagnosisHIABJ, cl.PrincipalProcedureHIBBR, cl.DRGHIDR)
	c.healthInfos(cl.PatientReasonsHIAPR...)
	c.healthInfos(cl.ExternalCausesHIABN...)
	c.healthInfos(cl.OtherDiagnosesHIABF...)
	c.healthInfos(cl.OtherProceduresHIBBQ...)
	c.healthInfos(cl.OccurrenceSpansHIBI...)
	c.healthInfos(cl.OccurrenceCodesHIBH...)
	c.healthInfos(cl.ValueCodesHIBE...)
	c.healthInfos(cl.ConditionCodesHIBG...)
	c.healthInfos(cl.TreatmentCodesHITC...)
	for _, h := range cl.Pricing {
		c.hcp(h)
	}
	c.providerLoop(cl.AttendingNM171.Name, cl.AttendingNM171.Specialty, cl.AttendingNM171.References, cl.AttendingNM171.Raw)
	c.providerLoop(cl.OperatingNM172.Name, ProviderPRV{}, cl.OperatingNM172.References, cl.OperatingNM172.Raw)
	c.providerLoop(cl.OtherOperatingNM1ZZ.Name, ProviderPRV{}, cl.OtherOperatingNM1ZZ.References, cl.OtherOperatingNM1ZZ.Raw)
	c.providerLoop(cl.RenderingNM182.Name, ProviderPRV{}, cl.RenderingNM182.References, cl.RenderingNM182.Raw)
	c.providerLoop(cl.ServiceFacilityNM177.Name, ProviderPRV{}, cl.ServiceFacilityNM177.References, cl.ServiceFacilityNM177.Raw)
	c.providerLoop(cl.ReferringNM1DN.Name, ProviderPRV{}, cl.ReferringNM1DN.References, cl.ReferringNM1DN.Raw)
	for _, o := range cl.OtherSubscribers {
		c.otherSubscriber(o)
	}
	for _, l := range cl.ServiceLines {
		c.serviceLine(l)
	}
	c.raw(cl.Raw)
}

func (c *collector) otherSubscriber(o OtherSubscriberSBR) {
	c.sbr(o.SBR)
	c.cas(o.Adjustments)
	if o.MOA.Set {
		m := o.MOA
		c.add(m.Index, "MOA", amountOut(m.ReimbursementRateMOA01), amountOut(m.HCPCSPayableMOA02), string(m.RemarkMOA03),
			string(m.RemarkMOA04), string(m.RemarkMOA05), string(m.RemarkMOA06), string(m.RemarkMOA07),
			amountOut(m.ESRDPaymentMOA08), amountOut(m.NonpayableProfCompMOA09))
	}
	if o.MIA.Set {
		m := o.MIA
		c.add(m.Index, "MIA", amountOut(m.CoveredDaysMIA01), amountOut(m.OperatingOutlierMIA02), amountOut(m.PsychDaysMIA03),
			amountOut(m.DRGAmountMIA04), string(m.RemarkMIA05), amountOut(m.DSHAmountMIA06), amountOut(m.MSPPassThruMIA07),
			amountOut(m.PPSCapitalMIA08), amountOut(m.CapitalFSPMIA09), amountOut(m.CapitalHSPMIA10), amountOut(m.CapitalDSHMIA11),
			amountOut(m.OldCapitalMIA12), amountOut(m.CapitalIMEMIA13), amountOut(m.OperHSPSpecMIA14), amountOut(m.CostReportDaysMIA15),
			amountOut(m.OperFSPSpecMIA16), amountOut(m.PPSOutlierMIA17), amountOut(m.IndirectTeachingMIA18), amountOut(m.NonPayProfCompMIA19),
			string(m.RemarkMIA20), string(m.RemarkMIA21), string(m.RemarkMIA22), string(m.RemarkMIA23), amountOut(m.CapitalExceptionMIA24))
	}
	if o.OtherInsurance.Set {
		i := o.OtherInsurance
		c.add(i.Index, "OI", i.FilingIndicatorOI01, i.SubmissionReasonOI02, string(i.BenefitsAssignOI03), i.SignatureSourceOI04,
			i.ProviderAgreementOI05, string(i.ReleaseOfInfoOI06))
	}
	c.amt(o.OtherPayerPaidAMTD)
	c.amt(o.RemainingLiabilityAMTEAF)
	c.amt(o.NonCoveredAMTA8)
	c.amts(o.Amounts)
	c.providerLoop(o.OtherInsuredNM1IL.Name, ProviderPRV{}, o.OtherInsuredNM1IL.References, o.OtherInsuredNM1IL.Raw)
	if p := o.OtherPayerNM1PR; placed(p.Name) {
		c.nm1(p.Name)
		c.dtp(p.CheckOrRemittanceDTP573)
		c.refs([]ReferenceREF{p.SecondaryIDREF2U, p.SecondaryIDREFEI, p.SecondaryIDREFFY, p.SecondaryIDREFNF,
			p.PriorAuthREFG1, p.ReferralNumberREF9F, p.AdjustmentIndicatorREFT4, p.ClaimControlREFF8})
		c.refs(p.References)
		c.raw(p.Raw)
	}
	c.providerLoop(o.OtherPayerAttendingNM171.Name, ProviderPRV{}, o.OtherPayerAttendingNM171.References, o.OtherPayerAttendingNM171.Raw)
	c.providerLoop(o.OtherPayerOperatingNM172.Name, ProviderPRV{}, o.OtherPayerOperatingNM172.References, o.OtherPayerOperatingNM172.Raw)
	c.providerLoop(o.OtherPayerOtherOperatingNM1ZZ.Name, ProviderPRV{}, o.OtherPayerOtherOperatingNM1ZZ.References, o.OtherPayerOtherOperatingNM1ZZ.Raw)
	c.providerLoop(o.OtherPayerServiceFacilityNM177.Name, ProviderPRV{}, o.OtherPayerServiceFacilityNM177.References, o.OtherPayerServiceFacilityNM177.Raw)
	c.providerLoop(o.OtherPayerRenderingNM182.Name, ProviderPRV{}, o.OtherPayerRenderingNM182.References, o.OtherPayerRenderingNM182.Raw)
	c.providerLoop(o.OtherPayerReferringNM1DN.Name, ProviderPRV{}, o.OtherPayerReferringNM1DN.References, o.OtherPayerReferringNM1DN.Raw)
	c.providerLoop(o.OtherPayerBillingNM185.Name, ProviderPRV{}, o.OtherPayerBillingNM185.References, o.OtherPayerBillingNM185.Raw)
	c.raw(o.Raw)
}

func (c *collector) serviceLine(l ServiceLineLX) {
	c.add(l.Index, "LX", l.NumberLX01)
	if l.Service.Set {
		v := l.Service
		c.add(v.Index, "SV2", v.RevenueCodeSV201, c.procedureText(v.ProcedureSV202), amountOut(v.ChargeSV203), string(v.UnitSV204),
			amountOut(v.UnitsSV205), amountOut(v.UnitRateSV206), amountOut(v.DeniedChargeSV207), v.ConditionSV208,
			v.ResidentialStatusSV209, v.LevelOfCareSV210)
	}
	c.pwks(l.SupplementalInfo)
	c.dtp(l.ServiceDateDTP472)
	c.dtps(l.Dates)
	c.refs([]ReferenceREF{l.LineControlREF6R, l.RepricedLineREF9B, l.AdjustedRepricedLineREF9D})
	c.refs(l.References)
	c.amt(l.ServiceTaxAMTGT)
	c.amt(l.FacilityTaxAMTN8)
	c.amts(l.Amounts)
	c.hcp(l.Pricing)
	for _, d := range l.Drugs {
		c.drug(d)
	}
	for _, p := range []LineProviderNM1{l.OperatingNM172, l.OtherOperatingNM1ZZ, l.RenderingNM182, l.ReferringNM1DN} {
		c.providerLoop(p.Name, ProviderPRV{}, p.References, p.Raw)
	}
	for _, a := range l.Adjudications {
		c.adjudication(a)
	}
	c.raw(l.Raw)
}

func (c *collector) drug(d DrugLIN) {
	c.add(d.Index, "LIN", append([]string{d.AssignedIDLIN01, string(d.QualifierLIN02), d.NDCLIN03}, d.ExtraElements...)...)
	if d.Quantity.Set {
		q := d.Quantity
		unit := c.join(append([]string{string(q.UnitCTP05_1)}, q.UnitCTP05Rest...)...)
		c.add(q.Index, "CTP", append([]string{q.ClassOfTradeCTP01, q.PriceIdentifierCTP02, amountOut(q.UnitPriceCTP03),
			amountOut(q.UnitCountCTP04), unit}, q.ExtraElements...)...)
	}
	c.ref(d.LinkSequenceREFVY)
	c.ref(d.PrescriptionREFXZ)
	c.refs(d.References)
	c.raw(d.Raw)
}

func (c *collector) adjudication(a LineAdjudicationSVD) {
	c.add(a.Index, "SVD", a.OtherPayerIDSVD01, amountOut(a.PaymentSVD02), c.procedureText(a.ProcedureSVD03), a.RevenueCodeSVD04,
		amountOut(a.PaidUnitsSVD05), a.BundledLineSVD06)
	c.cas(a.Adjustments)
	c.dtp(a.CheckOrRemittanceDateDTP573)
	c.dtps(a.Dates)
	c.amt(a.RemainingLiabilityAMTEAF)
	c.amts(a.Amounts)
	c.raw(a.Raw)
}
