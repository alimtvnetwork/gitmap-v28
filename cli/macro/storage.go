package macro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func getMacroDir() (string, error) {
	return resolveWritableMacroDir()
}

// SaveMacro writes a macro to disk atomically.
func SaveMacro(m *Macro) error {
	dir, err := resolveWritableMacroDir()
	if err != nil {
		return err
	}

	m.UpdatedAt = time.Now()
	m.TotalSteps = len(m.Steps)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return writeMacroFileAtomic(dir, m.Name, data)
}

func writeMacroFileAtomic(dir, name string, data []byte) error {
	targetPath := filepath.Join(dir, name+".json")
	tmpPath := fmt.Sprintf("%s.%d.tmp", targetPath, time.Now().UnixNano())
	if err := writeAndSyncMacroFile(tmpPath, data); err != nil {
		_ = os.Remove(tmpPath)

		return err
	}

	return finishMacroFileWrite(tmpPath, targetPath)
}

func finishMacroFileWrite(tmpPath, targetPath string) error {
	if err := replaceMacroFile(tmpPath, targetPath); err != nil {
		return err
	}
	restoreSudoOwnership(targetPath)

	return nil
}

func writeAndSyncMacroFile(tmpPath string, data []byte) error {
	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, writeErr := file.Write(data); writeErr != nil {
		return writeErr
	}

	return file.Sync()
}

func replaceMacroFile(tmpPath, targetPath string) error {
	_ = os.Remove(targetPath)
	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)

		return err
	}

	return nil
}

// LoadMacro loads a named macro from candidate directories.
func LoadMacro(name string) (*Macro, error) {
	filename := strings.TrimSuffix(name, ".json") + ".json"
	for _, dir := range candidateMacroDirs() {
		path := filepath.Join(dir, filename)
		m, isFound := readMacroIfExists(path)
		if isFound {
			return m, nil
		}
	}

	return nil, fmt.Errorf("macro %q not found", name)
}

func readMacroIfExists(path string) (*Macro, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var m Macro
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, false
	}

	return &m, true
}

// ListMacros returns all saved macros across candidate directories.
func ListMacros() MacroSliceResult {
	seen := make(map[string]bool)
	var out []Macro
	for _, dir := range candidateMacroDirs() {
		collectMacrosFromDir(dir, seen, &out)
	}

	return result.OkSlice(out)
}

func collectMacrosFromDir(dir string, seen map[string]bool, out *[]Macro) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		appendMacroIfNew(dir, e, seen, out)
	}
}

func appendMacroIfNew(dir string, e os.DirEntry, seen map[string]bool, out *[]Macro) {
	if !strings.HasSuffix(e.Name(), ".json") {
		return
	}
	name := strings.TrimSuffix(e.Name(), ".json")
	_, isSeen := seen[name]
	if isSeen {
		return
	}
	m, isFound := readMacroIfExists(filepath.Join(dir, e.Name()))
	if isFound && m != nil {
		seen[name] = true
		*out = append(*out, *m)
	}
}

// DeleteMacro removes a saved macro across candidate directories.
func DeleteMacro(name string) error {
	filename := strings.TrimSuffix(name, ".json") + ".json"
	var lastErr error
	for _, dir := range candidateMacroDirs() {
		path := filepath.Join(dir, filename)
		if err := removeMacroFileIfExists(path); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func removeMacroFileIfExists(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// MacroExists reports whether a named macro file exists in any candidate directory.
func MacroExists(name string) bool {
	filename := strings.TrimSuffix(name, ".json") + ".json"
	for _, dir := range candidateMacroDirs() {
		path := filepath.Join(dir, filename)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return true
		}
	}

	return false
}

