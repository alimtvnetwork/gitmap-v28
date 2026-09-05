package examples_test

import (
	"os"
	"path/filepath"
	"testing"

	"coding-guidelines/common/examples"
	"coding-guidelines/common/pkg/errtype"
)

func TestRunBytesLifecycleExample(t *testing.T) {
	env, appErr := examples.RunBytesLifecycleExample()
	if appErr != nil {
		t.Fatalf("bytes lifecycle failed: %v", appErr)
	}

	if !env.IsSuccess() {
		t.Fatalf("expected successful env")
	}

	if env.Payload().Name != "Alice Smith" {
		t.Fatalf("unexpected payload name: %s", env.Payload().Name)
	}
}

func TestRunJsonResultExample(t *testing.T) {
	pretty, appErr := examples.RunJsonResultExample()
	if appErr != nil {
		t.Fatalf("json result example failed: %v", appErr)
	}

	if len(pretty) == 0 {
		t.Fatalf("expected non-empty pretty string")
	}
}

func TestRunReflectConverterExample(t *testing.T) {
	profile, appErr := examples.RunReflectConverterExample()
	if appErr != nil {
		t.Fatalf("reflect converter example failed: %v", appErr)
	}

	if profile == nil || profile.ID != "cust-300" || profile.Status != errtype.ProcessStateCompleted {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestRunEnumOperationsExample(t *testing.T) {
	appErr := examples.RunEnumOperationsExample()
	if appErr != nil {
		t.Fatalf("enum operations failed: %v", appErr)
	}
}

func TestRunFileWriterAndAppenderExample(t *testing.T) {
	tempDir := t.TempDir()

	appErr := examples.RunFileWriterAndAppenderExample(tempDir)
	if appErr != nil {
		t.Fatalf("file writer and appender example failed: %v", appErr)
	}

	// Verify writer output
	writerFile := filepath.Join(tempDir, "output-strategy.txt")
	content, err := os.ReadFile(writerFile)
	if err != nil || string(content) != "Version 3: Truncated Clean State\n" {
		t.Fatalf("unexpected writer output: %s", string(content))
	}

	// Verify appender output
	appenderFile := filepath.Join(tempDir, "journal", "audit.log")
	journal, err := os.ReadFile(appenderFile)
	expectedJournal := "LOG ENTRY 1: Service started\nLOG ENTRY 2: Config loaded\n"
	if err != nil || string(journal) != expectedJournal {
		t.Fatalf("unexpected journal output: %s", string(journal))
	}
}

func TestRunBoundFileWriterExample(t *testing.T) {
	tempDir := t.TempDir()

	appErr := examples.RunBoundFileWriterExample(tempDir)
	if appErr != nil {
		t.Fatalf("bound file writer example failed: %v", appErr)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "bound-state.txt"))
	if err != nil {
		t.Fatalf("failed to read bound file: %v", err)
	}

	expected := "State: Initialized\nEvent: Connection established\nAudit: Checkpoint recorded\n--- Batch Header ---\nAction: Sync A\nAction: Sync B\n--- Batch Footer ---\nFinal: Terminated\n"
	if string(content) != expected {
		t.Fatalf("unexpected bound file content:\n%s", string(content))
	}
}
