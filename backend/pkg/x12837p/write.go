package x12837p

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Bytes returns isa as X12 text, the same bytes Write writes.
func Bytes(isa InterchangeISA) ([]byte, error) {
	d := isa.Delimiters
	if len(d.Element) != 1 {
		return nil, fmt.Errorf("x12837p: writing interchange: the element separator is not one byte")
	}
	if len(d.Segment) != 1 {
		return nil, fmt.Errorf("x12837p: writing interchange: the segment terminator is not one byte")
	}
	if len(isa.ComponentSepISA16) != 1 {
		return nil, fmt.Errorf("x12837p: writing interchange: ISA16, the component separator, is not one byte")
	}
	out := encoder{element: d.Element[0], terminator: d.Segment[0], after: d.AfterSegment}
	c := collector{component: isa.ComponentSepISA16[0]}

	out.segment(c.isaSegment(isa))

	var blocks []block
	for _, g := range isa.FunctionalGroups {
		blocks = append(blocks, block{index: g.Index, segments: c.group(g)})
	}
	for _, r := range isa.Raw {
		blocks = append(blocks, block{index: r.Index, segments: []outSegment{{index: r.Index, id: r.ID, elements: r.Elements}}})
	}
	sort.SliceStable(blocks, func(i, j int) bool { return blocks[i].index < blocks[j].index })
	for _, b := range blocks {
		for _, s := range b.segments {
			out.segment(s)
		}
	}
	out.segment(outSegment{id: "IEA", elements: []string{strconv.Itoa(isa.GroupCountIEA01), isa.ControlNumberIEA02}})
	return out.buf.Bytes(), nil
}

// outSegment is one segment on its way out: its position in the file (0
// for a model built in code) and its elements, composites already joined.
type outSegment struct {
	index    int
	id       string
	elements []string
}

// block is a run of segments that moves as one when the writer sorts: a
// functional group or a transaction, keyed by its opening segment's index,
// or one raw segment.
type block struct {
	index    int
	segments []outSegment
}

// encoder turns segments into bytes.
type encoder struct {
	element    byte
	terminator byte
	after      string
	buf        bytes.Buffer
}

func (e *encoder) segment(s outSegment) {
	e.buf.WriteString(s.id)
	for _, el := range s.elements {
		e.buf.WriteByte(e.element)
		e.buf.WriteString(el)
	}
	e.buf.WriteByte(e.terminator)
	e.buf.WriteString(e.after)
}

// collector gathers the segments of one envelope level in traversal order,
// to be sorted by index before they are written.
type collector struct {
	component byte
	segments  []outSegment
}

// add appends a segment with its trailing empty elements left off.
func (c *collector) add(index int, id string, elements ...string) {
	c.segments = append(c.segments, outSegment{index: index, id: id, elements: trimTrailingEmpty(elements)})
}

// raw appends raw segments as they were received.
func (c *collector) raw(raw []RawSegment) {
	for _, r := range raw {
		c.segments = append(c.segments, outSegment{index: r.Index, id: r.ID, elements: r.Elements})
	}
}

// sorted returns the segments in file order: by index, with the traversal
// order kept for equal indexes.
func (c *collector) sorted() []outSegment {
	sort.SliceStable(c.segments, func(i, j int) bool { return c.segments[i].index < c.segments[j].index })
	return c.segments
}

// composite joins components with the component separator, trailing empty
// components left off.
func (c *collector) composite(parts ...string) string {
	return strings.Join(trimTrailingEmpty(parts), string(c.component))
}

// isaSegment builds the ISA with all sixteen elements, padding included,
// since the ISA is fixed width and never drops a trailing element.
func (c *collector) isaSegment(isa InterchangeISA) outSegment {
	return outSegment{id: "ISA", elements: []string{
		string(isa.AuthQualifierISA01), isa.AuthInfoISA02,
		string(isa.SecQualifierISA03), isa.SecInfoISA04,
		string(isa.SenderIDQualifierISA05), isa.SenderIDISA06,
		string(isa.ReceiverIDQualifierISA07), isa.ReceiverIDISA08,
		shortDateOut(isa.DateISA09), timeOut(isa.TimeISA10),
		isa.RepetitionSepISA11, isa.VersionISA12, isa.ControlNumberISA13,
		string(isa.AckRequestedISA14), string(isa.UsageISA15), isa.ComponentSepISA16,
	}}
}

