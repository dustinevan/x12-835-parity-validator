package x12837p

// ===========================================================================
// 2000B / 2010BA / 2010BB: Subscriber and Payer
// ===========================================================================
//
// NameNM1, AddressN3, LocationN4, ContactPER, ReferenceREF, and
// SubscriberSBR are defined elsewhere in this package (name.go,
// address.go, contact.go, reference.go, subscriber_sbr.go) and reused here
// unchanged.

// DemographicDMG holds one DMG segment. Used at 2010BA (Subscriber
// Demographic Information, this file) and 2010CA (patient.go); both loops
// send the identical three-position shape (DMG_EL_NAME). DMG01 (Date Time
// Period Format Qualifier) is D8 at both loops (Transmittal 579, PDF p.18
// and p.25), and the parser reads DMG02 as CCYYMMDD; DMG01 is still kept as
// received, so a file that sends a different qualifier, or none, keeps that
// fact. DMG04 onward are Not Used, and a file that sends them keeps their
// text in Extra. Set is false when the segment did not appear at all.
type DemographicDMG struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	FormatDMG01 string `json:"formatDMG01,omitempty"` // Date Time Period Format Qualifier, D8

	BirthDateDMG02 Date   `json:"birthDateDMG02,omitempty"` // Birth Date; Transmittal 907, PDF p.145-146 (footer 13-14)
	GenderDMG03    string `json:"genderDMG03,omitempty"`    // Gender Code; external code list, not covered by an enum here

	Extra []string `json:"extra,omitempty"` // elements after DMG03, kept as sent
}

// PatientInfoPAT holds one PAT segment: Patient Information. Used at 2000B
// (this file; PAT01-04 Not Used there, because Medicare's subscriber is
// always the patient) and 2000C (patient.go; the whole loop, PAT included,
// is not used by Medicare — NOT_USED_BY_MEDICARE). PAT05 (format
// qualifier, D8; Transmittal 579, PDF p.17) is kept as received, and the
// parser reads PAT06 as CCYYMMDD, as DemographicDMG does for DMG01/DMG02.
// Set is false when the segment did not appear at all.
type PatientInfoPAT struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	RelationshipPAT01     RelationshipPAT01             `json:"relationshipPAT01,omitempty"`     // Individual Relationship Code; Not Used at 2000B
	LocationPAT02         string                        `json:"locationPAT02,omitempty"`         // Patient Location Code; external code list
	EmploymentStatusPAT03 string                        `json:"employmentStatusPAT03,omitempty"` // Employment Status Code; external code list
	StudentStatusPAT04    string                        `json:"studentStatusPAT04,omitempty"`    // Student Status Code; external code list
	FormatPAT05           string                        `json:"formatPAT05,omitempty"`           // Date Time Period Format Qualifier, D8
	DeathDatePAT06        Date                          `json:"deathDatePAT06,omitempty"`        // Patient Death Date
	WeightUOMPAT07        UnitOrBasisForMeasurementCode `json:"weightUOMPAT07,omitempty"`        // Unit or Basis for Measurement Code
	WeightPAT08           Amount                        `json:"weightPAT08,omitempty"`           // Patient Weight
	PregnancyPAT09        YesNoConditionResponseCode    `json:"pregnancyPAT09,omitempty"`        // Pregnancy Indicator

	Extra []string `json:"extra,omitempty"` // elements after PAT09, kept as sent
} // Transmittal 907, PDF p.144 (footer 12)

// SubscriberNM1IL holds the 2010BA loop: the subscriber's NM1 (NM101
// confirmed as IL, LOOP_NM1_NAME/NM101_CONFIRMED_NAME), its N3/N4, its
// demographic information (DMG), its Property and Casualty Claim Number
// REF, its Subscriber Secondary Identification REF (not used by Medicare),
// and its contact PER.
type SubscriberNM1IL struct {
	NM1          NameNM1        `json:"nm1,omitempty"`          // Transmittal 907, PDF p.144-145 (footer 12-13)
	Address      AddressN3      `json:"address,omitempty"`      // Transmittal 907, PDF p.145 (footer 13)
	Location     LocationN4     `json:"location,omitempty"`     // Transmittal 907, PDF p.145 (footer 13)
	Demographics DemographicDMG `json:"demographics,omitempty"` // Transmittal 907, PDF p.145-146 (footer 13-14)

	// References holds the Property and Casualty Claim Number REF
	// (SUBSCRIBER_PC_REF_NAME, PDF p.146 footer 14) and the Subscriber
	// Secondary Identification REF (SUBSCRIBER_SECONDARY_REF_NAME, PDF p.146
	// footer 14, not used by Medicare); neither REF01 qualifier is confirmed
	// in either reviewed source.
	References []ReferenceREF `json:"references,omitempty"`

	Contact ContactPER `json:"contact,omitempty"` // Property and Casualty Subscriber Contact Information; Transmittal 907, PDF p.146-147 (footer 14-15)

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop the model does not place
}

// PayerNM1 holds the 2010BB loop: the payer's NM1 (role code PR, confirmed
// by frontend/src/data/x12Names837p.ts's NM101_CONFIRMED_NAME), its N3/N4, its Payer Secondary
// Identification REF, and the Billing Provider Secondary Identification
// REF nested in this loop (both not used by Medicare).
type PayerNM1 struct {
	NM1      NameNM1    `json:"nm1,omitempty"`      // Transmittal 907, PDF p.147 (footer 15)
	Address  AddressN3  `json:"address,omitempty"`  // Transmittal 907, PDF p.148 (footer 16)
	Location LocationN4 `json:"location,omitempty"` // Transmittal 907, PDF p.148 (footer 16)

	// References holds the Payer Secondary Identification REF
	// (PAYER_SECONDARY_REF_NAME, PDF p.148 footer 16) and the Billing
	// Provider Secondary Identification REF that appears inside this loop
	// (PAYER_LOOP_BILLING_PROVIDER_REF_NAME, PDF p.148-149 footer 16-17);
	// both are not used by Medicare and neither qualifier is confirmed.
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"`
}

// SubscriberHL holds the 2000B loop: the HL, the SBR, the PAT (mostly Not
// Used at this loop under Medicare, since the subscriber is always the
// patient), the 2010BA/2010BB name loops, the claims for this subscriber as
// the patient (2300), and the 2000C patient loops whose HL02 names this HL.
type SubscriberHL struct {
	HL HierarchyHL `json:"hl,omitempty"` // Transmittal 907, PDF p.143 (footer 11)

	Info        SubscriberSBR  `json:"info,omitempty"`        // SBR
	PatientInfo PatientInfoPAT `json:"patientInfo,omitempty"` // PAT

	SubscriberNM1IL SubscriberNM1IL `json:"subscriberNM1IL,omitempty"` // 2010BA
	PayerNM1PR      PayerNM1        `json:"payerNM1PR,omitempty"`      // 2010BB

	Claims   []ClaimCLM  `json:"claims,omitempty"`   // 2300 loops that follow this HL before the next HL
	Patients []PatientHL `json:"patients,omitempty"` // 2000C loops whose HL02 is this HL's HL01

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop, before 2010BA, that the model does not place
}
