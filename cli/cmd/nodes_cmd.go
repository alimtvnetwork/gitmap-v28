package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
		if !strings.HasPrefix(arg, "-") && opts.targetFilter == "" {
			opts.targetFilter = arg
		}
	}
	return opts
}

func isNodesHelpRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func printUnifiedNodesHelp() error {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║       gitmap nodes - Unified Fleet & Infrastructure Nodes        ║")
	fmt.Println("  ╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap nodes [flags] [target]")
	fmt.Println("    gitmap nodes help")
	fmt.Println()
	fmt.Println("  Description:")
	fmt.Println("    Aggregates and displays all registered machines across:")
	fmt.Println("      • SSH Cluster Nodes       (gitmap ssh nodes / SSHConnection)")
	fmt.Println("      • Cluster Fleet Database  (gitmap cluster nodes / ClusterNode)")
	fmt.Println("      • Server-Clients Network  (gitmap sc / servers-clients)")
	fmt.Println()
	fmt.Println("  Commands & Filtering:")
	fmt.Println("    gitmap nodes                Display all unified fleet nodes with live liveness")
	fmt.Println("    gitmap nodes <alias|ip>     Filter output to a specific node or host")
	fmt.Println("    gitmap nodes --ssh          Filter only SSH-enrolled nodes")
	fmt.Println("    gitmap nodes --cluster      Filter only Cluster DB registered nodes")
	fmt.Println("    gitmap nodes --sc           Filter only Server-Client broadcast nodes")
	fmt.Println("    gitmap nodes --fast         Skip network liveness check for instant display")
	fmt.Println("    gitmap nodes --json         Output machine-readable JSON telemetry")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    gitmap nodes")
	fmt.Println("    gitmap nodes main")
	fmt.Println("    gitmap nodes --json")
	fmt.Println("    gitmap nodes --fast")
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Println("    gitmap node, gitmap allnodes, gitmap all-nodes, gitmap fleet-nodes")
	fmt.Println()
	return nil
}

