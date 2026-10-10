// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdaudit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcommittransfer"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconsole"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddiff"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdip"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmigrate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtemplates"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmduser"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvariable"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautomation"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func isVersionCommand(cmd string) bool {
	return cmd == constants.CmdVersion || cmd == constants.CmdVersionAlias || cmd == "--version" || cmd == "-version" || cmd == "-v"
}

// Run is the main entry point for the CLI.
func Run() {
	cmdconsole.InitConsole()

	if len(os.Args) < 2 {
		printUsageCompact()

		return
	}

	if os.Args[1] == "__complete" || os.Args[1] == "__completeNoDesc" {
		executeCobraCompletion(os.Args[1:])
		return
	}

	// Strip the global `--theme` palette selector first so it is
	// honored even when no subcommand-specific args are present.
	// The writer build below must run AFTER the env var is set but
	// BEFORE any subcommand writes colored output.
	os.Args = append(os.Args[:1], stripThemeFlag(os.Args[1:])...)

	// Strip the global `--glyphs` switch (rich | safe | auto).
	os.Args = append(os.Args[:1], stripGlyphsFlag(os.Args[1:])...)

	// Build the synchronous UI filter writers (see outputwriters.go
	// for the canonical theme-then-glyphs composition order) and
	// carry them in the dispatch context. Byte-faithful commands
	// (cat/view/type) bypass the filter explicitly at the call site:
	// printCatContent writes the decorative banner through the
	// filtered writer and the file bytes through output.Raw().
	//
	// os.Stdout/os.Stderr are never reassigned. Writes are
	// synchronous, so there is no pipe to drain and no flusher to
	// register before os.Exit — the "lost bytes on Windows" class
	// is gone by construction.
	uiOut, uiErr := filteredWriters(termout.Resolve(), glyphs.Resolve())
	output.SetDispatchWriters(uiOut, uiErr)

	// Strip the global `--vscode-sync-disabled` kill switch from argv
	// (and flip the env var) before any subcommand sees its flagset.
	// Done first so even URL-shortcut and alias rewrites operate on
	// already-cleaned args.
	os.Args = append(os.Args[:1], cmdvscode.StripVSCodeSyncDisabledFlag(os.Args[1:])...)
	// Strip the global tag-customization flags too, persisting their
	// values into GITMAP_VSCODE_TAG_{ADD,SKIP,MARKER} so every
	// DetectTagsCustom caller — present and future — picks them up.
	os.Args = append(os.Args[:1], cmdvscode.StripVSCodeTagFlags(os.Args[1:])...)
	// Strip global `--ai` execution recording flag and track in env.
	os.Args = append(os.Args[:1], cmdai.StripAiFlag(os.Args[1:])...)
	if len(os.Args) < 2 {
		printUsageCompact()

		return
	}

	// Skip migration for commands that must produce clean stdout
	cmd := os.Args[1]
	if isVersionCommand(cmd) == false {
		cmdmigrate.MigrateLegacyDirs()
	}

	// URL shortcut: `gitmap <git-url> [<url2> ...]` (and variants with
	// leading flags like `gitmap --verbose <url>`) is rewritten to
	// `gitmap clone <args...>` so users don't have to remember the
	// subcommand for the most common operation. We trigger when any
	// positional arg looks like an HTTPS / SSH git URL or a comma-list
	// containing one — covering all the forms users actually type:
	//
	//   gitmap https://...                       → gitmap clone https://...
	//   gitmap https://a,https://b,https://c     → gitmap clone https://a,https://b,https://c
	//   gitmap https://a, https://b https://c    → gitmap clone https://a, https://b https://c
	//   gitmap --verbose https://...             → gitmap clone --verbose https://...
	if isCloneRewriteRequired(os.Args[1:]) {
		os.Args = append([]string{os.Args[0], constants.CmdClone}, os.Args[1:]...)
	}

	// Multi-word command rewrites:
	//   gitmap full summary [N]    → gitmap full-summary [N]
	//   gitmap full status [N]     → gitmap full-status [N]
	//   gitmap full summary+pe     → gitmap full-summary+pe
	//   gitmap full status+pe      → gitmap full-status+pe
	if len(os.Args) >= 3 && strings.EqualFold(os.Args[1], "full") {
		sub := strings.ToLower(os.Args[2])
		switch sub {
		case "summary":
			os.Args = append([]string{os.Args[0], constants.CmdFullSummary}, os.Args[3:]...)
		case "status":
			os.Args = append([]string{os.Args[0], constants.CmdFullStatus}, os.Args[3:]...)
		case "summary+pe":
			os.Args = append([]string{os.Args[0], constants.CmdFullSummaryPE}, os.Args[3:]...)
		case "status+pe":
			os.Args = append([]string{os.Args[0], constants.CmdFullStatusPE}, os.Args[3:]...)
		}
	}

	os.Args = applyAliasContextIfPresent(os.Args[1], os.Args)

	runDispatch(os.Args[1])
}

