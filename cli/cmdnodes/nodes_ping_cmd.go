package cmdnodes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var (
	reWinPackets   = lazyregex.New(`Sent\s*=\s*(\d+),\s*Received\s*=\s*(\d+),\s*Lost\s*=\s*(\d+)\s*\(([\d.]+)%\s*loss\)`)
	reWinRTT       = lazyregex.New(`Minimum\s*=\s*(\d+)ms,\s*Maximum\s*=\s*(\d+)ms,\s*Average\s*=\s*(\d+)ms`)
	reLinuxPackets = lazyregex.New(`(\d+)\s+(?:packets\s+)?transmitted,\s*(\d+)\s+(?:packets\s+)?received,\s*([\d.]+)%\s*packet\s*loss`)
	reLinuxRTT     = lazyregex.New(`rtt\s+min/avg/max/mdev\s*=\s*([\d.]+)/([\d.]+)/([\d.]+)/`)
)

// NodePingResult encapsulates both machine ICMP command and TCP reachability telemetry.
type NodePingResult struct {
	Alias             string        `json:"alias"`
	Role              string        `json:"role"`
	Host              string        `json:"host"`
	Port              int           `json:"port"`
	User              string        `json:"user"`
	Subsystems        []string      `json:"subsystems"`
	Command           string        `json:"command"`
	PacketsSent       int           `json:"packets_sent"`
	PacketsReceived   int           `json:"packets_received"`
	PacketsLost       int           `json:"packets_lost"`
	PacketLossPercent int           `json:"packet_loss_percent"`
	MinRTT            string        `json:"min_rtt"`
	AvgRTT            string        `json:"avg_rtt"`
	MaxRTT            string        `json:"max_rtt"`
	ICMPOk            bool          `json:"icmp_ok"`
	TCPOk             bool          `json:"tcp_ok"`
	TCPRTT            string        `json:"tcp_rtt"`
	Status            string        `json:"status"`
	RawOutput         string        `json:"raw_output,omitempty"`
	Duration          time.Duration `json:"duration_ms"`
}

type nodesPingOptions struct {
	count        int
	timeoutMs    int
	isJSON       bool
	isRaw        bool
	isTCPOnly    bool
	filterSSH    bool
	filterClust  bool
	filterSC     bool
	targetFilter string
}

func parseNodesPingOptions(args []string) nodesPingOptions {
	opts := nodesPingOptions{
		count:     2,
		timeoutMs: 1500,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if parsePingFlag(&opts, args, &i) {
			continue
		}
		if !strings.HasPrefix(arg, "-") && !isPingRoutingWord(arg) && opts.targetFilter == "" {
			opts.targetFilter = arg
		}
	}
	return opts
}

func isPingRoutingWord(arg string) bool {
	low := strings.ToLower(arg)
	return low == "ping" || low == "probe" || low == "nodes" || low == "node"
}

func parsePingFlag(opts *nodesPingOptions, args []string, index *int) bool {
	arg := args[*index]
	if arg == "--json" || arg == "-j" {
		opts.isJSON = true
		return true
	}
	if arg == "--raw" {
		opts.isRaw = true
		return true
	}
	if arg == "--tcp" {
		opts.isTCPOnly = true
		return true
	}
	if arg == "--ssh" {
		opts.filterSSH = true
		return true
	}
	if arg == "--cluster" || arg == "--clst" {
		opts.filterClust = true
		return true
	}
	if arg == "--sc" || arg == "--server-clients" {
		opts.filterSC = true
		return true
	}
	if parsePingCountFlag(opts, args, index) {
		return true
	}
	return parsePingTimeoutFlag(opts, args, index)
}

func parsePingCountFlag(opts *nodesPingOptions, args []string, index *int) bool {
	arg := args[*index]
	if (arg == "--count" || arg == "-c" || arg == "-n") && *index+1 < len(args) {
		*index++
		applyPositiveCount(opts, args[*index])
		return true
	}
	if strings.HasPrefix(arg, "--count=") {
		applyPositiveCount(opts, strings.TrimPrefix(arg, "--count="))
		return true
	}
	return false
}

func applyPositiveCount(opts *nodesPingOptions, raw string) {
	val, err := strconv.Atoi(raw)
	if err == nil && val > 0 {
		opts.count = val
	}
}

