package helptext

// stamp.go — optional fixture-stamp validation for embedded help
// markdown. Extends the fixtureversion contract (originally built
// for test goldens) to the help corpus so a stale help file caught
// by a downstream test surfaces here with an actionable message
// instead of a low-level diff.
//
// The stamp is OPTIONAL: unstamped files are treated as "no
// contract" and pass through unchanged. New or churn-prone help
// files can opt in by adding a `<!-- fixture-stamp: ... -->` HTML
// comment at the top; the validator strips the comment wrapper and
// hands the inner marker to fixtureversion.ParseMarker.
//
// See mem://features/fixture-version-stamping for the base
// contract this file mirrors.

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/fixtureversion"
)

// stampCommentPrefix is the HTML-comment form the stamp takes
// inside a markdown help file. Kept as a constant so producers and
// this consumer cannot drift.
const stampCommentPrefix = "<!-- fixture-stamp:"

// ValidateStamp inspects the given help file body for a
// `<!-- fixture-stamp: ... -->` marker and, when present, runs it
// through fixtureversion.Validate against want. Files without a
// marker return nil (opt-in contract).
//
// Returns an error only when a marker IS present and it fails the
// staleness contract — same actionable message shape as the base
// fixtureversion package produces for test goldens.
func ValidateStamp(command string, body string, want fixtureversion.Expectation) error {
	marker, ok := extractHelpMarker(body)
	if !ok {
		return nil
	}
	stamp, ok := fixtureversion.ParseMarker(marker)
	if !ok {
		return fmt.Errorf("helptext %q has a fixture-stamp comment but it failed to parse; "+
			"expected form: %s name=<name> generation=<N> min-current=<N> for=<note> -->",
			command, stampCommentPrefix)
	}

	return fixtureversion.Validate(stamp, want)
}

// extractHelpMarker peels the HTML-comment wrapper off a stamped
// help file and returns the inner `// fixture-stamp: ...` line
// that fixtureversion.ParseMarker understands. Returns ("", false)
// when the file has no stamp comment.
func extractHelpMarker(body string) (string, bool) {
	// Only scan the first 1 KiB — stamps are always at the top,
	// and this keeps the check cheap for large help pages.
	head := body
	if len(head) > 1024 {
		head = head[:1024]
	}
	for _, line := range strings.Split(head, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, stampCommentPrefix) {
			continue
		}
		// Convert `<!-- fixture-stamp: X -->` → `// fixture-stamp: X`
		inner := strings.TrimPrefix(trimmed, "<!--")
		inner = strings.TrimSuffix(inner, "-->")
		inner = strings.TrimSpace(inner)

		return "// " + inner, true
	}

	return "", false
}
