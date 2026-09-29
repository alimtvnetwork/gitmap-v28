// Package cmdssh — ssh_sync_nodes.go provides smart match-and-update node synchronization logic.
package cmdssh

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// NodeSyncStats records metrics for a match-and-update node synchronization operation.
type NodeSyncStats struct {
	Total     int `json:"total"`
	Matched   int `json:"matched"`
	Updated   int `json:"updated"`
	Inserted  int `json:"inserted"`
	Unchanged int `json:"unchanged"`
}

// SyncSSHConnectionsLocally merges incoming SSH connections into local Split-DB.
// Existing connections are matched by Alias or IPAddress and updated in place.
// Missing connections are inserted. Existing passwords are never overwritten by empty strings.
func SyncSSHConnectionsLocally(conns []db.SSHConnection) (NodeSyncStats, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return NodeSyncStats{}, err
	}
	defer dbConn.Close()
	_ = store.EnsureSSHTables(dbConn.SQL())

	stats := NodeSyncStats{Total: len(conns)}
	ctx := dbConn.Context()
	sqlDB := dbConn.SQL()
	now := time.Now().UTC()

	for _, incoming := range conns {
		syncSingleConnection(&stats, ctx, sqlDB, incoming, now)
	}
	return stats, nil
}

func syncSingleConnection(stats *NodeSyncStats, ctx context.Context, sqlDB *sql.DB, incoming db.SSHConnection, now time.Time) {
	existing := findExistingLocalConnection(ctx, sqlDB, incoming)
	if existing != nil {
		stats.Matched++
		applyMatchedConnectionUpdate(stats, ctx, sqlDB, *existing, incoming, now)
		return
	}

	incoming = enrichConnectionPassword(incoming)
	if err := db.InsertOrUpdateSSHConnection(ctx, sqlDB, incoming); err == nil {
		stats.Inserted++
	}
	upsertConnToSSHHost(ctx, sqlDB, incoming, now)
}

func findExistingLocalConnection(ctx context.Context, sqlDB *sql.DB, incoming db.SSHConnection) *db.SSHConnection {
	foundAlias, errAlias := db.GetSSHConnectionByAlias(ctx, sqlDB, incoming.Alias)
	if errAlias == nil && foundAlias != nil && foundAlias.Alias != "" {
		return foundAlias
	}
	if incoming.IPAddress == "" {
		return nil
	}
	foundIP, errIP := db.GetSSHConnectionByIP(ctx, sqlDB, incoming.IPAddress)
	if errIP == nil && foundIP != nil && foundIP.Alias != "" {
		return foundIP
	}
	return nil
}

func applyMatchedConnectionUpdate(stats *NodeSyncStats, ctx context.Context, sqlDB *sql.DB, existing, incoming db.SSHConnection, now time.Time) {
	isChanged, merged := mergeConnectionChanges(existing, incoming)
	if !isChanged {
		stats.Unchanged++
		return
	}
	if err := db.InsertOrUpdateSSHConnection(ctx, sqlDB, merged); err == nil {
		stats.Updated++
	}
	upsertConnToSSHHost(ctx, sqlDB, merged, now)
}

func mergeConnectionChanges(existing, incoming db.SSHConnection) (bool, db.SSHConnection) {
	merged := existing
	isChanged := false

	if incoming.IPAddress != "" && incoming.IPAddress != existing.IPAddress {
		merged.IPAddress = incoming.IPAddress
		isChanged = true
	}
	if incoming.Username != "" && incoming.Username != existing.Username {
		merged.Username = incoming.Username
		isChanged = true
	}
	if passChanged, newPass := resolveMergedPassword(existing, incoming); passChanged {
		merged.EncryptedPassword = newPass
		isChanged = true
	}
	if incoming.KeyPath != "" && incoming.KeyPath != existing.KeyPath {
		merged.KeyPath = incoming.KeyPath
		isChanged = true
	}
	if incoming.OS != "" && incoming.OS != existing.OS {
		merged.OS = incoming.OS
		isChanged = true
	}
	if incoming.OSGroup != "" && incoming.OSGroup != existing.OSGroup {
		merged.OSGroup = incoming.OSGroup
		isChanged = true
	}
	if incoming.OSVersion != "" && incoming.OSVersion != existing.OSVersion {
		merged.OSVersion = incoming.OSVersion
		isChanged = true
	}
	if incoming.BuildVersion != "" && incoming.BuildVersion != existing.BuildVersion {
		merged.BuildVersion = incoming.BuildVersion
		isChanged = true
	}

	return isChanged, merged
}

func resolveMergedPassword(existing, incoming db.SSHConnection) (bool, string) {
	if incoming.EncryptedPassword != "" && incoming.EncryptedPassword != existing.EncryptedPassword {
		incomingDec := tryDecryptCandidate(incoming.EncryptedPassword)
		existingDec := tryDecryptCandidate(existing.EncryptedPassword)
		if incomingDec != "" || existingDec == "" {
			return true, incoming.EncryptedPassword
		}
		return false, existing.EncryptedPassword
	}
	if existing.EncryptedPassword != "" && tryDecryptCandidate(existing.EncryptedPassword) != "" {
		return false, existing.EncryptedPassword
	}
	enriched := enrichConnectionPassword(existing)
	if enriched.EncryptedPassword != "" && enriched.EncryptedPassword != existing.EncryptedPassword {
		return true, enriched.EncryptedPassword
	}
	return false, existing.EncryptedPassword
}

func portableEncryptConnections(conns []db.SSHConnection) []db.SSHConnection {
	out := make([]db.SSHConnection, len(conns))
	copy(out, conns)
	for i := range out {
		out[i].EncryptedPassword = portableEncryptSingle(out[i])
	}
	return out
}

func portableEncryptSingle(c db.SSHConnection) string {
	plain := resolveCandidatePassword(c)
	if plain == "" {
		return c.EncryptedPassword
	}
	enc, err := encryptWithFallbackAES(plain)
	if err != nil {
		return c.EncryptedPassword
	}
	return enc
}
