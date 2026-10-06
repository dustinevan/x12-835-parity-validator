// The 837P's loop-rule table: every box this viewer draws, which segment
// (and, where the same segment ID is reused for several loops, which
// qualifier value in which element) opens it, which box types are valid
// directly inside which other box type, and which of the four tint
// families (or neutral) it belongs to — see colors.ts for the actual
// classes. frontend/src/docs/viewer837p/parse.ts reads this table to decide
// what the next segment opens; frontend/src/docs/viewer837p/describe.ts and
// colors.ts read it for display. Grounded in the 24 worked examples under
// frontend/src/data/viewer837p/files/ and docs/design/837p-examples.md.
//
// Loop nesting (⊃ = contains): ISA ⊃ GS ⊃ ST ⊃ (BHT, submitter, receiver,
// billing provider HL ⊃ subscriber HL ⊃ [patient HL] ⊃ claim CLM ⊃
// (claim-level loops, other subscribers, service line LX ⊃ (drug, line
// providers, line adjudication, forms))). The 2000C Patient Hierarchical
// Level / 2010CA Patient Name loop is never built by Medicare (CEDI 837P
// CG, PDF p.26-27, printed p.18-19, Table 5: "Must not be present — the
// subscriber is always the same as the patient") and appears in no example
// this viewer has, so this table defines no rule for it; an HL with
// HL03='23', if one ever appeared, would render as a plain, unboxed row
// rather than a guessed box shape.
export type BoxType =
  | 'ISA'
  | 'GS'
  | 'ST'
  | 'SUBMITTER'
  | 'RECEIVER'
  | 'BILLING_HL'
  | 'BILLING_PROVIDER'
  | 'PAYTO'
  | 'SUBSCRIBER_HL'
  | 'SUBSCRIBER'
  | 'PAYER'
  | 'CLM'
  | 'REFERRING'
  | 'RENDERING'
  | 'SERVICE_FACILITY'
  | 'SUPERVISING'
  | 'AMB_PICKUP'
  | 'AMB_DROPOFF'
  | 'OTHER_SUBSCRIBER'
  | 'OTHER_SUBSCRIBER_NAME'
  | 'OTHER_PAYER'
  | 'LX'
  | 'DRUG'
  | 'RENDERING_LINE'
  | 'PURCHASED_SERVICE'
  | 'SERVICE_FACILITY_LINE'
  | 'SUPERVISING_LINE'
  | 'ORDERING'
  | 'REFERRING_LINE'
  | 'AMB_PICKUP_LINE'
  | 'AMB_DROPOFF_LINE'
  | 'LINE_ADJUDICATION'
  | 'FORMS'

// Which of the four marked tint families a box belongs to, or 'neutral' for
// the envelope and the two HL wrapper levels — see colors.ts. 1 = submitter/
// receiver/billing provider, 2 = subscriber/patient, 3 = claim (and its
// 2310A-F provider loops and 2320/2330A-B other-subscriber loops), 4 =
// service line (and everything nested inside it: 2410, 2420A-H, 2430, 2440).
export type Tint = 'neutral' | 1 | 2 | 3 | 4

export interface LoopRule {
  boxType: BoxType
  segmentId: string
  // Narrows the trigger beyond the segment id alone, for the segment IDs
  // (NM1, HL) this file reuses across several loops — 0-based index into
  // `elements` and the code(s) that identify this loop. A rule with no
  // qualifier matches on segment id alone (LIN, SVD, LQ, CLM, SBR, LX each
  // open exactly one kind of box).
  qualifierIndex?: number
  qualifierCodes?: string[]
  tint: Tint
}

