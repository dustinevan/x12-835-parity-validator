package x12835

// RefQualifierREF01 is REF01, the Reference Identification Qualifier, shared
// by every REF segment in the transaction (header, payer/payee loop, claim
// loop, and service line).
//
// TJ has two sourced names: at the payee loop
// (1000B) CMS 835 Flat File, PDF p.4 and frontend/src/data/x12Names835.ts's
// PAYEE_LOOP_NAME name REF qualifier TJ "Payee Additional Identification"
// (confirmed as the only REF01 TJ usage in this repo's example 835 files),
// while docs/design/835-examples.md's S-REF-RENDERING row lists TJ as one of
// four service-line (2110) qualifiers - HPI, SY, TJ, 1C - sharing the name
// "Rendering Provider Information" (CEDI 835 CG Table 21). One code, two
// sourced meanings depending on loop; as with ContactFunctionPER01CX,
// RefQualifierREF01TJ is named by its code rather than picking one sourced
// meaning over the other.
type RefQualifierREF01 string

const (
	RefQualifierREF01NotSent RefQualifierREF01 = ""
	RefQualifierREF01EV      RefQualifierREF01 = "EV"
	RefQualifierREF01F2      RefQualifierREF01 = "F2"
	RefQualifierREF012U      RefQualifierREF01 = "2U"
	RefQualifierREF01TJ      RefQualifierREF01 = "TJ"
	RefQualifierREF01EA      RefQualifierREF01 = "EA"
	RefQualifierREF016P      RefQualifierREF01 = "6P"
	RefQualifierREF0128      RefQualifierREF01 = "28"
	RefQualifierREF01F8      RefQualifierREF01 = "F8"
	RefQualifierREF01CE      RefQualifierREF01 = "CE"
	RefQualifierREF016R      RefQualifierREF01 = "6R"
	RefQualifierREF01LU      RefQualifierREF01 = "LU"
	RefQualifierREF011S      RefQualifierREF01 = "1S"
	RefQualifierREF01APC     RefQualifierREF01 = "APC"
	RefQualifierREF01RB      RefQualifierREF01 = "RB"
	RefQualifierREF01HPI     RefQualifierREF01 = "HPI"
	RefQualifierREF01SY      RefQualifierREF01 = "SY"
	RefQualifierREF011C      RefQualifierREF01 = "1C"
	RefQualifierREF010K      RefQualifierREF01 = "0K"
	RefQualifierREF010B      RefQualifierREF01 = "0B"
	RefQualifierREF011A      RefQualifierREF01 = "1A"
	RefQualifierREF011B      RefQualifierREF01 = "1B"
	RefQualifierREF011D      RefQualifierREF01 = "1D"
	RefQualifierREF011G      RefQualifierREF01 = "1G"
	RefQualifierREF011H      RefQualifierREF01 = "1H"
	RefQualifierREF011J      RefQualifierREF01 = "1J"
	RefQualifierREF011L      RefQualifierREF01 = "1L"
	RefQualifierREF011W      RefQualifierREF01 = "1W"
	RefQualifierREF019A      RefQualifierREF01 = "9A"
	RefQualifierREF019C      RefQualifierREF01 = "9C"
	RefQualifierREF01BB      RefQualifierREF01 = "BB"
	RefQualifierREF01D3      RefQualifierREF01 = "D3"
	RefQualifierREF01E9      RefQualifierREF01 = "E9"
	RefQualifierREF01EO      RefQualifierREF01 = "EO"
	RefQualifierREF01G1      RefQualifierREF01 = "G1"
	RefQualifierREF01G2      RefQualifierREF01 = "G2"
	RefQualifierREF01G3      RefQualifierREF01 = "G3"
	RefQualifierREF01HI      RefQualifierREF01 = "HI"
	RefQualifierREF01IG      RefQualifierREF01 = "IG"
	RefQualifierREF01NF      RefQualifierREF01 = "NF"
	RefQualifierREF01PQ      RefQualifierREF01 = "PQ"
)

