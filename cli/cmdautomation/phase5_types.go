package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// HelpViolation records an undocumented CLI command or missing description.
type HelpViolation struct {
	File       string `json:"file"`
	Command    string `json:"command"`
	Issue      string `json:"issue"`
	LineNumber int    `json:"lineNumber"`
}

// HelpAuditOptions configures CLI command discovery and help parity checks.
type HelpAuditOptions struct {
	Dir        string   `json:"dir"`
	Extensions []string `json:"extensions"`
	IsStrict   bool     `json:"isStrict"`
	IsJson     bool     `json:"isJson"`
}

// HelpAuditResult aggregates CLI help parity findings.
type HelpAuditResult struct {
	ScannedFiles  int             `json:"scannedFiles"`
	TotalCommands int             `json:"totalCommands"`
	Violations    []HelpViolation `json:"violations"`
	Duration      time.Duration   `json:"duration"`
	IsClean       bool            `json:"isClean"`
}

// PlanCluster represents a grouped cluster of completed plans.
type PlanCluster struct {
	Id           string   `json:"id"`
	Title        string   `json:"title"`
	Domain       string   `json:"domain"`
	MergedPlans  []string `json:"mergedPlans"`
	SubtaskCount int      `json:"subtaskCount"`
}

// PlanConsolidateOptions configures plan clustering and memory consolidation.
type PlanConsolidateOptions struct {
	Dir       string `json:"dir"`
	Threshold int    `json:"threshold"`
	IsDryRun  bool   `json:"isDryRun"`
	IsForce   bool   `json:"isForce"`
	IsJson    bool   `json:"isJson"`
}

// PlanConsolidateResult stores summary of plan consolidation.
type PlanConsolidateResult struct {
	TotalPlans      int           `json:"totalPlans"`
	CompletedPlans  int           `json:"completedPlans"`
	ClustersCount   int           `json:"clustersCount"`
	SubtasksCleaned int           `json:"subtasksCleaned"`
	Clusters        []PlanCluster `json:"clusters"`
	Duration        time.Duration `json:"duration"`
	IsSuccess       bool          `json:"isSuccess"`
}

// DocLinkViolation records a broken or outdated relative markdown link.
type DocLinkViolation struct {
	File           string `json:"file"`
	LineNumber     int    `json:"lineNumber"`
	RawLink        string `json:"rawLink"`
	ResolvedTarget string `json:"resolvedTarget"`
	Issue          string `json:"issue"`
}

// DocLinksOptions configures markdown link integrity checks and autofixing.
type DocLinksOptions struct {
	Dir    string `json:"dir"`
	IsFix  bool   `json:"isFix"`
	IsJson bool   `json:"isJson"`
}

// DocLinksResult aggregates outcomes of markdown link integrity scans.
type DocLinksResult struct {
	ScannedFiles int                `json:"scannedFiles"`
	TotalLinks   int                `json:"totalLinks"`
	BrokenLinks  int                `json:"brokenLinks"`
	FixedLinks   int                `json:"fixedLinks"`
	Violations   []DocLinkViolation `json:"violations"`
	Duration     time.Duration      `json:"duration"`
	IsClean      bool               `json:"isClean"`
}

// SpecMigrateChange records an updated spec file and cross-link count.
type SpecMigrateChange struct {
	SourceFile string `json:"sourceFile"`
	TargetFile string `json:"targetFile"`
	References int    `json:"references"`
}

// SpecMigrateOptions configures spec re-sequencing and reference updating.
type SpecMigrateOptions struct {
	Dir      string `json:"dir"`
	FromNum  int    `json:"fromNum"`
	ToNum    int    `json:"toNum"`
	IsDryRun bool   `json:"isDryRun"`
	IsJson   bool   `json:"isJson"`
}

// SpecMigrateResult stores the outcome of a spec migration execution.
type SpecMigrateResult struct {
	TotalSpecs   int                 `json:"totalSpecs"`
	Migrated     []SpecMigrateChange `json:"migrated"`
	FilesUpdated int                 `json:"filesUpdated"`
	Duration     time.Duration       `json:"duration"`
	IsSuccess    bool                `json:"isSuccess"`
}

type (
	// HelpAuditResultMonad wraps HelpAuditResult with AppError.
	HelpAuditResultMonad = result.Result[HelpAuditResult]

	// PlanConsolidateResultMonad wraps PlanConsolidateResult with AppError.
	PlanConsolidateResultMonad = result.Result[PlanConsolidateResult]

	// DocLinksResultMonad wraps DocLinksResult with AppError.
	DocLinksResultMonad = result.Result[DocLinksResult]

	// SpecMigrateResultMonad wraps SpecMigrateResult with AppError.
	SpecMigrateResultMonad = result.Result[SpecMigrateResult]
)
