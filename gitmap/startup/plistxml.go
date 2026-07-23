package startup

// Tiny XML token-stream helpers used only by the plist parser. Split
// out of plist.go so that file stays under the per-file budget; both
// helpers are pure (no startup-package state) and could in principle
// move to a shared util pkg, but they're scoped here to avoid
// exporting plist-shaped helpers that don't belong in a general util.

import (
	"encoding/xml"
	"strings"
)

// readElementText reads CharData until the matching end element.
// Used for <key>NAME</key> and <string>VALUE</string>. Returns the
// concatenated text trimmed of surrounding whitespace.
func readElementText(dec *xml.Decoder, start xml.StartElement) string {
	var b strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if cd, ok := tok.(xml.CharData); ok {
			b.Write(cd)
			continue
		}
		if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			break
		}
	}

	return strings.TrimSpace(b.String())
}

// readStringArray collects all <string>...</string> values inside an
// <array>...</array> until the array closes. Anything else inside
// the array (nested arrays, dicts) is ignored — LaunchAgent
// ProgramArguments is spec'd as a flat array of strings.
func readStringArray(dec *xml.Decoder) []string {
	var out []string
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if start, ok := tok.(xml.StartElement); ok && start.Name.Local == "string" {
			out = append(out, readElementText(dec, start))
			continue
		}
		if end, ok := tok.(xml.EndElement); ok && end.Name.Local == "array" {
			break
		}
	}

	return out
}
