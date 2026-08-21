package macro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

// SaveMacro writes a macro to disk.
func SaveMacro(m *Macro) error {
	dir, err := getMacroDir()
	if err != nil {
		return err
	}
	m.UpdatedAt = time.Now()
	m.TotalSteps = len(m.Steps)
	path := filepath.Join(dir, m.Name+".json")
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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
		if strings.HasSuffix(e.Name(), ".json") {
			name := strings.TrimSuffix(e.Name(), ".json")
			m, err := LoadMacro(name)
			if err == nil && m != nil {
				out = append(out, *m)
			}
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
	path := filepath.Join(dir, strings.TrimSuffix(name, ".json")+".json")
	return os.Remove(path)
}
