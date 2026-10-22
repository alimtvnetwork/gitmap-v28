package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ActionLogFilter specifies query filters for action logs.
type ActionLogFilter struct {
	TaskId     string
	AgentRole  string
	ActionType string
	TouchType  string
	SubtaskId  string
	Search     string
	SortBy     string
	Order      string
	Limit      int
	Offset     int
}

// StartAgentUIServer binds HTTP server and serves the visualizer dashboard.
func StartAgentUIServer(opts UIOptions) *appfault.AppError {
	ln, boundPort, scanErr := scanAvailableListener(opts.Port)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return scanErr
	}
	defer ln.Close()

	srv := buildHTTPServer(opts, boundPort)
	printServerBanner(boundPort, opts.Dir)
	triggerBrowserIfRequested(opts.IsBrowse, boundPort)
	setupGracefulShutdown(srv)

	return handleServerExit(srv.Serve(ln))
}

func handleServerExit(err error) *appfault.AppError {
	isClosed := err == nil || err == http.ErrServerClosed
	if isClosed {
		return nil
	}

	return appfault.WrapExecution(err, "agent visualizer server encountered an error")
}

func buildHTTPServer(opts UIOptions, boundPort int) *http.Server {
	return &http.Server{
		Handler:           buildAgentUIMux(opts, boundPort),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func scanAvailableListener(basePort int) (net.Listener, int, *appfault.AppError) {
	const maxScan = 50
	for i := 0; i < maxScan; i++ {
		candidatePort := basePort + i
		ln, isBound := tryBindPort(candidatePort)
		if isBound {
			return ln, candidatePort, nil
		}
	}

	return tryBindFallbackPort()
}

func tryBindPort(port int) (net.Listener, bool) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	hasLn := err == nil && ln != nil

	return ln, hasLn
}

func tryBindFallbackPort() (net.Listener, int, *appfault.AppError) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	hasLn := err == nil && ln != nil
	if hasLn {
		boundPort := ln.Addr().(*net.TCPAddr).Port
		return ln, boundPort, nil
	}

	return nil, 0, appfault.WrapExecution(err, "failed to bind agent UI on any port")
}

func triggerBrowserIfRequested(isBrowse bool, port int) {
	if isBrowse {
		go func() {
			time.Sleep(150 * time.Millisecond)
			launchBrowser(port)
		}()
	}
}

func launchBrowser(port int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	cmd := resolveBrowserCmd(url)
	hasCmd := cmd != nil
	if hasCmd {
		_ = cmd.Start()
	}
}

func resolveBrowserCmd(url string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd.exe", "/c", "start", url)
	case "darwin":
		return exec.Command("open", url)
	default:
		return exec.Command("xdg-open", url)
	}
}

func printServerBanner(port int, tempDir string) {
	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s║ GITMAP AI AGENT VISUALIZER DASHBOARD                                         ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Dashboard URL:  %shttp://127.0.0.1:%d%s\n", constants.ColorGreen, port, constants.ColorReset)
	fmt.Printf("  • Temp Directory: %s%s%s\n", constants.ColorYellow, tempDir, constants.ColorReset)
	fmt.Printf("  • Press %sCtrl+C%s to stop server\n\n", constants.ColorCyan, constants.ColorReset)
}

func setupGracefulShutdown(srv *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\nShutting down GitMap Agent Visualizer...\n")
		_ = srv.Close()
	}()
}

func buildAgentUIMux(opts UIOptions, boundPort int) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ServeDashboardHTML)
	mux.HandleFunc("/api/agent/status", makeStatusHandler(opts, boundPort))
	mux.HandleFunc("/api/agent/tasks", makeTasksHandler(opts))
	mux.HandleFunc("/api/agent/logs", makeLogsHandler(opts))
	mux.HandleFunc("/api/agent/crashes", makeCrashesHandler(opts))
	mux.HandleFunc("/api/agent/clear", makeClearHandler(opts))

	return mux
}

