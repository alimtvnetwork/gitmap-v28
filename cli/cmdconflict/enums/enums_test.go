package enums

import (
	"testing"
)

func TestConflictActionType(t *testing.T) {
	testConflictActionAllAndString(t)
	testConflictActionParse(t)
	testConflictActionSideAliases(t)
	testConflictActionIDEAliases(t)
	testConflictActionAbortAndRerereAliases(t)
	testConflictActionSidePredicates(t)
	testConflictActionControlPredicates(t)
}

func testConflictActionAllAndString(t *testing.T) {
	all := AllConflictActionTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 action variants, got %d", len(all))
	}

	for _, act := range all {
		if !act.IsValid() || !act.HasAction() {
			t.Errorf("expected valid for %s", act)
		}
		if act.String() != string(act) {
			t.Errorf("expected string representation to match %s", act)
		}
	}
}

func testConflictActionParse(t *testing.T) {
	parsed, err := ParseConflictActionType("  THEIRS ")
	if err != nil || parsed != ConflictActionTheirs {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseConflictActionType("unknown"); err == nil {
		t.Error("expected error parsing unknown action")
	}

	if ConflictActionType("invalid").IsValid() {
		t.Error("expected invalid action to return false on IsValid()")
	}
}

func testConflictActionSideAliases(t *testing.T) {
	aliases := map[string]ConflictActionType{
		"-t":       ConflictActionTheirs,
		"--theirs": ConflictActionTheirs,
		"-o":       ConflictActionOurs,
		"--ours":   ConflictActionOurs,
	}
	for alias, expected := range aliases {
		act, err := ParseConflictActionType(alias)
		if err != nil || act != expected {
			t.Errorf("alias %q: expected %s, got %s", alias, expected, act)
		}
	}
}

func testConflictActionIDEAliases(t *testing.T) {
	for _, alias := range []string{"-i", "--ide", "code", "--code", "resolve", "v"} {
		act, err := ParseConflictActionType(alias)
		if err != nil || act != ConflictActionIDE {
			t.Errorf("alias %q: expected %s, got %s", alias, ConflictActionIDE, act)
		}
	}
}

func testConflictActionAbortAndRerereAliases(t *testing.T) {
	aliases := map[string]ConflictActionType{
		"-a":       ConflictActionAbort,
		"--abort":  ConflictActionAbort,
		"--rerere": ConflictActionRerere,
	}
	for alias, expected := range aliases {
		act, err := ParseConflictActionType(alias)
		if err != nil || act != expected {
			t.Errorf("alias %q: expected %s, got %s", alias, expected, act)
		}
	}
}

func testConflictActionSidePredicates(t *testing.T) {
	if !ConflictActionTheirs.IsTheirs() || ConflictActionOurs.IsTheirs() {
		t.Error("IsTheirs predicate failure")
	}
	if !ConflictActionOurs.IsOurs() || ConflictActionIDE.IsOurs() {
		t.Error("IsOurs predicate failure")
	}
}

func testConflictActionControlPredicates(t *testing.T) {
	if !ConflictActionIDE.IsIDE() || ConflictActionAbort.IsIDE() {
		t.Error("IsIDE predicate failure")
	}
	if !ConflictActionAbort.IsAbort() || ConflictActionRerere.IsAbort() {
		t.Error("IsAbort predicate failure")
	}
	if !ConflictActionRerere.IsRerere() || ConflictActionTheirs.IsRerere() {
		t.Error("IsRerere predicate failure")
	}
}

func TestConflictToolType(t *testing.T) {
	testConflictToolAllAndString(t)
	testConflictToolParse(t)
	testConflictToolCodeAndCursorNames(t)
	testConflictToolWindsurfAndCodiumNames(t)
	testConflictToolPredicates(t)
}

func testConflictToolAllAndString(t *testing.T) {
	all := AllConflictToolTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 tool variants, got %d", len(all))
	}

	for _, tool := range all {
		if !tool.IsValid() || !tool.HasTool() {
			t.Errorf("expected valid for %s", tool)
		}
		if tool.String() != string(tool) {
			t.Errorf("expected string representation to match %s", tool)
		}
	}
}

