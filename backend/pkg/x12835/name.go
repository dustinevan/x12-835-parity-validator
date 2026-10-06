package x12835

// NameNM1 holds one NM1 segment from the claim loop (2100). The same shape
// is reused for every role the 835 sends there — patient, corrected
// patient/insured, service provider (rendering), crossover carrier, and
// corrected priority payer — and NM103-09 carry a different public name
// under each role (NM1_QC_EL_NAME, NM1_74_EL_NAME, NM1_82_EL_NAME,
// NM1_TT_EL_NAME, NM1_PR_EL_NAME in frontend/src/data/x12Names835.ts). The
// qualifier that selects which role applies is carried by the field holding
// this struct (PatientNM1QC, CorrectedPatientNM174, RenderingNM182,
// CrossoverNM1TT, CorrectedPayerNM1PR in claim.go), not by a field here.
// NM110 to NM112 are the positions the standard defines after the
// identifier; the 835 guide marks them not used, and a file that sends them
// keeps them.
type NameNM1 struct {
	Index             int              `json:"index,omitempty"`
	EntityIDNM101     EntityIDNM101    `json:"entityIDNM101,omitempty"`     // Entity Identifier Code; CMS 835 Flat File, PDF p.6-8
	EntityTypeNM102   EntityTypeNM102  `json:"entityTypeNM102,omitempty"`   // Entity Type Qualifier; CMS 835 Flat File, PDF p.6-8
	LastNM103         string           `json:"lastNM103,omitempty"`         // last name, or organization name, depending on the role (NM102)
	FirstNM104        string           `json:"firstNM104,omitempty"`        // first name
	MiddleNM105       string           `json:"middleNM105,omitempty"`       // middle name
	PrefixNM106       string           `json:"prefixNM106,omitempty"`       // name prefix
	SuffixNM107       string           `json:"suffixNM107,omitempty"`       // name suffix
	IDQualifierNM108  IDQualifierNM108 `json:"idQualifierNM108,omitempty"`  // ID Code Qualifier
	IDNM109           string           `json:"idNM109,omitempty"`           // the identifier itself (NPI, HICN/MBI, ...); kept as text even when all digits
	RelationshipNM110 string           `json:"relationshipNM110,omitempty"` // NM110 Entity Relationship Code; not used by the 835 guide, kept as sent
	EntityIDNM111     string           `json:"entityIDNM111,omitempty"`     // NM111 Entity Identifier Code; not used by the 835 guide, kept as sent
	LastNM112         string           `json:"lastNM112,omitempty"`         // NM112 Name Last or Organization Name; not used by the 835 guide, kept as sent
}
