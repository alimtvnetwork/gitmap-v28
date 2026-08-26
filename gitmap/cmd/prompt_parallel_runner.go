// Package cmd — prompt_parallel_runner.go manages serial or parallel prompt execution.
package cmd

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/installer"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func ExecuteSinglePromptInstall(targetDir string, isDryRun bool) model.PromptInstallResult {
	start := time.Now()
	res := model.PromptInstallResult{
		RepoPath: targetDir,
	}

	if isDryRun {
		SimulatePromptInstallation(targetDir)
		res.IsSuccess = true
		res.Duration = "0.1s"
		return res
	}

	err := installer.RunPromptInstallerForHost(targetDir, 90*time.Second)
	res.Duration = time.Since(start).Round(time.Millisecond).String()

	if err != nil {
		res.IsSuccess = false
		res.Error = err.Error()
		return res
	}

	res.IsSuccess = true
	meta, _ := ReadPromptArchitectMetadata(targetDir)
	res.Version = meta.Version
	return res
}
