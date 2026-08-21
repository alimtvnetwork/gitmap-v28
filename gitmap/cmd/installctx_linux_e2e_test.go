//go:build linux || darwin

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// withHome redirects $HOME to a fresh temp dir for the duration of the
// test. os.UserHomeDir() reads $HOME on linux and darwin, so this is
// the cleanest seam for a true E2E without touching the real home.
func withHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	prev := ctxExplainEnabled
	t.Cleanup(func() { ctxExplainEnabled = prev })

	return dir
}

// TestCtxLinuxInstallCreatesAllManagerArtifacts drives the real
// runInstallCtxLinux into a sandboxed $HOME and asserts every file
// manager backend (Nautilus / Dolphin / Thunar) emits the expected
// artifacts: one Nautilus script per leaf, one Dolphin .desktop with
// every Action ID, and a marker-wrapped Thunar uca.xml block.
func TestCtxLinuxInstallCreatesAllManagerArtifacts(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxLinux()

	assertNautilusScripts(t, home, leaves)
	assertDolphinDesktop(t, home, leaves)
	assertThunarXML(t, home, leaves)
}

func assertNautilusScripts(t *testing.T, home string, leaves []ctxFlatLeaf) {
	t.Helper()
	dir := filepath.Join(home, constants.CtxLinuxNautilusRel)
	for _, leaf := range leaves {
		assertSingleNautilusScript(t, dir, leaf)
	}
}

func assertSingleNautilusScript(t *testing.T, dir string, leaf ctxFlatLeaf) {
	t.Helper()
	path := filepath.Join(dir, leaf.Label)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("nautilus script %q missing: %v", path, err)
		return
	}
	hasShebang := strings.HasPrefix(string(body), "#!/bin/sh")
	isShebangMissing := !hasShebang
	if isShebangMissing == true {
		t.Errorf("%s missing shebang. body starts: %q", path, firstN(string(body), 40))
	}
	assertLinuxBodyMatchesMode(t, path, string(body), leaf)
}

func assertLinuxBodyMatchesMode(t *testing.T, path, body string, leaf ctxFlatLeaf) {
	t.Helper()
	switch leaf.Mode {
	case constants.CtxModePrefill:
		assertLinuxPrefillBody(t, path, body)
	case constants.CtxModeSilent:
		assertLinuxSilentBody(t, path, body)
	default:
		assertLinuxTerminalBody(t, path, body, leaf)
	}
}

func assertLinuxPrefillBody(t *testing.T, path, body string) {
	t.Helper()
	hasPrompt := strings.Contains(body, `printf "gitmap "`)
	isPromptMissing := !hasPrompt
	if isPromptMissing == true {
		t.Errorf("%s prefill missing prompt. body=%s", path, body)
	}
}

func assertLinuxSilentBody(t *testing.T, path, body string) {
	t.Helper()
	hasNotify := strings.Contains(body, "notify-send")
	isNotifyMissing := !hasNotify
	if isNotifyMissing == true {
		t.Errorf("%s silent missing notify-send. body=%s", path, body)
	}
}

func assertLinuxTerminalBody(t *testing.T, path, body string, leaf ctxFlatLeaf) {
	t.Helper()
	hasTerminal := strings.Contains(body, "x-terminal-emulator")
	isTerminalMissing := !hasTerminal
	if isTerminalMissing == true {
		t.Errorf("%s terminal missing x-terminal-emulator. body=%s", path, body)
	}
	joined := strings.Join(leaf.Args, " ")
	hasJoinedArgv := joined != ""
	containsArgv := strings.Contains(body, joined)
	isArgvMissing := !containsArgv
	if hasJoinedArgv == true && isArgvMissing == true {
		t.Errorf("%s missing argv %q. body=%s", path, joined, body)
	}
}

func assertDolphinDesktop(t *testing.T, home string, leaves []ctxFlatLeaf) {
	t.Helper()
	path := filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("dolphin desktop missing: %v", err)
	}
	desktopContent := string(body)
	assertDolphinHeaders(t, desktopContent)
	assertDolphinActions(t, desktopContent, leaves)
}

