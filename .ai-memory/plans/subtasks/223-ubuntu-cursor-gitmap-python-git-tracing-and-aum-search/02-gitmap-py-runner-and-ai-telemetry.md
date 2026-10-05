# Subtask Plan 02: Native GitMap Python Runner & AI Telemetry Tracing

- **Subtask Slug:** `02-gitmap-py-runner-and-ai-telemetry`
- **Parent Task:** `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search`
- **Target Files:**
  - `cli/cmdpy/py_cmd.go`
  - `cli/cmdpy/py_exec.go`
  - `cli/cmdpy/py_record.go`
  - `cli/cmdpy/py_cmd_test.go`
  - `cli/cmd/rootutility.go`
  - `cli/cmd/rootcore.go`

---

## 1. Overview & Objectives

In modern development and multi-agent workflows, automation scripts (such as version bumping, spec synchronizing, database migration, and test runner scripts) are authored in Python. Currently, executing them directly outside GitMap bypasses developer execution history, command timing benchmarks, and AI instruction telemetry.

This subtask implements:
1. Native GitMap Python runner package `cli/cmdpy/` supporting both `gitmap py` and `gitmap python`.
2. First-class support for:
   - Inline script execution: `gitmap py -c "print('hello')"`
   - File execution with arguments: `gitmap py script.py arg1 arg2`
   - Arbitrary flags & REPL: `gitmap py --version`, `gitmap py -m pytest`
3. Cross-platform Python discovery (resolving active virtual environment `$VIRTUAL_ENV`, `python3`, `python`, or Windows `py`).
4. Real-time streaming of `stdin`, `stdout`, and `stderr`, with transparent propagation of child process exit codes.
5. Dual SQLite telemetry recording:
   - Split-DB `commands.db` (`CommandHistory` table) for developer activity heatmap visualization.
   - AI instruction DB (`store.RecordAiExecution`) for agent audit tracking.
6. Integration into root dispatch tables in `cli/cmd/rootutility.go` and `cli/cmd/rootcore.go`.
7. Comprehensive unit tests in `cli/cmdpy/py_cmd_test.go`.

---

## 2. Step-by-Step Implementation Instructions

### Step 1: Create `cli/cmdpy/py_cmd.go`

Create the CLI dispatch and argument parsing entrypoint (`<= 95` lines):
```go
package cmdpy

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintPyHelp renders help information for gitmap py / gitmap python.
func PrintPyHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║             gitmap py / python - Python Execution Runner         ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap py -c \"<code>\"             Execute inline Python code")
	fmt.Println("    gitmap py <script.py> [args...]   Execute Python script file")
	fmt.Println("    gitmap py [args...]               Pass arguments directly to Python")
	fmt.Println("    gitmap python [args...]           Canonical alias for gitmap py")
	fmt.Println()
}

func isPyHelpRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	arg := strings.ToLower(args[0])
	return arg == "help" || arg == "--help" || arg == "-h"
}

// RunPy executes the Python runner CLI command.
func RunPy(args []string) error {
	if isPyHelpRequest(args) {
		PrintPyHelp()
		return nil
	}
	binary, err := resolvePythonBinary()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	exitCode, runErr := executePythonStreaming(binary, args, cwd)
	if exitCode != 0 && runErr == nil {
		os.Exit(exitCode)
	}
	if runErr != nil {
		return apperror.WrapSimple(runErr, "python execution failed")
	}
	return nil
}
```

### Step 2: Create `cli/cmdpy/py_exec.go`

