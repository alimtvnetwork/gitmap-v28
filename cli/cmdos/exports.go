package cmdos

var (
	RunPowerNeverSleepFn func() error
	RunPowerSetFn        func([]string) error
	RunPowerResetFn      func() error
	CheckHelpFn          func(subcmd string, args []string)
)

// RunOS dispatches gitmap os subcommands through canonical CLI router.
func RunOS(args []string) error {
	return RunOSCLI(args)
}

// RunOSFixLink inspects and repairs broken symlinks and shared directories.
func RunOSFixLink(args []string) error {
	return runOSFixLink(args)
}

// RunChangePassword changes the OS user password across Windows, Linux, and macOS.
func RunChangePassword(args []string) error {
	return RunChangePasswordCLI(args)
}

// RunCPUCLI displays CPU model, cores, and current load percentage.
func RunCPUCLI(args []string) error {
	return runCPUCLI(args)
}

// RunCPU is an alias for RunCPUCLI.
func RunCPU(args []string) error {
	return runCPUCLI(args)
}

// RunRAMCLI displays physical memory capacity, usage, and percentage.
func RunRAMCLI(args []string) error {
	return runRAMCLI(args)
}

// RunRAM is an alias for RunRAMCLI.
func RunRAM(args []string) error {
	return runRAMCLI(args)
}

// RunDiskCLI displays disk volume capacities and utilization.
func RunDiskCLI(args []string) error {
	return runDiskCLI(args)
}

// RunDisk is an alias for RunDiskCLI.
func RunDisk(args []string) error {
	return runDiskCLI(args)
}
