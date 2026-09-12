package cmdzip

// ZipCmd runs the main zip command.
func ZipCmd(args []string) error {
	return runZip(args)
}

// UnzipCmd runs the main unzip command.
func UnzipCmd(args []string) error {
	return runUnzipCompact(args)
}

// RunZip runs the zip CLI entrypoint.
func RunZip(args []string) error {
	return runZip(args)
}

// RunUnzipCompact runs the unzip CLI entrypoint.
func RunUnzipCompact(args []string) error {
	return runUnzipCompact(args)
}

// RunZipGroup runs the zip-group CLI entrypoint.
func RunZipGroup(args []string) error {
	return runZipGroup(args)
}
