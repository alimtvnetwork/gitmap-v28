package cmd

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// checkDuplicateBinaries detects multiple gitmap binaries on PATH.
// When >1 entry exists the uninstaller's Get-Command / which returns an
// array, producing a cryptic "not recognized as a cmdlet" error.
func checkDuplicateBinaries() int {
	paths := findAllBinaries()
	if len(paths) <= 1 {
		printOK(constants.DoctorDupBinOK)

		return 0
	}

	printIssue(constants.DoctorDupBinTitle, formatDupList(paths))
	printFix(formatDupFix(paths))

	return 1
}

// findAllBinaries returns every resolved gitmap binary path on PATH.
func findAllBinaries() []string {
	if runtime.GOOS == "windows" {
		return findAllBinariesWindows()
	}

	return findAllBinariesUnix()
}

// findAllBinariesWindows uses PowerShell Get-Command to list all matches.
func findAllBinariesWindows() []string {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-Command gitmap -All -ErrorAction SilentlyContinue).Source")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parseMultiline(string(out))
}

// findAllBinariesUnix uses `which -a` (or type -a as fallback) to list all matches.
func findAllBinariesUnix() []string {
	cmd := exec.Command("which", "-a", constants.GitMapBin)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parseMultiline(string(out))
}

// parseMultiline splits output into non-empty trimmed lines.
func parseMultiline(output string) []string {
	lines := strings.Split(output, "\n")
	var results []string
	seen := make(map[string]struct{})
	for _, line := range lines {
		p := strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if len(p) == 0 {
			continue
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		lower := strings.ToLower(abs)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		results = append(results, abs)
	}

	return results
}

// formatDupList formats the duplicate binary paths for display.
func formatDupList(paths []string) string {
	var b strings.Builder
	b.WriteString("Multiple gitmap binaries found on PATH:\n")
	for i, p := range paths {
		v := getBinaryVersion(p)
		if i == 0 {
			b.WriteString("       [active] " + p + " (" + v + ")\n")
		} else {
			b.WriteString("       [stale]  " + p + " (" + v + ")\n")
		}
	}

	return b.String()
}

// formatDupFix returns a one-shot removal command for each stale binary.
func formatDupFix(paths []string) string {
	stale := paths[1:]
	if runtime.GOOS == "windows" {
		return formatDupFixWindows(stale)
	}

	return formatDupFixUnix(stale)
}

// formatDupFixWindows returns a PowerShell one-liner to remove stale binaries.
func formatDupFixWindows(stale []string) string {
	if len(stale) == 1 {
		return "Remove-Item '" + stale[0] + "' -Force"
	}
	var b strings.Builder
	for i, p := range stale {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString("Remove-Item '" + p + "' -Force")
	}

	return b.String()
}

// formatDupFixUnix returns a shell command to remove stale binaries.
func formatDupFixUnix(stale []string) string {
	if len(stale) == 1 {
		return "sudo rm '" + stale[0] + "'"
	}
	var b strings.Builder
	b.WriteString("sudo rm")
	for _, p := range stale {
		b.WriteString(" '" + p + "'")
	}

	return b.String()
}
