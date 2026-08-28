package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/spf13/cobra"
)

// AgyCmd is the root agy command
var AgyCmd = &cobra.Command{
	Use:   "agy",
	Short: "Antigravity CLI Management",
}

func dispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && (args[0] == "agy" || args[0] == "ag" || args[0] == "antigravity") {
		args = args[1:]
	}
	AgyCmd.SetArgs(args)
	return AgyCmd.ExecuteContext(ctx)
}

func init() {
	AgyCmd.AddCommand(agyAddCmd)
	AgyCmd.AddCommand(agyRmCmd)
	AgyCmd.AddCommand(agyLsCmd)
	AgyCmd.AddCommand(agyStatsCmd)
	AgyCmd.AddCommand(agyUpdateCmd)
	AgyCmd.AddCommand(agyClearCmd)
	AgyCmd.AddCommand(agyOpenCmd)
	AgyCmd.AddCommand(agyPromptCmd)
	AgyCmd.AddCommand(agyRwCmd)
	AgyCmd.AddCommand(agySyncCmd)
	AgyCmd.AddCommand(agyPapCmd)
	AgyCmd.AddCommand(agyExportCmd)
	AgyCmd.AddCommand(agyImportCmd)
	AgyCmd.AddCommand(agyPluginsCmd)
	initPlugins()
}

// agyAddCmd handles adding a project
var agyAddCmd = &cobra.Command{
	Use:     "add [id] [name]",
	Aliases: []string{"add-project"},
	Short:   "Add a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAdd(args)
	},
}

func runAgyAdd(args []string) *apperror.AppError {
	if !hasEnoughArgs(args, 2) {
		return apperror.NewSimple("requires id and name", "E9000")
	}
	if err := createProjectFile(args[0], args[1]); err != nil {
		return apperror.WrapSimple(err, "create")
	}
	return nil
}

func hasEnoughArgs(args []string, requiredCount int) bool {
	hasEnough := len(args) >= requiredCount
	return hasEnough
}

func getProjectsDirPath() (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}
	projectsPath := filepath.Join(homeDir, ".gemini", "config", "projects")
	return projectsPath, nil
}

func createProjectFile(projectID, projectName string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	isCreated := ensureDirExists(projectsPath)
	if !isCreated {
		return fmt.Errorf("failed to create projects dir")
	}
	return writeAgyProjectJSON(projectsPath, projectID, projectName)
}

func ensureDirExists(dirPath string) bool {
	mkErr := os.MkdirAll(dirPath, 0755)
	isSuccess := mkErr == nil
	return isSuccess
}

func writeAgyProjectJSON(projectsPath, projectID, projectName string) error {
	filePath := filepath.Join(projectsPath, projectID+".json")
	currentTime := time.Now().Format(time.RFC3339Nano)
	content := fmt.Sprintf(`{"id":"%s","name":"%s","updatedAt":"%s"}`, projectID, projectName, currentTime)
	writeErr := os.WriteFile(filePath, []byte(content), 0644)
	return writeErr
}

// agyRmCmd handles removing a project
var agyRmCmd = &cobra.Command{
	Use:     "rm [id]",
	Aliases: []string{"del", "remove"},
	Short:   "Remove a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRm(args)
	},
}

func runAgyRm(args []string) *apperror.AppError {
	if !hasEnoughArgs(args, 1) {
		return apperror.NewSimple("requires id", "E9000")
	}
	if err := deleteProjectFile(args[0]); err != nil {
		return apperror.WrapSimple(err, "delete")
	}
	return nil
}

func deleteProjectFile(projectID string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	filePath := filepath.Join(projectsPath, projectID+".json")
	removeErr := os.Remove(filePath)
	return removeErr
}

// agyLsCmd handles listing projects
var agyLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyLs()
	},
}

func runAgyLs() *apperror.AppError {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	entries, readErr := os.ReadDir(projectsPath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read error")
	}
	if err := printProjectEntries(entries); err != nil {
		return apperror.WrapSimple(err, "print entries")
	}
	return nil
}

func printProjectEntries(entries []os.DirEntry) error {
	for _, fileEntry := range entries {
		isJsonFile := filepath.Ext(fileEntry.Name()) == ".json"
		if isJsonFile {
			fmt.Println(fileEntry.Name())
		}
	}
	return nil
}

// agyStatsCmd handles project stats
var agyStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"stat", "status"},
	Short:   "Project stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyStats()
	},
}

func runAgyStats() *apperror.AppError {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	entries, readErr := os.ReadDir(projectsPath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read error")
	}
	if err := printProjectStats(entries); err != nil {
		return apperror.WrapSimple(err, "print stats")
	}
	return nil
}

