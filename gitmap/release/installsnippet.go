// Package release — installsnippet.go appends a pinned-version installer
// snippet to the GitHub release body so users who copy the snippet from
// the release page install EXACTLY that tag, never auto-resolving "latest".
//
// Spec: spec/04-release/08-pinned-version-install-snippet.md
package release

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// AppendPinnedInstallSnippet returns the release body with a markdown
// section appended that contains PowerShell + bash one-liners pinned to
// the given tag. Idempotent: if the body already ends with a snippet for
// the same tag, it is returned unchanged.
func AppendPinnedInstallSnippet(body, tag string) string {
	if len(tag) == 0 {
		return body
	}

	marker := fmt.Sprintf(constants.ReleaseSnippetMarker, tag)
	if strings.Contains(body, marker) {
		return body
	}

	snippet := buildPinnedInstallSnippet(tag)
	if len(body) == 0 {
		return snippet
	}

	return strings.TrimRight(body, "\n") + "\n\n" + snippet
}

// buildPinnedInstallSnippet renders the markdown block for a tag.
func buildPinnedInstallSnippet(tag string) string {
	return fmt.Sprintf(
		constants.ReleaseSnippetTemplate,
		tag, // marker
		tag, // header
		tag, // ps1 -Version
		tag, // sh --version
	)
}
