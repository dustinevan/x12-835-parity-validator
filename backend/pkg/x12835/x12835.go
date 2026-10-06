// Package x12835 is a typed Go model of the X12 835 Health Care Claim
// Payment/Advice, version 005010X221A1, as Medicare sends it. Parse reads
// one interchange of bytes into an InterchangeISA: the ISA envelope, its
// functional groups, their transactions, and inside each transaction the
// payment header, the payer and payee loops, the provider summaries, every
// claim with its service lines, and the provider-level adjustments. Write
// puts an InterchangeISA back out as X12 with the file's own delimiters,
// byte for byte. ParseClaim and WriteClaim do the same for one claim's
// bytes, which is how the validator stores a claim.
//
// The model is straight definition. Every position the standard defines
// has a field, named by a short name and its element reference
// (ChargeCLP03, LineControlREF6R), with the public name and its citation in
// the field's comment. An amount is an exact decimal with the text as sent
// beside it; a date is a Date; a coded element is a string type whose value
// is the text the file carried, with Known reporting whether it is on the
// guide's list and Description naming it when it is (enums.go). A code off
// the list is kept, never dropped, and never a parse failure; a malformed
// amount, date, or count is. A segment the model has no slot for is kept
// whole, in file order, as a RawSegment. Values, not pointers, throughout.
// Element names come from frontend/src/data/x12Names835.ts; the code lists
// cite the implementation guide and the CMS documents under docs/reference/.
package x12835

import (
	"fmt"
	"io"
)

