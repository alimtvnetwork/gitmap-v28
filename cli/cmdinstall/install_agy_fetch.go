package cmdinstall

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isHtmlContent(body []byte) bool {
	str := strings.ToLower(strings.TrimSpace(string(body)))

	return strings.HasPrefix(str, "<!") || strings.HasPrefix(str, "<html")
}

func recordToolInDatabases(tool, ver, manager string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err == nil {
		defer splitDB.Close()
		_ = splitDB.SaveInstalledTool(tool, ver, manager)
	}
}
