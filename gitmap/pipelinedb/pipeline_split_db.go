package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	_ "modernc.org/sqlite"
)

// PipelineSplitDb encapsulates an isolated SQLite database connection for a single repository's pipeline data.
type PipelineSplitDb struct {
	conn     *sql.DB
	RepoSlug string
	Path     string
}

// PipelineSplitDB is an alias to PipelineSplitDb for backward compatibility.
type PipelineSplitDB = PipelineSplitDb

// PipelineRunRecord represents a recorded workflow run in the pipeline database.
type PipelineRunRecord struct {
	RunId           uint64 `json:"runId"`
	RepoSlug        string `json:"repoSlug"`
	WorkflowName    string `json:"workflowName"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	Branch          string `json:"branch"`
	Sha             string `json:"sha"`
	EtaSeconds      int    `json:"etaSeconds"`
	DurationSeconds int    `json:"durationSeconds"`
	RunUrl          string `json:"runUrl"`
	IsSuccess       bool   `json:"isSuccess"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// PipelineErrorRecord represents an isolated error diagnostic record.
type PipelineErrorRecord struct {
	RunId        uint64 `json:"runId"`
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	StepName     string `json:"stepName"`
	ErrorText    string `json:"errorText"`
	RawLogs      string `json:"rawLogs,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Comments     string `json:"comments,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

// PipelineDbStats encapsulates health and sizing metrics for the pipeline split database.
type PipelineDbStats struct {
	Path          string `json:"path"`
	Size          uint64 `json:"size"`
	TotalRuns     int    `json:"totalRuns"`
	SuccessRuns   int    `json:"successRuns"`
	FailedRuns    int    `json:"failedRuns"`
	ErrorLogCount int    `json:"errorLogCount"`
	SegmentCount  int    `json:"segmentCount"`
	LastUpdated   string `json:"lastUpdated"`
}

// PipelineDBStats is an alias to PipelineDbStats for backward compatibility.
type PipelineDBStats = PipelineDbStats

// SanitizeRepoSlug converts a repository slug into a valid safe filesystem name.
func SanitizeRepoSlug(repo string) string {
	lower := strings.ToLower(strings.TrimSpace(repo))
	reg := regexp.MustCompile(`[^a-z0-9_-]+`)
	slug := reg.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "pipeline-default"
	}
	return slug
}

// PipelineDbDir returns the dedicated directory where pipeline split DBs live.
func PipelineDbDir() string {
	dir := filepath.Join(store.BinaryDataDir(), "pipeline_db")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// PipelineDBDir is an alias to PipelineDbDir.
var PipelineDBDir = PipelineDbDir

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	return filepath.Join(PipelineDbDir(), "pipeline_"+slug+".db")
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath

// OpenPipelineSplitDb opens or initializes the split SQLite database for a repo.
func OpenPipelineSplitDb(repoSlug string) (*PipelineSplitDb, error) {
	dbPath := PipelineDbPath(repoSlug)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}
	conn.SetMaxOpenConns(1)
	p := &PipelineSplitDb{conn: conn, RepoSlug: repoSlug, Path: dbPath}
	if err := p.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return p, nil
}

// OpenPipelineSplitDB is an alias to OpenPipelineSplitDb.
var OpenPipelineSplitDB = OpenPipelineSplitDb

// InitSchema ensures all pipeline tables exist.
func (p *PipelineSplitDb) InitSchema() error {
	queries := []string{
		sqlCreatePipelineRun,
		sqlCreatePipelineErrorLog,
		sqlCreatePipelineSegment,
	}
	for _, q := range queries {
		if _, err := p.conn.Exec(q); err != nil {
			return apperror.WrapSimple(err, "init pipeline db schema")
		}
	}
	return nil
}

// Close closes the underlying SQLite connection.
func (p *PipelineSplitDb) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
