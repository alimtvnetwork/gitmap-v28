package cmdpipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPEFormatProfileDefaults(t *testing.T) {
	def := DefaultPEFormatProfile()
	if def.Name != "default" || !def.IsSkipEmpty {
		t.Fatalf("unexpected default profile: %+v", def)
	}
	if len(def.StripPrefixes) == 0 || len(def.ErrorMarkers) == 0 {
		t.Fatalf("expected non-empty default prefixes and error markers")
	}

	tauri := BuiltinTauriProfile()
	if tauri.Name != "tauri" || tauri.Alias != "tauri" {
		t.Fatalf("unexpected tauri profile: %+v", tauri)
	}
	if len(tauri.StripPrefixes) == 0 || len(tauri.CaptureUntilMarkers) == 0 {
		t.Fatalf("expected tauri strip prefixes and until markers")
	}
}

func TestPEFormatStoreSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	origFormatsDir := resolveFormatsDir()
	_ = origFormatsDir

	p := PEFormatProfile{
		Name:          "testcustom",
		Alias:         "tc",
		Description:   "Test custom profile",
		StripPrefixes: []string{"NoisePrefix:"},
		IsSkipEmpty:   true,
		ErrorMarkers:  []string{"FATAL_TEST_ERROR"},
	}

	testJsonPath := filepath.Join(tmpDir, "custom.json")
	if err := os.WriteFile(testJsonPath, []byte(`{
		"name": "fromfile",
		"alias": "ff",
		"strip_prefixes": ["file_noise:"],
		"error_markers": ["FILE_ERR"]
	}`), 0644); err != nil {
		t.Fatalf("failed to write test json: %v", err)
	}

	loadedFromFile, err := LoadPEFormatProfile(testJsonPath)
	if err != nil {
		t.Fatalf("failed to load from file: %v", err)
	}
	if loadedFromFile.Name != "fromfile" || loadedFromFile.Alias != "ff" {
		t.Fatalf("expected name 'fromfile', got '%s'", loadedFromFile.Name)
	}

	builtinTauri, err := LoadPEFormatProfile("tauri")
	if err != nil || builtinTauri.Name != "tauri" {
		t.Fatalf("failed to load builtin tauri profile: %v", err)
	}

	builtinRust, err := LoadPEFormatProfile("rust")
	if err != nil || builtinRust.Name != "tauri" {
		t.Fatalf("failed to load builtin rust profile: %v", err)
	}

	_ = p
}

func TestPEFormatFilterEvaluation(t *testing.T) {
	tauri := BuiltinTauriProfile()

	rawLog := strings.Join([]string{
		"Browserslist: browsers data is 8 months old. Please run: npx update-browserslist-db",
		"transforming...",
		"vite v7.3.1 building client environment for production...",
		"warning: unused import: `std::path::Path`",
		"  --> src-tauri/src/main.rs:12:5",
		"   |",
		"12 | use std::path::Path;",
		"   |     ^^^^^^^^^^^^^^^",
		"   |",
		"   = note: `#[warn(unused_imports)]` on by default",
		"Compiling agm v0.1.0 (/Users/runner/work/agm-desktop/agm-desktop/src-tauri)",
		"Finished `release` profile [optimized] target(s) in 2m 14s",
		`failed to bundle project: Failed to copy binary from "/target/release/agm" to "/bundle/macos/agm.app/Contents/MacOS/agm": "/target/release/agm" does not exist`,
		"Error: Process completed with exit code 1.",
	}, "\n")

	filtered := FilterLogWithProfile(rawLog, &tauri)

	if strings.Contains(filtered, "Browserslist:") {
		t.Errorf("expected Browserslist to be filtered out")
	}
	if strings.Contains(filtered, "transforming...") {
		t.Errorf("expected transforming... to be filtered out")
	}
	if !strings.Contains(filtered, "warning: unused import:") {
		t.Errorf("expected warning header to be retained")
	}
	if !strings.Contains(filtered, "--> src-tauri/src/main.rs:12:5") {
		t.Errorf("expected warning line pointer to be retained")
	}
	if !strings.Contains(filtered, "use std::path::Path;") {
		t.Errorf("expected code line to be retained")
	}
	if !strings.Contains(filtered, "= note:") {
		t.Errorf("expected note to be retained")
	}
	if !strings.Contains(filtered, "failed to bundle project:") {
		t.Errorf("expected bundler failure to be retained")
	}
	if !strings.Contains(filtered, "does not exist") {
		t.Errorf("expected does not exist to be retained")
	}
}

func TestParsePipelineErrorFlagsFormat(t *testing.T) {
	f1 := ParsePipelineErrorFlags([]string{"-f", "tauri"})
	if f1.FormatProfile != "tauri" {
		t.Fatalf("expected FormatProfile 'tauri', got '%s'", f1.FormatProfile)
	}
	if f1.HasFix {
		t.Fatalf("expected HasFix to be false when -f is followed by format name")
	}

	f2 := ParsePipelineErrorFlags([]string{"--format", "custom.json"})
	if f2.FormatProfile != "custom.json" {
		t.Fatalf("expected FormatProfile 'custom.json', got '%s'", f2.FormatProfile)
	}

	f3 := ParsePipelineErrorFlags([]string{"-f"})
	if f3.FormatProfile != "" {
		t.Fatalf("expected empty FormatProfile when -f is standalone, got '%s'", f3.FormatProfile)
	}
	if !f3.HasFix {
		t.Fatalf("expected HasFix to be true when -f is standalone")
	}

	f4 := ParsePipelineErrorFlags([]string{"--fix"})
	if !f4.HasFix {
		t.Fatalf("expected HasFix to be true for --fix")
	}

	f5 := ParsePipelineErrorFlags([]string{"-f", "-c"})
	if f5.FormatProfile != "" {
		t.Fatalf("expected empty FormatProfile when followed by another flag, got '%s'", f5.FormatProfile)
	}
	if !f5.HasFix || !f5.HasCheck {
		t.Fatalf("expected both HasFix and HasCheck to be true")
	}
}

func TestHandlePEFormatCommandsRouting(t *testing.T) {
	handled, err := HandlePEFormatCommands([]string{"list-formats"})
	if !handled || err != nil {
		t.Fatalf("expected list-formats to be handled successfully: %v", err)
	}

	handled, err = HandlePEFormatCommands([]string{"preview-format", "tauri"})
	if !handled || err != nil {
		t.Fatalf("expected preview-format to be handled successfully: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "test.log")
	_ = os.WriteFile(tmpFile, []byte("warning: sample\nfailed to bundle: err\n"), 0644)

	handled, err = HandlePEFormatCommands([]string{"-f", "tauri", "-test", tmpFile})
	if !handled || err != nil {
		t.Fatalf("expected format test on file to be handled: %v", err)
	}

	handled, _ = HandlePEFormatCommands([]string{"ee4a694"})
	if handled {
		t.Fatalf("expected commit SHA not to be treated as format command")
	}
}
