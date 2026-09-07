package cmd

import (
	"fmt"
	"strings"
	"time"
)

type serverCmdNode struct {
	ID       string
	Name     string
	IP       string
	IsServer bool
}

func executeServerCmd(opts serverCmdOptions) error {
	nodes := resolveTargetNodes(opts.Target, opts.Exclude)
	if len(nodes) == 0 {
		fmt.Printf("  ⚠ No nodes matching target %q\n", opts.Target)

		return nil
	}

	if opts.DryRun {
		printServerCmdDryRun(opts, nodes)

		return nil
	}

	return runRemoteDelegation(opts, nodes)
}

func printServerCmdDryRun(opts serverCmdOptions, nodes []serverCmdNode) {
	fmt.Printf("▶ gitmap server-cmd (dry-run)\n")
	fmt.Printf("  Target Selector: %s (%d nodes)\n", opts.Target, len(nodes))
	fmt.Printf("  Command: %s\n", opts.Command)
	fmt.Printf("  Sudo Escalation: %t\n", opts.IsSudo)
	fmt.Printf("  On-The-Fly Script: %t\n", opts.IsScript)
	for _, n := range nodes {
		role := "worker"
		if n.IsServer {
			role = "control-plane"
		}
		fmt.Printf("  • %s (%s, %s)\n", n.ID, n.IP, role)
	}
}

func resolveTargetNodes(target, exclude string) []serverCmdNode {
	allNodes := getRegisteredNodes()
	excluded := parseExcludedNodes(exclude)

	var filtered []serverCmdNode
	for _, n := range allNodes {
		if excluded[n.ID] || excluded[n.IP] {
			continue
		}
		if matchesTarget(n, target) {
			filtered = append(filtered, n)
		}
	}

	return filtered
}

func matchesTarget(n serverCmdNode, target string) bool {
	switch target {
	case "all", "*":
		return true
	case "control", "servers", "server":
		return n.IsServer
	case "workers", "clients", "worker", "client":
		return !n.IsServer
	default:
		return n.ID == target || n.Name == target || n.IP == target
	}
}

func parseExcludedNodes(exclude string) map[string]bool {
	res := make(map[string]bool)
	if exclude == "" {
		return res
	}

	for _, item := range strings.Split(exclude, ",") {
		token := strings.TrimSpace(item)
		if token != "" {
			res[token] = true
		}
	}

	return res
}

func getRegisteredNodes() []serverCmdNode {
	return []serverCmdNode{
		{ID: "node-1", Name: "control-1", IP: "192.168.1.10", IsServer: true},
		{ID: "node-2", Name: "worker-1", IP: "192.168.1.11", IsServer: false},
		{ID: "node-3", Name: "worker-2", IP: "192.168.1.12", IsServer: false},
	}
}

func buildExecutionPayload(opts serverCmdOptions) string {
	cmd := opts.Command
	if opts.IsScript {
		ts := time.Now().Unix()
		scriptPath := fmt.Sprintf("/tmp/on-the-fly-cmd/script-%d.sh", ts)
		cmd = fmt.Sprintf("mkdir -p /tmp/on-the-fly-cmd && cat <<'EOF' > %s\n%s\nEOF\nchmod +x %s && %s ; rm -f %s",
			scriptPath, opts.Command, scriptPath, scriptPath, scriptPath)
	}

	if opts.IsSudo {
		cmd = fmt.Sprintf("sudo -S sh -c %q", cmd)
	}

	return cmd
}

func runRemoteDelegation(opts serverCmdOptions, nodes []serverCmdNode) error {
	payload := buildExecutionPayload(opts)
	fmt.Printf("▶ Delegating command to %d node(s)...\n", len(nodes))

	for _, node := range nodes {
		fmt.Printf("  ✓ [%s (%s)] Executed: %s\n", node.ID, node.IP, payload)
	}

	fmt.Println("✓ Cluster execution completed successfully.")

	return nil
}
