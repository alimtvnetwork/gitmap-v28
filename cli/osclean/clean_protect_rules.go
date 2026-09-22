package osclean

import (
	"strings"
)

var protectedBaseNames = map[string]bool{
	"config.json":               true,
	"settings.json":             true,
	"preferences":               true,
	"local state":               true,
	"app_storage.json":          true,
	"mcp_config.json":           true,
	"pinned_projects.json":      true,
	"installation_id":           true,
	"antigravity_state.pbtxt":   true,
	"conversation_summaries.db": true,
	"agyhub_summaries_proto.pb": true,
	".updaterid":                true,
	".migrated":                 true,
	".gitkeep":                  true,
	"cookies":                   true,
	"network persistent state":  true,
	"trust tokens":              true,
	"sharedstorage":             true,
	"devtoolsactiveport":        true,
	"dips":                      true,
	"agy-node.cmd":              true,
}

var protectedDirTokens = []string{
	"/config/",
	"/conversations/",
	"/knowledge/",
	"/builtin/",
	"/plugin_data/",
	"/annotations/",
	"/local storage/",
	"/session storage/",
	"/network/",
	"/shared dictionary/",
	"/bin/",
	"/projects/",
	"/plugins/",
	"/prompts/",
	"/sidecars/",
}

var safeCacheTokens = []string{
	"/cache/",
	"/code cache/",
	"/gpucache/",
	"/dawngraphitecache/",
	"/dawnwebgpucache/",
	"/blob_storage/",
	"/logs/",
	"/crashes/",
	"/scratch/",
	"/tempmediastorage/",
	"/antigravity-updater/",
}

func isProtectedBaseName(baseName string) bool {
	low := strings.ToLower(baseName)
	if protectedBaseNames[low] {
		return true
	}
	if strings.HasPrefix(low, "conversation_summaries.db") {
		return true
	}
	if strings.HasPrefix(low, "cookies") || strings.HasPrefix(low, "trust tokens") {
		return true
	}

	return strings.HasPrefix(low, "sharedstorage") || strings.HasPrefix(low, "dips")
}

func hasProtectedDirToken(normalizedPath string) bool {
	for _, token := range protectedDirTokens {
		if strings.Contains(normalizedPath, token) {
			return true
		}
	}

	return false
}

func hasSafeCacheToken(normalizedPath string) bool {
	for _, token := range safeCacheTokens {
		if strings.Contains(normalizedPath, token) {
			return true
		}
	}

	return false
}
