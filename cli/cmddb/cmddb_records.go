package cmddb

import (
	"fmt"
	"io"
	"strings"
)

const defaultRowLimit = 20

// RunDbList displays rows from a specified table with limit paging.
func RunDbList(tableArg, dbArg string, limit int, isJson bool, out io.Writer) error {
	cleanTable := strings.TrimSpace(tableArg)
	if cleanTable == "" {
		return fmt.Errorf("table name is required (e.g. gitmap db list <table-name>)")
	}

	effectiveLimit := resolveEffectiveLimit(limit)
	conn, dbPath, err := OpenResolvedDb(dbArg)
	if err != nil {
		return err
	}
	defer conn.Close()

	escaped := strings.ReplaceAll(cleanTable, `"`, `""`)
	query := fmt.Sprintf(`SELECT * FROM "%s" LIMIT %d`, escaped, effectiveLimit)

	res, queryErr := executeDynamicQuery(conn, query)
	if queryErr != nil {
		return fmt.Errorf("list rows from table '%s': %w", cleanTable, queryErr)
	}

	if isJson {
		return outputJson(out, res)
	}

	fmt.Fprintf(out, "Table: %s in %s (limit %d)\n\n", cleanTable, dbPath, effectiveLimit)
	renderQueryTerminal(out, res)

	return nil
}

func resolveEffectiveLimit(limit int) int {
	if limit <= 0 {
		return defaultRowLimit
	}

	return limit
}
