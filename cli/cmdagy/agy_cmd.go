// Package cmd — agy_cmd.go is the root command for Antigravity workspace management.
package cmdagy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// AgyCmd is the root agy command
var AgyCmd = &cobra.Command{
	Use:   "agy",
	Short: "Antigravity CLI Management",
}

// DispatchAgy routes CLI arguments to agy commands.
func DispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	args = stripAgyPrefix(args)
	if len(args) > 0 && isAgyOpenPathArg(args[0]) {
		return RunAgyOpen(args[0])
	}
	args = normalizeAgyArgs(args)
	if err, isHandled := tryDispatchAgyShortcut(args); isHandled {
		return err
	}
	AgyCmd.SetArgs(args)

	return AgyCmd.ExecuteContext(ctx)
}

func tryDispatchAgyShortcut(args []string) (error, bool) {
	if len(args) > 0 && isAgyFindDuplicatesArg(args[0]) {
		return RunFindDuplicates(), true
	}
	if isAgyLsEmptyConvsArg(args) {
		return runAgyLsEmptyConvs(args[1:]), true
	}
	if isDirectAgyPromptInjection(args) {
		return runAgyPrompt(args[1:]), true
	}

	return nil, false
}

func isDirectAgyPromptInjection(args []string) bool {
	if len(args) < 2 {
		return false
	}
	isPrompt := args[0] == "prompt" || args[0] == "pr"
	if isPrompt == false {
		return false
	}
	sub := strings.ToLower(args[1])
	isSubcmd := sub == "read" || sub == "show" || sub == "view" || sub == "cat" || sub == "ls" || sub == "list" || sub == "-h" || sub == "--help"

	return isSubcmd == false
}

func stripAgyPrefix(args []string) []string {
	if len(args) > 0 && (args[0] == "agy" || args[0] == "ag" || args[0] == "antigravity") {
		return args[1:]
	}

	return args
}

func normalizeAgyArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	if isCompoundAgyFix(args) {
		return rewriteCompoundAgyFix(args)
	}

	args[0] = normalizeAgySubcommand(args[0])

	return args
}

func isCompoundAgyFix(args []string) bool {
	if len(args) < 2 {
		return false
	}

	first := strings.ToLower(args[0])
	second := strings.ToLower(args[1])
	if first == "errors" || first == "error" || first == "err" {
		return second == "fix" || second == "aef"
	}
	if first == "fix" {
		return second == "errors" || second == "error" || second == "pipeline" || second == "agy"
	}

	return false
}

func rewriteCompoundAgyFix(args []string) []string {
	return append([]string{"fix-pipeline"}, args[2:]...)
}

func isAgyOpenPathArg(arg string) bool {
	if arg == "." || arg == ".." || strings.HasPrefix(arg, "./") || strings.HasPrefix(arg, ".\\") {
		return true
	}
	if strings.HasPrefix(arg, "/") || strings.HasPrefix(arg, "\\") || strings.Contains(arg, string(filepath.Separator)) {
		return true
	}
	if info, err := os.Stat(arg); err == nil && info.IsDir() {
		return true
	}

	return false
}

func normalizeAgySubcommand(sub string) string {
	low := strings.ToLower(sub)
	if match := normalizeProjectSubcommands(low); len(match) > 0 {
		return match
	}
	if match := normalizeMaintenanceSubcommands(low); len(match) > 0 {
		return match
	}
	if match := normalizeWorkflowSubcommands(low); len(match) > 0 {
		return match
	}

	return sub
}

func normalizeProjectSubcommands(low string) string {
	if isCureDupsAlias(low) {
		return "optimize-projects"
	}
	if isRemoveMissingAlias(low) {
		return "remove-missing-projects"
	}
	if isReadMemoryAlias(low) {
		return "all-projects-read-memory-prompt"
	}
	if isRprpAlias(low) {
		return "read-all-projects-with-read-prompts"
	}
	if low == "reconcile" || low == "recon" || low == "reconcile-projects" {
		return "reconcile"
	}
	if low == "find-duplicate-projects" || low == "fdp" {
		return "find-duplicate-projects"
	}
	if low == "pin-projects" || low == "pin-project" || low == "pinned-projects" || low == "pinned" || low == "pins" {
		return "pin-projects"
	}

	return ""
}

