package cmdupdate

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func executeDefaultRemoteUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	if opts.IsDryRun {
		return `{"success": true, "details": "dry-run simulated"}`, nil
	}
	if opts.IsZip {
		return ExecuteFleetZipUpdateFn(target, opts)
	}
	if out, ok := tryRestFleetUpdate(target, opts); ok {
		return out, nil
	}
	return executeSSHFleetUpdate(target, opts)
}

func tryRestFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, bool) {
	url := fmt.Sprintf("http://%s:49152/api/v1/update", target.IP)
	payload := map[string]any{
		"package": opts.Pkg,
		"force":   opts.IsForce,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}
	client := http.Client{Timeout: 3000 * time.Millisecond}
	resp, reqErr := client.Post(url, "application/json", strings.NewReader(string(b)))
	if reqErr != nil {
		return "", false
	}
	defer resp.Body.Close()
	isOk := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
	if !isOk {
		return "", false
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", false
	}
	return string(body), true
}

func executeSSHFleetUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	cmd := resolveFleetUpdateCommand(osType, opts.Pkg)
	shell := resolveFleetShell(osType)
	out, err := secrets.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

func executeSSHFleetZipUpdate(target FleetTarget, opts FleetUpdateOptions) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	osType := target.OS
	if probed := cmdssh.ProbeRemoteOSType(client); probed != "" {
		osType = probed
	}

	zipData, err := getCachedUpdateZip(opts.Pkg, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "getCachedUpdateZip")
	}

	destZipPath := resolveRemoteZipPath(osType, opts.Pkg)
	err = StreamFileToRemoteFn(client, destZipPath, zipData, osType)
	if err != nil {
		return "", apperror.WrapSimple(err, "StreamFileToRemote")
	}

	cmd := resolveFleetZipInstallCommand(osType, opts.Pkg, destZipPath)
	shell := resolveFleetShell(osType)
	out, err := secrets.RunCommand(client, cmd, shell)
	if err == nil && isAgmPkg(opts.Pkg) {
		collectAgmUpdateZip(client, osType)
	}
	return out, err
}

var (
	zipCacheMu sync.Mutex
	zipCache   = make(map[string][]byte)
)

func remoteExportZipProbe(osType string) string {
	if isWindowsOS(osType) {
		return `powershell -NoProfile -Command "Join-Path $env:USERPROFILE '.antigravity_tools\update-export\agm-update.zip'"`
	}
	return `sh -c 'printf %s "$HOME/.antigravity_tools/update-export/agm-update.zip"'`
}

func collectAgmUpdateZip(client *ssh.Client, osType string) {
	pathOut, err := secrets.RunCommand(client, remoteExportZipProbe(osType), "")
	if err != nil {
		return
	}
	data, err := StreamFileFromRemoteFn(client, strings.TrimSpace(pathOut), osType)
	if err != nil || len(data) < 4 || string(data[:2]) != "PK" {
		return
	}
	rememberCollectedZip("agm", osType, data)
}

func rememberCollectedZip(pkg, osType string, data []byte) {
	cacheKey := fmt.Sprintf("%s_%s", strings.ToLower(pkg), strings.ToLower(osType))
	zipCacheMu.Lock()
	zipCache[cacheKey] = data
	zipCacheMu.Unlock()
}

func resolveRemoteZipPath(osType, pkg string) string {
	pkgName := "gitmap"
	if isAgmPkg(pkg) {
		pkgName = "agm"
	}
	if isWindowsOS(osType) {
		return fmt.Sprintf(`C:\Windows\Temp\gitmap_update_%s.zip`, pkgName)
	}
	return fmt.Sprintf(`/tmp/gitmap_update_%s.zip`, pkgName)
}

func resolveFleetZipInstallCommand(osType, pkg, zipPath string) string {
	isWin := isWindowsOS(osType)
	targetBin := "gitmap"
	if isAgmPkg(pkg) {
		targetBin = "agm"
	}
	if isWin {
		return resolveWindowsZipInstallCommand(targetBin, zipPath)
	}
	return resolvePOSIXZipInstallCommand(targetBin, zipPath)
}

func resolveWindowsZipInstallCommand(targetBin, zipPath string) string {
	return fmt.Sprintf(`powershell -NoProfile -ExecutionPolicy Bypass -Command "& { $ErrorActionPreference = 'SilentlyContinue'; $zipPath = '%s'; $destDir = Split-Path (Get-Command %s -ErrorAction SilentlyContinue).Path; if (-not $destDir) { $destDir = [System.IO.Path]::Combine($env:LOCALAPPDATA, 'Programs', '%s') }; if (-not (Test-Path $destDir)) { New-Item -ItemType Directory -Path $destDir -Force | Out-Null }; Expand-Archive -Path $zipPath -DestinationPath $destDir -Force; Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue; $curr = (%s version 2>$null | Out-String).Trim(); if ($curr) { @{ success = $true; current_version = $curr; details = ('Zip installed: ' + $curr) } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'Zip update failed to verify binary' } | ConvertTo-Json -Compress } }"`,
		zipPath, targetBin, targetBin, targetBin)
}

