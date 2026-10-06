package x12835

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// claimEnvelope wraps one claim's segments, written with the standard
// delimiters, in the smallest envelope Parse accepts: ISA, GS, ST, and an LX
// holding the claim, with SE, GE, and IEA after it.
func claimEnvelope(claim string) string {
	count := strings.Count(claim, "~") + 3 // ST, LX, the claim's segments, SE
	return "ISA*00*          *00*          *ZZ*SYNMAC91101    *ZZ*SYNCLEARING01  *260304*0600*^*00501*100000001*0*T*:~" +
		"GS*HP*SYNMAC91101*SYNCLEARING01*20260304*0600*1*X*005010X221A1~ST*835*0001~LX*1~" +
		claim + "SE*" + itoa(count) + "*0001~GE*1*1~IEA*1*100000001~"
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

// TestParseClaim checks that ParseClaim reads a claim's segments exactly as
// Parse reads the same claim inside a file, apart from every Index, which
// counts from the CLP, and that WriteClaim writes them back as they came.
func TestParseClaim(t *testing.T) {
	claim := "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~NM1*QC*1*SYNTHETIC*ALICE*E***MI*1AA0AA0AA00~" +
		"NM1*QC*1*SYNTHETIC*ALICE*E***MI*2BB1BB1BB11~MOA***MA01*N130~K3*TEXT~DTM*232*20260105~DTM*472*20260105~" +
		"SVC*HC:99213:25*100.00*80.00**1~DTM*472*20260105~CAS*CO*45*0.00**253*0.01*2~REF*6R*L1~AMT*B6*100.00~LQ*HE*N19~QTY*ZL*1~" +
		"SVC*HC:85025*49.00*7.32*0300*1~CAS*PI*45*41.53~REF*6R*L2~"
	want := func() ClaimCLP {
		isa, err := Parse([]byte(claimEnvelope(claim)))
		if err != nil {
			t.Fatal(err)
		}
		return isa.FunctionalGroups[0].Transactions[0].HeaderNumbers[0].Claims[0]
	}()
	got, err := ParseClaim([]byte(claim), '*', ':', '^', '~')
	if err != nil {
		t.Fatal(err)
	}
	// The file's claim starts at segment 4 (ISA, GS, ST, LX before it).
	shiftIndexes(&got, 4)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseClaim:\n got %+v\nwant %+v", got, want)
	}
	if len(got.Raw) != 3 || got.Raw[0].ID != "NM1" || got.Raw[1].ID != "K3" || got.Raw[2].ID != "DTM" {
		t.Fatalf("Raw = %+v, want the second NM1*QC, the K3, and the claim-level DTM*472", got.Raw)
	}
	if got.ServiceLines[0].Adjustments[0].Reasons[1].Quantity != (Amount{Value: Decimal{Coefficient: 2}, Set: true, Text: "2"}) {
		t.Fatalf("CAS quantity = %+v, want 2", got.ServiceLines[0].Adjustments[0].Reasons[1].Quantity)
	}

	// The claim writes back out as it came, raw segments in place.
	var out strings.Builder
	if err := WriteClaim(&out, got, Delimiters{Element: '*', Component: ':', Repetition: '^', Terminator: '~'}); err != nil {
		t.Fatal(err)
	}
	if out.String() != claim {
		t.Fatalf("WriteClaim:\n got %q\nwant %q", out.String(), claim)
	}

	// Other delimiters read the same claim.
	other := strings.NewReplacer("*", "|", ":", ">", "^", "!", "~", "?").Replace(claim)
	again, err := ParseClaim([]byte(other), '|', '>', '!', '?')
	if err != nil {
		t.Fatal(err)
	}
	shiftIndexes(&again, 4)
	if !reflect.DeepEqual(again, want) {
		t.Fatalf("ParseClaim with other delimiters:\n got %+v\nwant %+v", again, want)
	}

	// Line breaks between segments and a newline terminator.
	broken := strings.ReplaceAll(claim, "~", "~\r\n")
	if _, err := ParseClaim([]byte(broken), '*', ':', '^', '~'); err != nil {
		t.Fatalf("ParseClaim with line breaks: %v", err)
	}
}

