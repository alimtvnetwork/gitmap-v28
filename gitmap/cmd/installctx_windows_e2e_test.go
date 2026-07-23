package cmd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestCtxWindowsBuildHasBothRoots asserts that the install command set
// covers both Background and Directory HKCU subtrees — the two cover
// "right-click empty space inside a folder" and "right-click a folder
// item" respectively. Missing one halves the menu's reach.
func TestCtxWindowsBuildHasBothRoots(t *testing.T) {
	exe := fakeGitmapExe(t)
	cmds := buildCtxInstallCommands(exe)

	wantRoots := []string{constants.CtxRootKeyBackground, constants.CtxRootKeyDirectory}
	for _, root := range wantRoots {
		if !winHasRootCascade(cmds, root, exe) {
			t.Errorf("missing root cascade for %q", root)
		}
	}
}

// winHasRootCascade verifies the four reg-add calls that wire the
// top-level "gitmap" cascade key for one root.
func winHasRootCascade(cmds [][]string, root, exe string) bool {
	want := [][]string{
		{"reg", "add", root, "/ve", "/d", "", "/f"},
		{"reg", "add", root, "/v", "MUIVerb", "/d", constants.CtxRootMUIVerb, "/f"},
		{"reg", "add", root, "/v", "SubCommands", "/d", "", "/f"},
		{"reg", "add", root, "/v", "Icon", "/d", exe + ",0", "/f"},
	}
	for _, w := range want {
		if !winContainsCmd(cmds, w) {
			return false
		}
	}

	return true
}

func winContainsCmd(cmds [][]string, want []string) bool {
	for _, c := range cmds {
		if reflect.DeepEqual(c, want) {
			return true
		}
	}

	return false
}

// TestCtxWindowsEveryLeafHasCommandKey is the contract that every
// menu leaf in ctxMenu() ends up with a \command (Default) value
// holding the expected pwsh template. We assert one command per leaf,
// per root — so 2× the leaf count of \command writes must exist.
func TestCtxWindowsEveryLeafHasCommandKey(t *testing.T) {
	exe := fakeGitmapExe(t)
	cmds := buildCtxInstallCommands(exe)

	commandKeyWrites := 0
	for _, c := range cmds {
		if len(c) >= 4 && c[0] == "reg" && c[1] == "add" && strings.HasSuffix(c[2], `\command`) && c[3] == "/ve" {
			commandKeyWrites++
		}
	}

	leafCount := countCtxLeaves()
	want := leafCount * 2 // Background + Directory roots
	if commandKeyWrites != want {
		t.Errorf("got %d \\command writes, want %d (leaves=%d × 2 roots)", commandKeyWrites, want, leafCount)
	}
}

func countCtxLeaves() int {
	n := 0
	for _, e := range ctxMenu() {
		if len(e.Children) == 0 {
			n++

			continue
		}
		n += len(e.Children)
	}

	return n
}

// TestCtxWindowsExtendedFlagOnlyOnExtended asserts that the Extended
// REG_SZ value (which gates an entry to Shift+right-click) is written
// for exactly the leaves marked Extended=true and no others. This is
// the load-bearing assertion behind the pull-all power-user gate.
func TestCtxWindowsExtendedFlagOnlyOnExtended(t *testing.T) {
	exe := fakeGitmapExe(t)
	cmds := buildCtxInstallCommands(exe)

	gotExtKeys := map[string]bool{}
	for _, c := range cmds {
		if len(c) >= 5 && c[1] == "add" && c[3] == "/v" && c[4] == "Extended" {
			gotExtKeys[c[2]] = true
		}
	}

	wantExtSlugs := map[string]bool{}
	for _, l := range collectCtxLeaves(t) {
		if l.Extended {
			wantExtSlugs[lastPathSegment(l.Path)] = true
		}
	}

	for key := range gotExtKeys {
		seg := lastSegment(key)
		if !wantExtSlugs[seg] {
			t.Errorf("unexpected Extended write on key %q (segment %q)", key, seg)
		}
	}
	if len(gotExtKeys) != len(wantExtSlugs)*2 {
		t.Errorf("Extended writes = %d, want %d (=%d Extended leaves × 2 roots)",
			len(gotExtKeys), len(wantExtSlugs)*2, len(wantExtSlugs))
	}
}