// group returns a functional group's segments: the GS, its transactions
// and raw segments in file order, and the GE.
func (c *collector) group(g FunctionalGroupGS) []outSegment {
	head := collector{component: c.component}
	head.add(g.Index, "GS", append([]string{string(g.FunctionalIDGS01), g.SenderCodeGS02, g.ReceiverCodeGS03, dateOut(g.DateGS04), timeOut(g.TimeGS05), g.ControlNumberGS06, string(g.AgencyGS07), g.VersionGS08}, g.Extra...)...)
	segments := head.segments

	var blocks []block
	for _, t := range g.Transactions {
		blocks = append(blocks, block{index: t.Index, segments: c.transaction(t)})
	}
	for _, r := range g.Raw {
		blocks = append(blocks, block{index: r.Index, segments: []outSegment{{index: r.Index, id: r.ID, elements: r.Elements}}})
	}
	sort.SliceStable(blocks, func(i, j int) bool { return blocks[i].index < blocks[j].index })
	for _, b := range blocks {
		segments = append(segments, b.segments...)
	}
	segments = append(segments, outSegment{id: "GE", elements: trimTrailingEmpty([]string{strconv.Itoa(g.TransactionCountGE01), g.ControlNumberGE02})})
	return segments
}

// transaction returns a transaction set's segments: the ST, everything
// inside it in file order, and the SE.
func (c *collector) transaction(t TransactionST) []outSegment {
	head := collector{component: c.component}
	head.add(t.Index, "ST", append([]string{string(t.TransactionSetIDST01), t.ControlNumberST02, t.ImplementationGuideVersionST03}, t.Extra...)...)

	body := collector{component: c.component}
	body.bht(t.Header)
	body.submitter(t.SubmitterNM141)
	body.receiver(t.ReceiverNM140)
	body.raw(t.Raw)
	for _, b := range t.BillingProviders {
		body.billingProvider(b)
	}

	segments := append(head.segments, body.sorted()...)
	return append(segments, outSegment{id: "SE", elements: trimTrailingEmpty([]string{strconv.Itoa(t.SegmentCountSE01), t.ControlNumberSE02})})
}

// ---------------------------------------------------------------- header loops

func (c *collector) bht(h HeaderBHT) {
	if !h.Set {
		return
	}
	c.add(h.Index, "BHT", append([]string{string(h.StructureBHT01), string(h.PurposeBHT02), h.OriginatorTransactionIDBHT03, dateOut(h.DateBHT04), timeOut(h.TimeBHT05), string(h.TransactionTypeBHT06)}, h.Extra...)...)
}

func (c *collector) submitter(s SubmitterNM141) {
	if !c.nm1(s.Name) {
		return
	}
	for _, p := range s.Contacts {
		c.perAlways(p)
	}
	c.raw(s.Raw)
}

func (c *collector) receiver(r ReceiverNM140) {
	if !c.nm1(r.Name) {
		return
	}
	c.raw(r.Raw)
}

// ---------------------------------------------------------------- hierarchy

func (c *collector) billingProvider(b BillingProviderHL) {
	c.hl(b.HL)
	c.prv(b.Specialty)
	c.cur(b.Currency)
	if c.nm1(b.ProviderNM185.NM1) {
		p := b.ProviderNM185
		c.n3(p.Address)
		c.n4(p.Location)
		c.refSlot(p.TaxIDREFEI)
		c.refs(p.References)
		c.per(p.Contact)
		c.raw(p.Raw)
	}
	if c.nm1(b.PayToProviderNM187.NM1) {
		p := b.PayToProviderNM187
		c.n3(p.Address)
		c.n4(p.Location)
		c.raw(p.Raw)
	}
	if c.nm1(b.PayToPlanNM1PE.NM1) {
		p := b.PayToPlanNM1PE
		c.n3(p.Address)
		c.n4(p.Location)
		c.refs(p.References)
		c.raw(p.Raw)
	}
	c.raw(b.Raw)
	for _, s := range b.Subscribers {
		c.subscriber(s)
	}
}

func (c *collector) subscriber(s SubscriberHL) {
	c.hl(s.HL)
	c.sbr(s.Info)
	c.pat(s.PatientInfo)
	if c.nm1(s.SubscriberNM1IL.NM1) {
		n := s.SubscriberNM1IL
		c.n3(n.Address)
		c.n4(n.Location)
		c.dmg(n.Demographics)
		c.refs(n.References)
		c.per(n.Contact)
		c.raw(n.Raw)
	}
	if c.nm1(s.PayerNM1PR.NM1) {
		n := s.PayerNM1PR
		c.n3(n.Address)
		c.n4(n.Location)
		c.refs(n.References)
		c.raw(n.Raw)
	}
	c.raw(s.Raw)
	for _, cl := range s.Claims {
		c.claim(cl)
	}
	for _, p := range s.Patients {
		c.patient(p)
	}
}

