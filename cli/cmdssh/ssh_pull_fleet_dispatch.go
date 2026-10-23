package cmdssh

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/crypto/ssh"
)

const (
	// MaxPASWorkers is the maximum concurrency for remote nodes under the PAS Formula.
	MaxPASWorkers = 2
	// MaxAsyncOpsPerNode is the maximum parallel operations per remote node.
	MaxAsyncOpsPerNode = 2
)

// RunLocalPullAllJSONFn executes pull-all in-process and returns structured JSON output.
var RunLocalPullAllJSONFn func(args []string) (string, error)

func executeFleetPullAll(probes []fleetNodeLiveness, cleanArgs []string) []FleetNodePullOutcome {
	var online []db.SSHConnection
	var skipped []FleetNodePullOutcome
	for _, p := range probes {
		if p.isLocalHost {
			skipped = append(skipped, buildLocalSkippedFleetOutcome(p.conn.Alias, p.conn.IPAddress))
			continue
		}
		if p.isOnline {
			online = append(online, p.conn)
			continue
		}
		skipped = append(skipped, buildSkippedFleetOutcome(p.conn.Alias, p.conn.IPAddress))
	}
	concurrencyFlags := ParsePASConcurrencyFlags(cleanArgs)
	dispatched := dispatchOnlineFleetPull(online, cleanArgs, concurrencyFlags)

	return append(dispatched, skipped...)
}

func dispatchOnlineFleetPull(conns []db.SSHConnection, cleanArgs, concurrencyFlags []string) []FleetNodePullOutcome {
	total := len(conns) + 1
	outcomes := make([]FleetNodePullOutcome, total)
	var wg sync.WaitGroup
	dispatchLocalVMPullAsync(&wg, outcomes, cleanArgs, concurrencyFlags)
	dispatchRemoteFleetNodes(&wg, outcomes, conns, concurrencyFlags)

	stopHeartbeat := startFleetPullHeartbeat(len(conns), hasFleetJSONFlag(cleanArgs))
	wg.Wait()
	stopHeartbeat()

	return outcomes
}

func dispatchLocalVMPullAsync(wg *sync.WaitGroup, outcomes []FleetNodePullOutcome, cleanArgs, concurrencyFlags []string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		outcomes[0] = executeLocalVMPull(cleanArgs, concurrencyFlags)
	}()
}

func dispatchRemoteFleetNodes(wg *sync.WaitGroup, outcomes []FleetNodePullOutcome, conns []db.SSHConnection, concurrencyFlags []string) {
	workerCount := resolvePASWorkerCount()
	sem := make(chan struct{}, workerCount)
	for idx, c := range conns {
		wg.Add(1)
		slot := idx + 1
		targetConn := c
		go executeBoundedRemotePull(wg, sem, outcomes, slot, targetConn, concurrencyFlags)
	}
}

func executeBoundedRemotePull(wg *sync.WaitGroup, sem chan struct{}, outcomes []FleetNodePullOutcome, slot int, c db.SSHConnection, concurrencyFlags []string) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()
	outcomes[slot] = executeRemoteNodePull(c, concurrencyFlags)
}

func executeLocalVMPull(cleanArgs, concurrencyFlags []string) FleetNodePullOutcome {
	if RunLocalPullAllJSONFn != nil {
		return executeLocalVMPullInProcess(cleanArgs)
	}
	return executeLocalVMPullSubprocess(cleanArgs, concurrencyFlags)
}

func executeLocalVMPullInProcess(cleanArgs []string) FleetNodePullOutcome {
	out, err := RunLocalPullAllJSONFn(cleanArgs)
	if err != nil && len(strings.TrimSpace(out)) == 0 {
		return buildFailedFleetOutcome("Local VM", "127.0.0.1", true, err.Error())
	}
	return parseFleetPullOutcome("Local VM", "127.0.0.1", true, out)
}

func executeLocalVMPullSubprocess(cleanArgs, concurrencyFlags []string) FleetNodePullOutcome {
	cmdName := resolveLocalGitmapExecutable()
	subArgs := []string{"pa", "--json"}
	if len(concurrencyFlags) > 0 {
		subArgs = append(subArgs, concurrencyFlags...)
	}
	cmd := exec.Command(cmdName, subArgs...)
	out, runErr := cmd.CombinedOutput()
	if runErr != nil && len(out) == 0 {
		return buildFailedFleetOutcome("Local VM", "127.0.0.1", true, runErr.Error())
	}

	return parseFleetPullOutcome("Local VM", "127.0.0.1", true, string(out))
}

func resolveLocalGitmapExecutable() string {
	exePath, err := os.Executable()
	if err == nil && exePath != "" {
		return exePath
	}
	return "gitmap"
}

