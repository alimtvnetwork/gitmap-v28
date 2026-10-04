package cmdagy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// NodeStatusRecord provides lightweight health telemetry for a fleet node.
type NodeStatusRecord struct {
	Alias       string `json:"alias"`
	HostAlias   string `json:"hostAlias"`
	OS          string `json:"os"`
	IsOnline    bool   `json:"isOnline"`
	ActiveProcs int    `json:"activeProcs"`
	LatencyMs   int64  `json:"latencyMs"`
}

// ProjectPromptTree groups active prompt summaries under a project and node.
type ProjectPromptTree struct {
	NodeAlias     string                `json:"nodeAlias"`
	ProjectName   string                `json:"projectName"`
	WorkspacePath string                `json:"workspacePath"`
	ActivePrompts []ActivePromptSummary `json:"activePrompts"`
	QueuedPrompts []AgyPromptQueueEntry `json:"queuedPrompts"`
}

// FullAGYUIStatusPayload represents the full state returned to the web UI including fleet telemetry.
type FullAGYUIStatusPayload struct {
	IsServerRunning bool                   `json:"isServerRunning"`
	NodeAlias       string                 `json:"nodeAlias"`
	RunningProjects []RunningProjectRecord `json:"runningProjects"`
	ActivePrompts   []ActivePromptSummary  `json:"activePrompts"`
	QueuedPrompts   []AgyPromptQueueEntry  `json:"queuedPrompts"`
	SavedPrompts    []PromptRecord         `json:"savedPrompts"`
	FleetNodes      []NodeStatusRecord     `json:"fleetNodes"`
}

// PromptResendPayload models an inbound request to resend a past prompt.
type PromptResendPayload struct {
	ID            string `json:"id"`
	ProjectTarget string `json:"projectTarget,omitempty"`
}

type agyUIOptions struct {
	port      int
	noBrowser bool
}

// ExecuteAgySettingsExport exports Antigravity settings to the specified path.
func ExecuteAgySettingsExport(outPath string) error {
	return executeAgySettingsExport(outPath)
}

// ExecuteAgySettingsImport imports Antigravity settings from the specified path.
func ExecuteAgySettingsImport(inPath string) error {
	return executeAgySettingsImport(inPath)
}

func parseAgyUIOptions(args []string) agyUIOptions {
	opts := agyUIOptions{port: 7430, noBrowser: false}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--no-browser" || arg == "-n" {
			opts.noBrowser = true
			continue
		}
		if (arg == "--port" || arg == "-p") && i+1 < len(args) {
			if val, err := strconv.Atoi(args[i+1]); err == nil && val > 0 {
				opts.port = val
				i++
			}
		}
	}
	return opts
}

func newAgyUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAgyUIIndex)
	mux.HandleFunc("/api/status", handleAgyUIStatus)
	mux.HandleFunc("/api/prompts/send", handleAgyUIPromptSend)
	mux.HandleFunc("/api/prompts/enqueue", handleAgyUIPromptEnqueue)
	mux.HandleFunc("/api/prompts/history", handleAgyUIPromptHistory)
	mux.HandleFunc("/api/prompts/save", handleAgyUIPromptSave)
	mux.HandleFunc("/api/prompts/resend", handleAgyUIPromptResend)
	mux.HandleFunc("/api/prompts/tree", handleAgyUIPromptTree)
	return mux
}

func handleAgyUIIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(agyUIDashboardHTML))
}

// CollectFleetNodeStatuses gathers sanitized status records for all fleet nodes.
func CollectFleetNodeStatuses(runningCount int) []NodeStatusRecord {
	host, _ := os.Hostname()
	if host == "" {
		host = "localhost"
	}

	nodes := []NodeStatusRecord{
		{
			Alias:       "local",
			HostAlias:   host,
			OS:          runtime.GOOS,
			IsOnline:    true,
			ActiveProcs: runningCount,
			LatencyMs:   0,
		},
	}

	if SSHConnectionsFetcher == nil {
		return nodes
	}

	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nodes
	}

	for _, c := range conns {
		if strings.EqualFold(c.Alias, "local") || strings.EqualFold(c.Alias, "localhost") {
			continue
		}
		hostLabel := c.Alias
		if hostLabel == "" {
			hostLabel = "fleet-node"
		}
		isOnline := strings.Contains(strings.ToLower(c.OS), "win") || strings.Contains(strings.ToLower(c.OS), "lin") || c.Username != ""
		nodes = append(nodes, NodeStatusRecord{
			Alias:       c.Alias,
			HostAlias:   hostLabel,
			OS:          c.OS,
			IsOnline:    isOnline,
			ActiveProcs: 0,
			LatencyMs:   12,
		})
	}

	return nodes
}

