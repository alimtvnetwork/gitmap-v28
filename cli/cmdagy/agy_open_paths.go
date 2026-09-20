package cmdagy

import (
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ResolveAntigravityIDE locates the desktop GUI editor executable (Antigravity.exe / Antigravity).
func ResolveAntigravityIDE() result.Result[string] {
	pathRes := findBinaryInPath([]string{"antigravity", "Antigravity"})
	if pathRes.IsSuccess() {
		return pathRes
	}

	return checkCandidatePaths(getCandidateAntigravityIDEPaths())
}

// ResolveAntigravityCLI locates the command-line assistant executable (agy / agy.exe).
func ResolveAntigravityCLI() result.Result[string] {
	pathRes := findBinaryInPath([]string{"agy", "agy.exe", "antigravity-cli"})
	if pathRes.IsSuccess() {
		return pathRes
	}

	return checkCandidatePaths(getCandidateAntigravityCLIPaths())
}

// ResolveAntigravityBinary locates either the Antigravity IDE or CLI executable.
func ResolveAntigravityBinary() result.Result[string] {
	ideRes := ResolveAntigravityIDE()
	if ideRes.IsSuccess() {
		return ideRes
	}

	return ResolveAntigravityCLI()
}

// IsAntigravityIDERunning reports whether an Antigravity IDE process is active.
func IsAntigravityIDERunning() bool {
	procRes := DetectRunningAntigravityIDE()

	return procRes.IsSuccess()
}

// DetectRunningAntigravityIDE inspects running processes via tasklist or ps.
func DetectRunningAntigravityIDE() result.Result[AgyProcessInfo] {
	if runtime.GOOS == "windows" {
		return detectRunningIDEWindows()
	}

	return detectRunningIDEUnix()
}

func detectRunningIDEWindows() result.Result[AgyProcessInfo] {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	out, err := cmd.Output()
	if err != nil {
		return result.Fail[AgyProcessInfo](apperror.WrapSimple(err, "detectRunningIDEWindows"))
	}

	return parseRunningIDEFromTasklist(string(out))
}

func parseRunningIDEFromTasklist(output string) result.Result[AgyProcessInfo] {
	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return result.Fail[AgyProcessInfo](apperror.WrapSimple(err, "parseRunningIDEFromTasklist"))
	}

	return findIDEProcessInRecords(records)
}

func findIDEProcessInRecords(records [][]string) result.Result[AgyProcessInfo] {
	for _, rec := range records {
		matchRes := extractIDEProcessRecord(rec)
		if matchRes.IsSuccess() {
			return matchRes
		}
	}

	return result.Fail[AgyProcessInfo](apperror.NewSimple("Antigravity IDE process is not running", "E9002"))
}

func extractIDEProcessRecord(rec []string) result.Result[AgyProcessInfo] {
	if len(rec) < 2 {
		return result.Fail[AgyProcessInfo](apperror.NewSimple("invalid record", "E9003"))
	}

	name := strings.TrimSpace(rec[0])
	pid, convErr := strconv.Atoi(strings.TrimSpace(rec[1]))
	hasValidPid := convErr == nil && pid > 0 && strings.EqualFold(filepath.Base(name), "antigravity.exe")

	if hasValidPid {
		return result.Ok(AgyProcessInfo{PID: pid, Name: name})
	}

	return result.Fail[AgyProcessInfo](apperror.NewSimple("not antigravity IDE", "E9004"))
}

func detectRunningIDEUnix() result.Result[AgyProcessInfo] {
	cmd := exec.Command("ps", "-eo", "pid,comm")
	out, err := cmd.Output()
	if err != nil {
		return result.Fail[AgyProcessInfo](apperror.WrapSimple(err, "detectRunningIDEUnix"))
	}

	return parseRunningIDEFromPs(string(out))
}

