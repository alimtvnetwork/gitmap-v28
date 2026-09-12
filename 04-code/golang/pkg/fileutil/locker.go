package fileutil

import (
	"path/filepath"
	"sync"
)

type lockEntry struct {
	mu       *sync.RWMutex
	refCount int
}

var (
	// fileLocksMap stores the read/write mutexes and their reference counts mapped by absolute path.
	fileLocksMap = make(map[string]*lockEntry)
	// fileLocksMu protects the map itself from concurrent mutation.
	fileLocksMu sync.Mutex
)

func getAbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}

	return absPath
}

// GetFileLock retrieves or initializes a sync.RWMutex for the specified file path,
// incrementing its reference count. You MUST call ReleaseFileLock when done.
func GetFileLock(path string) *sync.RWMutex {
	absPath := getAbsPath(path)

	fileLocksMu.Lock()
	defer fileLocksMu.Unlock()

	if entry, exists := fileLocksMap[absPath]; exists {
		entry.refCount++

		return entry.mu
	}

	mu := &sync.RWMutex{}
	fileLocksMap[absPath] = &lockEntry{
		mu:       mu,
		refCount: 1,
	}

	return mu
}

// ReleaseFileLock decrements the reference count for a file's lock.
// If the count reaches 0, the lock is evicted from memory to prevent leaks.
func ReleaseFileLock(path string) {
	absPath := getAbsPath(path)

	fileLocksMu.Lock()
	defer fileLocksMu.Unlock()

	if entry, exists := fileLocksMap[absPath]; exists {
		entry.refCount--
		if entry.refCount <= 0 {
			delete(fileLocksMap, absPath)
		}
	}
}
