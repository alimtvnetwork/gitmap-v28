package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func loadSSHConnectionsForTarget(target string) ([]db.SSHConnection, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf("open store db: %w", err)
	}
	defer dbConn.Close()

	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil, fmt.Errorf("retrieve ssh connections: %w", res.AppError())
	}

	return filterConnectionsByTarget(res.Data, target), nil
}

func filterConnectionsByTarget(conns []db.SSHConnection, target string) []db.SSHConnection {
	if target == "" || target == "all" {
		return conns
	}

	targets := ParseMultiIPList(target)
	if len(targets) <= 1 {
		return filterSingleTarget(conns, target)
	}

	return filterMultipleTargets(conns, targets)
}

func filterSingleTarget(conns []db.SSHConnection, target string) []db.SSHConnection {
	var matched []db.SSHConnection
	for _, c := range conns {
		if strings.EqualFold(c.Alias, target) || strings.EqualFold(c.IPAddress, target) {
			matched = append(matched, c)
		}
	}
	return matched
}

func filterMultipleTargets(conns []db.SSHConnection, targets []string) []db.SSHConnection {
	var matched []db.SSHConnection
	for _, c := range conns {
		if matchesAnyTarget(c, targets) {
			matched = append(matched, c)
		}
	}
	return matched
}

func matchesAnyTarget(c db.SSHConnection, targets []string) bool {
	for _, t := range targets {
		if strings.EqualFold(c.Alias, t) || strings.EqualFold(c.IPAddress, t) {
			return true
		}
	}
	return false
}

