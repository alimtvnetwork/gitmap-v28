package cmdinstall

// runInstallAgyWithOpts routes Antigravity CLI (agy) installation directly into the unified Antigravity installer engine.
func runInstallAgyWithOpts(opts installOptions) error {
	return runInstallAntigravityWithOpts(opts)
}
