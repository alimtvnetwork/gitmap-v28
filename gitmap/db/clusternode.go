package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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

type ClusterStats struct {
	TotalRuns        int64
	TotalCommands    int64
	SuccessCommands  int64
	MostTargetedNode string
	MostUsedSubCmd   string
}

const (
	sqlUpsertClusterNode = `
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
	sqlSelectListClusterNodes = `
		SELECT 
			NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
		FROM ClusterNode
		ORDER BY DisplayId ASC
	`
	sqlSelectClusterNodeById = `
		SELECT 
			NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, JoinedAt, LastHeartbeat, Status, PasswordHash, PackageManager
		FROM ClusterNode
		WHERE NodeId = ?
	`
)

func InsertOrUpdateClusterNode(ctx context.Context, db *sql.DB, node ClusterNode) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlUpsertClusterNode,
		node.NodeId, node.Alias, node.DisplayId, node.IPAddress, node.NodeRole, node.OS,
		node.JoinedAt, node.LastHeartbeat, node.Status, node.PasswordHash, node.PackageManager,
	)
	if err != nil {
		return apperror.WrapSimple(err, "InsertOrUpdateClusterNode.Exec")
	}

	return nil
}

func ListClusterNodes(ctx context.Context, db *sql.DB) ([]ClusterNode, *apperror.AppError) {
	rows, err := db.QueryContext(ctx, sqlSelectListClusterNodes)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListClusterNodes.Query")
	}

	defer rows.Close()

	return scanClusterNodeRows(rows)
}

func scanClusterNodeRows(rows *sql.Rows) ([]ClusterNode, *apperror.AppError) {
	var nodes []ClusterNode
	for rows.Next() {
		var n ClusterNode
		err := rows.Scan(
			&n.NodeId, &n.Alias, &n.DisplayId, &n.IPAddress, &n.NodeRole, &n.OS,
			&n.JoinedAt, &n.LastHeartbeat, &n.Status, &n.PasswordHash, &n.PackageManager,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scanClusterNodeRows.Scan")
		}

		nodes = append(nodes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "scanClusterNodeRows.Rows")
	}

	return nodes, nil
}

func GetClusterNode(ctx context.Context, db *sql.DB, id string) (ClusterNode, *apperror.AppError) {
	row := db.QueryRowContext(ctx, sqlSelectClusterNodeById, id)
	var n ClusterNode
	err := row.Scan(
		&n.NodeId, &n.Alias, &n.DisplayId, &n.IPAddress, &n.NodeRole, &n.OS,
		&n.JoinedAt, &n.LastHeartbeat, &n.Status, &n.PasswordHash, &n.PackageManager,
	)
	if err == sql.ErrNoRows {
		return ClusterNode{}, apperror.NewWithDetails(
			"GetClusterNode", "E4004", "ClusterNode not found: "+id, "db",
			apperror.ErrorTypeNotFound, apperror.SeverityError, map[string]any{"id": id},
		)
	}

	if err != nil {
		return ClusterNode{}, apperror.WrapSimple(err, "GetClusterNode.Scan")
	}

	return n, nil
}

func UpdateClusterNodePassword(ctx context.Context, db *sql.DB, id string, hash *string) *apperror.AppError {
	query := `UPDATE ClusterNode SET PasswordHash = ? WHERE NodeId = ?`
	_, err := db.ExecContext(ctx, query, hash, id)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateClusterNodePassword.Exec")
	}

	return nil
}

func DeleteClusterNode(ctx context.Context, db *sql.DB, id string) *apperror.AppError {
	query := `DELETE FROM ClusterNode WHERE NodeId = ?`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteClusterNode.Exec")
	}

	return nil
}

func DeleteClusterRunsBefore(ctx context.Context, db *sql.DB, before time.Time) (int64, *apperror.AppError) {
	query := `DELETE FROM ClusterRun WHERE StartedAt < ?`
	res, err := db.ExecContext(ctx, query, before)
	if err != nil {
		return 0, apperror.WrapSimple(err, "DeleteClusterRunsBefore.Exec")
	}

	affected, errRows := res.RowsAffected()
	if errRows != nil {
		return 0, apperror.WrapSimple(errRows, "DeleteClusterRunsBefore.RowsAffected")
	}

	return affected, nil
}

func GetClusterStats(ctx context.Context, db *sql.DB) (ClusterStats, *apperror.AppError) {
	var stats ClusterStats
	if countErr := queryClusterCounts(ctx, db, &stats); countErr != nil {
		return stats, countErr
	}

	if aggErr := queryClusterAggregates(ctx, db, &stats); aggErr != nil {
		return stats, aggErr
	}

	return stats, nil
}

func queryClusterCounts(ctx context.Context, db *sql.DB, stats *ClusterStats) *apperror.AppError {
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterRun`).Scan(&stats.TotalRuns)
	if err != nil {
		return apperror.WrapSimple(err, "GetClusterStats.TotalRuns")
	}

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterExecResult`).Scan(&stats.TotalCommands)
	if err != nil {
		return apperror.WrapSimple(err, "GetClusterStats.TotalCommands")
	}

	return nil
}

func queryClusterAggregates(ctx context.Context, db *sql.DB, stats *ClusterStats) *apperror.AppError {
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ClusterExecResult WHERE ResultStatus = 1`).Scan(&stats.SuccessCommands)
	if err != nil && err != sql.ErrNoRows {
		return apperror.WrapSimple(err, "GetClusterStats.SuccessCommands")
	}

	err = db.QueryRowContext(ctx, `SELECT NodeId FROM ClusterExecResult GROUP BY NodeId ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&stats.MostTargetedNode)
	if err != nil && err != sql.ErrNoRows {
		return apperror.WrapSimple(err, "GetClusterStats.MostTargetedNode")
	}

	err = db.QueryRowContext(ctx, `SELECT SubCommand FROM ClusterExecResult GROUP BY SubCommand ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&stats.MostUsedSubCmd)
	if err != nil && err != sql.ErrNoRows {
		return apperror.WrapSimple(err, "GetClusterStats.MostUsedSubCmd")
	}

	return nil
}
