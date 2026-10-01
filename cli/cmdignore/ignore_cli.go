package cmdignore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IgnoreExportPayload models exported group definitions and repo bindings.
type IgnoreExportPayload struct {
	Groups   map[string]IgnoreGroup          `json:"groups"`
	Bindings []store.IgnoreRepoBindingRecord `json:"bindings,omitempty"`
}

// RunIgnoreCLI routes gitmap ignore / gitmap ig commands.
func RunIgnoreCLI(args []string) error {
	if len(args) == 0 {
		return PrintIgnoreHelp()
	}
	cleanCmd, cleanArgs := resolveSubcommandArgs(args)
	return dispatchIgnoreCommand(cleanCmd, cleanArgs)
}

func resolveSubcommandArgs(args []string) (string, []string) {
	subCmd := strings.ToLower(args[0])
	if isIgnoreFlag(subCmd) {
		return routeLeadingFlag(args)
	}
	return subCmd, args[1:]
}

func isIgnoreFlag(arg string) bool {
	return arg == "--force" || arg == "--interval" || arg == "-i"
}

func routeLeadingFlag(args []string) (string, []string) {
	for i, a := range args {
		if !strings.HasPrefix(a, "-") {
			remaining := append([]string{}, args[:i]...)
			remaining = append(remaining, args[i+1:]...)
			return strings.ToLower(a), remaining
		}
	}
	return "cache", args
}

func dispatchIgnoreCommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case "help", "-h", "--help":
		return PrintIgnoreHelp()
	case "ls", "list":
		return runIgnoreList()
	case "config", "cfg":
		return RunIgnoreConfig(subArgs)
	case "cache":
		return RunIgnoreCache(subArgs)
	default:
		return dispatchRuleAndScan(subCmd, subArgs)
	}
}

func dispatchRuleAndScan(subCmd string, subArgs []string) error {
	switch subCmd {
	case "add":
		return runIgnoreAdd(subArgs)
	case "remove", "rm":
		return runIgnoreRemove(subArgs)
	default:
		return dispatchScanAndEdit(subCmd, subArgs)
	}
}

func dispatchScanAndEdit(subCmd string, subArgs []string) error {
	switch subCmd {
	case "scan":
		return runIgnoreScan(subArgs)
	case "scan-ssh", "ss":
		return runIgnoreScanSSH(subArgs)
	case "edit":
		return runIgnoreEdit()
	case "ui", "app":
		return runIgnoreUI()
	default:
		return dispatchGroupManagement(subCmd, subArgs)
	}
}

func dispatchGroupManagement(subCmd string, subArgs []string) error {
	switch subCmd {
	case "add-group":
		return runAddGroup(subArgs)
	case "remove-group", "rm-grp":
		return runRemoveGroup(subArgs)
	case "set-default-group":
		return runSetDefaultGroup(subArgs)
	case "add-grp-to-default", "agtd":
		return runAddGrpToDefault(subArgs)
	default:
		return dispatchRepoOps(subCmd, subArgs)
	}
}

func dispatchRepoOps(subCmd string, subArgs []string) error {
	switch subCmd {
	case "apply":
		return runApply(subArgs)
	case "connect-group-with-repo", "cgwp":
		return runConnectGroupWithRepo(subArgs)
	default:
		return dispatchTransferOps(subCmd, subArgs)
	}
}

func dispatchTransferOps(subCmd string, subArgs []string) error {
	switch subCmd {
	case "export":
		return runExportGroups(subArgs)
	case "import":
		return runImportGroups(subArgs)
	case "action":
		return RunFixIgnoreAll(subArgs)
	default:
		return apperror.NewSimple("unknown ignore subcommand: "+subCmd, "E1003")
	}
}

func runIgnoreList() error {
	groups, _ := LoadGroups()
	fmt.Printf("\n%s  === GITMAP IGNORE GROUPS (%d groups) ===%s\n",
		constants.ColorCyan, len(groups), constants.ColorReset)
	totalPatterns := printAllGroups(groups)
	fmt.Printf("  Total patterns across all groups: %d\n", totalPatterns)
	printRepositoryBindings()
	return nil
}

func printAllGroups(groups map[string]IgnoreGroup) int {
	total := 0
	for name, grp := range groups {
		defTag := ""
		if grp.IsDefault {
			defTag = fmt.Sprintf(" %s[default]%s", constants.ColorGreen, constants.ColorReset)
		}
		total += len(grp.Patterns)
		fmt.Printf("  • %s%s (%d patterns):\n", name, defTag, len(grp.Patterns))
		printGroupPatterns(grp.Patterns)
	}
	return total
}

