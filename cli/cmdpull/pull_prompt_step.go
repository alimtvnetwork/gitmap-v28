package cmdpull

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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
	fmt.Printf("        %s[1/s]%s %sSend Commits%s        - gitmap sends cp %s \"wip: save changes\"\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset, repoName)
	fmt.Printf("        %s[2/a]%s %sSend All Repos%s      - gitmap sends cp all \"wip: save changes\"\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorGreen, constants.ColorReset)
	fmt.Printf("        %s[3/t]%s %sStash Changes%s       - git stash working tree changes\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
	fmt.Printf("        %s[4/p]%s %sRetry Pull%s          - Pull remote changes with autostash\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("        %s[n/k]%s %sSkip%s                - Leave repository as-is\n",
		constants.ColorCyan, constants.ColorReset, constants.ColorDim, constants.ColorReset)
	fmt.Printf("        %s[q]%s   %sQuit%s                - Exit step-by-step resolution\n",
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
