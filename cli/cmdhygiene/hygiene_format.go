package cmdhygiene

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene/enums"

// HygieneFormatType is the value of --format; "" means default table output.
// HygieneFormatType aliases the strongly-typed enum from the enums package.
type HygieneFormatType = enums.HygieneFormatType

const (
	HygieneFormatTypeTable = enums.HygieneFormatTypeTable
	HygieneFormatTypeJSON  = enums.HygieneFormatTypeJSON
	HygieneFormatTypeCSV   = enums.HygieneFormatTypeCSV
)

type hygieneFormat = enums.HygieneFormatType

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

	return "", apperror.NewValidationError(fmt.Sprintf("invalid --format %q (want table|json|csv)", s))
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
