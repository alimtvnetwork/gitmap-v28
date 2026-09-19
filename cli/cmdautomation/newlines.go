package cmdautomation

import (
	"bytes"
	"strings"
)

var defaultPolyglotExts = map[string]bool{
	".ts": true, ".tsx": true, ".js": true, ".jsx": true, ".mjs": true, ".cjs": true,
	".go": true, ".rs": true, ".cs": true, ".java": true, ".kt": true, ".cpp": true,
	".c": true, ".h": true, ".hpp": true, ".py": true, ".php": true, ".rb": true,
	".sh": true, ".bash": true, ".zsh": true, ".ps1": true, ".bat": true,
	".md": true, ".markdown": true, ".json": true, ".yaml": true, ".yml": true,
	".toml": true, ".sql": true, ".xml": true, ".html": true, ".css": true, ".scss": true,
}

// IsPolyglotTextExtension checks if an extension is a known text/source format.
func IsPolyglotTextExtension(ext string) bool {
	low := strings.ToLower(ext)
	return defaultPolyglotExts[low]
}

// HasBinaryContent probes the first 8KB chunk for null bytes.
func HasBinaryContent(data []byte) bool {
	probeLimit := 8192
	isShort := len(data) < probeLimit
	if isShort {
		probeLimit = len(data)
	}

	return bytes.IndexByte(data[:probeLimit], 0x00) != -1
}

// StripUtf8Bom removes leading UTF-8 BOM bytes if present.
func StripUtf8Bom(data []byte) []byte {
	bom := []byte{0xEF, 0xBB, 0xBF}
	hasBom := bytes.HasPrefix(data, bom)
	if hasBom {
		return data[len(bom):]
	}

	return data
}

// NormalizeContent converts CRLF to LF, trims line whitespace, and enforces 1 trailing LF.
func NormalizeContent(data []byte) ([]byte, int, bool) {
	stripped := StripUtf8Bom(data)
	crlfCount := bytes.Count(stripped, []byte("\r\n"))
	unified := bytes.ReplaceAll(stripped, []byte("\r\n"), []byte("\n"))
	unified = bytes.ReplaceAll(unified, []byte("\r"), []byte("\n"))

	cleaned := cleanTrailingLines(unified)
	isModified := !bytes.Equal(data, cleaned)

	return cleaned, crlfCount, isModified
}

func cleanTrailingLines(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	var outLines [][]byte
	for _, l := range lines {
		outLines = append(outLines, bytes.TrimRight(l, " \t"))
	}

	for len(outLines) > 0 && len(outLines[len(outLines)-1]) == 0 {
		outLines = outLines[:len(outLines)-1]
	}

	if len(outLines) == 0 {
		return []byte{}
	}

	return append(bytes.Join(outLines, []byte("\n")), '\n')
}
