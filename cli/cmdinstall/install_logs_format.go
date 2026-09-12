package cmdinstall

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func printInstallLogsTable(out io.Writer, logs []store.InstallationLogRecord) error {
	if len(logs) == 0 {
		fmt.Fprintln(out, "No installation logs found.")

		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOOL\tACTION\tVERSION\tMANAGER\tDURATION\tSTATUS\tCREATED_AT")
	for _, r := range logs {
		printSingleLogRow(w, r)
	}

	if err := w.Flush(); err != nil {
		return apperror.WrapSimple(err, "install_logs.flush")
	}

	return nil
}

func printSingleLogRow(w io.Writer, r store.InstallationLogRecord) {
	durStr := formatLogDuration(r.DurationMs)
	statusStr := formatLogStatus(r.IsSuccess)
	verStr := formatLogDisplayVal(r.Version)
	mgrStr := formatLogDisplayVal(r.PackageManager)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		r.Tool, r.Action, verStr, mgrStr, durStr, statusStr, r.CreatedAt)
}

func formatLogDisplayVal(v string) string {
	if v == "" {
		return "-"
	}

	return v
}

func formatLogStatus(isSuccess bool) string {
	if isSuccess {
		return "success"
	}

	return "failed"
}

func formatLogDuration(ms int64) string {
	if ms <= 0 {
		return "0ms"
	}

	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}

	return formatSecondsDuration(ms)
}

func formatSecondsDuration(ms int64) string {
	sec := float64(ms) / 1000.0
	if sec < 60 {
		return fmt.Sprintf("%.1fs", sec)
	}

	return fmt.Sprintf("%dm%ds", int(sec)/60, int(sec)%60)
}

// isInstallLogsCommand reports whether the args invoke the installation logs sub-command.
func isInstallLogsCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if args[0] == "logs" || args[0] == "log" {
		return true
	}

	return hasArgFlag(args, "--logs") || hasArgFlag(args, "-logs")
}

// extractInstallLogsArgs removes the 'logs' or '--logs' trigger from args.
func extractInstallLogsArgs(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}

	if args[0] == "logs" || args[0] == "log" {
		return args[1:]
	}

	return filterOutLogsFlag(args)
}

func filterOutLogsFlag(args []string) []string {
	var filtered []string
	for _, a := range args {
		if a != "--logs" && a != "-logs" {
			filtered = append(filtered, a)
		}
	}

	return filtered
}
