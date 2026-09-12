package cmd

import (
	"os"
	"path/filepath"
	"strings"
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
	lower := strings.ToLower(strings.TrimSpace(name))
	switch lower {
	case "vscode", "code", "vcode":

		return "vscode"
	case "qtorrent", "qbittorrent", "qbit":

		return "qtorrent"
	case "utorrent", "uttorrent", "u-torrent":

		return "utorrent"
	case "all", "*":

		return "all"
	default:

		return lower
	}
}

// resolveDefaultConfigFileName returns the default file name matching user input.
func resolveDefaultConfigFileName(inputTool string, canonicalTool string) string {
	lower := strings.ToLower(strings.TrimSpace(inputTool))
	if lower == "uttorrent" {
		return "uttorrent.json"
	}

	if lower == "qtorrent" {
		return "qtorrent.json"
	}

	return canonicalTool + ".json"
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
