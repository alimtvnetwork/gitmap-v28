// Package heatmap — heatmap_view_render.go formats tables, JSON, and YAML output for heatmap view.
package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"gopkg.in/yaml.v3"
)

func marshalAndWriteYAML(out io.Writer, payload any) error {
	b, err := yaml.Marshal(payload)
	if err != nil {
		return err
	}
	_, writeErr := out.Write(b)
	return writeErr
}

func formatOutput(out io.Writer, format string, payload any, printTable func()) error {
	if format == "json" {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}
	if format == "yaml" || format == "yml" {
		return marshalAndWriteYAML(out, payload)
	}
	printTable()
	return nil
}

func formatTablesOutput(out io.Writer, dbPath string, tables []TableInfo, format string) error {
	absDb, relDb := resolveDbPathDetails(dbPath)
	payload := map[string]any{
		"database":       relDb,
		"relativeDbPath": relDb,
		"absoluteDbPath": absDb,
		"repoSecretsUrl": ResolveRepoSecretsURL(),
		"totalTables":    len(tables),
		"tables":         tables,
	}
	return formatOutput(out, format, payload, func() {
		renderTablesOutput(out, relDb, tables)
	})
}

func renderTablesOutput(out io.Writer, relDb string, tables []TableInfo) {
	fmt.Fprintf(out, "Heatmap Database Tables: %s\n\n", relDb)
	rows := make([]termout.Row, len(tables))
	for i, t := range tables {
		rows[i] = termout.Row{Cells: []string{t.Name, t.Type, strconv.Itoa(t.ColumnCount), strconv.Itoa(t.RowCount)}}
	}
	cols := []termout.Column{
		{Title: "NAME", MinWidth: 26},
		{Title: "TYPE", MinWidth: 8},
		{Title: "COLUMNS", MinWidth: 8, Align: termout.AlignRight},
		{Title: "ROWS", MinWidth: 8, Align: termout.AlignRight},
	}
	fmt.Fprint(out, termout.RenderTable(termout.TableConfig{Columns: cols, Rows: rows}))
}

func formatViewOutput(out io.Writer, dbPath, table string, total, limit int, cols []string, records []map[string]any, format string) error {
	absDb, relDb := resolveDbPathDetails(dbPath)
	payload := map[string]any{
		"database":       relDb,
		"relativeDbPath": relDb,
		"absoluteDbPath": absDb,
		"repoSecretsUrl": ResolveRepoSecretsURL(),
		"table":          table,
		"totalCount":     total,
		"returnedCount":  len(records),
		"limit":          limit,
		"columns":        cols,
		"records":        records,
	}
	return formatOutput(out, format, payload, func() {
		renderViewOutput(out, relDb, table, total, cols, records)
	})
}

func renderViewOutput(out io.Writer, relDb, table string, total int, cols []string, records []map[string]any) {
	if len(records) == 0 {
		fmt.Fprintf(out, "No records found in table '%s' matching criteria.\n", table)
		return
	}
	fmt.Fprintf(out, "Table: %s (Showing %d of %d records, DB: %s)\n\n", table, len(records), total, relDb)
	renderViewTableRows(out, cols, records)
}

func renderViewTableRows(out io.Writer, cols []string, records []map[string]any) {
	disp := cols
	if len(cols) > 6 {
		disp = cols[:6]
	}
	var tCols []termout.Column
	for _, c := range disp {
		tCols = append(tCols, termout.Column{Title: strings.ToUpper(c), MaxWidth: 30, MinWidth: 10})
	}
	var rows []termout.Row
	for _, rec := range records {
		cells := make([]string, len(disp))
		for i, c := range disp {
			cells[i] = formatCell(rec[c])
		}
		rows = append(rows, termout.Row{Cells: cells})
	}
	fmt.Fprint(out, termout.RenderTable(termout.TableConfig{Columns: tCols, Rows: rows}))
}

func formatCell(v any) string {
	s := fmt.Sprintf("%v", v)
	if len(s) > 28 {
		return s[:25] + "..."
	}
	return s
}

func resolveFormat(flags map[string]string) string {
	f := strings.ToLower(strings.TrimSpace(flags["format"]))
	if flags["json"] == "true" || flags["j"] == "true" || f == "json" {
		return "json"
	}
	if flags["yaml"] == "true" || flags["y"] == "true" || f == "yaml" || f == "yml" {
		return "yaml"
	}
	return "table"
}
