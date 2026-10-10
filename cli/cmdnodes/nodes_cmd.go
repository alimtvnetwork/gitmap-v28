package cmdnodes

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpullerror"
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
	fmt.Println("    gitmap nodes                          Display all unified fleet nodes with live liveness")
	fmt.Println("    gitmap nodes ping [target]            Run machine ping command against all fleet nodes")
	fmt.Println("    gitmap nodes clone [flags] <targets> [dest]   Clone repository/manifest locally and across fleet async")
	fmt.Println("    gitmap nodes clone except-self <targets>      Clone strictly to remote fleet nodes (skipping local)")
	fmt.Println("    gitmap nodes cfr [flags] [targets] [dest]     Clone, fix, and auto-setup across all fleet nodes")
	fmt.Println("    gitmap nodes cfrp [flags] [targets] [dest]    Clone, fix, and promote public across all fleet nodes")
	fmt.Println("    gitmap nodes push-settings <node>     Export and push Antigravity settings to target node")
	fmt.Println("    gitmap nodes sync-settings            Broadcast Antigravity settings across all fleet nodes")
	fmt.Println("    gitmap nodes deploy agm-accounts [flags] Broadcast AGM accounts & credentials across fleet nodes")
	fmt.Println("    gitmap nodes sync-agm-accounts [flags]   Alias for deploy agm-accounts")
	fmt.Println("    gitmap nodes deploy repo <slug> [flags]  Deploy repository & register across VS Code, Cursor, Antigravity, GitHub Desktop")
	fmt.Println("    gitmap nodes deploy repos [targets]      Batch deploy multiple repositories across fleet nodes")
	fmt.Println("    gitmap nodes scan [target] [flags]       Broadcast remote repository scanner across fleet nodes")
	fmt.Println("    gitmap nodes rescan [target] [flags]     Broadcast remote repository rescan across fleet nodes")
	fmt.Println("    gitmap nodes send-projects [node]     Forward VS Code / Cursor Project Manager workspaces")
	fmt.Println("    gitmap nodes agy prompt <node> <p>    Dispatch prompt directly to target project on node")
	fmt.Println("    gitmap nodes agy query [--ssh]        Query active Antigravity instances and prompts")
	fmt.Println("    gitmap nodes agy [ui]                 Launch Antigravity web studio dashboard")
	fmt.Println("    gitmap nodes pending-commits [flags]     Inspect uncommitted & unpushed changes across fleet")
	fmt.Println("    gitmap nodes pc [flags]                  Alias for pending-commits")
	fmt.Println("    gitmap nodes commits <target> <msg>      Standard atomic commit across fleet nodes")
	fmt.Println("    gitmap nodes cpf <target> <msg>          Feature commit (Feature: <msg>) across fleet nodes")
	fmt.Println("    gitmap nodes cpb <target> <msg>          Bug fix commit (Bug: <msg>) across fleet nodes")
	fmt.Println("    gitmap nodes cpr <target> <msg>          Release commit (Release: <msg>) across fleet nodes")
	fmt.Println("    gitmap nodes commit-fix <target> <msg>   Fix commit (Fix: <msg>) across fleet nodes")
	fmt.Println("    gitmap nodes <alias|ip>               Filter output to a specific node or host")
	fmt.Println("    gitmap nodes --ssh                    Filter only SSH-enrolled nodes")
	fmt.Println("    gitmap nodes --cluster                Filter only Cluster DB registered nodes")
	fmt.Println("    gitmap nodes --sc                     Filter only Server-Client broadcast nodes")
	fmt.Println("    gitmap nodes --fast                   Skip network liveness check for instant display")
	fmt.Println("    gitmap nodes --json                   Output machine-readable JSON telemetry")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    gitmap nodes")
	fmt.Println("    gitmap nodes pc")
	fmt.Println("    gitmap nodes pc --sort=count --detail")
	fmt.Println("    gitmap nodes cpf gitmap \"add nodes commit suite\"")
	fmt.Println("    gitmap nodes cpb all \"fix nil pointer exception\"")
	fmt.Println("    gitmap nodes commit-fix my-service \"resolve merge conflict\"")
	fmt.Println("    gitmap nodes ping")
	fmt.Println("    gitmap nodes push-settings worker-1")
	fmt.Println("    gitmap nodes sync-settings")
	fmt.Println("    gitmap nodes deploy agm-accounts")
	fmt.Println("    gitmap nodes deploy agm-accounts --target worker-1")
	fmt.Println("    gitmap nodes deploy agm-accounts --except worker-3 --dry-run")
	fmt.Println("    gitmap nodes sync-agm-accounts --include-main")
	fmt.Println("    gitmap nodes deploy repo gitmap --dry-run")
	fmt.Println("    gitmap nodes deploy repo my-app --target worker-1 --with-pinned")
	fmt.Println("    gitmap nodes deploy repo api-server --with-conversations")
	fmt.Println("    gitmap nodes scan --open-only")
	fmt.Println("    gitmap nodes send-projects worker-1")
	fmt.Println("    gitmap nodes agy prompt worker-1 gitmap \"Run test suite\"")
	fmt.Println("    gitmap nodes agy query --ssh")
	fmt.Println("    gitmap nodes agy ui")
	fmt.Println("    gitmap nodes clone ChrisTitusTech/winutil")
	fmt.Println("    gitmap nodes clone https://github.com/user/repo D:\\work\\custom")
	fmt.Println("    gitmap nodes clone except-self https://github.com/user/repo")
	fmt.Println("    gitmap nodes clone repo1,repo2")
	fmt.Println("    gitmap nodes cfr gitmap.json")
	fmt.Println("    gitmap nodes cfr")
	fmt.Println("    gitmap nodes cfrp gitmap.json")
	fmt.Println("    gitmap nodes --json")
	fmt.Println("    gitmap nodes --fast")
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Println("    gitmap node, gitmap allnodes, gitmap all-nodes, gitmap fleet-nodes")
	fmt.Println()
	return nil
}

