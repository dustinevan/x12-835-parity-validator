package x12837p

import (
	"errors"
	"fmt"
	"strings"
)

// inputSegment is one segment of the input, split on the element separator.
type inputSegment struct {
	index    int      // position in the file, 0-based; the ISA is 0
	id       string   // the segment ID, e.g. "CLM"
	elements []string // elements[0] is element 01
}

// el returns element n (1-based), or "" when the segment ends before it.
func (s inputSegment) el(n int) string {
	if n < 1 || n > len(s.elements) {
		return ""
	}
	return s.elements[n-1]
}

// name returns element n's reference, e.g. "CLM02".
func (s inputSegment) name(n int) string {
	return fmt.Sprintf("%s%02d", s.id, n)
}

// extra returns a copy of the elements after element last, the highest one
// the model names for this segment, so text past the model is kept rather
// than dropped. splitSegments removes trailing empty elements, so the last
// one returned carries text.
func (s inputSegment) extra(last int) []string {
	if len(s.elements) <= last {
		return nil
	}
	return append([]string(nil), s.elements[last:]...)
}

// components splits element n on the component separator and returns the
// components the model names (the first last of them) and, separately, any
// after those, so text past the model's composite is kept.
func (s inputSegment) components(n, last int, component byte) ([]string, []string) {
	parts := trimTrailingEmpty(strings.Split(s.el(n), string(component)))
	if len(parts) <= last {
		return parts, nil
	}
	return parts[:last], append([]string(nil), parts[last:]...)
}

func rawSegment(s inputSegment) RawSegment {
	return RawSegment{Index: s.index, ID: s.id, Elements: s.elements}
}

func rawSegments(loop []inputSegment) []RawSegment {
	raw := make([]RawSegment, 0, len(loop))
	for _, s := range loop {
		raw = append(raw, rawSegment(s))
	}
	return raw
}

// splitSegments reads the delimiters from the ISA and splits data into
// segments. It returns the file's Delimiters and the component separator
// (ISA16) for the composite elements. Whitespace before a segment is
// skipped, a blank segment is skipped, trailing empty elements are removed
// (they carry no text), and when the terminator is a newline a carriage
// return before it is removed. The whitespace after each terminator is
// compared across the file: when every segment, the last included, is
// followed by the same text, that text is the Delimiters' AfterSegment.
func splitSegments(data []byte) ([]inputSegment, Delimiters, byte, error) {
	text := string(data)
	start := 0
	for start < len(text) && isSpace(text[start]) {
		start++
	}
	if !strings.HasPrefix(text[start:], "ISA") || len(text) < start+4 {
		return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting input: expected ISA")
	}
	element := text[start+3]

	// ISA16 is the one byte after the 16th element separator, and the
	// segment terminator is the byte after ISA16.
	isa16 := -1
	separators := 0
	for i := start + 3; i < len(text); i++ {
		if text[i] != element {
			continue
		}
		separators++
		if separators == 16 {
			isa16 = i + 1
			break
		}
	}
	if isa16 < 0 || isa16+1 >= len(text) {
		return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting ISA: incomplete header")
	}
	component := text[isa16]
	terminator := text[isa16+1]
	isaFields := strings.Split(text[start:isa16+1], string(element))
	repetition := isaFields[11]
	if len(repetition) != 1 {
		return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting ISA: ISA11 is not one character")
	}
	if element == component || element == terminator || element == repetition[0] ||
		component == terminator || component == repetition[0] || terminator == repetition[0] {
		return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting ISA: the four separators must differ")
	}
	delimiters := Delimiters{Element: string(element), Segment: string(terminator)}

	segments := []inputSegment{{index: 0, id: "ISA", elements: isaFields[1:]}}
	pos := isa16 + 2
	after := ""       // the whitespace after the first terminator
	sameAfter := true // whether every terminator is followed by that same whitespace
	for {
		gapStart := pos
		for pos < len(text) && isSpace(text[pos]) {
			pos++
		}
		gap := text[gapStart:pos]
		if len(segments) == 1 {
			after = gap
		} else if gap != after {
			sameAfter = false
		}
		if pos == len(text) {
			if sameAfter {
				delimiters.AfterSegment = after
			}
			return segments, delimiters, component, nil
		}
		if text[pos] == terminator {
			pos++
			continue
		}
		end := strings.IndexByte(text[pos:], terminator)
		if end < 0 {
			return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting segment %d: missing segment terminator", len(segments))
		}
		content := text[pos : pos+end]
		pos += end + 1
		if terminator == '\n' {
			content = strings.TrimSuffix(content, "\r")
		}
		fields := strings.Split(content, string(element))
		if !validID(fields[0]) {
			return nil, Delimiters{}, 0, fmt.Errorf("x12837p: splitting segment %d: segment ID is not 2 or 3 capital letters or digits", len(segments))
		}
		segments = append(segments, inputSegment{index: len(segments), id: fields[0], elements: trimTrailingEmpty(fields[1:])})
	}
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}

func validID(id string) bool {
	if len(id) < 2 || len(id) > 3 {
		return false
	}
	for i := 0; i < len(id); i++ {
		if (id[i] < 'A' || id[i] > 'Z') && (id[i] < '0' || id[i] > '9') {
			return false
		}
	}
	return true
}

// parser walks the segments once, building the tree. Each parse method
// consumes one loop; each loop struct's add method places one segment
// inside it.
type parser struct {
	segments  []inputSegment
	i         int  // the next segment to read
	component byte // ISA16

	// transactionLengths holds, in file order, the number of segments from
	// each ST through its SE, for the SE01 check.
	transactionLengths []int
}

// peek returns the next segment without consuming it. At the end of the
// input it returns a segment whose id is "".
func (p *parser) peek() inputSegment {
	if p.i >= len(p.segments) {
		return inputSegment{index: len(p.segments)}
	}
	return p.segments[p.i]
}

func (p *parser) take() inputSegment {
	s := p.segments[p.i]
	p.i++
	return s
}

// loopMembers lists the segment IDs a loop takes after the segment that
// starts it.
type loopMembers []string

func (m loopMembers) has(id string) bool {
	for _, member := range m {
		if member == id {
			return true
		}
	}
	return false
}

// The segments each name loop takes after its NM1, from the boxes
// Transmittal 907 lists for that loop. anyNameMembers is the set used for an
// NM1 loop that has no slot and is kept raw.
var (
	submitterMembers        = loopMembers{"PER"}
	receiverMembers         = loopMembers{}
	billingProviderMembers  = loopMembers{"N3", "N4", "REF", "PER"}
	payToProviderMembers    = loopMembers{"N3", "N4"}
	payToPlanMembers        = loopMembers{"N3", "N4", "REF"}
	subscriberMembers       = loopMembers{"N3", "N4", "DMG", "REF", "PER"}
	payerMembers            = loopMembers{"N3", "N4", "REF"}
	patientMembers          = loopMembers{"N3", "N4", "DMG", "REF", "PER"}
	referenceMembers        = loopMembers{"REF"}
	renderingMembers        = loopMembers{"PRV", "REF"}
	serviceFacilityMembers  = loopMembers{"N3", "N4", "REF", "PER"}
	locationMembers         = loopMembers{"N3", "N4"}
	otherSubscriberMembers  = loopMembers{"N3", "N4", "REF"}
	otherPayerMembers       = loopMembers{"N3", "N4", "DTP", "REF"}
	lineFacilityMembers     = loopMembers{"N3", "N4", "REF"}
	orderingProviderMembers = loopMembers{"N3", "N4", "REF", "PER"}
	anyNameMembers          = loopMembers{"N3", "N4", "REF", "PER", "PRV", "DMG", "DTP"}
	drugMembers             = loopMembers{"CTP", "REF"}
	lineAdjudicationMembers = loopMembers{"CAS", "DTP", "AMT"}
	supportingFormMembers   = loopMembers{"FRM"}
)

// takeLoop takes the segment that starts a loop and every following segment
// whose ID is one of members.
func (p *parser) takeLoop(members loopMembers) []inputSegment {
	loop := []inputSegment{p.take()}
	for members.has(p.peek().id) {
		loop = append(loop, p.take())
	}
	return loop
}

