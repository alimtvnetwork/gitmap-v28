package cmdvscode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// VSCodeProfileInfo represents an inspected VSCode profile.
type VSCodeProfileInfo struct {
	Name           string `json:"name"`
	ID             string `json:"id"`
	Location       string `json:"location"`
	HasSettings    bool   `json:"has_settings"`
	HasKeybindings bool   `json:"has_keybindings"`
	HasExtensions  bool   `json:"has_extensions"`
}

// VSCodeProfileBundle stores exportable profile configuration.
type VSCodeProfileBundle struct {
	Name        string          `json:"name"`
	ID          string          `json:"id"`
	Settings    json.RawMessage `json:"settings,omitempty"`
	Keybindings json.RawMessage `json:"keybindings,omitempty"`
	Extensions  json.RawMessage `json:"extensions,omitempty"`
}

func runVSCodeProfiles(args []string) error {
	if len(args) == 0 || isVSCodeHelp(args[0]) {
		printVSCodeProfilesUsage()
		return nil
	}
	return dispatchVSCodeProfiles(strings.ToLower(args[0]), args[1:])
}

func isVSCodeHelp(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func dispatchVSCodeProfiles(sub string, args []string) error {
	switch sub {
	case "ls", "list":
		return runProfilesList()
	case "export":
		return runProfileExport(args)
	case "export-all":
		return runProfilesExportAll(args)
	case "import":
		return runProfileImport(args)
	case "import-all":
		return runProfilesImportAll(args)
	default:
		printVSCodeProfilesUsage()
		return apperror.NewSimple("unknown vscode profiles subcommand: "+sub, "E_INVALID_PROFILE_SUBCMD")
	}
}

func runProfilesList() error {
	profiles := discoverVSCodeProfiles()
	if len(profiles) == 0 {
		fmt.Println("No VSCode profiles discovered.")
		return nil
	}
	fmt.Printf("%-18s %-16s %-12s %-12s %s\n", "NAME", "ID", "SETTINGS", "EXTENSIONS", "LOCATION")
	fmt.Println(strings.Repeat("-", 80))
	for _, p := range profiles {
		setStr := "no"
		if p.HasSettings {
			setStr = "yes"
		}
		extStr := "no"
		if p.HasExtensions {
			extStr = "yes"
		}
		fmt.Printf("%-18s %-16s %-12s %-12s %s\n", p.Name, p.ID, setStr, extStr, p.Location)
	}
	return nil
}

func runProfileExport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap vscode profiles export <name> [file.json]", "E_MISSING_ARG")
	}
	name := args[0]
	bundle, err := exportSingleProfile(name)
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(bundle, "", "  ")
	if len(args) >= 2 {
		return writeProfileExportFile(args[1], data, name)
	}
	fmt.Println(string(data))
	return nil
}

func writeProfileExportFile(destPath string, data []byte, name string) error {
	if writeErr := os.WriteFile(destPath, data, 0o644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write profile export")
	}
	fmt.Printf("✔ Exported VSCode profile %s to %s\n", name, destPath)
	return nil
}

func runProfilesExportAll(args []string) error {
	profiles := discoverVSCodeProfiles()
	var bundles []VSCodeProfileBundle
	for _, p := range profiles {
		if b, err := exportSingleProfile(p.Name); err == nil {
			bundles = append(bundles, b)
		}
	}
	data, _ := json.MarshalIndent(bundles, "", "  ")
	if len(args) >= 1 {
		return writeAllProfilesExportFile(args[0], data, len(bundles))
	}
	fmt.Println(string(data))
	return nil
}

func writeAllProfilesExportFile(destPath string, data []byte, count int) error {
	if writeErr := os.WriteFile(destPath, data, 0o644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write all profiles export")
	}
	fmt.Printf("✔ Exported %d VSCode profiles to %s\n", count, destPath)
	return nil
}

func runProfileImport(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap vscode profiles import <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read profile import")
	}
	var bundle VSCodeProfileBundle
	if unmarshalErr := json.Unmarshal(content, &bundle); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal profile bundle")
	}
	if importErr := applyProfileBundle(bundle); importErr != nil {
		return importErr
	}
	fmt.Printf("✔ Successfully imported VSCode profile: %s\n", bundle.Name)
	return nil
}

