package cmdinstall

import (
	"io"
	"os"
	"strings"
)

// reorderFlagsBeforeArgs moves flag-like arguments before positional arguments.
func reorderFlagsBeforeArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
			continue
		}

		flags = append(flags, arg)
		if hasTrailingFlagValue(args, arg, i) {
			flags = append(flags, args[i+1])
			i++
		}
	}

	return append(flags, positional...)
}

func hasTrailingFlagValue(args []string, arg string, i int) bool {
	if strings.Contains(arg, "=") || i+1 >= len(args) {
		return false
	}

	return !strings.HasPrefix(args[i+1], "-")
}

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func checkHelp(cmd string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(cmd, args)
	}
}

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}
