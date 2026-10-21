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
		isDirty := s.IsDirty || s.Changes == "dirty"
		if !isDirty && s.RepoPath != "" {
			diag := gitutil.InspectDirtyState(s.RepoPath)
			if diag.IsDirty {
				isDirty = true
			}
		}
		hasProblem := isDirty || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.ErrorMsg != ""
		if hasProblem {
			unresolved = append(unresolved, UnresolvedPullRepo{
				RepoName: s.RepoName,
				RepoPath: s.RepoPath,
				Branch:   s.Branch,
				ErrorMsg: s.ErrorMsg,
				IsDirty:  isDirty,
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
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return true
	}
	if termout.IsStdoutTTY() || term.IsTerminal(int(os.Stdout.Fd())) {
		fi, err := os.Stdin.Stat()
		if err == nil && (fi.Mode()&os.ModeNamedPipe) == 0 {
			return true
		}
	}

	return false
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
