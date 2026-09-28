// Package cmd — root_cobra_completion.go coordinates the Cobra completion tree and shell scripts.
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var rootCompletionCmd *cobra.Command

func init() {
	completion.CustomGenerator = GenerateCobraCompletionScript
}

// GetRootCompletionCmd builds and caches the top-level Cobra command tree for completion.
func GetRootCompletionCmd() *cobra.Command {
	if rootCompletionCmd != nil {
		return rootCompletionCmd
	}

	root := &cobra.Command{
		Use:   "gitmap",
		Short: "GitMap — Developer CLI for repository scanning, SSH fleet delegation, and AGY automation",
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
	root.AddCommand(makeTopLevelHelpCmd())
	root.AddCommand(makeTopLevelPECmd())

	populateRemainingCommands(root)

	rootCompletionCmd = root
	return rootCompletionCmd
}

func makeTopLevelRerunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rerun [target] [flags]",
		Aliases: []string{"rr", "rra", "rrq", "rerun-restart", "rerun-all", "rerun-queue"},
		Short:   "Rerun active prompt(s) with IDE restart, image re-injection, and queued prefix check",
		RunE: func(c *cobra.Command, args []string) error {
			return cmdagy.RunRerunTopLevelCLI(args)
		},
	}
	cmd.Flags().AddFlagSet(cmdagy.GetAgyRerunCmd().Flags())
	cmd.ValidArgsFunction = cmdagy.GetAgyRerunCmd().ValidArgsFunction
	return cmd
}

func makeTopLevelSugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shutdown-until-green [command]",
		Aliases: []string{"sug", "shutdown-until"},
		Short:   "Monitor designated projects and trigger OS shutdown when all CI/CD pipelines turn green",
		RunE: func(c *cobra.Command, args []string) error {
			return cmdagy.RunSUGCLI(args)
		},
	}
	cmd.Flags().AddFlagSet(cmdagy.AgySUGCmd.Flags())
	cmd.ValidArgsFunction = cmdagy.AgySUGCmd.ValidArgsFunction
	return cmd
}

func makeTopLevelRunningPromptsCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "running-prompts [command]",
		Aliases:           []string{"running-prompt", "rp-prompts", "backup-running-prompts", "restore-running-prompts"},
		Short:             "Manage running and queued Antigravity prompts",
		ValidArgsFunction: cmdagy.RunningPromptsCmd.ValidArgsFunction,
	}
}

func makeTopLevelRunningProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "running-projects [ls]",
		Aliases:           []string{"runningprojects"},
		Short:             "List projects hosting active or queued Antigravity prompts",
		ValidArgsFunction: cmdagy.AgyRunningProjectsCmd.ValidArgsFunction,
	}
	cmd.Flags().AddFlagSet(cmdagy.AgyRunningProjectsCmd.Flags())
	return cmd
}

func makeTopLevelFPUGCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "finish-prompts-until-green [targets...]",
		Aliases:           []string{"fpug"},
		Short:             "Monitor target projects until prompts finish and CI/CD pipelines turn green",
		ValidArgsFunction: cmdagy.AgyFPUGCmd.ValidArgsFunction,
	}
	cmd.Flags().AddFlagSet(cmdagy.AgyFPUGCmd.Flags())
	return cmd
}

func makeTopLevelSSHCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "ssh [target-node]",
		Short:             "Connect to or manage SSH fleet nodes",
		ValidArgsFunction: cmdssh.SSHValidArgsFunction,
	}
}

func makeTopLevelUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [flags]",
		Aliases: []string{"ua", "update-all", "updateall"},
		Short:   "Update gitmap binary, dependencies, and cluster fleet",
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
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return []string{
					"ls\tList recent command executions",
					"suggest\tSuggest past commands matching prefix",
					"clear\tClear command execution history",
					"help\tShow command usage guide",
				}, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().IntP("limit", "l", 20, "Number of commands to show")
	cmd.Flags().BoolP("json", "j", false, "Output history in JSON format")
	return cmd
}

func makeTopLevelCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [shell|install]",
		Short: "Generate shell completion scripts or install profile hooks",
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
		buf.WriteString("\n# GitMap Command History & PSReadLine Predictive IntelliSense\n")
		buf.WriteString("if ((Get-Module -ListAvailable -Name PSReadLine) -and -not [Console]::IsOutputRedirected) {\n")
		buf.WriteString("    try {\n")
		buf.WriteString("        Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue\n")
		buf.WriteString("        Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue\n")
		buf.WriteString("    } catch {}\n")
		buf.WriteString("}\n")
		return buf.String(), nil
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
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

func makeTopLevelHelpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "help [command|group]",
		Short: "Display comprehensive help for commands or functional groups",
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return buildHelpCompletions(), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	return cmd
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

func makeTopLevelPECmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pe [path|alias|url] [flags]",
		Aliases: []string{"pipeline-errors", "pipeline_errors"},
		Short:   "Inspect CI/CD pipeline error logs and status for target repository",
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return getDBRepoCompletions(), cobra.ShellCompDirectiveDefault
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