// takeUntilHierarchy takes segments up to the next HL or envelope segment,
// for a loop that has no slot.
func (p *parser) takeUntilHierarchy() []RawSegment {
	var raw []RawSegment
	for !endsHierarchy(p.peek().id) {
		raw = append(raw, rawSegment(p.take()))
	}
	return raw
}

// envelope reports whether a segment with this ID opens or closes an
// envelope. Every envelope segment ends every loop inside ST so that the
// envelope level can close the loop or report the segment out of place.
func envelope(id string) bool {
	switch id {
	case "ISA", "IEA", "GS", "GE", "ST", "SE":
		return true
	}
	return false
}

// endsHierarchy reports whether a segment with this ID ends an HL loop and
// everything in it; "" is the end of the input.
func endsHierarchy(id string) bool {
	return id == "" || id == "HL" || envelope(id)
}

func endsClaim(id string) bool {
	return id == "CLM" || endsHierarchy(id)
}

func endsOtherSubscriber(id string) bool {
	return id == "SBR" || id == "LX" || endsClaim(id)
}

func endsServiceLine(id string) bool {
	return id == "LX" || endsClaim(id)
}

func (p *parser) parseInterchange() (InterchangeISA, error) {
	isa, err := parseISA(p.take()) // splitSegments always returns the ISA first
	if err != nil {
		return InterchangeISA{}, err
	}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return InterchangeISA{}, fmt.Errorf("x12837p: parsing interchange: missing IEA")
		case "GS":
			group, err := p.parseGroup()
			if err != nil {
				return InterchangeISA{}, err
			}
			isa.FunctionalGroups = append(isa.FunctionalGroups, group)
		case "IEA":
			count, control, err := parseTrailer(p.take())
			if err != nil {
				return InterchangeISA{}, err
			}
			isa.GroupCountIEA01 = count
			isa.ControlNumberIEA02 = control
			if p.i < len(p.segments) {
				return InterchangeISA{}, fmt.Errorf("x12837p: parsing segment %d: input continues after the IEA", p.peek().index)
			}
			return isa, nil
		case "ISA", "GE", "ST", "SE":
			return InterchangeISA{}, fmt.Errorf("x12837p: parsing segment %d: %s outside a functional group", next.index, next.id)
		default:
			isa.Raw = append(isa.Raw, rawSegment(p.take()))
		}
	}
}

func (p *parser) parseGroup() (FunctionalGroupGS, error) {
	gs := p.take()
	group, err := parseGS(gs)
	if err != nil {
		return FunctionalGroupGS{}, err
	}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return FunctionalGroupGS{}, fmt.Errorf("x12837p: parsing functional group at segment %d: missing GE", gs.index)
		case "ST":
			transaction, err := p.parseTransaction()
			if err != nil {
				return FunctionalGroupGS{}, err
			}
			group.Transactions = append(group.Transactions, transaction)
		case "GE":
			count, control, err := parseTrailer(p.take())
			if err != nil {
				return FunctionalGroupGS{}, err
			}
			group.TransactionCountGE01 = count
			group.ControlNumberGE02 = control
			return group, nil
		case "ISA", "IEA", "GS", "SE":
			return FunctionalGroupGS{}, fmt.Errorf("x12837p: parsing segment %d: %s inside the functional group at segment %d (missing GE)", next.index, next.id, gs.index)
		default:
			group.Raw = append(group.Raw, rawSegment(p.take()))
		}
	}
}

// hierarchy maps each placed HL's HL01 to where its loop sits in the
// transaction, so a later HL can find its parent by HL02.
type hierarchy struct {
	billingProviders map[string]int    // HL01 → index in TransactionST.BillingProviders
	subscribers      map[string][2]int // HL01 → billing provider index, subscriber index
}

func (p *parser) parseTransaction() (TransactionST, error) {
	st := p.take()
	t, err := parseST(st)
	if err != nil {
		return TransactionST{}, err
	}
	tree := hierarchy{billingProviders: map[string]int{}, subscribers: map[string][2]int{}}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return TransactionST{}, fmt.Errorf("x12837p: parsing transaction at segment %d: missing SE", st.index)
		case "SE":
			se := p.take()
			count, control, err := parseTrailer(se)
			if err != nil {
				return TransactionST{}, err
			}
			t.SegmentCountSE01 = count
			t.ControlNumberSE02 = control
			p.transactionLengths = append(p.transactionLengths, se.index-st.index+1)
			return t, nil
		case "ISA", "IEA", "GS", "GE", "ST":
			return TransactionST{}, fmt.Errorf("x12837p: parsing segment %d: %s inside the transaction at segment %d (missing SE)", next.index, next.id, st.index)
		case "HL":
			if err := p.parseHierarchyLoop(&t, tree); err != nil {
				return TransactionST{}, err
			}
		case "NM1":
			if err := p.parseHeaderName(&t); err != nil {
				return TransactionST{}, err
			}
		default:
			if err := t.add(p.take()); err != nil {
				return TransactionST{}, err
			}
		}
	}
}