func (c *collector) patient(p PatientHL) {
	c.hl(p.HL)
	c.pat(p.PatientInfo)
	if c.nm1(p.PatientNM1QC.NM1) {
		n := p.PatientNM1QC
		c.n3(n.Address)
		c.n4(n.Location)
		c.dmg(n.Demographics)
		c.refs(n.References)
		c.per(n.Contact)
		c.raw(n.Raw)
	}
	c.raw(p.Raw)
	for _, cl := range p.Claims {
		c.claim(cl)
	}
}

// ---------------------------------------------------------------- claim

func (c *collector) claim(cl ClaimCLM) {
	location := c.composite(append([]string{cl.PlaceOfServiceCLM05_1, string(cl.FacilityQualifierCLM05_2), string(cl.FrequencyCLM05_3)}, cl.LocationExtraCLM05...)...)
	causes := c.composite(append([]string{string(cl.RelatedCausesCLM11[0]), string(cl.RelatedCausesCLM11[1]), cl.RelatedCausesCLM11_3, cl.AutoAccidentStateCLM11_4, cl.CountryCLM11_5}, cl.RelatedCausesExtraCLM11...)...)
	c.add(cl.Index, "CLM", append([]string{
		cl.PatientControlCLM01, amountOut(cl.TotalChargeCLM02), cl.ClaimFilingIndicatorCLM03, cl.ClaimTypeCLM04, location,
		string(cl.SignatureCLM06), string(cl.AssignmentCLM07), string(cl.BenefitsAssignmentCLM08), string(cl.ReleaseOfInfoCLM09), string(cl.SignatureSourceCLM10),
		causes, cl.SpecialProgramCLM12, cl.YesNoCLM13, cl.LevelOfServiceCLM14, cl.YesNoCLM15, cl.ParticipationCLM16, cl.ClaimStatusCLM17, cl.YesNoCLM18, cl.SubmissionReasonCLM19,
		string(cl.DelayReasonCLM20),
	}, cl.Extra...)...)

	for _, d := range []DateDTP{
		cl.OnsetDTP431, cl.InitialTreatmentDTP454, cl.LastSeenDTP304, cl.AcuteManifestationDTP453, cl.AccidentDTP439,
		cl.LastMenstrualPeriodDTP484, cl.LastXRayDTP455, cl.HearingVisionPrescriptionDTP471, cl.DisabilityDTP314, cl.LastWorkedDTP297,
		cl.AuthorizedReturnToWorkDTP296, cl.AdmissionDTP435, cl.DischargeDTP096, cl.AssumedOrRelinquishedCareDTP090[0],
		cl.AssumedOrRelinquishedCareDTP090[1], cl.PropertyCasualtyFirstContactDTP444, cl.RepricerReceivedDTP050,
	} {
		c.dtpSlot(d)
	}
	for _, p := range cl.SupplementalInfoPWK {
		c.pwk(p)
	}
	c.cn1(cl.ContractInfoCN1)
	c.amtSlot(cl.PatientAmountPaidAMTF5)
	for _, r := range []ReferenceREF{
		cl.ServiceAuthExceptionREF4N, cl.MedicareCrossoverIndicatorREFF5, cl.MammographyCertificationREFEW, cl.ReferralREF9F,
		cl.PriorAuthorizationREFG1, cl.PayerClaimControlNumberREFF8, cl.CLIAREFX4, cl.RepricedClaimREF9A, cl.AdjustedRepricedClaimREF9C,
		cl.InvestigationalDeviceExemptionREFLX, cl.TransmissionIntermediaryREFD9, cl.MedicalRecordNumberREFEA,
		cl.DemonstrationProjectREFP4, cl.CarePlanOversightREF1J,
	} {
		c.refSlot(r)
	}
	c.refs(cl.References)
	for _, k := range cl.FileInfoK3 {
		c.k3(k)
	}
	c.nte(cl.NoteNTE)
	c.cr1(cl.AmbulanceTransportCR1)
	c.cr2(cl.SpinalManipulationCR2)
	for _, crc := range cl.AmbulanceCertificationCRC07 {
		c.crc(crc)
	}
	for _, crc := range cl.VisionConditionCRCE1 {
		c.crc(crc)
	}
	c.crc(cl.HomeboundCRC75)
	c.crc(cl.EPSDTReferralCRCZZ)
	c.hi(cl.DiagnosisCodesHIABK)
	c.hi(cl.AnesthesiaProceduresHIBP)
	for _, h := range cl.ConditionCodesHIBG {
		c.hi(h)
	}
	c.hcp(cl.PricingHCP)
	c.raw(cl.Raw)

	if c.nm1(cl.ReferringProvider.Name) {
		c.refs(cl.ReferringProvider.References)
		c.raw(cl.ReferringProvider.Raw)
	}
	if c.nm1(cl.RenderingProviderNM182.Name) {
		c.prv(cl.RenderingProviderNM182.Specialty)
		c.refs(cl.RenderingProviderNM182.References)
		c.raw(cl.RenderingProviderNM182.Raw)
	}
	if c.nm1(cl.ServiceFacilityNM177.Name) {
		f := cl.ServiceFacilityNM177
		c.n3(f.Address)
		c.n4(f.Location)
		c.refs(f.References)
		c.per(f.Contact)
		c.raw(f.Raw)
	}
	if c.nm1(cl.SupervisingProviderNM1DQ.Name) {
		c.refs(cl.SupervisingProviderNM1DQ.References)
		c.raw(cl.SupervisingProviderNM1DQ.Raw)
	}
	if c.nm1(cl.AmbulancePickupNM1PW.Name) {
		c.n3(cl.AmbulancePickupNM1PW.Address)
		c.n4(cl.AmbulancePickupNM1PW.Location)
		c.raw(cl.AmbulancePickupNM1PW.Raw)
	}
	if c.nm1(cl.AmbulanceDropoffNM145.Name) {
		c.n3(cl.AmbulanceDropoffNM145.Address)
		c.n4(cl.AmbulanceDropoffNM145.Location)
		c.raw(cl.AmbulanceDropoffNM145.Raw)
	}
	for _, o := range cl.OtherSubscribers {
		c.otherSubscriber(o)
	}
	for _, l := range cl.ServiceLines {
		c.serviceLine(l)
	}
}

