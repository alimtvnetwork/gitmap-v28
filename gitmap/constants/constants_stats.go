package constants

// gitmap:cmd top-level
// Stats CLI commands.
const (
	CmdStats      = "stats"
	CmdStatsAlias = "ss"
)

// Stats help text.
const (
	HelpStats = "  stats (ss)          Show aggregated command usage statistics (--json, --command)"
)

// Stats flag descriptions.
const (
	FlagDescStatsCommand = "Show stats for a specific command only"
)

// Stats SQL queries.
const (
	SQLStatsPerCommand = `SELECT Command,
		COUNT(*) AS TotalRuns,
		SUM(CASE WHEN ExitCode = 0 THEN 1 ELSE 0 END) AS SuccessCount,
		SUM(CASE WHEN ExitCode != 0 THEN 1 ELSE 0 END) AS FailCount,
		ROUND(SUM(CASE WHEN ExitCode != 0 THEN 1.0 ELSE 0.0 END) / COUNT(*) * 100, 1) AS FailRate,
		COALESCE(AVG(DurationMs), 0) AS AvgDuration,
		COALESCE(MIN(DurationMs), 0) AS MinDuration,
		COALESCE(MAX(DurationMs), 0) AS MaxDuration,
		MAX(StartedAt) AS LastUsed
		FROM CommandHistory
		GROUP BY Command
		ORDER BY TotalRuns DESC`

	SQLStatsForCommand = `SELECT Command,
		COUNT(*) AS TotalRuns,
		SUM(CASE WHEN ExitCode = 0 THEN 1 ELSE 0 END) AS SuccessCount,
		SUM(CASE WHEN ExitCode != 0 THEN 1 ELSE 0 END) AS FailCount,
		ROUND(SUM(CASE WHEN ExitCode != 0 THEN 1.0 ELSE 0.0 END) / COUNT(*) * 100, 1) AS FailRate,
		COALESCE(AVG(DurationMs), 0) AS AvgDuration,
		COALESCE(MIN(DurationMs), 0) AS MinDuration,
		COALESCE(MAX(DurationMs), 0) AS MaxDuration,
		MAX(StartedAt) AS LastUsed
		FROM CommandHistory
		WHERE Command = ?
		GROUP BY Command`

	SQLStatsOverall = `SELECT
		COUNT(*) AS TotalCommands,
		COUNT(DISTINCT Command) AS UniqueCommands,
		SUM(CASE WHEN ExitCode = 0 THEN 1 ELSE 0 END) AS TotalSuccess,
		SUM(CASE WHEN ExitCode != 0 THEN 1 ELSE 0 END) AS TotalFail,
		ROUND(SUM(CASE WHEN ExitCode != 0 THEN 1.0 ELSE 0.0 END) / MAX(COUNT(*), 1) * 100, 1) AS OverallFailRate,
		COALESCE(AVG(DurationMs), 0) AS AvgDuration
		FROM CommandHistory`
)

// Stats terminal formatting.
const (
	MsgStatsHeader     = "Command Usage Statistics"
	MsgStatsSeparator  = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	MsgStatsOverallFmt = "Total: %d executions (%d unique commands) | Success: %d | Fail: %d (%.1f%%) | Avg: %dms\n"
	MsgStatsColumns    = "COMMAND         RUNS   SUCCESS  FAIL  FAIL%%   AVG(ms)  MIN(ms)  MAX(ms)  LAST USED"
	MsgStatsRowFmt     = "%-15s %-6d %-8d %-5d %-7.1f %-8d %-8d %-8d %s\n"
	MsgStatsEmpty      = "No command history found. Run some commands first.\n"
	ErrStatsQuery      = "failed to query stats: %v"
)
