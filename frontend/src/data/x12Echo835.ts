// Elements that a public document states carry the value submitted on the
// 837 claim — i.e. the 835 is reporting back what the provider sent, not a
// payer-determined (adjudicated) value. This is the sourced mapping behind
// the 835 viewer's "carries the 837 value" tint. The full research —
// every element examined, its decision, and the quoted source text — is in
// docs/design/835-echo.md; this file holds only the elements that cleared
// the bar there. Nothing here is assumed from how values usually match:
// every row has a public-document quote behind it in that doc.
//
// Key format:
//   - A bare element reference ('SVC06', 'SVC07') applies wherever that
//     element appears.
//   - 'ELEMENT@QUALIFIER' qualifies the key by the segment's own qualifier
//     value, when the element's meaning (and sourcing) depends on it — e.g.
//     'REF02@6R' is REF02 only when REF01 is the "6R" qualifier (a line's
//     Line Item Control Number); 'REF02@EA' is REF02 only when REF01 is
//     "EA" (a claim's Other Claim Related Identification); 'NM109@QC' is
//     NM109 only inside the NM1*QC (Patient Name) loop.
//
// Each entry's `source` is the document, page, and a quoted phrase that
// supports the call. `from837` is filled in only when the cited source
// itself names the 837 element — none of the sources cited here do, so it
// is omitted throughout rather than guessed.

export interface Echo837Entry {
  /** The 837 element this value comes from, only when the source names it. */
  from837?: string
  /** Document + page (or URL + page/section) and a supporting quote. */
  source: string
}

export const X12_ECHO_837: Partial<Record<string, Echo837Entry>> = {
  // Claim level (loop 2100) — a claim-level "Other Claim Related
  // Identification" REF, qualifier EA, is used only when the provider put
  // its own proprietary ID (e.g. a medical record number) on the claim.
  'REF02@EA': {
    source:
      'CMS, Medicare Claims Processing Manual, Chapter 22 — Remittance Advice ' +
      '(https://www.cms.gov/regulations-and-guidance/guidance/manuals/downloads/clm104c22_remit_phase3_apr24-03_r2.pdf), ' +
      'p. 26, REF (040.A): "Situational, but required for Medicare if provider ' +
      'submitted a proprietary identification number on the claim. REF01 ' +
      'Required. Only “EA” applies to Medicare."',
  },

  // Service line level (loop 2110) — the Line Item Control Number the
  // provider assigned to the line on the claim, returned so the provider
  // can re-associate the paid line with the line it billed.
  'REF02@6R': {
    source:
      'CMS, Medicare Claims Processing Manual, Chapter 22 — Remittance Advice ' +
      '(https://www.cms.gov/regulations-and-guidance/guidance/manuals/downloads/clm104c22_remit_phase3_apr24-03_r2.pdf), ' +
      'p. 45, REF (100.A): "REF02 Required. Note: The provider line item ' +
      'control number (6R) is not used by and will not be retained by the ' +
      'Medicare core system. As with a 20-digit patient account number, use ' +
      'the COB data repository to populate REF02 for 6R."',
  },

  // Service line level — the composite "original" procedure/modifier
  // identification (SVC06), populated only when the paid procedure or drug
  // code differs from what was billed.
  SVC06: {
    source:
      'CMS, Medicare Claims Processing Manual, Chapter 22 — Remittance Advice, ' +
      'p. 28, SVC06-1: "Situational, but required if the procedure or drug ' +
      'code has been changed during adjudication." Corroborated by CMS 835 ' +
      'Flat File (5010A1, updated 1/10/2023) ' +
      '(https://www.cms.gov/medicare/billing/electronicbillingeditrans/downloads/835-flatfile.pdf), ' +
      'p. 12, the composite’s Procedure Code Description sub-element: ' +
      '"Medicare will populate if received on the claim."',
  },

  // Service line level — the originally billed unit count, populated only
  // when the paid units differ from the billed units.
  SVC07: {
    source:
      'CMS, Medicare Claims Processing Manual, Chapter 22 — Remittance Advice, ' +
      'p. 29, SVC07: "Situational, but required for Medicare if the paid ' +
      'units of service is different than the billed units of service."',
  },

  // Claim level (loop 2100, NM1*QC — Patient Name) — the patient
  // identifier (NM109, e.g. the HIC/MBI under qualifier MI), reported only
  // when it was reported on the incoming claim.
  'NM109@QC': {
    source:
      'CMS, Medicare Claims Processing Manual, Chapter 22 — Remittance Advice, ' +
      'p. 22 (institutional) and p. 38 (carrier), NM1 (030.A): "NM109 ' +
      'Situational, but required for Medicare if reported on the incoming claim."',
  },
}
