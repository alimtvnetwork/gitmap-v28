package cmdrun

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// Run routes gitmap run invocations between subcommands and script execution.
func Run(args []string) error {
	if len(args) == 0 {
		PrintRunHelp()
		return apperror.NewValidationError("missing script or command for run")
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "-h", "--help", "help":
		PrintRunHelp()
		return nil
	case "errors", "err":
		return RunErrorsCmd(args[1:])
	case "clear-errors":
		return RunClearErrorsCmd(args[1:])
	case "history":
		return RunHistoryCmd(args[1:])
	default:
		return RunFile(args[0], args[1:])
	}
}

// RunFile resolves, audits, and executes a script target with error tracking.
func RunFile(targetArg string, userArgs []string) error {
	target, resolveErr := ResolveRunTarget(targetArg)
	if resolveErr != nil {
		return resolveErr
	}

	taskId, _ := EnqueueRunTaskAudit(target, userArgs)
	opts := RunOptions{Args: userArgs, WorkingDir: "."}
	result, execErr := ExecuteTarget(target, opts)

	if result.ExitCode != 0 || execErr != nil {
		return handleRunFailure(taskId, target, userArgs, result, execErr)
	}

	_ = CompleteRunTaskAudit(taskId, target.ResolvedPath, strings.Join(userArgs, " "), result.DurationMs)

	return nil
}

func handleRunFailure(taskId string, target *RunTarget, userArgs []string, result *RunResult, execErr error) error {
	errRec := RunErrorRecord{
		ErrorID:       taskId,
		FilePath:      target.ResolvedPath,
		FileExtension: target.Extension,
		Interpreter:   target.Interpreter,
		ExitCode:      result.ExitCode,
		DurationMs:    result.DurationMs,
		ErrorMessage:  resolveRunErrorMessage(execErr, result.ExitCode),
		StdoutSnippet: result.Stdout,
		StderrSnippet: result.Stderr,
		ExecutedArgs:  strings.Join(userArgs, " "),
	}

	_ = RecordRunError(errRec)

	if execErr != nil {
		return apperror.WrapSimple(execErr, "script execution failed")
	}

	return apperror.NewValidationError(fmt.Sprintf("script exited with code %d", result.ExitCode))
}

func resolveRunErrorMessage(execErr error, exitCode int) string {
	if execErr != nil {
		return execErr.Error()
	}

	return fmt.Sprintf("process terminated with non-zero exit code %d", exitCode)
}

// RunHistoryCmd queries and displays past script runs from TaskHistory.
func RunHistoryCmd(args []string) error {
	records, err := QueryRunTaskHistory(50)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Printf("  %sINFO%s No script execution history recorded.\n", constants.ColorCyan, constants.ColorReset)
		return nil
	}

	cfg := buildRunHistoryTableConfig(records)
	termout.PrintTable(cfg)

	return nil
}

func buildRunHistoryTableConfig(records []model.TaskHistoryRecord) termout.TableConfig {
	columns := []termout.Column{
		{Title: "TASK ID", MaxWidth: 26, Align: termout.AlignLeft},
		{Title: "ACTION", MaxWidth: 12, Align: termout.AlignLeft},
		{Title: "TARGET", MaxWidth: 36, Align: termout.AlignLeft},
		{Title: "STATUS", MaxWidth: 12, Align: termout.AlignLeft},
		{Title: "EXECUTED AT", MaxWidth: 22, Align: termout.AlignLeft},
	}

	rows := make([]termout.Row, 0, len(records))
	for _, r := range records {
		rows = append(rows, termout.Row{
			Cells: []string{r.TaskId, r.Action, r.Target, r.Status, r.ExecutedAt},
			Color: constants.ColorGreen,
		})
	}

	return termout.TableConfig{
		Columns:     columns,
		Rows:        rows,
		HeaderColor: constants.ColorCyan,
		BorderColor: constants.ColorDim,
		HasBorders:  true,
	}
}

// PrintRunHelp renders CLI usage help for gitmap run.
func PrintRunHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║        gitmap run - Polyglot Script & Macro Execution Suite       ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap run <script-file> [args...]    Auto-detect interpreter and execute")
	fmt.Println("    gitmap run <macro-name> [flags]       Execute saved interactive macro")
	fmt.Println("    gitmap run errors                     List recent script execution failures")
	fmt.Println("    gitmap run clear-errors               Purge execution error records")
	fmt.Println("    gitmap run history                    Inspect execution audit history")
	fmt.Println()
	fmt.Println("  Supported Script Types:")
	fmt.Println("    .py (Python), .ps1 (PowerShell), .sh (Bash), .js (Node), .ts (Bun/tsx), .go (Go run)")
	fmt.Println()
}
