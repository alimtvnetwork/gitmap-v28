package cmdignore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
)

// RunIgnoreCLI routes gitmap ignore / gitmap ig commands.
func RunIgnoreCLI(args []string) error {
	if len(args) == 0 {
		return PrintIgnoreHelp()
	}
	subCmd := strings.ToLower(args[0])
	subArgs := args[1:]

	switch subCmd {
	case "help", "-h", "--help":
		return PrintIgnoreHelp()
	case "ls", "list":
		return runIgnoreList()
	case "add":
		return runIgnoreAdd(subArgs)
	case "remove", "rm":
		return runIgnoreRemove(subArgs)
	case "scan":
		return runIgnoreScan(subArgs)
	case "scan-ssh", "ss":
		return runIgnoreScanSSH(subArgs)
	case "edit":
		return runIgnoreEdit()
	case "ui", "app":
		return runIgnoreUI()
	case "add-group":
		return runAddGroup(subArgs)
	case "remove-group", "rm-grp":
		return runRemoveGroup(subArgs)
	case "set-default-group":
		return runSetDefaultGroup(subArgs)
	case "add-grp-to-default", "agtd":
		return runAddGrpToDefault(subArgs)
	case "apply":
		return runApplyGroup(subArgs)
	case "connect-group-with-repo", "cgwp":
		return runConnectGroupWithRepo(subArgs)
	case "export":
		return runExportGroups(subArgs)
	case "import":
		return runImportGroups(subArgs)
	case "action":
		return RunFixIgnoreAll(subArgs)
	default:
		return apperror.NewSimple("unknown ignore subcommand: "+subCmd, "E1003")
	}
}

func runIgnoreList() error {
	groups, _ := LoadGroups()
	fmt.Printf("\n%s  === GITMAP IGNORE GROUPS ===%s\n", constants.ColorCyan, constants.ColorReset)
	for name, grp := range groups {
		defTag := ""
		if grp.IsDefault {
			defTag = fmt.Sprintf(" %s[default]%s", constants.ColorGreen, constants.ColorReset)
		}
		fmt.Printf("  • %s%s (%d patterns):\n", name, defTag, len(grp.Patterns))
		for _, pat := range grp.Patterns {
			fmt.Printf("      %s\n", pat)
		}
	}
	fmt.Println()
	return nil
}

func runIgnoreAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore add <pattern>", "E1004")
	}
	cwd, err := os.Getwd()
	if err != nil || !gitignoreagm.IsGitRepository(cwd) {
		return apperror.NewSimple("current directory is not a git repository", "E1005")
	}
	ignorePath := filepath.Join(cwd, ".gitignore")
	data, _ := os.ReadFile(ignorePath)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(string(data), args...)
	if isModified {
		_ = os.WriteFile(ignorePath, []byte(cleaned), 0644)
		fmt.Printf("%s✓ Added patterns to %s%s\n", constants.ColorGreen, ignorePath, constants.ColorReset)
	}
	return nil
}

func runIgnoreRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore remove <pattern>", "E1006")
	}
	cwd, _ := os.Getwd()
	ignorePath := filepath.Join(cwd, ".gitignore")
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		return apperror.WrapSimple(err, "read .gitignore")
	}
	lines := strings.Split(string(data), "\n")
	var kept []string
	removeSet := make(map[string]bool)
	for _, a := range args {
		removeSet[strings.TrimSpace(a)] = true
	}
	for _, l := range lines {
		if !removeSet[strings.TrimSpace(l)] {
			kept = append(kept, l)
		}
	}
	return os.WriteFile(ignorePath, []byte(strings.Join(kept, "\n")), 0644)
}

