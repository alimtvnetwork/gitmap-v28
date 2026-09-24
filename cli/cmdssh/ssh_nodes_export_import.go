// Package cmdssh — ssh_nodes_export_import.go handles SSH nodes JSON export/import, one-liner generation, and fleet node-config deployment.
package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// DefaultSSHNodesJSONFile is the default filename for SSH node JSON exports and imports.
const (
	DefaultSSHNodesJSONFile    = "gitmap-ssh-nodes.json"
	DefaultSSHNodesAltJSONFile = "gitmap-ssh.json"
)

// SSHNodeExportItem represents an exported SSH node with deterministic worker ID and metadata.
type SSHNodeExportItem struct {
	WorkerID   string `json:"worker_id"`
	NumericID  int    `json:"id"`
	Alias      string `json:"alias"`
	IPAddress  string `json:"ip_address"`
	Username   string `json:"username"`
	Port       int    `json:"port"`
	OS         string `json:"os"`
	AuthMethod string `json:"auth_method"`
	KeyPath    string `json:"key_path,omitempty"`
}

// SSHNodesExportEnvelope wraps the exported SSH nodes list with schema and timestamp metadata.
type SSHNodesExportEnvelope struct {
	SchemaVersion string              `json:"schema_version"`
	ExportedAt    string              `json:"exported_at"`
	TotalNodes    int                 `json:"total_nodes"`
	Nodes         []SSHNodeExportItem `json:"nodes"`
	Connections   []db.SSHConnection  `json:"connections"`
}

// RunSSHNodesExportJSON exports all registered SSH nodes to a JSON file (default: gitmap-ssh-nodes.json and gitmap-ssh.json).
func RunSSHNodesExportJSON(args []string) error {
	outPath, isDefaultFile, toStdout := parseNodesExportPathArg(args)
	envelope, err := BuildSSHNodesExportEnvelope()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return err
	}
	if toStdout {
		fmt.Println(string(data))
		return nil
	}
	if dir := filepath.Dir(outPath); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return err
	}
	if isDefaultFile {
		_ = os.WriteFile(DefaultSSHNodesAltJSONFile, data, 0644)
	}
	fmt.Printf("✓ Exported %d SSH node(s) to %s\n", len(envelope.Connections), outPath)
	return nil
}

// BuildSSHNodesExportEnvelope constructs the portable SSH nodes export envelope.
func BuildSSHNodesExportEnvelope() (*SSHNodesExportEnvelope, error) {
	conns, err := fetchAllSSHConnections()
	if err != nil {
		conns = []db.SSHConnection{}
	}
	items := make([]SSHNodeExportItem, 0, len(conns))
	for idx, c := range conns {
		numID := idx + 1
		authMethod := "key"
		if c.EncryptedPassword != "" && c.KeyPath == "" {
			authMethod = "password"
		}
		items = append(items, SSHNodeExportItem{
			WorkerID:   fmt.Sprintf("worker-%d", idx+1),
			NumericID:  numID,
			Alias:      c.Alias,
			IPAddress:  c.IPAddress,
			Username:   c.Username,
			Port:       22,
			OS:         c.OS,
			AuthMethod: authMethod,
			KeyPath:    c.KeyPath,
		})
	}
	return &SSHNodesExportEnvelope{
		SchemaVersion: "1.0",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalNodes:    len(conns),
		Nodes:         items,
		Connections:   conns,
	}, nil
}

// RunSSHNodesImportJSON imports SSH nodes from a JSON file (default: gitmap-ssh-nodes.json or gitmap-ssh.json) or --base64 payload.
func RunSSHNodesImportJSON(args []string) error {
	b64Payload, inPath := parseImportJSONArgs(args)
	var raw []byte
	var err error
	sourceLabel := inPath
	if b64Payload != "" {
		raw, err = base64.StdEncoding.DecodeString(strings.TrimSpace(b64Payload))
		sourceLabel = "inline-base64-oneliner"
	} else {
		raw, sourceLabel, err = readNodesImportFileWithFallback(inPath)
	}
	if err != nil {
		return fmt.Errorf("failed to read SSH nodes JSON from %s: %w", sourceLabel, err)
	}
	conns, parseErr := decodeConnectionsFromJSON(raw)
	if parseErr != nil {
		return parseErr
	}
	imported := importConnectionsLocally(conns)
	fmt.Printf("✓ Imported %d/%d SSH node(s) from %s\n", imported, len(conns), sourceLabel)
	return nil
}