func printGroupPatterns(patterns []string) {
	for _, pat := range patterns {
		fmt.Printf("      %s\n", pat)
	}
}

func printRepositoryBindings() {
	bindings, err := LoadRepoBindings()
	if err != nil {
		printEmptyBindings()
		return
	}
	if len(bindings) == 0 {
		printEmptyBindings()
		return
	}
	printActiveBindings(bindings)
}

func printEmptyBindings() {
	fmt.Printf("\n%s  === REPOSITORY BINDINGS (0 active) ===%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}

func printActiveBindings(bindings []store.IgnoreRepoBindingRecord) {
	fmt.Printf("\n%s  === REPOSITORY BINDINGS (%d active) ===%s\n",
		constants.ColorCyan, len(bindings), constants.ColorReset)
	for _, b := range bindings {
		fmt.Printf("  • %s -> group: %s\n", b.RepoPath, b.GroupName)
	}
	fmt.Println()
}

func runIgnoreAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore add <pattern>", "E1004")
	}
	addErr := AddPatternsToDefault(args)
	if addErr != nil {
		return addErr
	}
	fmt.Printf("%s✓ Added %d pattern(s) to default group%s\n",
		constants.ColorGreen, len(args), constants.ColorReset)
	applyToCurrentRepoIfGit(args)
	return nil
}

func applyToCurrentRepoIfGit(patterns []string) {
	cwd, err := os.Getwd()
	if err != nil || !gitignoreagm.IsGitRepository(cwd) {
		return
	}
	writeGitignorePatterns(cwd, patterns)
}

func writeGitignorePatterns(cwd string, patterns []string) {
	ignorePath := filepath.Join(cwd, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data), patterns...)
	if isModified {
		_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
		fmt.Printf("%s✓ Added patterns to %s%s\n", constants.ColorGreen, ignorePath, constants.ColorReset)
	}
}

func runIgnoreRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore remove <pattern>", "E1006")
	}
	cwd, _ := os.Getwd()
	ignorePath := filepath.Join(cwd, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		return apperror.WrapSimple(err, "read .gitignore")
	}
	removeSet := filterImmutableRemoveArgs(args)
	kept := filterRemovedLines(string(data), removeSet)
	return os.WriteFile(ignorePath, []byte(strings.Join(kept, "\n")), 0644)
}

func filterImmutableRemoveArgs(args []string) map[string]bool {
	removeSet := make(map[string]bool)
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		isImmutableGitmap := (trimmed == ".gitmap/" || trimmed == ".gitmap")
		isImmutableBackup := (trimmed == ".gitmap/backup/" || trimmed == ".gitmap/backup")
		if isImmutableGitmap || isImmutableBackup {
			fmt.Printf("%s⚠ Cannot remove immutable rule: %s%s\n", constants.ColorYellow, trimmed, constants.ColorReset)
			continue
		}
		removeSet[trimmed] = true
	}
	return removeSet
}

func filterRemovedLines(content string, removeSet map[string]bool) []string {
	lines := strings.Split(content, "\n")
	var kept []string
	for _, l := range lines {
		isRemoved := removeSet[strings.TrimSpace(l)]
		if isRemoved {
			continue
		}
		kept = append(kept, l)
	}
	return kept
}

func runIgnoreScan(args []string) error {
	records := resolveTargetRepos()
	issues := scanReposForIgnoreIssues(records)
	if len(issues) == 0 {
		fmt.Printf("%s✓ No ignore issues found across %d repositories.%s\n",
			constants.ColorGreen, len(records), constants.ColorReset)
		return nil
	}
	printIgnoreScanIssues(issues)
	return nil
}

