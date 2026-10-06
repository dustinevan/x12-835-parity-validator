package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sync"

	"github.com/dustinevan/x12-835-parity-validator/backend/lib/sys"
)

var ErrBusy = errors.New("service: a diff is already running")
var ErrInput = errors.New("service: invalid input")

// Service compares two named input files. Inputs come from the filesystem or,
// when enabled, read-only S3; the difference file is always written to the
// local filesystem. Set is the compared set the comparison runs under.
type Service struct {
	Inputs          sys.Reader
	Outputs         sys.FileSystem
	OutputDirectory string
	Set             CompareSet
	mu              sync.Mutex
}

type Result struct {
	File        string `json:"file"`
	Differences int    `json:"differences"`
}

// diffFile is the JSON file one diff writes: the two inputs and their
// SHA-256s, the compared set, how many claims each file holds and how many
// of them had no pair, and every difference of every compared pair. Claims
// pair by match key within the two files, as submissions do
// (CompareSubmissions): a key on one side only, or on more than one claim
// on either side, is unpaired and compared with nothing.
type diffFile struct {
	Original          string           `json:"original"`
	Candidate         string           `json:"candidate"`
	OriginalSHA256    string           `json:"original_sha256"`
	CandidateSHA256   string           `json:"candidate_sha256"`
	ComparedSet       string           `json:"compared_set"`
	OriginalClaims    int              `json:"original_claims"`
	CandidateClaims   int              `json:"candidate_claims"`
	UnpairedOriginal  int              `json:"unpaired_original"`
	UnpairedCandidate int              `json:"unpaired_candidate"`
	Differences       []diffDifference `json:"differences"`
}

// diffDifference is one difference of one pair, as the difference table
// holds it, with the pair's match key in hex. An empty line control number
// means a claim-level element; an empty value means the element is absent
// on that side (an absent dollar amount counts as zero in the amount
// difference); an empty amount difference means the element is not a
// dollar amount.
type diffDifference struct {
	MatchKey          string `json:"match_key"`
	Element           string `json:"element"`
	Occurrence        int    `json:"occurrence"`
	LineControlNumber string `json:"line_control_number"`
	OriginalValue     string `json:"original_value"`
	CandidateValue    string `json:"candidate_value"`
	AmountDifference  string `json:"amount_difference"`
}

// Diff compares the claims of two files and records every value difference
// under the compared set, amounts under the five-cent rule
// (AmountMatchCents) and everything else exactly. Explanations, any other
// tolerance, and scoring filters belong in later analysis of the saved
// evidence. It runs
// synchronously. Concurrent HTTP requests are rejected before reading
// either input, so only one pair of files is held by this service.
func (s *Service) Diff(ctx context.Context, originalName, candidateName string) (Result, error) {
	if !filepath.IsLocal(originalName) || !filepath.IsLocal(candidateName) {
		return Result{}, fmt.Errorf("%w: filenames must be relative to the store", ErrInput)
	}
	if !s.mu.TryLock() {
		return Result{}, ErrBusy
	}
	defer s.mu.Unlock()

	originalBytes, err := s.Inputs.Read(ctx, originalName)
	if err != nil {
		return Result{}, fmt.Errorf("service: reading original: %w", err)
	}
	candidateBytes, err := s.Inputs.Read(ctx, candidateName)
	if err != nil {
		return Result{}, fmt.Errorf("service: reading candidate: %w", err)
	}
	original, _, err := Split("original", originalName, originalBytes)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %w", ErrInput, err)
	}
	candidate, _, err := Split("candidate", candidateName, candidateBytes)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %w", ErrInput, err)
	}
	_, differences, summary, err := CompareSubmissions(s.Set, "candidate", original, candidate)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %w", ErrInput, err)
	}
	pairs := summary.Equal + summary.Different

	// Input hashes make each artifact traceable to the exact bytes compared.
	file := diffFile{
		Original:          originalName,
		Candidate:         candidateName,
		OriginalSHA256:    fmt.Sprintf("%x", sha256.Sum256(originalBytes)),
		CandidateSHA256:   fmt.Sprintf("%x", sha256.Sum256(candidateBytes)),
		ComparedSet:       s.Set.Name,
		OriginalClaims:    len(original),
		CandidateClaims:   len(candidate),
		UnpairedOriginal:  len(original) - pairs,
		UnpairedCandidate: len(candidate) - pairs,
		Differences:       make([]diffDifference, 0, len(differences)),
	}
	for _, d := range differences {
		file.Differences = append(file.Differences, diffDifference{
			MatchKey:          hex.EncodeToString(d.MatchKey),
			Element:           d.Element,
			Occurrence:        d.Occurrence,
			LineControlNumber: d.LineControlNumber,
			OriginalValue:     d.MACValue,
			CandidateValue:    d.VendorValue,
			AmountDifference:  d.AmountDifference,
		})
	}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return Result{}, fmt.Errorf("service: encoding differences: %w", err)
	}
	name := path.Join(s.OutputDirectory, fmt.Sprintf("%x.json", sha256.Sum256(data)))
	if err := s.Outputs.Write(ctx, name, data); err != nil {
		return Result{}, fmt.Errorf("service: writing differences: %w", err)
	}
	return Result{File: name, Differences: len(file.Differences)}, nil
}
