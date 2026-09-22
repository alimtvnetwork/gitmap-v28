// Package cmdagy — agy_clean_cache_purge.go purges candidate conversations during cleanup.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func purgeCandidateConversations(pruned []ConvPruneCandidate) int {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return 0
	}

	convDir := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainDir := filepath.Join(home, ".gemini", "antigravity", "brain")
	dbPath, _ := getConversationSummariesDBPath()
	conn, _ := store.OpenSQLiteDB(dbPath)
	if conn != nil {
		defer conn.Close()
	}

	deleted := 0
	for _, c := range pruned {
		removeSingleConvFiles(convDir, brainDir, c.ID)
		if conn != nil {
			_ = deleteSummaryRecord(conn, c.ID)
		}
		deleted++
	}
	return deleted
}

// RenderCleanPreviewWithRetention renders preview including retention candidates.
func RenderCleanPreviewWithRetention(targets []AgyCacheTarget, procs []AgyProcessInfo, kept, pruned []ConvPruneCandidate) {
	RenderCleanPreview(targets, procs)
	fmt.Printf("\n%sConversation Retention Analysis:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("  • %sPreserving top %d active/pinned conversation(s)%s\n",
		constants.ColorGreen, len(kept), constants.ColorReset)
	if len(pruned) == 0 {
		fmt.Printf("  • %sNo conversations exceed retention threshold%s\n",
			constants.ColorDim, constants.ColorReset)
		return
	}

	fmt.Printf("  • %sPruning %d conversation(s) exceeding retention:%s\n",
		constants.ColorYellow, len(pruned), constants.ColorReset)
	for i, c := range pruned {
		if i >= 5 {
			fmt.Printf("    %s... and %d more conversation(s)%s\n",
				constants.ColorDim, len(pruned)-5, constants.ColorReset)
			break
		}
		fmt.Printf("    %s• %s %s(%s, %s)%s\n",
			constants.ColorDim, c.ID[:minLen(8, len(c.ID))], c.Title, c.ProjectID, c.LastModified, constants.ColorReset)
	}
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}
