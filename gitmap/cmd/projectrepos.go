// Package cmd — projectrepos.go handles project type query commands.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runProjectRepos handles go-repos, node-repos, react-repos, cpp-repos, csharp-repos.
func runProjectRepos(typeKey string, args []string) {
	checkHelp(typeKey+"-repos", args)
	jsonOut, countOnly := parseProjectReposFlags(args)
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.MsgProjectNoDB)
		os.Exit(1)
	}
	defer db.Close()

	if countOnly {
		printProjectCount(db, typeKey)

		return
	}
	printProjectList(db, typeKey, jsonOut)
	if !jsonOut {
		printHints(projectReposHints())
	}
}

// parseProjectReposFlags parses --json and --count flags.
func parseProjectReposFlags(args []string) (bool, bool) {
	fs := flag.NewFlagSet("project-repos", flag.ExitOnError)
	jsonOut := fs.Bool(constants.FlagProjectJSON, false, "Output as JSON")
	countOnly := fs.Bool(constants.FlagProjectCount, false, "Print count only")
	_ = fs.Parse(args)

	return *jsonOut, *countOnly
}

// printProjectCount prints the count of projects for a type.
func printProjectCount(db *store.DB, typeKey string) {
	count, err := db.CountProjectsByTypeKey(typeKey)
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrProjectQuery, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgProjectCount, count)
}

// printProjectList queries and displays projects for a type.
func printProjectList(db *store.DB, typeKey string, jsonOut bool) {
	projects, err := db.SelectProjectsByTypeKey(typeKey)
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrProjectQuery, err)
		os.Exit(1)
	}
	if len(projects) == 0 {
		fmt.Printf(constants.MsgProjectNoneFound, typeKey)

		return
	}
	if jsonOut {
		printProjectsJSON(projects)

		return
	}
	printProjectsTerminal(projects)
	printProjectsSummary(projects)
}
