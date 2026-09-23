//go:build !windows

package cmdos

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func probeLocalPlatformOS() (string, string, string) {
	if runtime.GOOS == "darwin" {
		out, _ := exec.Command("sw_vers", "-productVersion").Output()
		ver := strings.TrimSpace(string(out))
		return "macos", "macOS " + ver, ""
	}

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return runtime.GOOS, runtime.GOOS, ""
	}

	lines := strings.Split(string(data), "\n")
	osType := "linux"
	prettyName := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ID=") {
			osType = strings.ToLower(strings.Trim(strings.TrimPrefix(trimmed, "ID="), "\""))
		}
		if strings.HasPrefix(trimmed, "PRETTY_NAME=") {
			prettyName = strings.Trim(strings.TrimPrefix(trimmed, "PRETTY_NAME="), "\"")
		}
	}

	if prettyName == "" {
		prettyName = "Linux"
	}

	unameOut, _ := exec.Command("uname", "-r").Output()
	kernel := strings.TrimSpace(string(unameOut))

	return osType, prettyName, kernel
}