// refQualifierREF01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var refQualifierREF01Codes = []codeEntry{
	{"EV", "Receiver Identification", "Receiver Identification — CMS 835 Flat File, PDF p.2 (header); frontend/src/data/x12Names835.ts HEADER_REF_NAME"},
	{"F2", "Version Identification", "Version Identification — CMS 835 Flat File, PDF p.2 (header); frontend/src/data/x12Names835.ts HEADER_REF_NAME"},
	{"2U", "Additional Payer Identification", "Additional Payer Identification — CMS 835 Flat File, PDF p.3 (payer loop, 1000A); CEDI 835 CG, PDF p.24 (printed p.18; Table 7); docs/design/835-examples.md PR-REF"},
	{"TJ", "Federal Taxpayer's Identification Number", "Federal Taxpayer's Identification Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide); ambiguous, see the type doc comment; CMS 835 Flat File, PDF p.4"},
	{"EA", "Other Claim Related Identification", "Other Claim Related Identification — CMS 835 Flat File, PDF p.10 (claim loop, 2100); CEDI 835 CG, PDF p.26 (printed p.20; Table 14); frontend/src/data/x12Names835.ts CLAIM_REF_NAME"},
	{"6P", "Other Claim Related Identification", "Other Claim Related Identification — same grouping and citation as EA (frontend/src/data/x12Names835.ts's comment: \"qualifiers EA/6P/28/F8/CE share one combined name\")"},
	{"28", "Other Claim Related Identification", "Other Claim Related Identification — same grouping and citation as EA"},
	{"F8", "Other Claim Related Identification", "Other Claim Related Identification — same grouping and citation as EA; CEDI 835 CG's own gloss: Original Reference Number, used to link a correction to the claim it corrects (docs/design/835-examples.md's reversal-and-correction example)"},
	{"CE", "Other Claim Related Identification", "Other Claim Related Identification — same grouping and citation as EA"},
	{"6R", "Line Item Control Number", "Line Item Control Number — CMS 835 Flat File, PDF p.13 (service line, 2110); frontend/src/data/x12Names835.ts SERVICE_REF_NAME"},
	{"LU", "Service Identification", "Service Identification — CMS 835 Flat File, PDF p.12 (service line, 2110); CEDI 835 CG, PDF p.27 (printed p.21; Table 20); frontend/src/data/x12Names835.ts SERVICE_REF_NAME"},
	{"1S", "Service Identification", "Service Identification — same grouping and citation as LU (docs/design/835-examples.md S-REF-SVCID lists qualifiers LU, 1S, APC, RB together)"},
	{"APC", "Service Identification", "Service Identification — same grouping and citation as LU"},
	{"RB", "Service Identification", "Service Identification — same grouping and citation as LU"},
	{"HPI", "Rendering Provider Information", "Rendering Provider Information — CMS 835 Flat File, PDF p.13 (service line, 2110); CEDI 835 CG, PDF p.28 (printed p.22; Table 21); frontend/src/data/x12Names835.ts SERVICE_REF_NAME"},
	{"SY", "Rendering Provider Information", "Rendering Provider Information — same grouping and citation as HPI (docs/design/835-examples.md S-REF-RENDERING lists qualifiers HPI, SY, TJ, 1C together)"},
	{"1C", "Rendering Provider Information", "Rendering Provider Information — same grouping and citation as HPI"},
	{"0K", "Health Care Policy Identification", "Health Care Policy Identification — CMS 835 Flat File, PDF p.13 (service line, 2110); CEDI 835 CG, PDF p.28 (printed p.22; Table 22); frontend/src/data/x12Names835.ts SERVICE_REF_NAME"},
	{"0B", "State License Number", "State License Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1A", "Blue Cross Provider Number", "Blue Cross Provider Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1B", "Blue Shield Provider Number", "Blue Shield Provider Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1D", "Medicaid Provider Number", "Medicaid Provider Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1G", "Provider UPIN Number", "Provider UPIN Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1H", "CHAMPUS Identification Number", "CHAMPUS Identification Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1J", "Facility ID Number", "Facility ID Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1L", "Group or Policy Number", "Group or Policy Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"1W", "Member Identification Number", "Member Identification Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"9A", "Repriced Claim Reference Number", "Repriced Claim Reference Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"9C", "Adjusted Repriced Claim Reference Number", "Adjusted Repriced Claim Reference Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"BB", "Authorization Number", "Authorization Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"D3", "National Council for Prescription Drug Programs Pharmacy Number", "National Council for Prescription Drug Programs Pharmacy Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"E9", "Attachment Code", "Attachment Code — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"EO", "Submitter Identification Number", "Submitter Identification Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"G1", "Prior Authorization Number", "Prior Authorization Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"G2", "Provider Commercial Number", "Provider Commercial Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"G3", "Predetermination of Benefits Identification Number", "Predetermination of Benefits Identification Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"HI", "Health Industry Number (HIN)", "Health Industry Number (HIN) — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"IG", "Insurance Policy Number", "Insurance Policy Number — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"NF", "National Association of Insurance Commissioners (NAIC) Code", "National Association of Insurance Commissioners (NAIC) Code — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
	{"PQ", "Payee Identification", "Payee Identification — X12 005010X221A1 TR3, REF01 (transcribed from the standard; verify against the printed guide)"},
}

