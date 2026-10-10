// Package cmdfixreleasetags provides flag parsing and validation for the
// release tags cleanup command.
package cmdfixreleasetags

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// FixReleaseTagsFlags defines parsed CLI options for release tag cleanup.
type FixReleaseTagsFlags struct {
	IsDryRun        bool   `json:"isDryRun"`
	IsConfirmed     bool   `json:"isConfirmed"`
	IsJSON          bool   `json:"isJSON"`
	IsLocalOnly     bool   `json:"isLocalOnly"`
	IsRemoteOnly    bool   `json:"isRemoteOnly"`
	IsVerbose       bool   `json:"isVerbose"`
	IsHelpRequested bool   `json:"isHelpRequested"`
	TargetDirectory string `json:"targetDirectory"`
}

// ParseFixReleaseTagsFlags parses arguments into a validated FixReleaseTagsFlags model.
func ParseFixReleaseTagsFlags(args []string) (FixReleaseTagsFlags, error) {
	flags := FixReleaseTagsFlags{
		TargetDirectory: ".",
	}

	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if arg == "" {
			continue
		}

		if isHelpFlag(arg) {
			flags.IsHelpRequested = true
			continue
		}

		if isDryRunFlag(arg) {
			flags.IsDryRun = true
			continue
		}

		if isConfirmFlag(arg) {
			flags.IsConfirmed = true
			continue
		}

		if isJSONFlag(arg) {
			flags.IsJSON = true
			continue
		}

		if isLocalOnlyFlag(arg) {
			flags.IsLocalOnly = true
			continue
		}

		if isRemoteOnlyFlag(arg) {
			flags.IsRemoteOnly = true
			continue
		}

		if isVerboseFlag(arg) {
			flags.IsVerbose = true
			continue
		}

		if isRepoFlag(arg) {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				flags.TargetDirectory = strings.TrimSpace(parts[1])
				continue
			}
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags.TargetDirectory = strings.TrimSpace(args[i+1])
				i++
				continue
			}
		}

		// Positional argument treated as target directory if not set
		if !strings.HasPrefix(arg, "-") && flags.TargetDirectory == "." {
			flags.TargetDirectory = arg
		}
	}

	if err := validateFlags(flags); err != nil {
		return flags, err
	}

	return flags, nil
}

// ParseFlags is an alias for ParseFixReleaseTagsFlags.
func ParseFlags(args []string) (FixReleaseTagsFlags, error) {
	return ParseFixReleaseTagsFlags(args)
}

func parseFlags(args []string) (FixReleaseTagsFlags, error) {
	return ParseFixReleaseTagsFlags(args)
}

func isHelpFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-h" || low == "--help" || low == "help"
}

func isDryRunFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-n" || low == "--dry-run" || low == "--dryrun" || low == "-dry-run"
}

func isConfirmFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-y" || low == "--yes" || low == "-yes" || low == "--confirm" || low == "-confirm"
}

func isJSONFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--json" || low == "-json"
}

func isLocalOnlyFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--local-only" || low == "--local"
}

func isRemoteOnlyFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--remote-only" || low == "--remote"
}

func isVerboseFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-v" || low == "--verbose"
}

func isRepoFlag(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-r" || low == "--repo" || strings.HasPrefix(low, "--repo=") || strings.HasPrefix(low, "-r=")
}

func validateFlags(flags FixReleaseTagsFlags) error {
	if flags.IsLocalOnly && flags.IsRemoteOnly {
		return apperror.NewSimple("cannot specify both --local-only and --remote-only flags simultaneously", "E1024")
	}

	if flags.TargetDirectory != "." && flags.TargetDirectory != "" {
		info, err := os.Stat(flags.TargetDirectory)
		if err != nil || !info.IsDir() {
			return apperror.NewSimple(fmt.Sprintf("specified repository directory does not exist or is inaccessible: %s", flags.TargetDirectory), "E1026")
		}
	}

	return nil
}
