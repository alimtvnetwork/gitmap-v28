// Package cmd — pull_table_format.go provides formatting and middle-truncation for table cells.
package cmdpull

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
)

func stripBranchPrefix(branch string) string {
	candidatePrefixes := []string{
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
	if len(repo) == 0 {
		return ""
	}

	formattedRepo := middleTruncate(repo, maxLength, 5)

	return formattedRepo
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
