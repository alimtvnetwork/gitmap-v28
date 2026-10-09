package cmdconfig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

// runWhatConfigsCLI handles `gitmap what-configs [globs...] [flags]`.
func RunWhatConfigsCLI(args []string) error {
	opts := parseWhatConfigsArgs(args)
	files, err := resolveConfigInspectionFiles(opts.FileArgs, opts.TargetDir)
	if err != nil {
		fmt.Printf("Error resolving configuration files: %v\n", err)
		return nil
	}

	if len(files) == 0 {
		printWhatConfigsNoFilesBanner(opts.TargetDir)
		return nil
	}

	results := inspectConfigFiles(files)
	renderWhatConfigsReport(results)
	return nil
}

type whatConfigsOptions struct {
	TargetDir string
	FileArgs  []string
}

func parseWhatConfigsArgs(args []string) whatConfigsOptions {
	opts := whatConfigsOptions{TargetDir: "."}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--dir" || a == "-d") && i+1 < len(args) {
			opts.TargetDir = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--dir=") {
			opts.TargetDir = strings.TrimPrefix(a, "--dir=")
			continue
		}
		if strings.HasPrefix(a, "-") {
			continue
		}
		if isWhatConfigsKeyword(a) {
			continue
		}
		opts.FileArgs = append(opts.FileArgs, a)
	}
	return opts
}

func isWhatConfigsKeyword(arg string) bool {
	low := strings.ToLower(arg)
	return low == "what-configs" || low == "what-config" || low == "whatconfigs" || low == "wc"
}

func resolveConfigInspectionFiles(fileArgs []string, targetDir string) ([]string, error) {
	if len(fileArgs) > 0 {
		return expandConfigFileGlobs(fileArgs), nil
	}
	return scanDirectoryAllConfigFiles(targetDir)
}

func expandConfigFileGlobs(patterns []string) []string {
	var collected []string
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err == nil && len(matches) > 0 {
			collected = append(collected, matches...)
			continue
		}
		collected = append(collected, p)
	}
	return dedupeAndFilterConfigFiles(collected)
}

func scanDirectoryAllConfigFiles(dir string) ([]string, error) {
	var matches []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isSupportedConfigFile(e.Name()) {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func isSupportedConfigFile(name string) bool {
	low := strings.ToLower(name)
	ext := filepath.Ext(low)
	if ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" || ext == ".toml" || ext == ".ini" || ext == ".cfg" || ext == ".conf" {
		return true
	}
	return strings.HasPrefix(low, ".env") || low == "dockerfile" || low == "makefile" || low == ".editorconfig"
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
		res := inspectSingleConfigFile(f)
		results = append(results, res)
	}
	return results
}

func inspectSingleConfigFile(path string) ConfigInspectionResult {
	normPath := filepath.ToSlash(path)
	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))

	content, err := os.ReadFile(path)
	if err != nil {
		return ConfigInspectionResult{
			FilePath:     normPath,
			Format:       strings.ToUpper(strings.TrimPrefix(ext, ".")),
			Category:     "Unreadable File",
			ManageCmd:    "-",
			SystemImpact: fmt.Sprintf("Error reading file: %v", err),
		}
	}

	if ext == ".json" {
		return inspectJSONConfigFile(normPath, base, content)
	}
	if ext == ".yaml" || ext == ".yml" {
		return inspectYAMLConfigFile(normPath, base, content)
	}
	if ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" {
		return inspectDBConfigFile(normPath, base, content)
	}
	if ext == ".toml" {
		return inspectTOMLConfigFile(normPath, base)
	}
	if strings.HasPrefix(base, ".env") {
		return ConfigInspectionResult{
			FilePath:     normPath,
			Format:       "ENV",
			Category:     "Environment Variables & Secrets",
			ManageCmd:    fmt.Sprintf("gitmap env --file %s", normPath),
			SystemImpact: "Injects environment credentials and service runtime variables.",
		}
	}
	return ConfigInspectionResult{
		FilePath:     normPath,
		Format:       strings.ToUpper(strings.TrimPrefix(ext, ".")),
		Category:     "Generic Configuration",
		ManageCmd:    fmt.Sprintf("gitmap cat %s", normPath),
		SystemImpact: "Application configuration settings.",
	}
}

func inspectJSONConfigFile(path, base string, content []byte) ConfigInspectionResult {
	desc, _, isMatched := jsonx.DetectFormat(content)
	if isMatched {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "JSON (Typed Envelope)",
			Category:     desc.Name,
			ManageCmd:    fmt.Sprintf(desc.SuggestedImportCmd, path),
			SystemImpact: desc.SystemImpact,
		}
	}
	if base == "package.json" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "JSON",
			Category:     "Node.js Package Manifest & Dependencies",
			ManageCmd:    "npm install / pnpm install",
			SystemImpact: "Defines JavaScript/TypeScript project dependencies and fspath.",
		}
	}
	if base == "version.json" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "JSON",
			Category:     "GitMap Canonical Version Pin Site",
			ManageCmd:    "gitmap version",
			SystemImpact: "Defines project release semantic version and build metadata.",
		}
	}
	if base == "tsconfig.json" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "JSON",
			Category:     "TypeScript Compiler Options",
			ManageCmd:    "tsc --noEmit",
			SystemImpact: "Governs TypeScript build rules, target ECMAScript version, and path aliases.",
		}
	}
	return ConfigInspectionResult{
		FilePath:     path,
		Format:       "JSON",
		Category:     "Generic JSON Document",
		ManageCmd:    fmt.Sprintf("gitmap which-format %s", path),
		SystemImpact: "Generic JSON payload (inspect with 'gitmap which-format').",
	}
}

