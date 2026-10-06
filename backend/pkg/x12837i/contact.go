package x12837i

// ContactPER is the PER contact segment: identical shape at every instance
// in this transaction (submitter 1000A, billing provider 2010AA). PER01 is
// always "Contact Function Code" (value "IC", Information Contact); PER02
// onward repeats name/communication-number triplets. Element names and
// citations are frontend/src/data/x12Names837i.ts's PER_EL_NAME array
// (source: Transmittal 907, PDF p.9-10, representative of every PER
// instance). Copied in shape from backend/pkg/x12835/party.go's ContactPER.
//
// Unlike x12835's ContactPER, FunctionPER01 and the three communication
// qualifier positions are kept as plain strings, not coded types:
// frontend/src/data/x12Names837i.ts's own list of the coded types this
// file's data implies does not include a PER01 or PER-qualifier type (PER01
// has one observed value, "IC"; PER_QUALIFIER_VALUES lists EM/EX/FX/TE with
// no per-code public names in either reviewed source).
type ContactPER struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	FunctionPER01 string `json:"functionPER01,omitempty"` // Contact Function Code — Transmittal 907, PDF p.9-10 ("IC", Information Contact)

	NamePER02 string `json:"namePER02,omitempty"` // contact name — Transmittal 907, PDF p.9-10

	QualifierPER03 string `json:"qualifierPER03,omitempty"` // Communication Number Qualifier — Transmittal 907, PDF p.9-10
	CommPER04      string `json:"commPER04,omitempty"`      // Communication Number — Transmittal 907, PDF p.9-10

	QualifierPER05 string `json:"qualifierPER05,omitempty"` // Communication Number Qualifier — Transmittal 907, PDF p.9-10
	CommPER06      string `json:"commPER06,omitempty"`      // Communication Number — Transmittal 907, PDF p.9-10

	QualifierPER07 string `json:"qualifierPER07,omitempty"` // Communication Number Qualifier — Transmittal 907, PDF p.9-10
	CommPER08      string `json:"commPER08,omitempty"`      // Communication Number — Transmittal 907, PDF p.9-10

	// ExtraElements holds PER09 (Contact Inquiry Reference) as sent. X12
	// defines it; the 837I guide marks it Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}
