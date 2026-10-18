// Package heatmap — heatmap_query_render.go formats query results as table, JSON, or YAML.
package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"gopkg.in/yaml.v3"
)

func renderQueryResult(dbPath, query string, cols []string, rows [][]any, durationMs float64, format string, out io.Writer) error {
	switch format {
	case "json":
		return renderQueryJSON(dbPath, query, cols, rows, durationMs, out)
	case "yaml", "yml":
		return renderQueryYAML(dbPath, query, cols, rows, durationMs, out)
	default:
		return renderQueryTable(cols, rows, durationMs, out)
	}
}

func renderQueryTable(cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	if len(rows) == 0 {
		fmt.Fprintf(out, "Empty set (%.2fms)\n", durationMs)
		return nil
	}
	var tCols []termout.Column
	for _, c := range cols {
		tCols = append(tCols, termout.Column{Title: strings.ToUpper(c), MinWidth: 10, MaxWidth: 35})
	}
	var tRows []termout.Row
	for _, row := range rows {
		cells := make([]string, len(cols))
		for i := range cols {
			cells[i] = formatCellValue(row[i])
		}
		tRows = append(tRows, termout.Row{Cells: cells})
	}
	fmt.Fprint(out, termout.RenderTable(termout.TableConfig{Columns: tCols, Rows: tRows}))
	fmt.Fprintf(out, "\n%d row(s) in set (%.2fms)\n", len(rows), durationMs)
	return nil
}

func renderQueryJSON(dbPath, query string, cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	absDb, relDb := resolveDbPathDetails(dbPath)
	payload := map[string]any{
		"database":       relDb,
		"relativeDbPath": relDb,
		"absoluteDbPath": absDb,
		"repoSecretsUrl": ResolveRepoSecretsURL(),
		"query":          query,
		"durationMs":     durationMs,
		"columns":        cols,
		"rowCount":       len(rows),
		"rows":           buildRowMaps(cols, rows),
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func renderQueryYAML(dbPath, query string, cols []string, rows [][]any, durationMs float64, out io.Writer) error {
	absDb, relDb := resolveDbPathDetails(dbPath)
	payload := map[string]any{
		"database":       relDb,
		"relativeDbPath": relDb,
		"absoluteDbPath": absDb,
		"repoSecretsUrl": ResolveRepoSecretsURL(),
		"query":          query,
		"durationMs":     durationMs,
		"columns":        cols,
		"rowCount":       len(rows),
		"rows":           buildRowMaps(cols, rows),
	}
	bytes, err := yaml.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = out.Write(bytes)
	return err
}

func buildRowMaps(cols []string, rows [][]any) []map[string]any {
	maps := make([]map[string]any, len(rows))
	for i, r := range rows {
		m := make(map[string]any, len(cols))
		for j, c := range cols {
			m[c] = r[j]
		}
		maps[i] = m
	}
	return maps
}

func formatCellValue(val any) string {
	if val == nil {
		return "NULL"
	}
	if b, ok := val.([]byte); ok {
		return string(b)
	}
	return fmt.Sprintf("%v", val)
}
