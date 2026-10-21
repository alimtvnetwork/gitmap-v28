package cmdfix

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func isFixActionWord(s string) bool {
	switch s {
	case "1", "stash", "s", "clone", "c", "2", "wip", "w", "init", "i", "link", "3", "discard", "clean", "d", "rm", "remove",
		"theirs", "t", "ours", "o", "resolve", "ide", "code", "v", "abort", "a", "rerere", "case", "casing":
		return true
	}
	return false
}

func tryExecuteTargetedFixAction(targetPath, action string) (bool, error) {
	if targetPath == "" {
		return false, nil
	}
	return executeTargetedFixAction(targetPath, action)
}

func applyFixRecipe(item *cmdremediation.RemediationItem, action string) error {
	ensureCaseCollisionRecipe(item)
	handled, err := executeTargetedFixAction(item.RepoPath, action)
	if handled {
		return err
	}

	if action == "" && len(item.Recipes) > 0 {
		return executeFixRecipe(item, item.Recipes[0])
	}
	idx := parseRecipeIndex(action, item.Recipes)
	if idx < 0 || idx >= len(item.Recipes) {
		return applyFallbackOrInvalid(item, action)
	}

	return executeFixRecipe(item, item.Recipes[idx])
}

func ensureCaseCollisionRecipe(item *cmdremediation.RemediationItem) {
	if item == nil || item.RepoPath == "" {
		return
	}
	cols, err := gitutil.DetectCaseCollisions(item.RepoPath)
	if err != nil || len(cols) == 0 {
		return
	}
	if hasExistingCaseRecipe(item.Recipes) {
		return
	}
	caseRecipes := buildCaseRecipes(item.RepoPath, cols)
	item.Recipes = append(caseRecipes, item.Recipes...)
	_ = cmdremediation.SaveRemediationState([]cmdremediation.RemediationItem{*item})
}

func hasExistingCaseRecipe(recipes []gitutil.RemediationRecipe) bool {
	for _, r := range recipes {
		if strings.Contains(strings.ToLower(r.Title), "case collision") {
			return true
		}
	}

	return false
}

func buildCaseRecipes(repoPath string, cols []gitutil.CaseCollision) []gitutil.RemediationRecipe {
	var recipes []gitutil.RemediationRecipe
	for _, col := range cols {
		recipes = append(recipes, gitutil.GenerateCaseNormalizationRecipe(repoPath, col))
	}

	return recipes
}

func applyFallbackOrInvalid(item *cmdremediation.RemediationItem, action string) error {
	if len(item.Recipes) > 0 {
		return executeFixRecipe(item, item.Recipes[0])
	}
	return apperror.New("fix", "E_INVALID_OPTION", map[string]any{
		"msg": fmt.Sprintf("Invalid fix option: %s", action),
	})
}

func resolveFixTarget(args []string, aliasOverride string, items []cmdremediation.RemediationItem) (*cmdremediation.RemediationItem, string, error) {
	repoQuery, action := cmdreconcile.ParseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}

	if isFixActionWord(repoQuery) && (action == "stash" || action == "") {
		action = repoQuery
		repoQuery = ""
	}

	if repoQuery != "" {
		return findItemOrError(items, repoQuery, action)
	}

	if len(items) == 1 {
		return &items[0], action, nil
	}

	cmdremediation.PrintRemediationSummary(items)
	return nil, "", apperror.New("fix", "E_AMBIGUOUS", map[string]any{
		"msg": "Multiple repositories need remediation. Specify repo: gitmap fix <repo> <action>",
	})
}

func parseRecipeIndex(option string, recipes []gitutil.RemediationRecipe) int {
	low := strings.ToLower(strings.TrimSpace(option))
	for i, r := range recipes {
		t := strings.ToLower(r.Title)
		if (low == "case" || low == "casing") && strings.Contains(t, "case") {
			return i
		}
		if (low == "init" || low == "i") && strings.Contains(t, "init") {
			return i
		}
		if (low == "clone" || low == "c") && strings.Contains(t, "clone") {
			return i
		}
		if (low == "stash" || low == "s") && strings.Contains(t, "stash") {
			return i
		}
		if (low == "wip" || low == "w") && strings.Contains(t, "wip") {
			return i
		}
		if (low == "discard" || low == "d") && strings.Contains(t, "discard") {
			return i
		}
		if (low == "rm" || low == "remove") && strings.Contains(t, "remove") {
			return i
		}
	}

	switch low {
	case "1", "stash", "s", "clone", "c":
		return 0
	case "2", "wip", "w", "init", "i", "link":
		return 1
	case "3", "discard", "clean", "d", "rm", "remove":
		return 2
	}

	if i, err := strconv.Atoi(low); err == nil && i > 0 && i <= len(recipes) {
		return i - 1
	}

	return -1
}

