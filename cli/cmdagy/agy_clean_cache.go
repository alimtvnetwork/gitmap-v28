// Package cmdagy — agy_clean_cache.go provides core cache cleaning and process management logic.
package cmdagy

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// AgyCacheTarget represents a cache directory target to be cleaned.
type AgyCacheTarget struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Exists      bool   `json:"exists"`
	SizeBytes   int64  `json:"sizeBytes"`
	FileCount   int    `json:"fileCount"`
}

// AgyProcessInfo represents a discovered Antigravity/browser process.
type AgyProcessInfo struct {
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

// CleanCacheOptions holds CLI flags and execution parameters.
type CleanCacheOptions struct {
	DryRun      bool `json:"dryRun"`
	Force       bool `json:"force"`
	Yes         bool `json:"yes"`
	NoKill      bool `json:"noKill"`
	JSON        bool `json:"json"`
	IncludeTemp bool `json:"includeTemp"`
}

// CleanCacheReport summarizes the results of the clean-cache execution.
type CleanCacheReport struct {
	Timestamp           string           `json:"timestamp"`
	DryRun              bool             `json:"dryRun"`
	Targets             []AgyCacheTarget `json:"targets"`
	Processes           []AgyProcessInfo `json:"processes"`
	ProcessesTerminated int              `json:"processesTerminated"`
	BytesFreed          int64            `json:"bytesFreed"`
	HumanFreed          string           `json:"humanFreed"`
	FilesDeleted        int              `json:"filesDeleted"`
	Warnings            []string         `json:"warnings,omitempty"`
	DurationMs          int64            `json:"durationMs"`
}

// CalculateDirStats recursively walks a directory to compute total size and file count.
func CalculateDirStats(dirPath string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	walkErr := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue scanning accessible files
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	return totalSize, fileCount, walkErr
}

// CleanDirectoryContents removes all items inside dirPath while preserving the parent directory.
func CleanDirectoryContents(dirPath string) (int64, int, []string) {
	var freedBytes int64
	var deletedFiles int
	var warnings []string

	entries, readErr := os.ReadDir(dirPath)
	if readErr != nil {
		warnings = append(warnings, fmt.Sprintf("read dir %s: %v", dirPath, readErr))
		return freedBytes, deletedFiles, warnings
	}

	for _, entry := range entries {
		subPath := filepath.Join(dirPath, entry.Name())
		subSize, subFiles, _ := CalculateDirStats(subPath)

		removeErr := os.RemoveAll(subPath)
		if removeErr != nil {
			warnings = append(warnings, fmt.Sprintf("failed removing %s: %v", subPath, removeErr))
			continue
		}

		freedBytes += subSize
		if subFiles > 0 {
			deletedFiles += subFiles
		} else {
			deletedFiles++
		}
	}

	return freedBytes, deletedFiles, warnings
}

// FormatBytes converts a byte count to a human-readable string (KB, MB, GB).
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// RenderCleanPreview prints target directories and detected processes to standard output.
func RenderCleanPreview(targets []AgyCacheTarget, procs []AgyProcessInfo) {
	fmt.Printf("\n%s⚡ Antigravity Cache & Process Cleanup%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%sTarget cleanup locations:%s\n", constants.ColorWhite, constants.ColorReset)

	for _, t := range targets {
		if t.Exists {
			fmt.Printf("  %s✓%s %-48s %s%s, %d files%s\n",
				constants.ColorGreen, constants.ColorReset,
				t.Path,
				constants.ColorDim, FormatBytes(t.SizeBytes), t.FileCount, constants.ColorReset)
		} else {
			fmt.Printf("  %s•%s %-48s %s(not found)%s\n",
				constants.ColorDim, constants.ColorReset,
				t.Path,
				constants.ColorDim, constants.ColorReset)
		}
	}

	fmt.Printf("\n%sProcesses to terminate:%s\n", constants.ColorWhite, constants.ColorReset)
	if len(procs) == 0 {
		fmt.Printf("  %s• None running%s\n", constants.ColorDim, constants.ColorReset)
		return
	}

	for _, p := range procs {
		fmt.Printf("  %s• PID %-6d %s%s\n", constants.ColorYellow, p.PID, p.Name, constants.ColorReset)
	}
}

// AskProceedConfirmation prompts the user interactively to confirm cleanup.
func AskProceedConfirmation() bool {
	fmt.Printf("\n%sProceed with cleanup? [y/N]: %s", constants.ColorYellow, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes"
}

// RenderCleanSuccess prints execution metrics and safety confirmation.
func RenderCleanSuccess(report CleanCacheReport) {
	fmt.Printf("\n%s✓ Cleanup completed successfully!%s\n", constants.ColorGreen, constants.ColorReset)
	if report.ProcessesTerminated > 0 {
		fmt.Printf("  • Terminated %d process(es)\n", report.ProcessesTerminated)
	}
	fmt.Printf("  • Removed %d file(s) (%s freed)\n", report.FilesDeleted, report.HumanFreed)
	fmt.Printf("  • Completed in %dms\n", report.DurationMs)

	if len(report.Warnings) > 0 {
		fmt.Printf("\n%sWarnings (locked files skipped):%s\n", constants.ColorYellow, constants.ColorReset)
		for _, w := range report.Warnings {
			fmt.Printf("  %s• %s%s\n", constants.ColorDim, w, constants.ColorReset)
		}
	}

	fmt.Printf("\n%sProtected: Projects, conversations, transcripts, and settings were preserved.%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

// RenderCleanJSON marshals and prints the CleanCacheReport in indented JSON format.
func RenderCleanJSON(report CleanCacheReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

// BuildBaseCacheReport constructs an initial report before execution.
func BuildBaseCacheReport(opts CleanCacheOptions, targets []AgyCacheTarget, procs []AgyProcessInfo) CleanCacheReport {
	return CleanCacheReport{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		DryRun:    opts.DryRun,
		Targets:   targets,
		Processes: procs,
	}
}

var targetProcRegex = regexp.MustCompile(`(?i)^(antigravity|electron|msedge|msedgewebview2)(\.exe)?$`)

func parseTasklistCSV(output string, currentPid int) []AgyProcessInfo {
	var procs []AgyProcessInfo
	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	if err != nil {
		return procs
	}

	seen := make(map[int]bool)
	for _, rec := range records {
		if len(rec) < 2 {
			continue
		}

		name := strings.TrimSpace(rec[0])
		pid, convErr := strconv.Atoi(strings.TrimSpace(rec[1]))
		if convErr != nil || pid == currentPid || seen[pid] {
			continue
		}

		if targetProcRegex.MatchString(name) {
			seen[pid] = true
			procs = append(procs, AgyProcessInfo{PID: pid, Name: name})
		}
	}

	return procs
}

func parsePsOutput(output string, currentPid int) []AgyProcessInfo {
	var procs []AgyProcessInfo
	lines := strings.Split(output, "\n")
	seen := make(map[int]bool)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		pid, convErr := strconv.Atoi(fields[0])
		if convErr != nil || pid == currentPid || seen[pid] {
			continue
		}

		name := filepath.Base(fields[1])
		if targetProcRegex.MatchString(name) {
			seen[pid] = true
			procs = append(procs, AgyProcessInfo{PID: pid, Name: name})
		}
	}

	return procs
}
