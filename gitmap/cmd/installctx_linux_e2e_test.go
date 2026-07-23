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
	for _, l := range leaves {
		path := filepath.Join(dir, l.Label)
		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("nautilus script %q missing: %v", path, err)

			continue
		}
		if !strings.HasPrefix(string(body), "#!/bin/sh") {
			t.Errorf("%s missing shebang. body starts: %q", path, firstN(string(body), 40))
		}
		assertLinuxBodyMatchesMode(t, path, string(body), l)
	}
}

func assertLinuxBodyMatchesMode(t *testing.T, path, body string, l ctxFlatLeaf) {
	t.Helper()
	switch l.Mode {
	case constants.CtxModePrefill:
		if !strings.Contains(body, `printf "gitmap "`) {
			t.Errorf("%s prefill missing prompt. body=%s", path, body)
		}
	case constants.CtxModeSilent:
		if !strings.Contains(body, "notify-send") {
			t.Errorf("%s silent missing notify-send. body=%s", path, body)
		}
	default:
		if !strings.Contains(body, "x-terminal-emulator") {
			t.Errorf("%s terminal missing x-terminal-emulator. body=%s", path, body)
		}
		joined := strings.Join(l.Args, " ")
		if joined != "" && !strings.Contains(body, joined) {
			t.Errorf("%s missing argv %q. body=%s", path, joined, body)
		}
	}
}

func assertDolphinDesktop(t *testing.T, home string, leaves []ctxFlatLeaf) {
	t.Helper()
	path := filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("dolphin desktop missing: %v", err)
	}
	s := string(body)
	if !strings.HasPrefix(s, "[Desktop Entry]") {
		t.Errorf("dolphin desktop missing [Desktop Entry] header")
	}
	if !strings.Contains(s, "X-KDE-Submenu=gitmap") {
		t.Errorf("dolphin desktop missing X-KDE-Submenu=gitmap")
	}
	for _, l := range leaves {
		if !strings.Contains(s, "[Desktop Action "+l.Slug+"]") {
			t.Errorf("dolphin missing action section for slug %q", l.Slug)
		}
		if !strings.Contains(s, "Name="+l.Label) {
			t.Errorf("dolphin missing Name=%s", l.Label)
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
	s := string(body)
	if !strings.Contains(s, constants.CtxThunarMarkBegin) || !strings.Contains(s, constants.CtxThunarMarkEnd) {
		t.Errorf("thunar xml missing marker block")
	}
	for _, l := range leaves {
		if !strings.Contains(s, "<unique-id>"+l.Slug+"</unique-id>") {
			t.Errorf("thunar xml missing unique-id for %q", l.Slug)
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

	if _, err := os.Stat(filepath.Join(home, constants.CtxLinuxNautilusRel)); !os.IsNotExist(err) {
		t.Errorf("nautilus dir still present after uninstall: err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(home, constants.CtxLinuxDolphinRel, constants.CtxLinuxDolphinFile)); !os.IsNotExist(err) {
		t.Errorf("dolphin desktop still present after uninstall: err=%v", err)
	}
	body, _ := os.ReadFile(filepath.Join(home, constants.CtxLinuxThunarRel))
	if strings.Contains(string(body), constants.CtxThunarMarkBegin) {
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
	if count != 1 {
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
	for _, l := range leaves {
		if l.Mode == constants.CtxModePrefill {
			continue
		}
		body, err := os.ReadFile(filepath.Join(dir, l.Label))
		if err != nil {
			t.Errorf("read %s: %v", l.Label, err)

			continue
		}
		marker := "> " + l.resolvedTarget(exe) + " " + strings.Join(l.Args, " ")
		if !strings.Contains(string(body), marker) {
			t.Errorf("%s explain marker missing %q. body=%s", l.Label, marker, string(body))
		}
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
	for _, l := range leaves {
		body, err := os.ReadFile(filepath.Join(dir, l.Label))
		if err != nil {
			continue
		}
		hasZenity := strings.Contains(string(body), "zenity --question")
		if l.Extended && !hasZenity {
			t.Errorf("Extended leaf %q missing zenity guard. body=%s", l.Label, string(body))
		}
		if !l.Extended && hasZenity {
			t.Errorf("non-Extended leaf %q has zenity guard. body=%s", l.Label, string(body))
		}
	}
}

func firstN(s string, n int) string {
	if len(s) < n {
		return s
	}

	return s[:n]
}