func (c *collector) otherSubscriber(o OtherSubscriberSBR) {
	c.sbr(o.Subscriber)
	for _, a := range o.Adjustments {
		c.cas(a)
	}
	c.amtSlot(o.PayerPaidAMTD)
	c.amtSlot(o.NonCoveredAMTA8)
	c.amtSlot(o.RemainingLiabilityAMTEAF)
	c.oi(o.OtherInsurance)
	c.moa(o.OutpatientAdjudication)
	c.raw(o.Raw)
	if c.nm1(o.OtherSubscriberNM1IL.Name) {
		s := o.OtherSubscriberNM1IL
		c.n3(s.Address)
		c.n4(s.Location)
		c.refSlot(s.SecondaryID)
		c.raw(s.Raw)
	}
	if c.nm1(o.OtherPayerNM1PR.Name) {
		p := o.OtherPayerNM1PR
		c.n3(p.Address)
		c.n4(p.Location)
		c.dtpSlot(p.RemittanceDateDTP573)
		c.refs(p.SecondaryID)
		c.refSlot(p.PriorAuthorizationREFG1)
		c.refSlot(p.ReferralREF9F)
		c.refSlot(p.ClaimAdjustmentIndicatorREFT4)
		c.refSlot(p.ClaimControlNumberREFF8)
		c.raw(p.Raw)
	}
	for _, r := range []OtherPayerRoleNM1{o.OtherPayerReferringProvider, o.OtherPayerRenderingProviderNM182,
		o.OtherPayerServiceFacilityNM177, o.OtherPayerSupervisingProviderNM1DQ, o.OtherPayerBillingProviderNM185} {
		if c.nm1(r.Name) {
			c.refs(r.References)
			c.raw(r.Raw)
		}
	}
}

// ---------------------------------------------------------------- service line

