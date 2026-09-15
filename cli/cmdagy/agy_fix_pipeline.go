// Package cmdagy — agy_fix_pipeline.go feeds failing pipeline error logs and CI/CD fix prompt to Antigravity IDE.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	agyFixDetailed     bool
	agyFixNoRelease    bool
	agyFixCustomPrompt string
	agyFixNoClipboard  bool
	agyFixOutputFile   string
	agyFixDryRun       bool
)

const activeAgyPromptRelativePath = ".lovable/temp/active-agy-pipeline-fix-prompt.txt"

// agyFixPipelineCmd represents the agy fix-pipeline CLI command.
var agyFixPipelineCmd = &cobra.Command{
	Use:     "fix-pipeline [repo]",
	Aliases: []string{"fix", "fp", "pipeline-fix", "fixpipeline"},
	Short:   "Extract latest pipeline error logs and CI/CD fix prompt into clipboard and active temp prompt for Antigravity IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyFixPipelineCLI(args)
	},
}

func init() {
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixDetailed, "detailed", "v", false, "Include verbose passing lines in error logs")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixNoRelease, "no-release", false, "Use CI/CD fix prompt without automated release")
	agyFixPipelineCmd.Flags().StringVarP(&agyFixCustomPrompt, "prompt", "p", "", "Path to custom prompt template")
	agyFixPipelineCmd.Flags().BoolVar(&agyFixNoClipboard, "no-clipboard", false, "Skip writing to system clipboard")
	agyFixPipelineCmd.Flags().StringVar(&agyFixOutputFile, "file", "", "Optional destination file path for prompt payload")
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixDryRun, "dry-run", "d", false, "Preview payload statistics without saving or copying")
}

// RunAgyFixPipelineCLI parses arguments and executes the pipeline fix feed assembly.
func RunAgyFixPipelineCLI(args []string) error {
	repo := resolveTargetRepoArg(args)
	errorReport, hasFailures := cmdpipeline.FetchLatestPipelineErrorReport(repo, agyFixDetailed)
	promptContent, promptSource := loadCicdFixPrompt(agyFixCustomPrompt, agyFixNoRelease)
	payload := AssembleFixPipelinePayload(errorReport, promptContent)

	if agyFixDryRun {
		renderAgyFixDryRun(repo, promptSource, errorReport, promptContent, payload, hasFailures)

		return nil
	}

	writeErr := persistFixPromptPayload(payload, agyFixOutputFile, agyFixNoClipboard)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "persist prompt payload")
	}

	renderAgyFixFeedback(repo, promptSource, errorReport, promptContent, payload, hasFailures)

	return nil
}

func resolveTargetRepoArg(args []string) string {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return args[0]
	}

	return ""
}

// AssembleFixPipelinePayload combines error logs and fix prompt with mandatory two-line gap.
func AssembleFixPipelinePayload(errorLogs, fixPrompt string) string {
	cleanLogs := strings.TrimSpace(errorLogs)
	cleanPrompt := strings.TrimSpace(fixPrompt)
	if len(cleanLogs) == 0 {
		return cleanPrompt
	}
	if len(cleanPrompt) == 0 {
		return cleanLogs
	}

	return cleanLogs + "\n\n" + cleanPrompt
}

func loadCicdFixPrompt(customPath string, isNoRelease bool) (string, string) {
	if len(customPath) > 0 {
		data, err := os.ReadFile(customPath)
		if err == nil {
			return string(data), customPath
		}
	}

	targetPath := selectPromptPath(isNoRelease)
	data, err := os.ReadFile(targetPath)
	if err == nil {
		return string(data), targetPath
	}

	return defaultCicdFixWithReleasePromptFallback, "embedded-fallback"
}

func selectPromptPath(isNoRelease bool) string {
	if isNoRelease {
		return "01-prompts/16-ci-cd/01-ci-cd-fix.md"
	}

	return "01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md"
}

func persistFixPromptPayload(payload, customFile string, skipClipboard bool) error {
	tempPath := filepath.Join(resolveProjectRootDir(), activeAgyPromptRelativePath)
	_ = os.MkdirAll(filepath.Dir(tempPath), 0755)
	writeErr := os.WriteFile(tempPath, []byte(payload), 0644)
	if writeErr != nil {
		return writeErr
	}

	if len(customFile) > 0 {
		_ = os.WriteFile(customFile, []byte(payload), 0644)
	}

	if !skipClipboard {
		_ = clipboard.WriteAll(payload)
	}

	return nil
}

func resolveProjectRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return cwd
}

func renderAgyFixFeedback(repo, promptSource, logs, prompt, payload string, hasFailures bool) {
	statusBanner := constants.ColorGreen + "✔" + constants.ColorReset
	fmt.Printf("\n  %s %sPrepared CI/CD pipeline fix prompt for Antigravity IDE!%s\n\n",
		statusBanner, constants.ColorCyan, constants.ColorReset)
	renderPayloadMetrics(repo, promptSource, logs, prompt, payload, hasFailures)
	renderClipboardNotice()
}

func renderPayloadMetrics(repo, promptSource, logs, prompt, payload string, hasFailures bool) {
	fmt.Printf("    • Target Repo:    %s\n", formatDisplayRepo(repo))
	fmt.Printf("    • Pipeline Logs:  %.1f KB (%s)\n", float64(len(logs))/1024.0, formatFailureStatus(hasFailures))
	fmt.Printf("    • Fix Prompt:     %.1f KB (%s)\n", float64(len(prompt))/1024.0, promptSource)
	fmt.Printf("    • Total Payload:  %.1f KB\n", float64(len(payload))/1024.0)
	fmt.Printf("    • Saved Payload:  %s\n", activeAgyPromptRelativePath)
}

func formatDisplayRepo(repo string) string {
	if len(repo) == 0 {
		return "(auto-detected current repo)"
	}

	return repo
}

func formatFailureStatus(hasFailures bool) string {
	if hasFailures {
		return constants.ColorRed + "failing error logs" + constants.ColorReset
	}

	return constants.ColorGreen + "clean status" + constants.ColorReset
}

func renderClipboardNotice() {
	if !agyFixNoClipboard {
		fmt.Printf("    • Clipboard:      %sCopied to OS clipboard%s ✅\n",
			constants.ColorGreen, constants.ColorReset)
	}
	fmt.Printf("\n  %sReady! Paste into Antigravity IDE chat window (Ctrl+V / Cmd+V) to start the fix loop.%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func renderAgyFixDryRun(repo, promptSource, logs, prompt, payload string, hasFailures bool) {
	fmt.Printf("\n  %s[dry-run] CI/CD Pipeline Fix Feed Preview:%s\n\n",
		constants.ColorYellow, constants.ColorReset)
	renderPayloadMetrics(repo, promptSource, logs, prompt, payload, hasFailures)
	fmt.Println("\n  [dry-run] Skipping clipboard copy and disk persistence.")
}

const defaultCicdFixWithReleasePromptFallback = `# Release-Triggered CI/CD Fix Loop — Workflow (must follow)

Trigger Keywords & Aliases: fix and release, ci release, fix CI/CD and release, cicd fix release

N = 200

### Master Task Checklist
1. First N/2 steps: Review central CI/CD pipelines and local Python runner (03-ai-scripts/06-cicd-local-runner.py).
2. Second N/2 steps: Run local runner (python 03-ai-scripts/06-cicd-local-runner.py --all) to catch all errors via 4-part RCA.
3. Finalize CI/CD: Loop until all registered cases exit 0.
4. Release: Perform version bump, changelog update, git tag, and release orchestration.
`