func parsePingTimeoutFlag(opts *nodesPingOptions, args []string, index *int) bool {
	arg := args[*index]
	if (arg == "--timeout" || arg == "-t" || arg == "-w") && *index+1 < len(args) {
		*index++
		opts.timeoutMs = parseTimeoutValue(args[*index])
		return true
	}
	if strings.HasPrefix(arg, "--timeout=") {
		opts.timeoutMs = parseTimeoutValue(strings.TrimPrefix(arg, "--timeout="))
		return true
	}
	return false
}

func parseTimeoutValue(val string) int {
	dur, err := time.ParseDuration(val)
	if err == nil && dur > 0 {
		return int(dur.Milliseconds())
	}
	msVal, err := strconv.Atoi(val)
	if err == nil && msVal > 0 {
		return msVal
	}
	return 1500
}

func isNodesPingHelpRequest(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func printUnifiedPingHelp() error {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║       gitmap nodes ping / gitmap ping - Fleet Machine Command Ping           ║")
	fmt.Println("  ╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap ping [target] [flags]")
	fmt.Println("    gitmap nodes ping [target] [flags]")
	fmt.Println("    gitmap ping help")
	fmt.Println()
	fmt.Println("  Description:")
	fmt.Println("    Runs the machine's native 'ping' command concurrently against registered")
	fmt.Println("    fleet nodes across SSH, Cluster DB, and Server-Clients networks, augmented")
	fmt.Println("    with port TCP reachability for nodes with ICMP firewall restrictions.")
	fmt.Println()
	fmt.Println("  Commands & Filtering:")
	fmt.Println("    gitmap ping                 Ping all registered nodes concurrently")
	fmt.Println("    gitmap ping <alias|ip>      Ping a specific fleet node or ad-hoc target")
	fmt.Println("    gitmap nodes ping           Ping all fleet nodes from the nodes subsystem")
	fmt.Println("    gitmap ping --ssh           Ping only SSH-enrolled nodes")
	fmt.Println("    gitmap ping --cluster       Ping only Cluster DB registered nodes")
	fmt.Println("    gitmap ping --sc            Ping only Server-Clients broadcast nodes")
	fmt.Println()
	fmt.Println("  Options & Flags:")
	fmt.Println("    -c, --count <n>             Number of ping packets to transmit (default: 2)")
	fmt.Println("    -t, --timeout <dur|ms>      Per-node probe timeout (default: 1500ms, e.g. 2s, 1000)")
	fmt.Println("    --raw                       Display exact machine command lines and stdout")
	fmt.Println("    --json, -j                  Output machine-readable JSON telemetry")
	fmt.Println("    -h, --help                  Show this formatted help banner")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ping")
	fmt.Println("    gitmap ping w1")
	fmt.Println("    gitmap ping 192.168.1.20 -c 3")
	fmt.Println("    gitmap ping --raw")
	fmt.Println("    gitmap nodes ping --json")
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Println("    gitmap ping, gitmap nodes ping, gitmap node ping, gitmap fleet-ping")
	fmt.Println()
	return nil
}

// runUnifiedNodesPingCLI executes the machine ping command across all or targeted fleet nodes.
func RunUnifiedNodesPingCLI(args []string) error {
	if isNodesPingHelpRequest(args) {
		return printUnifiedPingHelp()
	}

	opts := parseNodesPingOptions(args)
	ctx := context.Background()

	nodes, err := collectUnifiedNodes(ctx, true)
	if err != nil {
		return err
	}

	targetNodes := resolvePingTargetNodes(nodes, opts)
	if len(targetNodes) == 0 {
		return printNoPingTargetsNotice(os.Stdout, opts.targetFilter)
	}

	results := executeFleetPingConcurrent(ctx, targetNodes, opts)

	if opts.isJSON {
		return printPingJSON(results)
	}
	if opts.isRaw {
		printPingRaw(os.Stdout, results)
	}

	return renderPingTable(os.Stdout, results, opts)
}

func resolvePingTargetNodes(all []UnifiedFleetNode, opts nodesPingOptions) []UnifiedFleetNode {
	filterOpts := nodesFilterOptions{
		filterSSH:    opts.filterSSH,
		filterClust:  opts.filterClust,
		filterSC:     opts.filterSC,
		targetFilter: opts.targetFilter,
	}
	filtered := applyNodesFilter(all, filterOpts)
	if len(filtered) > 0 {
		return filtered
	}
	if opts.targetFilter != "" {
		return []UnifiedFleetNode{
			{
				Alias:      opts.targetFilter,
				Role:       "target",
				Host:       opts.targetFilter,
				Port:       22,
				User:       "operator",
				Subsystems: []string{"Ad-hoc"},
			},
		}
	}
	return nil
}

func printNoPingTargetsNotice(out io.Writer, target string) error {
	fmt.Fprintln(out)
	if target != "" {
		fmt.Fprintf(out, "  ● No registered fleet nodes match filter: %s\n", target)
		fmt.Fprintln(out, "    Verify available nodes with: gitmap nodes")
	} else {
		fmt.Fprintln(out, "  ● No nodes registered across SSH, Cluster, or Server-Client networks.")
		fmt.Fprintln(out, "    Enroll nodes with: gitmap ssh join <user@ip|ip> [alias]")
	}
	fmt.Fprintln(out)
	return nil
}

func executeFleetPingConcurrent(ctx context.Context, nodes []UnifiedFleetNode, opts nodesPingOptions) []NodePingResult {
	results := make([]NodePingResult, len(nodes))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for i, n := range nodes {
		wg.Add(1)
		go func(idx int, node UnifiedFleetNode) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[idx] = pingSingleNode(ctx, node, opts)
		}(i, n)
	}

	wg.Wait()
	return results
}

