package macro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func getMacroDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, constants.GitMapDir, "macros")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

// SaveMacro writes a macro to disk atomically.
func SaveMacro(m *Macro) error {
	dir, err := getMacroDir()
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

	return replaceMacroFile(tmpPath, targetPath)
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

// LoadMacro loads a named macro from disk.
func LoadMacro(name string) (*Macro, error) {
	dir, err := getMacroDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, strings.TrimSuffix(name, ".json")+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("macro %q not found: %w", name, err)
	}

	var m Macro
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// ListMacros returns all saved macros.
func ListMacros() ([]Macro, error) {
	dir, err := getMacroDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []Macro
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".json")
		m, err := LoadMacro(name)
		if err == nil && m != nil {
			out = append(out, *m)
		}
	}

	return out, nil
}

// DeleteMacro removes a saved macro.
func DeleteMacro(name string) error {
	dir, err := getMacroDir()
	if err != nil {
		return err
	}

	cleanName := strings.TrimSpace(strings.TrimSuffix(name, ".json"))
	path := filepath.Join(dir, cleanName+".json")
	if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
		return removeErr
	}

	return nil
}

// MacroExists reports whether a named macro file exists on disk.
func MacroExists(name string) bool {
	dir, err := getMacroDir()
	if err != nil {
		return false
	}

	cleanName := strings.TrimSpace(strings.TrimSuffix(name, ".json"))
	path := filepath.Join(dir, cleanName+".json")
	info, statErr := os.Stat(path)

	return statErr == nil && !info.IsDir()
}
