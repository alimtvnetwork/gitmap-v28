package vscodepm

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// DetectTagsCustom is the env-aware sibling of DetectTags. It runs
// the same shallow filesystem inspection but layers three user
// customizations on top, sourced from env vars set by the global
// CLI flag stripper (cmd/vscodecustomtags.go):
//
//  1. Extra marker→tag rules (GITMAP_VSCODE_TAG_MARKER) extend the
//     marker map for this call only — built-in rules win on conflict.
//  2. Skip set (GITMAP_VSCODE_TAG_SKIP) drops named tags from the
//     detected set BEFORE always-add tags are appended, so a user
//     can simultaneously skip "git" and force-add "git" if they
//     want every entry to carry it regardless of detection.
//  3. Always-add set (GITMAP_VSCODE_TAG_ADD) is appended last and
//     union'd with whatever survived steps 1-2.
//
// Output order: detected tags in canonical order first, then the
// always-add tags in user-supplied order. Determinism guaranteed
// for any given env state.
func DetectTagsCustom(rootPath string) []string {
	detected := detectTagsWithExtraMarkers(rootPath, parseMarkerEnv())
	branded := prependGitmapBrand(detected)
	filtered := dropSkipped(branded, parseListEnv(constants.EnvVSCodeTagSkip))

	return appendAlwaysAdd(filtered, parseListEnv(constants.EnvVSCodeTagAdd))
}

// prependGitmapBrand inserts the canonical "gitmap" brand tag at the
// head of the tag list so every projects.json entry written by
// gitmap is self-identifying in the VS Code Project Manager UI. The
// tag is added BEFORE the skip filter runs, so users who genuinely
// don't want it can opt out via `--vscode-tag-skip gitmap`. If the
// tag already exists (e.g. carried in via --vscode-tag gitmap or
// already on disk via unionTags upstream) it is not duplicated.
func prependGitmapBrand(tags []string) []string {
	for _, t := range tags {
		if t == constants.AutoTagGitmap {
			return tags
		}
	}

	return append([]string{constants.AutoTagGitmap}, tags...)
}

// detectTagsWithExtraMarkers runs the same scan as DetectTags but
// with an extended marker map. Extra rules whose tag collides with
// a built-in (same key) are ignored — built-in canonical order wins.
func detectTagsWithExtraMarkers(rootPath string, extra map[string]string) []string {
	if rootPath == "" {
		return nil
	}
	info, err := os.Stat(rootPath)
	if err != nil || !info.IsDir() {
		return nil
	}

	hits := map[string]struct{}{}
	customOrder := make([]string, 0)
	for marker, tag := range constants.AutoTagMarkers {
		if markerExists(rootPath, marker) {
			hits[tag] = struct{}{}
		}
	}
	for marker, tag := range extra {
		if _, builtin := constants.AutoTagMarkers[marker]; builtin {
			continue
		}
		if markerExists(rootPath, marker) {
			if _, dup := hits[tag]; !dup {
				hits[tag] = struct{}{}
				customOrder = append(customOrder, tag)
			}
		}
	}

	return append(tagsInCanonicalOrder(hits), customOnly(hits, customOrder)...)
}

// customOnly returns the user-marker tags that aren't already in the
// canonical built-in list, preserving first-seen order so two runs
// with the same env produce identical output.
func customOnly(hits map[string]struct{}, order []string) []string {
	known := map[string]struct{}{}
	for _, t := range constants.AutoTagOrder {
		known[t] = struct{}{}
	}
	out := make([]string, 0, len(order))
	seen := map[string]struct{}{}
	for _, t := range order {
		if _, isBuiltin := known[t]; isBuiltin {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		if _, ok := hits[t]; !ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	return out
}

// dropSkipped removes any tag in the skip set from tags.
func dropSkipped(tags, skip []string) []string {
	if len(skip) == 0 {
		return tags
	}
	skipSet := map[string]struct{}{}
	for _, s := range skip {
		skipSet[s] = struct{}{}
	}
	out := tags[:0:0]
	for _, t := range tags {
		if _, drop := skipSet[t]; drop {
			continue
		}
		out = append(out, t)
	}

	return out
}

// appendAlwaysAdd union-appends add to base, preserving base order.
func appendAlwaysAdd(base, add []string) []string {
	if len(add) == 0 {
		return base
	}
	seen := map[string]struct{}{}
	for _, t := range base {
		seen[t] = struct{}{}
	}
	out := append([]string{}, base...)
	for _, t := range add {
		if t == "" {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	return out
}

// parseListEnv reads an env var holding ASCII-unit-separator-joined
// tokens and returns the trimmed, non-empty list. Empty env → nil.
func parseListEnv(name string) []string {
	raw := os.Getenv(name)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, constants.EnvVSCodeTagSeparator)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

// parseMarkerEnv reads GITMAP_VSCODE_TAG_MARKER and returns the
// marker→tag map. Malformed entries (missing `=`, empty key/value)
// are silently skipped — bad CLI input must never break a sync.
func parseMarkerEnv() map[string]string {
	raw := parseListEnv(constants.EnvVSCodeTagMarker)
	if len(raw) == 0 {
		return nil
	}
	out := make(map[string]string, len(raw))
	for _, kv := range raw {
		k, v, ok := strings.Cut(kv, constants.TagMarkerKVSeparator)
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if !ok || k == "" || v == "" {
			continue
		}
		out[k] = v
	}

	return out
}
