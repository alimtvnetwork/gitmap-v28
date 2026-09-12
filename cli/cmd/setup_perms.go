package cmd

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PermsAppType defines the web application profile for permission auditing.
type PermsAppType string

const (
	PermsAppTypeGeneric   PermsAppType = "generic"
	PermsAppTypeWordPress PermsAppType = "wordpress"
	PermsAppTypeLaravel   PermsAppType = "laravel"
)

var currentOS = runtime.GOOS

// PermsOptions defines options for permission auditing and fixing.
type PermsOptions struct {
	TargetDir string
	AppType   PermsAppType
	Owner     string
	IsFix     bool
	IsDryRun  bool
}

// PermsReport summarizes the scanned, fixed, and violating items.
type PermsReport struct {
	ScannedCount int
	FixedCount   int
	Violations   []string
}

// ParsePermsAppType normalizes strings into a valid PermsAppType.
func ParsePermsAppType(raw string) PermsAppType {
	norm := strings.ToLower(strings.TrimSpace(raw))
	switch norm {
	case "wordpress", "wp":

		return PermsAppTypeWordPress
	case "laravel", "art", "artisan":

		return PermsAppTypeLaravel
	default:

		return PermsAppTypeGeneric
	}
}

// AuditAndFixPermissions routes permission auditing and fixing based on OS platform.
func AuditAndFixPermissions(opts PermsOptions) (*PermsReport, *apperror.AppError) {
	isWindows := currentOS == "windows"
	if isWindows {
		return ApplyWindowsPermissions(opts)
	}

	return ApplyUnixPermissions(opts)
}

// ApplyWordPressPermissions configures standard WordPress permissions.
func ApplyWordPressPermissions(targetDir string) *apperror.AppError {
	opts := PermsOptions{
		TargetDir: targetDir,
		AppType:   PermsAppTypeWordPress,
		Owner:     "www-data:www-data",
		IsFix:     true,
		IsDryRun:  false,
	}

	_, err := AuditAndFixPermissions(opts)

	return err
}

// ApplyLaravelPermissions configures standard Laravel permissions.
func ApplyLaravelPermissions(targetDir string) *apperror.AppError {
	opts := PermsOptions{
		TargetDir: targetDir,
		AppType:   PermsAppTypeLaravel,
		Owner:     "www-data:www-data",
		IsFix:     true,
		IsDryRun:  false,
	}

	_, err := AuditAndFixPermissions(opts)

	return err
}

func defaultPermsOptions() PermsOptions {
	return PermsOptions{
		TargetDir: ".",
		AppType:   PermsAppTypeGeneric,
		Owner:     "www-data:www-data",
		IsFix:     false,
		IsDryRun:  false,
	}
}

func parsePermsFlags(args []string) (PermsOptions, *apperror.AppError) {
	opts := defaultPermsOptions()
	var rawApp string
	fs := flag.NewFlagSet("perms", flag.ContinueOnError)
	fs.StringVar(&opts.TargetDir, "dir", opts.TargetDir, "Target directory")
	fs.StringVar(&opts.TargetDir, "path", opts.TargetDir, "Target directory alias")
	fs.StringVar(&rawApp, "app", string(opts.AppType), "Application type (generic, wordpress, laravel)")
	fs.StringVar(&rawApp, "type", string(opts.AppType), "Application type alias")
	fs.StringVar(&opts.Owner, "owner", opts.Owner, "Unix user:group ownership")
	fs.BoolVar(&opts.IsFix, "fix", opts.IsFix, "Apply fixes to permissions")
	fs.BoolVar(&opts.IsDryRun, "dry-run", opts.IsDryRun, "Simulate without modifying files")
	parseErr := fs.Parse(reorderFlagsBeforeArgs(args))
	if parseErr == flag.ErrHelp {
		cliexit.Exit(0)
	}

	if parseErr != nil {
		return opts, apperror.WrapSimple(parseErr, "flag.Parse")
	}

	opts.AppType = ParsePermsAppType(rawApp)
	hasPositional := len(fs.Args()) > 0
	if hasPositional {
		opts.TargetDir = fs.Args()[0]
	}

	return opts, nil
}

func printPermsSummary(report *PermsReport, isFix, isDryRun bool) {
	fmt.Printf("\n  %sPermissions Summary:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    Scanned:    %d items\n", report.ScannedCount)
	fmt.Printf("    Violations: %d items\n", len(report.Violations))
	if isFix {
		fmt.Printf("    Fixed:      %d items\n", report.FixedCount)
	}

	if isDryRun {
		fmt.Printf("    %s[dry-run] simulated execution only%s\n", constants.ColorDim, constants.ColorReset)
	}
}

func runSetupPerms(args []string) error {
	opts, parseErr := parsePermsFlags(args)
	if parseErr != nil {
		return parseErr
	}

	absPath, absErr := filepath.Abs(opts.TargetDir)
	if absErr != nil {
		return apperror.WrapSimple(absErr, "filepath.Abs")
	}

	opts.TargetDir = absPath
	report, appErr := AuditAndFixPermissions(opts)
	if appErr != nil {
		return appErr
	}

	printPermsSummary(report, opts.IsFix, opts.IsDryRun)

	return nil
}
