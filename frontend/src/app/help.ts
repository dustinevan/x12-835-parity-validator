// Every piece of explanatory text in the validator app, in one place, so a
// sentence can be changed without touching layout. Three layers per page:
//
//   purpose  one sentence under the page title: what the page answers
//   terms    the definitions behind the small ⓘ beside numbers, measures,
//            chart titles, and column headers, keyed by element id
//   guide    the "How to read this page" drawer opened from the ? button
//
// Definitions describe what the validator computes (backend/lib/core/query/
// summary.go, lib/core/service/match.go and compare.go, lib/core/control/
// views/views.go). When the code changes, change the sentence here.
// Each term's `label` names the thing it defines; `text` is the definition.

export type PageId = 'overview' | 'differences' | 'claims' | 'workloads' | 'unmatched' | 'claim'

export interface Term {
  label: string
  text: string
}

export interface GuideSection {
  heading: string
  text: string
}

export interface PageHelp {
  title: string
  purpose: string
  terms: Record<string, Term>
  guide: GuideSection[]
}

// ---------------------------------------------------------------- shared definitions
// Terms used on several pages. A page lists the ones it shows.

const outcomeMatch: Term = {
  label: 'Outcome-state match',
  text: 'Matched pairs whose claim status code (CLP02) is the same on both sides ÷ MAC (Medicare Administrative Contractor) claims not set aside; a MAC claim missing from the vendor submission, or sent more than once by it (duplicated), counts as a mismatch, and extra vendor claims are left out. The target, ≥ 95%, is the solicitation’s published objective. Defined by this project; the program has not published its formula.',
}

const exactlyEqual: Term = {
  label: 'Exactly equal',
  text: 'Equal claims ÷ MAC (Medicare Administrative Contractor) claims not set aside. Equal means every compared element agrees: amounts as exact decimals (10.0 equals 10.00) and everything else exactly; missing and duplicated claims count against the rate and extra vendor claims are not in it. A project measure; the program has not set a target for it.',
}

const equal: Term = {
  label: 'Equal',
  text: 'The match key found exactly one MAC (Medicare Administrative Contractor) claim and one vendor claim, and every compared element agrees.',
}

const different: Term = {
  label: 'Different',
  text: 'The match key found exactly one MAC (Medicare Administrative Contractor) claim and one vendor claim, and at least one compared element differs.',
}

const missing: Term = {
  label: 'Missing from vendor',
  text: 'A MAC (Medicare Administrative Contractor) claim whose match key has no claim in the vendor submission.',
}

const extra: Term = {
  label: 'Extra on vendor',
  text: 'A vendor claim whose match key has no MAC (Medicare Administrative Contractor) claim, counted once per key, however many times the vendor sent it.',
}

const matchKey: Term = {
  label: 'Match key',
  text: 'A MAC (Medicare Administrative Contractor) claim and a vendor claim are paired when six fields agree: patient MBI (Medicare Beneficiary Identifier, NM1*QC), date of service, billed charge (CLP03), tax ID (REF*TJ), payee NPI (National Provider Identifier, N1*PE), and MAC workload (N1*PR). Each is put in one spelling first (MBI upper case, tax ID without hyphens, charge as an exact decimal); the patient control number (CLP01) and payer claim number (CLP07) are not used.',
}

const serviceDate: Term = {
  label: 'Service date',
  text: 'The statement from date (DTM*232) when the claim has one, else the earliest service-line date (DTM*472, or DTM*150 for a period). One of the six match-key fields.',
}

const claimType: Term = {
  label: 'Claim type',
  text: 'Read from the claim: institutional when it has statement dates (DTM*232 or DTM*233), an MIA segment, or a line revenue code (SVC04); otherwise DME (durable medical equipment) when the payer name (N1*PR) contains “DME”; otherwise professional. The MAC (Medicare Administrative Contractor) claim decides, except for extra claims.',
}

const workload: Term = {
  label: 'Workload',
  text: 'The MAC (Medicare Administrative Contractor) workload: the payer identifier in the 835 header (loop 1000A, N1*PR, element 04). It is a match-key field, so claims are never paired across workloads.',
}

const claimsAffected: Term = {
  label: 'Claims affected',
  text: 'Different claims (matched pairs with at least one difference) that carry this element as a difference, in the chosen filters.',
}

