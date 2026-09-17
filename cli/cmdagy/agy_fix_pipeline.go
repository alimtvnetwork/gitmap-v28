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
	agyFixForce        bool
)

const activeAgyPromptRelativePath = ".lovable/temp/active-agy-pipeline-fix-prompt.txt"

// agyFixPipelineCmd represents the agy fix-pipeline CLI command.
var agyFixPipelineCmd = &cobra.Command{
	Use:     "fix-pipeline [repo]",
	Aliases: []string{"fix", "fp", "pipeline-fix", "fixpipeline", "aef", "agy-errors-fix"},
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
	agyFixPipelineCmd.Flags().BoolVarP(&agyFixForce, "force", "f", false, "Force resending even if previously sent")
}

func resolveTargetRepoArg(args []string) string {
	opts := parseAgyFixArgs(args)

	return opts.Repo
}

// RunAgyFixPipelineCLI parses arguments and executes the pipeline fix feed assembly.
func RunAgyFixPipelineCLI(args []string) error {
	return RunPipelineFixAgyCLI(args)
}

func checkAgyFixDuplicate(opts AgyFixOptions, payload cmdpipeline.PipelineErrorLogsPayload, errorReport string) (bool, string, string, string) {
	sig, errHash := ComputeErrorSignature(payload.Repo, payload.RunId, payload.Sha, errorReport)
	storePath := sentAgyErrorsStorePath()
	store := LoadSentAgyErrorsStore(storePath)
	isDuplicate, existingRecord := CheckSentErrorDuplicate(sig, store, opts.IsForce)
	if isDuplicate {
		renderDuplicateNotice(payload.Repo, payload.RunId, payload.Sha, existingRecord)

		return true, storePath, sig, errHash
	}

	return false, storePath, sig, errHash
}

// RunPipelineFixAgyCLI is the unified entrypoint for pipeline fix errors agy / aef.
func RunPipelineFixAgyCLI(args []string) error {
	opts := parseAgyFixArgs(args)
	payload, errorReport, hasFailures := cmdpipeline.FetchPipelineErrorReportWithMeta(opts.Repo, opts.IsDetailed)
	isDup, storePath, sig, errHash := checkAgyFixDuplicate(opts, payload, errorReport)
	if isDup {
		return nil
	}

	params := AgyFixDispatchParams{
		Opts: opts, StorePath: storePath, Sig: sig, ErrHash: errHash,
		Payload: payload, ErrorReport: errorReport, HasFailures: hasFailures,
	}

	return dispatchAgyFixPrepared(params)
}

func assemblePrimaryAndFollowup(opts AgyFixOptions, payload cmdpipeline.PipelineErrorLogsPayload, errorReport string) (string, string, string) {
	gitLog := ExtractGitLog("", 5)
	promptContent, promptSource := LoadCanonicalRcaPrompt(opts.CustomPrompt, opts.IsNoRelease)
	primary := AssembleRcaFixPayload(payload.Repo, payload.RunId, payload.Sha, gitLog, errorReport, promptContent)
	followup := BuildVerificationFollowupPrompt(payload.Repo, payload.RunId, payload.Sha)

	return primary, followup, promptSource
}

func dispatchAgyFixPrepared(p AgyFixDispatchParams) error {
	primary, followup, promptSource := assemblePrimaryAndFollowup(p.Opts, p.Payload, p.ErrorReport)
	if p.Opts.IsDryRun {
		renderAgyFixDryRun(p.Payload.Repo, promptSource, p.ErrorReport, primary, primary, p.HasFailures)

		return nil
	}

	return executeFixPayloadDispatch(p, promptSource, primary, followup)
}

func finalizeFixFeedback(repo, promptSource, errorReport, promptContent, primaryPayload string, hasFailures bool) {
	renderAgyFixFeedback(repo, promptSource, errorReport, promptContent, primaryPayload, hasFailures)
	renderQueuedVerificationNotice()
}

func persistAndRecordAgyFix(p AgyFixDispatchParams, primary, followup string) error {
	if writeErr := persistFixPromptPayload(primary, p.Opts.OutputFile, p.Opts.IsNoClipboard); writeErr != nil {
		return apperror.WrapSimple(writeErr, "persist prompt payload")
	}

	_ = StageVerificationFollowupPrompt(primary, followup)
	_ = RecordSentErrorSignature(p.StorePath, p.Sig, p.Payload.Repo, p.Payload.RunId, p.Payload.Sha, p.ErrHash)

	return nil
}

func executeFixPayloadDispatch(p AgyFixDispatchParams, promptSource, primary, followup string) error {
	if err := persistAndRecordAgyFix(p, primary, followup); err != nil {
		return err
	}

	finalizeFixFeedback(p.Payload.Repo, promptSource, p.ErrorReport, primary, primary, p.HasFailures)

	return nil
}

func defaultAgyFixOptions() AgyFixOptions {
	return AgyFixOptions{
		IsDetailed:    agyFixDetailed,
		IsNoRelease:   agyFixNoRelease,
		CustomPrompt:  agyFixCustomPrompt,
		IsNoClipboard: agyFixNoClipboard,
		OutputFile:    agyFixOutputFile,
		IsDryRun:      agyFixDryRun,
		IsForce:       agyFixForce,
	}
}

func parseAgyFixArgs(args []string) AgyFixOptions {
	opts := defaultAgyFixOptions()
	for i := 0; i < len(args); i++ {
		parseSingleArg(args, &i, &opts)
	}

	return opts
}

