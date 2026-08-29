package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runStats handles the "stats" subcommand.
func runStats(args []string) error {
	checkHelp("stats", args)
	cmdFilter, jsonOut := parseStatsFlags(args)
	overall, commands := loadStats(cmdFilter)

	if jsonOut {
		printStatsJSON(overall, commands)

		return nil
	}

	printStatsTerminal(overall, commands)
	return nil
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
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.stats.openDB",
			"E1154",
			fmt.Sprintf(constants.ErrStatsQuery, err),
			"cmd.stats",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return model.OverallStats{}, nil
	}
	defer db.Close()

	overall, err := db.QueryOverallStats()
	handleStatsError(err)

	commands := loadStatsCommands(db, cmdFilter)

	return overall, commands
}

// loadStatsCommands loads per-command stats with optional filter.
func loadStatsCommands(db interface {
	QueryCommandStats() ([]model.CommandStats, error)
	QueryCommandStatsFor(string) ([]model.CommandStats, error)
}, cmdFilter string) []model.CommandStats {
	if cmdFilter != "" {
		records, err := db.QueryCommandStatsFor(cmdFilter)
		handleStatsError(err)

		return records
	}

	records, err := db.QueryCommandStats()
	handleStatsError(err)

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

// printStatsJSON outputs stats as JSON.
func printStatsJSON(overall model.OverallStats, commands []model.CommandStats) {
	overall.Commands = commands
	data, err := json.MarshalIndent(overall, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal stats to JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func handleStatsError(err error) {
	if err == nil {
		return
	}
	if isLegacyDataError(err) {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.stats.legacyData",
			"E1155",
			constants.MsgLegacyProjectData,
			"cmd.stats",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return
	}
	appErr := apperror.WrapWithDetails(
		err,
		"cmd.stats.query",
		"E1156",
		fmt.Sprintf(constants.ErrStatsQuery, err),
		"cmd.stats",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	cliexit.HandleError(appErr, 1)
}
