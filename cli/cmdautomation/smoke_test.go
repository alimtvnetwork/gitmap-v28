package cmdautomation

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunSmokeTest validates installer scripts or installed tools across platforms.
func RunSmokeTest(opts SmokeTestOptions) SmokeTestResultMonad {
	start := time.Now()
	items := resolveSmokeTargets(opts)
	res := aggregateSmokeResult(items, start)
	return result.Ok(res)
}

func resolveSmokeTargets(opts SmokeTestOptions) []SmokeTestItem {
	if len(opts.Tools) > 0 {
		return testAllTools(opts.Tools)
	}
	root := resolveAuditDir(opts.Dir)
	files, _ := collectTextFiles(root, []string{".sh", ".ps1"})
	return testAllScripts(files, opts.Filter)
}

func testAllTools(tools []string) []SmokeTestItem {
	var items []SmokeTestItem
	for _, t := range tools {
		items = append(items, smokeTestTool(t))
	}
	return items
}

func testAllScripts(files []string, filter string) []SmokeTestItem {
	var items []SmokeTestItem
	filt := strings.ToLower(filter)
	for _, f := range files {
		matches := len(filt) == 0 || strings.Contains(strings.ToLower(f), filt)
		if matches {
			items = append(items, smokeTestScript(f))
		}
	}
	return items
}

func smokeTestScript(path string) SmokeTestItem {
	start := time.Now()
	data, err := os.ReadFile(path)
	if err != nil {
		return SmokeTestItem{
			Name:     path,
			Target:   path,
			Kind:     "installer_script",
			IsFail:   true,
			Duration: time.Since(start),
			Message:  "Failed to read script file",
		}
	}
	issues := checkScriptIssues(path, string(data))
	isFail := len(issues) > 0
	msg := "Valid installer patterns (clean tokens, SHA256 check, safe rename)"
	if isFail {
		msg = fmt.Sprintf("Found %d installer pattern issue(s)", len(issues))
	}
	return SmokeTestItem{
		Name:     path,
		Target:   path,
		Kind:     "installer_script",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
		Issues:   issues,
	}
}

func checkScriptIssues(path, content string) []string {
	var issues []string
	if strings.Contains(content, "PLACEHOLDER") {
		issues = append(issues, "Contains unreplaced PLACEHOLDER tokens")
	}
	hasSha := strings.Contains(strings.ToLower(content), "sha256") ||
		strings.Contains(strings.ToLower(content), "shasum") ||
		strings.Contains(strings.ToLower(content), "get-filehash")
	if !hasSha {
		issues = append(issues, "Missing SHA256 checksum verification pattern")
	}
	hasSafeRename := strings.Contains(content, "mv ") ||
		strings.Contains(content, "install -m") ||
		strings.Contains(strings.ToLower(content), "move-item") ||
		strings.Contains(strings.ToLower(content), "rename-item")
	if !hasSafeRename {
		issues = append(issues, "Missing non-destructive replacement logic (rename-first)")
	}
	return issues
}

func smokeTestTool(tool string) SmokeTestItem {
	start := time.Now()
	path, err := exec.LookPath(tool)
	if err != nil {
		return SmokeTestItem{
			Name:     tool,
			Target:   tool,
			Kind:     "tool_binary",
			IsFail:   true,
			Duration: time.Since(start),
			Message:  "Tool not found in PATH",
		}
	}
	cmd := exec.Command(path, "--version")
	out, runErr := cmd.Output()
	isFail := runErr != nil
	msg := strings.TrimSpace(string(out))
	if isFail {
		msg = fmt.Sprintf("Failed to run --version: %v", runErr)
	}
	return SmokeTestItem{
		Name:     tool,
		Target:   path,
		Kind:     "tool_binary",
		IsPass:   !isFail,
		IsFail:   isFail,
		Duration: time.Since(start),
		Message:  msg,
	}
}

func aggregateSmokeResult(items []SmokeTestItem, start time.Time) SmokeTestResult {
	passed := 0
	failed := 0
	for _, it := range items {
		if it.IsPass {
			passed++
		}
		if it.IsFail {
			failed++
		}
	}
	return SmokeTestResult{
		TotalItems:  len(items),
		PassedItems: passed,
		FailedItems: failed,
		Items:       items,
		Duration:    time.Since(start),
		IsPass:      failed == 0,
	}
}
