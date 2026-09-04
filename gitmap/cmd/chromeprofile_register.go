// Package cmd — chromeprofile_register.go: registers a freshly copied
// destination profile inside Chrome's `Local State` so the profile
// shows up in Chrome's profile picker. Without this step, copying
// `Profile 15` → `lv2` lands the files on disk but Chrome ignores
// the new directory because it's not listed in
// `profile.info_cache[<dir>]`.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// registerChromeProfileInLocalState clones the source profile's
// `info_cache` entry into the destination dir slot, scrubs the
// signed-in identity fields, sets the visible name, and appends the
// dir to `profiles_order` when present. Soft-fails: any error is
// returned so the caller can warn without aborting the whole copy.
func registerChromeProfileInLocalState(srcDir, dstDir, displayName string) error {
	path := filepath.Join(chromeUserDataDir(), constants.ChromeLocalStateFile)
	root, err := readOrCreateLocalStateRoot(path)
	if err != nil {
		return err
	}
	profile := ensureChromeLocalStateProfile(root)
	infoCache := ensureChromeLocalStateInfoCache(profile)
	infoCache[dstDir] = buildChromeDestinationInfoEntry(infoCache, srcDir, displayName)
	appendChromeProfileToOrder(profile, dstDir)
	return writeChromeLocalState(path, root)
}

// registerChromeProfileWithFullSchema registers a profile in Local State
// ensuring all 13 Chromium UI attributes are populated so Chrome's
// profile picker UI displays the profile tile without dropping it.
func registerChromeProfileWithFullSchema(dstDir, displayName, email string) error {
	path := filepath.Join(chromeUserDataDir(), constants.ChromeLocalStateFile)
	root, err := readOrCreateLocalStateRoot(path)
	if err != nil {
		return err
	}
	profile := ensureChromeLocalStateProfile(root)
	infoCache := ensureChromeLocalStateInfoCache(profile)
	entry, ok := infoCache[dstDir].(map[string]any)
	if !ok {
		entry = map[string]any{}
	}
	for _, k := range chromeInfoCacheGAIAFields {
		delete(entry, k)
	}
	applyChromeInfoEntryDefaults(entry, displayName, email)
	infoCache[dstDir] = entry
	appendChromeProfileToOrder(profile, dstDir)
	return writeChromeLocalState(path, root)
}

func readOrCreateLocalStateRoot(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	if len(raw) == 0 {
		return map[string]any{}, nil
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return root, nil
}

func ensureChromeLocalStateProfile(root map[string]any) map[string]any {
	if p, ok := root["profile"].(map[string]any); ok {
		return p
	}
	p := map[string]any{}
	root["profile"] = p
	return p
}

func ensureChromeLocalStateInfoCache(profile map[string]any) map[string]any {
	if c, ok := profile["info_cache"].(map[string]any); ok {
		return c
	}
	c := map[string]any{}
	profile["info_cache"] = c
	return c
}

// buildChromeDestinationInfoEntry clones the source entry (preserving
// avatar/theme so the picker tile looks identical) and scrubs every
// GAIA / signed-in account field — the destination profile must look
// signed-out to Chrome, matching the cookies/login data we deliberately
// excluded from the copy.
func buildChromeDestinationInfoEntry(
	infoCache map[string]any,
	srcDir,
	displayName string,
) map[string]any {
	entry := map[string]any{}
	if src, ok := infoCache[srcDir].(map[string]any); ok {
		for k, v := range src {
			entry[k] = v
		}
	}
	for _, k := range chromeInfoCacheGAIAFields {
		delete(entry, k)
	}
	applyChromeInfoEntryDefaults(entry, displayName, "")
	return entry
}

func applyChromeInfoEntryDefaults(entry map[string]any, displayName, email string) {
	if displayName != "" {
		entry["name"] = displayName
		entry["shortcut_name"] = displayName
	}
	entry["is_using_default_name"] = false
	entry["is_ephemeral"] = false
	entry["is_consented_primary_account"] = false
	entry["signin.with_credential_provider"] = false
	if email != "" {
		entry["user_name"] = email
	}
	if entry["avatar_icon"] == nil || entry["avatar_icon"] == "" {
		entry["avatar_icon"] = "chrome://theme/IDR_PROFILE_AVATAR_26"
		entry["is_using_default_avatar"] = true
	}
	if entry["default_avatar_fill_color"] == nil {
		entry["default_avatar_fill_color"] = -13625057
	}
	if entry["default_avatar_stroke_color"] == nil {
		entry["default_avatar_stroke_color"] = -1786428
	}
	if entry["profile_highlight_color"] == nil {
		entry["profile_highlight_color"] = -13625057
	}
	if entry["profile_color_seed"] == nil {
		entry["profile_color_seed"] = -4385188
	}
	if entry["active_time"] == nil {
		entry["active_time"] = float64(time.Now().Unix())
	}
}

// chromeInfoCacheGAIAFields are the signed-in identity keys we strip
// from the cloned entry so the new tile does not inherit the source
// account's email / avatar URL / hosted domain.
var chromeInfoCacheGAIAFields = []string{
	"gaia_id",
	"gaia_name",
	"gaia_given_name",
	"gaia_picture_file_name",
	"user_name",
	"hosted_domain",
	"managed_user_id",
}

func appendChromeProfileToOrder(profile map[string]any, dstDir string) {
	order, ok := profile["profiles_order"].([]any)
	if !ok {
		profile["profiles_order"] = []any{dstDir}
		return
	}
	for _, v := range order {
		if s, _ := v.(string); s == dstDir {
			return
		}
	}
	profile["profiles_order"] = append(order, dstDir)
}

func writeChromeLocalState(path string, root map[string]any) error {
	out, err := json.MarshalIndent(root, "", constants.JSONIndent)
	if err != nil {
		return fmt.Errorf("encode Local State: %w", err)
	}
	tmp := path + constants.ChromeLocalStateTmpSuffix
	if err := os.WriteFile(tmp, out, constants.FilePermission); err != nil {
		return fmt.Errorf("write %s: %w", tmp, err)
	}
	return replaceChromeLocalState(path, tmp)
}

func replaceChromeLocalState(path, tmp string) error {
	if !chromeProfilePathExists(path) {
		return os.Rename(tmp, path)
	}
	bak := path + constants.ChromeLocalStateBakSuffix
	if err := os.Remove(bak); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.WarnChromeProfileBakRm, bak, err)
	}
	if err := os.Rename(path, bak); err != nil {
		return cleanupChromeLocalStateTmp(tmp, err)
	}
	return finishChromeLocalStateReplace(path, tmp, bak)
}

func cleanupChromeLocalStateTmp(tmp string, cause error) error {
	_ = os.Remove(tmp)
	return fmt.Errorf("backup Local State: %w", cause)
}

func finishChromeLocalStateReplace(path, tmp, bak string) error {
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Rename(bak, path)
		_ = os.Remove(tmp)
		return fmt.Errorf("replace %s: %w", path, err)
	}
	if err := os.Remove(bak); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.WarnChromeProfileBakRm, bak, err)
	}
	return nil
}
