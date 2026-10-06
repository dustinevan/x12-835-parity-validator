package x12835

// Delimiters are the characters a file puts between its parts, read from
// its ISA: the element separator (the byte after "ISA"), the component
// separator (ISA16), the repetition separator (ISA11), and the segment
// terminator (the byte after ISA16). Suffix is the whitespace the file
// put after the ISA's terminator ("\n", "\r\n", or nothing), which Write
// puts after every segment, so a file that puts each segment on its own
// line writes back out that way.
type Delimiters struct {
	Element    byte   `json:"element,omitempty"`
	Component  byte   `json:"component,omitempty"`
	Repetition byte   `json:"repetition,omitempty"`
	Terminator byte   `json:"terminator,omitempty"`
	Suffix     string `json:"suffix,omitempty"`
}
