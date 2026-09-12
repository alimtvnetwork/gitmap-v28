package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runCDLookup finds a repo by name and prints its path to stdout.
func runCDLookup(name string, args []string) error {
	if HasAlias() {
		fmt.Print(GetAliasPath())

		return nil
	}

	pick, rest := parseCDPickFlag(args)
	records, err := lookupCDRecords(name)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return handleWorkDirOrNotFound(name, rest)
	}

	return dispatchCDPath(name, records, rest, pick)
}

func dispatchCDPath(name string, records []model.ScanRecord, rest []string, pick bool) error {
	path, err := resolveCDPath(name, records, pick)
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

// parseCDPickFlag extracts --pick and returns it plus any remaining positional args.
func parseCDPickFlag(args []string) (bool, []string) {
	fs := flag.NewFlagSet("cd-lookup", flag.ContinueOnError)
	pick := fs.Bool("pick", false, constants.FlagDescCDPick)
	_ = fs.Parse(args)

	return *pick, fs.Args()
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
		return repos
	}

	all, listErr := db.ListRepos()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)
	}

	return findBySlug(all, cleanName)
}

// resolveCDPath picks the correct path from matches.
func resolveCDPath(name string, records []model.ScanRecord, pick bool) (string, error) {
	if len(records) == 1 {
		return records[0].AbsolutePath, nil
	}

	dflt := loadCDDefault(name)
	if len(dflt) > 0 && !pick {
		return dflt, nil
	}

	return promptCDPick(name, records)
}

// promptCDPick shows a numbered list and reads user selection.
func promptCDPick(name string, records []model.ScanRecord) (string, error) {
	fmt.Fprintf(os.Stderr, constants.MsgCDMultipleHeader, name)

	for i, r := range records {
		fmt.Fprintf(os.Stderr, constants.MsgCDMultipleRowFmt, i+1, r.AbsolutePath)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))

	return readCDSelection(records)
}

// readCDSelection reads and validates the user's numeric choice.
func readCDSelection(records []model.ScanRecord) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)

		return "", fmt.Errorf("%s", constants.ErrCDInvalidPick)
	}

	idx, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || idx < 1 || idx > len(records) {
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)

		return "", fmt.Errorf("%s", constants.ErrCDInvalidPick)
	}

	return records[idx-1].AbsolutePath, nil
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
