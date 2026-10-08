package cmdrun

import (
	"time"
)

// SupportedRunExtensions lists priority extensions probed when no extension is specified.
var SupportedRunExtensions = []string{".py", ".ps1", ".sh", ".js", ".ts", ".go"}

// RunTarget represents a resolved executable script file.
type RunTarget struct {
	RawInput      string
	ResolvedPath  string
	Extension     string
	Interpreter   string
	IsDirectMatch bool
}

// RunOptions configures process execution parameters.
type RunOptions struct {
	Args       []string
	Timeout    time.Duration
	WorkingDir string
}

// RunResult captures process termination output and telemetry.
type RunResult struct {
	Target     RunTarget
	ExitCode   int
	DurationMs int64
	Stdout     string
	Stderr     string
	Err        error
}

// RunErrorRecord represents a persisted execution failure in SQLite.
type RunErrorRecord struct {
	ErrorID       string    `json:"error_id"`
	FilePath      string    `json:"file_path"`
	FileExtension string    `json:"file_extension"`
	Interpreter   string    `json:"interpreter"`
	ExitCode      int       `json:"exit_code"`
	DurationMs    int64     `json:"duration_ms"`
	ErrorMessage  string    `json:"error_message"`
	StdoutSnippet string    `json:"stdout_snippet"`
	StderrSnippet string    `json:"stderr_snippet"`
	ExecutedArgs  string    `json:"executed_args"`
	CreatedAt     time.Time `json:"created_at"`
}
