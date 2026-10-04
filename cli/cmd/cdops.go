package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runCDLookup finds a repo by name and prints its path to stdout.
func runCDLookup(name string, args []string) error {
	if HasAlias() {
		fmt.Print(GetAliasPath())

		return nil
	}

	pick, yes, rest := parseCDPickFlag(args)
	records, err := lookupCDRecords(name)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return handleWorkDirOrNotFound(name, rest)
	}

	return dispatchCDPath(name, records, rest, pick, yes)
}

func dispatchCDPath(name string, records []model.ScanRecord, rest []string, pick, yes bool) error {
	path, err := resolveCDPath(name, records, pick, yes)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return runCDInner(path, rest)
	}

	fmt.Print(path)
	WriteShellHandoff(path)
	warnIfNoWrapper()

	return nil
}

// runCDInner chdirs into path and dispatches the inner subcommand.
func runCDInner(path string, innerArgs []string) error {
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCDChdirFmt, path, err)

		return fmt.Errorf(constants.ErrCDChdirFmt, path, err)
	}

	os.Args = append([]string{os.Args[0]}, innerArgs...)
	dispatch(innerArgs[0])

	return nil
}

// parseCDPickFlag extracts --pick and --yes flags plus any remaining positional args.
func parseCDPickFlag(args []string) (bool, bool, []string) {
	fs := flag.NewFlagSet("cd-lookup", flag.ContinueOnError)
	pick := fs.Bool("pick", false, constants.FlagDescCDPick)
	fs.BoolVar(pick, "p", false, constants.FlagDescCDPick)
	yes := fs.Bool("yes", false, "Auto-select first location without prompting")
	fs.BoolVar(yes, "y", false, "Auto-select first location without prompting")
	_ = fs.Parse(args)

	return *pick, *yes, fs.Args()
}

// lookupCDRecords finds repos matching the given name via DB.
func lookupCDRecords(name string) ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return nil, fmt.Errorf(constants.ErrListDBFailed, err)
	}

	defer db.Close()

	return findCDRecords(db, name), nil
}

func findCDRecords(db *store.DB, name string) []model.ScanRecord {
	cleanName := strings.TrimRight(name, "/\\")
	repos, err := db.FindBySlug(strings.ToLower(cleanName))
	hasValidRepos := err == nil && len(repos) > 0
	if hasValidRepos {
		return deduplicateCDRecords(repos)
	}

	all, listErr := db.ListRepos()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)
	}

	return deduplicateCDRecords(findBySlug(all, cleanName))
}

func deduplicateCDRecords(records []model.ScanRecord) []model.ScanRecord {
	if len(records) <= 1 {
		return records
	}

	seen := make(map[string]bool)
	var deduped []model.ScanRecord
	for _, r := range records {
		norm := filepath.Clean(filepath.FromSlash(r.AbsolutePath))
		if runtime.GOOS == "windows" {
			norm = strings.ToLower(norm)
		}
		norm = strings.TrimRight(norm, "\\/")
		if !seen[norm] {
			seen[norm] = true
			deduped = append(deduped, r)
		}
	}
	return deduped
}

// resolveCDPath picks the correct path from matches.
func resolveCDPath(name string, records []model.ScanRecord, pick, yes bool) (string, error) {
	records = deduplicateCDRecords(records)
	if len(records) == 1 {
		return records[0].AbsolutePath, nil
	}

	exact := filterExactMatches(records, name)
	if len(exact) == 1 {
		return exact[0].AbsolutePath, nil
	}
	if len(exact) > 1 {
		records = exact
	}

	if yes || isAutoPickActive() {
		return records[0].AbsolutePath, nil
	}

	if pick {
		return promptCDPick(name, records)
	}

	dflt := loadCDDefault(name)
	if len(dflt) > 0 {
		return dflt, nil
	}

	return promptCDPick(name, records)
}

func isAutoPickActive() bool {
	return os.Getenv("GITMAP_YES") == "1" || os.Getenv("GITMAP_NON_INTERACTIVE") == "1"
}

func filterExactMatches(records []model.ScanRecord, name string) []model.ScanRecord {
	clean := strings.ToLower(strings.TrimRight(name, "/\\"))
	var exact []model.ScanRecord
	for _, r := range records {
		if isExactMatch(r, clean) {
			exact = append(exact, r)
		}
	}
	return exact
}

