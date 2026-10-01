package cmdos

import (
	"fmt"
	"strings"
)

// DevCleanOptions configures developer tools cache sweeping.
type DevCleanOptions struct {
	IsDryRun       bool     `json:"isDryRun"`
	HasAutoYes     bool     `json:"hasAutoYes"`
	IsJSON         bool     `json:"isJson"`
	IsVerbose      bool     `json:"isVerbose"`
	IsForce        bool     `json:"isForce"`
	IsTree         bool     `json:"isTree"`
	OnlyCategories []string `json:"onlyCategories,omitempty"`
}

// ParseDevCleanOptions parses CLI arguments into DevCleanOptions.
func ParseDevCleanOptions(args []string) DevCleanOptions {
	opts := DevCleanOptions{}
	for i := 0; i < len(args); i++ {
		parseDevCleanFlagToken(args[i], args, &i, &opts)
	}
	return opts
}

func parseDevCleanOptions(args []string) DevCleanOptions {
	return ParseDevCleanOptions(args)
}

func parseDevCleanFlagToken(raw string, args []string, idx *int, opts *DevCleanOptions) {
	a := strings.ToLower(raw)
	if isDevCleanIgnoredToken(a) || parseDevCleanBooleans(a, opts) {
		return
	}
	parseDevCleanComplex(a, args, idx, opts)
}

func parseDevCleanBooleans(a string, opts *DevCleanOptions) bool {
	switch a {
	case "--force", "-f", "force":
		opts.IsForce = true
	case "--tree", "-t", "tree":
		opts.IsTree = true
	case "--dry-run", "-n", "-d", "dry-run":
		opts.IsDryRun = true
	case "--yes", "-y", "/y", "yes":
		opts.HasAutoYes = true
	case "--json":
		opts.IsJSON = true
	case "--verbose", "-v":
		opts.IsVerbose = true
	default:
		return false
	}
	return true
}

func parseDevCleanComplex(a string, args []string, idx *int, opts *DevCleanOptions) {
	if strings.HasPrefix(a, "--only=") {
		opts.OnlyCategories = append(opts.OnlyCategories, parseCategoryTokens(strings.TrimPrefix(a, "--only="))...)
		return
	}
	if a == "--only" && *idx+1 < len(args) {
		*idx++
		opts.OnlyCategories = append(opts.OnlyCategories, parseCategoryTokens(args[*idx])...)
	}
}

func parseCategoryTokens(raw string) []string {
	tokens := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == ';'
	})
	var cats []string
	for _, t := range tokens {
		if clean := strings.TrimSpace(t); len(clean) > 0 {
			cats = append(cats, clean)
		}
	}
	return cats
}

func isDevCleanIgnoredToken(a string) bool {
	tokens := " clear clean cleanup cache caches dev devs developer devtool devtools dev-tool dev-tools tool tools devtool-cache devtools-cache "
	return strings.Contains(tokens, " "+a+" ")
}

func isHelpDevClean(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" || a == "/?" {
			return true
		}
	}
	return false
}

func printOSDevCleanUsage() {
	fmt.Println("Usage: gitmap clear devtools [flags]")
	fmt.Println("       gitmap devtools clear [flags]")
	fmt.Println("       gitmap clean-dev [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -f, --force        Invalidate Split-DB cache and force fresh dynamic discovery")
	fmt.Println("  -t, --tree         Render hierarchical directory tree view with file counts and sizes")
	fmt.Println("  -n, --dry-run      Preview space reclaimed without deleting")
	fmt.Println("  -y, --yes          Bypass confirmation prompt")
	fmt.Println("      --json         Output structured JSON summary")
	fmt.Println("  -v, --verbose      Display individual subpaths and command notes")
	fmt.Println("      --only <cats>  Comma-separated categories to clean (e.g. go,npm,pnpm)")
	fmt.Println("  -h, --help         Show this documentation")
}
