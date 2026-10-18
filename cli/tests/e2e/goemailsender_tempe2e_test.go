//go:build e2e || tempe2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

func requireE2EEnv(t *testing.T) {
	t.Helper()
	if os.Getenv("GITMAP_RUN_E2E") != "1" && os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping e2e integration test; run with GITMAP_RUN_E2E=1")
	}
}

func TestE2E_GoEmailReaderBinaryLookupAndHelp(t *testing.T) {
	requireE2EEnv(t)
	binPath, isFound := locateGoEmailReaderBinary()
	if !isFound {
		t.Skip("go-email-reader binary not found")
	}
	out, err := exec.Command(binPath, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("go-email-reader --help execution failed: %v, out: %s", err, string(out))
	}
	assertHelpOutputContent(t, string(out))
	uiHelpErr := cmdspace.RunSpaceSendEmailUI([]string{"--help"})
	if uiHelpErr != nil {
		t.Fatalf("gitmap space email ui --help failed: %v", uiHelpErr)
	}
}

func locateGoEmailReaderBinary() (string, bool) {
	if p, err := exec.LookPath("go-email-reader"); err == nil {
		return p, true
	}
	if p, err := exec.LookPath("go-email-reader.exe"); err == nil {
		return p, true
	}
	return checkHomeGoEmailBinary()
}

func checkHomeGoEmailBinary() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	cand := filepath.Join(home, "go", "bin", "go-email-reader.exe")
	_, statErr := os.Stat(cand)
	if statErr == nil {
		return cand, true
	}
	return "", false
}

func assertHelpOutputContent(t *testing.T, text string) {
	t.Helper()
	for _, expected := range []string{"go-email-reader", "send", "list"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected --help to contain %q, got: %s", expected, text)
		}
	}
}

func TestE2E_GitMapEmailSchedulingRegistrationAndTriggerMock(t *testing.T) {
	requireE2EEnv(t)
	setupE2ETriggerTestDB(t)

	regArgs := []string{"--to", "alerts@company.org", "--task", "release", "--name", "e2e-rel-trigger"}
	if err := cmdspace.RunScheduleSendEmail(regArgs); err != nil {
		t.Fatalf("register trigger error: %v", err)
	}

	outcome := cmdspace.InferTaskOutcome(cmdspace.EmailTaskRelease, nil)
	if outcome.IsSuccess == false || outcome.Operation != "release" {
		t.Fatalf("unexpected task outcome: %+v", outcome)
	}
	cmdspace.FireEmailTaskTriggers(cmdspace.EmailTaskRelease, outcome)

	if listErr := cmdspace.RunScheduleSendEmail([]string{"--list"}); listErr != nil {
		t.Fatalf("list triggers error: %v", listErr)
	}

	if remErr := cmdspace.RunScheduleSendEmail([]string{"--remove", "--name", "e2e-rel-trigger"}); remErr != nil {
		t.Fatalf("remove trigger error: %v", remErr)
	}
}

func setupE2ETriggerTestDB(t *testing.T) {
	t.Helper()
	tempDir := t.TempDir()
	store.SetBinaryDataDirForTesting(tempDir)
	t.Cleanup(func() { store.SetBinaryDataDirForTesting("") })
	db, err := store.OpenDefault()
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("failed to migrate store: %v", err)
	}
	_ = db.Close()
}

func TestE2E_MultiProfileFallbackMock(t *testing.T) {
	requireE2EEnv(t)
	tempDir := t.TempDir()
	store.SetBinaryDataDirForTesting(tempDir)
	defer store.SetBinaryDataDirForTesting("")

	seedE2EMultiProfiles(t)
	testRepoDir := filepath.Join(tempDir, "mock-repo")
	_ = os.MkdirAll(testRepoDir, 0755)

	winningProfile := simulateMultiProfileProbe(t, "https://github.com/org/private-repo.git")
	if winningProfile != "backup-profile" {
		t.Fatalf("expected winning profile backup-profile, got: %s", winningProfile)
	}

	if err := usercontext.BindProject(testRepoDir, winningProfile); err != nil {
		t.Fatalf("usercontext.BindProject failed: %v", err)
	}

	resolved, hasProfile := usercontext.ResolveProjectUser(testRepoDir)
	if !hasProfile || resolved == nil || resolved.Name != "backup-profile" {
		t.Fatalf("unexpected bound profile resolution: has=%t prof=%+v", hasProfile, resolved)
	}
}

func seedE2EMultiProfiles(t *testing.T) {
	t.Helper()
	cfg := model.GitProfileConfig{
		Default: "default-profile",
		Profiles: []model.GitProfile{
			{Name: "default-profile", Email: "default@test.org", IsDefault: true, Provider: "github", Type: "user"},
			{Name: "backup-profile", Email: "backup@test.org", IsDefault: false, Provider: "github", Type: "user"},
		},
	}
	if err := store.SaveGitProfiles(cfg); err != nil {
		t.Fatalf("seed profiles failed: %v", err)
	}
}

func simulateMultiProfileProbe(t *testing.T, targetURL string) string {
	t.Helper()
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		t.Fatalf("failed to load profiles: %v", err)
	}
	candidates := filterSecondaryProfiles(cfg)
	if len(candidates) > 0 {
		return candidates[0].Name
	}
	return ""
}

func filterSecondaryProfiles(cfg model.GitProfileConfig) []model.GitProfile {
	var secondaries []model.GitProfile
	for _, p := range cfg.Profiles {
		isDefault := p.IsDefault || p.Name == cfg.Default
		if !isDefault {
			secondaries = append(secondaries, p)
		}
	}
	return secondaries
}
