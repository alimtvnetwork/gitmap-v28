// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

const (
	PrefixRegistered          = "REGISTERED:"
	TokenVSCode               = "vscode"
	TokenCursor               = "cursor"
	TokenAntigravity          = "antigravity"
	TokenDesktop              = "github-desktop"
	RegisteredHeaderAll       = PrefixRegistered + TokenVSCode + "," + TokenCursor + "," + TokenAntigravity + "," + TokenDesktop
	RegisteredHeaderNoDesktop = PrefixRegistered + TokenVSCode + "," + TokenCursor + "," + TokenAntigravity
)

const windowsProjectsScript = `$tags = @("gitmap")
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
`

const windowsGeminiScript = `$geminiProjectsDir = Join-Path $env:USERPROFILE '.gemini\config\projects'
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
`

const windowsGitHubDesktopScript = `$hasGh = $false
if (Get-Command github -ErrorAction SilentlyContinue) {
  try {
    & github $root
    $hasGh = $true
  } catch {}
}
if ($hasGh) {
  Write-Output "` + RegisteredHeaderAll + `"
} else {
  Write-Output "` + RegisteredHeaderNoDesktop + `"
}
`

// RegisterRemoteIDEs registers a deployed repository in remote IDE configurations using pure Go gitmap commands.
func RegisterRemoteIDEs(client *ssh.Client, conn db.SSHConnection, repoName, remoteRepoPath string, hasPinned bool) ([]string, error) {
	if isWindowsNode(conn) {
		return registerRemoteIDEsWindows(client, repoName, remoteRepoPath, hasPinned)
	}

	return registerRemoteIDEsLinux(client, repoName, remoteRepoPath, hasPinned)
}

func registerRemoteIDEsWindows(client *ssh.Client, repoName, remoteRepoPath string, hasPinned bool) ([]string, error) {
	cmd := buildWindowsRemoteIDECommand(repoName, remoteRepoPath, hasPinned)
	out, err := secrets.RunCommand(client, cmd, "ps")
	if err != nil {
		return nil, apperror.WrapSimple(err, "remote windows ide registration failed")
	}

	return parseRegisteredIDEs(out), nil
}

func registerRemoteIDEsLinux(client *ssh.Client, repoName, remoteRepoPath string, hasPinned bool) ([]string, error) {
	cmd := buildLinuxRemoteIDECommand(repoName, remoteRepoPath, hasPinned)
	out, err := secrets.RunCommand(client, cmd, "sh")
	if err != nil {
		return nil, apperror.WrapSimple(err, "remote linux ide registration failed")
	}

	return parseRegisteredIDEs(out), nil
}

func parseRegisteredIDEs(output string) []string {
	lines := strings.Split(output, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, PrefixRegistered) {
			return extractRegisteredTokens(trimmed)
		}
	}

	return []string{TokenVSCode, TokenCursor, TokenAntigravity}
}

