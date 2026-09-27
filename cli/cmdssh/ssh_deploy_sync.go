// Package cmdssh — ssh_deploy_sync.go handles remote file probing and conflict resolution.
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

// DeploySyncMode defines bidirectional synchronization conflict behavior.
type DeploySyncMode string

const (
	// SyncModeNone prompts interactively on conflict.
	SyncModeNone DeploySyncMode = ""
	// SyncModeOverwrite forces overwriting the remote file unconditionally.
	SyncModeOverwrite DeploySyncMode = "overwrite"
	// SyncModeSkip skips transfer if remote file exists.
	SyncModeSkip DeploySyncMode = "skip"
	// SyncModeSync synchronizes bidirectionally based on latest mtime.
	SyncModeSync DeploySyncMode = "sync"
	// SyncModeSyncRight pushes local to remote only if local is newer.
	SyncModeSyncRight DeploySyncMode = "sync-right"
	// SyncModeSyncLeft pulls remote to local only if remote is newer.
	SyncModeSyncLeft DeploySyncMode = "sync-left"
)

// DeployAction represents the decided action for a specific file transfer.
type DeployAction string

const (
	// ActionTransferToRemote copies local file to remote target.
	ActionTransferToRemote DeployAction = "transfer-to-remote"
	// ActionTransferToLocal copies remote file to local workspace.
	ActionTransferToLocal DeployAction = "transfer-to-local"
	// ActionSkip skips file without transfer.
	ActionSkip DeployAction = "skip"
	// ActionConflictPrompt requests user decision interactively.
	ActionConflictPrompt DeployAction = "conflict-prompt"
)

// RemoteFileInfo contains metadata about a remote file or directory.
type RemoteFileInfo struct {
	Exists  bool      `json:"exists"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Hash    string    `json:"hash,omitempty"`
}

// LocalFileInfo contains metadata about a local file or directory.
type LocalFileInfo struct {
	RelPath string    `json:"relPath"`
	AbsPath string    `json:"absPath"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Hash    string    `json:"hash,omitempty"`
}

type rawRemoteProbe struct {
	Exists   bool   `json:"exists"`
	IsDir    bool   `json:"isDir"`
	Size     int64  `json:"size"`
	ModTime  string `json:"modTime"`
	ModEpoch int64  `json:"modEpoch"`
}

// probeRemoteFileInfo probes existence, size, directory status, and modification time of a remote path.
func probeRemoteFileInfo(client *ssh.Client, remotePath string, isWin bool) (RemoteFileInfo, error) {
	if client == nil {
		return RemoteFileInfo{}, apperror.NewInvalidInput("client", "E_SSH_NIL", "ssh client cannot be nil")
	}
	cmd := buildRemoteProbeCmd(remotePath, isWin)
	rawOut, err := crypto.RunCommand(client, cmd, "")
	if err != nil {
		return probeFallbackOnErr(rawOut, err)
	}
	return parseRemoteProbeOutput(rawOut)
}

func buildRemoteProbeCmd(remotePath string, isWin bool) string {
	if isWin {
		return buildWindowsProbeCmd(remotePath)
	}
	return buildUnixProbeCmd(remotePath)
}

func buildWindowsProbeCmd(path string) string {
	escaped := strings.ReplaceAll(path, "'", "''")
	return fmt.Sprintf("powershell.exe -NoProfile -Command \"$p='%s'; if (Test-Path -LiteralPath $p) { $it=Get-Item -LiteralPath $p; $sz=0; if (-not $it.PSIsContainer) { $sz=$it.Length }; [PSCustomObject]@{exists=$true; isDir=$it.PSIsContainer; size=$sz; modTime=$it.LastWriteTimeUtc.ToString('o')} | ConvertTo-Json -Compress } else { [PSCustomObject]@{exists=$false; isDir=$false; size=0; modTime=''} | ConvertTo-Json -Compress }\"", escaped)
}

func buildUnixProbeCmd(path string) string {
	escaped := strings.ReplaceAll(path, "'", "'\\''")
	return fmt.Sprintf("sh -c 'p=\"%s\"; if [ -e \"$p\" ]; then d=false; [ -d \"$p\" ] && d=true; sz=$(stat -c \"%%s\" \"$p\" 2>/dev/null || stat -f \"%%z\" \"$p\" 2>/dev/null || echo 0); mt=$(stat -c \"%%Y\" \"$p\" 2>/dev/null || stat -f \"%%m\" \"$p\" 2>/dev/null || echo 0); echo \"{\\\"exists\\\":true,\\\"isDir\\\":$d,\\\"size\\\":$sz,\\\"modEpoch\\\":$mt}\"; else echo \"{\\\"exists\\\":false,\\\"isDir\\\":false,\\\"size\\\":0,\\\"modEpoch\\\":0}\"; fi'", escaped)
}

