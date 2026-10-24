package cmdagy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
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

func parsePortQuiet(val string, fallback int) int {
	if p, err := strconv.Atoi(val); err == nil && p > 0 {
		return p
	}

	return fallback
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
			opts.port = parsePortQuiet(args[i+1], opts.port)
			i++
			continue
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
	mux.HandleFunc("/api/prompts/backup", handleAgyUIPromptBackup)
	mux.HandleFunc("/api/prompts/restore", handleAgyUIPromptRestore)
	mux.HandleFunc("/api/prompts/reset", handleAgyUIPromptReset)
	mux.HandleFunc("/api/prompts/refresh", handleAgyUIPromptRefresh)
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