// parseHeaderName places an NM1 before the first HL: NM1*41 opens the
// 1000A submitter loop and NM1*40 the 1000B receiver loop, each the first
// time. Any other NM1 loop is kept raw on the transaction.
func (p *parser) parseHeaderName(t *TransactionST) error {
	entity := EntityIdentifierNM101(p.peek().el(1))
	switch {
	case entity == EntityIdentifierNM10141 && t.SubmitterNM141.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		t.SubmitterNM141 = submitterNameLoop(p.takeLoop(submitterMembers))
	case entity == EntityIdentifierNM10140 && t.ReceiverNM140.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		receiver, err := receiverNameLoop(p.takeLoop(receiverMembers))
		if err != nil {
			return err
		}
		t.ReceiverNM140 = receiver
	default:
		t.Raw = append(t.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
	}
	return nil
}

// parseHierarchyLoop reads one HL loop and attaches it by its HL03 level and
// HL02 parent, or keeps it raw on the transaction when it has no slot.
func (p *parser) parseHierarchyLoop(t *TransactionST, tree hierarchy) error {
	hl := p.peek()
	level := HierarchicalLevelHL03(hl.el(3))
	parent := hl.el(2)
	switch level {
	case HierarchicalLevelHL0320:
		provider, err := p.parseBillingProvider()
		if err != nil {
			return err
		}
		t.BillingProviders = append(t.BillingProviders, provider)
		tree.billingProviders[provider.HL.IDNumberHL01] = len(t.BillingProviders) - 1
		return nil
	case HierarchicalLevelHL0322:
		b, ok := tree.billingProviders[parent]
		if !ok {
			break
		}
		subscriber, err := p.parseSubscriber()
		if err != nil {
			return err
		}
		t.BillingProviders[b].Subscribers = append(t.BillingProviders[b].Subscribers, subscriber)
		tree.subscribers[subscriber.HL.IDNumberHL01] = [2]int{b, len(t.BillingProviders[b].Subscribers) - 1}
		return nil
	case HierarchicalLevelHL0323:
		at, ok := tree.subscribers[parent]
		if !ok {
			break
		}
		patient, err := p.parsePatient()
		if err != nil {
			return err
		}
		b, s := at[0], at[1]
		t.BillingProviders[b].Subscribers[s].Patients = append(t.BillingProviders[b].Subscribers[s].Patients, patient)
		return nil
	}
	t.Raw = append(t.Raw, rawSegment(p.take()))
	t.Raw = append(t.Raw, p.takeUntilHierarchy()...)
	return nil
}

// parseBillingProvider reads a 2000A loop: the HL, its PRV and CUR, and the
// 2010AA (NM1*85), 2010AB (NM1*87), and 2010AC (NM1*PE) loops. A claim has
// no slot in a 2000A loop, so a CLM there is kept raw with every segment up
// to the next HL.
func (p *parser) parseBillingProvider() (BillingProviderHL, error) {
	b := BillingProviderHL{HL: parseHL(p.take())}
	for !endsHierarchy(p.peek().id) {
		next := p.peek()
		switch next.id {
		case "CLM":
			b.Raw = append(b.Raw, p.takeUntilHierarchy()...)
		case "NM1":
			entity := EntityIdentifierNM101(next.el(1))
			var err error
			switch {
			case entity == EntityIdentifierNM10185 && b.ProviderNM185.NM1.EntityIDNM101 == EntityIdentifierNM101NotSent:
				b.ProviderNM185, err = billingProviderNameLoop(p.takeLoop(billingProviderMembers))
			case entity == EntityIdentifierNM10187 && b.PayToProviderNM187.NM1.EntityIDNM101 == EntityIdentifierNM101NotSent:
				b.PayToProviderNM187, err = payToProviderNameLoop(p.takeLoop(payToProviderMembers))
			case entity == EntityIdentifierNM101PE && b.PayToPlanNM1PE.NM1.EntityIDNM101 == EntityIdentifierNM101NotSent:
				b.PayToPlanNM1PE, err = payToPlanNameLoop(p.takeLoop(payToPlanMembers))
			default:
				b.Raw = append(b.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
			}
			if err != nil {
				return BillingProviderHL{}, err
			}
		default:
			b.add(p.take())
		}
	}
	return b, nil
}

// parseSubscriber reads a 2000B loop: the HL, its SBR and PAT, the 2010BA
// (NM1*IL) and 2010BB (NM1*PR) loops, and its claims.
func (p *parser) parseSubscriber() (SubscriberHL, error) {
	sub := SubscriberHL{HL: parseHL(p.take())}
	for !endsHierarchy(p.peek().id) {
		next := p.peek()
		var err error
		switch next.id {
		case "CLM":
			var claim ClaimCLM
			claim, err = p.parseClaim()
			sub.Claims = append(sub.Claims, claim)
		case "NM1":
			entity := EntityIdentifierNM101(next.el(1))
			switch {
			case entity == EntityIdentifierNM101IL && sub.SubscriberNM1IL.NM1.EntityIDNM101 == EntityIdentifierNM101NotSent:
				sub.SubscriberNM1IL, err = subscriberNameLoop(p.takeLoop(subscriberMembers))
			case entity == EntityIdentifierNM101PR && sub.PayerNM1PR.NM1.EntityIDNM101 == EntityIdentifierNM101NotSent:
				sub.PayerNM1PR, err = payerNameLoop(p.takeLoop(payerMembers))
			default:
				sub.Raw = append(sub.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
			}
		default:
			err = sub.add(p.take())
		}
		if err != nil {
			return SubscriberHL{}, err
		}
	}
	return sub, nil
}

// parsePatient reads a 2000C loop: the HL, its PAT, the 2010CA (NM1*QC)
// loop, and its claims.
func (p *parser) parsePatient() (PatientHL, error) {
	pat := PatientHL{HL: parseHL(p.take())}
	for !endsHierarchy(p.peek().id) {
		next := p.peek()
		var err error
		switch next.id {
		case "CLM":
			var claim ClaimCLM
			claim, err = p.parseClaim()
			pat.Claims = append(pat.Claims, claim)
		case "NM1":
			entity := EntityIdentifierNM101(next.el(1))
			if entity != EntityIdentifierNM101QC || pat.PatientNM1QC.NM1.EntityIDNM101 != EntityIdentifierNM101NotSent {
				pat.Raw = append(pat.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
				continue
			}
			pat.PatientNM1QC, err = patientNameLoop(p.takeLoop(patientMembers))
		default:
			err = pat.add(p.take())
		}
		if err != nil {
			return PatientHL{}, err
		}
	}
	return pat, nil
}

// parseClaim reads a 2300 loop: the CLM, its own segments, the 2310A-F
// provider loops, the 2320 other-subscriber loops, and the 2400 service
// lines.
func (p *parser) parseClaim() (ClaimCLM, error) {
	c, err := parseCLM(p.take(), p.component)
	if err != nil {
		return ClaimCLM{}, err
	}
	for !endsClaim(p.peek().id) {
		var err error
		switch p.peek().id {
		case "NM1":
			err = p.parseClaimProvider(&c)
		case "SBR":
			var other OtherSubscriberSBR
			other, err = p.parseOtherSubscriber()
			c.OtherSubscribers = append(c.OtherSubscribers, other)
		case "LX":
			var line ServiceLineLX
			line, err = p.parseServiceLine()
			c.ServiceLines = append(c.ServiceLines, line)
		default:
			err = c.add(p.take(), p.component)
		}
		if err != nil {
			return ClaimCLM{}, err
		}
	}
	return c, nil
}

// parseClaimProvider places a claim-level NM1 loop by NM101: DN or P3 for
// 2310A, 82 for 2310B, 77 for 2310C, DQ for 2310D, PW for 2310E, and 45 for
// 2310F, each the first time. Any other NM1 loop is kept raw on the claim.
func (p *parser) parseClaimProvider(c *ClaimCLM) error {
	entity := EntityIdentifierNM101(p.peek().el(1))
	var err error
	switch {
	case (entity == EntityIdentifierNM101DN || entity == EntityIdentifierNM101P3) && c.ReferringProvider.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.ReferringProvider, err = claimReferringLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM10182 && c.RenderingProviderNM182.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.RenderingProviderNM182, err = claimRenderingLoop(p.takeLoop(renderingMembers))
	case entity == EntityIdentifierNM10177 && c.ServiceFacilityNM177.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.ServiceFacilityNM177, err = claimServiceFacilityLoop(p.takeLoop(serviceFacilityMembers))
	case entity == EntityIdentifierNM101DQ && c.SupervisingProviderNM1DQ.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.SupervisingProviderNM1DQ, err = claimSupervisingLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM101PW && c.AmbulancePickupNM1PW.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.AmbulancePickupNM1PW, err = claimAmbulancePickupLoop(p.takeLoop(locationMembers))
	case entity == EntityIdentifierNM10145 && c.AmbulanceDropoffNM145.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		c.AmbulanceDropoffNM145, err = claimAmbulanceDropoffLoop(p.takeLoop(locationMembers))
	default:
		c.Raw = append(c.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
	}
	return err
}

// parseOtherSubscriber reads a 2320 loop: the SBR, its CAS, AMT, OI, and
// MOA, and the 2330A-G name loops.
func (p *parser) parseOtherSubscriber() (OtherSubscriberSBR, error) {
	o := OtherSubscriberSBR{Subscriber: parseSBR(p.take())}
	for !endsOtherSubscriber(p.peek().id) {
		var err error
		if p.peek().id == "NM1" {
			err = p.parseOtherPayerName(&o)
		} else {
			err = o.add(p.take())
		}
		if err != nil {
			return OtherSubscriberSBR{}, err
		}
	}
	return o, nil
}

// parseOtherPayerName places a 2320 NM1 loop by NM101: IL for 2330A, PR for
// 2330B, DN or P3 for 2330C, 82 for 2330D, 77 for 2330E, DQ for 2330F, and
// 85 for 2330G, each the first time. Any other NM1 loop is kept raw on the
// 2320 loop.
func (p *parser) parseOtherPayerName(o *OtherSubscriberSBR) error {
	entity := EntityIdentifierNM101(p.peek().el(1))
	var err error
	switch {
	case entity == EntityIdentifierNM101IL && o.OtherSubscriberNM1IL.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherSubscriberNM1IL = otherSubscriberNameLoop(p.takeLoop(otherSubscriberMembers))
	case entity == EntityIdentifierNM101PR && o.OtherPayerNM1PR.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerNM1PR, err = otherPayerNameLoop(p.takeLoop(otherPayerMembers))
	case (entity == EntityIdentifierNM101DN || entity == EntityIdentifierNM101P3) && o.OtherPayerReferringProvider.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerReferringProvider, err = otherPayerRoleLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM10182 && o.OtherPayerRenderingProviderNM182.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerRenderingProviderNM182, err = otherPayerRoleLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM10177 && o.OtherPayerServiceFacilityNM177.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerServiceFacilityNM177, err = otherPayerRoleLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM101DQ && o.OtherPayerSupervisingProviderNM1DQ.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerSupervisingProviderNM1DQ, err = otherPayerRoleLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM10185 && o.OtherPayerBillingProviderNM185.Name.EntityIDNM101 == EntityIdentifierNM101NotSent:
		o.OtherPayerBillingProviderNM185, err = otherPayerRoleLoop(p.takeLoop(referenceMembers))
	default:
		o.Raw = append(o.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
	}
	return err
}

// parseServiceLine reads a 2400 loop: the LX, its own segments, and the
// 2410, 2420A-H, 2430, and 2440 loops.
func (p *parser) parseServiceLine() (ServiceLineLX, error) {
	l := parseLX(p.take())
	for !endsServiceLine(p.peek().id) {
		var err error
		switch p.peek().id {
		case "LIN":
			if l.Drug.Set {
				l.Raw = append(l.Raw, rawSegments(p.takeLoop(drugMembers))...)
				continue
			}
			l.Drug, err = drugIdentificationLoop(p.takeLoop(drugMembers), p.component)
		case "NM1":
			err = p.parseLineProvider(&l)
		case "SVD":
			var adjudication LineAdjudicationSVD
			adjudication, err = lineAdjudicationLoop(p.takeLoop(lineAdjudicationMembers), p.component)
			l.Adjudications = append(l.Adjudications, adjudication)
		case "LQ":
			var form FormLQ
			form, err = supportingFormLoop(p.takeLoop(supportingFormMembers))
			l.Forms = append(l.Forms, form)
		default:
			err = l.add(p.take(), p.component)
		}
		if err != nil {
			return ServiceLineLX{}, err
		}
	}
	return l, nil
}

// parseLineProvider places a line-level NM1 loop by NM101: 82 for 2420A, QB
// for 2420B, 77 for 2420C, DQ for 2420D, DK for 2420E, DN or P3 for 2420F,
// PW for 2420G, and 45 for 2420H, each the first time. Any other NM1 loop is
// kept raw on the service line.
func (p *parser) parseLineProvider(l *ServiceLineLX) error {
	entity := EntityIdentifierNM101(p.peek().el(1))
	var err error
	switch {
	case entity == EntityIdentifierNM10182 && !l.RenderingNM182.Set:
		l.RenderingNM182, err = lineProviderLoop(p.takeLoop(renderingMembers))
	case entity == EntityIdentifierNM101QB && !l.PurchasedServiceNM1QB.Set:
		l.PurchasedServiceNM1QB, err = lineProviderLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM10177 && !l.ServiceFacilityNM177.Set:
		l.ServiceFacilityNM177, err = lineProviderLoop(p.takeLoop(lineFacilityMembers))
	case entity == EntityIdentifierNM101DQ && !l.SupervisingNM1DQ.Set:
		l.SupervisingNM1DQ, err = lineProviderLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM101DK && !l.OrderingNM1DK.Set:
		l.OrderingNM1DK, err = lineProviderLoop(p.takeLoop(orderingProviderMembers))
	case (entity == EntityIdentifierNM101DN || entity == EntityIdentifierNM101P3) && !l.ReferringNM12420F.Set:
		l.ReferringNM12420F, err = lineProviderLoop(p.takeLoop(referenceMembers))
	case entity == EntityIdentifierNM101PW && !l.AmbulancePickupNM1PW.Set:
		l.AmbulancePickupNM1PW, err = lineProviderLoop(p.takeLoop(locationMembers))
	case entity == EntityIdentifierNM10145 && !l.AmbulanceDropoffNM145.Set:
		l.AmbulanceDropoffNM145, err = lineProviderLoop(p.takeLoop(locationMembers))
	default:
		l.Raw = append(l.Raw, rawSegments(p.takeLoop(anyNameMembers))...)
	}
	return err
}

// checkEnvelopes compares each envelope trailer with its header and with
// what the pass found, and returns every mismatch. transactionLengths holds,
// in file order, the number of segments from each ST through its SE.
func checkEnvelopes(isa InterchangeISA, transactionLengths []int) error {
	var errs []error
	if isa.ControlNumberISA13 != isa.ControlNumberIEA02 {
		errs = append(errs, fmt.Errorf("x12837p: checking IEA02: ISA13 is %q, IEA02 is %q", isa.ControlNumberISA13, isa.ControlNumberIEA02))
	}
	if isa.GroupCountIEA01 != len(isa.FunctionalGroups) {
		errs = append(errs, fmt.Errorf("x12837p: checking IEA01: IEA01 is %d, the interchange holds %d functional groups", isa.GroupCountIEA01, len(isa.FunctionalGroups)))
	}
	k := 0
	for _, g := range isa.FunctionalGroups {
		if g.ControlNumberGS06 != g.ControlNumberGE02 {
			errs = append(errs, fmt.Errorf("x12837p: checking GE02: GS06 is %q, GE02 is %q", g.ControlNumberGS06, g.ControlNumberGE02))
		}
		if g.TransactionCountGE01 != len(g.Transactions) {
			errs = append(errs, fmt.Errorf("x12837p: checking GE01 of group %q: GE01 is %d, the group holds %d transactions", g.ControlNumberGS06, g.TransactionCountGE01, len(g.Transactions)))
		}
		for _, t := range g.Transactions {
			if t.ControlNumberST02 != t.ControlNumberSE02 {
				errs = append(errs, fmt.Errorf("x12837p: checking SE02: ST02 is %q, SE02 is %q", t.ControlNumberST02, t.ControlNumberSE02))
			}
			if t.SegmentCountSE01 != transactionLengths[k] {
				errs = append(errs, fmt.Errorf("x12837p: checking SE01 of transaction %q: SE01 is %d, ST through SE is %d segments", t.ControlNumberST02, t.SegmentCountSE01, transactionLengths[k]))
			}
			k++
		}
	}
	return errors.Join(errs...)
}

// add places one segment that sits directly under ST outside any loop: the
// first BHT. Anything else is kept raw.
func (t *TransactionST) add(s inputSegment) error {
	if s.id != "BHT" || t.Header.Set {
		t.Raw = append(t.Raw, rawSegment(s))
		return nil
	}
	header, err := parseBHT(s)
	if err != nil {
		return err
	}
	t.Header = header
	return nil
}

// add places one segment of a 2000A loop outside its name loops: the first
// PRV that carries PRV01 and the first CUR. Anything else is kept raw.
func (b *BillingProviderHL) add(s inputSegment) {
	switch s.id {
	case "PRV":
		addSpecialty(s, &b.Specialty, &b.Raw)
		return
	case "CUR":
		if !b.Currency.Set {
			b.Currency = parseCUR(s)
			return
		}
	}
	b.Raw = append(b.Raw, rawSegment(s))
}

// add places one segment of a 2000B loop outside its name loops and claims:
// the first SBR and the first PAT. Anything else is kept raw.
func (sub *SubscriberHL) add(s inputSegment) error {
	switch s.id {
	case "SBR":
		if sub.Info.Set {
			break
		}
		sub.Info = parseSBR(s)
		return nil
	case "PAT":
		if sub.PatientInfo.Set {
			break
		}
		info, err := parsePAT(s)
		if err != nil {
			return err
		}
		sub.PatientInfo = info
		return nil
	}
	sub.Raw = append(sub.Raw, rawSegment(s))
	return nil
}

// add places one segment of a 2000C loop outside its name loop and claims:
// the first PAT. Anything else is kept raw.
func (pat *PatientHL) add(s inputSegment) error {
	if s.id != "PAT" || pat.PatientInfo.Set {
		pat.Raw = append(pat.Raw, rawSegment(s))
		return nil
	}
	info, err := parsePAT(s)
	if err != nil {
		return err
	}
	pat.PatientInfo = info
	return nil
}

// nameLoop holds what a name loop can carry after its NM1. Each loop type
// copies the parts its box holds; the segments a loop takes (its
// loopMembers) never include one it has no field for.
type nameLoop struct {
	name         NameNM1
	address      AddressN3
	location     LocationN4
	demographics DemographicDMG
	contact      ContactPER
	specialty    ProviderPRV
	references   []ReferenceREF
	raw          []RawSegment
}

// readNameLoop reads an NM1 and the segments of its loop: the first N3, N4,
// and DMG, the first PER that carries PER01, the first PRV that carries
// PRV01, and every REF. Anything else, a repeat included, is kept raw.
func readNameLoop(loop []inputSegment) (nameLoop, error) {
	n := nameLoop{name: parseNM1(loop[0])}
	for _, s := range loop[1:] {
		switch s.id {
		case "N3":
			addAddress(s, &n.address, &n.raw)
		case "N4":
			addLocation(s, &n.location, &n.raw)
		case "DMG":
			if err := addDemographic(s, &n.demographics, &n.raw); err != nil {
				return nameLoop{}, err
			}
		case "PER":
			addContact(s, &n.contact, &n.raw)
		case "PRV":
			addSpecialty(s, &n.specialty, &n.raw)
		case "REF":
			n.references = append(n.references, parseREF(s))
		default:
			n.raw = append(n.raw, rawSegment(s))
		}
	}
	return n, nil
}

// submitterNameLoop reads the 1000A loop: the NM1 and every PER.
func submitterNameLoop(loop []inputSegment) SubmitterNM141 {
	l := SubmitterNM141{Name: parseNM1(loop[0])}
	for _, s := range loop[1:] {
		if s.id != "PER" {
			l.Raw = append(l.Raw, rawSegment(s))
			continue
		}
		l.Contacts = append(l.Contacts, parsePER(s))
	}
	return l
}

func receiverNameLoop(loop []inputSegment) (ReceiverNM140, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ReceiverNM140{}, err
	}
	return ReceiverNM140{Name: n.name, Raw: n.raw}, nil
}

// billingProviderNameLoop reads the 2010AA loop. The first REF*EI with a
// non-empty REF02 fills TaxIDREFEI; every other REF goes to References.
func billingProviderNameLoop(loop []inputSegment) (BillingProviderNM185, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return BillingProviderNM185{}, err
	}
	b := BillingProviderNM185{NM1: n.name, Address: n.address, Location: n.location, Contact: n.contact, Raw: n.raw}
	for _, ref := range n.references {
		if ref.QualifierREF01 == ReferenceIdentificationQualifierREF01EI && b.TaxIDREFEI.QualifierREF01 == ReferenceIdentificationQualifierREF01NotSent && ref.Value != "" {
			b.TaxIDREFEI = ref
			continue
		}
		b.References = append(b.References, ref)
	}
	return b, nil
}

func payToProviderNameLoop(loop []inputSegment) (PayToProviderNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return PayToProviderNM1{}, err
	}
	return PayToProviderNM1{NM1: n.name, Address: n.address, Location: n.location, Raw: n.raw}, nil
}

func payToPlanNameLoop(loop []inputSegment) (PayToPlanNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return PayToPlanNM1{}, err
	}
	return PayToPlanNM1{NM1: n.name, Address: n.address, Location: n.location, References: n.references, Raw: n.raw}, nil
}

func subscriberNameLoop(loop []inputSegment) (SubscriberNM1IL, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return SubscriberNM1IL{}, err
	}
	return SubscriberNM1IL{NM1: n.name, Address: n.address, Location: n.location, Demographics: n.demographics, References: n.references, Contact: n.contact, Raw: n.raw}, nil
}

func payerNameLoop(loop []inputSegment) (PayerNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return PayerNM1{}, err
	}
	return PayerNM1{NM1: n.name, Address: n.address, Location: n.location, References: n.references, Raw: n.raw}, nil
}

