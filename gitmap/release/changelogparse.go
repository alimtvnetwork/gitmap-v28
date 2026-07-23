package release

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// parseChangelogStream consumes the Markdown changelog and returns the
// structured entries. It captures both the legacy flat Notes (top-level
// bullets only, trimmed) AND the new Bullets slice with depth + marker
// information for the pretty console renderer.
func parseChangelogStream(r io.Reader) ([]ChangelogEntry, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var entries []ChangelogEntry
	current := ChangelogEntry{}
	inSection := false

	for scanner.Scan() {
		raw := scanner.Text()
		if isVersionHeading(raw) {
			entries, current, inSection = startNewSection(entries, current, inSection, raw)
			continue
		}
		if !inSection {
			continue
		}
		current = appendChangelogBullet(current, raw)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if inSection {
		entries = append(entries, current)
	}

	return entries, nil
}

// isVersionHeading reports whether line starts a new "## " section.
func isVersionHeading(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "## ")
}

// startNewSection finalizes the previous section and opens a fresh one.
func startNewSection(entries []ChangelogEntry, current ChangelogEntry, inSection bool, raw string) ([]ChangelogEntry, ChangelogEntry, bool) {
	if inSection {
		entries = append(entries, current)
	}
	header := strings.TrimSpace(raw)
	version := parseVersionHeader(header)
	if len(version) == 0 {
		return entries, ChangelogEntry{}, false
	}

	return entries, ChangelogEntry{
		Version: version,
		Title:   parseVersionTitle(header),
		Notes:   []string{},
		Bullets: []ChangelogBullet{},
	}, true
}

// appendChangelogBullet parses a bullet line and appends it to current.
// Non-bullet lines (### sub-headings, blanks, paragraphs) are ignored so
// console output stays focused on the actionable list items.
func appendChangelogBullet(current ChangelogEntry, raw string) ChangelogEntry {
	depth, marker, text, ok := parseBulletLine(raw)
	if !ok {
		return current
	}

	current.Bullets = append(current.Bullets, ChangelogBullet{
		Depth:   depth,
		Ordered: isOrderedMarker(marker),
		Marker:  marker,
		Text:    text,
	})
	if depth == 0 {
		current.Notes = append(current.Notes, text)
	}

	return current
}

// parseBulletLine returns (depth, marker, text, ok) for a Markdown bullet.
// Depth is computed from the leading whitespace count divided by two
// (Markdown convention: 2 spaces per nesting level). Tabs count as 4 spaces.
func parseBulletLine(raw string) (int, string, string, bool) {
	leading, body := splitLeadingIndent(raw)
	if len(body) == 0 {
		return 0, "", "", false
	}

	marker, rest, ok := extractBulletMarker(body)
	if !ok {
		return 0, "", "", false
	}

	depth := indentToDepth(leading)
	text := strings.TrimSpace(rest)
	if len(text) == 0 {
		return 0, "", "", false
	}

	return depth, marker, text, true
}

// splitLeadingIndent returns the leading whitespace and the remaining body.
func splitLeadingIndent(raw string) (string, string) {
	for i, r := range raw {
		if r != ' ' && r != '\t' {
			return raw[:i], raw[i:]
		}
	}

	return raw, ""
}

// extractBulletMarker pulls "-", "*", or "<digits>." off the front of body.
func extractBulletMarker(body string) (string, string, bool) {
	if strings.HasPrefix(body, "- ") {
		return "-", body[2:], true
	}
	if strings.HasPrefix(body, "* ") {
		return "*", body[2:], true
	}

	return extractOrderedMarker(body)
}

// extractOrderedMarker recognizes "<digits>. " ordered list markers.
func extractOrderedMarker(body string) (string, string, bool) {
	end := 0
	for end < len(body) && unicode.IsDigit(rune(body[end])) {
		end++
	}
	if end == 0 || end+1 >= len(body) || body[end] != '.' || body[end+1] != ' ' {
		return "", "", false
	}

	return body[:end+1], body[end+2:], true
}

// isOrderedMarker reports whether marker came from an ordered list.
func isOrderedMarker(marker string) bool {
	return len(marker) > 1 && marker[len(marker)-1] == '.'
}

// indentToDepth converts leading whitespace into a nesting depth.
func indentToDepth(leading string) int {
	width := 0
	for _, r := range leading {
		if r == '\t' {
			width += 4
			continue
		}
		width++
	}

	return width / 2
}

// parseVersionTitle extracts the short title that follows the date in a
// heading like "## v3.32.1 — (2026-04-20) — Title goes here".
func parseVersionTitle(header string) string {
	raw := strings.TrimSpace(strings.TrimPrefix(header, "##"))
	parts := strings.Split(raw, "—")
	if len(parts) < 3 {
		return ""
	}

	return strings.TrimSpace(parts[len(parts)-1])
}
