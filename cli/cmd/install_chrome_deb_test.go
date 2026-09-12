package cmd

import (
	"slices"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestBuildChromeUpdateCmd(t *testing.T) {
	cmd := buildChromeUpdateCmd()
	expected := []string{"sudo", "apt-get", "update"}

	if !slices.Equal(cmd, expected) {
		t.Fatalf("expected %v, got %v", expected, cmd)
	}
}

func TestBuildChromeFetchUtilsCmd(t *testing.T) {
	cmd := buildChromeFetchUtilsCmd()
	expected := []string{"sudo", "apt-get", "install", "-y", "wget", "curl"}

	if !slices.Equal(cmd, expected) {
		t.Fatalf("expected %v, got %v", expected, cmd)
	}
}

func TestBuildChromeDownloadCmd(t *testing.T) {
	cmd := buildChromeDownloadCmd()

	if len(cmd) < 5 {
		t.Fatalf("expected at least 5 args in download cmd, got %v", cmd)
	}

	if !slices.Contains(cmd, chromeDebStage) {
		t.Errorf("expected cmd to contain stage path %q, got %v", chromeDebStage, cmd)
	}

	if !slices.Contains(cmd, chromeDebURL) {
		t.Errorf("expected cmd to contain deb URL %q, got %v", chromeDebURL, cmd)
	}
}

func TestBuildChromeInstallCmd(t *testing.T) {
	cmd := buildChromeInstallCmd()
	expected := []string{"sudo", "apt-get", "install", "-y", chromeDebStage}

	if !slices.Equal(cmd, expected) {
		t.Fatalf("expected %v, got %v", expected, cmd)
	}
}

func TestBuildChromeCleanupCmd(t *testing.T) {
	cmd := buildChromeCleanupCmd()
	expected := []string{"rm", "-f", chromeDebStage}

	if !slices.Equal(cmd, expected) {
		t.Fatalf("expected %v, got %v", expected, cmd)
	}
}

func TestBuildChromeVerifyCmd(t *testing.T) {
	cmd := buildChromeVerifyCmd()
	expected := []string{"google-chrome", "--version"}

	if !slices.Equal(cmd, expected) {
		t.Fatalf("expected %v, got %v", expected, cmd)
	}
}

func TestGetChromePipelineSteps(t *testing.T) {
	steps := getChromePipelineSteps()

	if len(steps) != 6 {
		t.Fatalf("expected 6 pipeline steps, got %d", len(steps))
	}

	verifyPipelineStepPhases(t, steps)
}

func verifyPipelineStepPhases(t *testing.T, steps []chromePhaseStep) {
	expected := []string{
		phaseChromeUpdate, phaseChromeFetchUtils, phaseChromeDownload,
		phaseChromeInstall, phaseChromeCleanup, phaseChromeVerify,
	}

	for i, phase := range expected {
		if steps[i].PhaseName != phase {
			t.Errorf("step %d: expected phase %s, got %s", i, phase, steps[i].PhaseName)
		}

		if len(steps[i].Command) == 0 {
			t.Errorf("step %d: command slice is empty", i)
		}
	}
}

func TestResolveChromeTool(t *testing.T) {
	if resolveChromeTool("") != constants.CmdChrome {
		t.Errorf("expected default tool %q, got %q", constants.CmdChrome, resolveChromeTool(""))
	}

	if resolveChromeTool("google-chrome") != "google-chrome" {
		t.Errorf("expected explicit tool to be preserved")
	}
}

func TestResolveChromeVersion(t *testing.T) {
	if resolveChromeVersion("") != "official-deb" {
		t.Errorf("expected default version 'official-deb', got %q", resolveChromeVersion(""))
	}

	if resolveChromeVersion("130.0.0") != "130.0.0" {
		t.Errorf("expected explicit version to be preserved")
	}
}

func TestCheckChromeDryRun_True(t *testing.T) {
	opts := installOptions{DryRun: true}

	if !checkChromeDryRun(opts) {
		t.Errorf("expected checkChromeDryRun to return true for DryRun=true")
	}
}

func TestCheckChromeDryRun_False(t *testing.T) {
	opts := installOptions{DryRun: false}

	if checkChromeDryRun(opts) {
		t.Errorf("expected checkChromeDryRun to return false for DryRun=false")
	}
}

func TestIsBinaryOnPath(t *testing.T) {
	if isBinaryOnPath("definitely-not-a-real-binary-xyz-98765") {
		t.Errorf("expected nonexistent binary to return false")
	}
}
