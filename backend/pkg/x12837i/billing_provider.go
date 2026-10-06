package x12837i

// ---------------------------------------------------------------------------
// Billing provider: 2000A (HL, PRV, CUR), 2010AA (NM1, N3, N4, REF, PER),
// 2010AB (pay-to address, NM1, N3, N4). The pay-to plan loop (2010AC,
// NM1*PE with its N3, N4, and REF; NM1_PAYTOPLAN_EL_NAME and
// REF_PAYTOPLAN_NAME in frontend/src/data/x12Names837i.ts) has no typed
// field yet; the parser keeps its segments in BillingProviderHL.Raw. The
// shared segment types (NameNM1, AddressN3, LocationN4, ContactPER,
// ReferenceREF, HierarchyHL, ProviderPRV) are defined in their own files.
// ---------------------------------------------------------------------------

// CurrencyCUR is the CUR segment: Foreign Currency Information. CUR03-21
// are Not Used in this flat file, unlike the 835's CUR, which also carries
// an Exchange Rate in CUR03. Set is false when the segment did not appear
// at all. Transmittal 907, PDF p.11 (2000A Foreign Currency Information).
type CurrencyCUR struct {
	Set   bool `json:"set,omitempty"`
	Index int  `json:"index,omitempty"` // the segment's position in the file

	EntityIDCUR01     string `json:"entityIDCUR01,omitempty"`     // Entity Identifier Code (85, Billing Provider, is the only value this flat file documents)
	CurrencyCodeCUR02 string `json:"currencyCodeCUR02,omitempty"` // Currency Code

	// ExtraElements holds CUR03 through CUR21 as sent. X12 defines them;
	// the 837I guide marks them Not Used.
	ExtraElements []string `json:"extraElements,omitempty"`
}

// BillingProviderHL is the 2000A loop: the HL and PRV segments (and,
// rarely, CUR), the billing provider's own NM1/N3/N4/REF/PER (2010AA), and
// the pay-to address's NM1/N3/N4 (2010AB).
type BillingProviderHL struct {
	Hierarchy HierarchyHL `json:"hierarchy,omitempty"` // HL; Transmittal 907, PDF p.10
	Provider  ProviderPRV `json:"provider,omitempty"`  // PRV; PDF p.10
	Currency  CurrencyCUR `json:"currency,omitempty"`  // CUR, rarely used; PDF p.11

	// Billing provider (2010AA). Transmittal 907, PDF p.12-13.
	BillingProviderNM185 NameNM1        `json:"billingProviderNM185,omitempty"` // NM1*85, Billing Provider; PDF p.12 (NM1_BILLING_EL_NAME)
	Address              AddressN3      `json:"address,omitempty"`              // N3
	Location             LocationN4     `json:"location,omitempty"`             // N4
	TaxIDREFEI           ReferenceREF   `json:"taxIDREFEI,omitempty"`           // REF*EI, Billing Provider Tax Identification; PDF p.13
	References           []ReferenceREF `json:"references,omitempty"`           // any other 2010AA REF, or a repeat
	Contacts             []ContactPER   `json:"contacts,omitempty"`             // PER; PDF p.13

	// Pay-to address (2010AB). NM1_PAYTO_EL_NAME: every name and ID position
	// is Not Used in this flat file; only the entity-type qualifier and
	// address flow, so PayToAddressNM187's LastNM103 through IDNM109 are
	// expected empty in practice even though the shared NameNM1 shape still
	// models them. Transmittal 907, PDF p.13-14.
	PayToAddressNM187 NameNM1    `json:"payToAddressNM187,omitempty"` // NM1*87, Pay-To Address
	PayToAddress      AddressN3  `json:"payToAddress,omitempty"`      // N3
	PayToLocation     LocationN4 `json:"payToLocation,omitempty"`     // N4

	// Subscribers holds each 2000B loop whose HL02 names this loop's HL01.
	Subscribers []SubscriberHL `json:"subscribers,omitempty"`

	Raw []RawSegment `json:"raw,omitempty"` // segments within 2000A/2010AA/2010AB this model does not place, including every segment of the 2010AC pay-to plan loop
}
