package cmdvhost

// VHostCmd runs the main vhost command.
func VHostCmd(args []string) error {
	return runVHost(args)
}

// RunVHost runs the vhost CLI entrypoint.
func RunVHost(args []string) error {
	return runVHost(args)
}

// RunVHostCreate runs the vhost create CLI entrypoint.
func RunVHostCreate(args []string) error {
	return runVHostCreate(args)
}

// RunVHostEnable runs the vhost enable CLI entrypoint.
func RunVHostEnable(args []string) error {
	return runVHostEnable(args)
}

// RunVHostDisable runs the vhost disable CLI entrypoint.
func RunVHostDisable(args []string) error {
	return runVHostDisable(args)
}

// RunVHostRm runs the vhost rm CLI entrypoint.
func RunVHostRm(args []string) error {
	return runNginxRm(args)
}

// RunVHostTest runs the vhost test CLI entrypoint.
func RunVHostTest(args []string) error {
	return runVHostTest(args)
}

// RunVHostReload runs the vhost reload CLI entrypoint.
func RunVHostReload(args []string) error {
	return runVHostReload(args)
}

// RunVHostList runs the vhost list CLI entrypoint.
func RunVHostList(args []string) error {
	return runVHostList(args)
}
