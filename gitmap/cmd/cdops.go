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
// If extra positional args follow the name, they are treated as a
// gitmap subcommand to execute inside the resolved repo directory
// (e.g. `gitmap cd myrepo cn v++` clones the next version of myrepo
// and hands off to the new folder via the inner command's handoff).
func runCDLookup(name string, args []string) {
	if HasAlias() {
		fmt.Print(GetAliasPath())

		return
	}

	pick, rest := parseCDPickFlag(args)
	records := lookupCDRecords(name)

	if len(records) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrCDNotFound, name)
		os.Exit(1)
	}

	path := resolveCDPath(name, records, pick)

	if len(rest) > 0 {
		runCDInner(path, rest)

		return
	}

	fmt.Print(path)
	WriteShellHandoff(path)
	warnIfNoWrapper()
}

// runCDInner chdirs into path and dispatches the inner subcommand
// (the args after `gitmap cd <name>`). The inner command is
// responsible for writing its own shell handoff if it relocates the
// caller — `cn`, `cfr`, `cfrp`, etc. already do this.
func runCDInner(path string, innerArgs []string) {
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCDChdirFmt, path, err)
		os.Exit(1)
	}

	os.Args = append([]string{os.Args[0]}, innerArgs...)
	dispatch(innerArgs[0])
}

// parseCDPickFlag extracts --pick and returns it plus any remaining
// positional args (the inner subcommand + its args, if any).
func parseCDPickFlag(args []string) (bool, []string) {
	fs := flag.NewFlagSet("cd-lookup", flag.ContinueOnError)
	pick := fs.Bool("pick", false, constants.FlagDescCDPick)
	_ = fs.Parse(args)

	return *pick, fs.Args()
}

// lookupCDRecords finds repos matching the given name via DB.
func lookupCDRecords(name string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.FindBySlug(strings.ToLower(name))
	if err != nil || len(repos) == 0 {
		all, listErr := db.ListRepos()
		if listErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)
		}

		return findBySlug(all, name)
	}

	return repos
}

// resolveCDPath picks the correct path from matches.
func resolveCDPath(name string, records []model.ScanRecord, pick bool) string {
	if len(records) == 1 {
		return records[0].AbsolutePath
	}

	dflt := loadCDDefault(name)
	if len(dflt) > 0 && !pick {
		return dflt
	}

	return promptCDPick(name, records)
}

// promptCDPick shows a numbered list and reads user selection.
func promptCDPick(name string, records []model.ScanRecord) string {
	fmt.Fprintf(os.Stderr, constants.MsgCDMultipleHeader, name)

	for i, r := range records {
		fmt.Fprintf(os.Stderr, constants.MsgCDMultipleRowFmt, i+1, r.AbsolutePath)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))

	return readCDSelection(records)
}

// readCDSelection reads and validates the user's numeric choice.
func readCDSelection(records []model.ScanRecord) string {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)
		os.Exit(1)
	}

	idx, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || idx < 1 || idx > len(records) {
		fmt.Fprint(os.Stderr, constants.ErrCDInvalidPick)
		os.Exit(1)
	}

	return records[idx-1].AbsolutePath
}

// runCDRepos shows an interactive numbered list of all repos.
func runCDRepos(args []string) {
	groupFilter := parseCDReposFlags(args)
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	records := loadCDReposList(db, groupFilter)
	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgListEmpty)
		os.Exit(1)
	}

	path := promptCDReposPick(records)
	fmt.Print(path)
	WriteShellHandoff(path)
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
	if len(group) > 0 {
		repos, grpErr := db.ShowGroup(group)
		if grpErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not load group %s: %v\n", group, grpErr)
		}

		return repos
	}

	repos, listErr := db.ListRepos()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not list repos: %v\n", listErr)
	}

	return repos
}

// promptCDReposPick shows all repos and reads user selection.
func promptCDReposPick(records []model.ScanRecord) string {
	fmt.Fprint(os.Stderr, constants.MsgCDReposHeader)

	for i, r := range records {
		fmt.Fprintf(os.Stderr, constants.MsgCDReposRowFmt, i+1, r.RepoName)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))

	return readCDSelection(records)
}
