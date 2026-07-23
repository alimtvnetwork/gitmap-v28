package startup

// macOS LaunchAgent reader. Mirrors the .desktop reader (one cheap
// filename pre-filter, one in-file marker re-check) but speaks
// Apple's plist XML grammar. Kept in its own file so the .desktop
// path stays untouched on Linux and the per-file budget holds.
//
// Plist parsing strategy: encoding/xml against a minimal schema
// (top-level <plist><dict>) instead of pulling in a third-party
// plist library. LaunchAgent plists are small (typically <2 KiB)
// and well-formed XML in practice; the binary plist format is rare
// for hand-authored agents and is intentionally not supported here
// — gitmap-MANAGED agents are always written by us in XML form, so
// a binary plist with our prefix is by definition not ours.
//
// Parsing model: we walk the dict's alternating <key>...</key> /
// value sequence as flat XML tokens. We look for two keys:
//
//   - XGitmapManaged → must be followed by <true/> for the file to
//     count as gitmap-managed.
//   - ProgramArguments → array of strings, joined with spaces for
//     display. If absent we fall back to Program (single string).

import (
	"encoding/xml"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// collectManagedPlist is the macOS analog of collectManagedDesktop.
// Same two-gate filter shape; different per-file reader.
func collectManagedPlist(dir string, files []os.DirEntry) []Entry {
	var out []Entry
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !looksLikeOursPlist(name) {
			continue
		}
		entry, ok := readManagedPlist(dir, name)
		if !ok {
			continue
		}
		out = append(out, entry)
	}

	return out
}

// looksLikeOursPlist is the cheap pre-filter for macOS: filename must
// end in `.plist` AND start with the gitmap. prefix. Same spoofing
// caveat as Linux — the marker check below is the real authority.
func looksLikeOursPlist(filename string) bool {
	if !strings.HasSuffix(filename, constants.StartupPlistExt) {
		return false
	}

	return strings.HasPrefix(filename, constants.StartupPlistPrefix)
}

// readManagedPlist opens the file and runs the marker + Exec parse.
// Returns ok=false on any I/O error or missing marker — both mean
// "not ours, skip it" from the caller's perspective. Identical
// contract to readManagedDesktop so collectManaged callers don't
// need to know which OS produced the entry.
func readManagedPlist(dir, filename string) (Entry, bool) {
	full := joinPath(dir, filename)
	f, err := os.Open(full)
	if err != nil {
		return Entry{}, false
	}
	defer f.Close()

	managed, exec := parsePlistFields(f)
	if !managed {
		return Entry{}, false
	}

	return Entry{
		Name: strings.TrimSuffix(filename, constants.StartupPlistExt),
		Path: full,
		Exec: exec,
	}, true
}

// parsePlistFields walks the plist XML token stream. Returns
// (managed, exec). We do NOT validate full plist DOCTYPE — the
// marker key + value being present is sufficient proof of
// gitmap-authored intent, and a strict DOCTYPE check would reject
// hand-edited but valid plists.
func parsePlistFields(r io.Reader) (bool, string) {
	dec := xml.NewDecoder(r)
	state := plistParseState{decoder: dec}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		state.consume(tok)
	}

	return state.managed, state.execLine()
}

// plistParseState carries the streaming parser's accumulators. Kept
// as a method receiver (not free vars) so consume() reads cleanly
// without a long parameter list.
type plistParseState struct {
	decoder *xml.Decoder
	// pendingKey holds the most recent <key>NAME</key> text so the
	// next non-key element knows which dict key it belongs to.
	pendingKey string
	managed    bool
	program    string
	progArgs   []string
}

// consume advances the parser by one token. The plist grammar we
// care about is simple enough that a flat key→value matcher works:
// keys we don't recognize are ignored, value elements unrelated to a
// pending key are dropped on the floor.
func (s *plistParseState) consume(tok xml.Token) {
	start, ok := tok.(xml.StartElement)
	if !ok {
		return
	}
	switch start.Name.Local {
	case "key":
		s.pendingKey = readElementText(s.decoder, start)
	case "true":
		if s.pendingKey == constants.StartupPlistMarker {
			s.managed = true
		}
		s.pendingKey = ""
	case "string":
		s.handleString(readElementText(s.decoder, start))
	case "array":
		if s.pendingKey == "ProgramArguments" {
			s.progArgs = readStringArray(s.decoder)
			s.pendingKey = ""
		}
	}
}

// handleString routes <string> values to the right field based on
// the pending key. Unknown keys cause the value to be discarded.
func (s *plistParseState) handleString(val string) {
	if s.pendingKey == "Program" {
		s.program = val
	}
	s.pendingKey = ""
}

// execLine joins ProgramArguments with spaces (the canonical display
// form) or falls back to Program. Both empty → empty string, which
// the renderer turns into "(no Exec line)".
func (s *plistParseState) execLine() string {
	if len(s.progArgs) > 0 {
		return strings.Join(s.progArgs, " ")
	}

	return s.program
}

// XML token-stream helpers (readElementText / readStringArray) live
// in plistxml.go to keep this file under the per-file budget.
