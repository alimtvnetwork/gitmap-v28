// Package cmd — ct.go implements the top-level ct (Custom Tools / Prompt Architect) command dispatcher.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runCT handles `gitmap ct [install-prompts|update-prompts|status|version]`.
//
//nolint:unused
func runCT(args []string) error {
	opts := parsePromptArgs(args)

	switch opts.Action {
	case "status", "prompts-status":
		runPromptStatus(opts.Targets)
		return nil
	case "version", "prompts-version":
		runPromptVersion(opts.Targets)
		return nil
	}

	targetDirs := resolveCTTargetDirs(opts)

	targetDirs = FilterPromptExclusions(targetDirs, opts.Exclude)
	if len(targetDirs) == 0 {
		fmt.Println("No target repositories found to install Prompt Architect.")
		return nil
	}

	fmt.Printf("→ Installing Prompt Architect v2 in %d repository(ies)...\n", len(targetDirs))
	var results []model.PromptInstallResult

	for _, dir := range targetDirs {
		name := filepath.Base(dir)
		fmt.Printf("  • %s (%s)... ", name, dir)
		res := ExecuteSinglePromptInstall(dir, opts.IsDryRun)
		results = append(results, res)
		if res.IsSuccess {
			fmt.Println("✓ Done")
		} else {
			fmt.Printf("✖ Failed: %s\n", res.Error)
		}
	}

	RenderPromptInstallSummary(results)
	ReportPromptFailures(results)
	return nil
}

//nolint:unused
func resolveCTTargetDirs(opts promptInstallOptions) []string {
	if len(opts.Targets) > 0 {
		return resolvePromptTargetsList(opts.Targets)
	}
	if opts.IsAll {
		resolved, _ := ResolveAllWorkDirPromptTargets()
		return resolved
	}
	resolved, _ := ResolvePromptTarget("")
	return resolved
}

//nolint:unused
func resolvePromptTargetsList(targets []string) []string {
	var targetDirs []string
	for _, t := range targets {
		resolved, err := ResolvePromptTarget(t)
		if err == nil {
			targetDirs = append(targetDirs, resolved...)
		}
	}
	return targetDirs
}
