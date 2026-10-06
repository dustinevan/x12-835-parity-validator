package x12837i

// AddressN3 is the N3 segment, the same two-line address shape at every
// loop that carries one in this transaction. Element names and citations
// are frontend/src/data/x12Names837i.ts's N3_EL_NAME array (source:
// Transmittal 907, PDF p.12, representative of every N3 instance; the box
// caption's role word, e.g. "Billing Provider Address", varies by loop and
// is not repeated here). Set is false when the segment did not appear at
// all. Copied in shape from backend/pkg/x12835/party.go's AddressN3.
type AddressN3 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	Line1N301 string `json:"line1N301,omitempty"` // Address Line — Transmittal 907, PDF p.12
	Line2N302 string `json:"line2N302,omitempty"` // Address Line — Transmittal 907, PDF p.12
}

// LocationN4 is the N4 segment, the same shape at every loop that carries
// one. Element names and citations are frontend/src/data/x12Names837i.ts's
// N4_EL_NAME array (source: Transmittal 907, PDF p.12, representative of
// every N4 instance). N405 and N406 are Not Used throughout this
// transaction per the flat file; N407 (Country Subdivision Code) is kept,
// unlike backend/pkg/x12835's LocationN4, which has no equivalent position.
// Set is false when the segment did not appear at all.
type LocationN4 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	CityN401    string `json:"cityN401,omitempty"`    // City Name — Transmittal 907, PDF p.12
	StateN402   string `json:"stateN402,omitempty"`   // State or Province Code — Transmittal 907, PDF p.12
	ZipN403     string `json:"zipN403,omitempty"`     // Postal Zone or ZIP Code — Transmittal 907, PDF p.12
	CountryN404 string `json:"countryN404,omitempty"` // Country Code — Transmittal 907, PDF p.12

	LocationQualifierN405  string `json:"locationQualifierN405,omitempty"`  // Location Qualifier (N/U) — Transmittal 907, PDF p.12
	LocationIdentifierN406 string `json:"locationIdentifierN406,omitempty"` // Location Identifier (N/U) — Transmittal 907, PDF p.12
	SubdivisionN407        string `json:"subdivisionN407,omitempty"`        // Country Subdivision Code — Transmittal 907, PDF p.12
}
