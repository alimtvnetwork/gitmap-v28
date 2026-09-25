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

func decodeEnvelopeConnections(raw []byte) ([]db.SSHConnection, bool) {
	var env SSHNodesExportEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, false
	}
	if len(env.Connections) > 0 {
		return env.Connections, true
	}
	if len(env.Nodes) > 0 {
		return convertExportItemsToConnections(env.Nodes), true
	}
	return nil, false
}

func decodeConnectionsFromJSON(raw []byte) ([]db.SSHConnection, error) {
	if conns, ok := decodeEnvelopeConnections(raw); ok {
		return conns, nil
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
	if inPath != DefaultSSHNodesJSONFile {
		return nil, inPath, err
	}
	altData, altErr := os.ReadFile(DefaultSSHNodesAltJSONFile)
	if altErr == nil {
		return altData, DefaultSSHNodesAltJSONFile, nil
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
