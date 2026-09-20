package cmdagy

import (
	"database/sql"
	"path/filepath"
)

func resolveConvWorkspace(convID string) string {
	ws := readWorkspaceFromDB(convID)
	hasWs := ws != ""
	if hasWs {
		return ws
	}

	return extractWorkspaceFromTranscript(convID)
}

func readWorkspaceFromDB(convID string) string {
	conn, err := openConvDB(convID)
	hasErr := err != nil
	if hasErr {
		return ""
	}
	defer conn.Close()

	return extractWorkspaceFromConv(conn)
}

func openConvDB(convID string) (*sql.DB, error) {
	convDir, err := getConversationsDirPath()
	hasDirErr := err != nil
	if hasDirErr {
		return nil, err
	}
	dbPath := filepath.Join(convDir, convID+".db")
	dsn := "file:" + filepath.ToSlash(dbPath) + "?mode=ro"

	return sql.Open("sqlite", dsn)
}