func (v RefQualifierREF01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v RefQualifierREF01) Known() bool { return known(refQualifierREF01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v RefQualifierREF01) ShortDescription() string {
	return shortOf(refQualifierREF01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v RefQualifierREF01) Description() string {
	return descriptionOf(refQualifierREF01Codes, string(v))
}

// ParseRefQualifierREF01 reads one REF01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseRefQualifierREF01(code string) (RefQualifierREF01, error) {
	return RefQualifierREF01(code), checkListed("REF01", refQualifierREF01Codes, code)
}

// DateQualifierDTM01 is DTM01, the Date/Time Qualifier, shared by every DTM
// segment in the transaction (header, claim loop, and service line) - CMS
// 835 Flat File, PDF p.2, 10-12 gives every DTM01 instance the same position
// name.
type DateQualifierDTM01 string

const (
	DateQualifierDTM01NotSent DateQualifierDTM01 = ""
	DateQualifierDTM01405     DateQualifierDTM01 = "405"
	DateQualifierDTM01050     DateQualifierDTM01 = "050"
	DateQualifierDTM01232     DateQualifierDTM01 = "232"
	DateQualifierDTM01233     DateQualifierDTM01 = "233"
	DateQualifierDTM01036     DateQualifierDTM01 = "036"
	DateQualifierDTM01472     DateQualifierDTM01 = "472"
	DateQualifierDTM01150     DateQualifierDTM01 = "150"
	DateQualifierDTM01151     DateQualifierDTM01 = "151"
)

// dateQualifierDTM01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var dateQualifierDTM01Codes = []codeEntry{
	{"405", "Production Date", "Production Date — CMS 835 Flat File, PDF p.2 (header); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"050", "Claim Received Date", "Claim Received Date — CMS 835 Flat File, PDF p.11 (claim loop, 2100); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"232", "Statement From-To Date", "Statement From-To Date — CMS 835 Flat File, PDF p.10 (claim loop, 2100); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"233", "Statement From-To Date", "Statement From-To Date — CMS 835 Flat File, PDF p.10 (claim loop, 2100); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"036", "Coverage Expiration Date", "Coverage Expiration Date — CMS 835 Flat File, PDF p.10 (claim loop, 2100); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"472", "Service Date", "Service Date — CMS 835 Flat File, PDF p.12 (service line, 2110); frontend/src/data/x12Names835.ts DTM_NAME"},
	{"150", "Service Date", "Service Date — CMS 835 Flat File, PDF p.12 (service line, 2110) names the DTM segment and DTM02 \"Service Date\" and lists 150, 151, 472 as its DTM01 values; the page names no qualifier individually"},
	{"151", "Service Date", "Service Date — same page, segment, and wording as 150"},
}

func (v DateQualifierDTM01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v DateQualifierDTM01) Known() bool { return known(dateQualifierDTM01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v DateQualifierDTM01) ShortDescription() string {
	return shortOf(dateQualifierDTM01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v DateQualifierDTM01) Description() string {
	return descriptionOf(dateQualifierDTM01Codes, string(v))
}

// ParseDateQualifierDTM01 reads one DTM01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseDateQualifierDTM01(code string) (DateQualifierDTM01, error) {
	return DateQualifierDTM01(code), checkListed("DTM01", dateQualifierDTM01Codes, code)
}

// AmountQualifierAMT01 is AMT01, the Amount Qualifier Code, shared by the
// claim-loop (2100) and service-line (2110) AMT segments.
//
// docs/design/835-examples.md's C-AMT row lists the claim-level set (AU, DY,
// F5, I, NL, ZK, ZL, ZM, ZN, ZO) and its S-AMT row the service-level set
// (B6, KH, ZK, ZL, ZM, ZN, ZO); this type is their union, since both levels
// use the same AMT01 code list. Neither reviewed document names any value in
// text.
type AmountQualifierAMT01 string

const (
	AmountQualifierAMT01NotSent AmountQualifierAMT01 = ""
	AmountQualifierAMT01AU      AmountQualifierAMT01 = "AU"
	AmountQualifierAMT01B6      AmountQualifierAMT01 = "B6"
	AmountQualifierAMT01DY      AmountQualifierAMT01 = "DY"
	AmountQualifierAMT01F5      AmountQualifierAMT01 = "F5"
	AmountQualifierAMT01I       AmountQualifierAMT01 = "I"
	AmountQualifierAMT01KH      AmountQualifierAMT01 = "KH"
	AmountQualifierAMT01NL      AmountQualifierAMT01 = "NL"
	AmountQualifierAMT01ZK      AmountQualifierAMT01 = "ZK"
	AmountQualifierAMT01ZL      AmountQualifierAMT01 = "ZL"
	AmountQualifierAMT01ZM      AmountQualifierAMT01 = "ZM"
	AmountQualifierAMT01ZN      AmountQualifierAMT01 = "ZN"
	AmountQualifierAMT01ZO      AmountQualifierAMT01 = "ZO"
	AmountQualifierAMT01D8      AmountQualifierAMT01 = "D8"
	AmountQualifierAMT01T       AmountQualifierAMT01 = "T"
	AmountQualifierAMT01T2      AmountQualifierAMT01 = "T2"
)

// amountQualifierAMT01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var amountQualifierAMT01Codes = []codeEntry{
	{"AU", "Coverage Amount", "Coverage Amount — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop); not named in text (docs/design/835-examples.md C-AMT)"},
	{"B6", "Allowed - Actual", "Allowed - Actual — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.13 (service line); not named in text (docs/design/835-examples.md S-AMT; frontend/src/data/x12Names835.ts's own comment calls out \"AMT qualifier B6\" as unresolved)"},
	{"DY", "Per Day Limit", "Per Day Limit — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop); not named in text"},
	{"F5", "Patient Amount Paid", "Patient Amount Paid — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop); not named in text"},
	{"I", "Interest", "Interest — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop); not named in text. docs/design/835-examples.md's C-AMT-I row glosses this as interest \"by the element's own context,\" but says so explicitly is \"not spelled out in text\" - that gloss is this project's own inference, not a published name, so it is not used as this constant's name"},
	{"KH", "Deduction Amount", "Deduction Amount — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.13 (service line); not named in text"},
	{"NL", "Negative Ledger Balance", "Negative Ledger Balance — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop); not named in text"},
	{"ZK", "Federal Medicare or Medicaid Payment Mandate - Category 1", "Federal Medicare or Medicaid Payment Mandate - Category 1 — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service line); not named in text"},
	{"ZL", "Federal Medicare or Medicaid Payment Mandate - Category 2", "Federal Medicare or Medicaid Payment Mandate - Category 2 — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service line); not named in text"},
	{"ZM", "Federal Medicare or Medicaid Payment Mandate - Category 3", "Federal Medicare or Medicaid Payment Mandate - Category 3 — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service line); not named in text"},
	{"ZN", "Federal Medicare or Medicaid Payment Mandate - Category 4", "Federal Medicare or Medicaid Payment Mandate - Category 4 — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service line); not named in text"},
	{"ZO", "Federal Medicare or Medicaid Payment Mandate - Category 5", "Federal Medicare or Medicaid Payment Mandate - Category 5 — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11 (claim loop), p.13 (service line); not named in text"},
	{"D8", "Discount Amount", "Discount Amount — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide)"},
	{"T", "Tax", "Tax — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide)"},
	{"T2", "Total Claim Before Taxes", "Total Claim Before Taxes — X12 005010X221A1 TR3, AMT01 (transcribed from the standard; verify against the printed guide)"},
}

func (v AmountQualifierAMT01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v AmountQualifierAMT01) Known() bool { return known(amountQualifierAMT01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v AmountQualifierAMT01) ShortDescription() string {
	return shortOf(amountQualifierAMT01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v AmountQualifierAMT01) Description() string {
	return descriptionOf(amountQualifierAMT01Codes, string(v))
}

// ParseAmountQualifierAMT01 reads one AMT01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseAmountQualifierAMT01(code string) (AmountQualifierAMT01, error) {
	return AmountQualifierAMT01(code), checkListed("AMT01", amountQualifierAMT01Codes, code)
}

// QuantityQualifierQTY01 is QTY01, the Quantity Qualifier on the claim-loop
// (2100) QTY segment (the service-line QTY segment, S-QTY, is documented as
// not used by Medicare).
type QuantityQualifierQTY01 string

const (
	QuantityQualifierQTY01NotSent QuantityQualifierQTY01 = ""
	QuantityQualifierQTY01CA      QuantityQualifierQTY01 = "CA"
	QuantityQualifierQTY01CD      QuantityQualifierQTY01 = "CD"
	QuantityQualifierQTY01LA      QuantityQualifierQTY01 = "LA"
	QuantityQualifierQTY01DU      QuantityQualifierQTY01 = "DU"
	QuantityQualifierQTY01ZK      QuantityQualifierQTY01 = "ZK"
	QuantityQualifierQTY01ZL      QuantityQualifierQTY01 = "ZL"
	QuantityQualifierQTY01ZM      QuantityQualifierQTY01 = "ZM"
	QuantityQualifierQTY01ZN      QuantityQualifierQTY01 = "ZN"
	QuantityQualifierQTY01ZO      QuantityQualifierQTY01 = "ZO"
	QuantityQualifierQTY01LE      QuantityQualifierQTY01 = "LE"
	QuantityQualifierQTY01NE      QuantityQualifierQTY01 = "NE"
	QuantityQualifierQTY01NR      QuantityQualifierQTY01 = "NR"
	QuantityQualifierQTY01OU      QuantityQualifierQTY01 = "OU"
	QuantityQualifierQTY01PS      QuantityQualifierQTY01 = "PS"
	QuantityQualifierQTY01VS      QuantityQualifierQTY01 = "VS"
)

// quantityQualifierQTY01Codes lists the codes the element may carry, each with its public
// name ("" when the sources give none) and the name with its citation.
var quantityQualifierQTY01Codes = []codeEntry{
	{"CA", "Covered - Actual", "Covered - Actual — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text (docs/design/835-examples.md C-QTY)"},
	{"CD", "Co-insured - Actual", "Co-insured - Actual — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"LA", "Life-time Reserve - Actual", "Life-time Reserve - Actual — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"DU", "", "CMS 835 Flat File, PDF p.11; not named in text"},
	{"ZK", "Federal Medicare or Medicaid Payment Mandate - Category 1", "Federal Medicare or Medicaid Payment Mandate - Category 1 — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"ZL", "Federal Medicare or Medicaid Payment Mandate - Category 2", "Federal Medicare or Medicaid Payment Mandate - Category 2 — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"ZM", "Federal Medicare or Medicaid Payment Mandate - Category 3", "Federal Medicare or Medicaid Payment Mandate - Category 3 — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"ZN", "Federal Medicare or Medicaid Payment Mandate - Category 4", "Federal Medicare or Medicaid Payment Mandate - Category 4 — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"ZO", "Federal Medicare or Medicaid Payment Mandate - Category 5", "Federal Medicare or Medicaid Payment Mandate - Category 5 — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide); CMS 835 Flat File, PDF p.11; not named in text"},
	{"LE", "Life-time Reserve - Estimated", "Life-time Reserve - Estimated — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
	{"NE", "Non-Covered - Estimated", "Non-Covered - Estimated — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
	{"NR", "Not Replaced Blood Units", "Not Replaced Blood Units — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
	{"OU", "Outlier Days", "Outlier Days — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
	{"PS", "Prescription", "Prescription — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
	{"VS", "Visits", "Visits — X12 005010X221A1 TR3, QTY01 (transcribed from the standard; verify against the printed guide)"},
}

func (v QuantityQualifierQTY01) String() string { return string(v) }

// Known reports whether the value is on the implementation guide's list.
func (v QuantityQualifierQTY01) Known() bool { return known(quantityQualifierQTY01Codes, string(v)) }

// ShortDescription returns the value's public name, or "" when the
// reviewed sources give it none, or when the value is not on the list.
func (v QuantityQualifierQTY01) ShortDescription() string {
	return shortOf(quantityQualifierQTY01Codes, string(v))
}

// Description returns the value's public name with its citation, or the
// citation note alone when the value has no public name; "" off the list.
func (v QuantityQualifierQTY01) Description() string {
	return descriptionOf(quantityQualifierQTY01Codes, string(v))
}

// ParseQuantityQualifierQTY01 reads one QTY01 value: the text as it came, and an error
// when it is not empty and not on the list. The value is kept either way.
func ParseQuantityQualifierQTY01(code string) (QuantityQualifierQTY01, error) {
	return QuantityQualifierQTY01(code), checkListed("QTY01", quantityQualifierQTY01Codes, code)
}
