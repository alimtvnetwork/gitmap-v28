package cmdagy

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/atotto/clipboard"
)

const activeAgyPromptRelativePath = ".ai-memory/temp/active-agy-pipeline-fix-prompt.txt"

func resolveProjectRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return resolveDefaultGitmapRepoRoot()
	}
	root, rErr := gitutil.RepoRoot(cwd)
	if rErr == nil && len(root) > 0 {
		return root
	}

	return resolveDefaultGitmapRepoRoot()
}

func resolveActiveAgyPromptPath() string {
	return filepath.Join(resolveProjectRootDir(), activeAgyPromptRelativePath)
}

func saveCustomFileIfRequested(customFile, payload string) {
	if len(customFile) > 0 {
		_ = os.WriteFile(customFile, []byte(payload), 0644)
	}
}

func copyClipboardIfNotSkipped(payload string, skipClipboard bool) {
	if !skipClipboard {
		_ = clipboard.WriteAll(payload)
	}
}

func persistFixPromptPayload(payload, customFile string, skipClipboard bool) error {
	tempPath := resolveActiveAgyPromptPath()
	_ = os.MkdirAll(filepath.Dir(tempPath), 0755)
	if writeErr := os.WriteFile(tempPath, []byte(payload), 0644); writeErr != nil {
		return writeErr
	}

	saveCustomFileIfRequested(customFile, payload)
	copyClipboardIfNotSkipped(payload, skipClipboard)

	return nil
}
