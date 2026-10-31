package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestBuildPowerShellCmd_Basic(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"Write-Host", "hello"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	hasScript := strings.Contains(joined, "Write-Host hello")
	if !hasScript {
		t.Fatalf("expected command to contain script, got: %s", joined)
	}
}

func TestBuildPowerShellCmd_WithFlagC(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"-c", "$a = 123; Write-Host $a"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	hasCommand := strings.Contains(joined, "-Command")
	hasScript := strings.Contains(joined, "$a = 123; Write-Host $a")
	if !hasCommand || !hasScript {
		t.Fatalf("expected -Command and script preserved, got: %s", joined)
	}
}

func TestBuildPowerShellCmd_Empty(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", nil)
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	isLengthOne := len(cmd.Args) == 1
	isBinaryPwsh := len(cmd.Args) > 0 && cmd.Args[0] == "pwsh"
	if !isLengthOne || !isBinaryPwsh {
		t.Fatalf("expected single argument pwsh, got: %v", cmd.Args)
	}
}

func TestBuildPowerShellCmd_EmptySlice(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	isLengthOne := len(cmd.Args) == 1
	if !isLengthOne {
		t.Fatalf("expected single argument for empty slice, got: %v", cmd.Args)
	}
}

func TestBuildPowerShellCmd_OnlyForceFlags(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"--force-all", "-f", "force-all"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	isLengthOne := len(cmd.Args) == 1
	if !isLengthOne {
		t.Fatalf("expected single argument when all args are force flags, got: %v", cmd.Args)
	}
}

func TestBuildPowerShellCmd_StripForceInArgs(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"--force-all", "Get-Process"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	hasForce := strings.Contains(joined, "--force-all")
	hasCmd := strings.Contains(joined, "Get-Process")
	if hasForce || !hasCmd {
		t.Fatalf("expected --force-all stripped and Get-Process present, got: %s", joined)
	}
}

