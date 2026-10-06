// Package x12837i is a typed Go model of the Medicare 837I (institutional
// claim, HIPAA version 005010X223A2). Field and type names, and every
// citation below, come from frontend/src/data/x12Names837i.ts, built from
// Transmittal 907 (the CMS 837/835 flat file layout, CR 7409) and, where
// noted, the CMS Claims Processing Manual Ch. 25. "PDF p.N" means
// Transmittal 907's own PDF viewer page; see x12Names837i.ts's header for
// the full citation key, including why several NUBC-maintained elements
// (Type of Bill, Admission Type/Source, Patient Status, Condition,
// Occurrence, and Value Code) carry no code list here.
//
// Parse reads one interchange of bytes into an InterchangeISA: the ISA
// envelope, its functional groups, their transactions, and inside each
// transaction the submitter and receiver, each billing provider with its
// subscribers and patients, and every claim with its providers, its
// other-payer loops, and its service lines. Write puts an InterchangeISA
// back out as X12 with the file's own delimiters, byte for byte. Every
// position the standard defines has a field; a coded element is a string
// type whose value is the text the file carried, with Known reporting
// whether it is on the guide's list (enums.go); a code off the list is kept,
// never dropped; a segment the model has no slot for is kept whole as a
// RawSegment. Values, not pointers, throughout.
package x12837i

import (
	"fmt"
	"io"
	"strings"
)

