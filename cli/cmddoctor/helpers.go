package cmddoctor

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
)

const gofmtArgvOverhead = cmdfixrepo.GofmtArgvOverhead

func chunkPathsForGofmt(paths []string, maxCmdLen int) [][]string {
	return cmdfixrepo.ChunkPathsForGofmt(paths, maxCmdLen)
}

func batchCmdLen(batch []string) int {
	return cmdfixrepo.BatchCmdLen(batch)
}

//nolint:unused
func isWrapperActive() bool {
	return cmdsetup.IsWrapperActive()
}

//nolint:unused
func resolveSetupConfigPath(configPath string, hasConfig bool) string {
	return cmdsetup.ResolveSetupConfigPath(configPath, hasConfig)
}

// extractJSONString extracts a string value from JSON bytes by key.
//
//nolint:unused
func extractJSONString(data []byte, key string) string {
	s := string(data)
	needle := `"` + key + `"`
	idx := findKeyValue(s, needle)
	if idx < 0 {
		return ""
	}

	return extractQuotedValue(s, idx)
}

// findKeyValue finds the position after a JSON key and colon.
//
//nolint:unused
func findKeyValue(s, needle string) int {
	idx := strings.Index(s, needle)
	if idx < 0 {
		return -1
	}

	idx += len(needle)
	for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '\t') {
		idx++
	}

	return idx
}

// extractQuotedValue extracts a quoted string starting at idx.
//
//nolint:unused
func extractQuotedValue(s string, idx int) string {
	if idx >= len(s) || s[idx] != '"' {
		return ""
	}

	end := strings.Index(s[idx+1:], `"`)
	if end >= 0 {
		return s[idx+1 : idx+1+end]
	}

	return ""
}
