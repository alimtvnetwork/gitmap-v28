// Command ncgen regenerates cli/cmdpending/nc_catalog_gen.go at build time.
//
// It parses the dispatcher registration in ../cmd (dispatchEntry routing
// tables and cobra.Command literals) plus the Cmd* string constants in
// ../constants, then merges every registered command with the curated
// editorial metadata in nc_annotations.json.
//
// Names and aliases come from the registration, so the catalog cannot
// drift from the real CLI surface; descriptions, examples, versions, and
// categories are editorial and live in the annotations file. Entries whose
// base command is not found in the registration (subcommand forms such as
// "spec next", or not-yet-wired commands) keep their annotated alias and
// are reported on stderr — informational only, never fatal.
//
// Usage (invoked via go:generate from package cmdpending):
//
//	go run ./ncgen/main.go ./ncgen/registry.go ./ncgen/constants.go ./ncgen/annotations.go ./ncgen/emit.go
//
// NOTE: the directive lists the .go files explicitly (instead of the
// ./ncgen package path) because `go generate` exports GOROOT into the
// generator's environment, and on machines where GOROOT == GOPATH (like
// this one) that breaks `go run` package-path resolution. File arguments
// run in ad-hoc mode and skip package resolution. If you add a .go file
// to this directory, add it to the //go:generate line as well.
package main

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

const (
	cmdDir          = "../cmd"
	constantsDir    = "../constants"
	annotationsFile = "ncgen/nc_annotations.json"
	outputFile      = "nc_catalog_gen.go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ncgen:", err)
		os.Exit(1)
	}
}

func run() error {
	if !fsutil.FileOrDirExists(annotationsFile) {
		return apperror.NewNotFoundError(fmt.Sprintf("annotations file %s not found (run from cli/cmdpending)", annotationsFile))
	}

	consts, err := loadConstants(constantsDir)
	if err != nil {
		return err
	}

	reg, err := buildRegistry(cmdDir, consts)
	if err != nil {
		return err
	}

	anns, err := loadAnnotations(annotationsFile)
	if err != nil {
		return err
	}

	return emitCatalog(outputFile, anns, reg)
}