func isNodesPingCommand(arg string) bool {
	low := strings.ToLower(arg)
	return low == "ping" || low == "probe"
}

func isNodesCloneRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	_, isClone := IsNodesCloneCommand(args[0])
	return isClone
}

func runNodesHistoryCLI(args []string) error {
	return cmdssh.RunFleetPASCommand("nodes history", "gitmap history --limit 10", func() error {
		return cmdhistory.RunHistory(args)
	})
}

func runNodesErrorsCLI(args []string) error {
	cmdStr := "gitmap errors " + strings.Join(args, " ")
	cmdStr = strings.TrimSpace(cmdStr)

	return cmdssh.RunFleetPASCommand("nodes errors", cmdStr, func() error {
		return cmderrors.RunErrorsCLI(args)
	})
}

func runNodesPullErrorsCLI(args []string) error {
	cmdStr := "gitmap pull errors " + strings.Join(args, " ")
	cmdStr = strings.TrimSpace(cmdStr)

	return cmdssh.RunFleetPASCommand("nodes pull errors", cmdStr, func() error {
		return cmdpullerror.RunPullErrorCLI(args)
	})
}

func isNodesAgyRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	tok := strings.ToLower(args[0])
	return tok == "agy" || tok == "agy-ui" || tok == "nodes-agy-ui" || tok == "studio" || tok == "dashboard"
}

func runNodesAgyDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesAgyUI(nil)
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "query":
		return RunNodesAgyQuery(args[1:])
	case "prompt":
		return RunNodesAgyPrompt(args[1:])
	case "push-settings", "pushsettings":
		return RunNodesPushSettings(args[1:])
	case "sync-settings", "syncsettings":
		return RunNodesSyncSettings(args[1:])
	case "send-projects", "sendprojects", "sync-projects", "syncprojects":
		return RunNodesSendProjects(args[1:])
	case "ui", "dashboard", "studio":
		return RunNodesAgyUI(args[1:])
	default:
		return RunNodesAgyUI(args)
	}
}

func isNodesDeployAGMRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "sync-agm-accounts" || first == "sync-agm" || first == "deploy-agm-accounts" || first == "deploy-agm" {
		return true
	}
	if first != "deploy" || len(args) < 2 {
		return false
	}
	sub := strings.ToLower(args[1])
	return sub == "agm-accounts" || sub == "agm" || sub == "accounts" || sub == "agm-account"
}

func runNodesDeployAGMDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesDeployAGMAccounts(nil)
	}
	first := strings.ToLower(args[0])
	if first == "deploy" && len(args) > 1 {
		return RunNodesDeployAGMAccounts(args[2:])
	}
	return RunNodesDeployAGMAccounts(args[1:])
}

func isNodesDeployRepoRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "deploy-repo" || first == "deploy-repos" || first == "deployrepo" || first == "deployrepos" {
		return true
	}
	if first != "deploy" || len(args) < 2 {
		return false
	}
	sub := strings.ToLower(args[1])
	return sub == "repo" || sub == "repos" || sub == "repository" || sub == "repositories"
}

func runNodesDeployRepoDispatch(args []string) error {
	if len(args) == 0 {
		return RunNodesDeployRepo(nil)
	}
	first := strings.ToLower(args[0])
	if first == "deploy-repos" || first == "deployrepos" {
		return RunNodesDeployRepos(args[1:])
	}
	if first == "deploy-repo" || first == "deployrepo" {
		return RunNodesDeployRepo(args[1:])
	}
	if first == "deploy" && len(args) > 1 {
		return dispatchDeploySubcommand(args)
	}
	return RunNodesDeployRepo(args[1:])
}

func dispatchDeploySubcommand(args []string) error {
	sub := strings.ToLower(args[1])
	if sub == "repos" || sub == "repositories" {
		return RunNodesDeployRepos(args[2:])
	}
	return RunNodesDeployRepo(args[2:])
}

func isNodesScanRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	return first == "scan" || first == "rescan"
}

func runNodesScanDispatch(args []string) error {
	first := strings.ToLower(args[0])
	if first == "rescan" {
		return RunNodesRescan(args[1:])
	}
	return RunNodesScan(args[1:])
}

