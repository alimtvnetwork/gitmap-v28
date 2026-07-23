package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// maxDocsSiteSize is the maximum total extraction size for docs-site.zip (100 MB).
const maxDocsSiteSize = 100 * 1024 * 1024

// runHelpDashboard serves the docs site locally.
func runHelpDashboard(args []string) {
	checkHelp("help-dashboard", args)

	port := parseHelpDashboardFlags(args)
	binaryDir := resolveBinaryDir()
	docsDir := filepath.Join(binaryDir, constants.HDDocsDir)

	// Auto-extract docs-site.zip if docs-site/ directory doesn't exist.
	// If the zip is also missing (older installer, or `gitmap update` not yet
	// run after a docs-site release), try to download it from GitHub first.
	// If that also fails (release didn't bundle docs-site.zip), gracefully
	// fall back to opening the hosted docs URL instead of hard-exiting.
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		zipPath := filepath.Join(binaryDir, constants.DocsSiteArchive)
		if _, zipErr := os.Stat(zipPath); os.IsNotExist(zipErr) {
			_, n, dlErr := downloadDocsSiteArchive(zipPath)
			if dlErr != nil {
				fmt.Fprintf(os.Stderr, constants.ErrDocsSiteDownload, 2, dlErr, zipPath)
				openHostedDocsFallback()
				return
			}
			fmt.Printf(constants.MsgDocsSiteDownloaded, n)
		}
		fmt.Printf("  Extracting %s...\n", constants.DocsSiteArchive)
		if mkErr := os.MkdirAll(docsDir, constants.DirPermission); mkErr != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to create docs-site dir: %v\n", mkErr)
			openHostedDocsFallback()
			return
		}
		extractTarget := chooseDocsExtractTarget(zipPath, binaryDir, docsDir)
		if extractErr := extractDocsSiteZip(zipPath, extractTarget); extractErr != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to extract docs-site.zip: %v\n", extractErr)
			openHostedDocsFallback()
			return
		}
		fmt.Printf("  ✓ Docs site extracted to %s\n", docsDir)
	}

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.ErrHDNoDocsDir, docsDir)
		openHostedDocsFallback()
		return
	}

	distDir := filepath.Join(docsDir, constants.HDDistDir)

	if info, err := os.Stat(distDir); err == nil && info.IsDir() {
		serveStatic(distDir, port)
	} else {
		fmt.Print(constants.MsgHDNoDistFallback)
		serveDev(docsDir, port)
	}
}

// parseHelpDashboardFlags parses the --port flag.
func parseHelpDashboardFlags(args []string) int {
	fs := flag.NewFlagSet(constants.CmdHelpDashboard, flag.ExitOnError)
	port := fs.Int("port", constants.HDDefaultPort, constants.FlagDescHDPort)
	fs.Parse(args)

	return *port
}

// resolveBinaryDir returns the directory containing the gitmap binary.
func resolveBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Dir(exe)
	}

	return filepath.Dir(resolved)
}

// serveStatic serves pre-built dist/ files over HTTP with SPA fallback.
func serveStatic(distDir string, port int) {
	fmt.Printf(constants.MsgHDServingStatic, distDir, port)
	openBrowser(port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           spaHandler(distDir),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go handleShutdown(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, constants.ErrHDServe, err)
		os.Exit(1)
	}

	fmt.Print(constants.MsgHDStopped)
}

// spaHandler serves static files from distDir and falls back to
// index.html for unknown routes so client-side routers (React Router,
// TanStack Router) resolve deep links. Also forces text/html on the
// fallback because Windows' MIME registry can be broken and cause
// browsers to download index.html instead of rendering it.
func spaHandler(distDir string) http.Handler {
	fs := http.FileServer(http.Dir(distDir))
	indexPath := distDir + "/index.html"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the exact file if it exists on disk.
		requested := distDir + r.URL.Path
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			// Explicit HTML content type for .html assets (Windows fix).
			if strings.HasSuffix(r.URL.Path, ".html") {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
			}
			fs.ServeHTTP(w, r)
			return
		}
		// Fallback to SPA index.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, indexPath)
	})
}

// serveDev runs npm install + npm run dev as a fallback.
func serveDev(docsDir string, port int) {
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrHDNPMNotFound)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgHDRunningNPM)

	install := exec.Command(npmPath, "install")
	install.Dir = docsDir
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr

	if err := install.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHDNPMInstall, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgHDStartingDev, docsDir)

	dev := exec.Command(npmPath, "run", "dev", "--", "--port", fmt.Sprintf("%d", port))
	dev.Dir = docsDir
	dev.Stdout = os.Stdout
	dev.Stderr = os.Stderr

	if err := dev.Start(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHDDevServer, err)
		os.Exit(1)
	}

	openBrowser(port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	_ = dev.Process.Kill()
	fmt.Print(constants.MsgHDStopped)
}

// openBrowser opens the local dev/static URL in the default browser.
func openBrowser(port int) {
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf(constants.MsgHDOpening, port)
	openURL(url)
}

// handleShutdown gracefully stops the static server on Ctrl+C.
func handleShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	server.Close()
}
