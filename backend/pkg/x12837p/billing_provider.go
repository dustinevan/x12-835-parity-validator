package x12837p

// ===========================================================================
// 2000A / 2010AA-AC: Billing Provider
// ===========================================================================
//
// NameNM1, AddressN3, LocationN4, ContactPER, ReferenceREF, ProviderPRV,
// and HierarchyHL are defined elsewhere in this package (name.go,
// address.go, contact.go, reference.go, provider.go, hierarchy.go) and
// reused here unchanged.

// CurrencyCUR is the CUR segment: Foreign Currency Information. CUR_NAME
// (frontend/src/data/x12Names837p.ts) notes Medicare does not support it;
// CUR_EL_NAME gives only two positions (CUR03-21 Not Used), unlike
// backend/pkg/x12835's three-element CUR; a file that sends CUR03 onward
// keeps that text in Extra. Set is false when the segment did not appear
// at all. Neither CUR01 nor CUR02 has a confirmed code table in
// frontend/src/data/x12Names837p.ts, so both stay plain strings.
type CurrencyCUR struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file, 0-based; the ISA is 0

	EntityIDCUR01     string `json:"entityIDCUR01,omitempty"`     // CUR01, Entity Identifier Code; Transmittal 907, PDF p.137 (footer 5); not used by Medicare
	CurrencyCodeCUR02 string `json:"currencyCodeCUR02,omitempty"` // CUR02, Currency Code

	Extra []string `json:"extra,omitempty"` // elements after CUR02, kept as sent
}

// BillingProviderNM185 holds the 2010AA loop: the billing provider's NM1
// (NM101 confirmed as 85, LOOP_NM1_NAME/NM101_CONFIRMED_NAME), its N3/N4,
// its tax ID REF, its UPIN/License REF (not used by Medicare), and its
// contact PER.
type BillingProviderNM185 struct {
	NM1      NameNM1    `json:"nm1,omitempty"`      // Transmittal 907, PDF p.138 (footer 6)
	Address  AddressN3  `json:"address,omitempty"`  // Transmittal 907, PDF p.138 (footer 6)
	Location LocationN4 `json:"location,omitempty"` // Transmittal 907, PDF p.138-139 (footer 6-7)

	// TaxIDREFEI is the Billing Provider Tax Identification box
	// (BILLING_PROVIDER_REF_NAME/EL_NAME); REF01 = EI is confirmed (CEDI
	// 837P CG, PDF p.50, printed p.42, Figure 3, "REF*EI*123456789"). The
	// parser fills it from the first REF*EI with a non-empty REF02.
	TaxIDREFEI ReferenceREF `json:"taxIDREFEI,omitempty"` // Transmittal 907, PDF p.139 (footer 7)

	// References holds the UPIN/License REF (BILLING_PROVIDER_UPIN_REF_NAME,
	// not used by Medicare; Transmittal 907, PDF p.139, footer 7) and any
	// repeat of the tax ID REF.
	References []ReferenceREF `json:"references,omitempty"`

	Contact ContactPER `json:"contact,omitempty"` // Transmittal 907, PDF p.139-140 (footer 7-8)

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop the model does not place
}

// PayToProviderNM1 holds the 2010AB loop: the pay-to provider's NM1 (role
// code 87, confirmed by frontend/src/data/x12Names837p.ts's
// NM101_CONFIRMED_NAME). NM105-09 are Not Used at this loop
// (NM1_ROLE_EL_NAME '2010AB'); no REF or PER box is named for this loop in
// frontend/src/data/x12Names837p.ts.
type PayToProviderNM1 struct {
	NM1      NameNM1    `json:"nm1,omitempty"`      // Transmittal 907, PDF p.140 (footer 8)
	Address  AddressN3  `json:"address,omitempty"`  // Transmittal 907, PDF p.140-141 (footer 8-9)
	Location LocationN4 `json:"location,omitempty"` // Transmittal 907, PDF p.141 (footer 9)

	Raw []RawSegment `json:"raw,omitempty"`
}

// PayToPlanNM1 holds the 2010AC loop: the pay-to plan's NM1. The loop is
// not used by Medicare in its entirety (NOT_USED_BY_MEDICARE,
// frontend/src/data/x12Names837p.ts, citing CEDI 837P CG, PDF p.26-27,
// printed p.18-19, Table 5) but is still modeled, matching how
// backend/pkg/x12835 keeps its own Medicare-unsupported CUR segment typed.
// Its NM101 role code is PE, confirmed by NM101_CONFIRMED_NAME; both its
// REF boxes' qualifiers are not confirmed in either reviewed source.
type PayToPlanNM1 struct {
	NM1      NameNM1    `json:"nm1,omitempty"`      // Transmittal 907, PDF p.141-142 (footer 9-10); loop not used by Medicare
	Address  AddressN3  `json:"address,omitempty"`  // Transmittal 907, PDF p.142 (footer 10)
	Location LocationN4 `json:"location,omitempty"` // Transmittal 907, PDF p.142 (footer 10)

	// References holds the Pay-To Plan Secondary Identification REF
	// (PAY_TO_PLAN_SECONDARY_REF_NAME, PDF p.142 footer 10) and the Pay-To
	// Plan Tax Identification REF (PAY_TO_PLAN_TAX_REF_NAME, PDF p.143
	// footer 11); neither qualifier is confirmed, and the loop is not used
	// by Medicare at all.
	References []ReferenceREF `json:"references,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"`
}

// BillingProviderHL holds the 2000A loop: the HL, the billing provider's
// specialty (PRV) and foreign currency (CUR, not used by Medicare)
// segments, the 2010AA/AB/AC name loops it contains, and the 2000B
// subscriber loops whose HL02 names this HL.
type BillingProviderHL struct {
	HL HierarchyHL `json:"hl,omitempty"` // Transmittal 907, PDF p.136 (footer 4)

	Specialty ProviderPRV `json:"specialty,omitempty"` // PRV, "Billing Provider Specialty Information"; Transmittal 907, PDF p.137 (footer 5)
	Currency  CurrencyCUR `json:"currency,omitempty"`  // CUR, "Foreign Currency Information"; not used by Medicare

	ProviderNM185      BillingProviderNM185 `json:"providerNM185,omitempty"`      // 2010AA
	PayToProviderNM187 PayToProviderNM1     `json:"payToProviderNM187,omitempty"` // 2010AB, optional
	PayToPlanNM1PE     PayToPlanNM1         `json:"payToPlanNM1PE,omitempty"`     // 2010AC, optional; not used by Medicare

	Subscribers []SubscriberHL `json:"subscribers,omitempty"` // 2000B loops whose HL02 is this HL's HL01

	Raw []RawSegment `json:"raw,omitempty"` // segments in this loop, before 2010AA, that the model does not place
}
