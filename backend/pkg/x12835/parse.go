package x12835

import (
	"errors"
	"fmt"
	"strings"
)

// inputSegment is one segment of the input, split on the element separator.
type inputSegment struct {
	index    int      // position in the file, 0-based; the ISA is 0
	id       string   // the segment ID, e.g. "CLP"
	elements []string // elements[0] is element 01
}

// el returns element n (1-based), or "" when the segment ends before it.
func (s inputSegment) el(n int) string {
	if n < 1 || n > len(s.elements) {
		return ""
	}
	return s.elements[n-1]
}

// name returns element n's reference, e.g. "CLP04".
func (s inputSegment) name(n int) string {
	return fmt.Sprintf("%s%02d", s.id, n)
}

// checkLast returns an error when the segment has an element after element
// last, the highest one the standard defines for this segment. splitSegments
// removes trailing empty elements, so any element past last carries text.
func (s inputSegment) checkLast(last int) error {
	if len(s.elements) > last {
		return fmt.Errorf("x12835: parsing %s in segment %d: the segment ends at %s", s.name(len(s.elements)), s.index, s.name(last))
	}
	return nil
}

func rawSegment(s inputSegment) RawSegment {
	return RawSegment{Index: s.index, ID: s.id, Elements: s.elements}
}

// splitSegments reads the delimiters from the ISA and splits data into
// segments. It returns the delimiters with them: the four separators, and
// the whitespace that followed the ISA's terminator, which Write puts after
// every segment. Whitespace before a segment is skipped, a blank segment is
// skipped, trailing empty elements are removed (they carry no text), and
// when the terminator is a newline a carriage return before it is removed.
func splitSegments(data []byte) ([]inputSegment, Delimiters, error) {
	text := string(data)
	start := 0
	for start < len(text) && isSpace(text[start]) {
		start++
	}
	if !strings.HasPrefix(text[start:], "ISA") || len(text) < start+4 {
		return nil, Delimiters{}, fmt.Errorf("x12835: splitting input: expected ISA")
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
		return nil, Delimiters{}, fmt.Errorf("x12835: splitting ISA: incomplete header")
	}
	component := text[isa16]
	terminator := text[isa16+1]
	isaFields := strings.Split(text[start:isa16+1], string(element))
	repetition := isaFields[11]
	if len(repetition) != 1 {
		return nil, Delimiters{}, fmt.Errorf("x12835: splitting ISA: ISA11 is not one character")
	}
	if element == component || element == terminator || element == repetition[0] ||
		component == terminator || component == repetition[0] || terminator == repetition[0] {
		return nil, Delimiters{}, fmt.Errorf("x12835: splitting ISA: the four separators must differ")
	}
	suffixEnd := isa16 + 2
	for suffixEnd < len(text) && isSpace(text[suffixEnd]) {
		suffixEnd++
	}
	delimiters := Delimiters{Element: element, Component: component, Repetition: repetition[0], Terminator: terminator, Suffix: text[isa16+2 : suffixEnd]}

	segments, err := appendSegments([]inputSegment{{index: 0, id: "ISA", elements: isaFields[1:]}}, text, isa16+2, element, terminator)
	if err != nil {
		return nil, Delimiters{}, err
	}
	return segments, delimiters, nil
}