// runUnifiedNodesCLI handles `gitmap nodes [flags] [target]`.
func runUnifiedNodesCLI(args []string) error {
	if isNodesHelpRequest(args) {
		return printUnifiedNodesHelp()
	}

	opts := parseNodesFilterOptions(args)
	ctx := context.Background()

	nodes, err := collectUnifiedNodes(ctx, opts.isFast)
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

func collectUnifiedNodes(ctx context.Context, isFast bool) ([]UnifiedFleetNode, error) {
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

	if !isFast {
		nodes = probeUnifiedNodesLive(ctx, nodes, hosts)
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

func populateClusterNodes(m map[string]*UnifiedFleetNode, nodes []dbpkg.ClusterNode) {
	for _, n := range nodes {
		key := resolveNodeKey(n.IPAddress, n.Alias)
		existing, ok := m[key]
		if ok {
			ensureSubsystem(existing, "Cluster")
			ensureSubsystem(existing, "SC")
			if n.NodeRole != "" {
				existing.Role = n.NodeRole
			}
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

func probeUnifiedNodesLive(ctx context.Context, nodes []UnifiedFleetNode, hosts []store.SSHHost) []UnifiedFleetNode {
	if len(nodes) == 0 {
		return nodes
	}

	tempHosts := make([]store.SSHHost, len(nodes))
	for i, n := range nodes {
		tempHosts[i] = store.SSHHost{
			Alias:             n.Alias,
			IP:                n.Host,
			Port:              n.Port,
			Username:          n.User,
			EncryptedPassword: lookupHostPassword(n.Host, hosts),
		}
	}

	probed := cmdssh.ProbeHostsLiveStatus(ctx, tempHosts)

	for i := range nodes {
		nodes[i].Status = probed[i].Status
		nodes[i].IsOnline = strings.Contains(probed[i].Status, "ready") || strings.Contains(probed[i].Status, "●")
	}

	return nodes
}

func lookupHostPassword(ip string, hosts []store.SSHHost) string {
	for _, h := range hosts {
		if h.IP == ip && h.EncryptedPassword != "" {
			return h.EncryptedPassword
		}
	}
	return ""
}

func applyNodesFilter(nodes []UnifiedFleetNode, opts nodesFilterOptions) []UnifiedFleetNode {
	var filtered []UnifiedFleetNode
	for _, n := range nodes {
		if !matchesTargetFilter(n, opts.targetFilter) {
			continue
		}
		if !matchesSubsystemFilter(n, opts) {
			continue
		}
		filtered = append(filtered, n)
	}
	return filtered
}

func matchesTargetFilter(n UnifiedFleetNode, target string) bool {
	if target == "" {
		return true
	}
	low := strings.ToLower(target)
	return strings.Contains(strings.ToLower(n.Alias), low) || strings.Contains(strings.ToLower(n.Host), low)
}

func matchesSubsystemFilter(n UnifiedFleetNode, opts nodesFilterOptions) bool {
	hasSpecificFilter := opts.filterSSH || opts.filterClust || opts.filterSC
	if !hasSpecificFilter {
		return true
	}
	if opts.filterSSH && hasSubsystem(n, "SSH") {
		return true
	}
	if opts.filterClust && hasSubsystem(n, "Cluster") {
		return true
	}
	if opts.filterSC && hasSubsystem(n, "SC") {
		return true
	}
	return false
}

func hasSubsystem(n UnifiedFleetNode, sys string) bool {
	for _, s := range n.Subsystems {
		if strings.EqualFold(s, sys) {
			return true
		}
	}
	return false
}

func printNodesJSON(nodes []UnifiedFleetNode) error {
	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "printNodesJSON")
	}
	fmt.Println(string(data))
	return nil
}

func padCell(s string, width int) string {
	count := utf8.RuneCountInString(s)
	if count >= width {
		return s
	}
	return s + strings.Repeat(" ", width-count)
}

func renderUnifiedNodesTable(out io.Writer, nodes []UnifiedFleetNode) error {
	if len(nodes) == 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "  ● No nodes registered across SSH, Cluster, or Server-Client networks.")
		fmt.Fprintln(out, "    Enroll nodes with: gitmap ssh join <user@ip|ip> [alias]")
		fmt.Fprintln(out)
		return nil
	}

	renderTableHeaderBox(out, len(nodes))
	renderTableColumnHeadings(out)

	for _, n := range nodes {
		renderUnifiedNodeRow(out, n)
	}

	renderTableFooter(out, len(nodes))
	return nil
}

func renderTableHeaderBox(out io.Writer, count int) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(out, "║ GITMAP UNIFIED FLEET NODES (SSH, CLUSTER & SERVER-CLIENTS)                                                       ║")
	fmt.Fprintln(out, "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintf(out, "  Discovered: %d registered node(s) across SSH, Cluster DB & Server-Client (SC) networks\n\n", count)
}

func renderTableColumnHeadings(out io.Writer) {
	colAlias := padCell("ALIAS", 16)
	colRole := padCell("ROLE", 14)
	colHost := padCell("HOST (IP:PORT)", 22)
	colUser := padCell("USER", 14)
	colSys := padCell("SUBSYSTEMS", 18)
	colStatus := padCell("STATUS", 21)
	colEnrolled := padCell("ENROLLED", 19)

	fmt.Fprintf(out, "  %s%s %s %s %s %s %s %s%s\n",
		constants.ColorCyan, colAlias, colRole, colHost, colUser, colSys, colStatus, colEnrolled, constants.ColorReset)
	divider := strings.Repeat("─", 126)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
}

func renderUnifiedNodeRow(out io.Writer, n UnifiedFleetNode) {
	alias := formatCellAlias(n.Alias, 16)
	role := formatCellRole(n.Role, 14)
	hostPort := formatCellHost(fmt.Sprintf("%s:%d", n.Host, n.Port), 22)
	user := formatCellDim(n.User, 14)
	sys := formatCellSubsystems(strings.Join(n.Subsystems, ", "), 18)
	status := formatCellStatus(n.Status, 21)
	enrolled := formatCellDim(n.EnrolledAt, 19)

	fmt.Fprintf(out, "  %s %s %s %s %s %s %s\n",
		alias, role, hostPort, user, sys, status, enrolled)
}

func formatCellAlias(alias string, width int) string {
	return constants.ColorBold + constants.ColorWhite + padCell(alias, width) + constants.ColorReset
}

func formatCellRole(role string, width int) string {
	plain := padCell(role, width)
	if role == "control" || role == "master" || role == "server" {
		return constants.ColorYellow + plain + constants.ColorReset
	}
	return constants.ColorCyan + plain + constants.ColorReset
}

func formatCellHost(host string, width int) string {
	return constants.ColorWhite + padCell(host, width) + constants.ColorReset
}

func formatCellDim(val string, width int) string {
	return constants.ColorDim + padCell(val, width) + constants.ColorReset
}

func formatCellSubsystems(sys string, width int) string {
	return constants.ColorCyan + padCell(sys, width) + constants.ColorReset
}

func formatCellStatus(status string, width int) string {
	if status == "ready" || strings.HasPrefix(status, "●") {
		text := ensurePrefix(status, "●")
		return constants.ColorGreen + padCell(text, width) + constants.ColorReset
	}
	if strings.Contains(status, "offline") || strings.HasPrefix(status, "○") {
		text := ensurePrefix(status, "○")
		return constants.ColorDim + padCell(text, width) + constants.ColorReset
	}
	if strings.Contains(status, "auth failed") || strings.HasPrefix(status, "▲") {
		text := ensurePrefix(status, "▲")
		return constants.ColorYellow + padCell(text, width) + constants.ColorReset
	}
	return constants.ColorDim + padCell(status, width) + constants.ColorReset
}

func ensurePrefix(s, pfx string) string {
	if strings.HasPrefix(s, pfx) {
		return s
	}
	return pfx + " " + s
}

func renderTableFooter(out io.Writer, count int) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "==============================================================================================================================")
	fmt.Fprintln(out, " TIP: Run 'gitmap ssh <alias>' for shell access, or 'gitmap sc exec <cmd>' to broadcast across all nodes.")
	fmt.Fprintln(out, "==============================================================================================================================")
	fmt.Fprintln(out)
}
