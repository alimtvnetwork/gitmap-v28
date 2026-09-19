package cmdautomation

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RuntimeRecord caches discovered CLI interpreters, compilers, and shells.
type RuntimeRecord struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	BinaryName        string `json:"binaryName"`
	BinaryPath        string `json:"binaryPath"`
	Version           string `json:"version"`
	Status            string `json:"status"`
	InstallCmd        string `json:"installCmd"`
	ProfileSuggestion string `json:"profileSuggestion"`
	FallbackCmd       string `json:"fallbackCmd"`
	DiscoveredAt      string `json:"discoveredAt"`
	LastVerifiedAt    string `json:"lastVerifiedAt"`
	IsValid           bool   `json:"isValid"`
}

// FileManifestItem represents an indexed file in the look-ahead cache.
type FileManifestItem struct {
	ID                       int64  `json:"id"`
	RelativePath             string `json:"relativePath"`
	FileName                 string `json:"fileName"`
	FileExtension            string `json:"fileExtension"`
	ParentFolderPath         string `json:"parentFolderPath"`
	AbsoluteFilePath         string `json:"absoluteFilePath"`
	AbsoluteParentFolderPath string `json:"absoluteParentFolderPath"`
	FileSize                 int64  `json:"fileSize"`
	ModifiedTimestamp        int64  `json:"modifiedTimestamp"`
	ContentHash              string `json:"contentHash"`
	IsBinary                 bool   `json:"isBinary"`
	ScannedAt                string `json:"scannedAt"`
}

// ExecutionHistoryRecord records execution audits, throughput metrics, and exit statuses.
type ExecutionHistoryRecord struct {
	ID               int64     `json:"id"`
	Runtime          string    `json:"runtime"`
	CommandType      string    `json:"commandType"`
	ScriptPreview    string    `json:"scriptPreview"`
	FilesMatched     int       `json:"filesMatched"`
	FilesProcessed   int       `json:"filesProcessed"`
	WorkersUsed      int       `json:"workersUsed"`
	ThreadsPerWorker int       `json:"threadsPerWorker"`
	Encoding         string    `json:"encoding"`
	DurationMs       int64     `json:"durationMs"`
	ExitCode         int       `json:"exitCode"`
	ExecutedAt       time.Time `json:"executedAt"`
}

// WorkerRunOptions configures polyglot worker pool execution.
type WorkerRunOptions struct {
	Runtime          string `json:"runtime"`
	CommandType      string `json:"commandType"`
	Script           string `json:"script"`
	Workers          int    `json:"workers"`
	Threads          int    `json:"threads"`
	Encoding         string `json:"encoding"`
	TimeoutSec       int    `json:"timeoutSec"`
	IsPreRead        bool   `json:"isPreRead"`
	Dir              string `json:"dir"`
	IncludeBinaries  bool   `json:"includeBinaries"`
	IncludeLargeJson bool   `json:"includeLargeJson"`
}

// WorkerRunResult captures execution throughput, outputs, and exit status.
type WorkerRunResult struct {
	FilesMatched   int           `json:"filesMatched"`
	FilesProcessed int           `json:"filesProcessed"`
	WorkersUsed    int           `json:"workersUsed"`
	ThreadsUsed    int           `json:"threadsUsed"`
	Duration       time.Duration `json:"duration"`
	ExitCode       int           `json:"exitCode"`
	Stdout         string        `json:"stdout"`
	Stderr         string        `json:"stderr"`
	HasErrors      bool          `json:"hasErrors"`
}

// FileContext represents the rich metadata injected into polyglot workers.
type FileContext struct {
	FilePath                 string `json:"filePath"`
	FileName                 string `json:"fileName"`
	FileExtension            string `json:"fileExtension"`
	ParentFolderPath         string `json:"parentFolderPath"`
	AbsoluteFilePath         string `json:"absoluteFilePath"`
	AbsoluteParentFolderPath string `json:"absoluteParentFolderPath"`
	FileSize                 int64  `json:"fileSize"`
	ModifiedTimestamp        int64  `json:"modifiedTimestamp"`
	LineNumber               int    `json:"lineNumber"`
	MatchedContent           string `json:"matchedContent"`
	Content                  string `json:"content"`
	Encoding                 string `json:"encoding"`
}

type runtimeMeta struct {
	InstallCmd        string
	ProfileSuggestion string
	FallbackCmd       string
}

var packageFallbacks = map[string][4]string{
	"python": {"Python.Python.3.12", "python3", "python@3.12", "python3 python3-pip"},
	"node":   {"OpenJS.NodeJS", "nodejs", "node", "nodejs npm"},
	"go":     {"GoLang.Go", "golang", "go", "golang-go"},
	"rust":   {"Rustlang.Rustup", "rust", "rust", "cargo rustc"},
	"pwsh":   {"Microsoft.PowerShell", "powershell-core", "--cask powershell", "powershell"},
	"bash":   {"Git.Git", "git", "bash", "bash"},
}

func formatFallbackCmd(norm string) string {
	p, ok := packageFallbacks[norm]
	if !ok {
		return ""
	}
	return fmt.Sprintf("  Windows (winget) : winget install %s\n  Windows (choco)  : choco install %s\n  macOS (brew)     : brew install %s\n  Linux (Ubuntu)   : sudo apt update && sudo apt install -y %s\n", p[0], p[1], p[2], p[3])
}

func normalizeRuntimeKey(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	switch lower {
	case "py", "python3", "py-file":
		return "python"
	case "js", "ts", "node-file":
		return "node"
	case "golang", "go-file":
		return "go"
	case "rs", "cargo", "rs-file":
		return "rust"
	case "ps", "powershell", "ps-file":
		return "pwsh"
	case "sh", "sh-file":
		return "bash"
	}
	return lower
}

type (
	// WorkerRunResultMonad wraps WorkerRunResult with an AppError.
	WorkerRunResultMonad = result.Result[WorkerRunResult]

	// RuntimeListResultMonad wraps a slice of RuntimeRecord with an AppError.
	RuntimeListResultMonad = result.Result[[]RuntimeRecord]

	// RuntimeRecordMonad wraps a single RuntimeRecord with an AppError.
	RuntimeRecordMonad = result.Result[RuntimeRecord]
)
