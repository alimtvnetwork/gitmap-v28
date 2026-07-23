// Package dashboard provides data collection and output for the HTML dashboard.
package dashboard

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

//go:embed templates/dashboard.html
var templateFS embed.FS

// WriteJSON writes dashboard data to a JSON file.
func WriteJSON(outDir string, data model.DashboardData) (string, error) {
	if err := os.MkdirAll(outDir, constants.DirPermission); err != nil {
		return "", err
	}

	path := filepath.Join(outDir, constants.DashboardJSONFile)

	buf, err := json.MarshalIndent(data, "", constants.JSONIndent)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, buf, constants.FilePermission); err != nil {
		return "", err
	}

	return path, nil
}

// WriteHTML writes the self-contained HTML dashboard with embedded data.
func WriteHTML(outDir string, data model.DashboardData) (string, error) {
	if err := os.MkdirAll(outDir, constants.DirPermission); err != nil {
		return "", err
	}

	tmpl, err := templateFS.ReadFile("templates/dashboard.html")
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Replace the {{.}} placeholder with actual JSON data.
	html := strings.Replace(string(tmpl), "{{.}}", string(jsonBytes), 1)

	path := filepath.Join(outDir, constants.DashboardHTMLFile)

	if err := os.WriteFile(path, []byte(html), constants.FilePermission); err != nil {
		return "", err
	}

	return path, nil
}

// formatSize returns a human-readable file size string.
func formatSize(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	size := info.Size()

	const kb = 1024

	if size < kb {
		return fmt.Sprintf("%d B", size)
	}

	return fmt.Sprintf("%.1f KB", float64(size)/float64(kb))
}

// Summary returns a formatted summary line for a written file.
func Summary(path string) string {
	var buf bytes.Buffer
	buf.WriteString(path)

	if s := formatSize(path); s != "" {
		buf.WriteString(" (")
		buf.WriteString(s)
		buf.WriteString(")")
	}

	return buf.String()
}
