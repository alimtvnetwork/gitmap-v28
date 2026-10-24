// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
	"golang.org/x/crypto/ssh"
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
