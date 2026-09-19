package cmdinstall

import "strings"

func isAgManagerInstalled() (string, bool) {
	for _, name := range []string{"ag-manager", "Antigravity.Tools", "Antigravity-Manager"} {
		if bin := resolveToolBinaryPath(name); bin != "" {
			return "installed", true
		}
	}

	if isAgManagerWindowsInstalled() {
		return "installed", true
	}

	return "", false
}

func isAppImageFile(path string) bool {
	low := strings.ToLower(path)

	return strings.HasSuffix(low, ".appimage")
}
