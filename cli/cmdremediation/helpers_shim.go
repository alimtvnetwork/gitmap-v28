package cmdremediation

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)
	la, lb := len(ar), len(br)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	dp := initLevenshteinDP(lb)
	fillLevenshteinDP(dp, ar, br, la, lb)
	return dp[lb]
}

func RunInteractiveRemediation(items []RemediationItem) error {
	reader := bufio.NewReader(os.Stdin)
	var applyAllAction string
	for i := 0; i < len(items); i++ {
		nextAction, shouldQuit := processPromptStep(reader, i+1, len(items), &items[i], applyAllAction)
		if shouldQuit {
			return nil
		}

		applyAllAction = nextAction
	}

	return nil
}

func initLevenshteinDP(lb int) []int {
	dp := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		dp[j] = j
	}
	return dp
}

func fillLevenshteinDP(dp []int, ar, br []rune, la, lb int) {
	for i := 1; i <= la; i++ {
		prev := dp[0]
		dp[0] = i
		updateLevenshteinRow(dp, ar[i-1], br, lb, &prev)
	}
}

func updateLevenshteinRow(dp []int, charA rune, br []rune, lb int, prev *int) {
	for j := 1; j <= lb; j++ {
		temp := dp[j]
		cost := 0
		if charA != br[j-1] {
			cost = 1
		}
		dp[j] = min3(dp[j]+1, dp[j-1]+1, *prev+cost)
		*prev = temp
	}
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

var promptChoiceMap = map[string]string{
	"1":         "stash",
	"stash":     "stash",
	"2":         "wip",
	"wip":       "wip",
	"3":         "discard",
	"discard":   "discard",
	"s":         "skip",
	"skip":      "skip",
	"a":         "all-stash",
	"all":       "all-stash",
	"all-stash": "all-stash",
}

func executeAllAction(item *RemediationItem, action string) {
	idx := ParseRecipeIndex(action, item.Recipes)
	isInvalidIndex := idx < 0 || idx >= len(item.Recipes)
	if isInvalidIndex {
		return
	}

	err := ExecuteFixRecipe(item, item.Recipes[idx])
	if err != nil {
		fmt.Printf("Warning: fix action failed on %s: %v\n", item.RepoName, err)
	}
}

func handlePromptAction(item *RemediationItem, action string) string {
	if action == "skip" {
		return ""
	}

	if strings.HasPrefix(action, "all-") {
		allAction := strings.TrimPrefix(action, "all-")
		executeAllAction(item, allAction)

		return allAction
	}

	executeAllAction(item, action)

	return ""
}

func resolveDirtyFiles(item *RemediationItem) []string {
	if len(item.RepoPath) == 0 {
		return item.Files
	}

	diag := gitutil.InspectDirtyState(item.RepoPath)
	if len(diag.AllFiles) > 0 {
		item.Files = diag.AllFiles
		return diag.AllFiles
	}

	return item.Files
}

func printOverflowNote(totalCount, maxShowCount int) {
	if totalCount <= maxShowCount {
		return
	}

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	remainingCount := totalCount - maxShowCount
	fmt.Printf("      %s\n", dimStyle.Render(fmt.Sprintf("... and %d more files", remainingCount)))
}

func formatDirtyFileEntry(raw string) string {
	if strings.HasPrefix(raw, "staged: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("[staged]   ")
		return tag + " " + strings.TrimPrefix(raw, "staged: ")
	}
	if strings.HasPrefix(raw, "modified: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c")).Render("[modified] ")
		return tag + " " + strings.TrimPrefix(raw, "modified: ")
	}
	if strings.HasPrefix(raw, "untracked: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd")).Render("[untracked]")
		return tag + " " + strings.TrimPrefix(raw, "untracked: ")
	}
	if strings.HasPrefix(raw, "deleted: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555")).Render("[deleted]  ")
		return tag + " " + strings.TrimPrefix(raw, "deleted: ")
	}

	return raw
}

