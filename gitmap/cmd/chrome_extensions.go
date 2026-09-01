// Package cmd — chrome_extensions.go: manages Chrome extensions/plugins,
// including listing, enabling, disabling by pattern or in bulk, and local injection.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"gopkg.in/yaml.v3"
)

type chromeExtensionInfo struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	IsEnabled   bool   `json:"isEnabled" yaml:"isEnabled"`
	Description string `json:"description" yaml:"description"`
	Profile     string `json:"profile" yaml:"profile"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
}

type extensionFilterOpts struct {
	Profile string
	Format  string
	Pattern string
	IsAll   bool
}

func runChromeExtensions(args []string) error {
	checkHelp(constants.SubCmdChromeExtensions, args)
	opts := parseExtensionFilterArgs(args)
	profiles := resolveTargetProfiles(opts.Profile, opts.IsAll)
	var allExts []chromeExtensionInfo
	for _, prof := range profiles {
		exts, _ := scanExtensionsForProfile(prof)
		allExts = append(allExts, exts...)
	}
	return renderExtensionsOutput(allExts, opts.Format)
}

func parseExtensionFilterArgs(args []string) extensionFilterOpts {
	var opts extensionFilterOpts
	opts.Format = constants.OutputTerminal
	for _, a := range args {
		if a == "--json" || a == "-json" {
			opts.Format = constants.OutputJSON
		} else if a == "--yaml" || a == "-yaml" {
			opts.Format = constants.OutputYAML
		} else if a == "--all" || a == "-a" {
			opts.IsAll = true
		} else if strings.HasPrefix(a, "--profile=") {
			opts.Profile = strings.TrimPrefix(a, "--profile=")
		} else if !strings.HasPrefix(a, "-") && opts.Profile == "" {
			opts.Profile = a
		}
	}
	return opts
}

func resolveTargetProfiles(explicit string, isAll bool) []string {
	if isAll || explicit == "all" {
		return availableChromeProfileNames()
	}
	if explicit != "" {
		return []string{explicit}
	}
	names := availableChromeProfileNames()
	if len(names) > 0 {
		return []string{names[0]}
	}
	return []string{constants.ChromeDefaultProfileDir}
}

func scanExtensionsForProfile(profName string) ([]chromeExtensionInfo, error) {
	srcPath, hasDir := resolveChromeProfileDir(profName)
	if !hasDir {
		return nil, apperror.NewSimple(fmt.Sprintf("profile %s not found", profName), "E4201")
	}
	prefPath := filepath.Join(srcPath, "Preferences")
	extSettings := readExtensionSettingsMap(prefPath)
	extBaseDir := filepath.Join(srcPath, "Extensions")
	return collectExtensionDetails(extBaseDir, profName, extSettings), nil
}

func readExtensionSettingsMap(prefPath string) map[string]map[string]any {
	raw, err := os.ReadFile(prefPath)
	if err != nil {
		return map[string]map[string]any{}
	}
	var doc struct {
		Extensions struct {
			Settings map[string]map[string]any `json:"settings"`
		} `json:"extensions"`
	}
	_ = json.Unmarshal(raw, &doc)
	if doc.Extensions.Settings == nil {
		return map[string]map[string]any{}
	}
	return doc.Extensions.Settings
}

func collectExtensionDetails(extBaseDir, profName string, settings map[string]map[string]any) []chromeExtensionInfo {
	entries, err := os.ReadDir(extBaseDir)
	if err != nil {
		return collectSettingsOnlyExtensions(profName, settings)
	}
	var list []chromeExtensionInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		info := inspectExtensionFolder(extBaseDir, e.Name(), profName, settings[e.Name()])
		list = append(list, info)
	}
	return list
}

func inspectExtensionFolder(base, extID, profName string, setting map[string]any) chromeExtensionInfo {
	info := chromeExtensionInfo{ID: extID, Profile: profName, IsEnabled: isExtensionEnabled(setting)}
	verEntries, err := os.ReadDir(filepath.Join(base, extID))
	if err != nil || len(verEntries) == 0 {
		info.Name = extID
		return info
	}
	latestVer := verEntries[len(verEntries)-1].Name()
	info.Version = latestVer
	info.Path = filepath.Join(base, extID, latestVer)
	manifestPath := filepath.Join(info.Path, "manifest.json")
	populateManifestInfo(&info, manifestPath)
	return info
}

func populateManifestInfo(info *chromeExtensionInfo, manifestPath string) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		info.Name = info.ID
		return
	}
	var m struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
	}
	_ = json.Unmarshal(raw, &m)
	info.Name = m.Name
	if info.Name == "" || strings.HasPrefix(info.Name, "__MSG_") {
		info.Name = info.ID
	}
	info.Description = m.Description
	if m.Version != "" {
		info.Version = m.Version
	}
}

func isExtensionEnabled(setting map[string]any) bool {
	if setting == nil {
		return true
	}
	stateVal, hasState := setting["state"]
	if !hasState {
		return true
	}
	flt, ok := stateVal.(float64)
	if !ok {
		return true
	}
	return int(flt) == 1
}

func collectSettingsOnlyExtensions(profName string, settings map[string]map[string]any) []chromeExtensionInfo {
	var list []chromeExtensionInfo
	for id, s := range settings {
		list = append(list, chromeExtensionInfo{
			ID:        id,
			Name:      id,
			Profile:   profName,
			IsEnabled: isExtensionEnabled(s),
		})
	}
	return list
}

func renderExtensionsOutput(exts []chromeExtensionInfo, format string) error {
	if format == constants.OutputJSON {
		return printJSON(exts)
	}
	if format == constants.OutputYAML {
		return printYAML(exts)
	}
	return printExtensionsTable(exts)
}

func printExtensionsTable(exts []chromeExtensionInfo) error {
	if len(exts) == 0 {
		fmt.Println("No extensions found in the selected profile(s).")
		return nil
	}
	fmt.Printf("\n%-34s %-12s %-10s %-14s %s\n", "NAME", "VERSION", "STATUS", "PROFILE", "ID")
	fmt.Println(strings.Repeat("-", 100))
	for _, e := range exts {
		status := "\033[1;92menabled\033[0m"
		if !e.IsEnabled {
			status = "\033[1;90mdisabled\033[0m"
		}
		dispName := e.Name
		if len(dispName) > 32 {
			dispName = dispName[:29] + "..."
		}
		fmt.Printf("%-34s %-12s %-10s %-14s %s\n", dispName, e.Version, status, e.Profile, e.ID)
	}
	fmt.Printf("\nTotal: %d extension(s)\n", len(exts))
	return nil
}

func printYAML(data any) error {
	raw, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(raw))
	return nil
}
