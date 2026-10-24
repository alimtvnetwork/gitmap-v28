package cmdconfig

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
)

// loadConfigBundle reads and parses a JSON config bundle file.
func loadConfigBundle(filePath string) (*ConfigBundle, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.Wrap(err, "loadConfigBundle.ReadFile", map[string]any{"path": filePath})
	}

	var bundle ConfigBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, apperror.Wrap(err, "loadConfigBundle.Unmarshal", map[string]any{"path": filePath})
	}

	return &bundle, nil
}

// resolveTargetFileName maps filenames across OS differences (e.g. .ini vs .conf).
func resolveTargetFileName(tool, origName string) string {
	toolType, _ := cfgenum.ParseConfigToolType(tool)
	if !toolType.IsQTorrent() {
		return origName
	}

	if runtime.GOOS == cfgconst.OSWindows {
		return cfgconst.FileQBittorrentIni
	}

	return cfgconst.FileQBittorrentConf
}

// decodeConfigFileBytes converts payload content to raw bytes.
func decodeConfigFileBytes(payload ConfigFilePayload) ([]byte, error) {
	if payload.Encoding != cfgconst.EncodingBase64 {
		return []byte(payload.Content), nil
	}

	data, err := base64.StdEncoding.DecodeString(payload.Content)
	if err != nil {
		return nil, apperror.Wrap(err, "decodeConfigFileBytes.base64", map[string]any{"file": payload.Name})
	}

	return data, nil
}

// writeSingleRestoredFile writes a decoded file to destination directory.
func writeSingleRestoredFile(destDir, fileName string, content []byte) error {
	destPath := filepath.Join(destDir, fileName)
	err := os.WriteFile(destPath, content, 0o644)
	if err != nil {
		return apperror.Wrap(err, "writeSingleRestoredFile", map[string]any{"dest": destPath})
	}

	return nil
}

// restoreFilePayloadEntries iterates and writes individual file payloads.
func restoreFilePayloadEntries(tool, destDir string, files map[string]ConfigFilePayload) (int, error) {
	count := 0
	for origName, payload := range files {
		if err := restoreSingleFilePayload(tool, destDir, origName, payload); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

// restoreBundleFiles writes all files in the bundle to the target directory.
func restoreBundleFiles(bundle *ConfigBundle, destDir string) (int, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return 0, apperror.Wrap(err, "restoreBundleFiles.MkdirAll", map[string]any{"dest": destDir})
	}

	return restoreFilePayloadEntries(bundle.Tool, destDir, bundle.Files)
}

// restoreSingleFilePayload decodes and restores one file entry.
func restoreSingleFilePayload(tool, destDir, origName string, payload ConfigFilePayload) error {
	targetName := resolveTargetFileName(tool, origName)
	data, decodeErr := decodeConfigFileBytes(payload)
	if decodeErr != nil {
		return decodeErr
	}

	return writeSingleRestoredFile(destDir, targetName, data)
}

// writeExtensionsReference writes extensions.txt if VS Code extensions are in the bundle.
func writeExtensionsReference(destDir string, extensions []string) {
	if len(extensions) == 0 {
		return
	}

	extPath := filepath.Join(destDir, cfgconst.FileExtensionsTxt)
	content := strings.Join(extensions, "\n") + "\n"
	_ = os.WriteFile(extPath, []byte(content), 0o644)
}
