// Package cmdssh — ssh_hosts_sync.go synchronizes SSHConnection records and ssh_hosts records bidirectionally.
package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// syncSSHHostsFromConnections synchronizes records from SSHConnection into ssh_hosts.
func syncSSHHostsFromConnections(ctx context.Context, sqlDB *sql.DB) {
	connsRes := db.GetSSHConnections(ctx, sqlDB)
	if connsRes.IsFailure() || len(connsRes.Data) == 0 {
		return
	}
	existing, err := store.ListHosts(ctx, sqlDB)
	if err != nil {
		existing = []store.SSHHost{}
	}
	known := buildKnownHostsMap(existing)
	now := time.Now().UTC()
	for _, c := range connsRes.Data {
		syncSingleConnToHost(ctx, sqlDB, c, known, now)
	}
}

func buildKnownHostsMap(hosts []store.SSHHost) map[string]bool {
	known := make(map[string]bool, len(hosts)*2)
	for _, h := range hosts {
		known[strings.ToLower(h.Alias)] = true
		known[strings.ToLower(h.IP)] = true
	}
	return known
}

func syncSingleConnToHost(ctx context.Context, sqlDB *sql.DB, c db.SSHConnection, known map[string]bool, now time.Time) {
	aliasKey := strings.ToLower(c.Alias)
	ipKey := strings.ToLower(c.IPAddress)
	if known[aliasKey] || known[ipKey] {
		return
	}
	host := store.SSHHost{
		ID:                fmt.Sprintf("host-%s", c.IPAddress),
		Alias:             c.Alias,
		IP:                c.IPAddress,
		Username:          c.Username,
		Port:              22,
		EncryptedPassword: c.EncryptedPassword,
		ClusterRole:       "worker",
		CreatedAt:         now,
	}
	_ = store.UpsertSSHHost(ctx, host, sqlDB)
	known[aliasKey] = true
	known[ipKey] = true
}

func upsertConnToSSHHost(ctx context.Context, sqlDB *sql.DB, c db.SSHConnection, now time.Time) {
	host := store.SSHHost{
		ID:                fmt.Sprintf("host-%s", c.IPAddress),
		Alias:             c.Alias,
		IP:                c.IPAddress,
		Username:          c.Username,
		Port:              22,
		EncryptedPassword: c.EncryptedPassword,
		ClusterRole:       "worker",
		CreatedAt:         now,
	}
	_ = store.UpsertSSHHost(ctx, host, sqlDB)
}
