package service

// Benchmarks for the in-memory comparison work. Nothing here touches the
// database, the store, or the filesystem after setup: the generated test set
// (make gen, under data/generated) is read into memory before the timer
// starts, so file I/O is not timed. mac/ is the MAC submission;
// vendor-1/<date> and vendor-2/<date> are the vendor submissions.
//
// They do not run under a plain go test. Run them with:
//
//	go test -run '^$' -bench . -benchmem -benchtime 3x ./backend/lib/core/service/

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/data/claim"
)

const benchGenerated = "../../../../data/generated"

// benchFile is one input file held in memory.
type benchFile struct {
	name string
	data []byte
}

// benchSubmission is one submission's files in file-name order. The prefix
// is the submission's directory under data/generated.
type benchSubmission struct {
	prefix string
	files  []benchFile
}

func readBenchSubmission(b *testing.B, prefix string) benchSubmission {
	b.Helper()
	dir := filepath.Join(benchGenerated, prefix)
	entries, err := os.ReadDir(dir)
	if err != nil {
		b.Skipf("service: reading %s: %v (run make gen)", dir, err)
	}
	s := benchSubmission{prefix: prefix}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			b.Fatalf("service: reading %s: %v", entry.Name(), err)
		}
		s.files = append(s.files, benchFile{name: entry.Name(), data: data})
	}
	return s
}

// readBenchCorpus returns the MAC submission and the vendor submissions:
// vendor-1's dates, then vendor-2's.
func readBenchCorpus(b *testing.B) (benchSubmission, []benchSubmission) {
	b.Helper()
	mac := readBenchSubmission(b, "mac")
	var vendors []benchSubmission
	for _, side := range []string{"vendor-1", "vendor-2"} {
		dates, err := os.ReadDir(filepath.Join(benchGenerated, side))
		if err != nil {
			b.Skipf("service: reading %s: %v (run make gen)", side, err)
		}
		for _, date := range dates {
			if date.IsDir() {
				vendors = append(vendors, readBenchSubmission(b, side+"/"+date.Name()))
			}
		}
	}
	if len(vendors) == 0 {
		b.Skipf("service: no vendor submissions under %s (run make gen)", benchGenerated)
	}
	return mac, vendors
}

// splitBenchSubmission is the first half of Controller.Load: Split every
// file of the submission into claim rows. The database insert that follows
// in Load is not run.
func splitBenchSubmission(s benchSubmission) ([]claim.Claim, error) {
	var claims []claim.Claim
	for _, f := range s.files {
		split, _, err := Split(s.prefix, f.name, f.data)
		if err != nil {
			return nil, err
		}
		claims = append(claims, split...)
	}
	return claims, nil
}

// memory is the process's allocation counters at one moment.
type memory struct {
	mallocs    uint64 // objects allocated so far
	totalAlloc uint64 // bytes allocated so far
}

func readMemory() memory {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return memory{mallocs: m.Mallocs, totalAlloc: m.TotalAlloc}
}

// reportPer adds time, allocation count, and bytes allocated per unit (a
// claim, a pair) for the work between two readings. Call it after StopTimer.
func reportPer(b *testing.B, unit string, units int, start, end memory) {
	n := float64(units)
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/n, "ns/"+unit)
	b.ReportMetric(float64(end.mallocs-start.mallocs)/n, "allocs/"+unit)
	b.ReportMetric(float64(end.totalAlloc-start.totalAlloc)/n, "B/"+unit)
}

// BenchmarkSplit runs the engine's file reader, Split, over every generated
// file: the typed parse of the whole file, the claim spans, the match key,
// and the copy of each claim's bytes.
func BenchmarkSplit(b *testing.B) {
	mac, vendors := readBenchCorpus(b)
	submissions := append([]benchSubmission{mac}, vendors...)
	totalBytes, files, claims := 0, 0, 0
	for _, s := range submissions {
		rows, err := splitBenchSubmission(s)
		if err != nil {
			b.Fatal(err)
		}
		claims += len(rows)
		files += len(s.files)
		for _, f := range s.files {
			totalBytes += len(f.data)
		}
	}

	b.ReportAllocs()
	b.SetBytes(int64(totalBytes))
	start := readMemory()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range submissions {
			if _, err := splitBenchSubmission(s); err != nil {
				b.Fatal(err)
			}
		}
	}
	b.StopTimer()
	end := readMemory()
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*files), "ns/file")
	reportPer(b, "claim", b.N*claims, start, end)
	b.ReportMetric(float64(files), "files/op")
	b.ReportMetric(float64(claims), "claims/op")
}

// benchPair returns the pair of median MAC claim size among the pairs with
// differences (or, when wantDifferences is false, without). A pair is a MAC
// claim whose match key is on one claim on each side. Median size, because
// the first pair in file order may be much larger or smaller than most.
func benchPair(b *testing.B, macClaims, vendorClaims []claim.Claim, wantDifferences bool) (claim.Claim, claim.Claim) {
	b.Helper()
	macCount := map[string]int{}
	for _, c := range macClaims {
		macCount[string(c.MatchKey)]++
	}
	vendorByKey := map[string][]claim.Claim{}
	for _, c := range vendorClaims {
		vendorByKey[string(c.MatchKey)] = append(vendorByKey[string(c.MatchKey)], c)
	}
	var macs, vendors []claim.Claim
	for _, m := range macClaims {
		k := string(m.MatchKey)
		if macCount[k] != 1 || len(vendorByKey[k]) != 1 {
			continue
		}
		diffs, err := DiffClaims(compareSet1, m, vendorByKey[k][0])
		if err != nil {
			b.Fatal(err)
		}
		if (len(diffs) > 0) == wantDifferences {
			macs = append(macs, m)
			vendors = append(vendors, vendorByKey[k][0])
		}
	}
	if len(macs) == 0 {
		b.Fatalf("service: no pair with differences=%v in the first vendor submission", wantDifferences)
	}
	order := make([]int, len(macs))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool { return len(macs[order[i]].Segments) < len(macs[order[j]].Segments) })
	median := order[len(order)/2]
	return macs[median], vendors[median]
}

// BenchmarkComparePair times ComparePair, the call the compare stage makes
// for each pair its per-workload stream returns, on one pair from the first
// vendor submission: the median-size pair that has differences, and the
// median-size pair that has none. Each call parses both claims' bytes into
// the typed model and compares them under set 1.
func BenchmarkComparePair(b *testing.B) {
	mac, vendors := readBenchCorpus(b)
	macClaims, err := splitBenchSubmission(mac)
	if err != nil {
		b.Fatal(err)
	}
	vendorClaims, err := splitBenchSubmission(vendors[0])
	if err != nil {
		b.Fatal(err)
	}
	tests := []struct {
		name            string
		wantDifferences bool
	}{
		{"different", true},
		{"equal", false},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			m, c := benchPair(b, macClaims, vendorClaims, test.wantDifferences)
			row, _, err := ComparePair(compareSet1, vendors[0].prefix, m, c)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(m.Segments) + len(c.Segments)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := ComparePair(compareSet1, vendors[0].prefix, m, c); err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(row.DifferenceCount), "differences/op")
		})
	}
}
