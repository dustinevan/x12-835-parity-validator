package x12835

// PartyN1 is the N1 organizational loop, used for both the payer (N1*PR)
// and the payee (N1*PE): the N1 segment itself, an optional N3 address, an
// optional N4 city/state/zip, any REF identifiers, and any PER contacts.
// Element names and citations are frontend/src/data/x12Names835.ts's
// N1_PR_EL_NAME and N1_PE_EL_NAME arrays, which name the same four elements
// slightly differently by loop (source: CMS 835 Flat File, PDF p.2 for the
// payer loop, p.3 for the payee loop); field comments give both. N105 and
// N106 are the positions the standard defines after the identifier; the
// 835 guide marks them not used, and a file that sends them keeps them.
type PartyN1 struct {
	Index            int             `json:"index,omitempty"`
	EntityIDN101     EntityIDN101    `json:"entityIDN101,omitempty"`     // Entity Identifier Code — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	NameN102         string          `json:"nameN102,omitempty"`         // Payer Name / Payee Name — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	IDQualifierN103  IDQualifierN103 `json:"idQualifierN103,omitempty"`  // ID Code Qualifier — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	IdentifierN104   string          `json:"identifierN104,omitempty"`   // Payer Identifier / Payee ID Code — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	RelationshipN105 string          `json:"relationshipN105,omitempty"` // N105 Entity Relationship Code; not used by the 835 guide, kept as sent
	EntityIDN106     string          `json:"entityIDN106,omitempty"`     // N106 Entity Identifier Code; not used by the 835 guide, kept as sent

	Address  AddressN3  `json:"address,omitempty"`  // N3, optional
	Location LocationN4 `json:"location,omitempty"` // N4, optional

	References []ReferenceREF `json:"references,omitempty"` // REF — CMS 835 Flat File, PDF p.3 (payer), p.4 (payee)
	Contacts   []ContactPER   `json:"contacts,omitempty"`   // PER — CMS 835 Flat File, PDF p.3; the payee loop carries none in these element definitions

	// DeliveryRDM is the RDM delivery-method segment. It belongs to the
	// payee loop only; a PartyN1 built for the payer leaves it zero.
	DeliveryRDM DeliveryRDM `json:"deliveryRDM,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // a repeat of N3, N4, or RDM within the loop, and any segment the loop has no slot for
}

// AddressN3 is the N3 segment, shared by the payer and payee loops. Element
// names and citations are frontend/src/data/x12Names835.ts's N3_PR_EL_NAME
// and N3_PE_EL_NAME arrays (source: CMS 835 Flat File, PDF p.2 payer, p.3
// payee). Set is false when the segment did not appear at all.
type AddressN3 struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	Line1N301 string `json:"line1N301,omitempty"` // Payer Address Line / Payee Address Line — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	Line2N302 string `json:"line2N302,omitempty"` // Payer Address Line / Payee Address Line — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
}

// LocationN4 is the N4 segment, shared by the payer and payee loops.
// Element names and citations are frontend/src/data/x12Names835.ts's
// N4_PR_EL_NAME and N4_PE_EL_NAME arrays (source: CMS 835 Flat File, PDF
// p.2 payer, p.3-4 payee). N405 to N407 are the positions the standard
// defines after the country code; the 835 guide uses N407 for a country
// subdivision and marks N405 and N406 not used, and a file that sends any
// of them keeps it. Set is false when the segment did not appear at all.
type LocationN4 struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	CityN401    string `json:"cityN401,omitempty"`    // Payer City Name / Payee City Name — CMS 835 Flat File, PDF p.2 (payer), p.3 (payee)
	StateN402   string `json:"stateN402,omitempty"`   // Payer State Code / Payee State Code — CMS 835 Flat File, PDF p.2 (payer), p.4 (payee)
	ZipN403     string `json:"zipN403,omitempty"`     // Payer Postal Zone or ZIP Code / Payee Postal Zone or ZIP Code — CMS 835 Flat File, PDF p.2 (payer), p.4 (payee)
	CountryN404 string `json:"countryN404,omitempty"` // Country Code — CMS 835 Flat File, PDF p.2 (payer; usage S)

	LocationQualifierN405  string `json:"locationQualifierN405,omitempty"`  // N405 Location Qualifier; kept as sent
	LocationN406           string `json:"locationN406,omitempty"`           // N406 Location Identifier; kept as sent
	CountrySubdivisionN407 string `json:"countrySubdivisionN407,omitempty"` // N407 Country Subdivision Code; kept as sent
}

// ContactPER is the PER contact segment. It appears in the payer loop
// (1000A) and the claim loop (2100). PER01 and the three communication
// qualifier positions carry the same element name in every instance;
// PER02's own label (whose name it is) and PER08's (a communication number
// in the payer loop, an extension in the claim loop) vary by loop - see
// PAYER_PER_NAME/CLAIM_PER_NAME and PER_PAYER_EL_NAME/PER_CLAIM_EL_NAME in
// frontend/src/data/x12Names835.ts for the per-loop text (source: CMS 835
// Flat File, PDF p.3 payer, p.11 claim). PER09 is the position the
// standard defines after the third number; the 835 guide marks it not
// used, and a file that sends it keeps it.
type ContactPER struct {
	Index         int                  `json:"index,omitempty"`
	FunctionPER01 ContactFunctionPER01 `json:"functionPER01,omitempty"` // Contact Function Code — CMS 835 Flat File, PDF p.3, 11

	NamePER02 string `json:"namePER02,omitempty"` // contact name — CMS 835 Flat File, PDF p.3 (payer), p.11 (claim)

	QualifierPER03 CommQualifierPER03 `json:"qualifierPER03,omitempty"` // Communication # Qualifier — CMS 835 Flat File, PDF p.3, 11
	CommPER04      string             `json:"commPER04,omitempty"`      // Communication Number — CMS 835 Flat File, PDF p.3, 11

	QualifierPER05 CommQualifierPER03 `json:"qualifierPER05,omitempty"` // Communication # Qualifier — CMS 835 Flat File, PDF p.3, 11
	CommPER06      string             `json:"commPER06,omitempty"`      // Communication Number — CMS 835 Flat File, PDF p.3, 11

	QualifierPER07 CommQualifierPER03 `json:"qualifierPER07,omitempty"` // Communication # Qualifier — CMS 835 Flat File, PDF p.3, 11
	CommPER08      string             `json:"commPER08,omitempty"`      // Communication Number (payer loop) / Communication # Extension (claim loop) — CMS 835 Flat File, PDF p.3, 11

	InquiryReferencePER09 string `json:"inquiryReferencePER09,omitempty"` // PER09 Contact Inquiry Reference; kept as sent
}

// DeliveryRDM is the RDM segment: how the payee receives remittance detail.
// It belongs to the payee loop only. Element names and citations are
// frontend/src/data/x12Names835.ts's RDM_EL_NAME array (source: CMS 835
// Flat File, PDF p.4). RDM04 and RDM05 are the positions the standard
// defines after the communication number; the 835 guide marks them not
// used, and a file that sends them keeps them. Set is false when the
// segment did not appear at all.
type DeliveryRDM struct {
	Index int  `json:"index,omitempty"`
	Set   bool `json:"set,omitempty"`

	ReportTransmissionRDM01 DeliveryMethodRDM01 `json:"reportTransmissionRDM01,omitempty"` // Report Transmission Code — CMS 835 Flat File, PDF p.4
	NameRDM02               string              `json:"nameRDM02,omitempty"`               // Name — CMS 835 Flat File, PDF p.4
	CommNumberRDM03         string              `json:"commNumberRDM03,omitempty"`         // Communication Number — CMS 835 Flat File, PDF p.4
	ReferenceRDM04          string              `json:"referenceRDM04,omitempty"`          // RDM04 Reference Identifier (composite C040); kept as sent with its component separators
	ReferenceRDM05          string              `json:"referenceRDM05,omitempty"`          // RDM05 Reference Identifier (composite C040); kept as sent with its component separators
}
