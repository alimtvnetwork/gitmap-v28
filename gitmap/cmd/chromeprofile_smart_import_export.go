package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
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
	return arg == "--file" || arg == "-f" || arg == "-o" || arg == "--output" || arg == "--out"
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
	return arg == "--tempfile" || arg == "--temp"
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
	if arg == "--json" || arg == "-j" {
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
		return "", fmt.Errorf("marshal candidates: %w", err)
	}

	return string(raw), nil
}
