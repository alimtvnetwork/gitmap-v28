// Package cmd — seowritecsv.go handles CSV parsing for seo-write.
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// loadCSVMessages reads commit messages from a CSV file.
func loadCSVMessages(path string) []commitMessage {
	records := readCSVFile(path)
	if len(records) == 0 {
		appErr := apperror.NewWithDetails(
			"seo.loadCSV",
			"E2019",
			constants.ErrSEOCSVEmpty,
			"cmd.seowritecsv",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
		cliexit.HandleError(appErr, 1)
	}

	return csvToMessages(records)
}

// readCSVFile opens and parses the CSV file.
func readCSVFile(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(constants.ErrSEOCSVRead)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOCSVRead, path, err)
		_ = f.Close()
		exitWith(1)
	}

	return records
}

// csvToMessages converts CSV rows into commit message pairs.
func csvToMessages(records [][]string) []commitMessage {
	var messages []commitMessage

	for _, row := range records {
		if len(row) < 2 {
			continue
		}
		messages = append(messages, commitMessage{
			title:       row[0],
			description: row[1],
		})
	}

	return messages
}