func patientNameLoop(loop []inputSegment) (PatientNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return PatientNM1{}, err
	}
	return PatientNM1{NM1: n.name, Address: n.address, Location: n.location, Demographics: n.demographics, References: n.references, Contact: n.contact, Raw: n.raw}, nil
}

func claimReferringLoop(loop []inputSegment) (ClaimReferringProviderNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimReferringProviderNM1{}, err
	}
	return ClaimReferringProviderNM1{Name: n.name, References: n.references, Raw: n.raw}, nil
}

func claimRenderingLoop(loop []inputSegment) (ClaimRenderingProviderNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimRenderingProviderNM1{}, err
	}
	return ClaimRenderingProviderNM1{Name: n.name, Specialty: n.specialty, References: n.references, Raw: n.raw}, nil
}

func claimServiceFacilityLoop(loop []inputSegment) (ClaimServiceFacilityNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimServiceFacilityNM1{}, err
	}
	return ClaimServiceFacilityNM1{Name: n.name, Address: n.address, Location: n.location, References: n.references, Contact: n.contact, Raw: n.raw}, nil
}

func claimSupervisingLoop(loop []inputSegment) (ClaimSupervisingProviderNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimSupervisingProviderNM1{}, err
	}
	return ClaimSupervisingProviderNM1{Name: n.name, References: n.references, Raw: n.raw}, nil
}

