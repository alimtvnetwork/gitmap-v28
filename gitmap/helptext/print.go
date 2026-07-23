package helptext

import (
	"embed"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

//go:embed *.md
var files embed.FS

// Print reads and prints the help file for the given command using the
// default PrettyAuto mode (TTY auto-detect + GITMAP_NO_PRETTY opt-out).
// Kept as a thin wrapper for callers that don't parse a --pretty flag.
func Print(command string) {
	PrintWithMode(command, render.PrettyAuto)
}

// PrintWithMode reads and prints the help file for `command`, routing the
// markdown through render.RenderANSI when render.Decide says so for the
// caller-supplied PrettyMode. This is the preferred entry point for
// command surfaces that parse --pretty / --no-pretty so user intent
// flows all the way to the renderer.
//
// The decision is delegated to render.Decide so help, templates show,
// and changelog all answer "should I emit ANSI?" identically.
func PrintWithMode(command string, mode render.PrettyMode) {
	data, err := files.ReadFile(command + ".md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No help available for '%s'\n", command)
		os.Exit(1)
	}

	if render.Decide(mode, render.StdoutIsTerminal(), true) {
		fmt.Print(render.RenderANSI(string(data)))

		return
	}

	fmt.Print(string(data))
}

// PrintRaw bypasses the pretty renderer and prints the embedded
// markdown verbatim. Useful for callers that pipe help into a pager
// or other tooling that handles its own formatting. Equivalent to
// PrintWithMode(command, render.PrettyOff) but spelled out for clarity
// at call sites that always want raw output regardless of TTY state.
func PrintRaw(command string) {
	data, err := ReadRaw(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No help available for '%s'\n", command)
		os.Exit(1)
	}
	fmt.Print(string(data))
}

// ReadRaw returns the embedded help markdown for the given command
// without exiting the process on miss. Test-friendly counterpart to
// PrintRaw — callers that want to assert on help contents (alias
// coverage, link integrity, etc.) should use this so the test binary
// is not torn down by os.Exit when a lookup fails.
func ReadRaw(command string) ([]byte, error) {
	return files.ReadFile(command + ".md")
}
