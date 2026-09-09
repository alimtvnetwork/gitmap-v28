package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runImportConfig handles `gitmap import-config <tool|all> [path]`.
func runImportConfig(args []string) error {
	if len(args) == 0 || isConfigHelpRequested(args[0]) {
		printImportConfigUsage()

		return nil
	}
	canonicalTool, inputTool, targetArg := parseImportArgs(args)
	if canonicalTool == "all" {

		return importAllToolConfigs(targetArg)
	}

	return importSingleToolConfig(canonicalTool, inputTool, targetArg)
}

// parseImportArgs extracts tool and path arguments for import.
func parseImportArgs(args []string) (string, string, string) {
	inputTool := args[0]
	canonicalTool := normalizeConfigTool(inputTool)
	targetArg := ""
	if len(args) > 1 {
		targetArg = args[1]
	}

	return canonicalTool, inputTool, targetArg
}

// printImportConfigUsage prints standard import command usage.
func printImportConfigUsage() {
	fmt.Println("Usage: gitmap import-config <tool|all> [path]")
	fmt.Println()
	fmt.Println("Import tool configuration files from a JSON bundle.")
	fmt.Println()
	fmt.Println("Supported tools: vscode, qtorrent (qbittorrent), utorrent (uttorrent), all")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap import-config qtorrent")
	fmt.Println("  gitmap import-config vscode")
	fmt.Println("  gitmap import-config uttorrent")
	fmt.Println("  gitmap import-config all ./backup-configs")
}

// locateExistingConfigBundle finds existing bundle file with alias fallbacks.
func locateExistingConfigBundle(targetArg, inputTool, canonicalTool string) string {
	defaultName := resolveDefaultConfigFileName(inputTool, canonicalTool)
	directPath := resolveConfigFilePath(targetArg, defaultName)
	if _, err := os.Stat(directPath); err == nil {

		return directPath
	}
	aliasName := canonicalTool + ".json"
	aliasPath := resolveConfigFilePath(targetArg, aliasName)
	if _, err := os.Stat(aliasPath); err == nil {

		return aliasPath
	}

	return directPath
}

// importSingleToolConfig imports a single tool bundle into its target config directory.
func importSingleToolConfig(canonicalTool, inputTool, targetArg string) error {
	canonicalTool = normalizeConfigTool(canonicalTool)
	sourcePath := locateExistingConfigBundle(targetArg, inputTool, canonicalTool)
	bundle, err := loadConfigBundle(sourcePath)
	if err != nil {
		cliexit.Reportf("import-config", "read", sourcePath, err)

		return err
	}

	return restoreAndReportToolConfig(canonicalTool, sourcePath, bundle)
}

// restoreAndReportToolConfig restores the bundle files and logs success.
func restoreAndReportToolConfig(canonicalTool, sourcePath string, bundle *ConfigBundle) error {
	destDir := resolveToolConfigDir(canonicalTool)
	if destDir == "" {
		err := fmt.Errorf("cannot resolve config directory for %s", canonicalTool)
		cliexit.Reportf("import-config", "resolve", canonicalTool, err)

		return err
	}
	count, restoreErr := restoreBundleFiles(bundle, destDir)
	if restoreErr != nil {
		cliexit.Reportf("import-config", "restore", destDir, restoreErr)

		return restoreErr
	}
	logImportSuccess(canonicalTool, sourcePath, destDir, count, bundle.Extensions)

	return nil
}

// logImportSuccess prints success message and writes extensions reference if present.
func logImportSuccess(tool, sourcePath, destDir string, count int, extensions []string) {
	writeExtensionsReference(destDir, extensions)
	fmt.Printf("%s Imported %s configuration from %s (%d file(s) restored into %s)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, tool, sourcePath, count, destDir)
}

// importAllToolConfigs imports all discoverable *.json bundles from a directory.
func importAllToolConfigs(targetFolder string) error {
	folder := resolveExportFolder(targetFolder)
	entries, err := os.ReadDir(folder)
	if err != nil {

		return apperror.Wrap(err, "importAllToolConfigs.ReadDir", map[string]any{"folder": folder})
	}
	total := 0
	for _, entry := range entries {
		if isImportableJSON(entry) {
			imported := importEntryIfValid(folder, entry.Name())
			total += imported
		}
	}
	fmt.Printf("%s Batch import complete: %d tool configuration(s) processed from %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, total, folder)

	return nil
}

// isImportableJSON returns true if entry is a JSON file.
func isImportableJSON(entry os.DirEntry) bool {
	return !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json")
}

// importEntryIfValid attempts to import a single bundle and returns 1 if successful.
func importEntryIfValid(folder, fileName string) int {
	fullPath := filepath.Join(folder, fileName)
	bundle, loadErr := loadConfigBundle(fullPath)
	if loadErr != nil || bundle.Tool == "" {

		return 0
	}
	if err := importSingleToolConfig(bundle.Tool, bundle.Tool, fullPath); err != nil {

		return 0
	}

	return 1
}
