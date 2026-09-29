// Package cmdssh — ssh_deploy_config_ssh.go deploys SSH fleet configuration across machines.
package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// DeployConfigSSHOptions holds arguments for deploying SSH fleet configuration.
type DeployConfigSSHOptions struct {
	Target   string
	FilePath string
	Except   string
	IsDryRun bool
	IsJSON   bool
	IsForce  bool
}

// RunSSHDeployConfigSSHCLI handles `gitmap deploy config ssh [target] [flags]`.
func RunSSHDeployConfigSSHCLI(args []string) error {
	args = stripLeadingSSHArg(args)
	if isHelpDeployConfigSSHRequest(args) {
		RenderDeployConfigSSHHelp()
		return nil
	}
	opts := parseDeployConfigSSHFlags(args)
	envelope, err := loadDeployConfigEnvelope(opts.FilePath)
	if err != nil {
		return err
	}
	envelope.Connections = portableEncryptConnections(envelope.Connections)

	allTargets, err := resolveDeployNodes(opts.Target)
	if err != nil {
		return err
	}
	filtered := filterConfigDeployTargets(allTargets, opts)

	if opts.IsDryRun || len(filtered) == 0 {
		return renderNodeConfigDeploySummary(filtered, opts.Except, len(envelope.Connections), opts.IsDryRun, opts.IsJSON)
	}
	return executeDeployConfigToFleet(filtered, envelope, opts)
}

func isHelpDeployConfigSSHRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "help" || low == "--help" || low == "-h" {
			return true
		}
	}
	return false
}

func stripLeadingSSHArg(args []string) []string {
	if len(args) > 0 && strings.EqualFold(args[0], "ssh") {
		return args[1:]
	}
	return args
}

func parseDeployConfigSSHFlags(args []string) DeployConfigSSHOptions {
	opts := DeployConfigSSHOptions{Target: "all"}
	for i := 0; i < len(args); i++ {
		i = consumeDeployConfigFlag(args, i, &opts)
	}
	return opts
}

func consumeDeployConfigFlag(args []string, i int, opts *DeployConfigSSHOptions) int {
	a := args[i]
	low := strings.ToLower(a)
	if isConfigSkipToken(low) {
		return i
	}
	if low == "--dry-run" || low == "-n" {
		opts.IsDryRun = true
		return i
	}
	if low == "--json" || low == "-j" {
		opts.IsJSON = true
		return i
	}
	if low == "--force" {
		opts.IsForce = true
		return i
	}
	return consumeValueOrTargetDeployConfig(args, i, opts)
}

func consumeValueOrTargetDeployConfig(args []string, i int, opts *DeployConfigSSHOptions) int {
	a := args[i]
	low := strings.ToLower(a)
	if isFileFlag(low) && i+1 < len(args) {
		opts.FilePath = args[i+1]
		return i + 1
	}
	if strings.HasPrefix(low, "--file=") || strings.HasPrefix(low, "-f=") {
		opts.FilePath = a[strings.IndexByte(a, '=')+1:]
		return i
	}
	if isExceptOrExcepFlag(low) && i+1 < len(args) {
		opts.Except = appendToken(opts.Except, args[i+1])
		return i + 1
	}
	if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--exclude=") {
		opts.Except = appendToken(opts.Except, a[strings.IndexByte(a, '=')+1:])
		return i
	}
	if !strings.HasPrefix(a, "-") && (opts.Target == "all" || opts.Target == "") {
		opts.Target = a
	}
	return i
}

func isConfigSkipToken(s string) bool {
	return s == "config" || s == "ssh" || s == "deploy" || s == "node-config" || s == "nc" || s == "all"
}

func isFileFlag(s string) bool {
	return s == "--file" || s == "-f" || s == "--from-file"
}

func appendToken(existing, addition string) string {
	if existing == "" {
		return addition
	}
	return existing + "," + addition
}

func loadDeployConfigEnvelope(filePath string) (*SSHNodesExportEnvelope, error) {
	if filePath != "" {
		raw, _, err := readNodesImportFileWithFallback(filePath)
		if err != nil {
			return nil, err
		}
		conns, err := decodeConnectionsFromJSON(raw)
		if err != nil {
			return nil, err
		}
		for i := range conns {
			conns[i] = enrichConnectionPassword(conns[i])
		}
		return buildEnvelopeFromConnections(conns), nil
	}
	return BuildSSHNodesExportEnvelope()
}

func buildEnvelopeFromConnections(conns []db.SSHConnection) *SSHNodesExportEnvelope {
	items := make([]SSHNodeExportItem, 0, len(conns))
	for idx, c := range conns {
		items = append(items, SSHNodeExportItem{
			WorkerId:   fmt.Sprintf("worker-%d", idx+1),
			NumericId:  idx + 1,
			Alias:      c.Alias,
			IPAddress:  c.IPAddress,
			Username:   c.Username,
			Port:       22,
			OS:         c.OS,
			AuthMethod: "key",
			KeyPath:    c.KeyPath,
		})
	}
	return &SSHNodesExportEnvelope{
		SchemaVersion: "1.0",
		TotalNodes:    len(conns),
		Nodes:         items,
		Connections:   conns,
	}
}

func filterConfigDeployTargets(targets []db.SSHConnection, opts DeployConfigSSHOptions) []db.SSHConnection {
	filtered := FilterSSHConnectionsByExcept(targets, opts.Except)
	if !isDeployAllTarget(opts.Target) {
		return filtered
	}
	var remoteOnly []db.SSHConnection
	for _, c := range filtered {
		if isLocalMachineIP(c.IPAddress) {
			continue
		}
		remoteOnly = append(remoteOnly, c)
	}
	return remoteOnly
}

func isLocalMachineIP(ip string) bool {
	if ip == "" || ip == "127.0.0.1" || ip == "localhost" {
		return true
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.String() == ip {
				return true
			}
		}
	}
	return false
}

func executeDeployConfigToFleet(targets []db.SSHConnection, envelope *SSHNodesExportEnvelope, opts DeployConfigSSHOptions) error {
	compactBytes, _ := json.Marshal(envelope)
	b64 := base64.StdEncoding.EncodeToString(compactBytes)
	remoteCmd := fmt.Sprintf("gitmap ssh nodes import-json --base64 \"%s\"", b64)

	fmt.Printf("\n%s🚀 GitMap Deploy Config SSH (Fleet Node Synchronization)%s\n", constants.ColorCyan, constants.ColorReset)
	sourceDesc := "CL Active Config"
	if opts.FilePath != "" {
		sourceDesc = opts.FilePath
	}
	fmt.Printf("  • Source Topology:  %s (%d node(s))\n", sourceDesc, len(envelope.Connections))
	fmt.Printf("  • Target Fleet:     %d remote node(s)\n", len(targets))
	fmt.Printf("  • Encryption:       AES-256 Symmetric Fleet Key (aes:...)\n\n")

	fleetOpts := FleetParallelOptions{
		Target:   "all",
		Except:   opts.Except,
		TaskName: "Deploy SSH Fleet Config (deploy config ssh)",
	}
	RunParallelFleetExecution(targets, fleetOpts, func(c db.SSHConnection) (string, error) {
		return executeNodeConfigDeployWorker(c, remoteCmd)
	})
	return renderNodeConfigDeploySummary(targets, opts.Except, len(envelope.Connections), opts.IsDryRun, opts.IsJSON)
}
