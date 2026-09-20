package cmdagy

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	isPingJSON    bool
	pingWorkspace string
)

var agyPingCmd = &cobra.Command{
	Use:     "ping",
	Aliases: []string{"check"},
	Short:   "Ping Antigravity IDE and inspect environment health",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPing()
	},
}

func init() {
	agyPingCmd.Flags().BoolVarP(&isPingJSON, "json", "j", false, "Output ping report in JSON format")
	agyPingCmd.Flags().StringVarP(&pingWorkspace, "workspace", "w", "", "Target workspace path to inspect")
}

func runAgyPing() error {
	report := ExecuteAgyPing(pingWorkspace)
	if isPingJSON {
		return renderPingJSON(report)
	}
	renderPingReport(report)

	return nil
}

// ExecuteAgyPing runs all diagnostic checks against the Antigravity IDE environment.
func ExecuteAgyPing(wsPath string) AgyPingReport {
	report := AgyPingReport{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Executable: checkIDEExecutable(),
		Process:    checkIDEProcess(),
		Filesystem: checkFilesystemHealth(),
		Workspace:  checkWorkspaceAndConversation(wsPath),
		Queue:      checkPromptQueue(),
	}
	report.IsHealthy = computePingHealth(report)

	return report
}

func computePingHealth(report AgyPingReport) bool {
	hasFs := report.Filesystem.IsAccessible
	hasBinary := report.Executable.IsFound || report.Process.IsRunning

	return hasFs && hasBinary
}

func checkIDEExecutable() AgyPingExecutableCheck {
	ideRes := ResolveAntigravityIDE()
	if ideRes.IsSuccess() {
		return AgyPingExecutableCheck{
			IsFound: true,
			Path:    ideRes.Value,
		}
	}

	return AgyPingExecutableCheck{
		IsFound: false,
		Error:   extractResultError(ideRes, "IDE executable not found"),
	}
}

func checkIDEProcess() AgyPingProcessCheck {
	procRes := DetectRunningAntigravityIDE()
	if procRes.IsSuccess() {
		return AgyPingProcessCheck{
			IsRunning: true,
			PID:       procRes.Value.PID,
			Name:      procRes.Value.Name,
		}
	}

	return AgyPingProcessCheck{
		IsRunning: false,
		Error:     extractResultError(procRes, "IDE process not running"),
	}
}

func extractResultError[T any](res result.Result[T], fallback string) string {
	if res.Error != nil {
		return res.Error.Error()
	}

	return fallback
}

func checkFilesystemHealth() AgyFilesystemHealth {
	brainDir, brainErr := GetBrainLogsDirPath()
	projDir, projErr := getProjectsDirPath()
	convDir, convErr := getConversationsDirPath()

	hasBrain := brainErr == nil && checkDirExists(brainDir)
	hasProj := projErr == nil && checkDirExists(projDir)
	hasConv := convErr == nil && checkDirExists(convDir)
	isAccessible := hasBrain && hasProj && hasConv

	return buildFilesystemHealth(isAccessible, brainDir, projDir, convDir, brainErr)
}

func buildFilesystemHealth(isAccessible bool, brain, proj, conv string, err error) AgyFilesystemHealth {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	return AgyFilesystemHealth{
		IsAccessible:     isAccessible,
		BrainDir:         brain,
		ProjectsDir:      proj,
		ConversationsDir: conv,
		Error:            errStr,
	}
}

func checkWorkspaceAndConversation(wsPath string) AgyPingWorkspaceCheck {
	target := resolvePingWorkspace(wsPath)
	state := DetectConversationExecutionState(target)
	convInfo := findMatchingConvInfo(target)

	return buildWorkspaceCheck(target, state, convInfo)
}

func resolvePingWorkspace(wsPath string) string {
	if len(wsPath) > 0 {
		return wsPath
	}
	cwd, err := os.Getwd()
	if err == nil && len(cwd) > 0 {
		return cwd
	}

	return resolveProjectRootDir()
}

func findMatchingConvInfo(wsPath string) AgyConvInfo {
	convs, err := scanAllConversations()
	if err != nil || len(convs) == 0 {
		return AgyConvInfo{}
	}
	pClean := cleanProjectWorkspace(wsPath)
	for _, c := range convs {
		if isConvPathMatch(pClean, c.CleanPath) {
			return c
		}
	}

	return AgyConvInfo{}
}

func buildWorkspaceCheck(target string, state AgyConversationExecutionState, info AgyConvInfo) AgyPingWorkspaceCheck {
	statusStr := string(state.Status)
	if len(statusStr) == 0 {
		statusStr = "unknown"
	}

	return AgyPingWorkspaceCheck{
		TargetWorkspace: target,
		ConvID:          state.ConvID,
		ConvStatus:      statusStr,
		StepCount:       info.StepCount,
		TranscriptPath:  state.TranscriptPath,
	}
}

func checkPromptQueue() AgyPingQueueCheck {
	q, err := GetQueueStatus()
	if err != nil {
		return AgyPingQueueCheck{
			HasQueue: false,
			Error:    err.Error(),
		}
	}

	return AgyPingQueueCheck{
		HasQueue:          true,
		ActivePromptCount: countActivePrompt(q.Active),
		QueuedCount:       len(q.Queued),
		LastUpdated:       q.UpdatedAt,
	}
}

func countActivePrompt(active *AgyPromptQueueEntry) int {
	if active != nil {
		return 1
	}

	return 0
}