// appendSegments splits text from pos to its end into segments and appends
// them to segments, numbering each by its position in the result. Parse
// starts it after the ISA; ParseClaim starts it at a claim's CLP. Whitespace
// before a segment is skipped, a blank segment is skipped, trailing empty
// elements are removed (they carry no text), and when the terminator is a
// newline a carriage return before it is removed.
func appendSegments(segments []inputSegment, text string, pos int, element, terminator byte) ([]inputSegment, error) {
	for {
		for pos < len(text) && isSpace(text[pos]) {
			pos++
		}
		if pos == len(text) {
			return segments, nil
		}
		if text[pos] == terminator {
			pos++
			continue
		}
		end := strings.IndexByte(text[pos:], terminator)
		if end < 0 {
			return nil, fmt.Errorf("x12835: splitting segment %d: missing segment terminator", len(segments))
		}
		content := text[pos : pos+end]
		pos += end + 1
		if terminator == '\n' {
			content = strings.TrimSuffix(content, "\r")
		}
		fields := strings.Split(content, string(element))
		if !validID(fields[0]) {
			return nil, fmt.Errorf("x12835: splitting segment %d: segment ID is not 2 or 3 capital letters or digits", len(segments))
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

func (p *parser) parseInterchange() (InterchangeISA, error) {
	isa, err := parseISA(p.take()) // splitSegments always returns the ISA first
	if err != nil {
		return InterchangeISA{}, err
	}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return InterchangeISA{}, fmt.Errorf("x12835: parsing interchange: missing IEA")
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
				return InterchangeISA{}, fmt.Errorf("x12835: parsing segment %d: input continues after the IEA", p.peek().index)
			}
			return isa, nil
		case "ISA", "GE", "ST", "SE":
			return InterchangeISA{}, fmt.Errorf("x12835: parsing segment %d: %s outside a functional group", next.index, next.id)
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
			return FunctionalGroupGS{}, fmt.Errorf("x12835: parsing functional group at segment %d: missing GE", gs.index)
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
			return FunctionalGroupGS{}, fmt.Errorf("x12835: parsing segment %d: %s inside the functional group at segment %d (missing GE)", next.index, next.id, gs.index)
		default:
			group.Raw = append(group.Raw, rawSegment(p.take()))
		}
	}
}

func (p *parser) parseTransaction() (TransactionST, error) {
	st := p.take()
	t, err := parseST(st)
	if err != nil {
		return TransactionST{}, err
	}
	for {
		next := p.peek()
		switch next.id {
		case "":
			return TransactionST{}, fmt.Errorf("x12835: parsing transaction at segment %d: missing SE", st.index)
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
			return TransactionST{}, fmt.Errorf("x12835: parsing segment %d: %s inside the transaction at segment %d (missing SE)", next.index, next.id, st.index)
		case "N1":
			// An N1 whose N101 is off the list, or a second PR or PE loop,
			// has no slot: its whole loop is kept raw, text included.
			entity := EntityIDN101(next.el(1))
			payer := entity == EntityIDN101PR && t.PayerN1.EntityIDN101 == EntityIDN101NotSent
			payee := entity == EntityIDN101PE && t.PayeeN1.EntityIDN101 == EntityIDN101NotSent
			if !payer && !payee {
				t.Raw = append(t.Raw, p.takePartyLoop()...)
				continue
			}
			party, err := p.parseParty()
			if err != nil {
				return TransactionST{}, err
			}
			if payer {
				t.PayerN1 = party
			} else {
				t.PayeeN1 = party
			}
		case "LX":
			lx, err := p.parseHeaderNumber()
			if err != nil {
				return TransactionST{}, err
			}
			t.HeaderNumbers = append(t.HeaderNumbers, lx)
		case "PLB":
			plb, err := parsePLB(p.take(), p.component)
			if err != nil {
				return TransactionST{}, err
			}
			t.ProviderAdjustments = append(t.ProviderAdjustments, plb)
		default:
			if err := t.add(p.take()); err != nil {
				return TransactionST{}, err
			}
		}
	}
}

// inPartyLoop reports whether a segment with this ID continues an N1 loop.
func inPartyLoop(id string) bool {
	return id == "N3" || id == "N4" || id == "REF" || id == "PER" || id == "RDM"
}

func (p *parser) parseParty() (PartyN1, error) {
	party, err := parseN1(p.take())
	if err != nil {
		return PartyN1{}, err
	}
	for inPartyLoop(p.peek().id) {
		if err := party.add(p.take()); err != nil {
			return PartyN1{}, err
		}
	}
	return party, nil
}

func (p *parser) takePartyLoop() []RawSegment {
	raw := []RawSegment{rawSegment(p.take())}
	for inPartyLoop(p.peek().id) {
		raw = append(raw, rawSegment(p.take()))
	}
	return raw
}

