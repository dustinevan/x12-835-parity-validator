package x12837i

import (
	"fmt"
	"math"
)

// Decimal is an exact scalar: Coefficient / 10^Scale. Trailing fractional
// zeros are removed when parsing, so equal values have identical
// representations. Copied in shape from backend/pkg/x12835 (decimal.go):
// no floating point anywhere in this package.
type Decimal struct {
	Coefficient int64 `json:"coefficient,omitempty"`
	Scale       int   `json:"scale,omitempty"`
}

// ParseDecimal parses one X12 numeric element (for example CLM02 or SV203)
// into an exact Decimal. It keeps no state.
func ParseDecimal(data []byte) (Decimal, error) {
	if len(data) == 0 {
		return Decimal{}, fmt.Errorf("x12837i: parsing decimal: empty value")
	}
	negative := data[0] == '-'
	if negative || data[0] == '+' {
		data = data[1:]
	}
	decimalPoint := -1
	digits := 0
	for i, b := range data {
		if b == '.' && decimalPoint == -1 {
			decimalPoint = i
			continue
		}
		if b < '0' || b > '9' {
			return Decimal{}, fmt.Errorf("x12837i: parsing decimal: invalid character")
		}
		digits++
	}
	if digits == 0 {
		return Decimal{}, fmt.Errorf("x12837i: parsing decimal: missing digits")
	}
	var value Decimal
	if decimalPoint >= 0 {
		for len(data) > decimalPoint+1 && data[len(data)-1] == '0' {
			data = data[:len(data)-1]
		}
		value.Scale = len(data) - decimalPoint - 1
	}
	for _, b := range data {
		if b == '.' {
			continue
		}
		digit := int64(b - '0')
		if value.Coefficient > (math.MaxInt64-digit)/10 {
			return Decimal{}, fmt.Errorf("x12837i: parsing decimal: coefficient exceeds int64")
		}
		value.Coefficient = value.Coefficient*10 + digit
	}
	if negative {
		value.Coefficient = -value.Coefficient
	}
	if value.Coefficient == 0 {
		value.Scale = 0
	}
	return value, nil
}

// String writes the decimal in plain form: the digits with a decimal point
// when Scale is more than zero, a leading minus when negative, and no
// trailing fractional zeros (12500, 412.5, -0.5). Write uses it only for an
// Amount built without the file's own text.
func (d Decimal) String() string {
	negative := d.Coefficient < 0
	digits := fmt.Sprintf("%d", d.Coefficient)
	if negative {
		digits = digits[1:]
	}
	for len(digits) <= d.Scale {
		digits = "0" + digits
	}
	if d.Scale > 0 {
		digits = digits[:len(digits)-d.Scale] + "." + digits[len(digits)-d.Scale:]
	}
	if negative {
		return "-" + digits
	}
	return digits
}

// Amount is an optional exact decimal: many 837I numeric elements are not
// sent at all, and that absence is a different fact than a sent 0.00 (a
// claim with a true zero line-item charge still sends "0"). Set is false
// when the element was not sent; Value and Text are meaningless in that
// case. Set is true otherwise, with Value holding the parsed Decimal
// (including an exact zero) and Text the element exactly as the file sent
// it ("12500.00" stays "12500.00", while Value reads 12500), so the file
// is written back with the spelling it came with. Copied in shape from
// backend/pkg/x12835/amount.go, with Text added.
type Amount struct {
	Value Decimal `json:"value,omitempty"`
	Text  string  `json:"text,omitempty"`
	Set   bool    `json:"set,omitempty"`
}
