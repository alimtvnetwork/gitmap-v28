package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const defaultWpConfigTemplate = `<?php
/**
 * The base configuration for WordPress
 */

define( 'DB_NAME', '%s' );
define( 'DB_USER', '%s' );
define( 'DB_PASSWORD', '%s' );
define( 'DB_HOST', '%s' );
define( 'DB_CHARSET', 'utf8mb4' );
define( 'DB_COLLATE', '' );

// >>> gitmap:wp-salts >>>
// <<< gitmap:wp-salts <<<

$table_prefix = '%s';

define( 'WP_DEBUG', %s );

if ( ! defined( 'ABSPATH' ) ) {
	define( 'ABSPATH', __DIR__ . '/' );
}

require_once ABSPATH . 'wp-settings.php';
`

// WordPressOptions holds configuration settings for WordPress setup.
type WordPressOptions struct {
	TargetDir   string
	DBName      string
	DBUser      string
	DBPassword  string
	DBHost      string
	TablePrefix string
	Domain      string
	Port        int
	IsDebug     bool
	IsVHost     bool
	IsFixPerms  bool
	IsDryRun    bool
}

func defaultWordPressOptions() WordPressOptions {
	return WordPressOptions{
		TargetDir:   ".",
		DBName:      "wordpress",
		DBUser:      "wp_user",
		DBPassword:  "",
		DBHost:      "localhost",
		TablePrefix: "wp_",
		Domain:      "wordpress.local",
		Port:        80,
		IsDebug:     false,
		IsVHost:     false,
		IsFixPerms:  false,
		IsDryRun:    false,
	}
}

func bindWpDatabaseFlags(fs *flag.FlagSet, opts *WordPressOptions) {
	fs.StringVar(&opts.TargetDir, "dir", opts.TargetDir, "Target directory")
	fs.StringVar(&opts.TargetDir, "path", opts.TargetDir, "Target directory alias")
	fs.StringVar(&opts.DBName, "db-name", opts.DBName, "Database name")
	fs.StringVar(&opts.DBUser, "db-user", opts.DBUser, "Database user")
	fs.StringVar(&opts.DBPassword, "db-pass", opts.DBPassword, "Database password")
	fs.StringVar(&opts.DBPassword, "db-password", opts.DBPassword, "Database password alias")
	fs.StringVar(&opts.DBHost, "db-host", opts.DBHost, "Database host")
	fs.StringVar(&opts.TablePrefix, "table-prefix", opts.TablePrefix, "Table prefix")
	fs.StringVar(&opts.TablePrefix, "prefix", opts.TablePrefix, "Table prefix alias")
}

func bindWpVHostFlags(fs *flag.FlagSet, opts *WordPressOptions) {
	fs.BoolVar(&opts.IsDebug, "debug", opts.IsDebug, "Enable WP_DEBUG")
	fs.BoolVar(&opts.IsVHost, "vhost", opts.IsVHost, "Generate Nginx virtual host")
	fs.StringVar(&opts.Domain, "domain", opts.Domain, "Domain name for virtual host")
	fs.IntVar(&opts.Port, "port", opts.Port, "Port for virtual host")
	fs.BoolVar(&opts.IsFixPerms, "fix-perms", opts.IsFixPerms, "Fix permissions")
	fs.BoolVar(&opts.IsDryRun, "dry-run", opts.IsDryRun, "Dry-run mode")
}

