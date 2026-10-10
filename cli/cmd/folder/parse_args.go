package folder

import (
	"path/filepath"
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
		case arg == "--tree" || arg == "-tree":
			opts.Format = FormatTree
			hasExplicitFormat = true
		case arg == "--md" || arg == "-md" || arg == "--markdown":
			opts.Format = FormatMd
			hasExplicitFormat = true
		case arg == "--json" || arg == "-json":
			opts.Format = FormatJson
			hasExplicitFormat = true
		case arg == "--yaml" || arg == "-yaml" || arg == "--yml":
			opts.Format = FormatYaml
			hasExplicitFormat = true
		case arg == "--flat" || arg == "-flat" || arg == "--list":
			opts.Format = FormatFlat
			hasExplicitFormat = true
		case arg == "--details" || arg == "-l" || arg == "-details":
			opts.IsDetailed = true
		case arg == "--only-text" || arg == "-only-text":
			opts.Filter.OnlyText = true
		case arg == "--only-binary" || arg == "-only-binary":
			opts.Filter.OnlyBinary = true
		case (arg == "--except" || arg == "--exclude" || arg == "-exclude" || arg == "-except") && i+1 < len(args):
			opts.Filter.ExceptGlobs = append(opts.Filter.ExceptGlobs, ParseExceptGlobs(args[i+1])...)
			i++
		case (arg == "--ext" || arg == "-ext") && i+1 < len(args):
			opts.Filter.Extensions = append(opts.Filter.Extensions, strings.Split(args[i+1], ",")...)
			i++
		case (arg == "--max-depth" || arg == "-max-depth") && i+1 < len(args):
			if val, errParse := strconv.Atoi(args[i+1]); errParse == nil {
				opts.Filter.MaxDepth = val
			}

			i++
		case (arg == "-o" || arg == "--out" || arg == "--output" || arg == "--file" || arg == "-file" || arg == "-f") && i+1 < len(args):
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
		if isGlobPattern(item) {
			opts.Filter.IncludeGlobs = append(opts.Filter.IncludeGlobs, item)
			if ext := extractGlobExtension(item); ext != "" {
				opts.Filter.Extensions = append(opts.Filter.Extensions, ext)
			}
		} else {
			dirCandidates = append(dirCandidates, item)
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
	if strings.HasPrefix(pattern, "*.") {
		ext := strings.TrimPrefix(pattern, "*.")
		if ext != "" && !strings.ContainsAny(ext, "*?") {
			return ext
		}
	}

	return ""
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
