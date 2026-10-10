package cmdfix

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdignore"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdreconcile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isFixAgyRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}

	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "agy" || low == "aef" || low == "agy-errors-fix" {
			return true
		}
	}

	return false
}

func isFixAllRequested(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "all" || low == "--all" || low == "-all" || low == "-a" {
			return true
		}
	}

	return false
}

func isFixPromptRequested(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "--prompt" || low == "-p" || low == "prompt" || low == "--interactive" || low == "-i" || low == "interactive" {
			return true
		}
	}

	return false
}

func isFixActionWord(s string) bool {
	switch s {
	case "1", "stash", "s", "clone", "c", "2", "wip", "w", "init", "i", "link", "3", "discard", "clean", "d", "rm", "remove":
		return true
	}

	return false
}

func resolveFixAllAction(args []string, aliasOverride string) string {
	if aliasOverride != "" {
		return aliasOverride
	}

	for _, a := range args {
		low := strings.ToLower(a)
		if isFixActionWord(low) {
			return low
		}
	}

	return "stash"
}

func runFixAll(action string, items []cmdremediation.RemediationItem) error {
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)

		return nil
	}

	fmt.Printf("%s Remediating %d repository(ies) with action: %s\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items), action)

	for i := range items {
		idx := parseRecipeIndex(action, items[i].Recipes)
		if idx < 0 || idx >= len(items[i].Recipes) {
			idx = 0
		}

		err := executeFixRecipe(&items[i], items[i].Recipes[idx])
		if err != nil {
			return err
		}
		if i < len(items)-1 {
			fmt.Println()
		}
	}

	fmt.Printf("\n%s All %d repository(ies) remediated successfully.\n",
		constants.ColorGreen+"✓"+constants.ColorReset, len(items))

	return nil
}

func RunFix(args []string, aliasOverride string) error {
	cmdName := "fix"
	if aliasOverride != "" {
		cmdName = aliasOverride
	}
	checkHelp(cmdName, args)

	if isFixAgyRequest(args) {
		return cmdagy.RunPipelineFixAgyCLI(args)
	}

	if isFixIgnoreRequest(args) {
		return runFixIgnoreDispatch(args)
	}

	if isFixLsRequest(args) {
		return runFixLs(args)
	}

	items := cmdremediation.LoadRemediationState()
	if isFixAllRequested(args) {
		action := resolveFixAllAction(args, aliasOverride)
		allItems := collectAllRemediationItems(items)
		return runFixAll(action, allItems)
	}

	if len(items) == 0 {
		return handleEmptyRemediationState(args, aliasOverride)
	}

	if isFixPromptRequested(args) {
		return cmdremediation.RunInteractiveRemediation(items)
	}

	if len(args) == 0 && aliasOverride == "" {
		cmdremediation.PrintRemediationSummary(items)
		return nil
	}

	return executeFixTarget(args, aliasOverride, items)
}

func executeFixTarget(args []string, aliasOverride string, items []cmdremediation.RemediationItem) error {
	item, action, err := resolveFixTarget(args, aliasOverride, items)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}

	return applyFixRecipe(item, action)
}

func handleEmptyRemediationState(args []string, aliasOverride string) error {
	if len(args) > 0 && !isFixAllRequested(args) && !isFixPromptRequested(args) {
		return runFixDirect(args, aliasOverride)
	}

	if handleLiveDiscoveredIssues() {
		return nil
	}

	printNoPendingRemediationHelp()

	return nil
}

func handleLiveDiscoveredIssues() bool {
	liveItems := scanTrackedReposForIssues()
	if len(liveItems) == 0 {
		return false
	}

	_ = cmdremediation.SaveRemediationState(liveItems)
	cmdremediation.PrintRemediationSummary(liveItems)

	return true
}

func printNoPendingRemediationHelp() {
	fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
	fmt.Println("  Run 'gitmap fix ls' to inspect repositories with issues.")
	fmt.Println("  Run 'gitmap pull' to pull all tracked repositories.")
}

func runFixDirect(args []string, aliasOverride string) error {
	repoQuery, action := cmdreconcile.ParseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
	}

	item, cleanMsg := cmdremediation.InspectLocalRepoItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)

		return nil
	}
	if item == nil {
		nonRepoItem, err := resolveNonRepoRemediationItem(repoQuery, action)
		if err == nil && nonRepoItem != nil {
			return applyFixRecipe(nonRepoItem, action)
		}
		return apperror.NewNotFoundError(fmt.Sprintf("repository %q not found or not a git repository", repoQuery))
	}

	return applyFixRecipe(item, action)
}

