package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

var primaryTopCommands = []string{
	"scan", "clone", "create", "clone-sync", "pull", "push", "pull-all",
	"status", "reconcile", "fix", "stash", "wip", "discard", "exec",
	"release", "pull-release", "changelog", "cd", "group", "open",
	"history", "stats", "export", "import", "profile", "bookmark",
	"dashboard", "version", "help", "diff", "amend", "sync",
	"add", "rm", "mv", "prune", "revert", "ip", "zsh", "user",
	"service", "schedule", "macro", "os",
}

type commandScore struct {
	cmd  string
	dist int
}

func collectTopCommandCandidates() []string {
	candidates := append([]string{}, primaryTopCommands...)
	candidates = append(candidates, "run", "run-until")
	macroList := macro.ListMacros()
	if macroList.IsSuccess() {
		for _, m := range macroList.Data {
			candidates = append(candidates, m.Name)
		}
	}

	return candidates
}

func suggestTopLevelCommands(command string) []string {
	norm := strings.ToLower(strings.TrimSpace(command))
	if norm == "" {
		return nil
	}

	if norm == "pleae" || norm == "please" {
		return []string{"pull", "release", "pull-release"}
	}

	candidates := collectTopCommandCandidates()
	scores := rankCandidateCommands(norm, candidates)

	return selectBestSuggestions(scores)
}

func rankCandidateCommands(input string, candidates []string) []commandScore {
	var scores []commandScore
	for _, c := range candidates {
		d := levenshtein(input, c)
		isSub := strings.Contains(c, input)
		if isSub {
			d = 1
		}

		if d <= 3 {
			scores = append(scores, commandScore{cmd: c, dist: d})
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].dist < scores[j].dist
	})

	return scores
}

func selectBestSuggestions(scores []commandScore) []string {
	var result []string
	seen := make(map[string]bool)
	for _, sc := range scores {
		if seen[sc.cmd] == false {
			seen[sc.cmd] = true
			result = append(result, sc.cmd)
		}

		if len(result) >= 3 {
			break
		}
	}

	return result
}

func buildUnknownCommandMessage(command string, suggestions []string) string {
	if looksLikeURLToken(command) {
		return fmt.Sprintf(constants.ErrUnknownCommandURLHint, command)
	}

	msg := fmt.Sprintf(constants.ErrUnknownCommand, command)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return msg
}

func printCommandSuggestions(suggestions []string) {
	if len(suggestions) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("  Did you mean one of these?")
	for _, s := range suggestions {
		fmt.Printf("    gitmap %s\n", s)
	}

	fmt.Println()
}

func handleUnknownCommand(command string) {
	suggestions := suggestTopLevelCommands(command)
	printCommandSuggestions(suggestions)
	msg := buildUnknownCommandMessage(command, suggestions)
	printUsage()
	dispatchErr := apperror.NewWithDetails(
		"cmd.dispatch",
		"E1001",
		msg,
		"cmd.root",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"command": command},
	)
	cliexit.HandleError(dispatchErr, 1)
}
