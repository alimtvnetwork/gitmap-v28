// Package cmd — macro_retry.go: executes macros or commands in a loop until success,
// with configurable sleep intervals, backoff strategies, and AI-ready failure diagnostics.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

type macroRetryConfig struct {
	Target     string
	Delay      time.Duration
	MaxRetries int
	Timeout    time.Duration
	Backoff    string
	IsAI       bool
	AIFilePath string
	IsJSON     bool
	IsYAML     bool
	IsVerbose  bool
}

type macroRunResult struct {
	Attempt   int           `json:"attempt" yaml:"attempt"`
	IsSuccess bool          `json:"isSuccess" yaml:"isSuccess"`
	Duration  time.Duration `json:"duration" yaml:"duration"`
	Output    string        `json:"output" yaml:"output"`
	ErrorMsg  string        `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`
	ExitCode  int           `json:"exitCode" yaml:"exitCode"`
}

func runMacroUntilSuccess(args []string) error {
	checkHelp("macro-retry", args)
	cfg := parseMacroRetryConfig(args)
	if cfg.Target == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro run-until-succeed <macro-name|\"cmd\"> [--sleep <sec>] [--max-retries <N>] [--ai]\n")
		return apperror.NewSimple("target macro or command string required", "E5003")
	}
	return executeRetryLoop(cfg)
}

func parseMacroRetryConfig(args []string) macroRetryConfig {
	cfg := macroRetryConfig{Delay: 5 * time.Second, Backoff: "fixed"}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case matchFlagWithVal(a, "--sleep", "--delay", "-d"):
			cfg.Delay = parseDurationArg(extractFlagValue(&i, args), 5*time.Second)
		case matchFlagWithVal(a, "--max-retries", "-n", "--limit"):
			cfg.MaxRetries, _ = strconv.Atoi(extractFlagValue(&i, args))
		case matchFlagWithVal(a, "--timeout", "-t"):
			cfg.Timeout = parseDurationArg(extractFlagValue(&i, args), 0)
		case matchFlagWithVal(a, "--backoff"):
			cfg.Backoff = extractFlagValue(&i, args)
		case matchFlagWithVal(a, "--error-file", "--ai-file"):
			cfg.AIFilePath = extractFlagValue(&i, args)
		case a == "--ai":
			cfg.IsAI = true
		case a == "--json":
			cfg.IsJSON = true
		case a == "--yaml":
			cfg.IsYAML = true
		case !strings.HasPrefix(a, "-") && cfg.Target == "":
			cfg.Target = a
		}
	}
	return cfg
}

func matchFlagWithVal(arg string, names ...string) bool {
	for _, n := range names {
		if arg == n || strings.HasPrefix(arg, n+"=") {
			return true
		}
	}
	return false
}

func extractFlagValue(idx *int, args []string) string {
	a := args[*idx]
	if strings.Contains(a, "=") {
		parts := strings.SplitN(a, "=", 2)
		return parts[1]
	}
	if *idx+1 < len(args) {
		*idx++
		return args[*idx]
	}
	return ""
}

func parseDurationArg(val string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	if sec, err := strconv.Atoi(val); err == nil {
		return time.Duration(sec) * time.Second
	}
	return fallback
}

func executeRetryLoop(cfg macroRetryConfig) error {
	startTime := time.Now()
	attempt := 1
	fmt.Printf("\n\033[1;96m▸ macro run-until-succeed\033[0m  target=\033[1m%q\033[0m (sleep=%v, backoff=%s)\n", cfg.Target, cfg.Delay, cfg.Backoff)
	for {
		res := executeSingleAttempt(cfg.Target, attempt)
		if res.IsSuccess {
			printRetrySuccess(attempt, time.Since(startTime))
			return nil
		}
		handleAttemptFailure(res, cfg, attempt)
		if isRetryExhausted(attempt, cfg.MaxRetries, startTime, cfg.Timeout) {
			return apperror.NewSimple(fmt.Sprintf("target %q failed after %d attempt(s)", cfg.Target, attempt), "E5004")
		}
		sleepNextDuration(cfg.Delay, cfg.Backoff, attempt)
		attempt++
	}
}

func isRetryExhausted(attempt, maxRetries int, start time.Time, timeout time.Duration) bool {
	if maxRetries > 0 && attempt >= maxRetries {
		return true
	}
	if timeout > 0 && time.Since(start) >= timeout {
		return true
	}
	return false
}

