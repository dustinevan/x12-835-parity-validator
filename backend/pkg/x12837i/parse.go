package x12837i

import (
	"errors"
	"fmt"
	"strconv"
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

// checkLast returns an error when the segment has an element after element
// last, the highest one X12 defines for this segment. splitSegments
// removes trailing empty elements, so any element past last carries text.
func (s inputSegment) checkLast(last int) error {
	if len(s.elements) > last {
		return fmt.Errorf("x12837i: parsing %s in segment %d: the model ends at %s", s.name(len(s.elements)), s.index, s.name(last))
	}
	return nil
}

func rawSegment(s inputSegment) RawSegment {
	return RawSegment{Index: s.index, ID: s.id, Elements: s.elements}
}

// splitInput is the input split into segments, with the delimiters and the
// layout the file uses.
type splitInput struct {
	segments     []inputSegment
	element      byte   // the element separator, the byte after "ISA"
	component    byte   // the component separator, ISA16
	terminator   byte   // the segment terminator, the byte after ISA16
	segmentBreak string // the whitespace after the ISA's terminator
}

// splitSegments reads the delimiters from the ISA and splits data into
// segments. Whitespace before a segment is skipped, a blank segment is
// skipped, trailing empty elements are removed (they carry no text), and
// when the terminator is a newline a carriage return before it is removed.
func splitSegments(data []byte) (splitInput, error) {
	text := string(data)
	start := 0
	for start < len(text) && isSpace(text[start]) {
		start++
	}
	if !strings.HasPrefix(text[start:], "ISA") || len(text) < start+4 {
		return splitInput{}, fmt.Errorf("x12837i: splitting input: expected ISA")
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
		return splitInput{}, fmt.Errorf("x12837i: splitting ISA: incomplete header")
	}
	component := text[isa16]
	terminator := text[isa16+1]
	isaFields := strings.Split(text[start:isa16+1], string(element))
	repetition := isaFields[11]
	if len(repetition) != 1 {
		return splitInput{}, fmt.Errorf("x12837i: splitting ISA: ISA11 is not one character")
	}
	if element == component || element == terminator || element == repetition[0] ||
		component == terminator || component == repetition[0] || terminator == repetition[0] {
		return splitInput{}, fmt.Errorf("x12837i: splitting ISA: the four separators must differ")
	}

	in := splitInput{element: element, component: component, terminator: terminator}
	in.segments = []inputSegment{{index: 0, id: "ISA", elements: isaFields[1:]}}
	pos := isa16 + 2
	for pos < len(text) && isSpace(text[pos]) {
		pos++
	}
	in.segmentBreak = text[isa16+2 : pos]
	for {
		for pos < len(text) && isSpace(text[pos]) {
			pos++
		}
		if pos == len(text) {
			return in, nil
		}
		if text[pos] == terminator {
			pos++
			continue
		}
		end := strings.IndexByte(text[pos:], terminator)
		if end < 0 {
			return splitInput{}, fmt.Errorf("x12837i: splitting segment %d: missing segment terminator", len(in.segments))
		}
		content := text[pos : pos+end]
		pos += end + 1
		if terminator == '\n' {
			content = strings.TrimSuffix(content, "\r")
		}
		fields := strings.Split(content, string(element))
		if !validID(fields[0]) {
			return splitInput{}, fmt.Errorf("x12837i: splitting segment %d: segment ID is not 2 or 3 capital letters or digits", len(in.segments))
		}
		in.segments = append(in.segments, inputSegment{index: len(in.segments), id: fields[0], elements: trimTrailingEmpty(fields[1:])})
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

// parser reads the segments once, building the tree. Each parse method
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

func (p *parser) parseInterchange() (InterchangeISA, error) {
	isa, err := parseISA(p.take()) // splitSegments always returns the ISA first
	if err != nil {
		return InterchangeISA{}, err
	}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return InterchangeISA{}, fmt.Errorf("x12837i: parsing interchange: missing IEA")
		case "GS":
			group, err := p.parseGroup()
			if err != nil {
				return InterchangeISA{}, err
			}
			isa.FunctionalGroups = append(isa.FunctionalGroups, group)
		case "IEA":
			iea := p.take()
			count, control, err := parseTrailer(iea)
			if err != nil {
				return InterchangeISA{}, err
			}
			isa.IndexIEA = iea.index
			isa.GroupCountIEA01 = count
			isa.ControlNumberIEA02 = control
			if p.i < len(p.segments) {
				return InterchangeISA{}, fmt.Errorf("x12837i: parsing segment %d: input continues after the IEA", p.peek().index)
			}
			return isa, nil
		case "ISA", "GE", "ST", "SE":
			return InterchangeISA{}, fmt.Errorf("x12837i: parsing segment %d: %s outside a functional group", next.index, next.id)
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
			return FunctionalGroupGS{}, fmt.Errorf("x12837i: parsing functional group at segment %d: missing GE", gs.index)
		case "ST":
			transaction, err := p.parseTransaction()
			if err != nil {
				return FunctionalGroupGS{}, err
			}
			group.Transactions = append(group.Transactions, transaction)
		case "GE":
			ge := p.take()
			count, control, err := parseTrailer(ge)
			if err != nil {
				return FunctionalGroupGS{}, err
			}
			group.IndexGE = ge.index
			group.TransactionCountGE01 = count
			group.ControlNumberGE02 = control
			return group, nil
		case "ISA", "IEA", "GS", "SE":
			return FunctionalGroupGS{}, fmt.Errorf("x12837i: parsing segment %d: %s inside the functional group at segment %d (missing GE)", next.index, next.id, gs.index)
		default:
			group.Raw = append(group.Raw, rawSegment(p.take()))
		}
	}
}

// hierarchy records where each placed HL loop of a transaction sits, keyed
// by its HL01, so that a later HL is placed under the loop its HL02 names.
type hierarchy struct {
	providers   map[string]int    // a 2000A loop's HL01: its index in BillingProviders
	subscribers map[string][2]int // a 2000B loop's HL01: its billing provider's index, then its own index in Subscribers
}

func (p *parser) parseTransaction() (TransactionST, error) {
	st := p.take()
	t, err := parseST(st)
	if err != nil {
		return TransactionST{}, err
	}
	h := hierarchy{providers: map[string]int{}, subscribers: map[string][2]int{}}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return TransactionST{}, fmt.Errorf("x12837i: parsing transaction at segment %d: missing SE", st.index)
		case "SE":
			se := p.take()
			count, control, err := parseTrailer(se)
			if err != nil {
				return TransactionST{}, err
			}
			t.IndexSE = se.index
			t.SegmentCountSE01 = count
			t.ControlNumberSE02 = control
			p.transactionLengths = append(p.transactionLengths, se.index-st.index+1)
			return t, nil
		case "ISA", "IEA", "GS", "GE", "ST":
			return TransactionST{}, fmt.Errorf("x12837i: parsing segment %d: %s inside the transaction at segment %d (missing SE)", next.index, next.id, st.index)
		case "NM1":
			if err := p.parseHeaderName(&t); err != nil {
				return TransactionST{}, err
			}
		case "HL":
			if err := p.parseHierarchy(&t, h); err != nil {
				return TransactionST{}, err
			}
		default:
			if err := t.add(p.take()); err != nil {
				return TransactionST{}, err
			}
		}
	}
}

