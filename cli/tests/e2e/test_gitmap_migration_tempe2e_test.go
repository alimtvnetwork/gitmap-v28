//go:build tempe2e

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isTestGitmapE2ESkipped() bool {
	return os.Getenv("RUN_TEMP_E2E") != "1"
}

func TestGitmapTestRepo_ConfigDeclaration_TempE2E(t *testing.T) {
	if isTestGitmapE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	cfgPath := filepath.Join("..", "..", "..", ".ai-memory", "temp", "commit-pull-config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("failed to read commit-pull-config.json: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to parse commit-pull-config.json: %v", err)
	}

	target := raw["target"].(string)
	if target != "${target_path}" && !strings.Contains(target, "test-gitmap") {
		t.Fatalf("unexpected target in config: %v", target)
	}

	inputs, ok := raw["inputs"].([]any)
	if !ok || len(inputs) < 2 {
		t.Fatalf("expected at least 2 input specifications, got: %v", raw["inputs"])
	}
}

func TestGitmapTestRepo_StateTemplatesImport_TempE2E(t *testing.T) {
	if isTestGitmapE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	tplPath := filepath.Join("..", "..", "..", ".ai-memory", "temp", "seo-templates.json")
	if _, err := os.Stat(tplPath); err != nil {
		t.Fatalf("seo-templates.json missing: %v", err)
	}

	count, isSkipped, exportID, err := store.ImportTemplatesFromFile(tplPath, false)
	if err != nil {
		t.Fatalf("failed to import templates: %v", err)
	}
	t.Logf("Imported %d templates (skipped=%v, exportID=%s)", count, isSkipped, exportID)

	compiled, err := store.PrecompileTemplates("seo", map[string]string{
		"files.2.names": "cli/cmdagy, cli/store",
	})
	if err != nil || len(compiled) == 0 {
		t.Fatalf("expected compiled templates for seo category, got %d (err=%v)", len(compiled), err)
	}
}

func TestGitmapTestRepo_LiveRepoIntegrity_TempE2E(t *testing.T) {
	if isTestGitmapE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	targetDir := filepath.Join("D:", string(filepath.Separator), "test-gitmap", "test-gitmap")
	if _, err := os.Stat(filepath.Join(targetDir, ".git")); err != nil {
		t.Skipf("target repo %s does not exist on this machine; skipping live check", targetDir)
	}

	verifyWorkingTreeClean(t, targetDir)
	verifyCommitCount(t, targetDir)
	verifyRemoteURL(t, targetDir)
	verifyCommitMessageFormat(t, targetDir)
}

func verifyWorkingTreeClean(t *testing.T, targetDir string) {
	cmd := exec.Command("git", "-C", targetDir, "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git status failed: %v", err)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Fatalf("working tree is not clean:\n%s", string(out))
	}
}

func verifyCommitCount(t *testing.T, targetDir string) {
	cmd := exec.Command("git", "-C", targetDir, "rev-list", "--count", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("rev-list failed: %v", err)
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	if count < 1000 {
		t.Fatalf("expected at least 1000 commits in target repo, got %d", count)
	}
}

func verifyRemoteURL(t *testing.T, targetDir string) {
	cmd := exec.Command("git", "-C", targetDir, "remote", "get-url", "origin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("remote get-url failed: %v", err)
	}
	url := strings.TrimSpace(string(out))
	if !strings.Contains(url, "test-gitmap") {
		t.Fatalf("unexpected remote url: %s", url)
	}
}

func verifyCommitMessageFormat(t *testing.T, targetDir string) {
	cmd := exec.Command("git", "-C", targetDir, "log", "-n", "10", "--format=%s%n%b")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git log failed: %v", err)
	}
	body := string(out)
	if !strings.Contains(body, "RISEUP ASIA") && !strings.Contains(body, "Release") {
		t.Fatalf("expected commit messages to contain sponsor templates or releases")
	}
}

func TestGitmapTestRepo_SelfContainedSEOVariables_TempE2E(t *testing.T) {
	if isTestGitmapE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	tplPath := filepath.Join("..", "..", "..", ".ai-memory", "temp", "seo-templates.json")
	_, _, _, err := store.ImportTemplatesFromFile(tplPath, true)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}

	compiled, err := store.PrecompileTemplates("seo", nil)
	if err != nil || len(compiled) == 0 {
		t.Fatalf("expected compiled templates, got %d, err: %v", len(compiled), err)
	}

	first := compiled[0]
	if strings.Contains(first.Title, "$COMPANY") || strings.Contains(first.Text, "$COMPANY") {
		t.Fatalf("template still contains unexpanded $COMPANY: %s", first.Text)
	}
	upper := strings.ToUpper(first.Text)
	if !strings.Contains(upper, "RISEUP ASIA") && !strings.Contains(upper, "RISE UP ASIA") {
		t.Fatalf("expected self-contained variable expansion of COMPANY: %s", first.Text)
	}
}
