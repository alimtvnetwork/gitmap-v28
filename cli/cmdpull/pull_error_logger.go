package cmdpull

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var pullErrorLogMu sync.Mutex

// PullErrorLogEntry represents a structured JSON log entry for a pull error.
type PullErrorLogEntry struct {
	ErrorID        string    `json:"error_id"`
	RepoSlug       string    `json:"repo_slug"`
	RepoPath       string    `json:"repo_path"`
	NodeID         string    `json:"node_id"`
	NodeVersion    string    `json:"node_version"`
	ErrorType      string    `json:"error_type"`
	ErrorText      string    `json:"error_text"`
	StackTrace     string    `json:"stack_trace,omitempty"`
	RemediationCmd string    `json:"remediation_cmd,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// AppendPullErrorLog appends a structured error entry to .gitmap/logs/pull-errors.log as a single-line JSON.
func AppendPullErrorLog(entry PullErrorLogEntry) error {
	pullErrorLogMu.Lock()
	defer pullErrorLogMu.Unlock()

	prepared := preparePullLogEntry(entry)
	logDir := filepath.Join(".gitmap", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create pull log dir: %w", err)
	}

	logPath := filepath.Join(logDir, "pull-errors.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open pull errors log: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(prepared)
	if err != nil {
		return fmt.Errorf("marshal pull error entry: %w", err)
	}

	_, writeErr := f.Write(append(data, '\n'))
	return writeErr
}

// LogPullErrorToFile maps a store.PullErrorRecord and appends it to .gitmap/logs/pull-errors.log.
func LogPullErrorToFile(rec store.PullErrorRecord) error {
	entry := PullErrorLogEntry{
		ErrorID:        rec.ErrorID,
		RepoSlug:       rec.RepoSlug,
		RepoPath:       rec.RepoPath,
		NodeID:         rec.NodeID,
		NodeVersion:    rec.NodeVersion,
		ErrorType:      rec.ErrorType,
		ErrorText:      rec.ErrorText,
		StackTrace:     rec.StackTrace,
		RemediationCmd: rec.RemediationCmd,
		CreatedAt:      rec.CreatedAt,
	}
	return AppendPullErrorLog(entry)
}

func preparePullLogEntry(entry PullErrorLogEntry) PullErrorLogEntry {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	} else {
		entry.CreatedAt = entry.CreatedAt.UTC()
	}
	if entry.ErrorID == "" {
		entry.ErrorID = fmt.Sprintf("err-%d", time.Now().UnixNano())
	}
	return entry
}
