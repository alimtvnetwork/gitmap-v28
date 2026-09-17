package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runGitLogStat(repoDir string, count int) (string, bool) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", count), "--stat", "--no-merges")
	if len(repoDir) > 0 {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()

	return strings.TrimSpace(string(out)), err == nil && len(out) > 0
}

// ExtractGitLog extracts recent commit history from local git repository.
func ExtractGitLog(repoDir string, count int) string {
	if count <= 0 {
		count = 5
	}

	if out, hasStat := runGitLogStat(repoDir, count); hasStat {
		return out
	}

	return fallbackGitLog(repoDir, count)
}

func fallbackGitLog(repoDir string, count int) string {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", count), "--oneline")
	if len(repoDir) > 0 {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	return "(recent git commit log unavailable)"
}

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

// LoadCanonicalRcaPrompt loads RCA bug fix or CI/CD fix prompt template.
func LoadCanonicalRcaPrompt(customPath string, isNoRelease bool) (string, string) {
	if prompt, path, hasCustom := readCustomPrompt(customPath); hasCustom {
		return prompt, path
	}

	for _, candidate := range selectRcaPromptCandidates(isNoRelease) {
		data, err := os.ReadFile(candidate)
		if err == nil && len(data) > 0 {
			return string(data), candidate
		}
	}

	return defaultRcaFixPromptFallback, "embedded-rca-fallback"
}

func selectRcaPromptCandidates(isNoRelease bool) []string {
	rootDir := resolveProjectRootDir()
	if isNoRelease {
		return []string{
			filepath.Join(rootDir, "01-prompts/07-bug-fix/01-fix-with-rca.md"),
			filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md"),
		}
	}

	return []string{
		filepath.Join(rootDir, "01-prompts/07-bug-fix/01-fix-with-rca.md"),
		filepath.Join(rootDir, "01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md"),
		filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md"),
	}
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

func selectPromptPath(isNoRelease bool) string {
	rootDir := resolveProjectRootDir()
	if isNoRelease {
		return filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md")
	}

	return filepath.Join(rootDir, "01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md")
}

func loadCicdFixPrompt(customPath string, isNoRelease bool) (string, string) {
	if prompt, path, hasCustom := readCustomPrompt(customPath); hasCustom {
		return prompt, path
	}

	targetPath := selectPromptPath(isNoRelease)
	data, err := os.ReadFile(targetPath)
	if err == nil {
		return string(data), targetPath
	}

	return defaultCicdFixWithReleasePromptFallback, "embedded-fallback"
}

const defaultCicdFixWithReleasePromptFallback = `# Release-Triggered CI/CD Fix Loop — Workflow (must follow)

N = 200
`

const defaultRcaFixPromptFallback = `# Bug Fix with 4-Part RCA & Regression Verification — Workflow (must follow)

/goal Autonomously fix the failing CI/CD pipeline errors, strictly enforcing coding guidelines, and document the complete RCA before pushing.

## The 4-Part RCA Requirement (Mandatory Memory File)
Before modifying code, document the issue in .lovable/memory/issues/xx-<slug>.md:
1. Why it happened: High-level architectural breakdown of the failure.
2. How it happened: Technical execution flow that triggered the error.
3. Root Cause: Exact file, line, and dependency responsible.
4. Code Fix: Exact code snippets showing the remediation.

## Non-Negotiable Rules
- Functions <= 8-15 lines.
- Affirmative booleans only (is*, has*).
- Zero naked panics or unhandled errors; wrap with *apperror.AppError.
- Group all changes into a single atomic commit at the final step.
`