func parseWordPressFlags(args []string) (WordPressOptions, *apperror.AppError) {
	opts := defaultWordPressOptions()
	fs := flag.NewFlagSet("setup-wordpress", flag.ContinueOnError)
	bindWpDatabaseFlags(fs, &opts)
	bindWpVHostFlags(fs, &opts)
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

func replaceConfigDefine(content, key, val string) string {
	pattern := fmt.Sprintf(`(?m)define\s*\(\s*['"]%s['"]\s*,\s*['"].*?['"]\s*\);`, regexp.QuoteMeta(key))
	re := regexp.MustCompile(pattern)
	replacement := fmt.Sprintf("define( '%s', '%s' );", key, val)
	hasMatch := re.MatchString(content)
	if hasMatch {
		return re.ReplaceAllLiteralString(content, replacement)
	}

	return content
}

func replaceTablePrefix(content, prefix string) string {
	re := regexp.MustCompile(`(?m)\$table_prefix\s*=\s*['"].*?['"]\s*;`)
	replacement := fmt.Sprintf("$table_prefix = '%s';", prefix)
	hasMatch := re.MatchString(content)
	if hasMatch {
		return re.ReplaceAllLiteralString(content, replacement)
	}

	return content
}

func replaceWpDebug(content string, isDebug bool) string {
	debugStr := "false"
	if isDebug {
		debugStr = "true"
	}
	re := regexp.MustCompile(`(?m)define\s*\(\s*['"]WP_DEBUG['"]\s*,\s*.*?\s*\);`)
	replacement := fmt.Sprintf("define( 'WP_DEBUG', %s );", debugStr)
	hasMatch := re.MatchString(content)
	if hasMatch {
		return re.ReplaceAllLiteralString(content, replacement)
	}

	return content
}

func resolveInitialWpContent(baseContent string, opts WordPressOptions) string {
	hasContent := strings.TrimSpace(baseContent) != ""
	if hasContent {
		return baseContent
	}
	debugStr := "false"
	if opts.IsDebug {
		debugStr = "true"
	}

	return fmt.Sprintf(defaultWpConfigTemplate,
		opts.DBName, opts.DBUser, opts.DBPassword, opts.DBHost,
		opts.TablePrefix, debugStr)
}

// GenerateWpConfigContent generates or updates wp-config.php content with salts.
func GenerateWpConfigContent(baseContent string, opts WordPressOptions) (string, *apperror.AppError) {
	content := resolveInitialWpContent(baseContent, opts)
	content = replaceConfigDefine(content, "DB_NAME", opts.DBName)
	content = replaceConfigDefine(content, "DB_USER", opts.DBUser)
	content = replaceConfigDefine(content, "DB_PASSWORD", opts.DBPassword)
	content = replaceConfigDefine(content, "DB_HOST", opts.DBHost)
	content = replaceTablePrefix(content, opts.TablePrefix)
	content = replaceWpDebug(content, opts.IsDebug)
	salts, saltsErr := GenerateWpSalts()
	if saltsErr != nil {
		return "", saltsErr
	}

	return InjectWpSalts(content, salts), nil
}

func readFileIfExists(path string) (string, bool, *apperror.AppError) {
	hasFile := fileExists(path)
	if !hasFile {
		return "", false, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", true, apperror.WrapSimple(err, "os.ReadFile")
	}

	return string(b), true, nil
}

func resolveWpBaseContent(targetDir string) (string, *apperror.AppError) {
	cfgContent, hasCfg, cfgErr := readFileIfExists(filepath.Join(targetDir, "wp-config.php"))
	if cfgErr != nil {
		return "", cfgErr
	}
	if hasCfg {
		return cfgContent, nil
	}
	sampleContent, hasSample, sampleErr := readFileIfExists(filepath.Join(targetDir, "wp-config-sample.php"))
	if sampleErr != nil {
		return "", sampleErr
	}
	if hasSample {
		return sampleContent, nil
	}

	return "", nil
}

func writeWpConfigFile(targetDir, content string) *apperror.AppError {
	filePath := filepath.Join(targetDir, "wp-config.php")
	writeErr := os.WriteFile(filePath, []byte(content), 0600)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "os.WriteFile")
	}
	fmt.Printf("  %s✓%s WordPress config written: %s\n", constants.ColorGreen, constants.ColorReset, filePath)

	return nil
}

func writeWpConfigIfActive(targetDir, content string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	return writeWpConfigFile(targetDir, content)
}

func writeVHostOutputFile(targetDir, domain, content string) *apperror.AppError {
	outFile := filepath.Join(targetDir, fmt.Sprintf("nginx-%s.conf", domain))
	writeErr := os.WriteFile(outFile, []byte(content), 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "os.WriteFile")
	}
	fmt.Printf("  %s✓%s VHost generated: %s\n", constants.ColorGreen, constants.ColorReset, outFile)

	return nil
}

func generateWpVHost(opts WordPressOptions) *apperror.AppError {
	absPath, pathErr := filepath.Abs(opts.TargetDir)
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "filepath.Abs")
	}
	vhostCfg := VHostConfig{
		SiteType:     VHostSiteTypeWordpress,
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

func runOptionalWpVHost(opts WordPressOptions) *apperror.AppError {
	if !opts.IsVHost {
		return nil
	}

	return generateWpVHost(opts)
}

func runOptionalWpPerms(targetDir string, isFix bool) *apperror.AppError {
	if !isFix {
		return nil
	}

	return ApplyWordPressPermissions(targetDir)
}

func runPostWpSetup(opts WordPressOptions) *apperror.AppError {
	vhostErr := runOptionalWpVHost(opts)
	if vhostErr != nil {
		return vhostErr
	}
	permsErr := runOptionalWpPerms(opts.TargetDir, opts.IsFixPerms)
	if permsErr != nil {
		return permsErr
	}

	return nil
}

// SetupWordPress orchestrates configuration, salts, vhost, and permissions for WordPress.
func SetupWordPress(opts WordPressOptions) *apperror.AppError {
	base, baseErr := resolveWpBaseContent(opts.TargetDir)
	if baseErr != nil {
		return baseErr
	}
	content, genErr := GenerateWpConfigContent(base, opts)
	if genErr != nil {
		return genErr
	}
	writeErr := writeWpConfigIfActive(opts.TargetDir, content, opts.IsDryRun)
	if writeErr != nil {
		return writeErr
	}

	return runPostWpSetup(opts)
}

func runSetupWordpress(args []string) error {
	opts, parseErr := parseWordPressFlags(args)
	if parseErr != nil {
		return parseErr
	}
	appErr := SetupWordPress(opts)
	if appErr != nil {
		return appErr
	}

	return nil
}
