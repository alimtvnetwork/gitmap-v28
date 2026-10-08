package cmdhistory

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func openDb() (*store.DB, error) {
	return store.OpenDefault()
}

func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}

func splitOnComma(s string) []string {
	out := make([]string, 0, 4)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			out = appendTrimmedPiece(out, s, start, i)
			start = i + 1
		}
	}
	return out
}

func appendTrimmedPiece(out []string, s string, start, end int) []string {
	piece := strings.TrimSpace(s[start:end])
	if len(piece) > 0 {
		out = append(out, piece)
	}
	return out
}
