package x12837i

import "fmt"

// PartyRoleNM101 is NM101, the Entity Identifier Code, shared by every NM1
// segment in the transaction regardless of loop (1000A/B, 2010AA/AB/AC,
// 2010BA/BB/CA, 2310A-F, 2330A-I, 2420A-D) —
// frontend/src/data/x12Names837i.ts's NM1_ROLE_NAME table, which states
// the same code means the same role everywhere it appears.
type PartyRoleNM101 string

const (
	PartyRoleNM10141 PartyRoleNM101 = "41"
	PartyRoleNM10140 PartyRoleNM101 = "40"
	PartyRoleNM10185 PartyRoleNM101 = "85"
	PartyRoleNM10187 PartyRoleNM101 = "87"
	PartyRoleNM101PE PartyRoleNM101 = "PE"
	PartyRoleNM101IL PartyRoleNM101 = "IL"
	PartyRoleNM101PR PartyRoleNM101 = "PR"
	PartyRoleNM101QC PartyRoleNM101 = "QC"
	PartyRoleNM10171 PartyRoleNM101 = "71"
	PartyRoleNM10172 PartyRoleNM101 = "72"
	PartyRoleNM101ZZ PartyRoleNM101 = "ZZ"
	PartyRoleNM10182 PartyRoleNM101 = "82"
	PartyRoleNM10177 PartyRoleNM101 = "77"
	PartyRoleNM101DN PartyRoleNM101 = "DN"
)

// The codes the guide allows for NM101, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var partyRoleNM101Codes = []codeEntry{
	{"41", "Submitter", "Submitter — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"40", "Receiver", "Receiver — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"85", "Billing Provider", "Billing Provider — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"87", "Pay-To Address", "Pay-To Address — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"PE", "Pay-To Plan", "Pay-To Plan — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"IL", "Insured/Subscriber", "Insured/Subscriber — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"PR", "Payer", "Payer — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"QC", "Patient", "Patient — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"71", "Attending Provider", "Attending Provider — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"72", "Operating Physician", "Operating Physician — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"ZZ", "Other Operating Physician", "Other Operating Physician — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"82", "Rendering Provider", "Rendering Provider — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"77", "Service Facility Location", "Service Facility Location — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
	{"DN", "Referring Provider", "Referring Provider — Transmittal 907, PDF p.9 (41, 1000A); p.10 (40, 1000B); p.12 (85, 2010AA; also 2330I); p.13 (87, 2010AB); p.14 (PE, 2010AC); p.16 (IL, 2010BA; also 2330A); p.18 (PR, 2010BB; also 2330B); p.20 (QC, 2010CA); p.92 (71, 2310A; also 2330C); p.93 (72, 2310B; also 2330D, 2420A); p.94 (ZZ, 2310C; also 2330E, 2420B); p.94 (82, 2310D; also 2330G, 2420C); p.95 (77, 2310E; also 2330F); p.96 (DN, 2310F; also 2330H, 2420D)"},
}

