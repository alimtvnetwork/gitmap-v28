package cmdcopy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/enums"
)

func TestFilterFlagArgs(t *testing.T) {
	args := []string{"foo", constants.FlagFile, "sample.txt", "bar"}
	got := filterFlagArgs(args)
	if len(got) != 2 || got[0] != "foo" || got[1] != "bar" {
		t.Fatalf("expected [foo bar], got %v", got)
	}
}

func TestExtractCopyFilePath(t *testing.T) {
	if got := extractCopyFilePath([]string{constants.FlagFile, "test.txt"}); got != "test.txt" {
		t.Fatalf("expected test.txt, got %s", got)
	}

	if got := extractCopyFilePath([]string{constants.FlagFileShort, "test.txt"}); got != "test.txt" {
		t.Fatalf("expected test.txt, got %s", got)
	}

	if got := extractCopyFilePath([]string{constants.FlagPrefixFile + "inline.txt"}); got != "inline.txt" {
		t.Fatalf("expected inline.txt, got %s", got)
	}
}

func TestExtractPasteOutFile(t *testing.T) {
	if got := extractPasteOutFile([]string{constants.FlagFile, "out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}

	if got := extractPasteOutFile([]string{constants.FlagOutShort, "out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}

	if got := extractPasteOutFile([]string{constants.FlagPrefixOut + "out.txt"}); got != "out.txt" {
		t.Fatalf("expected out.txt, got %s", got)
	}
}

func TestDetermineCopySource(t *testing.T) {
	if src := determineCopySource([]string{constants.FlagFile, "foo.txt"}); src != enums.CopySourceFile || !src.IsFile() {
		t.Fatalf("expected file source, got %v", src)
	}

	if src := determineCopySource([]string{"some", "words"}); src != enums.CopySourceArgs || !src.IsArgs() {
		t.Fatalf("expected args source, got %v", src)
	}

	if src := determineCopySource([]string{}); src != enums.CopySourceStdin || !src.IsStdin() {
		t.Fatalf("expected stdin source, got %v", src)
	}
}

func TestDeterminePasteTarget(t *testing.T) {
	if tgt := determinePasteTarget("destination.txt"); tgt != enums.CopyTargetFile || !tgt.IsFile() {
		t.Fatalf("expected file target, got %v", tgt)
	}

	if tgt := determinePasteTarget(""); tgt != enums.CopyTargetStdout || !tgt.IsStdout() {
		t.Fatalf("expected stdout target, got %v", tgt)
	}
}

func TestIsFileFlagWithArg(t *testing.T) {
	for _, flag := range constants.AllCopyFileFlags() {
		if !isFileFlagWithArg(flag) {
			t.Fatalf("expected %q to be recognized file flag", flag)
		}
	}

	if isFileFlagWithArg("--invalid-flag") {
		t.Fatal("expected invalid flag to return false")
	}
}

func TestSaveContentToMemory(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	testText := "gitmap persistent clipboard content test"
	if err := saveContentToMemoryAndClipboard(testText); err != nil {
		t.Fatalf("saveContentToMemoryAndClipboard failed: %v", err)
	}

	memPath, err := resolveMemoryFilePath()
	if err != nil {
		t.Fatalf("resolveMemoryFilePath failed: %v", err)
	}

	data, readErr := os.ReadFile(memPath)
	if readErr != nil || string(data) != testText {
		t.Fatalf("expected %q, got %q (err: %v)", testText, string(data), readErr)
	}
}

func TestWritePastedContentToFile(t *testing.T) {
	tempHome := t.TempDir()
	testText := "gitmap paste output file test"
	outPath := filepath.Join(tempHome, "pasted_output.txt")

	if err := writePastedContentToFile(outPath, testText); err != nil {
		t.Fatalf("writePastedContentToFile failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil || string(data) != testText {
		t.Fatalf("expected %q, got %q (err: %v)", testText, string(data), err)
	}
}

func TestRunCopyCmdEmpty(t *testing.T) {
	err := RunCopyCmd([]string{})
	if err == nil {
		t.Fatal("expected error when running copy with empty args in non-piped stdin")
	}
}

func TestRunPasteCmdEmpty(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	err := RunPasteCmd([]string{})
	if err != nil {
		t.Fatalf("expected no error when pasting empty content, got %v", err)
	}
}
