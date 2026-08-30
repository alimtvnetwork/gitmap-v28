package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runUI starts the web dashboard server and opens the UI settings or target page in the default browser.
func runUI(args []string) error {
	checkHelp("ui", args)

	port := parseUIFlags(args)
	targetRoute := resolveUIRoute(args)

	binaryDir := resolveBinaryDir()
	docsDir := filepath.Join(binaryDir, constants.HDDocsDir)

	_, err := os.Stat(docsDir)
	isMissing := os.IsNotExist(err)

	if isMissing {
		_ = ensureDocsSite(binaryDir, docsDir)
	}

	distDir := filepath.Join(docsDir, constants.HDDistDir)
	info, err := os.Stat(distDir)

	if err == nil && info.IsDir() {
		return serveUIStatic(distDir, port, targetRoute)
	}

	return serveUIDev(docsDir, port, targetRoute)
}

func parseUIFlags(args []string) int {
	fset := flag.NewFlagSet("ui", flag.ContinueOnError)
	port := fset.Int("port", 5173, "Port to serve the UI on")

	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			_ = fset.Parse([]string{a})
		}
	}

	return *port
}

func resolveUIRoute(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			switch strings.ToLower(a) {
			case "settings", "setting", "config":
				return "/settings"
			case "pipeline", "pl", "ci":
				return "/pipeline"
			case "terminal", "term", "t":
				return "/terminal"
			case "commands", "cmd":
				return "/commands"
			}
		}
	}

	return "/settings"
}

func serveUIStatic(distDir string, port int, targetRoute string) error {
	url := fmt.Sprintf("http://localhost:%d%s", port, targetRoute)
	fmt.Printf(constants.ColorCyan+"Starting GitMap UI at %s"+constants.ColorReset+"\n", url)

	mux := http.NewServeMux()
	mountTerminalHandlers(mux)
	mux.Handle("/", spaHandler(distDir))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go handleShutdown(server)
	openURL(url)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	fmt.Println("\n" + constants.ColorCyan + "UI server stopped." + constants.ColorReset)
	return nil
}

func serveUIDev(docsDir string, port int, targetRoute string) error {
	url := fmt.Sprintf("http://localhost:%d%s", port, targetRoute)
	fmt.Printf(constants.ColorYellow+"Serving UI in dev mode at %s"+constants.ColorReset+"\n", url)

	mux := http.NewServeMux()
	mountTerminalHandlers(mux)
	mux.Handle("/", spaHandler(docsDir))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go handleShutdown(server)
	openURL(url)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	fmt.Println("\n" + constants.ColorCyan + "UI server stopped." + constants.ColorReset)
	return nil
}
