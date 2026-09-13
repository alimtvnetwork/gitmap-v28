package cmdinstall

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func verifyLinuxArchivePlatform() error {
	if runtime.GOOS != "linux" {
		msg := fmt.Sprintf("Archive installation (gitmap install tar) is supported on Linux / Ubuntu only (current OS: %s).", runtime.GOOS)
		return apperror.NewWithDetails("cmd.install.tar.platform", "E_LINUX_ONLY", msg, "cmd.install.tar", apperror.ErrorTypeValidation, apperror.SeverityError, map[string]any{"os": runtime.GOOS})
	}
	return nil
}

func parseInstallTarArgs(args []string) (ArchiveInstallOptions, error) {
	fs := flag.NewFlagSet("install-tar", flag.ContinueOnError)
	var opts ArchiveInstallOptions
	fs.StringVar(&opts.AppName, "name", "", "Custom application name")
	fs.StringVar(&opts.AppName, "n", "", "Custom application name")
	fs.BoolVar(&opts.Verbose, "v", false, "Verbose output")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Preview installation steps")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	if fs.NArg() == 0 {
		return opts, apperror.NewSimple("missing archive path. Usage: gitmap install tar <file.tar.gz|file.zip|file.gz>", "E9000")
	}
	opts.ArchivePath = fs.Arg(0)
	if opts.AppName == "" {
		opts.AppName = deriveArchiveAppName(opts.ArchivePath)
	}
	return opts, nil
}

func printArchiveStep(step, total int, message string) {
	fmt.Printf("  %s[%d/%d]%s %s\n", constants.ColorCyan, step, total, constants.ColorReset, message)
}

func runInstallTar(args []string) error {
	if err := verifyLinuxArchivePlatform(); err != nil {
		cliexit.HandleError(err, 1)
		return nil
	}

	opts, err := parseInstallTarArgs(args)
	if err != nil {
		cliexit.HandleError(err, 1)
		return nil
	}

	return executeArchiveInstallPipeline(opts)
}

func executeArchiveInstallPipeline(opts ArchiveInstallOptions) error {
	ctx := context.Background()
	absPath, err := validateArchiveSource(opts.ArchivePath)
	if err != nil {
		cliexit.HandleError(err, 1)
		return nil
	}
	opts.ArchivePath = absPath

	fmt.Printf("\nInstalling archive package: %s%s%s\n", constants.ColorWhite, filepath.Base(opts.ArchivePath), constants.ColorReset)
	printArchiveStep(1, 5, fmt.Sprintf("Verifying Linux environment and resolving archive: %s", filepath.Base(opts.ArchivePath)))
	printArchiveStep(2, 5, "Inspecting archive format and extracting to staging...")
	stagingDir, extractErr := extractArchiveStaging(ctx, opts.ArchivePath, opts.AppName)
	if extractErr != nil {
		cliexit.HandleError(extractErr, 1)
		return nil
	}
	defer os.RemoveAll(stagingDir)

	return completeArchiveInstallPipeline(ctx, opts, stagingDir)
}

func completeArchiveInstallPipeline(ctx context.Context, opts ArchiveInstallOptions, stagingDir string) error {
	printArchiveStep(3, 5, "Analyzing package structure and detecting installation strategy...")
	inspection := inspectExtractedPackage(stagingDir, opts.AppName)
	fmt.Printf("      -> Detected strategy: %s%s%s\n", constants.ColorCyan, inspection.Strategy, constants.ColorReset)

	printArchiveStep(4, 5, fmt.Sprintf("Deploying package files to destination for %s...", opts.AppName))
	deployParams := ArchiveDeployParams{Ctx: ctx, Inspection: inspection, Opts: opts, Progress: printArchiveStep}
	if err := deployArchiveApp(deployParams); err != nil {
		cliexit.HandleError(err, 1)
		return nil
	}

	printArchiveStep(5, 5, "Finalizing PATH symlinks, desktop launcher, and database records...")
	fmt.Printf("\n  %s✓%s %s successfully installed!\n\n", constants.ColorGreen, constants.ColorReset, opts.AppName)
	return nil
}

func validateArchiveSource(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", apperror.WrapSimple(err, "archive.absPath")
	}
	info, statErr := os.Stat(abs)
	if statErr != nil {
		return "", apperror.NewSimple(fmt.Sprintf("archive file '%s' does not exist", path), "E_NOT_FOUND")
	}
	if info.IsDir() {
		return "", apperror.NewSimple(fmt.Sprintf("path '%s' is a directory, expected an archive file", path), "E_NOT_FILE")
	}
	if info.Size() == 0 {
		return "", apperror.NewSimple(fmt.Sprintf("archive file '%s' is empty", path), "E_EMPTY_FILE")
	}
	return abs, nil
}
