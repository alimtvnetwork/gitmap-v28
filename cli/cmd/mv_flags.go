package cmd

type moveOpts struct {
	yes           bool
	dryRun        bool
	isSkipVSCode  bool
	isSkipDesktop bool
}

func parseMoveFlags(args []string) (moveOpts, []string) {
	opts := moveOpts{}
	var positional []string
	for _, a := range args {
		switch a {
		case "-y", "--yes":
			opts.yes = true
		case "--dry-run":
			opts.dryRun = true
		case "--no-vscode":
			opts.isSkipVSCode = true
		case "--no-desktop":
			opts.isSkipDesktop = true
		default:
			positional = append(positional, a)
		}
	}

	return opts, positional
}