func claimAmbulancePickupLoop(loop []inputSegment) (ClaimAmbulancePickupNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimAmbulancePickupNM1{}, err
	}
	return ClaimAmbulancePickupNM1{Name: n.name, Address: n.address, Location: n.location, Raw: n.raw}, nil
}

func claimAmbulanceDropoffLoop(loop []inputSegment) (ClaimAmbulanceDropoffNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return ClaimAmbulanceDropoffNM1{}, err
	}
	return ClaimAmbulanceDropoffNM1{Name: n.name, Address: n.address, Location: n.location, Raw: n.raw}, nil
}

// otherSubscriberNameLoop reads the 2330A loop: the NM1, the first N3 and
// N4, and the first REF that carries REF01, which fills SecondaryID.
// Anything else is kept raw.
func otherSubscriberNameLoop(loop []inputSegment) OtherSubscriberNM1 {
	o := OtherSubscriberNM1{Name: parseNM1(loop[0])}
	for _, s := range loop[1:] {
		switch {
		case s.id == "N3":
			addAddress(s, &o.Address, &o.Raw)
		case s.id == "N4":
			addLocation(s, &o.Location, &o.Raw)
		case s.id == "REF" && s.el(1) != "" && o.SecondaryID.QualifierREF01 == ReferenceIdentificationQualifierREF01NotSent:
			o.SecondaryID = parseREF(s)
		default:
			o.Raw = append(o.Raw, rawSegment(s))
		}
	}
	return o
}

// otherPayerNameLoop reads the 2330B loop: the NM1, the first N3 and N4,
// the first DTP*573, every REF*2U, EI, FY, and NF in SecondaryID, and the
// first REF*G1, 9F, T4, and F8 with a non-empty REF02 in their own fields.
// Anything else is kept raw.
func otherPayerNameLoop(loop []inputSegment) (OtherPayerNM1, error) {
	o := OtherPayerNM1{Name: parseNM1(loop[0])}
	for _, s := range loop[1:] {
		switch s.id {
		case "N3":
			addAddress(s, &o.Address, &o.Raw)
		case "N4":
			addLocation(s, &o.Location, &o.Raw)
		case "DTP":
			qualifier := DateTimeQualifierDTP01(s.el(1))
			if qualifier != DateTimeQualifierDTP01573 || o.RemittanceDateDTP573.QualifierDTP01 != DateTimeQualifierDTP01NotSent {
				o.Raw = append(o.Raw, rawSegment(s))
				continue
			}
			date, err := parseDTP(s)
			if err != nil {
				return OtherPayerNM1{}, err
			}
			o.RemittanceDateDTP573 = date
		case "REF":
			o.addReference(s)
		default:
			o.Raw = append(o.Raw, rawSegment(s))
		}
	}
	return o, nil
}

func (o *OtherPayerNM1) addReference(s inputSegment) {
	ref := parseREF(s)
	var slot *ReferenceREF
	switch ref.QualifierREF01 {
	case ReferenceIdentificationQualifierREF012U, ReferenceIdentificationQualifierREF01EI,
		ReferenceIdentificationQualifierREF01FY, ReferenceIdentificationQualifierREF01NF:
		o.SecondaryID = append(o.SecondaryID, ref)
		return
	case ReferenceIdentificationQualifierREF01G1:
		slot = &o.PriorAuthorizationREFG1
	case ReferenceIdentificationQualifierREF019F:
		slot = &o.ReferralREF9F
	case ReferenceIdentificationQualifierREF01T4:
		slot = &o.ClaimAdjustmentIndicatorREFT4
	case ReferenceIdentificationQualifierREF01F8:
		slot = &o.ClaimControlNumberREFF8
	}
	if slot == nil || slot.QualifierREF01 != ReferenceIdentificationQualifierREF01NotSent || ref.Value == "" {
		o.Raw = append(o.Raw, rawSegment(s))
		return
	}
	*slot = ref
}

