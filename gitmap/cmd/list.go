package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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

// runList handles the "list" subcommand.
func runList(args []string) {
	checkHelp("list", args)

	if len(args) > 0 && isListTypeOrGroups(args[0]) {
		handleListSpecial(args[0], args[1:])

		return
	}

	groupFilter, verboseMode := parseListFlags(args)
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	records, err := loadListRecords(db, groupFilter)
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	printListOutput(records, verboseMode)
	printHints(listHints())
}

// isListTypeOrGroups checks if the arg is a type keyword or "groups".
func isListTypeOrGroups(arg string) bool {
	lower := strings.ToLower(arg)
	if lower == "groups" {
		return true
	}
	_, ok := typeKeywords[lower]

	return ok
}

// handleListSpecial handles gitmap ls <type> or gitmap ls groups.
func handleListSpecial(keyword string, args []string) {
	lower := strings.ToLower(keyword)
	if lower == "groups" {
		runGroupList()
		printHints(listGroupsHints())

		return
	}
	typeKey := typeKeywords[lower]
	runProjectRepos(typeKey, args)
}

// parseListFlags parses flags for the list command.
func parseListFlags(args []string) (group string, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdList, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	vFlag := fs.Bool("verbose", false, constants.FlagDescListVerbose)
	fs.Parse(args)

	return *gFlag, *vFlag
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
	if verbose {
		fmt.Printf(constants.MsgListDBPath, store.DefaultDBPath())
	}
	if len(records) == 0 {
		if !verbose {
			fmt.Printf(constants.MsgListDBPath, store.DefaultDBPath())
		}
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

// openDB opens the gitmap database from the binary's data directory.
func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		return nil, err
	}

	return db, nil
}
