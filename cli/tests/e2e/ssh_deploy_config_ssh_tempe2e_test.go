//go:build tempe2e

package e2e

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

func isDeployConfigTempE2ESkipped() bool {
	return os.Getenv("RUN_TEMP_E2E") != "1"
}

func TestDeployConfigSSH_DryRun_TempE2E(t *testing.T) {
	if isDeployConfigTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	err := cmdssh.RunSSHDeployConfigSSHCLI([]string{"all", "--dry-run"})
	if err != nil {
		t.Fatalf("expected dry-run deploy config ssh to succeed, got: %v", err)
	}
}

func TestDeployConfigSSH_FileResolve_TempE2E(t *testing.T) {
	if isDeployConfigTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	err := cmdssh.RunSSHDeployConfigSSHCLI([]string{"--file", "01-gitmap", "--dry-run"})
	if err != nil {
		t.Fatalf("expected dry-run deploy with repo-secrets token '01-gitmap' to succeed, got: %v", err)
	}
}

func TestDeployConfigSSH_LiveExecution_TempE2E(t *testing.T) {
	if isDeployConfigTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	err := cmdssh.RunSSHDeployConfigSSHCLI([]string{"w1"})
	if err != nil {
		t.Fatalf("expected live deploy config ssh to w1 to succeed, got: %v", err)
	}
}

func TestCloneSecretsResolver_TempE2E(t *testing.T) {
	if isDeployConfigTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	for _, token := range []string{"w1", "w2", "w3"} {
		manifest := secrets.ResolveRepoSecretsManifest(token, "gitmap.json")
		if manifest == "" {
			t.Fatalf("expected manifest to resolve for token %q, got empty string", token)
		}
		if _, statErr := os.Stat(manifest); statErr != nil {
			t.Fatalf("resolved manifest does not exist on disk for %q: %s", token, manifest)
		}
	}
}
