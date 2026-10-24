package constants

import (
	"strings"
	"testing"
)

func TestConstantsNonEmpty(t *testing.T) {
	if CmdCopy == "" || CmdPaste == "" {
		t.Fatal("expected CmdCopy and CmdPaste to be non-empty")
	}

	if CmdCopyAliasMem == "" || CmdCopyAliasToMem == "" || CmdPasteAliasMem == "" {
		t.Fatal("expected command aliases to be non-empty")
	}

	if DirGitmap == "" || DirMemory == "" || FileClipboard == "" {
		t.Fatal("expected directory and file constants to be non-empty")
	}

	if GitBin != "git" || GitRevParse != "rev-parse" || GitArgShowTopLevel != "--show-toplevel" {
		t.Fatal("unexpected git constants")
	}

	if ColorGreen == "" || ColorYellow == "" || ColorReset == "" {
		t.Fatal("expected color constants to be non-empty")
	}

	if GlyphClipboard == "" || GlyphSuccess == "" || GlyphWarning == "" {
		t.Fatal("expected glyph constants to be non-empty")
	}

	if MsgCopiedSuccessFmt == "" || MsgPastedSuccessFmt == "" || MsgEmptyPasteNoticeFmt == "" {
		t.Fatal("expected message format constants to be non-empty")
	}

	if HelpCopyUsage == "" || HelpCopyDesc == "" || HelpPasteUsage == "" || HelpPasteDesc == "" {
		t.Fatal("expected help constants to be non-empty")
	}

	if ErrCodeCopyEmpty == "" || ErrMsgNothingToCopy == "" || ErrMsgEmptyTopLevel == "" {
		t.Fatal("expected error constants to be non-empty")
	}

	if ExitSuccess != 0 || ExitGeneralError != 1 || ExitUsageError != 2 {
		t.Fatal("unexpected exit codes")
	}
}

func TestAllCopyFileFlags(t *testing.T) {
	flags := AllCopyFileFlags()
	if len(flags) != 6 {
		t.Fatalf("expected 6 flags, got %d", len(flags))
	}

	for _, flag := range flags {
		if !strings.HasPrefix(flag, "-") {
			t.Fatalf("expected flag %q to start with '-'", flag)
		}
	}
}
