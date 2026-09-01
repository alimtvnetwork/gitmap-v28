// Package cmd — schedule_import.go: imports scheduled tasks and their split DB run logs
// from JSON, YAML, SQLite databases, or ZIP archives with format auto-inference and except filters.
package cmd

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

func runScheduleImport(args []string) error {
	opts := parseScheduleExportOpts(args)
	if opts.FilePath == "" && opts.TargetName != "" && fileExists(opts.TargetName) {
		opts.FilePath = opts.TargetName
		opts.TargetName = ""
	}
	if opts.FilePath == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap schedule import <file> [-except \"name1, name2\"]\n")
		return apperror.NewSimple("import file path required", "E6011")
	}
	bundles, err := parseImportFileBundles(opts.FilePath)
	if err != nil {
		return err
	}
	return importBundlesIntoStore(bundles, opts.ExceptList)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func parseImportFileBundles(filePath string) ([]scheduleExportBundle, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return parseImportYAML(filePath)
	case ".db", ".sqlite", ".sqlite3":
		return parseImportSQLite(filePath)
	case ".zip":
		return parseImportZIP(filePath)
	default:
		return parseImportJSON(filePath)
	}
}

func parseImportJSON(filePath string) ([]scheduleExportBundle, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read import json")
	}
	var bundles []scheduleExportBundle
	if err := json.Unmarshal(raw, &bundles); err == nil && len(bundles) > 0 {
		return bundles, nil
	}
	var single scheduleExportBundle
	if err := json.Unmarshal(raw, &single); err == nil && single.Task.Name != "" {
		return []scheduleExportBundle{single}, nil
	}
	return nil, apperror.NewSimple("invalid json schedule export format", "E6012")
}

func parseImportYAML(filePath string) ([]scheduleExportBundle, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read import yaml")
	}
	var bundles []scheduleExportBundle
	if err := yaml.Unmarshal(raw, &bundles); err == nil && len(bundles) > 0 {
		return bundles, nil
	}
	var single scheduleExportBundle
	if err := yaml.Unmarshal(raw, &single); err == nil && single.Task.Name != "" {
		return []scheduleExportBundle{single}, nil
	}
	return nil, apperror.NewSimple("invalid yaml schedule export format", "E6013")
}

func parseImportSQLite(filePath string) ([]scheduleExportBundle, error) {
	conn, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open import sqlite db")
	}
	defer conn.Close()
	tasks, err := queryImportTasksFromDB(conn)
	if err != nil {
		return nil, err
	}
	var bundles []scheduleExportBundle
	for _, t := range tasks {
		runs := queryImportRunsFromDB(conn, t.Name)
		bundles = append(bundles, scheduleExportBundle{Task: t, Runs: runs})
	}
	return bundles, nil
}

func queryImportTasksFromDB(conn *sql.DB) ([]store.SchedulerTask, error) {
	q := `SELECT id, name, COALESCE(slug,''), COALESCE(db_path,''), COALESCE(macro_name,''), COALESCE(command_line,''), interval_val, delay_val, is_enabled, is_scheduled, has_delay, is_startup, run_count, COALESCE(last_run_at,''), created_at 
	      FROM scheduler_tasks`
	rows, err := conn.Query(q)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query import scheduler_tasks")
	}
	defer rows.Close()
	return parseImportTaskRows(rows), nil
}

func parseImportTaskRows(rows *sql.Rows) []store.SchedulerTask {
	var list []store.SchedulerTask
	for rows.Next() {
		var t store.SchedulerTask
		var isEnabledInt int
		err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.DBPath, &t.MacroName, &t.CommandLine, &t.IntervalVal, &t.DelayVal, &isEnabledInt, &t.IsScheduled, &t.HasDelay, &t.IsStartup, &t.RunCount, &t.LastRunAt, &t.CreatedAt)
		if err == nil {
			t.IsEnabled = isEnabledInt == 1
			list = append(list, t)
		}
	}
	return list
}

func queryImportRunsFromDB(conn *sql.DB, taskName string) []store.ScheduleRunRecord {
	q := `SELECT id, run_number, trigger_type, runner_user, started_at, finished_at, duration_ms, is_success, exit_code, output, error_msg, created_at 
	      FROM schedule_logs WHERE schedule_name = ? ORDER BY id ASC`
	rows, err := conn.Query(q, taskName)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return parseImportRunRows(rows)
}

func parseImportRunRows(rows *sql.Rows) []store.ScheduleRunRecord {
	var list []store.ScheduleRunRecord
	for rows.Next() {
		var r store.ScheduleRunRecord
		var isSuccessInt int
		err := rows.Scan(&r.ID, &r.RunNumber, &r.TriggerType, &r.RunnerUser, &r.StartedAt, &r.FinishedAt, &r.DurationMS, &isSuccessInt, &r.ExitCode, &r.Output, &r.ErrorMsg, &r.CreatedAt)
		if err == nil {
			r.IsSuccess = isSuccessInt == 1
			list = append(list, r)
		}
	}
	return list
}

func parseImportZIP(filePath string) ([]scheduleExportBundle, error) {
	zr, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open zip file")
	}
	defer zr.Close()
	var bundles []scheduleExportBundle
	for _, f := range zr.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".json") {
			continue
		}
		b := readBundleFromZipFile(f)
		if b.Task.Name != "" {
			bundles = append(bundles, b)
		}
	}
	return bundles, nil
}

func readBundleFromZipFile(f *zip.File) scheduleExportBundle {
	rc, err := f.Open()
	if err != nil {
		return scheduleExportBundle{}
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	var b scheduleExportBundle
	_ = json.Unmarshal(data, &b)
	return b
}

func importBundlesIntoStore(bundles []scheduleExportBundle, exceptList []string) error {
	db, err := openSchedulerDB()
	if err != nil {
		return err
	}
	defer db.Close()
	importedCount := 0
	for _, b := range bundles {
		if isNameInExceptList(b.Task.Name, exceptList) {
			continue
		}
		if err := persistImportedBundle(db, b); err == nil {
			importedCount++
		}
	}
	fmt.Printf("\n  \033[1;92m✔ Successfully imported %d schedule(s)\033[0m into root and split databases\n\n", importedCount)
	return nil
}

func isNameInExceptList(name string, exceptList []string) bool {
	for _, ex := range exceptList {
		if strings.EqualFold(name, ex) {
			return true
		}
	}
	return false
}

func persistImportedBundle(db *store.DB, b scheduleExportBundle) error {
	if b.Task.Slug == "" {
		b.Task.Slug = store.ScheduleSlug(b.Task.Name)
	}
	b.Task.DBPath = store.ScheduleDBPath(b.Task.Slug)
	if err := db.InsertSchedule(b.Task); err != nil {
		return err
	}
	splitDB, err := store.OpenScheduleSplitDB(b.Task.Slug)
	if err != nil {
		return err
	}
	defer splitDB.Close()
	_ = splitDB.SaveConfig(store.ScheduleConfig{
		Name:        b.Task.Name,
		Slug:        b.Task.Slug,
		MacroName:   b.Task.MacroName,
		CommandLine: b.Task.CommandLine,
		IntervalVal: b.Task.IntervalVal,
		DelayVal:    b.Task.DelayVal,
		IsEnabled:   b.Task.IsEnabled,
		IsStartup:   b.Task.IsStartup,
	})
	for _, r := range b.Runs {
		_ = splitDB.RecordRun(r)
	}
	return nil
}
