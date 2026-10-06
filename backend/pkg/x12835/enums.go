package x12835

// Coded elements
//
// Every coded element in this package (ClaimStatusCLP02, AuthQualifierISA01,
// AdjGroupCAS01, and the rest of enums_envelope.go, enums_payment.go,
// enums_party.go, enums_claim.go, enums_reference.go, enums_adjustment.go,
// and enums_service.go) follows one shape:
//
//   - The type is a string, and a value is the code's text exactly as the
//     file carried it. The zero value, "", is the one empty: the element
//     was not sent (nothing between its separators, or the segment ended
//     before it; X12 gives those two the same meaning). Every type names
//     it with a <Type>NotSent constant.
//   - Each code the 835 implementation guide (X12 005010X221A1) allows is a
//     named constant whose value is the code ("1", "HP", "CO"). Codes the
//     reviewed CMS documents name are named by that name; the rest by
//     their code. Where the guide allows a code and Medicare's documents
//     say Medicare does not send it, the constant is there with the note.
//   - One list beside each type, <type>Codes, holds every allowed code
//     with its public name ("" when no reviewed source names it) and that
//     name with its citation (codes.go, codeEntry). Known reports whether
//     a value is on the list; ShortDescription and Description read the
//     list and return "" for a value that is not.
//   - A value off the list is kept as it came. Nothing here drops,
//     rewrites, or replaces a received code, so a file reads in and writes
//     out byte for byte (write.go); whether a code is one the guide allows
//     is a question for the caller (Known), never a parse failure.
//
// Parse<Type> reads received text into the type: the text as it came, and
// an error naming the element when the text is not empty and not on the
// list. The parser (parse_segments.go) converts directly, without Parse,
// since it keeps every value. String is the identity, and encoding/json
// reads and writes the code's text with no methods of its own.
//
// CARC and RARC (CASReason.Reason, the MIA/MOA remark positions, RemarkLQ)
// are coded elements too, but their lists live in their own packages,
// backend/pkg/carc and backend/pkg/rarc, not here: X12 and CMS publish and
// retire hundreds of those codes on their own schedule, three times a year,
// so each package's entries table is updated from backend/pkg/carc.json or
// rarc.json by tools/codegen (make codes) rather than by hand. carc.Code and
// rarc.Code follow this same shape, except that their codes have no named
// constants; CodeNotSent is the one constant each package names.