const amountDiffering: Term = {
  label: 'Amount differing',
  text: 'The sum of |vendor − MAC (Medicare Administrative Contractor)| for this element over those claims, every occurrence and service line included; overpayments and underpayments add rather than cancel. Code, date, and count elements have no amount.',
}

const notCompared: Term = {
  label: 'What is compared',
  text: 'Each claim runs from its CLP segment to the segment before the next claim, and every element inside it is compared: amounts as exact decimals (80.00 and 80 are the same value); everything else exactly, adjustment quantities included; and any segment the comparison has no field for compared whole. An amount sent on one side only is a difference however small. Left out, under the compared set recorded on the submission: the match-key fields where they repeat as values (CLP03, the patient ID in NM1*QC, DTM*232), the control numbers each system assigns (CLP07, REF*F8), the two CLP elements Medicare does not use (CLP10, CLP14), the claim contact (PER), the line control number that pairs the lines (REF*6R), and service-line QTY. Nothing outside the claim is compared: envelope and control numbers (ISA, GS, ST), payment and payee headers (BPR, TRN, N1 loops), payment-level adjustments (PLB), and provider summaries (LX, TS3, TS2).',
}

// The submission pickers in the navigation and the page toggle
// (VendorPicker.tsx, selection.tsx), the same on every page.
export const SELECTION = {
  picker: {
    label: 'Submissions shown',
    text: 'Any vendor submission for each of two slots, 1 and 2, chosen by vendor and date (oldest first). Slot 1 defaults to the first vendor’s latest submission and slot 2 to the second vendor’s; “No submission” leaves a slot empty. The two can be two vendors, two dates of one vendor, or the same submission twice. The overview shows both, and the claim page shows both beside the MAC (Medicare Administrative Contractor) claim; on both, an empty slot’s column stays empty. Differences, Claims, Workloads, and Unmatched show one of the two, chosen with the toggle under the page title, or the other slot when the chosen one is empty. With both empty, a page asks for a submission instead. The choice is kept in the address (v1= and v2=) and carried by every link, so a copied link opens the same view.',
  },
  side: {
    label: 'Vendor shown',
    text: 'Which of the two submissions picked in the navigation this page shows. The choice is kept in the address (side=) and carried to the other pages. A claim opened from this page shows both submissions beside the MAC (Medicare Administrative Contractor) claim.',
  },
} satisfies Record<string, Term>

// The MAC map beside the pages (MacTree.tsx), from backend/pkg/macs.
export const MAC_MAP: Term = {
  label: 'Workloads and roll-ups',
  text: 'A workload is the claims of one Part (A, B, or DME) for a state, part of a state, or group of states that one Medicare Administrative Contractor (MAC) processes, and CMS gives each a five-digit number. The validator reads a claim’s workload from the payer identifier in the 835 header (N1*PR, element 04) and pairs claims only within one workload. A roll-up is a number CMS gives the workloads of one Part in a jurisdiction that one processing region handles together; some also carry the home health and hospice (HHH) claims. CMS’s change requests differ on whether roll-up numbers appear in files sent outside CMS, so the map shows each roll-up with the workloads under it.',
}

const picker = 'The navigation picks the two submissions (any vendor submission for each of slots 1 and 2, defaulting to the vendors’ latest, or “No submission” to leave one empty), and the toggle under the title chooses which of the two this page shows; both choices stay in the address and in every link.'

// ---------------------------------------------------------------- pages

