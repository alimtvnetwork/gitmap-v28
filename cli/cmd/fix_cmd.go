package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func runFix(args []string, aliasOverride string) error {
	items := LoadRemediationState()
	if len(items) == 0 {
		if len(args) > 0 {
			return runFixDirect(args, aliasOverride)
		}

		fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		fmt.Println("  Run 'gitmap pull' to pull all tracked repositories.")

		return nil
	}

	if len(args) == 0 && aliasOverride == "" {
		PrintRemediationSummary(items)

		return nil
	}

	item, action, err := resolveFixTarget(args, aliasOverride, items)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}

	return applyFixRecipe(item, action)
}

func runFixDirect(args []string, aliasOverride string) error {
	repoQuery, action := parseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}

	item, cleanMsg := InspectLocalRepoItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)

		return nil
	}
	if item == nil {
		return apperror.New("fix", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Repository %q not found or not a git repository.", repoQuery),
		})
	}

	return applyFixRecipe(item, action)
}

func applyFixRecipe(item *RemediationItem, action string) error {
	idx := parseRecipeIndex(action, item.Recipes)
	if idx < 0 || idx >= len(item.Recipes) {
		return apperror.New("fix", "E_INVALID_OPTION", map[string]any{
			"msg": fmt.Sprintf("Invalid fix option: %s", action),
		})
	}

	return executeFixRecipe(item, item.Recipes[idx])
}

func resolveFixTarget(args []string, aliasOverride string, items []RemediationItem) (*RemediationItem, string, error) {
	repoQuery, action := parseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}

	if repoQuery != "" {
		return findItemOrError(items, repoQuery, action)
	}

	if len(items) == 1 {
		return &items[0], action, nil
	}

	PrintRemediationSummary(items)

	return nil, "", apperror.New("fix", "E_AMBIGUOUS", map[string]any{
		"msg": "Multiple repositories need remediation. Specify repo: gitmap fix <repo> <action>",
	})
}

func parseRecipeIndex(option string, recipes []gitutil.RemediationRecipe) int {
	switch option {
	case "1", "stash", "s":
		return 0
	case "2", "wip", "w":
		return 1
	case "3", "discard", "clean", "d":
		return 2
	}

	if i, err := strconv.Atoi(option); err == nil && i > 0 && i <= len(recipes) {
		return i - 1
	}

	return -1
}

func findItemOrError(items []RemediationItem, repoQuery, action string) (*RemediationItem, string, error) {
	matched := FindRemediationItem(items, repoQuery)
	if matched != nil {
		return matched, action, nil
	}

	localItem, cleanMsg := InspectLocalRepoItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)

		return nil, action, nil
	}
	if localItem != nil {
		return localItem, action, nil
	}

	return nil, "", buildFixNotFoundError(items, repoQuery)
}

func buildFixNotFoundError(items []RemediationItem, repoQuery string) error {
	suggestions := FindRemediationSuggestions(items, repoQuery)
	msg := fmt.Sprintf("Repository %q not found in pending remediation list.", repoQuery)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return apperror.New("fix", "E_NOT_FOUND", map[string]any{"msg": msg})
}
