package cmdconfig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// ConfigInspectionResult holds the classification metadata of a configuration file.
type ConfigInspectionResult struct {
	FilePath     string
	Format       string // "JSON", "YAML", "SQLite DB", "TOML", "ENV", "INI"
	Category     string // Subsystem description
	ManageCmd    string // Suggested command
	SystemImpact string // Where it impacts
}

// inspectFilesOrBanner prints a banner when empty or returns inspection results.
func inspectFilesOrBanner(files []string, targetDir string) ([]ConfigInspectionResult, bool) {
	if len(files) == 0 {
		printWhatConfigsNoFilesBanner(targetDir)

		return nil, false
	}

	return inspectConfigFiles(files), true
}

// resolveAndInspectFiles resolves targets and processes inspection results.
func resolveAndInspectFiles(opts whatConfigsOptions) ([]ConfigInspectionResult, bool) {
	files, err := resolveConfigInspectionFiles(opts.FileArgs, opts.TargetDir)
	if err != nil {
		fmt.Printf("Error resolving configuration files: %v\n", err)

		return nil, false
	}

	return inspectFilesOrBanner(files, opts.TargetDir)
}

// RunWhatConfigsCLI handles `gitmap what-configs [globs...] [flags]`.
func RunWhatConfigsCLI(args []string) error {
	opts := parseWhatConfigsArgs(args)
	results, hasResults := resolveAndInspectFiles(opts)
	if hasResults {
		renderWhatConfigsReport(results)
	}

	return nil
}

type whatConfigsOptions struct {
	TargetDir string
	FileArgs  []string
}

func isDirFlag(a string) bool {
	return a == cfgconst.FlagArgDir || a == cfgconst.FlagArgDirShort
}

func tryConsumeDirArg(arg, nextArg string, hasNext bool, targetDir *string) (bool, int) {
	if isDirFlag(arg) && hasNext {
		*targetDir = nextArg

		return true, 1
	}

	if strings.HasPrefix(arg, cfgconst.FlagArgDirPrefix) {
		*targetDir = strings.TrimPrefix(arg, cfgconst.FlagArgDirPrefix)

		return true, 0
	}

	return false, 0
}

func applyArgToWhatConfigsOptions(arg string, nextArg string, hasNext bool, opts *whatConfigsOptions) int {
	if consumed, skip := tryConsumeDirArg(arg, nextArg, hasNext, &opts.TargetDir); consumed {
		return skip
	}

	if !strings.HasPrefix(arg, cfgconst.FlagPrefixDash) && !isWhatConfigsKeyword(arg) {
		opts.FileArgs = append(opts.FileArgs, arg)
	}

	return 0
}

func parseWhatConfigsArgs(args []string) whatConfigsOptions {
	opts := whatConfigsOptions{TargetDir: cfgconst.DefaultTargetDir}
	for i := 0; i < len(args); i++ {
		next := ""
		if i+1 < len(args) {
			next = args[i+1]
		}

		i += applyArgToWhatConfigsOptions(args[i], next, i+1 < len(args), &opts)
	}

	return opts
}

func isWhatConfigsKeyword(arg string) bool {
	low := strings.ToLower(arg)

	return low == cfgconst.CmdWhatConfigs || low == cfgconst.CmdWhatConfig || low == cfgconst.CmdWhatConfigsAlt || low == cfgconst.CmdWC
}

func resolveConfigInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandConfigFileGlobs(fileArgs), nil
	}

	return scanDirectoryAllConfigFiles(targetDir)
}

func appendPatternMatches(collected []string, pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err == nil && len(matches) > 0 {
		return append(collected, matches...)
	}

	return append(collected, pattern)
}

func expandConfigFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		collected = appendPatternMatches(collected, p)
	}

	return dedupeAndFilterConfigFiles(collected)
}

func filterSupportedConfigFiles(dir string, entries []os.DirEntry) []string {
	var matches []string
	for _, e := range entries {
		if !e.IsDir() && isSupportedConfigFile(e.Name()) {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}

	return matches
}

func scanDirectoryAllConfigFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	matches := filterSupportedConfigFiles(dir, entries)
	sort.Strings(matches)

	return matches, nil
}

func isSupportedConfigFile(name string) bool {
	low := strings.ToLower(name)
	ext := filepath.Ext(low)
	if _, err := cfgenum.ParseConfigFileFormatType(ext); err == nil {
		return true
	}

	return strings.HasPrefix(low, cfgconst.FileDotEnv) || low == cfgconst.FileDockerfile || low == cfgconst.FileMakefile || low == cfgconst.FileEditorConfig
}

func dedupeAndFilterConfigFiles(files []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, f := range files {
		clean := filepath.Clean(f)
		if seen[clean] {
			continue
		}

		seen[clean] = true
		out = append(out, clean)
	}

	return out
}