func applyAliasContextIfPresent(command string, args []string) []string {
	aliasName, cleaned := extractAliasFlag(args[2:])
	if len(aliasName) == 0 {
		return args
	}

	if err := resolveAliasContext(aliasName); err != nil {
		handleGlobalError(command, err)
	}

	return append(args[:2], cleaned...)
}

// runDispatch is the single entry point every CLI invocation flows
// through. The dispatch writers (built in Run) are already installed
// in the output package; writes are synchronous, so no pipe drainers
// remain — the Windows "last stdout line is lost" class (the
// version-mismatch smoke failure fixed in v6.74.0) is gone by
// construction. Any new dispatch surface MUST call runDispatch
// rather than dispatch directly.
func runDispatch(command string) {
	if tryInterceptCommandHelp(command, os.Args[2:]) {
		return
	}
	dispatch(command)
}

// isPipelineErrorsAlias reports whether the command routes to the pipeline
// errors inspector (pe/te family), which owns its own rich help printer.
func isPipelineErrorsAlias(command string) bool {
	switch command {
	case "pe", "te", "ee", "pipeline", "pipeline-errors", "pipeline_errors":
		return true
	}

	return false
}

func tryInterceptCommandHelp(command string, args []string) bool {
	if len(args) == 0 {
		return false
	}
	if !IsHelpFlag(args[0]) {
		return false
	}
	// pe/te are pipeline-errors aliases with their own rich help printer.
	// Route them there directly instead of the generic dynamic placeholder,
	// which knows nothing about -t, all, -f, or the format system.
	if isPipelineErrorsAlias(command) {
		cmdpipeline.PrintPipelineErrorsHelp()
		cliexit.Exit(0)

		return true
	}
	// Pilot commands (spec 243.3) render through the HelpDisplay registry
	// before the legacy rich-topic/helptext paths.
	if TryPrintDisplayerHelp(command, args) {
		return true
	}
	if tryRenderRichTopic(command) {
		cliexit.Exit(0)
		return true
	}
	if helpdoc.HasTopic(command) {
		printHelpAndExit(command, args)
		return true
	}
	if RenderDynamicCommandHelp(command) {
		cliexit.Exit(0)
		return true
	}

	return false
}

