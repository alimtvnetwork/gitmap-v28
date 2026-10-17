package cmdfleet

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// FleetNode is one reachable machine in the operator's fleet.
type FleetNode struct {
	Name string
	Host string
	User string
}

// loadFleetNodes returns fleet nodes from the SSH connection registry.
// When the registry is unreachable it returns an empty slice plus a note
// instead of failing, so callers degrade gracefully.
func loadFleetNodes() ([]FleetNode, string, *appfault.AppError) {
	if envVal := strings.TrimSpace(os.Getenv("GITMAP_FLEET_NODES")); envVal != "" {
		return parseFleetEnv(envVal), "", nil
	}
	dbConn, openErr := store.OpenDefault()
	if openErr != nil {
		return nil, "node registry unreachable (store open failed)", nil
	}
	defer dbConn.Close()
	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil, "node registry unreachable (ssh connections query failed)", nil
	}
	var nodes []FleetNode
	for _, c := range res.Data {
		name := strings.TrimSpace(c.Alias)
		if name == "" {
			name = strings.TrimSpace(c.IPAddress)
		}
		if name == "" {
			continue
		}
		nodes = append(nodes, FleetNode{
			Name: name,
			Host: strings.TrimSpace(c.IPAddress),
			User: strings.TrimSpace(c.Username),
		})
	}

	return nodes, "", nil
}

// parseFleetEnv parses GITMAP_FLEET_NODES entries like
// "build-01=deploy@10.0.0.2,mac-mini=alim@mac-mini.local".
func parseFleetEnv(raw string) []FleetNode {
	var nodes []FleetNode
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		name, target, hasEq := strings.Cut(entry, "=")
		if !hasEq {
			name, target = entry, entry
		}
		user, host, hasAt := strings.Cut(strings.TrimSpace(target), "@")
		if !hasAt {
			host = user
			user = ""
		}
		nodes = append(nodes, FleetNode{
			Name: strings.TrimSpace(name),
			Host: strings.TrimSpace(host),
			User: strings.TrimSpace(user),
		})
	}

	return nodes
}

// findFleetNode matches a node by alias or host (case-insensitive).
func findFleetNode(nodes []FleetNode, token string) (FleetNode, bool) {
	for _, n := range nodes {
		if strings.EqualFold(n.Name, token) || strings.EqualFold(n.Host, token) {
			return n, true
		}
	}

	return FleetNode{}, false
}

// sshTarget renders user@host, falling back to bare host.
func sshTarget(n FleetNode) string {
	if n.User != "" {
		return n.User + "@" + n.Host
	}

	return n.Host
}
