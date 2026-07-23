package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runVersionHistory handles the "version-history" subcommand.
func runVersionHistory(args []string) {
	checkHelp("version-history", args)
	limit, jsonOut := parseVersionHistoryFlags(args)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVersionHistoryCwd, err)
		os.Exit(1)
	}

	absPath := resolveVersionHistoryPath(cwd)
	records := loadVersionHistory(absPath, limit)

	if jsonOut {
		printVersionHistoryJSON(records)

		return
	}

	printVersionHistoryTerminal(records, absPath)
}

// parseVersionHistoryFlags parses --limit and --json flags.
func parseVersionHistoryFlags(args []string) (int, bool) {
	fs := flag.NewFlagSet(constants.CmdVersionHistory, flag.ExitOnError)
	limit := fs.Int("limit", 0, constants.FlagDescLimit)
	jsonFlag := fs.Bool("json", false, constants.FlagDescLBJSON)
	fs.Parse(args)

	return *limit, *jsonFlag
}

// resolveVersionHistoryPath resolves the repo path for version history lookup.
func resolveVersionHistoryPath(cwd string) string {
	remoteURL, err := gitutil.RemoteURL(cwd)
	if err != nil {
		return cwd
	}

	repoName := extractRepoName(remoteURL)
	parsed := clonenext.ParseRepoName(repoName)

	return filepath.Join(filepath.Dir(cwd), parsed.BaseName)
}

// loadVersionHistory fetches version history from the database.
func loadVersionHistory(absPath string, limit int) []model.RepoVersionHistoryRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVersionHistoryDB, err)
		os.Exit(1)
	}
	defer db.Close()

	repoID, findErr := db.GetRepoIDByPath(absPath)
	if findErr != nil {
		fmt.Print(constants.MsgVersionHistoryEmpty)
		os.Exit(0)
	}

	records, queryErr := db.ListVersionHistory(repoID)
	if queryErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrVersionHistoryDB, queryErr)
		os.Exit(1)
	}

	if limit > 0 && limit < len(records) {
		records = records[:limit]
	}

	return records
}

// printVersionHistoryTerminal prints version history in table format.
func printVersionHistoryTerminal(records []model.RepoVersionHistoryRecord, absPath string) {
	if len(records) == 0 {
		fmt.Print(constants.MsgVersionHistoryEmpty)

		return
	}

	fmt.Printf(constants.MsgVersionHistoryHeader, absPath)
	fmt.Println(constants.MsgVersionHistoryColumns)
	for _, r := range records {
		fmt.Printf(constants.MsgVersionHistoryRowFmt,
			r.FromVersionTag, r.ToVersionTag, r.FlattenedPath, r.CreatedAt)
	}
	fmt.Printf(constants.MsgVersionHistoryCount, len(records))
}

// printVersionHistoryJSON outputs version history as JSON.
func printVersionHistoryJSON(records []model.RepoVersionHistoryRecord) {
	if len(records) == 0 {
		fmt.Println("[]")

		return
	}

	if err := encodeVersionHistoryJSON(os.Stdout, records); err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to encode version history to JSON: %v\n", err)
	}
}