Implement interpreter resolution and streaming process runner (`<= 90` lines):
```go
package cmdpy

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolvePythonBinary() (string, error) {
	if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
		venvBin := filepath.Join(venv, "bin", "python")
		if runtime.GOOS == "windows" {
			venvBin = filepath.Join(venv, "Scripts", "python.exe")
		}
		if _, err := os.Stat(venvBin); err == nil {
			return venvBin, nil
		}
	}
	candidates := []string{"python3", "python"}
	if runtime.GOOS == "windows" {
		candidates = append(candidates, "py")
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil && path != "" {
			return path, nil
		}
	}
	return "", apperror.NewWithDetails("cmd.py.resolve", "E1032", "no Python interpreter found on PATH; please install Python 3", "cmdpy", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func executePythonStreaming(binary string, args []string, dir string) (int, error) {
	cmd := exec.CommandContext(context.Background(), binary, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	durationMs := time.Since(start).Milliseconds()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	recordPyTelemetry(binary, args, dir, durationMs, exitCode, err)
	return exitCode, err
}
```

### Step 3: Create `cli/cmdpy/py_record.go`

Implement dual telemetry recording (`<= 85` lines):
```go
package cmdpy

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func recordPyTelemetry(binary string, args []string, dir string, durationMs int64, exitCode int, err error) {
	cmdLine := "py " + strings.Join(args, " ")
	recordHistoryTelemetry(cmdLine, exitCode, durationMs)
	recordAiTelemetry(cmdLine, args, dir, durationMs, exitCode, err)
}

func recordHistoryTelemetry(cmdLine string, exitCode int, durationMs int64) {
	histDB, err := store.OpenCommandHistorySplitDB("")
	if err != nil {
		return
	}
	defer histDB.Close()
	_ = histDB.InsertCommandRecord(cmdLine, "py", exitCode, durationMs)
}

func recordAiTelemetry(cmdLine string, args []string, dir string, durationMs int64, exitCode int, err error) {
	argsJSON, _ := json.Marshal(args)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	_ = store.RecordAiExecution(
		"python_runner",
		cmdLine,
		string(argsJSON),
		dir,
		"127.0.0.1",
		int(durationMs),
		exitCode,
		"",
		errMsg,
		exitCode == 0,
	)
}
```

### Step 4: Create `cli/cmdpy/py_cmd_test.go`

Create comprehensive unit tests verifying CLI dispatch, argument resolution, and help text:
```go
package cmdpy

import (
	"testing"
)

func TestIsPyHelpRequest(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{[]string{}, false},
		{[]string{"help"}, true},
		{[]string{"--help"}, true},
		{[]string{"-h"}, true},
		{[]string{"-c", "print(1)"}, false},
		{[]string{"script.py"}, false},
	}
	for _, tc := range tests {
		actual := isPyHelpRequest(tc.args)
		if actual != tc.expected {
			t.Errorf("isPyHelpRequest(%v) = %v, expected %v", tc.args, actual, tc.expected)
		}
	}
}

func TestResolvePythonBinary(t *testing.T) {
	bin, err := resolvePythonBinary()
	if err != nil {
		t.Logf("resolvePythonBinary notice: %v", err)
	} else if bin == "" {
		t.Errorf("expected non-empty binary string")
	}
}
```

### Step 5: Update Root Dispatch Tables

1. In `cli/cmd/rootutility.go`:
   Add `[]string{"py", "python"}` to `utilityToolEntries()`:
   ```go
   {[]string{"py", "python"}, func() error { return cmdpy.RunPy(argsTail()) }},
   ```

2. In `cli/cmd/rootcore.go`:
   Ensure `py` and `python` are registered in routing validation and help menus.

---

## 3. Verification & Acceptance Checklist

- [ ] `cli/cmdpy/py_cmd.go` created and adheres to `<= 100` lines standard.
- [ ] `cli/cmdpy/py_exec.go` created with cross-platform Python binary resolution and streaming stdio.
- [ ] `cli/cmdpy/py_record.go` created with dual SQLite persistence to `commands.db` and `sql.db`.
- [ ] `cli/cmdpy/py_cmd_test.go` passes all test cases.
- [ ] `cli/cmd/rootutility.go` routes `gitmap py` and `gitmap python`.
- [ ] `gitmap py -c "print(1234)"` outputs `1234`.
- [ ] Telemetry entries successfully logged to `CommandHistory` with `CommandName = "py"`.
- [ ] Process exit codes are propagated faithfully.
