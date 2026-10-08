package cmdupdate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"golang.org/x/crypto/ssh"
)

func remoteExportZipProbe(osType string) string {
	if isWindowsOS(osType) {
		return `powershell -NoProfile -Command "Join-Path $env:USERPROFILE '.antigravity_tools\update-export\agm-update.zip'"`
	}
	return `sh -c 'printf %s "$HOME/.antigravity_tools/update-export/agm-update.zip"'`
}

func collectAgmUpdateZip(client *ssh.Client, osType string) {
	pathOut, err := crypto.RunCommand(client, remoteExportZipProbe(osType), "")
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