func runProfilesImportAll(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap vscode profiles import-all <file.json>", "E_MISSING_ARG")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read all profiles import")
	}
	var bundles []VSCodeProfileBundle
	if unmarshalErr := json.Unmarshal(content, &bundles); unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "unmarshal profiles bundle")
	}
	count := 0
	for _, b := range bundles {
		if applyProfileBundle(b) == nil {
			count++
		}
	}
	fmt.Printf("✔ Successfully imported %d VSCode profiles from %s\n", count, args[0])
	return nil
}

func exportSingleProfile(name string) (VSCodeProfileBundle, *apperror.AppError) {
	profiles := discoverVSCodeProfiles()
	var target *VSCodeProfileInfo
	for _, p := range profiles {
		if strings.EqualFold(p.Name, name) || strings.EqualFold(p.ID, name) {
			target = &p
			break
		}
	}
	if target == nil {
		return VSCodeProfileBundle{}, apperror.NewSimple("profile not found: "+name, "E_PROFILE_NOT_FOUND")
	}
	bundle := VSCodeProfileBundle{Name: target.Name, ID: target.ID}
	bundle.Settings = readJSONFile(filepath.Join(target.Location, "settings.json"))
	bundle.Keybindings = readJSONFile(filepath.Join(target.Location, "keybindings.json"))
	bundle.Extensions = readJSONFile(filepath.Join(target.Location, "extensions.json"))
	return bundle, nil
}

func applyProfileBundle(bundle VSCodeProfileBundle) *apperror.AppError {
	targetDir := resolveProfileDirectory(bundle)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return apperror.WrapSimple(err, "create profile directory")
	}
	writeIfNonEmpty(filepath.Join(targetDir, "settings.json"), bundle.Settings)
	writeIfNonEmpty(filepath.Join(targetDir, "keybindings.json"), bundle.Keybindings)
	writeIfNonEmpty(filepath.Join(targetDir, "extensions.json"), bundle.Extensions)
	return nil
}

func resolveProfileDirectory(bundle VSCodeProfileBundle) string {
	base := resolveVSCodeUserDir()
	if strings.EqualFold(bundle.Name, "default") || strings.EqualFold(bundle.ID, "default") {
		return base
	}
	id := bundle.ID
	if id == "" {
		id = strings.ToLower(bundle.Name)
	}
	return filepath.Join(base, "profiles", id)
}

func writeIfNonEmpty(path string, data json.RawMessage) {
	if len(data) > 0 {
		_ = os.WriteFile(path, data, 0o644)
	}
}

func readJSONFile(path string) json.RawMessage {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return json.RawMessage(content)
}

func discoverVSCodeProfiles() []VSCodeProfileInfo {
	userDir := resolveVSCodeUserDir()
	var profiles []VSCodeProfileInfo
	if defaultProfile := inspectProfileDir("Default", "default", userDir); defaultProfile.Location != "" {
		profiles = append(profiles, defaultProfile)
	}
	profilesDir := filepath.Join(userDir, "profiles")
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return profiles
	}
	for _, e := range entries {
		if e.IsDir() {
			pDir := filepath.Join(profilesDir, e.Name())
			profiles = append(profiles, inspectProfileDir(e.Name(), e.Name(), pDir))
		}
	}
	return profiles
}

func inspectProfileDir(name, id, dir string) VSCodeProfileInfo {
	if _, err := os.Stat(dir); err != nil {
		return VSCodeProfileInfo{}
	}
	return VSCodeProfileInfo{
		Name:           name,
		ID:             id,
		Location:       dir,
		HasSettings:    hasRegularFile(filepath.Join(dir, "settings.json")),
		HasKeybindings: hasRegularFile(filepath.Join(dir, "keybindings.json")),
		HasExtensions:  hasRegularFile(filepath.Join(dir, "extensions.json")),
	}
}

func hasRegularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func resolveVSCodeUserDir() string {
	if runtime.GOOS == constants.OSWindows {
		return filepath.Join(os.Getenv("APPDATA"), "Code", "User")
	}
	if runtime.GOOS == constants.OSDarwin {
		home := os.Getenv("HOME")
		return filepath.Join(home, "Library", "Application Support", "Code", "User")
	}
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", "Code", "User")
}

func printVSCodeProfilesUsage() {
	fmt.Println("Usage: gitmap vscode profiles <subcommand> [arguments]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  ls, list                   List all installed VSCode user profiles")
	fmt.Println("  export <name> [file.json]  Export a profile configuration")
	fmt.Println("  export-all [file.json]     Export all profiles to a JSON bundle")
	fmt.Println("  import <file.json>         Import a profile configuration")
	fmt.Println("  import-all <file.json>     Import multiple profiles from a bundle")
}
