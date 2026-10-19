package cmdagy

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	transcriptCacheMu sync.RWMutex
	transcriptCache   = make(map[string]transcriptCacheEntry)
)

type transcriptCacheEntry struct {
	workspace string
	cachedAt  time.Time
}

func invalidateTranscriptCache() {
	transcriptCacheMu.Lock()
	transcriptCache = make(map[string]transcriptCacheEntry)
	transcriptCacheMu.Unlock()
}

func cleanURIStringToPath(rawURI string) string {
	trimmed := strings.TrimPrefix(rawURI, "file:///")
	decoded, err := url.PathUnescape(trimmed)
	hasErr := err != nil
	if hasErr {
		decoded = trimmed
	}
	decoded = normalizeDecodedPath(decoded, rawURI)
	clean := filepath.Clean(filepath.FromSlash(decoded))

	return strings.ToLower(clean)
}

func normalizeDecodedPath(decoded, rawURI string) string {
	isDrivePath := len(decoded) > 1 && decoded[1] == ':'
	if isDrivePath {
		return decoded
	}
	hasFilePrefix := strings.HasPrefix(rawURI, "file:///")
	if hasFilePrefix {
		return "/" + decoded
	}

	return decoded
}

func extractWorkspaceFromTranscript(convID string) string {
	transcriptCacheMu.RLock()
	entry, ok := transcriptCache[convID]
	transcriptCacheMu.RUnlock()
	if ok && time.Since(entry.cachedAt) < convCacheTTL {
		return entry.workspace
	}

	f, err := openConvTranscriptFile(convID)
	hasErr := err != nil
	if hasErr {
		return ""
	}
	defer f.Close()

	ws := scanTranscriptForWorkspace(f)

	transcriptCacheMu.Lock()
	transcriptCache[convID] = transcriptCacheEntry{
		workspace: ws,
		cachedAt:  time.Now(),
	}
	transcriptCacheMu.Unlock()

	return ws
}

func scanTranscriptForWorkspace(f *os.File) string {
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		ws := extractWorkspaceFromLine(scanner.Bytes())
		hasWs := ws != ""
		if hasWs {
			return ws
		}
	}

	return ""
}

func extractWorkspaceFromLine(line []byte) string {
	match := fileURIRegex.Find(line)
	hasMatch := len(match) > 0
	if hasMatch == false {
		return ""
	}

	uriStr := string(match)
	isInternal := isInternalAgyPath(uriStr)
	if isInternal {
		return ""
	}

	return cleanURIStringToPath(uriStr)
}

func isInternalAgyPath(uri string) bool {
	lower := strings.ToLower(uri)
	hasBrain := strings.Contains(lower, ".gemini") || strings.Contains(lower, "antigravity/brain")

	return hasBrain
}
