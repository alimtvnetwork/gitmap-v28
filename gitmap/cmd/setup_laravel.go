package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const defaultLaravelEnvTemplate = `APP_NAME=Laravel
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_URL=http://localhost

LOG_CHANNEL=stack
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=

BROADCAST_DRIVER=log
CACHE_STORE=database
FILESYSTEM_DISK=local
QUEUE_CONNECTION=database
SESSION_DRIVER=database
SESSION_LIFETIME=120
`

// LaravelSetupOptions holds configuration settings for Laravel setup.
type LaravelSetupOptions struct {
	TargetDir     string
	EnvOpts       LaravelEnvOptions
	Domain        string
	Port          int
	IsStorageLink bool
	IsVHost       bool
	IsFixPerms    bool
	IsDryRun      bool
}

func defaultLaravelSetupOptions() LaravelSetupOptions {
	return LaravelSetupOptions{
		TargetDir: ".",
		EnvOpts: LaravelEnvOptions{
			AppName:       "Laravel",
			AppEnv:        "local",
			AppDebug:      "true",
			AppUrl:        "http://localhost",
			DBConnection:  "mysql",
			DBHost:        "127.0.0.1",
			DBPort:        "3306",
			DBDatabase:    "laravel",
			DBUsername:    "root",
			DBPassword:    "",
			CacheStore:    "database",
			SessionDriver: "database",
		},
		IsStorageLink: true,
		IsVHost:       false,
		Domain:        "laravel.local",
		Port:          80,
		IsFixPerms:    false,
		IsDryRun:      false,
	}
}

func bindLaravelEnvFlags(fs *flag.FlagSet, opts *LaravelSetupOptions) {
	fs.StringVar(&opts.TargetDir, "dir", opts.TargetDir, "Target directory")
	fs.StringVar(&opts.TargetDir, "path", opts.TargetDir, "Target directory alias")
	fs.StringVar(&opts.EnvOpts.AppName, "app-name", opts.EnvOpts.AppName, "Application name")
	fs.StringVar(&opts.EnvOpts.AppEnv, "app-env", opts.EnvOpts.AppEnv, "Application environment")
	fs.StringVar(&opts.EnvOpts.AppDebug, "app-debug", opts.EnvOpts.AppDebug, "Application debug mode")
	fs.StringVar(&opts.EnvOpts.AppUrl, "app-url", opts.EnvOpts.AppUrl, "Application URL")
	fs.StringVar(&opts.EnvOpts.DBConnection, "db-connection", opts.EnvOpts.DBConnection, "DB connection")
	fs.StringVar(&opts.EnvOpts.DBHost, "db-host", opts.EnvOpts.DBHost, "Database host")
	fs.StringVar(&opts.EnvOpts.DBPort, "db-port", opts.EnvOpts.DBPort, "Database port")
	fs.StringVar(&opts.EnvOpts.DBDatabase, "db-name", opts.EnvOpts.DBDatabase, "Database name")
	fs.StringVar(&opts.EnvOpts.DBUsername, "db-user", opts.EnvOpts.DBUsername, "Database user")
	fs.StringVar(&opts.EnvOpts.DBPassword, "db-pass", opts.EnvOpts.DBPassword, "Database password")
}

func bindLaravelVHostFlags(fs *flag.FlagSet, opts *LaravelSetupOptions) {
	fs.BoolVar(&opts.IsStorageLink, "storage-link", opts.IsStorageLink, "Run storage:link")
	fs.BoolVar(&opts.IsVHost, "vhost", opts.IsVHost, "Generate Nginx virtual host")
	fs.StringVar(&opts.Domain, "domain", opts.Domain, "Domain name for virtual host")
	fs.IntVar(&opts.Port, "port", opts.Port, "Port for virtual host")
	fs.BoolVar(&opts.IsFixPerms, "fix-perms", opts.IsFixPerms, "Fix permissions")
	fs.BoolVar(&opts.IsDryRun, "dry-run", opts.IsDryRun, "Dry-run mode")
}

func parseLaravelFlags(args []string) (LaravelSetupOptions, *apperror.AppError) {
	opts := defaultLaravelSetupOptions()
	fs := flag.NewFlagSet("setup-laravel", flag.ContinueOnError)
	bindLaravelEnvFlags(fs, &opts)
	bindLaravelVHostFlags(fs, &opts)
	parseErr := fs.Parse(reorderFlagsBeforeArgs(args))
	if parseErr == flag.ErrHelp {
		cliexit.Exit(0)
	}

	if parseErr != nil {
		return opts, apperror.WrapSimple(parseErr, "flag.Parse")
	}

	hasPositional := len(fs.Args()) > 0
	if hasPositional {
		opts.TargetDir = fs.Args()[0]
	}

	return opts, nil
}

