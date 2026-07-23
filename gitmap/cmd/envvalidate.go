package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// envNamePattern validates environment variable names.
var envNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// validateEnvName checks variable name is valid.
func validateEnvName(name string) {
	if name == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)
		os.Exit(1)
	}

	if envNamePattern.MatchString(name) {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrEnvInvalidName, name)
	os.Exit(1)
}

// validateEnvValue checks value is provided.
func validateEnvValue(value string) {
	if value == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvValueRequired)
		os.Exit(1)
	}
}

// validateEnvPathDir checks the directory exists.
func validateEnvPathDir(dir string) {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)
		os.Exit(1)
	}

	_, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrEnvPathNotExist, dir)
		os.Exit(1)
	}
}
