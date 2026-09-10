// Package cmd — install_custom_exec_test.go: unit tests for custom installer execution.
package cmd

import (
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func TestParseInstructionsMap(t *testing.T) {
	rawJSON := `{"win":"echo win","unix":"echo unix","ubuntu":"echo ubuntu"}`
	parsed := parseInstructionsMap(rawJSON)
	if parsed["win"] != "echo win" || parsed["unix"] != "echo unix" || parsed["ubuntu"] != "echo ubuntu" {
		t.Fatalf("unexpected parsed instructions: %+v", parsed)
	}

	plainText := "echo single"
	parsedPlain := parseInstructionsMap(plainText)
	if parsedPlain["all"] != "echo single" {
		t.Fatalf("expected all key for plain text, got %+v", parsedPlain)
	}
}

func TestResolveOSScriptForPlatform(t *testing.T) {
	script := &model.InstallerScript{
		Name:         "test-tool",
		Slug:         "test-tool",
		TargetOS:     "all",
		Instructions: `{"win":"Write-Output 'windows'","unix":"echo 'unix'","ubuntu":"echo 'ubuntu'"}`,
	}

	text, runner, err := resolveOSScriptForPlatform(script)
	if err != nil {
		t.Fatalf("resolveOSScriptForPlatform failed: %v", err)
	}

	if runtime.GOOS == "windows" {
		assertWindowsResolution(t, text, runner)
	}
}

func assertWindowsResolution(t *testing.T, text string, runner string) {
	if runner != "powershell" || text != "Write-Output 'windows'" {
		t.Fatalf("unexpected windows resolution: runner=%s, text=%s", runner, text)
	}
}

func TestResolveOSScriptForPlatform_UnsupportedOS(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping Windows-specific unsupported check on non-Windows host")
	}

	script := &model.InstallerScript{
		Name:         "unix-only",
		Slug:         "unix-only",
		TargetOS:     "unix",
		Instructions: `{"unix":"echo 'unix'"}`,
	}
	_, _, err := resolveOSScriptForPlatform(script)
	if err == nil {
		t.Fatal("expected error on windows for unix-only script")
	}
}

func TestExecuteCustomInstaller_DryRun(t *testing.T) {
	script := &model.InstallerScript{
		Name:         "dry-tool",
		Slug:         "dry-tool",
		TargetOS:     "all",
		Instructions: `{"win":"Write-Output dry","unix":"echo dry"}`,
	}

	opts := installOptions{
		Tool:   "dry-tool",
		DryRun: true,
	}

	if err := executeCustomInstaller(script, opts); err != nil {
		t.Fatalf("expected dry-run to succeed, got %v", err)
	}
}

func TestRunInstall_CustomInstallerIntegration(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Fatalf("open store failed: %v", errDB)
	}
	defer db.Close()
	if errMig := db.MigrateInstallers(); errMig != nil {
		t.Fatalf("migrate installers failed: %v", errMig)
	}

	script := &model.InstallerScript{
		Name:         "custom-runner",
		Slug:         "custom-runner",
		Description:  "Sample custom tool",
		TargetOS:     "all",
		Version:      "v1.0.0",
		Instructions: `{"win":"Write-Output 'hello from custom'","unix":"echo 'hello from custom'"}`,
	}
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		t.Fatalf("create installer failed: %v", errCreate)
	}

	assertCustomInstallLookupAndDryRun(t)
}

func assertCustomInstallLookupAndDryRun(t *testing.T) {
	if !isKnownInstallTool("custom-runner") {
		t.Fatal("expected isKnownInstallTool to return true for custom-runner")
	}

	if err := runInstall([]string{"custom-runner", "--dry-run"}); err != nil {
		t.Fatalf("runInstall custom-runner failed: %v", err)
	}
}

func TestExecuteCustomInstaller_LiveExecution(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	script := &model.InstallerScript{
		Name:         "live-echo",
		Slug:         "live-echo",
		TargetOS:     "all",
		Version:      "v1.0.0",
		Instructions: `{"win":"Write-Output 'live win'","unix":"echo 'live unix'"}`,
	}

	if err := executeCustomInstaller(script, installOptions{Tool: "live-echo"}); err != nil {
		t.Fatalf("expected live execution to succeed, got %v", err)
	}

	assertRecordedInSplitDB(t, "live-echo")
}

func assertRecordedInSplitDB(t *testing.T, slug string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		t.Fatalf("open split db failed: %v", err)
	}
	defer splitDB.Close()

	tool, errGet := splitDB.GetInstalledTool(slug)
	if errGet != nil || tool.Tool != slug {
		t.Fatalf("expected tool %q to be recorded in installation.db, got err: %v", slug, errGet)
	}
}
