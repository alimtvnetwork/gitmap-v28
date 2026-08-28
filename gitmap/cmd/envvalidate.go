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
		return apperror.New("fatal error", "E9000", nil)
	}

	if envNamePattern.MatchString(name) {
		return
	}

	return apperror.New(constants.ErrEnvInvalidName, "E9000", nil)
}

// validateEnvValue checks value is provided.
func validateEnvValue(value string) {
	if value == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvValueRequired)
		return apperror.New("fatal error", "E9000", nil)
	}
}

// validateEnvPathDir checks the directory exists.
func validateEnvPathDir(dir string) {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)
		return apperror.New("fatal error", "E9000", nil)
	}

	_, err := os.Stat(dir)
	if err != nil {
		return apperror.New(constants.ErrEnvPathNotExist, "E9000", nil)
	}
}
