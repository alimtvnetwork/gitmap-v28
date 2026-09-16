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
	Short:   "Remove a project configuration (files on disk preserved)",
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

	fmt.Printf("  %s✓ Project '%s' removed from Antigravity configuration (files on disk preserved).%s\n",
		constants.ColorGreen, args[0], constants.ColorReset)

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