func printIgnoreScanIssues(issues []IgnoreScanIssue) {
	fmt.Printf("\n%s  === GITIGNORE SCAN ISSUES (%d found) ===%s\n",
		constants.ColorYellow, len(issues), constants.ColorReset)
	for _, iss := range issues {
		fmt.Printf("  • %-30s | resume=%v | duplicates=%v | missingGitmap=%v\n",
			iss.RepoName, iss.HasResumeTask, iss.HasDuplicate, iss.MissingGitmapDir)
	}
	fmt.Printf("\nRun '%sgitmap fix-ignore-all%s' to automatically remediate.\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func runIgnoreScanSSH(args []string) error {
	return RunFixIgnoresAllSSH(append([]string{"--dry-run"}, args...))
}

func runIgnoreEdit() error {
	cwd, _ := os.Getwd()
	ignorePath := filepath.Join(cwd, ".gitignore")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}
	cmd := exec.Command(editor, ignorePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runIgnoreUI() error {
	fmt.Println("Opening GitMap Repository Ignore UI...")
	cmd := exec.Command("gitmap", "ui")
	return cmd.Start()
}

func runAddGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore add-group <name> [pattern1 pattern2...]", "E1007")
	}
	groupName := args[0]
	patterns := args[1:]
	addErr := AddGroup(groupName, patterns)
	if addErr != nil {
		return addErr
	}
	fmt.Printf("%s✓ Created ignore group '%s' with %d patterns%s\n",
		constants.ColorGreen, groupName, len(patterns), constants.ColorReset)
	return nil
}

func runRemoveGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore remove-group <name>", "E1008")
	}
	groupName := args[0]
	rmErr := RemoveGroup(groupName)
	if rmErr != nil {
		return rmErr
	}
	fmt.Printf("%s✓ Removed ignore group '%s'%s\n", constants.ColorGreen, groupName, constants.ColorReset)
	return nil
}

func runSetDefaultGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore set-default-group <name>", "E1009")
	}
	return SetDefaultGroup(args[0])
}

func runAddGrpToDefault(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore add-grp-to-default <name>", "E1010")
	}
	groupName := args[0]
	err := AddGroupToDefault(groupName)
	if err != nil {
		return err
	}
	fmt.Printf("%s✓ Group '%s' chained after default ignore group%s\n",
		constants.ColorGreen, groupName, constants.ColorReset)
	return nil
}

func runApply(args []string) error {
	if len(args) == 0 {
		cwd, _ := os.Getwd()
		return applyTargetDir(cwd)
	}
	firstArg := args[0]
	groups, _ := LoadGroups()
	cleanFirst := strings.ToLower(strings.TrimSpace(firstArg))
	if _, isGroup := groups[cleanFirst]; isGroup {
		targetDir := resolveApplyTargetDir(args)
		return ApplyGroupToRepo(cleanFirst, targetDir)
	}
	return applyTargetDir(firstArg)
}

func resolveApplyTargetDir(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	cwd, _ := os.Getwd()
	return cwd
}

func applyTargetDir(targetDir string) error {
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		absPath = targetDir
	}
	applyErr := ApplyAllToRepo(absPath)
	if applyErr != nil {
		return applyErr
	}
	fmt.Printf("%s✓ Applied ignore rules to %s%s\n", constants.ColorGreen, absPath, constants.ColorReset)
	return nil
}

func runConnectGroupWithRepo(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("Usage: gitmap ignore connect-group-with-repo <group> <repo-path> [--add-with-default|--awd]", "E1012")
	}
	hasAwd, cleanArgs := extractAwdFlag(args)
	if len(cleanArgs) < 2 {
		return apperror.NewSimple("Usage: gitmap ignore connect-group-with-repo <group> <repo-path> [--add-with-default|--awd]", "E1012")
	}
	groupName := cleanArgs[0]
	targetPaths := strings.Split(cleanArgs[1], ",")
	return connectGroupToPaths(groupName, targetPaths, hasAwd)
}

func extractAwdFlag(args []string) (bool, []string) {
	hasAwd := false
	var clean []string
	for _, a := range args {
		trimmed := strings.ToLower(strings.TrimSpace(a))
		isAwdFlag := (trimmed == "--add-with-default" || trimmed == "--awd" || trimmed == "-awd")
		if isAwdFlag {
			hasAwd = true
			continue
		}
		clean = append(clean, a)
	}
	return hasAwd, clean
}

func connectGroupToPaths(groupName string, paths []string, hasAwd bool) error {
	for _, p := range paths {
		repoPath := strings.TrimSpace(p)
		if err := ConnectGroupWithRepo(groupName, repoPath, hasAwd); err != nil {
			return err
		}
	}
	fmt.Printf("%s✓ Connected group '%s' to %d repositories (with default: %v)%s\n",
		constants.ColorGreen, groupName, len(paths), hasAwd, constants.ColorReset)
	return nil
}

