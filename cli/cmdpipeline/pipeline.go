package cmdpipeline

import (
	"fmt"
	"strings"
)

// PipelineStatusPayload represents the status output.
type PipelineStatusPayload struct {
	IsRunning              bool            `json:"isRunning"`
	EtaSeconds             int             `json:"etaSeconds"`
	SleepSeconds           int             `json:"sleepSeconds,omitempty"`
	NextAiCommand          string          `json:"nextAiCommand,omitempty"`
	LastTagRelease         string          `json:"lastTagRelease"`
	PendingPipelines       int             `json:"pendingPipelines"`
	PendingTasks           int             `json:"pendingTasks"`
	PendingPRs             int             `json:"pendingPRs"`
	Repo                   string          `json:"repo"`
	ActiveWorkflow         string          `json:"activeWorkflow,omitempty"`
	LastStatus             string          `json:"lastStatus,omitempty"`
	LastConclusion         string          `json:"lastConclusion,omitempty"`
	LastRunId              uint64          `json:"lastRunId,omitempty"`
	LastRunUrl             string          `json:"lastRunUrl,omitempty"`
	UpdatedAt              string          `json:"updatedAt"`
	HasErrors              bool            `json:"hasErrors"`
	IsStopWaiting          bool            `json:"isStopWaiting"`
	RecommendedAction      string          `json:"recommendedAction,omitempty"`
	FailedJobCount         int             `json:"failedJobCount,omitempty"`
	ErrorSummary           string          `json:"errorSummary,omitempty"`
	ErrorLogs              string          `json:"errorLogs,omitempty"`
	ActionableErrorSnippet string          `json:"actionableErrorSnippet,omitempty"`
	FailedJobs             []FailedJobItem `json:"failedJobs,omitempty"`
	DbPath                 string          `json:"dbPath,omitempty"`
	DbSize                 string          `json:"dbSize,omitempty"`
	IsFromCache            bool            `json:"isFromCache,omitempty"`
	Runs                   []ghRunItem     `json:"runs,omitempty"`
}

// PipelineErrorLogsPayload represents error log outputs.
type PipelineErrorLogsPayload struct {
	Repo               string            `json:"repo"`
	RepoUrl            string            `json:"repoUrl,omitempty"`
	LastHash           string            `json:"lastHash,omitempty"`
	LastReleaseVersion string            `json:"lastReleaseVersion,omitempty"`
	LatestBranch       string            `json:"latestBranch,omitempty"`
	OpenPRsCount       int               `json:"openPRsCount"`
	WorkflowName       string            `json:"workflowName"`
	RunId              uint64            `json:"runId"`
	Status             string            `json:"status"`
	Conclusion         string            `json:"conclusion"`
	Branch             string            `json:"branch,omitempty"`
	Sha                string            `json:"sha,omitempty"`
	CreatedAt          string            `json:"createdAt,omitempty"`
	UpdatedAt          string            `json:"updatedAt,omitempty"`
	DurationSeconds    int               `json:"durationSeconds,omitempty"`
	SavedLogFile       string            `json:"savedLogFile,omitempty"`
	SavedReportFile    string            `json:"savedReportFile,omitempty"`
	DbPath             string            `json:"dbPath,omitempty"`
	EtaSeconds         int               `json:"etaSeconds,omitempty"`
	RerunEtaSeconds    int               `json:"rerunEtaSeconds,omitempty"`
	IsRunning          bool              `json:"isRunning"`
	ActiveRunName      string            `json:"activeRunName,omitempty"`
	ActiveRunId        uint64            `json:"activeRunId,omitempty"`
	ActiveRunUrl       string            `json:"activeRunUrl,omitempty"`
	ErrorLogs          string            `json:"errorLogs"`
	CombinedErrors     string            `json:"combinedErrors,omitempty"`
	Url                string            `json:"url,omitempty"`
	Notes              string            `json:"notes,omitempty"`
	FailedRuns         []FailedRunItem   `json:"failedRuns,omitempty"`
	FailedTree         string            `json:"failedTree,omitempty"`
	SectionFailures    []SectionFailure  `json:"sectionFailures,omitempty"`
	CICDChecks         []CICDCheckResult `json:"cicdChecks,omitempty"`
	Runs               []ghRunItem       `json:"runs,omitempty"`
	IsFromCache        bool              `json:"isFromCache,omitempty"`
	CacheSource        string            `json:"cacheSource,omitempty"`
}

// SectionFailure represents a discrete failing section or step across pipeline runs.
type SectionFailure struct {
	WorkflowName   string   `json:"workflowName"`
	RunId          uint64   `json:"runId"`
	JobName        string   `json:"jobName"`
	StepName       string   `json:"stepName"`
	FailureSummary string   `json:"failureSummary"`
	ErrorLines     []string `json:"errorLines"`
	Warnings       []string `json:"warnings,omitempty"`
	SavedLogFile   string   `json:"savedLogFile,omitempty"`
	CreatedAt      string   `json:"createdAt,omitempty"`
	StackTrace     string   `json:"stackTrace,omitempty"`
}

