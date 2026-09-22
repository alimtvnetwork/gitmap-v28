package cmdssh

import (
	"sync"
	"time"
)

// SSHTraceStep captures a discrete action executed during an SSH operation.
type SSHTraceStep struct {
	Index         int    `json:"index"`
	Name          string `json:"name"`
	Details       string `json:"details"`
	Status        string `json:"status"`
	InternalError string `json:"internal_error,omitempty"`
}

// SSHExecutionTrace aggregates the end-to-end execution path of an SSH connection.
type SSHExecutionTrace struct {
	Timestamp     time.Time      `json:"timestamp"`
	Operation     string         `json:"operation"`
	Target        string         `json:"target"`
	User          string         `json:"user"`
	Host          string         `json:"host"`
	Port          int            `json:"port"`
	Steps         []SSHTraceStep `json:"steps"`
	InternalError string         `json:"internal_error,omitempty"`
	RawError      string         `json:"raw_error,omitempty"`
	LogPath       string         `json:"log_path,omitempty"`
	Suggestion    string         `json:"suggestion,omitempty"`
}

var (
	traceMu       sync.Mutex
	activeTrace   *SSHExecutionTrace
	lastFailTrace *SSHExecutionTrace
)

// BeginSSHTrace initializes a new active execution trace for an SSH operation.
func BeginSSHTrace(op, target, user, host string, port int) *SSHExecutionTrace {
	traceMu.Lock()
	defer traceMu.Unlock()

	activeTrace = &SSHExecutionTrace{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Target:    target,
		User:      user,
		Host:      host,
		Port:      port,
		Steps:     make([]SSHTraceStep, 0, 8),
	}

	return activeTrace
}

// GetActiveSSHTrace retrieves the currently active trace, or nil if none.
func GetActiveSSHTrace() *SSHExecutionTrace {
	traceMu.Lock()
	defer traceMu.Unlock()

	return activeTrace
}

// GetLastSSHTrace retrieves the last failed SSH execution trace.
func GetLastSSHTrace() *SSHExecutionTrace {
	traceMu.Lock()
	defer traceMu.Unlock()

	return lastFailTrace
}

// AddStep appends a milestone or diagnostic step to the execution trace.
func (t *SSHExecutionTrace) AddStep(name, details, status string, err error) {
	if t == nil {
		return
	}

	traceMu.Lock()
	defer traceMu.Unlock()

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	step := SSHTraceStep{
		Index:         len(t.Steps) + 1,
		Name:          name,
		Details:       details,
		Status:        status,
		InternalError: errStr,
	}

	t.Steps = append(t.Steps, step)
}

// SetInternalError records the root cause internal error and actionable suggestion.
func (t *SSHExecutionTrace) SetInternalError(rawErr error, suggestion string) {
	if t == nil {
		return
	}

	traceMu.Lock()
	defer traceMu.Unlock()

	if rawErr != nil {
		t.InternalError = rawErr.Error()
		t.RawError = rawErr.Error()
	}

	t.Suggestion = suggestion
	lastFailTrace = t
}
