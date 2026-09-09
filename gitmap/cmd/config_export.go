package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runExportConfig handles `gitmap export-config <tool> [path]`.
func runExportConfig(args []string) error {
	if len(args) == 0 || isConfigHelpRequested(args[0]) {
		printExportConfigUsage()

		return nil
	}
	canonicalTool, inputTool, targetArg := parseExportArgs(args)
	if canonicalTool == "all" {

		return exportAllToolConfigs(targetArg)
	}

	return exportSingleToolConfig(canonicalTool, inputTool, targetArg)
}

// parseExportArgs extracts tool and target arguments.
func parseExportArgs(args []string) (string, string, string) {
	inputTool := args[0]
	canonicalTool := normalizeConfigTool(inputTool)
	targetArg := ""
	if len(args) > 1 {
		targetArg = args[1]
	}

	return canonicalTool, inputTool, targetArg
}

// isConfigHelpRequested checks if arg indicates a request for help text.
func isConfigHelpRequested(arg string) bool {
	lower := strings.ToLower(strings.TrimSpace(arg))

	return lower == "help" || lower == "--help" || lower == "-h"
}

// printExportConfigUsage prints standard command usage.
func printExportConfigUsage() {
	fmt.Println("Usage: gitmap export-config <tool|all> [path]")
	fmt.Println()
	fmt.Println("Export tool configuration files into a portable JSON bundle.")
	fmt.Println()
	fmt.Println("Supported tools: vscode, qtorrent (qbittorrent), utorrent (uttorrent), all")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap export-config qtorrent")
	fmt.Println("  gitmap export-config vscode")
	fmt.Println("  gitmap export-config uttorrent")
	fmt.Println("  gitmap export-config all ./backup-configs")
}

// exportSingleToolConfig exports a single tool's configuration to a JSON file.
func exportSingleToolConfig(canonicalTool, inputTool, targetArg string) error {
	defaultName := resolveDefaultConfigFileName(inputTool, canonicalTool)
	finalPath := resolveConfigFilePath(targetArg, defaultName)
	srcDir := resolveToolConfigDir(canonicalTool)
	files, extensions := resolveExportFilesAndExtensions(canonicalTool, srcDir)
	bundle := buildConfigBundle(canonicalTool, srcDir, files, extensions)
	err := writeBundleJSON(bundle, finalPath)
	if err != nil {
		cliexit.Reportf("export-config", "export", canonicalTool, err)

		return err
	}
	reportExportSuccess(canonicalTool, finalPath, len(files))

	return nil
}

// resolveExportFilesAndExtensions gathers files or fallback templates.
func resolveExportFilesAndExtensions(tool, srcDir string) (map[string]ConfigFilePayload, []string) {
	files, extensions := collectToolFiles(tool, srcDir)
	if len(files) == 0 {
		files = collectDefaultToolFiles(tool)
		fmt.Printf("  ℹ No local %s config found at %s; generated default template\n", tool, srcDir)
	}

	return files, extensions
}

// reportExportSuccess prints single tool export confirmation.
func reportExportSuccess(tool, path string, fileCount int) {
	fmt.Printf("%s Exported %s configuration to %s (%d file(s))\n",
		constants.ColorGreen+"✓"+constants.ColorReset, tool, path, fileCount)
}

// exportAllToolConfigs exports vscode, qtorrent, and utorrent configurations to a directory.
func exportAllToolConfigs(targetFolder string) error {
	folder := resolveExportFolder(targetFolder)
	tools := []string{"vscode", "qtorrent", "utorrent"}
	for _, tool := range tools {
		targetFile := filepath.Join(folder, tool+".json")
		if err := exportSingleToolConfig(tool, tool, targetFile); err != nil {

			return apperror.Wrap(err, "exportAllToolConfigs", map[string]any{"tool": tool})
		}
	}
	reportAllExportSuccess(folder)

	return nil
}

// resolveExportFolder resolves folder path, defaulting to current working directory.
func resolveExportFolder(targetFolder string) string {
	folder := strings.TrimSpace(targetFolder)
	if folder == "" {

		return "."
	}

	return folder
}

// reportAllExportSuccess prints batch export confirmation.
func reportAllExportSuccess(folder string) {
	fmt.Printf("%s Successfully exported all configurations to %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, folder)
}