func testConflictToolParse(t *testing.T) {
	parsed, err := ParseConflictToolType("  CURSOR ")
	if err != nil || parsed != ConflictToolCursor {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseConflictToolType("notepad"); err == nil {
		t.Error("expected error parsing invalid tool")
	}

	if ConflictToolType("vim").IsValid() {
		t.Error("expected invalid tool to return false on IsValid()")
	}
}

func testConflictToolCodeAndCursorNames(t *testing.T) {
	if ConflictToolCode.DisplayName() != "VS Code" || ConflictToolCode.BinaryName() != "code" {
		t.Error("code tool display/binary name mismatch")
	}
	if ConflictToolCursor.DisplayName() != "Cursor" || ConflictToolCursor.BinaryName() != "cursor" {
		t.Error("cursor tool display/binary name mismatch")
	}
}

func testConflictToolWindsurfAndCodiumNames(t *testing.T) {
	if ConflictToolWindsurf.DisplayName() != "Windsurf" || ConflictToolWindsurf.BinaryName() != "windsurf" {
		t.Error("windsurf tool display/binary name mismatch")
	}
	if ConflictToolCodium.DisplayName() != "VSCodium" || ConflictToolCodium.BinaryName() != "codium" {
		t.Error("codium tool display/binary name mismatch")
	}
	if ConflictToolType("custom").DisplayName() != "custom" {
		t.Error("fallback display name mismatch")
	}
}

func testConflictToolPredicates(t *testing.T) {
	if !ConflictToolCode.IsCode() || ConflictToolCursor.IsCode() {
		t.Error("IsCode predicate failure")
	}
	if !ConflictToolCursor.IsCursor() || ConflictToolWindsurf.IsCursor() {
		t.Error("IsCursor predicate failure")
	}
	if !ConflictToolWindsurf.IsWindsurf() || ConflictToolCodium.IsWindsurf() {
		t.Error("IsWindsurf predicate failure")
	}
	if !ConflictToolCodium.IsCodium() || ConflictToolCode.IsCodium() {
		t.Error("IsCodium predicate failure")
	}
}

func TestConflictStatusCodeType(t *testing.T) {
	testConflictStatusAllAndString(t)
	testConflictStatusParse(t)
	testConflictStatusPrimaryDescriptions(t)
	testConflictStatusUserDescriptions(t)
	testConflictStatusOtherDescriptions(t)
	testConflictStatusPredicates(t)
}

func testConflictStatusAllAndString(t *testing.T) {
	all := AllConflictStatusCodeTypes()
	if len(all) != 7 {
		t.Errorf("expected 7 status code variants, got %d", len(all))
	}

	for _, code := range all {
		if !code.IsValid() || !code.HasConflict() || !code.IsUnmerged() {
			t.Errorf("expected valid unmerged for %s", code)
		}
		if code.String() != string(code) {
			t.Errorf("expected string representation to match %s", code)
		}
	}
}

func testConflictStatusParse(t *testing.T) {
	parsed, err := ParseConflictStatusCodeType("  uu ")
	if err != nil || parsed != ConflictStatusBothModified {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseConflictStatusCodeType("M "); err == nil {
		t.Error("expected error parsing non-conflict status")
	}

	if ConflictStatusCodeType("??").IsValid() {
		t.Error("expected non-conflict code to return false on IsValid()")
	}
}

func testConflictStatusPrimaryDescriptions(t *testing.T) {
	if ConflictStatusBothModified.Description() != "both modified" {
		t.Error("UU description mismatch")
	}
	if ConflictStatusBothAdded.Description() != "both added" {
		t.Error("AA description mismatch")
	}
	if ConflictStatusBothDeleted.Description() != "both deleted" {
		t.Error("DD description mismatch")
	}
}

func testConflictStatusUserDescriptions(t *testing.T) {
	if ConflictStatusUserDeleted.Description() != "updated by us, deleted by them" {
		t.Error("UD description mismatch")
	}
	if ConflictStatusDeletedUser.Description() != "deleted by us, updated by them" {
		t.Error("DU description mismatch")
	}
}

func testConflictStatusOtherDescriptions(t *testing.T) {
	if ConflictStatusAddedUser.Description() != "added by us, updated by them" {
		t.Error("AU description mismatch")
	}
	if ConflictStatusUserAdded.Description() != "updated by us, added by them" {
		t.Error("UA description mismatch")
	}
	if ConflictStatusCodeType("XX").Description() != "unknown conflict status" {
		t.Error("fallback description mismatch")
	}
}

func testConflictStatusPredicates(t *testing.T) {
	if !ConflictStatusBothModified.IsBothModified() || ConflictStatusBothAdded.IsBothModified() {
		t.Error("IsBothModified predicate failure")
	}
	if !ConflictStatusBothAdded.IsBothAdded() || ConflictStatusBothDeleted.IsBothAdded() {
		t.Error("IsBothAdded predicate failure")
	}
	if !ConflictStatusBothDeleted.IsBothDeleted() || ConflictStatusBothModified.IsBothDeleted() {
		t.Error("IsBothDeleted predicate failure")
	}
}

func TestConflictScopeType(t *testing.T) {
	testConflictScopeAllAndString(t)
	testConflictScopeParse(t)
	testConflictScopeResolution(t)
	testConflictScopePredicates(t)
}

func testConflictScopeAllAndString(t *testing.T) {
	all := AllConflictScopeTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 scope variants, got %d", len(all))
	}

	for _, scope := range all {
		if !scope.IsValid() || !scope.HasScope() {
			t.Errorf("expected valid for %s", scope)
		}
		if scope.String() != string(scope) {
			t.Errorf("expected string representation to match %s", scope)
		}
	}
}

func testConflictScopeParse(t *testing.T) {
	parsed, err := ParseConflictScopeType("  TARGET ")
	if err != nil || parsed != ConflictScopeTarget {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseConflictScopeType("none"); err == nil {
		t.Error("expected error parsing invalid scope")
	}

	if ConflictScopeType("other").IsValid() {
		t.Error("expected invalid scope to return false on IsValid()")
	}
}

func testConflictScopeResolution(t *testing.T) {
	if ResolveConflictScope("my-repo", false) != ConflictScopeTarget {
		t.Error("expected target scope when target is given")
	}
	if ResolveConflictScope("", true) != ConflictScopeCurrent {
		t.Error("expected current scope when inside git repo")
	}
	if ResolveConflictScope("", false) != ConflictScopeAll {
		t.Error("expected all scope when not inside git repo and no target")
	}
}

func testConflictScopePredicates(t *testing.T) {
	if !ConflictScopeTarget.IsTarget() || ConflictScopeCurrent.IsTarget() {
		t.Error("IsTarget predicate failure")
	}
	if !ConflictScopeCurrent.IsCurrent() || ConflictScopeAll.IsCurrent() {
		t.Error("IsCurrent predicate failure")
	}
	if !ConflictScopeAll.IsAll() || ConflictScopeTarget.IsAll() {
		t.Error("IsAll predicate failure")
	}
}
