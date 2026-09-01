// Package cmd — schedule_export.go: exports scheduled tasks and their split DB logs
// to JSON, YAML, SQLite database, or ZIP archives with format auto-inference and except filters.
package cmd

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

type scheduleExportBundle struct {
	Task store.SchedulerTask       `json:"task" yaml:"task"`
	Runs []store.ScheduleRunRecord `json:"runs" yaml:"runs"`
}

type scheduleExportOpts struct {
	TargetName string
	FilePath   string
	Format     string
	ExceptList []string
	IsAll      bool
}

func runScheduleExport(args []string) error {
	opts := parseScheduleExportOpts(args)
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	bundles, err := collectExportBundles(db, opts)
	if err != nil {
		return err
	}
	return writeScheduleExportOutput(bundles, opts)
}

func parseScheduleExportOpts(args []string) scheduleExportOpts {
	opts := scheduleExportOpts{Format: constants.OutputJSON}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if matchFlagWithVal(a, "-f", "--file", "--filepath", "-o", "--output") {
			opts.FilePath = extractFlagValue(&i, args)
		} else if matchFlagWithVal(a, "-except", "--except", "--exclude") {
			opts.ExceptList = parseExceptTokens(extractFlagValue(&i, args))
		} else if matchFlagWithVal(a, "--format") {
			opts.Format = strings.ToLower(extractFlagValue(&i, args))
		} else if a == "--json" {
			opts.Format = constants.OutputJSON
		} else if a == "--yaml" || a == "--yml" || a == "-y" {
			opts.Format = constants.OutputYAML
		} else if a == "--sqlite" || a == "--db" {
			opts.Format = "sqlite"
		} else if a == "--zip" {
			opts.Format = "zip"
		} else if !strings.HasPrefix(a, "-") && opts.TargetName == "" {
			opts.TargetName = a
		}
	}
	inferScheduleExportFormatFromPath(&opts)
	opts.IsAll = opts.TargetName == "" || opts.TargetName == "*" || opts.TargetName == "all" || opts.TargetName == "export-all" || opts.TargetName == "import-all"
	return opts
}

func inferScheduleExportFormatFromPath(opts *scheduleExportOpts) {
	if opts.FilePath == "" {
		return
	}
	ext := strings.ToLower(filepath.Ext(opts.FilePath))
	if ext == ".yaml" || ext == ".yml" {
		opts.Format = constants.OutputYAML
	}
	if ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" {
		opts.Format = "sqlite"
	}
	if ext == ".zip" {
		opts.Format = "zip"
	}
}

func parseExceptTokens(raw string) []string {
	var list []string
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			list = append(list, trimmed)
		}
	}
	return list
}

func collectExportBundles(db *store.DB, opts scheduleExportOpts) ([]scheduleExportBundle, error) {
	tasks, err := db.ListSchedules()
	if err != nil {
		return nil, apperror.WrapSimple(err, "list schedules for export")
	}
	var bundles []scheduleExportBundle
	for _, t := range tasks {
		if shouldSkipSchedule(t.Name, opts) {
			continue
		}
		runs := fetchScheduleRunsSafe(t.Slug)
		bundles = append(bundles, scheduleExportBundle{Task: t, Runs: runs})
	}
	if len(bundles) == 0 && !opts.IsAll {
		return nil, apperror.NewSimple("schedule "+opts.TargetName+" not found", "E6010")
	}
	return bundles, nil
}

func shouldSkipSchedule(name string, opts scheduleExportOpts) bool {
	if !opts.IsAll && !strings.EqualFold(name, opts.TargetName) {
		return true
	}
	for _, ex := range opts.ExceptList {
		if strings.EqualFold(name, ex) {
			return true
		}
	}
	return false
}

func fetchScheduleRunsSafe(slug string) []store.ScheduleRunRecord {
	splitDB, err := store.OpenScheduleSplitDB(slug)
	if err != nil {
		return nil
	}
	defer splitDB.Close()
	runs, _ := splitDB.GetRuns(1000)
	return runs
}

func writeScheduleExportOutput(bundles []scheduleExportBundle, opts scheduleExportOpts) error {
	switch opts.Format {
	case constants.OutputYAML:
		return writeScheduleExportYAML(bundles, opts.FilePath)
	case "sqlite":
		return writeScheduleExportSQLite(bundles, opts.FilePath)
	case "zip":
		return writeScheduleExportZIP(bundles, opts.FilePath)
	default:
		return writeScheduleExportJSON(bundles, opts.FilePath)
	}
}

