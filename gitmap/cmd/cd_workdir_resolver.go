// Package cmd — cd_workdir_resolver.go: resolves work directories for cd navigation.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func handleWorkDirOrNotFound(name string, rest []string) error {
	workPath, hasWorkDir := resolveCDWorkDirPath(name)
	if hasWorkDir {
		return dispatchCDWorkPath(workPath, rest)
	}

	fmt.Fprintf(os.Stderr, constants.ErrCDNotFound, name)
	cliexit.HandleError(nil, 1)

	return nil
}

func dispatchCDWorkPath(workPath string, rest []string) error {
	if len(rest) > 0 {
		return runCDInner(workPath, rest)
	}

	fmt.Print(workPath)
	WriteShellHandoff(workPath)
	warnIfNoWrapper()

	return nil
}

func resolveDefaultWorkDirPath() (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()

	wd, errGet := db.GetDefaultWorkDir()
	if errGet != nil || wd == nil || wd.AbsolutePath == "" {
		return "", false
	}

	return wd.AbsolutePath, true
}

func resolveCDWorkDirPath(name string) (string, bool) {
	lower := strings.ToLower(name)
	if isWorkDirKeyword(lower) {
		return resolveDefaultWorkDirPath()
	}

	return findWorkDirByNameOrLabel(name)
}

func isWorkDirKeyword(name string) bool {
	return name == "work" || name == "workdir" || name == "wd" || name == "default"
}

func findWorkDirByNameOrLabel(target string) (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil || len(dirs) == 0 {
		return "", false
	}

	for _, d := range dirs {
		if matchesWorkDir(d.AbsolutePath, d.Label, target) {
			return d.AbsolutePath, true
		}
	}

	return "", false
}

func matchesWorkDir(absPath, label, target string) bool {
	lowerTarget := strings.ToLower(target)
	if label != "" && strings.EqualFold(label, lowerTarget) {
		return true
	}

	baseName := strings.ToLower(filepath.Base(absPath))
	if baseName == lowerTarget {
		return true
	}

	return strings.Contains(strings.ToLower(absPath), lowerTarget)
}