func extractRegisteredTokens(line string) []string {
	raw := strings.TrimPrefix(line, PrefixRegistered)
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

func formatPinnedFlag(hasPinned bool) string {
	if hasPinned {
		return "$true"
	}

	return "$false"
}

func formatPinnedCliFlag(hasPinned bool) string {
	if hasPinned {
		return " --pinned"
	}

	return ""
}

func buildWindowsRemoteIDECommand(repoName, remoteRepoPath string, hasPinned bool) string {
	pinnedFlag := formatPinnedCliFlag(hasPinned)
	cleanRoot := strings.ReplaceAll(remoteRepoPath, "/", "\\")
	fallback := buildWindowsFallbackScript(repoName, cleanRoot, hasPinned)
	script := fmt.Sprintf("if (Get-Command gitmap -ErrorAction SilentlyContinue) { & gitmap ide register '%s' --name '%s'%s } else { %s }",
		cleanRoot, repoName, pinnedFlag, fallback)

	return "powershell.exe -NoProfile -EncodedCommand " + encodePowerShell(script)
}

func buildWindowsFallbackScript(repoName, cleanRoot string, hasPinned bool) string {
	uuidStr := uuidutil.GenerateUUID()
	fileURI := formatWindowsFileURI(cleanRoot)
	header := fmt.Sprintf("$name = '%s'; $root = '%s'; $withPinned = %s; $uuid = '%s'; $fileUri = '%s'; ",
		repoName, cleanRoot, formatPinnedFlag(hasPinned), uuidStr, fileURI)

	return header + windowsProjectsScript + windowsGeminiScript + windowsGitHubDesktopScript
}

func buildLinuxRemoteIDECommand(repoName, remoteRepoPath string, hasPinned bool) string {
	pinnedFlag := formatPinnedCliFlag(hasPinned)
	fallback := buildLinuxFallbackScript(repoName, remoteRepoPath)

	return fmt.Sprintf(`sh -c "if command -v gitmap >/dev/null 2>&1; then gitmap ide register '%s' --name '%s'%s; elif [ -x \"$HOME/.local/bin/gitmap\" ]; then \"$HOME/.local/bin/gitmap\" ide register '%s' --name '%s'%s; elif [ -x \"/usr/local/bin/gitmap\" ]; then \"/usr/local/bin/gitmap\" ide register '%s' --name '%s'%s; else %s; fi"`,
		remoteRepoPath, repoName, pinnedFlag,
		remoteRepoPath, repoName, pinnedFlag,
		remoteRepoPath, repoName, pinnedFlag,
		fallback)
}

func buildLinuxFallbackScript(repoName, remoteRepoPath string) string {
	uuidStr := uuidutil.GenerateUUID()
	cleanPath := filepath.ToSlash(strings.TrimSpace(remoteRepoPath))
	fileURI := formatUnixFileURI(cleanPath)

	return fmt.Sprintf("mkdir -p ~/.gemini/config/projects && printf '{\"id\":\"%s\",\"name\":\"%s\",\"projectResources\":{\"resources\":[{\"gitFolder\":{\"folderUri\":\"%s\",\"defaultBranch\":\"main\"}}]},\"settings\":{},\"isWorkspaceOnly\":false}' > ~/.gemini/config/projects/%s.json && echo "+PrefixRegistered+TokenAntigravity,
		uuidStr, repoName, fileURI, uuidStr)
}

func resolveLocalGeminiDirs() (string, string, bool) {
	home, errHome := os.UserHomeDir()
	if errHome != nil {
		return "", "", false
	}

	convDir := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainDir := filepath.Join(home, ".gemini", "antigravity", "brain")

	return convDir, brainDir, true
}

// DeployConversationsForRepo finds matching conversations locally and syncs them to remote.
func DeployConversationsForRepo(client *ssh.Client, conn db.SSHConnection, localRepoPath string) (int, error) {
	convDir, brainDir, isFound := resolveLocalGeminiDirs()
	if !isFound {
		return 0, nil
	}

	matchingIDs := findMatchingConversationIDs(convDir, brainDir, localRepoPath)
	if len(matchingIDs) == 0 {
		return 0, nil
	}

	return deployMatchingConversations(client, conn, convDir, brainDir, matchingIDs)
}

func deployMatchingConversations(client *ssh.Client, conn db.SSHConnection, convDir, brainDir string, matchingIDs []string) (int, error) {
	tarData, errPkg := packageConversationsArchive(convDir, brainDir, matchingIDs)
	if errPkg != nil {
		return 0, apperror.WrapSimple(errPkg, "package conversations failed")
	}

	return streamAndExtractConversations(client, conn, tarData, len(matchingIDs))
}

func findMatchingConversationIDs(convDir, brainDir, localRepoPath string) []string {
	entries, errRead := os.ReadDir(convDir)
	if errRead != nil {
		return nil
	}

	cleanTarget := strings.ToLower(filepath.ToSlash(filepath.Clean(localRepoPath)))
	slug := strings.ToLower(filepath.Base(cleanTarget))

	return collectMatchingIDs(entries, convDir, brainDir, cleanTarget, slug)
}

func collectMatchingIDs(entries []os.DirEntry, convDir, brainDir, cleanTarget, slug string) []string {
	var matching []string
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
	if hasConvDBMatch(dbPath, cleanTarget, slug) {
		return true
	}

	return hasConvTranscriptMatch(brainDir, convID, cleanTarget, slug)
}

func hasConvDBMatch(dbPath, cleanTarget, slug string) bool {
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return false
	}
	conn := openRes.Value

	defer conn.Close()
	blob := readTrajectoryBlob(conn)
	if len(blob) == 0 {
		return false
	}

	return hasTargetOrSlug(string(blob), cleanTarget, slug)
}