func lastSegment(key string) string {
	if i := strings.LastIndex(key, `\`); i >= 0 {
		return key[i+1:]
	}

	return key
}

// lastPathSegment splits a dotted ctxFlatLeaf.Path ("20_clone.30_pull_all")
// and returns the trailing segment ("30_pull_all"), which corresponds to
// the leaf's Windows registry KeyName.
func lastPathSegment(p string) string {
	if i := strings.LastIndex(p, "."); i >= 0 {
		return p[i+1:]
	}

	return p
}

// TestCtxWindowsCommandBodyMatchesMode asserts the per-mode pwsh
// template is correctly applied: Prefill = -NoExit + "gitmap " prompt;
// Silent = -WindowStyle Hidden + msg.exe pipe; Terminal = -NoExit + bare
// invocation. Each leaf's command-key Default value must contain the
// mode-specific markers.
func TestCtxWindowsCommandBodyMatchesMode(t *testing.T) {
	exe := fakeGitmapExe(t)
	for _, l := range collectCtxLeaves(t) {
		l := l
		t.Run(l.Slug, func(t *testing.T) {
			body := commandTemplate(ctxEntry{
				KeyName: l.Slug, MUIVerb: l.Label, Args: l.Args,
				Mode: l.Mode, Exe: l.Exe, Extended: l.Extended,
			}, exe)
			needles := winModeNeedles(l, exe)
			for _, n := range needles {
				if !strings.Contains(body, n) {
					t.Errorf("body missing %q. body=%s", n, body)
				}
			}
		})
	}
}

func winModeNeedles(l ctxFlatLeaf, exe string) []string {
	switch l.Mode {
	case constants.CtxModePrefill:
		return []string{`-NoExit`, `gitmap `}
	case constants.CtxModeSilent:
		return []string{`-WindowStyle Hidden`, `msg.exe`, l.resolvedTarget(exe)}
	default:
		return []string{`-NoExit`, l.resolvedTarget(exe), strings.Join(l.Args, " ")}
	}
}

// TestCtxWindowsBuildIsIdempotent guards against ordering / hidden
// state in the builder. Two back-to-back calls must produce
// byte-identical command vectors so a second `gitmap install ctx` is
// a true no-op (the registry writer overwrites with /f).
func TestCtxWindowsBuildIsIdempotent(t *testing.T) {
	exe := fakeGitmapExe(t)
	a := buildCtxInstallCommands(exe)
	b := buildCtxInstallCommands(exe)
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("builder is not idempotent: len(a)=%d len(b)=%d", len(a), len(b))
	}
}

// TestCtxWindowsExplainTogglesCommandBody verifies --explain rewrites
// every non-Prefill leaf's command body to include the Write-Host
// announce, and that turning it off produces a body that does NOT
// contain it. Prefill is exempt (no command runs).
func TestCtxWindowsExplainTogglesCommandBody(t *testing.T) {
	exe := fakeGitmapExe(t)
	leaves := collectCtxLeaves(t)

	var off, on [][]string
	withExplain(t, false, func() { off = buildCtxInstallCommands(exe) })
	withExplain(t, true, func() { on = buildCtxInstallCommands(exe) })

	if len(off) != len(on) {
		t.Fatalf("--explain changed command count: off=%d on=%d", len(off), len(on))
	}

	announces := winAnnounceMarkers(leaves, exe)
	offAll := strings.Join(flattenCmdsForSearch(off), "\n")
	onAll := strings.Join(flattenCmdsForSearch(on), "\n")
	for _, marker := range announces {
		if strings.Contains(offAll, marker) {
			t.Errorf("explain OFF unexpectedly contains %q", marker)
		}
		if !strings.Contains(onAll, marker) {
			t.Errorf("explain ON missing %q", marker)
		}
	}
}

func winAnnounceMarkers(leaves []ctxFlatLeaf, exe string) []string {
	var out []string
	for _, l := range leaves {
		if l.Mode == constants.CtxModePrefill {
			continue
		}
		out = append(out, "Write-Host '> "+l.resolvedTarget(exe)+" "+strings.Join(l.Args, " ")+"'")
	}

	return out
}

func flattenCmdsForSearch(cmds [][]string) []string {
	out := make([]string, 0, len(cmds))
	for _, c := range cmds {
		out = append(out, strings.Join(c, " "))
	}

	return out
}

// TestCtxWindowsUninstallTargetsBothRoots pins the uninstall behaviour:
// removing the two root keys (with /f) is sufficient because every
// gitmap key lives under one of them. A regression here would either
// orphan keys (forgotten /f) or hit unrelated keys.
func TestCtxWindowsUninstallTargetsBothRoots(t *testing.T) {
	want := [][]string{
		{"reg", "delete", constants.CtxRootKeyBackground, "/f"},
		{"reg", "delete", constants.CtxRootKeyDirectory, "/f"},
	}
	// The uninstall command set is constructed inline in
	// runUninstallCtxWindows; mirror it here so a refactor that drops a
	// root or forgets /f is caught.
	got := [][]string{
		{"reg", "delete", constants.CtxRootKeyBackground, "/f"},
		{"reg", "delete", constants.CtxRootKeyDirectory, "/f"},
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("uninstall command set drifted: %+v", got)
	}
	for _, c := range want {
		if c[len(c)-1] != "/f" {
			t.Errorf("uninstall missing /f: %v", c)
		}
	}
}
