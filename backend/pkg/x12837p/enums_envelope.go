package x12837p

// Submitter/receiver-level coded elements (1000A Submitter, 1000B Receiver):
// the NM1 identity qualifiers shared by every NM1 segment in the
// transaction. The ISA, GS, ST, and BHT codes sit beside their segments in
// interchange.go, functional_group.go, transaction.go, and header.go.

// EntityIdentifierNM101 is NM101, the Entity Identifier Code on an NM1 segment,
// shared by every NM1 role across the transaction (submitter, receiver, billing
// provider, subscriber, ordering provider, and the rest).
//
// The list is every NM101 code the 837P implementation guide allows at one of
// its NM1 loops. Transmittal 907 names every NM1 instance's loop/box but,
// unlike the 835 flat file, does not restate NM101's own code character at
// most instances (frontend/src/data/x12Names837p.ts's header note). The first
// five values below are the ones the CEDI 837P Companion Guide's own worked
// transaction shows (section 11.2, Figure 3). The rest come from Transmittal
// 579, the original 005010 837P flat file, whose Values column prints NM101's
// code at every NM1 box; the parser places each NM1 loop by this code. A box
// whose Values column lists one code gives that code its box's role as its
// short name; the Referring Provider boxes list two codes (DN, P3) without
// saying which is which, so those two take the guide's own names.
type EntityIdentifierNM101 string

const (
	EntityIdentifierNM101NotSent EntityIdentifierNM101 = ""
	EntityIdentifierNM10141      EntityIdentifierNM101 = "41"
	EntityIdentifierNM10140      EntityIdentifierNM101 = "40"
	EntityIdentifierNM10185      EntityIdentifierNM101 = "85"
	EntityIdentifierNM101IL      EntityIdentifierNM101 = "IL"
	EntityIdentifierNM101DK      EntityIdentifierNM101 = "DK"
	EntityIdentifierNM10187      EntityIdentifierNM101 = "87"
	EntityIdentifierNM101PE      EntityIdentifierNM101 = "PE"
	EntityIdentifierNM101PR      EntityIdentifierNM101 = "PR"
	EntityIdentifierNM101QC      EntityIdentifierNM101 = "QC"
	EntityIdentifierNM101DN      EntityIdentifierNM101 = "DN"
	EntityIdentifierNM101P3      EntityIdentifierNM101 = "P3"
	EntityIdentifierNM10182      EntityIdentifierNM101 = "82"
	EntityIdentifierNM10177      EntityIdentifierNM101 = "77"
	EntityIdentifierNM101DQ      EntityIdentifierNM101 = "DQ"
	EntityIdentifierNM101PW      EntityIdentifierNM101 = "PW"
	EntityIdentifierNM10145      EntityIdentifierNM101 = "45"
	EntityIdentifierNM101QB      EntityIdentifierNM101 = "QB"
)

