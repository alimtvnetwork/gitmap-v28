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
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--tree" || arg == "-tree":
			opts.Format = FormatTree
		case arg == "--md" || arg == "-md" || arg == "--markdown":
			opts.Format = FormatMd
		case arg == "--json" || arg == "-json":
			opts.Format = FormatJson
		case arg == "--yaml" || arg == "-yaml" || arg == "--yml":
			opts.Format = FormatYaml
		case arg == "--flat" || arg == "-flat" || arg == "--list":
			opts.Format = FormatFlat
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
		case (arg == "-o" || arg == "--out") && i+1 < len(args):
			opts.OutFile = args[i+1]
			i++
		case !strings.HasPrefix(arg, "-"):
			nonFlags = append(nonFlags, arg)
		}
	}

	return resolveNonFlags(opts, nonFlags)
}

func resolveNonFlags(opts Options, nonFlags []string) (Options, error) {
	if len(nonFlags) > 0 {
		opts.TargetDir = nonFlags[0]
	}
	if len(nonFlags) > 1 {
		opts.OutFile = nonFlags[1]
		opts.Format = deduceFormatFromExt(opts.OutFile, opts.Format)
	}
	return opts, nil
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
	case ".txt":
		return FormatFlat
	default:
		return fallback
	}
}
