// Package cmd - install_add.go handles interactive and scripted installer creation.
package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// InstallAddFlags encapsulates CLI and interactive arguments for installer creation.
type InstallAddFlags struct {
	Name         string
	Version      string
	Description  string
	WinScript    string
	UnixScript   string
	UbuntuScript string
	AddLater     bool
	Yes          bool
}

// isInstallAddCommand reports if args invoke installer add or create.
func isInstallAddCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	subCmd := strings.ToLower(strings.TrimSpace(args[0]))

	return subCmd == "add" || subCmd == "create"
}

func runInstallAdd(args []string) error {
	flags, errParse := parseInstallAddFlags(args)
	if errParse != nil {
		return errParse
	}

	if errPrompt := promptIfInteractive(flags); errPrompt != nil {
		return errPrompt
	}

	return executeInstallAdd(flags)
}

func promptIfInteractive(flags *InstallAddFlags) error {
	if !isInteractivePromptNeeded(flags) {
		return nil
	}

	return promptInteractiveInstallAdd(flags)
}

func isInteractivePromptNeeded(flags *InstallAddFlags) bool {
	if flags.Yes || !isTerminalInput() {
		return false
	}

	return true
}

func parseInstallAddFlags(args []string) (*InstallAddFlags, error) {
	fs := flag.NewFlagSet("install-add", flag.ContinueOnError)
	flags := &InstallAddFlags{}
	fs.StringVar(&flags.Description, "desc", "", "Installer description")
	fs.StringVar(&flags.Description, "d", "", "Installer description shorthand")
	fs.StringVar(&flags.WinScript, "win", "", "Windows install command/script")
	fs.StringVar(&flags.WinScript, "w", "", "Windows shorthand")
	fs.StringVar(&flags.UnixScript, "unix", "", "Unix install command/script")
	fs.StringVar(&flags.UnixScript, "u", "", "Unix shorthand")
	fs.StringVar(&flags.UbuntuScript, "ubuntu", "", "Ubuntu install command/script")
	fs.BoolVar(&flags.Yes, "yes", false, "Skip interactive prompt")
	fs.BoolVar(&flags.Yes, "y", false, "Skip interactive shorthand")
	fs.StringVar(&flags.Version, "version", "", "Installer version")
	fs.StringVar(&flags.Version, "v", "", "Version shorthand")

	return parseInstallAddPositional(fs, args, flags)
}

func parseInstallAddPositional(fs *flag.FlagSet, args []string, flags *InstallAddFlags) (*InstallAddFlags, error) {
	flagArgs, positional := separateFlagAndPositionalArgs(args)
	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseInstallAddFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"

		return nil, appErr
	}

	if len(positional) > 0 {
		flags.Name = strings.TrimSpace(positional[0])
	}

	if len(positional) > 1 && flags.Version == "" {
		flags.Version = strings.TrimSpace(positional[1])
	}

	if flags.Version == "" {
		flags.Version = "v1.0.0"
	}

	return flags, nil
}

func promptInteractiveInstallAdd(flags *InstallAddFlags) error {
	printInteractiveAddHeader(flags.Name, flags.Version)
	scanner := bufio.NewScanner(os.Stdin)
	if flags.Description == "" {
		flags.Description = promptAddLine(scanner, "Installer description: ")
	}

	if flags.WinScript == "" {
		flags.WinScript = promptAddLine(scanner, "Windows install script (PowerShell) [empty to skip]: ")
	}

	if flags.UnixScript == "" {
		flags.UnixScript = promptAddLine(scanner, "Unix install script (sh/bash) [empty to skip]: ")
	}

	if flags.UbuntuScript == "" {
		flags.UbuntuScript = promptAddLine(scanner, "Ubuntu install script (apt/bash) [empty to skip]: ")
	}

	askLater := promptAddLine(scanner, "Add/edit Unix or Ubuntu instructions later? (y/n) [n]: ")
	flags.AddLater = strings.EqualFold(askLater, "y") || strings.EqualFold(askLater, "yes")

	return scanner.Err()
}

