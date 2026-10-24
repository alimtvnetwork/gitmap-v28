package cmdplans

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"

// PlanStatusType aliases the strongly-typed enum from the enums package.
type PlanStatusType = enums.PlanStatusType

// PlansFilterModeType aliases the filter mode enum.
type PlansFilterModeType = enums.PlansFilterModeType

// PlansOutputFormatType aliases the output format enum.
type PlansOutputFormatType = enums.PlansOutputFormatType

const (
	// StatusCompleted indicates a plan whose tasks have been verified and finished.
	StatusCompleted = enums.PlanStatusCompleted
	// StatusPending indicates a plan queued for future execution.
	StatusPending = enums.PlanStatusPending
	// StatusInProgress indicates a plan currently under active execution.
	StatusInProgress = enums.PlanStatusInProgress
	// StatusActive indicates a plan active in execution.
	StatusActive = enums.PlanStatusActive
	// StatusArchived indicates an archived plan.
	StatusArchived = enums.PlanStatusArchived
	// StatusFailed indicates a failed plan.
	StatusFailed = enums.PlanStatusFailed
)

// PlanItem represents a single discovered plan metadata record.
type PlanItem struct {
	ID        int            `json:"id" yaml:"id"`
	Slug      string         `json:"slug" yaml:"slug"`
	Title     string         `json:"title" yaml:"title"`
	Status    PlanStatusType `json:"status" yaml:"status"`
	FilePath  string         `json:"filePath" yaml:"filePath"`
	UpdatedAt string         `json:"updatedAt,omitempty" yaml:"updatedAt,omitempty"`
}

// PlansSummary contains aggregate metrics and overflow flags for plan inventory.
type PlansSummary struct {
	CompletedCount   int    `json:"completedCount" yaml:"completedCount"`
	PendingCount     int    `json:"pendingCount" yaml:"pendingCount"`
	TotalCount       int    `json:"totalCount" yaml:"totalCount"`
	DisplayedLimit   int    `json:"displayedLimit" yaml:"displayedLimit"`
	HasMorePending   bool   `json:"hasMorePending" yaml:"hasMorePending"`
	HasMoreCompleted bool   `json:"hasMoreCompleted" yaml:"hasMoreCompleted"`
	PendingOverflow  string `json:"pendingOverflow,omitempty" yaml:"pendingOverflow,omitempty"`
}

// PlansReport represents the complete serializable payload for JSON and YAML output.
type PlansReport struct {
	Summary   PlansSummary `json:"summary" yaml:"summary"`
	Completed []PlanItem   `json:"completed" yaml:"completed"`
	Pending   []PlanItem   `json:"pending" yaml:"pending"`
}

// PlansOptions configures the execution parameters for the plans command.
type PlansOptions struct {
	Limit           int
	Format          enums.PlansOutputFormatType
	FilterMode      enums.PlansFilterModeType
	IsJSON          bool
	IsYAML          bool
	IsAll           bool
	IsPendingOnly   bool
	IsCompletedOnly bool
	RepoPath        string
}
