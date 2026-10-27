package cmdproc

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fspath"
)

// Run routes gitmap proc commands (kill, list).
func Run(args []string) error {
	if len(args) == 0 {
		printUsage()

		return nil
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "kill", "stop", "terminate":
		return RunKill(args[1:])
	case "help", "-h", "--help":
		printUsage()

		return nil
	default:
		if _, err := strconv.Atoi(sub); err == nil {
			return RunKill(args)
		}

		printUsage()

		return apperror.NewValidationError(fmt.Sprintf("unknown proc subcommand: %s", sub))
	}
}

func printUsage() {
	fmt.Println("GitMap Process Management:")
	fmt.Println("  gitmap proc kill <pid> [--force]   Gracefully or forcefully terminate a process by PID")
	fmt.Println("  gitmap proc <pid> [--force]        Shorthand to terminate process by PID")
}

// RunKill terminates a process by PID with safety checks and two-stage escalation.
func RunKill(args []string) error {
	fs := flag.NewFlagSet("proc kill", flag.ContinueOnError)
	isForce := fs.Bool("force", false, "force termination without graceful timeout")
	fs.BoolVar(isForce, "f", false, "force termination (shorthand)")

	if err := fs.Parse(args); err != nil {
		return apperror.WrapValidation(err, "parse proc kill flags")
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		return apperror.NewValidationError("missing PID argument (usage: gitmap proc kill <pid> [--force])")
	}

	pid, err := strconv.Atoi(strings.TrimSpace(remaining[0]))
	if err != nil || pid <= 0 {
		return apperror.NewValidationError(fmt.Sprintf("invalid PID: %q (must be a positive integer)", remaining[0]))
	}

	return TerminateProcess(pid, *isForce)
}

// TerminateProcess terminates target process with self-protection and two-stage kill.
func TerminateProcess(pid int, isForce bool) error {
	myPid := os.Getpid()
	if pid == myPid {
		return apperror.NewValidationError("cannot terminate current GitMap process (self-kill protection)")
	}

	parentPid := os.Getppid()
	if pid == parentPid {
		return apperror.NewValidationError(fmt.Sprintf("refusing to terminate parent terminal process (PID %d)", pid))
	}

	if pid <= 4 && runtime.GOOS == "windows" {
		return apperror.NewValidationError(fmt.Sprintf("refusing to terminate critical OS System process (PID %d)", pid))
	}

	if pid <= 1 && runtime.GOOS != "windows" {
		return apperror.NewValidationError(fmt.Sprintf("refusing to terminate init/system process (PID %d)", pid))
	}

	procName := LookupProcessName(pid)
	if fspath.IsProtectedProcess(procName, pid) && !isForce {
		fmt.Fprintf(os.Stderr, "⚠️ Refused to terminate protected process %q (PID %d).\n", procName, pid)
		fmt.Fprintln(os.Stderr, "💡 This process appears to be an IDE, editor, or shell. Use '--force' to terminate anyway.")

		return apperror.NewWithDetails(
			"proc",
			"E9001",
			fmt.Sprintf("process %q (PID %d) is protected", procName, pid),
			"cli",
			apperror.ErrorTypeAbort,
			apperror.SeverityWarn,
			map[string]any{"reported": true},
		)
	}

	return executeTermination(pid, procName, isForce)
}

func executeTermination(pid int, procName string, isForce bool) error {
	if isForce {
		return forceKill(pid, procName)
	}

	if err := gracefulKill(pid); err != nil {
		return forceKill(pid, procName)
	}

	for i := 0; i < 6; i++ {
		time.Sleep(250 * time.Millisecond)
		if !IsProcessRunning(pid) {
			fmt.Printf("✓ Process %d (%s) terminated gracefully.\n", pid, procName)

			return nil
		}
	}

	fmt.Printf("▸ Process %d did not terminate gracefully; escalating to forceful termination...\n", pid)

	return forceKill(pid, procName)
}

func gracefulKill(pid int) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid))

		return cmd.Run()
	}

	cmd := exec.Command("kill", "-15", strconv.Itoa(pid))

	return cmd.Run()
}

func forceKill(pid int, procName string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	} else {
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		if strings.Contains(strings.ToLower(outStr), "not found") || !IsProcessRunning(pid) {
			fmt.Printf("✓ Process %d is not running or already terminated.\n", pid)

			return nil
		}

		return apperror.WrapSimple(err, fmt.Sprintf("kill process %d: %s", pid, outStr))
	}

	fmt.Printf("✓ Process %d (%s) forcefully terminated.\n", pid, procName)

	return nil
}

// IsProcessRunning checks if the given PID is currently active.
func IsProcessRunning(pid int) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH")
		out, err := cmd.Output()
		if err != nil {
			return false
		}

		return strings.Contains(string(out), strconv.Itoa(pid))
	}

	cmd := exec.Command("kill", "-0", strconv.Itoa(pid))

	return cmd.Run() == nil
}

// LookupProcessName resolves the executable name for a given PID.
func LookupProcessName(pid int) string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
		out, err := cmd.Output()
		if err != nil {
			return "unknown"
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				name := strings.Trim(parts[0], `" `)
				if name != "" && !strings.HasPrefix(name, "INFO:") {
					return name
				}
			}
		}

		return "unknown"
	}

	out, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return "unknown"
}
