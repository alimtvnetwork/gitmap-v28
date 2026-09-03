package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// CICDCheckResult holds the outcome of an internal CI/CD diagnostic probe.
type CICDCheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Fixed   bool   `json:"fixed"`
	Detail  string `json:"detail"`
	FixHint string `json:"fixHint,omitempty"`
}

// runInternalCICDChecks orchestrates all local verification scripts and optional auto-repairs.
func runInternalCICDChecks(autoFix bool) []CICDCheckResult {
	fmt.Printf("\n%s● Running Internal CI/CD Diagnostic & Repair Suite...%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s\n", strings.Repeat("─", 74))

	results := []CICDCheckResult{
		probeGofmtFormat(autoFix),
		probeScriptCheck("linter-scripts/check-nested-ifs.py", "Nested If Linter", "Flatten nested ifs per spec/02-coding-guidelines/"),
		probeScriptCheck("linter-scripts/check-enum-and-boolean.py", "Boolean & Enum Linter", "Audit boolean conventions and enum Type suffixes"),
		probeScriptCheck("linter-scripts/check-relative-paths.py", "Relative Paths Linter", "Replace absolute drive letters with relative paths"),
		probeScriptCheck("linter-scripts/check-error-management.py", "Error Management Linter", "Wrap errors with AppError and avoid swallowed errors"),
		probeLegacyRefsCheck(),
		probeSpellCheck(),
		probeCompileGate(),
	}

	for _, r := range results {
		printCheckResultLine(r)
	}
	fmt.Printf("  %s\n\n", strings.Repeat("─", 74))

	return results
}

func printCheckResultLine(r CICDCheckResult) {
	if r.Fixed {
		fmt.Printf("  %s✔ FIXED%s  %-28s %s\n", constants.ColorYellow, constants.ColorReset, r.Name, r.Detail)
		return
	}
	if r.Passed {
		fmt.Printf("  %s✔ PASS %s  %-28s %s\n", constants.ColorGreen, constants.ColorReset, r.Name, r.Detail)
		return
	}
	fmt.Printf("  %s✖ FAIL %s  %-28s %s\n", constants.ColorRed, constants.ColorReset, r.Name, r.Detail)
	if len(r.FixHint) > 0 {
		fmt.Printf("           %sfix:%s %s\n", constants.ColorYellow, constants.ColorReset, r.FixHint)
	}
}

func resolveGitmapRoot() (string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		return ".", "."
	}
	root := findRepoRoot(cwd)
	if root != "" {
		return resolveGitmapDirFromRoot(root)
	}
	if _, err := os.Stat("version.json"); err == nil {
		return ".", "gitmap"
	}
	if _, err := os.Stat("../version.json"); err == nil {
		return "..", "."
	}
	return ".", "."
}

func resolveGitmapDirFromRoot(root string) (string, string) {
	gDir := filepath.Join(root, "gitmap")
	if _, err := os.Stat(gDir); err == nil {
		return root, gDir
	}
	return root, root
}

func probeGofmtFormat(autoFix bool) CICDCheckResult {
	_, gitmapDir := resolveGitmapRoot()
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = gitmapDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CICDCheckResult{
			Name:    "gofmt formatting",
			Passed:  false,
			Detail:  "gofmt not found on PATH",
			FixHint: "Install Go toolchain from https://go.dev/dl/",
		}
	}

	lines := strings.TrimSpace(string(out))
	if len(lines) == 0 {
		return CICDCheckResult{
			Name:   "gofmt formatting",
			Passed: true,
			Detail: "All Go files are gofmt-clean",
		}
	}

	unformatted := strings.Split(lines, "\n")
	if autoFix {
		return applyGofmtFix(gitmapDir, unformatted)
	}

	return CICDCheckResult{
		Name:    "gofmt formatting",
		Passed:  false,
		Detail:  fmt.Sprintf("%d unformatted .go file(s)", len(unformatted)),
		FixHint: "Run 'gitmap pipeline errorlogs --fix' or 'cd gitmap && gofmt -w .'",
	}
}

