package constants

import (
	"testing"
)

func TestCommandConstants(t *testing.T) {
	if CmdAdd != "add" {
		t.Errorf("expected add, got %s", CmdAdd)
	}
	if CmdAddAlias != "ad" {
		t.Errorf("expected ad, got %s", CmdAddAlias)
	}
}

func TestSubcommandConstants(t *testing.T) {
	if SubCmdCommonAttr != "common-attr" {
		t.Errorf("expected common-attr, got %s", SubCmdCommonAttr)
	}
	if SubCmdCommonIgnore != "common-ignore" {
		t.Errorf("expected common-ignore, got %s", SubCmdCommonIgnore)
	}
	if SubCmdAttributes != "attributes" {
		t.Errorf("expected attributes, got %s", SubCmdAttributes)
	}
	if SubCmdIgnore != "ignore" {
		t.Errorf("expected ignore, got %s", SubCmdIgnore)
	}
	if SubCmdLfsInstall != "lfs-install" {
		t.Errorf("expected lfs-install, got %s", SubCmdLfsInstall)
	}
}

func TestFileAndPermissionConstants(t *testing.T) {
	if FileGitAttributes != ".gitattributes" {
		t.Errorf("expected .gitattributes, got %s", FileGitAttributes)
	}
	if FileGitIgnore != ".gitignore" {
		t.Errorf("expected .gitignore, got %s", FileGitIgnore)
	}
	if DefaultFileMode != 0644 {
		t.Errorf("expected 0644, got %o", DefaultFileMode)
	}
}

func TestTemplateConstants(t *testing.T) {
	if len(TemplateCommonAttr) == 0 {
		t.Error("expected non-empty TemplateCommonAttr")
	}
	if len(TemplateCommonIgnore) == 0 {
		t.Error("expected non-empty TemplateCommonIgnore")
	}
}

func TestErrorCodeAndOperationConstants(t *testing.T) {
	if ErrCodeAddMissingArg != "E_ADD_MISSING_ARG" {
		t.Errorf("expected E_ADD_MISSING_ARG, got %s", ErrCodeAddMissingArg)
	}
	if ErrCodeAddUnknownTarget != "E_ADD_UNKNOWN_TARGET" {
		t.Errorf("expected E_ADD_UNKNOWN_TARGET, got %s", ErrCodeAddUnknownTarget)
	}
	if ErrPrefixError != "Error:" {
		t.Errorf("expected Error:, got %s", ErrPrefixError)
	}
	if OpAddRun != "AddRun" {
		t.Errorf("expected AddRun, got %s", OpAddRun)
	}
}
