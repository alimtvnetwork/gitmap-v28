package cmdbackup

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"strconv"
	"strings"
)

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}
func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}
func pickProfileBySequenceOrName(
	profiles []model.GitProfile,
	val string,
) (int, model.GitProfile, error) {
	num, err := strconv.Atoi(val)

	if err == nil && num >= 1 && num <= len(profiles) {
		return num - 1, profiles[num-1], nil
	}

	for i, p := range profiles {
		if strings.EqualFold(p.Name, val) || strings.EqualFold(p.ID, val) {
			return i, p, nil
		}
	}

	return -1, model.GitProfile{}, apperror.NewNotFoundError(fmt.Sprintf("profile not found: %s", val))
}
