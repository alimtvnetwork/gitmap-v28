package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// AgyCmd is the root agy command
var AgyCmd = &cobra.Command{
	Use:   "agy",
	Short: "Antigravity CLI Management",
}

func dispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	AgyCmd.SetArgs(args)
	return AgyCmd.ExecuteContext(ctx)
}

func init() {
	AgyCmd.AddCommand(agyAddCmd)
	AgyCmd.AddCommand(agyRmCmd)
	AgyCmd.AddCommand(agyLsCmd)
	AgyCmd.AddCommand(agyStatsCmd)
	AgyCmd.AddCommand(agyUpdateCmd)
}

// agyAddCmd handles adding a project
var agyAddCmd = &cobra.Command{
	Use:   "add [id] [name]",
	Short: "Add a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAdd(args)
	},
}

func runAgyAdd(args []string) error {
	if !hasEnoughArgs(args, 2) {
		return fmt.Errorf("requires id and name")
	}
	return createProjectFile(args[0], args[1])
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
		return pathErr
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
	Aliases: []string{"del"},
	Short:   "Remove a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRm(args)
	},
}

func runAgyRm(args []string) error {
	if !hasEnoughArgs(args, 1) {
		return fmt.Errorf("requires id")
	}
	return deleteProjectFile(args[0])
}

func deleteProjectFile(projectID string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return pathErr
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

func runAgyLs() error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return pathErr
	}
	entries, readErr := os.ReadDir(projectsPath)
	if readErr != nil {
		return readErr
	}
	return printProjectEntries(entries)
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
	Use:   "stats",
	Short: "Project stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyStats()
	},
}

func runAgyStats() error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return pathErr
	}
	entries, readErr := os.ReadDir(projectsPath)
	if readErr != nil {
		return readErr
	}
	return printProjectStats(entries)
}

func printProjectStats(entries []os.DirEntry) error {
	jsonFileCount := countJsonFiles(entries)
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

func runAgyUpdate(args []string) error {
	if !hasEnoughArgs(args, 1) {
		return fmt.Errorf("requires id")
	}
	return updateProjectFile(args[0])
}

func updateProjectFile(projectID string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return pathErr
	}
	filePath := filepath.Join(projectsPath, projectID+".json")
	return modifyProjectFile(filePath)
}

func modifyProjectFile(filePath string) error {
	fileBytes, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return readErr
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
