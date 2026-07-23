//go:build darwin || linux

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// macServicesRoot returns the absolute ~/Library/Services dir for the
// HOME-redirected sandbox the harness sets up.
func macServicesRoot(t *testing.T, home string) string {
	t.Helper()

	return filepath.Join(home, constants.CtxMacServicesRel)
}

// TestCtxMacInstallCreatesOneBundlePerLeaf drives runInstallCtxMac()
// against a sandboxed $HOME and asserts every leaf produces a
// .workflow bundle under ~/Library/Services with the two required
// payload files (Info.plist + document.wflow).
func TestCtxMacInstallCreatesOneBundlePerLeaf(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxMac()

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		bundle := filepath.Join(root, l.Slug+".workflow", "Contents")
		for _, f := range []string{"Info.plist", "document.wflow"} {
			path := filepath.Join(bundle, f)
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s/%s missing: %v", l.Slug, f, err)
			}
		}
	}
}

// TestCtxMacInfoPlistIsWellFormed parses every Info.plist as XML and
// asserts the structural contract Finder requires: NSServices array
// with NSMessage=runWorkflowAsService, NSSendFileTypes=public.folder,
// and an NSMenuItem default string equal to the leaf label.
func TestCtxMacInfoPlistIsWellFormed(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxMac()

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		path := filepath.Join(root, l.Slug+".workflow", "Contents", "Info.plist")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)

			continue
		}
		if !strings.HasPrefix(strings.TrimSpace(string(body)), `<?xml`) {
			t.Errorf("%s Info.plist missing XML declaration", l.Slug)
		}
		assertPlistContract(t, l, string(body))
	}
}

func assertPlistContract(t *testing.T, l ctxFlatLeaf, body string) {
	t.Helper()
	want := []string{
		"<key>NSServices</key>",
		"<key>NSMessage</key><string>runWorkflowAsService</string>",
		"<key>NSSendFileTypes</key>",
		"<string>public.folder</string>",
		"<string>" + l.Label + "</string>",
	}
	for _, n := range want {
		if !strings.Contains(body, n) {
			t.Errorf("%s Info.plist missing %q", l.Slug, n)
		}
	}
}

// TestCtxMacWflowEmbedsResolvedArgv asserts every document.wflow
// embeds the COMMAND_STRING that the platform shell-template would
// produce — composing target + joined argv. Catches any drift in
// macShellFor or the wflow XML escaper.
func TestCtxMacWflowEmbedsResolvedArgv(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)
	exe := resolveCtxExe()

	runInstallCtxMac()

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		path := filepath.Join(root, l.Slug+".workflow", "Contents", "document.wflow")
		body, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)

			continue
		}
		s := string(body)
		if !strings.Contains(s, "<key>COMMAND_STRING</key>") {
			t.Errorf("%s missing COMMAND_STRING key", l.Slug)
		}
		assertWflowMode(t, l, s, exe)
	}
}

func assertWflowMode(t *testing.T, l ctxFlatLeaf, body, exe string) {
	t.Helper()
	target := l.resolvedTarget(exe)
	switch l.Mode {
	case constants.CtxModePrefill:
		if !strings.Contains(body, `printf \"gitmap \"`) {
			t.Errorf("%s prefill missing prompt", l.Slug)
		}
	case constants.CtxModeSilent:
		if !strings.Contains(body, "display notification") || !strings.Contains(body, target) {
			t.Errorf("%s silent missing display notification + %s", l.Slug, target)
		}
	default:
		if !strings.Contains(body, "Terminal") || !strings.Contains(body, target) {
			t.Errorf("%s terminal missing Terminal + %s", l.Slug, target)
		}
		joined := strings.Join(l.Args, " ")
		if joined != "" && !strings.Contains(body, joined) {
			t.Errorf("%s missing argv %q", l.Slug, joined)
		}
	}
}

// TestCtxMacExtendedInjectsConfirmDialog asserts every Extended leaf
// (pull-all today) prepends the osascript display-dialog confirm
// guard, and non-Extended leaves do not.
func TestCtxMacExtendedInjectsConfirmDialog(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxMac()

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		body, err := os.ReadFile(filepath.Join(root, l.Slug+".workflow", "Contents", "document.wflow"))
		if err != nil {
			continue
		}
		hasGuard := strings.Contains(string(body), "display dialog")
		if l.Extended && !hasGuard {
			t.Errorf("Extended leaf %q missing osascript display-dialog guard", l.Slug)
		}
		if !l.Extended && hasGuard {
			t.Errorf("non-Extended leaf %q has unexpected display-dialog guard", l.Slug)
		}
	}
}

// TestCtxMacExplainInjectsAnnounce asserts --explain bakes the
// `echo "> <target> <args>"` (terminal) or `printf '...'` (silent)
// prefix into every non-Prefill bundle's COMMAND_STRING.
func TestCtxMacExplainInjectsAnnounce(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)
	exe := resolveCtxExe()

	withExplain(t, true, runInstallCtxMac)

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		if l.Mode == constants.CtxModePrefill {
			continue
		}
		body, err := os.ReadFile(filepath.Join(root, l.Slug+".workflow", "Contents", "document.wflow"))
		if err != nil {
			t.Errorf("read %s: %v", l.Slug, err)

			continue
		}
		marker := "> " + l.resolvedTarget(exe) + " " + strings.Join(l.Args, " ")
		if !strings.Contains(string(body), marker) {
			t.Errorf("%s explain marker missing %q", l.Slug, marker)
		}
	}
}

// TestCtxMacUninstallRemovesEveryBundle asserts the install→uninstall
// round-trip leaves ~/Library/Services empty of gitmap bundles. Any
// foreign Service is left alone (the test never creates one).
func TestCtxMacUninstallRemovesEveryBundle(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxMac()
	runUninstallCtxMac()

	root := macServicesRoot(t, home)
	for _, l := range leaves {
		path := filepath.Join(root, l.Slug+".workflow")
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("bundle %s still present after uninstall: err=%v", path, err)
		}
	}
}

// TestCtxMacInstallIsIdempotent runs install twice and asserts the
// second pass does not duplicate bundles or corrupt their bodies — a
// byte-equal Info.plist + document.wflow proves the writer is
// overwrite-safe.
func TestCtxMacInstallIsIdempotent(t *testing.T) {
	home := withHome(t)
	leaves := collectCtxLeaves(t)

	runInstallCtxMac()
	first := snapshotMacBundles(t, home, leaves)
	runInstallCtxMac()
	second := snapshotMacBundles(t, home, leaves)

	for k, v := range first {
		if second[k] != v {
			t.Errorf("idempotency drift at %q", k)
		}
	}
	if len(first) != len(second) {
		t.Errorf("bundle count drift: first=%d second=%d", len(first), len(second))
	}
}

func snapshotMacBundles(t *testing.T, home string, leaves []ctxFlatLeaf) map[string]string {
	t.Helper()
	out := map[string]string{}
	root := macServicesRoot(t, home)
	for _, l := range leaves {
		for _, f := range []string{"Info.plist", "document.wflow"} {
			path := filepath.Join(root, l.Slug+".workflow", "Contents", f)
			body, _ := os.ReadFile(path)
			out[l.Slug+"/"+f] = string(body)
		}
	}

	return out
}
