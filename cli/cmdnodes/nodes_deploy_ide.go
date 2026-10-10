// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RegisterRemoteIDEs registers a deployed repository in remote IDE configurations.
func RegisterRemoteIDEs(client *ssh.Client, conn db.SSHConnection, repoName, remoteRepoPath string, withPinned bool) ([]string, error) {
	isWin := isWindowsNode(conn)
	uuidStr := generateDeployUUID()
	fileURI := formatFolderFileURI(remoteRepoPath, isWin)

	if isWin {
		return registerRemoteIDEsWindows(client, repoName, remoteRepoPath, uuidStr, fileURI, withPinned)
	}

	return registerRemoteIDEsLinux(client, repoName, remoteRepoPath, uuidStr, fileURI, withPinned)
}

func registerRemoteIDEsWindows(client *ssh.Client, repoName, remoteRepoPath, uuidStr, fileURI string, withPinned bool) ([]string, error) {
	script := buildWindowsIDEScript(repoName, remoteRepoPath, uuidStr, fileURI, withPinned)
	encoded := encodePowerShell(script)
	cmd := "powershell.exe -NoProfile -EncodedCommand " + encoded
	out, err := secrets.RunCommand(client, cmd, "ps")
	if err != nil {
		return nil, fmt.Errorf("remote windows ide registration failed: %w", err)
	}

	return parseRegisteredIDEs(out), nil
}

func registerRemoteIDEsLinux(client *ssh.Client, repoName, remoteRepoPath, uuidStr, fileURI string, withPinned bool) ([]string, error) {
	script := buildLinuxIDEScript(repoName, remoteRepoPath, uuidStr, fileURI, withPinned)
	out, err := secrets.RunCommand(client, script, "sh")
	if err != nil {
		return nil, fmt.Errorf("remote linux ide registration failed: %w", err)
	}

	return parseRegisteredIDEs(out), nil
}

func parseRegisteredIDEs(output string) []string {
	lines := strings.Split(output, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "REGISTERED:") {
			return extractRegisteredTokens(trimmed)
		}
	}

	return []string{"vscode", "cursor", "antigravity"}
}

func extractRegisteredTokens(line string) []string {
	raw := strings.TrimPrefix(line, "REGISTERED:")
	parts := strings.Split(raw, ",")
	var res []string
	for _, p := range parts {
		pClean := strings.TrimSpace(p)
		if pClean != "" {
			res = append(res, pClean)
		}
	}

	return res
}

func encodePowerShell(script string) string {
	utf16Units := utf16.Encode([]rune(script))
	buf := make([]byte, len(utf16Units)*2)
	for i, r := range utf16Units {
		buf[i*2] = byte(r)
		buf[i*2+1] = byte(r >> 8)
	}

	return base64.StdEncoding.EncodeToString(buf)
}

func generateDeployUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	dst := make([]byte, 36)
	hex.Encode(dst, b[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], b[10:])

	return string(dst)
}

func formatFolderFileURI(remotePath string, isWin bool) string {
	cleanPath := filepath.ToSlash(strings.TrimSpace(remotePath))
	if isWin {
		return formatWindowsFileURI(cleanPath)
	}

	return formatUnixFileURI(cleanPath)
}

func formatWindowsFileURI(cleanPath string) string {
	if len(cleanPath) > 1 && cleanPath[1] == ':' {
		return "file:///" + string(cleanPath[0]) + "%3A" + url.PathEscape(cleanPath[2:])
	}

	return "file:///" + url.PathEscape(cleanPath)
}

func formatUnixFileURI(cleanPath string) string {
	if strings.HasPrefix(cleanPath, "/") {
		return "file://" + url.PathEscape(cleanPath)
	}

	return "file:///" + url.PathEscape(cleanPath)
}

func buildWindowsIDEScript(repoName, remoteRepoPath, uuidStr, fileURI string, withPinned bool) string {
	pinnedFlag := "$false"
	if withPinned {
		pinnedFlag = "$true"
	}
	cleanRoot := strings.ReplaceAll(remoteRepoPath, "/", "\\")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("$name = '%s'\n", repoName))
	sb.WriteString(fmt.Sprintf("$root = '%s'\n", cleanRoot))
	sb.WriteString(fmt.Sprintf("$withPinned = %s\n", pinnedFlag))
	sb.WriteString(fmt.Sprintf("$uuid = '%s'\n", uuidStr))
	sb.WriteString(fmt.Sprintf("$fileUri = '%s'\n", fileURI))
	sb.WriteString(windowsIDEScriptBody())

	return sb.String()
}

