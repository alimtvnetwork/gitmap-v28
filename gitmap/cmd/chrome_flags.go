// Package cmd — chrome_flags.go: inspect, enable, disable, and reset
// Chrome feature flags and experiments stored in Local State.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type chromeFlagReport struct {
	EnabledCount int      `json:"enabledCount" yaml:"enabledCount"`
	Flags        []string `json:"flags" yaml:"flags"`
	UserDataDir  string   `json:"userDataDir" yaml:"userDataDir"`
}

func runChromeFlags(args []string) error {
	checkHelp(constants.SubCmdChromeFlags, args)
	if len(args) == 0 {
		return displayChromeFlags(constants.OutputTerminal)
	}
	sub := args[0]
	if sub == "--json" || sub == "-json" {
		return displayChromeFlags(constants.OutputJSON)
	}
	if sub == "--yaml" || sub == "-yaml" {
		return displayChromeFlags(constants.OutputYAML)
	}
	return dispatchFlagAction(sub, args[1:])
}

func dispatchFlagAction(action string, args []string) error {
	switch action {
	case "list", "ls":
		fmtStr := constants.OutputTerminal
		if len(args) > 0 && (args[0] == "--json" || args[0] == "-json") {
			fmtStr = constants.OutputJSON
		}
		if len(args) > 0 && (args[0] == "--yaml" || args[0] == "-yaml") {
			fmtStr = constants.OutputYAML
		}
		return displayChromeFlags(fmtStr)
	case "enable", "on", "set":
		return setChromeFlagState(args, true)
	case "disable", "off", "rm":
		return setChromeFlagState(args, false)
	case "reset", "clear":
		return resetChromeFlags()
	default:
		return apperror.NewSimple(fmt.Sprintf("unknown flags command %q (use list, enable, disable, reset)", action), "E4301")
	}
}

func displayChromeFlags(format string) error {
	flags, dir, err := readEnabledChromeFlags()
	if err != nil {
		return err
	}
	report := chromeFlagReport{EnabledCount: len(flags), Flags: flags, UserDataDir: dir}
	if format == constants.OutputJSON {
		return printJSON(report)
	}
	if format == constants.OutputYAML {
		return printYAML(report)
	}
	return printFlagsTable(flags, dir)
}

func printFlagsTable(flags []string, dir string) error {
	fmt.Printf("\n\033[1;96mChrome Experimental Feature Flags\033[0m (%s)\n\n", dir)
	if len(flags) == 0 {
		fmt.Println("  (no custom experimental flags enabled in Local State)")
		return nil
	}
	for i, f := range flags {
		fmt.Printf("  %2d. \033[1;92m●\033[0m %s\n", i+1, f)
	}
	fmt.Printf("\nTotal: %d flag(s) active\n", len(flags))
	return nil
}

func readEnabledChromeFlags() ([]string, string, error) {
	root := chromeUserDataDir()
	statePath := filepath.Join(root, constants.ChromeLocalStateFile)
	raw, err := os.ReadFile(statePath)
	if err != nil {
		return nil, root, apperror.WrapSimple(err, "read Local State")
	}
	var doc struct {
		Browser struct {
			EnabledLabs []string `json:"enabled_labs_experiments"`
		} `json:"browser"`
	}
	_ = json.Unmarshal(raw, &doc)
	return doc.Browser.EnabledLabs, root, nil
}

func setChromeFlagState(args []string, isEnable bool) error {
	if len(args) == 0 {
		return apperror.NewSimple("flag name required (e.g. enable-gpu-rasterization)", "E4302")
	}
	flagName := args[0]
	root := chromeUserDataDir()
	statePath := filepath.Join(root, constants.ChromeLocalStateFile)
	raw, err := os.ReadFile(statePath)
	if err != nil {
		return apperror.WrapSimple(err, "read Local State")
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return apperror.WrapSimple(err, "parse Local State")
	}
	return commitFlagState(statePath, doc, flagName, isEnable)
}

func commitFlagState(statePath string, doc map[string]any, flagName string, isEnable bool) error {
	browser, ok := doc["browser"].(map[string]any)
	if !ok {
		browser = map[string]any{}
		doc["browser"] = browser
	}
	existing := extractStringSlice(browser["enabled_labs_experiments"])
	updated := mutateFlagInSlice(existing, flagName, isEnable)
	browser["enabled_labs_experiments"] = updated
	newRaw, err := json.MarshalIndent(doc, "", constants.JSONIndent)
	if err != nil {
		return err
	}
	_ = os.WriteFile(statePath+".bak", rawToBytes(doc), constants.FilePermission)
	if err := os.WriteFile(statePath, newRaw, constants.FilePermission); err != nil {
		return err
	}
	action := "disabled"
	if isEnable {
		action = "enabled"
	}
	fmt.Printf("\033[1;92m✓ flag %s\033[0m  %q in Chrome Local State\n", action, flagName)
	return nil
}

func rawToBytes(doc map[string]any) []byte {
	b, _ := json.MarshalIndent(doc, "", constants.JSONIndent)
	return b
}

func extractStringSlice(val any) []string {
	slice, ok := val.([]any)
	if !ok {
		return []string{}
	}
	var list []string
	for _, item := range slice {
		if s, isStr := item.(string); isStr {
			list = append(list, s)
		}
	}
	return list
}

func mutateFlagInSlice(list []string, target string, isEnable bool) []string {
	var filtered []string
	cleanTarget := strings.Split(target, "@")[0]
	for _, item := range list {
		if !strings.HasPrefix(item, cleanTarget) {
			filtered = append(filtered, item)
		}
	}
	if !isEnable {
		return filtered
	}
	flagWithState := target
	if !strings.Contains(target, "@") {
		flagWithState = target + "@1"
	}
	return append(filtered, flagWithState)
}

func resetChromeFlags() error {
	root := chromeUserDataDir()
	statePath := filepath.Join(root, constants.ChromeLocalStateFile)
	raw, err := os.ReadFile(statePath)
	if err != nil {
		return apperror.WrapSimple(err, "read Local State")
	}
	var doc map[string]any
	_ = json.Unmarshal(raw, &doc)
	if browser, ok := doc["browser"].(map[string]any); ok {
		browser["enabled_labs_experiments"] = []string{}
	}
	newRaw, _ := json.MarshalIndent(doc, "", constants.JSONIndent)
	_ = os.WriteFile(statePath, newRaw, constants.FilePermission)
	fmt.Println("\033[1;92m✓ reset\033[0m all Chrome experimental flags in Local State")
	return nil
}