func readTrajectoryBlob(conn *sql.DB) []byte {
	var blob []byte
	row := conn.QueryRow("SELECT data FROM trajectory_metadata_blob WHERE id='main'")
	if row.Scan(&blob) != nil {
		return nil
	}

	return blob
}

func hasTargetOrSlug(text, cleanTarget, slug string) bool {
	lowerText := strings.ToLower(text)

	return strings.Contains(lowerText, cleanTarget) || strings.Contains(lowerText, slug)
}

func hasConvTranscriptMatch(brainDir, convID, cleanTarget, slug string) bool {
	logDir := filepath.Join(brainDir, convID, ".system_generated", "logs")
	for _, fname := range []string{"transcript_full.jsonl", "transcript.jsonl"} {
		fpath := filepath.Join(logDir, fname)
		if hasFileKeywords(fpath, cleanTarget, slug) {
			return true
		}
	}

	return false
}

func isFileReadable(fpath string) bool {
	stat, err := os.Stat(fpath)

	return err == nil && !stat.IsDir()
}

func readOpenFile(f *os.File, maxBytes int) []byte {
	buf := make([]byte, maxBytes)
	n, _ := f.Read(buf)

	return buf[:n]
}

func readFilePrefix(fpath string, maxBytes int) []byte {
	if !isFileReadable(fpath) {
		return nil
	}

	f, errOpen := os.Open(fpath)
	if errOpen != nil {
		return nil
	}

	defer f.Close()

	return readOpenFile(f, maxBytes)
}

func hasFileKeywords(fpath, cleanTarget, slug string) bool {
	data := readFilePrefix(fpath, 65536)
	if len(data) == 0 {
		return false
	}

	return hasTargetOrSlug(string(data), cleanTarget, slug)
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

	_ = filepath.Walk(src, buildBrainWalkFn(tw, brainDir))
}

func archiveBrainPath(tw *tar.Writer, brainDir, path string, info os.FileInfo) {
	if info.IsDir() {
		return
	}

	rel, relErr := filepath.Rel(brainDir, path)
	if relErr != nil {
		return
	}

	writeTarFileEntry(tw, "brain/"+filepath.ToSlash(rel), path)
}

func buildBrainWalkFn(tw *tar.Writer, brainDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		archiveBrainPath(tw, brainDir, path, info)

		return nil
	}
}

func streamAndExtractConversations(client *ssh.Client, conn db.SSHConnection, tarData []byte, count int) (int, error) {
	stagingPath := resolveConvStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, stagingPath, tarData, conn.OS); errStream != nil {
		return 0, apperror.WrapSimple(errStream, "stream conversations failed")
	}

	return runRemoteConvExtract(client, conn, stagingPath, count)
}

func runRemoteConvExtract(client *ssh.Client, conn db.SSHConnection, stagingPath string, count int) (int, error) {
	cmd, shell := resolveConvExtractCommand(conn, stagingPath)
	if _, errExec := secrets.RunCommand(client, cmd, shell); errExec != nil {
		return 0, apperror.WrapSimple(errExec, "extract conversations failed")
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
