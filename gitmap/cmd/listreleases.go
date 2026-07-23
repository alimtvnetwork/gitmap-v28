package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runListReleases handles the "list-releases" command.
func runListReleases(args []string) {
	checkHelp("list-releases", args)
	asJSON := hasListReleasesJSONFlag(args)
	limit := parseListReleasesLimit(args)

	// New in v3.20.0: --all-repos pivots to the multi-repo joined view.
	if hasAllReposFlag(args) {
		runListReleasesAllRepos(asJSON, limit)

		return
	}

	source := parseListReleasesSource(args)
	releases := loadReleases()
	releases = filterBySource(releases, source)
	releases = applyReleaseLimit(releases, limit)

	if asJSON {
		printReleasesJSON(releases)

		return
	}

	printReleasesTerminal(releases)
}

// parseListReleasesSource extracts the --source value from args.
func parseListReleasesSource(args []string) string {
	for i, arg := range args {
		if arg == constants.FlagSource && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// filterBySource keeps only releases matching the given source (empty = all).
func filterBySource(releases []model.ReleaseRecord, source string) []model.ReleaseRecord {
	if source == "" {
		return releases
	}

	var filtered []model.ReleaseRecord
	for _, r := range releases {
		if r.Source == source {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

// hasListReleasesJSONFlag checks if --json is present in args.
func hasListReleasesJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagJSON {
			return true
		}
	}

	return false
}

// parseListReleasesLimit extracts the --limit N value from args.
func parseListReleasesLimit(args []string) int {
	for i, arg := range args {
		if arg == constants.FlagLimit && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil && n > 0 {
				return n
			}
		}
	}

	return 0
}

// applyReleaseLimit trims releases to at most n items (0 means no limit).
func applyReleaseLimit(releases []model.ReleaseRecord, n int) []model.ReleaseRecord {
	if n <= 0 || n >= len(releases) {
		return releases
	}

	return releases[:n]
}

// loadReleases builds a unified release list from repo metadata, git tags,
// and the database. Results are cached to the DB on every invocation.
func loadReleases() []model.ReleaseRecord {
	records := loadReleasesFromRepo()
	tagRecords := loadReleasesFromTags(records)
	records = append(records, tagRecords...)

	if len(records) == 0 {
		records = loadReleasesFromDB()
	}

	sortRecordsByDate(records)
	cacheReleasesToDB(records)

	return records
}

// printReleasesTerminal renders releases as a table to stdout.
func printReleasesTerminal(releases []model.ReleaseRecord) {
	if len(releases) == 0 {
		fmt.Println(constants.MsgListReleasesEmpty)

		return
	}

	fmt.Printf(constants.MsgListReleasesHeader, len(releases))
	fmt.Println(constants.MsgListReleasesSeparator)
	fmt.Println(constants.MsgListReleasesColumns)
	for _, r := range releases {
		printReleaseRow(r)
	}
}

// printReleaseRow prints a single release row.
func printReleaseRow(r model.ReleaseRecord) {
	draft := constants.MsgNo
	if r.IsDraft {
		draft = constants.MsgYes
	}
	latest := constants.MsgNo
	if r.IsLatest {
		latest = constants.MsgYes
	}

	fmt.Printf(constants.MsgListReleasesRowFmt, r.Version, r.Tag, r.Branch, draft, latest, r.Source, r.CreatedAt)
}

// printReleasesJSON renders releases as JSON to stdout via the
// stablejson encoder so key order is contractual (pinned by
// spec/08-json-schemas/list-releases.schema.json + the matching
// jsonschema_contract test). Keeps the legacy `_ = constants.JSONIndent`
// 2-space layout byte-for-byte compatible with the prior MarshalIndent
// output, so existing consumers see no diff.
func printReleasesJSON(releases []model.ReleaseRecord) {
	if err := encodeListReleasesJSON(os.Stdout, releases); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal releases to JSON: %v\n", err)
	}
}
