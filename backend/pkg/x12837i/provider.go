package x12837i

// ProviderPRV is the PRV segment: a provider specialty/taxonomy
// declaration, the same shape at every loop that carries one (2000A
// billing provider, 2310A attending provider). Element names and citations
// are frontend/src/data/x12Names837i.ts's PRV_EL_NAME array (source:
// Transmittal 907, PDF p.10, representative of every PRV instance). Set is
// false when the segment did not appear at all.
//
// RefQualifierPRV02 is kept as a plain string, not a coded type: its only
// observed value is "PXC" (Health Care Provider Taxonomy Code) and no PRV02
// type is in frontend/src/data/x12Names837i.ts's list of implied coded
// types.
type ProviderPRV struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	ProviderCodePRV01 ProviderRoleCodePRV01 `json:"providerCodePRV01,omitempty"` // Provider Code — Transmittal 907, PDF p.10
	RefQualifierPRV02 string                `json:"refQualifierPRV02,omitempty"` // Reference Identification Qualifier — Transmittal 907, PDF p.10 ("PXC", Health Care Provider Taxonomy Code)
	TaxonomyCodePRV03 string                `json:"taxonomyCodePRV03,omitempty"` // Provider Taxonomy Code — Transmittal 907, PDF p.10

	// ExtraElements holds PRV04 through PRV06 (State or Province Code, the
	// Provider Specialty Information composite as one text, Provider
	// Organization Code) as sent. X12 defines them; the 837I guide marks
	// them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}
