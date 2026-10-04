# Subtask 216.3: Shell Completion Restoration & Dynamic Argument Suggestions

- **Parent Plan:** `pending/216-gitmap-prompting-freeze-and-suggestion-engine-fix.md`
- **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/216-gitmap-prompting-freeze-and-suggestion-engine-fix/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmd/root_cobra_completion.go`, `cli/completion/catalog.go`, `cli/cmd/cd.go`, `cli/cmd/apps.go`, `cli/cmd/install.go`

---

## 1. Objective

Restore full shell autocompletion and tab suggestions across GitMap's entire 500+ command catalog and implement dynamic argument completion for high-frequency CLI operations (`cd`, `clone`, `apps`, `install`):
1. Fix the Cobra command pruning bug in `cli/cmd/root_cobra_completion.go` by assigning non-nil dummy runners (`Run: func(cmd *cobra.Command, args []string) {}`) to all commands in `populateRemainingCommands`, ensuring `c.Runnable() == true` and `c.IsAvailableCommand() == true`.
2. Ensure all top-level commands registered on `root` are runnable or have active subcommands.
3. Wire dynamic `ValidArgsFunction` handlers for:
   - `gitmap cd` / `go`: Suggest indexed repository slugs/names from SQLite database and subcommands (`repos`, `set-default`, `clear-default`).
   - `gitmap clone`: Suggest repository names/URLs and flag completions (`--https`, `--ssh`).
   - `gitmap apps` / `app`: Suggest subcommands (`list`, `uninstall`, `help`) and dynamically complete installed application names when sub-command is `uninstall`.
   - `gitmap install` / `in`: Suggest installable developer tools (`antigravity`, `chrome`, `vscode`, `flameshot`, `git`, `docker`, etc.) and support comma-separated tool completions.

---

## 2. Root Cause Analysis (RCA)

### 2.1 The Cobra Command Pruning Bug
- In `cli/cmd/root_cobra_completion.go`, `populateRemainingCommands` loops through `completion.AllCommands()` (512 entries) and registers each command on `root`:
  ```go
  root.AddCommand(&cobra.Command{
      Use:   cmdName,
      Short: desc,
  })
  ```
- In `github.com/spf13/cobra`, command availability for completion is governed by:
  ```go
  func (c *Command) IsAvailableCommand() bool {
      if len(c.Commands()) > 0 {
          return true
      }
      return c.Runnable() && !c.Hidden
  }

  func (c *Command) Runnable() bool {
      return c.Run != nil || c.RunE != nil
  }
  ```
- Because neither `Run` nor `RunE` was set on the commands added by `populateRemainingCommands`, and they had no subcommands, `c.Runnable()` returned `false`.
- Cobra treated all 503 commands as non-runnable disabled commands and pruned them from `gitmap __complete ""` and shell completion generation. Only 9 commands survived.

### 2.2 Missing Dynamic Argument Completion
- Typing `gitmap cd <Tab>` or `gitmap clone <Tab>` previously produced no suggestions or defaulted to raw filesystem paths rather than indexed repository slugs.
- `gitmap apps uninstall <Tab>` lacked dynamic inspection of currently installed applications.
- `gitmap install <Tab>` lacked suggestion of supported tool identifiers from `constants_install.go`.

---

## 3. Step-by-Step Implementation Details

### Step 3.1: Attach Non-Nil Runners in `populateRemainingCommands`
1. In `cli/cmd/root_cobra_completion.go`:
2. Modify `populateRemainingCommands(root *cobra.Command)`:
   ```go
   for _, cmdName := range completion.AllCommands() {
       if existing[cmdName] || strings.TrimSpace(cmdName) == "" {
           continue
       }
       desc := resolveCommandHelpShort(cmdName)
       root.AddCommand(&cobra.Command{
           Use:   cmdName,
           Short: desc,
           Run:   func(cmd *cobra.Command, args []string) {}, // Guarantees c.Runnable() == true
       })
       existing[cmdName] = true
   }
   ```
3. Verify that `root` itself and any top-level commands without subcommands have a non-nil `Run` or `RunE`.

### Step 3.2: Wire Dynamic Argument Completion for `cd` / `go`
1. In `cli/cmd/root_cobra_completion.go`, implement `makeTopLevelCDCmd()`:
   ```go
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
   ```
2. Register `root.AddCommand(makeTopLevelCDCmd())` in `GetRootCompletionCmd()`.

### Step 3.3: Wire Dynamic Argument Completion for `clone`
1. Ensure `cfr` and `clone` are registered with `makeTopLevelCFRCmd()`.
2. Wire `filterRepoCompletions(toComplete)` for repository name suggestions.
3. Register flags: `--https`, `--ssh`, `--branch`.

### Step 3.4: Wire Dynamic Argument Completion for `apps` / `app`
1. In `cli/cmd/root_cobra_completion.go`, implement `makeTopLevelAppsCmd()`:
   ```go
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
   ```
2. Implement `getInstalledAppCompletions`: queries `cmdapps.ListApps` or returns known installed tool identifiers.

### Step 3.5: Wire Dynamic Argument Completion for `install` / `in`
1. In `cli/cmd/root_cobra_completion.go`, implement `makeTopLevelInstallCmd()`:
   ```go
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
   ```
2. Implement `getSupportedToolCompletions`: returns catalog of supported tools (`antigravity`, `chrome`, `vscode`, `flameshot`, `git`, `docker`, `python`, `go`, etc.) with brief descriptions.

---

## 4. Verification & Acceptance Criteria

- **All Commands Present in Completion:**
  - Running `gitmap __complete ""` yields 500+ suggestions (matching the full set of registered commands in `completion.AllCommands()`).
  - No commands are pruned due to `c.Runnable() == false`.
- **Dynamic Argument Autocompletion Verified:**
  - `gitmap cd <Tab>` suggests database-tracked repositories (e.g. `gitmap\td:\work\gitmap`).
  - `gitmap apps <Tab>` suggests `list`, `uninstall`, `help`.
  - `gitmap apps uninstall <Tab>` suggests installed applications.
  - `gitmap install <Tab>` suggests supported tool names (`antigravity`, `chrome`, `vscode`, etc.).
- **Script Generation Parity:**
  - `gitmap completion powershell`, `bash`, `zsh`, `fish` all generate scripts containing the complete 500+ command catalog.
