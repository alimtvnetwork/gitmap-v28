package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runListReleasesAllRepos executes the multi-repo batch view that reads
// EVERY Release row joined with its owning Repo. Triggered by
// `gitmap releases --all-repos` (or `gitmap lr --all-repos`).
func runListReleasesAllRepos(asJSON bool, limit int) {
	records := loadReleasesAcrossRepos()
	records = applyAllReposLimit(records, limit)

	if asJSON {
		printAllReposJSON(records)

		return
	}
	printAllReposTerminal(records)
}

// loadReleasesAcrossRepos opens the DB and queries the joined view.
// Errors are written to stderr and the process exits non-zero.
func loadReleasesAcrossRepos() []store.ReleaseAcrossRepos {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrNoDatabase)
		os.Exit(1)
	}
	defer db.Close()

	records, err := db.ListReleasesAcrossRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListReleasesFailed, err)
		os.Exit(1)
	}

	return records
}

// applyAllReposLimit trims the slice to at most n records (0 = no limit).
func applyAllReposLimit(records []store.ReleaseAcrossRepos, n int) []store.ReleaseAcrossRepos {
	if n <= 0 || n >= len(records) {
		return records
	}

	return records[:n]
}

// printAllReposTerminal renders the joined view as a wide table.
func printAllReposTerminal(records []store.ReleaseAcrossRepos) {
	if len(records) == 0 {
		fmt.Println(constants.MsgListReleasesAllReposEmpty)

		return
	}

	fmt.Printf(constants.MsgListReleasesAllReposHeader, len(records))
	fmt.Println(constants.MsgListReleasesAllReposSeparator)
	fmt.Println(constants.MsgListReleasesAllReposColumns)
	for _, r := range records {
		printAllReposRow(r)
	}
}

// printAllReposRow prints a single joined-row line.
func printAllReposRow(r store.ReleaseAcrossRepos) {
	latest := constants.MsgNo
	if r.IsLatest {
		latest = constants.MsgYes
	}
	fmt.Printf(constants.MsgListReleasesAllReposRowFmt,
		r.RepoSlug, r.Version, r.Tag, r.Branch, latest, r.Source, r.CreatedAt)
}

// printAllReposJSON renders the joined --all-repos view via the
// stablejson encoder. PascalCase wire keys preserve the legacy
// MarshalIndent surface (store.ReleaseAcrossRepos has no `json:`
// tags). Schema:
// spec/08-json-schemas/list-releases-all-repos.schema.json.
func printAllReposJSON(records []store.ReleaseAcrossRepos) {
	if err := encodeListReleasesAllReposJSON(os.Stdout, records); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal releases to JSON: %v\n", err)
	}
}

// hasAllReposFlag reports whether --all-repos appears in args.
func hasAllReposFlag(args []string) bool {
	for _, a := range args {
		if a == constants.FlagAllRepos {
			return true
		}
	}

	return false
}