func decodeConnectionsFromJSON(raw []byte) ([]db.SSHConnection, error) {
	var env SSHNodesExportEnvelope
	if err := json.Unmarshal(raw, &env); err == nil && (len(env.Connections) > 0 || len(env.Nodes) > 0) {
		if len(env.Connections) > 0 {
			return env.Connections, nil
		}
		return convertExportItemsToConnections(env.Nodes), nil
	}
	var directConns []db.SSHConnection
	if err := json.Unmarshal(raw, &directConns); err != nil {
		return nil, err
	}
	return directConns, nil
}

func convertExportItemsToConnections(items []SSHNodeExportItem) []db.SSHConnection {
	out := make([]db.SSHConnection, 0, len(items))
	for _, it := range items {
		out = append(out, db.SSHConnection{
			Alias:     it.Alias,
			IPAddress: it.IPAddress,
			Username:  it.Username,
			OS:        it.OS,
			KeyPath:   it.KeyPath,
		})
	}
	return out
}

// RunSSHExportOnelinerCLI generates a single-line command containing all SSH nodes encoded in Base64.
func RunSSHExportOnelinerCLI(args []string) error {
	envelope, err := BuildSSHNodesExportEnvelope()
	if err != nil {
		return err
	}
	compactBytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	b64 := base64.StdEncoding.EncodeToString(compactBytes)
	fmt.Printf("gitmap ssh nodes import-json --base64 \"%s\"\n", b64)
	if !hasOnelinerRawFlag(args) {
		fmt.Printf("\n%s# PowerShell File + Import One-Liner:%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("[IO.File]::WriteAllBytes('%s',[Convert]::FromBase64String('%s')); gitmap ssh nodes import-json %s\n", DefaultSSHNodesJSONFile, b64, DefaultSSHNodesJSONFile)
		fmt.Printf("\n%s# POSIX sh/bash File + Import One-Liner:%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("printf '%%s' '%s' | base64 -d > %s && gitmap ssh nodes import-json %s\n", b64, DefaultSSHNodesJSONFile, DefaultSSHNodesJSONFile)
	}
	return nil
}

// RunSSHDeployNodeConfigCLI deploys the local SSH node configuration across all online remote fleet nodes (--except id,ip,alias).
func RunSSHDeployNodeConfigCLI(args []string) error {
	except, isDryRun, isJSON := parseNodeConfigDeployFlags(args)
	envelope, err := BuildSSHNodesExportEnvelope()
	if err != nil {
		return err
	}
	compactBytes, _ := json.Marshal(envelope)
	b64 := base64.StdEncoding.EncodeToString(compactBytes)
	remoteCmd := fmt.Sprintf("gitmap ssh nodes import-json --base64 \"%s\"", b64)

	conns, _ := fetchAllSSHConnections()
	filtered := FilterSSHConnectionsByExcept(conns, except)
	if isDryRun || len(filtered) == 0 {
		return renderNodeConfigDeploySummary(filtered, except, len(envelope.Connections), isDryRun, isJSON)
	}
	opts := FleetParallelOptions{
		Target:   "all",
		Except:   except,
		TaskName: "Deploy SSH Node-Config (nc)",
	}
	RunParallelFleetExecution(filtered, opts, func(c db.SSHConnection) (string, error) {
		return executeNodeConfigDeployWorker(c, remoteCmd)
	})
	return renderNodeConfigDeploySummary(filtered, except, len(envelope.Connections), isDryRun, isJSON)
}

func executeNodeConfigDeployWorker(c db.SSHConnection, remoteCmd string) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return "", fmt.Errorf("node %s is offline", header)
	}
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return "", fmt.Errorf("auth failed for %s", header)
	}
	defer client.Close()
	return crypto.RunCommand(client, remoteCmd, "")
}