func inspectConfigFiles(files []string) []ConfigInspectionResult {
	results := make([]ConfigInspectionResult, 0, len(files))
	for _, f := range files {
		results = append(results, inspectSingleConfigFile(f))
	}

	return results
}

func buildUnreadableResult(path, ext string, err error) ConfigInspectionResult {
	return ConfigInspectionResult{
		FilePath:     path,
		Format:       strings.ToUpper(strings.TrimPrefix(ext, ".")),
		Category:     cfgconst.CategoryUnreadable,
		ManageCmd:    cfgconst.DashPlaceholder,
		SystemImpact: fmt.Sprintf("Error reading file: %v", err),
	}
}

func inspectSingleConfigFile(path string) ConfigInspectionResult {
	normPath := filepath.ToSlash(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return buildUnreadableResult(normPath, filepath.Ext(path), err)
	}

	return classifyConfigFile(normPath, filepath.Base(path), filepath.Ext(path), content)
}

func classifyDataOrDocConfigFile(normPath, base string, formatType cfgenum.ConfigFileFormatType, content []byte) (ConfigInspectionResult, bool) {
	if formatType.IsDB() || formatType.IsSQLite() {
		return inspectDBConfigFile(normPath, strings.ToLower(base), content), true
	}

	if formatType.IsTOML() {
		return inspectTOMLConfigFile(normPath, strings.ToLower(base)), true
	}

	return ConfigInspectionResult{}, false
}

func classifyStructuredConfigFile(normPath, base string, formatType cfgenum.ConfigFileFormatType, content []byte) (ConfigInspectionResult, bool) {
	if formatType.IsJSON() {
		return inspectJSONConfigFile(normPath, strings.ToLower(base), content), true
	}

	if formatType.IsYAML() {
		return inspectYAMLConfigFile(normPath, strings.ToLower(base), content), true
	}

	return classifyDataOrDocConfigFile(normPath, base, formatType, content)
}

func classifyConfigFile(normPath, base, ext string, content []byte) ConfigInspectionResult {
	formatType, _ := cfgenum.ParseConfigFileFormatType(ext)
	if res, hasMatch := classifyStructuredConfigFile(normPath, base, formatType, content); hasMatch {
		return res
	}

	return classifyFallbackConfigFile(normPath, strings.ToLower(base), ext)
}

func buildEnvResult(normPath string) ConfigInspectionResult {
	return ConfigInspectionResult{
		FilePath:     normPath,
		Format:       cfgconst.FormatENV,
		Category:     cfgconst.CategoryEnvVars,
		ManageCmd:    fmt.Sprintf("gitmap env --file %s", normPath),
		SystemImpact: cfgconst.ImpactEnvVars,
	}
}

func classifyFallbackConfigFile(normPath, base, ext string) ConfigInspectionResult {
	if strings.HasPrefix(base, cfgconst.FileDotEnv) {
		return buildEnvResult(normPath)
	}

	return ConfigInspectionResult{
		FilePath:     normPath,
		Format:       strings.ToUpper(strings.TrimPrefix(ext, ".")),
		Category:     cfgconst.CategoryGenericConfig,
		ManageCmd:    fmt.Sprintf("gitmap cat %s", normPath),
		SystemImpact: cfgconst.ImpactGenericConfig,
	}
}

func matchKnownJSONConfigFile(path, base string) ConfigInspectionResult {
	switch base {
	case cfgconst.FilePackageJSON:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatJSON, Category: cfgconst.CategoryNodeManifest, ManageCmd: "npm install / pnpm install", SystemImpact: cfgconst.ImpactNodeManifest}
	case cfgconst.FileVersionJSON:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatJSON, Category: cfgconst.CategoryVersionPin, ManageCmd: "gitmap version", SystemImpact: cfgconst.ImpactVersionPin}
	case cfgconst.FileTSConfigJSON:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatJSON, Category: cfgconst.CategoryTSConfig, ManageCmd: "tsc --noEmit", SystemImpact: cfgconst.ImpactTSConfig}
	default:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatJSON, Category: cfgconst.CategoryGenericJSON, ManageCmd: fmt.Sprintf("gitmap which-format %s", path), SystemImpact: cfgconst.ImpactGenericJSON}
	}
}

func inspectJSONConfigFile(path, base string, content []byte) ConfigInspectionResult {
	desc, _, isMatched := jsonx.DetectFormat(content)
	if isMatched {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       cfgconst.FormatJSONTypedEnvelope,
			Category:     desc.Name,
			ManageCmd:    fmt.Sprintf(desc.SuggestedImportCmd, path),
			SystemImpact: desc.SystemImpact,
		}
	}

	return matchKnownJSONConfigFile(path, base)
}

