package cmdcursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secretsresolver"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const migrationScriptRelative = "repo-secrets/05-scripts/migrate-cursor-memories-conversations.py"

func resolveMigrationScriptPath() (string, error) {
	candidates := []string{
		migrationScriptRelative,
		filepath.Join(secretsresolver.ResolveRepoSecretsRoot(), "05-scripts", "migrate-cursor-memories-conversations.py"),
		filepath.Join("..", migrationScriptRelative),
	}
	for _, cand := range candidates {
		info, err := os.Stat(cand)
		if err == nil && !info.IsDir() {
			return cand, nil
		}
	}

	return "", apperror.NewNotFound("resolveMigrationScriptPath", "E1036", "migration script not found: "+migrationScriptRelative)
}

func buildMigrationPyArgs(script string, opts cursorMigrateOptions) []string {
	pyArgs := []string{script, "--sync", "--node", opts.targetNode}
	if opts.exclude != "" {
		pyArgs = append(pyArgs, "--exclude", opts.exclude)
	}
	if opts.outputPath != "" {
		pyArgs = append(pyArgs, "--output", opts.outputPath)
	}
	if opts.isSkipSettings {
		pyArgs = append(pyArgs, "--skip-settings")
	}
	if opts.isSettingsOnly {
		pyArgs = append(pyArgs, "--settings-only")
	}
	if opts.isDryRun {
		pyArgs = append(pyArgs, "--dry-run")
	}
	if opts.isNoBackup {
		pyArgs = append(pyArgs, "--no-backup")
	}
	if opts.isForce {
		pyArgs = append(pyArgs, "--force")
	}

	return pyArgs
}

func recordMigrateTelemetry(cmdLine string, args []string, durationMs int64, runErr error) {
	exitCode, errMsg := 0, ""
	if runErr != nil {
		exitCode, errMsg = 1, runErr.Error()
	}
	if histDB, err := store.OpenCommandHistorySplitDB(""); err == nil {
		_ = histDB.InsertCommandRecord(cmdLine, "cursor-migrate", exitCode, durationMs)
		histDB.Close()
	}

	argsJSON, _ := json.Marshal(args)
	cwd, _ := os.Getwd()
	isSuccess := exitCode == 0

	_ = store.RecordAiExecution("cursor_migrate", cmdLine, string(argsJSON), cwd, "127.0.0.1", int(durationMs), exitCode, "", errMsg, isSuccess)
}

func executeCursorMigration(opts cursorMigrateOptions, rawArgs []string) error {
	scriptPath, err := resolveMigrationScriptPath()
	if err != nil {
		return err
	}

	pyArgs := buildMigrationPyArgs(scriptPath, opts)
	cmdLine := "gitmap cursor migrate " + strings.Join(rawArgs, " ")

	fmt.Printf("\n%s● Launching Cursor IDE Migration to node %s...%s\n", constants.ColorCyan, opts.targetNode, constants.ColorReset)

	start := time.Now()
	runErr := cmdpy.RunPy(pyArgs)
	durationMs := time.Since(start).Milliseconds()

	recordMigrateTelemetry(cmdLine, pyArgs, durationMs, runErr)
	return runErr
}
