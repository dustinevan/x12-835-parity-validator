package service

import (
	"encoding/json"
	"testing"
)

func TestValidTag(t *testing.T) {
	for _, tag := range []string{"CLP04", "SVC03", "MOA", "LQ", "CAS*CO*253", "NM1*TT", "DTM*050", "AMT*B6", "MIA04", "CAS*CO*45*QTY", "RAW*K3", "RAW*NM1"} {
		if !ValidTag(tag) {
			t.Errorf("ValidTag(%q) = false, want true", tag)
		}
	}
	for _, tag := range []string{"", "clp04", "CLP04,", "CAS*CO*253*1*2", "line order", "RAW*", ";drop", "C"} {
		if ValidTag(tag) {
			t.Errorf("ValidTag(%q) = true, want false", tag)
		}
	}
}

func TestParseMatchKey(t *testing.T) {
	good := "0229c4c77346a4a2a57935eab2663262943511bf1bb94cbf30aec5eae5837408"
	if k, err := ParseMatchKey(good); err != nil || len(k) != 32 {
		t.Fatalf("ParseMatchKey(good) = %x, %v", k, err)
	}
	for _, bad := range []string{"", "abc", good[:62], "0229C4C77346A4A2A57935EAB2663262943511BF1BB94CBF30AEC5EAE5837408", good[:63] + "g"} {
		if _, err := ParseMatchKey(bad); err == nil {
			t.Errorf("ParseMatchKey(%q) succeeded, want an error", bad)
		}
	}
}

func TestTagOf(t *testing.T) {
	tests := []struct {
		tag  string
		want string
	}{
		{"CAS*CO*253", `{"tag":"CAS*CO*253","label":"Sequestration adjustment"}`},
		{"CAS*OA*23", `{"tag":"CAS*OA*23","label":"CAS*OA*23"}`},
		{"CAS*CO*253*QTY", `{"tag":"CAS*CO*253*QTY","label":"Sequestration adjustment, quantity"}`},
		{"RAW*K3", `{"tag":"RAW*K3","label":"Unplaced segment K3"}`},
	}
	for _, test := range tests {
		got, err := json.Marshal(TagOf(test.tag))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != test.want {
			t.Errorf("TagOf(%q) = %s, want %s", test.tag, got, test.want)
		}
	}
}
