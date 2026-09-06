package fileutil

import (
	"path/filepath"
	"sync"
)

var (
	// fileLocksMap stores the read/write mutexes mapped by their absolute paths.
	fileLocksMap = make(map[string]*sync.RWMutex)
	// fileLocksMu protects the map itself from concurrent mutation.
	fileLocksMu sync.Mutex
)

// GetFileLock retrieves or initializes a sync.RWMutex for the specified file path.
// The path is automatically cleaned/absolved for consistency.
func GetFileLock(path string) *sync.RWMutex {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = filepath.Clean(path)
	}

	fileLocksMu.Lock()
	defer fileLocksMu.Unlock()

	if mu, exists := fileLocksMap[absPath]; exists {
		return mu
	}

	mu := &sync.RWMutex{}
	fileLocksMap[absPath] = mu
	return mu
}
