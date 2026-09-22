// Package cmdagy — agy_clean_cache_retention.go evaluates conversation retention candidates.
package cmdagy

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func evaluateConvRetention(keep int) ([]ConvPruneCandidate, []ConvPruneCandidate, error) {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "get db path")
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil, nil, apperror.WrapSimple(openErr, "open summaries db")
	}
	defer conn.Close()

	return queryAndPartitionConvs(conn, keep)
}

func queryAndPartitionConvs(conn *sql.DB, keep int) ([]ConvPruneCandidate, []ConvPruneCandidate, error) {
	all, err := queryAllConvCandidates(conn)
	if err != nil {
		return nil, nil, err
	}

	var kept, pruned []ConvPruneCandidate
	pinnedMap := getPinnedProjectsMap()

	for _, c := range all {
		c.IsPinned = isConvCandidatePinned(c, pinnedMap)
		if c.IsPinned || len(kept) < keep {
			kept = append(kept, c)
			continue
		}
		pruned = append(pruned, c)
	}

	return kept, pruned, nil
}

func queryAllConvCandidates(conn *sql.DB) ([]ConvPruneCandidate, error) {
	query := "SELECT conversation_id, title, project_id, last_modified_time FROM conversation_summaries ORDER BY last_modified_time DESC"
	rows, err := conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query conv summaries")
	}
	defer rows.Close()

	var result []ConvPruneCandidate
	for rows.Next() {
		var c ConvPruneCandidate
		if scanErr := rows.Scan(&c.ID, &c.Title, &c.ProjectID, &c.LastModified); scanErr == nil {
			result = append(result, c)
		}
	}
	return result, nil
}

func isConvCandidatePinned(c ConvPruneCandidate, pinnedMap map[string]bool) bool {
	if len(pinnedMap) == 0 {
		return false
	}
	if pinnedMap[c.ProjectID] {
		return true
	}
	return pinnedMap[c.ID]
}
