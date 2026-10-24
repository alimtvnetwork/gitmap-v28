package cmdpull

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

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

	if err := gitutil.CheckoutConflictSide(item.RepoPath, enums.ConflictSideTheirs.String(), files); err != nil {
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

	if err := gitutil.CheckoutConflictSide(item.RepoPath, enums.ConflictSideOurs.String(), files); err != nil {
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
	renderVSCodeMergeGuideHeader(repoPath)
	renderVSCodeMergeGuideSteps()
}

func renderVSCodeMergeGuideHeader(repoPath string) {
	fmt.Printf("\n      %s┌──────────────────────────────────────────────────────────────┐%s\n", constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  %s📖 VS Code 3-Way Merge Editor Guide%s                          %s│%s\n",
		constants.ColorMagenta, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s├──────────────────────────────────────────────────────────────┤%s\n", constants.ColorMagenta, constants.ColorReset)
	fmt.Printf("      %s│%s  1. Open repository: %scode %s%s\n",
		constants.ColorMagenta, constants.ColorReset, constants.ColorCyan, repoPath, constants.ColorReset)
}

func renderVSCodeMergeGuideSteps() {
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
	isConflict, _ := detectItemConflicts(item)
	if isConflict {
		dispatchConflictStepAction(ans, item)

		return
	}

	dispatchStandardStepAction(ans, item, bin)
}

func dispatchConflictStepAction(ans string, item UnresolvedPullRepo) {
	if handleConflictInspectAction(ans, item) {
		return
	}

	handleConflictCheckoutAction(ans, item)
}

func handleConflictInspectAction(ans string, item UnresolvedPullRepo) bool {
	switch ans {
	case "v", "code":
		handleVSCodeOpen(item)
		return true
	case "h", "help":
		renderVSCodeMergeGuide(item.RepoPath)
		return true
	case "n", "k", "skip", "no":
		fmt.Printf("      %sSkipped.%s\n", constants.ColorDim, constants.ColorReset)
		return true
	}

	return false
}

func handleConflictCheckoutAction(ans string, item UnresolvedPullRepo) {
	switch ans {
	case "t", enums.ConflictSideTheirs.String():
		handleCheckoutTheirs(item)
	case "o", enums.ConflictSideOurs.String():
		handleCheckoutOurs(item)
	case "a", "abort":
		handleSafeAbort(item)
	}
}

func dispatchStandardStepAction(ans string, item UnresolvedPullRepo, bin string) {
	switch ans {
	case "1", "s", "send", "c", "commit":
		runSendSingle(item, bin)
	case "2", "a", "all":
		runSendAll(bin)
	case "3", "t", "stash":
		runStash(item)
	case "4", "p", "pull":
		runSingleRepoStepPull(item)
	case "n", "k", "skip", "no":
		fmt.Printf("      %sSkipped.%s\n", constants.ColorDim, constants.ColorReset)
	}
}

func runSendSingle(item UnresolvedPullRepo, bin string) {
	cmd := exec.Command(bin, "sends", "cp", item.RepoName, "wip: save changes")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	_ = promptAndSendAllFleet()
}

func runSendAll(bin string) {
	cmd := exec.Command(bin, "sends", "cp", "all", "wip: save changes")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func runStash(item UnresolvedPullRepo) {
	cmd := exec.Command("git", "-C", item.RepoPath, "stash")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("      %s✖ Stash failed: %s%s\n", constants.ColorRed, strings.TrimSpace(string(out)), constants.ColorReset)

		return
	}

	fmt.Printf("      %s✓ Stashed working tree.%s\n", constants.ColorGreen, constants.ColorReset)
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
