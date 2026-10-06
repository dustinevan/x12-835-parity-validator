package x12837p

import "fmt"

// The coded-element scheme
//
// Every coded element in this package (EntityIdentifierNM101,
// PayerResponsibilitySBR01, ClaimFrequencyCLM05_3, and the rest of
// enums_envelope.go, enums_provider.go, enums_subscriber.go, enums_claim.go,
// enums_service.go, enums_reference.go, enums_cob.go, and the envelope
// types beside their segments) is a string type whose value is the code
// text exactly as the file sent it:
//
//   - The zero value, "", is the one empty value: the element carried no
//     text, either nothing between its separators or the segment ending
//     before it. X12 gives those two the same meaning. Each type names it
//     <Type>NotSent.
//   - Every code the X12 005010X222A1 implementation guide allows for the
//     element is a named constant whose value is its code text, with the
//     guide's name and a citation in the type's list. A code Medicare's
//     own documents restrict or confirm keeps that note in its citation.
//   - A code off the list is still the value: the text is kept as it was
//     sent, String returns it, and the JSON form carries it, so a file
//     reads in and writes out without loss. Known reports whether the code
//     is on the list; ShortDescription and Description return "" for one
//     that is not.
//
// One list beside each type (a codeList) is the one place its text lives:
// the code each constant parses from, its public short name ("" when the
// sources give none), and that name with its citation. Parse<Type> returns
// the text as the type and an error when the text is neither empty nor on
// the list; the caller keeps the value either way. A string type needs no
// MarshalText or MarshalJSON of its own: encoding/json writes and reads
// the code text.
//
// Where a list is Medicare's rather than the guide's, the type's comment
// says so. An element whose list is an external code source the repository
// does not hold (place of service codes, taxonomy codes, the NUBC lists)
// stays a plain string.
//
// CARC and RARC (CASReason.Reason, the MOA remark positions) are coded
// elements too, but their types are carc.Code and rarc.Code
// (backend/pkg/carc, backend/pkg/rarc), whose entries tables tools/codegen
// updates from the published lists; both are string types with the same
// shape.

// codeEntry is one code of a coded element's list: the code text, its
// public short name ("" when the sources give none), and that name with
// its citation.
type codeEntry struct {
	code        string
	short       string
	description string
}

// codeList is the allowed codes of one coded element, in the order the
// citation lists them. Lists are short, so a lookup reads them in order.
type codeList []codeEntry

func (l codeList) find(code string) (codeEntry, bool) {
	for _, e := range l {
		if e.code == code {
			return e, true
		}
	}
	return codeEntry{}, false
}

func (l codeList) known(code string) bool {
	_, ok := l.find(code)
	return ok
}

// short returns the public name of a code: "not sent" for the empty value,
// the list's short name for a listed code, and "" for a code off the list
// or a listed code the sources give no name.
func (l codeList) short(code string) string {
	if code == "" {
		return "not sent"
	}
	e, _ := l.find(code)
	return e.short
}

// description returns the name of a code with its citation: a fixed note
// for the empty value, the list's text for a listed code, and "" for a
// code off the list.
func (l codeList) description(code string) string {
	if code == "" {
		return "the element carried no value"
	}
	e, _ := l.find(code)
	return e.description
}

// check returns nil when code is empty or on the list, and otherwise an
// error naming the element.
func (l codeList) check(element, code string) error {
	if code == "" || l.known(code) {
		return nil
	}
	return fmt.Errorf("x12837p: parsing %s: a value not on the list", element)
}

// tr3 is the citation for a code or a name transcribed from the 837P
// implementation guide, which is not in the repository.
const tr3 = "X12 005010X222A1 TR3 (transcribed from the standard; verify against the printed guide)"

// guide returns the description of a code the implementation guide allows
// at element: the guide's name for it, then the citation.
func guide(name, element string) string {
	return name + " — " + tr3 + ", " + element
}
