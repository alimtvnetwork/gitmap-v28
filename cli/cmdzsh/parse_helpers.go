// Package cmdzsh provides CLI flag and argument parsing helpers.
package cmdzsh

import (
	"strings"
)

// HasFlag checks if an argument list contains a given flag name.
func HasFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

// GetFlagValue extracts the value for a given flag name from args.
func GetFlagValue(args []string, flagName string) string {
	prefix := flagName + "="
	for i, a := range args {
		if strings.HasPrefix(a, prefix) {
			return strings.TrimPrefix(a, prefix)
		}

		if a == flagName && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// ResolveThemeParam resolves theme name from positional argument or --theme flag.
func ResolveThemeParam(args []string) string {
	theme := GetFlagValue(args, "--theme")
	if theme != "" {
		return theme
	}

	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}

	return string(ZshThemeRobbyrussell)
}

// ParseInstallOptions extracts ZshOptions from command arguments.
func ParseInstallOptions(args []string) ZshOptions {
	theme := ResolveThemeParam(args)

	return ZshOptions{
		Theme:         theme,
		TargetHome:    GetFlagValue(args, "--home"),
		IsAppendZshrc: HasFlag(args, "--append-zshrc") || HasFlag(args, "--append"),
		IsInstallZsh:  true,
		IsChangeShell: HasFlag(args, "--change-shell"),
		IsDryRun:      HasFlag(args, "--dry-run"),
		IsUnattended:  HasFlag(args, "--unattended"),
	}
}

// ParseCleanOptions extracts ZshCleanOptions from command arguments.
func ParseCleanOptions(args []string) ZshCleanOptions {
	return ZshCleanOptions{
		Theme:       GetFlagValue(args, "--theme"),
		TargetHome:  GetFlagValue(args, "--home"),
		IsBackup:    true,
		IsReinstall: HasFlag(args, "--reinstall"),
		IsDryRun:    HasFlag(args, "--dry-run"),
	}
}

// ParseSwitchOptions extracts ZshSwitchOptions from command arguments.
func ParseSwitchOptions(args []string) ZshSwitchOptions {
	return ZshSwitchOptions{
		TargetUser: GetFlagValue(args, "--user"),
		IsDryRun:   HasFlag(args, "--dry-run"),
	}
}

// ParseProfileOptions extracts ZshProfileOptions from command arguments.
func ParseProfileOptions(args []string) ZshProfileOptions {
	return ZshProfileOptions{
		TargetHome:        GetFlagValue(args, "--home"),
		AuthorizedKeyFile: GetFlagValue(args, "--keys"),
		IsDryRun:          HasFlag(args, "--dry-run"),
	}
}
