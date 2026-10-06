package x12837i

// NameNM1 holds one generic NM1 segment, the shape shared by every loop in
// this transaction that carries one: NM101-02 (role and entity type),
// NM103-07 (name components, with role-specific public names - a role's own
// field-level citations live on the dedicated loop struct that embeds this
// shape, e.g. SubmitterNM141, ReceiverNM140, and the per-role loops named
// by frontend/src/data/x12Names837i.ts's per-loop *_EL_NAME arrays, such as
// NM1_ATTENDING_EL_NAME), and NM108-09 (the identifier and
// its qualifier). Copied in shape from backend/pkg/x12835/name.go, whose
// NameNM1 is reused the same way for every claim-loop NM1 role.
//
// The role that selects which loop this occurrence belongs to is
// EntityIDNM101 itself here (unlike x12835's NameNM1, which leaves the role
// to be implied by the field name holding the struct): with 14 distinct
// NM101 codes in this transaction, against x12835's handful, the role is
// recorded directly rather than only implied.
type NameNM1 struct {
	Index int `json:"index,omitempty"` // the segment's position in the file

	EntityIDNM101   PartyRoleNM101  `json:"entityIDNM101,omitempty"`   // Entity Identifier Code — Transmittal 907, role varies by loop; see PartyRoleNM101
	EntityTypeNM102 EntityTypeNM102 `json:"entityTypeNM102,omitempty"` // Entity Type Qualifier — Transmittal 907, PDF p.9

	LastNM103   string `json:"lastNM103,omitempty"`   // last name, or organization name, depending on the role (NM102); public name varies by loop
	FirstNM104  string `json:"firstNM104,omitempty"`  // first name
	MiddleNM105 string `json:"middleNM105,omitempty"` // middle name
	PrefixNM106 string `json:"prefixNM106,omitempty"` // name prefix
	SuffixNM107 string `json:"suffixNM107,omitempty"` // name suffix

	IDQualifierNM108 IDQualifierNM108 `json:"idQualifierNM108,omitempty"` // ID Code Qualifier
	IDNM109          string           `json:"idNM109,omitempty"`          // the identifier itself (NPI, MBI, ...); kept as text even when all digits

	// ExtraElements holds NM110 through NM112 (Entity Relationship Code,
	// Entity Identifier Code, Name Last or Organization Name) as sent, in
	// order, trailing empties dropped. X12 defines them; the 837I guide
	// marks them Not Used; a file that sends them loses nothing.
	ExtraElements []string `json:"extraElements,omitempty"`
}
