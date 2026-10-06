# Taxonomy fixtures

Each taxonomy entry (`docs/taxonomy/`) has one directory, `testdata/taxonomy/<ID>/`,
with these files. All data is synthetic.

- `source.837`: the submitted claim or claims.
- `mac.835`: the MAC-processed remittance, the reference side.
- `vendor.835`: the vendor-processed remittance.
- `expected.json`: what a comparison of the two 835s should find.

T1.03, T1.04, and T5.38 hold only `expected.json`, marked
`"representable": false`: neither side has an 835 that shows the entry.
T5.46 has no `mac.835`, because the missing MAC remittance is the subject of
that entry.

`KNOWN-DIFFERENCES.md` lists the cases where the comparison engine and a
fixture disagree.

## The test

`backend/lib/core/service/fixtures_test.go` (`TestTaxonomyFixtures`) splits
both 835s with the engine's `Split`, runs `CompareSubmissions`, and compares
the engine's outcome for each key and its claim- and line-scope differences
with `expected.json`:

```
cd backend && go test ./lib/core/service/ -run TestTaxonomyFixtures -v
```

The engine compares claims only, so the test does not compare
`file_differences`. The engine labels a value by element (CLP04, SVC03) or by
segment and qualifier (AMT*B6, NM1*TT, CAS*CO*253). It pairs the n-th value
with a label on one side with the n-th on the other. The fixtures name
positions instead: a segment, its occurrence, and an element number. The test
translates each fixture entry into the engine differences it implies. It
reads the files to find qualifiers and the values the engine pairs. The
mapping is written out above `impliedFacts` in the test. The test also checks
that every fixture value is the value in its file.

## expected.json

Fields appear in this order. Optional fields are omitted when they do not
apply.

| Field | Meaning |
|---|---|
| `taxonomy_id`, `title` | The taxonomy entry. |
| `representable` | False when no pair of 835 files can show the entry. `reason` then says why. |
| `reason` | Only when `representable` is false. |
| `notes` | What the fixture shows, and which cause it follows. |
| `imbalance` | Optional. A list of deliberate imbalances: `side` (mac or vendor), `level` (line, claim, or payment), and `detail`. Every file balances unless it is listed here. |
| `system_assigned_differences_ignored` | Elements each system assigns for itself. They are never listed as differences. `REF02 (qualifier F8)` appears only where REF*F8 cites each system's own CLP07. |
| `claims` | One entry for each key the files contain. |
| `file_differences` | Differences outside claims, in payment, provider, and envelope scope. The list is always present and can be empty. |

### Claims

| Field | Meaning |
|---|---|
| `key` | The six match-key fields: `mbi`, `date_of_service`, `total_charge`, `tax_id`, `npi`, `workload`. A field is null when the file has no value for it. |
| `patient_control_number` | CLP01. It comes from the MAC claim, or from the vendor claim for an extra claim. |
| `outcome` | `matched`, `missing` (MAC only), `extra` (vendor only), or `set_aside`. |
| `occurrences` | Optional, for set-aside keys: how many claims share the key on each side. |
| `differences` | Claim- and line-scope differences of a matched pair. Empty for other outcomes. |

The key fields are the MBI (NM1*QC NM109); the date of service (DTM*232,
else the earliest line DTM*472 or DTM*150); the total charge (CLP03); the tax
ID (REF*TJ); the payee NPI (N1*PE N104 when N103 is XX); and the workload
(N1*PR N104). Claims pair by key. A key that repeats on the MAC side is set
aside. A key that repeats only on the vendor side pairs its first vendor
claim and lists the rest as extra (T1.06).

### Differences

| Field | Meaning |
|---|---|
| `segment` | The segment tag. |
| `element` | The element position: CLP04 is 4. `0` means the segment is present on one side only. |
| `occurrence` | The segment's position among segments with the same tag in the same scope, counting from 1. Always present. |
| `scope` | `claim`, `line`, `payment`, `provider`, or `envelope` (below). |
| `line`, `vendor_line` | Line scope only: the service line's position in the MAC claim and in the vendor claim. Each is present when that side has the line. |
| `transaction`, `vendor_transaction` | `file_differences` in files with more than one transaction only: the ST-SE transaction's position on each side. |
| `kind` | Optional: `order` or `resolved_target` (below). Absent for an ordinary difference. |
| `system_assigned` | Optional, in `file_differences`: true for a system-assigned element listed because it is the subject of the entry (T5.06, T5.40, T5.41). |
| `mac`, `vendor` | The values as the files write them. A value is null when it is absent, and never the empty string. For `element: 0`, the side that has the segment holds the whole segment text without the terminator, including any system-assigned elements. The other side is null. |

Scopes:

- `claim`: CLP up to the first SVC.
- `line`: an SVC and the segments after it, up to the next SVC or the end of the claim.
- `payment`: BPR, TRN, CUR, the header REF and DTM segments, and the 1000A payer loop (N1*PR, N3, N4, REF, PER).
- `provider`: the 1000B payee loop (N1*PE, N3, N4, REF), the header-number loop (LX, TS3, TS2), and PLB.
- `envelope`: ISA, GS, ST, SE, GE, and IEA.

### Pairing

- Claims pair by key.
- Service lines pair by REF*6R when the value is present and unique on both
  sides. Lines left over pair by position, in order. A line with no partner
  is present on one side only, and every one of its segments is listed as
  `element: 0`.
- Within a scope, segments pair by tag and by position within the tag: the
  second CAS pairs with the second CAS. Elements pair by position.
- Numeric elements compare as exact decimals, so `10.0` equals `10.00`.
  Identifiers are text.
- Line order: when paired lines appear in a different relative order, each
  paired line whose position differs gets one entry. The entry has `segment`
  SVC, `element` 0, `kind` `order`, and the values `"position i"` and
  `"position j"`. A line that moves only because an earlier line was added or
  removed has no order entry.
- Claim order (T1.12): `kind: order` in claim scope, with the claim's
  position in each file.
- `kind: resolved_target` (T1.11): REF*F8 is compared by the claim it points
  to, because its raw value is system-assigned.
