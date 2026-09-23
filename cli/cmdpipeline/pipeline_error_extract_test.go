package cmdpipeline

import (
	"strings"
	"testing"
)

func TestExtractTauriBundlerAndCompilerWarnings(t *testing.T) {
	rawLogs := strings.Join([]string{
		"build\tBuild Desktop App\t2026-03-23T10:00:00.0000000Z Browserslist: browsers data (caniuse-lite) is 8 months old. Please run: npx update-browserslist-db",
		"build\tBuild Desktop App\t2026-03-23T10:00:01.0000000Z vite v7.3.1 building client environment for production...",
		"build\tBuild Desktop App\t2026-03-23T10:00:02.0000000Z transforming...",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z warning: unused import: `std::path::Path`",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z   --> src-tauri/src/main.rs:12:5",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z 12 | use std::path::Path;",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |     ^^^^^^^^^^^^^^^",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    = note: `#[warn(unused_imports)]` on by default",
		"build\tBuild Desktop App\t2026-03-23T10:00:06.0000000Z Info Looking up installed tauri packages...",
		"build\tBuild Desktop App\t2026-03-23T10:00:07.0000000Z failed to bundle project: Failed to copy binary from \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/agm\" to \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/bundle/macos/agm.app/Contents/MacOS/agm\": \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/agm\" does not exist",
		"build\tBuild Desktop App\t2026-03-23T10:00:08.0000000Z ##[error]Process completed with exit code 1.",
	}, "\n")

	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		t.Fatalf("expected at least 1 failed job item, got 0")
	}

	job := jobs[0]
	if !strings.Contains(job.FailureSummary, "failed to bundle project") {
		t.Errorf("expected FailureSummary to contain 'failed to bundle project', got: %s", job.FailureSummary)
	}
	if !strings.Contains(job.FailureSummary, "does not exist") {
		t.Errorf("expected FailureSummary to contain 'does not exist', got: %s", job.FailureSummary)
	}
	if strings.Contains(job.FailureSummary, "Process completed with exit code 1") {
		t.Errorf("FailureSummary should NOT be generic exit code when bundler error exists, got: %s", job.FailureSummary)
	}

	hasWarningLinePointer := false
	hasWarningCodeSnippet := false
	for _, w := range job.Warnings {
		if strings.Contains(w, "--> src-tauri/src/main.rs:12:5") {
			hasWarningLinePointer = true
		}
		if strings.Contains(w, "use std::path::Path;") {
			hasWarningCodeSnippet = true
		}
	}

	if !hasWarningLinePointer {
		t.Errorf("expected job.Warnings to capture line pointer '--> src-tauri/src/main.rs:12:5', got: %v", job.Warnings)
	}
	if !hasWarningCodeSnippet {
		t.Errorf("expected job.Warnings to capture code snippet 'use std::path::Path;', got: %v", job.Warnings)
	}
}

func TestIsStrongerSummaryBundlerPriority(t *testing.T) {
	bundlerErr := `failed to bundle project: Failed to copy binary: file does not exist`
	genericExit := `Error: Process completed with exit code 1.`

	if !isStrongerSummary(bundlerErr, genericExit) {
		t.Errorf("bundler error should be stronger summary than generic exit code")
	}

	if isStrongerSummary(genericExit, bundlerErr) {
		t.Errorf("generic exit code should NOT override existing bundler error summary")
	}
}

func TestIsToolProgressNoiseTauri(t *testing.T) {
	noiseLines := []string{
		"Browserslist: browsers data (caniuse-lite) is 8 months old. Please run: npx update-browserslist-db",
		"transforming...",
		"vite v7.3.1 building client environment for production...",
		"Info Looking up installed tauri packages to check mismatched versions...",
		"Running beforeBuildCommand `npm run build`",
	}

	for _, line := range noiseLines {
		if !isToolProgressNoise(line) {
			t.Errorf("expected line to be recognized as tool progress noise: %s", line)
		}
	}
}
