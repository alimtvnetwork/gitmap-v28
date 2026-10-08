package cmdrepo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RunCreateFn is wired by cmd/di_hooks.go.
var RunCreateFn func(args []string) error

// RunRecreateFn is wired by cmd/di_hooks.go.
var RunRecreateFn func(args []string) error

// RunCreateLocalFn is wired by cmd/di_hooks.go.
var RunCreateLocalFn func(args []string) error

func runCreate(args []string) error {
	if RunCreateFn != nil {
		return RunCreateFn(args)
	}
	return nil
}

func runRecreate(args []string) error {
	if RunRecreateFn != nil {
		return RunRecreateFn(args)
	}
	return nil
}

func runCreateLocal(args []string) error {
	if RunCreateLocalFn != nil {
		return RunCreateLocalFn(args)
	}
	return nil
}

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || len(a) > len(flagName) && a[:len(flagName)+1] == flagName+"=" {
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
		if len(arg) > len(flagName) && arg[:len(flagName)+1] == flagName+"=" {
			return arg[len(flagName)+1:]
		}
	}
	return ""
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func pickProfileBySequenceOrName(profiles []model.GitProfile, val string) (int, model.GitProfile, error) {
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

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}
