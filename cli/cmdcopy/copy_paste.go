// Package cmdcopy — copy_paste.go: cross-platform clipboard and memory buffer copy-paste operations.
package cmdcopy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/enums"
	"github.com/atotto/clipboard"
)

func RunCopyCmd(args []string) error {
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

	return apperror.NewValidationError(constants.ErrMsgNothingToCopy)
}

func persistAndAnnounceCopy(content string) error {
	if err := saveContentToMemoryAndClipboard(content); err != nil {
		return err
	}

	fmt.Printf(constants.MsgCopiedSuccessFmt,
		constants.ColorGreen, len(content), constants.ColorReset)

	return nil
}

func printCopyUsage() {
	fmt.Println(constants.HelpCopyUsage)
	fmt.Println()
	fmt.Println(constants.HelpCopyDesc)
	fmt.Println()
}

func determineCopySource(args []string) enums.CopySourceType {
	if extractCopyFilePath(args) != "" {
		return enums.CopySourceFile
	}

	filtered := filterFlagArgs(args)
	if len(filtered) == 1 && isExistingFile(filtered[0]) {
		return enums.CopySourceFile
	}

	if len(filtered) > 0 {
		return enums.CopySourceArgs
	}

	return enums.CopySourceStdin
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
		if (a == constants.FlagFile || a == constants.FlagFileShort) && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(a, constants.FlagPrefixFile) {
			return strings.TrimPrefix(a, constants.FlagPrefixFile)
		}
	}

	return ""
}

func filterFlagArgs(args []string) []string {
	var out []string
	for i := 0; i < len(args); i++ {
		if (args[i] == constants.FlagFile || args[i] == constants.FlagFileShort) && i+1 < len(args) {
			i++
			continue
		}

		if !strings.HasPrefix(args[i], "-") {
			out = append(out, args[i])
		}
	}

	return out
}

func readFileContentForCopy(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf(constants.ErrFmtReadFile, path))
	}

	return string(data), nil
}

func readPipedStdinForCopy() (string, error) {
	if isTerminalInput() {
		return "", nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", apperror.WrapSimple(err, constants.ErrFmtReadStdin)
	}

	return string(data), nil
}

func saveContentToMemoryAndClipboard(content string) error {
	_ = clipboard.WriteAll(content)
	memPath, err := resolveMemoryFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(memPath, []byte(content), constants.DefaultFilePerm); err != nil {
		return apperror.WrapSimple(err, constants.ErrFmtWriteMemoryClipboard)
	}

	return nil
}

func resolveMemoryFilePath() (string, error) {
	dir := resolveMemoryBaseDir()
	if err := os.MkdirAll(dir, constants.DefaultDirPerm); err != nil {
		return "", apperror.WrapSimple(err, constants.ErrFmtCreateMemoryDir)
	}

	return filepath.Join(dir, constants.FileClipboard), nil
}

func resolveMemoryBaseDir() string {
	if root, err := gitTopLevel(); err == nil && len(root) > 0 {
		return filepath.Join(root, constants.DirGitmap, constants.DirMemory)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(constants.DotDirSentinel, constants.DirGitmap, constants.DirMemory)
	}

	return filepath.Join(home, constants.DirGitmap, constants.DirMemory)
}

func RunPasteCmd(args []string) error {
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
	fmt.Printf(constants.MsgEmptyPasteNoticeFmt, constants.ColorYellow, constants.ColorReset)
}

func determinePasteTarget(outFile string) enums.CopyTargetType {
	if outFile != "" {
		return enums.CopyTargetFile
	}

	return enums.CopyTargetStdout
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

		if strings.HasPrefix(a, constants.FlagPrefixFile) || strings.HasPrefix(a, constants.FlagPrefixOut) {
			return strings.TrimPrefix(strings.TrimPrefix(a, constants.FlagPrefixFile), constants.FlagPrefixOut)
		}
	}

	return ""
}

func writePastedContentToFile(outFile, content string) error {
	dir := filepath.Dir(outFile)
	if dir != "" && dir != constants.DotDirSentinel {
		_ = os.MkdirAll(dir, constants.DefaultDirPerm)
	}

	if err := os.WriteFile(outFile, []byte(content), constants.DefaultFilePerm); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrFmtWritePastedContent, outFile))
	}

	fmt.Printf(constants.MsgPastedSuccessFmt,
		constants.ColorGreen, len(content), outFile, constants.ColorReset)

	return nil
}

func fetchClipboardOrMemoryContent() (string, error) {
	clipText, err := clipboard.ReadAll()
	if err == nil && len(clipText) > 0 {
		return clipText, nil
	}

	return readMemoryFallbackContent()
}

func readMemoryFallbackContent() (string, error) {
	memPath, err := resolveMemoryFilePath()
	if err != nil {
		return "", err
	}

	data, readErr := os.ReadFile(memPath)
	if readErr == nil {
		return string(data), nil
	}

	return "", nil
}
