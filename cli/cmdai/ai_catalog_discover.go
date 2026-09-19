package cmdai

import (
	"os"
	"path/filepath"
	"strings"
)

// GetAllScripts returns master catalog scripts merged with dynamic scripts from 03-ai-scripts.
func GetAllScripts() []ScriptMetadata {
	base := append([]ScriptMetadata{}, masterCatalog...)
	dir := findAiScriptsDir()
	isDir := dir != ""
	if !isDir {
		return base
	}

	return mergeDynamicScripts(base, dir)
}

func findAiScriptsDir() string {
	candidates := []string{"03-ai-scripts", filepath.Join("..", "03-ai-scripts")}
	for _, c := range candidates {
		info, err := os.Stat(c)
		isMatch := err == nil && info.IsDir()
		if isMatch {
			return c
		}
	}

	return ""
}

func mergeDynamicScripts(base []ScriptMetadata, dir string) []ScriptMetadata {
	entries, err := os.ReadDir(dir)
	hasErr := err != nil
	if hasErr {
		return base
	}

	known := makeKnownMap(base)
	for _, entry := range entries {
		base = appendIfNewScript(base, entry, dir, known)
	}

	return base
}

func makeKnownMap(scripts []ScriptMetadata) map[string]bool {
	known := make(map[string]bool)
	for _, s := range scripts {
		known[s.Filename] = true
	}

	return known
}

func appendIfNewScript(base []ScriptMetadata, entry os.DirEntry, dir string, known map[string]bool) []ScriptMetadata {
	isPy := !entry.IsDir() && strings.HasSuffix(entry.Name(), ".py")
	isUnknown := isPy && !known[entry.Name()]
	if !isUnknown {
		return base
	}

	meta := buildDynamicMetadata(entry.Name(), dir)
	known[entry.Name()] = true

	return append(base, meta)
}