func inspectYAMLConfigFile(path, base string, content []byte) ConfigInspectionResult {
	if strings.Contains(path, ".github/workflows") {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "YAML",
			Category:     "GitHub Actions CI/CD Pipeline Workflow",
			ManageCmd:    "gitmap pipeline details",
			SystemImpact: "Orchestrates remote GitHub Actions runners, linters, tests, and releases.",
		}
	}
	if strings.HasPrefix(base, "docker-compose") {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "YAML",
			Category:     "Docker Compose Container Infrastructure",
			ManageCmd:    "docker compose up -d",
			SystemImpact: "Defines multi-container application services, networks, and storage volumes.",
		}
	}
	return ConfigInspectionResult{
		FilePath:     path,
		Format:       "YAML",
		Category:     "YAML Configuration Manifest",
		ManageCmd:    fmt.Sprintf("gitmap cat %s", path),
		SystemImpact: "Application declarative configuration.",
	}
}

func inspectDBConfigFile(path, base string, content []byte) ConfigInspectionResult {
	isSQLite := len(content) >= 16 && bytes.HasPrefix(content, []byte("SQLite format 3"))
	formatLabel := "SQLite DB"
	if !isSQLite {
		formatLabel = "Binary DB"
	}
	if base == "installation.db" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       formatLabel,
			Category:     "GitMap System & Tools Installation Registry",
			ManageCmd:    "gitmap installer ls",
			SystemImpact: "Tracks installed packages, tool versions, and system telemetry.",
		}
	}
	if base == "gitmap.db" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       formatLabel,
			Category:     "GitMap Core System & SSH Nodes Split Database",
			ManageCmd:    "gitmap ssh nodes",
			SystemImpact: "Stores cluster fleet nodes, SSH host definitions, credentials, and settings.",
		}
	}
	if strings.Contains(base, "pipeline") {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       formatLabel,
			Category:     "GitMap Local CI/CD Pipeline Database",
			ManageCmd:    "gitmap pipeline status",
			SystemImpact: "Maintains local CI/CD pipeline task executions, outputs, and status logs.",
		}
	}
	return ConfigInspectionResult{
		FilePath:     path,
		Format:       formatLabel,
		Category:     "SQLite Relational Database",
		ManageCmd:    fmt.Sprintf("gitmap storage info %s", path),
		SystemImpact: "SQLite storage database (schema tables, rows, indexes).",
	}
}

func inspectTOMLConfigFile(path, base string) ConfigInspectionResult {
	if base == "cargo.toml" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "TOML",
			Category:     "Rust Cargo Package & Workspace Manifest",
			ManageCmd:    "cargo check",
			SystemImpact: "Defines Rust crate dependencies, compilation targets, and features.",
		}
	}
	if base == "pyproject.toml" {
		return ConfigInspectionResult{
			FilePath:     path,
			Format:       "TOML",
			Category:     "Python PEP 518/621 Build & Packaging Manifest",
			ManageCmd:    "pip install .",
			SystemImpact: "Defines Python build dependencies, package configuration, and linter settings.",
		}
	}
	return ConfigInspectionResult{
		FilePath:     path,
		Format:       "TOML",
		Category:     "TOML Configuration File",
		ManageCmd:    fmt.Sprintf("gitmap cat %s", path),
		SystemImpact: "Structured TOML settings file.",
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

func renderWhatConfigsReport(results []ConfigInspectionResult) {
	fmt.Printf("\n╔══════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║ GITMAP CONFIGURATION FILE CLASSIFIER (WHAT-CONFIGS)                          ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("  Inspected: %d configuration file(s)\n\n", len(results))

	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: "FILE", Align: termout.AlignLeft, MinWidth: 24},
			{Title: "FORMAT", Align: termout.AlignLeft, MinWidth: 10},
			{Title: "PURPOSE / SUBSYSTEM", Align: termout.AlignLeft, MinWidth: 26},
			{Title: "RECOMMENDED CLI COMMAND", Align: termout.AlignLeft, MinWidth: 28},
		},
		Rows: buildWhatConfigsSummaryRows(results),
	}
	termout.PrintTable(cfg)
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" TIP: Run 'gitmap import-all-json *' to automatically import all detected JSONs.\n")
	fmt.Printf("%s================================================================================%s\n\n", constants.ColorCyan, constants.ColorReset)
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
