package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// VersionCheckItem records the audit status of a single version file.
type VersionCheckItem struct {
	Name            string `json:"name"`
	TargetFile      string `json:"targetFile"`
	CurrentVersion  string `json:"currentVersion"`
	ExpectedVersion string `json:"expectedVersion"`
	IsMatch         bool   `json:"isMatch"`
	Message         string `json:"message"`
}

// VersionSyncOptions configures version synchronization and auditing.
type VersionSyncOptions struct {
	Dir           string `json:"dir"`
	IsFixMode     bool   `json:"isFixMode"`
	TargetVersion string `json:"targetVersion"`
	IsJson        bool   `json:"isJson"`
}

// VersionSyncResult aggregates the outcomes of all checked version files.
type VersionSyncResult struct {
	CanonicalVersion string             `json:"canonicalVersion"`
	TotalChecked     int                `json:"totalChecked"`
	MismatchCount    int                `json:"mismatchCount"`
	FixedCount       int                `json:"fixedCount"`
	Items            []VersionCheckItem `json:"items"`
	Duration         time.Duration      `json:"duration"`
	IsClean          bool               `json:"isClean"`
}

// ReleaseBumpOptions configures SemVer version bumping and manifest updates.
type ReleaseBumpOptions struct {
	Dir           string   `json:"dir"`
	BumpType      string   `json:"bumpType"`
	TargetVersion string   `json:"targetVersion"`
	Scope         string   `json:"scope"`
	Bullets       []string `json:"bullets"`
	IsDryRun      bool     `json:"isDryRun"`
	IsTag         bool     `json:"isTag"`
	IsPush        bool     `json:"isPush"`
	IsSkipTests   bool     `json:"isSkipTests"`
	IsJson        bool     `json:"isJson"`
}

// ReleaseBumpResult stores the result of a release bump execution.
type ReleaseBumpResult struct {
	PreviousVersion string        `json:"previousVersion"`
	NewVersion      string        `json:"newVersion"`
	BumpType        string        `json:"bumpType"`
	UpdatedFiles    []string      `json:"updatedFiles"`
	ReleaseBranch   string        `json:"releaseBranch"`
	TagName         string        `json:"tagName"`
	IsTagged        bool          `json:"isTagged"`
	IsPushed        bool          `json:"isPushed"`
	Duration        time.Duration `json:"duration"`
	IsSuccess       bool          `json:"isSuccess"`
}

// MilestoneIssue represents an issue or PR associated with a GitHub milestone.
type MilestoneIssue struct {
	Number        int    `json:"number"`
	Title         string `json:"title"`
	State         string `json:"state"`
	Category      string `json:"category"`
	Url           string `json:"url"`
	IsPullRequest bool   `json:"isPullRequest"`
	IsClosed      bool   `json:"isClosed"`
}

// MilestonesOptions configures milestone querying and release notes generation.
type MilestonesOptions struct {
	Dir         string `json:"dir"`
	MilestoneId string `json:"milestoneId"`
	OutputPath  string `json:"outputPath"`
	IsJson      bool   `json:"isJson"`
}

// MilestonesResult captures consolidated milestone issues and formatted notes.
type MilestonesResult struct {
	MilestoneId  string           `json:"milestoneId"`
	Title        string           `json:"title"`
	TotalItems   int              `json:"totalItems"`
	ClosedCount  int              `json:"closedCount"`
	OpenCount    int              `json:"openCount"`
	Items        []MilestoneIssue `json:"items"`
	ReleaseNotes string           `json:"releaseNotes"`
	Duration     time.Duration    `json:"duration"`
	IsSuccess    bool             `json:"isSuccess"`
}

type (
	// VersionSyncResultMonad wraps VersionSyncResult with AppError.
	VersionSyncResultMonad = result.Result[VersionSyncResult]

	// ReleaseBumpResultMonad wraps ReleaseBumpResult with AppError.
	ReleaseBumpResultMonad = result.Result[ReleaseBumpResult]

	// MilestonesResultMonad wraps MilestonesResult with AppError.
	MilestonesResultMonad = result.Result[MilestonesResult]
)