var entityIdentifierNM101List = codeList{
	{"41", "Submitter", "Submitter — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"NM1*41*2*EDIT TEST SUBMITTER\"); Transmittal 579, PDF p.8 (1000A); frontend/src/data/x12Names837p.ts NM101_CONFIRMED_NAME"},
	{"40", "Receiver", "Receiver — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"NM1*40*2*DME MAC JURISDICTION B\"); Transmittal 579, PDF p.9 (1000B); frontend/src/data/x12Names837p.ts NM101_CONFIRMED_NAME"},
	{"85", "Billing Provider", "Billing Provider — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"NM1*85*2*BILLING PROVIDER\"); Transmittal 579, PDF p.11 (2010AA) and p.77 (2330G, Other Payer Billing Provider); frontend/src/data/x12Names837p.ts NM101_CONFIRMED_NAME"},
	{"IL", "Insured or Subscriber", "Insured or Subscriber — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"NM1*IL*1*LASTNAME*FIRSTNAME\"); Transmittal 579, PDF p.17 (2010BA) and p.68 (2330A, Other Subscriber); same code the 835 side of this project does not name (the 835 uses N1/NM1 PR/PE/QC instead); frontend/src/data/x12Names837p.ts NM101_CONFIRMED_NAME"},
	{"DK", "Ordering Provider", "Ordering Provider — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, \"NM1*DK*1*PROVIDER*ORDERING\"); Transmittal 579, PDF p.102 (2420E); frontend/src/data/x12Names837p.ts NM101_CONFIRMED_NAME"},
	{"87", "Pay-to Provider", "Pay-to Provider — the 2010AB Pay-to Address Name box's only NM101 value, Transmittal 579, PDF p.13 (Values column)"},
	{"PE", "Pay-to Plan", "Pay-to Plan — the 2010AC Pay-to Plan Name box's only NM101 value, Transmittal 579, PDF p.14 (Values column); the loop is not used by Medicare"},
	{"PR", "Payer", "Payer — the only NM101 value of the 2010BB Payer Name box, Transmittal 579, PDF p.20, and of the 2330B Other Payer Name box, PDF p.69 (Values column)"},
	{"QC", "Patient", "Patient — the 2010CA Patient Name box's only NM101 value, Transmittal 579, PDF p.23 (Values column); the loop is not used by Medicare"},
	{"DN", "Referring Provider", "Referring Provider — one of the two NM101 values (DN, P3) of the Referring Provider Name boxes, 2310A, 2330C, and 2420F, Transmittal 579, PDF p.57, 73, 104 (Values column), which does not say which code is which; the name is the guide's: " + tr3},
	{"P3", "Primary Care Provider", "Primary Care Provider — one of the two NM101 values (DN, P3) of the Referring Provider Name boxes, 2310A, 2330C, and 2420F, Transmittal 579, PDF p.57, 73, 104 (Values column), which does not say which code is which; the name is the guide's: " + tr3},
	{"82", "Rendering Provider", "Rendering Provider — the only NM101 value of the 2310B and 2420A Rendering Provider Name boxes and the 2330D Other Payer Rendering Provider box, Transmittal 579, PDF p.58, 97, 74 (Values column)"},
	{"77", "Service Facility Location", "Service Facility Location — the only NM101 value of the 2310C and 2420C Service Facility Location Name boxes and the 2330E Other Payer Service Facility Location box, Transmittal 579, PDF p.59, 100, 75 (Values column)"},
	{"DQ", "Supervising Provider", "Supervising Provider — the only NM101 value of the 2310D and 2420D Supervising Provider Name boxes and the 2330F Other Payer Supervising Provider box, Transmittal 579, PDF p.61, 101, 76 (Values column)"},
	{"PW", "Ambulance Pick Up Location", "Ambulance Pick Up Location — the only NM101 value of the 2310E and 2420G Ambulance Pick Up Location boxes, Transmittal 579, PDF p.62, 105 (Values column)"},
	{"45", "Ambulance Drop Off Location", "Ambulance Drop Off Location — the only NM101 value of the 2310F and 2420H Ambulance Drop Off Location boxes, Transmittal 579, PDF p.63, 106 (Values column)"},
	{"QB", "Purchased Service Provider", "Purchased Service Provider — the 2420B Purchased Service Provider Name box's only NM101 value, Transmittal 579, PDF p.99 (Values column)"},
}

func (v EntityIdentifierNM101) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v EntityIdentifierNM101) Known() bool { return entityIdentifierNM101List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v EntityIdentifierNM101) ShortDescription() string {
	return entityIdentifierNM101List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v EntityIdentifierNM101) Description() string {
	return entityIdentifierNM101List.description(string(v))
}

// ParseEntityIdentifierNM101 returns text as an EntityIdentifierNM101, with an
// error when the text is neither empty nor on the list. The value is kept
// either way.
func ParseEntityIdentifierNM101(text string) (EntityIdentifierNM101, error) {
	return EntityIdentifierNM101(text), entityIdentifierNM101List.check("NM101", text)
}

// EntityTypeQualifierNM102 is NM102, the Entity Type Qualifier on an NM1
// segment (1 = person, 2 = non-person entity, by the segment's own structure),
// shared by every NM1 instance.
//
// Neither Transmittal 907 nor the CEDI 837P Companion Guide names NM102's
// individual codes in text. Transmittal 579's Values column lists "1, 2"
// (PDF p.8, 1000A), and the CEDI worked transaction sends 1 with a person's
// name and 2 with an organization's. The names are the guide's.
type EntityTypeQualifierNM102 string

const (
	EntityTypeQualifierNM102NotSent EntityTypeQualifierNM102 = ""
	EntityTypeQualifierNM1021       EntityTypeQualifierNM102 = "1"
	EntityTypeQualifierNM1022       EntityTypeQualifierNM102 = "2"
)

var entityTypeQualifierNM102List = codeList{
	{"1", "Person", `Person — Transmittal 579, PDF p.8 (Values column, "1, 2"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3) sends it with a person's name, "NM1*IL*1*LASTNAME*FIRSTNAME"; the name is the guide's: ` + tr3},
	{"2", "Non-Person Entity", `Non-Person Entity — Transmittal 579, PDF p.8 (Values column, "1, 2"); CEDI 837P CG, PDF p.50 (printed p.42; Figure 3) sends it with an organization's name, "NM1*85*2*BILLING PROVIDER"; the name is the guide's: ` + tr3},
}

func (v EntityTypeQualifierNM102) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v EntityTypeQualifierNM102) Known() bool { return entityTypeQualifierNM102List.known(string(v)) }

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v EntityTypeQualifierNM102) ShortDescription() string {
	return entityTypeQualifierNM102List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v EntityTypeQualifierNM102) Description() string {
	return entityTypeQualifierNM102List.description(string(v))
}