func normalizeMaintenanceSubcommands(low string) string {
	if low == "install" || low == "in" || low == "i" {
		return "install"
	}
	if low == "clean-cache" || low == "cleancache" || low == "clean_cache" || low == "cc" {
		return "clean-cache"
	}
	if low == "ping" || low == "check" {
		return "ping"
	}
	if isFixPipelineAlias(low) {
		return "fix-pipeline"
	}

	return ""
}

func normalizeWorkflowSubcommands(low string) string {
	if low == "list-prompts" || low == "listprompts" || low == "lp" || low == "list-prompt" {
		return "list-prompts"
	}
	if low == "rerun" || low == "replay" || low == "rr" {
		return "rerun"
	}
	if low == "queue" || low == "q" {
		return "queue"
	}
	if low == "prompt" || low == "pr" {
		return "prompt"
	}

	return ""
}

func isFixPipelineAlias(low string) bool {
	return low == "fix-pipeline" || low == "fix" || low == "fp" ||
		low == "pipeline-fix" || low == "fixpipeline" ||
		low == "aef" || low == "agy-errors-fix" || low == "errors-fix" ||
		low == "fix-errors" || low == "fix-agy"
}

func isCureDupsAlias(low string) bool {
	return low == "--repeat-fix" || low == "-r" ||
		low == "cure-duplicate-projects" || low == "cdp" ||
		low == "cure-duplicates" || low == "cure-duplicate"
}

func isRemoveMissingAlias(low string) bool {
	return low == "remove-misisng-projects" || low == "remove-missing-projects" ||
		low == "rm-missing-projects" || low == "rm-missing" || low == "clean-missing"
}

func isReadMemoryAlias(low string) bool {
	return low == "all-projects-read-memory-prompt" || low == "aprmp" ||
		low == "read-memory-all" || low == "rm-all-prompt"
}

func isRprpAlias(low string) bool {
	return low == "read-all-projects-with-read-prompts" || low == "rprp" ||
		low == "rapwrp" || low == "read-all-with-prompts"
}

func isAgyFindDuplicatesArg(sub string) bool {
	low := strings.ToLower(sub)

	return low == "find-duplicates" || low == "duplicates" || low == "dups" || low == "find-dups"
}

func isAgyLsEmptyConvsArg(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if args[0] == "show-projects-with-empty-conversations" || args[0] == "show-proects-with-empty-conversations" {
		return true
	}

	if args[0] == "ls" && len(args) > 1 {
		sub := strings.ToLower(args[1])

		return sub == "show-projects-with-empty-conversations" ||
			sub == "show-proects-with-empty-conversations" ||
			sub == "empty-conversations" ||
			sub == "--empty-conversations" ||
			sub == "empty-convs"
	}

	return false
}

func init() {
	registerAgyBaseCommands()
	registerAgyProjectCommands()
	registerAgyUtilityCommands()
	initPlugins()
	initAgyGroup()
	initAgySettings()
	initAgyPinProjects()
	initAgyQueueCommands()
	initAgyPromptAndStatusCommands()
	AgyCmd.SetHelpFunc(renderAgyHelp)
}

func registerAgyBaseCommands() {
	AgyCmd.AddCommand(agyAddCmd)
	AgyCmd.AddCommand(agyRmCmd)
	AgyCmd.AddCommand(agyLsCmd)
	AgyCmd.AddCommand(agyStatusCmd)
	AgyCmd.AddCommand(agyPingCmd)
	AgyCmd.AddCommand(agyOptimizeCmd)
	AgyCmd.AddCommand(agyScanCmd)
	AgyCmd.AddCommand(agyStatsCmd)
	AgyCmd.AddCommand(agyUpdateCmd)
	AgyCmd.AddCommand(agyClearCmd)
	AgyCmd.AddCommand(agyOpenCmd)
	AgyCmd.AddCommand(agyPromptCmd)
}

