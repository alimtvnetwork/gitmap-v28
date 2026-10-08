package cmdlist

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// typeKeywords maps list filter keywords to project type keys.
var typeKeywords = map[string]string{
	"go":     constants.ProjectKeyGo,
	"node":   constants.ProjectKeyNode,
	"nodejs": constants.ProjectKeyNode,
	"react":  constants.ProjectKeyReact,
	"cpp":    constants.ProjectKeyCpp,
	"csharp": constants.ProjectKeyCsharp,
}

func RunList(args []string) error {
	checkHelp(constants.CmdList, args)

	if len(args) > 0 && isListTypeOrGroups(args[0]) {
		handleListSpecial(args[0], args[1:])

		return nil
	}

	executeList(args)

	return nil
}

type listOptions struct {
	group       string
	verbose     bool
	preview     bool
	tree        bool
	showNumbers bool
	showEmoji   bool
	useRelative bool
}

func executeList(args []string) {
	opts := parseListFlags(args)
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()
	records := fetchListRecords(db, opts.group)
	if opts.preview {
		previewOpts := ListPreviewOptions{
			ShowNumbers: opts.showNumbers,
			ShowEmoji:   opts.showEmoji,
			UseRelative: opts.useRelative,
		}
		PrintListPreviewWithOptions(records, previewOpts)
		return
	}
	if opts.tree {
		PrintListTree(records)
		return
	}
	printListOutput(records, opts.verbose)
	printHints(listHints())
}

// parseListFlags parses flags and positional modes for the list command.
func parseListFlags(args []string) listOptions {
	opts := listOptions{
		showNumbers: true,
		showEmoji:   true,
	}
	for _, a := range args {
		switch strings.ToLower(a) {
		case "preview", "-p", "--preview", "--gap", "gap":
			opts.preview = true
		case "tree", "-t", "--tree":
			opts.tree = true
		case "-v", "--verbose", "verbose":
			opts.verbose = true
		case "--no-numbers":
			opts.showNumbers = false
		case "--no-emoji":
			opts.showEmoji = false
		case "--relative", "-r", "--rel":
			opts.useRelative = true
		}
	}
	opts.group = extractGroupFlag(args)
	return opts
}

func extractGroupFlag(args []string) string {
	g := extractFlagVal(args, "-group")
	if g == "" {
		return extractFlagVal(args, "--group")
	}
	return g
}

func fetchListRecords(db *store.DB, groupFilter string) []model.ScanRecord {
	records, err := loadListRecords(db, groupFilter)
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		return nil
	}
	if err != nil {
		return nil
	}
	return records
}

// isListTypeOrGroups checks if the arg is a type keyword or "groups".
func isListTypeOrGroups(arg string) bool {
	lower := strings.ToLower(arg)
	if lower == constants.SubCmdGroups {
		return true
	}
	_, ok := typeKeywords[lower]
	return ok
}

// handleListSpecial handles gitmap ls <type> or gitmap ls groups.
func handleListSpecial(keyword string, args []string) {
	lower := strings.ToLower(keyword)
	if lower == constants.SubCmdGroups {
		runGroupList()
		printHints(listGroupsHints())
		return
	}
	typeKey := typeKeywords[lower]
	runProjectRepos(typeKey, args)
}

// loadListRecords loads repos, optionally filtered by group.
func loadListRecords(db *store.DB, group string) ([]model.ScanRecord, error) {
	if len(group) > 0 {
		return db.ShowGroup(group)
	}

	return db.ListRepos()
}

// printListOutput renders the list table to stdout.
func printListOutput(records []model.ScanRecord, verbose bool) {
	if verbose || len(records) == 0 {
		fmt.Printf(constants.MsgListDBPath, store.DefaultDBPath())
	}

	if len(records) == 0 {
		fmt.Println(constants.MsgListEmpty)

		return
	}

	fmt.Println(constants.MsgListHeader)
	fmt.Println(constants.MsgListSeparator)
	for _, r := range records {
		printListRow(r, verbose)
	}
}

// printListRow prints a single row in list output.
func printListRow(r model.ScanRecord, verbose bool) {
	if verbose {
		fmt.Printf(constants.MsgListVerboseFmt, r.Slug, r.RepoName, r.AbsolutePath)

		return
	}

	fmt.Printf(constants.MsgListRowFmt, r.Slug, r.RepoName)
}

// openDb opens the gitmap database from the binary's data directory.
func openDb() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

// openDB is a backwards-compatible alias for openDb.
var openDB = openDb
