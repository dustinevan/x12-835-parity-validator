// Package macs maps a Medicare payer identifier to the MAC jurisdiction, the
// claim Part, and the contractor behind it.
//
// The identifiers are the five-digit numbers CMS assigns to Medicare
// Administrative Contractor work: a workload number (CMS also calls it a
// contractor number) for each state, or part of a state, and each Part, and
// a roll-up number for a group of workloads that one claims-processing
// region handles. The rows come from CMS's Pub 100-20 change requests (CR
// 5651 of 2007, "New Numbers for All MAC Jurisdictions", sets the numbering;
// each award's CR lists its numbers) and CMS's "Who are the MACs"
// jurisdiction pages, fetched 2026-10-02. This table is hand-maintained
// from those documents, not generated. Tree (tree.go) places every number
// under its jurisdiction and Part, roll-ups with the workloads under them.
//
// The public documents do not settle which 835 element carries one of these
// numbers (the CEDI 835 companion guide prints the DME numbers in REF*2U,
// ISA06, and GS02). Lookup reports false for a number the table does not
// hold, and that is an ordinary result.
package macs

// Workload is one identifier CMS has published and what it stands for.
type Workload struct {
	ID           string   // five digits, for example "05101"
	Jurisdiction string   // "J5" through "JN" for A/B MACs, "JA" through "JD" for DME MACs
	Part         string   // "A", "B", "DME", or "" when unknown; home health and hospice numbers are "A"
	Contractor   string   // the contractor holding the jurisdiction's contract
	States       []string // USPS codes the number covers; nil when CMS names none
}

// Contractor names as CMS's jurisdiction pages print them (fetched 2026-10-02).
const (
	wps        = "Wisconsin Physicians Service Government Health Administrators"
	ngs        = "National Government Services, Inc."
	cgs        = "CGS Administrators, LLC"
	noridian   = "Noridian Healthcare Solutions, LLC"
	novitas    = "Novitas Solutions, Inc."
	palmetto   = "Palmetto GBA, LLC"
	firstCoast = "First Coast Service Options, Inc."
)