func registerAgyProjectCommands() {
	AgyCmd.AddCommand(agyRwCmd)
	AgyCmd.AddCommand(agySyncCmd)
	AgyCmd.AddCommand(agyPapCmd)
	AgyCmd.AddCommand(agyExportCmd)
	AgyCmd.AddCommand(agyImportCmd)
	AgyCmd.AddCommand(agyPluginsCmd)
	AgyCmd.AddCommand(agyRemoveEmptyConvsCmd)
	AgyCmd.AddCommand(agyFindDupsCmd)
	AgyCmd.AddCommand(agyRemoveMissingCmd)
	AgyCmd.AddCommand(agyReconcileCmd)
	AgyCmd.AddCommand(agyAllProjectsReadMemoryCmd)
}

func registerAgyUtilityCommands() {
	AgyCmd.AddCommand(agyReadAllProjectsWithReadPromptsCmd)
	AgyCmd.AddCommand(agyGroupCmd)
	AgyCmd.AddCommand(agyUndoCmd)
	AgyCmd.AddCommand(agyRedoCmd)
	AgyCmd.AddCommand(agySettingsCmd)
	AgyCmd.AddCommand(agyPinProjectsCmd)
	AgyCmd.AddCommand(agyCleanCacheCmd)
	AgyCmd.AddCommand(agyFixPipelineCmd)
	AgyCmd.AddCommand(agyRerunCmd)
	AgyCmd.AddCommand(agyListPromptsCmd)
}

var agyQueueCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"q"},
	Short:   "Manage Antigravity prompt queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueStatus()
	},
}

var agyQueueStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current prompt queue status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueStatus()
	},
}

var agyQueueLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all queued prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueLs()
	},
}

var agyQueueClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all prompts from the queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueClear()
	},
}

var agyQueuePopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop and stage the next queued prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueuePop()
	},
}

func initAgyQueueCommands() {
	agyQueueCmd.AddCommand(agyQueueStatusCmd)
	agyQueueCmd.AddCommand(agyQueueLsCmd)
	agyQueueCmd.AddCommand(agyQueueClearCmd)
	agyQueueCmd.AddCommand(agyQueuePopCmd)
	AgyCmd.AddCommand(agyQueueCmd)
}

func initAgyPromptAndStatusCommands() {
	initAgyPromptSubcommands()
	agyPromptCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runAgyPrompt(args)
	}
	agyStatusCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runAgyStatusWithQueue()
	}
}

func runAgyQueueStatus() error {
	q, err := GetQueueStatus()
	if err != nil {
		return apperror.WrapSimple(err, "get queue status")
	}

	fmt.Printf("%s● Antigravity Prompt Queue Status%s\n", constants.ColorCyan, constants.ColorReset)
	renderQueueActiveItem(q.Active)
	fmt.Printf("  Queued Prompts: %d\n", len(q.Queued))
	if len(q.UpdatedAt) > 0 {
		fmt.Printf("  Last Updated:   %s\n", q.UpdatedAt)
	}

	return nil
}

func renderQueueActiveItem(item *AgyPromptQueueEntry) {
	if item == nil {
		fmt.Println("  Active Prompt:  None")

		return
	}

	fmt.Printf("  Active Prompt:  #%d [%s] %s (%s)\n", item.ID, item.Type, item.Title, item.Status)
}

func runAgyQueueLs() error {
	q, err := GetQueueStatus()
	if err != nil {
		return apperror.WrapSimple(err, "list queue")
	}
	if len(q.Queued) == 0 {
		fmt.Println("No prompts currently queued.")

		return nil
	}

	fmt.Printf("%sQueued Prompts (%d):%s\n", constants.ColorCyan, len(q.Queued), constants.ColorReset)
	for _, item := range q.Queued {
		fmt.Printf("  #%d [%s] %s (Created: %s)\n", item.ID, item.Type, item.Title, item.CreatedAt)
	}

	return nil
}