func resolveLaravelBaseContent(targetDir string) (string, *apperror.AppError) {
	envContent, hasEnv, envErr := readFileIfExists(filepath.Join(targetDir, ".env"))
	if envErr != nil {
		return "", envErr
	}

	if hasEnv {
		return envContent, nil
	}

	sampleContent, hasSample, sampleErr := readFileIfExists(filepath.Join(targetDir, ".env.example"))
	if sampleErr != nil {
		return "", sampleErr
	}

	if hasSample {
		return sampleContent, nil
	}

	return defaultLaravelEnvTemplate, nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func ensurePublicStorageTarget(targetPath string) *apperror.AppError {
	hasTarget := dirExists(targetPath)
	if hasTarget {
		return nil
	}

	err := os.MkdirAll(targetPath, 0775)
	if err != nil {
		return apperror.WrapSimple(err, "os.MkdirAll")
	}

	return nil
}

func createWindowsStorageJunction(linkPath, targetPath string) {
	cmd := exec.Command("cmd", "/c", "mklink", "/J", linkPath, targetPath)
	_ = cmd.Run()
}

func ensureStorageParentDirs(linkPath, targetPath string) *apperror.AppError {
	_ = os.MkdirAll(filepath.Dir(linkPath), 0755)

	return ensurePublicStorageTarget(targetPath)
}

func handleWindowsStorageFallback(linkPath, targetPath string, err error) *apperror.AppError {
	if currentOS == "windows" {
		createWindowsStorageJunction(linkPath, targetPath)

		return nil
	}

	return apperror.WrapSimple(err, "os.Symlink")
}

func createStorageSymlink(targetDir string) *apperror.AppError {
	linkPath := filepath.Join(targetDir, "public", "storage")
	targetPath := filepath.Join(targetDir, "storage", "app", "public")
	ensureErr := ensureStorageParentDirs(linkPath, targetPath)
	if ensureErr != nil {
		return ensureErr
	}

	err := os.Symlink(targetPath, linkPath)
	isSafeErr := err == nil || os.IsExist(err)
	if isSafeErr {
		return nil
	}

	return handleWindowsStorageFallback(linkPath, targetPath, err)
}

func runArtisanStorageLink(targetDir string) *apperror.AppError {
	artisanPath := filepath.Join(targetDir, "artisan")
	hasArtisan := fileExists(artisanPath)
	if !hasArtisan {
		return nil
	}

	_, lookErr := exec.LookPath("php")
	if lookErr != nil {
		return createStorageSymlink(targetDir)
	}

	cmd := exec.Command("php", "artisan", "storage:link")
	cmd.Dir = targetDir
	err := cmd.Run()
	if err != nil {
		return createStorageSymlink(targetDir)
	}

	return nil
}

func generateLaravelVHost(opts LaravelSetupOptions) *apperror.AppError {
	absPath, pathErr := filepath.Abs(opts.TargetDir)
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "filepath.Abs")
	}

	vhostCfg := VHostConfig{
		SiteType:     VHostSiteTypeLaravel,
		Domain:       opts.Domain,
		DocumentRoot: absPath,
		Port:         opts.Port,
	}

	rendered, renderErr := RenderVHostConfig(vhostCfg)
	if renderErr != nil {
		return renderErr
	}

	return writeVHostOutputFile(opts.TargetDir, opts.Domain, rendered)
}

func writeLaravelEnvFile(targetDir, content string) *apperror.AppError {
	filePath := filepath.Join(targetDir, ".env")
	writeErr := os.WriteFile(filePath, []byte(content), 0600)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "os.WriteFile")
	}

	fmt.Printf("  %s✓%s Laravel .env written: %s\n", constants.ColorGreen, constants.ColorReset, filePath)

	return nil
}

func writeLaravelEnvIfActive(targetDir, content string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	return writeLaravelEnvFile(targetDir, content)
}

func runOptionalStorageLink(targetDir string, isLink bool) *apperror.AppError {
	if !isLink {
		return nil
	}

	return runArtisanStorageLink(targetDir)
}

func runOptionalLaravelVHost(opts LaravelSetupOptions) *apperror.AppError {
	if !opts.IsVHost {
		return nil
	}

	return generateLaravelVHost(opts)
}

func runOptionalLaravelPerms(targetDir string, isFix bool) *apperror.AppError {
	if !isFix {
		return nil
	}

	return ApplyLaravelPermissions(targetDir)
}

func runPostLaravelSetup(opts LaravelSetupOptions) *apperror.AppError {
	linkErr := runOptionalStorageLink(opts.TargetDir, opts.IsStorageLink)
	if linkErr != nil {
		return linkErr
	}

	vhostErr := runOptionalLaravelVHost(opts)
	if vhostErr != nil {
		return vhostErr
	}

	permsErr := runOptionalLaravelPerms(opts.TargetDir, opts.IsFixPerms)
	if permsErr != nil {
		return permsErr
	}

	return nil
}

// SetupLaravel orchestrates configuration, .env synthesis, vhost, and storage linking for Laravel.
func SetupLaravel(opts LaravelSetupOptions) *apperror.AppError {
	base, baseErr := resolveLaravelBaseContent(opts.TargetDir)
	if baseErr != nil {
		return baseErr
	}

	synthesized, synErr := SynthesizeLaravelEnv(base, opts.EnvOpts)
	if synErr != nil {
		return synErr
	}

	writeErr := writeLaravelEnvIfActive(opts.TargetDir, synthesized, opts.IsDryRun)
	if writeErr != nil {
		return writeErr
	}

	return runPostLaravelSetup(opts)
}

func runSetupLaravel(args []string) error {
	opts, parseErr := parseLaravelFlags(args)
	if parseErr != nil {
		return parseErr
	}

	appErr := SetupLaravel(opts)
	if appErr != nil {
		return appErr
	}

	return nil
}
