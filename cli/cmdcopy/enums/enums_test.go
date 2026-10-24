package enums

import (
	"testing"
)

func TestCopyActionType(t *testing.T) {
	if !CopyActionCopy.IsValid() || !CopyActionPaste.IsValid() {
		t.Fatal("expected CopyActionCopy and CopyActionPaste to be valid")
	}

	invalid := CopyActionType("unknown")
	if invalid.IsValid() {
		t.Fatal("expected unknown action to be invalid")
	}

	if CopyActionCopy.String() != "copy" {
		t.Fatalf("expected 'copy', got %q", CopyActionCopy.String())
	}

	parsed, err := ParseCopyActionType("COPY")
	if err != nil || !parsed.IsCopy() {
		t.Fatalf("expected CopyActionCopy, got %v, err: %v", parsed, err)
	}

	parsedPaste, err := ParseCopyActionType("paste")
	if err != nil || !parsedPaste.IsPaste() {
		t.Fatalf("expected CopyActionPaste, got %v, err: %v", parsedPaste, err)
	}

	if _, err := ParseCopyActionType("invalid"); err == nil {
		t.Fatal("expected error parsing invalid copy action")
	}

	if len(AllCopyActionTypes()) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(AllCopyActionTypes()))
	}

	if !CopyActionCopy.HasAction() {
		t.Fatal("expected HasAction to return true for CopyActionCopy")
	}
}

func TestCopyTargetType(t *testing.T) {
	targets := []CopyTargetType{
		CopyTargetClipboard,
		CopyTargetMemory,
		CopyTargetStdout,
		CopyTargetFile,
	}

	for _, target := range targets {
		if !target.IsValid() {
			t.Fatalf("expected %v to be valid", target)
		}
	}

	invalid := CopyTargetType("network")
	if invalid.IsValid() {
		t.Fatal("expected invalid target to not be valid")
	}

	if !CopyTargetClipboard.IsClipboard() {
		t.Fatal("expected IsClipboard to be true")
	}

	if !CopyTargetMemory.IsMemory() {
		t.Fatal("expected IsMemory to be true")
	}

	if !CopyTargetStdout.IsStdout() {
		t.Fatal("expected IsStdout to be true")
	}

	if !CopyTargetFile.IsFile() {
		t.Fatal("expected IsFile to be true")
	}

	parsed, err := ParseCopyTargetType("STDOUT")
	if err != nil || !parsed.IsStdout() {
		t.Fatalf("expected stdout target, got %v, err: %v", parsed, err)
	}

	if _, err := ParseCopyTargetType("bad"); err == nil {
		t.Fatal("expected error for bad target")
	}

	if len(AllCopyTargetTypes()) != 4 {
		t.Fatalf("expected 4 targets, got %d", len(AllCopyTargetTypes()))
	}

	if !CopyTargetFile.HasTarget() {
		t.Fatal("expected HasTarget to be true")
	}
}

func TestCopyModeType(t *testing.T) {
	modes := []CopyModeType{
		CopyModeText,
		CopyModeFile,
		CopyModePipe,
	}

	for _, mode := range modes {
		if !mode.IsValid() {
			t.Fatalf("expected %v to be valid", mode)
		}
	}

	invalid := CopyModeType("binary")
	if invalid.IsValid() {
		t.Fatal("expected binary mode to be invalid")
	}

	if !CopyModeText.IsText() {
		t.Fatal("expected IsText to be true")
	}

	if !CopyModeFile.IsFile() {
		t.Fatal("expected IsFile to be true")
	}

	if !CopyModePipe.IsPipe() {
		t.Fatal("expected IsPipe to be true")
	}

	parsed, err := ParseCopyModeType("PIPE")
	if err != nil || !parsed.IsPipe() {
		t.Fatalf("expected pipe mode, got %v, err: %v", parsed, err)
	}

	if _, err := ParseCopyModeType("invalid"); err == nil {
		t.Fatal("expected error for invalid mode")
	}

	if len(AllCopyModeTypes()) != 3 {
		t.Fatalf("expected 3 modes, got %d", len(AllCopyModeTypes()))
	}

	if !CopyModeText.HasMode() {
		t.Fatal("expected HasMode to be true")
	}
}

func TestCopySourceType(t *testing.T) {
	sources := []CopySourceType{
		CopySourceArgs,
		CopySourceFile,
		CopySourceStdin,
		CopySourceClipboard,
		CopySourceMemory,
	}

	for _, source := range sources {
		if !source.IsValid() {
			t.Fatalf("expected %v to be valid", source)
		}
	}

	invalid := CopySourceType("unknown")
	if invalid.IsValid() {
		t.Fatal("expected invalid source to be invalid")
	}

	if !CopySourceArgs.IsArgs() {
		t.Fatal("expected IsArgs to be true")
	}

	if !CopySourceFile.IsFile() {
		t.Fatal("expected IsFile to be true")
	}

	if !CopySourceStdin.IsStdin() {
		t.Fatal("expected IsStdin to be true")
	}

	if !CopySourceClipboard.IsClipboard() {
		t.Fatal("expected IsClipboard to be true")
	}

	if !CopySourceMemory.IsMemory() {
		t.Fatal("expected IsMemory to be true")
	}

	parsed, err := ParseCopySourceType("MEMORY")
	if err != nil || !parsed.IsMemory() {
		t.Fatalf("expected memory source, got %v, err: %v", parsed, err)
	}

	if _, err := ParseCopySourceType("bad"); err == nil {
		t.Fatal("expected error for bad source")
	}

	if len(AllCopySourceTypes()) != 5 {
		t.Fatalf("expected 5 sources, got %d", len(AllCopySourceTypes()))
	}

	if !CopySourceArgs.HasSource() {
		t.Fatal("expected HasSource to be true")
	}
}
