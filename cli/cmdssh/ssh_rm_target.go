package cmdssh

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveRmTarget(args []string) string {
	for _, arg := range args {
		isAll := isAllFlag(arg)
		if isAll {
			return "all"
		}
	}

	for _, arg := range args {
		clean := strings.TrimSpace(arg)
		isValid := clean != "" && !strings.HasPrefix(clean, "-")
		if isValid {
			return clean
		}
	}

	return ""
}

func validateRmTarget(args []string) (string, error) {
	target := resolveRmTarget(args)
	hasTarget := target != ""
	if !hasTarget {
		return "", apperror.NewValidationError("missing machine alias, IP, or --all to remove")
	}

	return target, nil
}
