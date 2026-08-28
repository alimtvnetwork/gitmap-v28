package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// envNamePattern validates environment variable names.
var envNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// validateEnvName checks variable name is valid.
func validateEnvName(name string) {
	if name == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)
		apperror.NewSimple("fatal error", "E9000")
		return
	}

	if envNamePattern.MatchString(name) {
		return
	}

	apperror.NewSimple(constants.ErrEnvInvalidName, "E9000")
	return
}

// validateEnvValue checks value is provided.
func validateEnvValue(value string) {
	if value == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvValueRequired)
		apperror.NewSimple("fatal error", "E9000")
		return
	}
}

// validateEnvPathDir checks the directory exists.
func validateEnvPathDir(dir string) {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)
		apperror.NewSimple("fatal error", "E9000")
		return
	}

	_, err := os.Stat(dir)
	if err != nil {
		apperror.NewSimple(constants.ErrEnvPathNotExist, "E9000")
		return
	}
}
