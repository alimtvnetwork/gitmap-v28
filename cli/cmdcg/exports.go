package cmdcg

// CGOptions holds options for coding guideline commands.
type CGOptions = cgOptions

// RunCG executes coding guideline commands with args.
func RunCG(args []string) error {
	return runCG(args)
}

// ParseCGFlags parses CLI arguments into CGOptions.
func ParseCGFlags(args []string) CGOptions {
	return parseCGFlags(args)
}
