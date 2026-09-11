package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runMkdir(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap mkdir [-p] [-f] <path>", "E9000")
	}

	createParents, createFiles, pathArg := parseMkdirArgs(args)
	if pathArg == "" {
		return apperror.NewSimple("Error: missing path argument", "E9000")
	}

	absPath, err := resolveMkdirAbsPath(pathArg)
	if err != nil {
		return apperror.WrapSimple(err, "Error resolving path:")
	}

	if err := createTargetDirectory(absPath, createParents, createFiles); err != nil {
		return apperror.WrapSimple(err, "Error creating directory:")
	}

	return nil
}

func resolveMkdirAbsPath(pathArg string) (string, error) {
	expanded := expandTilde(pathArg)
	return filepath.Abs(expanded)
}

func createTargetDirectory(absPath string, createParents, createFiles bool) error {
	if createFiles {
		return touchFile(absPath, createParents)
	}
	return makeDir(absPath, createParents)
}

func makeDir(absPath string, createParents bool) error {
	if !createParents {
		return makeDirSingle(absPath)
	}
	return makeDirDeep(absPath)
}

func makeDirSingle(absPath string) error {
	if err := os.Mkdir(absPath, 0755); err != nil {
		return err
	}
	fmt.Printf("  %s✓%s [DIR] Created: %s\n", constants.ColorGreen, constants.ColorReset, absPath)
	return nil
}

func makeDirDeep(absPath string) error {
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return err
	}
	fmt.Printf("  %s✓%s [DIR] Created deeply: %s\n", constants.ColorGreen, constants.ColorReset, absPath)
	return nil
}

func touchFile(absPath string, createParents bool) error {
	if err := touchFilePrepareParent(absPath, createParents); err != nil {
		return err
	}
	f, err := os.OpenFile(absPath, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_ = f.Close()
	fmt.Printf("  %s✓%s [FILE] Touched/Created: %s\n", constants.ColorGreen, constants.ColorReset, absPath)
	return nil
}

func touchFilePrepareParent(absPath string, createParents bool) error {
	if !createParents {
		return nil
	}
	parent := filepath.Dir(absPath)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	fmt.Printf("  %s✓%s [DIR] Ensured parent: %s\n", constants.ColorGreen, constants.ColorReset, parent)
	return nil
}

func parseMkdirArgs(args []string) (bool, bool, string) {
	createParents := false
	createFiles := false
	pathArg := ""

	for _, arg := range args {
		if arg == "-p" {
			createParents = true
		} else if arg == "-f" || arg == "--file" {
			createFiles = true
		} else {
			pathArg = arg
		}
	}
	return createParents, createFiles, pathArg
}