func handleGlobalError(command string, err error) {
	if err == nil {
		return
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if isAppErr && appErr == nil {
		return
	}

	cfg, cfgErr := config.LoadFromFile(constants.DefaultConfigPath)
	display := "full"
	showStack := true

	if cfgErr == nil {
		display = cfg.ErrorDisplay
		showStack = cfg.ShowStackTrace
	}

	persistLastError(command, err)

	if isAppErr && appErr != nil && appErr.HasSuggestions() {
		RenderErrorSuggestions(output.UIErr(), appErr)
	}

	if isAbortOrReportedError(err) {
		cliexit.HandleError(nil, 1)

		return
	}

	if isAppErr && appErr != nil && appErr.Type == apperror.ErrorTypeValidation {
		msg := getValidationErrorMessage(appErr)
		cliexit.Reportf(command, "validation", "", fmt.Errorf("%s", msg))
		cliexit.HandleError(nil, 1)

		return
	}

	if isAppErr && appErr != nil && appErr.Type == apperror.ErrorTypeNotFound {
		msg := getNotFoundErrorMessage(appErr)
		cliexit.Reportf(command, "not found", "", fmt.Errorf("%s", msg))
		cliexit.HandleError(nil, 1)

		return
	}

	if display == "simple" && isAppErr && appErr != nil {
		cliexit.Reportf(command, "execute", "", fmt.Errorf("%s: %w", appErr.Op, getRootCause(err)))
		cliexit.HandleError(nil, 1)

		return
	}

	cliexit.Reportf(command, "execute", "", err)

	stack := resolveErrorStackTrace(err)
	if stack != "" && showStack {
		fmt.Fprintf(output.UIErr(), "Stack Trace:%s\n", stack)
	}

	cliexit.HandleError(nil, 1)
}

func isAbortOrReportedError(err error) bool {
	if err == nil {
		return false
	}

	appErr, isApp := err.(*apperror.AppError)
	if isApp && appErr != nil {
		return checkRootAppErrorAbort(appErr)
	}

	multi, isMulti := err.(interface{ Unwrap() []error })
	if isMulti && multi != nil {
		return checkRootMultiErrorAbort(multi.Unwrap())
	}

	single, isSingle := err.(interface{ Unwrap() error })
	if isSingle && single != nil {
		return isAbortOrReportedError(single.Unwrap())
	}

	return false
}

func checkRootAppErrorAbort(appErr *apperror.AppError) bool {
	isAbort := appErr.Type == apperror.ErrorTypeAbort || isReportedRootContext(appErr.Ctx)
	if isAbort {
		return true
	}

	return isAbortOrReportedError(appErr.Cause)
}

func checkRootMultiErrorAbort(errs []error) bool {
	for _, e := range errs {
		if isAbortOrReportedError(e) {
			return true
		}
	}

	return false
}

func isReportedRootContext(ctx map[string]any) bool {
	if ctx == nil {
		return false
	}

	val, ok := ctx["reported"]
	if !ok {
		return false
	}

	isReported, isBool := val.(bool)

	return isBool && isReported
}

func resolveErrorStackTrace(err error) string {
	appErr, isAppErr := err.(*apperror.AppError)
	if isAppErr && appErr != nil && appErr.Stack != "" {
		return appErr.Stack
	}

	return apperror.CaptureStackTrace(2)
}

func getValidationErrorMessage(appErr *apperror.AppError) string {
	if appErr == nil {
		return ""
	}

	if appErr.Message != "" {
		return appErr.Message
	}

	return appErr.Op
}

func getNotFoundErrorMessage(appErr *apperror.AppError) string {
	if appErr == nil {
		return "item not found"
	}

	raw := resolveNotFoundRawMessage(appErr)
	prefix := appErr.Op + ": "
	if strings.HasPrefix(raw, prefix) {
		return strings.TrimPrefix(raw, prefix)
	}

	return raw
}

func resolveNotFoundRawMessage(appErr *apperror.AppError) string {
	if appErr.Message != "" {
		return strings.TrimSpace(appErr.Message)
	}
	msg, hasCtx := getContextMessage(appErr)
	if hasCtx {
		return msg
	}
	if appErr.Cause != nil {
		return strings.TrimSpace(appErr.Cause.Error())
	}

	return appErr.Op + " not found"
}

func getContextMessage(appErr *apperror.AppError) (string, bool) {
	if appErr.Ctx == nil {
		return "", false
	}
	ctxMsg, ok := appErr.Ctx["msg"].(string)
	if ok && ctxMsg != "" {
		return strings.TrimSpace(ctxMsg), true
	}

	return "", false
}

func persistLastError(command string, err error) {
	if err == nil {
		return
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if isAppErr && appErr == nil {
		return
	}

	writeLastErrorFile(command, err, appErr)
	persistToErrorsDB(command, err, appErr)
}

func writeLastErrorFile(command string, err error, appErr *apperror.AppError) {
	if mkdirErr := os.MkdirAll(".gitmap", 0755); mkdirErr != nil {
		warnLastErrorPersist(mkdirErr)
		return
	}

	report := map[string]any{
		"command":   command,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"error":     err.Error(),
	}

	if appErr != nil {
		populateAppErrorReport(report, appErr)
	}

	b, marshalErr := json.MarshalIndent(report, "", "  ")
	if marshalErr == nil {
		if writeErr := os.WriteFile(".gitmap/last_error.log", b, 0644); writeErr != nil {
			warnLastErrorPersist(writeErr)
		}
	} else {
		if writeErr := os.WriteFile(".gitmap/last_error.log", []byte(err.Error()), 0644); writeErr != nil {
			warnLastErrorPersist(writeErr)
		}
	}
}

// warnLastErrorPersist surfaces last-error-log persistence failures instead of
// swallowing them. It never obscures the original error being handled.
func warnLastErrorPersist(err error) {
	fmt.Fprintf(os.Stderr, "warning: could not persist last-error log: %v\n", err)
}

func persistToErrorsDB(command string, err error, appErr *apperror.AppError) {
	rec := store.InternalErrorRecord{
		ErrorCode:  "E_CLI_ERROR",
		ErrorType:  "CLI_ERROR",
		Command:    command,
		Message:    err.Error(),
		StackTrace: resolveErrorStackTrace(err),
	}

	if appErr != nil {
		rec.ErrorCode = string(appErr.Code)
		rec.ErrorType = string(appErr.Type)
		rec.ContextJson = resolveAppContextJson(appErr)
	}

	store.LogInternalErrorRecord(rec)
}

func resolveAppContextJson(appErr *apperror.AppError) string {
	if appErr == nil || len(appErr.Ctx) == 0 {
		return "{}"
	}

	ctxBytes, err := json.Marshal(appErr.Ctx)
	if err != nil {
		return "{}"
	}

	return string(ctxBytes)
}

func populateAppErrorReport(report map[string]any, appErr *apperror.AppError) {
	if appErr == nil {
		return
	}

	report["code"] = appErr.Code
	report["type"] = string(appErr.Type)
	report["severity"] = string(appErr.Severity)
	report["op"] = appErr.Op
	report["caller"] = appErr.Caller
	report["creator"] = appErr.Creator
	report["context"] = appErr.Ctx

	if appErr.Stack != "" {
		report["stack"] = appErr.Stack
	}

	if appErr.Cause != nil {
		report["cause"] = appErr.Cause.Error()
	}
}

func getRootCause(err error) error {
	for {
		appErr, isAppErr := err.(*apperror.AppError)
		if !isAppErr || appErr == nil || appErr.Cause == nil {
			return err
		}

		err = appErr.Cause
	}
}

func handleDispatchResult(
	command string,
	found bool,
	err error,
	shouldAudit bool,
	auditID int64,
	auditStart time.Time,
) bool {
	if !found {
		return false
	}

	cmdaudit.FinishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)
	if err != nil {
		handleGlobalError(command, err)
	}

	return true
}

// dispatch routes to the correct subcommand handler with audit tracking.
func dispatch(command string) {
	auditID, auditStart, shouldAudit := cmdaudit.BeginCommandAudit(command, os.Args[2:])

	found, err := cmduser.DispatchUser(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdfoldertree.DispatchFolderTree(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchCore(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdcommittransfer.DispatchCommitTransfer(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchRelease(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchUtility(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchData(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchTooling(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchProjectRepos(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmddiff.DispatchDiff(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmddiff.DispatchCompare(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdmerge.DispatchMoveMerge(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchAdd(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdsync.DispatchCommon(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = dispatchCommons(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdspace.DispatchSpace(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdspec.DispatchSpec(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	found, err = cmdtemplates.DispatchTemplates(command)
	if handleDispatchResult(command, found, err, shouldAudit, auditID, auditStart) {
		return
	}

	if dispatchExtraCommand(command, shouldAudit, auditID, auditStart) {
		return
	}

	if cmdmacro.DispatchMacroDynamic(command, shouldAudit, auditID, auditStart) {
		return
	}

	handleUnknownCommand(command)
}

// isCloneRewriteRequired returns true when the args (excluding argv[0])
// describe a bare-URL invocation that should be redirected to `clone`.
// It accepts URLs in any positional slot — not just the first — so
// invocations with leading flags (e.g. `gitmap --verbose <url>`) and
// PowerShell's silent comma-splitting both work.
func isCloneRewriteRequired(args []string) bool {
	if len(args) == 0 {
		return false
	}

	// Never rewrite if the first token is already a known subcommand.
	if !looksLikeFlag(args[0]) && !looksLikeURLToken(args[0]) {
		return false
	}

	for _, a := range args {
		if looksLikeURLToken(a) {
			return true
		}
	}

	return false
}

// looksLikeFlag reports whether the token starts with "-" or "--".
func looksLikeFlag(s string) bool {
	return len(s) > 1 && s[0] == '-'
}

// looksLikeURLToken reports whether a token (or any comma-split piece
// of it) is shaped like a git URL. Used by both the shortcut and the
// unknown-command hint so they agree on what counts as a URL.
func looksLikeURLToken(s string) bool {
	for _, part := range splitOnComma(s) {
		if isLikelyURL(part) {
			return true
		}
	}

	return false
}

// splitOnComma is a tiny strings.Split wrapper so root.go doesn't
// need its own strings import for this single use; trims whitespace
// around each piece and drops empties.
func splitOnComma(s string) []string {
	out := make([]string, 0, 4)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			out = appendTrimmedPiece(out, s, start, i)
			start = i + 1
		}
	}

	return out
}

func appendTrimmedPiece(out []string, s string, start, i int) []string {
	piece := trimSpaces(s[start:i])
	if len(piece) > 0 {
		return append(out, piece)
	}

	return out
}

// trimSpaces removes ASCII whitespace from both ends without pulling
// the strings package into root.go's import surface.
func trimSpaces(s string) string {
	i, j := 0, len(s)
	for i < j && isSpace(s[i]) {
		i++
	}

	for j > i && isSpace(s[j-1]) {
		j--
	}

	return s[i:j]
}

// isSpace reports whether b is ASCII whitespace.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func dispatchIP(ctx context.Context, args []string, parent *cobra.Command) error {
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "ip":
		// cmdip.IPCmd expects to parse args itself, so let Cobra do its thing.
		cmdip.IPCmd.SetArgs(args[1:])

		return cmdip.IPCmd.ExecuteContext(ctx)
	case "ip-change":
		cmdip.IPChangeCmd.SetArgs(args[1:])

		return cmdip.IPChangeCmd.ExecuteContext(ctx)
	}

	return nil
}

func runSJ(args []string) error {
	if err := dispatchSJ(context.Background(), os.Args[1:], nil); err != nil {
		handleGlobalError("ssh-join", err)
	}

	return nil
}

func dispatchSJ(ctx context.Context, args []string, root *cobra.Command) error {
	cleanArgs := stripLeadingSJCommand(args)
	return cmdssh.RunSSHJoinCLI(cleanArgs)
}

func isSJCommand(arg string) bool {
	return arg == "sj" || arg == "ssh-join" || arg == "ssh-joined" || arg == "ssh-joiner"
}

func stripLeadingSJCommand(args []string) []string {
	if len(args) == 0 {
		return args
	}

	if isSJCommand(args[0]) {
		return args[1:]
	}

	return args
}

func isSJCCommand(arg string) bool {
	return arg == "sjc" || arg == "ssh-join-common" || arg == "ssh-join-c" || arg == "join-common"
}

func stripLeadingSJCCommand(args []string) []string {
	if len(args) == 0 {
		return args
	}

	if isSJCCommand(args[0]) {
		return args[1:]
	}

	return args
}

func dispatchSJC(ctx context.Context, args []string, root *cobra.Command) error {
	cleanArgs := stripLeadingSJCCommand(args)
	return cmdssh.RunSSHJoinCommonCLI(cleanArgs)
}

func runSJC(args []string) error {
	if err := dispatchSJC(context.Background(), os.Args[1:], nil); err != nil {
		handleGlobalError("ssh-join-common", err)
	}

	return nil
}

func dispatchExtraCommand(cmd string, shouldAudit bool, id int64, start time.Time) bool {
	if dispatchSmartTestSubsystem(cmd, shouldAudit, id, start) {
		return true
	}
	if dispatchAgySubsystem(cmd, shouldAudit, id, start) {
		return true
	}
	if dispatchPromptSubsystem(cmd, shouldAudit, id, start) {
		return true
	}
	if dispatchAgentSubsystem(cmd, shouldAudit, id, start) {
		return true
	}

	return dispatchGeneralCommands(cmd, shouldAudit, id, start)
}

func dispatchAgentSubsystem(
	command string,
	shouldAudit bool,
	auditID int64,
	auditStart time.Time,
) bool {
	switch command {
	case constants.CmdAgent, constants.CmdAgentAlias, constants.CmdAgentAlias2:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdagent.DispatchAgent(args)
		}, shouldAudit, auditID, auditStart)

		return true
	default:
		return false
	}
}

func dispatchSmartTestSubsystem(
	command string,
	shouldAudit bool,
	auditID int64,
	auditStart time.Time,
) bool {
	switch command {
	case "run-smart", "run-incremental", "smart", "test", "smart-test":
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return runSmartTestCLI(args)
		}, shouldAudit, auditID, auditStart)

		return true
	default:
		return false
	}
}

func dispatchGeneralCommands(cmd string, shouldAudit bool, id int64, start time.Time) bool {
	switch cmd {
	case constants.CmdAutomation, constants.CmdAutomationAlias, constants.CmdAutomationAumAlias, constants.CmdAutomationPyAlias:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdautomation.DispatchAutomation(args)
		}, shouldAudit, id, start)
		return true
	case constants.CmdAi, constants.CmdAiAlias:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdai.DispatchAi(args)
		}, shouldAudit, id, start)
		return true
	case "ai-analysis", "ai-ana", "aia":
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdai.DispatchAiAnalysis(argsTail())
		}, shouldAudit, id, start)
		return true
	case "lowercase", "lower", "lower-case-fix", "lowercase-fix", "lcf", "lower-case-readme", "lowercase-readme", "readme-lower", "readme-lowercase", "lcr", "lc-fix":
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return runLowerCaseFixCLI(args)
		}, shouldAudit, id, start)
		return true
	case "safe-rm", "rm-safe":
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return RunSafeRmCLI(argsTail())
		}, shouldAudit, id, start)
		return true
	case constants.CmdVar, constants.CmdVarAlias:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdvariable.RunVariableCmd(argsTail())
		}, shouldAudit, id, start)
		return true
	case constants.CmdRestEnable:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return RunRestEnable(argsTail())
		}, shouldAudit, id, start)
		return true
	case constants.CmdMigrate:
		executeAndAudit(func(_ context.Context, args []string, _ *cobra.Command) error {
			return cmdmigrate.RunMigrate(argsTail())
		}, shouldAudit, id, start)
		return true
	default:
		return false
	}
}