// InterchangeISA is the outermost envelope: the ISA header, every functional
// group it carries, and the matching IEA trailer. Element names and their
// citations are frontend/src/data/x12Names835.ts's ISA_EL_NAME and
// IEA_EL_NAME arrays (source: CMS 835 Flat File, PDF p.1, 14-15). The ISA
// is always segment 0 of the file.
type InterchangeISA struct {
	// Delimiters are the file's own, read by Parse; Write writes with them.
	Delimiters Delimiters `json:"delimiters,omitempty"`

	AuthQualifierISA01 AuthQualifierISA01 `json:"authQualifierISA01,omitempty"` // Authorization Information Qualifier — CMS 835 Flat File, PDF p.1
	AuthInfoISA02      string             `json:"authInfoISA02,omitempty"`      // Authorization Information — CMS 835 Flat File, PDF p.1

	SecQualifierISA03 SecQualifierISA03 `json:"secQualifierISA03,omitempty"` // Security Information Qualifier — CMS 835 Flat File, PDF p.1
	SecInfoISA04      string            `json:"secInfoISA04,omitempty"`      // Security Information — CMS 835 Flat File, PDF p.1

	SenderIDQualifierISA05 InterchangeIDQualifierISA05 `json:"senderIDQualifierISA05,omitempty"` // Interchange ID Qualifier — CMS 835 Flat File, PDF p.1
	SenderIDISA06          string                      `json:"senderIDISA06,omitempty"`          // Interchange Sender ID — CMS 835 Flat File, PDF p.1

	ReceiverIDQualifierISA07 InterchangeIDQualifierISA05 `json:"receiverIDQualifierISA07,omitempty"` // Interchange ID Qualifier — CMS 835 Flat File, PDF p.1
	ReceiverIDISA08          string                      `json:"receiverIDISA08,omitempty"`          // Interchange Receiver ID — CMS 835 Flat File, PDF p.1

	// DateISA09 is sent as YYMMDD. The flat file's comment on ISA09 reads
	// "Format is YYMMDD with CC (20) appended", so Year is 2000 plus YY.
	DateISA09 Date `json:"dateISA09,omitempty"` // Interchange Date — CMS 835 Flat File, PDF p.1
	TimeISA10 Time `json:"timeISA10,omitempty"` // Interchange Time — CMS 835 Flat File, PDF p.1

	// RepetitionSepISA11 is the single character the interchange uses as its
	// repetition separator. It is kept as received, not interpreted here -
	// the parser is the one place that uses it as a delimiter.
	RepetitionSepISA11 string `json:"repetitionSepISA11,omitempty"` // Repetition Separator — CMS 835 Flat File, PDF p.1

	VersionISA12       string `json:"versionISA12,omitempty"`       // Interchange Control Version Number — CMS 835 Flat File, PDF p.1
	ControlNumberISA13 string `json:"controlNumberISA13,omitempty"` // Interchange Control Number — CMS 835 Flat File, PDF p.1

	AckRequestedISA14 AckRequestedISA14 `json:"ackRequestedISA14,omitempty"` // Acknowledgement Requested — CMS 835 Flat File, PDF p.1
	UsageISA15        UsageISA15        `json:"usageISA15,omitempty"`        // Usage Indicator — CMS 835 Flat File, PDF p.1

	// ComponentSepISA16 is the single character the interchange uses as its
	// component element separator, kept as received for the same reason as
	// RepetitionSepISA11.
	ComponentSepISA16 string `json:"componentSepISA16,omitempty"` // Component Element Separator — CMS 835 Flat File, PDF p.1

	FunctionalGroups []FunctionalGroupGS `json:"functionalGroups,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments between ISA and IEA, outside any GS, that the model does not place

	// Trailer: IEA. GroupCountIEA01 and ControlNumberIEA02 are envelope
	// counts, not business data; the parser checks them against the groups
	// actually present and against ControlNumberISA13.
	GroupCountIEA01    int    `json:"groupCountIEA01,omitempty"`    // # Included Functional Groups — CMS 835 Flat File, PDF p.14-15
	ControlNumberIEA02 string `json:"controlNumberIEA02,omitempty"` // Interchange Control # — CMS 835 Flat File, PDF p.14-15
}

// Parse reads one interchange of 835 transactions into the typed model in
// one pass over its segments.
//
// The delimiters come from the ISA: the element separator is the byte after
// "ISA", the component separator is ISA16, the segment terminator is the
// byte after ISA16, and the repetition separator is ISA11. Whitespace
// between segments is skipped.
//
// Loops follow the 835 viewer's rules (frontend/src/docs/viewer835/parse.ts):
// ISA holds GS, GS holds ST, and ST holds its header segments, the N1 PR
// loop, the N1 PE loop, LX loops, and PLB. An N1 loop takes the N3, N4, REF,
// PER, and RDM segments that follow it. An LX holds its claims (CLP), and a
// claim holds its service lines (SVC). A loop ends at the segment that
// starts its next sibling or an enclosing loop, at PLB, or at any envelope
// segment.
//
// Every element goes into its typed field, every position the standard
// defines for its segment. A coded value is kept as the text it carried
// whether or not it is on its list (its Known method says which). A segment
// the model has no slot for is kept whole, in file order, as a RawSegment on
// the struct that encloses it. Each segment struct records its position in
// the file (Index), and the ISA records the file's delimiters, so Write can
// put the file back together byte for byte.
//
// Parse returns an error naming the element and the segment's 0-based
// position in the file (the ISA is 0) when the model cannot hold the text: a
// malformed amount, date, time, or count; a non-empty element or component
// past the last one the standard defines for its segment; an envelope
// segment out of place; or input after the IEA. After the pass it checks each envelope
// trailer: SE01 against the segments from ST through SE, SE02 against ST02,
// GE01 against the transactions in the group, GE02 against GS06, IEA01
// against the groups, and IEA02 against ISA13.
func Parse(data []byte) (InterchangeISA, error) {
	segments, delimiters, err := splitSegments(data)
	if err != nil {
		return InterchangeISA{}, err
	}
	p := parser{segments: segments, component: delimiters.Component}
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

// Write writes an interchange back out as X12, with the file's own
// delimiters (isa.Delimiters): every segment Parse read, typed or raw, in
// the order it was read (each struct's Index), every element in its
// position, trailing empty elements and components left off the way X12
// leaves them off, the ISA at its fixed width, and the delimiters' Suffix
// after every segment. A file Parse read writes back out byte for byte,
// except that whitespace before the ISA and trailing empty elements a
// sender transmitted anyway are not reproduced. An interchange built by
// hand writes in the implementation guide's order wherever its Index
// values are all zero, and needs its Delimiters set.
func Write(w io.Writer, isa InterchangeISA) error {
	d := isa.Delimiters
	if d.Element == 0 || d.Component == 0 || d.Terminator == 0 {
		return fmt.Errorf("x12835: writing: the delimiters are not set")
	}
	x := writer{d: d}
	x.interchange(isa)
	_, err := w.Write(x.buf.Bytes())
	if err != nil {
		return fmt.Errorf("x12835: writing: %w", err)
	}
	return nil
}
