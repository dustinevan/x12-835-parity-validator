package x12837i

// The coded-element scheme
//
// Every coded element in this package (TransactionTypeGS01_ST01 and the rest
// of enums_envelope.go, enums_provider.go, enums_subscriber.go,
// enums_claim.go, enums_institutional.go, enums_service.go,
// enums_reference.go, and enums_cob.go) is a string type whose value is the
// code exactly as the file sent it:
//
//   - The zero value "" is the one empty: the element carried no value,
//     whether nothing stood between its separators or the segment ended
//     before it. X12 gives those two the same meaning.
//   - Each code the 837I implementation guide allows gets a named constant
//     whose value is its code text, in the order its citation lists it. A
//     value with no single, unambiguous public name is named by its own
//     code instead of a guessed English name - see DateQualifierDTP01573,
//     present on both the 2330B and 2430 DTP with two different sourced
//     names for the same code.
//   - A code off the list is kept as itself. Nothing is lost and nothing
//     is renamed: String returns the text, JSON carries the text, and
//     Known reports false. The parser never rejects a file for a code it
//     does not know; the guide's lists, Medicare's narrower usage, and a
//     sender's mistakes are all left for the caller to judge.
//
// One list beside each type, <type>Codes, is the one place its text lives:
// the code each constant stands for, its public name ("" when the reviewed
// sources give none), and that name with its citation. Known,
// ShortDescription, and Description read it; a code off the list has no
// name and no description.
//
// Parse<Type> turns received text into the type: "" and a known code come
// back with no error, and any other text comes back as itself with an
// error, so a caller that wants to enforce the list can, and a caller that
// wants the text keeps it either way. The parser converts directly and
// keeps every value.
//
// Five elements - CLM05-1 (facility type), CLM05-3 (claim frequency), CL101
// (admission type), CL102 (admission source), and CL103 (patient status) -
// are National Uniform Billing Committee (NUBC) code sets that the guide
// defers to an external list: Transmittal 907's "Values" column is blank
// for them, and CMS Claims Processing Manual Ch. 25 defers to Medicare
// contractors and the NUBC's Official UB-04 Data Specifications Manual
// (see the NUBC note at the top of frontend/src/data/x12Names837i.ts).
// Those five are plain string fields with no type of their own and no
// list, so a value there is kept as sent and never judged.
//
// CARC and RARC (the Claim Adjustment Reason Codes on this transaction's
// CAS segments, and the Remittance Advice Remark Codes on MOA's and MIA's
// Claim Payment Remark Code positions - both confined to the 2320/2430
// coordination-of-benefits loops in the 837I) are coded elements too, but
// their lists are not duplicated here: this package reuses
// backend/pkg/carc.Code and backend/pkg/rarc.Code directly, the same
// string types the 835 side uses, rather than maintaining a
// second copy of either code list.

// codeEntry is one code of a coded element's list: the code text, its
// public name ("" when the reviewed sources give none), and that name with
// its citation.
type codeEntry struct {
	code        string
	short       string
	description string
}

// knownCode reports whether code is on the list.
func knownCode(list []codeEntry, code string) bool {
	for _, e := range list {
		if e.code == code {
			return true
		}
	}
	return false
}

// shortOf returns the public name of code, or "" when the list has no
// entry for it or gives it no name.
func shortOf(list []codeEntry, code string) string {
	for _, e := range list {
		if e.code == code {
			return e.short
		}
	}
	return ""
}

// descriptionOf returns the public name of code with its citation, or ""
// when the list has no entry for it.
func descriptionOf(list []codeEntry, code string) string {
	for _, e := range list {
		if e.code == code {
			return e.description
		}
	}
	return ""
}
