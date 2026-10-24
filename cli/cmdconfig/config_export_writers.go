package cmdconfig

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
)

// filterNonEmptyLines trims each string and retains non-empty lines.
func filterNonEmptyLines(lines []string) []string {
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}

	return results
}

// collectVSCodeExtensions probes installed extensions if the code CLI is available.
func collectVSCodeExtensions() []string {
	cmd := exec.Command(cfgconst.ToolCode, cfgconst.FlagListExtensions)
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	return filterNonEmptyLines(strings.Split(strings.TrimSpace(string(out)), "\n"))
}

// resolveEncodingAndContent produces appropriate encoding token and serialized content.
func resolveEncodingAndContent(data []byte, isBinary bool) (string, string) {
	if isBinary {
		return cfgconst.EncodingBase64, base64.StdEncoding.EncodeToString(data)
	}

	return cfgconst.EncodingUTF8, string(data)
}

// readSingleConfigFile reads a file and packages it into ConfigFilePayload.
func readSingleConfigFile(filePath, fileName string, isBinary bool) (*ConfigFilePayload, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.Wrap(err, "readSingleConfigFile", map[string]any{"path": filePath})
	}

	enc, content := resolveEncodingAndContent(data, isBinary)

	return &ConfigFilePayload{Name: fileName, Encoding: enc, Content: content}, nil
}

// populateDefaultToolFiles assigns fallback template payloads per tool type.
func populateDefaultToolFiles(toolType cfgenum.ConfigToolType, tool string, files map[string]ConfigFilePayload) {
	switch toolType {
	case cfgenum.ConfigToolVSCode:
		files[cfgconst.FileSettingsJSON] = ConfigFilePayload{Name: cfgconst.FileSettingsJSON, Encoding: cfgconst.EncodingUTF8, Content: cfgconst.DefaultVSCodeSettings}
		files[cfgconst.FileKeybindingsJSON] = ConfigFilePayload{Name: cfgconst.FileKeybindingsJSON, Encoding: cfgconst.EncodingUTF8, Content: cfgconst.DefaultVSCodeKeybindings}
	case cfgenum.ConfigToolQTorrent:
		cfgName := resolveDefaultConfigFileNameForOS(tool)
		files[cfgName] = ConfigFilePayload{Name: cfgName, Encoding: cfgconst.EncodingUTF8, Content: cfgconst.DefaultQBittorrentConfig}
	case cfgenum.ConfigToolUTorrent:
		files[cfgconst.FileSettingsDat] = ConfigFilePayload{Name: cfgconst.FileSettingsDat, Encoding: cfgconst.EncodingUTF8, Content: cfgconst.DefaultUTorrentSettings}
	}
}

// collectDefaultToolFiles supplies starter templates when no local files exist.
func collectDefaultToolFiles(tool string) map[string]ConfigFilePayload {
	files := make(map[string]ConfigFilePayload)
	toolType, _ := cfgenum.ParseConfigToolType(tool)
	populateDefaultToolFiles(toolType, tool, files)

	return files
}

// collectToolFiles reads files from the host tool config directory.
func collectToolFiles(tool, srcDir string) (map[string]ConfigFilePayload, []string) {
	files := make(map[string]ConfigFilePayload)
	var extensions []string
	toolType, _ := cfgenum.ParseConfigToolType(tool)
	if toolType.IsVSCode() {
		extensions = collectVSCodeExtensions()
		appendFileIfExists(files, filepath.Join(srcDir, cfgconst.FileSettingsJSON), cfgconst.FileSettingsJSON, false)
		appendFileIfExists(files, filepath.Join(srcDir, cfgconst.FileKeybindingsJSON), cfgconst.FileKeybindingsJSON, false)
	}

	collectTorrentFiles(toolType, srcDir, files)

	return files, extensions
}

// collectTorrentFiles collects qBittorrent and uTorrent configuration files.
func collectTorrentFiles(toolType cfgenum.ConfigToolType, srcDir string, files map[string]ConfigFilePayload) {
	if toolType.IsQTorrent() {
		appendFileIfExists(files, filepath.Join(srcDir, cfgconst.FileQBittorrentIni), cfgconst.FileQBittorrentIni, false)
		appendFileIfExists(files, filepath.Join(srcDir, cfgconst.FileQBittorrentConf), cfgconst.FileQBittorrentConf, false)
	}

	if toolType.IsUTorrent() {
		appendFileIfExists(files, filepath.Join(srcDir, cfgconst.FileSettingsDat), cfgconst.FileSettingsDat, true)
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

// ensureBundleParentDir creates parent folders for exported bundle.
func ensureBundleParentDir(targetFile string) {
	parentDir := filepath.Dir(targetFile)
	if parentDir != cfgconst.DefaultExportFolder && parentDir != "" {
		_ = os.MkdirAll(parentDir, 0o755)
	}
}

// writeBundleJSON marshals the bundle to JSON and writes it to disk.
func writeBundleJSON(bundle ConfigBundle, targetFile string) error {
	ensureBundleParentDir(targetFile)
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
		SchemaVersion: cfgconst.DefaultSchemaVersion,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		SourceOS:      runtime.GOOS,
		Files:         files,
		Extensions:    extensions,
		Metadata: map[string]string{
			"sourceDir": srcDir,
		},
	}
}
