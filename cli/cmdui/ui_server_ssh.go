package cmdui

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SSHPingResult encapsulates the JSON response for node ping tests.
type SSHPingResult struct {
	Success   bool   `json:"success"`
	IsOnline  bool   `json:"isOnline"`
	LatencyMs int64  `json:"latencyMs"`
	Alias     string `json:"alias"`
	Error     string `json:"error,omitempty"`
}

func init() {
	http.HandleFunc("/api/ssh/test", handleAPISSHTest)
	http.HandleFunc("/api/ssh/ping", handleAPISSHPing)
}

// RegisterAPISSHRoutes mounts the SSH endpoints to an explicit ServeMux.
func RegisterAPISSHRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ssh/nodes", handleAPISSHNodes)
	mux.HandleFunc("/api/ssh/test", handleAPISSHTest)
	mux.HandleFunc("/api/ssh/ping", handleAPISSHPing)
	mux.HandleFunc("/api/ssh/export", handleAPISSHExport)
	mux.HandleFunc("/api/ssh/import", handleAPISSHImport)
	mux.HandleFunc("/api/ssh/deploy-keys", handleAPISSHDeployKeys)
}

func handleAPISSHNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.URL.Query().Get("action") == "test" || r.URL.Query().Get("action") == "ping" {
		handleAPISSHTest(w, r)
		return
	}

	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		_ = json.NewEncoder(w).Encode([]NodeSummary{})
		return
	}

	summaries := buildNodeSummaries(conns)
	_ = json.NewEncoder(w).Encode(summaries)
}

func buildNodeSummaries(conns []db.SSHConnection) []NodeSummary {
	summaries := make([]NodeSummary, len(conns))
	for i, c := range conns {
		account := c.Username
		if account == "" {
			account = c.Alias
		}
		summaries[i] = NodeSummary{
			NodeId:       c.Alias,
			Alias:        c.Alias,
			Host:         c.IPAddress,
			Port:         22,
			IsOnline:     false,
			OS:           c.OS,
			Account:      account,
			Instances:    1,
			InstanceList: []string{fmt.Sprintf("%s-instance-1", c.Alias)},
		}
	}
	return summaries
}

func extractAliasFromRequestBody(r *http.Request) string {
	if r.Method != http.MethodPost {
		return ""
	}
	var req struct {
		Alias     string `json:"alias"`
		NodeAlias string `json:"nodeAlias"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ""
	}
	if req.Alias != "" {
		return strings.TrimSpace(req.Alias)
	}
	return strings.TrimSpace(req.NodeAlias)
}

func handleAPISSHTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)

	alias := strings.TrimSpace(r.URL.Query().Get("alias"))
	if alias == "" {
		alias = strings.TrimSpace(r.URL.Query().Get("nodeAlias"))
	}

	if alias == "" {
		alias = extractAliasFromRequestBody(r)
	}

	if alias == "" {
		alias = "local"
	}

	res := executeNodePingCheck(alias)
	logFailedSSHPing(alias, res)
	_ = json.NewEncoder(w).Encode(res)
}

func logFailedSSHPing(alias string, res SSHPingResult) {
	isPingOk := res.Success
	if isPingOk {
		return
	}

	store.LogInternalErrorRecord(store.InternalErrorRecord{
		ErrorType:  "ssh",
		ErrorCode:  "ssh_ping_failed",
		Command:    "api/ssh/test",
		Message:    sshPingFailureMessage(res),
		Details:    alias,
		SourceFile: "cli/cmdui/ui_server_ssh.go",
	})
}

func sshPingFailureMessage(res SSHPingResult) string {
	message := strings.TrimSpace(res.Error)
	hasMessage := message != ""
	if hasMessage {
		return message
	}

	return "ssh ping failed"
}

func handleAPISSHPing(w http.ResponseWriter, r *http.Request) {
	handleAPISSHTest(w, r)
}

func executeNodePingCheck(alias string) SSHPingResult {
	start := time.Now()

	if strings.EqualFold(alias, "local") || strings.EqualFold(alias, constants.Localhost) || alias == "." {
		return SSHPingResult{
			Success:   true,
			IsOnline:  true,
			LatencyMs: 0,
			Alias:     alias,
		}
	}

	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return SSHPingResult{
			Success:   false,
			IsOnline:  false,
			LatencyMs: time.Since(start).Milliseconds(),
			Alias:     alias,
			Error:     err.Error(),
		}
	}

	var targetConn *db.SSHConnection
	for _, c := range conns {
		if strings.EqualFold(c.Alias, alias) || strings.EqualFold(c.IPAddress, alias) {
			targetConn = &c
			break
		}
	}

	if targetConn == nil {
		return SSHPingResult{
			Success:   false,
			IsOnline:  false,
			LatencyMs: time.Since(start).Milliseconds(),
			Alias:     alias,
			Error:     fmt.Sprintf("node %q not found in SSH fleet", alias),
		}
	}

	dialAddr := net.JoinHostPort(targetConn.IPAddress, "22")
	conn, dialErr := net.DialTimeout("tcp", dialAddr, 2*time.Second)
	latency := time.Since(start).Milliseconds()

	if dialErr != nil {
		return SSHPingResult{
			Success:   false,
			IsOnline:  false,
			LatencyMs: latency,
			Alias:     alias,
			Error:     dialErr.Error(),
		}
	}
	_ = conn.Close()

	return SSHPingResult{
		Success:   true,
		IsOnline:  true,
		LatencyMs: latency,
		Alias:     alias,
	}
}

func handleAPISSHExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	envelope, err := cmdssh.BuildSSHNodesExportEnvelope()
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}
	typedEnv := jsonx.NewEnvelope(
		jsonx.TypeSSHNodes,
		"gitmap-ssh-nodes.json",
		"gitmap ssh export",
		"1.0",
		envelope,
	)
	_ = json.NewEncoder(w).Encode(typedEnv)
}

func parseConnectionsFromRaw(raw []byte) ([]db.SSHConnection, error) {
	payload, _, extractErr := jsonx.ExtractPayload(raw)
	if extractErr == nil && len(payload) > 0 {
		raw = payload
	}
	var env cmdssh.SSHNodesExportEnvelope
	if err := json.Unmarshal(raw, &env); err == nil && len(env.Connections) > 0 {
		return env.Connections, nil
	}
	var conns []db.SSHConnection
	if err := json.Unmarshal(raw, &conns); err != nil {
		return nil, apperror.NewValidationError("invalid SSH nodes JSON format")
	}
	return conns, nil
}

func handleAPISSHImport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		JSONContent string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}
	conns, parseErr := parseConnectionsFromRaw([]byte(req.JSONContent))
	if parseErr != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": parseErr.Error()})
		return
	}
	stats, syncErr := cmdssh.SyncSSHConnectionsLocally(conns)
	if syncErr != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": syncErr.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "stats": stats})
}

func handleAPISSHDeployKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	err := cmdssh.RunSSHDeployKeysCLI([]string{"all"})
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Public keys deployed to fleet"})
}