func runAgyQueueClear() error {
	if err := ClearPromptQueue(); err != nil {
		return apperror.WrapSimple(err, "clear prompt queue")
	}
	fmt.Printf("%s✔ Cleared Antigravity prompt queue%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func runAgyQueuePop() error {
	popped, err := PopNextQueuedPrompt()
	if err != nil {
		return apperror.WrapSimple(err, "pop queued prompt")
	}
	if popped == nil {
		fmt.Println("Queue is empty. No prompts to pop.")

		return nil
	}
	fmt.Printf("%s✔ Popped prompt #%d [%s] %s%s\n", constants.ColorGreen, popped.ID, popped.Type, popped.Title, constants.ColorReset)
	fmt.Printf("  Staged to: %s\n", resolveActiveAgyPromptPath())
	fmt.Println("  Copied prompt content to clipboard.")

	return nil
}

func runAgyPrompt(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("requires prompt text: gitmap agy prompt [slug/id/path] [prompt-text]", "E9010")
	}

	target, promptText := parseAgyPromptArgs(args)
	if len(promptText) == 0 {
		return apperror.NewSimple("prompt text cannot be empty", "E9011")
	}

	res := InjectAgyPrompt(target, promptText, "User Prompt", "user_prompt", false)
	renderInjectionFeedback(res, false)

	return nil
}

func parseAgyPromptArgs(args []string) (string, string) {
	if len(args) == 1 {
		return "", args[0]
	}

	if isAgyOpenPathArg(args[0]) || isKnownProjectSlug(args[0]) {
		return args[0], strings.Join(args[1:], " ")
	}

	return "", strings.Join(args, " ")
}

func isKnownProjectSlug(slug string) bool {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return false
	}
	pPath := filepath.Join(dirPath, slug+".json")

	return checkFileExists(pPath)
}

func runAgyStatusWithQueue() error {
	cwd, _ := os.Getwd()
	ideProc := DetectRunningAntigravityIDE()
	convState := DetectConversationExecutionState(cwd)
	q, _ := GetQueueStatus()

	renderStatusHeader(ideProc, convState, len(q.Queued))

	return runAgyLs()
}

func renderStatusHeader(ideProc result.Result[AgyProcessInfo], state AgyConversationExecutionState, queueCount int) {
	fmt.Printf("%s● Antigravity System Status%s\n", constants.ColorCyan, constants.ColorReset)
	renderIDEProcessStatus(ideProc)
	renderConversationStatus(state)
	fmt.Printf("  Prompt Queue:  %d pending prompt(s)\n\n", queueCount)
}

func renderIDEProcessStatus(ideProc result.Result[AgyProcessInfo]) {
	if ideProc.IsSuccess() {
		fmt.Printf("  IDE Process:   %sRUNNING%s (PID: %d)\n", constants.ColorGreen, constants.ColorReset, ideProc.Value.PID)

		return
	}
	fmt.Printf("  IDE Process:   %sNOT RUNNING%s\n", constants.ColorYellow, constants.ColorReset)
}

func renderConversationStatus(state AgyConversationExecutionState) {
	color := constants.ColorGreen
	if state.Status == AgyConvStatusRunning {
		color = constants.ColorYellow
	}
	if len(state.ConvID) > 0 {
		fmt.Printf("  Conversation:  %s%s%s (ID: %s)\n", color, strings.ToUpper(string(state.Status)), constants.ColorReset, state.ConvID)

		return
	}
	fmt.Printf("  Conversation:  %s%s%s\n", color, strings.ToUpper(string(state.Status)), constants.ColorReset)
}

func getProjectsDirPath() (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}

	return filepath.Join(homeDir, ".gemini", "config", "projects"), nil
}

func ensureDirExists(dirPath string) bool {
	mkErr := os.MkdirAll(dirPath, 0755)

	return mkErr == nil
}
