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
func validateEnvName(name string) *apperror.AppError {
	if name == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)

		return apperror.NewSimple(constants.ErrEnvNameRequired, "E9000")
	}

	if envNamePattern.MatchString(name) {
		return nil
	}

	return apperror.NewSimple(constants.ErrEnvInvalidName, "E9000")
}

// validateEnvValue checks value is provided.
func validateEnvValue(value string) *apperror.AppError {
	if value == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvValueRequired)

		return apperror.NewSimple(constants.ErrEnvValueRequired, "E9000")
	}

	return nil
}

// validateEnvPathDir checks the directory exists.
func validateEnvPathDir(dir string) *apperror.AppError {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)

		return apperror.NewSimple(constants.ErrEnvPathRequired, "E9000")
	}

	_, err := os.Stat(dir)
	if err != nil {
		return apperror.NewSimple(constants.ErrEnvPathNotExist, "E9000")
	}

	return nil
}
