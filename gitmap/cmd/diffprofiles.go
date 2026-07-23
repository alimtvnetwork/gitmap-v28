package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runDiffProfiles handles the "diff-profiles" command.
func runDiffProfiles(args []string) {
	checkHelp("diff-profiles", args)
	nameA, nameB, showAll, jsonMode := parseDPFlags(args)
	validateDPProfiles(nameA, nameB)

	reposA := loadProfileRepos(nameA)
	reposB := loadProfileRepos(nameB)

	result := compareDPRepos(reposA, reposB)

	if jsonMode {
		printDPJSON(nameA, nameB, result)

		return
	}

	printDPOutput(nameA, nameB, result, showAll)
}

// parseDPFlags parses flags for the diff-profiles command.
func parseDPFlags(args []string) (string, string, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdDiffProfiles, flag.ExitOnError)
	allFlag := fs.Bool("all", false, "Include identical repos")
	jsonFlag := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprint(os.Stderr, constants.ErrDPUsage)
		os.Exit(1)
	}

	return fs.Arg(0), fs.Arg(1), *allFlag, *jsonFlag
}

// validateDPProfiles checks both profiles exist.
func validateDPProfiles(nameA, nameB string) {
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	for _, name := range []string{nameA, nameB} {
		if !profileExists(cfg.Profiles, name) {
			fmt.Fprintf(os.Stderr, constants.ErrDPProfileMissing, name)
			os.Exit(1)
		}
	}
}

// loadProfileRepos opens a profile's DB and returns all repos.
func loadProfileRepos(name string) []model.ScanRecord {
	db, err := store.OpenDefaultProfile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDPOpenFailed, name, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ DB migration failed: %v\n", err)
	}

	repos, err := db.ListRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDPOpenFailed, name, err)
		os.Exit(1)
	}

	return repos
}

// printDPJSON outputs the comparison result as JSON.
func printDPJSON(nameA, nameB string, result dpResult) {
	if err := encodeDiffProfilesJSON(os.Stdout, nameA, nameB, result); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode diff result to JSON: %v\n", err)
	}
}

// dpRepoSummaries converts records to simple name+path maps.
func dpRepoSummaries(records []model.ScanRecord) []map[string]string {
	result := make([]map[string]string, 0, len(records))

	for _, r := range records {
		result = append(result, map[string]string{
			"name": r.RepoName,
			"path": r.AbsolutePath,
		})
	}

	return result
}
