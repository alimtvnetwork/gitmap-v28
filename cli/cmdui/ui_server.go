package cmdui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbrowse"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdexclude"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const maxPortScanAttempts = 50

var activeUIServer *http.Server

// RunUI launches the embedded web server and opens the browser.
// The bind host and auto-open behavior come from the parsed `gitmap ui` flags
// (uiListenHost/uiNoBrowser in ui_cmd.go); programmatic callers such as
// `gitmap <module> ui` keep the secure loopback default.
func RunUI(page string, preferredPort int) error {
	listener, port, err := bindAvailablePort(uiListenHost, preferredPort)
	if err != nil {
		return apperror.WrapSimple(err, "ui_bind_port")
	}

	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
	}
	activeUIServer = server

	mux := http.NewServeMux()
	mountAPIRoutes(mux)
	mountSPARoutes(mux)
	server.Handler = mux

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
	mountSSHAndEditorRoutes(mux)
	mountToolAndScanRoutes(mux)
	mountContextMenuRoutes(mux)
}

func mountSSHAndEditorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ssh/nodes", handleAPISSHNodes)
	mux.HandleFunc("/api/ssh/export", handleAPISSHExport)
	mux.HandleFunc("/api/ssh/import", handleAPISSHImport)
	mux.HandleFunc("/api/ssh/deploy-keys", handleAPISSHDeployKeys)
	mux.HandleFunc("/api/editor/read", handleAPIEditorRead)
	mux.HandleFunc("/api/editor/save", handleAPIEditorSave)
	mux.HandleFunc("/api/commitin/exec", handleAPICommitinExec)
	mux.HandleFunc("/api/terminal/exec", handleAPITerminalExec)
}

func mountToolAndScanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/search/repo", handleAPISearchRepo)
	mux.HandleFunc("/api/project/open-antigravity", handleAPIOpenAntigravity)
	mux.HandleFunc("/api/command/exec", handleAPITerminalExec)
	mux.HandleFunc("/api/help/topics", handleAPIHelpTopics)
	mux.HandleFunc("/api/help/doc", handleAPIHelpDoc)
	mux.HandleFunc("/api/settings", handleAPISettings)
	mux.HandleFunc("/api/instances", handleAPIInstances)
	mux.HandleFunc("/api/prompts/instances", handleAPIPromptsInstances)
	mux.HandleFunc("/api/scan/repos", handleAPIScanRepos)
	mux.HandleFunc("/api/scan/disable", handleAPIScanDisable)
	mux.HandleFunc("/api/scan/enable", handleAPIScanEnable)
	mux.HandleFunc("/api/ui/exit", handleAPIUIExit(activeUIServer))
}

func handleAPIUIExit(server *http.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "shutting_down"})

		go func() {
			time.Sleep(100 * time.Millisecond)
			target := server
			if target == nil {
				target = activeUIServer
			}
			if target != nil {
				_ = target.Close()
			}
			os.Exit(0)
		}()
	}
}

func handleAPIEditorRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
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
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
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
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
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
	if cmdbrowse.LaunchBrowserAppWindow(targetURL) {
		return
	}

	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("cmd.exe", "/c", "start", "", targetURL).Start()
	case "darwin":
		_ = exec.Command("open", targetURL).Start()
	default:
		_ = exec.Command("xdg-open", targetURL).Start()
	}
}

// RepoScanItem represents a repository's scan status for the Web UI.
type RepoScanItem struct {
	Id           int    `json:"id"`
	RepoName     string `json:"repoName"`
	RelativePath string `json:"relativePath"`
	Branch       string `json:"branch"`
	IsExcluded   bool   `json:"isExcluded"`
	LastScanAt   string `json:"lastScanAt"`
}

// ScanReposResponse is the envelope returned by GET /api/scan/repos.
type ScanReposResponse struct {
	IsSuccess  bool           `json:"isSuccess"`
	TotalCount int            `json:"totalCount"`
	Repos      []RepoScanItem `json:"repos"`
}

// ToggleScanReposRequest defines the payload for disable/enable scan endpoints.
type ToggleScanReposRequest struct {
	Repos []string `json:"repos"`
}

// ToggleScanReposResponse is the envelope returned by POST /api/scan/disable and /api/scan/enable.
type ToggleScanReposResponse struct {
	IsSuccess     bool     `json:"isSuccess"`
	ModifiedCount int      `json:"modifiedCount"`
	AffectedRepos []string `json:"affectedRepos"`
	CurrentTotal  int      `json:"currentTotal"`
}

func handleAPIScanRepos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}

	repos, err := fetchScanRepos()
	if err != nil {
		_ = json.NewEncoder(w).Encode(ScanReposResponse{IsSuccess: false})
		return
	}

	resp := ScanReposResponse{IsSuccess: true, TotalCount: len(repos), Repos: repos}
	_ = json.NewEncoder(w).Encode(resp)
}

func fetchScanRepos() ([]RepoScanItem, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open_default_store")
	}

	defer db.Close()

	records, listErr := db.ListRepos()
	if listErr != nil {
		return nil, apperror.WrapSimple(listErr, "list_repos")
	}

	return mapRecordsToScanItems(records), nil
}

func mapRecordsToScanItems(records []model.ScanRecord) []RepoScanItem {
	items := make([]RepoScanItem, 0, len(records))
	now := time.Now().UTC().Format(time.RFC3339)

	for _, rec := range records {
		items = append(items, buildRepoScanItem(rec, now))
	}

	return items
}

func buildRepoScanItem(rec model.ScanRecord, timestamp string) RepoScanItem {
	isExcl := cmdexclude.IsCustomExcluded(rec.RelativePath) || cmdexclude.IsCustomExcluded(rec.Slug)

	return RepoScanItem{
		Id:           int(rec.ID),
		RepoName:     rec.RepoName,
		RelativePath: rec.RelativePath,
		Branch:       rec.Branch,
		IsExcluded:   isExcl,
		LastScanAt:   timestamp,
	}
}

func handleAPIScanDisable(w http.ResponseWriter, r *http.Request) {
	processScanToggleAPI(w, r, true)
}

func handleAPIScanEnable(w http.ResponseWriter, r *http.Request) {
	processScanToggleAPI(w, r, false)
}

func processScanToggleAPI(w http.ResponseWriter, r *http.Request, isDisable bool) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req ToggleScanReposRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(ToggleScanReposResponse{IsSuccess: false})
		return
	}

	resp := executeScanToggleAPI(req.Repos, isDisable)
	_ = json.NewEncoder(w).Encode(resp)
}

func executeScanToggleAPI(repos []string, isDisable bool) ToggleScanReposResponse {
	selection := strings.Join(repos, ",")
	res, err := cmdexclude.ToggleRepoExclusions(selection, isDisable)
	if err != nil {
		return ToggleScanReposResponse{IsSuccess: false}
	}

	return ToggleScanReposResponse{
		IsSuccess:     res.IsSuccess,
		ModifiedCount: res.ModifiedCount,
		AffectedRepos: res.AffectedRepos,
		CurrentTotal:  res.TotalExcluded,
	}
}
