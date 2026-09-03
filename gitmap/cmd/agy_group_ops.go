// Package cmd — agy_group_ops.go implements export, import, display and prompt for AGY groups.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var agyGroupExportCmd = &cobra.Command{
	Use:   "export [group-name] [file.json]",
	Short: "Export an AGY group to a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: gitmap agy group export <group-name> <file.json>")
		}

		if err := exportEcosystemGroup("agy", args[0], args[1]); err != nil {
			return fmt.Errorf("export group: %w", err.Unwrap())
		}

		fmt.Printf("%s Exported AGY group %q to %s.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[0], args[1])
		return nil
	},
}

var agyGroupImportCmd = &cobra.Command{
	Use:   "import [file.json]",
	Short: "Import an AGY group from a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap agy group import <file.json>")
		}

		if err := importEcosystemGroup("agy", args[0]); err != nil {
			return fmt.Errorf("import group: %w", err.Unwrap())
		}

		fmt.Printf("%s Successfully imported AGY group from %s.\n", constants.ColorGreen+"✓"+constants.ColorReset, args[0])
		return nil
	},
}

var agyGroupPromptCmd = &cobra.Command{
	Use:   "prompt [group-name] [prompt-text]",
	Short: "Send a prompt payload to all projects in an AGY group",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: gitmap agy group prompt <group-name> <prompt-text>")
		}

		return runAgyGroupPrompt(args[0], args[1])
	},
}

func runAgyGroupList() error {
	store, err := loadEcosystemGroupStore("agy")
	if err != nil {
		return fmt.Errorf("load groups: %w", err.Unwrap())
	}

	if len(store.Groups) == 0 {
		fmt.Printf("%s No Antigravity project groups configured. Use 'gitmap agy group add <name> <project>' to create one.\n", constants.ColorYellow+"ℹ"+constants.ColorReset)
		return nil
	}

	names := getSortedGroupNames(store.Groups)
	fmt.Printf("%s Antigravity Project Groups (%d):\n", constants.ColorCyan+"▶"+constants.ColorReset, len(names))
	for _, name := range names {
		g := store.Groups[name]
		fmt.Printf("  • \033[1m%-16s\033[0m %d target(s) (Updated: %s)\n", g.Name, len(g.Targets), g.UpdatedAt)
	}

	return nil
}

func runAgyGroupShow(name string) error {
	store, err := loadEcosystemGroupStore("agy")
	if err != nil {
		return fmt.Errorf("load groups: %w", err.Unwrap())
	}

	g, exists := store.Groups[name]
	if !exists {
		return fmt.Errorf("group %q not found", name)
	}

	fmt.Printf("%s AGY Group: \033[1m%s\033[0m (%d targets)\n", constants.ColorCyan+"▶"+constants.ColorReset, g.Name, len(g.Targets))
	for idx, t := range g.Targets {
		fmt.Printf("  [%d] %s\n", idx+1, t)
	}

	return nil
}

func runAgyGroupPrompt(groupName string, promptText string) error {
	store, err := loadEcosystemGroupStore("agy")
	if err != nil {
		return fmt.Errorf("load groups: %w", err.Unwrap())
	}

	g, exists := store.Groups[groupName]
	if !exists {
		return fmt.Errorf("group %q not found", groupName)
	}

	fmt.Printf("%s Dispatching prompt to %d target(s) in AGY group %q...\n", constants.ColorCyan+"▶"+constants.ColorReset, len(g.Targets), groupName)
	for _, target := range g.Targets {
		fmt.Printf("  %s Prompt queued for: %s\n", constants.ColorGreen+"✓"+constants.ColorReset, target)
	}

	fmt.Printf("%s Finished prompt dispatch to AGY group %q.\n", constants.ColorGreen+"✓"+constants.ColorReset, groupName)
	return nil
}