func (c *collector) serviceLine(l ServiceLineLX) {
	c.add(l.Index, "LX", append([]string{l.AssignedNumberLX01}, l.Extra...)...)
	c.sv1(l.Service)
	c.sv5(l.DME)
	for _, p := range l.Attachments {
		c.pwk(p)
	}
	c.cr1(l.AmbulanceTransport)
	c.cr3(l.DMECertification)
	for _, crc := range l.AmbulanceCertificationCRC07 {
		c.crc(crc)
	}
	c.crc(l.HospiceEmployeeCRC70)
	c.crc(l.DMEConditionCRC09)
	for _, d := range []DateDTP{l.ServiceDateDTP472, l.PrescriptionDTP471, l.CertificationRevisionDTP607, l.BeginTherapyDTP463,
		l.LastCertificationDTP461, l.LastTreatmentOrTherapyDTP304} {
		c.dtpSlot(d)
	}
	for _, d := range l.TestPerformedDTP738 {
		c.dtp(d)
	}
	for _, d := range []DateDTP{l.ShippedDTP011, l.LastXRayDTP455, l.InitialTreatmentDTP454} {
		c.dtpSlot(d)
	}
	c.qtySlot(l.AmbulancePatientCountQTYPT)
	c.qtySlot(l.ObstetricAdditionalUnitsQTYFL)
	for _, m := range l.TestResults {
		c.add(m.Index, "MEA", append([]string{m.ReferenceIDMEA01, m.QualifierMEA02, amountOut(m.ResultMEA03)}, m.Extra...)...)
	}
	c.cn1(l.Contract)
	c.refSlot(l.RepricedLineItemREF9B)
	c.refSlot(l.AdjustedRepricedLineItemREF9D)
	c.refs(l.PriorAuthorizationREFG1)
	c.refSlot(l.LineControlREF6R)
	c.refSlot(l.MammographyCertificationREFEW)
	c.refSlot(l.CLIAREFX4)
	c.refSlot(l.ReferringCLIAREFF4)
	c.refSlot(l.ImmunizationBatchREFBT)
	c.refs(l.ReferralREF9F)
	c.amtSlot(l.SalesTaxAMTT)
	c.amtSlot(l.PostageClaimedAMTF4)
	for _, k := range l.FileInfoK3 {
		c.k3(k)
	}
	c.nte(l.LineNote)
	c.nte(l.ThirdPartyNote)
	if l.PurchasedService.Set {
		p := l.PurchasedService
		c.add(p.Index, "PS1", append([]string{p.ProviderIDPS101, amountOut(p.ChargeAmountPS102)}, p.Extra...)...)
	}
	c.hcp(l.Pricing)
	c.raw(l.Raw)

	if l.Drug.Set {
		d := l.Drug
		c.add(d.Index, "LIN", append([]string{d.AssignedIDLIN01, d.QualifierLIN02, d.NDCLIN03}, d.Extra...)...)
		if d.Pricing.Set {
			p := d.Pricing
			unit := c.composite(append([]string{string(p.UnitCTP05_1)}, p.UnitCTP05Extra...)...)
			c.add(p.Index, "CTP", append([]string{p.ClassOfTradeCTP01, p.PriceIdentifierCTP02, amountOut(p.UnitPriceCTP03), amountOut(p.UnitCountCTP04), unit}, p.Extra...)...)
		}
		c.refSlot(d.AssociationREFXZ)
		c.raw(d.Raw)
	}
	for _, p := range []LineProviderNM1{l.RenderingNM182, l.PurchasedServiceNM1QB, l.ServiceFacilityNM177,
		l.SupervisingNM1DQ, l.OrderingNM1DK, l.ReferringNM12420F, l.AmbulancePickupNM1PW, l.AmbulanceDropoffNM145} {
		if !p.Set || !c.nm1(p.Name) {
			continue
		}
		c.prv(p.Specialty)
		c.n3(p.Address)
		c.n4(p.Location)
		c.per(p.Contact)
		c.refs(p.References)
		c.raw(p.Raw)
	}
	for _, a := range l.Adjudications {
		c.add(a.Index, "SVD", append([]string{a.OtherPayerSVD01, amountOut(a.PaidAmountSVD02), c.procedure(a.ProcedureSVD03), a.ProductServiceIDSVD04, amountOut(a.PaidUnitsSVD05), a.BundledLineSVD06}, a.Extra...)...)
		for _, cas := range a.Adjustments {
			c.cas(cas)
		}
		c.dtpSlot(a.AdjudicationDateDTP573)
		c.amtSlot(a.RemainingLiabilityAMTEAF)
		c.raw(a.Raw)
	}
	for _, f := range l.Forms {
		c.add(f.Index, "LQ", append([]string{f.QualifierLQ01, f.FormIDLQ02}, f.Extra...)...)
		for _, a := range f.Answers {
			c.add(a.Index, "FRM", append([]string{a.QuestionFRM01, a.ResponseFRM02, a.ResponseFRM03, dateOut(a.ResponseFRM04), amountOut(a.ResponseFRM05)}, a.Extra...)...)
		}
		c.raw(f.Raw)
	}
}

func (c *collector) sv1(v ProfessionalServiceSV1) {
	if !v.Set {
		return
	}
	pointers := make([]string, 0, 4+len(v.PointersExtraSV107))
	for _, p := range v.DiagnosisPointersSV107 {
		pointers = append(pointers, pointerOut(p))
	}
	pointers = append(pointers, v.PointersExtraSV107...)
	c.add(v.Index, "SV1", append([]string{
		c.procedure(v.ProcedureSV101), amountOut(v.ChargeSV102), string(v.UnitSV103), amountOut(v.UnitCountSV104), v.PlaceOfServiceSV105,
		v.ServiceTypeSV106, c.composite(pointers...), amountOut(v.MonetaryAmountSV108), string(v.EmergencySV109), v.MultipleProcedureSV110,
		string(v.EPSDTSV111), string(v.FamilyPlanningSV112), v.ReviewCodeSV113, v.ReviewValueSV114, v.CoPayStatusSV115,
	}, v.Extra...)...)
}