func TestParseClaimErrors(t *testing.T) {
	tests := []struct {
		name  string
		claim string
		want  string
	}{
		{"not a claim", "SVC*HC:99213*100.00*80.00**1~", "do not start with CLP"},
		{"empty", "", "do not start with CLP"},
		{"second claim", "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~CLP*A2*1*100.00*80.00*20.00*MB*ICN2*11*1~", "segment 1 (CLP) is outside the claim"},
		{"LX after the claim", "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~LX*2~", "segment 1 (LX) is outside the claim"},
		{"malformed amount", "CLP*A1*1*100.00*eighty*20.00*MB*ICN1*11*1~", "CLP04"},
		{"malformed date", "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1~DTM*232*2026010~", "DTM02"},
		{"missing terminator", "CLP*A1*1*100.00*80.00*20.00*MB*ICN1*11*1", "missing segment terminator"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseClaim([]byte(test.claim), '*', ':', '^', '~')
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ParseClaim = %v, want an error containing %q", err, test.want)
			}
		})
	}
	if _, err := ParseClaim([]byte("CLP*A1~"), '*', '*', '^', '~'); err == nil || !strings.Contains(err.Error(), "must differ") {
		t.Fatalf("ParseClaim with equal delimiters = %v, want an error", err)
	}
}

// shiftIndexes adds by to every Index in a claim, so a claim ParseClaim read
// (counting from its CLP) compares with the same claim Parse read inside a
// file.
func shiftIndexes(c *ClaimCLP, by int) {
	c.Index += by
	for i := range c.Adjustments {
		c.Adjustments[i].Index += by
	}
	for _, n := range []*NameNM1{&c.PatientNM1QC, &c.CorrectedPatientNM174, &c.RenderingNM182, &c.CrossoverNM1TT, &c.CorrectedPayerNM1PR} {
		if n.EntityIDNM101 != EntityIDNM101NotSent {
			n.Index += by
		}
	}
	if c.InpatientMIA.Set {
		c.InpatientMIA.Index += by
	}
	if c.OutpatientMOA.Set {
		c.OutpatientMOA.Index += by
	}
	for i := range c.References {
		c.References[i].Index += by
	}
	for _, d := range []*DateTimeDTM{&c.StatementFromDTM232, &c.StatementToDTM233, &c.ReceivedDTM050, &c.CoverageExpirationDTM036} {
		if d.DateQualifierDTM01 != DateQualifierDTM01NotSent {
			d.Index += by
		}
	}
	if c.ContactPER.FunctionPER01 != ContactFunctionPER01NotSent {
		c.ContactPER.Index += by
	}
	for i := range c.Amounts {
		c.Amounts[i].Index += by
	}
	for i := range c.Quantities {
		c.Quantities[i].Index += by
	}
	for i := range c.Raw {
		c.Raw[i].Index += by
	}
	for i := range c.ServiceLines {
		l := &c.ServiceLines[i]
		l.Index += by
		for _, d := range []*DateTimeDTM{&l.ServiceDateDTM472, &l.ServiceDateDTM150, &l.ServiceDateDTM151} {
			if d.DateQualifierDTM01 != DateQualifierDTM01NotSent {
				d.Index += by
			}
		}
		for j := range l.Adjustments {
			l.Adjustments[j].Index += by
		}
		for _, r := range []*ReferenceREF{&l.LineControlREF6R, &l.ServiceIDREFLU, &l.RenderingREFHPI, &l.PolicyREF0K} {
			if r.Value != "" {
				r.Index += by
			}
		}
		for j := range l.References {
			l.References[j].Index += by
		}
		for j := range l.Amounts {
			l.Amounts[j].Index += by
		}
		for j := range l.Quantities {
			l.Quantities[j].Index += by
		}
		for j := range l.Remarks {
			l.Remarks[j].Index += by
		}
		for j := range l.Raw {
			l.Raw[j].Index += by
		}
	}
}
