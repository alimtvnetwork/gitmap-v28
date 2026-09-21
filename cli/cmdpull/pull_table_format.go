// Package cmd — pull_table_format.go provides formatting and middle-truncation for table cells.
package cmdpull

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
)

func stripBranchPrefix(branch string) string {
	candidatePrefixes := []string{
		"backup/",
		"feature/",
		"feat/",
		"release/",
		"bugfix/",
		"hotfix/",
		"fix/",
		"dependabot/",
	}

	for _, prefix := range candidatePrefixes {
		hasMatch := strings.HasPrefix(strings.ToLower(branch), prefix)
		if hasMatch {
			cleanBranch := branch[len(prefix):]

			return cleanBranch
		}
	}

	return branch
}

// leadTruncate truncates input from the front, retaining ending text with leading "...".
func leadTruncate(input string, maxLength int) string {
	runes := []rune(input)
	if len(runes) <= maxLength {
		return input
	}
	if maxLength <= 0 {
		return ""
	}
	if maxLength <= 3 {
		return string(runes[len(runes)-maxLength:])
	}

	keepLen := maxLength - 3
	return "..." + string(runes[len(runes)-keepLen:])
}

// formatLatestBranchName strips prefixes and preserves the ending text if truncated.
func formatLatestBranchName(branch string, maxLength int) string {
	if len(branch) == 0 {
		return ""
	}

	cleanedBranch := stripBranchPrefix(branch)
	return leadTruncate(cleanedBranch, maxLength)
}

// formatPRCell formats PR count as a 2-digit number (e.g. "01", "02", "30") or "-" if 0/untracked.
func formatPRCell(prStatus string) string {
	prStatus = strings.TrimSpace(prStatus)
	if prStatus == "" || prStatus == "—" || prStatus == "-" {
		return "-"
	}

	numStr := extractLeadingDigits(prStatus)
	if numStr == "" {
		return "-"
	}

	count, err := strconv.Atoi(numStr)
	if err != nil || count <= 0 {
		return "-"
	}

	return fmt.Sprintf("%02d", count)
}

func extractLeadingDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
			continue
		}
		if b.Len() > 0 {
			break
		}
	}

	return b.String()
}

func middleTruncate(input string, maxLength int, endLength int) string {
	runes := []rune(input)
	if len(runes) <= maxLength {
		return input
	}

	minRequired := 3 + endLength + 1
	if maxLength < minRequired {
		return string(runes[:maxLength])
	}

	startLength := maxLength - 3 - endLength
	startPart := string(runes[:startLength])
	endPart := string(runes[len(runes)-endLength:])

	return startPart + "..." + endPart
}

func formatBranchName(branch string, maxLength int) string {
	if len(branch) == 0 {
		return ""
	}

	cleanedBranch := stripBranchPrefix(branch)
	formattedBranch := middleTruncate(cleanedBranch, maxLength, 5)

	return formattedBranch
}

func formatCombinedBranch(branch, latest string, maxLength int) string {
	if len(branch) == 0 {
		return ""
	}

	cleanedBranch := stripBranchPrefix(branch)
	cleanedLatest := stripBranchPrefix(latest)

	isSame := cleanedLatest == "" || strings.EqualFold(cleanedBranch, cleanedLatest)
	if isSame {
		return middleTruncate(cleanedBranch, maxLength, 4)
	}

	arrow := resolveBranchArrow()
	combined := cleanedBranch + arrow + cleanedLatest

	return middleTruncate(combined, maxLength, 4)
}

func resolveBranchArrow() string {
	if glyphs.Resolve() == glyphs.ModeSafe {
		return "->"
	}

	return "→"
}

func formatRepoName(repo string, maxLength int) string {
	isZeroLen := len(repo) == 0
	if isZeroLen {
		return ""
	}

	return leadTruncate(repo, maxLength)
}

func calcAnsiPadding(renderedText string, visibleWidth int) int {
	plainText := stripANSI(renderedText)
	extraAnsiBytes := len(renderedText) - len(plainText)
	targetPadding := visibleWidth + extraAnsiBytes

	return targetPadding
}

// PadVisual returns string s padded with trailing spaces until visual terminal width is targetWidth.
func PadVisual(s string, targetWidth int) string {
	filtered := glyphs.FilterString(s)
	plain := stripANSI(filtered)
	visWidth := runewidth.StringWidth(plain)
	if visWidth == targetWidth {
		return filtered
	}
	if visWidth > targetWidth {
		return truncateVisual(filtered, targetWidth)
	}

	return filtered + strings.Repeat(" ", targetWidth-visWidth)
}

func truncateVisual(s string, targetWidth int) string {
	if targetWidth <= 0 {
		return ""
	}
	plain := stripANSI(s)
	runes := []rune(plain)
	var curWidth int
	var cutIdx int
	for i, r := range runes {
		rw := runewidth.RuneWidth(r)
		if curWidth+rw > targetWidth {
			break
		}
		curWidth += rw
		cutIdx = i + 1
	}

	return string(runes[:cutIdx]) + strings.Repeat(" ", targetWidth-curWidth)
}

func formatReleaseCell(rel string, maxLength int) string {
	cleaned := strings.TrimSpace(rel)
	isEmpty := len(cleaned) == 0 || cleaned == "—" || cleaned == "-"
	if isEmpty {
		return "-"
	}

	return leadTruncate(cleaned, maxLength)
}

func formatSHACell(sha string, maxLength int) string {
	cleaned := strings.TrimSpace(sha)
	isEmpty := len(cleaned) == 0 || cleaned == "—" || cleaned == "-"
	if isEmpty {
		return "-"
	}

	limit := resolveSHACellLimit(maxLength)
	hasOverflow := len(cleaned) > limit
	if hasOverflow {
		return cleaned[:limit]
	}

	return cleaned
}

func resolveSHACellLimit(maxLength int) int {
	target := 7
	hasSmallerMax := maxLength < target && maxLength > 0
	if hasSmallerMax {
		return maxLength
	}

	return target
}