func resolvePOSIXZipInstallCommand(targetBin, zipPath string) string {
	return fmt.Sprintf(`sh -c 'ZIP="%s"; DEST=$(dirname "$(command -v %s 2>/dev/null || echo /usr/local/bin/%s)"); mkdir -p "$DEST"; unzip -o "$ZIP" -d "$DEST" >/dev/null 2>&1; chmod +x "$DEST/%s"; rm -f "$ZIP"; CURR=$(%s version 2>/dev/null | head -n1); if [ -n "$CURR" ]; then printf "{\"success\":true,\"current_version\":\"%%s\",\"details\":\"Zip installed: %%s\"}" "$CURR" "$CURR"; else printf "{\"success\":false,\"details\":\"Zip update failed to verify binary\"}"; fi'`,
		zipPath, targetBin, targetBin, targetBin, targetBin)
}

func isWindowsOS(osType string) bool {
	return strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
}

func isAgmPkg(pkg string) bool {
	low := strings.ToLower(pkg)
	return low == "agm" || low == "agy" || low == "ag-manager" || low == "antigravity-manager"
}

func resolveFleetUpdateCommand(osType, pkg string) string {
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	switch strings.ToLower(pkg) {
	case "agm", "ag-manager", "antigravity-manager":
		return resolveAgmUpdateCommand(isWin)
	default:
		return resolveGitmapUpdateCommand(isWin)
	}
}

func resolveGitmapUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (gitmap version 2>$null | Out-String).Trim(); & { $env:GITMAP_UPDATING='1'; irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.ps1 | iex } *>$null; $curr = (gitmap version 2>$null | Out-String).Trim(); if ($curr) { $msg = if ($prev -and $prev -ne $curr) { 'Upgraded: ' + $prev + ' -> ' + $curr } else { 'Version: ' + $curr }; @{ success = $true; current_version = $curr; previous_version = $prev; details = $msg } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'Update failed to verify binary' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'PREV=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); GITMAP_UPDATING=1 curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.sh | bash >/dev/null 2>&1; CURR=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); if [ -n \"$CURR\" ]; then if [ -n \"$PREV\" ] && [ \"$PREV\" != \"$CURR\" ]; then DET=\"Upgraded: $PREV -> $CURR\"; else DET=\"Version: $CURR\"; fi; printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"previous_version\\\":\\\"%s\\\",\\\"details\\\":\\\"%s\\\"}\" \"$CURR\" \"$PREV\" \"$DET\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"Update failed to verify binary\\\"}\"; fi'"
}

func resolveAgmUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (agm version 2>$null | Select-Object -First 1); & { irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex } *>$null; $curr = (agm version 2>$null | Select-Object -First 1); if ($curr) { @{ success = $true; current_version = $curr; details = ('Version: ' + $curr) } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'AGM update completed' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash >/dev/null 2>&1; CURR=$(agm version 2>/dev/null | head -n1); if [ -n \"$CURR\" ]; then printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"details\\\":\\\"Version: %s\\\"}\" \"$CURR\" \"$CURR\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"AGM update completed\\\"}\"; fi'"
}

func resolveFleetShell(osType string) string {
	if strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win") {
		return ""
	}
	return "sh"
}

func dialFleetSSH(target FleetTarget) (*ssh.Client, error) {
	conn := db.SSHConnection{
		Alias:             target.Alias,
		IPAddress:         target.IP,
		Username:          target.Username,
		EncryptedPassword: target.Password,
		KeyPath:           target.KeyPath,
		OS:                target.OS,
	}
	client, isOk := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s|%s]", target.Alias, target.IP))
	if isOk && client != nil {
		return client, nil
	}
	if c, ok := tryDialFleetPassword(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetKey(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetCandidateKeys(target); ok {
		return c, nil
	}
	return nil, fmt.Errorf("ssh dial failed for %s@%s", target.Username, target.IP)
}

func tryDialFleetPassword(target FleetTarget) (*ssh.Client, bool) {
	if target.Password == "" {
		return nil, false
	}
	plain, err := secrets.DecryptStoredPassword(target.Password)
	if err != nil || plain == "" {
		plain = target.Password
	}
	c, connErr := secrets.ConnectWithPassword(target.IP, target.Username, plain)
	return c, connErr == nil
}

func tryDialFleetKey(target FleetTarget) (*ssh.Client, bool) {
	if target.KeyPath == "" {
		return nil, false
	}
	c, err := secrets.ConnectWithKey(target.IP, target.Username, target.KeyPath)
	return c, err == nil
}

func tryDialFleetCandidateKeys(target FleetTarget) (*ssh.Client, bool) {
	home, _ := os.UserHomeDir()
	candidateKeys := []string{
		fmt.Sprintf("%s/.ssh/id_ed25519", home),
		fmt.Sprintf("%s/.ssh/id_rsa", home),
	}
	for _, k := range candidateKeys {
		c, err := secrets.ConnectWithKey(target.IP, target.Username, k)
		if err == nil {
			return c, true
		}
	}
	return nil, false
}
