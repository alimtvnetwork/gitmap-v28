// Package cmd — agy_projects.go handles adding, removing, and updating Antigravity projects.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var agyAddCmd = &cobra.Command{
	Use:     "add [id] [name]",
	Aliases: []string{"add-project"},
	Short:   "Add a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAdd(args)
	},
}

func runAgyAdd(args []string) error {
	if !hasEnoughArgs(args, 2) {
		return apperror.NewSimple("requires id and name", "E9000")
	}
	if err := createProjectFile(args[0], args[1]); err != nil {
		return apperror.WrapSimple(err, "create")
	}
	return nil
}

func hasEnoughArgs(args []string, requiredCount int) bool {
	return len(args) >= requiredCount
}

func createProjectFile(projectID, projectName string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	if !ensureDirExists(projectsPath) {
		return fmt.Errorf("failed to create projects dir")
	}
	return writeAgyProjectJson(projectsPath, projectID, projectName)
}

func writeAgyProjectJson(projectsPath, projectID, projectName string) error {
	filePath := filepath.Join(projectsPath, projectID+".json")
	currentTime := time.Now().Format(time.RFC3339Nano)
	content := fmt.Sprintf(`{"id":"%s","name":"%s","updatedAt":"%s"}`, projectID, projectName, currentTime)
	return os.WriteFile(filePath, []byte(content), 0644)
}

var agyRmCmd = &cobra.Command{
	Use:     "rm [id]",
	Aliases: []string{"del", "remove"},
	Short:   "Remove a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRm(args)
	},
}

func runAgyRm(args []string) error {
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
	return os.Remove(filePath)
}

var agyUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update project updatedAt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyUpdate(args)
	},
}

func runAgyUpdate(args []string) error {
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
	return os.WriteFile(filePath, newBytes, 0644)
}
