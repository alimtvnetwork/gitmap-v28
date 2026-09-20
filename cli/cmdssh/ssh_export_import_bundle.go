package cmdssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GitmapExportBundle encapsulates all configuration, macros, connections, and host keys.
type GitmapExportBundle struct {
	Version     string             `json:"version"`
	ExportedAt  string             `json:"exported_at"`
	SourceNode  string             `json:"source_node"`
	Config      map[string]any     `json:"config,omitempty"`
	Macros      []macro.Macro      `json:"macros,omitempty"`
	Connections []db.SSHConnection `json:"connections,omitempty"`
	KnownHosts  string             `json:"known_hosts,omitempty"`
}

func resolveSourceHostname() string {
	name, err := os.Hostname()
	if err == nil && name != "" {
		return name
	}
	return "localhost"
}

func collectLocalConfig() map[string]any {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".gitmap", "config.json")
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil
	}
	var cfg map[string]any
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

func collectLocalKnownHosts() string {
	khPath, err := DefaultKnownHostsPath()
	if err != nil {
		return ""
	}
	data, readErr := os.ReadFile(khPath)
	if readErr != nil {
		return ""
	}
	return string(data)
}

func collectLocalConnections() []db.SSHConnection {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer dbConn.Close()
	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil
	}
	return res.Data
}

func collectExportBundle() (*GitmapExportBundle, *apperror.AppError) {
	macrosRes := macro.ListMacros()
	var macros []macro.Macro
	if macrosRes.IsSuccess() {
		macros = macrosRes.Data
	}
	bundle := &GitmapExportBundle{
		Version:     constants.Version,
		ExportedAt:  time.Now().UTC().Format(time.RFC3339),
		SourceNode:  resolveSourceHostname(),
		Config:      collectLocalConfig(),
		Macros:      macros,
		Connections: collectLocalConnections(),
		KnownHosts:  collectLocalKnownHosts(),
	}
	return bundle, nil
}

func buildRemoteExportWriteCmd(subDir, fileName, b64 string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "$d=[IO.Path]::Combine($env:USERPROFILE, '%s'); if (-not (Test-Path $d)) { [IO.Directory]::CreateDirectory($d) | Out-Null }; $p=[IO.Path]::Combine($d, '%s'); [IO.File]::WriteAllBytes($p, [Convert]::FromBase64String('%s'))"`, subDir, fileName, b64)
	}
	return fmt.Sprintf(`mkdir -p "$HOME/%s" && printf '%%s' '%s' | base64 -d > "$HOME/%s/%s.tmp" && mv -f "$HOME/%s/%s.tmp" "$HOME/%s/%s"`, subDir, b64, subDir, fileName, subDir, fileName, subDir, fileName)
}

func buildRemoteExportAppendCmd(subDir, fileName, b64 string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "$d=[IO.Path]::Combine($env:USERPROFILE, '%s'); if (-not (Test-Path $d)) { [IO.Directory]::CreateDirectory($d) | Out-Null }; $p=[IO.Path]::Combine($d, '%s'); [IO.File]::AppendAllText($p, [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String('%s')))"`, subDir, fileName, b64)
	}
	return fmt.Sprintf(`mkdir -p "$HOME/%s" && printf '%%s' '%s' | base64 -d >> "$HOME/%s/%s"`, subDir, b64, subDir, fileName)
}

func buildRemoteReadBundleCmd(isWin bool) string {
	if isWin {
		return `powershell -NoProfile -Command "$p=[IO.Path]::Combine($env:USERPROFILE, '.gitmap', 'export_bundle.json'); if (Test-Path $p) { [Convert]::ToBase64String([IO.File]::ReadAllBytes($p)) } else { exit 1 }"`
	}
	return `cat "$HOME/.gitmap/export_bundle.json" 2>/dev/null | base64`
}
