package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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

func executeSendCommitsCurrentMachine(unresolved []UnresolvedPullRepo) error {
	bin := resolveGitmapExecutable()
	fmt.Printf("\n  %s⚡ Sending commits for dirty repositories on current machine...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	count := 0
	for _, item := range unresolved {
		if item.IsDirty && runSendCpTarget(bin, item) {
			count++
		}
	}
	fmt.Printf("\n  %s✓ Sent commits on current machine: %d repository(ies) pushed.%s\n\n",
		constants.ColorGreen, count, constants.ColorReset)

	return promptAndSendAllFleet()
}

func runSendCpTarget(bin string, item UnresolvedPullRepo) bool {
	target := resolveDirtyItemTarget(item)
	fmt.Printf("  • %s%s%s (branch: %s)\n", constants.ColorBold, target, constants.ColorReset, item.Branch)
	cmd := exec.Command(bin, "sends", "cp", target, "wip: save changes")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("    %s %s✖%s sends cp failed for %s: %v\n", constants.TreeCorner, constants.ColorRed, constants.ColorReset, target, err)
		return false
	}

	return true
}

func resolveDirtyItemTarget(item UnresolvedPullRepo) string {
	if item.RepoName != "" {
		return item.RepoName
	}

	return filepath.Base(item.RepoPath)
}

func promptAndSendAllFleet() error {
	if !isInteractivePromptEligible() {
		return nil
	}
	if askSendAllFleetConfirmation() {
		return executeSendAllCommits()
	}

	return nil
}

func askSendAllFleetConfirmation() bool {
	fmt.Printf("  %sWould you also like to send commits across all repositories (Fleet Multi-Repo)? [y/N]:%s ",
		constants.ColorCyan, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))

	return ans == "y" || ans == "yes"
}

func executeSendAllCommits() error {
	bin := resolveGitmapExecutable()
	fmt.Printf("\n  %s⚡ Sending commits across all repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := exec.Command(bin, "sends", "cp", "all", "wip: save changes")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "sends cp all failed")
	}

	return nil
}

func executeStashAll(unresolved []UnresolvedPullRepo) error {
	fmt.Printf("\n  %s📦 Stashing changes across unresolved repositories...%s\n\n",
		constants.ColorCyan, constants.ColorReset)
	count := 0
	for _, item := range unresolved {
		if item.IsDirty && stashSingleRepoItem(item) {
			count++
		}
	}
	fmt.Printf("\n  %s✓ Stash complete: %d repository(ies) stashed.%s\n\n",
		constants.ColorGreen, count, constants.ColorReset)

	return nil
}

func stashSingleRepoItem(item UnresolvedPullRepo) bool {
	fmt.Printf("  • %s%s%s (branch: %s)\n", constants.ColorBold, item.RepoName, constants.ColorReset, item.Branch)
	cmd := exec.Command("git", "-C", item.RepoPath, "stash")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("    %s %s✖%s git stash failed: %s\n", constants.TreeCorner, constants.ColorRed, constants.ColorReset, strings.TrimSpace(string(out)))
		return false
	}
	fmt.Printf("    %s %s✓%s Stashed working tree.\n", constants.TreeCorner, constants.ColorGreen, constants.ColorReset)

	return true
}