func applyGofmtFix(gitmapDir string, files []string) CICDCheckResult {
	cmd := exec.Command("gofmt", "-w", ".")
	cmd.Dir = gitmapDir
	err := cmd.Run()
	if err != nil {
		return CICDCheckResult{
			Name:    "gofmt formatting",
			Passed:  false,
			Detail:  fmt.Sprintf("Failed to auto-format %d file(s): %v", len(files), err),
			FixHint: "Check file permissions and run 'gofmt -w .' manually",
		}
	}
	return CICDCheckResult{
		Name:   "gofmt formatting",
		Passed: true,
		Fixed:  true,
		Detail: fmt.Sprintf("Auto-formatted %d .go file(s)", len(files)),
	}
}

func probeScriptCheck(scriptRelPath, checkName, fixHint string) CICDCheckResult {
	rootDir, _ := resolveGitmapRoot()
	fullPath := filepath.Join(rootDir, scriptRelPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return CICDCheckResult{
			Name:   checkName,
			Passed: true,
			Detail: "Script omitted in current environment (skipped)",
		}
	}

	cmd := exec.Command("python", fullPath)
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		firstLine := extractFirstLine(string(out))
		return CICDCheckResult{
			Name:    checkName,
			Passed:  false,
			Detail:  firstLine,
			FixHint: fixHint,
		}
	}

	return CICDCheckResult{
		Name:   checkName,
		Passed: true,
		Detail: "Zero violations found",
	}
}

func probeLegacyRefsCheck() CICDCheckResult {
	rootDir, _ := resolveGitmapRoot()
	scriptPath := filepath.Join(rootDir, ".github", "scripts", "check-legacy-refs.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return CICDCheckResult{
			Name:   "Legacy Refs Check",
			Passed: true,
			Detail: "Script omitted (skipped)",
		}
	}

	cmd := exec.Command("python", scriptPath, ".")
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CICDCheckResult{
			Name:    "Legacy Refs Check",
			Passed:  false,
			Detail:  extractFirstLine(string(out)),
			FixHint: "Remove forbidden legacy repository references",
		}
	}

	return CICDCheckResult{
		Name:   "Legacy Refs Check",
		Passed: true,
		Detail: "No forbidden legacy refs found",
	}
}

func probeSpellCheck() CICDCheckResult {
	rootDir, _ := resolveGitmapRoot()
	scriptPath := filepath.Join(rootDir, ".github", "scripts", "misspell-changed.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return CICDCheckResult{
			Name:   "Spell Check (US locale)",
			Passed: true,
			Detail: "Script omitted (skipped)",
		}
	}

	cmd := exec.Command("python", scriptPath)
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CICDCheckResult{
			Name:    "Spell Check (US locale)",
			Passed:  false,
			Detail:  extractFirstLine(string(out)),
			FixHint: "Correct spelling to standard American English",
		}
	}

	return CICDCheckResult{
		Name:   "Spell Check (US locale)",
		Passed: true,
		Detail: "All changed files passed US spell check",
	}
}

func probeCompileGate() CICDCheckResult {
	_, gitmapDir := resolveGitmapRoot()
	cmd := exec.Command("go", "test", "-run=^$", "./...", "-count=1")
	cmd.Dir = gitmapDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CICDCheckResult{
			Name:    "Go Compile Gate",
			Passed:  false,
			Detail:  extractFirstLine(string(out)),
			FixHint: "Fix package or test compilation errors",
		}
	}

	return CICDCheckResult{
		Name:   "Go Compile Gate",
		Passed: true,
		Detail: "All packages and test suites compiled cleanly",
	}
}

func extractFirstLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	trimmed := strings.TrimSpace(lines[0])
	hasContent := len(lines) > 0 && len(trimmed) > 0
	if hasContent {
		return trimmed
	}
	return "Check failed"
}
