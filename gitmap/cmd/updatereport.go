package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// reportErrorsConfig captures the resolved --report-errors and
// --debug-repo-detect settings. Both share the same flag-propagation and
// env-injection lifecycle, so they're bundled into one carrier struct.
type reportErrorsConfig struct {
	reportEnabled   bool
	format          string
	file            string
	debugRepoDetect bool
}

// resolveReportErrors reads --report-errors / --report-errors-file and
// --debug-repo-detect from os.Args[2:]. Returns a zero config when none
// are present.
func resolveReportErrors() reportErrorsConfig {
	cfg := reportErrorsConfig{debugRepoDetect: hasFlag(constants.FlagDebugRepoDetect)}

	value := getFlagValue(constants.FlagReportErrors)
	if len(value) == 0 {
		return cfg
	}

	if value != constants.ReportErrorsJSON {
		fmt.Fprintf(os.Stderr, constants.ErrReportErrorsValue, value)
		os.Exit(1)
	}

	path := getFlagValue(constants.FlagReportErrorsFile)
	if len(path) == 0 {
		path = defaultReportPath()
	}

	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if err := ensureReportFile(path); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnReportErrorsCreate, path, err)

		return cfg
	}

	cfg.reportEnabled = true
	cfg.format = value
	cfg.file = path

	return cfg
}

// defaultReportPath builds a timestamped report path under the system temp dir.
func defaultReportPath() string {
	name := fmt.Sprintf("%s%s%s",
		constants.ReportErrorsFilePrefix,
		time.Now().Format("20060102-150405"),
		constants.ReportErrorsFileSuffix,
	)

	return filepath.Join(os.TempDir(), name)
}

// ensureReportFile creates (or truncates) the JSONL report file so the worker
// scripts can append entries without first checking existence.
func ensureReportFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	return f.Close()
}

// applyToHandoffArgs appends flags so the handoff worker sees them.
func (c reportErrorsConfig) applyToHandoffArgs(args []string) []string {
	if c.reportEnabled {
		args = append(args, constants.FlagReportErrors, c.format)
		args = append(args, constants.FlagReportErrorsFile, c.file)
	}

	if c.debugRepoDetect {
		args = append(args, constants.FlagDebugRepoDetect)
	}

	return args
}

// applyToEnv injects env vars consumed by run.ps1 / run.sh.
func (c reportErrorsConfig) applyToEnv(env []string) []string {
	if c.reportEnabled {
		env = append(env,
			fmt.Sprintf("%s=%s", constants.EnvReportErrorsFormat, c.format),
			fmt.Sprintf("%s=%s", constants.EnvReportErrorsFile, c.file),
		)
	}

	if c.debugRepoDetect {
		env = append(env, fmt.Sprintf("%s=1", constants.EnvDebugRepoDetect))
	}

	return env
}

// announce prints a one-line notice when reporting or debug is active.
func (c reportErrorsConfig) announce() {
	if c.reportEnabled {
		fmt.Printf(constants.MsgReportErrorsEnabled, c.file)
	}

	if c.debugRepoDetect {
		fmt.Print(constants.MsgDebugRepoDetectOn)
	}
}

// summarize prints a short summary of the report file after the update finishes.
func (c reportErrorsConfig) summarize() {
	if !c.reportEnabled {
		return
	}

	count := countReportEntries(c.file)
	suffix := "ies"
	if count == 1 {
		suffix = "y"
	}

	fmt.Printf(constants.MsgReportErrorsSummary, count, suffix, c.file)
}

// countReportEntries counts non-empty lines in the JSONL report file.
func countReportEntries(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	count := 0
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] != '\n' {
			continue
		}
		if i > start {
			count++
		}
		start = i + 1
	}
	if start < len(data) {
		count++
	}

	return count
}