func otherPayerRoleLoop(loop []inputSegment) (OtherPayerRoleNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return OtherPayerRoleNM1{}, err
	}
	return OtherPayerRoleNM1{Name: n.name, References: n.references, Raw: n.raw}, nil
}

func lineProviderLoop(loop []inputSegment) (LineProviderNM1, error) {
	n, err := readNameLoop(loop)
	if err != nil {
		return LineProviderNM1{}, err
	}
	return LineProviderNM1{Set: true, Name: n.name, Specialty: n.specialty, Address: n.address, Location: n.location, Contact: n.contact, References: n.references, Raw: n.raw}, nil
}

// add places one segment of a claim outside its nested loops: every PWK and
// K3; every REF, the first with a non-empty REF02 for each claim REF box
// going to its field; the first DTP for each claim DTP box (two for Assumed
// and Relinquished Care); the first AMT*F5 that carries an amount; every CRC
// for a repeating CRC box and the first for the others; the first HI for the
// diagnosis and anesthesia boxes and every HI for the condition box; and the
// first CN1, NTE, CR1, CR2, and HCP. Anything else is kept raw.
func (c *ClaimCLM) add(s inputSegment, component byte) error {
	switch s.id {
	case "DTP":
		return c.addDate(s)
	case "REF":
		c.addReference(s)
		return nil
	case "CRC":
		c.addCondition(s)
		return nil
	case "HI":
		return c.addHealthCare(s, component)
	case "PWK":
		pwk, err := parsePWK(s)
		if err != nil {
			return err
		}
		c.SupplementalInfoPWK = append(c.SupplementalInfoPWK, pwk)
		return nil
	case "K3":
		c.FileInfoK3 = append(c.FileInfoK3, parseK3(s))
		return nil
	case "AMT":
		if AmountQualifierAMT01(s.el(1)) != AmountQualifierAMT01F5 || c.PatientAmountPaidAMTF5.Amount.Set || s.el(2) == "" {
			break
		}
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		c.PatientAmountPaidAMTF5 = amount
		return nil
	case "CN1":
		if c.ContractInfoCN1.Set {
			break
		}
		contract, err := parseCN1(s)
		if err != nil {
			return err
		}
		c.ContractInfoCN1 = contract
		return nil
	case "NTE":
		if c.NoteNTE.Set {
			break
		}
		c.NoteNTE = parseNTE(s)
		return nil
	case "CR1":
		if c.AmbulanceTransportCR1.Set {
			break
		}
		transport, err := parseCR1(s)
		if err != nil {
			return err
		}
		c.AmbulanceTransportCR1 = transport
		return nil
	case "CR2":
		if c.SpinalManipulationCR2.Set {
			break
		}
		spinal, err := parseCR2(s)
		if err != nil {
			return err
		}
		c.SpinalManipulationCR2 = spinal
		return nil
	case "HCP":
		if c.PricingHCP.Set {
			break
		}
		pricing, err := parseHCP(s)
		if err != nil {
			return err
		}
		c.PricingHCP = pricing
		return nil
	}
	c.Raw = append(c.Raw, rawSegment(s))
	return nil
}

