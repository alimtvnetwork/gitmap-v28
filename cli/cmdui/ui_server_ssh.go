package cmdui

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
)

func handleAPISSHNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		summaries[i] = NodeSummary{
			NodeId:   c.Alias,
			Alias:    c.Alias,
			Host:     c.IPAddress,
			Port:     22,
			IsOnline: true,
			OS:       c.OS,
		}
	}
	return summaries
}

func handleAPISSHExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
		return nil, fmt.Errorf("invalid SSH nodes JSON format")
	}
	return conns, nil
}

func handleAPISSHImport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