func windowsIDEScriptBody() string {
	return `$tags = @("gitmap")
if ($withPinned) {
  $tags += "pinned"
}
$newEntry = @{ name = $name; rootPath = $root; paths = @($root); tags = $tags; enabled = $true }

$appdata = [Environment]::GetFolderPath('ApplicationData')
$paths = @(
  (Join-Path $appdata 'Code\User\globalStorage\alefragnani.project-manager\projects.json'),
  (Join-Path $appdata 'Cursor\User\globalStorage\alefragnani.project-manager\projects.json')
)

foreach ($pj in $paths) {
  $dir = Split-Path $pj -Parent
  if (-not (Test-Path $dir)) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
  }
  $list = @()
  if (Test-Path $pj) {
    try {
      $raw = Get-Content $pj -Raw | ConvertFrom-Json
      if ($raw) {
        $list = @($raw)
      }
    } catch {}
  }
  $idx = -1
  for ($i = 0; $i -lt $list.Count; $i++) {
    if ($list[$i].rootPath -eq $root -or $list[$i].name -eq $name) {
      $idx = $i
      break
    }
  }
  if ($idx -ge 0) {
    $list[$idx] = $newEntry
  } else {
    $list += $newEntry
  }
  $list | ConvertTo-Json -Depth 5 | Set-Content $pj -Encoding UTF8
}

$geminiProjectsDir = Join-Path $env:USERPROFILE '.gemini\config\projects'
if (-not (Test-Path $geminiProjectsDir)) {
  New-Item -ItemType Directory -Force -Path $geminiProjectsDir | Out-Null
}
$nowStr = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ss.fffffffZ")
$agyConfig = @{
  id = $uuid
  name = $name
  projectResources = @{
    resources = @(
      @{
        gitFolder = @{
          folderUri = $fileUri
          defaultBranch = "main"
        }
      }
    )
  }
  settings = @{}
  updatedAt = $nowStr
  isWorkspaceOnly = $false
}
$agyConfig | ConvertTo-Json -Depth 5 | Set-Content (Join-Path $geminiProjectsDir "$uuid.json") -Encoding UTF8

if ($withPinned) {
  $pinnedPath = Join-Path $env:USERPROFILE '.gemini\config\pinned_projects.json'
  $pinnedDir = Split-Path $pinnedPath -Parent
  if (-not (Test-Path $pinnedDir)) {
    New-Item -ItemType Directory -Force -Path $pinnedDir | Out-Null
  }
  $pinnedStore = @{ version = "1.0.0"; updatedAt = $nowStr; projects = @() }
  if (Test-Path $pinnedPath) {
    try {
      $rawPinned = Get-Content $pinnedPath -Raw | ConvertFrom-Json
      if ($rawPinned) {
        $pinnedStore = $rawPinned
      }
    } catch {}
  }
  $pList = @($pinnedStore.projects)
  $pEntry = @{ id = $uuid; name = $name; path = $root; branch = "main"; pinnedAt = $nowStr }
  $pIdx = -1
  for ($i = 0; $i -lt $pList.Count; $i++) {
    if ($pList[$i].path -eq $root -or $pList[$i].name -eq $name) {
      $pIdx = $i
      break
    }
  }
  if ($pIdx -ge 0) {
    $pList[$pIdx] = $pEntry
  } else {
    $pList += $pEntry
  }
  $pinnedStore.projects = $pList
  $pinnedStore | ConvertTo-Json -Depth 5 | Set-Content $pinnedPath -Encoding UTF8
}

$hasGh = $false
if (Get-Command github -ErrorAction SilentlyContinue) {
  try {
    & github $root
    $hasGh = $true
  } catch {}
}
if ($hasGh) {
  Write-Output "REGISTERED:vscode,cursor,antigravity,github-desktop"
} else {
  Write-Output "REGISTERED:vscode,cursor,antigravity"
}
`
}

func buildLinuxIDEScript(repoName, remoteRepoPath, uuidStr, fileURI string, withPinned bool) string {
	pinnedStr := "False"
	if withPinned {
		pinnedStr = "True"
	}
	pyCode := fmt.Sprintf(linuxPythonScriptTemplate(), repoName, remoteRepoPath, pinnedStr, uuidStr, fileURI)
	b64 := base64.StdEncoding.EncodeToString([]byte(pyCode))

	return fmt.Sprintf(`sh -c "python3 -c \"import base64; exec(base64.b64decode('%s').decode('utf-8'))\""`, b64)
}

