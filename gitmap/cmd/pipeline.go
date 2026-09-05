package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// PipelineStatusPayload represents the status output.
type PipelineStatusPayload struct {
	IsRunning        bool   `json:"isRunning"`
	EtaSeconds       int    `json:"etaSeconds"`
	SleepSeconds     int    `json:"sleepSeconds,omitempty"`
	NextAiCommand    string `json:"nextAiCommand,omitempty"`
	LastTagRelease   string `json:"lastTagRelease"`
	PendingPipelines int    `json:"pendingPipelines"`
	PendingTasks     int    `json:"pendingTasks"`
	PendingPRs       int    `json:"pendingPRs"`
	Repo             string `json:"repo"`
	ActiveWorkflow   string `json:"activeWorkflow,omitempty"`
	LastStatus       string `json:"lastStatus,omitempty"`
	LastConclusion   string `json:"lastConclusion,omitempty"`
	LastRunId        uint64 `json:"lastRunId,omitempty"`
	LastRunUrl       string `json:"lastRunUrl,omitempty"`
	UpdatedAt        string `json:"updatedAt"`
}

// PipelineErrorLogsPayload represents error log outputs.
type PipelineErrorLogsPayload struct {
	Repo            string            `json:"repo"`
	WorkflowName    string            `json:"workflowName"`
	RunId           uint64            `json:"runId"`
	Status          string            `json:"status"`
	Conclusion      string            `json:"conclusion"`
	EtaSeconds      int               `json:"etaSeconds,omitempty"`
	RerunEtaSeconds int               `json:"rerunEtaSeconds,omitempty"`
	IsRunning       bool              `json:"isRunning"`
	ErrorLogs       string            `json:"errorLogs"`
	Url             string            `json:"url,omitempty"`
	CICDChecks      []CICDCheckResult `json:"cicdChecks,omitempty"`
}

// ErrorLogOutputParams encapsulates parameters for outputting error logs.
type ErrorLogOutputParams struct {
	Payload  PipelineErrorLogsPayload
	IsJSON   bool
	FilePath string
	TempFile string
}

type ghRunItem struct {
	DatabaseId uint64 `json:"databaseId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	HeadBranch string `json:"headBranch"`
	HeadSha    string `json:"headSha"`
	Url        string `json:"url"`
}

// runPipeline is the entry point for the gitmap pipeline command group.
func runPipeline(args []string) error {
	if len(args) == 0 {
		return handlePipelineStatus(nil)
	}

	subcmd := strings.ToLower(args[0])

	switch subcmd {
	case "status", "st", "s":
		return handlePipelineStatus(args[1:])
	case "waittime", "wait", "eta", "wt", "wait-time":
		return handlePipelineWaitTime(args[1:])
	case "error-logs", "errorlogs", "error-log", "errorlog", "errors", "err":
		return handlePipelineErrorLogs(args[1:])
	case "logs", "log", "l":
		return handlePipelineLogs(args[1:])
	case "db":
		return handlePipelineDB(args[1:])
	case "help", "-h", "--help":
		printPipelineHelp()

		return nil
	default:
		if strings.HasPrefix(subcmd, "-") {
			return handlePipelineStatus(args)
		}

		printPipelineHelp()

		return fmt.Errorf("unknown pipeline subcommand: %s", subcmd)
	}
}

func printPipelineHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap pipeline [command] [flags]")
	fmt.Println("  gitmap pipeline-ai [status|eta] [-t <seconds>] [--json]")
	fmt.Println("  gitmap pl [command] [flags]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Commands:" + constants.ColorReset)
	fmt.Println("  status                 Check live CI/CD pipeline status, ETA, and pending PRs")
	fmt.Println("  waittime               Output remaining ETA seconds for active pipeline (alias: eta)")
	fmt.Println("  eta                    Output remaining ETA seconds for active pipeline")
	fmt.Println("  error-logs             Display failure logs, rerun ETA, and internal CI/CD fix suite (alias: errorlogs)")
	fmt.Println("  logs                   Display all workflow logs in terminal")
	fmt.Println("  pipeline-ai status     Auto-delay (default: 20s or -t <seconds>) then query status")
	fmt.Println("  db                     Inspect or manage isolated pipeline split SQLite database")
	fmt.Println("  help                   Show this pipeline command suite documentation")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Flags:" + constants.ColorReset)
	fmt.Println("  -t, --timeline          Watch pipeline run until completion with dynamic ETA timeline")
	fmt.Println("  -f, --fix               Run internal CI/CD issue diagnostic and auto-repair scripts")
	fmt.Println("  -c, --check             Run internal CI/CD issue diagnostic checks without modifying files")
	fmt.Println("  -y, --yes               Auto-confirm prompts non-interactively")
	fmt.Println("  --json                  Output data in structured JSON format")
	fmt.Println("  --file <path>           Write error logs to specified file path")
	fmt.Println("  --tempfile <filename>   Write error logs to .lovable/temp/<filename>")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap pipeline status")
	fmt.Println("  gitmap pipeline status -t")
	fmt.Println("  gitmap pipeline errorlogs -t")
	fmt.Println("  gitmap pipeline error-logs -t --fix")
	fmt.Println("  gitmap pipeline errorlogs --check")
	fmt.Println("  gitmap pipeline error-logs --json")
}
