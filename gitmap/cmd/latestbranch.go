// Package cmd — latest-branch command handler.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
)

// latestBranchConfig holds parsed flags for the latest-branch command.
type latestBranchConfig struct {
	remote           string
	filterByRemote   bool
	containsFallback bool
	top              int
	format           string
	shouldFetch      bool
	sortBy           string
	filter           string
	// shouldSwitch toggles the post-report `git checkout` performed by
	// maybeSwitchToLatest. Wired to `--switch` and its short form `-s`.
	shouldSwitch bool
}

// runLatestBranch handles the 'latest-branch' / 'lb' command.
func runLatestBranch(args []string) {
	checkHelp("latest-branch", args)
	cfg := parseLatestBranchFlags(args)
	validateLatestBranchRepo()
	fetchLatestBranchRefs(cfg)
	refs := loadFilteredRefs(cfg)
	items := readAndSortBranches(refs, cfg.sortBy)
	result := resolveLatestResult(items, cfg)
	dispatchLatestOutput(result, items, cfg)
	maybeSwitchToLatest(result, cfg)
}

// validateLatestBranchRepo exits if the current directory is outside a git repo.
func validateLatestBranchRepo() {
	if gitutil.IsInsideWorkTree() {

		return
	}
	fmt.Fprintln(os.Stderr, constants.ErrLatestBranchNotRepo)
	os.Exit(1)
}

// fetchLatestBranchRefs fetches remotes when shouldFetch is enabled.
func fetchLatestBranchRefs(cfg latestBranchConfig) {
	if cfg.shouldFetch {
		isTerminal := cfg.format == constants.OutputTerminal
		if isTerminal {
			fmt.Println(constants.MsgLatestBranchFetching)
		}
		err := gitutil.FetchAllPrune()
		if err != nil && isTerminal {
			fmt.Fprintf(os.Stderr, constants.MsgLatestBranchFetchWarning, err)
		}
	}
}

// loadFilteredRefs lists remote branches and applies remote + pattern filters.
func loadFilteredRefs(cfg latestBranchConfig) []string {
	refs, err := gitutil.ListRemoteBranches()
	if err != nil || len(refs) == 0 {
		printNoRefsError(cfg)
		os.Exit(1)
	}
	refs = applyRemoteFilter(refs, cfg)
	refs = applyPatternFilter(refs, cfg)

	return refs
}

// printNoRefsError prints the appropriate "no refs" error message.
func printNoRefsError(cfg latestBranchConfig) {
	if cfg.filterByRemote {
		fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoRefs, cfg.remote)

		return
	}
	fmt.Fprintln(os.Stderr, constants.ErrLatestBranchNoRefsAll)
}

// applyRemoteFilter filters refs by remote when filterByRemote is set.
func applyRemoteFilter(refs []string, cfg latestBranchConfig) []string {
	if cfg.filterByRemote {
		filtered := gitutil.FilterByRemote(refs, cfg.remote)
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoRefs, cfg.remote)
			os.Exit(1)
		}

		return filtered
	}

	return refs
}

// applyPatternFilter filters refs by glob/substring when filter is set.
func applyPatternFilter(refs []string, cfg latestBranchConfig) []string {
	if len(cfg.filter) > 0 {
		filtered := gitutil.FilterByPattern(refs, cfg.filter)
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoMatch, cfg.filter)
			os.Exit(1)
		}

		return filtered
	}

	return refs
}

// readAndSortBranches reads tip commits and sorts by the given order.
func readAndSortBranches(refs []string, sortBy string) []gitutil.RemoteBranchInfo {
	items, err := gitutil.ReadBranchTips(refs)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLatestBranchNoCommits+"\n")
		os.Exit(1)
	}
	if sortBy == constants.SortByName {
		gitutil.SortByNameAsc(items)
	} else {
		gitutil.SortByDateDesc(items)
	}

	return items
}

// parseLatestBranchFlags parses flags and returns a config struct.
func parseLatestBranchFlags(args []string) latestBranchConfig {
	fs := flag.NewFlagSet(constants.CmdLatestBranch, flag.ExitOnError)
	var cfg latestBranchConfig
	var allRemotes, noFetch, jsonOut, switchLong, switchShort bool
	fs.StringVar(&cfg.remote, "remote", "origin", constants.FlagDescLBRemote)
	fs.BoolVar(&allRemotes, "all-remotes", false, constants.FlagDescLBAllRemotes)
	fs.BoolVar(&cfg.containsFallback, "contains-fallback", false, constants.FlagDescLBContains)
	fs.IntVar(&cfg.top, "top", 0, constants.FlagDescLBTop)
	fs.StringVar(&cfg.format, "format", constants.OutputTerminal, constants.FlagDescLBFormat)
	fs.BoolVar(&jsonOut, "json", false, constants.FlagDescLBJSON)
	fs.BoolVar(&noFetch, "no-fetch", false, constants.FlagDescLBNoFetch)
	fs.StringVar(&cfg.sortBy, "sort", constants.SortByDate, constants.FlagDescLBSort)
	fs.StringVar(&cfg.filter, "filter", "", constants.FlagDescLBFilter)
	// --switch / -s. Both registered against the same effect; either
	// being true flips cfg.shouldSwitch on. Go's flag package doesn't
	// natively support aliases so we OR them in resolveLatestBranchConfig.
	fs.BoolVar(&switchLong, "switch", false, constants.FlagDescLBSwitch)
	fs.BoolVar(&switchShort, "s", false, constants.FlagDescLBSwitchShort)
	fs.Parse(args)
	cfg.shouldSwitch = switchLong || switchShort

	return resolveLatestBranchConfig(fs, cfg, allRemotes, noFetch, jsonOut)
}

// resolveLatestBranchConfig converts parsed flags into positive-logic config.
func resolveLatestBranchConfig(fs *flag.FlagSet, cfg latestBranchConfig, allRemotes, noFetch, jsonOut bool) latestBranchConfig {
	cfg.filterByRemote = !allRemotes
	cfg.shouldFetch = !noFetch
	if jsonOut {
		cfg.format = constants.OutputJSON
	}
	cfg.top = resolvePositionalTop(fs, cfg.top)

	return cfg
}

// resolvePositionalTop checks for a bare integer positional argument.
func resolvePositionalTop(fs *flag.FlagSet, current int) int {
	if current > 0 || fs.NArg() == 0 {

		return current
	}
	n, err := strconv.Atoi(fs.Arg(0))
	if err == nil && n > 0 {

		return n
	}

	return current
}
