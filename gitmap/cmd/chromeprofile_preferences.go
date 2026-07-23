// Package cmd — chromeprofile_preferences.go: post-copy patcher for the
// destination profile's `Preferences` file. Chrome's profile picker
// reads BOTH `Local State` AND the per-profile `Preferences` JSON; if
// the destination Preferences still carries the source GAIA/signed-in
// fields (or a stale profile.name), Chrome may hide the new tile or
// silently merge it back into the source identity on next launch.
//
// This file scrubs those fields and stamps the picker-visible name so
// the freshly copied dir reliably appears in chrome://settings/manageProfile.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// patchCopiedChromeProfilePreferences rewrites <dst>/Preferences so the
// new profile is signed-out and named after the destination slug.
// Soft-fails: any error is returned to the caller as a warning, never
// fatal — the on-disk copy succeeded regardless.
func patchCopiedChromeProfilePreferences(dstPath, displayName string) error {
	prefPath := filepath.Join(dstPath, constants.ChromePreferencesFile)
	raw, err := os.ReadFile(prefPath) //nolint:gosec // chrome user-data path
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", prefPath, err)
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("parse %s: %w", prefPath, err)
	}
	scrubChromePreferencesIdentity(root, displayName)
	out, err := json.MarshalIndent(root, "", constants.JSONIndent)
	if err != nil {
		return fmt.Errorf("encode Preferences: %w", err)
	}
	if err := os.WriteFile(prefPath, out, constants.FilePermission); err != nil {
		return fmt.Errorf("write %s: %w", prefPath, err)
	}
	return nil
}

func scrubChromePreferencesIdentity(root map[string]any, displayName string) {
	delete(root, "account_info")
	delete(root, "signin")
	delete(root, "google")
	delete(root, "gaia_cookie")
	prof, ok := root["profile"].(map[string]any)
	if !ok {
		prof = map[string]any{}
		root["profile"] = prof
	}
	prof["name"] = displayName
	prof["using_default_name"] = false
	delete(prof, "gaia_info_picture_url")
	delete(prof, "gaia_given_name")
	delete(prof, "gaia_name")
	delete(prof, "last_used")
}