// envelope reports whether a segment with this ID opens or closes an
// envelope. Every envelope segment ends every loop inside ST, so that the
// envelope level can close the loop or report the segment out of place.
func envelope(id string) bool {
	switch id {
	case "ISA", "IEA", "GS", "GE", "ST", "SE":
		return true
	}
	return false
}

// endsHierarchy reports whether a segment with this ID ends an HL loop;
// "" is the end of the input.
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

// inNameLoop reports whether a segment with this ID continues an NM1 loop.
func inNameLoop(id string) bool {
	switch id {
	case "N3", "N4", "REF", "PER", "PRV", "DMG", "DTP":
		return true
	}
	return false
}

// inDrugLoop reports whether a segment with this ID continues a LIN loop
// (2410).
func inDrugLoop(id string) bool {
	return id == "CTP" || id == "REF"
}

// inAdjudicationLoop reports whether a segment with this ID continues an SVD
// loop (2430).
func inAdjudicationLoop(id string) bool {
	return id == "CAS" || id == "DTP" || id == "AMT"
}

// takeNameLoop takes an NM1 and the segments of its loop, all raw.
func (p *parser) takeNameLoop() []RawSegment {
	raw := []RawSegment{rawSegment(p.take())}
	for inNameLoop(p.peek().id) {
		raw = append(raw, rawSegment(p.take()))
	}
	return raw
}

// takeHierarchyLoop takes an HL and every segment up to the next HL or
// envelope segment, all raw.
func (p *parser) takeHierarchyLoop() []RawSegment {
	raw := []RawSegment{rawSegment(p.take())}
	for !endsHierarchy(p.peek().id) {
		raw = append(raw, rawSegment(p.take()))
	}
	return raw
}

// placed reports whether a name loop's NM1 has been read: an NM1 always
// carries its role in NM101, so an empty role means the field is still
// open.
func placed(n NameNM1) bool {
	return n.EntityIDNM101 != ""
}

// parseHeaderName reads an NM1 loop directly under ST: the submitter
// (NM1*41, 1000A) with its PER contacts, or the receiver (NM1*40, 1000B). A
// second loop for either role, or any other role, is kept raw.
func (p *parser) parseHeaderName(t *TransactionST) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	switch {
	case role == PartyRoleNM10141 && !placed(t.SubmitterNM141.Name):
		name, err := parseNM1(p.take())
		if err != nil {
			return err
		}
		t.SubmitterNM141.Name = name
		for inNameLoop(p.peek().id) {
			s := p.take()
			if s.id != "PER" {
				t.SubmitterNM141.Raw = append(t.SubmitterNM141.Raw, rawSegment(s))
				continue
			}
			contact, err := parsePER(s)
			if err != nil {
				return err
			}
			t.SubmitterNM141.Contacts = append(t.SubmitterNM141.Contacts, contact)
		}
	case role == PartyRoleNM10140 && !placed(t.ReceiverNM140.Name):
		name, err := parseNM1(p.take())
		if err != nil {
			return err
		}
		t.ReceiverNM140.Name = name
		for inNameLoop(p.peek().id) {
			t.ReceiverNM140.Raw = append(t.ReceiverNM140.Raw, rawSegment(p.take()))
		}
	default:
		t.Raw = append(t.Raw, p.takeNameLoop()...)
	}
	return nil
}

// parseHierarchy reads one HL loop and places it by its level (HL03) and
// parent (HL02); see Parse. A loop with no place is kept raw on t.
func (p *parser) parseHierarchy(t *TransactionST, h hierarchy) error {
	hl := p.peek()
	level := HierarchicalLevelHL03(hl.el(3)) // an off-list level has no place and is kept raw below
	parent := hl.el(2)
	switch level {
	case HierarchicalLevelHL0320:
		if parent != "" {
			break
		}
		b, err := p.parseBillingProvider()
		if err != nil {
			return err
		}
		h.providers[b.Hierarchy.IDNumberHL01] = len(t.BillingProviders)
		t.BillingProviders = append(t.BillingProviders, b)
		return nil
	case HierarchicalLevelHL0322:
		i, ok := h.providers[parent]
		if !ok {
			break
		}
		s, err := p.parseSubscriber()
		if err != nil {
			return err
		}
		h.subscribers[s.Hierarchy.IDNumberHL01] = [2]int{i, len(t.BillingProviders[i].Subscribers)}
		t.BillingProviders[i].Subscribers = append(t.BillingProviders[i].Subscribers, s)
		return nil
	case HierarchicalLevelHL0323:
		at, ok := h.subscribers[parent]
		if !ok {
			break
		}
		patient, err := p.parsePatient()
		if err != nil {
			return err
		}
		subscriber := &t.BillingProviders[at[0]].Subscribers[at[1]]
		subscriber.Patients = append(subscriber.Patients, patient)
		return nil
	}
	t.Raw = append(t.Raw, p.takeHierarchyLoop()...)
	return nil
}

func (p *parser) parseBillingProvider() (BillingProviderHL, error) {
	hl, err := parseHL(p.take())
	if err != nil {
		return BillingProviderHL{}, err
	}
	b := BillingProviderHL{Hierarchy: hl}
	for !endsHierarchy(p.peek().id) {
		if p.peek().id == "NM1" {
			if err := p.parseBillingProviderName(&b); err != nil {
				return BillingProviderHL{}, err
			}
			continue
		}
		if err := b.add(p.take()); err != nil {
			return BillingProviderHL{}, err
		}
	}
	return b, nil
}

