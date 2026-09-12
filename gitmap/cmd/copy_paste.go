// Package cmd — copy_paste.go: cross-platform clipboard and memory buffer copy-paste operations.
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/atotto/clipboard"
)

func runCopyCmd(args []string) error {
	content, err := resolveCopyContent(args)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return reportEmptyCopy()
	}

	return persistAndAnnounceCopy(content)
}

func reportEmptyCopy() error {
	printCopyUsage()

	return apperror.NewValidationError("nothing to copy: text or file required")
}

func persistAndAnnounceCopy(content string) error {
	if err := saveContentToMemoryAndClipboard(content); err != nil {
		return err
	}

	fmt.Printf("  %s📋 Copied %d bytes to memory & clipboard ✅%s\n",
		constants.ColorGreen, len(content), constants.ColorReset)

	return nil
}

func printCopyUsage() {
	fmt.Println("Usage: gitmap copy <text...> [--file <path>]")
	fmt.Println()
	fmt.Println("Copies text, items, or file contents to OS clipboard and persistent repository memory.")
	fmt.Println()
}

func resolveCopyContent(args []string) (string, error) {
	filePath := extractCopyFilePath(args)
	if filePath != "" {
		return readFileContentForCopy(filePath)
	}

	filtered := filterFlagArgs(args)
	if len(filtered) == 1 && isExistingFile(filtered[0]) {
		return readFileContentForCopy(filtered[0])
	}

	return resolveCopyTextOrStdin(filtered)
}

func resolveCopyTextOrStdin(filtered []string) (string, error) {
	if len(filtered) > 0 {
		return strings.Join(filtered, " "), nil
	}

	return readPipedStdinForCopy()
}

func extractCopyFilePath(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--file" || a == "-f") && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(a, "--file=") {
			return strings.TrimPrefix(a, "--file=")
		}
	}

	return ""
}

func filterFlagArgs(args []string) []string {
	var out []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--file" || a == "-f") && i+1 < len(args) {
			i++
			continue
		}

		if !strings.HasPrefix(a, "-") {
			out = append(out, a)
		}
	}

	return out
}

func readFileContentForCopy(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("read file %q for copy", path))
	}

	return string(data), nil
}

func readPipedStdinForCopy() (string, error) {
	if isTerminalInput() {
		return "", nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", apperror.WrapSimple(err, "read stdin for copy")
	}

	return string(data), nil
}

func saveContentToMemoryAndClipboard(content string) error {
	_ = clipboard.WriteAll(content)
	memPath, err := resolveMemoryFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(memPath, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "write memory clipboard buffer")
	}

	return nil
}

func resolveMemoryFilePath() (string, error) {
	dir := resolveMemoryBaseDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", apperror.WrapSimple(err, "create memory directory")
	}

	return filepath.Join(dir, "clipboard.txt"), nil
}

func resolveMemoryBaseDir() string {
	if root, err := gitTopLevel(); err == nil && len(root) > 0 {
		return filepath.Join(root, ".gitmap", "memory")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".gitmap", "memory")
	}

	return filepath.Join(home, ".gitmap", "memory")
}

func runPasteCmd(args []string) error {
	content, err := fetchClipboardOrMemoryContent()
	if err != nil {
		return err
	}

	if len(content) == 0 {
		printEmptyPasteNotice()

		return nil
	}

	return dispatchPasteOutput(args, content)
}

func printEmptyPasteNotice() {
	fmt.Printf("  %s⚠️  Clipboard and memory buffer are empty.%s\n", constants.ColorYellow, constants.ColorReset)
}

func dispatchPasteOutput(args []string, content string) error {
	outFile := extractPasteOutFile(args)
	if outFile != "" {
		return writePastedContentToFile(outFile, content)
	}

	printPastedContent(content)

	return nil
}

func printPastedContent(content string) {
	fmt.Print(content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Println()
	}
}

func extractPasteOutFile(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if isFileFlagWithArg(a) && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(a, "--file=") || strings.HasPrefix(a, "--out=") {
			return strings.TrimPrefix(strings.TrimPrefix(a, "--file="), "--out=")
		}
	}

	return ""
}

func writePastedContentToFile(outFile, content string) error {
	dir := filepath.Dir(outFile)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}

	if err := os.WriteFile(outFile, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("write pasted content to %q", outFile))
	}

	fmt.Printf("  %s📋 Pasted %d bytes to %s ✅%s\n",
		constants.ColorGreen, len(content), outFile, constants.ColorReset)

	return nil
}

func fetchClipboardOrMemoryContent() (string, error) {
	clipText, err := clipboard.ReadAll()
	if err == nil && len(clipText) > 0 {
		return clipText, nil
	}

	memPath, memErr := resolveMemoryFilePath()
	if memErr != nil {
		return "", memErr
	}

	data, readErr := os.ReadFile(memPath)
	if readErr == nil {
		return string(data), nil
	}

	return "", nil
}
