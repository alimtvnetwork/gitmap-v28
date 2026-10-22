// Package cmdcg provides coding guideline audit verification and reporting.
package cmdcg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

var (
	cgAuditTitleStyle = termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true)
	cgPassStyle       = termout.SafeTokyoStyle(termout.TokyoOK).Bold(true)
	cgWarnStyle       = termout.SafeTokyoStyle(termout.TokyoWarn).Bold(true)
	cgFailStyle       = termout.SafeTokyoStyle(termout.TokyoError).Bold(true)
	cgMutedStyle      = termout.SafeTokyoStyle(termout.TokyoMuted)
)

// CGAuditCheck represents the result of a single audit verification item.
type CGAuditCheck struct {
	Category string
	Name     string
	IsPass   bool
	Detail   string
}

// CGAuditSummary aggregates audit checks and compliance score for a repository.
type CGAuditSummary struct {
	RepoPath string
	Checks   []CGAuditCheck
	Score    int
}

func auditCheckVersionJSON(repoPath string) CGAuditCheck {
	data, err := os.ReadFile(filepath.Join(repoPath, "version.json"))
	if err != nil {
		return CGAuditCheck{Category: "Configuration", Name: "version.json valid", IsPass: false, Detail: "file missing"}
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return CGAuditCheck{Category: "Configuration", Name: "version.json valid", IsPass: false, Detail: "invalid JSON syntax"}
	}

	return CGAuditCheck{Category: "Configuration", Name: "version.json valid", IsPass: true, Detail: "valid JSON structure"}
}

func auditCheckCGMetadata(repoPath string) CGAuditCheck {
	meta, err := ReadCGMetadata(repoPath)
	if err != nil || meta == nil {
		return CGAuditCheck{Category: "Configuration", Name: "coding-guidelines section", IsPass: false, Detail: "metadata missing in version.json"}
	}

	return CGAuditCheck{Category: "Configuration", Name: "coding-guidelines section", IsPass: true, Detail: fmt.Sprintf("installed version %s", meta.Version)}
}

func auditCheckSpecDir(repoPath string) CGAuditCheck {
	dirPath := filepath.Join(repoPath, "02-spec", "02-coding-guidelines")
	stat, err := os.Stat(dirPath)
	if err != nil || !stat.IsDir() {
		return CGAuditCheck{Category: "Specifications", Name: "02-coding-guidelines directory", IsPass: false, Detail: "directory missing"}
	}

	return CGAuditCheck{Category: "Specifications", Name: "02-coding-guidelines directory", IsPass: true, Detail: "specification directory found"}
}

func hasRegularFile(filePath string) bool {
	stat, err := os.Stat(filePath)

	return err == nil && !stat.IsDir()
}

func auditCheckMemoryFiles(repoPath string) CGAuditCheck {
	hasAvoid := hasRegularFile(filepath.Join(repoPath, ".ai-memory", "strictly-avoid.md"))
	hasCG := hasRegularFile(filepath.Join(repoPath, ".ai-memory", "coding-guidelines.md"))
	if !hasAvoid && !hasCG {
		return CGAuditCheck{Category: "Memory", Name: "memory guidelines", IsPass: false, Detail: "strictly-avoid.md or coding-guidelines.md missing"}
	}

	return CGAuditCheck{Category: "Memory", Name: "memory guidelines", IsPass: true, Detail: "memory guidelines file found"}
}

func hasDirectory(dirPath string) bool {
	stat, err := os.Stat(dirPath)

	return err == nil && stat.IsDir()
}

func auditCheckRulesDir(repoPath string) CGAuditCheck {
	if !hasDirectory(filepath.Join(repoPath, ".agents", "rules")) {
		return CGAuditCheck{Category: "Rules", Name: ".agents/rules directory", IsPass: false, Detail: "rules directory missing"}
	}

	return CGAuditCheck{Category: "Rules", Name: ".agents/rules directory", IsPass: true, Detail: "rules directory found"}
}

func calculateAuditScore(checks []CGAuditCheck) int {
	if len(checks) == 0 {
		return 0
	}

	passCount := 0
	for _, c := range checks {
		if c.IsPass {
			passCount++
		}
	}

	return (passCount * 100) / len(checks)
}

func runAllAuditChecks(repoPath string) []CGAuditCheck {
	return []CGAuditCheck{
		auditCheckVersionJSON(repoPath),
		auditCheckCGMetadata(repoPath),
		auditCheckSpecDir(repoPath),
		auditCheckMemoryFiles(repoPath),
		auditCheckRulesDir(repoPath),
	}
}

func auditSingleRepo(repoPath string) CGAuditSummary {
	checks := runAllAuditChecks(repoPath)

	return CGAuditSummary{
		RepoPath: repoPath,
		Checks:   checks,
		Score:    calculateAuditScore(checks),
	}
}

// AuditSingleRepo evaluates coding guideline compliance for a single repo.
func AuditSingleRepo(repoPath string) CGAuditSummary {
	return auditSingleRepo(repoPath)
}

func renderAuditCheck(c CGAuditCheck) {
	tag := termout.TokyoOKPill("PASS")
	if !c.IsPass {
		tag = termout.TokyoErrorPill("FAIL")
	}

	fmt.Printf("    %s %-32s %s\n", tag, c.Name, cgMutedStyle.Render("("+c.Detail+")"))
}

func renderAuditScoreBadge(score int) string {
	if score >= 100 {
		return termout.TokyoOKPill(fmt.Sprintf("%d/100", score))
	}
	if score >= 60 {
		return termout.TokyoWarnPill(fmt.Sprintf("%d/100", score))
	}

	return termout.TokyoErrorPill(fmt.Sprintf("%d/100", score))
}

func renderAuditSummary(summary CGAuditSummary) {
	badge := renderAuditScoreBadge(summary.Score)
	name := filepath.Base(summary.RepoPath)
	progress := termout.RenderTokyoProgress(summary.Score, 100, 16)
	fmt.Printf("\n  %s %s [%s] %s\n", cgAuditTitleStyle.Render("● "+name), cgMutedStyle.Render("("+summary.RepoPath+")"), badge, progress)

	for _, check := range summary.Checks {
		renderAuditCheck(check)
	}
}

func maybeFixAuditRepo(summary CGAuditSummary, isFix bool) {
	if !isFix || summary.Score == 100 {
		return
	}

	fmt.Printf("    %s Auto-remediation triggered for %s\n", termout.TokyoWarnPill("FIX"), filepath.Base(summary.RepoPath))
	runCGUpdateAction([]string{summary.RepoPath})
}

// RunCGAudit iterates through repositories and prints compliance reports.
func RunCGAudit(repos []string, isFix bool) error {
	fmt.Printf("%s\n", cgAuditTitleStyle.Render("=== Coding Guidelines Multi-Repo Audit ==="))

	for _, repo := range repos {
		summary := auditSingleRepo(repo)
		renderAuditSummary(summary)
		maybeFixAuditRepo(summary, isFix)
	}

	fmt.Printf("\n%s\n", cgPassStyle.Render("Audit run completed successfully."))

	return nil
}
