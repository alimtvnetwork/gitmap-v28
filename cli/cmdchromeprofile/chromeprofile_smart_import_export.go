package cmdchromeprofile

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type ChromePreviewOutputParams struct {
	Candidates []DiscoveredProfileCandidate
	Target     string
	IsJSON     bool
	FilePath   string
	TempFile   string
	Fnf        bool
}

func parseOptionExport(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if parseOptionFile(arg, args, i, opts) {
		return true
	}

	if parseOptionTempFile(arg, args, i, opts) {
		return true
	}

	if parseOptionFnf(arg, args, i, opts) {
		return true
	}

	return parseOptionJSON(arg, opts)
}

func parseOptionFile(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--file=") {
		opts.FilePath = strings.TrimPrefix(arg, "--file=")

		return true
	}

	if isFlag := isFileFlagToken(arg); isFlag && *i+1 < len(args) {
		*i++
		opts.FilePath = args[*i]

		return true
	}

	return false
}

func isFileFlagToken(arg string) bool {
	return slices.Contains([]string{"--file", "-f", "-o", "--output", "--out"}, arg)
}

func parseOptionTempFile(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--tempfile=") {
		opts.TempFile = strings.TrimPrefix(arg, "--tempfile=")

		return true
	}

	if isTemp := isTempFlagToken(arg); isTemp && *i+1 < len(args) {
		*i++
		opts.TempFile = args[*i]

		return true
	}

	return false
}

func isTempFlagToken(arg string) bool {
	return slices.Contains([]string{"--tempfile", "--temp"}, arg)
}

func parseOptionFnf(arg string, args []string, i *int, opts *chromeTransferOptions) bool {
	if strings.HasPrefix(arg, "--fnf=") {
		opts.FilePath = strings.TrimPrefix(arg, "--fnf=")
		opts.Fnf = true

		return true
	}

	if arg != "--fnf" && arg != "--fail-not-found" {
		return false
	}

	opts.Fnf = true
	assignFnfFilePath(args, i, opts)

	return true
}

func assignFnfFilePath(args []string, i *int, opts *chromeTransferOptions) {
	if *i+1 < len(args) && !strings.HasPrefix(args[*i+1], "-") {
		*i++
		opts.FilePath = args[*i]
	}
}

func parseOptionJSON(arg string, opts *chromeTransferOptions) bool {
	if slices.Contains([]string{"--json", "-j"}, arg) {
		opts.IsJSON = true

		return true
	}

	return false
}

func isValuedFlag(arg string) bool {
	prefixes := []string{
		"--file", "-f", "-o", "--output", "--out",
		"--tempfile", "--temp", "--fnf", "--limit", "-n", "-l",
		"--except", "--exclude", "--skip", "--email", "--format", "--profile", "-p",
	}

	for _, p := range prefixes {
		if arg == p || strings.HasPrefix(arg, p+"=") {
			return true
		}
	}

	return false
}

func dispatchPreviewOutput(params ChromePreviewOutputParams) error {
	destPath := resolvePreviewDestPath(params)
	if len(destPath) > 0 {
		return writePreviewToFile(destPath, params)
	}

	if params.IsJSON {
		return renderProfileCandidatesJSON(params.Candidates)
	}

	renderProfileCandidatesTable(params.Target, params.Candidates)

	return nil
}

func resolvePreviewDestPath(params ChromePreviewOutputParams) string {
	if len(params.TempFile) > 0 {
		return filepath.Join(resolveTempDir(), params.TempFile)
	}

	return params.FilePath
}

func writePreviewToFile(destPath string, params ChromePreviewOutputParams) error {
	content, err := formatPreviewContent(params)
	if err != nil {
		return err
	}

	return writeContentToFile(destPath, content)
}

func formatPreviewContent(params ChromePreviewOutputParams) (string, error) {
	raw, err := json.MarshalIndent(params.Candidates, "", "  ")
	if err != nil {
		return "", apperror.WrapSimple(err, "marshal candidates")
	}

	return string(raw), nil
}