func makeStatusHandler(opts UIOptions, boundPort int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := buildStatusPayload(opts, boundPort)
		writeJSON(w, http.StatusOK, payload)
	}
}

func makeTasksHandler(opts UIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := resolveRequestTaskId(r, opts.TaskId)
		payload := buildTasksPayload(opts.Dir, taskId)
		writeJSON(w, http.StatusOK, payload)
	}
}

func makeLogsHandler(opts UIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := parseLogFilter(r, opts.TaskId)
		logs := queryFilteredActionLogs(opts.Dir, filter)
		writeJSON(w, http.StatusOK, map[string]any{"logs": logs, "count": len(logs)})
	}
}

func makeCrashesHandler(opts UIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := resolveRequestTaskId(r, opts.TaskId)
		report := buildTaskCrashReport(opts.Dir, taskId)
		writeJSON(w, http.StatusOK, report)
	}
}

func makeClearHandler(opts UIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := resolveClearTaskId(r, opts.TaskId)
		resp := executeServerClear(opts.Dir, taskId)
		writeJSON(w, http.StatusOK, resp)
	}
}

func resolveRequestTaskId(r *http.Request, defaultId string) string {
	qTask := strings.TrimSpace(r.URL.Query().Get("taskId"))
	hasQTask := len(qTask) > 0
	if hasQTask {
		return qTask
	}

	return defaultId
}

func resolveClearTaskId(r *http.Request, defaultId string) string {
	qTask := strings.TrimSpace(r.URL.Query().Get("taskId"))
	hasQTask := len(qTask) > 0
	if hasQTask {
		return qTask
	}

	bodyTask := parseBodyTaskId(r.Body)
	hasBodyTask := len(bodyTask) > 0
	if hasBodyTask {
		return bodyTask
	}

	return defaultId
}

func parseBodyTaskId(body io.ReadCloser) string {
	hasBody := body != nil
	if !hasBody {
		return ""
	}
	defer body.Close()

	var req map[string]string
	decodeErr := json.NewDecoder(body).Decode(&req)
	hasDecode := decodeErr == nil
	if hasDecode {
		return strings.TrimSpace(req["taskId"])
	}

	return ""
}

func executeServerClear(tempDir, taskId string) map[string]any {
	hasTaskId := len(strings.TrimSpace(taskId)) > 0
	if hasTaskId {
		return clearSpecificTask(tempDir, taskId)
	}

	return clearAllCompletedTasks(tempDir)
}

func clearSpecificTask(tempDir, taskId string) map[string]any {
	taskDir, findErr := ResolveTaskDir(tempDir, taskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return map[string]any{"status": "ERROR", "message": findErr.Message}
	}

	removeErr := os.RemoveAll(taskDir)
	hasRemoveErr := removeErr != nil
	if hasRemoveErr {
		return map[string]any{"status": "ERROR", "message": removeErr.Error()}
	}
	removeRegistryEntries(tempDir, taskId)

	return map[string]any{"status": "CLEARED", "taskId": taskId, "dir": filepath.ToSlash(taskDir)}
}

func clearAllCompletedTasks(tempDir string) map[string]any {
	runDirs := scanRunDirs(tempDir)
	var cleared []string
	for _, dir := range runDirs {
		isDone := isTaskRunCompleted(dir)
		if isDone {
			_ = os.RemoveAll(dir)
			removeRegistryEntries(tempDir, filepath.Base(dir))
			cleared = append(cleared, filepath.ToSlash(dir))
		}
	}

	return map[string]any{"status": "CLEARED", "clearedCount": len(cleared), "dirs": cleared}
}