func runIgnoreScan(args []string) error {
	records := resolveTargetRepos()
	issues := scanReposForIgnoreIssues(records)
	if len(issues) == 0 {
		fmt.Printf("%s✓ No ignore issues found across %d repositories.%s\n",
			constants.ColorGreen, len(records), constants.ColorReset)
		return nil
	}
	fmt.Printf("\n%s  === GITIGNORE SCAN ISSUES (%d found) ===%s\n", constants.ColorYellow, len(issues), constants.ColorReset)
	for _, iss := range issues {
		fmt.Printf("  • %-30s | resume=%v | duplicates=%v | missingGitmap=%v\n",
			iss.RepoName, iss.HasResumeTask, iss.HasDuplicate, iss.MissingGitmapDir)
	}
	fmt.Printf("\nRun '%sgitmap fix-ignore-all%s' to automatically remediate.\n\n",
		constants.ColorCyan, constants.ColorReset)
	return nil
}

func runIgnoreScanSSH(args []string) error {
	return RunFixIgnoresAllSSH(append([]string{"--dry-run"}, args...))
}

func runIgnoreEdit() error {
	cwd, _ := os.Getwd()
	ignorePath := filepath.Join(cwd, ".gitignore")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}
	cmd := exec.Command(editor, ignorePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runIgnoreUI() error {
	fmt.Println("Opening GitMap Repository Ignore UI...")
	cmd := exec.Command("gitmap", "ui")
	return cmd.Start()
}

func runAddGroup(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("Usage: gitmap ignore add-group <name> <pattern1> [pattern2...]", "E1007")
	}
	return AddGroup(args[0], args[1:])
}

func runRemoveGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore remove-group <name>", "E1008")
	}
	return RemoveGroup(args[0])
}

func runSetDefaultGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore set-default-group <name>", "E1009")
	}
	return SetDefaultGroup(args[0])
}

func runAddGrpToDefault(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore add-grp-to-default <name>", "E1010")
	}
	return AddGroupToDefault(args[0])
}

func runApplyGroup(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore apply <group-name>", "E1011")
	}
	cwd, _ := os.Getwd()
	return ApplyGroupToRepo(args[0], cwd)
}

func runConnectGroupWithRepo(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple("Usage: gitmap ignore connect-group-with-repo <group> <path1,alias> [--add-with-default]", "E1012")
	}
	groupName := args[0]
	targetPaths := strings.Split(args[1], ",")
	for _, p := range targetPaths {
		_ = ApplyGroupToRepo(groupName, strings.TrimSpace(p))
	}
	fmt.Printf("%s✓ Connected group '%s' to %d repositories.%s\n",
		constants.ColorGreen, groupName, len(targetPaths), constants.ColorReset)
	return nil
}

func runExportGroups(args []string) error {
	path := "gitmap-ignore-groups.json"
	if len(args) > 0 {
		path = args[0]
	}
	groups, _ := LoadGroups()
	data, _ := json.MarshalIndent(groups, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func runImportGroups(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap ignore import <file.json>", "E1013")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "read import file")
	}
	var groups map[string]IgnoreGroup
	if jsonErr := json.Unmarshal(data, &groups); jsonErr != nil {
		return apperror.WrapSimple(jsonErr, "parse json groups")
	}
	return SaveGroups(groups)
}

// PrintIgnoreHelp renders the usage menu for gitmap ignore commands.
func PrintIgnoreHelp() error {
	fmt.Printf(`
  GitMap Ignore Management Suite (ig / ignore)
    • gitmap ignore ls                             - List configured ignore groups
    • gitmap ignore add <pattern>                  - Add pattern to current repository
    • gitmap ignore scan                           - Scan all repos for ignore issues
    • gitmap ignore scan-ssh (ss)                  - Scan fleet nodes for ignore issues
    • gitmap ignore edit                           - Edit current repository .gitignore
    • gitmap ignore add-group <name> <pats...>     - Create named ignore group
    • gitmap ignore remove-group (rm-grp) <name>   - Remove named ignore group
    • gitmap ignore set-default-group <name>       - Set group as default
    • gitmap ignore add-grp-to-default (agtd) <n>  - Add group to defaults
    • gitmap ignore apply <group>                  - Apply group to current repo
    • gitmap ignore connect-group-with-repo (cgwp) - Connect group with repo
    • gitmap fix-ignore-all (fia) [-y]             - Fix & sanitize ignore across all repos
    • gitmap fix-ignores-all-ssh (fias) [-y]       - Fix & sanitize ignore across SSH fleet
`)
	return nil
}
