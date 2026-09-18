package cmdagy

import (
	"fmt"
	"strings"
)

func appendMetaLine(b *strings.Builder, label, val string) {
	if len(val) > 0 {
		b.WriteString(fmt.Sprintf("> **%s:** `%s`\n", label, val))
	}
}

// FormatRCAHeader formats the top-level RCA directive header.
func FormatRCAHeader(repo string, runID uint64, sha string) string {
	var b strings.Builder
	b.WriteString("# Fix CI/CD Pipeline Errors with 4-Part Root Cause Analysis (RCA)\n\n")
	appendMetaLine(&b, "Target Repository", repo)
	if runID > 0 {
		appendMetaLine(&b, "Failing Workflow Run", fmt.Sprintf("#%d", runID))
	}
	appendMetaLine(&b, "Commit SHA", sha)
	b.WriteString("\n/goal Autonomously analyze and fix the failing CI/CD pipeline errors documented below using 4-part Root Cause Analysis (RCA), strictly adhere to coding guidelines, and verify all quality gates pass.\n")

	return b.String()
}

// FormatGitLogSection formats recent commit logs into markdown.
func FormatGitLogSection(gitLog string) string {
	trimmed := strings.TrimSpace(gitLog)
	if len(trimmed) == 0 {
		return ""
	}

	return "## Recent Git Commit History\n\n```text\n" + trimmed + "\n```"
}

// FormatPipelineErrorSection formats pipeline failure diagnostics into markdown.
func FormatPipelineErrorSection(errorReport string) string {
	trimmed := strings.TrimSpace(errorReport)
	if len(trimmed) == 0 {
		return ""
	}

	return "## Failing CI/CD Pipeline Error Logs\n\n```text\n" + trimmed + "\n```"
}

func appendOptionalSection(sections *[]string, content string) {
	if clean := strings.TrimSpace(content); len(clean) > 0 {
		*sections = append(*sections, clean)
	}
}

// AssembleRcaFixPayload combines header, git log, error logs, and prompt template.
func AssembleRcaFixPayload(repo string, runID uint64, sha, gitLog, errorLogs, fixPrompt string) string {
	sections := make([]string, 0, 4)
	sections = append(sections, FormatRCAHeader(repo, runID, sha))
	appendOptionalSection(&sections, FormatGitLogSection(gitLog))
	appendOptionalSection(&sections, FormatPipelineErrorSection(errorLogs))
	appendOptionalSection(&sections, fixPrompt)

	return strings.Join(sections, "\n\n")
}

// AssembleFixPipelinePayload combines error logs and fix prompt with mandatory two-line gap.
func AssembleFixPipelinePayload(errorLogs, fixPrompt string) string {
	cleanLogs := strings.TrimSpace(errorLogs)
	cleanPrompt := strings.TrimSpace(fixPrompt)
	if len(cleanLogs) == 0 {
		return cleanPrompt
	}
	if len(cleanPrompt) == 0 {
		return cleanLogs
	}

	return cleanLogs + "\n\n" + cleanPrompt
}
