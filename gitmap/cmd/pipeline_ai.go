package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runPipelineAI handles gitmap pipeline-ai commands with automatic delays.
func runPipelineAI(args []string) error {
	checkHelp("pipeline-ai", args)
	delaySeconds, subArgs := parsePipelineAIDelay(args)
	executePipelineAIDelay(delaySeconds, subArgs)

	repo := resolveCurrentRepoSlug()
	runs := queryWorkflowRuns(repo)
	pendingPRs := queryPendingPRs(repo)
	lastTag := queryLatestTagRelease(repo)

	payload := buildStatusPayload(repo, lastTag, pendingPRs, runs)
	payload.SleepSeconds = delaySeconds
	if payload.IsRunning {
		payload.NextAiCommand = fmt.Sprintf("gitmap pipeline-ai status -t %d", payload.EtaSeconds)
	}
	recordPipelineInDB(payload, runs)

	return outputPipelineAIResult(payload, subArgs)
}

func parsePipelineAIDelay(args []string) (int, []string) {
	delaySeconds := 20
	subArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if isPipelineAIDelayFlag(args[i]) && i+1 < len(args) {
			if v, err := strconv.Atoi(args[i+1]); err == nil && v >= 0 {
				delaySeconds = v
			}
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

func isPipelineAIDelayFlag(arg string) bool {
	return arg == "-t" || arg == "--time" || arg == "-d" || arg == "--delay"
}

func executePipelineAIDelay(delaySeconds int, subArgs []string) {
	if delaySeconds <= 0 {
		return
	}
	isJSON := hasArgFlag(subArgs, "--json")
	msg := fmt.Sprintf("⏳ Auto-delaying %ds before checking pipeline status...\n", delaySeconds)
	if isJSON {
		fmt.Fprint(os.Stderr, msg)
	} else {
		fmt.Print(msg)
	}
	_ = checkSkipDelayEnv(delaySeconds)
}

func checkSkipDelayEnv(delaySeconds int) bool {
	if os.Getenv("GITMAP_SKIP_DELAY") == "1" {
		return true
	}
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	return false
}

func outputPipelineAIResult(payload PipelineStatusPayload, subArgs []string) error {
	isJSON := hasArgFlag(subArgs, "--json")
	if isJSON {
		return printJSON(payload)
	}
	renderPipelineAIStatusTerminal(payload)
	return nil
}

func renderPipelineAIStatusTerminal(p PipelineStatusPayload) {
	renderPipelineStatusTerminal(p)
	if p.IsRunning {
		fmt.Println()
		fmt.Printf("  %s🞠 AI Automation Next Action:%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("     Run: %s%s%s\n", constants.ColorGreen, p.NextAiCommand, constants.ColorReset)
		fmt.Printf("     (Automatically delays %ds then queries status)\n", p.EtaSeconds)
		fmt.Println()
	}
}
