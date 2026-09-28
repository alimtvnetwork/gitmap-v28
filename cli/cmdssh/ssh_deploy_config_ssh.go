// Package cmdssh — ssh_deploy_config_ssh.go implements 'gitmap deploy config ssh'.
package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHDeployConfigSSHCLI deploys SSH node topology and credentials across fleet nodes.
func RunSSHDeployConfigSSHCLI(args []string) error {
	target, filePath, except, isDryRun, isJSON, isHelp := parseDeployConfigSSHFlags(args)
	if isHelp {
		printDeployConfigSSHHelp()
		return nil
	}

	envelope, err := resolveSSHDeployEnvelope(filePath)
	if err != nil {
		return err
	}

	compactBytes, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to encode SSH nodes envelope: %w", err)
	}
	remoteCmd := fmt.Sprintf("gitmap ssh nodes import-json --base64 \"%s\"", base64.StdEncoding.EncodeToString(compactBytes))

	conns, _ := fetchAllSSHConnections()
	filtered := resolveDeployConfigTargets(conns, target, except)
	if isDryRun || len(filtered) == 0 {
		return renderNodeConfigDeploySummary(filtered, except, len(envelope.Connections), isDryRun, isJSON)
	}

	opts := FleetParallelOptions{
		Target:   target,
		Except:   except,
		TaskName: "Deploy SSH Config",
	}
	RunParallelFleetExecution(filtered, opts, func(c db.SSHConnection) (string, error) {
		return executeNodeConfigDeployWorker(c, remoteCmd)
	})
	return renderNodeConfigDeploySummary(filtered, except, len(envelope.Connections), isDryRun, isJSON)
}

func resolveSSHDeployEnvelope(filePath string) (*SSHNodesExportEnvelope, error) {
	if filePath == "" {
		return BuildSSHNodesExportEnvelope()
	}
	raw, sourceLabel, readErr := readNodesImportFileWithFallback(filePath)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read SSH config from %s: %w", sourceLabel, readErr)
	}
	conns, parseErr := decodeConnectionsFromJSON(raw)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse SSH config from %s: %w", sourceLabel, parseErr)
	}
	return &SSHNodesExportEnvelope{
		SchemaVersion: "1.0",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalNodes:    len(conns),
		Nodes:         convertConnectionsToExportItems(conns),
		Connections:   conns,
	}, nil
}

func convertConnectionsToExportItems(conns []db.SSHConnection) []SSHNodeExportItem {
	items := make([]SSHNodeExportItem, 0, len(conns))
	for idx, c := range conns {
		items = append(items, SSHNodeExportItem{
			WorkerID:   fmt.Sprintf("worker-%d", idx+1),
			NumericID:  idx + 1,
			Alias:      c.Alias,
			IPAddress:  c.IPAddress,
			Username:   c.Username,
			Port:       22,
			OS:         c.OS,
			AuthMethod: "key",
			KeyPath:    c.KeyPath,
		})
	}
	return items
}