func (c *collector) sv5(d LineDMEServiceSV5) {
	if !d.Set {
		return
	}
	c.add(d.Index, "SV5", append([]string{c.procedure(d.ProcedureSV501), string(d.UnitSV502), amountOut(d.LengthOfNecessitySV503), amountOut(d.RentalPriceSV504), amountOut(d.PurchasePriceSV505), d.RentalUnitPriceIndicatorSV506}, d.Extra...)...)
}

func (c *collector) cr3(r LineDMECertificationCR3) {
	if !r.Set {
		return
	}
	c.add(r.Index, "CR3", append([]string{r.CertificationTypeCR301, string(r.UnitCR302), amountOut(r.DurationCR303)}, r.Extra...)...)
}

func (c *collector) procedure(p ProcedureSV101) string {
	return c.composite(append([]string{p.Qualifier, p.Code, p.Modifiers[0], p.Modifiers[1], p.Modifiers[2], p.Modifiers[3], p.Description}, p.Extra...)...)
}

// ---------------------------------------------------------------- shared segments

// nm1 writes an NM1 and reports whether the loop it opens was present: a
// placed NM1 always carries NM101.
func (c *collector) nm1(n NameNM1) bool {
	if n.EntityIDNM101 == EntityIdentifierNM101NotSent {
		return false
	}
	c.add(n.Index, "NM1", append([]string{string(n.EntityIDNM101), string(n.EntityTypeNM102), n.LastNM103, n.FirstNM104, n.MiddleNM105, n.PrefixNM106, n.SuffixNM107, string(n.IDQualifierNM108), n.IDNM109}, n.Extra...)...)
	return true
}

func (c *collector) hl(h HierarchyHL) {
	c.add(h.Index, "HL", append([]string{h.IDNumberHL01, h.ParentIDNumberHL02, string(h.LevelHL03), string(h.ChildCodeHL04)}, h.Extra...)...)
}

func (c *collector) n3(a AddressN3) {
	if !a.Set {
		return
	}
	c.add(a.Index, "N3", append([]string{a.Line1N301, a.Line2N302}, a.Extra...)...)
}

func (c *collector) n4(l LocationN4) {
	if !l.Set {
		return
	}
	c.add(l.Index, "N4", append([]string{l.CityN401, l.StateN402, l.ZipN403, l.CountryN404, l.LocationQualifierN405, l.LocationIdentifierN406, l.SubdivisionN407}, l.Extra...)...)
}

func (c *collector) dmg(d DemographicDMG) {
	if !d.Set {
		return
	}
	c.add(d.Index, "DMG", append([]string{d.FormatDMG01, dateOut(d.BirthDateDMG02), d.GenderDMG03}, d.Extra...)...)
}

// per writes a placed PER, one whose PER01 carries a value.
func (c *collector) per(p ContactPER) {
	if p.FunctionPER01 == ContactFunctionCodePER01NotSent {
		return
	}
	c.perAlways(p)
}

func (c *collector) perAlways(p ContactPER) {
	c.add(p.Index, "PER", append([]string{string(p.FunctionPER01), p.NamePER02, string(p.QualifierPER03), p.CommPER04, string(p.QualifierPER05), p.CommPER06, string(p.QualifierPER07), p.CommPER08, p.InquiryReferencePER09}, p.Extra...)...)
}

func (c *collector) prv(p ProviderPRV) {
	if p.CodePRV01 == ProviderTypePRV01NotSent {
		return
	}
	c.add(p.Index, "PRV", append([]string{string(p.CodePRV01), string(p.QualifierPRV02), p.TaxonomyPRV03}, p.Extra...)...)
}

func (c *collector) cur(u CurrencyCUR) {
	if !u.Set {
		return
	}
	c.add(u.Index, "CUR", append([]string{u.EntityIDCUR01, u.CurrencyCodeCUR02}, u.Extra...)...)
}

func (c *collector) sbr(s SubscriberSBR) {
	if !s.Set {
		return
	}
	c.add(s.Index, "SBR", append([]string{string(s.PayerResponsibilitySBR01), string(s.RelationshipSBR02), s.GroupOrPolicyNumberSBR03, s.GroupNameSBR04, string(s.InsuranceTypeSBR05), s.CoordinationSBR06, s.ResponseSBR07, s.EmploymentSBR08, string(s.ClaimFilingIndicatorSBR09)}, s.Extra...)...)
}

func (c *collector) pat(p PatientInfoPAT) {
	if !p.Set {
		return
	}
	c.add(p.Index, "PAT", append([]string{string(p.RelationshipPAT01), p.LocationPAT02, p.EmploymentStatusPAT03, p.StudentStatusPAT04, p.FormatPAT05, dateOut(p.DeathDatePAT06), string(p.WeightUOMPAT07), amountOut(p.WeightPAT08), string(p.PregnancyPAT09)}, p.Extra...)...)
}

