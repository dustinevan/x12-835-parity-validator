package x12837p

// NameNM1 holds one NM1 segment. The same shape is reused for every role
// the 837P sends an NM1 for - submitter, receiver, billing/pay-to provider,
// subscriber, payer, patient, every claim-level and line-level provider
// role (2310A-F, 2420A-H), and every other-subscriber/other-payer role
// (2330A-G) - and NM103-09 carry a different public name under each role;
// LOOP_NM1_NAME and NM1_ROLE_EL_NAME in frontend/src/data/x12Names837p.ts
// give the per-loop box title and element names (source: Transmittal 907,
// PDF p.133 onward). The loop that selects which role applies is carried by
// the field holding this struct, not by a field here.
//
// NM106 (Name Prefix) is Not Used at every 837P NM1 instance
// (frontend/src/data/x12Names837p.ts's NM1_ROLE_EL_NAME arrays skip it), but
// it sits between positions the model holds, so it keeps a field, as in
// backend/pkg/x12835's NameNM1: a file that sends it keeps the text. NM110-112
// are Not Used too and come after NM109, so a file that sends them keeps
// their text in Extra.
//
// EntityIdentifierNM101, EntityTypeQualifierNM102, and
// IdentificationCodeQualifierNM108 are defined in enums_envelope.go.
type NameNM1 struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	EntityIDNM101    EntityIdentifierNM101            `json:"entityIDNM101,omitempty"`    // Entity Identifier Code — Transmittal 907, every NM1 instance; NM101_EL_NAME
	EntityTypeNM102  EntityTypeQualifierNM102         `json:"entityTypeNM102,omitempty"`  // Entity Type Qualifier — Transmittal 907, every NM1 instance; NM102_EL_NAME
	LastNM103        string                           `json:"lastNM103,omitempty"`        // last name, or organization name, depending on the role (NM102)
	FirstNM104       string                           `json:"firstNM104,omitempty"`       // first name
	MiddleNM105      string                           `json:"middleNM105,omitempty"`      // middle name or initial
	PrefixNM106      string                           `json:"prefixNM106,omitempty"`      // Name Prefix; Not Used at every 837P instance
	SuffixNM107      string                           `json:"suffixNM107,omitempty"`      // name suffix
	IDQualifierNM108 IdentificationCodeQualifierNM108 `json:"idQualifierNM108,omitempty"` // Identification Code Qualifier
	IDNM109          string                           `json:"idNM109,omitempty"`          // the identifier itself (NPI, EDI Transmitter ID, Member ID, ...); kept as text even when all digits

	Extra []string `json:"extra,omitempty"` // elements after NM109 (NM110-112 and beyond), kept as sent
}
