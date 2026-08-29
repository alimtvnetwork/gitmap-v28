package cmd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestCtxArgvCoversEveryLeaf is the cross-platform argv contract:
// every ctxMenu() leaf must appear in expectedCtxArgv, and every
// expectedCtxArgv entry must resolve to a real leaf. This is the
// inverse of installctxentries_argv_test.go's per-path check; here we
// run through the *flatten* path (the one macOS/Linux actually use)
// to prove both representations agree.
func TestCtxArgvCoversEveryLeaf(t *testing.T) {
	leaves := collectCtxLeaves(t)
	if len(leaves) == 0 {
		t.Fatal("flattenCtxMenu returned no leaves")
	}

	wantPaths := map[string]bool{}
	for path := range expectedCtxArgv {
		wantPaths[path] = true
	}

	seen := map[string]bool{}
	for _, l := range leaves {
		key := pathFromKeyName(l.Path)
		seen[key] = true
		want, ok := expectedCtxArgv[key]
		if !ok {
			t.Errorf("flat leaf %q (slug=%s) has no expectedCtxArgv entry", key, l.Slug)

			continue
		}
		if !reflect.DeepEqual(l.Args, want.argv) {
			t.Errorf("%s: flat Args = %v, want %v", key, l.Args, want.argv)
		}
		if l.Exe != want.exe {
			t.Errorf("%s: flat Exe = %q, want %q", key, l.Exe, want.exe)
		}
	}
	for path := range wantPaths {
		if !seen[path] {
			t.Errorf("expectedCtxArgv has %q but flatten produced no matching leaf", path)
		}
	}
}

// pathFromKeyName converts a "parent.child" KeyName path (as recorded
// by buildLeafPaths) into the "parent/child" form used as the key in
// expectedCtxArgv. Top-level leaves have no "." and pass through.
func pathFromKeyName(p string) string {
	return strings.ReplaceAll(p, ".", "/")
}

// TestCtxExplainAffectsEveryNonPrefillLeaf drives the platform
// rendering helpers with explain=on / explain=off and asserts the
// announce/echo prefix is present iff explain is on, for every leaf
// that is *not* Prefill (Prefill mode never runs a command, so it has
// nothing to announce).
func TestCtxExplainAffectsEveryNonPrefillLeaf(t *testing.T) {
	exe := fakeGitmapExe(t)
	leaves := collectCtxLeaves(t)

	for _, l := range leaves {
		if l.Mode == constants.CtxModePrefill {
			continue
		}
		l := l
		t.Run(l.Slug, func(t *testing.T) {
			target := l.resolvedTarget(exe)
			expect := "> " + target + " " + strings.Join(l.Args, " ")

			var off, on string
			withExplain(t, false, func() {
				off = renderAllPlatformsForLeaf(l, exe)
			})
			withExplain(t, true, func() {
				on = renderAllPlatformsForLeaf(l, exe)
			})

			if strings.Contains(off, expect) {
				t.Errorf("explain OFF unexpectedly contains announce %q in:\n%s", expect, off)
			}
			if !strings.Contains(on, expect) {
				t.Errorf("explain ON missing announce %q in:\n%s", expect, on)
			}
		})
	}
}

// renderAllPlatformsForLeaf concatenates the Windows pwsh template,
// the Linux shell-script body and the macOS Automator shell payload
// so a single substring assertion proves every platform honours the
// explain toggle. ctxEntry / flatCtxEntry are reconstructed from the
// harness leaf so the renderers see exactly what the install path
// would feed them.
func renderAllPlatformsForLeaf(l ctxFlatLeaf, exe string) string {
	cE := ctxEntry{KeyName: l.Slug, MUIVerb: l.Label, Args: l.Args, Mode: l.Mode, Exe: l.Exe, Extended: l.Extended}
	fE := flatCtxEntry{Label: l.Label, Slug: l.Slug, Args: l.Args, Mode: l.Mode, Exe: l.Exe, Extended: l.Extended}

	return commandTemplate(cE, exe) + "\n" + linuxShellScript(fE, exe) + "\n" + macShellFor(fE, exe)
}

// TestCtxReleaseNextE2EArgvComposes pins the highest-stakes entry —
// release-next — through every renderer. Any drift in FlagBumpDash
// or BumpMinor (or accidental quoting) fails this test on every OS.
func TestCtxReleaseNextE2EArgvComposes(t *testing.T) {
	exe := fakeGitmapExe(t)
	want := []string{constants.CmdRelease, constants.FlagBumpDash, constants.BumpMinor}
	wantJoined := strings.Join(want, " ")

	for _, l := range collectCtxLeaves(t) {
		if l.Path != "30_release.20_release_next" {
			continue
		}
		if !reflect.DeepEqual(l.Args, want) {
			t.Fatalf("release-next args = %v, want %v", l.Args, want)
		}
		body := renderAllPlatformsForLeaf(l, exe)
		if !strings.Contains(body, wantJoined) {
			t.Fatalf("release-next render missing composed argv %q. Body:\n%s", wantJoined, body)
		}

		return
	}
	t.Fatal("release-next leaf not found")
}

// TestCtxPullAllIsExtendedEverywhere proves the Shift-click /
// confirm-prompt gating reaches every platform: Windows leafCommands
// must emit the Extended REG_SZ; Linux/macOS must inject the guard
// (zenity-chain or osascript dialog) into the script body.
func TestCtxPullAllIsExtendedEverywhere(t *testing.T) {
	exe := fakeGitmapExe(t)
	leaves := collectCtxLeaves(t)
	var got *ctxFlatLeaf
	for i := range leaves {
		if leaves[i].Path == "20_clone.30_pull_all" {
			got = &leaves[i]

			break
		}
	}
	if got == nil {
		t.Fatal("pull-all leaf not found")
	}
	if !got.Extended {
		t.Fatal("pull-all must be Extended=true")
	}

	winCmds := leafCommands(`HKCU\Software\Classes\Directory\Background\shell\gitmap\shell\20_clone\shell\30_pull_all`,
		ctxEntry{KeyName: got.Slug, MUIVerb: got.Label, Args: got.Args, Mode: got.Mode, Extended: got.Extended}, exe)
	if !winCmdsHaveExtended(winCmds) {
		t.Errorf("Windows leafCommands missing Extended REG_SZ for pull-all: %#v", winCmds)
	}

	fE := flatCtxEntry{Label: got.Label, Slug: got.Slug, Args: got.Args, Mode: got.Mode, Extended: true}
	if guard := extendedGuard(fE); guard == "" || !strings.Contains(guard, "zenity") {
		t.Errorf("Linux extendedGuard returned no zenity-chain: %q", guard)
	}
	macBody := macShellFor(fE, exe)
	if !strings.Contains(macBody, "display dialog") {
		t.Errorf("macOS shell missing osascript confirm dialog:\n%s", macBody)
	}
}

func winCmdsHaveExtended(cmds [][]string) bool {
	for _, c := range cmds {
		joined := strings.Join(c, " ")
		if strings.Contains(joined, "/v Extended") {
			return true
		}
	}

	return false
}
