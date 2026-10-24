// Package cmdpending provides command implementations for pending commits and discovery.
//
//go:generate go run ./ncgen/main.go ./ncgen/registry.go ./ncgen/constants.go ./ncgen/annotations.go ./ncgen/emit.go
package cmdpending

import (
	"slices"
	"strconv"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunNewCommands executes the new-commands discovery tool.
func RunNewCommands(args []string) error {
	opts, appErr := parseNewCommandsOptions(args)
	if appErr != nil {
		return appErr
	}
	if opts.IsHelp {
		printNewCommandsHelp()
		return nil
	}
	return executeNewCommands(opts)
}

func executeNewCommands(opts NewCommandsOptions) error {
	catalog := buildNewCommandsCatalog()
	payload := buildNewCommandsPayload(catalog, opts)
	if opts.IsJSON {
		return renderNewCommandsJSON(payload)
	}
	return renderNewCommandsTerminal(payload)
}

func parseNewCommandsOptions(args []string) (NewCommandsOptions, *appfault.AppError) {
	opts := DefaultNewCommandsOptions()
	for i := 0; i < len(args); i++ {
		nextIdx, appErr := parseNextCommandOption(args, i, &opts)
		if appErr != nil {
			return opts, appErr
		}
		i = nextIdx
	}
	clampNewCommandsLimit(&opts)
	return opts, nil
}

func parseNextCommandOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isHelpFlag(arg) {
		opts.IsHelp = true
		return idx, nil
	}
	if isJSONFlag(arg) {
		opts.IsJSON = true
		return idx, nil
	}
	return parseValuedOption(args, idx, opts)
}

func parseValuedOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isCategoryFlag(arg) {
		return parseCategoryOption(args, idx, opts)
	}
	if isFilterFlag(arg) {
		return parseFilterOption(args, idx, opts)
	}
	if isLimitFlag(arg) {
		return parseLimitOption(args, idx, opts)
	}
	return idx, nil
}

func parseCategoryOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--category=") {
		opts.Category = strings.TrimSpace(strings.TrimPrefix(arg, "--category="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-c=") {
		opts.Category = strings.TrimSpace(strings.TrimPrefix(arg, "-c="))
		return idx, nil
	}
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --category")
	}
	opts.Category = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseFilterOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--filter=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "--filter="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-f=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "-f="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-q=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "-q="))
		return idx, nil
	}
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --filter")
	}
	opts.Filter = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseLimitOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	val, nextIdx, appErr := resolveLimitValue(args, idx, arg)
	if appErr != nil {
		return idx, appErr
	}
	lim, err := strconv.Atoi(val)
	if err != nil {
		return idx, appfault.NewValidation("new-commands", "E9003", "invalid numeric value for --limit: "+val)
	}
	opts.Limit = lim
	return nextIdx, nil
}

func resolveLimitValue(args []string, idx int, arg string) (string, int, *appfault.AppError) {
	if strings.HasPrefix(arg, "--limit=") {
		return strings.TrimPrefix(arg, "--limit="), idx, nil
	}
	if strings.HasPrefix(arg, "-n=") {
		return strings.TrimPrefix(arg, "-n="), idx, nil
	}
	if idx+1 >= len(args) {
		return "", idx, appfault.NewValidation("new-commands", "E9003", "missing value for --limit")
	}
	return args[idx+1], idx + 1, nil
}

func clampNewCommandsLimit(opts *NewCommandsOptions) {
	if opts.Limit <= 0 {
		opts.Limit = 100
	}
}

func isHelpFlag(arg string) bool {
	return slices.Contains([]string{"-h", "--help", "help"}, arg)
}

func isJSONFlag(arg string) bool {
	return slices.Contains([]string{"--json", "-j"}, arg)
}

func isCategoryFlag(arg string) bool {
	return slices.Contains([]string{"--category", "-c"}, arg) || strings.HasPrefix(arg, "--category=") || strings.HasPrefix(arg, "-c=")
}

func isFilterFlag(arg string) bool {
	return slices.Contains([]string{"--filter", "-f", "-q"}, arg) ||
		strings.HasPrefix(arg, "--filter=") || strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "-q=")
}

func isLimitFlag(arg string) bool {
	return slices.Contains([]string{"--limit", "-n"}, arg) || strings.HasPrefix(arg, "--limit=") || strings.HasPrefix(arg, "-n=")
}

// buildNewCommandsCatalog returns the new-commands catalog generated at
// build time by ncgen (see nc_catalog_gen.go). Regenerate with:
// go generate ./cli/cmdpending/
func buildNewCommandsCatalog() []NewCommandEntry {
	return generatedCatalogEntries()
}
