// Package cmdagy — agy_rm_cmd.go defines Cobra command for removing Antigravity projects.
package cmdagy

import (
	"github.com/spf13/cobra"
)

var agyRmFolderFlag string

var agyRmCmd = &cobra.Command{
	Use:     "rm [id/seq/slug...]",
	Aliases: []string{"remove", "delete", "del"},
	Short:   "Remove Antigravity project configuration (files on disk preserved)",
	Long: `Remove one or more Antigravity project configurations.
IMPORTANT: Project files and git repositories on disk are strictly PRESERVED.

Target formats:
  • Sequence number: 3-digit table index (e.g. 001, 002, 1, 2)
  • Project slug:    Prefix or exact name (e.g. gitmap)
  • Project ID:      Full UUID or short ID
  • Multiple:        Comma-separated list (e.g. 001,002 or proj-a,proj-b)
  • Folder batch:    --folder <dir> or rm folder <dir> removes all nested projects`,
	Example: `  gitmap agy rm 001
  gitmap agy rm 001,002
  gitmap agy rm gitmap
  gitmap agy rm 418bd745
  gitmap agy rm --folder /path/to/repos
  gitmap agy rm folder /path/to/repos`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeAgyRmArgs(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			renderAgyRmHelp()
			return nil
		}
		return runAgyRmEnhanced(args, agyRmFolderFlag)
	},
}

var agyRmHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show detailed help on removing Antigravity projects",
	Run: func(cmd *cobra.Command, args []string) {
		renderAgyRmHelp()
	},
}

func initAgyRmCmd() {
	agyRmCmd.Flags().StringVarP(&agyRmFolderFlag, "folder", "f", "", "Batch remove all projects under folder root")
	agyRmCmd.AddCommand(agyRmHelpCmd)
}
