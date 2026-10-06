package x12835

// Amount is an optional exact decimal: many 835 numeric elements are not
// sent at all, and that absence is a different fact than a sent 0.00 (an
// 835 that pays a claim $0.00 still sends "0"). Set is false when the
// element was not sent; Value is meaningless in that case. Set is true
// with Value holding the parsed Decimal otherwise, including an exact
// zero. Decimal itself (backend/pkg/x12835/decimal.go) is an exact
// scalar, Coefficient / 10^Scale, with no floating point.
//
// Text is the element's text as the file carried it ("180.00", "-0.50"),
// kept so the file writes back out byte for byte; Value is what compares.
// An Amount built by hand may leave Text empty, and then it writes as
// Value.String().
type Amount struct {
	Value Decimal `json:"value,omitempty"`
	Set   bool    `json:"set,omitempty"`
	Text  string  `json:"text,omitempty"`
}
