package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// SearchMatch represents a single pattern match in a file.
type SearchMatch struct {
	Path        string `json:"path"`
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
}

// SearchOptions configures the multi-core search engine.
type SearchOptions struct {
	Pattern           string   `json:"pattern"`
	Dir               string   `json:"dir"`
	IsRegex           bool     `json:"isRegex"`
	IsCaseInsensitive bool     `json:"isCaseInsensitive"`
	Extensions        []string `json:"extensions"`
	Workers           int      `json:"workers"`
}

// SearchResult wraps the collected matches and execution metrics.
type SearchResult struct {
	Matches     []SearchMatch `json:"matches"`
	TotalFiles  int           `json:"totalFiles"`
	Duration    time.Duration `json:"duration"`
	TotalHits   int           `json:"totalHits"`
}

// NewlineOptions configures the polyglot newline normalizer.
type NewlineOptions struct {
	Paths      []string `json:"paths"`
	IsFixMode  bool     `json:"isFixMode"`
	IsDryRun   bool     `json:"isDryRun"`
	Extensions []string `json:"extensions"`
}

// NewlineResult holds statistics from newline normalization.
type NewlineResult struct {
	ScannedFiles  int           `json:"scannedFiles"`
	ModifiedFiles int           `json:"modifiedFiles"`
	CrlfCount     int           `json:"crlfCount"`
	Duration      time.Duration `json:"duration"`
}

// BenchmarkMetric captures timing and resource metrics for an operation.
type BenchmarkMetric struct {
	Name           string        `json:"name"`
	GoDuration     time.Duration `json:"goDuration"`
	PyDuration     time.Duration `json:"pyDuration"`
	GoSpeedup      float64       `json:"goSpeedup"`
	GoTempBytes    int64         `json:"goTempBytes"`
	PyTempBytes    int64         `json:"pyTempBytes"`
	SampleFiles    int           `json:"sampleFiles"`
}

type (
	// SearchResultMonad wraps SearchResult with an AppError.
	SearchResultMonad = result.Result[SearchResult]

	// NewlineResultMonad wraps NewlineResult with an AppError.
	NewlineResultMonad = result.Result[NewlineResult]

	// BenchmarkMetricMonad wraps BenchmarkMetric with an AppError.
	BenchmarkMetricMonad = result.Result[BenchmarkMetric]
)
