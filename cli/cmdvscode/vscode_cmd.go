package cmdvscode

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func isVSCodeHelpRequested(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}

// runVSCode handles `gitmap vscode` CLI commands.
func runVSCode(args []string) error {
	if len(args) == 0 || isVSCodeHelpRequested(args) {
		RenderVSCodeHelp()

		return nil
	}

	return dispatchVSCodeAction(args)
}

func dispatchVSCodeAction(args []string) error {
	sub := strings.ToLower(args[0])
	if isVSCodeProjectSubcommand(sub) {
		return result.AsError(routeVSCodeProjectAction(sub, args))
	}

	return result.AsError(routeVSCodeMaintenanceAction(sub, args))
}

func isVSCodeProjectSubcommand(sub string) bool {
	return sub == "ls" || sub == "list" || sub == "add" || sub == "add-project" || sub == "ap" || sub == "rm" || sub == "remove" || sub == "delete" || sub == "del"
}

func routeVSCodeProjectAction(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "ls", "list":
		return result.MatchWrapper(runVSCodeLs())
	case "add", "add-project", "ap":
		handleVSCodeAdd(args)

		return result.SuccessWrapper()
	default:
		handleVSCodeRm(args)

		return result.SuccessWrapper()
	}
}

func routeVSCodeMaintenanceAction(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "pap", "prompt-all-project", "plugins", "plugin":
		fmt.Printf("Feature [vscode %s] is not yet implemented\n", sub)

		return result.SuccessWrapper()
	case "profiles", "profile":
		return result.MatchWrapper(runVSCodeProfiles(args[1:]))
	case "optimize-projects", "optimize", "--repeat-fix", "-r", "dedupe", "dedup":
		return result.MatchWrapper(runVSCodeOptimize(args[1:]))
	case "clear", "clean":
		return result.MatchWrapper(runVSCodeClear(args[1:]))
	case "group", "groups", "grp":
		return result.MatchWrapper(runVSCodeGroup(args[1:]))
	case "find-duplicates", "duplicates", "dups", "find-dups":
		return result.MatchWrapper(runFindDuplicatesVSCode())
	case "repair", "fix", "doctor":
		return result.MatchWrapper(runVSCodeRepair(args[1:]))
	case "remote":
		return result.MatchWrapper(runVSCodeRemote(args[1:]))
	default:
		printVSCodeUsage()

		return result.FailureWrapper(apperror.NewWithDetails(
			"cmd.vscode.dispatch",
			"E1024",
			fmt.Sprintf("unknown vscode subcommand '%s'", sub),
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"subcommand": sub},
		))
	}
}

func handleVSCodeAdd(args []string) {
	if len(args) < 2 {
		printVSCodeUsage()
		err := apperror.NewWithDetails(
			"cmd.vscode.add",
			"E1025",
			"missing path argument for vscode add",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)

		return
	}

	for _, p := range strings.Split(args[1], ",") {
		if p = strings.TrimSpace(p); p != "" {
			_ = runVSCodeAdd(p)
		}
	}
}

func handleVSCodeRm(args []string) {
	if len(args) < 2 {
		printVSCodeUsage()
		err := apperror.NewWithDetails(
			"cmd.vscode.rm",
			"E1026",
			"missing target argument for vscode rm",
			"cmd.vscode",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)

		return
	}

	for _, t := range strings.Split(args[1], ",") {
		if t = strings.TrimSpace(t); t != "" {
			_ = runVSCodeRm(t)
		}
	}
}

func printVSCodeUsage() {
	RenderVSCodeHelp()
}

func runVSCodeLs() error {
	entries, err := vscodepm.ListEntries()

	if isMissingUserData(err) {
		fmt.Println("No VS Code projects registered.")

		return nil
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading projects.json: %v\n", err)

		return err
	}

	if len(entries) == 0 {
		fmt.Println("No VS Code projects registered.")

		return nil
	}

	printVSCodeEntries(entries)

	return nil
}

func isMissingUserData(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, vscodepm.ErrUserDataMissing) || errors.Is(err, vscodepm.ErrExtensionMissing) || os.IsNotExist(err)
}

func runVSCodeAdd(target string) error {
	absPath, err := resolveVSCodePath(target)
	if err != nil {
		return err
	}

	name := filepath.Base(absPath)
	pair := vscodepm.Pair{RootPath: absPath, Name: name}
	summary, err := vscodepm.Sync([]vscodepm.Pair{pair})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding project: %v\n", err)

		return err
	}

	reportVSCodeAdd(summary, name, absPath)

	return nil
}

func resolveVSCodePath(target string) (string, error) {
	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", apperror.Wrap(err, "ResolveVSCodePath", map[string]any{"target": target})
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", apperror.Wrap(err, "ResolveVSCodePath", map[string]any{"path": absPath})
	}

	if !info.IsDir() {
		return "", apperror.New("ResolveVSCodePath", "INVALID_DIR", map[string]any{"path": absPath})
	}

	return absPath, nil
}

func reportVSCodeAdd(summary vscodepm.SyncSummary, name, absPath string) {
	if summary.Added > 0 || summary.Updated > 0 {
		fmt.Printf("Added/Updated %s (%s) in projects.json\n", name, absPath)

		return
	}

	fmt.Printf("Project %s already exists in projects.json\n", absPath)
}

func runVSCodeRm(target string) error {
	targetPath, err := findVSCodeTarget(target)
	if err != nil {
		return err
	}

	if targetPath == "" {
		fmt.Printf("Project not found: %s\n", target)

		return nil
	}

	if err := vscodepm.RemoveEntry(targetPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing project: %v\n", err)

		return err
	}

	fmt.Printf("Removed %s from projects.json\n", targetPath)

	return nil
}

func findVSCodeTarget(target string) (string, error) {
	entries, err := vscodepm.ListEntries()
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		isMatch := strings.EqualFold(e.Name, target) || strings.EqualFold(e.RootPath, target)
		if isMatch {
			return e.RootPath, nil
		}
	}

	return "", nil
}