func findItemOrError(items []cmdremediation.RemediationItem, repoQuery, action string) (*cmdremediation.RemediationItem, string, error) {
	matched := cmdremediation.FindRemediationItem(items, repoQuery)
	if matched != nil {
		ensureCaseCollisionRecipe(matched)
		return matched, action, nil
	}
	handled, err := cmdpull.RemediateRepoByQuery(repoQuery)
	if handled {
		return nil, action, err
	}
	localItem, cleanMsg := resolveLocalOrCaseItem(repoQuery)
	if localItem != nil {
		return localItem, action, nil
	}
	if cleanMsg != "" {
		fmt.Println(cleanMsg)
		return nil, action, nil
	}

	return resolveNonRepoOrError(items, repoQuery, action)
}

func resolveLocalOrCaseItem(repoQuery string) (*cmdremediation.RemediationItem, string) {
	item, cleanMsg := cmdremediation.InspectLocalRepoItem(repoQuery)
	if item != nil {
		ensureCaseCollisionRecipe(item)
		return item, ""
	}
	caseItem := inspectTargetCaseCollision(repoQuery)
	if caseItem != nil {
		return caseItem, ""
	}

	return nil, cleanMsg
}

func inspectTargetCaseCollision(repoQuery string) *cmdremediation.RemediationItem {
	targetPath := resolveRepoPathForFix(repoQuery)
	if targetPath == "" {
		return nil
	}
	cols, err := gitutil.DetectCaseCollisions(targetPath)
	if err != nil || len(cols) == 0 {
		return nil
	}
	item := &cmdremediation.RemediationItem{
		RepoPath:      targetPath,
		RepoName:      filepath.Base(targetPath),
		SummaryReason: "case collision in git index",
		Recipes:       buildCaseRecipes(targetPath, cols),
	}
	_ = cmdremediation.SaveRemediationState([]cmdremediation.RemediationItem{*item})

	return item
}

func resolveNonRepoOrError(items []cmdremediation.RemediationItem, repoQuery, action string) (*cmdremediation.RemediationItem, string, error) {
	nonRepoItem, err := resolveNonRepoRemediationItem(repoQuery, action)
	if err == nil && nonRepoItem != nil {
		return nonRepoItem, action, nil
	}

	return nil, "", buildFixNotFoundError(items, repoQuery)
}

func buildFixNotFoundError(items []cmdremediation.RemediationItem, repoQuery string) error {
	suggestions := cmdremediation.FindRemediationSuggestions(items, repoQuery)
	msg := fmt.Sprintf("Repository %q not found in pending remediation list.", repoQuery)
	if len(suggestions) > 0 {
		msg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
	}

	return apperror.New("fix", "E_NOT_FOUND", map[string]any{"msg": msg})
}

func runFixDirect(args []string, aliasOverride string) error {
	repoQuery, action := cmdreconcile.ParseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}

	if isFixActionWord(repoQuery) && (action == "stash" || action == "") {
		action = repoQuery
		repoQuery = ""
	}

	targetPath := resolveRepoPathForFix(repoQuery)
	if targetPath == "" && repoQuery == "" {
		targetPath = "."
	}
	handled, err := tryExecuteTargetedFixAction(targetPath, action)
	if handled {
		return err
	}

	handled, pullErr := cmdpull.RemediateRepoByQuery(repoQuery)
	if handled {
		return pullErr
	}

	localItem, cleanMsg := resolveLocalOrCaseItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)
		return nil
	}
	if localItem != nil {
		return applyFixRecipe(localItem, action)
	}

	return fixNonRepoDirectOrError(repoQuery, action)
}

func fixNonRepoDirectOrError(repoQuery, action string) error {
	nonRepoItem, err := resolveNonRepoRemediationItem(repoQuery, action)
	if err == nil && nonRepoItem != nil {
		return applyFixRecipe(nonRepoItem, action)
	}
	return apperror.NewNotFoundError(fmt.Sprintf("repository %q not found or not a git repository", repoQuery))
}
