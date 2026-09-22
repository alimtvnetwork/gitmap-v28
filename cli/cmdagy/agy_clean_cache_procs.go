// Package cmdagy — agy_clean_cache_procs.go parses running processes for cleanup.
package cmdagy

import (
	"encoding/csv"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

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
