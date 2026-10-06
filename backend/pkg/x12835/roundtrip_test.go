package x12835

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRoundTrip parses every 835 under testdata/taxonomy and every example
// the docs site uses and requires Write to give back the input bytes
// exactly: every segment, every element, in the file's own order and
// spelling. The one allowance is a file whose sender transmitted trailing
// empty elements ("CUR**~"), which X12 says are not transmitted and the
// writer leaves off; such a file must match once those are removed from the
// input, and nothing else may differ.
func TestRoundTrip(t *testing.T) {
	examples, err := filepath.Glob(filepath.Join(repoRoot, "frontend/src/data/viewer835/files/*.835"))
	if err != nil {
		t.Fatal(err)
	}
	taxonomy, err := filepath.Glob(filepath.Join(repoRoot, "testdata/taxonomy/*/*.835"))
	if err != nil {
		t.Fatal(err)
	}
	paths := append(examples, taxonomy...)
	if len(paths) == 0 {
		t.Fatal("no 835 files found")
	}
	exact, trailing := 0, 0
	for _, path := range paths {
		t.Run(strings.TrimPrefix(path, repoRoot+"/"), func(t *testing.T) {
			input, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			isa, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse = %v", err)
			}
			var out bytes.Buffer
			if err := Write(&out, isa); err != nil {
				t.Fatalf("Write = %v", err)
			}
			if bytes.Equal(out.Bytes(), input) {
				exact++
				return
			}
			stripped := withoutTrailingEmpties(input, isa.Delimiters)
			if bytes.Equal(out.Bytes(), stripped) {
				trailing++
				t.Logf("round-trips once the sender's trailing empty elements are left off")
				return
			}
			t.Errorf("Write(Parse(file)) differs from the file at byte %d:\n got %q\nwant %q", firstDifference(out.Bytes(), input), excerpt(out.Bytes(), input), excerpt(input, out.Bytes()))
		})
	}
	t.Logf("%d files byte for byte, %d once trailing empty elements are left off", exact, trailing)
}

// withoutTrailingEmpties removes element separators that sit right before
// a segment terminator, the trailing empty elements X12 does not transmit.
func withoutTrailingEmpties(data []byte, d Delimiters) []byte {
	segments := bytes.Split(data, []byte{d.Terminator})
	for i, s := range segments {
		segments[i] = bytes.TrimRight(s, string(d.Element))
	}
	return bytes.Join(segments, []byte{d.Terminator})
}

func firstDifference(a, b []byte) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

// excerpt is the 60 bytes of a around the first difference with b.
func excerpt(a, b []byte) string {
	i := firstDifference(a, b)
	start := max(i-20, 0)
	end := min(i+40, len(a))
	return string(a[start:end])
}
