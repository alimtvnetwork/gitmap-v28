// Package cmdssh — ssh_nodes_export_import.go handles SSH nodes JSON export/import, one-liner generation, and fleet node-config deployment.
package cmdssh

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
)

// DefaultSSHNodesJSONFile is the default filename for SSH node JSON exports and imports.
const (
	DefaultSSHNodesJSONFile    = "gitmap-ssh-nodes.json"
	DefaultSSHNodesAltJSONFile = "gitmap-ssh.json"
)

// SSHNodeExportItem represents an exported SSH node with deterministic worker ID and metadata.
type SSHNodeExportItem struct {
	WorkerId          string `json:"workerId"`
	NumericId         int    `json:"id"`
	Alias             string `json:"alias"`
	IPAddress         string `json:"ipAddress"`
	Username          string `json:"username"`
	Port              int    `json:"port"`
	OS                string `json:"os"`
	AuthMethod        string `json:"authMethod"`
	AuthType          string `json:"authType,omitempty"`
	KeyPath           string `json:"keyPath,omitempty"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"encryptedPassword,omitempty"`
}

// UnmarshalJSON implements custom decoding for both camelCase and legacy snake_case SSH node items.
func (item *SSHNodeExportItem) UnmarshalJSON(data []byte) error {
	type Alias SSHNodeExportItem
	var raw struct {
		Alias
		LegacyWorkerId          string `json:"worker_id"`
		LegacyIPAddress         string `json:"ip_address"`
		LegacyAuthMethod        string `json:"auth_method"`
		LegacyAuthType          string `json:"auth_type"`
		LegacyKeyPath           string `json:"key_path"`
		LegacyEncryptedPassword string `json:"encrypted_password"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*item = SSHNodeExportItem(raw.Alias)
	item.applyLegacyFallbacks(raw.LegacyWorkerId, raw.LegacyIPAddress, raw.LegacyAuthMethod, raw.LegacyAuthType, raw.LegacyKeyPath, raw.LegacyEncryptedPassword)
	return nil
}

func (item *SSHNodeExportItem) applyLegacyFallbacks(wId, ip, aMethod, aType, kPath, encPass string) {
	if item.WorkerId == "" && wId != "" {
		item.WorkerId = wId
	}
	if item.IPAddress == "" && ip != "" {
		item.IPAddress = ip
	}
	if item.AuthMethod == "" && aMethod != "" {
		item.AuthMethod = aMethod
	}
	if item.AuthMethod == "" && aType != "" {
		item.AuthMethod = aType
	}
	if item.AuthType == "" && item.AuthMethod != "" {
		item.AuthType = item.AuthMethod
	}
	if item.KeyPath == "" && kPath != "" {
		item.KeyPath = kPath
	}
	if item.EncryptedPassword == "" && encPass != "" {
		item.EncryptedPassword = encPass
	}
}

// SSHNodesExportEnvelope wraps the exported SSH nodes list with schema and timestamp metadata.
type SSHNodesExportEnvelope struct {
	SchemaVersion string              `json:"schemaVersion"`
	ExportedAt    string              `json:"exportedAt"`
	TotalNodes    int                 `json:"totalNodes"`
	Nodes         []SSHNodeExportItem `json:"nodes"`
	Connections   []db.SSHConnection  `json:"connections,omitempty"`
}

// UnmarshalJSON implements custom decoding for both camelCase and legacy snake_case envelope headers.
func (env *SSHNodesExportEnvelope) UnmarshalJSON(data []byte) error {
	type Alias SSHNodesExportEnvelope
	var raw struct {
		Alias
		LegacySchemaVersion string `json:"schema_version"`
		LegacyExportedAt    string `json:"exported_at"`
		LegacyTotalNodes    int    `json:"total_nodes"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*env = SSHNodesExportEnvelope(raw.Alias)
	if env.SchemaVersion == "" && raw.LegacySchemaVersion != "" {
		env.SchemaVersion = raw.LegacySchemaVersion
	}
	if env.ExportedAt == "" && raw.LegacyExportedAt != "" {
		env.ExportedAt = raw.LegacyExportedAt
	}
	if env.TotalNodes == 0 && raw.LegacyTotalNodes != 0 {
		env.TotalNodes = raw.LegacyTotalNodes
	}
	return nil
}

// RunSSHNodesExportJSON exports all registered SSH nodes to a JSON file (default: gitmap-ssh-nodes.json and gitmap-ssh.json).
func RunSSHNodesExportJSON(args []string) error {
	outPath, isDefaultFile, toStdout := parseNodesExportPathArg(args)
	envelope, err := BuildSSHNodesExportEnvelope()
	if err != nil {
		return err
	}
	typedEnv := jsonenvelope.NewEnvelope(
		jsonenvelope.TypeSSHNodes,
		outPath,
		"gitmap ssh export",
		"2.0",
		envelope,
	)
	data, err := json.MarshalIndent(typedEnv, "", "  ")
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
		numId := idx + 1
		authMethod := "key"
		if c.EncryptedPassword != "" && c.KeyPath == "" {
			authMethod = "password"
		}
		items = append(items, SSHNodeExportItem{
			WorkerId:   fmt.Sprintf("worker-%d", idx+1),
			NumericId:  numId,
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
		SchemaVersion: "2.0",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalNodes:    len(conns),
		Nodes:         items,
		Connections:   nil,
	}, nil
}

// BuildCompactSSHNodesExportEnvelope builds an envelope without raw db connections for clean, single-line exports.
func BuildCompactSSHNodesExportEnvelope() (*SSHNodesExportEnvelope, error) {
	env, err := BuildSSHNodesExportEnvelope()
	if err != nil {
		return nil, err
	}
	env.Connections = nil
	return env, nil
}

// RunSSHNodesImportJSON imports SSH nodes from a JSON file (default: gitmap-ssh-nodes.json or gitmap-ssh.json) or --base64 payload.
func RunSSHNodesImportJSON(args []string) error {
	b64Payload, inPath := parseImportJSONArgs(args)
	var raw []byte
	var err error
	sourceLabel := inPath
	if b64Payload != "" {
		cleaned := cleanBase64Payload(b64Payload)
		raw, err = base64.StdEncoding.DecodeString(cleaned)
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
	stats, syncErr := SyncSSHConnectionsLocally(conns)
	if syncErr != nil {
		return syncErr
	}
	fmt.Printf("✓ Synced %d SSH node(s) from %s: %d matched (%d updated, %d unchanged), %d inserted\n",
		stats.Total, sourceLabel, stats.Matched, stats.Updated, stats.Unchanged, stats.Inserted)
	_ = printSJList(context.Background(), os.Stdout, 0)
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
	payload, _, err := jsonenvelope.ExtractPayload(raw)
	if err == nil && len(payload) > 0 {
		raw = payload
	}
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
			Alias:             it.Alias,
			IPAddress:         it.IPAddress,
			Username:          it.Username,
			EncryptedPassword: resolveExportItemPassword(it.EncryptedPassword, it.Password),
			OS:                it.OS,
			KeyPath:           it.KeyPath,
		})
	}
	return out
}