// refSlot writes a REF held in a qualifier-specific field, which is
// present when its qualifier is.
func (c *collector) refSlot(r ReferenceREF) {
	if r.QualifierREF01 == ReferenceIdentificationQualifierREF01NotSent {
		return
	}
	c.ref(r)
}

// refs writes every REF of a list; a listed REF was sent whatever it
// carries.
func (c *collector) refs(refs []ReferenceREF) {
	for _, r := range refs {
		c.ref(r)
	}
}

func (c *collector) ref(r ReferenceREF) {
	c.add(r.Index, "REF", append([]string{string(r.QualifierREF01), r.Value}, r.Extra...)...)
}

// dtpSlot writes a DTP held in a qualifier-specific field, which is present
// when its qualifier is.
func (c *collector) dtpSlot(d DateDTP) {
	if d.QualifierDTP01 == DateTimeQualifierDTP01NotSent {
		return
	}
	c.dtp(d)
}

func (c *collector) dtp(d DateDTP) {
	c.add(d.Index, "DTP", append([]string{string(d.QualifierDTP01), d.FormatDTP02, d.RawDTP03}, d.Extra...)...)
}

// amtSlot writes an AMT held in a qualifier-specific field, which is
// present when its qualifier is.
func (c *collector) amtSlot(a AmountAMT) {
	if a.QualifierAMT01 == AmountQualifierAMT01NotSent {
		return
	}
	c.add(a.Index, "AMT", append([]string{string(a.QualifierAMT01), amountOut(a.Amount)}, a.Extra...)...)
}

// qtySlot writes a QTY held in a qualifier-specific field, which is
// present when its qualifier is.
func (c *collector) qtySlot(q QuantityQTY) {
	if q.QualifierQTY01 == QuantityQualifierQTY01NotSent {
		return
	}
	c.add(q.Index, "QTY", append([]string{string(q.QualifierQTY01), amountOut(q.Quantity)}, q.Extra...)...)
}

func (c *collector) pwk(p AttachmentPWK) {
	c.add(p.Index, "PWK", append([]string{string(p.ReportTypePWK01), string(p.TransmissionPWK02), amountOut(p.ReportCopiesPWK03), p.EntityIDPWK04, p.IDQualifierPWK05, p.ControlNumberPWK06}, p.Extra...)...)
}

func (c *collector) cn1(n ContractCN1) {
	if !n.Set {
		return
	}
	c.add(n.Index, "CN1", append([]string{n.TypeCN101, amountOut(n.AmountCN102), amountOut(n.PercentageCN103), n.CodeCN104, amountOut(n.TermsDiscountPercentCN105), n.VersionCN106}, n.Extra...)...)
}

func (c *collector) k3(k FileInfoK3) {
	c.add(k.Index, "K3", append([]string{k.TextK301}, k.Extra...)...)
}

func (c *collector) nte(n NoteNTE) {
	if !n.Set {
		return
	}
	c.add(n.Index, "NTE", append([]string{n.ReferenceCodeNTE01, n.TextNTE02}, n.Extra...)...)
}

func (c *collector) cr1(t AmbulanceTransportCR1) {
	if !t.Set {
		return
	}
	c.add(t.Index, "CR1", append([]string{string(t.WeightUnitCR101), amountOut(t.WeightCR102), t.TransportCodeCR103, t.ReasonCR104, string(t.DistanceUnitCR105), amountOut(t.DistanceCR106), t.AddressInfoCR107, t.AddressInfoCR108, t.RoundTripPurposeCR109, t.StretcherPurposeCR110}, t.Extra...)...)
}

func (c *collector) cr2(s ClaimSpinalManipulationCR2) {
	if !s.Set {
		return
	}
	c.add(s.Index, "CR2", append([]string{amountOut(s.TreatmentSeriesCR201), amountOut(s.TreatmentCountCR202), s.SubluxationLevelCR203, s.SubluxationLevelCR204, s.UnitCR205, amountOut(s.TreatmentPeriodCountCR206), amountOut(s.MonthlyTreatmentCountCR207), s.ConditionCodeCR208, s.ComplicationIndicatorCR209, s.ConditionDescriptionCR210, s.ConditionDescriptionCR211}, s.Extra...)...)
}

func (c *collector) crc(r ConditionCRC) {
	if !r.Set {
		return
	}
	c.add(r.Index, "CRC", append([]string{string(r.CodeCategoryCRC01), string(r.CertConditionCRC02), string(r.ConditionsCRC0x[0]), string(r.ConditionsCRC0x[1]), string(r.ConditionsCRC0x[2]), string(r.ConditionsCRC0x[3]), string(r.ConditionsCRC0x[4])}, r.Extra...)...)
}

