package cmdtoolchain

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// KillOptions captures execution options for gitmap kill.
type KillOptions struct {
	Target   string
	HasForce bool
	HasTree  bool
}

// RunKill is the main entry point for gitmap kill / terminate.
func RunKill(args []string) error {
	checkToolchainHelp("kill", args)
	opts, err := parseKillArgs(args)
	if err != nil {
		return err
	}

	return executeKill(opts)
}

func parseKillArgs(args []string) (KillOptions, error) {
	opts, pos := extractKillFlags(args)
	if len(pos) == 0 {
		return opts, apperror.NewValidationError("gitmap kill requires a target (<pid>, <name>, or :<port>)")
	}

	opts.Target = pos[0]

	return opts, nil
}

func extractKillFlags(args []string) (KillOptions, []string) {
	var opts KillOptions
	var pos []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isFlagProcessed(&opts, arg) {
			continue
		}

		if isPortFlag(arg) && i+1 < len(args) {
			i++
			pos = append(pos, ":"+args[i])
			continue
		}

		pos = append(pos, arg)
	}

	return opts, pos
}

var (
	killForceFlagTokens = []string{"-f", "--force"}
	killTreeFlagTokens  = []string{"-t", "--tree"}
	killPortFlagTokens  = []string{"-p", "--port"}
)

func isFlagProcessed(opts *KillOptions, arg string) bool {
	if slices.Contains(killForceFlagTokens, arg) {
		opts.HasForce = true

		return true
	}

	if slices.Contains(killTreeFlagTokens, arg) {
		opts.HasTree = true

		return true
	}

	return false
}

func isPortFlag(arg string) bool {
	return slices.Contains(killPortFlagTokens, arg)
}

func executeKill(opts KillOptions) error {
	if isPortTarget(opts.Target) {
		return killProcessByPort(opts)
	}

	pid, err := strconv.Atoi(opts.Target)
	if err == nil {
		return killProcessByPID(pid, opts)
	}

	return killProcessByName(opts.Target, opts)
}

func isPortTarget(target string) bool {
	return strings.HasPrefix(target, ":")
}

func isProtectedPID(pid int) bool {
	if pid <= 0 || pid == 1 || pid == 4 {
		return true
	}

	return pid == os.Getpid()
}

func isProtectedProcessName(name string) bool {
	lower := strings.ToLower(name)
	protected := []string{
		"antigravity", "antigravity.exe",
		"code", "code.exe",
		"pwsh", "pwsh.exe",
		"bash", "bash.exe",
		"gitmap", "gitmap.exe",
	}
	for _, p := range protected {
		if lower == p || strings.HasPrefix(lower, p) {
			return true
		}
	}

	return false
}

func killProcessByPID(pid int, opts KillOptions) error {
	if isProtectedPID(pid) {
		return apperror.NewValidationError(fmt.Sprintf("refusing to terminate protected system/self PID %d", pid))
	}

	return terminateProcessPID(pid, opts.HasTree)
}

func terminateProcessPID(pid int, hasTree bool) error {
	if runtime.GOOS == "windows" {
		return runWindowsTaskkillPID(pid, hasTree)
	}

	return runPosixKillPID(pid)
}

func runWindowsTaskkillPID(pid int, hasTree bool) error {
	args := []string{"/F"}
	if hasTree {
		args = append(args, "/T")
	}

	args = append(args, "/PID", strconv.Itoa(pid))
	cmd := exec.Command("taskkill", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("taskkill PID %d failed: %s (%v)", pid, string(out), err))
	}

	fmt.Printf("Process PID %d terminated successfully.\n", pid)

	return nil
}

func runPosixKillPID(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("kill -9 %d failed: %s (%v)", pid, string(out), err))
	}

	fmt.Printf("Process PID %d terminated successfully.\n", pid)

	return nil
}

func killProcessByName(name string, opts KillOptions) error {
	if !opts.HasForce && isProtectedProcessName(name) {
		return apperror.NewValidationError(fmt.Sprintf("refusing to terminate protected process %q without --force", name))
	}

	if runtime.GOOS == "windows" {
		return runWindowsTaskkillName(name, opts.HasTree)
	}

	return runPosixKillName(name)
}

func runWindowsTaskkillName(name string, hasTree bool) error {
	args := []string{"/F"}
	if hasTree {
		args = append(args, "/T")
	}

	args = append(args, "/IM", name)
	cmd := exec.Command("taskkill", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("taskkill %q failed: %s (%v)", name, string(out), err))
	}

	fmt.Printf("Process %q terminated successfully.\n", name)

	return nil
}

func runPosixKillName(name string) error {
	cmd := exec.Command("pkill", "-9", "-f", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("pkill -9 %q failed: %s (%v)", name, string(out), err))
	}

	fmt.Printf("Process %q terminated successfully.\n", name)

	return nil
}

func killProcessByPort(opts KillOptions) error {
	portStr := strings.TrimPrefix(opts.Target, ":")
	pid, err := FindListeningPIDByPort(portStr)
	if err != nil {
		return err
	}

	return killProcessByPID(pid, opts)
}

// FindListeningPIDByPort discovers the process ID listening on the given port.
func FindListeningPIDByPort(portStr string) (int, error) {
	if runtime.GOOS == "windows" {
		return findWindowsPortPID(portStr)
	}

	return findPosixPortPID(portStr)
}

func findWindowsPortPID(portStr string) (int, error) {
	out, err := exec.Command("netstat", "-ano", "-p", "tcp").Output()
	if err != nil {
		return 0, apperror.NewExecutionError(fmt.Sprintf("netstat failed: %v", err))
	}

	pid, isFound := parseNetstatOutput(string(out), portStr)
	if !isFound {
		return 0, apperror.NewNotFoundError(fmt.Sprintf("no listening process found on port %s", portStr))
	}

	return pid, nil
}

func parseNetstatOutput(output, portStr string) (int, bool) {
	targetSuffix := ":" + portStr
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if pid, isMatch := matchNetstatLine(line, targetSuffix); isMatch {
			return pid, true
		}
	}

	return 0, false
}

func matchNetstatLine(line, targetSuffix string) (int, bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, false
	}

	if fields[0] != "TCP" || fields[3] != "LISTENING" {
		return 0, false
	}

	if !strings.HasSuffix(fields[1], targetSuffix) {
		return 0, false
	}

	pid, err := strconv.Atoi(fields[4])

	return pid, err == nil
}

func findPosixPortPID(portStr string) (int, error) {
	out, err := exec.Command("lsof", "-ti", ":"+portStr).Output()
	if err != nil {
		return 0, apperror.NewNotFoundError(fmt.Sprintf("no listening process found on port %s", portStr))
	}

	firstLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	pid, err := strconv.Atoi(strings.TrimSpace(firstLine))
	if err != nil {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid PID returned for port %s: %s", portStr, firstLine))
	}

	return pid, nil
}