func linuxPythonScriptTemplate() string {
	return `import os, sys, json, subprocess
from datetime import datetime, timezone

home = os.path.expanduser('~')
name = %q
root = %q
with_pinned = %s
uuid_str = %q
file_uri = %q

tags = ["gitmap"]
if with_pinned:
    tags.append("pinned")

entry = {
    "name": name,
    "rootPath": root,
    "paths": [root],
    "tags": tags,
    "enabled": True
}

paths = [
    os.path.join(home, ".config/Code/User/globalStorage/alefragnani.project-manager/projects.json"),
    os.path.join(home, ".config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json")
]

for p in paths:
    os.makedirs(os.path.dirname(p), exist_ok=True)
    entries = []
    if os.path.exists(p):
        try:
            with open(p, 'r', encoding='utf-8') as f:
                content = json.load(f)
                if isinstance(content, list):
                    entries = content
        except Exception:
            pass
    matched = False
    for i, e in enumerate(entries):
        if e.get("rootPath") == root or e.get("name") == name:
            entries[i] = entry
            matched = True
            break
    if not matched:
        entries.append(entry)
    with open(p, 'w', encoding='utf-8') as f:
        json.dump(entries, f, indent=2)

now_iso = datetime.now(timezone.utc).isoformat()
agy_dir = os.path.join(home, ".gemini/config/projects")
os.makedirs(agy_dir, exist_ok=True)
agy_cfg = {
    "id": uuid_str,
    "name": name,
    "projectResources": {
        "resources": [
            {
                "gitFolder": {
                    "folderUri": file_uri,
                    "defaultBranch": "main"
                }
            }
        ]
    },
    "settings": {},
    "updatedAt": now_iso,
    "isWorkspaceOnly": False
}
with open(os.path.join(agy_dir, f"{uuid_str}.json"), 'w', encoding='utf-8') as f:
    json.dump(agy_cfg, f, indent=2)

if with_pinned:
    pinned_path = os.path.join(home, ".gemini/config/pinned_projects.json")
    os.makedirs(os.path.dirname(pinned_path), exist_ok=True)
    pinned_store = {"version": "1.0.0", "updatedAt": now_iso, "projects": []}
    if os.path.exists(pinned_path):
        try:
            with open(pinned_path, 'r', encoding='utf-8') as f:
                loaded = json.load(f)
                if isinstance(loaded, dict):
                    pinned_store = loaded
        except Exception:
            pass
    p_list = pinned_store.get("projects", [])
    p_entry = {"id": uuid_str, "name": name, "path": root, "branch": "main", "pinnedAt": now_iso}
    p_matched = False
    for i, item in enumerate(p_list):
        if item.get("path") == root or item.get("name") == name:
            p_list[i] = p_entry
            p_matched = True
            break
    if not p_matched:
        p_list.append(p_entry)
    pinned_store["projects"] = p_list
    with open(pinned_path, 'w', encoding='utf-8') as f:
        json.dump(pinned_store, f, indent=2)

has_gh = False
try:
    if subprocess.call("command -v github", shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL) == 0:
        subprocess.call(["github", root], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        has_gh = True
except Exception:
    pass

if has_gh:
    print("REGISTERED:vscode,cursor,antigravity,github-desktop")
else:
    print("REGISTERED:vscode,cursor,antigravity")
`
}

// DeployConversationsForRepo finds matching conversations locally and syncs them to remote.
func DeployConversationsForRepo(client *ssh.Client, conn db.SSHConnection, localRepoPath string) (int, error) {
	home, errHome := os.UserHomeDir()
	if errHome != nil {
		return 0, nil
	}

	convDir := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainDir := filepath.Join(home, ".gemini", "antigravity", "brain")
	matchingIDs := findMatchingConversationIDs(convDir, brainDir, localRepoPath)
	if len(matchingIDs) == 0 {
		return 0, nil
	}

	tarData, errPkg := packageConversationsArchive(convDir, brainDir, matchingIDs)
	if errPkg != nil {
		return 0, fmt.Errorf("package conversations failed: %w", errPkg)
	}

	return streamAndExtractConversations(client, conn, tarData, len(matchingIDs))
}

func findMatchingConversationIDs(convDir, brainDir, localRepoPath string) []string {
	entries, errRead := os.ReadDir(convDir)
	if errRead != nil {
		return nil
	}

	var matching []string
	cleanTarget := strings.ToLower(filepath.ToSlash(filepath.Clean(localRepoPath)))
	slug := strings.ToLower(filepath.Base(cleanTarget))

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
			continue
		}
		convID := strings.TrimSuffix(e.Name(), ".db")
		if isConvMatchingRepo(convDir, brainDir, convID, cleanTarget, slug) {
			matching = append(matching, convID)
		}
	}

	return matching
}