export const LOOP_RULES: LoopRule[] = [
  { boxType: 'SUBMITTER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['41'], tint: 1 }, // 1000A — CEDI 837P CG, PDF p.50 (printed p.42; Figure 3)
  { boxType: 'RECEIVER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['40'], tint: 1 }, // 1000B
  { boxType: 'BILLING_PROVIDER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['85'], tint: 1 }, // 2010AA
  { boxType: 'PAYTO', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['87'], tint: 1 }, // 2010AB — frontend/src/data/viewer837p/files/pay-to-address.837
  { boxType: 'SUBSCRIBER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['IL'], tint: 2 }, // 2010BA
  { boxType: 'PAYER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['PR'], tint: 2 }, // 2010BB
  { boxType: 'REFERRING', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['DN'], tint: 3 }, // 2310A — frontend/src/data/viewer837p/files/specialist-referral.837
  { boxType: 'RENDERING', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['82'], tint: 3 }, // 2310B
  { boxType: 'SERVICE_FACILITY', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['77'], tint: 3 }, // 2310C
  { boxType: 'SUPERVISING', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['DQ'], tint: 3 }, // 2310D — frontend/src/data/viewer837p/files/incident-to-service.837
  { boxType: 'AMB_PICKUP', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['PW'], tint: 3 }, // 2310E — frontend/src/data/viewer837p/files/ambulance.837
  { boxType: 'AMB_DROPOFF', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['45'], tint: 3 }, // 2310F
  { boxType: 'OTHER_SUBSCRIBER', segmentId: 'SBR', tint: 3 }, // 2320 — only valid as a direct child of CLM; the 2000B SBR is consumed as a plain row before this rule table is consulted (see parse.ts)
  { boxType: 'OTHER_SUBSCRIBER_NAME', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['IL'], tint: 3 }, // 2330A — same code as SUBSCRIBER, valid only inside OTHER_SUBSCRIBER
  { boxType: 'OTHER_PAYER', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['PR'], tint: 3 }, // 2330B — same code as PAYER, valid only inside OTHER_SUBSCRIBER
  { boxType: 'LX', segmentId: 'LX', tint: 4 }, // 2400, repeats
  { boxType: 'DRUG', segmentId: 'LIN', tint: 4 }, // 2410
  { boxType: 'RENDERING_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['82'], tint: 4 }, // 2420A — same code as RENDERING, valid only inside LX
  { boxType: 'PURCHASED_SERVICE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['QB'], tint: 4 }, // 2420B
  { boxType: 'SERVICE_FACILITY_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['77'], tint: 4 }, // 2420C — same code as SERVICE_FACILITY, valid only inside LX
  { boxType: 'SUPERVISING_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['DQ'], tint: 4 }, // 2420D — same code as SUPERVISING, valid only inside LX
  { boxType: 'ORDERING', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['DK'], tint: 4 }, // 2420E
  { boxType: 'REFERRING_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['DN'], tint: 4 }, // 2420F — same code as REFERRING, valid only inside LX
  { boxType: 'AMB_PICKUP_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['PW'], tint: 4 }, // 2420G — same code as AMB_PICKUP, valid only inside LX
  { boxType: 'AMB_DROPOFF_LINE', segmentId: 'NM1', qualifierIndex: 0, qualifierCodes: ['45'], tint: 4 }, // 2420H — same code as AMB_DROPOFF, valid only inside LX
  { boxType: 'LINE_ADJUDICATION', segmentId: 'SVD', tint: 4 }, // 2430
  { boxType: 'FORMS', segmentId: 'LQ', tint: 4 }, // 2440
]

// Which box types may appear as a direct child of which — the grammar a
// segment's context (not its code alone) disambiguates. A box type absent
// here (every leaf: SUBMITTER, RECEIVER, BILLING_PROVIDER, PAYTO,
// SUBSCRIBER, PAYER, the 2310/2420 provider loops, OTHER_SUBSCRIBER_NAME,
// OTHER_PAYER, DRUG, LINE_ADJUDICATION, FORMS) opens no further nested
// boxes of its own; everything inside it is a plain row.
export const CHILDREN: Partial<Record<BoxType, BoxType[]>> = {
  BILLING_HL: ['BILLING_PROVIDER', 'PAYTO', 'SUBSCRIBER_HL'],
  SUBSCRIBER_HL: ['SUBSCRIBER', 'PAYER', 'CLM'],
  CLM: ['REFERRING', 'RENDERING', 'SERVICE_FACILITY', 'SUPERVISING', 'AMB_PICKUP', 'AMB_DROPOFF', 'OTHER_SUBSCRIBER', 'LX'],
  OTHER_SUBSCRIBER: ['OTHER_SUBSCRIBER_NAME', 'OTHER_PAYER'],
  LX: [
    'DRUG',
    'RENDERING_LINE',
    'PURCHASED_SERVICE',
    'SERVICE_FACILITY_LINE',
    'SUPERVISING_LINE',
    'ORDERING',
    'REFERRING_LINE',
    'AMB_PICKUP_LINE',
    'AMB_DROPOFF_LINE',
    'LINE_ADJUDICATION',
    'FORMS',
  ],
}

export function matches(rule: LoopRule, segmentId: string, elements: string[]): boolean {
  if (rule.segmentId !== segmentId) return false
  if (rule.qualifierIndex === undefined) return true
  return (rule.qualifierCodes ?? []).includes(elements[rule.qualifierIndex] ?? '')
}

// Six box types open with no qualifier disambiguation at all (the parser,
// frontend/src/docs/viewer837p/parse.ts, opens them positionally — ISA/GS/ST
// each hold exactly one child box of the next type down, and CLM/BILLING_HL/
// SUBSCRIBER_HL are each the one HL or CLM segment expected at that point in
// the file) and so carry no LOOP_RULES row of their own. Their tint is
// given directly here; every other box type's tint comes from its row
// above. tintOf, in colors.ts, is the one place both are combined.
export const UNGOVERNED_TINT: Partial<Record<BoxType, Tint>> = {
  ISA: 'neutral',
  GS: 'neutral',
  ST: 'neutral',
  BILLING_HL: 'neutral',
  SUBSCRIBER_HL: 'neutral',
  CLM: 3,
}
