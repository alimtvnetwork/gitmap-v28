package cmdconfig

import (
	"os"
	"path/filepath"
	"strings"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
)

// ConfigFilePayload holds a single configuration file payload.
type ConfigFilePayload struct {
	Name     string `json:"name"`
	Encoding string `json:"encoding"` // "utf-8" or "base64"
	Content  string `json:"content"`
}

// ConfigBundle represents the universal JSON export format for tool configurations.
type ConfigBundle struct {
	Tool          string                       `json:"tool"`
	SchemaVersion int                          `json:"schemaVersion"`
	ExportedAt    string                       `json:"exportedAt"`
	SourceOS      string                       `json:"sourceOS"`
	Files         map[string]ConfigFilePayload `json:"files"`
	Extensions    []string                     `json:"extensions,omitempty"`
	Metadata      map[string]string            `json:"metadata,omitempty"`
}

// normalizeConfigTool maps input tool identifiers to canonical config names.
func normalizeConfigTool(name string) string {
	toolType, err := cfgenum.ParseConfigToolType(name)
	if err == nil {
		return toolType.String()
	}

	return strings.ToLower(strings.TrimSpace(name))
}

// resolveDefaultConfigFileName returns the default file name matching user input.
func resolveDefaultConfigFileName(inputTool string, canonicalTool string) string {
	lower := strings.ToLower(strings.TrimSpace(inputTool))
	if lower == cfgconst.ToolUTTorrent {
		return cfgconst.ToolUTTorrent + cfgconst.ExtJSON
	}

	if lower == cfgconst.ToolQTorrent {
		return cfgconst.ToolQTorrent + cfgconst.ExtJSON
	}

	return canonicalTool + cfgconst.ExtJSON
}

// isConfigDirectoryPath checks if path represents a directory.
func isConfigDirectoryPath(path string) bool {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return true
	}

	fi, err := os.Stat(path)

	return err == nil && fi.IsDir()
}

// resolveConfigFilePath determines the final output/input file path.
func resolveConfigFilePath(specifiedPath, defaultFilename string) string {
	cleanPath := strings.TrimSpace(specifiedPath)
	if cleanPath == "" {
		return defaultFilename
	}

	if isConfigDirectoryPath(cleanPath) {
		return filepath.Join(cleanPath, defaultFilename)
	}

	return cleanPath
}
