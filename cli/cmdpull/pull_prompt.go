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
	fmt.Println("    [1/a] Fleet Multi-Repo Control: Run `gitmap sends cp all \"wip: save changes\"`")
	fmt.Println("    [2/s] Step through each repository individually")
	fmt.Println("    [3/q/k] Skip and leave as-is")
	fmt.Printf("  %sChoice [1/2/3/q]:%s ", constants.ColorCyan, constants.ColorReset)
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
	if choice == "1" || choice == "a" || choice == "all" {
		return executeFleetCommitWIP(unresolved)
	}
	if choice == "2" || choice == "s" || choice == "single" || choice == "step" {
		return executeStepByStepRemediation(unresolved)
	}

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
		fmt.Printf("\n  [%d/%d] Repository: %s%s%s (branch: %s)\n",
			i+1, len(unresolved),
			constants.ColorBold, item.RepoName, constants.ColorReset,
			item.Branch)
		if item.ErrorMsg != "" {
			fmt.Printf("      Error: %s%s%s\n", constants.ColorRed, item.ErrorMsg, constants.ColorReset)
		}
		if item.IsDirty {
			fmt.Printf("      State: %sworking tree dirty%s\n", constants.ColorYellow, constants.ColorReset)
		}
		fmt.Printf("      Path:  %s\n", item.RepoPath)
		fmt.Printf("  %sAction for %s? [1=commit WIP / 2=pull / 3=skip / q=quit]:%s ",
			constants.ColorCyan, item.RepoName, constants.ColorReset)
		line, _ := reader.ReadString('\n')
		ans := strings.ToLower(strings.TrimSpace(line))
		if ans == "q" || ans == "quit" {
			break
		}
		if ans == "1" || ans == "c" || ans == "commit" {
			cmd := exec.Command(bin, "sends", "cp", item.RepoName, "wip: save changes")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		}
		if ans == "2" || ans == "p" || ans == "pull" {
			runSingleRepoStepPull(item)
		}
	}

	return nil
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
