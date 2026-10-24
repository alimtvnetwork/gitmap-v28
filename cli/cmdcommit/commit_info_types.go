package cmdcommit

import "time"

// OriginType designates the origin classification of a commit.
type OriginType string

const (
	OriginLocalAI         OriginType = "LocalAI"
	OriginForeignAI       OriginType = "ForeignAI"
	OriginHumanOrExternal OriginType = "HumanOrExternal"
)

// CommitFileDelta captures changes to a single file in a commit.
type CommitFileDelta struct {
	FilePath   string `json:"filePath"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	IsBinary   bool   `json:"isBinary"`
}

// CommitProvenance holds origin and attribution metadata.
type CommitProvenance struct {
	OriginType      OriginType `json:"originType"`
	AttributionText string     `json:"attributionText"`
	AiName          string     `json:"aiName,omitempty"`
	AiModel         string     `json:"aiModel,omitempty"`
	ToolName        string     `json:"toolName,omitempty"`
	AgentRole       string     `json:"agentRole,omitempty"`
	SessionUuid     string     `json:"sessionUuid,omitempty"`
	TaskSlug        string     `json:"taskSlug,omitempty"`
	MachineHostname string     `json:"machineHostname,omitempty"`
	TransmissionId  string     `json:"transmissionId,omitempty"`
	AgentId         string     `json:"agentId,omitempty"`
	HeatmapId       string     `json:"heatmapId,omitempty"`
	IsLocalMachine  bool       `json:"isLocalMachine"`
	IsAiCommitted   bool       `json:"isAiCommitted"`
}

// CommitDbCrossReferences captures linkages across SQLite databases.
type CommitDbCrossReferences struct {
	HasSplitDbRecord   bool   `json:"hasSplitDbRecord"`
	SplitDbPath        string `json:"splitDbPath,omitempty"`
	SplitDbRecordId    int64  `json:"splitDbRecordId,omitempty"`
	HasHeatmapRecord   bool   `json:"hasHeatmapRecord"`
	HeatmapDbPath      string `json:"heatmapDbPath,omitempty"`
	HeatmapActionCount int    `json:"heatmapActionCount"`
	HasAgentTask       bool   `json:"hasAgentTask"`
	AgentTaskDbPath    string `json:"agentTaskDbPath,omitempty"`
	SubtaskCode        string `json:"subtaskCode,omitempty"`
	SubtaskStatus      string `json:"subtaskStatus,omitempty"`
}

// CommitInfoReport represents the unified commit report payload.
type CommitInfoReport struct {
	FullSha            string                  `json:"fullSha"`
	ShortSha           string                  `json:"shortSha"`
	AuthorName         string                  `json:"authorName"`
	AuthorEmail        string                  `json:"authorEmail"`
	CommittedAt        time.Time               `json:"committedAt"`
	FormattedDate      string                  `json:"formattedDate"`
	RelativeDate       string                  `json:"relativeDate"`
	CommitMessage      string                  `json:"commitMessage"`
	BranchName         string                  `json:"branchName"`
	ParentShas         []string                `json:"parentShas"`
	RepoSlug           string                  `json:"repoSlug"`
	RepoPath           string                  `json:"repoPath"`
	Provenance         CommitProvenance        `json:"provenance"`
	CrossReferences    CommitDbCrossReferences `json:"crossReferences"`
	FileDeltas         []CommitFileDelta       `json:"fileDeltas"`
	TotalFiles         int                     `json:"totalFiles"`
	TotalInsertions    int                     `json:"totalInsertions"`
	TotalDeletions     int                     `json:"totalDeletions"`
	CognitiveIntention string                  `json:"cognitiveIntention,omitempty"`
	IsExpectationMet   bool                    `json:"isExpectationMet"`
	SpecPath           string                  `json:"specPath,omitempty"`
	SubtaskCount       int                     `json:"subtaskCount,omitempty"`
	PlanCount          int                     `json:"planCount,omitempty"`
}

// CommitInfoOptions encapsulates parsed CLI options for commit info inspection.
type CommitInfoOptions struct {
	Query        string
	RepoPath     string
	IsJsonFormat bool
	IsVerbose    bool
}

// UniversalShaQuery encapsulates a query for commit SHA resolution.
type UniversalShaQuery struct {
	Query    string
	RepoPath string
}
