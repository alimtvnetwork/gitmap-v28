package cmdautomation

import (
	"os"
	"path/filepath"
	"strings"
)

// parseSearchPositionalArgs parses pattern and optional target dir from positional arguments.
func parseSearchPositionalArgs(args []string) (string, string) {
	if len(args) == 0 {
		return "", ""
	}
	if len(args) == 1 {
		return cleanSearchPattern(args[0]), ""
	}
	if isPowerShellEscapedSplit(args) {
		return reconstructPowerShellSplit(args)
	}
	lastIdx := len(args) - 1
	if isTargetDirectoryOrFile(args[lastIdx]) {
		pattern := strings.Join(args[:lastIdx], " ")
		return cleanSearchPattern(pattern), args[lastIdx]
	}
	return cleanSearchPattern(strings.Join(args, " ")), ""
}

func isPowerShellEscapedSplit(args []string) bool {
	if len(args) < 2 {
		return false
	}
	isLeadingSlash := args[0] == `\` || args[0] == `\"`
	if !isLeadingSlash {
		return false
	}
	for i := 1; i < len(args); i++ {
		if strings.HasSuffix(args[i], `\`) || strings.HasSuffix(args[i], `\"`) {
			return true
		}
	}
	return false
}

func reconstructPowerShellSplit(args []string) (string, string) {
	endIdx := findSplitEndIndex(args)
	if endIdx == -1 {
		return cleanSearchPattern(args[0]), ""
	}
	patternParts := collectSplitPatternParts(args, endIdx)
	pattern := cleanSearchPattern(strings.Join(patternParts, " "))
	targetDir := ""
	if len(args) > endIdx+1 {
		targetDir = strings.Join(args[endIdx+1:], " ")
	}
	return pattern, targetDir
}

func findSplitEndIndex(args []string) int {
	for i := 1; i < len(args); i++ {
		if strings.HasSuffix(args[i], `\`) || strings.HasSuffix(args[i], `\"`) {
			return i
		}
	}
	return -1
}

func collectSplitPatternParts(args []string, endIdx int) []string {
	var parts []string
	for i := 1; i <= endIdx; i++ {
		item := args[i]
		if i == endIdx {
			item = strings.TrimSuffix(strings.TrimSuffix(item, `\"`), `\`)
		}
		if len(item) > 0 {
			parts = append(parts, item)
		}
	}
	return parts
}

func isTargetDirectoryOrFile(target string) bool {
	clean := strings.TrimSpace(target)
	if len(clean) == 0 {
		return false
	}
	if _, err := os.Stat(clean); err == nil {
		return true
	}
	if _, isMatch := resolveSearchCandidateFiles(clean); isMatch {
		return true
	}
	hasSlash := strings.Contains(clean, "/") || strings.Contains(clean, "\\")
	return hasSlash || hasCommonFileExtension(clean)
}

func hasCommonFileExtension(target string) bool {
	ext := strings.ToLower(filepath.Ext(target))
	if len(ext) == 0 {
		return false
	}
	for _, candidate := range commonFileExtensions {
		if ext == candidate {
			return true
		}
	}
	return false
}