func runExportGroups(args []string) error {
	path := resolveExportPath(args)
	payload := buildExportPayload()
	data, _ := json.MarshalIndent(payload, "", "  ")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return apperror.WrapSimple(err, "write export file")
	}
	fmt.Printf("%s✓ Exported %d groups and %d bindings to %s%s\n",
		constants.ColorGreen, len(payload.Groups), len(payload.Bindings), path, constants.ColorReset)
	return nil
}

func resolveExportPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "gitmap-ignore-groups.json"
}

func buildExportPayload() IgnoreExportPayload {
	groups, _ := LoadGroups()
	bindings, _ := LoadRepoBindings()
	return IgnoreExportPayload{Groups: groups, Bindings: bindings}
}

func runImportGroups(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore import <file.json>", "E1013")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import file")
	}
	payload, parseErr := parseImportData(data)
	if parseErr != nil {
		return parseErr
	}
	return applyImportedPayload(payload, args[0])
}

func parseImportData(data []byte) (IgnoreExportPayload, *apperror.AppError) {
	payload, isExportOk := parseExportPayload(data)
	if isExportOk {
		return payload, nil
	}
	legacyPayload, isLegacyOk := parseLegacyPayload(data)
	if isLegacyOk {
		return legacyPayload, nil
	}
	return payload, apperror.NewSimple("invalid json ignore groups format", "E1015")
}

func parseExportPayload(data []byte) (IgnoreExportPayload, bool) {
	var payload IgnoreExportPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, false
	}
	hasGroups := (len(payload.Groups) > 0)
	return payload, hasGroups
}

func parseLegacyPayload(data []byte) (IgnoreExportPayload, bool) {
	var legacyGroups map[string]IgnoreGroup
	if err := json.Unmarshal(data, &legacyGroups); err != nil {
		return IgnoreExportPayload{}, false
	}
	return IgnoreExportPayload{Groups: legacyGroups}, true
}

func applyImportedPayload(payload IgnoreExportPayload, sourceFile string) error {
	if err := SaveGroups(payload.Groups); err != nil {
		return err
	}
	importBindings(payload.Bindings)
	fmt.Printf("%s✓ Imported %d groups and %d bindings from %s%s\n",
		constants.ColorGreen, len(payload.Groups), len(payload.Bindings), sourceFile, constants.ColorReset)
	return nil
}

func importBindings(bindings []store.IgnoreRepoBindingRecord) {
	for _, b := range bindings {
		_ = saveRepoBindingRecord(b.RepoPath, b.GroupName)
	}
}

const ignoreHelpMenu = `
  GitMap Ignore Management Suite (ig / ignore)
    • gitmap ignore ls                             - List configured ignore groups & bindings
    • gitmap ignore config                         - Show check interval, Split-DB path & stats
    • gitmap ignore config set interval <dur>      - Configure audit check frequency
    • gitmap ignore cache                          - List cached repo audit statuses
    • gitmap ignore cache clear (--force)          - Invalidate all cached ignore records
    • gitmap ignore add <pattern>                  - Add pattern to default group
    • gitmap ignore remove (rm) <pattern>          - Remove pattern from current repo
    • gitmap ignore scan                           - Scan all repos for ignore issues
    • gitmap ignore scan-ssh (ss)                  - Scan fleet nodes for ignore issues
    • gitmap ignore edit                           - Edit current repository .gitignore
    • gitmap ignore add-group <name> [pats...]     - Create named ignore group
    • gitmap ignore remove-group (rm-grp) <name>   - Remove named ignore group
    • gitmap ignore set-default-group <name>       - Set group as default
    • gitmap ignore add-grp-to-default (agtd) <n>  - Chain named group after default group
    • gitmap ignore connect-group-with-repo (cgwp) - Connect group with repo (--awd)
    • gitmap ignore apply [path|group]             - Apply ignore patterns to directory/repo
    • gitmap ignore export [file]                  - Export ignore groups to JSON
    • gitmap ignore import <file>                  - Import ignore groups from JSON
    • gitmap fix-ignore-all (fia) [-y]             - Fix & sanitize ignore across all repos
    • gitmap fix-ignores-all-ssh (fias) [-y]       - Fix & sanitize ignore across SSH fleet
`

// PrintIgnoreHelp renders the usage menu for gitmap ignore commands.
func PrintIgnoreHelp() error {
	fmt.Print(ignoreHelpMenu)
	return nil
}
