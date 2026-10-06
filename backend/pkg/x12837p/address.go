package x12837p

// AddressN3 is the N3 segment, shared by every address-carrying loop in the
// 837P (billing provider, pay-to provider/plan, subscriber, payer, service
// facility location, ambulance pick-up/drop-off, other subscriber/payer, ...).
// Element names and citations are frontend/src/data/x12Names837p.ts's
// N3_EL_NAME array (source: Transmittal 907, every N3 instance) and
// LOOP_N3_NAME for each loop's own box title. Set is false when the segment
// did not appear at all.
type AddressN3 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	Line1N301 string `json:"line1N301,omitempty"` // Address Line — Transmittal 907, every N3 instance
	Line2N302 string `json:"line2N302,omitempty"` // Address Line — Transmittal 907, every N3 instance

	Extra []string `json:"extra,omitempty"` // elements after N302, kept as sent
}

// LocationN4 is the N4 segment, shared by every address-carrying loop in
// the 837P (see AddressN3's doc comment for the loop list). Element names
// and citations are frontend/src/data/x12Names837p.ts's N4_EL_NAME array
// (source: Transmittal 907, every N4 instance; N404-407 are usually Not
// Used except N404/N407 situational) and LOOP_N4_NAME for each loop's own
// box title. N405 and N406 are Not Used and kept as plain text. Set is
// false when the segment did not appear at all.
type LocationN4 struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	CityN401    string `json:"cityN401,omitempty"`    // City Name
	StateN402   string `json:"stateN402,omitempty"`   // State or Province Code
	ZipN403     string `json:"zipN403,omitempty"`     // Postal Zone or ZIP Code
	CountryN404 string `json:"countryN404,omitempty"` // Country Code; situational

	LocationQualifierN405  string `json:"locationQualifierN405,omitempty"`  // Location Qualifier (N/U) — Transmittal 907, PDF p.139 (2010AA N4)
	LocationIdentifierN406 string `json:"locationIdentifierN406,omitempty"` // Location Identifier (N/U) — Transmittal 907, PDF p.139 (2010AA N4)

	SubdivisionN407 string `json:"subdivisionN407,omitempty"` // Country Subdivision Code; situational — Transmittal 907, PDF p.139 (2010AA N4)

	Extra []string `json:"extra,omitempty"` // elements after N407, kept as sent
}