func isConvMatchingRepo(convDir, brainDir, convID, cleanTarget, slug string) bool {
	dbPath := filepath.Join(convDir, convID+".db")
	if checkConvDBMatch(dbPath, cleanTarget, slug) {
		return true
	}

	return checkConvTranscriptMatch(brainDir, convID, cleanTarget, slug)
}

func checkConvDBMatch(dbPath, cleanTarget, slug string) bool {
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return false
	}
	defer conn.Close()

	var blob []byte
	row := conn.QueryRow("SELECT data FROM trajectory_metadata_blob WHERE id='main'")
	if row.Scan(&blob) != nil || len(blob) == 0 {
		return false
	}

	lowerBlob := strings.ToLower(string(blob))
	if strings.Contains(lowerBlob, cleanTarget) || strings.Contains(lowerBlob, slug) {
		return true
	}

	return false
}

func checkConvTranscriptMatch(brainDir, convID, cleanTarget, slug string) bool {
	logDir := filepath.Join(brainDir, convID, ".system_generated", "logs")
	for _, fname := range []string{"transcript_full.jsonl", "transcript.jsonl"} {
		fpath := filepath.Join(logDir, fname)
		if checkFileContainsKeywords(fpath, cleanTarget, slug) {
			return true
		}
	}

	return false
}

func checkFileContainsKeywords(fpath, cleanTarget, slug string) bool {
	stat, err := os.Stat(fpath)
	if err != nil || stat.IsDir() {
		return false
	}

	f, errOpen := os.Open(fpath)
	if errOpen != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 65536)
	n, _ := f.Read(buf)
	if n == 0 {
		return false
	}

	lowerData := strings.ToLower(string(buf[:n]))
	if strings.Contains(lowerData, cleanTarget) || strings.Contains(lowerData, slug) {
		return true
	}

	return false
}

func packageConversationsArchive(convDir, brainDir string, convIDs []string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for _, id := range convIDs {
		archiveConvDBFiles(tw, convDir, id)
		archiveBrainFolder(tw, brainDir, id)
	}

	_ = tw.Close()
	_ = gw.Close()

	return buf.Bytes(), nil
}

func archiveConvDBFiles(tw *tar.Writer, convDir, id string) {
	for _, ext := range []string{".db", ".db-wal", ".db-shm"} {
		fPath := filepath.Join(convDir, id+ext)
		stat, err := os.Stat(fPath)
		if err == nil && !stat.IsDir() {
			writeTarFileEntry(tw, "conversations/"+id+ext, fPath)
		}
	}
}

func archiveBrainFolder(tw *tar.Writer, brainDir, id string) {
	src := filepath.Join(brainDir, id)
	stat, err := os.Stat(src)
	if err != nil || !stat.IsDir() {
		return
	}

	_ = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(brainDir, path)
		if relErr != nil {
			return nil
		}
		writeTarFileEntry(tw, "brain/"+filepath.ToSlash(rel), path)
		return nil
	})
}

func streamAndExtractConversations(client *ssh.Client, conn db.SSHConnection, tarData []byte, count int) (int, error) {
	stagingPath := resolveConvStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, stagingPath, tarData, conn.OS); errStream != nil {
		return 0, fmt.Errorf("stream conversations failed: %w", errStream)
	}

	cmd, shell := resolveConvExtractCommand(conn, stagingPath)
	if _, errExec := secrets.RunCommand(client, cmd, shell); errExec != nil {
		return 0, fmt.Errorf("extract conversations failed: %w", errExec)
	}

	return count, nil
}

func resolveConvStagingPath(conn db.SSHConnection) string {
	if isWindowsNode(conn) {
		return `C:\Windows\Temp\convs_deploy.tar.gz`
	}

	return "/tmp/convs_deploy.tar.gz"
}

func resolveConvExtractCommand(conn db.SSHConnection, stagingPath string) (string, string) {
	if isWindowsNode(conn) {
		cmd := fmt.Sprintf(`powershell.exe -NoProfile -Command "$dest = Join-Path $env:USERPROFILE '.gemini\antigravity'; if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Force -Path $dest | Out-Null }; tar.exe -xzf '%s' -C $dest; Remove-Item '%s' -Force -ErrorAction SilentlyContinue"`, stagingPath, stagingPath)
		return cmd, "ps"
	}

	cmd := fmt.Sprintf(`sh -c "mkdir -p ~/.gemini/antigravity && tar -xzf %s -C ~/.gemini/antigravity && rm -f %s"`, stagingPath, stagingPath)
	return cmd, "sh"
}
