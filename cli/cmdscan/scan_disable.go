package cmdscan

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexclude"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func init() {
	hookScanDisableDispatch()
}

// ParseNumberedSelection forwards to cmdexclude.ParseNumberedSelection.
func ParseNumberedSelection(input string, maxCount int) ([]int, error) {
	return cmdexclude.ParseNumberedSelection(input, maxCount)
}

// RunScanDisableRepos handles disabling scanned repositories.
func RunScanDisableRepos(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return renderInteractiveRepoTable(true)
	}

	return executeScanSelection(args, true)
}

// RunScanEnableRepos handles re-enabling scanned repositories.
func RunScanEnableRepos(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return renderInteractiveRepoTable(false)
	}

	return executeScanSelection(args, false)
}

func extractFlagsAndSelection(args []string) (string, bool) {
	var tokens []string
	isJSON := false

	for _, a := range args {
		clean := strings.TrimSpace(a)
		if clean == "--json" || clean == "-j" {
			isJSON = true
			continue
		}
		tokens = append(tokens, clean)
	}

	return strings.Join(tokens, " "), isJSON
}

func executeScanSelection(args []string, isDisable bool) error {
	selection, isJSON := extractFlagsAndSelection(args)
	res, err := cmdexclude.ToggleRepoExclusions(selection, isDisable)
	if err != nil {
		return err
	}

	if isJSON {
		return json.NewEncoder(os.Stdout).Encode(res)
	}

	printConfirmationSummary(res, isDisable)

	return nil
}

func renderInteractiveRepoTable(isDisable bool) error {
	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "open_default_store")
	}

	defer db.Close()

	repos, listErr := db.ListRepos()
	if listErr != nil {
		return apperror.WrapSimple(listErr, "list_repos")
	}

	return displayAndPromptRepos(repos, isDisable)
}

func displayAndPromptRepos(repos []model.ScanRecord, isDisable bool) error {
	fmt.Printf("Scanning indexed repositories... Found %d repositories.\n\n", len(repos))
	for i, r := range repos {
		printRepoRow(i+1, r)
	}

	input, err := readInteractiveInput(isDisable)
	if err != nil || strings.TrimSpace(input) == "" {
		return err
	}

	return executeScanSelection([]string{input}, isDisable)
}

func printRepoRow(idx int, repo model.ScanRecord) {
	pill := resolveRepoPill(repo)
	slug := repo.Slug
	if slug == "" {
		slug = repo.RepoName
	}

	path := repo.RelativePath
	if path == "" {
		path = repo.AbsolutePath
	}

	fmt.Printf("  [%d] %s (%s) %s\n", idx, slug, path, pill)
}

func resolveRepoPill(repo model.ScanRecord) string {
	isExcl := cmdexclude.IsCustomExcluded(repo.RelativePath) || cmdexclude.IsCustomExcluded(repo.Slug)
	if isExcl {
		return "[DISABLED]"
	}

	return "[ENABLED]"
}

func readInteractiveInput(isDisable bool) (string, error) {
	action := "disable"
	if !isDisable {
		action = "enable"
	}

	fmt.Printf("\nEnter numbers/ranges to %s (e.g. 1-3, 6) or repo names: ", action)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", nil
}

func printConfirmationSummary(res *cmdexclude.ToggleResult, isDisable bool) {
	if isDisable {
		printDisableSummary(res)
		return
	}

	printEnableSummary(res)
}

func printDisableSummary(res *cmdexclude.ToggleResult) {
	fmt.Println("[GITMAP-SCAN-SETTINGS] Updating repository exclusions...")
	for _, repo := range res.AffectedRepos {
		fmt.Printf("  ✓ Disabled repository: %s\n", repo)
	}

	fmt.Printf("\nSaved %d exclusion pattern(s) to Split-DB setting 'custom_repo_exclusions'.\n", res.ModifiedCount)
	fmt.Printf("Total excluded repositories: %d\n", res.TotalExcluded)
}

func printEnableSummary(res *cmdexclude.ToggleResult) {
	fmt.Println("[GITMAP-SCAN-SETTINGS] Re-enabling repositories...")
	for _, repo := range res.AffectedRepos {
		fmt.Printf("  ✓ Re-enabled repository: %s\n", repo)
	}

	fmt.Printf("\nUpdated custom exclusions: %d removed, %d remaining excluded repositories.\n", res.ModifiedCount, res.TotalExcluded)
}

func hookScanDisableDispatch() {
	if len(os.Args) < 2 {
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(os.Args[1]))
	if dispatchRootScanAliases(cmd) {
		return
	}

	dispatchScanSubcommands(cmd)
}

func dispatchRootScanAliases(cmd string) bool {
	if cmd == "disable-scan" || cmd == "ds" {
		executeScanToggle(os.Args[2:], true)
		return true
	}

	if cmd == "enable-scan" || cmd == "es" {
		executeScanToggle(os.Args[2:], false)
		return true
	}

	return false
}

func dispatchScanSubcommands(cmd string) bool {
	isScan := cmd == "scan" || cmd == "s"
	if !isScan || len(os.Args) < 3 {
		return false
	}

	sub := strings.ToLower(strings.TrimSpace(os.Args[2]))

	return matchScanSubcommandVerb(sub)
}

func matchScanSubcommandVerb(sub string) bool {
	if sub == "disable-repos" || sub == "disable-scan" || sub == "disable" {
		executeScanToggle(os.Args[3:], true)
		return true
	}
	if sub == "enable-repos" || sub == "enable-scan" || sub == "enable" {
		executeScanToggle(os.Args[3:], false)
		return true
	}

	return false
}

func executeScanToggle(args []string, isDisable bool) {
	var err error
	if isDisable {
		err = RunScanDisableRepos(args)
	} else {
		err = RunScanEnableRepos(args)
	}

	if err != nil {
		cliexit.HandleError(err, int(cliexit.ExitCodeGeneralError))
		return
	}

	cliexit.Exit(int(cliexit.ExitCodeSuccess))
}
