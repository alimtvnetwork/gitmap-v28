package macro

import (
	"context"
	"testing"
)

func TestMacroSaveLoadListDelete(t *testing.T) {
	m := &Macro{
		Name: "test-macro",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo hello", TimeoutSeconds: 10},
			{StepNum: 2, CommandLine: "echo world", TimeoutSeconds: 10},
		},
	}

	if err := SaveMacro(m); err != nil {
		t.Fatalf("SaveMacro failed: %v", err)
	}

	loaded, err := LoadMacro("test-macro")
	if err != nil {
		t.Fatalf("LoadMacro failed: %v", err)
	}
	if loaded.Name != "test-macro" || len(loaded.Steps) != 2 {
		t.Fatalf("Loaded macro mismatch: %+v", loaded)
	}

	list, err := ListMacros()
	if err != nil {
		t.Fatalf("ListMacros failed: %v", err)
	}
	found := false
	for _, item := range list {
		if item.Name == "test-macro" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("test-macro not found in list: %+v", list)
	}

	// Dry run execute
	err = Execute(context.Background(), loaded, ExecOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Dry run execution failed: %v", err)
	}

	if err := DeleteMacro("test-macro"); err != nil {
		t.Fatalf("DeleteMacro failed: %v", err)
	}
}