func applyFixRecipe(item *cmdremediation.RemediationItem, action string) error {
	if action == "" && len(item.Recipes) > 0 {
		return executeFixRecipe(item, item.Recipes[0])
	}
	idx := parseRecipeIndex(action, item.Recipes)
	if idx < 0 || idx >= len(item.Recipes) {
		if len(item.Recipes) > 0 {
			return executeFixRecipe(item, item.Recipes[0])
		}
		return apperror.New("fix", "E_INVALID_OPTION", map[string]any{
			"msg": fmt.Sprintf("Invalid fix option: %s", action),
		})
	}

	return executeFixRecipe(item, item.Recipes[idx])
}

func resolveFixTarget(args []string, aliasOverride string, items []cmdremediation.RemediationItem) (*cmdremediation.RemediationItem, string, error) {
	repoQuery, action := cmdreconcile.ParseReconcileArgs(args)
	if aliasOverride != "" {
		action = aliasOverride
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
		return matched, action, nil
	}

	localItem, cleanMsg := cmdremediation.InspectLocalRepoItem(repoQuery)
	if cleanMsg != "" {
		fmt.Println(cleanMsg)

		return nil, action, nil
	}
	if localItem != nil {
		return localItem, action, nil
	}

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

func isFixIgnoreRequest(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if isIgnoreRoutingKeyword(low) {
			return true
		}
	}
	return false
}

func isIgnoreRoutingKeyword(low string) bool {
	switch low {
	case "ignore", "ignores", "gitignore", "fia", "fias", "fix-ignore-all", "fix-ignores-all":
		return true
	}
	return false
}

func isSSHArgKeyword(low string) bool {
	if low == "ssh" || low == "--ssh" || low == "fias" {
		return true
	}
	return false
}

func isIgnoreFilterKeyword(low string) bool {
	switch low {
	case "ignore", "ignores", "gitignore", "all", "fia", "fias", "ssh", "--ssh":
		return true
	}
	return false
}

func filterIgnoreArgs(args []string) []string {
	var remaining []string
	for _, a := range args {
		low := strings.ToLower(a)
		if isIgnoreFilterKeyword(low) {
			continue
		}
		remaining = append(remaining, a)
	}
	return remaining
}

func isSSHRequested(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if isSSHArgKeyword(low) {
			return true
		}
	}
	return false
}

func runFixIgnoreDispatch(args []string) error {
	hasSSH := isSSHRequested(args)
	remaining := filterIgnoreArgs(args)
	if hasSSH {
		return runFixIgnoreSSH(remaining)
	}
	return runFixIgnoreLocal(remaining)
}

func runFixIgnoreSSH(args []string) error {
	appErr := cmdignore.RunFixIgnoresAllSSH(args)
	if appErr != nil {
		return appErr
	}
	return nil
}

func runFixIgnoreLocal(args []string) error {
	appErr := cmdignore.RunFixIgnoreAll(args)
	if appErr != nil {
		return appErr
	}
	return nil
}

func resolveRepoPathForFix(repoQuery string) string {
	cleanQuery := strings.TrimSpace(repoQuery)
	if cleanQuery == "" {
		return ""
	}
	if info, err := os.Stat(cleanQuery); err == nil && info.IsDir() {
		return cleanQuery
	}
	// Check special repos split db
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, err := db.GetSpecialRepo(cleanQuery); err == nil && rec.LocalPath != "" {
			if info, err := os.Stat(rec.LocalPath); err == nil && info.IsDir() {
				return rec.LocalPath
			}
		}
	}
	// Check standard repos DB
	if db, err := store.OpenDefault(); err == nil {
		defer db.Close()
		if recs, err := db.FindBySlug(cleanQuery); err == nil && len(recs) > 0 {
			for _, r := range recs {
				if r.AbsolutePath != "" {
					if info, err := os.Stat(r.AbsolutePath); err == nil && info.IsDir() {
						return r.AbsolutePath
					}
				}
			}
		}
	}
	return cleanQuery
}

