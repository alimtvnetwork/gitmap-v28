package cmdstats

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats/enums"
	rootconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

// isSSHForwardCandidate reports whether args should be forwarded to SSH.
func isSSHForwardCandidate(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return strutil.EqualFoldAny(args[0], constants.AllSSHForwardCandidates()...)
}

// RunStats handles the "stats" subcommand.
func RunStats(args []string) error {
	if isSSHForwardCandidate(args) {
		return cmdssh.RunSSH(args)
	}

	checkHelp(constants.CmdStats, args)
	cmdFilter, isJSON := parseStatsFlags(args)
	overall, commands := loadStats(cmdFilter)
	printStatsOutput(overall, commands, isJSON)

	return nil
}

// resolveStatsFormat maps flag state to a typed output format enum.
func resolveStatsFormat(isJSON bool) enums.StatsOutputFormatType {
	if isJSON {
		return enums.StatsOutputFormatJSON
	}

	return enums.StatsOutputFormatTable
}

// printStatsOutput routes stats to JSON or terminal format.
func printStatsOutput(overall model.OverallStats, commands []model.CommandStats, isJSON bool) {
	outputFormat := resolveStatsFormat(isJSON)
	if outputFormat.IsJSON() {
		printStatsJSON(overall, commands)

		return
	}

	printStatsTerminal(overall, commands)
}

// parseStatsFlags parses --command and --json flags.
func parseStatsFlags(args []string) (string, bool) {
	fs := flag.NewFlagSet(constants.CmdStats, flag.ExitOnError)
	command := fs.String(constants.FlagCommand, "", constants.FlagDescStatsCommand)
	jsonFlag := fs.Bool(constants.FlagJSON, false, rootconstants.FlagDescLBJSON)
	fs.Parse(args)

	return *command, *jsonFlag
}

func handleOpenDBError(err error) {
	appErr := apperror.WrapWithDetails(
		err,
		constants.OpStatsOpenDB,
		constants.ErrCodeE1154,
		fmt.Sprintf(rootconstants.ErrStatsQuery, err),
		constants.DomainStats,
		apperror.ErrorTypeExecution,
		apperror.SeverityFatal,
		nil,
	)
	cliexit.HandleError(appErr, constants.ExitGeneralError)
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
		fmt.Printf(rootconstants.MsgStatsRowFmt, s.Command, s.TotalRuns, s.SuccessCount,
			s.FailCount, s.FailRate, s.AvgDuration, s.MinDuration, s.MaxDuration, s.LastUsed)
	}
}

// printStatsHeader prints the overall summary and table header.
func printStatsHeader(overall model.OverallStats) {
	fmt.Println(rootconstants.MsgStatsHeader)
	fmt.Println(rootconstants.MsgStatsSeparator)
	fmt.Printf(rootconstants.MsgStatsOverallFmt, overall.TotalCommands, overall.UniqueCommands,
		overall.TotalSuccess, overall.TotalFail, overall.OverallFailRate, overall.AvgDuration)
	fmt.Println(rootconstants.MsgStatsSeparator)
	fmt.Println(rootconstants.MsgStatsColumns)
}

// printStatsTerminal prints stats in table format.
func printStatsTerminal(overall model.OverallStats, commands []model.CommandStats) {
	if overall.TotalCommands == 0 {
		fmt.Print(rootconstants.MsgStatsEmpty)

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
		fmt.Fprintf(os.Stderr, constants.ErrStatsJSONMarshalFmt, err)

		return
	}

	fmt.Println(string(data))
}

func handleStatsLegacyError(err error) {
	appErr := apperror.WrapWithDetails(
		err,
		constants.OpStatsLegacyData,
		constants.ErrCodeE1155,
		rootconstants.MsgLegacyProjectData,
		constants.DomainStats,
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	cliexit.HandleError(appErr, constants.ExitGeneralError)
}

func handleStatsQueryError(err error) {
	appErr := apperror.WrapWithDetails(
		err,
		constants.OpStatsQuery,
		constants.ErrCodeE1156,
		fmt.Sprintf(rootconstants.ErrStatsQuery, err),
		constants.DomainStats,
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	cliexit.HandleError(appErr, constants.ExitGeneralError)
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
