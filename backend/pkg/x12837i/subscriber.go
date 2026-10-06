package x12837i

// ---------------------------------------------------------------------------
// Subscriber: 2000B (HL, SBR), 2010BA (subscriber NM1, N3, N4, DMG, REF),
// 2010BB (payer NM1, N3, N4, REF). SubscriberSBR is in subscriber_sbr.go
// and the SBR and DMG coded types in enums_subscriber.go.
// ---------------------------------------------------------------------------

// DemographicDMG holds one DMG segment, the same three-element shape at
// every instance (subscriber 2010BA; patient 2010CA). Set is false when
// the segment did not appear at all. Transmittal 907, PDF p.17 (subscriber,
// 2010BA); identical shape at patient 2010CA, PDF p.22.
type DemographicDMG struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	DateFormatDMG01 string   `json:"dateFormatDMG01,omitempty"` // Date Time Period Format Qualifier ("D8", Date Expressed in Format CCYYMMDD, is the only value this flat file documents); not enumerated, no DMG01 type is in this package's shared type list
	BirthDateDMG02  Date     `json:"birthDateDMG02,omitempty"`  // Birth Date
	SexDMG03        SexDMG03 `json:"sexDMG03,omitempty"`        // Gender Code

	// ExtraElements holds DMG04 through DMG11 (Marital Status Code, the
	// Composite Race or Ethnicity Information as one text, Citizenship
	// Status Code, Country Code, Basis of Verification Code, Quantity, Code
	// List Qualifier Code, Industry Code) as sent. X12 defines them; the
	// 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// SubscriberHL is the 2000B loop: the HL and SBR segments, the
// subscriber's own NM1/N3/N4/DMG/REF (2010BA), and the payer's own
// NM1/N3/N4/REF (2010BB).
type SubscriberHL struct {
	Hierarchy  HierarchyHL   `json:"hierarchy,omitempty"`  // HL; Transmittal 907, PDF p.16
	Subscriber SubscriberSBR `json:"subscriber,omitempty"` // SBR; PDF p.16

	// Subscriber (2010BA). Transmittal 907, PDF p.16-18.
	SubscriberNM1IL NameNM1        `json:"subscriberNM1IL,omitempty"` // NM1*IL, Insured/Subscriber; PDF p.16-17 (NM1_SUBSCRIBER_EL_NAME)
	Address         AddressN3      `json:"address,omitempty"`         // N3
	Location        LocationN4     `json:"location,omitempty"`        // N4
	Demographics    DemographicDMG `json:"demographics,omitempty"`    // DMG; PDF p.17

	// REF (2010BA): SY and Y4 each carry their own distinct name in
	// REF_SUBSCRIBER_NAME (unlike the payer loop below, where every
	// qualifier collapses to one combined name), so each gets its own
	// field here. PDF p.18.
	SecondaryIDREFSY ReferenceREF   `json:"secondaryIDREFSY,omitempty"` // REF*SY, Subscriber Secondary Identification (REF02 itself labeled Social Security Number)
	ClaimNumberREFY4 ReferenceREF   `json:"claimNumberREFY4,omitempty"` // REF*Y4, Property and Casualty Claim Number
	References       []ReferenceREF `json:"references,omitempty"`       // any other 2010BA REF, or a repeat

	// Payer (2010BB). Transmittal 907, PDF p.18-20.
	PayerNM1PR    NameNM1    `json:"payerNM1PR,omitempty"`    // NM1*PR, Payer; PDF p.18-19 (NM1_PAYER_EL_NAME)
	PayerAddress  AddressN3  `json:"payerAddress,omitempty"`  // N3
	PayerLocation LocationN4 `json:"payerLocation,omitempty"` // N4

	// PayerReferences (2010BB REF): 2U, EI, FY, and NF all carry the same
	// combined name "Payer Secondary Identification" in REF_PAYER_NAME, as
	// do G2 and LU (a second REF box at this loop) — frontend/src/data/
	// x12Names837i.ts's comment above REF_PAYER_NAME records that the box
	// captions and REF02 element labels for these two REFs appear swapped
	// in the source spreadsheet, so all six qualifiers are kept generically
	// here rather than split into individually named fields. PDF p.19-20.
	PayerReferences []ReferenceREF `json:"payerReferences,omitempty"`

	// Patients holds each 2000C loop whose HL02 names this loop's HL01.
	// Claims holds the 2300 loops that follow this loop directly, when the
	// subscriber is the patient (the Medicare case).
	Patients []PatientHL `json:"patients,omitempty"`
	Claims   []ClaimCLM  `json:"claims,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within 2000B/2010BA/2010BB this model does not place
}
