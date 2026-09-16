// Package cmdagy — agy_projects_update.go handles updating Antigravity project timestamps.
package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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