// addDate places a claim DTP in the field for its DTP01 code when that field
// is still empty; otherwise the DTP is kept raw.
func (c *ClaimCLM) addDate(s inputSegment) error {
	var slot *DateDTP
	switch DateTimeQualifierDTP01(s.el(1)) {
	case DateTimeQualifierDTP01431:
		slot = &c.OnsetDTP431
	case DateTimeQualifierDTP01454:
		slot = &c.InitialTreatmentDTP454
	case DateTimeQualifierDTP01304:
		slot = &c.LastSeenDTP304
	case DateTimeQualifierDTP01453:
		slot = &c.AcuteManifestationDTP453
	case DateTimeQualifierDTP01439:
		slot = &c.AccidentDTP439
	case DateTimeQualifierDTP01484:
		slot = &c.LastMenstrualPeriodDTP484
	case DateTimeQualifierDTP01455:
		slot = &c.LastXRayDTP455
	case DateTimeQualifierDTP01471:
		slot = &c.HearingVisionPrescriptionDTP471
	case DateTimeQualifierDTP01314, DateTimeQualifierDTP01360, DateTimeQualifierDTP01361:
		slot = &c.DisabilityDTP314
	case DateTimeQualifierDTP01297:
		slot = &c.LastWorkedDTP297
	case DateTimeQualifierDTP01296:
		slot = &c.AuthorizedReturnToWorkDTP296
	case DateTimeQualifierDTP01435:
		slot = &c.AdmissionDTP435
	case DateTimeQualifierDTP01096:
		slot = &c.DischargeDTP096
	case DateTimeQualifierDTP01090, DateTimeQualifierDTP01091:
		slot = &c.AssumedOrRelinquishedCareDTP090[0]
		if slot.QualifierDTP01 != DateTimeQualifierDTP01NotSent {
			slot = &c.AssumedOrRelinquishedCareDTP090[1]
		}
	case DateTimeQualifierDTP01444:
		slot = &c.PropertyCasualtyFirstContactDTP444
	case DateTimeQualifierDTP01050:
		slot = &c.RepricerReceivedDTP050
	}
	if slot == nil || slot.QualifierDTP01 != DateTimeQualifierDTP01NotSent {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	date, err := parseDTP(s)
	if err != nil {
		return err
	}
	*slot = date
	return nil
}

// addReference places a claim REF in the field for its REF01 code when that
// field is still empty and REF02 carries text; otherwise the REF goes to
// References.
func (c *ClaimCLM) addReference(s inputSegment) {
	ref := parseREF(s)
	var slot *ReferenceREF
	switch ref.QualifierREF01 {
	case ReferenceIdentificationQualifierREF014N:
		slot = &c.ServiceAuthExceptionREF4N
	case ReferenceIdentificationQualifierREF01F5:
		slot = &c.MedicareCrossoverIndicatorREFF5
	case ReferenceIdentificationQualifierREF01EW:
		slot = &c.MammographyCertificationREFEW
	case ReferenceIdentificationQualifierREF019F:
		slot = &c.ReferralREF9F
	case ReferenceIdentificationQualifierREF01G1:
		slot = &c.PriorAuthorizationREFG1
	case ReferenceIdentificationQualifierREF01F8:
		slot = &c.PayerClaimControlNumberREFF8
	case ReferenceIdentificationQualifierREF01X4:
		slot = &c.CLIAREFX4
	case ReferenceIdentificationQualifierREF019A:
		slot = &c.RepricedClaimREF9A
	case ReferenceIdentificationQualifierREF019C:
		slot = &c.AdjustedRepricedClaimREF9C
	case ReferenceIdentificationQualifierREF01LX:
		slot = &c.InvestigationalDeviceExemptionREFLX
	case ReferenceIdentificationQualifierREF01D9:
		slot = &c.TransmissionIntermediaryREFD9
	case ReferenceIdentificationQualifierREF01EA:
		slot = &c.MedicalRecordNumberREFEA
	case ReferenceIdentificationQualifierREF01P4:
		slot = &c.DemonstrationProjectREFP4
	case ReferenceIdentificationQualifierREF011J:
		slot = &c.CarePlanOversightREF1J
	}
	if slot == nil || slot.QualifierREF01 != ReferenceIdentificationQualifierREF01NotSent || ref.Value == "" {
		c.References = append(c.References, ref)
		return
	}
	*slot = ref
}

// addCondition places a claim CRC by its CRC01 code: every CRC*07 and every
// CRC*E1, E2, and E3 in their repeating boxes, and the first CRC*75 and
// CRC*ZZ in theirs. Anything else is kept raw.
func (c *ClaimCLM) addCondition(s inputSegment) {
	category := CodeCategoryCRC01(s.el(1))
	switch category {
	case CodeCategoryCRC0107:
		c.AmbulanceCertificationCRC07 = append(c.AmbulanceCertificationCRC07, parseCRC(s))
		return
	case CodeCategoryCRC01E1, CodeCategoryCRC01E2, CodeCategoryCRC01E3:
		c.VisionConditionCRCE1 = append(c.VisionConditionCRCE1, parseCRC(s))
		return
	case CodeCategoryCRC0175, CodeCategoryCRC01ZZ:
		slot := &c.HomeboundCRC75
		if category == CodeCategoryCRC01ZZ {
			slot = &c.EPSDTReferralCRCZZ
		}
		if !slot.Set {
			*slot = parseCRC(s)
			return
		}
	}
	c.Raw = append(c.Raw, rawSegment(s))
}

// addHealthCare places a claim HI by HI01-1: ABK or BK opens the diagnosis
// box, BP the anesthesia box, each once, and BG the condition box, every
// time. Anything else is kept raw.
func (c *ClaimCLM) addHealthCare(s inputSegment, component byte) error {
	first, _, _ := strings.Cut(s.el(1), string(component))
	qualifier := DiagnosisTypeCodeHI0x_1(first)
	open := false
	switch qualifier {
	case DiagnosisTypeCodeHI0x_1ABK, DiagnosisTypeCodeHI0x_1BK:
		open = !c.DiagnosisCodesHIABK.Set
	case DiagnosisTypeCodeHI0x_1BP:
		open = !c.AnesthesiaProceduresHIBP.Set
	case DiagnosisTypeCodeHI0x_1BG:
		open = true
	}
	if !open {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	info, err := parseHI(s, component)
	if err != nil {
		return err
	}
	switch qualifier {
	case DiagnosisTypeCodeHI0x_1ABK, DiagnosisTypeCodeHI0x_1BK:
		c.DiagnosisCodesHIABK = info
	case DiagnosisTypeCodeHI0x_1BP:
		c.AnesthesiaProceduresHIBP = info
	default:
		c.ConditionCodesHIBG = append(c.ConditionCodesHIBG, info)
	}
	return nil
}

// add places one segment of a 2320 loop outside its name loops: every CAS,
// the first AMT*D, AMT*A8, and AMT*EAF that carries an amount, and the
// first OI and MOA. Anything else is kept raw.
func (o *OtherSubscriberSBR) add(s inputSegment) error {
	switch s.id {
	case "CAS":
		adjustment, err := parseCAS(s)
		if err != nil {
			return err
		}
		o.Adjustments = append(o.Adjustments, adjustment)
		return nil
	case "AMT":
		var slot *AmountAMT
		switch AmountQualifierAMT01(s.el(1)) {
		case AmountQualifierAMT01D:
			slot = &o.PayerPaidAMTD
		case AmountQualifierAMT01A8:
			slot = &o.NonCoveredAMTA8
		case AmountQualifierAMT01EAF:
			slot = &o.RemainingLiabilityAMTEAF
		}
		if slot == nil || slot.Amount.Set || s.el(2) == "" {
			break
		}
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		*slot = amount
		return nil
	case "OI":
		if o.OtherInsurance.Set {
			break
		}
		o.OtherInsurance = parseOI(s)
		return nil
	case "MOA":
		if o.OutpatientAdjudication.Set {
			break
		}
		adjudication, err := parseMOA(s)
		if err != nil {
			return err
		}
		o.OutpatientAdjudication = adjudication
		return nil
	}
	o.Raw = append(o.Raw, rawSegment(s))
	return nil
}

// add places one segment of a service line outside its nested loops: the
// first SV1, SV5, CR1, CR3, CN1, PS1, and HCP; every PWK, MEA, and K3; the
// first DTP, QTY, AMT, and CRC for each of the line's boxes (every DTP*738
// and 739, and every CRC*07); every REF*G1 and 9F, the first REF for each
// other line REF box, and the first REF*6R with a non-empty REF02; and the
// first NTE for each of the two note boxes. Anything else is kept raw.
func (l *ServiceLineLX) add(s inputSegment, component byte) error {
	switch s.id {
	case "SV1":
		if l.Service.Set {
			break
		}
		service, err := parseSV1(s, component)
		if err != nil {
			return err
		}
		l.Service = service
		return nil
	case "SV5":
		if l.DME.Set {
			break
		}
		dme, err := parseSV5(s, component)
		if err != nil {
			return err
		}
		l.DME = dme
		return nil
	case "PWK":
		pwk, err := parsePWK(s)
		if err != nil {
			return err
		}
		l.Attachments = append(l.Attachments, pwk)
		return nil
	case "CR1":
		if l.AmbulanceTransport.Set {
			break
		}
		transport, err := parseCR1(s)
		if err != nil {
			return err
		}
		l.AmbulanceTransport = transport
		return nil
	case "CR3":
		if l.DMECertification.Set {
			break
		}
		certification, err := parseCR3(s)
		if err != nil {
			return err
		}
		l.DMECertification = certification
		return nil
	case "CRC":
		l.addCondition(s)
		return nil
	case "DTP":
		return l.addDate(s)
	case "QTY":
		var slot *QuantityQTY
		switch QuantityQualifierQTY01(s.el(1)) {
		case QuantityQualifierQTY01PT:
			slot = &l.AmbulancePatientCountQTYPT
		case QuantityQualifierQTY01FL:
			slot = &l.ObstetricAdditionalUnitsQTYFL
		}
		if slot == nil || slot.QualifierQTY01 != QuantityQualifierQTY01NotSent {
			break
		}
		quantity, err := parseQTY(s)
		if err != nil {
			return err
		}
		*slot = quantity
		return nil
	case "MEA":
		result, err := parseMEA(s)
		if err != nil {
			return err
		}
		l.TestResults = append(l.TestResults, result)
		return nil
	case "CN1":
		if l.Contract.Set {
			break
		}
		contract, err := parseCN1(s)
		if err != nil {
			return err
		}
		l.Contract = contract
		return nil
	case "REF":
		l.addReference(s)
		return nil
	case "AMT":
		var slot *AmountAMT
		switch AmountQualifierAMT01(s.el(1)) {
		case AmountQualifierAMT01T:
			slot = &l.SalesTaxAMTT
		case AmountQualifierAMT01F4:
			slot = &l.PostageClaimedAMTF4
		}
		if slot == nil || slot.QualifierAMT01 != AmountQualifierAMT01NotSent {
			break
		}
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		*slot = amount
		return nil
	case "K3":
		l.FileInfoK3 = append(l.FileInfoK3, parseK3(s))
		return nil
	case "NTE":
		slot := &l.LineNote
		if s.el(1) == "TPO" {
			slot = &l.ThirdPartyNote
		}
		if slot.Set {
			break
		}
		*slot = parseNTE(s)
		return nil
	case "PS1":
		if l.PurchasedService.Set {
			break
		}
		purchased, err := parsePS1(s)
		if err != nil {
			return err
		}
		l.PurchasedService = purchased
		return nil
	case "HCP":
		if l.Pricing.Set {
			break
		}
		pricing, err := parseHCP(s)
		if err != nil {
			return err
		}
		l.Pricing = pricing
		return nil
	}
	l.Raw = append(l.Raw, rawSegment(s))
	return nil
}

// addCondition places a service line CRC by its CRC01 code: every CRC*07,
// and the first CRC*70 and CRC*09. Anything else is kept raw.
func (l *ServiceLineLX) addCondition(s inputSegment) {
	var slot *ConditionCRC
	switch CodeCategoryCRC01(s.el(1)) {
	case CodeCategoryCRC0107:
		l.AmbulanceCertificationCRC07 = append(l.AmbulanceCertificationCRC07, parseCRC(s))
		return
	case CodeCategoryCRC0170:
		slot = &l.HospiceEmployeeCRC70
	case CodeCategoryCRC0109:
		slot = &l.DMEConditionCRC09
	}
	if slot == nil || slot.Set {
		l.Raw = append(l.Raw, rawSegment(s))
		return
	}
	*slot = parseCRC(s)
}

// addDate places a service line DTP in the field for its DTP01 code when
// that field is still empty, appending every DTP*738 and 739; otherwise the
// DTP is kept raw.
func (l *ServiceLineLX) addDate(s inputSegment) error {
	var slot *DateDTP
	switch DateTimeQualifierDTP01(s.el(1)) {
	case DateTimeQualifierDTP01472:
		slot = &l.ServiceDateDTP472
	case DateTimeQualifierDTP01471:
		slot = &l.PrescriptionDTP471
	case DateTimeQualifierDTP01607:
		slot = &l.CertificationRevisionDTP607
	case DateTimeQualifierDTP01463:
		slot = &l.BeginTherapyDTP463
	case DateTimeQualifierDTP01461:
		slot = &l.LastCertificationDTP461
	case DateTimeQualifierDTP01304:
		slot = &l.LastTreatmentOrTherapyDTP304
	case DateTimeQualifierDTP01011:
		slot = &l.ShippedDTP011
	case DateTimeQualifierDTP01455:
		slot = &l.LastXRayDTP455
	case DateTimeQualifierDTP01454:
		slot = &l.InitialTreatmentDTP454
	case DateTimeQualifierDTP01738, DateTimeQualifierDTP01739:
		date, err := parseDTP(s)
		if err != nil {
			return err
		}
		l.TestPerformedDTP738 = append(l.TestPerformedDTP738, date)
		return nil
	}
	if slot == nil || slot.QualifierDTP01 != DateTimeQualifierDTP01NotSent {
		l.Raw = append(l.Raw, rawSegment(s))
		return nil
	}
	date, err := parseDTP(s)
	if err != nil {
		return err
	}
	*slot = date
	return nil
}

// addReference places a service line REF by its REF01 code: every REF*G1
// and 9F in their repeating boxes, the first REF*6R with a non-empty REF02
// in LineControlREF6R, and the first REF*9B, 9D, EW, X4, F4, and BT in
// theirs. Anything else is kept raw.
func (l *ServiceLineLX) addReference(s inputSegment) {
	ref := parseREF(s)
	var slot *ReferenceREF
	switch ref.QualifierREF01 {
	case ReferenceIdentificationQualifierREF01G1:
		l.PriorAuthorizationREFG1 = append(l.PriorAuthorizationREFG1, ref)
		return
	case ReferenceIdentificationQualifierREF019F:
		l.ReferralREF9F = append(l.ReferralREF9F, ref)
		return
	case ReferenceIdentificationQualifierREF016R:
		if l.LineControlREF6R.QualifierREF01 != ReferenceIdentificationQualifierREF01NotSent || ref.Value == "" {
			break
		}
		l.LineControlREF6R = ref
		return
	case ReferenceIdentificationQualifierREF019B:
		slot = &l.RepricedLineItemREF9B
	case ReferenceIdentificationQualifierREF019D:
		slot = &l.AdjustedRepricedLineItemREF9D
	case ReferenceIdentificationQualifierREF01EW:
		slot = &l.MammographyCertificationREFEW
	case ReferenceIdentificationQualifierREF01X4:
		slot = &l.CLIAREFX4
	case ReferenceIdentificationQualifierREF01F4:
		slot = &l.ReferringCLIAREFF4
	case ReferenceIdentificationQualifierREF01BT:
		slot = &l.ImmunizationBatchREFBT
	}
	if slot == nil || slot.QualifierREF01 != ReferenceIdentificationQualifierREF01NotSent {
		l.Raw = append(l.Raw, rawSegment(s))
		return
	}
	*slot = ref
}

// drugIdentificationLoop reads the 2410 loop: the LIN, the first CTP, and
// the first REF that carries REF01. Anything else is kept raw.
func drugIdentificationLoop(loop []inputSegment, component byte) (DrugLIN, error) {
	d := parseLIN(loop[0])
	for _, s := range loop[1:] {
		switch {
		case s.id == "CTP" && !d.Pricing.Set:
			pricing, err := parseCTP(s, component)
			if err != nil {
				return DrugLIN{}, err
			}
			d.Pricing = pricing
		case s.id == "REF" && s.el(1) != "" && d.AssociationREFXZ.QualifierREF01 == ReferenceIdentificationQualifierREF01NotSent:
			d.AssociationREFXZ = parseREF(s)
		default:
			d.Raw = append(d.Raw, rawSegment(s))
		}
	}
	return d, nil
}

// lineAdjudicationLoop reads one 2430 loop: the SVD, every CAS, the first
// DTP*573, and the first AMT*EAF. Anything else is kept raw.
func lineAdjudicationLoop(loop []inputSegment, component byte) (LineAdjudicationSVD, error) {
	a, err := parseSVD(loop[0], component)
	if err != nil {
		return LineAdjudicationSVD{}, err
	}
	for _, s := range loop[1:] {
		var err error
		switch s.id {
		case "CAS":
			var adjustment AdjustmentCAS
			adjustment, err = parseCAS(s)
			a.Adjustments = append(a.Adjustments, adjustment)
		case "DTP":
			if DateTimeQualifierDTP01(s.el(1)) != DateTimeQualifierDTP01573 || a.AdjudicationDateDTP573.QualifierDTP01 != DateTimeQualifierDTP01NotSent {
				a.Raw = append(a.Raw, rawSegment(s))
				continue
			}
			a.AdjudicationDateDTP573, err = parseDTP(s)
		case "AMT":
			if AmountQualifierAMT01(s.el(1)) != AmountQualifierAMT01EAF || a.RemainingLiabilityAMTEAF.QualifierAMT01 != AmountQualifierAMT01NotSent {
				a.Raw = append(a.Raw, rawSegment(s))
				continue
			}
			a.RemainingLiabilityAMTEAF, err = parseAMT(s)
		default:
			a.Raw = append(a.Raw, rawSegment(s))
		}
		if err != nil {
			return LineAdjudicationSVD{}, err
		}
	}
	return a, nil
}

// supportingFormLoop reads one 2440 loop: the LQ and every FRM.
func supportingFormLoop(loop []inputSegment) (FormLQ, error) {
	f := parseLQ(loop[0])
	for _, s := range loop[1:] {
		if s.id != "FRM" {
			f.Raw = append(f.Raw, rawSegment(s))
			continue
		}
		answer, err := parseFRM(s)
		if err != nil {
			return FormLQ{}, err
		}
		f.Answers = append(f.Answers, answer)
	}
	return f, nil
}

// addAddress places an N3 in address when it is still empty; otherwise the
// N3 is kept raw.
func addAddress(s inputSegment, address *AddressN3, raw *[]RawSegment) {
	if address.Set {
		*raw = append(*raw, rawSegment(s))
		return
	}
	*address = parseN3(s)
}

// addLocation places an N4 in location when it is still empty; otherwise
// the N4 is kept raw.
func addLocation(s inputSegment, location *LocationN4, raw *[]RawSegment) {
	if location.Set {
		*raw = append(*raw, rawSegment(s))
		return
	}
	*location = parseN4(s)
}

// addDemographic places a DMG in demographics when it is still empty;
// otherwise the DMG is kept raw.
func addDemographic(s inputSegment, demographics *DemographicDMG, raw *[]RawSegment) error {
	if demographics.Set {
		*raw = append(*raw, rawSegment(s))
		return nil
	}
	parsed, err := parseDMG(s)
	if err != nil {
		return err
	}
	*demographics = parsed
	return nil
}

// addContact places a PER in contact when contact is still empty and PER01
// carries a value; otherwise the PER is kept raw, since a placed PER with no
// PER01 could not show that it was sent.
func addContact(s inputSegment, contact *ContactPER, raw *[]RawSegment) {
	if contact.FunctionPER01 != ContactFunctionCodePER01NotSent || s.el(1) == "" {
		*raw = append(*raw, rawSegment(s))
		return
	}
	*contact = parsePER(s)
}

// addSpecialty places a PRV in specialty when specialty is still empty and
// PRV01 carries a value; otherwise the PRV is kept raw, for the same reason
// as addContact.
func addSpecialty(s inputSegment, specialty *ProviderPRV, raw *[]RawSegment) {
	if specialty.CodePRV01 != ProviderTypePRV01NotSent || s.el(1) == "" {
		*raw = append(*raw, rawSegment(s))
		return
	}
	*specialty = parsePRV(s)
}