func resolveExportItemPassword(encPass, plainPass string) string {
	if encPass != "" || plainPass == "" {
		return encPass
	}
	if encrypted, err := EncryptSSHPassword(plainPass); err == nil {
		return encrypted
	}
	return ""
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
	workerId := fmt.Sprintf("worker-%d", oneBasedIdx)
	idxStr := strconv.Itoa(oneBasedIdx)
	for _, tok := range tokens {
		low := strings.ToLower(strings.TrimSpace(tok))
		if low == "" {
			continue
		}
		if low == idxStr || low == strings.ToLower(workerId) ||
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
	trimmed := strings.TrimSpace(inPath)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return []byte(trimmed), "inline-json", nil
	}
	resolvedPath := ResolveRepoSecretsNodesPath(trimmed)
	if data, err := os.ReadFile(resolvedPath); err == nil {
		return data, resolvedPath, nil
	}
	if data, dirPath, isOk := tryReadDefaultNodesInDir(inPath); isOk {
		return data, dirPath, nil
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

func tryReadDefaultNodesInDir(inPath string) ([]byte, string, bool) {
	info, err := os.Stat(inPath)
	if err != nil || !info.IsDir() {
		return nil, "", false
	}
	dirPath := filepath.Join(inPath, DefaultSSHNodesJSONFile)
	data, dirErr := os.ReadFile(dirPath)
	if dirErr != nil {
		return nil, "", false
	}
	return data, dirPath, true
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
	if b64 == "" && isLikelyBase64Payload(pathArg) {
		b64 = pathArg
		pathArg = ""
	}
	if pathArg == "" {
		pathArg = DefaultSSHNodesJSONFile
	}
	return b64, pathArg
}

func isLikelyBase64Payload(s string) bool {
	clean := strings.Trim(strings.TrimSpace(s), "\"'`")
	return strings.HasPrefix(clean, "eyJ") && len(clean) > 20
}

func cleanBase64Payload(raw string) string {
	trimmed := strings.Trim(strings.TrimSpace(raw), "\"'`")
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' || r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, trimmed)
}