func executeRemoteNodePull(c db.SSHConnection, concurrencyFlags []string) FleetNodePullOutcome {
	isOnline, _ := CheckConnLiveness(context.Background(), c.IPAddress, 22, 500*time.Millisecond)
	if !isOnline {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, "offline: node unreachable")
	}
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, "connection or auth failed")
	}
	defer client.Close()
	osType := resolveTargetNodeOSFast(client, c)
	ensureRemoteNodeGitmapFast(client, c, osType)
	remoteCmd := buildRemoteFleetPullCmd(osType, concurrencyFlags)
	shellType := resolveRemoteShell(osType)
	out, err := runRemotePullJSON(client, c, remoteCmd, shellType)
	if err != nil && len(out) == 0 {
		return buildFailedFleetOutcome(c.Alias, c.IPAddress, false, err.Error())
	}

	return parseFleetPullOutcome(c.Alias, c.IPAddress, false, out)
}

func resolveTargetNodeOSFast(client *ssh.Client, c db.SSHConnection) string {
	if c.OS != "" {
		return c.OS
	}

	return resolveTargetNodeOS(client, c)
}

func ensureRemoteNodeGitmapFast(client *ssh.Client, c db.SSHConnection, osType string) {
	if c.BuildVersion != "" && CompareSemverStrings(c.BuildVersion, "6.349.0") >= 0 {
		return
	}
	ensureRemoteNodeGitmap(client, c)
}

func buildRemoteFleetPullCmd(osType string, concurrencyFlags []string) string {
	cmd := "gitmap pa --json"
	if len(concurrencyFlags) > 0 {
		cmd += " " + strings.Join(concurrencyFlags, " ")
	} else {
		cmd += " --parallel 4"
	}
	if strings.EqualFold(osType, "windows") {
		prefix := "$env:GCM_CREDENTIAL_STORE='none'; $env:GCM_NO_PERSIST='1'; $env:GITMAP_HEADLESS_SSH='1'; "
		return prefix + cmd
	}
	prefix := "export GCM_CREDENTIAL_STORE=none GCM_NO_PERSIST=1 GITMAP_HEADLESS_SSH=1 && "

	return prefix + cmd
}

func ensureRemoteNodeGitmap(client *ssh.Client, c db.SSHConnection) {
	osType := resolveTargetNodeOS(client, c)
	ver, isInstalled := queryNodeVersionViaSSH(client, osType)
	if !isInstalled {
		cmd := BuildGitmapInstallOneLiner(osType, "latest")
		_, _ = secrets.RunCommand(client, cmd, resolveRemoteShell(osType))
		return
	}
	if CompareSemverStrings(ver, "6.349.0") < 0 {
		updateTargetNodeGitmap(client, osType)
	}
}

func updateTargetNodeGitmap(client *ssh.Client, osType string) {
	cmd := resolveRemoteUpdateCommand(osType, "gitmap")
	_, _ = secrets.RunCommand(client, cmd, resolveRemoteShell(osType))
}

func resolveTargetNodeOS(client *ssh.Client, c db.SSHConnection) string {
	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		return probed
	}

	return osType
}

func runRemotePullJSON(client *ssh.Client, c db.SSHConnection, remoteCmd, shellType string) (string, error) {
	stdout, stderr, err := secrets.RunCommandSeparated(client, remoteCmd, shellType)
	combined := stdout + "\n" + stderr
	if strings.Contains(combined, "flag provided but not defined: -json") {
		return handleLegacyRemotePull(client, c, remoteCmd, shellType)
	}

	if strings.Contains(combined, "pending task already exists for pa") {
		return recoverRemotePendingTask(client, combined, remoteCmd, shellType)
	}

	return stdout, err
}

func handleLegacyRemotePull(client *ssh.Client, c db.SSHConnection, remoteCmd, shellType string) (string, error) {
	osType := resolveTargetNodeOS(client, c)
	updateTargetNodeGitmap(client, osType)
	resolvedShell := resolveRemoteShell(osType)
	if shellType != "" {
		resolvedShell = shellType
	}
	stdout, _, err := secrets.RunCommandSeparated(client, remoteCmd, resolvedShell)

	return stdout, err
}

func recoverRemotePendingTask(client *ssh.Client, rawOutput, remoteCmd, shellType string) (string, error) {
	taskID := extractPendingTaskID(rawOutput)
	if taskID != "" {
		_, _ = secrets.RunCommand(client, "gitmap task cancel "+taskID, shellType)
	}
	stdout, _, err := secrets.RunCommandSeparated(client, remoteCmd, shellType)

	return stdout, err
}

func extractPendingTaskID(s string) string {
	lower := strings.ToLower(s)
	idx := strings.Index(lower, "(id ")
	if idx == -1 {
		return ""
	}

	sub := s[idx+4:]
	end := strings.Index(sub, ")")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(sub[:end])
}

func resolvePASWorkerCount() int {
	if isHighCPUPressure() {
		return 1
	}

	return MaxPASWorkers
}

func isHighCPUPressure() bool {
	if os.Getenv("GITMAP_HIGH_CPU") == "1" {
		return true
	}

	return runtime.NumCPU() <= 2
}
