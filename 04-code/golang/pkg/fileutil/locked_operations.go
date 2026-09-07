package fileutil

import (
	"coding-guidelines/common/pkg/result"
)

// ReadTextLocked acquires a read lock for the file before executing ReadText.
func ReadTextLocked(path string) result.Wrap[string] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.RLock()
	defer mu.RUnlock()

	return ReadText(path)
}

// ReadLinesLocked acquires a read lock for the file before executing ReadLines.
func ReadLinesLocked(path string) result.Wrap[[]string] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.RLock()
	defer mu.RUnlock()

	return ReadLines(path)
}

// ReadJSONLocked acquires a read lock for the file before executing ReadJSON.
func ReadJSONLocked[T any](path string) result.Wrap[T] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.RLock()
	defer mu.RUnlock()

	return ReadJSON[T](path)
}

// ReadYAMLLocked acquires a read lock for the file before executing ReadYAML.
func ReadYAMLLocked[T any](path string) result.Wrap[T] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.RLock()
	defer mu.RUnlock()

	return ReadYAML[T](path)
}

// ExportTextLocked acquires an exclusive write lock for the file before executing ExportText.
func ExportTextLocked(path string, content string, perm FilePermType) result.Wrap[bool] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.Lock()
	defer mu.Unlock()

	return ExportText(path, content, perm)
}

// ExportLinesLocked acquires an exclusive write lock for the file before executing ExportLines.
func ExportLinesLocked(path string, lines []string, perm FilePermType) result.Wrap[bool] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.Lock()
	defer mu.Unlock()

	return ExportLines(path, lines, perm)
}

// ExportJSONLocked acquires an exclusive write lock for the file before executing ExportJSON.
func ExportJSONLocked(path string, data any, perm FilePermType) result.Wrap[bool] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.Lock()
	defer mu.Unlock()

	return ExportJSON(path, data, perm)
}

// ExportYAMLLocked acquires an exclusive write lock for the file before executing ExportYAML.
func ExportYAMLLocked(path string, data any, perm FilePermType) result.Wrap[bool] {
	mu := GetFileLock(path)
	defer ReleaseFileLock(path)
	mu.Lock()
	defer mu.Unlock()

	return ExportYAML(path, data, perm)
}
