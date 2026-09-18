package cmdagy

import (
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveConvWorkspace(convID string) string {
	convDir, err := getConversationsDirPath()
	if err != nil {
		return ""
	}
	dbPath := filepath.Join(convDir, convID+".db")
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return ""
	}
	defer conn.Close()

	return extractWorkspaceFromConv(conn)
}
