package x12837p

// ProviderPRV holds one PRV segment: the provider-specialty information
// carried at the billing provider (2000A), rendering provider (2310B,
// 2420A) loops. PRV_EL_NAME and LOOP_PRV_NAME in
// frontend/src/data/x12Names837p.ts give the position names (source:
// Transmittal 907, every PRV instance) and each loop's own box title;
// PRV04-06 are Not Used at every instance, and a file that sends them keeps
// their text in Extra.
//
// QualifierPRV02 reuses ReferenceIdentificationQualifierREF01
// (enums_reference.go): PRV02 is X12 element 128, the same Reference
// Identification Qualifier code list REF01 carries, at a different segment
// position - not a separately sourced code list.
//
// ProviderTypePRV01 is defined in enums_provider.go.
type ProviderPRV struct {
	Index int `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	CodePRV01      ProviderTypePRV01                     `json:"codePRV01,omitempty"`      // Provider Code
	QualifierPRV02 ReferenceIdentificationQualifierREF01 `json:"qualifierPRV02,omitempty"` // Reference Identification Qualifier
	TaxonomyPRV03  string                                `json:"taxonomyPRV03,omitempty"`  // Provider Taxonomy Code

	Extra []string `json:"extra,omitempty"` // elements after PRV03, kept as sent
}
