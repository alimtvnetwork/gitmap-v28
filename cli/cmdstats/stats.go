package cmdstats

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

// isSSHForwardCandidate reports whether args should be forwarded to SSH.
func isSSHForwardCandidate(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return strutil.EqualFoldAny(
		args[0],
		"fix-auth", "login", "host", "hosts", "nodes",
		"enable", "port", "troubleshoot", "key", "copy-id", "auth-key",
	)
}

// runStats handles the "stats" subcommand.
func RunStats(args []string) error {
	if isSSHForwardCandidate(args) {
		return cmdssh.RunSSH(args)
	}

	checkHelp("stats", args)
	cmdFilter, isJSON := parseStatsFlags(args)
	overall, commands := loadStats(cmdFilter)
	printStatsOutput(overall, commands, isJSON)

	return nil
}

// printStatsOutput routes stats to JSON or terminal format.
func printStatsOutput(overall model.OverallStats, commands []model.CommandStats, isJSON bool) {
	if isJSON {
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

func handleOpenDBError(err error) {
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
}

// loadStats fetches aggregated stats from the database.
func loadStats(cmdFilter string) (model.OverallStats, []model.CommandStats) {
	db, err := openDB()

	if err != nil {
		handleOpenDBError(err)

		return model.OverallStats{}, nil
	}

	defer db.Close()

	overall, err := db.QueryOverallStats()
	handleStatsError(err)

	commands := loadStatsCommands(db, cmdFilter)

	return overall, commands
}

type statsQuerier interface {
	QueryCommandStats() ([]model.CommandStats, error)
	QueryCommandStatsFor(string) ([]model.CommandStats, error)
}

// loadStatsCommands loads per-command stats with optional filter.
func loadStatsCommands(db statsQuerier, cmdFilter string) []model.CommandStats {
	if cmdFilter != "" {
		records, err := db.QueryCommandStatsFor(cmdFilter)
		handleStatsError(err)

		return records
	}

	records, err := db.QueryCommandStats()
	handleStatsError(err)

	return records
}

// printStatsRows prints all per-command stats rows.
func printStatsRows(commands []model.CommandStats) {
	for _, s := range commands {
		fmt.Printf(constants.MsgStatsRowFmt, s.Command, s.TotalRuns, s.SuccessCount,
			s.FailCount, s.FailRate, s.AvgDuration, s.MinDuration, s.MaxDuration, s.LastUsed)
	}
}

// printStatsHeader prints the overall summary and table header.
func printStatsHeader(overall model.OverallStats) {
	fmt.Println(constants.MsgStatsHeader)
	fmt.Println(constants.MsgStatsSeparator)
	fmt.Printf(constants.MsgStatsOverallFmt, overall.TotalCommands, overall.UniqueCommands,
		overall.TotalSuccess, overall.TotalFail, overall.OverallFailRate, overall.AvgDuration)
	fmt.Println(constants.MsgStatsSeparator)
	fmt.Println(constants.MsgStatsColumns)
}

// printStatsTerminal prints stats in table format.
func printStatsTerminal(overall model.OverallStats, commands []model.CommandStats) {
	if overall.TotalCommands == 0 {
		fmt.Print(constants.MsgStatsEmpty)

		return
	}

	printStatsHeader(overall)
	printStatsRows(commands)
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

func handleStatsLegacyError(err error) {
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
}

func handleStatsQueryError(err error) {
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

func handleStatsError(err error) {
	if err == nil {
		return
	}

	if isLegacyDataError(err) {
		handleStatsLegacyError(err)

		return
	}

	handleStatsQueryError(err)
}
