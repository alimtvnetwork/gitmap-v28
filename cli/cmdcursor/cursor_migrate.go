package cmdcursor

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type cursorMigrateOptions struct {
	targetNode     string
	exclude        string
	outputPath     string
	isSkipSettings bool
	isSettingsOnly bool
	isDryRun       bool
	isNoBackup     bool
	isForce        bool
}

type CursorMigrateOptions = cursorMigrateOptions

func renderCursorMigrateHelp() {
	fmt.Printf("\n  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║        gitmap cursor migrate - Cursor Fleet IDE Migration        ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Print("  Usage: gitmap cursor migrate [flags] (aliases: cur m, delegate, del)\n" +
		"  Flags:\n" +
		"    -n, --node <alias>      Target fleet node (default: u1)\n" +
		"        --exclude <list>    Comma-separated project/repo names to exclude\n" +
		"    -o, --output <path>     Custom destination archive bundle path\n" +
		"        --skip-settings     Exclude settings.json from migration\n" +
		"        --settings-only     Package and sync only settings.json and themes\n" +
		"    -d, --dry-run           Simulate operations without making changes\n" +
		"        --no-backup         Skip remote pre-flight backup snapshot\n" +
		"    -f, --force             Force overwrite existing remote configurations\n\n")
}

func isMigrateHelpRequested(args []string) bool {
	for _, arg := range args {
		lower := strings.ToLower(arg)
		if lower == "help" || lower == "--help" || lower == "-h" {
			return true
		}
	}

	return false
}

func parseCursorMigrateArgs(args []string) (cursorMigrateOptions, error) {
	opts := cursorMigrateOptions{targetNode: "u1"}

	for i := 0; i < len(args); i++ {
		arg := strings.ToLower(args[i])
		switch {
		case (arg == "--node" || arg == "-n") && i+1 < len(args):
			opts.targetNode, i = args[i+1], i+1
		case arg == "--exclude" && i+1 < len(args):
			opts.exclude, i = args[i+1], i+1
		case (arg == "--output" || arg == "-o") && i+1 < len(args):
			opts.outputPath, i = args[i+1], i+1
		case arg == "--skip-settings":
			opts.isSkipSettings = true
		case arg == "--settings-only":
			opts.isSettingsOnly = true
		case arg == "--dry-run" || arg == "-d":
			opts.isDryRun = true
		case arg == "--no-backup":
			opts.isNoBackup = true
		case arg == "--force" || arg == "-f":
			opts.isForce = true
		}
	}

	if opts.isSkipSettings && opts.isSettingsOnly {
		return opts, apperror.NewWithDetails("cmd.cursor.migrate", "E1035", "cannot specify both --skip-settings and --settings-only", "cmdcursor", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}

	return opts, nil
}

// RunCursorMigrate handles `gitmap cursor migrate` / `gitmap cur m`.
func RunCursorMigrate(args []string) error {
	if isMigrateHelpRequested(args) {
		renderCursorMigrateHelp()

		return nil
	}

	opts, err := parseCursorMigrateArgs(args)
	if err != nil {
		return err
	}

	return executeCursorMigration(opts, args)
}
