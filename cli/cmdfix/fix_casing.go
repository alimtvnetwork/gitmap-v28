// Package cmdfix — fix_casing.go handles automated 2-step lowercase rename and Git index casing reconciliation.
package cmdfix

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CasingDiscrepancy represents a case divergence between disk/index and canonical conventions.
type CasingDiscrepancy struct {
	RepoPath      string `json:"repoPath"`
	CurrentPath   string `json:"currentPath"`
	CanonicalPath string `json:"canonicalPath"`
	IsDualTracked bool   `json:"isDualTracked"`
}

// RemediationRecipeRecord logs the execution details of a remediation recipe.
type RemediationRecipeRecord struct {
	Timestamp     string   `json:"timestamp"`
	RecipeID      string   `json:"recipeId"`
	Repository    string   `json:"repository"`
	RepoPath      string   `json:"repoPath"`
	Issue         string   `json:"issue"`
	Strategy      string   `json:"strategy"`
	Steps         []string `json:"steps"`
	Result        string   `json:"result"`
	CanonicalPath string   `json:"canonicalPath"`
}

// DetectCasingDiscrepancies inspects the repository Git index and disk for readme casing collisions.
func DetectCasingDiscrepancies(repoPath string) ([]CasingDiscrepancy, error) {
	if !isGitDirectory(repoPath) {
		return nil, nil
	}

	trackedFiles, err := listTrackedFiles(repoPath)
	if err != nil {
		return detectCasingFromDisk(repoPath), nil
	}

	discrepancies := inspectCasingDiscrepancies(repoPath, trackedFiles)
	if len(discrepancies) > 0 {
		return discrepancies, nil
	}

	return detectCasingFromDisk(repoPath), nil
}

func listTrackedFiles(repoPath string) ([]string, error) {
	out, err := exec.Command("git", "-C", repoPath, "ls-files").Output()
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(string(out), "\n") {
		if trimmed := strings.TrimSpace(line); len(trimmed) > 0 {
			files = append(files, trimmed)
		}
	}

	return files, nil
}

func inspectCasingDiscrepancies(repoPath string, tracked []string) []CasingDiscrepancy {
	hasLower, hasUpper, upperPath := inspectReadmeDuality(tracked)
	if !hasUpper {
		return nil
	}

	return []CasingDiscrepancy{{
		RepoPath:      repoPath,
		CurrentPath:   upperPath,
		CanonicalPath: buildCanonicalPath(upperPath),
		IsDualTracked: hasLower && hasUpper,
	}}
}

func inspectReadmeDuality(tracked []string) (bool, bool, string) {
	hasLower, hasUpper, upperPath := false, false, ""
	for _, f := range tracked {
		if isVendorOrDependencyPath(f) {
			continue
		}
		base := filepath.Base(f)
		isReadme := strings.EqualFold(base, "readme.md")
		if isReadme && base == "readme.md" {
			hasLower = true
		}
		if isReadme && base != "readme.md" {
			hasUpper, upperPath = true, f
		}
	}

	return hasLower, hasUpper, upperPath
}

func isVendorOrDependencyPath(f string) bool {
	clean := filepath.ToSlash(f)
	if strings.HasPrefix(clean, "vendor/") || strings.Contains(clean, "/vendor/") {
		return true
	}
	if strings.HasPrefix(clean, "node_modules/") || strings.Contains(clean, "/node_modules/") {
		return true
	}
	if strings.HasPrefix(clean, ".git/") || strings.Contains(clean, "/.git/") {
		return true
	}

	return false
}

func buildCanonicalPath(rawPath string) string {
	dir := filepath.Dir(rawPath)
	if dir == "." {
		return "readme.md"
	}

	return filepath.Join(dir, "readme.md")
}

func detectCasingFromDisk(repoPath string) []CasingDiscrepancy {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if strings.EqualFold(e.Name(), "readme.md") && e.Name() != "readme.md" {
			return []CasingDiscrepancy{{
				RepoPath: repoPath, CurrentPath: e.Name(), CanonicalPath: "readme.md",
			}}
		}
	}

	return nil
}

// ExecuteCasingNormalization resolves casing collisions via 2-step rename or index purge.
func ExecuteCasingNormalization(repoPath string, item CasingDiscrepancy) error {
	if !isGitDirectory(repoPath) {
		return apperror.NewValidationError("cannot normalize casing: directory is not a Git repository")
	}

	if item.IsDualTracked {
		return executeDualTrackedPurge(repoPath, item)
	}

	return executeTwoStepCasingRename(repoPath, item)
}

func executeDualTrackedPurge(repoPath string, item CasingDiscrepancy) error {
	diskPath := filepath.Join(repoPath, item.CanonicalPath)
	content, _ := os.ReadFile(diskPath)
	_ = exec.Command("git", "-C", repoPath, "rm", "--cached", "-f", item.CurrentPath).Run()
	if len(content) > 0 {
		_ = os.WriteFile(diskPath, content, 0644)
	}

	return exec.Command("git", "-C", repoPath, "add", "-A").Run()
}

func executeTwoStepCasingRename(repoPath string, item CasingDiscrepancy) error {
	tmpPath := item.CurrentPath + ".tmp-lcf"
	errGit := executeGitTwoStepMv(repoPath, item.CurrentPath, tmpPath, item.CanonicalPath)
	if errGit == nil {
		stageGitAddAll(repoPath)

		return nil
	}

	return executeFSTwoStepRename(repoPath, item.CurrentPath, tmpPath, item.CanonicalPath)
}