func isExactMatch(r model.ScanRecord, clean string) bool {
	base := strings.ToLower(filepath.Base(filepath.Clean(filepath.FromSlash(r.AbsolutePath))))
	isMatchBase := base == clean
	isMatchSlug := strings.ToLower(r.Slug) == clean
	isMatchRepo := strings.ToLower(r.RepoName) == clean
	return isMatchBase || isMatchSlug || isMatchRepo
}

// promptCDPick shows a numbered list and reads user selection.
func promptCDPick(name string, records []model.ScanRecord) (string, error) {
	if len(records) == 0 {
		return "", fmt.Errorf(constants.ErrCDNotFound, name)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDMultipleHeader, name)

	for i, r := range records {
		fmt.Fprintf(os.Stderr, constants.MsgCDMultipleRowFmt, i+1, r.AbsolutePath)
	}

	fmt.Fprintf(os.Stderr, "\nPick [1-%d] (default 1): ", len(records))

	return readCDSelection(records, name)
}

// readCDSelection reads and validates the user's numeric choice.
func readCDSelection(records []model.ScanRecord, name ...string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if len(records) > 0 {
			fmt.Fprintf(os.Stderr, "  (auto-selected default: %s)\n", records[0].AbsolutePath)
			return records[0].AbsolutePath, nil
		}
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)
		return "", fmt.Errorf("%s", constants.ErrCDInvalidPick)
	}

	text := strings.TrimSpace(scanner.Text())
	if text == "" {
		selected := records[0].AbsolutePath
		saveCDDefaultChoice(name, selected)
		return selected, nil
	}

	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(records) {
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)
		return "", fmt.Errorf("%s", constants.ErrCDInvalidPick)
	}

	selected := records[idx-1].AbsolutePath
	saveCDDefaultChoice(name, selected)
	return selected, nil
}

func saveCDDefaultChoice(name []string, path string) {
	if len(name) == 0 {
		return
	}
	cleanName := strings.TrimSpace(name[0])
	cleanPath := strings.TrimSpace(path)
	if cleanName == "" || cleanPath == "" {
		return
	}
	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)
	if defaults == nil {
		defaults = make(map[string]string)
	}
	defaults[cleanName] = cleanPath
	_ = store.SaveCDDefaults(constants.DefaultOutputFolder, defaults)
}

// runCDRepos shows an interactive numbered list of all repos.
func runCDRepos(args []string) error {
	groupFilter := parseCDReposFlags(args)
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)

		return fmt.Errorf(constants.ErrListDBFailed, err)
	}

	defer db.Close()

	records := loadCDReposList(db, groupFilter)
	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgListEmpty)

		return fmt.Errorf("%s", constants.MsgListEmpty)
	}

	return executeCDReposPick(records)
}

func executeCDReposPick(records []model.ScanRecord) error {
	path, pickErr := promptCDReposPick(records)
	if pickErr != nil {
		return pickErr
	}

	fmt.Print(path)
	WriteShellHandoff(path)

	return nil
}

// parseCDReposFlags parses the --group flag for the repos subcommand.
func parseCDReposFlags(args []string) string {
	fs := flag.NewFlagSet("cd-repos", flag.ContinueOnError)
	group := fs.String("group", "", constants.FlagDescCDGroup)
	fs.StringVar(group, "g", "", constants.FlagDescCDGroup)
	_ = fs.Parse(args)

	return *group
}

// loadCDReposList loads repos optionally filtered by group.
func loadCDReposList(db *store.DB, group string) []model.ScanRecord {
	hasGroup := len(group) > 0
	if hasGroup {
		return loadCDGroupRepos(db, group)
	}

	repos, listErr := db.ListRepos()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)
	}

	return repos
}

func loadCDGroupRepos(db *store.DB, group string) []model.ScanRecord {
	repos, grpErr := db.ShowGroup(group)
	if grpErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load group %s: %v\n", group, grpErr)
	}

	return repos
}

// promptCDReposPick shows all repos and reads user selection.
func promptCDReposPick(records []model.ScanRecord) (string, error) {
	fmt.Fprint(os.Stderr, constants.MsgCDReposHeader)

	for i, r := range records {
		fmt.Fprintf(os.Stderr, constants.MsgCDReposRowFmt, i+1, r.RepoName)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))

	return readCDSelection(records)
}
