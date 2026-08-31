package macro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ExecutionReport captures structured execution details for JSON/YAML serialization.
type ExecutionReport struct {
	Macro          string          `json:"macro" yaml:"macro"`
	Status         string          `json:"status" yaml:"status"`
	TotalSteps     int             `json:"totalSteps" yaml:"totalSteps"`
	ExecutedSteps  int             `json:"executedSteps" yaml:"executedSteps"`
	FailedSteps    int             `json:"failedSteps" yaml:"failedSteps"`
	ElapsedSeconds float64         `json:"elapsedSeconds" yaml:"elapsedSeconds"`
	StartedAt      time.Time       `json:"startedAt" yaml:"startedAt"`
	CompletedAt    time.Time       `json:"completedAt" yaml:"completedAt"`
	OutputFile     string          `json:"outputFile,omitempty" yaml:"outputFile,omitempty"`
	Steps          []StepExecution `json:"steps" yaml:"steps"`
}

// StepExecution represents the execution outcome of an individual step.
type StepExecution struct {
	StepNum        int      `json:"stepNum" yaml:"stepNum"`
	CommandLine    string   `json:"commandLine" yaml:"commandLine"`
	WorkingDir     string   `json:"workingDir" yaml:"workingDir"`
	Status         string   `json:"status" yaml:"status"`
	ExitCode       int      `json:"exitCode" yaml:"exitCode"`
	ElapsedSeconds float64  `json:"elapsedSeconds" yaml:"elapsedSeconds"`
	Logs           []string `json:"logs" yaml:"logs"`
	Error          string   `json:"error,omitempty" yaml:"error,omitempty"`
	ErrorLogs      []string `json:"errorLogs,omitempty" yaml:"errorLogs,omitempty"`
}

// NewExecutionReport initializes an execution report.
func NewExecutionReport(macroName string, totalSteps int, start time.Time) *ExecutionReport {
	return &ExecutionReport{
		Macro:         macroName,
		Status:        "success",
		TotalSteps:    totalSteps,
		ExecutedSteps: 0,
		FailedSteps:   0,
		StartedAt:     start,
		Steps:         make([]StepExecution, 0, totalSteps),
	}
}

// Finalize completes the report metrics.
func (r *ExecutionReport) Finalize(completedAt time.Time, hasErrors bool) {
	r.CompletedAt = completedAt
	r.ElapsedSeconds = completedAt.Sub(r.StartedAt).Seconds()
	if hasErrors {
		r.Status = "failed"
	}
}

// FormatReport serializes the report to string in either JSON or YAML format.
func FormatReport(r *ExecutionReport, isYAML bool) (string, error) {
	if isYAML {
		return formatYAMLReport(r)
	}
	return formatJSONReport(r)
}

func formatYAMLReport(r *ExecutionReport) (string, error) {
	data, err := yaml.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func formatJSONReport(r *ExecutionReport) (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SaveReportToFile writes the formatted content to a target file.
func SaveReportToFile(filePath, content string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}
	dir := filepath.Dir(absPath)
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return "", mkErr
	}
	if writeErr := os.WriteFile(absPath, []byte(content), 0644); writeErr != nil {
		return "", writeErr
	}
	return absPath, nil
}

// HandleReportOutput outputs the report to terminal and/or saves to file.
func HandleReportOutput(r *ExecutionReport, opts ExecOptions) error {
	isYAML := determineIsYAML(opts)
	content, err := FormatReport(r, isYAML)
	if err != nil {
		return fmt.Errorf("failed formatting report: %w", err)
	}
	if len(opts.FilePath) > 0 {
		return writeReportAndNotify(r, opts.FilePath, content)
	}
	if opts.JSON || opts.YAML {
		fmt.Println(content)
	}
	return nil
}

func determineIsYAML(opts ExecOptions) bool {
	if opts.YAML {
		return true
	}
	lower := strings.ToLower(opts.FilePath)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}

func writeReportAndNotify(r *ExecutionReport, filePath, content string) error {
	savedPath, saveErr := SaveReportToFile(filePath, content)
	if saveErr != nil {
		return fmt.Errorf("failed saving report to %s: %w", filePath, saveErr)
	}
	r.OutputFile = savedPath
	fmt.Println(content)
	fmt.Printf("\n  %s✔ Macro execution report saved to:%s %s%s%s\n\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorCyan, savedPath, constants.ColorReset)
	return nil
}
