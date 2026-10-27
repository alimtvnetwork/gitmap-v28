// Package cmd — root_cobra_completion.go coordinates the Cobra completion tree and shell fspath.
package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var rootCompletionCmd *cobra.Command

func init() {
	completion.CustomGenerator = GenerateCobraCompletionScript
	completion.DynamicRepoSupplier = getDynamicRepoCompletions
}

// GetRootCompletionCmd builds and caches the top-level Cobra command tree for completion.
func GetRootCompletionCmd() *cobra.Command {
	if rootCompletionCmd != nil {
		return rootCompletionCmd
	}

	root := &cobra.Command{
		Use:   "gitmap",
		Short: "GitMap — Developer CLI for repository scanning, SSH fleet delegation, and AGY automation",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	root.AddCommand(cmdagy.AgyCmd)
	root.AddCommand(makeTopLevelRerunCmd())
	root.AddCommand(makeTopLevelSugCmd())
	root.AddCommand(makeTopLevelRunningPromptsCmd())
	root.AddCommand(makeTopLevelRunningProjectsCmd())
	root.AddCommand(makeTopLevelFPUGCmd())
	root.AddCommand(cmdssh.MakeDeployCobraCmd("deploy"))
	root.AddCommand(cmdssh.MakeDeployCobraCmd("deploy-right"))
	root.AddCommand(cmdssh.MakeDeployCobraCmd("deploy-left"))
	root.AddCommand(makeTopLevelSSHCmd())
	root.AddCommand(makeTopLevelUpdateCmd())
	root.AddCommand(makeTopLevelCompletionCmd())
	root.AddCommand(makeTopLevelHistoryCmd())
	root.AddCommand(makeTopLevelHPCmd())
	root.AddCommand(makeTopLevelHUCmd())
	root.AddCommand(makeTopLevelHelpCmd())
	root.AddCommand(makeTopLevelPECmd())
	root.AddCommand(makeTopLevelCFRCmd())
	root.AddCommand(makeTopLevelWPRCmd())
	root.AddCommand(makeTopLevelFolderTreeCmd())
	root.AddCommand(makeTopLevelCDCmd())
	root.AddCommand(makeTopLevelAppsCmd())
	root.AddCommand(makeTopLevelInstallCmd())
	root.AddCommand(makeTopLevelExcludeCmd())
	root.AddCommand(makeTopLevelCleanReposCmd())

	populateRemainingCommands(root)

	rootCompletionCmd = root
	return rootCompletionCmd
}

func populateRemainingCommands(root *cobra.Command) {
	existing := make(map[string]bool)
	for _, c := range root.Commands() {
		existing[c.Name()] = true
		for _, alias := range c.Aliases {
			existing[alias] = true
		}
	}

	for _, cmdName := range completion.AllCommands() {
		if existing[cmdName] || strings.TrimSpace(cmdName) == "" {
			continue
		}
		desc := resolveCommandHelpShort(cmdName)
		root.AddCommand(&cobra.Command{
			Use:   cmdName,
			Short: desc,
			Run:   func(cmd *cobra.Command, args []string) {},
		})
		existing[cmdName] = true
	}
}

func resolveCommandHelpShort(name string) string {
	switch name {
	case "clone":
		return "Clone git repositories with caching and desktop sync"
	case "pull", "pa":
		return "Fast parallel pull across all repositories"
	case "commit-in", "cin":
		return "Transfer and replay commits into target repository"
	case "status":
		return "Repository git status and branch overview"
	case "cd", "go":
		return "Navigate directly to repository directory"
	case "scan":
		return "Discover and index local git repositories"
	case "version":
		return "Display gitmap version and build information"
	case "help":
		return "Show gitmap help and documentation"
	default:
		return "GitMap command: " + name
	}
}

func executeCobraCompletion(args []string) {
	cmd := GetRootCompletionCmd()
	cmd.SetArgs(args)
	_ = cmd.Execute()
}

// GenerateCobraCompletionScript generates modern shell completion scripts using Cobra.
func GenerateCobraCompletionScript(shell string) (string, error) {
	rootCmd := GetRootCompletionCmd()
	var buf strings.Builder

	switch shell {
	case constants.ShellPowerShell:
		if err := rootCmd.GenPowerShellCompletionWithDesc(&buf); err != nil {
			return "", err
		}
		rawScript := buf.String()
		// Harden directive parsing against non-directive error output
		rawScript = strings.Replace(rawScript,
			`[int]$Directive = $Out[-1].TrimStart(':')`,
			`$Directive = 0; if ($Out -and $Out.Count -gt 0 -and $Out[-1].StartsWith(':')) { [int]::TryParse($Out[-1].TrimStart(':'), [ref]$Directive) | Out-Null }`,
			1)

		var hardened strings.Builder
		hardened.WriteString(rawScript)
		hardened.WriteString("\n# GitMap Command History & PSReadLine Predictive IntelliSense\n")
		hardened.WriteString("Import-Module PSReadLine -MinimumVersion 2.2.0 -ErrorAction SilentlyContinue\n")
		hardened.WriteString("if (Get-Command Set-PSReadLineOption -ErrorAction SilentlyContinue) {\n")
		hardened.WriteString("    try {\n")
		hardened.WriteString("        Set-PSReadLineOption -PredictionSource HistoryAndPlugin -ErrorAction SilentlyContinue\n")
		hardened.WriteString("    } catch {\n")
		hardened.WriteString("        try {\n")
		hardened.WriteString("            Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue\n")
		hardened.WriteString("        } catch {}\n")
		hardened.WriteString("    }\n")
		hardened.WriteString("    try {\n")
		hardened.WriteString("        Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue\n")
		hardened.WriteString("    } catch {}\n")
		hardened.WriteString("}\n")
		return hardened.String(), nil
	case constants.ShellBash:
		if err := rootCmd.GenBashCompletionV2(&buf, true); err != nil {
			return "", err
		}
		return buf.String(), nil
	case constants.ShellZsh:
		if err := rootCmd.GenZshCompletion(&buf); err != nil {
			return "", err
		}
		return buf.String(), nil
	case constants.ShellFish:
		if err := rootCmd.GenFishCompletion(&buf, true); err != nil {
			return "", err
		}
		return buf.String(), nil
	default:
		return "", apperror.NewValidationError(fmt.Sprintf("unsupported shell: %s", shell))
	}
}

func makeTopLevelExcludeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exclude [add|rm|list|clean] [path-or-pattern]",
		Aliases: []string{"repo-exclude", "repo-exclusions", "exclusions"},
		Short:   "Manage custom repository and directory scan exclusions",
		Run:     func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{
					"list\tList all configured repository and directory exclusions",
					"add\tAdd a path or pattern to scan exclusions",
					"rm\tRemove a path or pattern from scan exclusions",
					"clean\tPurge all excluded or corrupt repositories from gitmap.db",
				}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
	}
	return cmd
}

func makeTopLevelCleanReposCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "clean-repos",
		Aliases: []string{"cleanrepos", "repo-clean", "repoclean"},
		Short:   "Purge all excluded, cache, and corrupt repositories from gitmap.db",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
}
