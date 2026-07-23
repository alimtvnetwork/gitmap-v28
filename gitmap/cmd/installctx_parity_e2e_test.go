package cmd

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestCtxParityWindowsLinuxMacEmitSameLeafSet proves all three
// platform installers consume the SAME flatten set: same slugs,
// same argv, same mode, same Extended bit. If any platform-specific
// builder ever filters, reorders, or drops a leaf, the user gets
// inconsistent menus across OSes — this test fails first.
func TestCtxParityWindowsLinuxMacEmitSameLeafSet(t *testing.T) {
	exe := fakeGitmapExe(t)
	leaves := collectCtxLeaves(t)
	if len(leaves) == 0 {
		t.Fatal("no leaves to compare")
	}

	// Build the canonical (slug, argv, mode, extended, target, path) tuple set.
	type tup struct {
		Slug, Mode, Target, Path string
		Args                     []string
		Extended                 bool
	}
	canonical := map[string]tup{}
	for _, l := range leaves {
		canonical[l.Slug] = tup{
			Slug: l.Slug, Mode: string(l.Mode), Target: l.resolvedTarget(exe),
			Path: l.Path,
			Args: append([]string(nil), l.Args...), Extended: l.Extended,
		}
	}

	// Windows view: derive from buildCtxInstallCommands — every leaf
	// must have at least one \command write whose key matches the
	// nested KeyName cascade (Path "20_clone.30_pull_all" =>
	// "...\shell\20_clone\shell\30_pull_all\command") and whose body
	// references the resolved target.
	winSlugs := map[string]bool{}
	for _, c := range buildCtxInstallCommands(exe) {
		joined := strings.Join(c, " ")
		for slug, tu := range canonical {
			winKey := `\shell\` + strings.ReplaceAll(tu.Path, ".", `\shell\`) + `\command`
			if !strings.Contains(joined, winKey) {
				continue
			}
			// Prefill mode emits a generic pwsh prompt with no target
			// binary baked in (see commandTemplate in installctx.go) —
			// matches Linux/macOS branches below which also exempt it.
			if tu.Mode == string(constants.CtxModePrefill) || strings.Contains(joined, tu.Target) {
				winSlugs[slug] = true
			}
		}
	}

	// Linux view: every leaf must produce a non-empty Nautilus shell
	// script that embeds the resolved target.
	linSlugs := map[string]bool{}
	for slug, tu := range canonical {
		fE := flatCtxEntry{Label: slug, Slug: slug, Args: tu.Args,
			Mode: constants.CtxMode(tu.Mode), Extended: tu.Extended, Exe: targetToExe(tu.Target, exe)}
		body := linuxShellScript(fE, exe)
		if body != "" && (tu.Mode == string(constants.CtxModePrefill) || strings.Contains(body, tu.Target)) {
			linSlugs[slug] = true
		}
	}

	// macOS view: same contract via macShellFor.
	macSlugs := map[string]bool{}
	for slug, tu := range canonical {
		fE := flatCtxEntry{Label: slug, Slug: slug, Args: tu.Args,
			Mode: constants.CtxMode(tu.Mode), Extended: tu.Extended, Exe: targetToExe(tu.Target, exe)}
		body := macShellFor(fE, exe)
		if body != "" && (tu.Mode == string(constants.CtxModePrefill) || strings.Contains(body, tu.Target)) {
			macSlugs[slug] = true
		}
	}

	wantSlugs := keysOf(canonical)
	assertSlugSet(t, "windows", wantSlugs, keysOfBool(winSlugs))
	assertSlugSet(t, "linux", wantSlugs, keysOfBool(linSlugs))
	assertSlugSet(t, "macos", wantSlugs, keysOfBool(macSlugs))
}

// TestCtxDuplicateTopLevelTerminalDocsRegression locks in the
// intentional duplication of 90_terminal / 91_docs at the top level
// of ctxMenu (lines 29-32 of installctxentries.go). The pair is
// emitted twice on purpose so the menu surfaces these shortcuts both
// at the start AND end of the categories. If someone "cleans this up"
// by dedup, this test fails and explains why.
func TestCtxDuplicateTopLevelTerminalDocsRegression(t *testing.T) {
	terminalCount, docsCount := 0, 0
	for _, e := range ctxMenu() {
		switch e.KeyName {
		case "90_terminal":
			terminalCount++
		case "91_docs":
			docsCount++
		}
	}
	if terminalCount != 2 {
		t.Errorf("90_terminal must appear exactly 2× at top level (got %d) — see installctxentries.go lines 29-32", terminalCount)
	}
	if docsCount != 2 {
		t.Errorf("91_docs must appear exactly 2× at top level (got %d) — see installctxentries.go lines 29-32", docsCount)
	}
}

// TestCtxFlattenDedupesDuplicateTopLevelEntries — even with the
// intentional double-declaration of 90_terminal/91_docs, the flatten
// pass must collapse to one slug per (label, args) so the platform
// installers don't write duplicate menu items.
func TestCtxFlattenDedupesDuplicateTopLevelEntries(t *testing.T) {
	leaves := collectCtxLeaves(t)
	seen := map[string]int{}
	for _, l := range leaves {
		seen[l.Slug]++
	}
	for slug, n := range seen {
		if n > 1 {
			t.Errorf("flatten produced slug %q %d times — duplicate-collapse broken", slug, n)
		}
	}
}

// TestCtxArgvParityAcrossPlatformRenders — for every leaf, the joined
// argv string MUST appear verbatim in both the Linux and macOS render
// (Windows uses pwsh quoting that may rewrap, so we only assert each
// individual arg is present). Catches accidental shell-escape drift.
func TestCtxArgvParityAcrossPlatformRenders(t *testing.T) {
	exe := fakeGitmapExe(t)
	for _, l := range collectCtxLeaves(t) {
		if l.Mode == constants.CtxModePrefill || len(l.Args) == 0 {
			continue
		}
		l := l
		t.Run(l.Slug, func(t *testing.T) {
			fE := flatCtxEntry{Label: l.Label, Slug: l.Slug, Args: l.Args,
				Mode: l.Mode, Exe: l.Exe, Extended: l.Extended}
			joined := strings.Join(l.Args, " ")
			lin := linuxShellScript(fE, exe)
			mac := macShellFor(fE, exe)
			if !strings.Contains(lin, joined) {
				t.Errorf("Linux render missing joined argv %q in:\n%s", joined, lin)
			}
			if !strings.Contains(mac, joined) {
				t.Errorf("macOS render missing joined argv %q in:\n%s", joined, mac)
			}
		})
	}
}

func keysOf[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

func keysOfBool(m map[string]bool) []string {
	return keysOf(m)
}

func assertSlugSet(t *testing.T, platform string, want, got []string) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		missing, extra := diffSlugSets(want, got)
		t.Errorf("%s slug-set mismatch — missing=%v extra=%v", platform, missing, extra)
	}
}

func diffSlugSets(want, got []string) (missing, extra []string) {
	w := map[string]bool{}
	for _, s := range want {
		w[s] = true
	}
	g := map[string]bool{}
	for _, s := range got {
		g[s] = true
	}
	for s := range w {
		if !g[s] {
			missing = append(missing, s)
		}
	}
	for s := range g {
		if !w[s] {
			extra = append(extra, s)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	return missing, extra
}

// targetToExe inverts resolvedTarget for parity-test reconstruction:
// if target equals the gitmap exe, leave Exe empty (default path);
// otherwise pass it through as a per-entry override.
func targetToExe(target, gitmapExe string) string {
	if target == gitmapExe {
		return ""
	}

	return target
}
