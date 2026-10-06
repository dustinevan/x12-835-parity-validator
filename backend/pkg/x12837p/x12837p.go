// Package x12837p is a typed Go model of the X12 837 Health Care Claim:
// Professional, version 005010X222A1, as Medicare receives it. Parse reads
// one interchange of bytes into an InterchangeISA: the ISA envelope, its
// functional groups, their transactions, and inside each transaction the
// submitter and receiver, each billing provider with its subscribers and
// patients, and every claim with its providers, its other-payer loops, and
// its service lines. Write puts an InterchangeISA back out as X12 with the
// file's own delimiters, byte for byte; Bytes returns the same as a slice.
//
// The model is straight definition. Every position the standard defines
// has a field, named by a short name and its element reference
// (ChargeCLM02, ClaimFrequencyCLM05_3), with the public name and its
// citation in the field's comment. An amount is an exact decimal with the
// text as sent beside it; a date is a Date; a coded element is a string
// type whose value is the text the file carried, with Known reporting
// whether it is on the guide's list and Description naming it when it is
// (enums.go). A code off the list is kept, never dropped, and never a parse
// failure; a malformed amount, date, or count is. A segment the model has
// no slot for is kept whole, in file order, as a RawSegment. Values, not
// pointers, throughout. Element names come from
// frontend/src/data/x12Names837p.ts; the code lists cite the implementation
// guide, Transmittal 907, and the CEDI 837P guide under docs/reference/.
package x12837p

import (
	"fmt"
	"io"
)

