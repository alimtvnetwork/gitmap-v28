package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdapps"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

func makeTopLevelHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help [command|group]",
		Short: "Display comprehensive help for commands or functional groups",
		Run:   func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return buildHelpCompletions(), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
}

func buildHelpCompletions() []string {
	var completions []string
	groups := []string{
		"scanning\tRepository discovery and indexing commands",
		"cloning\tHigh-speed git clone and desktop sync",
		"gitops\tGit operations (pull, status, watch, etc.)",
		"navigation\tDirectory navigation and repository aliases",
		"release\tSemantic version bumping and release workflows",
		"release-info\tChangelogs, tags, and release history",
		"data\tDatabase management, exports, and profiles",
		"import-export\tJSON config and summary import/export",
		"history\tCommand execution and commit history",
		"ssh\tSSH fleet management and node delegation",
		"integrations\tIDE, AGY, AGM, and pipeline integrations",
		"templates\tPre-compiled variables, ignore, and SEO templates",
		"search-find\tCode search, AUM search, and file finding",
		"tasks\tMacro automation and background task execution",
		"cluster\tMulti-node cluster join and fleet coordination",
	}
	completions = append(completions, groups...)
	for _, cmdName := range completion.AllCommands() {
		completions = append(completions, cmdName+"\t"+resolveCommandHelpShort(cmdName))
	}
	return completions
}

func makeTopLevelAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apps [list|uninstall|help] [flags]",
		Aliases: []string{"app"},
		Short:   "Manage, inspect, and uninstall desktop and system applications",
		Run:     func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{
					"list\tList installed applications across user and system scopes",
					"ls\tAlias for list",
					"uninstall\tUninstall an application by name or ID",
					"rm\tAlias for uninstall",
					"purge\tPurge application with user data",
					"help\tShow apps usage guide",
				}, cobra.ShellCompDirectiveNoFileComp
			}
			if len(args) == 1 && (args[0] == "uninstall" || args[0] == "rm" || args[0] == "remove" || args[0] == "purge") {
				return getInstalledAppCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().Bool("json", false, "Output apps in JSON format")
	cmd.Flags().Bool("system", false, "Filter only system-wide applications")
	cmd.Flags().Bool("user", false, "Filter only user-scoped applications")
	cmd.Flags().BoolP("all", "a", false, "Include utilities and hidden apps")
	cmd.Flags().StringP("filter", "f", "", "Fuzzy match name or package identifier")
	cmd.Flags().Bool("purge", false, "Purge configuration and user data")
	cmd.Flags().Bool("force", false, "Force uninstall without confirmation")
	cmd.Flags().Bool("dry-run", false, "Simulate uninstall without changes")
	return cmd
}

func getInstalledAppCompletions(toComplete string) []string {
	resp, err := cmdapps.ListApps(cmdapps.ListOptions{IsAll: true})
	hasValidData := err == nil && len(resp.Data) > 0
	if !hasValidData {
		return nil
	}
	clean := strings.ToLower(strings.TrimSpace(toComplete))
	var out []string
	seen := make(map[string]bool)
	for _, app := range resp.Data {
		targetName := app.Name
		if targetName == "" {
			targetName = app.ID
		}
		if seen[targetName] {
			continue
		}
		seen[targetName] = true
		desc := string(app.Manager)
		if app.Version != "" {
			desc += " " + app.Version
		}
		entry := targetName + "\t" + desc
		isMatch := clean == "" || strings.Contains(strings.ToLower(targetName), clean) || strings.Contains(strings.ToLower(app.ID), clean)
		if isMatch {
			out = append(out, entry)
		}
	}
	return out
}

func makeTopLevelInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [tool] [flags]",
		Aliases: []string{"in"},
		Short:   "Install developer tools and workstation software packages",
		Run:     func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return getSupportedToolCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().String("tools", "", "Comma-separated list of tools to install")
	cmd.Flags().Bool("dry-run", false, "Simulate installation without changes")
	cmd.Flags().BoolP("yes", "y", false, "Auto-confirm all installation prompts")
	return cmd
}

func getSupportedToolCompletions(toComplete string) []string {
	prefix := ""
	clean := strings.ToLower(strings.TrimSpace(toComplete))
	idx := strings.LastIndex(toComplete, ",")
	hasComma := idx >= 0
	if hasComma {
		prefix = toComplete[:idx+1]
		clean = strings.ToLower(strings.TrimSpace(toComplete[idx+1:]))
	}

	var out []string
	seen := make(map[string]bool)

	for tool, desc := range constants.InstallToolDescriptions {
		if seen[tool] {
			continue
		}
		seen[tool] = true
		isMatch := clean == "" || strings.HasPrefix(strings.ToLower(tool), clean) || strings.Contains(strings.ToLower(tool), clean)
		if isMatch {
			out = append(out, prefix+tool+"\t"+desc)
		}
	}

	extras := []struct {
		name string
		desc string
	}{
		{"copilot", "Windows Copilot assistant tool"},
		{"cc", "Conventional commits assistant / cache cleaner"},
	}
	for _, extra := range extras {
		if seen[extra.name] {
			continue
		}
		seen[extra.name] = true
		isMatch := clean == "" || strings.HasPrefix(strings.ToLower(extra.name), clean) || strings.Contains(strings.ToLower(extra.name), clean)
		if isMatch {
			out = append(out, prefix+extra.name+"\t"+extra.desc)
		}
	}

	if len(out) > 0 {
		return out
	}

	for tool, desc := range constants.InstallToolDescriptions {
		out = append(out, prefix+tool+"\t"+desc)
	}
	return out
}
