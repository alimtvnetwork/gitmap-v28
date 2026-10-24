package cmdschedule

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	// DefaultScheduleUIPort starts sequential probing for the schedule dashboard.
	DefaultScheduleUIPort = 8095
	// SchedulePortScanAttempts specifies the maximum port search iterations.
	SchedulePortScanAttempts = 50
)

// ScheduleUIOptions configures the schedule UI HTTP server.
type ScheduleUIOptions struct {
	Host string
	Port int
}

// ScheduleTaskItem represents a cron or interval scheduled task.
type ScheduleTaskItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	CommandLine string `json:"commandLine"`
	Interval    string `json:"interval"`
	IsEnabled   bool   `json:"isEnabled"`
	IsScheduled bool   `json:"isScheduled"`
	RunCount    int    `json:"runCount"`
	LastRunAt   string `json:"lastRunAt,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

// TaskToggleRequest represents payload to toggle task active state.
type TaskToggleRequest struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"isEnabled"`
}

// TaskToggleResponse returns outcome of task state change.
type TaskToggleResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Name      string `json:"name"`
	IsEnabled bool   `json:"isEnabled"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

// FolderWatcherItem defines a filesystem watcher trigger.
type FolderWatcherItem struct {
	ID              string   `json:"id"`
	Path            string   `json:"path"`
	Events          []string `json:"events"`
	DebounceSeconds int      `json:"debounceSeconds"`
	ActionCommand   string   `json:"actionCommand"`
	IsEnabled       bool     `json:"isEnabled"`
}

// WatcherCreateRequest specifies payload to register a folder watcher.
type WatcherCreateRequest struct {
	Path            string   `json:"path"`
	Events          []string `json:"events"`
	DebounceSeconds int      `json:"debounceSeconds"`
	ActionCommand   string   `json:"actionCommand"`
	IsEnabled       bool     `json:"isEnabled"`
}

// WatcherCreateResponse returns watcher creation result.
type WatcherCreateResponse struct {
	IsSuccess bool              `json:"isSuccess"`
	Watcher   FolderWatcherItem `json:"watcher"`
	Message   string            `json:"message"`
	Error     string            `json:"error,omitempty"`
}

// NotifierRecipientItem defines a recipient for alerts.
type NotifierRecipientItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	Events    string `json:"events"`
	IsEnabled bool   `json:"isEnabled"`
}

var (
	uiWatcherLock  sync.RWMutex
	uiWatchers     = []FolderWatcherItem{{ID: "w-01", Path: "./cli", Events: []string{"modify", "create"}, DebounceSeconds: 3, ActionCommand: "gitmap pe -t", IsEnabled: true}}
	uiNotifierLock sync.RWMutex
	uiNotifiers    = []NotifierRecipientItem{{ID: "n-01", Type: "email", Target: "devops-alerts@internal.company.com", Events: "failure,alert", IsEnabled: true}, {ID: "n-02", Type: "webhook", Target: "https://hooks.slack.com/services/T00/B00/X00", Events: "all", IsEnabled: true}}
	uiTaskLock     sync.RWMutex
	uiTasks        = []ScheduleTaskItem{{ID: 1, Name: "daily-fleet-sync", Slug: "daily-fleet-sync", CommandLine: "gitmap sync --fleet", Interval: "@daily", IsEnabled: true, IsScheduled: true, RunCount: 48, LastRunAt: "2026-10-20 04:00:00", CreatedAt: "2026-09-01 00:00:00"}, {ID: 2, Name: "hourly-cache-clean", Slug: "hourly-cache-clean", CommandLine: "gitmap clean --os-caches", Interval: "@hourly", IsEnabled: true, IsScheduled: true, RunCount: 216, LastRunAt: "2026-10-20 15:00:00", CreatedAt: "2026-09-10 12:00:00"}}
)

// BuildScheduleUIMux registers all dashboard and API routes.
func BuildScheduleUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleScheduleIndex)
	mux.HandleFunc("/ui", handleScheduleIndex)
	mux.HandleFunc("/api/schedule/tasks", handleAPIScheduleTasks)
	mux.HandleFunc("/api/schedule/watchers", handleAPIScheduleWatchers)
	mux.HandleFunc("/api/schedule/toggle", handleAPIScheduleToggle)
	mux.HandleFunc("/api/schedule/notifiers", handleAPIScheduleNotifiers)
	return mux
}

// ScanAvailableScheduleListener probes 50 sequential ports starting from preferred port.
func ScanAvailableScheduleListener(host string, startPort int) (net.Listener, int, error) {
	bindHost := "127.0.0.1"
	if host != "" {
		bindHost = host
	}
	targetPort := DefaultScheduleUIPort
	if startPort > 0 {
		targetPort = startPort
	}
	for port := targetPort; port < targetPort+SchedulePortScanAttempts; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindHost, port))
		if err == nil && ln != nil {
			return ln, port, nil
		}
	}
	return bindFallbackSchedulePort(bindHost)
}

func bindFallbackSchedulePort(host string) (net.Listener, int, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	if err == nil && ln != nil {
		return ln, ln.Addr().(*net.TCPAddr).Port, nil
	}
	return nil, 0, apperror.WrapSimple(err, "failed to bind schedule UI on any port")
}

// StartScheduleUIServer launches the schedule dashboard HTTP server.
func StartScheduleUIServer(opts ScheduleUIOptions) (*http.Server, net.Listener, int, error) {
	listener, port, err := ScanAvailableScheduleListener(opts.Host, opts.Port)
	if err != nil {
		return nil, nil, 0, err
	}
	server := &http.Server{Handler: BuildScheduleUIMux(), ReadHeaderTimeout: 5 * time.Second}
	return server, listener, port, nil
}

func handleScheduleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/ui" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(BuildScheduleDashboardHTML()))
}

func handleAPIScheduleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	uiTaskLock.RLock()
	defer uiTaskLock.RUnlock()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(uiTasks)
}

func handleAPIScheduleToggle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		writeScheduleJSONError(w, http.StatusMethodNotAllowed, "POST method required")
		return
	}
	var req TaskToggleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeScheduleJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	resp := executeTaskToggle(req)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func executeTaskToggle(req TaskToggleRequest) TaskToggleResponse {
	uiTaskLock.Lock()
	defer uiTaskLock.Unlock()
	for i := range uiTasks {
		if uiTasks[i].Name == req.Name || uiTasks[i].Slug == req.Name {
			uiTasks[i].IsEnabled = req.IsEnabled
			return TaskToggleResponse{IsSuccess: true, Name: req.Name, IsEnabled: req.IsEnabled, Message: fmt.Sprintf("Task '%s' state set to %t", req.Name, req.IsEnabled)}
		}
	}
	return TaskToggleResponse{IsSuccess: true, Name: req.Name, IsEnabled: req.IsEnabled, Message: fmt.Sprintf("Task '%s' updated", req.Name)}
}

func handleAPIScheduleWatchers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		uiWatcherLock.RLock()
		defer uiWatcherLock.RUnlock()
		_ = json.NewEncoder(w).Encode(uiWatchers)
		return
	}
	if r.Method != http.MethodPost {
		writeScheduleJSONError(w, http.StatusMethodNotAllowed, "GET or POST required")
		return
	}
	processCreateWatcher(w, r)
}

func processCreateWatcher(w http.ResponseWriter, r *http.Request) {
	var req WatcherCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeScheduleJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Path == "" {
		writeScheduleJSONError(w, http.StatusBadRequest, "Folder path is required")
		return
	}
	watcher := appendWatcher(req)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(WatcherCreateResponse{IsSuccess: true, Watcher: watcher, Message: "Folder watcher trigger registered"})
}

func appendWatcher(req WatcherCreateRequest) FolderWatcherItem {
	uiWatcherLock.Lock()
	defer uiWatcherLock.Unlock()
	debounce := req.DebounceSeconds
	if debounce <= 0 {
		debounce = 3
	}
	item := FolderWatcherItem{
		ID:              fmt.Sprintf("w-%d", time.Now().UnixNano()),
		Path:            req.Path,
		Events:          req.Events,
		DebounceSeconds: debounce,
		ActionCommand:   req.ActionCommand,
		IsEnabled:       req.IsEnabled,
	}
	uiWatchers = append(uiWatchers, item)
	return item
}

func handleAPIScheduleNotifiers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	uiNotifierLock.RLock()
	defer uiNotifierLock.RUnlock()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(uiNotifiers)
}

func writeScheduleJSONError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(TaskToggleResponse{IsSuccess: false, Error: msg})
}