// ParseEntityTypeQualifierNM102 returns text as an EntityTypeQualifierNM102,
// with an error when the text is neither empty nor on the list. The value
// is kept either way.
func ParseEntityTypeQualifierNM102(text string) (EntityTypeQualifierNM102, error) {
	return EntityTypeQualifierNM102(text), entityTypeQualifierNM102List.check("NM102", text)
}

// IdentificationCodeQualifierNM108 is NM108, the Identification Code Qualifier
// on an NM1 segment, shared by every NM1 instance that carries an identifier.
//
// The list is every NM108 code the 837P implementation guide allows at one of
// its NM1 loops: 46 for the submitter and receiver, XX for every provider, MI
// and II for the subscriber and patient, and PI and XV for a payer or a
// pay-to plan. The four the CEDI 837P Companion Guide's own worked
// transaction and Medicare-restriction tables confirm keep those citations
// (frontend/src/data/x12Names837p.ts NM108_CONFIRMED_NAME); II and XV are the
// two further codes Transmittal 579's Values column lists (II beside MI at
// 2010BA and 2330A, XV beside PI at 2010AC, 2010BB, and 2330B), named here
// from the guide.
type IdentificationCodeQualifierNM108 string

const (
	IdentificationCodeQualifierNM108NotSent IdentificationCodeQualifierNM108 = ""
	IdentificationCodeQualifierNM10846      IdentificationCodeQualifierNM108 = "46"
	IdentificationCodeQualifierNM108XX      IdentificationCodeQualifierNM108 = "XX"
	IdentificationCodeQualifierNM108MI      IdentificationCodeQualifierNM108 = "MI"
	IdentificationCodeQualifierNM108PI      IdentificationCodeQualifierNM108 = "PI"
	IdentificationCodeQualifierNM108II      IdentificationCodeQualifierNM108 = "II"
	IdentificationCodeQualifierNM108XV      IdentificationCodeQualifierNM108 = "XV"
)

var identificationCodeQualifierNM108List = codeList{
	{"46", "Electronic Transmitter Identification Number", "Electronic Transmitter Identification Number — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, submitter/receiver \"*****46*B08XXXXXX\"); frontend/src/data/x12Names837p.ts NM108_CONFIRMED_NAME"},
	{"XX", "National Provider Identifier", "National Provider Identifier — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, billing/ordering provider \"*****XX*1234567893\"); same code x12Names835.ts ID_QUALIFIER_NAME.XX names on the 835 side; frontend/src/data/x12Names837p.ts NM108_CONFIRMED_NAME"},
	{"MI", "Member Identification Number", "Member Identification Number — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3, subscriber \"****MI*##########\"); matches CEDI 837P CG, PDF p.33 (printed p.25; Table 15) \"NM108 Subscriber Identification Code Qualifier MI\"; frontend/src/data/x12Names837p.ts NM108_CONFIRMED_NAME"},
	{"PI", "Payer Identification", "Payer Identification — CEDI 837P CG, PDF p.35 (printed p.27; Table 17) \"NM108 Payer Identification Code Qualifier PI\"; frontend/src/data/x12Names837p.ts NM108_CONFIRMED_NAME"},
	{"II", "Standard Unique Health Identifier for each Individual in the United States", "Standard Unique Health Identifier for each Individual in the United States — Transmittal 579, PDF p.18 (2010BA) and p.68 (2330A), Values column \"II, MI\"; the name is the guide's: " + tr3},
	{"XV", "Centers for Medicare and Medicaid Services PlanID", "Centers for Medicare and Medicaid Services PlanID — Transmittal 579, PDF p.14 (2010AC), p.20 (2010BB), and p.70 (2330B), Values column \"PI, XV\"; the name is the guide's: " + tr3},
}

func (v IdentificationCodeQualifierNM108) String() string { return string(v) }

// Known reports whether the value is on the list.
func (v IdentificationCodeQualifierNM108) Known() bool {
	return identificationCodeQualifierNM108List.known(string(v))
}

// ShortDescription returns the value's public name, or "" when the sources
// give it none or the value is not on the list.
func (v IdentificationCodeQualifierNM108) ShortDescription() string {
	return identificationCodeQualifierNM108List.short(string(v))
}

// Description returns the value's public name with its citation, or ""
// when the value is not on the list.
func (v IdentificationCodeQualifierNM108) Description() string {
	return identificationCodeQualifierNM108List.description(string(v))
}

// ParseIdentificationCodeQualifierNM108 returns text as an
// IdentificationCodeQualifierNM108, with an error when the text is neither
// empty nor on the list. The value is kept either way.
func ParseIdentificationCodeQualifierNM108(text string) (IdentificationCodeQualifierNM108, error) {
	return IdentificationCodeQualifierNM108(text), identificationCodeQualifierNM108List.check("NM108", text)
}
