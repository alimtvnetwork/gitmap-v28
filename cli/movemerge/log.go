package movemerge

import (
	"fmt"
	"os"
)

// logf prints a single structured log line to stdout.
func logf(prefix, format string, a ...any) {
	fmt.Printf("%s %s\n", prefix, fmt.Sprintf(format, a...))
}

// logIndent prints an indented sub-line.
func logIndent(prefix, format string, a ...any) {
	fmt.Printf("%s   %s\n", prefix, fmt.Sprintf(format, a...))
}

// logErr writes a structured error line to stderr.
func logErr(prefix, msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, msg)
}
