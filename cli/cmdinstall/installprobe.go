package cmdinstall

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type toolProbeConfig struct {
	bins []string
	args []string
}

var versionRegex = regexp.MustCompile(`v?(\d+\.\d+(?:\.\d+)*)`)

var toolProbeMap = map[string]toolProbeConfig{
	constants.ToolVSCode:         {bins: []string{"code", "code.cmd", "vscode"}, args: []string{"--version"}},
	constants.ToolCPP:            {bins: []string{"g++", "clang++", "gcc", "cpp"}, args: []string{"--version"}},
	constants.ToolPowerShell:     {bins: []string{"pwsh", "powershell"}, args: []string{"--version"}},
	constants.ToolPython:         {bins: []string{"python", "python3"}, args: []string{"--version"}},
	constants.ToolNodeJS:         {bins: []string{"node", "nodejs"}, args: []string{"-v"}},
	constants.ToolGo:             {bins: []string{"go"}, args: []string{"version"}},
	constants.ToolBun:            {bins: []string{"bun"}, args: []string{"-v"}},
	constants.ToolPnpm:           {bins: []string{"pnpm"}, args: []string{"-v"}},
	constants.ToolYarn:           {bins: []string{"yarn"}, args: []string{"-v"}},
	constants.ToolPHP:            {bins: []string{"php"}, args: []string{"-v"}},
	constants.ToolChrome:         {bins: []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "chrome"}, args: []string{"--version"}},
	constants.ToolGoogleChrome:   {bins: []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "chrome"}, args: []string{"--version"}},
	constants.ToolGitLFS:         {bins: []string{"git-lfs"}, args: []string{"--version"}},
	constants.ToolJava:           {bins: []string{"java"}, args: []string{"-version"}},
	constants.ToolRust:           {bins: []string{"rustc", "cargo"}, args: []string{"--version"}},
	constants.ToolDotnet:         {bins: []string{"dotnet"}, args: []string{"--version"}},
	constants.ToolKubernetes:     {bins: []string{"kubectl"}, args: []string{"version", "--client"}},
	constants.ToolOpenVmTools:    {bins: []string{"vmtoolsd"}, args: []string{"-v"}},
	constants.ToolVMware:         {bins: []string{"vmware", "vmtoolsd", "vmhgfs-fuse"}, args: []string{"-v"}},
	constants.ToolNginx:          {bins: []string{"nginx"}, args: []string{"-v"}},
	constants.ToolWpCli:          {bins: []string{"wp"}, args: []string{"--version"}},
	constants.ToolWordPress:      {bins: []string{"wp"}, args: []string{"--version"}},
	constants.ToolSQLite:         {bins: []string{"sqlite3", "sqlite"}, args: []string{"--version"}},
	constants.ToolPostgreSQL:     {bins: []string{"psql"}, args: []string{"--version"}},
	constants.ToolMySQL:          {bins: []string{"mysql", "mariadb"}, args: []string{"--version"}},
	constants.ToolMariaDB:        {bins: []string{"mariadb", "mysql"}, args: []string{"--version"}},
	constants.ToolRedis:          {bins: []string{"redis-server", "redis-cli"}, args: []string{"--version"}},
	constants.ToolMongoDB:        {bins: []string{"mongod", "mongosh"}, args: []string{"--version"}},
	constants.ToolBuildEssential: {bins: []string{"gcc", "make", "g++"}, args: []string{"--version"}},
	constants.ToolQBittorrent:    {bins: []string{"qbittorrent", "qbittorrent-nox"}, args: []string{"--version"}},
	constants.ToolUTorrent:       {bins: []string{"utorrent", "uTorrent", "utserver"}, args: []string{"--version"}},
	constants.ToolZsh:            {bins: []string{"zsh"}, args: []string{"--version"}},
	constants.ToolAntigravity:    {bins: []string{"antigravity", "Antigravity"}, args: []string{"--version"}},
	constants.ToolAgy:            {bins: []string{"agy"}, args: []string{"--version"}},
	constants.ToolAgManager:      {bins: []string{"ag-manager", "Antigravity.Tools", "Antigravity-Manager"}, args: []string{"--version"}},
}

func resolveToolCandidates(tool string) ([]string, []string) {
	if cfg, ok := toolProbeMap[tool]; ok {
		return cfg.bins, cfg.args
	}

	return []string{tool}, []string{"--version"}
}

func resolvePowerShellArgs(bin string, defaultArgs []string) []string {
	if bin == "powershell" && !isBinaryInPath("pwsh") {
		return []string{"-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()"}
	}

	return defaultArgs
}

func probeCandidate(bin string, args []string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	target := resolveToolBinaryPath(bin)
	if target == "" {
		target = bin
	}

	cmd := exec.CommandContext(ctx, target, args...)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return ""
	}

	outStr := string(out)
	if strings.Contains(outStr, "Python was not found") || strings.Contains(outStr, "Microsoft Store") {
		return ""
	}

	return parseVersionFromOutput(outStr)
}

func probeGitLFSFallback() (string, string) {
	if isBinaryInPath("git") {
		ver := probeCandidate("git", []string{"lfs", "version"})

		return "git", ver
	}

	return "", ""
}

func probeSingleCandidate(bin string, defaultArgs []string) (string, bool) {
	if !isBinaryInPath(bin) {
		return "", false
	}

	args := resolvePowerShellArgs(bin, defaultArgs)
	ver := probeCandidate(bin, args)

	return ver, true
}

func resolveToolProbeCommand(tool string) (string, string) {
	candidates, defaultArgs := resolveToolCandidates(tool)
	var fallbackBin string
	for _, bin := range candidates {
		ver, isPresent := probeSingleCandidate(bin, defaultArgs)
		if !isPresent {
			continue
		}

		if fallbackBin == "" {
			fallbackBin = bin
		}

		if ver != "" {
			return bin, ver
		}
	}

	if tool == constants.ToolGitLFS {
		return probeGitLFSFallback()
	}

	if tool == constants.ToolAntigravity {
		if path, isFound := findInstalledAntigravityDesktopPath(); isFound {
			return path, "installed"
		}
	}

	if fallbackBin != "" {
		return fallbackBin, "installed"
	}

	return "", ""
}

func parseVersionFromOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		m := versionRegex.FindString(line)
		if m != "" {
			return m
		}
	}

	return ""
}

func isBinaryInPath(tool string) bool {
	path := resolveToolBinaryPath(tool)

	return path != ""
}
