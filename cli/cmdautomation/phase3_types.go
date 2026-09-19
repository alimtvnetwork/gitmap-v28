package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PreflightOptions configures local CI/CD preflight multi-core checks.
type PreflightOptions struct {
	Dir        string `json:"dir"`
	IsFailFast bool   `json:"isFailFast"`
	IsJson     bool   `json:"isJson"`
	Workers    int    `json:"workers"`
	Phase      string `json:"phase"`
	Filter     string `json:"filter"`
}

// PreflightCheck represents the result of a single preflight check.
type PreflightCheck struct {
	Name     string        `json:"name"`
	Category string        `json:"category"`
	IsPass   bool          `json:"isPass"`
	IsFail   bool          `json:"isFail"`
	Duration time.Duration `json:"duration"`
	Message  string        `json:"message"`
	Details  []string      `json:"details,omitempty"`
}

// PreflightResult aggregates all preflight check outcomes.
type PreflightResult struct {
	TotalChecks  int              `json:"totalChecks"`
	PassedChecks int              `json:"passedChecks"`
	FailedChecks int              `json:"failedChecks"`
	Checks       []PreflightCheck `json:"checks"`
	Duration     time.Duration    `json:"duration"`
	IsPass       bool             `json:"isPass"`
}

// TestInventoryOptions configures test discovery and inventory generation.
type TestInventoryOptions struct {
	Dir           string   `json:"dir"`
	OutPath       string   `json:"outPath"`
	IsRefresh     bool     `json:"isRefresh"`
	IsJson        bool     `json:"isJson"`
	SlowThreshold float64  `json:"slowThreshold"`
	IsForceRunAll bool     `json:"isForceRunAll"`
	RecordPaths   []string `json:"recordPaths"`
	IsQueryRecent bool     `json:"isQueryRecent"`
	IsClear       bool     `json:"isClear"`
}

// TestInventoryItem records metadata for an individual test function.
type TestInventoryItem struct {
	Id          string  `json:"id"`
	Package     string  `json:"package"`
	TestFile    string  `json:"test_file"`
	TestFunc    string  `json:"test_func"`
	TestHash    string  `json:"test_hash"`
	TargetFile  string  `json:"target_file"`
	TargetFunc  string  `json:"target_func"`
	CodeHash    string  `json:"code_hash"`
	DurationSec float64 `json:"duration_sec"`
	Tier        string  `json:"tier"`
	IsSlow      bool    `json:"is_slow"`
	LastStatus  string  `json:"last_status"`
	NeedsRun    bool    `json:"needs_run"`
}

// TestInventorySummary summarizes test counts, durations, and packages.
type TestInventorySummary struct {
	Total            int     `json:"total"`
	Cached           int     `json:"cached"`
	Dirty            int     `json:"dirty"`
	Packages         int     `json:"packages"`
	SlowTests        int     `json:"slow_tests"`
	FastTests        int     `json:"fast_tests"`
	SlowThresholdSec float64 `json:"slow_threshold_sec"`
	EstimatedSlowSec float64 `json:"estimated_slow_sec"`
	EstimatedFastSec float64 `json:"estimated_fast_sec"`
}

// TestInventoryResult represents the serialized inventory output.
type TestInventoryResult struct {
	Version         int                          `json:"version"`
	UpdatedAt       string                       `json:"updated_at"`
	TotalTests      int                          `json:"total_tests"`
	Summary         TestInventorySummary         `json:"summary"`
	Tests           map[string]TestInventoryItem `json:"tests"`
	OutPath         string                       `json:"outPath,omitempty"`
	RecordedFiles   []string                     `json:"recorded_files,omitempty"`
	AssociatedTests []string                     `json:"associated_tests,omitempty"`
	Duration        time.Duration                `json:"duration"`
	IsPass          bool                         `json:"isPass"`
}

// PurgeActionsOptions configures GitHub Actions artifact cleanup.
type PurgeActionsOptions struct {
	Repo          string `json:"repo"`
	OlderThanDays int    `json:"olderThanDays"`
	IsDryRun      bool   `json:"isDryRun"`
	IsJson        bool   `json:"isJson"`
	Workers       int    `json:"workers"`
}

// PurgeArtifactItem records an artifact identified or deleted.
type PurgeArtifactItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_in_bytes"`
	CreatedAt string `json:"created_at,omitempty"`
	IsPurged  bool   `json:"isPurged"`
}

// PurgeActionsResult aggregates GitHub Actions purge outcomes.
type PurgeActionsResult struct {
	Repo            string              `json:"repo"`
	TotalArtifacts  int                 `json:"totalArtifacts"`
	PurgedArtifacts int                 `json:"purgedArtifacts"`
	FreedBytes      int64               `json:"freedBytes"`
	FreedMB         float64             `json:"freedMB"`
	Artifacts       []PurgeArtifactItem `json:"artifacts"`
	Duration        time.Duration       `json:"duration"`
	IsPass          bool                `json:"isPass"`
}

// SmokeTestOptions configures installer and tool validation.
type SmokeTestOptions struct {
	Dir     string   `json:"dir"`
	Tools   []string `json:"tools"`
	Workers int      `json:"workers"`
	IsJson  bool     `json:"isJson"`
	Filter  string   `json:"filter"`
}

// SmokeTestItem records an individual smoke check outcome.
type SmokeTestItem struct {
	Name     string        `json:"name"`
	Target   string        `json:"target"`
	Kind     string        `json:"kind"`
	IsPass   bool          `json:"isPass"`
	IsFail   bool          `json:"isFail"`
	Duration time.Duration `json:"duration"`
	Message  string        `json:"message"`
	Issues   []string      `json:"issues,omitempty"`
}

// SmokeTestResult aggregates smoke test findings.
type SmokeTestResult struct {
	TotalItems  int             `json:"totalItems"`
	PassedItems int             `json:"passedItems"`
	FailedItems int             `json:"failedItems"`
	Items       []SmokeTestItem `json:"items"`
	Duration    time.Duration   `json:"duration"`
	IsPass      bool            `json:"isPass"`
}

type (
	// PreflightResultMonad wraps PreflightResult with AppError.
	PreflightResultMonad = result.Result[PreflightResult]

	// TestInventoryResultMonad wraps TestInventoryResult with AppError.
	TestInventoryResultMonad = result.Result[TestInventoryResult]

	// PurgeActionsResultMonad wraps PurgeActionsResult with AppError.
	PurgeActionsResultMonad = result.Result[PurgeActionsResult]

	// SmokeTestResultMonad wraps SmokeTestResult with AppError.
	SmokeTestResultMonad = result.Result[SmokeTestResult]
)
