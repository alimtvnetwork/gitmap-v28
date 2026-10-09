package cmdui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
)

const maxPortScanAttempts = 50

// RunUI launches the embedded web server and opens the browser.
// The bind host and auto-open behavior come from the parsed `gitmap ui` flags
// (uiListenHost/uiNoBrowser in ui_cmd.go); programmatic callers such as
// `gitmap <module> ui` keep the secure loopback default.
func RunUI(page string, preferredPort int) error {
	listener, port, err := bindAvailablePort(uiListenHost, preferredPort)

	if err != nil {
		return apperror.WrapSimple(err, "ui_bind_port")
	}

	mux := http.NewServeMux()
	mountAPIRoutes(mux)
	mountSPARoutes(mux)

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	targetURL := fmt.Sprintf("http://localhost:%d/%s", port, strings.TrimPrefix(page, "/"))
	fmt.Printf("⚡ GitMap Fleet Web UI running at: %s\n", targetURL)
	if !uiNoBrowser {
		openBrowserURL(targetURL)
	}

	return server.Serve(listener)
}

func bindAvailablePort(host string, startPort int) (net.Listener, int, error) {
	for port := startPort; port < startPort+maxPortScanAttempts; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

		if err == nil {
			return ln, port, nil
		}
	}

	return nil, 0, apperror.NewSimple("bind_port", "no_ports_available")
}

func mountSPARoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(IndexHTML))
	})
}

func mountAPIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ssh/nodes", handleAPISSHNodes)
	mux.HandleFunc("/api/ssh/export", handleAPISSHExport)
	mux.HandleFunc("/api/ssh/import", handleAPISSHImport)
	mux.HandleFunc("/api/ssh/deploy-keys", handleAPISSHDeployKeys)
	mux.HandleFunc("/api/editor/read", handleAPIEditorRead)
	mux.HandleFunc("/api/editor/save", handleAPIEditorSave)
	mux.HandleFunc("/api/commitin/exec", handleAPICommitinExec)
	mux.HandleFunc("/api/terminal/exec", handleAPITerminalExec)
	mux.HandleFunc("/api/settings", handleAPISettings)
	mux.HandleFunc("/api/instances", handleAPIInstances)
	mux.HandleFunc("/api/prompts/instances", handleAPIPromptsInstances)
}

func handleAPIInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     "method not allowed, GET required",
		})
		return
	}

	instances, err := cmdagy.DiscoverAllAgyInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess":      false,
			"totalInstances": 0,
			"instances":      []cmdagy.AgyInstanceInfo{},
			"error":          err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"isSuccess":      true,
		"totalInstances": len(instances),
		"instances":      instances,
	})
}

func handleAPIPromptsInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     "method not allowed, GET required",
		})
		return
	}

	q := r.URL.Query()
	instID := q.Get("instance")
	status := q.Get("status")
	limitStr := q.Get("limit")
	maxWordsStr := q.Get("max_words")
	includeConvsStr := q.Get("include_convs")

	limit := parsePositiveInt(limitStr, 10)
	maxWords := parsePositiveInt(maxWordsStr, 100)

	includeConvs := true
	if includeConvsStr != "" {
		includeConvs = includeConvsStr == "true" || includeConvsStr == "1"
	}

	opts := cmdagy.AgyInstancePromptQueryOptions{
		InstanceID:   instID,
		IsAll:        instID == "" || strings.EqualFold(instID, "all"),
		Status:       status,
		Limit:        limit,
		MaxWords:     maxWords,
		IncludeConvs: includeConvs,
	}

	resp, err := cmdagy.QueryInstancePrompts(opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"isSuccess": false,
			"error":     err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

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
	typedEnv := jsonenvelope.NewEnvelope(
		jsonenvelope.TypeSSHNodes,
		"gitmap-ssh-nodes.json",
		"gitmap ssh export",
		"1.0",
		envelope,
	)
	_ = json.NewEncoder(w).Encode(typedEnv)
}

