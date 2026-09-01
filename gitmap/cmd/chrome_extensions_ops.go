// Package cmd — chrome_extensions_ops.go: enable, disable, and install operations
// for Chrome extensions.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runChromeExtensionEnable(args []string) error {
	checkHelp(constants.SubCmdChromeExtEnable, args)
	return mutateExtensionsState(args, true, false)
}

func runChromeExtensionDisable(args []string) error {
	checkHelp(constants.SubCmdChromeExtDisable, args)
	return mutateExtensionsState(args, false, false)
}

func runChromeExtensionDisableAll(args []string) error {
	checkHelp(constants.SubCmdChromeExtDisableAll, args)
	return mutateExtensionsState(args, false, true)
}

func mutateExtensionsState(args []string, isEnable, isAll bool) error {
	opts := parseExtensionFilterArgs(args)
	if len(args) == 0 && !isAll {
		return apperror.NewSimple("extension pattern or ID required", "E4202")
	}
	pattern := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		pattern = args[0]
	}
	profiles := resolveTargetProfiles(opts.Profile, opts.IsAll)
	for _, prof := range profiles {
		if err := applyExtensionStateChange(prof, pattern, isEnable, isAll); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[1;91m✗ %s:\033[0m %v\n", prof, err)
		}
	}
	return nil
}

func applyExtensionStateChange(profName, pattern string, isEnable, isAll bool) error {
	srcPath, hasDir := resolveChromeProfileDir(profName)
	if !hasDir {
		return apperror.NewSimple(fmt.Sprintf("profile %s not found", profName), "E4203")
	}
	exts, _ := scanExtensionsForProfile(profName)
	matched := filterExtensionsByPattern(exts, pattern, isAll)
	if len(matched) == 0 {
		fmt.Printf("No matching extensions found in profile %q\n", profName)
		return nil
	}
	return updateProfilePreferencesExtensionState(srcPath, profName, matched, isEnable)
}

func filterExtensionsByPattern(exts []chromeExtensionInfo, pattern string, isAll bool) []chromeExtensionInfo {
	if isAll || pattern == "*" || pattern == "all" {
		return exts
	}
	var matched []chromeExtensionInfo
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(pattern))
	for _, e := range exts {
		isMatch := e.ID == pattern || (err == nil && (re.MatchString(e.Name) || re.MatchString(e.ID)))
		if isMatch {
			matched = append(matched, e)
		}
	}
	return matched
}

func updateProfilePreferencesExtensionState(srcPath, profName string, matched []chromeExtensionInfo, isEnable bool) error {
	prefPath := filepath.Join(srcPath, "Preferences")
	raw, err := os.ReadFile(prefPath)
	if err != nil {
		return apperror.WrapSimple(err, "read Preferences")
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return apperror.WrapSimple(err, "parse Preferences")
	}
	settings := ensureSettingsMap(doc)
	targetState := 0
	actionStr := "disabled"
	if isEnable {
		targetState = 1
		actionStr = "enabled"
	}
	return commitExtensionStateChanges(prefPath, profName, doc, settings, matched, targetState, actionStr)
}

func ensureSettingsMap(doc map[string]any) map[string]any {
	extVal, ok := doc["extensions"].(map[string]any)
	if !ok {
		extVal = map[string]any{}
		doc["extensions"] = extVal
	}
	setVal, ok := extVal["settings"].(map[string]any)
	if !ok {
		setVal = map[string]any{}
		extVal["settings"] = setVal
	}
	return setVal
}

func commitExtensionStateChanges(prefPath, profName string, doc map[string]any, settings map[string]any, matched []chromeExtensionInfo, targetState int, actionStr string) error {
	for _, e := range matched {
		entry, ok := settings[e.ID].(map[string]any)
		if !ok {
			entry = map[string]any{}
			settings[e.ID] = entry
		}
		entry["state"] = targetState
		fmt.Printf("  \033[1;92m✓\033[0m extension %q (%s) %s in profile %q\n", e.Name, e.ID, actionStr, profName)
	}
	newRaw, err := json.MarshalIndent(doc, "", constants.JSONIndent)
	if err != nil {
		return err
	}
	_ = os.WriteFile(prefPath+".bak", docToBytes(doc), constants.FilePermission)
	return os.WriteFile(prefPath, newRaw, constants.FilePermission)
}

func docToBytes(doc map[string]any) []byte {
	b, _ := json.MarshalIndent(doc, "", constants.JSONIndent)
	return b
}

func runChromeExtensionInstall(args []string) error {
	checkHelp(constants.SubCmdChromeExtInstall, args)
	if len(args) == 0 {
		return apperror.NewSimple("extension folder path or .crx file required", "E4204")
	}
	extPath := args[0]
	profName := "Default"
	for _, a := range args[1:] {
		if strings.HasPrefix(a, "--profile=") {
			profName = strings.TrimPrefix(a, "--profile=")
		}
	}
	return installExtensionToProfile(extPath, profName)
}

func installExtensionToProfile(extPath, profName string) error {
	manifestPath := filepath.Join(extPath, "manifest.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("invalid extension directory %s (missing manifest.json)", extPath))
	}
	var m struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	_ = json.Unmarshal(raw, &m)
	srcProfile, hasDir := resolveChromeProfileDir(profName)
	if !hasDir {
		return apperror.NewSimple(fmt.Sprintf("profile %s not found", profName), "E4205")
	}
	extID := generateExtensionID(extPath, m.Name)
	dstDir := filepath.Join(srcProfile, "Extensions", extID, m.Version)
	if _, copyErr := copyEntry(extPath, dstDir); copyErr != nil {
		return copyErr
	}
	fmt.Printf("\033[1;92m✓ installed\033[0m extension %q (v%s, id: %s) into profile %q\n", m.Name, m.Version, extID, profName)
	return nil
}

func generateExtensionID(path, name string) string {
	base := filepath.Base(path)
	clean := strings.ToLower(strings.ReplaceAll(base, "-", ""))
	if len(clean) >= 32 {
		return clean[:32]
	}
	return clean + strings.Repeat("a", 32-len(clean))
}
