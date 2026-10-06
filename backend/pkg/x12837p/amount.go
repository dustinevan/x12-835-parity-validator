package x12837p

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Decimal is an exact scalar: Coefficient / 10^Scale. Trailing fractional
// zeros are removed when parsing, so equal values have identical
// representations. Shape copied from backend/pkg/x12835 (decimal.go);
// kept independent here rather than imported so this package does not
// depend on x12835.
type Decimal struct {
	Coefficient int64 `json:"coefficient,omitempty"`
	Scale       int   `json:"scale,omitempty"`
}

// String returns the decimal in its canonical spelling: an optional minus
// sign, the digits, and a decimal point only when the scale is positive
// ("80", "-0.5", "104.63").
func (d Decimal) String() string {
	negative := d.Coefficient < 0
	digits := strconv.FormatInt(d.Coefficient, 10)
	if negative {
		digits = digits[1:]
	}
	if d.Scale > 0 {
		if len(digits) <= d.Scale {
			digits = strings.Repeat("0", d.Scale-len(digits)+1) + digits
		}
		digits = digits[:len(digits)-d.Scale] + "." + digits[len(digits)-d.Scale:]
	}
	if negative {
		return "-" + digits
	}
	return digits
}

// Amount is an optional exact decimal: many 837P numeric elements (CLM02,
// the CAS amount/quantity positions, AMT02, SV102, ...) are not sent at
// all, and that absence is a different fact than a sent 0.00. Set is false
// when the element was not sent; Value and Text are meaningless in that
// case. Set is true with Value holding the parsed Decimal otherwise,
// including an exact zero, and Text holding the element's text exactly as
// the file sent it ("175.50" where Value reads 175.5), so the file writes
// back with its own spelling. An Amount built in code may leave Text empty;
// the writer then spells Value.
type Amount struct {
	Value Decimal `json:"value,omitempty"`
	Text  string  `json:"text,omitempty"`
	Set   bool    `json:"set,omitempty"`
}

// ParseDecimal reads an X12 decimal (an optional sign, digits, and at most
// one decimal point) into an exact Decimal, removing trailing fractional
// zeros so that equal values have identical representations. Copied from
// backend/pkg/x12835's ParseDecimal.
func ParseDecimal(data []byte) (Decimal, error) {
	if len(data) == 0 {
		return Decimal{}, fmt.Errorf("x12837p: parsing decimal: empty value")
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
			return Decimal{}, fmt.Errorf("x12837p: parsing decimal: invalid character")
		}
		digits++
	}
	if digits == 0 {
		return Decimal{}, fmt.Errorf("x12837p: parsing decimal: missing digits")
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
			return Decimal{}, fmt.Errorf("x12837p: parsing decimal: coefficient exceeds int64")
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
