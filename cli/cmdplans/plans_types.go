package cmdplans

// PlanStatusType represents the lifecycle state of an architectural plan.
type PlanStatusType string

const (
	// StatusCompleted indicates a plan whose tasks have been verified and finished.
	StatusCompleted PlanStatusType = "Completed"
	// StatusPending indicates a plan queued for future execution.
	StatusPending PlanStatusType = "Pending"
	// StatusInProgress indicates a plan currently under active execution.
	StatusInProgress PlanStatusType = "In Progress"
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
	IsJSON          bool
	IsYAML          bool
	IsAll           bool
	IsPendingOnly   bool
	IsCompletedOnly bool
	RepoPath        string
}
