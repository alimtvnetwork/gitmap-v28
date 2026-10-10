package cmdreconcile

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunReconcileCmd executes the reconcile command workflow.
func RunReconcileCmd(args []string) error {
	checkHelp(constants.CmdReconcile, args)
	items := cmdremediation.LoadRemediationState()
	if isReconcileAllRequested(args) {
		return runReconcileAll(args, items)
	}

	if len(args) == 0 {
		return runReconcileList(items)
	}

	return runReconcileTarget(args, items)
}

func runReconcileTarget(args []string, items []cmdremediation.RemediationItem) error {
	repoQuery, action := ParseReconcileArgs(args)
	matched, err := resolveReconcileItem(items, repoQuery)
	if err != nil {
		return err
	}
	if matched == nil {
		return nil
	}

	idx := cmdremediation.ParseRecipeIndex(action, matched.Recipes)
	if idx < 0 || idx >= len(matched.Recipes) {
		idx = 0
	}

	return cmdremediation.ExecuteFixRecipe(matched, matched.Recipes[idx])
}

func resolveReconcileItem(items []cmdremediation.RemediationItem, repoQuery string) (*cmdremediation.RemediationItem, error) {
	matched := cmdremediation.FindRemediationItem(items, repoQuery)
	if matched != nil {
		return matched, nil
	}

	localItem, cleanMsg := cmdremediation.InspectLocalRepoItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)

		return nil, nil
	}
	if localItem != nil {
		return localItem, nil
	}

	return nil, buildReconcileNotFoundError(items, repoQuery)
}

func buildReconcileNotFoundError(items []cmdremediation.RemediationItem, repoQuery string) error {
	suggestions := cmdremediation.FindRemediationSuggestions(items, repoQuery)
	msg := fmt.Sprintf("Repository %q not found in pending reconciliation list.", repoQuery)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return apperror.NewNotFoundError(msg)
}

func isReconcileAllRequested(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-all" || a == "-a" || a == "all" {
			return true
		}
	}

	return false
}

func ParseReconcileArgs(args []string) (string, string) {
	if len(args) == 0 {
		return "", "stash"
	}

	if len(args) == 1 && isNamedAction(args[0]) {
		return "", args[0]
	}

	if len(args) == 1 {
		return args[0], "stash"
	}

	if isNamedAction(args[0]) {
		return args[1], args[0]
	}

	return args[0], args[1]
}

func isNamedAction(s string) bool {
	norm := strings.ToLower(s)

	return norm == "stash" || norm == "wip" || norm == "discard" || norm == "clean"
}

func runReconcileList(items []cmdremediation.RemediationItem) error {
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require reconciliation.\n", constants.ColorGreen+"✓"+constants.ColorReset)

		return nil
	}

	cmdremediation.PrintRemediationSummary(items)

	return nil
}

func runReconcileAll(args []string, items []cmdremediation.RemediationItem) error {
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require reconciliation.\n", constants.ColorGreen+"✓"+constants.ColorReset)

		return nil
	}

	action := resolveReconcileAllAction(args)
	fmt.Printf("%s Reconciling %d repository(ies) with action: %s\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items), action)

	for i := range items {
		idx := cmdremediation.ParseRecipeIndex(action, items[i].Recipes)
		if idx < 0 || idx >= len(items[i].Recipes) {
			idx = 0
		}

		err := cmdremediation.ExecuteFixRecipe(&items[i], items[i].Recipes[idx])
		if err != nil {
			return err
		}
	}

	return nil
}

func resolveReconcileAllAction(args []string) string {
	action := "stash"
	for _, a := range args {
		norm := strings.ToLower(a)
		if norm == "wip" || norm == "discard" || norm == "clean" || norm == "stash" {
			return norm
		}
	}

	return action
}
