// Package cmdcg provides diagnostic health checks and repair suggestions.
package cmdcg

import (
	"fmt"
	"path/filepath"
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
		fmt.Printf("  %s %s: Healthy (all checks passed)\n", cgPassStyle.Render("✔"), cgAuditTitleStyle.Render(name))

		return
	}

	fmt.Printf("  %s %s: %d issue(s) detected\n", cgFailStyle.Render("✖"), cgAuditTitleStyle.Render(name), len(diag.Issues))
	for i, issue := range diag.Issues {
		fmt.Printf("    • %s %s\n", cgWarnStyle.Render("Issue:"), issue)
		fmt.Printf("      %s %s\n", cgPassStyle.Render("Fix:"), diag.Suggestions[i])
	}
}

func renderDoctorOverall(totalRepos, healthyRepos int) {
	issuesCount := totalRepos - healthyRepos
	fmt.Printf("\n%s\n", cgAuditTitleStyle.Render("── Doctor Summary ──"))
	fmt.Printf("  Total scanned: %d | Healthy: %d | Needs action: %d\n", totalRepos, healthyRepos, issuesCount)
	if issuesCount == 0 {
		fmt.Printf("  %s All repositories are in optimal state.\n", cgPassStyle.Render("✔"))

		return
	}

	fmt.Printf("  %s Action required on %d repository(ies).\n", cgWarnStyle.Render("!"), issuesCount)
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