func resolveNonRepoRemediationItem(repoQuery, action string) (*cmdremediation.RemediationItem, error) {
	targetPath := resolveRepoPathForFix(repoQuery)
	if targetPath == "" {
		return nil, nil
	}

	diag := cloner.ClassifyNonRepoFolder(targetPath, repoQuery)
	if !diag.IsNonRepoFolder {
		return nil, nil
	}

	return buildNonRepoRemediationItem(diag, action)
}

func buildNonRepoRemediationItem(diag cloner.NonRepoDiagnosis, action string) (*cmdremediation.RemediationItem, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(diag.Path))
	var recipes []gitutil.RemediationRecipe

	if diag.HasRemote {
		entries, readErr := os.ReadDir(diag.Path)
		isEmpty := readErr == nil && len(entries) == 0

		var cloneSteps []gitutil.RemediationStep
		if isEmpty {
			cloneSteps = []gitutil.RemediationStep{
				{Name: "git", Args: []string{"clone", diag.RemoteURL, cleanPath}},
			}
		} else {
			cloneSteps = []gitutil.RemediationStep{
				{Name: "git", Args: []string{"-C", cleanPath, "init"}},
				{Name: "git", Args: []string{"-C", cleanPath, "remote", "add", "origin", diag.RemoteURL}},
				{Name: "git", Args: []string{"-C", cleanPath, "fetch", "origin"}},
			}
		}

		recipes = append(recipes, gitutil.RemediationRecipe{
			Title:       "Clone from Remote",
			Command:     diag.Option1,
			Description: fmt.Sprintf("Clone %s into %s", diag.RemoteURL, cleanPath),
			Steps:       cloneSteps,
		})

		linkSteps := []gitutil.RemediationStep{
			{Name: "git", Args: []string{"-C", cleanPath, "init"}},
			{Name: "git", Args: []string{"-C", cleanPath, "remote", "add", "origin", diag.RemoteURL}},
			{Name: "git", Args: []string{"-C", cleanPath, "fetch", "origin"}},
		}
		recipes = append(recipes, gitutil.RemediationRecipe{
			Title:       "Init & Link Remote",
			Command:     diag.Option2,
			Description: fmt.Sprintf("Initialize local git repo in %s and link origin to %s", cleanPath, diag.RemoteURL),
			Steps:       linkSteps,
		})
	} else {
		initSteps := []gitutil.RemediationStep{
			{Name: "git", Args: []string{"-C", cleanPath, "init"}},
		}
		recipes = append(recipes, gitutil.RemediationRecipe{
			Title:       "Initialize Local Repo",
			Command:     diag.Option1,
			Description: fmt.Sprintf("Initialize local git repository in %s", cleanPath),
			Steps:       initSteps,
		})

		removeSteps := []gitutil.RemediationStep{
			{Name: "gitmap", Args: []string{"rm", diag.RepoName, "--db-only"}},
		}
		recipes = append(recipes, gitutil.RemediationRecipe{
			Title:       "Remove from Registry",
			Command:     diag.Option2,
			Description: fmt.Sprintf("Remove %s from database registry", diag.RepoName),
			Steps:       removeSteps,
		})
	}

	return &cmdremediation.RemediationItem{
		RepoName:      diag.RepoName,
		RepoPath:      diag.Path,
		SummaryReason: diag.Reason,
		Recipes:       recipes,
	}, nil
}

func collectAllRemediationItems(existing []cmdremediation.RemediationItem) []cmdremediation.RemediationItem {
	seen := make(map[string]bool)
	var all []cmdremediation.RemediationItem
	for _, item := range existing {
		cleanPath := filepath.Clean(item.RepoPath)
		if !seen[cleanPath] {
			seen[cleanPath] = true
			all = append(all, item)
		}
	}

	liveItems := scanTrackedReposForIssues()
	for _, item := range liveItems {
		cleanPath := filepath.Clean(item.RepoPath)
		if !seen[cleanPath] {
			seen[cleanPath] = true
			all = append(all, item)
		}
	}

	specialKeys := []string{"repo-cache", "repo-secrets"}
	for _, key := range specialKeys {
		if path := resolveRepoPathForFix(key); path != "" {
			cleanPath := filepath.Clean(path)
			if !seen[cleanPath] {
				diag := cloner.ClassifyNonRepoFolder(path, key)
				if diag.IsNonRepoFolder {
					if item, err := buildNonRepoRemediationItem(diag, ""); err == nil && item != nil {
						seen[cleanPath] = true
						all = append(all, *item)
					}
				}
			}
		}
	}

	return all
}