// InterchangeISA is the outermost envelope: the ISA header, every functional
// group it carries, and the matching IEA trailer. Element names and their
// citations are frontend/src/data/x12Names837i.ts's ISA_EL_NAME and
// IEA_EL_NAME arrays (source: Transmittal 907, PDF p.7, 125).
//
// Unlike backend/pkg/x12835's InterchangeISA, the individual ISA qualifier
// and indicator elements (ISA01, ISA03, ISA05/07, ISA14, ISA15) are kept as
// plain strings here rather than as dedicated coded types: Transmittal 907
// does not name their individual codes (the same gap x12835 notes for its
// own ISA types), and the 837I's coded types were written from the
// sources, not from the three example files under testdata/x12837i.
//
// Every segment struct in this package carries Index, its 0-based position
// in the file (the ISA is 0), so Write can put the segments back in the
// order they came, typed and raw alike.
type InterchangeISA struct {
	Index int `json:"index,omitempty"` // the ISA's position in the file: always 0

	// The file's own delimiters and layout, read from the ISA and kept so
	// Write puts the file back exactly as it came: the element separator is
	// the byte after "ISA", the segment terminator the byte after ISA16, and
	// SegmentBreak the whitespace that follows each segment terminator ("\n"
	// in a file with one segment per line, "" in a file with none). The
	// component and repetition separators are ISA16 and ISA11 themselves.
	ElementSeparator  byte   `json:"elementSeparator,omitempty"`
	SegmentTerminator byte   `json:"segmentTerminator,omitempty"`
	SegmentBreak      string `json:"segmentBreak,omitempty"`

	AuthQualifierISA01 string `json:"authQualifierISA01,omitempty"` // Authorization Information Qualifier — Transmittal 907, PDF p.7
	AuthInfoISA02      string `json:"authInfoISA02,omitempty"`      // Authorization Information — Transmittal 907, PDF p.7

	SecQualifierISA03 string `json:"secQualifierISA03,omitempty"` // Security Information Qualifier — Transmittal 907, PDF p.7
	SecInfoISA04      string `json:"secInfoISA04,omitempty"`      // Security Information — Transmittal 907, PDF p.7

	SenderIDQualifierISA05 string `json:"senderIDQualifierISA05,omitempty"` // Interchange ID Qualifier — Transmittal 907, PDF p.7
	SenderIDISA06          string `json:"senderIDISA06,omitempty"`          // Interchange Sender ID — Transmittal 907, PDF p.7

	ReceiverIDQualifierISA07 string `json:"receiverIDQualifierISA07,omitempty"` // Interchange ID Qualifier — Transmittal 907, PDF p.7
	ReceiverIDISA08          string `json:"receiverIDISA08,omitempty"`          // Interchange Receiver ID — Transmittal 907, PDF p.7

	// DateISA09 is sent as YYMMDD; Year is 2000 plus YY, following the same
	// convention noted in backend/pkg/x12835's InterchangeISA.
	DateISA09 Date `json:"dateISA09,omitempty"` // Interchange Date — Transmittal 907, PDF p.7
	TimeISA10 Time `json:"timeISA10,omitempty"` // Interchange Time — Transmittal 907, PDF p.7

	// RepetitionSepISA11 is the single character the interchange uses as its
	// repetition separator. It is kept as received, not interpreted here -
	// the parser is the one place that uses it as a delimiter.
	RepetitionSepISA11 string `json:"repetitionSepISA11,omitempty"` // Repetition Separator — Transmittal 907, PDF p.7

	VersionISA12       string `json:"versionISA12,omitempty"`       // Interchange Control Version Number — Transmittal 907, PDF p.7
	ControlNumberISA13 string `json:"controlNumberISA13,omitempty"` // Interchange Control Number — Transmittal 907, PDF p.7

	AckRequestedISA14 string `json:"ackRequestedISA14,omitempty"` // Acknowledgement Requested — Transmittal 907, PDF p.7
	UsageISA15        string `json:"usageISA15,omitempty"`        // Interchange Usage Indicator — Transmittal 907, PDF p.7

	// ComponentSepISA16 is the single character the interchange uses as its
	// component element separator, kept as received for the same reason as
	// RepetitionSepISA11.
	ComponentSepISA16 string `json:"componentSepISA16,omitempty"` // Component Element Separator — Transmittal 907, PDF p.7

	FunctionalGroups []FunctionalGroupGS `json:"functionalGroups,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between ISA and IEA, outside any GS, that the model does not place

	// Trailer: IEA. GroupCountIEA01 and ControlNumberIEA02 are envelope
	// counts, not business data; the parser checks them against the groups
	// actually present and against ControlNumberISA13. The count is kept as
	// the digits the file sent, so it is written back the same way.
	IndexIEA           int    `json:"indexIEA,omitempty"`           // the IEA's position in the file
	GroupCountIEA01    string `json:"groupCountIEA01,omitempty"`    // Number of Included Functional Groups — Transmittal 907, PDF p.125
	ControlNumberIEA02 string `json:"controlNumberIEA02,omitempty"` // Interchange Control Number — Transmittal 907, PDF p.125
}

// Parse reads one interchange of 837I transactions into the typed model in
// one pass over its segments.
//
// The delimiters come from the ISA: the element separator is the byte after
// "ISA", the component separator is ISA16, the segment terminator is the
// byte after ISA16, and the repetition separator is ISA11. Whitespace
// between segments is skipped; the whitespace after the ISA's terminator
// is kept on the interchange as its SegmentBreak, with the two delimiters,
// so Write can put the file back as it came.
//
// Loops follow Transmittal 907's 837I layout. ISA holds GS, GS holds ST, and
// ST holds the BHT, the submitter (NM1*41, 1000A) and receiver (NM1*40,
// 1000B) loops, and the hierarchical loops. Each HL starts a loop that runs
// to the next HL or the SE. HL03 gives the loop's level and HL02 its parent:
// a billing provider (20) has no parent and goes under ST; a subscriber (22)
// goes under the billing provider whose HL01 equals its HL02; a patient (23)
// goes under the subscriber whose HL01 equals its HL02. An HL whose level is
// off the list, or whose parent is not a placed loop one level up, is kept
// raw on the transaction together with every segment of its loop.
//
// Within an HL loop, each NM1 starts a name loop that takes the N3, N4, REF,
// PER, PRV, DMG, and DTP segments after it (2010AA, 2010AB, 2010BA, 2010BB,
// 2010CA). In a subscriber or patient loop, each CLM starts a claim (2300)
// that runs to the next CLM, HL, or SE. Within a claim, an NM1 starts a
// claim-level provider loop (2310A-F); each SBR starts an other-subscriber
// loop (2320) whose NM1s start loops 2330A-I; and each LX starts a service
// line (2400) that runs to the next LX. Within a service line, LIN starts a
// drug loop (2410) that takes CTP and REF, NM1 starts a line provider loop
// (2420A-D), and SVD starts a line adjudication loop (2430) that takes CAS,
// DTP, and AMT. Every envelope segment ends every loop.
//
// Every element goes into its typed field, and every typed segment records
// its position in the file (Index). A coded value is kept as the text the
// file sent, on the list or off it (see enums.go). A segment the model has
// no slot for is kept whole, in file order, as a RawSegment on the struct
// that encloses it. That includes a whole NM1 loop whose NM101 is off the
// list, has no field at its level, or repeats a role already placed (the
// 2010AC pay-to plan loop among them), and an HI segment with a composite
// whose qualifier has no field.
//
// The parser keeps every element and component position X12 defines for a
// segment, whether or not Medicare uses it: a "Not Used" or situational
// position still gets a field and is parsed, never rejected, so text found
// there is never silently dropped; a position past the ones the model
// names goes to the struct's ExtraElements as sent. Medicare's own usage
// rules (Not Used, situational, required) are a restriction on top of the
// X12 definition, not part of it, so this package does not enforce them;
// the caller checks usage against Transmittal 907/579 if it needs to.
//
// Parse returns an error naming the element and the segment's 0-based
// position in the file (the ISA is 0) only when the model cannot hold the
// text at all: a malformed amount, date, time, or count; a non-empty
// element or component past the last position X12 defines for its segment;
// an envelope segment out of place; or input after the IEA. After the pass
// it checks each envelope trailer: SE01 against the segments from ST
// through SE, SE02 against ST02, GE01 against the transactions in the
// group, GE02 against GS06, IEA01 against the groups, and IEA02 against
// ISA13.
func Parse(data []byte) (InterchangeISA, error) {
	in, err := splitSegments(data)
	if err != nil {
		return InterchangeISA{}, err
	}
	p := parser{segments: in.segments, component: in.component}
	isa, err := p.parseInterchange()
	if err != nil {
		return InterchangeISA{}, err
	}
	isa.ElementSeparator = in.element
	isa.SegmentTerminator = in.terminator
	isa.SegmentBreak = in.segmentBreak
	if err := checkEnvelopes(isa, p.transactionLengths); err != nil {
		return InterchangeISA{}, err
	}
	return isa, nil
}

// Write puts an interchange back as X12, byte for byte as Parse read it:
// every typed segment and every raw segment in the order of its Index,
// each element in its position with the file's own element separator,
// composites joined with the file's own component separator, trailing
// empty elements and components left off the way X12 does, each segment
// closed with the file's segment terminator and followed by its
// SegmentBreak. The ISA is written with its elements as stored, so its
// fixed widths come back as they went in.
//
// Write is the proof that the model holds everything the file sent: with
// the delimiters and layout the parser kept, Write(Parse(bytes)) is bytes
// again, and the round-trip test holds that for every example file. A tree
// built by hand needs an Index on every segment, 0 for the ISA and one
// more for each segment after it; a missing or repeated position is an
// error, since a gap would mean a segment was lost. An interchange whose
// delimiters are zero is written with '*', ':', and '~'.
func Write(w io.Writer, isa InterchangeISA) error {
	element := isa.ElementSeparator
	if element == 0 {
		element = '*'
	}
	terminator := isa.SegmentTerminator
	if terminator == 0 {
		terminator = '~'
	}
	component := byte(':')
	if len(isa.ComponentSepISA16) == 1 {
		component = isa.ComponentSepISA16[0]
	}
	c := collector{component: component, healthInfo: map[int][]HealthInfoHI{}}
	c.interchange(isa)
	segments, err := c.ordered()
	if err != nil {
		return err
	}
	var b strings.Builder
	for _, s := range segments {
		b.WriteString(s.id)
		for _, e := range s.elements {
			b.WriteByte(element)
			b.WriteString(e)
		}
		b.WriteByte(terminator)
		b.WriteString(isa.SegmentBreak)
	}
	if _, err := io.WriteString(w, b.String()); err != nil {
		return fmt.Errorf("x12837i: writing: %w", err)
	}
	return nil
}
