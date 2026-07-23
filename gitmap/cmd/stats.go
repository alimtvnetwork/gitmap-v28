package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runStats handles the "stats" subcommand.
func runStats(args []string) {
	checkHelp("stats", args)
	cmdFilter, jsonOut := parseStatsFlags(args)
	overall, commands := loadStats(cmdFilter)

	if jsonOut {
		printStatsJSON(overall, commands)

		return
	}

	printStatsTerminal(overall, commands)
}

// parseStatsFlags parses --command and --json flags.
func parseStatsFlags(args []string) (string, bool) {
	fs := flag.NewFlagSet(constants.CmdStats, flag.ExitOnError)
	command := fs.String("command", "", constants.FlagDescStatsCommand)
	jsonFlag := fs.Bool("json", false, constants.FlagDescLBJSON)
	fs.Parse(args)

	return *command, *jsonFlag
}

// loadStats fetches aggregated stats from the database.
func loadStats(cmdFilter string) (model.OverallStats, []model.CommandStats) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrStatsQuery+"\n", err)
		os.Exit(1)
	}
	defer db.Close()

	overall, err := db.QueryOverallStats()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrStatsQuery+"\n", err)
		os.Exit(1)
	}

	commands := loadStatsCommands(db, cmdFilter)
	overall.Commands = commands

	return overall, commands
}

// loadStatsCommands loads per-command stats with optional filter.
func loadStatsCommands(db interface {
	QueryCommandStats() ([]model.CommandStats, error)
	QueryCommandStatsFor(string) ([]model.CommandStats, error)
}, cmdFilter string) []model.CommandStats {
	if cmdFilter != "" {
		records, err := db.QueryCommandStatsFor(cmdFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrStatsQuery+"\n", err)
			os.Exit(1)
		}

		return records
	}

	records, err := db.QueryCommandStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrStatsQuery+"\n", err)
		os.Exit(1)
	}

	return records
}

// printStatsTerminal prints stats in table format.
func printStatsTerminal(overall model.OverallStats, commands []model.CommandStats) {
	if overall.TotalCommands == 0 {
		fmt.Print(constants.MsgStatsEmpty)

		return
	}

	fmt.Println(constants.MsgStatsHeader)
	fmt.Println(constants.MsgStatsSeparator)
	fmt.Printf(constants.MsgStatsOverallFmt, overall.TotalCommands, overall.UniqueCommands,
		overall.TotalSuccess, overall.TotalFail, overall.OverallFailRate, overall.AvgDuration)
	fmt.Println(constants.MsgStatsSeparator)
	fmt.Println(constants.MsgStatsColumns)

	for _, s := range commands {
		fmt.Printf(constants.MsgStatsRowFmt, s.Command, s.TotalRuns, s.SuccessCount,
			s.FailCount, s.FailRate, s.AvgDuration, s.MinDuration, s.MaxDuration, s.LastUsed)
	}
}

// printStatsJSON outputs stats as stable JSON (see statsrender.go).
func printStatsJSON(overall model.OverallStats, commands []model.CommandStats) {
	overall.Commands = commands
	if err := encodeStatsJSON(os.Stdout, overall, commands); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal stats to JSON: %v\n", err)
	}
}
