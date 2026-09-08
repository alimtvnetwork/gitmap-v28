package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type installLogsOptions struct {
	Failed bool
	Tool   string
	Limit  int
}

func parseInstallLogsFlags(args []string) (installLogsOptions, error) {
	fs := flag.NewFlagSet("install logs", flag.ContinueOnError)
	var opts installLogsOptions
	fs.BoolVar(&opts.Failed, "failed", false, "Filter failed installations only")
	fs.StringVar(&opts.Tool, "tool", "", "Filter logs by tool name")
	fs.IntVar(&opts.Limit, "limit", 50, "Limit number of logs to display")
	if err := fs.Parse(args); err != nil {

		return opts, err
	}
	if opts.Tool == "" && len(fs.Args()) > 0 {
		opts.Tool = fs.Arg(0)
	}

	return opts, nil
}

func resolveLogsLimit(limit int) int {
	if limit > 0 {

		return limit
	}

	return 50
}

func queryInstallationLogs(db *store.InstallationSplitDB, opts installLogsOptions) ([]store.InstallationLogRecord, error) {
	limit := resolveLogsLimit(opts.Limit)
	if opts.Tool != "" && opts.Failed {

		return queryFailedLogsForTool(db, opts.Tool, limit)
	}
	if opts.Tool != "" {

		return db.GetLogsByTool(opts.Tool, limit)
	}
	if opts.Failed {

		return db.GetFailedLogs(limit)
	}

	return db.GetLogs(limit)
}

func queryFailedLogsForTool(db *store.InstallationSplitDB, tool string, limit int) ([]store.InstallationLogRecord, error) {
	logs, err := db.GetLogsByTool(tool, limit*2)
	if err != nil {

		return nil, err
	}

	return filterFailedLogs(logs, limit), nil
}

func filterFailedLogs(logs []store.InstallationLogRecord, limit int) []store.InstallationLogRecord {
	var filtered []store.InstallationLogRecord
	for _, r := range logs {
		if !r.IsSuccess && len(filtered) < limit {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func executeInstallLogs(out io.Writer, db *store.InstallationSplitDB, args []string) error {
	opts, err := parseInstallLogsFlags(args)
	if err != nil {

		return err
	}
	logs, err := queryInstallationLogs(db, opts)
	if err != nil {

		return err
	}

	return printInstallLogsTable(out, logs)
}

func runInstallLogsWithOutput(out io.Writer, args []string) error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {

		return apperror.WrapSimple(err, "install_logs.openDB")
	}
	defer splitDB.Close()

	return executeInstallLogs(out, splitDB, args)
}

// runInstallLogs prints a table of recent installation execution logs.
func runInstallLogs(args []string) error {

	return runInstallLogsWithOutput(os.Stdout, args)
}
