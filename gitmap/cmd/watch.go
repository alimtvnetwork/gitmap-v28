package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runWatch handles the "watch" subcommand.
func runWatch(args []string) {
	checkHelp("watch", args)
	interval, groupName, noFetch, jsonMode := parseWatchFlags(args)
	records := loadWatchRecords(groupName)

	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrWatchNoRepos)
		os.Exit(1)
	}

	if jsonMode {
		printWatchJSON(records, noFetch)

		return
	}

	runWatchLoop(records, interval, noFetch)
}

// parseWatchFlags parses flags for the watch command.
func parseWatchFlags(args []string) (int, string, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdWatch, flag.ExitOnError)
	interval := fs.Int("interval", constants.WatchDefaultInterval, constants.FlagDescWatchInterval)
	group := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(group, "g", "", constants.FlagDescGroup)
	noFetch := fs.Bool("no-fetch", false, constants.FlagDescWatchNoFetch)
	jsonFlag := fs.Bool("json", false, constants.FlagDescWatchJSON)
	fs.Parse(args)

	if *interval < constants.WatchMinInterval {
		*interval = constants.WatchMinInterval
	}

	return *interval, *group, *noFetch, *jsonFlag
}

// loadWatchRecords loads repos for watching.
func loadWatchRecords(groupName string) []model.ScanRecord {
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}

	db, err := openDB()
	if err != nil {
		return loadRecordsJSONFallback()
	}
	defer db.Close()

	repos, err := db.ListRepos()
	if err != nil {
		return loadRecordsJSONFallback()
	}

	return repos
}

// runWatchLoop runs the refresh loop until interrupted.
func runWatchLoop(records []model.ScanRecord, interval int, noFetch bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	for {
		printWatchDashboard(records, interval, noFetch)

		select {
		case <-stop:
			fmt.Println(constants.WatchStoppedMsg)

			return
		case <-time.After(time.Duration(interval) * time.Second):
		}
	}
}

// printWatchJSON outputs a single snapshot as stable JSON and exits.
func printWatchJSON(records []model.ScanRecord, noFetch bool) {
	snapshots := collectAllStatuses(records, noFetch)
	summary := buildWatchSummary(snapshots)

	if err := encodeWatchJSON(os.Stdout, snapshots, summary, time.Now().UTC().Format(time.RFC3339)); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode watch result to JSON: %v\n", err)
	}
}
