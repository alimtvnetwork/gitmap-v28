package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ClusterNode struct {
	NodeId         string
	Alias          string
	DisplayId      int
	IPAddress      string
	NodeRole       string
	OS             string
	JoinedAt       time.Time
	LastHeartbeat  *time.Time
	Status         string
	PasswordHash   *string
	PackageManager *string
}

func InsertOrUpdateClusterNode(ctx context.Context, db *sql.DB, node ClusterNode) error {
	query := `
		INSERT INTO ClusterNode (
			NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(NodeId) DO UPDATE SET
			Alias = excluded.Alias,
			DisplayId = excluded.DisplayId,
			IPAddress = excluded.IPAddress,
			NodeRole = excluded.NodeRole,
			OS = excluded.OS,
			LastHeartbeat = excluded.LastHeartbeat,
			Status = excluded.Status,
			PasswordHash = excluded.PasswordHash,
			PackageManager = excluded.PackageManager
	`
	_, err := db.ExecContext(ctx, query,
		node.NodeId, node.Alias, node.DisplayId, node.IPAddress, node.NodeRole, node.OS,
		node.JoinedAt, node.LastHeartbeat, node.Status, node.PasswordHash, node.PackageManager,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert ClusterNode: %w", err)
	}
	return nil
}

func ListClusterNodes(ctx context.Context, db *sql.DB) ([]ClusterNode, error) {
	query := `
		SELECT 
			NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
		FROM ClusterNode
		ORDER BY DisplayId ASC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ClusterNodes: %w", err)
	}
	defer rows.Close()

	var nodes []ClusterNode
	for rows.Next() {
		var n ClusterNode
		err := rows.Scan(
			&n.NodeId, &n.Alias, &n.DisplayId, &n.IPAddress, &n.NodeRole, &n.OS,
			&n.JoinedAt, &n.LastHeartbeat, &n.Status, &n.PasswordHash, &n.PackageManager,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ClusterNode: %w", err)
		}
		nodes = append(nodes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetClusterNode(ctx context.Context, db *sql.DB, id string) (ClusterNode, error) {
	query := `
		SELECT 
			NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
		FROM ClusterNode
		WHERE NodeId = ?
	`
	row := db.QueryRowContext(ctx, query, id)
	var n ClusterNode
	err := row.Scan(
		&n.NodeId, &n.Alias, &n.DisplayId, &n.IPAddress, &n.NodeRole, &n.OS,
		&n.JoinedAt, &n.LastHeartbeat, &n.Status, &n.PasswordHash, &n.PackageManager,
	)
	if err == sql.ErrNoRows {
		return ClusterNode{}, fmt.Errorf("ClusterNode not found: %s", id)
	}
	if err != nil {
		return ClusterNode{}, fmt.Errorf("failed to scan ClusterNode: %w", err)
	}
	return n, nil
}

func UpdateClusterNodePassword(ctx context.Context, db *sql.DB, id string, hash *string) error {
	query := `UPDATE ClusterNode SET PasswordHash = ? WHERE NodeId = ?`
	_, err := db.ExecContext(ctx, query, hash, id)
	if err != nil {
		return fmt.Errorf("failed to update password hash for node %s: %w", id, err)
	}
	return nil
}

func DeleteClusterNode(ctx context.Context, db *sql.DB, id string) error {
	query := `DELETE FROM ClusterNode WHERE NodeId = ?`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete node %s: %w", id, err)
	}
	return nil
}

func DeleteClusterRunsBefore(ctx context.Context, db *sql.DB, before time.Time) (int64, error) {
	query := `DELETE FROM ClusterRun WHERE StartedAt < ?`
	res, err := db.ExecContext(ctx, query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete ClusterRuns: %w", err)
	}
	return res.RowsAffected()
}

type ClusterStats struct {
	TotalRuns        int64
	TotalCommands    int64
	SuccessCommands  int64
	MostTargetedNode string
	MostUsedSubCmd   string
}

func GetClusterStats(ctx context.Context, db *sql.DB) (ClusterStats, error) {
	var stats ClusterStats

	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterRun`).Scan(&stats.TotalRuns)
	if err != nil {
		return stats, err
	}

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterExecResult`).Scan(&stats.TotalCommands)
	if err != nil {
		return stats, err
	}

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterExecResult WHERE ResultStatus = 1`).Scan(&stats.SuccessCommands)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = db.QueryRowContext(ctx, `SELECT NodeId FROM ClusterExecResult GROUP BY NodeId ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&stats.MostTargetedNode)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	err = db.QueryRowContext(ctx, `SELECT SubCommand FROM ClusterExecResult GROUP BY SubCommand ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&stats.MostUsedSubCmd)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	return stats, nil
}