func assertDolphinHeaders(t *testing.T, desktopContent string) {
	t.Helper()
	hasHeader := strings.HasPrefix(desktopContent, "[Desktop Entry]")
	isHeaderMissing := !hasHeader
	if isHeaderMissing == true {
		t.Errorf("dolphin desktop missing [Desktop Entry] header")
	}
	hasSubmenu := strings.Contains(desktopContent, "X-KDE-Submenu=gitmap")
	isSubmenuMissing := !hasSubmenu
	if isSubmenuMissing == true {
		t.Errorf("dolphin desktop missing X-KDE-Submenu=gitmap")
	}
}

func assertDolphinActions(t *testing.T, desktopContent string, leaves []ctxFlatLeaf) {
	t.Helper()
	for _, leaf := range leaves {
		hasAction := strings.Contains(desktopContent, "[Desktop Action "+leaf.Slug+"]")
		isActionMissing := !hasAction
		if isActionMissing == true {
			t.Errorf("dolphin missing action section for slug %q", leaf.Slug)
		}
		hasName := strings.Contains(desktopContent, "Name="+leaf.Label)
		isNameMissing := !hasName
		if isNameMissing == true {
			t.Errorf("dolphin missing Name=%s", leaf.Label)
		}
	}
}

func assertThunarXML(t *testing.T, home string, leaves []ctxFlatLeaf) {
	t.Helper()
	path := filepath.Join(home, constants.CtxLinuxThunarRel)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("thunar uca.xml missing: %v", err)
	}
	xmlContent := string(body)
	assertThunarMarkers(t, xmlContent)
	assertThunarUniqueIDs(t, xmlContent, leaves)
}

func assertThunarMarkers(t *testing.T, xmlContent string) {
	t.Helper()
	hasBeginMark := strings.Contains(xmlContent, constants.CtxThunarMarkBegin)
	hasEndMark := strings.Contains(xmlContent, constants.CtxThunarMarkEnd)
	isBeginMissing := !hasBeginMark
	isEndMissing := !hasEndMark
	if isBeginMissing == true || isEndMissing == true {
		t.Errorf("thunar xml missing marker block")
	}
}

func assertThunarUniqueIDs(t *testing.T, xmlContent string, leaves []ctxFlatLeaf) {
	t.Helper()
	for _, leaf := range leaves {
		hasUniqueID := strings.Contains(xmlContent, "<unique-id>"+leaf.Slug+"</unique-id>")
		isUniqueIDMissing := !hasUniqueID
		if isUniqueIDMissing == true {
			t.Errorf("thunar xml missing unique-id for %q", leaf.Slug)
		}
	}
}

// TestCtxLinuxUninstallRemovesEverythingInstall added asserts the
// reverse direction: after install→uninstall, the Nautilus dir is
// gone, the Dolphin .desktop is gone, and the Thunar marker block has
// been stripped (leaving any user-managed entries — none in this test
// — alone).
func TestCtxLinuxUninstallRemovesEverythingInstallAdded(t *testing.T) {
	home := withHome(t)

	runInstallCtxLinux()
	runUninstallCtxLinux()

	assertUninstallPathsRemoved(t, home)
	assertUninstallThunarCleaned(t, home)
}

func assertUninstallPathsRemoved(t *testing.T, home string) {
	t.Helper()
	_, nautilusErr := os.Stat(filepath.Join(home, constants.CtxLinuxNautilusRel))
	isNautilusNotExist := os.IsNotExist(nautilusErr)
	isNautilusPresent := !isNautilusNotExist
	if isNautilusPresent == true {
		t.Errorf("nautilus dir still present after uninstall: err=%v", nautilusErr)
	}
	_, dolphinErr := os.Stat(filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile))
	isDolphinNotExist := os.IsNotExist(dolphinErr)
	isDolphinPresent := !isDolphinNotExist
	if isDolphinPresent == true {
		t.Errorf("dolphin desktop still present after uninstall: err=%v", dolphinErr)
	}
}

func assertUninstallThunarCleaned(t *testing.T, home string) {
	t.Helper()
	body, _ := os.ReadFile(filepath.Join(home, constants.CtxLinuxThunarRel))
	hasBeginMark := strings.Contains(string(body), constants.CtxThunarMarkBegin)
	if hasBeginMark == true {
		t.Errorf("thunar marker block not stripped: %s", string(body))
	}
}

