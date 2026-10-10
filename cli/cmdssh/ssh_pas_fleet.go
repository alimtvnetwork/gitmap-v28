package cmdssh

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// FleetPASOutcome captures command execution on a single local or remote node.
type FleetPASOutcome struct {
	NodeName  string
	IP        string
	IsLocal   bool
	Success   bool
	IsSkipped bool
	Output    string
	ErrorMsg  string
}

// RunFleetPASCommand executes a task following the GitMap PAS Formula.
func RunFleetPASCommand(label, remoteCmd string, localFn func() error) error {
	conns, _ := fetchAllSSHConnections()
	probes := probeFleetLiveness(conns)
	printFleetPASBanner(label, probes)

	var online []db.SSHConnection
	var skipped []FleetPASOutcome
	for _, p := range probes {
		if p.isOnline {
			online = append(online, p.conn)
			continue
		}
		skipped = append(skipped, FleetPASOutcome{
			NodeName:  p.conn.Alias,
			IP:        p.conn.IPAddress,
			IsSkipped: true,
			ErrorMsg:  "offline (skipped, no task enqueued)",
		})
	}

	outcomes := dispatchPASFleetWork(online, remoteCmd, localFn)
	outcomes = append(outcomes, skipped...)
	renderPASFleetSummary(label, outcomes)
	return nil
}

func printFleetPASBanner(label string, probes []fleetNodeLiveness) {
	fmt.Printf("\n%s  Enqueuing '%s' across SSH fleet:%s\n", constants.ColorCyan, label, constants.ColorReset)
	for _, p := range probes {
		c := p.conn
		if !p.isOnline {
			fmt.Printf("    • Remote Node [%s] (%s): %sOffline (skipped, no task enqueued)%s\n",
				c.Alias, c.IPAddress, constants.ColorYellow, constants.ColorReset)
			continue
		}
		verStr := p.version
		if verStr != "" && verStr != "unknown" && !strings.HasPrefix(verStr, "v") {
			verStr = "v" + verStr
		}
		fmt.Printf("    • Remote Node [%s] (%s) [GitMap %s]: %sOnline → Enqueued (async)%s\n",
			c.Alias, c.IPAddress, verStr, constants.ColorGreen, constants.ColorReset)
	}
	host := resolveLocalHostname()
	localVer := constants.Version
	if !strings.HasPrefix(localVer, "v") {
		localVer = "v" + localVer
	}
	fmt.Printf("    • Current Machine [%s (127.0.0.1)] [GitMap %s]: %sRunning locally (direct execution, not enqueued)%s\n\n",
		host, localVer, constants.ColorGreen, constants.ColorReset)
}

func dispatchPASFleetWork(conns []db.SSHConnection, remoteCmd string, localFn func() error) []FleetPASOutcome {
	total := len(conns) + 1
	outcomes := make([]FleetPASOutcome, total)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		localErr := localFn()
		errMsg := ""
		if localErr != nil {
			errMsg = localErr.Error()
		}
		outcomes[0] = FleetPASOutcome{
			NodeName: "Local VM",
			IP:       "127.0.0.1",
			IsLocal:  true,
			Success:  localErr == nil,
			ErrorMsg: errMsg,
		}
	}()

	for idx, c := range conns {
		wg.Add(1)
		slot := idx + 1
		targetConn := c
		go func() {
			defer wg.Done()
			outcomes[slot] = executeRemotePASCommand(targetConn, remoteCmd)
		}()
	}
	wg.Wait()
	return outcomes
}

func executeRemotePASCommand(c db.SSHConnection, remoteCmd string) FleetPASOutcome {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return FleetPASOutcome{
			NodeName: c.Alias,
			IP:       c.IPAddress,
			Success:  false,
			ErrorMsg: "connection or auth failed",
		}
	}
	defer client.Close()

	osType := resolveTargetNodeOS(client, c)
	out, err := secrets.RunCommand(client, remoteCmd, resolveRemoteShell(osType))
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return FleetPASOutcome{
		NodeName: c.Alias,
		IP:       c.IPAddress,
		Success:  err == nil,
		Output:   strings.TrimSpace(out),
		ErrorMsg: errMsg,
	}
}

func renderPASFleetSummary(label string, outcomes []FleetPASOutcome) {
	fmt.Printf("\n%s  ▶ Fleet '%s' Summary:%s\n", constants.ColorCyan, label, constants.ColorReset)
	for _, o := range outcomes {
		if o.IsSkipped {
			fmt.Printf("    • [%s] (%s): %sSkipped%s (%s)\n",
				o.NodeName, o.IP, constants.ColorYellow, constants.ColorReset, o.ErrorMsg)
			continue
		}
		if o.Success {
			fmt.Printf("    • [%s] (%s): %sSuccess%s\n",
				o.NodeName, o.IP, constants.ColorGreen, constants.ColorReset)
			printIndentedOutputIfPresent(o.Output)
			continue
		}
		fmt.Printf("    • [%s] (%s): %sFailed%s (%s)\n",
			o.NodeName, o.IP, constants.ColorRed, constants.ColorReset, o.ErrorMsg)
		printIndentedOutputIfPresent(o.Output)
	}
	fmt.Println()
}

func printIndentedOutputIfPresent(output string) {
	if output == "" {
		return
	}
	printIndentedOutput(output)
}

func printIndentedOutput(output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			fmt.Printf("        %s\n", trimmed)
		}
	}
}

// ParsePASConcurrencyFlags extracts worker and hand concurrency flags from command args.
func ParsePASConcurrencyFlags(args []string) []string {
	var forwarded []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--wwoh" || a == "--worker-with-one-hand":
			forwarded = append(forwarded, a)
		case a == "--w" || a == "--worker" || a == "--workers" || a == "-w":
			forwarded = append(forwarded, a)
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				forwarded = append(forwarded, args[i])
			}
		case a == "--hand" || a == "--hands" || a == "--h":
			forwarded = append(forwarded, a)
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				forwarded = append(forwarded, args[i])
			}
		case a == "-h":
			if i+1 < len(args) && isIntegerArg(args[i+1]) {
				forwarded = append(forwarded, "--hand", args[i+1])
				i++
			}
		case strings.HasPrefix(a, "--w=") || strings.HasPrefix(a, "--worker=") || strings.HasPrefix(a, "--workers="):
			forwarded = append(forwarded, a)
		case strings.HasPrefix(a, "--hand=") || strings.HasPrefix(a, "--hands=") || strings.HasPrefix(a, "--h="):
			forwarded = append(forwarded, a)
		}
	}
	return forwarded
}

// DispatchOnlineFleetPAS constructs remote commands forwarding concurrency flags and dispatches fleet work.
func DispatchOnlineFleetPAS(conns []db.SSHConnection, baseRemoteCmd string, args []string, localFn func() error) []FleetPASOutcome {
	forwardedFlags := ParsePASConcurrencyFlags(args)
	remoteCmd := baseRemoteCmd
	if len(forwardedFlags) > 0 {
		remoteCmd = remoteCmd + " " + strings.Join(forwardedFlags, " ")
	}
	return dispatchPASFleetWork(conns, remoteCmd, localFn)
}

// RunSSHPASFleet executes pull-all across SSH fleet nodes, forwarding concurrency flags to remote nodes.
func RunSSHPASFleet(cleanArgs []string) error {
	return RunSSHPullAllFleet(cleanArgs)
}

func isIntegerArg(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
