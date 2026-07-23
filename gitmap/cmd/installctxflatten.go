package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// flatCtxEntry is the macOS/Linux representation of a menu item — those
// platforms do not support arbitrary nested cascades inside Finder
// Services / Nautilus scripts, so we flatten "Category ▸ Child" into a
// single labelled entry. Single source of truth: ctxMenu().
type flatCtxEntry struct {
	Label    string   // "gitmap: Release — Release next (bump minor)"
	Slug     string   // filesystem-safe id derived from label, "gitmap-release-release-next"
	Args     []string // {"release", "--bump", "minor"}
	Mode     constants.CtxMode
	Exe      string // override executable; empty => use the gitmap binary
	Extended bool   // power-user action: confirm-gated on macOS/Linux (Shift+click on Windows)
}

// flattenCtxMenu walks ctxMenu() into a flat list. Categories with
// children become "<prefix>: <Category> — <Child>"; top-level leaves
// become "<prefix>: <Label>". Order is preserved. Duplicate slugs
// (e.g. the intentionally double-listed 90_terminal / 91_docs Windows
// shortcuts) are collapsed to the FIRST occurrence so the per-leaf
// destinations on macOS/Linux (one .workflow / Nautilus script /
// Dolphin Action / Thunar <unique-id>) never collide.
func flattenCtxMenu() []flatCtxEntry {
	var out []flatCtxEntry
	seen := map[string]bool{}
	add := func(e flatCtxEntry) {
		if seen[e.Slug] {
			return
		}
		seen[e.Slug] = true
		out = append(out, e)
	}
	for _, e := range ctxMenu() {
		if len(e.Children) > 0 {
			for _, c := range e.Children {
				add(flatEntry(e.MUIVerb, c))
			}

			continue
		}
		add(flatEntry("", e))
	}

	return out
}

// flatEntry builds one flatCtxEntry from a category + child (or empty
// category for top-level leaves).
func flatEntry(category string, e ctxEntry) flatCtxEntry {
	label := constants.CtxFlatPrefix + constants.CtxFlatSeparator
	if category != "" {
		label += category + constants.CtxFlatChildJoiner
	}
	label += e.MUIVerb

	return flatCtxEntry{
		Label:    label,
		Slug:     slugifyCtx(label),
		Args:     append([]string(nil), e.Args...),
		Mode:     e.Mode,
		Exe:      e.Exe,
		Extended: e.Extended,
	}
}

// slugifyCtx returns a filesystem-safe id: lowercase alphanumerics
// joined by "-". Used as workflow folder, .desktop file, and Nautilus
// script base names. Common disambiguating glyphs (`+`, `#`) are
// transliterated to letters BEFORE the strip pass so labels like
// "C++ projects" vs "C# projects" do not collide on the same slug.
func slugifyCtx(s string) string {
	s = strings.ReplaceAll(s, "++", "pp")
	s = strings.ReplaceAll(s, "+", "p")
	s = strings.ReplaceAll(s, "#", "sharp")
	out := make([]byte, 0, len(s))
	dashed := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
			dashed = false
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			out = append(out, c)
			dashed = false
		default:
			if !dashed && len(out) > 0 {
				out = append(out, '-')
				dashed = true
			}
		}
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}

	return string(out)
}
