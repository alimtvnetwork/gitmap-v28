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

	conns := filterConnectionsByTarget(res.Data, target)
	if target == "" || target == "all" {
		return conns, nil
	}
	return augmentMissingTargets(dbConn, target, conns), nil
}

func isTargetMatched(token string, conns []db.SSHConnection) bool {
	for _, c := range conns {
		if strings.EqualFold(c.Alias, token) || strings.EqualFold(c.IPAddress, token) {
			return true
		}
	}
	return false
}

func augmentMissingTargets(dbConn *store.DB, rawTarget string, existing []db.SSHConnection) []db.SSHConnection {
	targets := ParseMultiIPList(rawTarget)
	result := append([]db.SSHConnection{}, existing...)
	for _, token := range targets {
		if isTargetMatched(token, result) {
			continue
		}
		conn, isFound := resolveSingleMissingTarget(dbConn, token)
		if isFound {
			result = append(result, conn)
		}
	}
	return result
}

func findHostInStore(dbConn *store.DB, token string) (store.SSHHost, bool) {
	ctx, sqlDB := dbConn.Context(), dbConn.SQL()
	if host, err := store.GetHostByID(ctx, token, sqlDB); err == nil {
		return host, true
	}
	if host, err := store.GetHostByAlias(ctx, token, sqlDB); err == nil {
		return host, true
	}
	if host, err := store.GetHostByIP(ctx, token, sqlDB); err == nil {
		return host, true
	}
	return store.SSHHost{}, false
}

func convertHostToConnection(h store.SSHHost) db.SSHConnection {
	user := h.Username
	if user == "" {
		user = "root"
	}
	return db.SSHConnection{
		Alias:             h.Alias,
		IPAddress:         h.IP,
		Username:          user,
		EncryptedPassword: h.EncryptedPassword,
		OS:                "linux",
	}
}

func parseAdHocConnection(token string) (db.SSHConnection, bool) {
	t, err := ParseSSHTarget(token, "root", 22)
	isParsed := err == nil && t != nil && t.IP != ""
	if isParsed == false {
		return db.SSHConnection{}, false
	}
	user := t.Username
	if user == "" {
		user = "root"
	}
	return db.SSHConnection{
		Alias:     token,
		IPAddress: t.IP,
		Username:  user,
		OS:        "linux",
	}, true
}

func resolveSingleMissingTarget(dbConn *store.DB, token string) (db.SSHConnection, bool) {
	host, isHost := findHostInStore(dbConn, token)
	if isHost {
		return convertHostToConnection(host), true
	}
	return parseAdHocConnection(token)
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
