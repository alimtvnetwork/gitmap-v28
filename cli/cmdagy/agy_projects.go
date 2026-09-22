// Package cmdagy — agy_projects.go handles adding and removing Antigravity projects.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
		return apperror.WrapSimple(pathErr, "get projects dir path")
	}

	if err := ensureDirExistsErr(projectsPath); err != nil {
		return apperror.WrapSimple(err, "ensure projects dir exists")
	}

	return writeAgyProjectJson(projectsPath, projectID, projectName)
}

func writeAgyProjectJson(projectsPath, projectID, projectName string) error {
	filePath := filepath.Join(projectsPath, projectID+".json")
	currentTime := time.Now().Format(time.RFC3339Nano)
	content := fmt.Sprintf(`{"id":"%s","name":"%s","updatedAt":"%s"}`, projectID, projectName, currentTime)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "write agy project json")
	}
	return nil
}


func deleteProjectFile(projectID string) error {
	projectsPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "get projects dir path")
	}

	filePath := filepath.Join(projectsPath, projectID+".json")

	if err := os.Remove(filePath); err != nil {
		return apperror.WrapSimple(err, "remove project json file")
	}
	return nil
}
