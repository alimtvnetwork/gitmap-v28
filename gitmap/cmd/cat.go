package cmd

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runCat(args []string) error {
	if len(args) == 0 {
		return apperror.New("Usage: gitmap cat <file>", "E9000", nil)
	}

	file, err := os.Open(args[0])
	if err != nil {
		return apperror.Wrap(err, "Error opening file:", nil)
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return apperror.Wrap(err, "Error reading file:", nil)
	}
	return nil
}
