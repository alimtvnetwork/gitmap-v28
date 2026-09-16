package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// HygieneFormatType is the value of --format; "" means default table output.
type HygieneFormatType string

type hygieneFormat = HygieneFormatType

const (
	HygieneFormatTypeTable HygieneFormatType = ""
	HygieneFormatTypeJSON  HygieneFormatType = "json"
	HygieneFormatTypeCSV   HygieneFormatType = "csv"
)

const (
	hygieneFormatTable = HygieneFormatTypeTable
	hygieneFormatJSON  = HygieneFormatTypeJSON
	hygieneFormatCSV   = HygieneFormatTypeCSV
)

// parseHygieneFormat normalizes and validates a --format value.
func parseHygieneFormat(s string) (HygieneFormatType, error) {
	switch s {
	case "", "table":
		return HygieneFormatTypeTable, nil
	case "json":
		return HygieneFormatTypeJSON, nil
	case "csv":
		return HygieneFormatTypeCSV, nil
	}

	return "", fmt.Errorf("invalid --format %q (want table|json|csv)", s)
}

// emitJSON writes v as indented JSON to stdout.
func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "  json encode: %v\n", err)
	}
}

// emitCSV writes a header + rows to stdout via encoding/csv.
func emitCSV(header []string, rows [][]string) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write(header)
	for _, r := range rows {
		_ = w.Write(r)
	}

	w.Flush()
}
