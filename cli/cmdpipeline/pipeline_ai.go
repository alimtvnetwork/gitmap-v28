package cmdpipeline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// runPipelineAI handles gitmap pipeline-ai commands with automatic delays.
func runPipelineAI(args []string) error {
	checkHelp("pipeline-ai", args)
	subcmd, subArgs := extractPipelineAISubcmd(args)
	if isErrorLogsSubcmd(subcmd) {
		return handlePipelineAIErrors(subArgs)
	}

	return handlePipelineAIStatus(subArgs)
}

func extractPipelineAISubcmd(args []string) (string, []string) {
	if len(args) == 0 {
		return "status", args
	}

	first := strings.ToLower(args[0])
	if strings.HasPrefix(first, "-") {
		return "status", args
	}

	return first, args[1:]
}

func handlePipelineAIStatus(args []string) error {
	delaySeconds, subArgs := parsePipelineAIDelay(args)
	executePipelineAIDelay(delaySeconds, subArgs)

	repo := resolveCurrentRepoSlug()
	runs := queryWorkflowRuns(repo)
	payload := buildStatusPayload(repo, queryLatestTagRelease(repo), queryPendingPRs(repo), runs)
	payload.SleepSeconds = delaySeconds
	resolvePipelineAINextCommand(&payload)

	recordPipelineInDB(payload, runs)

	return outputPipelineAIResult(payload, subArgs)
}

func resolvePipelineAINextCommand(p *PipelineStatusPayload) {
	if p.HasErrors {
		p.NextAiCommand = "gitmap pipeline fix agy"
		p.IsStopWaiting = true
		p.RecommendedAction = "fix_errors"

		return
	}

	if p.IsRunning {
		p.NextAiCommand = fmt.Sprintf("gitmap pipeline-ai status -t %d", p.EtaSeconds)
		p.IsStopWaiting = false
		p.RecommendedAction = "wait"

		return
	}

	p.NextAiCommand = ""
	p.RecommendedAction = "none"
}

func handlePipelineAIErrors(args []string) error {
	repo := resolveCurrentRepoSlug()
	payload, report, hasFailures := FetchPipelineErrorReportWithMeta(repo, false)
	if hasArgFlag(args, "--json") {
		return outputPipelineAIErrorsJSON(repo, hasFailures, report, payload)
	}

	renderPipelineAIErrorsTerminal(report, hasFailures)

	return nil
}

func outputPipelineAIErrorsJSON(repo string, hasFailures bool, report string, payload PipelineErrorLogsPayload) error {
	action := resolveTernaryAction(hasFailures, "fix_errors", "none")
	cmd := resolveTernaryAction(hasFailures, "gitmap pipeline fix agy", "")
	out := map[string]any{
		"repo":              repo,
		"hasErrors":         hasFailures,
		"recommendedAction": action,
		"nextAiCommand":     cmd,
		"errorReport":       report,
		"payload":           payload,
	}

	return printJSON(out)
}

func resolveTernaryAction(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}

	return falseVal
}

func parsePipelineAIDelay(args []string) (int, []string) {
	delaySeconds := 20
	subArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		isDelayFlag := isPipelineAIDelayFlag(args[i])
		hasNextArg := i+1 < len(args)
		if isDelayFlag && hasNextArg {
			delaySeconds = extractDelaySeconds(args[i+1], delaySeconds)
			i++
			continue
		}

		subArgs = append(subArgs, args[i])
	}

	if delaySeconds < 20 {
		delaySeconds = 20
	}

	return delaySeconds, subArgs
}

func extractDelaySeconds(argVal string, defaultDelay int) int {
	v, err := strconv.Atoi(argVal)
	if err == nil && v >= 0 {
		return v
	}

	return defaultDelay
}

func isPipelineAIDelayFlag(arg string) bool {
	return arg == "-t" || arg == "--time" || arg == "-d" || arg == "--delay"
}

func isSkipDelayRequested() bool {
	if os.Getenv("GITMAP_SKIP_DELAY") == "1" {
		return true
	}

	if len(os.Getenv("CI")) > 0 {
		return true
	}

	if len(os.Getenv("GITHUB_ACTIONS")) > 0 {
		return true
	}

	return false
}

func executePipelineAIDelay(delaySeconds int, subArgs []string) {
	if delaySeconds <= 0 || isSkipDelayRequested() {
		return
	}

	isJSON := hasArgFlag(subArgs, "--json")
	msg := fmt.Sprintf("⏳ Auto-delaying %ds before checking pipeline status...\n", delaySeconds)
	if isJSON {
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Print(msg)
	}

	time.Sleep(time.Duration(delaySeconds) * time.Second)
}

func outputPipelineAIResult(payload PipelineStatusPayload, subArgs []string) error {
	if hasArgFlag(subArgs, "--json") {
		return printJSON(payload)
	}

	renderPipelineAIStatusTerminal(payload)

	return nil
}
