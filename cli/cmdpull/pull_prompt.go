package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// UnresolvedPullRepo encapsulates repository details that require user attention.
type UnresolvedPullRepo struct {
	RepoName string
	RepoPath string
	Branch   string
	ErrorMsg string
	IsDirty  bool
	IsFailed bool
}

// ExtractUnresolvedRepos extracts repositories that failed or suffered unresolvable conflicts.
func ExtractUnresolvedRepos(states []*PullRepoState) []UnresolvedPullRepo {
	var unresolved []UnresolvedPullRepo
	for _, s := range states {
		hasProblem := s.IsDirty || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.ErrorMsg != ""
		if hasProblem {
			unresolved = append(unresolved, UnresolvedPullRepo{
				RepoName: s.RepoName,
				RepoPath: s.RepoPath,
				Branch:   s.Branch,
				ErrorMsg: s.ErrorMsg,
				IsDirty:  s.IsDirty,
				IsFailed: s.Step == PullStepTypeError || s.ErrorMsg != "",
			})
		}
	}

	return unresolved
}

// PromptPostPullRemediation presents an interactive resolution menu when unresolved repos exist.
func PromptPostPullRemediation(unresolved []UnresolvedPullRepo) error {
	hasUnresolved := len(unresolved) > 0
	if !hasUnresolved || !isInteractivePromptEligible() {
		return nil
	}

	renderPostPullUnresolvedList(unresolved)
	renderPostPullMenu()

	choice := readPromptUserChoice()
	return dispatchPostPullChoice(choice, unresolved)
}

