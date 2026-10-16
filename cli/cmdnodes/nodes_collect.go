package cmdnodes

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	dbpkg "github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// UnifiedFleetNode represents a unified node aggregated across SSH, Cluster, and Server-Client (SC) networks.
type UnifiedFleetNode struct {
	Alias      string   `json:"alias"`
	Role       string   `json:"role"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	User       string   `json:"user"`
	Status     string   `json:"status"`
	Subsystems []string `json:"subsystems"`
	EnrolledAt string   `json:"enrolled_at"`
	IsOnline   bool     `json:"is_online"`
}

type nodesFilterOptions struct {
	isJSON       bool
	isFast       bool
	openOnly     bool
	isNoCache    bool
	isRefresh    bool
	filterSSH    bool
	filterClust  bool
	filterSC     bool
	targetFilter string
}

func parseNodesFilterOptions(args []string) nodesFilterOptions {
	var opts nodesFilterOptions
	for _, arg := range args {
		if arg == "--json" || arg == "-j" {
			opts.isJSON = true
			continue
		}
		if arg == "--fast" || arg == "--no-probe" || arg == "-f" {
			opts.isFast = true
			continue
		}
		if arg == "--open-only" || arg == "--running-only" || arg == "-ro" || arg == "--online-only" {
			opts.openOnly = true
			continue
		}
		if arg == "--no-cache" {
			opts.isNoCache = true
			continue
		}
		if arg == "--refresh" || arg == "-r" {
			opts.isRefresh = true
			InvalidateNodesLivenessCache()
			continue
		}
		if arg == "--ssh" {
			opts.filterSSH = true
			continue
		}
		if arg == "--cluster" || arg == "--clst" {
			opts.filterClust = true
			continue
		}
		if arg == "--sc" || arg == "--server-clients" || arg == "--servers-clients" {
			opts.filterSC = true
			continue
		}
		hasDashPrefix := strings.HasPrefix(arg, "-")
		isTargetCandidate := !hasDashPrefix && opts.targetFilter == ""
		if isTargetCandidate {
			opts.targetFilter = arg
		}
	}
	return opts
}

func runUnifiedNodesList(args []string) error {
	opts := parseNodesFilterOptions(args)
	ctx := context.Background()

	canUseCache := !opts.isNoCache && !opts.isRefresh
	nodes, err := collectUnifiedNodes(ctx, opts.isFast, canUseCache)
	if err != nil {
		return err
	}

	filtered := applyNodesFilter(nodes, opts)

	if opts.isJSON {
		return printNodesJSON(filtered)
	}

	return renderUnifiedNodesTable(os.Stdout, filtered)
}

func openUnifiedStore() (*store.DB, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "openUnifiedStore")
	}
	if err := dbConn.Migrate(); err != nil {
		dbConn.Close()
		return nil, apperror.WrapSimple(err, "openUnifiedStore.Migrate")
	}
	return dbConn, nil
}

func resolveProbeCachePreference(canUseCache []bool) bool {
	hasCacheParam := len(canUseCache) > 0
	if hasCacheParam {
		return canUseCache[0]
	}
	return true
}

func collectUnifiedNodes(ctx context.Context, isFast bool, canUseCache ...bool) ([]UnifiedFleetNode, error) {
	dbConn, err := openUnifiedStore()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	sqlDB := dbConn.SQL()
	hosts, _ := store.ListHosts(ctx, sqlDB)
	connsRes := dbpkg.GetSSHConnections(ctx, sqlDB)
	var conns []dbpkg.SSHConnection
	if connsRes.IsSuccess() {
		conns = connsRes.Data
	}
	clusterNodesRes := dbpkg.ListClusterNodes(ctx, sqlDB)
	var clusterNodes []dbpkg.ClusterNode
	if clusterNodesRes.IsSuccess() {
		clusterNodes = clusterNodesRes.Data
	}

	nodeMap := make(map[string]*UnifiedFleetNode)

	populateSSHHosts(nodeMap, hosts)
	populateSSHConnections(nodeMap, conns)
	populateClusterNodes(nodeMap, clusterNodes)

	nodes := flattenNodeMap(nodeMap)

	isProbeNeeded := !isFast
	if isProbeNeeded {
		useCache := resolveProbeCachePreference(canUseCache)
		nodes = probeUnifiedNodesLive(ctx, nodes, hosts, useCache)
	}

	return nodes, nil
}

func populateSSHHosts(m map[string]*UnifiedFleetNode, hosts []store.SSHHost) {
	for _, h := range hosts {
		key := resolveNodeKey(h.IP, h.Alias)
		m[key] = &UnifiedFleetNode{
			Alias:      resolveNodeAlias(h.Alias),
			Role:       resolveNodeRole(h.ClusterRole),
			Host:       h.IP,
			Port:       resolveNodePort(h.Port),
			User:       resolveNodeUser(h.Username),
			Status:     h.Status,
			Subsystems: []string{"SSH", "Cluster", "SC"},
			EnrolledAt: formatNodeTime(h.CreatedAt),
			IsOnline:   strings.Contains(h.Status, "ready"),
		}
	}
}

func populateSSHConnections(m map[string]*UnifiedFleetNode, conns []dbpkg.SSHConnection) {
	for _, c := range conns {
		key := resolveNodeKey(c.IPAddress, c.Alias)
		existing, ok := m[key]
		if ok {
			ensureSubsystem(existing, "SSH")
			continue
		}
		m[key] = &UnifiedFleetNode{
			Alias:      resolveNodeAlias(c.Alias),
			Role:       "worker",
			Host:       c.IPAddress,
			Port:       22,
			User:       resolveNodeUser(c.Username),
			Status:     "ready",
			Subsystems: []string{"SSH", "Cluster", "SC"},
			EnrolledAt: formatNodeTime(c.CreatedAt),
			IsOnline:   true,
		}
	}
}

func updateExistingClusterNode(existing *UnifiedFleetNode, nodeRole string) {
	ensureSubsystem(existing, "Cluster")
	ensureSubsystem(existing, "SC")
	if nodeRole != "" {
		existing.Role = nodeRole
	}
}

func populateClusterNodes(m map[string]*UnifiedFleetNode, nodes []dbpkg.ClusterNode) {
	for _, n := range nodes {
		key := resolveNodeKey(n.IPAddress, n.Alias)
		existing, ok := m[key]
		if ok {
			updateExistingClusterNode(existing, n.NodeRole)
			continue
		}
		m[key] = &UnifiedFleetNode{
			Alias:      resolveNodeAlias(n.Alias),
			Role:       resolveNodeRole(n.NodeRole),
			Host:       n.IPAddress,
			Port:       22,
			User:       "root",
			Status:     n.Status,
			Subsystems: []string{"Cluster", "SC"},
			EnrolledAt: formatNodeTime(n.JoinedAt),
			IsOnline:   strings.Contains(n.Status, "ready") || strings.Contains(n.Status, "active"),
		}
	}
}

func ensureSubsystem(node *UnifiedFleetNode, sys string) {
	for _, s := range node.Subsystems {
		if s == sys {
			return
		}
	}
	node.Subsystems = append(node.Subsystems, sys)
}

func resolveNodeKey(ip, alias string) string {
	if ip != "" {
		return strings.ToLower(ip)
	}
	return strings.ToLower(alias)
}

func resolveNodeAlias(alias string) string {
	if alias == "" {
		return "-"
	}
	return alias
}

func resolveNodeRole(role string) string {
	if role == "" {
		return "worker"
	}
	return role
}

func resolveNodePort(port int) int {
	if port <= 0 {
		return 22
	}
	return port
}

func resolveNodeUser(user string) string {
	if user == "" {
		return "root"
	}
	return user
}

func formatNodeTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

func flattenNodeMap(m map[string]*UnifiedFleetNode) []UnifiedFleetNode {
	nodes := make([]UnifiedFleetNode, 0, len(m))
	for _, node := range m {
		nodes = append(nodes, *node)
	}
	return nodes
}