func pingSingleNode(ctx context.Context, node UnifiedFleetNode, opts nodesPingOptions) NodePingResult {
	start := time.Now()
	timeoutDur := time.Duration(opts.timeoutMs) * time.Millisecond
	nodeCtx, cancel := context.WithTimeout(ctx, timeoutDur+500*time.Millisecond)
	defer cancel()

	cmd := buildMachinePingCmd(nodeCtx, node.Host, opts.count, opts.timeoutMs)
	rawOut, _ := cmd.CombinedOutput()
	rawStr := string(rawOut)

	sent, recv, lost, lossPct, minR, avgR, maxR := parseMachinePingOutput(rawStr, opts.count)
	tcpOk, tcpRTT := probeNodeTCPPort(nodeCtx, node.Host, node.Port, timeoutDur)

	icmpOk := recv > 0
	status := resolvePingStatus(icmpOk, tcpOk, lossPct)

	return NodePingResult{
		Alias:             node.Alias,
		Role:              node.Role,
		Host:              node.Host,
		Port:              node.Port,
		User:              node.User,
		Subsystems:        node.Subsystems,
		Command:           cmd.String(),
		PacketsSent:       sent,
		PacketsReceived:   recv,
		PacketsLost:       lost,
		PacketLossPercent: lossPct,
		MinRTT:            minR,
		AvgRTT:            avgR,
		MaxRTT:            maxR,
		ICMPOk:            icmpOk,
		TCPOk:             tcpOk,
		TCPRTT:            tcpRTT,
		Status:            status,
		RawOutput:         strings.TrimSpace(rawStr),
		Duration:          time.Since(start),
	}
}

func buildMachinePingCmd(ctx context.Context, host string, count, timeoutMs int) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "ping", "-n", strconv.Itoa(count), "-w", strconv.Itoa(timeoutMs), host)
	}
	timeoutSec := timeoutMs / 1000
	if timeoutSec < 1 {
		timeoutSec = 1
	}
	return exec.CommandContext(ctx, "ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeoutSec), host)
}

func parseMachinePingOutput(raw string, count int) (sent, recv, lost, lossPct int, minR, avgR, maxR string) {
	sent = count
	lost = count
	lossPct = 100
	minR = "-"
	avgR = "-"
	maxR = "-"

	if parseWinPingOutput(raw, &sent, &recv, &lost, &lossPct, &minR, &avgR, &maxR) {
		return
	}
	parseLinuxPingOutput(raw, &sent, &recv, &lost, &lossPct, &minR, &avgR, &maxR)
	return
}

