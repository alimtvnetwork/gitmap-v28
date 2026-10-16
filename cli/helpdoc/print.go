package helpdoc

import (
	"embed"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

//go:embed *.md
var files embed.FS

// Print reads and prints the help file for the given command using the
// PrettyOn mode to ensure modern terminal formatting is displayed.
// Kept as a thin wrapper for callers that don't parse a --pretty flag.
func Print(command string) {
	PrintWithMode(command, render.PrettyOn)
}

// PrintWithMode reads and prints the help file for `command`, routing the
// markdown through render.RenderANSI when render.Decide says so for the
// caller-supplied PrettyMode. This is the preferred entry point for
// command surfaces that parse --pretty / --no-pretty so user intent
// flows all the way to the renderer.
//
// The decision is delegated to render.Decide so help, templates show,
// and changelog all answer "should I emit ANSI?" identically.
func PrintWithMode(command string, mode render.PrettyModeType) {
	data, err := ReadRaw(command)
	if err != nil {
		cliexit.HandleError(wrapHelpUnavailableError(err, "helptext.PrintWithMode", "E1071", command), 1)

		return
	}

	if mode != render.PrettyOff {
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
		cliexit.HandleError(wrapHelpUnavailableError(err, "helptext.PrintRaw", "E1072", command), 1)

		return
	}

	fmt.Print(string(data))
}

func wrapHelpUnavailableError(err error, caller, code, command string) *apperror.AppError {
	return apperror.WrapWithDetails(err, caller, code, fmt.Sprintf("No help available for '%s'", command), "helptext", apperror.ErrorTypeValidation, apperror.SeverityError, map[string]any{"command": command})
}

// ReadRaw returns the embedded help markdown for the given command
// without exiting the process on miss. Test-friendly counterpart to
// PrintRaw — callers that want to assert on help contents (alias
// coverage, link integrity, etc.) should use this so the test binary
// is not torn down by os.Exit when a lookup fails.
func ReadRaw(command string) ([]byte, error) {
	command = resolveHelpAlias(command)
	data, err := files.ReadFile(command + ".md")
	if err == nil {
		return data, nil
	}

	low := strings.ToLower(command)
	data, err = files.ReadFile(low + ".md")
	if err == nil {
		return data, nil
	}

	dashed := strings.ReplaceAll(low, " ", "-")

	return files.ReadFile(dashed + ".md")
}

// HasTopic reports whether help documentation exists for the specified topic.
func HasTopic(command string) bool {
	_, err := ReadRaw(command)

	return err == nil
}
