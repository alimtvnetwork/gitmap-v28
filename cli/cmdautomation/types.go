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
	MaxJsonKb         int      `json:"maxJsonKb"`
	IncludeBinaries   bool     `json:"includeBinaries"`
	IncludeLargeJson  bool     `json:"includeLargeJson"`
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
	Name        string        `json:"name"`
	GoDuration  time.Duration `json:"goDuration"`
	PyDuration  time.Duration `json:"pyDuration"`
	GoSpeedup   float64       `json:"goSpeedup"`
	GoTempBytes int64         `json:"goTempBytes"`
	PyTempBytes int64         `json:"pyTempBytes"`
	SampleFiles int           `json:"sampleFiles"`
}

// OversizedFile records an audited file that exceeds limits or is a binary.
type OversizedFile struct {
	Path            string `json:"path"`
	SizeBytes       int64  `json:"sizeBytes"`
	MaxAllowedBytes int64  `json:"maxAllowedBytes"`
	Kind            string `json:"kind"`
}

// GuardOptions configures repository blob and binary audit limits.
type GuardOptions struct {
	Dir         string `json:"dir"`
	MaxFileKb   int    `json:"maxFileKb"`
	MaxJsonKb   int    `json:"maxJsonKb"`
	Interactive bool   `json:"interactive"`
	AutoExclude bool   `json:"autoExclude"`
	AsJson      bool   `json:"asJson"`
}

// GuardResult captures discovered oversized and binary files.
type GuardResult struct {
	TotalFiles         int             `json:"totalFiles"`
	OversizedFiles     []OversizedFile `json:"oversizedFiles"`
	BinaryFiles        []OversizedFile `json:"binaryFiles"`
	ExcludedLargeJsons []OversizedFile `json:"excludedLargeJsons"`
	ExcludedCount      int             `json:"excludedCount"`
	Duration           time.Duration   `json:"duration"`
}

// SequenceAuditorOptions configures the numbering and title audit.
type SequenceAuditorOptions struct {
	Dir       string `json:"dir"`
	IsFixMode bool   `json:"isFixMode"`
	AsJson    bool   `json:"asJson"`
}

// SequenceDirResult captures sequence and title health for a folder.
type SequenceDirResult struct {
	Dir                 string   `json:"dir"`
	NumberedFilesCount  int      `json:"numberedFilesCount"`
	SequenceGaps        []string `json:"sequenceGaps"`
	TitleMismatches     []string `json:"titleMismatches"`
	FixedTitles         []string `json:"fixedTitles"`
	IsClean             bool     `json:"isClean"`
}

// SequenceAuditorResult aggregates findings across all audited folders.
type SequenceAuditorResult struct {
	ScannedDirs        int                 `json:"scannedDirs"`
	TotalNumberedFiles int                 `json:"totalNumberedFiles"`
	DirResults         []SequenceDirResult `json:"dirResults"`
	TotalGaps          int                 `json:"totalGaps"`
	TotalMismatches    int                 `json:"totalMismatches"`
	TotalFixed         int                 `json:"totalFixed"`
	Duration           time.Duration       `json:"duration"`
}

type (
	// SearchResultMonad wraps SearchResult with an AppError.
	SearchResultMonad = result.Result[SearchResult]

	// NewlineResultMonad wraps NewlineResult with an AppError.
	NewlineResultMonad = result.Result[NewlineResult]

	// BenchmarkMetricMonad wraps BenchmarkMetric with an AppError.
	BenchmarkMetricMonad = result.Result[BenchmarkMetric]

	// GuardResultMonad wraps GuardResult with an AppError.
	GuardResultMonad = result.Result[GuardResult]

	// SequenceAuditorResultMonad wraps SequenceAuditorResult with an AppError.
	SequenceAuditorResultMonad = result.Result[SequenceAuditorResult]
)