// TestCtxLinuxThunarIsIdempotent installs twice and asserts the marker
// block appears exactly once — proving the splice path in
// thunarMerged() replaces rather than duplicating. Catches the class
// of bug where re-running the installer doubles the menu.
func TestCtxLinuxThunarIsIdempotent(t *testing.T) {
	home := withHome(t)

	runInstallCtxLinux()
	runInstallCtxLinux()

	body, err := os.ReadFile(filepath.Join(home, constants.CtxLinuxThunarRel))
	if err != nil {
		t.Fatalf("thunar uca.xml missing: %v", err)
	}
	count := strings.Count(string(body), constants.CtxThunarMarkBegin)
	isNotSingleBlock := count != 1
	if isNotSingleBlock == true {
		t.Fatalf("thunar marker block appears %d times after double-install, want 1", count)
	}
}

// TestCtxLinuxExplainInjectsAnnounce drives install with --explain
// enabled and asserts every non-Prefill Nautilus script contains the
// `echo '> <target> <args>'` (terminal) or printf-announce (silent)
// prefix. Reuses the harness's withExplain() guard.
func TestCtxLinuxExplainInjectsAnnounce(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	withExplain(t, true, runInstallCtxLinux)

	dir := filepath.Join(home, constants.CtxLinuxNautilusRel)
	exe := resolveCtxExe()
	assertLeavesExplainAnnounce(t, dir, exe, leaves)
}

func assertLeavesExplainAnnounce(t *testing.T, dir, exe string, leaves []ctxFlatLeaf) {
	t.Helper()
	for _, leaf := range leaves {
		isPrefill := leaf.Mode == constants.CtxModePrefill
		if isPrefill == true {
			continue
		}
		assertSingleExplainAnnounce(t, dir, exe, leaf)
	}
}

func assertSingleExplainAnnounce(t *testing.T, dir, exe string, leaf ctxFlatLeaf) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, leaf.Label))
	if err != nil {
		t.Errorf("read %s: %v", leaf.Label, err)
		return
	}
	marker := "> " + leaf.resolvedTarget(exe) + " " + strings.Join(leaf.Args, " ")
	hasMarker := strings.Contains(string(body), marker)
	isMarkerMissing := !hasMarker
	if isMarkerMissing == true {
		t.Errorf("%s explain marker missing %q. body=%s", leaf.Label, marker, string(body))
	}
}

// TestCtxLinuxExtendedGuardOnlyOnExtended asserts the zenity/kdialog/
// xmessage confirm-prompt chain appears in the body of Extended
// entries (pull-all today) and is absent from non-Extended ones.
func TestCtxLinuxExtendedGuardOnlyOnExtended(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxLinux()

	dir := filepath.Join(home, constants.CtxLinuxNautilusRel)
	for _, leaf := range leaves {
		assertSingleExtendedGuard(t, dir, leaf)
	}
}

func assertSingleExtendedGuard(t *testing.T, dir string, leaf ctxFlatLeaf) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, leaf.Label))
	if err != nil {
		return
	}
	hasZenity := strings.Contains(string(body), "zenity --question")
	checkExtendedGuardMatches(t, leaf, string(body), hasZenity)
}

func checkExtendedGuardMatches(t *testing.T, leaf ctxFlatLeaf, body string, hasZenity bool) {
	t.Helper()
	isZenityMissing := !hasZenity
	isExtendedMissing := !leaf.Extended
	isExtendedWithoutGuard := leaf.Extended == true && isZenityMissing == true
	if isExtendedWithoutGuard == true {
		t.Errorf("Extended leaf %q missing zenity guard. body=%s", leaf.Label, body)
	}
	isNonExtendedWithGuard := isExtendedMissing == true && hasZenity == true
	if isNonExtendedWithGuard == true {
		t.Errorf("non-Extended leaf %q has zenity guard. body=%s", leaf.Label, body)
	}
}

func firstN(str string, limit int) string {
	isWithinLimit := len(str) < limit
	if isWithinLimit == true {
		return str
	}

	return str[:limit]
}
