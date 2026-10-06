package x12835

import (
	"math"
	"testing"
)

func TestDecimalEquality(t *testing.T) {
	for _, pair := range [][2]string{{"10.0", "10.00"}, {"-0", "0.00"}, {".50", "0.5"}, {"-10.00", "-10"}, {"00012.030", "12.03"}, {"987654321098765432.0", "987654321098765432"}} {
		a, ea := ParseDecimal([]byte(pair[0]))
		b, eb := ParseDecimal([]byte(pair[1]))
		if ea != nil || eb != nil || a != b {
			t.Fatal(pair, a, b, ea, eb)
		}
	}
	for _, value := range []string{"", "-", ".", "1.2.3", "1e2", " 10", "99999999999999999999"} {
		if _, err := ParseDecimal([]byte(value)); err == nil {
			t.Fatal("accepted invalid decimal", value)
		}
	}
}

func TestDecimalString(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{"447.00", "447"},
		{"1469.50", "1469.5"},
		{"0.50", "0.5"},
		{".5", "0.5"},
		{"-0.01", "-0.01"},
		{"0.00", "0"},
		{"-0", "0"},
		{"00012.030", "12.03"},
		{"100", "100"},
		{"-1234.5678", "-1234.5678"},
	}
	for _, test := range tests {
		d, err := ParseDecimal([]byte(test.text))
		if err != nil {
			t.Fatal(err)
		}
		if got := d.String(); got != test.want {
			t.Errorf("ParseDecimal(%q).String() = %q, want %q", test.text, got, test.want)
		}
	}
	if got := (Decimal{Coefficient: 1000, Scale: 1}).String(); got != "100.0" {
		t.Errorf("Decimal{1000, 1}.String() = %q, want %q (an unreduced value prints what it holds)", got, "100.0")
	}
}

func TestDecimalSub(t *testing.T) {
	tests := []struct {
		a, b string
		want string
	}{
		{"80", "79.99", "0.01"},
		{"79.99", "80", "-0.01"},
		{"80", "80", "0"},
		{"10.5", "0.5", "10"},
		{"100", "0.25", "99.75"},
		{"-1.5", "2", "-3.5"},
		{"0.001", "1", "-0.999"},
	}
	for _, test := range tests {
		a, err := ParseDecimal([]byte(test.a))
		if err != nil {
			t.Fatal(err)
		}
		b, err := ParseDecimal([]byte(test.b))
		if err != nil {
			t.Fatal(err)
		}
		got, err := a.Sub(b)
		if err != nil {
			t.Fatalf("%s - %s: %v", test.a, test.b, err)
		}
		want, err := ParseDecimal([]byte(test.want))
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s - %s = %+v (%s), want %+v (%s)", test.a, test.b, got, got, want, test.want)
		}
	}

	big := Decimal{Coefficient: math.MaxInt64}
	if _, err := big.Sub(Decimal{Coefficient: 1, Scale: 1}); err == nil {
		t.Error("Sub: scaling MaxInt64 by 10 succeeded, want an error")
	}
	if _, err := big.Sub(Decimal{Coefficient: -1}); err == nil {
		t.Error("Sub: MaxInt64 - (-1) succeeded, want an error")
	}
}

func TestDecimalAdd(t *testing.T) {
	tests := []struct {
		a, b string
		want string
	}{
		{"79.99", "0.01", "80"},
		{"0", "40.00", "40"},
		{"392.04", "90.48", "482.52"},
		{"10.5", "-0.5", "10"},
		{"-1.5", "2", "0.5"},
		{"0.001", "1", "1.001"},
		{"-40", "-2.25", "-42.25"},
		{"0", "0", "0"},
	}
	for _, test := range tests {
		a, err := ParseDecimal([]byte(test.a))
		if err != nil {
			t.Fatal(err)
		}
		b, err := ParseDecimal([]byte(test.b))
		if err != nil {
			t.Fatal(err)
		}
		got, err := a.Add(b)
		if err != nil {
			t.Fatalf("%s + %s: %v", test.a, test.b, err)
		}
		want, err := ParseDecimal([]byte(test.want))
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s + %s = %+v (%s), want %+v (%s)", test.a, test.b, got, got, want, test.want)
		}
	}

	big := Decimal{Coefficient: math.MaxInt64}
	if _, err := big.Add(Decimal{Coefficient: 1, Scale: 1}); err == nil {
		t.Error("Add: scaling MaxInt64 by 10 succeeded, want an error")
	}
	if _, err := big.Add(Decimal{Coefficient: 1}); err == nil {
		t.Error("Add: MaxInt64 + 1 succeeded, want an error")
	}
}
