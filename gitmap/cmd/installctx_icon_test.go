package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestCtxWindowsHelpEntryHasIcon pins the contract that the new
// 92_help prefill entry renders with a per-entry Icon registry value
// resolved from constants.CtxIconGitmap (the {exe},0 form). One Icon
// write must appear under each of the two HKCU roots.
func TestCtxWindowsHelpEntryHasIcon(t *testing.T) {
	exe := fakeGitmapExe(t)
	cmds := buildCtxInstallCommands(exe)

	wantIcon := strings.ReplaceAll(constants.CtxIconGitmap, constants.CtxIconExeToken, exe)

	hits := 0
	for _, c := range cmds {
		if len(c) < 7 || c[0] != "reg" || c[1] != "add" || c[3] != "/v" || c[4] != "Icon" {
			continue
		}
		if !strings.HasSuffix(c[2], `\shell\92_help`) {
			continue
		}
		if c[6] != wantIcon {
			t.Errorf("92_help Icon = %q, want %q (key=%s)", c[6], wantIcon, c[2])
		}
		hits++
	}

	if hits != 2 {
		t.Errorf("got %d Icon writes for 92_help, want 2 (Background + Directory roots)", hits)
	}
}

// TestCtxWindowsNoIconWhenUnset guards the inverse: entries without an
// Icon field must not generate any Icon REG_SZ writes. Otherwise we
// would silently inherit a stale default and mask icon-resolution
// bugs (e.g. a missing pwsh.exe path).
func TestCtxWindowsNoIconWhenUnset(t *testing.T) {
	exe := fakeGitmapExe(t)
	cmds := buildCtxInstallCommands(exe)

	for _, c := range cmds {
		if len(c) < 5 || c[3] != "/v" || c[4] != "Icon" {
			continue
		}
		// The root cascade keys legitimately carry an Icon; every other
		// Icon write must correspond to an entry that opted in.
		if c[2] == constants.CtxRootKeyBackground || c[2] == constants.CtxRootKeyDirectory {
			continue
		}
		if !ctxKeyHasOptedInIcon(c[2]) {
			t.Errorf("unexpected Icon write on key without Icon opt-in: %s", c[2])
		}
	}
}

// ctxKeyHasOptedInIcon returns true when the trailing key segment
// matches a ctxEntry (top-level or child) whose Icon field is set.
func ctxKeyHasOptedInIcon(key string) bool {
	seg := key
	if i := strings.LastIndex(key, `\`); i >= 0 {
		seg = key[i+1:]
	}
	for _, e := range ctxMenu() {
		if e.KeyName == seg && e.Icon != "" {
			return true
		}
		for _, c := range e.Children {
			if c.KeyName == seg && c.Icon != "" {
				return true
			}
		}
	}

	return false
}