func renderDirtyFileList(files []string) {
	fmt.Printf("    Pending Changes (%d files):\n", len(files))
	maxShowCount := 10
	displayCount := len(files)
	if displayCount > maxShowCount {
		displayCount = maxShowCount
	}

	for i := 0; i < displayCount; i++ {
		fmt.Printf("      • %s\n", formatDirtyFileEntry(files[i]))
	}

	printOverflowNote(len(files), maxShowCount)
}

func printRepoDirtyFiles(item *RemediationItem) {
	files := resolveDirtyFiles(item)
	if len(files) == 0 {
		return
	}

	renderDirtyFileList(files)
}

func printPromptOptions() {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Println("    Remediation Options:")
	fmt.Printf("      [1] %-18s %s\n", "Stash & Re-apply", dimStyle.Render("Stash local changes (-u), pull remote, then pop stash"))
	fmt.Printf("      [2] %-18s %s\n", "Commit WIP", dimStyle.Render("Commit all modified/untracked files, then pull --rebase"))
	fmt.Printf("      [3] %-18s %s\n", "Discard Local", dimStyle.Render("Permanently discard changes (reset --hard & clean -fd), pull"))
	fmt.Printf("      [s] %-18s %s\n", "Skip", dimStyle.Render("Skip this repository for now"))
	fmt.Printf("      [a] %-18s %s\n", "Apply to All", dimStyle.Render("Apply stash to this and all remaining repositories"))
	fmt.Printf("      [q] %-18s %s\n", "Quit", dimStyle.Render("Exit interactive prompt"))
}

func resolvePromptChoice(choice string) (string, bool) {
	isQuit := choice == "q" || choice == "quit" || choice == "exit"
	if isQuit {
		return "", true
	}

	action, hasAction := promptChoiceMap[choice]
	if hasAction {
		return action, false
	}

	return "stash", false
}

func promptSingleRepo(reader *bufio.Reader, idx, total int, item *RemediationItem) (string, bool) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	fmt.Printf("\n[%d/%d] %s %s\n", idx, total, titleStyle.Render(item.RepoName), dimStyle.Render("("+item.SummaryReason+")"))
	if len(item.RepoPath) > 0 {
		fmt.Printf("    Path: %s\n", dimStyle.Render(item.RepoPath))
	}
	printRepoDirtyFiles(item)
	printPromptOptions()
	fmt.Printf("  %s ", promptStyle.Render("Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]:"))

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", true
	}

	return resolvePromptChoice(strings.TrimSpace(strings.ToLower(input)))
}

func processPromptStep(reader *bufio.Reader, idx, total int, item *RemediationItem, applyAll string) (string, bool) {
	if len(applyAll) > 0 {
		executeAllAction(item, applyAll)

		return applyAll, false
	}

	action, shouldQuit := promptSingleRepo(reader, idx, total, item)
	if shouldQuit {
		return "", true
	}

	return handlePromptAction(item, action), false
}

func ParseRecipeIndex(option string, recipes []gitutil.RemediationRecipe) int {
	switch option {
	case "1", "stash", "s":
		return 0
	case "2", "wip", "w":
		return 1
	case "3", "discard", "clean", "d":
		return 2
	}

	if i, err := strconv.Atoi(option); err == nil && i > 0 && i <= len(recipes) {
		return i - 1
	}

	return -1
}

func ExecuteFixRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, item.RepoName)
	if recipe.Description != "" {
		fmt.Printf("  Plan:    %s\n\n", recipe.Description)
	}

	if len(recipe.Steps) == 0 && item.RepoPath != "" {
		recipe.Steps = synthesizeRecipeSteps(recipe, item.RepoPath)
	}

	if len(recipe.Steps) > 0 {
		return executeStructuredRecipe(item, recipe)
	}

	return executeShellFallback(item, recipe)
}

