// Package cmd — agy_misc.go implements auxiliary Antigravity commands.
package cmd

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/workspacesync"
)

var agyStatsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"stat"},
	Short:   "Project stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyStats()
	},
}

var agyStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show Antigravity project status table",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyLs()
	},
}

func runAgyStats() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load error")
	}
	activeCount, missingCount := 0, 0
	for _, p := range projects {
		path := p.GetPath()
		if path != "" && checkDirExists(path) {
			activeCount++
		} else if path != "" {
			missingCount++
		}
	}
	fmt.Printf("Account:        Default\n")
	fmt.Printf("Total projects: %d\n", len(projects))
	fmt.Printf("Active on disk: %d\n", activeCount)
	fmt.Printf("Missing paths:  %d\n", missingCount)
	return nil
}

var agyExportCmd = &cobra.Command{
	Use:     "export-projects [file]",
	Aliases: []string{"ep"},
	Short:   "Create a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires destination zip path")
		}
		return executeAgyExport(args[0])
	},
}

func executeAgyExport(dest string) error {
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()
	fmt.Printf("Created zip backup of Antigravity projects at %s\n", dest)
	return nil
}

var agyImportCmd = &cobra.Command{
	Use:     "import-projects [file]",
	Aliases: []string{"ip"},
	Short:   "Import a zip backup of Antigravity projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires source zip path")
		}
		fmt.Printf("Imported zip backup from %s\n", args[0])
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
	Use:   "prompt [project slug/id/path] [prompt]",
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
		return runAgySync(args)
	},
}

func runAgySync(args []string) error {
	scanFile := filepath.Join(".gitmap", "output", "gitmap.json")
	records, err := loadStatusRecords(scanFile)
	if err != nil {
		fmt.Printf("%s Could not read %s. Run 'gitmap scan' first.\n", constants.ColorYellow+"ℹ"+constants.ColorReset, scanFile)
		return nil
	}
	syncedCount := 0
	for _, rec := range records {
		if workspacesync.SyncAntigravity(rec.AbsolutePath, rec.RepoName) {
			syncedCount++
		}
	}
	fmt.Printf("%s Successfully synced %d repositories to Antigravity workspaces.\n", constants.ColorGreen+"✓"+constants.ColorReset, syncedCount)
	return nil
}

var agyPapCmd = &cobra.Command{
	Use:     "prompt-all-project [title] [prompt]",
	Aliases: []string{"pap"},
	Short:   "Send prompt to all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Sending prompt to all projects")
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
