// Package cmdcg provides diagnostic health checks and repair suggestions.
package cmdcg

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// CGDoctorDiagnosis encapsulates health issues and suggested fixes for a repository.
type CGDoctorDiagnosis struct {
	RepoPath    string
	Issues      []string
	Suggestions []string
	IsHealthy   bool
}

var cgCheckSuggestions = map[string][2]string{
	"version.json valid":             {"version.json missing or invalid", "Run `gitmap cg install-version-json`"},
	"coding-guidelines section":      {"coding-guidelines section missing", "Run `gitmap cg install`"},
	"02-coding-guidelines directory": {"02-spec/02-coding-guidelines directory missing", "Run `gitmap cg update`"},
	"memory guidelines":              {"Memory guidelines (.ai-memory/) missing", "Run `gitmap cg sync`"},
	".agents/rules directory":        {".agents/rules directory missing", "Run `gitmap cg install-prompts`"},
}

func mapCheckToSuggestion(name string) (string, string) {
	pair, hasSuggestion := cgCheckSuggestions[name]
	if !hasSuggestion {
		return "Configuration check failed", "Run `gitmap cg sync`"
	}

	return pair[0], pair[1]
}

func appendAuditIssue(check CGAuditCheck, issues, suggestions []string) ([]string, []string) {
	if check.IsPass {
		return issues, suggestions
	}

	issue, sugg := mapCheckToSuggestion(check.Name)

	return append(issues, issue), append(suggestions, sugg)
}

func checkDoctorAuditIssues(summary CGAuditSummary) ([]string, []string) {
	var issues []string
	var suggestions []string

	for _, check := range summary.Checks {
		issues, suggestions = appendAuditIssue(check, issues, suggestions)
	}

	return issues, suggestions
}

func checkDoctorDrift(repoPath string) (string, string, bool) {
	meta, err := ReadCGMetadata(repoPath)
	if err != nil || meta == nil || meta.Version == "" {
		return "", "", false
	}

	if meta.Version == CanonicalCGVersion {
		return "", "", false
	}

	issue := fmt.Sprintf("Version drift detected: %s -> canonical %s", meta.Version, CanonicalCGVersion)
	sugg := "Run `gitmap cg update` to upgrade guidelines"

	return issue, sugg, true
}

func appendDriftIssue(repoPath string, issues, suggestions []string) ([]string, []string) {
	driftIssue, driftSugg, hasDrift := checkDoctorDrift(repoPath)
	if !hasDrift {
		return issues, suggestions
	}

	return append(issues, driftIssue), append(suggestions, driftSugg)
}

func diagnoseSingleRepo(repoPath string) CGDoctorDiagnosis {
	summary := auditSingleRepo(repoPath)
	issues, suggestions := checkDoctorAuditIssues(summary)
	issues, suggestions = appendDriftIssue(repoPath, issues, suggestions)

	return CGDoctorDiagnosis{
		RepoPath:    repoPath,
		Issues:      issues,
		Suggestions: suggestions,
		IsHealthy:   len(issues) == 0,
	}
}

// DiagnoseSingleRepo performs doctor diagnostics on a single repository.
func DiagnoseSingleRepo(repoPath string) CGDoctorDiagnosis {
	return diagnoseSingleRepo(repoPath)
}

func renderDoctorDiagnosis(diag CGDoctorDiagnosis) {
	name := filepath.Base(diag.RepoPath)
	if diag.IsHealthy {
		fmt.Printf("  %s %s: Healthy (all checks passed)\n", termout.TokyoOKPill("HEALTHY"), cgAuditTitleStyle.Render(name))

		return
	}

	fmt.Printf("  %s %s: %d issue(s) detected\n", termout.TokyoWarnPill("ACTION"), cgAuditTitleStyle.Render(name), len(diag.Issues))
	for i, issue := range diag.Issues {
		fmt.Printf("    • %s %s\n", cgWarnStyle.Render("Issue:"), issue)
		fmt.Printf("      %s %s\n", termout.TokyoOKPill("FIX"), diag.Suggestions[i])
	}
}

func buildDoctorSummaryContent(totalRepos, healthyRepos int) string {
	issuesCount := totalRepos - healthyRepos
	bar := termout.RenderTokyoProgress(healthyRepos, totalRepos, 24)
	header := fmt.Sprintf("Doctor Summary | Health: %s (%d/%d healthy)", bar, healthyRepos, totalRepos)
	if issuesCount == 0 {
		return fmt.Sprintf("%s\nAll %d repositories are in optimal state.", header, totalRepos)
	}

	return fmt.Sprintf("%s\nAction required on %d repository(ies).", header, issuesCount)
}

func renderDoctorOverall(totalRepos, healthyRepos int) {
	fmt.Println()
	content := buildDoctorSummaryContent(totalRepos, healthyRepos)
	fmt.Println(termout.TokyoBoxStyle(content))
}

func runDoctorCheck(repo string) int {
	diag := diagnoseSingleRepo(repo)
	renderDoctorDiagnosis(diag)
	if diag.IsHealthy {
		return 1
	}

	return 0
}

// RunCGDoctor diagnoses health issues and outputs actionable remediation commands.
func RunCGDoctor(repos []string) error {
	fmt.Printf("%s\n\n", cgAuditTitleStyle.Render("=== Coding Guidelines Doctor Diagnostics ==="))
	healthyCount := 0

	for _, repo := range repos {
		healthyCount += runDoctorCheck(repo)
	}

	renderDoctorOverall(len(repos), healthyCount)

	return nil
}