export const HELP = {
  differences: {
    title: 'Differences by kind',
    purpose: 'Which kinds of difference affect the most claims and dollars in this submission.',
    terms: {
      pairsCompared: {
        label: 'Claim pairs compared',
        text: 'Equal plus different results in the chosen claim type and workload: pairs of one MAC (Medicare Administrative Contractor) claim and one vendor claim with the same match key.',
      },
      exactlyEqual,
      withAnyDifference: {
        label: 'With any difference',
        text: 'Pairs with at least one differing element; the percentage is of pairs compared.',
      },
      kindsSeen: {
        label: 'Kinds seen',
        text: 'Distinct elements that differ in at least one pair, in the chosen filters.',
      },
      paidMoved: {
        label: 'Medicare paid moved',
        text: 'The sum of |vendor − MAC (Medicare Administrative Contractor)| claim payment (CLP04) over the different pairs; overpayments and underpayments add rather than cancel. Missing, duplicated, and set-aside claims are not included.',
      },
      concentration: {
        label: 'Claims with a difference, by kind',
        text: 'Each segment is one kind’s claim count as a share of all kind counts added together. A claim with three kinds counts three times here, so these shares are smaller than the table’s share of different claims.',
      },
      whatDiffers: {
        label: 'What differs',
        text: 'The X12 element that differs, such as CLP04 (claim payment) or CAS*CO*253 (an adjustment, named by group and reason code), with its plain label.',
      },
      claimsAffected,
      shareOfDifferent: {
        label: 'Share of different claims',
        text: 'Claims affected ÷ pairs with any difference, in the same filters. A claim with several kinds counts under each, so the column can add to more than 100%.',
      },
      amountDiffering,
      effects: {
        label: 'Kind groups',
        text: 'The colored square groups kinds by element, a grouping this page applies: Claim outcome is CLP02; Patient owes is CLP05 or any patient-responsibility (PR) adjustment; Amounts is any other element with a dollar difference; Codes only is an element without a dollar amount, such as a code, date, identifier, unit count, or weight.',
      },
    },
    guide: [
      { heading: 'Colors and marks', text: 'Amber bars count claims; blue bars add dollars. The small colored square beside each kind is its group, explained in the legend under the table; “no amount” marks an element without a dollar value.' },
      { heading: 'What clicking does', text: `${picker} Claim type and workload tabs filter, and the address keeps them. Rank by switches between claims and dollars. A row opens Claims with that kind in Must have; the long-tail row expands the rest.` },
      { heading: 'Example', text: 'Of 1,200 different claims, 480 carry CAS*CO*253: 480 ÷ 1,200 = 40%. If 300 of them differ by +$1.10 on that element and 180 by −$0.90, amount differing is $330.00 + $162.00 = $492.00; the signs do not cancel.' },
    ],
  },

  claims: {
    title: 'Claims by tag combination',
    purpose: 'Which claims in this submission carry a chosen combination of differences, and how far their Medicare payment moved.',
    terms: {
      filter: {
        label: 'Tag combination',
        text: 'Must have: the claim carries every tag listed; Any of: at least one; Must not have: none. The boxes combine with AND. With Must have and Any of both empty, the list is every result in the submission without the excluded tags, including equal, missing, duplicated, extra, and set-aside ones.',
      },
      byOutcome: {
        label: 'Result',
        text: 'Set by a cell of the Overview’s agreement grid: keeps the claims with this result in each vendor’s picked submission (equal, different, missing, or duplicated). The toggle under the title shows the same claims for the other vendor. Clear removes it.',
      },
      claims: {
        label: 'Claims',
        text: 'Results that pass the tag combination, the result filter when set, and the claim-type tab.',
      },
      paidDiffers: {
        label: 'Paid differs',
        text: 'Compared pairs in the selection whose claim payment (CLP04) differs: the amber bars of the histogram.',
      },
      samePaid: {
        label: 'Same paid',
        text: 'Compared pairs in the selection with the same claim payment (CLP04); they can still differ in other elements. Results without a pair (missing, duplicated, extra, set aside) are in neither count.',
      },
      paidTotal: {
        label: 'Paid difference, total',
        text: 'The sum of |vendor − MAC (Medicare Administrative Contractor)| claim payment (CLP04) over the selection; overpayments and underpayments add rather than cancel.',
      },
      histogram: {
        label: 'Difference in Medicare paid, per claim',
        text: 'Vendor paid − MAC (Medicare Administrative Contractor) paid (CLP04) for each compared pair in the selection. Below zero and above zero are each cut into up to five equal-width ranges; each label is the smallest and largest difference actually in the range, and empty ranges are left out.',
      },
      patientId: {
        label: 'Patient ID',
        text: 'The patient’s MBI (Medicare Beneficiary Identifier), from NM1*QC; a match-key field. Opens the claim.',
      },
      serviceDate,
      billed: {
        label: 'Billed',
        text: 'Total claim charge (CLP03). It is a match-key field, so it is the same on both sides.',
      },
      difference: {
        label: 'Difference',
        text: 'Vendor claim payment minus MAC (Medicare Administrative Contractor) claim payment (CLP04); positive means the vendor paid more.',
      },
      claimType,
      workload,
    },
    guide: [
      { heading: 'Colors and marks', text: 'Solid amber chips are Must have, outlined amber chips Any of, and struck-through chips Must not have; a chip’s number is its different-claim count in the submission. In the example, the MAC value is on top and the vendor value below, with vendor minus MAC in amber; a dashed “absent” box means the element is not on that side. Histogram bars are amber where paid differs and gray where it is the same.' },
      { heading: 'What clicking does', text: `${picker} Click a box to choose it, then click tags below to add them; click a placed tag to remove it. Column headers with arrows sort. A row opens the claim, the MAC claim beside both vendors. The filter lives in the address, so a copied link shows the same claims.` },
      { heading: 'Example', text: 'Must have CLP04 with Must not have CAS*CO*253 lists pairs whose claim payment differs for some reason other than sequestration. A pair paid $100.00 by the MAC and $98.00 by the vendor shows −$2.00 and falls in a bar below zero.' },
    ],
  },

  workloads: {
    title: 'By workload',
    purpose: 'Which MAC workloads fall short of the 95% outcome-state match target in this submission, and by how much.',
    terms: {
      workloadsMeeting: {
        label: 'Workloads at or above 95%',
        text: 'Workload rows whose outcome-state match is 95% or higher, out of the rows shown. A workload with claims of more than one claim type has a row for each.',
      },
      outcomeMatch,
      exactlyEqual: { ...exactlyEqual, label: 'No difference at all' },
      macClaims: {
        label: 'MAC claims',
        text: 'MAC (Medicare Administrative Contractor) claims in the workloads shown, set-aside claims included.',
      },
      tiles: {
        label: 'Workload tiles',
        text: 'One tile per workload: color is the outcome-state match band, and width grows with MAC (Medicare Administrative Contractor) claims. Part A, Part B, and DME name the claim type the validator read: institutional, professional, and durable medical equipment.',
      },
      workload,
      perThousand: {
        label: 'Different outcome, per 1,000 claims',
        text: '(1 − outcome-state match) × 1,000: per 1,000 MAC (Medicare Administrative Contractor) claims not set aside, those whose claim status (CLP02) differs, that are missing from the vendor, or that the vendor sent more than once. The dashed line at 50 is the 95% target.',
      },
      claims: {
        label: 'Claims',
        text: 'MAC (Medicare Administrative Contractor) claims in the workload row, set-aside claims included.',
      },
      largestKind: {
        label: 'Largest difference kind',
        text: 'The element carried by the most different claims in the workload, across claim types.',
      },
    },
    guide: [
      { heading: 'Colors and marks', text: 'Tile and bar colors are outcome-state match bands: dark green ≥ 98%, light green 95–98%, light amber 90–95%, amber below 90%. In the ranked list, rows below the target have an amber edge and tint, and the dashed line marks the target.' },
      { heading: 'What clicking does', text: `${picker} A tile or workload name opens Differences for that workload and claim type; a difference kind opens its claims. A claim-type tab’s dot is green when every workload of that type meets 95%.` },
      { heading: 'Example', text: 'A workload has 2,000 MAC claims not set aside, and 1,880 match on claim status: 1,880 ÷ 2,000 = 94.0%, which is (1 − 0.94) × 1,000 = 60 per 1,000, past the dashed line at 50.' },
    ],
  },

  unmatched: {
    title: 'Missing, extra, duplicated, and set aside',
    purpose: 'Which claims were not compared in this submission: missing from the vendor, extra on the vendor side, duplicated by the vendor, or set aside because the MAC repeats their match key.',
    terms: {
      missing: {
        label: 'Missing',
        text: 'MAC (Medicare Administrative Contractor) claims whose match key has no claim in the vendor submission.',
      },
      extra: {
        label: 'Extra',
        text: 'Vendor claims whose match key has no MAC (Medicare Administrative Contractor) claim. The tab counts keys; this list and its total count vendor claims, so they differ when a vendor repeats a key.',
      },
      duplicated: {
        label: 'Duplicated',
        text: 'The vendor sent this claim more than once; the MAC (Medicare Administrative Contractor) has it once; no copy can be chosen, so it counts as a miss for the vendor. The list shows the MAC claim, once per key; the claim page shows the first copy the vendor sent, for reference, and how many copies it sent.',
      },
      set_aside: {
        label: 'Set aside',
        text: 'MAC (Medicare Administrative Contractor) claims whose match key the MAC submission itself repeats; none are compared, whatever the vendor sent for the key. A key repeated only on the vendor side is duplicated, and a key repeated with no MAC claim is extra.',
      },
      share: {
        label: 'Share',
        text: 'Claims on this tab ÷ every MAC (Medicare Administrative Contractor) claim in the submission, set aside included; rate is the same per 1,000. The overview’s missing percentage divides by MAC claims not set aside, so the two can differ slightly.',
      },
      matchKey,
      sides: {
        label: 'MAC-processed and vendor-processed',
        text: 'Paid is the claim payment (CLP04); below it is the patient control number (CLP01), shown for lookup and not used in matching. The dashed “absent” box is the side with no claim for this key, and the “sent more than once” box a vendor that repeated the key; at most 1,000 claims are listed.',
      },
      collisions: {
        label: 'Repeated keys',
        text: 'Each box is one repeated key: its six key fields, then every MAC (Medicare Administrative Contractor) claim carrying it, with its patient control number (CLP01) and paid amount. ×N counts those MAC claims; at most 1,000 claims are listed.',
      },
    },
    guide: [
      { heading: 'Colors and marks', text: 'A dashed “absent” box marks the side with no claim; a dashed “sent more than once” box marks a vendor that repeated the key; a vendor claim with no MAC match has an amber edge. On a set-aside box, ×N counts the MAC claims sharing the key. A tab’s dot is green when its count is zero.' },
      { heading: 'What clicking does', text: `${picker} Tabs switch between missing, extra, duplicated, and set aside. A row or claim opens the claim page for that key, the MAC claim beside both vendors.` },
      { heading: 'Example', text: 'With 10,000 MAC claims, 25 missing claims are 25 ÷ 10,000 = 0.25%, or 2.5 per 1,000. A key the vendor sent twice and the MAC has once is duplicated: the MAC claim is listed once, no copy is compared, and it counts as a miss for the vendor. A key carried by two MAC claims is set aside: both appear in one ×2 box, and neither is compared.' },
    ],
  },

  overview: {
    title: 'Overview',
    purpose: 'Every count and sum the comparison produced for the two picked submissions, beside the legacy (MAC-processed) figures they are compared with, over the claims the selection at the top picks.',
    terms: {
      vendors: {
        label: 'Legacy, 1, and 2',
        text: 'Legacy is the MAC-processed (Medicare Administrative Contractor) submission, the reference both vendors are compared with. Slots 1 and 2 hold the two submissions picked in the navigation, by default each vendor’s latest; they can also be two dates of one vendor. Every table on the page has the same columns at the same places: the label, then Legacy, 1, 2. Legacy has figures only where the legacy side has its own; elsewhere its column is empty.',
      },
      selection: {
        label: 'Selection',
        text: 'Which claims the numbers are over. Each box is one way to cut the claims. With nothing picked, a box counts every value (All); a click on a value picks it, so the numbers count only the picked values, and a second click unpicks it. None leaves every value of a box out, and Reset selection clears every box. A box lists only the values that have claims under what the other boxes pick; a picked value with none stays, dimmed, so it can be unpicked. A value left out drops its claims from every figure except two: the agreement grid, which covers the whole submissions, and the value’s own breakdown table, which still lists it, struck through. The legacy figures follow only the boxes for what the MAC (Medicare Administrative Contractor) claim carries; the vendor claim status and outcome boxes do not change them. The selection is kept in the address, so a copied link opens the same view.',
      },
      workloads: {
        label: 'Workload',
        text: 'The MAC (Medicare Administrative Contractor) workloads with claims under what the other boxes pick, each under its jurisdiction, so picking a jurisdiction leaves only its workloads. A roll-up, the number CMS gives the workloads of one Part that one processing region handles together, heads the workloads under it, as on the MAC map. The box starts closed, showing how many workloads it lists or how many are picked.',
      },
      addressOnly: {
        label: 'Narrowed by the address',
        text: 'Part, MAC (Medicare Administrative Contractor) paid band, and payment difference band have no box, but a link can still narrow them, and then every figure is narrowed by them too. This line names them; Clear puts every value of them back.',
      },
      columns: {
        label: 'Legacy, 1, and 2 columns',
        text: 'Legacy: the MAC (Medicare Administrative Contractor) claims and what the MAC paid on them (CLP04), counted from the MAC submission itself, each claim once, so a key the MAC repeats counts each of its claims. It is cut only by what the MAC claim carries (jurisdiction, Part, workload, claim type, inpatient or outpatient, MAC status, frequency, paid band, service month), never by a vendor’s status, outcome, or payment difference, so it is the same for both slots and has figures only in the total, beside each MAC status, and in the breakdowns by those; elsewhere its column is empty. 1 and 2: the vendor’s results in the row and what it paid (CLP04); its results are each MAC claim’s key once, whatever the vendor did with it, plus the claims only it sent (extra). A dash is a figure the row cannot have: no vendor amount where the comparison kept no claim of the vendor’s (missing, duplicated, set aside), no legacy figure for the claims only a vendor sent.',
      },
      outcomes: {
        label: 'Every claim, by outcome',
        text: 'Each MAC (Medicare Administrative Contractor) claim counted once by its result for each vendor: equal, different, missing from the vendor, duplicated by the vendor, or set aside; a claim only the vendor sent is extra. An outcome is a vendor’s result, so Legacy has a figure only on the total line: the MAC claims. The bars share one scale: legacy’s is the MAC claims, and each vendor’s is its results by outcome with its extra claims last. Legacy counts claims and a vendor counts keys, so a key the MAC repeats is one set-aside result on a vendor’s bar and two or more claims on legacy’s.',
      },
      equal,
      different,
      missing,
      duplicated: {
        label: 'Duplicated by vendor',
        text: 'The vendor sent this claim more than once and the MAC (Medicare Administrative Contractor) has it once; no copy can be chosen, so the comparison keeps no claim of the vendor’s for it. Counted once per key, as its one MAC claim.',
      },
      extra,
      setAside: {
        label: 'Set aside',
        text: 'A match key the MAC (Medicare Administrative Contractor) submission itself repeats; none of its claims are compared, whatever the vendor sent for it. Decided by the MAC submission alone, so it is the same for both vendors. Counted once per key on this page.',
      },
      agreement: {
        label: 'How the two agree',
        text: 'Each MAC (Medicare Administrative Contractor) claim counted once, by its result in slot 1 (rows) and in slot 2 (columns): equal, different, missing, or duplicated. Claims set aside and extra claims are left out, so the cells add to the MAC claims not set aside, printed above the grid; a cell’s share is its count ÷ that total. The grid covers both whole submissions: the selection does not narrow it.',
      },
      sameElements: {
        label: 'Differ in the same elements',
        text: 'Claims different in both submissions whose differing elements are exactly the same list in both, for example CLP04 and CAS*CO*253 in each. The values may still differ between the two. When both vendors differ from the MAC claim in the same places, the MAC result, the input, or the validator is worth checking first.',
      },
      payment: {
        label: 'Payment agreement',
        text: 'Among pairs (one MAC claim and one vendor claim with the same key, equal or different): how many pay the same amount (CLP04) and how many do not, with the net and gross differences. Payment agreement is each vendor’s own comparison with the MAC (Medicare Administrative Contractor) claims, so the Legacy column is empty here. No rate is computed; the pair count is printed so any share can be read off it.',
      },
      samePayment: {
        label: 'Same payment',
        text: 'Pairs whose two payments (CLP04) match, the first payment band; these pairs have no CLP04 difference.',
      },
      pairsPaid: {
        label: 'Pairs',
        text: 'Each vendor’s pairs and what the vendor paid on them (CLP04).',
      },
      netGross: {
        label: 'Net and gross difference',
        text: 'Net: the sum of vendor paid minus MAC (Medicare Administrative Contractor) paid over the pairs, signed, so overpayments and underpayments cancel. Gross: the sum of the absolute differences, so they add.',
      },
      sameStatus: {
        label: 'Same status',
        text: 'Pairs whose claim status code (CLP02) is the same on both sides, printed with the pair count it is out of.',
      },
      bands: {
        label: 'Payment difference bands',
        text: 'Pairs by how far apart the two payments are: a match; up to a dollar; up to 1% of the MAC (Medicare Administrative Contractor) payment; up to 5%; more. Legacy has no difference of its own, so its column is empty and the chart is 1 and 2 side by side. A percentile is read off the bars: the band a given share of the pairs falls in is where that percentile lies. A band a link has left out of the selection is struck through, its bar gray.',
      },
      matrix: {
        label: 'Status agreement',
        text: 'Pairs counted by the MAC (Medicare Administrative Contractor) claim status code (CLP02), down the side, against each vendor’s own, across: slot 1’s matrix in its column, slot 2’s in its. The Legacy column beside each MAC status is the MAC claims with that status, all of them, not only those in pairs. The diagonal is pairs with the same status; a cell off it is a different outcome state. Which cells count as the same state is a reading of these matrices, not a rule of the tool.',
      },
      differs: {
        label: 'What differs',
        text: 'Among different pairs, the elements that differ, each as the comparison tags it (CLP04, CAS*CO*45, DTM*050), with the claims that carry the tag and the dollars on it (the sum of the absolute amount differences for that element). A claim with several differing elements counts under each, so the rows add to more than the different claims. Legacy has no count here, so its column is empty: a difference is between the legacy claim and a vendor’s.',
      },
      breakdowns: {
        label: 'Breakdowns',
        text: 'One table for each way the claims can be cut, every row one value with its Legacy, 1, and 2 figures. The tables by a vendor’s claim status, payment difference band, and outcome cut a vendor’s results, so their Legacy column is empty. Each table is added under the selection on every other way of cutting the claims and not its own, so a value left out of the selection is still listed, struck through, with what it would add back; the rows not struck through add up to the totals. A long table shows its first rows until opened.',
      },
    },
    guide: [
      { heading: 'Colors and marks', text: 'L, 1, and 2 mark the legacy (MAC-processed) column and the two picked submissions; every table has those columns at the same places, and the Legacy column is empty where the legacy side has no figure of its own. In the bars, the outlined bar is the legacy MAC claims; green is equal, amber different, the dashed box missing, the dotted box duplicated, light gray set aside, and dark gray extra. In the agreement grid, green is equal in both, light amber different in one, and dark amber different in both. In the status matrices, green is the diagonal (the same status) and amber a cell off it. In the selection, a picked value is dark, and a picked value with no claims under the rest of the selection is dimmed. A struck-through value is left out of the selection. Nothing on the page is a pass or a fail.' },
      { heading: 'What clicking does', text: 'A value’s tab picks it or unpicks it; All clears a box’s picks, None leaves every value of it out, and Reset selection clears every box. The workload box opens and closes on its name. A count, or a bar segment, opens Claims or Unmatched for that outcome or element, for that submission, carrying the claim type when exactly one is picked; an agreement cell opens Claims for exactly those claims. Show all opens a long breakdown table.' },
      { heading: 'Example', text: 'Picking Different among the outcomes narrows every vendor figure to the different pairs; the legacy figures stay the MAC claims, since an outcome is a vendor’s result, not something the MAC claim carries. The payment bands show how far apart those pairs’ payments are, and What differs lists the elements behind them. Picking one jurisdiction narrows everything, the legacy figures included, to that jurisdiction’s claims; By jurisdiction still lists the others, struck through, and the workload box lists only that jurisdiction’s workloads.' },
    ],
  },

  claim: {
    title: 'Compare claim',
    purpose: 'One MAC-processed claim beside each vendor’s claim with the same match key, one row per element in file order, marked where a value differs from the MAC’s.',
    terms: {
      matchKey,
      notCompared,
    },
    guide: [
      { heading: 'Colors and marks', text: 'The MAC column is the reference. Amber marks a vendor’s cell whose value differs from the MAC’s, with the signed difference for an amount; a faint dot is an element the claim does not carry. A row the compared set leaves out is dimmed and tagged “left out”; a watched element is tagged “watched”. The marks are computed on this page from the values; the counts elsewhere use the comparison’s own saved differences.' },
      { heading: 'What clicking does', text: 'The tree on the left hides or shows elements, one box for a group or an element in every claim and line. Only differences keeps the rows where a value differs. A preset is a named hidden set: pick one, save one, or copy its text to paste in another browser; presets are shared with the docs site’s 835 compare page.' },
      { heading: 'Example', text: 'The MAC claim pays $80.00 with a $20.00 CO-45 adjustment. One vendor pays $79.98 with $20.02 of CO-45: two cents apart, within five, so neither cell is marked. The other pays $79.90 with $20.10 of CO-45: its CLP04 cell is marked with −$0.10 and its CO-45 cell with +$0.10.' },
    ],
  },
} satisfies Record<PageId, PageHelp>

// The page each app route explains; routes without an entry have no guide.
export const PAGE_BY_PATH: Partial<Record<string, PageId>> = {
  '/': 'overview',
  '/claims': 'claims',
  '/differences': 'differences',
  '/unmatched': 'unmatched',
  '/workloads': 'workloads',
  '/claim': 'claim',
}
