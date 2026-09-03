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

// PipelineSplitDB encapsulates an isolated SQLite database connection for a single repository's pipeline data.
type PipelineSplitDB struct {
	conn     *sql.DB
	RepoSlug string
	Path     string
}

// PipelineRunRecord represents a recorded workflow run in the pipeline database.
type PipelineRunRecord struct {
	RunID           int64  `json:"runId"`
	RepoSlug        string `json:"repoSlug"`
	WorkflowName    string `json:"workflowName"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	Branch          string `json:"branch"`
	Sha             string `json:"sha"`
	EtaSeconds      int    `json:"etaSeconds"`
	DurationSeconds int    `json:"durationSeconds"`
	RunURL          string `json:"runUrl"`
	IsSuccess       bool   `json:"isSuccess"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// PipelineErrorRecord represents an isolated error diagnostic record.
type PipelineErrorRecord struct {
	RunID        int64  `json:"runId"`
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	StepName     string `json:"stepName"`
	ErrorText    string `json:"errorText"`
	RawLogs      string `json:"rawLogs,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Comments     string `json:"comments,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

// PipelineDBStats encapsulates health and sizing metrics for the pipeline split database.
type PipelineDBStats struct {
	Path          string `json:"path"`
	Size          int64  `json:"size"`
	TotalRuns     int    `json:"totalRuns"`
	SuccessRuns   int    `json:"successRuns"`
	FailedRuns    int    `json:"failedRuns"`
	ErrorLogCount int    `json:"errorLogCount"`
	SegmentCount  int    `json:"segmentCount"`
	LastUpdated   string `json:"lastUpdated"`
}

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

// PipelineDBDir returns the dedicated directory where pipeline split DBs live.
func PipelineDBDir() string {
	dir := filepath.Join(store.BinaryDataDir(), "pipeline_db")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// PipelineDBPath returns the full SQLite database file path for a repository.
func PipelineDBPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	return filepath.Join(PipelineDBDir(), "pipeline_"+slug+".db")
}

// OpenPipelineSplitDB opens or initializes the split SQLite database for a repo.
func OpenPipelineSplitDB(repoSlug string) (*PipelineSplitDB, error) {
	dbPath := PipelineDBPath(repoSlug)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open pipeline split db "+repoSlug)
	}
	conn.SetMaxOpenConns(1)
	p := &PipelineSplitDB{conn: conn, RepoSlug: repoSlug, Path: dbPath}
	if err := p.InitSchema(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return p, nil
}

// InitSchema ensures all pipeline tables exist.
func (p *PipelineSplitDB) InitSchema() error {
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
func (p *PipelineSplitDB) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
