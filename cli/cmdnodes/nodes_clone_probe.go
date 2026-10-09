// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type gitmapVersionPayload struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	OS        string `json:"os,omitempty"`
	Arch      string `json:"arch,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
}

var semverRegex = regexp.MustCompile(`v?[0-9]+\.[0-9]+(?:\.[0-9]+)?(?:-[a-zA-Z0-9.]+)?`)

// probeFleetNodesPreFlight probes all remote candidate nodes concurrently.
func probeFleetNodesPreFlight(ctx context.Context, conns []db.SSHConnection, opts NodesCloneOptions) FleetPreFlightReport {
	candidates := getProbeCandidates(conns, opts)
	if len(candidates) == 0 {
		return FleetPreFlightReport{Nodes: []PreFlightNodeInfo{}}
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	nodes := make([]PreFlightNodeInfo, 0, len(candidates))
	for _, c := range candidates {
		wg.Add(1)
		go spawnNodeProbeWorker(ctx, &wg, &mu, c, opts, &nodes)
	}
	wg.Wait()
	return buildFleetReport(candidates, nodes)
}

// ProbeFleetReadiness is an exported wrapper conforming to the preflight probe specification.
func ProbeFleetReadiness(conns []db.SSHConnection, opts NodesCloneOptions, timeout time.Duration) FleetPreFlightReport {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return probeFleetNodesPreFlight(ctx, conns, opts)
}

func getProbeCandidates(conns []db.SSHConnection, opts NodesCloneOptions) []db.SSHConnection {
	if len(conns) == 0 {
		return nil
	}
	var res []db.SSHConnection
	for _, c := range conns {
		if isConnectionExcluded(c, opts) {
			continue
		}
		res = append(res, c)
	}
	return res
}

func buildFleetReport(candidates []db.SSHConnection, nodes []PreFlightNodeInfo) FleetPreFlightReport {
	onlineCount := countOnlineNodes(nodes)
	total := len(candidates)
	offlineCount := total - onlineCount
	if offlineCount < 0 {
		offlineCount = 0
	}
	return FleetPreFlightReport{
		TotalCount:   total,
		OnlineCount:  onlineCount,
		OfflineCount: offlineCount,
		Nodes:        nodes,
	}
}

func countOnlineNodes(nodes []PreFlightNodeInfo) int {
	count := 0
	for _, n := range nodes {
		if n.IsOnline {
			count++
		}
	}
	return count
}

func spawnNodeProbeWorker(ctx context.Context, wg *sync.WaitGroup, mu *sync.Mutex, conn db.SSHConnection, opts NodesCloneOptions, nodes *[]PreFlightNodeInfo) {
	defer wg.Done()
	info := probeSingleNodeWithTimeout(ctx, conn, opts)
	mu.Lock()
	*nodes = append(*nodes, info)
	mu.Unlock()
}

func probeSingleNodeWithTimeout(ctx context.Context, conn db.SSHConnection, opts NodesCloneOptions) PreFlightNodeInfo {
	nodeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	start := time.Now()
	destDir := resolveRemoteDisplayTargetDir(conn, opts)
	return awaitNodeProbe(nodeCtx, conn, destDir, start)
}

func awaitNodeProbe(ctx context.Context, conn db.SSHConnection, destDir string, start time.Time) PreFlightNodeInfo {
	infoCh := make(chan PreFlightNodeInfo, 1)
	go func() {
		infoCh <- probeSingleNodeConnection(conn, destDir, start)
	}()
	select {
	case <-ctx.Done():
		return buildTimeoutNodeInfo(conn, destDir, time.Since(start))
	case info := <-infoCh:
		return info
	}
}

func buildTimeoutNodeInfo(conn db.SSHConnection, destDir string, dur time.Duration) PreFlightNodeInfo {
	errStr := "connection timed out (3s limit exceeded)"
	return PreFlightNodeInfo{
		Alias:           conn.Alias,
		Host:            conn.IPAddress,
		Role:            "worker",
		OS:              conn.OS,
		IsOnline:        false,
		StatusBadge:     resolvePreFlightBadge(false, errStr),
		RemoteTargetDir: destDir,
		ProbeDuration:   dur,
		Error:           errStr,
	}
}

func probeSingleNodeConnection(conn db.SSHConnection, destDir string, start time.Time) PreFlightNodeInfo {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return buildOfflineNodeInfo(conn, destDir, time.Since(start), errConnect.Error())
	}
	defer client.Close()
	return inspectRemoteNodeDetails(client, conn, destDir, start)
}

func buildOfflineNodeInfo(conn db.SSHConnection, destDir string, dur time.Duration, errStr string) PreFlightNodeInfo {
	return PreFlightNodeInfo{
		Alias:           conn.Alias,
		Host:            conn.IPAddress,
		Role:            "worker",
		OS:              conn.OS,
		IsOnline:        false,
		StatusBadge:     resolvePreFlightBadge(false, errStr),
		RemoteTargetDir: destDir,
		ProbeDuration:   dur,
		Error:           errStr,
	}
}

func resolveVersionCommand(isWin bool) (string, string) {
	if isWin {
		return `powershell -NoProfile -Command "gitmap version --json"`, "ps"
	}
	return "gitmap version --json 2>/dev/null || gitmap --version 2>/dev/null", "bash"
}

func queryRemoteVersionRaw(client *ssh.Client, isWin bool) (string, error) {
	cmd, shell := resolveVersionCommand(isWin)
	out, err := secrets.RunCommand(client, cmd, shell)
	if err != nil && !isWin && isBashMissingError(out, err) {
		out, err = secrets.RunCommand(client, cmd, "sh")
	}
	return out, err
}

func inspectRemoteNodeDetails(client *ssh.Client, conn db.SSHConnection, destDir string, start time.Time) PreFlightNodeInfo {
	isWin := isWindowsNode(conn)
	out, errCmd := queryRemoteVersionRaw(client, isWin)
	dur := time.Since(start)
	if errCmd != nil && strings.TrimSpace(out) == "" {
		return buildOfflineNodeInfo(conn, destDir, dur, errCmd.Error())
	}
	return buildOnlineNodeFromOutput(conn, destDir, out, dur)
}

func parseVersionJSON(raw string) (gitmapVersionPayload, bool) {
	var payload gitmapVersionPayload
	startIdx := strings.Index(raw, "{")
	endIdx := strings.LastIndex(raw, "}")
	if startIdx == -1 || endIdx <= startIdx {
		return payload, false
	}
	jsonStr := raw[startIdx : endIdx+1]
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil || payload.Version == "" {
		return payload, false
	}
	return payload, true
}

func extractFallbackVersion(raw string) string {
	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		clean := strings.TrimSpace(l)
		if strings.Contains(strings.ToLower(clean), "not found") || strings.Contains(strings.ToLower(clean), "not recognized") {
			return "-"
		}
		if match := semverRegex.FindString(clean); match != "" {
			return match
		}
	}
	return "-"
}

func recordDiscoveredVersion(alias, ver string) {
	if ver != "" && ver != "-" {
		recordNodeVersion(alias, ver)
	}
}

func buildOnlineNodeFromOutput(conn db.SSHConnection, destDir, out string, dur time.Duration) PreFlightNodeInfo {
	payload, hasJSON := parseVersionJSON(out)
	ver, commit, arch, osType := resolveVersionFields(conn, payload, out, hasJSON)
	recordDiscoveredVersion(conn.Alias, ver)
	return createOnlineNodeInfo(conn, destDir, ver, commit, arch, osType, dur)
}

func resolveVersionFields(conn db.SSHConnection, p gitmapVersionPayload, out string, hasJSON bool) (string, string, string, string) {
	if !hasJSON {
		return extractFallbackVersion(out), "", "", conn.OS
	}

	osType := conn.OS
	if p.OS != "" {
		osType = p.OS
	}

	return p.Version, p.Commit, p.Arch, osType
}

func createOnlineNodeInfo(conn db.SSHConnection, destDir, ver, commit, arch, osType string, dur time.Duration) PreFlightNodeInfo {
	return PreFlightNodeInfo{
		Alias:           conn.Alias,
		Host:            conn.IPAddress,
		Role:            "worker",
		OS:              osType,
		Arch:            arch,
		Version:         ver,
		Commit:          commit,
		IsOnline:        true,
		StatusBadge:     resolvePreFlightBadge(true, ""),
		RemoteTargetDir: destDir,
		ProbeDuration:   dur,
	}
}