// parseBillingProviderName reads an NM1 loop in a 2000A loop: the billing
// provider (NM1*85, 2010AA) or the pay-to address (NM1*87, 2010AB), each
// the first time it appears. Any other loop, the 2010AC pay-to plan
// (NM1*PE) included, is kept raw.
func (p *parser) parseBillingProviderName(b *BillingProviderHL) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	billing := role == PartyRoleNM10185 && !placed(b.BillingProviderNM185)
	payTo := role == PartyRoleNM10187 && !placed(b.PayToAddressNM187)
	if !billing && !payTo {
		b.Raw = append(b.Raw, p.takeNameLoop()...)
		return nil
	}
	name, err := parseNM1(p.take())
	if err != nil {
		return err
	}
	if billing {
		b.BillingProviderNM185 = name
	} else {
		b.PayToAddressNM187 = name
	}
	for inNameLoop(p.peek().id) {
		s := p.take()
		var err error
		if billing {
			err = b.addBilling(s)
		} else {
			err = b.addPayTo(s)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseSubscriber() (SubscriberHL, error) {
	hl, err := parseHL(p.take())
	if err != nil {
		return SubscriberHL{}, err
	}
	sub := SubscriberHL{Hierarchy: hl}
	for !endsHierarchy(p.peek().id) {
		switch p.peek().id {
		case "CLM":
			claim, err := p.parseClaim()
			if err != nil {
				return SubscriberHL{}, err
			}
			sub.Claims = append(sub.Claims, claim)
		case "NM1":
			if err := p.parseSubscriberName(&sub); err != nil {
				return SubscriberHL{}, err
			}
		default:
			if err := sub.add(p.take()); err != nil {
				return SubscriberHL{}, err
			}
		}
	}
	return sub, nil
}

// parseSubscriberName reads an NM1 loop in a 2000B loop: the subscriber
// (NM1*IL, 2010BA) or the payer (NM1*PR, 2010BB), each the first time it
// appears. Any other loop is kept raw.
func (p *parser) parseSubscriberName(sub *SubscriberHL) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	subscriber := role == PartyRoleNM101IL && !placed(sub.SubscriberNM1IL)
	payer := role == PartyRoleNM101PR && !placed(sub.PayerNM1PR)
	if !subscriber && !payer {
		sub.Raw = append(sub.Raw, p.takeNameLoop()...)
		return nil
	}
	name, err := parseNM1(p.take())
	if err != nil {
		return err
	}
	if subscriber {
		sub.SubscriberNM1IL = name
	} else {
		sub.PayerNM1PR = name
	}
	for inNameLoop(p.peek().id) {
		s := p.take()
		var err error
		if subscriber {
			err = sub.addSubscriber(s)
		} else {
			err = sub.addPayer(s)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parsePatient() (PatientHL, error) {
	hl, err := parseHL(p.take())
	if err != nil {
		return PatientHL{}, err
	}
	pt := PatientHL{Hierarchy: hl}
	for !endsHierarchy(p.peek().id) {
		switch p.peek().id {
		case "CLM":
			claim, err := p.parseClaim()
			if err != nil {
				return PatientHL{}, err
			}
			pt.Claims = append(pt.Claims, claim)
		case "NM1":
			if err := p.parsePatientName(&pt); err != nil {
				return PatientHL{}, err
			}
		default:
			if err := pt.add(p.take()); err != nil {
				return PatientHL{}, err
			}
		}
	}
	return pt, nil
}

// parsePatientName reads an NM1 loop in a 2000C loop: the patient (NM1*QC,
// 2010CA) the first time it appears. Any other loop is kept raw.
func (p *parser) parsePatientName(pt *PatientHL) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	if role != PartyRoleNM101QC || placed(pt.PatientNM1QC) {
		pt.Raw = append(pt.Raw, p.takeNameLoop()...)
		return nil
	}
	name, err := parseNM1(p.take())
	if err != nil {
		return err
	}
	pt.PatientNM1QC = name
	for inNameLoop(p.peek().id) {
		if err := pt.addPatient(p.take()); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseClaim() (ClaimCLM, error) {
	c, err := parseCLM(p.take(), p.component)
	if err != nil {
		return ClaimCLM{}, err
	}
	for !endsClaim(p.peek().id) {
		switch p.peek().id {
		case "LX":
			line, err := p.parseServiceLine()
			if err != nil {
				return ClaimCLM{}, err
			}
			c.ServiceLines = append(c.ServiceLines, line)
		case "SBR":
			other, err := p.parseOtherSubscriber()
			if err != nil {
				return ClaimCLM{}, err
			}
			c.OtherSubscribers = append(c.OtherSubscribers, other)
		case "NM1":
			if err := p.parseClaimProvider(&c); err != nil {
				return ClaimCLM{}, err
			}
		default:
			if err := c.add(p.take(), p.component); err != nil {
				return ClaimCLM{}, err
			}
		}
	}
	return c, nil
}

// providerLoop is one NM1-led provider loop read into the parts the
// provider loop types hold: the name, its PRV, its REF segments, and the
// segments it does not place.
type providerLoop struct {
	name       NameNM1
	specialty  ProviderPRV
	references []ReferenceREF
	raw        []RawSegment
}

// parseProviderLoop reads an NM1 and the segments of its loop: every REF,
// and the first PRV when withPRV is true. Anything else is kept raw.
func (p *parser) parseProviderLoop(withPRV bool) (providerLoop, error) {
	name, err := parseNM1(p.take())
	if err != nil {
		return providerLoop{}, err
	}
	l := providerLoop{name: name}
	for inNameLoop(p.peek().id) {
		s := p.take()
		switch {
		case s.id == "REF":
			ref, err := parseREF(s)
			if err != nil {
				return providerLoop{}, err
			}
			l.references = append(l.references, ref)
		case s.id == "PRV" && withPRV && !l.specialty.Set:
			prv, err := parsePRV(s)
			if err != nil {
				return providerLoop{}, err
			}
			l.specialty = prv
		default:
			l.raw = append(l.raw, rawSegment(s))
		}
	}
	return l, nil
}

// parseClaimProvider reads an NM1 loop in a claim: 2310A-F by NM101 (71,
// 72, ZZ, 82, 77, DN), each the first time it appears. Any other loop is
// kept raw.
func (p *parser) parseClaimProvider(c *ClaimCLM) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	open := false
	switch role {
	case PartyRoleNM10171:
		open = !placed(c.AttendingNM171.Name)
	case PartyRoleNM10172:
		open = !placed(c.OperatingNM172.Name)
	case PartyRoleNM101ZZ:
		open = !placed(c.OtherOperatingNM1ZZ.Name)
	case PartyRoleNM10182:
		open = !placed(c.RenderingNM182.Name)
	case PartyRoleNM10177:
		open = !placed(c.ServiceFacilityNM177.Name)
	case PartyRoleNM101DN:
		open = !placed(c.ReferringNM1DN.Name)
	}
	if !open {
		c.Raw = append(c.Raw, p.takeNameLoop()...)
		return nil
	}
	l, err := p.parseProviderLoop(role == PartyRoleNM10171)
	if err != nil {
		return err
	}
	switch role {
	case PartyRoleNM10171:
		c.AttendingNM171 = AttendingNM171{Name: l.name, Specialty: l.specialty, References: l.references, Raw: l.raw}
	case PartyRoleNM10172:
		c.OperatingNM172 = OperatingNM172{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM101ZZ:
		c.OtherOperatingNM1ZZ = OtherOperatingNM1ZZ{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10182:
		c.RenderingNM182 = RenderingNM182{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10177:
		c.ServiceFacilityNM177 = ServiceFacilityNM177{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM101DN:
		c.ReferringNM1DN = ReferringNM1DN{Name: l.name, References: l.references, Raw: l.raw}
	}
	return nil
}

func (p *parser) parseOtherSubscriber() (OtherSubscriberSBR, error) {
	sbr, err := parseSBR(p.take())
	if err != nil {
		return OtherSubscriberSBR{}, err
	}
	o := OtherSubscriberSBR{SBR: sbr}
	for !endsOtherSubscriber(p.peek().id) {
		if p.peek().id == "NM1" {
			if err := p.parseOtherPayerName(&o); err != nil {
				return OtherSubscriberSBR{}, err
			}
			continue
		}
		if err := o.add(p.take()); err != nil {
			return OtherSubscriberSBR{}, err
		}
	}
	return o, nil
}

// parseOtherPayerName reads an NM1 loop in a 2320 loop: 2330A-I by NM101
// (IL, PR, 71, 72, ZZ, 77, 82, DN, 85), each the first time it appears. Any
// other loop is kept raw.
func (p *parser) parseOtherPayerName(o *OtherSubscriberSBR) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	open := false
	switch role {
	case PartyRoleNM101IL:
		open = !placed(o.OtherInsuredNM1IL.Name)
	case PartyRoleNM101PR:
		open = !placed(o.OtherPayerNM1PR.Name)
	case PartyRoleNM10171:
		open = !placed(o.OtherPayerAttendingNM171.Name)
	case PartyRoleNM10172:
		open = !placed(o.OtherPayerOperatingNM172.Name)
	case PartyRoleNM101ZZ:
		open = !placed(o.OtherPayerOtherOperatingNM1ZZ.Name)
	case PartyRoleNM10177:
		open = !placed(o.OtherPayerServiceFacilityNM177.Name)
	case PartyRoleNM10182:
		open = !placed(o.OtherPayerRenderingNM182.Name)
	case PartyRoleNM101DN:
		open = !placed(o.OtherPayerReferringNM1DN.Name)
	case PartyRoleNM10185:
		open = !placed(o.OtherPayerBillingNM185.Name)
	}
	if !open {
		o.Raw = append(o.Raw, p.takeNameLoop()...)
		return nil
	}
	if role == PartyRoleNM101PR {
		name, err := parseNM1(p.take())
		if err != nil {
			return err
		}
		o.OtherPayerNM1PR = OtherPayerNM1PR{Name: name}
		for inNameLoop(p.peek().id) {
			if err := o.OtherPayerNM1PR.add(p.take()); err != nil {
				return err
			}
		}
		return nil
	}
	l, err := p.parseProviderLoop(false)
	if err != nil {
		return err
	}
	switch role {
	case PartyRoleNM101IL:
		o.OtherInsuredNM1IL = OtherInsuredNM1IL{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10171:
		o.OtherPayerAttendingNM171 = OtherPayerAttendingNM171{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10172:
		o.OtherPayerOperatingNM172 = OperatingNM172{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM101ZZ:
		o.OtherPayerOtherOperatingNM1ZZ = OtherOperatingNM1ZZ{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10177:
		o.OtherPayerServiceFacilityNM177 = ServiceFacilityNM177{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10182:
		o.OtherPayerRenderingNM182 = RenderingNM182{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM101DN:
		o.OtherPayerReferringNM1DN = ReferringNM1DN{Name: l.name, References: l.references, Raw: l.raw}
	case PartyRoleNM10185:
		o.OtherPayerBillingNM185 = OtherPayerBillingNM185{Name: l.name, References: l.references, Raw: l.raw}
	}
	return nil
}

func (p *parser) parseServiceLine() (ServiceLineLX, error) {
	s := p.take()
	if err := s.checkLast(1); err != nil {
		return ServiceLineLX{}, err
	}
	if _, err := parseWholeNumber(s, 1); err != nil {
		return ServiceLineLX{}, err
	}
	l := ServiceLineLX{Index: s.index, NumberLX01: s.el(1)}
	for !endsServiceLine(p.peek().id) {
		switch p.peek().id {
		case "LIN":
			drug, err := p.parseDrug()
			if err != nil {
				return ServiceLineLX{}, err
			}
			l.Drugs = append(l.Drugs, drug)
		case "NM1":
			if err := p.parseLineProvider(&l); err != nil {
				return ServiceLineLX{}, err
			}
		case "SVD":
			adjudication, err := p.parseAdjudication()
			if err != nil {
				return ServiceLineLX{}, err
			}
			l.Adjudications = append(l.Adjudications, adjudication)
		default:
			if err := l.add(p.take(), p.component); err != nil {
				return ServiceLineLX{}, err
			}
		}
	}
	return l, nil
}

func (p *parser) parseDrug() (DrugLIN, error) {
	drug, err := parseLIN(p.take())
	if err != nil {
		return DrugLIN{}, err
	}
	for inDrugLoop(p.peek().id) {
		if err := drug.add(p.take(), p.component); err != nil {
			return DrugLIN{}, err
		}
	}
	return drug, nil
}

// parseLineProvider reads an NM1 loop in a service line: 2420A-D by NM101
// (72, ZZ, 82, DN), each the first time it appears. Any other loop is kept
// raw.
func (p *parser) parseLineProvider(l *ServiceLineLX) error {
	role := PartyRoleNM101(p.peek().el(1)) // an off-list role has no field and is kept raw below
	open := false
	switch role {
	case PartyRoleNM10172:
		open = !placed(l.OperatingNM172.Name)
	case PartyRoleNM101ZZ:
		open = !placed(l.OtherOperatingNM1ZZ.Name)
	case PartyRoleNM10182:
		open = !placed(l.RenderingNM182.Name)
	case PartyRoleNM101DN:
		open = !placed(l.ReferringNM1DN.Name)
	}
	if !open {
		l.Raw = append(l.Raw, p.takeNameLoop()...)
		return nil
	}
	loop, err := p.parseProviderLoop(false)
	if err != nil {
		return err
	}
	provider := LineProviderNM1{Name: loop.name, References: loop.references, Raw: loop.raw}
	switch role {
	case PartyRoleNM10172:
		l.OperatingNM172 = provider
	case PartyRoleNM101ZZ:
		l.OtherOperatingNM1ZZ = provider
	case PartyRoleNM10182:
		l.RenderingNM182 = provider
	case PartyRoleNM101DN:
		l.ReferringNM1DN = provider
	}
	return nil
}

func (p *parser) parseAdjudication() (LineAdjudicationSVD, error) {
	a, err := parseSVD(p.take(), p.component)
	if err != nil {
		return LineAdjudicationSVD{}, err
	}
	for inAdjudicationLoop(p.peek().id) {
		if err := a.add(p.take()); err != nil {
			return LineAdjudicationSVD{}, err
		}
	}
	return a, nil
}

// add places one segment that sits directly under ST, outside any loop:
// the first BHT. Anything else is kept raw.
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

// add places one segment of a 2000A loop outside its NM1 loops: the first
// PRV and CUR. Anything else is kept raw.
func (b *BillingProviderHL) add(s inputSegment) error {
	switch s.id {
	case "PRV":
		if b.Provider.Set {
			break
		}
		prv, err := parsePRV(s)
		if err != nil {
			return err
		}
		b.Provider = prv
		return nil
	case "CUR":
		if b.Currency.Set {
			break
		}
		currency, err := parseCUR(s)
		if err != nil {
			return err
		}
		b.Currency = currency
		return nil
	}
	b.Raw = append(b.Raw, rawSegment(s))
	return nil
}

// addBilling places one segment of the 2010AA loop after its NM1: the first
// N3 and N4, every REF (the first REF*EI that carries a value to
// TaxIDREFEI), and every PER. Anything else is kept raw.
func (b *BillingProviderHL) addBilling(s inputSegment) error {
	switch s.id {
	case "N3":
		if b.Address.Set {
			break
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		b.Address = address
		return nil
	case "N4":
		if b.Location.Set {
			break
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		b.Location = location
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		if ref.QualifierREF01 == "EI" && b.TaxIDREFEI.QualifierREF01 == "" && ref.Value != "" {
			b.TaxIDREFEI = ref
			return nil
		}
		b.References = append(b.References, ref)
		return nil
	case "PER":
		contact, err := parsePER(s)
		if err != nil {
			return err
		}
		b.Contacts = append(b.Contacts, contact)
		return nil
	}
	b.Raw = append(b.Raw, rawSegment(s))
	return nil
}

// addPayTo places one segment of the 2010AB loop after its NM1: the first
// N3 and N4. Anything else is kept raw.
func (b *BillingProviderHL) addPayTo(s inputSegment) error {
	switch s.id {
	case "N3":
		if b.PayToAddress.Set {
			break
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		b.PayToAddress = address
		return nil
	case "N4":
		if b.PayToLocation.Set {
			break
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		b.PayToLocation = location
		return nil
	}
	b.Raw = append(b.Raw, rawSegment(s))
	return nil
}

// add places one segment of a 2000B loop outside its NM1 loops and claims:
// the first SBR. Anything else is kept raw.
func (sub *SubscriberHL) add(s inputSegment) error {
	if s.id != "SBR" || sub.Subscriber.Set {
		sub.Raw = append(sub.Raw, rawSegment(s))
		return nil
	}
	sbr, err := parseSBR(s)
	if err != nil {
		return err
	}
	sub.Subscriber = sbr
	return nil
}

// addSubscriber places one segment of the 2010BA loop after its NM1: the
// first N3, N4, and DMG (a DMG only when DMG01 is D8 or not sent, since
// DemographicDMG holds a single date), and every REF (the first REF*SY and
// REF*Y4 that carry a value to their fields). Anything else is kept raw.
func (sub *SubscriberHL) addSubscriber(s inputSegment) error {
	switch s.id {
	case "N3":
		if sub.Address.Set {
			break
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		sub.Address = address
		return nil
	case "N4":
		if sub.Location.Set {
			break
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		sub.Location = location
		return nil
	case "DMG":
		if sub.Demographics.Set || !singleDateFormat(s.el(1)) {
			break
		}
		dmg, err := parseDMG(s)
		if err != nil {
			return err
		}
		sub.Demographics = dmg
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		switch {
		case ref.QualifierREF01 == "SY" && sub.SecondaryIDREFSY.QualifierREF01 == "" && ref.Value != "":
			sub.SecondaryIDREFSY = ref
		case ref.QualifierREF01 == "Y4" && sub.ClaimNumberREFY4.QualifierREF01 == "" && ref.Value != "":
			sub.ClaimNumberREFY4 = ref
		default:
			sub.References = append(sub.References, ref)
		}
		return nil
	}
	sub.Raw = append(sub.Raw, rawSegment(s))
	return nil
}

// addPayer places one segment of the 2010BB loop after its NM1: the first
// N3 and N4, and every REF. Anything else is kept raw.
func (sub *SubscriberHL) addPayer(s inputSegment) error {
	switch s.id {
	case "N3":
		if sub.PayerAddress.Set {
			break
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		sub.PayerAddress = address
		return nil
	case "N4":
		if sub.PayerLocation.Set {
			break
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		sub.PayerLocation = location
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		sub.PayerReferences = append(sub.PayerReferences, ref)
		return nil
	}
	sub.Raw = append(sub.Raw, rawSegment(s))
	return nil
}

// add places one segment of a 2000C loop outside its NM1 loops and claims:
// the first PAT. Anything else is kept raw.
func (pt *PatientHL) add(s inputSegment) error {
	if s.id != "PAT" || pt.Patient.Set {
		pt.Raw = append(pt.Raw, rawSegment(s))
		return nil
	}
	pat, err := parsePAT(s)
	if err != nil {
		return err
	}
	pt.Patient = pat
	return nil
}

// addPatient places one segment of the 2010CA loop after its NM1: the first
// N3, N4, and DMG (as addSubscriber places them), and every REF. Anything
// else is kept raw.
func (pt *PatientHL) addPatient(s inputSegment) error {
	switch s.id {
	case "N3":
		if pt.Address.Set {
			break
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		pt.Address = address
		return nil
	case "N4":
		if pt.Location.Set {
			break
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		pt.Location = location
		return nil
	case "DMG":
		if pt.Demographics.Set || !singleDateFormat(s.el(1)) {
			break
		}
		dmg, err := parseDMG(s)
		if err != nil {
			return err
		}
		pt.Demographics = dmg
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		pt.References = append(pt.References, ref)
		return nil
	}
	pt.Raw = append(pt.Raw, rawSegment(s))
	return nil
}

// singleDateFormat reports whether a date format qualifier (DMG01) names
// the single CCYYMMDD date the model holds, or was not sent.
func singleDateFormat(format string) bool {
	return format == "" || format == "D8"
}

// add places one segment of a claim outside its sub-loops: every PWK, CN1,
// K3, NTE, CRC, and HCP; the first CL1; the DTP, AMT, REF, and HI segments
// as addDate, addAmount, addReference, and addHealthInfo describe.
// Anything else is kept raw.
func (c *ClaimCLM) add(s inputSegment, component byte) error {
	switch s.id {
	case "DTP":
		return c.addDate(s)
	case "CL1":
		if c.InstitutionalCL1.Set {
			break
		}
		cl1, err := parseCL1(s)
		if err != nil {
			return err
		}
		c.InstitutionalCL1 = cl1
		return nil
	case "PWK":
		pwk, err := parsePWK(s)
		if err != nil {
			return err
		}
		c.SupplementalInfo = append(c.SupplementalInfo, pwk)
		return nil
	case "CN1":
		contract, err := parseCN1(s)
		if err != nil {
			return err
		}
		c.Contracts = append(c.Contracts, contract)
		return nil
	case "AMT":
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		c.addAmount(amount)
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		c.addReference(ref)
		return nil
	case "K3":
		file, err := parseK3(s)
		if err != nil {
			return err
		}
		c.Files = append(c.Files, file)
		return nil
	case "NTE":
		note, err := parseNTE(s)
		if err != nil {
			return err
		}
		c.Notes = append(c.Notes, note)
		return nil
	case "CRC":
		referral, err := parseCRC(s)
		if err != nil {
			return err
		}
		c.EPSDTReferralsCRCZZ = append(c.EPSDTReferralsCRCZZ, referral)
		return nil
	case "HI":
		return c.addHealthInfo(s, component)
	case "HCP":
		pricing, err := parseHCP(s)
		if err != nil {
			return err
		}
		c.Pricing = append(c.Pricing, pricing)
		return nil
	}
	c.Raw = append(c.Raw, rawSegment(s))
	return nil
}

// addDate places a claim DTP in the field for its qualifier (DTP01) when
// that field is still empty, DTP03 carries a value, and DTP02 names a format
// the field holds: 096 TM, 434 RD8 or D8, 435 DT or D8, and 050 D8
// (Transmittal 907, PDF p.24). Otherwise the DTP is kept raw.
func (c *ClaimCLM) addDate(s inputSegment) error {
	qualifier := DateQualifierDTP01(s.el(1)) // a qualifier with no claim field is kept raw below
	format := s.el(2)
	if s.el(3) == "" {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	d := DateDTP{Set: true, Index: s.index, QualifierDTP01: qualifier, FormatDTP02: format}
	switch {
	case qualifier == DateQualifierDTP01096 && format == "TM" && !c.DischargeHourDTP096.Set:
		if err := s.checkLast(3); err != nil {
			return err
		}
		hour, err := parseTime(s, 3)
		if err != nil {
			return err
		}
		d.Time = hour
		c.DischargeHourDTP096 = d
	case qualifier == DateQualifierDTP01434 && (format == "RD8" || format == "D8") && !c.StatementDatesDTP434.Set:
		from, to, err := parseDTPPeriod(s)
		if err != nil {
			return err
		}
		d.Date, d.DateThrough = from, to
		c.StatementDatesDTP434 = d
	case qualifier == DateQualifierDTP01435 && (format == "DT" || format == "D8") && !c.AdmissionDTP435.Set:
		if err := s.checkLast(3); err != nil {
			return err
		}
		date, hour, err := parseDateHour(s, 3, format)
		if err != nil {
			return err
		}
		d.Date, d.Time = date, hour
		c.AdmissionDTP435 = d
	case qualifier == DateQualifierDTP01050 && format == "D8" && !c.RepricerReceivedDTP050.Set:
		date, _, err := parseDTPPeriod(s)
		if err != nil {
			return err
		}
		d.Date = date
		c.RepricerReceivedDTP050 = d
	default:
		c.Raw = append(c.Raw, rawSegment(s))
	}
	return nil
}

// addAmount places a claim AMT: the first AMT*F3 that carries an amount to
// PatientEstimatedAmountAMTF3, and every other AMT to Amounts.
func (c *ClaimCLM) addAmount(a AmountAMT) {
	if a.QualifierAMT01 == "F3" && c.PatientEstimatedAmountAMTF3.QualifierAMT01 == "" && a.Amount.Set {
		c.PatientEstimatedAmountAMTF3 = a
		return
	}
	c.Amounts = append(c.Amounts, a)
}

// addReference places a claim REF: the first REF with a value for each
// qualifier REF_CLAIM_NAME names goes to its field, and every other REF to
// References.
func (c *ClaimCLM) addReference(ref ReferenceREF) {
	var field *ReferenceREF
	switch ref.QualifierREF01 {
	case "4N":
		field = &c.ServiceAuthExceptionREF4N
	case "9F":
		field = &c.ReferralNumberREF9F
	case "G1":
		field = &c.PriorAuthREFG1
	case "F8":
		field = &c.PayerClaimControlREFF8
	case "9A":
		field = &c.RepricedClaimREF9A
	case "9C":
		field = &c.AdjustedRepricedClaimREF9C
	case "LX":
		field = &c.InvestigationalDeviceREFLX
	case "D9":
		field = &c.TransmissionIntermediaryREFD9
	case "LU":
		field = &c.AutoAccidentStateREFLU
	case "EA":
		field = &c.MedicalRecordREFEA
	case "P4":
		field = &c.DemonstrationProjectREFP4
	case "G4":
		field = &c.PROApprovalREFG4
	}
	if field == nil || field.QualifierREF01 != "" || ref.Value == "" {
		c.References = append(c.References, ref)
		return
	}
	*field = ref
}

// addHealthInfo places the composites of one HI segment in the fields for
// their qualifiers (HI0n-1). The singular fields (principal diagnosis,
// admitting diagnosis, principal procedure, DRG) take one composite; the
// others take every composite. A segment is placed whole or kept raw whole:
// when any composite's qualifier is off the list or names a singular field
// already filled, or when the segment sends no composite at all, the
// segment is kept raw.
func (c *ClaimCLM) addHealthInfo(s inputSegment, component byte) error {
	qualifiers := hiQualifiers(s, component)
	taken := len(qualifiers) == 0
	principal, admitting, procedure, drg := c.PrincipalDiagnosisHIABK.Set, c.AdmittingDiagnosisHIABJ.Set, c.PrincipalProcedureHIBBR.Set, c.DRGHIDR.Set
	for _, qualifier := range qualifiers {
		switch qualifier {
		case DiagnosisQualifierHI01_1ABK, DiagnosisQualifierHI01_1BK:
			taken, principal = taken || principal, true
		case DiagnosisQualifierHI01_1ABJ, DiagnosisQualifierHI01_1BJ:
			taken, admitting = taken || admitting, true
		case DiagnosisQualifierHI01_1BBR, DiagnosisQualifierHI01_1BR:
			taken, procedure = taken || procedure, true
		case DiagnosisQualifierHI01_1DR:
			taken, drg = taken || drg, true
		default:
			if !qualifier.Known() {
				taken = true
			}
		}
	}
	if taken {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	infos, err := parseHI(s, component)
	if err != nil {
		return err
	}
	for _, hi := range infos {
		switch hi.QualifierHI01_1 {
		case DiagnosisQualifierHI01_1ABK, DiagnosisQualifierHI01_1BK:
			c.PrincipalDiagnosisHIABK = hi
		case DiagnosisQualifierHI01_1ABJ, DiagnosisQualifierHI01_1BJ:
			c.AdmittingDiagnosisHIABJ = hi
		case DiagnosisQualifierHI01_1APR, DiagnosisQualifierHI01_1PR:
			c.PatientReasonsHIAPR = append(c.PatientReasonsHIAPR, hi)
		case DiagnosisQualifierHI01_1ABN, DiagnosisQualifierHI01_1BN:
			c.ExternalCausesHIABN = append(c.ExternalCausesHIABN, hi)
		case DiagnosisQualifierHI01_1ABF, DiagnosisQualifierHI01_1BF:
			c.OtherDiagnosesHIABF = append(c.OtherDiagnosesHIABF, hi)
		case DiagnosisQualifierHI01_1BBR, DiagnosisQualifierHI01_1BR:
			c.PrincipalProcedureHIBBR = hi
		case DiagnosisQualifierHI01_1BBQ, DiagnosisQualifierHI01_1BQ:
			c.OtherProceduresHIBBQ = append(c.OtherProceduresHIBBQ, hi)
		case DiagnosisQualifierHI01_1BI:
			c.OccurrenceSpansHIBI = append(c.OccurrenceSpansHIBI, hi)
		case DiagnosisQualifierHI01_1BH:
			c.OccurrenceCodesHIBH = append(c.OccurrenceCodesHIBH, hi)
		case DiagnosisQualifierHI01_1BE:
			c.ValueCodesHIBE = append(c.ValueCodesHIBE, hi)
		case DiagnosisQualifierHI01_1BG:
			c.ConditionCodesHIBG = append(c.ConditionCodesHIBG, hi)
		case DiagnosisQualifierHI01_1TC:
			c.TreatmentCodesHITC = append(c.TreatmentCodesHITC, hi)
		case DiagnosisQualifierHI01_1DR:
			c.DRGHIDR = hi
		}
	}
	return nil
}

// add places one segment of a 2320 loop outside its NM1 loops: every CAS;
// the first OI, MIA, and MOA; and every AMT, the first AMT*D, AMT*EAF, and
// AMT*A8 that carries an amount going to their fields. Anything else is kept
// raw.
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
		a, err := parseAMT(s)
		if err != nil {
			return err
		}
		switch {
		case a.QualifierAMT01 == "D" && o.OtherPayerPaidAMTD.QualifierAMT01 == "" && a.Amount.Set:
			o.OtherPayerPaidAMTD = a
		case a.QualifierAMT01 == "EAF" && o.RemainingLiabilityAMTEAF.QualifierAMT01 == "" && a.Amount.Set:
			o.RemainingLiabilityAMTEAF = a
		case a.QualifierAMT01 == "A8" && o.NonCoveredAMTA8.QualifierAMT01 == "" && a.Amount.Set:
			o.NonCoveredAMTA8 = a
		default:
			o.Amounts = append(o.Amounts, a)
		}
		return nil
	case "OI":
		if o.OtherInsurance.Set {
			break
		}
		oi, err := parseOI(s)
		if err != nil {
			return err
		}
		o.OtherInsurance = oi
		return nil
	case "MIA":
		if o.MIA.Set {
			break
		}
		mia, err := parseMIA(s)
		if err != nil {
			return err
		}
		o.MIA = mia
		return nil
	case "MOA":
		if o.MOA.Set {
			break
		}
		moa, err := parseMOA(s)
		if err != nil {
			return err
		}
		o.MOA = moa
		return nil
	}
	o.Raw = append(o.Raw, rawSegment(s))
	return nil
}

// add places one segment of the 2330B loop after its NM1: the first DTP*573
// sent in D8 format, and every REF (the first REF with a value for each
// qualifier REF_OTHER_PAYER_NAME names going to its field). Anything else is
// kept raw.
func (o *OtherPayerNM1PR) add(s inputSegment) error {
	switch s.id {
	case "DTP":
		qualifier := DateQualifierDTP01(s.el(1)) // a qualifier with no field here is kept raw below
		if qualifier != DateQualifierDTP01573 || s.el(2) != "D8" || s.el(3) == "" || o.CheckOrRemittanceDTP573.Set {
			break
		}
		date, err := parseDTP(s)
		if err != nil {
			return err
		}
		o.CheckOrRemittanceDTP573 = date
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		var field *ReferenceREF
		switch ref.QualifierREF01 {
		case "2U":
			field = &o.SecondaryIDREF2U
		case "EI":
			field = &o.SecondaryIDREFEI
		case "FY":
			field = &o.SecondaryIDREFFY
		case "NF":
			field = &o.SecondaryIDREFNF
		case "G1":
			field = &o.PriorAuthREFG1
		case "9F":
			field = &o.ReferralNumberREF9F
		case "T4":
			field = &o.AdjustmentIndicatorREFT4
		case "F8":
			field = &o.ClaimControlREFF8
		}
		if field == nil || field.QualifierREF01 != "" || ref.Value == "" {
			o.References = append(o.References, ref)
			return nil
		}
		*field = ref
		return nil
	}
	o.Raw = append(o.Raw, rawSegment(s))
	return nil
}

// add places one segment of a service line outside its sub-loops: the first
// SV2 and HCP; every PWK; every DTP (the first DTP*472 in D8 or RD8 format
// to ServiceDateDTP472, any other D8 or RD8 DTP to Dates); every REF (the
// first REF*6R, REF*9B, and REF*9D with a value to their fields); and every
// AMT (the first AMT*GT and AMT*N8 that carry an amount to their fields).
// Anything else is kept raw, the line-level NTE included.
func (l *ServiceLineLX) add(s inputSegment, component byte) error {
	switch s.id {
	case "SV2":
		if l.Service.Set {
			break
		}
		sv2, err := parseSV2(s, component)
		if err != nil {
			return err
		}
		l.Service = sv2
		return nil
	case "PWK":
		pwk, err := parsePWK(s)
		if err != nil {
			return err
		}
		l.SupplementalInfo = append(l.SupplementalInfo, pwk)
		return nil
	case "DTP":
		if !periodFormat(s.el(2)) || s.el(3) == "" {
			break
		}
		date, err := parseDTP(s)
		if err != nil {
			return err
		}
		if date.QualifierDTP01 == DateQualifierDTP01472 && !l.ServiceDateDTP472.Set {
			l.ServiceDateDTP472 = date
			return nil
		}
		l.Dates = append(l.Dates, date)
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		var field *ReferenceREF
		switch ref.QualifierREF01 {
		case "6R":
			field = &l.LineControlREF6R
		case "9B":
			field = &l.RepricedLineREF9B
		case "9D":
			field = &l.AdjustedRepricedLineREF9D
		}
		if field == nil || field.QualifierREF01 != "" || ref.Value == "" {
			l.References = append(l.References, ref)
			return nil
		}
		*field = ref
		return nil
	case "AMT":
		a, err := parseAMT(s)
		if err != nil {
			return err
		}
		switch {
		case a.QualifierAMT01 == "GT" && l.ServiceTaxAMTGT.QualifierAMT01 == "" && a.Amount.Set:
			l.ServiceTaxAMTGT = a
		case a.QualifierAMT01 == "N8" && l.FacilityTaxAMTN8.QualifierAMT01 == "" && a.Amount.Set:
			l.FacilityTaxAMTN8 = a
		default:
			l.Amounts = append(l.Amounts, a)
		}
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

// periodFormat reports whether a date format qualifier (DTP02) names a
// format DateDTP holds as a date or a range: D8, one date, or RD8, a range.
func periodFormat(format string) bool {
	return format == "D8" || format == "RD8"
}

// add places one segment of a 2410 loop after its LIN: the first CTP, and
// every REF (the first REF*VY and REF*XZ with a value to their fields).
// Anything else is kept raw.
func (d *DrugLIN) add(s inputSegment, component byte) error {
	switch s.id {
	case "CTP":
		if d.Quantity.Set {
			break
		}
		ctp, err := parseCTP(s, component)
		if err != nil {
			return err
		}
		d.Quantity = ctp
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		switch {
		case ref.QualifierREF01 == "VY" && d.LinkSequenceREFVY.QualifierREF01 == "" && ref.Value != "":
			d.LinkSequenceREFVY = ref
		case ref.QualifierREF01 == "XZ" && d.PrescriptionREFXZ.QualifierREF01 == "" && ref.Value != "":
			d.PrescriptionREFXZ = ref
		default:
			d.References = append(d.References, ref)
		}
		return nil
	}
	d.Raw = append(d.Raw, rawSegment(s))
	return nil
}

// add places one segment of a 2430 loop after its SVD: every CAS; every DTP
// in D8 or RD8 format (the first DTP*573 in D8 format to its field, any
// other to Dates); and every AMT (the first AMT*EAF that carries an amount
// to its field). Anything else is kept raw.
func (a *LineAdjudicationSVD) add(s inputSegment) error {
	switch s.id {
	case "CAS":
		adjustment, err := parseCAS(s)
		if err != nil {
			return err
		}
		a.Adjustments = append(a.Adjustments, adjustment)
		return nil
	case "DTP":
		if !periodFormat(s.el(2)) || s.el(3) == "" {
			break
		}
		date, err := parseDTP(s)
		if err != nil {
			return err
		}
		if date.QualifierDTP01 == DateQualifierDTP01573 && date.FormatDTP02 == "D8" && !a.CheckOrRemittanceDateDTP573.Set {
			a.CheckOrRemittanceDateDTP573 = date
			return nil
		}
		a.Dates = append(a.Dates, date)
		return nil
	case "AMT":
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		if amount.QualifierAMT01 == "EAF" && a.RemainingLiabilityAMTEAF.QualifierAMT01 == "" && amount.Amount.Set {
			a.RemainingLiabilityAMTEAF = amount
			return nil
		}
		a.Amounts = append(a.Amounts, amount)
		return nil
	}
	a.Raw = append(a.Raw, rawSegment(s))
	return nil
}

// checkEnvelopes compares each envelope trailer with its header and with
// what the pass found, and returns every mismatch. transactionLengths holds,
// in file order, the number of segments from each ST through its SE. The
// counts are the digits the file sent; parseTrailer has already checked
// that each is a whole number.
func checkEnvelopes(isa InterchangeISA, transactionLengths []int) error {
	var errs []error
	if isa.ControlNumberISA13 != isa.ControlNumberIEA02 {
		errs = append(errs, fmt.Errorf("x12837i: checking IEA02: ISA13 is %q, IEA02 is %q", isa.ControlNumberISA13, isa.ControlNumberIEA02))
	}
	if count(isa.GroupCountIEA01) != len(isa.FunctionalGroups) {
		errs = append(errs, fmt.Errorf("x12837i: checking IEA01: IEA01 is %s, the interchange holds %d functional groups", isa.GroupCountIEA01, len(isa.FunctionalGroups)))
	}
	k := 0
	for _, g := range isa.FunctionalGroups {
		if g.ControlNumberGS06 != g.ControlNumberGE02 {
			errs = append(errs, fmt.Errorf("x12837i: checking GE02: GS06 is %q, GE02 is %q", g.ControlNumberGS06, g.ControlNumberGE02))
		}
		if count(g.TransactionCountGE01) != len(g.Transactions) {
			errs = append(errs, fmt.Errorf("x12837i: checking GE01 of group %q: GE01 is %s, the group holds %d transactions", g.ControlNumberGS06, g.TransactionCountGE01, len(g.Transactions)))
		}
		for _, t := range g.Transactions {
			if t.ControlNumberST02 != t.ControlNumberSE02 {
				errs = append(errs, fmt.Errorf("x12837i: checking SE02: ST02 is %q, SE02 is %q", t.ControlNumberST02, t.ControlNumberSE02))
			}
			if count(t.SegmentCountSE01) != transactionLengths[k] {
				errs = append(errs, fmt.Errorf("x12837i: checking SE01 of transaction %q: SE01 is %s, ST through SE is %d segments", t.ControlNumberST02, t.SegmentCountSE01, transactionLengths[k]))
			}
			k++
		}
	}
	return errors.Join(errs...)
}

// count reads a trailer count that parseTrailer has checked is a whole
// number; -1 for anything else, which no trailer count equals.
func count(digits string) int {
	n, err := strconv.Atoi(digits)
	if err != nil {
		return -1
	}
	return n
}
