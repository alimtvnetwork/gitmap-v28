package cmdai

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunAiCreate scaffolds and registers a new AI automation script.
func RunAiCreate(opts CreateScriptOptions) *apperror.AppError {
	valErr := validateCreateOptions(opts)
	if valErr != nil {
		return valErr
	}

	dir := resolveScriptsDir()
	fname := resolveDestinationFilename(opts.Name, dir)
	fullPath := filepath.Join(dir, fname)

	checkErr := verifyTargetFile(fullPath, opts)
	if checkErr != nil {
		return checkErr
	}

	return executeScriptCreation(fullPath, fname, opts)
}

func validateCreateOptions(opts CreateScriptOptions) *apperror.AppError {
	isEmpty := len(opts.Name) == 0
	if isEmpty {
		return apperror.NewValidationError("script name is required (e.g. 'my-linter')")
	}

	return nil
}

func resolveScriptsDir() string {
	dir := findAiScriptsDir()
	hasDir := dir != ""
	if hasDir {
		return dir
	}

	return "03-ai-scripts"
}

func verifyTargetFile(fullPath string, opts CreateScriptOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}

	_, err := os.Stat(fullPath)
	exists := err == nil
	if exists && !opts.IsForce {
		ctx := map[string]any{"path": fullPath}
		return apperror.New("create_script", "E_SCRIPT_EXISTS", ctx)
	}

	return nil
}

func executeScriptCreation(fullPath, fname string, opts CreateScriptOptions) *apperror.AppError {
	res := GenerateScriptContent(opts)
	if res.IsFailed() {
		return res.Err
	}

	content := res.Value
	if opts.IsDryRun {
		renderDryRunPreview(fname, content)
		return nil
	}

	writeErr := writeScriptToDisk(fullPath, content)
	if writeErr != nil {
		return writeErr
	}

	renderCreationSummary(fname, opts)

	return nil
}

func writeScriptToDisk(path string, content string) *apperror.AppError {
	err := os.WriteFile(path, []byte(content), 0755)
	hasErr := err != nil
	if hasErr {
		ctx := map[string]any{"path": path, "err": err.Error()}
		return apperror.New("write_script", "E_WRITE_FAILED", ctx)
	}

	return nil
}

func renderDryRunPreview(fname, content string) {
	fmt.Printf("\n%s [DRY RUN] Generated Script: %s%s\n", constants.ColorCyan, fname, constants.ColorReset)
	fmt.Println(constants.ColorDim + "----------------------------------------" + constants.ColorReset)
	fmt.Println(content)
	fmt.Println(constants.ColorDim + "----------------------------------------" + constants.ColorReset)
}

func renderCreationSummary(fname string, opts CreateScriptOptions) {
	slug := sanitizeScriptSlug(opts.Name)
	fmt.Printf("\n%s✔ Successfully created AI script:%s %s\n", constants.ColorGreen, constants.ColorReset, fname)
	fmt.Printf("  %sType:%s     %s\n", constants.ColorDim, constants.ColorReset, opts.Type)
	fmt.Printf("  %sExecute:%s  gitmap ai run %s\n", constants.ColorBold, constants.ColorReset, slug)
	fmt.Printf("  %sDirect:%s   python 03-ai-scripts/%s\n\n", constants.ColorDim, constants.ColorReset, fname)
}