func executeSingleAttempt(target string, attempt int) macroRunResult {
	start := time.Now()
	fmt.Printf("\n  \033[1;94m● Attempt #%d\033[0m (%s)...\n", attempt, time.Now().Format("15:04:05"))
	loadedMacro, err := macro.LoadMacro(target)
	if err == nil && loadedMacro != nil {
		return runMacroDirect(loadedMacro, attempt, start)
	}
	return runShellCmdDirect(target, attempt, start)
}

func runMacroDirect(m *macro.Macro, attempt int, start time.Time) macroRunResult {
	execErr := macro.Execute(context.Background(), m, macro.ExecOptions{})
	dur := time.Since(start)
	if execErr == nil {
		return macroRunResult{Attempt: attempt, IsSuccess: true, Duration: dur}
	}
	return macroRunResult{
		Attempt:   attempt,
		IsSuccess: false,
		Duration:  dur,
		ErrorMsg:  execErr.Error(),
		ExitCode:  1,
	}
}

func runShellCmdDirect(cmdStr string, attempt int, start time.Time) macroRunResult {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	var outBuf bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = ioTeeAndBuffer(os.Stderr, &outBuf)

	err := cmd.Run()
	dur := time.Since(start)
	if err == nil {
		return macroRunResult{Attempt: attempt, IsSuccess: true, Duration: dur, Output: outBuf.String()}
	}
	exitCode := 1
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}
	return macroRunResult{
		Attempt:   attempt,
		IsSuccess: false,
		Duration:  dur,
		Output:    outBuf.String(),
		ErrorMsg:  err.Error(),
		ExitCode:  exitCode,
	}
}

func ioTeeAndBuffer(w *os.File, buf *bytes.Buffer) *os.File {
	return w
}

func handleAttemptFailure(res macroRunResult, cfg macroRetryConfig, attempt int) {
	fmt.Printf("  \033[1;91m✗ Failed\033[0m (exit code: %d, elapsed: %v)\n", res.ExitCode, res.Duration.Round(time.Millisecond))
	aiReport := buildAIFailureReport(cfg.Target, attempt, res)
	if cfg.AIFilePath != "" {
		_ = os.WriteFile(cfg.AIFilePath, []byte(aiReport), constants.FilePermission)
		fmt.Printf("  \033[1;93m📝 AI Diagnostic saved to:\033[0m %s\n", cfg.AIFilePath)
	}
	if cfg.IsAI {
		fmt.Printf("\n%s\n", aiReport)
	}
}

func buildAIFailureReport(target string, attempt int, res macroRunResult) string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf(`### 🤖 AI Failure Diagnostic Report
- **Command / Macro**: %s
- **Attempt**: #%d
- **Exit Code**: %d
- **Duration**: %v
- **Working Dir**: %s
- **Timestamp**: %s
- **Error Summary**: %s

#### Context & Suggested Next Step
The execution of %q did not complete successfully. Review the exit code and runtime logs above to fix source errors before the next retry iteration.`,
		target, attempt, res.ExitCode, res.Duration.Round(time.Millisecond), cwd, time.Now().UTC().Format(time.RFC3339), res.ErrorMsg, target)
}

func sleepNextDuration(baseDelay time.Duration, strategy string, attempt int) {
	sleepTime := calculateBackoff(baseDelay, strategy, attempt)
	fmt.Printf("  \033[2;37m⏳ Sleeping %v before retry #%d...\033[0m\n", sleepTime, attempt+1)
	time.Sleep(sleepTime)
}

func calculateBackoff(base time.Duration, strategy string, attempt int) time.Duration {
	switch strings.ToLower(strategy) {
	case "linear":
		return base * time.Duration(attempt)
	case "exponential":
		multiplier := math.Pow(2, float64(attempt-1))
		if multiplier > 30 {
			multiplier = 30
		}
		return time.Duration(float64(base) * multiplier)
	default:
		return base
	}
}

func printRetrySuccess(attempt int, totalDur time.Duration) {
	fmt.Printf("\n\033[1;92m🎉 ── Success on attempt #%d ── 🎉\033[0m\n", attempt)
	fmt.Printf("Total execution time: %v\n\n", totalDur.Round(time.Millisecond))
}