func parseConnectionsFromRaw(raw []byte) ([]db.SSHConnection, error) {
	payload, _, extractErr := jsonenvelope.ExtractPayload(raw)
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

func handleAPIEditorRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req RemoteFileReadReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		failResult := RemoteFileContent{IsSuccess: false, Error: err.Error()}
		_ = json.NewEncoder(w).Encode(failResult)
		return
	}

	res, readErr := ReadRemoteOrLocalFile(req.NodeAlias, req.FilePath)

	if readErr != nil {
		res.IsSuccess = false
		res.Error = readErr.Error()
	}

	_ = json.NewEncoder(w).Encode(res)
}

func handleAPIEditorSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req RemoteFileSaveReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	saveErr := SaveRemoteOrLocalFile(req.NodeAlias, req.FilePath, req.Content)

	if saveErr != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": saveErr.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func handleAPICommitinExec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var opts CommitinOptions
	_ = json.NewDecoder(r.Body).Decode(&opts)

	gitCmd := exec.Command("git", "commit", "-m", opts.Message)
	out, err := gitCmd.CombinedOutput()

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "output": string(out), "error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "output": string(out)})
}

func handleAPITerminalExec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req TerminalExecReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Command) == "" {
		http.Error(w, "command is required", http.StatusBadRequest)
		return
	}

	timeoutSec := req.TimeoutSec

	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	if timeoutSec > 120 {
		timeoutSec = 120
	}

	isRemote := req.NodeAlias != "" && !strings.EqualFold(req.NodeAlias, "local")

	if isRemote {
		resp := executeRemoteTerminal(req.NodeAlias, req.Command, timeoutSec)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := executeLocalTerminal(req.Command, req.Cwd, timeoutSec)
	_ = json.NewEncoder(w).Encode(resp)
}

func executeRemoteTerminal(nodeAlias, command string, timeoutSec int) TerminalExecResp {
	conn, hasConn := findSSHNodeByAlias(nodeAlias)

	if !hasConn {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     "remote node not found: " + nodeAlias,
		}
	}

	client, isConnected := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s]", nodeAlias))

	if !isConnected {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     "ssh connection failed to " + nodeAlias,
		}
	}
	defer client.Close()

	session, sessErr := client.NewSession()

	if sessErr != nil {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     sessErr.Error(),
		}
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	start := time.Now()
	done := make(chan error, 1)

	go func() {
		done <- session.Run(command)
	}()

	var runErr error

	select {
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		_ = session.Close()
		return TerminalExecResp{
			Success:    false,
			NodeAlias:  nodeAlias,
			ExitCode:   124,
			Error:      fmt.Sprintf("command timed out after %d seconds", timeoutSec),
			DurationMs: time.Since(start).Milliseconds(),
		}
	case runErr = <-done:
	}

	durationMs := time.Since(start).Milliseconds()
	exitCode := resolveSSHExitCode(runErr)
	errMsg := ""

	if runErr != nil {
		errMsg = runErr.Error()
	}

	return TerminalExecResp{
		Success:    runErr == nil && exitCode == 0,
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		ExitCode:   exitCode,
		Error:      errMsg,
		NodeAlias:  nodeAlias,
		DurationMs: durationMs,
	}
}

func executeLocalTerminal(command, cwd string, timeoutSec int) TerminalExecResp {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	if cwd != "" {
		cmd.Dir = cwd
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return TerminalExecResp{
			Success:    false,
			ExitCode:   124,
			Error:      fmt.Sprintf("command timed out after %d seconds", timeoutSec),
			NodeAlias:  "local",
			DurationMs: durationMs,
		}
	}

	exitCode := resolveExitCode(runErr)
	errMsg := ""

	if runErr != nil {
		errMsg = runErr.Error()
	}

	return TerminalExecResp{
		Success:    runErr == nil && exitCode == 0,
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		ExitCode:   exitCode,
		Error:      errMsg,
		NodeAlias:  "local",
		DurationMs: durationMs,
	}
}

func resolveExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func resolveSSHExitCode(err error) int {
	if err == nil {
		return 0
	}

	var sshExitErr *ssh.ExitError

	if errors.As(err, &sshExitErr) {
		return sshExitErr.ExitStatus()
	}

	return 1
}

var customSettingsPath string

func getSettingsFilePath() string {
	if customSettingsPath != "" {
		return customSettingsPath
	}

	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return filepath.Join(home, ".gitmap", "ui_settings.json")
}

func loadSettings() SettingsData {
	def := SettingsData{
		Theme:           "dark",
		DefaultRemote:   "origin",
		ClusterPort:     49152,
		AutoDeployKey:   true,
		GraphicsMode:    "high",
		AutoOpenBrowser: true,
		CommitInLayout:  "split",
		PullDirection:   "pull-left",
		PRReplayMode:    "merges",
		Attributes:      make(map[string]string),
	}
	path := getSettingsFilePath()

	if path == "" {
		return def
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return def
	}

	payload, _, extractErr := jsonenvelope.ExtractPayload(data)

	if extractErr != nil {
		payload = data
	}

	var loaded SettingsData
	err = json.Unmarshal(payload, &loaded)

	if err != nil {
		return def
	}

	if loaded.Attributes == nil {
		loaded.Attributes = make(map[string]string)
	}

	return loaded
}

func saveSettingsData(s SettingsData) error {
	path := getSettingsFilePath()

	if path == "" {
		return nil
	}

	if s.Attributes == nil {
		s.Attributes = make(map[string]string)
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	envelope := jsonenvelope.NewEnvelope(
		jsonenvelope.TypeUISettings,
		path,
		"gitmap ui settings",
		"1.0",
		s,
	)
	data, err := json.MarshalIndent(envelope, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ImportSettingsFromFile loads settings from a file (supporting JSON envelope) and persists them to ui_settings.json.
func ImportSettingsFromFile(filePath string) error {
	raw, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	payload, _, err := jsonenvelope.ExtractPayload(raw)

	if err == nil && len(payload) > 0 {
		raw = payload
	}

	var loaded SettingsData
	err = json.Unmarshal(raw, &loaded)

	if err != nil {
		return err
	}

	return saveSettingsData(loaded)
}

func mergeSettings(current, req SettingsData) SettingsData {
	if req.Theme != "" {
		current.Theme = req.Theme
	}

	if req.DefaultRemote != "" {
		current.DefaultRemote = req.DefaultRemote
	}

	if req.ClusterPort > 0 {
		current.ClusterPort = req.ClusterPort
	}

	current.AutoDeployKey = req.AutoDeployKey

	if req.GraphicsMode != "" {
		current.GraphicsMode = req.GraphicsMode
	}

	current.AutoOpenBrowser = req.AutoOpenBrowser

	if req.CommitInLayout != "" {
		current.CommitInLayout = req.CommitInLayout
	}

	if req.PullDirection != "" {
		current.PullDirection = req.PullDirection
	}

	if req.PRReplayMode != "" {
		current.PRReplayMode = req.PRReplayMode
	}

	if current.Attributes == nil {
		current.Attributes = make(map[string]string)
	}

	if req.Attributes != nil {
		for k, v := range req.Attributes {
			current.Attributes[k] = v
		}
	}

	return current
}

func handleAPISettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		_ = json.NewEncoder(w).Encode(loadSettings())
		return
	}

	var req SettingsData
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	merged := mergeSettings(loadSettings(), req)
	saveErr := saveSettingsData(merged)

	if saveErr != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": saveErr.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func openBrowserURL(targetURL string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("cmd.exe", "/c", "start", "", targetURL).Start()
	case "darwin":
		_ = exec.Command("open", targetURL).Start()
	default:
		_ = exec.Command("xdg-open", targetURL).Start()
	}
}

func parsePositiveInt(str string, fallback int) int {
	if str == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(str)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
