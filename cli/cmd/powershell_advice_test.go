package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestDetectPowerShellTarget_Directory(t *testing.T) {
	tempDir := t.TempDir()
	kind, resolved := detectPowerShellTarget(tempDir)
	if kind != TargetKindDirectory {
		t.Fatalf("expected TargetKindDirectory, got: %s", kind)
	}

	if resolved != tempDir {
		t.Fatalf("expected resolved path %s, got: %s", tempDir, resolved)
	}
}

func TestDetectPowerShellTarget_DocumentFile(t *testing.T) {
	tempDir := t.TempDir()
	docFile := filepath.Join(tempDir, "spec.md")
	_ = os.WriteFile(docFile, []byte("# Spec"), 0644)

	kind, resolved := detectPowerShellTarget(docFile)
	if kind != TargetKindDocument {
		t.Fatalf("expected TargetKindDocument, got: %s", kind)
	}

	if resolved != docFile {
		t.Fatalf("expected resolved path %s, got: %s", docFile, resolved)
	}
}

func TestDetectPowerShellTarget_ExecutableScript(t *testing.T) {
	tempDir := t.TempDir()
	psScript := filepath.Join(tempDir, "run.ps1")
	_ = os.WriteFile(psScript, []byte("Write-Host 'hello'"), 0644)

	kind, resolved := detectPowerShellTarget(psScript)
	if kind != TargetKindValidCommand {
		t.Fatalf("expected TargetKindValidCommand, got: %s", kind)
	}

	if resolved != psScript {
		t.Fatalf("expected resolved path %s, got: %s", psScript, resolved)
	}
}

func TestDetectPowerShellTarget_ValidCommand(t *testing.T) {
	kind, _ := detectPowerShellTarget("Get-Process")
	if kind != TargetKindValidCommand {
		t.Fatalf("expected TargetKindValidCommand, got: %s", kind)
	}
}

func TestDetectPowerShellTarget_Flags(t *testing.T) {
	kind, _ := detectPowerShellTarget("-NoProfile")
	if kind != TargetKindValidCommand {
		t.Fatalf("expected TargetKindValidCommand for flags, got: %s", kind)
	}
}

func TestIsDocumentFileExtension(t *testing.T) {
	docExts := []string{".md", ".json", ".txt", ".yml", ".yaml", ".sql", "html"}
	for _, ext := range docExts {
		if !isDocumentFileExtension(ext) {
			t.Fatalf("expected true for doc ext: %s", ext)
		}
	}
}

func TestIsExecutableScriptExtension(t *testing.T) {
	execExts := []string{".ps1", ".psm1", ".bat", ".cmd", ".exe", ".sh"}
	for _, ext := range execExts {
		if !isExecutableScriptExtension(ext) {
			t.Fatalf("expected true for exec ext: %s", ext)
		}
	}
}

func TestBuildPowerShellDiagnosticCard_ContainsInquiriesAndLearning(t *testing.T) {
	opts := PowerShellAdviceOptions{
		RawTarget:  ".ai-memory/plans",
		TargetKind: TargetKindDirectory,
	}
	card := buildPowerShellDiagnosticCard(opts)

	hasInquiry := strings.Contains(card, "Why did you run this? What was your expectation?")
	hasLearning := strings.Contains(card, "3-Step Immediate Learning Action")
	hasGuardTitle := strings.Contains(card, "GITMAP AI GUARD")
	if !hasInquiry || !hasLearning || !hasGuardTitle {
		t.Fatalf("expected card to contain inquiry, learning, and title, got:\n%s", card)
	}
}

func TestBuildPowerShellDiagnosticCard_DocumentReplacements(t *testing.T) {
	opts := PowerShellAdviceOptions{
		RawTarget:  "readme.md",
		TargetKind: TargetKindDocument,
	}
	card := buildPowerShellDiagnosticCard(opts)

	hasHeatmap := strings.Contains(card, "gitmap heatmap readme.md")
	hasCat := strings.Contains(card, "gitmap cat readme.md")
	if !hasHeatmap || !hasCat {
		t.Fatalf("expected card to recommend heatmap and cat, got:\n%s", card)
	}
}

func TestBuildPowerShellDiagnosticCard_UnknownCmdlet(t *testing.T) {
	opts := PowerShellAdviceOptions{
		RawTarget:  "foobar-badcmdlet",
		TargetKind: TargetKindUnknownCmdlet,
	}
	card := buildPowerShellDiagnosticCard(opts)

	hasHelp := strings.Contains(card, "gitmap help")
	hasUnknown := strings.Contains(card, "UNKNOWN_CMDLET")
	if !hasHelp || !hasUnknown {
		t.Fatalf("expected card to recommend help and show UNKNOWN_CMDLET, got:\n%s", card)
	}
}

func TestRenderPowerShellDiagnosticCard_OutputStderr(t *testing.T) {
	opts := PowerShellAdviceOptions{RawTarget: "some-dir", TargetKind: TargetKindDirectory}
	out := captureStderr(t, func() {
		renderPowerShellDiagnosticCard(opts)
	})

	hasTarget := strings.Contains(out, "some-dir")
	hasCard := strings.Contains(out, "GITMAP AI GUARD")
	if !hasTarget || !hasCard {
		t.Fatalf("expected stderr output to contain target and card, got:\n%s", out)
	}
}

func TestBuildPreflightAbortError_Details(t *testing.T) {
	err := buildPreflightAbortError(".ai-memory/plans", TargetKindDirectory)
	appErr, isApp := err.(*apperror.AppError)
	if !isApp || appErr == nil {
		t.Fatalf("expected *apperror.AppError, got: %T", err)
	}

	isAbort := appErr.Type == apperror.ErrorTypeAbort
	isReported := appErr.Ctx["reported"] == true
	isExitCodeUsage := appErr.Ctx["exitCode"] == int(cliexit.ExitCodeUsageError)
	if !isAbort || !isReported || !isExitCodeUsage {
		t.Fatalf("unexpected appErr details: %+v", appErr)
	}
}