func writeScheduleExportJSON(bundles []scheduleExportBundle, filePath string) error {
	raw, err := json.MarshalIndent(bundles, "", constants.JSONIndent)
	if err != nil {
		return err
	}
	if filePath == "" {
		fmt.Println(string(raw))
		return nil
	}
	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	if err := os.WriteFile(filePath, raw, constants.FilePermission); err != nil {
		return apperror.WrapSimple(err, "write json export")
	}
	printExportSuccessBanner(filePath, len(bundles), "JSON")
	return nil
}

func writeScheduleExportYAML(bundles []scheduleExportBundle, filePath string) error {
	raw, err := yaml.Marshal(bundles)
	if err != nil {
		return err
	}
	if filePath == "" {
		fmt.Println(string(raw))
		return nil
	}
	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	if err := os.WriteFile(filePath, raw, constants.FilePermission); err != nil {
		return apperror.WrapSimple(err, "write yaml export")
	}
	printExportSuccessBanner(filePath, len(bundles), "YAML")
	return nil
}

func writeScheduleExportSQLite(bundles []scheduleExportBundle, filePath string) error {
	if filePath == "" {
		filePath = "schedules_export.db"
	}
	_ = os.Remove(filePath)
	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	conn, err := sql.Open("sqlite", filePath)
	if err != nil {
		return apperror.WrapSimple(err, "create export sqlite db")
	}
	defer conn.Close()
	if err := populateExportSQLite(conn, bundles); err != nil {
		return err
	}
	printExportSuccessBanner(filePath, len(bundles), "SQLite")
	return nil
}

func populateExportSQLite(conn *sql.DB, bundles []scheduleExportBundle) error {
	_, _ = conn.Exec(store.SQLCreateSchedulerTasksTable)
	_, _ = conn.Exec(`CREATE TABLE IF NOT EXISTS schedule_logs (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    schedule_name TEXT,
	    run_number INTEGER,
	    trigger_type TEXT,
	    runner_user TEXT,
	    started_at TEXT,
	    finished_at TEXT,
	    duration_ms INTEGER,
	    is_success INTEGER,
	    exit_code INTEGER,
	    output TEXT,
	    error_msg TEXT,
	    created_at TEXT
	);`)
	for _, b := range bundles {
		insertBundleIntoSQLite(conn, b)
	}
	return nil
}

func insertBundleIntoSQLite(conn *sql.DB, b scheduleExportBundle) {
	qTask := `INSERT INTO scheduler_tasks (name, slug, db_path, macro_name, command_line, interval_val, delay_val, is_enabled, is_scheduled, has_delay, is_startup, run_count, last_run_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, _ = conn.Exec(qTask, b.Task.Name, b.Task.Slug, b.Task.DBPath, b.Task.MacroName, b.Task.CommandLine, b.Task.IntervalVal, b.Task.DelayVal, b.Task.IsEnabled, b.Task.IsScheduled, b.Task.HasDelay, b.Task.IsStartup, b.Task.RunCount, b.Task.LastRunAt)

	qRun := `INSERT INTO schedule_logs (schedule_name, run_number, trigger_type, runner_user, started_at, finished_at, duration_ms, is_success, exit_code, output, error_msg, created_at)
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, r := range b.Runs {
		isSuccessInt := 0
		if r.IsSuccess {
			isSuccessInt = 1
		}
		_, _ = conn.Exec(qRun, b.Task.Name, r.RunNumber, r.TriggerType, r.RunnerUser, r.StartedAt, r.FinishedAt, r.DurationMS, isSuccessInt, r.ExitCode, r.Output, r.ErrorMsg, r.CreatedAt)
	}
}

func writeScheduleExportZIP(bundles []scheduleExportBundle, filePath string) error {
	if filePath == "" {
		filePath = "schedules_export.zip"
	}
	_ = os.MkdirAll(filepath.Dir(filePath), constants.DirPermission)
	outFile, err := os.Create(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "create zip file")
	}
	defer outFile.Close()
	zw := zip.NewWriter(outFile)
	defer zw.Close()
	return addBundlesToZIP(zw, bundles, filePath)
}

func addBundlesToZIP(zw *zip.Writer, bundles []scheduleExportBundle, filePath string) error {
	for _, b := range bundles {
		raw, _ := json.MarshalIndent(b, "", constants.JSONIndent)
		f, err := zw.Create(b.Task.Slug + ".json")
		if err == nil {
			_, _ = f.Write(raw)
		}
	}
	printExportSuccessBanner(filePath, len(bundles), "ZIP")
	return nil
}

func printExportSuccessBanner(filePath string, count int, format string) {
	fmt.Printf("\n  \033[1;92m✔ Exported %d schedule(s)\033[0m to %s (\033[1m%s\033[0m)\n\n", count, format, filePath)
}