// runUnifiedNodesCLI handles `gitmap nodes [flags] [target]`.
func RunUnifiedNodesCLI(args []string) error {
	for len(args) > 0 && strings.EqualFold(args[0], "nodes") {
		args = args[1:]
	}
	if len(args) == 0 {
		// Proceed to default nodes listing
	} else if isNodesAgyRequest(args) {
		return runNodesAgyDispatch(args[1:])
	} else if strings.EqualFold(args[0], "push-settings") || strings.EqualFold(args[0], "pushsettings") {
		return RunNodesPushSettings(args[1:])
	} else if strings.EqualFold(args[0], "sync-settings") || strings.EqualFold(args[0], "syncsettings") {
		return RunNodesSyncSettings(args[1:])
	} else if isNodesDeployAGMRequest(args) {
		return runNodesDeployAGMDispatch(args)
	} else if isNodesDeployRepoRequest(args) {
		return runNodesDeployRepoDispatch(args)
	} else if isNodesScanRequest(args) {
		return runNodesScanDispatch(args)
	} else if strings.EqualFold(args[0], "send-projects") || strings.EqualFold(args[0], "sendprojects") ||
		strings.EqualFold(args[0], "sync-projects") || strings.EqualFold(args[0], "syncprojects") {
		return RunNodesSendProjects(args[1:])
	} else if strings.EqualFold(args[0], "history") || strings.EqualFold(args[0], "histories") {
		return runNodesHistoryCLI(args[1:])
	} else if strings.EqualFold(args[0], "errors") {
		return runNodesErrorsCLI(args[1:])
	} else if strings.EqualFold(args[0], "pull") && len(args) > 1 && strings.EqualFold(args[1], "errors") {
		return runNodesPullErrorsCLI(args[2:])
	} else if isNodesPingCommand(args[0]) {
		return RunUnifiedNodesPingCLI(args[1:])
	} else if strings.EqualFold(args[0], "pending-commits") || strings.EqualFold(args[0], "pc") || strings.EqualFold(args[0], "pendingcommits") {
		return RunNodesPendingCommits(args[1:])
	} else if strings.EqualFold(args[0], "commits") {
		return RunNodesCommits(args[1:])
	} else if strings.EqualFold(args[0], "cpf") {
		return RunNodesCommitPushFeature(args[1:])
	} else if strings.EqualFold(args[0], "cpb") {
		return RunNodesCommitPushBug(args[1:])
	} else if strings.EqualFold(args[0], "cpr") {
		return RunNodesCommitPushRelease(args[1:])
	} else if strings.EqualFold(args[0], "commit-fix") || strings.EqualFold(args[0], "commitfix") {
		return RunNodesCommitFix(args[1:])
	} else if isNodesFullSummary(args) {
		if strings.EqualFold(args[0], "full") && len(args) > 1 {
			return RunNodesFullSummary(args[2:])
		}
		return RunNodesFullSummary(args[1:])
	} else if strings.EqualFold(args[0], "summary") {
		return RunNodesSummary(args[1:])
	} else if isNodesPipelineAll(args) {
		if strings.EqualFold(args[0], "pe") && len(args) > 1 {
			return RunNodesPipelineErrorsAll(args[2:])
		}
		return RunNodesPipelineErrorsAll(args[1:])
	} else if isNodesCloneRequest(args) {
		return RunNodesClone(args)
	} else if isNodesHelpRequest(args) {
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
		return renderEmptyUnifiedNodesNotice(out)
	}
	if err := renderCommandsMatrixTable(out, nodes); err != nil {
		return err
	}
	if err := renderPrimaryNodesTable(out, nodes); err != nil {
		return err
	}
	renderUnifiedNodesFooter(out)
	return nil
}

