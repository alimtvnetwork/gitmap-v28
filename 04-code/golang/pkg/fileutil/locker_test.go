package fileutil

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestLockerMechanisms(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "locked.txt")

	var wg sync.WaitGroup
	// Concurrently write and read using locked operations
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			ExportTextLocked(path, "locked content", FilePermStandard)
		}()
		go func() {
			defer wg.Done()
			res := ReadTextLocked(path)
			if res.IsSuccess() && res.Data() != "locked content" {
				t.Errorf("Unexpected content: %s", res.Data())
			}
		}()
	}
	wg.Wait()

	// Just ensuring it didn't panic or crash due to concurrent access
}