func parseWinPingOutput(raw string, sent, recv, lost, lossPct *int, minR, avgR, maxR *string) bool {
	matchPkt := reWinPackets.Compiled().FindStringSubmatch(raw)
	if len(matchPkt) < 5 {
		return false
	}

	*sent, _ = strconv.Atoi(matchPkt[1])
	*recv, _ = strconv.Atoi(matchPkt[2])
	*lost, _ = strconv.Atoi(matchPkt[3])
	*lossPct, _ = strconv.Atoi(matchPkt[4])

	matchRTT := reWinRTT.Compiled().FindStringSubmatch(raw)
	if len(matchRTT) >= 4 {
		*minR = matchRTT[1] + "ms"
		*maxR = matchRTT[2] + "ms"
		*avgR = resolveWinAvgRTT(matchRTT[3]+"ms", raw)
	}
	return true
}

func resolveWinAvgRTT(avg, raw string) string {
	if avg == "0ms" && strings.Contains(raw, "time<1ms") {
		return "<1ms"
	}
	return avg
}

func parseLinuxPingOutput(raw string, sent, recv, lost, lossPct *int, minR, avgR, maxR *string) bool {
	matchPkt := reLinuxPackets.Compiled().FindStringSubmatch(raw)
	if len(matchPkt) < 4 {
		return false
	}

	*sent, _ = strconv.Atoi(matchPkt[1])
	*recv, _ = strconv.Atoi(matchPkt[2])
	*lost = *sent - *recv
	val, _ := strconv.ParseFloat(matchPkt[3], 64)
	*lossPct = int(val)

	matchRTT := reLinuxRTT.Compiled().FindStringSubmatch(raw)
	if len(matchRTT) >= 4 {
		*minR = matchRTT[1] + "ms"
		*avgR = matchRTT[2] + "ms"
		*maxR = matchRTT[3] + "ms"
	}
	return true
}

func probeNodeTCPPort(ctx context.Context, host string, port int, timeout time.Duration) (bool, string) {
	if port <= 0 {
		port = 22
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := net.Dialer{Timeout: timeout}
	start := time.Now()

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, "timeout"
	}
	_ = conn.Close()
	elapsed := time.Since(start)
	return true, formatDurationShort(elapsed)
}

func formatDurationShort(d time.Duration) string {
	ms := d.Milliseconds()
	if ms < 1 {
		return "<1ms"
	}
	return fmt.Sprintf("%dms", ms)
}

func resolvePingStatus(icmpOk, tcpOk bool, lossPct int) string {
	if icmpOk && lossPct == 0 {
		return "● ONLINE"
	}
	if icmpOk && lossPct > 0 {
		return "▲ DEGRADED"
	}
	if tcpOk {
		return "● REACHABLE (TCP)"
	}
	return "○ OFFLINE"
}

func printPingJSON(results []NodePingResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "printPingJSON")
	}
	fmt.Println(string(data))
	return nil
}

func printPingRaw(out io.Writer, results []NodePingResult) {
	fmt.Fprintln(out)
	for _, r := range results {
		fmt.Fprintf(out, "--- [%s | %s] command: %s ---\n", r.Alias, r.Host, r.Command)
		if r.RawOutput != "" {
			fmt.Fprintln(out, r.RawOutput)
		} else {
			fmt.Fprintln(out, "(no output received)")
		}
		fmt.Fprintln(out)
	}
}

func renderPingTable(out io.Writer, results []NodePingResult, opts nodesPingOptions) error {
	renderPingHeaderBox(out, len(results), opts)
	renderPingColumnHeaders(out)

	var onlineCount int
	var tcpCount int
	var offlineCount int
	var maxDur time.Duration

	for _, r := range results {
		if r.Duration > maxDur {
			maxDur = r.Duration
		}
		switch {
		case r.ICMPOk:
			onlineCount++
		case r.TCPOk:
			tcpCount++
		default:
			offlineCount++
		}
		renderPingRow(out, r)
	}

	renderPingSummaryFooter(out, onlineCount, tcpCount, offlineCount, len(results), maxDur)
	return nil
}

func renderPingHeaderBox(out io.Writer, count int, opts nodesPingOptions) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(out, "║ GITMAP FLEET NODES PING (MACHINE ICMP & TCP PROBE)                                                               ║")
	fmt.Fprintln(out, "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintf(out, "  Probing %d node(s) across fleet with machine command 'ping' (packets: %d, timeout: %dms)...\n\n",
		count, opts.count, opts.timeoutMs)
}