func synthesizeRecipeSteps(recipe gitutil.RemediationRecipe, repoPath string) []gitutil.RemediationStep {
	titleLower := strings.ToLower(recipe.Title)
	cmdLower := strings.ToLower(recipe.Command)
	if strings.Contains(titleLower, "wip") || strings.Contains(titleLower, "commit") || strings.Contains(cmdLower, "commit") {
		return gitutil.GenerateCommitRecipe(repoPath).Steps
	}

	if strings.Contains(titleLower, "discard") || strings.Contains(titleLower, "clean") || strings.Contains(cmdLower, "reset --hard") {
		return gitutil.GenerateDiscardRecipe(repoPath).Steps
	}

	if strings.Contains(titleLower, "stash") || strings.Contains(cmdLower, "stash") {
		return gitutil.GenerateStashRecipe(repoPath).Steps
	}

	return nil
}

func executeStructuredRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	total := len(recipe.Steps)
	for i, step := range recipe.Steps {
		err := executeSingleStep(item.RepoName, i+1, total, step)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	RemoveRemediationItem(item.RepoName)

	return nil
}

func formatStepCommand(step gitutil.RemediationStep) string {
	if len(step.Args) >= 2 && step.Args[0] == "-C" {
		return step.Name + " " + strings.Join(step.Args[2:], " ")
	}

	return step.Name + " " + strings.Join(step.Args, " ")
}

