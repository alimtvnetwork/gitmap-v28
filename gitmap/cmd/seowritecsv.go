// Package cmd — seowritecsv.go handles CSV parsing for seo-write.
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// loadCSVMessages reads commit messages from a CSV file.
func loadCSVMessages(path string) []commitMessage {
	records := readCSVFile(path)
	if len(records) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrSEOCSVEmpty)
		os.Exit(1)
	}

	return csvToMessages(records)
}

// readCSVFile opens and parses the CSV file.
func readCSVFile(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOCSVRead, path, err)
		os.Exit(1)
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