// String returns the code as the file sent it.
func (v PartyRoleNM101) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v PartyRoleNM101) Known() bool { return knownCode(partyRoleNM101Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v PartyRoleNM101) ShortDescription() string { return shortOf(partyRoleNM101Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v PartyRoleNM101) Description() string { return descriptionOf(partyRoleNM101Codes, string(v)) }

// ParsePartyRoleNM101 reads one NM101 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParsePartyRoleNM101(code string) (PartyRoleNM101, error) {
	v := PartyRoleNM101(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing NM101: a value not on the list")
}

// EntityTypeNM102 is NM102, the Entity Type Qualifier on an NM1 segment (1
// = person, 2 = non-person entity, by the segment's own structure), shared
// across every NM1 instance in the transaction.
type EntityTypeNM102 string

const (
	EntityTypeNM1021 EntityTypeNM102 = "1"
	EntityTypeNM1022 EntityTypeNM102 = "2"
)

// The codes the guide allows for NM102, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var entityTypeNM102Codes = []codeEntry{
	{"1", "Person", "Person — Transmittal 907, PDF p.9"},
	{"2", "Non-Person Entity", "Non-Person Entity — Transmittal 907, PDF p.9"},
}

// String returns the code as the file sent it.
func (v EntityTypeNM102) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v EntityTypeNM102) Known() bool { return knownCode(entityTypeNM102Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v EntityTypeNM102) ShortDescription() string { return shortOf(entityTypeNM102Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v EntityTypeNM102) Description() string { return descriptionOf(entityTypeNM102Codes, string(v)) }

// ParseEntityTypeNM102 reads one NM102 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseEntityTypeNM102(code string) (EntityTypeNM102, error) {
	v := EntityTypeNM102(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing NM102: a value not on the list")
}

// IDQualifierNM108 is NM108, the Identification Code Qualifier, shared
// across every NM1 instance that carries an ID (submitter/receiver,
// billing provider, pay-to plan, payer, subscriber).
type IDQualifierNM108 string

const (
	IDQualifierNM108XX IDQualifierNM108 = "XX"
	IDQualifierNM10846 IDQualifierNM108 = "46"
	IDQualifierNM108PI IDQualifierNM108 = "PI"
	IDQualifierNM108XV IDQualifierNM108 = "XV"
	IDQualifierNM108II IDQualifierNM108 = "II"
	IDQualifierNM108MI IDQualifierNM108 = "MI"
)

// The codes the guide allows for NM108, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var idQualifierNM108Codes = []codeEntry{
	{"XX", "National Provider Identifier", "National Provider Identifier — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
	{"46", "Electronic Transmitter Identification Number", "Electronic Transmitter Identification Number — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
	{"PI", "Payer Identification", "Payer Identification — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
	{"XV", "Health Care Financing Administration National PlanID", "Health Care Financing Administration National PlanID — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
	{"II", "Standard Unique Health Identifier for each Individual in the United States", "Standard Unique Health Identifier for each Individual in the United States — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
	{"MI", "Member Identification Number", "Member Identification Number — Transmittal 907, PDF p.12 (XX, NM109 described as an NPI wherever XX is used); PDF p.9 (46, submitter/receiver); PDF p.9 (PI, pay-to plan/payer, alternate with XV); PDF p.17 (II, subscriber, alternate with MI); PDF p.17 (MI, subscriber — carries the Medicare Beneficiary Identifier for Medicare claims)"},
}

// String returns the code as the file sent it.
func (v IDQualifierNM108) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v IDQualifierNM108) Known() bool { return knownCode(idQualifierNM108Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v IDQualifierNM108) ShortDescription() string { return shortOf(idQualifierNM108Codes, string(v)) }

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v IDQualifierNM108) Description() string {
	return descriptionOf(idQualifierNM108Codes, string(v))
}

// ParseIDQualifierNM108 reads one NM108 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseIDQualifierNM108(code string) (IDQualifierNM108, error) {
	v := IDQualifierNM108(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing NM108: a value not on the list")
}

// DateQualifierDTP01 is DTP01, the Date/Time Qualifier, shared by every DTP
// segment in the transaction (claim-level 2300, other-payer 2330B,
// service-line 2400, line-level COB 2430) — Transmittal 907, PDF p.24
// states the element name is consistent across every DTP instance.
//
// 573 has two sourced names: at 2330B (other
// payer) Transmittal 907 names it "Date Claim Check or Remittance Date"
// (PDF p.105), while at 2430 (line-level COB) it names the identical code
// "Line Check or Remittance Date" (PDF p.124) — the same code, two
// sourced names depending on loop, the same situation as x12835's
// RefQualifierREF01TJ. 573 is code-named rather than picking one sourced
// meaning over the other.
type DateQualifierDTP01 string

const (
	DateQualifierDTP01096 DateQualifierDTP01 = "096"
	DateQualifierDTP01434 DateQualifierDTP01 = "434"
	DateQualifierDTP01435 DateQualifierDTP01 = "435"
	DateQualifierDTP01050 DateQualifierDTP01 = "050"
	DateQualifierDTP01573 DateQualifierDTP01 = "573"
	DateQualifierDTP01472 DateQualifierDTP01 = "472"
)

// The codes the guide allows for DTP01, each with its public name ("" when
// the reviewed sources give none) and that name with its citation.
var dateQualifierDTP01Codes = []codeEntry{
	{"096", "Discharge Hour", "Discharge Hour — Transmittal 907, PDF p.24 (claim-level DTP)"},
	{"434", "Statement From-Through Dates", "Statement From-Through Dates — Transmittal 907, PDF p.24 (claim-level DTP)"},
	{"435", "Admission Date and Hour", "Admission Date and Hour — Transmittal 907, PDF p.24 (claim-level DTP)"},
	{"050", "Repricer Received Date", "Repricer Received Date — Transmittal 907, PDF p.24 (claim-level DTP)"},
	{"573", "", "ambiguous, see the type doc comment; Transmittal 907, PDF p.105 (2330B); PDF p.124 (2430)"},
	{"472", "Service Date", "Service Date — Transmittal 907, PDF p.113 (2400 service line)"},
}

// String returns the code as the file sent it.
func (v DateQualifierDTP01) String() string { return string(v) }

// Known reports whether the code is on the list.
func (v DateQualifierDTP01) Known() bool { return knownCode(dateQualifierDTP01Codes, string(v)) }

// ShortDescription returns the code's public name, or "" when the
// reviewed sources give it none or the code is not on the list.
func (v DateQualifierDTP01) ShortDescription() string {
	return shortOf(dateQualifierDTP01Codes, string(v))
}

// Description returns the code's public name with its citation, or ""
// when the code is not on the list.
func (v DateQualifierDTP01) Description() string {
	return descriptionOf(dateQualifierDTP01Codes, string(v))
}

// ParseDateQualifierDTP01 reads one DTP01 value.
// The text comes back as sent, with an error when it is not on the list,
// so a caller keeps what the file sent; "" is the element not sent.
func ParseDateQualifierDTP01(code string) (DateQualifierDTP01, error) {
	v := DateQualifierDTP01(code)
	if code == "" || v.Known() {
		return v, nil
	}
	return v, fmt.Errorf("x12837i: parsing DTP01: a value not on the list")
}
