package cmdssh

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// migrateLegacySSHHistoryTables migrates legacy tables with detect-then-act checks.
func migrateLegacySSHHistoryTables(conn *sql.DB) {
	if checkTableExists(conn, "ssh_task_history") {
		migrateLegacyHistoryRows(conn)
		_ = store.ExecWrapper(conn, `DROP TABLE IF EXISTS ssh_task_history;`)
	}

	if checkTableExists(conn, "ssh_task_queue") {
		migrateLegacyQueueRows(conn)
		_ = store.ExecWrapper(conn, `DROP TABLE IF EXISTS ssh_task_queue;`)
	}
}

func migrateLegacyHistoryRows(conn *sql.DB) {
	if checkColumnExists(conn, "ssh_task_history", "forward_payload") {
		_ = store.ExecWrapper(conn, `INSERT OR IGNORE INTO TaskHistory 
			(TaskId, Section, Action, Target, ForwardPayload, InversePayload, RestoredAt, CreatedAt)
			SELECT task_id, 'ssh', action, target, forward_payload, inverse_payload, COALESCE(restored_at, ''), created_at 
			FROM ssh_task_history;`)
		return
	}

	payloadCol := "''"
	if checkColumnExists(conn, "ssh_task_history", "payload_json") {
		payloadCol = "payload_json"
	}
	query := fmt.Sprintf(`INSERT OR IGNORE INTO TaskHistory 
		(TaskId, Section, Action, Target, ForwardPayload, InversePayload, RestoredAt, CreatedAt)
		SELECT task_id, 'ssh', action, target, target, %s, COALESCE(restored_at, ''), created_at 
		FROM ssh_task_history;`, payloadCol)
	_ = store.ExecWrapper(conn, query)
}

func migrateLegacyQueueRows(conn *sql.DB) {
	if checkColumnExists(conn, "ssh_task_queue", "forward_payload") {
		_ = store.ExecWrapper(conn, `INSERT OR IGNORE INTO TaskQueue 
			(QueueId, Section, Action, Target, ForwardPayload, InversePayload, Status, CreatedAt, UpdatedAt)
			SELECT task_id, 'ssh', action, target, forward_payload, inverse_payload, status, created_at, updated_at 
			FROM ssh_task_queue;`)
		return
	}

	payloadCol := "''"
	if checkColumnExists(conn, "ssh_task_queue", "payload_json") {
		payloadCol = "payload_json"
	}
	query := fmt.Sprintf(`INSERT OR IGNORE INTO TaskQueue 
		(QueueId, Section, Action, Target, ForwardPayload, InversePayload, Status, CreatedAt, UpdatedAt)
		SELECT task_id, 'ssh', action, target, target, %s, status, created_at, updated_at 
		FROM ssh_task_queue;`, payloadCol)
	_ = store.ExecWrapper(conn, query)
}

func checkTableExists(conn *sql.DB, tableName string) bool {
	query := `SELECT 1 FROM sqlite_master WHERE type='table' AND name=? LIMIT 1`
	row := conn.QueryRow(query, tableName)
	var dummy int
	return row.Scan(&dummy) == nil
}

func checkColumnExists(conn *sql.DB, tableName, columnName string) bool {
	rows, err := conn.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt any
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			continue
		}
		if strings.EqualFold(name, columnName) {
			return true
		}
	}
	return false
}