func parseRunningIDEFromPs(output string) result.Result[AgyProcessInfo] {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		matchRes := extractIDEProcessLine(line)
		if matchRes.IsSuccess() {
			return matchRes
		}
	}

	return result.Fail[AgyProcessInfo](apperror.NewSimple("Antigravity IDE process is not running", "E9002"))
}

func extractIDEProcessLine(line string) result.Result[AgyProcessInfo] {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return result.Fail[AgyProcessInfo](apperror.NewSimple("invalid line", "E9003"))
	}

	return buildUnixIDEProcessInfo(fields[0], fields[1])
}

func buildUnixIDEProcessInfo(pidStr, rawName string) result.Result[AgyProcessInfo] {
	pid, convErr := strconv.Atoi(pidStr)
	name := filepath.Base(rawName)
	hasValidPid := convErr == nil && pid > 0 && strings.EqualFold(name, "antigravity")
	if hasValidPid {
		return result.Ok(AgyProcessInfo{PID: pid, Name: name})
	}

	return result.Fail[AgyProcessInfo](apperror.NewSimple("not antigravity IDE", "E9004"))
}

func isValidBinaryPath(path string) bool {
	if len(path) == 0 {
		return false
	}

	_, err := os.Stat(path)

	return err == nil
}

func checkCandidatePaths(candidates []string) result.Result[string] {
	for _, candidate := range candidates {
		if isValidBinaryPath(candidate) {
			return result.Ok(candidate)
		}
	}

	return result.Fail[string](apperror.NewSimple("binary not found in candidate paths", "E9006"))
}

func findBinaryInPath(names []string) result.Result[string] {
	for _, name := range names {
		path, err := exec.LookPath(name)
		hasPath := err == nil && len(path) > 0

		if hasPath {
			return result.Ok(path)
		}
	}

	return result.Fail[string](apperror.NewSimple("binary not found in PATH", "E9005"))
}

func getCandidateAntigravityIDEPaths() []string {
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	progFiles := os.Getenv("ProgramFiles")

	windowsPaths := getWindowsIDEPaths(home, localApp, progFiles)
	unixPaths := getUnixIDEPaths(home)

	return append(windowsPaths, unixPaths...)
}

func getWindowsIDEPaths(home, localApp, progFiles string) []string {
	progFilesX86 := os.Getenv("ProgramFiles(x86)")

	return []string{
		filepath.Join(localApp, "Programs", "Antigravity", "Antigravity.exe"),
		filepath.Join(localApp, "Programs", "antigravity", "Antigravity.exe"),
		filepath.Join(progFiles, "Antigravity", "Antigravity.exe"),
		filepath.Join(progFilesX86, "Antigravity", "Antigravity.exe"),
		filepath.Join(home, "AppData", "Local", "Programs", "antigravity", "Antigravity.exe"),
		filepath.Join(home, "AppData", "Local", "Programs", "Antigravity", "Antigravity.exe"),
		"C:\\Program Files\\Antigravity\\Antigravity.exe",
	}
}

func getUnixIDEPaths(home string) []string {
	paths := []string{
		filepath.Join(home, ".local", "share", "antigravity", "Antigravity"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		"/opt/antigravity/Antigravity", "/opt/Antigravity/antigravity",
		"/usr/share/antigravity/Antigravity", "/usr/local/bin/antigravity",
		"/usr/bin/antigravity", "/snap/bin/antigravity",
		"/var/lib/flatpak/exports/bin/antigravity",
		"/Applications/Antigravity.app/Contents/MacOS/Antigravity",
		filepath.Join(home, "Applications", "Antigravity.app", "Contents", "MacOS", "Antigravity"),
	}

	return paths
}

func getCandidateAntigravityCLIPaths() []string {
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")

	return []string{
		filepath.Join(localApp, "agy", "bin", "agy.exe"),
		filepath.Join(appData, "npm", "agy.cmd"),
		filepath.Join(home, ".local", "bin", "agy"),
		filepath.Join(home, ".gemini", "antigravity", "bin", "agy"),
		"/usr/local/bin/agy",
		"/usr/bin/agy",
	}
}
