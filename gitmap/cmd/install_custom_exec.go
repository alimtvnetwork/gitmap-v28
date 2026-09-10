// Package cmd — install_custom_exec.go: execution engine for user-defined custom installers.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func findCustomInstaller(slug string) *model.InstallerScript {
	cleanSlug := strings.ToLower(strings.TrimSpace(slug))
	if cleanSlug == "" {
		return nil
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return nil
	}
	defer db.Close()

	_ = db.MigrateInstallers()
	script, errGet := db.GetInstallerBySlug(cleanSlug)
	if errGet != nil || script == nil {
		return nil
	}

	return script
}

func hasCustomInstaller(slug string) bool {
	return findCustomInstaller(slug) != nil
}

func executeCustomInstaller(script *model.InstallerScript, opts installOptions) error {
	printCustomInstallerStartBanner(script)

	if opts.DryRun {
		printCustomInstallerDryRun(script)

		return nil
	}

	scriptText, runtimeName, errResolve := resolveOSScriptForPlatform(script)
	if errResolve != nil {
		return handleCustomInstallerResolveError(script, errResolve)
	}

	if errExec := runCustomInstallerProcess(scriptText, runtimeName); errExec != nil {
		return handleCustomInstallerExecError(script, errExec)
	}

	recordCustomInstallSuccess(script)

	return nil
}

func printCustomInstallerStartBanner(script *model.InstallerScript) {
	fmt.Printf("\n  %s⚙ Running Custom Installer: %q (%s)...%s\n",
		constants.ColorCyan, script.Name, script.Version, constants.ColorReset)
	fmt.Printf("    Target OS:  %s\n", script.TargetOS)
	if script.Description != "" {
		fmt.Printf("    Desc:       %s\n", script.Description)
	}
	fmt.Println()
}

func printCustomInstallerDryRun(script *model.InstallerScript) {
	fmt.Printf("  %s(dry run) Would execute custom installer for %q on OS %s%s\n\n",
		constants.ColorYellow, script.Name, runtime.GOOS, constants.ColorReset)
}

func parseInstructionsMap(raw string) map[string]string {
	out := make(map[string]string)
	if err := json.Unmarshal([]byte(raw), &out); err == nil && len(out) > 0 {
		return out
	}

	trimmed := strings.TrimSpace(raw)
	if trimmed != "" {
		out["all"] = trimmed
	}

	return out
}

func resolveOSScriptForPlatform(script *model.InstallerScript) (string, string, error) {
	instructionsMap := parseInstructionsMap(script.Instructions)

	switch runtime.GOOS {
	case "windows":
		return resolveWindowsScript(instructionsMap, script.TargetOS)
	case "linux":
		return resolveLinuxScript(instructionsMap, script.TargetOS)
	case "darwin":
		return resolveDarwinScript(instructionsMap, script.TargetOS)
	default:
		return resolveGenericUnixScript(instructionsMap, script.TargetOS)
	}
}

func resolveWindowsScript(instructions map[string]string, targetOS string) (string, string, error) {
	if script, exists := instructions["win"]; exists && strings.TrimSpace(script) != "" {
		return script, "powershell", nil
	}

	if script, exists := instructions["all"]; exists && strings.TrimSpace(script) != "" {
		return script, "powershell", nil
	}

	return "", "", apperror.NewValidationError("no Windows instructions configured for this installer")
}

func resolveLinuxScript(instructions map[string]string, targetOS string) (string, string, error) {
	if uScript := resolveUbuntuScript(instructions); uScript != "" {
		return uScript, "bash", nil
	}

	if script, exists := instructions["unix"]; exists && strings.TrimSpace(script) != "" {
		return script, "bash", nil
	}

	if script, exists := instructions["all"]; exists && strings.TrimSpace(script) != "" {
		return script, "bash", nil
	}

	return "", "", apperror.NewValidationError("no Linux/Ubuntu instructions configured for this installer")
}

func resolveUbuntuScript(instructions map[string]string) string {
	if !isUbuntuHost() {
		return ""
	}

	return strings.TrimSpace(instructions["ubuntu"])
}


func isUbuntuHost() bool {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}

	content := strings.ToLower(string(data))

	return strings.Contains(content, "ubuntu") || strings.Contains(content, "debian")
}

func resolveDarwinScript(instructions map[string]string, targetOS string) (string, string, error) {
	if script, exists := instructions["unix"]; exists && strings.TrimSpace(script) != "" {
		return script, "bash", nil
	}

	if script, exists := instructions["all"]; exists && strings.TrimSpace(script) != "" {
		return script, "bash", nil
	}

	return "", "", apperror.NewValidationError("no macOS/Unix instructions configured for this installer")
}

func resolveGenericUnixScript(instructions map[string]string, targetOS string) (string, string, error) {
	if script, exists := instructions["unix"]; exists && strings.TrimSpace(script) != "" {
		return script, "sh", nil
	}

	if script, exists := instructions["all"]; exists && strings.TrimSpace(script) != "" {
		return script, "sh", nil
	}

	return "", "", apperror.NewValidationError("no instructions configured for this operating system")
}

func runCustomInstallerProcess(scriptText string, runtimeName string) error {
	cmd := buildInstallerCommand(scriptText, runtimeName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("execute %s script", runtimeName))
	}

	return nil
}

func buildInstallerCommand(scriptText string, runtimeName string) *exec.Cmd {
	if runtimeName == "powershell" {
		return exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", scriptText)
	}

	if runtimeName == "bash" {
		return exec.Command("bash", "-c", scriptText)
	}

	return exec.Command("sh", "-c", scriptText)
}

func recordCustomInstallSuccess(script *model.InstallerScript) {
	splitDB, errDB := store.OpenInstallationSplitDB()
	if errDB == nil {
		defer splitDB.Close()
		_ = splitDB.SaveInstalledTool(script.Slug, script.Version, "custom")
	}

	fmt.Printf("\n  %s✔ Custom installer %q installed successfully!%s\n",
		constants.ColorGreen, script.Name, constants.ColorReset)
	fmt.Println("    Status:  Recorded as installed in installation database.")
	fmt.Printf("    Listing: Run 'gitmap install ls' to view updated status.\n\n")
}

func handleCustomInstallerResolveError(script *model.InstallerScript, err error) error {
	fmt.Fprintf(os.Stderr, "\n  %s✖ Custom installer %q cannot run on OS %q: %v%s\n",
		constants.ColorRed, script.Name, runtime.GOOS, err, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "    Target OS configured: %s\n", script.TargetOS)
	fmt.Fprintf(os.Stderr, "    Run 'gitmap install add %s' to configure instructions for this OS.\n\n", script.Slug)

	return err
}

func handleCustomInstallerExecError(script *model.InstallerScript, err error) error {
	fmt.Fprintf(os.Stderr, "\n  %s✖ Custom installer %q failed during execution: %v%s\n\n",
		constants.ColorRed, script.Name, err, constants.ColorReset)

	return err
}
