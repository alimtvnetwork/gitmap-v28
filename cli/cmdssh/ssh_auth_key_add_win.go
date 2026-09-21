package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func resolveWinAdminKeysPath() string {
	programData := os.Getenv("ProgramData")
	hasData := programData != ""
	if hasData {
		return filepath.Join(programData, "ssh", "administrators_authorized_keys")
	}

	return filepath.Join(`C:\ProgramData`, "ssh", "administrators_authorized_keys")
}

func resolveWinUserKeysPath() string {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return ""
	}

	return filepath.Join(home, ".ssh", "authorized_keys")
}

func applyWindowsAdminACL(path string) {
	cmd := exec.Command("icacls", path, "/inheritance:r", "/grant", "SYSTEM:(F)", "BUILTIN\\Administrators:(F)")
	_ = cmd.Run()
}

func applyWindowsUserACL(path string) {
	user := os.Getenv("USERNAME")
	hasUser := user != ""
	if hasUser {
		cmd := exec.Command("icacls", path, "/inheritance:r", "/grant", user+":(F)", "SYSTEM:(F)")
		_ = cmd.Run()
	}
}

func appendSingleAuthFile(path, key, keyBlob string, isAdmin bool) bool {
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	data, _ := os.ReadFile(path)
	hasKey := isKeyInAuthorizedKeys(string(data), keyBlob)
	if hasKey {
		fmt.Printf("  ℹ Key already present in: %s\n", path)

		return false
	}

	updated := formatKeyAppend(string(data), key)
	writeErr := os.WriteFile(path, []byte(updated), 0600)
	if writeErr != nil {
		return false
	}

	applyWinACL(path, isAdmin)
	fmt.Printf("  ✓ Key added to: %s\n", path)

	return true
}

func applyWinACL(path string, isAdmin bool) {
	if isAdmin {
		applyWindowsAdminACL(path)

		return
	}

	applyWindowsUserACL(path)
}

func installWindowsAuthKeys(key, keyBlob string) int {
	addedCount := 0
	adminPath := resolveWinAdminKeysPath()
	if appendSingleAuthFile(adminPath, key, keyBlob, true) {
		addedCount++
	}

	userPath := resolveWinUserKeysPath()
	hasUserPath := userPath != ""
	if hasUserPath && appendSingleAuthFile(userPath, key, keyBlob, false) {
		addedCount++
	}

	return addedCount
}