func (c *collector) hi(h HealthCareInfoHI) {
	if !h.Set {
		return
	}
	elements := make([]string, 0, len(h.Codes)+len(h.Extra))
	for _, code := range h.Codes {
		elements = append(elements, c.composite(append([]string{
			string(code.TypeHI0x1), code.CodeHI0x2, code.DateQualifierHI0x3, periodOut(code),
			amountOut(code.AmountHI0x5), amountOut(code.QuantityHI0x6), code.VersionIdentifierHI0x7, code.IndustryCodeHI0x8, string(code.YesNoHI0x9),
		}, code.Extra...)...))
	}
	c.add(h.Index, "HI", append(elements, h.Extra...)...)
}

func (c *collector) hcp(p PricingHCP) {
	if !p.Set {
		return
	}
	c.add(p.Index, "HCP", append([]string{
		string(p.MethodologyHCP01), amountOut(p.AllowedAmountHCP02), amountOut(p.SavingAmountHCP03), p.OrgIDHCP04, amountOut(p.PerDiemOrFlatRateHCP05),
		p.APGCodeHCP06, amountOut(p.APGAmountHCP07), p.ProductServiceIDHCP08, p.IDQualifierHCP09, p.ApprovedHCPCSHCP10,
		string(p.UnitHCP11), amountOut(p.ApprovedUnitsHCP12), p.RejectReasonHCP13, p.PolicyComplianceHCP14, p.ExceptionHCP15,
	}, p.Extra...)...)
}

func (c *collector) cas(a AdjustmentCAS) {
	elements := make([]string, 0, 1+3*len(a.Reasons)+len(a.Extra))
	elements = append(elements, string(a.GroupCAS01))
	for _, r := range a.Reasons {
		elements = append(elements, string(r.Reason), amountOut(r.Amount), amountOut(r.Quantity))
	}
	c.add(a.Index, "CAS", append(elements, a.Extra...)...)
}

func (c *collector) oi(o OtherInsuranceOI) {
	if !o.Set {
		return
	}
	c.add(o.Index, "OI", append([]string{o.ClaimFilingIndicatorOI01, o.ClaimSubmitReasonOI02, string(o.BenefitsAssignmentOI03), string(o.SignatureSourceOI04), o.ProviderAgreementOI05, string(o.ReleaseOfInfoOI06)}, o.Extra...)...)
}

func (c *collector) moa(m OutpatientAdjudicationMOA) {
	if !m.Set {
		return
	}
	c.add(m.Index, "MOA", append([]string{amountOut(m.ReimbursementRateMOA01), amountOut(m.HCPCSPayableAmountMOA02), string(m.RemarkMOA03), string(m.RemarkMOA04), string(m.RemarkMOA05), string(m.RemarkMOA06), string(m.RemarkMOA07), amountOut(m.ESRDPaymentMOA08), amountOut(m.NonPayableProfCompMOA09)}, m.Extra...)...)
}

// ---------------------------------------------------------------- element text

// amountText returns an amount as the file sent it, or, for an amount built
// in code, in its canonical spelling; "" when not sent.
func amountOut(a Amount) string {
	if !a.Set {
		return ""
	}
	if a.Text != "" {
		return a.Text
	}
	return a.Value.String()
}

// dateText returns a date as CCYYMMDD, or "" when not sent.
func dateOut(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%04d%02d%02d", d.Year, d.Month, d.Day)
}

// shortDateText returns a date as YYMMDD, the ISA09 form, or "" when not
// sent.
func shortDateOut(d Date) string {
	if !d.Set {
		return ""
	}
	return fmt.Sprintf("%02d%02d%02d", d.Year%100, d.Month, d.Day)
}

// timeText returns a time as HHMM, or HHMMSS when the file sent seconds;
// "" when not sent.
func timeOut(t Time) string {
	if !t.Set {
		return ""
	}
	if t.HasSeconds {
		return fmt.Sprintf("%02d%02d%02d", t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d%02d", t.Hour, t.Minute)
}

// periodTextOf returns an HI composite's date component as its format
// qualifier says: one CCYYMMDD date, or a CCYYMMDD-CCYYMMDD range; "" when
// not sent.
func periodOut(code HealthCareCodeHI) string {
	if !code.DateHI0x4.Set {
		return ""
	}
	if code.DateThroughHI0x4.Set {
		return dateOut(code.DateHI0x4) + "-" + dateOut(code.DateThroughHI0x4)
	}
	return dateOut(code.DateHI0x4)
}

// pointerText returns a diagnosis pointer as digits, or "" for a position
// not sent.
func pointerOut(p int) string {
	if p == 0 {
		return ""
	}
	return strconv.Itoa(p)
}
