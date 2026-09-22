//go:build linux

package cmdos

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const gdm3ConfigPath = "/etc/gdm3/custom.conf"

func configureGDM3(user string) error {
	content, err := os.ReadFile(gdm3ConfigPath)
	if err != nil {
		return apperror.WrapSimple(err, "failed reading GDM3 config")
	}
	updated := patchGDM3Content(string(content), user)
	if err := os.WriteFile(gdm3ConfigPath, []byte(updated), 0644); err != nil {
		return apperror.WrapSimple(err, "failed writing GDM3 config (run with sudo)")
	}
	return nil
}

func patchGDM3Content(text, user string) string {
	lines := strings.Split(text, "\n")
	var result []string
	hasDaemonSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "[daemon]" {
			hasDaemonSection = true
			result = append(result, line, "AutomaticLoginEnable=true", fmt.Sprintf("AutomaticLogin=%s", user))
			continue
		}
		if isGDM3AutoLoginKey(trimmed) {
			continue
		}
		result = append(result, line)
	}
	if !hasDaemonSection {
		result = append(result, "[daemon]", "AutomaticLoginEnable=true", fmt.Sprintf("AutomaticLogin=%s", user))
	}
	return strings.Join(result, "\n")
}

func isGDM3AutoLoginKey(line string) bool {
	return strings.HasPrefix(line, "AutomaticLoginEnable") || strings.HasPrefix(line, "AutomaticLogin")
}

func disableGDM3() error {
	content, err := os.ReadFile(gdm3ConfigPath)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		if isGDM3AutoLoginKey(strings.TrimSpace(line)) {
			continue
		}
		result = append(result, line)
	}
	return os.WriteFile(gdm3ConfigPath, []byte(strings.Join(result, "\n")), 0644)
}

func readGDM3Status() (AutoLoginStatus, error) {
	content, err := os.ReadFile(gdm3ConfigPath)
	if err != nil {
		return AutoLoginStatus{DisplayManager: "GDM3"}, nil
	}
	enabled, user := parseGDM3Lines(strings.Split(string(content), "\n"))
	return AutoLoginStatus{IsEnabled: enabled, Username: user, DisplayManager: "GDM3"}, nil
}

func parseGDM3Lines(lines []string) (bool, string) {
	enabled := false
	user := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "AutomaticLoginEnable=true") {
			enabled = true
		}
		if strings.HasPrefix(trimmed, "AutomaticLogin=") && !strings.HasPrefix(trimmed, "AutomaticLoginEnable") {
			user = strings.TrimPrefix(trimmed, "AutomaticLogin=")
		}
	}
	return enabled, user
}