func TestBuildPowerShellCmd_WithFlagCommand(t *testing.T) {
	cmd := buildPowerShellCmd("pwsh", []string{"-Command", "Get-Location"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	joined := strings.Join(cmd.Args, " ")
	hasCommand := strings.Contains(joined, "-Command")
	hasCmdlet := strings.Contains(joined, "Get-Location")
	if !hasCommand || !hasCmdlet {
		t.Fatalf("expected -Command and Get-Location in args, got: %s", joined)
	}
}

func TestAssemblePowerShellCmd(t *testing.T) {
	cmd := assemblePowerShellCmd("pwsh", []string{"Write-Host", "test"})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	isBinaryMatch := cmd.Path == "pwsh" || cmd.Args[0] == "pwsh"
	hasCommand := strings.Contains(strings.Join(cmd.Args, " "), "-Command")
	if !isBinaryMatch || !hasCommand {
		t.Fatalf("expected pwsh and -Command in args: %v", cmd.Args)
	}
}

func TestAssemblePowerShellCmd_WindowsProfile(t *testing.T) {
	isWin := constants.IsWindowsPlatform()
	if !isWin {
		t.Skip("skipping Windows profile check on non-Windows")
	}

	cmd := assemblePowerShellCmd("pwsh", []string{"Write-Host", "test"})
	hasNoProfile := strings.Contains(strings.Join(cmd.Args, " "), "-NoProfile")
	if !hasNoProfile {
		t.Fatalf("expected -NoProfile on Windows, got: %v", cmd.Args)
	}
}

func TestExtractPowerShellScript_FlagCommand(t *testing.T) {
	script := extractPowerShellScript([]string{"-Command", "Get-Date", "-Format", "yyyy"})
	if script != "Get-Date -Format yyyy" {
		t.Fatalf("unexpected extracted script: %s", script)
	}
}

func TestExtractPowerShellScript_FlagC(t *testing.T) {
	script := extractPowerShellScript([]string{"-c", "$x", "=", "10"})
	if script != "$x = 10" {
		t.Fatalf("unexpected extracted script: %s", script)
	}
}

func TestExtractPowerShellScript_DirectArgs(t *testing.T) {
	script := extractPowerShellScript([]string{"Get-ChildItem", "-Path", "."})
	if script != "Get-ChildItem -Path ." {
		t.Fatalf("unexpected extracted script: %s", script)
	}
}

func TestStripForceFlag_Variations(t *testing.T) {
	input := []string{"--force-all", "first", "-f", "second", "force-all", "third"}
	cleaned := stripForceFlag(input)
	joined := strings.Join(cleaned, " ")
	if joined != "first second third" {
		t.Fatalf("expected 'first second third', got: %s", joined)
	}
}

func TestStripForceFlag_NoFlags(t *testing.T) {
	input := []string{"echo", "clean", "tokens"}
	cleaned := stripForceFlag(input)
	joined := strings.Join(cleaned, " ")
	if joined != "echo clean tokens" {
		t.Fatalf("expected 'echo clean tokens', got: %s", joined)
	}
}

func TestStripForceFlag_NilAndEmpty(t *testing.T) {
	cleanNil := stripForceFlag(nil)
	if len(cleanNil) != 0 {
		t.Fatalf("expected 0 items for nil, got: %d", len(cleanNil))
	}

	cleanEmpty := stripForceFlag([]string{})
	if len(cleanEmpty) != 0 {
		t.Fatalf("expected 0 items for empty, got: %d", len(cleanEmpty))
	}
}

func TestFindPowerShellExecutable_Contract(t *testing.T) {
	psPath, hasPS := findPowerShellExecutable()
	isWin := constants.IsWindowsPlatform()
	if isWin && !hasPS {
		t.Fatal("expected powershell executable on Windows platform")
	}

	if hasPS && psPath == "" {
		t.Fatal("expected non-empty path when hasPS is true")
	}
}

func TestPrintPowerShellMissingAdvice_Output(t *testing.T) {
	out := captureStdoutForTest(t, printPowerShellMissingAdvice)
	hasMissingMsg := strings.Contains(out, "PowerShell (pwsh) is not installed")
	if !hasMissingMsg {
		t.Errorf("expected missing message in advice, got: %s", out)
	}

	hasInstallCmd := strings.Contains(out, "gitmap install powershell")
	if !hasInstallCmd {
		t.Errorf("expected gitmap install powershell in advice, got: %s", out)
	}
}

func TestPrintPowerShellMissingHeader_Output(t *testing.T) {
	out := captureStdoutForTest(t, printPowerShellMissingHeader)
	hasHeader := strings.Contains(out, "PowerShell (pwsh) is not installed on this system")
	if !hasHeader {
		t.Errorf("expected header text, got: %s", out)
	}
}

func TestPrintPowerShellInstallOptions_Output(t *testing.T) {
	out := captureStdoutForTest(t, printPowerShellInstallOptions)
	hasLocal := strings.Contains(out, "Local Installation:")
	hasWinget := strings.Contains(out, "winget install")
	hasSSH := strings.Contains(out, "Remote SSH / Unix Installation:")
	if !hasLocal || !hasWinget || !hasSSH {
		t.Errorf("expected local, winget, and remote SSH in options, got: %s", out)
	}
}

func TestRunPowerShell_ValidExecution(t *testing.T) {
	_, hasPS := findPowerShellExecutable()
	if !hasPS {
		t.Skip("skipping execution test: powershell not installed")
	}

	err := RunPowerShell([]string{"-Command", "Write-Output 'gitmap-test-ok'"})
	if err != nil {
		t.Fatalf("expected successful execution, got: %v", err)
	}
}

func TestRunPowerShell_ExitError(t *testing.T) {
	_, hasPS := findPowerShellExecutable()
	if !hasPS {
		t.Skip("skipping execution test: powershell not installed")
	}

	err := RunPowerShell([]string{"-Command", "exit 42"})
	if err == nil {
		t.Fatal("expected error on non-zero exit code")
	}
}

func TestPowerShell_PreflightDirectory(t *testing.T) {
	tempDir := t.TempDir()
	err := RunPowerShell([]string{tempDir})
	if err == nil {
		t.Fatal("expected error on directory target")
	}

	assertAbortUsageError(t, err)
}

func assertAbortUsageError(t *testing.T, err error) {
	appErr, isApp := err.(*apperror.AppError)
	if !isApp || appErr == nil {
		t.Fatalf("expected *apperror.AppError, got: %T", err)
	}

	if appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("expected ErrorTypeAbort, got: %s", appErr.Type)
	}

	code, hasCode := appErr.Ctx["exitCode"]
	if !hasCode || code != int(cliexit.ExitCodeUsageError) {
		t.Errorf("expected exitCode=%d, got: %v", cliexit.ExitCodeUsageError, code)
	}
}

func TestPowerShell_PreflightDocumentFile(t *testing.T) {
	testDoc := createSampleDocument(t, "sample.md")
	err := RunPowerShell([]string{testDoc})
	if err == nil {
		t.Fatal("expected error on document target")
	}

	assertAbortDocumentError(t, err)
}

func createSampleDocument(t *testing.T, name string) string {
	docPath := filepath.Join(t.TempDir(), name)
	if writeErr := os.WriteFile(docPath, []byte("# Test"), 0644); writeErr != nil {
		t.Fatalf("failed to create sample doc: %v", writeErr)
	}

	return docPath
}

func assertAbortDocumentError(t *testing.T, err error) {
	appErr, isApp := err.(*apperror.AppError)
	if !isApp || appErr == nil {
		t.Fatalf("expected *apperror.AppError, got: %T", err)
	}

	if appErr.Type != apperror.ErrorTypeAbort {
		t.Errorf("expected ErrorTypeAbort, got: %s", appErr.Type)
	}

	kindVal, hasKind := appErr.Ctx["kind"]
	if !hasKind || kindVal != string(TargetKindDocument) {
		t.Errorf("expected kind=%s, got: %v", TargetKindDocument, kindVal)
	}
}

func TestPowerShell_AllowValidPs1Script(t *testing.T) {
	skipIfNoPowerShell(t)
	scriptPath := createSampleDocument(t, "test.ps1")
	preflightErr := inspectPowerShellPreflight([]string{scriptPath})
	if preflightErr != nil {
		t.Fatalf("expected preflight to allow .ps1, got: %v", preflightErr)
	}
}

func skipIfNoPowerShell(t *testing.T) {
	_, hasPS := findPowerShellExecutable()
	if !hasPS {
		t.Skip("skipping test: powershell not installed")
	}
}

func TestPowerShell_AllowPowerShellExpression(t *testing.T) {
	expr := "Get-ChildItem -Path '.'"
	preflightErr := inspectPowerShellPreflight([]string{expr})
	if preflightErr != nil {
		t.Fatalf("expected preflight to allow expression, got: %v", preflightErr)
	}
}

func TestPowerShell_CleanErrorBubbling(t *testing.T) {
	tempDir := t.TempDir()
	err := RunPowerShell([]string{tempDir})
	if err == nil {
		t.Fatal("expected preflight error")
	}

	assertCleanErrorBubblingReport(t, err)
}

func assertCleanErrorBubblingReport(t *testing.T, err error) {
	exitCode := resolveAbortExitCode(err)
	if exitCode != int(cliexit.ExitCodeUsageError) {
		t.Errorf("expected exitCode=%d, got: %d", cliexit.ExitCodeUsageError, exitCode)
	}

	appErr := err.(*apperror.AppError)
	var buf bytes.Buffer
	cliexit.WriteAppErrorReport(&buf, appErr)
	hasStack := strings.Contains(buf.String(), "stack trace:")
	if hasStack {
		t.Errorf("expected stack trace suppressed, got report:\n%s", buf.String())
	}
}

func TestPowerShell_UnrecognizedCmdletInterception(t *testing.T) {
	skipIfNoPowerShell(t)
	err := RunPowerShell([]string{"nonexistent_cmdlet_xyz_12345"})
	if err == nil {
		t.Fatal("expected error on unrecognized cmdlet")
	}

	assertAbortUsageError(t, err)
}