func renderEmptyUnifiedNodesNotice(out io.Writer) error {
	fmt.Fprintln(out)
	fmt.Fprintf(out, "  %s● No nodes registered across SSH, Cluster, or Server-Client networks.%s\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Fprintln(out, "    Enroll nodes with: gitmap ssh join <user@ip|ip> [alias]")
	fmt.Fprintln(out)
	renderUnifiedNodesFooter(out)
	return nil
}

func renderPrimaryNodesTable(out io.Writer, nodes []UnifiedFleetNode) error {
	if err := renderHostsHeader(out); err != nil {
		return err
	}
	for _, n := range nodes {
		renderPrimaryNodeRow(out, n)
	}
	fmt.Fprintf(out, "\n  %sTotal: %d registered node(s)%s\n\n", constants.ColorDim, len(nodes), constants.ColorReset)
	return nil
}

func renderHostsHeader(out io.Writer) error {
	colAlias := padCell("ALIAS", 16)
	colRole := padCell("ROLE", 14)
	colHost := padCell("HOST (IP:PORT)", 22)
	colUser := padCell("USER", 14)
	colStatus := padCell("STATUS", 21)
	colEnrolled := padCell("ENROLLED", 19)

	headerLine := fmt.Sprintf("  %s %s %s %s %s %s\n",
		colAlias, colRole, colHost, colUser, colStatus, colEnrolled)
	fmt.Fprintf(out, "\n%s%s%s", constants.ColorCyan, headerLine, constants.ColorReset)
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
	return nil
}

func renderPrimaryNodeRow(out io.Writer, n UnifiedFleetNode) {
	alias := formatCellAlias(n.Alias, 16)
	role := formatCellRole(n.Role, 14)
	hostPort := formatCellHost(fmt.Sprintf("%s:%d", n.Host, n.Port), 22)
	user := formatCellDim(n.User, 14)
	status := formatCellStatus(n.Status, 21)
	enrolled := formatCellDim(n.EnrolledAt, 19)

	fmt.Fprintf(out, "  %s %s %s %s %s %s\n", alias, role, hostPort, user, status, enrolled)
}

func renderCommandsMatrixTable(out io.Writer, nodes []UnifiedFleetNode) error {
	renderCommandsMatrixHeader(out)
	for _, n := range nodes {
		renderCommandsMatrixRow(out, n)
	}
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n\n", constants.ColorDim, divider, constants.ColorReset)
	return nil
}

func renderCommandsMatrixHeader(out io.Writer) {
	colAlias := padCell("NODE (ALIAS)", 16)
	colSys := padCell("SUBSYSTEMS", 22)
	colCmds := padCell("SUPPORTED COMMANDS & CLUSTERS", 70)

	headerLine := fmt.Sprintf("  %s %s %s\n", colAlias, colSys, colCmds)
	fmt.Fprintf(out, "\n%s%s%s", constants.ColorCyan, headerLine, constants.ColorReset)
	divider := strings.Repeat("-", 110)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
}

func renderCommandsMatrixRow(out io.Writer, n UnifiedFleetNode) {
	alias := formatCellAlias(n.Alias, 16)
	subsys := formatCellSubsystems(strings.Join(n.Subsystems, ", "), 22)
	cmds := formatSupportedCommands(n)

	fmt.Fprintf(out, "  %s %s %s\n", alias, subsys, cmds)
}

func formatSupportedCommands(n UnifiedFleetNode) string {
	var parts []string
	if hasSubsystem(n, "SSH") {
		parts = append(parts, fmt.Sprintf("gitmap ssh %s", n.Alias))
		parts = append(parts, fmt.Sprintf("gitmap exec %s", n.Alias))
	}
	if hasSubsystem(n, "Cluster") {
		parts = append(parts, fmt.Sprintf("gitmap cluster exec %s", n.Alias))
	}
	if hasSubsystem(n, "SC") {
		parts = append(parts, "gitmap sc exec")
	}
	return constants.ColorDim + strings.Join(parts, ", ") + constants.ColorReset
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

func renderUnifiedNodesFooter(out io.Writer) {
	fmt.Fprintf(out, "  %s💡 Fleet Operations & Supported Command Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Fprintln(out, "    • Shell Access:            gitmap ssh <alias>                (or: gitmap ssh login <alias>)")
	fmt.Fprintln(out, "    • Command Execution:       gitmap exec <alias> \"<cmd>\"       (or: gitmap ssh exec <alias> \"<cmd>\")")
	fmt.Fprintln(out, "    • Cluster Remote Exec:     gitmap cluster exec <alias> \"<cmd>\"")
	fmt.Fprintln(out, "    • Broadcast to All Nodes:  gitmap sc exec \"<cmd>\"            (or: gitmap ssh exec all \"<cmd>\")")
	fmt.Fprintln(out, "    • Sync Keys & Ping Nodes:  gitmap ssh deploy-keys            | gitmap nodes ping")
	fmt.Fprintln(out, "    • Fleet Clones & CFR:      gitmap nodes clone <repo|file>    | gitmap nodes cfr <repo|file>")
	fmt.Fprintln(out)
}

func isNodesFullSummary(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "fs" || first == "fs+pe" || first == "fspe" ||
		first == "full-status" || first == "full-summary" ||
		first == "full-status+pe" || first == "full-summary+pe" {
		return true
	}
	if first == "full" && len(args) > 1 {
		second := strings.ToLower(args[1])
		return second == "status" || second == "summary" || second == "status+pe" || second == "summary+pe"
	}
	return false
}

func isNodesPipelineAll(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "pe" && len(args) > 1 && strings.EqualFold(args[1], "all") {
		return true
	}
	return first == "pipe-error-all" || first == "pipeline-error-all"
}