// FilterSSHConnectionsByExcept excludes connections matching any comma/space-separated ID (1, worker-1), IP, or alias.
func FilterSSHConnectionsByExcept(conns []db.SSHConnection, exceptRaw string) []db.SSHConnection {
	tokens := splitExceptTokens(exceptRaw)
	if len(tokens) == 0 {
		return conns
	}
	var out []db.SSHConnection
	for idx, c := range conns {
		if isConnectionExcluded(c, idx+1, tokens) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func isConnectionExcluded(c db.SSHConnection, oneBasedIdx int, tokens []string) bool {
	workerID := fmt.Sprintf("worker-%d", oneBasedIdx)
	idxStr := strconv.Itoa(oneBasedIdx)
	for _, tok := range tokens {
		low := strings.ToLower(strings.TrimSpace(tok))
		if low == "" {
			continue
		}
		if low == idxStr || low == strings.ToLower(workerID) ||
			strings.EqualFold(c.Alias, low) || strings.EqualFold(c.IPAddress, low) {
			return true
		}
	}
	return false
}

func splitExceptTokens(raw string) []string {
	var out []string
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	}) {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseNodesExportPathArg(args []string) (string, bool, bool) {
	var pathArg string
	toStdout := false
	for _, a := range args {
		if a == "--stdout" || a == "-o-" {
			toStdout = true
			continue
		}
		if !strings.HasPrefix(a, "-") && pathArg == "" {
			pathArg = a
		}
	}
	if pathArg == "" {
		return DefaultSSHNodesJSONFile, true, toStdout
	}
	if info, err := os.Stat(pathArg); err == nil && info.IsDir() {
		return filepath.Join(pathArg, DefaultSSHNodesJSONFile), false, toStdout
	}
	return pathArg, false, toStdout
}

func readNodesImportFileWithFallback(inPath string) ([]byte, string, error) {
	if info, err := os.Stat(inPath); err == nil && info.IsDir() {
		inPath = filepath.Join(inPath, DefaultSSHNodesJSONFile)
	}
	data, err := os.ReadFile(inPath)
	if err == nil {
		return data, inPath, nil
	}
	if inPath == DefaultSSHNodesJSONFile {
		if altData, altErr := os.ReadFile(DefaultSSHNodesAltJSONFile); altErr == nil {
			return altData, DefaultSSHNodesAltJSONFile, nil
		}
	}
	return nil, inPath, err
}

func parseImportJSONArgs(args []string) (string, string) {
	var b64, pathArg string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--base64" || a == "-b64") && i+1 < len(args) {
			b64 = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--base64=") {
			b64 = strings.TrimPrefix(a, "--base64=")
			continue
		}
		if !strings.HasPrefix(a, "-") && pathArg == "" {
			pathArg = a
		}
	}
	if pathArg == "" {
		pathArg = DefaultSSHNodesJSONFile
	}
	return b64, pathArg
}

func parseNodeConfigDeployFlags(args []string) (string, bool, bool) {
	var exceptParts []string
	isDryRun, isJSON := false, false
	for i := 0; i < len(args); i++ {
		a := args[i]
		low := strings.ToLower(a)
		if low == "node-config" || low == "nc" || low == "deploy" || low == "ssh" {
			continue
		}
		if low == "--dry-run" || low == "-n" {
			isDryRun = true
			continue
		}
		if low == "--json" {
			isJSON = true
			continue
		}
		if isExceptOrExcepFlag(low) {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				exceptParts = append(exceptParts, args[i+1])
				i++
			}
			continue
		}
		if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--excep=") || strings.HasPrefix(low, "--accept=") || strings.HasPrefix(low, "--exclude=") {
			idx := strings.IndexByte(a, '=')
			exceptParts = append(exceptParts, a[idx+1:])
		}
	}
	return strings.Join(exceptParts, ","), isDryRun, isJSON
}

func isExceptOrExcepFlag(low string) bool {
	return low == "--except" || low == "--excep" || low == "--accept" || low == "--exclude" || low == "-e"
}

func hasOnelinerRawFlag(args []string) bool {
	for _, a := range args {
		if a == "--raw" || a == "-q" || a == "--quiet" {
			return true
		}
	}
	return false
}

func renderNodeConfigDeploySummary(targets []db.SSHConnection, except string, totalConfigNodes int, isDryRun, isJSON bool) error {
	var names []string
	for _, t := range targets {
		names = append(names, fmt.Sprintf("%s(%s)", t.Alias, t.IPAddress))
	}
	if isJSON {
		payload := map[string]any{
			"command":             "ssh deploy node-config",
			"config_nodes_synced": totalConfigNodes,
			"target_nodes":        names,
			"target_count":        len(targets),
			"except":              except,
			"dry_run":             isDryRun,
		}
		b, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(b))
		return nil
	}
	fmt.Printf("✓ Deployed SSH node-config (%d node definitions) to %d target machine(s) [except=%q]: %s\n",
		totalConfigNodes, len(targets), except, strings.Join(names, ", "))
	return nil
}
