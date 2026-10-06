package x12835

import "fmt"

// ParseClaim reads one claim's segments into a ClaimCLP: the bytes from
// the claim's CLP through the last segment of its last service line, as the
// claim table's segments column holds them, with the file's four delimiters,
// since the bytes carry no ISA. The segments go through the same claim
// parsing as Parse, with the same Raw behavior; every Index, a RawSegment's
// included, counts from the CLP, which is 0. Input that does not start with
// CLP, or that continues past the claim (a second CLP, an LX, a PLB, or an
// envelope segment), is an error.
func ParseClaim(data []byte, element, component, repetition, terminator byte) (ClaimCLP, error) {
	if element == component || element == repetition || element == terminator ||
		component == repetition || component == terminator || repetition == terminator {
		return ClaimCLP{}, fmt.Errorf("x12835: parsing claim: the four delimiters must differ")
	}
	segments, err := appendSegments(nil, string(data), 0, element, terminator)
	if err != nil {
		return ClaimCLP{}, err
	}
	if len(segments) == 0 || segments[0].id != "CLP" {
		return ClaimCLP{}, fmt.Errorf("x12835: parsing claim: the segments do not start with CLP")
	}
	p := parser{segments: segments, component: component}
	claim, err := p.parseClaim()
	if err != nil {
		return ClaimCLP{}, err
	}
	if next := p.peek(); next.id != "" {
		return ClaimCLP{}, fmt.Errorf("x12835: parsing claim: segment %d (%s) is outside the claim", next.index, next.id)
	}
	return claim, nil
}