func buildTaskCrashReport(tempDir, taskId string) map[string]any {
	taskDir, dirErr := ResolveTaskDir(tempDir, taskId)
	hasDirErr := dirErr != nil
	if hasDirErr {
		return emptyCrashReport()
	}

	report, diagErr := buildDiagnosticReport(taskDir)
	hasDiagErr := diagErr != nil
	if hasDiagErr {
		return emptyCrashReport()
	}

	return formatCrashReport(report, taskDir)
}

func emptyCrashReport() map[string]any {
	return map[string]any{"hasCrashesDetected": false, "crashedAgents": []any{}, "counts": map[string]int{}}
}

func formatCrashReport(report *DiagnosticReport, taskDir string) map[string]any {
	return map[string]any{
		"hasCrashesDetected": report.HasCrashesDetected,
		"crashedAgents":      report.CrashedAgents,
		"counts":             report.Counts,
		"taskDirectory":      filepath.ToSlash(taskDir),
	}
}

func buildStatusPayload(opts UIOptions, boundPort int) map[string]any {
	tempDir := ResolveAgentTempDir(opts.Dir)
	activeTask, taskDir := resolveActiveTaskAndDir(tempDir, opts.TaskId)
	metrics := queryGlobalMetricsOrCompute(tempDir)
	crashes := checkTaskCrashes(taskDir)
	counts := queryTaskSubtaskCounts(taskDir)

	return formatStatusPayload(activeTask, tempDir, metrics, crashes, counts, boundPort)
}

func formatStatusPayload(task map[string]any, dir string, m map[string]string, cr []CrashAutopsy, cnt SubtaskCounts, port int) map[string]any {
	status := computeOverallStatus(task, cr, cnt)
	hasCrashes := len(cr) > 0

	return map[string]any{
		"status":               status,
		"tempDirectory":        filepath.ToSlash(dir),
		"activeParentTask":     task,
		"globalMetrics":        m,
		"crashesDetectedCount": len(cr),
		"hasCrashesDetected":   hasCrashes,
		"serverPort":           port,
	}
}

func computeOverallStatus(task map[string]any, crashes []CrashAutopsy, counts SubtaskCounts) string {
	hasCrashes := len(crashes) > 0
	if hasCrashes {
		return "CRASH_DETECTED"
	}
	hasRunning := counts.InProgress > 0
	if hasRunning {
		return "RUNNING"
	}
	hasTask := task != nil && len(task) > 0
	if hasTask {
		return resolveTaskStatusString(task, counts)
	}

	return "IDLE"
}

func resolveTaskStatusString(task map[string]any, counts SubtaskCounts) string {
	rawStatus, hasStatus := task["status"].(string)
	if hasStatus && (rawStatus == "DONE" || rawStatus == "COMPLETED") {
		return "COMPLETED"
	}
	hasRemaining := counts.Pending > 0 || counts.InProgress > 0
	if hasRemaining {
		return "ACTIVE"
	}

	return "IDLE"
}

func checkTaskCrashes(taskDir string) []CrashAutopsy {
	hasDir := len(taskDir) > 0
	if !hasDir {
		return nil
	}
	report, err := buildDiagnosticReport(taskDir)
	hasErr := err != nil || report == nil
	if hasErr {
		return nil
	}

	return report.CrashedAgents
}

func queryTaskSubtaskCounts(taskDir string) SubtaskCounts {
	hasDir := len(taskDir) > 0
	if !hasDir {
		return SubtaskCounts{}
	}
	db, err := OpenTier2DB(taskDir)
	hasErr := err != nil
	if hasErr {
		return SubtaskCounts{}
	}
	defer db.Close()
	counts, _ := querySubtaskCounts(db)

	return counts
}

func resolveActiveTaskAndDir(tempDir, taskFlag string) (map[string]any, string) {
	taskDir, err := ResolveTaskDir(tempDir, taskFlag)
	hasDir := err == nil && len(taskDir) > 0
	if !hasDir {
		return map[string]any{}, ""
	}
	task := readParentTaskFromDir(taskDir)

	return task, taskDir
}
