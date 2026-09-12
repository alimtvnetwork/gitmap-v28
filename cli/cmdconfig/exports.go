package cmdconfig

// RunExportConfig handles gitmap export-config.
func RunExportConfig(args []string) error {
	return runExportConfig(args)
}

// RunImportConfig handles gitmap import-config.
func RunImportConfig(args []string) error {
	return runImportConfig(args)
}