func parseBehaviorToggle(arg string, opts *AgyFixOptions) bool {
	if isAgyForceFlag(arg) {
		opts.IsForce = true
		return true
	}
	if isAgyDetailedFlag(arg) {
		opts.IsDetailed = true
		return true
	}

	return false
}

func parseOutputToggle(arg string, opts *AgyFixOptions) bool {
	switch {
	case arg == "--no-release":
		opts.IsNoRelease = true
		return true
	case arg == "--no-clipboard":
		opts.IsNoClipboard = true
		return true
	case isAgyDryRunFlag(arg):
		opts.IsDryRun = true
		return true
	default:
		return false
	}
}

func parseToggleArg(arg string, opts *AgyFixOptions) bool {
	return parseBehaviorToggle(arg, opts) || parseOutputToggle(arg, opts)
}

func parseSingleArg(args []string, idx *int, opts *AgyFixOptions) {
	arg := args[*idx]
	if parseToggleArg(arg, opts) {
		return
	}

	parseArgWithParam(args, idx, opts, arg)
}

func isAgyForceFlag(arg string) bool {
	return arg == "--force" || arg == "-f"
}

func isAgyDetailedFlag(arg string) bool {
	return arg == "--detailed" || arg == "-v"
}

func isAgyDryRunFlag(arg string) bool {
	return arg == "--dry-run" || arg == "-d"
}

func parseParamFlag(args []string, idx *int, opts *AgyFixOptions, arg string) bool {
	hasNext := *idx+1 < len(args)
	if (arg == "--prompt" || arg == "-p") && hasNext {
		opts.CustomPrompt = args[*idx+1]
		*idx++
		return true
	}
	if arg == "--file" && hasNext {
		opts.OutputFile = args[*idx+1]
		*idx++
		return true
	}

	return false
}

func parseArgWithParam(args []string, idx *int, opts *AgyFixOptions, arg string) {
	if parseParamFlag(args, idx, opts, arg) || strings.HasPrefix(arg, "-") {
		return
	}

	if !isSubcommandKeyword(strings.ToLower(arg)) && len(opts.Repo) == 0 {
		opts.Repo = arg
	}
}

func isSubcommandKeyword(word string) bool {
	switch word {
	case "fix", "errors", "error", "err", "agy", "aef",
		"pipeline", "pipeline-fix", "fix-agy", "agy-errors-fix",
		"fix-pipeline", "fixpipeline", "fp":
		return true
	}

	return false
}

func renderDuplicateNotice(repo string, runID uint64, sha string, rec *SentAgyErrorRecord) {
	desc := formatErrorRunDesc(repo, runID, sha)
	fmt.Printf("\n  %s⚠ Pipeline errors for %s have already been sent to Antigravity!%s\n",
		constants.ColorYellow, desc, constants.ColorReset)
	if rec != nil && len(rec.SentAt) > 0 {
		fmt.Printf("    Previously sent at: %s (dispatch count: %d)\n", rec.SentAt, rec.SentCount)
	}
	fmt.Printf("    %sDo you want to send again? Use --force or -f to send again.%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func formatErrorRunDesc(repo string, runID uint64, sha string) string {
	if runID > 0 && len(sha) > 0 {
		return fmt.Sprintf("run #%d (commit %s)", runID, truncateSHA(sha))
	}
	if runID > 0 {
		return fmt.Sprintf("run #%d", runID)
	}
	if len(sha) > 0 {
		return fmt.Sprintf("commit %s", truncateSHA(sha))
	}

	return formatDisplayRepo(repo)
}

func truncateSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}

	return sha
}

func renderQueuedVerificationNotice() {
	fmt.Printf("  %s✓ Follow-up Verification Prompt Queued: %s (\"Is it fixed?\")%s\n",
		constants.ColorGreen, queuedAgyPromptRelativePath, constants.ColorReset)
	fmt.Printf("    • Queue Ledger:   %s\n", agyPromptQueueRelativePath)
	if agyPath, hasAgy := resolveAntigravityBinary(); hasAgy {
		fmt.Printf("    • Antigravity CLI: Detected at %s (run: agy -c)\n", agyPath)
	}
	fmt.Println()
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

func readCustomPrompt(customPath string) (string, string, bool) {
	if len(customPath) == 0 {
		return "", "", false
	}
	data, err := os.ReadFile(customPath)
	if err != nil {
		return "", "", false
	}

	return string(data), customPath, true
}

func saveCustomFileIfRequested(customFile, payload string) {
	if len(customFile) > 0 {
		_ = os.WriteFile(customFile, []byte(payload), 0644)
	}
}

func copyClipboardIfNotSkipped(payload string, skipClipboard bool) {
	if !skipClipboard {
		_ = clipboard.WriteAll(payload)
	}
}

func persistFixPromptPayload(payload, customFile string, skipClipboard bool) error {
	tempPath := filepath.Join(resolveProjectRootDir(), activeAgyPromptRelativePath)
	_ = os.MkdirAll(filepath.Dir(tempPath), 0755)
	if writeErr := os.WriteFile(tempPath, []byte(payload), 0644); writeErr != nil {
		return writeErr
	}

	saveCustomFileIfRequested(customFile, payload)
	copyClipboardIfNotSkipped(payload, skipClipboard)

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
