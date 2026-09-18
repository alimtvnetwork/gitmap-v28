package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// RunnerETAInfo stores execution countdown metadata from the local runner.
type RunnerETAInfo struct {
	Status                string `json:"status"`
	EtaSeconds            int    `json:"eta_seconds"`
	ElapsedSeconds        int    `json:"elapsed_seconds"`
	TotalEstimatedSeconds int    `json:"total_estimated_seconds"`
	UpdatedAt             int64  `json:"updated_at"`
}

// ReadRunnerETA attempts to read .ai-memory/temp/runner-eta.json from repo root.
func ReadRunnerETA() (*RunnerETAInfo, bool) {
	root, err := gitutil.RepoRoot(".")
	if err != nil {
		root = "."
	}

	etaPath := filepath.Join(root, ".ai-memory", "temp", "runner-eta.json")
	data, err := os.ReadFile(etaPath)
	if err != nil {
		return nil, false
	}

	var info RunnerETAInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, false
	}

	return &info, true
}

// WaitForRunnerETAIfActive polls and displays countdown if runner-eta.json is active.
func WaitForRunnerETAIfActive() {
	info, hasInfo := ReadRunnerETA()
	if !hasInfo || !isRunnerActive(info) {
		return
	}

	runRunnerCountdownLoop(info)
}

func isRunnerActive(info *RunnerETAInfo) bool {
	return info.Status == "running" && info.EtaSeconds > 0
}

func runRunnerCountdownLoop(info *RunnerETAInfo) {
	fmt.Printf("%s⏳ [runner-eta] Waiting for runner completion (ETA ~%ds)...%s\n",
		constants.ColorYellow, info.EtaSeconds, constants.ColorReset)

	remaining := info.EtaSeconds
	for remaining > 0 {
		if isTimelineTestMode() {
			break
		}

		time.Sleep(1 * time.Second)
		remaining--
		if remaining%5 == 0 && remaining > 0 {
			fmt.Printf("  %s⏳ ETA remaining: %ds...%s\n", constants.ColorYellow, remaining, constants.ColorReset)
		}
	}

	fmt.Printf("%s✓ Runner countdown finished. Surface error logs:%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}