func dispatchAgySubsystem(
	command string,
	shouldAudit bool,
	auditID int64,
	auditStart time.Time,
) bool {
	switch command {
	case "ip", "ip-change":
		executeAndAudit(dispatchIP, shouldAudit, auditID, auditStart)

		return true
	case "agy", "ag", "antigravity", "rrr", "rm-rejoin-read", "rejoin-read", "rrpr", "rrbr", "rm-rejoin-pin-read", "pins", "ccko", "cckf", "cache-clear", "clean-cache", "recreate-project", "recreate-proj", "recreate-p", "recreate", "rec", "rcp",
		"backup-running-prompts", "backup-running-prompt", "brp",
		"restore-running-prompts", "restore-running-prompt", "rrp",
		"running-prompts", "running-prompt", "rp-prompts",
		"running-projects", "runningprojects", "rp",
		"fpug", "finish-prompts-until-green",
		"sug", "shutdown-until", "shutdown-until-green",
		"account-switch", "asw", "switch-account", "fast-forward", "ff",
		"lap", "last-active-projects",
		"add-read", "ar", "add-and-read",
		"rwi", "rerun-with-id",
		"rwc", "rerun-with-convid", "rwp", "rerun-with-prompt-id",
		"telegram", "email", "settings":
		cmdagy.MachineCLIRunner = func(a []string) error { return cmdos.RunOSCLI(append([]string{"machine"}, a...)) }
		cmdagy.AliasCLIRunner = func(a []string) error { return cmdos.RunOSCLI(append([]string{"alias"}, a...)) }
		executeAndAudit(cmdagy.DispatchAgy, shouldAudit, auditID, auditStart)

		return true
	case "fix-pipeline", "fixpipeline", "pipeline-fix", "aef", "agy-errors-fix", "fix-agy":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdagy.RunPipelineFixAgyCLI(args)
		}, shouldAudit, auditID, auditStart)

		return true
	case "agm", "ag-manager", "antigravity-manager":
		executeAndAudit(cmdagy.DispatchAgm, shouldAudit, auditID, auditStart)

		return true
	default:
		return false
	}
}

