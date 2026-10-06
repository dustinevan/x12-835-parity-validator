package x12835

import (
	"fmt"
	"math"
	"strconv"
)

// Decimal is an exact scalar: Coefficient / 10^Scale, with no floating
// point. ParseDecimal removes trailing fractional zeros and Sub reduces its
// result the same way, so equal values have identical representations and
// two Decimals compare with ==: 10.0 and 10.00 are both {10, 0}.
type Decimal struct {
	Coefficient int64
	Scale       int
}

// ParseDecimal reads an X12 decimal: an optional sign, digits, and at most
// one decimal point. Trailing fractional zeros are removed, and a zero has
// scale 0, so equal values parse to equal Decimals. Identifiers and codes
// are never parsed with it: whether an element is numeric is decided by the
// element's definition, not by the look of its value.
func ParseDecimal(data []byte) (Decimal, error) {
	if len(data) == 0 {
		return Decimal{}, fmt.Errorf("x12: parsing decimal: empty value")
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
			return Decimal{}, fmt.Errorf("x12: parsing decimal: invalid character")
		}
		digits++
	}
	if digits == 0 {
		return Decimal{}, fmt.Errorf("x12: parsing decimal: missing digits")
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
			return Decimal{}, fmt.Errorf("x12: parsing decimal: coefficient exceeds int64")
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

// String writes the value as X12 writes a number: the digits, a decimal
// point when Scale is above zero, and a leading minus when negative. A
// parsed value has no trailing fractional zeros, so 447.00 prints "447",
// 1469.50 "1469.5", 0.50 "0.5", and -0.01 "-0.01".
func (d Decimal) String() string {
	negative := d.Coefficient < 0
	digits := strconv.FormatInt(d.Coefficient, 10)
	if negative {
		digits = digits[1:]
	}
	if d.Scale > 0 {
		for len(digits) <= d.Scale {
			digits = "0" + digits
		}
		digits = digits[:len(digits)-d.Scale] + "." + digits[len(digits)-d.Scale:]
	}
	if negative {
		return "-" + digits
	}
	return digits
}

// Sub returns d minus e, exactly. The subtraction runs at the larger of the
// two scales, so 80 minus 79.99 is -0.01, and the result is reduced as a
// parsed value is: 10.5 minus 0.5 is {10, 0}, "10". It returns an error
// when a coefficient brought to the larger scale, or the difference, does
// not fit in int64.
func (d Decimal) Sub(e Decimal) (Decimal, error) {
	scale := max(d.Scale, e.Scale)
	a, err := scaled(d, scale)
	if err != nil {
		return Decimal{}, err
	}
	b, err := scaled(e, scale)
	if err != nil {
		return Decimal{}, err
	}
	difference := a - b
	if (b > 0 && difference > a) || (b < 0 && difference < a) {
		return Decimal{}, fmt.Errorf("x12: subtracting decimals: difference exceeds int64")
	}
	return reduce(Decimal{Coefficient: difference, Scale: scale}), nil
}

// Add returns d plus e, exactly: the sum at the larger of the two scales,
// reduced as Sub's result is, so 392.04 plus 90.48 is {48252, 2}, "482.52",
// and 0.5 plus 10.5 is {11, 0}, "11". It returns an error when a
// coefficient brought to the larger scale, or the sum, does not fit in
// int64.
func (d Decimal) Add(e Decimal) (Decimal, error) {
	scale := max(d.Scale, e.Scale)
	a, err := scaled(d, scale)
	if err != nil {
		return Decimal{}, err
	}
	b, err := scaled(e, scale)
	if err != nil {
		return Decimal{}, err
	}
	sum := a + b
	if (b > 0 && sum < a) || (b < 0 && sum > a) {
		return Decimal{}, fmt.Errorf("x12: adding decimals: sum exceeds int64")
	}
	return reduce(Decimal{Coefficient: sum, Scale: scale}), nil
}

// scaled is d's coefficient written at a scale at least as large as its own.
func scaled(d Decimal, scale int) (int64, error) {
	c := d.Coefficient
	for s := d.Scale; s < scale; s++ {
		if c > math.MaxInt64/10 || c < math.MinInt64/10 {
			return 0, fmt.Errorf("x12: scaling decimal: coefficient exceeds int64")
		}
		c *= 10
	}
	return c, nil
}

// reduce removes trailing fractional zeros, so the result reads as
// ParseDecimal would parse its text.
func reduce(d Decimal) Decimal {
	for d.Scale > 0 && d.Coefficient%10 == 0 {
		d.Coefficient /= 10
		d.Scale--
	}
	if d.Coefficient == 0 {
		d.Scale = 0
	}
	return d
}