func inspectYAMLConfigFile(path, base string, content []byte) ConfigInspectionResult {
	if strings.Contains(path, ".github/workflows") {
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatYAML, Category: cfgconst.CategoryGitHubWorkflow, ManageCmd: "gitmap pipeline details", SystemImpact: cfgconst.ImpactGitHubWorkflow}
	}

	if strings.HasPrefix(base, "docker-compose") {
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatYAML, Category: cfgconst.CategoryDockerCompose, ManageCmd: "docker compose up -d", SystemImpact: cfgconst.ImpactDockerCompose}
	}

	return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatYAML, Category: cfgconst.CategoryYAMLManifest, ManageCmd: fmt.Sprintf("gitmap cat %s", path), SystemImpact: cfgconst.ImpactYAMLManifest}
}

func resolveDBFormatLabel(content []byte) string {
	if len(content) >= cfgconst.SQLiteHeaderMinLen && bytes.HasPrefix(content, []byte(cfgconst.SQLiteFormat3)) {
		return cfgconst.FormatSQLiteDB
	}

	return cfgconst.FormatBinaryDB
}

func inspectDBConfigFile(path, base string, content []byte) ConfigInspectionResult {
	formatLabel := resolveDBFormatLabel(content)
	switch {
	case base == cfgconst.FileInstallationDB:
		return ConfigInspectionResult{FilePath: path, Format: formatLabel, Category: cfgconst.CategoryInstallationDB, ManageCmd: "gitmap installer ls", SystemImpact: cfgconst.ImpactInstallationDB}
	case base == cfgconst.FileGitmapDB:
		return ConfigInspectionResult{FilePath: path, Format: formatLabel, Category: cfgconst.CategoryGitmapDB, ManageCmd: "gitmap ssh nodes", SystemImpact: cfgconst.ImpactGitmapDB}
	case strings.Contains(base, "pipeline"):
		return ConfigInspectionResult{FilePath: path, Format: formatLabel, Category: cfgconst.CategoryPipelineDB, ManageCmd: "gitmap pipeline status", SystemImpact: cfgconst.ImpactPipelineDB}
	default:
		return ConfigInspectionResult{FilePath: path, Format: formatLabel, Category: cfgconst.CategorySQLiteDB, ManageCmd: fmt.Sprintf("gitmap storage info %s", path), SystemImpact: cfgconst.ImpactSQLiteDB}
	}
}

func inspectTOMLConfigFile(path, base string) ConfigInspectionResult {
	switch base {
	case cfgconst.FileCargoTOML:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatTOML, Category: cfgconst.CategoryCargoTOML, ManageCmd: "cargo check", SystemImpact: cfgconst.ImpactCargoTOML}
	case cfgconst.FilePyprojectTOML:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatTOML, Category: cfgconst.CategoryPyprojectTOML, ManageCmd: "pip install .", SystemImpact: cfgconst.ImpactPyprojectTOML}
	default:
		return ConfigInspectionResult{FilePath: path, Format: cfgconst.FormatTOML, Category: cfgconst.CategoryTOMLFile, ManageCmd: fmt.Sprintf("gitmap cat %s", path), SystemImpact: cfgconst.ImpactTOMLFile}
	}
}

func printWhatConfigsNoFilesBanner(dir string) {
	fmt.Printf("\nNo configuration files found in: %s\n", dir)
	fmt.Printf("Usage: gitmap what-configs [pattern...] [--dir <directory>]\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  gitmap what-configs *\n")
	fmt.Printf("  gitmap what-configs *.json\n")
	fmt.Printf("  gitmap what-configs specific.json\n")
	fmt.Printf("  gitmap wc *\n\n")
}

func printWhatConfigsTip() {
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" TIP: Run 'gitmap import-all-json *' to automatically import all detected JSONs.\n")
	fmt.Printf("%s================================================================================%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func printWhatConfigsHeader(count int) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║ GITMAP CONFIGURATION FILE CLASSIFIER (WHAT-CONFIGS)                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Printf("  Inspected: %d configuration file(s)\n\n", count)
}

func renderWhatConfigsReport(results []ConfigInspectionResult) {
	printWhatConfigsHeader(len(results))
	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: cfgconst.ColFile, Align: termout.AlignLeft, MinWidth: 24},
			{Title: cfgconst.ColFormat, Align: termout.AlignLeft, MinWidth: 10},
			{Title: cfgconst.ColPurpose, Align: termout.AlignLeft, MinWidth: 26},
			{Title: cfgconst.ColRecommendedCmd, Align: termout.AlignLeft, MinWidth: 28},
		},
		Rows: buildWhatConfigsSummaryRows(results),
	}
	termout.PrintTable(cfg)
	printWhatConfigsTip()
}

func buildWhatConfigsSummaryRows(results []ConfigInspectionResult) []termout.Row {
	rows := make([]termout.Row, 0, len(results))
	for _, r := range results {
		rows = append(rows, termout.Row{
			Cells: []string{
				r.FilePath,
				r.Format,
				r.Category,
				r.ManageCmd,
			},
		})
	}

	return rows
}
