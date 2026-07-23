package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// runVSCodePMPath implements `gitmap vscode-pm-path` (alias `vpath`).
//
// It exercises the cross-platform resolver and prints exactly one line:
//
//   - The fully-resolved projects.json path on success.
//   - A sentinel-mapped diagnostic line on stderr + non-zero exit when the
//     user-data root or extension storage dir is missing.
//
// This is the user-facing way to confirm where gitmap will write
// projects.json before running `scan` or `code`.
func runVSCodePMPath(args []string) {
	checkHelp("vscode-pm-path", args)

	path, err := vscodepm.ProjectsJSONPath()
	if err != nil {
		printVSCodePMPathError(path, err)
		os.Exit(1)
	}

	fmt.Println(path)
}

// printVSCodePMPathError surfaces the two soft-fail sentinels with
// actionable hints. Anything else is forwarded verbatim.
func printVSCodePMPathError(path string, err error) {
	switch {
	case errors.Is(err, vscodepm.ErrUserDataMissing):
		fmt.Fprintf(os.Stderr, "%s\n", constants.MsgVSCodePMPathRootMissing)
	case errors.Is(err, vscodepm.ErrExtensionMissing):
		fmt.Fprintf(os.Stderr, constants.MsgVSCodePMPathExtMissing, path)
	default:
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