func executeGitTwoStepMv(repoPath, src, tmp, dst string) error {
	_ = os.Remove(filepath.Join(repoPath, tmp))
	if err := runGitMv(repoPath, src, tmp); err != nil {
		return err
	}
	if err := runGitMv(repoPath, tmp, dst); err != nil {
		rollbackGitMv(repoPath, tmp, src)

		return apperror.WrapSimple(err, "git mv step 2 failed")
	}

	return nil
}

func runGitMv(repoPath, src, dst string) error {
	return exec.Command("git", "-C", repoPath, "mv", src, dst).Run()
}

func rollbackGitMv(repoPath, tmp, src string) {
	_ = runGitMv(repoPath, tmp, src)
}

func executeFSTwoStepRename(repoPath, src, tmp, dst string) error {
	s, t, d := filepath.Join(repoPath, src), filepath.Join(repoPath, tmp), filepath.Join(repoPath, dst)
	_ = os.Remove(t)
	if err := os.Rename(s, t); err != nil {
		return apperror.WrapSimple(err, "fs rename step 1 failed")
	}
	if err := os.Rename(t, d); err != nil {
		_ = os.Rename(t, s)

		return apperror.WrapSimple(err, "fs rename step 2 failed")
	}
	stageGitAddAll(repoPath)

	return nil
}

func stageGitAddAll(repoPath string) {
	_ = exec.Command("git", "-C", repoPath, "add", "-A").Run()
}

func commitCasingNormalization(repoPath string, d CasingDiscrepancy) {
	out, err := exec.Command("git", "-C", repoPath, "status", "--porcelain").Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return
	}

	msg := fmt.Sprintf("chore: normalize casing %s to %s", d.CurrentPath, d.CanonicalPath)
	_ = exec.Command("git", "-C", repoPath, "commit", "-m", msg).Run()
}

// RecordCasingRemediationRecipe persists the executed casing recipe to repo-cache and local state.
func RecordCasingRemediationRecipe(item CasingDiscrepancy) error {
	recordLocalOutputState(item)
	recipePath := resolveRemediationRecipesFile()
	records := loadExistingRemediationRecipes(recipePath)
	repoName := filepath.Base(item.RepoPath)
	records = append(records, buildRemediationRecipeRecord(item, repoName))

	return persistRemediationRecipes(recipePath, records)
}

func buildRemediationRecipeRecord(item CasingDiscrepancy, repoName string) RemediationRecipeRecord {
	s1 := fmt.Sprintf("git mv %s %s.tmp-lcf", item.CurrentPath, item.CurrentPath)
	s2 := fmt.Sprintf("git mv %s.tmp-lcf %s", item.CurrentPath, item.CanonicalPath)

	return RemediationRecipeRecord{
		Timestamp: time.Now().UTC().Format(time.RFC3339), RecipeID: fmt.Sprintf("rcp-case-readme-%d", time.Now().Unix()),
		Repository: repoName, RepoPath: item.RepoPath, Issue: "NTFS core.ignorecase README.md vs readme.md collision",
		Strategy: "two_step_git_mv", Steps: []string{s1, s2, "git add -A"}, Result: "success", CanonicalPath: item.CanonicalPath,
	}
}

func recordLocalOutputState(item CasingDiscrepancy) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".gitmap", "output")
	_ = os.MkdirAll(dir, 0755)
	b, _ := json.MarshalIndent(buildRemediationRecipeRecord(item, filepath.Base(item.RepoPath)), "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "last_remediation.json"), b, 0644)
}

func resolveRemediationRecipesFile() string {
	for _, c := range []string{"repo-cache", "../repo-cache", "../../repo-cache"} {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)

			return filepath.Join(abs, "01-gitmap", "remediation_recipes.json")
		}
	}

	return filepath.Join("repo-cache", "01-gitmap", "remediation_recipes.json")
}

func loadExistingRemediationRecipes(path string) []RemediationRecipeRecord {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var records []RemediationRecipeRecord
	_ = json.Unmarshal(b, &records)

	return records
}

func persistRemediationRecipes(path string, records []RemediationRecipeRecord) error {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}

func resolveCasingTargetPath(args []string) string {
	if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		return args[1]
	}

	return "."
}

// RunFixCasing runs targeted casing remediation on a directory or repository.
func RunFixCasing(args []string) error {
	target := resolveCasingTargetPath(args)
	discrepancies, err := DetectCasingDiscrepancies(target)
	if err != nil {
		return err
	}
	if len(discrepancies) == 0 {
		fmt.Printf("%s No casing discrepancies found in %s\n", constants.ColorGreen+"✓"+constants.ColorReset, target)

		return nil
	}

	return ReconcileDiscrepancyList(discrepancies)
}

// ReconcileDiscrepancyList processes and records casing normalizations.
func ReconcileDiscrepancyList(discrepancies []CasingDiscrepancy) error {
	for _, d := range discrepancies {
		if err := ExecuteCasingNormalization(d.RepoPath, d); err != nil {
			return err
		}

		commitCasingNormalization(d.RepoPath, d)
		_ = RecordCasingRemediationRecipe(d)
		repoName := filepath.Base(d.RepoPath)
		fmt.Printf("  • Normalized casing: %s -> %s (%s)\n",
			d.CurrentPath, d.CanonicalPath, repoName)
	}

	return nil
}