func printInteractiveAddHeader(name, ver string) {
	header := fmt.Sprintf("\n  %s⚙ Custom Installer Setup: %q (%s)%s\n", constants.ColorCyan, name, ver, constants.ColorReset)
	fmt.Println(header)
}

func promptAddLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print("  " + prompt)
	if !scanner.Scan() {
		return ""
	}

	return strings.TrimSpace(scanner.Text())
}

func executeInstallAdd(flags *InstallAddFlags) error {
	if strings.TrimSpace(flags.Name) == "" {
		appErr := apperror.NewValidationError("installer name is required")
		appErr.Code = "E_INSTALLER_INVALID_INPUT"

		return appErr
	}

	db, errDB := openAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}

	defer db.Close()

	script := buildInstallerScriptFromFlags(flags)

	return persistInstallerScript(db, script)
}

func buildInstallerScriptFromFlags(flags *InstallAddFlags) *model.InstallerScript {
	slug := slugify(flags.Name)
	scriptsMap := buildOSScriptsMap(flags)
	targetOS := resolveTargetOSFromFlags(flags)
	instructions := marshalMultiOSInstructions(flags)

	return &model.InstallerScript{
		Name:         flags.Name,
		Slug:         slug,
		Description:  flags.Description,
		TargetOS:     targetOS,
		Version:      flags.Version,
		Instructions: instructions,
		Scripts:      scriptsMap,
	}
}

func buildOSScriptsMap(flags *InstallAddFlags) map[string]model.OSScript {
	out := make(map[string]model.OSScript)
	if flags.WinScript != "" {
		out["win"] = model.OSScript{Runtime: "powershell", Instructions: flags.WinScript}
	}

	if flags.UnixScript != "" {
		out["unix"] = model.OSScript{Runtime: "bash", Instructions: flags.UnixScript}
	}

	if flags.UbuntuScript != "" {
		out["ubuntu"] = model.OSScript{Runtime: "bash", Instructions: flags.UbuntuScript}
	}

	return out
}

func resolveTargetOSFromFlags(flags *InstallAddFlags) string {
	hasWin := flags.WinScript != ""
	hasUnix := flags.UnixScript != "" || flags.UbuntuScript != ""
	if hasWin && hasUnix {
		return "all"
	}

	if hasWin {
		return "win"
	}

	if flags.UbuntuScript != "" {
		return "ubuntu"
	}

	if flags.UnixScript != "" {
		return "unix"
	}

	return "all"
}

func marshalMultiOSInstructions(flags *InstallAddFlags) string {
	payload := map[string]string{}
	if flags.WinScript != "" {
		payload["win"] = flags.WinScript
	}

	if flags.UnixScript != "" {
		payload["unix"] = flags.UnixScript
	}

	if flags.UbuntuScript != "" {
		payload["ubuntu"] = flags.UbuntuScript
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	return string(data)
}

func persistInstallerScript(db *store.DB, script *model.InstallerScript) error {
	existing, _ := db.GetInstallerBySlug(script.Slug)
	if existing != nil {
		script.ID = existing.ID

		return updateExistingInstaller(db, script)
	}

	return createNewInstaller(db, script)
}

func updateExistingInstaller(db *store.DB, script *model.InstallerScript) error {
	if errUpdate := db.UpdateInstaller(script); errUpdate != nil {
		return errUpdate
	}

	printInstallAddSuccess(script, true)

	return nil
}

func createNewInstaller(db *store.DB, script *model.InstallerScript) error {
	if errCreate := db.CreateInstaller(script); errCreate != nil {
		return errCreate
	}

	printInstallAddSuccess(script, false)

	return nil
}

func printInstallAddSuccess(script *model.InstallerScript, isUpdate bool) {
	action := "created"
	if isUpdate {
		action = "updated"
	}

	msg := fmt.Sprintf("\n  %s✔ Custom installer %q (%s) %s successfully!%s\n",
		constants.ColorGreen, script.Name, script.Version, action, constants.ColorReset)
	fmt.Println(msg)
	fmt.Printf("    Slug:       %s\n", script.Slug)
	fmt.Printf("    Target OS:  %s\n", script.TargetOS)
	fmt.Printf("    Available:  Run 'gitmap install ls' to view in Custom Tools section.\n\n")
}
