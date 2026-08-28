package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runMkdir(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap mkdir [-p] <path>", "E9000")
	}

	createParents := false
	pathArg := ""

	if args[0] == "-p" {
		createParents = true
		if len(args) > 1 {
			pathArg = args[1]
		}
	} else {
		pathArg = args[0]
	}

	if pathArg == "" {
		return apperror.NewSimple("Error: missing path argument", "E9000")
	}

	absPath, err := filepath.Abs(pathArg)
	if err != nil {
		return apperror.WrapSimple(err, "Error resolving path:")
	}

	if createParents {
		err = os.MkdirAll(absPath, 0755)
	} else {
		err = os.Mkdir(absPath, 0755)
	}

	if err != nil {
		return apperror.WrapSimple(err, "Error creating directory:")
	}
	fmt.Printf("Created directory: %s\n", absPath)
	return nil
}
