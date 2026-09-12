// Package cmd — ct.go implements the top-level ct (Custom Tools / Prompt Architect) command dispatcher.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// runCT handles `gitmap ct [install-prompts|update-prompts|status|version]`.
//
//nolint:unused
func runCT(args []string) error {
	opts := cmdprompt.ParsePromptArgs(args)

	switch opts.Action {
	case "status", "prompts-status":
		cmdprompt.RunPromptStatus(opts.Targets)

		return nil
	case "version", "prompts-version":
		cmdprompt.RunPromptVersion(opts.Targets)

		return nil
	}

	targetDirs := resolveCTTargetDirs(opts)

	targetDirs = cmdprompt.FilterPromptExclusions(targetDirs, opts.Exclude)
	if len(targetDirs) == 0 {
		fmt.Println("No target repositories found to install Prompt Architect.")

		return nil
	}

	fmt.Printf("→ Installing Prompt Architect v2 in %d repository(ies)...\n", len(targetDirs))
	var results []model.PromptInstallResult

	for _, dir := range targetDirs {
		name := filepath.Base(dir)
		fmt.Printf("  • %s (%s)... ", name, dir)
		res := cmdprompt.ExecuteSinglePromptInstall(dir, opts.IsDryRun)
		results = append(results, res)
		if res.IsSuccess {
			fmt.Println("✓ Done")
		} else {
			fmt.Printf("✖ Failed: %s\n", res.Error)
		}
	}

	cmdprompt.RenderPromptInstallSummary(results)
	cmdprompt.ReportPromptFailures(results)

	return nil
}

//nolint:unused
func resolveCTTargetDirs(opts cmdprompt.PromptInstallOptions) []string {
	if len(opts.Targets) > 0 {
		return resolvePromptTargetsList(opts.Targets)
	}

	if opts.IsAll {
		resolved, _ := cmdprompt.ResolveAllWorkDirPromptTargets()

		return resolved
	}

	resolved, _ := cmdprompt.ResolvePromptTarget("")

	return resolved
}

//nolint:unused
func resolvePromptTargetsList(targets []string) []string {
	var targetDirs []string
	for _, t := range targets {
		resolved, err := cmdprompt.ResolvePromptTarget(t)
		if err == nil {
			targetDirs = append(targetDirs, resolved...)
		}
	}

	return targetDirs
}
