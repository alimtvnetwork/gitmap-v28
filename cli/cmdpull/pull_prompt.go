package cmdpull

import (
	"bufio"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"golang.org/x/term"
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

// ExtractUnresolvedRepos extracts repositories that failed or suffered unresolvable conflicts or are dirty.
func ExtractUnresolvedRepos(states []*PullRepoState) []UnresolvedPullRepo {
	var unresolved []UnresolvedPullRepo
	for _, s := range states {
		if s == nil {
			continue
		}
		isDirty := resolveRepoDirtyFlag(s)
		hasProblem := isDirty || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.ErrorMsg != ""
		if hasProblem {
			unresolved = append(unresolved, buildUnresolvedItem(s, isDirty))
		}
	}

	return unresolved
}

func resolveRepoDirtyFlag(s *PullRepoState) bool {
	if s.IsDirty || s.Changes == "dirty" {
		return true
	}
	if s.RepoPath == "" {
		return false
	}
	diag := gitutil.InspectDirtyState(s.RepoPath)

	return diag.IsDirty
}

func buildUnresolvedItem(s *PullRepoState, isDirty bool) UnresolvedPullRepo {
	return UnresolvedPullRepo{
		RepoName: s.RepoName,
		RepoPath: s.RepoPath,
		Branch:   s.Branch,
		ErrorMsg: s.ErrorMsg,
		IsDirty:  isDirty,
		IsFailed: s.Step == PullStepTypeError || s.ErrorMsg != "",
	}
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
	if isNonInteractiveEnv() {
		return false
	}

	return isTerminalAttached()
}

func isTerminalAttached() bool {
	if term.IsTerminal(int(os.Stdin.Fd())) || term.IsTerminal(int(os.Stdout.Fd())) {
		return true
	}
	if termout.IsStdoutTTY() {
		return true
	}
	info, err := os.Stdin.Stat()

	return err == nil && (info.Mode()&os.ModeCharDevice) != 0
}

func isNonInteractiveEnv() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") != "" || os.Getenv("GITMAP_JSON") != ""
}

func readPromptUserChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "q"
	}

	return strings.ToLower(strings.TrimSpace(line))
}

func resolveGitmapExecutable() string {
	self, err := os.Executable()
	if err == nil && self != "" {
		return self
	}

	return "gitmap"
}
