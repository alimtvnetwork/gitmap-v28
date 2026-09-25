// Package cmdagy — agy_ls_discovery.go dynamically discovers running Antigravity Language Server address and CSRF token.
package cmdagy

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	cachedLSAddress string
	cachedCSRFToken string
	cachedLSEnvMu   sync.RWMutex
)

// ResolveAntigravityLSEnv returns the discovered address and CSRF token for the running language_server.
func ResolveAntigravityLSEnv() (string, string, bool) {
	envAddr := os.Getenv("ANTIGRAVITY_LS_ADDRESS")
	envToken := os.Getenv("ANTIGRAVITY_CSRF_TOKEN")
	hasDirectEnv := len(envAddr) > 0 && len(envToken) > 0
	if hasDirectEnv {
		return envAddr, envToken, true
	}

	cachedAddr, cachedToken, hasCached := getCachedLSEnv()
	if hasCached {
		return cachedAddr, cachedToken, true
	}

	addr, token, ok := discoverLSEnv()
	if ok {
		setCachedLSEnv(addr, token)
		return addr, token, true
	}

	return "", "", false
}

// InvalidateAntigravityLSEnvCache clears cached LS discovery on connection failure.
func InvalidateAntigravityLSEnvCache() {
	cachedLSEnvMu.Lock()
	cachedLSAddress = ""
	cachedCSRFToken = ""
	cachedLSEnvMu.Unlock()
}

func getCachedLSEnv() (string, string, bool) {
	cachedLSEnvMu.RLock()
	defer cachedLSEnvMu.RUnlock()
	hasCached := len(cachedLSAddress) > 0 && len(cachedCSRFToken) > 0
	return cachedLSAddress, cachedCSRFToken, hasCached
}

func setCachedLSEnv(addr, token string) {
	cachedLSEnvMu.Lock()
	cachedLSAddress = addr
	cachedCSRFToken = token
	cachedLSEnvMu.Unlock()
}

func discoverLSEnv() (string, string, bool) {
	if runtime.GOOS == "windows" {
		return discoverLSEnvWindows()
	}

	return discoverLSEnvUnix()
}

func discoverLSEnvWindows() (string, string, bool) {
	psScript := `$p = Get-CimInstance Win32_Process -Filter "name='language_server.exe'" | Select-Object -First 1; if ($p) { $token = ''; if ($p.CommandLine -match '--csrf_token\s+([a-zA-Z0-9\-]+)') { $token = $matches[1] }; $ports = Get-NetTCPConnection -OwningProcess $p.ProcessId -State Listen | Where-Object { $_.LocalAddress -eq '127.0.0.1' } | Select-Object -ExpandProperty LocalPort; "$($p.ProcessId)|$token|$($ports -join ',')" }`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	out, err := cmd.Output()
	hasErr := err != nil
	if hasErr {
		return "", "", false
	}

	return parseLSEnvWindowsOutput(strings.TrimSpace(string(out)))
}

func parseLSEnvWindowsOutput(raw string) (string, string, bool) {
	hasEmpty := len(raw) == 0
	if hasEmpty {
		return "", "", false
	}
	parts := strings.Split(raw, "|")
	hasValidParts := len(parts) >= 3
	if hasValidParts == false {
		return "", "", false
	}
	token := strings.TrimSpace(parts[1])
	portsStr := strings.TrimSpace(parts[2])

	return selectWorkingLSPort(portsStr, token)
}

func selectWorkingLSPort(portsStr, token string) (string, string, bool) {
	portList := strings.Split(portsStr, ",")
	for _, port := range portList {
		trimmedPort := strings.TrimSpace(port)
		if len(trimmedPort) == 0 {
			continue
		}
		if probeLSPort(trimmedPort, token) {
			return "127.0.0.1:" + trimmedPort, token, true
		}
	}

	return "", "", false
}

func probeLSPort(port, token string) bool {
	binPath, baseArgs, isResolved := ResolveAgentAPI()
	if isResolved == false {
		return false
	}
	addr := "127.0.0.1:" + port
	fullArgs := append([]string{}, baseArgs...)
	fullArgs = append(fullArgs, "get-conversation-metadata", "probe-check")
	cmd := exec.Command(binPath, fullArgs...)
	cmd.Env = append(os.Environ(),
		"ANTIGRAVITY_LS_ADDRESS="+addr,
		"ANTIGRAVITY_CSRF_TOKEN="+token,
	)
	out, _ := cmd.CombinedOutput()

	return isProbeOutputValid(string(out))
}

func isProbeOutputValid(outStr string) bool {
	hasConnErr := strings.Contains(outStr, "connection error") ||
		strings.Contains(outStr, "connection refused") ||
		strings.Contains(outStr, "forcibly closed") ||
		strings.Contains(outStr, "server preface")

	return !hasConnErr
}

func discoverLSEnvUnix() (string, string, bool) {
	return "", "", false
}
