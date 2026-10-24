package cmdworkgroup

import (
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func tryExtractFlagValue(arg, flag string) (string, bool) {
	prefix := flag + "="
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix), true
	}

	return "", false
}

func parseKeyValueFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagSSHKey)); ok {
		opts.SSHKeyPath = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagSSHKey) && idx+1 < len(args) {
		opts.SSHKeyPath = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagDesc)); ok {
		opts.Description = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagDesc) && idx+1 < len(args) {
		opts.Description = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagHost)); ok {
		opts.Host = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagHost) && idx+1 < len(args) {
		opts.Host = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseUserHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagUserName)); ok {
		opts.UserName = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagUserName) && idx+1 < len(args) {
		opts.UserName = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagUserEmail)); ok {
		opts.UserEmail = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagUserEmail) && idx+1 < len(args) {
		opts.UserEmail = args[idx+1]

		return idx + 1
	}

	return parseHostFlag(args, idx, opts)
}

var workgroupJSONFlagTokens = []string{string(constants.WorkgroupFlagJSONShort), string(constants.WorkgroupFlagJSON)}

func parseBoolFlag(arg string, opts *WorkgroupOptions) bool {
	if slices.Contains(workgroupJSONFlagTokens, arg) {
		opts.IsJSON = true

		return true
	}
	if arg == string(constants.WorkgroupFlagDefault) {
		opts.IsDefault = true

		return true
	}

	return false
}

func parseFlagsLoop(args []string, opts *WorkgroupOptions) []string {
	pos := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if parseBoolFlag(args[i], opts) {
			continue
		}
		newI := parseKeyValueFlag(args, i, opts)
		if newI != i {
			i = newI

			continue
		}
		newI2 := parseUserHostFlag(args, i, opts)
		if newI2 != i {
			i = newI2

			continue
		}
		pos = append(pos, args[i])
	}

	return pos
}

func assignPositionalWorkgroup(pos []string, opts *WorkgroupOptions) {
	if len(pos) == 0 {
		return
	}
	opts.Subcommand = strings.ToLower(pos[0])
	if len(pos) > 1 {
		opts.Name = pos[1]
		opts.Path = pos[1]
	}
	if len(pos) > 2 {
		assignSecondPositional(pos[2], opts)
	}
}

func assignSecondPositional(val string, opts *WorkgroupOptions) {
	if opts.Subcommand == string(constants.WorkgroupCmdBind) && opts.SSHKeyPath == "" {
		opts.SSHKeyPath = val

		return
	}
	opts.Path = val
}

// ParseWorkgroupOptions extracts flags and positional arguments.
func ParseWorkgroupOptions(args []string) WorkgroupOptions {
	var opts WorkgroupOptions
	pos := parseFlagsLoop(args, &opts)
	assignPositionalWorkgroup(pos, &opts)

	return opts
}
