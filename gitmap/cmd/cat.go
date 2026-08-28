package cmd

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runCat(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap cat <file>", "E9000")
	}

	file, err := os.Open(args[0])
	if err != nil {
		return apperror.WrapSimple(err, "Error opening file:")
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return apperror.WrapSimple(err, "Error reading file:")
	}
	return nil
}