// isInteractivePromptEligible verifies terminal TTY and standard input availability.
func isInteractivePromptEligible() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" || os.Getenv("GITMAP_JSON") != "" {
		return false
	}
	if !termout.IsStdoutTTY() {
		return false
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func renderPostPullUnresolvedList(unresolved []UnresolvedPullRepo) {
	fmt.Printf("\n  %s[!]%s %s%d repository(ies) require attention after pull:%s\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorBold, len(unresolved), constants.ColorReset)

	for _, item := range unresolved {
		reason := "working tree dirty / autostash conflict"
		if item.IsFailed && item.ErrorMsg != "" {
			reason = item.ErrorMsg
		}
		fmt.Printf("      • %s%s%s: %s%s%s\n",
			constants.ColorBold, item.RepoName, constants.ColorReset,
			constants.ColorDim, reason, constants.ColorReset)
	}
}

func renderPostPullMenu() {
	fmt.Printf("\n  %sSome repositories could not be pulled automatically. How would you like to resolve?%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[1/f/fix]%s    %sFix All%s                         - Run %sgitmap fix --all%s across failed repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[2/c/commit]%s %sCommit All%s                      - Commit WIP locally on this machine only\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %s[3/t/stash]%s  %sStash All%s                       - Run %sgit stash%s across dirty repositories\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[4/s/step]%s   %sStep through individually%s       - Tree-view diff & colored action previews\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s[5/fleet]%s    %sFleet Multi-Repo Control%s        - Run %sgitmap sends cp all \"wip: save changes\"%s\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	fmt.Printf("    %s[6/q/skip]%s   %sSkip and leave as-is%s            - No further action taken\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
	fmt.Printf("  %sChoice [1/2/3/4/5/6/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func readPromptUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

func dispatchPostPullChoice(choice string, unresolved []UnresolvedPullRepo) error {
	switch choice {
	case "1", "f", "fix", "fix-all":
		return executeFixAll()
	case "2", "c", "commit", "commit-all":
		return executeLocalCommitAllWIP(unresolved)
	case "3", "t", "stash", "stash-all":
		return executeStashAll(unresolved)
	case "4", "s", "step", "single":
		return executeStepByStepRemediation(unresolved)
	case "5", "fleet", "sends":
		return executeFleetCommitWIP(unresolved)
	case "6", "q", "quit", "skip", "n", "no":
		return nil
	default:
		return nil
	}
}

func executeFixAll() error {
	bin := resolveGitmapExecutable()
	fmt.Printf("\n  %s🔧 Running gitmap fix --all across failed repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := exec.Command(bin, "fix", "--all")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "fix --all failed")
	}

	return nil
}

func executeLocalCommitAllWIP(unresolved []UnresolvedPullRepo) error {
	fmt.Printf("\n  %s💾 Committing WIP locally across unresolved repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	count := 0
	for _, item := range unresolved {
		if !item.IsDirty {
			continue
		}
		fmt.Printf("  • %s%s%s (branch: %s)\n", constants.ColorBold, item.RepoName, constants.ColorReset, item.Branch)
		addCmd := exec.Command("git", "-C", item.RepoPath, "add", "-A")
		if err := addCmd.Run(); err != nil {
			fmt.Printf("    %s %s✖%s git add failed: %v\n", constants.TreeCorner, constants.ColorRed, constants.ColorReset, err)
			continue
		}
		commitCmd := exec.Command("git", "-C", item.RepoPath, "commit", "-m", "wip: save changes")
		out, err := commitCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("    %s %s✖%s commit failed: %s\n", constants.TreeCorner, constants.ColorRed, constants.ColorReset, strings.TrimSpace(string(out)))
		} else {
			fmt.Printf("    %s %s✓%s Committed WIP locally.\n", constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
			count++
		}
	}
	fmt.Printf("\n  %s✓ Local WIP commit complete: %d repository(ies) committed.%s\n\n",
		constants.ColorGreen, count, constants.ColorReset)
	return nil
}

func executeStashAll(unresolved []UnresolvedPullRepo) error {
	fmt.Printf("\n  %s📦 Stashing changes across unresolved repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	count := 0
	for _, item := range unresolved {
		if !item.IsDirty {
			continue
		}
		fmt.Printf("  • %s%s%s (branch: %s)\n", constants.ColorBold, item.RepoName, constants.ColorReset, item.Branch)
		cmd := exec.Command("git", "-C", item.RepoPath, "stash")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("    %s %s✖%s git stash failed: %s\n", constants.TreeCorner, constants.ColorRed, constants.ColorReset, strings.TrimSpace(string(out)))
		} else {
			fmt.Printf("    %s %s✓%s Stashed working tree.\n", constants.TreeCorner, constants.ColorGreen, constants.ColorReset)
			count++
		}
	}
	fmt.Printf("\n  %s✓ Stash complete: %d repository(ies) stashed.%s\n\n",
		constants.ColorGreen, count, constants.ColorReset)
	return nil
}

func executeFleetCommitWIP(unresolved []UnresolvedPullRepo) error {
	bin := resolveGitmapExecutable()
	fmt.Printf("\n  %s⚡ Running fleet commit & push across repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := exec.Command(bin, "sends", "cp", "all", "wip: save changes")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "fleet commit failed")
	}

	return nil
}

func resolveGitmapExecutable() string {
	self, err := os.Executable()
	if err == nil && self != "" {
		return self
	}

	return "gitmap"
}

func executeStepByStepRemediation(unresolved []UnresolvedPullRepo) error {
	reader := bufio.NewReader(os.Stdin)
	bin := resolveGitmapExecutable()
	for i, item := range unresolved {
		fmt.Printf("\n  [%d/%d] Repository: %s%s%s (branch: %s%s%s)\n",
			i+1, len(unresolved),
			constants.ColorBold, item.RepoName, constants.ColorReset,
			constants.ColorCyan, item.Branch, constants.ColorReset)
		fmt.Printf("      Path:  %s%s%s\n", constants.ColorDim, item.RepoPath, constants.ColorReset)

		isConflict, _ := detectItemConflicts(item)
		if item.ErrorMsg != "" && !isConflict {
			fmt.Printf("      Error: %s%s%s\n", constants.ColorRed, item.ErrorMsg, constants.ColorReset)
		}
		if !isConflict && item.IsDirty {
			fmt.Printf("      State: %sWorking tree dirty%s\n", constants.ColorYellow, constants.ColorReset)
			renderRepoShortDiff(item.RepoPath)
		}

		for {
			renderStepByStepActions(item)
			line, _ := reader.ReadString('\n')
			ans := strings.ToLower(strings.TrimSpace(line))
			if ans == "q" || ans == "quit" {
				return nil
			}
			if ans == "h" || ans == "help" {
				renderVSCodeMergeGuide(item.RepoPath)
				continue
			}
			if ans == "v" || ans == "code" {
				dispatchStepAction(ans, item, bin)
				continue
			}
			dispatchStepAction(ans, item, bin)
			break
		}
	}

	return nil
}

func renderRepoShortDiff(repoPath string) {
	statusOut, err := exec.Command("git", "-C", repoPath, "status", "--short").Output()
	if err != nil {
		return
	}
	rawLines := strings.Split(strings.TrimSpace(string(statusOut)), "\n")
	var lines []string
	for _, l := range rawLines {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	if len(lines) == 0 {
		return
	}
	fmt.Printf("      %s %sUncommitted files (%d):%s\n", constants.TreeBranch, constants.ColorCyan, len(lines), constants.ColorReset)
	maxShow := 4
	if len(lines) < maxShow {
		maxShow = len(lines)
	}
	for idx := 0; idx < maxShow; idx++ {
		connector := constants.TreeBranch
		if idx == maxShow-1 && len(lines) <= 4 {
			connector = constants.TreeCorner
		}
		fmt.Printf("      │   %s %s\n", connector, lines[idx])
	}
	if len(lines) > 4 {
		fmt.Printf("      │   %s ... and %d more modified file(s)\n", constants.TreeCorner, len(lines)-4)
	}
}

func detectItemConflicts(item UnresolvedPullRepo) (bool, []string) {
	hasConflictText := strings.Contains(strings.ToLower(item.ErrorMsg), "conflict") ||
		strings.Contains(item.ErrorMsg, "E_MERGE_CONFLICT")

	conflictedFiles := cloner.ExtractConflictedFilesFromRepo(item.RepoPath)
	if len(conflictedFiles) == 0 && item.ErrorMsg != "" {
		conflictedFiles = cloner.ExtractConflictedFilesFromOutput(item.ErrorMsg)
	}

	if hasConflictText || len(conflictedFiles) > 0 {
		return true, conflictedFiles
	}

	return false, nil
}

func renderConflictBox(repoName string, conflictedFiles []string) {
	fmt.Printf("\n      %s╔══════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorRed, constants.ColorReset)
	fmt.Printf("      %s║%s  %s%s⚠ MERGE CONFLICT DETECTED IN %s%s\n",
		constants.ColorRed, constants.ColorReset,
		constants.ColorBold, constants.ColorYellow, repoName, constants.ColorReset)
	fmt.Printf("      %s╠══════════════════════════════════════════════════════════════════════╣%s\n", constants.ColorRed, constants.ColorReset)
	if len(conflictedFiles) > 0 {
		fmt.Printf("      %s║%s  %sConflicted file(s) (%d):%s\n",
			constants.ColorRed, constants.ColorReset,
			constants.ColorBold, len(conflictedFiles), constants.ColorReset)
		for _, f := range conflictedFiles {
			fmt.Printf("      %s║%s    %s•%s %s%s%s\n",
				constants.ColorRed, constants.ColorReset,
				constants.ColorYellow, constants.ColorReset,
				constants.ColorCyan, f, constants.ColorReset)
		}
	} else {
		fmt.Printf("      %s║%s    %sUnresolved merge conflicts detected during pull.%s\n",
			constants.ColorRed, constants.ColorReset,
			constants.ColorYellow, constants.ColorReset)
	}
	fmt.Printf("      %s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorRed, constants.ColorReset)
}

func renderConflictStepActions(repoName string, repoPath string) {
	fmt.Printf("      %s %sConflict Resolution Actions for %s:%s\n", constants.TreeCorner, constants.ColorBold, repoName, constants.ColorReset)
	fmt.Printf("        %s[v/code]%s   %sOpen in VS Code%s          - code %s\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorBold, constants.ColorReset, repoPath)
	fmt.Printf("        %s[t/theirs]%s %sAccept Incoming%s          - Remote changes win (git checkout --theirs)\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("        %s[o/ours]%s   %sAccept Current%s           - Local changes win (git checkout --ours)\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("        %s[a/abort]%s  %sSafe Abort Merge%s         - Reset repository to pre-merge state\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("        %s[h/help]%s   %sVS Code Merge Guide%s      - How to use 3-way merge editor\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("        %s[n/k]%s      %sSkip%s                     - Leave repository in conflict state\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
	fmt.Printf("        %s[q]%s        %sQuit%s                     - Exit step-by-step resolution\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorRed, constants.ColorReset)
	fmt.Printf("      %sChoice [v/t/o/a/h/n/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func renderStandardStepActions(repoName string) {
	fmt.Printf("      %s %sAction for %s:%s\n", constants.TreeCorner, constants.ColorBold, repoName, constants.ColorReset)
	fmt.Printf("        %s[1/c]%s %sCommit WIP locally%s    - Commit \"wip: save changes\" locally\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("        %s[2/p]%s %sRetry Pull%s           - Pull remote changes with autostash\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("        %s[3/s]%s %sStash Changes%s        - git stash working tree changes\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("        %s[4/f]%s %sFleet Commit & Push%s  - gitmap sends cp %s\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorMagenta, constants.ColorReset, repoName)
	fmt.Printf("        %s[n/k]%s %sSkip%s                 - Leave repository as-is\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
	fmt.Printf("        %s[q]%s   %sQuit%s                 - Exit step-by-step resolution\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorRed, constants.ColorReset)
	fmt.Printf("      %sChoice [1/2/3/4/n/q]:%s ", constants.ColorCyan, constants.ColorReset)
}

func renderStepByStepActions(item UnresolvedPullRepo) {
	isConflict, conflictedFiles := detectItemConflicts(item)
	if isConflict {
		renderConflictBox(item.RepoName, conflictedFiles)
		renderConflictStepActions(item.RepoName, item.RepoPath)
		return
	}

	renderStandardStepActions(item.RepoName)
}

func handleVSCodeOpen(item UnresolvedPullRepo) {
	fmt.Printf("      %sLaunching VS Code for %s...%s\n", constants.ColorCyan, item.RepoPath, constants.ColorReset)
	cmd := exec.Command("code", item.RepoPath)
	if err := cmd.Start(); err != nil {
		fmt.Printf("      %sℹ VS Code CLI ('code') could not be launched (%v). Run manually:%s\n        code %s\n",
			constants.ColorYellow, err, constants.ColorReset, item.RepoPath)
		return
	}

	fmt.Printf("      %s✓ Opened %s in VS Code.%s\n", constants.ColorGreen, item.RepoPath, constants.ColorReset)
}

func resolveConflictFiles(item UnresolvedPullRepo) []string {
	files := cloner.ExtractConflictedFilesFromRepo(item.RepoPath)
	if len(files) == 0 && item.ErrorMsg != "" {
		files = cloner.ExtractConflictedFilesFromOutput(item.ErrorMsg)
	}

	return files
}

func handleCheckoutTheirs(item UnresolvedPullRepo) {
	files := resolveConflictFiles(item)
	if len(files) == 0 {
		fmt.Printf("      %s⚠ No conflicted files identified to accept with incoming changes.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}

	if err := gitutil.CheckoutConflictSide(item.RepoPath, "theirs", files); err != nil {
		fmt.Printf("      %s✖ Accept Incoming failed: %v%s\n", constants.ColorRed, err, constants.ColorReset)
		return
	}

	fmt.Printf("      %s✓ Accepted incoming (remote) changes for %d file(s) and staged.%s\n",
		constants.ColorGreen, len(files), constants.ColorReset)
}

func handleCheckoutOurs(item UnresolvedPullRepo) {
	files := resolveConflictFiles(item)
	if len(files) == 0 {
		fmt.Printf("      %s⚠ No conflicted files identified to accept with current changes.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}

	if err := gitutil.CheckoutConflictSide(item.RepoPath, "ours", files); err != nil {
		fmt.Printf("      %s✖ Accept Current failed: %v%s\n", constants.ColorRed, err, constants.ColorReset)
		return
	}

	fmt.Printf("      %s✓ Accepted current (local) changes for %d file(s) and staged.%s\n",
		constants.ColorGreen, len(files), constants.ColorReset)
}

func handleSafeAbort(item UnresolvedPullRepo) {
	if err := gitutil.SafeAbortMergeOrRebase(item.RepoPath); err != nil {
		fmt.Printf("      %s✖ Safe abort failed: %v%s\n", constants.ColorRed, err, constants.ColorReset)
		return
	}

	fmt.Printf("      %s✓ Safe abort completed: repository restored to pre-merge state.%s\n",
		constants.ColorGreen, constants.ColorReset)
}

func renderVSCodeMergeGuide(repoPath string) {
	fmt.Printf("\n      %s┌──────────────────────────────────────────────────────────────┐%s\n", constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  %s📖 VS Code 3-Way Merge Editor Guide%s                          %s│%s\n",
		constants.ColorMagenta, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s├──────────────────────────────────────────────────────────────┤%s\n", constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  1. Open repository: %scode %s%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorCyan, repoPath, constants.ColorReset)
	fmt.Printf("      %s│%s  2. Switch to Source Control tab (Ctrl+Shift+G / Cmd+Shift+G)  %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  3. Click a file under 'Merge Changes'                         %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  4. Click 'Resolve in Merge Editor' (bottom right)             %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  5. Compare Incoming (Remote) vs Current (Local) changes        %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  6. Check checkboxes to accept changes or edit Result below    %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  7. Click 'Complete Merge' and stage changes                    %s│%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s└──────────────────────────────────────────────────────────────┘%s\n\n", constants.ColorMagenta, constants.ColorReset)
}

func dispatchStepAction(ans string, item UnresolvedPullRepo, bin string) {
	switch ans {
	case "v", "code":
		handleVSCodeOpen(item)
	case "t", "theirs":
		handleCheckoutTheirs(item)
	case "o", "ours":
		handleCheckoutOurs(item)
	case "a", "abort":
		handleSafeAbort(item)
	case "h", "help":
		renderVSCodeMergeGuide(item.RepoPath)
	case "1", "c", "commit":
		addCmd := exec.Command("git", "-C", item.RepoPath, "add", "-A")
		_ = addCmd.Run()
		commitCmd := exec.Command("git", "-C", item.RepoPath, "commit", "-m", "wip: save changes")
		out, err := commitCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("      %s✖ Commit failed: %s%s\n", constants.ColorRed, strings.TrimSpace(string(out)), constants.ColorReset)
		} else {
			fmt.Printf("      %s✓ Committed WIP locally.%s\n", constants.ColorGreen, constants.ColorReset)
		}
	case "2", "p", "pull":
		runSingleRepoStepPull(item)
	case "3", "s", "stash":
		cmd := exec.Command("git", "-C", item.RepoPath, "stash")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("      %s✖ Stash failed: %s%s\n", constants.ColorRed, strings.TrimSpace(string(out)), constants.ColorReset)
		} else {
			fmt.Printf("      %s✓ Stashed working tree.%s\n", constants.ColorGreen, constants.ColorReset)
		}
	case "4", "f", "fleet":
		cmd := exec.Command(bin, "sends", "cp", item.RepoName, "wip: save changes")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	case "n", "k", "skip", "no":
		fmt.Printf("      %sSkipped.%s\n", constants.ColorDim, constants.ColorReset)
	}
}

func runSingleRepoStepPull(item UnresolvedPullRepo) {
	rec := model.ScanRecord{AbsolutePath: item.RepoPath, RepoName: item.RepoName, Branch: item.Branch}
	res := cloner.SafePullOne(rec, item.RepoPath)
	if res.IsSuccess {
		fmt.Printf("      %s✓ Pull successful%s\n", constants.ColorGreen, constants.ColorReset)
		return
	}

	fmt.Printf("      %s✖ Pull failed: %s%s\n", constants.ColorRed, res.Error, constants.ColorReset)
}
