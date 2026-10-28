package cmdlogviewer

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbrowse"
)

// RunShowLog starts the interactive log viewer server and opens the UI.
func RunShowLog(args []string) error {
	fs := flag.NewFlagSet("show-log", flag.ContinueOnError)
	port := fs.Int("port", 8090, "HTTP listen port")
	fs.IntVar(port, "p", 8090, "HTTP listen port (shorthand)")
	noBrowser := fs.Bool("no-browser", false, "disable auto opening browser")
	branchCompare := fs.Bool("branch-compare", false, "open branch comparison directly")

	if err := fs.Parse(args); err != nil {
		return err
	}

	repoRoot, relPath, err := ResolveTargetRepoAndPath(fs.Args())
	if err != nil {
		cleanTarget := sanitizeTargetArg(fs.Args())
		fmt.Fprintf(os.Stderr, "  ✗ '%s' is not a Git repository.\n", cleanTarget)
		childRepos := findChildRepos(cleanTarget)
		if len(childRepos) > 0 {
			fmt.Fprintln(os.Stderr, "  📁 Detected Git repositories in this directory:")
			for _, cr := range childRepos {
				rel, _ := filepath.Rel(cleanTarget, cr)
				if rel == "" {
					rel = filepath.Base(cr)
				}
				fmt.Fprintf(os.Stderr, "    • %s  (run: gitmap show-log %s)\n", filepath.Base(cr), rel)
			}
		} else {
			fmt.Fprintln(os.Stderr, "  💡 Tip: Run 'gitmap show-log' inside a Git repository or specify target: gitmap show-log <path>")
		}
		return nil
	}

	return startLogServer(repoRoot, relPath, *port, !*noBrowser, *branchCompare)
}

func findChildRepos(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var found []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		candidate := filepath.Join(dir, e.Name())
		gitDir := filepath.Join(candidate, ".git")
		if info, err := os.Stat(gitDir); err == nil && (info.IsDir() || !info.IsDir()) {
			found = append(found, candidate)
			if len(found) >= 10 {
				break
			}
		}
	}

	return found
}

// ResolveTargetRepoAndPath resolves repository root and relative file/dir path.
func ResolveTargetRepoAndPath(args []string) (string, string, error) {
	target := sanitizeTargetArg(args)

	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", "", apperror.WrapSimple(err, "resolve abs path")
	}

	searchDir := resolveProbeDirectory(absPath)
	repoRoot, err := findRepoRoot(searchDir)
	if err != nil {
		return "", "", err
	}

	relPath := computeRelativePath(repoRoot, absPath)

	return repoRoot, relPath, nil
}

func sanitizeTargetArg(args []string) string {
	if len(args) == 0 {
		return "."
	}

	clean := strings.TrimRight(strings.TrimSpace(args[0]), `/\`)
	if clean == "" {
		return "."
	}

	return clean
}

func resolveProbeDirectory(absPath string) string {
	info, err := os.Stat(absPath)
	if err == nil && !info.IsDir() {
		return filepath.Dir(absPath)
	}

	return absPath
}

func findRepoRoot(dir string) (string, error) {
	out, err := runGitInDir(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", newNonGitRepoShowLogError(dir)
	}

	return strings.TrimSpace(out), nil
}

func newNonGitRepoShowLogError(dir string) *apperror.AppError {
	return apperror.NewWithDetails(
		"show-log",
		"E9001",
		"not a git repository: "+dir,
		"cli",
		apperror.ErrorTypeAbort,
		apperror.SeverityWarn,
		map[string]any{"reported": true},
	)
}

func computeRelativePath(repoRoot, absPath string) string {
	rel, err := filepath.Rel(repoRoot, absPath)
	if err != nil || rel == "." || rel == "" {
		return "."
	}

	return filepath.ToSlash(rel)
}

func startLogServer(repoRoot, relPath string, preferredPort int, isAutoOpen, isBranchCompare bool) error {
	ln, port, err := bindFreePort(preferredPort)
	if err != nil {
		return err
	}

	server := &http.Server{ReadHeaderTimeout: 5 * time.Second}
	mux := buildServerMux(repoRoot, relPath, server)
	server.Handler = mux
	targetURL := formatTargetURL(port, relPath, isBranchCompare)

	fmt.Printf("⚡ GitMap Show Log running at: %s\n", targetURL)
	if isAutoOpen {
		openBrowser(targetURL)
	}

	return server.Serve(ln)
}

func formatTargetURL(port int, relPath string, isBranchCompare bool) string {
	base := fmt.Sprintf("http://127.0.0.1:%d/?path=%s", port, url.QueryEscape(relPath))
	if isBranchCompare {
		return base + "&compare=1"
	}

	return base
}

func buildServerMux(repoRoot, initialPath string, server *http.Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndexHTML)
	mux.HandleFunc("/show-log", serveIndexHTML)
	mux.HandleFunc("/favicon.ico", serveFaviconICO)
	mux.HandleFunc("/favicon.svg", serveFaviconSVG)
	mux.HandleFunc("/favicon.png", serveFaviconICO)
	mux.HandleFunc("/manifest.json", serveManifestJSON)
	mux.HandleFunc("/api/log/commits", handleAPICommits(repoRoot))
	mux.HandleFunc("/api/log/commit", handleAPICommitDetail(repoRoot))
	mux.HandleFunc("/api/log/diff", handleAPIDiff(repoRoot))
	mux.HandleFunc("/api/log/compare-commits", handleAPICompareCommits(repoRoot))
	mux.HandleFunc("/api/log/diff-commits", handleAPIDiffCommits(repoRoot))
	mux.HandleFunc("/api/log/branches", handleAPIBranches(repoRoot))
	mux.HandleFunc("/api/log/branch-compare", handleAPIBranchCompare(repoRoot))
	mux.HandleFunc("/api/log/ping", handleAPIPing)
	mux.HandleFunc("/api/log/exit", handleAPIExit(server))

	return mux
}

func handleAPIExit(server *http.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "shutting_down"})
		go func() {
			time.Sleep(100 * time.Millisecond)
			if server != nil {
				_ = server.Close()
			}
			os.Exit(0) // lint-allow: os.Exit
		}()
	}
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(ShowLogHTML))
}

func handleAPIPing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func bindFreePort(startPort int) (net.Listener, int, error) {
	const maxScan = 50
	for p := startPort; p < startPort+maxScan; p++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err == nil {
			return ln, p, nil
		}
	}

	return nil, 0, apperror.NewSimple("bind_free_port", "no free port found")
}

func openBrowser(targetURL string) {
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
