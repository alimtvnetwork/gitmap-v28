package macro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func tryResolveBaseLogsDir() (string, bool) {
	baseDir, err := resolveWritableMacroDir()
	if err != nil {
		return "", false
	}

	target := filepath.Join(baseDir, "logs")
	if probeAndPrepareDir(target) {
		return target, true
	}

	return "", false
}

// ResolveMacroLogsDir finds or creates a writable directory for macro logs.
func ResolveMacroLogsDir() MacroLogResult {
	if target, isFound := tryResolveBaseLogsDir(); isFound {
		return result.Ok(target)
	}

	tempLogs := filepath.Join(resolveTempMacroDir(), "logs")
	if probeAndPrepareDir(tempLogs) {
		return result.Ok(tempLogs)
	}

	appErr := apperror.NewWithDetails(
		"macro.logs",
		"E5002",
		"no writable macro logs directory found",
		"macro.logger",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)

	return result.Fail[string](appErr)
}

func sanitizeMacroName(name string) string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}

		return '_'
	}, name)

	if len(cleaned) == 0 {
		return "macro"
	}

	return cleaned
}

// BuildFailureLogPath calculates the target file path for a macro failure log.
func BuildFailureLogPath(macroName, explicitPath string, now time.Time) (string, error) {
	if len(explicitPath) > 0 {
		return filepath.Abs(explicitPath)
	}

	dirRes := ResolveMacroLogsDir()
	if dirRes.IsFailure() {
		return "", dirRes.AppError()
	}

	sanitized := sanitizeMacroName(macroName)
	fname := fmt.Sprintf("%s-%s.log", sanitized, now.Format("20060102-150405"))

	return filepath.Join(dirRes.Data, fname), nil
}

func writeFailureHeader(sb *strings.Builder) {
	sb.WriteString("================================================================================\n")
	sb.WriteString("GITMAP MACRO FAILURE DIAGNOSTIC REPORT\n")
	sb.WriteString("================================================================================\n")
}

func writeFailureMetadata(sb *strings.Builder, ctx MacroFailureContext) {
	sb.WriteString(fmt.Sprintf("Timestamp:       %s\n", ctx.Timestamp.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Macro Name:      %s\n", ctx.MacroName))
	sb.WriteString(fmt.Sprintf("Step Number:     %d of %d\n", ctx.Step.StepNum, ctx.TotalSteps))
	sb.WriteString(fmt.Sprintf("Command Line:    %s\n", ctx.Step.CommandLine))
	sb.WriteString(fmt.Sprintf("Working Dir:     %s\n", ctx.Step.WorkingDir))
	sb.WriteString(fmt.Sprintf("Exit Code:       %d\n", ctx.Step.ExitCode))
	sb.WriteString(fmt.Sprintf("Duration:        %.2fs\n", ctx.Step.ElapsedSeconds))
	sb.WriteString(fmt.Sprintf("Error:           %s\n\n", ctx.Step.Error))
}

func writeFailureDiagnostics(sb *strings.Builder, errLogs []string) {
	sb.WriteString("--- STDERR DIAGNOSTICS ---\n")
	for _, line := range errLogs {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
}

func writeFailureStdoutContext(sb *strings.Builder, logs []string) {
	sb.WriteString("--- RECENT STDOUT CONTEXT ---\n")
	start := 0
	if len(logs) > 20 {
		start = len(logs) - 20
	}

	for i := start; i < len(logs); i++ {
		sb.WriteString(logs[i])
		sb.WriteString("\n")
	}

	sb.WriteString("================================================================================\n")
}

// FormatFailureLogEntry constructs the formatted diagnostic entry string.
func FormatFailureLogEntry(ctx MacroFailureContext) string {
	var sb strings.Builder
	writeFailureHeader(&sb)
	writeFailureMetadata(&sb, ctx)
	writeFailureDiagnostics(&sb, ctx.Step.ErrorLogs)
	writeFailureStdoutContext(&sb, ctx.Step.Logs)

	return sb.String()
}

func appendLogFile(logPath, content string) error {
	dir := filepath.Dir(logPath)
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "mkdir logs dir")
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return apperror.WrapSimple(err, "open log file")
	}

	defer f.Close()
	_, writeErr := f.WriteString(content)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "write log content")
	}

	return nil
}

// RecordStepFailureLog formats, appends, and returns the log file path.
func RecordStepFailureLog(ctx MacroFailureContext, explicitPath string) (string, error) {
	logPath, err := BuildFailureLogPath(ctx.MacroName, explicitPath, ctx.Timestamp)
	if err != nil {
		return "", apperror.WrapSimple(err, "BuildFailureLogPath")
	}

	content := FormatFailureLogEntry(ctx)
	if writeErr := appendLogFile(logPath, content); writeErr != nil {
		return "", apperror.WrapSimple(writeErr, "appendLogFile")
	}

	return logPath, nil
}