func probeFallbackOnErr(rawOut string, err error) (RemoteFileInfo, error) {
	trimmed := strings.TrimSpace(rawOut)
	if strings.Contains(trimmed, `"exists":false`) || strings.Contains(trimmed, `"exists": false`) {
		return RemoteFileInfo{Exists: false}, nil
	}
	return RemoteFileInfo{}, apperror.WrapSimple(err, "probeRemoteFileInfo")
}

func parseRemoteProbeOutput(rawOut string) (RemoteFileInfo, error) {
	trimmed := strings.TrimSpace(rawOut)
	idxStart := strings.Index(trimmed, "{")
	idxEnd := strings.LastIndex(trimmed, "}")
	if idxStart < 0 || idxEnd <= idxStart {
		return RemoteFileInfo{}, apperror.NewExecutionError(fmt.Sprintf("invalid probe output: %s", rawOut))
	}
	jsonStr := trimmed[idxStart : idxEnd+1]
	return parseProbeJSON([]byte(jsonStr))
}

func parseProbeJSON(data []byte) (RemoteFileInfo, error) {
	var raw rawRemoteProbe
	if err := json.Unmarshal(data, &raw); err != nil {
		return RemoteFileInfo{}, apperror.WrapSimple(err, "parseProbeJSON")
	}
	return RemoteFileInfo{
		Exists:  raw.Exists,
		IsDir:   raw.IsDir,
		Size:    raw.Size,
		ModTime: parseProbeModTime(raw.ModTime, raw.ModEpoch),
	}, nil
}

func parseProbeModTime(modStr string, epoch int64) time.Time {
	if epoch > 0 {
		return time.Unix(epoch, 0)
	}
	if modStr == "" {
		return time.Time{}
	}
	if parsed, err := time.Parse(time.RFC3339Nano, modStr); err == nil {
		return parsed
	}
	parsed, _ := time.Parse(time.RFC3339, modStr)
	return parsed
}

// evaluateFileConflict determines the sync action based on local and remote metadata.
func evaluateFileConflict(local LocalFileInfo, remote RemoteFileInfo, mode DeploySyncMode) (DeployAction, error) {
	if !remote.Exists {
		return ActionTransferToRemote, nil
	}
	if mode == SyncModeOverwrite {
		return ActionTransferToRemote, nil
	}
	if mode == SyncModeSkip {
		return ActionSkip, nil
	}
	return evaluateSyncModeAction(local.ModTime, remote.ModTime, mode), nil
}

func evaluateSyncModeAction(localMTime, remoteMTime time.Time, mode DeploySyncMode) DeployAction {
	if mode == SyncModeSync {
		return evaluateSyncBidi(localMTime, remoteMTime)
	}
	if mode == SyncModeSyncRight {
		return evaluateSyncRight(localMTime, remoteMTime)
	}
	if mode == SyncModeSyncLeft {
		return evaluateSyncLeft(localMTime, remoteMTime)
	}
	return ActionConflictPrompt
}

func evaluateSyncBidi(localMTime, remoteMTime time.Time) DeployAction {
	diff := localMTime.Sub(remoteMTime)
	if diff > time.Second {
		return ActionTransferToRemote
	}
	if diff < -time.Second {
		return ActionTransferToLocal
	}
	return ActionSkip
}

func evaluateSyncRight(localMTime, remoteMTime time.Time) DeployAction {
	if localMTime.Sub(remoteMTime) > time.Second {
		return ActionTransferToRemote
	}
	return ActionSkip
}

func evaluateSyncLeft(localMTime, remoteMTime time.Time) DeployAction {
	if remoteMTime.Sub(localMTime) > time.Second {
		return ActionTransferToLocal
	}
	return ActionSkip
}

func resolveDeploySyncMode(opts DeployOptions) DeploySyncMode {
	if opts.IsOverwrite {
		return SyncModeOverwrite
	}
	if opts.IsSkip {
		return SyncModeSkip
	}
	if opts.IsSync {
		return SyncModeSync
	}
	return resolveDirectionalSyncMode(opts)
}

func resolveDirectionalSyncMode(opts DeployOptions) DeploySyncMode {
	if opts.IsSyncRight {
		return SyncModeSyncRight
	}
	if opts.IsSyncLeft {
		return SyncModeSyncLeft
	}
	return SyncModeNone
}