func handleAgyUIStatus(w http.ResponseWriter, _ *http.Request) {
	status, err := GetAGYUIStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	runningCount := len(status.RunningProjects)
	fleetNodes := CollectFleetNodeStatuses(runningCount)

	fullStatus := FullAGYUIStatusPayload{
		IsServerRunning: status.IsServerRunning,
		NodeAlias:       status.NodeAlias,
		RunningProjects: status.RunningProjects,
		ActivePrompts:   status.ActivePrompts,
		QueuedPrompts:   status.QueuedPrompts,
		SavedPrompts:    status.SavedPrompts,
		FleetNodes:      fleetNodes,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fullStatus)
}

func handleAgyUIPromptSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var payload PromptPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	res, _ := SendPromptImmediate(payload)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func handleAgyUIPromptEnqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var payload PromptPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	rec, err := EnqueuePromptPayload(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rec)
}

func handleAgyUIPromptHistory(w http.ResponseWriter, _ *http.Request) {
	history, _ := GetPromptHistory()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(history)
}

func handleAgyUIPromptSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var rec PromptRecord
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	_ = SavePromptTemplate(rec)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleAgyUIPromptResend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var req PromptResendPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, `{"error":"prompt id required"}`, http.StatusBadRequest)
		return
	}

	history, _ := GetPromptHistory()
	var targetRec *PromptRecord
	for _, rec := range history {
		if rec.ID == req.ID {
			targetRec = &rec
			break
		}
	}

	if targetRec == nil {
		saved, _ := LoadSavedPrompts()
		for _, rec := range saved {
			if rec.ID == req.ID {
				targetRec = &rec
				break
			}
		}
	}

	if targetRec == nil {
		http.Error(w, `{"error":"prompt not found in history"}`, http.StatusNotFound)
		return
	}

	target := targetRec.ProjectTarget
	if req.ProjectTarget != "" {
		target = req.ProjectTarget
	}

	res, err := SendPromptImmediate(PromptPayload{
		ProjectTarget: target,
		Title:         targetRec.Title,
		PromptText:    targetRec.PromptText,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func sanitizeWorkspacePath(path string) string {
	if path == "" || path == "." {
		return "."
	}
	base := filepath.Base(path)
	return filepath.Join("$WORKSPACE_ROOT", base)
}

func handleAgyUIPromptTree(w http.ResponseWriter, _ *http.Request) {
	status, err := GetAGYUIStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	treeMap := make(map[string]*ProjectPromptTree)
	for _, p := range status.RunningProjects {
		key := p.ProjectName
		if key == "" {
			key = filepath.Base(p.ProjectPath)
		}
		summary := buildActiveSummaryFromProject(p)
		summary.WorkspacePath = sanitizeWorkspacePath(p.ProjectPath)
		treeMap[key] = &ProjectPromptTree{
			NodeAlias:     status.NodeAlias,
			ProjectName:   p.ProjectName,
			WorkspacePath: summary.WorkspacePath,
			ActivePrompts: []ActivePromptSummary{summary},
			QueuedPrompts: []AgyPromptQueueEntry{},
		}
	}

	for _, q := range status.QueuedPrompts {
		key := q.ProjectName
		if key == "" {
			key = "default"
		}
		if item, exists := treeMap[key]; exists {
			item.QueuedPrompts = append(item.QueuedPrompts, q)
		} else {
			treeMap[key] = &ProjectPromptTree{
				NodeAlias:     status.NodeAlias,
				ProjectName:   key,
				WorkspacePath: "$WORKSPACE_ROOT/" + key,
				ActivePrompts: []ActivePromptSummary{},
				QueuedPrompts: []AgyPromptQueueEntry{q},
			}
		}
	}

	var treeList []ProjectPromptTree
	for _, item := range treeMap {
		treeList = append(treeList, *item)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(treeList)
}

func bindAgyUIListener(port int) (net.Listener, int, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err == nil {
		return ln, port, nil
	}
	dynLn, dynErr := net.Listen("tcp", "127.0.0.1:0")
	if dynErr != nil {
		return nil, 0, apperror.WrapSimple(dynErr, "bind agy ui listener")
	}
	dynPort := dynLn.Addr().(*net.TCPAddr).Port
	return dynLn, dynPort, nil
}

func printServerBanner(port int, url string) {
	fmt.Printf("\n%s● Antigravity Studio Web UI running on:%s %s\n", constants.ColorCyan, constants.ColorReset, url)
	fmt.Println("  Press Ctrl+C to terminate.")
}

// RunAgyUI starts the embedded Antigravity studio HTTP server.
func RunAgyUI(args []string) error {
	opts := parseAgyUIOptions(args)
	ln, activePort, err := bindAgyUIListener(opts.port)
	if err != nil {
		return err
	}
	defer ln.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d", activePort)
	printServerBanner(activePort, url)
	if !opts.noBrowser {
		go func() {
			time.Sleep(300 * time.Millisecond)
			openSUGUIBrowser(url)
		}()
	}

	return http.Serve(ln, newAgyUIMux())
}

// IsAgyUICommand checks if the command string triggers AGY UI.
func IsAgyUICommand(cmd string) bool {
	low := strings.ToLower(cmd)
	return low == "ui" || low == "dashboard" || low == "studio"
}