func renderPingColumnHeaders(out io.Writer) {
	cAlias := padCell("ALIAS", 14)
	cRole := padCell("ROLE", 12)
	cHost := padCell("HOST (IP:PORT)", 22)
	cICMP := padCell("ICMP STATS", 15)
	cLoss := padCell("LOSS", 8)
	cAvg := padCell("AVG RTT", 10)
	cTCP := padCell("TCP (PORT)", 16)
	cStatus := padCell("STATUS", 20)

	fmt.Fprintf(out, "  %s%s %s %s %s %s %s %s %s%s\n",
		constants.ColorCyan, cAlias, cRole, cHost, cICMP, cLoss, cAvg, cTCP, cStatus, constants.ColorReset)
	divider := strings.Repeat("─", 124)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)
}

func renderPingRow(out io.Writer, r NodePingResult) {
	alias := formatCellAlias(r.Alias, 14)
	role := formatCellRole(r.Role, 12)
	hostPort := formatCellHost(fmt.Sprintf("%s:%d", r.Host, r.Port), 22)
	icmpStats := formatCellICMP(r.PacketsReceived, r.PacketsSent, 15)
	loss := formatCellLoss(r.PacketLossPercent, 8)
	avgRTT := formatCellDim(r.AvgRTT, 10)
	tcpPort := formatCellTCP(r.TCPOk, r.Port, r.TCPRTT, 16)
	status := formatCellPingStatus(r.Status, 20)

	fmt.Fprintf(out, "  %s %s %s %s %s %s %s %s\n",
		alias, role, hostPort, icmpStats, loss, avgRTT, tcpPort, status)
}

func formatCellICMP(recv, sent, width int) string {
	val := fmt.Sprintf("%d/%d rcvd", recv, sent)
	if recv > 0 {
		return constants.ColorGreen + padCell(val, width) + constants.ColorReset
	}
	return constants.ColorDim + padCell(val, width) + constants.ColorReset
}

func formatCellLoss(lossPct, width int) string {
	val := fmt.Sprintf("%d%%", lossPct)
	if lossPct == 0 {
		return constants.ColorGreen + padCell(val, width) + constants.ColorReset
	}
	if lossPct < 100 {
		return constants.ColorYellow + padCell(val, width) + constants.ColorReset
	}
	return constants.ColorDim + padCell(val, width) + constants.ColorReset
}

func formatCellTCP(tcpOk bool, port int, rtt string, width int) string {
	if tcpOk {
		val := fmt.Sprintf("● %d (%s)", port, rtt)
		return constants.ColorGreen + padCell(val, width) + constants.ColorReset
	}
	val := "○ timeout"
	return constants.ColorDim + padCell(val, width) + constants.ColorReset
}

func formatCellPingStatus(status string, width int) string {
	if strings.Contains(status, "ONLINE") {
		return constants.ColorGreen + padCell(status, width) + constants.ColorReset
	}
	if strings.Contains(status, "REACHABLE") {
		return constants.ColorCyan + padCell(status, width) + constants.ColorReset
	}
	if strings.Contains(status, "DEGRADED") {
		return constants.ColorYellow + padCell(status, width) + constants.ColorReset
	}
	return constants.ColorDim + padCell(status, width) + constants.ColorReset
}

func renderPingSummaryFooter(out io.Writer, online, tcp, offline, total int, dur time.Duration) {
	fmt.Fprintln(out)
	divider := strings.Repeat("─", 124)
	fmt.Fprintf(out, "  %s%s%s\n", constants.ColorDim, divider, constants.ColorReset)

	reachableTotal := online + tcp
	durSec := float64(dur.Milliseconds()) / 1000.0

	fmt.Fprintf(out, "  Fleet Ping: %d/%d reachable (%d online ICMP, %d reachable via TCP) | %d offline | Max Elapsed: %.2fs\n",
		reachableTotal, total, online, tcp, offline, durSec)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "==============================================================================================================================")
	fmt.Fprintln(out, " TIP: Run 'gitmap ssh <alias>' for shell access, or 'gitmap ping <alias>' to ping an individual node.")
	fmt.Fprintln(out, "==============================================================================================================================")
	fmt.Fprintln(out)
}
