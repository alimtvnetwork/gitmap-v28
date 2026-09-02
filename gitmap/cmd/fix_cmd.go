package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

func runFix(args []string, aliasOverride string) error {
	items := LoadRemediationState()
	if len(items) == 0 {
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
	idx := parseRecipeIndex(action, item.Recipes)
	if idx < 0 || idx >= len(item.Recipes) {
		return apperror.New("fix", "E_INVALID_OPTION", map[string]any{"msg": fmt.Sprintf("Invalid fix option: %s", action)})
	}
	return executeFixRecipe(item, item.Recipes[idx])
}

func resolveFixTarget(args []string, aliasOverride string, items []RemediationItem) (*RemediationItem, string, error) {
	repoQuery, action := parseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}
	if repoQuery != "" {
		matched := FindRemediationItem(items, repoQuery)
		if matched == nil {
			return nil, "", apperror.New("fix", "E_NOT_FOUND", map[string]any{"msg": fmt.Sprintf("Repository %q not found in pending remediation list.", repoQuery)})
		}
		return matched, action, nil
	}
	if len(items) == 1 {
		return &items[0], action, nil
	}
	PrintRemediationSummary(items)
	return nil, "", apperror.New("fix", "E_AMBIGUOUS", map[string]any{"msg": "Multiple repositories need remediation. Specify repo: gitmap fix <repo> <action>"})
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

func executeFixRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, item.RepoName)
	fmt.Printf("  Running: %s\n\n", recipe.Command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", recipe.Command)
	} else {
		cmd = exec.Command("sh", "-c", recipe.Command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s Fix failed: %v\n", constants.ColorRed+"✗"+constants.ColorReset, err)
		return nil
	}
	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	RemoveRemediationItem(item.RepoName)
	return nil
}
