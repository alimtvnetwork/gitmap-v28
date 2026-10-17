package cmdui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
	mux.HandleFunc("/api/command/exec", handleAPITerminalExec)
	mux.HandleFunc("/api/help/topics", handleAPIHelpTopics)
	mux.HandleFunc("/api/help/doc", handleAPIHelpDoc)
	mux.HandleFunc("/api/settings", handleAPISettings)
	mux.HandleFunc("/api/instances", handleAPIInstances)
	mux.HandleFunc("/api/prompts/instances", handleAPIPromptsInstances)
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
