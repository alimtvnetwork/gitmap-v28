// Package cmderrors provides the CLI command implementation for inspecting internal errors.
package cmderrors

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunErrorsCLI is the entry point for 'gitmap errors' and 'gitmap e'.
func RunErrorsCLI(args []string) error {
	if isErrorsHelpRequested(args) {
		renderErrorsHelp()

		return nil
	}

	if isErrorsClearRequested(args) {
		return handleClearErrors()
	}

	if id, isShow := extractShowID(args); isShow {
		return handleShowError(id)
	}

	return handleListErrors(args)
}

func isErrorsHelpRequested(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(args[0])

	return first == "help" || first == "--help" || first == "-h"
}

func isErrorsClearRequested(args []string) bool {
	for _, arg := range args {
		lower := strings.ToLower(arg)
		if lower == "clear" || lower == "clean" || lower == "--clear" || lower == "-c" {
			return true
		}
	}

	return false
}

func extractShowID(args []string) (int64, bool) {
	if len(args) == 0 {
		return 0, false
	}

	target := args[0]
	if (target == "show" || target == "get" || target == "view") && len(args) > 1 {
		target = args[1]
	}

	id, err := strconv.ParseInt(target, 10, 64)
	if err == nil && id > 0 {
		return id, true
	}

	return 0, false
}

func handleClearErrors() error {
	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	if err := db.ClearErrors(); err != nil {
		return fmt.Errorf("clear errors: %w", err)
	}

	fmt.Println("  ✓ Cleared all recorded internal errors from gitmap-errors.db.")

	return nil
}

func handleShowError(id int64) error {
	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	rec, err := db.GetError(id)
	if err != nil {
		return apperror.NewNotFound("internal_error", "E404", fmt.Sprintf("internal error #%d not found", id))
	}

	renderSingleError(rec)

	return nil
}

func handleListErrors(args []string) error {
	isJSON := hasFlag(args, "--json")
	limit := parseLimitArg(args, 25)

	db, err := store.OpenErrorsSplitDB()
	if err != nil {
		return fmt.Errorf("open errors db: %w", err)
	}

	defer db.Close()

	records, err := db.ListErrors(limit, false)
	if err != nil {
		return fmt.Errorf("list errors: %w", err)
	}

	if isJSON {
		return emitJSONErrors(records)
	}

	renderErrorsTable(records, db.Path)

	return nil
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if strings.EqualFold(arg, flag) {
			return true
		}
	}

	return false
}

func parseLimitArg(args []string, defaultLimit int) int {
	for i, arg := range args {
		if !isLimitFlag(arg) || i+1 >= len(args) {
			continue
		}

		if parsed, err := strconv.Atoi(args[i+1]); err == nil && parsed > 0 {
			return parsed
		}
	}

	return defaultLimit
}

func isLimitFlag(arg string) bool {
	return arg == "--limit" || arg == "-l"
}

func emitJSONErrors(records []store.InternalErrorRecord) error {
	if records == nil {
		records = []store.InternalErrorRecord{}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(records)
}