func dispatchPromptSubsystem(
	command string,
	shouldAudit bool,
	auditID int64,
	auditStart time.Time,
) bool {
	switch command {
	case "sj", "ssh-join", "ssh-joined", "ssh-joiner":
		executeAndAudit(dispatchSJ, shouldAudit, auditID, auditStart)

		return true
	case "sjc", "ssh-join-common", "ssh-join-c", "join-common":
		executeAndAudit(dispatchSJC, shouldAudit, auditID, auditStart)

		return true
	case "prompt", "prompts", "pmt":
		executeAndAudit(cmdprompt.DispatchPrompt, shouldAudit, auditID, auditStart)

		return true
	case "prompts-template", "prompt-template", "prompts-templates", "prompt-templates", "pt":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdprompttemplate.RunPromptsTemplateCLI(args[1:])
		}, shouldAudit, auditID, auditStart)

		return true
	case "rerun", "rr":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			if len(args) > 1 && cmdagy.IsRerunHelpRequested(args[1:]) {
				cmdagy.RenderAgyRerunHelp()
				return nil
			}
			return cmdagy.RunRerunTopLevelCLI(args[1:])
		}, shouldAudit, auditID, auditStart)

		return true
	case "rra", "rerun-all":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdagy.RunRerunTopLevelCLI(append([]string{"all"}, args[1:]...))
		}, shouldAudit, auditID, auditStart)

		return true
	case "rrq", "rerun-queue":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdagy.RunRerunTopLevelCLI(append([]string{"queue"}, args[1:]...))
		}, shouldAudit, auditID, auditStart)

		return true
	case "rerun-restart":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdagy.RunRerunTopLevelCLI(append([]string{"--restart"}, args[1:]...))
		}, shouldAudit, auditID, auditStart)

		return true
	case "list-prompts", "listprompts", "lp":
		executeAndAudit(func(ctx context.Context, args []string, root *cobra.Command) error {
			return cmdagy.RunListPromptsTopLevelCLI(args[1:])
		}, shouldAudit, auditID, auditStart)

		return true
	default:
		return false
	}
}

func executeAndAudit(
	fn func(context.Context,
		[]string,
		*cobra.Command,
	) error, shouldAudit bool, auditID int64, auditStart time.Time) {
	err := fn(context.Background(), os.Args[1:], nil)
	if err != nil {
		cliexit.HandleError(err, 1)
	}

	cmdaudit.FinishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)
}
