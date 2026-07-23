// Package formatter — text.go writes a plain text file with one git clone command per line.
package formatter

import (
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// WriteText writes one git clone instruction per line to a plain text file.
func WriteText(w io.Writer, records []model.ScanRecord) {
	for _, r := range records {
		fmt.Fprintln(w, r.CloneInstruction)
	}
}