// envelope reports whether a segment with this ID opens or closes an
// envelope. Every envelope segment ends LX, CLP, and SVC loops so that the
// envelope level can close the loop or report the segment out of place; the
// 835 viewer ends loops at SE, GE, and IEA only.
func envelope(id string) bool {
	switch id {
	case "ISA", "IEA", "GS", "GE", "ST", "SE":
		return true
	}
	return false
}

// endsHeaderNumber reports whether a segment with this ID ends an LX loop;
// "" is the end of the input.
func endsHeaderNumber(id string) bool {
	return id == "" || id == "LX" || id == "PLB" || envelope(id)
}

func endsClaim(id string) bool {
	return id == "CLP" || endsHeaderNumber(id)
}

func endsServiceLine(id string) bool {
	return id == "SVC" || endsClaim(id)
}

func (p *parser) parseHeaderNumber() (HeaderNumberLX, error) {
	s := p.take()
	if err := s.checkLast(1); err != nil {
		return HeaderNumberLX{}, err
	}
	number, err := parseWholeNumber(s, 1)
	if err != nil {
		return HeaderNumberLX{}, err
	}
	lx := HeaderNumberLX{Index: s.index, NumberLX01: number}
	for !endsHeaderNumber(p.peek().id) {
		if p.peek().id == "CLP" {
			claim, err := p.parseClaim()
			if err != nil {
				return HeaderNumberLX{}, err
			}
			lx.Claims = append(lx.Claims, claim)
			continue
		}
		if err := lx.add(p.take()); err != nil {
			return HeaderNumberLX{}, err
		}
	}
	return lx, nil
}

func (p *parser) parseClaim() (ClaimCLP, error) {
	claim, err := parseCLP(p.take())
	if err != nil {
		return ClaimCLP{}, err
	}
	for !endsClaim(p.peek().id) {
		if p.peek().id == "SVC" {
			line, err := p.parseServiceLine()
			if err != nil {
				return ClaimCLP{}, err
			}
			claim.ServiceLines = append(claim.ServiceLines, line)
			continue
		}
		if err := claim.add(p.take()); err != nil {
			return ClaimCLP{}, err
		}
	}
	return claim, nil
}

func (p *parser) parseServiceLine() (ServiceLineSVC, error) {
	line, err := parseSVC(p.take(), p.component)
	if err != nil {
		return ServiceLineSVC{}, err
	}
	for !endsServiceLine(p.peek().id) {
		if err := line.add(p.take()); err != nil {
			return ServiceLineSVC{}, err
		}
	}
	return line, nil
}

// add places one segment that sits directly under ST, outside any loop: the
// first BPR, TRN, and CUR, every REF, and the first DTM*405. Anything else
// is kept raw.
func (t *TransactionST) add(s inputSegment) error {
	switch s.id {
	case "BPR":
		if t.Payment.Set {
			t.Raw = append(t.Raw, rawSegment(s))
			return nil
		}
		payment, err := parseBPR(s)
		if err != nil {
			return err
		}
		t.Payment = payment
		return nil
	case "TRN":
		if t.Trace.Set {
			t.Raw = append(t.Raw, rawSegment(s))
			return nil
		}
		trace, err := parseTRN(s)
		if err != nil {
			return err
		}
		t.Trace = trace
		return nil
	case "CUR":
		if t.Currency.Set {
			t.Raw = append(t.Raw, rawSegment(s))
			return nil
		}
		currency, err := parseCUR(s)
		if err != nil {
			return err
		}
		t.Currency = currency
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		t.HeaderRefs = append(t.HeaderRefs, ref)
		return nil
	case "DTM":
		// Any qualifier but 405, and a second 405, is kept raw.
		if DateQualifierDTM01(s.el(1)) != DateQualifierDTM01405 || t.ProductionDateDTM405.DateQualifierDTM01 != DateQualifierDTM01NotSent {
			t.Raw = append(t.Raw, rawSegment(s))
			return nil
		}
		dtm, err := parseDTM(s)
		if err != nil {
			return err
		}
		t.ProductionDateDTM405 = dtm
		return nil
	}
	t.Raw = append(t.Raw, rawSegment(s))
	return nil
}

