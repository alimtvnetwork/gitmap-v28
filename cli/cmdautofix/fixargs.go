package cmdautofix

import "strings"

// Argument helpers for `gitmap fix`: subcommand splitting, interspersed
// flag reordering, and --ext/--category normalization.

// reorderFixArgs moves flag tokens ahead of positional args so flags work
// in any position after the subcommand (`fix all [path] [flags]`).
// valueFlags are the flags that consume the following token.
func reorderFixArgs(args []string) []string {
	valueFlags := map[string]bool{
		"-w": true, "--workers": true, "--ext": true, "--uri-pattern": true,
	}
	flags := []string{}
	positional := []string{}
	i := 0
	for i < len(args) {
		a := args[i]
		if !strings.HasPrefix(a, "-") || a == "-" {
			positional = append(positional, a)
			i++
			continue
		}
		name := a
		if eq := strings.IndexByte(a, '='); eq >= 0 {
			name = a[:eq]
		}
		flags = append(flags, a)
		i++
		if valueFlags[name] && !strings.Contains(a, "=") && i < len(args) {
			flags = append(flags, args[i])
			i++
		}
	}
	return append(flags, positional...)
}

// splitSubcommand separates the first positional arg (the subcommand) from
// the rest. A leading flag means no subcommand was given.
func splitSubcommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	if strings.HasPrefix(args[0], "-") {
		return "", args
	}
	return args[0], args[1:]
}

func isFixSubcommand(name string) bool {
	for _, s := range fixSubcommands {
		if s == name {
			return true
		}
	}
	return false
}

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h"
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if isHelpFlag(a) {
			return true
		}
	}
	return false
}

// renderCategoryOrParentHelp handles `gitmap fix --help <sub>` ordering.
func renderCategoryOrParentHelp(rest []string) {
	if len(rest) > 0 && isFixSubcommand(rest[0]) {
		RenderFixCategoryHelp(rest[0])
		return
	}
	RenderFixHelp()
}

// normalizeExts lowercases extensions and ensures a leading dot.
func normalizeExts(raw string) []string {
	parts := splitCSV(raw)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(p)
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		out = append(out, p)
	}
	return out
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
