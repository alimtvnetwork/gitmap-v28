package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runReconcileCmd(args []string) error {
	checkHelp(constants.CmdReconcile, args)
	items := LoadRemediationState()
	if isReconcileAllRequested(args) {
		return runReconcileAll(args, items)
	}
	if len(args) == 0 {
		return runReconcileList(items)
	}
	repoQuery, action := parseReconcileArgs(args)
	matched := FindRemediationItem(items, repoQuery)
	if matched == nil {
		return apperror.New("reconcile", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Repository %q not found in pending reconciliation list.", repoQuery),
		})
	}
	idx := parseRecipeIndex(action, matched.Recipes)
	if idx < 0 || idx >= len(matched.Recipes) {
		idx = 0
	}
	return executeFixRecipe(matched, matched.Recipes[idx])
}

func isReconcileAllRequested(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-all" || a == "-a" || a == "all" {
			return true
		}
	}
	return false
}

func parseReconcileArgs(args []string) (string, string) {
	if len(args) == 0 {
		return "", "stash"
	}
	if len(args) == 1 {
		if isNamedAction(args[0]) {
			return "", args[0]
		}
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

func runReconcileList(items []RemediationItem) error {
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require reconciliation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	PrintRemediationSummary(items)
	return nil
}

func runReconcileAll(args []string, items []RemediationItem) error {
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require reconciliation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	action := "stash"
	for _, a := range args {
		norm := strings.ToLower(a)
		if norm == "wip" || norm == "discard" || norm == "clean" || norm == "stash" {
			action = norm
			break
		}
	}
	fmt.Printf("%s Reconciling %d repository(ies) with action: %s\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items), action)
	for i := range items {
		idx := parseRecipeIndex(action, items[i].Recipes)
		if idx < 0 || idx >= len(items[i].Recipes) {
			idx = 0
		}
		_ = executeFixRecipe(&items[i], items[i].Recipes[idx])
	}
	return nil
}