// add places one segment of an N1 loop after the N1: the first N3, N4, and
// RDM, and every REF and PER. Anything else is kept raw.
func (n *PartyN1) add(s inputSegment) error {
	switch s.id {
	case "N3":
		if n.Address.Set {
			n.Raw = append(n.Raw, rawSegment(s))
			return nil
		}
		address, err := parseN3(s)
		if err != nil {
			return err
		}
		n.Address = address
		return nil
	case "N4":
		if n.Location.Set {
			n.Raw = append(n.Raw, rawSegment(s))
			return nil
		}
		location, err := parseN4(s)
		if err != nil {
			return err
		}
		n.Location = location
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		n.References = append(n.References, ref)
		return nil
	case "PER":
		contact, err := parsePER(s)
		if err != nil {
			return err
		}
		n.Contacts = append(n.Contacts, contact)
		return nil
	case "RDM":
		if n.DeliveryRDM.Set {
			n.Raw = append(n.Raw, rawSegment(s))
			return nil
		}
		delivery, err := parseRDM(s)
		if err != nil {
			return err
		}
		n.DeliveryRDM = delivery
		return nil
	}
	n.Raw = append(n.Raw, rawSegment(s))
	return nil
}

// add places one segment of an LX loop before its first claim: the first
// TS3 and TS2. Anything else is kept raw.
func (lx *HeaderNumberLX) add(s inputSegment) error {
	switch s.id {
	case "TS3":
		if lx.ProviderSummary.Set {
			lx.Raw = append(lx.Raw, rawSegment(s))
			return nil
		}
		summary, err := parseTS3(s)
		if err != nil {
			return err
		}
		lx.ProviderSummary = summary
		return nil
	case "TS2":
		if lx.InstitutionalSummary.Set {
			lx.Raw = append(lx.Raw, rawSegment(s))
			return nil
		}
		summary, err := parseTS2(s)
		if err != nil {
			return err
		}
		lx.InstitutionalSummary = summary
		return nil
	}
	lx.Raw = append(lx.Raw, rawSegment(s))
	return nil
}

// add places one segment of a claim before its first service line: every
// CAS, REF, AMT, and QTY; the first NM1 for each entity code the claim
// holds; the first MIA and MOA; the first DTM for each qualifier the claim
// holds, when it carries a date; and the first PER that carries PER01.
// Anything else is kept raw.
func (c *ClaimCLP) add(s inputSegment) error {
	switch s.id {
	case "CAS":
		adjustment, err := parseCAS(s)
		if err != nil {
			return err
		}
		c.Adjustments = append(c.Adjustments, adjustment)
		return nil
	case "NM1":
		return c.addName(s)
	case "MIA":
		if c.InpatientMIA.Set {
			c.Raw = append(c.Raw, rawSegment(s))
			return nil
		}
		inpatient, err := parseMIA(s)
		if err != nil {
			return err
		}
		c.InpatientMIA = inpatient
		return nil
	case "MOA":
		if c.OutpatientMOA.Set {
			c.Raw = append(c.Raw, rawSegment(s))
			return nil
		}
		outpatient, err := parseMOA(s)
		if err != nil {
			return err
		}
		c.OutpatientMOA = outpatient
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		c.References = append(c.References, ref)
		return nil
	case "DTM":
		return c.addDate(s)
	case "PER":
		if c.ContactPER.FunctionPER01 != ContactFunctionPER01NotSent || s.el(1) == "" {
			c.Raw = append(c.Raw, rawSegment(s))
			return nil
		}
		contact, err := parsePER(s)
		if err != nil {
			return err
		}
		c.ContactPER = contact
		return nil
	case "AMT":
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		c.Amounts = append(c.Amounts, amount)
		return nil
	case "QTY":
		quantity, err := parseQTY(s)
		if err != nil {
			return err
		}
		c.Quantities = append(c.Quantities, quantity)
		return nil
	}
	c.Raw = append(c.Raw, rawSegment(s))
	return nil
}