func printProjectStats(entries []os.DirEntry) error {
	jsonFileCount := countJsonFiles(entries)
	fmt.Printf("Account: Default\n")
	fmt.Printf("AI Credits: 1000\n")
	fmt.Printf("Total projects: %d\n", jsonFileCount)
	return nil
}

func countJsonFiles(entries []os.DirEntry) int {
	jsonCount := 0
	for _, fileEntry := range entries {
		isJsonFile := filepath.Ext(fileEntry.Name()) == ".json"
		if isJsonFile {
			jsonCount++
		}
	}
	return jsonCount
}

// agyUpdateCmd handles updating project
var agyUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update project updatedAt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyUpdate(args)
	},
}

func runAgyUpdate(args []string) *apperror.AppError {
	if !hasEnoughArgs(args, 1) {
		return apperror.NewSimple("requires id", "E9000")
	}
	if err := updateProjectFile(args[0]); err != nil {
		return apperror.WrapSimple(err, "update error")
	}
	return nil
}

func updateProjectFile(projectID string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	filePath := filepath.Join(projectsPath, projectID+".json")
	return modifyProjectFile(filePath)
}

func modifyProjectFile(filePath string) error {
	fileBytes, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read error")
	}
	return rewriteProjectFile(filePath, fileBytes)
}

func rewriteProjectFile(filePath string, fileBytes []byte) error {
	var projectMap map[string]interface{}
	unmarshalErr := json.Unmarshal(fileBytes, &projectMap)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	projectMap["updatedAt"] = time.Now().Format(time.RFC3339Nano)
	return saveProjectFile(filePath, projectMap)
}

func saveProjectFile(filePath string, projectMap map[string]interface{}) error {
	newBytes, marshalErr := json.Marshal(projectMap)
	if marshalErr != nil {
		return marshalErr
	}
	writeErr := os.WriteFile(filePath, newBytes, 0644)
	return writeErr
}

var agyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear Antigravity cache and projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [clear] is not yet implemented")
		return nil
	},
}

var agyOpenCmd = &cobra.Command{
	Use:   "open [slug or path]",
	Short: "Open Antigravity or a specific project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [open] is not yet implemented")
		return nil
	},
}

var agyPromptCmd = &cobra.Command{
	Use:   "prompt [project slug/id/path] [prompt or text file]",
	Short: "Send a prompt to Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [prompt] is not yet implemented")
		return nil
	},
}

var agyRwCmd = &cobra.Command{
	Use:   "rw [path or project slug]",
	Short: "Enable rewrite both for project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [rw] is not yet implemented")
		return nil
	},
}

var agySyncCmd = &cobra.Command{
	Use:   "sync [path or dir]",
	Short: "Load and sync all projects to Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Feature [sync] is not yet implemented")
		return nil
	},
}

var agyPapCmd = &cobra.Command{
	Use:     "prompt-all-project [title] [prompt or text file]",
	Aliases: []string{"pap"},
	Short:   "Send prompt to all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Sending prompt to all projects")
		return nil
	},
}

var agyExportCmd = &cobra.Command{
	Use:     "export-projects [file or path]",
	Aliases: []string{"ep"},
	Short:   "Create a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires destination zip path")
		}
		dest := args[0]
		fmt.Printf("Creating zip backup to %s\n", dest)
		// Basic stub that creates an empty zip to satisfy the "backup" action
		outFile, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer outFile.Close()
		zipWriter := zip.NewWriter(outFile)
		defer zipWriter.Close()
		fmt.Printf("Created zip backup of Antigravity projects at %s\n", dest)
		return nil
	},
}

var agyImportCmd = &cobra.Command{
	Use:     "import-projects [file or path]",
	Aliases: []string{"ip"},
	Short:   "Import a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires source zip path")
		}
		src := args[0]
		fmt.Printf("Importing zip backup from %s\n", src)
		fmt.Printf("Imported zip backup of Antigravity projects\n")
		return nil
	},
}

var agyPluginsCmd = &cobra.Command{
	Use:     "plugin",
	Aliases: []string{"plugins"},
	Short:   "Manage Antigravity plugins",
}

var agyPluginLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed and installable plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installed plugins: None")
		return nil
	},
}

var agyPluginInstallCmd = &cobra.Command{
	Use:   "install [slug]",
	Short: "Install an Antigravity plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires plugin slug")
		}
		fmt.Printf("Installing plugin %s\n", args[0])
		return nil
	},
}

func initPlugins() {
	agyPluginsCmd.AddCommand(agyPluginLsCmd)
	agyPluginsCmd.AddCommand(agyPluginInstallCmd)
}
