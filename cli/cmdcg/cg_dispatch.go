// Package cmdcg - cg_dispatch.go: action dispatch for cg subcommands.
// Split from cg.go on 31-Oct-26 so cg.go stays under the 300-line file-size gate.
package cmdcg

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
)

func hasDispatchedCGStatus(action string, repos []string) bool {
	switch action {
	case "help":
		RenderCGHelp()
	case "version", "ver", "-v":
		PrintCGVersion(repos)
	case "status", "stat", "ls":
		PrintCGStatus(repos)
	default:

		return false
	}

	return true
}

func hasDispatchedCGDiagnostic(action string, repos []string) bool {
	switch action {
	case "audit", "check":
		RunCGAudit(repos, false)
	case "doctor", "doc":
		RunCGDoctor(repos)
	case "diff", "drift":
		RunCGDiff(repos)
	default:

		return false
	}

	return true
}

func hasDispatchedCGMutation(action string, repos []string, opts cgOptions) bool {
	switch action {
	case "sync", "sy":
		RunCGSync(repos, opts.DryRun)
	case "prompt-sync", "cg-prompt-sync", "prompts-sync":
		RunCGPromptSync(repos, opts.DryRun)
	case "skill-sync", "cg-skill-sync", "skills-sync":
		RunCGSkillSync(repos, opts.DryRun)
	case "update":
		runCGUpdateAction(repos)
	case "install":
		executeCGWorkers(repos)
	default:

		return false
	}

	return true
}

func hasDispatchedCGVersionInit(action string, repos []string, opts cgOptions) bool {
	if action != "install-version-json" && action != "install-version" && action != "init-version" && action != "version-init" {
		return false
	}

	runCGInstallVersionJSON(repos, opts.VersionValue, opts.DryRun)

	return true
}

func hasDispatchedCGPrompt(action string, repos []string, opts cgOptions) bool {
	switch action {
	case "install-prompts", "install-prompt", "prompts-install", "update-prompts", "update-prompt", "prompts-update":
		runCGInstallPrompts(repos, opts.DryRun)
	case "prompts-status", "prompts-ls", "prompt-status":
		cmdprompt.RunPromptStatus(repos)
	case "prompts-version", "prompt-version":
		cmdprompt.RunPromptVersion(repos)
	default:

		return false
	}

	return true
}

func hasDispatchedAnyCGAction(opts cgOptions, repos []string) bool {
	hasStatus := hasDispatchedCGStatus(opts.Action, repos)
	hasDiag := hasDispatchedCGDiagnostic(opts.Action, repos)
	hasMut := hasDispatchedCGMutation(opts.Action, repos, opts)
	hasVer := hasDispatchedCGVersionInit(opts.Action, repos, opts)
	hasPrompt := hasDispatchedCGPrompt(opts.Action, repos, opts)

	return hasStatus || hasDiag || hasMut || hasVer || hasPrompt
}

func dispatchCGAction(opts cgOptions, repos []string, args []string) error {
	if hasDispatchedAnyCGAction(opts, repos) {
		return nil
	}

	if opts.Action == "general" || opts.Action == "general-install" || opts.Action == "install-general" {
		return RunCGGeneral(args)
	}

	RenderCGHelp()

	return apperror.NewValidationError(fmt.Sprintf("unknown action: %s", opts.Action))
}