// InterchangeISA is the outermost envelope: the ISA header, every
// functional group it carries, and the matching IEA trailer. Element names
// and citations are frontend/src/data/x12Names837p.ts's ISA_EL_NAME and
// IEA_EL_NAME arrays (source: Transmittal 907, PDF p.133 footer 1, PDF
// p.239 footer 107). The ISA/IEA envelope is the same X12 005010 structure
// the 835 side of this project already models (backend/pkg/x12835's
// InterchangeISA); this type is independent of that package.
//
// The ISA is always the file's first segment, so it carries no Index; its
// sixteen elements have fixed widths and the parser keeps each one's text
// exactly, padding included.
type InterchangeISA struct {
	// Delimiters is what the file fixes by position rather than by element,
	// kept so the file writes back with its own layout (see Write).
	Delimiters Delimiters `json:"delimiters,omitempty"`

	AuthQualifierISA01 AuthQualifierISA01 `json:"authQualifierISA01,omitempty"` // Authorization Information Qualifier
	AuthInfoISA02      string             `json:"authInfoISA02,omitempty"`      // Authorization Information

	SecQualifierISA03 SecQualifierISA03 `json:"secQualifierISA03,omitempty"` // Security Information Qualifier
	SecInfoISA04      string            `json:"secInfoISA04,omitempty"`      // Security Information

	SenderIDQualifierISA05 InterchangeIDQualifierISA05 `json:"senderIDQualifierISA05,omitempty"` // Interchange ID Qualifier
	SenderIDISA06          string                      `json:"senderIDISA06,omitempty"`          // Interchange Sender ID

	ReceiverIDQualifierISA07 InterchangeIDQualifierISA05 `json:"receiverIDQualifierISA07,omitempty"` // Interchange ID Qualifier
	ReceiverIDISA08          string                      `json:"receiverIDISA08,omitempty"`          // Interchange Receiver ID

	// DateISA09 is sent as YYMMDD, with Year taken as 2000 plus YY, the same
	// convention as the 835 side (backend/pkg/x12835's InterchangeISA).
	DateISA09 Date `json:"dateISA09,omitempty"` // Interchange Date
	TimeISA10 Time `json:"timeISA10,omitempty"` // Interchange Time

	// RepetitionSepISA11 is the single character the interchange uses as its
	// repetition separator. It is kept as received, not interpreted here -
	// the parser is the one place that uses it as a delimiter.
	RepetitionSepISA11 string `json:"repetitionSepISA11,omitempty"` // Repetition Separator

	VersionISA12       string `json:"versionISA12,omitempty"`       // Interchange Control Version Number
	ControlNumberISA13 string `json:"controlNumberISA13,omitempty"` // Interchange Control Number

	AckRequestedISA14 AckRequestedISA14 `json:"ackRequestedISA14,omitempty"` // Acknowledgement Requested
	UsageISA15        UsageISA15        `json:"usageISA15,omitempty"`        // Usage Indicator

	// ComponentSepISA16 is the single character the interchange uses as its
	// component element separator, kept as received for the same reason as
	// RepetitionSepISA11.
	ComponentSepISA16 string `json:"componentSepISA16,omitempty"` // Component Element Separator

	FunctionalGroups []FunctionalGroupGS `json:"functionalGroups,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between ISA and IEA, outside any GS, that the model does not place

	// Trailer: IEA. GroupCountIEA01 and ControlNumberIEA02 are envelope
	// counts, not business data; the parser checks them against the groups
	// actually present and against ControlNumberISA13.
	GroupCountIEA01    int    `json:"groupCountIEA01,omitempty"`    // Number of Included Functional Groups
	ControlNumberIEA02 string `json:"controlNumberIEA02,omitempty"` // Interchange Control Number
}

// Parse reads one interchange of 837P transactions into the typed model in
// one pass over its segments.
//
// The delimiters come from the ISA: the element separator is the byte after
// "ISA", the component separator is ISA16, the segment terminator is the
// byte after ISA16, and the repetition separator is ISA11. Whitespace
// between segments is skipped, and recorded on the interchange's
// Delimiters when every segment is followed by the same whitespace, so
// Write can put the file back the way it came.
//
// Loops follow Transmittal 907's 837P layout. ISA holds GS, GS holds ST, and
// ST holds the BHT, the 1000A submitter (NM1*41) and 1000B receiver (NM1*40)
// loops, and the HL loops. Each HL loop runs to the next HL or SE and is
// placed by its HL03 code and HL02 parent link: 20 opens a 2000A billing
// provider; 22 opens a 2000B subscriber under the 2000A whose HL01 equals
// its HL02; 23 opens a 2000C patient under the 2000B whose HL01 equals its
// HL02. An HL whose HL03 is none of those, or whose HL02 names no HL of the
// level above, has no slot: it and every segment up to the next HL are kept
// raw on the transaction, and so are its children, whose parent is then not
// found. When two HLs share an HL01, children attach to the later one.
//
// Within an HL loop an NM1 starts a name loop (2010AA-AC, 2010BA-BB, 2010CA)
// chosen by NM101, and a CLM starts a 2300 claim that runs to the next CLM,
// HL, or SE. Within a claim an NM1 starts a 2310A-F loop, an SBR a 2320
// other-subscriber loop (whose NM1s start 2330A-G loops), and an LX a 2400
// service line, which runs to the next LX, CLM, HL, or SE. Within a service
// line a LIN starts the 2410 drug loop, an NM1 a 2420A-H loop, an SVD a 2430
// adjudication loop, and an LQ a 2440 form loop. A name loop takes the
// segments its box holds in Transmittal 907 (some of N3, N4, REF, PER, PRV,
// DMG, and DTP); 2410 takes CTP and REF; 2430 takes CAS, DTP, and AMT; 2440
// takes FRM. Any other segment ends the nested loop and goes to the loop
// around it.
//
// Every element goes into its typed field, and every typed segment records
// its position in the file (Index), so the model can be written back in
// file order. A box that a qualifier selects (an NM1 role, a REF, DTP, AMT,
// QTY, CRC, or HI box) is chosen by the code Transmittal 579's Values column
// gives it (see enums.go). A coded element is kept as its text whether or
// not the code is on its list, so nothing a file sends is lost; Known on
// the value says whether it is listed. A segment the model has no slot for
// is kept whole, in file order, as a RawSegment on the struct that encloses
// it: an unmodeled segment ID, a repeat of a box that occurs once, a
// qualifier that matches no box, or a whole loop whose NM101 or HL03 matches
// no slot. An element after the last position the model names for a
// segment goes to that struct's Extra, and a component after the last a
// composite names goes to the composite's own extra field, kept as sent.
//
// The parser keeps every element and component position the 837P standard
// defines for a segment, whether or not Medicare uses it: a "Not Used" or
// situational position still gets a field and is parsed, never rejected, so
// text found there is never silently dropped. Medicare's own usage rules
// (Not Used, situational, required) are a restriction on top of the X12
// definition, not part of it, so this package does not enforce them; the
// caller checks usage against Transmittal 907/579 if it needs to.
//
// Parse returns an error naming the element and the segment's 0-based
// position in the file (the ISA is 0) only when the model cannot hold the
// text at all: a malformed amount, date, time, count, or diagnosis pointer;
// a GS01 or ST01 that says the file is not an 837; an envelope segment out
// of place; or input after the IEA. After the pass it checks each envelope
// trailer: SE01 against the segments from ST through SE, SE02 against ST02,
// GE01 against the transactions in the group, GE02 against GS06, IEA01
// against the groups, and IEA02 against ISA13.
func Parse(data []byte) (InterchangeISA, error) {
	segments, delimiters, component, err := splitSegments(data)
	if err != nil {
		return InterchangeISA{}, err
	}
	p := parser{segments: segments, component: component}
	isa, err := p.parseInterchange()
	if err != nil {
		return InterchangeISA{}, err
	}
	isa.Delimiters = delimiters
	if err := checkEnvelopes(isa, p.transactionLengths); err != nil {
		return InterchangeISA{}, err
	}
	return isa, nil
}

// Write writes isa to w as X12, using the interchange's own delimiters: the
// element separator and segment terminator in Delimiters, the component
// separator in ISA16, and the line break Delimiters.AfterSegment names
// after every segment terminator. Every typed segment is written with its
// elements in position, trailing empty elements and components left off
// the way X12 sends them, and the text a segment carried past the model's
// last position (Extra) after them; a raw segment is written as it was
// received; the ISA keeps its fixed widths because the parser kept its
// padding. Segments come out in file order: each typed segment and raw
// segment carries its Index from the parse, and the writer sorts a
// transaction's segments by it, so a parsed file writes back byte for
// byte. A model built in code, whose indexes are all zero, writes in the
// implementation guide's order instead.
//
// Write returns an error when the delimiters are not set: the element
// separator, segment terminator, and ISA16 must each be one byte.
func Write(w io.Writer, isa InterchangeISA) error {
	data, err := Bytes(isa)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("x12837p: writing interchange: %w", err)
	}
	return nil
}
