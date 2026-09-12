package cmdos

var (
	RunPowerNeverSleepFn func() error
	RunPowerSetFn        func([]string) error
	RunPowerResetFn      func() error
	CheckHelpFn          func(subcmd string, args []string)
)

// RunOS dispatches gitmap os subcommands.
func RunOS(args []string) error {
	return runOS(args)
}

// RunOSFixLink inspects and repairs broken symlinks and shared directories.
func RunOSFixLink(args []string) error {
	return runOSFixLink(args)
}
