# Known differences between the engine and the taxonomy fixtures

`TestTaxonomyFixtures` (backend/lib/core/service/fixtures_test.go) runs each
fixture through the comparison engine under compared set 1 and compares the
result with `expected.json`. This file lists the cases where the two disagree
and the reason for each disagreement. A listed case passes only when its
disagreements are exactly the ones in its row. A new disagreement fails the
test, and so does a listed one that no longer occurs; remove it from the row
when the engine changes.

Two rules apply to every case and are not repeated in the rows:

- The engine compares claims only (CLP through the claim's last segment), so
  the test does not compare `file_differences` (payment, provider, and
  envelope scope). The test logs how many it skipped for each case.
- The engine writes a numeric value as an exact decimal without trailing
  fractional zeros (`80.00` reads `80`); a fixture holds the file's text. The
  test puts a fixture's numeric values in the engine's form before matching,
  so a spelling difference is not a disagreement.

The second column uses the test's wording:

- `outcome <fixture> vs <engine>`: the key's outcome differs.
- `key not found (<outcome>)`: no engine row has the fixture's key.
- `engine key not in fixture (<outcome>)`: an engine row's key is not in the fixture.
- `not compared <what>`: the fixture lists a difference in something the
  engine does not compare, or that compared set 1 leaves out
  (`backend/lib/core/service/compareset.go`: CLP03, CLP07, CLP10, CLP14,
  NM1*QC's NM109, DTM*232, REF*F8, the claim-level PER, line REF*6R, and
  line QTY). The engine does not produce a `line order` difference.
- `unreported <element>`, `unexpected <element>`, `value <element>`: element-level disagreements. The known ones are all `unreported` amounts within five cents of each other: the fixture records them as differences, as the taxonomy describes them, and the engine counts two amounts within five cents as the same amount (the five-cent rule, `AmountMatchCents` in `backend/lib/core/service/compare.go`).

Each row lists the case, its disagreements separated by `; `, and the reason.

| Case | Disagreements | Why |
|---|---|---|
| T1.06 | outcome extra+matched vs duplicated | The key occurs once in mac.835 and twice in vendor.835. The fixture pairs the first vendor claim and lists the second as extra; the engine reports a key the MAC has once and the vendor sent more than once as duplicated and compares nothing for it, since no copy can be chosen. |
| T1.11 | not compared REF*F8 target | The fixture records which of two same-key claims each reversal's REF*F8 points to. The engine does not read REF*F8: set 1 leaves it out as system-assigned. |
| T1.12 | not compared claim order | The fixture records claim positions within the file. The engine pairs claims by key and does not compare their order. |
| T1.15 | not compared REF*F8 | REF*F8 is present on the MAC side only. Set 1 leaves REF*F8 out, so its presence on one side is not reported. The other differences (payment, adjustments, MOA) agree. |
| T1.25 | not compared REF*6R | The MAC side reports no service lines and the vendor side three, so each vendor line is one-sided, and the fixture lists each line's REF*6R segment among its segments. Set 1 leaves REF*6R out as a value: it pairs the lines, and every line difference records it. The lines' other segments agree. |
| T1.26 | not compared REF*6R | The MAC's second line has no vendor partner, and the fixture lists its REF*6R segment as one-sided; set 1 leaves REF*6R out as a value. The line's other segments and the claim-level differences agree. |
| T1.27 | not compared REF*6R | The vendor splits a line into two pieces with the same REF*6R, so the second piece is an extra line, and the fixture lists its REF*6R segment as one-sided; set 1 leaves REF*6R out as a value. The rest agrees. |
| T1.28 | not compared REF*6R | The vendor bundles two MAC lines into one, so the MAC's second line is one-sided, and the fixture lists its REF*6R segment; set 1 leaves REF*6R out as a value. The rest agrees. |
| T1.29 | not compared line order | The three lines pair by REF*6R and differ only in their order, which the fixture lists as two position entries. The engine does not produce a `line order` difference: line order bears on no measure. |
| T1.31 | not compared REF*6R | The MAC lines carry REF*6R and the vendor lines none, so the lines pair by position and the only fixture entries are the two one-sided REF*6R segments, which set 1 leaves out as values. |
| T1.32 | not compared REF*6R | Two different claims share a key. The lines pair by position (no REF*6R value matches across the sides), and the fixture lists the REF*6R value difference on the paired line and the one-sided REF*6R of the MAC's third line; set 1 leaves REF*6R out as a value. Every other difference agrees. |
| T2.09 | unreported CLP04; unreported SVC03; unreported CAS*CO*253 | The MAC truncates sequestration to 2.13 and the vendor rounds it to 2.14, so CLP04, SVC03, and CAS*CO*253 each differ by one cent. The engine counts two amounts within five cents as the same amount, so it reports none of the three and the pair is equal. |
| T2.10 | unreported CAS*CO*253 | Sequestration (CAS*CO*253) is 1.97 against 2.02, five cents apart; the five-cent rule includes its edge, so the engine counts them as the same amount. The fee-schedule differences in the payment, coinsurance, CO-45, and allowed amount agree. |
| T2.14 | not compared REF*F8 | Each side's REF*F8 cites its own number for the original claim. The fixture lists the value difference; set 1 leaves REF*F8 out as system-assigned. |
| T2.18 | not compared CLP10 | The vendor sends a patient status code (CLP10) the MAC omits. Set 1 leaves CLP10 out: the CMS 835 flat file marks it not used by Medicare. |
| T2.37 | not compared DTM*232 | The MAC side sends statement dates and the vendor side a received date, so the fixture lists the paired DTM's qualifier and date as two differences. Set 1 leaves DTM*232 out, the key's date of service; the engine reports the vendor's DTM*050 as one-sided, which agrees with the fixture. |
| T2.40 | not compared REF*F8 | The entry is about the REF*F8 value difference. Set 1 leaves REF*F8 out as system-assigned, so the engine reports no difference for this claim. |
| T2.42 | not compared NM1 name or ID qualifier | The patient name differs in case and middle initial (NM103 to NM105). The engine compares only NM109 of an NM1, so it reports no difference for this claim. |
| T2.45 | not compared NM1 name or ID qualifier | The crossover carrier's name (NM1*TT NM103) differs as well as its ID. The engine reports the NM109 difference; it does not compare names. |
| T3.02 | unreported CAS*CO*253 | Sequestration (CAS*CO*253) is 1.92 against 1.90, two cents apart, which the engine counts as the same amount under the five-cent rule. The 1.50 pricing difference in the payment, CO-45, PR-2, and allowed amount agrees. |
| T3.16 | unreported CLP04; unreported CLP05; unreported SVC03; unreported CAS*PR*2 | Coinsurance rounded per line on the MAC side and per claim on the vendor side leaves CLP04, CLP05, and line 3's SVC03 and CAS*PR*2 one cent apart. The engine counts two amounts within five cents as the same amount, so it reports none of them and the pair is equal. |
| T3.18 | unreported CAS*CO*253 | Sequestration (CAS*CO*253) is 1.88 against 1.92, four cents apart, which the engine counts as the same amount under the five-cent rule. The missing CO-237 adjustment, its remark, and the payment it moves agree. |
| T3.43 | unreported CAS*CO*253 | Sequestration (CAS*CO*253) is 1.25 against 1.23, two cents apart, which the engine counts as the same amount under the five-cent rule. The CO-45, PR-2, payment, and allowed amount differences agree. |
| T4.01 | not compared REF*6R | The vendor omits the MAC's zero-paid third line, so it is one-sided, and the fixture lists its REF*6R segment; set 1 leaves REF*6R out as a value. The rest agrees. |
| T4.03 | not compared REF*6R | The vendor reports a rental span as two lines sharing REF*6R L1, so the second is an extra line, and the fixture lists its REF*6R segment; set 1 leaves REF*6R out as a value. The rest agrees. |
| T4.04 | not compared REF*6R | The entry is about the REF*6R value: the vendor carries the first 20 of 30 characters. The lines pair by position, and set 1 leaves REF*6R out as a value, so the engine reports no difference for this claim. |
| T4.13 | not compared REF*6R | The vendor splits a line into two with the same REF*6R, so the second is an extra line, and the fixture lists its REF*6R segment; set 1 leaves REF*6R out as a value. The SVC02 and other line differences agree. |
| T4.19 | unreported CLP04; unreported CLP05; unreported SVC03; unreported CAS*PR*2 | The MAC rounds 20 percent coinsurance half up to 14.67 and the vendor truncates it to 14.66, so CLP04, CLP05, SVC03, and CAS*PR*2 each differ by one cent. The engine counts two amounts within five cents as the same amount, so it reports none of them and the pair is equal. |
| T4.25 | not compared QTY (line) | The vendor writes the claim-level QTY*ZL at line level. Set 1 leaves service-line QTY out: the CMS 835 flat file marks it not used by Medicare. The missing claim-level QTY on the vendor side agrees. |
| T5.13 | key not found (extra); engine key not in fixture (extra) | The vendor file names the payee by tax ID (N1*PE with N103 FI) and has no REF*TJ. The fixture's key leaves NPI and tax ID empty. The engine uses N1*PE N104 as the NPI whatever N103 says. Both list the two vendor claims as extra, under different keys. |
