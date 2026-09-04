// Package cmd — pull_table_format.go provides formatting and middle-truncation for table cells.
package cmd

import (
	"strings"
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
	if len(input) <= maxLength {
		return input
	}
	minRequired := 3 + endLength + 1
	if maxLength < minRequired {
		return input[:maxLength]
	}
	startLength := maxLength - 3 - endLength
	startPart := input[:startLength]
	endPart := input[len(input)-endLength:]
	truncated := startPart + "..." + endPart

	return truncated
}

func formatBranchName(branch string, maxLength int) string {
	if len(branch) <= 0 {
		return ""
	}
	cleanedBranch := stripBranchPrefix(branch)
	formattedBranch := middleTruncate(cleanedBranch, maxLength, 5)

	return formattedBranch
}

func formatRepoName(repo string, maxLength int) string {
	if len(repo) <= 0 {
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