// Workloads holds every identifier in the table, grouped by jurisdiction.
// Read it; do not modify it. Columns: ID, Jurisdiction, Part, Contractor,
// States. Each group's comment names the CMS change request (CR) it came
// from, with its transmittal number and date; a CR is a Pub 100-20 One-Time
// Notification.
//
// A number ending in 4 belongs to home health and hospice (HHH): an HHH
// workload (06004, 11004, 14014, 15004) or, for J6, the internal processing
// number 06014. The J6, J15, JK, and JM MACs process HHH claims for areas that
// differ from their A/B states, and an HHH row's States are its HHH area. A
// roll-up number lists only the A/B states its CR names, even where the CR
// says it also carries HHH claims.
var Workloads = []Workload{
	// J5. CR 11555 (Transmittal 10055, 2020-04-17) lists the workload
	// numbers and says the 2020 award does not change them. 05901 and the
	// Part A roll-up 05001: CR 8059 (Transmittal 1119, 2012-09-14). CR 8059
	// names no states for 05901, the WPS legacy (formerly Mutual of Omaha)
	// Part A workload.
	{"05101", "J5", "A", wps, []string{"IA"}},
	{"05201", "J5", "A", wps, []string{"KS"}},
	{"05301", "J5", "A", wps, []string{"MO"}},
	{"05401", "J5", "A", wps, []string{"NE"}},
	{"05901", "J5", "A", wps, nil},
	{"05001", "J5", "A", wps, []string{"IA", "KS", "MO", "NE"}},
	{"05102", "J5", "B", wps, []string{"IA"}},
	{"05202", "J5", "B", wps, []string{"KS"}},
	{"05302", "J5", "B", wps, []string{"MO"}},
	{"05402", "J5", "B", wps, []string{"NE"}},

	// J6. CR 11976 (Transmittal 10452, 2020-11-06). 06004 is the HHH
	// workload for both HHH Region D areas the CR lists. 06001 is the roll-up
	// for the Part A region that holds Wisconsin Part A and the first HHH
	// area; 06014 is what the CR calls an internal processing number, for
	// the region that holds the second HHH area (AK, AS, AZ, CA, GU, HI, ID,
	// MP, NV, OR, WA).
	{"06101", "J6", "A", ngs, []string{"IL"}},
	{"06201", "J6", "A", ngs, []string{"MN"}},
	{"06301", "J6", "A", ngs, []string{"WI"}},
	{"06004", "J6", "A", ngs, []string{"AK", "AS", "AZ", "CA", "GU", "HI", "ID", "MI", "MN", "MP", "NJ", "NV", "NY", "OR", "PR", "VI", "WA", "WI"}},
	{"06001", "J6", "A", ngs, []string{"WI"}},
	{"06014", "J6", "A", ngs, []string{"AK", "AS", "AZ", "CA", "GU", "HI", "ID", "MP", "NV", "OR", "WA"}},
	{"06102", "J6", "B", ngs, []string{"IL"}},
	{"06202", "J6", "B", ngs, []string{"MN"}},
	{"06302", "J6", "B", ngs, []string{"WI"}},

	// J8. CR 11153 (Transmittal 2263, 2019-02-22): the 2019 award does not
	// change the workload numbers.
	{"08101", "J8", "A", wps, []string{"IN"}},
	{"08201", "J8", "A", wps, []string{"MI"}},
	{"08102", "J8", "B", wps, []string{"IN"}},
	{"08202", "J8", "B", wps, []string{"MI"}},

	// J15. CR 9456 (Transmittal 1600, 2016-01-29) lists the workload numbers,
	// the HHH workload (Region B), and the roll-ups, which it calls "roll-up
	// number/payer ID"; CR 13562 (Transmittal 12574, 2024-04-11) says the
	// 2024 award changes none of them. The Part A roll-up 15001 also carries
	// the HHH claims.
	{"15101", "J15", "A", cgs, []string{"KY"}},
	{"15201", "J15", "A", cgs, []string{"OH"}},
	{"15004", "J15", "A", cgs, []string{"CO", "DC", "DE", "IA", "KS", "MD", "MO", "MT", "ND", "NE", "PA", "SD", "UT", "VA", "WV", "WY"}},
	{"15001", "J15", "A", cgs, []string{"KY", "OH"}},
	{"15102", "J15", "B", cgs, []string{"KY"}},
	{"15202", "J15", "B", cgs, []string{"OH"}},
	{"15002", "J15", "B", cgs, []string{"KY", "OH"}},

	// JE. CR 8226 (Transmittal 1201, 2013-03-22) assigned these numbers when
	// the workloads moved to the JE A/B MAC, and named the new Part A
	// roll-up 01011; CR 12306 (Transmittal 10856, 2021-07-13) says the 2021
	// award changes none of them. 01911 is the Part A "J1 WPS Legacy"
	// workload; no later CR repeats it, and CR 8226 names no states for it.
	// CR 8226 splits California Part B: 01112 is Northern California and
	// 01182 is Southern California (the CMS locality key page lists the
	// counties of each).
	{"01111", "JE", "A", noridian, []string{"CA"}},
	{"01211", "JE", "A", noridian, []string{"AS", "GU", "HI", "MP"}},
	{"01311", "JE", "A", noridian, []string{"NV"}},
	{"01911", "JE", "A", noridian, nil},
	{"01011", "JE", "A", noridian, []string{"AS", "CA", "GU", "HI", "MP", "NV"}},
	{"01112", "JE", "B", noridian, []string{"CA"}},
	{"01182", "JE", "B", noridian, []string{"CA"}},
	{"01212", "JE", "B", noridian, []string{"AS", "GU", "HI", "MP"}},
	{"01312", "JE", "B", noridian, []string{"NV"}},

	// JF. CR 14277 (Transmittal 13477, 2025-12-05) lists the workload
	// numbers and the two Part A roll-ups (02001 for AK, ID, OR, WA; 03001
	// for AZ, MT, ND, SD, UT, WY) and says the 2025 award changes none of
	// them. The CR names no Part B roll-up.
	{"02101", "JF", "A", noridian, []string{"AK"}},
	{"03101", "JF", "A", noridian, []string{"AZ"}},
	{"02201", "JF", "A", noridian, []string{"ID"}},
	{"03201", "JF", "A", noridian, []string{"MT"}},
	{"03301", "JF", "A", noridian, []string{"ND"}},
	{"02301", "JF", "A", noridian, []string{"OR"}},
	{"03401", "JF", "A", noridian, []string{"SD"}},
	{"03501", "JF", "A", noridian, []string{"UT"}},
	{"02401", "JF", "A", noridian, []string{"WA"}},
	{"03601", "JF", "A", noridian, []string{"WY"}},
	{"02001", "JF", "A", noridian, []string{"AK", "ID", "OR", "WA"}},
	{"03001", "JF", "A", noridian, []string{"AZ", "MT", "ND", "SD", "UT", "WY"}},
	{"02102", "JF", "B", noridian, []string{"AK"}},
	{"03102", "JF", "B", noridian, []string{"AZ"}},
	{"02202", "JF", "B", noridian, []string{"ID"}},
	{"03202", "JF", "B", noridian, []string{"MT"}},
	{"03302", "JF", "B", noridian, []string{"ND"}},
	{"02302", "JF", "B", noridian, []string{"OR"}},
	{"03402", "JF", "B", noridian, []string{"SD"}},
	{"03502", "JF", "B", noridian, []string{"UT"}},
	{"02402", "JF", "B", noridian, []string{"WA"}},
	{"03602", "JF", "B", noridian, []string{"WY"}},

	// JH. CR 14584 (Transmittal 13928, 2026-09-15) lists the workload
	// numbers and the Part A roll-ups (04011 for CO, NM, TX; 07001 for AR,
	// LA, MS) and says the award that takes effect 2026-12-01 changes none of
	// them. 04911 is the Part A workload for the legacy Mutual of Omaha
	// providers in CO, NM, OK, TX. The CR names no Part B roll-up and none
	// for Oklahoma Part A.
	{"04111", "JH", "A", novitas, []string{"CO"}},
	{"04211", "JH", "A", novitas, []string{"NM"}},
	{"04311", "JH", "A", novitas, []string{"OK"}},
	{"04411", "JH", "A", novitas, []string{"TX"}},
	{"04911", "JH", "A", novitas, []string{"CO", "NM", "OK", "TX"}},
	{"07101", "JH", "A", novitas, []string{"AR"}},
	{"07201", "JH", "A", novitas, []string{"LA"}},
	{"07301", "JH", "A", novitas, []string{"MS"}},
	{"04011", "JH", "A", novitas, []string{"CO", "NM", "TX"}},
	{"07001", "JH", "A", novitas, []string{"AR", "LA", "MS"}},
	{"04112", "JH", "B", novitas, []string{"CO"}},
	{"04212", "JH", "B", novitas, []string{"NM"}},
	{"04312", "JH", "B", novitas, []string{"OK"}},
	{"04412", "JH", "B", novitas, []string{"TX"}},
	{"07102", "JH", "B", novitas, []string{"AR"}},
	{"07202", "JH", "B", novitas, []string{"LA"}},
	{"07302", "JH", "B", novitas, []string{"MS"}},

	// JJ. CR 10316 (Transmittal 1960, 2017-11-03) renumbered the workloads
	// for the incoming JJ A/B MAC and named the Part A roll-up 10011 (FISS)
	// and the Part B roll-up 10012 (MCS); CR 13841 (Transmittal 12932,
	// 2024-10-30) says the 2024 award changes none of them.
	{"10111", "JJ", "A", palmetto, []string{"AL"}},
	{"10211", "JJ", "A", palmetto, []string{"GA"}},
	{"10311", "JJ", "A", palmetto, []string{"TN"}},
	{"10011", "JJ", "A", palmetto, []string{"AL", "GA", "TN"}},
	{"10112", "JJ", "B", palmetto, []string{"AL"}},
	{"10212", "JJ", "B", palmetto, []string{"GA"}},
	{"10312", "JJ", "B", palmetto, []string{"TN"}},
	{"10012", "JJ", "B", palmetto, []string{"AL", "GA", "TN"}},

	// JK. CT and NY: CR 5843 (Transmittal 343, 2008-05-23); the three NY
	// Part B workloads are named in CR 6694 (Transmittal 576, 2009-10-16):
	// 13202 Downstate NY, 13282 Upstate NY, 13292 Queens County. ME, MA, NH,
	// RI, VT, the HHH workload for Region A, and the Part A roll-ups 14011
	// (ME and the HHH area) and 14013 (NH, VT): CR 8303 (Transmittal 1314,
	// 2013-11-13). CR 12695 (Transmittal 11390, 2022-04-28) says the 2022
	// award changes none of them. No CR read names a roll-up for MA or RI.
	{"13101", "JK", "A", ngs, []string{"CT"}},
	{"13201", "JK", "A", ngs, []string{"NY"}},
	{"14111", "JK", "A", ngs, []string{"ME"}},
	{"14211", "JK", "A", ngs, []string{"MA"}},
	{"14311", "JK", "A", ngs, []string{"NH"}},
	{"14411", "JK", "A", ngs, []string{"RI"}},
	{"14511", "JK", "A", ngs, []string{"VT"}},
	{"14014", "JK", "A", ngs, []string{"CT", "MA", "ME", "NH", "RI", "VT"}},
	{"14011", "JK", "A", ngs, []string{"ME"}},
	{"14013", "JK", "A", ngs, []string{"NH", "VT"}},
	{"13102", "JK", "B", ngs, []string{"CT"}},
	{"13202", "JK", "B", ngs, []string{"NY"}},
	{"13282", "JK", "B", ngs, []string{"NY"}},
	{"13292", "JK", "B", ngs, []string{"NY"}},
	{"14112", "JK", "B", ngs, []string{"ME"}},
	{"14212", "JK", "B", ngs, []string{"MA"}},
	{"14312", "JK", "B", ngs, []string{"NH"}},
	{"14412", "JK", "B", ngs, []string{"RI"}},
	{"14512", "JK", "B", ngs, []string{"VT"}},

	// JL. CR 8327 (Transmittal 1245, 2013-06-07) lists the numbers when the
	// jurisdiction was renamed from J12 and says they stay the same; CR
	// 12448 (Transmittal 11033, 2021-10-01) says the 2021 award changes
	// none of them. The CR gives no state or Part for each number; the
	// states and Parts here come from CR 5651's table (2007, Part A ends in
	// 1 and Part B in 2, states in alphabetical order) and CR 5843 (12101 is
	// Delaware Part A). 12901 is the Part A workload for the legacy WPS
	// (formerly Mutual of Omaha) providers in DE, DC, MD, NJ, PA: CR 7135
	// (Transmittal 782, 2010-10-08).
	{"12101", "JL", "A", novitas, []string{"DE"}},
	{"12201", "JL", "A", novitas, []string{"DC"}},
	{"12301", "JL", "A", novitas, []string{"MD"}},
	{"12401", "JL", "A", novitas, []string{"NJ"}},
	{"12501", "JL", "A", novitas, []string{"PA"}},
	{"12901", "JL", "A", novitas, []string{"DC", "DE", "MD", "NJ", "PA"}},
	{"12102", "JL", "B", novitas, []string{"DE"}},
	{"12202", "JL", "B", novitas, []string{"DC"}},
	{"12302", "JL", "B", novitas, []string{"MD"}},
	{"12402", "JL", "B", novitas, []string{"NJ"}},
	{"12502", "JL", "B", novitas, []string{"PA"}},

	// JM. CR 9171 (Transmittal 1510, 2015-06-12) lists the workload numbers,
	// the HHH workload (Region C), and the Part A roll-ups, which it calls
	// "roll-up number/payer ID": 11001 for SC and the HHH claims, 11003 for
	// VA and WV. CR 12971 (Transmittal 11730, 2022-12-08) says the 2022
	// award changes none of them. The CR names no roll-up for NC Part A and
	// none for Part B.
	{"11201", "JM", "A", palmetto, []string{"SC"}},
	{"11301", "JM", "A", palmetto, []string{"VA"}},
	{"11401", "JM", "A", palmetto, []string{"WV"}},
	{"11501", "JM", "A", palmetto, []string{"NC"}},
	{"11004", "JM", "A", palmetto, []string{"AL", "AR", "FL", "GA", "IL", "IN", "KY", "LA", "MS", "NC", "NM", "OH", "OK", "SC", "TN", "TX"}},
	{"11001", "JM", "A", palmetto, []string{"SC"}},
	{"11003", "JM", "A", palmetto, []string{"VA", "WV"}},
	{"11202", "JM", "B", palmetto, []string{"SC"}},
	{"11302", "JM", "B", palmetto, []string{"VA"}},
	{"11402", "JM", "B", palmetto, []string{"WV"}},
	{"11502", "JM", "B", palmetto, []string{"NC"}},

	// JN. CR 6285 (Transmittal 423, 2008-12-24) assigned these numbers: 09201
	// is one Part A workload for Puerto Rico and the Virgin Islands, and Part
	// B has a separate workload for each of Florida, Puerto Rico, and the
	// Virgin Islands. CR 8759 (Transmittal 1390, 2014-06-06) lists the same
	// five numbers when the jurisdiction was renamed from J9; CR 12811
	// (Transmittal 11539, 2022-08-05) says the 2022 award changes none of
	// them. The CMS locality key page (2024 localities) lists the Virgin
	// Islands under 09202, not 09302.
	{"09101", "JN", "A", firstCoast, []string{"FL"}},
	{"09201", "JN", "A", firstCoast, []string{"PR", "VI"}},
	{"09102", "JN", "B", firstCoast, []string{"FL"}},
	{"09202", "JN", "B", firstCoast, []string{"PR"}},
	{"09302", "JN", "B", firstCoast, []string{"VI"}},

	// DME MACs. The CEDI 835 companion guide v5.6 (January 2025, Table 7)
	// prints all four as "DME MAC contractor ID" in REF*2U, ISA06, and GS02.
	// JA 16013: CR 9546 (Transmittal 1642, 2016-04-19) changed it from 16003
	// for the incoming JA DME MAC; CR 14040 (Transmittal 13240, 2025-05-22)
	// lists it. JB 17013: CR 9526 (Transmittal 1636, 2016-03-11) changed it
	// from 17003. JC 18003: CR 8235 (Transmittal 1210, 2013-04-19) and CR
	// 11985 (Transmittal 10462, 2020-11-13). JD 19003: CR 9453 (Transmittal
	// 1592, 2016-01-15). The states are from CMS's jurisdiction pages.
	{"16013", "JA", "DME", noridian, []string{"CT", "DC", "DE", "MA", "MD", "ME", "NH", "NJ", "NY", "PA", "RI", "VT"}},
	{"17013", "JB", "DME", cgs, []string{"IL", "IN", "KY", "MI", "MN", "OH", "WI"}},
	{"18003", "JC", "DME", cgs, []string{"AL", "AR", "CO", "FL", "GA", "LA", "MS", "NC", "NM", "OK", "PR", "SC", "TN", "TX", "VA", "VI", "WV"}},
	{"19003", "JD", "DME", noridian, []string{"AK", "AS", "AZ", "CA", "GU", "HI", "IA", "ID", "KS", "MO", "MP", "MT", "ND", "NE", "NV", "OR", "SD", "UT", "WA", "WY"}},
}

// byID indexes Workloads by ID.
var byID = indexByID()

func indexByID() map[string]int {
	m := make(map[string]int, len(Workloads))
	for i, w := range Workloads {
		m[w.ID] = i
	}
	return m
}

// Lookup returns the row for id, the five-digit number exactly as the file
// carries it; trim padding first. It reports false for a number the table
// does not hold. The returned States slice is shared with the table; do not
// modify it.
func Lookup(id string) (Workload, bool) {
	i, ok := byID[id]
	if !ok {
		return Workload{}, false
	}
	return Workloads[i], true
}
