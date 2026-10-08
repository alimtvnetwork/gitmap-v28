package cmdrun

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// RunErrorsCmd displays recent execution failures in an ANSI table.
func RunErrorsCmd(args []string) error {
	limit := 50
	records, err := QueryRecentRunErrors(limit)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Printf("  %sINFO%s No execution errors recorded.\n", constants.ColorCyan, constants.ColorReset)
		return nil
	}

	cfg := buildRunErrorsTableConfig(records)
	termtable.PrintTable(cfg)

	return nil
}

func buildRunErrorsTableConfig(records []RunErrorRecord) termtable.TableConfig {
	columns := []termtable.Column{
		{Title: "ID", MaxWidth: 22, Align: termtable.AlignLeft},
		{Title: "FILE", MaxWidth: 32, Align: termtable.AlignLeft},
		{Title: "INTERPRETER", MaxWidth: 14, Align: termtable.AlignLeft},
		{Title: "EXIT", MaxWidth: 6, Align: termtable.AlignRight},
		{Title: "DURATION", MaxWidth: 10, Align: termtable.AlignRight},
		{Title: "ERROR", MaxWidth: 40, Align: termtable.AlignLeft},
	}

	rows := make([]termtable.Row, 0, len(records))
	for _, r := range records {
		rows = append(rows, buildRunErrorRow(r))
	}

	return termtable.TableConfig{
		Columns:     columns,
		Rows:        rows,
		HeaderColor: constants.ColorCyan,
		BorderColor: constants.ColorDim,
		HasBorders:  true,
	}
}

func buildRunErrorRow(r RunErrorRecord) termtable.Row {
	durStr := fmt.Sprintf("%dms", r.DurationMs)
	exitStr := fmt.Sprintf("%d", r.ExitCode)

	return termtable.Row{
		Cells: []string{
			r.ErrorID,
			r.FilePath,
			r.Interpreter,
			exitStr,
			durStr,
			r.ErrorMessage,
		},
		Color: constants.ColorRed,
	}
}

// RunClearErrorsCmd purges recorded execution failures.
func RunClearErrorsCmd(args []string) error {
	if err := ClearRunErrors(); err != nil {
		return err
	}

	fmt.Printf("  %sSUCCESS%s Cleared all recorded run execution errors.\n", constants.ColorGreen, constants.ColorReset)

	return nil
}
