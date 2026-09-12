package cmd

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// collectVSCodeExtensions probes installed extensions if the code CLI is available.
func collectVSCodeExtensions() []string {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}

	return results
}

// readSingleConfigFile reads a file and packages it into ConfigFilePayload.
func readSingleConfigFile(filePath, fileName string, isBinary bool) (*ConfigFilePayload, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.Wrap(err, "readSingleConfigFile", map[string]any{"path": filePath})
	}

	if isBinary {
		return &ConfigFilePayload{
			Name:     fileName,
			Encoding: "base64",
			Content:  base64.StdEncoding.EncodeToString(data),
		}, nil
	}

	return &ConfigFilePayload{
		Name:     fileName,
		Encoding: "utf-8",
		Content:  string(data),
	}, nil
}

// collectDefaultToolFiles supplies starter templates when no local files exist.
func collectDefaultToolFiles(tool string) map[string]ConfigFilePayload {
	files := make(map[string]ConfigFilePayload)
	switch tool {
	case "vscode":
		files["settings.json"] = ConfigFilePayload{Name: "settings.json", Encoding: "utf-8", Content: "{\n  \"editor.tabSize\": 2,\n  \"editor.formatOnSave\": true\n}\n"}
		files["keybindings.json"] = ConfigFilePayload{Name: "keybindings.json", Encoding: "utf-8", Content: "[]\n"}
	case "qtorrent":
		cfgName := resolveDefaultConfigFileNameForOS(tool)
		files[cfgName] = ConfigFilePayload{Name: cfgName, Encoding: "utf-8", Content: "[LegalNotice]\nAccepted=true\n\n[Preferences]\nConnection\\PortRangeMin=6881\nDownloads\\SavePath=Downloads\n"}
	case "utorrent":
		files["settings.dat"] = ConfigFilePayload{Name: "settings.dat", Encoding: "utf-8", Content: "d4:name8:uTorrente"}
	}

	return files
}

// collectToolFiles reads files from the host tool config directory.
func collectToolFiles(tool, srcDir string) (map[string]ConfigFilePayload, []string) {
	files := make(map[string]ConfigFilePayload)
	var extensions []string
	if tool == "vscode" {
		extensions = collectVSCodeExtensions()
		appendFileIfExists(files, filepath.Join(srcDir, "settings.json"), "settings.json", false)
		appendFileIfExists(files, filepath.Join(srcDir, "keybindings.json"), "keybindings.json", false)
	}

	collectTorrentFiles(tool, srcDir, files)

	return files, extensions
}

// collectTorrentFiles collects qBittorrent and uTorrent configuration files.
func collectTorrentFiles(tool, srcDir string, files map[string]ConfigFilePayload) {
	if tool == "qtorrent" {
		appendFileIfExists(files, filepath.Join(srcDir, "qBittorrent.ini"), "qBittorrent.ini", false)
		appendFileIfExists(files, filepath.Join(srcDir, "qBittorrent.conf"), "qBittorrent.conf", false)
	}

	if tool == "utorrent" {
		appendFileIfExists(files, filepath.Join(srcDir, "settings.dat"), "settings.dat", true)
	}
}

// appendFileIfExists checks existence and appends payload if present.
func appendFileIfExists(dest map[string]ConfigFilePayload, filePath, fileName string, isBinary bool) {
	_, statErr := os.Stat(filePath)
	if statErr != nil {
		return
	}

	payload, readErr := readSingleConfigFile(filePath, fileName, isBinary)
	if readErr != nil {
		return
	}

	dest[fileName] = *payload
}

// writeBundleJSON marshals the bundle to JSON and writes it to disk.
func writeBundleJSON(bundle ConfigBundle, targetFile string) error {
	parentDir := filepath.Dir(targetFile)
	if parentDir != "." && parentDir != "" {
		_ = os.MkdirAll(parentDir, 0o755)
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return apperror.Wrap(err, "writeBundleJSON.Marshal", map[string]any{"target": targetFile})
	}

	data = append(data, '\n')
	writeErr := os.WriteFile(targetFile, data, 0o644)
	if writeErr != nil {
		return apperror.Wrap(writeErr, "writeBundleJSON.WriteFile", map[string]any{"target": targetFile})
	}

	return nil
}

// buildConfigBundle builds the full ConfigBundle struct ready for export.
func buildConfigBundle(tool, srcDir string, files map[string]ConfigFilePayload, extensions []string) ConfigBundle {
	return ConfigBundle{
		Tool:          tool,
		SchemaVersion: 1,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		SourceOS:      runtime.GOOS,
		Files:         files,
		Extensions:    extensions,
		Metadata: map[string]string{
			"sourceDir": srcDir,
		},
	}
}
