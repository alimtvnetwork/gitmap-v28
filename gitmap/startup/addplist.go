package startup

// macOS LaunchAgent renderer for `gitmap startup-add`. Mirrors
// addrender.go (the .desktop renderer) one-for-one so the rest of
// add.go can stay OS-agnostic:
//
//   - prefixedFilenamePlist  → analog of prefixedFilename
//   - renderPlist            → analog of renderDesktop
//
// Why a separate file: keeps both renderers under the per-file
// budget, and each one only has to import the format-specific
// constants. The atomic-write helper (atomicWrite) is shared — the
// .plist body is just bytes from launchd's perspective.
//
// Marker contract: the first <key> in the dict is XGitmapManaged
// with <true/>. Same value the plist parser (parsePlistFields) keys
// on, so List/Remove pick the file up immediately after Add writes
// it. Label uses the on-disk basename (without extension) so a user
// running `launchctl list` sees the gitmap. prefix and can
// correlate it with `gitmap startup-list` output.
//
// What we deliberately DO NOT emit:
//
//   - KeepAlive / RunAtLoad combined: only RunAtLoad=true. Autostart
//     means "run once at login"; KeepAlive would respawn the
//     process forever, which is a different feature with different
//     uninstall consequences (would keep restarting after the
//     plist is removed mid-session). Users who want that can edit
//     the file by hand — gitmap stays out of the supervision
//     business.
//   - StandardOutPath / StandardErrorPath: choosing a log location
//     is policy that varies per host (~/Library/Logs vs syslog vs
//     /dev/null). Leaving them unset means launchd's default (no
//     redirection beyond what the process inherits) which is the
//     least surprising for an autostart entry that runs a CLI.
//   - launchctl load: the file is written but NOT loaded. Same
//     rationale as the package doc: `launchctl load` requires a
//     live user GUI session and would make CI / SSH provisioning
//     scripts brittle. The entry takes effect at the next login,
//     which matches user expectation for an "autostart" command.

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// prefixedFilenamePlist returns "<gitmap.prefix><clean>.plist".
// macOS LaunchAgent convention is reverse-DNS-style labels, hence
// the dot-separated `gitmap.` prefix (vs Linux's `gitmap-`). A name
// that already starts with the prefix is NOT double-prefixed, same
// rule as the .desktop variant.
func prefixedFilenamePlist(clean string) string {
	if strings.HasPrefix(clean, constants.StartupPlistPrefix) {
		return clean + constants.StartupPlistExt
	}

	return constants.StartupPlistPrefix + clean + constants.StartupPlistExt
}

// renderPlist builds a LaunchAgent plist body. Field order in the
// dict matches what `plutil -lint` and most hand-authored agents
// use: Label first (identity), then ProgramArguments (what runs),
// then RunAtLoad (when), then the gitmap marker LAST so it's
// visible at the bottom of `cat ~/Library/LaunchAgents/...` for
// quick eyeballing — same convention as the .desktop renderer.
//
// Exec is split on whitespace to populate ProgramArguments because
// launchd does NOT shell-parse a Program string the way Linux's
// Exec= line gets glued back together. Callers that need a single-
// argument command with embedded spaces should pre-quote and we
// will pass it through as-is (one element); see splitExecArgs.
func renderPlist(clean string, opts AddOptions) []byte {
	label := prefixedFilenamePlist(clean)
	label = strings.TrimSuffix(label, constants.StartupPlistExt)
	display := opts.DisplayName
	if len(display) == 0 {
		display = clean
	}
	dict := buildPlistDict(label, display, opts)

	return encodePlist(dict)
}

// plistEntry is one <key>/<value> pair in the rendered output.
// Modeled as a slice (not a map) because dict order in plist XML is
// the canonical convention readers expect — and Go map iteration is
// randomized, which would break golden-byte tests.
type plistEntry struct {
	key   string
	value plistValue
}

// plistValue is the discriminated-union of the four value shapes
// renderPlist needs. Anything else (data, date, real) is out of
// scope for an autostart entry.
type plistValue struct {
	kind     string // "string" | "bool" | "stringArray"
	str      string
	boolean  bool
	strArray []string
}

// buildPlistDict assembles the ordered key/value pairs. Comment is
// emitted as a top-of-file XML comment (NOT a dict key) because
// LaunchAgent has no standard "Comment" key — embedding it as a
// real key would surface in `launchctl list` output as garbage.
func buildPlistDict(label, display string, opts AddOptions) []plistEntry {
	out := []plistEntry{
		{"Label", plistValue{kind: "string", str: label}},
		{"ProgramArguments", plistValue{kind: "stringArray", strArray: splitExecArgs(opts.Exec)}},
		{"RunAtLoad", plistValue{kind: "bool", boolean: true}},
	}
	// WorkingDirectory is the canonical LaunchAgent key launchd
	// chdirs to before exec'ing ProgramArguments. Inserted right
	// after RunAtLoad so the "what runs / where it runs" pair sits
	// together in the rendered file.
	if len(opts.WorkingDir) > 0 {
		out = append(out, plistEntry{"WorkingDirectory", plistValue{kind: "string", str: opts.WorkingDir}})
	}
	if len(display) > 0 && display != label {
		out = append(out, plistEntry{"GitmapDisplayName", plistValue{kind: "string", str: display}})
	}
	out = append(out, plistEntry{constants.StartupPlistMarker, plistValue{kind: "bool", boolean: true}})

	return out
}

// splitExecArgs turns a free-form Exec string into the array launchd
// requires. Whitespace-split with no shell escaping — same contract
// the .desktop Exec= line documents (callers pre-quote complex
// commands). Empty input returns a single empty string so launchd
// rejects the file at load time with a clear error rather than
// gitmap silently writing an empty array (which launchd treats as
// "no program" and refuses to load).
func splitExecArgs(exec string) []string {
	parts := strings.Fields(exec)
	if len(parts) == 0 {
		return []string{""}
	}

	return parts
}

// encodePlist serializes the dict to the XML form launchd expects.
// Hand-rolled (not encoding/xml on a struct) because we need exact
// control over ordering and the doctype line — encoding/xml inserts
// no doctype and would emit elements in struct-field order, not
// plist-key order.
func encodePlist(entries []plistEntry) []byte {
	var b strings.Builder
	b.WriteString(xml.Header) // <?xml version="1.0" encoding="UTF-8"?>\n

	b.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")
	b.WriteString(`<plist version="1.0">` + "\n")
	b.WriteString("<dict>\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "  <key>%s</key>\n", xmlEscape(e.key))
		writePlistValue(&b, e.value)
	}
	b.WriteString("</dict>\n")
	b.WriteString("</plist>\n")

	return []byte(b.String())
}

// writePlistValue emits one value element with 2-space indent so
// the file passes `plutil -lint` and is human-skimmable.
func writePlistValue(b *strings.Builder, v plistValue) {
	switch v.kind {
	case "string":
		fmt.Fprintf(b, "  <string>%s</string>\n", xmlEscape(v.str))
	case "bool":
		if v.boolean {
			b.WriteString("  <true/>\n")
		} else {
			b.WriteString("  <false/>\n")
		}
	case "stringArray":
		b.WriteString("  <array>\n")
		for _, s := range v.strArray {
			fmt.Fprintf(b, "    <string>%s</string>\n", xmlEscape(s))
		}
		b.WriteString("  </array>\n")
	}
}

// xmlEscape replaces the five XML special characters. We do NOT
// reuse xml.EscapeText because it writes to an io.Writer and would
// require error-handling boilerplate at every call site for a
// renderer that can't actually fail.
func xmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)

	return r.Replace(s)
}