// FailedRunItem represents an individual failed workflow run.
type FailedRunItem struct {
	WorkflowName    string          `json:"workflowName"`
	RunId           uint64          `json:"runId"`
	Status          string          `json:"status,omitempty"`
	Conclusion      string          `json:"conclusion"`
	Branch          string          `json:"branch,omitempty"`
	Sha             string          `json:"sha,omitempty"`
	CreatedAt       string          `json:"createdAt,omitempty"`
	UpdatedAt       string          `json:"updatedAt,omitempty"`
	DurationSeconds int             `json:"durationSeconds,omitempty"`
	SavedLogFile    string          `json:"savedLogFile,omitempty"`
	SavedMetaFile   string          `json:"savedMetaFile,omitempty"`
	Url             string          `json:"url"`
	FailedJobs      []FailedJobItem `json:"failedJobs"`
	RawErrors       string          `json:"rawErrors,omitempty"`
	StackTrace      string          `json:"stackTrace,omitempty"`
}

// FailedJobItem represents a failed job and step within a workflow run.
type FailedJobItem struct {
	JobName        string   `json:"jobName"`
	StepName       string   `json:"stepName"`
	FailureSummary string   `json:"failureSummary"`
	ErrorLines     []string `json:"errorLines"`
	Warnings       []string `json:"warnings,omitempty"`
	StackTrace     string   `json:"stackTrace,omitempty"`
}

// ErrorLogOutputParams encapsulates parameters for outputting error logs.
type ErrorLogOutputParams struct {
	Payload              PipelineErrorLogsPayload
	IsJSON               bool
	HasSuppressOutputLog bool
	FilePath             string
	TempFile             string
}

type ghRunItem struct {
	DatabaseId   uint64 `json:"databaseId"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	HeadBranch   string `json:"headBranch"`
	HeadSha      string `json:"headSha"`
	Url          string `json:"url"`
	DisplayTitle string `json:"displayTitle,omitempty"`
	Event        string `json:"event,omitempty"`
}

// runPipeline is the entry point for the gitmap pipeline command group.
func runPipeline(args []string) error {
	if len(args) == 0 {
		return handlePipelineStatus(nil)
	}

	if IsPipelineFixAgyArgs(args) && PipelineAgyFixRunner != nil {
		return PipelineAgyFixRunner(args)
	}

	subcmd := strings.ToLower(args[0])
	if handled, err := checkErrorLogsSubcmd(subcmd, args); handled {
		return err
	}

	return dispatchPipelineSubcmd(subcmd, args)
}

func checkErrorLogsSubcmd(subcmd string, args []string) (bool, error) {
	if isDetailsSubcmd(subcmd) {
		return true, HandlePipelineDetails(args[1:])
	}
	if subcmd == "last-failed-logs" {
		return true, HandlePipelineLastFailedLogs(args[1:])
	}
	if IsNegativeIndexToken(subcmd) {
		return true, handlePipelineErrorLogs(args)
	}

	return checkErrorLogsOrClear(subcmd, args)
}

func checkErrorLogsOrClear(subcmd string, args []string) (bool, error) {
	if !isErrorLogsSubcmd(subcmd) {
		return false, nil
	}
	if len(args) > 1 && isPipelineClearAction(args[1]) {
		return true, handlePipelineDB(append([]string{"clear"}, args[2:]...))
	}

	return true, handlePipelineErrorLogs(args[1:])
}

func isDetailsSubcmd(subcmd string) bool {
	switch subcmd {
	case "details", "detail", "pd", "dt", "pipeline-details", "pipeline_details":
		return true
	}

	return false
}

func isPipelineClearAction(action string) bool {
	switch strings.ToLower(action) {
	case "clear", "reset", "clean", "clear-errors", "reset-errors":
		return true
	default:
		return false
	}
}

func isErrorLogsSubcmd(subcmd string) bool {
	switch subcmd {
	case "error-logs", "errorlogs", "error-log", "errorlog", "errors", "err",
		"errorslogs", "errors-log", "errors-logs", "pe":
		return true
	}

	return false
}

func isPipelineClearDbSubcmd(subcmd string) bool {
	switch subcmd {
	case "clear-db", "cleardb", "db-clear", "clear-errors", "errors-clear", "reset-errors", "clear", "clean", "reset":
		return true
	default:
		return false
	}
}

func dispatchPipelineSubcmd(subcmd string, args []string) error {
	if isPipelineClearDbSubcmd(subcmd) {
		return handlePipelineDB(append([]string{"clear"}, args[1:]...))
	}

	return dispatchCorePipelineSubcmd(subcmd, args)
}

func dispatchCorePipelineSubcmd(subcmd string, args []string) error {
	switch subcmd {
	case "status", "st", "s":
		return handlePipelineStatus(args[1:])
	case "details", "detail", "pd", "dt":
		return HandlePipelineDetails(args[1:])
	case "waittime", "wait", "eta", "wt", "wait-time":
		return handlePipelineWaitTime(args[1:])
	case "logs", "log", "l":
		return handlePipelineLogs(args[1:])
	case "history", "hist", "h":
		return handlePipelineHistory(args[1:])
	case "db":
		return handlePipelineDB(args[1:])
	case "help", "-h", "--help":
		return showPipelineHelp()
	default:
		return dispatchPipelineFallback(subcmd, args)
	}
}

func showPipelineHelp() error {
	printPipelineHelp()

	return nil
}

func dispatchPipelineFallback(subcmd string, args []string) error {
	if strings.HasPrefix(subcmd, "-") {
		return handlePipelineStatus(args)
	}

	printPipelineHelp()

	return fmt.Errorf("unknown pipeline subcommand: %s", subcmd)
}

func printPipelineHelp() {
	RenderPipelineHelp()
}