// addName places an NM1 in the field for its entity code (NM101) when that
// field is still empty; otherwise the NM1 is kept raw.
func (c *ClaimCLP) addName(s inputSegment) error {
	entity := EntityIDNM101(s.el(1)) // an NM101 with no field here is kept raw below
	open := false
	switch entity {
	case EntityIDNM101QC:
		open = c.PatientNM1QC.EntityIDNM101 == EntityIDNM101NotSent
	case EntityIDNM10174:
		open = c.CorrectedPatientNM174.EntityIDNM101 == EntityIDNM101NotSent
	case EntityIDNM10182:
		open = c.RenderingNM182.EntityIDNM101 == EntityIDNM101NotSent
	case EntityIDNM101TT:
		open = c.CrossoverNM1TT.EntityIDNM101 == EntityIDNM101NotSent
	case EntityIDNM101PR:
		open = c.CorrectedPayerNM1PR.EntityIDNM101 == EntityIDNM101NotSent
	}
	if !open {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	name, err := parseNM1(s)
	if err != nil {
		return err
	}
	switch entity {
	case EntityIDNM101QC:
		c.PatientNM1QC = name
	case EntityIDNM10174:
		c.CorrectedPatientNM174 = name
	case EntityIDNM10182:
		c.RenderingNM182 = name
	case EntityIDNM101TT:
		c.CrossoverNM1TT = name
	case EntityIDNM101PR:
		c.CorrectedPayerNM1PR = name
	}
	return nil
}

// addDate places a claim DTM in the field for its qualifier (DTM01) when
// that field is still empty and DTM02 carries a date; otherwise the DTM is
// kept raw, so a DTM without a date, or a second one, still shows it was
// sent.
func (c *ClaimCLP) addDate(s inputSegment) error {
	qualifier := DateQualifierDTM01(s.el(1)) // a qualifier with no claim field is kept raw below
	open := false
	switch qualifier {
	case DateQualifierDTM01232:
		open = c.StatementFromDTM232.DateQualifierDTM01 == DateQualifierDTM01NotSent
	case DateQualifierDTM01233:
		open = c.StatementToDTM233.DateQualifierDTM01 == DateQualifierDTM01NotSent
	case DateQualifierDTM01050:
		open = c.ReceivedDTM050.DateQualifierDTM01 == DateQualifierDTM01NotSent
	case DateQualifierDTM01036:
		open = c.CoverageExpirationDTM036.DateQualifierDTM01 == DateQualifierDTM01NotSent
	}
	if !open || s.el(2) == "" {
		c.Raw = append(c.Raw, rawSegment(s))
		return nil
	}
	dtm, err := parseDTM(s)
	if err != nil {
		return err
	}
	switch qualifier {
	case DateQualifierDTM01232:
		c.StatementFromDTM232 = dtm
	case DateQualifierDTM01233:
		c.StatementToDTM233 = dtm
	case DateQualifierDTM01050:
		c.ReceivedDTM050 = dtm
	case DateQualifierDTM01036:
		c.CoverageExpirationDTM036 = dtm
	}
	return nil
}

// add places one segment of a service line after the SVC: every CAS, AMT,
// QTY, and LQ; every REF, the first REF*6R, REF*LU, REF*HPI, and REF*0K
// with a value going to their named fields; and the first DTM*472, DTM*150,
// and DTM*151 that carries a date. Anything else is kept raw.
func (l *ServiceLineSVC) add(s inputSegment) error {
	switch s.id {
	case "DTM":
		return l.addDate(s)
	case "CAS":
		adjustment, err := parseCAS(s)
		if err != nil {
			return err
		}
		l.Adjustments = append(l.Adjustments, adjustment)
		return nil
	case "REF":
		ref, err := parseREF(s)
		if err != nil {
			return err
		}
		switch {
		case ref.RefQualifierREF01 == RefQualifierREF016R && l.LineControlREF6R.Value == "" && ref.Value != "":
			l.LineControlREF6R = ref
		case ref.RefQualifierREF01 == RefQualifierREF01LU && l.ServiceIDREFLU.Value == "" && ref.Value != "":
			l.ServiceIDREFLU = ref
		case ref.RefQualifierREF01 == RefQualifierREF01HPI && l.RenderingREFHPI.Value == "" && ref.Value != "":
			l.RenderingREFHPI = ref
		case ref.RefQualifierREF01 == RefQualifierREF010K && l.PolicyREF0K.Value == "" && ref.Value != "":
			l.PolicyREF0K = ref
		default:
			l.References = append(l.References, ref)
		}
		return nil
	case "AMT":
		amount, err := parseAMT(s)
		if err != nil {
			return err
		}
		l.Amounts = append(l.Amounts, amount)
		return nil
	case "QTY":
		quantity, err := parseQTY(s)
		if err != nil {
			return err
		}
		l.Quantities = append(l.Quantities, quantity)
		return nil
	case "LQ":
		remark, err := parseLQ(s)
		if err != nil {
			return err
		}
		l.Remarks = append(l.Remarks, remark)
		return nil
	}
	l.Raw = append(l.Raw, rawSegment(s))
	return nil
}

// addDate places a service line DTM the same way ClaimCLP.addDate places a
// claim DTM, for qualifiers 472, 150, and 151.
func (l *ServiceLineSVC) addDate(s inputSegment) error {
	qualifier := DateQualifierDTM01(s.el(1)) // a qualifier with no service line field is kept raw below
	open := false
	switch qualifier {
	case DateQualifierDTM01472:
		open = l.ServiceDateDTM472.DateQualifierDTM01 == DateQualifierDTM01NotSent
	case DateQualifierDTM01150:
		open = l.ServiceDateDTM150.DateQualifierDTM01 == DateQualifierDTM01NotSent
	case DateQualifierDTM01151:
		open = l.ServiceDateDTM151.DateQualifierDTM01 == DateQualifierDTM01NotSent
	}
	if !open || s.el(2) == "" {
		l.Raw = append(l.Raw, rawSegment(s))
		return nil
	}
	dtm, err := parseDTM(s)
	if err != nil {
		return err
	}
	switch qualifier {
	case DateQualifierDTM01472:
		l.ServiceDateDTM472 = dtm
	case DateQualifierDTM01150:
		l.ServiceDateDTM150 = dtm
	case DateQualifierDTM01151:
		l.ServiceDateDTM151 = dtm
	}
	return nil
}

// checkEnvelopes compares each envelope trailer with its header and with
// what the pass found, and returns every mismatch. transactionLengths holds,
// in file order, the number of segments from each ST through its SE.
func checkEnvelopes(isa InterchangeISA, transactionLengths []int) error {
	var errs []error
	if isa.ControlNumberISA13 != isa.ControlNumberIEA02 {
		errs = append(errs, fmt.Errorf("x12835: checking IEA02: ISA13 is %q, IEA02 is %q", isa.ControlNumberISA13, isa.ControlNumberIEA02))
	}
	if isa.GroupCountIEA01 != len(isa.FunctionalGroups) {
		errs = append(errs, fmt.Errorf("x12835: checking IEA01: IEA01 is %d, the interchange holds %d functional groups", isa.GroupCountIEA01, len(isa.FunctionalGroups)))
	}
	k := 0
	for _, g := range isa.FunctionalGroups {
		if g.ControlNumberGS06 != g.ControlNumberGE02 {
			errs = append(errs, fmt.Errorf("x12835: checking GE02: GS06 is %q, GE02 is %q", g.ControlNumberGS06, g.ControlNumberGE02))
		}
		if g.TransactionCountGE01 != len(g.Transactions) {
			errs = append(errs, fmt.Errorf("x12835: checking GE01 of group %q: GE01 is %d, the group holds %d transactions", g.ControlNumberGS06, g.TransactionCountGE01, len(g.Transactions)))
		}
		for _, t := range g.Transactions {
			if t.ControlNumberST02 != t.ControlNumberSE02 {
				errs = append(errs, fmt.Errorf("x12835: checking SE02: ST02 is %q, SE02 is %q", t.ControlNumberST02, t.ControlNumberSE02))
			}
			if t.SegmentCountSE01 != transactionLengths[k] {
				errs = append(errs, fmt.Errorf("x12835: checking SE01 of transaction %q: SE01 is %d, ST through SE is %d segments", t.ControlNumberST02, t.SegmentCountSE01, transactionLengths[k]))
			}
			k++
		}
	}
	return errors.Join(errs...)
}
