// Package cmdssh — ssh_deploy_config_ssh.go deploys SSH fleet configuration across machines.
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func loadDeployConfigEnvelope(filePath string) (*SSHNodesExportEnvelope, error) {
	if filePath == "" {
		return BuildSSHNodesExportEnvelope()
	}
	return loadDeployConfigFromFile(filePath)
}

func loadDeployConfigFromFile(filePath string) (*SSHNodesExportEnvelope, error) {
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
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalNodes:    len(conns),
		Nodes:         items,
		Connections:   conns,
	}
}

func executeDeployConfigToFleet(targets []db.SSHConnection, envelope *SSHNodesExportEnvelope, opts DeployConfigSSHOptions) error {
	compactBytes, _ := json.Marshal(envelope)

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
		return executeNodeConfigDeployWorker(c, compactBytes)
	})
	return renderNodeConfigDeploySummary(targets, opts.Except, len(envelope.Connections), opts.IsDryRun, opts.IsJSON)
}
