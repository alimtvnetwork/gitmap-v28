// Package cmdzsh provides subcommand dispatching handlers.
package cmdzsh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func dispatchInstall(args []string) error {
	opts := ParseInstallOptions(args)
	appErr := InstallZshSuite(opts)
	if appErr != nil {
		return appErr
	}

	fmt.Printf("ZSH suite installed successfully with theme %s%s", opts.Theme, constants.NewLineUnix)

	return nil
}

func dispatchTheme(args []string) error {
	theme := ResolveThemeParam(args)
	home := GetFlagValue(args, "--home")
	isDryRun := HasFlag(args, "--dry-run")
	appErr := ApplyTheme(home, theme, isDryRun)
	if appErr != nil {
		return appErr
	}

	fmt.Printf("ZSH theme updated to %s%s", theme, constants.NewLineUnix)

	return nil
}

func dispatchClean(args []string) error {
	opts := ParseCleanOptions(args)
	appErr := CleanOhMyZsh(opts)
	if appErr != nil {
		return appErr
	}

	fmt.Printf("Oh-My-Zsh cleaned successfully%s", constants.NewLineUnix)

	return nil
}

func dispatchSwitch(args []string) error {
	opts := ParseSwitchOptions(args)
	appErr := ChangeDefaultShell(opts)
	if appErr != nil {
		return appErr
	}

	fmt.Printf("Default shell switched to ZSH%s", constants.NewLineUnix)

	return nil
}

func dispatchProfile(args []string) error {
	opts := ParseProfileOptions(args)
	appErr := SetupUserProfile(opts)
	if appErr != nil {
		return appErr
	}

	fmt.Printf("User workspace directories provisioned%s", constants.NewLineUnix)

	return nil
}

func dispatchStatus(args []string) error {
	home := GetFlagValue(args, "--home")
	res := InspectStatus(home)
	if res.IsFailure() {
		return res.AppError()
	}

	report := FormatStatusReport(res.Value)
	fmt.Print(report)

	return nil
}
