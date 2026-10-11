package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfoldertree"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhistory"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

func makeTopLevelSSHCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "ssh [target-node]",
		Short:             "Connect to or manage SSH fleet nodes",
		Run:               func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: cmdssh.SSHValidArgsFunction,
	}
}

func makeTopLevelUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [flags]",
		Aliases: []string{"ua", "update-all", "updateall"},
		Short:   "Update gitmap binary, dependencies, and cluster fleet",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().String("remote", "", "Update target remote cluster node")
	cmd.Flags().BoolP("all", "a", false, "Update all cluster fleet nodes")
	cmd.Flags().BoolP("json", "j", false, "Output update results in JSON format")
	cmd.Flags().Bool("rebuild", false, "Force local binary rebuild")
	cmd.Flags().Bool("self", false, "Self-update gitmap executable")
	_ = cmd.RegisterFlagCompletionFunc("remote", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cmdssh.GetFleetNodeCompletions(), cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func makeTopLevelHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "history [command]",
		Aliases: []string{"hist"},
		Short:   "View past command execution history and reuse suggestions",
		Run:     func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{
					"ls\tList recent command executions",
					"suggest\tSuggest past commands matching prefix",
					"clear\tClear command execution history",
					"purge\tPurge files, folders, or commits from Git history",
					"undo\tUndo a previous history purge operation and restore history",
					"help\tShow command usage guide",
				}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().IntP("limit", "l", 20, "Number of commands to show")
	cmd.Flags().BoolP("json", "j", false, "Output history in JSON format")
	cmd.AddCommand(cmdhistory.NewHistoryPurgeCmd())
	cmd.AddCommand(cmdhistory.NewHistoryUndoCmd())
	return cmd
}

func makeTopLevelHPCmd() *cobra.Command {
	cmd := cmdhistory.NewHistoryPurgeCmd()
	cmd.Use = "hp [flags]"
	cmd.Aliases = []string{"history-purge"}
	return cmd
}

func makeTopLevelHUCmd() *cobra.Command {
	cmd := cmdhistory.NewHistoryUndoCmd()
	cmd.Use = "hu [operation-id] [flags]"
	cmd.Aliases = []string{"history-undo", "undo-history"}
	return cmd
}

func makeTopLevelCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [shell|install]",
		Short: "Generate shell completion scripts or install profile hooks",
		Run:   func(cmd *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{
					"powershell\tGenerate PowerShell completion script",
					"bash\tGenerate Bash completion script",
					"zsh\tGenerate Zsh completion script",
					"fish\tGenerate Fish completion script",
					"install\tInstall tab completion and predictive suggestions into shell profile",
				}, cobra.ShellCompDirectiveNoFileComp
			}
			if len(args) == 1 && args[0] == "install" {
				return []string{"powershell", "bash", "zsh", "fish"}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
}

func makeTopLevelPECmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pe [path|alias|url] [flags]",
		Aliases: []string{"pipeline-errors", "pipeline_errors", "ee", "te"},
		Short:   "Inspect CI/CD pipeline error logs and status for target repository",
		Run:     func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return filterRepoCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().Bool("json", false, "Output error report in JSON format")
	cmd.Flags().BoolP("fix", "f", false, "Generate actionable fix suggestions")
	cmd.Flags().BoolP("check", "c", false, "Check live pipeline status")
	cmd.Flags().BoolP("detailed", "v", false, "Display detailed stack traces")
	return cmd
}

func makeTopLevelCFRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cfr [repo|alias|url] [flags]",
		Aliases: []string{"clone-fix-repo", "cfrp", "clone-fix-repo-pub", "clone"},
		Short:   "Clone, inspect, and fix repository with desktop sync",
		Run:     func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return filterRepoCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().Bool("https", false, "Force HTTPS clone")
	cmd.Flags().Bool("ssh", false, "Force SSH clone")
	cmd.Flags().String("branch", "", "Target branch to clone")
	return cmd
}

func filterRepoCompletions(toComplete string) []string {
	all := getDBRepoCompletions()
	if toComplete == "" {
		return all
	}
	clean := strings.ToLower(strings.TrimSpace(toComplete))
	var filtered []string
	for _, entry := range all {
		parts := strings.Split(entry, "\t")
		name := strings.ToLower(parts[0])
		if strings.HasPrefix(name, clean) || strings.Contains(name, clean) {
			filtered = append(filtered, entry)
		}
	}
	if len(filtered) > 0 {
		return filtered
	}
	return all
}

func getDBRepoCompletions() []string {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer db.Close()
	repos, errList := db.ListRepos()
	if errList != nil {
		return nil
	}
	out := make([]string, 0, len(repos))
	for _, r := range repos {
		out = append(out, r.Slug+"\t"+r.AbsolutePath)
		if r.RepoName != r.Slug {
			out = append(out, r.RepoName+"\t"+r.AbsolutePath)
		}
	}
	return out
}

func getDynamicRepoCompletions() []string {
	db, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer db.Close()
	repos, errList := db.ListRepos()
	if errList != nil {
		return nil
	}
	return extractRepoNamesAndSlugs(repos)
}

func extractRepoNamesAndSlugs(repos []model.ScanRecord) []string {
	seen := make(map[string]bool, len(repos)*2)
	out := make([]string, 0, len(repos)*2)
	for _, r := range repos {
		if r.RepoName != "" && !seen[r.RepoName] {
			seen[r.RepoName] = true
			out = append(out, r.RepoName)
		}
		if r.Slug != "" && !seen[r.Slug] {
			seen[r.Slug] = true
			out = append(out, r.Slug)
		}
	}
	return out
}

func makeTopLevelFolderTreeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "folder-tree [path]",
		Aliases: []string{"ft", "foldertree"},
		Short:   "Render folder & repo tree with emoji structure, sequence numbers, export and import",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdfoldertree.RunFolderTree(args)
		},
	}
}

func makeTopLevelCDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cd [repo|alias] [flags]",
		Aliases: []string{"go"},
		Short:   "Navigate directly to a tracked repository directory",
		Run:     func(c *cobra.Command, args []string) {},
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				subcommands := []string{
					"repos\tList all tracked repository locations",
					"set-default\tSet a default directory for ambiguous repository",
					"clear-default\tClear default directory mapping",
				}
				repos := filterRepoCompletions(toComplete)
				return append(subcommands, repos...), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().BoolP("pick", "p", false, "Force interactive picker even if default set")
	cmd.Flags().String("group", "", "Filter repos list by group")
	return cmd
}
