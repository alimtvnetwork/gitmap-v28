// Package cmd — chromeprofile_reconcile.go: scans on-disk Chrome profile
// directories and reconciles unlinked profiles into Chrome's Local State
// (`profile.info_cache` and `profile.profiles_order`).
//
// Spec: spec/25-chrome-profile-management/01-profile-registration-and-picker.md.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runChromeProfileReconcile scans on-disk profiles and registers any
// directories missing from Local State so they appear in Chrome's UI.
func runChromeProfileReconcile(args []string) error {
	checkHelp(constants.SubCmdChromeReconcile, args)
	warnIfChromeRunningForReconcile()

	localStatePath := filepath.Join(chromeUserDataDir(), constants.ChromeLocalStateFile)
	root, err := loadChromeLocalStateMap(localStatePath)
	if err != nil {
		return fmt.Errorf("reconcile: %w", err)
	}

	profile := ensureChromeLocalStateProfile(root)
	infoCache := ensureChromeLocalStateInfoCache(profile)
	dirs := availableChromeProfileNames()

	fmt.Printf("\n\033[1;96mReconciling Chrome Profiles with Local State (%s)\033[0m\n\n", chromeUserDataDir())
	reconciledCount := reconcileOnDiskProfiles(dirs, infoCache, profile)

	if reconciledCount == 0 {
		fmt.Printf("\033[1;92m✓ All %d on-disk Chrome profile(s) are already synchronized in Local State.\033[0m\n\n", len(dirs))
		return nil
	}

	if err := writeChromeLocalState(localStatePath, root); err != nil {
		return fmt.Errorf("reconcile write Local State: %w", err)
	}

	fmt.Printf("\n\033[1;92m✓ Chrome Profile Reconciliation Complete:\033[0m %d profile(s) reconciled in Local State.\n\n", reconciledCount)
	return nil
}

func warnIfChromeRunningForReconcile() {
	isRunning, err := isChromeRunning(runtime.GOOS)
	if err != nil || !isRunning {
		return
	}
	fmt.Println("\n\033[1;93m⚠ Notice: Google Chrome is currently running.\033[0m")
	fmt.Println("  Chrome holds Local State in memory and flushes periodically.")
	fmt.Println("  For changes to persist permanently in the Profile Picker,")
	fmt.Println("  please restart or close Google Chrome.")
}

func loadChromeLocalStateMap(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return root, nil
}

func reconcileOnDiskProfiles(dirs []string, infoCache map[string]any, profile map[string]any) int {
	reconciledCount := 0
	for _, dir := range dirs {
		if reconcileSingleProfileDir(dir, infoCache, profile) {
			reconciledCount++
		}
	}
	return reconciledCount
}

func reconcileSingleProfileDir(dir string, infoCache map[string]any, profile map[string]any) bool {
	entry, exists := infoCache[dir].(map[string]any)
	if !exists {
		return restoreUnlinkedProfileEntry(dir, infoCache, profile)
	}
	return ensureExistingProfileInOrder(dir, entry, infoCache, profile)
}

func restoreUnlinkedProfileEntry(dir string, infoCache map[string]any, profile map[string]any) bool {
	dirPath := chromeProfilePath(dir)
	name, email := readProfilePreferencesMetadata(dirPath)
	if name == "" {
		name = dir
	}

	entry := map[string]any{}
	applyChromeInfoEntryDefaults(entry, name, email)
	infoCache[dir] = entry
	appendChromeProfileToOrder(profile, dir)
	_ = patchImportedChromeProfilePreferences(dirPath, name)

	emailStr := email
	if emailStr == "" {
		emailStr = "(none)"
	}
	fmt.Printf("  \033[1;92m+ Reconciled unlinked profile:\033[0m %-12s (display: %q, email: %s)\n", dir, name, emailStr)
	return true
}

func ensureExistingProfileInOrder(dir string, entry map[string]any, infoCache map[string]any, profile map[string]any) bool {
	changed := false
	if entry["avatar_icon"] == nil || entry["avatar_icon"] == "" {
		displayName, _ := entry["name"].(string)
		email, _ := entry["user_name"].(string)
		applyChromeInfoEntryDefaults(entry, displayName, email)
		infoCache[dir] = entry
		changed = true
	}
	if !isProfileInOrder(profile, dir) {
		appendChromeProfileToOrder(profile, dir)
		changed = true
	}
	return changed
}

func isProfileInOrder(profile map[string]any, dir string) bool {
	order, ok := profile["profiles_order"].([]any)
	if !ok {
		return false
	}
	for _, v := range order {
		if s, _ := v.(string); s == dir {
			return true
		}
	}
	return false
}

func readProfilePreferencesMetadata(dirPath string) (string, string) {
	prefPath := filepath.Join(dirPath, constants.ChromePreferencesFile)
	raw, err := os.ReadFile(prefPath)
	if err != nil {
		return "", ""
	}

	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return "", ""
	}

	displayName := extractPreferencesDisplayName(root)
	email := extractEmailFromPreferences(raw)
	return displayName, email
}

func extractPreferencesDisplayName(root map[string]any) string {
	accountName := extractAccountFullName(root)
	if accountName != "" {
		return accountName
	}
	return extractProfileBlockName(root)
}

func extractAccountFullName(root map[string]any) string {
	accounts, ok := root["account_info"].([]any)
	if !ok || len(accounts) == 0 {
		return ""
	}
	acc, ok := accounts[0].(map[string]any)
	if !ok {
		return ""
	}
	fn, _ := acc["full_name"].(string)
	return fn
}

func extractProfileBlockName(root map[string]any) string {
	prof, ok := root["profile"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := prof["name"].(string)
	if name == "Person 1" {
		return ""
	}
	return name
}