func executeSingleStep(repoName string, idx, total int, step gitutil.RemediationStep) error {
	stepCmd := formatStepCommand(step)
	fmt.Printf("  [%d/%d] ➜ %s ... ", idx, total, stepCmd)

	cmd := exec.Command(step.Name, step.Args...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err == nil {
		fmt.Printf("%s ok\n", constants.ColorGreen+"✔"+constants.ColorReset)

		return nil
	}

	return handleStepFailure(repoName, step, outBuf.String(), err)
}

func handleStepFailure(repoName string, step gitutil.RemediationStep, outStr string, err error) error {
	if isBenignCommitClean(step, outStr) || isBenignStashClean(step, outStr) {
		fmt.Printf("%s clean (nothing to commit/apply)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	isRecovered := recoverStashCollision(step, outStr)
	if isRecovered {
		fmt.Printf("%s recovered (dropped stashed untracked collision)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
	printBluntRemediationFailure(repoName, step, outStr, err)

	return err
}

func isBenignStashClean(step gitutil.RemediationStep, output string) bool {
	hasPop := false
	for _, arg := range step.Args {
		if arg == "pop" {
			hasPop = true
			break
		}
	}
	if !hasPop {
		return false
	}

	lower := strings.ToLower(output)

	return strings.Contains(lower, "no stash entries found") || strings.Contains(lower, "no stash found")
}

func isBenignCommitClean(step gitutil.RemediationStep, output string) bool {
	hasCommit := false
	for _, arg := range step.Args {
		if arg == "commit" {
			hasCommit = true
			break
		}
	}
	if !hasCommit {
		return false
	}

	lower := strings.ToLower(output)

	return strings.Contains(lower, "nothing to commit") || strings.Contains(lower, "working tree clean")
}

func recoverStashCollision(step gitutil.RemediationStep, output string) bool {
	hasPop := checkHasPop(step.Args)
	if !hasPop {
		return false
	}

	lowerOut := strings.ToLower(output)
	hasUntracked := strings.Contains(lowerOut, "could not restore untracked files from stash")
	hasCollisionNotice := strings.Contains(lowerOut, "already exists, no checkout")

	if hasUntracked && hasCollisionNotice {
		return dropStashCollision(step.Args)
	}

	return false
}

func checkHasPop(args []string) bool {
	hasPop := false
	for _, arg := range args {
		if arg == "pop" {
			hasPop = true
		}
	}

	return hasPop
}

func dropStashCollision(args []string) bool {
	hasArgs := len(args) >= 2
	if !hasArgs {
		return false
	}

	isDirFlag := args[0] == "-C"
	if !isDirFlag {
		return false
	}

	err := exec.Command("git", "-C", args[1], "stash", "drop").Run()
	isSuccess := err == nil

	return isSuccess
}

func executeShellFallback(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", recipe.Command)
	} else {
		cmd = exec.Command("sh", "-c", recipe.Command)
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err != nil {
		step := gitutil.RemediationStep{Name: "shell", Args: []string{recipe.Command}}
		printBluntRemediationFailure(item.RepoName, step, outBuf.String(), err)

		return err
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	RemoveRemediationItem(item.RepoName)

	return nil
}

type gitDiagnosticInfo struct {
	RCA       string
	Solutions []string
}

func printBluntRemediationFailure(repoName string, step gitutil.RemediationStep, output string, err error) {
	diag := analyzeGitErrorOutput(output, err)
	cmdStr := step.Name + " " + strings.Join(step.Args, " ")

	fmt.Printf("\n%s Remediation Step Failed on %s: %v\n", constants.ColorRed+"✖"+constants.ColorReset, repoName, err)
	fmt.Printf("  Command:   %s\n", cmdStr)
	if len(strings.TrimSpace(output)) > 0 {
		fmt.Printf("  Output:\n    %s\n", strings.ReplaceAll(strings.TrimSpace(output), "\n", "\n    "))
	}

	fmt.Printf("  RCA (Root Cause): %s\n", diag.RCA)
	fmt.Printf("  Known Solutions:\n")
	for _, sol := range diag.Solutions {
		fmt.Printf("    • %s\n", sol)
	}

	fmt.Println()
}

func analyzeGitErrorOutput(output string, err error) gitDiagnosticInfo {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "pathspec") {
		return pathspecDiagnostic()
	}

	if strings.Contains(lower, "would be overwritten by merge") || strings.Contains(lower, "would be overwritten by rebase") {
		return overwrittenDiagnostic()
	}

	if strings.Contains(lower, "conflict") {
		return conflictDiagnostic()
	}

	if strings.Contains(lower, "permission denied") || strings.Contains(lower, "could not read username") {
		return authDiagnostic()
	}

	if strings.Contains(lower, "could not resolve host") || strings.Contains(lower, "unable to access") {
		return networkDiagnostic()
	}

	return defaultDiagnostic(err)
}

func pathspecDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Command arguments were misquoted or split incorrectly by shell wrappers, causing git to interpret words in commit messages as file paths.",
		Solutions: []string{
			"Run native command: git -C <repo> commit -m \"wip: local changes\"",
			"Verify argument quoting to ensure spaces do not fragment commit messages into pathspecs",
		},
	}
}

func overwrittenDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Incoming remote commits conflict with local modified or untracked files.",
		Solutions: []string{
			"Stash untracked files: git -C <repo> stash -u",
			"Or discard untracked changes: git -C <repo> clean -fd && git -C <repo> reset --hard HEAD",
		},
	}
}

func conflictDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Merge or rebase conflict detected between local commits and remote branch.",
		Solutions: []string{
			"Check status: git -C <repo> status",
			"Resolve conflicts in editor, then: git -C <repo> rebase --continue (or git rebase --abort)",
		},
	}
}

func authDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Authentication failure connecting to remote repository. SSH key or credentials missing.",
		Solutions: []string{
			"Verify SSH agent keys: ssh-add -l",
			"Verify remote URL: git -C <repo> remote -v",
		},
	}
}

func networkDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Network connectivity failure. Unable to reach remote Git host.",
		Solutions: []string{
			"Verify internet connection and DNS resolution",
			"Retry pull operation once network connectivity is restored",
		},
	}
}

func defaultDiagnostic(err error) gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: fmt.Sprintf("Git process exited with failure: %v", err),
		Solutions: []string{
			"Inspect working tree status: git -C <repo> status",
			"Check git log: git -C <repo> log -n 3 --oneline",
		},
	}
}
