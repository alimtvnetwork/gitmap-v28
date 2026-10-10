package cmdscan

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhygiene"
	"os"
)

func parseHygieneFormat(s string) (cmdhygiene.HygieneFormatType, error) {
	switch s {
	case "", "table":
		return cmdhygiene.HygieneFormatTypeTable, nil
	case "json":
		return cmdhygiene.HygieneFormatTypeJSON, nil
	case "csv":
		return cmdhygiene.HygieneFormatTypeCSV, nil
	}

	return "", fmt.Errorf("invalid --format %q (want table|json|csv)", s)
}

func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "  json encode: %v\n", err)
	}
}
func emitCSV(header []string, rows [][]string) {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write(header)
	for _, r := range rows {
		_ = w.Write(r)
	}

	w.Flush()
}
