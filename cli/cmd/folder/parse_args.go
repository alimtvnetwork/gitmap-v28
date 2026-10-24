package folder

import (
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// ParseArgs parses CLI arguments into strongly typed Options.
func ParseArgs(args []string) (Options, error) {
	opts := Options{
		TargetDir: ".",
		Format:    FormatTree,
		Filter:    FilterConfig{},
	}

	var hasExplicitFormat bool
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case slices.Contains([]string{"--tree", "-tree"}, arg):
			opts.Format = FormatTree
			hasExplicitFormat = true
		case slices.Contains([]string{"--md", "-md", "--markdown"}, arg):
			opts.Format = FormatMd
			hasExplicitFormat = true
		case slices.Contains([]string{"--json", "-json"}, arg):
			opts.Format = FormatJson
			hasExplicitFormat = true
		case slices.Contains([]string{"--yaml", "-yaml", "--yml"}, arg):
			opts.Format = FormatYaml
			hasExplicitFormat = true
		case slices.Contains([]string{"--flat", "-flat", "--list"}, arg):
			opts.Format = FormatFlat
			hasExplicitFormat = true
		case slices.Contains([]string{"--details", "-l", "-details"}, arg):
			opts.IsDetailed = true
		case slices.Contains([]string{"--only-text", "-only-text"}, arg):
			opts.Filter.OnlyText = true
		case slices.Contains([]string{"--only-binary", "-only-binary"}, arg):
			opts.Filter.OnlyBinary = true
		case slices.Contains([]string{"--except", "--exclude", "-exclude", "-except"}, arg) && i+1 < len(args):
			opts.Filter.ExceptGlobs = append(opts.Filter.ExceptGlobs, ParseExceptGlobs(args[i+1])...)
			i++
		case slices.Contains([]string{"--ext", "-ext"}, arg) && i+1 < len(args):
			opts.Filter.Extensions = append(opts.Filter.Extensions, strings.Split(args[i+1], ",")...)
			i++
		case slices.Contains([]string{"--max-depth", "-max-depth"}, arg) && i+1 < len(args):
			if val, errParse := strconv.Atoi(args[i+1]); errParse == nil {
				opts.Filter.MaxDepth = val
			}

			i++
		case slices.Contains([]string{"-o", "--out", "--output", "--file", "-file", "-f"}, arg) && i+1 < len(args):
			opts.OutFile = args[i+1]
			if !hasExplicitFormat {
				opts.Format = deduceFormatFromExt(opts.OutFile, opts.Format)
			}
			i++
		case !strings.HasPrefix(arg, "-"):
			nonFlags = append(nonFlags, arg)
		}
	}

	return resolveNonFlags(opts, nonFlags)
}

func resolveNonFlags(opts Options, nonFlags []string) (Options, error) {
	var dirCandidates []string

	for _, item := range nonFlags {
		if !isGlobPattern(item) {
			dirCandidates = append(dirCandidates, item)
			continue
		}

		opts.Filter.IncludeGlobs = append(opts.Filter.IncludeGlobs, item)
		ext := extractGlobExtension(item)
		if ext != "" {
			opts.Filter.Extensions = append(opts.Filter.Extensions, ext)
		}
	}

	if len(dirCandidates) > 0 {
		opts.TargetDir = dirCandidates[0]
	}

	if len(dirCandidates) > 1 && opts.OutFile == "" {
		opts.OutFile = dirCandidates[1]
		opts.Format = deduceFormatFromExt(opts.OutFile, opts.Format)
	}

	return opts, nil
}

func isGlobPattern(s string) bool {
	return strings.ContainsAny(s, "*?")
}

func extractGlobExtension(pattern string) string {
	if !strings.HasPrefix(pattern, "*.") {
		return ""
	}

	ext := strings.TrimPrefix(pattern, "*.")
	if ext == "" || strings.ContainsAny(ext, "*?") {
		return ""
	}

	return ext
}

func deduceFormatFromExt(outFile string, fallback OutputFormat) OutputFormat {
	ext := strings.ToLower(filepath.Ext(outFile))
	switch ext {
	case ".json":
		return FormatJson
	case ".yaml", ".yml":
		return FormatYaml
	case ".md":
		return FormatMd
	case ".txt", ".tree":
		if fallback == FormatFlat {
			return FormatFlat
		}

		return FormatTree
	default:
		return fallback
	}
}
